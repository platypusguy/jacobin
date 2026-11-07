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
	"strconv"
	"strings"
	"sync"
	"time"
)

var classNameLocalTime = "java/time/LocalTime"

const (
	hoursPerDay      int64 = 24
	minutesPerHour   int64 = 60
	minutesPerDay    int64 = 1440
	secondsPerMinute int64 = 60
	secondsPerHour   int64 = 3600
	secondsPerDay    int64 = 86400
	millisPerDay     int64 = 86400000
	microsPerDay     int64 = 86400000000
	nanosPerMinute   int64 = 60000000000
	nanosPerHour     int64 = 3600000000000
	nanosPerDay      int64 = 86400000000000
	nanosPerHalfDay  int64 = 43200000000000
)

var localTimeMutex sync.Mutex
var localTimeInitialized bool = false

var (
	localTimeMin      *object.Object
	localTimeMax      *object.Object
	localTimeMidnight *object.Object
	localTimeNoon     *object.Object
)

func ensureLocalTimeInited() {
	localTimeMutex.Lock()
	defer localTimeMutex.Unlock()
	if localTimeInitialized {
		return
	}

	localTimeMin = createLocalTimeInternal(0, 0, 0, 0)
	localTimeMax = createLocalTimeInternal(23, 59, 59, 999999999)
	localTimeMidnight = createLocalTimeInternal(0, 0, 0, 0)
	localTimeNoon = createLocalTimeInternal(12, 0, 0, 0)

	_ = statics.AddStatic(classNameLocalTime+".MIN", statics.Static{
		Type:  "Ljava/time/LocalTime;",
		Value: localTimeMin,
	})
	_ = statics.AddStatic(classNameLocalTime+".MAX", statics.Static{
		Type:  "Ljava/time/LocalTime;",
		Value: localTimeMax,
	})
	_ = statics.AddStatic(classNameLocalTime+".MIDNIGHT", statics.Static{
		Type:  "Ljava/time/LocalTime;",
		Value: localTimeMidnight,
	})
	_ = statics.AddStatic(classNameLocalTime+".NOON", statics.Static{
		Type:  "Ljava/time/LocalTime;",
		Value: localTimeNoon,
	})
	_ = statics.AddStatic(classNameLocalTime+".HOURS_PER_DAY", statics.Static{
		Type:  "I",
		Value: int64(24),
	})
	_ = statics.AddStatic(classNameLocalTime+".MINUTES_PER_HOUR", statics.Static{
		Type:  "I",
		Value: int64(60),
	})
	_ = statics.AddStatic(classNameLocalTime+".MINUTES_PER_DAY", statics.Static{
		Type:  "I",
		Value: int64(1440),
	})
	_ = statics.AddStatic(classNameLocalTime+".SECONDS_PER_MINUTE", statics.Static{
		Type:  "I",
		Value: int64(60),
	})
	_ = statics.AddStatic(classNameLocalTime+".SECONDS_PER_HOUR", statics.Static{
		Type:  "I",
		Value: int64(3600),
	})
	_ = statics.AddStatic(classNameLocalTime+".SECONDS_PER_DAY", statics.Static{
		Type:  "I",
		Value: int64(86400),
	})
	_ = statics.AddStatic(classNameLocalTime+".MILLIS_PER_DAY", statics.Static{
		Type:  "J",
		Value: int64(86400000),
	})
	_ = statics.AddStatic(classNameLocalTime+".MICROS_PER_DAY", statics.Static{
		Type:  "J",
		Value: int64(86400000000),
	})
	_ = statics.AddStatic(classNameLocalTime+".NANOS_PER_SECOND", statics.Static{
		Type:  "J",
		Value: int64(1000000000),
	})
	_ = statics.AddStatic(classNameLocalTime+".NANOS_PER_MINUTE", statics.Static{
		Type:  "J",
		Value: int64(60000000000),
	})
	_ = statics.AddStatic(classNameLocalTime+".NANOS_PER_HOUR", statics.Static{
		Type:  "J",
		Value: int64(3600000000000),
	})
	_ = statics.AddStatic(classNameLocalTime+".NANOS_PER_DAY", statics.Static{
		Type:  "J",
		Value: int64(86400000000000),
	})

	localTimeInitialized = true
}

func createLocalTimeInternal(hour, minute, second, nano int64) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameLocalTime)
	obj.FieldTable["hour"] = object.Field{Ftype: types.Int, Fvalue: hour}
	obj.FieldTable["minute"] = object.Field{Ftype: types.Int, Fvalue: minute}
	obj.FieldTable["second"] = object.Field{Ftype: types.Int, Fvalue: second}
	obj.FieldTable["nano"] = object.Field{Ftype: types.Int, Fvalue: nano}
	obj.FieldTable["HOUR_OF_DAY"] = object.Field{Ftype: types.Long, Fvalue: hour}
	obj.FieldTable["MINUTE_OF_HOUR"] = object.Field{Ftype: types.Long, Fvalue: minute}
	obj.FieldTable["SECOND_OF_MINUTE"] = object.Field{Ftype: types.Long, Fvalue: second}
	obj.FieldTable["NANO_OF_SECOND"] = object.Field{Ftype: types.Long, Fvalue: nano}
	nanoOfDay := hour*nanosPerHour + minute*nanosPerMinute + second*nanosPerSecond + nano
	obj.FieldTable["nanoOfDay"] = object.Field{Ftype: types.Long, Fvalue: nanoOfDay}
	return obj
}

// createLocalTime creates a LocalTime object, ensuring statics are initialized.
func createLocalTime(hour, minute, second, nano int64) *object.Object {
	ensureLocalTimeInited()
	return createLocalTimeInternal(hour, minute, second, nano)
}

// CreateLocalTime is an exported helper for other packages.
func CreateLocalTime(hour, minute, second, nano int64) *object.Object {
	return createLocalTime(hour, minute, second, nano)
}

