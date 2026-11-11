/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"math"
	"testing"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
)

func assertInt64Eq(t *testing.T, got any, want int64, msg string) {
	t.Helper()
	gi, ok := got.(int64)
	if !ok {
		t.Fatalf("%s: expected int64 result, got %T (%v)", msg, got, got)
	}
	if gi != want {
		t.Fatalf("%s: expected %d, got %d", msg, want, gi)
	}
}

func assertErr(t *testing.T, got any, expectedExc int, msg string) {
	t.Helper()
	geb, ok := got.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("%s: expected *ghelpers.GErrBlk, got %T (%v)", msg, got, got)
	}
	if geb.ExceptionType != expectedExc {
		t.Fatalf("%s: expected exception type %d, got %d (msg: %s)", msg, expectedExc, geb.ExceptionType, geb.ErrMsg)
	}
}

func TestTimeUnit_MethodRegistration(t *testing.T) {
	Load_Util_Concurrent_TimeUnit()

	expectedMethods := []struct {
		sig   string
		slots int
	}{
		{"java/util/concurrent/TimeUnit.<clinit>()V", 0},
		{"java/util/concurrent/TimeUnit.values()[Ljava/util/concurrent/TimeUnit;", 0},
		{"java/util/concurrent/TimeUnit.valueOf(Ljava/lang/String;)Ljava/util/concurrent/TimeUnit;", 1},
		{"java/util/concurrent/TimeUnit.of(Ljava/time/temporal/ChronoUnit;)Ljava/util/concurrent/TimeUnit;", 1},
		{"java/util/concurrent/TimeUnit.convert(JLjava/util/concurrent/TimeUnit;)J", 2},
		{"java/util/concurrent/TimeUnit.convert(Ljava/time/Duration;)J", 1},
		{"java/util/concurrent/TimeUnit.toNanos(J)J", 1},
		{"java/util/concurrent/TimeUnit.toMicros(J)J", 1},
		{"java/util/concurrent/TimeUnit.toMillis(J)J", 1},
		{"java/util/concurrent/TimeUnit.toSeconds(J)J", 1},
		{"java/util/concurrent/TimeUnit.toMinutes(J)J", 1},
		{"java/util/concurrent/TimeUnit.toHours(J)J", 1},
		{"java/util/concurrent/TimeUnit.toDays(J)J", 1},
		{"java/util/concurrent/TimeUnit.toChronoUnit()Ljava/time/temporal/ChronoUnit;", 0},
		{"java/util/concurrent/TimeUnit.sleep(J)V", 1},
		{"java/util/concurrent/TimeUnit.timedWait(Ljava/lang/Object;J)V", 2},
		{"java/util/concurrent/TimeUnit.timedJoin(Ljava/lang/Thread;J)V", 2},
		{"java/util/concurrent/TimeUnit.name()Ljava/lang/String;", 0},
		{"java/util/concurrent/TimeUnit.ordinal()I", 0},
		{"java/util/concurrent/TimeUnit.toString()Ljava/lang/String;", 0},
		{"java/util/concurrent/TimeUnit.compareTo(Ljava/lang/Enum;)I", 1},
		{"java/util/concurrent/TimeUnit.compareTo(Ljava/lang/Object;)I", 1},
		{"java/util/concurrent/TimeUnit.equals(Ljava/lang/Object;)Z", 1},
		{"java/util/concurrent/TimeUnit.hashCode()I", 0},
	}

	for _, m := range expectedMethods {
		meth, exists := ghelpers.MethodSignatures[m.sig]
		if !exists {
			t.Errorf("Method signature missing: %s", m.sig)
			continue
		}
		if meth.ParamSlots != m.slots {
			t.Errorf("Method %s ParamSlots = %d, want %d", m.sig, meth.ParamSlots, m.slots)
		}
		if meth.GFunction == nil {
			t.Errorf("Method %s GFunction is nil", m.sig)
		}
	}
}

