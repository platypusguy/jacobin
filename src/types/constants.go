/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-4 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package types

// Grab bag of constants used in Jacobin

// ---- <clInit> status bytes ----
const NoClInit byte = 0x00
const ClInitNotRun byte = 0x01
const ClInitInProgress byte = 0x02
const ClInitRun byte = 0x03

// ---- invalid index into string pool ----
const InvalidStringIndex uint32 = 0xffffffff

// ---- default superclass ----
var ObjectClassName = "java/lang/Object"
var ObjectArrayClassName = "[java/lang/Object"
var PtrToJavaLangObject = &ObjectClassName
var StringPoolObjectIndex = uint32(2) // points to the string pool slice for "java/lang/Object"

// Constants related to "java/lang/String":
var StringClassName = "java/lang/String"
var StringArrayClassName = "[java/lang/String"
var StringClassRef = "Ljava/lang/String;"
var StringPoolStringIndex = uint32(1) // points to the string pool slice for "java/lang/String"

// Misc
var ModuleClassRef = "Ljava/lang/Module;"
var EmptyString = ""
var NullString = "null"

// Other class names
var ClassNameThread = "java/lang/Thread"
var ClassNameThreadGroup = "java/lang/ThreadGroup"
var ClassNameThreadState = "java/lang/Thread$State"
var StringPoolThreadIndex = uint32(3) // points to the string pool slice for "java/lang/Thread"
var ClassNameLinkedList = "java/util/LinkedList"
var ClassNameProperties = "java/util/Properties"
var FieldNameProperties = "map"
var ClassNameOptional string = "java/util/Optional"
var ClassNameBigDecimal = "java/math/BigDecimal"
var ClassNameBigInteger = "java/math/BigInteger"
var ClassNameSecurityProvider = "java/security/Provider"

// Security-related
const SecurityProviderName = "GoSecurityProvider"
const SecurityProviderInfo = "Security + Cryptography"

// File system
var FileSystemProviderValue = &struct{}{}
var FileSystemProviderType = "Ljava/nio/file/spi/FileSystemProvider;"

// ---- experimental values ----
var StackInflator = 2 // for toying with whether to increase # of stack entries
