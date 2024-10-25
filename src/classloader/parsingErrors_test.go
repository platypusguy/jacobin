/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"io"
	"jacobin/globals"
	"jacobin/trace"
	"os"
	"strings"
	"testing"
)

// These are the bytecodes for a minimal class, Barebones.java, which simply prints a one-line greeting and exits.
var classBytes = []byte{
	/* 0000-0003: Magic number                  */ 0xCA, 0xFE, 0xBA, 0xBE,
	/* 0004-0005: Minor number of Java version  */ 0x00, 0x00,
	/* 0006-0007: Major number of Java version  */ 0x00, 0x3D, /* 0x3D = Java 17 */
	/* 0008-0009: # of entries in constant pool */ 0x00, 0x22,
	/* 0010-0392: the constant pool entries     */ 0x0A, 0x00, 0x02, 0x00, 0x03, 0x07,
	0x00, 0x04, 0x0C, 0x00, 0x05, 0x00, 0x06, 0x01, 0x00, 0x10, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x6C,
	0x61, 0x6E, 0x67, 0x2F, 0x4F, 0x62, 0x6A, 0x65, 0x63, 0x74, 0x01, 0x00, 0x06, 0x3C, 0x69, 0x6E,
	0x69, 0x74, 0x3E, 0x01, 0x00, 0x03, 0x28, 0x29, 0x56, 0x09, 0x00, 0x08, 0x00, 0x09, 0x07, 0x00,
	0x0A, 0x0C, 0x00, 0x0B, 0x00, 0x0C, 0x01, 0x00, 0x10, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x6C, 0x61,
	0x6E, 0x67, 0x2F, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6D, 0x01, 0x00, 0x03, 0x6F, 0x75, 0x74, 0x01,
	0x00, 0x15, 0x4C, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x69, 0x6F, 0x2F, 0x50, 0x72, 0x69, 0x6E, 0x74,
	0x53, 0x74, 0x72, 0x65, 0x61, 0x6D, 0x3B, 0x08, 0x00, 0x0E, 0x01, 0x00, 0x1A, 0x48, 0x65, 0x6C,
	0x6C, 0x6F, 0x20, 0x66, 0x72, 0x6F, 0x6D, 0x20, 0x62, 0x61, 0x72, 0x65, 0x62, 0x6F, 0x6E, 0x65,
	0x73, 0x20, 0x63, 0x6C, 0x61, 0x73, 0x73, 0x0A, 0x00, 0x10, 0x00, 0x11, 0x07, 0x00, 0x12, 0x0C,
	0x00, 0x13, 0x00, 0x14, 0x01, 0x00, 0x13, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x69, 0x6F, 0x2F, 0x50,
	0x72, 0x69, 0x6E, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6D, 0x01, 0x00, 0x07, 0x70, 0x72, 0x69,
	0x6E, 0x74, 0x6C, 0x6E, 0x01, 0x00, 0x15, 0x28, 0x4C, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x6C, 0x61,
	0x6E, 0x67, 0x2F, 0x53, 0x74, 0x72, 0x69, 0x6E, 0x67, 0x3B, 0x29, 0x56, 0x07, 0x00, 0x16, 0x01,
	0x00, 0x09, 0x42, 0x61, 0x72, 0x65, 0x62, 0x6F, 0x6E, 0x65, 0x73, 0x01, 0x00, 0x04, 0x43, 0x6F,
	0x64, 0x65, 0x01, 0x00, 0x0F, 0x4C, 0x69, 0x6E, 0x65, 0x4E, 0x75, 0x6D, 0x62, 0x65, 0x72, 0x54,
	0x61, 0x62, 0x6C, 0x65, 0x01, 0x00, 0x12, 0x4C, 0x6F, 0x63, 0x61, 0x6C, 0x56, 0x61, 0x72, 0x69,
	0x61, 0x62, 0x6C, 0x65, 0x54, 0x61, 0x62, 0x6C, 0x65, 0x01, 0x00, 0x04, 0x74, 0x68, 0x69, 0x73,
	0x01, 0x00, 0x0B, 0x4C, 0x42, 0x61, 0x72, 0x65, 0x62, 0x6F, 0x6E, 0x65, 0x73, 0x3B, 0x01, 0x00,
	0x04, 0x6D, 0x61, 0x69, 0x6E, 0x01, 0x00, 0x16, 0x28, 0x5B, 0x4C, 0x6A, 0x61, 0x76, 0x61, 0x2F,
	0x6C, 0x61, 0x6E, 0x67, 0x2F, 0x53, 0x74, 0x72, 0x69, 0x6E, 0x67, 0x3B, 0x29, 0x56, 0x01, 0x00,
	0x04, 0x61, 0x72, 0x67, 0x73, 0x01, 0x00, 0x13, 0x5B, 0x4C, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x6C,
	0x61, 0x6E, 0x67, 0x2F, 0x53, 0x74, 0x72, 0x69, 0x6E, 0x67, 0x3B, 0x01, 0x00, 0x0A, 0x53, 0x6F,
	0x75, 0x72, 0x63, 0x65, 0x46, 0x69, 0x6C, 0x65, 0x01, 0x00, 0x0E, 0x42, 0x61, 0x72, 0x65, 0x62,
	0x6F, 0x6E, 0x65, 0x73, 0x2E, 0x6A, 0x61, 0x76, 0x61,

	/* 0393-0394: Access flags                             */ 0x00, 0x21,
	/* 0395-0396: Pointer to CP record for class name      */ 0x00, 0x15,
	/* 0397-0398: Pointer to CP record for superclass name */ 0x00, 0x02,
	/* 0399-0400: Count of interfaces  (should = 0)        */ 0x00, 0x00,
	/* 0401-0402: Count of fields (should = 0)             */ 0x00, 0x00,
	/* 0403-0404: Count of methods                         */ 0x00, 0x02,
	/* 0405-0534: Contents of methods                      */
	0x00, 0x01, 0x00, 0x05, 0x00, 0x06, 0x00, 0x01, 0x00, 0x17, 0x00,
	0x00, 0x00, 0x2F, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x05, 0x2A, 0xB7, 0x00, 0x01, 0xB1,
	0x00, 0x00, 0x00, 0x02, 0x00, 0x18, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x00, 0x00, 0x04,
	0x00, 0x19, 0x00, 0x00, 0x00, 0x0C, 0x00, 0x01, 0x00, 0x00, 0x00, 0x05, 0x00, 0x1A, 0x00, 0x1B,
	0x00, 0x00, 0x00, 0x09, 0x00, 0x1C, 0x00, 0x1D, 0x00, 0x01, 0x00, 0x17, 0x00, 0x00, 0x00, 0x37,
	0x00, 0x02, 0x00, 0x01, 0x00, 0x00, 0x00, 0x09, 0xB2, 0x00, 0x07, 0x12, 0x0D, 0xB6, 0x00, 0x0F,
	0xB1, 0x00, 0x00, 0x00, 0x02, 0x00, 0x18, 0x00, 0x00, 0x00, 0x0A, 0x00, 0x02, 0x00, 0x00, 0x00,
	0x06, 0x00, 0x08, 0x00, 0x07, 0x00, 0x19, 0x00, 0x00, 0x00, 0x0C, 0x00, 0x01, 0x00, 0x00, 0x00,
	0x09, 0x00, 0x1E, 0x00, 0x1F, 0x00, 0x00,
	/* 0535-0536: Class attribute count                   */ 0x00, 0x01,
	/* 0537-0544: Class attributes                        */
	0x00, 0x20, 0x00, 0x00, 0x00, 0x02, 0x00, 0x21,
}

