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
func ResolveCPmethRefs(k *Klass) error {
	if k == nil || k.Data == nil || &k.Data.CP == nil {
		return errors.New("invalid class or class data in ResolveCPmethRefs")
	}
	cp := k.Data.CP
	resEntry := ResolvedMethodRefEntry{}

	for _, methEntry := range cp.MethodRefs {
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

		k.Data.CP.ResolvedMethodRefs = append(k.Data.CP.ResolvedMethodRefs, resEntry)
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
