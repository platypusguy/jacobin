/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"errors"
	"jacobin/log"
	"jacobin/stringPool"
	"strconv"
	"strings"
)

// Performs the format check on a fully parsed class. The requirements are listed
// here: https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.8
// They are:
//  1. must start with 0xCAFEBABE -- this is verified in the parsing, so not done here
//  2. most predefined attributes must be the right length -- verified during parsing
//     However, some additional attribute checking done here in formatCheckClassAttributes()
//  3. class must not be truncated or have extra bytes -- verified during parsing
//  4. CP must fulfill all constraints. This is done in formatCheckConstantPool() below
//  5. Fields must have valid names, classes, and descriptions. Partially done in
//     the parsing, but entirely done in formatCheckFields() below
func formatCheckClass(klass *ParsedClass) error {
	if formatCheckConstantPool(klass) != nil {
		return errors.New("") // whatever error occurs, the user will have been notified
	}

	if formatCheckFields(klass) != nil {
		return errors.New("") // whatever error occurs, the user will have been notified
	}

	if formatCheckClassAttributes(klass) != nil {
		return errors.New("") // whatever error occurs, the user will have been notified
	}

	return formatCheckStructure(klass)
}

// validates that the CP fits all the requirements enumerated in:
// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4
// some of these checks were performed perforce in the parsing. Here, however,
// we verify them all. This is a requirement of all classes loaded in the JVM
// Note that this is *not* part of the larger class verification process.
func formatCheckConstantPool(klass *ParsedClass) error {
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
		case LongConst:
			// there are complex bit patterns that can be enforced for longs, but for the
			// nonce, we'll just make sure that there is an actual value pointed to and
			// that the long is followed in the CP by a dummy entry. Consult:
			// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4.5
			whichLong := entry.slot
			if whichLong < 0 || whichLong >= len(klass.longConsts) {
				return cfe("Long constant at CP entry #" + strconv.Itoa(j) +
					" points to an invalid entry in CP longConsts")
			}

			nextEntry := klass.cpIndex[j+1]
			if nextEntry.entryType != Dummy {
				return cfe("Missing dummy entry after long constant at CP entry#" +
					strconv.Itoa(j))
			}
			j += 1
		case DoubleConst:
			// see the comments on the LongConst. They apply exactly to the following code.
			whichDouble := entry.slot
			if whichDouble < 0 || whichDouble >= len(klass.doubles) {
				return cfe("Double constant at CP entry #" + strconv.Itoa(j) +
					" points to an invalid entry in CP doubless")
			}

			nextEntry := klass.cpIndex[j+1]
			if nextEntry.entryType != Dummy {
				return cfe("Missing dummy entry after double constant at CP entry#" +
					strconv.Itoa(j))
			}
			j += 1
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
		case FieldRef:
			// the requirements are that the class index points to a valid Class entry
			// and the name_and_type index points to a valid NameAndType entry. Consult
			// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4.2
			// Here we just make sure they point to entries of the correct type and that
			// they exist. The pointed-to entries are themselves validated as this loop
			// picks them up going through the CP.
			whichFieldRef := entry.slot
			if whichFieldRef < 0 || whichFieldRef >= len(klass.fieldRefs) {
				return cfe("Field Ref at CP entry #" + strconv.Itoa(j) +
					" points to an invalid entry in CP fieldRefs")
			}
			fieldRef := klass.fieldRefs[whichFieldRef]
			classIndex := fieldRef.classIndex
			class := klass.cpIndex[classIndex]
			if class.entryType != ClassRef ||
				class.slot < 0 || class.slot >= len(klass.classRefs) {
				return cfe("Field Ref at CP entry #" + strconv.Itoa(j) +
					" has a class index that points to an invalid entry in ClassRefs. " +
					strconv.Itoa(classIndex))
			}

			nameAndType := klass.cpIndex[fieldRef.nameAndTypeIndex]
			if nameAndType.entryType != NameAndType ||
				nameAndType.slot < 0 || nameAndType.slot >= len(klass.nameAndTypes) {
				return cfe("Field Ref at CP entry #" + strconv.Itoa(j) +
					" has a nameAndType index that points to an invalid entry in nameAndTypes. " +
					strconv.Itoa(fieldRef.nameAndTypeIndex))
			}
		case MethodRef:
			// the MethodRef must have a class index that points to a Class_info entry
			// which itself must point to a class, not an interface. The MethodRef also has
			// an index to a NameAndType entry. If the name of the latter entry begins with
			// and <, then the name can only be <init>. Consult:
			// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4.2
			whichMethodRef := entry.slot
			methodRef := klass.methodRefs[whichMethodRef]

			classIndex := methodRef.classIndex
			class := klass.cpIndex[classIndex]
			if class.entryType != ClassRef ||
				class.slot < 0 || class.slot >= len(klass.classRefs) {
				return cfe("Method Ref at CP entry #" + strconv.Itoa(j) +
					" holds an invalid class index: " +
					strconv.Itoa(class.slot))
			}

			nAndTIndex := methodRef.nameAndTypeIndex
			nAndT := klass.cpIndex[nAndTIndex]
			if nAndT.entryType != NameAndType ||
				nAndT.slot < 0 || nAndT.slot >= len(klass.nameAndTypes) {
				return cfe("Method Ref at CP entry #" + strconv.Itoa(j) +
					" holds an invalid NameAndType index: " +
					strconv.Itoa(nAndT.slot))
			}

			nAndTentry := klass.nameAndTypes[nAndT.slot]
			methodNameIndex := nAndTentry.nameIndex
			name, err := FetchUTF8string(klass, methodNameIndex)
			if err != nil {
				return cfe("Method Ref (at CP entry #" + strconv.Itoa(j) +
					") has a Name and Type entry does not have a name that is a valid UTF8 entry")
			}

			nameBytes := []byte(name)
			if nameBytes[0] == '<' && name != "<init>" {
				return cfe("Method Ref at CP entry #" + strconv.Itoa(j) +
					" holds an NameAndType index to an entry with an invalid method name " +
					name)
			}
		case Interface:
			// the Interface entries are almost identical to the class entries (see above),
			// except that the class index must point to an interface class, and the requirement
			// re naming < and <init> does not apply.
			whichInterface := entry.slot
			interfaceRef := klass.interfaceRefs[whichInterface]

			classIndex := interfaceRef.classIndex
			class := klass.cpIndex[classIndex]
			if class.entryType != ClassRef ||
				class.slot < 0 || class.slot >= len(klass.classRefs) {
				return cfe("Interface Ref at CP entry #" + strconv.Itoa(j) +
					" holds an invalid class index: " + strconv.Itoa(class.slot))
			}

			clRef := klass.classRefs[class.slot]
			clName := stringPool.GetStringPointer(clRef)
			if clName == nil {
				return cfe("Interface Ref at CP entry #" + strconv.Itoa(j) +
					// " holds an invalid UTF8 index to the interface name: " +
					"holds an invalid stringPool index for interface: " +
					strconv.FormatUint(uint64(clRef), 10))
			}

			/* TODO: REVISIT: with java.lang.String the following code works OK
			   with the three interfaces defined in klass.interfaces[], but Iterable
			   is not among those classes and yet it's got a interfaceRef CP entry.
			   So, not presently sure how you validate that the interfaceRef CP entry
			   points to an interface. So for the nonce, the following code is commented out.

			   // now that we have the UTF8 index for the interface reference,
			   // check whether it's in our list of interfaces for this class.
			   matchesInterface := false
			   for i := range klass.interfaces {
			   	if klass.interfaces[i] == utfIndex {
			   		matchesInterface = true
			   	}
			   }

			   if ! matchesInterface {
			   	return cfe("Interface Ref at CP entry #"+ strconv.Itoa(j) +
			   		" does not match to any interface in this class.")
			   }
			*/

			nAndTIndex := interfaceRef.nameAndTypeIndex
			nAndT := klass.cpIndex[nAndTIndex]
			if nAndT.entryType != NameAndType ||
				nAndT.slot < 0 || nAndT.slot >= len(klass.nameAndTypes) {
				return cfe("Method Ref at CP entry #" + strconv.Itoa(j) +
					" holds an invalid NameAndType index: " +
					strconv.Itoa(nAndT.slot))
			}
		case NameAndType:
			// a NameAndType entry points to two UTF8 entries: name and description. Consult
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
			_, err := FetchUTF8string(klass, nAndTentry.nameIndex)
			if err != nil {
				return cfe("Name and Type at CP entry #" + strconv.Itoa(j) +
					" has a name index that points to an invalid UTF8 entry: " +
					strconv.Itoa(nAndTentry.nameIndex))
			}

			desc, err2 := FetchUTF8string(klass, nAndTentry.descriptorIndex)
			if err2 != nil {
				return cfe("Name and Type at CP entry #" + strconv.Itoa(j) +
					" has a description index that points to an invalid UTF8 entry: " +
					strconv.Itoa(nAndTentry.nameIndex))
			}

			err = validateFieldDesc(desc)
			if err != nil {
				return cfe("Name and Type at CP entry #" + strconv.Itoa(j) +
					" has an invalid description string: " + desc)
			}
		case MethodHandle:
			// Method handles have complex validation logic. It's entirely enforced here. See:
			// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4.8
			// CONSTANT_MethodHandle_info {
			//    u1 tag;
			//    u1 reference_kind;
			//    u2 reference_index; }
			whichMethHandle := entry.slot
			mhe := klass.methodHandles[whichMethHandle]
			refKind := mhe.referenceKind
			if refKind < 1 || refKind > 9 {
				return cfe("MethodHandle at CP entry #" + strconv.Itoa(j) +
					" has an invalid reference kind: " + strconv.Itoa(refKind))
			}
			refIndex := mhe.referenceIndex

			switch refKind {
			// if refKind is 1-4, the reference_index must point to a fieldRef
			case 1, 2, 3, 4:
				if klass.cpIndex[refIndex].entryType != FieldRef {
					return cfe("MethodHandle at CP entry #" + strconv.Itoa(j) +
						" has an reference kind between 1-4 ( " + strconv.Itoa(refKind) +
						") which does not point to a FieldRef")
				}
			// if refKind is 5 or 8, the reference_index must point to a methodRef
			case 5, 8:
				if klass.cpIndex[refIndex].entryType != MethodRef {
					return cfe("MethodHandle at CP entry #" + strconv.Itoa(j) +
						" has an reference kind between of 5 or 8 ( " + strconv.Itoa(refKind) +
						") which does not point to a MethodRef")
				}
			case 6, 7:
				// if refKind is 6 or 7, the reference_index must point to a methodRef or if the
				// class version # is >= 52, it can point to an Interface. To make the logic readable,
				// we test for the positive here, rather than the negative as in the other cases
				if klass.cpIndex[refIndex].entryType == MethodRef ||
					(klass.javaVersion >= 52 && klass.cpIndex[refIndex].entryType == Interface) {
					break
				} else {
					return cfe("MethodHandle at CP entry #" + strconv.Itoa(j) +
						" has an reference kind between of 6 or 7 ( " + strconv.Itoa(refKind) +
						") which does not point to a MethodRef or in Java version 52 or later " +
						"does not point to an Interface.")
				}
			case 9:
				if klass.cpIndex[refIndex].entryType != Interface {
					return cfe("MethodHandle at CP entry #" + strconv.Itoa(j) +
						" has an reference kind  of 9 which does not point to an interface")
				}
			}

			// get the class name pointed to by the MethodRef pointed to by the MethodHandle
			var methodName string
			var err error
			if klass.cpIndex[refIndex].entryType == MethodRef {
				methodName, _, _, err = resolveCPmethodRef(refIndex, klass)
				if err != nil {
					return errors.New("") // the error messsage is already displayed
				}
			}

			// if the reference_kind is 5-7 the name of the method pointed to
			// by the nameAndType entry in the method handle cannot be <init> or <clinit>
			if refKind >= 5 && refKind <= 7 && klass.cpIndex[refIndex].entryType == MethodRef {
				methRefIndex := klass.cpIndex[refIndex].slot
				if methRefIndex < 0 || methRefIndex >= len(klass.methodRefs) {
					return cfe("Reference index for MethodHandle at CP entry #" + strconv.Itoa(j) +
						" points to an invalid MethodRef: " + strconv.Itoa(methRefIndex))
				}

				if methodName == "<init>" || methodName == "<clinit>" {
					return cfe("Invalid class name for MethodHandle at CP entry #" + strconv.Itoa(j) +
						" : " + methodName)
				}
			}
			// The following code is commented out b/c it was emitting errors when Jacobin moved from
			// JDK 11 to JDK 17. The cause of these errors is unclear--whether something changed in
			// JDK 17 w.r.t MethodHandles with refKind = 8. Issue #JACOBIN-183 is the reference for this.
			// else if refKind == 8 {
			//	if methodName != "<init>" {
			//		return cfe("Class name for MethodHandle at CP entry #" + strconv.Itoa(j) +
			//			" should be <init>, but is: " + methodName)
			//	}
			// }

			_ = log.Log("ClassName in MethodRef of MethodHandle at CP entry #"+strconv.Itoa(j)+
				" is:"+methodName, log.FINEST)
		case MethodType:
			// Method types consist of an integer pointing to a CP entry that's a UTF8 description
			// of the method type, which appears to require an initial opening parenthesis. See
			// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4.9
			whichMethType := entry.slot
			mte := klass.methodTypes[whichMethType]
			utf8 := klass.cpIndex[mte]
			if utf8.entryType != UTF8 || utf8.slot < 0 || utf8.slot > len(klass.utf8Refs)-1 {
				return cfe("MethodType at CP entry #" + strconv.Itoa(j) +
					" has an invalid description index: " + strconv.Itoa(utf8.slot))
			}
			methType := klass.utf8Refs[utf8.slot]
			if !strings.HasPrefix(methType.content, "(") {
				return cfe("MethodType at CP entry #" + strconv.Itoa(j) +
					" does not point to a type that starts with an open parenthesis. Got: " +
					methType.content)
			}
		case Dynamic:
			// Like InvokeDynamic, Dynamic is a unique kind of entry. The first field,
			// bootstrapIndex, must be a "valid index into the bootstrap_methods array
			// of the bootstrap method table of this this class file" (specified in §4.7.23).
			// The document spec for InvokeDynamic entries is found at:
			// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4.10
			// The second field is a nameAndType record describing the bootstrap method.
			// Here we just make sure, the field points to the right kind of entry and that
			// the descriptor in the nameAndType points to a field.
			whichDyn := entry.slot
			if whichDyn >= len(klass.dynamics) {
				return cfe("The dynamic entry at CP[" + strconv.Itoa(j) + "] " +
					"points to a non-existent dynamic slot: " + strconv.Itoa(entry.slot))
			}
			dyn := klass.dynamics[whichDyn]

			bootstrap := dyn.bootstrapIndex
			if bootstrap >= klass.bootstrapCount {
				return cfe("The bootstrap index in dynamic at CP[" + strconv.Itoa(j) +
					"] is invalid: " + strconv.Itoa(bootstrap))
			}

			// just trying to access it to make sure it's actually there and accessible.
			bse := klass.bootstraps[bootstrap]
			if !(bse.methodRef > 0) {
				return cfe("Invalid methodRef in bootstrap method[" + strconv.Itoa(bootstrap) + "]")
			}

			nAndT := dyn.nameAndType
			if nAndT < 1 || nAndT > len(klass.cpIndex)-1 {
				return cfe("The entry number into klass.dynamics[] at CP entry #" +
					strconv.Itoa(j) + " is invalid: " + strconv.Itoa(nAndT))
			}
			if klass.cpIndex[nAndT].entryType != NameAndType {
				return cfe("NameAndType index at CP entry #" + strconv.Itoa(j) +
					" (dynamic) points to an entry that's not NameAndType: " +
					strconv.Itoa(klass.cpIndex[nAndT].entryType))
			}

			natSlot := klass.cpIndex[nAndT].slot
			nat := klass.nameAndTypes[natSlot] // gets the actual nameAndType entry
			desc, err := FetchUTF8string(klass, nat.descriptorIndex)
			if err != nil {
				return cfe("Descriptor in nameAndType entry of dynamic CP entry #" +
					strconv.Itoa(j) + " is invalid: " + strconv.Itoa(nat.descriptorIndex))
			}

			if validateFieldDesc(desc) != nil {
				return cfe("Descriptor in nameAndType entry of dynamic CP entry #" +
					strconv.Itoa(j) + " is an invalid field descriptor: " + desc)
			}

		case InvokeDynamic:
			// InvokeDynamic is a unique kind of entry. The first field, bootstrapIndex, must be a
			// "valid index into the bootstrap_methods array of the bootstrap method table of this
			// this class file" (specified in §4.7.23). The document spec for InvokeDynamic entries is:
			// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4.10
			// The second field is a nameAndType record describing the bootstrap method.
			// Here we just make sure, the field points to the right kind of entry. That entry
			// will be checked later/earlier in this format check.
			whichInvDyn := entry.slot
			if whichInvDyn >= len(klass.invokeDynamics) {
				return cfe("The invokeDynamic entry at CP[" + strconv.Itoa(j) + "] " +
					"points to a non-existent invokeDynamic slot: " + strconv.Itoa(entry.slot))
			}
			invDyn := klass.invokeDynamics[whichInvDyn]

			bootstrap := invDyn.bootstrapIndex
			if bootstrap >= klass.bootstrapCount {
				return cfe("The bootstrap index in InvokeDynamic at CP[" + strconv.Itoa(j) +
					"] is invalid: " + strconv.Itoa(bootstrap))
			}

			// just trying to access it to make sure it's actually there and accessible.
			bse := klass.bootstraps[bootstrap]
			if !(bse.methodRef > 0) {
				return cfe("Invalid methodRef in bootstrap method[" + strconv.Itoa(bootstrap) + "]")
			}

			nAndTslot := invDyn.nameAndType
			if nAndTslot < 1 || nAndTslot > len(klass.cpIndex)-1 {
				return cfe("The entry number into klass.InvokeDynamics[] at CP entry #" +
					strconv.Itoa(j) + " is invalid: " + strconv.Itoa(nAndTslot))
			}
			if klass.cpIndex[nAndTslot].entryType != NameAndType {
				return cfe("NameAndType index at CP entry #" + strconv.Itoa(j) +
					" (InvokeDynamic) points to an entry that's not NameAndType: " +
					strconv.Itoa(klass.cpIndex[nAndTslot].entryType))
			}

			natSlot := klass.cpIndex[nAndTslot].slot
			nat := klass.nameAndTypes[natSlot] // gets the actual nameAndType entry
			desc, err := FetchUTF8string(klass, nat.descriptorIndex)
			if err != nil {
				return cfe("Descriptor in nameAndType entry of dynamic CP entry #" +
					strconv.Itoa(j) + " is invalid: " + strconv.Itoa(nat.descriptorIndex))
			}

			if validateMethodDesc(desc) != nil {
				return cfe("Descriptor in nameAndType entry of dynamic CP entry #" +
					strconv.Itoa(j) + " is an invalid method descriptor: " + desc)
			}
		case Module:
			// if there's a module entry, the module name has already been fetched and
			// placed into klass.moduleName. So, here we verify this module name rather
			// than the CP entry that got it. We also check access permissions, as required
			// in: https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4.11
			// Note: the test for minimum Java 9 version and the limit of at most one
			// Module entry is enforced in the original CP parsing (see cpParser.go)
			if !klass.classIsModule {
				return cfe("Module CP entry must appear only in class with ACC_MODULE set.")
			}
			if checkModuleName(klass.moduleName) != nil {
				return errors.New("") // the error message will already have been displayed
			}
		case Package:
			// if there's a package entry, the package name has already been fetched and
			// placed into klass.packageName. So, here we verify this package name rather
			// than the CP entry that got it. We also check access permissions, as required
			// in: https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4.12
			// Note: the test for minimum Java 9 version and the limit of at most one
			// Package entry is enforced in the original CP parsing (see cpParser.go)
			if !klass.classIsModule {
				return cfe("Package CP entry must appear only in class with ACC_MODULE set.")
			}

			// packages have the same restrictions on the names as modules.
			if checkPackageName(klass.packageName) != nil {
				return errors.New("") // the error message will already have been displayed
			}
		default:
			continue
		}
	}

	return nil
}

