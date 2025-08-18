/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/object"
	"math"
	"math/big"
	"math/rand"
)

/*
 Each object or library that has Go methods contains a reference to MethodSignatures,
 which contain data needed to insert the go method into the MTable of the currently
 executing JVM. MethodSignatures is a map whose key is the fully qualified name and
 type of the method (that is, the method's full signature) and a value consisting of
 a struct of an int (the number of slots to pop off the caller's operand stack when
 creating the new frame and a function. All methods have the same signature, regardless
 of the signature of their Java counterparts. That signature is that it accepts a slice
 of interface{} and returns an interface{}. The accepted slice can be empty and the
 return interface can be nil. This covers all Java functions. (Objects are returned
 as a 64-bit address in this scheme (as they are in the JVM).

 The passed-in slice contains one entry for every parameter passed to the method (which
 could mean an empty slice).
*/

/*
Math constant references:

	Java: https://docs.oracle.com/en/java/javase/17/docs/api/constant-values.html#java.math
	Golang: https://pkg.go.dev/math#pkg-constants

*/

const MAX_DOUBLE_EXPONENT = 1023
const PI = 3.14159265358979323846

func Load_Lang_Math() {

	MethodSignatures["java/lang/Math.abs(D)D"] = GMeth{ParamSlots: 1, GFunction: absFloat64}
	MethodSignatures["java/lang/Math.abs(F)F"] = GMeth{ParamSlots: 1, GFunction: absFloat64}
	MethodSignatures["java/lang/Math.abs(I)I"] = GMeth{ParamSlots: 1, GFunction: absInt64}
	MethodSignatures["java/lang/Math.abs(J)J"] = GMeth{ParamSlots: 1, GFunction: absInt64}
	MethodSignatures["java/lang/Math.absExact(I)I"] = GMeth{ParamSlots: 1, GFunction: absInt64}
	MethodSignatures["java/lang/Math.absExact(J)J"] = GMeth{ParamSlots: 1, GFunction: absInt64}
	MethodSignatures["java/lang/Math.acos(D)D"] = GMeth{ParamSlots: 1, GFunction: acosFloat64}
	MethodSignatures["java/lang/Math.addExact(II)I"] = GMeth{ParamSlots: 2, GFunction: addExactII}
	MethodSignatures["java/lang/Math.addExact(JJ)J"] = GMeth{ParamSlots: 2, GFunction: addExactJJ}
	MethodSignatures["java/lang/Math.asin(D)D"] = GMeth{ParamSlots: 1, GFunction: asinFloat64}
	MethodSignatures["java/lang/Math.atan(D)D"] = GMeth{ParamSlots: 1, GFunction: atanFloat64}
	MethodSignatures["java/lang/Math.atan2(DD)D"] = GMeth{ParamSlots: 2, GFunction: atan2Float64}
	MethodSignatures["java/lang/Math.cbrt(D)D"] = GMeth{ParamSlots: 1, GFunction: cbrtFloat64}
	MethodSignatures["java/lang/Math.ceil(D)D"] = GMeth{ParamSlots: 1, GFunction: ceilFloat64}
	MethodSignatures["java/lang/Math.copySign(DD)D"] = GMeth{ParamSlots: 2, GFunction: copySignDD}
	MethodSignatures["java/lang/Math.copySign(FF)F"] = GMeth{ParamSlots: 2, GFunction: copySignFF}
	MethodSignatures["java/lang/Math.cos(D)D"] = GMeth{ParamSlots: 1, GFunction: cosFloat64}
	MethodSignatures["java/lang/Math.cosh(D)D"] = GMeth{ParamSlots: 1, GFunction: coshFloat64}
	MethodSignatures["java/lang/Math.decrementExact(I)I"] = GMeth{ParamSlots: 1, GFunction: decrementExactInt64}
	MethodSignatures["java/lang/Math.decrementExact(J)J"] = GMeth{ParamSlots: 1, GFunction: decrementExactInt64}
	MethodSignatures["java/lang/Math.exp(D)D"] = GMeth{ParamSlots: 1, GFunction: expFloat64}
	MethodSignatures["java/lang/Math.expm1(D)D"] = GMeth{ParamSlots: 1, GFunction: expm1Float64}
	MethodSignatures["java/lang/Math.floor(D)D"] = GMeth{ParamSlots: 1, GFunction: floorFloat64}
	MethodSignatures["java/lang/Math.floorDiv(II)I"] = GMeth{ParamSlots: 2, GFunction: floorDivII}
	MethodSignatures["java/lang/Math.floorDiv(JI)J"] = GMeth{ParamSlots: 2, GFunction: floorDivJx}
	MethodSignatures["java/lang/Math.floorDiv(JJ)J"] = GMeth{ParamSlots: 2, GFunction: floorDivJx}
	MethodSignatures["java/lang/Math.floorMod(II)I"] = GMeth{ParamSlots: 2, GFunction: floorModII}
	MethodSignatures["java/lang/Math.floorMod(JI)I"] = GMeth{ParamSlots: 2, GFunction: floorModJx}
	MethodSignatures["java/lang/Math.floorMod(JJ)J"] = GMeth{ParamSlots: 2, GFunction: floorModJx}
	MethodSignatures["java/lang/Math.fma(DDD)D"] = GMeth{ParamSlots: 3, GFunction: fmaDDD}
	MethodSignatures["java/lang/Math.fma(FFF)F"] = GMeth{ParamSlots: 3, GFunction: fmaFFF}
	MethodSignatures["java/lang/Math.getExponent(D)I"] = GMeth{ParamSlots: 1, GFunction: getExponentFloat64}
	MethodSignatures["java/lang/Math.getExponent(F)I"] = GMeth{ParamSlots: 1, GFunction: getExponentFloat64}
	MethodSignatures["java/lang/Math.hypot(DD)D"] = GMeth{ParamSlots: 2, GFunction: hypotFloat64}
	MethodSignatures["java/lang/Math.IEEEremainder(DD)D"] = GMeth{ParamSlots: 2, GFunction: IEEEremainderFloat64}
	MethodSignatures["java/lang/Math.incrementExact(I)I"] = GMeth{ParamSlots: 1, GFunction: incrementExactInt64}
	MethodSignatures["java/lang/Math.incrementExact(J)J"] = GMeth{ParamSlots: 1, GFunction: incrementExactInt64}
	MethodSignatures["java/lang/Math.log(D)D"] = GMeth{ParamSlots: 1, GFunction: logFloat64}
	MethodSignatures["java/lang/Math.log10(D)D"] = GMeth{ParamSlots: 1, GFunction: log10Float64}
	MethodSignatures["java/lang/Math.log1p(D)D"] = GMeth{ParamSlots: 1, GFunction: log1pFloat64}
	MethodSignatures["java/lang/Math.max(DD)D"] = GMeth{ParamSlots: 2, GFunction: maxDD}
	MethodSignatures["java/lang/Math.max(FF)F"] = GMeth{ParamSlots: 2, GFunction: maxFF}
	MethodSignatures["java/lang/Math.max(II)I"] = GMeth{ParamSlots: 2, GFunction: maxII}
	MethodSignatures["java/lang/Math.max(JJ)J"] = GMeth{ParamSlots: 2, GFunction: maxJJ}
	MethodSignatures["java/lang/Math.min(DD)D"] = GMeth{ParamSlots: 2, GFunction: minDD}
	MethodSignatures["java/lang/Math.min(FF)F"] = GMeth{ParamSlots: 2, GFunction: minFF}
	MethodSignatures["java/lang/Math.min(II)I"] = GMeth{ParamSlots: 2, GFunction: minII}
	MethodSignatures["java/lang/Math.min(JJ)J"] = GMeth{ParamSlots: 2, GFunction: minJJ}
	MethodSignatures["java/lang/Math.multiplyExact(II)I"] = GMeth{ParamSlots: 2, GFunction: multiplyExactII}
	MethodSignatures["java/lang/Math.multiplyExact(JI)I"] = GMeth{ParamSlots: 2, GFunction: multiplyExactJx}
	MethodSignatures["java/lang/Math.multiplyExact(JJ)J"] = GMeth{ParamSlots: 2, GFunction: multiplyExactJx}
	MethodSignatures["java/lang/Math.multiplyHigh(JJ)J"] = GMeth{ParamSlots: 2, GFunction: multiplyHighJJ}
	MethodSignatures["java/lang/Math.negateExact(I)I"] = GMeth{ParamSlots: 1, GFunction: negateExactInt64}
	MethodSignatures["java/lang/Math.negateExact(J)J"] = GMeth{ParamSlots: 1, GFunction: negateExactInt64}
	MethodSignatures["java/lang/Math.nextAfter(DD)D"] = GMeth{ParamSlots: 2, GFunction: nextAfterDD}
	MethodSignatures["java/lang/Math.nextAfter(FD)F"] = GMeth{ParamSlots: 2, GFunction: nextAfterFD}
	MethodSignatures["java/lang/Math.nextDown(D)D"] = GMeth{ParamSlots: 1, GFunction: nextDownFloat64}
	MethodSignatures["java/lang/Math.nextDown(F)F"] = GMeth{ParamSlots: 1, GFunction: nextDownFloat64}
	MethodSignatures["java/lang/Math.nextUp(D)D"] = GMeth{ParamSlots: 1, GFunction: nextUpFloat64}
	MethodSignatures["java/lang/Math.nextUp(F)F"] = GMeth{ParamSlots: 1, GFunction: nextUpFloat64}
	MethodSignatures["java/lang/Math.pow(DD)D"] = GMeth{ParamSlots: 2, GFunction: powFloat64}
	MethodSignatures["java/lang/Math.random()D"] = GMeth{ParamSlots: 0, GFunction: randomFloat64}
	MethodSignatures["java/lang/Math.rint(D)D"] = GMeth{ParamSlots: 1, GFunction: rintFloat64}
	MethodSignatures["java/lang/Math.round(D)J"] = GMeth{ParamSlots: 1, GFunction: roundInt64}
	MethodSignatures["java/lang/Math.round(F)I"] = GMeth{ParamSlots: 1, GFunction: roundInt64}
	MethodSignatures["java/lang/Math.scalb(DI)D"] = GMeth{ParamSlots: 2, GFunction: scalbDI}
	MethodSignatures["java/lang/Math.scalb(FI)F"] = GMeth{ParamSlots: 2, GFunction: scalbFI}
	MethodSignatures["java/lang/Math.signum(D)D"] = GMeth{ParamSlots: 1, GFunction: signumFloat64}
	MethodSignatures["java/lang/Math.signum(F)F"] = GMeth{ParamSlots: 1, GFunction: signumFloat64}
	MethodSignatures["java/lang/Math.sin(D)D"] = GMeth{ParamSlots: 1, GFunction: sinFloat64}
	MethodSignatures["java/lang/Math.sinh(D)D"] = GMeth{ParamSlots: 1, GFunction: sinhFloat64}
	MethodSignatures["java/lang/Math.sqrt(D)D"] = GMeth{ParamSlots: 1, GFunction: sqrtFloat64}
	MethodSignatures["java/lang/Math.subtractExact(II)I"] = GMeth{ParamSlots: 2, GFunction: subtractExactII}
	MethodSignatures["java/lang/Math.subtractExact(JJ)J"] = GMeth{ParamSlots: 2, GFunction: subtractExactJJ}
	MethodSignatures["java/lang/Math.tan(D)D"] = GMeth{ParamSlots: 1, GFunction: tanFloat64}
	MethodSignatures["java/lang/Math.tanh(D)D"] = GMeth{ParamSlots: 1, GFunction: tanhFloat64}
	MethodSignatures["java/lang/Math.toDegrees(D)D"] = GMeth{ParamSlots: 1, GFunction: toDegreesFloat64}
	MethodSignatures["java/lang/Math.toIntExact(J)I"] = GMeth{ParamSlots: 1, GFunction: toIntExactInt64}
	MethodSignatures["java/lang/Math.toRadians(D)D"] = GMeth{ParamSlots: 1, GFunction: toRadiansFloat64}
	MethodSignatures["java/lang/Math.ulp(D)D"] = GMeth{ParamSlots: 1, GFunction: ulpFloat64}
	MethodSignatures["java/lang/Math.ulp(F)F"] = GMeth{ParamSlots: 1, GFunction: ulpFloat64}

	MethodSignatures["java/lang/StrictMath.abs(D)D"] = GMeth{ParamSlots: 1, GFunction: absFloat64}
	MethodSignatures["java/lang/StrictMath.abs(F)F"] = GMeth{ParamSlots: 1, GFunction: absFloat64}
	MethodSignatures["java/lang/StrictMath.abs(I)I"] = GMeth{ParamSlots: 1, GFunction: absInt64}
	MethodSignatures["java/lang/StrictMath.abs(J)J"] = GMeth{ParamSlots: 1, GFunction: absInt64}
	MethodSignatures["java/lang/StrictMath.absExact(I)I"] = GMeth{ParamSlots: 1, GFunction: absInt64}
	MethodSignatures["java/lang/StrictMath.absExact(J)J"] = GMeth{ParamSlots: 1, GFunction: absInt64}
	MethodSignatures["java/lang/StrictMath.acos(D)D"] = GMeth{ParamSlots: 1, GFunction: acosFloat64}
	MethodSignatures["java/lang/StrictMath.addExact(II)I"] = GMeth{ParamSlots: 2, GFunction: addExactII}
	MethodSignatures["java/lang/StrictMath.addExact(JJ)J"] = GMeth{ParamSlots: 2, GFunction: addExactJJ}
	MethodSignatures["java/lang/StrictMath.asin(D)D"] = GMeth{ParamSlots: 1, GFunction: asinFloat64}
	MethodSignatures["java/lang/StrictMath.atan(D)D"] = GMeth{ParamSlots: 1, GFunction: atanFloat64}
	MethodSignatures["java/lang/StrictMath.atan2(DD)D"] = GMeth{ParamSlots: 2, GFunction: atan2Float64}
	MethodSignatures["java/lang/StrictMath.cbrt(D)D"] = GMeth{ParamSlots: 1, GFunction: cbrtFloat64}
	MethodSignatures["java/lang/StrictMath.ceil(D)D"] = GMeth{ParamSlots: 1, GFunction: ceilFloat64}
	MethodSignatures["java/lang/StrictMath.copySign(DD)D"] = GMeth{ParamSlots: 2, GFunction: copySignDD}
	MethodSignatures["java/lang/StrictMath.copySign(FF)F"] = GMeth{ParamSlots: 2, GFunction: copySignFF}
	MethodSignatures["java/lang/StrictMath.cos(D)D"] = GMeth{ParamSlots: 1, GFunction: cosFloat64}
	MethodSignatures["java/lang/StrictMath.cosh(D)D"] = GMeth{ParamSlots: 1, GFunction: coshFloat64}
	MethodSignatures["java/lang/StrictMath.decrementExact(I)I"] = GMeth{ParamSlots: 1, GFunction: decrementExactInt64}
	MethodSignatures["java/lang/StrictMath.decrementExact(J)J"] = GMeth{ParamSlots: 1, GFunction: decrementExactInt64}
	MethodSignatures["java/lang/StrictMath.exp(D)D"] = GMeth{ParamSlots: 1, GFunction: expFloat64}
	MethodSignatures["java/lang/StrictMath.expm1(D)D"] = GMeth{ParamSlots: 1, GFunction: expm1Float64}
	MethodSignatures["java/lang/StrictMath.floor(D)D"] = GMeth{ParamSlots: 1, GFunction: floorFloat64}
	MethodSignatures["java/lang/StrictMath.floorDiv(II)I"] = GMeth{ParamSlots: 2, GFunction: floorDivII}
	MethodSignatures["java/lang/StrictMath.floorDiv(JI)J"] = GMeth{ParamSlots: 2, GFunction: floorDivJx}
	MethodSignatures["java/lang/StrictMath.floorDiv(JJ)J"] = GMeth{ParamSlots: 2, GFunction: floorDivJx}
	MethodSignatures["java/lang/StrictMath.floorMod(II)I"] = GMeth{ParamSlots: 2, GFunction: floorModII}
	MethodSignatures["java/lang/StrictMath.floorMod(JI)I"] = GMeth{ParamSlots: 2, GFunction: floorModJx}
	MethodSignatures["java/lang/StrictMath.floorMod(JJ)J"] = GMeth{ParamSlots: 2, GFunction: floorModJx}
	MethodSignatures["java/lang/StrictMath.fma(DDD)D"] = GMeth{ParamSlots: 3, GFunction: fmaDDD}
	MethodSignatures["java/lang/StrictMath.fma(FFF)F"] = GMeth{ParamSlots: 3, GFunction: fmaFFF}
	MethodSignatures["java/lang/StrictMath.getExponent(D)I"] = GMeth{ParamSlots: 1, GFunction: getExponentFloat64}
	MethodSignatures["java/lang/StrictMath.getExponent(F)I"] = GMeth{ParamSlots: 1, GFunction: getExponentFloat64}
	MethodSignatures["java/lang/StrictMath.hypot(DD)D"] = GMeth{ParamSlots: 2, GFunction: hypotFloat64}
	MethodSignatures["java/lang/StrictMath.IEEEremainder(DD)D"] = GMeth{ParamSlots: 2, GFunction: IEEEremainderFloat64}
	MethodSignatures["java/lang/StrictMath.incrementExact(I)I"] = GMeth{ParamSlots: 1, GFunction: incrementExactInt64}
	MethodSignatures["java/lang/StrictMath.incrementExact(J)J"] = GMeth{ParamSlots: 1, GFunction: incrementExactInt64}
	MethodSignatures["java/lang/StrictMath.log(D)D"] = GMeth{ParamSlots: 1, GFunction: logFloat64}
	MethodSignatures["java/lang/StrictMath.log10(D)D"] = GMeth{ParamSlots: 1, GFunction: log10Float64}
	MethodSignatures["java/lang/StrictMath.log1p(D)D"] = GMeth{ParamSlots: 1, GFunction: log1pFloat64}
	MethodSignatures["java/lang/StrictMath.max(DD)D"] = GMeth{ParamSlots: 2, GFunction: maxDD}
	MethodSignatures["java/lang/StrictMath.max(FF)F"] = GMeth{ParamSlots: 2, GFunction: maxFF}
	MethodSignatures["java/lang/StrictMath.max(II)I"] = GMeth{ParamSlots: 2, GFunction: maxII}
	MethodSignatures["java/lang/StrictMath.max(JJ)J"] = GMeth{ParamSlots: 2, GFunction: maxJJ}
	MethodSignatures["java/lang/StrictMath.min(DD)D"] = GMeth{ParamSlots: 2, GFunction: minDD}
	MethodSignatures["java/lang/StrictMath.min(FF)F"] = GMeth{ParamSlots: 2, GFunction: minFF}
	MethodSignatures["java/lang/StrictMath.min(II)I"] = GMeth{ParamSlots: 2, GFunction: minII}
	MethodSignatures["java/lang/StrictMath.min(JJ)J"] = GMeth{ParamSlots: 2, GFunction: minJJ}
	MethodSignatures["java/lang/StrictMath.multiplyExact(II)I"] = GMeth{ParamSlots: 2, GFunction: multiplyExactII}
	MethodSignatures["java/lang/StrictMath.multiplyExact(JI)I"] = GMeth{ParamSlots: 2, GFunction: multiplyExactJx}
	MethodSignatures["java/lang/StrictMath.multiplyExact(JJ)J"] = GMeth{ParamSlots: 2, GFunction: multiplyExactJx}
	MethodSignatures["java/lang/StrictMath.multiplyHigh(JJ)J"] = GMeth{ParamSlots: 2, GFunction: multiplyHighJJ}
	MethodSignatures["java/lang/StrictMath.negateExact(I)I"] = GMeth{ParamSlots: 1, GFunction: negateExactInt64}
	MethodSignatures["java/lang/StrictMath.negateExact(J)J"] = GMeth{ParamSlots: 1, GFunction: negateExactInt64}
	MethodSignatures["java/lang/StrictMath.nextAfter(DD)D"] = GMeth{ParamSlots: 2, GFunction: nextAfterDD}
	MethodSignatures["java/lang/StrictMath.nextAfter(FD)F"] = GMeth{ParamSlots: 2, GFunction: nextAfterFD}
	MethodSignatures["java/lang/StrictMath.nextDown(D)D"] = GMeth{ParamSlots: 1, GFunction: nextDownFloat64}
	MethodSignatures["java/lang/StrictMath.nextDown(F)F"] = GMeth{ParamSlots: 1, GFunction: nextDownFloat64}
	MethodSignatures["java/lang/StrictMath.nextUp(D)D"] = GMeth{ParamSlots: 1, GFunction: nextUpFloat64}
	MethodSignatures["java/lang/StrictMath.nextUp(F)F"] = GMeth{ParamSlots: 1, GFunction: nextUpFloat64}
	MethodSignatures["java/lang/StrictMath.pow(DD)D"] = GMeth{ParamSlots: 2, GFunction: powFloat64}
	MethodSignatures["java/lang/StrictMath.random()D"] = GMeth{ParamSlots: 0, GFunction: randomFloat64}
	MethodSignatures["java/lang/StrictMath.rint(D)D"] = GMeth{ParamSlots: 1, GFunction: rintFloat64}
	MethodSignatures["java/lang/StrictMath.round(D)J"] = GMeth{ParamSlots: 1, GFunction: roundInt64}
	MethodSignatures["java/lang/StrictMath.round(F)I"] = GMeth{ParamSlots: 1, GFunction: roundInt64}
	MethodSignatures["java/lang/StrictMath.scalb(DI)D"] = GMeth{ParamSlots: 2, GFunction: scalbDI}
	MethodSignatures["java/lang/StrictMath.scalb(FI)F"] = GMeth{ParamSlots: 2, GFunction: scalbFI}
	MethodSignatures["java/lang/StrictMath.signum(D)D"] = GMeth{ParamSlots: 1, GFunction: signumFloat64}
	MethodSignatures["java/lang/StrictMath.signum(F)F"] = GMeth{ParamSlots: 1, GFunction: signumFloat64}
	MethodSignatures["java/lang/StrictMath.sin(D)D"] = GMeth{ParamSlots: 1, GFunction: sinFloat64}
	MethodSignatures["java/lang/StrictMath.sinh(D)D"] = GMeth{ParamSlots: 1, GFunction: sinhFloat64}
	MethodSignatures["java/lang/StrictMath.sqrt(D)D"] = GMeth{ParamSlots: 1, GFunction: sqrtFloat64}
	MethodSignatures["java/lang/StrictMath.subtractExact(II)I"] = GMeth{ParamSlots: 2, GFunction: subtractExactII}
	MethodSignatures["java/lang/StrictMath.subtractExact(JJ)J"] = GMeth{ParamSlots: 2, GFunction: subtractExactJJ}
	MethodSignatures["java/lang/StrictMath.tan(D)D"] = GMeth{ParamSlots: 1, GFunction: tanFloat64}
	MethodSignatures["java/lang/StrictMath.tanh(D)D"] = GMeth{ParamSlots: 1, GFunction: tanhFloat64}
	MethodSignatures["java/lang/StrictMath.toDegrees(D)D"] = GMeth{ParamSlots: 1, GFunction: toDegreesFloat64}
	MethodSignatures["java/lang/StrictMath.toIntExact(J)I"] = GMeth{ParamSlots: 1, GFunction: toIntExactInt64}
	MethodSignatures["java/lang/StrictMath.toRadians(D)D"] = GMeth{ParamSlots: 1, GFunction: toRadiansFloat64}
	MethodSignatures["java/lang/StrictMath.ulp(D)D"] = GMeth{ParamSlots: 1, GFunction: ulpFloat64}
	MethodSignatures["java/lang/StrictMath.ulp(F)F"] = GMeth{ParamSlots: 1, GFunction: ulpFloat64}

	MethodSignatures["java/lang/Math.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  mathClinit,
		}

}

