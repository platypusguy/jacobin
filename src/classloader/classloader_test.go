/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import "testing"

// Most of the functionality in classloader package is tested in other files, such as
// * cpParser_test.go (constant pool parser)
// * formatCheck_test.go (the format checking)
// * parser_test.go (the class parsing)
// etc.
// This files tests remaining routines.

// remove leading [L and delete trailing;, eliminate all other entries with [prefix
func TestNormalizingClassReference(t *testing.T) {
	s := normalizeClassReference("[Ljava/test/java.String;")
	if s != "java/test/java.String" {
		t.Error("Unexpected normalized class reference: " + s)
	}

	s = normalizeClassReference("[B")
	if s != "" {
		t.Error("Unexpected normalized class reference: " + s)
	}

	s = normalizeClassReference("java/lang/Object")
	if s != "java/lang/Object" {
		t.Error("Unexpected normalized class reference: " + s)
	}
}
