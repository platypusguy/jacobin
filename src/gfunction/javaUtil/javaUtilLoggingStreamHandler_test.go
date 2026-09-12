package javaUtil

import (
	"io"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func newStreamHandlerObj() *object.Object {
	return object.MakeEmptyObjectWithClassName(&streamHandlerClassName)
}

func funcName(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func TestLoad_Util_Logging_StreamHandler_Registration(t *testing.T) {
	globals.InitStringPool()
	Load_Util_Logging_StreamHandler()

	trapEntries := []string{
		"java/util/logging/StreamHandler.<init>(Ljava/io/OutputStream;)V",
		"java/util/logging/StreamHandler.<init>(Ljava/io/OutputStream;Ljava/util/logging/Formatter;)V",
		"java/util/logging/StreamHandler.isLoggable(Ljava/util/logging/LogRecord;)Z",
		"java/util/logging/StreamHandler.setEncoding(Ljava/lang/String;)V",
		"java/util/logging/StreamHandler.setOutputStream(Ljava/io/OutputStream;)V",
	}
	for _, key := range trapEntries {
		gm, ok := ghelpers.MethodSignatures[key]
		if !ok {
			t.Fatalf("expected %s to be registered", key)
		}
		if funcName(gm.GFunction) != funcName(ghelpers.TrapFunction) {
			t.Fatalf("expected %s to use ghelpers.TrapFunction", key)
		}
	}

	implementedEntries := map[string]int{
		"java/util/logging/StreamHandler.<clinit>()V":                             0,
		"java/util/logging/StreamHandler.<init>()V":                               0,
		"java/util/logging/StreamHandler.close()V":                                0,
		"java/util/logging/StreamHandler.flush()V":                                0,
		"java/util/logging/StreamHandler.publish(Ljava/util/logging/LogRecord;)V": 1,
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

func TestLoggingStreamHandlerInit_Defaults(t *testing.T) {
	globals.InitStringPool()

	h := newStreamHandlerObj()
	ret := loggingStreamHandlerInit([]interface{}{h})
	if ret != nil {
		t.Fatalf("loggingStreamHandlerInit returned error: %v", ret)
	}

	level := h.FieldTable[fieldNameHandlerLevel].Fvalue.(*object.Object)
	if got := loggingLevelIntValue([]interface{}{level}).(int64); got != standardLevels["INFO"] {
		t.Fatalf("expected default level INFO, got %d", got)
	}
	if h.FieldTable[fieldNameHandlerFilter].Fvalue != object.Null {
		t.Fatalf("expected null filter by default")
	}
	if h.FieldTable[fieldNameStreamHandlerOutputStream].Fvalue != object.Null {
		t.Fatalf("expected null outputStream by default")
	}
}

func TestLoggingStreamHandlerPublish_WritesToOutputStream(t *testing.T) {
	globals.InitStringPool()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	defer r.Close()

	h := newStreamHandlerObj()
	_ = loggingStreamHandlerInit([]interface{}{h})
	h.FieldTable[fieldNameStreamHandlerOutputStream] = object.Field{Ftype: types.Ref, Fvalue: w}

	record := makeLogRecord("INFO", standardLevels["INFO"], "hello stream")
	ret := loggingStreamHandlerPublish([]interface{}{h, record})
	if ret != nil {
		t.Fatalf("loggingStreamHandlerPublish returned error: %v", ret)
	}

	_ = w.Close()
	buf, _ := io.ReadAll(r)
	got := string(buf)
	if !strings.Contains(got, "INFO: hello stream") {
		t.Fatalf("unexpected published output: %q", got)
	}
}

func TestLoggingStreamHandlerPublish_NotLoggable_NoOutput(t *testing.T) {
	globals.InitStringPool()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	defer r.Close()

	h := newStreamHandlerObj()
	_ = loggingStreamHandlerInit([]interface{}{h})
	h.FieldTable[fieldNameStreamHandlerOutputStream] = object.Field{Ftype: types.Ref, Fvalue: w}
	h.FieldTable[fieldNameHandlerLevel] = object.Field{Ftype: types.Ref, Fvalue: makeLevelObject("SEVERE", standardLevels["SEVERE"], "")}

	record := makeLogRecord("INFO", standardLevels["INFO"], "should not appear")
	ret := loggingStreamHandlerPublish([]interface{}{h, record})
	if ret != nil {
		t.Fatalf("loggingStreamHandlerPublish returned error: %v", ret)
	}

	_ = w.Close()
	buf, _ := io.ReadAll(r)
	if len(buf) != 0 {
		t.Fatalf("expected no output for non-loggable record, got %q", string(buf))
	}
}

func TestLoggingStreamHandlerFlushAndClose(t *testing.T) {
	globals.InitStringPool()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	defer r.Close()

	h := newStreamHandlerObj()
	_ = loggingStreamHandlerInit([]interface{}{h})
	h.FieldTable[fieldNameStreamHandlerOutputStream] = object.Field{Ftype: types.Ref, Fvalue: w}

	if ret := loggingStreamHandlerFlush([]interface{}{h}); ret != nil {
		t.Fatalf("loggingStreamHandlerFlush returned error: %v", ret)
	}
	if ret := loggingStreamHandlerClose([]interface{}{h}); ret != nil {
		t.Fatalf("loggingStreamHandlerClose returned error: %v", ret)
	}
}

func TestLoggingStreamHandlerFlushAndClose_NoStream_NoOp(t *testing.T) {
	globals.InitStringPool()

	h := newStreamHandlerObj()
	_ = loggingStreamHandlerInit([]interface{}{h})

	if ret := loggingStreamHandlerFlush([]interface{}{h}); ret != nil {
		t.Fatalf("expected no-op success, got: %v", ret)
	}
	if ret := loggingStreamHandlerClose([]interface{}{h}); ret != nil {
		t.Fatalf("expected no-op success, got: %v", ret)
	}
}

func TestLoggingStreamHandler_ErrorPaths(t *testing.T) {
	globals.InitStringPool()

	if ret := loggingStreamHandlerInit([]interface{}{"not an object"}); ret == nil {
		t.Fatalf("expected error for non-object param")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}

	if ret := loggingStreamHandlerClose([]interface{}{"not an object"}); ret == nil {
		t.Fatalf("expected error for non-object param")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}

	if ret := loggingStreamHandlerFlush([]interface{}{"not an object"}); ret == nil {
		t.Fatalf("expected error for non-object param")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}

	h := newStreamHandlerObj()
	_ = loggingStreamHandlerInit([]interface{}{h})
	if ret := loggingStreamHandlerPublish([]interface{}{h}); ret == nil {
		t.Fatalf("expected error when LogRecord param is missing")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}

	if ret := loggingStreamHandlerPublish([]interface{}{"not an object", "x"}); ret == nil {
		t.Fatalf("expected error for non-object receiver param")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}
}
