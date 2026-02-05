/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-5 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package ghelpers

import (
	"container/list"
	"crypto/rand"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/exceptions"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/trace"
	"jacobin/src/types"
	"math/big"
	"os"
)

// Map repository of method signatures for all G functions:
var MethodSignatures = make(map[string]GMeth)
var TestMethodSignatures = make(map[string]GMeth) // used only for the test gfunctions
var TestGfunctionsLoaded = false

// TrapClass is a generic Trap for classes
func TrapClass([]interface{}) interface{} {
	errMsg := "TRAP: The requested class is not yet supported"
	return GetGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// TrapDeprecated is a generic Trap for deprecated classes and functions
func TrapDeprecated([]interface{}) interface{} {
	errMsg := "TRAP: The requested class or function is deprecated and, therefore, not supported"
	return GetGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// TrapUndocumented is a generic Trap for deprecated classes and functions
func TrapUndocumented([]interface{}) interface{} {
	errMsg := "TRAP: The requested class or function is undocumented and, therefore, not supported"
	return GetGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// TrapFunction is a generic Trap for functions
func TrapFunction([]interface{}) interface{} {
	errMsg := "TRAP: The requested function is not yet supported"
	return GetGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// TrapProtected is a generic Trap for functions
func TrapProtected([]interface{}) interface{} {
	errMsg := "TRAP: The requested function is protected"
	return GetGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

func TrapUnicode([]interface{}) interface{} {
	return GetGErrBlk(
		excNames.UnsupportedOperationException,
		"Character Unicode method not yet implemented",
	)
}

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

// GetGErrBlk constructs a G function error block. Return a ptr to it.
func GetGErrBlk(exceptionType int, errMsg string) *GErrBlk {
	var gErrBlk GErrBlk
	gErrBlk.ExceptionType = exceptionType
	gErrBlk.ErrMsg = errMsg
	return &gErrBlk
}

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

// int64 value boundaries:
var MaxIntValue int64 = 2147483647
var MinIntValue int64 = -2147483648

// ClinitGeneric is a do-nothing Go function shared by several source files
func ClinitGeneric([]interface{}) interface{} {
	return nil
}

// JustReturn is a do-nothing Go function shared by several source files
func JustReturn([]interface{}) interface{} {
	return nil
}

// ReturnNull returns a Java null object.
func ReturnNull([]interface{}) interface{} {
	return object.Null
}

// ReturnNullObject returns a null object.
func ReturnNullObject([]interface{}) interface{} {
	return object.Null
}

func ReturnCharsetName([]interface{}) interface{} {
	return object.StringObjectFromGoString(globals.GetCharsetName())
}

// ReturnFalse returns false.
func ReturnFalse([]interface{}) interface{} {
	return types.JavaBoolFalse
}

// ReturnTrue returns true.
func ReturnTrue([]interface{}) interface{} {
	return types.JavaBoolTrue
}

// EofSet sets File EOF condition.
func EofSet(obj *object.Object, value bool) {
	obj.FieldTable[FileAtEOF] = object.Field{Ftype: types.Bool, Fvalue: value}
}

// EofGet gets File EOF boolean.
func EofGet(obj *object.Object) bool {
	value, ok := obj.FieldTable[FileAtEOF].Fvalue.(bool)
	if !ok {
		return false
	}
	return value
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

// ReturnRandomLong returns a random long.
func ReturnRandomLong([]interface{}) interface{} {
	// Generate random int64.
	var result int64
	byteArray := make([]byte, 8) // int64 is 8 bytes
	_, err := rand.Read(byteArray)
	if err != nil {
		trace.Warning(fmt.Sprintf("ReturnRandomLong: Failed to generate random int64: %v", err))
		return int64(42)
	}

	// Convert bytes to int64.
	for i := 0; i < 8; i++ {
		result = (result << 8) | int64(byteArray[i])
	}

	return result
}

// DefaultSecurityProvider is the single security provider for Jacobin.
// No other security providers are entertained.
var DefaultSecurityProvider *object.Object

// GetDefaultSecurityProvider returns the default security provider, initializing it if needed.
func GetDefaultSecurityProvider() *object.Object {
	return DefaultSecurityProvider
}

// getLinkedListFromObject (internal function) extracts the *list.List from the object
func GetLinkedListFromObject(self *object.Object) (*list.List, interface{}) {
	field, exists := self.FieldTable["value"]
	if !exists {
		return nil, GetGErrBlk(excNames.NullPointerException, "getLinkedListFromObject: LinkedList not initialized")
	}
	llst, ok := field.Fvalue.(*list.List)
	if !ok {
		return nil, GetGErrBlk(excNames.VirtualMachineError, "getLinkedListFromObject: Invalid LinkedList storage")
	}
	return llst, nil
}

// Invoke invokes a G function without setting up a frame. It is used primarily for
// calling constructors on libraries that are loaded as gfunctions.
func Invoke(whichFunc string, params []interface{}) interface{} {
	_, ret := MethodSignatures[whichFunc]
	if !ret {
		errMsg := fmt.Sprintf("Invoke: G function %s not found", whichFunc)
		exceptions.ThrowExNil(excNames.NoSuchMethodException, errMsg)
	}
	return MethodSignatures[whichFunc].GFunction(params)
}