func getLocalTimeHMSN(obj *object.Object) (int64, int64, int64, int64) {
	if obj == nil || obj.FieldTable == nil {
		return 0, 0, 0, 0
	}
	var hour, minute, second, nano int64
	var found bool
	if hFld, ok := obj.FieldTable["hour"]; ok {
		if h, ok2 := getInt64Val(hFld.Fvalue); ok2 {
			hour = h
			found = true
		}
	} else if hFld, ok := obj.FieldTable["HOUR_OF_DAY"]; ok {
		if h, ok2 := getInt64Val(hFld.Fvalue); ok2 {
			hour = h
			found = true
		}
	}

	if mFld, ok := obj.FieldTable["minute"]; ok {
		if m, ok2 := getInt64Val(mFld.Fvalue); ok2 {
			minute = m
		}
	} else if mFld, ok := obj.FieldTable["MINUTE_OF_HOUR"]; ok {
		if m, ok2 := getInt64Val(mFld.Fvalue); ok2 {
			minute = m
		}
	}

	if sFld, ok := obj.FieldTable["second"]; ok {
		if s, ok2 := getInt64Val(sFld.Fvalue); ok2 {
			second = s
		}
	} else if sFld, ok := obj.FieldTable["SECOND_OF_MINUTE"]; ok {
		if s, ok2 := getInt64Val(sFld.Fvalue); ok2 {
			second = s
		}
	}

	if nFld, ok := obj.FieldTable["nano"]; ok {
		if n, ok2 := getInt64Val(nFld.Fvalue); ok2 {
			nano = n
		}
	} else if nFld, ok := obj.FieldTable["NANO_OF_SECOND"]; ok {
		if n, ok2 := getInt64Val(nFld.Fvalue); ok2 {
			nano = n
		}
	}

	if !found {
		if nanoOfDayFld, ok := obj.FieldTable["nanoOfDay"]; ok {
			if nod, ok2 := getInt64Val(nanoOfDayFld.Fvalue); ok2 {
				hour = (nod / nanosPerHour) % 24
				nod %= nanosPerHour
				minute = nod / nanosPerMinute
				nod %= nanosPerMinute
				second = nod / nanosPerSecond
				nano = nod % nanosPerSecond
			}
		}
	}
	return hour, minute, second, nano
}

func getLocalTimeNanoOfDay(obj *object.Object) int64 {
	h, m, s, n := getLocalTimeHMSN(obj)
	return h*nanosPerHour + m*nanosPerMinute + s*nanosPerSecond + n
}

func compareLocalTime(t1, t2 *object.Object) int64 {
	nod1 := getLocalTimeNanoOfDay(t1)
	nod2 := getLocalTimeNanoOfDay(t2)
	if nod1 < nod2 {
		return int64(-1)
	} else if nod1 > nod2 {
		return int64(1)
	}
	return int64(0)
}

func validateLocalTime(hour, minute, second, nano int64) error {
	if hour < 0 || hour > 23 {
		return fmt.Errorf("Invalid value for HourOfDay (valid values 0 - 23): %d", hour)
	}
	if minute < 0 || minute > 59 {
		return fmt.Errorf("Invalid value for MinuteOfHour (valid values 0 - 59): %d", minute)
	}
	if second < 0 || second > 59 {
		return fmt.Errorf("Invalid value for SecondOfMinute (valid values 0 - 59): %d", second)
	}
	if nano < 0 || nano > 999999999 {
		return fmt.Errorf("Invalid value for NanoOfSecond (valid values 0 - 999999999): %d", nano)
	}
	return nil
}

func formatLocalTimeString(h, m, s, n int64) string {
	if s == 0 && n == 0 {
		return fmt.Sprintf("%02d:%02d", h, m)
	}
	if n == 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	if n%1000000 == 0 {
		return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, n/1000000)
	} else if n%1000 == 0 {
		return fmt.Sprintf("%02d:%02d:%02d.%06d", h, m, s, n/1000)
	}
	return fmt.Sprintf("%02d:%02d:%02d.%09d", h, m, s, n)
}

func parseLocalTimeString(text string) (int64, int64, int64, int64, error) {
	text = strings.TrimPrefix(text, "T")
	parts := strings.Split(text, ":")
	if len(parts) == 2 {
		h, errH := strconv.ParseInt(parts[0], 10, 64)
		m, errM := strconv.ParseInt(parts[1], 10, 64)
		if errH != nil || errM != nil || len(parts[0]) != 2 || len(parts[1]) != 2 {
			return 0, 0, 0, 0, fmt.Errorf("text could not be parsed: %s", text)
		}
		if err := validateLocalTime(h, m, 0, 0); err != nil {
			return 0, 0, 0, 0, err
		}
		return h, m, 0, 0, nil
	} else if len(parts) == 3 {
		h, errH := strconv.ParseInt(parts[0], 10, 64)
		m, errM := strconv.ParseInt(parts[1], 10, 64)
		if errH != nil || errM != nil || len(parts[0]) != 2 || len(parts[1]) != 2 {
			return 0, 0, 0, 0, fmt.Errorf("text could not be parsed: %s", text)
		}
		secPart := parts[2]
		var s, n int64
		if dotIdx := strings.Index(secPart, "."); dotIdx >= 0 {
			sStr := secPart[:dotIdx]
			nStr := secPart[dotIdx+1:]
			if len(sStr) != 2 || len(nStr) == 0 || len(nStr) > 9 {
				return 0, 0, 0, 0, fmt.Errorf("text could not be parsed: %s", text)
			}
			var errS error
			s, errS = strconv.ParseInt(sStr, 10, 64)
			if errS != nil {
				return 0, 0, 0, 0, fmt.Errorf("text could not be parsed: %s", text)
			}
			for len(nStr) < 9 {
				nStr += "0"
			}
			var errN error
			n, errN = strconv.ParseInt(nStr, 10, 64)
			if errN != nil {
				return 0, 0, 0, 0, fmt.Errorf("text could not be parsed: %s", text)
			}
		} else {
			if len(secPart) != 2 {
				return 0, 0, 0, 0, fmt.Errorf("text could not be parsed: %s", text)
			}
			var errS error
			s, errS = strconv.ParseInt(secPart, 10, 64)
			if errS != nil {
				return 0, 0, 0, 0, fmt.Errorf("text could not be parsed: %s", text)
			}
		}
		if err := validateLocalTime(h, m, s, n); err != nil {
			return 0, 0, 0, 0, err
		}
		return h, m, s, n, nil
	}
	return 0, 0, 0, 0, fmt.Errorf("text could not be parsed: %s", text)
}

