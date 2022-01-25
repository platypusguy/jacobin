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
func TestInvalidMagicNumber(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()

	// redirect stderr to inspect output
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	bytesToTest := []byte{0xCA, 0xFE, 0xBA, 0xBA, 0x00, 0x00, 0xFF, 0xF0}
	err := parseMagicNumber(bytesToTest)

	// restore stderr to what it was before
	_ = w.Close()
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

func TestTooShortMagicNumber(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()

	// redirect stderr to inspect output
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	bytesToTest := []byte{0xCA, 0xFE, 0xBA}
	err := parseMagicNumber(bytesToTest)

	// restore stderr to what it was before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	if err == nil {
		t.Error("Too short Java magic number did not generate an error")
	}

	if !strings.Contains(msg, "invalid magic number") {
		t.Error("Did not get expected error msg for invalid magic number. Got: " + msg)
	}
}

func TestValidMagicNumber(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()

	// redirect stderr to inspect output
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	bytesToTest := []byte{0xCA, 0xFE, 0xBA, 0xBE, 0x00, 0x00, 0xFF, 0xF0}
	err := parseMagicNumber(bytesToTest)

	// restore stderr to what it was before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	if err != nil {
		t.Error("Valid Java magic number generate an unexpected error: " + msg)
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
	err := parseJavaVersionNumber(bytesToTest, &ParsedClass{})

	// restore stderr to what it was before
	_ = w.Close()
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
	err := parseJavaVersionNumber(bytesToTest, &ParsedClass{})
	if err != nil {
		t.Error("valid Java version # generated an error in version # parser")
	}
}

