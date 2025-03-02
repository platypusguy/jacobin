/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import (
	// "io"
	"jacobin/classloader"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/opcodes"
	"jacobin/stringPool"
	"jacobin/testutil"
	"jacobin/trace"
	"jacobin/types"
	"strings"
	"testing"
	"time"
	"unsafe"
)

// tests for runUtils.go. Note that most functions are tested inside the tests for run.go,
// but several benefit from standalone testing. Those are tested here

func TestConvertBoolByteToInt64(t *testing.T) {
	var bite = byte(0x01)
	res := convertInterfaceToInt64(bite)
	if res != 1 {
		t.Errorf("convertBoolByteToInt64(byte), expected = 1, got %d", res)
	}

	yesNo := true
	if convertInterfaceToInt64(yesNo) != 1 {
		t.Errorf("convertBoolByteToInt64(bool) != 1 (true), got %d", res)
	}
}

// test conversion of valid interfaces representing numeric values to int64
func TestConvertRemainingUntestedTypesToInt64(t *testing.T) {
	globals.InitGlobals("test")

	i8 := int8(42)
	val := convertInterfaceToInt64(i8)
	if val != 42 {
		t.Errorf("convertInterfaceToInt64(int8), expected = 42, got %d", val)
	}

	u8 := int8(42)
	val = convertInterfaceToInt64(u8)
	if val != 42 {
		t.Errorf("convertInterfaceToInt64(int8), expected = 42, got %d", val)
	}

	i16 := int16(-42)
	val = convertInterfaceToInt64(i16)
	if val != -42 {
		t.Errorf("convertInterfaceToInt64(int16), expected = -42, got %d", val)
	}

	u16 := uint16(142)
	val = convertInterfaceToInt64(u16)
	if val != 142 {
		t.Errorf("convertInterfaceToInt64(uint16), expected = 142, got %d", val)
	}

	i := int(-1042)
	val = convertInterfaceToInt64(i)
	if val != -1042 {
		t.Errorf("convertInterfaceToInt6(int), expected = -1042, got %d", val)
	}

	i32 := int32(-104232)
	val = convertInterfaceToInt64(i32)
	if val != -104232 {
		t.Errorf("convertInterfaceToInt64(int32), expected = -104232, got %d", val)
	}

	u32 := uint32(104232)
	val = convertInterfaceToInt64(u32)
	if val != 104232 {
		t.Errorf("convertInterfaceToInt64(uint32), expected = 104232, got %d", val)
	}
}

// test conversion of invalid type to int64
func TestConvertInvalidTypeToInt64(t *testing.T) {
	globals.InitGlobals("test")
	val := convertInterfaceToInt64(nil)
	if val != 0 {
		t.Errorf("convertInterfaceToInt64, expected = 0, got %d", val)
	}
}

// convert to uint64
func TestConvertFloatToInt64RoundDown(t *testing.T) {
	f := float64(5432.10)
	val := convertInterfaceToUint64(f)
	if val != 5432 {
		t.Errorf("convertFloatToInt64(float64), expected = 5432, got %d", val)
	}

}

func TestConvertFloatToInt64RoundUp(t *testing.T) {
	f := float64(5432.501)
	val := convertInterfaceToUint64(f)
	if val != 5433 {
		t.Errorf("convertFloatToInt64(float64), expected = 5433, got %d", val)
	}
}

// golang bytes are unsigned 8-bit fields. However, when a byte is part of a
// larger number (i.e., a 32-bit field) the most significant bit can indeed
// represent a sign. This test makes sure we convert such a data byte to a
// negative number.
func TestByteToInt64(t *testing.T) {
	b := byte(0xA0)
	val := byteToInt64(b)
	if val != -96 {
		t.Errorf("TestByteToInt64: byteToInt64(0xA0) expected -96, got %d", val)
	}
	b = 0x7F
	val = byteToInt64(b)
	if val != 127 {
		t.Errorf("TestByteToInt64: byteToInt64(0x7F) expected 127, got %d", val)
	}
}

