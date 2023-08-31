/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"errors"
	"fmt"
	"jacobin/log"
	"jacobin/shutdown"
)

// the definition of the class as it's stored in the method area
type Klass struct {
	Status byte // I=Initializing,F=formatChecked,V=verified,L=linked,N=instantiated
	Loader string
	Data   *ClData
}

type ClData struct {
	Name       string
	Superclass string
	Module     string
	Pkg        string   // package name, if any. (so named, b/c 'package' is a golang keyword)
	Interfaces []uint16 // indices into UTF8Refs
	Fields     []Field
	Methods    []Method
	Attributes []Attr
	SourceFile string
	Bootstraps []BootstrapMethod
	CP         CPool
	Access     AccessFlags
	ClInit     byte // 0 = no clinit, 1 = clinit not run, 2 clinit run
}

type CPool struct {
	CpIndex        []CpEntry // the constant pool index to entries
	ClassRefs      []uint16  // points to a UTF8 entry in the CP bearing class name
	Doubles        []float64
	Dynamics       []DynamicEntry
	FieldRefs      []FieldRefEntry
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
	Utf8Refs []string
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
	Name        uint16 // index of the UTF-8 entry in the CP
	Desc        uint16 // index of the UTF-8 entry in the CP
	IsStatic    bool   // is the field static?
	Attributes  []Attr
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
	MaxStack   int
	MaxLocals  int
	Code       []byte
	Exceptions []CodeException // exception entries for this method
	Attributes []Attr          // the code attributes has its own sub-attributes(!)
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

// the boostrap methods, specified in the bootstrap class attribute
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

type InterfaceRefEntry struct { // type: 11 (interface reference)
	ClassIndex  uint16
	NameAndType uint16
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

// FetchMethodAndCP gets the method and the CP for the class of the method.
// It searches for the method first by checking the MTable (that is, the method table).
// If it doesn't find it there, then it looks for it in the class entry in MethArea.
// If it finds it there, then it loads that class into the MTable and returns that
// entry as the Method it's returning.
func FetchMethodAndCP(class, meth string, methType string) (MTentry, error) {
	origClassName := class
	for {
	startSearch:

		// has the class been loaded? If not, then load it now.
		if MethAreaFetch(class) == nil {
			err := LoadClassFromNameOnly(class)
			if err != nil {
				if meth == "main" {
					// the starting class is always loaded, so if main() isn't found
					// right away, don't go up superclasses, just bail.
					noMainError(origClassName)
					shutdown.Exit(shutdown.JVM_EXCEPTION)
				}
				_ = log.Log("FetchMethodAndCP: LoadClassFromNameOnly for "+class+" failed: "+err.Error(), log.WARNING)
				_ = log.Log(err.Error(), log.SEVERE)
				shutdown.Exit(shutdown.JVM_EXCEPTION)
			}
		}

		methFQN := class + "." + meth + methType // FQN = fully qualified name
		methEntry := MTable[methFQN]

		if methEntry.Meth != nil { // we found the entry in the MTable
			if methEntry.MType == 'J' {
				return MTentry{Meth: methEntry.Meth, MType: 'J'}, nil
			} else if methEntry.MType == 'G' {
				return MTentry{Meth: methEntry.Meth, MType: 'G'}, nil
			}
		}

		// method is not in the MTable, so find it and put it there
		err := WaitForClassStatus(class)
		if err != nil {
			errMsg := fmt.Sprintf("FetchMethodAndCP: %s", err.Error())
			_ = log.Log(errMsg, log.SEVERE)
			shutdown.Exit(shutdown.JVM_EXCEPTION)
			return MTentry{}, errors.New(errMsg) // dummy return needed for tests
		}

		k := MethAreaFetch(class)
		if k == nil {
			errMsg := fmt.Sprintf("FetchMethodAndCP: MethAreaFetch could not find class {%s}", class)
			_ = log.Log(errMsg, log.SEVERE)
			shutdown.Exit(shutdown.JVM_EXCEPTION)
			return MTentry{}, errors.New(errMsg) // dummy return needed for tests
		}

		if k.Loader == "" { // if class is not found, the zero value struct is returned
			// TODO: check superclasses if method not found
			errMsg := "FetchMethodAndCP: Null Loader in class: " + class
			_ = log.Log(errMsg, log.SEVERE)
			return MTentry{}, errors.New(errMsg) // dummy return needed for tests
		}

		// the class has been found (k) so now go down the list of methods until
		// we find one that matches the name we're looking for. Then return that
		// method along with a pointer to the CP
		for i := 0; i < len(k.Data.Methods); i++ {
			if k.Data.CP.Utf8Refs[k.Data.Methods[i].Name] == meth &&
				k.Data.CP.Utf8Refs[k.Data.Methods[i].Desc] == methType {
				m := k.Data.Methods[i]
				jme := JmEntry{
					accessFlags: m.AccessFlags,
					MaxStack:    m.CodeAttr.MaxStack,
					MaxLocals:   m.CodeAttr.MaxLocals,
					Code:        m.CodeAttr.Code,
					exceptions:  m.CodeAttr.Exceptions,
					attribs:     m.CodeAttr.Attributes,
					params:      m.Parameters,
					deprecated:  m.Deprecated,
					Cp:          &k.Data.CP,
				}
				MTable[methFQN] = MTentry{
					Meth:  jme,
					MType: 'J',
				}
				return MTentry{Meth: jme, MType: 'J'}, nil
			}
		}

		// if we're searching for main(), don't go up the list of superclasses
		if meth == "main" { // to be consistent with the JDK, we print this peculiar error message when main() is missing
			noMainError(origClassName)
			break
			// } else {
			// 	_ = log.Log("FetchMethodAndCP: Found class "+class+", but it did not contain method: "+meth, log.SEVERE)
		}

		// if we got this far, the method was not found, so check the superclass(es)
		if class == "java/lang/Object" { // if we're already at the topmost superclass, then stop the loop
			break
		} else {
			class = k.Data.Superclass
			goto startSearch
		}
	}

	// if we got this far, something went wrong with locating the method
	_ = log.Log("FetchMethodAndCP: Found class "+class+", but it did not contain method: "+meth, log.SEVERE)
	shutdown.Exit(shutdown.JVM_EXCEPTION)
	return MTentry{}, errors.New("method not found") // dummy return needed for tests
}

// error message when main() can't be found
func noMainError(className string) {
	_ = log.Log("Error: main() method not found in class "+className+"\n"+
		"Please define the main method as:\n"+
		"   public static void main(String[] args)", log.SEVERE)
}

// FetchUTF8stringFromCPEntryNumber fetches the UTF8 string using the CP entry number
// for that string in the designated ClData.CP. Returns "" on error.
func FetchUTF8stringFromCPEntryNumber(cp *CPool, entry uint16) string {
	if entry < 1 || entry >= uint16(len(cp.CpIndex)) {
		msg := fmt.Sprintf("FetchUTF8stringFromCPEntryNumber: entry=%d is out of bounds(1, %d)", entry, uint16(len(cp.CpIndex)))
		_ = log.Log(msg, log.SEVERE)
		return ""
	}

	u := cp.CpIndex[entry]
	if u.Type != UTF8 {
		msg := fmt.Sprintf("FetchUTF8stringFromCPEntryNumber: cp.CpIndex[%d].Type=%d, expected UTF8", entry, u.Type)
		_ = log.Log(msg, log.SEVERE)
		return ""
	}

	return cp.Utf8Refs[u.Slot]
}
