/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"math/bits"
	"strings"
)

// java.util.BitSet implementation for Jacobin.
// We store the bits as a slice of uint64 in the object's FieldTable.

func Load_Util_BitSet() {
	ghelpers.MethodSignatures["java/util/BitSet.<init>()V"] = ghelpers.GMeth{ParamSlots: 0, GFunction: bitsetInit}
	ghelpers.MethodSignatures["java/util/BitSet.<init>(I)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: bitsetInitWithCapacity}
	ghelpers.MethodSignatures["java/util/BitSet.and(Ljava/util/BitSet;)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: bitsetAnd}
	ghelpers.MethodSignatures["java/util/BitSet.andNot(Ljava/util/BitSet;)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: bitsetAndNot}
	ghelpers.MethodSignatures["java/util/BitSet.cardinality()I"] = ghelpers.GMeth{ParamSlots: 0, GFunction: bitsetCardinality}
	ghelpers.MethodSignatures["java/util/BitSet.clear()V"] = ghelpers.GMeth{ParamSlots: 0, GFunction: bitsetClearAll}
	ghelpers.MethodSignatures["java/util/BitSet.clear(I)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: bitsetClear}
	ghelpers.MethodSignatures["java/util/BitSet.clear(II)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: bitsetClearRange}
	ghelpers.MethodSignatures["java/util/BitSet.clone()Ljava/lang/Object;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: bitsetClone}
	ghelpers.MethodSignatures["java/util/BitSet.equals(Ljava/lang/Object;)Z"] = ghelpers.GMeth{ParamSlots: 1, GFunction: bitsetEquals}
	ghelpers.MethodSignatures["java/util/BitSet.flip(I)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: bitsetFlip}
	ghelpers.MethodSignatures["java/util/BitSet.flip(II)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: bitsetFlipRange}
	ghelpers.MethodSignatures["java/util/BitSet.get(I)Z"] = ghelpers.GMeth{ParamSlots: 1, GFunction: bitsetGet}
	ghelpers.MethodSignatures["java/util/BitSet.get(II)Ljava/util/BitSet;"] = ghelpers.GMeth{ParamSlots: 2, GFunction: bitsetGetRange}
	ghelpers.MethodSignatures["java/util/BitSet.hashCode()I"] = ghelpers.GMeth{ParamSlots: 0, GFunction: bitsetHashCode}
	ghelpers.MethodSignatures["java/util/BitSet.intersects(Ljava/util/BitSet;)Z"] = ghelpers.GMeth{ParamSlots: 1, GFunction: bitsetIntersects}
	ghelpers.MethodSignatures["java/util/BitSet.isEmpty()Z"] = ghelpers.GMeth{ParamSlots: 0, GFunction: bitsetIsEmpty}
	ghelpers.MethodSignatures["java/util/BitSet.length()I"] = ghelpers.GMeth{ParamSlots: 0, GFunction: bitsetLength}
	ghelpers.MethodSignatures["java/util/BitSet.nextClearBit(I)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: bitsetNextClearBit}
	ghelpers.MethodSignatures["java/util/BitSet.nextSetBit(I)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: bitsetNextSetBit}
	ghelpers.MethodSignatures["java/util/BitSet.or(Ljava/util/BitSet;)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: bitsetOr}
	ghelpers.MethodSignatures["java/util/BitSet.previousClearBit(I)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: bitsetPreviousClearBit}
	ghelpers.MethodSignatures["java/util/BitSet.previousSetBit(I)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: bitsetPreviousSetBit}
	ghelpers.MethodSignatures["java/util/BitSet.set(I)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: bitsetSet}
	ghelpers.MethodSignatures["java/util/BitSet.set(IZ)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: bitsetSetWithValue}
	ghelpers.MethodSignatures["java/util/BitSet.set(II)V"] = ghelpers.GMeth{ParamSlots: 2, GFunction: bitsetSetRange}
	ghelpers.MethodSignatures["java/util/BitSet.set(IIZ)V"] = ghelpers.GMeth{ParamSlots: 3, GFunction: bitsetSetRangeWithValue}
	ghelpers.MethodSignatures["java/util/BitSet.size()I"] = ghelpers.GMeth{ParamSlots: 0, GFunction: bitsetSize}
	ghelpers.MethodSignatures["java/util/BitSet.toByteArray()[B"] = ghelpers.GMeth{ParamSlots: 0, GFunction: bitsetToByteArray}
	ghelpers.MethodSignatures["java/util/BitSet.toLongArray()[J"] = ghelpers.GMeth{ParamSlots: 0, GFunction: bitsetToLongArray}
	ghelpers.MethodSignatures["java/util/BitSet.toString()Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: bitsetToString}
	ghelpers.MethodSignatures["java/util/BitSet.xor(Ljava/util/BitSet;)V"] = ghelpers.GMeth{ParamSlots: 1, GFunction: bitsetXor}

	// Static valueOf methods
	ghelpers.MethodSignatures["java/util/BitSet.valueOf([B)Ljava/util/BitSet;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: bitsetValueOfBytes}
	ghelpers.MethodSignatures["java/util/BitSet.valueOf([J)Ljava/util/BitSet;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: bitsetValueOfLongs}

	// Traps
	ghelpers.MethodSignatures["java/util/BitSet.stream()Ljava/util/stream/IntStream;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/BitSet.valueOf(Ljava/nio/ByteBuffer;)Ljava/util/BitSet;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/BitSet.valueOf(Ljava/nio/LongBuffer;)Ljava/util/BitSet;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
}

const bitsetStorage = "bits"

func getBitsFromObject(self *object.Object) ([]uint64, interface{}) {
	field, exists := self.FieldTable[bitsetStorage]
	if !exists {
		return nil, ghelpers.GetGErrBlk(excNames.NullPointerException, "BitSet storage not initialized")
	}
	bits, ok := field.Fvalue.([]uint64)
	if !ok {
		return nil, ghelpers.GetGErrBlk(excNames.VirtualMachineError, "Invalid BitSet storage")
	}
	return bits, nil
}

func setBitsToObject(self *object.Object, bits []uint64) {
	self.FieldTable[bitsetStorage] = object.Field{Ftype: types.RawGoPointer, Fvalue: bits}
}

func bitsetInit(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	setBitsToObject(self, make([]uint64, 1)) // Default capacity: 1 long (64 bits)
	return nil
}

func bitsetInitWithCapacity(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	nbits, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Invalid nbits")
	}
	if nbits < 0 {
		return ghelpers.GetGErrBlk(excNames.NegativeArraySizeException, fmt.Sprintf("nbits < 0: %d", nbits))
	}
	words := (nbits + 63) / 64
	if words == 0 {
		words = 1
	}
	setBitsToObject(self, make([]uint64, words))
	return nil
}

func bitsetGet(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	bitIndex, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Invalid bitIndex")
	}
	if bitIndex < 0 {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("bitIndex < 0: %d", bitIndex))
	}

	bits, err := getBitsFromObject(self)
	if err != nil {
		return err
	}

	wordIndex := bitIndex / 64
	if wordIndex >= int64(len(bits)) {
		return types.JavaBoolFalse
	}

	val := (bits[wordIndex] & (uint64(1) << (bitIndex % 64))) != 0
	return types.ConvertGoBoolToJavaBool(val)
}

func ensureCapacity(bits []uint64, bitIndex int64) []uint64 {
	wordIndex := bitIndex / 64
	if wordIndex >= int64(len(bits)) {
		newBits := make([]uint64, max(int64(len(bits))*2, wordIndex+1))
		copy(newBits, bits)
		return newBits
	}
	return bits
}

func bitsetSet(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	bitIndex, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Invalid bitIndex")
	}
	if bitIndex < 0 {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("bitIndex < 0: %d", bitIndex))
	}

	bits, err := getBitsFromObject(self)
	if err != nil {
		return err
	}

	bits = ensureCapacity(bits, bitIndex)
	wordIndex := bitIndex / 64
	bits[wordIndex] |= (uint64(1) << (bitIndex % 64))
	setBitsToObject(self, bits)
	return nil
}

func bitsetSetWithValue(params []interface{}) interface{} {
	value, ok2 := params[2].(int64)
	if !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Invalid value argument")
	}
	if value != 0 {
		return bitsetSet(params[:2])
	} else {
		return bitsetClear(params[:2])
	}
}

func bitsetClear(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	bitIndex, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Invalid bitIndex")
	}
	if bitIndex < 0 {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("bitIndex < 0: %d", bitIndex))
	}

	bits, err := getBitsFromObject(self)
	if err != nil {
		return err
	}

	wordIndex := bitIndex / 64
	if wordIndex < int64(len(bits)) {
		bits[wordIndex] &= ^(uint64(1) << (bitIndex % 64))
		setBitsToObject(self, bits)
	}
	return nil
}

func bitsetClearAll(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	bits, err := getBitsFromObject(self)
	if err != nil {
		return err
	}
	for i := range bits {
		bits[i] = 0
	}
	setBitsToObject(self, bits)
	return nil
}

func bitsetSize(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	bits, err := getBitsFromObject(self)
	if err != nil {
		return err
	}
	return int64(len(bits) * 64)
}

func bitsetLength(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	bitsData, err := getBitsFromObject(self)
	if err != nil {
		return err
	}
	for i := len(bitsData) - 1; i >= 0; i-- {
		if bitsData[i] != 0 {
			return int64(i)*64 + int64(64-bits.LeadingZeros64(bitsData[i]))
		}
	}
	return int64(0)
}

func bitsetCardinality(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	bitsData, err := getBitsFromObject(self)
	if err != nil {
		return err
	}
	count := 0
	for _, b := range bitsData {
		count += bits.OnesCount64(b)
	}
	return int64(count)
}

func bitsetIsEmpty(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	bits, err := getBitsFromObject(self)
	if err != nil {
		return err
	}
	for _, b := range bits {
		if b != 0 {
			return types.JavaBoolFalse
		}
	}
	return types.JavaBoolTrue
}

func bitsetFlip(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	bitIndex, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Invalid bitIndex")
	}
	if bitIndex < 0 {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("bitIndex < 0: %d", bitIndex))
	}

	bits, err := getBitsFromObject(self)
	if err != nil {
		return err
	}

	bits = ensureCapacity(bits, bitIndex)
	wordIndex := bitIndex / 64
	bits[wordIndex] ^= (uint64(1) << (bitIndex % 64))
	setBitsToObject(self, bits)
	return nil
}

func bitsetSetRange(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	fromIndex, ok1 := params[1].(int64)
	toIndex, ok2 := params[2].(int64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Invalid indices")
	}
	if fromIndex < 0 || toIndex < 0 || fromIndex > toIndex {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("fromIndex: %d, toIndex: %d", fromIndex, toIndex))
	}
	if fromIndex == toIndex {
		return nil
	}

	bits, err := getBitsFromObject(self)
	if err != nil {
		return err
	}

	bits = ensureCapacity(bits, toIndex-1)
	startWord := fromIndex / 64
	endWord := (toIndex - 1) / 64

	firstWordMask := uint64(0xffffffffffffffff) << (fromIndex % 64)
	lastWordMask := uint64(0xffffffffffffffff) >> (63 - ((toIndex - 1) % 64))

	if startWord == endWord {
		bits[startWord] |= (firstWordMask & lastWordMask)
	} else {
		bits[startWord] |= firstWordMask
		for i := startWord + 1; i < endWord; i++ {
			bits[i] = uint64(0xffffffffffffffff)
		}
		bits[endWord] |= lastWordMask
	}

	setBitsToObject(self, bits)
	return nil
}

func bitsetSetRangeWithValue(params []interface{}) interface{} {
	value := params[3].(int64)

	if value != 0 {
		return bitsetSetRange(params[:3])
	} else {
		return bitsetClearRange(params[:3])
	}
}

func bitsetClearRange(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	fromIndex, ok1 := params[1].(int64)
	toIndex, ok2 := params[2].(int64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Invalid indices")
	}
	if fromIndex < 0 || toIndex < 0 || fromIndex > toIndex {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("fromIndex: %d, toIndex: %d", fromIndex, toIndex))
	}
	if fromIndex == toIndex {
		return nil
	}

	bits, err := getBitsFromObject(self)
	if err != nil {
		return err
	}

	if fromIndex >= int64(len(bits))*64 {
		return nil
	}

	startWord := fromIndex / 64
	endWord := (toIndex - 1) / 64

	firstWordMask := uint64(0xffffffffffffffff) << (fromIndex % 64)
	lastWordMask := uint64(0xffffffffffffffff) >> (63 - ((toIndex - 1) % 64))

	if endWord >= int64(len(bits)) {
		endWord = int64(len(bits)) - 1
		lastWordMask = uint64(0xffffffffffffffff)
	}

	if startWord == endWord {
		bits[startWord] &= ^(firstWordMask & lastWordMask)
	} else {
		bits[startWord] &= ^firstWordMask
		for i := startWord + 1; i <= endWord; i++ {
			if i == endWord {
				bits[i] &= ^lastWordMask
			} else {
				bits[i] = 0
			}
		}
	}

	setBitsToObject(self, bits)
	return nil
}

func bitsetFlipRange(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	fromIndex, ok1 := params[1].(int64)
	toIndex, ok2 := params[2].(int64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Invalid indices")
	}
	if fromIndex < 0 || toIndex < 0 || fromIndex > toIndex {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("fromIndex: %d, toIndex: %d", fromIndex, toIndex))
	}
	if fromIndex == toIndex {
		return nil
	}

	bits, err := getBitsFromObject(self)
	if err != nil {
		return err
	}

	bits = ensureCapacity(bits, toIndex-1)
	startWord := fromIndex / 64
	endWord := (toIndex - 1) / 64

	firstWordMask := uint64(0xffffffffffffffff) << (fromIndex % 64)
	lastWordMask := uint64(0xffffffffffffffff) >> (63 - ((toIndex - 1) % 64))

	if startWord == endWord {
		bits[startWord] ^= (firstWordMask & lastWordMask)
	} else {
		bits[startWord] ^= firstWordMask
		for i := startWord + 1; i < endWord; i++ {
			bits[i] ^= uint64(0xffffffffffffffff)
		}
		bits[endWord] ^= lastWordMask
	}

	setBitsToObject(self, bits)
	return nil
}

func bitsetGetRange(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	fromIndex, ok1 := params[1].(int64)
	toIndex, ok2 := params[2].(int64)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Invalid indices")
	}
	if fromIndex < 0 || toIndex < 0 || fromIndex > toIndex {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("fromIndex: %d, toIndex: %d", fromIndex, toIndex))
	}

	newBitSetRaw, _ := globals.GetGlobalRef().FuncInstantiateClass("java/util/BitSet", nil)
	newBitSet := newBitSetRaw.(*object.Object)
	bitsetInitWithCapacity([]interface{}{newBitSet, toIndex - fromIndex})

	if fromIndex == toIndex {
		return newBitSet
	}

	bits, _ := getBitsFromObject(self)
	newBits, _ := getBitsFromObject(newBitSet)

	for i := fromIndex; i < toIndex; i++ {
		wordIdx := i / 64
		if wordIdx < int64(len(bits)) {
			if (bits[wordIdx] & (uint64(1) << (i % 64))) != 0 {
				targetIdx := i - fromIndex
				newBits[targetIdx/64] |= (uint64(1) << (targetIdx % 64))
			}
		}
	}
	// Optimization possible for word-aligned copies but this is simpler for now.
	setBitsToObject(newBitSet, newBits)
	return newBitSet
}

func bitsetAnd(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)
	bits, _ := getBitsFromObject(self)
	otherBits, _ := getBitsFromObject(other)

	for i := range bits {
		if i < len(otherBits) {
			bits[i] &= otherBits[i]
		} else {
			bits[i] = 0
		}
	}
	setBitsToObject(self, bits)
	return nil
}

func bitsetAndNot(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)
	bits, _ := getBitsFromObject(self)
	otherBits, _ := getBitsFromObject(other)

	for i := range bits {
		if i < len(otherBits) {
			bits[i] &= ^otherBits[i]
		}
	}
	setBitsToObject(self, bits)
	return nil
}

func bitsetOr(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)
	bits, _ := getBitsFromObject(self)
	otherBits, _ := getBitsFromObject(other)

	if len(otherBits) > len(bits) {
		bits = ensureCapacity(bits, int64(len(otherBits))*64-1)
	}

	for i := range otherBits {
		bits[i] |= otherBits[i]
	}
	setBitsToObject(self, bits)
	return nil
}

func bitsetXor(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)
	bits, _ := getBitsFromObject(self)
	otherBits, _ := getBitsFromObject(other)

	if len(otherBits) > len(bits) {
		bits = ensureCapacity(bits, int64(len(otherBits))*64-1)
	}

	for i := range otherBits {
		bits[i] ^= otherBits[i]
	}
	setBitsToObject(self, bits)
	return nil
}

