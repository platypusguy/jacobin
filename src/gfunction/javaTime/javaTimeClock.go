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
	"jacobin/src/types"
	"math"
	"time"
)

var (
	classNameClockVal       = "java/time/Clock"
	classNameSystemClockVal  = "java/time/Clock$SystemClock"
	classNameFixedClockVal   = "java/time/Clock$FixedClock"
	classNameOffsetClockVal  = "java/time/Clock$OffsetClock"
	classNameTickClockVal    = "java/time/Clock$TickClock"
	clockInited             bool
)

func ensureClockInited() {
	if clockInited {
		return
	}
	ensureZoneIdInited()
	ensureZoneOffsetInited()
	ensureInstantInited()
	clockInited = true
}

// Load_Time_Clock registers method signatures for java/time/Clock and its inner classes.
func Load_Time_Clock() {
	ensureClockInited()

	// java/time/Clock
	ghelpers.MethodSignatures["java/time/Clock.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  clockClinit,
		}

	ghelpers.MethodSignatures["java/time/Clock.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  clockInit,
		}

	ghelpers.MethodSignatures["java/time/Clock.systemUTC()Ljava/time/Clock;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  clockSystemUTC,
		}

	ghelpers.MethodSignatures["java/time/Clock.systemDefaultZone()Ljava/time/Clock;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  clockSystemDefaultZone,
		}

	ghelpers.MethodSignatures["java/time/Clock.system(Ljava/time/ZoneId;)Ljava/time/Clock;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  clockSystem,
		}

	ghelpers.MethodSignatures["java/time/Clock.tickMillis(Ljava/time/ZoneId;)Ljava/time/Clock;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  clockTickMillis,
		}

	ghelpers.MethodSignatures["java/time/Clock.tickSeconds(Ljava/time/ZoneId;)Ljava/time/Clock;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  clockTickSeconds,
		}

	ghelpers.MethodSignatures["java/time/Clock.tickMinutes(Ljava/time/ZoneId;)Ljava/time/Clock;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  clockTickMinutes,
		}

	ghelpers.MethodSignatures["java/time/Clock.tick(Ljava/time/Clock;Ljava/time/Duration;)Ljava/time/Clock;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  clockTick,
		}

	ghelpers.MethodSignatures["java/time/Clock.fixed(Ljava/time/Instant;Ljava/time/ZoneId;)Ljava/time/Clock;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  clockFixed,
		}

	ghelpers.MethodSignatures["java/time/Clock.offset(Ljava/time/Clock;Ljava/time/Duration;)Ljava/time/Clock;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  clockOffset,
		}

	ghelpers.MethodSignatures["java/time/Clock.getZone()Ljava/time/ZoneId;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  clockGetZone,
		}

	ghelpers.MethodSignatures["java/time/Clock.withZone(Ljava/time/ZoneId;)Ljava/time/Clock;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  clockWithZone,
		}

	ghelpers.MethodSignatures["java/time/Clock.millis()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  clockMillis,
		}

	ghelpers.MethodSignatures["java/time/Clock.instant()Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  clockInstant,
		}

	ghelpers.MethodSignatures["java/time/Clock.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  clockEquals,
		}

	ghelpers.MethodSignatures["java/time/Clock.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  clockHashCode,
		}

	ghelpers.MethodSignatures["java/time/Clock.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  clockToString,
		}

	// Register signatures for inner classes
	registerInnerClockSignatures(classNameSystemClockVal)
	registerInnerClockSignatures(classNameFixedClockVal)
	registerInnerClockSignatures(classNameOffsetClockVal)
	registerInnerClockSignatures(classNameTickClockVal)
}

func registerInnerClockSignatures(kName string) {
	ghelpers.MethodSignatures[kName+".<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  clockClinit,
		}

	ghelpers.MethodSignatures[kName+".<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  clockInit,
		}

	ghelpers.MethodSignatures[kName+".getZone()Ljava/time/ZoneId;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  clockGetZone,
		}

	ghelpers.MethodSignatures[kName+".withZone(Ljava/time/ZoneId;)Ljava/time/Clock;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  clockWithZone,
		}

	ghelpers.MethodSignatures[kName+".millis()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  clockMillis,
		}

	ghelpers.MethodSignatures[kName+".instant()Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  clockInstant,
		}

	ghelpers.MethodSignatures[kName+".equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  clockEquals,
		}

	ghelpers.MethodSignatures[kName+".hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  clockHashCode,
		}

	ghelpers.MethodSignatures[kName+".toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  clockToString,
		}
}