// field entries consist of two string indexes, one of which points to the name, the other
// to a string containing a description of the type. Here we grab the strings and check that
// they fulfill the requirements: name doesn't start with a digit or contain a space, and the
// type begins with one of the required letters/symbols
func formatCheckFields(klass *ParsedClass) error {
	for i, f := range klass.fields {
		// f.name points to a UTF8 entry in klass.utf8refs, so check it's in a valid range
		if f.name < 0 || f.name >= len(klass.utf8Refs) {
			return cfe("Invalid index to UTF8 string for field name in field #" + strconv.Itoa(i))
		}
		fName := klass.utf8Refs[f.name].content

		// f.description points to a UTF8 entry in klass.utf8refs, so check it's in a valid range
		if f.description < 0 || f.description >= len(klass.utf8Refs) {
			return cfe("Invalid index for UTF8 string containing description of field " + fName)
		}
		fDesc := klass.utf8Refs[f.description].content

		fNameBytes := []byte(fName)
		if fNameBytes[0] >= '0' && fNameBytes[0] <= '9' {
			return cfe("Invalid field name in format check (starts with a digit): " + fName)
		}

		// check that there is no leading, trailing, or embedded whitespace
		for _, c := range fNameBytes {
			switch c {
			case
				'\u0009', // horizontal tab
				'\u000A', // line feed
				'\u000B', // vertical tab
				'\u000C', // form feed
				'\u000D', // carriage return
				'\u0020', // space
				'\u0085', // next line
				'\u00A0': // no-break space
				return cfe("Invalid field name in format check (contains whitespace): " + fName)
			default:
				continue
			}
		}

		if validateFieldDesc(fDesc) != nil {
			return cfe("Field " + fName + " has an invalid description string: " + fDesc)
		}
	}
	return nil
}

