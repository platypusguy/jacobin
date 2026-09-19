package jdkInternal

import (
	"jacobin/src/gfunction/ghelpers"
	"reflect"
	"testing"
)

func TestLoad_Traps_Jdk_Internal_RegistersSomeMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Traps_Jdk_Internal()

	checks := []struct {
		key   string
		slots int
		fn    func([]interface{}) interface{}
	}{
		{"jdk/internal/access/SharedSecrets.<clinit>()V", 0, ghelpers.ClinitGeneric},
		{"jdk/internal/misc/VM.initialize()V", 0, ghelpers.JustReturn},
		{"jdk/internal/misc/CDS.getRandomSeedForDumping()J", 0, ghelpers.ReturnRandomLong},
		{"jdk/internal/misc/CDS.initializeFromArchive(Ljava/lang/Class;)V", 1, ghelpers.JustReturn},
		{"jdk/internal/misc/CDS.isDumpingArchive0()Z", 0, ghelpers.ReturnFalse},
		{"jdk/internal/misc/CDS.isDumpingClassList0()Z", 0, ghelpers.ReturnFalse},
		{"jdk/internal/misc/CDS.isSharingEnabled0()Z", 0, ghelpers.ReturnFalse},
		{"jdk/internal/util/ArraysSupport.<clinit>()V", 0, ghelpers.ClinitGeneric},
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
