package javaNet

import (
	"jacobin/src/gfunction/ghelpers"
	"testing"
)

func TestLoad_Net_Socket_RegistersAllMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Net_Socket()

	if len(ghelpers.MethodSignatures) == 0 {
		t.Fatalf("expected method signatures to be registered")
	}

	clinit, ok := ghelpers.MethodSignatures["java/net/Socket.<clinit>()V"]
	if !ok {
		t.Fatalf("expected <clinit> to be registered")
	}
	if clinit.ParamSlots != 0 {
		t.Errorf("expected <clinit> ParamSlots 0, got %d", clinit.ParamSlots)
	}

	for key, gm := range ghelpers.MethodSignatures {
		if key == "java/net/Socket.<clinit>()V" {
			continue
		}
		if gm.GFunction == nil {
			t.Errorf("entry %q has nil GFunction", key)
		}
	}

	closeMeth, ok := ghelpers.MethodSignatures["java/net/Socket.close()V"]
	if !ok {
		t.Fatalf("expected close() to be registered")
	}
	if closeMeth.ParamSlots != 0 {
		t.Errorf("expected close() ParamSlots 0, got %d", closeMeth.ParamSlots)
	}
}
