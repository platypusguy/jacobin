/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"jacobin/shutdown"
	"time"
)

/*
 Each object or library that has Go methods contains a reference to MethodSignatures,
 which contain data needed to insert the go method into the MTable of the currently
 executing JVM. MethodSignatures is a map whose key is the fully qualified name and
 type of the method (that is, the method's full signature) and a value consisting of
 a struct of an int (the number of slots to pop off the caller's operand stack when
 creating the new frame and a function. All methods have the same signature, regardless
 of the signature of their Java counterparts. That signature is that it accepts a slice
 of interface{} and returns an interface{}. The accepted slice can be empty and the
 return interface can be nil. This covers all Java functions. (Objects are returned
 as a 64-bit address in this scheme (as they are in the JVM).

 The passed-in slice contains one entry for every parameter passed to the method (which
 could mean an empty slice).
*/

func Load_Lang_System() map[string]GMeth {

	MethodSignatures["java/lang/System.currentTimeMillis()J"] = // get time in ms since Jan 1, 1970, returned as long
		GMeth{
			ParamSlots: 0,
			GFunction:  currentTimeMillis,
		}

	MethodSignatures["java/lang/System.nanoTime()J"] = // get nanoseconds time, returned as long
		GMeth{
			ParamSlots: 0,
			GFunction:  nanoTime,
		}

	MethodSignatures["java/lang/System.exit(I)V"] = // shutdown the app
		GMeth{
			ParamSlots: 1,
			GFunction:  exitI,
		}

	return MethodSignatures
}

// Return time in milliseconds, measured since midnight of Jan 1, 1970
func currentTimeMillis([]interface{}) interface{} {
	return int64(time.Now().UnixMilli())
}

// Return time in nanoseconds. Note that in golang this function has a lower (that is, less good)
// resolution than Java: two successive calls often return the same value.
func nanoTime([]interface{}) interface{} {
	return int64(time.Now().UnixNano())
}

// Exits the program directly, returning the passed in value
func exitI(params []interface{}) interface{} {
	exitCode := params[0].(int64) // int64
	var exitStatus = int(exitCode)
	shutdown.Exit(exitStatus)
	return 0 // this code is not executed as previous line ends Jacobin
}
