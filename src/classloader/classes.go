/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-5 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"errors"
	"fmt"
	"jacobin/src/object"
	"jacobin/src/shutdown"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"jacobin/src/types"
)

// the definition of the class as it's stored in the method area
type Klass struct {
	Status      byte // I=Initializing,F=formatChecked,V=verified,L=linked,N=instantiated
	Loader      string
	Data        *ClData
	CodeChecked bool // has the code been checked for this class?
	Resolved    bool // has the CP been resolved for this class?
}

// the class data, as it's posted to the method area
type ClData struct {
	Name            string
	NameIndex       uint32 // index into StringPool
	SuperclassIndex uint32 // index into StringPool
	Module          string
	Pkg             string   // package name, if any. (so named, b/c 'package' is a golang keyword)
	Interfaces      []uint16 // indices into UTF8Refs
	Fields          []Field
	MethodList      map[string]string  // maps method names including superclass methods to FQN, which is the key to GMT
	MethodTable     map[string]*Method // the methods defined in this class
	Attributes      []Attr
	SourceFile      string
	CP              CPool
	Access          AccessFlags
	ClInit          byte           // 0 = no clinit, 1 = clinit not run, 2 clinit
	ClassObject     *object.Object // the java/lang/Class object for this class
}

// the CP of the loaded class (see above)
type CPool struct {
	CpIndex        []CpEntry // the constant pool index to entries
	ClassRefs      []uint32  // points to a string pool entry = class name
	Doubles        []float64
	Dynamics       []DynamicEntry
	FieldRefs      []ResolvedFieldEntry
	Floats         []float32
	IntConsts      []int32 // 32-bit int containing the actual int value
	InterfaceRefs  []InterfaceRefEntry
	InvokeDynamics []InvokeDynamicEntry
	LongConsts     []int64
	MethodHandles  []MethodHandleEntry
	MethodRefs     []MethodRefEntry
	MethodTypes    []uint16
	NameAndTypes   []NameAndTypeEntry
	//	StringRefs     []uint16 // all StringRefs are converted into utf8Refs
	Utf8Refs              []string
	Bootstraps            []BootstrapMethod           // not technically part of the CP, but convenient to store here
	ResolvedInterfaceRefs []ResolvedInterfaceRefEntry // resolved interface references
	ResolvedMethodRefs    []ResolvedMethodRefEntry    // resolved method references
}

type AccessFlags struct {
	ClassIsPublic     bool
	ClassIsFinal      bool
	ClassIsSuper      bool
	ClassIsInterface  bool
	ClassIsAbstract   bool
	ClassIsSynthetic  bool
	ClassIsAnnotation bool
	ClassIsEnum       bool
	ClassIsModule     bool
}

// For the nonce, these definitions are similar to corresponding items in
// classloader.go. The biggest difference is that ints there often become uint16
// here (where correct to do so). This greatly reduces memory consumption.
// Likewise certain fields needed there (counts) are not used here.

type Field struct {
	AccessFlags int
	NameStr     string
	Name        uint16      // index of the UTF-8 entry in the CP
	Desc        uint16      // index of the UTF-8 entry in the CP
	DescStr     string      // the type of the field, as a string (using Desc) JACOBIN-720
	IsStatic    bool        // is the field static?
	ConstValue  interface{} // if static and has constant value, it's stored here.
	Attributes  []Attr      // all attributes for this field other than ConstantValue
}

// the methods of the class, including the constructors
type Method struct {
	AccessFlags int
	Name        uint16 // index of the UTF-8 entry in the CP
	Desc        uint16 // index of the UTF-8 entry in the CP
	CodeAttr    CodeAttrib
	Attributes  []Attr
	Exceptions  []uint16 // indexes into Utf8Refs in the CP
	Parameters  []ParamAttrib
	Deprecated  bool // is the method deprecated?
}

type CodeAttrib struct {
	MaxStack          int
	MaxLocals         int
	Code              []byte
	Exceptions        []CodeException // exception entries for this method
	Attributes        []Attr          // the code attributes has its own sub-attributes(!)
	BytecodeSourceMap []BytecodeToSourceLine
}