// Clock representation:
// clockType: 0 = SystemClock, 1 = FixedClock, 2 = OffsetClock, 3 = TickClock
const (
	clockTypeSystem = 0
	clockTypeFixed  = 1
	clockTypeOffset = 2
	clockTypeTick   = 3
)

func createSystemClock(zoneObj *object.Object) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameSystemClockVal)
	obj.FieldTable["zone"] = object.Field{Ftype: types.Ref, Fvalue: zoneObj}
	obj.FieldTable["clockType"] = object.Field{Ftype: types.Int, Fvalue: int64(clockTypeSystem)}
	return obj
}

func createFixedClock(instantObj *object.Object, zoneObj *object.Object) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameFixedClockVal)
	obj.FieldTable["instant"] = object.Field{Ftype: types.Ref, Fvalue: instantObj}
	obj.FieldTable["zone"] = object.Field{Ftype: types.Ref, Fvalue: zoneObj}
	obj.FieldTable["clockType"] = object.Field{Ftype: types.Int, Fvalue: int64(clockTypeFixed)}
	return obj
}

func createOffsetClock(baseClock *object.Object, offsetDuration *object.Object) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameOffsetClockVal)
	obj.FieldTable["baseClock"] = object.Field{Ftype: types.Ref, Fvalue: baseClock}
	obj.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: offsetDuration}
	obj.FieldTable["clockType"] = object.Field{Ftype: types.Int, Fvalue: int64(clockTypeOffset)}
	return obj
}

func createTickClock(baseClock *object.Object, tickNanos int64) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameTickClockVal)
	obj.FieldTable["baseClock"] = object.Field{Ftype: types.Ref, Fvalue: baseClock}
	obj.FieldTable["tickNanos"] = object.Field{Ftype: types.Long, Fvalue: tickNanos}
	obj.FieldTable["clockType"] = object.Field{Ftype: types.Int, Fvalue: int64(clockTypeTick)}
	return obj
}

func getClockType(obj *object.Object) int {
	if obj == nil || obj.FieldTable == nil {
		return clockTypeSystem
	}
	if fld, ok := obj.FieldTable["clockType"]; ok {
		if val, ok2 := getInt64Val(fld.Fvalue); ok2 {
			return int(val)
		}
	}
	kName := object.GoStringFromStringPoolIndex(obj.KlassName)
	switch kName {
	case classNameFixedClockVal:
		return clockTypeFixed
	case classNameOffsetClockVal:
		return clockTypeOffset
	case classNameTickClockVal:
		return clockTypeTick
	default:
		return clockTypeSystem
	}
}

func clockClinit([]any) any {
	ensureClockInited()
	return nil
}

func clockInit([]any) any {
	return nil
}

func clockSystemUTC([]any) any {
	ensureClockInited()
	zoneUTC := makeZoneIdFromID("UTC")
	return createSystemClock(zoneUTC)
}

func clockSystemDefaultZone([]any) any {
	ensureClockInited()
	zoneDef := zoneIdSystemDefault(nil)
	if zObj, ok := zoneDef.(*object.Object); ok {
		return createSystemClock(zObj)
	}
	return createSystemClock(makeZoneIdFromID("UTC"))
}

func clockSystem(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Clock.system: zone is null")
	}
	zoneObj := params[0].(*object.Object)
	return createSystemClock(zoneObj)
}

func clockTickMillis(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Clock.tickMillis: zone is null")
	}
	zoneObj := params[0].(*object.Object)
	return createSystemClock(zoneObj)
}

func clockTickSeconds(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Clock.tickSeconds: zone is null")
	}
	zoneObj := params[0].(*object.Object)
	base := createSystemClock(zoneObj)
	return createTickClock(base, 1000000000)
}

