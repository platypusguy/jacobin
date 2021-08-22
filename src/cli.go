/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package main

import (
	"errors"
	"fmt"
	"jacobin/globals"
	"jacobin/log"
	"os"
	"strings"
)

var global = globals.GetInstance()

// handle all the args from the command line, including those from the enviroment
// variables that the JVM recognizes and prepends to the command-line options
func HandleCli(osArgs []string, Global *globals.Globals) (err error) {
	var javaEnvOptions = getEnvArgs()
	log.Log("Java environment variables: "+javaEnvOptions, log.FINE)

	// add command-line args to those extracted from the enviroment (if any)
	cliArgs := javaEnvOptions + " "
	for _, v := range osArgs[1:] {
		//		fmt.Printf("\t%q\n", v)
		cliArgs += v + " "
	}
	Global.CommandLine = strings.TrimSpace(cliArgs)
	log.Log("Commandline: "+Global.CommandLine, log.FINE)

	// pull out all the arguments into an array of strings. Note that an arg with spaces but
	// within quotes is treated as a single arg
	args := strings.Fields(javaEnvOptions)
	for _, v := range osArgs[1:] {
		//		fmt.Printf("\t%q\n", v)
		args = append(args, v)
	}
	Global.Args = args
	showCopyright()

	for i := 0; i < len(args); i++ {
		// break the option into the option and any embedded arg values
		option, arg, err := getOptionRootAndArgs(args[i])

		if err != nil {
			continue // skip the arg if there was a problem. (Might want to revisit this.)
		}

		// if the option is the name of the class to execute, note that then get
		// all successive arguments and store them as app args in Global
		if strings.HasSuffix(option, ".class") {
			Global.StartingClass = option
			for i = i + 1; i < len(args); i++ {
				Global.AppArgs = append(Global.AppArgs, args[i])
			}
			break
		}

		opt, ok := Global.Options[option]
		if ok {
			i, _ = opt.Action(i, arg, Global)
		} else {
			fmt.Fprintf(os.Stderr, "%s is not a recognized option. Ignored.\n", args[i])
		}

		// TODO: check for JAR specified and process the JAR. At present, it will
		// recognize the JAR file and insert it into Global, and copy all succeeding args
		// to app args. However, it does not recognize the JAR file as an executable.

		// if len(arg) > 0 {
		// 	fmt.Printf("Option %s has argument value: %s\n", option, arg)
		// }
	}
	return nil
}

// pass in the option potentially with embedded arguments and get back
// the option name and the embedded argument(s), if any
func getOptionRootAndArgs(option string) (string, string, error) {
	if len(option) == 0 {
		return "", "", errors.New("empty option error")
	}

	// if the option has an embedded arg value, it'll come after a : or an =
	argMarker := strings.Index(option, ":")
	if argMarker == -1 {
		argMarker = strings.Index(option, "=")
	}

	// if there's no embedded : or = then the option doesn't contain an arg value
	if argMarker == -1 {
		return option, "", nil
	}

	return option[:argMarker], option[argMarker+1:], nil

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
	-client       to select the "client" VM
	-verbose:[class|info|fine|finest]  enable verbose output
                  info, fine, finest are Jacobin-specific options providing
                    increasing amounts of detail. The finest level is used
                    primarily for performance analysis.
	-? -h -help   print this help message to the error stream
	--help        print this help message to the output stream
	-version      print product version to the error stream and exit
	--version     print product version to the output stream and exit
	-showversion  print product version to the error stream and continue
	--show-version
				  print product version to the output stream and continue`

	fmt.Fprintln(outStream, userMessage)
}

// show the Jacobin version and minor associated data
func showVersion(outStream *os.File) {
	// get the build date of the presently executing Jacobin executable
	exeDate := ""
	file, err := os.Stat(global.JacobinName)
	if err == nil {
		date := file.ModTime()
		exeDate = fmt.Sprintf("%d-%02d-%02d", date.Year(), date.Month(), date.Day())
	}

	ver := fmt.Sprintf(
		"Jacobin VM v. %s (Java 11.0.10) %s\n64-bit %s VM", global.Version, exeDate, global.VmModel)
	fmt.Fprintln(outStream, ver)
}

// show the copyright. Because the various -version commands show much the
// same data, rather than printing it twice, we skip showing the copyright
// info when the -version option variants are specified
func showCopyright() {
	global = globals.GetInstance()
	if !strings.Contains(global.CommandLine, "-showversion") &&
		!strings.Contains(global.CommandLine, "--show-version") &&
		!strings.Contains(global.CommandLine, "-version") &&
		!strings.Contains(global.CommandLine, "--version") {
		fmt.Println("Jacobin VM v. " + global.Version +
			", © 2021 by Andrew Binstock. All rights reserved. MPL 2.0 License.")
	}
}
