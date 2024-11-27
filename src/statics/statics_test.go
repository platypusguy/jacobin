package statics

import (
	"fmt"
	"io"
	// "jacobin/classloader"
	"jacobin/globals"
	"jacobin/trace"
	"jacobin/types"
	"os"
	"strings"
	"testing"
)

func tAddStatic(t *testing.T, sname string, stype string, svalue any) {
	ss := Static{stype, svalue}
	err := AddStatic(sname, ss)
	if sname == "" {
		// sname is nil (abnormal)
		if err == nil {
			t.Errorf("tAddStatic: AddStatic failed to diagnose nil name, type=%s, value=%v\n", stype, svalue)
		} else {
			t.Log("Successful diagnosis of a nil statics name.")
		}
		return
	}
	// NoN-NIL VALUE FOR SNAME
	if err != nil {
		t.Errorf("tAddStatic: AddStatic(name=%s, type=%s, value=%v) failed, err=%s\n", sname, stype, svalue, err.Error())
	}
}

func tCheckStatic(t *testing.T, className string, fieldName string, expValue any) {
	retValue := GetStaticValue(className, fieldName)
	switch retValue.(type) {
	case error:
		t.Errorf("tGetStatic: statics.GetStaticValue diagnosed an error (previous message)\n")
		return
	}
	switch expValue.(type) {
	case bool:
		if expValue.(bool) {
			expValue = types.JavaBoolTrue
		} else {
			expValue = types.JavaBoolFalse
		}
	}
	if expValue == retValue {
		fmt.Printf("tGetStatic: %s.%s --> %v\n", className, fieldName, retValue)
	} else {
		t.Errorf("tGetStatic: Expected GetStaticValue(%s.%s) return is %v.(%T) but observed %v.(%T)\n",
			className, fieldName, expValue, expValue, retValue, retValue)
	}
}

func TestStatics1(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	Statics = make(map[string]Static)
	/***
	PreloadStatics()
	classloader.MethArea = &sync.Map{}
	k := classloader.Klass{}
	k.Status = 'F'
	k.Loader = "application"
	clData := classloader.ClData{}
	clData.Name = "AlphaBetaGamma"
	k.Data = &clData
	classloader.MethAreaInsert("AlphaBetaGamma", &k)
	ref := classloader.MethAreaFetch("AlphaBetaGamma")
	***/
	ii := 42
	ref := &ii

	/**
	Set statics values.
	*/
	LoadProgramStatics()
	tAddStatic(t, "AlphaBetaGamma.ONE", types.Byte, uint8(0x31))
	tAddStatic(t, "AlphaBetaGamma.QM", types.Char, '?')
	tAddStatic(t, "AlphaBetaGamma.PI", types.Double, 3.14159265)
	tAddStatic(t, "AlphaBetaGamma.TEN", types.Float, 10.0)
	tAddStatic(t, "AlphaBetaGamma.D-ADAMS", types.Int, 42)
	tAddStatic(t, "AlphaBetaGamma.BILLION", types.Long, 2000000000)
	tAddStatic(t, "AlphaBetaGamma.WILLIE", "LAlphaBetaGamma;", ref)
	tAddStatic(t, "AlphaBetaGamma.THIRTEEN", types.Short, 13)
	tAddStatic(t, "AlphaBetaGamma.TRUE", types.Bool, true)
	// Omitted: [x
	// Omitted: G
	// Omitted: T to avoid a cycle (object >> statics >> object)

	/**
	Check statics values.
	*/
	tCheckStatic(t, "main", "$assertionsDisabled", int64(1))
	tCheckStatic(t, "AlphaBetaGamma", "ONE", int64(0x31))
	tCheckStatic(t, "AlphaBetaGamma", "QM", int64('?'))
	tCheckStatic(t, "AlphaBetaGamma", "PI", float64(3.14159265))
	tCheckStatic(t, "AlphaBetaGamma", "TEN", float64(10.0))
	tCheckStatic(t, "AlphaBetaGamma", "D-ADAMS", int64(42))
	tCheckStatic(t, "AlphaBetaGamma", "BILLION", int64(2000000000))
	tCheckStatic(t, "AlphaBetaGamma", "WILLIE", ref)
	tCheckStatic(t, "AlphaBetaGamma", "THIRTEEN", int64(13))
	tCheckStatic(t, "AlphaBetaGamma", "TRUE", true)
	DumpStatics("TestStatics1", SelectUser, "")

}

