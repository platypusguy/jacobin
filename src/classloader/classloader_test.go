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

// Most of the functionality in classloader package is tested in other files, such as
// * cpParser_test.go (constant pool parser)
// * formatCheck_test.go (the format checking)
// * parser_test.go (the class parsing)
// etc.
// This files tests remaining routines.

func TestInitOfClassloaders(t *testing.T) {
	_ = Init()

	// check that the classloader hierarchy is set up correctly
	if BootstrapCL.Parent != "" {
		t.Errorf("Expecting parent of Boostrap classloader to be empty, got: %s",
			BootstrapCL.Parent)
	}

	if ExtensionCL.Parent != "bootstrap" {
		t.Errorf("Expecting parent of Extension classloader to be Boostrap, got: %s",
			ExtensionCL.Parent)
	}

	if AppCL.Parent != "extension" {
		t.Errorf("Expecting parent of Application classloader to be Extension, got: %s",
			AppCL.Parent)
	}

	// check that the classloaders have empty tables ready
	if len(BootstrapCL.Classes) != 0 {
		t.Errorf("Expected size of boostrap CL's table to be 0, got: %d", len(BootstrapCL.Classes))
	}

	if len(ExtensionCL.Classes) != 0 {
		t.Errorf("Expected size of extension CL's table to be 0, got: %d", len(ExtensionCL.Classes))
	}

	if len(AppCL.Classes) != 0 {
		t.Errorf("Expected size of application CL's table to be 0, got: %d", len(AppCL.Classes))
	}
}

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
	_ = log.SetLogLevel(log.CLASS)

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

func TestInsertionIntoMethodArea(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.CLASS)

	// redirect stderr & stdout to capture results from stderr
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	k := Klass{}
	k.Status = 'F'
	k.Loader = "application"
	clData := ClData{}
	clData.Name = "WillyWonkaClass"
	k.Data = &clData
	_ = insert("WillyWonkaClass", k)

	// restore stderr and stdout to what they were before
	_ = w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stderr = normalStderr

	msg := string(out[:])

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(msg, "WillyWonkaClass") || !strings.Contains(msg, "application") {
		t.Error("Got unexpected logging message for insertion of Klass into method area: " + msg)
	}

	if len(Classes) != 1 {
		t.Errorf("Expecting method area to have a size of 1, got: %d", len(Classes))
	}
}
