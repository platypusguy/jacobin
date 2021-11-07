/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"encoding/binary"
	"fmt"
	"jacobin/log"
	"math"
	"os"
)

// this file contains the parser for the constant pool and the verifier.
// Refer to: https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4-140

// the various types of entries in the constant pool. Similar list in exec.classes
const (
	Dummy         = 0 // used for initialization and for dummy entries (viz. for longs, doubles)
	UTF8          = 1
	IntConst      = 3
	FloatConst    = 4
	LongConst     = 5
	DoubleConst   = 6
	ClassRef      = 7
	StringConst   = 8
	FieldRef      = 9
	MethodRef     = 10
	Interface     = 11
	NameAndType   = 12
	MethodHandle  = 15
	MethodType    = 16
	Dynamic       = 17
	InvokeDynamic = 18
	Module        = 19
	Package       = 20
)

// the constant pool, which is an array of different record types. Each entry in the table
// consists of an identifying integer (see enums above)
// all entries are defined at the end of this file
type cpEntry struct {
	entryType int
	slot      int
}

// parse the CP entries in the class file and put references to their data in klass.cpIndex,
// where appropriate. (Some entries, such as invokeDynamic, Module, etc. require other actions
// performed here. Returns location through last parsed byte and any error.
func parseConstantPool(rawBytes []byte, klass *ParsedClass) (int, error) {
	klass.cpIndex = make([]cpEntry, klass.cpCount)
	pos := 9 // position of the last byte before the constant pool

	klass.moduleName = ""

	klass.classRefs = []int{}
	klass.fieldRefs = []fieldRefEntry{}
	klass.intConsts = []int{}
	klass.invokeDynamics = []invokeDynamic{}
	klass.methodHandles = []methodHandleEntry{}
	klass.methodRefs = []methodRefEntry{}
	klass.nameAndTypes = []nameAndTypeEntry{}
	klass.stringRefs = []stringConstantEntry{}
	klass.utf8Refs = []utf8Entry{}

	// the first entry in the CP is a dummy entry, so that all references are 1-based
	klass.cpIndex[0] = cpEntry{Dummy, 0}

	var i int
	for i = 1; i <= klass.cpCount-1; { // i starts at 1 due to the dummy entry at CP[0]
		pos += 1
		entryType := int(rawBytes[pos])
		switch entryType {
		case UTF8:
			var content string
			length, _ := intFrom2Bytes(rawBytes, pos+1)
			pos += 2
			if length == 0 {
				content = ""
			} else {
				content = string(rawBytes[pos+1 : pos+length+1])
			}
			pos += length
			utfe := utf8Entry{content}
			klass.utf8Refs = append(klass.utf8Refs, utfe)
			klass.cpIndex[i] = cpEntry{UTF8, len(klass.utf8Refs) - 1}
			i += 1
		case IntConst:
			intValue, _ := intFrom4Bytes(rawBytes, pos+1)
			pos += 4
			klass.intConsts = append(klass.intConsts, intValue)
			klass.cpIndex[i] = cpEntry{IntConst, len(klass.intConsts) - 1}
			i += 1
		case FloatConst:
			bytes := make([]byte, 4)
			for j := 0; j < 4; j++ {
				bytes[j] = rawBytes[pos+1+j]
			}
			pos += 4
			bits := binary.BigEndian.Uint32(bytes)
			floatValue := math.Float32frombits(bits)
			klass.floats = append(klass.floats, floatValue)
			klass.cpIndex[i] = cpEntry{FloatConst, len(klass.floats) - 1}
			i++
		case LongConst:
			highBytes, _ := intFrom4Bytes(rawBytes, pos+1)
			lowBytes, _ := intFrom4Bytes(rawBytes, pos+5)
			pos += 8
			longValue := int64((highBytes << 32) + lowBytes)
			klass.longConsts = append(klass.longConsts, longValue)
			klass.cpIndex[i] = cpEntry{LongConst, len(klass.longConsts) - 1}
			i++
			// long ints take up two slots in the CP, of which the second is just a dummy slot.
			klass.cpIndex[i] = cpEntry{Dummy, 0}
			i++
		case DoubleConst:
			bytes := make([]byte, 8)
			for j := 0; j < 8; j++ {
				bytes[j] = rawBytes[pos+1+j]
			}
			pos += 8
			bits := binary.BigEndian.Uint64(bytes)
			doubleValue := math.Float64frombits(bits)
			klass.doubles = append(klass.doubles, doubleValue)
			klass.cpIndex[i] = cpEntry{DoubleConst, len(klass.doubles) - 1}
			i++
			// doubles take up two slots in the CP, of which the second is just a dummy slot.
			klass.cpIndex[i] = cpEntry{Dummy, 0}
			i++
		case ClassRef:
			index, _ := intFrom2Bytes(rawBytes, pos+1)
			// cre := classRefEntry{index}
			klass.classRefs = append(klass.classRefs, index)
			klass.cpIndex[i] = cpEntry{ClassRef, len(klass.classRefs) - 1}
			pos += 2
			i += 1
		case StringConst:
			index, _ := intFrom2Bytes(rawBytes, pos+1)
			sce := stringConstantEntry{index}
			klass.stringRefs = append(klass.stringRefs, sce)
			klass.cpIndex[i] = cpEntry{StringConst, len(klass.stringRefs) - 1}
			pos += 2
			i += 1
		case FieldRef:
			classIndex, _ := intFrom2Bytes(rawBytes, pos+1)
			nameAndTypeIndex, _ := intFrom2Bytes(rawBytes, pos+3)
			fre := fieldRefEntry{classIndex, nameAndTypeIndex}
			klass.fieldRefs = append(klass.fieldRefs, fre)
			klass.cpIndex[i] = cpEntry{FieldRef, len(klass.fieldRefs) - 1}
			pos += 4
			i += 1
		case MethodRef:
			classIndex, _ := intFrom2Bytes(rawBytes, pos+1)
			nameAndTypeIndex, _ := intFrom2Bytes(rawBytes, pos+3)
			mre := methodRefEntry{classIndex, nameAndTypeIndex}
			klass.methodRefs = append(klass.methodRefs, mre)
			klass.cpIndex[i] = cpEntry{MethodRef, len(klass.methodRefs) - 1}
			pos += 4
			i += 1
		case Interface:
			classIndex, _ := intFrom2Bytes(rawBytes, pos+1)
			nameAndTypeIndex, _ := intFrom2Bytes(rawBytes, pos+3)
			ire := interfaceRefEntry{classIndex, nameAndTypeIndex}
			klass.interfaceRefs = append(klass.interfaceRefs, ire)
			klass.cpIndex[i] = cpEntry{Interface, len(klass.interfaceRefs) - 1}
			pos += 4
			i += 1
		case NameAndType:
			nameIndex, _ := intFrom2Bytes(rawBytes, pos+1)
			descriptorIndex, _ := intFrom2Bytes(rawBytes, pos+3)
			nte := nameAndTypeEntry{nameIndex, descriptorIndex}
			klass.nameAndTypes = append(klass.nameAndTypes, nte)
			klass.cpIndex[i] = cpEntry{NameAndType, len(klass.nameAndTypes) - 1}
			pos += 4
			i += 1
		case MethodHandle:
			refKind := int(rawBytes[pos+1])
			refIndex, _ := intFrom2Bytes(rawBytes, pos+2)
			mhe := methodHandleEntry{refKind, refIndex}
			klass.methodHandles = append(klass.methodHandles, mhe)
			klass.cpIndex[i] = cpEntry{MethodHandle, len(klass.methodHandles) - 1}
			pos += 3
			i += 1
		case MethodType:
			descIndex, _ := intFrom2Bytes(rawBytes, pos+1)
			klass.methodTypes = append(klass.methodTypes, descIndex)
			klass.cpIndex[i] = cpEntry{MethodType, len(klass.methodTypes) - 1}
			pos += 2
			i += 1
		case Dynamic:
			bootstrap, _ := intFrom2Bytes(rawBytes, pos+1)
			nAndT, _ := intFrom2Bytes(rawBytes, pos+3)
			dyn := dynamic{
				bootstrapIndex: bootstrap,
				nameAndType:    nAndT,
			}
			klass.dynamics = append(klass.dynamics, dyn)
			klass.cpIndex[i] = cpEntry{Dynamic, len(klass.dynamics) - 1}
			pos += 4
			i += 1
		case InvokeDynamic:
			bootstrap, _ := intFrom2Bytes(rawBytes, pos+1)
			nAndT, _ := intFrom2Bytes(rawBytes, pos+3)
			ide := invokeDynamic{
				bootstrapIndex: bootstrap,
				nameAndType:    nAndT,
			}
			klass.invokeDynamics = append(klass.invokeDynamics, ide)
			klass.cpIndex[i] = cpEntry{InvokeDynamic, len(klass.invokeDynamics) - 1}
			pos += 4
			i += 1
		case Module:
			if klass.javaVersion < 53 {
				return pos, cfe("Java module record requires Java 9 or later version")
			}
			nameIndex, _ := intFrom2Bytes(rawBytes, pos+1)
			moduleName, err := fetchUTF8string(klass, nameIndex)
			if err != nil {
				break // error message will already have been shown
			}
			if klass.moduleName != "" {
				return pos + 2, cfe("Class " + klass.className + " has two module names: " + klass.moduleName +
					" and " + moduleName)
			}
			klass.moduleName = moduleName
			klass.cpIndex[i] = cpEntry{Module, nameIndex}
			pos += 2
			i += 1
		case Package:
			if klass.javaVersion < 53 {
				return pos, cfe("Java package entry requires Java 9 or later version")
			}
			nameIndex, _ := intFrom2Bytes(rawBytes, pos+1)
			packageName, err := fetchUTF8string(klass, nameIndex)
			if err != nil {
				break // error message will already have been shown
			}
			if klass.packageName != "" {
				return pos + 2, cfe("Class " + klass.className + " has two package names: " + klass.packageName +
					" and " + packageName)
			}
			klass.packageName = packageName
			klass.cpIndex[i] = cpEntry{Package, nameIndex}
			pos += 2
			i += 1

		default:
			klass.cpCount = i // just to get it over with for the moment
		}
	}

	if log.Level == log.FINEST {
		printCP(klass)

	}

	return pos, nil
}

