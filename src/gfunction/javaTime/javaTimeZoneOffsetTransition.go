/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"sync"
)

var (
	classNameZoneOffsetTransition = "java/time/zone/ZoneOffsetTransition"
	zoneOffsetTransitionMutex     sync.Mutex
	zoneOffsetTransitionInited    bool = false
)

func ensureZoneOffsetTransitionInited() {
	zoneOffsetTransitionMutex.Lock()
	defer zoneOffsetTransitionMutex.Unlock()
	if zoneOffsetTransitionInited {
		return
	}
	ensureZoneOffsetInited()
	zoneOffsetTransitionInited = true
}

// Load_Time_Zone_ZoneOffsetTransition registers all java/time/zone/ZoneOffsetTransition method signatures.
func Load_Time_Zone_ZoneOffsetTransition() {
	ensureZoneOffsetTransitionInited()

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.<init>(Ljava/time/LocalDateTime;Ljava/time/ZoneOffset;Ljava/time/ZoneOffset;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  zoneOffsetTransitionInitLdt,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.<init>(JLjava/time/ZoneOffset;Ljava/time/ZoneOffset;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  zoneOffsetTransitionInitEpoch,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.of(Ljava/time/LocalDateTime;Ljava/time/ZoneOffset;Ljava/time/ZoneOffset;)Ljava/time/zone/ZoneOffsetTransition;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  zoneOffsetTransitionOf,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.getInstant()Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetTransitionGetInstant,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.toEpochSecond()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetTransitionToEpochSecond,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.getDateTimeBefore()Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetTransitionGetDateTimeBefore,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.getDateTimeAfter()Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetTransitionGetDateTimeAfter,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.getOffsetBefore()Ljava/time/ZoneOffset;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetTransitionGetOffsetBefore,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.getOffsetAfter()Ljava/time/ZoneOffset;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetTransitionGetOffsetAfter,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.getDuration()Ljava/time/Duration;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetTransitionGetDuration,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.isGap()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetTransitionIsGap,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.isOverlap()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetTransitionIsOverlap,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.isValidOffset(Ljava/time/ZoneOffset;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetTransitionIsValidOffset,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.compareTo(Ljava/time/zone/ZoneOffsetTransition;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetTransitionCompareTo,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetTransitionCompareTo,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetTransitionEquals,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetTransitionHashCode,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetTransitionToString,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.writeExternal(Ljava/io/DataOutput;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetTransitionWriteExternal,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.readExternal(Ljava/io/DataInput;)Ljava/time/zone/ZoneOffsetTransition;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetTransitionReadExternal,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneOffsetTransition.write(Ljava/io/DataOutput;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetTransitionWriteExternal,
		}
}

func createZoneOffsetTransition(transition, offsetBefore, offsetAfter *object.Object) *object.Object {
	epochSecVal := localDateTimeToEpochSecond([]any{transition, offsetBefore})
	var epochSec int64
	if s, ok := epochSecVal.(int64); ok {
		epochSec = s
	}

	obj := object.MakeEmptyObjectWithClassName(&classNameZoneOffsetTransition)
	obj.FieldTable["epochSecond"] = object.Field{Ftype: types.Long, Fvalue: epochSec}
	obj.FieldTable["transition"] = object.Field{Ftype: types.Ref, Fvalue: transition}
	obj.FieldTable["offsetBefore"] = object.Field{Ftype: types.Ref, Fvalue: offsetBefore}
	obj.FieldTable["offsetAfter"] = object.Field{Ftype: types.Ref, Fvalue: offsetAfter}
	return obj
}

func createZoneOffsetTransitionFromEpoch(epochSec int64, offsetBefore, offsetAfter *object.Object) *object.Object {
	secBefore, _ := getZoneOffsetTotalSecondsVal(offsetBefore)
	transLdt := localDateTimeOfEpochSecond([]any{epochSec + secBefore, int64(0), offsetBefore})
	var ldtObj *object.Object
	if l, ok := transLdt.(*object.Object); ok {
		ldtObj = l
	}

	obj := object.MakeEmptyObjectWithClassName(&classNameZoneOffsetTransition)
	obj.FieldTable["epochSecond"] = object.Field{Ftype: types.Long, Fvalue: epochSec}
	obj.FieldTable["transition"] = object.Field{Ftype: types.Ref, Fvalue: ldtObj}
	obj.FieldTable["offsetBefore"] = object.Field{Ftype: types.Ref, Fvalue: offsetBefore}
	obj.FieldTable["offsetAfter"] = object.Field{Ftype: types.Ref, Fvalue: offsetAfter}
	return obj
}

func getTransitionEpochSecondVal(obj *object.Object) (int64, bool) {
	if obj == nil || obj.FieldTable == nil {
		return 0, false
	}
	if fld, ok := obj.FieldTable["epochSecond"]; ok {
		if val, okVal := getInt64Val(fld.Fvalue); okVal {
			return val, true
		}
	}
	return 0, false
}

func getTransitionEpochSecond(obj *object.Object) int64 {
	val, _ := getTransitionEpochSecondVal(obj)
	return val
}

func getTransitionLocalDateTime(obj *object.Object) *object.Object {
	if obj == nil || obj.FieldTable == nil {
		return nil
	}
	if fld, ok := obj.FieldTable["transition"]; ok {
		if ldt, okLdt := fld.Fvalue.(*object.Object); okLdt {
			return ldt
		}
	}
	return nil
}

func getTransitionOffsetBefore(obj *object.Object) *object.Object {
	if obj == nil || obj.FieldTable == nil {
		return nil
	}
	if fld, ok := obj.FieldTable["offsetBefore"]; ok {
		if off, okOff := fld.Fvalue.(*object.Object); okOff {
			return off
		}
	}
	return nil
}

func getTransitionOffsetAfter(obj *object.Object) *object.Object {
	if obj == nil || obj.FieldTable == nil {
		return nil
	}
	if fld, ok := obj.FieldTable["offsetAfter"]; ok {
		if off, okOff := fld.Fvalue.(*object.Object); okOff {
			return off
		}
	}
	return nil
}

func zoneOffsetTransitionInitLdt(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.<init>: self is null")
	}
	if len(params) < 4 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneOffsetTransition.<init>: invalid argument count")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.<init>: transition is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.<init>: offsetBefore is null")
	}
	if params[3] == nil || object.IsNull(params[3].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.<init>: offsetAfter is null")
	}

	self := params[0].(*object.Object)
	transition := params[1].(*object.Object)
	offsetBefore := params[2].(*object.Object)
	offsetAfter := params[3].(*object.Object)

	epochSecVal := localDateTimeToEpochSecond([]any{transition, offsetBefore})
	var epochSec int64
	if s, ok := epochSecVal.(int64); ok {
		epochSec = s
	}

	self.FieldTable["epochSecond"] = object.Field{Ftype: types.Long, Fvalue: epochSec}
	self.FieldTable["transition"] = object.Field{Ftype: types.Ref, Fvalue: transition}
	self.FieldTable["offsetBefore"] = object.Field{Ftype: types.Ref, Fvalue: offsetBefore}
	self.FieldTable["offsetAfter"] = object.Field{Ftype: types.Ref, Fvalue: offsetAfter}
	return nil
}

func zoneOffsetTransitionInitEpoch(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.<init>: self is null")
	}
	if len(params) < 4 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneOffsetTransition.<init>: invalid argument count")
	}
	epochSec, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneOffsetTransition.<init>: invalid epochSecond")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.<init>: offsetBefore is null")
	}
	if params[3] == nil || object.IsNull(params[3].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.<init>: offsetAfter is null")
	}

	self := params[0].(*object.Object)
	offsetBefore := params[2].(*object.Object)
	offsetAfter := params[3].(*object.Object)

	secBefore, _ := getZoneOffsetTotalSecondsVal(offsetBefore)
	transLdt := localDateTimeOfEpochSecond([]any{epochSec + secBefore, int64(0), offsetBefore})
	var ldtObj *object.Object
	if l, ok := transLdt.(*object.Object); ok {
		ldtObj = l
	}

	self.FieldTable["epochSecond"] = object.Field{Ftype: types.Long, Fvalue: epochSec}
	self.FieldTable["transition"] = object.Field{Ftype: types.Ref, Fvalue: ldtObj}
	self.FieldTable["offsetBefore"] = object.Field{Ftype: types.Ref, Fvalue: offsetBefore}
	self.FieldTable["offsetAfter"] = object.Field{Ftype: types.Ref, Fvalue: offsetAfter}
	return nil
}

