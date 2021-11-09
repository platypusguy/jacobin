/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"jacobin/exec"
	"jacobin/globals"
	"jacobin/log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// Classloader holds the parsed bytecode in classes, where they can be retrieved
// and moved to an execution role. Most of the comments and code presuppose some
// familiarity with the role of classloaders. More information can be found at:
// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-5.html#jvms-5.3
type Classloader struct {
	Name    string
	Parent  string
	Classes map[string]ParsedClass
}

// AppCL is the application classloader, which loads most of the app's classes
var AppCL Classloader

// BootstrapCL is the classloader that loads most of the standard libraries
var BootstrapCL Classloader

// ExtensionCL is the classloader typically used for loading custom agents
var ExtensionCL Classloader

// the parsed class
type ParsedClass struct {
	javaVersion    int
	className      string // name of class without path and without .class
	superClass     string // name of superclass for this class
	moduleName     string
	packageName    string
	interfaceCount int   // number of interfaces this class implements
	interfaces     []int // the interfaces this class implements, as indices into utf8Refs
	fieldCount     int   // number of fields in this class
	fields         []field
	methodCount    int
	methods        []method
	attribCount    int
	attributes     []attr
	sourceFile     string
	bootstrapCount int // the number of bootstrap methods
	bootstraps     []bootstrapMethod

	deprecated bool

	// ---- constant pool data items ----
	cpCount        int       // count of constant pool entries
	cpIndex        []cpEntry // the constant pool index to entries
	classRefs      []int     // points to a UTF-8 entry in the CP
	doubles        []float64
	dynamics       []dynamic
	fieldRefs      []fieldRefEntry
	floats         []float32
	intConsts      []int // 32-bit int containing the actual int value
	interfaceRefs  []interfaceRefEntry
	invokeDynamics []invokeDynamic
	longConsts     []int64
	methodHandles  []methodHandleEntry
	methodRefs     []methodRefEntry
	methodTypes    []int
	nameAndTypes   []nameAndTypeEntry
	stringRefs     []stringConstantEntry // integer index into utf8Refs
	utf8Refs       []utf8Entry

	// ---- access flags items ----
	accessFlags       int // the following booleans interpret the access flags
	classIsPublic     bool
	classIsFinal      bool
	classIsSuper      bool
	classIsInterface  bool
	classIsAbstract   bool
	classIsSynthetic  bool
	classIsAnnotation bool
	classIsEnum       bool
	classIsModule     bool

	// ---- field attributes ----
}

// the fields defined in the class
type field struct {
	accessFlags int
	name        int // index of the UTF-8 entry in the CP
	description int // index of the UTF-8 entry in the CP
	attributes  []attr
}

// the methods of the class, including the constructors
type method struct {
	accessFlags int
	name        int // index of the UTF-8 entry in the CP
	description int // index of the UTF-8 entry in the CP
	codeAttr    codeAttrib
	attributes  []attr
	exceptions  []int // indexes into Utf8Refs in the CP
	parameters  []paramAttrib
	deprecated  bool // is the method deprecated?
}

type codeAttrib struct {
	maxStack   int
	maxLocals  int
	code       []byte
	exceptions []exception // exception entries for this method
	attributes []attr      // the code attributes has its own sub-attributes(!)
}

// the MethodParameters method attribute
type paramAttrib struct {
	name        string // string, rather than index into utf8Refs b/c the name could be ""
	accessFlags int
}

// the structure of many attributes (field, class, etc.) The content is just the raw bytes.
type attr struct {
	attrName    int    // index of the UTF-8 entry in the CP
	attrSize    int    // length of the following array of raw bytes
	attrContent []byte // the raw data of the attribute
}

// the exception-related data for each exception in the Code attribute of a given method
type exception struct {
	startPc   int // first instruction covered by this exception (pc = program counter)
	endPc     int // the last instruction covered by this exception
	handlerPc int // the place in the method code that has the exception instructions
	catchType int // the type of exception, index to CP, which must point a ClassFref entry
}

