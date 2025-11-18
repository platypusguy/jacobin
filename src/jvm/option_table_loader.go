/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-5 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"fmt"
	"jacobin/src/globals"
	"jacobin/src/statics"
	"jacobin/src/trace"
	"jacobin/src/types"
	"jacobin/src/util"
	"os"
	"path/filepath"
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
//     globals data structure, which contains the Options table. The function returns
//     an int showing the last arg processed, and an error if any.
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
	Global.Options["-enableassertions"] = ea

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

	newThread := globals.Option{true, false, 0, useOldThread}
	Global.Options["-732"] = newThread

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

	// because the -cp and -classpath options override the default classpath as well
	// as the one set in the environment variable CLASSPATH, we need to clear the
	// classpath in the globals structure.
	gl.ClasspathRaw = ""
	gl.Classpath = make([]string, 0) // reset the slice

    // Decide whether the parameter was embedded in the option token itself
    // (e.g., "--class-path=..." or "-cp:...") or supplied as the next arg
    // (e.g., "-cp ..."). We can detect embedded form by inspecting the
    // original option token at args[pos] for '=' or ':'.
    hasEmbedded := false
    if pos >= 0 && pos < len(gl.Args) {
        token := gl.Args[pos]
        hasEmbedded = strings.Contains(token, "=") || strings.Contains(token, ":")
    }

    // If we have an embedded param and a non-empty param string, use it and do not
    // consume the next argument element.
    if hasEmbedded && strings.TrimSpace(param) != "" {
        gl.ClasspathRaw = param
        expandClasspth(gl)
        return pos, nil
    }

    // Otherwise, expect the classpath in the following argument (space-separated form)
    if len(gl.Args) > pos+1 {
        gl.ClasspathRaw = gl.Args[pos+1]
        expandClasspth(gl)  // expand the classpath to its components
        return pos + 1, nil // return pos+1 to indicate that the next arg has been consumed
    }

    // Maintain legacy error text expected by tests
    return pos, fmt.Errorf("missing classpath after -cp or -classpath option")
}

func expandClasspth(gl *globals.Globals) {
	// if the classpath is not set, then set it to the current directory
	if gl.ClasspathRaw == "" {
		gl.ClasspathRaw, _ = os.Getwd()
		gl.Classpath[0] = gl.ClasspathRaw
		checkForPreJDK9(gl)
	}

	// if the classpath is set by env variable or CLI, then split it into its components and expand them
	classpaths := strings.Split(gl.ClasspathRaw, string(os.PathListSeparator))

	jarFiles := make([]string, 0, 10) // for the JAR files, if any, specified in the classpath or via wildcard
	for _, path := range classpaths {
		var entry string
		if strings.HasPrefix(path, `"`) && strings.HasSuffix(path, `"`) {
			entry = path[1 : len(path)-1] // remove the quotes
		}

		if entry == "." { // expand the . to the present working directory
			entry, _ = os.Getwd()
			gl.Classpath = append(gl.Classpath, entry)
			continue
		}

		// expand paths that end with a wildcard
		// (per JVM spec, only the * wildcard is allowed and it must be at end)
		wildcard := string(os.PathSeparator) + "*"
		if strings.HasSuffix(path, wildcard) {
			// if the path ends with a wildcard, then we need to expand it
			// to all files in that directory
			loweriles, _ := filepath.Glob(path + ".jar")
			upperFiles, _ := filepath.Glob(path + ".JAR")
			jarFiles = append(jarFiles, loweriles...)  // add the lower-case jar filenames
			jarFiles = append(jarFiles, upperFiles...) // add the upper-case JAR filenames
			if len(jarFiles) > 0 {
				// if there are JAR files, then add them to the classpath
				gl.Classpath = append(gl.Classpath, jarFiles...)
				continue
			}
		}

		if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".JAR") {
			gl.Classpath = append(gl.Classpath, path)
			continue
		} else if !strings.HasSuffix(path, string(os.PathSeparator)) { // make sure each path ends w/ a path separator
			entry = path + string(os.PathSeparator)
			gl.Classpath = append(gl.Classpath, entry)
			continue
		}
	}

	checkForPreJDK9(gl)
}

// checkForPreJDK9 checks if the JDK version is pre-JDK9 and adds the jar files in the JRE's
// jre/lib/ext directory to the classpath. This option was discontinued in JDK9
func checkForPreJDK9(gl *globals.Globals) {
	if globals.JavaVersion() == "" {
		globals.GetJDKmajorVersion() // if JDKmajorVersion is 0, then set it to the JDK version
	}

	// if JDK is pre-JDK9, then we need to add the JRE lib directory to the classpath
	if globals.GetGlobalRef().JDKmajorVersion != 0 || globals.GetGlobalRef().JDKmajorVersion < 9 {
		jreLibExt := "jre" + string(os.PathSeparator) + "lib" + string(os.PathSeparator) +
			"ext" + string(os.PathSeparator)
		if !strings.HasSuffix(gl.JavaHome, string(os.PathListSeparator)) {
			jreLibExt += string(os.PathListSeparator)
		}
		jreLibExtPath := filepath.Join(gl.JavaHome, jreLibExt) // full path to the JDK's jre/lib/ext directory
		jars, err := util.ListJarFiles(jreLibExtPath)
		if err != nil || len(jars) == 0 {
			return
		} else {
			// add the JRE lib directory to the classpath
			for _, jar := range jars {
				if globals.TraceVerbose {
					trace.Trace("Adding JRE lib jar to classpath: " + jar)
				}
				gl.Classpath = append(gl.Classpath, jar)
			}
			gl.ClasspathRaw = gl.ClasspathRaw + string(os.PathListSeparator) + jreLibExtPath
		}
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

func useOldThread(pos int, name string, gl *globals.Globals) (int, error) {
	gl.UseOldThread = true
	setOptionToSeen("-732", gl)
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

// Marks the given option as having been 'set' that is, specified on the command line
func setOptionToSeen(optionKey string, gl *globals.Globals) {
	o := gl.Options[optionKey]
	o.Set = true
	gl.Options[optionKey] = o
}
