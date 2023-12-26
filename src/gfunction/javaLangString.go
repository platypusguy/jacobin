/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/exceptions"
	"jacobin/object"
	"jacobin/types"
	"strconv"
	"strings"
)

// We don't run String's static initializer block because the initialization
// is already handled in String creation

func Load_Lang_String() map[string]GMeth {

	// === OBJECT INSTANTIATION ===

	// String instantiation without parameters i.e. String string = new String();
	// need to replace eventually by enabling the Java initializer to run
	MethodSignatures["java/lang/String.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringClinit,
		}

	// String(byte[] bytes) - instantiate a String from a byte array
	MethodSignatures["java/lang/String.<init>([B)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  newStringFromBytes,
		}

	// String(byte[] ascii, int hibyte) *** DEPRECATED
	MethodSignatures["java/lang/String.<init>([BI)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  noSupportYetInString,
		}

	// String(byte[] bytes, int offset, int length)	- instantiate a String from a byte array SUBSET
	MethodSignatures["java/lang/String.<init>([BII)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  newSubstringFromBytes,
		}

	// String(byte[] ascii, int hibyte, int offset, int count) *** DEPRECATED
	MethodSignatures["java/lang/String.<init>([BIII)V"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  noSupportYetInString,
		}

	// String(byte[] bytes, int offset, int length, String charsetName) *********** CHARSET
	MethodSignatures["java/lang/String.<init>([BIILjava/lang/String;)V"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  noSupportYetInString,
		}

	// String(byte[] bytes, int offset, int length, Charset charset) ************** CHARSET
	MethodSignatures["java/lang/String.<init>([BIILjava/nio/charset/Charset;)V"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  noSupportYetInString,
		}

	// String(byte[] bytes, String charsetName) *********************************** CHARSET
	MethodSignatures["java/lang/String.<init>([BLjava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  noSupportYetInString,
		}

	// String(byte[] bytes, Charset charset) ************************************** CHARSET
	MethodSignatures["java/lang/String.<init>([BLjava/nio/charset/Charset;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  noSupportYetInString,
		}

	// String(char[] value) *************************************************** works fine in Java

	// String(char[] value, int offset, int count) ***************************- works fine in Java

	// String(int[] codePoints, int offset, int count) ************************ CODEPOINTS
	MethodSignatures["java/lang/String.<init>([III)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  noSupportYetInString,
		}

	// String(String original) - works fine in Java

	// String(StringBuffer buffer) ********************************************* StringBuffer
	MethodSignatures["java/lang/String.<init>(Ljava/lang/StringBuffer;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  noSupportYetInString,
		}

	// String(StringBuilder builder) ******************************************* StringBuilder
	MethodSignatures["java/lang/String.<init>(Ljava/lang/StringBuilder;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  noSupportYetInString,
		}

	// === METHOD FUNCTIONS ===

	// get the bytes from a string
	MethodSignatures["java/lang/String.getBytes()[B"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  getBytesVoid,
		}

	// get the bytes from a string, given the Charset string name ************************ CHARSET
	MethodSignatures["java/lang/String.getBytes(Ljava/lang/String;)[B"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  noSupportYetInString,
		}

	// get the bytes from a string, given the specified Charset object ******************* CHARSET
	MethodSignatures["java/lang/String.getBytes(Ljava/nio/charset/Charset;)[B"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  noSupportYetInString,
		}

	// Return a formatted string using the reference object string as the format string
	// and the supplied arguments as input object arguments.
	// E.g. String string = String.format("%s %i", "ABC", 42);
	MethodSignatures["java/lang/String.format(Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  sprintf,
		}

	// This method is equivalent to String.format(this, args).
	MethodSignatures["java/lang/String.formatted([Ljava/lang/Object;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  sprintf,
		}

	// Return a formatted string using the specified locale, format string, and arguments.
	MethodSignatures["java/lang/String.format(Ljava/util/Locale;Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  noSupportYetInString,
		}

	// Return the length of a String..
	MethodSignatures["java/lang/String.length()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringLength,
		}

	// Return a string in all lower case, using the reference object string as input.
	MethodSignatures["java/lang/String.toLowerCase()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  toLowerCase,
		}

	// Return a string in all lower case, using the reference object string as input.
	MethodSignatures["java/lang/String.toUpperCase()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  toUpperCase,
		}

	// Return a string representing a boolean value.
	MethodSignatures["java/lang/String.valueOf(Z)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  valueOfBoolean,
		}

	// Return a string representing a char value.
	MethodSignatures["java/lang/String.valueOf(C)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  valueOfChar,
		}

	// Return a string representing a char array.
	MethodSignatures["java/lang/String.valueOf([C)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  valueOfCharArray,
		}

	// Return a string representing a char subarray.
	MethodSignatures["java/lang/String.valueOf([CII)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  valueOfCharSubarray,
		}

	// Return a string representing a double value.
	MethodSignatures["java/lang/String.valueOf(D)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  valueOfDouble,
		}

	// Return a string representing a float value.
	MethodSignatures["java/lang/String.valueOf(F)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  valueOfFloat,
		}

	// Return a string representing an int value.
	MethodSignatures["java/lang/String.valueOf(I)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  valueOfInt,
		}

	// Return a string representing an int value.
	MethodSignatures["java/lang/String.valueOf(J)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  valueOfLong,
		}

	// Return a string representing the value of an Object.
	MethodSignatures["java/lang/String.valueOf(Ljava/lang/Object;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  valueOfObject,
		}

	// Compare 2 strings lexicographically, case-sensitive (upper/lower).
	// The return value is a negative integer, zero, or a positive integer
	// as the String argument is greater than, equal to, or less than this String,
	// case-sensitive.
	MethodSignatures["java/lang/String.compareTo(Ljava/lang/String;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  compareToCaseSensitive,
		}

	// Compare 2 strings lexicographically, ignoring case (upper/lower).
	// The return value is a negative integer, zero, or a positive integer
	// as the String argument is greater than, equal to, or less than this String,
	// ignoring case considerations.
	MethodSignatures["java/lang/String.compareToIgnoreCase(Ljava/lang/String;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  compareToIgnoreCase,
		}

	MethodSignatures["java/lang/String.concat(Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringConcat,
		}

	return MethodSignatures

}

