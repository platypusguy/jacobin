/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"strconv"
)

// Implementation of some of the functions in Java/lang/Class.

func Load_Lang_StringBuilder() {

	// === Instantiation ===

	ghelpers.MethodSignatures["java/lang/StringBuilder.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringBuilderInit,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.<init>(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderInit,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.<init>(Ljava/lang/CharSequence;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderInitString,
		}

	// === Methods ===

	ghelpers.MethodSignatures["java/lang/StringBuilder.append(Z)Ljava/lang/StringBuilder;"] = // append boolean
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppendBoolean,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.append(C)Ljava/lang/StringBuilder;"] = // append char
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppendChar,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.append([C)Ljava/lang/StringBuilder;"] = // append char array
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.append([CII)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  stringBuilderAppend,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.append(D)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.append(F)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.append(I)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.append(J)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.append(Ljava/lang/CharSequence;)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.append(Ljava/lang/CharSequence;II)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.append(Ljava/lang/Object;)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.append(Ljava/lang/String;)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.append(Ljava/lang/StringBuffer;)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.appendCodePoint(I)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.capacity()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringBuilderCapacity,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.charAt(I)C"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderCharAt,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.chars()Ljava/util/stream/IntStream;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.codePointAt(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.codePointBefore(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.codePointCount(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.codePoints()Ljava/util/stream/IntStream;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.compareTo(Ljava/lang/StringBuilder;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringCompareToCaseSensitive,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.delete(II)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderDelete,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.deleteCharAt(I)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderDelete,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.ensureCapacity(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.getChars(II[CI)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.indexOf(Ljava/lang/String;)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.indexOf(Ljava/lang/String;I)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.insert(IZ)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsertBoolean,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.insert(IC)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsertChar,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.insert(I[C)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.insert(I[CII)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  stringBuilderInsert,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.insert(ID)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.insert(IF)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.insert(II)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.insert(IJ)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.insert(ILjava/lang/CharSequence;)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.insert(ILjava/lang/CharSequence;II)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.insert(ILjava/lang/Object;)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.insert(ILjava/lang/String;)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}
	ghelpers.MethodSignatures["java/lang/StringBuilder.isLatin1()Z"] = // internal member function, not in API
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnTrue,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.lastIndexOf(Ljava/lang/String;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.lastIndexOf(Ljava/lang/String;I)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.length()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringBuilderLength,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.offsetByCodePoints(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.replace(IILjava/lang/String;)Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  stringBuilderReplace,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.reverse()Ljava/lang/StringBuilder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringBuilderReverse,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.setCharAt(IC)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderSetCharAt,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.setLength(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderSetLength,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.subSequence(II)Ljava/lang/CharSequence;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	// Return a substring starting at the given index of the byte array.
	ghelpers.MethodSignatures["java/lang/StringBuilder.substring(I)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  substringToTheEnd, // javaLangString.go
		}

	// Return a substring starting at the given index of the byte array of the given length.
	ghelpers.MethodSignatures["java/lang/StringBuilder.substring(II)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  substringStartEnd, // javaLangString.go
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringBuilderToString,
		}

	ghelpers.MethodSignatures["java/lang/StringBuilder.trimToSize()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

}

var classStringBuilder = "java/lang/StringBuilder"

// Initialise StringBuilder with or without a capacity integer.
func stringBuilderInit(params []any) any {
	// Get File object and initialise the field map.
	obj := params[0].(*object.Object)
	obj.FieldTable = make(map[string]object.Field)

	// Set the count = 0.
	fld := object.Field{Ftype: types.Int, Fvalue: int64(0)}
	obj.FieldTable["count"] = fld

	// Set the value = nil byte array.
	fld = object.Field{Ftype: types.ByteArray, Fvalue: make([]types.JavaByte, 0)}
	obj.FieldTable["value"] = fld

	// Set the capacity field value.
	var capacity int64
	if len(params) > 1 { // Was a capacity parameter supplied?
		capacity = params[1].(int64)
	} else {
		capacity = 16 // default capacity value per API
	}
	fld = object.Field{Ftype: types.Int, Fvalue: capacity}
	obj.FieldTable["capacity"] = fld

	return nil
}

// Initialise StringBuilder with a String object.
func stringBuilderInitString(params []any) any {
	// Get File object and initialise the field map.
	obj := params[0].(*object.Object)
	obj.FieldTable = make(map[string]object.Field)

	var byteArray []types.JavaByte
	switch params[1].(type) {
	case *object.Object: // String
		rawArray := params[1].(*object.Object).FieldTable["value"].Fvalue
		switch rawArray.(type) {
		case []types.JavaByte:
			byteArray = rawArray.([]types.JavaByte)
		case []byte:
			byteArray = object.JavaByteArrayFromGoByteArray(rawArray.([]byte))
		case string:
			byteArray = object.JavaByteArrayFromGoString(rawArray.(string))
		}
	default:
		errMsg := fmt.Sprintf("stringBuilderInitString: Parameter type (%T) is illegal", params[1])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Append parmArray to the byteArray.
	// Set the byte count.
	count := int64(len(byteArray))
	capacity := count + 16
	obj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: byteArray}
	obj.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: count}
	obj.FieldTable["capacity"] = object.Field{Ftype: types.Int, Fvalue: capacity}

	return nil
}

// stringBuilderAppend appends the second parameter to the bytes in the StringBuilder
// that is passed in the objectRef parameter (the first param).
//
// If a character array with offset and size parameters, there is special handling.
//
// Method parameter types:
// [C                          int64 array
// [CII                        int64 array, offset, size
// D                           float64
// F                           float64
// I                           int64
// J                           int64
// Ljava/lang/Object;          *object.Object
// Ljava/lang/String;          *object.Object
// Ljava/lang/StringBuffer;    *object.Object
func stringBuilderAppend(params []any) any {
	// Get base object and its value field, byteArray.
	objBase := params[0].(*object.Object)
	byteArray := objBase.FieldTable["value"].Fvalue.([]types.JavaByte)

	// Resolved parameter byte array, regardless of parameters:
	var parmArray []types.JavaByte

	// Process based primarily on the params[1] type.
	switch params[1].(type) {
	case *object.Object: // char array, Object, String, StringBuffer, or StringBuilder
		fvalue := params[1].(*object.Object).FieldTable["value"].Fvalue
		switch fvalue.(type) {
		case []types.JavaByte: // JavaByte array, String, StringBuffer, or StringBuilder
			parmArray = fvalue.([]types.JavaByte)
		case []byte: // byte array
			parmArray = object.JavaByteArrayFromGoByteArray(fvalue.([]byte))
		case []int64: // char array, int array
			if len(params) == 4 {
				int64Array := fvalue.([]int64)
				len64Array := int64(len(int64Array))
				start := params[2].(int64)
				length := params[3].(int64)
				end := start + length
				if start < 0 || start > len64Array || end <= start || end > len64Array {
					errMsg := fmt.Sprintf("stringBuilderAppend: Invalid offset (%d) or length (%d)", start, length)
					return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
				}
				for ix := start; ix < start+length; ix++ {
					parmArray = append(parmArray, types.JavaByte(int64Array[ix]))
				}
			} else { // Append the entire char array.
				for _, elem := range fvalue.([]int64) {
					parmArray = append(parmArray, types.JavaByte(elem))
				}
			}
		default:
			str := object.StringifyAnythingGo(fvalue)
			parmArray = object.JavaByteArrayFromGoString(str)
		}
	case int64: // int, long, short
		str := fmt.Sprintf("%d", params[1].(int64))
		parmArray = object.JavaByteArrayFromGoString(str)
	case float64: // float, double
		ff := params[1].(float64)
		str := strconv.FormatFloat(ff, 'f', -1, 64)
		parmArray = object.JavaByteArrayFromGoString(str)
	default:
		str := object.StringifyAnythingGo(params[1])
		parmArray = object.JavaByteArrayFromGoString(str)
	}

	// Append parmArray to the byteArray.
	// Set the byte count.
	byteArray = append(byteArray, parmArray...)
	count := int64(len(byteArray))
	objBase.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: byteArray}
	objBase.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: count}
	expandCapacity(objBase, count)

	return objBase
}

// Append the second parameter (boolean) to the bytes in the StringBuilder that is
// passed in the objectRef parameter (the first param).
func stringBuilderAppendBoolean(params []any) any {
	// Get base object and its value field, byteArray.
	objBase := params[0].(*object.Object)
	byteArray := objBase.FieldTable["value"].Fvalue.([]types.JavaByte)

	var parmArray []types.JavaByte
	switch params[1].(type) {
	case int64: // boolean
		var str string
		if params[1].(int64) == types.JavaBoolTrue {
			str = "true"
		} else {
			str = "false"
		}
		parmArray = object.JavaByteArrayFromGoString(str)
	default:
		errMsg := fmt.Sprintf("stringBuilderAppendBoolean: Parameter type (%T) is illegal", params[1])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Append parmArray to the byteArray.
	// Set the byte count.
	byteArray = append(byteArray, parmArray...)
	count := int64(len(byteArray))
	objBase.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: byteArray}
	objBase.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: count}
	expandCapacity(objBase, count)

	return objBase
}

// Append the second parameter (char) to the bytes in the StringBuilder that is
// passed in the objectRef parameter (the first param).
func stringBuilderAppendChar(params []any) any {
	// Get base object and its value field, byteArray.
	objBase := params[0].(*object.Object)
	byteArray := objBase.FieldTable["value"].Fvalue.([]types.JavaByte)

	var parmArray = make([]types.JavaByte, 1)
	switch params[1].(type) {
	case int64: // char
		bb := types.JavaByte(params[1].(int64))
		parmArray[0] = bb
	default:
		errMsg := fmt.Sprintf("stringBuilderAppendChar: Parameter type (%T) is illegal", params[1])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Append parmArray to the byteArray.
	// Set the byte count.
	byteArray = append(byteArray, parmArray...)
	count := int64(len(byteArray))
	objBase.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: byteArray}
	objBase.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: count}
	expandCapacity(objBase, count)

	return objBase
}

// Extract a character at the given index.
func stringBuilderCharAt(params []any) any {
	obj := params[0].(*object.Object)
	ix := params[1].(int64)
	bytes := obj.FieldTable["value"].Fvalue.([]types.JavaByte)
	if ix >= int64(len(bytes)) {
		errMsg := fmt.Sprintf("stringBuilderCharAt: Index value (%d) exceeds the byte array size (%d)", ix, len(bytes))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	return int64(bytes[ix])
}

// Removes the characters in a substring of the StringBuilder object. The substring begins at the specified start
// and extends to the character at index end - 1 or to the end of the sequence if no such character exists.
// If start is equal to end, no changes are made.
func stringBuilderDelete(params []any) any {
	objBase := params[0].(*object.Object)
	initBytes := objBase.FieldTable["value"].Fvalue.([]types.JavaByte)
	initLen := int64(len(initBytes))
	start := params[1].(int64)
	var end int64
	if len(params) == 3 {
		end = params[2].(int64) // delete(start, end)
	} else {
		end = start + 1 // deleteCharAt(offset)
	}

	// Validate start and end.
	if start < 0 || start > initLen {
		errMsg := fmt.Sprintf("stringBuilderDelete: Start value (%d) < 0 or exceeds the byte array size (%d)", start, initLen)
		return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
	}
	if end < start {
		errMsg := fmt.Sprintf("stringBuilderDelete: End value (%d) < Start value (%d)", start, end)
		return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
	}
	if end > initLen {
		end = initLen
	}

	// If start = end, return object as-is.
	if start == end {
		return objBase
	}

	// Copy retained bytes to a new byte array.
	newArray := make([]types.JavaByte, start)
	if start > 0 {
		copy(newArray, initBytes[0:start])
	}
	newArray = append(newArray, initBytes[end:]...)

	// New length of byte array --> count.
	count := int64(len(newArray))

	// Finalize output object.
	objBase.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: newArray}
	objBase.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: count}
	expandCapacity(objBase, count)

	return objBase
}

