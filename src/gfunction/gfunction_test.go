/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"io"
	"jacobin/src/classloader"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"jacobin/src/types"
	"math/big"
	"os"
	"strings"
	"testing"
)

func f1([]interface{}) interface{} { return nil }
func f2([]interface{}) interface{} { return nil }
func f3([]interface{}) interface{} { return nil }

func TestMTableLoadLib(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	libMeths := make(map[string]ghelpers.GMeth)
	libMeths["test.f1()V"] = ghelpers.GMeth{ParamSlots: 0, GFunction: f1}
	libMeths["test.f2(I)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: f2}
	libMeths["test.f3(Ljava/lang/String;JZ)D"] = ghelpers.GMeth{ParamSlots: 3, GFunction: f3}
	mtbl := make(classloader.MT)
	loadlib(&mtbl, libMeths)
	if len(mtbl) != 3 {
		t.Errorf("ERROR, Expecting MTable with 3 entries, got: %d\n", len(mtbl))
	}
	mte := libMeths["test.f1()V"]
	if mte.ParamSlots != 0 {
		t.Errorf("ERROR, Expecting f1 MTable entry to have 0 param slots, got: %d\n",
			mte.ParamSlots)
	}
	mte = libMeths["test.f2(I)V"]
	if mte.ParamSlots != 1 {
		t.Errorf("ERROR, Expecting f2 MTable entry to have 1 param slots, got: %d\n",
			mte.ParamSlots)
	}
	mte = libMeths["test.f3(Ljava/lang/String;JZ)D"]
	if mte.ParamSlots != 3 {
		t.Errorf("ERROR, Expecting f3 MTable entry to have 3 param slots, got: %d\n",
			mte.ParamSlots)
	}

	if mte.NeedsContext {
		t.Errorf("ERROR, Expecting MTable entry's NeedContext to be false\n")
	}
}

// test loading of native functions

func TestMTableLoadGFunctions(t *testing.T) {
	classloader.MTable = make(map[string]classloader.MTentry)
	MTableLoadGFunctions(&classloader.MTable)
	mte, exists := classloader.MTable["java/lang/Object.<init>()V"]
	if !exists {
		t.Errorf("Expecting MTable entry for java/lang/Object.<init>()V, but it does not exist")
	}

	if mte.MType != 'G' {
		t.Errorf("Expecting java/lang/Object.<init>()V to be of type 'G', but got type: %c",
			mte.MType)
	}

	toolkitMethod := "java/awt/Toolkit.loadLibraries()V"
	toolkitEntry, exists := classloader.MTable[toolkitMethod]
	if !exists {
		t.Fatalf("Expecting MTable entry for %s, but it does not exist", toolkitMethod)
	}
	if toolkitEntry.MType != 'G' {
		t.Errorf("Expecting %s to be of type 'G', but got type: %c", toolkitMethod, toolkitEntry.MType)
	}
}

func TestCheckKey(t *testing.T) {
	if checkKey("java/lang/Object") != false {
		t.Errorf("invalid key %s was allowed in gfunction", "java/lang/Object")
	}

	if checkKey("java/lang/Object.toString") != false {
		t.Errorf("invalid key %s was allowed in gfunction",
			"java/lang/Object.toString")
	}

	if checkKey("java/lang/Object.toString(") != false {
		t.Errorf("invalid key %s was allowed in gfunction",
			"java/lang/Object.toString(")
	}

	if checkKey("java/lang/Object.toString()") != false {
		t.Errorf("invalid key %s was allowed in gfunction",
			"java/lang/Object.toString()")
	}

	if checkKey("java/lang/Object.toString(I)Ljava/lang/String;") != true {
		t.Errorf("got unexpected error checking key %s",
			"java/lang/Object.toString(I)Ljava/lang/String;")
	}
}

func TestPopulate_PrimitivesAndString(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()

	// Integer primitive object
	iobj := object.MakePrimitiveObject("java/lang/Integer", "I", int64(42))
	if iobj == nil {
		t.Fatalf("object.MakePrimitiveObject returned nil for Integer")
	}
	fld, ok := iobj.FieldTable["value"]
	if !ok || fld.Ftype != "I" {
		t.Fatalf("Integer value field missing or wrong type: %#v", iobj.FieldTable["value"])
	}
	if v, ok := fld.Fvalue.(int64); !ok || v != 42 {
		t.Fatalf("Integer value mismatch: %v", fld.Fvalue)
	}

	// String via StringIndex path returns a proper String object
	sobj := object.MakePrimitiveObject("java/lang/String", "T", "hello")
	if sobj == nil {
		t.Fatalf("object.MakePrimitiveObject returned nil for String")
	}
	if !object.IsStringObject(sobj) {
		t.Fatalf("object.MakePrimitiveObject did not create a String object")
	}
	if s := object.GoStringFromStringObject(sobj); s != "hello" {
		t.Fatalf("String content mismatch: %q", s)
	}
}

func TestReturnNullTrueFalse(t *testing.T) {
	if v := ghelpers.ReturnNull(nil); v != object.Null {
		t.Fatalf("ghelpers.ReturnNull did not return object.Null: %v", v)
	}
	if v := ghelpers.ReturnTrue(nil); v != types.JavaBoolTrue {
		t.Fatalf("ghelpers.ReturnTrue != true: %v", v)
	}
	if v := ghelpers.ReturnFalse(nil); v != types.JavaBoolFalse {
		t.Fatalf("ghelpers.ReturnFalse != false: %v", v)
	}
}

func TestEOFSetGet(t *testing.T) {
	obj := object.MakeEmptyObject()
	ghelpers.EofSet(obj, true)
	if !ghelpers.EofGet(obj) {
		t.Fatalf("ghelpers.EofGet expected true")
	}
	ghelpers.EofSet(obj, false)
	if ghelpers.EofGet(obj) {
		t.Fatalf("ghelpers.EofGet expected false")
	}
}

func TestReturnRandomLong_Type(t *testing.T) {
	v := ghelpers.ReturnRandomLong(nil)
	if _, ok := v.(int64); !ok {
		t.Fatalf("ghelpers.ReturnRandomLong did not return int64, got %T", v)
	}
}

func TestGetGErrBlk(t *testing.T) {
	errBlk := ghelpers.GetGErrBlk(123, "test error")
	if errBlk.ExceptionType != 123 {
		t.Errorf("Expected ExceptionType 123, got %d", errBlk.ExceptionType)
	}
	if errBlk.ErrMsg != "test error" {
		t.Errorf("Expected ErrMsg 'test error', got '%s'", errBlk.ErrMsg)
	}
}

func TestSimpleReturnFunctions(t *testing.T) {
	if ghelpers.ClinitGeneric(nil) != nil {
		t.Error("ghelpers.ClinitGeneric should return nil")
	}
	if ghelpers.JustReturn(nil) != nil {
		t.Error("ghelpers.JustReturn should return nil")
	}
	if ghelpers.ReturnNull(nil) != object.Null {
		t.Error("ghelpers.ReturnNullObject should return object.Null")
	}
}

func TestReturnCharsetName(t *testing.T) {
	globals.InitGlobals("test")
	res := ghelpers.ReturnCharsetName(nil)
	obj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("ghelpers.ReturnCharsetName should return *object.Object, got %T", res)
	}
	if !object.IsStringObject(obj) {
		t.Error("ghelpers.ReturnCharsetName should return a String object")
	}
	charset := object.GoStringFromStringObject(obj)
	if charset != globals.GetCharsetName() {
		t.Errorf("Expected charset %s, got %s", globals.GetCharsetName(), charset)
	}
}

