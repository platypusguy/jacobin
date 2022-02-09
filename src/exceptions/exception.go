/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package exceptions

import (
	"fmt"
	"jacobin/log"
)

const (
	Arithmetic_DivideByZero = 0
)

var literals = []string{
	"Arithmetic Exception, Divide by Zero",
}

// Throw duplicates the exception mechanism in Java. Right now, it displays the
// error message. Will add: catch logic, stack trace, and halt of execution
func Throw(excType int, clName string, methName string, cp int) {
	msg := fmt.Sprintf(
		"%s%sin %s, in%s, at bytecode[]: %d", literals[excType], ": ", clName, methName, cp)
	_ = log.Log(msg, log.SEVERE)
}
