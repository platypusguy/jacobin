/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import (
	"encoding/binary"
	"fmt"
	"jacobin/excNames"
	"jacobin/exceptions"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/object"
	"jacobin/opcodes"
	"jacobin/types"
	"math"
	"runtime/debug"
	"unsafe"
)

// This file contains many support functions for the interpreter in run.go.
// These notably include push, pop, and peek operations on the operand stack,
// as well as some formatting functions for tracing, and utility functions for
// conversions of interfaces and data types.

// Convert a byte to an int64 by extending the sign-bit
func byteToInt64(bite byte) int64 {
	if (bite & 0x80) == 0x80 { // Negative bite value (left-most bit on)?
		// Negative byte - need to extend the sign (left-most) bit
		var wbytes = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00}
		wbytes[7] = bite
		// Form an int64 from the wbytes array
		// If you know C, this is equivalent to memcpy(&wint64, &wbytes, 8)
		return int64(binary.BigEndian.Uint64(wbytes))
	}

	// Not negative (left-most bit off) : just cast bite as an int64
	return int64(bite)
}

// converts four bytes into a signed 64-bit integer
func fourBytesToInt64(b1, b2, b3, b4 byte) int64 {
	wbytes := make([]byte, 8)
	wbytes[4] = b1
	wbytes[5] = b2
	wbytes[6] = b3
	wbytes[7] = b4

	if (b1 & 0x80) == 0x80 { // Negative bite value (left-most bit on)?
		// Negative byte - need to extend the sign (left-most) bit
		wbytes[0] = 0xff
		wbytes[1] = 0xff
		wbytes[2] = 0xff
		wbytes[3] = 0xff
	}
	return int64(binary.BigEndian.Uint64(wbytes))
}

// converts an interface{} value to int8. Used for BASTORE
func convertInterfaceToByte(val interface{}) byte {
	switch t := val.(type) {
	case int64:
		return byte(t)
	case int:
		return byte(t)
	case int8:
		return byte(t)
	case byte:
		return t
	}
	return 0
}

// converts an interface{} value into uint64
func convertInterfaceToUint64(val interface{}) uint64 {
	// in theory, the only types passed to this function are those
	// found on the operand stack: ints, floats, pointers
	switch t := val.(type) {
	case int64:
		return uint64(t)
	case float64:
		return uint64(math.Round(t))
	case unsafe.Pointer:
		intVal := uintptr(t)
		return uint64(intVal)
	}
	return 0
}

// Convert an interface{} consisting of some integral value to int64.
// Appears primarily in the runFrame{} IF* bytecodes.
func convertIntegralValueToInt64(arg interface{}) int64 {
	var value int64
	// fmt.Printf("convertIntegralValueToInt64 *TRAP*: argument type: %T\n", arg)
	switch arg.(type) {
	case int64:
		value = arg.(int64)
	case bool:
		if arg == true {
			value = int64(1)
		} else {
			value = int64(0)
		}
	case uint8:
		if arg == uint8(1) {
			value = int64(1)
		} else {
			value = int64(0)
		}
	default:
		errMsg := fmt.Sprintf("convertIntegralValueToInt64 *TRAP*: Invalid argument type: %T", arg)
		exceptions.ThrowEx(excNames.InvalidTypeException, errMsg, nil)
	}

	return value
}

// Log the existing stack
// Could be called for tracing -or- supply info for an error section
func logTraceStack(f *frames.Frame) {
	var traceInfo, output string
	if f.TOS == -1 {
		traceInfo = fmt.Sprintf("%55s %s.%s stack <empty>", "", f.ClName, f.MethName)
		_ = log.Log(traceInfo, log.WARNING)
		return
	}
	for ii := 0; ii <= f.TOS; ii++ {
		switch f.OpStack[ii].(type) {
		case *object.Object:
			if object.IsNull(f.OpStack[ii].(*object.Object)) {
				output = fmt.Sprintf("<null>")
			} else {
				objPtr := f.OpStack[ii].(*object.Object)
				output = objPtr.FormatField("")
			}
		case *[]uint8:
			value := f.OpStack[ii]
			strPtr := value.(*[]byte)
			str := string(*strPtr)
			output = fmt.Sprintf("*[]byte: %-10s", str)
		case []uint8:
			value := f.OpStack[ii]
			bytes := value.([]byte)
			str := string(bytes)
			output = fmt.Sprintf("[]byte: %-10s", str)
		default:
			output = fmt.Sprintf("%T %v ", f.OpStack[ii], f.OpStack[ii])
		}
		if f.TOS == ii {
			traceInfo = fmt.Sprintf("%55s %s.%s TOS   [%d] %s", "", f.ClName, f.MethName, ii, output)
		} else {
			traceInfo = fmt.Sprintf("%55s %s.%s stack [%d] %s", "", f.ClName, f.MethName, ii, output)
		}
		_ = log.Log(traceInfo, log.WARNING)
	}
}

