/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"strconv"
	"unicode"
)

var classNameCharacter = "java/lang/Character"

func Load_Lang_Character() {

	// ---- class initialization ----

	ghelpers.MethodSignatures["java/lang/Character.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	// ---- constructors ----

	// Deprecated since Java 9
	ghelpers.MethodSignatures["java/lang/Character.<init>(C)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapDeprecated,
		}

	// ---- methods (alphabetical by FQN) ----

	ghelpers.MethodSignatures["java/lang/Character.charCount(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  charCount,
		}

	ghelpers.MethodSignatures["java/lang/Character.charValue()C"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  charValue,
		}

	ghelpers.MethodSignatures["java/lang/Character.codePointAt([CI)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapUnicode,
		}

	ghelpers.MethodSignatures["java/lang/Character.codePointBefore([CI)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapUnicode,
		}

	ghelpers.MethodSignatures["java/lang/Character.codePointCount([CII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapUnicode,
		}

	ghelpers.MethodSignatures["java/lang/Character.compare(CC)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  charCompare,
		}

	ghelpers.MethodSignatures["java/lang/Character.compareTo(Ljava/lang/Character;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Character.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Character.describeConstable()Ljava/util/Optional;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Character.digit(CI)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  charDigit,
		}

	ghelpers.MethodSignatures["java/lang/Character.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  charEquals,
		}

	ghelpers.MethodSignatures["java/lang/Character.forDigit(II)C"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  charForDigit,
		}

	ghelpers.MethodSignatures["java/lang/Character.getNumericValue(C)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Character.getType(C)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Character.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  charHashCode,
		}

	ghelpers.MethodSignatures["java/lang/Character.hashCode(C)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  charHashCodeStatic,
		}

	ghelpers.MethodSignatures["java/lang/Character.isAlphabetic(I)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapUnicode,
		}

	ghelpers.MethodSignatures["java/lang/Character.isDigit(C)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  charIsDigit,
		}

	ghelpers.MethodSignatures["java/lang/Character.isDigit(I)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapUnicode,
		}

	ghelpers.MethodSignatures["java/lang/Character.isHighSurrogate(C)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Character.isJavaIdentifierPart(C)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Character.isJavaIdentifierStart(C)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Character.isLetter(C)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  charIsLetter,
		}

	ghelpers.MethodSignatures["java/lang/Character.isLetter(I)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapUnicode,
		}

	ghelpers.MethodSignatures["java/lang/Character.isLetterOrDigit(C)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Character.isLetterOrDigit(I)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapUnicode,
		}

	ghelpers.MethodSignatures["java/lang/Character.isLowerCase(C)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  charIsLowerCase,
		}

	ghelpers.MethodSignatures["java/lang/Character.isLowSurrogate(C)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Character.isSpaceChar(C)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Character.isSurrogatePair(CC)Z"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Character.isUpperCase(C)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  charIsUpperCase,
		}

	ghelpers.MethodSignatures["java/lang/Character.isWhitespace(C)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  charIsWhitespace,
		}

	ghelpers.MethodSignatures["java/lang/Character.isWhitespace(I)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapUnicode,
		}

	ghelpers.MethodSignatures["java/lang/Character.toLowerCase(C)C"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  charToLowerCase,
		}

	ghelpers.MethodSignatures["java/lang/Character.toLowerCase(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapUnicode,
		}

	ghelpers.MethodSignatures["java/lang/Character.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  charToString,
		}

	ghelpers.MethodSignatures["java/lang/Character.toString(C)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  charToStringStatic,
		}

	ghelpers.MethodSignatures["java/lang/Character.toUpperCase(C)C"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  charToUpperCase,
		}

	ghelpers.MethodSignatures["java/lang/Character.toUpperCase(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapUnicode,
		}

	ghelpers.MethodSignatures["java/lang/Character.valueOf(C)Ljava/lang/Character;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  characterValueOf,
		}
}

func charCount([]interface{}) interface{} {
	// TODO: We only support UTF-8.
	return int64(1)
}

// "java/lang/Character.compare(CC)I"
func charCompare(params []interface{}) interface{} {
	c1 := params[0].(int64)
	c2 := params[1].(int64)
	return c1 - c2
}

// "java/lang/Character.digit(CI)I"
func charDigit(params []interface{}) interface{} {
	codePoint := rune(params[0].(int64))
	radix := int(params[1].(int64))

	if radix < int(ghelpers.MinRadix) || radix > int(ghelpers.MaxRadix) {
		return int64(-1)
	}

	val := -1
	if codePoint >= '0' && codePoint <= '9' {
		val = int(codePoint - '0')
	} else if codePoint >= 'a' && codePoint <= 'z' {
		val = int(codePoint - 'a' + 10)
	} else if codePoint >= 'A' && codePoint <= 'Z' {
		val = int(codePoint - 'A' + 10)
	} else if codePoint >= 0xFF21 && codePoint <= 0xFF3A { // Fullwidth Latin Capital Letter
		val = int(codePoint - 0xFF21 + 10)
	} else if codePoint >= 0xFF41 && codePoint <= 0xFF5A { // Fullwidth Latin Small Letter
		val = int(codePoint - 0xFF41 + 10)
	}

	if val >= 0 && val < radix {
		return int64(val)
	}

	// For other unicode digits, we can use unicode package but it's complex for all radices.

	// Java's Character.digit handles any Unicode digit.
	if unicode.IsDigit(codePoint) {
		// unicode.Digit is not a function in Go's unicode package that returns the numeric value.

		// It's a table. We can use fmt.Sprintf or similar but that's overkill.
		// Actually, for simple cases:
		str := string(codePoint)
		if d, err := strconv.Atoi(str); err == nil {
			if d >= 0 && d < radix {
				return int64(d)
			}
		}
	}

	return int64(-1)
}

// "java/lang/Character.equals(Ljava/lang/Object;)Z"
func charEquals(params []interface{}) interface{} {
	if len(params) != 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "charEquals requires exactly 2 arguments")
	}

	charObj, ok1 := params[0].(*object.Object)
	otherObj, ok2 := params[1].(*object.Object)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "charEquals: Invalid argument types")
	}

	if object.GoStringFromStringPoolIndex(otherObj.KlassName) != classNameCharacter {
		return types.JavaBoolFalse
	}

	charValue := charObj.FieldTable["value"].Fvalue.(int64)
	otherValue := otherObj.FieldTable["value"].Fvalue.(int64)

	if charValue == otherValue {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// "java/lang/Character.forDigit(II)C"
func charForDigit(params []interface{}) interface{} {
	digit := int(params[0].(int64))
	radix := int(params[1].(int64))

	if radix < int(ghelpers.MinRadix) || radix > int(ghelpers.MaxRadix) {
		return int64(0)
	}
	if digit < 0 || digit >= radix {
		return int64(0)
	}
	if digit < 10 {
		return int64('0' + digit)
	}
	return int64('a' + digit - 10)
}

// "java/lang/Character.hashCode()I"
func charHashCode(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	val := self.FieldTable["value"].Fvalue.(int64)
	return charHashCodeStatic([]interface{}{val})
}

// "java/lang/Character.hashCode(C)I"
func charHashCodeStatic(params []interface{}) interface{} {
	val := params[0].(int64)
	return val
}

// "java/lang/Character.isDigit(C)Z"
func charIsDigit(params []interface{}) interface{} {
	ii := params[0].(int64)
	if unicode.IsDigit(rune(ii)) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// "java/lang/Character.isLetter(C)Z"
func charIsLetter(params []interface{}) interface{} {
	ii := params[0].(int64)
	if unicode.IsLetter(rune(ii)) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// "java/lang/Character.isLowerCase(C)Z"
func charIsLowerCase(params []interface{}) interface{} {
	ii := params[0].(int64)
	if unicode.IsLower(rune(ii)) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// "java/lang/Character.isUpperCase(C)Z"
func charIsUpperCase(params []interface{}) interface{} {
	ii := params[0].(int64)
	if unicode.IsUpper(rune(ii)) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// "java/lang/Character.isWhitespace(C)Z"
func charIsWhitespace(params []interface{}) interface{} {
	ii := params[0].(int64)
	// Java Character.isWhitespace has a specific set of characters,
	// Go's unicode.IsSpace is similar but not identical.
	// Java: space, tab, line feed, carriage return, form feed, etc.
	r := rune(ii)
	if unicode.IsSpace(r) {
		return types.JavaBoolTrue
	}
	// Java's isWhitespace also includes some specific non-breaking spaces etc.
	// This is a reasonable approximation for now.
	return types.JavaBoolFalse
}

// "java/lang/Character.toString()Ljava/lang/String;"
func charToString(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	val := self.FieldTable["value"].Fvalue.(int64)
	return charToStringStatic([]interface{}{val})
}

// "java/lang/Character.toString(C)Ljava/lang/String;"
func charToStringStatic(params []interface{}) interface{} {
	val := params[0].(int64)
	str := string(rune(val))
	return object.StringObjectFromGoString(str)
}

// "java/lang/Character.charValue()C"
func charValue(params []interface{}) interface{} {
	var ch int64
	parmObj := params[0].(*object.Object)
	ch = parmObj.FieldTable["value"].Fvalue.(int64)
	return ch
}

// "java/lang/Character.toLowerCase(C)C"
func charToLowerCase(params []interface{}) interface{} {
	ii := params[0].(int64)
	rr := unicode.ToLower(rune(ii))
	return int64(rr)
}

// "java/lang/Character.toUpperCase(C)C"
func charToUpperCase(params []interface{}) interface{} {
	ii := params[0].(int64)
	rr := unicode.ToUpper(rune(ii))
	return int64(rr)
}

// "java/lang/Character.valueOf(C)Ljava/lang/Character;"
func characterValueOf(params []interface{}) interface{} {
	int64Value := params[0].(int64)
	return object.MakePrimitiveObject("java/lang/Character", types.Char, int64Value)
}
