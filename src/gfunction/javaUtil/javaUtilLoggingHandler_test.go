package javaUtil

import (
	"io"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"os"
	"strings"
	"testing"
)

func newHandlerObj() *object.Object {
	return object.MakeEmptyObjectWithClassName(&handlerClassName)
}

func makeLogRecord(levelName string, levelValue int64, message string) *object.Object {
	rec := object.MakeEmptyObject()
	rec.FieldTable[fieldNameHandlerLevel] = object.Field{Ftype: types.Ref, Fvalue: makeLevelObject(levelName, levelValue, "")}
	rec.FieldTable["message"] = object.Field{Ftype: types.Ref, Fvalue: object.StringObjectFromGoString(message)}
	return rec
}

func TestLoggingHandlerInit_DefaultsToLevelAll(t *testing.T) {
	globals.InitStringPool()

	h := newHandlerObj()
	ret := loggingHandlerInit([]interface{}{h})
	if ret != nil {
		t.Fatalf("loggingHandlerInit returned error: %v", ret)
	}

	level := loggingHandlerGetLevel([]interface{}{h}).(*object.Object)
	if got := loggingLevelIntValue([]interface{}{level}).(int64); got != standardLevels["ALL"] {
		t.Fatalf("expected default level ALL, got %d", got)
	}
	if ret := loggingHandlerGetFilter([]interface{}{h}); ret != object.Null {
		t.Fatalf("expected null filter by default")
	}
	if ret := loggingHandlerGetEncoding([]interface{}{h}); ret != object.Null {
		t.Fatalf("expected null encoding by default")
	}
}

func TestLoggingHandlerSettersAndGetters(t *testing.T) {
	globals.InitStringPool()

	h := newHandlerObj()
	_ = loggingHandlerInit([]interface{}{h})

	// setLevel/getLevel
	newLevel := makeLevelObject("WARNING", standardLevels["WARNING"], "")
	_ = loggingHandlerSetLevel([]interface{}{h, newLevel})
	if got := loggingHandlerGetLevel([]interface{}{h}).(*object.Object); got != newLevel {
		t.Fatalf("expected getLevel to return the set level object")
	}

	// setEncoding/getEncoding
	encStr := object.StringObjectFromGoString("UTF-8")
	_ = loggingHandlerSetEncoding([]interface{}{h, encStr})
	got := object.GoStringFromStringObject(loggingHandlerGetEncoding([]interface{}{h}).(*object.Object))
	if got != "UTF-8" {
		t.Fatalf("expected encoding UTF-8, got %q", got)
	}

	// setFilter/getFilter
	filterObj := object.MakeEmptyObject()
	_ = loggingHandlerSetFilter([]interface{}{h, filterObj})
	if got := loggingHandlerGetFilter([]interface{}{h}); got != filterObj {
		t.Fatalf("expected getFilter to return the set filter object")
	}

	// setFormatter/getFormatter
	formatterObj := object.MakeEmptyObject()
	_ = loggingHandlerSetFormatter([]interface{}{h, formatterObj})
	if got := loggingHandlerGetFormatter([]interface{}{h}); got != formatterObj {
		t.Fatalf("expected getFormatter to return the set formatter object")
	}

	// setErrorManager/getErrorManager
	errMgrObj := object.MakeEmptyObject()
	_ = loggingHandlerSetErrorManager([]interface{}{h, errMgrObj})
	if got := loggingHandlerGetErrorManager([]interface{}{h}); got != errMgrObj {
		t.Fatalf("expected getErrorManager to return the set error manager object")
	}
}

