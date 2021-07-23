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
	jto := os.Getenv("JAVA_TOOL_OPTIONS")
	jo := os.Getenv("_JAVA_OPTIONS")
	jjo := os.Getenv("JDK_JAVA_OPTIONS")

	envArgs := ""
	if len(jto) > 0 {
		envArgs += jto
		if !strings.HasSuffix(envArgs, " ") {
			envArgs += " "
		}
	}

	if len(jo) > 0 {
		envArgs += jo
		if !strings.HasSuffix(envArgs, " ") {
			envArgs += " "
		}
	}

	if len(jjo) > 0 {
		envArgs += jjo
	}

	return strings.TrimSpace(envArgs)
}
