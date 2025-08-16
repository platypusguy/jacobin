/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
    "jacobin/globals"
    "jacobin/object"
    "jacobin/util"
    "testing"
)

func TestNewField_Basics_Getters(t *testing.T) {
    globals.InitGlobals("test")

    cls := object.MakeEmptyObject()
    f := NewField(cls, "count")

    if f == nil {
        t.Fatalf("NewField returned nil")
    }
    if f.Class != cls {
        t.Fatalf("NewField Class mismatch")
    }
    if f.Name != "count" {
        t.Fatalf("NewField Name mismatch: %q", f.Name)
    }

    // default zero values
    if f.Modifiers != 0 {
        t.Fatalf("expected default Modifiers=0, got %d", f.Modifiers)
    }
    if f.Type != nil {
        t.Fatalf("expected default Type=nil")
    }

    // Set Modifiers and Type and validate getters
    f.Modifiers = 0x0010 // arbitrary example (final) value constant-like
    typ := object.MakeEmptyObject()
    f.Type = typ

    if got := f.GetDeclaringClass(); got != cls {
        t.Fatalf("GetDeclaringClass mismatch")
    }
    if got := f.GetName(); got != "count" {
        t.Fatalf("GetName mismatch: %q", got)
    }
    if got := f.GetModifiers(); got != 0x0010 {
        t.Fatalf("GetModifiers mismatch: %d", got)
    }
    if got := f.GetType(); got != typ {
        t.Fatalf("GetType mismatch")
    }
}

func TestField_Equals(t *testing.T) {
    globals.InitGlobals("test")

    cls := object.MakeEmptyObject()
    typ := object.MakeEmptyObject()

    f1 := NewField(cls, "name")
    f1.Type = typ

    // Same underlying references
    f2 := NewField(cls, "name")
    f2.Type = typ

    if !f1.Equals(f2) {
        t.Fatalf("Fields with same Class/Name/Type should be equal")
    }

    // Different name
    f3 := NewField(cls, "other")
    f3.Type = typ
    if f1.Equals(f3) {
        t.Fatalf("Fields with different Name should not be equal")
    }

    // Different class
    cls2 := object.MakeEmptyObject()
    f4 := NewField(cls2, "name")
    f4.Type = typ
    if f1.Equals(f4) {
        t.Fatalf("Fields with different Class should not be equal")
    }

    // Different type
    typ2 := object.MakeEmptyObject()
    f5 := NewField(cls, "name")
    f5.Type = typ2
    if f1.Equals(f5) {
        t.Fatalf("Fields with different Type should not be equal")
    }
}

func TestField_HashCode_DelegatesToClassHash(t *testing.T) {
    globals.InitGlobals("test")

    cls := object.MakeEmptyObject()
    f := NewField(cls, "value")

    want, _ := util.HashAnything(cls)
    if got := f.HashCode(); got != want {
        t.Fatalf("HashCode mismatch: got %d, want %d", got, want)
    }

    // stability across calls
    if got2 := f.HashCode(); got2 != want {
        t.Fatalf("HashCode not stable across calls: got %d, want %d", got2, want)
    }
}
