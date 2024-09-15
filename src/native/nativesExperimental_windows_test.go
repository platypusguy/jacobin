/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package native

import (
	"io"
	"jacobin/globals"
	// "jacobin/log"
	"log"
	"testing"
)

// simple test to exercise the code, rather than validating.
// validation tests will come in time.
func TestExports(t *testing.T) {
	globals.InitGlobals("test")
	log.SetOutput(io.Discard) // turn off purego logging

	// os.Setenv("JAVA_HOME", ".\\testdata")
	// expects to be running in the top Jacobin directory
	err := CreateNativeFunctionTable("testdata")
	if err != nil {
		t.Error(err)
	}
}
