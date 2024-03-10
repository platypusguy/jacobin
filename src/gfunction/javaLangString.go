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
			GFunction:  newStringFromBytesSubset,
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

	// String(char[] value) ****************************** works fine with JDK libraries
	// String(char[] value, int offset, int count) ******* works fine with JDK libraries

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

	// // Are 2 strings equal?
	MethodSignatures["java/lang/String.equals(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  stringEquals,
		}

	// get the bytes from a string
	MethodSignatures["java/lang/String.getBytes()[B"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  getBytesFromString,
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
	klass := classloader.MethAreaFetch(object.StringClassName)
	if klass == nil {
		errMsg := fmt.Sprintf("stringClinit: Could not find class %s in the MethodArea", object.StringClassName)
		return getGErrBlk(exceptions.ClassNotLoadedException, errMsg)
	}
	klass.Data.ClInit = types.ClInitRun // just mark that String.<clinit>() has been run
	return nil
}

// No support YET for references to Charset objects nor for Unicode code point arrays
func noSupportYetInString([]interface{}) interface{} {
	errMsg := fmt.Sprintf("%s: No support yet for user-specified character sets and Unicode code point arrays", object.StringClassName)
	return getGErrBlk(exceptions.UnsupportedEncodingException, errMsg)
}

// Are 2 strings equal?
func stringEquals(params []interface{}) interface{} {
	// params[0]: reference string object
	// params[1]: compare-to string Object

	// Unpack the reference string.
	ptrObj := params[0].(*object.Object)
	fld := ptrObj.FieldTable["value"]
	if fld.Ftype != types.StringIndex {
		errMsg := "stringEquals: reference object must be a String"
		return getGErrBlk(exceptions.VirtualMachineError, errMsg)
	}
	str1 := object.GetGoStringFromObject(ptrObj)

	// Unpack the compare-to string
	ptrObj = params[1].(*object.Object)
	fld = ptrObj.FieldTable["value"]
	if fld.Ftype != types.StringIndex {
		return int64(0) // Not a string, return false
	}
	str2 := object.GetGoStringFromObject(ptrObj)

	// Are they equal in value?
	if str1 == str2 {
		return int64(1) // true
	} else {
		return int64(0) // false
	}
}

// Construct a compact string object (usable by Java) from a Go byte array.
func newStringFromBytes(params []interface{}) interface{} {
	// params[1] = byte array object
	// params[0] = target object for string (updated)
	klass := classloader.MethAreaFetch(object.StringClassName)
	if klass == nil {
		errMsg := fmt.Sprintf("newStringFromBytes: Expected %s to be in the MethodArea, but it was not", object.StringClassName)
		return getGErrBlk(exceptions.VirtualMachineError, errMsg)
	}

	// Mark that String.<clinit>() has been run.
	klass.Data.ClInit = types.ClInitRun

	var str string
	switch params[1].(type) {
	case *object.Object:
		bytes := params[1].(*object.Object).FieldTable["value"].Fvalue.([]byte)
		str = string(bytes)
	case []byte:
		str = string(params[1].([]byte))
	}
	object.UpdateObjectFromGoString(params[0].(*object.Object), str)
	return nil
}

// Construct a compact string object (usable by Java) from a Go byte array.
func newStringFromBytesSubset(params []interface{}) interface{} {
	klass := classloader.MethAreaFetch(object.StringClassName)
	if klass == nil {
		errMsg := fmt.Sprintf("newStringFromBytes: Expected %s to be in the MethodArea, but it was not", object.StringClassName)
		return getGErrBlk(exceptions.VirtualMachineError, errMsg)
	}

	// Mark that String.<clinit>() has been run.
	klass.Data.ClInit = types.ClInitRun

	var bytes []byte
	switch params[1].(type) {
	case *object.Object:
		bytes = params[1].(*object.Object).FieldTable["value"].Fvalue.([]byte)
	case []byte:
		bytes = params[1].([]byte)
	}

	// Get substring offset and length
	ssOffset := params[2].(int64)
	ssLength := params[3].(int64)

	// Validate boundaries.
	totalLength := int64(len(bytes))
	if totalLength < 1 || ssOffset < 0 || ssLength < 1 || ssOffset > (totalLength-1) || (ssOffset+ssLength) > totalLength {
		errMsg1 := "newStringFromBytesSubset: Either: nil input byte array, invalid substring offset, or invalid substring length"
		errMsg2 := fmt.Sprintf("\n\twhole='%s' wholelen=%d, offset=%d, sslen=%d\n\n", string(bytes), totalLength, ssOffset, ssLength)
		return getGErrBlk(exceptions.StringIndexOutOfBoundsException, errMsg1+errMsg2)
	}

	// Compute substring.
	bytes = bytes[ssOffset : ssOffset+ssLength]

	// Copy subset of byte array from params[1] to the whole byte array in params[0].
	str := string(bytes)
	object.UpdateObjectFromGoString(params[0].(*object.Object), str)
	return nil

}

