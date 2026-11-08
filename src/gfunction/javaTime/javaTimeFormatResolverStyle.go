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

var classNameResolverStyle = "java/time/format/ResolverStyle"

var resolverStyleDefinitions = []string{
	"STRICT",
	"SMART",
	"LENIENT",
}

var resolverStyleMutex sync.Mutex
var resolverStyleInitialized bool = false
var resolverStyleInstances []*object.Object

func ensureResolverStyleInited() {
	resolverStyleMutex.Lock()
	defer resolverStyleMutex.Unlock()
	if resolverStyleInitialized {
		return
	}

	resolverStyleInstances = make([]*object.Object, len(resolverStyleDefinitions))
	for index, enumName := range resolverStyleDefinitions {
		styleObj := object.MakeEmptyObjectWithClassName(&classNameResolverStyle)
		styleObj.FieldTable["name"] = object.Field{
			Ftype:  types.StringClassRef,
			Fvalue: object.StringObjectFromGoString(enumName),
		}
		styleObj.FieldTable["ordinal"] = object.Field{
			Ftype:  types.Int,
			Fvalue: int64(index),
		}

		resolverStyleInstances[index] = styleObj
		_ = statics.AddStatic(classNameResolverStyle+"."+enumName, statics.Static{
			Type:  "Ljava/time/format/ResolverStyle;",
			Value: styleObj,
		})
	}

	resolverStyleInitialized = true
}

func getResolverStyleOrdinal(targetObject *object.Object) (int64, bool) {
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
			for index, enumName := range resolverStyleDefinitions {
				if enumName == goName {
					return int64(index), true
				}
			}
		}
	}
	return 0, false
}

// Load_Time_Format_ResolverStyle registers all java/time/format/ResolverStyle methods.
func Load_Time_Format_ResolverStyle() {
	ghelpers.MethodSignatures["java/time/format/ResolverStyle.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  resolverStyleClinit,
		}

	ghelpers.MethodSignatures["java/time/format/ResolverStyle.<init>(Ljava/lang/String;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  resolverStyleInit,
		}

	ghelpers.MethodSignatures["java/time/format/ResolverStyle.compareTo(Ljava/lang/Enum;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  resolverStyleCompareTo,
		}

	ghelpers.MethodSignatures["java/time/format/ResolverStyle.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  resolverStyleCompareTo,
		}

	ghelpers.MethodSignatures["java/time/format/ResolverStyle.compareTo(Ljava/time/format/ResolverStyle;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  resolverStyleCompareTo,
		}

	ghelpers.MethodSignatures["java/time/format/ResolverStyle.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  resolverStyleEquals,
		}

	ghelpers.MethodSignatures["java/time/format/ResolverStyle.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  resolverStyleHashCode,
		}

	ghelpers.MethodSignatures["java/time/format/ResolverStyle.name()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  resolverStyleName,
		}

	ghelpers.MethodSignatures["java/time/format/ResolverStyle.ordinal()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  resolverStyleOrdinal,
		}

	ghelpers.MethodSignatures["java/time/format/ResolverStyle.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  resolverStyleToString,
		}

	ghelpers.MethodSignatures["java/time/format/ResolverStyle.valueOf(Ljava/lang/String;)Ljava/time/format/ResolverStyle;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  resolverStyleValueOf,
		}

	ghelpers.MethodSignatures["java/time/format/ResolverStyle.values()[Ljava/time/format/ResolverStyle;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  resolverStyleValues,
		}
}

func resolverStyleClinit([]any) any {
	ensureResolverStyleInited()
	return nil
}

func resolverStyleValues([]any) any {
	ensureResolverStyleInited()
	arrayObject := object.Make1DimRefArray("Ljava/time/format/ResolverStyle;", int64(len(resolverStyleInstances)))
	elementsSlot := arrayObject.FieldTable["value"].Fvalue.([]*object.Object)
	copy(elementsSlot, resolverStyleInstances)
	return arrayObject
}

func resolverStyleValueOf(params []any) any {
	ensureResolverStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ResolverStyle.valueOf: name is null")
	}

	targetNameObj := params[0].(*object.Object)
	if !object.IsStringObject(targetNameObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ResolverStyle.valueOf: name is not a string")
	}

	targetName := object.GoStringFromStringObject(targetNameObj)
	for index, enumName := range resolverStyleDefinitions {
		if enumName == targetName {
			return resolverStyleInstances[index]
		}
	}

	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ResolverStyle.valueOf: unknown enum constant "+targetName)
}

func resolverStyleInit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ResolverStyle.<init>: self is null")
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

func resolverStyleName(params []any) any {
	ensureResolverStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ResolverStyle.name: self is null")
	}
	self := params[0].(*object.Object)
	if field, exists := self.FieldTable["name"]; exists {
		if nameObj, ok := field.Fvalue.(*object.Object); ok && !object.IsNull(nameObj) {
			return nameObj
		}
	}
	ordinal, ok := getResolverStyleOrdinal(self)
	if !ok || ordinal < 0 || ordinal >= int64(len(resolverStyleDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ResolverStyle.name: invalid instance")
	}
	return object.StringObjectFromGoString(resolverStyleDefinitions[ordinal])
}

func resolverStyleOrdinal(params []any) any {
	ensureResolverStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ResolverStyle.ordinal: self is null")
	}
	self := params[0].(*object.Object)
	ordinal, ok := getResolverStyleOrdinal(self)
	if !ok || ordinal < 0 || ordinal >= int64(len(resolverStyleDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ResolverStyle.ordinal: invalid instance")
	}
	return ordinal
}

func resolverStyleToString(params []any) any {
	return resolverStyleName(params)
}

func resolverStyleCompareTo(params []any) any {
	ensureResolverStyleInited()
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ResolverStyle.compareTo: null parameter")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	selfOrdinal, ok1 := getResolverStyleOrdinal(self)
	otherOrdinal, ok2 := getResolverStyleOrdinal(other)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "ResolverStyle.compareTo: incompatible type")
	}

	return selfOrdinal - otherOrdinal
}

func resolverStyleEquals(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	if self == other {
		return types.JavaBoolTrue
	}

	selfOrdinal, ok1 := getResolverStyleOrdinal(self)
	otherOrdinal, ok2 := getResolverStyleOrdinal(other)
	if ok1 && ok2 && selfOrdinal == otherOrdinal {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func resolverStyleHashCode(params []any) any {
	ensureResolverStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ResolverStyle.hashCode: self is null")
	}
	self := params[0].(*object.Object)
	ordinal, ok := getResolverStyleOrdinal(self)
	if !ok || ordinal < 0 || ordinal >= int64(len(resolverStyleDefinitions)) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ResolverStyle.hashCode: invalid instance")
	}
	return ordinal
}
