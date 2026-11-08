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
	"sync"
)

var (
	isoChronologyMutex       sync.Mutex
	isoChronologyInitialized bool           = false
	isoChronologySingleton   *object.Object = nil
)

func getIsoChronologySingleton() *object.Object {
	isoChronologyMutex.Lock()
	defer isoChronologyMutex.Unlock()
	if !isoChronologyInitialized {
		isoChronologySingleton = GetIsoChronology()
		_ = statics.AddStatic(classNameIsoChronology+".INSTANCE", statics.Static{
			Type:  "Ljava/time/chrono/IsoChronology;",
			Value: isoChronologySingleton,
		})
		isoChronologyInitialized = true
	}
	return isoChronologySingleton
}

// Load_Time_Chrono_IsoChronology loads method signatures for java/time/chrono/IsoChronology.
func Load_Time_Chrono_IsoChronology() {
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  isoChronologyClinit,
		}

	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.compareTo(Ljava/time/chrono/Chronology;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyCompareTo,
		}

	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyCompareTo,
		}

	// date(III)
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.date(III)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  chronologyDateYMD,
		}
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.date(III)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  chronologyDateYMD,
		}

	// date(Era, III)
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.date(Ljava/time/chrono/Era;III)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  chronologyDateEra,
		}
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.date(Ljava/time/chrono/Era;III)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  chronologyDateEra,
		}

	// date(TemporalAccessor)
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.date(Ljava/time/temporal/TemporalAccessor;)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyDateTemporal,
		}
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.date(Ljava/time/temporal/TemporalAccessor;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyDateTemporal,
		}

	// dateEpochDay(J)
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.dateEpochDay(J)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronologyDateEpochDay,
		}
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.dateEpochDay(J)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronologyDateEpochDay,
		}

	// dateNow()
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.dateNow()Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronologyDateNow,
		}
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.dateNow()Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronologyDateNow,
		}

	// dateNow(Clock)
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.dateNow(Ljava/time/Clock;)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyDateNowClock,
		}
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.dateNow(Ljava/time/Clock;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyDateNowClock,
		}

	// dateNow(ZoneId)
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.dateNow(Ljava/time/ZoneId;)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyDateNowZone,
		}
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.dateNow(Ljava/time/ZoneId;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyDateNowZone,
		}

	// dateYearDay(II)
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.dateYearDay(II)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronologyDateYearDay,
		}
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.dateYearDay(II)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronologyDateYearDay,
		}

	// dateYearDay(Era, II)
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.dateYearDay(Ljava/time/chrono/Era;II)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  chronologyDateYearDayEra,
		}
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.dateYearDay(Ljava/time/chrono/Era;II)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  chronologyDateYearDayEra,
		}

	// epochSecond(IIIIIILjava/time/ZoneOffset;)J
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.epochSecond(IIIIIILjava/time/ZoneOffset;)J"] =
		ghelpers.GMeth{
			ParamSlots: 7,
			GFunction:  chronologyEpochSecond,
		}

	// epochSecond(Era, IIIIIILjava/time/ZoneOffset;)J
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.epochSecond(Ljava/time/chrono/Era;IIIIIILjava/time/ZoneOffset;)J"] =
		ghelpers.GMeth{
			ParamSlots: 8,
			GFunction:  chronologyEpochSecondEra,
		}

	// equals(Object)
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyEquals,
		}

	// eraOf(I)
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.eraOf(I)Ljava/time/chrono/IsoEra;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  isoChronologyEraOf,
		}
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.eraOf(I)Ljava/time/chrono/Era;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  isoChronologyEraOf,
		}

	// eras()
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.eras()Ljava/util/List;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronologyEras,
		}

	// getCalendarType()
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.getCalendarType()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronologyGetCalendarType,
		}

	// getId()
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.getId()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronologyGetId,
		}

	// hashCode()
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronologyHashCode,
		}

	// isLeapYear(J)
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.isLeapYear(J)Z"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  isoChronologyIsLeapYear,
		}

	// localDateTime(TemporalAccessor)
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.localDateTime(Ljava/time/temporal/TemporalAccessor;)Ljava/time/LocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyLocalDateTime,
		}
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.localDateTime(Ljava/time/temporal/TemporalAccessor;)Ljava/time/chrono/ChronoLocalDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyLocalDateTime,
		}

	// period(III)
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.period(III)Ljava/time/Period;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  chronologyPeriod,
		}
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.period(III)Ljava/time/chrono/ChronoPeriod;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  chronologyPeriod,
		}

	// prolepticYear(Era, I)
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.prolepticYear(Ljava/time/chrono/Era;I)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  isoChronologyProlepticYear,
		}

	// range(ChronoField)
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.range(Ljava/time/temporal/ChronoField;)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  isoChronologyRange,
		}

	// resolveDate(Map, ResolverStyle)
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.resolveDate(Ljava/util/Map;Ljava/time/format/ResolverStyle;)Ljava/time/LocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronologyResolveDate,
		}
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.resolveDate(Ljava/util/Map;Ljava/time/format/ResolverStyle;)Ljava/time/chrono/ChronoLocalDate;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronologyResolveDate,
		}

	// toString()
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  chronologyToString,
		}

	// zonedDateTime(Instant, ZoneId)
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.zonedDateTime(Ljava/time/Instant;Ljava/time/ZoneId;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronologyZonedDateTimeInstant,
		}
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.zonedDateTime(Ljava/time/Instant;Ljava/time/ZoneId;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  chronologyZonedDateTimeInstant,
		}

	// zonedDateTime(TemporalAccessor)
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.zonedDateTime(Ljava/time/temporal/TemporalAccessor;)Ljava/time/ZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyZonedDateTimeTemporal,
		}
	ghelpers.MethodSignatures["java/time/chrono/IsoChronology.zonedDateTime(Ljava/time/temporal/TemporalAccessor;)Ljava/time/chrono/ChronoZonedDateTime;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  chronologyZonedDateTimeTemporal,
		}
}

