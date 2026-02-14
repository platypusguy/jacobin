/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaMath

import (
	"crypto/rand"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"math/big"
	"math/bits"
)

/*
The BigInteger object is implemented using Golang package math/big.

*/

func Load_Math_Big_Integer() {

	ghelpers.MethodSignatures["java/math/BigInteger.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerClinit,
		}

	// <init> functions

	ghelpers.MethodSignatures["java/math/BigInteger.<init>([B)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerInitByteArray,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.<init>([BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.<init>(I[B)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.<init>(I[BII)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.<init>(IILjava/util/Random;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  bigIntegerInitProbablyPrime,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.<init>(ILjava/util/Random;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigIntegerInitRandom,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerInitString,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.<init>(Ljava/lang/String;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigIntegerInitStringRadix,
		}

	// Member functions

	ghelpers.MethodSignatures["java/math/BigInteger.abs()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerAbs,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.add(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerAdd,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.and(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerAnd,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.andNot(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerAndNot,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.bitCount()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerBitCount,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.bitLength()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerBitLength,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.byteValueExact()B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerByteValueExact,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.clearBit(I)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.compareTo(Ljava/math/BigInteger;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerCompareTo,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.divide(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerDivide,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.divideAndRemainder(Ljava/math/BigInteger;)[Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerDivideAndRemainder,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.doubleValue()D"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerFloat64Value,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerEquals,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.flipBit(I)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.floatValue()F"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerFloat64Value,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.gcd(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerGCD,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.getLowestSetBit()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerHashCode,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.intValue()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerInt64Value,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.intValueExact()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerInt64Value,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.isProbablePrime(I)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerIsProbablePrime,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.longValue()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerInt64Value,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.longValueExact()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerInt64Value,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.max(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerMax,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.min(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerMin,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.mod(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerMod,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.modInverse(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerModInverse,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.modPow(Ljava/math/BigInteger;Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigIntegerModPow,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.multiply(J)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerMultiply,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.multiply(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerMultiply,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.negate()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerNegate,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.nextProbablePrime()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.not()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerNot,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.or(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerOr,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.pow(I)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerPow,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.probablePrime(ILjava/util/Random;)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigIntegerProbablyPrime,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.remainder(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerRemainder,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.setBit(I)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerSetBit,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.shiftLeft(I)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerShiftLeft,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.shiftRight(I)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerShiftRight,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.shortValueExact()S"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerInt64Value,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.signum()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerSignum,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.sqrt()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerSqrt,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.sqrtAndRemainder()[Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.subtract(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerSubtract,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.testBit(I)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerTestBit,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.toByteArray()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigIntegerToByteArray,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  BigIntegerToString,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.toString(I)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerToStringRadix,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.valueOf(J)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerValueOf,
		}

	ghelpers.MethodSignatures["java/math/BigInteger.xor(Ljava/math/BigInteger;)Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigIntegerXor,
		}

}

// Get a prime number formatted as a big.Int.
// Uses crypto/rand.Prime which already returns a probable prime,
// avoiding any custom unbounded loop here.
func getPrime(bitLength int) (*big.Int, string) {
	zz, err := rand.Prime(rand.Reader, bitLength)
	if err != nil {
		errMsg := fmt.Sprintf("getPrime: rand.Reader(bitLength=%d) failed, reason: %s", bitLength, err.Error())
		return nil, errMsg
	}
	return zz, ""
}

// addStaticBigInteger: Form a BigInteger object based on the parameter value.
func addStaticBigInteger(argName string, argValue int64) {
	name := fmt.Sprintf("%s.%s", types.ClassNameBigInteger, argName)
	obj := object.MakeEmptyObjectWithClassName(&types.ClassNameBigInteger)
	ghelpers.InitBigIntegerField(obj, argValue)
	_ = statics.AddStatic(name, statics.Static{Type: "Ljava/math/BigInteger;", Value: obj})
}

// "java/math/BigInteger.<clinit>()V"
func bigIntegerClinit([]interface{}) interface{} {
	addStaticBigInteger("ONE", int64(1))
	addStaticBigInteger("TEN", int64(10))
	addStaticBigInteger("TWO", int64(2))
	addStaticBigInteger("ZERO", int64(0))
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
	object.ClearFieldTable(obj)
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
	object.ClearFieldTable(obj)
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
	return ghelpers.GetGErrBlk(excNames.ArithmeticException, errMsg)
}

// "java/math/BigInteger.<init>(IILjava/util/Random;)V"
func bigIntegerInitRandom(params []interface{}) interface{} {
	// params[0]: base object
	// params[1]: int64 holding numbits such that the base object value field
	//            will be set to a random value in the rang given by [0 : 2**(numbits) - 1].
	// params[2]: Random object
	obj := params[0].(*object.Object)
	object.ClearFieldTable(obj)
	fld := obj.FieldTable["value"]
	numBits := params[1].(int64)
	// TODO: Ignored for now: objRandom := params[2].(*object.Object)

	// Compute upperBound = 2**(numBits) based on numBits.
	upperBound := new(big.Int).Lsh(big.NewInt(1), uint(numBits))

	// Get a big.Int in the randge of [0, upperBound].
	zz, err := rand.Int(rand.Reader, upperBound)
	if err != nil {
		errMsg := fmt.Sprintf("bigIntegerInitRandom: rand.Int(numBits=%d) failed, reason: %s", numBits, err.Error())
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, errMsg)
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
	object.ClearFieldTable(obj)
	fld := obj.FieldTable["value"]
	str := object.GoStringFromStringObject(params[1].(*object.Object))
	var zz = new(big.Int)
	_, ok := zz.SetString(str, 10)
	if !ok {
		errMsg := fmt.Sprintf("bigIntegerInitString: string (%s) not all numerics", str)
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, errMsg)
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
	object.ClearFieldTable(obj)
	fld := obj.FieldTable["value"]
	str := object.GoStringFromStringObject(params[1].(*object.Object))
	rdx := params[2].(int64)
	var zz = new(big.Int)
	_, ok := zz.SetString(str, int(rdx))
	if !ok {
		errMsg := fmt.Sprintf("bigIntegerInitStringRadix: string (%s) not all numerics or the radix (%d) is invalid", str, rdx)
		return ghelpers.GetGErrBlk(excNames.NumberFormatException, errMsg)
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
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

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

	objThis, ok := params[0].(*object.Object)
	if !ok || object.IsNull(objThis) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "bigIntegerAdd: this-object is null")
	}
	objArg, ok := params[1].(*object.Object)
	if !ok || object.IsNull(objThis) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "bigIntegerAdd: arg-object is null")
	}
	xx := objThis.FieldTable["value"].Fvalue.(*big.Int)
	yy := objArg.FieldTable["value"].Fvalue.(*big.Int)

	// BigInteger operation
	var zz = new(big.Int)
	zz.Add(xx, yy)

	// Create return object
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

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
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

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
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

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
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, errMsg)
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
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, errMsg)
	}

	// BigInteger operation
	var zz = new(big.Int)
	zz.Div(xx, yy)

	// Create return object
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

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
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, errMsg)
	}

	// BigInteger operation
	var zz = new(big.Int)
	var rr = new(big.Int)
	zz.Div(xx, yy)
	rr.Rem(xx, yy)

	// Create xx / yy and xx % yy objects
	obj1 := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)
	obj2 := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, rr)

	// Create the return object with the object-array
	var objectArray = []*object.Object{obj1, obj2}
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, objectArray)

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
	objArg, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "bigIntegerEquals: argument not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	if objArg.FieldTable["value"].Ftype != types.BigInteger {
		return types.JavaBoolFalse
	}
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	yy := objArg.FieldTable["value"].Fvalue.(*big.Int)
	if xx.Cmp(yy) != 0 {
		return types.JavaBoolFalse
	}
	return types.JavaBoolTrue
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
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// Compute the hash code based on the big.Int value.
func bigIntegerHashCode(params []interface{}) interface{} {
	objBase := params[0].(*object.Object)
	bi := objBase.FieldTable["value"].Fvalue.(*big.Int)

	// Java's BigInteger.hashCode() uses mag array which is big-endian.
	// big.Int.Bytes() returns big-endian bytes.
	bytes := bi.Bytes()

	// Convert bytes to int32 array (mag array in Java)
	var mag []int32
	// If the byte array length is not a multiple of 4, we need to handle the first few bytes.
	firstIntLen := len(bytes) % 4
	if firstIntLen > 0 {
		var firstInt int32
		for i := 0; i < firstIntLen; i++ {
			firstInt = (firstInt << 8) | int32(bytes[i])
		}
		mag = append(mag, firstInt)
	}

	for i := firstIntLen; i < len(bytes); i += 4 {
		var val int32
		for j := 0; j < 4; j++ {
			val = (val << 8) | int32(bytes[i+j])
		}
		mag = append(mag, val)
	}

	var hashCode int32 = 0
	for _, val := range mag {
		hashCode = 31*hashCode + val
	}

	return int64(hashCode * int32(bi.Sign()))
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
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

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
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

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
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, errMsg)
	}

	// BigInteger operation
	var zz = new(big.Int)
	zz.Mod(xx, yy)

	// Create return object
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

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
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, errMsg)
	}

	// BigInteger operation
	var zz = new(big.Int)
	ret := zz.ModInverse(xx, mm)
	if ret == nil {
		errMsg := "bigIntegerModInverse: BigInteger not invertible"
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, errMsg)
	}

	// Create return object
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

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
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, errMsg)
	}

	// BigInteger operation
	var zz = new(big.Int)
	zz.Exp(xx, ee, mm)

	// Create return object
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

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
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

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
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.not()Ljava/math/BigInteger;"
func bigIntegerNot(params []interface{}) interface{} {
	// params[0]:  base object (xx)
	// zz = not xx

	objBase := params[0].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)

	// BigInteger operation
	var zz = new(big.Int)
	zz.Not(xx)

	// Create return object
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

	// Set signum field to the sign.
	signum := int64(zz.Sign())
	fld := object.Field{Ftype: types.Int, Fvalue: signum}
	obj.FieldTable["signum"] = fld

	return obj
}

