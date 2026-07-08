/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package misc

import (
	"bytes"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/trace"
	"jacobin/src/types"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

// jj (Jacobin JVM) functions are functions that can be inserted inside Java programs
// for diagnostic purposes. They simply return when run in the JDK, but do what they're
// supposed to do when run under Jacobin.
//
// Note this is a rough first design that will surely be refined. (JACOBIN-624)

func Load_jj() {

	ghelpers.MethodSignatures["jj._dumpStatics(Ljava/lang/String;ILjava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  jjDumpStatics,
		}

	ghelpers.MethodSignatures["jj._dumpObject(Ljava/lang/Object;Ljava/lang/String;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  jjDumpObject,
		}

	ghelpers.MethodSignatures["jj._getStaticString(Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  jjGetStaticString,
		}

	ghelpers.MethodSignatures["jj._getFieldString(Ljava/lang/Object;Ljava/lang/String;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  jjGetFieldString,
		}

	ghelpers.MethodSignatures["jj._subProcess(LjjSubProcessObject;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  jjSubProcess,
		}

	ghelpers.MethodSignatures["jj._getProgramName()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  jjGetProgramName,
		}

	ghelpers.MethodSignatures["jj._panic()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  jjPanic,
		}

	ghelpers.MethodSignatures["jj._traceInst(Z)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  jjTraceInst,
		}

	ghelpers.MethodSignatures["jj._traceVerbose(Z)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  jjTraceVerbose,
		}
}

func jjStringifyScalar(ftype string, fvalue any) *object.Object {
	var str string
	switch ftype {
	case types.Bool:
		if fvalue.(int64) == 1 {
			str = "true"
		} else {
			str = "false"
		}
	case types.Byte: // uint8, int8
		switch fvalue.(type) {
		case byte:
			str = fmt.Sprintf("0x%02x", fvalue.(byte))
		case types.JavaByte:
			str = fmt.Sprintf("0x%02x", fvalue.(types.JavaByte))
		case int64:
			str = fmt.Sprintf("0x%02x", fvalue.(int64))
		default:
			str = fmt.Sprintf("%v", fvalue)
		}
	case types.Char, types.Rune:
		str = fmt.Sprintf("%c", fvalue.(int64))
	case types.Double:
		str = strconv.FormatFloat(fvalue.(float64), 'g', -1, 64)
	case types.Float:
		str = strconv.FormatFloat(float64(fvalue.(float64)), 'g', -1, 64)
	case types.Int:
		str = fmt.Sprintf("%d", fvalue.(int64))
	case types.Long:
		str = fmt.Sprintf("%d", fvalue.(int64))
	case "Ljava/lang/String;":
		str = object.GoStringFromStringObject(fvalue.(*object.Object))
	case types.Short:
		str = fmt.Sprintf("%d", fvalue.(int64))
	case types.Ref, types.ByteArray:
		if object.IsNull(fvalue.(*object.Object)) {
			str = types.NullString
		} else {
			obj := fvalue.(*object.Object)
			if obj.KlassName == types.StringPoolStringIndex {
				// It is a Java String object. Return it as-is.
				return obj
			}
			// Not a Java String object.
			str = fmt.Sprintf("%v", fvalue)
		}
	default:
		str = fmt.Sprintf("%v", fvalue)
	}
	return object.StringObjectFromGoString(str)
}

func jjStringifyVector(thing any) *object.Object {
	var result string = ""
	var anArray reflect.Value
	switch thing.(type) {
	case *object.Object:
		anArray = reflect.ValueOf(thing.(*object.Object).FieldTable["value"].Fvalue)
	default:
		anArray = reflect.ValueOf(thing)
	}
	for ix := 0; ix < anArray.Len(); ix++ {
		if ix > 0 {
			result += "," // comma as a separator between elements
		}
		element := anArray.Index(ix).Interface() // Get the element as an interface{}
		result += fmt.Sprintf("%v", element)
	}
	return object.StringObjectFromGoString(result)
}

func jjGetStaticString(params []interface{}) interface{} {

	if len(params) == 0 || params[0] == nil {
		errMsg := fmt.Sprintf("jjGetStaticString: No class object")
		return object.StringObjectFromGoString(errMsg)
	}

	// Get class name.
	classObj := params[0].(*object.Object)
	if classObj == nil || classObj.KlassName == types.InvalidStringIndex {
		errMsg := fmt.Sprintf("jjGetStaticString: Invalid class object: %T", params[0])
		return object.StringObjectFromGoString(errMsg)
	}
	className := object.ObjectFieldToString(classObj, "value")

	// Get field name.
	if len(params) < 2 || params[1] == nil {
		errMsg := fmt.Sprintf("jjGetStaticString: Invalid field is missing or nil")
		return object.StringObjectFromGoString(errMsg)
	}
	fieldObj := params[1].(*object.Object)
	if fieldObj == nil || fieldObj.KlassName == types.InvalidStringIndex {
		errMsg := fmt.Sprintf("jjGetStaticString: Invalid field object: %T", params[1])
		return object.StringObjectFromGoString(errMsg)
	}
	fieldName := object.ObjectFieldToString(fieldObj, "value")

	// Convert statics entry to a string object.
	sme, ok := statics.QueryStatic(className, fieldName)
	if !ok {
		errMsg := fmt.Sprintf("jjGetStaticString: statics.QueryStatic(%s, %s) failed", className, fieldName)
		return object.StringObjectFromGoString(errMsg)
	}

	return object.StringifyAnythingJava(object.Field{
		Ftype:  sme.Type,
		Fvalue: sme.Value,
	})
}

func jjGetFieldString(params []interface{}) interface{} {

	// Get this object.
	thisObj := params[0].(*object.Object)

	// Get field name.
	if len(params) < 2 || params[1] == nil {
		errMsg := fmt.Sprintf("jjGetFieldString: Invalid field is missing or nil")
		return object.StringObjectFromGoString(errMsg)
	}

	fieldObj := params[1].(*object.Object)
	if fieldObj == nil || fieldObj.KlassName == types.InvalidStringIndex {
		errMsg := fmt.Sprintf("jjGetFieldString: Invalid field object: %T", params[1])
		return object.StringObjectFromGoString(errMsg)
	}
	fieldName := object.ObjectFieldToString(fieldObj, "value")

	// Convert field entry to a string object.
	fld, ok := thisObj.FieldTable[fieldName]
	if !ok {
		errMsg := fmt.Sprintf("jjGetFieldString: No such field name: %s", fieldName)
		return object.StringObjectFromGoString(errMsg)
	}

	return object.StringifyAnythingJava(fld)
}

func jjDumpStatics(params []interface{}) interface{} {
	if len(params) < 1 || params[0] == nil {
		errMsg := "jjDumpStatics: Missing from object"
		return object.StringObjectFromGoString(errMsg)
	}
	fromObj, ok := params[0].(*object.Object)
	if !ok || fromObj == nil || fromObj.KlassName == types.InvalidStringIndex {
		errMsg := fmt.Sprintf("jjDumpStatics: Invalid from object: %T", params[0])
		return object.StringObjectFromGoString(errMsg)
	}
	from := object.ObjectFieldToString(fromObj, "value")

	if len(params) < 2 || params[1] == nil {
		errMsg := "jjDumpStatics: Missing selection"
		return object.StringObjectFromGoString(errMsg)
	}
	selection, ok := params[1].(int64)
	if !ok {
		errMsg := "jjDumpStatics: Invalid selection"
		return object.StringObjectFromGoString(errMsg)
	}

	if len(params) < 3 || params[2] == nil {
		errMsg := "jjDumpStatics: Missing className object"
		return object.StringObjectFromGoString(errMsg)
	}
	classNameObj, ok := params[2].(*object.Object)
	if !ok || classNameObj == nil || classNameObj.KlassName == types.InvalidStringIndex {
		errMsg := fmt.Sprintf("jjDumpStatics: Invalid className object: %T", params[2])
		return object.StringObjectFromGoString(errMsg)
	}
	className := object.ObjectFieldToString(classNameObj, "value")

	statics.DumpStatics(from, selection, className)
	return nil
}

func jjDumpObject(params []interface{}) interface{} {
	if len(params) < 1 || params[0] == nil {
		trace.Error("jjDumpObject: Missing object")
		return nil
	}
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		trace.Error("jjDumpObject: Invalid object")
		return nil
	}

	if len(params) < 2 || params[1] == nil {
		trace.Error("jjDumpObject: Missing title")
		return nil
	}
	objTitle, ok := params[1].(*object.Object)
	if !ok || objTitle == nil {
		trace.Error("jjDumpObject: Invalid title")
		return nil
	}
	title := object.ObjectFieldToString(objTitle, "value")

	if len(params) < 3 || params[2] == nil {
		trace.Error("jjDumpObject: Missing indent")
		return nil
	}
	indent, ok := params[2].(int64)
	if !ok {
		trace.Error("jjDumpObject: Invalid indent")
		return nil
	}

	this.DumpObject(title, int(indent))
	return nil
}

