/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"sync"
)

var classNameSignStyle = "java/time/format/SignStyle"

var signStyleDefinitions = []string{
	"NORMAL",
	"ALWAYS",
	"NEVER",
	"NOT_NEGATIVE",
	"EXCEEDS_PAD",
}

var signStyleMutex sync.Mutex
var signStyleInitialized bool = false
var signStyleInstances []*object.Object

func ensureSignStyleInited() {
	signStyleMutex.Lock()
	defer signStyleMutex.Unlock()
	if signStyleInitialized {
		return
	}

	signStyleInstances = make([]*object.Object, len(signStyleDefinitions))
	for index, enumName := range signStyleDefinitions {
		styleObj := object.MakeEmptyObjectWithClassName(&classNameSignStyle)
		styleObj.FieldTable["name"] = object.Field{
			Ftype:  types.StringClassRef,
			Fvalue: object.StringObjectFromGoString(enumName),
		}
		styleObj.FieldTable["ordinal"] = object.Field{
			Ftype:  types.Int,
			Fvalue: int64(index),
		}

		signStyleInstances[index] = styleObj
		_ = statics.AddStatic(classNameSignStyle+"."+enumName, statics.Static{
			Type:  "Ljava/time/format/SignStyle;",
			Value: styleObj,
		})
	}

	signStyleInitialized = true
}

func getSignStyleOrdinal(targetObject *object.Object) (int64, bool) {
	if targetObject == nil || targetObject.FieldTable == nil {
		return 0, false
	}
	if ordinalField, exists := targetObject.FieldTable["ordinal"]; exists {
		if ordinalValue, ok := ghelpers.GetInt64Val(ordinalField.Fvalue); ok {
			return ordinalValue, true
		}
	}
	if nameField, exists := targetObject.FieldTable["name"]; exists {
		if nameStringObject, ok := nameField.Fvalue.(*object.Object); ok && object.IsStringObject(nameStringObject) {
			goName := object.GoStringFromStringObject(nameStringObject)
			for index, enumName := range signStyleDefinitions {
				if enumName == goName {
					return int64(index), true
				}
			}
		}
	}
	return 0, false
}

// Load_Time_Format_SignStyle registers all java/time/format/SignStyle methods.
func Load_Time_Format_SignStyle() {
	ghelpers.MethodSignatures["java/time/format/SignStyle.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  signStyleClinit,
		}

	ghelpers.MethodSignatures["java/time/format/SignStyle.<init>(Ljava/lang/String;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  signStyleInit,
		}

	ghelpers.MethodSignatures["java/time/format/SignStyle.compareTo(Ljava/lang/Enum;)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  signStyleCompareTo,
		}

	ghelpers.MethodSignatures["java/time/format/SignStyle.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  signStyleCompareTo,
		}

	ghelpers.MethodSignatures["java/time/format/SignStyle.compareTo(Ljava/time/format/SignStyle;)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  signStyleCompareTo,
		}

	ghelpers.MethodSignatures["java/time/format/SignStyle.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  signStyleEquals,
		}

	ghelpers.MethodSignatures["java/time/format/SignStyle.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  signStyleHashCode,
		}

	ghelpers.MethodSignatures["java/time/format/SignStyle.name()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  signStyleName,
		}

	ghelpers.MethodSignatures["java/time/format/SignStyle.ordinal()I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  signStyleOrdinal,
		}

	ghelpers.MethodSignatures["java/time/format/SignStyle.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  signStyleToString,
		}

	ghelpers.MethodSignatures["java/time/format/SignStyle.valueOf(Ljava/lang/String;)Ljava/time/format/SignStyle;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  signStyleValueOf,
		}

	ghelpers.MethodSignatures["java/time/format/SignStyle.values()[Ljava/time/format/SignStyle;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  signStyleValues,
		}

	// SignStyle package-private parse/print helper methods in Java time format
	ghelpers.MethodSignatures["java/time/format/SignStyle.parse(ZZZ)Z"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  signStyleParse,
		}
}

func signStyleClinit([]any) any {
	ensureSignStyleInited()
	return nil
}

func signStyleValues([]any) any {
	ensureSignStyleInited()
	arrayObject := object.Make1DimRefArray("Ljava/time/format/SignStyle;", int64(len(signStyleInstances)))
	elementsSlot := arrayObject.FieldTable["value"].Fvalue.([]*object.Object)
	copy(elementsSlot, signStyleInstances)
	return arrayObject
}

