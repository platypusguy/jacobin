/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import "jacobin/object"

// Implement the minimum number of gfunctions to be able to run the java/lang/StringBuffer class,
// which is the younger brother of java/lang/Stringbuilder. Both classes enable you to create
// a String from an array of characters, but only StringBuffer is thread-safe.
// see: https://docs.oracle.com/en/java/javase/17/docs/api/java.base/java/lang/StringBuffer.html

func Load_Lang_StringBuffer() {

	// === OBJECT INSTANTIATION ===

	MethodSignatures["java/lang/StringBuffer.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/StringBuffer.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/StringBuffer.<init>(I)V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.<init>(Ljava/lang/CharSequence;)V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.<init>(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/StringBuffer.append(Ljava/lang/String;)Ljava/lang/StringBuffer;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  appendStringToStringBuffer,
		}
}

// append the string in the first parameter to the chars in the StringBuffer that's
// passed in the objectRef parameter (the second param)
func appendStringToStringBuffer(params []any) any {
	stringBufferObject := params[0].(*object.Object)

	strObjectToAppend := params[1].(*object.Object)
	strToAppend := strObjectToAppend.FieldTable["value"].Fvalue.([]byte)

	switch stringBufferObject.FieldTable["value"].Fvalue.(type) {
	case []byte:
		stringBufferContent := stringBufferObject.FieldTable["value"].Fvalue.([]byte)
		stringBufferContent = append(stringBufferContent, strToAppend...)
	case nil: // a raw StringBuffer
		stringBufferObject.FieldTable["value"] = object.Field{
			Ftype:  "[B",
			Fvalue: strToAppend,
		}
	}
	return stringBufferObject
}
