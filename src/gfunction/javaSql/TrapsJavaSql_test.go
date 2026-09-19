package javaSql

import (
	"jacobin/src/gfunction/ghelpers"
	"reflect"
	"testing"
)

func TestLoad_Traps_Java_Sql_RegistersSomeMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Traps_Java_Sql()

	checks := []struct {
		key   string
		slots int
		fn    func([]interface{}) interface{}
	}{
		{"java/sql/Driver.<clinit>()V", 0, ghelpers.TrapClass},
		{"java/sql/DriverAction.<clinit>()V", 0, ghelpers.TrapClass},
		{"java/sql/DriverPropertyInfo.<clinit>()V", 0, ghelpers.TrapClass},
		{"java/sql/DriverPropertyInfo.<init>(Ljava/lang/String;Ljava/lang/String;)V", 2, ghelpers.TrapClass},
		{"java/sql/DriverManager.<clinit>()V", 0, ghelpers.TrapClass},
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