func mathClinit([]interface{}) interface{} {
	klass := classloader.MethAreaFetch("java/lang/Math")
	if klass == nil {
		errMsg := "mathClinit: Expected java/lang/Math to be in the MethodArea, but it was not"
		return getGErrBlk(excNames.InternalException, errMsg)
	}
	return object.StringObjectFromGoString("mathClinit")
}

// Absolute value function for Java float and double
func absFloat64(params []interface{}) interface{} {
	return math.Abs(params[0].(float64))
}

// Absolute value function for Java int and long
func absInt64(params []interface{}) interface{} {
	xx := params[0].(int64)
	if xx < 0 {
		return -xx
	}
	return xx
}

// Arc cosine of a value; the returned angle is in the range 0.0 through pi.
func acosFloat64(params []interface{}) interface{} {
	return math.Acos(params[0].(float64))
}

// Sum of its arguments
func addExactII(params []interface{}) interface{} {
	return params[0].(int64) + params[1].(int64)
}
func addExactJJ(params []interface{}) interface{} {
	return params[0].(int64) + params[1].(int64)
}

// Arc sine of a value; the returned angle is in the range -pi/2 through pi/2.
func asinFloat64(params []interface{}) interface{} {
	return math.Asin(params[0].(float64))
}

// Arc tangent of a value; the returned angle is in the range -pi/2 through pi/2.
func atanFloat64(params []interface{}) interface{} {
	return math.Atan(params[0].(float64))
}

