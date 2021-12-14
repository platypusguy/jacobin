/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package util

import "testing"

func TestConvertInternalClassNameToFilename(t *testing.T) {
	// var s string
	if ConvertInternalClassNameToFilename("sponge") != "sponge.class" {
		t.Error("Unexpected result in call ConvertInternalClassNameToFilename()")
	}

	s := ConvertInternalClassNameToFilename("sponge/bob")
	if s != "sponge\\bob.class" {
		t.Error("Unexpected result in call ConvertInternalClassNameToFilename(): " + s)
	}

	s = ConvertInternalClassNameToFilename("sponge/bob\\square.pants")
	if s != "sponge\\bob\\square\\pants.class" {
		t.Error("Unexpected result in call ConvertInternalClassNameToFilename(): " + s)
	}
}

func TestConvertClassFilenameToInternalFormat(t *testing.T) {
	// var s string
	if ConvertClassFilenameToInternalFormat("sponge") != "sponge" {
		t.Error("Unexpected result in call ConvertClassFilenameToInternalFormat()")
	}

	s := ConvertClassFilenameToInternalFormat("sponge.bob.class")
	if s != "sponge/bob" {
		t.Error("Unexpected result in call ConvertInternalClassNameToFilename(): " + s)
	}

	s = ConvertClassFilenameToInternalFormat("sponge/bob/square.Pants.class")
	if s != "sponge/bob/square/Pants" {
		t.Error("Unexpected result in call ConvertInternalClassNameToFilename(): " + s)
	}
}
