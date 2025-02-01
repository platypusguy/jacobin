/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-5 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package trace

// The principal logging function. Note it currently logs to stderr.
// At some future point, might allow the user to specify where logging should go.
import (
	"fmt"
	"jacobin/excNames"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/opcodes"
	"jacobin/types"
	"os"
	"sync"
	"time"
)

// Mutex for protecting the Log function during multithreading.
var mutex = sync.Mutex{}

// StartTime is the start time of this instance of the Jacoby VM.
var StartTime time.Time

// Identical to shutdown.UNKNOWN_ERROR (avoiding a cycle)
const UNKNOWN_ERROR = 5

// Initialize the trace frame.
func Init() {
	StartTime = time.Now()
}

// Trace is the principal tracing function. Note that it currently
// writes to stderr. At some future point, this might become an option.
func Trace(argMsg string) {

	var err error

	// if the message is more low-level than a WARNING,
	// prefix it with the elapsed time in millisecs.
	// check duration accuracy: time.Sleep(100 * time.Millisecond)
	duration := time.Since(StartTime)
	var millis = duration.Milliseconds()

	// Lock access to the logging stream to prevent inter-thread overwrite issues
	mutex.Lock()
	_, err = fmt.Fprintf(os.Stderr, "[%3d.%03ds] %s\n", millis/1000, millis%1000, argMsg)
	mutex.Unlock()
	if err != nil {
		errMsg := fmt.Sprintf("Trace: *** stderr failed, err: %v", err)
		rawAbort(excNames.IOError, errMsg)
	}
}

// An error message is a prefix-decorated message that has no time-stamp.
func Error(argMsg string) {
	var err error
	errMsg := "ERROR: " + argMsg
	mutex.Lock()
	_, err = fmt.Fprintf(os.Stderr, "%s\n", errMsg)
	mutex.Unlock()
	if err != nil {
		errMsg = fmt.Sprintf("Error: *** stderr failed, err: %v", err)
		rawAbort(excNames.IOError, errMsg)
	}
}

// Similar to Error, except it's a warning, not an error.
func Warning(argMsg string) {
	errMsg := "WARNING: " + argMsg
	mutex.Lock()
	_, err := fmt.Fprintf(os.Stderr, "%s\n", errMsg)
	mutex.Unlock()
	if err != nil {
		errMsg = fmt.Sprintf("Error: *** stderr failed, err: %v", err)
		rawAbort(excNames.IOError, errMsg)
	}
}

// the generation and formatting of trace data for each executed bytecode.
// Returns the formatted data for output to logging, console, or other uses.
func EmitTraceData(f *frames.Frame) string {
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
		case []types.JavaByte:
			value := f.OpStack[f.TOS]
			bytes := value.([]types.JavaByte)
			str := object.GoStringFromJavaByteArray(bytes)
			stackTop = fmt.Sprintf("[]javaByte: %-10s", str)
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

// Generate a trace of a field ID (static or non-static).
func EmitTraceFieldID(opcode, fld string) {
	traceInfo := fmt.Sprintf("%65s fieldName: %s", opcode, fld)
	Trace(traceInfo)
}

// Log the existing stack
// Could be called for tracing -or- supply info for an error section
func LogTraceStack(f *frames.Frame) {
	var traceInfo, output string
	if f.TOS == -1 {
		traceInfo = fmt.Sprintf("%55s %s.%s stack <empty>", "", f.ClName, f.MethName)
		Trace(traceInfo)
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
		case []types.JavaByte:
			value := f.OpStack[ii]
			bytes := value.([]types.JavaByte)
			str := object.GoStringFromJavaByteArray(bytes)
			output = fmt.Sprintf("[]javaByte: %-10s", str)
		default:
			output = fmt.Sprintf("%T %v ", f.OpStack[ii], f.OpStack[ii])
		}
		if f.TOS == ii {
			traceInfo = fmt.Sprintf("%55s %s.%s TOS   [%d] %s", "", f.ClName, f.MethName, ii, output)
		} else {
			traceInfo = fmt.Sprintf("%55s %s.%s stack [%d] %s", "", f.ClName, f.MethName, ii, output)
		}
		Trace(traceInfo)
	}
}

// TraceObject : Used by push, pop, and peek in tracing an object.
func TraceObject(f *frames.Frame, opStr string, obj *object.Object) {
	var traceInfo string
	prefix := fmt.Sprintf(" %4s          TOS:", opStr)

	// Nil pointer to object?
	if obj == nil {
		traceInfo = fmt.Sprintf("%74s", prefix) + fmt.Sprintf("%3d null", f.TOS)
		Trace(traceInfo)
		return
	}

	// The object pointer is not nil.
	klass := object.GoStringFromStringPoolIndex(obj.KlassName)
	traceInfo = fmt.Sprintf("%74s", prefix) + fmt.Sprintf("%3d, class: %s", f.TOS, klass)
	Trace(traceInfo)

	// Trace field table.
	prefix = " "
	if len(obj.FieldTable) > 0 {
		for fieldName := range obj.FieldTable {
			fld := obj.FieldTable[fieldName]
			if klass == types.StringClassName && fieldName == "value" {
				var str string
				switch fld.Fvalue.(type) {
				case []types.JavaByte:
					str = object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
				default:
					str = string(fld.Fvalue.([]byte))
				}

				traceInfo = fmt.Sprintf("%74s", prefix) + fmt.Sprintf("field: %s %s %v \"%s\"", fieldName, fld.Ftype, fld.Fvalue, str)
			} else {
				traceInfo = fmt.Sprintf("%74s", prefix) + fmt.Sprintf("field: %s %s %v", fieldName, fld.Ftype, fld.Fvalue)
			}
			Trace(traceInfo)
		}
	} else { // nil FieldTable
		traceInfo = fmt.Sprintf("%74s", prefix) + fmt.Sprintf("no fields")
		Trace(traceInfo)
	}
}

// Perform a minimal abort, which is a direct call to the global minimal abort function.
// Clearly, if trace is not working, then something is grievously wrong and the abort
// must be immediate.
func rawAbort(whichException int, msg string) {
	globals.GetGlobalRef().FuncMinimalAbort(whichException, msg)
}