func getBytesFromString(params []interface{}) interface{} {
	switch params[0].(type) {
	case *object.Object:
		str := object.GetGoStringFromObject(params[0].(*object.Object))
		return []byte(str)
	default:
		errMsg := fmt.Sprintf("getBytesFromString: Unexpected params[0] type=%T, value=%v", params[0], params[0])
		return getGErrBlk(exceptions.VirtualMachineError, errMsg)
	}
}

func sprintf(params []interface{}) interface{} {
	// params[0]: format string
	// params[1]: object slice
	return StringFormatter(params)
}

func StringFormatter(params []interface{}) interface{} {
	lenParams := len(params)
	//fmt.Printf("DEBUG StringFormatter lenParams=%d\n", lenParams)
	if lenParams < 1 || lenParams > 2 {
		errMsg := fmt.Sprintf("StringFormatter: Invalid parameter count: %d", lenParams)
		return getGErrBlk(exceptions.IllegalArgumentException, errMsg)
	}
	if lenParams == 1 { // No parameters beyond the format string
		formatStringObj := params[0].(*object.Object) // the format string is passed as a pointer to a string object
		return formatStringObj
	}

	var formatString string
	switch params[0].(type) {
	case *object.Object:
		formatStringObj := params[0].(*object.Object) // the format string is passed as a pointer to a string object
		switch formatStringObj.FieldTable["value"].Ftype {
		case types.StringIndex:
			formatString = object.GetGoStringFromObject(formatStringObj)
		case types.ByteArray:
			formatString = string(formatStringObj.FieldTable["value"].Fvalue.([]byte))
		default:
			errMsg := fmt.Sprintf("StringFormatter: Invalid Ftype in format string object: %s",
				formatStringObj.FieldTable["value"].Ftype)
			return getGErrBlk(exceptions.IllegalArgumentException, errMsg)
		}
		//fmt.Printf("DEBUG StringFormatter formatString=%s\n", formatString)
	default:
		errMsg := fmt.Sprintf("StringFormatter: Invalid variable type for format string: %T", params[0])
		return getGErrBlk(exceptions.IllegalArgumentException, errMsg)
	}

	valuesOut := []any{}
	fmt.Printf("DEBUG StringFormatter params[1]=%T\n", params[1])
	fld := params[1].(*object.Object).FieldTable["value"]
	//fmt.Printf("DEBUG StringFormatter fld.Ftype=%s, fld.Fvalue=%v\n", fld.Ftype, fld.Fvalue)
	if !strings.HasPrefix(fld.Ftype, types.RefArray) {
		errMsg := fmt.Sprintf("StringFormatter: Expected Ftype=%s for params[1]: fld.Ftype=%s, fld.Fvalue=%v",
			types.RefArray, fld.Ftype, fld.Fvalue)
		return getGErrBlk(exceptions.IllegalArgumentException, errMsg)
	}
	valuesIn := fld.Fvalue.([]*object.Object) // ptr to slice of pointers to 1 or more objects

	for ii := 0; ii < len(valuesIn); ii++ {
		//fmt.Printf("DEBUG StringFormatter ii: %d of %d\n", ii+1, len(valuesIn))
		//fmt.Printf("DEBUG StringFormatter valuesIn[ii] field type: %s, field value: %v\n",
		//	valuesIn[ii].FieldTable["value"].Ftype, valuesIn[ii].FieldTable["value"].Fvalue)
		if valuesIn[ii].FieldTable["value"].Ftype == types.StringIndex {
			str := object.GetGoStringFromObject(valuesIn[ii])
			valuesOut = append(valuesOut, str)
			//fmt.Printf("DEBUG StringFormatter got a string: %s\n", str)
		} else {
			//str := valuesIn[ii].FormatField("value")
			//fmt.Printf("DEBUG StringFormatter valuesIn[%d] FormatField: %s\n", ii, str)

			// Extract the field.
			fld := valuesIn[ii].FieldTable["value"]

			// Process depending on field type
			switch fld.Ftype {
			case types.ByteArray:
				str := string(fld.Fvalue.([]byte))
				valuesOut = append(valuesOut, str)
			case types.Byte:
				valuesOut = append(valuesOut, fld.Fvalue.(int64))
			case types.Bool:
				// fmt.Printf("DEBUG %T %v\n", fvalue, fvalue)
				var zz bool
				if fld.Fvalue.(int64) == 0 {
					zz = false
				} else {
					zz = true
				}
				valuesOut = append(valuesOut, zz)
			case types.Char:
				cc := fmt.Sprint(fld.Fvalue.(int64))
				valuesOut = append(valuesOut, cc)
			case types.Double:
				valuesOut = append(valuesOut, fld.Fvalue.(float64))
			case types.Float:
				valuesOut = append(valuesOut, fld.Fvalue.(float64))
			case types.Int:
				valuesOut = append(valuesOut, fld.Fvalue.(int64))
			case types.Long:
				valuesOut = append(valuesOut, fld.Fvalue.(int64))
			case types.Short:
				valuesOut = append(valuesOut, fld.Fvalue.(int64))
			default:
				errMsg := fmt.Sprintf("StringFormatter: Invalid parameter %d type %s", ii+1, fld.Ftype)
				return getGErrBlk(exceptions.IllegalArgumentException, errMsg)
			}
		}
	}

	// Use golang fmt.Sprintf to do the heavy lifting.
	str := fmt.Sprintf(formatString, valuesOut...)

	// Return a pointer to an object.Object that wraps the string byte array.
	return object.NewPoolStringFromGoString(str)
}

