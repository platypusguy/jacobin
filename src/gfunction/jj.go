/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"bytes"
	"fmt"
	"jacobin/src/excNames"
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

	MethodSignatures["jj._dumpStatics(Ljava/lang/String;ILjava/lang/String;)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  jjDumpStatics,
		}

	MethodSignatures["jj._dumpObject(Ljava/lang/Object;Ljava/lang/String;I)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  jjDumpObject,
		}

	MethodSignatures["jj._getStaticString(Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  jjGetStaticString,
		}

	MethodSignatures["jj._getFieldString(Ljava/lang/Object;Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  jjGetFieldString,
		}

	MethodSignatures["jj._subProcess(LjjSubProcessObject;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  jjSubProcess,
		}

	MethodSignatures["jj._getProgramName()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  jjGetProgramName,
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
	if len(params) < 0 || params[1] == nil {
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
	static := statics.Statics[className+"."+fieldName]

	// Handle vectors.
	if strings.HasPrefix(static.Type, types.Array) {
		return jjStringifyVector(static.Value.(*object.Object))
	}

	// Handle a scalar.
	return jjStringifyScalar(static.Type, static.Value)
}

func jjGetFieldString(params []interface{}) interface{} {

	// Get this object.
	thisObj := params[0].(*object.Object)

	// Get field name.
	if len(params) < 0 || params[1] == nil {
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
	if fld.Ftype == "Ljava/lang/String;" {
		return object.StringObjectFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
	}

	// Handle vectors.
	if strings.HasPrefix(fld.Ftype, types.Array) {
		return jjStringifyVector(fld.Fvalue)
	}

	// Handle a scalar.
	return jjStringifyScalar(fld.Ftype, fld.Fvalue)
}

func jjDumpStatics(params []interface{}) interface{} {
	fromObj := params[0].(*object.Object)
	if fromObj == nil || fromObj.KlassName == types.InvalidStringIndex {
		errMsg := fmt.Sprintf("jjDumpStatics: Invalid from object: %T", params[0])
		return object.StringObjectFromGoString(errMsg)
	}
	from := object.ObjectFieldToString(fromObj, "value")
	selection := params[1].(int64)
	classNameObj := params[2].(*object.Object)
	className := object.ObjectFieldToString(classNameObj, "value")

	statics.DumpStatics(from, selection, className)
	return nil
}

func jjDumpObject(params []interface{}) interface{} {
	this := params[0].(*object.Object)
	objTitle := params[1].(*object.Object)
	title := object.ObjectFieldToString(objTitle, "value")
	indent := params[2].(int64)
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
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Replace the CLASSPATH environment variable if the classpath field is non-empty.
	cpArray, ok := subpObj.FieldTable["classpath"].Fvalue.([]*object.Object)
	if !ok {
		errMsg := "jjSubProcess: Missing/Misformatted subprocess classpath field"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
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

	// Build command line.
	objArray, ok := subpObj.FieldTable["commandLine"].Fvalue.([]*object.Object)
	if !ok {
		errMsg := "jjSubProcess: Missing/Misformatted subprocess commandLine field"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	strArray := object.GoStringArrayFromStringObjectArray(objArray)
	switch len(strArray) {
	case 0:
		errMsg := "jjSubProcess: Nil subprocess commandLine field"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	case 1:
		cmd = exec.Command(strArray[0])
	default:
		cmd = exec.Command(strArray[0], strArray[1:]...)
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
		stderr = fmt.Sprintf("jjSubProcess: Process %s failed, err: %s", cmdString, err.Error())
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
