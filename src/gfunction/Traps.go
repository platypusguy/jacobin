/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/excNames"
)

func Load_Traps() {

	MethodSignatures["java/io/BufferedInputStream.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/io/BufferedOutputStream.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/io/BufferedWriter.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/io/CharArrayReader.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/io/CharArrayWriter.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/io/DefaultFileSystem.getFileSystem()Ljava/io/FileSystem;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/io/FileSystem.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/io/FilterReader.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/io/PipedReader.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/io/StringReader.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/io/FilterOutputStream.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/io/FilterOutputStream.<init>(Ljava/io/OutputStream;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/io/FilterWriter.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/io/PipedWriter.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/io/PrintWriter.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/io/StringWriter.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/lang/SecurityManager.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/lang/SecurityManager.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/nio/ByteBuffer.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/nio/charset/StandardCharsets.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/nio/channels/AsynchronousFileChannel.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/nio/channels/FileChannel.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

	MethodSignatures["java/rmi/RMISecurityManager.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/rmi/RMISecurityManager.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["sun/security/util/Debug.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapClass,
		}

}

// Generic trap for classes
func trapClass([]interface{}) interface{} {
	errMsg := "TRAP: The requested class is not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Generic trap for deprecated classes and functions
func trapDeprecated([]interface{}) interface{} {
	errMsg := "TRAP: The requested class or function is deprecated and, therefore, not supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Generic trap for deprecated classes and functions
func trapUndocumented([]interface{}) interface{} {
	errMsg := "TRAP: The requested class or function is undocumented and, therefore, not supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Generic trap for functions
func trapFunction([]interface{}) interface{} {
	errMsg := "TRAP: The requested function is not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}
