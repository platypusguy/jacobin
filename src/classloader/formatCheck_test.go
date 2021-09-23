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

// Get an error if the klass.cpCount of entries does not match the actual number
func TestInvalidCPsize(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.FINEST)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// variables we'll need.
	klass := parsedClass{}
	klass.cpIndex = append(klass.cpIndex, cpEntry{})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 0})

	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"Exceptions"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"testMethod"})

	klass.cpCount = 4 // the error we're testing. There are only two entries, not 4

	err := validateConstantPool(&klass)
	if err == nil {
		t.Error("Did not get error for mismatch between CP count field and actual number of CP entries")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr
	msg := string(out[:])

	if !strings.Contains(msg, "Error in size of constant pool") {
		t.Error("Did not get expected error msg for invalid CP count. Got: " + msg)
	}

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestInvalidIndexInUTF8Entry(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.FINEST)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// variables we'll need.
	klass := parsedClass{}
	klass.cpIndex = append(klass.cpIndex, cpEntry{})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 4}) // the error: there are only 2 UTF8 entries (see below)

	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"Exceptions"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"testMethod"})

	klass.cpCount = 2

	err := validateConstantPool(&klass)
	if err == nil {
		t.Error("Expected error for incorrect ut8Refs index, but got none.")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr
	msg := string(out[:])

	if !strings.Contains(msg, "points to invalid UTF8 entry") {
		t.Error("Did not get expected error msg. Got: " + msg)
	}

	_ = wout.Close()
	os.Stdout = normalStdout
}
