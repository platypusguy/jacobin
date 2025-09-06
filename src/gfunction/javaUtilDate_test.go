package gfunction

import (
	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
	"time"
)

// Helpers
func newDateObj(t *testing.T) *object.Object {
	// Create an empty Date object with correct class name
	className := "java/util/Date"
	obj := object.MakeEmptyObjectWithClassName(&className)
	if obj == nil {
		t.Fatalf("failed to allocate Date object")
	}
	return obj
}

func assertJavaBoolDate(t *testing.T, got interface{}, want int64, msg string) {
	t.Helper()
	b, ok := got.(int64)
	if !ok {
		if geb, ok := got.(*GErrBlk); ok {
			t.Fatalf("%s: expected Java boolean, got GErrBlk %d (%s)", msg, geb.ExceptionType, geb.ErrMsg)
		}
		t.Fatalf("%s: expected Java boolean (int64), got %T", msg, got)
	}
	if b != want {
		t.Fatalf("%s: expected %d, got %d", msg, want, b)
	}
}

func TestJavaUtilDate_MethodRegistration(t *testing.T) {
	globals.InitStringPool()
	MethodSignatures = make(map[string]GMeth)
	Load_Util_Date()

	cases := []struct{
		key   string
		slots int
	}{
		{"java/util/Date.<init>()V", 0},
		{"java/util/Date.<init>(J)V", 1},
		{"java/util/Date.after(Ljava/util/Date;)Z", 1},
		{"java/util/Date.before(Ljava/util/Date;)Z", 1},
		{"java/util/Date.equals(Ljava/lang/Object;)Z", 1},
		{"java/util/Date.getTime()J", 0},
		{"java/util/Date.setTime(J)V", 1},
		{"java/util/Date.hashCode()I", 0},
		{"java/util/Date.clone()Ljava/lang/Object;", 0},
		{"java/util/Date.toString()Ljava/lang/String;", 0},
	}
	for _, c := range cases {
		gm, ok := MethodSignatures[c.key]
		if !ok {
			t.Fatalf("method not registered: %s", c.key)
		}
		if gm.ParamSlots != c.slots {
			t.Fatalf("ParamSlots mismatch for %s: want %d got %d", c.key, c.slots, gm.ParamSlots)
		}
		if gm.GFunction == nil {
			t.Fatalf("GFunction is nil for %s", c.key)
		}
	}
}

func TestJavaUtilDate_Init_CurrentTime(t *testing.T) {
	globals.InitStringPool()

	obj := newDateObj(t)
	now := time.Now().UnixMilli()
	if ret := udateInit([]interface{}{obj}); ret != nil {
		if geb, ok := ret.(*GErrBlk); ok {
			t.Fatalf("udateInit returned error: %d %s", geb.ExceptionType, geb.ErrMsg)
		}
		t.Fatalf("udateInit returned non-nil: %v", ret)
	}
	// getTime should be within a reasonable window around now
	gt := udateGetTime([]interface{}{obj})
	millis, ok := gt.(int64)
	if !ok {
		t.Fatalf("getTime type: expected int64, got %T", gt)
	}
	// allow for a 10-second skew (generous for CI)
	if millis < now-10_000 || millis > now+10_000 {
		t.Fatalf("constructed time %d not within 10s of now %d", millis, now)
	}
}

func TestJavaUtilDate_InitLong_And_GetSetTime(t *testing.T) {
	globals.InitStringPool()

	obj := newDateObj(t)
	const base int64 = 1_700_000_000_000 // a stable fixed millis value
	if ret := udateInitLong([]interface{}{obj, base}); ret != nil {
		t.Fatalf("udateInitLong returned error: %v", ret)
	}
	if got := udateGetTime([]interface{}{obj}).(int64); got != base {
		t.Fatalf("getTime after initLong: got %d want %d", got, base)
	}
	// setTime then verify
	const newer int64 = base + 12345
	if ret := udateSetTime([]interface{}{obj, newer}); ret != nil {
		t.Fatalf("setTime returned error: %v", ret)
	}
	if got := udateGetTime([]interface{}{obj}).(int64); got != newer {
		t.Fatalf("getTime after setTime: got %d want %d", got, newer)
	}
}

