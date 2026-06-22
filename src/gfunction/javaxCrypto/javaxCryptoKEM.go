/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaxCrypto

import (
	"jacobin/src/gfunction/ghelpers"
)

func Load_Crypto_KEM() {

	// KEM class
	ghelpers.MethodSignatures["javax/crypto/KEM.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM.getInstance(Ljava/lang/String;)Ljavax/crypto/KEM;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM.getInstance(Ljava/lang/String;Ljava/lang/String;)Ljavax/crypto/KEM;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM.getInstance(Ljava/lang/String;Ljava/security/Provider;)Ljavax/crypto/KEM;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM.getAlgorithm()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM.getProvider()Ljava/security/Provider;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM.newDecapsulator(Ljava/security/PrivateKey;)Ljavax/crypto/KEM$Decapsulator;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM.newDecapsulator(Ljava/security/PrivateKey;Ljava/security/spec/AlgorithmParameterSpec;)Ljavax/crypto/KEM$Decapsulator;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM.newEncapsulator(Ljava/security/PublicKey;)Ljavax/crypto/KEM$Encapsulator;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM.newEncapsulator(Ljava/security/PublicKey;Ljava/security/spec/AlgorithmParameterSpec;Ljava/security/SecureRandom;)Ljavax/crypto/KEM$Encapsulator;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	// KEM$Encapsulated class
	ghelpers.MethodSignatures["javax/crypto/KEM$Encapsulated.<init>(Ljavax/crypto/SecretKey;[B[B)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM$Encapsulated.encapsulation()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM$Encapsulated.key()Ljavax/crypto/SecretKey;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM$Encapsulated.params()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	// KEM$Encapsulator class
	ghelpers.MethodSignatures["javax/crypto/KEM$Encapsulator.encapsulate()Ljavax/crypto/KEM$Encapsulated;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM$Encapsulator.encapsulate(IIILjava/lang/String;)Ljavax/crypto/KEM$Encapsulated;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM$Encapsulator.encapsulationSize()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM$Encapsulator.providerName()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM$Encapsulator.secretSize()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	// KEM$Decapsulator class
	ghelpers.MethodSignatures["javax/crypto/KEM$Decapsulator.decapsulate([B)Ljavax/crypto/SecretKey;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM$Decapsulator.decapsulate([BIIILjava/lang/String;)Ljavax/crypto/SecretKey;"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM$Decapsulator.encapsulationSize()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM$Decapsulator.providerName()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/crypto/KEM$Decapsulator.secretSize()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}
}