func ofNanoOfDay(nanoOfDay int64) *object.Object {
	hours := (nanoOfDay / nanosPerHour) % 24
	nanoOfDay %= nanosPerHour
	minutes := nanoOfDay / nanosPerMinute
	nanoOfDay %= nanosPerMinute
	seconds := nanoOfDay / nanosPerSecond
	nanos := nanoOfDay % nanosPerSecond
	return createLocalTime(hours, minutes, seconds, nanos)
}

func localTimePlusNanosInternal(h, m, s, n, nanosToAdd int64) *object.Object {
	if nanosToAdd == 0 {
		return createLocalTime(h, m, s, n)
	}
	nod := h*nanosPerHour + m*nanosPerMinute + s*nanosPerSecond + n
	newNod := ((nod+nanosToAdd)%nanosPerDay + nanosPerDay) % nanosPerDay
	return ofNanoOfDay(newNod)
}

func localTimePlusSecondsInternal(h, m, s, n, secondsToAdd int64) *object.Object {
	if secondsToAdd == 0 {
		return createLocalTime(h, m, s, n)
	}
	totSec := h*secondsPerHour + m*secondsPerMinute + s
	newSec := ((totSec+secondsToAdd)%secondsPerDay + secondsPerDay) % secondsPerDay
	newH := newSec / secondsPerHour
	newM := (newSec % secondsPerHour) / secondsPerMinute
	newS := newSec % secondsPerMinute
	return createLocalTime(newH, newM, newS, n)
}

func localTimePlusMinutesInternal(h, m, s, n, minutesToAdd int64) *object.Object {
	if minutesToAdd == 0 {
		return createLocalTime(h, m, s, n)
	}
	totMin := h*minutesPerHour + m
	newMin := ((totMin+minutesToAdd)%minutesPerDay + minutesPerDay) % minutesPerDay
	newH := newMin / minutesPerHour
	newM := newMin % minutesPerHour
	return createLocalTime(newH, newM, s, n)
}

func localTimePlusHoursInternal(h, m, s, n, hoursToAdd int64) *object.Object {
	if hoursToAdd == 0 {
		return createLocalTime(h, m, s, n)
	}
	newH := ((h+hoursToAdd)%hoursPerDay + hoursPerDay) % hoursPerDay
	return createLocalTime(newH, m, s, n)
}

// Load_Time_LocalTime registers all java/time/LocalTime methods.
func Load_Time_LocalTime() {
	ghelpers.MethodSignatures["java/time/LocalTime.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localTimeClinit,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.adjustInto(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeAdjustInto,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.atDate(Ljava/time/LocalDate;)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeAtDate,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.atOffset(Ljava/time/ZoneOffset;)Ljava/time/OffsetTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeAtOffset,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeCompareTo,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.compareTo(Ljava/time/LocalTime;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeCompareTo,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeEquals,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.format(Ljava/time/format/DateTimeFormatter;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeFormat,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.from(Ljava/time/temporal/TemporalAccessor;)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeFrom,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.get(Ljava/time/temporal/TemporalField;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeGet,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.getHour()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localTimeGetHour,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.getLong(Ljava/time/temporal/TemporalField;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeGetLong,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.getMinute()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localTimeGetMinute,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.getNano()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localTimeGetNano,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.getSecond()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localTimeGetSecond,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localTimeHashCode,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.isAfter(Ljava/time/LocalTime;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeIsAfter,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.isBefore(Ljava/time/LocalTime;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeIsBefore,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.isSupported(Ljava/time/temporal/TemporalField;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeIsSupportedField,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.isSupported(Ljava/time/temporal/TemporalUnit;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeIsSupportedUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  localTimeMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.minus(JLjava/time/temporal/TemporalUnit;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  localTimeMinusUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.minus(Ljava/time/temporal/TemporalAmount;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeMinusAmount,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.minusHours(J)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localTimeMinusHours,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.minusMinutes(J)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localTimeMinusMinutes,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.minusNanos(J)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localTimeMinusNanos,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.minusSeconds(J)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localTimeMinusSeconds,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.now()Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localTimeNow,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.now(Ljava/time/Clock;)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeNowClock,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.now(Ljava/time/ZoneId;)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeNowZone,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.of(II)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localTimeOfHM,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.of(III)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  localTimeOfHMS,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.of(IIII)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  localTimeOfHMSN,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.ofInstant(Ljava/time/Instant;Ljava/time/ZoneId;)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localTimeOfInstant,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.ofNanoOfDay(J)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localTimeOfNanoOfDay,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.ofSecondOfDay(J)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localTimeOfSecondOfDay,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.parse(Ljava/lang/CharSequence;)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeParse,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.parse(Ljava/lang/CharSequence;Ljava/time/format/DateTimeFormatter;)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localTimeParseFormatter,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  localTimePlusUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.plus(JLjava/time/temporal/TemporalUnit;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  localTimePlusUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimePlusAmount,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.plus(Ljava/time/temporal/TemporalAmount;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimePlusAmount,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.plusHours(J)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localTimePlusHours,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.plusMinutes(J)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localTimePlusMinutes,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.plusNanos(J)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localTimePlusNanos,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.plusSeconds(J)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localTimePlusSeconds,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.query(Ljava/time/temporal/TemporalQuery;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeQuery,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.range(Ljava/time/temporal/TemporalField;)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeRange,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.toEpochSecond(Ljava/time/LocalDate;Ljava/time/ZoneOffset;)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localTimeToEpochSecond,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.toNanoOfDay()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localTimeToNanoOfDay,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.toSecondOfDay()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localTimeToSecondOfDay,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localTimeToString,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.truncatedTo(Ljava/time/temporal/TemporalUnit;)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeTruncatedTo,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.until(Ljava/time/temporal/Temporal;Ljava/time/temporal/TemporalUnit;)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localTimeUntilUnit,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.with(Ljava/time/temporal/TemporalAdjuster;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeWithAdjuster,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.with(Ljava/time/temporal/TemporalField;J)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  localTimeWithField,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.with(Ljava/time/temporal/TemporalField;J)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  localTimeWithField,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.withHour(I)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeWithHour,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.withMinute(I)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeWithMinute,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.withNano(I)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeWithNano,
		}

	ghelpers.MethodSignatures["java/time/LocalTime.withSecond(I)Ljava/time/LocalTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localTimeWithSecond,
		}
}

