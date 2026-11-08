/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
)

// Load_Time_Temporal_TemporalAdjuster loads method signatures for java/time/temporal/TemporalAdjuster.
func Load_Time_Temporal_TemporalAdjuster() {
	ghelpers.MethodSignatures["java/time/temporal/TemporalAdjuster.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAdjuster.adjustInto(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalAdjusterAdjustInto,
		}
}

// temporalAdjusterAdjustInto adjusts the specified temporal object.
// Temporal adjustInto(Temporal temporal)
func temporalAdjusterAdjustInto(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAdjuster.adjustInto: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAdjuster.adjustInto: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)

	return applyTemporalAdjuster(self, temporal)
}

func getAdjusterType(adjuster *object.Object) string {
	if adjuster == nil || adjuster.FieldTable == nil {
		return ""
	}
	if typeFld, ok := adjuster.FieldTable["adjusterType"]; ok {
		if strObj, okStr := typeFld.Fvalue.(*object.Object); okStr && object.IsStringObject(strObj) {
			return object.GoStringFromStringObject(strObj)
		}
	}
	return ""
}

func getAdjusterDowValue(adjuster *object.Object) (int64, bool) {
	if adjuster == nil || adjuster.FieldTable == nil {
		return 0, false
	}
	if dowFld, ok := adjuster.FieldTable["dowValue"]; ok {
		if val, okVal := getInt64Val(dowFld.Fvalue); okVal {
			return val, true
		}
	}
	if dowObjFld, ok := adjuster.FieldTable["dayOfWeek"]; ok {
		if dowObj, okObj := dowObjFld.Fvalue.(*object.Object); okObj && !object.IsNull(dowObj) {
			return getDayOfWeekValueFromObject(dowObj)
		}
	}
	return 0, false
}

func getAdjusterOrdinal(adjuster *object.Object) (int64, bool) {
	if adjuster == nil || adjuster.FieldTable == nil {
		return 0, false
	}
	if ordFld, ok := adjuster.FieldTable["ordinal"]; ok {
		return getInt64Val(ordFld.Fvalue)
	}
	return 0, false
}