func TestConvertInterfaceToJavaByte(t *testing.T) {
	argUint8 := byte(0xA0)
	jb := convertInterfaceToByte(argUint8)
	if jb != -96 {
		t.Errorf("TestConvertInterfaceToJavaByte: convertInterfaceToByte(0xA0) expected -96, got %d", jb)
	}
	argUint8 = 0x7F
	jb = convertInterfaceToByte(argUint8)
	if jb != 127 {
		t.Errorf("TestConvertInterfaceToJavaByte: convertInterfaceToByte(0x7F) expected 127, got %d", jb)
	}
	argInt := int(0xA0)
	jb = convertInterfaceToByte(argInt)
	if jb != -96 {
		t.Errorf("TestConvertInterfaceToJavaByte: convertInterfaceToByte(0xA0) expected -96, got %d", jb)
	}
	argInt = 0x7F
	jb = convertInterfaceToByte(argInt)
	if jb != 127 {
		t.Errorf("TestConvertInterfaceToJavaByte: convertInterfaceToByte(0x7F) expected 127, got %d", jb)
	}
	argJavaByte := types.JavaByte(-1)
	jb = convertInterfaceToByte(argJavaByte)
	if jb != -1 {
		t.Errorf("TestConvertInterfaceToJavaByte: convertInterfaceToByte(-1) expected -1, got %d", jb)
	}
	argJavaByte = 0x7F
	jb = convertInterfaceToByte(argJavaByte)
	if jb != 127 {
		t.Errorf("TestConvertInterfaceToJavaByte: convertInterfaceToByte(0x7F) expected 127, got %d", jb)
	}
	argInt64 := int64(32767)
	jb = convertInterfaceToByte(argInt64)
	if jb != -1 {
		t.Errorf("TestConvertInterfaceToJavaByte: convertInterfaceToByte(32767) expected -1, got %d", jb)
	}
	argInt64 = 0x7F
	jb = convertInterfaceToByte(argInt64)
	if jb != 127 {
		t.Errorf("TestConvertInterfaceToJavaByte: convertInterfaceToByte(0x7F) expected 127, got %d", jb)
	}
	argRubbish := "ABC"
	jb = convertInterfaceToByte(argRubbish)
	if jb != 0 {
		t.Errorf("TestConvertInterfaceToJavaByte: convertInterfaceToByte(\"ABC\") expected 0, got %d", jb)
	}
}

func TestConvertInterfaceToUint64(t *testing.T) {
	var i64 int64 = 200
	var f64 float64 = 345.0
	var i64ptr = unsafe.Pointer(&i64)

	ret := convertInterfaceToUint64(i64)
	if ret != 200 {
		t.Errorf("TestConvertInterfaceToUint64: Expected convertInterfaceToUint64(200) to return 200, got %d\n", ret)
	}

	ret = convertInterfaceToUint64(f64)
	if ret != 345 {
		t.Errorf("TestConvertInterfaceToUint64: Expected convertInterfaceToUint64(345.0) to return 345, got %d\n", ret)
	}

	ret = convertInterfaceToUint64(i64ptr)
	if ret == 200 {
		t.Errorf("TestConvertInterfaceToUint64: Expected convertInterfaceToUint64(ptr to 200) to return 200, got %d\n", ret)
	}

	ret = convertInterfaceToUint64("ABC")
	if ret != 0 {
		t.Errorf("TestConvertInterfaceToUint64: convertInterfaceToUint64(\"ABC\") expected 0, got %d", ret)
	}

}

func TestIfClassAisAsubclassOfBool(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// Initialize classloaders and method area
	err := classloader.Init()
	if err != nil {
		t.Errorf("Failure to load classes in TestInvokeSpecialJavaLangObject")
	}
	classloader.LoadBaseClasses() // must follow classloader.Init()
	classAname := "java/lang/ClassNotFoundException"
	classA := stringPool.GetStringIndex(&classAname)

	classBname := "java/lang/Throwable"
	classB := stringPool.GetStringIndex(&classBname)

	isIt := isClassAaSublclassOfB(classA, classB)
	if !isIt {
		t.Errorf("%s is a subclass of %s, but result said not",
			classAname, classBname)
	}
}

// check that a class is not a subclass of itself
func TestIfClassAisAsubclassOfItself(t *testing.T) {

}
func TestIfClassAisAsubclassOfBoolInvalid(t *testing.T) {
	globals.InitGlobals("test")

	isIt := isClassAaSublclassOfB(127, 127)
	if !isIt {
		t.Errorf("Expecting identical classes to return true, but returned false")
	}
}

