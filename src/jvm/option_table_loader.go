/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-5 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"fmt"
	"jacobin/globals"
	"jacobin/statics"
	"jacobin/trace"
	"jacobin/types"
	"os"
	"strings"
)

// This set of routines loads the globPtr.Options table with the various
// JVM command-line options for use later by the CLI processing logic.
//
// The table is initially created in globals.go and its declaration contains a
// key consisting of a string with the option as typed on the command line, and
// a value consisting of an Option struct (also defined in global.go), having
// this layout:
//     type Option struct {
//	        supported bool      // is this option supported in Jacobin?
//	        set       bool      // has this option previously been set on the command line?
//	        argStyle  int16     // what is the format for the argument values to this option?
//                              // 0 = no argument      1 = value follows a :
//                              // 2 = value follows =  4 = value follows a space
//                              // 10 = option has multiple values separated by a single character (such as in -trace and -cp)
//	        action  func(position int, name string, gl pointer to globasl) error
//                              // which is the action to perform when this option found.
//      }
//
// Every option that Jacobin responds to (even if just to say it's not supported) requires
// an entry in the Option table, except for these options:
// 		-h, -help, --help, and -?
// because these have been handled prior to the use of this table.
//
// ==== How to add new options to Jacobin:
// 1) Create an entry in LoadOptionsTable:
//    * x := globalOptions {
//             where param1 = is a boolean: is the option supported? s/be true
//							  Setting it to false avoids an error message to the
//							  user that the option is unrecognized while still
//							  having it be unsupported
//                   param2 = boolean: has the options been set yet? s/be false
//					 param3 = integer as explained in the previous paragraphs
//                   param3 = the function to perform
//  2) Add x to the GlobalOptions table, using the string of the option as the key
//     Note that in options with parameters after an : or an = (types 1 or 2 in
//     param3 in step 1), you enter only the root as the key. For example, see
//     the -verbose entry below.
//  3) create the function referred to in param 3 in step 1. This function accepts
//     the position in the command line where the present option is located (first
//     option is at position zero), a string which contains any parameters (if it has
//     no parameters an empty string is passed in), and finally a pointer to the
//     globals data structure, which contains the Options table.
//

// LoadOptionsTable loads the table with all the options Jacobin recognizes.
func LoadOptionsTable(Global globals.Globals) {

	classpath := globals.Option{true, false, 4, getClasspath}
	Global.Options["-classpath"] = classpath
	Global.Options["--class-path"] = classpath
	Global.Options["-cp"] = classpath
	classpath.Set = true

	client := globals.Option{true, false, 0, clientVM}
	Global.Options["-client"] = client
	client.Set = true

	// --dry-run option is a valid HotSpot option, but not supported in Jacobin.
	// including it here so that we can test the unsupported option.
	// in Hotpot, it is used to run the VM without actually running the main method.
	dryRun := globals.Option{false, false, 0, notSupported}
	Global.Options["--dry-run"] = dryRun
	dryRun.Set = true

	ea := globals.Option{false, false, 0, enableAssertions}
	Global.Options["-ea"] = ea

	help := globals.Option{true, false, 0, showHelpStderrAndExit}
	Global.Options["-h"] = help
	Global.Options["-help"] = help
	Global.Options["-?"] = help

	helpp := globals.Option{true, false, 0, showHelpStdoutAndExit}
	Global.Options["--help"] = helpp

	jarFile := globals.Option{true, false, 4, getJarFilename}
	Global.Options["-jar"] = jarFile
	jarFile.Set = true

	showversion := globals.Option{true, false, 0, showVersionStderr}
	Global.Options["-showversion"] = showversion

	show_Version := globals.Option{true, false, 0, showVersionStdout}
	Global.Options["--show-version"] = show_Version

	strictJdk := globals.Option{true, false, 0, strictJDK}
	Global.Options["-strictJDK"] = strictJdk

	traceInstruction := globals.Option{true, false, 10, enableTrace}
	Global.Options["-trace"] = traceInstruction

	JJ := globals.Option{true, false, 10, enableJJ}
	Global.Options["-JJ"] = JJ

	version := globals.Option{true, false, 1, versionStderrThenExit}
	Global.Options["-version"] = version

	vversion := globals.Option{true, false, 1, versionStdoutThenExit}
	Global.Options["--version"] = vversion
}

// ---- the functions for the supported CLI options, in alphabetic order ----

// client VM function, simply changes the wording of the version
// info. (This is the same behavior as the OpenJDK JVM.)
func clientVM(pos int, name string, gl *globals.Globals) (int, error) {
	gl.VmModel = "client"
	setOptionToSeen("-client", gl)
	return pos, nil
}

