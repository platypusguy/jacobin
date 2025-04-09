/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"crypto/rand"
	"fmt"
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/statics"
	"jacobin/trace"
	"jacobin/types"
	"math/big"
	"math/bits"
)

/*
The BigInteger object is implemented using Golang package math/big.
*/

func Load_Math_Big_Integer() {

	MethodSignatures["java/math/BigInteger.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerClinit,
		}

	// <init> functions

	MethodSignatures["java/math/BigInteger.<init>([B)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerInitByteArray,
		}

	MethodSignatures["java/math/BigInteger.<init>([BII)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/math/BigInteger.<init>(I[B)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/math/BigInteger.<init>(I[BII)V"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/math/BigInteger.<init>(IILjava/util/Random;)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  bigIntegerInitProbablyPrime,
		}

	MethodSignatures["java/math/BigInteger.<init>(ILjava/util/Random;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigIntegerInitRandom,
		}

	MethodSignatures["java/math/BigInteger.<init>(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerInitString,
		}

	MethodSignatures["java/math/BigInteger.<init>(Ljava/lang/String;I)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigIntegerInitStringRadix,
		}

	// Member functions

	MethodSignatures["java/math/BigInteger.abs()Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerAbs,
		}

	MethodSignatures["java/math/BigInteger.add(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerAdd,
		}

	MethodSignatures["java/math/BigInteger.and(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerAnd,
		}

	MethodSignatures["java/math/BigInteger.andNot(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerAndNot,
		}

	MethodSignatures["java/math/BigInteger.bitCount()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerBitCount,
		}

	MethodSignatures["java/math/BigInteger.bitLength()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerBitLength,
		}

	MethodSignatures["java/math/BigInteger.byteValueExact()B"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerByteValueExact,
		}

	MethodSignatures["java/math/BigInteger.clearBit(I)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/math/BigInteger.compareTo(Ljava/math/BigInteger;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerCompareTo,
		}

	MethodSignatures["java/math/BigInteger.divide(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerDivide,
		}

	MethodSignatures["java/math/BigInteger.divideAndRemainder(Ljava/math/BigInteger;)[Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerDivideAndRemainder,
		}

	MethodSignatures["java/math/BigInteger.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerFloat64Value,
		}

	MethodSignatures["java/math/BigInteger.equals(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerEquals,
		}

	MethodSignatures["java/math/BigInteger.flipBit(I)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/math/BigInteger.floatValue()F"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerFloat64Value,
		}

	MethodSignatures["java/math/BigInteger.gcd(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerGCD,
		}

	MethodSignatures["java/math/BigInteger.getLowestSetBit()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/math/BigInteger.hashCode()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/math/BigInteger.intValue()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerInt64Value,
		}

	MethodSignatures["java/math/BigInteger.intValueExact()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerInt64Value,
		}

	MethodSignatures["java/math/BigInteger.isProbablePrime(I)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerIsProbablePrime,
		}

	MethodSignatures["java/math/BigInteger.longValue()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerInt64Value,
		}

	MethodSignatures["java/math/BigInteger.longValueExact()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerInt64Value,
		}

	MethodSignatures["java/math/BigInteger.max(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerMax,
		}

	MethodSignatures["java/math/BigInteger.min(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerMin,
		}

	MethodSignatures["java/math/BigInteger.mod(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerMod,
		}

	MethodSignatures["java/math/BigInteger.modInverse(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerModInverse,
		}

	MethodSignatures["java/math/BigInteger.modPow(Ljava/math/BigInteger;Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigIntegerModPow,
		}

	MethodSignatures["java/math/BigInteger.multiply(J)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerMultiply,
		}

	MethodSignatures["java/math/BigInteger.multiply(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerMultiply,
		}

	MethodSignatures["java/math/BigInteger.negate()Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerNegate,
		}

	MethodSignatures["java/math/BigInteger.nextProbablePrime()Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/math/BigInteger.not()Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerNot,
		}

	MethodSignatures["java/math/BigInteger.or(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerOr,
		}

	MethodSignatures["java/math/BigInteger.pow(I)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerPow,
		}

	MethodSignatures["java/math/BigInteger.probablePrime(ILjava/util/Random;)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigIntegerProbablyPrime,
		}

	MethodSignatures["java/math/BigInteger.remainder(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerRemainder,
		}

	MethodSignatures["java/math/BigInteger.setBit(I)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerSetBit,
		}

	MethodSignatures["java/math/BigInteger.shiftLeft(I)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerShiftLeft,
		}

	MethodSignatures["java/math/BigInteger.shiftRight(I)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerShiftRight,
		}

	MethodSignatures["java/math/BigInteger.shortValueExact()S"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerInt64Value,
		}

	MethodSignatures["java/math/BigInteger.signum()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerSignum,
		}

	MethodSignatures["java/math/BigInteger.sqrt()Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerSqrt,
		}

	MethodSignatures["java/math/BigInteger.sqrtAndRemainder()[Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/math/BigInteger.subtract(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerSubtract,
		}

	MethodSignatures["java/math/BigInteger.testBit(I)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerTestBit,
		}

	MethodSignatures["java/math/BigInteger.toByteArray()[B"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerToByteArray,
		}

	MethodSignatures["java/math/BigInteger.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerToString,
		}

	MethodSignatures["java/math/BigInteger.toString(I)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerToStringRadix,
		}

	MethodSignatures["java/math/BigInteger.valueOf(J)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerValueOf,
		}

	MethodSignatures["java/math/BigInteger.xor(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerXor,
		}

}

