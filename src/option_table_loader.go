/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package main

import "fmt"

// This set of routines loads the Global.options table with the various
// JVM command-line options for use later by the CLI processing logic.
//
// The table is initially created in globals.go and its declaration contains a
// key consisting of a string with the option as typed on the command line, and
// a value concisting of an Option struct (also defined in global.go), having
// this layout:
//     type Option struct {
//	        supported bool      // is this option supported in Jacobin?
//	        set       bool      // has this option previously been set on the command line?
//	        argStyle  int16     // what is the format for the argument values to this option?
//                              // 0 = no argument      1 = value follows a :
//                              // 2 = value follows =  4 = value follows a space
//                              // 8 = option has multiple values separated by a ; (such as -cp)
//	        action  func(position int, name string) error  // the action to perform when this option found.
//      }
//
// Every option Jacobin responds to (even if just to say it's not supported) requires an Option
// entry in the Option table.

func LoadOptionsTable(Global *Globals) error {

	dryRun := Option{false, false, 0, notSupported}
	Global.options["--dry-run"] = dryRun

	return nil
}

// generic notification function that an option is not supported
func notSupported(pos int, name string) error {
	fmt.Printf("%s is not currently supported in Jacobin\n", name)
	return nil
}
