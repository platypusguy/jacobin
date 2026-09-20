package javaNet

import (
	"jacobin/src/gfunction/ghelpers"
	"testing"
)

func TestLoad_Net_InetAddress_RegistersAllMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Net_InetAddress()

	if len(ghelpers.MethodSignatures) == 0 {
		t.Fatalf("expected method signatures to be registered")
	}

	if _, ok := ghelpers.MethodSignatures["java/net/InetAddress.<clinit>()V"]; !ok {
		t.Fatalf("expected <clinit> to be registered")
	}

	getByName, ok := ghelpers.MethodSignatures["java/net/InetAddress.getByName(Ljava/lang/String;)Ljava/net/InetAddress;"]
	if !ok {
		t.Fatalf("expected getByName() to be registered")
	}
	if getByName.ParamSlots != 1 {
		t.Errorf("expected getByName() ParamSlots 1, got %d", getByName.ParamSlots)
	}
}
