package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/types"
	"testing"
)

func TestLoadUtilLoggingFilter_RegistersIsLoggable(t *testing.T) {
	Load_Util_Logging_Filter()

	key := "java/util/logging/Filter.isLoggable(Ljava/util/logging/LogRecord;)Z"
	gmeth, ok := ghelpers.MethodSignatures[key]
	if !ok {
		t.Fatalf("expected method signature %q to be registered", key)
	}
	if gmeth.ParamSlots != 1 {
		t.Fatalf("expected ParamSlots 1, got %d", gmeth.ParamSlots)
	}
}

func TestLoggingFilterIsLoggable_AlwaysTrue(t *testing.T) {
	if ret := loggingFilterIsLoggable(nil); ret != types.JavaBoolTrue {
		t.Fatalf("expected JavaBoolTrue for nil params, got %v", ret)
	}
	if ret := loggingFilterIsLoggable([]interface{}{"anything"}); ret != types.JavaBoolTrue {
		t.Fatalf("expected JavaBoolTrue regardless of params, got %v", ret)
	}
}
