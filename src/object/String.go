/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package object

import "jacobin/classloader"

// Strings are so commonly used in Java, that it makes sense
// to have a means of creating them quickly, rather than building
// them from scratch each time by walking through the constant pool.

func NewString() *Object {
    s := new(Object)
    s.Mark.Hash = 0
    k := classloader.Classes["java/lang/String"]
    s.Klass = &k

    // ==== now the fields ====

    // field 01 -- value: the content of the string as array of chars
    // Note: Post JDK9, this field is an array of bytes, so as to
    // enable compact strings. Here, for better compatibility with
    // go, for the nonce we make it an array of Java chars's
    // equivalent in go: runes.
    s.Fields = append(s.Fields, Field{Ftype: "[C", Fvalue: make([]rune, 0)})

    // field 02 -- coder LATIN(=bytes, for compact strings) is 0; UTF16 is 1
    s.Fields = append(s.Fields, Field{Ftype: "B", Fvalue: 1})

    // field 03 -- string hash
    s.Fields = append(s.Fields, Field{Ftype: "I", Fvalue: 0})

    // field 04 -- COMPACT_STRINGS (always true for JDK >= 9)
    s.Fields = append(s.Fields, Field{Ftype: "Z", Fvalue: true})

    // field 05 -- UTF_8.INSTANCE ptr to encoder
    s.Fields = append(s.Fields, Field{Ftype: "L", Fvalue: nil})

    // field 06 -- ISO_8859_1.INSTANCE ptr to encoder
    s.Fields = append(s.Fields, Field{Ftype: "L", Fvalue: nil})

    // field 07 -- sun/nio/cs/US_ASCII.INSTANCE
    s.Fields = append(s.Fields, Field{Ftype: "L", Fvalue: nil})

    // field 08 -- java/nio/charset/CodingErrorAction.REPLACE
    s.Fields = append(s.Fields, Field{Ftype: "L", Fvalue: nil})

    // field 09 -- java/lang/String.CASE_INSENSITIVE_ORDER
    // points to a comparator. Will be useful to fill in later
    s.Fields = append(s.Fields, Field{Ftype: "L", Fvalue: nil})

    // field 10 -- hashIsZero (only true in rare case where hash is 0
    s.Fields = append(s.Fields, Field{Ftype: "Z", Fvalue: false})

    // field 11 -- serialPersistentFields
    s.Fields = append(s.Fields, Field{Ftype: "L", Fvalue: nil})

    return s
}

// NewStringFromGoString converts a go string to a Java string
// converting the individual chars from runees to bytes (if
// compact strings are enabled) or UTF-16 values if not.
func NewStringFromGoString(in string) *Object {
    s := NewString()
    s.Fields[0].Fvalue = in // test for compact strings and use GoStringToBytes() if on
    return s
}

// Convert go string (consiting of 32-bit runes aka chars) to
// single-byte values--for use in compact strings
func GoStringToBytes(in string) []byte {
    bytes := make([]byte, len(in))
    runes := []rune(in)
    for i := 0; i < len(in); i++ {
        r := runes[i]
        b := byte(r)
        bytes[i] = b
    }
    return bytes
}
