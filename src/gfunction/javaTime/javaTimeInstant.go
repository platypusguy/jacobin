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
	"jacobin/src/statics"
	"jacobin/src/types"
	"math"
	"strings"
	"sync"
	"time"
)

var (
	classNameInstantVal        = "java/time/Instant"
	classNameOffsetDateTimeVal = "java/time/OffsetDateTime"
	classNameZonedDateTimeVal  = "java/time/ZonedDateTime"
)

const (
	instantMinSecond   int64 = -31557014167219200
	instantMaxSecond   int64 = 31556889864403199
	instantNanosPerSec int64 = 1000000000
)

var (
	instantMutex       sync.Mutex
	instantInitialized bool = false
	instantEpoch       *object.Object
	instantMin         *object.Object
	instantMax         *object.Object
)

func ensureInstantInited() {
	instantMutex.Lock()
	defer instantMutex.Unlock()
	if instantInitialized {
		return
	}

	instantEpoch = createInstant(0, 0)
	instantMin = createInstant(instantMinSecond, 0)
	instantMax = createInstant(instantMaxSecond, 999999999)

	_ = statics.AddStatic(classNameInstantVal+".EPOCH", statics.Static{
		Type:  "Ljava/time/Instant;",
		Value: instantEpoch,
	})
	_ = statics.AddStatic(classNameInstantVal+".MIN", statics.Static{
		Type:  "Ljava/time/Instant;",
		Value: instantMin,
	})
	_ = statics.AddStatic(classNameInstantVal+".MAX", statics.Static{
		Type:  "Ljava/time/Instant;",
		Value: instantMax,
	})

	instantInitialized = true
}

func getChronoFieldForInstant(name string) *object.Object {
	ensureChronoFieldInited()
	for index, definition := range chronoFieldDefinitions {
		if definition.enumName == name {
			return chronoFieldInstances[index]
		}
	}
	return nil
}

// Load_Time_Instant registers method signatures for java/time/Instant.
func Load_Time_Instant() {
	ensureInstantInited()

	ghelpers.MethodSignatures["java/time/Instant.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  instantClinit,
		}

	ghelpers.MethodSignatures["java/time/Instant.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  instantInit,
		}

	ghelpers.MethodSignatures["java/time/Instant.now()Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  instantNow,
		}

	ghelpers.MethodSignatures["java/time/Instant.now(Ljava/time/Clock;)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantNowClock,
		}

	ghelpers.MethodSignatures["java/time/Instant.ofEpochSecond(J)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantOfEpochSecond1,
		}

	ghelpers.MethodSignatures["java/time/Instant.ofEpochSecond(JJ)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  instantOfEpochSecond2,
		}

	ghelpers.MethodSignatures["java/time/Instant.ofEpochMilli(J)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantOfEpochMilli,
		}

	ghelpers.MethodSignatures["java/time/Instant.from(Ljava/time/temporal/TemporalAccessor;)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantFrom,
		}

	ghelpers.MethodSignatures["java/time/Instant.parse(Ljava/lang/CharSequence;)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantParse,
		}

	ghelpers.MethodSignatures["java/time/Instant.isSupported(Ljava/time/temporal/TemporalField;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantIsSupportedField,
		}

	ghelpers.MethodSignatures["java/time/Instant.isSupported(Ljava/time/temporal/TemporalUnit;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantIsSupportedUnit,
		}

	ghelpers.MethodSignatures["java/time/Instant.range(Ljava/time/temporal/TemporalField;)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantRange,
		}

	ghelpers.MethodSignatures["java/time/Instant.get(Ljava/time/temporal/TemporalField;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantGet,
		}

	ghelpers.MethodSignatures["java/time/Instant.getLong(Ljava/time/temporal/TemporalField;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantGetLong,
		}

	ghelpers.MethodSignatures["java/time/Instant.getEpochSecond()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  instantGetEpochSecond,
		}

	ghelpers.MethodSignatures["java/time/Instant.getNano()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  instantGetNano,
		}

	ghelpers.MethodSignatures["java/time/Instant.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/Instant.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/Instant.with(Ljava/time/temporal/TemporalField;J)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  instantWithField,
		}

	ghelpers.MethodSignatures["java/time/Instant.with(Ljava/time/temporal/TemporalField;J)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  instantWithField,
		}

	ghelpers.MethodSignatures["java/time/Instant.truncatedTo(Ljava/time/temporal/TemporalUnit;)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantTruncatedTo,
		}

	ghelpers.MethodSignatures["java/time/Instant.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantPlusAmount,
		}

	ghelpers.MethodSignatures["java/time/Instant.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantPlusAmount,
		}

	ghelpers.MethodSignatures["java/time/Instant.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  instantPlusUnit,
		}

	ghelpers.MethodSignatures["java/time/Instant.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  instantPlusUnit,
		}

	ghelpers.MethodSignatures["java/time/Instant.plusSeconds(J)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantPlusSeconds,
		}

	ghelpers.MethodSignatures["java/time/Instant.plusMillis(J)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantPlusMillis,
		}

	ghelpers.MethodSignatures["java/time/Instant.plusNanos(J)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantPlusNanos,
		}

	ghelpers.MethodSignatures["java/time/Instant.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/Instant.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/Instant.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  instantMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/Instant.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  instantMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/Instant.minusSeconds(J)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantMinusSeconds,
		}

	ghelpers.MethodSignatures["java/time/Instant.minusMillis(J)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantMinusMillis,
		}

	ghelpers.MethodSignatures["java/time/Instant.minusNanos(J)Ljava/time/Instant;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantMinusNanos,
		}

	ghelpers.MethodSignatures["java/time/Instant.query(Ljava/time/temporal/TemporalQuery;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantQuery,
		}

	ghelpers.MethodSignatures["java/time/Instant.adjustInto(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantAdjustInto,
		}

	ghelpers.MethodSignatures["java/time/Instant.until(Ljava/time/temporal/Temporal;Ljava/time/temporal/TemporalUnit;)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  instantUntil,
		}

	ghelpers.MethodSignatures["java/time/Instant.atOffset(Ljava/time/ZoneOffset;)Ljava/time/OffsetDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantAtOffset,
		}

	ghelpers.MethodSignatures["java/time/Instant.atZone(Ljava/time/ZoneId;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantAtZone,
		}

	ghelpers.MethodSignatures["java/time/Instant.toEpochMilli()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  instantToEpochMilli,
		}

	ghelpers.MethodSignatures["java/time/Instant.compareTo(Ljava/time/Instant;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantCompareTo,
		}

	ghelpers.MethodSignatures["java/time/Instant.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantCompareTo,
		}

	ghelpers.MethodSignatures["java/time/Instant.isAfter(Ljava/time/Instant;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantIsAfter,
		}

	ghelpers.MethodSignatures["java/time/Instant.isBefore(Ljava/time/Instant;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantIsBefore,
		}

	ghelpers.MethodSignatures["java/time/Instant.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  instantEquals,
		}

	ghelpers.MethodSignatures["java/time/Instant.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  instantHashCode,
		}

	ghelpers.MethodSignatures["java/time/Instant.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  instantToString,
		}
}

