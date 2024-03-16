/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"errors"
	"jacobin/stringPool"
	"strconv"
)

// various utilities frequently used in parsing classfiles

// read two bytes in big endian order and convert to an int
func intFrom2Bytes(bytes []byte, pos int) (int, error) {
	if len(bytes) < pos+2 {
		return 0, cfe("invalid offset into file")
	}

	value := (uint16(bytes[pos]) << 8) + uint16(bytes[pos+1])
	return int(value), nil
}

// the same as inFrom2Bytes(), but returns a uint16.
func u16From2bytes(bytes []byte, pos int) (uint16, error) {
	i, err := intFrom2Bytes(bytes, pos)
	if err != nil {
		return 0, err
	}

	return uint16(i), nil
}

// read four bytes in big endian order and convert to an int
func intFrom4Bytes(bytes []byte, pos int) (int, error) {
	if len(bytes) < pos+4 {
		return 0, cfe("invalid offset into file")
	}

	value1 := (uint32(bytes[pos]) << 8) + uint32(bytes[pos+1])
	value2 := (uint32(bytes[pos+2]) << 8) + uint32(bytes[pos+3])
	retVal := int(value1<<16) + int(value2)
	return retVal, nil
}

// finds and returns a UTF8 string when handed an index into the CP that points
// to a UTF8 entry. Does extensive checking of values.
func FetchUTF8string(klass *ParsedClass, index int) (string, error) {
	if index < 1 || index > klass.cpCount-1 {
		return "", cfe("attempt to fetch invalid UTF8 at CP entry #" + strconv.Itoa(index))
	}

	if klass.cpIndex[index].entryType != UTF8 {
		return "", cfe("attempt to fetch UTF8 string from non-UTF8 CP entry #" + strconv.Itoa(index))
	}

	i := klass.cpIndex[index].slot
	if i < 0 || i > len(klass.utf8Refs)-1 {
		return "", cfe("invalid index into UTF8 array of CP: " + strconv.Itoa(i))
	}

	return klass.utf8Refs[i].content, nil
}

// finds and returns a UTF8 string when handed an index into the CP that points
// to a UTF8 entry in a loaded class. Similar to previous method, except that it
// works with a loaded class, rather than a parsed (i.e., pre-loaded) class.
// Does extensive checking of values.
func FetchUTF8stringInLoadedClass(klass *Klass, index int) (string, error) {
	if index < 1 || index > len(klass.Data.CP.CpIndex)-1 {
		return "", cfe("attempt to fetch invalid UTF8 at CP entry #" + strconv.Itoa(index))
	}

	if klass.Data.CP.CpIndex[index].Type != UTF8 {
		return "", cfe("attempt to fetch UTF8 string from non-UTF8 CP entry #" + strconv.Itoa(index))
	}

	i := int(klass.Data.CP.CpIndex[index].Slot)
	if i < 0 || i > len(klass.Data.CP.Utf8Refs)-1 {
		return "", cfe("invalid index into UTF8 array of CP: " + strconv.Itoa(i))
	}

	return klass.Data.CP.Utf8Refs[i], nil
}

// like the preceding function, except this returns the slot number in the utf8Refs
// rather than the string that's in that slot.
func fetchUTF8slot(klass *ParsedClass, index int) (int, error) {
	if index < 1 || index > klass.cpCount-1 {
		return -1, cfe("attempt to fetch invalid UTF8 at CP entry #" + strconv.Itoa(index))
	}

	if klass.cpIndex[index].entryType != UTF8 {
		return -1, cfe("attempt to fetch UTF8 string from non-UTF8 CP entry #" + strconv.Itoa(index))
	}

	slot := klass.cpIndex[index].slot
	if slot < 0 || slot > len(klass.utf8Refs)-1 {
		return -1, cfe("invalid index into UTF8 array of CP: " + strconv.Itoa(slot))
	}
	return slot, nil
}

