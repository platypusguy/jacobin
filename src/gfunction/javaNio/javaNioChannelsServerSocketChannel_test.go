package javaNio

import (
	"jacobin/src/gfunction/ghelpers"
	"testing"
)

func TestLoad_Nio_Channels_ServerSocketChannel_RegistersAllMethods(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Nio_Channels_ServerSocketChannel()

	if len(ghelpers.MethodSignatures) == 0 {
		t.Fatalf("expected method signatures to be registered")
	}

	if _, ok := ghelpers.MethodSignatures["java/nio/channels/ServerSocketChannel.<clinit>()V"]; !ok {
		t.Fatalf("expected <clinit> to be registered")
	}

	open, ok := ghelpers.MethodSignatures["java/nio/channels/ServerSocketChannel.open()Ljava/nio/channels/ServerSocketChannel;"]
	if !ok {
		t.Fatalf("expected open() to be registered")
	}
	if open.GFunction == nil {
		t.Errorf("open() GFunction should not be nil")
	}
}
