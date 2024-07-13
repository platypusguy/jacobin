/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import "jacobin/object"

// Implement the minimum number of gfunctions to be able to run the java/lang/StringBuffer class,
// which is the younger brother of java/lang/Stringbuilder. Both classes enable you to create
// a String from an array of characters, but only StringBuffer is thread-safe.
// see: https://docs.oracle.com/en/java/javase/17/docs/api/java.base/java/lang/StringBuffer.html

func Load_Lang_StringBuffer() {

	// === Instantiation ===

	MethodSignatures["java/lang/StringBuffer.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	// === Methods ===

	MethodSignatures["java/lang/StringBuffer.append(Ljava/lang/String;)Ljava/lang/StringBuffer;"] = // append string
		GMeth{
			ParamSlots: 1,
			GFunction:  appendStringToStringBuffer,
		}

	MethodSignatures["java/lang/StringBuffer.append([C)Ljava/lang/StringBuffer"] = // append char array
		GMeth{
			ParamSlots: 1,
			GFunction:  appendStringToStringBuffer,
		}

	MethodSignatures["java/lang/StringBuffer.append(Ljava/lang/CharSequence;)Ljava/lang/StringBuffer"] = // append char seq
		GMeth{
			ParamSlots: 1,
			GFunction:  appendStringToStringBuffer,
		}
	MethodSignatures["java/lang/StringBuffer.append(Ljava/lang/Object;)Ljava/lang/StringBuffer;"] = // append object
		GMeth{
			ParamSlots: 1,
			GFunction:  appendStringToStringBuffer,
		}

	MethodSignatures["java/lang/StringBuffer.append(Ljava/lang/StringBuffer;)Ljava/lang/StringBuffer"] = // append stringBuffer
		GMeth{
			ParamSlots: 1,
			GFunction:  appendStringToStringBuffer,
		}
}

// append the string in the second parameter to the chars in the StringBuffer that's
// passed in the objectRef parameter (the first param)
func appendStringToStringBuffer(params []any) any {
	stringBufferObject := params[0].(*object.Object)
	stringBufferStringLen := stringBufferObject.FieldTable["count"].Fvalue.(int64)

	strObjectToAppend := params[1].(*object.Object)
	strToAppend := strObjectToAppend.FieldTable["value"].Fvalue.([]byte)

	switch stringBufferObject.FieldTable["value"].Fvalue.(type) {
	case []byte: // the usual case
		if stringBufferStringLen == 0 {
			stringBufferObject.FieldTable["value"] = object.Field{
				Ftype:  "[B",
				Fvalue: strToAppend,
			}
			stringBufferObject.FieldTable["count"] = object.Field{
				Ftype:  "I",
				Fvalue: int64(len(strToAppend)),
			}
		} else {
			stringBufferContent := stringBufferObject.FieldTable["value"].Fvalue.([]byte)
			stringBufferContent = append(stringBufferContent, strToAppend...)
			stringBufferObject.FieldTable["count"] = object.Field{
				Ftype:  "I",
				Fvalue: int64(len(stringBufferContent)),
			}
		}
	case nil: // a raw StringBuffer
		stringBufferObject.FieldTable["value"] = object.Field{
			Ftype:  "[B",
			Fvalue: strToAppend,
		}
		stringBufferObject.FieldTable["count"] = object.Field{
			Ftype:  "I",
			Fvalue: int64(len(strToAppend)),
		}
	}
	return stringBufferObject
}
