/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-5 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/exceptions"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/trace"
	"jacobin/types"
	"math/big"
	"os"
	"strings"
)

// Map repository of method signatures for all G functions:
var MethodSignatures = make(map[string]GMeth)
var TestMethodSignatures = make(map[string]GMeth) // used only for the test gfunctions
var TestGfunctionsLoaded = false

// File I/O and stream Field keys:
var FileStatus string = "status"     // using this value in case some member function is looking at it
var FilePath string = "FilePath"     // full absolute path of a file aka canonical path
var FileHandle string = "FileHandle" // *os.File
var FileMark string = "FileMark"     // file position relative to beginning (0)
var FileAtEOF string = "FileAtEOF"   // file at EOF

// File I/O constants:
var CreateFilePermissions os.FileMode = 0664 // When creating, read and write for user and group, others read-only

// Radix boundaries:
var MinRadix int64 = 2
var MaxRadix int64 = 36
var MaxIntValue int64 = 2147483647
var MinIntValue int64 = -2147483648

// GMeth is the entry in the MTable for Go functions. See MTable comments for details.
//   - ParamSlots - the number of user parameters in a G function. E.g. For atan2, this would be 2.
//   - GFunction - a go function. All go functions accept a possibly empty slice of interface{} and
//     return an interface{} which might be nil (E.g. Java void).
//   - NeedsContext - does this method need a pointer to the frame stack? Defaults to false.
type GMeth struct {
	ParamSlots   int
	GFunction    func([]interface{}) interface{}
	NeedsContext bool
	ThreadSafe   bool
}

// G function error block.
type GErrBlk struct {
	ExceptionType int
	ErrMsg        string
}

// Construct a G function error block. Return a ptr to it.
func getGErrBlk(exceptionType int, errMsg string) *GErrBlk {
	var gErrBlk GErrBlk
	gErrBlk.ExceptionType = exceptionType
	gErrBlk.ErrMsg = errMsg
	return &gErrBlk
}

// do-nothing Go function shared by several source files
func clinitGeneric([]interface{}) interface{} {
	return object.StringObjectFromGoString("clinitGeneric")
}

// do-nothing Go function shared by several source files
func justReturn([]interface{}) interface{} {
	return object.StringObjectFromGoString("justReturn")
}

// return a Java null object.
func returnNull([]interface{}) interface{} {
	return object.Null
}

func returnCharsetName([]interface{}) interface{} {
	return object.StringObjectFromGoString(globals.GetCharsetName())
}

// MTableLoadGFunctions loads the Go methods from files that contain them. It does this
// by calling the Load_* function in each of those files to load whatever Go functions
// they make available.
func MTableLoadGFunctions(MTable *classloader.MT) {

	// java/awt/*
	Load_Awt_Graphics_Environment()

	// java/io/*
	Load_Io_BufferedReader()
	Load_Io_Console()
	Load_Io_File()
	Load_Io_FileInputStream()
	Load_Io_FileOutputStream()
	Load_Io_FileReader()
	Load_Io_FileWriter()
	Load_Io_FilterInputStream()
	Load_Io_InputStreamReader()
	Load_Io_OutputStreamWriter()
	Load_Io_PrintStream()
	Load_Io_RandomAccessFile()

	// java/lang/*
	Load_Lang_Boolean()
	Load_Lang_Byte()
	Load_Lang_Character()
	Load_Lang_Class()
	Load_Lang_Double()
	Load_Lang_Float()
	Load_Lang_Integer()
	Load_Lang_Long()
	Load_Lang_Math()
	Load_Lang_Object()
	Load_Lang_Runtime()
	Load_Lang_Short()
	Load_Lang_StackTraceELement()
	Load_Lang_String()
	Load_Lang_StringBuffer()
	Load_Lang_StringBuilder()
	Load_Lang_System()
	Load_Lang_Thread()
	Load_Lang_Throwable()
	Load_Lang_UTF16()

	// java/math/*
	Load_Math_Big_Integer()

	// java/security/*
	Load_Security_SecureRandom()
	Load_Security_AccessController()

	// java/util/*
	Load_Util_Arrays()
	Load_Util_Concurrent_Atomic_AtomicInteger()
	Load_Util_Concurrent_Atomic_Atomic_Long()
	Load_Util_Hash_Map()
	Load_Util_HexFormat()
	Load_Util_Locale()
	Load_Util_Objects()
	Load_Util_Random()
	Load_Util_Zip_Adler32()
	Load_Util_Zip_Crc32_Crc32c()

	// jdk/internal/misc/*
	Load_Jdk_Internal_Misc_Unsafe()
	Load_Jdk_Internal_Misc_ScopedMemoryAccess()

	// Load functions that invoke clinitGeneric() and do nothing else.
	Load_Other_methods()

	// Load traps that lead to unconditional error returns.
	Load_Traps()

	// jacobin JVM diagnostic routines
	Load_jj()

	//	now, with the accumulated MethodSignatures maps, load MTable.
	loadlib(MTable, MethodSignatures)

}