func localTimeClinit([]any) any {
	ensureLocalTimeInited()
	return nil
}

func localTimeNow([]any) any {
	now := time.Now()
	return CreateLocalTime(int64(now.Hour()), int64(now.Minute()), int64(now.Second()), int64(now.Nanosecond()))
}

func localTimeNowClock(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.now: clock is null")
	}
	now := time.Now()
	return CreateLocalTime(int64(now.Hour()), int64(now.Minute()), int64(now.Second()), int64(now.Nanosecond()))
}

func localTimeNowZone(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.now: zone is null")
	}
	now := time.Now()
	return CreateLocalTime(int64(now.Hour()), int64(now.Minute()), int64(now.Second()), int64(now.Nanosecond()))
}

func localTimeOfHM(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.of: missing arguments")
	}
	h, okH := getInt64Val(params[0])
	m, okM := getInt64Val(params[1])
	if !okH || !okM {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.of: invalid arguments")
	}
	if err := validateLocalTime(h, m, 0, 0); err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
	}
	return createLocalTime(h, m, 0, 0)
}

func localTimeOfHMS(params []any) any {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.of: missing arguments")
	}
	h, okH := getInt64Val(params[0])
	m, okM := getInt64Val(params[1])
	s, okS := getInt64Val(params[2])
	if !okH || !okM || !okS {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.of: invalid arguments")
	}
	if err := validateLocalTime(h, m, s, 0); err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
	}
	return createLocalTime(h, m, s, 0)
}

func localTimeOfHMSN(params []any) any {
	if len(params) < 4 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.of: missing arguments")
	}
	h, okH := getInt64Val(params[0])
	m, okM := getInt64Val(params[1])
	s, okS := getInt64Val(params[2])
	n, okN := getInt64Val(params[3])
	if !okH || !okM || !okS || !okN {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.of: invalid arguments")
	}
	if err := validateLocalTime(h, m, s, n); err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
	}
	return createLocalTime(h, m, s, n)
}

func localTimeOfInstant(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.ofInstant: null argument")
	}
	instantObj := params[0].(*object.Object)
	zoneObj := params[1].(*object.Object)

	var epochSec int64 = 0
	var nanoVal int64 = 0
	if fld, ok := instantObj.FieldTable["seconds"]; ok {
		if s, ok2 := getInt64Val(fld.Fvalue); ok2 {
			epochSec = s
		}
	}
	if fld, ok := instantObj.FieldTable["nanos"]; ok {
		if n, ok2 := getInt64Val(fld.Fvalue); ok2 {
			nanoVal = n
		}
	}

	offsetSec := getZoneOffsetTotalSeconds(zoneObj)
	secOfDay := ((epochSec%secondsPerDay+secondsPerDay)%secondsPerDay + offsetSec) % secondsPerDay
	if secOfDay < 0 {
		secOfDay += secondsPerDay
	}
	nod := secOfDay*nanosPerSecond + nanoVal
	return ofNanoOfDay(nod)
}

func localTimeOfSecondOfDay(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.ofSecondOfDay: missing argument")
	}
	secOfDay, ok := getInt64Val(params[0])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.ofSecondOfDay: invalid argument")
	}
	if secOfDay < 0 || secOfDay >= secondsPerDay {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for SecondOfDay (valid values 0 - 86399): %d", secOfDay))
	}
	h := secOfDay / secondsPerHour
	secOfDay %= secondsPerHour
	m := secOfDay / secondsPerMinute
	s := secOfDay % secondsPerMinute
	return createLocalTime(h, m, s, 0)
}

func localTimeOfNanoOfDay(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.ofNanoOfDay: missing argument")
	}
	nanoOfDay, ok := getInt64Val(params[0])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.ofNanoOfDay: invalid argument")
	}
	if nanoOfDay < 0 || nanoOfDay >= nanosPerDay {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for NanoOfDay (valid values 0 - 86399999999999): %d", nanoOfDay))
	}
	return ofNanoOfDay(nanoOfDay)
}

func localTimeFrom(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.from: null argument")
	}
	obj := params[0].(*object.Object)
	className := object.GoStringFromStringPoolIndex(obj.KlassName)

	if className == classNameLocalTime {
		return obj
	}
	if strings.Contains(className, "LocalDateTime") {
		timeObj := getChronoLocalDateTimeTime(obj)
		return timeObj
	}
	if strings.Contains(className, "OffsetTime") {
		if tFld, ok := obj.FieldTable["time"]; ok {
			if tObj, ok2 := tFld.Fvalue.(*object.Object); ok2 && tObj != nil {
				return tObj
			}
		}
	}
	if strings.Contains(className, "ZonedDateTime") {
		ldt := getChronoZonedDateTimeLocalDateTime(obj)
		return getChronoLocalDateTimeTime(ldt)
	}

	h, m, s, n := getLocalTimeHMSN(obj)
	if _, okH := obj.FieldTable["hour"]; okH {
		return createLocalTime(h, m, s, n)
	}
	if _, okH := obj.FieldTable["HOUR_OF_DAY"]; okH {
		return createLocalTime(h, m, s, n)
	}
	if _, okN := obj.FieldTable["nanoOfDay"]; okN {
		return ofNanoOfDay(getLocalTimeNanoOfDay(obj))
	}

	return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unable to obtain LocalTime from TemporalAccessor")
}

func localTimeParse(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.parse: text is null")
	}
	strObj := params[0].(*object.Object)
	if !object.IsStringObject(strObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.parse: invalid string")
	}
	goStr := strings.TrimSpace(object.GoStringFromStringObject(strObj))
	h, m, s, n, err := parseLocalTimeString(goStr)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeParseException, err.Error())
	}
	return createLocalTime(h, m, s, n)
}

func localTimeParseFormatter(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.parse: null argument")
	}
	return localTimeParse(params[:1])
}

func isSupportedLocalTimeField(fieldName string) bool {
	switch strings.ToUpper(fieldName) {
	case "NANO_OF_SECOND", "NANO_OF_DAY", "MICRO_OF_SECOND", "MICRO_OF_DAY",
		"MILLI_OF_SECOND", "MILLI_OF_DAY", "SECOND_OF_MINUTE", "SECOND_OF_DAY",
		"MINUTE_OF_HOUR", "MINUTE_OF_DAY", "HOUR_OF_AMPM", "CLOCK_HOUR_OF_AMPM",
		"HOUR_OF_DAY", "CLOCK_HOUR_OF_DAY", "AMPM_OF_DAY":
		return true
	default:
		return false
	}
}