func stringClinit([]interface{}) interface{} {
	klass := classloader.MethAreaFetch("java/lang/String")
	if klass == nil {
		errMsg := "In stringClinit, expected java/lang/String to be in the MethodArea, but it was not"
		exceptions.Throw(exceptions.VirtualMachineError, errMsg)
	}
	klass.Data.ClInit = types.ClInitRun // just mark that String.<clinit>() has been run
	return nil
}

// No support YET for references to Charset objects nor for Unicode code point arrays
func noSupportYetInString([]interface{}) interface{} {
	errMsg := "No support yet for user-specified character sets and Unicode code point arrays"
	exceptions.Throw(exceptions.UnsupportedEncodingException, errMsg)
	return nil
}

// Given a Go interface parameter from caller, compute the associated Go string.
func getGoString(param0 interface{}) string {
	var bptr *[]uint8
	switch param0.(type) {
	case *[]uint8:
		bptr = param0.(*[]uint8)
	case *object.Object:
		parmObj := param0.(*object.Object)
		bptr = parmObj.Fields[0].Fvalue.(*[]byte)
	default:
		errMsg := fmt.Sprintf("In getGoString, unexpected param[0] type = %T", param0)
		exceptions.Throw(exceptions.VirtualMachineError, errMsg)
		bptr = nil
	}
	return string(*bptr)
}

