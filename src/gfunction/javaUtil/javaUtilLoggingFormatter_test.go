package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"testing"
)

func newFormatterObj() *object.Object {
	return object.MakeEmptyObjectWithClassName(&formatterClassName)
}

func TestLoad_Util_Logging_Formatter_Registration(t *testing.T) {
	globals.InitStringPool()
	Load_Util_Logging_Formatter()

	implementedEntries := map[string]int{
		"java/util/logging/Formatter.<clinit>()V":                                                0,
		"java/util/logging/Formatter.<init>()V":                                                   0,
		"java/util/logging/Formatter.format(Ljava/util/logging/LogRecord;)Ljava/lang/String;":     1,
		"java/util/logging/Formatter.formatMessage(Ljava/util/logging/LogRecord;)Ljava/lang/String;": 1,
		"java/util/logging/Formatter.getHead(Ljava/util/logging/Handler;)Ljava/lang/String;":      1,
		"java/util/logging/Formatter.getTail(Ljava/util/logging/Handler;)Ljava/lang/String;":      1,
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

func TestLoggingFormatterInit(t *testing.T) {
	globals.InitStringPool()

	f := object.MakeEmptyObject()
	if ret := loggingFormatterInit([]interface{}{f}); ret != nil {
		t.Fatalf("loggingFormatterInit returned error: %v", ret)
	}
	if object.GoStringFromStringPoolIndex(f.KlassName) != formatterClassName {
		t.Fatalf("expected class name %s, got %s", formatterClassName, object.GoStringFromStringPoolIndex(f.KlassName))
	}
}

func TestLoggingFormatterFormatMessage(t *testing.T) {
	globals.InitStringPool()

	record := makeLogRecord("INFO", standardLevels["INFO"], "hello world")
	ret := loggingFormatterFormatMessage([]interface{}{record})
	msgObj, ok := ret.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", ret)
	}
	if got := object.GoStringFromStringObject(msgObj); got != "hello world" {
		t.Fatalf("expected 'hello world', got %q", got)
	}
}

func TestLoggingFormatterFormat(t *testing.T) {
	globals.InitStringPool()

	f := newFormatterObj()
	_ = loggingFormatterInit([]interface{}{f})

	levelObj := makeLevelObject("SEVERE", standardLevels["SEVERE"], "")
	record := object.MakeEmptyObject()
	record.FieldTable[fieldNameHandlerLevel] = object.Field{Fvalue: levelObj}
	record.FieldTable[fieldNameLogRecordMessage] = object.Field{Fvalue: object.StringObjectFromGoString("boom")}
	record.FieldTable[fieldNameLogRecordLoggerName] = object.Field{Fvalue: object.StringObjectFromGoString("myLogger")}

	ret := loggingFormatterFormat([]interface{}{record})
	strObj, ok := ret.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", ret)
	}
	got := object.GoStringFromStringObject(strObj)
	want := "myLogger SEVERE: boom\n"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestLoggingFormatterGetHeadAndTail(t *testing.T) {
	globals.InitStringPool()

	h := newHandlerObj()
	_ = loggingHandlerInit([]interface{}{h})

	headRet := loggingFormatterGetHead([]interface{}{h})
	headObj, ok := headRet.(*object.Object)
	if !ok || object.GoStringFromStringObject(headObj) != "" {
		t.Fatalf("expected empty head string, got %v", headRet)
	}

	tailRet := loggingFormatterGetTail([]interface{}{h})
	tailObj, ok := tailRet.(*object.Object)
	if !ok || object.GoStringFromStringObject(tailObj) != "" {
		t.Fatalf("expected empty tail string, got %v", tailRet)
	}
}

func TestLoggingFormatter_ErrorPaths(t *testing.T) {
	globals.InitStringPool()

	if ret := loggingFormatterInit([]interface{}{"not an object"}); ret == nil {
		t.Fatalf("expected error for non-object param")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}

	if ret := loggingFormatterFormat([]interface{}{"not an object"}); ret == nil {
		t.Fatalf("expected error for non-object param")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}

	if ret := loggingFormatterFormatMessage([]interface{}{"not an object"}); ret == nil {
		t.Fatalf("expected error for non-object param")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}

	if ret := loggingFormatterGetHead([]interface{}{"not an object"}); ret == nil {
		t.Fatalf("expected error for non-object param")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}

	if ret := loggingFormatterGetTail([]interface{}{"not an object"}); ret == nil {
		t.Fatalf("expected error for non-object param")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}
}
