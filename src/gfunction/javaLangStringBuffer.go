/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

// Implement the minimum number of gfunctions to be able to run the java/lang/StringBuffer class,
// which is the younger brother of java/lang/Stringbuilder. Both classes enable you to create
// a String from an array of characters.
// see: https://docs.oracle.com/en/java/javase/17/docs/api/java.base/java/lang/StringBuffer.html

func Load_Lang_StringBuffer() {

	// === OBJECT INSTANTIATION ===

	MethodSignatures["java/lang/StringBuffer.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/lang/StringBuffer.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
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
}