func instantClinit([]any) any {
	ensureInstantInited()
	return nil
}

func instantInit([]any) any {
	return nil
}

// CreateInstant creates an Instant object
func CreateInstant(seconds int64, nanos int64) *object.Object {
	return createInstant(seconds, nanos)
}

func createInstant(seconds int64, nanos int64) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameInstantVal)
	obj.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: seconds}
	obj.FieldTable["nanos"] = object.Field{Ftype: types.Int, Fvalue: nanos}
	return obj
}

func createInstantValidated(seconds int64, nanos int64) any {
	if seconds < instantMinSecond || seconds > instantMaxSecond {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Instant exceeds minimum or maximum instant: %d", seconds))
	}
	if nanos < 0 || nanos > 999999999 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid nanosecond value: %d", nanos))
	}
	return createInstant(seconds, nanos)
}

func createInstantFromNanos(epochSecond int64, nanoAdjustment int64) any {
	secs := floorDiv(nanoAdjustment, instantNanosPerSec)
	nos := floorMod(nanoAdjustment, instantNanosPerSec)

	totalSec, ok := addExactSafe(epochSecond, secs)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Instant exceeds minimum or maximum instant")
	}
	return createInstantValidated(totalSec, nos)
}

func getInstantSecondsAndNanos(obj *object.Object) (int64, int64) {
	if obj == nil || obj.FieldTable == nil {
		return 0, 0
	}
	var seconds, nanos int64
	if sFld, ok := obj.FieldTable["seconds"]; ok {
		if s, ok2 := getInt64Val(sFld.Fvalue); ok2 {
			seconds = s
		}
	} else if sFld, ok := obj.FieldTable["epochSecond"]; ok {
		if s, ok2 := getInt64Val(sFld.Fvalue); ok2 {
			seconds = s
		}
	}
	if nFld, ok := obj.FieldTable["nanos"]; ok {
		if n, ok2 := getInt64Val(nFld.Fvalue); ok2 {
			nanos = n
		}
	} else if nFld, ok := obj.FieldTable["nano"]; ok {
		if n, ok2 := getInt64Val(nFld.Fvalue); ok2 {
			nanos = n
		}
	}
	return seconds, nanos
}

func addExactSafe(x, y int64) (int64, bool) {
	r := x + y
	if ((x ^ r) & (y ^ r)) < 0 {
		return 0, false
	}
	return r, true
}

func multiplyExactSafe(x, y int64) (int64, bool) {
	if y == 0 || x == 0 {
		return 0, true
	}
	r := x * y
	if r/y != x {
		return 0, false
	}
	return r, true
}

func instantNow([]any) any {
	now := time.Now().UTC()
	return createInstant(now.Unix(), int64(now.Nanosecond()))
}

func instantNowClock(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.now: clock is null")
	}
	clock := params[0].(*object.Object)
	return clockInstant([]any{clock})
}

