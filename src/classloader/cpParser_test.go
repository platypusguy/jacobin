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
		0xCA, 0xFE, 0xBA, 0xBA, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x05, 0x00, 0x00, 0x00, 0x01, // first four bytes of long
		0x00, 0x00, 0x00, 0x02, // second four bytes of long
	}

	pc := parsedClass{}
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
