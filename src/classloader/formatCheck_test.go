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
	"strings"
	"testing"
)

// Get an error if the klass.cpCount of entries does not match the actual number
func TestInvalidCPsize(t *testing.T) {
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

	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"Exceptions"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"testMethod"})

	klass.cpCount = 4 // the error we're testing. There are only two entries, not 4

	err := validateConstantPool(&klass)
	if err == nil {
		t.Error("Did not get error for mismatch between CP count field and actual number of CP entries")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr
	msg := string(out[:])

	if !strings.Contains(msg, "Error in size of constant pool") {
		t.Error("Did not get expected error msg for invalid CP count. Got: " + msg)
	}

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestInvalidIndexInUTF8Entry(t *testing.T) {
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
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 4}) // the error: there are only 2 UTF8 entries (see below)

	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"Exceptions"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"testMethod"})

	klass.cpCount = 2

	err := validateConstantPool(&klass)
	if err == nil {
		t.Error("Expected error for incorrect ut8Refs index, but got none.")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr
	msg := string(out[:])

	if !strings.Contains(msg, "points to invalid UTF8 entry") {
		t.Error("Did not get expected error msg. Got: " + msg)
	}

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestInvalidStringInUTF8Entry(t *testing.T) {
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

	invalidUtf8bytes := []byte{'B', 'a', 'd', 0xFA} // the last char is disallowed in UTF8 entries
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{string(invalidUtf8bytes)})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"testMethod"})

	klass.cpCount = 2

	err := validateConstantPool(&klass)
	if err == nil {
		t.Error("Expected error for invalid UTF8 string, but got none.")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr
	msg := string(out[:])

	if !strings.Contains(msg, "contains an invalid character") {
		t.Error("Did not get expected error msg. Got: " + msg)
	}

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestMissingDummyEntryAfterLongConst(t *testing.T) {
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
	klass.cpIndex = append(klass.cpIndex, cpEntry{LongConst, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 0}) // this should be a dummy entry

	klass.longConsts = append(klass.longConsts, int64(123))

	klass.cpCount = 3

	err := validateConstantPool(&klass)
	if err == nil {
		t.Error("Expected error for missing dummy entry after long, but got none.")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr
	msg := string(out[:])

	if !strings.Contains(msg, "Missing dummy entry") {
		t.Error("Did not get expected error msg. Got: " + msg)
	}

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestInvalidFieldRef(t *testing.T) {
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
	klass.cpIndex = append(klass.cpIndex, cpEntry{FieldRef, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 0}) // unimportant entry

	klass.fieldRefs = append(klass.fieldRefs, fieldRefEntry{
		classIndex:       1, // this points to a non-existent class ref
		nameAndTypeIndex: 0,
	})

	klass.cpCount = 3

	err := validateConstantPool(&klass)
	if err == nil {
		t.Error("Expected error for invalid class index in FieldRef entry, but got none.")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr
	msg := string(out[:])

	if !strings.Contains(msg, "points to an invalid entry in ClassRefs") {
		t.Error("Did not get expected error msg. Got: " + msg)
	}

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestFieldRefWithInvalidNameAndTypeIndex(t *testing.T) {
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
	klass.cpIndex = append(klass.cpIndex, cpEntry{FieldRef, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{ClassRef, 0})

	klass.fieldRefs = append(klass.fieldRefs, fieldRefEntry{
		classIndex:       2, // this correctly points to the ClassRef entry at klass.cpIndex[2]
		nameAndTypeIndex: 1, // this points to a non-existent class ref, causing the tested error
	})
	klass.classRefs = append(klass.classRefs, 0)

	klass.cpCount = 3

	err := validateConstantPool(&klass)
	if err == nil {
		t.Error("Expected error for invalid nameAndType index in FieldRef entry, but got none.")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr
	msg := string(out[:])

	if !strings.Contains(msg, "points to an invalid entry in nameAndType") {
		t.Error("Did not get expected error msg. Got: " + msg)
	}

	_ = wout.Close()
	os.Stdout = normalStdout
}

// a MethodRef points to a class index and a nameAndType index. The name in
// nameAndType must point to a valid class name. If that class name begins with
// a < then it must be <init>. This test makes sure of this latter part.
func TestMethodRefWithInvalidMethodName(t *testing.T) {
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
	klass.cpIndex = append(klass.cpIndex, cpEntry{MethodRef, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{ClassRef, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{NameAndType, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 0})

	klass.methodRefs = append(klass.methodRefs, methodRefEntry{
		classIndex:       2, // this correctly points to the ClassRef entry at klass.cpIndex[2]
		nameAndTypeIndex: 3, // this points to a nameAndType entry that points to an invalid class name
	})

	klass.classRefs = append(klass.classRefs, 3)

	klass.nameAndTypes = append(klass.nameAndTypes, nameAndTypeEntry{
		nameIndex:       4, // points to cpIndex[4], which is UTF8 rec w/ invalid name
		descriptorIndex: 0,
	})

	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"<invalidName>"})

	klass.cpCount = 5

	err := validateConstantPool(&klass)
	if err == nil {
		t.Error("Expected error for invalid method name in MethodRef's nameAndType entry, but got none.")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr
	msg := string(out[:])

	if !strings.Contains(msg, "an entry with an invalid method name") {
		t.Error("Did not get expected error msg. Got: " + msg)
	}

	_ = wout.Close()
	os.Stdout = normalStdout
}

// this test validates both InterfaceRefs and NameAndType refs.
func TestValidInterfaceRefEntry(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.CLASS)

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
	klass.cpIndex = append(klass.cpIndex, cpEntry{Interface, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{ClassRef, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{NameAndType, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 1})

	ire := interfaceRefEntry{classIndex: 2, nameAndTypeIndex: 3}
	klass.interfaceRefs = append(klass.interfaceRefs, ire)

	klass.classRefs = append(klass.classRefs, 4)

	klass.nameAndTypes = append(klass.nameAndTypes, nameAndTypeEntry{
		nameIndex:       4, // points to cpIndex[4], which is UTF8
		descriptorIndex: 5,
	})

	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"interface"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"B"})

	klass.cpCount = 6

	err := validateConstantPool(&klass)
	if err != nil {
		t.Error("Got but did not expect error in test of valid InterfaceRef.")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr
	msg := string(out[:])

	if len(msg) != 0 {
		t.Error("Got unexpected output to stderr: " + msg)
	}

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestInvalidFieldNameContainingWhitepace(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.CLASS)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// variables we'll need.
	klass := parsedClass{}
	klass.cpIndex = append(klass.cpIndex, cpEntry{})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 1})

	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"bad name"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"B"})

	klass.cpCount = 3

	klass.fieldCount = 1
	klass.fields = append(klass.fields, field{
		accessFlags: 0,
		name:        0, // points to the first utf8Refs entry
		description: 1, // points to the 2nd utf8Refs entry
		attributes:  nil,
	})

	err := validateFields(&klass)
	if err == nil {
		t.Error("Did not get expected error for invalid field name.")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	// out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr
	// msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout
}

// the field description must start with one only a few characters, of which
// 's' (our test value) is not one. We also test for an empty description
func TestInvalidFieldDescription(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.CLASS)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	// variables we'll need.
	klass := parsedClass{}
	klass.cpIndex = append(klass.cpIndex, cpEntry{})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 1})

	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"validName"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"s"})

	klass.cpCount = 3

	klass.fieldCount = 1
	klass.fields = append(klass.fields, field{
		accessFlags: 0,
		name:        0,
		description: 1,
		attributes:  nil,
	})

	err := validateFields(&klass)
	if err == nil {
		t.Error("Did not get expected error for invalid field description for " +
			"field: validName")
	}

	// now test for empty description string
	klass.utf8Refs[1] = utf8Entry{""}
	err = validateFields(&klass)
	if err == nil {
		t.Error("Did not get expected error for empty field description for " +
			"field: validName")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	// out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr
	// msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout
}
