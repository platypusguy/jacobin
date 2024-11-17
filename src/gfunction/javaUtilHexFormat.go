/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/object"
	"jacobin/statics"
	"jacobin/types"
)

// Implementation of some of the functions in in Java/lang/Class.

func Load_Util_HexFormat() {

	MethodSignatures["java/util/HexFormat.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  hexMapClinit,
		}

	MethodSignatures["java/util/HexFormat.fromHexDigits(Ljava/lang/CharSequence;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/HexFormat.fromHexDigits(Ljava/lang/CharSequence;II)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/HexFormat.fromHexDigitsToLong(Ljava/lang/CharSequence;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/HexFormat.fromHexDigitsToLong(Ljava/lang/CharSequence;II)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/HexFormat.parseHex(Ljava/lang/CharSequence;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/HexFormat.parseHex(Ljava/lang/CharSequence;II)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

}

// <clinit> for class HexFormat
func hexMapClinit(params []interface{}) interface{} {
	DIGITS := []uint8{
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 255, 255, 255, 255, 255, 255,
		255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	}
	sDigits := statics.Static{Type: "[B", Value: DIGITS}
	statics.AddStatic("java/util/HexFormat.DIGITS", sDigits)

	UPPERCASE_DIGITS := []uint8{
		'0', '1', '2', '3', '4', '5', '6', '7',
		'8', '9', 'A', 'B', 'C', 'D', 'E', 'F',
	}
	sUppercaseDigits := statics.Static{Type: "[B", Value: UPPERCASE_DIGITS}
	statics.AddStatic("java/util/HexFormat.UPPERCASE_DIGITS", sUppercaseDigits)

	LOWERCASE_DIGITS := []uint8{
		'0', '1', '2', '3', '4', '5', '6', '7',
		'8', '9', 'a', 'b', 'c', 'd', 'e', 'f',
	}
	sLowercaseDigits := statics.Static{Type: "[B", Value: LOWERCASE_DIGITS}
	statics.AddStatic("java/util/HexFormat.LOWERCASE_DIGITS", sLowercaseDigits)

	obj := mkHexFormatObject("", "", "", LOWERCASE_DIGITS)
	sHexFormat := statics.Static{Type: "java/util/HexFormat", Value: obj}
	statics.AddStatic("java/util/HexFormat.HEX_FORMAT", sHexFormat)

	sEmptyBytes := statics.Static{Type: "[B", Value: []byte{}}
	statics.AddStatic("java/util/HexFormat.EMPTY_BYTES", sEmptyBytes)

	return nil
}

// Make a new HexFormat object and return it to caller.
func mkHexFormatObject(delimiter, prefix, suffix string, digits []byte) *object.Object {
	var fld object.Field
	className := "java/util/HexFormat"
	obj := object.MakeEmptyObjectWithClassName(&className)

	fld = object.Field{Ftype: types.ByteArray, Fvalue: []byte(delimiter)}
	obj.FieldTable["delimiter"] = fld

	fld = object.Field{Ftype: types.ByteArray, Fvalue: []byte(prefix)}
	obj.FieldTable["prefix"] = fld

	fld = object.Field{Ftype: types.ByteArray, Fvalue: []byte(suffix)}
	obj.FieldTable["suffix"] = fld

	fld = object.Field{Ftype: types.ByteArray, Fvalue: digits}
	obj.FieldTable["digits"] = fld

	return obj
}
