/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
	"strconv"
	"strings"
)

func Load_Math_Math_Context() {

	MethodSignatures["java/math/MathContext.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/math/MathContext.<init>(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  mconInitInt,
		}

	MethodSignatures["java/math/MathContext.<init>(ILjava/math/RoundingMode;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  mconInitIntRoundingMode,
		}

	MethodSignatures["java/math/MathContext.<init>(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  mconInitString,
		}

	MethodSignatures["java/math/MathContext.equals(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/math/MathContext.getPrecision()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  mconGetPrecision,
		}

	MethodSignatures["java/math/MathContext.getRoundingMode()Ljava/math/RoundingMode;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  mconGetRoundingMode,
		}

	MethodSignatures["java/math/MathContext.hashCode()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/math/MathContext.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  mconToString,
		}

}

// ---------------- Implementations -----------------

// mconInitInt implements MathContext.<init>(int)
// Default rounding mode per JDK is HALF_UP when not specified.
func mconInitInt(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	prec := params[1].(int64)
	if prec < 0 {
		return getGErrBlk(excNames.IllegalArgumentException, "MathContext.<init>(int): negative precision")
	}
	// Ensure RoundingMode constants exist
	ensureRoundingModeInited()
	// Default HALF_UP is ordinal 4
	rm := rmodeInstances[4]
	self.FieldTable["precision"] = object.Field{Ftype: types.Int, Fvalue: prec}
	self.FieldTable["roundingMode"] = object.Field{Ftype: "Ljava/math/RoundingMode;", Fvalue: rm}
	return nil
}

// mconInitIntRoundingMode implements MathContext.<init>(int, RoundingMode)
func mconInitIntRoundingMode(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	prec := params[1].(int64)
	rmode := params[2].(*object.Object)
	if prec < 0 {
		return getGErrBlk(excNames.IllegalArgumentException, "MathContext.<init>(int,RoundingMode): negative precision")
	}
	if object.IsNull(rmode) {
		return getGErrBlk(excNames.NullPointerException, "MathContext.<init>(int,RoundingMode): roundingMode is null")
	}
	self.FieldTable["precision"] = object.Field{Ftype: types.Int, Fvalue: prec}
	self.FieldTable["roundingMode"] = object.Field{Ftype: "Ljava/math/RoundingMode;", Fvalue: rmode}
	return nil
}

// mconInitString implements MathContext.<init>(String)
// Accepts strings like: "precision=3 roundingMode=HALF_UP" (case-insensitive keys).
// If roundingMode is omitted, defaults to HALF_UP.
func mconInitString(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	strObj := params[1].(*object.Object)
	if object.IsNull(strObj) {
		return getGErrBlk(excNames.NullPointerException, "MathContext.<init>(String): null")
	}
	s := object.GoStringFromStringObject(strObj)
	s = strings.TrimSpace(s)
	if s == "" {
		return getGErrBlk(excNames.IllegalArgumentException, "MathContext.<init>(String): empty string")
	}
	// Tokenize on spaces and commas
	replacer := strings.NewReplacer(",", " ")
	s = replacer.Replace(s)
	parts := strings.Fields(s)
	var (
		precSet bool
		precVal int64
		rmName string
	)
	for _, p := range parts {
		if !strings.Contains(p, "=") {
			continue
		}
		kv := strings.SplitN(p, "=", 2)
		key := strings.ToLower(strings.TrimSpace(kv[0]))
		val := strings.TrimSpace(kv[1])
		switch key {
		case "precision":
			if val == "" {
				return getGErrBlk(excNames.IllegalArgumentException, "MathContext.<init>(String): missing precision value")
			}
			i, err := strconv.Atoi(val)
			if err != nil {
				return getGErrBlk(excNames.IllegalArgumentException, "MathContext.<init>(String): invalid precision")
			}
			if i < 0 {
				return getGErrBlk(excNames.IllegalArgumentException, "MathContext.<init>(String): negative precision")
			}
			precVal = int64(i)
			precSet = true
		case "roundingmode":
			rmName = val
		}
	}
	if !precSet {
		return getGErrBlk(excNames.IllegalArgumentException, "MathContext.<init>(String): precision not specified")
	}
	ensureRoundingModeInited()
	var rmodeObj *object.Object
	if rmName == "" {
		// default
		rmodeObj = rmodeInstances[4] // HALF_UP
	} else {
		// Use valueOf(String) to resolve and validate name
		ret := rmodeValueOfString([]interface{}{object.StringObjectFromGoString(rmName)})
		if blk, ok := ret.(*GErrBlk); ok {
			return blk
		}
		rmodeObj = ret.(*object.Object)
	}
	self.FieldTable["precision"] = object.Field{Ftype: types.Int, Fvalue: precVal}
	self.FieldTable["roundingMode"] = object.Field{Ftype: "Ljava/math/RoundingMode;", Fvalue: rmodeObj}
	return nil
}

// mconGetPrecision returns the precision field
func mconGetPrecision(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	fld := self.FieldTable["precision"]
	if v, ok := fld.Fvalue.(int64); ok {
		return v
	}
	// default 0 if absent/malformed
	return int64(0)
}

// mconGetRoundingMode returns the roundingMode field
func mconGetRoundingMode(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	fld, ok := self.FieldTable["roundingMode"]
	if !ok {
		return object.Null
	}
	if obj, ok := fld.Fvalue.(*object.Object); ok {
		return obj
	}
	return object.Null
}

// mconToString returns "precision=<n> roundingMode=<NAME>"
func mconToString(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	prec := mconGetPrecision([]interface{}{self}).(int64)
	rmObj := mconGetRoundingMode([]interface{}{self}).(*object.Object)
	name := "null"
	if !object.IsNull(rmObj) {
		if nameFld, ok := rmObj.FieldTable["name"]; ok {
			if sObj, ok := nameFld.Fvalue.(*object.Object); ok {
				name = object.GoStringFromStringObject(sObj)
			}
		}
	}
	str := "precision=" + strconv.FormatInt(prec, 10) + " roundingMode=" + name
	return object.StringObjectFromGoString(str)
}