func bitsetIntersects(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)
	bits, _ := getBitsFromObject(self)
	otherBits, _ := getBitsFromObject(other)

	limit := min(len(bits), len(otherBits))
	for i := range limit {
		if (bits[i] & otherBits[i]) != 0 {
			return types.JavaBoolTrue
		}
	}
	return types.JavaBoolFalse
}

func bitsetNextSetBit(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	fromIndex, _ := params[1].(int64)
	if fromIndex < 0 {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("fromIndex < 0: %d", fromIndex))
	}

	bitsData, _ := getBitsFromObject(self)
	wordIndex := fromIndex / 64
	if wordIndex >= int64(len(bitsData)) {
		return int64(-1)
	}

	word := bitsData[wordIndex] & (uint64(0xffffffffffffffff) << (fromIndex % 64))
	for {
		if word != 0 {
			return wordIndex*64 + int64(bits.TrailingZeros64(word))
		}
		wordIndex++
		if wordIndex >= int64(len(bitsData)) {
			return int64(-1)
		}
		word = bitsData[wordIndex]
	}
}

func bitsetNextClearBit(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	fromIndex, _ := params[1].(int64)
	if fromIndex < 0 {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("fromIndex < 0: %d", fromIndex))
	}

	bitsData, _ := getBitsFromObject(self)
	wordIndex := fromIndex / 64
	if wordIndex >= int64(len(bitsData)) {
		return fromIndex
	}

	word := (^bitsData[wordIndex]) & (uint64(0xffffffffffffffff) << (fromIndex % 64))
	for {
		if word != 0 {
			return wordIndex*64 + int64(bits.TrailingZeros64(word))
		}
		wordIndex++
		if wordIndex >= int64(len(bitsData)) {
			return int64(len(bitsData)) * 64
		}
		word = ^bitsData[wordIndex]
	}
}

