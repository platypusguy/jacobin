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

func TestCPvalidFieldRef(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{0xCA, 0xFE, 0xBA, 0xBA, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x09, 0x00, 0x14, 0x01, 0x01,
	}

	pc := parsedClass{}
	pc.cpCount = 2
	loc, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Error("Parsing valid CP FieldRef generated an unexpected error")
	}

	if loc != 14 {
		t.Error("Was expecting a new position of 14, but got: " + strconv.Itoa(loc))
	}

	if len(fieldRefs) != 1 {
		t.Error("Was expecting the field ref array to have 1 entry, but it has: " + strconv.Itoa(len(fieldRefs)))
	}

	fre := fieldRefs[0]
	if fre.classIndex != 20 {
		t.Error("Was expecting a field ref classIndex of 20, but got: " + strconv.Itoa(fre.classIndex))
	}

	if fre.nameAndTypeIndex != 257 {
		t.Error("Was expecting a field ref classIndex of 257, but got: " + strconv.Itoa(fre.nameAndTypeIndex))
	}

	if len(cpool) != 2 {
		t.Error("Was expecting cpool to have 2 entries, but instead got: " + strconv.Itoa(len(cpool)))
	}
}

func TestCPvalidMethodRef(t *testing.T) {

	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.WARNING)

	bytesToTest := []byte{0xCA, 0xFE, 0xBA, 0xBA, 0x00,
		0x00, 0xFF, 0xF0, 0x00, 0x00,
		0x0A, 0x00, 0x15, 0x01, 0x06,
	}

	pc := parsedClass{}
	pc.cpCount = 2
	loc, err := parseConstantPool(bytesToTest, &pc)

	if err != nil {
		t.Error("Parsing valid CP FieldRef generated an unexpected error")
	}

	if loc != 14 {
		t.Error("Was expecting a new position of 14, but got: " + strconv.Itoa(loc))
	}

	if len(methodRefs) != 1 {
		t.Error("Was expecting the method ref array to have 1 entry, but it has: " + strconv.Itoa(len(methodRefs)))
	}

	mre := methodRefs[0]
	if mre.classIndex != 21 {
		t.Error("Was expecting a method ref classIndex of 21, but got: " + strconv.Itoa(mre.classIndex))
	}

	if mre.nameAndTypeIndex != 262 {
		t.Error("Was expecting a method ref classIndex of 262, but got: " + strconv.Itoa(mre.nameAndTypeIndex))
	}

	if len(cpool) != 2 {
		t.Error("Was expecting cpool to have 2 entries, but instead got: " + strconv.Itoa(len(cpool)))
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
		t.Error("Parsing valid CP StringConstRef generated an unexpected error")
	}

	if loc != 12 {
		t.Error("Was expecting a new position of 12, but got: " + strconv.Itoa(loc))
	}

	if len(stringRefs) != 1 {
		t.Error("Was expecting the string const ref array to have 1 entry, but it has: " + strconv.Itoa(len(methodRefs)))
	}

	sre := stringRefs[0]
	if sre.index != 32 {
		t.Error("Was expecting a method ref classIndex of 32, but got: " + strconv.Itoa(sre.index))
	}

	if len(cpool) != 2 {
		t.Error("Was expecting cpool to have 2 entries, but instead got: " + strconv.Itoa(len(cpool)))
	}
}
