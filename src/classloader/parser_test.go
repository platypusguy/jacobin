/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"io/ioutil"
	"jacobin/globals"
	"jacobin/log"
	"os"
	"strings"
	"testing"
)

func TestParseOfInvalidJavaVersionNumber(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()

	// redirect stderr to inspect output
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	bytesToTest := []byte{0xCA, 0xFE, 0xBA, 0xBE, 0x00, 0x00, 0xFF, 0xF0}
	err := parseJavaVersionNumber(bytesToTest, parsedClass{})

	// restore stderr to what it was before
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	if err == nil {
		t.Error("Invalid Java version number did not generate an error")
	}

	if !strings.Contains(msg, "supports only Java versions") {
		t.Error("Did not get expected error msg for invalid Java version. Got: " + msg)
	}
}

func TestParseValidJavaVersion(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	bytesToTest := []byte{0xCA, 0xFE, 0xBA, 0xBE, 0x00, 0x00, 0x00, 0x30}
	err := parseJavaVersionNumber(bytesToTest, parsedClass{})
	if err != nil {
		t.Error("valid Java version # generated an error in version # parser")
	}
}
