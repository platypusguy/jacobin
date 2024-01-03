/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package object

import "testing"

func TestIsNull(t *testing.T) {
	if !IsNull(nil) {
		t.Errorf("nil should be null")
	}

	var op *Object
	if !IsNull(op) {
		t.Errorf("pointer to non-allocated object should be null")
	}
}
