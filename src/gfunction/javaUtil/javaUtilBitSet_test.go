/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"container/list"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func setupTestGlobals() *globals.Globals {
	globals.InitGlobals("test")
	g := globals.GetGlobalRef()
	g.FuncInstantiateClass = func(classname string, frameStack *list.List) (any, error) {
		return object.MakeEmptyObjectWithClassName(&classname), nil
	}
	return g
}

func TestBitSet_Basic(t *testing.T) {
	setupTestGlobals()

	self := object.MakeEmptyObjectWithClassName(new("java/util/BitSet"))
	bitsetInit([]interface{}{self})

	// Test isEmpty
	if bitsetIsEmpty([]interface{}{self}) != types.JavaBoolTrue {
		t.Errorf("Expected empty BitSet")
	}

	// Test set and get
	bitsetSet([]interface{}{self, int64(10)})
	if bitsetGet([]interface{}{self, int64(10)}) != types.JavaBoolTrue {
		t.Errorf("Expected bit 10 to be set")
	}
	if bitsetGet([]interface{}{self, int64(11)}) != types.JavaBoolFalse {
		t.Errorf("Expected bit 11 to be false")
	}

	// Test cardinality
	if bitsetCardinality([]interface{}{self}) != int64(1) {
		t.Errorf("Expected cardinality 1, got %d", bitsetCardinality([]interface{}{self}))
	}

	// Test clear
	bitsetClear([]interface{}{self, int64(10)})
	if bitsetGet([]interface{}{self, int64(10)}) != types.JavaBoolFalse {
		t.Errorf("Expected bit 10 to be cleared")
	}
	if bitsetIsEmpty([]interface{}{self}) != types.JavaBoolTrue {
		t.Errorf("Expected empty BitSet after clear")
	}
}

func TestBitSet_Range(t *testing.T) {
	setupTestGlobals()
	self := object.MakeEmptyObjectWithClassName(new("java/util/BitSet"))
	bitsetInit([]interface{}{self})

	bitsetSetRange([]interface{}{self, int64(10), int64(20)})
	for i := int64(10); i < 20; i++ {
		if bitsetGet([]interface{}{self, i}) != types.JavaBoolTrue {
			t.Errorf("Expected bit %d to be set", i)
		}
	}
	if bitsetGet([]interface{}{self, int64(9)}) != types.JavaBoolFalse {
		t.Errorf("Expected bit 9 to be false")
	}
	if bitsetGet([]interface{}{self, int64(20)}) != types.JavaBoolFalse {
		t.Errorf("Expected bit 20 to be false")
	}

	if bitsetCardinality([]interface{}{self}) != int64(10) {
		t.Errorf("Expected cardinality 10, got %d", bitsetCardinality([]interface{}{self}))
	}

	bitsetClearRange([]interface{}{self, int64(15), int64(25)})
	for i := int64(10); i < 15; i++ {
		if bitsetGet([]interface{}{self, i}) != types.JavaBoolTrue {
			t.Errorf("Expected bit %d to remain set", i)
		}
	}
	for i := int64(15); i < 25; i++ {
		if bitsetGet([]interface{}{self, i}) != types.JavaBoolFalse {
			t.Errorf("Expected bit %d to be cleared", i)
		}
	}
}

