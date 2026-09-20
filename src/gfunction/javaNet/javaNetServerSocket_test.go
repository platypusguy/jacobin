package javaNet

import (
	"jacobin/src/gfunction/ghelpers"
	"testing"
)

func TestLoad_Net_ServerSocket_RegistersAllMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Net_ServerSocket()

	if len(ghelpers.MethodSignatures) == 0 {
		t.Fatalf("expected method signatures to be registered")
	}

	if _, ok := ghelpers.MethodSignatures["java/net/ServerSocket.<clinit>()V"]; !ok {
		t.Fatalf("expected <clinit> to be registered")
	}

	accept, ok := ghelpers.MethodSignatures["java/net/ServerSocket.accept()Ljava/net/Socket;"]
	if !ok {
		t.Fatalf("expected accept() to be registered")
	}
	if accept.GFunction == nil {
		t.Errorf("accept() GFunction should not be nil")
	}
}
