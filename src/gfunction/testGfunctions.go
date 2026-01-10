/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
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

	ghelpers.TestMethodSignatures["jacobin/src/test/Object.test()D"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  vd,
		}

	ghelpers.TestMethodSignatures["jacobin/src/test/Object.test()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  vl,
		}

	// ==== accepting params ====

	// === returning void
	ghelpers.TestMethodSignatures["jacobin/src/test/Object.test(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  iv,
		}

	ghelpers.TestMethodSignatures["jacobin/src/test/Object.test(D)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  dv,
		}

	ghelpers.TestMethodSignatures["jacobin/src/test/Object.test(Ljava/lang/Object;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  lv,
		}

	// === returning int or double

	ghelpers.TestMethodSignatures["jacobin/src/test/Object.test(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ii,
		}

	ghelpers.TestMethodSignatures["jacobin/src/test/Object.test(I)D"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  id,
		}

	ghelpers.TestMethodSignatures["jacobin/src/test/Object.test(I)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  il,
		}

	// === accepting reference to java/lang/Object and returning something

	ghelpers.TestMethodSignatures["jacobin/src/test/Object.test(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  li,
		}

	ghelpers.TestMethodSignatures["jacobin/src/test/Object.test(Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ll,
		}

	ghelpers.TestMethodSignatures["jacobin/src/test/Object.test(Ljava/lang/Object;)D"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ld,
		}

	// === return error block ===
	ghelpers.TestMethodSignatures["jacobin/src/test/Object.test(D)E"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ie,
		}
}

func vd(params []any) any {
	return float64(41.41)
}

func vl(params []any) any {
	obj := object.MakeEmptyObject()
	obj.KlassName = types.StringPoolObjectIndex
	return obj
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
	obj.KlassName = types.StringPoolObjectIndex
	return obj
}

func li(params []any) any {
	return 44
}

func ll(params []any) any {
	obj := object.MakeEmptyObject()
	obj.KlassName = types.StringPoolObjectIndex
	return obj
}

func ld(params []any) any {
	return float64(44.44)
}

func ie(params []any) any {
	geb := ghelpers.GErrBlk{excNames.InternalException, "intended return of test error"}
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
	classloader.MethAreaInsert("jacobin/src/test/Object", &klass)

	LoadTestGfunctions(&classloader.MTable)
}
