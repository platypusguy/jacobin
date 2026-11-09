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
	"math/bits"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var classNamePeriod = "java/time/Period"

var periodMutex sync.Mutex
var periodInitialized bool = false

var periodZero *object.Object

func ensurePeriodInited() {
	periodMutex.Lock()
	defer periodMutex.Unlock()
	if periodInitialized {
		return
	}

	periodZero = createPeriodInternal(0, 0, 0)

	_ = statics.AddStatic(classNamePeriod+".ZERO", statics.Static{
		Type:  "Ljava/time/Period;",
		Value: periodZero,
	})

	periodInitialized = true
}

func createPeriodInternal(years, months, days int64) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNamePeriod)
	obj.FieldTable["years"] = object.Field{Ftype: types.Int, Fvalue: int64(int32(years))}
	obj.FieldTable["months"] = object.Field{Ftype: types.Int, Fvalue: int64(int32(months))}
	obj.FieldTable["days"] = object.Field{Ftype: types.Int, Fvalue: int64(int32(days))}
	return obj
}

// createPeriod creates a Period object ensuring statics are initialized.
func createPeriod(years, months, days int64) *object.Object {
	ensurePeriodInited()
	if years == 0 && months == 0 && days == 0 && periodZero != nil {
		return periodZero
	}
	return createPeriodInternal(years, months, days)
}

// CreatePeriod is an exported helper for other packages.
func CreatePeriod(years, months, days int64) *object.Object {
	return createPeriod(years, months, days)
}

func getPeriodYMD(obj *object.Object) (years, months, days int64, ok bool) {
	if obj == nil || obj.FieldTable == nil {
		return 0, 0, 0, false
	}
	if yFld, okY := obj.FieldTable["years"]; okY {
		years, _ = getInt64Val(yFld.Fvalue)
	}
	if mFld, okM := obj.FieldTable["months"]; okM {
		months, _ = getInt64Val(mFld.Fvalue)
	}
	if dFld, okD := obj.FieldTable["days"]; okD {
		days, _ = getInt64Val(dFld.Fvalue)
	}
	return years, months, days, true
}

func checkInt32Overflow(val int64) bool {
	return val < math.MinInt32 || val > math.MaxInt32
}

