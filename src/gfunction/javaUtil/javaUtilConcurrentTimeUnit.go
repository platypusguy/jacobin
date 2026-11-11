/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"math"
	"sync"
	"time"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
)

var timeUnitClassName = "java/util/concurrent/TimeUnit"

// TimeUnit constant names
const (
	NANOSECONDS  = "NANOSECONDS"
	MICROSECONDS = "MICROSECONDS"
	MILLISECONDS = "MILLISECONDS"
	SECONDS      = "SECONDS"
	MINUTES      = "MINUTES"
	HOURS        = "HOURS"
	DAYS         = "DAYS"
)

var timeUnitNames = []string{
	NANOSECONDS,
	MICROSECONDS,
	MILLISECONDS,
	SECONDS,
	MINUTES,
	HOURS,
	DAYS,
}

var timeUnitScaleInNanos = map[string]int64{
	NANOSECONDS:  1,
	MICROSECONDS: 1000,
	MILLISECONDS: 1000000,
	SECONDS:      1000000000,
	MINUTES:      60000000000,
	HOURS:        3600000000000,
	DAYS:         86400000000000,
}

var chronoUnitToTimeUnit = map[string]string{
	"NANOS":   NANOSECONDS,
	"MICROS":  MICROSECONDS,
	"MILLIS":  MILLISECONDS,
	"SECONDS": SECONDS,
	"MINUTES": MINUTES,
	"HOURS":   HOURS,
	"DAYS":    DAYS,
}

var timeUnitToChronoUnitName = map[string]string{
	NANOSECONDS:  "NANOS",
	MICROSECONDS: "MICROS",
	MILLISECONDS: "MILLIS",
	SECONDS:      "SECONDS",
	MINUTES:      "MINUTES",
	HOURS:        "HOURS",
	DAYS:         "DAYS",
}

var (
	timeUnitMutex           sync.Mutex
	timeUnitOnceInitialized bool
	timeUnitInstances       []*object.Object
)

// Load_Util_Concurrent_TimeUnit registers all java.util.concurrent.TimeUnit methods.
func Load_Util_Concurrent_TimeUnit() {
	// <clinit>()V
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.<clinit>()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: timeUnitClinit}

	// values()[Ljava/util/concurrent/TimeUnit;
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.values()[Ljava/util/concurrent/TimeUnit;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: timeUnitValues}

	// valueOf(Ljava/lang/String;)Ljava/util/concurrent/TimeUnit;
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.valueOf(Ljava/lang/String;)Ljava/util/concurrent/TimeUnit;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: timeUnitValueOfString}

	// of(Ljava/time/temporal/ChronoUnit;)Ljava/util/concurrent/TimeUnit;
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.of(Ljava/time/temporal/ChronoUnit;)Ljava/util/concurrent/TimeUnit;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: timeUnitOfChronoUnit}

	// convert(JLjava/util/concurrent/TimeUnit;)J
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.convert(JLjava/util/concurrent/TimeUnit;)J"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: timeUnitConvert}

	// convert(Ljava/time/Duration;)J
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.convert(Ljava/time/Duration;)J"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: timeUnitConvertDuration}

	// toNanos(J)J
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.toNanos(J)J"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: timeUnitToNanos}

	// toMicros(J)J
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.toMicros(J)J"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: timeUnitToMicros}

	// toMillis(J)J
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.toMillis(J)J"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: timeUnitToMillis}

	// toSeconds(J)J
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.toSeconds(J)J"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: timeUnitToSeconds}

	// toMinutes(J)J
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.toMinutes(J)J"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: timeUnitToMinutes}

	// toHours(J)J
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.toHours(J)J"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: timeUnitToHours}

	// toDays(J)J
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.toDays(J)J"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: timeUnitToDays}

	// toChronoUnit()Ljava/time/temporal/ChronoUnit;
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.toChronoUnit()Ljava/time/temporal/ChronoUnit;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: timeUnitToChronoUnit}

	// sleep(J)V
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.sleep(J)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: timeUnitSleep}

	// timedWait(Ljava/lang/Object;J)V
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.timedWait(Ljava/lang/Object;J)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: timeUnitTimedWait}

	// timedJoin(Ljava/lang/Thread;J)V
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.timedJoin(Ljava/lang/Thread;J)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: timeUnitTimedJoin}

	// name()Ljava/lang/String;
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.name()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: timeUnitName}

	// ordinal()I
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.ordinal()I"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: timeUnitOrdinal}

	// toString()Ljava/lang/String;
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: timeUnitToString}

	// compareTo(Ljava/lang/Enum;)I
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.compareTo(Ljava/lang/Enum;)I"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: timeUnitCompareTo}

	// compareTo(Ljava/lang/Object;)I
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: timeUnitCompareTo}

	// equals(Ljava/lang/Object;)Z
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: timeUnitEquals}

	// hashCode()I
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.hashCode()I"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: timeUnitHashCode}

	// Standard Enum & Object lifecycle methods
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.getDeclaringClass()Ljava/lang/Class;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.clone()Ljava/lang/Object;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.finalize()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.getClass()Ljava/lang/Class;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.describeConstable()Ljava/util/Optional;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.notify()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.notifyAll()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.wait()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.wait(J)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/util/concurrent/TimeUnit.wait(JI)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
}

