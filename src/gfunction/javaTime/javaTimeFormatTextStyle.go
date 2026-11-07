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

var classNameTextStyle = "java/time/format/TextStyle"

type textStyleDefinition struct {
	enumName           string
	calendarStyle      int64
	zoneNameStyleIndex int64
}

var textStyleDefinitions = []textStyleDefinition{
	{"FULL", 2, 0},
	{"FULL_STANDALONE", 32770, 0},
	{"SHORT", 1, 1},
	{"SHORT_STANDALONE", 32769, 1},
	{"NARROW", 4, 2},
	{"NARROW_STANDALONE", 32772, 2},
}

var textStyleMutex sync.Mutex
var textStyleInitialized bool = false
var textStyleInstances []*object.Object

func getInt64Val(val any) (int64, bool) {
	return ghelpers.GetInt64Val(val)
}

func ensureTextStyleInited() {
	textStyleMutex.Lock()
	defer textStyleMutex.Unlock()
	if textStyleInitialized {
		return
	}

	textStyleInstances = make([]*object.Object, len(textStyleDefinitions))
	for index, definition := range textStyleDefinitions {
		styleObj := object.MakeEmptyObjectWithClassName(&classNameTextStyle)
		styleObj.FieldTable["name"] = object.Field{
			Ftype:  types.StringClassRef,
			Fvalue: object.StringObjectFromGoString(definition.enumName),
		}
		styleObj.FieldTable["ordinal"] = object.Field{
			Ftype:  types.Int,
			Fvalue: int64(index),
		}
		styleObj.FieldTable["calendarStyle"] = object.Field{
			Ftype:  types.Int,
			Fvalue: definition.calendarStyle,
		}
		styleObj.FieldTable["zoneNameStyleIndex"] = object.Field{
			Ftype:  types.Int,
			Fvalue: definition.zoneNameStyleIndex,
		}

		textStyleInstances[index] = styleObj
		_ = statics.AddStatic(classNameTextStyle+"."+definition.enumName, statics.Static{
			Type:  "Ljava/time/format/TextStyle;",
			Value: styleObj,
		})
	}

	textStyleInitialized = true
}

func getTextStyleOrdinal(targetObject *object.Object) (int64, bool) {
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
			for index, definition := range textStyleDefinitions {
				if definition.enumName == goName {
					return int64(index), true
				}
			}
		}
	}
	return 0, false
}

// Load_Time_Format_TextStyle registers all java/time/format/TextStyle methods.
func Load_Time_Format_TextStyle() {
	ghelpers.MethodSignatures["java/time/format/TextStyle.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  textStyleClinit,
		}

	ghelpers.MethodSignatures["java/time/format/TextStyle.<init>(Ljava/lang/String;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  textStyleInit,
		}

	ghelpers.MethodSignatures["java/time/format/TextStyle.<init>(Ljava/lang/String;III)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  textStyleInit,
		}

	ghelpers.MethodSignatures["java/time/format/TextStyle.asNormal()Ljava/time/format/TextStyle;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  textStyleAsNormal,
		}

	ghelpers.MethodSignatures["java/time/format/TextStyle.asStandalone()Ljava/time/format/TextStyle;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  textStyleAsStandalone,
		}

	ghelpers.MethodSignatures["java/time/format/TextStyle.compareTo(Ljava/lang/Enum;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  textStyleCompareTo,
		}

	ghelpers.MethodSignatures["java/time/format/TextStyle.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  textStyleCompareTo,
		}

	ghelpers.MethodSignatures["java/time/format/TextStyle.compareTo(Ljava/time/format/TextStyle;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  textStyleCompareTo,
		}

	ghelpers.MethodSignatures["java/time/format/TextStyle.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  textStyleEquals,
		}

	ghelpers.MethodSignatures["java/time/format/TextStyle.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  textStyleHashCode,
		}

	ghelpers.MethodSignatures["java/time/format/TextStyle.isStandalone()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  textStyleIsStandalone,
		}

	ghelpers.MethodSignatures["java/time/format/TextStyle.name()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  textStyleName,
		}

	ghelpers.MethodSignatures["java/time/format/TextStyle.ordinal()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  textStyleOrdinal,
		}

	ghelpers.MethodSignatures["java/time/format/TextStyle.toCalendarStyle()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  textStyleToCalendarStyle,
		}

	ghelpers.MethodSignatures["java/time/format/TextStyle.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  textStyleToString,
		}

	ghelpers.MethodSignatures["java/time/format/TextStyle.valueOf(Ljava/lang/String;)Ljava/time/format/TextStyle;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  textStyleValueOf,
		}

	ghelpers.MethodSignatures["java/time/format/TextStyle.values()[Ljava/time/format/TextStyle;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  textStyleValues,
		}

	ghelpers.MethodSignatures["java/time/format/TextStyle.zoneNameStyleIndex()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  textStyleZoneNameStyleIndex,
		}
}

func textStyleClinit([]any) any {
	ensureTextStyleInited()
	return nil
}

func textStyleInit(params []any) any {
	ensureTextStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TextStyle.<init>: self is null")
	}
	self := params[0].(*object.Object)
	if len(params) > 1 {
		if nameStrObj, ok := params[1].(*object.Object); ok && !object.IsNull(nameStrObj) {
			self.FieldTable["name"] = object.Field{Ftype: types.StringClassRef, Fvalue: nameStrObj}
		}
	}
	if len(params) > 2 {
		if ordVal, ok := ghelpers.GetInt64Val(params[2]); ok {
			self.FieldTable["ordinal"] = object.Field{Ftype: types.Int, Fvalue: ordVal}
		}
	}
	if len(params) > 3 {
		if calStyle, ok := ghelpers.GetInt64Val(params[3]); ok {
			self.FieldTable["calendarStyle"] = object.Field{Ftype: types.Int, Fvalue: calStyle}
		}
	}
	if len(params) > 4 {
		if zoneIdx, ok := ghelpers.GetInt64Val(params[4]); ok {
			self.FieldTable["zoneNameStyleIndex"] = object.Field{Ftype: types.Int, Fvalue: zoneIdx}
		}
	}
	return nil
}

