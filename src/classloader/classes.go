/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"errors"
	"jacobin/log"
	"time"
)

type Klass struct {
	Status byte // I=Initializing,F=formatChecked,V=verified,L=linked,N=instantiated
	Loader string
	Data   *ClData
}

type ClData struct {
	Name       string
	Superclass string
	Module     string
	Pkg        string   // package name, if any. ('package' is a golang keyword)
	Interfaces []uint16 // indices into UTF8Refs
	Fields     []Field
	Methods    []Method
	Attributes []Attr
	SourceFile string
	Bootstraps []BootstrapMethod
	CP         CPool
	Access     AccessFlags
}

type CPool struct {
	CpIndex        []CpEntry // the constant pool index to entries
	ClassRefs      []uint16  // points to a UTF8 entry in the CP
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
	methFQN := class + "." + meth + methType // FQN = fully qualified name
	methEntry := MTable[methFQN]
	if methEntry.Meth == nil { // method is not in the MTable, so find it and put it there
		k := MethAreaFetch(class)
		if k.Status == 'I' { // class is being initialized by a loader, so wait
			time.Sleep(15 * time.Millisecond) // TODO: must be a better way to do this
			k = MethAreaFetch(class)
		}

		if k.Loader == "" { // if class is not found, the zero value struct is returned
			// TODO: check superclasses if method not found
			_ = log.Log("Could not find class: "+class, log.SEVERE)
			return MTentry{}, errors.New("class not found")
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
	} else { // we found the entry in the MTable
		if methEntry.MType == 'J' {
			return MTentry{Meth: methEntry.Meth, MType: 'J'}, nil
		} else if methEntry.MType == 'G' {
			return MTentry{Meth: methEntry.Meth, MType: 'G'}, nil
		}
	}

	// if we got this far, the class was not found

	if meth == "main" { // to be consistent with the JDK, we print this peculiar error message when main() is missing
		_ = log.Log("Error: Main method not found in class "+class+", please define the main method as:\n"+
			"   public static void main(String[] args)", log.SEVERE)
	} else {
		_ = log.Log("Found class: "+class+", but it did not contain method: "+meth, log.SEVERE)
	}

	return MTentry{}, errors.New("method not found")
}

// FetchUTF8stringFromCPEntryNumber fetches the UTF8 string using the CP entry number
// for that string in the designated ClData.CP. Returns "" on error.
func FetchUTF8stringFromCPEntryNumber(cp *CPool, entry uint16) string {
	if entry < 1 || entry >= uint16(len(cp.CpIndex)) {
		return ""
	}

	u := cp.CpIndex[entry]
	if u.Type != UTF8 {
		return ""
	}

	return cp.Utf8Refs[u.Slot]
}