func TestLoadBarebonesClass(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	_, err := parse(classBytes)
	if err != nil {
		t.Errorf("Got unexpected error from parse of Class.class: %s", err.Error())
	}
}

func TestBarebonesClassWithInvalidMagicNumber(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	rout, wout, _ := os.Pipe()
	os.Stdout = wout

	classBytes[1] = 0xFF // this byte should be 0xFE, so this should generate an error.
	_, err := parse(classBytes)
	if err == nil {
		t.Errorf("Expected an error but did not get one: %s", err.Error())
	}

	classBytes[1] = 0xFE // reset the erroneous byte to its correct value

	_ = w.Close()
	msgOut, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	_ = wout.Close()
	_, _ = io.ReadAll(rout)
	os.Stdout = normalStdout

	if !strings.Contains(string(msgOut), "Class Format Error: invalid magic number") {
		t.Errorf("Expected msg: 'Class Format Error: invalid magic number', got: %s", string(msgOut))
	}
}

func TestBarebonesClassWithInvalidJavaVersion(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	rout, wout, _ := os.Pipe()
	os.Stdout = wout

	classBytes[7] = 0xFF // this byte should be 0x3D for Java 17
	_, err := parse(classBytes)
	if err == nil {
		t.Errorf("Expected an error but did not get one: %s", err.Error())
	}

	classBytes[7] = 0x3D // reset the erroneous byte to its correct value

	_ = w.Close()
	msgOut, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	_ = wout.Close()
	_, _ = io.ReadAll(rout)
	os.Stdout = normalStdout

	if !strings.Contains(string(msgOut), "Jacobin supports only Java versions through") {
		t.Errorf("Expected msg to contain: 'Jacobin supports only Java versions through', \n got: %s", string(msgOut))
	}
}

func TestBarebonesClassWithInvalidCPcount(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	rout, wout, _ := os.Pipe()
	os.Stdout = wout

	classBytes[9] = 0x00 // this byte should be 0x22 (0 CP entries is always an error in a class file)
	_, err := parse(classBytes)
	if err == nil {
		t.Errorf("Expected an error but did not get one: %s", err.Error())
	}

	classBytes[9] = 0x22 // reset the erroneous byte to its correct value

	_ = w.Close()
	msgOut, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	_ = wout.Close()
	_, _ = io.ReadAll(rout)
	os.Stdout = normalStdout

	if !strings.Contains(string(msgOut), "Invalid number of entries in constant pool") {
		t.Errorf("Expected msg to contain: 'Invalid number of entries in constant pool', \n got: %s", string(msgOut))
	}
}
