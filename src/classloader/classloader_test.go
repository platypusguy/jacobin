/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"jacobin/globals"
	"jacobin/log"
	"testing"
)

// Most of the functionality in classloader package is tested in other files, such as
// * cpParser_test.go (constant pool parser)
// * formatCheck_test.go (the format checking)
// * parser_test.go (the class parsing)
// etc.
// This files tests remaining routines.

// remove leading [L and delete trailing;, eliminate all other entries with [prefix
func TestNormalizingClassReference(t *testing.T) {
	s := normalizeClassReference("[Ljava/test/java.String;")
	if s != "java/test/java.String" {
		t.Error("Unexpected normalized class reference: " + s)
	}

	s = normalizeClassReference("[B")
	if s != "" {
		t.Error("Unexpected normalized class reference: " + s)
	}

	s = normalizeClassReference("java/lang/Object")
	if s != "java/lang/Object" {
		t.Error("Unexpected normalized class reference: " + s)
	}
}

func TestConvertToPostableClassStringRefs(t *testing.T) {
	// Testing the changes made as a result of JACOBIN-103
	globals.InitGlobals("test")
	log.Init()
	log.SetLogLevel(log.CLASS)

	// set up a class with a constant pool containing the one
	// StringConst we want to make sure is converted to a UTF8
	klass := ParsedClass{}
	klass.cpIndex = append(klass.cpIndex, cpEntry{})
	klass.cpIndex = append(klass.cpIndex, cpEntry{StringConst, 0})
	klass.cpIndex = append(klass.cpIndex, cpEntry{UTF8, 0})

	klass.stringRefs = append(klass.stringRefs, stringConstantEntry{index: 0})
	klass.utf8Refs = append(klass.utf8Refs, utf8Entry{content: "Hello string"})

	klass.cpCount = 3

	postableClass := convertToPostableClass(&klass)
	if len(postableClass.CP.Utf8Refs) != 1 {
		t.Errorf("Expecting a UTF8 slice of length 1, got %d",
			len(postableClass.CP.Utf8Refs))
	}

	// cpIndex[1] is a StringConst above, should now be a UTF8
	utf8 := postableClass.CP.CpIndex[1]
	if utf8.Type != UTF8 {
		t.Errorf("Expecting StringConst entry to have become UTF8 entry,"+
			"but instead is of type: %d", utf8.Type)
	}
}