func TestLoggingHandlerIsLoggable_LevelThreshold(t *testing.T) {
	globals.InitStringPool()

	h := newHandlerObj()
	_ = loggingHandlerInit([]interface{}{h})
	_ = loggingHandlerSetLevel([]interface{}{h, makeLevelObject("WARNING", standardLevels["WARNING"], "")})

	// A record below the handler's level should not be loggable
	lowRecord := makeLogRecord("INFO", standardLevels["INFO"], "info msg")
	if ret := loggingHandlerIsLoggable([]interface{}{h, lowRecord}); ret != types.JavaBoolFalse {
		t.Fatalf("expected isLoggable false for INFO record with WARNING handler level")
	}

	// A record at or above the handler's level should be loggable
	highRecord := makeLogRecord("SEVERE", standardLevels["SEVERE"], "severe msg")
	if ret := loggingHandlerIsLoggable([]interface{}{h, highRecord}); ret != types.JavaBoolTrue {
		t.Fatalf("expected isLoggable true for SEVERE record with WARNING handler level")
	}
}

func TestLoggingHandlerPublish_WritesToSystemErr(t *testing.T) {
	globals.InitStringPool()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	defer func() {
		_ = statics.AddStatic("java/lang/System.err", statics.Static{Type: "GS", Value: os.Stderr})
	}()
	_ = statics.AddStatic("java/lang/System.err", statics.Static{Type: "GS", Value: w})

	h := newHandlerObj()
	_ = loggingHandlerInit([]interface{}{h})

	record := makeLogRecord("INFO", standardLevels["INFO"], "hello handler")
	ret := loggingHandlerPublish([]interface{}{h, record})
	if ret != nil {
		t.Fatalf("loggingHandlerPublish returned error: %v", ret)
	}

	_ = w.Close()
	buf, _ := io.ReadAll(r)
	got := string(buf)
	if !strings.Contains(got, "INFO: hello handler") {
		t.Fatalf("unexpected published output: %q", got)
	}
}

func TestLoggingHandlerPublish_NotLoggable_NoOutput(t *testing.T) {
	globals.InitStringPool()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	defer func() {
		_ = statics.AddStatic("java/lang/System.err", statics.Static{Type: "GS", Value: os.Stderr})
	}()
	_ = statics.AddStatic("java/lang/System.err", statics.Static{Type: "GS", Value: w})

	h := newHandlerObj()
	_ = loggingHandlerInit([]interface{}{h})
	_ = loggingHandlerSetLevel([]interface{}{h, makeLevelObject("SEVERE", standardLevels["SEVERE"], "")})

	record := makeLogRecord("INFO", standardLevels["INFO"], "should not appear")
	ret := loggingHandlerPublish([]interface{}{h, record})
	if ret != nil {
		t.Fatalf("loggingHandlerPublish returned error: %v", ret)
	}

	_ = w.Close()
	buf, _ := io.ReadAll(r)
	if len(buf) != 0 {
		t.Fatalf("expected no output for non-loggable record, got %q", string(buf))
	}
}

func TestLoggingHandlerFlushAndClose_SyncSystemErr(t *testing.T) {
	globals.InitStringPool()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	defer func() {
		_ = statics.AddStatic("java/lang/System.err", statics.Static{Type: "GS", Value: os.Stderr})
	}()
	_ = statics.AddStatic("java/lang/System.err", statics.Static{Type: "GS", Value: w})

	if ret := loggingHandlerFlush(nil); ret != nil {
		t.Fatalf("loggingHandlerFlush returned error: %v", ret)
	}
	if ret := loggingHandlerClose(nil); ret != nil {
		t.Fatalf("loggingHandlerClose returned error: %v", ret)
	}
}

func TestLoggingHandler_ErrorPaths(t *testing.T) {
	globals.InitStringPool()

	if ret := loggingHandlerInit([]interface{}{"not an object"}); ret == nil {
		t.Fatalf("expected error for non-object param")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}

	h := newHandlerObj()
	_ = loggingHandlerInit([]interface{}{h})
	if ret := loggingHandlerSetLevel([]interface{}{h, "not a level"}); ret == nil {
		t.Fatalf("expected error for non-object level param")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}

	if ret := loggingHandlerPublish([]interface{}{h}); ret == nil {
		t.Fatalf("expected error when LogRecord param is missing")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}
}
