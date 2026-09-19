package sunSecurity

import (
	"jacobin/src/gfunction/ghelpers"
	"reflect"
	"testing"
)

func TestLoad_Traps_Sun_Security_RegistersSomeMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Traps_Sun_Security()

	checks := []struct {
		key   string
		slots int
		fn    func([]interface{}) interface{}
	}{
		{"sun/security/util/Debug.<clinit>()V", 0, ghelpers.ClinitGeneric},
		{"sun/security/util/Debug.getInstance(Ljava/lang/String;)Lsun/security/util/Debug;", 1, ghelpers.ReturnNull},
		{"sun/security/util/Debug.getInstance(Ljava/lang/String;Ljava/lang/String;)Lsun/security/util/Debug;", 2, ghelpers.ReturnNull},
		{"sun/security/util/Debug.isOn(Ljava/lang/String;)Z", 1, ghelpers.ReturnFalse},
	}

	for _, c := range checks {
		gm, ok := ghelpers.MethodSignatures[c.key]
		if !ok {
			t.Fatalf("missing MethodSignatures entry for %s", c.key)
		}
		if gm.ParamSlots != c.slots {
			t.Fatalf("%s ParamSlots expected %d, got %d", c.key, c.slots, gm.ParamSlots)
		}
		if gm.GFunction == nil {
			t.Fatalf("%s GFunction expected non-nil", c.key)
		}
		if reflect.ValueOf(gm.GFunction).Pointer() != reflect.ValueOf(c.fn).Pointer() {
			t.Fatalf("%s GFunction mismatch", c.key)
		}
	}
}