// the generation and formatting of trace data for each executed bytecode.
// Returns the formatted data for output to logging, console, or other uses.
func emitTraceData(f *frames.Frame) string {
	var tos = " -"
	var stackTop = ""
	if f.TOS != -1 {
		tos = fmt.Sprintf("%2d", f.TOS)
		switch f.OpStack[f.TOS].(type) {
		// if the value at TOS is a string, say so and print the first 10 chars of the string
		case *object.Object:
			if object.IsNull(f.OpStack[f.TOS].(*object.Object)) {
				stackTop = fmt.Sprintf("<null>")
			} else {
				objPtr := f.OpStack[f.TOS].(*object.Object)
				stackTop = objPtr.FormatField("")
			}
		case *[]uint8:
			value := f.OpStack[f.TOS]
			strPtr := value.(*[]byte)
			str := string(*strPtr)
			stackTop = fmt.Sprintf("*[]byte: %-10s", str)
		case []uint8:
			value := f.OpStack[f.TOS]
			bytes := value.([]byte)
			str := string(bytes)
			stackTop = fmt.Sprintf("[]byte: %-10s", str)
		default:
			stackTop = fmt.Sprintf("%T %v ", f.OpStack[f.TOS], f.OpStack[f.TOS])
		}
	}

	traceInfo :=
		"class: " + fmt.Sprintf("%-22s", f.ClName) +
			" meth: " + fmt.Sprintf("%-10s", f.MethName) +
			" PC: " + fmt.Sprintf("% 3d", f.PC) +
			", " + fmt.Sprintf("%-13s", opcodes.BytecodeNames[int(f.Meth[f.PC])]) +
			" TOS: " + tos +
			" " + stackTop +
			" "
	return traceInfo
}

// traceObject : Used by push, pop, and peek in tracing an object.
func traceObject(f *frames.Frame, opStr string, obj *object.Object) {
	var traceInfo string
	prefix := fmt.Sprintf(" %4s          TOS:", opStr)

	// Nil pointer to object?
	if obj == nil {
		traceInfo = fmt.Sprintf("%74s", prefix) + fmt.Sprintf("%3d null", f.TOS)
		_ = log.Log(traceInfo, log.TRACE_INST)
		return
	}

	// The object pointer is not nil.
	klass := object.GoStringFromStringPoolIndex(obj.KlassName)
	traceInfo = fmt.Sprintf("%74s", prefix) + fmt.Sprintf("%3d, class: %s", f.TOS, klass)
	_ = log.Log(traceInfo, log.TRACE_INST)

	// Trace field table.
	prefix = " "
	if len(obj.FieldTable) > 0 {
		for fieldName := range obj.FieldTable {
			fld := obj.FieldTable[fieldName]
			if klass == types.StringClassName && fieldName == "value" {
				str := string(fld.Fvalue.([]byte))
				traceInfo = fmt.Sprintf("%74s", prefix) + fmt.Sprintf("field: %s %s %v \"%s\"", fieldName, fld.Ftype, fld.Fvalue, str)
			} else {
				traceInfo = fmt.Sprintf("%74s", prefix) + fmt.Sprintf("field: %s %s %v", fieldName, fld.Ftype, fld.Fvalue)
			}
			_ = log.Log(traceInfo, log.TRACE_INST)
		}
	} else { // nil FieldTable
		traceInfo = fmt.Sprintf("%74s", prefix) + fmt.Sprintf("no fields")
		_ = log.Log(traceInfo, log.TRACE_INST)
	}
}

