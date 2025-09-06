/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/src/object"
	"jacobin/src/types"
)

func Load_Math_SimpleDateFormat() {

	MethodSignatures["java/text/SimpleDateFormat.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/text/SimpleDateFormat.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  sdfInit,
		}

	MethodSignatures["java/text/SimpleDateFormat.<init>(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  sdfInitString,
		}

	MethodSignatures["java/text/SimpleDateFormat.<init>(Ljava/lang/String;Ljava/text/DateFormatSymbols;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/text/SimpleDateFormat.applyLocalizedPattern(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  sdfApplyPattern,
		}

	MethodSignatures["java/text/SimpleDateFormat.applyPattern(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  sdfApplyPattern,
		}

	MethodSignatures["java/text/SimpleDateFormat.clone()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  sdfClone,
		}

	MethodSignatures["java/text/SimpleDateFormat.format(Ljava/util/Date;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/text/SimpleDateFormat.format(Ljava/util/Date;Ljava/lang/StringBuffer;Ljava/text/FieldPosition;)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/text/SimpleDateFormat.parse(Ljava/lang/String;)Ljava/util/Date;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/text/SimpleDateFormat.parse(Ljava/lang/String;Ljava/text/ParsePosition;)Ljava/util/Date;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/text/SimpleDateFormat.toLocalizedPattern()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  sdfToPattern,
		}

	MethodSignatures["java/text/SimpleDateFormat.toPattern()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  sdfToPattern,
		}
}

// === SimpleDateFormat minimal implementations ===

// sdfInit initializes a new SimpleDateFormat instance with no pattern.
// Per project constructor convention, returns nil (void).
func sdfInit(params []interface{}) interface{} {
	if len(params) < 1 {
		return nil
	}
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return nil
	}
	// Initialize empty field table.
	obj.FieldTable = make(map[string]object.Field)
	return nil
}

// sdfInitString initializes a new SimpleDateFormat with a pattern String.
// Stores the provided pattern object in a "pattern" field for potential later use.
func sdfInitString(params []interface{}) interface{} {
	if len(params) < 2 {
		return sdfInit(params)
	}
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return nil
	}
	obj.FieldTable = make(map[string]object.Field)
	// Accept either a Java String object or null; store as-is.
	s, _ := params[1].(*object.Object)
	obj.FieldTable["pattern"] = object.Field{Ftype: types.StringClassRef, Fvalue: s}
	return nil
}

// sdfClone returns a shallow clone; minimal implementation returns the same object reference.
// Many call sites only require an Object to be returned; deeper cloning is out of scope here.
func sdfClone(params []interface{}) interface{} {
	if len(params) < 1 {
		return object.Null
	}
	obj, _ := params[0].(*object.Object)
	if obj == nil {
		return object.Null
	}
	return obj
}

// sdfToPattern returns the stored pattern String if present.
// Behavior:
// - If a "pattern" field exists and is a non-nil String object, return it.
// - If the field exists but is null, return Java null.
// - If no pattern field exists (e.g., default ctor), return an empty String.
func sdfToPattern(params []interface{}) interface{} {
	if len(params) < 1 {
		return object.Null
	}
	obj, _ := params[0].(*object.Object)
	if obj == nil {
		return object.Null
	}
	if fld, ok := obj.FieldTable["pattern"]; ok {
		so, _ := fld.Fvalue.(*object.Object)
		if so == nil {
			return object.Null
		}
		return so
	}
	// No pattern stored; return empty string as minimal placeholder.
	return object.StringObjectFromGoString("")
}

// sdfApplyPattern sets or replaces the stored pattern String.
// Behavior:
// - Stores the provided String object (or null) into the "pattern" field.
// - Initializes FieldTable if needed.
// - Returns nil (void).
func sdfApplyPattern(params []interface{}) interface{} {
	if len(params) < 2 {
		return nil
	}
	obj, _ := params[0].(*object.Object)
	if obj == nil {
		return nil
	}
	if obj.FieldTable == nil {
		obj.FieldTable = make(map[string]object.Field)
	}
	s, _ := params[1].(*object.Object)
	obj.FieldTable["pattern"] = object.Field{Ftype: types.StringClassRef, Fvalue: s}
	return nil
}