// check that if an array is cast to an object, only java/lang/Object works.
func TestCheckCastArray1(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// Initialize classloaders and method area
	err := classloader.Init()
	if err != nil {
		t.Errorf("Failure to load classes in TestCheckCastArray1")
	}
	classloader.LoadBaseClasses()

	array := object.Make1DimArray(object.INT, 10)

	ret := checkcastArray(array, "java/lang/Object")
	if !ret {
		t.Errorf("checkcastArray(array, \"java/lang/Object\") shoud return true, got false")
	}

	ret = checkcastArray(array, "java/lang/Array")
	if ret {
		t.Errorf("checkcastArray(array, \"java/lang/Object\") shoud return false, got true")
	}
}

// check that two identical arrays come back as castable to each other
func TestCheckCastArray2(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// Initialize classloaders and method area
	err := classloader.Init()
	if err != nil {
		t.Errorf("Failure to load classes in TestCheckCastArray2")
	}
	classloader.LoadBaseClasses()

	array1 := object.Make1DimArray(object.INT, 10)
	array2 := object.Make1DimArray(object.INT, 10)

	ret := checkcastArray(array1, *(stringPool.GetStringPointer(array2.KlassName)))
	if !ret {
		t.Errorf("checkcastArray of two identical arrays should return true, got false")
	}
}

// check that two reference arrays are castable if one is a subclass of the other
func TestCheckCastArray3(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// Initialize classloaders and method area
	err := classloader.Init()
	if err != nil {
		t.Errorf("Failure to load classes in TestCheckCastArray3")
	}
	classloader.LoadBaseClasses()

	object := object.MakeEmptyObject()
	objectKlassName := "[Ljava/lang/NullPointerException;"
	object.KlassName = stringPool.GetStringIndex(&objectKlassName)

	ret := checkcastArray(object, "[java/lang/Throwable")
	if !ret {
		t.Errorf("checkcastArray of a subclass array should return true, got false")
	}
}

