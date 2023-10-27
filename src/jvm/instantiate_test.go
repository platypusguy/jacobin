/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package jvm

import (
	"jacobin/classloader"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/types"
	"testing"
)

// Arrays are preloaded, so this should only confirm the presence of the class
// in the method area--and make sure it has no fields.
func TestInstantiateArray(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)
	classloader.InitMethodArea()

	obj, err := InstantiateClass(types.ByteArray, nil)
	if err != nil {
		t.Errorf("Got unexpected error from instantiating array: %s", err.Error())
	}

	if len(obj.Fields) != 0 {
		t.Errorf("Expected 0 fields in array class, got %d fields", len(obj.Fields))
	}
}
