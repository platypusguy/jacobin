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
	"jacobin/src/statics"
	"jacobin/src/types"
	"math"
	"math/big"
	"strconv"
	"strings"
	"sync"
)

var classNameDuration = "java/time/Duration"

var (
	durationMutex       sync.Mutex
	durationInitialized bool
)

func Load_Time_Duration() {
	// Static methods
	ghelpers.MethodSignatures["java/time/Duration.<clinit>()V"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationClinit}
	ghelpers.MethodSignatures["java/time/Duration.abs()Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationAbs}
	ghelpers.MethodSignatures["java/time/Duration.between(Ljava/time/temporal/Temporal;Ljava/time/temporal/Temporal;)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 2, GFunction: durationBetween}
	ghelpers.MethodSignatures["java/time/Duration.compareTo(Ljava/time/Duration;)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationCompareTo}
	ghelpers.MethodSignatures["java/time/Duration.dividedBy(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationDividedByLong}
	ghelpers.MethodSignatures["java/time/Duration.dividedBy(Ljava/time/Duration;)J"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationDividedByDuration}
	ghelpers.MethodSignatures["java/time/Duration.equals(Ljava/lang/Object;)Z"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationEquals}
	ghelpers.MethodSignatures["java/time/Duration.from(Ljava/time/temporal/TemporalAmount;)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationFrom}
	ghelpers.MethodSignatures["java/time/Duration.get(Ljava/time/temporal/TemporalUnit;)J"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationGet}
	ghelpers.MethodSignatures["java/time/Duration.getNano()I"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationGetNano}
	ghelpers.MethodSignatures["java/time/Duration.getSeconds()J"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationGetSeconds}
	ghelpers.MethodSignatures["java/time/Duration.getUnits()Ljava/util/List;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationGetUnits}
	ghelpers.MethodSignatures["java/time/Duration.hashCode()I"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationHashCode}
	ghelpers.MethodSignatures["java/time/Duration.isNegative()Z"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationIsNegative}
	ghelpers.MethodSignatures["java/time/Duration.isPositive()Z"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationIsPositive}
	ghelpers.MethodSignatures["java/time/Duration.isZero()Z"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationIsZero}
	ghelpers.MethodSignatures["java/time/Duration.minus(Ljava/time/Duration;)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationMinus}
	ghelpers.MethodSignatures["java/time/Duration.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 2, GFunction: durationMinusUnit}
	ghelpers.MethodSignatures["java/time/Duration.minusDays(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationMinusDays}
	ghelpers.MethodSignatures["java/time/Duration.minusHours(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationMinusHours}
	ghelpers.MethodSignatures["java/time/Duration.minusMillis(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationMinusMillis}
	ghelpers.MethodSignatures["java/time/Duration.minusMinutes(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationMinusMinutes}
	ghelpers.MethodSignatures["java/time/Duration.minusNanos(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationMinusNanos}
	ghelpers.MethodSignatures["java/time/Duration.minusSeconds(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationMinusSeconds}
	ghelpers.MethodSignatures["java/time/Duration.multipliedBy(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationMultipliedBy}
	ghelpers.MethodSignatures["java/time/Duration.negated()Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationNegated}
	ghelpers.MethodSignatures["java/time/Duration.of(JLjava/time/temporal/TemporalUnit;)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 2, GFunction: durationOf}
	ghelpers.MethodSignatures["java/time/Duration.ofDays(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationOfDays}
	ghelpers.MethodSignatures["java/time/Duration.ofHours(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationOfHours}
	ghelpers.MethodSignatures["java/time/Duration.ofMillis(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationOfMillis}
	ghelpers.MethodSignatures["java/time/Duration.ofMinutes(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationOfMinutes}
	ghelpers.MethodSignatures["java/time/Duration.ofNanos(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationOfNanos}
	ghelpers.MethodSignatures["java/time/Duration.ofSeconds(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationOfSeconds}
	ghelpers.MethodSignatures["java/time/Duration.ofSeconds(JJ)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 2, GFunction: durationOfSecondsNanos}
	ghelpers.MethodSignatures["java/time/Duration.parse(Ljava/lang/CharSequence;)Ljava/time/Duration;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: durationParse}
	ghelpers.MethodSignatures["java/time/Duration.plus(Ljava/time/Duration;)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationPlus}
	ghelpers.MethodSignatures["java/time/Duration.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 2, GFunction: durationPlusUnit}
	ghelpers.MethodSignatures["java/time/Duration.plusDays(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationPlusDays}
	ghelpers.MethodSignatures["java/time/Duration.plusHours(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationPlusHours}
	ghelpers.MethodSignatures["java/time/Duration.plusMillis(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationPlusMillis}
	ghelpers.MethodSignatures["java/time/Duration.plusMinutes(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationPlusMinutes}
	ghelpers.MethodSignatures["java/time/Duration.plusNanos(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationPlusNanos}
	ghelpers.MethodSignatures["java/time/Duration.plusSeconds(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationPlusSeconds}
	ghelpers.MethodSignatures["java/time/Duration.toDays()J"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationToDays}
	ghelpers.MethodSignatures["java/time/Duration.toDaysPart()J"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationToDaysPart}
	ghelpers.MethodSignatures["java/time/Duration.toHours()J"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationToHours}
	ghelpers.MethodSignatures["java/time/Duration.toHoursPart()I"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationToHoursPart}
	ghelpers.MethodSignatures["java/time/Duration.toMillis()J"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationToMillis}
	ghelpers.MethodSignatures["java/time/Duration.toMillisPart()I"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationToMillisPart}
	ghelpers.MethodSignatures["java/time/Duration.toMinutes()J"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationToMinutes}
	ghelpers.MethodSignatures["java/time/Duration.toMinutesPart()I"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationToMinutesPart}
	ghelpers.MethodSignatures["java/time/Duration.toNanos()J"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationToNanos}
	ghelpers.MethodSignatures["java/time/Duration.toNanosPart()I"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationToNanosPart}
	ghelpers.MethodSignatures["java/time/Duration.toSeconds()J"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationToSeconds}
	ghelpers.MethodSignatures["java/time/Duration.toSecondsPart()I"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationToSecondsPart}
	ghelpers.MethodSignatures["java/time/Duration.toString()Ljava/lang/String;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationToString}
	ghelpers.MethodSignatures["java/time/Duration.truncatedTo(Ljava/time/temporal/TemporalUnit;)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationTruncatedTo}
	ghelpers.MethodSignatures["java/time/Duration.withNanos(I)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationWithNanos}
	ghelpers.MethodSignatures["java/time/Duration.withSeconds(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationWithSeconds}
	ghelpers.MethodSignatures["java/time/Duration.addTo(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationAddTo}
	ghelpers.MethodSignatures["java/time/Duration.subtractFrom(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationSubtractFrom}
}

// CreateDuration creates a new Duration object
func CreateDuration(seconds int64, nanos int32) *object.Object {
	return createDuration(seconds, nanos)
}

// Helper to create a new Duration object
func createDuration(seconds int64, nanos int32) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameDuration)
	obj.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: seconds}
	obj.FieldTable["nanos"] = object.Field{Ftype: types.Int, Fvalue: int64(nanos)}
	return obj
}

func getDurationSecondsNanos(dur *object.Object) (int64, int32, bool) {
	if dur == nil || dur.FieldTable == nil {
		return 0, 0, false
	}
	secFld, okSec := dur.FieldTable["seconds"]
	nanoFld, okNano := dur.FieldTable["nanos"]
	if !okSec || !okNano {
		return 0, 0, false
	}
	sec, ok1 := getInt64Val(secFld.Fvalue)
	nano, ok2 := getInt64Val(nanoFld.Fvalue)
	if !ok1 || !ok2 {
		return 0, 0, false
	}
	return sec, int32(nano), true
}

func addExact(a, b int64) (int64, *ghelpers.GErrBlk) {
	res := a + b
	if (a^res) < 0 && (a^b) >= 0 {
		return 0, ghelpers.GetGErrBlk(excNames.ArithmeticException, "long overflow")
	}
	return res, nil
}

func multiplyExact(a, b int64) (int64, *ghelpers.GErrBlk) {
	if b == 0 {
		return 0, nil
	}
	res := a * b
	if a != 0 && res/a != b || (a == -1 && b == math.MinInt64) {
		return 0, ghelpers.GetGErrBlk(excNames.ArithmeticException, "long overflow")
	}
	return res, nil
}

func durationClinit(params []interface{}) interface{} {
	obj := createDuration(0, 0)
	name := fmt.Sprintf("%s.ZERO", classNameDuration)
	_ = statics.AddStatic(name, statics.Static{Type: "Ljava/time/Duration;", Value: obj})
	return nil
}

const (
	nanosPerSecond = 1000000000
)

func durationOfSeconds(params []interface{}) interface{} {
	seconds := params[0].(int64)
	return createDuration(seconds, 0)
}

func durationOfSecondsNanos(params []interface{}) interface{} {
	seconds := params[0].(int64)
	nanos := params[1].(int64)
	secAdjustment := nanos / nanosPerSecond
	nanos = nanos % nanosPerSecond
	if nanos < 0 {
		nanos += nanosPerSecond
		secAdjustment--
	}
	resSec, err := addExact(seconds, secAdjustment)
	if err != nil {
		return err
	}
	return createDuration(resSec, int32(nanos))
}

func durationOfMillis(params []interface{}) interface{} {
	millis := params[0].(int64)
	secs := millis / 1000
	mos := int32(millis % 1000)
	if mos < 0 {
		mos += 1000
		secs--
	}
	return createDuration(secs, mos*1000000)
}

func durationOfNanos(params []interface{}) interface{} {
	nanos := params[0].(int64)
	secs := nanos / nanosPerSecond
	nos := int32(nanos % nanosPerSecond)
	if nos < 0 {
		nos += nanosPerSecond
		secs--
	}
	return createDuration(secs, nos)
}

func durationOfMinutes(params []interface{}) interface{} {
	minutes := params[0].(int64)
	secs, err := multiplyExact(minutes, 60)
	if err != nil {
		return err
	}
	return createDuration(secs, 0)
}

func durationOfHours(params []interface{}) interface{} {
	hours := params[0].(int64)
	secs, err := multiplyExact(hours, 3600)
	if err != nil {
		return err
	}
	return createDuration(secs, 0)
}

func durationOfDays(params []interface{}) interface{} {
	days := params[0].(int64)
	secs, err := multiplyExact(days, 86400)
	if err != nil {
		return err
	}
	return createDuration(secs, 0)
}

func durationGetSeconds(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	return self.FieldTable["seconds"].Fvalue.(int64)
}

func durationGetNano(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	return self.FieldTable["nanos"].Fvalue.(int64)
}

func durationIsNegative(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	seconds := self.FieldTable["seconds"].Fvalue.(int64)
	if seconds < 0 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func durationIsPositive(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	seconds := self.FieldTable["seconds"].Fvalue.(int64)
	nanos := self.FieldTable["nanos"].Fvalue.(int64)
	if seconds > 0 || (seconds == 0 && nanos > 0) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func durationIsZero(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	seconds := self.FieldTable["seconds"].Fvalue.(int64)
	nanos := self.FieldTable["nanos"].Fvalue.(int64)
	if seconds == 0 && nanos == 0 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func durationPlus(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)
	s1 := self.FieldTable["seconds"].Fvalue.(int64)
	n1 := int32(self.FieldTable["nanos"].Fvalue.(int64))
	s2 := other.FieldTable["seconds"].Fvalue.(int64)
	n2 := int32(other.FieldTable["nanos"].Fvalue.(int64))

	resSeconds, err := addExact(s1, s2)
	if err != nil {
		return err
	}
	resNanos := n1 + n2
	if resNanos >= nanosPerSecond {
		resSeconds, err = addExact(resSeconds, 1)
		if err != nil {
			return err
		}
		resNanos -= nanosPerSecond
	} else if resNanos < 0 {
		resSeconds, err = addExact(resSeconds, -1)
		if err != nil {
			return err
		}
		resNanos += nanosPerSecond
	}
	return createDuration(resSeconds, resNanos)
}

func durationPlusSeconds(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	secondsToAdd := params[1].(int64)
	s1 := self.FieldTable["seconds"].Fvalue.(int64)
	n1 := int32(self.FieldTable["nanos"].Fvalue.(int64))
	resSec, err := addExact(s1, secondsToAdd)
	if err != nil {
		return err
	}
	return createDuration(resSec, n1)
}

func durationPlusNanos(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	nanosToAdd := params[1].(int64)
	s1 := self.FieldTable["seconds"].Fvalue.(int64)
	n1 := int32(self.FieldTable["nanos"].Fvalue.(int64))

	secAdjustment := nanosToAdd / nanosPerSecond
	nanosToAdd = nanosToAdd % nanosPerSecond

	resSeconds, err := addExact(s1, secAdjustment)
	if err != nil {
		return err
	}
	resNanos := n1 + int32(nanosToAdd)

	if resNanos >= nanosPerSecond {
		resSeconds, err = addExact(resSeconds, 1)
		if err != nil {
			return err
		}
		resNanos -= nanosPerSecond
	} else if resNanos < 0 {
		resSeconds, err = addExact(resSeconds, -1)
		if err != nil {
			return err
		}
		resNanos += nanosPerSecond
	}
	return createDuration(resSeconds, resNanos)
}

func durationPlusMillis(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	millisToAdd := params[1].(int64)
	nanosToAdd, err := multiplyExact(millisToAdd, 1000000)
	if err != nil {
		// Java 21 Duration.plusMillis(long) says:
		// ArithmeticException - if numeric overflow occurs
		// But wait, if millisToAdd * 1M overflows long, we can still potentially add it
		// if we do it in terms of seconds and nanos.
		secs := millisToAdd / 1000
		mos := int32(millisToAdd % 1000)
		if mos < 0 {
			mos += 1000
			secs--
		}
		d2 := createDuration(secs, mos*1000000)
		return durationPlus([]interface{}{self, d2})
	}
	return durationPlusNanos([]interface{}{self, nanosToAdd})
}

func durationPlusMinutes(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	minutesToAdd := params[1].(int64)
	secsToAdd, err := multiplyExact(minutesToAdd, 60)
	if err != nil {
		d2 := durationOfMinutes([]interface{}{minutesToAdd})
		if e, ok := d2.(*ghelpers.GErrBlk); ok {
			return e
		}
		return durationPlus([]interface{}{self, d2})
	}
	return durationPlusSeconds([]interface{}{self, secsToAdd})
}

func durationPlusHours(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	hoursToAdd := params[1].(int64)
	secsToAdd, err := multiplyExact(hoursToAdd, 3600)
	if err != nil {
		d2 := durationOfHours([]interface{}{hoursToAdd})
		if e, ok := d2.(*ghelpers.GErrBlk); ok {
			return e
		}
		return durationPlus([]interface{}{self, d2})
	}
	return durationPlusSeconds([]interface{}{self, secsToAdd})
}

func durationPlusDays(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	daysToAdd := params[1].(int64)
	secsToAdd, err := multiplyExact(daysToAdd, 86400)
	if err != nil {
		d2 := durationOfDays([]interface{}{daysToAdd})
		if e, ok := d2.(*ghelpers.GErrBlk); ok {
			return e
		}
		return durationPlus([]interface{}{self, d2})
	}
	return durationPlusSeconds([]interface{}{self, secsToAdd})
}

func durationMinus(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)
	s2 := other.FieldTable["seconds"].Fvalue.(int64)
	n2 := int32(other.FieldTable["nanos"].Fvalue.(int64))

	negOther := createDuration(-s2, -n2)
	// Handle -n2 being -nanosPerSecond if we were strict but n2 is always 0..999999999
	if n2 > 0 {
		negOther = createDuration(-s2-1, nanosPerSecond-n2)
	} else {
		negOther = createDuration(-s2, 0)
	}

	return durationPlus([]interface{}{self, negOther})
}

func durationMinusSeconds(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	secondsToSub := params[1].(int64)
	if secondsToSub == math.MinInt64 {
		return durationPlusSeconds([]interface{}{durationPlusSeconds([]interface{}{self, int64(math.MaxInt64)}), int64(1)})
	}
	return durationPlusSeconds([]interface{}{self, -secondsToSub})
}

func durationMinusNanos(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	nanosToSub := params[1].(int64)
	if nanosToSub == math.MinInt64 {
		return durationPlusNanos([]interface{}{durationPlusNanos([]interface{}{self, int64(math.MaxInt64)}), int64(1)})
	}
	return durationPlusNanos([]interface{}{self, -nanosToSub})
}

func durationMinusMillis(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	millisToSub := params[1].(int64)
	if millisToSub == math.MinInt64 {
		return durationPlusMillis([]interface{}{durationPlusMillis([]interface{}{self, int64(math.MaxInt64)}), int64(1)})
	}
	return durationPlusMillis([]interface{}{self, -millisToSub})
}

func durationMinusMinutes(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	minutesToSub := params[1].(int64)
	if minutesToSub == math.MinInt64 {
		return durationPlusMinutes([]interface{}{durationPlusMinutes([]interface{}{self, int64(math.MaxInt64)}), int64(1)})
	}
	return durationPlusMinutes([]interface{}{self, -minutesToSub})
}

func durationMinusHours(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	hoursToSub := params[1].(int64)
	if hoursToSub == math.MinInt64 {
		return durationPlusHours([]interface{}{durationPlusHours([]interface{}{self, int64(math.MaxInt64)}), int64(1)})
	}
	return durationPlusHours([]interface{}{self, -hoursToSub})
}

func durationMinusDays(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	daysToSub := params[1].(int64)
	if daysToSub == math.MinInt64 {
		return durationPlusDays([]interface{}{durationPlusDays([]interface{}{self, int64(math.MaxInt64)}), int64(1)})
	}
	return durationPlusDays([]interface{}{self, -daysToSub})
}

func durationMultipliedBy(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Duration.multipliedBy: self is null")
	}
	self := params[0].(*object.Object)
	multiplicand, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Duration.multipliedBy: invalid multiplicand")
	}
	if multiplicand == 0 {
		return createDuration(0, 0)
	}
	if multiplicand == 1 {
		return self
	}
	s := self.FieldTable["seconds"].Fvalue.(int64)
	n := self.FieldTable["nanos"].Fvalue.(int64)

	b := big.NewInt(s)
	b.Mul(b, big.NewInt(nanosPerSecond))
	b.Add(b, big.NewInt(n))
	b.Mul(b, big.NewInt(multiplicand))

	secBig := new(big.Int).Quo(b, big.NewInt(nanosPerSecond))
	nanoBig := new(big.Int).Rem(b, big.NewInt(nanosPerSecond))

	if !secBig.IsInt64() {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "long overflow")
	}
	sec := secBig.Int64()
	nano := int32(nanoBig.Int64())
	if nano < 0 {
		nano += nanosPerSecond
		sec--
	}
	return createDuration(sec, nano)
}

func durationDividedByLong(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Duration.dividedBy: self is null")
	}
	self := params[0].(*object.Object)
	divisor, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Duration.dividedBy: invalid divisor")
	}
	if divisor == 0 {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Division by zero")
	}
	if divisor == 1 {
		return self
	}
	s := self.FieldTable["seconds"].Fvalue.(int64)
	n := self.FieldTable["nanos"].Fvalue.(int64)

	b := big.NewInt(s)
	b.Mul(b, big.NewInt(nanosPerSecond))
	b.Add(b, big.NewInt(n))

	resNanos := new(big.Int).Quo(b, big.NewInt(divisor))

	secBig := new(big.Int).Quo(resNanos, big.NewInt(nanosPerSecond))
	nanoBig := new(big.Int).Rem(resNanos, big.NewInt(nanosPerSecond))

	if !secBig.IsInt64() {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "long overflow")
	}
	sec := secBig.Int64()
	nano := int32(nanoBig.Int64())
	if nano < 0 {
		nano += nanosPerSecond
		sec--
	}
	return createDuration(sec, nano)
}

func durationDividedByDuration(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Duration.dividedBy: null argument")
	}
	self := params[0].(*object.Object)
	divisor := params[1].(*object.Object)
	s1 := self.FieldTable["seconds"].Fvalue.(int64)
	n1 := self.FieldTable["nanos"].Fvalue.(int64)
	s2 := divisor.FieldTable["seconds"].Fvalue.(int64)
	n2 := divisor.FieldTable["nanos"].Fvalue.(int64)

	if s2 == 0 && n2 == 0 {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Division by zero")
	}

	b1 := big.NewInt(s1)
	b1.Mul(b1, big.NewInt(nanosPerSecond))
	b1.Add(b1, big.NewInt(n1))

	b2 := big.NewInt(s2)
	b2.Mul(b2, big.NewInt(nanosPerSecond))
	b2.Add(b2, big.NewInt(n2))

	if b2.Sign() == 0 {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Division by zero")
	}
	res := new(big.Int).Quo(b1, b2)
	return res.Int64()
}

func durationNegated(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	s := self.FieldTable["seconds"].Fvalue.(int64)
	n := int32(self.FieldTable["nanos"].Fvalue.(int64))
	if s == math.MinInt64 {
		if n == 0 {
			return ghelpers.GetGErrBlk(excNames.ArithmeticException, "long overflow")
		}
		return createDuration(math.MaxInt64, nanosPerSecond-n)
	}
	if n == 0 {
		return createDuration(-s, 0)
	}
	return createDuration(-s-1, nanosPerSecond-n)
}

func durationAbs(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	s := self.FieldTable["seconds"].Fvalue.(int64)
	if s < 0 {
		return durationNegated([]interface{}{self})
	}
	return self
}

func durationToDays(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	s := self.FieldTable["seconds"].Fvalue.(int64)
	return s / 86400
}

func durationToHours(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	s := self.FieldTable["seconds"].Fvalue.(int64)
	return s / 3600
}

func durationToMinutes(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	s := self.FieldTable["seconds"].Fvalue.(int64)
	return s / 60
}

func durationToSeconds(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	return self.FieldTable["seconds"].Fvalue.(int64)
}

func durationToMillis(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	s := self.FieldTable["seconds"].Fvalue.(int64)
	n := self.FieldTable["nanos"].Fvalue.(int64)
	millis, err := multiplyExact(s, 1000)
	if err != nil {
		return err
	}
	res, err := addExact(millis, n/1000000)
	if err != nil {
		return err
	}
	return res
}

func durationToNanos(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	s := self.FieldTable["seconds"].Fvalue.(int64)
	n := self.FieldTable["nanos"].Fvalue.(int64)
	nanos, err := multiplyExact(s, nanosPerSecond)
	if err != nil {
		return err
	}
	res, err := addExact(nanos, n)
	if err != nil {
		return err
	}
	return res
}

func durationToDaysPart(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	s := self.FieldTable["seconds"].Fvalue.(int64)
	return s / 86400
}

func durationToHoursPart(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	s := self.FieldTable["seconds"].Fvalue.(int64)
	return int64((s / 3600) % 24)
}

func durationToMinutesPart(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	s := self.FieldTable["seconds"].Fvalue.(int64)
	return int64((s / 60) % 60)
}

func durationToSecondsPart(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	s := self.FieldTable["seconds"].Fvalue.(int64)
	return int64(s % 60)
}

func durationToMillisPart(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	n := self.FieldTable["nanos"].Fvalue.(int64)
	return int64(n / 1000000)
}

func durationToNanosPart(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	n := self.FieldTable["nanos"].Fvalue.(int64)
	return n
}

func durationWithSeconds(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	seconds := params[1].(int64)
	nanos := int32(self.FieldTable["nanos"].Fvalue.(int64))
	return createDuration(seconds, nanos)
}

func durationWithNanos(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	nanos := int32(params[1].(int64))
	seconds := self.FieldTable["seconds"].Fvalue.(int64)
	return createDuration(seconds, nanos)
}

func durationCompareTo(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)
	s1 := self.FieldTable["seconds"].Fvalue.(int64)
	n1 := self.FieldTable["nanos"].Fvalue.(int64)
	s2 := other.FieldTable["seconds"].Fvalue.(int64)
	n2 := other.FieldTable["nanos"].Fvalue.(int64)

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

func durationEquals(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	otherObj, ok := params[1].(*object.Object)
	if !ok || otherObj == nil {
		return types.JavaBoolFalse
	}
	if object.GoStringFromStringPoolIndex(otherObj.KlassName) != classNameDuration {
		return types.JavaBoolFalse
	}
	s1 := self.FieldTable["seconds"].Fvalue.(int64)
	n1 := self.FieldTable["nanos"].Fvalue.(int64)
	s2 := otherObj.FieldTable["seconds"].Fvalue.(int64)
	n2 := otherObj.FieldTable["nanos"].Fvalue.(int64)

	if s1 == s2 && n1 == n2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func durationHashCode(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	seconds := self.FieldTable["seconds"].Fvalue.(int64)
	nanos := self.FieldTable["nanos"].Fvalue.(int64)
	return int64(int32(seconds^(seconds>>32)) + 51*int32(nanos))
}

func durationToString(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	seconds := self.FieldTable["seconds"].Fvalue.(int64)
	nanos := self.FieldTable["nanos"].Fvalue.(int64)

	if seconds == 0 && nanos == 0 {
		return object.StringObjectFromGoString("PT0S")
	}

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	var sb strings.Builder
	sb.WriteString("PT")
	if hours != 0 {
		sb.WriteString(strconv.FormatInt(hours, 10))
		sb.WriteByte('H')
	}
	if minutes != 0 {
		sb.WriteString(strconv.FormatInt(minutes, 10))
		sb.WriteByte('M')
	}
	if secs != 0 || nanos != 0 || (hours == 0 && minutes == 0) {
		if secs == 0 && seconds < 0 {
			sb.WriteString("-0")
		} else {
			sb.WriteString(strconv.FormatInt(secs, 10))
		}
		if nanos != 0 {
			sb.WriteByte('.')
			nStr := fmt.Sprintf("%09d", int32(math.Abs(float64(nanos))))
			sb.WriteString(strings.TrimRight(nStr, "0"))
		}
		sb.WriteByte('S')
	}

	return object.StringObjectFromGoString(sb.String())
}

func durationParse(params []interface{}) interface{} {
	textObj := params[0].(*object.Object)
	text := object.GoStringFromStringObject(textObj)

	if text == "" {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Empty duration string")
	}

	negated := false
	if strings.HasPrefix(text, "-") {
		negated = true
		text = text[1:]
	} else if strings.HasPrefix(text, "+") {
		text = text[1:]
	}

	if !strings.HasPrefix(text, "P") {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Duration must start with P")
	}
	text = text[1:]

	var seconds int64
	var nanos int32

	// ISO-8601 parser for [n[D]]T[nH][nM][n[.n]S]
	// Each component can have an optional +/- sign.
	inTime := false
	i := 0
	for i < len(text) {
		if text[i] == 'T' {
			inTime = true
			i++
			continue
		}

		start := i
		// Handle optional sign for each component
		if i < len(text) && (text[i] == '-' || text[i] == '+') {
			i++
		}
		// Read numeric part
		for i < len(text) && ((text[i] >= '0' && text[i] <= '9') || text[i] == '.') {
			i++
		}
		if i == start || (i == start+1 && (text[start] == '-' || text[start] == '+')) {
			return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Invalid value in duration")
		}
		if i >= len(text) {
			return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Missing unit in duration")
		}

		valStr := text[start:i]
		unit := text[i]
		i++

		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Invalid value in duration")
		}

		switch unit {
		case 'D':
			if inTime {
				return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "D unit in time section")
			}
			seconds += int64(val * 86400)
		case 'H':
			if !inTime {
				return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "H unit before T")
			}
			seconds += int64(val * 3600)
		case 'M':
			if !inTime {
				return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "M unit before T")
			}
			seconds += int64(val * 60)
		case 'S':
			if !inTime {
				return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "S unit before T")
			}
			secPart := int64(val)
			nanoPart := int32(math.Round((val - float64(secPart)) * nanosPerSecond))
			seconds += secPart
			nanos += nanoPart
		default:
			return ghelpers.GetGErrBlk(excNames.DateTimeParseException, "Unknown duration unit")
		}
	}

	if negated {
		seconds = -seconds
		nanos = -nanos
	}

	// Standardize nanos to [0, 999_999_999]
	for nanos < 0 {
		seconds--
		nanos += nanosPerSecond
	}
	for nanos >= nanosPerSecond {
		seconds += int64(nanos / nanosPerSecond)
		nanos %= nanosPerSecond
	}

	return createDuration(seconds, nanos)
}

func getInstantSecNano(inst *object.Object) (int64, int64, bool) {
	if inst == nil || inst.FieldTable == nil {
		return 0, 0, false
	}
	secFld, ok1 := inst.FieldTable["seconds"]
	nanoFld, ok2 := inst.FieldTable["nanos"]
	if !ok1 || !ok2 {
		return 0, 0, false
	}
	sec, ok1 := getInt64Val(secFld.Fvalue)
	nano, ok2 := getInt64Val(nanoFld.Fvalue)
	if !ok1 || !ok2 {
		return 0, 0, false
	}
	return sec, nano, true
}

// Duration.between(Temporal startInclusive, Temporal endExclusive)
func durationBetween(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Duration.between: null argument")
	}
	start := params[0].(*object.Object)
	end := params[1].(*object.Object)

	startClass := object.GoStringFromStringPoolIndex(start.KlassName)
	endClass := object.GoStringFromStringPoolIndex(end.KlassName)

	// Instant
	if strings.Contains(startClass, "Instant") && strings.Contains(endClass, "Instant") {
		sSec, sNano, ok1 := getInstantSecNano(start)
		eSec, eNano, ok2 := getInstantSecNano(end)
		if ok1 && ok2 {
			secs, err := addExact(eSec, -sSec)
			if err != nil {
				return err
			}
			nanos := eNano - sNano
			return durationOfSecondsNanos([]any{secs, nanos})
		}
	}

	// LocalTime
	if strings.Contains(startClass, "LocalTime") && strings.Contains(endClass, "LocalTime") {
		sH, sM, sS, sN := getLocalTimeHMSN(start)
		eH, eM, eS, eN := getLocalTimeHMSN(end)
		startNanosOfDay := sH*3600_000_000_000 + sM*60_000_000_000 + sS*1_000_000_000 + sN
		endNanosOfDay := eH*3600_000_000_000 + eM*60_000_000_000 + eS*1_000_000_000 + eN
		nanosDiff := endNanosOfDay - startNanosOfDay
		return durationOfNanos([]any{nanosDiff})
	}

	// LocalDateTime
	if strings.Contains(startClass, "LocalDateTime") && strings.Contains(endClass, "LocalDateTime") {
		sDate := getLocalDateTimeDate(start)
		sTime := getLocalDateTimeTime(start)
		eDate := getLocalDateTimeDate(end)
		eTime := getLocalDateTimeTime(end)
		if sDate != nil && sTime != nil && eDate != nil && eTime != nil {
			sY, sM, sD, ok1 := getLocalDateYMD(sDate)
			eY, eM, eD, ok2 := getLocalDateYMD(eDate)
			if ok1 && ok2 {
				sEpochDay := ymdToEpochDay(sY, sM, sD)
				eEpochDay := ymdToEpochDay(eY, eM, eD)
				dayDiff := eEpochDay - sEpochDay
				secFromDays, err := multiplyExact(dayDiff, 86400)
				if err != nil {
					return err
				}
				sH, sMin, sS, sN := getLocalTimeHMSN(sTime)
				eH, eMin, eS, eN := getLocalTimeHMSN(eTime)
				sSecOfDay := sH*3600 + sMin*60 + sS
				eSecOfDay := eH*3600 + eMin*60 + eS
				secDiff := eSecOfDay - sSecOfDay
				totalSec, err := addExact(secFromDays, secDiff)
				if err != nil {
					return err
				}
				nanosDiff := eN - sN
				return durationOfSecondsNanos([]any{totalSec, nanosDiff})
			}
		}
	}

	// ZonedDateTime / OffsetDateTime
	var sSec, eSec int64
	var sNano, eNano int64
	var okInstant1, okInstant2 bool

	if strings.Contains(startClass, "ZonedDateTime") {
		if inst, ok := zonedDateTimeToInstant([]any{start}).(*object.Object); ok {
			sSec, sNano, okInstant1 = getInstantSecNano(inst)
		}
	} else if strings.Contains(startClass, "OffsetDateTime") {
		if inst, ok := offsetDateTimeToInstant([]any{start}).(*object.Object); ok {
			sSec, sNano, okInstant1 = getInstantSecNano(inst)
		}
	}

	if strings.Contains(endClass, "ZonedDateTime") {
		if inst, ok := zonedDateTimeToInstant([]any{end}).(*object.Object); ok {
			eSec, eNano, okInstant2 = getInstantSecNano(inst)
		}
	} else if strings.Contains(endClass, "OffsetDateTime") {
		if inst, ok := offsetDateTimeToInstant([]any{end}).(*object.Object); ok {
			eSec, eNano, okInstant2 = getInstantSecNano(inst)
		}
	}

	if okInstant1 && okInstant2 {
		secs, err := addExact(eSec, -sSec)
		if err != nil {
			return err
		}
		nanos := eNano - sNano
		return durationOfSecondsNanos([]any{secs, nanos})
	}

	// General fallback: try start.until(end, SECONDS) and start.until(end, NANOS)
	nanosUnit := ghelpers.CreateTemporalUnit("NANOS")
	secsUnit := ghelpers.CreateTemporalUnit("SECONDS")
	resNanos := temporalUntil([]any{start, end, nanosUnit})
	if nVal, ok := getInt64Val(resNanos); ok && resNanos != nil {
		if _, isErr := resNanos.(*ghelpers.GErrBlk); !isErr {
			return durationOfNanos([]any{nVal})
		}
	}

	resSecs := temporalUntil([]any{start, end, secsUnit})
	if sVal, ok := getInt64Val(resSecs); ok && resSecs != nil {
		if _, isErr := resSecs.(*ghelpers.GErrBlk); !isErr {
			return durationOfSeconds([]any{sVal})
		}
	}

	return ghelpers.GetGErrBlk(excNames.DateTimeException, "Duration.between: unable to calculate duration between temporal objects")
}

// Duration.from(TemporalAmount amount)
func durationFrom(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Duration.from: null argument")
	}
	amount := params[0].(*object.Object)
	className := object.GoStringFromStringPoolIndex(amount.KlassName)
	if className == classNameDuration {
		return amount
	}
	if sec, nano, ok := getDurationSecondsNanos(amount); ok {
		return createDuration(sec, nano)
	}

	unitsObj := temporalAmountGetUnits([]any{amount})
	if listObj, ok := unitsObj.(*object.Object); ok {
		if valFld, ok := listObj.FieldTable["value"]; ok {
			if unitsList, ok := valFld.Fvalue.([]any); ok {
				var durObj = createDuration(0, 0)
				for _, u := range unitsList {
					if uObj, ok := u.(*object.Object); ok {
						val := temporalAmountGet([]any{amount, uObj})
						if v, ok := getInt64Val(val); ok {
							unitName := getTemporalUnitName(uObj)
							switch unitName {
							case "SECONDS", "Seconds":
								durObj = durationPlusSeconds([]any{durObj, v}).(*object.Object)
							case "NANOS", "Nanos":
								durObj = durationPlusNanos([]any{durObj, v}).(*object.Object)
							default:
								return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Duration.from: unsupported unit %s", unitName))
							}
						} else {
							return val
						}
					}
				}
				return durObj
			}
		}
	}

	return ghelpers.GetGErrBlk(excNames.DateTimeException, "Duration.from: unable to obtain Duration from TemporalAmount")
}

// Duration.get(TemporalUnit unit)
func durationGet(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Duration.get: null argument")
	}
	self := params[0].(*object.Object)
	unit := params[1].(*object.Object)

	unitName := getTemporalUnitName(unit)
	switch unitName {
	case "SECONDS", "Seconds":
		return self.FieldTable["seconds"].Fvalue.(int64)
	case "NANOS", "Nanos":
		return self.FieldTable["nanos"].Fvalue.(int64)
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported unit: %s", unitName))
	}
}

// Duration.getUnits()
func durationGetUnits(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Duration.getUnits: self is null")
	}
	ensureChronoUnitInited()
	secUnit := ghelpers.CreateTemporalUnit("SECONDS")
	nanosUnit := ghelpers.CreateTemporalUnit("NANOS")

	listObj := object.MakeEmptyObjectWithClassName(&classNameArrayList)
	listObj.FieldTable["value"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: []any{secUnit, nanosUnit},
	}
	return listObj
}

// Duration.of(long amount, TemporalUnit unit)
func durationOf(params []any) any {
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Duration.of: unit is null")
	}
	amount, ok := getInt64Val(params[0])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Duration.of: invalid amount")
	}
	unit := params[1].(*object.Object)
	unitName := getTemporalUnitName(unit)

	switch unitName {
	case "NANOS", "Nanos":
		return durationOfNanos([]any{amount})
	case "MICROS", "Micros":
		sec := amount / 1000000
		nano := (amount % 1000000) * 1000
		if nano < 0 {
			nano += nanosPerSecond
			sec--
		}
		return durationOfSecondsNanos([]any{sec, nano})
	case "MILLIS", "Millis":
		return durationOfMillis([]any{amount})
	case "SECONDS", "Seconds":
		return durationOfSeconds([]any{amount})
	case "MINUTES", "Minutes":
		return durationOfMinutes([]any{amount})
	case "HOURS", "Hours":
		return durationOfHours([]any{amount})
	case "HALF_DAYS", "HalfDays":
		secs, err := multiplyExact(amount, 43200)
		if err != nil {
			return err
		}
		return durationOfSeconds([]any{secs})
	case "DAYS", "Days":
		return durationOfDays([]any{amount})
	default:
		if durFld, exists := unit.FieldTable["duration"]; exists {
			if durObj, ok := durFld.Fvalue.(*object.Object); ok && durObj != nil {
				if estFld, ok := unit.FieldTable["estimated"]; ok {
					if estVal, ok := getInt64Val(estFld.Fvalue); ok && estVal != 0 {
						return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unit must not have an estimated duration")
					}
				}
				return durationMultipliedBy([]any{durObj, amount})
			}
		}
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported unit: %s", unitName))
	}
}

// Duration.plus(long amountToAdd, TemporalUnit unit)
func durationPlusUnit(params []any) any {
	if len(params) < 3 || params[0] == nil || params[2] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Duration.plus: null argument")
	}
	self := params[0].(*object.Object)
	amountToAdd, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Duration.plus: invalid amount")
	}
	unit := params[2].(*object.Object)
	if amountToAdd == 0 {
		return self
	}
	toAdd := durationOf([]any{amountToAdd, unit})
	if err, ok := toAdd.(*ghelpers.GErrBlk); ok {
		return err
	}
	return durationPlus([]any{self, toAdd})
}

// Duration.minus(long amountToSubtract, TemporalUnit unit)
func durationMinusUnit(params []any) any {
	if len(params) < 3 || params[0] == nil || params[2] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Duration.minus: null argument")
	}
	self := params[0].(*object.Object)
	amountToSubtract, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Duration.minus: invalid amount")
	}
	unit := params[2].(*object.Object)
	if amountToSubtract == math.MinInt64 {
		p1 := durationPlusUnit([]any{self, int64(math.MaxInt64), unit})
		if err, ok := p1.(*ghelpers.GErrBlk); ok {
			return err
		}
		return durationPlusUnit([]any{p1, int64(1), unit})
	}
	return durationPlusUnit([]any{self, -amountToSubtract, unit})
}

// Duration.truncatedTo(TemporalUnit unit)
func durationTruncatedTo(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Duration.truncatedTo: null argument")
	}
	self := params[0].(*object.Object)
	unit := params[1].(*object.Object)
	unitName := getTemporalUnitName(unit)

	var unitNanos int64
	switch unitName {
	case "NANOS", "Nanos":
		return self
	case "MICROS", "Micros":
		unitNanos = 1000
	case "MILLIS", "Millis":
		unitNanos = 1_000_000
	case "SECONDS", "Seconds":
		unitNanos = 1_000_000_000
	case "MINUTES", "Minutes":
		unitNanos = 60 * 1_000_000_000
	case "HOURS", "Hours":
		unitNanos = 3600 * 1_000_000_000
	case "HALF_DAYS", "HalfDays":
		unitNanos = 43200 * 1_000_000_000
	case "DAYS", "Days":
		unitNanos = 86400 * 1_000_000_000
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unit is not supported for truncation: %s", unitName))
	}

	const nanosPerDay = int64(86400 * 1_000_000_000)
	if nanosPerDay%unitNanos != 0 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unit does not divide evenly into a day")
	}

	seconds := self.FieldTable["seconds"].Fvalue.(int64)
	nanos := self.FieldTable["nanos"].Fvalue.(int64)

	nod := (seconds%86400)*1_000_000_000 + nanos
	resultNod := (nod / unitNanos) * unitNanos
	return durationPlusNanos([]any{self, resultNod - nod})
}

// Duration.addTo(Temporal temporal)
func durationAddTo(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Duration.addTo: null argument")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)

	seconds := self.FieldTable["seconds"].Fvalue.(int64)
	nanos := self.FieldTable["nanos"].Fvalue.(int64)

	res := temporal
	if seconds != 0 {
		secUnit := ghelpers.CreateTemporalUnit("SECONDS")
		p := temporalPlusUnit([]any{res, seconds, secUnit})
		if err, ok := p.(*ghelpers.GErrBlk); ok {
			return err
		}
		if obj, ok := p.(*object.Object); ok {
			res = obj
		}
	}
	if nanos != 0 {
		nanoUnit := ghelpers.CreateTemporalUnit("NANOS")
		p := temporalPlusUnit([]any{res, nanos, nanoUnit})
		if err, ok := p.(*ghelpers.GErrBlk); ok {
			return err
		}
		if obj, ok := p.(*object.Object); ok {
			res = obj
		}
	}
	return res
}

// Duration.subtractFrom(Temporal temporal)
func durationSubtractFrom(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Duration.subtractFrom: null argument")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)

	seconds := self.FieldTable["seconds"].Fvalue.(int64)
	nanos := self.FieldTable["nanos"].Fvalue.(int64)

	res := temporal
	if seconds != 0 {
		secUnit := ghelpers.CreateTemporalUnit("SECONDS")
		p := temporalMinusUnit([]any{res, seconds, secUnit})
		if err, ok := p.(*ghelpers.GErrBlk); ok {
			return err
		}
		if obj, ok := p.(*object.Object); ok {
			res = obj
		}
	}
	if nanos != 0 {
		nanoUnit := ghelpers.CreateTemporalUnit("NANOS")
		p := temporalMinusUnit([]any{res, nanos, nanoUnit})
		if err, ok := p.(*ghelpers.GErrBlk); ok {
			return err
		}
		if obj, ok := p.(*object.Object); ok {
			res = obj
		}
	}
	return res
}