// Construct a compact string object (usable by Java) from a Go byte array.
func newStringFromBytes(params []interface{}) interface{} {
	klass := classloader.MethAreaFetch("java/lang/String")
	if klass == nil {
		errMsg := "In newStringFromBytes, expected java/lang/String to be in the MethodArea, but it was not"
		exceptions.Throw(exceptions.VirtualMachineError, errMsg)
	}
	klass.Data.ClInit = types.ClInitRun // just mark that String.<clinit>() has been run

	// Fetch a pointer to the raw slice of bytes from params[0].
	// Convert the raw slice of bytes to a Go string.
	wholeString := getGoString(params[1])

	// Convert the Go string to a compact string object, usable by Java. Return to caller.
	obj := object.CreateCompactStringFromGoString(&wholeString)
	return obj

}

// Construct a compact string object (usable by Java) from a Go byte array.
func newSubstringFromBytes(params []interface{}) interface{} {
	klass := classloader.MethAreaFetch("java/lang/String")
	if klass == nil {
		errMsg := "In newStringFromBytes, expected java/lang/String to be in the MethodArea, but it was not"
		exceptions.Throw(exceptions.VirtualMachineError, errMsg)
	}
	klass.Data.ClInit = types.ClInitRun // just mark that String.<clinit>() has been run

	// Fetch a pointer to the raw slice of bytes from params[0].
	// Convert the raw slice of bytes to a Go string.
	wholeString := getGoString(params[1])

	// Get substring offset and length
	ssOffset := params[2].(int64)
	ssLength := params[3].(int64)

	// Validate boundaries.
	wholeLength := int64(len(wholeString))
	if wholeLength < 1 || ssOffset < 0 || ssLength < 1 || ssOffset > (wholeLength-1) || (ssOffset+ssLength) > wholeLength {
		errMsg := "In newSubstringFromBytes, either: nil input byte array, invalid substring offset, or invalid substring length"
		exceptions.Throw(exceptions.StringIndexOutOfBoundsException, errMsg)
	}

	// Compute substring.
	ss := wholeString[ssOffset : ssOffset+ssLength]

	// Convert the Go string (ss) to a compact string object, usable by Java. Return to caller.
	obj := object.CreateCompactStringFromGoString(&ss)
	return obj

}

func getBytesVoid(params []interface{}) interface{} {
	switch params[0].(type) {
	case *object.Object:
		parmObj := params[0].(*object.Object)
		bytes := parmObj.Fields[0].Fvalue.(*[]byte)
		return bytes
	default:
		errMsg := fmt.Sprintf("In libs.GetBytesVoid, unexpected params[0] type=%T, value=%v", params[0], params[0])
		exceptions.Throw(exceptions.VirtualMachineError, errMsg)
		return nil
	}
}

func sprintf(params []interface{}) interface{} {
	// params[0]: format string
	// params[1]: object slice
	return StringFormatter(params)
}

