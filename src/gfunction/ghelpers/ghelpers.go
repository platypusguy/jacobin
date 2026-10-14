/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package ghelpers

import (
	"container/list"
	"crypto/rand"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/trace"
	"jacobin/src/types"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
)

// Map repository of method signatures for all G functions:
var MethodSignatures = make(map[string]GMeth)
var TestMethodSignatures = make(map[string]GMeth) // used only for the test gfunctions
var TestGfunctionsLoaded = false

// GMeth is the entry in the MTable for Go functions. See MTable comments for details.
//   - ParamSlots - the number of user parameters in a G function. E.g. For atan2, this would be 2.
//   - GFunction - a go function. All go functions accept a possibly empty slice of any and
//     return an any which might be nil (E.g. Java void).
//   - NeedsContext - does this method need a pointer to the frame stack? Defaults to false.
type GMeth struct {
	ParamSlots   int
	GFunction    func([]any) any
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
var FileStatus = "status"     // using this value in case some member function is looking at it
var FilePath = "FilePath"     // full absolute path of a file aka canonical path
var FileHandle = "FileHandle" // *os.File
var FileMark = "FileMark"     // file position relative to beginning (0)
var FileAtEOF = "FileAtEOF"   // file at EOF

// File I/O constants:
var CreateFilePermissions os.FileMode = 0664 // When creating, read and write for user and group, others read-only

// Radix boundaries:
var MinRadix int64 = 2
var MaxRadix int64 = 36

// int64 value boundaries:
var MaxIntValue int64 = 2147483647
var MinIntValue int64 = -2147483648

// ClinitGeneric is a do-nothing Go function shared by several source files
func ClinitGeneric([]any) any {
	return nil
}

// JustReturn is a do-nothing Go function shared by several source files
func JustReturn([]any) any {
	return nil
}

// ReturnNull returns a Java null reference.
func ReturnNull([]any) any {
	return object.Null
}

func ReturnCharsetName([]any) any {
	return object.StringObjectFromGoString(globals.GetCharsetName())
}

// ReturnFalse returns false.
func ReturnFalse([]any) any {
	return types.JavaBoolFalse
}

// ReturnTrue returns true.
func ReturnTrue([]any) any {
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
	fldValue := object.Field{Ftype: types.Ref, Fvalue: ptrBigInt}
	obj.FieldTable["value"] = fldValue
	var fldSign object.Field
	switch {
	case argValue == 0:
		fldSign = object.Field{Ftype: types.Ref, Fvalue: int64(0)}
	case argValue < 0:
		fldSign = object.Field{Ftype: types.Ref, Fvalue: int64(-1)}
	default:
		fldSign = object.Field{Ftype: types.Ref, Fvalue: int64(+1)}
	}
	obj.FieldTable["signum"] = fldSign
}

// ReturnRandomLong returns a random long.
func ReturnRandomLong([]any) any {
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
func GetLinkedListFromObject(self *object.Object) (*list.List, any) {
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
func Invoke(whichFunc string, params []any) any {
	_, ret := MethodSignatures[whichFunc]
	if !ret {
		errMsg := fmt.Sprintf("Invoke: G function %s not found", whichFunc)
		//exceptions.ThrowExNil(excNames.NoSuchMethodException, errMsg)
		return GetGErrBlk(excNames.NoSuchMethodException, errMsg)
	}
	return MethodSignatures[whichFunc].GFunction(params)
}

// Converts a variable number of arguments into a []any that can be used in gfunctions
func ConvertArgsToParams(args ...any) []any {
	if len(args) == 0 {
		return make([]any, 0)
	}

	params := make([]any, len(args))
	for i, arg := range args {
		params[i] = arg
	}
	return params
}

// JavaFloatingPointToString reproduces Float.toString / Double.toString.
// bitSize is 32 for float and 64 for double.
func JavaFloatingPointToString(value float64, bitSize int) string {
	switch {
	case math.IsNaN(value):
		return "NaN"
	case math.IsInf(value, 1):
		return "Infinity"
	case math.IsInf(value, -1):
		return "-Infinity"
	case value == 0:
		if math.Signbit(value) {
			return "-0.0"
		}
		return "0.0"
	}

	// Shortest round-trip digits, e.g. "-1.2345e+10".
	scientific := strconv.FormatFloat(value, 'e', -1, bitSize)
	sign := ""
	if scientific[0] == '-' {
		sign = "-"
		scientific = scientific[1:]
	}
	mantissa, exponentText, _ := strings.Cut(scientific, "e")
	exponent, _ := strconv.Atoi(exponentText)
	digits := strings.Replace(mantissa, ".", "", 1)

	magnitude := math.Abs(value)
	if magnitude >= 1e-3 && magnitude < 1e7 {
		// Plain decimal notation with at least one digit after the point.
		if exponent >= 0 {
			integerLength := exponent + 1
			if len(digits) <= integerLength {
				return fmt.Sprintf("%s%s%s.0", sign, digits, strings.Repeat("0", integerLength-len(digits)))
			}
			return fmt.Sprintf("%s%s.%s", sign, digits[:integerLength], digits[integerLength:])
		}
		return fmt.Sprintf("%s0.%s%s", sign, strings.Repeat("0", -exponent-1), digits)
	}

	// Computerized scientific notation: d.dddE<exponent>.
	fraction := digits[1:]
	if fraction == "" {
		fraction = "0"
	}
	return fmt.Sprintf("%s%s.%sE%d", sign, digits[:1], fraction, exponent)
}

func WrongArgType(funcName, expected string, observed any) any {
	errMsg := fmt.Sprintf("%s: expected %s, observed: %T", funcName, expected, observed)
	return GetGErrBlk(excNames.IllegalArgumentException, errMsg)
}
