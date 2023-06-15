/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package javaTypes

// bytes in Go are uint8, whereas in Java they are int8. Hence this type alias.
type JavaByte = int8

// booleans in Java are defined as integer values of 0 and 1
// in arrays, they're stored as bytes, everywhere else as 32-bit ints.
// Jacobin uses 64-bit ints.

var JavaBoolTrue int64 = 1
var JavaBoolFalse int64 = 0
var JavaBool int64
