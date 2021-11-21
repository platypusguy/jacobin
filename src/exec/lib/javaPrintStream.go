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

var MethodSignatures = make(map[string]method)

type method struct {
	paramSlots int
	fu         function
}

type function func([]interface{})

func load() {
	MethodSignatures["println"] = method{
		paramSlots: 1,
		fu:         Println,
	}
}

// a temporary stand-in for java\io\PrintStream
type stream *os.File

var Out stream

func PrintStream(out stream) {
	Out = out
}

func init() {
	Out = os.Stdout
}

func Println(i []interface{}) {
	fmt.Fprintln(os.Stderr, i[0])
}
