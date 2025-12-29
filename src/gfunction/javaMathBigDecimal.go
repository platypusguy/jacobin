/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

var classNameBigDecimal = "java/math/BigDecimal"

func Load_Math_Big_Decimal() {

	MethodSignatures[classNameBigDecimal+".<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures[classNameBigDecimal+".<init>([C)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures[classNameBigDecimal+".<init>([CII)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures[classNameBigDecimal+".<init>([CIILjava/math/MathContext;)V"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  trapFunction,
		}

	MethodSignatures[classNameBigDecimal+".<init>([CLjava/math/MathContext;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures[classNameBigDecimal+".<init>(D)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalInitDouble,
		}

	MethodSignatures[classNameBigDecimal+".<init>(DLjava/math/MathContext;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalInitDoubleContext,
		}

	MethodSignatures[classNameBigDecimal+".<init>(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalInitIntLong,
		}

	MethodSignatures[classNameBigDecimal+".<init>(ILjava/math/MathContext;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures[classNameBigDecimal+".<init>(J)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalInitIntLong,
		}

	MethodSignatures[classNameBigDecimal+".<init>(JLjava/math/MathContext;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures[classNameBigDecimal+".<init>(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalInitString,
		}

	MethodSignatures[classNameBigDecimal+".<init>(Ljava/lang/String;Ljava/math/MathContext;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalInitStringContext,
		}

	MethodSignatures[classNameBigDecimal+".<init>(Ljava/math/BigInteger;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalInitBigInteger,
		}

	MethodSignatures[classNameBigDecimal+".<init>(Ljava/math/BigInteger;I)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalInitBigIntegerScale,
		}

	MethodSignatures[classNameBigDecimal+".<init>(Ljava/math/BigInteger;Ljava/math/MathContext;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalInitBigIntegerContext,
		}

	// ---------------- end of <clinit> and <init> -------------------------------------------

	MethodSignatures[classNameBigDecimal+".abs()Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalAbs,
		}

	MethodSignatures[classNameBigDecimal+".add(Ljava/math/BigDecimal;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalAdd,
		}

	MethodSignatures[classNameBigDecimal+".add(Ljava/math/BigDecimal;Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalAddContext,
		}

	MethodSignatures[classNameBigDecimal+".byteValueExact()B"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalByteValueExact,
		}

	MethodSignatures[classNameBigDecimal+".compareTo(Ljava/math/BigDecimal;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalCompareTo,
		}

	MethodSignatures[classNameBigDecimal+".divide(Ljava/math/BigDecimal;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalDivide,
		}

	MethodSignatures[classNameBigDecimal+".divide(Ljava/math/BigDecimal;I)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapDeprecated,
		}

	MethodSignatures[classNameBigDecimal+".divide(Ljava/math/BigDecimal;II)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapDeprecated,
		}

	MethodSignatures[classNameBigDecimal+".divide(Ljava/math/BigDecimal;ILjava/math/RoundingMode;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  bigdecimalDivideScaleRoundingMode,
		}

	MethodSignatures[classNameBigDecimal+".divide(Ljava/math/BigDecimal;Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalDivideMathContext,
		}

	MethodSignatures[classNameBigDecimal+".divide(Ljava/math/BigDecimal;Ljava/math/RoundingMode;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalDivideRoundingMode,
		}

	MethodSignatures[classNameBigDecimal+".divideAndRemainder(Ljava/math/BigDecimal;)[Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalDivideAndRemainder,
		}

	MethodSignatures[classNameBigDecimal+".divideAndRemainder(Ljava/math/BigDecimal;Ljava/math/MathContext;)[Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalDivideAndRemainderContext,
		}

	MethodSignatures[classNameBigDecimal+".divideToIntegralValue(Ljava/math/BigDecimal;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalDivideToIntegralValue,
		}

	MethodSignatures[classNameBigDecimal+".divideToIntegralValue(Ljava/math/BigDecimal;Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalDivideToIntegralValueContext,
		}

	MethodSignatures[classNameBigDecimal+".doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalDoubleValue,
		}

	MethodSignatures[classNameBigDecimal+".equals(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalEquals,
		}

	MethodSignatures[classNameBigDecimal+".floatValue()F"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalFloatValue,
		}

	MethodSignatures[classNameBigDecimal+".hashCode()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures[classNameBigDecimal+".intValue()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalIntValue,
		}

	MethodSignatures[classNameBigDecimal+".intValueExact()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalIntValueExact,
		}

	MethodSignatures[classNameBigDecimal+".longValue()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalLongValue,
		}

	MethodSignatures[classNameBigDecimal+".longValueExact()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalLongValueExact,
		}

	MethodSignatures[classNameBigDecimal+".max(Ljava/math/BigDecimal;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalMax,
		}

	MethodSignatures[classNameBigDecimal+".min(Ljava/math/BigDecimal;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalMin,
		}

	MethodSignatures[classNameBigDecimal+".movePointLeft(I)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalMovePointLeft,
		}

	MethodSignatures[classNameBigDecimal+".movePointRight(I)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalMovePointRight,
		}

	MethodSignatures[classNameBigDecimal+".multiply(Ljava/math/BigDecimal;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalMultiply,
		}

	MethodSignatures[classNameBigDecimal+".multiply(Ljava/math/BigDecimal;Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalMultiplyContext,
		}

	MethodSignatures[classNameBigDecimal+".negate()Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalNegate,
		}

	MethodSignatures[classNameBigDecimal+".negate(Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalNegateContext,
		}

	MethodSignatures[classNameBigDecimal+".plus()Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalPlus,
		}

	MethodSignatures[classNameBigDecimal+".plus(Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalPlusContext,
		}

	MethodSignatures[classNameBigDecimal+".pow(I)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalPow,
		}

	MethodSignatures[classNameBigDecimal+".pow(ILjava/math/MathContext;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalPowContext,
		}

	MethodSignatures[classNameBigDecimal+".precision()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalPrecision,
		}

	MethodSignatures[classNameBigDecimal+".remainder(Ljava/math/BigDecimal;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalRemainder,
		}

	MethodSignatures[classNameBigDecimal+".remainder(Ljava/math/BigDecimal;Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalRemainderContext,
		}

	MethodSignatures[classNameBigDecimal+".round(Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalRoundContext,
		}

	MethodSignatures[classNameBigDecimal+".scale()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalScale,
		}

	MethodSignatures[classNameBigDecimal+".scaleByPowerOfTen(I)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalScaleByPowerOfTen,
		}

	MethodSignatures[classNameBigDecimal+".setScale(I)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalSetScale,
		}

	MethodSignatures[classNameBigDecimal+".setScale(II)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapDeprecated,
		}

	MethodSignatures[classNameBigDecimal+".setScale(ILjava/math/RoundingMode;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalSetScaleRoundingMode,
		}

	MethodSignatures[classNameBigDecimal+".shortValueExact()S"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalShortValueExact,
		}

	MethodSignatures[classNameBigDecimal+".signum()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalSignum,
		}

	MethodSignatures[classNameBigDecimal+".sqrt(Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalSqrtContext,
		}

	MethodSignatures[classNameBigDecimal+".stripTrailingZeros()Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalStripTrailingZeros,
		}

	MethodSignatures[classNameBigDecimal+".subtract(Ljava/math/BigDecimal;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalSubtract,
		}

	MethodSignatures[classNameBigDecimal+".subtract(Ljava/math/BigDecimal;Ljava/math/MathContext;)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalSubtractContext,
		}

	MethodSignatures[classNameBigDecimal+".toBigInteger()Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalToBigInteger,
		}

	MethodSignatures[classNameBigDecimal+".toBigIntegerExact()Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalToBigIntegerExact,
		}

	MethodSignatures[classNameBigDecimal+".toEngineeringString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalToString,
		}

	MethodSignatures[classNameBigDecimal+".toPlainString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalToString,
		}

	MethodSignatures[classNameBigDecimal+".toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalToString,
		}

	MethodSignatures[classNameBigDecimal+".ulp()Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalUlp,
		}

	MethodSignatures[classNameBigDecimal+".unscaledValue()Ljava/math/BigInteger;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  bigdecimalUnscaledValue,
		}

	MethodSignatures[classNameBigDecimal+".valueOf(D)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalValueOfDouble,
		}

	MethodSignatures[classNameBigDecimal+".valueOf(J)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  bigdecimalValueOfLong,
		}

	MethodSignatures[classNameBigDecimal+".valueOf(JI)Ljava/math/BigDecimal;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  bigdecimalValueOfLongInt,
		}

	// Avoiding a cycle issue by doing this here and not in package statics.
	loadStaticsBigDecimal()
}