// Insert the second parameter to the bytes into the StringBuilder
// at the given index.
func stringBuilderInsert(params []any) any {
	// Get base object and its value field, byteArray.
	objBase := params[0].(*object.Object)
	byteArray := objBase.FieldTable["value"].Fvalue.([]types.JavaByte)

	// Get the index value.
	ix := params[1].(int64)
	if ix < 0 || ix > int64(len(byteArray)) {
		errMsg := fmt.Sprintf("stringBuilderInsert: Index value (%d) is negative or exceeds the byte array size (%d)",
			ix, len(byteArray))
		return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
	}

	var parmArray []types.JavaByte
	switch params[2].(type) {
	case *object.Object: // char array, String, StringBuffer, or StringBuilder
		fvalue := params[2].(*object.Object).FieldTable["value"].Fvalue
		switch fvalue.(type) {
		case []types.JavaByte: // String, StringBuffer, or StringBuilder
			parmArray = fvalue.([]types.JavaByte)
		case []byte: // byte array
			parmArray = object.JavaByteArrayFromGoByteArray(fvalue.([]byte))
		case []int64: // char array
			if len(params) == 5 { // subset of char array
				int64Array := fvalue.([]int64)
				len64Array := int64(len(int64Array))
				start := params[3].(int64)
				length := params[4].(int64)
				end := start + length
				if start < 0 || start > len64Array || end <= start || end > len64Array {
					errMsg := fmt.Sprintf("stringBuilderInsert: Invalid offset (%d) or length (%d)", start, length)
					return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
				}
				for ix := start; ix < start+length; ix++ {
					parmArray = append(parmArray, types.JavaByte(int64Array[ix]))
				}
			} else { // Append the entire char array.
				for _, elem := range fvalue.([]int64) {
					parmArray = append(parmArray, types.JavaByte(elem))
				}
			}
		default:
			errMsg := fmt.Sprintf("stringBuilderInsert: Object value field value type (%T) is not a byte array nor a char array",
				params[1])
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	case int64: // integer, long
		str := fmt.Sprintf("%d", params[2].(int64))
		parmArray = object.JavaByteArrayFromGoString(str)
	case float64: // float, double
		ff := params[2].(float64)
		str := strconv.FormatFloat(ff, 'f', -1, 64)
		parmArray = object.JavaByteArrayFromGoString(str)
	default:
		errMsg := fmt.Sprintf("stringBuilderInsert: Parameter type (%T) is illegal", params[1])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Append parmArray to the byteArray.
	// Set the byte count.
	newArray := make([]types.JavaByte, ix)
	if ix > 0 {
		copy(newArray, byteArray[0:ix])
	}
	newArray = append(newArray, parmArray...)
	newArray = append(newArray, byteArray[ix:]...)
	count := int64(len(newArray))

	objBase.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: newArray}
	objBase.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: count}
	expandCapacity(objBase, count)

	return objBase
}

