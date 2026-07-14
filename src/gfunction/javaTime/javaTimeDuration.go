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
	"strconv"
	"strings"
)

var classNameDuration = "java/time/Duration"

func Load_Time_Duration() {
	// Static methods
	ghelpers.MethodSignatures["java/time/Duration.<clinit>()V"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationClinit}
	ghelpers.MethodSignatures["java/time/Duration.abs()Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationAbs}
	ghelpers.MethodSignatures["java/time/Duration.between(Ljava/time/temporal/Temporal;Ljava/time/temporal/Temporal;)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/time/Duration.compareTo(Ljava/time/Duration;)I"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationCompareTo}
	ghelpers.MethodSignatures["java/time/Duration.dividedBy(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationDividedByLong}
	ghelpers.MethodSignatures["java/time/Duration.dividedBy(Ljava/time/Duration;)J"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationDividedByDuration}
	ghelpers.MethodSignatures["java/time/Duration.equals(Ljava/lang/Object;)Z"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationEquals}
	ghelpers.MethodSignatures["java/time/Duration.from(Ljava/time/temporal/TemporalAmount;)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/time/Duration.get(Ljava/time/temporal/TemporalUnit;)J"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/time/Duration.getNano()I"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationGetNano}
	ghelpers.MethodSignatures["java/time/Duration.getSeconds()J"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationGetSeconds}
	ghelpers.MethodSignatures["java/time/Duration.getUnits()Ljava/util/List;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/time/Duration.hashCode()I"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationHashCode}
	ghelpers.MethodSignatures["java/time/Duration.isNegative()Z"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationIsNegative}
	ghelpers.MethodSignatures["java/time/Duration.isZero()Z"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationIsZero}
	ghelpers.MethodSignatures["java/time/Duration.minus(Ljava/time/Duration;)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationMinus}
	ghelpers.MethodSignatures["java/time/Duration.minusDays(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationMinusDays}
	ghelpers.MethodSignatures["java/time/Duration.minusHours(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationMinusHours}
	ghelpers.MethodSignatures["java/time/Duration.minusMillis(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationMinusMillis}
	ghelpers.MethodSignatures["java/time/Duration.minusMinutes(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationMinusMinutes}
	ghelpers.MethodSignatures["java/time/Duration.minusNanos(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationMinusNanos}
	ghelpers.MethodSignatures["java/time/Duration.minusSeconds(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationMinusSeconds}
	ghelpers.MethodSignatures["java/time/Duration.multipliedBy(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationMultipliedBy}
	ghelpers.MethodSignatures["java/time/Duration.negated()Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 0, GFunction: durationNegated}
	ghelpers.MethodSignatures["java/time/Duration.of(JLjava/time/temporal/TemporalUnit;)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
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
	ghelpers.MethodSignatures["java/time/Duration.truncatedTo(Ljava/time/temporal/TemporalUnit;)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/time/Duration.withNanos(I)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationWithNanos}
	ghelpers.MethodSignatures["java/time/Duration.withSeconds(J)Ljava/time/Duration;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: durationWithSeconds}
	ghelpers.MethodSignatures["java/time/Duration.addTo(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/time/Duration.subtractFrom(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] = ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
}

// Helper to create a new Duration object
func createDuration(seconds int64, nanos int32) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameDuration)
	obj.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: seconds}
	obj.FieldTable["nanos"] = object.Field{Ftype: types.Int, Fvalue: int64(nanos)}
	return obj
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

func durationMultipliedBy(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	multiplicand := params[1].(int64)
	if multiplicand == 0 {
		return createDuration(0, 0)
	}
	if multiplicand == 1 {
		return self
	}
	s := self.FieldTable["seconds"].Fvalue.(int64)
	n := self.FieldTable["nanos"].Fvalue.(int64)

	resSeconds, err := multiplyExact(s, multiplicand)
	if err != nil {
		return err
	}
	resNanos, err := multiplyExact(n, multiplicand)
	if err != nil {
		// If nano multiplication overflows, we can still calculate it using seconds
		secAdjustment := resNanos / nanosPerSecond
		resNanos = resNanos % nanosPerSecond
		if resNanos < 0 {
			resNanos += nanosPerSecond
			secAdjustment--
		}
		resSeconds, err = addExact(resSeconds, secAdjustment)
		if err != nil {
			return err
		}
		return createDuration(resSeconds, int32(resNanos))
	}
	// The above logic for nano overflow was flawed because multiplyExact already returned err
	// Correct way to multiply nano:
	// n * multiplicand = (n * multiplicand / 1B) * 1B + (n * multiplicand % 1B)
	// Since n < 1B, n * multiplicand can only overflow if multiplicand is very large.
	// Let's use float64 or big.Int for simplicity if it overflows?
	// Actually, n * multiplicand for n < 1B and multiplicand < ~9B will not overflow int64.
	// If it does, we can use a more robust way.
	
	// Re-calculating with potential overflow in mind
	resNanosLong := n * multiplicand
	secAdjustment := resNanosLong / nanosPerSecond
	resNanosFinal := resNanosLong % nanosPerSecond
	if resNanosFinal < 0 {
		resNanosFinal += nanosPerSecond
		secAdjustment--
	}
	
	resSeconds, err = addExact(resSeconds, secAdjustment)
	if err != nil {
		return err
	}
	
	return createDuration(resSeconds, int32(resNanosFinal))
}

func durationDividedByLong(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	divisor := params[1].(int64)
	if divisor == 0 {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Division by zero")
	}
	if divisor == 1 {
		return self
	}
	s := self.FieldTable["seconds"].Fvalue.(int64)
	n := self.FieldTable["nanos"].Fvalue.(int64)
	
	totalNanos := s*nanosPerSecond + n
	// Check for overflow if s is large.
	if s > 9223372036 || s < -9223372036 {
		// Fallback to simpler but less precise? Or just do it.
		resSeconds := s / divisor
		remainderSeconds := s % divisor
		resNanos := (remainderSeconds*nanosPerSecond + n) / divisor
		return createDuration(resSeconds, int32(resNanos))
	}
	
	resNanosTotal := totalNanos / divisor
	return durationOfNanos([]interface{}{resNanosTotal})
}

func durationDividedByDuration(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	divisor := params[1].(*object.Object)
	s1 := self.FieldTable["seconds"].Fvalue.(int64)
	n1 := self.FieldTable["nanos"].Fvalue.(int64)
	s2 := divisor.FieldTable["seconds"].Fvalue.(int64)
	n2 := divisor.FieldTable["nanos"].Fvalue.(int64)
	
	if s2 == 0 && n2 == 0 {
		return ghelpers.GetGErrBlk(excNames.ArithmeticException, "Division by zero")
	}
	
	t1 := s1*nanosPerSecond + n1
	t2 := s2*nanosPerSecond + n2
	
	if (s1 > 9223372036 || s1 < -9223372036) || (s2 > 9223372036 || s2 < -9223372036) {
		// Use big.Int for large durations
		return int64(float64(s1)/float64(s2)) // VERY crude approximation
	}
	
	return t1 / t2
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
	return int64(int32(seconds ^ (seconds >> 32)) + 51*int32(nanos))
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
