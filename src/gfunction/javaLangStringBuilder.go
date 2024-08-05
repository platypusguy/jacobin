/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/types"
	"strconv"
)

// Implementation of some of the functions in Java/lang/Class.

func Load_Lang_StringBuilder() {

	// === Instantiation ===

	MethodSignatures["java/lang/StringBuilder.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/StringBuilder.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringBuilderInit,
		}

	MethodSignatures["java/lang/StringBuilder.<init>(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderInit,
		}

	MethodSignatures["java/lang/StringBuilder.<init>(Ljava/lang/CharSequence;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.<init>(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderInitString,
		}

	// === Methods ===

	MethodSignatures["java/lang/StringBuilder.append(Z)Ljava/lang/StringBuilder;"] = // append boolean
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppendBoolean,
		}

	MethodSignatures["java/lang/StringBuilder.append(C)Ljava/lang/StringBuilder;"] = // append char
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppendChar,
		}

	MethodSignatures["java/lang/StringBuilder.append([C)Ljava/lang/StringBuilder;"] = // append char array
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	MethodSignatures["java/lang/StringBuilder.append([CII)Ljava/lang/StringBuilder;"] = // append subset of char array
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.append(D)Ljava/lang/StringBuilder;"] = // append double
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderAppend,
		}

	MethodSignatures["java/lang/StringBuilder.append(F)Ljava/lang/StringBuilder;"] = // append float
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	MethodSignatures["java/lang/StringBuilder.append(I)Ljava/lang/StringBuilder;"] = // append integer
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	MethodSignatures["java/lang/StringBuilder.append(J)Ljava/lang/StringBuilder;"] = // append long
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderAppend,
		}

	MethodSignatures["java/lang/StringBuilder.append(Ljava/lang/CharSequence;)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.append(Ljava/lang/CharSequence;II)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.append(Ljava/lang/Object;)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	MethodSignatures["java/lang/StringBuilder.append(Ljava/lang/String;)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	MethodSignatures["java/lang/StringBuilder.append(Ljava/lang/StringBuffer;)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	MethodSignatures["java/lang/StringBuilder.appendCodePoint(I)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.capacity()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringBuilderCapacity,
		}

	MethodSignatures["java/lang/StringBuilder.charAt(I)C"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderCharAt,
		}

	MethodSignatures["java/lang/StringBuilder.chars()Ljava/util/stream/IntStream;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.codePointAt(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.codePointBefore(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.codePointCount(II)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.codePoints()Ljava/util/stream/IntStream;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.compareTo(Ljava/lang/StringBuilder;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.delete(II)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.deleteCharAt(I)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.ensureCapacity(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/StringBuilder.getChars(II[CI)V"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.insert(IZ)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsertBoolean,
		}

	MethodSignatures["java/lang/StringBuilder.insert(IC)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	MethodSignatures["java/lang/StringBuilder.insert(I[C)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	MethodSignatures["java/lang/StringBuilder.insert(I[CII)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.insert(ID)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	MethodSignatures["java/lang/StringBuilder.insert(IF)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	MethodSignatures["java/lang/StringBuilder.insert(II)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	MethodSignatures["java/lang/StringBuilder.insert(IJ)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	MethodSignatures["java/lang/StringBuilder.insert(ILjava/lang/CharSequence;)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.insert(ILjava/lang/CharSequence;II)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.insert(ILjava/lang/Object;)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	MethodSignatures["java/lang/StringBuilder.insert(ILjava/lang/String;)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}
	MethodSignatures["java/lang/StringBuilder.isLatin1()Z"] = // internal member function, not in API
		GMeth{
			ParamSlots: 0,
			GFunction:  returnTrue,
		}

	MethodSignatures["java/lang/StringBuilder.lastIndexOf(Ljava/lang/String;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.lastIndexOf(Ljava/lang/String;I)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.length()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringBuilderLength,
		}

	MethodSignatures["java/lang/StringBuilder.offsetByCodePoints(II)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.replace(IILjava/lang/String;)Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.reverse()Ljava/lang/StringBuilder;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.setCharAt(IC)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.setLength(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuilder.subSequence(II)Ljava/lang/CharSequence;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	// Return a substring starting at the given index of the byte array.
	MethodSignatures["java/lang/StringBuilder.substring(I)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  substringToTheEnd, // javaLangString.go
		}

	// Return a substring starting at the given index of the byte array of the given length.
	MethodSignatures["java/lang/StringBuilder.substring(II)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  substringStartEnd, // javaLangString.go
		}

	MethodSignatures["java/lang/StringBuilder.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringBuilderToString,
		}

	MethodSignatures["java/lang/StringBuilder.trimToSize()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
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
	fld = object.Field{Ftype: types.ByteArray, Fvalue: make([]byte, 0)}
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

	var byteArray []byte
	var ok bool
	switch params[1].(type) {
	case *object.Object: // String
		byteArray, ok = params[1].(*object.Object).FieldTable["value"].Fvalue.([]byte)
		if !ok {
			errMsg := "Value field missing in <init> object argument or the field is not a byte array"
			return getGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	default:
		errMsg := fmt.Sprintf("Parameter type (%T) is illegal", params[1])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
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

// Append the second parameter to the bytes in the StringBuilder that is
// passed in the objectRef parameter (the first param).
func stringBuilderAppend(params []any) any {
	// Get base object and its value field, byteArray.
	objBase := params[0].(*object.Object)
	byteArray := objBase.FieldTable["value"].Fvalue.([]byte)

	// Initialise the output object.
	objOut := object.MakeEmptyObjectWithClassName(&classStringBuilder)
	objOut.FieldTable = objBase.FieldTable

	var parmArray []byte
	switch params[1].(type) {
	case *object.Object: // char array, String, StringBuffer, or StringBuilder
		fvalue := params[1].(*object.Object).FieldTable["value"].Fvalue
		switch fvalue.(type) {
		case []byte: // String, StringBuffer, or StringBuilder
			parmArray = fvalue.([]byte)
		case []int64: // char array
			for _, elem := range fvalue.([]int64) {
				parmArray = append(parmArray, byte(elem))
			}
		default:
			errMsg := fmt.Sprintf("Object value field value type (%T) is not a byte array nor a char array", params[1])
			return getGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	case int64: // integer, long
		str := fmt.Sprintf("%d", params[1].(int64))
		parmArray = []byte(str)
	case float64: // float, double
		ff := params[1].(float64)
		str := strconv.FormatFloat(ff, 'f', -1, 64)
		parmArray = []byte(str)
	default:
		errMsg := fmt.Sprintf("Parameter type (%T) is illegal", params[1])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Append parmArray to the byteArray.
	// Set the byte count.
	byteArray = append(byteArray, parmArray...)
	count := int64(len(byteArray))
	objOut.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: byteArray}
	objOut.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: count}
	expandCapacity(objOut, count)

	return objOut
}

// Append the second parameter (boolean) to the bytes in the StringBuilder that is
// passed in the objectRef parameter (the first param).
func stringBuilderAppendBoolean(params []any) any {
	// Get base object and its value field, byteArray.
	objBase := params[0].(*object.Object)
	byteArray := objBase.FieldTable["value"].Fvalue.([]byte)

	// Initialise the output object.
	objOut := object.MakeEmptyObjectWithClassName(&classStringBuilder)
	objOut.FieldTable = objBase.FieldTable

	var parmArray []byte
	switch params[1].(type) {
	case int64: // boolean
		var str string
		if params[1].(int64) == types.JavaBoolTrue {
			str = "true"
		} else {
			str = "false"
		}
		parmArray = []byte(str)
	default:
		errMsg := fmt.Sprintf("Parameter type (%T) is illegal", params[1])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Append parmArray to the byteArray.
	// Set the byte count.
	byteArray = append(byteArray, parmArray...)
	count := int64(len(byteArray))
	objOut.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: byteArray}
	objOut.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: count}
	expandCapacity(objOut, count)

	return objOut
}

// Append the second parameter (char) to the bytes in the StringBuilder that is
// passed in the objectRef parameter (the first param).
func stringBuilderAppendChar(params []any) any {
	// Get base object and its value field, byteArray.
	objBase := params[0].(*object.Object)
	byteArray := objBase.FieldTable["value"].Fvalue.([]byte)

	// Initialise the output object.
	objOut := object.MakeEmptyObjectWithClassName(&classStringBuilder)
	objOut.FieldTable = objBase.FieldTable

	var parmArray = make([]byte, 1)
	switch params[1].(type) {
	case int64: // char
		bb := byte(params[1].(int64))
		parmArray[0] = bb
	default:
		errMsg := fmt.Sprintf("Parameter type (%T) is illegal", params[1])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Append parmArray to the byteArray.
	// Set the byte count.
	byteArray = append(byteArray, parmArray...)
	count := int64(len(byteArray))
	objOut.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: byteArray}
	objOut.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: count}
	expandCapacity(objOut, count)

	return objOut
}

// Extract a character at the given index.
func stringBuilderCharAt(params []any) any {
	obj := params[0].(*object.Object)
	ix := params[1].(int64)
	bytes := obj.FieldTable["value"].Fvalue.([]byte)
	if ix >= int64(len(bytes)) {
		errMsg := fmt.Sprintf("Index value (%d) exceeds the byte array size (%d)", ix, len(bytes))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	return int64(bytes[ix])
}

// Insert the second parameter to the bytes into the StringBuilder
// at the given index.
func stringBuilderInsert(params []any) any {
	// Get base object and its value field, byteArray.
	objBase := params[0].(*object.Object)
	byteArray := objBase.FieldTable["value"].Fvalue.([]byte)

	// Get the index value.
	ix := params[1].(int64)
	if ix >= int64(len(byteArray)) {
		errMsg := fmt.Sprintf("Index value (%d) exceeds the byte array size (%d)", ix, len(byteArray))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Initialise the output object.
	objOut := object.MakeEmptyObjectWithClassName(&classStringBuilder)
	objOut.FieldTable = objBase.FieldTable

	var parmArray []byte
	switch params[1].(type) {
	case *object.Object: // char array, String, StringBuffer, or StringBuilder
		fvalue := params[1].(*object.Object).FieldTable["value"].Fvalue
		switch fvalue.(type) {
		case []byte: // String, StringBuffer, or StringBuilder
			parmArray = fvalue.([]byte)
		case []int64: // char array
			for _, elem := range fvalue.([]int64) {
				parmArray = append(parmArray, byte(elem))
			}
		default:
			errMsg := fmt.Sprintf("Object value field value type (%T) is not a byte array nor a char array", params[1])
			return getGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	case int64: // integer, long
		str := fmt.Sprintf("%d", params[1].(int64))
		parmArray = []byte(str)
	case float64: // float, double
		ff := params[1].(float64)
		str := strconv.FormatFloat(ff, 'f', -1, 64)
		parmArray = []byte(str)
	default:
		errMsg := fmt.Sprintf("Parameter type (%T) is illegal", params[1])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Append parmArray to the byteArray.
	// Set the byte count.
	newArray := byteArray[:(ix - 1)]
	newArray = append(newArray, parmArray...)
	newArray = append(newArray, byteArray[ix:]...)
	count := int64(len(newArray))

	objOut.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: newArray}
	objOut.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: count}
	expandCapacity(objOut, count)

	return objOut
}

// Insert the boolean parameter into the bytes into the StringBuilder
// at the given index.
func stringBuilderInsertBoolean(params []any) any {
	// Get base object and its value field, byteArray.
	objBase := params[0].(*object.Object)
	byteArray := objBase.FieldTable["value"].Fvalue.([]byte)

	// Get the index value.
	ix := params[1].(int64)
	if ix >= int64(len(byteArray)) {
		errMsg := fmt.Sprintf("Index value (%d) exceeds the byte array size (%d)", ix, len(byteArray))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Initialise the output object.
	objOut := object.MakeEmptyObjectWithClassName(&classStringBuilder)
	objOut.FieldTable = objBase.FieldTable

	var parmArray []byte
	switch params[1].(type) {
	case int64: // boolean
		var str string
		if params[1].(int64) == types.JavaBoolTrue {
			str = "true"
		} else {
			str = "false"
		}
		parmArray = []byte(str)
	default:
		errMsg := fmt.Sprintf("Parameter type (%T) is illegal", params[1])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Append parmArray to the byteArray.
	// Set the byte count.
	newArray := byteArray[:(ix - 1)]
	newArray = append(newArray, parmArray...)
	newArray = append(newArray, byteArray[ix:]...)
	count := int64(len(newArray))

	objOut.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: newArray}
	objOut.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: count}
	expandCapacity(objOut, count)

	return objOut
}

// Insert the char parameter into the bytes into the StringBuilder
// at the given index.
func stringBuilderInsertChar(params []any) any {
	// Get base object and its value field, byteArray.
	objBase := params[0].(*object.Object)
	byteArray := objBase.FieldTable["value"].Fvalue.([]byte)

	// Get the index value.
	ix := params[1].(int64)
	if ix >= int64(len(byteArray)) {
		errMsg := fmt.Sprintf("Index value (%d) exceeds the byte array size (%d)", ix, len(byteArray))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Initialise the output object.
	objOut := object.MakeEmptyObjectWithClassName(&classStringBuilder)
	objOut.FieldTable = objBase.FieldTable

	var parmArray = make([]byte, 1)
	switch params[1].(type) {
	case int64: // char
		bb := byte(params[1].(int64))
		parmArray[0] = bb
	default:
		errMsg := fmt.Sprintf("Parameter type (%T) is illegal", params[1])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Append parmArray to the byteArray.
	// Set the byte count.
	newArray := byteArray[:(ix - 1)]
	newArray = append(newArray, parmArray...)
	newArray = append(newArray, byteArray[ix:]...)
	count := int64(len(newArray))

	objOut.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: newArray}
	objOut.FieldTable["count"] = object.Field{Ftype: types.Int, Fvalue: count}
	expandCapacity(objOut, count)

	return objOut
}

// Convert the byte array of a StringBuilder object to a String object. Then, return it.
func stringBuilderToString(params []any) any {
	objBase := params[0].(*object.Object)
	byteArray := objBase.FieldTable["value"].Fvalue.([]byte)
	objOut := object.StringObjectFromGoString(string(byteArray))
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
