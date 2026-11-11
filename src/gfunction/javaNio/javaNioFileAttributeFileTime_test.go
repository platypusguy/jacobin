/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaNio

import (
	"math"
	"testing"
	"time"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaTime"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
)

func TestFileTime_Registration(t *testing.T) {
	globals.InitGlobals("test")
	Load_Nio_File_Attribute_FileTime()

	expectedMethods := []string{
		"java/nio/file/attribute/FileTime.<clinit>()V",
		"java/nio/file/attribute/FileTime.from(JLjava/util/concurrent/TimeUnit;)Ljava/nio/file/attribute/FileTime;",
		"java/nio/file/attribute/FileTime.from(Ljava/time/Instant;)Ljava/nio/file/attribute/FileTime;",
		"java/nio/file/attribute/FileTime.fromMillis(J)Ljava/nio/file/attribute/FileTime;",
		"java/nio/file/attribute/FileTime.to(Ljava/util/concurrent/TimeUnit;)J",
		"java/nio/file/attribute/FileTime.toMillis()J",
		"java/nio/file/attribute/FileTime.toInstant()Ljava/time/Instant;",
		"java/nio/file/attribute/FileTime.compareTo(Ljava/nio/file/attribute/FileTime;)I",
		"java/nio/file/attribute/FileTime.compareTo(Ljava/lang/Object;)I",
		"java/nio/file/attribute/FileTime.equals(Ljava/lang/Object;)Z",
		"java/nio/file/attribute/FileTime.hashCode()I",
		"java/nio/file/attribute/FileTime.toString()Ljava/lang/String;",
	}

	for _, m := range expectedMethods {
		if _, ok := ghelpers.MethodSignatures[m]; !ok {
			t.Errorf("expected method signature %s to be registered", m)
		}
	}
}

func TestFileTime_FromMillisAndToMillis(t *testing.T) {
	globals.InitGlobals("test")
	Load_Nio_File_Attribute_FileTime()

	// fromMillis(123456789)
	ftObj := fileTimeFromMillis([]any{int64(123456789)}).(*object.Object)
	res := fileTimeToMillis([]any{ftObj})
	if res.(int64) != 123456789 {
		t.Errorf("expected toMillis() to return 123456789, got %v", res)
	}

	// fromMillis negative
	ftNeg := fileTimeFromMillis([]any{int64(-500)}).(*object.Object)
	resNeg := fileTimeToMillis([]any{ftNeg})
	if resNeg.(int64) != -500 {
		t.Errorf("expected toMillis() to return -500, got %v", resNeg)
	}

	// null self
	errBlk := fileTimeToMillis([]any{nil})
	if eb, ok := errBlk.(*ghelpers.GErrBlk); !ok || eb.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException on null self, got %v", errBlk)
	}
}

