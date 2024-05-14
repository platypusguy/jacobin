/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/exceptions"
	"jacobin/log"
	"jacobin/object"
	"jacobin/types"
	"os"
	"strings"
)

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

// MTableLoadNatives loads the Go methods from files that contain them. It does this
// by calling the Load_* function in each of those files to load whatever Go functions
// they make available.
func MTableLoadNatives(MTable *classloader.MT) {

	loadlib(MTable, Load_Io_Console()) // load the java.io.Console golang functions

	loadlib(MTable, Load_Io_BufferedReader())
	loadlib(MTable, Load_Io_File())
	loadlib(MTable, Load_Io_FileInputStream())
	loadlib(MTable, Load_Io_FileOutputStream())
	loadlib(MTable, Load_Io_FileReader())
	loadlib(MTable, Load_Io_FileWriter())
	loadlib(MTable, Load_Io_InputStreamReader())
	loadlib(MTable, Load_Io_OutputStreamWriter())
	loadlib(MTable, Load_Io_PrintStream())
	loadlib(MTable, Load_Io_RandomAccessFile())

	loadlib(MTable, Load_Lang_Boolean())
	loadlib(MTable, Load_Lang_Byte())
	loadlib(MTable, Load_Lang_Character())
	loadlib(MTable, Load_Lang_Class()) // load the java.lang.Class golang functions
	loadlib(MTable, Load_Lang_Double())
	loadlib(MTable, Load_Lang_Float())
	loadlib(MTable, Load_Lang_Integer())
	loadlib(MTable, Load_Lang_Long())
	loadlib(MTable, Load_Lang_Math())   // load the java.lang.Math & StrictMath golang functions
	loadlib(MTable, Load_Lang_Object()) // load the java.lang.Class golang functions
	loadlib(MTable, Load_Lang_Short())
	loadlib(MTable, Load_Lang_String())            // load the java.lang.String golang functions
	loadlib(MTable, Load_Lang_StringBuilder())     // load the java.lang.StringBuilder golang functions
	loadlib(MTable, Load_Lang_System())            // load the java.lang.System golang functions
	loadlib(MTable, Load_Lang_StackTraceELement()) //  java.lang.StackTraceElement golang functions
	loadlib(MTable, Load_Lang_Thread())            // load the java.lang.Thread golang functions
	loadlib(MTable, Load_Lang_Throwable())         // load the java.lang.Throwable golang functions (errors & exceptions)
	loadlib(MTable, Load_Lang_UTF16())             // load the java.lang.UTF16 golang functions

	loadlib(MTable, Load_Nio_Charset_Charset()) // Zero Charset support

	loadlib(MTable, Load_Util_Concurrent_Atomic_AtomicInteger())
	loadlib(MTable, Load_Util_Concurrent_Atomic_Atomic_Long())
	loadlib(MTable, Load_Util_HashMap())
	loadlib(MTable, Load_Util_HexFormat())
	loadlib(MTable, Load_Util_Locale())
	loadlib(MTable, Load_Util_Random())

	loadlib(MTable, Load_Jdk_Internal_Misc_Unsafe())
	loadlib(MTable, Load_Jdk_Internal_Misc_ScopedMemoryAccess())

	loadlib(MTable, Load_Nil_Clinit()) // Load <clinit> functions that invoke justReturn()
	loadlib(MTable, Load_Traps())      // Load traps that lead to unconditional error returns

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
			log.Log(errMsg, log.SEVERE)
			ok = false
		}
		gme := GMeth{}
		gme.ParamSlots = val.ParamSlots
		gme.GFunction = val.GFunction
		gme.NeedsContext = val.NeedsContext

		tableEntry := classloader.MTentry{
			MType: 'G',
			Meth:  gme,
		}

		classloader.AddEntry(tbl, key, tableEntry)
	}
	if !ok {
		exceptions.ThrowExNil(excNames.InternalException, "loadlib: at least one key was invalid")
	}
}

// do-nothing Go function shared by several source files
func justReturn([]interface{}) interface{} {
	return nil
}

// Populate an object for a primitive type (Byte, Character, Double, Float, Integer, Long, Short, String).
func populator(classname string, fldtype string, fldvalue interface{}) interface{} {
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
