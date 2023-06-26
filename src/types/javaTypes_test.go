/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package types

import "testing"

func TestJavaBoolean(t *testing.T) {

	val := ConvertGoBoolToJavaBool(true)
	if val != JavaBoolTrue {
		t.Errorf("JavaBool: expected a result of 1, but got: %d", val)
	}

	val = ConvertGoBoolToJavaBool(false)
	if val != JavaBoolFalse {
		t.Errorf("JavaBool: expected a result of 0, but got: %d", val)
	}
}
