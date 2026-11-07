package ghelpers

import (
	"jacobin/src/object"
	"jacobin/src/types"
	"time"
)

// CreateTemporalUnit allocates a new ChronoUnit/TemporalUnit object with the given name.
func CreateTemporalUnit(name string) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&types.ClassNameChronoUnit)
	nameObj := object.StringObjectFromGoString(name)
	obj.FieldTable["name"] = object.Field{Ftype: types.Ref, Fvalue: nameObj}
	return obj
}

func GetTemporalAmountDays(amount *object.Object) int64 {
	if amount == nil || amount.FieldTable == nil {
		return 0
	}
	var days int64 = 0
	if daysFld, ok := amount.FieldTable["days"]; ok {
		if d, ok := GetInt64Val(daysFld.Fvalue); ok {
			days += d
		}
	}
	if secFld, ok := amount.FieldTable["seconds"]; ok {
		if s, ok := GetInt64Val(secFld.Fvalue); ok {
			days += s / 86400
		}
	}
	return days
}

func IsLeapYearHelper(year int64) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

func DaysInMonthHelper(year, month int64) int64 {
	switch month {
	case 2:
		if IsLeapYearHelper(year) {
			return 29
		}
		return 28
	case 4, 6, 9, 11:
		return 30
	default:
		return 31
	}
}

func GetChronoLocalDateEpochDay(obj *object.Object) int64 {
	if obj == nil || obj.FieldTable == nil {
		return 0
	}
	if epochFld, ok := obj.FieldTable["epochDay"]; ok {
		if epochVal, ok := GetInt64Val(epochFld.Fvalue); ok {
			return epochVal
		}
	}
	if yFld, ok := obj.FieldTable["year"]; ok {
		if yVal, ok := GetInt64Val(yFld.Fvalue); ok {
			var mVal int64 = 1
			var dVal int64 = 1
			if mFld, ok := obj.FieldTable["month"]; ok {
				if m, ok := GetInt64Val(mFld.Fvalue); ok {
					mVal = m
				}
			}
			if dFld, ok := obj.FieldTable["day"]; ok {
				if d, ok := GetInt64Val(dFld.Fvalue); ok {
					dVal = d
				}
			}
			t := time.Date(int(yVal), time.Month(mVal), int(dVal), 0, 0, 0, 0, time.UTC)
			return t.Unix() / 86400
		}
	}
	return 0
}

func GetChronoLocalDateYMD(obj *object.Object) (int64, int64, int64) {
	if obj == nil || obj.FieldTable == nil {
		return 1970, 1, 1
	}
	var year int64 = 1970
	var month int64 = 1
	var day int64 = 1

	if yFld, ok := obj.FieldTable["year"]; ok {
		if yVal, ok := GetInt64Val(yFld.Fvalue); ok {
			year = yVal
		}
	}
	if mFld, ok := obj.FieldTable["month"]; ok {
		if mVal, ok := GetInt64Val(mFld.Fvalue); ok {
			month = mVal
		}
	}
	if dFld, ok := obj.FieldTable["day"]; ok {
		if dVal, ok := GetInt64Val(dFld.Fvalue); ok {
			day = dVal
		}
	}
	if _, ok := obj.FieldTable["year"]; !ok {
		if epochFld, ok := obj.FieldTable["epochDay"]; ok {
			if epochDay, ok := GetInt64Val(epochFld.Fvalue); ok {
				t := time.Unix(epochDay*86400, 0).UTC()
				year = int64(t.Year())
				month = int64(t.Month())
				day = int64(t.Day())
			}
		}
	}
	return year, month, day
}

func CreateChronoLocalDateTimeFromDateAndTime(dateObj *object.Object, timeObj *object.Object) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&types.ClassNameChronoLocalDateTime)
	obj.FieldTable["date"] = object.Field{Ftype: types.Ref, Fvalue: dateObj}
	obj.FieldTable["time"] = object.Field{Ftype: types.Ref, Fvalue: timeObj}
	return obj
}

func CreateChronoPeriodYMD(chrono *object.Object, years, months, days int64) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&types.ClassNameChronoPeriod)
	obj.FieldTable["years"] = object.Field{Ftype: types.Int, Fvalue: years}
	obj.FieldTable["months"] = object.Field{Ftype: types.Int, Fvalue: months}
	obj.FieldTable["days"] = object.Field{Ftype: types.Int, Fvalue: days}
	obj.FieldTable["chronology"] = object.Field{Ftype: types.Ref, Fvalue: chrono}
	return obj
}

// CreateChronology creates a Chronology object with given id, calendarType, and className.
func CreateChronology(className, id, calendarType string) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&className)
	idObj := object.StringObjectFromGoString(id)
	calObj := object.StringObjectFromGoString(calendarType)
	obj.FieldTable["id"] = object.Field{Ftype: types.Ref, Fvalue: idObj}
	obj.FieldTable["calendarType"] = object.Field{Ftype: types.Ref, Fvalue: calObj}
	return obj
}
