package javaUtil

import (
	"io"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"os"
	"strings"
	"testing"
)

func newConsoleHandlerObj() *object.Object {
	return object.MakeEmptyObjectWithClassName(&consoleHandlerClassName)
}

func TestLoggingConsoleHandlerInit_DefaultsToLevelInfo(t *testing.T) {
	globals.InitStringPool()

	ch := newConsoleHandlerObj()
	ret := loggingConsoleHandlerInit([]interface{}{ch})
	if ret != nil {
		t.Fatalf("loggingConsoleHandlerInit returned error: %v", ret)
	}

	level := loggingHandlerGetLevel([]interface{}{ch}).(*object.Object)
	if got := loggingLevelIntValue([]interface{}{level}).(int64); got != standardLevels["INFO"] {
		t.Fatalf("expected default level INFO, got %d", got)
	}
	if ret := loggingHandlerGetFilter([]interface{}{ch}); ret != object.Null {
		t.Fatalf("expected null filter by default")
	}
}

func TestLoggingConsoleHandlerPublish_WritesToSystemErr(t *testing.T) {
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

	ch := newConsoleHandlerObj()
	_ = loggingConsoleHandlerInit([]interface{}{ch})

	record := makeLogRecord("INFO", standardLevels["INFO"], "hello console")
	ret := loggingConsoleHandlerPublish([]interface{}{ch, record})
	if ret != nil {
		t.Fatalf("loggingConsoleHandlerPublish returned error: %v", ret)
	}

	_ = w.Close()
	buf, _ := io.ReadAll(r)
	got := string(buf)
	if !strings.Contains(got, "INFO: hello console") {
		t.Fatalf("unexpected published output: %q", got)
	}
}

func TestLoggingConsoleHandlerPublish_NotLoggable_NoOutput(t *testing.T) {
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

	ch := newConsoleHandlerObj()
	_ = loggingConsoleHandlerInit([]interface{}{ch})
	_ = loggingHandlerSetLevel([]interface{}{ch, makeLevelObject("SEVERE", standardLevels["SEVERE"], "")})

	record := makeLogRecord("INFO", standardLevels["INFO"], "should not appear")
	ret := loggingConsoleHandlerPublish([]interface{}{ch, record})
	if ret != nil {
		t.Fatalf("loggingConsoleHandlerPublish returned error: %v", ret)
	}

	_ = w.Close()
	buf, _ := io.ReadAll(r)
	if len(buf) != 0 {
		t.Fatalf("expected no output for non-loggable record, got %q", string(buf))
	}
}

func TestLoggingConsoleHandlerClose_SyncSystemErr(t *testing.T) {
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

	if ret := loggingConsoleHandlerClose(nil); ret != nil {
		t.Fatalf("loggingConsoleHandlerClose returned error: %v", ret)
	}
}

func TestLoggingConsoleHandler_ErrorPaths(t *testing.T) {
	globals.InitStringPool()

	if ret := loggingConsoleHandlerInit([]interface{}{"not an object"}); ret == nil {
		t.Fatalf("expected error for non-object param")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}

	ch := newConsoleHandlerObj()
	_ = loggingConsoleHandlerInit([]interface{}{ch})
	if ret := loggingConsoleHandlerPublish([]interface{}{ch}); ret == nil {
		t.Fatalf("expected error when LogRecord param is missing")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}
}
