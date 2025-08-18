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

    checks := []struct{
        key   string
        slots int
        fn    func([]interface{}) interface{}
    }{
        {"java/lang/Character.<clinit>()V", 0, clinitGeneric},
        {"java/lang/Character.isDigit(C)Z", 1, charIsDigit},
        {"java/lang/Character.isLetter(C)Z", 1, charIsLetter},
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