func bitsetPreviousSetBit(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	fromIndex, _ := params[1].(int64)
	if fromIndex < 0 {
		if fromIndex == -1 {
			return int64(-1)
		}
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("fromIndex < -1: %d", fromIndex))
	}

	bitsData, _ := getBitsFromObject(self)
	wordIndex := fromIndex / 64
	if wordIndex >= int64(len(bitsData)) {
		wordIndex = int64(len(bitsData)) - 1
		for wordIndex >= 0 && bitsData[wordIndex] == 0 {
			wordIndex--
		}
		if wordIndex < 0 {
			return int64(-1)
		}
		return wordIndex*64 + int64(63-bits.LeadingZeros64(bitsData[wordIndex]))
	}

	word := bitsData[wordIndex] & (uint64(0xffffffffffffffff) >> (63 - (fromIndex % 64)))
	for {
		if word != 0 {
			return wordIndex*64 + int64(63-bits.LeadingZeros64(word))
		}
		wordIndex--
		if wordIndex < 0 {
			return int64(-1)
		}
		word = bitsData[wordIndex]
	}
}

func bitsetPreviousClearBit(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	fromIndex, _ := params[1].(int64)
	if fromIndex < 0 {
		if fromIndex == -1 {
			return int64(-1)
		}
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("fromIndex < -1: %d", fromIndex))
	}

	bitsData, _ := getBitsFromObject(self)
	wordIndex := fromIndex / 64
	if wordIndex >= int64(len(bitsData)) {
		return fromIndex
	}

	word := (^bitsData[wordIndex]) & (uint64(0xffffffffffffffff) >> (63 - (fromIndex % 64)))
	for {
		if word != 0 {
			return wordIndex*64 + int64(63-bits.LeadingZeros64(word))
		}
		wordIndex--
		if wordIndex < 0 {
			return int64(-1)
		}
		word = ^bitsData[wordIndex]
	}
}

