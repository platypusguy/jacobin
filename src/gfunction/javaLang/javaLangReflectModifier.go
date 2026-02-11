/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-5 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"strings"
)

// Java Modifier constants as defined in java.lang.reflect.Modifier
const (
	PUBLIC       = 0x0001
	PRIVATE      = 0x0002
	PROTECTED    = 0x0004
	STATIC       = 0x0008
	FINAL        = 0x0010
	SYNCHRONIZED = 0x0020
	VOLATILE     = 0x0040
	TRANSIENT    = 0x0080
	NATIVE       = 0x0100
	INTERFACE    = 0x0200
	ABSTRACT     = 0x0400
	STRICT       = 0x0800
	SYNTHETIC    = 0x1000
	ANNOTATION   = 0x2000
	ENUM         = 0x4000
	MANDATED     = 0x8000
)

func Load_Lang_Reflect_Modifier() {
	ghelpers.MethodSignatures["java/lang/reflect/Modifier.isAbstract(I)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: modifierIsAbstract}
	ghelpers.MethodSignatures["java/lang/reflect/Modifier.isFinal(I)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: modifierIsFinal}
	ghelpers.MethodSignatures["java/lang/reflect/Modifier.isInterface(I)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: modifierIsInterface}
	ghelpers.MethodSignatures["java/lang/reflect/Modifier.isNative(I)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: modifierIsNative}
	ghelpers.MethodSignatures["java/lang/reflect/Modifier.isPrivate(I)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: modifierIsPrivate}
	ghelpers.MethodSignatures["java/lang/reflect/Modifier.isProtected(I)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: modifierIsProtected}
	ghelpers.MethodSignatures["java/lang/reflect/Modifier.isPublic(I)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: modifierIsPublic}
	ghelpers.MethodSignatures["java/lang/reflect/Modifier.isStatic(I)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: modifierIsStatic}
	ghelpers.MethodSignatures["java/lang/reflect/Modifier.isStrict(I)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: modifierIsStrict}
	ghelpers.MethodSignatures["java/lang/reflect/Modifier.isSynchronized(I)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: modifierIsSynchronized}
	ghelpers.MethodSignatures["java/lang/reflect/Modifier.isTransient(I)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: modifierIsTransient}
	ghelpers.MethodSignatures["java/lang/reflect/Modifier.isVolatile(I)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: modifierIsVolatile}
	ghelpers.MethodSignatures["java/lang/reflect/Modifier.toString(I)Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: modifierToString}
}

// java/lang/reflect/Modifier.isPublic(I)Z
func modifierIsPublic(params []interface{}) interface{} {
	mod := params[0].(int64)
	if (mod & PUBLIC) != 0 {
		return int64(1) // true
	}
	return int64(0) // false
}

// java/lang/reflect/Modifier.isPrivate(I)Z
func modifierIsPrivate(params []interface{}) interface{} {
	mod := params[0].(int64)
	if (mod & PRIVATE) != 0 {
		return int64(1) // true
	}
	return int64(0) // false
}

// java/lang/reflect/Modifier.isProtected(I)Z
func modifierIsProtected(params []interface{}) interface{} {
	mod := params[0].(int64)
	if (mod & PROTECTED) != 0 {
		return int64(1) // true
	}
	return int64(0) // false
}

// java/lang/reflect/Modifier.isStatic(I)Z
func modifierIsStatic(params []interface{}) interface{} {
	mod := params[0].(int64)
	if (mod & STATIC) != 0 {
		return int64(1) // true
	}
	return int64(0) // false
}

// java/lang/reflect/Modifier.isFinal(I)Z
func modifierIsFinal(params []interface{}) interface{} {
	mod := params[0].(int64)
	if (mod & FINAL) != 0 {
		return int64(1) // true
	}
	return int64(0) // false
}

// java/lang/reflect/Modifier.isSynchronized(I)Z
func modifierIsSynchronized(params []interface{}) interface{} {
	mod := params[0].(int64)
	if (mod & SYNCHRONIZED) != 0 {
		return int64(1) // true
	}
	return int64(0) // false
}

// java/lang/reflect/Modifier.isVolatile(I)Z
func modifierIsVolatile(params []interface{}) interface{} {
	mod := params[0].(int64)
	if (mod & VOLATILE) != 0 {
		return int64(1) // true
	}
	return int64(0) // false
}

// java/lang/reflect/Modifier.isTransient(I)Z
func modifierIsTransient(params []interface{}) interface{} {
	mod := params[0].(int64)
	if (mod & TRANSIENT) != 0 {
		return int64(1) // true
	}
	return int64(0) // false
}

// java/lang/reflect/Modifier.isNative(I)Z
func modifierIsNative(params []interface{}) interface{} {
	mod := params[0].(int64)
	if (mod & NATIVE) != 0 {
		return int64(1) // true
	}
	return int64(0) // false
}

// java/lang/reflect/Modifier.isInterface(I)Z
func modifierIsInterface(params []interface{}) interface{} {
	mod := params[0].(int64)
	if (mod & INTERFACE) != 0 {
		return int64(1) // true
	}
	return int64(0) // false
}

// java/lang/reflect/Modifier.isAbstract(I)Z
func modifierIsAbstract(params []interface{}) interface{} {
	mod := params[0].(int64)
	if (mod & ABSTRACT) != 0 {
		return int64(1) // true
	}
	return int64(0) // false
}

// java/lang/reflect/Modifier.isStrict(I)Z
func modifierIsStrict(params []interface{}) interface{} {
	mod := params[0].(int64)
	if (mod & STRICT) != 0 {
		return int64(1) // true
	}
	return int64(0) // false
}

// java/lang/reflect/Modifier.isVolatile(I)Z
// Note: This is used for methods to check if they are bridge methods
// In the modifier bit representation, VOLATILE (0x0040) is also used as BRIDGE for methods

// java/lang/reflect/Modifier.toString(I)Ljava/lang/String;
// Returns a string describing the access modifier flags in the specified modifier
func modifierToString(params []interface{}) interface{} {
	mod := params[0].(int64)

	var modifiers []string

	// Order matters - must match Java's order
	if (mod & PUBLIC) != 0 {
		modifiers = append(modifiers, "public")
	}
	if (mod & PROTECTED) != 0 {
		modifiers = append(modifiers, "protected")
	}
	if (mod & PRIVATE) != 0 {
		modifiers = append(modifiers, "private")
	}
	if (mod & ABSTRACT) != 0 {
		modifiers = append(modifiers, "abstract")
	}
	if (mod & STATIC) != 0 {
		modifiers = append(modifiers, "static")
	}
	if (mod & FINAL) != 0 {
		modifiers = append(modifiers, "final")
	}
	if (mod & TRANSIENT) != 0 {
		modifiers = append(modifiers, "transient")
	}
	if (mod & VOLATILE) != 0 {
		modifiers = append(modifiers, "volatile")
	}
	if (mod & SYNCHRONIZED) != 0 {
		modifiers = append(modifiers, "synchronized")
	}
	if (mod & NATIVE) != 0 {
		modifiers = append(modifiers, "native")
	}
	if (mod & STRICT) != 0 {
		modifiers = append(modifiers, "strictfp")
	}
	if (mod & INTERFACE) != 0 {
		modifiers = append(modifiers, "interface")
	}

	result := strings.Join(modifiers, " ")
	return object.StringObjectFromGoString(result)
}
