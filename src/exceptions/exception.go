/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package exceptions

import (
	"fmt"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/thread"
)

const (
	ArithmeticDividebyzero = 0
	/*
	   Exception in thread "main" java.lang.ArithmeticException: / by zero
	      at Main.main(Main.java:4)
	*/
)

var literals = []string{
	"Arithmetic Exception, Divide by Zero",
}

// Throw duplicates the exception mechanism in Java. Right now, it displays the
// error message. Will add: catch logic, stack trace, and halt of execution
// TODO: use ThreadNum to find the right thread
func Throw(excType int, clName string, threadNum int, methName string, cp int) {
	thd := globals.GetGlobalRef().Threads.ThreadList.Front().Value.(*thread.ExecThread)
	frameStack := thd.Stack
	f := frames.PeekFrame(frameStack, 0)
	fmt.Println("class name: " + f.ClName)
	msg := fmt.Sprintf(
		"%s%sin %s, in%s, at bytecode[]: %d", literals[excType], ": ", clName, methName, cp)
	_ = log.Log(msg, log.SEVERE)
}
