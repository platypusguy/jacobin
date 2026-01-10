package ghelpers

import (
	"reflect"
	"testing"
)

func TestLoad_Traps_Java_Nio_RegistersSomeMethods(t *testing.T) {
	// Preserve global map and restore after test
	saved := MethodSignatures
	defer func() { MethodSignatures = saved }()
	MethodSignatures = make(map[string]GMeth)

	Load_Traps_Java_Nio()

	// Representative subset across java.nio/* and subpackages
	checks := []struct {
		key     string
		slots   int
		fn      func([]interface{}) interface{}
		checkFn bool
	}{
		{"java/nio/file/AccessMode.<clinit>()V", 0, TrapClass, true},
		{"java/nio/ByteBuffer.<clinit>()V", 0, TrapClass, true},
		{"java/nio/file/Files.<clinit>()V", 0, TrapClass, true},
		{"java/nio/charset/StandardCharsets.<clinit>()V", 0, TrapClass, true},
		{"java/nio/channels/FileChannel.<clinit>()V", 0, TrapClass, true},
	}

	for _, c := range checks {
		gm, ok := MethodSignatures[c.key]
		if !ok {
			t.Fatalf("missing MethodSignatures entry for %s", c.key)
		}
		if gm.ParamSlots != c.slots {
			t.Fatalf("%s ParamSlots expected %d, got %d", c.key, c.slots, gm.ParamSlots)
		}
		if gm.GFunction == nil {
			t.Fatalf("%s GFunction expected non-nil", c.key)
		}
		if c.checkFn {
			if reflect.ValueOf(gm.GFunction).Pointer() != reflect.ValueOf(c.fn).Pointer() {
				t.Fatalf("%s GFunction mismatch", c.key)
			}
		}
	}
}