// prints the entries in the CP. Accepts the number of entries for the nonce.
// func printCP(entries int, klass *ParsedClass) {
func printCP(klass *ParsedClass) {

	fmt.Fprintf(os.Stderr, "Number of CP entries parsed: %02d\n", len(klass.cpIndex))
	for j := 0; j < len(klass.cpIndex); j++ {
		entry := klass.cpIndex[j]
		fmt.Fprintf(os.Stderr, "CP entry: %02d, type %02d ", j, entry.entryType)
		switch entry.entryType {
		case Dummy:
			fmt.Fprintf(os.Stderr, "(dummy entry)\n")
		case UTF8:
			s := entry.slot
			fmt.Fprintf(os.Stderr, "(UTF-8 string)     %s\n", klass.utf8Refs[s].content)
		case IntConst:
			ic := entry.slot
			fmt.Fprintf(os.Stderr, "(int constant)     %d\n", klass.intConsts[ic])
		case FloatConst:
			fc := entry.slot
			fmt.Fprintf(os.Stderr, "(float constant)   %f\n", klass.floats[fc])
		case LongConst:
			lc := entry.slot
			fmt.Fprintf(os.Stderr, "(long constant)    %dL\n", klass.longConsts[lc])
		case DoubleConst:
			dc := entry.slot
			fmt.Fprintf(os.Stderr, "(double constant)  %f\n", klass.doubles[dc])
		case ClassRef:
			fmt.Fprintf(os.Stderr, "(class ref)        ")
			c := entry.slot
			fmt.Fprintf(os.Stderr, "index: %02d\n", klass.classRefs[c])
		case StringConst:
			fmt.Fprintf(os.Stderr, "(string const ref) ")
			s := entry.slot
			fmt.Fprintf(os.Stderr, "index: %02d\n", klass.stringRefs[s].index)
		case FieldRef:
			fmt.Fprintf(os.Stderr, "(field ref)        ")
			k := entry.slot
			fmt.Fprintf(os.Stderr, "class index: %02d, nameAndType index: %02d\n",
				klass.fieldRefs[k].classIndex, klass.fieldRefs[k].nameAndTypeIndex)
		case MethodRef:
			fmt.Fprintf(os.Stderr, "(method ref)       ")
			k := entry.slot
			fmt.Fprintf(os.Stderr, "class index: %02d, nameAndType index: %02d\n",
				klass.methodRefs[k].classIndex, klass.methodRefs[k].nameAndTypeIndex)
		case Interface:
			fmt.Fprintf(os.Stderr, "(interface ref)    ")
			k := entry.slot
			fmt.Fprintf(os.Stderr, "class index: %02d, nameAndType index: %02d\n",
				klass.interfaceRefs[k].classIndex, klass.interfaceRefs[k].nameAndTypeIndex)
		case NameAndType:
			fmt.Fprintf(os.Stderr, "(name and type)    ")
			n := entry.slot
			fmt.Fprintf(os.Stderr, "name index: %02d, descriptor index: %02d\n",
				klass.nameAndTypes[n].nameIndex, klass.nameAndTypes[n].descriptorIndex)
		case MethodHandle:
			fmt.Fprintf(os.Stderr, "(method handle)    ")
			m := entry.slot
			fmt.Fprintf(os.Stderr, "reference kind: %d, reference index: %02d\n",
				klass.methodHandles[m].referenceKind, klass.methodHandles[m].referenceIndex)
		case MethodType:
			fmt.Fprintf(os.Stderr, "(method type)      ")
			mt := entry.slot
			fmt.Fprintf(os.Stderr, "description index: %02d\n", klass.methodTypes[mt])
		case Dynamic:
			fmt.Fprintf(os.Stderr, "(dynamic)          ")
			n := entry.slot
			fmt.Fprintf(os.Stderr, "boostrap index: %02d, name and type: %02d\n",
				klass.dynamics[n].bootstrapIndex, klass.dynamics[n].nameAndType)
		case InvokeDynamic:
			fmt.Fprintf(os.Stderr, "(invokedynamic)    ")
			n := entry.slot
			fmt.Fprintf(os.Stderr, "boostrap index: %02d, name and type: %02d\n",
				klass.invokeDynamics[n].bootstrapIndex, klass.invokeDynamics[n].nameAndType)
		case Module:
			fmt.Fprintf(os.Stderr, "(module name)      ")
			fmt.Fprintf(os.Stderr, "%s\n", klass.moduleName)
		case Package:
			fmt.Fprintf(os.Stderr, "(package name)     ")
			fmt.Fprintf(os.Stderr, "%s\n", klass.packageName)
		default:
			fmt.Fprintf(os.Stderr, "invalid entry\n")
		}
	}
}

// ==== the various entry types in the constant pool (listed in order of the enums above) ====
// note that the DummyEntry, value 0, is never accessed and so has no corresponding struct

type utf8Entry struct { // type: 01 (UTF-8 string)
	content string
}

type stringConstantEntry struct { // type: 08 (string constant reference)
	index int
}

type fieldRefEntry struct { // type: 09 (field reference)
	classIndex       int
	nameAndTypeIndex int
}

type methodRefEntry struct { // type: 10 (method reference)
	classIndex       int
	nameAndTypeIndex int
}

type interfaceRefEntry struct { // type: 11 (interface reference)
	classIndex       int
	nameAndTypeIndex int
}

type nameAndTypeEntry struct { // type 12 (name and type reference)
	nameIndex       int
	descriptorIndex int
}

type methodHandleEntry struct { // type: 15 (method handle)
	referenceKind  int
	referenceIndex int
}

type dynamic struct { // type 17 (dynamic--similar to invokedynamic)
	bootstrapIndex int
	nameAndType    int
}

type invokeDynamic struct { // type 18 (invokedynamic data)
	bootstrapIndex int
	nameAndType    int
}
