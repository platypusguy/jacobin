/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import (
	"io"
	"jacobin/classloader"
	"jacobin/globals"
	"jacobin/opcodes"
	"os"
	"strings"
	"testing"
)

// GETFIELD : get field from object -- here testing for CP not pointing to a field ref
func TestCodeCheckForGetfield(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.GETFIELD)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.MethodRef, Slot: 0} // should be a field ref

	f.CP = &CP

	err := classloader.CheckCodeValidity(f.Meth, f.CP.(*classloader.CPool))
	if err == nil {
		t.Errorf("GETFIELD: Expected error but did not get one.")
	}

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if !strings.Contains(errMsg, "java.lang.VerifyError") || !strings.Contains(errMsg, "not a field reference") {
		t.Errorf("GETFIELD: Did not get expected error message, got: %s", errMsg)
	}
}

// INVOKEVIRTUAL : invoke method -- here testing for error
func TestNewInvokevirtualInvalidMethRef(t *testing.T) {
	globals.InitGlobals("test")

	// redirect stderr so as not to pollute the test output with the expected error message
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.INVOKEVIRTUAL)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0} // should be a method ref
	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		ClName:  "testClass",
		FldName: "testField",
		FldType: "I",
	}

	f.CP = &CP

	err := classloader.CheckCodeValidity(f.Meth, f.CP.(*classloader.CPool))
	if err == nil {
		t.Errorf("INVOKEVIRTUAL: Expected error but did not get one.")
	}

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)

	if errMsg == "" {
		t.Errorf("INVOKEVIRTUAL: Expected error but did not get one.")
	}

	if !strings.Contains(errMsg, "java.lang.VerifyError") || !strings.Contains(errMsg, "not a method reference") {
		t.Errorf("INVOKEVIRTUAL: Did not get expected error message, got:\n %s", errMsg)
	}
}
