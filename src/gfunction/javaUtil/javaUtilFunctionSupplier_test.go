package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"reflect"
	"testing"
)

func TestLoadUtilFunctionSupplier_RegistersGetMethod(t *testing.T) {
	Load_Util_Function_Supplier()

	key := "java/util/function/Supplier.get()Ljava/lang/Object;"
	gmeth, ok := ghelpers.MethodSignatures[key]
	if !ok {
		t.Fatalf("expected method signature %q to be registered", key)
	}
	if gmeth.ParamSlots != 0 {
		t.Fatalf("expected ParamSlots 0, got %d", gmeth.ParamSlots)
	}
	if reflect.ValueOf(gmeth.GFunction).Pointer() != reflect.ValueOf(ghelpers.TrapFunction).Pointer() {
		t.Fatalf("expected GFunction to be ghelpers.TrapFunction")
	}
}
