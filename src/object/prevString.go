/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package object

import (
	"fmt"
	"jacobin/exceptions"
	"jacobin/statics"
	"jacobin/stringPool"
	"jacobin/types"
	"unsafe"
)

// Strings are so commonly used in Java, that it makes sense
// to have a means of creating them quickly, rather than building
// them from scratch each time by walking through the constant pool.

// var StringClassName = "java/lang/String"
// var StringClassRef = "Ljava/lang/String;"
// var StringPoolStringIndex = uint32(1)
// var EmptyString = ""

// NewString creates an empty string.
func NewString() *Object {
	s := new(Object)
	s.Mark.Hash = 0
	s.KlassName = StringPoolStringIndex // =  java/lang/String
	// s.Klass = &StringClassName // java/lang/String
	s.FieldTable = make(map[string]Field)

	// ==== now the fields ====

	// value: the content of the string as array of runes or bytes
	// Note: Post JDK9, this field is an array of bytes, so as to
	// enable compact strings.

	// value := make([]byte, 0) // presently empty // commented out due to JACOBIN-463
	// valueField := Field{Ftype: types.ByteArray, Fvalue: value}
	valueField := Field{Ftype: types.StringIndex,
		Fvalue: types.InvalidStringIndex} // empty (i.e., non existent string)
	s.FieldTable["value"] = valueField

	// coder has two possible values:
	// LATIN(e.g., 0 = bytes for compact strings) or UTF16(e.g., 1 = UTF16)
	coderField := Field{Ftype: types.Byte, Fvalue: int64(0)}
	s.FieldTable["coder"] = coderField

	// the hash code, which is initialized to 0
	hash := Field{Ftype: types.Int, Fvalue: int64(0)}
	s.FieldTable["hash"] = hash

	// hashIsZero: only true in rare case where compute hash is 0
	hashIsZero := Field{Ftype: types.Bool, Fvalue: types.JavaBoolFalse}
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

// NewStringFromGoString converts a go string to a Java string-like
// entity, in which the chars are stored as runes, rather than chars
// or as bytes, depending on the status of COMPACT_STRINGS (which defaults
// to true for JDK >= 9). In the latter case, the string is interned in
// the string pool and the field is set to the index in the pool.

func NewStringFromGoString(in string) *Object {
	s := NewString()
	if statics.GetStaticValue("java/lang/String", "COMPACT_STRINGS") == types.JavaBoolFalse {
		s.FieldTable["value"] = Field{types.RuneArray, in}
	} else {
		s.FieldTable["value"] = Field{types.StringIndex, stringPool.GetStringIndex(&in)}
		// s.FieldTable["value"] = Field{types.ByteArray, []byte(in)} // changed in JACOBIN-463
	}
	return s
}

/* This function no longer is used due mostly to JACOBIN-463
// CreateCompactStringFromGoString creates a string in which the chars
// are stored as bytes--that is, a compact string.
func CreateCompactStringFromGoString(in *string) *Object {
	s := NewString()
	s.FieldTable["value"] = Field{types.ByteArray, []byte(*in)}
	return s
}
*/

// CreateStringPoolEntryFromGoString creates an object that points to an interned string
func CreateStringPoolEntryFromGoString(in *string) *Object {
	s := NewString()
	s.FieldTable["value"] = Field{types.StringIndex, stringPool.GetStringIndex(in)}
	return s
}

// convenience method to extract a Go string from a Java string
func GetGoStringFromJavaStringPtr(strPtr *Object) string {
	s := *strPtr
	bytes := s.FieldTable["value"].Fvalue.([]byte)
	return string(bytes)
}

// determine whether an object is a Java string
// assumes that any object whose Klass pointer points to java/lang/String
// is an instance of a Java string
func IsJavaString(unknown any) bool {
	var objPtr *Object

	if unknown == nil {
		return false
	}

	switch unknown.(type) {
	case *Object:
		objPtr = unknown.(*Object)
		break
	default:
		return false
	}

	return *stringPool.GetStringPointer(objPtr.KlassName) == StringClassName
}

/*
------------------------------------------
The string pool mid-level functions follow
------------------------------------------
*/

// MakeEmptyStringObject creates an empty object.Object.
// It is expected that the caller will fill in the FieldTable.
func MakeEmptyStringObject() *Object {
	object := Object{}
	ptrObject := uintptr(unsafe.Pointer(&object))
	object.Mark.Hash = uint32(ptrObject)
	object.KlassName = StringPoolStringIndex // = java/lang/String

	// initialize the map of this object's fields
	object.FieldTable = make(map[string]Field)
	return &object
}

func NewPoolStringFromGoString(str string) *Object {
	objPtr := MakeEmptyStringObject()
	/* TODO - Is ignoring the COMPACT_STRINGS flag valid?
	if statics.GetStaticValue("java/lang/String", "COMPACT_STRINGS") == types.JavaBoolFalse {
		objPtr.FieldTable["value"] = Field{types.RuneArray, in}
	} else {
		objPtr.FieldTable["value"] = Field{types.StringIndex, GetStringIndex(&in)}
	}
	*/
	objPtr.FieldTable["value"] = Field{types.StringIndex, stringPool.GetStringIndex(&str)}
	return objPtr
}

// GetGoStringFromObject : convenience method to extract a Go string from a Pool string
func GetGoStringFromObject(strPtr *Object) string {
	obj := *strPtr
	fld := obj.FieldTable["value"]
	if fld.Ftype != types.StringIndex {
		errMsg := fmt.Sprintf("GetGoStringFromObject: Expected Ftype=T, observed Ftype=%s", fld.Ftype)
		exceptions.Throw(exceptions.IllegalArgumentException, errMsg)
	}
	index := fld.Fvalue.(uint32)
	return *stringPool.GetStringPointer(index)
}

// UpdateObjectFromGoString : Set the value field of the given object to the given string
func UpdateObjectFromGoString(objPtr *Object, argString string) {
	index := stringPool.GetStringIndex(&argString)
	fld := Field{Ftype: types.StringIndex, Fvalue: index}
	objPtr.FieldTable["value"] = fld
}
