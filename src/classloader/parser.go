/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"encoding/hex"
	"errors"
	"fmt"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/stringPool"

	// "jacobin/object"
	"jacobin/types"
	"os"
	"strconv"
)

// reads in a class file, parses it, and puts the values into the fields of the
// class that will be loaded into the classloader. Basic verification performed.
// receives the rawBytes of the class that were previously read in
//
// ClassFormatError - if the parser finds anything unexpected
func parse(rawBytes []byte) (ParsedClass, error) {

	// the parsed class as we'll give it to the classloader
	var pClass = ParsedClass{}

	err := parseMagicNumber(rawBytes)
	if err != nil {
		return pClass, err
	}

	err = parseJavaVersionNumber(rawBytes, &pClass)
	if err != nil {
		return pClass, err
	}

	err = getConstantPoolCount(rawBytes, &pClass)
	if err != nil {
		return pClass, err
	}

	pos, err := parseConstantPool(rawBytes, &pClass)
	if err != nil || pos < 10 {
		return pClass, err
	}

	pos, err = parseAccessFlags(rawBytes, pos, &pClass)
	if err != nil {
		return pClass, err
	}

	pos, err = parseClassName(rawBytes, pos, &pClass)
	if err != nil {
		return pClass, err
	}

	pos, err = parseSuperClassName(rawBytes, pos, &pClass)
	if err != nil {
		return pClass, err
	}

	pos, err = parseInterfaceCount(rawBytes, pos, &pClass)
	if err != nil {
		return pClass, err
	}

	if pClass.interfaceCount > 0 {
		pos, err = parseInterfaces(rawBytes, pos, &pClass)
		if err != nil {
			return pClass, err
		}
	}

	pos, err = parseFieldCount(rawBytes, pos, &pClass)
	if err != nil {
		return pClass, err
	}

	if pClass.fieldCount > 0 {
		pos, err = parseFields(rawBytes, pos, &pClass)
		if err != nil {
			return pClass, err
		}
	}

	pos, err = parseMethodCount(rawBytes, pos, &pClass)
	if err != nil {
		return pClass, err
	}

	if pClass.methodCount > 0 {
		pos, err = parseMethods(rawBytes, pos, &pClass)
		if err != nil {
			return pClass, err
		}
	}

	pos, err = parseClassAttributeCount(rawBytes, pos, &pClass)
	if err != nil {
		return pClass, err
	}

	if pClass.attribCount > 0 {
		pos, err = parseClassAttributes(rawBytes, pos, &pClass)
	}
	if err != nil {
		return pClass, err
	}

	if pos != len(rawBytes)-1 {
		return pClass, cfe("Unexpected bytes found at end of class file: " + pClass.className)
	}
	return pClass, nil
}

// all bytecode files start with 0xCAFEBABE ( it was the 90s!)
// this checks for that.
func parseMagicNumber(bytes []byte) error {
	if len(bytes) < 4 {
		return cfe("invalid magic number")
	} else if (bytes[0] != 0xCA) || (bytes[1] != 0xFE) || (bytes[2] != 0xBA) || (bytes[3] != 0xBE) {
		return cfe("invalid magic number")
	} else {
		return nil
	}
}

// get the Java version number used in creating this class file. If it's higher than the
// version Jacobin presently supports, report an error.
func parseJavaVersionNumber(bytes []byte, klass *ParsedClass) error {
	version, err := intFrom2Bytes(bytes, 6)
	if err != nil {
		return err
	}

	if version > globals.GetGlobalRef().MaxJavaVersionRaw {
		errMsg := "Jacobin supports only Java versions through Java " +
			strconv.Itoa(globals.GetGlobalRef().MaxJavaVersion)
		return cfe(errMsg)
	}

	klass.javaVersion = version
	_ = log.Log("Java version: "+strconv.Itoa(version), log.FINEST)
	return nil
}