// ParamAttrib is the MethodParameters method attribute
type ParamAttrib struct {
	Name        string // string, rather than index into utf8Refs b/c the name could be ""
	AccessFlags int
}

// the structure of many attributes (field, class, etc.) The content is just the raw bytes.
type Attr struct {
	AttrName    uint16 // index of the UTF8 entry in the CP
	AttrSize    int    // length of the following array of raw bytes
	AttrContent []byte // the raw data of the attribute
}

// the exception-related data for each exception in the Code attribute of a given method
type CodeException struct {
	StartPc   int    // first instruction covered by this exception (pc = program counter)
	EndPc     int    // the last instruction covered by this exception
	HandlerPc int    // the place in the method code that has the exception instructions
	CatchType uint16 // the type of exception, index to CP, which must point a ClassFref entry
}

// the bootstrap methods, specified in the bootstrap class attribute
type BootstrapMethod struct {
	MethodRef uint16   // index pointing to a MethodHandle
	Args      []uint16 // arguments: indexes to loadable arguments from the CP
}

// ==== Constant Pool structs (in order by their numeric code) ====//
type CpEntry struct {
	Type uint16
	Slot uint16
}

type FieldRefEntry struct { // type: 09 (field reference)
	ClassIndex  uint16
	NameAndType uint16
}

type MethodRefEntry struct { // type: 10 (method reference)
	ClassIndex  uint16
	NameAndType uint16
}

type ResolvedMethodRefEntry struct { // type: 10 (method reference, resolved)
	ClassIndex  uint32 // all of these are indices into the StringPool
	NameIndex   uint32
	TypeIndex   uint32
	FQNameIndex uint32 // the three previous strings appended into one entry (the most common usage)
}

type InterfaceRefEntry struct { // type: 11 (interface reference)
	ClassIndex  uint16
	NameAndType uint16
}

type ResolvedInterfaceRefEntry struct { // type: 11 (interface reference, resolved)
	ClassIndex  uint32 // all of these are indices into the StringPool
	NameIndex   uint32
	TypeIndex   uint32
	FQNameIndex uint32 // the three previous strings appended into one entry (the most common usage)
}

type NameAndTypeEntry struct { // type 12 (name and type reference)
	NameIndex uint16
	DescIndex uint16
}

type MethodHandleEntry struct { // type: 15 (method handle)
	RefKind  uint16
	RefIndex uint16
}

type DynamicEntry struct { // type 17 (dynamic--similar to invokedynamic)
	BootstrapIndex uint16
	NameAndType    uint16
}

type InvokeDynamicEntry struct { // type 18 (invokedynamic data)
	BootstrapIndex uint16
	NameAndType    uint16
}

// // the various types of entries in the constant pool. These entries are duplicates
// // of the ones in cpParser.go. These lists should be kept in sync.
// const (
// 	Dummy              = 0 // used for initialization and for dummy entries (viz. for longs, doubles)
// 	UTF8               = 1
// 	IntConst           = 3
// 	FloatConst         = 4
// 	LongConst          = 5
// 	DoubleConst        = 6
// 	ClassRef           = 7
// 	StringConst        = 8
// 	FieldRef           = 9
// 	MethodRef          = 10
// 	Interface          = 11
// 	NameAndType        = 12
// 	MethodHandle       = 15
// 	MethodType         = 16
// 	DynamicEntry       = 17
// 	InvokeDynamicEntry = 18
// 	Module             = 19
// 	Package            = 20
// )

