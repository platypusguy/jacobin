/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"reflect"
	"testing"

	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
)

func TestLoad_Lang_Character_RegistersMethods(t *testing.T) {
	saved := MethodSignatures
	defer func() { MethodSignatures = saved }()
	MethodSignatures = make(map[string]GMeth)

	Load_Lang_Character()

	checks := []struct {
		key   string
		slots int
		fn    func([]interface{}) interface{}
	}{
		{"java/lang/Character.<clinit>()V", 0, clinitGeneric},
		{"java/lang/Character.compare(CC)I", 2, charCompare},
		{"java/lang/Character.digit(CI)I", 2, charDigit},
		{"java/lang/Character.equals(Ljava/lang/Object;)Z", 1, charEquals},
		{"java/lang/Character.forDigit(II)C", 2, charForDigit},
		{"java/lang/Character.hashCode()I", 0, charHashCode},
		{"java/lang/Character.hashCode(C)I", 1, charHashCodeStatic},
		{"java/lang/Character.isDigit(C)Z", 1, charIsDigit},
		{"java/lang/Character.isLetter(C)Z", 1, charIsLetter},
		{"java/lang/Character.isLowerCase(C)Z", 1, charIsLowerCase},
		{"java/lang/Character.isUpperCase(C)Z", 1, charIsUpperCase},
		{"java/lang/Character.isWhitespace(C)Z", 1, charIsWhitespace},
		{"java/lang/Character.toString()Ljava/lang/String;", 0, charToString},
		{"java/lang/Character.toString(C)Ljava/lang/String;", 1, charToStringStatic},
		{"java/lang/Character.charValue()C", 0, charValue},
		{"java/lang/Character.toLowerCase(C)C", 1, charToLowerCase},
		{"java/lang/Character.toUpperCase(C)C", 1, charToUpperCase},
		{"java/lang/Character.valueOf(C)Ljava/lang/Character;", 1, characterValueOf},
	}

	for _, c := range checks {
		got, ok := MethodSignatures[c.key]
		if !ok {
			t.Fatalf("missing MethodSignatures entry for %s", c.key)
		}
		if got.ParamSlots != c.slots {
			t.Fatalf("%s ParamSlots expected %d, got %d", c.key, c.slots, got.ParamSlots)
		}
		if got.GFunction == nil {
			t.Fatalf("%s GFunction expected non-nil", c.key)
		}
		if reflect.ValueOf(got.GFunction).Pointer() != reflect.ValueOf(c.fn).Pointer() {
			t.Fatalf("%s GFunction mismatch", c.key)
		}
	}
}

func TestCharacter_IsDigit_IsLetter(t *testing.T) {
	globals.InitGlobals("test")

	if v := charIsDigit([]interface{}{int64('0')}).(int64); v != types.JavaBoolTrue {
		t.Fatalf("isDigit('0') expected true")
	}
	if v := charIsDigit([]interface{}{int64('A')}).(int64); v != types.JavaBoolFalse {
		t.Fatalf("isDigit('A') expected false")
	}

	if v := charIsLetter([]interface{}{int64('A')}).(int64); v != types.JavaBoolTrue {
		t.Fatalf("isLetter('A') expected true")
	}
	if v := charIsLetter([]interface{}{int64('1')}).(int64); v != types.JavaBoolFalse {
		t.Fatalf("isLetter('1') expected false")
	}
}

func TestCharacter_ToLower_ToUpper_ValueOf_CharValue(t *testing.T) {
	globals.InitGlobals("test")

	if v := charToLowerCase([]interface{}{int64('Z')}).(int64); v != int64('z') {
		t.Fatalf("toLowerCase('Z') expected 'z'")
	}
	if v := charToUpperCase([]interface{}{int64('a')}).(int64); v != int64('A') {
		t.Fatalf("toUpperCase('a') expected 'A'")
	}

	obj := characterValueOf([]interface{}{int64('Q')}).(*object.Object)
	if vv := obj.FieldTable["value"].Fvalue.(int64); vv != int64('Q') {
		t.Fatalf("valueOf('Q') wrong: %v", vv)
	}

	if cv := charValue([]interface{}{obj}).(int64); cv != int64('Q') {
		t.Fatalf("charValue expected 'Q', got %v", cv)
	}
}

