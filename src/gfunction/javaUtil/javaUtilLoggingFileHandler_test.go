package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newFileHandlerObj() *object.Object {
	return object.MakeEmptyObjectWithClassName(&fileHandlerClassName)
}

func tempPattern(t *testing.T, name string) string {
	return filepath.Join(t.TempDir(), name)
}

func TestLoggingFileHandlerInit_Default(t *testing.T) {
	globals.InitStringPool()

	// cwd-relative default pattern -- confine it to a temp dir via os.Chdir.
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldWd) }()
	_ = os.Chdir(dir)

	h := newFileHandlerObj()
	ret := loggingFileHandlerInit([]interface{}{h})
	if ret != nil {
		t.Fatalf("loggingFileHandlerInit returned error: %v", ret)
	}

	ret = loggingFileHandlerClose([]interface{}{h})
	if ret != nil {
		t.Fatalf("loggingFileHandlerClose returned error: %v", ret)
	}
}

func TestLoggingFileHandlerInitPattern(t *testing.T) {
	globals.InitStringPool()

	path := tempPattern(t, "test1.log")
	h := newFileHandlerObj()
	patternObj := object.StringObjectFromGoString(path)
	ret := loggingFileHandlerInitPattern([]interface{}{h, patternObj})
	if ret != nil {
		t.Fatalf("loggingFileHandlerInitPattern returned error: %v", ret)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file %s to exist: %v", path, err)
	}

	ret = loggingFileHandlerClose([]interface{}{h})
	if ret != nil {
		t.Fatalf("loggingFileHandlerClose returned error: %v", ret)
	}
}

func TestLoggingFileHandlerInitPatternAppend(t *testing.T) {
	globals.InitStringPool()

	path := tempPattern(t, "test2.log")

	// Write once, close.
	h1 := newFileHandlerObj()
	patternObj := object.StringObjectFromGoString(path)
	_ = loggingFileHandlerInitPatternAppend([]interface{}{h1, patternObj, int64(0)})
	record1 := makeLogRecord("INFO", standardLevels["INFO"], "line1")
	_ = loggingFileHandlerPublish([]interface{}{h1, record1})
	_ = loggingFileHandlerClose([]interface{}{h1})

	// Re-open with append=true and write again.
	h2 := newFileHandlerObj()
	_ = loggingFileHandlerInitPatternAppend([]interface{}{h2, patternObj, int64(1)})
	record2 := makeLogRecord("INFO", standardLevels["INFO"], "line2")
	_ = loggingFileHandlerPublish([]interface{}{h2, record2})
	_ = loggingFileHandlerClose([]interface{}{h2})

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile: %v", err)
	}
	got := string(data)
	if !strings.Contains(got, "INFO: line1") || !strings.Contains(got, "INFO: line2") {
		t.Fatalf("expected appended content, got %q", got)
	}
}

func TestLoggingFileHandlerInitPatternLimitCount(t *testing.T) {
	globals.InitStringPool()

	path := tempPattern(t, "test3.log")
	h := newFileHandlerObj()
	patternObj := object.StringObjectFromGoString(path)
	ret := loggingFileHandlerInitPatternLimitCount([]interface{}{h, patternObj, int64(1000), int64(2)})
	if ret != nil {
		t.Fatalf("loggingFileHandlerInitPatternLimitCount returned error: %v", ret)
	}
	_ = loggingFileHandlerClose([]interface{}{h})

	// Invalid: count < 1
	h2 := newFileHandlerObj()
	ret2 := loggingFileHandlerInitPatternLimitCount([]interface{}{h2, patternObj, int64(1000), int64(0)})
	if _, ok := ret2.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected error for count < 1, got %v", ret2)
	}
}

func TestLoggingFileHandlerInitPatternLimitCountAppend(t *testing.T) {
	globals.InitStringPool()

	path := tempPattern(t, "test4.log")
	h := newFileHandlerObj()
	patternObj := object.StringObjectFromGoString(path)
	ret := loggingFileHandlerInitPatternLimitCountAppend([]interface{}{h, patternObj, int64(1000), int64(1), int64(0)})
	if ret != nil {
		t.Fatalf("loggingFileHandlerInitPatternLimitCountAppend returned error: %v", ret)
	}
	_ = loggingFileHandlerClose([]interface{}{h})
}

func TestLoggingFileHandlerPublish_WritesToFile(t *testing.T) {
	globals.InitStringPool()

	path := tempPattern(t, "test5.log")
	h := newFileHandlerObj()
	patternObj := object.StringObjectFromGoString(path)
	_ = loggingFileHandlerInitPattern([]interface{}{h, patternObj})

	record := makeLogRecord("INFO", standardLevels["INFO"], "hello file")
	ret := loggingFileHandlerPublish([]interface{}{h, record})
	if ret != nil {
		t.Fatalf("loggingFileHandlerPublish returned error: %v", ret)
	}
	_ = loggingFileHandlerClose([]interface{}{h})

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile: %v", err)
	}
	if got := string(data); !strings.Contains(got, "INFO: hello file") {
		t.Fatalf("unexpected file content: %q", got)
	}
}

func TestLoggingFileHandlerPublish_NotLoggable_NoOutput(t *testing.T) {
	globals.InitStringPool()

	path := tempPattern(t, "test6.log")
	h := newFileHandlerObj()
	patternObj := object.StringObjectFromGoString(path)
	_ = loggingFileHandlerInitPattern([]interface{}{h, patternObj})
	_ = loggingHandlerSetLevel([]interface{}{h, makeLevelObject("SEVERE", standardLevels["SEVERE"], "")})

	record := makeLogRecord("INFO", standardLevels["INFO"], "should not appear")
	ret := loggingFileHandlerPublish([]interface{}{h, record})
	if ret != nil {
		t.Fatalf("loggingFileHandlerPublish returned error: %v", ret)
	}
	_ = loggingFileHandlerClose([]interface{}{h})

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile: %v", err)
	}
	if got := string(data); got != "" {
		t.Fatalf("expected no output, got %q", got)
	}
}

func TestLoggingFileHandler_ErrorPaths(t *testing.T) {
	globals.InitStringPool()

	if ret := loggingFileHandlerInit([]interface{}{"not an object"}); ret == nil {
		t.Fatalf("expected error for non-object param")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}

	h := newFileHandlerObj()
	if ret := loggingFileHandlerInitPattern([]interface{}{h, object.Null}); ret == nil {
		t.Fatalf("expected error for null pattern")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}

	if ret := loggingFileHandlerClose([]interface{}{"not an object"}); ret == nil {
		t.Fatalf("expected error for non-object param on close")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}

	if ret := loggingFileHandlerPublish([]interface{}{"not an object", "x"}); ret == nil {
		t.Fatalf("expected error for non-object param on publish")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}

	h2 := newFileHandlerObj()
	if ret := loggingFileHandlerPublish([]interface{}{h2}); ret == nil {
		t.Fatalf("expected error for missing LogRecord parameter")
	} else if geb, ok := ret.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("expected IllegalArgumentException")
	}
}
