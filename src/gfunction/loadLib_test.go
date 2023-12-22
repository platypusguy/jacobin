/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/classloader"
	"testing"
)

func TestMTableLoadLib(t *testing.T) {
	libMeths := make(map[string]GMeth)
	libMeths["testG1"] = GMeth{ParamSlots: 1, GFunction: nil}
	libMeths["testG2"] = GMeth{ParamSlots: 2, GFunction: nil}
	libMeths["testG3"] = GMeth{ParamSlots: 3, GFunction: nil}
	mtbl := make(classloader.MT)
	loadlib(&mtbl, libMeths)
	if len(mtbl) != 3 {
		t.Errorf("Expecting MTable with 3 entries, got: %d", len(mtbl))
	}
	mte := libMeths["testG2"]
	if mte.ParamSlots != 2 {
		t.Errorf("Expecting MTable entry to have 2 param slots, got: %d",
			mte.ParamSlots)
	}

	if mte.NeedsContext != false {
		t.Errorf("Expecting MTable entry's NeedContext to be false")
	}
}

// test loading of native functions

func TestMTableLoadNatives(t *testing.T) {
	classloader.MTable = make(map[string]classloader.MTentry)
	MTableLoadNatives(&classloader.MTable)
	mte, exists := classloader.MTable["java/lang/Object.<init>()V"]
	if !exists {
		t.Errorf("Expecting MTable entry for java/lang/Object.<init>()V, but it does not exist")
	}

	if mte.MType != 'G' {
		t.Errorf("Expecting java/lang/Object.<init>()V to be of type 'G', but got type: %c",
			mte.MType)
	}
}

// make sure that JustReturn in fact does nothing
func TestJustReturn(t *testing.T) {
	retVal := justReturn(nil)
	if retVal != nil {
		t.Errorf("Expecting nil return value, got: %v", retVal)
	}
}
