/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"os/exec"
	"strings"
)

func Load_Lang_Process_Builder() {

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.<init>(Ljava/util/List;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  processBuilderInitList,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.<init>([Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  processBuilderInitStringArray,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.command(Ljava/util/List;)Ljava/lang/ProcessBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.command([Ljava/lang/String;)Ljava/lang/ProcessBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  processBuilderCommandArray,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.directory()Ljava/io/File;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.directory(Ljava/io/File;)Ljava/lang/ProcessBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.environment()Ljava/util/Map;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  processBuilderEnvironment,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.inheritIO()Ljava/lang/ProcessBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.redirectError()Ljava/lang/ProcessBuilder$Redirect;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.redirectError(Ljava/lang/ProcessBuilder$Redirect;)Ljava/lang/ProcessBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.redirectErrorStream()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  processBuilderRedirectErrorStreamGet,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.redirectErrorStream(Z)Ljava/lang/ProcessBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  processBuilderRedirectErrorStreamSet,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.redirectInput()Ljava/lang/ProcessBuilder$Redirect;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.redirectInput(Ljava/io/File;)Ljava/lang/ProcessBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.redirectInput(Ljava/lang/ProcessBuilder$Redirect;)Ljava/lang/ProcessBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.redirectOutput()Ljava/lang/ProcessBuilder$Redirect;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.redirectOutput(Ljava/io/File;)Ljava/lang/ProcessBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.redirectOutput(Ljava/lang/ProcessBuilder$Redirect;)Ljava/lang/ProcessBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.start()Ljava/lang/Process;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  processBuilderStart,
		}

	ghelpers.MethodSignatures["java/lang/ProcessBuilder.startPipeline(Ljava/util/List;)Ljava/util/List;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}
}

// "java/lang/ProcessBuilder.<init>(Ljava/util/List;)V"
func processBuilderInitList(params []any) any {
	res := processBuilderCommandList(params)
	if errBlk, ok := res.(*ghelpers.GErrBlk); ok {
		return errBlk
	}
	return nil
}

// "java/lang/ProcessBuilder.<init>([Ljava/lang/String;)V"
func processBuilderInitStringArray(params []any) any {
	res := processBuilderCommandArray(params)
	if errBlk, ok := res.(*ghelpers.GErrBlk); ok {
		return errBlk
	}
	return nil
}

// "java/lang/ProcessBuilder.command(Ljava/util/List;)Ljava/lang/ProcessBuilder;"
func processBuilderCommandList(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderCommandList: self is null")
	}

	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "processBuilderCommandList: self is not an object")
	}

	if len(params) < 2 || params[1] == nil || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderCommandList: command list is null")
	}

	var listElements []any
	switch arg := params[1].(type) {
	case *object.Object:
		if fld, ok := arg.FieldTable["value"]; ok {
			switch elems := fld.Fvalue.(type) {
			case []*object.Object:
				listElements = make([]any, len(elems))
				for i, e := range elems {
					if e == nil || object.IsNull(e) {
						return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderCommandList: command element is null")
					}
					listElements[i] = e
				}
			case []any:
				listElements = make([]any, len(elems))
				for i, e := range elems {
					if e == nil || object.IsNull(e) {
						return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderCommandList: command element is null")
					}
					listElements[i] = e
				}
			case []string:
				listElements = make([]any, len(elems))
				for i, s := range elems {
					listElements[i] = object.StringObjectFromGoString(s)
				}
			default:
				return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("processBuilderCommandList: unsupported list value type %T", fld.Fvalue))
			}
		} else {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "processBuilderCommandList: missing list value field")
		}
	case []*object.Object:
		listElements = make([]any, len(arg))
		for i, e := range arg {
			if e == nil || object.IsNull(e) {
				return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderCommandList: command element is null")
			}
			listElements[i] = e
		}
	case []any:
		listElements = make([]any, len(arg))
		for i, e := range arg {
			if e == nil || object.IsNull(e) {
				return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderCommandList: command element is null")
			}
			listElements[i] = e
		}
	case []string:
		listElements = make([]any, len(arg))
		for i, s := range arg {
			listElements[i] = object.StringObjectFromGoString(s)
		}
	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("processBuilderCommandList: unsupported command argument type %T", params[1]))
	}

	listObj := object.MakePrimitiveObject("java/util/ArrayList", types.Ref, listElements)

	self.ThMutex.Lock()
	if self.FieldTable == nil {
		self.FieldTable = make(map[string]object.Field)
	}
	self.FieldTable["command"] = object.Field{
		Ftype:  "java/util/List",
		Fvalue: listObj,
	}
	self.ThMutex.Unlock()

	return self
}

