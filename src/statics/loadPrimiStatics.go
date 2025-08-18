package statics

// Reference: https://docs.oracle.com/en/java/javase/17/docs/api/constant-values.html

import (
	"jacobin/src/types"
	"math"
)

func LoadStaticsByte() {
	_ = AddStatic("java/lang/Byte.BYTES", Static{Type: types.Int, Value: int64(1)})
	_ = AddStatic("java/lang/Byte.MAX_VALUE", Static{Type: types.Byte, Value: int64(0x7f)})
	_ = AddStatic("java/lang/Byte.MIN_VALUE", Static{Type: types.Byte, Value: int64(0x80)})
	_ = AddStatic("java/lang/Byte.SIZE", Static{Type: types.Int, Value: int64(8)})
}

func LoadStaticsCharacter() {
	_ = AddStatic("java/lang/Character.BYTES", Static{Type: types.Int, Value: int64(2)})
	_ = AddStatic("java/lang/Character.COMBINING_SPACING_MARK", Static{Type: types.Byte, Value: int64(0x8)})
	_ = AddStatic("java/lang/Character.CONNECTOR_PUNCTUATION", Static{Type: types.Byte, Value: int64(0x17)})
	_ = AddStatic("java/lang/Character.CONTROL", Static{Type: types.Byte, Value: int64(0xf)})
	_ = AddStatic("java/lang/Character.CURRENCY_SYMBOL", Static{Type: types.Byte, Value: int64(0x1a)})
	_ = AddStatic("java/lang/Character.DASH_PUNCTUATION", Static{Type: types.Byte, Value: int64(0x14)})
	_ = AddStatic("java/lang/Character.DECIMAL_DIGIT_NUMBER", Static{Type: types.Byte, Value: int64(0x9)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_ARABIC_NUMBER", Static{Type: types.Byte, Value: int64(0x6)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_BOUNDARY_NEUTRAL", Static{Type: types.Byte, Value: int64(0x9)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_COMMON_NUMBER_SEPARATOR", Static{Type: types.Byte, Value: int64(0x7)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_EUROPEAN_NUMBER", Static{Type: types.Byte, Value: int64(0x3)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_EUROPEAN_NUMBER_SEPARATOR", Static{Type: types.Byte, Value: int64(0x4)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_EUROPEAN_NUMBER_TERMINATOR", Static{Type: types.Byte, Value: int64(0x5)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_FIRST_STRONG_ISOLATE", Static{Type: types.Byte, Value: int64(0x15)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_LEFT_TO_RIGHT", Static{Type: types.Byte, Value: int64(0x0)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_LEFT_TO_RIGHT_EMBEDDING", Static{Type: types.Byte, Value: int64(0xe)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_LEFT_TO_RIGHT_ISOLATE", Static{Type: types.Byte, Value: int64(0x13)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_LEFT_TO_RIGHT_OVERRIDE", Static{Type: types.Byte, Value: int64(0xf)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_NONSPACING_MARK", Static{Type: types.Byte, Value: int64(0x8)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_OTHER_NEUTRALS", Static{Type: types.Byte, Value: int64(0xd)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_PARAGRAPH_SEPARATOR", Static{Type: types.Byte, Value: int64(0xa)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_POP_DIRECTIONAL_FORMAT", Static{Type: types.Byte, Value: int64(0x12)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_POP_DIRECTIONAL_ISOLATE", Static{Type: types.Byte, Value: int64(0x16)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_RIGHT_TO_LEFT", Static{Type: types.Byte, Value: int64(0x1)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_RIGHT_TO_LEFT_ARABIC", Static{Type: types.Byte, Value: int64(0x2)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_RIGHT_TO_LEFT_EMBEDDING", Static{Type: types.Byte, Value: int64(0x10)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_RIGHT_TO_LEFT_ISOLATE", Static{Type: types.Byte, Value: int64(0x14)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_RIGHT_TO_LEFT_OVERRIDE", Static{Type: types.Byte, Value: int64(0x11)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_SEGMENT_SEPARATOR", Static{Type: types.Byte, Value: int64(0xb)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_UNDEFINED", Static{Type: types.Byte, Value: int64(0xff)})
	_ = AddStatic("java/lang/Character.DIRECTIONALITY_WHITESPACE", Static{Type: types.Byte, Value: int64(0xc)})
	_ = AddStatic("java/lang/Character.ENCLOSING_MARK", Static{Type: types.Byte, Value: int64(0x7)})
	_ = AddStatic("java/lang/Character.END_PUNCTUATION", Static{Type: types.Byte, Value: int64(0x16)})
	_ = AddStatic("java/lang/Character.FINAL_QUOTE_PUNCTUATION", Static{Type: types.Byte, Value: int64(0x1e)})
	_ = AddStatic("java/lang/Character.FORMAT", Static{Type: types.Byte, Value: int64(0x10)})
	_ = AddStatic("java/lang/Character.INITIAL_QUOTE_PUNCTUATION", Static{Type: types.Byte, Value: int64(0x1d)})
	_ = AddStatic("java/lang/Character.LETTER_NUMBER", Static{Type: types.Byte, Value: int64(0xa)})
	_ = AddStatic("java/lang/Character.LINE_SEPARATOR", Static{Type: types.Byte, Value: int64(0xd)})
	_ = AddStatic("java/lang/Character.LOWERCASE_LETTER", Static{Type: types.Byte, Value: int64(0x2)})
	_ = AddStatic("java/lang/Character.MATH_SYMBOL", Static{Type: types.Byte, Value: int64(0x19)})
	_ = AddStatic("java/lang/Character.MAX_CODE_POINT", Static{Type: types.Int, Value: int64(1114111)})
	_ = AddStatic("java/lang/Character.MAX_HIGH_SURROGATE", Static{Type: types.Char, Value: rune(56319)})
	_ = AddStatic("java/lang/Character.MAX_LOW_SURROGATE", Static{Type: types.Char, Value: rune(57343)})
	_ = AddStatic("java/lang/Character.MAX_RADIX", Static{Type: types.Int, Value: int64(36)})
	_ = AddStatic("java/lang/Character.MAX_SURROGATE", Static{Type: types.Char, Value: rune(57343)})
	_ = AddStatic("java/lang/Character.MAX_VALUE", Static{Type: types.Char, Value: rune(65535)})
	_ = AddStatic("java/lang/Character.MIN_CODE_POINT", Static{Type: types.Int, Value: int64(0)})
	_ = AddStatic("java/lang/Character.MIN_HIGH_SURROGATE", Static{Type: types.Char, Value: rune(55296)})
	_ = AddStatic("java/lang/Character.MIN_LOW_SURROGATE", Static{Type: types.Char, Value: rune(56320)})
	_ = AddStatic("java/lang/Character.MIN_RADIX", Static{Type: types.Int, Value: int64(2)})
	_ = AddStatic("java/lang/Character.MIN_SUPPLEMENTARY_CODE_POINT", Static{Type: types.Int, Value: int64(65536)})
	_ = AddStatic("java/lang/Character.MIN_SURROGATE", Static{Type: types.Char, Value: rune(55296)})
	_ = AddStatic("java/lang/Character.MIN_VALUE", Static{Type: types.Char, Value: rune(0)})
	_ = AddStatic("java/lang/Character.MODIFIER_LETTER", Static{Type: types.Byte, Value: int64(0x4)})
	_ = AddStatic("java/lang/Character.MODIFIER_SYMBOL", Static{Type: types.Byte, Value: int64(0x1b)})
	_ = AddStatic("java/lang/Character.NON_SPACING_MARK", Static{Type: types.Byte, Value: int64(0x6)})
	_ = AddStatic("java/lang/Character.OTHER_LETTER", Static{Type: types.Byte, Value: int64(0x5)})
	_ = AddStatic("java/lang/Character.OTHER_NUMBER", Static{Type: types.Byte, Value: int64(0xb)})
	_ = AddStatic("java/lang/Character.OTHER_PUNCTUATION", Static{Type: types.Byte, Value: int64(0x18)})
	_ = AddStatic("java/lang/Character.OTHER_SYMBOL", Static{Type: types.Byte, Value: int64(0x1c)})
	_ = AddStatic("java/lang/Character.PARAGRAPH_SEPARATOR", Static{Type: types.Byte, Value: int64(0xe)})
	_ = AddStatic("java/lang/Character.PRIVATE_USE", Static{Type: types.Byte, Value: int64(0x12)})
	_ = AddStatic("java/lang/Character.SIZE", Static{Type: types.Int, Value: int64(16)})
	_ = AddStatic("java/lang/Character.SPACE_SEPARATOR", Static{Type: types.Byte, Value: int64(0xc)})
	_ = AddStatic("java/lang/Character.START_PUNCTUATION", Static{Type: types.Byte, Value: int64(0x15)})
	_ = AddStatic("java/lang/Character.SURROGATE", Static{Type: types.Byte, Value: int64(0x13)})
	_ = AddStatic("java/lang/Character.TITLECASE_LETTER", Static{Type: types.Byte, Value: int64(0x3)})
	_ = AddStatic("java/lang/Character.UNASSIGNED", Static{Type: types.Byte, Value: int64(0x0)})
	_ = AddStatic("java/lang/Character.UPPERCASE_LETTER", Static{Type: types.Byte, Value: int64(0x1)})
}

func LoadStaticsDouble() {
	_ = AddStatic("java/lang/Double.BYTES", Static{Type: types.Int, Value: int64(8)})
	_ = AddStatic("java/lang/Double.MAX_EXPONENT", Static{Type: types.Int, Value: int64(1023)})
	_ = AddStatic("java/lang/Double.MAX_VALUE", Static{Type: types.Double, Value: float64(1.7976931348623157e308)})
	_ = AddStatic("java/lang/Double.MIN_EXPONENT", Static{Type: types.Int, Value: int64(-1022)})
	_ = AddStatic("java/lang/Double.MIN_NORMAL", Static{Type: types.Double, Value: float64(2.2250738585072014e-308)})
	_ = AddStatic("java/lang/Double.MIN_VALUE", Static{Type: types.Double, Value: float64(4.9e-324)})
	_ = AddStatic("java/lang/Double.NaN", Static{Type: types.Double, Value: float64(math.NaN())})
	_ = AddStatic("java/lang/Double.NEGATIVE_INFINITY", Static{Type: types.Double, Value: float64(math.Inf(-1))})
	_ = AddStatic("java/lang/Double.POSITIVE_INFINITY", Static{Type: types.Double, Value: float64(math.Inf(+1))})
	_ = AddStatic("java/lang/Double.SIZE", Static{Type: types.Int, Value: int64(64)})
}

func LoadStaticsFloat() {
	_ = AddStatic("java/lang/Float.BYTES", Static{Type: types.Int, Value: int64(4)})
	_ = AddStatic("java/lang/Float.MAX_EXPONENT", Static{Type: types.Int, Value: int64(127)})
	_ = AddStatic("java/lang/Float.MAX_VALUE", Static{Type: types.Float, Value: float64(3.4028234663852886e38)})
	_ = AddStatic("java/lang/Float.MIN_EXPONENT", Static{Type: types.Int, Value: int64(-126)})
	_ = AddStatic("java/lang/Float.MIN_NORMAL", Static{Type: types.Float, Value: float64(1.1754943508222875e-38)})
	_ = AddStatic("java/lang/Float.MIN_VALUE", Static{Type: types.Float, Value: float64(1.401298464324817e-45)})
	_ = AddStatic("java/lang/Float.NaN", Static{Type: types.Float, Value: float64(math.NaN())})
	_ = AddStatic("java/lang/Float.NEGATIVE_INFINITY", Static{Type: types.Float, Value: float64(math.Inf(-1))})
	_ = AddStatic("java/lang/Float.POSITIVE_INFINITY", Static{Type: types.Float, Value: float64(math.Inf(+1))})
	_ = AddStatic("java/lang/Float.SIZE", Static{Type: types.Int, Value: int64(32)})
}

func LoadStaticsInteger() {
	_ = AddStatic("java/lang/Integer.BYTES", Static{Type: types.Int, Value: int64(4)})
	_ = AddStatic("java/lang/Integer.MAX_VALUE", Static{Type: types.Int, Value: int64(2147483647)})
	_ = AddStatic("java/lang/Integer.MIN_VALUE", Static{Type: types.Int, Value: int64(-2147483648)})
	_ = AddStatic("java/lang/Integer.SIZE", Static{Type: types.Int, Value: int64(32)})
}

func LoadStaticsLong() {
	_ = AddStatic("java/lang/Long.BYTES", Static{Type: types.Int, Value: int64(8)})
	_ = AddStatic("java/lang/Long.MAX_VALUE", Static{Type: types.Long, Value: int64(9223372036854775807)})
	_ = AddStatic("java/lang/Long.MIN_VALUE", Static{Type: types.Long, Value: int64(-9223372036854775808)})
	_ = AddStatic("java/lang/Long.SIZE", Static{Type: types.Int, Value: int64(64)})
}

func LoadStaticsMath() {
	_ = AddStatic("java/lang/Math.E", Static{Type: types.Double, Value: float64(2.718281828459045)})
	_ = AddStatic("java/lang/Math.PI", Static{Type: types.Double, Value: float64(3.141592653589793)})
}

func LoadStaticsShort() {
	_ = AddStatic("java/lang/Short.BYTES", Static{Type: types.Int, Value: int64(2)})
	_ = AddStatic("java/lang/Short.MAX_VALUE", Static{Type: types.Short, Value: int64(32767)})
	_ = AddStatic("java/lang/Short.MIN_VALUE", Static{Type: types.Short, Value: int64(-32768)})
	_ = AddStatic("java/lang/Short.SIZE", Static{Type: types.Int, Value: int64(16)})
}

func LoadStaticsStrictMath() {
	_ = AddStatic("java/lang/StrictMath.E", Static{Type: types.Double, Value: float64(2.718281828459045)})
	_ = AddStatic("java/lang/StrictMath.PI", Static{Type: types.Double, Value: float64(3.141592653589793)})
}
