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

// Implementation of java/util/logging/SimpleFormatter.
// Strategy: SimpleFormatter extends Formatter and produces a brief,
// human-readable two-line summary of a LogRecord, matching the format
// used by the real JDK (and HotSpot's default console output):
//
//	<date> <calling class> <calling method>
//	<LEVEL>: <message>
//
// e.g.
//
//	Sep 11, 2026 2:48:52 PM main main
//	INFO: Message #1

var simpleFormatterClassName = "java/util/logging/SimpleFormatter"

// simpleFormatterDateLayout mirrors the default JDK format
// "MMM d, yyyy h:mm:ss a" (e.g. "Sep 11, 2026 2:48:52 PM").
const simpleFormatterDateLayout = "Jan 2, 2006 3:04:05 PM"

func Load_Util_Logging_SimpleFormatter() {

	ghelpers.MethodSignatures["java/util/logging/SimpleFormatter.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/logging/SimpleFormatter.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingSimpleFormatterInit,
		}

	ghelpers.MethodSignatures["java/util/logging/SimpleFormatter.format(Ljava/util/logging/LogRecord;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingSimpleFormatterFormat,
		}
}

// makeDefaultSimpleFormatter creates a new java/util/logging/SimpleFormatter
// instance suitable for use as a Handler's default formatter, mirroring the
// real JDK behavior where every Handler is initialized with a SimpleFormatter
// unless a custom formatter is explicitly supplied.
func makeDefaultSimpleFormatter() *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&simpleFormatterClassName)
	_ = loggingSimpleFormatterInit([]interface{}{obj})
	return obj
}

// "java/util/logging/SimpleFormatter.<init>()V"
func loggingSimpleFormatterInit(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingSimpleFormatterInit: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	if obj.KlassName == 0 || obj.KlassName == types.InvalidStringIndex {
		obj.KlassName = object.StringPoolIndexFromGoString(simpleFormatterClassName)
	}
	return nil
}

// "java/util/logging/SimpleFormatter.format(Ljava/util/logging/LogRecord;)Ljava/lang/String;"
// Produces:
//
//	<date> <calling class> <calling method>
//	<LEVEL>: <message>
func loggingSimpleFormatterFormat(params []interface{}) interface{} {
	record, ok := params[0].(*object.Object)
	if !ok || record == nil {
		errMsg := "loggingSimpleFormatterFormat: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	record.ThMutex.RLock()
	millis, _ := record.FieldTable[fieldNameLogRecordMillis].Fvalue.(int64)

	className := ""
	if cls, ok := record.FieldTable[fieldNameLogRecordSourceClassName].Fvalue.(*object.Object); ok && cls != nil {
		className = object.GoStringFromStringObject(cls)
	}
	methodName := ""
	if mth, ok := record.FieldTable[fieldNameLogRecordSourceMethodName].Fvalue.(*object.Object); ok && mth != nil {
		methodName = object.GoStringFromStringObject(mth)
	}
	if className == "" {
		if lg, ok := record.FieldTable[fieldNameLogRecordLoggerName].Fvalue.(*object.Object); ok && lg != nil {
			className = object.GoStringFromStringObject(lg)
		}
	}

	levelName := "INFO"
	if level, ok := record.FieldTable[fieldNameHandlerLevel].Fvalue.(*object.Object); ok && level != nil {
		if name, ok := level.FieldTable[fieldNameLevelName].Fvalue.(string); ok && name != "" {
			levelName = name
		}
	}
	record.ThMutex.RUnlock()

	message := ""
	if msg, ok := loggingFormatterFormatMessage([]interface{}{record}).(*object.Object); ok && msg != nil {
		message = object.GoStringFromStringObject(msg)
	}

	dateStr := time.UnixMilli(millis).Local().Format(simpleFormatterDateLayout)

	firstLine := dateStr
	if className != "" {
		firstLine += " " + className
	}
	if methodName != "" {
		firstLine += " " + methodName
	}

	formatted := firstLine + "\n" + levelName + ": " + message + "\n"
	return object.StringObjectFromGoString(formatted)
}
