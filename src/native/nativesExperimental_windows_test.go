/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package native

import (
	"jacobin/globals"
	"jacobin/log"
	"testing"
)

// simple test to exercise the code, rather than validating.
// validation tests will come in time.
func TestExports(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.FINE)
	err := CreateNativeFunctionTable("")
	if err != nil {
		t.Error(err)
	}
}
