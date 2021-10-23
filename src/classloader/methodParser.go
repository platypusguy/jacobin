/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"jacobin/log"
	"strconv"
)

// Get the methods for this class. This can involve complex logic, but here
// we're just grabbing the info about the class and the actual method bytecodes
// as raw bytes. The description of the method entries in the spec is at:
// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.6
// The layout of the entries is:
// method_info {
//    u2             access_flags;
//    u2             name_index;
//    u2             descriptor_index;
//    u2             attributes_count;
//    attribute_info attributes[attributes_count];
// }
func parseMethods(bytes []byte, loc int, klass *parsedClass) (int, error) {
	pos := loc
	var meth method
	for i := 0; i < klass.methodCount; i++ {
		meth = method{}
		accessFlags, err := intFrom2Bytes(bytes, pos+1)
		pos += 2
		if err != nil {
			return pos, cfe("Invalid fetch of method access flags in class: " +
				klass.className)
		}

		nameIndex, err := intFrom2Bytes(bytes, pos+1)
		pos += 2
		if err != nil {
			return pos, cfe("Invalid fetch of method name index in class: " +
				klass.className)
		}
		nameSlot, err2 := fetchUTF8slot(klass, nameIndex)

		descIndex, err3 := intFrom2Bytes(bytes, pos+1)
		pos += 2
		if err2 != nil || err3 != nil {
			return pos, cfe("Invalid fetch of method description index in method: " +
				klass.utf8Refs[nameSlot].content)
		}
		descSlot, err4 := fetchUTF8slot(klass, descIndex)
		if err4 != nil {
			return pos, cfe("Invalid fetch of method description slot in method: " +
				klass.utf8Refs[nameSlot].content)
		}

		attrCount, err := intFrom2Bytes(bytes, pos+1)
		pos += 2
		if err != nil {
			return pos, cfe("Invalid fetch of method attribute count in method: " +
				klass.utf8Refs[nameSlot].content)
		}

		meth.accessFlags = accessFlags
		meth.name = nameSlot
		meth.description = descSlot

		// The Code attribute has sub-attributes that are important to right execution
		// The following code goes through those sub-attributes and processes them.

		if attrCount > 1 {
			log.Log(
				"Method: "+klass.utf8Refs[nameSlot].content+" Desc: "+
					klass.utf8Refs[descSlot].content+" has "+strconv.Itoa(attrCount)+" attributes",
				log.FINEST)
		}

		for j := 0; j < attrCount; j++ {
			attrib, location, err5 := fetchAttribute(klass, bytes, pos)
			pos = location
			if err5 == nil {
				meth.attributes = append(meth.attributes, attrib)
				// switch on the name of the attribute (listed here in alpha order)
				switch klass.utf8Refs[attrib.attrName].content {
				case "Code":
					if attrCount > 1 {
						log.Log("    Attribute: Code", log.FINEST)
					} else {
						log.Log("Method: "+klass.utf8Refs[nameSlot].content+" Desc: "+
							klass.utf8Refs[descSlot].content+" has "+strconv.Itoa(attrCount)+
							" attribute: Code", log.FINEST)
					}
					if parseCodeAttribute(attrib, &meth, klass) != nil {
						return pos, cfe("") // error msg will already have been shown to user
					}
				case "Deprecated":
					meth.deprecated = true
					log.Log("    Attribute: Deprecated", log.FINEST)
				case "Exceptions":
					log.Log("    Attribute: Exceptions", log.FINEST)
					if parseExceptionsMethodAttribute(attrib, &meth, klass) != nil {
						return pos, cfe("") // error msg will already have been shown to user
					}
				case "MethodParameters":
					log.Log("    Attribute: MethodParameters", log.FINEST)
					if parseMethodParametersAttribute(attrib, &meth, klass) != nil {
						return pos, cfe("") // error msg will already have been shown to user
					}
				default:
					log.Log("    Attribute: "+klass.utf8Refs[attrib.attrName].content, log.FINEST)
				}

			} else {
				return pos, cfe("Error fetching method attribute in method: " +
					klass.utf8Refs[nameSlot].content)
			}
		}
		klass.methods = append(klass.methods, meth)
	}

	return pos, nil
}

