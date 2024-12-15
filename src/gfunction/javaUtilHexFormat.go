/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/excNames"
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
			GFunction:  hfEquals,
		}

	MethodSignatures["java/util/HexFormat.formatHex([B)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hfFormatHexFromBytes,
		}

	MethodSignatures["java/util/HexFormat.formatHex([BII)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  hfFormatHexFromBytes,
		}

	MethodSignatures["java/util/HexFormat.fromHexDigits(Ljava/lang/CharSequence;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/HexFormat.fromHexDigit(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hfFromHexDigit,
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

	MethodSignatures["java/util/HexFormat.hashCode()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/HexFormat.of()Ljava/util/HexFormat;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  hfOf,
		}

	MethodSignatures["java/util/HexFormat.ofDelimiter(Ljava/lang/String;)Ljava/util/HexFormat;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hfOfDelimiter,
		}

	MethodSignatures["java/util/HexFormat.parseHex([CII)[B"] =
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

	MethodSignatures["java/util/HexFormat.prefix()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  hfPrefix,
		}

	MethodSignatures["java/util/HexFormat.suffix()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  hfSuffix,
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

	MethodSignatures["java/util/HexFormat.toHighHexDigit(I)C"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hfToHighHexDigit,
		}

	MethodSignatures["java/util/HexFormat.toLowHexDigit(I)C"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hfToLowHexDigit,
		}

	MethodSignatures["java/util/HexFormat.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  hfToString,
		}

	MethodSignatures["java/util/HexFormat.withDelimiter(Ljava/lang/String;)Ljava/util/HexFormat;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hfWithDelimiter,
		}

	MethodSignatures["java/util/HexFormat.withLowerCase()Ljava/util/HexFormat;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  hfWithLowerCase,
		}

	MethodSignatures["java/util/HexFormat.withPrefix(Ljava/lang/String;)Ljava/util/HexFormat;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hfWithPrefix,
		}

	MethodSignatures["java/util/HexFormat.withSuffix(Ljava/lang/String;)Ljava/util/HexFormat;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hfWithSuffix,
		}

	MethodSignatures["java/util/HexFormat.withUpperCase()Ljava/util/HexFormat;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  hfWithUpperCase,
		}

}

// <clinit> for class HexFormat
func hexFormatClinit(params []interface{}) interface{} {
	DIGITS := []types.JavaByte{
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, -1, -1, -1, -1, -1, -1,
		-1, 10, 11, 12, 13, 14, 15, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, 10, 11, 12, 13, 14, 15, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	}
	sDigits := statics.Static{Type: types.ByteArray, Value: DIGITS}
	statics.AddStatic("java/util/HexFormat.DIGITS", sDigits)

	UPPERCASE_DIGITS := []types.JavaByte{
		'0', '1', '2', '3', '4', '5', '6', '7',
		'8', '9', 'A', 'B', 'C', 'D', 'E', 'F',
	}
	sUppercaseDigits := statics.Static{Type: types.ByteArray, Value: UPPERCASE_DIGITS}
	statics.AddStatic("java/util/HexFormat.UPPERCASE_DIGITS", sUppercaseDigits)

	LOWERCASE_DIGITS := []types.JavaByte{
		'0', '1', '2', '3', '4', '5', '6', '7',
		'8', '9', 'a', 'b', 'c', 'd', 'e', 'f',
	}
	sLowercaseDigits := statics.Static{Type: types.ByteArray, Value: LOWERCASE_DIGITS}
	statics.AddStatic("java/util/HexFormat.LOWERCASE_DIGITS", sLowercaseDigits)

	obj := mkHexFormatObject([]types.JavaByte{}, []types.JavaByte{}, []types.JavaByte{}, LOWERCASE_DIGITS)
	sHexFormat := statics.Static{Type: "Ljava/util/HexFormat;", Value: obj}
	statics.AddStatic("java/util/HexFormat.HEX_FORMAT", sHexFormat)

	sEmptyBytes := statics.Static{Type: types.ByteArray, Value: []types.JavaByte{}}
	statics.AddStatic("java/util/HexFormat.EMPTY_BYTES", sEmptyBytes)

	sJavaLangAccess := statics.Static{Type: types.Int, Value: 42}
	statics.AddStatic("java/util/HexFormat.jla", sJavaLangAccess)

	return nil
}

// Make a new HexFormat object and return the object struct to caller.
func mkHexFormatObject(delimiter, prefix, suffix, digits []types.JavaByte) *object.Object {
	var fld object.Field
	className := "java/util/HexFormat"
	obj := object.MakeEmptyObjectWithClassName(&className)

	fld = object.Field{Ftype: types.ByteArray, Fvalue: delimiter}
	obj.FieldTable["delimiter"] = fld

	fld = object.Field{Ftype: types.ByteArray, Fvalue: prefix}
	obj.FieldTable["prefix"] = fld

	fld = object.Field{Ftype: types.ByteArray, Fvalue: suffix}
	obj.FieldTable["suffix"] = fld

	fld = object.Field{Ftype: types.ByteArray, Fvalue: digits}
	obj.FieldTable["digits"] = fld

	return obj
}

func hfToString(params []interface{}) interface{} {
	var delimiter, prefix, suffix, digits []types.JavaByte

	obj := params[0].(*object.Object)
	switch obj.FieldTable["delimiter"].Fvalue.(type) {
	case []types.JavaByte:
		delimiter = obj.FieldTable["delimiter"].Fvalue.([]types.JavaByte)
	case []byte:
		delimiter =
			object.JavaByteArrayFromGoByteArray(obj.FieldTable["delimiter"].Fvalue.([]byte))
	}

	switch obj.FieldTable["prefix"].Fvalue.(type) {
	case []types.JavaByte:

		prefix = obj.FieldTable["prefix"].Fvalue.([]types.JavaByte)
	case []byte:
		prefix = object.JavaByteArrayFromGoByteArray(obj.FieldTable["prefix"].Fvalue.([]byte))
	}

	switch obj.FieldTable["suffix"].Fvalue.(type) {
	case []types.JavaByte:
		suffix = obj.FieldTable["suffix"].Fvalue.([]types.JavaByte)
	case []byte:
		suffix =
			object.JavaByteArrayFromGoByteArray(obj.FieldTable["suffix"].Fvalue.([]byte))
	}

	switch obj.FieldTable["digits"].Fvalue.(type) {
	case []types.JavaByte:
		digits =
			obj.FieldTable["digits"].Fvalue.([]types.JavaByte)
	case []byte:
		digits =
			object.JavaByteArrayFromGoByteArray(obj.FieldTable["digits"].Fvalue.([]byte))
	}

	uppercase := (digits[15] == 'F')
	str := fmt.Sprintf("uppercase: %v, delimiter: \"%s\", prefix: \"%s\", suffix: \"%s\"",
		uppercase, object.GoStringFromJavaByteArray(delimiter),
		object.GoStringFromJavaByteArray(prefix),
		object.GoStringFromJavaByteArray(suffix))
	return object.StringObjectFromGoString(str)
}

func hfDelimiter(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	return object.StringObjectFromJavaByteArray(obj.FieldTable["delimiter"].Fvalue.([]types.JavaByte))
}

func hfPrefix(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	return object.StringObjectFromJavaByteArray(obj.FieldTable["prefix"].Fvalue.([]types.JavaByte))
}

func hfSuffix(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	return object.StringObjectFromJavaByteArray(obj.FieldTable["suffix"].Fvalue.([]types.JavaByte))
}

func hfByteToHexDigits(params []interface{}) interface{} {
	var str string
	obj := params[0].(*object.Object)
	primitive := params[1]
	input := primitive.(int64) % 256
	switch primitive.(type) {
	case int64:
		digits := obj.FieldTable["digits"].Fvalue.([]types.JavaByte)
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
		digits := obj.FieldTable["digits"].Fvalue.([]types.JavaByte)
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
		digits := obj.FieldTable["digits"].Fvalue.([]types.JavaByte)
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
		digits := obj.FieldTable["digits"].Fvalue.([]types.JavaByte)
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
		digits := obj.FieldTable["digits"].Fvalue.([]types.JavaByte)
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
// if the fromIndex and toIndex are given in params, use their values.
func hfFormatHexFromBytes(params []interface{}) interface{} {
	str := ""
	this := params[0].(*object.Object)
	digits := this.FieldTable["digits"].Fvalue.([]types.JavaByte)
	delimiter :=
		object.GoStringFromJavaByteArray(this.FieldTable["delimiter"].Fvalue.([]types.JavaByte))
	prefix :=
		object.GoStringFromJavaByteArray(this.FieldTable["prefix"].Fvalue.([]types.JavaByte))
	suffix :=
		object.GoStringFromJavaByteArray(this.FieldTable["suffix"].Fvalue.([]types.JavaByte))
	objBytes := params[1].(*object.Object)
	bytes := objBytes.FieldTable["value"].Fvalue.([]types.JavaByte)
	var fromIndex int
	var toIndex int
	if len(params) > 2 {
		fromIndex = int(params[2].(int64))
		if fromIndex < 0 || fromIndex > len(bytes) {
			errMsg := fmt.Sprintf("from index out of range: %d", fromIndex)
			return getGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
		}
		toIndex = int(params[3].(int64))
		if toIndex < 0 || toIndex > len(bytes) {
			errMsg := fmt.Sprintf("to index out of range: %d", fromIndex)
			return getGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
		}
		if toIndex <= fromIndex {
			errMsg := fmt.Sprintf("to index <= from index: %d", fromIndex)
			return getGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
		}
	} else {
		fromIndex = 0
		toIndex = len(bytes)
	}
	for ix := fromIndex; ix < toIndex; ix++ {
		if digits[15] == 'F' { // uppercase
			str += fmt.Sprintf("%s%02X%s%s", prefix, bytes[ix], suffix, delimiter)
		} else { // lowercase
			str += fmt.Sprintf("%s%02x%s%s", prefix, bytes[ix], suffix, delimiter)
		}
	}
	str = str[:len(str)-len(delimiter)]
	return object.StringObjectFromGoString(str)
}

func hfFromHexDigit(params []interface{}) interface{} {
	arg := params[0].(int64)
	if arg < 58 && arg > 47 { // range: '0' to '9'
		return arg - 48 // arg - '0'
	}
	if arg < 71 && arg > 64 { // range: 'A' to 'F'
		return arg - 55 // arg + 10 - 'A'
	}
	if arg < 103 && arg > 96 { // range: 'a' to 'f'
		return arg - 87 // arg + 10 - 'a'
	}
	errMsg := fmt.Sprintf("Out of range: %d", arg)
	return getGErrBlk(excNames.NumberFormatException, errMsg)

}

func hfOf(params []interface{}) interface{} {
	template := statics.GetStaticValue("java/util/HexFormat", "HEX_FORMAT").(*object.Object)
	obj := object.CloneObject(template)
	return obj
}

func hfOfDelimiter(params []interface{}) interface{} {
	delimiter := params[0].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	template := statics.GetStaticValue("java/util/HexFormat", "HEX_FORMAT").(*object.Object)
	obj := object.CloneObject(template)
	fld := obj.FieldTable["delimiter"]
	fld.Fvalue = delimiter
	obj.FieldTable["delimiter"] = fld

	return obj
}

func hfWithPrefix(params []interface{}) interface{} {
	obj1 := params[0].(*object.Object)
	prefix := params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	obj2 := object.CloneObject(obj1)
	fld := obj2.FieldTable["prefix"]
	fld.Fvalue = prefix
	obj2.FieldTable["prefix"] = fld
	return obj2
}

func hfWithSuffix(params []interface{}) interface{} {
	obj1 := params[0].(*object.Object)
	suffix := params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	obj2 := object.CloneObject(obj1)
	fld := obj2.FieldTable["suffix"]
	fld.Fvalue = suffix
	obj2.FieldTable["suffix"] = fld
	return obj2
}

func hfWithDelimiter(params []interface{}) interface{} {
	obj1 := params[0].(*object.Object)
	delimiter := params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	obj2 := object.CloneObject(obj1)
	fld := obj2.FieldTable["delimiter"]
	fld.Fvalue = delimiter
	obj2.FieldTable["delimiter"] = fld
	return obj2
}

func hfWithUpperCase(params []interface{}) interface{} {
	obj1 := params[0].(*object.Object)
	obj2 := object.CloneObject(obj1)
	fld := obj2.FieldTable["digits"]
	digits := statics.GetStaticValue("java/util/HexFormat", "UPPERCASE_DIGITS").([]types.JavaByte)
	fld.Fvalue = digits
	obj2.FieldTable["digits"] = fld
	return obj2
}

func hfWithLowerCase(params []interface{}) interface{} {
	obj1 := params[0].(*object.Object)
	obj2 := object.CloneObject(obj1)
	fld := obj2.FieldTable["digits"]
	digits := statics.GetStaticValue("java/util/HexFormat", "LOWERCASE_DIGITS").([]types.JavaByte)
	fld.Fvalue = digits
	obj2.FieldTable["digits"] = fld
	return obj2
}

func hfToHighHexDigit(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	digits := obj.FieldTable["digits"].Fvalue.([]types.JavaByte)
	arg := params[1].(int64) % 256
	return int64(digits[arg>>4])
}

func hfToLowHexDigit(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	digits := obj.FieldTable["digits"].Fvalue.([]types.JavaByte)
	arg := params[1].(int64) % 256
	return int64(digits[arg&0x0f])
}

// Helper function for hfEquals.
func hfEqualsHelper(this, that *object.Object, fieldName string) bool {
	f1, ok := this.FieldTable[fieldName].Fvalue.([]types.JavaByte)
	if !ok {
		return false
	}
	f2, ok := that.FieldTable[fieldName].Fvalue.([]types.JavaByte)
	if !ok {
		return false
	}
	if len(f1) != len(f2) {
		return false
	}

	if object.GoStringFromJavaByteArray(f1) != object.GoStringFromJavaByteArray(f2) {
		return false
	}

	return true
}

// Returns true if the other object is a HexFormat and the parameters uppercase, delimiter, prefix, and suffix are equal;
// otherwise false.
func hfEquals(params []interface{}) interface{} {

	this := params[0].(*object.Object)
	that := params[1].(*object.Object)
	if that.KlassName != this.KlassName {
		return types.JavaBoolFalse
	}
	ok := hfEqualsHelper(this, that, "digits")
	if !ok {
		return types.JavaBoolFalse
	}
	ok = hfEqualsHelper(this, that, "prefix")
	if !ok {
		return types.JavaBoolFalse
	}
	ok = hfEqualsHelper(this, that, "suffix")
	if !ok {
		return types.JavaBoolFalse
	}
	ok = hfEqualsHelper(this, that, "delimiter")
	if !ok {
		return types.JavaBoolFalse
	}

	return types.JavaBoolTrue

}
