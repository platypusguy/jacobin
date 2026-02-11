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

// Other useful class names in alpha order
var ClassNameBigDecimal = "java/math/BigDecimal"
var ClassNameBigInteger = "java/math/BigInteger"
var ClassNameJavaLangClass = "java/lang/Class"
var ClassNameJcaProviderList = "sun/security/jca/ProviderList"
var ClassNameLinkedList = "java/util/LinkedList"
var ClassNameMessageDigest = "java/security/MessageDigest"
var ClassNameOptional string = "java/util/Optional"
var ClassNameProperties = "java/util/Properties"
var ClassNameThread = "java/lang/Thread"
var ClassNameThreadGroup = "java/lang/ThreadGroup"
var ClassNameThreadState = "java/lang/Thread$State"

var StringPoolThreadIndex = uint32(3)        // points to the string pool entry for "java/lang/Thread"
var StringPoolJavaLangClassIndex = uint32(4) // points to the string entry for "java/lang/Class"
var FieldNameProperties = "map"

// Security-related
const SecurityProviderName = "GoSecurityProvider"
const SecurityProviderInfo = "Security + Cryptography"
const SecurityProviderVersion = 1.0

var ClassNameDHParameterSpec = "javax/crypto/spec/DHParameterSpec"
var ClassNameDHPublicKey = "javax/crypto/interfaces/DHPublicKey"
var ClassNameDHPrivateKey = "javax/crypto/interfaces/DHPrivateKey"
var ClassNameDSAParameterSpec = "java/security/spec/DSAParameterSpec"
var ClassNameDSAPublicKey = "java/security/interfaces/DSAPublicKey"
var ClassNameDSAPrivateKey = "java/security/interfaces/DSAPrivateKey"
var ClassNameECGenParameterSpec = "java/security/spec/ECGenParameterSpec"
var ClassNameECPublicKey = "java/security/interfaces/ECPublicKey"
var ClassNameECPrivateKey = "java/security/interfaces/ECPrivateKey"
var ClassNameECParameterSpec = "java/security/ECParameterSpec"
var ClassNameECPoint = "java/security/ECPoint"
var ClassNameEd25519PrivateKey = "java/security/interfaces/Ed25519PrivateKey"
var ClassNameEd25519PublicKey = "java/security/interfaces/Ed25519PublicKey"
var ClassNameEd448PrivateKey = "java/security/interfaces/Ed448PrivateKey"
var ClassNameEd448PublicKey = "java/security/interfaces/Ed448PublicKey"
var ClassNameEllipticCurve = "java/security/EllipticCurve"
var ClassNameKeyPairGenerator = "java/security/KeyPairGenerator"
var ClassNameKeyPair = "java/security/KeyPair"
var ClassNamePrivateKey = "java/security/PrivateKey"
var ClassNamePublicKey = "java/security/PublicKey"
var ClassNameRSAKeyGenParameterSpec = "java/security/spec/RSAKeyGenParameterSpec"
var ClassNameRSAPrivateKey = "java/security/interfaces/RSAPrivateKey"
var ClassNameRSAPublicKey = "java/security/interfaces/RSAPublicKey"
var ClassNameSecurityProvider = "java/security/Provider"
var ClassNameSecurityProviderService = "java/security/Provider$Service"
var ClassNameSignature = "java/security/Signature"
var ClassNameSecureRandom = "java/security/SecureRandom"
var ClassNameX25519PrivateKey = "java/security/interfaces/X25519PrivateKey"
var ClassNameX25519PublicKey = "java/security/interfaces/X25519PublicKey"
var ClassNameX448PrivateKey = "java/security/interfaces/X448PrivateKey"
var ClassNameX448PublicKey = "java/security/interfaces/X448PublicKey"

// File system
var FileSystemProviderValue = &struct{}{}
var FileSystemProviderType = "Ljava/nio/file/spi/FileSystemProvider;"

// ---- experimental values ----
var StackInflator = 2 // for toying with whether to increase # of stack entries
