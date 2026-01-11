/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"runtime"
)

func Load_Util_Concurrent_Atomic_Atomic_Long() {

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	// Native functions or caller to native functions

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.VMSupportsCS8()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  atomicLongVMSupportsCS8,
		}

}

// atomicLongVMSupportsCS8 simulates the behavior of the native VMSupportsCS8() method
// "java/util/concurrent/atomic/AtomicLong.VMSupportsCS8()Z"
func atomicLongVMSupportsCS8([]interface{}) interface{} {
	// Check if the current architecture supports 64-bit atomic operations
	arch := runtime.GOARCH

	// List of architectures that typically support 64-bit atomic operations
	supportedArchitectures := map[string]bool{
		"amd64":    true, // x86-64 (Intel/AMD)
		"arm64":    true, // ARM 64-bit
		"ppc64":    true, // PowerPC 64-bit
		"ppc64le":  true, // Little-endian PowerPC 64-bit
		"s390x":    true, // IBM Z (System z) 64-bit
		"sparc64":  true, // SPARC 64-bit
		"mips64":   true, // MIPS 64-bit
		"mips64le": true, // Little-endian MIPS 64-bit
	}

	// Check if the current architecture is in the supported list
	if supported, ok := supportedArchitectures[arch]; ok {
		return object.JavaBooleanFromGoBoolean(supported)
	}

	// If architecture is not recognized or supported, return false by default
	return object.JavaBooleanFromGoBoolean(false)
}