// get the number of entries in the constant pool. This number will
// be used later on to verify that the number of entries we fetch is
// correct. Note that this number is technically 1 greater than the
// number of actual entries, because the first entry in the constant
// pool is an empty placeholder, rather than an actual entry.
func getConstantPoolCount(bytes []byte, klass *ParsedClass) error {
	cpEntryCount, err := intFrom2Bytes(bytes, 8)
	if err != nil || cpEntryCount <= 2 {
		return cfe("Invalid number of entries in constant pool: " +
			strconv.Itoa(cpEntryCount))
	}

	klass.cpCount = cpEntryCount
	_ = log.Log("Number of CP entries: "+strconv.Itoa(cpEntryCount), log.FINEST)
	return nil
}

// decode the meaning of the class access flags and set the various getters
// in the class. FromTable 4.1-B in the spec:
// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.1-200-E.1
func parseAccessFlags(bytes []byte, loc int, klass *ParsedClass) (int, error) {
	pos := loc
	accessFlags, err := intFrom2Bytes(bytes, pos+1)
	pos += 2
	if err != nil {
		return pos, cfe("Invalid get of class access flags")
	} else {
		klass.accessFlags = accessFlags
		if accessFlags&0x0001 > 0 {
			klass.classIsPublic = true
		}
		if accessFlags&0x0010 > 0 {
			klass.classIsFinal = true
		}
		if accessFlags&0x0020 > 0 {
			klass.classIsSuper = true
		}
		if accessFlags&0x0200 > 0 {
			klass.classIsInterface = true
		}
		if accessFlags&0x0400 > 0 {
			klass.classIsAbstract = true
		}
		if accessFlags&0x1000 > 0 {
			klass.classIsSynthetic = true
		} // is generated by the JVM, is not in the program
		if accessFlags&0x2000 > 0 {
			klass.classIsAnnotation = true
		}
		if accessFlags&0x4000 > 0 {
			klass.classIsEnum = true
		}
		if accessFlags&0x8000 > 0 {
			klass.classIsModule = true
		}
		_ = log.Log("Access flags: 0x"+hex.EncodeToString(bytes[pos-1:pos+1]), log.FINEST)

		if log.Level == log.FINEST {
			if klass.classIsPublic {
				_, _ = fmt.Fprintf(os.Stderr, "access: public\n")
			}
			if klass.classIsFinal {
				_, _ = fmt.Fprintf(os.Stderr, "access: final\n")
			}
			if klass.classIsSuper {
				_, _ = fmt.Fprintf(os.Stderr, "access: super\n")
			}
			if klass.classIsInterface {
				_, _ = fmt.Fprintf(os.Stderr, "access: interface\n")
			}
			if klass.classIsAbstract {
				_, _ = fmt.Fprintf(os.Stderr, "access: abstract\n")
			}
			if klass.classIsSynthetic {
				_, _ = fmt.Fprintf(os.Stderr, "access: synthetic\n")
			}
			if klass.classIsAnnotation {
				_, _ = fmt.Fprintf(os.Stderr, "access: annotation\n")
			}
			if klass.classIsEnum {
				_, _ = fmt.Fprintf(os.Stderr, "access: enum\n")
			}
			if klass.classIsModule {
				_, _ = fmt.Fprintf(os.Stderr, "access: module\n")
			}
		}
		return pos, nil
	}
}