func bitsetHashCode(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	bits, _ := getBitsFromObject(self)
	h := uint64(1234)
	for i := len(bits) - 1; i >= 0; i-- {
		h ^= bits[i] * uint64(i+1)
	}
	return int64(int32(h>>32 ^ h))
}

func bitsetEquals(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	otherObj, ok := params[1].(*object.Object)
	if !ok || otherObj == nil {
		return types.JavaBoolFalse
	}
	// Simplified check: assume otherObj is BitSet if it has bitsetStorage
	bits, _ := getBitsFromObject(self)
	otherBits, err := getBitsFromObject(otherObj)
	if err != nil {
		return types.JavaBoolFalse
	}

	maxLen := max(len(bits), len(otherBits))
	for i := range maxLen {
		var b1, b2 uint64
		if i < len(bits) {
			b1 = bits[i]
		}
		if i < len(otherBits) {
			b2 = otherBits[i]
		}
		if b1 != b2 {
			return types.JavaBoolFalse
		}
	}
	return types.JavaBoolTrue
}

func bitsetClone(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	bitsData, _ := getBitsFromObject(self)
	newBitSetRaw, _ := globals.GetGlobalRef().FuncInstantiateClass("java/util/BitSet", nil)
	newBitSet := newBitSetRaw.(*object.Object)
	newBits := make([]uint64, len(bitsData))
	copy(newBits, bitsData)
	setBitsToObject(newBitSet, newBits)
	return newBitSet
}