// FetchMethodAndCP gets a method and the CP for the class of the method. It searches
// for the method first by checking the global MTable (that is, the global method table).
// If it doesn't find it there, then it looks for the method in the class entry in MethArea.
// If it finds it there, then it loads that class into the MTable and returns that
// entry as the Method it's returning.
//
// Note that if the given method is not found, the hierarchy of superclasses is ascended,
// in search for the method. The one exception is for main() which, if not found in the
// first class, will never be in one of the superclasses.
//
// Note: if the method is not in the class or its superclasses, an error is returned. This
// method does not check interfaces for the method.
func FetchMethodAndCP(className, methName, methType string) (MTentry, error) {
	origClassName := className

	// has the className been loaded? If not, then load it now.
	if MethAreaFetch(className) == nil {
		err := LoadClassFromNameOnly(className)
		if err != nil {
			if methName == "main" {
				// the starting className is always loaded, so if main() isn't found
				// something is seriously wrong, so show the specificerror and shutdown.
				noMainError(origClassName)
				// noMainError() calls shutdown.Exit(). However, in test mode, shutdown.Exit() doesn't exit,
				// so the following error return is needed to cover the test cases.
				return MTentry{}, errors.New("Error: main() method not found in class " + origClassName + "\n")
			} else {
				errMsg := fmt.Sprintf("FetchMethodAndCP: LoadClassFromNameOnly for %s failed: %s",
					className, err.Error())
				trace.Error(errMsg)
				shutdown.Exit(shutdown.JVM_EXCEPTION)
				return MTentry{}, errors.New(errMsg) // dummy return needed for tests
			}
		}
	}

	// --- at this point we know the class exists and has been loaded ---

	// look for the method in the MTable
	methFQN := className + "." + methName + methType // FQN = fully qualified name
	methEntry := GetMtableEntry(methFQN)

	if methEntry.Meth != nil { // we found the entry in the MTable
		if methEntry.MType == 'J' {
			return MTentry{Meth: methEntry.Meth, MType: 'J'}, nil
		}
		if methEntry.MType == 'G' {
			return MTentry{Meth: methEntry.Meth, MType: 'G'}, nil
		}
		errMsg := fmt.Sprintf("FetchMethodAndCP: methEntry.Meth != nil BUT methEntry.MType is neither J nor G for %s", methFQN)
		trace.Error(errMsg)
		shutdown.Exit(shutdown.JVM_EXCEPTION)
		return MTentry{}, errors.New(errMsg) // dummy return needed for tests
	}

	// --- at this point, the method is not in the MTable ---

	// While the class has been loaded, it might not have been initialized.
	// The method is not in the MTable, so, find it, and insert it there:
	// This is done by making sure the class has been initialized,
	// then fetching it from the MethArea, then searching its methodTable
	err := WaitForClassStatus(className)
	if err != nil {
		errMsg := fmt.Sprintf("FetchMethodAndCP: %s", err.Error())
		trace.Error(errMsg)
		shutdown.Exit(shutdown.JVM_EXCEPTION)
		return MTentry{}, errors.New(errMsg) // dummy return needed for tests
	}

	k := MethAreaFetch(className)
	if k == nil {
		errMsg := fmt.Sprintf("FetchMethodAndCP: MethAreaFetch could not find class %s", className)
		trace.Error(errMsg)
		shutdown.Exit(shutdown.JVM_EXCEPTION)
		return MTentry{}, errors.New(errMsg) // dummy return needed for tests
	}

	// the class, k, has been found, so check the method table for the method. Then return the
	// method along with a pointer to the CP
	var m Method
	searchName := methName + methType
	methRef, ok := k.Data.MethodTable[searchName]
	if ok {
		m = *methRef

		// create a Java method struct for this method. We know it's a Java method
		// because if it were a native method it would have been found in the initial
		// lookup in the MTable (as all native methods are loaded there before
		// program execution begins).
		jme := JmEntry{
			AccessFlags: m.AccessFlags,
			MaxStack:    m.CodeAttr.MaxStack,
			MaxLocals:   m.CodeAttr.MaxLocals,
			Code:        m.CodeAttr.Code,
			Exceptions:  m.CodeAttr.Exceptions, // just a list of CP indexes to exceptions thrown by this method
			Attribs:     m.CodeAttr.Attributes,
			params:      m.Parameters,
			deprecated:  m.Deprecated,
			Cp:          &k.Data.CP,
		}

		// add the method to the MTable and return it
		methodEntry := MTentry{Meth: jme, MType: 'J'}
		AddEntry(&MTable, methFQN, methodEntry)
		return methodEntry, nil
	}

	// if we're here, the className did not contain the searched-for method. So, go up the superclasses,
	// except if we're searching for main(), in which case, we don't go up the list of superclasses
	if methName == "main" {
		noMainError(origClassName)
		// even though noMainError() exits, in testing the exit is disabled,
		// so we add this return statement to test for correct operation
		return MTentry{}, errors.New("main() not found")
	} else {
		// go up the list of superclasses

	superclassLoop:

		// Get the superclass name.
		className = *stringPool.GetStringPointer(k.Data.SuperclassIndex)

		// Matching a special Jacobin class?
		if className == types.ClassNameThread || className == types.ClassNameThreadGroup {
			methFQN = className + "." + methName + methType
			methEntry = GetMtableEntry(methFQN)
			if methEntry.Meth != nil { // we found the entry in the MTable
				if methEntry.MType == 'G' {
					return methEntry, nil
				}
			}
		}

		// Get the method area of the class.
		k = MethAreaFetch(className)
		if k == nil {
			errMsg := fmt.Sprintf("FetchMethodAndCP: MethAreaFetch could not find superclass %s", className)
			trace.Error(errMsg)
			shutdown.Exit(shutdown.JVM_EXCEPTION)
			return MTentry{}, errors.New(errMsg) // dummy return needed for tests
		}

		// Search for the method in the class method table.
		methRef, ok = k.Data.MethodTable[searchName]
		if ok {
			m = *methRef

			// create a Java method struct for this method. We know it's a Java method
			// because if it were a native method it would have been found in the initial
			// lookup in the MTable (as all native methods are loaded there before
			// program execution begins).
			jme := JmEntry{
				AccessFlags: m.AccessFlags,
				MaxStack:    m.CodeAttr.MaxStack,
				MaxLocals:   m.CodeAttr.MaxLocals,
				Code:        m.CodeAttr.Code,
				Exceptions:  m.CodeAttr.Exceptions, // note that the CodeAttr.Code exceptions are placed here
				Attribs:     m.CodeAttr.Attributes,
				params:      m.Parameters,
				deprecated:  m.Deprecated,
				Cp:          &k.Data.CP,
			}

			// add the method to the MTable and return it
			methodEntry := MTentry{Meth: jme, MType: 'J'}
			AddEntry(&MTable, methFQN, methodEntry)
			return methodEntry, nil

		} else {

			// if we've ascended to Object and don't have the method, it ain't here (error).
			if className != types.ObjectClassName {
				goto superclassLoop
			} else {
				errMsg := fmt.Sprintf("FetchMethodAndCP: Neither %s nor its superclasses contain method %s",
					origClassName, methName)
				return MTentry{}, errors.New(errMsg)
			}
		}
	}
}