// The value for this item points to a CP entry of type Class_info. (See:
// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.4.1 )
// In turn, that entry points to the UTF-8 name of the class. This name includes
// the package name as a path, but not the extension of .class. So, for example,
// ParsePosition.class in the core Java string library has a class name of:
// java/text/ParsePosition
func parseClassName(bytes []byte, loc int, klass *ParsedClass) (int, error) {
	pos := loc
	index, err := intFrom2Bytes(bytes, pos+1)
	var classNameIndex int
	pos += 2
	if err != nil {
		return pos, cfe("error obtaining index for class name")
	}

	if index < 1 || index > (len(klass.cpIndex)-1) {
		return pos, cfe("invalid index into CP for class name: " +
			strconv.Itoa(index))
	}

	pointedToClassRef := klass.cpIndex[index]
	if pointedToClassRef.entryType != ClassRef {
		return pos, cfe("invalid entry for class name")
	}

	// the entry pointed to by pointedToClassRef holds an index to
	// a UTF-8 string that holds the class name
	classNameIndex = klass.classRefs[pointedToClassRef.slot]
	className, err := FetchUTF8string(klass, classNameIndex)
	if err != nil {
		return pos, errors.New("") // the error msg has already been show to user
	}

	_ = log.Log("class name: "+className, log.FINEST)

	if len(klass.className) > 0 {
		return pos, cfe("Class appears to have two names: " + klass.className + " and: " + className)
	}

	klass.className = className
	klass.classNameIndex = stringPool.GetStringIndex(&className)
	return pos, nil
}

// Get the name of the superclass. The logic is identical to that of parseClassName()
// All classes, except java/lang/Object have superclasses.
func parseSuperClassName(bytes []byte, loc int, klass *ParsedClass) (int, error) {
	pos := loc
	index, err := intFrom2Bytes(bytes, pos+1)
	var classNameIndex int
	pos += 2
	if err != nil {
		return pos, cfe("error obtaining index for superclass name")
	}

	if index == 0 {
		if klass.className != "java/lang/Object" {
			return pos, cfe("invaild index for superclass name. Got: 0," +
				" but class is not java/lang/Object")
		} else {
			_ = log.Log("superclass name: [none]", log.FINEST)
			klass.superClass = ""
			return pos, nil
		}
	}

	if index < 1 || index > (len(klass.cpIndex)-1) {
		return pos, cfe("invalid index into CP for superclass name")
	}

	pointedToClassRef := klass.cpIndex[index]
	if pointedToClassRef.entryType != ClassRef {
		return pos, cfe("invalid entry for superclass name")
	}

	// the entry pointed to by pointedToClassRef holds an index to
	// a UTF-8 string that holds the class name
	classNameIndex = klass.classRefs[pointedToClassRef.slot]

	superClassName, err := FetchUTF8string(klass, classNameIndex)
	if err != nil {
		return pos, errors.New("") // error has already been reported to user
	}

	if superClassName == "" { // only Object.class can have an empty superclass and it's handled above
		return pos, cfe("invalid empty string for superclass name")
	}

	_ = log.Log("superclass name: "+superClassName, log.FINEST)
	if len(klass.superClass) > 0 {
		return pos, cfe("Class can only have 1 superclass, found two: " + klass.superClass + " and: " + superClassName)
	}

	klass.superClass = superClassName
	klass.superClassIndex = stringPool.GetStringIndex(&superClassName)
	return pos, nil
}

// Get the count of the number of interfaces this class implements
func parseInterfaceCount(bytes []byte, loc int, klass *ParsedClass) (int, error) {
	pos := loc
	interfaceCount, err := intFrom2Bytes(bytes, pos+1)
	pos += 2
	if err != nil {
		return pos, cfe("Invalid fetch of interface count")
	}

	_ = log.Log("interface count: "+strconv.Itoa(interfaceCount), log.FINEST)
	klass.interfaceCount = interfaceCount
	return pos, nil
}

