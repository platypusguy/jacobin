/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package exec

import (
	"strings"
	"sync"
)

var Classes = make(map[string]Klass) // make this a sync.Map

type Klass struct {
	Status byte // I=Initializing,F=formatChecked,V=verified,L=linked
	Loader string
	Data   *ClData
}

var MethAreaMutex sync.RWMutex // All additions or updates to Classes map come through this mutex

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
	Dynamics       []Dynamic
	FieldRefs      []FieldRefEntry
	Floats         []float32
	IntConsts      []int32 // 32-bit int containing the actual int value
	InterfaceRefs  []InterfaceRefEntry
	InvokeDynamics []InvokeDynamic
	LongConsts     []int64
	MethodHandles  []MethodHandleEntry
	MethodRefs     []MethodRefEntry
	MethodTypes    []uint16
	NameAndTypes   []NameAndTypeEntry
	StringRefs     []uint16 // integer index into utf8Refs
	Utf8Refs       []string
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

// the MethodParameters method attribute
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

type Dynamic struct { // type 17 (dynamic--similar to invokedynamic)
	BootstrapIndex uint16
	NameAndType    uint16
}

type InvokeDynamic struct { // type 18 (invokedynamic data)
	BootstrapIndex uint16
	NameAndType    uint16
}

// fetches the UTF8 string using the CP entry number for that string in the
// designated ClData.CP. Returns "" on error.
func FetchUTF8stringFromCPEntryNumber(cp *CPool, entry uint16) string {
	if entry < 1 || entry >= uint16(len(cp.CpIndex)) {
		return ""
	}

	u := cp.CpIndex[entry]
	if u.Type != 1 { // cannot test for classloader.UTF8 b/c cannot create dependence
		return ""
	}

	return cp.Utf8Refs[u.Slot]
}

// accepts a class name with the JVM's internal format and converts
// it to a filename (with backslashes). Returns "" on error.
func ConvertInternalClassNameToFilename(clName string) string {
	name := strings.ReplaceAll(clName, "/", "\\")
	name = strings.ReplaceAll(clName, "\\", ".") + ".class"
	return name
}

func ConvertClassFilenameToInternalFormat(fName string) string {
	name := strings.TrimSuffix(fName, ".class")
	strings.ReplaceAll(name, ".", "\\")
	return name
}
