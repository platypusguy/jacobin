/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */
package classloader

import (
	"jacobin/globals"
	"jacobin/log"
	"os"
	"testing"
	"time"
)

func checkMap(t *testing.T, key string, expectedJmod string) {

	jmod := CJMapFetch(key)
	if len(jmod) < 1 {
		t.Errorf("checkMap: Nil jmod returned with key={%s}", key)
		return
	}
	if jmod != expectedJmod {
		t.Errorf("checkMap: Expected jmod={%s} but observed jmod={%s}", expectedJmod, jmod)
		return
	}

	t.Logf("checkMap: Key {%s} fetched jmod={%s}\n", key, jmod)

}

func TestJacobinHomeTempdir(t *testing.T) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("cannot get user home dir: %s", err.Error())
		return
	}
	tempDir := homeDir + string(os.PathSeparator) + "temp"
	defer os.RemoveAll(tempDir)
	t.Setenv("JACOBIN_HOME", tempDir)
	_ = os.RemoveAll(tempDir) // Make sure that JACOBIN_HOME does not yet exist
	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.FINEST)

	tStart := time.Now()
	CJMapInit()
	tStop := time.Now()
	elapsed := tStop.Sub(tStart)
	t.Logf("CJMapInit finished in %s seconds\n", elapsed.Round(time.Second).String())

	mapSize := CJMapSize()
	if mapSize < 1 {
		t.Errorf("CJMAP size < 1 (cjmap error)")
		return
	}
	t.Logf("Map size is %d\n", mapSize)

	if CJMapFoundGob() {
		t.Errorf("Expected gob not found but one was found")
	} else {
		t.Logf("Gob not found as expected")
	}

	checkMap(t, "java/lang/String.class", "java.base.jmod")
	checkMap(t, "com/sun/accessibility/internal/resources/accessibility.class", "java.desktop.jmod")

}

func TestJacobinHomeDefault(t *testing.T) {

	saved := os.Getenv("JACOBIN_HOME")
	if saved != "" {
		defer os.Setenv("JACOBIN_HOME", saved)
		os.Unsetenv("JACOBIN_HOME")
	}
	globals.InitGlobals("test")
	log.Init()
	CJMapInit() // Create gob file if it does not yet exist.

	globals.InitGlobals("test")
	log.Init()
	CJMapInit() // Process a pre-existing gob file.
	_ = log.SetLogLevel(log.FINEST)

	if !CJMapFoundGob() {
		t.Errorf("Expected gob found but one was not found")
	} else {
		t.Logf("Gob found as expected")
	}

	mapSize := CJMapSize()
	if mapSize < 1 {
		t.Errorf("CJMAP size < 1 (cjmap error)")
		return
	}
	t.Logf("Map size is %d\n", mapSize)

	checkMap(t, "java/lang/String.class", "java.base.jmod")
	checkMap(t, "com/sun/accessibility/internal/resources/accessibility.class", "java.desktop.jmod")

}
