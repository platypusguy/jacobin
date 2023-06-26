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

// NewString creates an empty string. However, it lacks an updated
// Klass field, which due to circularity issues, is updated in
// classloader.MakeString(). DO NOT CALL THIS FUNCTION DIRECTLYy.
// It should be called ONLY by classloader.MakeString
func NewString() *Object {
	s := new(Object)
	s.Mark.Hash = 0
	s.Klass = nil

	// ==== now the fields ====

	// field 00 -- value: the content of the string as array of chars
	// Note: Post JDK9, this field is an array of bytes, so as to
	// enable compact strings. Here, for better compatibility with
	// go, for the nonce we make it an array of Java chars's
	// equivalent in go: runes.
	array := make([]byte, 10)
	s.Fields = append(s.Fields,
		Field{Ftype: "[B", Fvalue: &array})

	// field 01 -- coder LATIN(=bytes, for compact strings) is 0; UTF16 is 1
	s.Fields = append(s.Fields, Field{Ftype: "B", Fvalue: int64(1)})

	// field 02 -- string hash
	s.Fields = append(s.Fields, Field{Ftype: "I", Fvalue: int64(0)})

	// field 03 -- COMPACT_STRINGS (always true for JDK >= 9)
	s.Fields = append(s.Fields, Field{Ftype: "XZ", Fvalue: types.JavaBoolTrue})

	// field 04 -- UTF_8.INSTANCE ptr to encoder
	s.Fields = append(s.Fields, Field{Ftype: "L", Fvalue: nil})

	// field 05 -- ISO_8859_1.INSTANCE ptr to encoder
	s.Fields = append(s.Fields, Field{Ftype: "L", Fvalue: nil})

	// field 06 -- sun/nio/cs/US_ASCII.INSTANCE
	s.Fields = append(s.Fields, Field{Ftype: "L", Fvalue: nil})

	// field 07 -- java/nio/charset/CodingErrorAction.REPLACE
	s.Fields = append(s.Fields, Field{Ftype: "L", Fvalue: nil})

	// field 08 -- java/lang/String.CASE_INSENSITIVE_ORDER
	// points to a comparator. Will be useful to fill in later
	s.Fields = append(s.Fields, Field{Ftype: "L", Fvalue: nil})

	// field 09 -- hashIsZero (only true in rare case where hash is 0
	s.Fields = append(s.Fields, Field{Ftype: "Z", Fvalue: types.JavaBoolFalse})

	// field 10 -- serialPersistentFields
	s.Fields = append(s.Fields, Field{Ftype: "L", Fvalue: nil})

	return s
}

// Functions for creating strings are in classloader/makeString.go due to Go's
// circularity prohibition.