// ensureTimeUnitInited lazily creates singleton instances for all TimeUnit constants.
func ensureTimeUnitInited() {
	timeUnitMutex.Lock()
	defer timeUnitMutex.Unlock()
	if timeUnitOnceInitialized {
		return
	}
	timeUnitInstances = make([]*object.Object, len(timeUnitNames))
	for i, nm := range timeUnitNames {
		obj := object.MakeEmptyObjectWithClassName(&timeUnitClassName)
		obj.FieldTable["name"] = object.Field{Ftype: types.StringClassRef, Fvalue: object.StringObjectFromGoString(nm)}
		obj.FieldTable["ordinal"] = object.Field{Ftype: types.Int, Fvalue: int64(i)}
		timeUnitInstances[i] = obj
		_ = statics.AddStatic(timeUnitClassName+"."+nm, statics.Static{Type: "Ljava/util/concurrent/TimeUnit;", Value: obj})
	}
	timeUnitOnceInitialized = true
}

// GetTimeUnitInstance returns the singleton Object for a given TimeUnit name.
func GetTimeUnitInstance(name string) *object.Object {
	ensureTimeUnitInited()
	for i, nm := range timeUnitNames {
		if nm == name {
			return timeUnitInstances[i]
		}
	}
	return nil
}

// GetTimeUnitByOrdinal returns the singleton Object for a given ordinal index (0..6).
func GetTimeUnitByOrdinal(ord int) *object.Object {
	ensureTimeUnitInited()
	if ord >= 0 && ord < len(timeUnitInstances) {
		return timeUnitInstances[ord]
	}
	return nil
}

func parseTimeUnit(arg any) (string, int, *ghelpers.GErrBlk) {
	if arg == nil {
		return "", -1, ghelpers.GetGErrBlk(excNames.NullPointerException, "TimeUnit: unit is null")
	}
	if obj, ok := arg.(*object.Object); ok {
		if object.IsNull(obj) {
			return "", -1, ghelpers.GetGErrBlk(excNames.NullPointerException, "TimeUnit: unit is null")
		}
		if object.IsStringObject(obj) {
			unitStr := object.GoStringFromStringObject(obj)
			for i, nm := range timeUnitNames {
				if nm == unitStr {
					return unitStr, i, nil
				}
			}
			return "", -1, ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit: invalid TimeUnit: "+unitStr)
		}
		if fld, exists := obj.FieldTable["name"]; exists {
			var name string
			if sObj, ok2 := fld.Fvalue.(*object.Object); ok2 && !object.IsNull(sObj) {
				name = object.GoStringFromStringObject(sObj)
			} else if s, ok2 := fld.Fvalue.(string); ok2 {
				name = s
			}
			for i, nm := range timeUnitNames {
				if nm == name {
					return name, i, nil
				}
			}
		}
		if fld, exists := obj.FieldTable["ordinal"]; exists {
			if ord, ok2 := ghelpers.GetInt64Val(fld.Fvalue); ok2 && ord >= 0 && int(ord) < len(timeUnitNames) {
				return timeUnitNames[ord], int(ord), nil
			}
		}
	} else if str, ok := arg.(string); ok {
		for i, nm := range timeUnitNames {
			if nm == str {
				return str, i, nil
			}
		}
		return "", -1, ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit: invalid TimeUnit: "+str)
	}
	return "", -1, ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit: invalid TimeUnit argument")
}

