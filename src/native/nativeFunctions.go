/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package native

func IsUnsupportedNativeMethod(methodName string) bool {
	return false
}

var unsupportedNativeMethodsList = []string{
	"test.entry",
	"java/util/zip/Adler32.update",
}

var UnsupportedNativeMethods = make(map[string]interface{})

func LoadUnsupportedNativeMethods() int64 {
	for _, methodName := range unsupportedNativeMethodsList {
		UnsupportedNativeMethods[methodName] = nil
	}
	return int64(len(UnsupportedNativeMethods))
}