// certain descriptions and type strings must start with one of the letters shown here.
// See: https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-FieldType
func validateFieldDesc(desc string) error {
	if len(desc) < 1 {
		return errors.New("invalid")
	}

	descBytes := []byte(desc)
	c := descBytes[0]
	if !(c == '(' || c == 'B' || c == 'C' || c == 'D' || c == 'F' ||
		c == 'I' || c == 'J' || c == 'L' || c == 'S' || c == 'Z' ||
		c == '[') {
		return errors.New("invalid")
	}
	return nil
}

// Method descriptors list the parameters and the return type of a method. The symbols
// for these are identical to field descriptors see alidateFieldDesc()with the addition
// of V for void. https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.3.3
func validateMethodDesc(desc string) error {
	if len(desc) < 1 {
		return errors.New("invalid")
	}

	descBytes := []byte(desc)
	c := descBytes[0]
	if !(c == '(' || c == 'B' || c == 'C' || c == 'D' || c == 'F' ||
		c == 'I' || c == 'J' || c == 'L' || c == 'S' || c == 'Z' ||
		c == '[' || c == 'V') {
		return errors.New("invalid")
	}
	return nil
}

// validates the unqualified names of fields and methods. "Unqualified" is a term of art, see:
// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.2.2
// the 'method' parameter indicates whether the string is the name of a method (which would
// necessitate additional checks. On error, returns false
func validateUnqualifiedName(name string, method bool) bool {
	if len(name) == 0 {
		return false
	}
	bytes := []byte(name)
	for _, v := range bytes { // check there are no embedded . ; [ / (
		if v == '.' || v == ';' || v == '[' || v == '/' || v == '(' {
			return false
		}
	}

	// only the methods <init> and <clinit> can contain a < or a >
	if method {
		if name != "<init>" && name != "<clinit>" {
			for _, v := range bytes {
				if v == '<' || v == '>' {
					return false
				}
			}
		}
	}
	return true
}

