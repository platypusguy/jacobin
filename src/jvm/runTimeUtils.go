/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"jacobin/classloader"
	"math"
)

// Utility routines for runtime operations

// Look up an entry in a CP and return its type and its value.
// The three return fields are:
//  1. The CP entry type. If it's a Dummy (0), this indidates that
//     an error occurred, as we would never look up a dummy entry.
//  2. an int that identifies the type of the returned value. The
//     options are:
//     0 = error
//     1 = int64 or address
//     2 = float64
//     3 = a string address
//  3. the value itself
func FetchCPentry(cpp *classloader.CPool, index int) (uint16, int, any) {
	if cpp == nil {
		return classloader.Dummy, 0, math.NaN()
	}

	cp := *cpp
	// if index is out of range, return error
	if index < 1 || index >= len(cp.CpIndex) {
		return classloader.Dummy, 0, math.NaN()
	}

	entry := cp.CpIndex[index]

	switch entry.Type {
	// integers
	case classloader.IntConst:
		retInt := int64(cp.IntConsts[entry.Slot])
		return entry.Type, 1, retInt
	case classloader.LongConst:
		retInt := cp.LongConsts[entry.Slot]
		return entry.Type, 1, retInt

	case classloader.MethodType:
		return entry.Type, 1, int64(cp.MethodTypes[entry.Slot])

	// floating point
	case classloader.FloatConst:
		retFloat := float64(cp.Floats[entry.Slot])
		return entry.Type, 2, retFloat
	case classloader.DoubleConst:
		retFloat := cp.Doubles[entry.Slot]
		return entry.Type, 2, retFloat

	// addresses and references
	case classloader.ClassRef: // points to a UTF-8 string
		return entry.Type, 3, &(cp.Utf8Refs[entry.Slot])

	case classloader.Dynamic:
		return entry.Type, 1, &(cp.Dynamics[entry.Slot])

	case classloader.Interface:
		return entry.Type, 1, &(cp.InterfaceRefs[entry.Slot])

	case classloader.InvokeDynamic:
		return entry.Type, 1, &(cp.InvokeDynamics[entry.Slot])

	case classloader.MethodHandle:
		return entry.Type, 1, &(cp.MethodHandles[entry.Slot])

	case classloader.MethodRef:
		return entry.Type, 1, &(cp.MethodRefs[entry.Slot])

	case classloader.NameAndType:
		return entry.Type, 1, &(cp.NameAndTypes[entry.Slot])

	// error: name of module would not normally be retrieved here
	case classloader.Module, classloader.Package:
		return 0, 0, math.NaN()
	}

	return classloader.Dummy, 0, math.NaN()
}
