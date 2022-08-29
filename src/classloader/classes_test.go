/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */
package classloader

import (
	"io"
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
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "Class: testClass") {
		t.Errorf("Expecting log message containing 'Class: testClass', got: %s", msg)
	}
}

func TestInvalidLookupOfMethod_Test0(t *testing.T) {
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
	k.Loader = ""
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

	_, err = FetchMethodAndCP("TestEntry", "main", "([L)V")
	if err == nil {
		t.Errorf("Expecting an err msg for invalid Fetch in MTable, but got none")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "Class: testClass") {
		t.Errorf("Expecting log message containing 'Class: testClass', got: %s", msg)
	}
}

func TestInvalidLookupOfMethod_Test1(t *testing.T) {
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

	k.Loader = "testloader"
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

	_, err = FetchMethodAndCP("TestEntry", "main", "([L)V")
	if err == nil {
		t.Errorf("Expecting an err msg for invalid Fetch of main() in MTable, but got none")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "Main method not found in class") {
		t.Errorf("Expecting log message containing 'Main method not found in class', got: %s", msg)
	}
}

func TestInvalidLookupOfMethod_Test2(t *testing.T) {
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

	k.Loader = "testloader"
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

	// fetch a non-existent class, called 'gherkin'
	_, err = FetchMethodAndCP("TestEntry", "gherkin", "([L)V")
	if err == nil {
		t.Errorf("Expecting an err msg for invalid Fetch of main() in MTable, but got none")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "it did not contain method: gherkin") {
		t.Errorf("Expecting log message containing 'it did not contain method: gherkin', got: %s", msg)
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
