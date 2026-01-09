/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
	"strings"
	"time"
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
			GFunction:  sdfFormat,
		}

	MethodSignatures["java/text/SimpleDateFormat.format(Ljava/util/Date;Ljava/lang/StringBuffer;Ljava/text/FieldPosition;)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/text/SimpleDateFormat.parse(Ljava/lang/String;)Ljava/util/Date;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  sdfParse,
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

// sdfFormat formats a Date into a date/time string.
func sdfFormat(params []interface{}) interface{} {
	if len(params) < 2 {
		return getGErrBlk(excNames.IllegalArgumentException, "sdfFormat: missing parameters")
	}
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return getGErrBlk(excNames.IllegalArgumentException, "sdfFormat: self is not an object")
	}
	dateObj, ok := params[1].(*object.Object)
	if !ok || dateObj == nil {
		return getGErrBlk(excNames.NullPointerException, "sdfFormat: date parameter is null")
	}

	// Get pattern
	javaPattern := ""
	if fld, exists := obj.FieldTable["pattern"]; exists {
		if so, ok := fld.Fvalue.(*object.Object); ok && so != nil {
			javaPattern = object.GoStringFromStringObject(so)
		}
	}
	goLayout := javaToGoDateFormat(javaPattern)

	// Get milliseconds from Date
	millis, err := dateGetMillis(dateObj)
	if err != nil {
		return err
	}

	t := time.UnixMilli(millis).UTC()
	formatted := t.Format(goLayout)
	return object.StringObjectFromGoString(formatted)
}

// sdfParse parses text from the beginning of the given string to produce a date.
func sdfParse(params []interface{}) interface{} {
	if len(params) < 2 {
		return getGErrBlk(excNames.IllegalArgumentException, "sdfParse: missing parameters")
	}
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return getGErrBlk(excNames.IllegalArgumentException, "sdfParse: self is not an object")
	}
	strObj, ok := params[1].(*object.Object)
	if !ok || strObj == nil {
		return getGErrBlk(excNames.NullPointerException, "sdfParse: string parameter is null")
	}

	// Get pattern
	javaPattern := ""
	if fld, exists := obj.FieldTable["pattern"]; exists {
		if so, ok := fld.Fvalue.(*object.Object); ok && so != nil {
			javaPattern = object.GoStringFromStringObject(so)
		}
	}
	goLayout := javaToGoDateFormat(javaPattern)

	inputStr := object.GoStringFromStringObject(strObj)
	t, err := time.Parse(goLayout, inputStr)
	if err != nil {
		return getGErrBlk(excNames.ParseException, fmt.Sprintf("sdfParse: failed to parse %q with layout %q: %v", inputStr, goLayout, err))
	}

	// Create new Date object
	return object.MakePrimitiveObject("java/util/Date", types.Long, t.UnixMilli())
}

func javaToGoDateFormat(javaPattern string) string {
	if javaPattern == "" {
		return time.RFC3339 // Default layout if no pattern provided
	}

	// Minimal mapper for common SimpleDateFormat patterns.
	// Longest patterns must come first in strings.Replacer.
	replacer := strings.NewReplacer(
		"yyyy", "2006",
		"yy", "06",
		"MMMM", "January",
		"MMM", "Jan",
		"MM", "01",
		"M", "1",
		"dd", "02",
		"d", "2",
		"HH", "15",
		"mm", "04",
		"m", "4",
		"ss", "05",
		"s", "5",
		"a", "PM",
		"SSS", "000",
		"z", "MST",
		"Z", "-0700",
	)
	return replacer.Replace(javaPattern)
}