func TestInitBigIntegerField(t *testing.T) {
	cases := []struct {
		val  int64
		sign int64
	}{
		{100, 1},
		{0, 0},
		{-100, -1},
	}

	for _, c := range cases {
		obj := object.MakeEmptyObject()
		ghelpers.InitBigIntegerField(obj, c.val)

		fldVal, ok := obj.FieldTable["value"]
		if !ok {
			t.Errorf("value field missing for %d", c.val)
			continue
		}
		bigInt, ok := fldVal.Fvalue.(*big.Int)
		if !ok || bigInt.Int64() != c.val {
			t.Errorf("value field mismatch for %d: got %v", c.val, fldVal.Fvalue)
		}

		fldSign, ok := obj.FieldTable["signum"]
		if !ok {
			t.Errorf("signum field missing for %d", c.val)
			continue
		}
		if fldSign.Fvalue.(int64) != c.sign {
			t.Errorf("signum mismatch for %d: expected %d, got %d", c.val, c.sign, fldSign.Fvalue)
		}
	}
}

func TestInvoke(t *testing.T) {
	key := "test/Invoke.run()V"
	ghelpers.MethodSignatures[key] = ghelpers.GMeth{
		ParamSlots: 0,
		GFunction: func(p []interface{}) interface{} {
			return "invoked"
		},
	}
	defer delete(ghelpers.MethodSignatures, key)

	res := ghelpers.Invoke(key, nil)
	if res != "invoked" {
		t.Errorf("Expected 'invoked', got %v", res)
	}

	// Test NoSuchMethodException
	res = ghelpers.Invoke("non/Existent", nil)
	switch res.(type) {
	case *ghelpers.GErrBlk:
	default:
		t.Errorf("Expected *ghelpers.GErrBlk, got %v", res)
	}
}