func zoneOffsetTransitionOf(params []any) any {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneOffsetTransition.of: invalid argument count")
	}
	if params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.of: transition is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.of: offsetBefore is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.of: offsetAfter is null")
	}

	transition := params[0].(*object.Object)
	offsetBefore := params[1].(*object.Object)
	offsetAfter := params[2].(*object.Object)

	eq := zoneOffsetEquals([]any{offsetBefore, offsetAfter})
	if eq == types.JavaBoolTrue {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneOffsetTransition.of: Offsets must not be equal")
	}

	timeObj := getLocalDateTimeTime(transition)
	if timeObj != nil {
		_, _, _, nano := getLocalTimeHMSN(timeObj)
		if nano != 0 {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneOffsetTransition.of: Nano-of-second must be zero")
		}
	}

	return createZoneOffsetTransition(transition, offsetBefore, offsetAfter)
}

func zoneOffsetTransitionGetInstant(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.getInstant: self is null")
	}
	self := params[0].(*object.Object)
	epochSec := getTransitionEpochSecond(self)
	return createInstant(epochSec, 0)
}

func zoneOffsetTransitionToEpochSecond(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.toEpochSecond: self is null")
	}
	self := params[0].(*object.Object)
	return getTransitionEpochSecond(self)
}

func zoneOffsetTransitionGetDateTimeBefore(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.getDateTimeBefore: self is null")
	}
	self := params[0].(*object.Object)
	return getTransitionLocalDateTime(self)
}