func instantOfEpochSecond1(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Instant.ofEpochSecond: missing argument")
	}
	sec, ok := getInt64Val(params[0])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Instant.ofEpochSecond: invalid epochSecond")
	}
	return createInstantValidated(sec, 0)
}

func instantOfEpochSecond2(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Instant.ofEpochSecond: missing arguments")
	}
	sec, ok1 := getInt64Val(params[0])
	nanoAdj, ok2 := getInt64Val(params[1])
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Instant.ofEpochSecond: invalid arguments")
	}
	return createInstantFromNanos(sec, nanoAdj)
}

func instantOfEpochMilli(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Instant.ofEpochMilli: missing argument")
	}
	milli, ok := getInt64Val(params[0])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Instant.ofEpochMilli: invalid epochMilli")
	}
	secs := floorDiv(milli, 1000)
	nos := floorMod(milli, 1000) * 1000000
	return createInstantValidated(secs, nos)
}

func instantFrom(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.from: temporal is null")
	}
	temporal := params[0].(*object.Object)
	kName := object.GoStringFromStringPoolIndex(temporal.KlassName)
	if kName == classNameInstantVal {
		return temporal
	}

	// Try extracting INSTANT_SECONDS and NANO_OF_SECOND directly if present
	if temporal.FieldTable != nil {
		if sFld, ok := temporal.FieldTable["seconds"]; ok {
			if sec, ok2 := getInt64Val(sFld.Fvalue); ok2 {
				var nano int64 = 0
				if nFld, okN := temporal.FieldTable["nanos"]; okN {
					if n, okN2 := getInt64Val(nFld.Fvalue); okN2 {
						nano = n
					}
				}
				return createInstant(sec, nano)
			}
		}
	}

	// Try ChronoZonedDateTime or ZonedDateTime
	if kName == "java/time/chrono/ChronoZonedDateTime" || kName == "java/time/ZonedDateTime" {
		toInstMeth, ok := ghelpers.MethodSignatures["java/time/chrono/ChronoZonedDateTime.toInstant()Ljava/time/Instant;"]
		if ok {
			return toInstMeth.GFunction([]any{temporal})
		}
	}

	// Check via getLong on INSTANT_SECONDS
	cfInstantSec := getChronoFieldForInstant("INSTANT_SECONDS")
	cfNanoOfSec := getChronoFieldForInstant("NANO_OF_SECOND")
	if cfInstantSec != nil {
		isSupp := temporalAccessorIsSupported([]any{temporal, cfInstantSec})
		if isSupp == types.JavaBoolTrue {
			secRes := temporalAccessorGetLong([]any{temporal, cfInstantSec})
			if secVal, ok := getInt64Val(secRes); ok {
				var nanoVal int64 = 0
				if cfNanoOfSec != nil {
					isSuppNano := temporalAccessorIsSupported([]any{temporal, cfNanoOfSec})
					if isSuppNano == types.JavaBoolTrue {
						nanoRes := temporalAccessorGetLong([]any{temporal, cfNanoOfSec})
						if n, okN := getInt64Val(nanoRes); okN {
							nanoVal = n
						}
					}
				}
				return createInstant(secVal, nanoVal)
			}
		}
	}

	return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unable to obtain Instant from TemporalAccessor: "+kName)
}

func instantParse(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.parse: text is null")
	}
	textObj := params[0].(*object.Object)
	text := object.GoStringFromStringObject(textObj)
	text = strings.TrimSpace(text)

	// Format: <date>T<time>[Z|+HH:MM|-HH:MM]
	tIdx := strings.IndexAny(text, "Tt")
	if tIdx < 0 {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Text '"+text+"' could not be parsed: missing 'T' separator")
	}

	datePart := text[:tIdx]
	timePart := text[tIdx+1:]

	y, m, d, errDate := parseLocalDateString(datePart)
	if errDate != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Text '"+text+"' could not be parsed: "+errDate.Error())
	}

	// Parse offset from timePart
	var offsetSec int64 = 0
	var timeOnly string

	if strings.HasSuffix(timePart, "Z") || strings.HasSuffix(timePart, "z") {
		timeOnly = timePart[:len(timePart)-1]
		offsetSec = 0
	} else if plusIdx := strings.LastIndex(timePart, "+"); plusIdx >= 0 {
		timeOnly = timePart[:plusIdx]
		offStr := timePart[plusIdx:]
		offObj := zoneOffsetOf([]any{offStr})
		if err, isErr := offObj.(*ghelpers.GErrBlk); isErr {
			return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Text '"+text+"' could not be parsed: "+err.Error())
		}
		offsetSec = getZoneOffsetTotalSeconds(offObj.(*object.Object))
	} else if minusIdx := strings.LastIndex(timePart, "-"); minusIdx >= 0 {
		timeOnly = timePart[:minusIdx]
		offStr := timePart[minusIdx:]
		offObj := zoneOffsetOf([]any{offStr})
		if err, isErr := offObj.(*ghelpers.GErrBlk); isErr {
			return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Text '"+text+"' could not be parsed: "+err.Error())
		}
		offsetSec = getZoneOffsetTotalSeconds(offObj.(*object.Object))
	} else {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Text '"+text+"' could not be parsed: missing zone offset")
	}

	h, min, s, n, errTime := parseLocalTimeString(timeOnly)
	if errTime != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Text '"+text+"' could not be parsed: "+errTime.Error())
	}

	epochDay := ymdToEpochDay(y, m, d)
	localSec := epochDay*86400 + h*3600 + min*60 + s
	epochSec := localSec - offsetSec

	return createInstantValidated(epochSec, n)
}

