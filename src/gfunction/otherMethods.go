/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

func Load_Other_methods() {

	MethodSignatures["java/awt/Color.initIDs()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/awt/Toolkit.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/awt/Toolkit.loadLibraries()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/io/FileDescriptor.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/lang/AbstractStringBuilder.ensureCapacityInternal(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/CharacterDataLatin1.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/lang/SecurityManager.checkRead(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/math/BigDecimal.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/math/MathContext.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/math/RoundingMode.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/nio/charset/Charset.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/nio/charset/Charset.defaultCharset()Ljava/nio/charset/Charset;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  returnCharsetName,
		}

	MethodSignatures["java/nio/charset/Charset.name()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  returnCharsetName,
		}

	MethodSignatures["java/util/Locale$Category.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/util/Locale$Category.<init>(Ljava/lang/String;ILjava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 7,
			GFunction:  justReturn,
		}

	MethodSignatures["java/util/Regex.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["jdk/internal/access/SharedSecrets.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["jdk/internal/misc/VM.initialize()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["jdk/internal/misc/CDS.getRandomSeedForDumping()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  returnRandomLong,
		}

	MethodSignatures["jdk/internal/misc/CDS.initializeFromArchive(Ljava/lang/Class;)V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["jdk/internal/misc/CDS.isDumpingArchive0()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  returnFalse,
		}

	MethodSignatures["jdk/internal/misc/CDS.isDumpingClassList0()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  returnFalse,
		}

	MethodSignatures["jdk/internal/misc/CDS.isSharingEnabled0()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  returnFalse,
		}

	MethodSignatures["jdk/internal/util/ArraysSupport.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["sun/security/util/Debug.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/io/ByteArrayInputStream.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/io/BufferedInputStream.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

}
