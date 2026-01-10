/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"time"
)

func Load_Util_Date() {

	ghelpers.MethodSignatures["java/util/Date.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/Date.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  udateInit,
		}

	ghelpers.MethodSignatures["java/util/Date.<init>(III)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  udateInit3Ints,
		}

	ghelpers.MethodSignatures["java/util/Date.<init>(IIIII)V"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  udateInit5Ints,
		}

	ghelpers.MethodSignatures["java/util/Date.<init>(IIIIII)V"] =
		ghelpers.GMeth{
			ParamSlots: 6,
			GFunction:  udateInit6Ints,
		}

	ghelpers.MethodSignatures["java/util/Date.<init>(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  udateInitLong,
		}

	ghelpers.MethodSignatures["java/util/Date.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  udateInitString,
		}

	ghelpers.MethodSignatures["java/util/Date.after(Ljava/util/Date;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  udateAfter,
		}

	ghelpers.MethodSignatures["java/util/Date.before(Ljava/util/Date;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  udateBefore,
		}

	ghelpers.MethodSignatures["java/util/Date.clone()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  udateClone,
		}

	ghelpers.MethodSignatures["java/util/Date.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  udateEquals,
		}

	ghelpers.MethodSignatures["java/util/Date.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  udateHashCode,
		}

	// --- Deprecated getters ---
	ghelpers.MethodSignatures["java/util/Date.getDate()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapDeprecated,
		}
	ghelpers.MethodSignatures["java/util/Date.getDay()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapDeprecated,
		}
	ghelpers.MethodSignatures["java/util/Date.getHours()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapDeprecated,
		}
	ghelpers.MethodSignatures["java/util/Date.getMinutes()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapDeprecated,
		}
	ghelpers.MethodSignatures["java/util/Date.getMonth()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapDeprecated,
		}
	ghelpers.MethodSignatures["java/util/Date.getSeconds()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures["java/util/Date.getTime()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  udateGetTime,
		}

	// --- Deprecated setters ---
	ghelpers.MethodSignatures["java/util/Date.setDate(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapDeprecated,
		}
	ghelpers.MethodSignatures["java/util/Date.setHours(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapDeprecated,
		}
	ghelpers.MethodSignatures["java/util/Date.setMinutes(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapDeprecated,
		}
	ghelpers.MethodSignatures["java/util/Date.setMonth(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapDeprecated,
		}
	ghelpers.MethodSignatures["java/util/Date.setSeconds(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures["java/util/Date.setTime(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  udateSetTime,
		}

	ghelpers.MethodSignatures["java/util/Date.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  udateToString,
		}
}

// === java/util/Date minimal implementation ===

const dateValueField = "value" // store milliseconds since epoch (UTC) as types.Long

// udateInit: no-arg constructor -> initialize to current time in milliseconds.
func udateInit(params []interface{}) interface{} {
	if len(params) < 1 {
		return nil
	}
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "udateInit: self is not an object")
	}
	object.ClearFieldTable(obj)
	millis := time.Now().UnixMilli()
	obj.FieldTable[dateValueField] = object.Field{Ftype: types.Long, Fvalue: millis}
	return nil
}

// Deprecated ctors: return a ghelpers.Trap indicating deprecation (unsupported).
func udateInit3Ints(params []interface{}) interface{}  { return ghelpers.TrapDeprecated(params) }
func udateInit5Ints(params []interface{}) interface{}  { return ghelpers.TrapDeprecated(params) }
func udateInit6Ints(params []interface{}) interface{}  { return ghelpers.TrapDeprecated(params) }
func udateInitString(params []interface{}) interface{} { return ghelpers.TrapDeprecated(params) }

// Construct from long milliseconds.
func udateInitLong(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "udateInitLong: missing millis parameter")
	}
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "udateInitLong: self is not an object")
	}
	millis, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "udateInitLong: millis parameter is not a long")
	}
	object.ClearFieldTable(obj)
	obj.FieldTable[dateValueField] = object.Field{Ftype: types.Long, Fvalue: millis}
	return nil
}

// Helper to get millis and maybe error.
func DateGetMillis(obj *object.Object) (int64, interface{}) {
	if obj == nil {
		return 0, ghelpers.GetGErrBlk(excNames.NullPointerException, "Date object is null")
	}
	fld, ok := obj.FieldTable[dateValueField]
	if !ok {
		return 0, ghelpers.GetGErrBlk(excNames.IllegalStateException, "Date value field missing")
	}
	millis, ok := fld.Fvalue.(int64)
	if !ok {
		return 0, ghelpers.GetGErrBlk(excNames.IllegalStateException, "Date value field not a long")
	}
	return millis, nil
}

func udateAfter(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "udateAfter: requires this and other Date")
	}
	this, _ := params[0].(*object.Object)
	other, _ := params[1].(*object.Object)
	m1, err := DateGetMillis(this)
	if err != nil {
		return err
	}
	m2, err2 := DateGetMillis(other)
	if err2 != nil {
		return err2
	}
	if m1 > m2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func udateBefore(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "udateBefore: requires this and other Date")
	}
	this, _ := params[0].(*object.Object)
	other, _ := params[1].(*object.Object)
	m1, err := DateGetMillis(this)
	if err != nil {
		return err
	}
	m2, err2 := DateGetMillis(other)
	if err2 != nil {
		return err2
	}
	if m1 < m2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func udateClone(params []interface{}) interface{} {
	if len(params) < 1 {
		return object.Null
	}
	obj, _ := params[0].(*object.Object)
	if obj == nil {
		return object.Null
	}
	return object.CloneObject(obj)
}

func udateEquals(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "udateEquals: requires this and other Object")
	}
	this, _ := params[0].(*object.Object)
	other, _ := params[1].(*object.Object)
	// If other is null, return false
	if object.IsNull(other) {
		return types.JavaBoolFalse
	}
	m1, err := DateGetMillis(this)
	if err != nil {
		return err
	}
	m2, err2 := DateGetMillis(other)
	if err2 != nil {
		return err2
	}
	if m1 == m2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func udateHashCode(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "udateHashCode: requires this")
	}
	this, _ := params[0].(*object.Object)
	m, err := DateGetMillis(this)
	if err != nil {
		return err
	}
	// Java Date.hashCode: (int)(value ^ (value >>> 32))
	h := int32(m ^ (m >> 32))
	return int64(h)
}

func udateGetTime(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "udateGetTime: requires this")
	}
	this, _ := params[0].(*object.Object)
	m, err := DateGetMillis(this)
	if err != nil {
		return err
	}
	return m
}

func udateSetTime(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "udateSetTime: requires this and millis")
	}
	this, _ := params[0].(*object.Object)
	if this == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "udateSetTime: this is null")
	}
	m, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "udateSetTime: millis is not a long")
	}
	fld := object.Field{Ftype: types.Long, Fvalue: m}
	this.FieldTable[dateValueField] = fld
	return nil
}

func udateToString(params []interface{}) interface{} {
	if len(params) < 1 {
		return object.StringObjectFromGoString("null")
	}
	this, _ := params[0].(*object.Object)
	m, err := DateGetMillis(this)
	if err != nil {
		// On error, return an informative string
		msg := fmt.Sprintf("Date[error:%T]", err)
		return object.StringObjectFromGoString(msg)
	}
	// Minimal readable format: RFC3339-like in local time for friendliness.
	t := time.UnixMilli(m).Local()
	// Use Go's default formatting; Java's exact format is not replicated.
	return object.StringObjectFromGoString(t.String())
}