func TestFileTime_FromAndToTimeUnit(t *testing.T) {
	globals.InitGlobals("test")
	Load_Nio_File_Attribute_FileTime()

	// from(10, SECONDS)
	ftSec := fileTimeFromUnit([]any{int64(10), "SECONDS"}).(*object.Object)

	// to(MILLISECONDS) -> 10000
	toMillis := fileTimeToUnit([]any{ftSec, "MILLISECONDS"})
	if toMillis.(int64) != 10000 {
		t.Errorf("expected 10000 millis, got %v", toMillis)
	}

	// to(MICROSECONDS) -> 10000000
	toMicros := fileTimeToUnit([]any{ftSec, "MICROSECONDS"})
	if toMicros.(int64) != 10000000 {
		t.Errorf("expected 10000000 micros, got %v", toMicros)
	}

	// to(NANOSECONDS) -> 10000000000
	toNanos := fileTimeToUnit([]any{ftSec, "NANOSECONDS"})
	if toNanos.(int64) != 10000000000 {
		t.Errorf("expected 10000000000 nanos, got %v", toNanos)
	}

	// to(SECONDS) -> 10
	toSec := fileTimeToUnit([]any{ftSec, "SECONDS"})
	if toSec.(int64) != 10 {
		t.Errorf("expected 10 seconds, got %v", toSec)
	}

	// to(MINUTES) -> 0
	toMin := fileTimeToUnit([]any{ftSec, "MINUTES"})
	if toMin.(int64) != 0 {
		t.Errorf("expected 0 minutes, got %v", toMin)
	}

	// from(120, SECONDS) to(MINUTES) -> 2
	ft120 := fileTimeFromUnit([]any{int64(120), "SECONDS"}).(*object.Object)
	toMin2 := fileTimeToUnit([]any{ft120, "MINUTES"})
	if toMin2.(int64) != 2 {
		t.Errorf("expected 2 minutes, got %v", toMin2)
	}

	// from(2, DAYS) to(HOURS) -> 48
	ftDays := fileTimeFromUnit([]any{int64(2), "DAYS"}).(*object.Object)
	toHours := fileTimeToUnit([]any{ftDays, "HOURS"})
	if toHours.(int64) != 48 {
		t.Errorf("expected 48 hours, got %v", toHours)
	}

	// Saturation test on down-conversion
	ftHugeDays := fileTimeFromUnit([]any{int64(math.MaxInt64), "DAYS"}).(*object.Object)
	hugeNanos := fileTimeToUnit([]any{ftHugeDays, "NANOSECONDS"})
	if hugeNanos.(int64) != math.MaxInt64 {
		t.Errorf("expected saturated MaxInt64, got %v", hugeNanos)
	}

	ftNegHugeDays := fileTimeFromUnit([]any{int64(math.MinInt64), "DAYS"}).(*object.Object)
	hugeNegNanos := fileTimeToUnit([]any{ftNegHugeDays, "NANOSECONDS"})
	if hugeNegNanos.(int64) != math.MinInt64 {
		t.Errorf("expected saturated MinInt64, got %v", hugeNegNanos)
	}
}

func TestFileTime_TimeUnitEnumObject(t *testing.T) {
	globals.InitGlobals("test")
	Load_Nio_File_Attribute_FileTime()

	// TimeUnit simulated enum object with name and ordinal
	unitObj := object.MakeEmptyObjectWithClassName(&fileTimeClassName)
	unitObj.FieldTable["name"] = object.Field{Ftype: types.StringClassRef, Fvalue: object.StringObjectFromGoString("SECONDS")}
	unitObj.FieldTable["ordinal"] = object.Field{Ftype: types.Int, Fvalue: int64(3)}

	ft := fileTimeFromUnit([]any{int64(42), unitObj}).(*object.Object)
	sec := fileTimeToUnit([]any{ft, unitObj})
	if sec.(int64) != 42 {
		t.Errorf("expected 42, got %v", sec)
	}

	// Null unit
	errBlk := fileTimeFromUnit([]any{int64(42), nil})
	if eb, ok := errBlk.(*ghelpers.GErrBlk); !ok || eb.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for null unit in from(), got %v", errBlk)
	}

	// Invalid unit
	errBlk2 := fileTimeFromUnit([]any{int64(42), "LIGHTYEARS"})
	if eb, ok := errBlk2.(*ghelpers.GErrBlk); !ok || eb.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException for invalid unit in from(), got %v", errBlk2)
	}
}

func TestFileTime_FromInstantAndToInstant(t *testing.T) {
	globals.InitGlobals("test")
	Load_Nio_File_Attribute_FileTime()

	inst := javaTime.CreateInstant(1600000000, 500000000)
	ft := fileTimeFromInstant([]any{inst}).(*object.Object)

	// toMillis
	milli := fileTimeToMillis([]any{ft})
	expectedMilli := int64(1600000000)*1000 + 500
	if milli.(int64) != expectedMilli {
		t.Errorf("expected %d milli, got %v", expectedMilli, milli)
	}

	// to(NANOSECONDS)
	nanos := fileTimeToUnit([]any{ft, "NANOSECONDS"})
	expectedNanos := int64(1600000000)*1000000000 + 500000000
	if nanos.(int64) != expectedNanos {
		t.Errorf("expected %d nanos, got %v", expectedNanos, nanos)
	}

	// toInstant
	instRes := fileTimeToInstant([]any{ft}).(*object.Object)
	s, n := extractInstantSecondsAndNanos(instRes)
	if s != 1600000000 || n != 500000000 {
		t.Errorf("expected (1600000000, 500000000), got (%d, %d)", s, n)
	}

	// toInstant from TimeUnit
	ftFromUnit := fileTimeFromUnit([]any{int64(1500), "MILLISECONDS"}).(*object.Object)
	instFromUnit := fileTimeToInstant([]any{ftFromUnit}).(*object.Object)
	sU, nU := extractInstantSecondsAndNanos(instFromUnit)
	if sU != 1 || nU != 500000000 {
		t.Errorf("expected (1, 500000000), got (%d, %d)", sU, nU)
	}

	// null instant check
	errBlk := fileTimeFromInstant([]any{nil})
	if eb, ok := errBlk.(*ghelpers.GErrBlk); !ok || eb.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException on null instant, got %v", errBlk)
	}
}