func TestGetDefaultSecurityProvider(t *testing.T) {
	globals.InitGlobals("test")
	p1 := ghelpers.GetDefaultSecurityProvider()
	if p1 == nil {
		t.Fatal("javaSecurity.GetDefaultSecurityProvider returned nil")
	}
	p2 := ghelpers.GetDefaultSecurityProvider()
	if p1 != p2 {
		t.Error("javaSecurity.GetDefaultSecurityProvider should return a singleton")
	}
}

func TestLoadTestGfunctions(t *testing.T) {
	mt := make(classloader.MT)
	LoadTestGfunctions(&mt)
	if !ghelpers.TestGfunctionsLoaded {
		t.Error("ghelpers.TestGfunctionsLoaded should be true after LoadTestGfunctions")
	}
	// Verify some test function is loaded.
	// Based on testGfunctions.go if available, but we can check if mt is not empty.
	if len(mt) == 0 {
		t.Error("MTable should not be empty after LoadTestGfunctions")
	}
}

func TestLoadlib_InvalidKey(t *testing.T) {
	globals.InitGlobals("test")
	// to inspect log messages, redirect stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	libMeths := make(map[string]ghelpers.GMeth)
	libMeths["invalidKey"] = ghelpers.GMeth{ParamSlots: 0, GFunction: f1}
	mtbl := make(classloader.MT)
	loadlib(&mtbl, libMeths)

	// restore stderr to what it was before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])
	if !strings.Contains(msg, "loadlib: at least one key was invalid") {
		t.Errorf("Expected error message about invalid key, got: %s", msg)
	}
}

func TestLoadlib_TracingEnabled(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()
	trace.Init()
	globals.TraceInst = true
	defer func() {
		globals.TraceInst = false
	}()

	libMeths := make(map[string]ghelpers.GMeth)
	key := "java/lang/Math.max(II)I"
	libMeths[key] = ghelpers.GMeth{ParamSlots: 2, GFunction: f1}

	mtbl := make(classloader.MT)
	loadlib(&mtbl, libMeths)

	mte, exists := mtbl[key]
	if !exists {
		t.Fatalf("expected key %q to exist in mtbl", key)
	}
	if mte.MType != 'G' {
		t.Errorf("expected MType 'G', got %c", mte.MType)
	}

	className := *(stringPool.GetStringPointer(mte.MethClass))
	methName := *(stringPool.GetStringPointer(mte.MethName))
	methType := *(stringPool.GetStringPointer(mte.MethType))

	if className != "java/lang/Math" {
		t.Errorf("expected className 'java/lang/Math', got %q", className)
	}
	if methName != "max" {
		t.Errorf("expected methName 'max', got %q", methName)
	}
	if methType != "(II)I" {
		t.Errorf("expected methType '(II)I', got %q", methType)
	}
}