func TestCharacter_AdditionalMethods(t *testing.T) {
	globals.InitGlobals("test")

	// compare
	if res := charCompare([]interface{}{int64('a'), int64('b')}).(int64); res >= 0 {
		t.Errorf("compare('a', 'b') expected < 0, got %d", res)
	}
	if res := charCompare([]interface{}{int64('b'), int64('a')}).(int64); res <= 0 {
		t.Errorf("compare('b', 'a') expected > 0, got %d", res)
	}
	if res := charCompare([]interface{}{int64('a'), int64('a')}).(int64); res != 0 {
		t.Errorf("compare('a', 'a') expected 0, got %d", res)
	}

	// equals
	c1 := characterValueOf([]interface{}{int64('x')}).(*object.Object)
	c2 := characterValueOf([]interface{}{int64('x')}).(*object.Object)
	c3 := characterValueOf([]interface{}{int64('y')}).(*object.Object)
	if res := charEquals([]interface{}{c1, c2}); res != types.JavaBoolTrue {
		t.Errorf("equals('x', 'x') expected true, got %v", res)
	}
	if res := charEquals([]interface{}{c1, c3}); res != types.JavaBoolFalse {
		t.Errorf("equals('x', 'y') expected false, got %v", res)
	}

	// hashCode
	if res := charHashCodeStatic([]interface{}{int64('A')}).(int64); res != int64('A') {
		t.Errorf("hashCodeStatic('A') expected %d, got %d", int64('A'), res)
	}
	if res := charHashCode([]interface{}{c1}).(int64); res != int64('x') {
		t.Errorf("hashCode('x') expected %d, got %d", int64('x'), res)
	}

	// isLowerCase / isUpperCase
	if res := charIsLowerCase([]interface{}{int64('a')}); res != types.JavaBoolTrue {
		t.Errorf("isLowerCase('a') expected true")
	}
	if res := charIsLowerCase([]interface{}{int64('A')}); res != types.JavaBoolFalse {
		t.Errorf("isLowerCase('A') expected false")
	}
	if res := charIsUpperCase([]interface{}{int64('A')}); res != types.JavaBoolTrue {
		t.Errorf("isUpperCase('A') expected true")
	}
	if res := charIsUpperCase([]interface{}{int64('a')}); res != types.JavaBoolFalse {
		t.Errorf("isUpperCase('a') expected false")
	}

	// isWhitespace
	if res := charIsWhitespace([]interface{}{int64(' ')}); res != types.JavaBoolTrue {
		t.Errorf("isWhitespace(' ') expected true")
	}
	if res := charIsWhitespace([]interface{}{int64('\t')}); res != types.JavaBoolTrue {
		t.Errorf("isWhitespace('\\t') expected true")
	}
	if res := charIsWhitespace([]interface{}{int64('a')}); res != types.JavaBoolFalse {
		t.Errorf("isWhitespace('a') expected false")
	}

	// digit
	if res := charDigit([]interface{}{int64('5'), int64(10)}).(int64); res != 5 {
		t.Errorf("digit('5', 10) expected 5, got %d", res)
	}
	if res := charDigit([]interface{}{int64('a'), int64(16)}).(int64); res != 10 {
		t.Errorf("digit('a', 16) expected 10, got %d", res)
	}
	if res := charDigit([]interface{}{int64('f'), int64(16)}).(int64); res != 15 {
		t.Errorf("digit('f', 16) expected 15, got %d", res)
	}
	if res := charDigit([]interface{}{int64('G'), int64(16)}).(int64); res != -1 {
		t.Errorf("digit('G', 16) expected -1, got %d", res)
	}

	// forDigit
	if res := charForDigit([]interface{}{int64(5), int64(10)}).(int64); res != int64('5') {
		t.Errorf("forDigit(5, 10) expected '5', got %c", rune(res))
	}
	if res := charForDigit([]interface{}{int64(15), int64(16)}).(int64); res != int64('f') {
		t.Errorf("forDigit(15, 16) expected 'f', got %c", rune(res))
	}
	if res := charForDigit([]interface{}{int64(16), int64(16)}).(int64); res != 0 {
		t.Errorf("forDigit(16, 16) expected 0, got %d", res)
	}

	// toString
	sObj := charToStringStatic([]interface{}{int64('A')}).(*object.Object)
	if s := object.GoStringFromStringObject(sObj); s != "A" {
		t.Errorf("toStringStatic('A') expected 'A', got '%s'", s)
	}
	sObj2 := charToString([]interface{}{c1}).(*object.Object)
	if s := object.GoStringFromStringObject(sObj2); s != "x" {
		t.Errorf("charToString('x') expected 'x', got '%s'", s)
	}
}