// ConvertTimeUnit converts a duration from srcUnit to dstUnit with saturation on overflow.
func ConvertTimeUnit(val int64, srcUnit, dstUnit string) int64 {
	if srcUnit == dstUnit {
		return val
	}
	srcScale := timeUnitScaleInNanos[srcUnit]
	dstScale := timeUnitScaleInNanos[dstUnit]
	if srcScale > dstScale {
		// Converting to a smaller unit: multiplication with saturation
		scale := srcScale / dstScale
		if val > math.MaxInt64/scale {
			return math.MaxInt64
		}
		if val < math.MinInt64/scale {
			return math.MinInt64
		}
		return val * scale
	}
	// Converting to a larger unit: division
	scale := dstScale / srcScale
	return val / scale
}

func timeUnitClinit([]interface{}) interface{} {
	ensureTimeUnitInited()
	return nil
}

func timeUnitValues([]interface{}) interface{} {
	ensureTimeUnitInited()
	arr := object.Make1DimRefArray("Ljava/util/concurrent/TimeUnit;", int64(len(timeUnitInstances)))
	slot := arr.FieldTable["value"].Fvalue.([]*object.Object)
	copy(slot, timeUnitInstances)
	arr.FieldTable["value"] = object.Field{Ftype: types.RefArray + "Ljava/util/concurrent/TimeUnit;", Fvalue: slot}
	return arr
}

func timeUnitValueOfString(params []interface{}) interface{} {
	ensureTimeUnitInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TimeUnit.valueOf: name is null")
	}
	strObj, ok := params[0].(*object.Object)
	if !ok || !object.IsStringObject(strObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.valueOf: argument is not a String")
	}
	name := object.GoStringFromStringObject(strObj)
	for i, nm := range timeUnitNames {
		if nm == name {
			return timeUnitInstances[i]
		}
	}
	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.valueOf: no enum constant "+name)
}

func timeUnitOfChronoUnit(params []interface{}) interface{} {
	ensureTimeUnitInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TimeUnit.of: chronoUnit is null")
	}
	cuObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.of: argument is not an object")
	}
	var name string
	if fld, ok2 := cuObj.FieldTable["name"]; ok2 {
		if sObj, ok3 := fld.Fvalue.(*object.Object); ok3 && !object.IsNull(sObj) {
			name = object.GoStringFromStringObject(sObj)
		} else if s, ok3 := fld.Fvalue.(string); ok3 {
			name = s
		}
	}
	tuName, found := chronoUnitToTimeUnit[name]
	if !found {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.of: No TimeUnit equivalent for ChronoUnit "+name)
	}
	return GetTimeUnitInstance(tuName)
}

func timeUnitConvert(params []interface{}) interface{} {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.convert: insufficient parameters")
	}
	dstUnit, _, gerr := parseTimeUnit(params[0])
	if gerr != nil {
		return gerr
	}
	dur, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.convert: invalid duration")
	}
	srcUnit, _, gerr2 := parseTimeUnit(params[2])
	if gerr2 != nil {
		return gerr2
	}
	return ConvertTimeUnit(dur, srcUnit, dstUnit)
}

