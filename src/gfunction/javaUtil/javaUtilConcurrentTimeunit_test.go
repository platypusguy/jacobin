package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"testing"
)

// helper to build a String object for the TimeUnit name
func tu(unit string) *object.Object {
	return object.StringObjectFromGoString(unit)
}

// helper to assert an int64 equality
func assertInt64Equal(t *testing.T, got interface{}, want int64, msg string) {
	t.Helper()
	gi, ok := got.(int64)
	if !ok {
		t.Fatalf("%s: expected int64 result, got %T", msg, got)
	}
	if gi != want {
		t.Fatalf("%s: expected %d, got %d", msg, want, gi)
	}
}

// helper to assert an error block with expected exception type
func assertErrType(t *testing.T, got interface{}, expected int) {
	t.Helper()
	geb, ok := got.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected error block, got %T", got)
	}
	if geb.ExceptionType != expected {
		t.Fatalf("expected exception type %d, got %d", expected, geb.ExceptionType)
	}
}

func TestTimeUnit_IdentityConversions(t *testing.T) {
	globals.InitStringPool()

	// toMillis with MILLISECONDS should be identity
	assertInt64Equal(t, toMillis([]interface{}{tu(MILLISECONDS), int64(12345)}), 12345, "toMillis identity")

	// toSeconds with SECONDS should be identity
	assertInt64Equal(t, toSeconds([]interface{}{tu(SECONDS), int64(-77)}), -77, "toSeconds identity")

	// toMinutes with MINUTES should be identity
	assertInt64Equal(t, toMinutes([]interface{}{tu(MINUTES), int64(0)}), 0, "toMinutes identity")

	// toHours with HOURS should be identity
	assertInt64Equal(t, toHours([]interface{}{tu(HOURS), int64(9)}), 9, "toHours identity")

	// toDays with DAYS should be identity
	assertInt64Equal(t, toDays([]interface{}{tu(DAYS), int64(1)}), 1, "toDays identity")
}

func TestTimeUnit_InvalidUnit_Errors(t *testing.T) {
	globals.InitStringPool()

	bad := object.StringObjectFromGoString("WEEKS")

	if res := toMillis([]interface{}{bad, int64(1)}); res != nil {
		assertErrType(t, res, excNames.IllegalArgumentException)
	} else {
		t.Fatalf("expected error for invalid unit in toMillis")
	}

	if res := toSeconds([]interface{}{bad, int64(1)}); res != nil {
		assertErrType(t, res, excNames.IllegalArgumentException)
	} else {
		t.Fatalf("expected error for invalid unit in toSeconds")
	}

	if res := toMinutes([]interface{}{bad, int64(1)}); res != nil {
		assertErrType(t, res, excNames.IllegalArgumentException)
	} else {
		t.Fatalf("expected error for invalid unit in toMinutes")
	}

	if res := toHours([]interface{}{bad, int64(1)}); res != nil {
		assertErrType(t, res, excNames.IllegalArgumentException)
	} else {
		t.Fatalf("expected error for invalid unit in toHours")
	}

	if res := toDays([]interface{}{bad, int64(1)}); res != nil {
		assertErrType(t, res, excNames.IllegalArgumentException)
	} else {
		t.Fatalf("expected error for invalid unit in toDays")
	}
}
