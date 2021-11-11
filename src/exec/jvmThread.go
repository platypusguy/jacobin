/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package exec

// Creates a JVM program execution thread. These threads are extremely limited.
// They basically hold a stack of frames. They push and pop frames as required.
// They begin execution; they exit when execution ends; and they emit diagnostic
// and performance data.

type execThread struct {
	id    int
	stack []frame
}

type frame struct {
}

func CreateThread() execThread {
	t := execThread{}
	return t
}