// Returns the angle theta from the conversion of rectangular coordinates (x, y)
// to polar coordinates (r, theta).
func atan2Float64(params []interface{}) interface{} {
	return math.Atan2(params[0].(float64), params[1].(float64))
}

// Cube root of a double value.
func cbrtFloat64(params []interface{}) interface{} {
	return math.Cbrt(params[0].(float64))
}

// Smallest (closest to negative infinity) double value that is
// greater than or equal to the argument and is equal to a mathematical integer.
func ceilFloat64(params []interface{}) interface{} {
	return math.Ceil(params[0].(float64))
}

// Amend the first argument with the sign of the second argument.
func copySignFF(params []interface{}) interface{} {
	return math.Copysign(params[0].(float64), params[1].(float64))
}
func copySignDD(params []interface{}) interface{} {
	return math.Copysign(params[0].(float64), params[1].(float64))
}

// Cosine of an angle expressed in radians.
func cosFloat64(params []interface{}) interface{} {
	return math.Cos(params[0].(float64))
}

// Hyperbolic cosine of an angle expressed in radians.
func coshFloat64(params []interface{}) interface{} {
	return math.Cosh(params[0].(float64))
}

// Decrement the argument by 1
func decrementExactInt64(params []interface{}) interface{} {
	return params[0].(int64) - 1
}

