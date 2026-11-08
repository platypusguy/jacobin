/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) Consult jacobin.org.
 */

package ghelpers

import (
	"fmt"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestCreateTemporalUnit(t *testing.T) {
	unit := CreateTemporalUnit("DAYS")
	if got := object.GoStringFromStringPoolIndex(unit.KlassName); got != types.ClassNameChronoUnit {
		t.Errorf("class = %q, want %q", got, types.ClassNameChronoUnit)
	}
	name, ok := unit.FieldTable["name"]
	if !ok {
		t.Fatal("missing name field")
	}
	if name.Ftype != types.Ref {
		t.Errorf("name field type = %q, want %q", name.Ftype, types.Ref)
	}
	nameObj, ok := name.Fvalue.(*object.Object)
	if !ok {
		t.Fatalf("name field value has type %T, want *object.Object", name.Fvalue)
	}
	if got := object.GoStringFromStringObject(nameObj); got != "DAYS" {
		t.Errorf("name = %q, want %q", got, "DAYS")
	}
}

func TestGetTemporalAmountDays(t *testing.T) {
	tests := []struct {
		name   string
		amount *object.Object
		want   int64
	}{
		{name: "nil amount", amount: nil, want: 0},
		{name: "nil field table", amount: &object.Object{}, want: 0},
		{name: "empty amount", amount: object.MakeEmptyObject(), want: 0},
		{name: "days only", amount: dateTimeTestObject(map[string]object.Field{"days": dateTimeTestField(int64(4))}), want: 4},
		{name: "seconds only", amount: dateTimeTestObject(map[string]object.Field{"seconds": dateTimeTestField(int64(172799))}), want: 1},
		{name: "combined days and seconds", amount: dateTimeTestObject(map[string]object.Field{
			"days": dateTimeTestField(int64(2)), "seconds": dateTimeTestField(int64(172800)),
		}), want: 4},
		{name: "integer representations", amount: dateTimeTestObject(map[string]object.Field{
			"days": dateTimeTestField(int32(2)), "seconds": dateTimeTestField(int(86400)),
		}), want: 3},
		{name: "negative seconds truncate toward zero", amount: dateTimeTestObject(map[string]object.Field{
			"seconds": dateTimeTestField(int64(-86401)),
		}), want: -1},
		{name: "invalid fields ignored", amount: dateTimeTestObject(map[string]object.Field{
			"days": dateTimeTestField(float64(3)), "seconds": dateTimeTestField("86400"),
		}), want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTemporalAmountDays(tt.amount); got != tt.want {
				t.Errorf("GetTemporalAmountDays() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestIsLeapYearHelper(t *testing.T) {
	tests := []struct {
		year int64
		want bool
	}{
		{year: 2000, want: true},
		{year: 1900, want: false},
		{year: 2024, want: true},
		{year: 2023, want: false},
		{year: 0, want: true},
		{year: -400, want: true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("year_%d", tt.year), func(t *testing.T) {
			if got := IsLeapYearHelper(tt.year); got != tt.want {
				t.Errorf("IsLeapYearHelper(%d) = %t, want %t", tt.year, got, tt.want)
			}
		})
	}
}

func TestDaysInMonthHelper(t *testing.T) {
	tests := []struct {
		name  string
		year  int64
		month int64
		want  int64
	}{
		{name: "leap February", year: 2000, month: 2, want: 29},
		{name: "common February", year: 1900, month: 2, want: 28},
		{name: "thirty-day month", year: 2024, month: 4, want: 30},
		{name: "November", year: 2024, month: 11, want: 30},
		{name: "thirty-one-day month", year: 2024, month: 1, want: 31},
		{name: "out-of-range month defaults to 31", year: 2024, month: 13, want: 31},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DaysInMonthHelper(tt.year, tt.month); got != tt.want {
				t.Errorf("DaysInMonthHelper(%d, %d) = %d, want %d", tt.year, tt.month, got, tt.want)
			}
		})
	}
}

func TestGetChronoLocalDateEpochDay(t *testing.T) {
	tests := []struct {
		name string
		obj  *object.Object
		want int64
	}{
		{name: "nil object", obj: nil, want: 0},
		{name: "nil field table", obj: &object.Object{}, want: 0},
		{name: "empty object", obj: object.MakeEmptyObject(), want: 0},
		{name: "epoch day", obj: dateTimeTestObject(map[string]object.Field{
			"epochDay": dateTimeTestField(int64(-1)),
		}), want: -1},
		{name: "epoch day takes precedence", obj: dateTimeTestObject(map[string]object.Field{
			"epochDay": dateTimeTestField(int32(42)), "year": dateTimeTestField(int64(2000)),
			"month": dateTimeTestField(int64(2)), "day": dateTimeTestField(int64(29)),
		}), want: 42},
		{name: "calendar date", obj: dateTimeTestObject(map[string]object.Field{
			"year": dateTimeTestField(int64(2000)), "month": dateTimeTestField(int64(2)),
			"day": dateTimeTestField(int64(29)),
		}), want: 11016},
		{name: "date defaults to January first", obj: dateTimeTestObject(map[string]object.Field{
			"year": dateTimeTestField(int(1970)),
		}), want: 0},
		{name: "invalid month and day default", obj: dateTimeTestObject(map[string]object.Field{
			"year": dateTimeTestField(int64(1970)), "month": dateTimeTestField("bad"),
			"day": dateTimeTestField(float64(2)),
		}), want: 0},
		{name: "invalid epoch day falls back to date", obj: dateTimeTestObject(map[string]object.Field{
			"epochDay": dateTimeTestField("bad"), "year": dateTimeTestField(int64(1969)),
			"month": dateTimeTestField(int64(12)), "day": dateTimeTestField(int64(31)),
		}), want: -1},
		{name: "invalid epoch day without year defaults to zero", obj: dateTimeTestObject(map[string]object.Field{
			"epochDay": dateTimeTestField("bad"),
		}), want: 0},
		{name: "valid epoch day without date", obj: dateTimeTestObject(map[string]object.Field{
			"epochDay": dateTimeTestField(int64(-2)),
		}), want: -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetChronoLocalDateEpochDay(tt.obj); got != tt.want {
				t.Errorf("GetChronoLocalDateEpochDay() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestGetChronoLocalDateYMD(t *testing.T) {
	tests := []struct {
		name  string
		obj   *object.Object
		year  int64
		month int64
		day   int64
	}{
		{name: "nil object", obj: nil, year: 1970, month: 1, day: 1},
		{name: "nil field table", obj: &object.Object{}, year: 1970, month: 1, day: 1},
		{name: "empty object", obj: object.MakeEmptyObject(), year: 1970, month: 1, day: 1},
		{name: "calendar date", obj: dateTimeTestObject(map[string]object.Field{
			"year": dateTimeTestField(int64(2024)), "month": dateTimeTestField(int64(2)),
			"day": dateTimeTestField(int64(29)),
		}), year: 2024, month: 2, day: 29},
		{name: "integer representations", obj: dateTimeTestObject(map[string]object.Field{
			"year": dateTimeTestField(int32(2001)), "month": dateTimeTestField(int(9)),
			"day": dateTimeTestField(int32(11)),
		}), year: 2001, month: 9, day: 11},
		{name: "invalid fields use defaults", obj: dateTimeTestObject(map[string]object.Field{
			"year": dateTimeTestField("bad"), "month": dateTimeTestField(float64(4)),
			"day": dateTimeTestField("bad"),
		}), year: 1970, month: 1, day: 1},
		{name: "positive epoch day", obj: dateTimeTestObject(map[string]object.Field{
			"epochDay": dateTimeTestField(int64(1)),
		}), year: 1970, month: 1, day: 2},
		{name: "negative epoch day", obj: dateTimeTestObject(map[string]object.Field{
			"epochDay": dateTimeTestField(int64(-1)),
		}), year: 1969, month: 12, day: 31},
		{name: "year takes precedence over epoch day", obj: dateTimeTestObject(map[string]object.Field{
			"epochDay": dateTimeTestField(int64(1)), "year": dateTimeTestField(int64(2000)),
			"month": dateTimeTestField(int64(3)), "day": dateTimeTestField(int64(4)),
		}), year: 2000, month: 3, day: 4},
		{name: "invalid present year suppresses epoch fallback", obj: dateTimeTestObject(map[string]object.Field{
			"epochDay": dateTimeTestField(int64(1)), "year": dateTimeTestField("bad"),
		}), year: 1970, month: 1, day: 1},
		{name: "invalid epoch day uses defaults", obj: dateTimeTestObject(map[string]object.Field{
			"epochDay": dateTimeTestField("bad"),
		}), year: 1970, month: 1, day: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			year, month, day := GetChronoLocalDateYMD(tt.obj)
			if year != tt.year || month != tt.month || day != tt.day {
				t.Errorf("GetChronoLocalDateYMD() = (%d, %d, %d), want (%d, %d, %d)",
					year, month, day, tt.year, tt.month, tt.day)
			}
		})
	}
}

func TestCreateChronoLocalDateTimeFromDateAndTime(t *testing.T) {
	date := object.MakeEmptyObject()
	time := object.MakeEmptyObject()
	dateTime := CreateChronoLocalDateTimeFromDateAndTime(date, time)
	if got := object.GoStringFromStringPoolIndex(dateTime.KlassName); got != types.ClassNameChronoLocalDateTime {
		t.Errorf("class = %q, want %q", got, types.ClassNameChronoLocalDateTime)
	}
	for name, want := range map[string]*object.Object{"date": date, "time": time} {
		field, ok := dateTime.FieldTable[name]
		if !ok {
			t.Errorf("missing %s field", name)
			continue
		}
		if field.Ftype != types.Ref {
			t.Errorf("%s field type = %q, want %q", name, field.Ftype, types.Ref)
		}
		if field.Fvalue != want {
			t.Errorf("%s field value = %v, want the original object", name, field.Fvalue)
		}
	}
}

func TestCreateChronoPeriodYMD(t *testing.T) {
	chronology := object.MakeEmptyObject()
	period := CreateChronoPeriodYMD(chronology, 2, 3, 5)
	if got := object.GoStringFromStringPoolIndex(period.KlassName); got != types.ClassNameChronoPeriod {
		t.Errorf("class = %q, want %q", got, types.ClassNameChronoPeriod)
	}
	for name, want := range map[string]int64{"years": 2, "months": 3, "days": 5} {
		field, ok := period.FieldTable[name]
		if !ok {
			t.Errorf("missing %s field", name)
			continue
		}
		if field.Ftype != types.Int {
			t.Errorf("%s field type = %q, want %q", name, field.Ftype, types.Int)
		}
		if got, ok := field.Fvalue.(int64); !ok || got != want {
			t.Errorf("%s field value = %v, want %d", name, field.Fvalue, want)
		}
	}
	chronologyField, ok := period.FieldTable["chronology"]
	if !ok {
		t.Fatal("missing chronology field")
	}
	if chronologyField.Ftype != types.Ref {
		t.Errorf("chronology field type = %q, want %q", chronologyField.Ftype, types.Ref)
	}
	if chronologyField.Fvalue != chronology {
		t.Errorf("chronology field value = %v, want the original object", chronologyField.Fvalue)
	}
}

func TestCreateChronology(t *testing.T) {
	chronology := CreateChronology("example/Chronology", "Example", "example")
	if got := object.GoStringFromStringPoolIndex(chronology.KlassName); got != "example/Chronology" {
		t.Errorf("class = %q, want %q", got, "example/Chronology")
	}
	for name, want := range map[string]string{"id": "Example", "calendarType": "example"} {
		field, ok := chronology.FieldTable[name]
		if !ok {
			t.Errorf("missing %s field", name)
			continue
		}
		if field.Ftype != types.Ref {
			t.Errorf("%s field type = %q, want %q", name, field.Ftype, types.Ref)
		}
		value, ok := field.Fvalue.(*object.Object)
		if !ok {
			t.Errorf("%s field value has type %T, want *object.Object", name, field.Fvalue)
			continue
		}
		if got := object.GoStringFromStringObject(value); got != want {
			t.Errorf("%s = %q, want %q", name, got, want)
		}
	}
}

func dateTimeTestObject(fields map[string]object.Field) *object.Object {
	obj := object.MakeEmptyObject()
	for name, field := range fields {
		obj.FieldTable[name] = field
	}
	return obj
}

func dateTimeTestField(value any) object.Field {
	return object.Field{Ftype: types.Int, Fvalue: value}
}
