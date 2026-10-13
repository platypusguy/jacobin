/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"strconv"
)

func Load_Java_Lang_AssertionError() {

	ghelpers.MethodSignatures["java/lang/AssertionError.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots:   0,
			GFunction:    initAssertionErrorNull,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["java/lang/AssertionError.<init>(Z)V"] =
		ghelpers.GMeth{
			ParamSlots:   1,
			GFunction:    initAssertionErrorBoolean,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["java/lang/AssertionError.<init>(Ljava/lang/Object;)V"] =
		ghelpers.GMeth{
			ParamSlots:   1,
			GFunction:    initAssertionErrorObject,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["java/lang/AssertionError.<init>(C)V"] =
		ghelpers.GMeth{
			ParamSlots:   1,
			GFunction:    initAssertionErrorInt64,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["java/lang/AssertionError.<init>(I)V"] =
		ghelpers.GMeth{
			ParamSlots:   1,
			GFunction:    initAssertionErrorInt64,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["java/lang/AssertionError.<init>(J)V"] =
		ghelpers.GMeth{
			ParamSlots:   2,
			GFunction:    initAssertionErrorInt64,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["java/lang/AssertionError.<init>(F)V"] =
		ghelpers.GMeth{
			ParamSlots:   1,
			GFunction:    initAssertionErrorFloat,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["java/lang/AssertionError.<init>(D)V"] =
		ghelpers.GMeth{
			ParamSlots:   2,
			GFunction:    initAssertionErrorDouble,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["java/lang/AssertionError.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots:   1,
			GFunction:    initAssertionErrorString,
			NeedsContext: true,
		}

	ghelpers.MethodSignatures["java/lang/AssertionError.<init>(Ljava/lang/String;Ljava/lang/Throwable;)V"] =
		ghelpers.GMeth{
			ParamSlots:   2,
			GFunction:    ghelpers.TrapFunction,
			NeedsContext: true,
		}

}

func initAssertionErrorNull(params []any) any {
	return ThrowableInitNull(params)
}

func initAssertionErrorBoolean(params []any) any {
	javaBool, ok := params[2].(types.JavaBool)
	if !ok {
		return ghelpers.WrongArgType("initAssertionErrorBoolean", "a boolean argument", params[2])
	}
	params[2] = object.StringObjectFromGoString(strconv.FormatBool(javaBool == types.JavaBoolTrue))
	return throwableInitString(params)
}

// AssertionError(char): String.valueOf(char). Chars are carried as int64.
func initAssertionErrorInt64(params []any) any {
	ii, ok := params[2].(int64)
	if !ok {
		return ghelpers.WrongArgType("initAssertionErrorInt64", "an int64 argument", params[2])
	}
	params[2] = object.StringObjectFromGoString(strconv.Itoa(int(ii)))
	return throwableInitString(params)
}

// AssertionError(float): formatted at float32 precision, Java style.
func initAssertionErrorFloat(params []any) any {
	ff, ok := params[2].(float64)
	if !ok {
		return ghelpers.WrongArgType("initAssertionErrorFloat", "a float64 argument", params[2])
	}
	params[2] = object.StringObjectFromGoString(ghelpers.JavaFloatingPointToString(ff, 32))
	return throwableInitString(params)
}

// AssertionError(double)
func initAssertionErrorDouble(params []any) any {
	ff, ok := params[2].(float64)
	if !ok {
		return ghelpers.WrongArgType("initAssertionErrorDouble", "a float64 argument", params[2])
	}
	params[2] = object.StringObjectFromGoString(ghelpers.JavaFloatingPointToString(ff, 64))
	return throwableInitString(params)
}

// AssertionError(Object)
func initAssertionErrorObject(params []any) any {
	obj, ok := params[2].(*object.Object)
	if !ok {
		return ghelpers.WrongArgType("initAssertionErrorObject", "a Java object argument", params[2])
	}
	params[2] = object.StringifyAnythingJava(obj)
	return throwableInitString(params)
}

// AssertionError(Object)
func initAssertionErrorString(params []any) any {
	obj, ok := params[2].(*object.Object)
	if !ok || obj.KlassName != types.StringPoolStringIndex {
		return ghelpers.WrongArgType("initAssertionErrorString", "a Java String argument", params[2])
	}
	return throwableInitString(params)
}
