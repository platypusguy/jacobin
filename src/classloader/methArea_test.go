/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
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

// Note: many MethArea functions are tested in classes_test,go
// These tests simply fill in untested functions, rather than duplicate those tests

func TestMethAreadDelete(t *testing.T) {
	MethArea = &sync.Map{}
	methAreaSize = 0
	currLen := MethAreaSize()
	if currLen != 0 {
		t.Errorf("Expecting MethArea size of 0, got: %d", currLen)
	}

	k := Klass{
		Status: 0,
		Loader: "",
		Data:   &ClData{},
	}
	k.Data.Name = "testClass1"
	k.Data.Superclass = "java/lang/Object"
	k.Loader = "testloader"
	k.Status = 'F'
	MethAreaInsert("TestEntry1", &k)
	MethAreaInsert("TestEntry2", &k)
	currLen = MethAreaSize()
	if MethAreaSize() != 2 {
		t.Errorf("Expecting MethArea size of 2, got: %d", currLen)
	}

	MethAreaDelete("TestEntry1")
	newLen := MethAreaSize()
	if newLen != 1 {
		t.Errorf("Expected post-deletion MethArea[] to have length of 1, got: %d",
			newLen)
	}
}

func TestMethAreadDeleteNonExistentEntry(t *testing.T) {
	MethArea = &sync.Map{}
	methAreaSize = 0
	currLen := MethAreaSize()
	if currLen != 0 {
		t.Errorf("Expecting MethArea size of 0, got: %d", currLen)
	}

	k := Klass{
		Status: 0,
		Loader: "",
		Data:   &ClData{},
	}
	k.Data.Name = "testClass1"
	k.Data.Superclass = "java/lang/Object"
	k.Loader = "testloader"
	k.Status = 'F'
	MethAreaInsert("TestEntry", &k)
	currLen = MethAreaSize()
	if MethAreaSize() != 1 {
		t.Errorf("Expecting MethArea size of 1, got: %d", currLen)
	}

	// deleting a non-entry should not cause an error or reduce MethArea size
	MethAreaDelete("NoSuchEntry")
	newLen := MethAreaSize()
	if newLen != 1 {
		t.Errorf("Expected post-deletion MethArea[] to have length of 1, got: %d",
			newLen)
	}
}

func TestMethAreadFetchNonExistentEntry(t *testing.T) {
	// Testing the changes made as a result of JACOBIN-103
	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.CLASS)

	// redirect stderr to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	MethArea = &sync.Map{}
	methAreaSize = 0
	currLen := MethAreaSize()
	if currLen != 0 {
		t.Errorf("Expecting MethArea size of 0, got: %d", currLen)
	}

	k := Klass{
		Status: 0,
		Loader: "",
		Data:   &ClData{},
	}
	k.Data.Name = "testClass1"
	k.Data.Superclass = "java/lang/Object"
	k.Loader = "testloader"
	k.Status = 'F'
	MethAreaInsert("TestEntry", &k)

	// fetching a non-entry should not cause an error, shiuld return nil
	me := MethAreaFetch("NoSuchEntry")
	if me != nil {
		t.Errorf("Expected nil return from MethAreaFetch(), got: %v", me)
	}

	// restore stderr
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])
	if !strings.Contains(msg, "--> nil") {
		t.Errorf("Expected different log message, got: %s", msg)
	}
}
