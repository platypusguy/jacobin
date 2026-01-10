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
)

// Implementation of some of the functions in Java/lang/Class.

func Load_Lang_StringBuffer() {

	// === Instantiation ===

	ghelpers.MethodSignatures["java/lang/StringBuffer.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringBufferInit,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.<init>(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBufferInit,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.<init>(Ljava/lang/CharSequence;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBufferInitString,
		}

	// === Methods ===

	ghelpers.MethodSignatures["java/lang/StringBuffer.append(Z)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppendBoolean,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.append(C)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppendChar,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.append([C)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.append([CII)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  stringBuilderAppend,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.append(D)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.append(F)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.append(I)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.append(J)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.append(Ljava/lang/CharSequence;)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.append(Ljava/lang/CharSequence;II)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.append(Ljava/lang/Object;)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.append(Ljava/lang/String;)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.append(Ljava/lang/StringBuffer;)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderAppend,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.appendCodePoint(I)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.capacity()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringBuilderCapacity,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.charAt(I)C"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderCharAt,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.chars()Ljava/util/stream/IntStream;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.codePointAt(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.codePointBefore(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.codePointCount(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.codePoints()Ljava/util/stream/IntStream;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.compareTo(Ljava/lang/StringBuffer;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringCompareToCaseSensitive,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.delete(II)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderDelete,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.deleteCharAt(I)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderDelete,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.ensureCapacity(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.getChars(II[CI)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.indexOf(Ljava/lang/String;)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.indexOf(Ljava/lang/String;I)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.insert(IZ)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsertBoolean,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.insert(IC)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsertChar,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.insert(I[C)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.insert(I[CII)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  stringBuilderInsert,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.insert(ID)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.insert(IF)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.insert(II)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.insert(IJ)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.insert(ILjava/lang/CharSequence;)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.insert(ILjava/lang/CharSequence;II)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.insert(ILjava/lang/Object;)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.insert(ILjava/lang/String;)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderInsert,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.isLatin1()Z"] = // internal member function, not in API
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnTrue,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.lastIndexOf(Ljava/lang/String;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.lastIndexOf(Ljava/lang/String;I)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.length()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringBuilderLength,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.offsetByCodePoints(II)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.replace(IILjava/lang/String;)Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  stringBuilderReplace,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.reverse()Ljava/lang/StringBuffer;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringBuilderReverse,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.setCharAt(IC)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringBuilderSetCharAt,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.setLength(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringBuilderSetLength,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.subSequence(II)Ljava/lang/CharSequence;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.substring(I)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  substringToTheEnd, // javaLangString.go
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.substring(II)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  substringStartEnd, // javaLangString.go
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringBuilderToString,
		}

	ghelpers.MethodSignatures["java/lang/StringBuffer.trimToSize()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
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

// Initialize StringBuffer with a String object.
func stringBufferInitString(params []any) any {
	// Get File object and initialise the field map.
	obj := params[0].(*object.Object)
	obj.FieldTable = make(map[string]object.Field)

	var byteArray []types.JavaByte
	var ok bool
	switch params[1].(type) {
	case *object.Object: // String
		byteArray, ok = params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
		if !ok {
			errMsg := "StringBufferInitString: value field missing in <init> object or the field is not a byte array"
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	default:
		errMsg := fmt.Sprintf("StringBufferInitString: Parameter type (%T) is illegal", params[1])
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