// Load_Time_Period registers all java/time/Period methods.
func Load_Time_Period() {
	ghelpers.MethodSignatures["java/time/Period.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  periodClinit,
		}

	ghelpers.MethodSignatures["java/time/Period.addTo(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodAddTo,
		}

	ghelpers.MethodSignatures["java/time/Period.between(Ljava/time/LocalDate;Ljava/time/LocalDate;)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  periodBetween,
		}

	ghelpers.MethodSignatures["java/time/Period.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodEquals,
		}

	ghelpers.MethodSignatures["java/time/Period.from(Ljava/time/temporal/TemporalAmount;)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodFrom,
		}

	ghelpers.MethodSignatures["java/time/Period.get(Ljava/time/temporal/TemporalUnit;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodGet,
		}

	ghelpers.MethodSignatures["java/time/Period.getChronology()Ljava/time/chrono/Chronology;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  periodGetChronology,
		}

	ghelpers.MethodSignatures["java/time/Period.getChronology()Ljava/time/chrono/IsoChronology;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  periodGetChronology,
		}

	ghelpers.MethodSignatures["java/time/Period.getDays()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  periodGetDays,
		}

	ghelpers.MethodSignatures["java/time/Period.getMonths()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  periodGetMonths,
		}

	ghelpers.MethodSignatures["java/time/Period.getUnits()Ljava/util/List;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  periodGetUnits,
		}

	ghelpers.MethodSignatures["java/time/Period.getYears()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  periodGetYears,
		}

	ghelpers.MethodSignatures["java/time/Period.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  periodHashCode,
		}

	ghelpers.MethodSignatures["java/time/Period.isNegative()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  periodIsNegative,
		}

	ghelpers.MethodSignatures["java/time/Period.isZero()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  periodIsZero,
		}

	ghelpers.MethodSignatures["java/time/Period.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodMinus,
		}

	ghelpers.MethodSignatures["java/time/Period.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/chrono/ChronoPeriod;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodMinus,
		}

	ghelpers.MethodSignatures["java/time/Period.minusDays(J)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodMinusDays,
		}

	ghelpers.MethodSignatures["java/time/Period.minusMonths(J)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodMinusMonths,
		}

	ghelpers.MethodSignatures["java/time/Period.minusYears(J)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodMinusYears,
		}

	ghelpers.MethodSignatures["java/time/Period.multipliedBy(I)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodMultipliedBy,
		}

	ghelpers.MethodSignatures["java/time/Period.multipliedBy(I)Ljava/time/chrono/ChronoPeriod;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodMultipliedBy,
		}

	ghelpers.MethodSignatures["java/time/Period.negated()Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  periodNegated,
		}

	ghelpers.MethodSignatures["java/time/Period.negated()Ljava/time/chrono/ChronoPeriod;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  periodNegated,
		}

	ghelpers.MethodSignatures["java/time/Period.normalized()Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  periodNormalized,
		}

	ghelpers.MethodSignatures["java/time/Period.normalized()Ljava/time/chrono/ChronoPeriod;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  periodNormalized,
		}

	ghelpers.MethodSignatures["java/time/Period.of(III)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  periodOf,
		}

	ghelpers.MethodSignatures["java/time/Period.ofDays(I)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodOfDays,
		}

	ghelpers.MethodSignatures["java/time/Period.ofMonths(I)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodOfMonths,
		}

	ghelpers.MethodSignatures["java/time/Period.ofWeeks(I)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodOfWeeks,
		}

	ghelpers.MethodSignatures["java/time/Period.ofYears(I)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodOfYears,
		}

	ghelpers.MethodSignatures["java/time/Period.parse(Ljava/lang/CharSequence;)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodParse,
		}

	ghelpers.MethodSignatures["java/time/Period.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodPlus,
		}

	ghelpers.MethodSignatures["java/time/Period.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/chrono/ChronoPeriod;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodPlus,
		}

	ghelpers.MethodSignatures["java/time/Period.plusDays(J)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodPlusDays,
		}

	ghelpers.MethodSignatures["java/time/Period.plusMonths(J)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodPlusMonths,
		}

	ghelpers.MethodSignatures["java/time/Period.plusYears(J)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodPlusYears,
		}

	ghelpers.MethodSignatures["java/time/Period.subtractFrom(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodSubtractFrom,
		}

	ghelpers.MethodSignatures["java/time/Period.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  periodToString,
		}

	ghelpers.MethodSignatures["java/time/Period.toTotalMonths()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  periodToTotalMonths,
		}

	ghelpers.MethodSignatures["java/time/Period.withDays(I)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodWithDays,
		}

	ghelpers.MethodSignatures["java/time/Period.withMonths(I)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodWithMonths,
		}

	ghelpers.MethodSignatures["java/time/Period.withYears(I)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  periodWithYears,
		}
}

func periodClinit(_ []any) any {
	ensurePeriodInited()
	return nil
}

// periodOf: static Period of(int years, int months, int days)
func periodOf(params []any) any {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.of requires 3 arguments")
	}
	years, ok1 := getInt64Val(params[0])
	months, ok2 := getInt64Val(params[1])
	days, ok3 := getInt64Val(params[2])
	if !ok1 || !ok2 || !ok3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.of: invalid arguments")
	}
	return createPeriod(years, months, days)
}

// periodOfYears: static Period ofYears(int years)
func periodOfYears(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.ofYears requires 1 argument")
	}
	years, ok := getInt64Val(params[0])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.ofYears: invalid argument")
	}
	return createPeriod(years, 0, 0)
}

// periodOfMonths: static Period ofMonths(int months)
func periodOfMonths(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.ofMonths requires 1 argument")
	}
	months, ok := getInt64Val(params[0])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.ofMonths: invalid argument")
	}
	return createPeriod(0, months, 0)
}

// periodOfWeeks: static Period ofWeeks(int weeks)
func periodOfWeeks(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.ofWeeks requires 1 argument")
	}
	weeks, ok := getInt64Val(params[0])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.ofWeeks: invalid argument")
	}
	days := weeks * 7
	if checkInt32Overflow(days) {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Period.ofWeeks: integer overflow")
	}
	return createPeriod(0, 0, days)
}

// periodOfDays: static Period ofDays(int days)
func periodOfDays(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.ofDays requires 1 argument")
	}
	days, ok := getInt64Val(params[0])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.ofDays: invalid argument")
	}
	return createPeriod(0, 0, days)
}

