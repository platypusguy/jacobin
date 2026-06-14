/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaNio

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"sync"
)

func Load_Nio_File_FileVisitResult() {
	ghelpers.MethodSignatures["java/nio/file/FileVisitResult.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  fvResultClinit,
		}

	ghelpers.MethodSignatures["java/nio/file/FileVisitResult.valueOf(Ljava/lang/String;)Ljava/nio/file/FileVisitResult;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  fvResultValueOfString,
		}

	ghelpers.MethodSignatures["java/nio/file/FileVisitResult.values()[Ljava/nio/file/FileVisitResult;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  fvResultValues,
		}
}

var fvResultMutex = sync.Mutex{}
var fvResultOnceInitialized bool = false
var fvResultClassName = "java/nio/file/FileVisitResult"
var fvResultNames = []string{"CONTINUE", "TERMINATE", "SKIP_SUBTREE", "SKIP_SIBLINGS"}
var fvResultInstances []*object.Object

func ensureFvResultInited() {
	fvResultMutex.Lock()
	defer fvResultMutex.Unlock()
	if fvResultOnceInitialized {
		return
	}
	fvResultInstances = make([]*object.Object, len(fvResultNames))
	for i, nm := range fvResultNames {
		obj := object.MakeEmptyObjectWithClassName(&fvResultClassName)
		obj.FieldTable["name"] = object.Field{Ftype: types.StringClassRef, Fvalue: object.StringObjectFromGoString(nm)}
		obj.FieldTable["ordinal"] = object.Field{Ftype: types.Int, Fvalue: int64(i)}
		fvResultInstances[i] = obj
		_ = statics.AddStatic(fvResultClassName+"."+nm, statics.Static{Type: "Ljava/nio/file/FileVisitResult;", Value: obj})
	}
	fvResultOnceInitialized = true
}

func fvResultClinit([]interface{}) interface{} {
	ensureFvResultInited()
	return nil
}

func fvResultValueOfString(params []interface{}) interface{} {
	ensureFvResultInited()
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "FileVisitResult.valueOf(String): missing argument")
	}
	strObj, ok := params[0].(*object.Object)
	if !ok || object.IsNull(strObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "FileVisitResult.valueOf(String): name is null")
	}
	if !object.IsStringObject(strObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "FileVisitResult.valueOf(String): argument is not a String")
	}
	name := object.GoStringFromStringObject(strObj)
	for i, nm := range fvResultNames {
		if nm == name {
			return fvResultInstances[i]
		}
	}
	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "FileVisitResult.valueOf(String): no enum constant "+name)
}

func fvResultValues(params []interface{}) interface{} {
	ensureFvResultInited()
	arr := object.Make1DimRefArray("Ljava/nio/file/FileVisitResult;", int64(len(fvResultInstances)))
	slot := arr.FieldTable["value"].Fvalue.([]*object.Object)
	copy(slot, fvResultInstances)
	arr.FieldTable["value"] = object.Field{Ftype: types.RefArray + "Ljava/nio/file/FileVisitResult;", Fvalue: slot}
	return arr
}
