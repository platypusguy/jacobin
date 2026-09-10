package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"reflect"
	"testing"
)

func TestLoadUtilLoggingSocketHandler_RegistersAllMethods(t *testing.T) {
	Load_Util_Logging_SocketHandler()

	expected := map[string]int{
		"java/util/logging/SocketHandler.close()V":                                0,
		"java/util/logging/SocketHandler.publish(Ljava/util/logging/LogRecord;)V": 1,
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