var classNameBigInteger = "java/math/BigInteger"

// Get a prime number formatted as a big.Int.
func getPrime(bitLength int) (*big.Int, string) {
	for { // TODO: Infinite loop ok?
		// Generate a random number given the bit length.
		zz, err := rand.Prime(rand.Reader, bitLength)
		if err != nil {
			errMsg := fmt.Sprintf("getPrime: rand.Reader(bitLength=%d) failed, reason: %s", bitLength, err.Error())
			return nil, errMsg
		}

		// Check if the number is probably prime
		if zz.ProbablyPrime(20) { // 20 is the number of Miller-Rabin tests
			return zz, ""
		}
	}
}

// addStaticBigInteger: Form a BigInteger object based on the parameter value.
func addStaticBigInteger(argName string, argValue int64) {
	name := fmt.Sprintf("%s.%s", classNameBigInteger, argName)
	obj := object.MakeEmptyObjectWithClassName(&classNameBigInteger)
	InitBigIntegerField(obj, argValue)
	_ = statics.AddStatic(name, statics.Static{Type: "Ljava/math/BigInteger;", Value: obj})
}

// "java/math/BigInteger.<clinit>()V"
func bigIntegerClinit([]interface{}) interface{} {
	klass := classloader.MethAreaFetch(classNameBigInteger)
	if klass == nil {
		errMsg := fmt.Sprintf("bigIntegerClinit: Expected %s to be in the MethodArea, but it was not", classNameBigInteger)
		trace.Error(errMsg)
		return getGErrBlk(excNames.ClassNotLoadedException, errMsg)
	}
	if klass.Data.ClInit != types.ClInitRun {
		addStaticBigInteger("ONE", int64(1))
		addStaticBigInteger("TEN", int64(10))
		addStaticBigInteger("TWO", int64(2))
		addStaticBigInteger("ZERO", int64(0))
		klass.Data.ClInit = types.ClInitRun
	}
	return nil
}

// Convert a byte slice into a signed big integer.
// Thank you, ChatGPT.
func BytesToBigInt(buf []byte) (*big.Int, int64) {
	if len(buf) == 0 {
		return big.NewInt(0), int64(0)
	}

	// Check if the most significant bit is set (indicating a negative number).
	negative := buf[0]&0x80 != 0

	if negative {
		// Create a copy of the buffer to avoid modifying the original byte slice.
		twosComplement := make([]byte, len(buf))
		copy(twosComplement, buf)

		// Invert the bits (two's complement step 1).
		for i := range twosComplement {
			twosComplement[i] = ^twosComplement[i]
		}

		// Add one to the result (two's complement step 2).
		one := big.NewInt(1)
		twoComplementBigInt := new(big.Int).SetBytes(twosComplement)
		twoComplementBigInt.Add(twoComplementBigInt, one)

		// Negate the result to get the original negative number.
		twoComplementBigInt.Neg(twoComplementBigInt)

		return twoComplementBigInt, int64(-1)
	}

	// Not negative.
	// Use SetBytes to convert to a big integer.
	bigInt := new(big.Int).SetBytes(buf)

	// Get sign (+ or 0).
	signum := int64(bigInt.Sign())

	// Return result.
	return bigInt, signum
}