func instantIsSupportedField(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.isSupported: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	fldName := getTemporalFieldName(params[1].(*object.Object))
	switch strings.ToUpper(fldName) {
	case "INSTANT_SECONDS", "NANO_OF_SECOND", "MICRO_OF_SECOND", "MILLI_OF_SECOND":
		return types.JavaBoolTrue
	default:
		return types.JavaBoolFalse
	}
}

func instantIsSupportedUnit(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.isSupported: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	unitName := getTemporalUnitName(params[1].(*object.Object))
	switch strings.ToUpper(unitName) {
	case "NANOS", "MICROS", "MILLIS", "SECONDS", "MINUTES", "HOURS", "HALF_DAYS", "DAYS":
		return types.JavaBoolTrue
	default:
		return types.JavaBoolFalse
	}
}

func instantRange(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.range: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.range: field is null")
	}
	fldName := getTemporalFieldName(params[1].(*object.Object))
	switch strings.ToUpper(fldName) {
	case "INSTANT_SECONDS":
		return createValueRange(instantMinSecond, instantMinSecond, instantMaxSecond, instantMaxSecond)
	case "NANO_OF_SECOND":
		return createValueRange(0, 0, 999999999, 999999999)
	case "MICRO_OF_SECOND":
		return createValueRange(0, 0, 999999, 999999)
	case "MILLI_OF_SECOND":
		return createValueRange(0, 0, 999, 999)
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported field: "+fldName)
	}
}

func instantGet(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.get: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.get: field is null")
	}
	self := params[0].(*object.Object)
	fldName := getTemporalFieldName(params[1].(*object.Object))
	_, nano := getInstantSecondsAndNanos(self)

	switch strings.ToUpper(fldName) {
	case "NANO_OF_SECOND":
		return int64(nano)
	case "MICRO_OF_SECOND":
		return int64(nano / 1000)
	case "MILLI_OF_SECOND":
		return int64(nano / 1000000)
	case "INSTANT_SECONDS":
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Invalid field INSTANT_SECONDS for get() method, use getLong() instead")
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported field: "+fldName)
	}
}

func instantGetLong(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.getLong: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.getLong: field is null")
	}
	self := params[0].(*object.Object)
	fldName := getTemporalFieldName(params[1].(*object.Object))
	sec, nano := getInstantSecondsAndNanos(self)

	switch strings.ToUpper(fldName) {
	case "INSTANT_SECONDS":
		return sec
	case "NANO_OF_SECOND":
		return int64(nano)
	case "MICRO_OF_SECOND":
		return int64(nano / 1000)
	case "MILLI_OF_SECOND":
		return int64(nano / 1000000)
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported field: "+fldName)
	}
}

func instantGetEpochSecond(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.getEpochSecond: self is null")
	}
	sec, _ := getInstantSecondsAndNanos(params[0].(*object.Object))
	return sec
}

func instantGetNano(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.getNano: self is null")
	}
	_, nano := getInstantSecondsAndNanos(params[0].(*object.Object))
	return int64(nano)
}

func instantWithAdjuster(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.with: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.with: adjuster is null")
	}
	self := params[0].(*object.Object)
	adjuster := params[1].(*object.Object)

	kName := object.GoStringFromStringPoolIndex(adjuster.KlassName)
	if kName == classNameInstantVal {
		return adjuster
	}

	// If adjuster has adjustInto
	adjustMeth, ok := ghelpers.MethodSignatures[kName+".adjustInto(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"]
	if ok {
		return adjustMeth.GFunction([]any{adjuster, self})
	}

	return ghelpers.GetGErrBlk(excNames.DateTimeException, "Instant.with: unsupported adjuster "+kName)
}

func instantWithField(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.with: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.with: field is null")
	}
	newValue, ok := getInt64Val(params[2])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Instant.with: invalid newValue")
	}
	self := params[0].(*object.Object)
	fldName := getTemporalFieldName(params[1].(*object.Object))
	sec, nano := getInstantSecondsAndNanos(self)

	switch strings.ToUpper(fldName) {
	case "INSTANT_SECONDS":
		if newValue < instantMinSecond || newValue > instantMaxSecond {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for InstantSeconds: %d", newValue))
		}
		return createInstant(newValue, nano)
	case "NANO_OF_SECOND":
		if newValue < 0 || newValue > 999999999 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for NanoOfSecond: %d", newValue))
		}
		return createInstant(sec, newValue)
	case "MICRO_OF_SECOND":
		if newValue < 0 || newValue > 999999 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for MicroOfSecond: %d", newValue))
		}
		newNano := newValue*1000 + (nano % 1000)
		return createInstant(sec, newNano)
	case "MILLI_OF_SECOND":
		if newValue < 0 || newValue > 999 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for MilliOfSecond: %d", newValue))
		}
		newNano := newValue*1000000 + (nano % 1000000)
		return createInstant(sec, newNano)
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported field: "+fldName)
	}
}