func TestPushPeekPop(t *testing.T) {
	// Let verbose trace messages go into a pipe that we will never see.
	// We only care about evaluating the return from pop() in the loop.
	// On Windows, this call must appear before anything else.
	// testutil.UTnewConsole(t)

	testutil.UTinit(t)
	globals.TraceVerbose = false
	var ret, thing interface{}
	flagDeepTracing := false

	// Create frame (fr).
	fr := frames.CreateFrame(13)
	fr.Thread = 0 // Mainthread
	// left nil: fr.FrameStack
	fr.ClName = "TestClass"
	fr.MethName = "TestMethod"
	fr.MethType = "()V"
	fr.Meth = []byte{byte(opcodes.NOP)}

	// Try a nil stack.
	t.Log("Trying a pop() with a nil stack. Ignore the following ThrowEx stack underflow warning.")
	ret = pop(fr)
	if ret != nil {
		t.Errorf("TestPushPeekPop(nil): fr.TOS = -1. Expected nil returned from pop(), got %v", ret)
	}

	// Setup loop.
	globals.TraceVerbose = true
	objstr42 := object.StringObjectFromGoString("42")
	barray := []byte{'A', 'B', 'C'}
	jba := []types.JavaByte{'A', 'B', 'C'}
	rubbish := "rubbish"

	// Push 8 / pop 8 in a loop.
	for ix := 0; ix < 3; ix++ {
		push(fr, nil)
		push(fr, int64(42))
		push(fr, float64(42.0))
		push(fr, objstr42)
		push(fr, &barray)
		push(fr, barray)
		push(fr, jba)
		push(fr, rubbish)

		// rubbish
		thing = peek(fr)
		ret = pop(fr)
		if ret != thing {
			t.Errorf("TestPushPeekPop(rubbish): Loop %d. ret != thing. thing=%v, ret=%v", ix, thing, ret)
		}
		if ret != rubbish {
			t.Errorf("TestPushPeekPop(rubbish): Loop %d. Expected \"rubbish\" returned from pop(), got %v", ix, ret)
			break
		}
		if flagDeepTracing {
			t.Log("#7 rubbish ok")
		}

		// jba
		thing = peek(fr)
		thingString := object.GoStringFromJavaByteArray(thing.([]types.JavaByte))
		ret = pop(fr)
		retString := object.GoStringFromJavaByteArray(ret.([]types.JavaByte))
		if retString != thingString {
			t.Errorf("TestPushPeekPop(jba): Loop %d. ret != thing. thing=%v, ret=%v", ix, thing, ret)
		}
		if retString != "ABC" {
			t.Errorf("TestPushPeekPop(jba): Loop %d. Expected \"ABC\" in a []types.JavaByte returned from pop(), got %v", ix, ret)
			break
		}
		if flagDeepTracing {
			t.Log("#6 jba ok")
		}

		// barray
		thing = peek(fr)
		thingString = string(thing.([]uint8))
		ret = pop(fr)
		retString = string(ret.([]uint8))
		if retString != thingString {
			t.Errorf("TestPushPeekPop(barray): Loop %d. ret != thing. thing=%v, ret=%v", ix, thing, ret)
		}
		if retString != "ABC" {
			t.Errorf("TestPushPeekPop(barray): Loop %d. Expected \"ABC\" in a []uint8 returned from pop(), got %v", ix, ret)
			break
		}
		if flagDeepTracing {
			t.Log("#5 barray ok")
		}

		// &barray
		thing = peek(fr)
		ret = pop(fr)
		if ret != thing {
			t.Errorf("TestPushPeekPop(&barray): Loop %d. ret != thing. thing=%v, ret=%v", ix, thing, ret)
		}
		retString = string(*ret.(*[]byte))
		if retString != "ABC" {
			t.Errorf("TestPushPeekPop(&barray): Loop %d. Expected \"ABC\" in a *[]uint8 returned from pop(), got %v", ix, ret)
			break
		}
		if flagDeepTracing {
			t.Log("#4 &barray ok")
		}

		// objstr42
		thing = peek(fr)
		ret = pop(fr)
		if ret != thing {
			t.Errorf("TestPushPeekPop(objstr42): Loop %d. ret != thing. thing=%v, ret=%v", ix, thing, ret)
		}
		retString = object.GoStringFromStringObject(ret.(*object.Object))
		if retString != "42" {
			t.Errorf("TestPushPeekPop(objstr42): Loop %d. Expected \"42\" in a String object returned from pop(), got %v", ix, ret)
			break
		}
		if flagDeepTracing {
			t.Log("#3 objstr42 ok")
		}

		// float64
		thing = peek(fr)
		ret = pop(fr)
		if ret != thing {
			t.Errorf("TestPushPeekPop(float64): Loop %d. ret != thing. thing=%v, ret=%v", ix, thing, ret)
		}
		if ret.(float64) != 42.0 {
			t.Errorf("TestPushPeekPop(float64): Loop %d. Expected 42.0 in a float64 from pop(), got %v", ix, ret)
			break
		}
		if flagDeepTracing {
			t.Log("#2 float64 ok")
		}

		// int64
		thing = peek(fr)
		ret = pop(fr)
		if ret != thing {
			t.Errorf("TestPushPeekPop(int64): Loop %d. ret != thing. thing=%v, ret=%v", ix, thing, ret)
		}
		if ret.(int64) != 42 {
			t.Errorf("TestPushPeekPop(int64): Loop %d. Expected 42 in an int64 from pop(), got %v", ix, ret)
			break
		}
		if flagDeepTracing {
			t.Log("#1 int64 ok")
		}

		// nil
		thing = peek(fr)
		ret = pop(fr)
		if ret != thing {
			t.Errorf("TestPushPeekPop(nil): Loop %d. ret != thing. thing=%v, ret=%v", ix, thing, ret)
		}
		if ret != nil {
			t.Errorf("TestPushPeekPop(nil): Loop %d. Expected nil returned from pop(), got %v", ix, ret)
			break
		}
		if flagDeepTracing {
			t.Log("#0 nil ok")
		}

	}

	// Restore console for go test.
	// testutil.UTrestoreConsole(t)

}

