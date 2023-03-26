/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"errors"
	"fmt"
	"jacobin/execdata"
	"jacobin/globals"
	"jacobin/log"
	"os"
	"strings"
)

// HandleCli handles all args from the command line, including those from environment
// variables that the JVM recognizes and prepends to the list of command-line options
// func HandleCli(osArgs []string, Global *globals.Globals) (err error) {
func HandleCli(osArgs []string, Global *globals.Globals) (err error) {
	var javaEnvOptions = getEnvArgs()
	_ = log.Log("Java environment variables: "+javaEnvOptions, log.FINE)

	// JAVA_HOME and JACOBIN_HOME were obtained in the init of globals.go. Here we just log them.
	showJavaHomeArgs(Global)

	// add command-line args to those extracted from the environment (if any)
	cliArgs := javaEnvOptions + " "
	for _, v := range osArgs[1:] {
		//		fmt.Printf("\t%q\n", v)
		cliArgs += v + " "
	}
	Global.CommandLine = strings.TrimSpace(cliArgs)
	_ = log.Log("Commandline: "+Global.CommandLine, log.FINE)

	// pull out all the arguments into an array of strings. Note that an arg with spaces but
	// within quotes is treated as a single arg
	args := strings.Fields(javaEnvOptions)
	for _, v := range osArgs[1:] {
		//		fmt.Printf("\t%q\n", v)
		args = append(args, v)
	}
	Global.Args = args
	showCopyright(Global)

	for i := 0; i < len(args); i++ {
		var option, arg string
		// if it's a JVM option (so, it begins with a hyphen)
		// break the option into the option and any embedded arg values, if any
		if strings.HasPrefix(args[i], "-") {
			option, arg, err = getOptionRootAndArgs(args[i])
		} else {
			option = args[i]
		}

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
			_, _ = fmt.Fprintf(os.Stderr, "%s is not a recognized option. Ignored.\n", args[i])
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

// log the two environmental variables from which we'll load base classes, if log level allows.
func showJavaHomeArgs(Global *globals.Globals) {
	if Global.JavaHome != "" {
		_ = log.Log("JAVA_HOME: "+Global.JavaHome, log.FINE)
	}
	if Global.JacobinHome != "" {
		_ = log.Log("JACOBIN_HOME: "+Global.JacobinHome, log.FINE)
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
				  print product version to the output stream and continue

Jacobin-specific options:
	-strictJDK    make user messages conform closely to the JDK's format
	-trace:inst   display instruction-level tracing data to the console`

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
			"© 2021-2 by the Jacobin authors. MPL 2.0 License. www.jacobin.org")
	}
}