func TestFileTime_CompareToAndEquals(t *testing.T) {
	globals.InitGlobals("test")
	Load_Nio_File_Attribute_FileTime()

	// Same unit
	ft1 := fileTimeFromMillis([]any{int64(1000)}).(*object.Object)
	ft2 := fileTimeFromMillis([]any{int64(1000)}).(*object.Object)
	ft3 := fileTimeFromMillis([]any{int64(2000)}).(*object.Object)

	if cmp := fileTimeCompareTo([]any{ft1, ft2}); cmp.(int64) != 0 {
		t.Errorf("expected ft1 == ft2 (0), got %v", cmp)
	}
	if cmp := fileTimeCompareTo([]any{ft1, ft3}); cmp.(int64) != -1 {
		t.Errorf("expected ft1 < ft3 (-1), got %v", cmp)
	}
	if cmp := fileTimeCompareTo([]any{ft3, ft1}); cmp.(int64) != 1 {
		t.Errorf("expected ft3 > ft1 (1), got %v", cmp)
	}

	if eq := fileTimeEquals([]any{ft1, ft2}); eq.(int64) != types.JavaBoolTrue {
		t.Errorf("expected ft1.equals(ft2) to be true")
	}
	if eq := fileTimeEquals([]any{ft1, ft3}); eq.(int64) != types.JavaBoolFalse {
		t.Errorf("expected ft1.equals(ft3) to be false")
	}

	// Cross-unit comparison: 1000 MILLISECONDS vs 1 SECONDS vs 1_000_000_000 NANOSECONDS vs Instant(1, 0)
	ftSec := fileTimeFromUnit([]any{int64(1), "SECONDS"}).(*object.Object)
	ftNanos := fileTimeFromUnit([]any{int64(1000000000), "NANOSECONDS"}).(*object.Object)
	ftInst := fileTimeFromInstant([]any{javaTime.CreateInstant(1, 0)}).(*object.Object)

	if cmp := fileTimeCompareTo([]any{ft1, ftSec}); cmp.(int64) != 0 {
		t.Errorf("expected 1000 millis == 1 second, got %v", cmp)
	}
	if eq := fileTimeEquals([]any{ft1, ftSec}); eq.(int64) != types.JavaBoolTrue {
		t.Errorf("expected 1000 millis equals 1 second")
	}
	if eq := fileTimeEquals([]any{ftSec, ftNanos}); eq.(int64) != types.JavaBoolTrue {
		t.Errorf("expected 1 second equals 10^9 nanos")
	}
	if eq := fileTimeEquals([]any{ftNanos, ftInst}); eq.(int64) != types.JavaBoolTrue {
		t.Errorf("expected 10^9 nanos equals Instant(1, 0)")
	}

	// Null and type checks
	if eq := fileTimeEquals([]any{ft1, nil}); eq.(int64) != types.JavaBoolFalse {
		t.Errorf("expected ft.equals(null) to be false")
	}
	strObj := object.StringObjectFromGoString("not a FileTime")
	if eq := fileTimeEquals([]any{ft1, strObj}); eq.(int64) != types.JavaBoolFalse {
		t.Errorf("expected ft.equals(string) to be false")
	}

	// compareToObject bridge
	if cmp := fileTimeCompareToObject([]any{ft1, ftSec}); cmp.(int64) != 0 {
		t.Errorf("expected compareToObject to succeed with 0")
	}
	errBlk := fileTimeCompareToObject([]any{ft1, strObj})
	if eb, ok := errBlk.(*ghelpers.GErrBlk); !ok || eb.ExceptionType != excNames.ClassCastException {
		t.Errorf("expected ClassCastException on compareToObject with non-FileTime, got %v", errBlk)
	}
}

