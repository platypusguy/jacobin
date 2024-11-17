/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/object"
	"jacobin/statics"
	"jacobin/types"
)

// Implementation of some of the functions in in Java/lang/Class.

func Load_Util_HexFormat() {

	MethodSignatures["java/util/HexFormat.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  hexFormatClinit,
		}

	MethodSignatures["java/util/HexFormat.delimiter()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  hfDelimiter,
		}

	MethodSignatures["java/util/HexFormat.equals(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/HexFormat.formatHex([B)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hfFormatHexFromBytes,
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

	MethodSignatures["java/util/HexFormat.toHexDigits(B)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hfByteToHexDigits,
		}

	MethodSignatures["java/util/HexFormat.toHexDigits(C)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hfCharToHexDigits,
		}

	MethodSignatures["java/util/HexFormat.toHexDigits(I)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hfIntToHexDigits,
		}

	MethodSignatures["java/util/HexFormat.toHexDigits(J)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hfLongToHexDigits,
		}

	MethodSignatures["java/util/HexFormat.toHexDigits(JI)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  hfLongToHexDigits,
		}

	MethodSignatures["java/util/HexFormat.toHexDigits(S)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hfShortToHexDigits,
		}

	MethodSignatures["java/util/HexFormat.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  hfToString,
		}

}

// <clinit> for class HexFormat
func hexFormatClinit(params []interface{}) interface{} {
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

	obj := mkHexFormatObject("", "", "", LOWERCASE_DIGITS, false)
	sHexFormat := statics.Static{Type: "java/util/HexFormat", Value: obj}
	statics.AddStatic("java/util/HexFormat.HEX_FORMAT", sHexFormat)

	sEmptyBytes := statics.Static{Type: "[B", Value: []byte{}}
	statics.AddStatic("java/util/HexFormat.EMPTY_BYTES", sEmptyBytes)

	return nil
}

// Make a new HexFormat object and return it to caller.
func mkHexFormatObject(delimiter, prefix, suffix string, digits []byte, flagUpperCase bool) *object.Object {
	var fld object.Field
	className := "java/util/HexFormat"
	ftypeString := "Ljava/lang/String;"
	obj := object.MakeEmptyObjectWithClassName(&className)

	fld = object.Field{Ftype: ftypeString, Fvalue: object.StringObjectFromGoString(delimiter)}
	obj.FieldTable["delimiter"] = fld

	fld = object.Field{Ftype: ftypeString, Fvalue: object.StringObjectFromGoString(prefix)}
	obj.FieldTable["prefix"] = fld

	fld = object.Field{Ftype: ftypeString, Fvalue: object.StringObjectFromGoString(suffix)}
	obj.FieldTable["suffix"] = fld

	fld = object.Field{Ftype: types.ByteArray, Fvalue: digits}
	obj.FieldTable["digits"] = fld

	return obj
}

func hfToString(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	delimiter := obj.FieldTable["delimiter"].Fvalue.([]byte)
	prefix := obj.FieldTable["prefix"].Fvalue.([]byte)
	suffix := obj.FieldTable["suffix"].Fvalue.([]byte)
	digits := obj.FieldTable["digits"].Fvalue.([]byte)
	uppercase := (digits[15] == 'F')
	str := fmt.Sprintf("uppercase: %v, delimiter: \"%s\", prefix: \"%s\", suffix: \"%s\"", uppercase, delimiter, prefix, suffix)
	return object.StringObjectFromGoString(str)
}

func hfDelimiter(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	return object.StringObjectFromByteArray(obj.FieldTable["delimiter"].Fvalue.([]byte))
}

func hfByteToHexDigits(params []interface{}) interface{} {
	var str string
	obj := params[0].(*object.Object)
	primitive := params[1]
	input := primitive.(int64) % 256
	switch primitive.(type) {
	case int64:
		digits := obj.FieldTable["digits"].Fvalue.([]byte)
		if digits[15] == 'F' { // uppercase
			str = fmt.Sprintf("%02X", input)
		} else { // lowercase
			str = fmt.Sprintf("%02x", input)
		}
	default:
		return trapFunction(params)
	}
	return object.StringObjectFromGoString(str)
}

func hfCharToHexDigits(params []interface{}) interface{} {
	var str string
	obj := params[0].(*object.Object)
	primitive := params[1]
	input := primitive.(int64) % 256
	switch primitive.(type) {
	case int64:
		digits := obj.FieldTable["digits"].Fvalue.([]byte)
		if digits[15] == 'F' { // uppercase
			str = fmt.Sprintf("%04X", input)
		} else { // lowercase
			str = fmt.Sprintf("%04x", input)
		}
	default:
		return trapFunction(params)
	}
	return object.StringObjectFromGoString(str)
}

func hfIntToHexDigits(params []interface{}) interface{} {
	var str string
	obj := params[0].(*object.Object)
	primitive := params[1]
	input := primitive.(int64)
	switch primitive.(type) {
	case int64:
		digits := obj.FieldTable["digits"].Fvalue.([]byte)
		if digits[15] == 'F' { // uppercase
			str = fmt.Sprintf("%08X", input)
		} else { // lowercase
			str = fmt.Sprintf("%08x", input)
		}
	default:
		return trapFunction(params)
	}
	return object.StringObjectFromGoString(str)
}

func hfLongToHexDigits(params []interface{}) interface{} {
	var str string
	obj := params[0].(*object.Object)
	primitive := params[1]
	input := primitive.(int64)
	switch primitive.(type) {
	case int64:
		digits := obj.FieldTable["digits"].Fvalue.([]byte)
		if digits[15] == 'F' { // uppercase
			str = fmt.Sprintf("%016X", input)
		} else { // lowercase
			str = fmt.Sprintf("%016x", input)
		}
		if len(params) > 2 {
			outlen := int(params[2].(int64))
			str = str[len(str)-outlen:]
		}
	default:
		return trapFunction(params)
	}
	return object.StringObjectFromGoString(str)
}

func hfShortToHexDigits(params []interface{}) interface{} {
	var str string
	obj := params[0].(*object.Object)
	primitive := params[1]
	input := primitive.(int64)
	switch primitive.(type) {
	case int64:
		digits := obj.FieldTable["digits"].Fvalue.([]byte)
		if digits[15] == 'F' { // uppercase
			str = fmt.Sprintf("%04X", input)
		} else { // lowercase
			str = fmt.Sprintf("%04x", input)
		}
	default:
		return trapFunction(params)
	}
	return object.StringObjectFromGoString(str)
}

// Format a hex string from a byte slice.
func hfFormatHexFromBytes(params []interface{}) interface{} {
	str := ""
	this := params[0].(*object.Object)
	digits := this.FieldTable["digits"].Fvalue.([]byte)
	delimiter := string(this.FieldTable["delimiter"].Fvalue.([]byte))
	objBytes := params[1].(*object.Object)
	bytes := objBytes.FieldTable["value"].Fvalue.([]byte)
	for ix := 0; ix < len(bytes); ix++ {
		if digits[15] == 'F' { // uppercase
			str += fmt.Sprintf("%02X%s", bytes[ix], delimiter)
		} else { // lowercase
			str += fmt.Sprintf("%02x%s", bytes[ix], delimiter)
		}
	}
	str = str[:len(str)-1]
	return object.StringObjectFromGoString(str)
}
