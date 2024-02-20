package gfunction

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/exceptions"
	"jacobin/object"
	"jacobin/types"
	"math"
)

func populator(classname string, fldtype string, value interface{}) interface{} {
	klass := classloader.MethAreaFetch(classname)
	if klass == nil {
		errMsg := fmt.Sprintf("populator: Could not find %s in the MethodArea", classname)
		return getGErrBlk(exceptions.VirtualMachineError, errMsg)
	}
	klass.Data.ClInit = types.ClInitRun // just mark that String.<clinit>() has been run
	objPtr := object.MakePrimitiveObject(classname, fldtype, value)
	(*objPtr).FieldTable["value"] = object.Field{fldtype, value}
	return objPtr
}

func populateByte(objPtr *object.Object, value int64) {
	(*objPtr).FieldTable["BYTES"] = object.Field{"I", 1}
	(*objPtr).FieldTable["MAX_VALUE"] = object.Field{"B", 0x7f}
	(*objPtr).FieldTable["MIN_VALUE"] = object.Field{"B", 0x80}
	(*objPtr).FieldTable["SIZE"] = object.Field{"I", 8}
	(*objPtr).FieldTable["value"] = object.Field{"B", value}
}

func populateCharacter(objPtr *object.Object, value int64) {
	(*objPtr).FieldTable["BYTES"] = object.Field{"I", 2}
	(*objPtr).FieldTable["COMBINING_SPACING_MARK"] = object.Field{"B", 0x8}
	(*objPtr).FieldTable["CONNECTOR_PUNCTUATION"] = object.Field{"B", 0x17}
	(*objPtr).FieldTable["CONTROL"] = object.Field{"B", 0xf}
	(*objPtr).FieldTable["CURRENCY_SYMBOL"] = object.Field{"B", 0x1a}
	(*objPtr).FieldTable["DASH_PUNCTUATION"] = object.Field{"B", 0x14}
	(*objPtr).FieldTable["DECIMAL_DIGIT_NUMBER"] = object.Field{"B", 0x9}
	(*objPtr).FieldTable["DIRECTIONALITY_ARABIC_NUMBER"] = object.Field{"B", 0x6}
	(*objPtr).FieldTable["DIRECTIONALITY_BOUNDARY_NEUTRAL"] = object.Field{"B", 0x9}
	(*objPtr).FieldTable["DIRECTIONALITY_COMMON_NUMBER_SEPARATOR"] = object.Field{"B", 0x7}
	(*objPtr).FieldTable["DIRECTIONALITY_EUROPEAN_NUMBER"] = object.Field{"B", 0x3}
	(*objPtr).FieldTable["DIRECTIONALITY_EUROPEAN_NUMBER_SEPARATOR"] = object.Field{"B", 0x4}
	(*objPtr).FieldTable["DIRECTIONALITY_EUROPEAN_NUMBER_TERMINATOR"] = object.Field{"B", 0x5}
	(*objPtr).FieldTable["DIRECTIONALITY_FIRST_STRONG_ISOLATE"] = object.Field{"B", 0x15}
	(*objPtr).FieldTable["DIRECTIONALITY_LEFT_TO_RIGHT"] = object.Field{"B", 0x0}
	(*objPtr).FieldTable["DIRECTIONALITY_LEFT_TO_RIGHT_EMBEDDING"] = object.Field{"B", 0xe}
	(*objPtr).FieldTable["DIRECTIONALITY_LEFT_TO_RIGHT_ISOLATE"] = object.Field{"B", 0x13}
	(*objPtr).FieldTable["DIRECTIONALITY_LEFT_TO_RIGHT_OVERRIDE"] = object.Field{"B", 0xf}
	(*objPtr).FieldTable["DIRECTIONALITY_NONSPACING_MARK"] = object.Field{"B", 0x8}
	(*objPtr).FieldTable["DIRECTIONALITY_OTHER_NEUTRALS"] = object.Field{"B", 0xd}
	(*objPtr).FieldTable["DIRECTIONALITY_PARAGRAPH_SEPARATOR"] = object.Field{"B", 0xa}
	(*objPtr).FieldTable["DIRECTIONALITY_POP_DIRECTIONAL_FORMAT"] = object.Field{"B", 0x12}
	(*objPtr).FieldTable["DIRECTIONALITY_POP_DIRECTIONAL_ISOLATE"] = object.Field{"B", 0x16}
	(*objPtr).FieldTable["DIRECTIONALITY_RIGHT_TO_LEFT"] = object.Field{"B", 0x1}
	(*objPtr).FieldTable["DIRECTIONALITY_RIGHT_TO_LEFT_ARABIC"] = object.Field{"B", 0x2}
	(*objPtr).FieldTable["DIRECTIONALITY_RIGHT_TO_LEFT_EMBEDDING"] = object.Field{"B", 0x10}
	(*objPtr).FieldTable["DIRECTIONALITY_RIGHT_TO_LEFT_ISOLATE"] = object.Field{"B", 0x14}
	(*objPtr).FieldTable["DIRECTIONALITY_RIGHT_TO_LEFT_OVERRIDE"] = object.Field{"B", 0x11}
	(*objPtr).FieldTable["DIRECTIONALITY_SEGMENT_SEPARATOR"] = object.Field{"B", 0xb}
	(*objPtr).FieldTable["DIRECTIONALITY_UNDEFINED"] = object.Field{"B", 0xff}
	(*objPtr).FieldTable["DIRECTIONALITY_WHITESPACE"] = object.Field{"B", 0xc}
	(*objPtr).FieldTable["ENCLOSING_MARK"] = object.Field{"B", 0x7}
	(*objPtr).FieldTable["END_PUNCTUATION"] = object.Field{"B", 0x16}
	(*objPtr).FieldTable["FINAL_QUOTE_PUNCTUATION"] = object.Field{"B", 0x1e}
	(*objPtr).FieldTable["FORMAT"] = object.Field{"B", 0x10}
	(*objPtr).FieldTable["INITIAL_QUOTE_PUNCTUATION"] = object.Field{"B", 0x1d}
	(*objPtr).FieldTable["LETTER_NUMBER"] = object.Field{"B", 0xa}
	(*objPtr).FieldTable["LINE_SEPARATOR"] = object.Field{"B", 0xd}
	(*objPtr).FieldTable["LOWERCASE_LETTER"] = object.Field{"B", 0x2}
	(*objPtr).FieldTable["MATH_SYMBOL"] = object.Field{"B", 0x19}
	(*objPtr).FieldTable["MAX_CODE_POINT"] = object.Field{"I", 1114111}
	(*objPtr).FieldTable["MAX_HIGH_SURROGATE"] = object.Field{"C", 56319} // '\udbff'
	(*objPtr).FieldTable["MAX_LOW_SURROGATE"] = object.Field{"C", 57343}  // '\udfff'
	(*objPtr).FieldTable["MAX_RADIX"] = object.Field{"I", 36}
	(*objPtr).FieldTable["MAX_SURROGATE"] = object.Field{"C", 57343} // '\udfff'
	(*objPtr).FieldTable["MAX_VALUE"] = object.Field{"C", 65535}     // '\uffff'
	(*objPtr).FieldTable["MIN_CODE_POINT"] = object.Field{"I", 0}
	(*objPtr).FieldTable["MIN_HIGH_SURROGATE"] = object.Field{"C", 55296} // '\ud800'
	(*objPtr).FieldTable["MIN_LOW_SURROGATE"] = object.Field{"C", 56320}  // '\udc00'
	(*objPtr).FieldTable["MIN_RADIX"] = object.Field{"I", 2}
	(*objPtr).FieldTable["MIN_SUPPLEMENTARY_CODE_POINT"] = object.Field{"I", 65536}
	(*objPtr).FieldTable["MIN_SURROGATE"] = object.Field{"C", 55296} // '\ud800'
	(*objPtr).FieldTable["MIN_VALUE"] = object.Field{"C", '\u0000'}
	(*objPtr).FieldTable["MODIFIER_LETTER"] = object.Field{"B", 0x4}
	(*objPtr).FieldTable["MODIFIER_SYMBOL"] = object.Field{"B", 0x1b}
	(*objPtr).FieldTable["NON_SPACING_MARK"] = object.Field{"B", 0x6}
	(*objPtr).FieldTable["OTHER_LETTER"] = object.Field{"B", 0x5}
	(*objPtr).FieldTable["OTHER_NUMBER"] = object.Field{"B", 0xb}
	(*objPtr).FieldTable["OTHER_PUNCTUATION"] = object.Field{"B", 0x18}
	(*objPtr).FieldTable["OTHER_SYMBOL"] = object.Field{"B", 0x1c}
	(*objPtr).FieldTable["PARAGRAPH_SEPARATOR"] = object.Field{"B", 0xe}
	(*objPtr).FieldTable["PRIVATE_USE"] = object.Field{"B", 0x12}
	(*objPtr).FieldTable["SIZE"] = object.Field{"I", 16}
	(*objPtr).FieldTable["SPACE_SEPARATOR"] = object.Field{"B", 0xc}
	(*objPtr).FieldTable["START_PUNCTUATION"] = object.Field{"B", 0x15}
	(*objPtr).FieldTable["SURROGATE"] = object.Field{"B", 0x13}
	(*objPtr).FieldTable["TITLECASE_LETTER"] = object.Field{"B", 0x3}
	(*objPtr).FieldTable["UNASSIGNED"] = object.Field{"B", 0x0}
	(*objPtr).FieldTable["UPPERCASE_LETTER"] = object.Field{"B", 0x1}
	(*objPtr).FieldTable["value"] = object.Field{"C", value}
}

