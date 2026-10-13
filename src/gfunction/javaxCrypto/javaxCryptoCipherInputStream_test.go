/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaxCrypto

import (
	"jacobin/src/gfunction/ghelpers"
	"reflect"
	"runtime"
	"testing"
)

func TestCipherInputStream_MethodSignatures(t *testing.T) {
	Load_Crypto_CipherInputStream()

	expectedMethods := []string{
		"javax/crypto/CipherInputStream.<clinit>()V",
		"javax/crypto/CipherInputStream.<init>(Ljava/io/InputStream;)V",
		"javax/crypto/CipherInputStream.<init>(Ljava/io/InputStream;Ljavax/crypto/Cipher;)V",
		"javax/crypto/CipherInputStream.available()I",
		"javax/crypto/CipherInputStream.close()V",
		"javax/crypto/CipherInputStream.mark(I)V",
		"javax/crypto/CipherInputStream.markSupported()Z",
		"javax/crypto/CipherInputStream.read()I",
		"javax/crypto/CipherInputStream.read([B)I",
		"javax/crypto/CipherInputStream.read([BII)I",
		"javax/crypto/CipherInputStream.reset()V",
		"javax/crypto/CipherInputStream.skip(J)J",
	}

	for _, m := range expectedMethods {
		if _, ok := ghelpers.MethodSignatures[m]; !ok {
			t.Errorf("expected method signature %q to be registered", m)
		}
	}
}

func TestCipherInputStream_AllAreTraps(t *testing.T) {
	Load_Crypto_CipherInputStream()

	clinitName := runtime.FuncForPC(reflect.ValueOf(ghelpers.ClinitGeneric).Pointer()).Name()
	trapName := runtime.FuncForPC(reflect.ValueOf(ghelpers.TrapFunction).Pointer()).Name()

	prefix := "javax/crypto/CipherInputStream."
	for key, meth := range ghelpers.MethodSignatures {
		if len(key) < len(prefix) || key[:len(prefix)] != prefix {
			continue
		}
		gotName := runtime.FuncForPC(reflect.ValueOf(meth.GFunction).Pointer()).Name()
		if key == "javax/crypto/CipherInputStream.<clinit>()V" {
			if gotName != clinitName {
				t.Errorf("expected %q to use ClinitGeneric, got %s", key, gotName)
			}
			continue
		}
		if gotName != trapName {
			t.Errorf("expected %q to be a TrapFunction, got %s", key, gotName)
		}
	}
}
