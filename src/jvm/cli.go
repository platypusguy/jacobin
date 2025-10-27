/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-5 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"errors"
	"fmt"
	"jacobin/src/execdata"
	"jacobin/src/globals"
	"jacobin/src/trace"
	"os"
	"strings"
)

// HandleCli handles all args from the command line, including those from environment
// variables that the JVM recognizes and prepends to the list of command-line options
func HandleCli(osArgs []string, Global *globals.Globals) (err error) {

	// Get Java OPTIONS from the environment.
	var javaEnvOptions = getEnvArgs()
	if globals.TraceInit {
		trace.Trace("HandleCli: Java environment variables: " + javaEnvOptions)
	}

	// JAVA_HOME and JACOBIN_HOME were obtained in the init of globals.go. Here we just log them.
	showJavaHomeArgs(Global)

	// add command-line args to those extracted from the environment (if any)
	// Store in Global.CommandLine as a scalar string.
	cliArgs := javaEnvOptions + " "
	for _, v := range osArgs[1:] {
		cliArgs += v + " "
	}
	Global.CommandLine = strings.TrimSpace(cliArgs)
	if globals.TraceInit {
		trace.Trace("HandleCli: Commandline: " + Global.CommandLine)
	}

	// pull out all the arguments into an array of strings. Note that an arg with spaces but
	// within quotes is treated as a single arg
	// Store in Global.Args as a string array.
	args := strings.Fields(javaEnvOptions)
	for _, v := range osArgs[1:] {
		args = append(args, v)
	}
	Global.Args = args

	// Make the lawyers happy.
	showCopyright(Global)

	// Begin main loop.
	// For each args element .....
	for i := 0; i < len(args); i++ {

		// Options look like one of these:
		//		-label			start of a classpath sequence or a flag
		//		-label:value	where value can be a scalar or a list of subvalues
		var optLabel, optValue string
		var dashed bool
		// if it's a JVM option (it begins with a hyphen),
		// 		break the option into the optLabel and any embedded arg values, if any, into optValue
		// else,
		//		just capture the string value in optLabel.
		if strings.HasPrefix(args[i], "-") {
			optLabel, optValue, err = getOptionRootAndArgs(args[i])
			dashed = true
		} else {
			optLabel = args[i]
			dashed = false
		}

		if err != nil {
			errMsg := fmt.Sprintf("HandleCli: getOptionRootAndArgs detected an error in %s, err: %v", args[i], err)
			trace.Error(errMsg)
			return err
		}

		// if the option is the name of the class to execute,
		// * get all successive arguments
		// * store them in the Global.AppArgs array
		// * break out of the outer for loop
		if !dashed {
			Global.StartingClass = optLabel
			for i = i + 1; i < len(args); i++ {
				Global.AppArgs = append(Global.AppArgs, args[i])
			}
			if !strings.HasSuffix(optLabel, ".class") {
				optLabel += ".class"
			}
			break
		}

		// Get the option value for this label.
		opt, ok := Global.Options[optLabel]
		if !ok {
			errMsg := fmt.Sprintf("HandleCli: Parameter %s is not a recognized option. Exiting.\n", args[i])
			trace.Error(errMsg)
			return err
		}

		// Process the option value with the action function.
		newPos, err := opt.Action(i, optValue, Global)
		if err != nil {
			errMsg := fmt.Sprintf("HandleCli: Parameter %s has errors, err: %v\n", args[i], err)
			trace.Error(errMsg)
			return err
		}

		// if the option is a JAR file, then
		// * get all successive arguments
		// * store them in the Global.AppArgs array
		// * break out of the outer for loop
		if optLabel == "-jar" {
			for i = i + 1; i < len(args); i++ {
				Global.AppArgs = append(Global.AppArgs, args[i])
			}
			break
		}
		i = newPos // advance the index by the number of args consumed by this option
	}

	// Finished with args array.
	return nil
}

// pass in the option potentially with embedded arguments and get back
// the option name and the embedded argument(s) as a single string, if any
//
// Return, 3 patterns:
// (1) Pattern is -key:value
// * 	option name (key) - string (E.g. "-cp")
// * 	option argument(s) - string (E.g. ".;C:\home\user\classes")
// * 	error struct - nil (indicates success)
// (2) Pattern is -key
// * 	option name (key) - string (E.g. "--help")
// * 	option argument(s) - ""
// * 	error struct - nil (indicates success)
// (3) Error
// * 	option name (key) - ""
// * 	option argument(s) - ""
// * 	error struct - !nil (indicates failure)
func getOptionRootAndArgs(option string) (string, string, error) {
	if len(option) == 0 {
		return "", "", errors.New("empty option error")
	}

	// if the option has an embedded arg value, it'll come after the first colon (:).
	argMarker := strings.Index(option, ":")

	// if there's no embedded colon (:), then the option doesn't contain an arg value
	if argMarker == -1 {
		return option, "", nil
	}

	return option[:argMarker], option[argMarker+1:], nil

}

