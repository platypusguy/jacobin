package gfunction

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

// Helpers
func newSDFObj() *object.Object {
	className := "java/text/SimpleDateFormat"
	return object.MakeEmptyObjectWithClassName(&className)
}

func sdfStr(s string) *object.Object { return object.StringObjectFromGoString(s) }

func TestSimpleDateFormat_MethodRegistration(t *testing.T) {
	// Ensure the string pool is initialized for class names and strings.
	globals.InitStringPool()

	// Clear and load just the SDF methods to MethodSignatures map for this test context.
	MethodSignatures = make(map[string]GMeth)
	Load_Math_SimpleDateFormat()

	// Verify a few critical registrations and their ParamSlots.
	tests := []struct {
		key   string
		slots int
	}{
		{"java/text/SimpleDateFormat.<init>()V", 0},
		{"java/text/SimpleDateFormat.<init>(Ljava/lang/String;)V", 1},
		{"java/text/SimpleDateFormat.clone()Ljava/lang/Object;", 0},
		{"java/text/SimpleDateFormat.toPattern()Ljava/lang/String;", 0},
		{"java/text/SimpleDateFormat.applyPattern(Ljava/lang/String;)V", 1},
		{"java/text/SimpleDateFormat.applyLocalizedPattern(Ljava/lang/String;)V", 1},
	}

	for _, tc := range tests {
		gm, ok := MethodSignatures[tc.key]
		if !ok {
			t.Fatalf("Method signature not registered: %s", tc.key)
		}
		if gm.ParamSlots != tc.slots {
			t.Fatalf("ParamSlots mismatch for %s: want %d, got %d", tc.key, tc.slots, gm.ParamSlots)
		}
		if gm.GFunction == nil {
			t.Fatalf("GFunction pointer is nil for %s", tc.key)
		}
	}
}

func TestSimpleDateFormat_Init_NoPattern(t *testing.T) {
	globals.InitStringPool()

	obj := newSDFObj()
	// Pre-populate a dummy field to ensure constructor clears/initializes FieldTable anew.
	obj.FieldTable["dummy"] = object.Field{Ftype: types.Int, Fvalue: int64(1)}

	ret := sdfInit([]interface{}{obj})
	if ret != nil {
		t.Fatalf("sdfInit returned non-nil: %v", ret)
	}
	if obj.FieldTable == nil {
		t.Fatalf("FieldTable is nil after sdfInit")
	}
	if _, exists := obj.FieldTable["dummy"]; exists {
		t.Fatalf("FieldTable was not reinitialized; unexpected 'dummy' key remains")
	}
}

func TestSimpleDateFormat_Init_WithPattern(t *testing.T) {
	globals.InitStringPool()

	obj := newSDFObj()
 pat := sdfStr("yyyy-MM-dd")

	ret := sdfInitString([]interface{}{obj, pat})
	if ret != nil {
		t.Fatalf("sdfInitString returned non-nil: %v", ret)
	}
	fld, ok := obj.FieldTable["pattern"]
	if !ok {
		t.Fatalf("pattern field not set by sdfInitString")
	}
	if fld.Ftype != types.StringClassRef {
		t.Fatalf("pattern field Ftype mismatch: want %s, got %s", types.StringClassRef, fld.Ftype)
	}
	if fld.Fvalue != pat {
		t.Fatalf("pattern field Fvalue should be the same String object reference")
	}
}

func TestSimpleDateFormat_Init_WithNullPattern(t *testing.T) {
	globals.InitStringPool()

	obj := newSDFObj()
	// Pass a nil *object.Object as Java null for the pattern.
	var nullStr *object.Object = nil

	ret := sdfInitString([]interface{}{obj, nullStr})
	if ret != nil {
		t.Fatalf("sdfInitString(null) returned non-nil: %v", ret)
	}
	fld, ok := obj.FieldTable["pattern"]
	if !ok {
		t.Fatalf("pattern field not set by sdfInitString with null pattern")
	}
	if fld.Ftype != types.StringClassRef {
		t.Fatalf("pattern field Ftype mismatch with null: want %s, got %s", types.StringClassRef, fld.Ftype)
	}
 if !object.IsNull(fld.Fvalue) {
		t.Fatalf("pattern field Fvalue should be Java null for null pattern; got %T", fld.Fvalue)
	}
}

