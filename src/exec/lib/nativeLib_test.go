/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */
package lib

import (
	"testing"
)

func TestLoadOfMethodSignature(t *testing.T) {
	ms := MethodSignatures
	load()
	m := ms["println"]
	ia := make([]interface{}, m.paramSlots)
	ia[0] = "hello!"
	m.fu(ia)
}

/*
Implementation notes:
To genericize the function calls, we are defining all functions as taking an array
of empty interfaces and returning nothing. We load the parameters into the array of
empty interfaces (they can be any type) and then call the function whose pointer is
held in the MethodSignatures array.
*/
