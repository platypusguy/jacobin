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
	"strings"
)

// Load_Time_Chrono_ChronoPeriod loads method signatures for java/time/chrono/ChronoPeriod.
func Load_Time_Chrono_ChronoPeriod() {
	ghelpers.MethodSignatures["java/time/chrono/ChronoPeriod.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoPeriod.between(Ljava/time/chrono/ChronoLocalDate;Ljava/time/chrono/ChronoLocalDate;)Ljava/time/chrono/ChronoPeriod;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronoPeriodBetween,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoPeriod.addTo(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoPeriodAddTo,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoPeriod.subtractFrom(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoPeriodSubtractFrom,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoPeriod.get(Ljava/time/temporal/TemporalUnit;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoPeriodGet,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoPeriod.getChronology()Ljava/time/chrono/Chronology;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoPeriodGetChronology,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoPeriod.getUnits()Ljava/util/List;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoPeriodGetUnits,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoPeriod.isNegative()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoPeriodIsNegative,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoPeriod.isZero()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoPeriodIsZero,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoPeriod.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/chrono/ChronoPeriod;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoPeriodMinus,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoPeriod.multipliedBy(I)Ljava/time/chrono/ChronoPeriod;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoPeriodMultipliedBy,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoPeriod.negated()Ljava/time/chrono/ChronoPeriod;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoPeriodNegated,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoPeriod.normalized()Ljava/time/chrono/ChronoPeriod;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoPeriodNormalized,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoPeriod.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/chrono/ChronoPeriod;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoPeriodPlus,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoPeriod.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronoPeriodEquals,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoPeriod.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoPeriodHashCode,
		}

	ghelpers.MethodSignatures["java/time/chrono/ChronoPeriod.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronoPeriodToString,
		}
}

// Helpers for ChronoPeriod

func getChronoPeriodValues(periodObject *object.Object) (yearsValue, monthsValue, daysValue int64) {
	if periodObject == nil || periodObject.FieldTable == nil {
		return 0, 0, 0
	}
	if yearsField, exists := periodObject.FieldTable["years"]; exists {
		yearsValue, _ = getInt64Val(yearsField.Fvalue)
	}
	if monthsField, exists := periodObject.FieldTable["months"]; exists {
		monthsValue, _ = getInt64Val(monthsField.Fvalue)
	}
	if daysField, exists := periodObject.FieldTable["days"]; exists {
		daysValue, _ = getInt64Val(daysField.Fvalue)
	}
	return yearsValue, monthsValue, daysValue
}

func getChronoPeriodChronologyObject(periodObject *object.Object) *object.Object {
	if periodObject == nil || periodObject.FieldTable == nil {
		return GetIsoChronology()
	}
	if chronologyField, exists := periodObject.FieldTable["chronology"]; exists {
		if chronologyObject, isObject := chronologyField.Fvalue.(*object.Object); isObject && chronologyObject != nil {
			return chronologyObject
		}
	}
	return GetIsoChronology()
}

// chronoPeriodBetween obtains a ChronoPeriod consisting of amount of time between two dates.
// static ChronoPeriod between(ChronoLocalDate startDateInclusive, ChronoLocalDate endDateExclusive)
func chronoPeriodBetween(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.between: startDateInclusive is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.between: endDateExclusive is null")
	}
	startDate := params[0].(*object.Object)
	endDate := params[1].(*object.Object)

	return chronoLocalDateUntilPeriod([]any{startDate, endDate})
}

