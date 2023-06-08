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

	if CJMapFoundGob() {
		t.Errorf("Expected gob not found but one was found")
	} else {
		t.Logf("Gob not found as expected")
	}

	key := "java/lang/String.class"
	jmod := CJMapFetch(key)
	mapSize := CJMapSize()
	if mapSize < 1 {
		t.Errorf("CJMAP size < 1 (cjmap error)")
		return
	}
	t.Logf("Map size is %d\n", mapSize)

	elapsed := tStop.Sub(tStart)
	t.Logf("Key %s fetched %s in %s seconds\n", key, jmod, elapsed.Round(time.Second).String())
	if jmod != "java.base.jmod" {
		t.Errorf("Expected jmod=java.base.jmod, observed jmod=%s", jmod)
	} else {
		t.Logf("jmod=java.base.jmod as expected")
	}

	key = "com/sun/accessibility/internal/resources/accessibility.class"
	jmod = CJMapFetch(key)
	mapSize = CJMapSize()
	if mapSize < 1 {
		t.Errorf("CJMAP size < 1 (cjmap error)")
		return
	}
	t.Logf("Key %s fetched %s\n", key, jmod)

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

	key := "java/lang/String.class"
	jmod := CJMapFetch(key)
	mapSize := CJMapSize()
	if mapSize < 1 {
		t.Errorf("CJMAP size < 1 (cjmap error)")
		return
	}
	t.Logf("Key %s fetched %s\n", key, jmod)
	t.Logf("Map size is %d\n", mapSize)

	key = "com/sun/accessibility/internal/resources/accessibility.class"
	jmod = CJMapFetch(key)
	mapSize = CJMapSize()
	if mapSize < 1 {
		t.Errorf("CJMAP size < 1 (cjmap error)")
		return
	}
	t.Logf("Key %s fetched %s\n", key, jmod)

}