func instantTruncatedTo(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.truncatedTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.truncatedTo: unit is null")
	}
	self := params[0].(*object.Object)
	unitName := getTemporalUnitName(params[1].(*object.Object))
	sec, nano := getInstantSecondsAndNanos(self)

	switch strings.ToUpper(unitName) {
	case "NANOS":
		return self
	case "MICROS":
		return createInstant(sec, (nano/1000)*1000)
	case "MILLIS":
		return createInstant(sec, (nano/1000000)*1000000)
	case "SECONDS":
		return createInstant(sec, 0)
	case "MINUTES":
		return createInstant(sec-floorMod(sec, 60), 0)
	case "HOURS":
		return createInstant(sec-floorMod(sec, 3600), 0)
	case "HALF_DAYS":
		return createInstant(sec-floorMod(sec, 43200), 0)
	case "DAYS":
		return createInstant(sec-floorMod(sec, 86400), 0)
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unit is not supported for truncation: "+unitName)
	}
}

func instantPlusAmount(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.plus: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.plus: amountToAdd is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)

	kName := object.GoStringFromStringPoolIndex(amount.KlassName)
	if kName == "java/time/Duration" {
		secToAdd, nanoToAdd := getInstantSecondsAndNanos(amount)
		return plusSecondsAndNanos(self, secToAdd, nanoToAdd)
	}

	// Try addTo method
	addToMeth, ok := ghelpers.MethodSignatures[kName+".addTo(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"]
	if ok {
		return addToMeth.GFunction([]any{amount, self})
	}

	return ghelpers.GetGErrBlk(excNames.DateTimeException, "Instant.plus: unsupported amount type "+kName)
}

func instantPlusUnit(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.plus: self is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.plus: unit is null")
	}
	amountToAdd, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Instant.plus: invalid amountToAdd")
	}
	if amountToAdd == 0 {
		return params[0]
	}
	self := params[0].(*object.Object)
	unitName := getTemporalUnitName(params[2].(*object.Object))

	switch strings.ToUpper(unitName) {
	case "NANOS":
		return plusSecondsAndNanos(self, 0, amountToAdd)
	case "MICROS":
		s := amountToAdd / 1000000
		n := (amountToAdd % 1000000) * 1000
		return plusSecondsAndNanos(self, s, n)
	case "MILLIS":
		s := amountToAdd / 1000
		n := (amountToAdd % 1000) * 1000000
		return plusSecondsAndNanos(self, s, n)
	case "SECONDS":
		return plusSecondsAndNanos(self, amountToAdd, 0)
	case "MINUTES":
		s, okM := multiplyExactSafe(amountToAdd, 60)
		if !okM {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "Instant.plus: overflow")
		}
		return plusSecondsAndNanos(self, s, 0)
	case "HOURS":
		s, okM := multiplyExactSafe(amountToAdd, 3600)
		if !okM {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "Instant.plus: overflow")
		}
		return plusSecondsAndNanos(self, s, 0)
	case "HALF_DAYS":
		s, okM := multiplyExactSafe(amountToAdd, 43200)
		if !okM {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "Instant.plus: overflow")
		}
		return plusSecondsAndNanos(self, s, 0)
	case "DAYS":
		s, okM := multiplyExactSafe(amountToAdd, 86400)
		if !okM {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "Instant.plus: overflow")
		}
		return plusSecondsAndNanos(self, s, 0)
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported unit: "+unitName)
	}
}

func plusSecondsAndNanos(self *object.Object, secondsToAdd, nanosToAdd int64) any {
	if secondsToAdd == 0 && nanosToAdd == 0 {
		return self
	}
	sec, nano := getInstantSecondsAndNanos(self)

	s1, ok1 := addExactSafe(sec, secondsToAdd)
	if !ok1 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Instant calculation exceeds supported range")
	}

	extraSec := floorDiv(nanosToAdd, instantNanosPerSec)
	remNano := floorMod(nanosToAdd, instantNanosPerSec)

	s2, ok2 := addExactSafe(s1, extraSec)
	if !ok2 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Instant calculation exceeds supported range")
	}

	totalNano := nano + remNano
	if totalNano >= instantNanosPerSec {
		s2, ok2 = addExactSafe(s2, 1)
		if !ok2 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "Instant calculation exceeds supported range")
		}
		totalNano -= instantNanosPerSec
	} else if totalNano < 0 {
		s2, ok2 = addExactSafe(s2, -1)
		if !ok2 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "Instant calculation exceeds supported range")
		}
		totalNano += instantNanosPerSec
	}

	return createInstantValidated(s2, totalNano)
}

