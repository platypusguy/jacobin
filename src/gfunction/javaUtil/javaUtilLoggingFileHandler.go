/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"strings"
)

// Implementation of java/util/logging/FileHandler.
// Strategy: FileHandler = jacobin Object with the same fields as Handler
// (level, filter, formatter, encoding, errorManager), since FileHandler
// extends Handler, plus its own fields:
//   pattern:  Go string  -- the file name pattern used to create the log file
//   limit:    int64      -- the maximum number of bytes to write to any one file
//   count:    int64      -- the number of output files to cycle through
//   append:   int64      -- whether to append to (vs. overwrite) an existing file
//   fileHandle: *os.File -- the actual file handle backing this Handler
// Published records are written to the underlying file, as in the real
// java.util.logging.FileHandler.

var fileHandlerClassName = "java/util/logging/FileHandler"
var fieldNameFileHandlerPattern = "pattern"
var fieldNameFileHandlerLimit = "limit"
var fieldNameFileHandlerCount = "count"
var fieldNameFileHandlerAppend = "append"
var fieldNameFileHandlerFile = "fileHandle"

const defaultFileHandlerPattern = "default_file_handler_%u.log"

func Load_Util_Logging_FileHandler() {

	ghelpers.MethodSignatures["java/util/logging/FileHandler.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/logging/FileHandler.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingFileHandlerInit,
		}

	ghelpers.MethodSignatures["java/util/logging/FileHandler.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingFileHandlerInitPattern,
		}

	ghelpers.MethodSignatures["java/util/logging/FileHandler.<init>(Ljava/lang/String;Z)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  loggingFileHandlerInitPatternAppend,
		}

	ghelpers.MethodSignatures["java/util/logging/FileHandler.<init>(Ljava/lang/String;II)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  loggingFileHandlerInitPatternLimitCount,
		}

	ghelpers.MethodSignatures["java/util/logging/FileHandler.<init>(Ljava/lang/String;IIZ)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  loggingFileHandlerInitPatternLimitCountAppend,
		}

	ghelpers.MethodSignatures["java/util/logging/FileHandler.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingFileHandlerClose,
		}

	ghelpers.MethodSignatures["java/util/logging/FileHandler.publish(Ljava/util/logging/LogRecord;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingFileHandlerPublish,
		}
}

// resolvePattern substitutes the "%u" and "%g" tokens used by the real
// java.util.logging.FileHandler with a fixed, simple value, since this
// implementation does not perform file rotation.
func resolveFileHandlerPattern(pattern string) string {
	if pattern == "" {
		pattern = defaultFileHandlerPattern
	}
	replacer := strings.NewReplacer("%u", "0", "%g", "0")
	return replacer.Replace(pattern)
}

// initFileHandlerCommon sets up the common state and opens the underlying
// file, given a jacobin Object, a file-name pattern, limit, count, and append flag.
func initFileHandlerCommon(obj *object.Object, pattern string, limit, count, appendFlag int64) interface{} {
	if obj.KlassName == 0 || obj.KlassName == types.InvalidStringIndex {
		obj.KlassName = object.StringPoolIndexFromGoString(fileHandlerClassName)
	}

	resolvedPath := resolveFileHandlerPattern(pattern)

	var osFile *os.File
	var err error
	if appendFlag != 0 {
		osFile, err = os.OpenFile(resolvedPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, ghelpers.CreateFilePermissions)
	} else {
		osFile, err = os.Create(resolvedPath)
	}
	if err != nil {
		errMsg := fmt.Sprintf("initFileHandlerCommon: could not open %s, reason: %s", resolvedPath, err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameHandlerLevel] = object.Field{Ftype: types.Ref, Fvalue: makeLevelObject("ALL", standardLevels["ALL"], "")}
	obj.FieldTable[fieldNameHandlerFilter] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameHandlerFormatter] = object.Field{Ftype: types.Ref, Fvalue: makeDefaultSimpleFormatter()}
	obj.FieldTable[fieldNameHandlerEncoding] = object.Field{Ftype: types.StringClassRef, Fvalue: ""}
	obj.FieldTable[fieldNameHandlerErrorManager] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameFileHandlerPattern] = object.Field{Ftype: types.StringClassRef, Fvalue: pattern}
	obj.FieldTable[fieldNameFileHandlerLimit] = object.Field{Ftype: types.Int, Fvalue: limit}
	obj.FieldTable[fieldNameFileHandlerCount] = object.Field{Ftype: types.Int, Fvalue: count}
	obj.FieldTable[fieldNameFileHandlerAppend] = object.Field{Ftype: types.Int, Fvalue: appendFlag}
	obj.FieldTable[fieldNameFileHandlerFile] = object.Field{Ftype: ghelpers.FileHandle, Fvalue: osFile}
	return nil
}

// "java/util/logging/FileHandler.<init>()V"
func loggingFileHandlerInit(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingFileHandlerInit: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	return initFileHandlerCommon(obj, defaultFileHandlerPattern, 0, 1, 0)
}