func TestEmitTraceData(t *testing.T) {
	// Let verbose trace messages go into a pipe that we will never see.
	// We only care about evaluating the return from pop() in the loop.
	// On Windows, this must be called before anything else.
	// testutil.UTnewConsole(t)

	testutil.UTinit(t)
	globals.TraceVerbose = true
	var ret interface{}
	flagDeepTracing := false

	// Create frame (fr).
	fr := frames.CreateFrame(13)
	fr.Thread = 0 // Mainthread
	// left nil: fr.FrameStack
	fr.ClName = "TestClass"
	fr.MethName = "TestMethod"
	fr.MethType = "()V"
	fr.Meth = []byte{byte(opcodes.NOP)}

	// Try a nil stack.
	ret = EmitTraceData(fr)
	if !strings.Contains(ret.(string), "TOS:  -") {
		t.Errorf("TestEmitTraceData(nil): Expected \"TOS: -\", got: %v", ret)
	}

	// Setup loop.
	objstr42 := object.StringObjectFromGoString("42")
	barray := []byte{'A', 'B', 'C'}
	jba := []types.JavaByte{'A', 'B', 'C'}
	rubbish := "rubbish"

	// EmitTraceData in a loop for 8 different top of stack variables.
	for ix := 0; ix < 3; ix++ {

		start := time.Now()

		// rubbish
		push(fr, object.StringObjectFromGoString(rubbish))
		ret = EmitTraceData(fr)
		_ = pop(fr)
		if !strings.Contains(ret.(string), "String: rubbish") {
			t.Errorf("TestEmitTraceData(rubbish): Loop %d. Expected \"String: rubbish\", got: %v", ix, ret)
			break
		}

		if flagDeepTracing {
			t.Log("#7 rubbish ok")
		}

		// jba
		push(fr, jba)
		ret = EmitTraceData(fr)
		_ = pop(fr)
		if !strings.Contains(ret.(string), "[]JavaByte: ABC") {
			t.Errorf("TestEmitTraceData(jba): Loop %d. Expected \"[]JavaByte: ABC\", got: %v", ix, ret)
			break
		}
		if flagDeepTracing {
			t.Log("#6 jba ok")
		}

		// barray
		push(fr, barray)
		ret = EmitTraceData(fr)
		_ = pop(fr)
		if !strings.Contains(ret.(string), "[]byte: ABC") {
			t.Errorf("TestEmitTraceData(barray): Loop %d. Expected \"[]byte: ABC\", got %v", ix, ret)
			break
		}
		if flagDeepTracing {
			t.Log("#5 barray ok")
		}

		// &barray
		push(fr, &barray)
		ret = EmitTraceData(fr)
		_ = pop(fr)
		if !strings.Contains(ret.(string), "[]byte: ABC") {
			t.Errorf("TestEmitTraceData(&barray): Loop %d. Expected \"[]byte: ABC\", got: %v", ix, ret)
			break
		}
		if flagDeepTracing {
			t.Log("#4 &barray ok")
		}

		// objstr42
		push(fr, objstr42)
		ret = EmitTraceData(fr)
		_ = pop(fr)
		if !strings.Contains(ret.(string), "String: 42") {
			t.Errorf("TestEmitTraceData(objstr42): Loop %d. Expected \"String: 42\", got: %v", ix, ret)
			break
		}

		if flagDeepTracing {
			t.Log("#3 objstr42 ok")
		}

		// float64
		push(fr, 42.0)
		ret = EmitTraceData(fr)
		_ = pop(fr)
		if !strings.Contains(ret.(string), "42") {
			t.Errorf("TestEmitTraceData(float64): Loop %d. Expected \"42\", got: %v", ix, ret)
			break
		}
		if flagDeepTracing {
			t.Log("#2 float64 ok")
		}

		// int64
		push(fr, int64(42))
		ret = EmitTraceData(fr)
		_ = pop(fr)
		if !strings.Contains(ret.(string), "42") {
			t.Errorf("TestEmitTraceData(int64): Loop %d. Expected \"42\", got: %v", ix, ret)
			break
		}
		if flagDeepTracing {
			t.Log("#1 int64 ok")
		}

		// nil
		push(fr, object.Null)
		ret = EmitTraceData(fr)
		_ = pop(fr)
		if !strings.Contains(ret.(string), "<null>") {
			t.Errorf("TestEmitTraceData(nil): Loop %d. Expected \"<null>\", got: %v", ix, ret)
			break
		}
		if flagDeepTracing {
			t.Log("#0 nil ok")
		}

		elapsed := time.Since(start)
		t.Logf("Loop %d consumed %s", ix, elapsed)

	}

	// Restore console for go test.
	// testutil.UTrestoreConsole(t)
}
