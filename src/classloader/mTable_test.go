/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"io"
	"jacobin/src/globals"
	"os"
	"strings"
	"testing"
)

// additional tests for loading native methods into an MTable
// are found in the gfunction package
func TestMtableAdd(t *testing.T) {
	mtbl := make(MT)
	AddEntry(&mtbl, "test1", MTentry{
		Meth:  nil,
		MType: 'G',
	})

	if len(mtbl) != 1 {
		t.Errorf("Expecting MTable size of 1, got: %d", len(mtbl))
	}

	if mtbl["test1"].MType != 'G' {
		t.Errorf("Expecting fetch of a 'G' MTable rec, but got type: %c",
			mtbl["test1"].MType)
	}
}

func TestMtableDump(t *testing.T) {
	globals.InitGlobals("test")
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	AddEntry(&MTable, "test1", MTentry{
		Meth:  nil,
		MType: 'G',
	})

	AddEntry(&MTable, "test0", MTentry{
		Meth:  nil,
		MType: 'J',
	})

	if len(MTable) != 2 {
		t.Errorf("Expecting MTable size of 2, got: %d", len(MTable))
	}

	DumpMTable()

	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr
	msg := string(out[:])
	if !strings.Contains(msg, "J   test0\nG") {
		t.Errorf("Expecting different content in dump of Mtable, got: %s", msg)
	}
}