// "java/lang/ProcessBuilder.command([Ljava/lang/String;)Ljava/lang/ProcessBuilder;"
func processBuilderCommandArray(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderCommandArray: self is null")
	}

	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "processBuilderCommandArray: self is not an object")
	}

	if len(params) < 2 || params[1] == nil || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderCommandArray: command array is null")
	}

	var listElements []any
	switch arg := params[1].(type) {
	case *object.Object:
		if fld, ok := arg.FieldTable["value"]; ok {
			switch elems := fld.Fvalue.(type) {
			case []*object.Object:
				listElements = make([]any, len(elems))
				for i, e := range elems {
					if e == nil || object.IsNull(e) {
						return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderCommandArray: command element is null")
					}
					listElements[i] = e
				}
			case []any:
				listElements = make([]any, len(elems))
				for i, e := range elems {
					if e == nil || object.IsNull(e) {
						return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderCommandArray: command element is null")
					}
					listElements[i] = e
				}
			case []string:
				listElements = make([]any, len(elems))
				for i, s := range elems {
					listElements[i] = object.StringObjectFromGoString(s)
				}
			default:
				return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("processBuilderCommandArray: unsupported array value type %T", fld.Fvalue))
			}
		} else {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "processBuilderCommandArray: missing array value field")
		}
	case []*object.Object:
		listElements = make([]any, len(arg))
		for i, e := range arg {
			if e == nil || object.IsNull(e) {
				return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderCommandArray: command element is null")
			}
			listElements[i] = e
		}
	case []any:
		listElements = make([]any, len(arg))
		for i, e := range arg {
			if e == nil || object.IsNull(e) {
				return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderCommandArray: command element is null")
			}
			listElements[i] = e
		}
	case []string:
		listElements = make([]any, len(arg))
		for i, s := range arg {
			listElements[i] = object.StringObjectFromGoString(s)
		}
	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("processBuilderCommandArray: unsupported command argument type %T", params[1]))
	}

	listObj := object.MakePrimitiveObject("java/util/ArrayList", types.Ref, listElements)

	self.ThMutex.Lock()
	if self.FieldTable == nil {
		self.FieldTable = make(map[string]object.Field)
	}
	self.FieldTable["command"] = object.Field{
		Ftype:  "java/util/List",
		Fvalue: listObj,
	}
	self.ThMutex.Unlock()

	return self
}

// "java/lang/ProcessBuilder.environment()Ljava/util/Map;"
func processBuilderEnvironment(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderEnvironment: self is null")
	}

	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "processBuilderEnvironment: self is not an object")
	}

	self.ThMutex.Lock()
	defer self.ThMutex.Unlock()

	if self.FieldTable == nil {
		self.FieldTable = make(map[string]object.Field)
	}

	if envFld, exists := self.FieldTable["environment"]; exists && envFld.Fvalue != nil && !object.IsNull(envFld.Fvalue) {
		return envFld.Fvalue
	}

	envMap := make(types.DefHashMap)
	for _, env := range os.Environ() {
		k, v, found := strings.Cut(env, "=")
		if found {
			envMap[k] = object.StringObjectFromGoString(v)
		}
	}

	classNameHashMap := "java/util/HashMap"
	mapObj := object.MakeEmptyObjectWithClassName(&classNameHashMap)
	mapObj.ThMutex.Lock()
	mapObj.FieldTable["map"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: envMap,
	}
	mapObj.ThMutex.Unlock()

	self.FieldTable["environment"] = object.Field{
		Ftype:  "java/util/Map",
		Fvalue: mapObj,
	}

	return mapObj
}

// "java/lang/ProcessBuilder.redirectErrorStream()Z"
func processBuilderRedirectErrorStreamGet(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderRedirectErrorStream: self is null")
	}

	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "processBuilderRedirectErrorStream: self is not an object")
	}

	self.ThMutex.RLock()
	defer self.ThMutex.RUnlock()

	if self.FieldTable != nil {
		if fld, exists := self.FieldTable["redirectErrorStream"]; exists && fld.Fvalue != nil && !object.IsNull(fld.Fvalue) {
			switch v := fld.Fvalue.(type) {
			case bool:
				if v {
					return types.JavaBoolTrue
				}
				return types.JavaBoolFalse
			case int64:
				if v != 0 {
					return types.JavaBoolTrue
				}
				return types.JavaBoolFalse
			case int:
				if v != 0 {
					return types.JavaBoolTrue
				}
				return types.JavaBoolFalse
			}
		}
	}

	return types.JavaBoolFalse
}

