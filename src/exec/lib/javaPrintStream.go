/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package lib

import (
	"fmt"
	"os"
)

// a temporary stand-in for java\io\PrintStream
type stream *os.File

var Out stream

func PrintStream(out stream) {
	Out = out
}

func init() {
	Out = os.Stdout
}

func Println(s string) {
	if Out == os.Stdout {
		println(s)
	} else if Out == os.Stderr {
		fmt.Fprintln(os.Stderr, s)
	}
}
