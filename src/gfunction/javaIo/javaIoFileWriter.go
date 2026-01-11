/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaIo

import (
	"jacobin/src/gfunction/ghelpers"
)

func Load_Io_FileWriter() {

	ghelpers.MethodSignatures["java/io/FileWriter.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/io/FileWriter.<init>(Ljava/io/File;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  initFileOutputStreamFile,
		}

	ghelpers.MethodSignatures["java/io/FileWriter.<init>(Ljava/io/File;Z)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  initFileOutputStreamFileBoolean,
		}

	ghelpers.MethodSignatures["java/io/FileWriter.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  initFileOutputStreamString,
		}

	ghelpers.MethodSignatures["java/io/FileWriter.<init>(Ljava/lang/String;Z)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  initFileOutputStreamStringBoolean,
		}

	ghelpers.MethodSignatures["java/io/FileWriter.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  oswClose,
		}

	ghelpers.MethodSignatures["java/io/FileWriter.flush()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  oswFlush,
		}

	ghelpers.MethodSignatures["java/io/FileWriter.write(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  oswWriteOneChar,
		}

	ghelpers.MethodSignatures["java/io/FileWriter.write([CII)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  oswWriteCharBuffer,
		}

	ghelpers.MethodSignatures["java/io/FileWriter.write(Ljava/lang/String;II)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  oswWriteStringBuffer,
		}

	// -----------------------------------------
	// traps that do nothing but return an error
	// -----------------------------------------

	ghelpers.MethodSignatures["java/io/FileWriter.<init>(Ljava/io/File;Ljava/lang.String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileWriter.<init>(Ljava/io/FileDescriptor;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileWriter.<init>(Ljava/io/File;Ljava/nio/charset/Charset;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileWriter.<init>(Ljava/io/File;Ljava/nio/charset/Charset;Z)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileWriter.<init>(Ljava/lang/String;Ljava/nio/charset/Charset;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/FileWriter.<init>(Ljava/lang/String;Ljava/nio/charset/Charset;Z)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

}