// "java/lang/ProcessBuilder.redirectErrorStream(Z)Ljava/lang/ProcessBuilder;"
func processBuilderRedirectErrorStreamSet(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderRedirectErrorStream: self is null")
	}

	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "processBuilderRedirectErrorStream: self is not an object")
	}

	if len(params) < 2 || params[1] == nil || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderRedirectErrorStream: redirectErrorStream parameter is null")
	}

	var redirect bool
	switch val := params[1].(type) {
	case bool:
		redirect = val
	case int64:
		redirect = (val != 0)
	case int:
		redirect = (val != 0)
	case int32:
		redirect = (val != 0)
	case int8:
		redirect = (val != 0)
	case int16:
		redirect = (val != 0)
	case *object.Object:
		if fld, exists := val.FieldTable["value"]; exists {
			switch v := fld.Fvalue.(type) {
			case bool:
				redirect = v
			case int64:
				redirect = (v != 0)
			case int:
				redirect = (v != 0)
			default:
				return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("processBuilderRedirectErrorStream: unsupported value in Boolean object: %T", fld.Fvalue))
			}
		} else {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "processBuilderRedirectErrorStream: Boolean object missing value field")
		}
	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("processBuilderRedirectErrorStream: unsupported boolean parameter type %T", params[1]))
	}

	self.ThMutex.Lock()
	if self.FieldTable == nil {
		self.FieldTable = make(map[string]object.Field)
	}
	self.FieldTable["redirectErrorStream"] = object.Field{
		Ftype:  types.Bool,
		Fvalue: redirect,
	}
	self.ThMutex.Unlock()

	return self
}