func bitsetToString(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	bits, _ := getBitsFromObject(self)
	var sb strings.Builder
	sb.WriteString("{")
	first := true
	for i, word := range bits {
		if word == 0 {
			continue
		}
		for b := range 64 {
			if (word & (uint64(1) << b)) != 0 {
				if !first {
					sb.WriteString(", ")
				}
				fmt.Fprintf(&sb, "%d", i*64+b)
				first = false
			}
		}
	}
	sb.WriteString("}")
	return sb.String()
}

func bitsetToByteArray(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	bitsData, _ := getBitsFromObject(self)
	// Find last non-zero word
	lastWord := -1
	for i := len(bitsData) - 1; i >= 0; i-- {
		if bitsData[i] != 0 {
			lastWord = i
			break
		}
	}
	if lastWord == -1 {
		return make([]byte, 0)
	}

	nbytes := lastWord*8 + (64-bits.LeadingZeros64(bitsData[lastWord])+7)/8
	res := make([]byte, nbytes)
	for i := range nbytes {
		res[i] = byte(bitsData[i/8] >> (8 * (i % 8)))
	}
	return res
}

func bitsetToLongArray(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	bits, _ := getBitsFromObject(self)
	// Trim trailing zeros
	lastWord := -1
	for i := len(bits) - 1; i >= 0; i-- {
		if bits[i] != 0 {
			lastWord = i
			break
		}
	}
	if lastWord == -1 {
		return make([]int64, 0)
	}
	res := make([]int64, lastWord+1)
	for i := range lastWord + 1 {
		res[i] = int64(bits[i])
	}
	return res
}

