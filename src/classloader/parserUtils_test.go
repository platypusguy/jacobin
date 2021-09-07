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
	"strconv"
	"strings"
	"testing"
)

func TestGetIntFrom2BytesInvalid(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	bytesToTest := []byte{0xCA, 0xFE}
	_, err := intFrom2Bytes(bytesToTest, 3)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout

	if err == nil {
		t.Error("intFrom2Bytes() did not return an error when given an invalid offset")
	}
}

func TestGetIntFrom2BytesValid(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	bytesToTest := []byte{0x01, 0x0B}
	i, err := intFrom2Bytes(bytesToTest, 0)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout

	if i != 267 || err != nil {
		t.Error("intFrom2Bytes() should have returned 267, but got: " + strconv.Itoa(i))
	}
}

func TestGetIntFrom4BytesValid(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	bytesToTest := []byte{0x01, 0x02, 0x03, 0x04}
	i, err := intFrom4Bytes(bytesToTest, 0)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout

	if i != 16909060 || err != nil {
		t.Error("intFrom4Bytes() should have returned 16909060, but got: " + strconv.Itoa(i))
	}
}

func TestFetchValidUTF8string(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	klass := parsedClass{}
	klass.cpIndex = append(klass.cpIndex, cpEntry{})
	klass.cpIndex = append(klass.cpIndex, cpEntry{1, 0}) // the UTF-8 reference
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"gherkin"})
	klass.cpCount = 2

	result, err := fetchUTF8string(&klass, 1)
	if err != nil {
		t.Error("Unexpected error testing fetch of UTF8 entry")
	}

	if result != "gherkin" {
		t.Error("Expecting fetch of UTF8 to return 'gherkin' but got: " + result)
	}
	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestFetchInvalidUTF8string(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	// redirect stderr & stdout to capture results from stderr and to
	// prevent error message from showing up in the test results
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	klass := parsedClass{}
	klass.cpIndex = append(klass.cpIndex, cpEntry{})
	klass.cpIndex = append(klass.cpIndex, cpEntry{7, 0}) // the invalid UTF-8 reference
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"gherkin"})
	klass.cpCount = 2

	_, err := fetchUTF8string(&klass, 1)
	if err == nil {
		t.Error("Expected error testing fetch of invalid UTF8 entry, but got none")
	}

	// restore stderr and stdout to what they were before
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "attempt to fetch UTF8 string from non-UTF8 CP entry") {
		t.Error("Expected different error msg on failed fetch of UTF-8 CP entry. Got: " + msg)
	}
}