// these are actually interface references, simply indexes into the CP that point to
// class name entries, which in turn point to the UTF-8 string holding the name of the
// interface class.
func parseInterfaces(bytes []byte, loc int, klass *ParsedClass) (int, error) {
	pos := loc
	for i := 0; i < klass.interfaceCount; i += 1 {
		interfaceIndex, err := intFrom2Bytes(bytes, pos+1)
		pos += 2
		if err != nil {
			return pos, cfe("Invalid fetch of interface index")
		}

		if interfaceIndex < 1 || interfaceIndex > klass.cpCount-1 {
			return pos, cfe("Interface index is out of range: " + strconv.Itoa(interfaceIndex))
		}

		// get the entry in the CP that the interface index points to,
		// which is a class reference entry that then points to a UTF-8 entry
		classref := klass.cpIndex[interfaceIndex]
		if classref.entryType != ClassRef {
			return pos, cfe("Interface index does not point to a class type. Got: " +
				strconv.Itoa(classref.entryType))
		}

		// get the class entry from classRefs slice
		classEntry := klass.classRefs[classref.slot]

		// use the class entry's index field to look up the UTF-8 string
		interfaceName, err := FetchUTF8string(klass, classEntry)
		if err != nil {
			return pos, errors.New("") // error msg has already been shown
		}

		_ = log.Log("Interface class: "+interfaceName, log.FINEST)

		// klass.interfaces is a slice that holds the index into utf8Refs for
		// each of the interface class names. This avoids duplicating the name
		// that's already in the CP and it allows the classloader to get the
		// interface name in a single dereference.
		klass.interfaces = append(klass.interfaces, klass.cpIndex[classEntry].slot)
	}
	return pos, nil
}

// Get the number of fields in this class
func parseFieldCount(bytes []byte, loc int, klass *ParsedClass) (int, error) {
	pos := loc
	fieldCount, err := intFrom2Bytes(bytes, pos+1)
	pos += 2
	if err != nil {
		return pos, cfe("Invalid fetch of field count")
	}

	_ = log.Log("field count: "+strconv.Itoa(fieldCount), log.FINEST)
	klass.fieldCount = fieldCount
	return pos, nil
}

// parse the fields in a class. The contents of each field is explained here:
// https://docs.oracle.com/javase/specs/jvms/se17/html/jvms-4.html#jvms-4.5
// The layout, per that spec:
// field_info {
//    u2             access_flags;
//    u2             name_index;
//    u2             descriptor_index;
//    u2             attributes_count;
//    attribute_info attributes[attributes_count];
// }

