package gfunction

import (
	"jacobin/object"
	"jacobin/types"
	"testing"
)

func TestMathSqrt(t *testing.T) {
	var params []interface{}
	params = append(params, 4.0)
	TGrunner(t, "java/lang/Math", "sqrt", "(D)D", 2.0, nil, params)
}

func TestStringLength(t *testing.T) {
	var params []interface{}
	obj := object.StringObjectFromGoString("ABCDEF")
	TGrunner(t, "java/lang/String", "length", "()I", int64(6), obj, params)
}

func TestStringRepeater(t *testing.T) {
	var params []interface{}
	objIn := object.StringObjectFromGoString("Beetlejuice")
	params = append(params, int64(3))
	expected := object.StringObjectFromGoString("BeetlejuiceBeetlejuiceBeetlejuice")
	TGrunner(t, "java/lang/String", "repeat", "(I)Ljava/lang/String;", expected, objIn, params)
}

func TestStringToCharArray(t *testing.T) {
	var params []interface{}
	// I need the string pool set up before calling populator
	if !TGinit(t) {
		return // failed, already reported
	}
	// Here's the base string and its bytes equivalent.
	str := "ABCDEF"
	bytes := []byte(str)
	// Make a Java object of the string.
	objIn := object.StringObjectFromGoString(str)
	// Make an int64 array of the bytes.
	var iArray []int64
	for ix := 0; ix < len(str); ix++ {
		iArray = append(iArray, int64(bytes[ix]))
	}
	// Make a Java object of the int64 array.
	expected := populator("[C", types.IntArray, iArray)
	// Try it.
	TGrunner(t, "java/lang/String", "toCharArray", "()[C", expected, objIn, params)
}

func TestBigIntegers(t *testing.T) {
	var params []interface{}
	biClassName := "java/math/BigInteger"
	// I need the string pool set up before calling populator
	if !TGinit(t) {
		return // failed, already reported
	}
	// Build the expected value object.
	expected := object.MakeEmptyObjectWithClassName(&biClassName)
	initBigIntegerField(expected, int64(42))
	// Initialise the base object.
	obj := object.MakeEmptyObjectWithClassName(&biClassName)
	initBigIntegerField(obj, int64(0))
	// Add the String value to the parameter.
	params = append(params, object.StringObjectFromGoString("42"))
	// Try <init>.
	TGrunner(t, "java/math/BigInteger", "<init>", "(Ljava/lang/String;)V",
		nil, obj, params)
	// Try intValue()
	params = nil
	TGrunner(t, "java/math/BigInteger", "intValue", "()I",
		int64(42), obj, params)
}

func TestHexFormatBytes(t *testing.T) {
	var params []interface{}
	UPPERCASE_DIGITS := []types.JavaByte{
		'0', '1', '2', '3', '4', '5', '6', '7',
		'8', '9', 'A', 'B', 'C', 'D', 'E', 'F',
	}
	// I need the string pool set up before calling populator
	if !TGinit(t) {
		return // failed, already reported
	}
	// Create a HexFormat base object.
	delimiter := []types.JavaByte{}
	prefix := []types.JavaByte{}
	suffix := []types.JavaByte{}
	objBase := mkHexFormatObject(delimiter, prefix, suffix, UPPERCASE_DIGITS)
	// Input bytes to format in hex. Append the object to params.
	bytes := object.JavaByteArrayFromGoString("ABCDEF")
	objBytes := object.MakePrimitiveObject("[B", types.ByteArray, bytes)
	params = append(params, objBytes)
	// Make an object with the expected output.
	objExp := object.StringObjectFromGoString("414243444546")
	// Try it.
	TGrunner(t, "java/util/HexFormat", "formatHex", "([B)Ljava/lang/String;",
		objExp, objBase, params)
}
