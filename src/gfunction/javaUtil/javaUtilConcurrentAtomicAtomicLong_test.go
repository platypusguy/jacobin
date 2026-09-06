package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"runtime"
	"testing"
)

func newAtomicLongObj() *object.Object {
	className := "java/util/concurrent/atomic/AtomicLong"
	return object.MakeEmptyObjectWithClassName(&className)
}

func alGet(t *testing.T, obj *object.Object) int64 {
	t.Helper()
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	f, ok := obj.FieldTable["value"]
	if !ok {
		t.Fatalf("expected 'value' field on AtomicLong object")
	}
	v, ok := f.Fvalue.(int64)
	if !ok {
		t.Fatalf("expected 'value' to be int64, got %T (%v)", f.Fvalue, f.Fvalue)
	}
	return v
}

func TestAtomicLong_Init_And_Set_Get(t *testing.T) {
	globals.InitStringPool()

	al := newAtomicLongObj()
	ret := atomicLongInitVoid([]interface{}{al})
	if ret != nil {
		t.Fatalf("expected initVoid to return nil, got %v", ret)
	}
	if v := alGet(t, al); v != 0 {
		t.Fatalf("expected initial value 0, got %d", v)
	}

	ret = atomicLongInitLong([]interface{}{al, int64(123456789012345)})
	if ret != nil {
		t.Fatalf("expected initLong to return nil, got %v", ret)
	}
	if v := alGet(t, al); v != 123456789012345 {
		t.Fatalf("expected initial value 123456789012345, got %d", v)
	}

	ret = atomicLongSet([]interface{}{al, int64(-98765432109876)})
	if ret != nil {
		t.Fatalf("expected set to return nil, got %v", ret)
	}
	if v := atomicLongGet([]interface{}{al}); v != int64(-98765432109876) {
		t.Fatalf("expected get -98765432109876, got %v", v)
	}
}

func TestAtomicLong_GetAndSet_And_CompareAndSet(t *testing.T) {
	globals.InitStringPool()

	al := newAtomicLongObj()
	_ = atomicLongInitLong([]interface{}{al, int64(100)})

	old := atomicLongGetAndSet([]interface{}{al, int64(200)})
	if old != int64(100) {
		t.Fatalf("expected old value 100, got %v", old)
	}
	if cur := alGet(t, al); cur != 200 {
		t.Fatalf("expected current value 200, got %d", cur)
	}

	casFail := atomicLongCompareAndSet([]interface{}{al, int64(999), int64(300)})
	if casFail != types.JavaBoolFalse {
		t.Fatalf("expected CAS fail (JavaBoolFalse), got %v", casFail)
	}
	if cur := alGet(t, al); cur != 200 {
		t.Fatalf("expected value unchanged at 200, got %d", cur)
	}

	casSuccess := atomicLongCompareAndSet([]interface{}{al, int64(200), int64(300)})
	if casSuccess != types.JavaBoolTrue {
		t.Fatalf("expected CAS success (JavaBoolTrue), got %v", casSuccess)
	}
	if cur := alGet(t, al); cur != 300 {
		t.Fatalf("expected value updated to 300, got %d", cur)
	}
}

func TestAtomicLong_IncDec_Add_Variants(t *testing.T) {
	globals.InitStringPool()

	al := newAtomicLongObj()
	_ = atomicLongInitLong([]interface{}{al, int64(10)})

	if v := atomicLongGetAndIncrement([]interface{}{al}); v != int64(10) {
		t.Fatalf("expected getAndIncrement 10, got %v", v)
	}
	if cur := alGet(t, al); cur != 11 {
		t.Fatalf("expected current 11, got %d", cur)
	}

	if v := atomicLongIncrementAndGet([]interface{}{al}); v != int64(12) {
		t.Fatalf("expected incrementAndGet 12, got %v", v)
	}

	if v := atomicLongGetAndDecrement([]interface{}{al}); v != int64(12) {
		t.Fatalf("expected getAndDecrement 12, got %v", v)
	}
	if cur := alGet(t, al); cur != 11 {
		t.Fatalf("expected current 11, got %d", cur)
	}

	if v := atomicLongDecrementAndGet([]interface{}{al}); v != int64(10) {
		t.Fatalf("expected decrementAndGet 10, got %v", v)
	}

	if v := atomicLongGetAndAdd([]interface{}{al, int64(5)}); v != int64(10) {
		t.Fatalf("expected getAndAdd 10, got %v", v)
	}
	if cur := alGet(t, al); cur != 15 {
		t.Fatalf("expected current 15, got %d", cur)
	}

	if v := atomicLongAddAndGet([]interface{}{al, int64(-20)}); v != int64(-5) {
		t.Fatalf("expected addAndGet -5, got %v", v)
	}
}