// Insert the boolean parameter into the bytes into the StringBuilder
// at the given index.
func stringBuilderInsertBoolean(params []any) any {
	// Get base object and its value field, byteArray.
	objBase := params[0].(*object.Object)
	byteArray := objBase.FieldTable["value"].Fvalue.([]types.JavaByte)

	// Get the index value.
	ix := params[1].(int64)
	if ix < 0 || ix > int64(len(byteArray)) {
		errMsg := fmt.Sprintf("stringBuilderInsertBoolean: Index value (%d) is negative or exceeds the byte array size (%d)",
			ix, len(byteArray))
		return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
	}

	var parmArray []types.JavaByte
	switch params[2].(type) {
	case int64: // boolean
		var str string
		if params[2].(int64) == types.JavaBoolTrue {
			str = "true"
		} else {
			str = "false"
		}
		parmArray = object.JavaByteArrayFromGoString(str)
	default:
		errMsg := fmt.Sprintf("stringBuilderInsertBoolean: Parameter type (%T) is illegal", params[1])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Append parmArray to the byteArray.
	// Set the byte count.
	newArray := make([]types.JavaByte, ix)
	if ix > 0 {
		copy(newArray, byteArray[0:ix])
	}
	newArray = append(newArray, parmArray...)
	newArray = append(newArray, byteArray[ix:]...)
	count := int64(len(newArray))

	objBase.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: newArray}
	objBase.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: count}
	expandCapacity(objBase, count)

	return objBase
}

