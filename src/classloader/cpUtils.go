/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-4 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

// This file contains utility routines for runtime operations involving the
// class's constant pool (CP). It was formerly in the jvm package, but moved
// here to avoid circular dependencies.

import (
	"jacobin/src/stringPool"
)

type CpType struct {
	EntryType int
	RetType   int
	IntVal    int64
	FloatVal  float64
	AddrVal   *CPuint16s
	StringVal *string
}

// all the multi-value structs use this layout. Depending on the type of
// the CPentry, the fields entry1 and entry2 will have different meanings.
type CPuint16s struct {
	entry1 uint16
	entry2 uint16
}

var IS_ERROR = 0
var IS_STRUCT_ADDR = 1
var IS_FLOAT64 = 2
var IS_INT64 = 3
var IS_STRING_ADDR = 4

// Utility routines for runtime operations

// FetchCPentry looks up an entry in a CP and return its type and
// its value. The returned value is a struct that serves as
// a substitute for a discriminated union.
// The fields are:
//  1. EntryType: The CP entry type. Equals 0 if an error occurred.
//     The five EntryType values are listed above: IS_ERROR, etc.
//  2. RetType: int that identifies the type of the returned value.
//     The options are:
//     0 = error
//     1 = address of item other than string
//     2 = float64
//     3 = int64
//     4 = address of string
//  3. three fields that hold an int64, float64, or 64-bit address, respectively.
//     The calling function checks the RetType field to determine which
//     of these three fields holds the returned value.
func FetchCPentry(cpp *CPool, index int) CpType {
	if cpp == nil {
		return CpType{EntryType: 0, RetType: IS_ERROR}
	}
	cp := *cpp
	// if index is out of range, return error
	if index < 1 || index >= len(cp.CpIndex) {
		return CpType{EntryType: 0, RetType: IS_ERROR}
	}

	entry := cp.CpIndex[index]

	switch entry.Type {
	// integers
	case IntConst:
		retInt := int64(cp.IntConsts[entry.Slot])
		return CpType{EntryType: int(entry.Type), RetType: IS_INT64, IntVal: retInt}

	case LongConst:
		retInt := cp.LongConsts[entry.Slot]
		return CpType{EntryType: int(entry.Type), RetType: IS_INT64, IntVal: retInt}

	case MethodType: // method type is an integer
		retInt := int64(cp.MethodTypes[entry.Slot])
		return CpType{EntryType: int(entry.Type), RetType: IS_INT64, IntVal: retInt}

	// floating point
	case FloatConst:
		retFloat := float64(cp.Floats[entry.Slot])
		return CpType{EntryType: int(entry.Type), RetType: IS_FLOAT64, FloatVal: retFloat}

	case DoubleConst:
		retFloat := cp.Doubles[entry.Slot]
		return CpType{EntryType: int(entry.Type), RetType: IS_FLOAT64, FloatVal: retFloat}

	// addresses of strings
	case ClassRef: // points to a CP entry, which is a string pool entry
		e := cp.ClassRefs[entry.Slot]
		classNamePtr := stringPool.GetStringPointer(uint32(e))

		className := ""
		if classNamePtr != nil {
			className = *classNamePtr
		}

		return CpType{EntryType: int(entry.Type),
			RetType: IS_STRING_ADDR, StringVal: &className}

	case StringConst: // points to a CP entry, which is a UTF-8 string constant
		e := cp.CpIndex[entry.Slot]
		// should point to a UTF-8
		if e.Type != UTF8 {
			return CpType{EntryType: 0, RetType: IS_ERROR}
		}

		str := cp.Utf8Refs[e.Slot]
		return CpType{EntryType: int(entry.Type),
			RetType: IS_STRING_ADDR, StringVal: &str}
	case UTF8: // points to a UTF-8 string
		v := &(cp.Utf8Refs[entry.Slot])
		return CpType{EntryType: int(entry.Type), RetType: IS_STRING_ADDR, StringVal: v}

	// addresses of structures or other elements
	case Dynamic:
		dyn := cp.Dynamics[entry.Slot]
		cpe := CPuint16s{
			entry1: dyn.BootstrapIndex,
			entry2: dyn.NameAndType,
		}
		return CpType{EntryType: int(entry.Type), RetType: IS_STRUCT_ADDR, AddrVal: &cpe}

	case Interface:
		iface := cp.InterfaceRefs[entry.Slot]
		cpe := CPuint16s{
			entry1: iface.ClassIndex,
			entry2: iface.NameAndType,
		}
		return CpType{EntryType: int(entry.Type), RetType: IS_STRUCT_ADDR, AddrVal: &cpe}

	case InvokeDynamic:
		idyn := cp.InvokeDynamics[entry.Slot]
		cpe := CPuint16s{
			entry1: idyn.BootstrapIndex,
			entry2: idyn.NameAndType,
		}
		return CpType{EntryType: int(entry.Type), RetType: IS_STRUCT_ADDR, AddrVal: &cpe}

	case MethodHandle:
		mh := cp.MethodHandles[entry.Slot]
		cpe := CPuint16s{
			entry1: mh.RefKind,
			entry2: mh.RefIndex,
		}
		return CpType{EntryType: int(entry.Type), RetType: IS_STRUCT_ADDR, AddrVal: &cpe}

	case MethodRef:
		mr := cp.MethodRefs[entry.Slot]
		cpe := CPuint16s{
			entry1: mr.ClassIndex,
			entry2: mr.NameAndType,
		}
		return CpType{EntryType: int(entry.Type), RetType: IS_STRUCT_ADDR, AddrVal: &cpe}

	case NameAndType:
		nat := cp.NameAndTypes[entry.Slot]
		cpe := CPuint16s{
			entry1: nat.NameIndex,
			entry2: nat.DescIndex,
		}
		return CpType{EntryType: int(entry.Type), RetType: IS_STRUCT_ADDR, AddrVal: &cpe}

	// error: name of module or package would
	// not normally be retrieved here
	case Module,
		Package:
		return CpType{EntryType: 0, RetType: IS_ERROR}
	}

	return CpType{EntryType: 0, RetType: IS_ERROR}
}