// "java/math/BigInteger.<init>([B)V"
func bigIntegerInitByteArray(params []interface{}) interface{} {
	// params[0]: base object
	// params[1]: byte array object
	obj := params[0].(*object.Object)
	fld := obj.FieldTable["value"]
	jba := params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	bytes := object.GoByteArrayFromJavaByteArray(jba)
	zz, signum := BytesToBigInt(bytes)

	// Set value to big integer.
	fld = object.Field{Ftype: types.BigInteger, Fvalue: zz}
	obj.FieldTable["value"] = fld

	// Set signum to sign.
	fld = object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	// Return void.
	return nil
}

// "java/math/BigInteger.<init>(ILjava/util/Random;)V"
func bigIntegerInitProbablyPrime(params []interface{}) interface{} {
	// params[0]: base object (to be updated).
	// params[1]: int64 holding bit length.
	// params[2]: int64 holding certainty (TODO: currently ignored).
	// params[3]: Random object (TODO: currently ignored).

	obj := params[0].(*object.Object)
	fld := obj.FieldTable["value"]
	bitLength := params[1].(int64)

	zz, errMsg := getPrime(int(bitLength))
	if zz != nil {
		// Set value to big integer.
		fld = object.Field{Ftype: types.BigInteger, Fvalue: zz}
		obj.FieldTable["value"] = fld

		// Set signum to sign.
		fld = object.Field{Ftype: types.BigInteger, Fvalue: int64(+1)}
		obj.FieldTable["signum"] = fld

		// Return void.
		return nil
	}
	return getGErrBlk(excNames.ArithmeticException, errMsg)
}

// "java/math/BigInteger.<init>(IILjava/util/Random;)V"
func bigIntegerInitRandom(params []interface{}) interface{} {
	// params[0]: base object
	// params[1]: int64 holding numbits such that the base object value field
	//            will be set to a random value in the rang given by [0 : 2**(numbits) - 1].
	// params[2]: Random object
	obj := params[0].(*object.Object)
	fld := obj.FieldTable["value"]
	numBits := params[1].(int64)
	// TODO: Ignored for now: objRandom := params[2].(*object.Object)

	// Compute upperBound = 2**(numBits) based on numBits.
	upperBound := new(big.Int).Lsh(big.NewInt(1), uint(numBits))

	// Get a big.Int in the randge of [0, upperBound].
	zz, err := rand.Int(rand.Reader, upperBound)
	if err != nil {
		errMsg := fmt.Sprintf("bigIntegerInitRandom: rand.Int(numBits=%d) failed, reason: %s", numBits, err.Error())
		return getGErrBlk(excNames.NumberFormatException, errMsg)
	}

	// Set value to big integer.
	fld = object.Field{Ftype: types.BigInteger, Fvalue: zz}
	obj.FieldTable["value"] = fld

	// Set signum to sign.
	fld = object.Field{Ftype: types.BigInteger, Fvalue: int64(+1)}
	obj.FieldTable["signum"] = fld

	// Return void.
	return nil
}

// "java/math/BigInteger.<init>(Ljava/lang/String;)V"
func bigIntegerInitString(params []interface{}) interface{} {
	// params[0]: base object
	// params[1]: String object
	obj := params[0].(*object.Object)
	fld := obj.FieldTable["value"]
	str := object.GoStringFromStringObject(params[1].(*object.Object))
	var zz = new(big.Int)
	_, ok := zz.SetString(str, 10)
	if !ok {
		errMsg := fmt.Sprintf("bigIntegerInitString: string (%s) not all numerics", str)
		return getGErrBlk(excNames.NumberFormatException, errMsg)
	}

	// Update base object and return nil
	fld = object.Field{Ftype: types.BigInteger, Fvalue: zz}
	obj.FieldTable["value"] = fld

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld = object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return nil
}