func timeUnitConvertDuration(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.convert: insufficient parameters")
	}
	dstUnit, _, gerr := parseTimeUnit(params[0])
	if gerr != nil {
		return gerr
	}
	if params[1] == nil || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TimeUnit.convert: duration is null")
	}
	durObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.convert: duration is not an object")
	}
	var seconds, nanos int64
	if sFld, ok2 := durObj.FieldTable["seconds"]; ok2 {
		if s, ok3 := ghelpers.GetInt64Val(sFld.Fvalue); ok3 {
			seconds = s
		}
	}
	if nFld, ok2 := durObj.FieldTable["nanos"]; ok2 {
		if n, ok3 := ghelpers.GetInt64Val(nFld.Fvalue); ok3 {
			nanos = n
		}
	}

	// Convert seconds and nanos to dstUnit
	switch dstUnit {
	case NANOSECONDS:
		if seconds > math.MaxInt64/1000000000 {
			return int64(math.MaxInt64)
		}
		if seconds < math.MinInt64/1000000000 {
			return int64(math.MinInt64)
		}
		secNanos := seconds * 1000000000
		if secNanos > 0 && math.MaxInt64-secNanos < nanos {
			return int64(math.MaxInt64)
		}
		if secNanos < 0 && math.MinInt64-secNanos > nanos {
			return int64(math.MinInt64)
		}
		return secNanos + nanos
	case MICROSECONDS:
		if seconds > math.MaxInt64/1000000 {
			return int64(math.MaxInt64)
		}
		if seconds < math.MinInt64/1000000 {
			return int64(math.MinInt64)
		}
		secMicros := seconds * 1000000
		microNanos := nanos / 1000
		if secMicros > 0 && math.MaxInt64-secMicros < microNanos {
			return int64(math.MaxInt64)
		}
		if secMicros < 0 && math.MinInt64-secMicros > microNanos {
			return int64(math.MinInt64)
		}
		return secMicros + microNanos
	case MILLISECONDS:
		if seconds > math.MaxInt64/1000 {
			return int64(math.MaxInt64)
		}
		if seconds < math.MinInt64/1000 {
			return int64(math.MinInt64)
		}
		secMillis := seconds * 1000
		milliNanos := nanos / 1000000
		if secMillis > 0 && math.MaxInt64-secMillis < milliNanos {
			return int64(math.MaxInt64)
		}
		if secMillis < 0 && math.MinInt64-secMillis > milliNanos {
			return int64(math.MinInt64)
		}
		return secMillis + milliNanos
	case SECONDS:
		return seconds
	case MINUTES:
		return seconds / 60
	case HOURS:
		return seconds / 3600
	case DAYS:
		return seconds / 86400
	default:
		return int64(0)
	}
}

func timeUnitToNanos(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.toNanos: insufficient parameters")
	}
	srcUnit, _, gerr := parseTimeUnit(params[0])
	if gerr != nil {
		return gerr
	}
	dur, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.toNanos: invalid duration")
	}
	return ConvertTimeUnit(dur, srcUnit, NANOSECONDS)
}

func timeUnitToMicros(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.toMicros: insufficient parameters")
	}
	srcUnit, _, gerr := parseTimeUnit(params[0])
	if gerr != nil {
		return gerr
	}
	dur, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.toMicros: invalid duration")
	}
	return ConvertTimeUnit(dur, srcUnit, MICROSECONDS)
}

func timeUnitToMillis(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.toMillis: insufficient parameters")
	}
	srcUnit, _, gerr := parseTimeUnit(params[0])
	if gerr != nil {
		return gerr
	}
	dur, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.toMillis: invalid duration")
	}
	return ConvertTimeUnit(dur, srcUnit, MILLISECONDS)
}

func timeUnitToSeconds(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.toSeconds: insufficient parameters")
	}
	srcUnit, _, gerr := parseTimeUnit(params[0])
	if gerr != nil {
		return gerr
	}
	dur, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.toSeconds: invalid duration")
	}
	return ConvertTimeUnit(dur, srcUnit, SECONDS)
}

func timeUnitToMinutes(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.toMinutes: insufficient parameters")
	}
	srcUnit, _, gerr := parseTimeUnit(params[0])
	if gerr != nil {
		return gerr
	}
	dur, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.toMinutes: invalid duration")
	}
	return ConvertTimeUnit(dur, srcUnit, MINUTES)
}