// Euler's number e raised to the power of a double value.
func expFloat64(params []interface{}) interface{} {
	return math.Exp(params[0].(float64))
}

// Euler's number e raised to the power of a double value minus 1.
func expm1Float64(params []interface{}) interface{} {
	return math.Expm1(params[0].(float64))
}

// Largest (closest to positive infinity) double value that is less than or equal to
// the argument and is equal to a mathematical integer.
func floorFloat64(params []interface{}) interface{} {
	return math.Floor(params[0].(float64))
}

// Largest (closest to positive infinity) int value that is less than or equal
// to the algebraic quotient.
func floorDivInt64(dividend int64, divisor int64) interface{} {
	if divisor == 0 {
		return getGErrBlk(excNames.ArithmeticException, "floorDivInt64: Divide by zero")
	}
	if dividend <= math.MinInt64 && divisor == -1 {
		return math.MinInt64
	}
	if (dividend <= 0 && divisor < 00) || (dividend >= 0 && divisor > 00) {
		return dividend / divisor
	}
	// At this point, (a) x and y are nonzero and (b) they have opposite signs.
	return (dividend / divisor) - 1
}
func floorDivII(params []interface{}) interface{} {
	dividend := params[0].(int64)
	divisor := params[1].(int64)
	return floorDivInt64(dividend, divisor)
}
func floorDivJx(params []interface{}) interface{} {
	dividend := params[0].(int64)
	divisor := params[1].(int64)
	return floorDivInt64(dividend, divisor)
}