// the boostrap methods, specified in the bootstrap class attribute
type bootstrapMethod struct {
	methodRef int   // index pointing to a MethodHandle
	args      []int // arguments: indexes to loadable arguments from the CP
}

// var lock = sync.RWMutex{}

// cfe = class format error, which is the error thrown by the parser for most
// of the errors arising from malformed bytecode. Prints out file and line# where
// the call to cfe() occurred.
func cfe(msg string) error {
	errMsg := "Class Format Error: " + msg

	// get the filename and line# of the function where the error occurred
	// implementation note: Caller(0) would be this function. (1) is the
	// previous function on the stack (so, the one calling this error routine)
	// To traverse all the way back to the start of the program, set up a loop
	// and exit when ok is no longer true.
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		fn := runtime.FuncForPC(pc)
		fileName, fileLine := fn.FileLine(pc)
		errMsg = errMsg + "\n  dectected by file: " + filepath.Base(fileName) +
			", line: " + strconv.Itoa(fileLine)
	}
	log.Log(errMsg, log.SEVERE)
	return errors.New(errMsg)
}

// LoadBaseClasses loads a basic set of classes that are specified in the file
// classes\baseclasslist.txt, which is found in JACOBIN_HOME. It's similar to
// classlist file in the JDK, except shorter (for the nonce)
func LoadBaseClasses(global *globals.Globals) {
	classList := global.JacobinHome + "classes\\baseclasslist.txt"
	file, err := os.Open(classList)
	if err != nil {
		log.Log("Did not find baseclasslist.txt in JACOBIN_HOME ("+classList+")",
			log.WARNING)
		file.Close()
	} else {
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			rawName := scanner.Text()
			name := exec.ConvertInternalClassNameToFilename(rawName)
			name = globals.JacobinHome() + "classes\\" + name
			LoadClassFromFile(BootstrapCL, name)
			// LoadReferencedClasses(BootstrapCL, rawName)
		}
		err = nil // used only to be able to add a breakpoint in debugger.
	}
}

// This loads the classes referenced in the loading of the class named clName.
// It does this by reading the class entries (7) in the CP and sending the class names
// it finds there to a go channel that will load the class.
func LoadReferencedClasses(classloader Classloader, clName string) {
	cpClassCP := &exec.Classes[clName].Data.CP
	classRefs := cpClassCP.ClassRefs

	loaderChannel := make(chan string, len(classRefs))
	for _, v := range classRefs {
		refClassName := exec.FetchUTF8stringFromCPEntryNumber(cpClassCP, v)
		name := normalizeClassReference(refClassName)
		if name == "" {
			continue
		}
		loaderChannel <- name
	}
	globals.LoaderWg.Add(1)
	go LoadFromLoaderChannel(loaderChannel)
	close(loaderChannel)
}

// receives a name of a class to load in /java/lang/String format, determines the
// classloader, checks if the class is already loaded, and loads it if not.
func LoadFromLoaderChannel(LoaderChannel <-chan string) {
	for name := range LoaderChannel {
		_, present := exec.Classes[name]
		if present { // if the class is already loaded, skip this.
			continue
		}

		if strings.HasPrefix(name, "java/") || strings.HasPrefix(name, "jdk/") ||
			strings.HasPrefix(name, "sun/") {
			name = exec.ConvertInternalClassNameToFilename(name)
			name = globals.JacobinHome() + "classes\\" + name
			LoadClassFromFile(BootstrapCL, name)
		}
		println("loading from channel: " + name)
	}
	globals.LoaderWg.Done()
}

