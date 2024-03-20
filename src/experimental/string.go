/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package experimental

import "jacobin/object"

// This experimental package is designed to codify the various kinds of strings and provide
// convenience methods to create and manipulate them.
// There exist only the following kinds of strings:
// 1) golang string -- this is commonly used inside the JVM
// 2) byte array -- less frequently used. Can represent either a string of bytes or a string of chars
// 3) String object -- used only when passing args to/from Java methods and gfunctions. Important note:
//    string objects are the *only* form of strings passed to/from Java methods and gfunctions.
//
// Implementation details:
// * the string pool stores only golang strings. This is done for performance reasons.
// * string objects' "value" field contains a byte array, which is required by Java methods

import (
	"jacobin/stringPool"
	"jacobin/types"
)

// NewString creates an empty string.
func NewStringObject() *object.Object {
	s := new(object.Object)
	s.Mark.Hash = 0
	s.KlassName = object.StringPoolStringIndex // =  java/lang/String
	s.FieldTable = make(map[string]object.Field)

	// ==== now the fields ====

	// value: the content of the string as array of runes or bytes
	// Note: Post JDK9, this field is an array of bytes, so as to
	// enable compact strings.

	// value := make([]byte, 0) // presently empty // commented out due to JACOBIN-463
	// valueField := Field{Ftype: types.ByteArray, Fvalue: value}
	valueField := object.Field{Ftype: types.ByteArray, Fvalue: ""} // empty string
	s.FieldTable["value"] = valueField

	// coder has two possible values:
	// LATIN(e.g., 0 = bytes for compact strings) or UTF16(e.g., 1 = UTF16)
	coderField := object.Field{Ftype: types.Byte, Fvalue: byte(0)}
	s.FieldTable["coder"] = coderField

	// the hash code, which is initialized to 0
	hash := object.Field{Ftype: types.Int, Fvalue: int32(0)}
	s.FieldTable["hash"] = hash

	// hashIsZero: only true in rare case where compute hash is 0
	hashIsZero := object.Field{Ftype: types.Byte, Fvalue: byte(0)}
	s.FieldTable["hashIsZero"] = hashIsZero

	// The following static fields are preloaded in statics/LoadStaticsString()
	//   COMPACT_STRINGS (always true for JDK >= 9)
	//   UTF_8.INSTANCE ptr to encoder
	//   ISO_8859_1.INSTANCE ptr to encoder
	//   US_ASCII.INSTANCE ptr to encoder
	//   CodingErrorAction.REPLACE
	//   CASE_INSENSITIVE_ORDER
	//   serialPersistentFields

	return s
}

// StringObjectFromGoString: convenience method to create a string object from a Golang string
func StringObjectFromGoString(str string) *object.Object {
	newStr := NewStringObject()
	newStr.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: []byte(str)}
	return newStr
}

// GoStringFromStringObject: convenience method to extract a Go string from a String object (Java string)
func GoStringFromStringObject(obj *object.Object) string {
	if obj != nil && obj.KlassName == object.StringPoolStringIndex {
		if obj.FieldTable["value"].Fvalue != nil {
			return obj.FieldTable["value"].Fvalue.(string)
		}
	}
	return ""
}

// ByteArrayFromStringObject: convenience method to extract a byte array from a String object (Java string)
func ByteArrayFromStringObject(obj *object.Object) []byte {
	if obj != nil && obj.KlassName == object.StringPoolStringIndex {
		return obj.FieldTable["value"].Fvalue.([]byte)
	} else {
		return nil
	}
}

// StringObjectFromByteArray: convenience method to create a string object from a byte array
func StringObjectFromByteArray(bytes []byte) *object.Object {
	newStr := NewStringObject()
	newStr.FieldTable["value"] = object.Field{Ftype: types.ByteArray, Fvalue: bytes}
	return newStr
}

// StringPoolIndexFromStringObject: convenience method to extract a string pool index from a String object
func StringPoolIndexFromStringObject(obj *object.Object) uint32 {
	if obj != nil && obj.KlassName == object.StringPoolStringIndex {
		str := string(obj.FieldTable["value"].Fvalue.([]byte))
		index := stringPool.GetStringIndex(&str)
		return index
	} else {
		return types.InvalidStringIndex
	}
}

// GoStringFromStringPoolIndex: convenience method to extract a Go string from a string pool index
func GoStringFromStringPoolIndex(index uint32) string {
	if index < stringPool.GetStringPoolSize() {
		return *stringPool.GetStringPointer(index)
	} else {
		return ""
	}
}

// StringObjectFromStringPoolIndex: convenience method to create a string object using a string pool index
func StringObjectFromPoolIndex(index uint32) *object.Object {
	if index < stringPool.GetStringPoolSize() {
		return StringObjectFromGoString(*stringPool.GetStringPointer(index))
	} else {
		return nil
	}
}

// ByteArrayFromStringPoolIndex: convenience method to get a byte array using a string pool index
func ByteArrayFromStringPoolIndex(index uint32) []byte {
	if index < stringPool.GetStringPoolSize() {
		return []byte(*stringPool.GetStringPointer(index))
	} else {
		return nil
	}
}

// determine whether an object is a string object (i.e., a Java string)
// assumes that any object whose Klass pointer points to java/lang/String
// is an instance of a Java string
func IsStringObject(unknown any) bool {
	if unknown == nil {
		return false
	}

	o, ok := unknown.(*object.Object)
	if !ok {
		return false
	}

	if o.KlassName == object.StringPoolStringIndex {
		return true
	}
	return false
}
