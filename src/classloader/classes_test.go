/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */
package classloader

import (
	"io"
	"jacobin/globals"
	"jacobin/log"
	"os"
	"strings"
	"sync"
	"testing"
)

// test insertion of klass into the method area (called MethArea[])
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

	MethArea = &sync.Map{}
	currLen := MethAreaSize()
	k := Klass{
		Status: 0,
		Loader: "",
		Data:   &ClData{},
	}
	k.Data.Name = "testClass"
	k.Loader = "testLoader"
	k.Status = 'F'
	MethAreaInsert("TestEntry", &k)

	newLen := MethAreaSize()
	if newLen != currLen+1 {
		t.Errorf("Expected post-insertion MethArea[] to have length of %d, got: %d",
			currLen+1, newLen)
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "Method area insert: testClass, loader: testLoader") {
		t.Errorf("Expecting log message containing 'testClass', got: %s", msg)
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

	MethArea = &sync.Map{}
	currLen := MethAreaSize()
	k := Klass{
		Status: 0,
		Loader: "",
		Data:   &ClData{},
	}
	k.Data.Name = "testClass"
	k.Loader = ""
	k.Status = 'F'
	MethAreaInsert("TestEntry", &k)

	newLen := MethAreaSize()
	if newLen != currLen+1 {
		t.Errorf("Expected post-insertion MethArea[] to have length of %d, got: %d",
			currLen+1, newLen)
	}

	_, err := FetchMethodAndCP("TestEntry", "main", "([L)V")
	if err == nil {
		t.Errorf("Expecting an err msg for invalid MethAreaFetch in MTable, but got none")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "Method area insert: testClass, loader:") {
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

	MethArea = &sync.Map{}
	currLen := MethAreaSize()
	k := Klass{
		Status: 0,
		Loader: "",
		Data:   &ClData{},
	}
	k.Data.Name = "testClass"
	k.Data.Superclass = "java/lang/Object"
	k.Loader = "testloader"
	k.Status = 'F'
	MethAreaInsert("TestEntry", &k)

	// we need a java/lang/Object instance, so just duplicate the entry
	// in the MethArea. It's only a placeholder
	MethAreaInsert("java/lang/Object", &k)

	newLen := MethAreaSize()
	if newLen != currLen+2 {
		t.Errorf("Expected post-insertion MethArea[] to have length of %d, got: %d",
			currLen+1, newLen)
	}

	_, err := FetchMethodAndCP("TestEntry", "main", "([L)V")
	if err == nil {
		t.Errorf("Expecting an err msg for invalid MethAreaFetch of main() in MTable, but got none")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "main() method not found in class") {
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

	MethArea = &sync.Map{}
	currLen := MethAreaSize()
	k := Klass{
		Status: 0,
		Loader: "",
		Data:   &ClData{},
	}
	k.Data.Name = "testClass"
	k.Data.Superclass = "java/lang/Object"
	k.Loader = "testloader"
	k.Status = 'F'
	MethAreaInsert("TestEntry", &k)

	// we need a java/lang/Object instance, so just duplicate the entry
	// in the MethArea. It's only a placeholder
	MethAreaInsert("java/lang/Object", &k)

	newLen := MethAreaSize()
	if newLen != currLen+2 {
		t.Errorf("Expected post-insertion MethArea[] to have length of %d, got: %d",
			currLen+1, newLen)
	}

	// fetch a non-existent class, called 'gherkin'
	_, err := FetchMethodAndCP("TestEntry", "gherkin", "([L)V")
	if err == nil {
		t.Errorf("Expecting an err msg for invalid MethAreaFetch of main() in MTable, but got none")
	}

	if err.Error() != "method not found" {
		t.Errorf("Expecting error of 'method not found', got %s", err.Error())
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "did not contain method: gherkin") {
		t.Errorf("Expecting log message containing 'did not contain method: gherkin', got: %s", msg)
	}
}

func TestFetchUTF8stringFromCPEntryNumber(t *testing.T) {
	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

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

	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}
