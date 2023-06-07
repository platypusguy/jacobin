/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package globals

import (
	"bufio"
	"container/list"
	"fmt"
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

	// ---- Java Home and Version ----
	JavaHome    string
	JavaVersion string

	// ---- Jacobin Home ----
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
		JavaVersion:       "",
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
	if global.JavaHome == "" || global.JavaVersion == "" {
		os.Exit(1)
	}
	InitJacobinHome()
	if global.JacobinHome == "" {
		os.Exit(1)
	}
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
	} else {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "InitJacobinHome: os.UserHomeDir() failed. Exiting.\n")
			_, _ = fmt.Fprintf(os.Stderr, err.Error()+"\n")
			return
		}
		jacobinHome = userHomeDir + string(os.PathSeparator) + "jacobin"
	}
	// 0755 (Unix octal): user(owner) can do anything, group and other can read and visit directory ("execute").
	// Ref: https://opensource.com/article/19/8/linux-permissions-101
	err := os.MkdirAll(jacobinHome, 0755)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "InitJacobinHome: os.MkDirAll(%s) failed. Exiting.\n", jacobinHome)
		_, _ = fmt.Fprintf(os.Stderr, err.Error()+"\n")
		return
	}

	// Success!
	global.JacobinHome = jacobinHome
}

func JacobinHome() string { return global.JacobinHome }

// InitJavaHome gets JAVA_HOME from the environment and formats it as expected
// Note: any trailing separator is removed from the retrieved string per JACOBIN-184
func InitJavaHome() {

	javaHome := os.Getenv("JAVA_HOME")
	if javaHome == "" {
		_, _ = fmt.Fprintf(os.Stderr, "InitJavaHome: Environment variable JAVA_HOME missing but is required. Exiting.\n")
		return
	}
	javaHome = strings.TrimRight(javaHome, "\\/") // remove any trailing separator
	javaHome = cleanupPath(javaHome)
	global.JavaHome = javaHome

	releasePath := javaHome + string(os.PathSeparator) + "release"
	handle, err := os.Open(releasePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "InitJavaHome: os.Open(%s) failed. Exiting.\n", releasePath)
		_, _ = fmt.Fprintf(os.Stderr, err.Error()+"\n")
		return
	}
	defer handle.Close()
	scanner := bufio.NewScanner(handle)
	for scanner.Scan() {
		// do something with a line
		line := scanner.Text()
		tokens := strings.Split(line, "=")
		if len(tokens) != 2 {
			_, _ = fmt.Fprintf(os.Stderr, "InitJavaHome: File format error in %s. Exiting.\n", releasePath)
			return
		}
		if tokens[0] == "JAVA_VERSION" {
			global.JavaVersion = strings.Trim(tokens[1], "\"")
			return
		}
	}

	// At this pint, we did not find a Java version record
	_, _ = fmt.Fprintf(os.Stderr, "InitJavaHome: Did not find the JAVA_VERSION record in %s. Exiting.\n", releasePath)
	global.JavaVersion = ""

}

func JavaHome() string    { return global.JavaHome }
func JavaVersion() string { return global.JavaVersion }

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