func TestTimeUnit_ValuesAndValueOf(t *testing.T) {
	globals.InitStringPool()
	Load_Util_Concurrent_TimeUnit()
	timeUnitClinit(nil)

	// values()
	valsRes := timeUnitValues(nil)
	arrObj, ok := valsRes.(*object.Object)
	if !ok {
		t.Fatalf("values() did not return an Object: %T", valsRes)
	}
	fld, ok := arrObj.FieldTable["value"]
	if !ok {
		t.Fatalf("values() array missing 'value' field")
	}
	slice, ok := fld.Fvalue.([]*object.Object)
	if !ok || len(slice) != 7 {
		t.Fatalf("values() slice length = %d, want 7", len(slice))
	}

	expectedNames := []string{
		"NANOSECONDS",
		"MICROSECONDS",
		"MILLISECONDS",
		"SECONDS",
		"MINUTES",
		"HOURS",
		"DAYS",
	}

	for i, expected := range expectedNames {
		inst := slice[i]
		if inst == nil {
			t.Fatalf("values()[%d] is nil", i)
		}
		nameRes := timeUnitName([]any{inst})
		strObj, ok := nameRes.(*object.Object)
		if !ok {
			t.Fatalf("name() did not return *object.Object")
		}
		if gotName := object.GoStringFromStringObject(strObj); gotName != expected {
			t.Errorf("values()[%d] name = %s, want %s", i, gotName, expected)
		}

		ordRes := timeUnitOrdinal([]any{inst})
		if ordRes.(int64) != int64(i) {
			t.Errorf("values()[%d] ordinal = %v, want %d", i, ordRes, i)
		}

		// Check valueOf(String)
		strArg := object.StringObjectFromGoString(expected)
		voRes := timeUnitValueOfString([]any{strArg})
		if voRes != inst {
			t.Errorf("valueOf(%s) did not return singleton instance", expected)
		}
	}

	// valueOf errors
	assertErr(t, timeUnitValueOfString([]any{nil}), excNames.NullPointerException, "valueOf(null)")
	assertErr(t, timeUnitValueOfString([]any{object.Null}), excNames.NullPointerException, "valueOf(Null)")
	assertErr(t, timeUnitValueOfString([]any{object.StringObjectFromGoString("INVALID_UNIT")}), excNames.IllegalArgumentException, "valueOf(INVALID)")
}

func TestTimeUnit_OfChronoUnit(t *testing.T) {
	globals.InitStringPool()
	Load_Util_Concurrent_TimeUnit()
	timeUnitClinit(nil)

	mappings := []struct {
		cuName string
		tuName string
	}{
		{"NANOS", "NANOSECONDS"},
		{"MICROS", "MICROSECONDS"},
		{"MILLIS", "MILLISECONDS"},
		{"SECONDS", "SECONDS"},
		{"MINUTES", "MINUTES"},
		{"HOURS", "HOURS"},
		{"DAYS", "DAYS"},
	}

	for _, m := range mappings {
		cuObj := ghelpers.CreateTemporalUnit(m.cuName)
		res := timeUnitOfChronoUnit([]any{cuObj})
		inst, ok := res.(*object.Object)
		if !ok || inst == nil {
			t.Fatalf("of(%s) failed to return object", m.cuName)
		}
		nameObj := timeUnitName([]any{inst}).(*object.Object)
		if got := object.GoStringFromStringObject(nameObj); got != m.tuName {
			t.Errorf("of(%s) = %s, want %s", m.cuName, got, m.tuName)
		}
	}

	// Invalid ChronoUnits
	weeksObj := ghelpers.CreateTemporalUnit("WEEKS")
	assertErr(t, timeUnitOfChronoUnit([]any{weeksObj}), excNames.IllegalArgumentException, "of(WEEKS)")
	assertErr(t, timeUnitOfChronoUnit([]any{nil}), excNames.NullPointerException, "of(null)")
	assertErr(t, timeUnitOfChronoUnit([]any{object.Null}), excNames.NullPointerException, "of(Null)")
}