// Largest (closest to positive infinity) int value that is less than or equal
// to the algebraic quotient.
// params[0]=dividend=x
// params[1]=divisor=y
// floorDiv(x, y) * y + floorMod(x, y) = x
// Therefore, floorMod(x, y) = x - floorDiv(x, y) * y
func floorModII(params []interface{}) interface{} {
	fldiv := floorDivII(params)
	switch fldiv.(type) {
	case *GErrBlk:
		// Return the G function error block
		return fldiv
	}
	result := params[0].(int64) - fldiv.(int64)*params[1].(int64)
	return result
}
func floorModJx(params []interface{}) interface{} {
	fldiv := floorDivJx(params)
	switch fldiv.(type) {
	case *GErrBlk:
		// Return the G function error block
		return fldiv
	}
	result := params[0].(int64) - fldiv.(int64)*params[1].(int64)
	return result
}

// FMA (fused multiply add) the three arguments; that is, returns the exact product
// of the first two arguments summed with the third argument and then rounded once to the nearest double.
func fmaDDD(params []interface{}) interface{} {
	xx := params[0].(float64)
	yy := params[1].(float64)
	zz := params[2].(float64)
	return math.FMA(xx, yy, zz)
}
func fmaFFF(params []interface{}) interface{} {
	xx := params[0].(float64)
	yy := params[1].(float64)
	zz := params[2].(float64)
	return math.FMA(xx, yy, zz)
}

