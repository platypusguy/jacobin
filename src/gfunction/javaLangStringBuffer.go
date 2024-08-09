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
)

// Implementation of some of the functions in Java/lang/Class.

func Load_Lang_StringBuffer() {

	// === Instantiation ===

	MethodSignatures["java/lang/StringBuffer.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/StringBuffer.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringBufferInit,
		}

	MethodSignatures["java/lang/StringBuffer.<init>(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBufferInit,
		}

	MethodSignatures["java/lang/StringBuffer.<init>(Ljava/lang/CharSequence;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.<init>(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBufferInitString,
		}

	// === Methods ===

	MethodSignatures["java/lang/StringBuffer.append(Z)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppendBoolean,
		}

	MethodSignatures["java/lang/StringBuffer.append(C)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppendChar,
		}

	MethodSignatures["java/lang/StringBuffer.append([C)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	MethodSignatures["java/lang/StringBuffer.append([CII)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  stringBuilderAppend,
		}

	MethodSignatures["java/lang/StringBuffer.append(D)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderAppend,
		}

	MethodSignatures["java/lang/StringBuffer.append(F)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	MethodSignatures["java/lang/StringBuffer.append(I)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	MethodSignatures["java/lang/StringBuffer.append(J)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderAppend,
		}

	MethodSignatures["java/lang/StringBuffer.append(Ljava/lang/CharSequence;)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.append(Ljava/lang/CharSequence;II)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.append(Ljava/lang/Object;)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	MethodSignatures["java/lang/StringBuffer.append(Ljava/lang/String;)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	MethodSignatures["java/lang/StringBuffer.append(Ljava/lang/StringBuffer;)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	MethodSignatures["java/lang/StringBuffer.appendCodePoint(I)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.capacity()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringBuilderCapacity,
		}

	MethodSignatures["java/lang/StringBuffer.charAt(I)C"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderCharAt,
		}

	MethodSignatures["java/lang/StringBuffer.chars()Ljava/util/stream/IntStream;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.codePointAt(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.codePointBefore(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.codePointCount(II)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.codePoints()Ljava/util/stream/IntStream;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.compareTo(Ljava/lang/StringBuffer;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringCompareToCaseSensitive,
		}

	MethodSignatures["java/lang/StringBuffer.delete(II)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderDelete,
		}

	MethodSignatures["java/lang/StringBuffer.deleteCharAt(I)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderDelete,
		}

	MethodSignatures["java/lang/StringBuffer.ensureCapacity(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/StringBuffer.getChars(II[CI)V"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.indexOf(Ljava/lang/String;)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.indexOf(Ljava/lang/String;I)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.insert(IZ)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsertBoolean,
		}

	MethodSignatures["java/lang/StringBuffer.insert(IC)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsertChar,
		}

	MethodSignatures["java/lang/StringBuffer.insert(I[C)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	MethodSignatures["java/lang/StringBuffer.insert(I[CII)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  stringBuilderInsert,
		}

	MethodSignatures["java/lang/StringBuffer.insert(ID)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  stringBuilderInsert,
		}

	MethodSignatures["java/lang/StringBuffer.insert(IF)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	MethodSignatures["java/lang/StringBuffer.insert(II)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	MethodSignatures["java/lang/StringBuffer.insert(IJ)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  stringBuilderInsert,
		}

	MethodSignatures["java/lang/StringBuffer.insert(ILjava/lang/CharSequence;)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.insert(ILjava/lang/CharSequence;II)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.insert(ILjava/lang/Object;)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	MethodSignatures["java/lang/StringBuffer.insert(ILjava/lang/String;)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}
	MethodSignatures["java/lang/StringBuffer.isLatin1()Z"] = // internal member function, not in API
		GMeth{
			ParamSlots: 0,
			GFunction:  returnTrue,
		}

	MethodSignatures["java/lang/StringBuffer.lastIndexOf(Ljava/lang/String;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.lastIndexOf(Ljava/lang/String;I)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.length()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringBuilderLength,
		}

	MethodSignatures["java/lang/StringBuffer.offsetByCodePoints(II)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.replace(IILjava/lang/String;)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  stringBuilderReplace,
		}

	MethodSignatures["java/lang/StringBuffer.reverse()Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringBuilderReverse,
		}

	MethodSignatures["java/lang/StringBuffer.setCharAt(IC)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderSetCharAt,
		}

	MethodSignatures["java/lang/StringBuffer.setLength(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderSetLength,
		}

	MethodSignatures["java/lang/StringBuffer.subSequence(II)Ljava/lang/CharSequence;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.substring(I)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  substringToTheEnd, // javaLangString.go
		}

	MethodSignatures["java/lang/StringBuffer.substring(II)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  substringStartEnd, // javaLangString.go
		}

	MethodSignatures["java/lang/StringBuffer.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringBuilderToString,
		}

	MethodSignatures["java/lang/StringBuffer.trimToSize()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

}

var classStringBuffer = "java/lang/StringBuffer"

// Initialise StringBuffer with or without a capacity integer.
func stringBufferInit(params []any) any {
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

// Initialise StringBuffer with a String object.
func stringBufferInitString(params []any) any {
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
