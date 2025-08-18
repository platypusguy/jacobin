/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package execdata

import (
	"io"
	"jacobin/src/globals"
	"os"
	"strings"
	"testing"
)

func TestGetExecBuildInfo(t *testing.T) {
	globals.InitGlobals("test")
	g := globals.GetGlobalRef()
	GetExecBuildInfo(g)
	if g.JacobinBuildData == nil || len(g.JacobinBuildData) == 0 {
		t.Error("Expected JacobinBuildData to be populated")
	}
	t.Log(g.JacobinBuildData)
}

func TestPrintJacobinBuildData(t *testing.T) {
	globals.InitGlobals("test")
	g := globals.GetGlobalRef()
	g.JacobinBuildData = nil

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	PrintJacobinBuildData(g)

	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr
	msg := string(out[:])

	if !strings.Contains(msg, "Jacobin executable:") {
		t.Errorf("Expecting different output from PrintJacobinBuildData(), got: %s", msg)
	}
}
