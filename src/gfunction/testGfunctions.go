/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/classloader"
	"jacobin/object"
	"jacobin/types"
)

// This file contains test gfunctions for use in unit tests. They're primarily designed such
// that you specify the variable types passed in and the return value. They do nothing
// but accept the params and return what the signature promises. In the case of a returned
// *object.Object, it's always a pointer to a bare object.
//
// *NOTE*: In your test, call gfunction.CheckTestGfunctionsLoaded() prior to using any
// of the tests below.
//
// For an example using these tests, see: jvm.TestInvokeSpecialGmethodNoParams()

func Load_TestGfunctions() {

	// ==== accepting no params ====
	TestMethodSignatures["java/lang/Object.test()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vl,
		}

	// ==== accepting params ====

	// === returning void
	TestMethodSignatures["java/lang/Object.test(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  iv,
		}

	TestMethodSignatures["java/lang/Object.test(D)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  dv,
		}

	TestMethodSignatures["java/lang/Object.test(Ljava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  lv,
		}

	// === returning int or double

	TestMethodSignatures["java/lang/Object.test(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  ii,
		}

	TestMethodSignatures["java/lang/Object.test(I)D"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  id,
		}

	TestMethodSignatures["java/lang/Object.test(I)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  il,
		}

	// === accepting reference to java/lang/Object and returning something

	TestMethodSignatures["java/lang/Object.test(Ljava/lang/Object;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  li,
		}

	TestMethodSignatures["java/lang/Object.test(Ljava/lang/Object;)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  ll,
		}

	TestMethodSignatures["java/lang/Object.test(Ljava/lang/Object;)D"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  ld,
		}
}

func vl(params []any) any {
	obj := object.MakeEmptyObject()
	obj.KlassName = types.ObjectPoolStringIndex
	return &obj
}

func iv(params []any) any {
	return nil
}

func dv(params []any) any {
	return nil
}

func lv(params []any) any {
	return nil
}

func ii(params []any) any {
	return int64(43)
}

func id(params []any) any {
	return float64(43.43)
}

func il(params []any) any {
	obj := object.MakeEmptyObject()
	obj.KlassName = types.ObjectPoolStringIndex
	return &obj
}

func li(params []any) any {
	return 44
}

func ll(params []any) any {
	obj := object.MakeEmptyObject()
	obj.KlassName = types.ObjectPoolStringIndex
	return &obj
}

func ld(params []any) any {
	return float64(44.44)
}

// Make sure that these test gfunctions have been loaded. Call this
// from the test that invokes one of the test gfunctions in this file
func CheckTestGfunctionsLoaded() {
	if !TestGfunctionsLoaded {
		if classloader.MTable == nil {
			classloader.MTable = make(map[string]classloader.MTentry)
		}
		LoadTestGfunctions(&classloader.MTable)
	}
}