// module names have multiple restrictions. Some UTF8 code points are disallowed. We don't
// check for those here, but certain characters are disallowed. Those are explained
// here: https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.2.3
// see checkPackageName() for the same check (but different error messages)
func checkModuleName(name string) error {
	if name == "" {
		return cfe("Expected a module/package name, but none was found.")
	}

	bArr := []byte(name)
	if bArr[0] == '@' || bArr[0] == ':' { // a @ or : must be escaped, so can't start name
		return cfe("Module/Package name " + name + " contains an illegal character")
	}

	invalidName := false
	for i := 1; i < len(bArr); i++ {
		switch bArr[i] {
		case '@', ':':
			if bArr[i-1] != '\\' {
				invalidName = true
			}
		case '\\':
			if i+1 >= len(bArr) { // name cannot end on a \
				invalidName = true
				break
			}
			// if a \ is encountered it can only escape a @, :, or \
			// if this is the case, we skip the escaped char, if not, it's an error
			if bArr[i+1] == '@' || bArr[i+1] == ':' || bArr[i+1] == '\\' {
				i += 1
			} else {
				invalidName = true
			}
		}
		if invalidName {
			return cfe("Module name " + name + " contains an illegal character")
		}
	}
	return nil
}

// package names have multiple restrictions. Some UTF8 code points are disallowed. We don't
// check for those here, but certain characters are disallowed. Those are explained
// here: https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.2.3
// see checkModuleName() for the same check (but different error messages)
func checkPackageName(name string) error {
	if name == "" {
		return cfe("Expected a package name, but none was found.")
	}

	bArr := []byte(name)
	if bArr[0] == '@' || bArr[0] == ':' { // a @ or : must be escaped, so can't start name
		return cfe("Package name " + name + " contains an illegal character")
	}

	invalidName := false
	for i := 1; i < len(bArr); i++ {
		switch bArr[i] {
		case '@', ':':
			if bArr[i-1] != '\\' {
				invalidName = true
			}
		case '\\':
			if i+1 >= len(bArr) { // name cannot end on a \
				invalidName = true
				break
			}
			// if a \ is encountered it can only escape a @, :, or \
			// if this is the case, we skip the escaped char, if not, it's an error
			if bArr[i+1] == '@' || bArr[i+1] == ':' || bArr[i+1] == '\\' {
				i += 1
			} else {
				invalidName = true
			}
		}
		if invalidName {
			return cfe("Package name " + name + " contains an illegal character")
		}
	}
	return nil
}