func TestAtomicLong_Conversations(t *testing.T) {
	globals.InitStringPool()

	al := newAtomicLongObj()
	// Set to a value that exceeds 32 bits: 0x100000005L = 4294967301
	_ = atomicLongInitLong([]interface{}{al, int64(4294967301)})

	// intValue should truncate to 32 bits signed: int32(4294967301) = 5
	iv := atomicLongToInt([]interface{}{al}).(int64)
	if iv != 5 {
		t.Fatalf("expected intValue to truncate to 5, got %d", iv)
	}

	fv := atomicLongToFloat([]interface{}{al}).(float32)
	if fv != float32(4294967301) {
		t.Fatalf("expected floatValue float32, got %v", fv)
	}

	dv := atomicLongToDouble([]interface{}{al}).(float64)
	if dv != float64(4294967301) {
		t.Fatalf("expected doubleValue float64, got %v", dv)
	}

	strObj := atomicLongToString([]interface{}{al}).(*object.Object)
	if strObj == nil {
		t.Fatalf("expected non-nil string object")
	}
	goStr := object.GoStringFromStringObject(strObj)
	if goStr != "4294967301" {
		t.Fatalf("expected toString '4294967301', got %q", goStr)
	}
}

func TestAtomicLong_MethodSignatures(t *testing.T) {
	globals.InitStringPool()
	Load_Util_Concurrent_Atomic_Atomic_Long()

	sigGetOpaque := "java/util/concurrent/atomic/AtomicLong.getOpaque()J"
	if _, ok := ghelpers.MethodSignatures[sigGetOpaque]; !ok {
		t.Fatalf("missing method signature: %s", sigGetOpaque)
	}

	sigWeakVolatile := "java/util/concurrent/atomic/AtomicLong.weakCompareAndSetVolatile(JJ)Z"
	if _, ok := ghelpers.MethodSignatures[sigWeakVolatile]; !ok {
		t.Fatalf("missing method signature: %s", sigWeakVolatile)
	}

	sigDouble := "java/util/concurrent/atomic/AtomicLong.doubleValue()D"
	if _, ok := ghelpers.MethodSignatures[sigDouble]; !ok {
		t.Fatalf("missing method signature: %s", sigDouble)
	}

	sigInt := "java/util/concurrent/atomic/AtomicLong.intValue()I"
	if _, ok := ghelpers.MethodSignatures[sigInt]; !ok {
		t.Fatalf("missing method signature: %s", sigInt)
	}

}

func TestAtomicLong_VMSupportsCS8_ReflectsArchitecture(t *testing.T) {
	globals.InitStringPool()

	// Determine expected support based on the same architecture map used in the implementation
	arch := runtime.GOARCH
	supportedArchitectures := map[string]bool{
		"amd64":    true,
		"arm64":    true,
		"ppc64":    true,
		"ppc64le":  true,
		"s390x":    true,
		"sparc64":  true,
		"mips64":   true,
		"mips64le": true,
	}
	expected := supportedArchitectures[arch]

	// Invoke the gfunction
	res := atomicLongVMSupportsCS8([]interface{}{})

	// Ensure it returns a Java boolean in int64 form
	val, ok := res.(int64)
	if !ok {
		t.Fatalf("VMSupportsCS8 did not return int64 (Java boolean), got %T", res)
	}

	// Compare to expected mapping
	if expected {
		if val != types.JavaBoolTrue {
			t.Fatalf("VMSupportsCS8 expected true for arch %s, got %d", arch, val)
		}
	} else {
		if val != types.JavaBoolFalse {
			t.Fatalf("VMSupportsCS8 expected false for arch %s, got %d", arch, val)
		}
	}
}
