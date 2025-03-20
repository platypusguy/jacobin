/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-5 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package util

import (
	"testing"
)

func essai(t *testing.T, arg interface{}) {
	value, err := HashAnything(arg)
	if err != nil {
		t.Errorf("TestHashPrimitives/essai: %v, err: %v", arg, err)
		return
	}
	t.Logf("TestHashPrimitives/essai ok: %v --> %v", arg, value)
}

type tMisc struct {
	ii int64
	ff float64
	bb [8]byte
}

func TestHashAnything(t *testing.T) {
	for ix := 0; ix < 10; ix++ {
		essai(t, ix)
	}

	essai(t, "a")
	essai(t, "ab")
	essai(t, "abc")

	var misc tMisc = tMisc{42, 3.14159265, [8]byte{1, 2, 3, 4, 5, 6, 7, 8}}
	essai(t, misc)
}
