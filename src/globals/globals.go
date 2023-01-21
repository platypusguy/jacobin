/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package globals

import (
	"container/list"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Globals contains variables that need to be globally accessible,
// such as VM and program args, etc.
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

	// ---- thread management ----
	Threads ThreadList // list of all app execution threads

	// ---- execution context ----
	JacobinBuildData map[string]string

	// ---- special switches ----
	StrictJDK bool // hew closely to actions and error messages of the JDK

	// ---- list of addresses of arrays, see jvm/arrays.go for info ----
	ArrayAddressList *list.List
}

// LoaderWg is a wait group for various channels used for parallel loading of classes.
var LoaderWg sync.WaitGroup

var global Globals

// InitGlobals initializes the global values that are known at start-up
func InitGlobals(progName string) Globals {
	global = Globals{
		Version:           "0.2.1",
		VmModel:           "server",
		ExitNow:           false,
		JacobinName:       progName,
		JacobinHome:       "",
		JavaHome:          "",
		Options:           make(map[string]Option),
		StartingClass:     "",
		StartingJar:       "",
		MaxJavaVersion:    17, // this value and MaxJavaVersionRaw must *always* be in sync
		MaxJavaVersionRaw: 61, // this value and MaxJavaVersion must *always* be in sync
		Threads:           ThreadList{list.New(), sync.Mutex{}},
		JacobinBuildData:  nil,
		StrictJDK:         false,
		ArrayAddressList:  InitArrayAddressList(),
	}

	InitJavaHome()
	InitJacobinHome()
	InitArrayAddressList()
	return global
}

// ThreadList contains a list of all app execution threads and a mutex for adding new threads to the list.
type ThreadList struct {
	ThreadsList  *list.List
	ThreadsMutex sync.Mutex
}

// GetGlobalRef returns a pointer to the singleton instance of Globals
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

// InitJacobinHome gets JACOBIN_HOME and formats it as expected
// Note: any trailing separator is removed from the retrieved string per JACOBIN-184
func InitJacobinHome() {
	jacobinHome := os.Getenv("JACOBIN_HOME")
	if jacobinHome != "" {
		jacobinHome = strings.TrimRight(jacobinHome, "\\/") // remove any trailing separator
		jacobinHome = cleanupPath(jacobinHome)
	}
	global.JacobinHome = jacobinHome
}

func JacobinHome() string { return global.JacobinHome }

// InitJavaHome gets JAVA_HOME from the environment and formats it as expected
// Note: any trailing separator is removed from the retrieved string per JACOBIN-184
func InitJavaHome() {

	javaHome := os.Getenv("JAVA_HOME")
	if javaHome != "" {
		javaHome = strings.TrimRight(javaHome, "\\/") // remove any trailing separator
		javaHome = cleanupPath(javaHome)
	}
	global.JavaHome = javaHome
}

func JavaHome() string { return global.JavaHome }

// Normalize a file path. Slashes are converted to the current platform's path separator if necessary.
func cleanupPath(path string) string {
	path = filepath.FromSlash(path)
	return path
}

// Array addresses must be kept in a list to avoid being GC'd.
// This creates that list.
func InitArrayAddressList() *list.List {
	return list.New()
}