// pop from the operand stack.
func pop(f *frames.Frame) interface{} {
	var value interface{}

	if f.TOS == -1 {
		glob := globals.GetGlobalRef()
		glob.ErrorGoStack = string(debug.Stack())
		exceptions.FormatStackUnderflowError(f)
		value = nil
	} else {
		value = f.OpStack[f.TOS]
	}

	// we show trace info of the TOS *before* we change its value--
	// all traces show TOS before the instruction is executed.
	if MainThread.Trace {
		var traceInfo string
		if f.TOS == -1 {
			traceInfo = fmt.Sprintf("%74s", "POP           TOS:  -")
			_ = log.Log(traceInfo, log.TRACE_INST)
		} else {
			if value == nil {
				traceInfo = fmt.Sprintf("%74s", "POP           TOS:") +
					fmt.Sprintf("%3d <nil>", f.TOS)
				_ = log.Log(traceInfo, log.TRACE_INST)
			} else {
				switch value.(type) {
				case *object.Object:
					obj := value.(*object.Object)
					traceObject(f, "POP", obj)
				case *[]uint8:
					strPtr := value.(*[]byte)
					str := string(*strPtr)
					traceInfo = fmt.Sprintf("%74s", "POP           TOS:") +
						fmt.Sprintf("%3d *[]byte: %-10s", f.TOS, str)
					_ = log.Log(traceInfo, log.TRACE_INST)
				case []uint8:
					bytes := value.([]byte)
					str := string(bytes)
					traceInfo = fmt.Sprintf("%74s", "POP           TOS:") +
						fmt.Sprintf("%3d []byte: %-10s", f.TOS, str)
					_ = log.Log(traceInfo, log.TRACE_INST)
				default:
					traceInfo = fmt.Sprintf("%74s", "POP           TOS:") +
						fmt.Sprintf("%3d %T %v", f.TOS, value, value)
					_ = log.Log(traceInfo, log.TRACE_INST)
				}
			}
		}
	}

	f.TOS -= 1 // adjust TOS
	if MainThread.Trace {
		logTraceStack(f)
	} // trace the resultant stack
	return value
}

// returns the value at the top of the stack without popping it off.
func peek(f *frames.Frame) interface{} {
	if f.TOS == -1 {
		glob := globals.GetGlobalRef()
		glob.ErrorGoStack = string(debug.Stack())
		exceptions.FormatStackUnderflowError(f)
		return nil
	}

	if MainThread.Trace {
		var traceInfo string
		value := f.OpStack[f.TOS]
		switch value.(type) {
		case *object.Object:
			obj := value.(*object.Object)
			traceObject(f, "PEEK", obj)
		default:
			traceInfo = fmt.Sprintf("                                                  "+
				"PEEK          TOS:%3d %T %v", f.TOS, value, value)
			_ = log.Log(traceInfo, log.TRACE_INST)
		}
	}
	if MainThread.Trace {
		logTraceStack(f)
	} // trace the stack
	return f.OpStack[f.TOS]
}

// push onto the operand stack
func push(f *frames.Frame, x interface{}) {
	if f.TOS == len(f.OpStack)-1 {
		// next step will set up error reporting and dump of frame stack
		// exceptions.FormatStackOverflowError(f)
		errMsg := fmt.Sprintf("in %s.%s, exceeded op stack size of %d",
			f.ClName, f.MethName, len(f.OpStack))
		_ = exceptions.ThrowEx(excNames.StackOverflowError, errMsg, f)
		return
	}
	// we show trace info of the TOS *before* we change its value--
	// all traces show TOS before the instruction is executed.
	if MainThread.Trace {
		var traceInfo string

		if f.TOS == -1 {
			traceInfo = fmt.Sprintf("%77s", "PUSH          TOS:  -")
			_ = log.Log(traceInfo, log.TRACE_INST)
		} else {
			if x == nil {
				traceInfo = fmt.Sprintf("%74s", "PUSH          TOS:") +
					fmt.Sprintf("%3d <nil>", f.TOS)
				_ = log.Log(traceInfo, log.TRACE_INST)
			} else {
				if x == object.Null {
					traceInfo = fmt.Sprintf("%74s", "PUSH          TOS:") +
						fmt.Sprintf("%3d null", f.TOS)
					_ = log.Log(traceInfo, log.TRACE_INST)
				} else {
					switch x.(type) {
					case *object.Object:
						obj := x.(*object.Object)
						traceObject(f, "PUSH", obj)
					case *[]uint8:
						strPtr := x.(*[]byte)
						str := string(*strPtr)
						traceInfo = fmt.Sprintf("%74s", "PUSH          TOS:") +
							fmt.Sprintf("%3d *[]byte: %-10s", f.TOS, str)
						_ = log.Log(traceInfo, log.TRACE_INST)
					case []uint8:
						bytes := x.([]byte)
						str := string(bytes)
						traceInfo = fmt.Sprintf("%74s", "PUSH          TOS:") +
							fmt.Sprintf("%3d []byte: %-10s", f.TOS, str)
						_ = log.Log(traceInfo, log.TRACE_INST)
					default:
						traceInfo = fmt.Sprintf("%56s", " ") +
							fmt.Sprintf("PUSH          TOS:%3d %T %v", f.TOS, x, x)
						_ = log.Log(traceInfo, log.TRACE_INST)
					}
				}
			}
		}
	}

	// the actual push
	f.TOS += 1
	f.OpStack[f.TOS] = x
	if MainThread.Trace {
		logTraceStack(f)
	} // trace the resultant stack

}
