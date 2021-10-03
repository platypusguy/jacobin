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

func TestValidExceptionsMethodAttribute(t *testing.T) {
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
	klass.cpIndex = append(klass.cpIndex, cpEntry{ClassRef, 0}) // points to classRef below, which points to the next CP entry
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 2})

	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"Exceptions"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"testMethod"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"java/io/IOException"})

	klass.classRefs = append(klass.classRefs, 3) // classRef[0] points to CP entry #4, which points to UTF #3

	klass.cpCount = 4

	// method
	meth := method{}
	meth.name = 2 // points to UTF8 entry: "testMethod"

	attrib := attr{}
	attrib.attrName = 1
	attrib.attrSize = 4
	attrib.attrContent = []byte{
		0, 1, // number of exceptions = 1
		0, 2, // points to 3rd CP entry, a classref that points to UTF8: java/io/IOException
	}

	err := parseExceptionsMethodAttribute(attrib, &meth, &klass)
	if err != nil {
		t.Error("Unexpected error in processing valid Exceptions attribute of method")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "java/io/IOException") {
		t.Error("Output did not contain name of exception. Got: " + msg)
	}

	if len(meth.exceptions) != 1 {
		t.Error("In test of Exceptions method attribute, attribute was not added to method struct")
	}

	me := meth.exceptions[0]
	if me != 2 {
		t.Error("The wrong value for the UTF8 record on Exceptions method attribute was stored. Got:" +
			strconv.Itoa(me))
	}
}

func TestValidMethodParametersAttribute(t *testing.T) {
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
	klass.cpIndex = append(klass.cpIndex, cpEntry{ClassRef, 0}) // points to classRef below, which points to the next CP entry
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 1})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 2})

	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"MethodParameters"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"param1"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"java/io/IOException"})

	klass.classRefs = append(klass.classRefs, 3) // classRef[0] points to CP entry #4, which points to UTF #3

	klass.cpCount = 5

	// method
	meth := method{}
	meth.name = 5 // points to UTF8 entry: "testMethod"

	attrib := attr{}
	attrib.attrName = 1 // CP[1] points to UTF8[0] -> "MethodParameters" (required)
	attrib.attrSize = 5 // 1 byte (param count) + 1 parameters of 2x2bytes = 5 bytes
	attrib.attrContent = []byte{
		0x01,       // just 1 attribute to process
		0x00, 0x03, // name index: CP[3] points to UTF8[1] -> name of parameter: "param1"
		0x80, 0x00, // access flags: ACC_MANDATED (a parameter from the language)
	}

	meth.accessFlags = 0x20

	err := parseMethodParametersAttribute(attrib, &meth, &klass)
	if err != nil {
		t.Error("Unexpected error in processing valid MethodParameter attribute of method")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "param1") {
		t.Error("Expected output containing 'param1' but got: " + msg)
	}

	if len(meth.parameters) != 1 {
		t.Error("In test of MethodParameters method attribute, attribute was not added to method struct")
	}

	mp := meth.parameters[0]
	if mp.name != "param1" {
		t.Error("The wrong value for the UTF8 record on MethodParams method attribute was stored. Got:" +
			mp.name)
	}
}
