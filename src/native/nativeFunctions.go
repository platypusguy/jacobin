/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package native

// In this file we handle how we flag native methods that Jacobin does not
// or cannot run. Eventually, many of these will be replaced by functions
// written in golang.

// list where we add unsupported native methods and delete them as we implement them
var unsupportedNativeMethodsList = []string{
	"test.entry", // for testing only
	"java/util/zip/Adler32.update",
}

// UnsupportedNativeMethods is the look up table. The value field is unused.
// We're only interested in looking up the key, whose presence in the table
// identifies the method as native and unsupported.
var UnsupportedNativeMethods = make(map[string]interface{})

// IsUnsupportedNativeMethod looks up the method in the table of
// unsupported native methods. If it's there it returns true, otherwise false.
// Note that the method name includes the class: /class/name.method
func IsUnsupportedNativeMethod(methodName string) bool {
	_, ok := UnsupportedNativeMethods[methodName]
	return ok
}

// LoadUnsupportedNativeMethods loads the list of unsupported methods into
// the lookup table. Called at JVM start-up.
func LoadUnsupportedNativeMethods() int64 {
	for _, methodName := range unsupportedNativeMethodsList {
		UnsupportedNativeMethods[methodName] = nil
	}
	return int64(len(UnsupportedNativeMethods))
}