func clockTickMinutes(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Clock.tickMinutes: zone is null")
	}
	zoneObj := params[0].(*object.Object)
	base := createSystemClock(zoneObj)
	return createTickClock(base, 60000000000)
}

func clockTick(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Clock.tick: null argument")
	}
	baseClock := params[0].(*object.Object)
	tickDuration := params[1].(*object.Object)

	sec := getDurationSecondsVal(tickDuration)
	nano := getDurationNanosVal(tickDuration)

	if sec < 0 || (sec == 0 && nano < 0) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Tick duration must not be negative")
	}
	if sec > math.MaxInt64/1000000000 {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "long overflow")
	}

	tickNanos := sec*1000000000 + int64(nano)
	if tickNanos <= 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Tick duration must not be negative or zero")
	}
	if tickNanos%1000000 != 0 && 1000000000%tickNanos != 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Invalid tick duration")
	}
	if tickNanos <= 1 {
		return baseClock
	}
	return createTickClock(baseClock, tickNanos)
}

func clockFixed(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Clock.fixed: null argument")
	}
	instant := params[0].(*object.Object)
	zone := params[1].(*object.Object)
	return createFixedClock(instant, zone)
}

func clockOffset(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Clock.offset: null argument")
	}
	baseClock := params[0].(*object.Object)
	offsetDuration := params[1].(*object.Object)

	sec := getDurationSecondsVal(offsetDuration)
	nano := getDurationNanosVal(offsetDuration)
	if sec == 0 && nano == 0 {
		return baseClock
	}
	return createOffsetClock(baseClock, offsetDuration)
}

func clockGetZone(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Clock.getZone: self is null")
	}
	self := params[0].(*object.Object)
	cType := getClockType(self)
	switch cType {
	case clockTypeSystem, clockTypeFixed:
		if zFld, ok := self.FieldTable["zone"]; ok && zFld.Fvalue != nil {
			return zFld.Fvalue.(*object.Object)
		}
		return makeZoneIdFromID("UTC")
	case clockTypeOffset, clockTypeTick:
		if bFld, ok := self.FieldTable["baseClock"]; ok && bFld.Fvalue != nil {
			return clockGetZone([]any{bFld.Fvalue.(*object.Object)})
		}
		return makeZoneIdFromID("UTC")
	default:
		if zFld, ok := self.FieldTable["zone"]; ok && zFld.Fvalue != nil {
			return zFld.Fvalue.(*object.Object)
		}
		return makeZoneIdFromID("UTC")
	}
}

func clockWithZone(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Clock.withZone: null argument")
	}
	self := params[0].(*object.Object)
	zone := params[1].(*object.Object)

	currZone := clockGetZone([]any{self})
	if currZoneObj, ok := currZone.(*object.Object); ok && zoneEquals(currZoneObj, zone) {
		return self
	}

	cType := getClockType(self)
	switch cType {
	case clockTypeSystem:
		return createSystemClock(zone)
	case clockTypeFixed:
		inst := self.FieldTable["instant"].Fvalue.(*object.Object)
		return createFixedClock(inst, zone)
	case clockTypeOffset:
		base := self.FieldTable["baseClock"].Fvalue.(*object.Object)
		offset := self.FieldTable["offset"].Fvalue.(*object.Object)
		newBaseRes := clockWithZone([]any{base, zone})
		if err, isErr := newBaseRes.(*ghelpers.GErrBlk); isErr {
			return err
		}
		return createOffsetClock(newBaseRes.(*object.Object), offset)
	case clockTypeTick:
		base := self.FieldTable["baseClock"].Fvalue.(*object.Object)
		tickNanos := self.FieldTable["tickNanos"].Fvalue.(int64)
		newBaseRes := clockWithZone([]any{base, zone})
		if err, isErr := newBaseRes.(*ghelpers.GErrBlk); isErr {
			return err
		}
		return createTickClock(newBaseRes.(*object.Object), tickNanos)
	default:
		return createSystemClock(zone)
	}
}