// "java/util/logging/FileHandler.<init>(Ljava/lang/String;)V"
func loggingFileHandlerInitPattern(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingFileHandlerInitPattern: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	patternObj, ok := params[1].(*object.Object)
	if !ok || patternObj == nil || object.IsNull(patternObj) {
		errMsg := "loggingFileHandlerInitPattern: The pattern parameter is not a valid string"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	pattern := object.GoStringFromStringObject(patternObj)
	return initFileHandlerCommon(obj, pattern, 0, 1, 0)
}

// "java/util/logging/FileHandler.<init>(Ljava/lang/String;Z)V"
func loggingFileHandlerInitPatternAppend(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingFileHandlerInitPatternAppend: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	patternObj, ok := params[1].(*object.Object)
	if !ok || patternObj == nil || object.IsNull(patternObj) {
		errMsg := "loggingFileHandlerInitPatternAppend: The pattern parameter is not a valid string"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	pattern := object.GoStringFromStringObject(patternObj)
	appendFlag, ok := params[2].(int64)
	if !ok {
		errMsg := "loggingFileHandlerInitPatternAppend: Missing append-boolean argument"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	return initFileHandlerCommon(obj, pattern, 0, 1, appendFlag)
}

// "java/util/logging/FileHandler.<init>(Ljava/lang/String;II)V"
func loggingFileHandlerInitPatternLimitCount(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingFileHandlerInitPatternLimitCount: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	patternObj, ok := params[1].(*object.Object)
	if !ok || patternObj == nil || object.IsNull(patternObj) {
		errMsg := "loggingFileHandlerInitPatternLimitCount: The pattern parameter is not a valid string"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	pattern := object.GoStringFromStringObject(patternObj)
	limit, ok := params[2].(int64)
	if !ok {
		errMsg := "loggingFileHandlerInitPatternLimitCount: Missing limit argument"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	count, ok := params[3].(int64)
	if !ok {
		errMsg := "loggingFileHandlerInitPatternLimitCount: Missing count argument"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	if limit < 0 || count < 1 {
		errMsg := "loggingFileHandlerInitPatternLimitCount: limit must be >= 0 and count must be >= 1"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	return initFileHandlerCommon(obj, pattern, limit, count, 0)
}

// "java/util/logging/FileHandler.<init>(Ljava/lang/String;IIZ)V"
func loggingFileHandlerInitPatternLimitCountAppend(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingFileHandlerInitPatternLimitCountAppend: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	patternObj, ok := params[1].(*object.Object)
	if !ok || patternObj == nil || object.IsNull(patternObj) {
		errMsg := "loggingFileHandlerInitPatternLimitCountAppend: The pattern parameter is not a valid string"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	pattern := object.GoStringFromStringObject(patternObj)
	limit, ok := params[2].(int64)
	if !ok {
		errMsg := "loggingFileHandlerInitPatternLimitCountAppend: Missing limit argument"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	count, ok := params[3].(int64)
	if !ok {
		errMsg := "loggingFileHandlerInitPatternLimitCountAppend: Missing count argument"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	if limit < 0 || count < 1 {
		errMsg := "loggingFileHandlerInitPatternLimitCountAppend: limit must be >= 0 and count must be >= 1"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	appendFlag, ok := params[4].(int64)
	if !ok {
		errMsg := "loggingFileHandlerInitPatternLimitCountAppend: Missing append-boolean argument"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	return initFileHandlerCommon(obj, pattern, limit, count, appendFlag)
}

// "java/util/logging/FileHandler.close()V"
// Closes the underlying file that this Handler publishes to.
func loggingFileHandlerClose(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingFileHandlerClose: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	obj.ThMutex.RLock()
	osFile, hasFile := obj.FieldTable[fieldNameFileHandlerFile].Fvalue.(*os.File)
	obj.ThMutex.RUnlock()
	if !hasFile || osFile == nil {
		errMsg := "loggingFileHandlerClose: FileHandler object lacks a file handle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	err := osFile.Close()
	if err != nil {
		errMsg := fmt.Sprintf("loggingFileHandlerClose: osFile.Close() failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	return nil
}

// "java/util/logging/FileHandler.publish(Ljava/util/logging/LogRecord;)V"
// Writes the LogRecord's message to the underlying file if it is loggable by this Handler.
func loggingFileHandlerPublish(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingFileHandlerPublish: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	if len(params) < 2 {
		errMsg := "loggingFileHandlerPublish: Requires a LogRecord parameter"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	result := loggingHandlerIsLoggable(params)
	if result != types.JavaBoolTrue {
		return nil
	}

	record, ok := params[1].(*object.Object)
	if !ok || record == nil || object.IsNull(record) {
		return nil
	}

	msg := formatLogRecordWithHandlerFormatter(obj, record)
	if msg == "" {
		return nil
	}

	obj.ThMutex.RLock()
	osFile, hasFile := obj.FieldTable[fieldNameFileHandlerFile].Fvalue.(*os.File)
	obj.ThMutex.RUnlock()
	if !hasFile || osFile == nil {
		errMsg := "loggingFileHandlerPublish: FileHandler object lacks a file handle field"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}

	_, err := fmt.Fprint(osFile, msg)
	if err != nil {
		errMsg := fmt.Sprintf("loggingFileHandlerPublish: write failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	return nil
}
