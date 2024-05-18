/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package types

// Grab bag of constants used in Jacobin

// ---- <clInit> status bytes ----
const NoClinit byte = 0x00
const ClInitNotRun byte = 0x01
const ClInitInProgress byte = 0x02
const ClInitRun byte = 0x03

// ---- invalid index into string pool ----
const InvalidStringIndex uint32 = 0xffffffff

// ---- default superclass ----
var ObjectClassName = "java/lang/Object"
var PtrToJavaLangObject = &ObjectClassName
var ObjectPoolStringIndex = uint32(2) // points to the string pool slice for "java/lang/Object"

// Constants related to "java/lang/String":
var StringClassName = "java/lang/String"
var StringClassRef = "Ljava/lang/String;"
var StringPoolStringIndex = uint32(1) // points to the string pool slice for "java/lang/String"
var EmptyString = ""
