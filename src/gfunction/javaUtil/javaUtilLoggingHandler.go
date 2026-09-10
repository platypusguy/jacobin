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
	"jacobin/src/statics"
	"jacobin/src/types"
	"os"
)

// Implementation of java/util/logging/Handler.
// Strategy: Handler = jacobin Object with fields:
//   level:       *object.Object (java/util/logging/Level) -- the level threshold for this Handler
//   filter:      *object.Object (java/util/logging/Filter) -- the filter for this Handler
//   formatter:   *object.Object (java/util/logging/Formatter) -- the formatter for this Handler
//   encoding:    Go string -- the character encoding used by this Handler
//   errorManager: *object.Object (java/util/logging/ErrorManager) -- the ErrorManager for this Handler

var handlerClassName = "java/util/logging/Handler"
var fieldNameHandlerLevel = "level"
var fieldNameHandlerFilter = "filter"
var fieldNameHandlerFormatter = "formatter"
var fieldNameHandlerEncoding = "encoding"
var fieldNameHandlerErrorManager = "errorManager"

func Load_Util_Logging_Handler() {

	ghelpers.MethodSignatures["java/util/logging/Handler.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingHandlerInit,
		}

	ghelpers.MethodSignatures["java/util/logging/Handler.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingHandlerClose,
		}

	ghelpers.MethodSignatures["java/util/logging/Handler.flush()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingHandlerFlush,
		}

	ghelpers.MethodSignatures["java/util/logging/Handler.getEncoding()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingHandlerGetEncoding,
		}

	ghelpers.MethodSignatures["java/util/logging/Handler.getErrorManager()Ljava/util/logging/ErrorManager;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingHandlerGetErrorManager,
		}

	ghelpers.MethodSignatures["java/util/logging/Handler.getFilter()Ljava/util/logging/Filter;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingHandlerGetFilter,
		}

	ghelpers.MethodSignatures["java/util/logging/Handler.getFormatter()Ljava/util/logging/Formatter;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingHandlerGetFormatter,
		}

	ghelpers.MethodSignatures["java/util/logging/Handler.getLevel()Ljava/util/logging/Level;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingHandlerGetLevel,
		}

	ghelpers.MethodSignatures["java/util/logging/Handler.isLoggable(Ljava/util/logging/LogRecord;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingHandlerIsLoggable,
		}

	ghelpers.MethodSignatures["java/util/logging/Handler.publish(Ljava/util/logging/LogRecord;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingHandlerPublish,
		}

	ghelpers.MethodSignatures["java/util/logging/Handler.setEncoding(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingHandlerSetEncoding,
		}

	ghelpers.MethodSignatures["java/util/logging/Handler.setErrorManager(Ljava/util/logging/ErrorManager;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingHandlerSetErrorManager,
		}

	ghelpers.MethodSignatures["java/util/logging/Handler.setFilter(Ljava/util/logging/Filter;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingHandlerSetFilter,
		}

	ghelpers.MethodSignatures["java/util/logging/Handler.setFormatter(Ljava/util/logging/Formatter;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingHandlerSetFormatter,
		}

	ghelpers.MethodSignatures["java/util/logging/Handler.setLevel(Ljava/util/logging/Level;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingHandlerSetLevel,
		}
}

// "java/util/logging/Handler.<init>()V"
// Default Handler state: level = Level.ALL, no filter, no formatter, no encoding, no error manager.
func loggingHandlerInit(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingHandlerInit: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	if obj.KlassName == 0 || obj.KlassName == types.InvalidStringIndex {
		obj.KlassName = object.StringPoolIndexFromGoString(handlerClassName)
	}

	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameHandlerLevel] = object.Field{Ftype: types.Ref, Fvalue: makeLevelObject("ALL", standardLevels["ALL"], "")}
	obj.FieldTable[fieldNameHandlerFilter] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameHandlerFormatter] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameHandlerEncoding] = object.Field{Ftype: types.StringClassRef, Fvalue: ""}
	obj.FieldTable[fieldNameHandlerErrorManager] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	return nil
}

// "java/util/logging/Handler.close()V"
// Flushes System.err, since that is the stream this Handler publishes to.
func loggingHandlerClose([]interface{}) interface{} {
	stderr := statics.GetStaticValue("java/lang/System", "err").(*os.File)
	_ = stderr.Sync()
	return nil
}

// "java/util/logging/Handler.flush()V"
// Flushes System.err, since that is the stream this Handler publishes to.
func loggingHandlerFlush([]interface{}) interface{} {
	stderr := statics.GetStaticValue("java/lang/System", "err").(*os.File)
	_ = stderr.Sync()
	return nil
}

// "java/util/logging/Handler.getEncoding()Ljava/lang/String;"
func loggingHandlerGetEncoding(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingHandlerGetEncoding: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	encoding, _ := obj.FieldTable[fieldNameHandlerEncoding].Fvalue.(string)
	if encoding == "" {
		return object.Null
	}
	return object.StringObjectFromGoString(encoding)
}