func TestJavaUtilDate_Deprecated_Constructors_Trap(t *testing.T) {
	globals.InitStringPool()
	obj := newDateObj(t)
	cases := []interface{}{
		udateInit3Ints([]interface{}{obj, int64(1), int64(2), int64(3)}),
		udateInit5Ints([]interface{}{obj, int64(1), int64(2), int64(3), int64(4), int64(5)}),
		udateInit6Ints([]interface{}{obj, int64(1), int64(2), int64(3), int64(4), int64(5), int64(6)}),
		udateInitString([]interface{}{obj, object.StringObjectFromGoString("2020-01-01")}),
	}
	for i, ret := range cases {
		geb, ok := ret.(*GErrBlk)
		if !ok {
			t.Fatalf("case %d: expected GErrBlk for deprecated ctor, got %T", i, ret)
		}
		if geb.ExceptionType != excNames.UnsupportedOperationException {
			t.Fatalf("case %d: expected UnsupportedOperationException, got %d", i, geb.ExceptionType)
		}
	}
}

func TestJavaUtilDate_After_Before_Equals(t *testing.T) {
	globals.InitStringPool()
	// d1 < d2
	d1 := newDateObj(t)
	d2 := newDateObj(t)
	_ = udateInitLong([]interface{}{d1, int64(1000)})
	_ = udateInitLong([]interface{}{d2, int64(2000)})

 assertJavaBoolDate(t, udateBefore([]interface{}{d1, d2}), types.JavaBoolTrue, "1000 before 2000")
 assertJavaBoolDate(t, udateAfter([]interface{}{d1, d2}), types.JavaBoolFalse, "1000 after 2000 should be false")
 assertJavaBoolDate(t, udateAfter([]interface{}{d2, d1}), types.JavaBoolTrue, "2000 after 1000")
 assertJavaBoolDate(t, udateBefore([]interface{}{d2, d1}), types.JavaBoolFalse, "2000 before 1000 should be false")

	// equals
 assertJavaBoolDate(t, udateEquals([]interface{}{d1, d1}), types.JavaBoolTrue, "equals self true")
 assertJavaBoolDate(t, udateEquals([]interface{}{d1, d2}), types.JavaBoolFalse, "equals different false")
	var nullObj *object.Object = nil
 assertJavaBoolDate(t, udateEquals([]interface{}{d1, nullObj}), types.JavaBoolFalse, "equals null false")
}

func TestJavaUtilDate_Clone_And_HashCode(t *testing.T) {
	globals.InitStringPool()
	obj := newDateObj(t)
	// Use a deterministic millis to compute expected hash
	const v int64 = 0x0123456789ABCDEF
	_ = udateInitLong([]interface{}{obj, v})
	cl := udateClone([]interface{}{obj})
	clObj, ok := cl.(*object.Object)
	if !ok || clObj == nil {
		t.Fatalf("clone returned invalid object: %T", cl)
	}
	if clObj == obj {
		t.Fatalf("clone should return a distinct object pointer")
	}
	// Millis should be equal
	mv := udateGetTime([]interface{}{clObj}).(int64)
	if mv != v {
		t.Fatalf("clone millis mismatch: got %d want %d", mv, v)
	}
	// HashCode per Java: (int)(v ^ (v >>> 32)) as int64 in our impl
	upper := int64(uint64(v) >> 32)
	x := v ^ upper
	expectedHash := int64(int32(x))
	gotHash := udateHashCode([]interface{}{obj}).(int64)
	if gotHash != expectedHash {
		t.Fatalf("hashCode mismatch: got %d want %d", gotHash, expectedHash)
	}
}

func TestJavaUtilDate_ToString_Sane(t *testing.T) {
	globals.InitStringPool()
	obj := newDateObj(t)
	_ = udateInitLong([]interface{}{obj, int64(1_700_000_000_000)})
	strObj := udateToString([]interface{}{obj})
	so, ok := strObj.(*object.Object)
	if !ok || so == nil {
		t.Fatalf("toString did not return a String object: %T", strObj)
	}
	gs := object.GoStringFromStringObject(so)
	if gs == "" {
		t.Fatalf("toString returned empty string")
	}
	// It should not be our explicit error marker prefix
	if len(gs) >= 5 && gs[:5] == "Date[" {
		t.Fatalf("toString returned an error marker: %q", gs)
	}
}