// Format checks various class attributes. Technically speaking, format checking is
// only supposed to check the length of attributes. (This is done when the attributes
// are initially parsed.) This adds a little more checking.
func formatCheckClassAttributes(klass *ParsedClass) error {

	// enforce basic checks of bootstrap entries (which are used by invokedynamic)
	if len(klass.bootstraps) > 0 {
		for i := 0; i < len(klass.bootstraps); i++ {
			bsm := klass.bootstraps[i]
			if klass.cpIndex[bsm.methodRef].entryType != MethodHandle {
				return cfe("MethodRef in bootstrapMethod[" + strconv.Itoa(i) + "] in class " +
					klass.className + "should but does not point to a MethodHandle")
			}

			if len(bsm.args) > 0 {
				for j := 0; j < len(bsm.args); j++ {
					if !validateItemIsLodable(klass, bsm.args[j]) {
						return cfe("Bootstrap method argument[" + strconv.Itoa(j) + "] in class " +
							klass.className + " bootstrap method #[" + strconv.Itoa(i) + "] " +
							"should be but is not a loadable constant")
					}
				}
			}
		}
	}
	return nil
}

// Certain types of items are loadable. This checks that an entry into the CP
// does in fact point to a loadable item. Returns false if not or on any error.
// See Table 4.4C: https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4
func validateItemIsLodable(klass *ParsedClass, arg int) bool {
	if arg < 1 || arg >= len(klass.cpIndex) {
		return false
	}

	t := klass.cpIndex[arg].entryType
	if t != IntConst && t != FloatConst && t != LongConst && t != DoubleConst &&
		t != ClassRef && t != StringConst && t != MethodHandle && t != MethodType &&
		t != Dynamic {
		return false
	}

	return true
}