// "java/util/logging/Handler.setEncoding(Ljava/lang/String;)V"
func loggingHandlerSetEncoding(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingHandlerSetEncoding: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	encoding := ""
	if encObj, ok := params[1].(*object.Object); ok && encObj != nil && !object.IsNull(encObj) {
		encoding = object.GoStringFromStringObject(encObj)
	}

	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameHandlerEncoding] = object.Field{Ftype: types.StringClassRef, Fvalue: encoding}
	return nil
}

// "java/util/logging/Handler.getErrorManager()Ljava/util/logging/ErrorManager;"
func loggingHandlerGetErrorManager(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingHandlerGetErrorManager: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	errMgr, ok := obj.FieldTable[fieldNameHandlerErrorManager].Fvalue.(*object.Object)
	if !ok || errMgr == nil {
		return object.Null
	}
	return errMgr
}

// "java/util/logging/Handler.setErrorManager(Ljava/util/logging/ErrorManager;)V"
func loggingHandlerSetErrorManager(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingHandlerSetErrorManager: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	errMgr := params[1]

	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameHandlerErrorManager] = object.Field{Ftype: types.Ref, Fvalue: errMgr}
	return nil
}

// "java/util/logging/Handler.getFilter()Ljava/util/logging/Filter;"
func loggingHandlerGetFilter(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingHandlerGetFilter: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	filter, ok := obj.FieldTable[fieldNameHandlerFilter].Fvalue.(*object.Object)
	if !ok || filter == nil {
		return object.Null
	}
	return filter
}

// "java/util/logging/Handler.setFilter(Ljava/util/logging/Filter;)V"
func loggingHandlerSetFilter(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingHandlerSetFilter: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	filter := params[1]

	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameHandlerFilter] = object.Field{Ftype: types.Ref, Fvalue: filter}
	return nil
}

// "java/util/logging/Handler.getFormatter()Ljava/util/logging/Formatter;"
func loggingHandlerGetFormatter(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingHandlerGetFormatter: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	formatter, ok := obj.FieldTable[fieldNameHandlerFormatter].Fvalue.(*object.Object)
	if !ok || formatter == nil {
		return object.Null
	}
	return formatter
}

// "java/util/logging/Handler.setFormatter(Ljava/util/logging/Formatter;)V"
func loggingHandlerSetFormatter(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingHandlerSetFormatter: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	formatter := params[1]

	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameHandlerFormatter] = object.Field{Ftype: types.Ref, Fvalue: formatter}
	return nil
}

// "java/util/logging/Handler.getLevel()Ljava/util/logging/Level;"
func loggingHandlerGetLevel(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingHandlerGetLevel: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	level, ok := obj.FieldTable[fieldNameHandlerLevel].Fvalue.(*object.Object)
	if !ok || level == nil {
		return makeLevelObject("ALL", standardLevels["ALL"], "")
	}
	return level
}

// "java/util/logging/Handler.setLevel(Ljava/util/logging/Level;)V"
func loggingHandlerSetLevel(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingHandlerSetLevel: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	level, ok := params[1].(*object.Object)
	if !ok || level == nil {
		errMsg := "loggingHandlerSetLevel: The level parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameHandlerLevel] = object.Field{Ftype: types.Ref, Fvalue: level}
	return nil
}

// "java/util/logging/Handler.isLoggable(Ljava/util/logging/LogRecord;)Z"
// A LogRecord is loggable if this Handler's filter (if any) accepts it, and
// the LogRecord's level is >= this Handler's level.
func loggingHandlerIsLoggable(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingHandlerIsLoggable: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	obj.ThMutex.RLock()
	filter, hasFilter := obj.FieldTable[fieldNameHandlerFilter].Fvalue.(*object.Object)
	level, hasLevel := obj.FieldTable[fieldNameHandlerLevel].Fvalue.(*object.Object)
	obj.ThMutex.RUnlock()

	if hasFilter && filter != nil && !object.IsNull(filter) {
		result := loggingFilterIsLoggable(params[1:])
		if result != types.JavaBoolTrue {
			return types.JavaBoolFalse
		}
	}

	if hasLevel && level != nil {
		record, ok := params[1].(*object.Object)
		if ok && record != nil && !object.IsNull(record) {
			if recLevel, ok := record.FieldTable[fieldNameLevelValue].Fvalue.(int64); ok {
				handlerLevelValue, _ := level.FieldTable[fieldNameLevelValue].Fvalue.(int64)
				if recLevel < handlerLevelValue {
					return types.JavaBoolFalse
				}
			}
		}
	}

	return types.JavaBoolTrue
}

// "java/util/logging/Handler.publish(Ljava/util/logging/LogRecord;)V"
// Writes the LogRecord's message to System.err if it is loggable by this Handler.
func loggingHandlerPublish(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingHandlerPublish: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	if len(params) < 2 {
		errMsg := "loggingHandlerPublish: Requires a LogRecord parameter"
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

	msgFld, exists := record.FieldTable["message"]
	if !exists {
		return nil
	}
	msgObj, ok := msgFld.Fvalue.(*object.Object)
	if !ok || msgObj == nil || object.IsNull(msgObj) {
		return nil
	}
	msg := object.GoStringFromStringObject(msgObj)

	stderr := statics.GetStaticValue("java/lang/System", "err").(*os.File)
	_, _ = fmt.Fprintln(stderr, msg)
	return nil
}