func clockMillis(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Clock.millis: self is null")
	}
	self := params[0].(*object.Object)
	cType := getClockType(self)
	switch cType {
	case clockTypeSystem:
		return time.Now().UnixMilli()
	case clockTypeFixed:
		inst := self.FieldTable["instant"].Fvalue.(*object.Object)
		sec := getInstantEpochSecondVal(inst)
		nano := getInstantNanoVal(inst)
		return sec*1000 + nano/1000000
	case clockTypeOffset:
		base := self.FieldTable["baseClock"].Fvalue.(*object.Object)
		baseMillisRes := clockMillis([]any{base})
		if err, isErr := baseMillisRes.(*ghelpers.GErrBlk); isErr {
			return err
		}
		baseMillis := baseMillisRes.(int64)
		offset := self.FieldTable["offset"].Fvalue.(*object.Object)
		offSec := getDurationSecondsVal(offset)
		offNano := getDurationNanosVal(offset)
		offMillis := offSec*1000 + int64(offNano)/1000000
		return baseMillis + offMillis
	case clockTypeTick:
		base := self.FieldTable["baseClock"].Fvalue.(*object.Object)
		tickNanos := self.FieldTable["tickNanos"].Fvalue.(int64)
		if tickNanos%1000000 == 0 {
			baseMillisRes := clockMillis([]any{base})
			if err, isErr := baseMillisRes.(*ghelpers.GErrBlk); isErr {
				return err
			}
			baseMillis := baseMillisRes.(int64)
			tickMillis := tickNanos / 1000000
			return baseMillis - floorMod(baseMillis, tickMillis)
		}
		instRes := clockInstant([]any{self})
		if err, isErr := instRes.(*ghelpers.GErrBlk); isErr {
			return err
		}
		inst := instRes.(*object.Object)
		sec := getInstantEpochSecondVal(inst)
		nano := getInstantNanoVal(inst)
		return sec*1000 + nano/1000000
	default:
		return time.Now().UnixMilli()
	}
}

func clockInstant(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Clock.instant: self is null")
	}
	self := params[0].(*object.Object)
	cType := getClockType(self)
	switch cType {
	case clockTypeSystem:
		now := time.Now().UTC()
		return createInstant(now.Unix(), int64(now.Nanosecond()))
	case clockTypeFixed:
		return self.FieldTable["instant"].Fvalue.(*object.Object)
	case clockTypeOffset:
		base := self.FieldTable["baseClock"].Fvalue.(*object.Object)
		baseInstRes := clockInstant([]any{base})
		if err, isErr := baseInstRes.(*ghelpers.GErrBlk); isErr {
			return err
		}
		baseInst := baseInstRes.(*object.Object)
		offset := self.FieldTable["offset"].Fvalue.(*object.Object)
		sec := getDurationSecondsVal(offset)
		nano := getDurationNanosVal(offset)
		baseSec := getInstantEpochSecondVal(baseInst)
		baseNano := getInstantNanoVal(baseInst)
		totSec := baseSec + sec
		totNano := baseNano + int64(nano)
		return createInstantFromNanos(totSec, totNano)
	case clockTypeTick:
		base := self.FieldTable["baseClock"].Fvalue.(*object.Object)
		tickNanos := self.FieldTable["tickNanos"].Fvalue.(int64)
		if tickNanos%1000000 == 0 {
			baseMillisRes := clockMillis([]any{base})
			if err, isErr := baseMillisRes.(*ghelpers.GErrBlk); isErr {
				return err
			}
			baseMillis := baseMillisRes.(int64)
			tickMillis := tickNanos / 1000000
			adjMillis := baseMillis - floorMod(baseMillis, tickMillis)
			secs := floorDiv(adjMillis, 1000)
			nos := floorMod(adjMillis, 1000) * 1000000
			return createInstantValidated(secs, nos)
		}
		baseInstRes := clockInstant([]any{base})
		if err, isErr := baseInstRes.(*ghelpers.GErrBlk); isErr {
			return err
		}
		baseInst := baseInstRes.(*object.Object)
		nanos := getInstantNanoVal(baseInst)
		adjust := floorMod(nanos, tickNanos)
		sec := getInstantEpochSecondVal(baseInst)
		adjNano := nanos - adjust
		return createInstant(sec, adjNano)
	default:
		now := time.Now().UTC()
		return createInstant(now.Unix(), int64(now.Nanosecond()))
	}
}