// extracts the classpath from the command line, and break it into it components
func getClasspath(pos int, param string, gl *globals.Globals) (int, error) {
	setOptionToSeen("-cp", gl)
	setOptionToSeen("-classpath", gl)
	setOptionToSeen("--class-path", gl)
	if len(gl.Args) > pos+1 {
		gl.ClasspathRaw = gl.Args[pos+1]
		gl.Classpath = strings.Split(gl.ClasspathRaw, string(os.PathListSeparator))
		return pos + 1, nil
	} else {
		return pos, fmt.Errorf("missing classpath after -cp or -classpath option")
	}
}

// for -jar option. Get the next arg, which must be the JAR filename, and then all remaining args
// are app args, which are duly added to globPtr.appArgs
func getJarFilename(pos int, name string, gl *globals.Globals) (int, error) {
	setOptionToSeen("-jar", gl)
	if len(gl.Args) > pos+1 {
		gl.StartingJar = gl.Args[pos+1]
		if globals.TraceVerbose {
			trace.Trace("Starting with JAR file: " + gl.StartingJar)
		}
		for i := pos + 2; i < len(gl.Args); i++ {
			gl.AppArgs = append(gl.AppArgs, gl.Args[i])
		}
		return len(gl.Args), nil
	} else {
		return pos, os.ErrInvalid
	}
}

// generic notification function that an option is not supported
func notSupported(pos int, arg string, gl *globals.Globals) (int, error) {
	name := gl.Args[pos]
	fmt.Fprintf(os.Stderr, "%s is not currently supported in Jacobin\n", name)
	return pos, nil
}

func showHelpStderrAndExit(pos int, name string, gl *globals.Globals) (int, error) {
	ShowUsage(os.Stderr)
	gl.ExitNow = true
	return pos, nil
}

func showHelpStdoutAndExit(pos int, name string, gl *globals.Globals) (int, error) {
	ShowUsage(os.Stdout)
	gl.ExitNow = true
	return pos, nil
}

func showVersionStderr(pos int, name string, gl *globals.Globals) (int, error) {
	showVersion(os.Stderr, gl)
	setOptionToSeen("-showversion", gl)
	return pos, nil
}

func showVersionStdout(pos int, name string, gl *globals.Globals) (int, error) {
	showVersion(os.Stdout, gl)
	setOptionToSeen("--show-version", gl)
	return pos, nil
}

func strictJDK(pos int, name string, gl *globals.Globals) (int, error) {
	gl.StrictJDK = true
	setOptionToSeen("-strictJDK", gl)
	return pos, nil
}

// note that the -version option prints the version then exits the VM
func versionStderrThenExit(pos int, name string, gl *globals.Globals) (int, error) {
	showVersion(os.Stderr, gl)
	gl.ExitNow = true
	return pos, nil
}

// note that the --version option prints the version info then exits the VM
func versionStdoutThenExit(pos int, name string, gl *globals.Globals) (int, error) {
	showVersion(os.Stdout, gl)
	gl.ExitNow = true
	return pos, nil
}

const TraceSep = ","

func enableTrace(pos int, argValue string, gl *globals.Globals) (int, error) {
	setOptionToSeen("-trace", gl)
	array := strings.Split(argValue, TraceSep)
	for i := 0; i < len(array); i++ {
		switch array[i] {
		case "class":
			globals.TraceClass = true
		case "cloadi":
			globals.TraceCloadi = true
		case "init":
			globals.TraceInit = true
		case "inst":
			globals.TraceInst = true
		case "verbose":
			globals.TraceVerbose = true
			globals.TraceInst = true
		default:
			return 0, fmt.Errorf("unknown -trace option: %s", array[i])
		}
	}
	return pos, nil
}

func enableAssertions(pos int, name string, gl *globals.Globals) (int, error) {
	setOptionToSeen("-ea", gl)
	statics.AddStatic("main.$assertionsDisabled",
		statics.Static{Type: types.Int, Value: types.JavaBoolFalse})
	return pos, nil
}

func enableJJ(pos int, argValue string, gl *globals.Globals) (int, error) {
	setOptionToSeen("-trace", gl)
	array := strings.Split(argValue, TraceSep)
	for i := 0; i < len(array); i++ {
		switch array[i] {
		case "galt":
			globals.Galt = true
		default:
			return 0, fmt.Errorf("unknown -JJ option: %s", array[i])
		}
	}
	return pos, nil
}

// func newInterpeter(pos int, name string, gl *globals.Globals) (int, error) {
// 	gl.NewInterpreter = true
// 	setOptionToSeen("-new", gl)
// 	return pos, nil
// }

// Marks the given option as having been 'set' that is, specified on the command line
func setOptionToSeen(optionKey string, gl *globals.Globals) {
	o := gl.Options[optionKey]
	o.Set = true
	gl.Options[optionKey] = o
}
