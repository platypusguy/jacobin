/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */
package exec

import (
	"testing"
)

func TestFetchUTF8stringFromCPEntryNumber(t *testing.T) {
	cp := CPool{}

	cp.CpIndex = append(cp.CpIndex, CpEntry{})
	cp.CpIndex = append(cp.CpIndex, CpEntry{UTF8, 0})
	cp.CpIndex = append(cp.CpIndex, CpEntry{ClassRef, 0}) // points to classRef below, which points to the next CP entry
	cp.CpIndex = append(cp.CpIndex, CpEntry{UTF8, 2})

	cp.Utf8Refs = append(cp.Utf8Refs, "Exceptions")
	cp.Utf8Refs = append(cp.Utf8Refs, "testMethod")
	cp.Utf8Refs = append(cp.Utf8Refs, "java/io/IOException")

	s := FetchUTF8stringFromCPEntryNumber(&cp, 0) // invalid CP entry
	if s != "" {
		t.Error("Unexpected result in call toFetchUTF8stringFromCPEntryNumber()")
	}

	s = FetchUTF8stringFromCPEntryNumber(&cp, 1)
	if s != "Exceptions" {
		t.Error("Unexpected result in call toFetchUTF8stringFromCPEntryNumber()")
	}

	s = FetchUTF8stringFromCPEntryNumber(&cp, 2) // not UTF8, so should be an error
	if s != "" {
		t.Error("Unexpected result in call toFetchUTF8stringFromCPEntryNumber()")
	}
}

func TestConvertInternalClassNameToFilename(t *testing.T) {
	// var s string
	if ConvertInternalClassNameToFilename("sponge") != "sponge.class" {
		t.Error("Unexpected result in call ConvertInternalClassNameToFilename()")
	}

	s := ConvertInternalClassNameToFilename("sponge/bob")
	if s != "sponge.bob.class" {
		t.Error("Unexpected result in call ConvertInternalClassNameToFilename(): " + s)
	}

	s = ConvertInternalClassNameToFilename("sponge/bob\\square.pants")
	if s != "sponge.bob.square.pants.class" {
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