func parseFields(bytes []byte, loc int, klass *ParsedClass) (int, error) {
	pos := loc
	for i := 0; i < klass.fieldCount; i += 1 {
		f := field{}
		f.constValue = nil

		accessFlags, err := intFrom2Bytes(bytes, pos+1)
		pos += 2
		if err != nil {
			return pos, cfe("error retrieving access flags for field " + strconv.Itoa(i))
		}
		f.accessFlags = accessFlags
		if (accessFlags & 0b1000) == 8 {
			f.isStatic = true
		}

		nameIndex, err := intFrom2Bytes(bytes, pos+1)
		pos += 2
		if err != nil || nameIndex < 1 || nameIndex > klass.cpCount-1 {
			return pos, cfe("error retrieving name index for field")
		}

		f.name, err = fetchUTF8slot(klass, nameIndex)
		if err != nil {
			return pos, cfe("error fetching UTF-8 string for name of field")
		}

		descIndex, err := intFrom2Bytes(bytes, pos+1)
		pos += 2
		if err != nil || descIndex < 1 || descIndex > klass.cpCount-1 {
			return pos, cfe("error retrieving description index for field: " +
				klass.utf8Refs[f.name].content)
		}
		f.description, err = fetchUTF8slot(klass, descIndex)
		if err != nil {
			return pos, cfe("error retrieving UTF8 slot for description of field: " +
				klass.utf8Refs[f.name].content)
		}

		attrCount, err := intFrom2Bytes(bytes, pos+1)
		pos += 2
		if err != nil {
			return pos, cfe("error retrieving attribute count for field: " +
				klass.utf8Refs[f.name].content)
		}

		for j := 0; j < attrCount; j++ {
			attribute, k, err := fetchAttribute(klass, bytes, pos)
			if err != nil {
				return pos, errors.New("") // error message will already have been displayed
			}
			attrName := klass.utf8Refs[attribute.attrName].content
			// if the attribute is a constant value (for initializing the field)
			// then stick the value into the field struct. That value is a pointer
			// into the CP and its value must be converted based on the type of
			// field we're dealing with (shown in the desc data item)
			if attrName == "ConstantValue" {
				desc := klass.utf8Refs[f.description].content
				switch desc {
				case types.Ref, types.Bool: // TODO: Find out how to process these
					f.constValue = nil
				case types.Byte: // byte--same logic as for types.Int, only error message is different
					indexIntoCP := int(attribute.attrContent[0])*256 +
						int(attribute.attrContent[1])
					entryInCp := klass.cpIndex[indexIntoCP]
					if entryInCp.entryType != IntConst {
						return pos, cfe("error: wrong type of constant value for byte " +
							klass.utf8Refs[f.name].content)
					}
					f.constValue = klass.intConsts[entryInCp.slot]
				case types.Char: // char--same logic as for types.Int, only error message is different
					indexIntoCP := int(attribute.attrContent[0])*256 +
						int(attribute.attrContent[1])
					entryInCp := klass.cpIndex[indexIntoCP]
					if entryInCp.entryType != IntConst {
						return pos, cfe("error: wrong type of constant value for char " +
							klass.utf8Refs[f.name].content)
					}
					f.constValue = klass.intConsts[entryInCp.slot]
				case types.Double: // double
					indexIntoCP := int(attribute.attrContent[0])*256 +
						int(attribute.attrContent[1])
					entryInCp := klass.cpIndex[indexIntoCP]
					if entryInCp.entryType != DoubleConst {
						return pos, cfe("error: wrong type of constant value for double " +
							klass.utf8Refs[f.name].content)
					}
					f.constValue = klass.doubles[entryInCp.slot]
				case types.Float: // float
					indexIntoCP := int(attribute.attrContent[0])*256 +
						int(attribute.attrContent[1])
					entryInCp := klass.cpIndex[indexIntoCP]
					if entryInCp.entryType != FloatConst {
						return pos, cfe("error: wrong type of constant value for float " +
							klass.utf8Refs[f.name].content)
					}
					f.constValue = klass.floats[entryInCp.slot]
				case types.Int: // integer
					indexIntoCP := int(attribute.attrContent[0])*256 +
						int(attribute.attrContent[1])
					entryInCp := klass.cpIndex[indexIntoCP]
					if entryInCp.entryType != IntConst {
						return pos, cfe("error: wrong type of constant value for integer " +
							klass.utf8Refs[f.name].content)
					}
					f.constValue = klass.intConsts[entryInCp.slot]
				case types.Long: // long
					indexIntoCP := int(attribute.attrContent[0])*256 +
						int(attribute.attrContent[1])
					entryInCp := klass.cpIndex[indexIntoCP]
					if entryInCp.entryType != LongConst {
						return pos, cfe("error: wrong type of constant value for long " +
							klass.utf8Refs[f.name].content)
					}
					f.constValue = klass.longConsts[entryInCp.slot]
				case types.Short: // short--same logic as int, only message is different
					indexIntoCP := int(attribute.attrContent[0])*256 +
						int(attribute.attrContent[1])
					entryInCp := klass.cpIndex[indexIntoCP]
					if entryInCp.entryType != IntConst {
						return pos, cfe("error: wrong type of constant value for short " +
							klass.utf8Refs[f.name].content)
					}
					f.constValue = klass.intConsts[entryInCp.slot]
				}
			} else { // append the attribute only if it's not ConstantValue
				f.attributes = append(f.attributes, attribute)
			}
			pos = k
		}

		klass.fields = append(klass.fields, f)

		if log.Level == log.FINEST {
			_, _ = fmt.Fprintf(os.Stderr, "\tField %s, desc: %s has %d attributes, access flags: %X.",
				klass.utf8Refs[f.name].content, klass.utf8Refs[f.description].content,
				len(f.attributes), accessFlags)
			if log.Level == log.FINEST && f.isStatic == true {
				_, _ = fmt.Fprintln(os.Stderr, " Field is static")
			}
			if len(f.attributes) > 0 {
				_, _ = fmt.Fprintf(os.Stderr, "First attrib: %s\n",
					klass.utf8Refs[f.attributes[0].attrName].content)
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "\n")
			}
		}
	}
	return pos, nil
}

