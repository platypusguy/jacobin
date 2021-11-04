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

// Tests for parsing of CP entries. These tests are sequenced according
// to the CP entry number for that record:
// 0 - Dummy entry					TestDummyEntry
// 1 - UTF							TestCPvalidUTF8Ref
// 3 - IntConst						TestCPvalidIntConst
// 4 - FloatConst					TestCPvalidFloatConst
// 5 - LongConst 		 			TestCPvalidLongConst
// 6 - DoubleConst					TestCPvalidDoubleConst
// 7 - ClassRef						TestCPvalidClassRef
// 8 - StringConst					TestCPvalidStringConstRef
// 9 - FieldRef						TestCPvalidFieldRef
// 10- MethodRef					TestCPvalidMethodRef
// 11- Interface					TestCPvalidInterface
// 12- NameAndTypeEntry				TestCPvalidNameAndTypeEntry
// 15- MethodHandle  	 			TestCPvalidMethodHandle
// 16- MethodType 		 			TestCPvalidMethodType
// 18- InvokeDynamic 	 			TestCPvalidInvokeDynamic
// 19- ModuleName					see TestPrintOfCP2
//
// Printing of CP contents			TestPrintOfCP

// Pass in a CP with a single UTF8 entry and make sure the first CP entry
// (CP[0]) is a dummy entry as it should be.
func TestDummyEntry(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x01, 0x00, 0x04, 'J', 'A',
		'C', 'O',
	}

	pc := ParsedClass{}
	pc.cpCount = 2
	_, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Error("Parsing valid CP UTF-8 entry (01) generated an unexpected error")
	}

	if pc.cpIndex[0].entryType != Dummy {
		t.Error("Parsing a valid CP did not result a dummy entry at CP[0]")
	}
}

func TestCPvalidUTF8Ref(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x01, 0x00, 0x04, 'J', 'A',
		'C', 'O',
	}

	pc := ParsedClass{}
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
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x03, 0x01, 0x05, 0x20, 0x44,
	}

	pc := ParsedClass{}
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
	if ice != 17113156 {
		t.Error("Was expecting an integer constant of 17113156, but got: " + strconv.Itoa(ice))

	}

	if len(pc.cpIndex) != 2 {
		t.Error("Was expecting pc.cpIndex to have 2 entries, but instead got: " + strconv.Itoa(len(pc.cpIndex)))
	}
}

func TestCPvalidLongConst(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x05, 0x00, 0x00, 0x00, 0x01, // first four bytes of long
		0x00, 0x00, 0x00, 0x02, // second four bytes of long
	}

	pc := ParsedClass{}
	pc.cpCount = 3 // it's 3 b/c the long constant takes up two slots
	loc, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Error("Parsing valid CP long constant generated an unexpected error")
	}

	if loc != 18 {
		t.Error("Was expecting a new position of 18, but got: " + strconv.Itoa(loc))
	}

	if len(pc.longConsts) != 1 {
		t.Error("Was expecting the long const array to have 1 entry, but it has: " + strconv.Itoa(len(pc.intConsts)))
	}

	long := pc.longConsts[0]
	if long != 4294967298 {
		longInt := int(long)
		t.Error("Was expecting an long constant of 4294967298, but got: " + strconv.Itoa(longInt))

	}

	if len(pc.cpIndex) != 3 { // the dummy entry + 2 slots for the long
		t.Error("Was expecting pc.cpIndex to have 3 entries, but instead got: " + strconv.Itoa(len(pc.cpIndex)))
	}
}

func TestCPvalidFloatConst(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x04, // Double constant
		// Big endian hex value of 40 09 21 F9 F0 1B 86 6E should be a double of value: 3.14159
		//
		0x40, 0x09, 0x21, 0xF9, // ffour bytes of float
	}

	pc := ParsedClass{}
	pc.cpCount = 2 //
	loc, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Error("Parsing valid CP float constant generated an unexpected error")
	}

	if loc != 14 {
		t.Error("Was expecting a new position of 14, but got: " + strconv.Itoa(loc))
	}

	if len(pc.floats) != 1 {
		t.Error("Was expecting the double const array to have 1 entry, but it has: " +
			strconv.Itoa(len(pc.floats)))
	}

	float := pc.floats[0]
	if float != 2.14269853 { // precision of value is low enough that exact match is possible.
		bigFloat := float64(float)
		t.Error("Was expecting a value of 2.14269853, but got: " +
			strconv.FormatFloat(bigFloat, 'E', -1, 32))
	}
}

