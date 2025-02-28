/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package config

import (
	"io"
	"os"
	"strings"
	"testing"
	"unicode"
)

func TestBuildNo(t *testing.T) {
	ver := BuildNo
	if ver < 2900 { // the first ever build # was the number of GH commits, which wa > 2900
		t.Errorf("BuildNo is too low: %d", ver)
	}
}

func TestVersionNo(t *testing.T) {
	ver := JacobinVersion

	// version # must begin and end with a digit
	verBytes := []rune(ver)
	if !unicode.IsDigit(verBytes[0]) || !unicode.IsDigit(verBytes[len(verBytes)-1]) {
		t.Errorf("Jacobin version using invalid vormat: %s", ver)
	}

	// there must be at least two dots in the version #
	count := 0
	for _, c := range ver {
		if c == '.' {
			count++
		}
	}
	if count < 2 {
		t.Errorf("Jacobin version using invalid vormat: %s", ver)
	}
}

// test dumping the config data to a file
func TestDumpConfig(t *testing.T) {
	file, err := os.CreateTemp("", "TestDumpConfig")
	if err != nil {
		t.Errorf("Error creating temporary file: %s", err)
	}

	defer os.Remove(file.Name())

	err = DumpConfig(file)

	if err != nil {
		errStr := err.Error()
		t.Errorf("%s", errStr)
	}

	_ = file.Close()
	file, err = os.Open(file.Name())
	output, err := io.ReadAll(file)

	if output == nil {
		t.Errorf("Expected output for DumpConfig(), but got none")
	}

	config := string(output)
	if !strings.Contains(config, "Version") || !strings.Contains(config, "OS") {
		t.Errorf("Got unexpected output for DumpConfig(), but got: %s", config)
	}
}