func TestTimeUnit_ToChronoUnit(t *testing.T) {
	globals.InitStringPool()
	Load_Util_Concurrent_TimeUnit()
	timeUnitClinit(nil)

	for _, name := range timeUnitNames {
		inst := GetTimeUnitInstance(name)
		res := timeUnitToChronoUnit([]any{inst})
		cuObj, ok := res.(*object.Object)
		if !ok || cuObj == nil {
			t.Fatalf("toChronoUnit() failed for %s", name)
		}
		fld := cuObj.FieldTable["name"]
		strObj := fld.Fvalue.(*object.Object)
		gotCU := object.GoStringFromStringObject(strObj)
		expectedCU := timeUnitToChronoUnitName[name]
		if gotCU != expectedCU {
			t.Errorf("%s.toChronoUnit() = %s, want %s", name, gotCU, expectedCU)
		}
	}
}

func TestTimeUnit_Conversions(t *testing.T) {
	globals.InitStringPool()
	Load_Util_Concurrent_TimeUnit()
	timeUnitClinit(nil)

	sec := GetTimeUnitInstance("SECONDS")
	milli := GetTimeUnitInstance("MILLISECONDS")
	micro := GetTimeUnitInstance("MICROSECONDS")
	nano := GetTimeUnitInstance("NANOSECONDS")
	min := GetTimeUnitInstance("MINUTES")
	hr := GetTimeUnitInstance("HOURS")
	day := GetTimeUnitInstance("DAYS")

	// SECONDS conversions
	assertInt64Eq(t, timeUnitToNanos([]any{sec, int64(5)}), 5000000000, "5 sec to nanos")
	assertInt64Eq(t, timeUnitToMicros([]any{sec, int64(5)}), 5000000, "5 sec to micros")
	assertInt64Eq(t, timeUnitToMillis([]any{sec, int64(5)}), 5000, "5 sec to millis")
	assertInt64Eq(t, timeUnitToSeconds([]any{sec, int64(5)}), 5, "5 sec to seconds")
	assertInt64Eq(t, timeUnitToMinutes([]any{sec, int64(125)}), 2, "125 sec to minutes")
	assertInt64Eq(t, timeUnitToHours([]any{sec, int64(7200)}), 2, "7200 sec to hours")
	assertInt64Eq(t, timeUnitToDays([]any{sec, int64(86400 * 3)}), 3, "3 days in sec to days")

	// DAYS conversions
	assertInt64Eq(t, timeUnitToHours([]any{day, int64(2)}), 48, "2 days to hours")
	assertInt64Eq(t, timeUnitToMinutes([]any{day, int64(1)}), 1440, "1 day to minutes")
	assertInt64Eq(t, timeUnitToSeconds([]any{day, int64(1)}), 86400, "1 day to seconds")
	assertInt64Eq(t, timeUnitToMillis([]any{day, int64(1)}), 86400000, "1 day to millis")

	// NANOSECONDS to larger units
	assertInt64Eq(t, timeUnitToMicros([]any{nano, int64(1500)}), 1, "1500 nanos to micros")
	assertInt64Eq(t, timeUnitToMillis([]any{nano, int64(1999999)}), 1, "1999999 nanos to millis")
	assertInt64Eq(t, timeUnitToSeconds([]any{nano, int64(2000000000)}), 2, "2000000000 nanos to seconds")

	// Negative durations
	assertInt64Eq(t, timeUnitToMillis([]any{sec, int64(-3)}), -3000, "-3 sec to millis")
	assertInt64Eq(t, timeUnitToSeconds([]any{milli, int64(-2500)}), -2, "-2500 millis to seconds")

	// convert(long, TimeUnit)
	assertInt64Eq(t, timeUnitConvert([]any{sec, int64(1000), milli}), 1, "convert 1000 millis to seconds")
	assertInt64Eq(t, timeUnitConvert([]any{milli, int64(2), sec}), 2000, "convert 2 seconds to millis")
	assertInt64Eq(t, timeUnitConvert([]any{hr, int64(180), min}), 3, "convert 180 min to hours")
	assertInt64Eq(t, timeUnitConvert([]any{micro, int64(5000), nano}), 5, "convert 5000 nanos to micros")
}

