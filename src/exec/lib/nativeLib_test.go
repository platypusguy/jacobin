/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */
package lib

import (
	"testing"
)

func TestLoadOfMethodSignature(t *testing.T) {
	// load the MethodSignature table
	ms := MethodSignatures
	load()

	// get the entry for println() (Java: System.out.println())
	m := ms["println"]

	// create a slice of interfaces and added the one param (a string to print
	ia := make([]interface{}, m.paramSlots)
	ia[0] = "hello!"

	// call the method body
	m.fu(ia)
}

/*
Implementation notes:
To genericize the function calls, we are defining all functions as taking an array
of empty interfaces and returning nothing. We load the parameters into the array of
empty interfaces (they can be any type) and then call the function whose pointer is
held in the MethodSignatures array.
*/