// Insert the char parameter into the bytes into the StringBuilder
// at the given index.
func stringBuilderInsertChar(params []any) any {
	// Get base object and its value field, byteArray.
	objBase := params[0].(*object.Object)
	byteArray := objBase.FieldTable["value"].Fvalue.([]types.JavaByte)

	// Get the index value.
	ix := params[1].(int64)
	if ix < 0 || ix > int64(len(byteArray)) {
		errMsg := fmt.Sprintf("stringBuilderInsertChar: Index value (%d) is negative or exceeds the byte array size (%d)",
			ix, len(byteArray))
		return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
	}

	var bb types.JavaByte
	switch params[2].(type) {
	case int64: // char
		bb = types.JavaByte(params[2].(int64))
	default:
		errMsg := fmt.Sprintf("stringBuilderInsertChar: Parameter type (%T) is illegal", params[1])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Append parmArray to the byteArray.
	// Set the byte count.
	newArray := make([]types.JavaByte, ix)
	if ix > 0 {
		copy(newArray, byteArray[0:ix])
	}
	newArray = append(newArray, bb)
	newArray = append(newArray, byteArray[ix:]...)
	count := int64(len(newArray))

	objBase.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: newArray}
	objBase.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: count}
	expandCapacity(objBase, count)

	return objBase
}

