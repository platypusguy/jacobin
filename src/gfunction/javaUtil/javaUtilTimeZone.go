/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"time"
)

const (
	fieldTZID       = "id"
	fieldRawOffset  = "rawOffset"  // milliseconds
	fieldDSTSavings = "dstSavings" // milliseconds
)

func Load_Util_TimeZone() {
	ghelpers.MethodSignatures["java/util/TimeZone.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  tzInit,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.clone()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  txClone,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.getAvailableIDs()[Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  tzGetAvailableIDs,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.getAvailableIDs(I)[Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  tzGetAvailableIDsInt,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.getDefault()Ljava/util/TimeZone;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.getDisplayName()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  tzGetDisplayName,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.getDisplayName(ZI)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.getDisplayName(ZILjava/util/Locale;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.getDisplayName(Ljava/util/Locale;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.getDSTSavings()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ttzGetDSTSavings,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.getID()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  tzGetID,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.getOffset(IIIII)I"] =
		ghelpers.GMeth{
			ParamSlots: 6,
			GFunction:  tzGetOffset,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.getOffset(J)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  tzGetOffsetLong,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.getRawOffset()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  tzGetRawOffset,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.getTimeZone(Ljava/lang/String;)Ljava/util/TimeZone;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  tzGetTimeZoneString,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.getTimeZone(Ljava/time/ZoneId;)Ljava/util/TimeZone;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.hasSameRules(Ljava/util/TimeZone;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  tzHasSameRules,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.inDaylightTime(Ljava/util/Date;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  tzInDaylightTime,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.observesDaylightTime()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  tzObservesDaylightTime,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.setDefault(Ljava/util/TimeZone;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  tzSetDefault,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.setID(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  tzSetID,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.setRawOffset(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  tzSetRawOffset,
		}

	ghelpers.MethodSignatures["java/util/TimeZone.toZoneId()Ljava/time/ZoneId;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}
}

// --- Helpers ---
func tzGetStringID(obj *object.Object) (string, interface{}) {
	if obj == nil {
		return "", ghelpers.GetGErrBlk(excNames.NullPointerException, "TimeZone object is null")
	}
	fld, ok := obj.FieldTable[fieldTZID]
	if !ok {
		return "", ghelpers.GetGErrBlk(excNames.IllegalStateException, "TimeZone id field missing")
	}
	so, _ := fld.Fvalue.(*object.Object)
	if so == nil {
		return "", nil // treat null as empty
	}
	return object.GoStringFromStringObject(so), nil
}

func tzComputeOffsetsForID(id string) (rawMS int64, dstMS int64) {
	if id == "" || id == "UTC" || id == "GMT" {
		return 0, 0
	}
	loc, err := time.LoadLocation(id)
	if err != nil {
		return 0, 0
	}
	now := time.Now()
	y := now.Year()
	w1 := time.Date(y, 1, 1, 0, 0, 0, 0, loc)
	w2 := time.Date(y, 7, 1, 0, 0, 0, 0, loc)
	_, off1 := w1.Zone()
	_, off2 := w2.Zone()
	std := off1
	if off2 < std {
		std = off2
	}
	// dst savings estimated as absolute difference
	diff := off1 - off2
	if diff < 0 {
		diff = -diff
	}
	return int64(std) * 1000, int64(diff) * 1000
}

func tzCurrentOffsetAt(id string, millis int64) int64 {
	if id == "" || id == "UTC" || id == "GMT" {
		return 0
	}
	loc, err := time.LoadLocation(id)
	if err != nil {
		// fallback to raw offset computed from ID
		raw, _ := tzComputeOffsetsForID(id)
		return raw
	}
	t := time.UnixMilli(millis).In(loc)
	_, off := t.Zone()
	return int64(off) * 1000
}

func tzEnsureFields(obj *object.Object) {
	if obj.FieldTable == nil {
		obj.FieldTable = make(map[string]object.Field)
	}
	if _, ok := obj.FieldTable[fieldTZID]; !ok {
		obj.FieldTable[fieldTZID] = object.Field{Ftype: types.StringClassRef, Fvalue: object.StringObjectFromGoString("UTC")}
	}
	if _, ok := obj.FieldTable[fieldRawOffset]; !ok {
		obj.FieldTable[fieldRawOffset] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
	}
	if _, ok := obj.FieldTable[fieldDSTSavings]; !ok {
		obj.FieldTable[fieldDSTSavings] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
	}
}

// --- Implementations ---

func tzInit(params []interface{}) interface{} {
	obj, _ := params[0].(*object.Object)
	if obj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "tzInit: self not object")
	}
	object.ClearFieldTable(obj)
	obj.FieldTable[fieldTZID] = object.Field{Ftype: types.StringClassRef, Fvalue: object.StringObjectFromGoString("UTC")}
	obj.FieldTable[fieldRawOffset] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
	obj.FieldTable[fieldDSTSavings] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
	return nil
}

// Note: loader references txClone; provide that symbol
func txClone(params []interface{}) interface{} {
	if len(params) < 1 {
		return object.Null
	}
	obj, _ := params[0].(*object.Object)
	if obj == nil {
		return object.Null
	}
	return object.CloneObject(obj)
}

func tzGetAvailableIDs(params []interface{}) interface{} {
	// Minimal list
	ids := []string{"UTC", "GMT"}
	arr := make([]*object.Object, 0, len(ids))
	for _, id := range ids {
		arr = append(arr, object.StringObjectFromGoString(id))
	}
	return object.MakeArrayFromRawArray(arr)
}

func tzGetAvailableIDsInt(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "tzGetAvailableIDsInt: need rawOffset")
	}
	off, _ := params[1].(int64)
	if off == 0 {
		return tzGetAvailableIDs(nil)
	}
	// No other known IDs in minimal implementation
	return object.MakeArrayFromRawArray([]*object.Object{})
}

func tzGetDisplayName(params []interface{}) interface{} {
	obj, _ := params[0].(*object.Object)
	if obj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "tzGetDisplayName: self not object")
	}
	id, err := tzGetStringID(obj)
	if err != nil {
		return err
	}
	return object.StringObjectFromGoString(id)
}

