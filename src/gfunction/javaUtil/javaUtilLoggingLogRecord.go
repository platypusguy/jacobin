/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"time"
)

// Implementation of java/util/logging/LogRecord.
// Strategy: LogRecord = jacobin Object with fields:
//   level:              *object.Object (java/util/logging/Level)
//   message:            *object.Object (java/lang/String)
//   loggerName:         *object.Object (java/lang/String)
//   sequenceNumber:     int64
//   sourceClassName:    *object.Object (java/lang/String)
//   sourceMethodName:   *object.Object (java/lang/String)
//   parameters:         *object.Object (Object[])
//   threadID:           int64
//   millis:             int64
//   thrown:             *object.Object (java/lang/Throwable)
//   resourceBundleName: *object.Object (java/lang/String)
//   resourceBundle:     *object.Object (java/util/ResourceBundle)
//
// getMillis()/setMillis(long) are deprecated as of Java 9 in favor of
// getInstant()/setInstant(Instant), so they use ghelpers.TrapDeprecated.

var logRecordClassName = "java/util/logging/LogRecord"

var fieldNameLogRecordMessage = "message"
var fieldNameLogRecordLoggerName = "loggerName"
var fieldNameLogRecordSequenceNumber = "sequenceNumber"
var fieldNameLogRecordSourceClassName = "sourceClassName"
var fieldNameLogRecordSourceMethodName = "sourceMethodName"
var fieldNameLogRecordParameters = "parameters"
var fieldNameLogRecordThreadID = "threadID"
var fieldNameLogRecordMillis = "millis"
var fieldNameLogRecordThrown = "thrown"
var fieldNameLogRecordResourceBundleName = "resourceBundleName"
var fieldNameLogRecordResourceBundle = "resourceBundle"

func Load_Util_Logging_LogRecord() {

	ghelpers.MethodSignatures["java/util/logging/LogRecord.<init>(Ljava/util/logging/Level;Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  loggingLogRecordInit,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.getLevel()Ljava/util/logging/Level;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingLogRecordGetLevel,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.getLoggerName()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingLogRecordGetLoggerName,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.getMessage()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingLogRecordGetMessage,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.getMillis()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.getParameters()[Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingLogRecordGetParameters,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.getResourceBundle()Ljava/util/ResourceBundle;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingLogRecordGetResourceBundle,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.getResourceBundleName()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingLogRecordGetResourceBundleName,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.getSequenceNumber()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingLogRecordGetSequenceNumber,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.getSourceClassName()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingLogRecordGetSourceClassName,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.getSourceMethodName()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingLogRecordGetSourceMethodName,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.getThreadID()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingLogRecordGetThreadID,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.getThrown()Ljava/lang/Throwable;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingLogRecordGetThrown,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.setLevel(Ljava/util/logging/Level;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingLogRecordSetLevel,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.setLoggerName(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingLogRecordSetLoggerName,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.setMessage(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingLogRecordSetMessage,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.setMillis(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.setParameters([Ljava/lang/Object;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingLogRecordSetParameters,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.setResourceBundle(Ljava/util/ResourceBundle;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingLogRecordSetResourceBundle,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.setResourceBundleName(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingLogRecordSetResourceBundleName,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.setSequenceNumber(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingLogRecordSetSequenceNumber,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.setSourceClassName(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingLogRecordSetSourceClassName,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.setSourceMethodName(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingLogRecordSetSourceMethodName,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.setThreadID(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingLogRecordSetThreadID,
		}

	ghelpers.MethodSignatures["java/util/logging/LogRecord.setThrown(Ljava/lang/Throwable;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingLogRecordSetThrown,
		}
}

