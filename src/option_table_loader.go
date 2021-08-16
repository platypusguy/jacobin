/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package main

import (
	"errors"
	"fmt"
	"os"
)

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
// Every option Jacobin responds to (even if just to say it's not supported) requires an entry in
// the Option table, except for these options:
// 		-h, -help, --help, and -?
// because these have been handled prior to the use of this table.

func LoadOptionsTable(Global *Globals) {

	client := Option{true, false, 0, clientVM}
	Global.options["-client"] = client
	client.set = true

	dryRun := Option{false, false, 0, notSupported}
	Global.options["--dry-run"] = dryRun
	dryRun.set = true

	help := Option{true, false, 0, showHelpStderrAndExit}
	Global.options["-h"] = help
	Global.options["-help"] = help
	Global.options["-?"] = help

	helpp := Option{true, false, 0, showHelpStdoutAndExit}
	Global.options["--help"] = helpp

	jarFile := Option{true, false, 4, getJarFilename}
	Global.options["-jar"] = jarFile
	jarFile.set = true

	showversion := Option{true, false, 0, showVersionStderr}
	Global.options["-showversion"] = showversion
	showversion.set = true

	show_Version := Option{true, false, 0, showVersionStdout}
	Global.options["--show-version"] = show_Version
	show_Version.set = true

	verboseClass := Option{true, false, 1, verbosityLevel}
	Global.options["-verbose"] = verboseClass
	verboseClass.set = true

	version := Option{true, false, 1, versionStderrThenExit}
	Global.options["-version"] = version
	version.set = true

	vversion := Option{true, false, 1, versionStdoutThenExit}
	Global.options["--version"] = vversion
	vversion.set = true

}

// ---- the functions for the supported CLI options, in alphabetic order ----

// client VM function, simply changes the wording of the version
// info. (This is the same behavior as the OpenJDK JVM.)
func clientVM(pos int, name string) (int, error) { Global.vmModel = "client"; return pos, nil }

// for -jar option. Get the next arg, which must be the JAR filename
func getJarFilename(pos int, name string) (int, error) {
	if len(Global.args) > pos+1 {
		Global.startingJar = Global.args[pos+1]
		Log("Starting with JAR file: "+Global.startingJar, FINE)
		return pos + 1, nil
	} else {
		return pos, os.ErrInvalid
	}
}

// generic notification function that an option is not supported
func notSupported(pos int, name string) (int, error) {
	fmt.Printf("%s is not currently supported in Jacobin\n", name)
	return pos, nil
}

func showHelpStderrAndExit(pos int, name string) (int, error) {
	showUsage(os.Stderr)
	Global.exitNow = true
	return pos, nil
}

func showHelpStdoutAndExit(pos int, name string) (int, error) {
	showUsage(os.Stdout)
	Global.exitNow = true
	return pos, nil
}

func showVersionStderr(pos int, name string) (int, error) {
	showVersion(os.Stderr)
	return pos, nil
}

func showVersionStdout(pos int, name string) (int, error) {
	showVersion(os.Stdout)
	return pos, nil
}

// note that the -version option prints the version then exits the VM
func versionStderrThenExit(pos int, name string) (int, error) {
	showVersion(os.Stderr)
	Global.exitNow = true
	return pos, nil
}

// note that the --version option prints the version info then exits the VM
func versionStdoutThenExit(pos int, name string) (int, error) {
	showVersion(os.Stdout)
	Global.exitNow = true
	return pos, nil
}

// set verbosity level. Note Jacobin starts up at WARNING level, so there is no
// need to set it to that level. You cannot set the level to coarser than WARNING
// which is why there is no way to set the verbosity to SEVERE only.
func verbosityLevel(pos int, argValue string) (int, error) {
	switch argValue {
	case "class":
		Global.logLevel = CLASS
		Log("Logging level set to CLASS", INFO)
	case "info":
		Global.logLevel = INFO
		Log("Logging level set to INFO", INFO)
	case "fine":
		Global.logLevel = FINE
		Log("Logging level set to FINE", INFO)
	case "finest":
		Global.logLevel = FINEST
		Log("Logging level set to FINEST", INFO)
	default:
		Log("Error: "+argValue+" is not a valid verbosity option. Ignored.", WARNING)
		return pos, errors.New("Invalid logging level specified: " + argValue)
	}
	return pos, nil
}
