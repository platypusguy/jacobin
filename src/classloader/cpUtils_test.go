/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"jacobin/globals"
	"jacobin/log"
	"testing"
)

func TestMeInfoFromMethRefInvalid(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.CLASS)

	// set up a class with a constant pool containing the one
	// errors in it
	klass := CPool{}
	klass.CpIndex = append(klass.CpIndex, CpEntry{})
	klass.CpIndex = append(klass.CpIndex, CpEntry{IntConst, 0})
	klass.CpIndex = append(klass.CpIndex, CpEntry{UTF8, 0})

	klass.IntConsts = append(klass.IntConsts, int32(26))
	klass.Utf8Refs = append(klass.Utf8Refs, "Hello string")

	s1, s2, s3 := GetMethInfoFromCPmethref(&klass, 0)
	if s1 != "" && s2 != "" && s3 != "" {
		t.Errorf("Did not get expected result for pointing to CPentry[0]")
	}

	s1, s2, s3 = GetMethInfoFromCPmethref(&klass, 999)
	if s1 != "" && s2 != "" && s3 != "" {
		t.Errorf("Did not get expected result for pointing to CPentry outside of CP")
	}

	s1, s2, s3 = GetMethInfoFromCPmethref(&klass, 1)
	if s1 != "" && s2 != "" && s3 != "" {
		t.Errorf("Did not get expected result for pointing to CPentry that's not a MethodRef")
	}
}
