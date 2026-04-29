/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaIo

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
)

func Load_Io_InputStream() {
	ghelpers.MethodSignatures["java/io/InputStream.available()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  inputStreamAvailable,
		}

	ghelpers.MethodSignatures["java/io/InputStream.close()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/InputStream.mark(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/InputStream.markSupported()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/InputStream.nullInputStream()Ljava/io/InputStream;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/InputStream.read()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/InputStream.read([B)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  inputStreamReadIntoByteArray,
		}

	ghelpers.MethodSignatures["java/io/InputStream.read([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/InputStream.readAllBytes()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/InputStream.readNBytes(I)[B"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/InputStream.readNBytes([BII)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/InputStream.reset()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/InputStream.skip(J)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/InputStream.skipNBytes(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/io/InputStream.transferTo(Ljava/io/OutputStream;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}
}

func inputStreamAvailable(params []any) any {
	return int64(0)
}

// java/io/InputStream.read([B)I reads stdin contents to a byte array
// returns number of bytes processed or -1 if EOF
func inputStreamReadIntoByteArray(params []any) any {
	if object.IsNull(params[1]) {
		errMsg := "java.lang.io.inputStream.read() called with null array"
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	byteArray := params[1].(*object.Object).FieldTable["value"].Fvalue.([]types.JavaByte)
	if len(byteArray) == 0 {
		return int64(0)
	}

	// for the time being, we return -1 to indicate EOF
	// TODO: add the system interface to read from stdin
	// for both this function and available()
	return int64(-1)
}