func stringLength(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	fld := obj.FieldTable["value"]
	if fld.Ftype != types.StringIndex {
		errMsg := "stringLength: reference object must be a String"
		return getGErrBlk(exceptions.VirtualMachineError, errMsg)
	}
	bytes := []byte(object.GetGoStringFromObject(obj))
	return int64(len(bytes))
}

func toLowerCase(params []interface{}) interface{} {
	// params[0]: input string
	str := strings.ToLower(object.GetGoStringFromObject(params[0].(*object.Object)))
	obj := object.NewPoolStringFromGoString(str)
	return obj
}

func toUpperCase(params []interface{}) interface{} {
	// params[0]: input string
	str := strings.ToUpper(object.GetGoStringFromObject(params[0].(*object.Object)))
	obj := object.NewPoolStringFromGoString(str)
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
	obj := object.NewPoolStringFromGoString(str)
	return obj
}

func valueOfChar(params []interface{}) interface{} {
	// params[0]: input char
	value := params[0].(int64)
	str := fmt.Sprintf("%c", value)
	obj := object.NewPoolStringFromGoString(str)
	return obj
}

func valueOfCharArray(params []interface{}) interface{} {
	// params[0]: input char array
	propObj := params[0].(*object.Object)
	intArray := propObj.FieldTable["value"].Fvalue.([]int64)
	var str string
	for _, ch := range intArray {
		str += fmt.Sprintf("%c", ch)
	}
	obj := object.NewPoolStringFromGoString(str)
	return obj
}

func valueOfCharSubarray(params []interface{}) interface{} {
	// params[0]: input char array
	// params[1]: input offset
	// params[2]: input count
	propObj := params[0].(*object.Object)
	intArray := propObj.FieldTable["value"].Fvalue.([]int64)
	var wholeString string
	for _, ch := range intArray {
		wholeString += fmt.Sprintf("%c", ch)
	}
	// Get substring offset and count
	ssOffset := params[1].(int64)
	ssCount := params[2].(int64)

	// Validate boundaries.
	wholeLength := int64(len(wholeString))
	if wholeLength < 1 || ssOffset < 0 || ssCount < 1 || ssOffset > (wholeLength-1) || (ssOffset+ssCount) > wholeLength {
		errMsg := "valueOfCharSubarray: Either: nil input byte array, invalid substring offset, or invalid substring length"
		return getGErrBlk(exceptions.StringIndexOutOfBoundsException, errMsg)
	}

	// Compute substring.
	str := wholeString[ssOffset : ssOffset+ssCount]

	obj := object.NewPoolStringFromGoString(str)
	return obj
}

func valueOfDouble(params []interface{}) interface{} {
	// params[0]: input double
	value := params[0].(float64)
	str := strconv.FormatFloat(value, 'f', -1, 64)
	if !strings.Contains(str, ".") {
		str += ".0"
	}
	obj := object.NewPoolStringFromGoString(str)
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
	obj := object.NewPoolStringFromGoString(str)
	return obj
}

func valueOfInt(params []interface{}) interface{} {
	// params[0]: input int
	value := params[0].(int64)
	str := fmt.Sprintf("%d", value)
	obj := object.NewPoolStringFromGoString(str)
	return obj
}

func valueOfLong(params []interface{}) interface{} {
	// params[0]: input long
	value := params[0].(int64)
	str := fmt.Sprintf("%d", value)
	obj := object.NewPoolStringFromGoString(str)
	return obj
}

func valueOfObject(params []interface{}) interface{} {
	// params[0]: input Object
	ptrObj := params[0].(*object.Object)
	str := ptrObj.FormatField("")
	obj := object.NewPoolStringFromGoString(str)
	return obj
}

func compareToCaseSensitive(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	str1 := object.GetGoStringFromObject(obj)
	obj = params[1].(*object.Object)
	str2 := object.GetGoStringFromObject(obj)
	if str2 == str1 {
		return int64(0)
	}
	if str1 < str2 {
		return int64(-1)
	}
	return int64(1)
}

func compareToIgnoreCase(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	str1 := strings.ToLower(object.GetGoStringFromObject(obj))
	obj = params[1].(*object.Object)
	str2 := strings.ToLower(object.GetGoStringFromObject(obj))
	if str2 == str1 {
		return int64(0)
	}
	if str1 < str2 {
		return int64(-1)
	}
	return int64(1)
}

func stringConcat(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	strRef := object.GetGoStringFromObject(obj)
	obj = params[1].(*object.Object)
	strArg := object.GetGoStringFromObject(obj)
	str := strRef + strArg
	obj = object.NewPoolStringFromGoString(str)
	return obj
}