func TestCPvalidDoubleConst(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x06, // Double constant
		// Big endian hex value of 40 09 21 F9 F0 1B 86 6E should be a double of value: 3.14159
		//
		0x40, 0x09, 0x21, 0xF9, // first four bytes of double
		0xF0, 0x1B, 0x86, 0x6E, // second four bytes of double
	}

	pc := ParsedClass{}
	pc.cpCount = 3 // it's 3 b/c the long constant takes up two slots
	loc, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Error("Parsing valid CP double constant generated an unexpected error")
	}

	if loc != 18 {
		t.Error("Was expecting a new position of 18, but got: " + strconv.Itoa(loc))
	}

	if len(pc.doubles) != 1 {
		t.Error("Was expecting the double const array to have 1 entry, but it has: " +
			strconv.Itoa(len(pc.doubles)))
	}

	double := pc.doubles[0]
	if double != 3.14159 { // because of the low precision of the value, a direct comparison should work
		t.Error("Was expecting a value of 3.14159, but got: " +
			strconv.FormatFloat(double, 'E', -1, 64))
	}
}

func TestCPvalidClassRef(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x07, 0x02, 0x05,
	}

	pc := ParsedClass{}
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
	if cre != 517 {
		t.Error("Was expecting a class ref index of 517, but got: " + strconv.Itoa(cre))
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
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x08, 0x00, 0x20,
	}

	pc := ParsedClass{}
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
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x09, 0x00, 0x14, 0x01, 0x01,
	}

	pc := ParsedClass{}
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
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x0A, 0x00, 0x15, 0x01, 0x06,
	}

	pc := ParsedClass{}
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

func TestCPvalidInterface(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x0C,       // name and type (12)
		0x00, 0x14, // name and type name index
		0x01, 0x01, // name and type descriptor index
		0x0B,       // interface entry (11)
		0x00, 0x20, // interface class index
		0x00, 0x01, // name and type entry index
	}

	pc := ParsedClass{}
	pc.cpCount = 3
	loc, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Error("Parsing valid CP Interface (11) generated an unexpected error")
	}

	if loc != 19 { // 20 bytes, but 0-based
		t.Error("Was expecting a new position of 19, but got: " + strconv.Itoa(loc))
	}

	if len(pc.interfaceRefs) != 1 {
		t.Error("Was expecting Interfaces to have 1 entry. Got: " +
			strconv.Itoa(len(pc.interfaces)))
	}

	ie := pc.interfaceRefs[0]
	if ie.classIndex != 32 {
		t.Error("Was expecting interface to have a class index of 32. Got: " +
			strconv.Itoa(ie.classIndex))
	}

	if ie.nameAndTypeIndex != 1 {
		t.Error("Was expecting interface to have a name-and-type index of 1. Got: " +
			strconv.Itoa(ie.nameAndTypeIndex))
	}
}

func TestCPvalidNameAndTypeEntry(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x0C, 0x00, 0x14, 0x01, 0x01,
	}

	pc := ParsedClass{}
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

func TestCPvalidMethodHandle(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x0C, // Name and Type
		0x00, 0x14,
		0x01, 0x01,
		0x0F,       // MethodHanlde (15)
		0x05,       // Ref kind (one byte)
		0x00, 0x01, // Ref index
	}

	pc := ParsedClass{}
	pc.cpCount = 3
	loc, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Error("Parsing valid CP MethodHandle (15) generated an unexpected error")
	}

	if loc != 18 {
		t.Error("Was expecting a new position of 18, but got: " + strconv.Itoa(loc))
	}

	if len(pc.methodHandles) != 1 {
		t.Error("Was expecting the methodHandles array to have 1 entry, but it has: " + strconv.Itoa(len(pc.nameAndTypes)))
	}

	mhe := pc.methodHandles[0]
	if mhe.referenceKind != 5 {
		t.Error("Was expecting a methodHandle kind of 5. Got: " + strconv.Itoa(mhe.referenceKind))
	}

	if mhe.referenceIndex != 1 {
		t.Error("Was expecting a methodHandle reference index of 1. Got: " + strconv.Itoa(mhe.referenceIndex))
	}

	if len(pc.cpIndex) != 3 {
		t.Error("Was expecting pc.cpIndex to have 3 entries, but instead got: " + strconv.Itoa(len(pc.cpIndex)))
	}
}