// error message when main() can't be found. Syntax mirrors OpenJDK HotSpot
func noMainError(className string) {
	errMsg := fmt.Sprintf(
		"Error: main() method not found in class %s\n"+
			"Please define the main method as:\n"+
			"   public static void main(String[] args)", className)
	trace.Error(errMsg)
	shutdown.Exit(shutdown.APP_EXCEPTION)
}

// FetchUTF8stringFromCPEntryNumber fetches the UTF8 string using the CP entry number
// for that string in the designated ClData.CP. Returns "" on error.
func FetchUTF8stringFromCPEntryNumber(cp *CPool, entry uint16) string {
	if entry < 1 || entry >= uint16(len(cp.CpIndex)) {
		errMsg := fmt.Sprintf("FetchUTF8stringFromCPEntryNumber: entry=%d is out of bounds(1, %d)",
			entry, uint16(len(cp.CpIndex)))
		trace.Error(errMsg)
		return ""
	}

	u := cp.CpIndex[entry]
	if u.Type != UTF8 {
		errMsg := fmt.Sprintf("FetchUTF8stringFromCPEntryNumber: cp.CpIndex[%d].Type=%d, expected UTF8", entry, u.Type)
		trace.Error(errMsg)
		return ""
	}

	return cp.Utf8Refs[u.Slot]
}