/*
params[0] *object.Object:

	class jjSubProcessObject {
		String[] commandLine; // input
		String[] classpath; // input; empty means use existing
		String stdout; // output
		String stderr; // output
	}

The returned int64 indicates the sub-process exit code.
*/
func jjSubProcess(params []interface{}) interface{} {

	// Subprocess execution handle.
	var cmd *exec.Cmd

	// params[0] should have subprocess object.
	subpObj, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "jjSubProcess: Missing/Misformatted subprocess object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Replace the CLASSPATH environment variable if the classpath field is non-empty.
	var cpArray []*object.Object
	if cpField, ok := subpObj.FieldTable["classpath"]; ok {
		cpArray, _ = cpField.Fvalue.([]*object.Object)
	}

	// Build command line.
	objArray, ok := subpObj.FieldTable["commandLine"].Fvalue.([]*object.Object)
	if !ok {
		errMsg := "jjSubProcess: Missing/Misformatted subprocess commandLine field"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	strArray := object.GoStringArrayFromStringObjectArray(objArray)
	if len(strArray) == 0 {
		errMsg := "jjSubProcess: Nil subprocess commandLine field"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	cmd = exec.Command(strArray[0], strArray[1:]...)

	if len(cpArray) > 0 {
		cpStrArray := object.GoStringArrayFromStringObjectArray(cpArray)

		// Join with platform-specific separator.
		sep := ":"
		if runtime.GOOS == "windows" {
			sep = ";"
		}
		classpath := strings.Join(cpStrArray, sep)

		// Clone environment and replace any existing CLASSPATH entry.
		env := os.Environ()
		newEnv := make([]string, 0, len(env)+1)
		for _, kv := range env {
			if !strings.HasPrefix(kv, "CLASSPATH=") {
				newEnv = append(newEnv, kv)
			}
		}
		newEnv = append(newEnv, "CLASSPATH="+classpath)

		// Update subprocess environment with new classpath.
		cmd.Env = newEnv
	}

	// Buffers to capture stdout and stderr
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Run the command and wait for it to finish
	err := cmd.Run()

	// Collect output from pipes.
	stdout := stdoutBuf.String()
	stderr := stderrBuf.String()
	subpObj.FieldTable["stdout"] = object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(stdout)}
	subpObj.FieldTable["stderr"] = object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(stderr)}

	// Handle exit code or POSIX signal.
	exitCode := int64(0)
	var signal syscall.Signal
	if err != nil {
		// Something went wrong. Indicate err.Error() in stderr.
		cmdString := strings.Join(strArray, " ")
		stderrLines := strings.Split(stderr, "\n")
		stderr = fmt.Sprintf("jjSubProcess: Process %s failed, err: %s\nstderr: %s", cmdString, err.Error(), stderrLines[0])
		trace.Error(stderr)
		subpObj.FieldTable["stderr"] = object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(stderr)}

		// Is err of type *exec.ExitError?
		if exitErr, ok := err.(*exec.ExitError); ok {
			// The subprocess exited with a non-zero exit status.
			// POSIX?
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				// POSIX system.
				// Signalled?
				if status.Signaled() {
					signal = status.Signal()
					return int64(signal)
				} else {
					// Not signalled.
					return int64(status.ExitStatus())
				}
			} else {
				// Not a POSIX system (e.g. Windows).
				return int64(exitErr.ExitCode())
			}
		} else {
			// The err is not of type *exec.ExitError.
			// Command probably did not even start (E.g. file not found).
			// stderr should already have been set with an error message.
			return int64(-1)
		}
	}

	return exitCode
}

func jjGetProgramName([]interface{}) interface{} {
	glob := globals.GetGlobalRef()
	str := glob.JacobinName
	return object.StringObjectFromGoString(str)
}

func jjPanic([]interface{}) interface{} {
	trace.Warning("jjPanic: Will cause a Go runtime divide by zero panic")
	var zero = 0
	zero = 1 / zero
	errMsg := "jjPanic: What??? No splash???"
	return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
}

func jjTraceInst(params []interface{}) interface{} {
	flag := params[0].(types.JavaBool)
	if flag == types.JavaBoolTrue {
		globals.TraceInst = true
		trace.Trace("jjTraceInst: begin")
	} else {
		globals.TraceInst = false
		trace.Trace("jjTraceInst: end")
	}
	return nil
}

func jjTraceVerbose(params []interface{}) interface{} {
	flag := params[0].(types.JavaBool)
	if flag == types.JavaBoolTrue {
		globals.TraceVerbose = true
		trace.Trace("jjTraceVerbose: begin")
	} else {
		globals.TraceVerbose = false
		trace.Trace("jjTraceVerbose: end")
	}
	return nil
}
