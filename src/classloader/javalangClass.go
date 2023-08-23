/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

// Implementation of some of the functions in in Java/lang/Class. Starting with getPrimitiveInsance()

func getPrimitiveClass(primitive string) Klass {
	if primitive == "int" {
		// do somethng
		return Klass{}
	}
	return Klass{} // should be a void class or the like.
}
