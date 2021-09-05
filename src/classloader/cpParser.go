/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"fmt"
	"jacobin/log"
	"os"
)

// this file contains the parser for the constant pool and the verifier.
// Refer to: https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4-140

// the various types of entries in the constant pool
const (
	Invalid       = -1 // used for initialization and for dummy entries (viz. for longs, doubles)
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
// // and a structure that has varying fields and fulfills the cpTyper interface.
// // all entries are defined at the end of this file
type cpEntry struct {
	entryType int
	slot      int
}

// parse the CP entries in the class file and put references to their data in klass.cpIndex
func parseConstantPool(rawBytes []byte, klass *parsedClass) (int, error) {
	klass.cpIndex = make([]cpEntry, klass.cpCount)
	pos := 9 // position of the last byte before the constant pool

	// the first entry in the CP is a dummy entry, so that all references are 1-based
	klass.cpIndex[0] = cpEntry{Invalid, 0}

	klass.classRefs = []classRefEntry{}
	klass.fieldRefs = []fieldRefEntry{}
	klass.intConsts = []intConst{}
	klass.methodRefs = []methodRefEntry{}
	klass.nameAndTypes = []nameAndTypeEntry{}
	klass.stringRefs = []stringConstantEntry{}
	klass.utf8Refs = []utf8Entry{}

	i := 0
	for i = 1; i <= klass.cpCount-1; {
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
			ice := intConst{intValue}
			klass.intConsts = append(klass.intConsts, ice)
			klass.cpIndex[i] = cpEntry{IntConst, len(klass.intConsts) - 1}
			i += 1
		case ClassRef:
			index, _ := intFrom2Bytes(rawBytes, pos+1)
			cre := classRefEntry{index}
			klass.classRefs = append(klass.classRefs, cre)
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
		case NameAndType:
			nameIndex, _ := intFrom2Bytes(rawBytes, pos+1)
			descriptorIndex, _ := intFrom2Bytes(rawBytes, pos+3)
			nte := nameAndTypeEntry{nameIndex, descriptorIndex}
			klass.nameAndTypes = append(klass.nameAndTypes, nte)
			klass.cpIndex[i] = cpEntry{NameAndType, len(klass.nameAndTypes) - 1}
			pos += 4
			i += 1
		default:
			klass.cpCount = i // just to get it over with for the moment
		}
	}

	if log.LogLevel == log.FINEST {
		printCP(i, klass)
	}

	return pos, nil
}

// prints the entries in the CP. Accepts the number of entries for the nonce.
func printCP(entries int, klass *parsedClass) {
	fmt.Fprintf(os.Stderr, "Number of CP entries parsed: %02d\n", len(klass.cpIndex))
	for j := 0; j < entries; j++ {
		entry := klass.cpIndex[j]
		fmt.Fprintf(os.Stderr, "CP entry: %02d, type %02d ", j, entry.entryType)
		switch entry.entryType {
		case Invalid:
			fmt.Fprintf(os.Stderr, "(dummy entry)\n")
		case IntConst:
			ic := entry.slot
			fmt.Fprintf(os.Stderr, "(int constant)     %d\n", klass.intConsts[ic].value)
		case UTF8:
			s := entry.slot
			fmt.Fprintf(os.Stderr, "(UTF-8 string)     %s\n", klass.utf8Refs[s].content)
		case ClassRef:
			fmt.Fprintf(os.Stderr, "(class ref)        ")
			c := entry.slot
			fmt.Fprintf(os.Stderr, "index: %02d\n", klass.classRefs[c].index)
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
		case NameAndType:
			fmt.Fprintf(os.Stderr, "(name and type)    ")
			n := entry.slot
			fmt.Fprintf(os.Stderr, "name index: %02d, descriptor index: %02d\n",
				klass.nameAndTypes[n].nameIndex, klass.nameAndTypes[n].descriptorIndex)
		default:
			fmt.Fprintf(os.Stderr, "invalid entry\n")
		}
	}
}

// ==== the various entry types in the constant pool (listed in order of the enums above) ====
type dummyEntry struct { // type -1 (invalid or dummy entry)
}

type utf8Entry struct { // type: 01 (UTF-8 string)
	content string
}

type intConst struct { // type 03 (integer constant)
	value int
}

type classRefEntry struct { // type: 07 (class refence -- points to UTF8 entry)
	index int
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

type nameAndTypeEntry struct { // type 12 (name and type reference)
	nameIndex       int
	descriptorIndex int
}