func zoneOffsetTransitionGetDateTimeAfter(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.getDateTimeAfter: self is null")
	}
	self := params[0].(*object.Object)
	trans := getTransitionLocalDateTime(self)
	offBefore := getTransitionOffsetBefore(self)
	offAfter := getTransitionOffsetAfter(self)

	secBefore, _ := getZoneOffsetTotalSecondsVal(offBefore)
	secAfter, _ := getZoneOffsetTotalSecondsVal(offAfter)
	diffSec := secAfter - secBefore

	return localDateTimePlusSeconds([]any{trans, diffSec})
}

func zoneOffsetTransitionGetOffsetBefore(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.getOffsetBefore: self is null")
	}
	self := params[0].(*object.Object)
	return getTransitionOffsetBefore(self)
}

func zoneOffsetTransitionGetOffsetAfter(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.getOffsetAfter: self is null")
	}
	self := params[0].(*object.Object)
	return getTransitionOffsetAfter(self)
}

func zoneOffsetTransitionGetDuration(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.getDuration: self is null")
	}
	self := params[0].(*object.Object)
	offBefore := getTransitionOffsetBefore(self)
	offAfter := getTransitionOffsetAfter(self)

	secBefore, _ := getZoneOffsetTotalSecondsVal(offBefore)
	secAfter, _ := getZoneOffsetTotalSecondsVal(offAfter)
	diffSec := secAfter - secBefore

	return createDuration(diffSec, 0)
}

