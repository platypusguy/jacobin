/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"io"
	"jacobin/src/globals"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
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
	trace.Init()

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
	trace.Init()

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
	trace.Init()

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
	trace.Init()

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
	trace.Init()

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
	trace.Init()

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
	trace.Init()

	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x07, 0x00, 0x02, // Class reference pointing to CP index 2
		0x01, 0x00, 0x05, byte('h'), byte('e'), byte('l'), byte('l'), byte('o'),
	}

	pc := ParsedClass{}
	pc.cpCount = 3
	loc, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Error("Parsing valid CP class reference (7) generated an unexpected error")
	}

	if loc != 20 {
		t.Error("Was expecting a new position of 20, but got: " + strconv.Itoa(loc))
	}

	if len(pc.classRefs) != 1 {
		t.Error("Was expecting the class ref array to have 1 entry, but it has: " + strconv.Itoa(len(pc.classRefs)))
	}

	cre := pc.classRefs[0]
	className := stringPool.GetStringPointer(cre)
	if *className != "hello" {
		t.Errorf("Was expecting a class ref to point to 'hello', but got: %s", *className)
	}

	if len(pc.cpIndex) != 3 {
		t.Error("Was expecting pc.cpIndex to have 3 entries, but instead got: " + strconv.Itoa(len(pc.cpIndex)))
	}
}

func TestCPvalidStringConstRef(t *testing.T) {

	globals.InitGlobals("test")
	trace.Init()

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
	trace.Init()

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
	trace.Init()

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
	trace.Init()

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
	trace.Init()

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
	trace.Init()

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
	trace.Init()

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
	trace.Init()

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
		t.Error("Was expecting an invokeDynamic bootstrap index of 8. Got: " + strconv.Itoa(ide.bootstrapIndex))
	}

	if ide.nameAndType != 1 {
		t.Error("Was expecing an invokeDynamic nameAndType index of 1. Got: " + strconv.Itoa(ide.nameAndType))
	}

	if len(pc.cpIndex) != 3 {
		t.Error("Was expecting pc.cpIndex to have 3 entries, but instead got: " + strconv.Itoa(len(pc.cpIndex)))
	}
}

// === the following tests were generated by JetBrains Junie to fill in testing gaps  ===

func TestParseConstantPool_InvalidEntryType(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// Use an invalid entry type (99) that doesn't exist
	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x63, // Invalid entry type 99
	}

	pc := ParsedClass{}
	pc.cpCount = 2
	_, err := parseConstantPool(bytesToTest, &pc)

	// Should not error but should handle gracefully
	if err != nil {
		t.Errorf("Unexpected error for invalid entry type: %v", err)
	}
}

func TestParseConstantPool_ModuleEntry_JavaVersionTooLow(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// Module entry with Java version < 53
	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x13, 0x00, 0x01, // Module entry pointing to index 1
	}

	pc := ParsedClass{}
	pc.cpCount = 2
	pc.javaVersion = 52 // Java 8

	_, err := parseConstantPool(bytesToTest, &pc)
	if err == nil {
		t.Error("Expected error for Module entry with Java version < 53")
	}
	if !strings.Contains(err.Error(), "Java module record requires Java 9 or later") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestParseConstantPool_PackageEntry_JavaVersionTooLow(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// Package entry with Java version < 53
	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x14, 0x00, 0x01, // Package entry pointing to index 1
	}

	pc := ParsedClass{}
	pc.cpCount = 2
	pc.javaVersion = 52 // Java 8

	_, err := parseConstantPool(bytesToTest, &pc)
	if err == nil {
		t.Error("Expected error for Package entry with Java version < 53")
	}
	if !strings.Contains(err.Error(), "Java package entry requires Java 9 or later") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestParseConstantPool_DuplicateModuleName(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// Create a UTF-8 entry first, then two Module entries
	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x01, 0x00, 0x04, 'T', 'e', 's', 't', // UTF-8 "Test"
		0x13, 0x00, 0x01, // Module entry pointing to index 1
		0x13, 0x00, 0x01, // Duplicate Module entry
	}

	pc := ParsedClass{}
	pc.cpCount = 4
	pc.javaVersion = 53
	pc.className = "TestClass"

	_, err := parseConstantPool(bytesToTest, &pc)
	if err == nil {
		t.Error("Expected error for duplicate module names")
	}
	if !strings.Contains(err.Error(), "has two module names") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestParseConstantPool_DuplicatePackageName(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// Create a UTF-8 entry first, then two Package entries
	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x01, 0x00, 0x08, 'c', 'o', 'm', '/', 't', 'e', 's', 't', // UTF-8 "com/test"
		0x14, 0x00, 0x01, // Package entry pointing to index 1
		0x14, 0x00, 0x01, // Duplicate Package entry
	}

	pc := ParsedClass{}
	pc.cpCount = 4
	pc.javaVersion = 53
	pc.className = "TestClass"

	_, err := parseConstantPool(bytesToTest, &pc)
	if err == nil {
		t.Error("Expected error for duplicate package names")
	} else if !strings.Contains(err.Error(), "has two package names") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestParseConstantPool_InvalidClassRefIndex_Zero(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// Class reference with index 0 (invalid)
	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x07, 0x00, 0x00, // ClassRef with index 0
	}

	pc := ParsedClass{}
	pc.cpCount = 2

	_, err := parseConstantPool(bytesToTest, &pc)
	if err == nil {
		t.Error("Expected error for ClassRef with index 0")
	}
	if !strings.Contains(err.Error(), "invalid index into CP for class name: 0") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestParseConstantPool_InvalidClassRefIndex_OutOfBounds(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// Class reference with out-of-bounds index
	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x07, 0x00, 0x05, // ClassRef with index 5 (out of bounds for cpCount=2)
	}

	pc := ParsedClass{}
	pc.cpCount = 2

	_, err := parseConstantPool(bytesToTest, &pc)
	if err == nil {
		t.Error("Expected error for ClassRef with out-of-bounds index")
	}
	if !strings.Contains(err.Error(), "invalid index into CP for class name: 5") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestParseConstantPool_UTF8_ZeroLength(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// UTF-8 entry with zero length
	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x01, 0x00, 0x00, // UTF-8 with zero length
	}

	pc := ParsedClass{}
	pc.cpCount = 2
	_, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Errorf("Unexpected error for zero-length UTF-8: %v", err)
	}

	if len(pc.utf8Refs) != 1 {
		t.Error("Expected 1 UTF-8 entry")
	}

	if pc.utf8Refs[0].content != "" {
		t.Errorf("Expected empty string, got: %q", pc.utf8Refs[0].content)
	}
}