func instantPlusSeconds(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.plusSeconds: self is null")
	}
	secondsToAdd, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Instant.plusSeconds: invalid secondsToAdd")
	}
	return plusSecondsAndNanos(params[0].(*object.Object), secondsToAdd, 0)
}

func instantPlusMillis(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.plusMillis: self is null")
	}
	millisToAdd, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Instant.plusMillis: invalid millisToAdd")
	}
	s := millisToAdd / 1000
	n := (millisToAdd % 1000) * 1000000
	return plusSecondsAndNanos(params[0].(*object.Object), s, n)
}

func instantPlusNanos(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.plusNanos: self is null")
	}
	nanosToAdd, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Instant.plusNanos: invalid nanosToAdd")
	}
	return plusSecondsAndNanos(params[0].(*object.Object), 0, nanosToAdd)
}

func instantMinusAmount(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.minus: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.minus: amountToSubtract is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)

	kName := object.GoStringFromStringPoolIndex(amount.KlassName)
	if kName == "java/time/Duration" {
		secToSub, nanoToSub := getInstantSecondsAndNanos(amount)
		if secToSub == math.MinInt64 {
			res := plusSecondsAndNanos(self, -(secToSub + 1), -nanoToSub)
			if err, isErr := res.(*ghelpers.GErrBlk); isErr {
				return err
			}
			return plusSecondsAndNanos(res.(*object.Object), 1, 0)
		}
		return plusSecondsAndNanos(self, -secToSub, -nanoToSub)
	}

	// Try subtractFrom
	subMeth, ok := ghelpers.MethodSignatures[kName+".subtractFrom(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"]
	if ok {
		return subMeth.GFunction([]any{amount, self})
	}

	return ghelpers.GetGErrBlk(excNames.DateTimeException, "Instant.minus: unsupported amount type "+kName)
}

func instantMinusUnit(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.minus: self is null")
	}
	amountToSubtract, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Instant.minus: invalid amountToSubtract")
	}
	if amountToSubtract == math.MinInt64 {
		res := instantPlusUnit([]any{params[0], -(amountToSubtract + 1), params[2]})
		if err, isErr := res.(*ghelpers.GErrBlk); isErr {
			return err
		}
		return instantPlusUnit([]any{res, int64(1), params[2]})
	}
	return instantPlusUnit([]any{params[0], -amountToSubtract, params[2]})
}

func instantMinusSeconds(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.minusSeconds: self is null")
	}
	secondsToSubtract, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Instant.minusSeconds: invalid secondsToSubtract")
	}
	if secondsToSubtract == math.MinInt64 {
		res := plusSecondsAndNanos(params[0].(*object.Object), -(secondsToSubtract + 1), 0)
		if err, isErr := res.(*ghelpers.GErrBlk); isErr {
			return err
		}
		return plusSecondsAndNanos(res.(*object.Object), 1, 0)
	}
	return plusSecondsAndNanos(params[0].(*object.Object), -secondsToSubtract, 0)
}

func instantMinusMillis(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.minusMillis: self is null")
	}
	millisToSubtract, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Instant.minusMillis: invalid millisToSubtract")
	}
	if millisToSubtract == math.MinInt64 {
		s := -(millisToSubtract + 1) / 1000
		n := (-(millisToSubtract + 1) % 1000) * 1000000
		res := plusSecondsAndNanos(params[0].(*object.Object), s, n)
		if err, isErr := res.(*ghelpers.GErrBlk); isErr {
			return err
		}
		return plusSecondsAndNanos(res.(*object.Object), 0, 1000000)
	}
	s := -millisToSubtract / 1000
	n := (-millisToSubtract % 1000) * 1000000
	return plusSecondsAndNanos(params[0].(*object.Object), s, n)
}

func instantMinusNanos(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.minusNanos: self is null")
	}
	nanosToSubtract, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Instant.minusNanos: invalid nanosToSubtract")
	}
	if nanosToSubtract == math.MinInt64 {
		res := plusSecondsAndNanos(params[0].(*object.Object), 0, -(nanosToSubtract + 1))
		if err, isErr := res.(*ghelpers.GErrBlk); isErr {
			return err
		}
		return plusSecondsAndNanos(res.(*object.Object), 0, 1)
	}
	return plusSecondsAndNanos(params[0].(*object.Object), 0, -nanosToSubtract)
}

func instantQuery(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.query: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.query: query is null")
	}
	self := params[0].(*object.Object)
	query := params[1].(*object.Object)

	queryName := ""
	if query.FieldTable != nil {
		if nameFld, okName := query.FieldTable["name"]; okName {
			if strObj, okStr := nameFld.Fvalue.(*object.Object); okStr && object.IsStringObject(strObj) {
				queryName = object.GoStringFromStringObject(strObj)
			}
		}
	}

	switch queryName {
	case "precision":
		return ghelpers.CreateTemporalUnit("NANOS")
	case "chronology", "zoneId", "zone", "offset", "localDate", "localTime":
		return nil
	default:
		return executeTemporalQuery(query, self)
	}
}