// LoadClassFromFile first canonicalizes the filename, checks whether
// the class is already loaded, and if not, then parses the class and loads it.
// Returns the class's internal name and error, if any.
// 1 TODO: canonicalize class name
// 2 TODO: search through classloaders for this class
// 3 TODO: determine which classloader should load the class, then
// 4 TODO: have *it* parse and load the class.
func LoadClassFromFile(cl Classloader, filename string) (string, error) {
	rawBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Log("Could not read file: "+filename+". Exiting.", log.SEVERE)
		return "", fmt.Errorf("file I/O error")
	}

	log.Log(filename+" read", log.FINE)

	fullyParsedClass, err := parse(rawBytes)
	if err != nil {
		log.Log("error parsing "+filename+". Exiting.", log.SEVERE)
		return "", fmt.Errorf("parsing error")
	}

	// add entry to the method area, indicating initialization of the load of this class
	eKI := exec.Klass{
		Status: 'I', // I = initializing the load
		Loader: cl.Name,
		Data:   nil,
	}
	insert(fullyParsedClass.className, eKI)

	// format check the class
	if formatCheckClass(&fullyParsedClass) != nil {
		log.Log("error format-checking "+filename+". Exiting.", log.SEVERE)
		return "", fmt.Errorf("format-checking error")
	}
	log.Log("Class "+fullyParsedClass.className+" has been format-checked.", log.FINEST)

	classToPost := convertToPostableClass(&fullyParsedClass)
	eKF := exec.Klass{
		Status: 'F', // F = format-checked
		Loader: cl.Name,
		Data:   &classToPost,
	}
	insert(fullyParsedClass.className, eKF)

	return fullyParsedClass.className, nil
}

// insert the fully parsed class into the method area (exec.Classes)
func insert(name string, klass exec.Klass) error {
	exec.MethAreaMutex.Lock()
	exec.Classes[name] = klass
	exec.MethAreaMutex.Unlock()

	if klass.Status == 'F' || klass.Status == 'V' || klass.Status == 'L' {
		log.Log("Class: "+klass.Data.Name+", loader: "+klass.Loader, log.CLASS)
	}
	return nil
}

