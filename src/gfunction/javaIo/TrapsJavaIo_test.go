package javaIo

import (
	"jacobin/src/gfunction/ghelpers"
	"reflect"
	"testing"
)

func TestLoad_Traps_Java_Io_RegistersSomeMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Traps_Java_Io()

	checks := []struct {
		key   string
		slots int
		// we won't check function pointer for all; for two representative entries we do
		fn      func([]interface{}) interface{}
		checkFn bool
	}{
		{"java/io/BufferedOutputStream.<clinit>()V", 0, ghelpers.TrapClass, true},
		{"java/io/DefaultFileSystem.getFileSystem()Ljava/io/FileSystem;", 0, ghelpers.TrapFunction, true},
		{"java/io/FileDescriptor.<clinit>()V", 0, ghelpers.TrapClass, true},
		{"java/io/FileDescriptor.valid()Z", 0, ghelpers.TrapFunction, true},
		{"java/io/FilterReader.<clinit>()V", 0, ghelpers.TrapClass, false},
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
		if c.checkFn {
			if reflect.ValueOf(gm.GFunction).Pointer() != reflect.ValueOf(c.fn).Pointer() {
				t.Fatalf("%s GFunction mismatch", c.key)
			}
		}
	}
}
