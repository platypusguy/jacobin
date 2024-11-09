/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/types"
)

// This file contains test gfunctions for use in unit tests. They're primarily designed such
// that you specify the variable types passed in and the returned value. They do nothing
// but accept the params and return what the signature promises. In the case of a returned
// *object.Object, it's always a pointer to a bare java/lang/Object.
//
// *IMPORTANT*: Before calling any of these functions, run gfunction.CheckTestGfunctionsLoaded()
// prior to using the tests below. It needs to be called only once per unit test.
//
// For an example using these tests, see: jvm.TestInvokeSpecialGmethodNoParams()

func Load_TestGfunctions() {

	// ==== accepting no params ====

	TestMethodSignatures["jacobin/test/Object.test()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vd,
		}

	TestMethodSignatures["jacobin/test/Object.test()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vl,
		}

	// ==== accepting params ====

	// === returning void
	TestMethodSignatures["jacobin/test/Object.test(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  iv,
		}

	TestMethodSignatures["jacobin/test/Object.test(D)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  dv,
		}

	TestMethodSignatures["jacobin/test/Object.test(Ljava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  lv,
		}

	// === returning int or double

	TestMethodSignatures["jacobin/test/Object.test(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  ii,
		}

	TestMethodSignatures["jacobin/test/Object.test(I)D"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  id,
		}

	TestMethodSignatures["jacobin/test/Object.test(I)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  il,
		}

	// === accepting reference to java/lang/Object and returning something

	TestMethodSignatures["jacobin/test/Object.test(Ljava/lang/Object;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  li,
		}

	TestMethodSignatures["jacobin/test/Object.test(Ljava/lang/Object;)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  ll,
		}

	TestMethodSignatures["jacobin/test/Object.test(Ljava/lang/Object;)D"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  ld,
		}

	// === return error block ===
	TestMethodSignatures["jacobin/test/Object.test(D)E"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  ie,
		}
}

func vd(params []any) any {
	return float64(41.41)
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

func ie(params []any) any {
	geb := GErrBlk{excNames.InternalException, "intended return of test error"}
	return &geb
}

// Make sure that these test gfunctions have been loaded. Call this
// from the test that invokes one of the test gfunctions in this file
func CheckTestGfunctionsLoaded() {
	if classloader.MTable == nil {
		classloader.MTable = make(map[string]classloader.MTentry)
	}

	// in order to load the test functions, there needs to be an object-like entry
	// in the method table, which is handled next. After which, we load the test functions
	klass := classloader.Klass{
		Status: 'Z',
		Loader: "test",
		Data:   nil,
	}
	classloader.MethAreaInsert("jacobin/test/Object", &klass)

	LoadTestGfunctions(&classloader.MTable)
}
