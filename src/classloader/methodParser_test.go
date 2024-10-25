/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"jacobin/globals"
	"jacobin/stringPool"
	"jacobin/trace"
	"strconv"
	"testing"
)

// test a valid Code attribute of a method
func TestValidCodeMethodAttribute(t *testing.T) {
	globals.InitGlobals("test")

	// variables we'll need.
	klass := ParsedClass{}
	klass.cpIndex = append(klass.cpIndex, cpEntry{})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{ClassRef, 0}) // points to classRef below
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 2})

	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"Exceptions"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"testMethod"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"java/io/IOException"}) // not used -- string pool instead

	name := "java/io/IOException"
	nameIndex := stringPool.GetStringIndex(&name)
	klass.classRefs = append(klass.classRefs, nameIndex) // classRef[0]

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
	trace.Init()

	// variables we'll need.
	klass := ParsedClass{}
	klass.cpIndex = append(klass.cpIndex, cpEntry{})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{ClassRef, 0}) // points to classRef below
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 2})

	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"Exceptions"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"testMethod"})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{"java/io/IOException"}) // not used -- string pool instead

	name := "java/io/IOException"
	nameIndex := stringPool.GetStringIndex(&name)

	klass.classRefs = append(klass.classRefs, nameIndex) // classRef[0] points to stringPool entry for "java/io/IOException"

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

	if klass.utf8Refs[2].content != name {
		t.Errorf("Expected %s but observed %s", name, klass.utf8Refs[2].content)
	}

	if len(meth.exceptions) != 1 {
		t.Error("In test of Exceptions method attribute, attribute was not added to method struct")
	}

	me := meth.exceptions[0]
	excName := stringPool.GetStringPointer(me)
	if *excName != name {
		t.Errorf("The wrong value for the UTF8 record on Exceptions method attribute was stored. Got: %s",
			*excName)
	}
}