func localTimeIsSupportedField(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return false
	}
	fieldObj := params[1].(*object.Object)
	fieldName := getTemporalFieldName(fieldObj)
	return isSupportedLocalTimeField(fieldName)
}

func localTimeIsSupportedUnit(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return false
	}
	unitObj := params[1].(*object.Object)
	unitName := strings.ToUpper(getTemporalUnitName(unitObj))
	switch unitName {
	case "NANOS", "MICROS", "MILLIS", "SECONDS", "MINUTES", "HOURS", "HALF_DAYS":
		return true
	default:
		return false
	}
}

func localTimeRange(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.range: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.range: field is null")
	}
	fieldObj := params[1].(*object.Object)
	fieldName := strings.ToUpper(getTemporalFieldName(fieldObj))

	switch fieldName {
	case "NANO_OF_SECOND":
		return createValueRange(0, 0, 999999999, 999999999)
	case "NANO_OF_DAY":
		return createValueRange(0, 0, 86399999999999, 86399999999999)
	case "MICRO_OF_SECOND":
		return createValueRange(0, 0, 999999, 999999)
	case "MICRO_OF_DAY":
		return createValueRange(0, 0, 86399999999, 86399999999)
	case "MILLI_OF_SECOND":
		return createValueRange(0, 0, 999, 999)
	case "MILLI_OF_DAY":
		return createValueRange(0, 0, 86399999, 86399999)
	case "SECOND_OF_MINUTE":
		return createValueRange(0, 0, 59, 59)
	case "SECOND_OF_DAY":
		return createValueRange(0, 0, 86399, 86399)
	case "MINUTE_OF_HOUR":
		return createValueRange(0, 0, 59, 59)
	case "MINUTE_OF_DAY":
		return createValueRange(0, 0, 1439, 1439)
	case "HOUR_OF_AMPM":
		return createValueRange(0, 0, 11, 11)
	case "CLOCK_HOUR_OF_AMPM":
		return createValueRange(1, 1, 12, 12)
	case "HOUR_OF_DAY":
		return createValueRange(0, 0, 23, 23)
	case "CLOCK_HOUR_OF_DAY":
		return createValueRange(1, 1, 24, 24)
	case "AMPM_OF_DAY":
		return createValueRange(0, 0, 1, 1)
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
	}
}

func localTimeGet(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.get: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.get: field is null")
	}
	fieldObj := params[1].(*object.Object)
	fieldName := strings.ToUpper(getTemporalFieldName(fieldObj))

	if fieldName == "NANO_OF_DAY" || fieldName == "MICRO_OF_DAY" {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid field %s for get() method, use getLong() instead", fieldName))
	}

	val := localTimeGetLong(params)
	if err, isErr := val.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if v, ok := getInt64Val(val); ok {
		return v
	}
	return val
}

func localTimeGetLong(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.getLong: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.getLong: field is null")
	}
	self := params[0].(*object.Object)
	fieldObj := params[1].(*object.Object)
	fieldName := strings.ToUpper(getTemporalFieldName(fieldObj))

	h, m, s, n := getLocalTimeHMSN(self)
	nod := h*nanosPerHour + m*nanosPerMinute + s*nanosPerSecond + n

	switch fieldName {
	case "NANO_OF_SECOND":
		return n
	case "NANO_OF_DAY":
		return nod
	case "MICRO_OF_SECOND":
		return n / 1000
	case "MICRO_OF_DAY":
		return nod / 1000
	case "MILLI_OF_SECOND":
		return n / 1000000
	case "MILLI_OF_DAY":
		return nod / 1000000
	case "SECOND_OF_MINUTE":
		return s
	case "SECOND_OF_DAY":
		return h*secondsPerHour + m*secondsPerMinute + s
	case "MINUTE_OF_HOUR":
		return m
	case "MINUTE_OF_DAY":
		return h*minutesPerHour + m
	case "HOUR_OF_AMPM":
		return h % 12
	case "CLOCK_HOUR_OF_AMPM":
		r := h % 12
		if r == 0 {
			return int64(12)
		}
		return r
	case "HOUR_OF_DAY":
		return h
	case "CLOCK_HOUR_OF_DAY":
		if h == 0 {
			return int64(24)
		}
		return h
	case "AMPM_OF_DAY":
		return h / 12
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
	}
}

func localTimeGetHour(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.getHour: self is null")
	}
	h, _, _, _ := getLocalTimeHMSN(params[0].(*object.Object))
	return h
}

func localTimeGetMinute(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.getMinute: self is null")
	}
	_, m, _, _ := getLocalTimeHMSN(params[0].(*object.Object))
	return m
}

func localTimeGetSecond(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.getSecond: self is null")
	}
	_, _, s, _ := getLocalTimeHMSN(params[0].(*object.Object))
	return s
}

func localTimeGetNano(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.getNano: self is null")
	}
	_, _, _, n := getLocalTimeHMSN(params[0].(*object.Object))
	return n
}

func localTimeWithHour(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.withHour: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.withHour: missing hour")
	}
	newH, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.withHour: invalid hour")
	}
	h, m, s, n := getLocalTimeHMSN(params[0].(*object.Object))
	if h == newH {
		return params[0]
	}
	if newH < 0 || newH > 23 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for HourOfDay (valid values 0 - 23): %d", newH))
	}
	return createLocalTime(newH, m, s, n)
}

func localTimeWithMinute(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.withMinute: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.withMinute: missing minute")
	}
	newM, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.withMinute: invalid minute")
	}
	h, m, s, n := getLocalTimeHMSN(params[0].(*object.Object))
	if m == newM {
		return params[0]
	}
	if newM < 0 || newM > 59 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for MinuteOfHour (valid values 0 - 59): %d", newM))
	}
	return createLocalTime(h, newM, s, n)
}