func TestFileTime_HashCode(t *testing.T) {
	globals.InitGlobals("test")
	Load_Nio_File_Attribute_FileTime()

	ft1 := fileTimeFromMillis([]any{int64(1000)}).(*object.Object)
	ftSec := fileTimeFromUnit([]any{int64(1), "SECONDS"}).(*object.Object)
	ftInst := fileTimeFromInstant([]any{javaTime.CreateInstant(1, 0)}).(*object.Object)

	h1 := fileTimeHashCode([]any{ft1}).(int64)
	h2 := fileTimeHashCode([]any{ftSec}).(int64)
	h3 := fileTimeHashCode([]any{ftInst}).(int64)

	if h1 != h2 || h2 != h3 {
		t.Errorf("expected matching hash codes for equal FileTimes, got h1=%d, h2=%d, h3=%d", h1, h2, h3)
	}
}

func TestFileTime_ToString(t *testing.T) {
	globals.InitGlobals("test")
	Load_Nio_File_Attribute_FileTime()

	testCases := []struct {
		name     string
		ft       *object.Object
		expected string
	}{
		{
			name:     "from(0, SECONDS)",
			ft:       fileTimeFromUnit([]any{int64(0), "SECONDS"}).(*object.Object),
			expected: "1970-01-01T00:00:00Z",
		},
		{
			name:     "from(0, MILLISECONDS)",
			ft:       fileTimeFromMillis([]any{int64(0)}).(*object.Object),
			expected: "1970-01-01T00:00:00.000Z",
		},
		{
			name:     "from(1234, MILLISECONDS)",
			ft:       fileTimeFromMillis([]any{int64(1234)}).(*object.Object),
			expected: "1970-01-01T00:00:01.234Z",
		},
		{
			name:     "from(-1, MILLISECONDS)",
			ft:       fileTimeFromMillis([]any{int64(-1)}).(*object.Object),
			expected: "1969-12-31T23:59:59.999Z",
		},
		{
			name:     "from(1234567, MICROSECONDS)",
			ft:       fileTimeFromUnit([]any{int64(1234567), "MICROSECONDS"}).(*object.Object),
			expected: "1970-01-01T00:00:01.234567Z",
		},
		{
			name:     "from(1234567890, NANOSECONDS)",
			ft:       fileTimeFromUnit([]any{int64(1234567890), "NANOSECONDS"}).(*object.Object),
			expected: "1970-01-01T00:00:01.234567890Z",
		},
		{
			name:     "from(Instant(0, 0))",
			ft:       fileTimeFromInstant([]any{javaTime.CreateInstant(0, 0)}).(*object.Object),
			expected: "1970-01-01T00:00:00Z",
		},
		{
			name:     "from(Instant(0, 123000000))",
			ft:       fileTimeFromInstant([]any{javaTime.CreateInstant(0, 123000000)}).(*object.Object),
			expected: "1970-01-01T00:00:00.123Z",
		},
		{
			name:     "from(Instant(0, 123456000))",
			ft:       fileTimeFromInstant([]any{javaTime.CreateInstant(0, 123456000)}).(*object.Object),
			expected: "1970-01-01T00:00:00.123456Z",
		},
		{
			name:     "from(Instant(0, 123456789))",
			ft:       fileTimeFromInstant([]any{javaTime.CreateInstant(0, 123456789)}).(*object.Object),
			expected: "1970-01-01T00:00:00.123456789Z",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			strObj := fileTimeToString([]any{tc.ft}).(*object.Object)
			str := object.GoStringFromStringObject(strObj)
			if str != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, str)
			}
		})
	}
}

func TestFileTime_NewFileTimeHelper(t *testing.T) {
	globals.InitGlobals("test")
	Load_Nio_File_Attribute_FileTime()

	goTime := time.Date(2024, 6, 15, 12, 30, 45, 123456789, time.UTC)
	ft := newFileTime(goTime)

	toMillis := fileTimeToMillis([]any{ft}).(int64)
	if toMillis != goTime.UnixMilli() {
		t.Errorf("expected %d millis, got %d", goTime.UnixMilli(), toMillis)
	}

	toNanos := fileTimeToUnit([]any{ft, "NANOSECONDS"}).(int64)
	if toNanos != goTime.UnixNano() {
		t.Errorf("expected %d nanos, got %d", goTime.UnixNano(), toNanos)
	}
}