func TestSimpleDateFormat_Clone_Minimal(t *testing.T) {
	globals.InitStringPool()

	obj := newSDFObj()
	_ = sdfInit([]interface{}{obj})

	cl := sdfClone([]interface{}{obj})
	clObj, ok := cl.(*object.Object)
	if !ok || clObj == nil {
						t.Fatalf("sdfClone should return a non-nil *object.Object, got %T", cl)
	}
	if clObj != obj {
		t.Fatalf("sdfClone minimal behavior should return same object reference")
	}
}

func TestSimpleDateFormat_ToPattern_Behavior(t *testing.T) {
	globals.InitStringPool()

	// 1) With explicit pattern
	obj1 := newSDFObj()
	pat1 := sdfStr("yyyy-MM-dd")
	_ = sdfInitString([]interface{}{obj1, pat1})
	out1 := sdfToPattern([]interface{}{obj1})
	if out1 == nil {
		t.Fatalf("toPattern returned nil for explicit pattern")
	}
	so1, ok := out1.(*object.Object)
	if !ok || so1 == nil {
		t.Fatalf("toPattern did not return a String object: %T", out1)
	}
	if gs := object.GoStringFromStringObject(so1); gs != "yyyy-MM-dd" {
		t.Fatalf("toPattern mismatch: got %q want %q", gs, "yyyy-MM-dd")
	}

	// 2) With null pattern -> expect Java null
	obj2 := newSDFObj()
	var nullStr *object.Object = nil
	_ = sdfInitString([]interface{}{obj2, nullStr})
	out2 := sdfToPattern([]interface{}{obj2})
	if !object.IsNull(out2) {
		t.Fatalf("toPattern with null pattern should return null; got %T", out2)
	}

	// 3) No pattern (default constructor) -> expect empty string
	obj3 := newSDFObj()
	_ = sdfInit([]interface{}{obj3})
	out3 := sdfToPattern([]interface{}{obj3})
	so3, ok := out3.(*object.Object)
	if !ok || so3 == nil {
		t.Fatalf("toPattern with no pattern should return a String object: %T", out3)
	}
	if gs := object.GoStringFromStringObject(so3); gs != "" {
		t.Fatalf("toPattern with no pattern: got %q want empty string", gs)
	}
}


func TestSimpleDateFormat_ApplyPattern_Behavior(t *testing.T) {
	globals.InitStringPool()

	// 1) Start with default ctor, then applyPattern("MM/dd/yyyy")
	obj := newSDFObj()
	_ = sdfInit([]interface{}{obj})
	_ = sdfApplyPattern([]interface{}{obj, sdfStr("MM/dd/yyyy")})
	out := sdfToPattern([]interface{}{obj})
	so, ok := out.(*object.Object)
	if !ok || so == nil {
		t.Fatalf("toPattern after applyPattern should return a String object: %T", out)
	}
	if gs := object.GoStringFromStringObject(so); gs != "MM/dd/yyyy" {
		t.Fatalf("applyPattern did not set pattern: got %q want %q", gs, "MM/dd/yyyy")
	}

	// 2) applyPattern(null) should set pattern to null
	_ = sdfApplyPattern([]interface{}{obj, (*object.Object)(nil)})
	out2 := sdfToPattern([]interface{}{obj})
	if !object.IsNull(out2) {
		t.Fatalf("applyPattern(null) should result in toPattern==null; got %T", out2)
	}

	// 3) applyLocalizedPattern minimal impl behaves the same (loader maps to same Go func)
	_ = sdfApplyPattern([]interface{}{obj, sdfStr("yyyyMMdd")}) // set to some value first
	_ = sdfApplyPattern([]interface{}{obj, sdfStr("yyyy-MM")})
	out3 := sdfToPattern([]interface{}{obj})
	so3, ok := out3.(*object.Object)
	if !ok || so3 == nil {
		t.Fatalf("toPattern after applyLocalizedPattern should return a String object: %T", out3)
	}
	if gs := object.GoStringFromStringObject(so3); gs != "yyyy-MM" {
		t.Fatalf("applyLocalizedPattern behavior mismatch: got %q want %q", gs, "yyyy-MM")
	}
}