// "java/lang/ProcessBuilder.start()Ljava/lang/Process;"
func processBuilderStart(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderStart: self is null")
	}

	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "processBuilderStart: self is not an object")
	}

	self.ThMutex.RLock()
	cmdField, ok := self.FieldTable["command"]
	self.ThMutex.RUnlock()
	if !ok || cmdField.Fvalue == nil || object.IsNull(cmdField.Fvalue) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderStart: command is null")
	}

	var cmdArgs []string
	switch val := cmdField.Fvalue.(type) {
	case *object.Object:
		if valField, exists := val.FieldTable["value"]; exists && valField.Fvalue != nil && !object.IsNull(valField.Fvalue) {
			switch elems := valField.Fvalue.(type) {
			case []*object.Object:
				for _, elem := range elems {
					if elem == nil || object.IsNull(elem) {
						return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderStart: command element is null")
					}
					cmdArgs = append(cmdArgs, object.GoStringFromStringObject(elem))
				}
			case []any:
				for _, elem := range elems {
					if elem == nil || object.IsNull(elem) {
						return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderStart: command element is null")
					}
					switch e := elem.(type) {
					case *object.Object:
						cmdArgs = append(cmdArgs, object.GoStringFromStringObject(e))
					case string:
						cmdArgs = append(cmdArgs, e)
					default:
						return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("processBuilderStart: invalid command element type %T", elem))
					}
				}
			case []string:
				for _, s := range elems {
					cmdArgs = append(cmdArgs, s)
				}
			default:
				return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("processBuilderStart: unsupported command list value type %T", valField.Fvalue))
			}
		} else {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "processBuilderStart: invalid command list object")
		}
	case []*object.Object:
		for _, elem := range val {
			if elem == nil || object.IsNull(elem) {
				return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderStart: command element is null")
			}
			cmdArgs = append(cmdArgs, object.GoStringFromStringObject(elem))
		}
	case []any:
		for _, elem := range val {
			if elem == nil || object.IsNull(elem) {
				return ghelpers.GetGErrBlk(excNames.NullPointerException, "processBuilderStart: command element is null")
			}
			switch e := elem.(type) {
			case *object.Object:
				cmdArgs = append(cmdArgs, object.GoStringFromStringObject(e))
			case string:
				cmdArgs = append(cmdArgs, e)
			default:
				return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("processBuilderStart: invalid command element type %T", elem))
			}
		}
	case []string:
		for _, s := range val {
			cmdArgs = append(cmdArgs, s)
		}
	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("processBuilderStart: unsupported command type %T", cmdField.Fvalue))
	}

	if len(cmdArgs) == 0 {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, "processBuilderStart: command list is empty")
	}

	var workingDir string
	self.ThMutex.RLock()
	dirField, hasDir := self.FieldTable["directory"]
	self.ThMutex.RUnlock()
	if hasDir && dirField.Fvalue != nil && !object.IsNull(dirField.Fvalue) {
		switch d := dirField.Fvalue.(type) {
		case *object.Object:
			if pathFld, hasPath := d.FieldTable[ghelpers.FilePath]; hasPath && pathFld.Fvalue != nil {
				switch pVal := pathFld.Fvalue.(type) {
				case []types.JavaByte:
					workingDir = object.GoStringFromJavaByteArray(pVal)
				case []byte:
					workingDir = string(pVal)
				case string:
					workingDir = pVal
				}
			} else if strVal := object.GoStringFromStringObject(d); strVal != "" {
				workingDir = strVal
			}
		case string:
			workingDir = d
		}
	}

	var envList []string
	self.ThMutex.RLock()
	envField, hasEnv := self.FieldTable["environment"]
	self.ThMutex.RUnlock()
	if hasEnv && envField.Fvalue != nil && !object.IsNull(envField.Fvalue) {
		envList = []string{}
		switch envObj := envField.Fvalue.(type) {
		case *object.Object:
			envObj.ThMutex.RLock()
			mapFld, ok := envObj.FieldTable["map"]
			envObj.ThMutex.RUnlock()
			if ok && mapFld.Fvalue != nil {
				if defMap, isDefMap := mapFld.Fvalue.(types.DefHashMap); isDefMap {
					for k, v := range defMap {
						var keyStr, valStr string
						switch kVal := k.(type) {
						case string:
							keyStr = kVal
						case *object.Object:
							keyStr = object.GoStringFromStringObject(kVal)
						default:
							keyStr = fmt.Sprintf("%v", kVal)
						}

						switch vVal := v.(type) {
						case string:
							valStr = vVal
						case *object.Object:
							valStr = object.GoStringFromStringObject(vVal)
						default:
							valStr = fmt.Sprintf("%v", vVal)
						}
						envList = append(envList, keyStr+"="+valStr)
					}
				}
			}
		case map[string]string:
			for k, v := range envObj {
				envList = append(envList, k+"="+v)
			}
		}
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	if workingDir != "" {
		cmd.Dir = workingDir
	}
	if envList != nil {
		cmd.Env = envList
	}

	var redirectErrorStream bool
	self.ThMutex.RLock()
	redirField, hasRedir := self.FieldTable["redirectErrorStream"]
	self.ThMutex.RUnlock()
	if hasRedir && redirField.Fvalue != nil {
		switch v := redirField.Fvalue.(type) {
		case bool:
			redirectErrorStream = v
		case int64:
			redirectErrorStream = (v != 0)
		case int:
			redirectErrorStream = (v != 0)
		}
	}

	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("processBuilderStart: %s", err.Error()))
	}
	cmd.Stdout = stdoutW
	if redirectErrorStream {
		cmd.Stderr = stdoutW
	}

	if err := cmd.Start(); err != nil {
		_ = stdoutR.Close()
		_ = stdoutW.Close()
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("processBuilderStart: %s", err.Error()))
	}
	_ = stdoutW.Close()

	procObj := object.MakeOneFieldObject("java/lang/Process", "pid", types.Int, int64(cmd.Process.Pid))
	inClassName := "java/io/FileInputStream"
	inObj := object.MakeEmptyObjectWithClassName(&inClassName)
	if inObj.FieldTable == nil {
		inObj.FieldTable = make(map[string]object.Field)
	}
	inObj.FieldTable[ghelpers.FileHandle] = object.Field{
		Ftype:  ghelpers.FileHandle,
		Fvalue: stdoutR,
	}
	inObj.FieldTable[ghelpers.FilePath] = object.Field{
		Ftype:  types.JavaByteArray,
		Fvalue: object.JavaByteArrayFromGoString("ProcessStdout"),
	}
	procObj.FieldTable["inputStream"] = object.Field{
		Ftype:  "java/io/InputStream",
		Fvalue: inObj,
	}
	procObj.FieldTable["stdout"] = object.Field{
		Ftype:  ghelpers.FileHandle,
		Fvalue: stdoutR,
	}

	return procObj
}
