package gfunction

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

// Helpers
func newTZObj(t *testing.T) *object.Object {
	t.Helper()
	className := "java/util/TimeZone"
	obj := object.MakeEmptyObjectWithClassName(&className)
	if obj == nil {
		t.Fatalf("failed to allocate TimeZone object")
	}
	return obj
}

func str(s string) *object.Object { return object.StringObjectFromGoString(s) }

func assertJavaBoolTZ(t *testing.T, got interface{}, want int64, msg string) {
	t.Helper()
	b, ok := got.(int64)
	if !ok {
		if geb, ok := got.(*GErrBlk); ok {
			t.Fatalf("%s: expected Java boolean, got GErrBlk %d (%s)", msg, geb.ExceptionType, geb.ErrMsg)
		}
		t.Fatalf("%s: expected Java boolean (int64), got %T", msg, got)
	}
	if b != want {
		t.Fatalf("%s: expected %d, got %d", msg, want, b)
	}
}

func TestTimeZone_MethodRegistration(t *testing.T) {
	globals.InitStringPool()
	MethodSignatures = make(map[string]GMeth)
	Load_Util_TimeZone()

	cases := []struct {
		key   string
		slots int
	}{
		{"java/util/TimeZone.<init>()V", 0},
		{"java/util/TimeZone.clone()Ljava/lang/Object;", 0},
		{"java/util/TimeZone.getAvailableIDs()[Ljava/lang/String;", 0},
		{"java/util/TimeZone.getAvailableIDs(I)[Ljava/lang/String;", 1},
		{"java/util/TimeZone.getDisplayName()Ljava/lang/String;", 0},
		{"java/util/TimeZone.getDSTSavings()I", 0},
		{"java/util/TimeZone.getID()Ljava/lang/String;", 0},
		{"java/util/TimeZone.getOffset(IIIII)I", 6},
		{"java/util/TimeZone.getOffset(J)I", 1},
		{"java/util/TimeZone.getRawOffset()I", 0},
		{"java/util/TimeZone.getTimeZone(Ljava/lang/String;)Ljava/util/TimeZone;", 1},
		{"java/util/TimeZone.hasSameRules(Ljava/util/TimeZone;)Z", 1},
		{"java/util/TimeZone.inDaylightTime(Ljava/util/Date;)Z", 1},
		{"java/util/TimeZone.observesDaylightTime()Z", 0},
		{"java/util/TimeZone.setDefault(Ljava/util/TimeZone;)V", 1},
		{"java/util/TimeZone.setID(Ljava/lang/String;)V", 1},
		{"java/util/TimeZone.setRawOffset(I)V", 1},
	}
	for _, c := range cases {
		gm, ok := MethodSignatures[c.key]
		if !ok {
			t.Fatalf("method not registered: %s", c.key)
		}
		if gm.ParamSlots != c.slots {
			t.Fatalf("ParamSlots mismatch for %s: want %d got %d", c.key, c.slots, gm.ParamSlots)
		}
		if gm.GFunction == nil {
			t.Fatalf("GFunction is nil for %s", c.key)
		}
	}
}

func TestTimeZone_Init_DefaultFields(t *testing.T) {
	globals.InitStringPool()
	obj := newTZObj(t)
	if ret := tzInit([]interface{}{obj}); ret != nil {
		t.Fatalf("tzInit returned error: %v", ret)
	}
	// id should be "UTC"
	id := tzGetID([]interface{}{obj}).(*object.Object)
	if object.GoStringFromStringObject(id) != "UTC" {
		t.Fatalf("getID expected UTC, got %q", object.GoStringFromStringObject(id))
	}
	// displayName should match id
	dn := tzGetDisplayName([]interface{}{obj}).(*object.Object)
	if object.GoStringFromStringObject(dn) != "UTC" {
		t.Fatalf("getDisplayName expected UTC, got %q", object.GoStringFromStringObject(dn))
	}
	// raw/dst should be 0
	if got := tzGetRawOffset([]interface{}{obj}).(int64); got != 0 {
		t.Fatalf("rawOffset expected 0, got %d", got)
	}
	if got := ttzGetDSTSavings([]interface{}{obj}).(int64); got != 0 {
		t.Fatalf("dstSavings expected 0, got %d", got)
	}
}