func timeUnitToHours(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.toHours: insufficient parameters")
	}
	srcUnit, _, gerr := parseTimeUnit(params[0])
	if gerr != nil {
		return gerr
	}
	dur, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.toHours: invalid duration")
	}
	return ConvertTimeUnit(dur, srcUnit, HOURS)
}

func timeUnitToDays(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.toDays: insufficient parameters")
	}
	srcUnit, _, gerr := parseTimeUnit(params[0])
	if gerr != nil {
		return gerr
	}
	dur, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.toDays: invalid duration")
	}
	return ConvertTimeUnit(dur, srcUnit, DAYS)
}

func timeUnitToChronoUnit(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.toChronoUnit: missing self parameter")
	}
	tuName, _, gerr := parseTimeUnit(params[0])
	if gerr != nil {
		return gerr
	}
	cuName, ok := timeUnitToChronoUnitName[tuName]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.toChronoUnit: invalid TimeUnit "+tuName)
	}
	return ghelpers.CreateTemporalUnit(cuName)
}

func timeUnitSleep(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.sleep: insufficient parameters")
	}
	srcUnit, _, gerr := parseTimeUnit(params[0])
	if gerr != nil {
		return gerr
	}
	dur, ok := ghelpers.GetInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.sleep: invalid timeout")
	}
	if dur > 0 {
		nanos := ConvertTimeUnit(dur, srcUnit, NANOSECONDS)
		if nanos > 0 {
			time.Sleep(time.Duration(nanos))
		}
	}
	return nil
}

func timeUnitTimedWait(params []interface{}) interface{} {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.timedWait: insufficient parameters")
	}
	_, _, gerr := parseTimeUnit(params[0])
	if gerr != nil {
		return gerr
	}
	if params[1] == nil || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TimeUnit.timedWait: object is null")
	}
	timeout, ok := ghelpers.GetInt64Val(params[2])
	if !ok || timeout < 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.timedWait: timeout value is negative")
	}
	return nil
}

func timeUnitTimedJoin(params []interface{}) interface{} {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.timedJoin: insufficient parameters")
	}
	_, _, gerr := parseTimeUnit(params[0])
	if gerr != nil {
		return gerr
	}
	if params[1] == nil || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TimeUnit.timedJoin: thread is null")
	}
	timeout, ok := ghelpers.GetInt64Val(params[2])
	if !ok || timeout < 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.timedJoin: timeout value is negative")
	}
	return nil
}

func timeUnitName(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.name: missing self parameter")
	}
	name, _, gerr := parseTimeUnit(params[0])
	if gerr != nil {
		return gerr
	}
	return object.StringObjectFromGoString(name)
}

func timeUnitOrdinal(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.ordinal: missing self parameter")
	}
	_, ord, gerr := parseTimeUnit(params[0])
	if gerr != nil {
		return gerr
	}
	return int64(ord)
}

func timeUnitToString(params []interface{}) interface{} {
	return timeUnitName(params)
}

func timeUnitCompareTo(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TimeUnit.compareTo: insufficient parameters")
	}
	_, selfOrd, gerr := parseTimeUnit(params[0])
	if gerr != nil {
		return gerr
	}
	_, otherOrd, gerr2 := parseTimeUnit(params[1])
	if gerr2 != nil {
		return gerr2
	}
	return int64(selfOrd - otherOrd)
}

func timeUnitEquals(params []interface{}) interface{} {
	if len(params) < 2 {
		return types.JavaBoolFalse
	}
	if params[0] == nil || params[1] == nil || object.IsNull(params[0]) || object.IsNull(params[1]) {
		return types.JavaBoolFalse
	}
	selfName, _, gerr1 := parseTimeUnit(params[0])
	otherName, _, gerr2 := parseTimeUnit(params[1])
	if gerr1 != nil || gerr2 != nil {
		return types.JavaBoolFalse
	}
	if selfName == otherName {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func timeUnitHashCode(params []interface{}) interface{} {
	if len(params) < 1 {
		return int64(0)
	}
	_, ord, gerr := parseTimeUnit(params[0])
	if gerr != nil {
		return int64(0)
	}
	return int64(ord)
}
