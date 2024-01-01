package statics

import (
	"fmt"
	"io"
	"jacobin/classloader"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/types"
	"os"
	"strings"
	"sync"
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
	log.Init()
	// _ = log.SetLogLevel(log.CLASS)
	StaticsPreload()
	classloader.MethArea = &sync.Map{}
	k := classloader.Klass{}
	k.Status = 'F'
	k.Loader = "application"
	clData := classloader.ClData{}
	clData.Name = "AlphaBetaGamma"
	k.Data = &clData
	classloader.MethAreaInsert("AlphaBetaGamma", &k)
	ref := classloader.MethAreaFetch("AlphaBetaGamma")

	/**
	Set statics values.
	*/
	LoadProgramStatics()
	tAddStatic(t, "AlphaBetaGamma.ONE", "B", 0x31)
	tAddStatic(t, "AlphaBetaGamma.QM", "C", '?')
	tAddStatic(t, "AlphaBetaGamma.PI", "D", 3.14159265)
	tAddStatic(t, "AlphaBetaGamma.TEN", "F", 10.0)
	tAddStatic(t, "AlphaBetaGamma.D-ADAMS", "I", 42)
	tAddStatic(t, "AlphaBetaGamma.BILLION", "J", 2000000000)
	tAddStatic(t, "AlphaBetaGamma.WILLIE", "LAlphaBetaGamma;", ref)
	tAddStatic(t, "AlphaBetaGamma.THIRTEEN", "S", 13)
	tAddStatic(t, "AlphaBetaGamma.TRUE", "Z", true)
	// Omitted: [x
	// Omitted: G
	// Omitted: T to avoid a cycle (object >> statics >> object)

	/**
	Check statics values.
	*/
	tCheckStatic(t, "main", "$assertionsDisabled", int64(1))
	tCheckStatic(t, "java/lang/String", "COMPACT_STRINGS", true)
	tCheckStatic(t, "AlphaBetaGamma", "ONE", int64(0x31))
	tCheckStatic(t, "AlphaBetaGamma", "QM", int64('?'))
	tCheckStatic(t, "AlphaBetaGamma", "PI", float64(3.14159265))
	tCheckStatic(t, "AlphaBetaGamma", "TEN", float64(10.0))
	tCheckStatic(t, "AlphaBetaGamma", "D-ADAMS", int64(42))
	tCheckStatic(t, "AlphaBetaGamma", "BILLION", int64(2000000000))
	tCheckStatic(t, "AlphaBetaGamma", "WILLIE", ref)
	tCheckStatic(t, "AlphaBetaGamma", "THIRTEEN", int64(13))
	tCheckStatic(t, "AlphaBetaGamma", "TRUE", true)
	DumpStatics()

}

func TestInvalidStaticAdd(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	Statics = make(map[string]Static)

	err := AddStatic("", Static{})
	if !strings.Contains(err.Error(), "Attempting to add static entry with a nil name") {
		t.Errorf("TestInvalidStaticAdd: got unexpected error message: %s\n", err.Error())
	}
}

func TestIntConversions(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
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
	log.Init()
	Statics = make(map[string]Static)

	StaticsPreload()
	s1 := GetStaticValue("java/lang/String", "COMPACT_STRINGS")
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

func TestDumpStatics(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	Statics = make(map[string]Static)

	err1 := AddStatic("test.1", Static{Type: types.Byte, Value: 'B'})
	err2 := AddStatic("test.2", Static{Type: types.Int, Value: int(42)})
	err3 := AddStatic("test.3", Static{Type: types.Double, Value: 24.0})
	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("TestIntConversions: got unexpected error adding statics for testing")
	}

	// redirect stderr, to avoid all the error msgs for a non-existent class
	normalStderr := os.Stderr
	rerr, werr, _ := os.Pipe()
	os.Stderr = werr

	DumpStatics()

	_ = werr.Close()
	out, _ := io.ReadAll(rerr)
	os.Stderr = normalStderr
	contents := string(out[:])

	os.Stderr = normalStderr

	if !strings.Contains(contents, "test.1") || !strings.Contains(contents, "test.2") || !strings.Contains(contents, "test.3") {
		t.Errorf("TestIntConversions: got unexpected error in DumpStatics: %s", contents)
	}
}