// parse the Code attribute and its sub-attributes. Details of the contents here:
// https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.7.3
func parseCodeAttribute(att attr, meth *method, klass *parsedClass) error {
	methodName := klass.utf8Refs[meth.name].content
	ca := codeAttrib{}

	pos := -1
	maxStack, err := intFrom2Bytes(att.attrContent, pos+1)
	pos += 2
	if err != nil {
		return cfe("Error getting maxStack value in Code attribute in " + klass.className)
	}

	maxLocals, err := intFrom2Bytes(att.attrContent, pos+1)
	pos += 2
	if err != nil {
		return cfe("Error getting maxLocals value in Code attribute in " + klass.className)
	}

	codeLength, err := intFrom4Bytes(att.attrContent, pos+1)
	pos += 4
	if err != nil {
		return cfe("Error getting code length in Code attribute in " + klass.className)
	}

	var code []byte
	for i := 0; i < codeLength; i++ {
		code = append(code, att.attrContent[pos+1+i])
	}
	pos += codeLength

	exceptionCount, err := intFrom2Bytes(att.attrContent, pos+1)
	pos += 2
	if err != nil {
		return cfe("Error getting count of exceptions in Code attribute in " + klass.className)
	}

	if exceptionCount > 0 {
		log.Log("        Method: "+methodName+" throws "+strconv.Itoa(exceptionCount)+" exception(s)",
			log.FINEST)
		for k := 0; k < exceptionCount; k++ {
			ex := exception{}
			ex.startPc, err = intFrom2Bytes(att.attrContent, pos+1)
			ex.endPc, err = intFrom2Bytes(att.attrContent, pos+3)
			ex.handlerPc, err = intFrom2Bytes(att.attrContent, pos+5)
			ex.catchType, err = intFrom2Bytes(att.attrContent, pos+7)
			pos += 8

			if err != nil {
				return cfe("Error getting catch type for exception in " + methodName +
					"() of " + klass.className + "\n at position: " + strconv.Itoa(pos) +
					" in the method (after parse of start/endPC, handlerPc, and catch type)")
			}

			if ex.catchType != 0 {
				catchType := klass.cpIndex[ex.catchType]
				if catchType.entryType != ClassRef {
					return cfe("Invalid catchType in method " + methodName +
						" in " + klass.className)
				} else {
					log.Log("        Method: "+methodName+
						" throws exception: "+klass.utf8Refs[catchType.slot].content,
						log.FINEST)
				}
			}
			ca.exceptions = append(ca.exceptions, ex)
		}
	}

	ca.attributes = []attr{}
	attrCount, err := intFrom2Bytes(att.attrContent, pos+1)
	pos += 2
	if err != nil {
		return cfe("Error getting attributes in Code attribute of " + methodName +
			"() of " + klass.className)
	}

	if attrCount > 0 {
		log.Log("        Code attribute has "+strconv.Itoa(attrCount)+
			" attributes: ", log.FINEST)
		for m := 0; m < attrCount; m++ {
			cat, loc, err2 := fetchAttribute(klass, att.attrContent, pos)
			if err2 != nil {
				return cfe("Error retrieving attributes in Code attribute of " + methodName +
					"() of " + klass.className)
			}
			pos = loc
			log.Log("        "+klass.utf8Refs[cat.attrName].content, log.FINEST)
			ca.attributes = append(ca.attributes, cat)
		}
	}

	ca.maxStack = maxStack
	ca.maxLocals = maxLocals
	ca.code = code
	meth.codeAttr = ca

	return nil
}