// "java/math/BigInteger.or(Ljava/math/BigInteger;)Ljava/math/BigInteger;"
func bigIntegerOr(params []interface{}) interface{} {
	// params[0]: base object (xx)
	// params[1]: argument object (yy)
	// zz = xx OR yy

	objBase := params[0].(*object.Object)
	objArg := params[1].(*object.Object)
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	yy := objArg.FieldTable["value"].Fvalue.(*big.Int)

	// BigInteger operation
	var zz = new(big.Int)
	zz.Or(xx, yy)

	// Create return object
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

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

	objBase, ok := params[0].(*object.Object)
	if !ok || object.IsNull(objBase) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "bigIntegerPow: this-object is null")
	}
	xx := objBase.FieldTable["value"].Fvalue.(*big.Int)
	pow := params[1].(int64)
	if pow < 0 {
		errMsg := fmt.Sprintf("bigIntegerPow: Power (%d) is negative", pow)
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, errMsg)
	}
	yy := big.NewInt(pow)

	// BigInteger operation
	var zz = new(big.Int)
	zz.Exp(xx, yy, nil)

	// Create return object
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

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
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, errMsg)
	}

	// BigInteger operation
	var zz = new(big.Int)
	zz.Rem(xx, yy)

	// Create return object
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

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
		obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

		// Set signum field to the sign.
		signum := int64(zz.Sign())
		fld := object.Field{Ftype: types.Int, Fvalue: signum}
		obj.FieldTable["signum"] = fld

		return obj
	}
	return ghelpers.GetGErrBlk(excNames.ArithmeticException, errMsg)
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
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, errMsg)
	}

	// BigInteger operation
	var zz = new(big.Int)
	zz.Sqrt(xx)

	// Create return object
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

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
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

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
func BigIntegerToString(params []interface{}) interface{} {
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
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	str := xx.Text(int(rdx))
	objOut := object.StringObjectFromGoString(str)

	return objOut
}

// "java/math/BigInteger.valueOf(J)Ljava/math/BigInteger;"
func bigIntegerValueOf(params []interface{}) interface{} {
	// params[0]:  long value for returned big.Int object

	argValue := params[0].(int64)
	obj := object.MakeEmptyObjectWithClassName(&types.ClassNameBigInteger)
	ghelpers.InitBigIntegerField(obj, argValue)

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
	obj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, zz)

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
	return object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, newBigInt)
}

// "java/math/BigInteger.shiftLeft(I)Ljava/math/BigInteger;"
func bigIntegerShiftLeft(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	bitN := uint(params[1].(int64))
	fld := obj.FieldTable["value"]
	bigInt := fld.Fvalue.(*big.Int)
	newBigInt := new(big.Int).Lsh(bigInt, bitN)
	return object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, newBigInt)
}

// "java/math/BigInteger.shiftLeft(I)Ljava/math/BigInteger;"
func bigIntegerShiftRight(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	bitN := uint(params[1].(int64))
	fld := obj.FieldTable["value"]
	bigInt := fld.Fvalue.(*big.Int)
	newBigInt := new(big.Int).Rsh(bigInt, bitN)
	return object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, newBigInt)
}
