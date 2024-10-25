/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */
package classloader

import (
	"jacobin/globals"
	"jacobin/types"
	"os"
	"strings"
	"sync"
	"testing"
)

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
		t.Error("Unexpected result in call to FetchUTF8stringFromCPEntryNumber()")
	}

	s = FetchUTF8stringFromCPEntryNumber(&cp, 1)
	if s != "Exceptions" {
		t.Error("Unexpected result in call to FetchUTF8stringFromCPEntryNumber()")
	}

	s = FetchUTF8stringFromCPEntryNumber(&cp, 2) // not UTF8, so should be an error
	if s != "" {
		t.Error("Unexpected result in call to FetchUTF8stringFromCPEntryNumber()")
	}

	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestInvalidMainMethod(t *testing.T) {
	// Testing the changes made as a result of JACOBIN-103
	globals.InitGlobals("test")

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	MethArea = &sync.Map{}
	k := Klass{
		Status: 0,
		Loader: "",
		Data:   &ClData{},
	}
	k.Data.Name = "testClass"
	k.Data.SuperclassIndex = types.ObjectPoolStringIndex
	k.Loader = "testloader"
	k.Status = 'F'
	MethAreaInsert("TestEntry", &k)

	// we need a java/lang/Object instance, so just duplicate the entry
	// in the MethArea. It's only a placeholder
	MethAreaInsert(types.ObjectClassName, &k)

	// fetch a non-existent main() method
	_, err := FetchMethodAndCP("java/lan/Object", "main", "([LString;)V")
	if err == nil {
		t.Errorf("Expecting an err msg for invalid MethAreaFetch of main(), but got none")
	}

	msg := err.Error()
	if !strings.Contains(msg, "main() method not found") {
		t.Errorf("TestInvalidLookupOfMethod: Expecting error of 'main() method not found', got %s", err.Error())
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr
}

func TestInvalidClassName(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	MethArea = &sync.Map{}

	// fetch a non-existent method in a non-existent class
	_, err := FetchMethodAndCP("gherkin", "mcMurtry", "([LString;)V")
	if err == nil {
		t.Errorf("Expecting an err msg for invalid MethAreaFetch of main(), but got none")
	}

	msg := err.Error()
	if !strings.HasPrefix(msg, "FetchMethodAndCP: LoadClassFromNameOnly for gherkin failed") {
		t.Errorf("TestInvalidClassName: Did not get expected error message', got %s", err.Error())
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr
}
