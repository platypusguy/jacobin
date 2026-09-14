/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"container/list"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
)

// Implementation of java/util/logging/Formatter.
// Strategy: Formatter is an abstract class with no instance state of its
// own. It is registered here as a jacobin Object so subclasses (such as
// SimpleFormatter/XMLFormatter, if ever implemented) can extend it, and so
// its own methods (format/formatMessage/getHead/getTail) never erroneously
// fall through to real JDK bytecode.

var formatterClassName = "java/util/logging/Formatter"

func Load_Util_Logging_Formatter() {

	ghelpers.MethodSignatures["java/util/logging/Formatter.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/logging/Formatter.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingFormatterInit,
		}

	ghelpers.MethodSignatures["java/util/logging/Formatter.format(Ljava/util/logging/LogRecord;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingFormatterFormat,
		}

	ghelpers.MethodSignatures["java/util/logging/Formatter.formatMessage(Ljava/util/logging/LogRecord;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingFormatterFormatMessage,
		}

	ghelpers.MethodSignatures["java/util/logging/Formatter.getHead(Ljava/util/logging/Handler;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingFormatterGetHead,
		}

	ghelpers.MethodSignatures["java/util/logging/Formatter.getTail(Ljava/util/logging/Handler;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingFormatterGetTail,
		}
}

// "java/util/logging/Formatter.<init>()V"
func loggingFormatterInit(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingFormatterInit: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	if obj.KlassName == 0 || obj.KlassName == types.InvalidStringIndex {
		obj.KlassName = object.StringPoolIndexFromGoString(formatterClassName)
	}
	return nil
}

// "java/util/logging/Formatter.formatMessage(Ljava/util/logging/LogRecord;)Ljava/lang/String;"
// Real Java resolves the record's message via its resource bundle (if any)
// as a MessageFormat pattern, substituting in the record's parameters. This
// simplified implementation returns the record's raw message, since resource
// bundle/message-format substitution isn't otherwise modeled here.
func loggingFormatterFormatMessage(params []interface{}) interface{} {
	record, ok := params[0].(*object.Object)
	if !ok || record == nil {
		errMsg := "loggingFormatterFormatMessage: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	record.ThMutex.RLock()
	defer record.ThMutex.RUnlock()
	msg, ok := record.FieldTable[fieldNameLogRecordMessage].Fvalue.(*object.Object)
	if !ok || msg == nil {
		return object.StringObjectFromGoString("")
	}
	return msg
}

// "java/util/logging/Formatter.format(Ljava/util/logging/LogRecord;)Ljava/lang/String;"
// A basic textual representation: "<loggerName> <level>: <message>\n".
func loggingFormatterFormat(params []interface{}) interface{} {
	record, ok := params[0].(*object.Object)
	if !ok || record == nil {
		errMsg := "loggingFormatterFormat: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	record.ThMutex.RLock()
	loggerName := ""
	if lg, ok := record.FieldTable[fieldNameLogRecordLoggerName].Fvalue.(*object.Object); ok && lg != nil {
		loggerName = object.GoStringFromStringObject(lg)
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

	var formatted string
	if loggerName != "" {
		formatted = loggerName + " " + levelName + ": " + message + "\n"
	} else {
		formatted = levelName + ": " + message + "\n"
	}
	return object.StringObjectFromGoString(formatted)
}

// "java/util/logging/Formatter.getHead(Ljava/util/logging/Handler;)Ljava/lang/String;"
// Real Java's default returns an empty string; subclasses (e.g. XMLFormatter)
// override it to emit a header.
func loggingFormatterGetHead(params []interface{}) interface{} {
	if _, ok := params[0].(*object.Object); !ok {
		errMsg := "loggingFormatterGetHead: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	return object.StringObjectFromGoString("")
}

// "java/util/logging/Formatter.getTail(Ljava/util/logging/Handler;)Ljava/lang/String;"
// Real Java's default returns an empty string; subclasses (e.g. XMLFormatter)
// override it to emit a footer.
func loggingFormatterGetTail(params []interface{}) interface{} {
	if _, ok := params[0].(*object.Object); !ok {
		errMsg := "loggingFormatterGetTail: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	return object.StringObjectFromGoString("")
}

// formatLogRecordWithHandlerFormatter dispatches the record to the correct
// native format() implementation based on the Handler's formatter field,
// mirroring what a real Handler.publish() implementation does: it invokes
// getFormatter().format(record) before writing to the underlying sink. If no
// formatter is set (should not normally happen, since every Handler is
// initialized with a default SimpleFormatter), it falls back to the raw
// message text so nothing is silently dropped.
func formatLogRecordWithHandlerFormatter(handlerObj *object.Object, record *object.Object, fs *list.List) string {
	handlerObj.ThMutex.RLock()
	formatter, hasFormatter := handlerObj.FieldTable[fieldNameHandlerFormatter].Fvalue.(*object.Object)
	handlerObj.ThMutex.RUnlock()

	if hasFormatter && formatter != nil && !object.IsNull(formatter) {
		var ret interface{}
		className := object.GoStringFromStringPoolIndex(formatter.KlassName)
		switch className {
		case simpleFormatterClassName:
			ret = loggingSimpleFormatterFormat([]interface{}{record})
		case formatterClassName:
			ret = loggingFormatterFormat([]interface{}{record})
		default:
			// Possibly a user-supplied Formatter subclass with its own
			// format(LogRecord) override; dispatch to the actual bytecode
			// implementation, mirroring real Handler.publish() behavior.
			if custom, ok := loggingInvokeUserJavaMethod(fs, formatter, "format",
				"(Ljava/util/logging/LogRecord;)Ljava/lang/String;", record); ok {
				ret = custom
			} else {
				ret = loggingFormatterFormat([]interface{}{record})
			}
		}
		if strObj, ok := ret.(*object.Object); ok && strObj != nil {
			return object.GoStringFromStringObject(strObj)
		}
	}

	msgFld, exists := record.FieldTable[fieldNameLogRecordMessage]
	if !exists {
		return ""
	}
	msgObj, ok := msgFld.Fvalue.(*object.Object)
	if !ok || msgObj == nil || object.IsNull(msgObj) {
		return ""
	}
	return object.GoStringFromStringObject(msgObj) + "\n"
}
