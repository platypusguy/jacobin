/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by Andrew Binstock. All rights reserved.
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

// test insertion of klass into the method area (called Classes[])
func TestInsertValid(t *testing.T) {
	// Testing the changes made as a result of JACOBIN-103
	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.CLASS)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	Classes = make(map[string]Klass)
	currLen := len(Classes)
	k := Klass{
		Status: 0,
		Loader: "",
		Data:   &ClData{},
	}
	k.Data.Name = "testClass"
	k.Loader = "testLoader"
	k.Status = 'F'
	err := insert("TestEntry", k)
	if err != nil {
		t.Errorf("Got unexpected error on valid insertion into Classes[]: %s", err.Error())
	}

	newLen := len(Classes)
	if newLen != currLen+1 {
		t.Errorf("Expected post-insertion Classes[] to have length of %d, got: %d",
			currLen+1, newLen)
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "Class: testClass") {
		t.Errorf("Expecting log message containing 'Class: testClass', got: %s", msg)
	}
}

func TestFetchUTF8stringFromCPEntryNumber(t *testing.T) {
	cp := CPool{}

	cp.CpIndex = append(cp.CpIndex, CpEntry{})
	cp.CpIndex = append(cp.CpIndex, CpEntry{UTF8, 0})
	cp.CpIndex = append(cp.CpIndex, CpEntry{ClassRef, 0}) // points to classRef below, which points to the next CP entry
	cp.CpIndex = append(cp.CpIndex, CpEntry{UTF8, 2})

	cp.Utf8Refs = append(cp.Utf8Refs, "Exceptions")
	cp.Utf8Refs = append(cp.Utf8Refs, "testMethod")
	cp.Utf8Refs = append(cp.Utf8Refs, "java/io/IOException")

	s := FetchUTF8stringFromCPEntryNumber(&cp, 0) // invalid CP entry
	if s != "" {
		t.Error("Unexpected result in call toFetchUTF8stringFromCPEntryNumber()")
	}

	s = FetchUTF8stringFromCPEntryNumber(&cp, 1)
	if s != "Exceptions" {
		t.Error("Unexpected result in call toFetchUTF8stringFromCPEntryNumber()")
	}

	s = FetchUTF8stringFromCPEntryNumber(&cp, 2) // not UTF8, so should be an error
	if s != "" {
		t.Error("Unexpected result in call toFetchUTF8stringFromCPEntryNumber()")
	}
}