// periodBetween: static Period between(LocalDate startDateInclusive, LocalDate endDateExclusive)
func periodBetween(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.between: startDateInclusive is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.between: endDateExclusive is null")
	}
	startDate := params[0].(*object.Object)
	endDate := params[1].(*object.Object)

	return localDateUntilPeriod([]any{startDate, endDate})
}

// periodFrom: static Period from(TemporalAmount amount)
func periodFrom(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.from: amount is null")
	}
	amountObj := params[0].(*object.Object)

	if object.GoStringFromStringPoolIndex(amountObj.KlassName) == classNamePeriod {
		return amountObj
	}

	chronoObj := getChronoPeriodChronologyObject(amountObj)
	if chronoObj != nil {
		chronoId := getChronologyId(chronoObj)
		if chronoId != "" && chronoId != "ISO" {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Chronology mismatch, expected: ISO, actual: %s", chronoId))
		}
	}

	years, months, days := getChronoPeriodValues(amountObj)
	return createPeriod(years, months, days)
}

var periodRegex = regexp.MustCompile(`^(?i)([-+]?)P(?:([-+]?[0-9]+)Y)?(?:([-+]?[0-9]+)M)?(?:([-+]?[0-9]+)W)?(?:([-+]?[0-9]+)D)?$`)

// periodParse: static Period parse(CharSequence text)
func periodParse(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.parse: text is null")
	}
	str := object.GoStringFromStringObject(params[0].(*object.Object))
	if str == "" {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Period.parse: empty string")
	}

	matches := periodRegex.FindStringSubmatch(str)
	if matches == nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Text '%s' cannot be parsed to a Period", str))
	}

	negate := matches[1] == "-"
	hasY := matches[2] != ""
	hasM := matches[3] != ""
	hasW := matches[4] != ""
	hasD := matches[5] != ""

	if !hasY && !hasM && !hasW && !hasD {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Text '%s' cannot be parsed to a Period", str))
	}

	var years, months, weeks, days int64
	var err error
	if hasY {
		years, err = strconv.ParseInt(matches[2], 10, 64)
		if err != nil || checkInt32Overflow(years) {
			return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Text '%s' cannot be parsed to a Period", str))
		}
	}
	if hasM {
		months, err = strconv.ParseInt(matches[3], 10, 64)
		if err != nil || checkInt32Overflow(months) {
			return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Text '%s' cannot be parsed to a Period", str))
		}
	}
	if hasW {
		weeks, err = strconv.ParseInt(matches[4], 10, 64)
		if err != nil || checkInt32Overflow(weeks) {
			return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Text '%s' cannot be parsed to a Period", str))
		}
	}
	if hasD {
		days, err = strconv.ParseInt(matches[5], 10, 64)
		if err != nil || checkInt32Overflow(days) {
			return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Text '%s' cannot be parsed to a Period", str))
		}
	}

	totDays := days + weeks*7
	if checkInt32Overflow(totDays) {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, fmt.Sprintf("Text '%s' cannot be parsed to a Period: overflow", str))
	}

	if negate {
		years = -years
		months = -months
		totDays = -totDays
	}

	return createPeriod(years, months, totDays)
}

// periodGetYears: int getYears()
func periodGetYears(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.getYears: self is null")
	}
	self := params[0].(*object.Object)
	years, _, _, _ := getPeriodYMD(self)
	return years
}

// periodGetMonths: int getMonths()
func periodGetMonths(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.getMonths: self is null")
	}
	self := params[0].(*object.Object)
	_, months, _, _ := getPeriodYMD(self)
	return months
}

// periodGetDays: int getDays()
func periodGetDays(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.getDays: self is null")
	}
	self := params[0].(*object.Object)
	_, _, days, _ := getPeriodYMD(self)
	return days
}

// periodIsZero: boolean isZero()
func periodIsZero(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.isZero: self is null")
	}
	self := params[0].(*object.Object)
	years, months, days, _ := getPeriodYMD(self)
	if years == 0 && months == 0 && days == 0 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// periodIsNegative: boolean isNegative()
func periodIsNegative(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.isNegative: self is null")
	}
	self := params[0].(*object.Object)
	years, months, days, _ := getPeriodYMD(self)
	if years < 0 || months < 0 || days < 0 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// periodWithYears: Period withYears(int years)