// The Exceptions attribute of a method indicates which checked exceptions a method
// can throw. See: https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.7.5
// The structure of the Exceptions attribute of a method is: {
// 		u2 attribute_name_index;
// 		u4 attribute_length;
// 		u2 number_of_exceptions;
// 		u2 exception_index_table[number_of_exceptions];
//   }
// The last two entries are in attrContent, which is a []byte. The last entry, per the spec,
// is a ClassRef entry, which consists of a CP index that points to UTF8 entry containing the
// name of the checked exception class, e.g., java/io/IOException
func parseExceptionsMethodAttribute(attrib attr, meth *method, klass *parsedClass) error {
	loc := -1
	exceptionCount, err := intFrom2Bytes(attrib.attrContent, loc+1)
	loc += 2
	if err != nil {
		return cfe("Error retrieving exception count in method " +
			klass.utf8Refs[meth.name].content)
	}

	for ex := 0; ex < exceptionCount; ex++ {
		// exception is an index into CP that points to a classRef
		cRefIndex, _ := intFrom2Bytes(attrib.attrContent, loc+1)
		loc += 2
		if klass.cpIndex[cRefIndex].entryType != ClassRef {
			return cfe("Exception attribute #" + strconv.Itoa(ex+1) +
				" in method " + klass.utf8Refs[meth.name].content +
				" does not point to a ClassRef CP entry")
		}

		// whichClassRef is the entry # in the classRefs array
		whichClassRef := klass.cpIndex[cRefIndex].slot
		// get the classRef from the slice of classRefs in the parsedClass
		classRef := klass.classRefs[whichClassRef]

		// the classRef should point to a UTF8 record with the name of the exception class
		exceptionName, err2 := fetchUTF8string(klass, classRef)
		if err2 != nil {
			return cfe("Exception attribute #" + strconv.Itoa(ex+1) +
				" in method " + klass.utf8Refs[meth.name].content +
				" has a ClassRef CP entry that does not point to a UTF8 string")
		}

		// if the previous fetch of the UTF8 record succeeded, this one shouldn't fail
		// so we don't check the error return
		whichUtf8Rec, _ := fetchUTF8slot(klass, classRef)

		// store the slot # of the utf8 entries into the method exceptions slice
		meth.exceptions = append(meth.exceptions, whichUtf8Rec)
		log.Log("        "+exceptionName, log.FINEST)
	}
	return nil
}

// Per the spec, 'A MethodParameters attribute records information about the formal parameters
// of a method, such as their names.' See: https://docs.oracle.com/javase/specs/jvms/se11/html/jvms-4.html#jvms-4.7.24
//    u2 attribute_name_index;
//    u4 attribute_length;
//    u1 parameters_count;
//    {   u2 name_index;
//        u2 access_flags;
//    } parameters[parameters_count];
// }
func parseMethodParametersAttribute(att attr, meth *method, klass *parsedClass) error {
	var err error
	pos := 0
	parametersCount := int(att.attrContent[pos])
	pos += 1
	if err != nil {
		return cfe("Error getting number of Parameter attributes in method: " +
			klass.utf8Refs[meth.name].content)
	}

	for k := 0; k < parametersCount; k++ {
		mpAttrib := paramAttrib{}
		paramNameIndex, err := intFrom2Bytes(att.attrContent, pos)
		pos += 2
		if err != nil {
			return cfe("Error getting name index for MethodParameters attribute #" +
				strconv.Itoa(k+1) + " in " + klass.utf8Refs[meth.name].content)
		}
		if paramNameIndex == 0 {
			mpAttrib.name = ""
		} else {
			mpAttrib.name, err = fetchUTF8string(klass, paramNameIndex)
		}
		if err != nil {
			return cfe("Error getting name of MethodParameters attribute #" +
				strconv.Itoa(k+1) + " in " + klass.utf8Refs[meth.name].content)
		}

		logName := "{none}"
		if mpAttrib.name != "" {
			logName = mpAttrib.name
		}
		log.Log("        "+logName, log.FINEST)

		accessFlags, err := intFrom2Bytes(att.attrContent, pos)
		if err != nil {
			return cfe("Error getting access flags of MethodParameters attribute #" +
				strconv.Itoa(k+1) + " in " + klass.utf8Refs[meth.name].content)
		}
		// do format check on the access flags here
		if accessFlags != 0x10 && accessFlags != 0x1000 && accessFlags != 0x8000 {
			return cfe("Invalid access flags of MethodParameters attribute #" +
				strconv.Itoa(k+1) + " in " + klass.utf8Refs[meth.name].content)
		}

		mpAttrib.accessFlags = accessFlags
		meth.parameters = append(meth.parameters, mpAttrib)
	}
	return nil
}
