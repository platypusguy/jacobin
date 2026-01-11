/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaMath

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/types"
)

func Load_Math_Big_Decimal() {

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".<init>([C)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".<init>([CII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".<init>([CIILjava/math/MathContext;)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".<init>([CLjava/math/MathContext;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".<init>(D)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalInitDouble,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".<init>(DLjava/math/MathContext;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalInitDoubleContext,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".<init>(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalInitIntLong,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".<init>(ILjava/math/MathContext;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".<init>(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalInitIntLong,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".<init>(JLjava/math/MathContext;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  BigdecimalInitString,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".<init>(Ljava/lang/String;Ljava/math/MathContext;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalInitStringContext,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".<init>(Ljava/math/BigInteger;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalInitBigInteger,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".<init>(Ljava/math/BigInteger;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalInitBigIntegerScale,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".<init>(Ljava/math/BigInteger;Ljava/math/MathContext;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalInitBigIntegerContext,
		}

	// ---------------- end of <clinit> and <init> -------------------------------------------

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".abs()Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalAbs,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".add(Ljava/math/BigDecimal;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalAdd,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".add(Ljava/math/BigDecimal;Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalAddContext,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".byteValueExact()B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalByteValueExact,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".compareTo(Ljava/math/BigDecimal;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalCompareTo,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".divide(Ljava/math/BigDecimal;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalDivide,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".divide(Ljava/math/BigDecimal;I)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".divide(Ljava/math/BigDecimal;II)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".divide(Ljava/math/BigDecimal;ILjava/math/RoundingMode;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  bigdecimalDivideScaleRoundingMode,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".divide(Ljava/math/BigDecimal;Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalDivideMathContext,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".divide(Ljava/math/BigDecimal;Ljava/math/RoundingMode;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalDivideRoundingMode,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".divideAndRemainder(Ljava/math/BigDecimal;)[Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalDivideAndRemainder,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".divideAndRemainder(Ljava/math/BigDecimal;Ljava/math/MathContext;)[Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalDivideAndRemainderContext,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".divideToIntegralValue(Ljava/math/BigDecimal;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalDivideToIntegralValue,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".divideToIntegralValue(Ljava/math/BigDecimal;Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalDivideToIntegralValueContext,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".doubleValue()D"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalDoubleValue,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalEquals,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".floatValue()F"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalFloatValue,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".intValue()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalIntValue,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".intValueExact()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalIntValueExact,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".longValue()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalLongValue,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".longValueExact()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalLongValueExact,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".max(Ljava/math/BigDecimal;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalMax,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".min(Ljava/math/BigDecimal;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalMin,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".movePointLeft(I)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalMovePointLeft,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".movePointRight(I)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalMovePointRight,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".multiply(Ljava/math/BigDecimal;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalMultiply,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".multiply(Ljava/math/BigDecimal;Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalMultiplyContext,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".negate()Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalNegate,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".negate(Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalNegateContext,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".plus()Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalPlus,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".plus(Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalPlusContext,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".pow(I)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalPow,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".pow(ILjava/math/MathContext;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalPowContext,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".precision()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalPrecision,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".remainder(Ljava/math/BigDecimal;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalRemainder,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".remainder(Ljava/math/BigDecimal;Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalRemainderContext,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".round(Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalRoundContext,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".scale()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalScale,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".scaleByPowerOfTen(I)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalScaleByPowerOfTen,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".setScale(I)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalSetScale,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".setScale(II)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".setScale(ILjava/math/RoundingMode;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalSetScaleRoundingMode,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".shortValueExact()S"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalShortValueExact,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".signum()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalSignum,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".sqrt(Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalSqrtContext,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".stripTrailingZeros()Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalStripTrailingZeros,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".subtract(Ljava/math/BigDecimal;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalSubtract,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".subtract(Ljava/math/BigDecimal;Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalSubtractContext,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".toBigInteger()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalToBigInteger,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".toBigIntegerExact()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalToBigIntegerExact,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".toEngineeringString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalToString,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".toPlainString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalToString,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalToString,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".ulp()Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalUlp,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".unscaledValue()Ljava/math/BigInteger;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalUnscaledValue,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".valueOf(D)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalValueOfDouble,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".valueOf(J)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalValueOfLong,
		}

	ghelpers.MethodSignatures[types.ClassNameBigDecimal+".valueOf(JI)Ljava/math/BigDecimal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalValueOfLongInt,
		}

	// Avoiding a cycle issue by doing this here and not in package statics.

	loadStaticsBigDecimal()
}