func periodWithYears(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.withYears: self is null")
	}
	self := params[0].(*object.Object)
	newYears, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.withYears: invalid years")
	}
	years, months, days, _ := getPeriodYMD(self)
	if newYears == years {
		return self
	}
	return createPeriod(newYears, months, days)
}

// periodWithMonths: Period withMonths(int months)
func periodWithMonths(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.withMonths: self is null")
	}
	self := params[0].(*object.Object)
	newMonths, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.withMonths: invalid months")
	}
	years, months, days, _ := getPeriodYMD(self)
	if newMonths == months {
		return self
	}
	return createPeriod(years, newMonths, days)
}

// periodWithDays: Period withDays(int days)
func periodWithDays(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.withDays: self is null")
	}
	self := params[0].(*object.Object)
	newDays, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.withDays: invalid days")
	}
	years, months, days, _ := getPeriodYMD(self)
	if newDays == days {
		return self
	}
	return createPeriod(years, months, newDays)
}

// periodPlus: Period plus(TemporalAmount amountToAdd)
func periodPlus(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.plus: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.plus: amountToAdd is null")
	}
	self := params[0].(*object.Object)
	amountToAdd := params[1].(*object.Object)

	otherObj := periodFrom([]any{amountToAdd})
	if err, isErr := otherObj.(*ghelpers.GErrBlk); isErr {
		return err
	}
	other := otherObj.(*object.Object)

	y1, m1, d1, _ := getPeriodYMD(self)
	y2, m2, d2, _ := getPeriodYMD(other)

	ny := y1 + y2
	nm := m1 + m2
	nd := d1 + d2

	if checkInt32Overflow(ny) || checkInt32Overflow(nm) || checkInt32Overflow(nd) {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Period.plus: integer overflow")
	}
	return createPeriod(ny, nm, nd)
}

// periodPlusYears: Period plusYears(long yearsToAdd)
func periodPlusYears(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.plusYears: self is null")
	}
	self := params[0].(*object.Object)
	yearsToAdd, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.plusYears: invalid yearsToAdd")
	}
	if yearsToAdd == 0 {
		return self
	}
	years, months, days, _ := getPeriodYMD(self)
	ny := years + yearsToAdd
	if checkInt32Overflow(ny) {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Period.plusYears: integer overflow")
	}
	return createPeriod(ny, months, days)
}

// periodPlusMonths: Period plusMonths(long monthsToAdd)
func periodPlusMonths(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.plusMonths: self is null")
	}
	self := params[0].(*object.Object)
	monthsToAdd, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.plusMonths: invalid monthsToAdd")
	}
	if monthsToAdd == 0 {
		return self
	}
	years, months, days, _ := getPeriodYMD(self)
	nm := months + monthsToAdd
	if checkInt32Overflow(nm) {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Period.plusMonths: integer overflow")
	}
	return createPeriod(years, nm, days)
}

// periodPlusDays: Period plusDays(long daysToAdd)
func periodPlusDays(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.plusDays: self is null")
	}
	self := params[0].(*object.Object)
	daysToAdd, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.plusDays: invalid daysToAdd")
	}
	if daysToAdd == 0 {
		return self
	}
	years, months, days, _ := getPeriodYMD(self)
	nd := days + daysToAdd
	if checkInt32Overflow(nd) {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Period.plusDays: integer overflow")
	}
	return createPeriod(years, months, nd)
}

// periodMinus: Period minus(TemporalAmount amountToSubtract)
func periodMinus(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.minus: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.minus: amountToSubtract is null")
	}
	self := params[0].(*object.Object)
	amountToSubtract := params[1].(*object.Object)

	otherObj := periodFrom([]any{amountToSubtract})
	if err, isErr := otherObj.(*ghelpers.GErrBlk); isErr {
		return err
	}
	other := otherObj.(*object.Object)

	y1, m1, d1, _ := getPeriodYMD(self)
	y2, m2, d2, _ := getPeriodYMD(other)

	ny := y1 - y2
	nm := m1 - m2
	nd := d1 - d2

	if checkInt32Overflow(ny) || checkInt32Overflow(nm) || checkInt32Overflow(nd) {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Period.minus: integer overflow")
	}
	return createPeriod(ny, nm, nd)
}