func StringFormatter(params []interface{}) *object.Object {
	lenParams := len(params)
	if lenParams < 1 || lenParams > 2 {
		errMsg := fmt.Sprintf("StringFormatter: Invalid parameter count: %d", lenParams)
		exceptions.Throw(exceptions.IllegalClassFormatException, errMsg)
	}
	if lenParams == 1 { // No parameters beyond the format string
		formatStringObj := params[1].(*object.Object) // the format string is passed as a pointer to a string object
		return formatStringObj
	}
	formatStringObj := params[0].(*object.Object) // the format string is passed as a pointer to a string object
	formatString := object.GetGoStringFromJavaStringPtr(formatStringObj)
	valuesIn := *(params[1].(*object.Object).Fields[0].Fvalue).(*[]*object.Object) // ptr to slice of pointers to 1 or more objects
	valuesOut := []any{}

	for i := 0; i < len(valuesIn); i++ {
		// fmt.Printf("DEBUG i: %d of %d\n", i+1, len(valuesIn))
		// fmt.Printf("DEBUG valuesIn[i] klass: %s, fields: %v\n", *valuesIn[i].Klass, valuesIn[i].Fields)
		if object.IsJavaString(valuesIn[i]) {
			valuesOut = append(valuesOut, object.GetGoStringFromJavaStringPtr(valuesIn[i]))
			// fmt.Printf("DEBUG got a string: %s\n", object.GetGoStringFromJavaStringPtr(valuesIn[i]))
		} else {
			// str := valuesIn[i].FormatField()
			// fmt.Printf("DEBUG StringFormatter valuesIn[%d] FormatField:\n%s", i, str)

			// Establish a pointer to the field.
			var fldPtr *object.Field
			if len(valuesIn[i].FieldTable) > 0 { // using FieldTable
				fldPtr = valuesIn[i].FieldTable["value"]
			} else { // using Fields slice
				fldPtr = &valuesIn[i].Fields[0]
			}

			// Get the field value.
			fvalue := (*fldPtr).Fvalue

			// Process depending on field type
			switch (*fldPtr).Ftype {
			case types.Byte:
				valuesOut = append(valuesOut, fvalue.(int64))
			case types.Bool:
				// fmt.Printf("DEBUG %T %v\n", fvalue, fvalue)
				var zz bool
				if fvalue.(int64) == 0 {
					zz = false
				} else {
					zz = true
				}
				valuesOut = append(valuesOut, zz)
			case types.Char:
				valuesOut = append(valuesOut, fvalue.(int64))
			case types.Double:
				valuesOut = append(valuesOut, fvalue.(float64))
			case types.Float:
				valuesOut = append(valuesOut, fvalue.(float64))
			case types.Int:
				valuesOut = append(valuesOut, fvalue.(int64))
			case types.Long:
				valuesOut = append(valuesOut, fvalue.(int64))
			case types.Short:
				valuesOut = append(valuesOut, fvalue.(int64))
			default:
				errMsg := fmt.Sprintf("StringFormatter: Invalid parameter %d type %s", i+1, valuesIn[i].Fields[0].Ftype)
				exceptions.Throw(exceptions.IllegalClassFormatException, errMsg)
			}
		}
	}

	// Use golang fmt.Sprintf to do the heavy lifting.
	str := fmt.Sprintf(formatString, valuesOut...)

	// Return a pointer to an object.Object that wraps the string byte array.
	return object.CreateCompactStringFromGoString(&str)
}

func stringLength(params []interface{}) interface{} {
	var bytesPtr *[]byte
	parmObj := params[0].(*object.Object)
	if len(parmObj.FieldTable) > 0 {
		bytesPtr = parmObj.FieldTable["value"].Fvalue.(*[]byte)
	} else {
		bytesPtr = parmObj.Fields[0].Fvalue.(*[]byte)
	}
	return int64(len(*bytesPtr))

}

func toLowerCase(params []interface{}) interface{} {
	// params[0]: input string
	propObj := params[0].(*object.Object)
	strPtr := propObj.Fields[0].Fvalue.(*[]byte)
	str := strings.ToLower(string(*strPtr))
	obj := object.CreateCompactStringFromGoString(&str)
	return obj
}

func toUpperCase(params []interface{}) interface{} {
	// params[0]: input string
	propObj := params[0].(*object.Object)
	strPtr := propObj.Fields[0].Fvalue.(*[]byte)
	str := strings.ToUpper(string(*strPtr))
	obj := object.CreateCompactStringFromGoString(&str)
	return obj
}

func valueOfBoolean(params []interface{}) interface{} {
	// params[0]: input boolean
	value := params[0].(int64)
	var str string
	if value != 0 {
		str = "true"
	} else {
		str = "false"
	}
	obj := object.CreateCompactStringFromGoString(&str)
	return obj
}

func valueOfChar(params []interface{}) interface{} {
	// params[0]: input char
	value := params[0].(int64)
	str := fmt.Sprintf("%c", value)
	obj := object.CreateCompactStringFromGoString(&str)
	return obj
}

func valueOfCharArray(params []interface{}) interface{} {
	// params[0]: input char array
	propObj := params[0].(*object.Object)
	caPtr := propObj.Fields[0].Fvalue.(*[]int64)
	var str string
	for _, ch := range *caPtr {
		str += fmt.Sprintf("%c", ch)
	}
	obj := object.CreateCompactStringFromGoString(&str)
	return obj
}