// Replace the characters in a substring of this StringBuilder object with characters in the specified String.
func stringBuilderReplace(params []any) any {
	// Get byteArray.
	objBase := params[0].(*object.Object)
	fld := objBase.FieldTable["value"]
	initBytes := fld.Fvalue.([]types.JavaByte)
	initLen := int64(len(initBytes))

	// Get start index, end index, and byte array to use as a replacment.
	start := params[1].(int64)
	end := params[2].(int64)
	repls := params[3].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)

	// Validate start and end.
	if start < 0 || start > initLen {
		errMsg := fmt.Sprintf("stringBuilderReplace: Start value (%d) < 0 or exceeds the byte array size (%d)", start, initLen)
		return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
	}
	if end < start {
		errMsg := fmt.Sprintf("stringBuilderReplace: End value (%d) < Start value (%d)", start, end)
		return ghelpers.GetGErrBlk(excNames.StringIndexOutOfBoundsException, errMsg)
	}
	if end > initLen {
		end = initLen
	}

	// If start = end, return object as-is.
	if start == end {
		return objBase
	}

	// Copy the left-most retained bytes to a new byte array.
	newArray := make([]types.JavaByte, start)
	if start > 0 {
		copy(newArray, initBytes[0:start])
	}

	// Append newArray with the replacement bytes.
	if len(repls) > 0 {
		newArray = append(newArray, repls...)
	}

	// Append newArray with the right-most retained bytes.
	newArray = append(newArray, initBytes[end:]...)
	newlen := int64(len(newArray))

	objBase.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: newArray}
	objBase.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: newlen}
	expandCapacity(objBase, newlen)

	return objBase
}

// Reverse the order of the byte array.
func stringBuilderReverse(params []any) any {
	// Get byteArray.
	objBase := params[0].(*object.Object)
	fld := objBase.FieldTable["value"]
	byteArray := fld.Fvalue.([]types.JavaByte)

	// Reverse the bytes in byteArray.
	for ii, jj := 0, len(byteArray)-1; ii < jj; ii, jj = ii+1, jj-1 {
		byteArray[ii], byteArray[jj] = byteArray[jj], byteArray[ii]
	}

	objBase.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: byteArray}
	return objBase
}

// Set the char parameter into the bytes into the StringBuilder
// at the given index.
func stringBuilderSetCharAt(params []any) any {
	obj := params[0].(*object.Object)
	fld := obj.FieldTable["value"]
	byteArray := fld.Fvalue.([]types.JavaByte)
	ix := params[1].(int64)
	ch := params[2].(int64)
	if ix < 0 || ix > int64(len(byteArray)) {
		errMsg := fmt.Sprintf("stringBuilderSetCharAt: Index value (%d) is illegal", ix)
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
	}
	byteArray[ix] = types.JavaByte(ch)
	obj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: byteArray}

	return nil
}

// Set the length of the character sequence.
func stringBuilderSetLength(params []any) any {
	obj := params[0].(*object.Object)
	fld := obj.FieldTable["value"]
	oldArray := fld.Fvalue.([]types.JavaByte)
	oldlen := int64(len(oldArray))
	newlen := params[1].(int64)
	newArray := make([]types.JavaByte, newlen)
	if newlen < 0 {
		errMsg := fmt.Sprintf("stringBuilderSetLength: Length value (%d) is negative", newlen)
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
	}
	if newlen == oldlen {
		return nil
	}
	if newlen > oldlen {
		copy(newArray, oldArray)
		for ix := oldlen; ix < newlen; ix++ {
			newArray[ix] = 0
		}
	} else { // truncation, newlen < oldlen
		copy(newArray, oldArray[:newlen])
	}
	obj.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: newArray}
	obj.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: newlen}
	expandCapacity(obj, newlen)

	return nil
}

// Convert the byte array of a StringBuilder object to a String object. Then, return it.
func stringBuilderToString(params []any) any {
	objBase := params[0].(*object.Object)
	byteArray := objBase.FieldTable["value"].Fvalue.([]types.JavaByte)
	goStr := object.GoStringFromJavaByteArray(byteArray)
	objOut := object.StringObjectFromGoString(goStr)
	return objOut
}

// Return the StringBuilder object capacity.
func stringBuilderCapacity(params []any) any {
	objBase := params[0].(*object.Object)
	return objBase.FieldTable["capacity"].Fvalue.(int64)
}

// Return the StringBuilder object length.
func stringBuilderLength(params []any) any {
	objBase := params[0].(*object.Object)
	return objBase.FieldTable["count"].Fvalue.(int64)
}

// Expand the capacity of a StringBuilder object.
func expandCapacity(obj *object.Object, count int64) {
	capField := obj.FieldTable["capacity"]
	capacity := capField.Fvalue.(int64)
	for count > capacity { // Expand capacity while count exceeds capacity.
		capacity = (capacity * 2) + 2
	}
	capField.Fvalue = capacity
	obj.FieldTable["capacity"] = capField
}
