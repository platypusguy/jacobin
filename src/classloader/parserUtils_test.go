/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"jacobin/globals"
	"jacobin/log"
	"os"
	"strconv"
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
