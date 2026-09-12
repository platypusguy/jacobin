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
)

// Implementation of java/util/logging/StreamHandler.
// Strategy: StreamHandler = jacobin Object with the same fields as Handler
// (level, filter, formatter, encoding, errorManager), since StreamHandler
// extends Handler, plus its own field:
//
//	outputStream: *os.File -- the underlying stream this Handler publishes to
//
// This class is registered here (rather than left to be loaded from real
// JDK bytecode) so that ancestor classes such as ConsoleHandler and
// FileHandler never erroneously execute StreamHandler's real <clinit> or
// other bytecode, which this jacobin implementation does not otherwise support.

var streamHandlerClassName = "java/util/logging/StreamHandler"
var fieldNameStreamHandlerOutputStream = "outputStream"

func Load_Util_Logging_StreamHandler() {

	ghelpers.MethodSignatures["java/util/logging/StreamHandler.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/logging/StreamHandler.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingStreamHandlerInit,
		}

	ghelpers.MethodSignatures["java/util/logging/StreamHandler.<init>(Ljava/io/OutputStream;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/logging/StreamHandler.<init>(Ljava/io/OutputStream;Ljava/util/logging/Formatter;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/logging/StreamHandler.isLoggable(Ljava/util/logging/LogRecord;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/logging/StreamHandler.setEncoding(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/logging/StreamHandler.setOutputStream(Ljava/io/OutputStream;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/logging/StreamHandler.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingStreamHandlerClose,
		}

	ghelpers.MethodSignatures["java/util/logging/StreamHandler.flush()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingStreamHandlerFlush,
		}

	ghelpers.MethodSignatures["java/util/logging/StreamHandler.publish(Ljava/util/logging/LogRecord;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingStreamHandlerPublish,
		}
}

// "java/util/logging/StreamHandler.<init>()V"
// Default StreamHandler state: level = Level.INFO, no filter, no formatter,
// no encoding, no error manager, and no output stream set yet.
func loggingStreamHandlerInit(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingStreamHandlerInit: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	if obj.KlassName == 0 || obj.KlassName == types.InvalidStringIndex {
		obj.KlassName = object.StringPoolIndexFromGoString(streamHandlerClassName)
	}

	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameHandlerLevel] = object.Field{Ftype: types.Ref, Fvalue: makeLevelObject("INFO", standardLevels["INFO"], "")}
	obj.FieldTable[fieldNameHandlerFilter] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameHandlerFormatter] = object.Field{Ftype: types.Ref, Fvalue: makeDefaultSimpleFormatter()}
	obj.FieldTable[fieldNameHandlerEncoding] = object.Field{Ftype: types.StringClassRef, Fvalue: ""}
	obj.FieldTable[fieldNameHandlerErrorManager] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameStreamHandlerOutputStream] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	return nil
}

// "java/util/logging/StreamHandler.close()V"
// Flushes and closes the underlying output stream, if one is set.
func loggingStreamHandlerClose(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingStreamHandlerClose: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	obj.ThMutex.RLock()
	osFile, hasFile := obj.FieldTable[fieldNameStreamHandlerOutputStream].Fvalue.(*os.File)
	obj.ThMutex.RUnlock()
	if !hasFile || osFile == nil {
		return nil
	}

	_ = osFile.Sync()
	err := osFile.Close()
	if err != nil {
		errMsg := fmt.Sprintf("loggingStreamHandlerClose: osFile.Close() failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	return nil
}

// "java/util/logging/StreamHandler.flush()V"
// Flushes the underlying output stream, if one is set.
func loggingStreamHandlerFlush(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingStreamHandlerFlush: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	obj.ThMutex.RLock()
	osFile, hasFile := obj.FieldTable[fieldNameStreamHandlerOutputStream].Fvalue.(*os.File)
	obj.ThMutex.RUnlock()
	if !hasFile || osFile == nil {
		return nil
	}

	_ = osFile.Sync()
	return nil
}

// "java/util/logging/StreamHandler.publish(Ljava/util/logging/LogRecord;)V"
// Writes the LogRecord's message to the underlying output stream if it is
// loggable by this Handler and a stream has been set.
func loggingStreamHandlerPublish(params []interface{}) interface{} {
	fs, args := loggingExtractFsAndArgs(params)
	obj, ok := args[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingStreamHandlerPublish: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	if len(args) < 2 {
		errMsg := "loggingStreamHandlerPublish: Requires a LogRecord parameter"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	result := loggingHandlerIsLoggable(params)
	if result != types.JavaBoolTrue {
		return nil
	}

	record, ok := args[1].(*object.Object)
	if !ok || record == nil || object.IsNull(record) {
		return nil
	}

	msg := formatLogRecordWithHandlerFormatter(obj, record, fs)
	if msg == "" {
		return nil
	}

	obj.ThMutex.RLock()
	osFile, hasFile := obj.FieldTable[fieldNameStreamHandlerOutputStream].Fvalue.(*os.File)
	obj.ThMutex.RUnlock()
	if !hasFile || osFile == nil {
		return nil
	}

	_, err := fmt.Fprint(osFile, msg)
	if err != nil {
		errMsg := fmt.Sprintf("loggingStreamHandlerPublish: write failed, reason: %s", err.Error())
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	return nil
}