func instantAdjustInto(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.adjustInto: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.adjustInto: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)
	sec, nano := getInstantSecondsAndNanos(self)

	cfInstantSec := getChronoFieldForInstant("INSTANT_SECONDS")
	cfNanoOfSec := getChronoFieldForInstant("NANO_OF_SECOND")

	res1 := temporalWithField([]any{temporal, cfInstantSec, sec})
	if err, isErr := res1.(*ghelpers.GErrBlk); isErr {
		return err
	}
	res2 := temporalWithField([]any{res1, cfNanoOfSec, nano})
	return res2
}

func instantUntil(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.until: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.until: endExclusive is null")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.until: unit is null")
	}
	self := params[0].(*object.Object)
	endRes := instantFrom([]any{params[1]})
	if err, isErr := endRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	end := endRes.(*object.Object)

	unitName := getTemporalUnitName(params[2].(*object.Object))
	startSec, startNano := getInstantSecondsAndNanos(self)
	endSec, endNano := getInstantSecondsAndNanos(end)

	secsDiff, ok := addExactSafe(endSec, -startSec)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Instant.until: arithmetic overflow")
	}
	nanosDiff := endNano - startNano

	switch strings.ToUpper(unitName) {
	case "NANOS":
		totalNanos, okN := multiplyExactSafe(secsDiff, instantNanosPerSec)
		if !okN {
			return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Instant.until: arithmetic overflow")
		}
		resN, okN2 := addExactSafe(totalNanos, nanosDiff)
		if !okN2 {
			return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Instant.until: arithmetic overflow")
		}
		return resN
	case "MICROS":
		totalNanos, okN := multiplyExactSafe(secsDiff, instantNanosPerSec)
		if !okN {
			return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Instant.until: arithmetic overflow")
		}
		resN, okN2 := addExactSafe(totalNanos, nanosDiff)
		if !okN2 {
			return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Instant.until: arithmetic overflow")
		}
		return resN / 1000
	case "MILLIS":
		startMilli := instantToEpochMilliHelper(startSec, startNano)
		endMilli := instantToEpochMilliHelper(endSec, endNano)
		resM, okM := addExactSafe(endMilli, -startMilli)
		if !okM {
			return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Instant.until: arithmetic overflow")
		}
		return resM
	case "SECONDS":
		if secsDiff > 0 && nanosDiff < 0 {
			secsDiff--
		} else if secsDiff < 0 && nanosDiff > 0 {
			secsDiff++
		}
		return secsDiff
	case "MINUTES":
		if secsDiff > 0 && nanosDiff < 0 {
			secsDiff--
		} else if secsDiff < 0 && nanosDiff > 0 {
			secsDiff++
		}
		return secsDiff / 60
	case "HOURS":
		if secsDiff > 0 && nanosDiff < 0 {
			secsDiff--
		} else if secsDiff < 0 && nanosDiff > 0 {
			secsDiff++
		}
		return secsDiff / 3600
	case "HALF_DAYS":
		if secsDiff > 0 && nanosDiff < 0 {
			secsDiff--
		} else if secsDiff < 0 && nanosDiff > 0 {
			secsDiff++
		}
		return secsDiff / 43200
	case "DAYS":
		if secsDiff > 0 && nanosDiff < 0 {
			secsDiff--
		} else if secsDiff < 0 && nanosDiff > 0 {
			secsDiff++
		}
		return secsDiff / 86400
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported unit: "+unitName)
	}
}

func instantAtOffset(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.atOffset: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.atOffset: offset is null")
	}
	self := params[0].(*object.Object)
	offsetObj := params[1].(*object.Object)
	sec, nano := getInstantSecondsAndNanos(self)
	offsetSec := getZoneOffsetTotalSeconds(offsetObj)

	localSec := sec + offsetSec
	epochDay := floorDiv(localSec, secondsPerDay)
	secOfDay := floorMod(localSec, secondsPerDay)

	y, m, d := epochDayToYMD(epochDay)
	dateObj := CreateLocalDate(y, m, d)
	timeObj := ofNanoOfDay(secOfDay*nanosPerSecond + nano)
	ldt := createLocalDateTime(dateObj, timeObj)

	odt := object.MakeEmptyObjectWithClassName(&classNameOffsetDateTimeVal)
	odt.FieldTable["dateTime"] = object.Field{Ftype: types.Ref, Fvalue: ldt}
	odt.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: offsetObj}
	return odt
}

