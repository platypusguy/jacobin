package javaNet

import (
	"jacobin/src/gfunction/ghelpers"
	"testing"
)

func TestLoad_Net_InetSocketAddress_RegistersAllMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Net_InetSocketAddress()

	if _, ok := ghelpers.MethodSignatures["java/net/InetSocketAddress.<clinit>()V"]; !ok {
		t.Fatalf("expected <clinit> to be registered")
	}
	if _, ok := ghelpers.MethodSignatures["java/net/InetSocketAddress.<init>(I)V"]; !ok {
		t.Fatalf("expected <init>(I)V to be registered")
	}
	if _, ok := ghelpers.MethodSignatures["java/net/InetSocketAddress.getPort()I"]; !ok {
		t.Fatalf("expected getPort()I to be registered")
	}
}
