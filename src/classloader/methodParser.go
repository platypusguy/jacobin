/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-3 by Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"fmt"
	"jacobin/log"
	"jacobin/stringPool"
	"jacobin/util"
	"sort"
	"strconv"
)

// Get the methods for this class. This can involve complex logic, but here
// we're just grabbing the info about the class and the actual method bytecodes
// as raw bytes. The description of the method entries in the spec is at:
// https://docs.oracle.com/javase/specs/jvms/se17/html/jvms-4.html#jvms-4.6
// The layout of the entries is:
//
//	method_info {
//	   u2             access_flags;
//	   u2             name_index;
//	   u2             descriptor_index;
//	   u2             attributes_count;
//	   attribute_info attributes[attributes_count];
//	}
func parseMethods(bytes []byte, loc int, klass *ParsedClass) (int, error) {
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
			_ = log.Log(
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
						_ = log.Log("    Attribute: Code", log.FINEST)
					} else {
						_ = log.Log("Method: "+klass.utf8Refs[nameSlot].content+" Desc: "+
							klass.utf8Refs[descSlot].content+" has "+strconv.Itoa(attrCount)+
							" attribute: Code", log.FINEST)
					}
					if parseCodeAttribute(attrib, &meth, klass) != nil {
						return pos, cfe("") // error msg will already have been shown to user
					}
				case "Deprecated":
					meth.deprecated = true
					_ = log.Log("    Attribute: Deprecated", log.FINEST)
				case "Exceptions":
					_ = log.Log("    Attribute: Exceptions", log.FINEST)
					if parseExceptionsMethodAttribute(attrib, &meth, klass) != nil {
						return pos, cfe("") // error msg will already have been shown to user
					}
				case "MethodParameters":
					_ = log.Log("    Attribute: MethodParameters", log.FINEST)
					// JACOBIN-577: Removed because JDK 21 causes something in the parsing
					// of this attribute to become unhinged. As this attribute is not needed in
					// the execution of a class, we've temporarily chosen to block this call.
					// When time permits or if the MethodParameters attribute assumes a new importance,
					// we'll return here and figure out what needs to be done.
					//
					// if parseMethodParametersAttribute(attrib, &meth, klass) != nil {
					// 	return pos, cfe("") // error msg will already have been shown to user
					// }
				default:
					_ = log.Log("    Attribute: "+klass.utf8Refs[attrib.attrName].content, log.FINEST)
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
// https://docs.oracle.com/javase/specs/jvms/se17/html/jvms-4.html#jvms-4.7.3
func parseCodeAttribute(att attr, meth *method, klass *ParsedClass) error {
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
		_ = log.Log("        Method: "+methodName+" throws "+strconv.Itoa(exceptionCount)+" exception(s)",
			log.FINEST)
		for k := 0; k < exceptionCount; k++ {
			ex := exception{}
			ex.startPc, _ = intFrom2Bytes(att.attrContent, pos+1)
			ex.endPc, _ = intFrom2Bytes(att.attrContent, pos+3)
			ex.handlerPc, _ = intFrom2Bytes(att.attrContent, pos+5)
			ex.catchType, err = intFrom2Bytes(att.attrContent, pos+7)
			// fmt.Printf("DEBUG parseCodeAttribute methodName=%s, ex = %v\n", methodName, ex)
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
					_ = log.Log("        Method: "+methodName+
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
		_ = log.Log("        Code attribute has "+strconv.Itoa(attrCount)+
			" attributes: ", log.FINEST)
		for m := 0; m < attrCount; m++ {
			subAttr, loc, err2 := fetchAttribute(klass, att.attrContent, pos)
			if err2 != nil {
				return cfe("Error retrieving attributes in Code attribute of " + methodName +
					"() of " + klass.className)
			}
			pos = loc
			_ = log.Log("        "+klass.utf8Refs[subAttr.attrName].content, log.FINEST)
			if klass.utf8Refs[subAttr.attrName].content == "LineNumberTable" &&
				!util.IsFilePartOfJDK(&klass.className) {
				buildLineNumberTable(&ca, &subAttr, methodName)
			}
			ca.attributes = append(ca.attributes, subAttr)
		}
	}

	ca.maxStack = maxStack
	ca.maxLocals = maxLocals
	ca.code = code
	meth.codeAttr = ca

	return nil
}

// build the table of line numbers (that map bytecode location to source line #)
// consult https://docs.oracle.com/javase/specs/jvms/se17/html/jvms-4.html#jvms-4.7.12
func buildLineNumberTable(codeAttr *codeAttrib, thisAttr *attr, methodName string) {
	entryCount := uint(thisAttr.attrContent[0])*256 + uint(thisAttr.attrContent[1])
	loc := 2 // we're two bytes into the attr.Content byte array
	if entryCount < 1 {
		(*codeAttr).sourceLineTable = nil
		return
	}

	var table []BytecodeToSourceLine
	if (*codeAttr).sourceLineTable != nil { // we could be adding to the table
		table = []BytecodeToSourceLine{}
		(*codeAttr).sourceLineTable = &table
	}
	var i uint
	for i = 0; i < entryCount; i++ {
		bytecodeNumber := uint16(thisAttr.attrContent[loc])*256 + uint16(thisAttr.attrContent[loc+1])
		sourceLineNumber := uint16(thisAttr.attrContent[loc+2])*256 + uint16(thisAttr.attrContent[loc+3])
		loc += 4

		tableEntry := BytecodeToSourceLine{bytecodeNumber, sourceLineNumber}
		table = append(table, tableEntry)
	}

	// now sort the table
	if len(table) > 1 {
		sort.Sort(b2sTable(table))
	}

	(*codeAttr).sourceLineTable = &table

	// if methodName == "main" {
	// 	fmt.Fprintf(os.Stderr, "%v\n", table)
	// }
}

// the following four lines are all needed for the call to Sort()
type b2sTable []BytecodeToSourceLine

func (t b2sTable) Len() int           { return len(t) }
func (t b2sTable) Swap(k, j int)      { (t)[k], (t)[j] = (t)[j], (t)[k] }
func (t b2sTable) Less(k, j int) bool { return (t)[k].BytecodePos < (t)[j].BytecodePos }

// BytecodeToSourceLine maps the PC in a method to the
// corresponding source line in the original source file.
// This data is captured in the method's attributes
type BytecodeToSourceLine struct {
	BytecodePos uint16
	SourceLine  uint16
}

// The Exceptions attribute of a method indicates which checked exceptions a method
// can throw. See: https://docs.oracle.com/javase/specs/jvms/se17/html/jvms-4.html#jvms-4.7.5
//
//	The structure of the Exceptions attribute of a method is: {
//			u2 attribute_name_index;
//			u4 attribute_length;
//			u2 number_of_exceptions;
//			u2 exception_index_table[number_of_exceptions];
//	  }
//
// The last two entries are in attrContent, which is a []byte. The last entry, per the spec,
// is a ClassRef entry, which consists of a CP index that points to UTF8 entry containing the
// name of the checked exception class, e.g., java/io/IOException
func parseExceptionsMethodAttribute(attrib attr, meth *method, klass *ParsedClass) error {
	loc := -1
	exceptionCount, err := intFrom2Bytes(attrib.attrContent, loc+1)
	loc += 2
	if err != nil {
		return cfe("Error retrieving exception count in method " +
			klass.utf8Refs[meth.name].content)
	}

	for ex := 0; ex < exceptionCount; ex++ {
		// exception is an index into CP that points to a exceptionClassRef
		cRefIndex, _ := intFrom2Bytes(attrib.attrContent, loc+1)
		loc += 2
		if klass.cpIndex[cRefIndex].entryType != ClassRef {
			return cfe("Exception attribute #" + strconv.Itoa(ex+1) +
				" in method " + klass.utf8Refs[meth.name].content +
				" does not point to a ClassRef CP entry")
		}

		// whichClassRef is the entry # in the classRefs array
		whichClassRef := klass.cpIndex[cRefIndex].slot
		// get the exceptionClassRef from the slice of classRefs in the ParsedClass
		exceptionClassRef := klass.classRefs[whichClassRef]

		// the exceptionClassRef should point to a UTF8 record with the name of the exception class
		// exceptionName, err2 := FetchUTF8string(klass, exceptionClassRef)
		exceptionName := stringPool.GetStringPointer(exceptionClassRef)
		if exceptionName == nil {
			return cfe("Exception attribute #" + strconv.Itoa(ex+1) +
				" in method " + klass.utf8Refs[meth.name].content +
				"  does not point to a valid stringPool entry")
			// return cfe("Exception attribute #" + strconv.Itoa(ex+1) +
			// 	" in method " + klass.utf8Refs[meth.name].content +
			// 	" has a ClassRef CP entry that does not point to a UTF8 string")
		}
		//
		// // if the previous fetch of the UTF8 record succeeded, this one shouldn't fail
		// // so we don't check the error return
		// whichUtf8Rec, _ := fetchUTF8slot(klass, exceptionClassRef)

		// store the slot # of the utf8 entries into the method exceptions slice
		meth.exceptions = append(meth.exceptions, exceptionClassRef)
		_ = log.Log("        "+*exceptionName, log.FINEST)
	}
	return nil
}

// Per the spec, 'A MethodParameters attribute records information about the formal parameters
// of a method, such as their names.' See: https://docs.oracle.com/javase/specs/jvms/se17/html/jvms-4.html#jvms-4.7.24
//
//	   u2 attribute_name_index;
//	   u4 attribute_length;
//	   u1 parameters_count;
//	   {   u2 name_index;
//	       u2 access_flags;
//	   } parameters[parameters_count];
//	}
func parseMethodParametersAttribute(att attr, meth *method, klass *ParsedClass) error {
	pos := 0
	parametersCount := int(att.attrContent[pos])
	pos += 1

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
			mpAttrib.name, err = FetchUTF8string(klass, paramNameIndex)
		}
		if err != nil {
			return cfe("Error getting name of MethodParameters attribute #" +
				strconv.Itoa(k+1) + " in " + klass.utf8Refs[meth.name].content)
		}

		logName := "{none}"
		if mpAttrib.name != "" {
			logName = mpAttrib.name
		}
		_ = log.Log("        "+logName, log.FINEST)

		accessFlags, err := intFrom2Bytes(att.attrContent, pos)
		if err != nil {
			return cfe("Error getting access flags of MethodParameters attribute #" +
				strconv.Itoa(k+1) + " in " + klass.utf8Refs[meth.name].content)
		}
		// do format check on the access flags here
		switch accessFlags {
		case 0x00, 0x10, 0x1000, 0x1010, 0x8000, 0x8010:
			break
		default:
			errMsg := fmt.Sprintf(
				"Invalid access flags of MethodParameters attribute #%s in method %s: %X",
				strconv.Itoa(k+1), klass.utf8Refs[meth.name].content, accessFlags)
			return cfe(errMsg)
		}

		mpAttrib.accessFlags = accessFlags
		meth.parameters = append(meth.parameters, mpAttrib)
	}
	return nil
}
