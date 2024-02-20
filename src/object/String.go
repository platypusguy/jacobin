/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package object

import (
	"jacobin/types"
)

// Strings are so commonly used in Java, that it makes sense
// to have a means of creating them quickly, rather than building
// them from scratch each time by walking through the constant pool.

var StringClassName = "java/lang/String"
var StringClassRef = "Ljava/lang/String;"
var EmptyString = ""

// NewString creates an empty string. However, it lacks an updated
// Klass field, which due to circularity issues, is updated in
// classloader.MakeString(). DO NOT CALL THIS FUNCTION DIRECTLYy.
// It should be called ONLY by classloader.MakeString
func NewString() *Object {
	s := new(Object)
	s.Mark.Hash = 0
	s.Klass = &StringClassName // java/lang/String
	s.FieldTable = make(map[string]Field)

	// ==== now the fields ====

	// value: the content of the string as array of bytes
	// Note: Post JDK9, this field is an array of bytes, so as to
	// enable compact strings.
	value := make([]byte, 10)
	// make value (the content of the string) Fields[0] and FieldTable["value"]
	valueField := Field{Ftype: types.ByteArray, Fvalue: &value}
	s.Fields = append(s.Fields, valueField)
	s.FieldTable["value"] = valueField

	// Field{Ftype: types.ByteArray, Fvalue: &value})

	// field 01 -- coder LATIN(=bytes, for compact strings) is 0; UTF16 is 1
	s.Fields = append(s.Fields, Field{Ftype: types.Byte, Fvalue: int64(1)})

	// field 02 -- string hash
	s.Fields = append(s.Fields, Field{Ftype: types.Int, Fvalue: int64(0)})

	// // field 03 -- COMPACT_STRINGS (always true for JDK >= 9)
	// s.Fields = append(s.Fields, Field{Ftype: "XZ", Fvalue: types.JavaBoolTrue})

	// // field 04 -- UTF_8.INSTANCE ptr to encoder
	// s.Fields = append(s.Fields, Field{Ftype: types.Ref, Fvalue: nil})

	// // field 05 -- ISO_8859_1.INSTANCE ptr to encoder
	// s.Fields = append(s.Fields, Field{Ftype: types.Ref, Fvalue: nil})

	// // field 06 -- sun/nio/cs/US_ASCII.INSTANCE
	// s.Fields = append(s.Fields, Field{Ftype: types.Ref, Fvalue: nil})

	// field 07 -- java/nio/charset/CodingErrorAction.REPLACE
	s.Fields = append(s.Fields, Field{Ftype: types.Ref, Fvalue: nil})

	// // field 08 -- java/lang/String.CASE_INSENSITIVE_ORDER
	// // points to a comparator. Will be useful to fill in later
	// s.Fields = append(s.Fields, Field{Ftype: types.Ref, Fvalue: nil})

	// field 09 -- hashIsZero (only true in rare case where hash is 0)
	s.Fields = append(s.Fields, Field{Ftype: types.Bool, Fvalue: types.JavaBoolFalse})

	// field 10 -- serialPersistentFields
	s.Fields = append(s.Fields, Field{Ftype: types.Ref, Fvalue: nil})

	return s
}

// NewStringFromGoString converts a go string to a Java string-like
// entity, in which the chars are stored as runes, rather than chars.
// TODO: it needs to determine whether a string can be stored as bytes or
// chars and set the flags in the String instance correctly.
func NewStringFromGoString(in string) *Object {
	s := NewString()
	s.Fields[0].Ftype = types.RuneArray
	s.Fields[0].Fvalue = in // test for compact strings and use GoStringToBytes() if on
	return s
}

// CreateCompactStringFromGoString creates a string in which the chars
// are stored as bytes--that is, a compact string.
func CreateCompactStringFromGoString(in *string) *Object {
	stringBytes := []byte(*in)
	s := NewString()

	// set the value of the string
	s.Fields[0].Ftype = types.ByteArray
	s.Fields[0].Fvalue = &stringBytes

	// set the string to LATIN
	s.Fields[1].Fvalue = int64(0)
	return s
}

// convenience method to extract a Go string from a Java string
func GetGoStringFromJavaStringPtr(strPtr *Object) string {
	s := *strPtr
	bytes := s.Fields[0].Fvalue.(*[]byte)
	return string(*bytes)
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

	return *objPtr.Klass == "java/lang/String"
}
