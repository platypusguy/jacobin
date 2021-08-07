/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

/****
TODO: Set up a table with all the supported switches:
	name(string), func
		name should remove leading + or - anything after : or =
	should first check for all the version / verbose / help options which just print out info
*/

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// handle all the args from the command line, including those from the enviroment
// variables that the JVM recognizes and prepends to the command-line options
func HandleCli(osArgs []string) (err error) {
	var javaEnvOptions = getEnvArgs()
	Log("Java environment variables: "+javaEnvOptions, FINE)

	// add command-line args to those extracted from the enviroment (if any)
	cliArgs := javaEnvOptions + " "
	for _, v := range osArgs[1:] {
		//		fmt.Printf("\t%q\n", v)
		cliArgs += v + " "
	}
	Global.commandLine = strings.TrimSpace(cliArgs)
	Log("Commandline: "+Global.commandLine, FINE)

	// handle options that request info but don't run the VM, such as:
	// show version, show help, etc.
	discontinue := handleUserMessages(cliArgs) // use cliArgs b/c we want the version with the final space (to ease search)

	// some user messages require a shutdown after message is displayed (see Usage text for examples)
	if discontinue == true {
		return errors.New("end of processing")
	}

	// pull out all the arguments into an array of strings. Note that an arg with spaces but within
	// quotes is treated as a single arg
	args := strings.Fields(javaEnvOptions)
	for _, v := range osArgs[1:] {
		fmt.Printf("\t%q\n", v)
		args = append(args, v)
	}

	for i := 0; i < len(osArgs); i++ {
		opt, ok := Global.options[osArgs[i]]
		if ok {
			opt.f(i, osArgs[i])
		}
	}

	//fmt.Printf("args are: %q\n", args)
	return
}

// you can can set JVM options using the three environment variables that are
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

// handle all the options that simply print messages for the user's benefit
func handleUserMessages(allArgs string) bool {
	const exitProcessing = true

	// the order of the messages is important b/c -showversion and -help can both be requested,
	// but in the reverse order, only -help is shown and then then Jacobin exits
	if strings.Contains(allArgs, "-showversion") {
		showVersion(os.Stderr)
	} else if strings.Contains(allArgs, "--show-version") {
		showVersion(os.Stdout)
	} else if strings.Contains(allArgs, "-version") {
		showVersion(os.Stderr)
		return exitProcessing
	} else if strings.Contains(allArgs, "--version") {
		showVersion(os.Stdout)
		return exitProcessing
	} else {
		showCopyright() // show copyright only if the version information is not requested
	}

	if strings.Contains(allArgs, "-h") || strings.Contains(allArgs, "-help") ||
		strings.Contains(allArgs, "-?") {
		showUsage(os.Stderr)
		return exitProcessing
	} else if strings.Contains(allArgs, "--help") {
		showUsage(os.Stdout)
		return exitProcessing
	}

	return !exitProcessing
}

// show the usage info to the user (in response to errors or java -help and
// similar command-line options). The text will be updated to conform closer
// to the OpenJDK message as features are added to Jacobin
func showUsage(outStream *os.File) {
	userMessage :=
		`
Usage: jacobin [options] <mainclass> [args...]
	        (to execute a class)
   or jacobin [options] -jar <jarfile> [args...]
	        (to execute a jar file)
Arguments following the main class, source file, -jar <jarfile>,
are passed as the arguments to main class.

where options include:
	-? -h -help   print this help message to the error stream
	--help        print this help message to the output stream
	-version      print product version to the error stream and exit
	--version     print product version to the output stream and exit
	-showversion  print product version to the error stream and continue
	--show-version
				  print product version to the output stream and continue`

	fmt.Fprintln(outStream, userMessage)
}

// show the Jacobin version
func showVersion(outStream *os.File) {
	ver := fmt.Sprintf(
		"Jacobin VM v. %s 2021\n64-bit server JVM", Global.version)
	fmt.Fprintln(outStream, ver)
}

func showCopyright() {
	fmt.Println("Jacobin VM v. " + Global.version +
		", © 2021 by Andrew Binstock. All rights reserved. MPL 2.0 License.")
}
