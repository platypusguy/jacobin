package statics

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/globals"
	"jacobin/log"
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
			expValue = int64(1)
		} else {
			expValue = int64(0)
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
	t.Log("Try a nil statics name. Expecting AddStatic to complain and return an error.")
	tAddStatic(t, "", "rubbishType", "rubbishValue")

	/**
	Check statics values.
	*/
	tCheckStatic(t, "main", "$assertionsDisabled", int64(1))
	tCheckStatic(t, "java/lang/String", "COMPACT_STRINGS", true)
	tCheckStatic(t, "AlphaBetaGamma", "ONE", int64(0x31))
	tCheckStatic(t, "AlphaBetaGamma", "QM", '?')
	tCheckStatic(t, "AlphaBetaGamma", "PI", float64(3.14159265))
	tCheckStatic(t, "AlphaBetaGamma", "TEN", float64(10.0))
	tCheckStatic(t, "AlphaBetaGamma", "D-ADAMS", int64(42))
	tCheckStatic(t, "AlphaBetaGamma", "BILLION", int64(2000000000))
	tCheckStatic(t, "AlphaBetaGamma", "WILLIE", ref)
	tCheckStatic(t, "AlphaBetaGamma", "THIRTEEN", int64(13))
	tCheckStatic(t, "AlphaBetaGamma", "TRUE", true)
	DumpStatics()

}
