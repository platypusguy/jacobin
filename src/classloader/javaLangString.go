/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package classloader

import (
	"fmt"
	"jacobin/exceptions"
	"jacobin/libs"
	"jacobin/object"
	"jacobin/types"
)

// IMPORTANT NOTE: Some String functions are placed in libs\javaLangStringMethods.go
// due to golang circularity concerns, alas.

/*
   We don't run String's static initializer block because the initialization
   is already handled in String creation
*/

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

	// String(byte[] ascii, int hibyte) ***************************************** DEPRECATED
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

	// String(byte[] ascii, int hibyte, int offset, int count) *****************- DEPRECATED
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

	// String(byte[] bytes, String charsetName) ******************************** CHARSET
	MethodSignatures["java/lang/String.<init>([BLjava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  noSupportYetInString,
		}

	// String(byte[] bytes, Charset charset) ********************************** CHARSET
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
			GFunction:  libs.GetBytesVoid,
		}

	// get the bytes from a string, given the Charset string name ************************ CHARSET
	MethodSignatures["java/lang/String.getBytes(Ljava/lang/String;)[B"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  noSupportYetInString,
		}

	// get the bytes from a string, given the specified Charset object *****************- CHARSET
	MethodSignatures["java/lang/String.getBytes(Ljava/nio/charset/Charset;)[B"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  noSupportYetInString,
		}

	return MethodSignatures

}

func stringClinit([]interface{}) interface{} {
	klass := MethAreaFetch("java/lang/String")
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
	klass := MethAreaFetch("java/lang/String")
	if klass == nil {
		errMsg := "In newStringFromBytes, expected java/lang/String to be in the MethodArea, but it was not"
		exceptions.Throw(exceptions.VirtualMachineError, errMsg)
	}
	klass.Data.ClInit = types.ClInitRun // just mark that String.<clinit>() has been run

	// Fetch a pointer to the raw slice of bytes from params[0].
	// Convert the raw slice of bytes to a Go string.
	wholeString := getGoString(params[0])

	// Convert the Go string to a compact string object, usable by Java. Return to caller.
	obj := object.CreateCompactStringFromGoString(&wholeString)
	return obj

}

// Construct a compact string object (usable by Java) from a Go byte array.
func newSubstringFromBytes(params []interface{}) interface{} {
	klass := MethAreaFetch("java/lang/String")
	if klass == nil {
		errMsg := "In newStringFromBytes, expected java/lang/String to be in the MethodArea, but it was not"
		exceptions.Throw(exceptions.VirtualMachineError, errMsg)
	}
	klass.Data.ClInit = types.ClInitRun // just mark that String.<clinit>() has been run

	// Fetch a pointer to the raw slice of bytes from params[0].
	// Convert the raw slice of bytes to a Go string.
	wholeString := getGoString(params[0])

	// Get substring offset and length
	ssOffset := params[1].(int64)
	ssLength := params[2].(int64)

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
