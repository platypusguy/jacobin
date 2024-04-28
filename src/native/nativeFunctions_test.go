/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package native

import "testing"

func TestLoadUnsupportedNativeMethods(t *testing.T) {
	size := LoadUnsupportedNativeMethods()
	if size < 1 {
		t.Errorf("LoadUnsupportedNativeMethods size should be > 0, got %d", size)
	}
}

func TestIsUnsupportedNativeMethod(t *testing.T) {
	LoadUnsupportedNativeMethods()
	if !IsUnsupportedNativeMethod("test.entry") {
		t.Errorf("IsUnsupportedNativeMethod(\"test.entry\") should be true")
	}

	if IsUnsupportedNativeMethod("no/such.entry") {
		t.Errorf("IsUnsupportedNativeMethod(\"no/such.entry\") should be false")
	}
}