func ttzGetDSTSavings(params []interface{}) interface{} {
	obj, _ := params[0].(*object.Object)
	if obj == nil {
		return int64(0)
	}
	fld := obj.FieldTable[fieldDSTSavings]
	val, ok := fld.Fvalue.(int64)
	if !ok {
		return int64(0)
	}
	return val
}

func tzGetID(params []interface{}) interface{} {
	obj, _ := params[0].(*object.Object)
	if obj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "tzGetID: self not object")
	}
	id, err := tzGetStringID(obj)
	if err != nil {
		return err
	}
	return object.StringObjectFromGoString(id)
}

func tzGetOffset(params []interface{}) interface{} {
	// Minimal: return raw offset only
	obj, _ := params[0].(*object.Object)
	if obj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "tzGetOffset: self not object")
	}
	fld := obj.FieldTable[fieldRawOffset]
	val, ok := fld.Fvalue.(int64)
	if !ok {
		return int64(0)
	}
	return val
}

func tzGetOffsetLong(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "tzGetOffsetLong: need millis")
	}
	obj, _ := params[0].(*object.Object)
	millis, _ := params[1].(int64)
	if obj == nil {
		return int64(0)
	}
	id, err := tzGetStringID(obj)
	if err != nil {
		return err
	}
	if id == "" { // fallback to raw
		fld := obj.FieldTable[fieldRawOffset]
		if v, ok := fld.Fvalue.(int64); ok {
			return v
		}
		return int64(0)
	}
	return tzCurrentOffsetAt(id, millis)
}

func tzGetRawOffset(params []interface{}) interface{} {
	obj, _ := params[0].(*object.Object)
	if obj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "tzGetRawOffset: self not object")
	}
	fld := obj.FieldTable[fieldRawOffset]
	if v, ok := fld.Fvalue.(int64); ok {
		return v
	}
	return int64(0)
}

