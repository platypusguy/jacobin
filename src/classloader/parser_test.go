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

// Magic number should be OxCAFEBABE in the first four bytes of the classfile
func TestMagicNumber(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()

	// redirect stderr to inspect output
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	bytesToTest := []byte{0xCA, 0xFE, 0xBA, 0xBA, 0x00, 0x00, 0xFF, 0xF0}
	err := parseMagicNumber(bytesToTest)

	// restore stderr to what it was before
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	if err == nil {
		t.Error("Invalid Java magic number did not generate an error")
	}

	if !strings.Contains(msg, "invalid magic number") {
		t.Error("Did not get expected error msg for invalid magic number. Got: " + msg)
	}
}

func TestParseOfInvalidJavaVersionNumber(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()

	// redirect stderr to inspect output
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	bytesToTest := []byte{0xCA, 0xFE, 0xBA, 0xBE, 0x00, 0x00, 0xFF, 0xF0}
	err := parseJavaVersionNumber(bytesToTest, &parsedClass{})

	// restore stderr to what it was before
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	if err == nil {
		t.Error("Invalid Java version number did not generate an error")
	}

	if !strings.Contains(msg, "supports only Java versions") {
		t.Error("Did not get expected error msg for invalid Java version. Got: " + msg)
	}
}

func TestParseValidJavaVersion(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	bytesToTest := []byte{0xCA, 0xFE, 0xBA, 0xBE, 0x00, 0x00, 0x00, 0x30}
	err := parseJavaVersionNumber(bytesToTest, &parsedClass{})
	if err != nil {
		t.Error("valid Java version # generated an error in version # parser")
	}
}

func TestConstantPoolCountValid(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	pClass := parsedClass{}

	bytesToTest := []byte{0xCA, 0xFE, 0xBA, 0xBE, 0x00, 0x00, 0x00, 0x30, 0x00, 0x20}
	err := getConstantPoolCount(bytesToTest, &pClass)
	if err != nil {
		t.Error("valid constant pool count generated an error in version # parser")
	}

	if pClass.cpCount != 32 {
		t.Error("expected a pool count of 32, instead got: " +
			strconv.Itoa(pClass.cpCount))
	}
}

func TestConstantPoolCountInvalid(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	// redirect stderr to inspect output
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	bytesToTest := []byte{0xCA, 0xFE, 0xBA, 0xBE, 0x00, 0x00, 0x00, 0x30, 0x00, 0x01}
	err := getConstantPoolCount(bytesToTest, &parsedClass{})

	// restore stderr to what it was before
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	if err == nil {
		t.Error("Invalid constant pool entry count did not generate an error")
	}

	if !strings.Contains(msg, "Invalid number of entries in constant pool") {
		t.Error("Did not get expected error msg for invalid number of entries in CP. Got: " + msg)
	}
}

// Access flags consist of a 2-byte integer. In the parsing, a variety of booleans are set in
// the parsed class to show what access is allowed by the access flags. Both the retrieval of
// the value and setting of the booleans is tested here.
func TestAccessFlags(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	pc := parsedClass{}
	bytes := []byte{0x00, 0x84, 0x21}
	loc, err := parseAccessFlags(bytes, 0, &pc)

	if err != nil {
		t.Error("Unexpected error occurred testing parse of Access flags")
	}

	if loc != 2 {
		t.Error("Expected location from parse of Access flags to be 2. Got: " + strconv.Itoa(loc))
	}

	if pc.classIsPublic == false ||
		pc.classIsSuper == false ||
		pc.classIsAbstract == false ||
		pc.classIsModule == false {
		t.Error("Access flags did not set expected values in the parsed class")
	}
}

func TestClassNameInvalidIndex(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	pc := parsedClass{}
	bytes := []byte{0x00, 0x00, 0x10}
	_, err := parseClassName(bytes, 0, &pc)

	if err == nil {
		t.Error("Should have returned an error for invalid value in class name item")
	}
}

// a complex test. It first parses a minimal constant pool that has the records we need
// for the actual test. It then passes bytes containing the class name entry and tests
// whether all the records and pointers in the CP structs point to the right entry.
func TestClassNameValidName(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	pc := parsedClass{}
	pc.cpCount = 3
	bytes := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00, // the required first 10 bytes
		0x00, 0x00, 0x37, 0x00, 0x03, // Java 8, CP with 3 entries (plus the dummy entry)
		0x07, 0x00, 0x02, // entry #1, a ClassRef that points to the following UTF-8 record
		0x01, 0x00, 0x05, 'H', 'e', 'l', 'l', 'o', // entry #2, the UTF-8 record containing "Hello"
	}

	_, err := parseConstantPool(bytes, &pc)
	if err != nil {
		t.Error("Error parsing test CP for setup in testing ClassName")
	}

	testBytes := []byte{0x00, 0x00, 0x01} // 3 bytes b/c first byte is skipped. So, this points to entry 1
	_, err = parseClassName(testBytes, 0, &pc)
	if err != nil {
		t.Error("Unexpected error in getting class name from the CP")
	}

	if pc.className != "Hello" {
		t.Error("Test of getting class name should get 'Hello' but got: " + pc.className)
	}
}

