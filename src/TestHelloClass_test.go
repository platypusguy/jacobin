/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package main

import (
	"jacobin/classloader"
	"jacobin/globals"
	"jacobin/log"
	"testing"
)

func TestHelloClassFromByteArray(t *testing.T) {

	fileBytes := []byte{0xCA, 0xFE}

	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.CLASS)
	err := classloader.Init()
	// CURR: Resume here.
	// TODO: CFE does not halt execution--we get an error message about no main()

	// classloader.ParseAndPostClass(classloader.BootstrapCL, "Hello", fileBytes)
	// err := StartExec("Hello", globals.GetGlobalRef())
	if err != nil {
		t.Errorf("Error in processing main: %s", error.Error(err)+" "+string(fileBytes)) // added fileBytes here to quiet compiler
	}
}