// "java/math/BigInteger.<init>(Ljava/lang/String;I)V"
func bigIntegerInitStringRadix(params []interface{}) interface{} {
	// params[0]: base object
	// params[1]: String object
	// params[2]: radix int64
	obj := params[0].(*object.Object)
	fld := obj.FieldTable["value"]
	str := object.GoStringFromStringObject(params[1].(*object.Object))
	rdx := params[2].(int64)
	var zz = new(big.Int)
	_, ok := zz.SetString(str, int(rdx))
	if !ok {
		errMsg := fmt.Sprintf("bigIntegerInitStringRadix: string (%s) not all numerics or the radix (%d) is invalid", str, rdx)
		return getGErrBlk(excNames.NumberFormatException, errMsg)
	}

	// Update base object and return nil
	fld = object.Field{Ftype: types.BigInteger, Fvalue: zz}
	obj.FieldTable["value"] = fld

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld = object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return nil
}

// "java/math/BigInteger.abs()Ljava/math/BigInteger;"
func bigIntegerAbs(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// zz = abs(xx)

	objIn := params[0].(*object.Object)
	fld := objIn.FieldTable["value"]
	xx := fld.Fvalue.(*big.Int)

	// BigInteger operation
	var zz = new(big.Int)
	zz.Abs(xx)

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld = object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.add(Ljava/math/BigInteger;)Ljava/math/BigInteger;"
func bigIntegerAdd(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: argument object (yy)
	// zz = xx + yy

	objBase := params[0].(*object.Object)
	objArg := params[1].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	yy := objArg.FieldTable["value"].Fvalue.(*big.Int)

	// BigInteger operation
	var zz = new(big.Int)
	zz.Add(xx, yy)

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.and(Ljava/math/BigInteger;)Ljava/math/BigInteger;"
func bigIntegerAnd(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: argument object (yy)
	// zz = xx && yy

	objBase := params[0].(*object.Object)
	objArg := params[1].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	yy := objArg.FieldTable["value"].Fvalue.(*big.Int)

	// BigInteger operation
	var zz = new(big.Int)
	zz.And(xx, yy)

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.andNot(Ljava/math/BigInteger;)Ljava/math/BigInteger;"
func bigIntegerAndNot(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: argument object (yy)
	// zz = xx && ~yy

	objBase := params[0].(*object.Object)
	objArg := params[1].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	yy := objArg.FieldTable["value"].Fvalue.(*big.Int)

	// BigInteger operation
	var zz = new(big.Int)
	zz.AndNot(xx, yy)

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.bitCount()I"
func bigIntegerBitCount(params []interface{}) interface{} {
	// params[0]: base object (xx)

	obj := params[0].(*object.Object)
	fld := obj.FieldTable["value"]
	xx := fld.Fvalue.(*big.Int)
	var count int
	for _, wd := range xx.Bits() {
		count += bits.OnesCount(uint(wd))
	}
	return int64(count)

}

// "java/math/BigInteger.bitLength()I"
func bigIntegerBitLength(params []interface{}) interface{} {
	// params[0]: base object (xx)

	obj := params[0].(*object.Object)
	fld := obj.FieldTable["value"]
	xx := fld.Fvalue.(*big.Int)
	return int64(xx.BitLen())

}

// "java/math/BigInteger.byteValueExact()B"
func bigIntegerByteValueExact(params []interface{}) interface{} {
	// params[0]: base object (xx)

	obj := params[0].(*object.Object)
	fld := obj.FieldTable["value"]
	xx := fld.Fvalue.(*big.Int)
	ii := xx.Int64()
	if ii < 0 || ii > 255 {
		errMsg := fmt.Sprintf("bigIntegerByteValueExact: Value (%d) will not fit in a byte", ii)
		return getGErrBlk(excNames.ArithmeticException, errMsg)
	}
	return ii & 0xFF

}

// "java/math/BigInteger.compareTo(Ljava/math/BigInteger;)I"
func bigIntegerCompareTo(params []interface{}) interface{} {
	// params[0]:  base object (xx)
	// params[1]:  argument object (yy)
	objBase := params[0].(*object.Object)
	objArg := params[1].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	yy := objArg.FieldTable["value"].Fvalue.(*big.Int)
	return int64(xx.Cmp(yy))
}

// "java/math/BigInteger.divide(Ljava/math/BigInteger;)Ljava/math/BigInteger;"
func bigIntegerDivide(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: argument object (yy)
	// zz = xx / yy

	objBase := params[0].(*object.Object)
	objArg := params[1].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	yy := objArg.FieldTable["value"].Fvalue.(*big.Int)
	zero := big.NewInt(int64(0))
	if yy.Cmp(zero) <= 0 {
		errMsg := "bigIntegerDivide: Divide by zero"
		return getGErrBlk(excNames.ArithmeticException, errMsg)
	}

	// BigInteger operation
	var zz = new(big.Int)
	zz.Div(xx, yy)

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.divideAndRemainder(Ljava/math/BigInteger;)[Ljava/math/BigInteger;"
func bigIntegerDivideAndRemainder(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: argument object (yy)
	// zz = xx / yy; rr = xx % yy

	objBase := params[0].(*object.Object)
	objArg := params[1].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	yy := objArg.FieldTable["value"].Fvalue.(*big.Int)
	zero := big.NewInt(int64(0))
	if yy.Cmp(zero) <= 0 {
		errMsg := "bigIntegerDivideAndRemainder: Divide by zero"
		return getGErrBlk(excNames.ArithmeticException, errMsg)
	}

	// BigInteger operation
	var zz = new(big.Int)
	var rr = new(big.Int)
	zz.Div(xx, yy)
	rr.Rem(xx, yy)

	// Create xx / yy and xx % yy objects
	obj1 := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)
	obj2 := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, rr)

	// Create the return object with the object-array
	var objectArray = []*object.Object{obj1, obj2}
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, objectArray)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.doubleValue()J"
func bigIntegerFloat64Value(params []interface{}) interface{} {
	// params[0]:  base object (xx)

	objBase := params[0].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	outDouble := float64(xx.Int64())

	return outDouble
}

// "java/math/BigInteger.equals(Ljava/math/BigInteger;)Z"
func bigIntegerEquals(params []interface{}) interface{} {
	// params[0]:  base object (xx)
	// params[1]:  argument object (yy)
	objBase := params[0].(*object.Object)
	objArg := params[1].(*object.Object)
	if objArg.FieldTable["value"].Ftype != types.BigInteger {
		return int64(0)
	}
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	yy := objArg.FieldTable["value"].Fvalue.(*big.Int)
	if xx.Cmp(yy) != 0 {
		return int64(0)
	}
	return int64(1)
}

// "java/math/BigInteger.gcd(Ljava/math/BigInteger;)Ljava/math/BigInteger;"
func bigIntegerGCD(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: argument object (yy)
	// zz = GCD of xx and yy

	objBase := params[0].(*object.Object)
	objArg := params[1].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	yy := objArg.FieldTable["value"].Fvalue.(*big.Int)

	// BigInteger operation
	var zz = new(big.Int)
	zz.GCD(nil, nil, xx, yy)

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.intValue()I"
// "java/math/BigInteger.intValueExact()I"
// "java/math/BigInteger.longValue()J"
// "java/math/BigInteger.longValueExact()J"
// "java/math/BigInteger.shortValueExact()S"
func bigIntegerInt64Value(params []interface{}) interface{} {
	// params[0]:  base object (xx)

	objBase := params[0].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	outInt64 := xx.Int64()

	return outInt64
}

// "java/math/BigInteger.isProbablePrime(I)Z"
func bigIntegerIsProbablePrime(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: certainty integer
	// Ref: https://docs.oracle.com/en/java/javase/17/docs/api/java.base/java/math/BigInteger.html#isProbablePrime(int)

	baseObj := params[0].(*object.Object)
	xx := baseObj.FieldTable["value"].Fvalue.(*big.Int)
	certaintyInt64 := params[1].(int64)
	if xx.ProbablyPrime(int(certaintyInt64)) {
		return int64(1)
	}

	return int64(0)
}

// "java/math/BigInteger.max(Ljava/math/BigInteger;)Ljava/math/BigInteger;"
func bigIntegerMax(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: argument object (yy)
	// zz = max(xx, yy)

	objBase := params[0].(*object.Object)
	objArg := params[1].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	yy := objArg.FieldTable["value"].Fvalue.(*big.Int)

	// BigInteger operation
	var zz = new(big.Int)
	if xx.Cmp(yy) > 0 {
		zz = xx
	} else {
		zz = yy
	}

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.min(Ljava/math/BigInteger;)Ljava/math/BigInteger;"
func bigIntegerMin(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: argument object (yy)
	// zz = min(xx, yy)

	objBase := params[0].(*object.Object)
	objArg := params[1].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	yy := objArg.FieldTable["value"].Fvalue.(*big.Int)

	// BigInteger operation
	var zz = new(big.Int)
	if xx.Cmp(yy) < 0 {
		zz = xx
	} else {
		zz = yy
	}

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.mod(Ljava/math/BigInteger;)Ljava/math/BigInteger;"
func bigIntegerMod(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: argument object (yy)
	// zz = xx mod yy

	objBase := params[0].(*object.Object)
	objArg := params[1].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	yy := objArg.FieldTable["value"].Fvalue.(*big.Int)
	zero := big.NewInt(int64(0))
	if yy.Cmp(zero) <= 0 {
		errMsg := "bigIntegerMod: modulus not positive"
		return getGErrBlk(excNames.ArithmeticException, errMsg)
	}

	// BigInteger operation
	var zz = new(big.Int)
	zz.Mod(xx, yy)

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.modInverse(Ljava/math/BigInteger;)Ljava/math/BigInteger;"
// The modInverse() method returns the modular multiplicative inverse of the base object, modulo the argument.
// Note that zz = the modular multiplicative inverse of (xx % mm) is such that
// (xx * zz) % mm = 1.
//
// This method throws an ArithmeticException if modulus <= 0
// or this has no multiplicative inverse modulo the modulus.
func bigIntegerModInverse(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: modulus object (mm)

	objBase := params[0].(*object.Object)
	objModulus := params[1].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	mm := objModulus.FieldTable["value"].Fvalue.(*big.Int)
	zero := big.NewInt(int64(0))
	if mm.Cmp(zero) <= 0 {
		errMsg := "bigIntegerModInverse: modulus not positive"
		return getGErrBlk(excNames.ArithmeticException, errMsg)
	}

	// BigInteger operation
	var zz = new(big.Int)
	ret := zz.ModInverse(xx, mm)
	if ret == nil {
		errMsg := "bigIntegerModInverse: BigInteger not invertible"
		return getGErrBlk(excNames.ArithmeticException, errMsg)
	}

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.modPow(Ljava/math/BigInteger;Ljava/math/BigInteger;)Ljava/math/BigInteger;"
// Compute a BigInteger whose value is (bb ^ ee modulo mm )
func bigIntegerModPow(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: exponent object (ee)
	// params[2]: modulus object (mm)
	// zz = (xx ** ee) modulo mm

	objBB := params[0].(*object.Object)
	objEE := params[1].(*object.Object)
	objMM := params[2].(*object.Object)
	xx := objBB.FieldTable["value"].Fvalue.(*big.Int)
	ee := objEE.FieldTable["value"].Fvalue.(*big.Int)
	mm := objMM.FieldTable["value"].Fvalue.(*big.Int)
	zero := big.NewInt(int64(0))
	if mm.Cmp(zero) <= 0 {
		errMsg := fmt.Sprintf("bigIntegerModPow: Modulus (%d) is negative", mm.Int64())
		return getGErrBlk(excNames.ArithmeticException, errMsg)
	}

	// BigInteger operation
	var zz = new(big.Int)
	zz.Exp(xx, ee, mm)

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.multiply(Ljava/math/BigInteger;)Ljava/math/BigInteger;"
func bigIntegerMultiply(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: argument object (yy)
	// zz = xx * yy

	objBase := params[0].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)

	// yy = BigInteger argument.
	var yy *big.Int
	switch params[1].(type) {
	case int64:
		argLong := params[1].(int64)
		yy = big.NewInt(argLong)
	default: // BigInteger object
		objArg := params[1].(*object.Object)
		yy = objArg.FieldTable["value"].Fvalue.(*big.Int)
	}

	// BigInteger operation
	var zz = new(big.Int)
	zz.Mul(xx, yy)

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.negate()Ljava/math/BigInteger;"
func bigIntegerNegate(params []interface{}) interface{} {
	// params[0]:  base object (xx)
	// zz = -xx

	objBase := params[0].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)

	// BigInteger operation
	var zz = new(big.Int)
	zz.Neg(xx)

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.negate()Ljava/math/BigInteger;"
func bigIntegerNot(params []interface{}) interface{} {
	// params[0]:  base object (xx)
	// zz = not xx

	objBase := params[0].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)

	// BigInteger operation
	var zz = new(big.Int)
	zz.Not(xx)

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.xor(Ljava/math/BigInteger;)Ljava/math/BigInteger;"
func bigIntegerOr(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: argument object (yy)
	// zz = xx XOR yy

	objBase := params[0].(*object.Object)
	objArg := params[1].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	yy := objArg.FieldTable["value"].Fvalue.(*big.Int)

	// BigInteger operation
	var zz = new(big.Int)
	zz.Or(xx, yy)

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.pow(I)Ljava/math/BigInteger;"
func bigIntegerPow(params []interface{}) interface{} {
	// params[0]:  base object (xx)
	// params[1] = int64 power (pow)
	// zz = xx ** pow

	objBase := params[0].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	pow := params[1].(int64)
	if pow < 0 {
		errMsg := fmt.Sprintf("bigIntegerPow: Power (%d) is negative", pow)
		return getGErrBlk(excNames.ArithmeticException, errMsg)
	}
	yy := big.NewInt(pow)

	// BigInteger operation
	var zz = new(big.Int)
	zz.Exp(xx, yy, nil)

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.remainder(Ljava/math/BigInteger;)Ljava/math/BigInteger;"
func bigIntegerRemainder(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: argument object (yy)
	// zz = remainder when dividing xx by yy

	objBase := params[0].(*object.Object)
	objArg := params[1].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	yy := objArg.FieldTable["value"].Fvalue.(*big.Int)
	zero := big.NewInt(int64(0))
	if yy.Cmp(zero) <= 0 {
		errMsg := "bigIntegerRemainder: Divide by zero"
		return getGErrBlk(excNames.ArithmeticException, errMsg)
	}

	// BigInteger operation
	var zz = new(big.Int)
	zz.Rem(xx, yy)

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.probablyPrime(ILjava/util/Random;)Ljava/math/BigInteger;"
func bigIntegerProbablyPrime(params []interface{}) interface{} {
	// params[0]: number of bits (bitLength)
	// params[1]: Random object (yy)

	bitLength := params[0].(int64)
	// TODO: Ignored for now: objRandom := params[2].(*object.Object)

	zz, errMsg := getPrime(int(bitLength))
	if zz != nil {
		// Create return object.
		obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

		// Set signum field to the sign.
		signum := int64(zz.Sign())
		fld := object.Field{Ftype: types.Int, Fvalue: signum}
		obj.FieldTable["signum"] = fld

		return obj
	}
	return getGErrBlk(excNames.ArithmeticException, errMsg)
}

// "java/math/BigInteger.signum()I"
func bigIntegerSignum(params []interface{}) interface{} {
	// params[0]:  base object (xx)

	objBase := params[0].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	return int64(xx.Sign())
}

// "java/math/BigInteger.sqrt()Ljava/math/BigInteger;"
func bigIntegerSqrt(params []interface{}) interface{} {
	// params[0]:  base object (xx)

	objBase := params[0].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	zero := big.NewInt(int64(0))
	if xx.Cmp(zero) < 0 {
		errMsg := fmt.Sprintf("bigIntegerSqrt: Argument (%d) is negative", xx.Int64())
		return getGErrBlk(excNames.ArithmeticException, errMsg)
	}

	// BigInteger operation
	var zz = new(big.Int)
	zz.Sqrt(xx)

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.subtract(Ljava/math/BigInteger;)Ljava/math/BigInteger;"
func bigIntegerSubtract(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: argument object (yy)
	// zz = xx - yy

	objBase := params[0].(*object.Object)
	objArg := params[1].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	yy := objArg.FieldTable["value"].Fvalue.(*big.Int)

	// BigInteger operation
	var zz = new(big.Int)
	zz.Sub(xx, yy)

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.toByteArray()[B"
func bigIntegerToByteArray(params []interface{}) interface{} {
	// params[0]: base object (xx)

	obj := params[0].(*object.Object)
	xx := obj.FieldTable["value"].Fvalue.(*big.Int)
	bytes := xx.Bytes()
	objOut :=
		object.StringObjectFromJavaByteArray(object.JavaByteArrayFromGoByteArray(bytes))

	return objOut
}

// "java/math/BigInteger.toString()Ljava/lang/String;"
func bigIntegerToString(params []interface{}) interface{} {
	// params[0]:  base object (xx)

	obj := params[0].(*object.Object)
	xx := obj.FieldTable["value"].Fvalue.(*big.Int)
	str := xx.String()
	objOut := object.StringObjectFromGoString(str)

	return objOut
}

// "java/math/BigInteger.toString(I)Ljava/lang/String;"
func bigIntegerToStringRadix(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: radix int64 (rdx)

	objBase := params[0].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	rdx := params[1].(int64)
	if rdx < 2 || rdx > 62 {
		errMsg := fmt.Sprintf("bigIntegerToStringRadix: Invalid radix value (%d)", rdx)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	str := xx.Text(int(rdx))
	objOut := object.StringObjectFromGoString(str)

	return objOut
}

// "java/math/BigInteger.valueOf(J)Ljava/math/BigInteger;"
func bigIntegerValueOf(params []interface{}) interface{} {
	// params[0]:  long value for returned big.Int object

	argValue := params[0].(int64)
	obj := object.MakeEmptyObjectWithClassName(&classNameBigInteger)
	InitBigIntegerField(obj, argValue)

	return obj
}

// "java/math/BigInteger.xor(Ljava/math/BigInteger;)Ljava/math/BigInteger;"
func bigIntegerXor(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: argument object (yy)
	// zz = xx XOR yy

	objBase := params[0].(*object.Object)
	objArg := params[1].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	yy := objArg.FieldTable["value"].Fvalue.(*big.Int)

	// BigInteger operation
	var zz = new(big.Int)
	zz.Xor(xx, yy)

	// Create return object
	obj := object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.testBit(I)Z"
func bigIntegerTestBit(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	bitN := int(params[1].(int64))
	bi := obj.FieldTable["value"].Fvalue.(*big.Int)
	if bi.Bit(bitN) == 1 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// "java/math/BigInteger.setBit(I)Ljava/math/BigInteger;"
func bigIntegerSetBit(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	bitN := int(params[1].(int64))
	fld := obj.FieldTable["value"]
	bigInt := fld.Fvalue.(*big.Int)
	newBigInt := new(big.Int).Set(bigInt)
	newBigInt.SetBit(newBigInt, bitN, 1)
	return object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, newBigInt)
}

// "java/math/BigInteger.shiftLeft(I)Ljava/math/BigInteger;"
func bigIntegerShiftLeft(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	bitN := uint(params[1].(int64))
	fld := obj.FieldTable["value"]
	bigInt := fld.Fvalue.(*big.Int)
	newBigInt := new(big.Int).Lsh(bigInt, bitN)
	return object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, newBigInt)
}

// "java/math/BigInteger.shiftLeft(I)Ljava/math/BigInteger;"
func bigIntegerShiftRight(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	bitN := uint(params[1].(int64))
	fld := obj.FieldTable["value"]
	bigInt := fld.Fvalue.(*big.Int)
	newBigInt := new(big.Int).Rsh(bigInt, bitN)
	return object.MakePrimitiveObject(classNameBigInteger, types.BigInteger, newBigInt)
}