// Unbiased exponent used in the representation of a double or float.
func getExponentFloat64(params []interface{}) interface{} {
	xx := params[0].(float64)

	// Check if the number is NaN or infinite
	if math.IsNaN(xx) || math.IsInf(xx, 0) {
		return MAX_DOUBLE_EXPONENT + 1
	}

	// Extract the exponent bits using math.Float64bits
	bits := math.Float64bits(xx)
	exponentBits := int64((bits >> 52) & 0x7FF)

	// Subtract the bias to get the actual exponent
	return exponentBits - MAX_DOUBLE_EXPONENT
}

// Sqrt(x^2 + y^2) without intermediate overflow or underflow.
func hypotFloat64(params []interface{}) interface{} {
	return math.Hypot(params[0].(float64), params[1].(float64))
}

// Remainder operation on two arguments as prescribed by the IEEE 754 standard.
func IEEEremainderFloat64(params []interface{}) interface{} {
	return math.Remainder(params[0].(float64), params[1].(float64))
}

// Increment the argument by 1
func incrementExactInt64(params []interface{}) interface{} {
	return params[0].(int64) + 1
}

// Natural logarithm (base e) of a double value.
func logFloat64(params []interface{}) interface{} {
	return math.Log(params[0].(float64))
}

// Base 10 logarithm of a double value.
func log10Float64(params []interface{}) interface{} {
	return math.Log10(params[0].(float64))
}