// applyTemporalAdjuster applies an adjuster object to a temporal target.
func applyTemporalAdjuster(adjuster *object.Object, temporal *object.Object) any {
	if adjuster == nil || object.IsNull(adjuster) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "applyTemporalAdjuster: adjuster is null")
	}
	if temporal == nil || object.IsNull(temporal) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "applyTemporalAdjuster: temporal is null")
	}

	adjusterClass := object.GoStringFromStringPoolIndex(adjuster.KlassName)
	switch adjusterClass {
	case classNameLocalDateVal:
		return localDateAdjustInto([]any{adjuster, temporal})
	case classNameLocalTime:
		return localTimeAdjustInto([]any{adjuster, temporal})
	case classNameLocalDateTime:
		return localDateTimeAdjustInto([]any{adjuster, temporal})
	case classNameDayOfWeekVal:
		return dayOfWeekAdjustInto([]any{adjuster, temporal})
	case classNameMonthVal:
		return monthAdjustInto([]any{adjuster, temporal})
	case classNameYearVal:
		return yearAdjustInto([]any{adjuster, temporal})
	case classNameChronoLocalDate:
		return chronoLocalDateAdjustInto([]any{adjuster, temporal})
	case classNameChronoLocalDateTime:
		return chronoLocalDateTimeAdjustInto([]any{adjuster, temporal})
	}

	adjType := getAdjusterType(adjuster)
	if adjType == "" {
		return cloneTemporalObject(temporal)
	}

	temporalClass := object.GoStringFromStringPoolIndex(temporal.KlassName)
	if temporalClass == classNameLocalTime {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unsupported field: DayOfMonth")
	}

	if temporalClass == classNameLocalDateTime {
		dObj := getLocalDateTimeDate(temporal)
		tObj := getLocalDateTimeTime(temporal)
		if dObj == nil || tObj == nil {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "applyTemporalAdjuster: invalid LocalDateTime")
		}
		newDate := applyTemporalAdjuster(adjuster, dObj)
		if err, isErr := newDate.(*ghelpers.GErrBlk); isErr {
			return err
		}
		return createLocalDateTime(newDate.(*object.Object), tObj)
	}

	if temporalClass == classNameChronoLocalDateTime {
		dObj := getChronoLocalDateTimeDate(temporal)
		tObj := getChronoLocalDateTimeTime(temporal)
		if dObj == nil || tObj == nil {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "applyTemporalAdjuster: invalid ChronoLocalDateTime")
		}
		newDate := applyTemporalAdjuster(adjuster, dObj)
		if err, isErr := newDate.(*ghelpers.GErrBlk); isErr {
			return err
		}
		return createChronoLocalDateTime(newDate.(*object.Object), tObj)
	}

	if temporalClass == classNameChronoLocalDate {
		y, m, d := ghelpers.GetChronoLocalDateYMD(temporal)
		chrono := GetChronoLocalDateChronology(temporal)
		fakeDate := CreateLocalDate(y, m, d)
		newDate := applyTemporalAdjuster(adjuster, fakeDate)
		if err, isErr := newDate.(*ghelpers.GErrBlk); isErr {
			return err
		}
		ny, nm, nd, _ := getLocalDateYMD(newDate.(*object.Object))
		return createChronoLocalDate(chrono, ny, nm, nd)
	}

	if temporalClass == classNameChronoZonedDateTime {
		ldt := getChronoZonedDateTimeLocalDateTime(temporal)
		zObj := getChronoZonedDateTimeZone(temporal)
		offObj := getChronoZonedDateTimeOffset(temporal)
		if ldt == nil {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "applyTemporalAdjuster: invalid ChronoZonedDateTime")
		}
		newDT := applyTemporalAdjuster(adjuster, ldt)
		if err, isErr := newDT.(*ghelpers.GErrBlk); isErr {
			return err
		}
		return createChronoZonedDateTimeFromLDT(newDT.(*object.Object), zObj, offObj)
	}

	y, m, d, okYMD := getLocalDateYMD(temporal)
	if !okYMD {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Temporal type %s does not support date adjustment", temporalClass))
	}

	var ny, nm, nd int64
	switch adjType {
	case "firstDayOfMonth":
		ny, nm, nd = y, m, 1

	case "lastDayOfMonth":
		ny, nm, nd = y, m, daysInMonthDate(y, m)

	case "firstDayOfNextMonth":
		nm = m + 1
		ny = y
		if nm > 12 {
			nm = 1
			ny++
		}
		nd = 1

	case "firstDayOfYear":
		ny, nm, nd = y, 1, 1

	case "lastDayOfYear":
		ny, nm, nd = y, 12, 31

	case "firstDayOfNextYear":
		ny, nm, nd = y+1, 1, 1

	case "dayOfWeekInMonth", "firstInMonth", "lastInMonth":
		ordinal, _ := getAdjusterOrdinal(adjuster)
		if adjType == "firstInMonth" {
			ordinal = 1
		} else if adjType == "lastInMonth" {
			ordinal = -1
		}
		dowVal, okDow := getAdjusterDowValue(adjuster)
		if !okDow {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "dayOfWeekInMonth: missing day of week")
		}

		if ordinal >= 0 {
			curDow := dayOfWeekDate(y, m, 1)
			dowDiff := (dowVal - curDow + 7) % 7
			dowDiff += (ordinal - 1) * 7
			ep := ymdToEpochDay(y, m, 1) + dowDiff
			ny, nm, nd = epochDayToYMD(ep)
		} else {
			lastDay := daysInMonthDate(y, m)
			curDow := dayOfWeekDate(y, m, lastDay)
			daysDiff := dowVal - curDow
			if daysDiff > 0 {
				daysDiff -= 7
			}
			daysDiff -= (-ordinal - 1) * 7
			ep := ymdToEpochDay(y, m, lastDay) + daysDiff
			ny, nm, nd = epochDayToYMD(ep)
		}

	case "next":
		dowVal, okDow := getAdjusterDowValue(adjuster)
		if !okDow {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "next: missing day of week")
		}
		curDow := dayOfWeekDate(y, m, d)
		daysDiff := curDow - dowVal
		var daysToAdd int64
		if daysDiff >= 0 {
			daysToAdd = 7 - daysDiff
		} else {
			daysToAdd = -daysDiff
		}
		ep := ymdToEpochDay(y, m, d) + daysToAdd
		ny, nm, nd = epochDayToYMD(ep)

	case "nextOrSame":
		dowVal, okDow := getAdjusterDowValue(adjuster)
		if !okDow {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "nextOrSame: missing day of week")
		}
		curDow := dayOfWeekDate(y, m, d)
		if curDow == dowVal {
			return temporal
		}
		daysDiff := curDow - dowVal
		var daysToAdd int64
		if daysDiff >= 0 {
			daysToAdd = 7 - daysDiff
		} else {
			daysToAdd = -daysDiff
		}
		ep := ymdToEpochDay(y, m, d) + daysToAdd
		ny, nm, nd = epochDayToYMD(ep)

	case "previous":
		dowVal, okDow := getAdjusterDowValue(adjuster)
		if !okDow {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "previous: missing day of week")
		}
		curDow := dayOfWeekDate(y, m, d)
		daysDiff := dowVal - curDow
		var daysToSub int64
		if daysDiff >= 0 {
			daysToSub = 7 - daysDiff
		} else {
			daysToSub = -daysDiff
		}
		ep := ymdToEpochDay(y, m, d) - daysToSub
		ny, nm, nd = epochDayToYMD(ep)

	case "previousOrSame":
		dowVal, okDow := getAdjusterDowValue(adjuster)
		if !okDow {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "previousOrSame: missing day of week")
		}
		curDow := dayOfWeekDate(y, m, d)
		if curDow == dowVal {
			return temporal
		}
		daysDiff := dowVal - curDow
		var daysToSub int64
		if daysDiff >= 0 {
			daysToSub = 7 - daysDiff
		} else {
			daysToSub = -daysDiff
		}
		ep := ymdToEpochDay(y, m, d) - daysToSub
		ny, nm, nd = epochDayToYMD(ep)

	case "ofDateAdjuster":
		if fnFld, okFn := adjuster.FieldTable["dateBasedAdjuster"]; okFn {
			if fnObj, okObj := fnFld.Fvalue.(*object.Object); okObj && !object.IsNull(fnObj) {
				if object.GoStringFromStringPoolIndex(fnObj.KlassName) == classNameLocalDateVal {
					ny, nm, nd, _ = getLocalDateYMD(fnObj)
				} else {
					return localDateAdjustInto([]any{fnObj, temporal})
				}
			} else {
				return ghelpers.GetGErrBlk(excNames.NullPointerException, "ofDateAdjuster: dateBasedAdjuster is null")
			}
		} else {
			return ghelpers.GetGErrBlk(excNames.NullPointerException, "ofDateAdjuster: dateBasedAdjuster missing")
		}

	default:
		return cloneTemporalObject(temporal)
	}

	if temporalClass == classNameLocalDateVal {
		return CreateLocalDate(ny, nm, nd)
	}

	cfEpochDay := getChronoField("EPOCH_DAY")
	return temporalWithField([]any{temporal, cfEpochDay, ymdToEpochDay(ny, nm, nd)})
}
