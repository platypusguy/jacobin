/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classdata

// To avoid circular dependencies, this file contains the fields definitions
// for classes used in metadata (as in the method area) and in statics. It's
// used by the classloader, jvm, and statics packages.
/*
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

// the structure of many attributes (field, class, etc.) The content is just the raw bytes.
type Attr struct {
	AttrName    uint16 // index of the UTF8 entry in the CP
	AttrSize    int    // length of the following array of raw bytes
	AttrContent []byte // the raw data of the attribute
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
	Bootstraps      []BootstrapMethod
	CP              CPool
	Access          AccessFlags
	ClInit          byte           // 0 = no clinit, 1 = clinit not run, 2 clinit
	ClassObject     *object.Object // the java/lang/Class object for this class
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

// the exception-related data for each exception in the Code attribute of a given method
type CodeException struct {
	StartPc   int    // first instruction covered by this exception (pc = program counter)
	EndPc     int    // the last instruction covered by this exception
	HandlerPc int    // the place in the method code that has the exception instructions
	CatchType uint16 // the type of exception, index to CP, which must point a ClassFref entry
}

// BytecodeToSourceLine maps the PC in a method to the
// corresponding source line in the original source file.
// This data is captured in the method's attributes
type BytecodeToSourceLine struct {
	BytecodePos uint16
	SourceLine  uint16
}

// the bootstrap methods, specified in the bootstrap class attribute
type BootstrapMethod struct {
	MethodRef uint16   // index pointing to a MethodHandle
	Args      []uint16 // arguments: indexes to loadable arguments from the CP
}

*/