func TestTimeUnit_OverflowSaturation(t *testing.T) {
	globals.InitStringPool()
	Load_Util_Concurrent_TimeUnit()
	timeUnitClinit(nil)

	day := GetTimeUnitInstance("DAYS")
	sec := GetTimeUnitInstance("SECONDS")

	// Days to Nanos overflow
	resMax := timeUnitToNanos([]any{day, math.MaxInt64})
	assertInt64Eq(t, resMax, math.MaxInt64, "Days.toNanos(MAX_VALUE) saturates at MAX_VALUE")

	resMin := timeUnitToNanos([]any{day, math.MinInt64})
	assertInt64Eq(t, resMin, math.MinInt64, "Days.toNanos(MIN_VALUE) saturates at MIN_VALUE")

	// Seconds to Nanos overflow
	resSecMax := timeUnitToNanos([]any{sec, math.MaxInt64})
	assertInt64Eq(t, resSecMax, math.MaxInt64, "Seconds.toNanos(MAX_VALUE) saturates at MAX_VALUE")

	resSecMin := timeUnitToNanos([]any{sec, math.MinInt64})
	assertInt64Eq(t, resSecMin, math.MinInt64, "Seconds.toNanos(MIN_VALUE) saturates at MIN_VALUE")
}

func TestTimeUnit_ConvertDuration(t *testing.T) {
	globals.InitStringPool()
	Load_Util_Concurrent_TimeUnit()
	timeUnitClinit(nil)

	sec := GetTimeUnitInstance("SECONDS")
	milli := GetTimeUnitInstance("MILLISECONDS")
	nano := GetTimeUnitInstance("NANOSECONDS")

	durClass := "java/time/Duration"
	durObj := object.MakeEmptyObjectWithClassName(&durClass)
	durObj.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: int64(12)}
	durObj.FieldTable["nanos"] = object.Field{Ftype: types.Long, Fvalue: int64(345000000)}

	assertInt64Eq(t, timeUnitConvertDuration([]any{sec, durObj}), 12, "convert(Duration) to seconds")
	assertInt64Eq(t, timeUnitConvertDuration([]any{milli, durObj}), 12345, "convert(Duration) to millis")
	assertInt64Eq(t, timeUnitConvertDuration([]any{nano, durObj}), 12345000000, "convert(Duration) to nanos")

	// Error checks
	assertErr(t, timeUnitConvertDuration([]any{sec, nil}), excNames.NullPointerException, "convert(null Duration)")
	assertErr(t, timeUnitConvertDuration([]any{sec, object.Null}), excNames.NullPointerException, "convert(Null Duration)")
}

func TestTimeUnit_CompareEqualsHashToString(t *testing.T) {
	globals.InitStringPool()
	Load_Util_Concurrent_TimeUnit()
	timeUnitClinit(nil)

	nano := GetTimeUnitInstance("NANOSECONDS")
	sec := GetTimeUnitInstance("SECONDS")
	sec2 := GetTimeUnitInstance("SECONDS")
	day := GetTimeUnitInstance("DAYS")

	// compareTo
	cmp1 := timeUnitCompareTo([]any{nano, sec})
	if cmp1.(int64) >= 0 {
		t.Errorf("NANOSECONDS.compareTo(SECONDS) expected < 0, got %v", cmp1)
	}

	cmp2 := timeUnitCompareTo([]any{sec, nano})
	if cmp2.(int64) <= 0 {
		t.Errorf("SECONDS.compareTo(NANOSECONDS) expected > 0, got %v", cmp2)
	}

	cmp3 := timeUnitCompareTo([]any{sec, sec2})
	if cmp3.(int64) != 0 {
		t.Errorf("SECONDS.compareTo(SECONDS) expected 0, got %v", cmp3)
	}

	// equals
	eqTrue := timeUnitEquals([]any{sec, sec2})
	if eqTrue.(int64) != types.JavaBoolTrue {
		t.Errorf("SECONDS.equals(SECONDS) expected true, got %v", eqTrue)
	}

	eqFalse := timeUnitEquals([]any{sec, day})
	if eqFalse.(int64) != types.JavaBoolFalse {
		t.Errorf("SECONDS.equals(DAYS) expected false, got %v", eqFalse)
	}

	// toString
	strObj := timeUnitToString([]any{sec}).(*object.Object)
	if s := object.GoStringFromStringObject(strObj); s != "SECONDS" {
		t.Errorf("SECONDS.toString() = %s, want SECONDS", s)
	}

	// hashCode
	h := timeUnitHashCode([]any{sec})
	if h.(int64) != 3 {
		t.Errorf("SECONDS.hashCode() = %v, want 3", h)
	}
}