// Get the number of methods in this class
func parseMethodCount(bytes []byte, loc int, klass *ParsedClass) (int, error) {
	pos := loc
	methodCount, err := intFrom2Bytes(bytes, pos+1)
	pos += 2
	if err != nil {
		return pos, cfe("Invalid fetch of method count")
	}

	_ = log.Log("method count: "+strconv.Itoa(methodCount), log.FINEST)
	klass.methodCount = methodCount
	return pos, nil
}

// get the count of the class attributes (which form the last group of elements in
// the class file).
func parseClassAttributeCount(bytes []byte, loc int, klass *ParsedClass) (int, error) {
	pos := loc
	attributeCount, err := intFrom2Bytes(bytes, pos+1)
	pos += 2
	if err != nil {
		return pos, cfe("Invalid fetch of class attribute count")
	}

	_ = log.Log("Class attribute count: "+strconv.Itoa(attributeCount), log.FINEST)
	klass.attribCount = attributeCount
	return pos, nil
}

func parseClassAttributes(bytes []byte, loc int, klass *ParsedClass) (int, error) {
	pos := loc
	for j := 0; j < klass.attribCount; j++ {
		attrib, location, err := fetchAttribute(klass, bytes, pos)
		pos = location
		if err == nil {
			klass.attributes = append(klass.attributes, attrib)
		} else {
			return pos, cfe("Error fetching class attribute in class: " +
				klass.className)
		}

		_ = log.Log("Class: "+klass.className+", attribute: "+klass.utf8Refs[attrib.attrName].content,
			log.FINEST)

		switch klass.utf8Refs[attrib.attrName].content {
		case "BootstrapMethods":
			// see: https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.7.23
			loc = 0
			bootstrapCount, err1 := u16From2bytes(attrib.attrContent, loc)
			loc += 2
			if err1 != nil {
				break // error msg will already have been shown
			} else {
				klass.bootstrapCount = int(bootstrapCount)
			}
			for m := 0; m < klass.bootstrapCount; m++ {
				bsm := bootstrapMethod{}
				methodRef, err2 := u16From2bytes(attrib.attrContent, loc)
				loc += 2
				if err2 != nil || klass.cpIndex[methodRef].entryType != MethodHandle {
					return pos, cfe("Invalid method reference in Boostrap method #" + strconv.Itoa(m))
				} else {
					bsm.methodRef = int(methodRef)
				}

				bootstrapArgCount, _ := intFrom2Bytes(attrib.attrContent, loc)
				loc += 2
				if bootstrapArgCount > 0 {
					for n := 0; n < bootstrapArgCount; n++ {
						arg, _ := intFrom2Bytes(attrib.attrContent, loc)
						loc += 2
						bsm.args = append(bsm.args, arg)
					}
				}
				klass.bootstraps = append(klass.bootstraps, bsm)
			}
			_ = log.Log("    "+strconv.Itoa(klass.bootstrapCount)+" bootstrap method(s)", log.FINEST)

		case "Deprecated":
			klass.deprecated = true

		case "SourceFile":
			sourceNameIndex, _ := intFrom2Bytes(attrib.attrContent, 0)
			utf8slot := klass.cpIndex[sourceNameIndex].slot
			sourceFile := klass.utf8Refs[utf8slot].content // points to the name of the source file
			klass.sourceFile = sourceFile
			_ = log.Log("Source file: "+sourceFile, log.FINEST)
		}
	}
	return pos, nil
}
