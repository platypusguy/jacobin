/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

// This file contains test gfunctions for unit tests. They're primarily designed such
// that you specify the variable types passed in and the return value. They do nothing
// but accept the params and return what the signature promises

func Load_TestGfunctions() {

	// === returning void
	MethodSignatures["java/lang/Object.test(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  arrayCopy,
		}

	MethodSignatures["java/lang/Object.test(D)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  currentTimeMillis,
		}

	MethodSignatures["java/lang/Object.test(Ljava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  nanoTime,
		}

	// === returning int or double

	MethodSignatures["java/lang/Object.test(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  exitI,
		}

	MethodSignatures["java/lang/Object.test(I)D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  forceGC,
		}

	// === accepting reference to java/lang/Object and returning something

	MethodSignatures["java/lang/Object.test(Ljava/lang/Object;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  getProperty,
		}

	MethodSignatures["java/lang/Object.test(Ljava/lang/Object;)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  getConsole,
		}

	MethodSignatures["java/lang/Object.test(Ljava/lang/Object;)D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinit,
		}

}