// GetMethInfoFromCPmethref receives a CP entry index that points to a method or interface
// and returns the class name, method name, method signature, and these three combined as a
// fully qualified name (FQN). Note that checks on the validity of the cpIndex are performed
// in codeCheck.go.
func GetMethInfoFromCPmethref(CP *CPool, cpIndex int) (string, string, string, string) {
	cp := *CP
	meth := cp.ResolvedMethodRefs[cp.CpIndex[cpIndex].Slot]
	cls := *stringPool.GetStringPointer(meth.ClassIndex)
	mth := *stringPool.GetStringPointer(meth.NameIndex)
	typ := *stringPool.GetStringPointer(meth.TypeIndex)
	fqn := *stringPool.GetStringPointer(meth.FQNameIndex)
	return cls, mth, typ, fqn
}

func GetMethInfoFromCPinterfaceRef(CP *CPool, cpIndex int) (string, string, string) {

	methodRef := CP.CpIndex[cpIndex].Slot
	classIndex := CP.InterfaceRefs[methodRef].ClassIndex

	classRefIdx := CP.CpIndex[classIndex].Slot
	classIdx := CP.ClassRefs[classRefIdx]
	classNamePtr := stringPool.GetStringPointer(uint32(classIdx))
	className := *classNamePtr

	// now get the method signature
	nameAndTypeCPindex := CP.InterfaceRefs[methodRef].NameAndType
	nameAndTypeIndex := CP.CpIndex[nameAndTypeCPindex].Slot
	nameAndType := CP.NameAndTypes[nameAndTypeIndex]
	methNameCPindex := nameAndType.NameIndex
	methNameUTF8index := CP.CpIndex[methNameCPindex].Slot
	methName := CP.Utf8Refs[methNameUTF8index]

	// and get the method signature/description
	methSigCPindex := nameAndType.DescIndex
	methSigUTF8index := CP.CpIndex[methSigCPindex].Slot
	methSig := CP.Utf8Refs[methSigUTF8index]

	// className, methName, methSig, _ := GetMethInfoFromCPmethref(CP, cpIndex)
	return className, methName, methSig
}

// accepts the index of a CP entry, which should point to a classref
// and resolves it to return a string containing the class name.
// Returns an empty string if an error occurred
func GetClassNameFromCPclassref(CP *CPool, cpIndex uint16) string {
	entry := FetchCPentry(CP, int(cpIndex))
	if entry.RetType != IS_STRING_ADDR {
		return ""
	} else {
		return *entry.StringVal
	}
}
