/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

// do-nothing Go function shared by several source files
func justReturn([]interface{}) interface{} {
	return nil
}

func Load_Just_Return() {

	MethodSignatures["java/awt/Color.initIDs()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/awt/Toolkit.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/awt/Toolkit.loadLibraries()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/CharacterDataLatin1.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/math/BigDecimal.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/io/FileDescriptor.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/math/MathContext.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/math/RoundingMode.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/AbstractStringBuilder.ensureCapacityInternal(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkRead(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/System.registerNatives()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/util/HexFormat.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/util/Locale$Category.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/util/Locale$Category.<init>(Ljava/lang/String;ILjava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 6,
			GFunction:  justReturn,
		}

	MethodSignatures["java/util/Regex.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["jdk/internal/access/SharedSecrets.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

}