// chronoPeriodAddTo adds this period to the specified temporal object.
// Temporal addTo(Temporal temporal)
func chronoPeriodAddTo(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.addTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.addTo: temporal is null")
	}
	self := params[0].(*object.Object)
	temporalTarget := params[1].(*object.Object)

	yearsValue, monthsValue, daysValue := getChronoPeriodValues(self)
	periodChronology := getChronoPeriodChronologyObject(self)
	temporalChronology := GetChronoLocalDateChronology(temporalTarget)
	if temporalChronology != nil && periodChronology != nil {
		periodId := getChronologyId(periodChronology)
		temporalId := getChronologyId(temporalChronology)
		if periodId != "" && temporalId != "" && periodId != temporalId {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Chronology mismatch, expected: %s, actual: %s", periodId, temporalId))
		}
	}

	resultTemporal := temporalTarget
	if yearsValue != 0 {
		plusResult := temporalPlusUnit([]any{resultTemporal, yearsValue, ghelpers.CreateTemporalUnit("YEARS")})
		if resultingObject, isObject := plusResult.(*object.Object); isObject {
			resultTemporal = resultingObject
		} else {
			return plusResult
		}
	}
	if monthsValue != 0 {
		plusResult := temporalPlusUnit([]any{resultTemporal, monthsValue, ghelpers.CreateTemporalUnit("MONTHS")})
		if resultingObject, isObject := plusResult.(*object.Object); isObject {
			resultTemporal = resultingObject
		} else {
			return plusResult
		}
	}
	if daysValue != 0 {
		plusResult := temporalPlusUnit([]any{resultTemporal, daysValue, ghelpers.CreateTemporalUnit("DAYS")})
		if resultingObject, isObject := plusResult.(*object.Object); isObject {
			resultTemporal = resultingObject
		} else {
			return plusResult
		}
	}
	return resultTemporal
}

// chronoPeriodSubtractFrom subtracts this period from the specified temporal object.
// Temporal subtractFrom(Temporal temporal)
func chronoPeriodSubtractFrom(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.subtractFrom: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.subtractFrom: temporal is null")
	}
	self := params[0].(*object.Object)
	temporalTarget := params[1].(*object.Object)

	yearsValue, monthsValue, daysValue := getChronoPeriodValues(self)
	periodChronology := getChronoPeriodChronologyObject(self)
	temporalChronology := GetChronoLocalDateChronology(temporalTarget)
	if temporalChronology != nil && periodChronology != nil {
		periodId := getChronologyId(periodChronology)
		temporalId := getChronologyId(temporalChronology)
		if periodId != "" && temporalId != "" && periodId != temporalId {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Chronology mismatch, expected: %s, actual: %s", periodId, temporalId))
		}
	}

	resultTemporal := temporalTarget
	if yearsValue != 0 {
		minusResult := temporalMinusUnit([]any{resultTemporal, yearsValue, ghelpers.CreateTemporalUnit("YEARS")})
		if resultingObject, isObject := minusResult.(*object.Object); isObject {
			resultTemporal = resultingObject
		} else {
			return minusResult
		}
	}
	if monthsValue != 0 {
		minusResult := temporalMinusUnit([]any{resultTemporal, monthsValue, ghelpers.CreateTemporalUnit("MONTHS")})
		if resultingObject, isObject := minusResult.(*object.Object); isObject {
			resultTemporal = resultingObject
		} else {
			return minusResult
		}
	}
	if daysValue != 0 {
		minusResult := temporalMinusUnit([]any{resultTemporal, daysValue, ghelpers.CreateTemporalUnit("DAYS")})
		if resultingObject, isObject := minusResult.(*object.Object); isObject {
			resultTemporal = resultingObject
		} else {
			return minusResult
		}
	}
	return resultTemporal
}

// chronoPeriodGet returns the value of the requested unit.
// long get(TemporalUnit unit)
func chronoPeriodGet(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.get: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.get: unit is null")
	}
	self := params[0].(*object.Object)
	unit := params[1].(*object.Object)

	unitName := strings.ToUpper(getTemporalUnitName(unit))
	yearsValue, monthsValue, daysValue := getChronoPeriodValues(self)

	switch unitName {
	case "YEARS":
		return yearsValue
	case "MONTHS":
		return monthsValue
	case "DAYS":
		return daysValue
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "ChronoPeriod.get: unsupported unit: "+unitName)
	}
}