// you can set JVM options using the three environment variables that are
// inspected in this function. Note: order is important because later options
// can override earlier ones. These are checked before any of the command-line
// options are processed.
func getEnvArgs() string {
	envArgs := ""
	javaEnvKeys := [3]string{"JAVA_TOOL_OPTIONS", "_JAVA_OPTIONS", "JDK_JAVA_OPTIONS"}

	for i := 0; i < 3; i++ { // if a string is found copy it and a trailing space
		envString := os.Getenv(javaEnvKeys[i])
		if len(envString) > 0 {
			envArgs += envString
			if !strings.HasSuffix(envArgs, " ") {
				envArgs += " "
			}
		}
	}
	return strings.TrimSpace(envArgs)
}

// log the two environmental variables from which we'll load base classes.
func showJavaHomeArgs(Global *globals.Globals) {
	if globals.TraceVerbose {
		if Global.JavaHome != "" {
			trace.Trace("JAVA_HOME: " + Global.JavaHome)
		} else {
			trace.Trace("JAVA_HOME: nil")
		}
		if Global.JacobinHome != "" {
			trace.Trace("JACOBIN_HOME: " + Global.JacobinHome)
		} else {
			trace.Trace("JACOBIN_HOME: nil")
		}
	}
}

// show the usage info to the user (in response to errors or java -help and
// similar command-line options). The text will be updated to conform closer
// to the OpenJDK message as features are added to Jacobin
func ShowUsage(outStream *os.File) {
	userMessage :=
		`
Usage: jacobin [options] <mainclass> [args...]
	        (to execute a class)
   or jacobin [options] -jar <jarfile> [args...]
	        (to execute a jar file)
Arguments following the main class, source file, -jar <jarfile>,
are passed as the arguments to main class.

where options include:
	-client         to select the "client" VM
	-? -h -help     print this help message to the error stream
	--help          print this help message to the output stream
	-version        print product version to the error stream and exit
	--version       print product version to the output stream and exit
	-showversion    print product version to the error stream and continue
	--show-version  print product version to the output stream and continue

Jacobin-specific options:
    -strictJDK            make user messages conform closely to the JDK's format
    -trace=<selections>   display selected tracing to the console
                          where the <selections> are one or more of the following separated by commas (,):
                          * init - process initilization
                          * cloadi - classloader initialization
                          * inst - bytecode interpreter trace
                          * class - class & method support for the interpreter
                          * verbose - inst, class, and more details of the interpreter
    -JJ:galt              Do not use this unless you are a Jacobin developer! `

	_, _ = fmt.Fprintln(outStream, userMessage)
}

// show the Jacobin version and minor associated data
func showVersion(outStream *os.File, global *globals.Globals) {
	// get the build date of the presently executing Jacobin executable
	exeDate := ""
	file, err := os.Stat(global.JacobinName)
	if err == nil {
		date := file.ModTime()
		exeDate = fmt.Sprintf("%d-%02d-%02d", date.Year(), date.Month(), date.Day())
	}

	ver := fmt.Sprintf(
		"Jacobin VM v. %s (Java %d) %s\n64-bit %s VM", global.Version, global.MaxJavaVersion, exeDate, global.VmModel)
	_, _ = fmt.Fprintln(outStream, ver)

	if !strings.Contains(global.CommandLine, "-strictJDK") {
		execdata.GetExecBuildInfo(global)
		vcsHash, exists := global.JacobinBuildData["vcs.revision"]
		if !exists {
			vcsHash = "n/a"
		}

		vcsDate, exists := global.JacobinBuildData["vcs.time"]
		if !exists {
			vcsDate = "n/a"
		}

		_, _ = fmt.Fprintf(outStream, "source: %s, dated %s\n",
			vcsHash, vcsDate)
	}
}

// show the copyright. This appears only in the -version family of options, and
// then only when -strictJDK is off.
func showCopyright(g *globals.Globals) {
	if !strings.Contains(g.CommandLine, "-strictJDK") &&
		(strings.Contains(g.CommandLine, "-showversion") ||
			strings.Contains(g.CommandLine, "--show-version") ||
			strings.Contains(g.CommandLine, "-version") ||
			strings.Contains(g.CommandLine, "--version")) {
		fmt.Println("Jacobin VM, Copyright " +
			"© 2021-5 by the Jacobin authors. MPL 2.0 License. www.jacobin.org")
	}
}
