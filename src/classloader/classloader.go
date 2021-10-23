/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"errors"
	"fmt"
	"jacobin/log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

// Classloader holds the parsed bytecode in classes, where they can be retrieved
// and moved to an execution role. Most of the comments and code presuppose some
// familiarity with the role of classloaders. More information can be found at:
// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-5.html#jvms-5.3
type Classloader struct {
	Name    string
	Parent  string
	Classes map[string]parsedClass
}

// AppCL is the application classloader, which loads most of the app's classes
var AppCL Classloader

// BootstrapCL is the classloader that loads most of the standard libraries
var BootstrapCL Classloader

// ExtensionCL is the classloader typically used for loading custom agents
var ExtensionCL Classloader

// the parsed class
type parsedClass struct {
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

// LoadClassFromFile first canonicalizes the filename, checks whether
// the class is already loaded, and if not, then parses the class and loads it.
//
// 1 TODO: canonicalize class name
// 2 TODO: search through classloaders for this class
// 3 TODO: determine which classloader should load the class, then
// 4 TODO: have *it* parse and load the class.
func (cl Classloader) LoadClassFromFile(filename string) error {
	rawBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Log("Could not read file: "+filename+". Exiting.", log.SEVERE)
		return fmt.Errorf("file I/O error")
	}

	log.Log(filename+" read", log.FINE)

	fullyParsedClass, err := parse(rawBytes)
	if err != nil {
		log.Log("error parsing "+filename+". Exiting.", log.SEVERE)
		return fmt.Errorf("parsing error")
	}

	err = formatCheckClass(&fullyParsedClass)
	if err != nil {
		log.Log("error format-checking "+filename+". Exiting.", log.SEVERE)
		return fmt.Errorf("format-checking error")
	}
	log.Log("Class "+fullyParsedClass.className+" has been format-checked.", log.FINEST)

	return insert(fullyParsedClass)

}

// Init simply initializes the three classloaders and points them to each other
// in the proper order. This function might be substantially revised later.
func Init() error {
	BootstrapCL.Name = "bootstrap"
	BootstrapCL.Parent = ""
	BootstrapCL.Classes = make(map[string]parsedClass)

	ExtensionCL.Name = "extension"
	ExtensionCL.Parent = "bootstrap"
	ExtensionCL.Classes = make(map[string]parsedClass)

	AppCL.Name = "app"
	AppCL.Parent = "system"
	AppCL.Classes = make(map[string]parsedClass)
	return nil
}

// insert the fully parsed class into the classloader
func insert(class parsedClass) error {
	return nil //TODO: fill out after finishing parser
}