// Natural logarithm (base e) of (double value + 1).
func log1pFloat64(params []interface{}) interface{} {
	return math.Log1p(params[0].(float64))
}

// Maximum functions.
func maxDD(params []interface{}) interface{} {
	return math.Max(params[0].(float64), params[1].(float64))
}
func maxFF(params []interface{}) interface{} {
	return math.Max(params[0].(float64), params[1].(float64))
}
func maxII(params []interface{}) interface{} {
	xx := params[0].(int64)
	yy := params[1].(int64)
	if xx > yy {
		return xx
	}
	return yy
}
func maxJJ(params []interface{}) interface{} {
	xx := params[0].(int64)
	yy := params[1].(int64)
	if xx > yy {
		return xx
	}
	return yy
}

// Minimum functions.
func minDD(params []interface{}) interface{} {
	return math.Min(params[0].(float64), params[1].(float64))
}
func minFF(params []interface{}) interface{} {
	return math.Min(params[0].(float64), params[1].(float64))
}
func minII(params []interface{}) interface{} {
	xx := params[0].(int64)
	yy := params[1].(int64)
	if xx < yy {
		return xx
	}
	return yy
}
func minJJ(params []interface{}) interface{} {
	xx := params[0].(int64)
	yy := params[1].(int64)
	if xx < yy {
		return xx
	}
	return yy
}

// Product of the arguments.
func multiplyExactII(params []interface{}) interface{} {
	return params[0].(int64) * params[1].(int64)
}
func multiplyExactJx(params []interface{}) interface{} {
	return params[0].(int64) * params[1].(int64)
}

// Most significant 64 bits of the 128-bit product of two 64-bit factors.
func multiplyHighJJ(params []interface{}) interface{} {
	xx := big.NewInt(params[0].(int64))
	yy := big.NewInt(params[1].(int64))
	zz := big.NewInt(0)
	zz.Mul(xx, yy)
	zz.Rsh(zz, 64)
	return zz.Int64()
}