func TestTimeUnit_SleepTimedWaitTimedJoin(t *testing.T) {
	globals.InitStringPool()
	Load_Util_Concurrent_TimeUnit()
	timeUnitClinit(nil)

	milli := GetTimeUnitInstance("MILLISECONDS")

	// sleep(0) should return immediately with void
	resSleep := timeUnitSleep([]any{milli, int64(0)})
	if resSleep != nil {
		t.Errorf("sleep(0) expected nil, got %v", resSleep)
	}

	// sleep(1)
	resSleep1 := timeUnitSleep([]any{milli, int64(1)})
	if resSleep1 != nil {
		t.Errorf("sleep(1) expected nil, got %v", resSleep1)
	}

	dummyObj := object.MakeEmptyObject()
	dummyThread := object.MakeEmptyObject()

	// timedWait
	resWait := timeUnitTimedWait([]any{milli, dummyObj, int64(10)})
	if resWait != nil {
		t.Errorf("timedWait expected nil, got %v", resWait)
	}
	assertErr(t, timeUnitTimedWait([]any{milli, nil, int64(10)}), excNames.NullPointerException, "timedWait null object")
	assertErr(t, timeUnitTimedWait([]any{milli, dummyObj, int64(-5)}), excNames.IllegalArgumentException, "timedWait negative timeout")

	// timedJoin
	resJoin := timeUnitTimedJoin([]any{milli, dummyThread, int64(10)})
	if resJoin != nil {
		t.Errorf("timedJoin expected nil, got %v", resJoin)
	}
	assertErr(t, timeUnitTimedJoin([]any{milli, nil, int64(10)}), excNames.NullPointerException, "timedJoin null thread")
	assertErr(t, timeUnitTimedJoin([]any{milli, dummyThread, int64(-5)}), excNames.IllegalArgumentException, "timedJoin negative timeout")
}

func TestTimeUnit_ErrorHandling(t *testing.T) {
	globals.InitStringPool()
	Load_Util_Concurrent_TimeUnit()
	timeUnitClinit(nil)

	badUnit := object.StringObjectFromGoString("INVALID_UNIT")

	assertErr(t, timeUnitToNanos([]any{badUnit, int64(10)}), excNames.IllegalArgumentException, "toNanos invalid unit")
	assertErr(t, timeUnitToMillis([]any{nil, int64(10)}), excNames.NullPointerException, "toMillis null unit")
	assertErr(t, timeUnitConvert([]any{badUnit, int64(10), "SECONDS"}), excNames.IllegalArgumentException, "convert invalid target unit")
	assertErr(t, timeUnitConvert([]any{"SECONDS", int64(10), badUnit}), excNames.IllegalArgumentException, "convert invalid source unit")
	assertErr(t, timeUnitName([]any{nil}), excNames.NullPointerException, "name null unit")
	assertErr(t, timeUnitOrdinal([]any{nil}), excNames.NullPointerException, "ordinal null unit")
}
