/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-3 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */
package classloader

import (
	"jacobin/src/globals"
	"jacobin/src/types"
	"testing"
)

func checkClass(t *testing.T, className string, expectedJmod string) bool {

	jmod := JmodMapFetch(className)
	if len(jmod) < 1 {
		t.Errorf("checkClass: Nil jmod returned with className={%s}\n", className)
		return false
	}
	if jmod != expectedJmod {
		t.Errorf("checkClass: Expected jmod = %s, but observed jmod = %s\n", expectedJmod, jmod)
		return false
	}

	t.Logf("checkClass: JmodMapFetch(%s) --> jmod={%s} ok\n", className, jmod)

	classBytes, err := GetClassBytes(expectedJmod, className)
	if err != nil {
		t.Errorf("checkClass: GetClassBytes expectedJmod=%s, className=%s failed\n",
			expectedJmod, className)
		return false
	}
	t.Logf("checkClass: classloader.GetClassBytes returned a byte array for class %s in jmod %s ok\n", className, expectedJmod)

	// Load class from bytes
	_, _, err = loadClassFromBytes(AppCL, className, classBytes)
	if err != nil {
		t.Errorf("checkClass: loadClassFromBytes returned an error: %s\n", error.Error(err))
		return false
	}

	t.Logf("checkClass: Success!\n")

	return true

}

func TestJmodToClass(t *testing.T) {

	// Initialise global and logging
	globals.InitGlobals("test")

	t.Logf("globals.InitGlobals(test) ok\n")

	// Initialise JMODMAP
	JmodMapInit()
	t.Logf("JmodMapInit ok\n")
	mapSize := JmodMapSize()
	if mapSize < 1 {
		t.Errorf("Oh, no! JMODMAP size < 1\n")
		return
	}
	t.Logf("JMODMAP size is %d\n", mapSize)

	// Initialise classloader
	Init()
	t.Logf("classloader.Init ok\n")

	_ = checkClass(t, "com/sun/accessibility/internal/resources/accessibility", "java.desktop.jmod")
	_ = checkClass(t, types.StringClassName, "java.base.jmod")

}
