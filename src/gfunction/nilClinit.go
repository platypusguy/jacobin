/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

func Load_Nil_Clinit() map[string]GMeth {

	MethodSignatures["java/math/BigDecimal.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/math/BigInteger.<clinit>()V"] =
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

	return MethodSignatures
}
