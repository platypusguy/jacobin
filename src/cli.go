/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package main

import (
	"os"
	"strings"
)

func HandleCli() (err error) {
	var javaEnvOptions = getEnvArgs()
	Log("Java environment variables: "+javaEnvOptions, FINEST)
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