// "java/util/logging/LogRecord.<init>(Ljava/util/logging/Level;Ljava/lang/String;)V"
func loggingLogRecordInit(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordInit: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	if obj.KlassName == 0 || obj.KlassName == types.InvalidStringIndex {
		obj.KlassName = object.StringPoolIndexFromGoString(logRecordClassName)
	}

	level := params[1]
	message := params[2]

	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameHandlerLevel] = object.Field{Ftype: types.Ref, Fvalue: level}
	obj.FieldTable[fieldNameLogRecordMessage] = object.Field{Ftype: types.Ref, Fvalue: message}
	obj.FieldTable[fieldNameLogRecordLoggerName] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameLogRecordSequenceNumber] = object.Field{Ftype: types.Long, Fvalue: int64(0)}
	obj.FieldTable[fieldNameLogRecordSourceClassName] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameLogRecordSourceMethodName] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameLogRecordParameters] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameLogRecordThreadID] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
	obj.FieldTable[fieldNameLogRecordMillis] = object.Field{Ftype: types.Long, Fvalue: time.Now().UnixMilli()}
	obj.FieldTable[fieldNameLogRecordThrown] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameLogRecordResourceBundleName] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameLogRecordResourceBundle] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	return nil
}

// "java/util/logging/LogRecord.getLevel()Ljava/util/logging/Level;"
func loggingLogRecordGetLevel(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordGetLevel: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	level, ok := obj.FieldTable[fieldNameHandlerLevel].Fvalue.(*object.Object)
	if !ok || level == nil {
		return object.Null
	}
	return level
}

// "java/util/logging/LogRecord.setLevel(Ljava/util/logging/Level;)V"
func loggingLogRecordSetLevel(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordSetLevel: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameHandlerLevel] = object.Field{Ftype: types.Ref, Fvalue: params[1]}
	return nil
}

// "java/util/logging/LogRecord.getMessage()Ljava/lang/String;"
func loggingLogRecordGetMessage(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordGetMessage: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	msg, ok := obj.FieldTable[fieldNameLogRecordMessage].Fvalue.(*object.Object)
	if !ok || msg == nil {
		return object.Null
	}
	return msg
}

// "java/util/logging/LogRecord.setMessage(Ljava/lang/String;)V"
func loggingLogRecordSetMessage(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordSetMessage: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameLogRecordMessage] = object.Field{Ftype: types.Ref, Fvalue: params[1]}
	return nil
}

// "java/util/logging/LogRecord.getLoggerName()Ljava/lang/String;"
func loggingLogRecordGetLoggerName(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordGetLoggerName: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	name, ok := obj.FieldTable[fieldNameLogRecordLoggerName].Fvalue.(*object.Object)
	if !ok || name == nil {
		return object.Null
	}
	return name
}

// "java/util/logging/LogRecord.setLoggerName(Ljava/lang/String;)V"
func loggingLogRecordSetLoggerName(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordSetLoggerName: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameLogRecordLoggerName] = object.Field{Ftype: types.Ref, Fvalue: params[1]}
	return nil
}

// "java/util/logging/LogRecord.getSequenceNumber()J"
func loggingLogRecordGetSequenceNumber(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordGetSequenceNumber: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	seq, _ := obj.FieldTable[fieldNameLogRecordSequenceNumber].Fvalue.(int64)
	return seq
}

// "java/util/logging/LogRecord.setSequenceNumber(J)V"
func loggingLogRecordSetSequenceNumber(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordSetSequenceNumber: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	seq, _ := params[1].(int64)
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameLogRecordSequenceNumber] = object.Field{Ftype: types.Long, Fvalue: seq}
	return nil
}

// "java/util/logging/LogRecord.getSourceClassName()Ljava/lang/String;"
func loggingLogRecordGetSourceClassName(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordGetSourceClassName: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	name, ok := obj.FieldTable[fieldNameLogRecordSourceClassName].Fvalue.(*object.Object)
	if !ok || name == nil {
		return object.Null
	}
	return name
}

// "java/util/logging/LogRecord.setSourceClassName(Ljava/lang/String;)V"
func loggingLogRecordSetSourceClassName(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordSetSourceClassName: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameLogRecordSourceClassName] = object.Field{Ftype: types.Ref, Fvalue: params[1]}
	return nil
}