func tzGetTimeZoneString(params []interface{}) interface{} {
	if len(params) < 1 {
		return object.Null
	}
	// First param is this? For static method, receiver is absent; in our calling convention, params[0] is the ID String
	idObj, _ := params[0].(*object.Object)
	if idObj == nil {
		return object.Null
	}
	id := object.GoStringFromStringObject(idObj)
	className := "java/util/TimeZone"
	obj := object.MakeEmptyObjectWithClassName(&className)
	object.ClearFieldTable(obj)
	raw, dst := tzComputeOffsetsForID(id)
	obj.FieldTable[fieldTZID] = object.Field{Ftype: types.StringClassRef, Fvalue: object.StringObjectFromGoString(id)}
	obj.FieldTable[fieldRawOffset] = object.Field{Ftype: types.Int, Fvalue: raw}
	obj.FieldTable[fieldDSTSavings] = object.Field{Ftype: types.Int, Fvalue: dst}
	return obj
}

func tzHasSameRules(params []interface{}) interface{} {
	obj, _ := params[0].(*object.Object)
	other, _ := params[1].(*object.Object)
	if obj == nil || other == nil {
		return types.JavaBoolFalse
	}
	r1, ok1 := obj.FieldTable[fieldRawOffset].Fvalue.(int64)
	r2, ok2 := other.FieldTable[fieldRawOffset].Fvalue.(int64)
	d1, ok3 := obj.FieldTable[fieldDSTSavings].Fvalue.(int64)
	d2, ok4 := other.FieldTable[fieldDSTSavings].Fvalue.(int64)
	if ok1 && ok2 && ok3 && ok4 && r1 == r2 && d1 == d2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func tzInDaylightTime(params []interface{}) interface{} {
	obj, _ := params[0].(*object.Object)
	dateObj, _ := params[1].(*object.Object)
	if obj == nil || dateObj == nil {
		return types.JavaBoolFalse
	}
	millisVal := udateGetTime([]interface{}{dateObj})
	millis, ok := millisVal.(int64)
	if !ok {
		return types.JavaBoolFalse
	}
	id, err := tzGetStringID(obj)
	if err != nil {
		return types.JavaBoolFalse
	}
	offsetAt := tzCurrentOffsetAt(id, millis)
	raw := tzGetRawOffset([]interface{}{obj}).(int64)
	if offsetAt != raw {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func tzObservesDaylightTime(params []interface{}) interface{} {
	obj, _ := params[0].(*object.Object)
	if obj == nil {
		return types.JavaBoolFalse
	}
	dst := ttzGetDSTSavings([]interface{}{obj}).(int64)
	if dst != 0 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

var defaultTZ *object.Object = nil

func tzSetDefault(params []interface{}) interface{} {
	// Set static default; ignore invalid
	if len(params) < 1 {
		return nil
	}
	tz, _ := params[0].(*object.Object)
	defaultTZ = tz
	return nil
}

func tzSetID(params []interface{}) interface{} {
	obj, _ := params[0].(*object.Object)
	idObj, _ := params[1].(*object.Object)
	if obj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "tzSetID: self not object")
	}
	if idObj == nil { // set to empty
		obj.FieldTable[fieldTZID] = object.Field{Ftype: types.StringClassRef, Fvalue: nil}
		obj.FieldTable[fieldRawOffset] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
		obj.FieldTable[fieldDSTSavings] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
		return nil
	}
	id := object.GoStringFromStringObject(idObj)
	raw, dst := tzComputeOffsetsForID(id)
	obj.FieldTable[fieldTZID] = object.Field{Ftype: types.StringClassRef, Fvalue: object.StringObjectFromGoString(id)}
	obj.FieldTable[fieldRawOffset] = object.Field{Ftype: types.Int, Fvalue: raw}
	obj.FieldTable[fieldDSTSavings] = object.Field{Ftype: types.Int, Fvalue: dst}
	return nil
}

func tzSetRawOffset(params []interface{}) interface{} {
	obj, _ := params[0].(*object.Object)
	if obj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "tzSetRawOffset: self not object")
	}
	val, _ := params[1].(int64)
	obj.FieldTable[fieldRawOffset] = object.Field{Ftype: types.Int, Fvalue: val}
	return nil
}
