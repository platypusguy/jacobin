package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"reflect"
	"testing"
)

func TestLoad_Traps_Java_Util_RegistersSomeMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Traps_Java_Util()

	checks := []struct {
		key   string
		slots int
		fn    func([]interface{}) interface{}
	}{
		{"java/util/concurrent/Executors.<clinit>()V", 0, ghelpers.TrapClass},
		{"java/util/concurrent/Executors.newCachedThreadPool()Ljava/util/concurrent/ExecutorService;", 0, ghelpers.TrapFunction},
		{"java/util/concurrent/Executors.newCachedThreadPool(Ljava/util/concurrent/ThreadFactory;)Ljava/util/concurrent/ExecutorService;", 1, ghelpers.TrapFunction},
		{"java/util/SimpleTimeZone.<clinit>()V", 0, ghelpers.TrapClass},
		{"java/util/random/RandomGenerator.<clinit>()V", 0, ghelpers.TrapClass},
		{"java/util/random/RandomGenerator.getDefault()Ljava/util/random/RandomGenerator;", 0, ghelpers.TrapFunction},
		{"java/util/random/RandomGenerator.of(Ljava/lang/String;)Ljava/util/random/RandomGenerator;", 1, ghelpers.TrapFunction},
		{"java/util/random/RandomGeneratorFactory.<clinit>()V", 0, ghelpers.TrapClass},
		{"java/util/random/RandomGeneratorFactory.of(Ljava/lang/String;)Ljava/util/random/RandomGeneratorFactory;", 1, ghelpers.TrapFunction},
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
