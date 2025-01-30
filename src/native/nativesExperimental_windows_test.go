/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package native

import (
	"io"
	"jacobin/globals"
	"log"
	"strings"
	"testing"
)

// simple test to exercise the code, rather than validating.
// validation tests will come in time.
func TestExports(t *testing.T) {
	globals.InitGlobals("test")
	log.SetOutput(io.Discard) // turn off purego logging, which is noisy at times

	// expects to be running in the top Jacobin directory
	err := CreateNativeFunctionTable("testdata")
	if !strings.Contains(err.Error(), "not found") {
		t.Error("did not get expected error for missing file")
	}
}
