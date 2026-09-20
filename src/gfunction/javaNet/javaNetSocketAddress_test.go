package javaNet

import (
	"jacobin/src/gfunction/ghelpers"
	"testing"
)

func TestLoad_Net_SocketAddress_RegistersAllMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Net_SocketAddress()

	if _, ok := ghelpers.MethodSignatures["java/net/SocketAddress.<clinit>()V"]; !ok {
		t.Fatalf("expected <clinit> to be registered")
	}
	if _, ok := ghelpers.MethodSignatures["java/net/SocketAddress.<init>()V"]; !ok {
		t.Fatalf("expected <init> to be registered")
	}
}
