/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import "strconv"

// Performs the format check on a fully parsed class. The requirements are listed
// here: https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.8
// They are:
// 1) must start with 0xCAFEBABE -- this is verified in the parsing, so not done here
// 2) most predefined attributes must be the right length -- verified during parsing
// 3) class must not be truncated or have extra bytes -- verified during parsing
// 4) CP must fulfill all constraints. This is done in this function
// 5) Fields must have valid names, classes, and descriptions. Partially done in
//    the parsing, but entirely done below
func formatCheckClass(klass *parsedClass) error {
	err := validateConstantPool(klass)
	if err != nil {
		return err // whatever occurs will there notify the user
	}

	err = validateFields(klass)
	return err
}

// validates that the CP fits all the requirements enumerated in:
// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4
// some of these checks were performed perforce in the parsing. Here, however,
// we verify them all. This is a requirement of all classes loaded in the JVM
// Note that this is *not* part of the larger class verification process.
func validateConstantPool(klass *parsedClass) error {
	cpSize := klass.cpCount
	if len(klass.cpIndex) != cpSize {
		return cfe("Error in size of constant pool discovered in format check." +
			"Expected: " + strconv.Itoa(cpSize) + ", got: " + strconv.Itoa(len(klass.cpIndex)))
	}

	if klass.cpIndex[0].entryType != Dummy {
		return cfe("Missing dummy entry in first slot of constant pool")
	}

	for j := 1; j < cpSize; j++ {
		entry := klass.cpIndex[j]
		switch entry.entryType {
		case UTF8:
			// points to an entry in utf8Refs, which holds a string. Check for:
			// * No byte may have the value (byte)0.
			// * No byte may lie in the range (byte)0xf0 to (byte)0xff
			whichUtf8 := entry.slot
			if whichUtf8 < 0 || whichUtf8 >= len(klass.utf8Refs) {
				return cfe("CP entry #" + strconv.Itoa(j) + "points to invalid UTF8 entry: " +
					strconv.Itoa(whichUtf8))
			}
			utf8string := klass.utf8Refs[whichUtf8].content
			utf8bytes := []byte(utf8string)
			for _, char := range utf8bytes {
				if char == 0x00 || (char >= 0xf0 && char <= 0xff) {
					return cfe("UTF8 string for CP entry #" + strconv.Itoa(j) +
						" contains an invalid character")
				}
			}
		case IntConst:
			// there are no specific format checks for integers, so we only check
			// that there is a valid entry pointed to in intConsts
			whichInt := entry.slot
			if whichInt < 0 || whichInt >= len(klass.intConsts) {
				return cfe("Integer at CP entry #" + strconv.Itoa(j) +
					" points to an invalid entry in CP intConsts")
			}
		case FloatConst:
			// there are complex bit patterns that can be enforced for floats, but
			// for the nonce, we'll just make sure that the float index points to an actual value
			whichFloat := entry.slot
			if whichFloat < 0 || whichFloat >= len(klass.floats) {
				return cfe("Float at CP entry #" + strconv.Itoa(j) +
					" points to an invalid entry in CP floats")
			}
		case ClassRef:
			// the only field of a ClassRef points to a UTF8 entry holding the class name
			// in the case of arrays, the UTF8 entry will describe the type and dimensions of the array
			whichClassRef := entry.slot
			if whichClassRef < 0 || whichClassRef >= len(klass.utf8Refs) {
				return cfe("ClassRef at CP entry #" + strconv.Itoa(j) +
					" points to an invalid entry in CP utf8Refs")
			}
		case StringConst:
			// a StringConst holds only an index into the utf8Refs. so we check this.
			// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4.3
			whichString := entry.slot
			if whichString < 0 || whichString >= len(klass.utf8Refs) {
				return cfe("Constant String at CP entry #" + strconv.Itoa(j) +
					" points to an invalid entry in CP utf8Refs")
			}
		case NameAndType:
			// a NameAndType entry points to two UTF8 entries: name and description
			// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4.6
			// the descriptor points either to a method, whose UTF8 should begin with a (
			// or to a field, which must start with one of the letter specified in:
			// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.3.2-200
			whichNandT := entry.slot
			if whichNandT < 0 || whichNandT >= len(klass.nameAndTypes) {
				return cfe("Name and Type at CP entry #" + strconv.Itoa(j) +
					" points to an invalid entry in CP nameAndTypes")
			}

			nAndTentry := klass.nameAndTypes[whichNandT]
			_, err := fetchUTF8string(klass, nAndTentry.nameIndex)
			if err != nil {
				return cfe("Name and Type at CP entry #" + strconv.Itoa(j) +
					"has a name index that points to an invalid UTF8 entry: " +
					strconv.Itoa(nAndTentry.nameIndex))
			}

			desc, err2 := fetchUTF8string(klass, nAndTentry.descriptorIndex)
			if err2 != nil {
				return cfe("Name and Type at CP entry #" + strconv.Itoa(j) +
					"has a description index that points to an invalid UTF8 entry: " +
					strconv.Itoa(nAndTentry.nameIndex))
			}

			descBytes := []byte(desc)
			c := descBytes[0]
			if !(c == '(' || c == 'B' || c == 'C' || c == 'D' || c == 'F' ||
				c == 'I' || c == 'J' || c == 'L' || c == 'S' || c == 'Z' ||
				c == '[') {
				return cfe("Name and Type at CP entry #" + strconv.Itoa(j) +
					"has an invalid description string: " + desc)
			}
			// CURR: continue format checking other CP entries
		default:
			continue
		}
	}

	return nil
}

func validateFields(klass *parsedClass) error {
	return nil
}
