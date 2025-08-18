/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */
package classloader

import (
	"jacobin/src/globals"
	"jacobin/src/trace"
	"jacobin/src/types"
	"os"
	"testing"
	"time"
)

func checkMap(t *testing.T, key string, expectedJmod string) {

	jmod := JmodMapFetch(key)
	if len(jmod) < 1 {
		t.Errorf("checkClass: Nil jmod returned with key={%s}", key)
		return
	}
	if jmod != expectedJmod {
		t.Errorf("checkClass: Expected jmod={%s} but observed jmod={%s}", expectedJmod, jmod)
		return
	}

	t.Logf("checkClass: Key {%s} fetched jmod={%s}\n", key, jmod)

}

func TestJmodMapHomeTempdir(t *testing.T) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("cannot get user home dir: %s", err.Error())
		return
	}
	tempDir := homeDir + string(os.PathSeparator) + "temp"
	_ = os.RemoveAll(tempDir) // in case that it pre-exists
	t.Setenv("JACOBIN_HOME", tempDir)
	globals.InitGlobals("test")
	trace.Init()

	tStart := time.Now()
	JmodMapInit()
	tStop := time.Now()
	elapsed := tStop.Sub(tStart)
	t.Logf("JmodMapInit finished in %s seconds\n", elapsed.Round(time.Second).String())

	mapSize := JmodMapSize()
	if mapSize < 1 {
		t.Errorf("map size < 1 (fatal error)")
		return
	}
	t.Logf("Map size is %d\n", mapSize)

	if JmodMapFoundGob() {
		t.Errorf("Expected gob not found but one was found")
	} else {
		t.Logf("Gob not found as expected")
	}

	checkMap(t, types.StringClassName, "java.base.jmod")
	checkMap(t, "com/sun/accessibility/internal/resources/accessibility", "java.desktop.jmod")

}

func TestJmodMapHomeDefault(t *testing.T) {

	saved := os.Getenv("JACOBIN_HOME")
	if saved != "" {
		defer os.Setenv("JACOBIN_HOME", saved)
		_ = os.Unsetenv("JACOBIN_HOME")
	}
	globals.InitGlobals("test")
	trace.Init()
	JmodMapInit() // Create gob file if it does not yet exist.

	globals.InitGlobals("test")
	trace.Init()
	JmodMapInit() // Process a pre-existing gob file.

	if !JmodMapFoundGob() {
		t.Errorf("Expected gob found but one was not found")
	} else {
		t.Logf("Gob found as expected")
	}

	mapSize := JmodMapSize()
	if mapSize < 1 {
		t.Errorf("map size < 1 (fatal error)")
		return
	}
	t.Logf("Map size is %d\n", mapSize)

	checkMap(t, types.StringClassName, "java.base.jmod")
	checkMap(t, "java/util/zip/ZipUtils", "java.base.jmod")
	checkMap(t, "com/sun/accessibility/internal/resources/accessibility", "java.desktop.jmod")

}
