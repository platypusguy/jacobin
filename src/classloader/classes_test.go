/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */
package classloader

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