func localTimeWithSecond(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.withSecond: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.withSecond: missing second")
	}
	newS, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.withSecond: invalid second")
	}
	h, m, s, n := getLocalTimeHMSN(params[0].(*object.Object))
	if s == newS {
		return params[0]
	}
	if newS < 0 || newS > 59 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for SecondOfMinute (valid values 0 - 59): %d", newS))
	}
	return createLocalTime(h, m, newS, n)
}

func localTimeWithNano(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.withNano: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.withNano: missing nano")
	}
	newN, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.withNano: invalid nano")
	}
	h, m, s, n := getLocalTimeHMSN(params[0].(*object.Object))
	if n == newN {
		return params[0]
	}
	if newN < 0 || newN > 999999999 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for NanoOfSecond (valid values 0 - 999999999): %d", newN))
	}
	return createLocalTime(h, m, s, newN)
}

func localTimeWithField(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.with: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.with: field is null")
	}
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.with: missing newValue")
	}
	self := params[0].(*object.Object)
	fieldObj := params[1].(*object.Object)
	newValue, okVal := getInt64Val(params[2])
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.with: invalid newValue")
	}
	fieldName := strings.ToUpper(getTemporalFieldName(fieldObj))
	h, m, s, n := getLocalTimeHMSN(self)

	switch fieldName {
	case "NANO_OF_SECOND":
		return localTimeWithNano([]any{self, newValue})
	case "NANO_OF_DAY":
		if newValue < 0 || newValue >= nanosPerDay {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for NanoOfDay (valid values 0 - 86399999999999): %d", newValue))
		}
		return ofNanoOfDay(newValue)
	case "MICRO_OF_SECOND":
		if newValue < 0 || newValue > 999999 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for MicroOfSecond (valid values 0 - 999999): %d", newValue))
		}
		return createLocalTime(h, m, s, newValue*1000+(n%1000))
	case "MICRO_OF_DAY":
		if newValue < 0 || newValue >= microsPerDay {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for MicroOfDay (valid values 0 - 86399999999): %d", newValue))
		}
		return ofNanoOfDay(newValue * 1000)
	case "MILLI_OF_SECOND":
		if newValue < 0 || newValue > 999 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for MilliOfSecond (valid values 0 - 999): %d", newValue))
		}
		return createLocalTime(h, m, s, newValue*1000000+(n%1000000))
	case "MILLI_OF_DAY":
		if newValue < 0 || newValue >= millisPerDay {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for MilliOfDay (valid values 0 - 86399999): %d", newValue))
		}
		return ofNanoOfDay(newValue * 1000000)
	case "SECOND_OF_MINUTE":
		return localTimeWithSecond([]any{self, newValue})
	case "SECOND_OF_DAY":
		if newValue < 0 || newValue >= secondsPerDay {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for SecondOfDay (valid values 0 - 86399): %d", newValue))
		}
		curSec := h*secondsPerHour + m*secondsPerMinute + s
		return localTimePlusSeconds([]any{self, newValue - curSec})
	case "MINUTE_OF_HOUR":
		return localTimeWithMinute([]any{self, newValue})
	case "MINUTE_OF_DAY":
		if newValue < 0 || newValue >= minutesPerDay {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for MinuteOfDay (valid values 0 - 1439): %d", newValue))
		}
		curMin := h*minutesPerHour + m
		return localTimePlusMinutes([]any{self, newValue - curMin})
	case "HOUR_OF_AMPM":
		if newValue < 0 || newValue > 11 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for HourOfAmPm (valid values 0 - 11): %d", newValue))
		}
		return localTimePlusHours([]any{self, newValue - (h % 12)})
	case "CLOCK_HOUR_OF_AMPM":
		if newValue < 1 || newValue > 12 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for ClockHourOfAmPm (valid values 1 - 12): %d", newValue))
		}
		targetH := newValue
		if targetH == 12 {
			targetH = 0
		}
		return localTimePlusHours([]any{self, targetH - (h % 12)})
	case "HOUR_OF_DAY":
		return localTimeWithHour([]any{self, newValue})
	case "CLOCK_HOUR_OF_DAY":
		if newValue < 1 || newValue > 24 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for ClockHourOfDay (valid values 1 - 24): %d", newValue))
		}
		targetH := newValue
		if targetH == 24 {
			targetH = 0
		}
		return localTimeWithHour([]any{self, targetH})
	case "AMPM_OF_DAY":
		if newValue < 0 || newValue > 1 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid value for AmPmOfDay (valid values 0 - 1): %d", newValue))
		}
		return localTimePlusHours([]any{self, (newValue - (h / 12)) * 12})
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
	}
}

func localTimeWithAdjuster(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.with: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.with: adjuster is null")
	}
	self := params[0].(*object.Object)
	adjuster := params[1].(*object.Object)
	className := object.GoStringFromStringPoolIndex(adjuster.KlassName)

	if className == classNameLocalTime {
		return adjuster
	}
	res := temporalAdjusterAdjustInto([]any{adjuster, self})
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if resObj, ok := res.(*object.Object); ok {
		return localTimeFrom([]any{resObj})
	}
	return res
}

func localTimeTruncatedTo(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.truncatedTo: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.truncatedTo: unit is null")
	}
	self := params[0].(*object.Object)
	unitObj := params[1].(*object.Object)
	unitName := strings.ToUpper(getTemporalUnitName(unitObj))

	h, m, s, n := getLocalTimeHMSN(self)

	switch unitName {
	case "NANOS":
		return self
	case "MICROS":
		return createLocalTime(h, m, s, n-(n%1000))
	case "MILLIS":
		return createLocalTime(h, m, s, n-(n%1000000))
	case "SECONDS":
		return createLocalTime(h, m, s, 0)
	case "MINUTES":
		return createLocalTime(h, m, 0, 0)
	case "HOURS":
		return createLocalTime(h, 0, 0, 0)
	case "HALF_DAYS":
		return createLocalTime((h/12)*12, 0, 0, 0)
	case "DAYS":
		return createLocalTime(0, 0, 0, 0)
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unit is too large to be used for truncation: %s", unitName))
	}
}

func localTimePlusHours(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.plusHours: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.plusHours: missing hoursToAdd")
	}
	hoursToAdd, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.plusHours: invalid argument")
	}
	h, m, s, n := getLocalTimeHMSN(params[0].(*object.Object))
	return localTimePlusHoursInternal(h, m, s, n, hoursToAdd)
}