// periodMinusYears: Period minusYears(long yearsToSubtract)
func periodMinusYears(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.minusYears: self is null")
	}
	self := params[0].(*object.Object)
	yearsToSubtract, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.minusYears: invalid yearsToSubtract")
	}
	if yearsToSubtract == 0 {
		return self
	}
	years, months, days, _ := getPeriodYMD(self)
	ny := years - yearsToSubtract
	if checkInt32Overflow(ny) {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Period.minusYears: integer overflow")
	}
	return createPeriod(ny, months, days)
}

// periodMinusMonths: Period minusMonths(long monthsToSubtract)
func periodMinusMonths(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.minusMonths: self is null")
	}
	self := params[0].(*object.Object)
	monthsToSubtract, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.minusMonths: invalid monthsToSubtract")
	}
	if monthsToSubtract == 0 {
		return self
	}
	years, months, days, _ := getPeriodYMD(self)
	nm := months - monthsToSubtract
	if checkInt32Overflow(nm) {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Period.minusMonths: integer overflow")
	}
	return createPeriod(years, nm, days)
}

// periodMinusDays: Period minusDays(long daysToSubtract)
func periodMinusDays(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.minusDays: self is null")
	}
	self := params[0].(*object.Object)
	daysToSubtract, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.minusDays: invalid daysToSubtract")
	}
	if daysToSubtract == 0 {
		return self
	}
	years, months, days, _ := getPeriodYMD(self)
	nd := days - daysToSubtract
	if checkInt32Overflow(nd) {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Period.minusDays: integer overflow")
	}
	return createPeriod(years, months, nd)
}

// periodMultipliedBy: Period multipliedBy(int scalar)
func periodMultipliedBy(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.multipliedBy: self is null")
	}
	self := params[0].(*object.Object)
	scalar, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Period.multipliedBy: invalid scalar")
	}
	if scalar == 1 {
		return self
	}
	years, months, days, _ := getPeriodYMD(self)
	ny := years * scalar
	nm := months * scalar
	nd := days * scalar
	if checkInt32Overflow(ny) || checkInt32Overflow(nm) || checkInt32Overflow(nd) {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Period.multipliedBy: integer overflow")
	}
	return createPeriod(ny, nm, nd)
}

// periodNegated: Period negated()
func periodNegated(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.negated: self is null")
	}
	self := params[0].(*object.Object)
	return periodMultipliedBy([]any{self, int64(-1)})
}

// periodNormalized: Period normalized()
func periodNormalized(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.normalized: self is null")
	}
	self := params[0].(*object.Object)
	years, months, days, _ := getPeriodYMD(self)
	totalMonths := years*12 + months
	splitYears := totalMonths / 12
	splitMonths := totalMonths % 12
	if splitYears == years && splitMonths == months {
		return self
	}
	if checkInt32Overflow(splitYears) {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Period.normalized: integer overflow")
	}
	return createPeriod(splitYears, splitMonths, days)
}

// periodToTotalMonths: long toTotalMonths()
func periodToTotalMonths(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.toTotalMonths: self is null")
	}
	self := params[0].(*object.Object)
	years, months, _, _ := getPeriodYMD(self)
	return years*12 + months
}

// periodGetChronology: IsoChronology getChronology()
func periodGetChronology(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.getChronology: self is null")
	}
	return GetIsoChronology()
}

// periodGetUnits: List<TemporalUnit> getUnits()
func periodGetUnits(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.getUnits: self is null")
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

// periodGet: long get(TemporalUnit unit)
func periodGet(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.get: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.get: unit is null")
	}
	self := params[0].(*object.Object)
	unitObj := params[1].(*object.Object)

	unitName := getTemporalUnitName(unitObj)
	years, months, days, _ := getPeriodYMD(self)

	switch unitName {
	case "YEARS":
		return years
	case "MONTHS":
		return months
	case "DAYS":
		return days
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported unit: %s", unitName))
	}
}

