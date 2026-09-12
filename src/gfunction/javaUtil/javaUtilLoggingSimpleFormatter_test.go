/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"strings"
	"testing"
	"time"
)

func newSimpleFormatterObj() *object.Object {
	return object.MakeEmptyObjectWithClassName(&simpleFormatterClassName)
}

func TestLoad_Util_Logging_SimpleFormatter_Registration(t *testing.T) {
	globals.InitStringPool()
	Load_Util_Logging_SimpleFormatter()

	implementedEntries := map[string]int{
		"java/util/logging/SimpleFormatter.<clinit>()V":                                          0,
		"java/util/logging/SimpleFormatter.<init>()V":                                             0,
		"java/util/logging/SimpleFormatter.format(Ljava/util/logging/LogRecord;)Ljava/lang/String;": 1,
	}
	for key, slots := range implementedEntries {
		gm, ok := ghelpers.MethodSignatures[key]
		if !ok {
			t.Fatalf("expected %s to be registered", key)
		}
		if gm.ParamSlots != slots {
			t.Fatalf("expected %s to have %d param slots, got %d", key, slots, gm.ParamSlots)
		}
	}
}

func TestLoggingSimpleFormatterInit(t *testing.T) {
	globals.InitStringPool()

	f := object.MakeEmptyObject()
	if ret := loggingSimpleFormatterInit([]interface{}{f}); ret != nil {
		t.Fatalf("loggingSimpleFormatterInit returned error: %v", ret)
	}
	if object.GoStringFromStringPoolIndex(f.KlassName) != simpleFormatterClassName {
		t.Fatalf("expected class name %s, got %s", simpleFormatterClassName, object.GoStringFromStringPoolIndex(f.KlassName))
	}
}

func TestLoggingSimpleFormatterFormat(t *testing.T) {
	globals.InitStringPool()

	f := newSimpleFormatterObj()
	_ = loggingSimpleFormatterInit([]interface{}{f})

	millis := time.Date(2026, time.September, 11, 14, 48, 52, 0, time.Local).UnixMilli()

	levelObj := makeLevelObject("INFO", standardLevels["INFO"], "")
	record := object.MakeEmptyObject()
	record.FieldTable[fieldNameHandlerLevel] = object.Field{Fvalue: levelObj}
	record.FieldTable[fieldNameLogRecordMessage] = object.Field{Fvalue: object.StringObjectFromGoString("Message #1")}
	record.FieldTable[fieldNameLogRecordLoggerName] = object.Field{Fvalue: object.StringObjectFromGoString("main")}
	record.FieldTable[fieldNameLogRecordSourceClassName] = object.Field{Fvalue: object.StringObjectFromGoString("main")}
	record.FieldTable[fieldNameLogRecordSourceMethodName] = object.Field{Fvalue: object.StringObjectFromGoString("main")}
	record.FieldTable[fieldNameLogRecordMillis] = object.Field{Ftype: types.Long, Fvalue: millis}

	ret := loggingSimpleFormatterFormat([]interface{}{record})
	strObj, ok := ret.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", ret)
	}
	got := object.GoStringFromStringObject(strObj)
	lines := strings.Split(got, "\n")
	if len(lines) < 2 {
		t.Fatalf("expected at least two lines, got %q", got)
	}
	if !strings.HasSuffix(lines[0], "main main") {
		t.Fatalf("expected first line to end with 'main main', got %q", lines[0])
	}
	if lines[1] != "INFO: Message #1" {
		t.Fatalf("expected second line 'INFO: Message #1', got %q", lines[1])
	}
}

func TestLoggingSimpleFormatterFormat_FallsBackToLoggerName(t *testing.T) {
	globals.InitStringPool()

	levelObj := makeLevelObject("WARNING", standardLevels["WARNING"], "")
	record := object.MakeEmptyObject()
	record.FieldTable[fieldNameHandlerLevel] = object.Field{Fvalue: levelObj}
	record.FieldTable[fieldNameLogRecordMessage] = object.Field{Fvalue: object.StringObjectFromGoString("hello")}
	record.FieldTable[fieldNameLogRecordLoggerName] = object.Field{Fvalue: object.StringObjectFromGoString("myLogger")}

	ret := loggingSimpleFormatterFormat([]interface{}{record})
	strObj, ok := ret.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", ret)
	}
	got := object.GoStringFromStringObject(strObj)
	if !strings.Contains(got, "myLogger") {
		t.Fatalf("expected fallback to loggerName, got %q", got)
	}
	if !strings.Contains(got, "WARNING: hello") {
		t.Fatalf("expected 'WARNING: hello' in output, got %q", got)
	}
}

func TestLoggingSimpleFormatter_ErrorPaths(t *testing.T) {
	globals.InitStringPool()

	if ret := loggingSimpleFormatterInit([]interface{}{"not an object"}); ret == nil {
		t.Fatalf("expected error for non-object param")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}

	if ret := loggingSimpleFormatterFormat([]interface{}{"not an object"}); ret == nil {
		t.Fatalf("expected error for non-object param")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}
}