// see notes about the setup in the previous test. Here, the record that is pointed to by the
// class name field is not a ClassRef, but instead a string constant entry. This should generate
// an error, for which we test.
func TestClassNameWhenDoesNotPointToClassRef(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	pc := parsedClass{}
	pc.cpCount = 3
	bytes := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00, // the required first 10 bytes
		0x00, 0x00, 0x37, 0x00, 0x03, // Java 8, CP with 3 entries (plus the dummy entry)
		0x08, 0x00, 0x02, // entry #1, should be a ClassRef entry, but is not
		0x01, 0x00, 0x05, 'H', 'e', 'l', 'l', 'o', // entry #2, the UTF-8 record containing "Hello"
	}

	_, err := parseConstantPool(bytes, &pc)
	if err != nil {
		t.Error("Error parsing test CP for setup in testing ClassName")
	}

	testBytes := []byte{0x00, 0x00, 0x01} // 3 bytes b/c first byte is skipped. So, this points to entry 1
	_, err = parseClassName(testBytes, 0, &pc)
	if err == nil {
		t.Error("Parse of class name field should have generated an error but it did not.")
	}
	if !strings.HasPrefix(err.Error(), "Class Format Error: invalid entry for class name") {
		t.Error("Expected error msg about invalid entry for class name. Got: " + err.Error())
	}
}

// see the previous tests for explanation of the setup. Here we test whether a class name entry
// that points to a valid ClassRef record, but when that ClassRef record does not itself point
// to an expected UTF-8 entry, that the right error is issued.
func TestClassNameWithMissingUTF8(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	pc := parsedClass{}
	pc.cpCount = 3
	bytes := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00, // the required first 10 bytes
		0x00, 0x00, 0x37, 0x00, 0x03, // Java 8, CP with 3 entries (plus the dummy entry)
		0x07, 0x00, 0x02, // entry #1, a ClassRef that should point to a UTF-8 entry
		0x07, 0x00, 0x01, // entry #2, this should be a UTF-8 entry, but it's not
	}

	_, err := parseConstantPool(bytes, &pc)
	if err != nil {
		t.Error("Error parsing test CP for setup in testing ClassName")
	}

	testBytes := []byte{0x00, 0x00, 0x01} // 3 bytes b/c first byte is skipped. So, this points to entry 1
	_, err = parseClassName(testBytes, 0, &pc)
	if err == nil {
		t.Error("Parse of class name field should have generated an error but it did not.")
	}

	if !strings.HasPrefix(err.Error(), "Class Format Error: error classRef in CP does not point to a UTF-8 string") {
		t.Error("Expected error msg about invalid UTF-8 entry for class name. Got: " + err.Error())
	}
}

func TestErrorOnEmptySuperclassName(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	pc := parsedClass{}
	pc.cpCount = 5
	bytes := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00, // the required first 10 bytes
		0x00, 0x00, 0x37, 0x00, 0x05, // Java 8, CP with 5 entries (plus the dummy entry)
		// entry #0, a dummy entry created by the JVM
		0x07, 0x00, 0x02, // entry #1, a ClassRef that points to the following UTF-8 record
		0x01, 0x00, 0x05, 'H', 'e', 'l', 'l', 'o', // entry #2, the UTF-8 record containing "Hello"
		0x07, 0x00, 0x04, // entry #3, a ClassRef that points to the following UTF-8 record
		0x01, 0x00, 0x00, // emtry #4 an empty string
	}

	_, err := parseConstantPool(bytes, &pc)
	if err != nil {
		t.Error("Error parsing test CP for setup in testing superclassName")
	}

	testBytes := []byte{0x00, 0x00, 0x01, // 3 bytes b/c first byte is skipped. So, this points to entry 1
		0x00, 0x03, // points to the superclass entry (entry #3)
	}

	_, err = parseSuperClassName(testBytes, 2, &pc)
	if err == nil {
		t.Error("Expected but did not get an error for superclass name that's empty")
	} else {
		if !strings.HasPrefix(err.Error(), "Class Format Error: invaild empty string for superclass name") {
			t.Error("Expected an invalid string for superclass error, but got: " + err.Error())
		}
	}
}