// chronoPeriodGetChronology returns the chronology of this period.
// Chronology getChronology()
func chronoPeriodGetChronology(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.getChronology: self is null")
	}
	self := params[0].(*object.Object)
	return getChronoPeriodChronologyObject(self)
}

// chronoPeriodGetUnits returns the list of units supported by this period.
// List<TemporalUnit> getUnits()
func chronoPeriodGetUnits(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.getUnits: self is null")
	}
	listObject := object.MakeEmptyObjectWithClassName(&classNameArrayList)
	unitsSlice := []any{
		ghelpers.CreateTemporalUnit("YEARS"),
		ghelpers.CreateTemporalUnit("MONTHS"),
		ghelpers.CreateTemporalUnit("DAYS"),
	}
	listObject.FieldTable["value"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: unitsSlice,
	}
	return listObject
}

// chronoPeriodIsNegative checks if any of the supported units of this period are negative.
// boolean isNegative()
func chronoPeriodIsNegative(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.isNegative: self is null")
	}
	self := params[0].(*object.Object)
	yearsValue, monthsValue, daysValue := getChronoPeriodValues(self)
	if yearsValue < 0 || monthsValue < 0 || daysValue < 0 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// chronoPeriodIsZero checks if all the supported units of this period are zero.
// boolean isZero()
func chronoPeriodIsZero(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.isZero: self is null")
	}
	self := params[0].(*object.Object)
	yearsValue, monthsValue, daysValue := getChronoPeriodValues(self)
	if yearsValue == 0 && monthsValue == 0 && daysValue == 0 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// chronoPeriodMinus returns a copy of this period with the specified amount subtracted.
// ChronoPeriod minus(TemporalAmount amountToSubtract)
func chronoPeriodMinus(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.minus: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.minus: amountToSubtract is null")
	}
	self := params[0].(*object.Object)
	amountToSubtract := params[1].(*object.Object)

	selfChronology := getChronoPeriodChronologyObject(self)
	otherChronology := getChronoPeriodChronologyObject(amountToSubtract)
	if selfChronology != nil && otherChronology != nil {
		selfId := getChronologyId(selfChronology)
		otherId := getChronologyId(otherChronology)
		if selfId != "" && otherId != "" && selfId != otherId {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Chronology mismatch, expected: %s, actual: %s", selfId, otherId))
		}
	}

	yearsSelf, monthsSelf, daysSelf := getChronoPeriodValues(self)
	yearsOther, monthsOther, daysOther := getChronoPeriodValues(amountToSubtract)

	return ghelpers.CreateChronoPeriodYMD(selfChronology, yearsSelf-yearsOther, monthsSelf-monthsOther, daysSelf-daysOther)
}

// chronoPeriodMultipliedBy returns a new instance with each amount in this period multiplied by the scalar.
// ChronoPeriod multipliedBy(int scalar)
func chronoPeriodMultipliedBy(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.multipliedBy: self is null")
	}
	self := params[0].(*object.Object)
	scalarValue, _ := getInt64Val(params[1])
	if scalarValue == 1 {
		return self
	}

	yearsValue, monthsValue, daysValue := getChronoPeriodValues(self)
	chronologyObject := getChronoPeriodChronologyObject(self)

	return ghelpers.CreateChronoPeriodYMD(chronologyObject, yearsValue*scalarValue, monthsValue*scalarValue, daysValue*scalarValue)
}

// chronoPeriodNegated returns a new instance with each amount in this period negated.
// ChronoPeriod negated()
func chronoPeriodNegated(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.negated: self is null")
	}
	self := params[0].(*object.Object)
	return chronoPeriodMultipliedBy([]any{self, int64(-1)})
}