func localTimePlusMinutes(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.plusMinutes: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.plusMinutes: missing minutesToAdd")
	}
	minutesToAdd, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.plusMinutes: invalid argument")
	}
	h, m, s, n := getLocalTimeHMSN(params[0].(*object.Object))
	return localTimePlusMinutesInternal(h, m, s, n, minutesToAdd)
}

func localTimePlusSeconds(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.plusSeconds: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.plusSeconds: missing secondsToAdd")
	}
	secondsToAdd, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.plusSeconds: invalid argument")
	}
	h, m, s, n := getLocalTimeHMSN(params[0].(*object.Object))
	return localTimePlusSecondsInternal(h, m, s, n, secondsToAdd)
}

func localTimePlusNanos(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.plusNanos: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.plusNanos: missing nanosToAdd")
	}
	nanosToAdd, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.plusNanos: invalid argument")
	}
	h, m, s, n := getLocalTimeHMSN(params[0].(*object.Object))
	return localTimePlusNanosInternal(h, m, s, n, nanosToAdd)
}

func localTimePlusUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.plus: self is null")
	}
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.plus: missing arguments")
	}
	amount, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.plus: invalid amount")
	}
	if params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.plus: unit is null")
	}
	unitObj := params[2].(*object.Object)
	unitName := strings.ToUpper(getTemporalUnitName(unitObj))
	h, m, s, n := getLocalTimeHMSN(params[0].(*object.Object))

	switch unitName {
	case "NANOS":
		return localTimePlusNanosInternal(h, m, s, n, amount)
	case "MICROS":
		return localTimePlusNanosInternal(h, m, s, n, (amount%microsPerDay)*1000)
	case "MILLIS":
		return localTimePlusNanosInternal(h, m, s, n, (amount%millisPerDay)*1000000)
	case "SECONDS":
		return localTimePlusSecondsInternal(h, m, s, n, amount)
	case "MINUTES":
		return localTimePlusMinutesInternal(h, m, s, n, amount)
	case "HOURS":
		return localTimePlusHoursInternal(h, m, s, n, amount)
	case "HALF_DAYS":
		return localTimePlusHoursInternal(h, m, s, n, (amount%2)*12)
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported unit: %s", unitName))
	}
}

func localTimePlusAmount(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.plus: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.plus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)
	className := object.GoStringFromStringPoolIndex(amount.KlassName)

	if strings.Contains(className, "Duration") {
		var sec int64 = 0
		var nano int64 = 0
		if sFld, ok := amount.FieldTable["seconds"]; ok {
			if s, ok2 := getInt64Val(sFld.Fvalue); ok2 {
				sec = s
			}
		}
		if nFld, ok := amount.FieldTable["nanos"]; ok {
			if n, ok2 := getInt64Val(nFld.Fvalue); ok2 {
				nano = n
			}
		}
		h, m, s, n := getLocalTimeHMSN(self)
		res := localTimePlusSecondsInternal(h, m, s, n, sec)
		rh, rm, rs, rn := getLocalTimeHMSN(res)
		return localTimePlusNanosInternal(rh, rm, rs, rn, nano)
	}

	res := temporalAmountAddTo([]any{amount, self})
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if resObj, ok := res.(*object.Object); ok {
		return localTimeFrom([]any{resObj})
	}
	return res
}

func localTimeMinusHours(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.minusHours: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.minusHours: missing hoursToSubtract")
	}
	hours, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.minusHours: invalid argument")
	}
	return localTimePlusHours([]any{params[0], -hours})
}

func localTimeMinusMinutes(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.minusMinutes: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.minusMinutes: missing minutesToSubtract")
	}
	minutes, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.minusMinutes: invalid argument")
	}
	return localTimePlusMinutes([]any{params[0], -minutes})
}

func localTimeMinusSeconds(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.minusSeconds: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.minusSeconds: missing secondsToSubtract")
	}
	seconds, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.minusSeconds: invalid argument")
	}
	return localTimePlusSeconds([]any{params[0], -seconds})
}

func localTimeMinusNanos(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.minusNanos: self is null")
	}
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.minusNanos: missing nanosToSubtract")
	}
	nanos, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.minusNanos: invalid argument")
	}
	return localTimePlusNanos([]any{params[0], -nanos})
}

func localTimeMinusUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.minus: self is null")
	}
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.minus: missing arguments")
	}
	amount, ok := getInt64Val(params[1])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "LocalTime.minus: invalid amount")
	}
	return localTimePlusUnit([]any{params[0], -amount, params[2]})
}

func localTimeMinusAmount(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.minus: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.minus: amount is null")
	}
	self := params[0].(*object.Object)
	amount := params[1].(*object.Object)
	className := object.GoStringFromStringPoolIndex(amount.KlassName)

	if strings.Contains(className, "Duration") {
		var sec int64 = 0
		var nano int64 = 0
		if sFld, ok := amount.FieldTable["seconds"]; ok {
			if s, ok2 := getInt64Val(sFld.Fvalue); ok2 {
				sec = s
			}
		}
		if nFld, ok := amount.FieldTable["nanos"]; ok {
			if n, ok2 := getInt64Val(nFld.Fvalue); ok2 {
				nano = n
			}
		}
		h, m, s, n := getLocalTimeHMSN(self)
		res := localTimePlusSecondsInternal(h, m, s, n, -sec)
		rh, rm, rs, rn := getLocalTimeHMSN(res)
		return localTimePlusNanosInternal(rh, rm, rs, rn, -nano)
	}

	res := temporalAmountSubtractFrom([]any{amount, self})
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if resObj, ok := res.(*object.Object); ok {
		return localTimeFrom([]any{resObj})
	}
	return res
}

func localTimeQuery(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.query: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.query: query is null")
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
		return ghelpers.CreateTemporalUnit("Nanos")
	case "localTime":
		return self
	case "chronology", "zoneId", "zone", "offset", "localDate":
		return nil
	default:
		return executeTemporalQuery(query, self)
	}
}

