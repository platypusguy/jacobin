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

// Implementation of java/util/logging/ConsoleHandler.
// Strategy: ConsoleHandler = jacobin Object with the same fields as Handler
// (level, filter, formatter, encoding, errorManager), since ConsoleHandler
// extends StreamHandler which extends Handler. Published records are
// written to System.err, as in the real java.util.logging.ConsoleHandler.

var consoleHandlerClassName = "java/util/logging/ConsoleHandler"

func Load_Util_Logging_ConsoleHandler() {

	ghelpers.MethodSignatures["java/util/logging/ConsoleHandler.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/logging/ConsoleHandler.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingConsoleHandlerInit,
		}

	ghelpers.MethodSignatures["java/util/logging/ConsoleHandler.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingConsoleHandlerClose,
		}

	ghelpers.MethodSignatures["java/util/logging/ConsoleHandler.publish(Ljava/util/logging/LogRecord;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingConsoleHandlerPublish,
		}
}

// "java/util/logging/ConsoleHandler.<init>()V"
// Default ConsoleHandler state: level = Level.INFO, no filter, no formatter, no encoding, no error manager.
func loggingConsoleHandlerInit(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingConsoleHandlerInit: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	if obj.KlassName == 0 || obj.KlassName == types.InvalidStringIndex {
		obj.KlassName = object.StringPoolIndexFromGoString(consoleHandlerClassName)
	}

	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameHandlerLevel] = object.Field{Ftype: types.Ref, Fvalue: makeLevelObject("INFO", standardLevels["INFO"], "")}
	obj.FieldTable[fieldNameHandlerFilter] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameHandlerFormatter] = object.Field{Ftype: types.Ref, Fvalue: makeDefaultSimpleFormatter()}
	obj.FieldTable[fieldNameHandlerEncoding] = object.Field{Ftype: types.StringClassRef, Fvalue: ""}
	obj.FieldTable[fieldNameHandlerErrorManager] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	return nil
}

// "java/util/logging/ConsoleHandler.close()V"
// Flushes System.err, since that is the stream this Handler publishes to.
func loggingConsoleHandlerClose([]interface{}) interface{} {
	stderr := statics.GetStaticValue("java/lang/System", "err").(*os.File)
	_ = stderr.Sync()
	return nil
}

// "java/util/logging/ConsoleHandler.publish(Ljava/util/logging/LogRecord;)V"
// Writes the LogRecord's message to System.err if it is loggable by this Handler.
func loggingConsoleHandlerPublish(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingConsoleHandlerPublish: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	if len(params) < 2 {
		errMsg := "loggingConsoleHandlerPublish: Requires a LogRecord parameter"
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

	stderr := statics.GetStaticValue("java/lang/System", "err").(*os.File)
	_, _ = fmt.Fprint(stderr, msg)
	return nil
}
