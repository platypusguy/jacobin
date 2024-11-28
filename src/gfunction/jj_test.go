package gfunction

import (
	"io"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/statics"
	"jacobin/trace"
	"jacobin/types"
	"os"
	"strings"
	"testing"
)

func fnTestJjDumpStatics(t *testing.T, selection int64, className string, threesome []string) {
	// Re-direct stderr.
	originalStderr := os.Stderr
	rerr, werr, _ := os.Pipe()
	os.Stderr = werr

	// Dump statics.
	objTitle := object.StringObjectFromGoString("TestDumpStatics")
	objClassName := object.StringObjectFromGoString(className)
	params := make([]interface{}, 3)
	params[0] = objTitle
	params[1] = selection
	params[2] = objClassName
	jjDumpStatics(params)

	// Close the working stderr, capture its contents, and restore the original stderr.
	_ = werr.Close()
	bytes, _ := io.ReadAll(rerr)
	contents := string(bytes[:])
	os.Stderr = originalStderr

	if !strings.Contains(contents, threesome[0]) || !strings.Contains(contents, threesome[1]) || !strings.Contains(contents, threesome[2]) {
		t.Errorf("fnTestDumpStatics(%d, \"%s\"): looking for these: %v", selection, className, threesome)
		t.Errorf("fnTestDumpStatics(%d, \"%s\"): didn't see them in DumpStatics output: %s", selection, className, contents)
	}

}

func TestJjDumpStatics(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	statics.Statics = make(map[string]statics.Static)

	err1 := statics.AddStatic("test.f1", statics.Static{Type: types.Byte, Value: 'B'})
	err2 := statics.AddStatic("test.f2", statics.Static{Type: types.Int, Value: int(42)})
	err3 := statics.AddStatic("test.f3", statics.Static{Type: types.Double, Value: 24.0})
	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("TestIntConversions: got unexpected error adding statics for testing")
	}

	fnTestJjDumpStatics(t, statics.SelectAll, "", []string{"test.f1", "test.f2", "test.f3"})
	fnTestJjDumpStatics(t, statics.SelectUser, "", []string{"test.f1", "test.f2", "test.f3"})
	fnTestJjDumpStatics(t, statics.SelectClass, "test", []string{"test.f1", "test.f2", "test.f3"})
}
