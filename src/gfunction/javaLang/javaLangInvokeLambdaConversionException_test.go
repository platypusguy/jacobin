/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"jacobin/src/gfunction/ghelpers"
	"testing"
)

func TestLoad_Java_Lang_Invoke_LambdaConversionException_Registration(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Java_Lang_Invoke_LambdaConversionException()

	expectedConstructors := []string{
		"java/lang/invoke/LambdaConversionException.<init>()V",
		"java/lang/invoke/LambdaConversionException.<init>(Ljava/lang/String;)V",
		"java/lang/invoke/LambdaConversionException.<init>(Ljava/lang/String;Ljava/lang/Throwable;)V",
		"java/lang/invoke/LambdaConversionException.<init>(Ljava/lang/Throwable;)V",
	}

	for _, sig := range expectedConstructors {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Expected constructor signature to be registered: %s", sig)
		}
	}
}
