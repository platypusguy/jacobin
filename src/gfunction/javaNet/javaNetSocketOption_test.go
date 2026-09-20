package javaNet

import (
	"jacobin/src/gfunction/ghelpers"
	"testing"
)

func TestLoad_Net_SocketOption_RegistersAllMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Net_SocketOption()

	if _, ok := ghelpers.MethodSignatures["java/net/SocketOption.<clinit>()V"]; !ok {
		t.Fatalf("expected <clinit> to be registered")
	}
	if _, ok := ghelpers.MethodSignatures["java/net/SocketOption.name()Ljava/lang/String;"]; !ok {
		t.Fatalf("expected name() to be registered")
	}
	if _, ok := ghelpers.MethodSignatures["java/net/SocketOption.type()Ljava/lang/Class;"]; !ok {
		t.Fatalf("expected type() to be registered")
	}
}