func signStyleValueOf(params []any) any {
	ensureSignStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "SignStyle.valueOf: name is null")
	}

	targetNameObj := params[0].(*object.Object)
	if !object.IsStringObject(targetNameObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "SignStyle.valueOf: name is not a string")
	}

	targetName := object.GoStringFromStringObject(targetNameObj)
	for index, enumName := range signStyleDefinitions {
		if enumName == targetName {
			return signStyleInstances[index]
		}
	}

	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "SignStyle.valueOf: unknown enum constant "+targetName)
}

func signStyleInit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "SignStyle.<init>: self is null")
	}
	self := params[0].(*object.Object)
	if len(params) > 1 && params[1] != nil {
		if nameObj, ok := params[1].(*object.Object); ok && !object.IsNull(nameObj) {
			self.FieldTable["name"] = object.Field{Ftype: types.StringClassRef, Fvalue: nameObj}
		}
	}
	if len(params) > 2 {
		if ordVal, ok := ghelpers.GetInt64Val(params[2]); ok {
			self.FieldTable["ordinal"] = object.Field{Ftype: types.Int, Fvalue: ordVal}
		}
	}
	return nil
}

func signStyleName(params []any) any {
	ensureSignStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "SignStyle.name: self is null")
	}
	self := params[0].(*object.Object)
	if field, exists := self.FieldTable["name"]; exists {
		if nameObj, ok := field.Fvalue.(*object.Object); ok && !object.IsNull(nameObj) {
			return nameObj
		}
	}
	ordinal, ok := getSignStyleOrdinal(self)
	if !ok || ordinal < 0 || ordinal >= int64(len(signStyleDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "SignStyle.name: invalid instance")
	}
	return object.StringObjectFromGoString(signStyleDefinitions[ordinal])
}

func signStyleOrdinal(params []any) any {
	ensureSignStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "SignStyle.ordinal: self is null")
	}
	self := params[0].(*object.Object)
	ordinal, ok := getSignStyleOrdinal(self)
	if !ok || ordinal < 0 || ordinal >= int64(len(signStyleDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "SignStyle.ordinal: invalid instance")
	}
	return ordinal
}

func signStyleToString(params []any) any {
	return signStyleName(params)
}

func signStyleCompareTo(params []any) any {
	ensureSignStyleInited()
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "SignStyle.compareTo: null parameter")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	selfOrdinal, ok1 := getSignStyleOrdinal(self)
	otherOrdinal, ok2 := getSignStyleOrdinal(other)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "SignStyle.compareTo: incompatible type")
	}

	return selfOrdinal - otherOrdinal
}

func signStyleEquals(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	if self == other {
		return types.JavaBoolTrue
	}

	selfOrdinal, ok1 := getSignStyleOrdinal(self)
	otherOrdinal, ok2 := getSignStyleOrdinal(other)
	if ok1 && ok2 && selfOrdinal == otherOrdinal {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func signStyleHashCode(params []any) any {
	ensureSignStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "SignStyle.hashCode: self is null")
	}
	self := params[0].(*object.Object)
	ordinal, ok := getSignStyleOrdinal(self)
	if !ok || ordinal < 0 || ordinal >= int64(len(signStyleDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "SignStyle.hashCode: invalid instance")
	}
	return ordinal
}

func toBool(val any) bool {
	if b, ok := val.(bool); ok {
		return b
	}
	if i, ok := ghelpers.GetInt64Val(val); ok {
		return i != types.JavaBoolFalse
	}
	return false
}

// signStyleParse implements boolean parse(boolean positive, boolean strict, boolean fixedWidth)
// In OpenJDK java.time.format.SignStyle:
//
//	boolean parse(boolean positive, boolean strict, boolean fixedWidth) {
//	    switch (ordinal()) {
//	        case 0: // NORMAL
//	            // normal prints sign only if negative, so any positive direction is not strictly parsed
//	            return !positive || !strict;
//	        case 1: // ALWAYS
//	        case 4: // EXCEEDS_PAD
//	            return true;
//	        default:
//	            // NEVER and NOT_NEGATIVE
//	            return !strict && !fixedWidth;
//	    }
//	}
func signStyleParse(params []any) any {
	ensureSignStyleInited()
	if len(params) < 4 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "SignStyle.parse: self is null")
	}
	self := params[0].(*object.Object)
	ordinal, ok := getSignStyleOrdinal(self)
	if !ok || ordinal < 0 || ordinal >= int64(len(signStyleDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "SignStyle.parse: invalid instance")
	}

	positive := toBool(params[1])
	strict := toBool(params[2])
	fixedWidth := toBool(params[3])

	var res bool
	switch ordinal {
	case 0: // NORMAL
		res = !positive || !strict
	case 1, 4: // ALWAYS, EXCEEDS_PAD
		res = true
	default: // NEVER (2), NOT_NEGATIVE (3)
		res = !strict && !fixedWidth
	}

	if res {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}