func textStyleIsStandalone(params []any) any {
	ensureTextStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TextStyle.isStandalone: self is null")
	}
	self := params[0].(*object.Object)
	ordinal, ok := getTextStyleOrdinal(self)
	if !ok || ordinal < 0 || ordinal >= int64(len(textStyleDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TextStyle.isStandalone: invalid instance")
	}
	if (ordinal & 1) == 1 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func textStyleAsStandalone(params []any) any {
	ensureTextStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TextStyle.asStandalone: self is null")
	}
	self := params[0].(*object.Object)
	ordinal, ok := getTextStyleOrdinal(self)
	if !ok || ordinal < 0 || ordinal >= int64(len(textStyleDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TextStyle.asStandalone: invalid instance")
	}
	standaloneOrdinal := ordinal | 1
	return textStyleInstances[standaloneOrdinal]
}

func textStyleAsNormal(params []any) any {
	ensureTextStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TextStyle.asNormal: self is null")
	}
	self := params[0].(*object.Object)
	ordinal, ok := getTextStyleOrdinal(self)
	if !ok || ordinal < 0 || ordinal >= int64(len(textStyleDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TextStyle.asNormal: invalid instance")
	}
	normalOrdinal := ordinal & ^int64(1)
	return textStyleInstances[normalOrdinal]
}

func textStyleToCalendarStyle(params []any) any {
	ensureTextStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TextStyle.toCalendarStyle: self is null")
	}
	self := params[0].(*object.Object)
	ordinal, ok := getTextStyleOrdinal(self)
	if !ok || ordinal < 0 || ordinal >= int64(len(textStyleDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TextStyle.toCalendarStyle: invalid instance")
	}
	return textStyleDefinitions[ordinal].calendarStyle
}

func textStyleZoneNameStyleIndex(params []any) any {
	ensureTextStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TextStyle.zoneNameStyleIndex: self is null")
	}
	self := params[0].(*object.Object)
	ordinal, ok := getTextStyleOrdinal(self)
	if !ok || ordinal < 0 || ordinal >= int64(len(textStyleDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TextStyle.zoneNameStyleIndex: invalid instance")
	}
	return textStyleDefinitions[ordinal].zoneNameStyleIndex
}

func textStyleName(params []any) any {
	ensureTextStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TextStyle.name: self is null")
	}
	self := params[0].(*object.Object)
	if self.FieldTable != nil {
		if fieldEntry, exists := self.FieldTable["name"]; exists {
			if stringObject, okStr := fieldEntry.Fvalue.(*object.Object); okStr && object.IsStringObject(stringObject) {
				return stringObject
			}
		}
	}
	return object.StringObjectFromGoString("")
}

func textStyleOrdinal(params []any) any {
	ensureTextStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TextStyle.ordinal: self is null")
	}
	self := params[0].(*object.Object)
	ordinalValue, okOrdinal := getTextStyleOrdinal(self)
	if okOrdinal {
		return ordinalValue
	}
	return int64(0)
}

func textStyleToString(params []any) any {
	ensureTextStyleInited()
	return textStyleName(params)
}

func textStyleValueOf(params []any) any {
	ensureTextStyleInited()
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TextStyle.valueOf: missing argument")
	}
	stringObject, okObj := params[0].(*object.Object)
	if !okObj {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TextStyle.valueOf: argument is not a String")
	}
	if stringObject == nil || object.IsNull(stringObject) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TextStyle.valueOf: name is null")
	}
	if !object.IsStringObject(stringObject) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TextStyle.valueOf: argument is not a String")
	}
	targetName := object.GoStringFromStringObject(stringObject)
	for index, definition := range textStyleDefinitions {
		if definition.enumName == targetName {
			return textStyleInstances[index]
		}
	}
	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "No enum constant java.time.format.TextStyle."+targetName)
}

func textStyleValues([]any) any {
	ensureTextStyleInited()
	arrayObject := object.Make1DimRefArray("Ljava/time/format/TextStyle;", int64(len(textStyleInstances)))
	elementsSlot := arrayObject.FieldTable["value"].Fvalue.([]*object.Object)
	copy(elementsSlot, textStyleInstances)
	return arrayObject
}

func textStyleCompareTo(params []any) any {
	ensureTextStyleInited()
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TextStyle.compareTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TextStyle.compareTo: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	selfOrdinal, okSelf := getTextStyleOrdinal(self)
	otherOrdinal, okOther := getTextStyleOrdinal(other)
	if !okSelf || !okOther {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "TextStyle.compareTo: invalid enum type")
	}
	if selfOrdinal < otherOrdinal {
		return int64(-1)
	} else if selfOrdinal > otherOrdinal {
		return int64(1)
	}
	return int64(0)
}

func textStyleEquals(params []any) any {
	ensureTextStyleInited()
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return types.JavaBoolFalse
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	if self == other {
		return types.JavaBoolTrue
	}

	selfOrdinal, okSelf := getTextStyleOrdinal(self)
	otherOrdinal, okOther := getTextStyleOrdinal(other)
	if okSelf && okOther && selfOrdinal == otherOrdinal {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func textStyleHashCode(params []any) any {
	ensureTextStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TextStyle.hashCode: self is null")
	}
	self := params[0].(*object.Object)
	ordinal, ok := getTextStyleOrdinal(self)
	if ok {
		return ordinal
	}
	return int64(0)
}