func TestTimeZone_Setters_And_Offsets(t *testing.T) {
	globals.InitStringPool()
	obj := newTZObj(t)
	_ = tzInit([]interface{}{obj})

	// Change raw offset and verify getRawOffset and getOffset reflect it
	_ = tzSetRawOffset([]interface{}{obj, int64(3600 * 1000)}) // +1 hour
	if got := tzGetRawOffset([]interface{}{obj}).(int64); got != 3600*1000 {
		t.Fatalf("getRawOffset got %d want %d", got, 3600*1000)
	}
	if got := tzGetOffset([]interface{}{obj}).(int64); got != 3600*1000 {
		t.Fatalf("getOffset(IIIII) minimal impl should return raw, got %d", got)
	}
	// getOffset(long) for UTC id remains 0; for non-empty id we rely on tz logic
	_ = tzSetID([]interface{}{obj, str("UTC")})
	if got := tzGetOffsetLong([]interface{}{obj, int64(1_700_000_000_000)}).(int64); got != 0 {
		t.Fatalf("getOffset(long) for UTC expected 0, got %d", got)
	}
}

func TestTimeZone_GetTimeZone_Factory_And_Rules(t *testing.T) {
	globals.InitStringPool()
	// Static factory: pass String id as single param
	utcObj := tzGetTimeZoneString([]interface{}{str("UTC")} )
	if utcObj == nil {
		t.Fatalf("getTimeZone(UTC) returned nil")
	}
	utc := utcObj.(*object.Object)
	id := tzGetID([]interface{}{utc}).(*object.Object)
	if object.GoStringFromStringObject(id) != "UTC" {
		t.Fatalf("factory getID expected UTC, got %q", object.GoStringFromStringObject(id))
	}
	// Same rules with another UTC instance
	utc2 := tzGetTimeZoneString([]interface{}{str("UTC")} ).(*object.Object)
	assertJavaBoolTZ(t, tzHasSameRules([]interface{}{utc, utc2}), types.JavaBoolTrue, "same rules for two UTCs")
	// Different rules after rawOffset change
	_ = tzSetRawOffset([]interface{}{utc2, int64(1234)})
	assertJavaBoolTZ(t, tzHasSameRules([]interface{}{utc, utc2}), types.JavaBoolFalse, "different rules after raw change")
}

func TestTimeZone_AvailableIDs_And_InDaylight(t *testing.T) {
	globals.InitStringPool()
	arr := tzGetAvailableIDs(nil).(*object.Object)
	// Expect at least 2 entries (UTC,GMT)
	if arr == nil {
		t.Fatalf("getAvailableIDs returned nil")
	}
	// Object array is stored as [](*object.Object) in String[]; we can convert via object utilities
	// However, we only ensure it is a non-null array by relying on not panicking while reading value.
	// For a stronger check, we create strings and compare by scanning the Go string conversion of each element.
	// Build simple presence flags.
	vals, _ := arr.FieldTable["value"].Fvalue.([]*object.Object)
	foundUTC, foundGMT := false, false
	for _, e := range vals {
		if e == nil { continue }
		v := object.GoStringFromStringObject(e)
		if v == "UTC" { foundUTC = true }
		if v == "GMT" { foundGMT = true }
	}
	if !foundUTC || !foundGMT {
		t.Fatalf("available IDs should contain UTC and GMT; got UTC=%v GMT=%v", foundUTC, foundGMT)
	}
	// Filter by raw offset: 0 should return non-empty; non-zero returns empty in minimal impl
	empty := tzGetAvailableIDsInt([]interface{}{nil, int64(123)}).(*object.Object)
	vals2, _ := empty.FieldTable["value"].Fvalue.([]*object.Object)
	if len(vals2) != 0 {
		t.Fatalf("getAvailableIDs(123) expected empty, got %d entries", len(vals2))
	}

	// inDaylightTime should be false for UTC regardless of date
	obj := newTZObj(t)
	_ = tzInit([]interface{}{obj})
	_ = tzSetID([]interface{}{obj, str("UTC")} )
	// Construct a Date for some millis
	d := object.MakeEmptyObjectWithClassName(&[]string{"java/util/Date"}[0])
	_ = udateInitLong([]interface{}{d, int64(1_700_000_000_000)})
	assertJavaBoolTZ(t, tzInDaylightTime([]interface{}{obj, d}), types.JavaBoolFalse, "UTC not in DST")
}

func TestTimeZone_Clone_Minimal(t *testing.T) {
	globals.InitStringPool()
	obj := newTZObj(t)
	_ = tzInit([]interface{}{obj})
	_ = tzSetRawOffset([]interface{}{obj, int64(999)})
	cl := txClone([]interface{}{obj})
	clObj, ok := cl.(*object.Object)
	if !ok || clObj == nil {
		t.Fatalf("clone should return *object.Object, got %T", cl)
	}
	if clObj == obj {
		t.Fatalf("clone should return a distinct object pointer")
	}
	// Fields copied
	if got := tzGetRawOffset([]interface{}{clObj}).(int64); got != 999 {
		t.Fatalf("cloned rawOffset got %d want %d", got, 999)
	}
}
