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
			GFunction:  trapBufferedInputStream,
		}

	MethodSignatures["java/io/BufferedOutputStream.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapBufferedOutputStream,
		}

	MethodSignatures["java/io/BufferedWriter.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapWriter,
		}

	MethodSignatures["java/io/CharArrayReader.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapReader,
		}

	MethodSignatures["java/io/CharArrayWriter.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapWriter,
		}

	MethodSignatures["java/io/DefaultFileSystem.getFileSystem()Ljava/io/FileSystem;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapGetDefaultFileSystem,
		}

	MethodSignatures["java/io/FileDescriptor.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFileDescriptor,
		}

	MethodSignatures["java/io/FileSystem.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFileSystem,
		}

	MethodSignatures["java/io/FilterReader.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapReader,
		}

	MethodSignatures["java/io/PipedReader.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapReader,
		}

	MethodSignatures["java/io/StringReader.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapReader,
		}

	MethodSignatures["java/io/FilterInputStream.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFilterInputStream,
		}

	MethodSignatures["java/io/FilterInputStream.<init>(Ljava/io/InputStream;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFilterInputStream,
		}

	MethodSignatures["java/io/FilterOutputStream.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFilterOutputStream,
		}

	MethodSignatures["java/io/FilterOutputStream.<init>(Ljava/io/OutputStream;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFilterOutputStream,
		}

	MethodSignatures["java/io/FilterWriter.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapWriter,
		}

	MethodSignatures["java/io/PipedWriter.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapWriter,
		}

	MethodSignatures["java/io/PrintWriter.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapWriter,
		}

	MethodSignatures["java/io/StringWriter.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapWriter,
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

	MethodSignatures["java/lang/StringBuffer.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapStringBuffer,
		}

	MethodSignatures["java/lang/StringBuffer.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapStringBuffer,
		}

	MethodSignatures["java/lang/StringBuffer.<init>(I)V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapStringBuffer,
		}

	MethodSignatures["java/lang/StringBuffer.<init>(Ljava/lang/CharSequence;)V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapStringBuffer,
		}

	MethodSignatures["java/lang/StringBuffer.<init>(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapStringBuffer,
		}

	MethodSignatures["java/lang/StringBuilder.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapStringBuilder,
		}

	MethodSignatures["java/lang/StringBuilder.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapStringBuilder,
		}

	MethodSignatures["java/lang/StringBuilder.<init>(I)V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapStringBuilder,
		}

	MethodSignatures["java/lang/StringBuilder.<init>(Ljava/lang/CharSequence;)V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapStringBuilder,
		}

	MethodSignatures["java/lang/StringBuilder.<init>(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapStringBuilder,
		}

	MethodSignatures["java/nio/charset/Charset.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapCharset,
		}

	MethodSignatures["java/nio/channels/AsynchronousFileChannel.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFileChannel,
		}

	MethodSignatures["java/nio/channels/FileChannel.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFileChannel,
		}

	MethodSignatures["java/security/AccessController.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/security/AccessController.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/security/SecureRandom.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapSecureRandom,
		}

	MethodSignatures["java/security/SecureRandom.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapSecureRandom,
		}

	MethodSignatures["java/security/SecureRandom.<init>([B)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapSecureRandom,
		}

	MethodSignatures["java/security/SecureRandom.<init>(Ljava.security.SecureRandomSpi;Ljava.security.Provider;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapSecureRandom,
		}

	MethodSignatures["java/util/Random.next(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapRandomNext,
		}

	MethodSignatures["jdk/internal/access/SharedSecrets.<clinit>()V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapSharedSecrets,
		}

}

// Trap for BufferedInputStream references
func trapBufferedInputStream([]interface{}) interface{} {
	errMsg := "trapBufferedInputStream: Class java/io/BufferedInputStream is not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Trap for BufferedOutputStream references
func trapBufferedOutputStream([]interface{}) interface{} {
	errMsg := "trapBufferedOutputStream: Class java/io/BufferedOutputStream is not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Trap for Charset references
func trapCharset([]interface{}) interface{} {
	errMsg := "trapCharset: Class java/nio/charset/Charset is not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Trap for deprecated functions
func trapDeprecated([]interface{}) interface{} {
	errMsg := "trapDeprecated: The class or function requested is deprecated and is not supported by jacobin"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Trap for FileChannel references
func trapFileChannel([]interface{}) interface{} {
	errMsg := "trapFileChannel: File Channels are not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Trap for FileDescriptor references
func trapFileDescriptor([]interface{}) interface{} {
	errMsg := "trapFileDescriptor: Class java/io/FileDescriptor is not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Trap for FileSystem references
func trapFileSystem([]interface{}) interface{} {
	errMsg := "trapFileSystem: Class java.io.FileSystem is not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// "java/io/DefaultFileSystem.getFileSystem()Ljava/io/FileSystem;"
func trapGetDefaultFileSystem([]interface{}) interface{} {
	errMsg := "trapGetDefaultFileSystem: DefaultFileSystem.getFileSystem() is not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Trap for unsupported readers
func trapReader([]interface{}) interface{} {
	errMsg := "trapReader: The requested reader is not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Trap for unsupported writers
func trapWriter([]interface{}) interface{} {
	errMsg := "trapWriter: The requested writer is not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Trap for StringBuilder
func trapStringBuilder([]interface{}) interface{} {
	errMsg := "trapStringBuilder: Class StringBuilder is not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Trap for StringBuffer
func trapStringBuffer([]interface{}) interface{} {
	errMsg := "trapStringBuffer: Class StringBuffer is not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Trap for FilterInputStream
func trapFilterInputStream([]interface{}) interface{} {
	errMsg := "trapFilterInputStream: Class FilterInputStream is not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Trap for FilterOutputStream
func trapFilterOutputStream([]interface{}) interface{} {
	errMsg := "trapFilterOutputStream: Class FilterOutputStream is not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Trap for Random.next()
func trapRandomNext([]interface{}) interface{} {
	errMsg := "trapRandomNext: Protected method Random.next should never be reached unless done by reflection"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Trap for Random.next()
func trapSharedSecrets([]interface{}) interface{} {
	errMsg := "trapSharedSecrets: Class jdk/internal/access/SharedSecrets is not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}

// Trap for SecureRandom
func trapSecureRandom([]interface{}) interface{} {
	errMsg := "trapSecureRandom: Class java.security.SecureRandom is not yet supported"
	return getGErrBlk(excNames.UnsupportedOperationException, errMsg)
}