func populateDouble(objPtr *object.Object, value float64) {
	(*objPtr).FieldTable["BYTES"] = object.Field{"I", 8}
	(*objPtr).FieldTable["MAX_EXPONENT"] = object.Field{"I", 1023}
	(*objPtr).FieldTable["MAX_VALUE"] = object.Field{"D", 1.7976931348623157e308}
	(*objPtr).FieldTable["MIN_EXPONENT"] = object.Field{"I", -1022}
	(*objPtr).FieldTable["MIN_NORMAL"] = object.Field{"D", 2.2250738585072014e-308}
	(*objPtr).FieldTable["MIN_VALUE"] = object.Field{"D", 4.9e-324}
	(*objPtr).FieldTable["NaN"] = object.Field{"D", math.NaN()}
	(*objPtr).FieldTable["NEGATIVE_INFINITY"] = object.Field{"D", math.Inf(-1)}
	(*objPtr).FieldTable["POSITIVE_INFINITY"] = object.Field{"D", math.Inf(+1)}
	(*objPtr).FieldTable["SIZE"] = object.Field{"I", 64}
	(*objPtr).FieldTable["value"] = object.Field{"D", value}
}

func populateFloat(objPtr *object.Object, value float64) {
	(*objPtr).FieldTable["BYTES"] = object.Field{"I", 4}
	(*objPtr).FieldTable["MAX_EXPONENT"] = object.Field{"I", 127}
	(*objPtr).FieldTable["MAX_VALUE"] = object.Field{"F", 3.4028234663852886e38}
	(*objPtr).FieldTable["MIN_EXPONENT"] = object.Field{"I", -126}
	(*objPtr).FieldTable["MIN_NORMAL"] = object.Field{"F", 1.1754943508222875e-38}
	(*objPtr).FieldTable["MIN_VALUE"] = object.Field{"F", 1.401298464324817e-45}
	(*objPtr).FieldTable["NaN"] = object.Field{"F", math.NaN()}
	(*objPtr).FieldTable["NEGATIVE_INFINITY"] = object.Field{"F", math.Inf(-1)}
	(*objPtr).FieldTable["POSITIVE_INFINITY"] = object.Field{"F", math.Inf(+1)}
	(*objPtr).FieldTable["SIZE"] = object.Field{"I", 32}
	(*objPtr).FieldTable["value"] = object.Field{"D", value}
}

func populateLong(objPtr *object.Object, value int64) {
	(*objPtr).FieldTable["BYTES"] = object.Field{"I", 8}
	(*objPtr).FieldTable["MAX_VALUE"] = object.Field{"J", 9223372036854775807}
	(*objPtr).FieldTable["MIN_VALUE"] = object.Field{"J", -9223372036854775808}
	(*objPtr).FieldTable["SIZE"] = object.Field{"I", 64}
	(*objPtr).FieldTable["value"] = object.Field{"J", value}
}

func populateShort(objPtr *object.Object, value int64) {
	(*objPtr).FieldTable["BYTES"] = object.Field{"I", 2}
	(*objPtr).FieldTable["MAX_VALUE"] = object.Field{"S", 32767}
	(*objPtr).FieldTable["MIN_VALUE"] = object.Field{"S", -32768}
	(*objPtr).FieldTable["SIZE"] = object.Field{"I", 16}
	(*objPtr).FieldTable["value"] = object.Field{"S", value}
}