func TestCPvalidMethodType(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x0C, // Name and Type
		0x00, 0x14,
		0x01, 0x01,
		0x10,       // MethodType (16)
		0x00, 0x05, // Desc Index
	}

	pc := ParsedClass{}
	pc.cpCount = 3
	loc, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Error("Parsing valid CP MethodType (16) generated an unexpected error")
	}

	if loc != 17 {
		t.Error("Was expecting a new position of 17, but got: " + strconv.Itoa(loc))
	}

	if len(pc.methodTypes) != 1 {
		t.Error("Was expecting the methodTypes array to have 1 entry, but it has: " + strconv.Itoa(len(pc.nameAndTypes)))
	}

	mte := pc.methodTypes[0]
	if mte != 5 {
		t.Error("Was expecting a methodType kind of 5. Got: " + strconv.Itoa(mte))
	}

	if len(pc.cpIndex) != 3 {
		t.Error("Was expecting pc.cpIndex to have 3 entries, but instead got: " + strconv.Itoa(len(pc.cpIndex)))
	}
}

func TestCPvalidInvokeDynamic(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x0C, // Name and Type entry
		0x00, 0x14,
		0x01, 0x01,
		0x12,       // InvokeDynamic (18)
		0x00, 0x08, // Bootstrap index
		0x00, 0x01, // name and type entry
	}

	pc := ParsedClass{}
	pc.cpCount = 3
	loc, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Error("Parsing valid CP InvokeDynamic (18) generated an unexpected error")
	}

	if loc != 19 {
		t.Error("Was expecting a new position of 19, but got: " + strconv.Itoa(loc))
	}

	if len(pc.invokeDynamics) != 1 {
		t.Error("Was expecting the invokeDynamics array to have 1 entry, but it has: " + strconv.Itoa(len(pc.nameAndTypes)))
	}

	ide := pc.invokeDynamics[0]
	if ide.bootstrapIndex != 8 {
		t.Error("Was expecting an invokeDynamic boostrap index of 8. Got: " + strconv.Itoa(ide.bootstrapIndex))
	}

	if ide.nameAndType != 1 {
		t.Error("Was expecing an invokeDynamic nameAndType index of 1. Got: " + strconv.Itoa(ide.nameAndType))
	}

	if len(pc.cpIndex) != 3 {
		t.Error("Was expecting pc.cpIndex to have 3 entries, but instead got: " + strconv.Itoa(len(pc.cpIndex)))
	}
}

// test whether the info logged to when logging set to FINEST is correct
// This test captures stderr and then does searches on the saved output for
// logging contents. The size of saved stderr output is quite limited and the
// CP logging of all possible entry types exceeds this space, so this test is
// broken up into two tests, of which this is the first.
func TestPrintOfCPpart1(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	log.SetLogLevel(log.FINEST)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x01,       // UTF8 entry
		0x00, 0x04, //		length of UTF8 string
		'J', 'A', //  	contents of UTF8 string
		'C', 'O',

		0x03,       // Integer constant
		0x00, 0x00, // 		value of int (four bytes)
		0x01, 0x02, //  	should = 258

		0x04,       // Float Constant
		0x41, 0x42, //		value of float comprises 4 bytes
		0x43, 0x44,

		0x05,       // Long constant
		0x50, 0x50, //  	value consists of 8 bytes
		0x50, 0x51,
		0x50, 0x52,
		0x50, 0x53,

		0x06,       // Double constant
		0x61, 0x60, //  	value consists of 8 bytes
		0x61, 0x61,
		0x61, 0x62,
		0x61, 0x63,

		0x07,       // Classref
		0x74, 0x75, // 		value consists of 2 bytes

		0x08,       // StringConst
		0x84, 0x85, // 		value consists of 2 bytes

		0x09,       // FieldRef
		0x91, 0x92, //		class index (2 bytes)
		0x93, 0x94, //		name and type index (2 bytes)

		0x0A,       // MethodRef (10)
		0xA0, 0xA1, //		class index (2 bytes)
		0xA2, 0xA3, //		name and type index (2 bytes)

		0x0B,       // Interface (11)
		0xB0, 0xB1, //		class index (2 bytes)
		0xB2, 0xB3, //		name and type index (2 bytes)

		0x0C,       // Name and Type entry (12)
		0xC0, 0xC4, //		value consists of 4 bytes
		0xC1, 0xC1,

		0x0F,       // Method Handle (15)
		0xF1, 0xF2, //		reference kind (2 bytes)
		0xF3, 0xF4, //		reference index (2bytes)

	}

	pc := ParsedClass{}
	pc.cpCount = 15 // Dummy entry/entries plus the number of entries above

	// this parses the CP and logs it to stderr b/c logging is set to FINEST
	_, err := parseConstantPool(bytesToTest, &pc)
	if err != nil {
		t.Error("Unexpected error in parsing CP in testPrintOfCP()")
	}

	// restore stderr and stdout to what they were before
	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	logMsg := string(out[:])

	if !strings.Contains(logMsg, "(dummy entry)") {
		t.Error("Dummy CP entry did not appear in logging of CP contents")
	}

	if !strings.Contains(logMsg, "(UTF-8 string) ") {
		t.Error("UTF8 string CP entry did not appear in logging of CP contents")
	}

	if !strings.Contains(logMsg, "(int constant)") || !strings.Contains(logMsg, "258") {
		t.Error("IntConst CP entry with value 258 did not appear correctly in CP logging")
	}

	if !strings.Contains(logMsg, "(long constant)") {
		t.Error("LongConst CP entry did not appear correctly in CP logging")
	}

	if !strings.Contains(logMsg, "(double constant)") {
		t.Error("DoubleConst CP entry did not appear correctly in CP logging")
	}

	if !strings.Contains(logMsg, "(class ref)") {
		t.Error("ClassRef CP entry did not appear correctly in CP logging")
	}

	if !strings.Contains(logMsg, "(string const ref)") {
		t.Error("StringConst CP entry did not appear correctly in CP logging")
	}

	if !strings.Contains(logMsg, "(field ref)") {
		t.Error("FieldRef CP entry did not appear correctly in CP logging")
	}

	if !strings.Contains(logMsg, "(method ref)") {
		t.Error("MethodRef CP entry did not appear correctly in CP logging")
	}

	if !strings.Contains(logMsg, "(interface ref)") {
		t.Error("InterfaceRef CP entry did not appear correctly in CP logging")
	}

	if !strings.Contains(logMsg, "(name and type) ") {
		t.Error("Name and type CP entry did not appear in logging of CP contents")
	}

	if !strings.Contains(logMsg, "(method handle)") {
		t.Error("Method Handle CP entry did not appear in logging of CP contents")
	}

	_ = wout.Close()
	os.Stdout = normalStdout
}

