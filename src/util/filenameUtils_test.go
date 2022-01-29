/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package util

import (
	"os"
	"strings"
	"testing"
)

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

func TestConverFilenameToPlatformPathSeparator(t *testing.T) {
	if os.PathSeparator == '\\' {
		s := ConvertToPlatformPathSeparators("snoop/dog/the/man")
		if strings.ContainsRune(s, '/') {
			t.Errorf("Expected a path with no / slashes, got: %s", s)
		}
		if strings.Count(s, string("\\")) != 3 {
			t.Error("Expected 3 backslashes in path, got: ", strings.Count(s, "\\"))
		}
	} else {
		var s string
		if os.PathSeparator == '/' {
			s = ConvertToPlatformPathSeparators("smoop\\dog\the\\man")
			if strings.ContainsRune(s, '\\') {
				t.Errorf("Expected a path with no \\ slashes, got: %s", s)
			}
		}
		if strings.Count(s, string("/")) != 3 {
			t.Error("Expected 3 forward slashes in path, got: ", strings.Count(s, "/"))
		}
	}
}