// format checks of structural elements outside of CP and fields. For example,
// checking that a count field holds the correct number, etc.
func formatCheckStructure(klass *ParsedClass) error {
	if klass.cpCount != len(klass.cpIndex) {
		return cfe("CP count: " + strconv.Itoa(klass.cpCount) +
			" is not equal to actual size of CP: " + strconv.Itoa(len(klass.cpIndex)))
	}

	if klass.interfaceCount != len(klass.interfaces) {
		return cfe("Expected " + strconv.Itoa(klass.interfaceCount) + " interfaces. Got: " +
			strconv.Itoa(len(klass.interfaces)))
	}

	if klass.methodCount != len(klass.methods) {
		return cfe("Expected " + strconv.Itoa(klass.methodCount) + " methods. Got: " +
			strconv.Itoa(len(klass.methods)))
	}

	if klass.attribCount != len(klass.attributes) {
		return cfe("Expected " + strconv.Itoa(klass.attribCount) + " class attributes. Got: " +
			strconv.Itoa(len(klass.attributes)))
	}

	if klass.bootstrapCount != len(klass.bootstraps) {
		return cfe("Expected " + strconv.Itoa(klass.bootstrapCount) + " bootstrap methods. Got: " +
			strconv.Itoa(len(klass.bootstraps)))
	}

	return nil
}