// "java/util/logging/LogRecord.getSourceMethodName()Ljava/lang/String;"
func loggingLogRecordGetSourceMethodName(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordGetSourceMethodName: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	name, ok := obj.FieldTable[fieldNameLogRecordSourceMethodName].Fvalue.(*object.Object)
	if !ok || name == nil {
		return object.Null
	}
	return name
}

// "java/util/logging/LogRecord.setSourceMethodName(Ljava/lang/String;)V"
func loggingLogRecordSetSourceMethodName(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordSetSourceMethodName: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameLogRecordSourceMethodName] = object.Field{Ftype: types.Ref, Fvalue: params[1]}
	return nil
}

// "java/util/logging/LogRecord.getParameters()[Ljava/lang/Object;"
func loggingLogRecordGetParameters(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordGetParameters: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	parameters, ok := obj.FieldTable[fieldNameLogRecordParameters].Fvalue.(*object.Object)
	if !ok || parameters == nil {
		return object.Null
	}
	return parameters
}

// "java/util/logging/LogRecord.setParameters([Ljava/lang/Object;)V"
func loggingLogRecordSetParameters(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordSetParameters: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameLogRecordParameters] = object.Field{Ftype: types.Ref, Fvalue: params[1]}
	return nil
}

// "java/util/logging/LogRecord.getThreadID()I"
func loggingLogRecordGetThreadID(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordGetThreadID: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	tid, _ := obj.FieldTable[fieldNameLogRecordThreadID].Fvalue.(int64)
	return tid
}

// "java/util/logging/LogRecord.setThreadID(I)V"
func loggingLogRecordSetThreadID(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordSetThreadID: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	tid, _ := params[1].(int64)
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameLogRecordThreadID] = object.Field{Ftype: types.Int, Fvalue: tid}
	return nil
}

// "java/util/logging/LogRecord.getThrown()Ljava/lang/Throwable;"
func loggingLogRecordGetThrown(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordGetThrown: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	thrown, ok := obj.FieldTable[fieldNameLogRecordThrown].Fvalue.(*object.Object)
	if !ok || thrown == nil {
		return object.Null
	}
	return thrown
}

// "java/util/logging/LogRecord.setThrown(Ljava/lang/Throwable;)V"
func loggingLogRecordSetThrown(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordSetThrown: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameLogRecordThrown] = object.Field{Ftype: types.Ref, Fvalue: params[1]}
	return nil
}

// "java/util/logging/LogRecord.getResourceBundleName()Ljava/lang/String;"
func loggingLogRecordGetResourceBundleName(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordGetResourceBundleName: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	name, ok := obj.FieldTable[fieldNameLogRecordResourceBundleName].Fvalue.(*object.Object)
	if !ok || name == nil {
		return object.Null
	}
	return name
}

// "java/util/logging/LogRecord.setResourceBundleName(Ljava/lang/String;)V"
func loggingLogRecordSetResourceBundleName(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordSetResourceBundleName: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameLogRecordResourceBundleName] = object.Field{Ftype: types.Ref, Fvalue: params[1]}
	return nil
}

// "java/util/logging/LogRecord.getResourceBundle()Ljava/util/ResourceBundle;"
func loggingLogRecordGetResourceBundle(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordGetResourceBundle: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	bundle, ok := obj.FieldTable[fieldNameLogRecordResourceBundle].Fvalue.(*object.Object)
	if !ok || bundle == nil {
		return object.Null
	}
	return bundle
}

// "java/util/logging/LogRecord.setResourceBundle(Ljava/util/ResourceBundle;)V"
func loggingLogRecordSetResourceBundle(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLogRecordSetResourceBundle: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameLogRecordResourceBundle] = object.Field{Ftype: types.Ref, Fvalue: params[1]}
	return nil
}