// see the comment for part 1 of this test
func TestPrintOfCPpart2(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	log.SetLogLevel(log.FINEST)

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,

		0x10,       // Method Type (16)
		0x11, 0x12, //		description index (2 bytes)

		0x12,       // InvokeDynamic (18)
		0x12, 0x08, // 		Bootstrap index
		0x12, 0x01, // 		name and type entry

		0x01,       // UTF-8 String (1)
		0x00, 0x06, //		length of UTF8 string
		'M', 'o', //  	contents of UTF8 string
		'd', 'u',
		'l', 'e',

		0x13,       // Module name (19)
		0x00, 0x03, // CP[3] -> UTF8 rec with name of module: "Module"

		0x14,       // Package name (20)
		0x00, 0x03, // CP[3] -> UTF8 rec with name of package: "Module"

		// The following entries and the assertions they correspond to below are
		// commented out due to the problem described in JACOBIN-86.

		0x11,       // Dynamic (17)
		0x12, 0x08, // 		Bootstrap index
		0x12, 0x01, // 		name and type entry
		//
		// The following UTF8 record was added during the attempt to diagnose the
		// problem with JACOBIN-86. It is not otherwise needed and can be deleted.
		0x01,       // UTF-8 String (1)
		0x00, 0x05, //		length of UTF8 string
		'H', 'e', //  	contents of UTF8 string
		'l', 'l', // added to see whether it solves the missing 'dynamic' entry on GitHub
		'o',
	}

	pc := ParsedClass{}
	pc.cpCount = 6 // Dummy entry/entries plus the number of entries above

	pc.javaVersion = 55 // Java 11
	pc.moduleName = "Module"

	// this parses the CP and logs it to stderr b/c logging is set to FINEST
	_, err := parseConstantPool(bytesToTest, &pc)
	if err != nil {
		t.Error("Unexpected error in parsing CP in testPrintOfCP()")
	}

	// restore stderr and stdout to what they were before
	w.Close()
	out, _ := ioutil.ReadAll(r)

	logMsg := string(out[:])

	if !strings.Contains(logMsg, "(method type) ") {
		t.Error("MethodType CP entry did not appear in logging of CP contents")
	}
	//
	// if !strings.Contains(logMsg, "(dynamic) ") {
	// 	t.Error("dynamic CP entry did not appear in logging of CP contents:\n" + logMsg)
	// }

	if !strings.Contains(logMsg, "(invokedynamic) ") {
		t.Error("invokedynamic CP entry did not appear in logging of CP contents")
	}

	if !strings.Contains(logMsg, "(package name) ") {
		t.Error("package name CP entry did not appear in logging of CP contents" +
			"Output: " + logMsg)
	}

	_ = wout.Close()
	os.Stdout = normalStdout
	os.Stderr = normalStderr
}
