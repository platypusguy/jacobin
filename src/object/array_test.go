/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package object

import "testing"

// This file tests array primitives. Array bytecodes are tested in
// jvm.arrayByetcodes.go

func TestArrayTypeConversions(t *testing.T) {
	if JdkArrayTypeToJacobinType(T_BOOLEAN) != BYTE {
		t.Errorf("did not get expected Jacobin type BOOLEAN")
	}

	if JdkArrayTypeToJacobinType(T_LONG) != INT {
		t.Errorf("did not get expected Jacobin type for LONG")
	}

	if JdkArrayTypeToJacobinType(T_DOUBLE) != FLOAT {
		t.Errorf("did not get expected Jacobin type for DOUBLE")
	}

	if JdkArrayTypeToJacobinType(T_REF) != REF {
		t.Errorf("did not get expected Jacobin type for REF")
	}

	if JdkArrayTypeToJacobinType(99) != 0 {
		t.Errorf("did not get expected Jacobin type for invalid value")
	}
}