func TestBitSet_Bitwise(t *testing.T) {
	setupTestGlobals()
	bs1 := object.MakeEmptyObjectWithClassName(new("java/util/BitSet"))
	bitsetInit([]interface{}{bs1})
	bitsetSetRange([]interface{}{bs1, int64(0), int64(10)})

	bs2 := object.MakeEmptyObjectWithClassName(new("java/util/BitSet"))
	bitsetInit([]interface{}{bs2})
	bitsetSetRange([]interface{}{bs2, int64(5), int64(15)})

	// Test Intersects
	if bitsetIntersects([]interface{}{bs1, bs2}) != types.JavaBoolTrue {
		t.Errorf("Expected intersection")
	}

	// Test And
	bsAnd := bitsetClone([]interface{}{bs1}).(*object.Object)
	bitsetAnd([]interface{}{bsAnd, bs2})
	if bitsetCardinality([]interface{}{bsAnd}) != int64(5) {
		t.Errorf("Expected cardinality 5 after AND, got %d", bitsetCardinality([]interface{}{bsAnd}))
	}

	// Test Or
	bsOr := bitsetClone([]interface{}{bs1}).(*object.Object)
	bitsetOr([]interface{}{bsOr, bs2})
	if bitsetCardinality([]interface{}{bsOr}) != int64(15) {
		t.Errorf("Expected cardinality 15 after OR, got %d", bitsetCardinality([]interface{}{bsOr}))
	}

	// Test Xor
	bsXor := bitsetClone([]interface{}{bs1}).(*object.Object)
	bitsetXor([]interface{}{bsXor, bs2})
	if bitsetCardinality([]interface{}{bsXor}) != int64(10) {
		t.Errorf("Expected cardinality 10 after XOR, got %d", bitsetCardinality([]interface{}{bsXor}))
	}

	// Test AndNot
	bsAndNot := bitsetClone([]interface{}{bs1}).(*object.Object)
	bitsetAndNot([]interface{}{bsAndNot, bs2})
	if bitsetCardinality([]interface{}{bsAndNot}) != int64(5) {
		t.Errorf("Expected cardinality 5 after ANDNOT, got %d", bitsetCardinality([]interface{}{bsAndNot}))
	}
}

func TestBitSet_Search(t *testing.T) {
	setupTestGlobals()
	self := object.MakeEmptyObjectWithClassName(new("java/util/BitSet"))
	bitsetInit([]interface{}{self})

	bitsetSet([]interface{}{self, int64(10)})
	bitsetSet([]interface{}{self, int64(20)})

	if bitsetNextSetBit([]interface{}{self, int64(0)}) != int64(10) {
		t.Errorf("Expected nextSetBit(0) = 10")
	}
	if bitsetNextSetBit([]interface{}{self, int64(11)}) != int64(20) {
		t.Errorf("Expected nextSetBit(11) = 20")
	}
	if bitsetNextSetBit([]interface{}{self, int64(21)}) != int64(-1) {
		t.Errorf("Expected nextSetBit(21) = -1")
	}

	if bitsetNextClearBit([]interface{}{self, int64(10)}) != int64(11) {
		t.Errorf("Expected nextClearBit(10) = 11")
	}

	if bitsetPreviousSetBit([]interface{}{self, int64(25)}) != int64(20) {
		t.Errorf("Expected previousSetBit(25) = 20")
	}
	if bitsetPreviousSetBit([]interface{}{self, int64(19)}) != int64(10) {
		t.Errorf("Expected previousSetBit(19) = 10")
	}
	if bitsetPreviousSetBit([]interface{}{self, int64(9)}) != int64(-1) {
		t.Errorf("Expected previousSetBit(9) = -1")
	}
}

func TestBitSet_Conversion(t *testing.T) {
	setupTestGlobals()
	self := object.MakeEmptyObjectWithClassName(new("java/util/BitSet"))
	bitsetInit([]interface{}{self})

	bitsetSet([]interface{}{self, int64(0)})
	bitsetSet([]interface{}{self, int64(10)})

	bytes := bitsetToByteArray([]interface{}{self}).([]byte)
	if len(bytes) != 2 {
		t.Errorf("Expected 2 bytes, got %d", len(bytes))
	}

	longs := bitsetToLongArray([]interface{}{self}).([]int64)
	if len(longs) != 1 {
		t.Errorf("Expected 1 long, got %d", len(longs))
	}

	s := bitsetToString([]interface{}{self}).(string)
	if s != "{0, 10}" {
		t.Errorf("Expected \"{0, 10}\", got %q", s)
	}
}