func TestConstantPoolCountValid(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	pClass := ParsedClass{}

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
	err := getConstantPoolCount(bytesToTest, &ParsedClass{})

	// restore stderr to what it was before
	_ = w.Close()
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
	_ = log.SetLogLevel(log.WARNING)

	pc := ParsedClass{}
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

func TestClassNameInvalidLocation(t *testing.T) {

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)

	pc := ParsedClass{}
	bytes := []byte{0x00, 0x00, 0x10}
	_, err := parseClassName(bytes, 4, &pc)

	if err == nil {
		t.Error("Should have returned an error for invalid value in class name item")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestClassNameInvalidIndex(t *testing.T) {

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)

	pc := ParsedClass{}
	bytes := []byte{0x00, 0x00, 0x10}
	_, err := parseClassName(bytes, 0, &pc)

	if err == nil {
		t.Error("Should have returned an error for invalid value in class name item")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

// a complex test. It first parses a minimal constant pool that has the records we need
// for the actual test. It then passes bytes containing the class name entry and tests
// whether all the records and pointers in the CP structs point to the right entry.
func TestClassNameValidName(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)

	pc := ParsedClass{}
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
	_ = log.SetLogLevel(log.WARNING)

	pc := ParsedClass{}
	pc.cpCount = 3
	bytes := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00, // the required first 10 bytes
		0x00, 0x00, 0x37, 0x00, 0x03, // Java 8, CP with 3 entries (plus the dummy entry)
		0x08, 0x00, 0x02, // entry #1, should be a ClassRef entry, but is not
		0x01, 0x00, 0x05, 'H', 'e', 'l', 'l', 'o', // entry #2, the UTF-8 record containing "Hello"
	}

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

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

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

// see the previous tests for explanation of the setup. Here we test whether a class name entry
// that points to a valid ClassRef record, but when that ClassRef record does not itself point
// to an expected UTF-8 entry, that the right error is issued.
func TestClassNameWithMissingUTF8(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)

	pc := ParsedClass{}
	pc.cpCount = 3
	bytes := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00, // the required first 10 bytes
		0x00, 0x00, 0x37, 0x00, 0x03, // Java 8, CP with 3 entries (plus the dummy entry)
		0x07, 0x00, 0x02, // entry #1, a ClassRef that should point to a UTF-8 entry
		0x07, 0x00, 0x01, // entry #2, this should be a UTF-8 entry, but it's not
	}

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	_, err := parseConstantPool(bytes, &pc)
	if err != nil {
		t.Error("Error parsing test CP for setup in testing ClassName")
	}

	testBytes := []byte{0x00, 0x00, 0x01} // 3 bytes b/c first byte is skipped. So, this points to entry 1
	_, err = parseClassName(testBytes, 0, &pc)
	if err == nil {
		t.Error("Parse of class name field should have generated an error but it did not.")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestSuperclassNameEmpty(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)

	pc := ParsedClass{}
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

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

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
		if !strings.HasPrefix(err.Error(), "Class Format Error: invalid empty string for superclass name") {
			t.Error("Expected an invalid string for superclass error, but got: " + err.Error())
		}
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestSuperclassNameInvalidIndex(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	testBytes := []byte{0xFF, 0xFF, 0xFF, 0xFF}

	pc := ParsedClass{} // create a CP with a single entry
	pc.cpIndex = append(pc.cpIndex, cpEntry{
		entryType: 1,
		slot:      1,
	})

	_, err := parseSuperClassName(testBytes, 0, &pc)
	if err == nil {
		t.Error("Expected but did not get an error for superclass name that's empty")
	} else {
		if !strings.HasPrefix(err.Error(), "Class Format Error: invalid index into CP for superclass name") {
			t.Error("Expected an invalid index for superclass name error, but got: " + err.Error())
		}
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestSuperclassNameNotClassref(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	testBytes := []byte{0x00, 0x00, 0x01, 0xFF}

	pc := ParsedClass{} // create a CP with a single valid entry to which the superclass points incorrectly
	pc.cpIndex = append(pc.cpIndex, cpEntry{
		entryType: 0,
		slot:      0,
	})
	pc.cpIndex = append(pc.cpIndex, cpEntry{
		entryType: 0,
		slot:      1,
	})

	_, err := parseSuperClassName(testBytes, 0, &pc)
	if err == nil {
		t.Error("Expected but did not get an error for superclass that's not a Classref")
	} else {
		if !strings.Contains(err.Error(), "invalid entry for superclass name") {
			t.Error("Expected error that superclass entry is not a Classref: " + err.Error())
		}
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestSuperclassNameNoneButNotObectClass(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	testBytes := []byte{0x00, 0x00, 0x00, 0xFF}

	pc := ParsedClass{} // create a CP with a single valid entry to which the superclass points incorrectly
	pc.cpIndex = append(pc.cpIndex, cpEntry{
		entryType: 0,
		slot:      0,
	})
	pc.cpIndex = append(pc.cpIndex, cpEntry{
		entryType: 0,
		slot:      1,
	})
	pc.className = "Allo!"

	_, err := parseSuperClassName(testBytes, 0, &pc)
	if err == nil {
		t.Error("Expected but did not get an error for missing superclass when class is not Object.class")
	} else {
		if !strings.Contains(err.Error(), "class is not java/lang/Object") {
			t.Error("Expected superclass name error for non-Object class, but got: " + err.Error())
		}
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestInterfaceCountValid(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)

	pc := ParsedClass{}

	bytesToTest := []byte{
		0x00, 0x00, 0x12,
	}

	_, err := parseInterfaceCount(bytesToTest, 0, &pc)
	if err != nil {
		t.Error("Unexpected error reading interface count")
	}

	if pc.interfaceCount != 18 {
		t.Error("Expecting interface count of 18. Got: " + strconv.Itoa(pc.interfaceCount))
	}
}

func TestInterfaceCountInvalid(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	pc := ParsedClass{}

	bytesToTest := []byte{ // too short to extract a count
		0x00,
	}

	_, err := parseInterfaceCount(bytesToTest, 0, &pc)
	if err == nil {
		t.Error("Did not get error fetching invalid interface count")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestInterfaceValid(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	klass := ParsedClass{}
	klass.cpIndex = append(klass.cpIndex, cpEntry{})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 0}) // the UTF-8 reference
	klass.cpIndex = append(klass.cpIndex, cpEntry{ClassRef, 0})
	klass.classRefs = append(klass.classRefs, 1) // -> cpIndex[1] -> UTF8 entry
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"gherkin"})
	klass.cpCount = 3
	klass.interfaceCount = 1

	testBytes := []byte{0x00, 0x00, 0x02} // 3 bytes b/c first byte is skipped.
	_, err := parseInterfaces(testBytes, 0, &klass)

	if err != nil {
		t.Error("Unexpected error testing parse of interface")
	}

	i := klass.interfaces[0]
	if klass.utf8Refs[i].content != "gherkin" {
		t.Error("Expecting interface of 'gherkin' but got: " + klass.utf8Refs[i].content)
	}
	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestInterfaceInvalid(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	klass := ParsedClass{}
	klass.cpIndex = append(klass.cpIndex, cpEntry{})
	klass.cpIndex = append(klass.cpIndex, cpEntry{ClassRef, 0}) // the invalid reference
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 0})
	klass.classRefs = append(klass.classRefs, 1) // -> cpIndex[1] -> UTF8 entry
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"gherkin"})
	klass.cpCount = 3
	klass.interfaceCount = 1

	testBytes := []byte{0x00, 0x00, 0x02} // 3 bytes b/c first byte is skipped.
	_, err := parseInterfaces(testBytes, 0, &klass)

	if err == nil {
		t.Error("Expected an error while testing parse of interface, but got none")
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr
	msg := string(out[:])

	if !strings.Contains(msg, "Interface index does not point to a class type") {
		t.Error("Did not get expected error msg. Got: " + msg)
	}

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestFieldCountValid(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	pc := ParsedClass{}
	pc.fieldCount = 200

	bytesToTest := []byte{ // count of 10 fields
		0x00, 0x00, 0x0A,
	}

	_, err := parseFieldCount(bytesToTest, 0, &pc)
	if err != nil {
		t.Error("Got unexpected error on valid field count: " + err.Error())
	}

	if pc.fieldCount != 10 {
		t.Errorf("expected field count to be 10, but got: %d", pc.fieldCount)
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestFieldCountInvalid(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	pc := ParsedClass{}
	pc.fieldCount = 200

	bytesToTest := []byte{ // too short to extract a count
		0x00,
	}

	_, err := parseFieldCount(bytesToTest, 0, &pc)
	if err == nil {
		t.Error("Did not get expected error on invalid field count")
	}

	if pc.fieldCount != 200 {
		t.Errorf("expected field count to be 200, but got: %d", pc.fieldCount) // original count should not change
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestFieldWithFaultyNameIndex(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	klass := ParsedClass{}
	klass.cpIndex = append(klass.cpIndex, cpEntry{})
	klass.cpIndex = append(klass.cpIndex, cpEntry{ClassRef, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 0})
	klass.classRefs = append(klass.classRefs, 1) // -> cpIndex[1] -> UTF8 entry
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"gherkin"})
	klass.cpCount = 3
	klass.fieldCount = 1

	testBytes := []byte{0x00, 0x00, 0x02, 0x00, 0x10} //  first byte is skipped.
	// last byte of testBytes says that the name index is in slot 16 (0x10), which
	// exceeds the size of the CP (stored in klass.cpCount). This should generate error
	_, err := parseFields(testBytes, 0, &klass)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout

	if err == nil {
		t.Error("Expected an error while testing parse of field, but got none")
	}
}

func TestFieldWithNoAttributes(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	klass := ParsedClass{}
	klass.cpIndex = append(klass.cpIndex, cpEntry{})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 1})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"gherkin"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"pickle"})
	klass.cpCount = 3
	klass.fieldCount = 1

	testBytes := []byte{
		0x00,       // first byte is skipped
		0x00, 0x02, // access flags
		0x00, 0x01, // UTF-8 record for the nameIndex, i.e., "gherkin"
		0x00, 0x02, // UTF-8 record for the descIndex, i.e., "pickle
		0x00, 0x00, // attribute count: there are none here
	}

	_, err := parseFields(testBytes, 0, &klass)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout

	if err != nil {
		t.Error("Expected no error in parsing field, but got one.")
	}

	if len(klass.fields) != 1 {
		t.Error("Expected 1 field entry in parsed class, but got none")
	}

	f := klass.fields[0]
	if f.accessFlags != 2 {
		t.Error("Expected a field access flag of 2, got: " + strconv.Itoa(f.accessFlags))
	}

	if f.name != 0 {
		t.Error("Expected a field name UTF entry at 0, got: " + strconv.Itoa(f.name))
	}

	if f.description != 1 {
		t.Error("Expectef a field description UTF entry at 1, got: " + strconv.Itoa(f.description))
	}

	if len(f.attributes) != 0 {
		t.Error("Expected a field attribute count of 0, got: " + strconv.Itoa(len(f.attributes)))
	}
}

func TestMethodCountValid(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	pc := ParsedClass{}
	pc.methodCount = 200

	bytesToTest := []byte{ // count of 10 methods
		0x00, 0x00, 0x0A,
	}

	_, err := parseMethodCount(bytesToTest, 0, &pc)
	if err != nil {
		t.Error("Got unexpected error on valid method count: " + err.Error())
	}

	if pc.methodCount != 10 {
		t.Errorf("expected method count to be 10, but got: %d", pc.methodCount)
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestMethodCountInvalid(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	pc := ParsedClass{}
	pc.methodCount = 200

	bytesToTest := []byte{ // too short to extract a count
		0x00,
	}

	_, err := parseMethodCount(bytesToTest, 0, &pc)
	if err == nil {
		t.Error("Did not get expected error on invalid method count")
	}

	if pc.methodCount != 200 {
		t.Errorf("expected method count to be 200, but got: %d", pc.methodCount) // original count should not change
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestParseClassAttributeCountFor2Attributes(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	// redirect stderr & stdout to prevent error message from showing up in the test results
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	klass := ParsedClass{}
	klass.attribCount = 0

	testBytes := []byte{
		0x00,       // first byte is skipped
		0x00, 0x02, // two Class attributes
	}

	pos, err := parseClassAttributeCount(testBytes, 0, &klass)
	if err != nil {
		t.Error("Unexpected error getting class attribute count")
	}

	if pos != 2 {
		t.Error("Returned position from class attribute count parse should be 2. Got: " +
			strconv.Itoa(pos))
	}

	if klass.attribCount != 2 {
		t.Error("Expecting a class attribute count of 2, but got: " +
			strconv.Itoa(pos))
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestValidBootstrapClassAttribute(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	klass := ParsedClass{}
	klass.cpIndex = append(klass.cpIndex, cpEntry{})
	klass.cpIndex = append(klass.cpIndex, cpEntry{LongConst, 0}) // LongConst (which is loadable)
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{MethodHandle, 0})

	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"BootstrapMethods"})
	klass.cpCount = 4
	klass.attribCount = 1

	// the attribute bytes. There's a leading dummy byte b/c the fetch routine starts
	// at 1 byte after the passed-in position.
	bytes := []byte{00, // dummy byte
		00, 02, // CP[2] -> UTF8[0] -> "BootstrapMethods"
		00, 00, 00, 0x08, // length of attribute (must be 0 for 'Deprecated')
		00, 01, // bootstrap count
		00, 03, // CP[3] -> MethodHandle
		00, 01, // arg count for this first boostrap method
		00, 01, // CP[1] -> LongConst (which is loadable, so all is good.
	}

	_, err := parseClassAttributes(bytes, 0, &klass)
	if err != nil {
		t.Error("Unexpected error in test of parseClassAttributes()")
	}

	if len(klass.bootstraps) != 1 {
		t.Error("Class should have 1 bootstrap methods. Got: " + strconv.Itoa(len(klass.bootstraps)))
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}

func TestDeprecatedClassAttribute(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	klass := ParsedClass{}
	klass.cpIndex = append(klass.cpIndex, cpEntry{})
	klass.cpIndex = append(klass.cpIndex, cpEntry{1, 0}) // UTF-8 rec w/ attribute name
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"Deprecated"})
	klass.cpCount = 2
	klass.attribCount = 1

	// the attribute bytes. There's a leading dummy byte b/c the fetch routine starts
	// at 1 byte after the passed-in position. So here we have a name index of 01, which
	// points to the first entry in the CP above. That entry points to the first UTF-8
	// record, which is in position 0 in the utf8Refs and has a value of "Deprecated"
	// The next four bytes are the length of the remaining
	// bytes in the attribute. For Deprecated the length must be 0.
	bytes := []byte{00, // dummy byte
		00, 01, // CP[1] -> UTF8[0] -> "Deprecated"
		00, 00, 00, 00} // length of attribute (must be 0 for 'Deprecated')

	_, err := parseClassAttributes(bytes, 0, &klass)
	if err != nil {
		t.Error("Unexpected error in test of parseClassAttributes()")
	}

	if !klass.deprecated {
		t.Error("klass.deprecated should be true, but it's false")
	}

	if len(klass.attributes) != 1 {
		t.Error("number of klass attributes should be 1, but got: " +
			strconv.Itoa(len(klass.attributes)))
	}

	// restore stderr and stdout to what they were before
	_ = w.Close()
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout
}
