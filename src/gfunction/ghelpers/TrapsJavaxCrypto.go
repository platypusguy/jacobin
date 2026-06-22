/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package ghelpers

func Load_Traps_Javax_Crypto() {

	MethodSignatures["javax/crypto/CipherSpi.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  ClinitGeneric,
		}

	MethodSignatures["javax/crypto/KeyAgreementSpi.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  ClinitGeneric,
		}

	MethodSignatures["javax/crypto/MacSpi.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  ClinitGeneric,
		}

	MethodSignatures["javax/crypto/SecretKeyFactorySpi.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  ClinitGeneric,
		}

	MethodSignatures["javax/crypto/ExemptionMechanismSpi.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  ClinitGeneric,
		}

	MethodSignatures["javax/crypto/KEMSpi.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  ClinitGeneric,
		}

	MethodSignatures["javax/crypto/KEMSpi$EncapsulatorSpi.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  ClinitGeneric,
		}

	MethodSignatures["javax/crypto/KEMSpi$DecapsulatorSpi.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  ClinitGeneric,
		}
}
