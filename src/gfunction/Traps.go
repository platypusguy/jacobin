/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/exceptions"
)

func Load_Traps() map[string]GMeth {

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

	MethodSignatures["java/nio/charset/Charset.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapCharset,
		}

	MethodSignatures["java/nio/channels/FileChannel.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFileChannel,
		}

	// Unsupported readers

	MethodSignatures["java/io/BufferedReader.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapReader,
		}

	MethodSignatures["java/io/CharArrayReader.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapReader,
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

	// Unsupported writers

	MethodSignatures["java/io/BufferedWriter.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapWriter,
		}

	MethodSignatures["java/io/CharArrayWriter.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapWriter,
		}

	MethodSignatures["java/io/FileSystem.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFileSystem,
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

	// String Builder

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

	// String Buffer

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

	return MethodSignatures
}

// Trap for Charset references
func trapCharset([]interface{}) interface{} {
	errMsg := "Class java/nio/charset/Charset is not yet supported"
	exceptions.Throw(exceptions.UnsupportedOperationException, errMsg)
	return getGErrBlk(exceptions.UnsupportedOperationException, errMsg)
}

// Trap for deprecated functions
func trapDeprecated([]interface{}) interface{} {
	errMsg := "The class or function requested is deprecated and is not supported by jacobin"
	exceptions.Throw(exceptions.UnsupportedOperationException, errMsg)
	return getGErrBlk(exceptions.UnsupportedOperationException, errMsg)
}

// Trap for FileChannel references
func trapFileChannel([]interface{}) interface{} {
	errMsg := "Class java.nio.channels.FileChannel is not yet supported"
	exceptions.Throw(exceptions.UnsupportedOperationException, errMsg)
	return getGErrBlk(exceptions.UnsupportedOperationException, errMsg)
}

// Trap for FileDescriptor references
func trapFileDescriptor([]interface{}) interface{} {
	errMsg := "Class java/io/FileDescriptor is not yet supported"
	exceptions.Throw(exceptions.UnsupportedOperationException, errMsg)
	return getGErrBlk(exceptions.UnsupportedOperationException, errMsg)
}

// Trap for FileSystem references
func trapFileSystem([]interface{}) interface{} {
	errMsg := "Class java.io.FileSystem is not yet supported"
	exceptions.Throw(exceptions.UnsupportedOperationException, errMsg)
	return getGErrBlk(exceptions.UnsupportedOperationException, errMsg)
}

// "java/io/DefaultFileSystem.getFileSystem()Ljava/io/FileSystem;"
func trapGetDefaultFileSystem([]interface{}) interface{} {
	errMsg := "DefaultFileSystem.getFileSystem() is not yet supported"
	exceptions.Throw(exceptions.UnsupportedOperationException, errMsg)
	return getGErrBlk(exceptions.UnsupportedOperationException, errMsg)
}

// Trap for unsupported readers
func trapReader([]interface{}) interface{} {
	errMsg := "The requested reader is not yet supported"
	exceptions.Throw(exceptions.UnsupportedOperationException, errMsg)
	return getGErrBlk(exceptions.UnsupportedOperationException, errMsg)
}

// Trap for unsupported writers
func trapWriter([]interface{}) interface{} {
	errMsg := "The requested writer is not yet supported"
	exceptions.Throw(exceptions.UnsupportedOperationException, errMsg)
	return getGErrBlk(exceptions.UnsupportedOperationException, errMsg)
}

// Trap for StringBuilder
func trapStringBuilder([]interface{}) interface{} {
	errMsg := "Class StringBuilder is not yet supported"
	exceptions.Throw(exceptions.UnsupportedOperationException, errMsg)
	return getGErrBlk(exceptions.UnsupportedOperationException, errMsg)
}

// Trap for StringBuilder
func trapStringBuffer([]interface{}) interface{} {
	errMsg := "Class StringBuffer is not yet supported"
	exceptions.Throw(exceptions.UnsupportedOperationException, errMsg)
	return getGErrBlk(exceptions.UnsupportedOperationException, errMsg)
}
