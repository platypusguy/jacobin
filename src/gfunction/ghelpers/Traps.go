/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package ghelpers

import (
	"jacobin/src/excNames"
)

func Load_Traps() {

	MethodSignatures["java/awt/image/BufferedImage.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  TrapClass,
		}

	MethodSignatures["java/awt/image/BufferedImage.<init>(III)Ljava/awt/image/BufferedImage;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  TrapFunction,
		}

	MethodSignatures["java/awt/image/BufferedImage.<init>(IIILjava/awt/image/IndexColorModel;)Ljava/awt/image/BufferedImage;"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  TrapFunction,
		}

	MethodSignatures["java/awt/image/BufferedImage.<init>(Ljava/awt/image/ColorModel;Ljava/awt/image/WritableRaster;ZLjava/util/Hashtable;)Ljava/awt/image/BufferedImage;"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  TrapFunction,
		}

	MethodSignatures["java/awt/Image.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  TrapClass,
		}

	MethodSignatures["java/awt/Image.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  TrapFunction,
		}

	MethodSignatures["java/awt/ImageCapabilities.<init>(Z)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  TrapFunction,
		}

	MethodSignatures["java/rmi/RMISecurityManager.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  TrapDeprecated,
		}

	MethodSignatures["java/rmi/RMISecurityManager.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  TrapDeprecated,
		}

	MethodSignatures["java/security/AccessController.doPrivileged(Ljava/security/PrivilegedAction;)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  TrapDeprecated,
		}

	MethodSignatures["java/sql/Driver.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  TrapClass,
		}

	MethodSignatures["java/sql/DriverAction.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  TrapClass,
		}

	MethodSignatures["java/sql/DriverPropertyInfo.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  TrapClass,
		}

	MethodSignatures["java/sql/DriverPropertyInfo.<init>(Ljava/lang/String;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  TrapClass,
		}

	MethodSignatures["java/sql/DriverManager.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  TrapClass,
		}

	MethodSignatures["java/util/zip/CheckedInputStream.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  TrapClass,
		}

	MethodSignatures["java/util/zip/CheckedInputStream.<init>(Ljava/io/InputStream;Ljava/util/zip/Checksum;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  TrapClass,
		}

}

// TrapClass is a generic Trap for classes
func TrapClass([]interface{}) interface{} {
	errMsg := "TRAP: The requested class is not yet supported"
	return GetGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// TrapDeprecated is a generic Trap for deprecated classes and functions
func TrapDeprecated([]interface{}) interface{} {
	errMsg := "TRAP: The requested class or function is deprecated and, therefore, not supported"
	return GetGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Generic trap for functions
func TrapFunction([]interface{}) interface{} {
	errMsg := "TRAP: The requested function is not yet supported"
	return GetGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

func TrapKeyPairGeneration([]interface{}) interface{} {
	errMsg := "TRAP: Use KeyPairGenerator to create public and private keys"
	return GetGErrBlk(excNames.SecurityException, errMsg)
}

// TrapProtected is a generic Trap for functions
func TrapProtected([]interface{}) interface{} {
	errMsg := "TRAP: The requested function is protected"
	return GetGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// TrapUndocumented is a generic Trap for deprecated classes and functions
func TrapUndocumented([]interface{}) interface{} {
	errMsg := "TRAP: The requested class or function is undocumented and, therefore, not supported"
	return GetGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

func TrapUnicode([]interface{}) interface{} {
	return GetGErrBlk(
		excNames.UnsupportedOperationException,
		"Character Unicode method not yet implemented",
	)
}
