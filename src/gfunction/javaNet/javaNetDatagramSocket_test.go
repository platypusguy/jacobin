package javaNet

import (
	"jacobin/src/gfunction/ghelpers"
	"testing"
)

func TestLoad_Net_DatagramSocket_RegistersAllMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Net_DatagramSocket()

	if _, ok := ghelpers.MethodSignatures["java/net/DatagramSocket.<clinit>()V"]; !ok {
		t.Fatalf("expected <clinit> to be registered")
	}
	if _, ok := ghelpers.MethodSignatures["java/net/DatagramSocket.<init>()V"]; !ok {
		t.Fatalf("expected <init>()V to be registered")
	}
	if _, ok := ghelpers.MethodSignatures["java/net/DatagramSocket.send(Ljava/net/DatagramPacket;)V"]; !ok {
		t.Fatalf("expected send(DatagramPacket) to be registered")
	}
}
