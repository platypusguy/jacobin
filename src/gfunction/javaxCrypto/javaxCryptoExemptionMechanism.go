/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaxCrypto

import (
	"jacobin/src/gfunction/ghelpers"
)

func Load_Crypto_ExemptionMechanism() {

	ghelpers.MethodSignatures["javax/crypto/ExemptionMechanism.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["javax/crypto/ExemptionMechanism.genExemptionBlob()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/ExemptionMechanism.genExemptionBlob([B)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/ExemptionMechanism.genExemptionBlob([BI)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/ExemptionMechanism.getAlgorithm()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/ExemptionMechanism.getInstance(Ljava/lang/String;)Ljavax/crypto/ExemptionMechanism;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/ExemptionMechanism.getInstance(Ljava/lang/String;Ljava/lang/String;)Ljavax/crypto/ExemptionMechanism;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/ExemptionMechanism.getInstance(Ljava/lang/String;Ljava/security/Provider;)Ljavax/crypto/ExemptionMechanism;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/ExemptionMechanism.getOutputSize(I)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/ExemptionMechanism.getProvider()Ljava/security/Provider;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/ExemptionMechanism.init(Ljava/security/Key;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/ExemptionMechanism.init(Ljava/security/Key;Ljava/security/AlgorithmParameters;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/ExemptionMechanism.init(Ljava/security/Key;Ljava/security/spec/AlgorithmParameterSpec;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/ExemptionMechanism.isCryptoAllowed(Ljava/security/Key;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}
}
