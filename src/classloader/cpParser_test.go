/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"jacobin/globals"
	"jacobin/log"
	"strconv"
	"testing"
)

func TestCPvalidUTF8Ref(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBA, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x01, 0x00, 0x04, 'J', 'A',
		'C', 'O',
	}

	pc := parsedClass{}
	pc.cpCount = 2
	loc, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Error("Parsing valid CP UTF-8 entry (01) generated an unexpected error")
	}

	if loc != 16 {
		t.Error("Was expecting a new position of 16, but got: " + strconv.Itoa(loc))
	}

	if len(pc.utf8Refs) != 1 {
		t.Error("Was expecting the UTF8 ref array to have 1 entry, but it has: " + strconv.Itoa(len(pc.utf8Refs)))
	}

	ute := pc.utf8Refs[0]
	if ute.content != "JACO" {
		t.Error("Was expecting a UTF-8 string of 'JACO', but got: " + ute.content)
	}

	if len(pc.cpIndex) != 2 {
		t.Error("Was expecting pc.cpIndex to have 2 entries, but instead got: " + strconv.Itoa(len(pc.cpIndex)))
	}
}

func TestCPvalidIntConst(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBA, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x03, 0x01, 0x05, 0x20, 0x44,
	}

	pc := parsedClass{}
	pc.cpCount = 2
	loc, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Error("Parsing valid CP integer constant generated an unexpected error")
	}

	if loc != 14 {
		t.Error("Was expecting a new position of 14, but got: " + strconv.Itoa(loc))
	}

	if len(pc.intConsts) != 1 {
		t.Error("Was expecting the int const array to have 1 entry, but it has: " + strconv.Itoa(len(pc.intConsts)))
	}

	ice := pc.intConsts[0]
	if ice.value != 17113156 {
		t.Error("Was expecting an integer constant of 17113156, but got: " + strconv.Itoa(ice.value))
	}

	if len(pc.cpIndex) != 2 {
		t.Error("Was expecting pc.cpIndex to have 2 entries, but instead got: " + strconv.Itoa(len(pc.cpIndex)))
	}
}

func TestCPvalidClassRef(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBA, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x07, 0x02, 0x05,
	}

	pc := parsedClass{}
	pc.cpCount = 2
	loc, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Error("Parsing valid CP class reference (7) generated an unexpected error")
	}

	if loc != 12 {
		t.Error("Was expecting a new position of 12, but got: " + strconv.Itoa(loc))
	}

	if len(pc.classRefs) != 1 {
		t.Error("Was expecting the class ref array to have 1 entry, but it has: " + strconv.Itoa(len(pc.classRefs)))
	}

	cre := pc.classRefs[0]
	if cre.index != 517 {
		t.Error("Was expecting a class ref index of 517, but got: " + strconv.Itoa(cre.index))
	}

	if len(pc.cpIndex) != 2 {
		t.Error("Was expecting pc.cpIndex to have 2 entries, but instead got: " + strconv.Itoa(len(pc.cpIndex)))
	}
}

func TestCPvalidStringConstRef(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBA, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x08, 0x00, 0x20,
	}

	pc := parsedClass{}
	pc.cpCount = 2
	loc, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Error("Parsing valid CP StringConstRef (8) generated an unexpected error")
	}

	if loc != 12 {
		t.Error("Was expecting a new position of 12, but got: " + strconv.Itoa(loc))
	}

	if len(pc.stringRefs) != 1 {
		t.Error("Was expecting the string const ref array to have 1 entry, but it has: " + strconv.Itoa(len(pc.stringRefs)))
	}

	sre := pc.stringRefs[0]
	if sre.index != 32 {
		t.Error("Was expecting a string ref index of 32, but got: " + strconv.Itoa(sre.index))
	}

	if len(pc.cpIndex) != 2 {
		t.Error("Was expecting pc.cpIndex to have 2 entries, but instead got: " + strconv.Itoa(len(pc.cpIndex)))
	}
}

