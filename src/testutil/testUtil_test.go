package testutil

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"math"
	"testing"
)

// Unit test the unit test utilities.

const loopCount = 20
const absTolerance = 1e-9
const relTolerance = 1e-6

// Close enough?
// If very small, use an absolute tolerance.
// Otherwise, use a relative tolerance.
func float64CloseEnough(tweedleDee, tweedleDum float64) bool {
	diff := math.Abs(tweedleDee - tweedleDum)

	// Absolute difference check for very small values.
	if diff <= absTolerance {
		return true
	}

	// Compute the maximum of the magnitudes.
	maxMagnitude := math.Max(math.Abs(tweedleDee), math.Abs(tweedleDum))

	// Relative difference check for larger values.
	return diff <= maxMagnitude*relTolerance
}

func TestUTinit(t *testing.T) {
	for ix := 0; ix < loopCount; ix++ {
		UTinit(t)
	}
}

func TestUTconsoles(t *testing.T) {
	for ix := 0; ix < loopCount; ix++ {
		UTnewConsole(t)
		UTrestoreConsole(t)
	}
}

func TestUTensemble(t *testing.T) {
	for ix := 0; ix < loopCount; ix++ {
		UTinit(t)
		UTnewConsole(t)
		UTrestoreConsole(t)
	}
}

func TestMathSqrt_D_D(t *testing.T) {
	var params []interface{}
	params = append(params, 4.0)
	result := UTgfunc(t, "java/lang/Math", "sqrt", "(D)D", nil, params)
	switch result.(type) {
	case float64:
		if !float64CloseEnough(result.(float64), 2.0) {
			t.Errorf("TestMathSqrt_D_D: Expected sqrt(4)=2, observed: %f", result.(float64))
		}
	case ghelpers.GErrBlk:
		t.Errorf("TestMathSqrt_D_D: ghelpers.GErrBlk.ErrMsg: %s", result.(ghelpers.GErrBlk).ErrMsg)
	default:
		t.Errorf("TestMathSqrt_D_D: Expected result of type float64, observed: %T", result)
	}
}

func TestStrictMathCbrt_D_D(t *testing.T) {
	var params []interface{}
	params = append(params, 27.0)
	result := UTgfunc(t, "java/lang/StrictMath", "cbrt", "(D)D", nil, params)
	switch result.(type) {
	case float64:
		if !float64CloseEnough(result.(float64), 3.0) {
			t.Errorf("TestStrictMathCbrt_D_D: Expected cbrt(27)=3, observed: %f", result.(float64))
		}
	case ghelpers.GErrBlk:
		t.Errorf("TestStrictMathCbrt_D_D: ghelpers.GErrBlk.ErrMsg: %s", result.(ghelpers.GErrBlk).ErrMsg)
	default:
		t.Errorf("TestStrictMathCbrt_D_D: Expected result of type float64, observed: %T", result)
	}
}

func TestStringLength(t *testing.T) {
	var params []interface{}
	obj := object.StringObjectFromGoString("ABCDEF")
	result := UTgfunc(t, "java/lang/String", "length", "()I", obj, params)
	switch result.(type) {
	case int64:
		if result.(int64) != 6 {
			t.Errorf("TestStringLength: Expected length(\"ABCDEF\")=6, observed: %d", result.(int64))
		}
	case ghelpers.GErrBlk:
		t.Errorf("TestStringLength: ghelpers.GErrBlk.ErrMsg: %s", result.(ghelpers.GErrBlk).ErrMsg)
	default:
		t.Errorf("TestStringLength: Expected result of type int64, observed: %T", result)
	}
}

func TestStringRepeater(t *testing.T) {
	var params []interface{}
	objIn := object.StringObjectFromGoString("Beetlejuice")
	params = append(params, int64(3))
	expected := "BeetlejuiceBeetlejuiceBeetlejuice"
	result := UTgfunc(t, "java/lang/String", "repeat", "(I)Ljava/lang/String;", objIn, params)
	switch result.(type) {
	case *object.Object:
		observed := object.GoStringFromStringObject(result.(*object.Object))
		if observed != expected {
			t.Errorf("TestStringRepeater: Expected: \"%s\", observed: \"%s\"", expected, observed)
		}
	case ghelpers.GErrBlk:
		t.Errorf("TestStringRepeater: ghelpers.GErrBlk.ErrMsg: %s", result.(ghelpers.GErrBlk).ErrMsg)
	default:
		t.Errorf("TestStringRepeater: Expected result of type String object, observed: %T", result)
	}
}

func TestStringToCharArray(t *testing.T) {
	var params []interface{}

	// Here's the base string and its bytes equivalent.
	expected := "ABCDEF"
	bytes := []byte(expected)

	// Make a Java object of the string.
	objIn := object.StringObjectFromGoString(expected)

	// Make an int64 array of the bytes.
	var iArray []int64
	for ix := 0; ix < len(expected); ix++ {
		iArray = append(iArray, int64(bytes[ix]))
	}

	// Try it.
	result := UTgfunc(t, "java/lang/String", "toCharArray", "()[C", objIn, params)
	switch result.(type) {
	case *object.Object:
		observed := object.ObjectFieldToString(result.(*object.Object), "value")
		if observed != expected {
			t.Errorf("TestStringToCharArray: Expected: \"%s\", observed: \"%s\"", expected, observed)
		}
	case ghelpers.GErrBlk:
		t.Errorf("TestStringToCharArray: ghelpers.GErrBlk.ErrMsg: %s", result.(ghelpers.GErrBlk).ErrMsg)
	default:
		t.Errorf("TestStringToCharArray: Expected result of type [C object, observed: %T", result)
	}
}

func TestBigInteger(t *testing.T) {

	var params []interface{}
	biClassName := "java/math/BigInteger"
	expected := int64(42)

	// Need to do an early initialisation before object.MakeEmptyObjectWithClassName).
	UTinit(t)

	// Initialise the base object.
	obj := object.MakeEmptyObjectWithClassName(&biClassName)
	ghelpers.InitBigIntegerField(obj, int64(0))

	// Add the String value to params.
	params = append(params, object.StringObjectFromGoString("42"))

	// Try java/math/BigInteger.<init>.(Ljava/lang/String;)V
	result := UTgfunc(t, "java/math/BigInteger", "<init>", "(Ljava/lang/String;)V", obj, params)
	if result != nil {
		t.Fatalf("TestBigInteger: <init>: Expected nil return")
	}

	// Try intValue()
	params = nil
	result = UTgfunc(t, "java/math/BigInteger", "intValue", "()I", obj, params)
	switch result.(type) {
	case int64:
		if result.(int64) != expected {
			t.Errorf("TestStringToCharArray: Expected: \"%d\", observed: \"%d\"", expected, result.(int64))
		}
	case ghelpers.GErrBlk:
		t.Errorf("TestStringToCharArray: ghelpers.GErrBlk.ErrMsg: %s", result.(ghelpers.GErrBlk).ErrMsg)
	default:
		t.Errorf("TestStringToCharArray: Expected result of type [C object, observed: %T", result)
	}
}