// Negation of the argument for int and long.
func negateExactInt64(params []interface{}) interface{} {
	return -params[0].(int64)
}

// Next after double of float value.
func nextAfterDD(params []interface{}) interface{} {
	return math.Nextafter(params[0].(float64), params[1].(float64))
}
func nextAfterFD(params []interface{}) interface{} {
	return math.Nextafter(params[0].(float64), params[1].(float64))
}

// Next down double of float value.
func nextDownFloat64(params []interface{}) interface{} {
	return math.Nextafter(params[0].(float64), math.Inf(-1))
}

// Next up double of float value.
func nextUpFloat64(params []interface{}) interface{} {
	return math.Nextafter(params[0].(float64), math.Inf(+1))
}

// Value of the first argument raised to the power of the second argument.
func powFloat64(params []interface{}) interface{} {
	return math.Pow(params[0].(float64), params[1].(float64))
}

// Generate a random number >= 0.0 and < 1.0
func randomFloat64(params []interface{}) interface{} {
	return rand.Float64()
}

// Computes a double-valued number that is closest in value to the argument and is equal to a mathematical integer.
func rintFloat64(params []interface{}) interface{} {
	return math.Round(params[0].(float64))
}

// Computes the closest long to the argument, with ties rounding towards positive infinity.
func roundInt64(params []interface{}) interface{} {
	return int64(math.Round(params[0].(float64)))
}

// Compute the product of the argument and 2 raised to the power of the scaleFactor.
func scalbFloat64I(xx float64, scaleFactor int64) float64 {
	if math.IsNaN(xx) {
		return math.NaN()
	}
	if math.IsInf(xx, -1) {
		return math.Inf(-1)
	}
	if math.IsInf(xx, 1) {
		return math.Inf(1)
	}
	result := xx * math.Pow(2.0, float64(scaleFactor))
	return result
}
func scalbDI(params []interface{}) interface{} {
	xx := params[0].(float64)
	scaleFactor := params[1].(int64)
	return scalbFloat64I(xx, scaleFactor)
}
func scalbFI(params []interface{}) interface{} {
	xx := params[0].(float64)
	scaleFactor := params[1].(int64)
	return scalbFloat64I(xx, scaleFactor)
}

// Compute the signum value of an argument.
func signumFloat64(params []interface{}) interface{} {
	xx := params[0].(float64)
	if math.IsNaN(xx) {
		return math.NaN()
	}
	if xx > 0 {
		return 1.0
	} else if xx < 0 {
		return -1.0
	}
	return 0.0
}

// Compute the sine of an angle expressed in radians.
func sinFloat64(params []interface{}) interface{} {
	return math.Sin(params[0].(float64))
}

// Compute the hyperbolic sine of an angle expressed in radians.
func sinhFloat64(params []interface{}) interface{} {
	return math.Sinh(params[0].(float64))
}

// Compute a square root.
func sqrtFloat64(params []interface{}) interface{} {
	return math.Sqrt(params[0].(float64))
}

// Difference of its arguments
func subtractExactII(params []interface{}) interface{} {
	return params[0].(int64) - params[1].(int64)
}
func subtractExactJJ(params []interface{}) interface{} {
	return params[0].(int64) - params[1].(int64)
}

// Compute the tangent of an angle expressed in radians.
func tanFloat64(params []interface{}) interface{} {
	return math.Tan(params[0].(float64))
}

// Compute the hyperbolic tangent of an angle expressed in radians.
func tanhFloat64(params []interface{}) interface{} {
	return math.Tanh(params[0].(float64))
}

// Convert radians to degrees.
func toDegreesFloat64(params []interface{}) interface{} {
	return params[0].(float64) * 180.0 / PI
}

// Not very interesting as long and its are both int64.
func toIntExactInt64(params []interface{}) interface{} {
	return params[0].(int64)
}

// Convert degrees to radians.
func toRadiansFloat64(params []interface{}) interface{} {
	return params[0].(float64) * PI / 180.0
}

// ULP: Unit of Least Precision.
func ulpFloat64(params []interface{}) interface{} {
	xx := params[0].(float64)
	if math.IsNaN(xx) {
		return xx
	}
	if math.IsInf(xx, -1) {
		return math.Inf(1) // "If the argument is positive or negative infinity, then the result is positive infinity."
	}
	if math.IsInf(xx, 1) {
		return math.Inf(1)
	}
	xx = math.Abs(xx)
	if math.IsInf(xx, +1) {
		return xx
	}
	next := math.Nextafter(xx, math.Inf(1))
	if math.IsInf(next, 1) {
		next = math.Nextafter(xx, math.Inf(-1))
		return xx - next
	}
	return next - xx
}