func localTimeAdjustInto(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.adjustInto: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.adjustInto: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)
	h, m, s, n := getLocalTimeHMSN(self)
	nanoOfDay := getLocalTimeNanoOfDay(self)

	if temporal.FieldTable != nil {
		targetClassName := object.GoStringFromStringPoolIndex(temporal.KlassName)
		if targetClassName == classNameLocalTime || temporal.FieldTable["hour"].Fvalue != nil {
			temporal.FieldTable["hour"] = object.Field{Ftype: types.Int, Fvalue: h}
			temporal.FieldTable["minute"] = object.Field{Ftype: types.Int, Fvalue: m}
			temporal.FieldTable["second"] = object.Field{Ftype: types.Int, Fvalue: s}
			temporal.FieldTable["nano"] = object.Field{Ftype: types.Int, Fvalue: n}
			temporal.FieldTable["HOUR_OF_DAY"] = object.Field{Ftype: types.Long, Fvalue: h}
			temporal.FieldTable["MINUTE_OF_HOUR"] = object.Field{Ftype: types.Long, Fvalue: m}
			temporal.FieldTable["SECOND_OF_MINUTE"] = object.Field{Ftype: types.Long, Fvalue: s}
			temporal.FieldTable["NANO_OF_SECOND"] = object.Field{Ftype: types.Long, Fvalue: n}
			temporal.FieldTable["nanoOfDay"] = object.Field{Ftype: types.Long, Fvalue: nanoOfDay}
			return temporal
		}
		if _, ok := temporal.FieldTable["time"]; ok {
			temporal.FieldTable["time"] = object.Field{Ftype: types.Ref, Fvalue: self}
			return temporal
		}
	}
	cfNanoOfDay := getChronoField("NANO_OF_DAY")
	return temporalWithField([]any{temporal, cfNanoOfDay, nanoOfDay})
}

func localTimeUntilUnit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.until: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.until: endExclusive is null")
	}
	if len(params) < 3 || params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.until: unit is null")
	}
	self := params[0].(*object.Object)
	endObj := params[1].(*object.Object)
	unitObj := params[2].(*object.Object)

	endLocalTime := localTimeFrom([]any{endObj})
	if err, isErr := endLocalTime.(*ghelpers.GErrBlk); isErr {
		return err
	}
	nod1 := getLocalTimeNanoOfDay(self)
	nod2 := getLocalTimeNanoOfDay(endLocalTime.(*object.Object))
	nanosUntil := nod2 - nod1

	unitName := strings.ToUpper(getTemporalUnitName(unitObj))
	switch unitName {
	case "NANOS":
		return nanosUntil
	case "MICROS":
		return nanosUntil / 1000
	case "MILLIS":
		return nanosUntil / 1000000
	case "SECONDS":
		return nanosUntil / nanosPerSecond
	case "MINUTES":
		return nanosUntil / nanosPerMinute
	case "HOURS":
		return nanosUntil / nanosPerHour
	case "HALF_DAYS":
		return nanosUntil / nanosPerHalfDay
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported unit: %s", unitName))
	}
}

func localTimeFormat(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.format: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.format: formatter is null")
	}
	return dateTimeFormatterFormat([]any{params[1], params[0]})
}

func localTimeAtDate(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.atDate: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.atDate: date is null")
	}
	self := params[0].(*object.Object)
	dateObj := params[1].(*object.Object)
	return createLocalDateTime(dateObj, self)
}

func localTimeAtOffset(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.atOffset: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.atOffset: offset is null")
	}
	self := params[0].(*object.Object)
	offsetObj := params[1].(*object.Object)
	classNameOffsetTime := "java/time/OffsetTime"
	ot := object.MakeEmptyObjectWithClassName(&classNameOffsetTime)
	ot.FieldTable["time"] = object.Field{Ftype: types.Ref, Fvalue: self}
	ot.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: offsetObj}
	return ot
}

func localTimeToSecondOfDay(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.toSecondOfDay: self is null")
	}
	h, m, s, _ := getLocalTimeHMSN(params[0].(*object.Object))
	return h*secondsPerHour + m*secondsPerMinute + s
}

func localTimeToNanoOfDay(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.toNanoOfDay: self is null")
	}
	return getLocalTimeNanoOfDay(params[0].(*object.Object))
}

func localTimeToEpochSecond(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.toEpochSecond: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.toEpochSecond: date is null")
	}
	if len(params) < 3 || params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.toEpochSecond: offset is null")
	}
	self := params[0].(*object.Object)
	dateObj := params[1].(*object.Object)
	offsetObj := params[2].(*object.Object)

	epochDay := ghelpers.GetChronoLocalDateEpochDay(dateObj)
	secOfDay := getLocalTimeNanoOfDay(self) / nanosPerSecond
	offsetSec := getZoneOffsetTotalSeconds(offsetObj)

	epochSec := epochDay*secondsPerDay + secOfDay - offsetSec
	return epochSec
}

func localTimeCompareTo(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.compareTo: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.compareTo: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)
	return compareLocalTime(self, other)
}

func localTimeIsAfter(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.isAfter: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.isAfter: other is null")
	}
	cmp := localTimeCompareTo(params)
	if err, isErr := cmp.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if val, ok := getInt64Val(cmp); ok {
		return val > 0
	}
	return false
}

func localTimeIsBefore(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.isBefore: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.isBefore: other is null")
	}
	cmp := localTimeCompareTo(params)
	if err, isErr := cmp.(*ghelpers.GErrBlk); isErr {
		return err
	}
	if val, ok := getInt64Val(cmp); ok {
		return val < 0
	}
	return false
}

func localTimeEquals(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil {
		return false
	}
	self, okS := params[0].(*object.Object)
	other, okO := params[1].(*object.Object)
	if !okS || !okO || object.IsNull(self) || object.IsNull(other) {
		return false
	}
	if self == other {
		return true
	}
	h1, m1, s1, n1 := getLocalTimeHMSN(self)
	h2, m2, s2, n2 := getLocalTimeHMSN(other)
	return h1 == h2 && m1 == m2 && s1 == s2 && n1 == n2
}

func localTimeHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.hashCode: self is null")
	}
	nod := getLocalTimeNanoOfDay(params[0].(*object.Object))
	nodU := uint64(nod)
	return int64(int32(nodU ^ (nodU >> 32)))
}

func localTimeToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "LocalTime.toString: self is null")
	}
	h, m, s, n := getLocalTimeHMSN(params[0].(*object.Object))
	return object.StringObjectFromGoString(formatLocalTimeString(h, m, s, n))
}