func valueOfCharSubarray(params []interface{}) interface{} {
	// params[0]: input char array
	// params[1]: input offset
	// params[2]: input count
	propObj := params[0].(*object.Object)
	caPtr := propObj.Fields[0].Fvalue.(*[]int64)
	var wholeString string
	for _, ch := range *caPtr {
		wholeString += fmt.Sprintf("%c", ch)
	}
	// Get substring offset and count
	ssOffset := params[1].(int64)
	ssCount := params[2].(int64)

	// Validate boundaries.
	wholeLength := int64(len(wholeString))
	if wholeLength < 1 || ssOffset < 0 || ssCount < 1 || ssOffset > (wholeLength-1) || (ssOffset+ssCount) > wholeLength {
		errMsg := "In valueOfCharSubarray, either: nil input byte array, invalid substring offset, or invalid substring length"
		exceptions.Throw(exceptions.StringIndexOutOfBoundsException, errMsg)
	}

	// Compute substring.
	str := wholeString[ssOffset : ssOffset+ssCount]

	obj := object.CreateCompactStringFromGoString(&str)
	return obj
}

func valueOfDouble(params []interface{}) interface{} {
	// params[0]: input double
	value := params[0].(float64)
	str := strconv.FormatFloat(value, 'f', -1, 64)
	if !strings.Contains(str, ".") {
		str += ".0"
	}
	obj := object.CreateCompactStringFromGoString(&str)
	return obj
}

func valueOfFloat(params []interface{}) interface{} {
	// params[0]: input double
	value := params[0].(float64)
	// str := fmt.Sprintf("%.0g", value)
	str := strconv.FormatFloat(value, 'f', -1, 64)
	if !strings.Contains(str, ".") {
		str += ".0"
	}
	obj := object.CreateCompactStringFromGoString(&str)
	return obj
}

func valueOfInt(params []interface{}) interface{} {
	// params[0]: input int
	value := params[0].(int64)
	str := fmt.Sprintf("%d", value)
	obj := object.CreateCompactStringFromGoString(&str)
	return obj
}

func valueOfLong(params []interface{}) interface{} {
	// params[0]: input long
	value := params[0].(int64)
	str := fmt.Sprintf("%d", value)
	obj := object.CreateCompactStringFromGoString(&str)
	return obj
}

func valueOfObject(params []interface{}) interface{} {
	// params[0]: input Object
	ptrObj := params[0].(*object.Object)
	str := ptrObj.FormatField("")
	obj := object.CreateCompactStringFromGoString(&str)
	return obj
}

func compareToCaseSensitive(params []interface{}) interface{} {
	propObj := params[0].(*object.Object)
	strPtr := propObj.Fields[0].Fvalue.(*[]byte)
	str1 := string(*strPtr)
	propObj = params[1].(*object.Object)
	strPtr = propObj.Fields[0].Fvalue.(*[]byte)
	str2 := string(*strPtr)
	if str2 == str1 {
		return int64(0)
	}
	if str1 < str2 {
		return int64(-1)
	}
	return int64(1)
}

func compareToIgnoreCase(params []interface{}) interface{} {
	propObj := params[0].(*object.Object)
	strPtr := propObj.Fields[0].Fvalue.(*[]byte)
	str1 := strings.ToLower(string(*strPtr))
	propObj = params[1].(*object.Object)
	strPtr = propObj.Fields[0].Fvalue.(*[]byte)
	str2 := strings.ToLower(string(*strPtr))
	if str2 == str1 {
		return int64(0)
	}
	if str1 < str2 {
		return int64(-1)
	}
	return int64(1)
}

func stringConcat(params []interface{}) interface{} {
	propObj := params[0].(*object.Object)
	strPtr := propObj.Fields[0].Fvalue.(*[]byte)
	strRef := strings.ToLower(string(*strPtr))
	propObj = params[1].(*object.Object)
	strPtr = propObj.Fields[0].Fvalue.(*[]byte)
	strArg := strings.ToLower(string(*strPtr))
	str := strRef + strArg
	obj := object.CreateCompactStringFromGoString(&str)
	return obj
}
