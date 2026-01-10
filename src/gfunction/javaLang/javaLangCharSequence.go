/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
)

func Load_Lang_CharSequence() {

	ghelpers.MethodSignatures["java/lang/CharSequence.compare(Ljava/lang/CharSequence;Ljava/lang/CharSequence;)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  charSequenceCompare,
		}

	ghelpers.MethodSignatures["java/lang/CharSequence.length()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  charSequenceLength,
		}

	ghelpers.MethodSignatures["java/lang/CharSequence.charAt(I)C"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  charSequenceCharAt,
		}

	ghelpers.MethodSignatures["java/lang/CharSequence.subSequence(II)Ljava/lang/CharSequence;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  charSequenceSubSequence,
		}

	ghelpers.MethodSignatures["java/lang/CharSequence.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  charSequenceToString,
		}
}

// charSequenceLength returns the length of the CharSequence.
func charSequenceLength(params []interface{}) interface{} {
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "charSequenceLength: self is nil")
	}
	obj := params[0].(*object.Object)
	// CharSequence is an interface. We need to handle known implementations.
	// In Jacobin, String objects have their value in a specific format.
	if object.IsStringObject(obj) {
		str := object.GoStringFromStringObject(obj)
		return int64(len(str))
	}

	className := object.GoStringFromStringPoolIndex(obj.KlassName)
	if className == "java/lang/StringBuilder" || className == "java/lang/StringBuffer" {
		countFld, ok := obj.FieldTable["count"]
		if ok {
			return countFld.Fvalue.(int64)
		}
	}

	return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, "charSequenceLength: unknown implementation")
}

// charSequenceCharAt returns the char value at the specified index.
func charSequenceCharAt(params []interface{}) interface{} {
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "charSequenceCharAt: self is nil")
	}
	obj := params[0].(*object.Object)
	index := params[1].(int64)

	if object.IsStringObject(obj) {
		str := object.GoStringFromStringObject(obj)
		if index < 0 || index >= int64(len(str)) {
			return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, "charSequenceCharAt: index out of bounds")
		}
		return int64(str[index])
	}

	className := object.GoStringFromStringPoolIndex(obj.KlassName)
	if className == "java/lang/StringBuilder" || className == "java/lang/StringBuffer" {
		fld, ok := obj.FieldTable["value"]
		if ok {
			byteArray := fld.Fvalue.([]types.JavaByte)
			if index < 0 || index >= int64(len(byteArray)) {
				return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, "charSequenceCharAt: index out of bounds")
			}
			return int64(byteArray[index])
		}
	}

	return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, "charSequenceCharAt: unknown implementation")
}

// charSequenceSubSequence returns a CharSequence that is a subsequence of this sequence.
func charSequenceSubSequence(params []interface{}) interface{} {
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "charSequenceSubSequence: self is nil")
	}
	obj := params[0].(*object.Object)
	start := params[1].(int64)
	end := params[2].(int64)

	if object.IsStringObject(obj) {
		str := object.GoStringFromStringObject(obj)
		if start < 0 || end > int64(len(str)) || start > end {
			return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, "charSequenceSubSequence: bounds out of range")
		}
		subStr := str[start:end]
		return object.StringObjectFromGoString(subStr)
	}

	return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, "charSequenceSubSequence: unknown implementation")
}

// charSequenceToString returns a string containing the characters in this sequence.
func charSequenceToString(params []interface{}) interface{} {
	if params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "charSequenceToString: self is nil")
	}
	obj := params[0].(*object.Object)
	if object.IsStringObject(obj) {
		return obj
	}

	return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, "charSequenceToString: unknown implementation")
}

// charSequenceCompare compares two CharSequence instances lexicographically.
func charSequenceCompare(params []interface{}) interface{} {
	if params[0] == nil || params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "charSequenceCompare: null parameter")
	}

	cs1Obj := params[0].(*object.Object)
	cs2Obj := params[1].(*object.Object)

	str1 := object.GoStringFromStringObject(cs1Obj)
	str2 := object.GoStringFromStringObject(cs2Obj)

	if str1 < str2 {
		return int64(-1)
	} else if str1 > str2 {
		return int64(1)
	}
	return int64(0)
}
