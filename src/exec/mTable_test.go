/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package exec

import "testing"

func TestMTableAdd(t *testing.T) {
	mtbl := make(MT)
	addEntry(&mtbl, "test1", MTentry{
		meth:  nil,
		mType: 'G',
	})

	if len(mtbl) != 1 {
		t.Errorf("Expecting MTable size of 1, got: %d", len(mtbl))
	}

	if mtbl["test1"].mType != 'G' {
		t.Errorf("Expecting fetch of a 'G' MTable rec, but got type: %c",
			mtbl["test1"].mType)
	}
}

func TestMTableLoadLib(t *testing.T) {
	libMeths := make(map[string]GMeth)
	libMeths["testG1"] = GMeth{ParamSlots: 1, GFunction: nil}
	libMeths["testG2"] = GMeth{ParamSlots: 2, GFunction: nil}
	libMeths["testG3"] = GMeth{ParamSlots: 3, GFunction: nil}

	mtbl := make(MT)
	loadlib(&mtbl, libMeths)
	if len(mtbl) != 3 {
		t.Errorf("Expecting MTable with 3 entries, got: %d", len(mtbl))
	}

	mte := libMeths["testG2"]
	if mte.ParamSlots != 2 {
		t.Errorf("Expecting MTable entry to have 2 param slots, got: %d",
			mte.ParamSlots)
	}
}