func zoneOffsetTransitionIsGap(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.isGap: self is null")
	}
	self := params[0].(*object.Object)
	offBefore := getTransitionOffsetBefore(self)
	offAfter := getTransitionOffsetAfter(self)

	secBefore, _ := getZoneOffsetTotalSecondsVal(offBefore)
	secAfter, _ := getZoneOffsetTotalSecondsVal(offAfter)

	if secAfter > secBefore {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func zoneOffsetTransitionIsOverlap(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.isOverlap: self is null")
	}
	self := params[0].(*object.Object)
	offBefore := getTransitionOffsetBefore(self)
	offAfter := getTransitionOffsetAfter(self)

	secBefore, _ := getZoneOffsetTotalSecondsVal(offBefore)
	secAfter, _ := getZoneOffsetTotalSecondsVal(offAfter)

	if secAfter < secBefore {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func zoneOffsetTransitionIsValidOffset(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.isValidOffset: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	self := params[0].(*object.Object)
	targetOffset := params[1].(*object.Object)

	offBefore := getTransitionOffsetBefore(self)
	offAfter := getTransitionOffsetAfter(self)

	secBefore, _ := getZoneOffsetTotalSecondsVal(offBefore)
	secAfter, _ := getZoneOffsetTotalSecondsVal(offAfter)

	if secAfter > secBefore {
		// Gap has no valid offsets
		return types.JavaBoolFalse
	}

	// Overlap has before and after as valid offsets
	if zoneOffsetEquals([]any{offBefore, targetOffset}) == types.JavaBoolTrue ||
		zoneOffsetEquals([]any{offAfter, targetOffset}) == types.JavaBoolTrue {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func zoneOffsetTransitionCompareTo(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.compareTo: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.compareTo: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	sec1 := getTransitionEpochSecond(self)
	sec2 := getTransitionEpochSecond(other)

	if sec1 < sec2 {
		return int64(-1)
	} else if sec1 > sec2 {
		return int64(1)
	}
	return int64(0)
}

func zoneOffsetTransitionEquals(params []any) any {
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

	sec1, ok1 := getTransitionEpochSecondVal(self)
	sec2, ok2 := getTransitionEpochSecondVal(other)
	if !ok1 || !ok2 || sec1 != sec2 {
		return types.JavaBoolFalse
	}

	b1 := getTransitionOffsetBefore(self)
	b2 := getTransitionOffsetBefore(other)
	if b1 == nil || b2 == nil || zoneOffsetEquals([]any{b1, b2}) != types.JavaBoolTrue {
		return types.JavaBoolFalse
	}

	a1 := getTransitionOffsetAfter(self)
	a2 := getTransitionOffsetAfter(other)
	if a1 == nil || a2 == nil || zoneOffsetEquals([]any{a1, a2}) != types.JavaBoolTrue {
		return types.JavaBoolFalse
	}

	return types.JavaBoolTrue
}

func zoneOffsetTransitionHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.hashCode: self is null")
	}
	self := params[0].(*object.Object)
	trans := getTransitionLocalDateTime(self)
	offBefore := getTransitionOffsetBefore(self)
	offAfter := getTransitionOffsetAfter(self)

	var h1, h2, h3 int64
	if res := localDateTimeHashCode([]any{trans}); res != nil {
		if val, ok := res.(int64); ok {
			h1 = val
		}
	}
	if res := zoneOffsetHashCode([]any{offBefore}); res != nil {
		if val, ok := res.(int64); ok {
			h2 = val
		}
	}
	if res := zoneOffsetHashCode([]any{offAfter}); res != nil {
		if val, ok := res.(int64); ok {
			h3 = val
		}
	}

	rotH3 := int64(uint32(h3)<<16 | uint32(h3)>>16)
	return int64(int32(h1 ^ h2 ^ rotH3))
}

func zoneOffsetTransitionToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffsetTransition.toString: self is null")
	}
	self := params[0].(*object.Object)
	trans := getTransitionLocalDateTime(self)
	offBefore := getTransitionOffsetBefore(self)
	offAfter := getTransitionOffsetAfter(self)

	secBefore, _ := getZoneOffsetTotalSecondsVal(offBefore)
	secAfter, _ := getZoneOffsetTotalSecondsVal(offAfter)

	transType := "Gap"
	if secAfter < secBefore {
		transType = "Overlap"
	}

	transStr := ""
	if strRes := localDateTimeToString([]any{trans}); strRes != nil {
		if sObj, ok := strRes.(*object.Object); ok {
			transStr = object.GoStringFromStringObject(sObj)
		}
	}

	beforeStr := ""
	if strRes := zoneOffsetToString([]any{offBefore}); strRes != nil {
		if sObj, ok := strRes.(*object.Object); ok {
			beforeStr = object.GoStringFromStringObject(sObj)
		}
	}

	afterStr := ""
	if strRes := zoneOffsetToString([]any{offAfter}); strRes != nil {
		if sObj, ok := strRes.(*object.Object); ok {
			afterStr = object.GoStringFromStringObject(sObj)
		}
	}

	res := fmt.Sprintf("Transition[%s at %s%s to %s]", transType, transStr, beforeStr, afterStr)
	return object.StringObjectFromGoString(res)
}

func zoneOffsetTransitionWriteExternal([]any) any {
	return nil
}

func zoneOffsetTransitionReadExternal([]any) any {
	return nil
}
