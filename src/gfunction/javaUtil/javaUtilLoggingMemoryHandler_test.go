package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"reflect"
	"testing"
)

func TestLoadUtilLoggingMemoryHandler_RegistersAllMethods(t *testing.T) {
	Load_Util_Logging_MemoryHandler()

	expected := map[string]int{
		"java/util/logging/MemoryHandler.close()V":                                      0,
		"java/util/logging/MemoryHandler.flush()V":                                      0,
		"java/util/logging/MemoryHandler.getPushLevel()Ljava/util/logging/Level;":       0,
		"java/util/logging/MemoryHandler.publish(Ljava/util/logging/LogRecord;)V":        1,
		"java/util/logging/MemoryHandler.push()V":                                       0,
		"java/util/logging/MemoryHandler.setPushLevel(Ljava/util/logging/Level;)V":       1,
	}

	for key, paramSlots := range expected {
		gmeth, ok := ghelpers.MethodSignatures[key]
		if !ok {
			t.Fatalf("expected method signature %q to be registered", key)
		}
		if gmeth.ParamSlots != paramSlots {
			t.Fatalf("key %q: expected ParamSlots %d, got %d", key, paramSlots, gmeth.ParamSlots)
		}
		if reflect.ValueOf(gmeth.GFunction).Pointer() != reflect.ValueOf(ghelpers.TrapFunction).Pointer() {
			t.Fatalf("key %q: expected GFunction to be ghelpers.TrapFunction", key)
		}
	}
}