func isoChronologyClinit([]any) any {
	getIsoChronologySingleton()
	return nil
}

func isoChronologyIsLeapYear(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "IsoChronology.isLeapYear: self is null")
	}
	prolepticYear, okYear := getInt64Val(params[1])
	if !okYear {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "IsoChronology.isLeapYear: invalid year parameter")
	}

	if (prolepticYear&3) == 0 && ((prolepticYear%100) != 0 || (prolepticYear%400) == 0) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func isoChronologyEraOf(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "IsoChronology.eraOf: self is null")
	}
	eraValue, okValue := getInt64Val(params[1])
	if !okValue {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "IsoChronology.eraOf: invalid eraValue parameter")
	}

	switch eraValue {
	case 0:
		return createEra(0, "BCE")
	case 1:
		return createEra(1, "CE")
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Invalid era: %d", eraValue))
	}
}

func isoChronologyProlepticYear(params []any) any {
	if len(params) < 3 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "IsoChronology.prolepticYear: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "IsoChronology.prolepticYear: era is null or not an Era")
	}
	eraObject := params[1].(*object.Object)
	yearOfEra, okYear := getInt64Val(params[2])
	if !okYear {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "IsoChronology.prolepticYear: invalid yearOfEra parameter")
	}

	var eraValue int64 = 1
	if eraObject.FieldTable != nil {
		if eraValueField, exists := eraObject.FieldTable["eraValue"]; exists {
			if parsedEraValue, okParsed := getInt64Val(eraValueField.Fvalue); okParsed {
				eraValue = parsedEraValue
			}
		}
	}

	if eraValue == 1 {
		return yearOfEra
	} else if eraValue == 0 {
		return int64(1) - yearOfEra
	}
	return ghelpers.GetGErrBlk(excNames.ClassCastException, "IsoChronology.prolepticYear: era must be IsoEra")
}

func isoChronologyRange(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "IsoChronology.range: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "IsoChronology.range: field is null")
	}
	fieldObject := params[1].(*object.Object)

	fieldName := ""
	if fieldObject.FieldTable != nil {
		if nameField, exists := fieldObject.FieldTable["name"]; exists {
			if nameStringObject, okString := nameField.Fvalue.(*object.Object); okString && object.IsStringObject(nameStringObject) {
				fieldName = object.GoStringFromStringObject(nameStringObject)
			}
		}
	}

	switch fieldName {
	case "DAY_OF_MONTH":
		return createValueRange(1, 1, 28, 31)
	case "DAY_OF_YEAR":
		return createValueRange(1, 1, 365, 366)
	case "ALIGNED_WEEK_OF_MONTH":
		return createValueRange(1, 1, 4, 5)
	case "YEAR_OF_ERA":
		return createValueRange(1, 1, 999_999_999, 1_000_000_000)
	case "PROLEPTIC_MONTH":
		return createValueRange(-11_999_999_988, -11_999_999_988, 11_999_999_999, 11_999_999_999)
	case "YEAR":
		return createValueRange(-999_999_999, -999_999_999, 999_999_999, 999_999_999)
	case "ERA":
		return createValueRange(0, 0, 1, 1)
	default:
		if fieldObject.FieldTable != nil {
			if rangeField, exists := fieldObject.FieldTable["range"]; exists {
				if rangeObject, okRange := rangeField.Fvalue.(*object.Object); okRange && !object.IsNull(rangeObject) {
					return rangeObject
				}
			}
		}
		return createValueRange(1, 1, math.MaxInt64, math.MaxInt64)
	}
}