func TestCPvalidFieldRef(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBA, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x09, 0x00, 0x14, 0x01, 0x01,
	}

	pc := parsedClass{}
	pc.cpCount = 2
	loc, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Error("Parsing valid CP FieldRef (09) generated an unexpected error")
	}

	if loc != 14 {
		t.Error("Was expecting a new position of 14, but got: " + strconv.Itoa(loc))
	}

	if len(pc.fieldRefs) != 1 {
		t.Error("Was expecting the field ref array to have 1 entry, but it has: " + strconv.Itoa(len(pc.fieldRefs)))
	}

	fre := pc.fieldRefs[0]
	if fre.classIndex != 20 {
		t.Error("Was expecting a field ref classIndex of 20, but got: " + strconv.Itoa(fre.classIndex))
	}

	if fre.nameAndTypeIndex != 257 {
		t.Error("Was expecting a field ref nameAndTypeIndex of 257, but got: " + strconv.Itoa(fre.nameAndTypeIndex))
	}

	if len(pc.cpIndex) != 2 {
		t.Error("Was expecting pc.cpIndex to have 2 entries, but instead got: " + strconv.Itoa(len(pc.cpIndex)))
	}
}

func TestCPvalidMethodRef(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBA, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x0A, 0x00, 0x15, 0x01, 0x06,
	}

	pc := parsedClass{}
	pc.cpCount = 2
	loc, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Error("Parsing valid CP MethodRef (10) generated an unexpected error")
	}

	if loc != 14 {
		t.Error("Was expecting a new position of 14, but got: " + strconv.Itoa(loc))
	}

	if len(pc.methodRefs) != 1 {
		t.Error("Was expecting the method ref array to have 1 entry, but it has: " + strconv.Itoa(len(pc.methodRefs)))
	}

	mre := pc.methodRefs[0]
	if mre.classIndex != 21 {
		t.Error("Was expecting a method ref classIndex of 21, but got: " + strconv.Itoa(mre.classIndex))
	}

	if mre.nameAndTypeIndex != 262 {
		t.Error("Was expecting a method ref nameAndType of 262, but got: " + strconv.Itoa(mre.nameAndTypeIndex))
	}

	if len(pc.cpIndex) != 2 {
		t.Error("Was expecting pc.cpIndex to have 2 entries, but instead got: " + strconv.Itoa(len(pc.cpIndex)))
	}
}

func TestCPvalidNameAndTypeEntry(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBA, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x0C, 0x00, 0x14, 0x01, 0x01,
	}

	pc := parsedClass{}
	pc.cpCount = 2
	loc, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Error("Parsing valid CP NameAndType (12) generated an unexpected error")
	}

	if loc != 14 {
		t.Error("Was expecting a new position of 14, but got: " + strconv.Itoa(loc))
	}

	if len(pc.nameAndTypes) != 1 {
		t.Error("Was expecting the nameAndTypes array to have 1 entry, but it has: " + strconv.Itoa(len(pc.nameAndTypes)))
	}

	nte := pc.nameAndTypes[0]
	if nte.nameIndex != 20 {
		t.Error("Was expecting a nameAndType nameIndex of 20, but got: " + strconv.Itoa(nte.nameIndex))
	}

	if nte.descriptorIndex != 257 {
		t.Error("Was expecting a nameAndType descriptor index of 257, but got: " + strconv.Itoa(nte.descriptorIndex))
	}

	if len(pc.cpIndex) != 2 {
		t.Error("Was expecting pc.cpIndex to have 2 entries, but instead got: " + strconv.Itoa(len(pc.cpIndex)))
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
		t.Error("Should have returned an error for invalid value in class name entry")
	}
}

// a complex test. It first parses a minimal constant pool that has the records we need
// for the actual test. It then passes bytes containing the class name entry and tests
// whether all the records and pointers in the CP structs point to the right entry.
func TestClassNameValidUTF8(t *testing.T) {

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
	if err.Error() != "Class Format Error: invalid entry for class name" {
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

	if err.Error() != "Class Format Error: error classRef in CP does not point to a UTF-8 string" {
		t.Error("Expected error msg about invalid UTF-8 entry for class name. Got: " + err.Error())
	}
}
