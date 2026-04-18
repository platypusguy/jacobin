package javaSecurity

import (
	"jacobin/src/gfunction/ghelpers"
)

func Load_Security_SecretKey() {
	ghelpers.MethodSignatures["javax/crypto/SecretKey.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["javax/crypto/SecretKey.getAlgorithm()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  keyGetAlgorithm,
		}

	ghelpers.MethodSignatures["javax/crypto/SecretKey.getEncoded()[B"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  keyGetEncoded,
		}

	ghelpers.MethodSignatures["javax/crypto/SecretKey.getFormat()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  keyGetFormat,
		}

}