func clockEquals(params []any) any {
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

	t1 := getClockType(self)
	t2 := getClockType(other)
	if t1 != t2 {
		return types.JavaBoolFalse
	}

	switch t1 {
	case clockTypeSystem:
		z1 := self.FieldTable["zone"].Fvalue.(*object.Object)
		z2 := other.FieldTable["zone"].Fvalue.(*object.Object)
		return types.ConvertGoBoolToJavaBool(zoneEquals(z1, z2))
	case clockTypeFixed:
		i1 := self.FieldTable["instant"].Fvalue.(*object.Object)
		i2 := other.FieldTable["instant"].Fvalue.(*object.Object)
		z1 := self.FieldTable["zone"].Fvalue.(*object.Object)
		z2 := other.FieldTable["zone"].Fvalue.(*object.Object)
		return types.ConvertGoBoolToJavaBool(instantEqualsObjects(i1, i2) && zoneEquals(z1, z2))
	case clockTypeOffset:
		b1 := self.FieldTable["baseClock"].Fvalue.(*object.Object)
		b2 := other.FieldTable["baseClock"].Fvalue.(*object.Object)
		baseEq := clockEquals([]any{b1, b2}) == types.JavaBoolTrue
		o1 := self.FieldTable["offset"].Fvalue.(*object.Object)
		o2 := other.FieldTable["offset"].Fvalue.(*object.Object)
		offEq := durationEqualsObjects(o1, o2)
		return types.ConvertGoBoolToJavaBool(baseEq && offEq)
	case clockTypeTick:
		b1 := self.FieldTable["baseClock"].Fvalue.(*object.Object)
		b2 := other.FieldTable["baseClock"].Fvalue.(*object.Object)
		baseEq := clockEquals([]any{b1, b2}) == types.JavaBoolTrue
		tickEq := self.FieldTable["tickNanos"].Fvalue.(int64) == other.FieldTable["tickNanos"].Fvalue.(int64)
		return types.ConvertGoBoolToJavaBool(baseEq && tickEq)
	default:
		return types.JavaBoolFalse
	}
}

func clockHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return int64(0)
	}
	self := params[0].(*object.Object)
	t := getClockType(self)
	switch t {
	case clockTypeSystem:
		z := self.FieldTable["zone"].Fvalue.(*object.Object)
		return zoneHashCodeVal(z) + 1
	case clockTypeFixed:
		inst := self.FieldTable["instant"].Fvalue.(*object.Object)
		z := self.FieldTable["zone"].Fvalue.(*object.Object)
		return instantHashCodeVal(inst) ^ zoneHashCodeVal(z)
	case clockTypeOffset:
		b := self.FieldTable["baseClock"].Fvalue.(*object.Object)
		off := self.FieldTable["offset"].Fvalue.(*object.Object)
		bHash := clockHashCode([]any{b}).(int64)
		return bHash ^ durationHashCodeVal(off)
	case clockTypeTick:
		b := self.FieldTable["baseClock"].Fvalue.(*object.Object)
		tn := self.FieldTable["tickNanos"].Fvalue.(int64)
		bHash := clockHashCode([]any{b}).(int64)
		return bHash ^ (tn ^ (tn >> 32))
	default:
		return int64(0)
	}
}

func clockToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return object.StringObjectFromGoString("null")
	}
	self := params[0].(*object.Object)
	t := getClockType(self)
	switch t {
	case clockTypeSystem:
		z := self.FieldTable["zone"].Fvalue.(*object.Object)
		return object.StringObjectFromGoString("SystemClock[" + getZoneString(z) + "]")
	case clockTypeFixed:
		inst := self.FieldTable["instant"].Fvalue.(*object.Object)
		z := self.FieldTable["zone"].Fvalue.(*object.Object)
		return object.StringObjectFromGoString("FixedClock[" + getInstantString(inst) + "," + getZoneString(z) + "]")
	case clockTypeOffset:
		b := self.FieldTable["baseClock"].Fvalue.(*object.Object)
		off := self.FieldTable["offset"].Fvalue.(*object.Object)
		bStr := object.GoStringFromStringObject(clockToString([]any{b}).(*object.Object))
		return object.StringObjectFromGoString("OffsetClock[" + bStr + "," + getDurationString(off) + "]")
	case clockTypeTick:
		b := self.FieldTable["baseClock"].Fvalue.(*object.Object)
		tn := self.FieldTable["tickNanos"].Fvalue.(int64)
		durObj := createDuration(tn/1000000000, int32(tn%1000000000))
		bStr := object.GoStringFromStringObject(clockToString([]any{b}).(*object.Object))
		return object.StringObjectFromGoString("TickClock[" + bStr + "," + getDurationString(durObj) + "]")
	default:
		return object.StringObjectFromGoString("Clock[]")
	}
}

