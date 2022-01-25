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

// test a valid Code attribute of a method
func TestValidCodeMethodAttribute(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.FINEST)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// variables we'll need.
	klass := ParsedClass{}
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
		0, 4, // maxstack = 4
		0, 3, // maxlocals = 3
		0, 0, 0, 2, // code length = 2
		0x11, 0x16, // the two code bytes (randomly chosen)
		0, 0, // number of exceptions = 0 (exception handling is done elsewhere)
		0, 0, // attribute count of Code attribute (line number, etc.) = 0
	}

	err := parseCodeAttribute(attrib, &meth, &klass)
	if err != nil {
		t.Error("Unexpected error in processing valid Exceptions attribute of method")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	// out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	// msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if len(meth.codeAttr.code) != 2 {
		t.Error("Expected code length of 2. Got: " + strconv.Itoa(len(meth.codeAttr.code)))
	}

	if meth.codeAttr.maxStack != 4 {
		t.Error("Expected maxStack of 4. Got: " + strconv.Itoa(meth.codeAttr.maxStack))
	}

	if meth.codeAttr.maxLocals != 3 {
		t.Error("Expected maxLocals of 3. Got: " + strconv.Itoa(meth.codeAttr.maxLocals))
	}
	if len(meth.codeAttr.attributes) != 0 {
		t.Error("Expected 0 attributes of Code attribute. Got: " + strconv.Itoa(len(meth.codeAttr.attributes)))
	}
}
func Test1ValidMethodExceptionsAttribute(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.FINEST)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// variables we'll need.
	klass := ParsedClass{}
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

func Test2ValidMethodExceptionAttributes(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.FINEST)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// variables we'll need.
	klass := ParsedClass{}
	klass.cpIndex = append(klass.cpIndex, cpEntry{})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{ClassRef, 0}) // points to classRef below, which points to the next CP entry
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 1})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 2})
	klass.cpIndex = append(klass.cpIndex, cpEntry{ClassRef, 1})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 3})

	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"Exceptions"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"java/io/IOError"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"java/io/IOException"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"testMethod"})

	klass.classRefs = append(klass.classRefs, 4) // classRef[0] -> CP[4] -> UTF8[2]
	klass.classRefs = append(klass.classRefs, 3) // classRef[1] -> CP[3] -> UTF8[1]
	klass.cpCount = 7

	// method
	meth := method{}
	meth.name = 6 // CP[6] -> UTF8[3]: "testMethod"

	attrib := attr{}
	attrib.attrName = 1 // CP[1] points to UTF8[0] -> "MethodParameters" (required)
	attrib.attrSize = 6 // 2 (exc count) + 2x2bytes = 6 bytes
	attrib.attrContent = []byte{
		0x00, 0x02, // 2 exceptions to process
		0x00, 0x02, // Exception1: CP[2] points to ClassRef[0] -> CP[4]->UTF[3]: name of exception
		0x00, 0x05, // Exception2: CP[5] points to ClassRef[1] -> CP[3]->UTF[1]: name of exception
	}

	meth.accessFlags = 0x20

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
		t.Error("Expected output containing 'java/io/IOException' but got: " + msg)
	}

	if len(meth.exceptions) != 2 {
		t.Error("Expected 2 exceptions in Method 'testMethod', got: " +
			strconv.Itoa(len(meth.exceptions)))
	}

}

func TestValidMethodParameterAttribute(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.FINEST)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// variables we'll need.
	klass := ParsedClass{}
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

	if !validateUnqualifiedName(mp.name, false) {
		t.Error("MethodParameter name: " + mp.name + " is not a valid unqualified name")
	}
}
