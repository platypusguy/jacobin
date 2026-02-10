/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/src/gfunction/ghelpers"
)

func Load_Other_Methods() {

	ghelpers.MethodSignatures["java/awt/Color.initIDs()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/awt/Toolkit.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/awt/Toolkit.loadLibraries()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/io/FileDescriptor.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/io/BufferedInputStream.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/io/ByteArrayInputStream.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/lang/AbstractStringBuilder.ensureCapacityInternal(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/lang/CharacterDataLatin1.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/lang/ThreadLocal.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/lang/ThreadLocal.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/math/MathContext.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/nio/charset/Charset.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/nio/charset/Charset.defaultCharset()Ljava/nio/charset/Charset;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnCharsetName,
		}

	ghelpers.MethodSignatures["java/nio/charset/Charset.name()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnCharsetName,
		}

	ghelpers.MethodSignatures["java/security/AccessController.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/Locale$Category.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/Locale$Category.<init>(Ljava/lang/String;ILjava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 7,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["java/util/Regex.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["jdk/internal/access/SharedSecrets.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/VM.initialize()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/CDS.getRandomSeedForDumping()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnRandomLong,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/CDS.initializeFromArchive(Ljava/lang/Class;)V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/CDS.isDumpingArchive0()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnFalse,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/CDS.isDumpingClassList0()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnFalse,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/CDS.isSharingEnabled0()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnFalse,
		}

	ghelpers.MethodSignatures["jdk/internal/util/ArraysSupport.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["sun/nio/ch/SocketDispatcher.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["sun/security/util/Debug.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["sun/security/util/Debug.getInstance(Ljava/lang/String;)Lsun/security/util/Debug;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.ReturnNull,
		}

	ghelpers.MethodSignatures["sun/security/util/Debug.getInstance(Ljava/lang/String;Ljava/lang/String;)Lsun/security/util/Debug;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.ReturnNull,
		}

	ghelpers.MethodSignatures["sun/security/util/Debug.isOn(Ljava/lang/String;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.ReturnFalse,
		}

}