func instantAtZone(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.atZone: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.atZone: zone is null")
	}
	self := params[0].(*object.Object)
	zoneObj := params[1].(*object.Object)

	rulesRes := zoneIdGetRules([]any{zoneObj})
	if err, isErr := rulesRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	rulesObj := rulesRes.(*object.Object)
	offsetRes := zoneRulesGetOffsetInstant([]any{rulesObj, self})
	if err, isErr := offsetRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	offsetObj := offsetRes.(*object.Object)

	sec, nano := getInstantSecondsAndNanos(self)
	offsetSec := getZoneOffsetTotalSeconds(offsetObj)

	localSec := sec + offsetSec
	epochDay := floorDiv(localSec, secondsPerDay)
	secOfDay := floorMod(localSec, secondsPerDay)

	y, m, d := epochDayToYMD(epochDay)
	dateObj := CreateLocalDate(y, m, d)
	timeObj := ofNanoOfDay(secOfDay*nanosPerSecond + nano)
	ldt := createLocalDateTime(dateObj, timeObj)

	zdt := object.MakeEmptyObjectWithClassName(&classNameZonedDateTimeVal)
	zdt.FieldTable["dateTime"] = object.Field{Ftype: types.Ref, Fvalue: ldt}
	zdt.FieldTable["zone"] = object.Field{Ftype: types.Ref, Fvalue: zoneObj}
	zdt.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: offsetObj}
	return zdt
}

func instantToEpochMilliHelper(seconds, nanos int64) int64 {
	if seconds >= 0 || nanos == 0 {
		return seconds*1000 + nanos/1000000
	}
	return (seconds+1)*1000 - 1000 + nanos/1000000
}

func instantToEpochMilli(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.toEpochMilli: self is null")
	}
	sec, nano := getInstantSecondsAndNanos(params[0].(*object.Object))
	if sec >= 0 || nano == 0 {
		milli, ok := multiplyExactSafe(sec, 1000)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Instant.toEpochMilli: overflow")
		}
		res, ok2 := addExactSafe(milli, nano/1000000)
		if !ok2 {
			return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Instant.toEpochMilli: overflow")
		}
		return res
	}

	milli, ok := multiplyExactSafe(sec+1, 1000)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Instant.toEpochMilli: overflow")
	}
	res, ok2 := addExactSafe(milli, (nano/1000000)-1000)
	if !ok2 {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Instant.toEpochMilli: overflow")
	}
	return res
}

func instantCompareTo(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.compareTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.compareTo: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)
	kName := object.GoStringFromStringPoolIndex(other.KlassName)
	if kName != classNameInstantVal {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "Instant.compareTo: cannot compare to "+kName)
	}

	s1, n1 := getInstantSecondsAndNanos(self)
	s2, n2 := getInstantSecondsAndNanos(other)

	if s1 < s2 {
		return int64(-1)
	}
	if s1 > s2 {
		return int64(1)
	}
	if n1 < n2 {
		return int64(-1)
	}
	if n1 > n2 {
		return int64(1)
	}
	return int64(0)
}

func instantIsAfter(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.isAfter: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.isAfter: other is null")
	}
	cmpRes := instantCompareTo(params)
	if err, isErr := cmpRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return types.ConvertGoBoolToJavaBool(cmpRes.(int64) > 0)
}

func instantIsBefore(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.isBefore: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.isBefore: other is null")
	}
	cmpRes := instantCompareTo(params)
	if err, isErr := cmpRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	return types.ConvertGoBoolToJavaBool(cmpRes.(int64) < 0)
}

func instantEquals(params []any) any {
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
	kName := object.GoStringFromStringPoolIndex(other.KlassName)
	if kName != classNameInstantVal {
		return types.JavaBoolFalse
	}
	s1, n1 := getInstantSecondsAndNanos(self)
	s2, n2 := getInstantSecondsAndNanos(other)
	return types.ConvertGoBoolToJavaBool(s1 == s2 && n1 == n2)
}

func instantHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.hashCode: self is null")
	}
	sec, nano := getInstantSecondsAndNanos(params[0].(*object.Object))
	uSec := uint64(sec)
	h := int32(uSec ^ (uSec >> 32)) + int32(51*nano)
	return int64(h)
}

func instantToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Instant.toString: self is null")
	}
	self := params[0].(*object.Object)
	sec, nano := getInstantSecondsAndNanos(self)

	epochDay := floorDiv(sec, 86400)
	secOfDay := floorMod(sec, 86400)

	y, m, d := epochDayToYMD(epochDay)
	h := secOfDay / 3600
	min := (secOfDay % 3600) / 60
	s := secOfDay % 60

	var yearStr string
	if y >= 0 && y <= 9999 {
		yearStr = fmt.Sprintf("%04d", y)
	} else if y > 9999 {
		yearStr = fmt.Sprintf("+%d", y)
	} else {
		yearStr = fmt.Sprintf("%05d", y) // includes negative sign
	}

	var nanoStr string
	if nano > 0 {
		if nano%1000000 == 0 {
			nanoStr = fmt.Sprintf(".%03d", nano/1000000)
		} else if nano%1000 == 0 {
			nanoStr = fmt.Sprintf(".%06d", nano/1000)
		} else {
			nanoStr = fmt.Sprintf(".%09d", nano)
		}
	}

	res := fmt.Sprintf("%s-%02d-%02dT%02d:%02d:%02d%sZ", yearStr, m, d, h, min, s, nanoStr)
	return object.StringObjectFromGoString(res)
}