func getDurationSecondsVal(d *object.Object) int64 {
	if d == nil || d.FieldTable == nil {
		return 0
	}
	if f, ok := d.FieldTable["seconds"]; ok {
		if v, ok2 := getInt64Val(f.Fvalue); ok2 {
			return v
		}
	}
	return 0
}

func getDurationNanosVal(d *object.Object) int32 {
	if d == nil || d.FieldTable == nil {
		return 0
	}
	if f, ok := d.FieldTable["nanos"]; ok {
		if v, ok2 := getInt64Val(f.Fvalue); ok2 {
			return int32(v)
		}
	}
	return 0
}

func getInstantEpochSecondVal(i *object.Object) int64 {
	if i == nil || i.FieldTable == nil {
		return 0
	}
	if f, ok := i.FieldTable["seconds"]; ok {
		if v, ok2 := getInt64Val(f.Fvalue); ok2 {
			return v
		}
	}
	return 0
}

func getInstantNanoVal(i *object.Object) int64 {
	if i == nil || i.FieldTable == nil {
		return 0
	}
	if f, ok := i.FieldTable["nanos"]; ok {
		if v, ok2 := getInt64Val(f.Fvalue); ok2 {
			return v
		}
	}
	return 0
}

func zoneEquals(z1, z2 *object.Object) bool {
	if z1 == z2 {
		return true
	}
	if z1 == nil || z2 == nil {
		return false
	}
	return getZoneIdString(z1) == getZoneIdString(z2)
}

func zoneHashCodeVal(z *object.Object) int64 {
	if z == nil {
		return 0
	}
	id := getZoneIdString(z)
	var h int32 = 0
	for i := 0; i < len(id); i++ {
		h = 31*h + int32(id[i])
	}
	return int64(h)
}

func getZoneString(z *object.Object) string {
	if z == nil {
		return "UTC"
	}
	return getZoneIdString(z)
}

func instantEqualsObjects(i1, i2 *object.Object) bool {
	if i1 == i2 {
		return true
	}
	if i1 == nil || i2 == nil {
		return false
	}
	return getInstantEpochSecondVal(i1) == getInstantEpochSecondVal(i2) && getInstantNanoVal(i1) == getInstantNanoVal(i2)
}

func instantHashCodeVal(i *object.Object) int64 {
	if i == nil {
		return 0
	}
	sec := getInstantEpochSecondVal(i)
	nano := getInstantNanoVal(i)
	return (sec ^ (sec >> 32)) + 51*nano
}

func getInstantString(i *object.Object) string {
	if i == nil {
		return ""
	}
	strObj := instantToString([]any{i}).(*object.Object)
	return object.GoStringFromStringObject(strObj)
}

func durationEqualsObjects(d1, d2 *object.Object) bool {
	if d1 == d2 {
		return true
	}
	if d1 == nil || d2 == nil {
		return false
	}
	return getDurationSecondsVal(d1) == getDurationSecondsVal(d2) && getDurationNanosVal(d1) == getDurationNanosVal(d2)
}

func durationHashCodeVal(d *object.Object) int64 {
	if d == nil {
		return 0
	}
	sec := getDurationSecondsVal(d)
	nano := int64(getDurationNanosVal(d))
	return (sec ^ (sec >> 32)) + 51*nano
}

func getDurationString(d *object.Object) string {
	if d == nil {
		return ""
	}
	strObj := durationToString([]any{d}).(*object.Object)
	return object.GoStringFromStringObject(strObj)
}