// fetches attribute info. Attributes are values associated with fields, methods, classes, and
// code attributes (yes, the word 'attribute' is overloaded in JVM parlance). The spec is here:
// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.7 and the general
// layout is:
//
//	attribute_info {
//	   u2 attribute_name_index;  // the name of the attribute
//	   u4 attribute_length;
//	   u1 info[attribute_length];
//	}
func fetchAttribute(klass *ParsedClass, bytes []byte, loc int) (attr, int, error) {
	pos := loc
	attribute := attr{}
	nameIndex, err := intFrom2Bytes(bytes, pos+1)
	pos += 2
	if err != nil {
		return attribute, pos, cfe("error fetching field attribute")
	}
	nameSlot, err := fetchUTF8slot(klass, nameIndex)
	if err != nil {
		return attribute, pos, cfe("error fetching name of field attribute")
	}

	attribute.attrName = nameSlot // slot in UTF-8 slice of CP

	length, err := intFrom4Bytes(bytes, pos+1)
	pos += 4
	if err != nil {
		return attribute, pos, cfe("error fetching length of field attribute")
	}
	attribute.attrSize = length

	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = bytes[pos+1+i]
	}

	attribute.attrContent = b

	return attribute, pos + length, nil
}

// returns all the elements of a methodRef (10) CP entry when given the CP entry #
//
//	classIndex       int
//	nameAndTypeIndex int
func resolveCPmethodRef(index int, klass *ParsedClass) (string, string, string, error) {
	if index < 1 || index >= len(klass.cpIndex) {
		return "", "", "", cfe("Invalid index into CP: " + strconv.Itoa(index))
	}
	cpEnt := klass.cpIndex[index]
	if cpEnt.entryType != MethodRef {
		return "", "", "", cfe("Expecting MethodRef (10) at CP entry #" + strconv.Itoa(index) +
			" but instead got CP type: " + strconv.Itoa(cpEnt.entryType))
	}

	methRef := klass.methodRefs[cpEnt.slot]
	pointedToClassRef := klass.cpIndex[methRef.classIndex]
	nameIndex := klass.classRefs[pointedToClassRef.slot]
	// className, err := FetchUTF8string(klass, nameIndex)
	classNamePtr := stringPool.GetStringPointer(nameIndex)
	if classNamePtr != nil {
		return "", "", "", cfe("ClassRef entry in MethodRef CP entry #" + strconv.Itoa(index) +
			" does not point to a valid string")
	}

	// pointedToNandT := klass.cpIndex[methRef.nameAndTypeIndex]
	methName, methType, err := ResolveCPnameAndType(klass, methRef.nameAndTypeIndex)
	if err != nil {
		return "", "", "", errors.New("error occurred") // the error msg is displayed in the called func.
	}

	return *classNamePtr, methName, methType, nil
}

func ResolveCPnameAndType(klass *ParsedClass, index int) (string, string, error) {
	if index < 1 || index >= len(klass.cpIndex) {
		return "", "", cfe("Invalid nameAndType index into CP: " +
			strconv.Itoa(index))
	}

	nAndTindex := klass.cpIndex[index]
	nAndT := klass.nameAndTypes[nAndTindex.slot]
	nameIndex := nAndT.nameIndex
	descIndex := nAndT.descriptorIndex

	if klass.cpIndex[nameIndex].entryType != UTF8 {
		return "", "", cfe("Name index in nameAndType entry (CP #" + strconv.Itoa(index) +
			") does not point to a UTF8 entry.")
	}

	name := klass.utf8Refs[klass.cpIndex[nameIndex].slot]

	if klass.cpIndex[descIndex].entryType != UTF8 {
		return "", "", cfe("Desc index in nameAndType entry (CP #" + strconv.Itoa(index) +
			") does not point to a UTF8 entry.")
	}

	desc := klass.utf8Refs[klass.cpIndex[descIndex].slot]
	return name.content, desc.content, nil
}