func bitsetValueOfBytes(params []interface{}) interface{} {
	bytes, ok := params[0].([]byte)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Invalid byte array")
	}
	nlongs := (len(bytes) + 7) / 8
	bitsData := make([]uint64, nlongs)
	for i, b := range bytes {
		bitsData[i/8] |= uint64(b) << (8 * (i % 8))
	}
	// We need a classloader to instantiate BitSet.
	// In GFunctions, usually params[0] is self, but for static methods it depends.
	// valueOf([B) is static. How to get classloader?
	// For now, use nil or try to get from somewhere if available.
	// Actually, static methods in Jacobin often don't have self.
	// Let's assume there is a way to get it or it's provided.
	// Standard Jacobin practice for static factory methods?

	newBitSetRaw, _ := globals.GetGlobalRef().FuncInstantiateClass("java/util/BitSet", nil)
	newBitSet := newBitSetRaw.(*object.Object)
	setBitsToObject(newBitSet, bitsData)
	return newBitSet
}

func bitsetValueOfLongs(params []interface{}) interface{} {
	longs, ok := params[0].([]int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Invalid long array")
	}
	bitsData := make([]uint64, len(longs))
	for i, l := range longs {
		bitsData[i] = uint64(l)
	}
	newBitSetRaw, _ := globals.GetGlobalRef().FuncInstantiateClass("java/util/BitSet", nil)
	newBitSet := newBitSetRaw.(*object.Object)
	setBitsToObject(newBitSet, bitsData)
	return newBitSet
}