func TestInvalidStaticAdd(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	Statics = make(map[string]Static)

	err := AddStatic("", Static{})
	if !strings.Contains(err.Error(), "Attempting to add static entry with a nil name") {
		t.Errorf("TestInvalidStaticAdd: got unexpected error message: %s\n", err.Error())
	}
}

func TestInvalidLookup(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	Statics = make(map[string]Static)

	err1 := AddStatic("test.1", Static{Type: types.Int, Value: int(42)})
	if err1 != nil {
		t.Errorf("TestIntConversions: got unexpected error adding static for testing")
	}

	// redirect stderr, to avoid all the error msgs for a non-existent class
	normalStderr := os.Stderr
	_, werr, _ := os.Pipe()
	os.Stderr = werr

	retVal := GetStaticValue("test", "noSuchEntry")

	_ = werr.Close()
	os.Stderr = normalStderr

	switch retVal.(type) {
	case error:
		if !strings.Contains(retVal.(error).Error(), "could not find static") {
			t.Errorf("Did not get expected error message for missing static: %v\n", retVal)
		}
	default:
		t.Errorf("Did not get an error for missing static")
	}
}

func TestIntConversions(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	Statics = make(map[string]Static)

	err1 := AddStatic("test.1", Static{Type: types.Byte, Value: 'B'})
	err2 := AddStatic("test.2", Static{Type: types.Int, Value: int(42)})
	err3 := AddStatic("test.3", Static{Type: types.Double, Value: 24.0})
	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("TestIntConversions: got unexpected error adding statics for testing")
	}

	retValue1 := GetStaticValue("test", "1")
	switch retValue1.(type) {
	case int64:
		if retValue1.(int64) != 'B' {
			t.Errorf("TestIntConversions: Expected 'B' but observed %v\n", retValue1)
		}
	default:
		t.Errorf("TestIntConversions: invalid type for test.1: %T\n", retValue1)
	}

	retValue2 := GetStaticValue("test", "2")
	switch retValue1.(type) {
	case int64:
		if retValue2.(int64) != 42 {
			t.Errorf("TestIntConversions: Expected 42 but observed %d\n", retValue2)
		}
	default:
		t.Errorf("TestIntConversions: invalid type for test.2: %T\n", retValue2)
	}

	retValue3 := GetStaticValue("test", "3")
	switch retValue3.(type) {
	case float64:
		if retValue3.(float64) != 24.0 {
			t.Errorf("TestIntConversions: Expected 24.0 but observed %fv\n", retValue3)
		}
	default:
		t.Errorf("TestIntConversions: invalid type for test.3: %T\n", retValue3)
	}
}

func TestStaticsPreload(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	Statics = make(map[string]Static)

	PreloadStatics()
	s1 := GetStaticValue(types.StringClassName, "COMPACT_STRINGS")
	switch s1.(type) {
	case int64:
		if s1.(int64) != types.JavaBoolTrue {
			t.Errorf("testStaticsPreload: Expected COMPACT_STRINGS to be true but observed false\n")
		}
	default:
		t.Errorf("testStaticsPreload: invalid value for java/lang/String.COMPACT_STRINGS")
	}

	s2 := GetStaticValue("main", "$assertionsDisabled")
	switch s2.(type) {
	case int64:
		if s2.(int64) != types.JavaBoolTrue {
			t.Errorf("testStaticsPreload: Expected main.$assertionsDisabled to be true but observed false\n")
		}
	default:
		t.Errorf("testStaticsPreload: invalid value for main.$assertionsDisabled")
	}
}

func fnTestDumpStatics(t *testing.T, selection int64, className string, threesome []string) {
	// Re-direct stderr.
	originalStderr := os.Stderr
	rerr, werr, _ := os.Pipe()
	os.Stderr = werr

	// Dump statics.
	DumpStatics("TestDumpStatics", selection, className)

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

func TestDumpStatics(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	Statics = make(map[string]Static)

	err1 := AddStatic("test.f1", Static{Type: types.Byte, Value: 'B'})
	err2 := AddStatic("test.f2", Static{Type: types.Int, Value: int(42)})
	err3 := AddStatic("test.f3", Static{Type: types.Double, Value: 24.0})
	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("TestIntConversions: got unexpected error adding statics for testing")
	}

	fnTestDumpStatics(t, SelectAll, "", []string{"test.f1", "test.f2", "test.f3"})
	fnTestDumpStatics(t, SelectUser, "", []string{"test.f1", "test.f2", "test.f3"})
	fnTestDumpStatics(t, SelectClass, "test", []string{"test.f1", "test.f2", "test.f3"})
}
