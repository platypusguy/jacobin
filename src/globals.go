/* Jacobin VM -- A Java virtual machine
 * (c) Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0
 */

package main

import (
	"fmt"
	"time"
)

// Globals contains variables that need to be globally accessible,
// such as VM and program args, pointers to classloaders, etc.
type Globals struct {
	// ---- jacobin version number ----
	// note: all references to version number must come from this literal
	version string

	// ---- logging items ----
	logLevel  int
	startTime time.Time

	// ---- command-line items ----
	jacobinName string
	args        []string
	commandLine string

	// ---- classloading items ----
	/*
		var bootstrapLoader = Classloader( name: "bootstrap", parent: "" )
		var systemLoader    = Classloader( name: "system", parent: "bootstrap" )
		var assertionStatus = true //default assertion status is that assertions are executed. This is only for start-up.
		var verifyBytecode  = verifyLevel.remote
	*/
	// ---- command-line items ----
	// commandLine: String = ""
	// var startingClass = ""
	// var appArgs: [String] = [""]
	options map[string]Option

	// ---- classloading items ----
	/*
	   // 0 = no verification, 1=remote (non-bootloader classes), 2=all classes
	   enum verifyLevel : Int { case none = 0, remote = 1, all = 2 }

	*/

	// ---- Command-line options ----
	/*
		Possibly set up a table with a key string: option name
		 and the value being a struct containing:

		 boolean: supported?
		 boolean: set?
		 int16: arguments it takes: 0 = none, 1 = value follows an :, 2= value follows an =,
		                            4= value follows a space, 8= value has multiple ;-separated values
		 function: processing routine (passing in the index to the arg)
	*/
}

// initialize the global values that are known at start-up
func initGlobals(progName string) *Globals {
	globals := new(Globals)
	globals.startTime = time.Now()
	globals.jacobinName = progName
	globals.version = "0.1.0"
	globals.logLevel = WARNING

	globals.options = make(map[string]Option)
	dryRun := Option{false, false, 0, notSupported}
	globals.options["--dry-run"] = dryRun

	return globals
}

type Option struct {
	supported bool
	set       bool
	argStyle  int16
	f         func(position int, name string) error
}

func notSupported(pos int, name string) error {
	fmt.Printf("%s not currently supported in Jacobin\n", name)
	return nil
}