// periodAddTo: Temporal addTo(Temporal temporal)
func periodAddTo(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.addTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.addTo: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)

	years, months, days, _ := getPeriodYMD(self)
	temporalChrono := GetChronoLocalDateChronology(temporal)
	if temporalChrono != nil {
		chronoId := getChronologyId(temporalChrono)
		if chronoId != "" && chronoId != "ISO" {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Chronology mismatch, expected: ISO, actual: %s", chronoId))
		}
	}

	res := temporal
	if years != 0 {
		if months != 0 {
			toTotal := years*12 + months
			plusRes := temporalPlusUnit([]any{res, toTotal, ghelpers.CreateTemporalUnit("MONTHS")})
			if obj, ok := plusRes.(*object.Object); ok {
				res = obj
			} else {
				return plusRes
			}
		} else {
			plusRes := temporalPlusUnit([]any{res, years, ghelpers.CreateTemporalUnit("YEARS")})
			if obj, ok := plusRes.(*object.Object); ok {
				res = obj
			} else {
				return plusRes
			}
		}
	} else if months != 0 {
		plusRes := temporalPlusUnit([]any{res, months, ghelpers.CreateTemporalUnit("MONTHS")})
		if obj, ok := plusRes.(*object.Object); ok {
			res = obj
		} else {
			return plusRes
		}
	}

	if days != 0 {
		plusRes := temporalPlusUnit([]any{res, days, ghelpers.CreateTemporalUnit("DAYS")})
		if obj, ok := plusRes.(*object.Object); ok {
			res = obj
		} else {
			return plusRes
		}
	}

	return res
}

// periodSubtractFrom: Temporal subtractFrom(Temporal temporal)
func periodSubtractFrom(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.subtractFrom: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.subtractFrom: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)

	years, months, days, _ := getPeriodYMD(self)
	temporalChrono := GetChronoLocalDateChronology(temporal)
	if temporalChrono != nil {
		chronoId := getChronologyId(temporalChrono)
		if chronoId != "" && chronoId != "ISO" {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Chronology mismatch, expected: ISO, actual: %s", chronoId))
		}
	}

	res := temporal
	if years != 0 {
		if months != 0 {
			toTotal := years*12 + months
			minusRes := temporalMinusUnit([]any{res, toTotal, ghelpers.CreateTemporalUnit("MONTHS")})
			if obj, ok := minusRes.(*object.Object); ok {
				res = obj
			} else {
				return minusRes
			}
		} else {
			minusRes := temporalMinusUnit([]any{res, years, ghelpers.CreateTemporalUnit("YEARS")})
			if obj, ok := minusRes.(*object.Object); ok {
				res = obj
			} else {
				return minusRes
			}
		}
	} else if months != 0 {
		minusRes := temporalMinusUnit([]any{res, months, ghelpers.CreateTemporalUnit("MONTHS")})
		if obj, ok := minusRes.(*object.Object); ok {
			res = obj
		} else {
			return minusRes
		}
	}

	if days != 0 {
		minusRes := temporalMinusUnit([]any{res, days, ghelpers.CreateTemporalUnit("DAYS")})
		if obj, ok := minusRes.(*object.Object); ok {
			res = obj
		} else {
			return minusRes
		}
	}

	return res
}

// periodEquals: boolean equals(Object obj)
func periodEquals(params []any) any {
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
	if object.GoStringFromStringPoolIndex(other.KlassName) != classNamePeriod {
		return types.JavaBoolFalse
	}

	y1, m1, d1, _ := getPeriodYMD(self)
	y2, m2, d2, _ := getPeriodYMD(other)

	if y1 == y2 && m1 == m2 && d1 == d2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// periodHashCode: int hashCode()
func periodHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.hashCode: self is null")
	}
	self := params[0].(*object.Object)
	y, m, d, _ := getPeriodYMD(self)

	hash := uint32(int32(y)) + bits.RotateLeft32(uint32(int32(m)), 8) + bits.RotateLeft32(uint32(int32(d)), 16)
	return int64(int32(hash))
}

// periodToString: String toString()
func periodToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Period.toString: self is null")
	}
	self := params[0].(*object.Object)
	years, months, days, _ := getPeriodYMD(self)

	if years == 0 && months == 0 && days == 0 {
		return object.StringObjectFromGoString("P0D")
	}

	var sb strings.Builder
	sb.WriteByte('P')
	if years != 0 {
		sb.WriteString(fmt.Sprintf("%dY", years))
	}
	if months != 0 {
		sb.WriteString(fmt.Sprintf("%dM", months))
	}
	if days != 0 {
		sb.WriteString(fmt.Sprintf("%dD", days))
	}

	return object.StringObjectFromGoString(sb.String())
}