// chronoPeriodNormalized returns a copy of this period with the amounts of years and months normalized.
// ChronoPeriod normalized()
func chronoPeriodNormalized(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.normalized: self is null")
	}
	self := params[0].(*object.Object)
	yearsValue, monthsValue, daysValue := getChronoPeriodValues(self)
	chronologyObject := getChronoPeriodChronologyObject(self)

	totalMonths := yearsValue*12 + monthsValue
	normalizedYears := totalMonths / 12
	normalizedMonths := totalMonths % 12

	return ghelpers.CreateChronoPeriodYMD(chronologyObject, normalizedYears, normalizedMonths, daysValue)
}

// chronoPeriodPlus returns a copy of this period with the specified amount added.
// ChronoPeriod plus(TemporalAmount amountToAdd)
func chronoPeriodPlus(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.plus: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.plus: amountToAdd is null")
	}
	self := params[0].(*object.Object)
	amountToAdd := params[1].(*object.Object)

	selfChronology := getChronoPeriodChronologyObject(self)
	otherChronology := getChronoPeriodChronologyObject(amountToAdd)
	if selfChronology != nil && otherChronology != nil {
		selfId := getChronologyId(selfChronology)
		otherId := getChronologyId(otherChronology)
		if selfId != "" && otherId != "" && selfId != otherId {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Chronology mismatch, expected: %s, actual: %s", selfId, otherId))
		}
	}

	yearsSelf, monthsSelf, daysSelf := getChronoPeriodValues(self)
	yearsOther, monthsOther, daysOther := getChronoPeriodValues(amountToAdd)

	return ghelpers.CreateChronoPeriodYMD(selfChronology, yearsSelf+yearsOther, monthsSelf+monthsOther, daysSelf+daysOther)
}

// chronoPeriodEquals checks if this period is equal to another period.
// boolean equals(Object obj)
func chronoPeriodEquals(params []any) any {
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

	yearsSelf, monthsSelf, daysSelf := getChronoPeriodValues(self)
	yearsOther, monthsOther, daysOther := getChronoPeriodValues(other)

	if yearsSelf != yearsOther || monthsSelf != monthsOther || daysSelf != daysOther {
		return types.JavaBoolFalse
	}

	selfChronology := getChronoPeriodChronologyObject(self)
	otherChronology := getChronoPeriodChronologyObject(other)
	selfId := getChronologyId(selfChronology)
	otherId := getChronologyId(otherChronology)

	if selfId == otherId {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// chronoPeriodHashCode returns a hash code for this period.
// int hashCode()
func chronoPeriodHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return int64(0)
	}
	self := params[0].(*object.Object)
	yearsValue, monthsValue, daysValue := getChronoPeriodValues(self)
	chronologyObject := getChronoPeriodChronologyObject(self)
	chronologyId := getChronologyId(chronologyObject)

	var chronologyHashCode int64 = 0
	for _, character := range chronologyId {
		chronologyHashCode = 31*chronologyHashCode + int64(character)
	}

	combinedHash := (yearsValue ^ (monthsValue << 8) ^ (daysValue << 16)) ^ chronologyHashCode
	return combinedHash
}

// chronoPeriodToString returns a string representation of this period.
// String toString()
func chronoPeriodToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ChronoPeriod.toString: self is null")
	}
	self := params[0].(*object.Object)
	yearsValue, monthsValue, daysValue := getChronoPeriodValues(self)
	chronologyObject := getChronoPeriodChronologyObject(self)
	chronologyId := getChronologyId(chronologyObject)
	if chronologyId == "" {
		chronologyId = "ISO"
	}

	if yearsValue == 0 && monthsValue == 0 && daysValue == 0 {
		return object.StringObjectFromGoString(chronologyId + " P0D")
	}

	formattedPeriod := chronologyId + " P"
	if yearsValue != 0 {
		formattedPeriod += fmt.Sprintf("%dY", yearsValue)
	}
	if monthsValue != 0 {
		formattedPeriod += fmt.Sprintf("%dM", monthsValue)
	}
	if daysValue != 0 {
		formattedPeriod += fmt.Sprintf("%dD", daysValue)
	}

	return object.StringObjectFromGoString(formattedPeriod)
}