// load the test gfunctions in testGfunctions.go
func LoadTestGfunctions(MTable *classloader.MT) {
	Load_TestGfunctions()
	loadlib(MTable, TestMethodSignatures)
	TestGfunctionsLoaded = true
}

func checkKey(key string) bool {
	if strings.Index(key, ".") == -1 || strings.Index(key, "(") == -1 || strings.Index(key, ")") == -1 {
		return false
	}
	if strings.HasSuffix(key, ")") {
		return false
	}
	return true
}

func loadlib(tbl *classloader.MT, libMeths map[string]GMeth) {
	ok := true
	for key, val := range libMeths {
		if !checkKey(key) {
			errMsg := fmt.Sprintf("loadlib: Invalid key=%s", key)
			trace.Error(errMsg)
			ok = false
		}
		gme := GMeth{}
		gme.ParamSlots = val.ParamSlots
		gme.GFunction = val.GFunction
		gme.NeedsContext = val.NeedsContext
		gme.ThreadSafe = val.ThreadSafe

		tableEntry := classloader.MTentry{
			MType: 'G',
			Meth:  gme,
		}

		classloader.AddEntry(tbl, key, tableEntry)
		classloader.GmtAddEntry(key, classloader.GmtEntry{ // JACOBIN-575 adding to GMT in parallel
			MethData: &gme,
			MType:    'G',
		})
	}
	if !ok {
		exceptions.ThrowExNil(excNames.InternalException, "loadlib: at least one key was invalid")
	}
}

// Populate an object for a primitive type (Byte, Character, Double, Float, Integer, Long, Short, String).
func Populator(classname string, fldtype string, fldvalue interface{}) *object.Object {
	var objPtr *object.Object
	if fldtype == types.StringIndex {
		objPtr = object.StringObjectFromGoString(fldvalue.(string))
	} else {
		objPtr = object.MakePrimitiveObject(classname, fldtype, fldvalue)
		(*objPtr).FieldTable["value"] = object.Field{fldtype, fldvalue}
	}
	return objPtr
}

// File set EOF condition.
func eofSet(obj *object.Object, value bool) {
	obj.FieldTable[FileAtEOF] = object.Field{Ftype: types.Bool, Fvalue: value}
}

// File get EOF boolean.
func eofGet(obj *object.Object) bool {
	value, ok := obj.FieldTable[FileAtEOF].Fvalue.(bool)
	if !ok {
		return false
	}
	return value
}

// Return a null object.
func returnNullObject(params []interface{}) interface{} {
	return object.Null
}

// Return false.
func returnFalse(params []interface{}) interface{} {
	return types.JavaBoolFalse
}

// Return true.
func returnTrue(params []interface{}) interface{} {
	return types.JavaBoolTrue
}

// InitBigIntegerField: Initialise the object field.
// Fvalue holds *big.Int (pointer).
func InitBigIntegerField(obj *object.Object, argValue int64) {
	ptrBigInt := big.NewInt(argValue)
	fldValue := object.Field{Ftype: types.BigInteger, Fvalue: ptrBigInt}
	obj.FieldTable["value"] = fldValue
	var fldSign object.Field
	switch {
	case argValue == 0:
		fldSign = object.Field{Ftype: types.BigInteger, Fvalue: int64(0)}
	case argValue < 0:
		fldSign = object.Field{Ftype: types.BigInteger, Fvalue: int64(-1)}
	default:
		fldSign = object.Field{Ftype: types.BigInteger, Fvalue: int64(+1)}
	}
	obj.FieldTable["signum"] = fldSign
}
