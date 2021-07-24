/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

/****
Set up a table with all the supported switches:
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

func HandleCli() (err error) {
	var javaEnvOptions = getEnvArgs()
	Log("Java environment variables: "+javaEnvOptions, FINEST)
	cliArgs := javaEnvOptions + " "
	for _, v := range os.Args[1:] {
		cliArgs += v + " "
	}
	Global.commandLine = strings.TrimSpace(cliArgs)
	Log("Commandline: "+Global.commandLine, FINE)

	// handle options that request info but don't run the VM, such as:
	// show version, show help, etc.
	discontinue := handleUserMessages(cliArgs) // use cliArgs b/c we want the version with the final space (to ease search)
	if discontinue == true {
		return errors.New("discontinue")
	}
	return
}

// you can can set JVM options using the three environment variables that are
// inspected in this function. Note: order is important because later options
// can override earlier ones. These are checked before any of the command-line
// options are processed.
func getEnvArgs() string {
	envArgs := ""
	javaEnvKeys := [3]string{"JAVA_TOOL_OPTIONS", "_JAVA_OPTIONS", "JDK_JAVA_OPTIONS"}

	for i := 0; i < 3; i++ { // if string is found copy it and a trailing space
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

func handleUserMessages(allArgs string) bool {
	if strings.Contains(allArgs, "-h") || strings.Contains(allArgs, "-help") ||
		strings.Contains(allArgs, "-?") {
		showUsage(os.Stderr)
		return true
	} else if strings.Contains(allArgs, "--help") {
		showUsage(os.Stdout)
		return true
	}

	return false
}

func showUsage(outStream *os.File) {
	userMessage :=
		`Usage: jacobin [options] <mainclass> [args...]
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

/**
          allArgs.contains( "?" ) ||
             allArgs.contains( "-h" ) ||
             allArgs.contains( "-help" ) {
      UserMsgs.showUsage( stream: Streams.serr );
      return execStop
  } else if allArgs.contains( "--help" ) {
      UserMsgs.showUsage( stream: Streams.sout )
      return execStop
  } else if allArgs.contains( "-version" ) {
      UserMsgs.showVersion( stream: Streams.serr )
      return execStop
  } else if allArgs.contains( "--version" ) {
      UserMsgs.showVersion( stream: Streams.sout )
      return execStop
  } else if allArgs.contains( "-showversion" ) {
      UserMsgs.showVersion( stream: Streams.serr )
      return execContinue
  } else if allArgs.contains( "--showversion" ) {
      UserMsgs.showVersion( stream: Streams.sout )
      return execContinue
  }
*/