// load the parse class into a form suitable for posting to the method area (which is
// exec.Classes. This mostly involves copying the data, converting most indexes to uint16
// and removing some fields we needed in parsing, but which are no longer required.
func convertToPostableClass(fullyParsedClass *ParsedClass) exec.ClData {

	kd := exec.ClData{}
	kd.Name = fullyParsedClass.className
	kd.Superclass = fullyParsedClass.superClass
	kd.Module = fullyParsedClass.moduleName
	kd.Pkg = fullyParsedClass.packageName
	for i := 0; i < len(fullyParsedClass.interfaces); i++ {
		kd.Interfaces = append(kd.Interfaces, uint16(fullyParsedClass.interfaces[i]))
	}
	if len(fullyParsedClass.fields) > 0 {
		for i := 0; i < len(fullyParsedClass.fields); i++ {
			kdf := exec.Field{}
			kdf.Name = uint16(fullyParsedClass.fields[i].name)
			kdf.Desc = uint16(fullyParsedClass.fields[i].description)
			if len(fullyParsedClass.fields[i].attributes) > 0 {
				for j := 0; j < len(fullyParsedClass.fields[i].attributes); j++ {
					kdfa := exec.Attr{}
					kdfa.AttrName = uint16(fullyParsedClass.fields[i].attributes[j].attrName)
					kdfa.AttrSize = fullyParsedClass.fields[i].attributes[j].attrSize
					kdfa.AttrContent = fullyParsedClass.fields[i].attributes[j].attrContent
					kdf.Attributes = append(kdf.Attributes, kdfa)
				}
			}
			kd.Fields = append(kd.Fields, kdf)
		}
	}

	if len(fullyParsedClass.methods) > 0 {
		for i := 0; i < len(fullyParsedClass.methods); i++ {
			kdm := exec.Method{}
			kdm.Name = uint16(fullyParsedClass.methods[i].name)
			kdm.Desc = uint16(fullyParsedClass.methods[i].description)
			kdm.AccessFlags = fullyParsedClass.methods[i].accessFlags
			kdm.CodeAttr.MaxStack = fullyParsedClass.methods[i].codeAttr.maxStack
			kdm.CodeAttr.MaxLocals = fullyParsedClass.methods[i].codeAttr.maxLocals
			kdm.CodeAttr.Code = fullyParsedClass.methods[i].codeAttr.code
			if len(fullyParsedClass.methods[i].codeAttr.exceptions) > 0 {
				for j := 0; j < len(fullyParsedClass.methods[i].codeAttr.exceptions); j++ {
					kdmce := exec.CodeException{}
					kdmce.StartPc = fullyParsedClass.methods[i].codeAttr.exceptions[j].startPc
					kdmce.EndPc = fullyParsedClass.methods[i].codeAttr.exceptions[j].endPc
					kdmce.HandlerPc = fullyParsedClass.methods[i].codeAttr.exceptions[j].handlerPc
					kdmce.CatchType = uint16(fullyParsedClass.methods[i].codeAttr.exceptions[j].catchType)
					kdm.CodeAttr.Exceptions = append(kdm.CodeAttr.Exceptions, kdmce)
				}
			}
			if len(fullyParsedClass.methods[i].codeAttr.attributes) > 0 {
				for m := 0; m < len(fullyParsedClass.methods[i].codeAttr.attributes); m++ {
					kdmca := exec.Attr{}
					kdmca.AttrName = uint16(fullyParsedClass.methods[i].codeAttr.attributes[m].attrName)
					kdmca.AttrSize = fullyParsedClass.methods[i].codeAttr.attributes[m].attrSize
					kdmca.AttrContent = fullyParsedClass.methods[i].codeAttr.attributes[m].attrContent
					kdm.CodeAttr.Attributes = append(kdm.CodeAttr.Attributes, kdmca)
				}
			}
			if len(fullyParsedClass.methods[i].attributes) > 0 {
				for n := 0; n < len(fullyParsedClass.methods[i].attributes); n++ {
					kdma := exec.Attr{
						AttrName:    uint16(fullyParsedClass.methods[i].attributes[n].attrName),
						AttrSize:    fullyParsedClass.methods[i].attributes[n].attrSize,
						AttrContent: fullyParsedClass.methods[i].attributes[n].attrContent,
					}
					kdm.Attributes = append(kdm.Attributes, kdma)
				}
			}
			if len(fullyParsedClass.methods[i].exceptions) > 0 {
				for p := 0; p < len(fullyParsedClass.methods[i].exceptions); p++ {
					kdm.Exceptions = append(kdm.Exceptions, uint16(fullyParsedClass.methods[i].exceptions[p]))
				}
			}
			if len(fullyParsedClass.methods[i].parameters) > 0 {
				for q := 0; q < len(fullyParsedClass.methods[i].parameters); q++ {
					kdmp := exec.ParamAttrib{
						Name:        fullyParsedClass.methods[i].parameters[q].name,
						AccessFlags: fullyParsedClass.methods[i].parameters[q].accessFlags,
					}
					kdm.Parameters = append(kdm.Parameters, kdmp)
				}
			}
			kdm.Deprecated = fullyParsedClass.methods[i].deprecated
			kd.Methods = append(kd.Methods, kdm)
		}
	}
	if len(fullyParsedClass.attributes) > 0 {
		for i := 0; i < len(fullyParsedClass.attributes); i++ {
			kda := exec.Attr{
				AttrName:    uint16(fullyParsedClass.attributes[i].attrName),
				AttrSize:    fullyParsedClass.attributes[i].attrSize,
				AttrContent: fullyParsedClass.attributes[i].attrContent,
			}
			kd.Attributes = append(kd.Attributes, kda)
		}
	}
	kd.SourceFile = fullyParsedClass.sourceFile
	if len(fullyParsedClass.bootstraps) > 0 {
		for j := 0; j < len(fullyParsedClass.bootstraps); j++ {
			kdbs := exec.BootstrapMethod{
				MethodRef: uint16(fullyParsedClass.bootstraps[j].methodRef),
				Args:      nil,
			}
			if len(fullyParsedClass.bootstraps[j].args) > 0 {
				for l := 0; l < len(fullyParsedClass.bootstraps[j].args); l++ {
					kdbs.Args = append(kdbs.Args, (uint16(fullyParsedClass.bootstraps[j].args[l])))
				}
			}
			kd.Bootstraps = append(kd.Bootstraps, kdbs)
		}
	}
	kd.Access.ClassIsPublic = fullyParsedClass.classIsPublic
	kd.Access.ClassIsFinal = fullyParsedClass.classIsFinal
	kd.Access.ClassIsSuper = fullyParsedClass.classIsSuper
	kd.Access.ClassIsInterface = fullyParsedClass.classIsInterface
	kd.Access.ClassIsAbstract = fullyParsedClass.classIsAbstract
	kd.Access.ClassIsSynthetic = fullyParsedClass.classIsSynthetic
	kd.Access.ClassIsAnnotation = fullyParsedClass.classIsAnnotation
	kd.Access.ClassIsEnum = fullyParsedClass.classIsEnum
	kd.Access.ClassIsModule = fullyParsedClass.classIsModule

	// ---- loading the CP ----
	for i := 0; i < fullyParsedClass.cpCount; i++ {
		cpE := exec.CpEntry{
			Type: uint16(fullyParsedClass.cpIndex[i].entryType),
			Slot: uint16(fullyParsedClass.cpIndex[i].slot),
		}
		kd.CP.CpIndex = append(kd.CP.CpIndex, cpE)
	}

	if len(fullyParsedClass.classRefs) > 0 {
		for i := 0; i < len(fullyParsedClass.classRefs); i++ {
			kd.CP.ClassRefs = append(kd.CP.ClassRefs, uint16(fullyParsedClass.classRefs[i]))
		}
	}

	if len(fullyParsedClass.doubles) > 0 {
		for i := 0; i < len(fullyParsedClass.doubles); i++ {
			kd.CP.Doubles = append(kd.CP.Doubles, fullyParsedClass.doubles[i])
		}
	}

	if len(fullyParsedClass.dynamics) > 0 {
		for i := 0; i < len(fullyParsedClass.dynamics); i++ {
			dyn := exec.Dynamic{
				BootstrapIndex: uint16(fullyParsedClass.dynamics[i].bootstrapIndex),
				NameAndType:    uint16(fullyParsedClass.dynamics[i].nameAndType),
			}
			kd.CP.Dynamics = append(kd.CP.Dynamics, dyn)
		}
	}

	if len(fullyParsedClass.fieldRefs) > 0 {
		for i := 0; i < len(fullyParsedClass.fieldRefs); i++ {
			fr := exec.FieldRefEntry{
				ClassIndex:  uint16(fullyParsedClass.fieldRefs[i].classIndex),
				NameAndType: uint16(fullyParsedClass.fieldRefs[i].nameAndTypeIndex),
			}
			kd.CP.FieldRefs = append(kd.CP.FieldRefs, fr)
		}
	}

	if len(fullyParsedClass.floats) > 0 {
		for i := 0; i < len(fullyParsedClass.floats); i++ {
			kd.CP.Floats = append(kd.CP.Floats, fullyParsedClass.floats[i])
		}
	}

	if len(fullyParsedClass.intConsts) > 0 {
		for i := 0; i < len(fullyParsedClass.intConsts); i++ {
			kd.CP.IntConsts = append(kd.CP.IntConsts, int32(fullyParsedClass.intConsts[i]))
		}
	}

	if len(fullyParsedClass.interfaceRefs) > 0 {
		for i := 0; i < len(fullyParsedClass.interfaceRefs); i++ {
			ir := exec.InterfaceRefEntry{
				ClassIndex:  uint16(fullyParsedClass.interfaceRefs[i].classIndex),
				NameAndType: uint16(fullyParsedClass.interfaceRefs[i].nameAndTypeIndex),
			}
			kd.CP.InterfaceRefs = append(kd.CP.InterfaceRefs, ir)
		}
	}

	if len(fullyParsedClass.invokeDynamics) > 0 {
		for i := 0; i < len(fullyParsedClass.invokeDynamics); i++ {
			id := exec.InvokeDynamic{
				BootstrapIndex: uint16(fullyParsedClass.invokeDynamics[i].bootstrapIndex),
				NameAndType:    uint16(fullyParsedClass.invokeDynamics[i].nameAndType),
			}
			kd.CP.InvokeDynamics = append(kd.CP.InvokeDynamics, id)
		}
	}

	if len(fullyParsedClass.longConsts) > 0 {
		for i := 0; i < len(fullyParsedClass.longConsts); i++ {
			kd.CP.LongConsts = append(kd.CP.LongConsts, fullyParsedClass.longConsts[i])
		}
	}

	if len(fullyParsedClass.methodHandles) > 0 {
		for i := 0; i < len(fullyParsedClass.methodHandles); i++ {
			mh := exec.MethodHandleEntry{
				RefKind:  uint16(fullyParsedClass.methodHandles[i].referenceKind),
				RefIndex: uint16(fullyParsedClass.methodHandles[i].referenceIndex),
			}
			kd.CP.MethodHandles = append(kd.CP.MethodHandles, mh)
		}
	}

	if len(fullyParsedClass.methodRefs) > 0 {
		for i := 0; i < len(fullyParsedClass.methodRefs); i++ {
			mr := exec.MethodRefEntry{
				ClassIndex:  uint16(fullyParsedClass.methodRefs[i].classIndex),
				NameAndType: uint16(fullyParsedClass.methodRefs[i].nameAndTypeIndex),
			}
			kd.CP.MethodRefs = append(kd.CP.MethodRefs, mr)
		}
	}

	if len(fullyParsedClass.methodTypes) > 0 {
		for i := 0; i < len(fullyParsedClass.methodTypes); i++ {
			kd.CP.MethodTypes = append(kd.CP.MethodTypes, uint16(fullyParsedClass.methodTypes[i]))
		}
	}

	if len(fullyParsedClass.nameAndTypes) > 0 {
		for i := 0; i < len(fullyParsedClass.nameAndTypes); i++ {
			nat := exec.NameAndTypeEntry{
				NameIndex: uint16(fullyParsedClass.nameAndTypes[i].nameIndex),
				DescIndex: uint16(fullyParsedClass.nameAndTypes[i].descriptorIndex),
			}
			kd.CP.NameAndTypes = append(kd.CP.NameAndTypes, nat)
		}
	}

	if len(fullyParsedClass.stringRefs) > 0 {
		for i := 0; i < len(fullyParsedClass.stringRefs); i++ {
			kd.CP.StringRefs = append(kd.CP.StringRefs, uint16(fullyParsedClass.stringRefs[i].index))
		}
	}

	if len(fullyParsedClass.utf8Refs) > 0 {
		for i := 0; i < len(fullyParsedClass.utf8Refs); i++ {
			kd.CP.Utf8Refs = append(kd.CP.Utf8Refs, fullyParsedClass.utf8Refs[i].content)
		}
	}

	if log.Level == log.FINEST {
		b := new(bytes.Buffer)
		if gob.NewEncoder(b).Encode(kd) == nil {
			log.Log("Size of loaded class: "+strconv.Itoa(b.Len()), log.FINEST)
		}
	}
	return kd
}

// accepts a string containing a class reference from a class file and converts
// it into a normalized z/y/x format. It converts references that start with [L
// and skips all array classes. For these latter cases or any errors, it returns ""
func normalizeClassReference(ref string) string {
	refClassName := ref
	if strings.HasPrefix(refClassName, "[L") {
		refClassName = strings.TrimPrefix(refClassName, "[L")
		if strings.HasSuffix(refClassName, ";") {
			refClassName = strings.TrimSuffix(refClassName, ";")
		}
	} else if strings.HasPrefix(refClassName, "[") {
		refClassName = ""
	}
	return refClassName
}

// Init simply initializes the three classloaders and points them to each other
// in the proper order. This function might be substantially revised later.
func Init(gl *globals.Globals) error {
	BootstrapCL.Name = "bootstrap"
	BootstrapCL.Parent = ""
	BootstrapCL.Classes = make(map[string]ParsedClass)

	ExtensionCL.Name = "extension"
	ExtensionCL.Parent = "bootstrap"
	ExtensionCL.Classes = make(map[string]ParsedClass)

	AppCL.Name = "app"
	AppCL.Parent = "system"
	AppCL.Classes = make(map[string]ParsedClass)

	gl.MethArea = &exec.Classes
	return nil
}
