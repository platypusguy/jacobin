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

var classNameFormatStyle = "java/time/format/FormatStyle"

var formatStyleDefinitions = []string{
	"FULL",
	"LONG",
	"MEDIUM",
	"SHORT",
}

var formatStyleMutex sync.Mutex
var formatStyleInitialized bool = false
var formatStyleInstances []*object.Object

func ensureFormatStyleInited() {
	formatStyleMutex.Lock()
	defer formatStyleMutex.Unlock()
	if formatStyleInitialized {
		return
	}

	formatStyleInstances = make([]*object.Object, len(formatStyleDefinitions))
	for index, enumName := range formatStyleDefinitions {
		styleObj := object.MakeEmptyObjectWithClassName(&classNameFormatStyle)
		styleObj.FieldTable["name"] = object.Field{
			Ftype:  types.StringClassRef,
			Fvalue: object.StringObjectFromGoString(enumName),
		}
		styleObj.FieldTable["ordinal"] = object.Field{
			Ftype:  types.Int,
			Fvalue: int64(index),
		}

		formatStyleInstances[index] = styleObj
		_ = statics.AddStatic(classNameFormatStyle+"."+enumName, statics.Static{
			Type:  "Ljava/time/format/FormatStyle;",
			Value: styleObj,
		})
	}

	formatStyleInitialized = true
}

func getFormatStyleOrdinal(targetObject *object.Object) (int64, bool) {
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
			for index, enumName := range formatStyleDefinitions {
				if enumName == goName {
					return int64(index), true
				}
			}
		}
	}
	return 0, false
}

// Load_Time_Format_FormatStyle registers all java/time/format/FormatStyle methods.
func Load_Time_Format_FormatStyle() {
	ghelpers.MethodSignatures["java/time/format/FormatStyle.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  formatStyleClinit,
		}

	ghelpers.MethodSignatures["java/time/format/FormatStyle.<init>(Ljava/lang/String;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  formatStyleInit,
		}

	ghelpers.MethodSignatures["java/time/format/FormatStyle.compareTo(Ljava/lang/Enum;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  formatStyleCompareTo,
		}

	ghelpers.MethodSignatures["java/time/format/FormatStyle.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  formatStyleCompareTo,
		}

	ghelpers.MethodSignatures["java/time/format/FormatStyle.compareTo(Ljava/time/format/FormatStyle;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  formatStyleCompareTo,
		}

	ghelpers.MethodSignatures["java/time/format/FormatStyle.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  formatStyleEquals,
		}

	ghelpers.MethodSignatures["java/time/format/FormatStyle.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  formatStyleHashCode,
		}

	ghelpers.MethodSignatures["java/time/format/FormatStyle.name()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  formatStyleName,
		}

	ghelpers.MethodSignatures["java/time/format/FormatStyle.ordinal()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  formatStyleOrdinal,
		}

	ghelpers.MethodSignatures["java/time/format/FormatStyle.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  formatStyleToString,
		}

	ghelpers.MethodSignatures["java/time/format/FormatStyle.valueOf(Ljava/lang/String;)Ljava/time/format/FormatStyle;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  formatStyleValueOf,
		}

	ghelpers.MethodSignatures["java/time/format/FormatStyle.values()[Ljava/time/format/FormatStyle;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  formatStyleValues,
		}
}

func formatStyleClinit([]any) any {
	ensureFormatStyleInited()
	return nil
}

func formatStyleValues([]any) any {
	ensureFormatStyleInited()
	arrayObject := object.Make1DimRefArray("Ljava/time/format/FormatStyle;", int64(len(formatStyleInstances)))
	elementsSlot := arrayObject.FieldTable["value"].Fvalue.([]*object.Object)
	copy(elementsSlot, formatStyleInstances)
	return arrayObject
}

func formatStyleValueOf(params []any) any {
	ensureFormatStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "FormatStyle.valueOf: name is null")
	}

	targetNameObj := params[0].(*object.Object)
	if !object.IsStringObject(targetNameObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "FormatStyle.valueOf: name is not a string")
	}

	targetName := object.GoStringFromStringObject(targetNameObj)
	for index, enumName := range formatStyleDefinitions {
		if enumName == targetName {
			return formatStyleInstances[index]
		}
	}

	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "FormatStyle.valueOf: unknown enum constant "+targetName)
}

func formatStyleInit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "FormatStyle.<init>: self is null")
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

func formatStyleName(params []any) any {
	ensureFormatStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "FormatStyle.name: self is null")
	}
	self := params[0].(*object.Object)
	if field, exists := self.FieldTable["name"]; exists {
		if nameObj, ok := field.Fvalue.(*object.Object); ok && !object.IsNull(nameObj) {
			return nameObj
		}
	}
	ordinal, ok := getFormatStyleOrdinal(self)
	if !ok || ordinal < 0 || ordinal >= int64(len(formatStyleDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "FormatStyle.name: invalid instance")
	}
	return object.StringObjectFromGoString(formatStyleDefinitions[ordinal])
}

func formatStyleOrdinal(params []any) any {
	ensureFormatStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "FormatStyle.ordinal: self is null")
	}
	self := params[0].(*object.Object)
	ordinal, ok := getFormatStyleOrdinal(self)
	if !ok || ordinal < 0 || ordinal >= int64(len(formatStyleDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "FormatStyle.ordinal: invalid instance")
	}
	return ordinal
}

func formatStyleToString(params []any) any {
	return formatStyleName(params)
}

func formatStyleCompareTo(params []any) any {
	ensureFormatStyleInited()
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "FormatStyle.compareTo: null parameter")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	selfOrdinal, ok1 := getFormatStyleOrdinal(self)
	otherOrdinal, ok2 := getFormatStyleOrdinal(other)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "FormatStyle.compareTo: incompatible type")
	}

	return selfOrdinal - otherOrdinal
}

func formatStyleEquals(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	if self == other {
		return types.JavaBoolTrue
	}

	selfOrdinal, ok1 := getFormatStyleOrdinal(self)
	otherOrdinal, ok2 := getFormatStyleOrdinal(other)
	if ok1 && ok2 && selfOrdinal == otherOrdinal {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func formatStyleHashCode(params []any) any {
	ensureFormatStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "FormatStyle.hashCode: self is null")
	}
	self := params[0].(*object.Object)
	ordinal, ok := getFormatStyleOrdinal(self)
	if !ok || ordinal < 0 || ordinal >= int64(len(formatStyleDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "FormatStyle.hashCode: invalid instance")
	}
	return ordinal
}
