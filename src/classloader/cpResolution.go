/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"errors"
	"fmt"
	"jacobin/globals"
	"jacobin/stringPool"
	"os"
)

// ResolveCPmethRefs resolves the method references in the constant pool of a class
func ResolveCPinterfaceRefs(cpp *CPool) error {
	cp := *cpp
	if cp.CpIndex == nil || cp.MethodRefs == nil {
		return errors.New("invalid constant pool or class data passed to classloader.ResolveCPmethRefs()")
	}

	for _, interfaceEntry := range cp.InterfaceRefs {
		resEntry := ResolvedInterfaceRefEntry{}
		// get the class name as an index into the string pool
		classIndex := interfaceEntry.ClassIndex
		classRefIdx := cp.CpIndex[classIndex].Slot
		classIdx := cp.ClassRefs[classRefIdx]
		resEntry.ClassIndex = classIdx

		// now get the method signature as an index into the string pool
		nameAndTypeCPindex := interfaceEntry.NameAndType
		nameAndTypeIndex := cp.CpIndex[nameAndTypeCPindex].Slot
		nameAndType := cp.NameAndTypes[nameAndTypeIndex]
		methNameCPindex := nameAndType.NameIndex
		methNameUTF8index := cp.CpIndex[methNameCPindex].Slot
		methName := cp.Utf8Refs[methNameUTF8index]
		resEntry.NameIndex = stringPool.GetStringIndex(&methName)

		// and get the method type
		methSigCPindex := nameAndType.DescIndex
		methSigUTF8index := cp.CpIndex[methSigCPindex].Slot
		methSig := cp.Utf8Refs[methSigUTF8index]
		resEntry.TypeIndex = stringPool.GetStringIndex(&methSig)

		// append them all together to get the fully qualified method signature
		fqn := *stringPool.GetStringPointer(resEntry.ClassIndex) + "." + methName + methSig
		resEntry.FQNameIndex = stringPool.GetStringIndex(&fqn)

		cpp.ResolvedInterfaceRefs = append(cpp.ResolvedInterfaceRefs, resEntry)
		if globals.TraceInst {
			fmt.Fprintf(os.Stderr, "Resolved interface ref: %s\n", fqn)
		}
	}
	return nil
}

// ResolveCPmethRefs resolves the method references in the constant pool of a class
func ResolveCPmethRefs(cpp *CPool) error {
	cp := *cpp
	if cp.CpIndex == nil || cp.MethodRefs == nil {
		return errors.New("invalid constant pool or class data passed to classloader.ResolveCPmethRefs()")
	}

	for _, methEntry := range cp.MethodRefs {
		resEntry := ResolvedMethodRefEntry{}
		// get the class name as an index into the string pool
		classIndex := methEntry.ClassIndex
		classRefIdx := cp.CpIndex[classIndex].Slot
		classIdx := cp.ClassRefs[classRefIdx]
		resEntry.ClassIndex = classIdx

		// now get the method signature as an index into the string pool
		nameAndTypeCPindex := methEntry.NameAndType
		nameAndTypeIndex := cp.CpIndex[nameAndTypeCPindex].Slot
		nameAndType := cp.NameAndTypes[nameAndTypeIndex]
		methNameCPindex := nameAndType.NameIndex
		methNameUTF8index := cp.CpIndex[methNameCPindex].Slot
		methName := cp.Utf8Refs[methNameUTF8index]
		resEntry.NameIndex = stringPool.GetStringIndex(&methName)

		// and get the method type
		methSigCPindex := nameAndType.DescIndex
		methSigUTF8index := cp.CpIndex[methSigCPindex].Slot
		methSig := cp.Utf8Refs[methSigUTF8index]
		resEntry.TypeIndex = stringPool.GetStringIndex(&methSig)

		// append them all together to get the fully qualified method signature
		fqn := *stringPool.GetStringPointer(resEntry.ClassIndex) + "." + methName + methSig
		resEntry.FQNameIndex = stringPool.GetStringIndex(&fqn)

		cpp.ResolvedMethodRefs = append(cpp.ResolvedMethodRefs, resEntry)
		if globals.TraceInst {
			fmt.Fprintf(os.Stderr, "Resolved method ref: %s\n", fqn)
		}
	}
	return nil
}

/*
methodRef := CP.CpIndex[cpIndex].Slot
	classIndex := CP.MethodRefs[methodRef].ClassIndex

	classRefIdx := CP.CpIndex[classIndex].Slot
	classIdx := CP.ClassRefs[classRefIdx]
	classNamePtr := stringPool.GetStringPointer(uint32(classIdx))
	className := *classNamePtr

	// now get the method signature
	nameAndTypeCPindex := CP.MethodRefs[methodRef].NameAndType
	nameAndTypeIndex := CP.CpIndex[nameAndTypeCPindex].Slot
	nameAndType := CP.NameAndTypes[nameAndTypeIndex]
	methNameCPindex := nameAndType.NameIndex
	methNameUTF8index := CP.CpIndex[methNameCPindex].Slot
	methName := CP.Utf8Refs[methNameUTF8index]

	// and get the method signature/description
	methSigCPindex := nameAndType.DescIndex
	methSigUTF8index := CP.CpIndex[methSigCPindex].Slot
	methSig := CP.Utf8Refs[methSigUTF8index]

*/