func TestParseConstantPool_LongConst_SlotHandling(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// Long constant (takes 2 slots)
	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x05, 0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0, // Long constant
	}

	pc := ParsedClass{}
	pc.cpCount = 3 // Long takes 2 slots + dummy entry
	_, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Errorf("Unexpected error for long constant: %v", err)
	}

	if len(pc.cpIndex) != 3 {
		t.Errorf("Expected 3 CP entries, got %d", len(pc.cpIndex))
	}

	// Check that the second slot is a dummy
	if pc.cpIndex[2].entryType != Dummy {
		t.Error("Expected second slot of long constant to be dummy entry")
	}
}

func TestParseConstantPool_DoubleConst_SlotHandling(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// Double constant (takes 2 slots)
	bytesToTest := []byte{
		0xCA, 0xFE, 0xBA, 0xBE, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x06, 0x40, 0x09, 0x21, 0xFB, 0x54, 0x44, 0x2D, 0x18, // Double constant (π)
	}

	pc := ParsedClass{}
	pc.cpCount = 3 // Double takes 2 slots + dummy entry
	_, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Errorf("Unexpected error for double constant: %v", err)
	}

	if len(pc.cpIndex) != 3 {
		t.Errorf("Expected 3 CP entries, got %d", len(pc.cpIndex))
	}

	// Check that the second slot is a dummy
	if pc.cpIndex[2].entryType != Dummy {
		t.Error("Expected second slot of double constant to be dummy entry")
	}
}

func TestPrintCP_EmptyCP(t *testing.T) {
	// Capture stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	pc := ParsedClass{}
	pc.cpIndex = []cpEntry{{Dummy, 0}} // Just the dummy entry

	printCP(&pc)

	w.Close()
	output, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	outputStr := string(output)
	if !strings.Contains(outputStr, "Number of CP entries parsed: 01") {
		t.Error("Expected output to show 1 CP entry")
	}
	if !strings.Contains(outputStr, "(dummy entry)") {
		t.Error("Expected output to show dummy entry")
	}
}

func TestPrintCP_InvalidEntryType(t *testing.T) {
	// Capture stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	pc := ParsedClass{}
	pc.cpIndex = []cpEntry{
		{Dummy, 0},
		{999, 0}, // Invalid entry type
	}

	printCP(&pc)

	w.Close()
	output, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	outputStr := string(output)
	if !strings.Contains(outputStr, "invalid entry") {
		t.Error("Expected output to show 'invalid entry' for unknown entry type")
	}
}

func TestPrintCP_AllEntryTypes(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()

	// Capture stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	pc := ParsedClass{}
	pc.cpIndex = []cpEntry{
		{Dummy, 0},
		{UTF8, 0},
		{IntConst, 0},
		{FloatConst, 0},
		{LongConst, 0},
		{DoubleConst, 0},
		{ClassRef, 0},
		{StringConst, 0},
		{FieldRef, 0},
		{MethodRef, 0},
		{Interface, 0},
		{NameAndType, 0},
		{MethodHandle, 0},
		{MethodType, 0},
		{Dynamic, 0},
		{InvokeDynamic, 0},
		{Module, 0},
		{Package, 0},
	}

	// Initialize the corresponding data structures
	pc.utf8Refs = []utf8Entry{{"test"}}
	pc.intConsts = []int{42}
	pc.floats = []float32{3.14}
	pc.longConsts = []int64{123456789}
	pc.doubles = []float64{2.718}
	pc.classRefs = []uint32{1}
	pc.stringRefs = []stringConstantEntry{{1}}
	pc.fieldRefs = []fieldRefEntry{{1, 2}}
	pc.methodRefs = []methodRefEntry{{1, 2}}
	pc.interfaceRefs = []interfaceRefEntry{{1, 2}}
	pc.nameAndTypes = []nameAndTypeEntry{{1, 2}}
	pc.methodHandles = []methodHandleEntry{{1, 1}}
	pc.methodTypes = []int{1}
	pc.dynamics = []dynamic{{1, 2}}
	pc.invokeDynamics = []invokeDynamic{{1, 2}}
	pc.moduleName = "testModule"
	pc.packageName = "testPackage"

	printCP(&pc)

	w.Close()
	output, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	outputStr := string(output)

	// Check that all entry types are printed
	expectedStrings := []string{
		"(dummy entry)",
		"(UTF-8 string)",
		"(int constant)",
		"(float constant)",
		"(long constant)",
		"(double constant)",
		"(class ref)",
		"(string const ref)",
		"(field ref)",
		"(method ref)",
		"(interface ref)",
		"(name and type)",
		"(method handle)",
		"(method type)",
		"(dynamic)",
		"(invokedynamic)",
		"(module name)",
		"(package name)",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Expected output to contain: %s", expected)
		}
	}
}

// === end of tests generated by JetBrains Junie ===