func TestCheckKey_EdgeCases(t *testing.T) {
	invalidKeys := []string{
		"",
		"nodotorparens",
		"java/lang/Object",
		"java/lang/Object.toString",
		"java/lang/Object.toString(",
		"java/lang/Object.toString)",
		"java/lang/Object.toString()",
		"java/lang/Math.abs(I)",
		"java/lang/Math.max(II)",
	}

	for _, k := range invalidKeys {
		if checkKey(k) {
			t.Errorf("checkKey(%q) = true, want false", k)
		}
	}

	validKeys := []string{
		"java/lang/Object.<init>()V",
		"java/lang/String.valueOf(I)Ljava/lang/String;",
		"java/lang/Math.max(DD)D",
		"java/time/ZoneOffset.ofTotalSeconds(I)Ljava/time/ZoneOffset;",
		"java/time/ZoneId.of(Ljava/lang/String;)Ljava/time/ZoneId;",
		"java/time/Instant.now()Ljava/time/Instant;",
		"java/time/ZonedDateTime.now()Ljava/time/ZonedDateTime;",
		"java/time/Clock.systemUTC()Ljava/time/Clock;",
		"java/time/ZoneRegion.ofId(Ljava/lang/String;Z)Ljava/time/ZoneRegion;",
		"java/time/zone/ZoneRulesProvider.getRules(Ljava/lang/String;Z)Ljava/time/zone/ZoneRules;",
		"java/time/zone/ZoneOffsetTransition.of(Ljava/time/LocalDateTime;Ljava/time/ZoneOffset;Ljava/time/ZoneOffset;)Ljava/time/zone/ZoneOffsetTransition;",
		"java/time/OffsetDateTime.now()Ljava/time/OffsetDateTime;",
		"javax/crypto/Cipher.getInstance(Ljava/lang/String;)Ljavax/crypto/Cipher;",
	}

	for _, k := range validKeys {
		if !checkKey(k) {
			t.Errorf("checkKey(%q) = false, want true", k)
		}
	}
}

func TestMTableLoadGFunctions_Comprehensive(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()
	trace.Init()

	classloader.MTable = make(map[string]classloader.MTentry)
	MTableLoadGFunctions(&classloader.MTable)

	expectedKeys := []string{
		"java/lang/Object.<init>()V",
		"java/lang/System.currentTimeMillis()J",
		"java/lang/Math.abs(I)I",
		"java/io/File.exists()Z",
		"java/util/HashMap.<init>()V",
		"java/time/ZoneOffset.ofTotalSeconds(I)Ljava/time/ZoneOffset;",
		"java/time/ZoneId.of(Ljava/lang/String;)Ljava/time/ZoneId;",
		"java/time/Instant.now()Ljava/time/Instant;",
		"java/time/ZonedDateTime.now()Ljava/time/ZonedDateTime;",
		"java/time/Clock.systemUTC()Ljava/time/Clock;",
		"java/time/ZoneRegion.ofId(Ljava/lang/String;Z)Ljava/time/ZoneRegion;",
		"java/time/zone/ZoneRulesProvider.getRules(Ljava/lang/String;Z)Ljava/time/zone/ZoneRules;",
		"java/time/zone/ZoneOffsetTransition.of(Ljava/time/LocalDateTime;Ljava/time/ZoneOffset;Ljava/time/ZoneOffset;)Ljava/time/zone/ZoneOffsetTransition;",
		"java/time/OffsetDateTime.now()Ljava/time/OffsetDateTime;",
		"javax/crypto/Cipher.getInstance(Ljava/lang/String;)Ljavax/crypto/Cipher;",
		"jdk/internal/misc/Unsafe.getUnsafe()Ljdk/internal/misc/Unsafe;",
		"sun/misc/Unsafe.getUnsafe()Lsun/misc/Unsafe;",
	}

	for _, key := range expectedKeys {
		mte, exists := classloader.MTable[key]
		if !exists {
			t.Errorf("expected MTable entry for %q, but not found", key)
			continue
		}
		if mte.MType != 'G' {
			t.Errorf("expected entry %q to have MType 'G', got %c", key, mte.MType)
		}
		gmeth, ok := mte.Meth.(ghelpers.GMeth)
		if !ok {
			t.Errorf("expected entry %q to have Meth of type ghelpers.GMeth, got %T", key, mte.Meth)
		} else if gmeth.GFunction == nil {
			t.Errorf("expected entry %q to have non-nil GFunction", key)
		}
	}
}
