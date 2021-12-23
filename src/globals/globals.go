/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package globals

import (
	"os"
	"strings"
	"sync"
)

// Globals contains variables that need to be globally accessible,
// such as VM and program args, pointers to classloaders, etc.
// Note: globals cannot depend on exec package to avoid circularity.
// As a result, exec contains its own globals
type Globals struct {
	// ---- jacobin version number ----
	// note: all references to version number must come from this literal
	Version string
	VmModel string // "client" or "server" (both the same acc. to JVM docs)

	// ---- processing stoppage? ----
	ExitNow bool

	// ---- command-line items ----
	JacobinName string // name of the executing Jacobin executable
	Args        []string
	CommandLine string

	StartingClass string
	StartingJar   string
	AppArgs       []string
	Options       map[string]Option

	// ---- classloading items ----
	MaxJavaVersion    int // the Java version as commonly known, i.e. Java 11
	MaxJavaVersionRaw int // the Java version as it appears in bytecode i.e., 55 (= Java 11)
	VerifyLevel       int

	// ---- paths for finding the base classes to load ----
	JavaHome    string
	JacobinHome string
}

// Wait group for various channels used for parallel loading of classes.
var LoaderWg sync.WaitGroup

var global Globals

// InitGlobals initializes the global values that are known at start-up
func InitGlobals(progName string) Globals {
	global = Globals{
		Version:           "0.1.0",
		VmModel:           "server",
		ExitNow:           false,
		JacobinName:       progName,
		JacobinHome:       "",
		JavaHome:          "",
		Options:           make(map[string]Option),
		StartingClass:     "",
		StartingJar:       "",
		MaxJavaVersion:    11, // this value and MaxJavaVersionRaw must *always* be in sync
		MaxJavaVersionRaw: 55, // this value and MaxJavaVersion must *always* be in sync
	}
	initJavaHome()
	initJacobinHome()
	return global
}

// GetInstance enables the singleton aspect of Globals
func GetInstance() Globals {
	return global
}

func GetGlobalRef() *Globals {
	return &global
}

// Option is the value portion of the globals.options table. This table is described in
// more detail in option_table_loader.go introductory comments
type Option struct {
	Supported bool
	Set       bool
	ArgStyle  int16
	Action    func(position int, name string, gl *Globals) (int, error)
}

// get JACOBIN_HOME and format it as expected
func initJacobinHome() {
	jacobinHome := os.Getenv("JACOBIN_HOME")
	if jacobinHome != "" {
		// if the JacobinHome doesn't end in a backward slash, add one.
		if !(strings.HasSuffix(jacobinHome, "\\") ||
			strings.HasSuffix(jacobinHome, "/")) {
			jacobinHome = jacobinHome + "\\"
		}
		// replace forward slashes in JacobinHome with backward slashes
		jacobinHome = strings.ReplaceAll(jacobinHome, "/", "\\")
	}
	global.JacobinHome = jacobinHome
}
func JacobinHome() string { return global.JacobinHome }

// get JAVA_HOME and format it as expected
func initJavaHome() {
	javaHome := os.Getenv("JAVA_HOME")
	if javaHome != "" {
		// if the JacobinHome doesn't end in a backward slash, add one.
		if !(strings.HasSuffix(javaHome, "\\") ||
			strings.HasSuffix(javaHome, "/")) {
			javaHome = javaHome + "\\"
		}
		// replace forward slashes in JacobinHome with backward slashes
		javaHome = strings.ReplaceAll(javaHome, "/", "\\")
	}
	global.JavaHome = javaHome
}
func JavaHome() string { return global.JavaHome }
