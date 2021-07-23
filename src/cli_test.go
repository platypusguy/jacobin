/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package main

import (
	"os"
	"testing"
)

// unset all of the JVM environment variables and make sure
// collecting them results in an empty string
func TestGetJVMenvVariablesWhenAbsent(t *testing.T) {
	os.Unsetenv("JAVA_TOOL_OPTIONS")
	os.Unsetenv("_JAVA_OPTIONS")
	os.Unsetenv("JDK_JAVA_OPTIONS")

	javaEnvVars := getEnvArgs()
	if javaEnvVars != "" {
		t.Error("getting non-existent Java enviroment options failed")
	}
}

// set two of the JVM environment variables and make sure
// they are fetched correctly and a space is inserted between them
func TestGetJVMenvVariablesWhenTwoArePresent(t *testing.T) {
	os.Unsetenv("JAVA_TOOL_OPTIONS")
	os.Setenv("_JAVA_OPTIONS", "Hello,")
	os.Setenv("JDK_JAVA_OPTIONS", "Jacobin!")

	javaEnvVars := getEnvArgs()
	if javaEnvVars != "Hello, Jacobin!" {
		t.Error("getting two set Java enviroment options failed: " + javaEnvVars)
	}

	// clean up the environment
	os.Unsetenv("_JAVA_OPTIONS")
	os.Unsetenv("JDK_JAVA_OPTIONS")
}
