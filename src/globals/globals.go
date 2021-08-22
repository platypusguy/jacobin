/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package globals

// Globals contains variables that need to be globally accessible,
// such as VM and program args, pointers to classloaders, etc.
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
	VerifyLevel int
}

// initialize the global values that are known at start-up
// listed in alpha order after the first two items
func InitGlobals(progName string) Globals {
	global := Globals{
		Version:       "0.1.0",
		VmModel:       "server",
		ExitNow:       false,
		JacobinName:   progName,
		Options:       make(map[string]Option),
		StartingClass: "",
		StartingJar:   "",
	}
	return global
}

// the value portion of the globals.ptions table. This is described in more detail in
// option_table_loader.go introductory comments
type Option struct {
	Supported bool
	Set       bool
	ArgStyle  int16
	Action    func(position int, name string) (int, error)
}
