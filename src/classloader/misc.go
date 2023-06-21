/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import "jacobin/object"

// Miscellaneous functions (often placed here because of circularity issues0

// MakeString retuns an empty String object. It would normally be placed
// in String.go, except that the updating of the Klass field as done here
// creates a circularity error
func MakeString() *object.Object {
	strPtr := object.NewString()
	strPtr.Klass = MethAreaFetch("java/lang/String")
	return strPtr
}

// NewStringFromGoString converts a go string to a Java string
// converting the individual chars from runees to bytes (if
// compact strings are enabled) or UTF-16 values if not.
func NewStringFromGoString(in string) *object.Object {
	s := MakeString()
	s.Fields[0].Fvalue = in // test for compact strings and use GoStringToBytes() if on
	return s
}

func CreateJavaStringFromGoString(in *string) *object.Object {
	// stringBytes := GoStringToJavaBytes(*ins)
	stringBytes := []byte(*in)
	s := MakeString()
	// set the value of the string
	s.Fields[0].Ftype = "[B"
	s.Fields[0].Fvalue = &stringBytes
	// set the string to LATIN
	s.Fields[1].Fvalue = int64(0)
	return s
}
