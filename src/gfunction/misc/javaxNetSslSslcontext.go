package misc

import (
	"jacobin/src/gfunction/ghelpers"
)

func Load_Javax_Net_Ssl_SSLContext() {

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.<init>(Ljavax/net/ssl/SSLContextSpi;Ljava/security/Provider;Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.<init>(Ljavax/net/ssl/SSLContextSpi;Ljava/security/Provider;Ljava/lang/String;Z)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.createSSLEngine()Ljavax/net/ssl/SSLEngine;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.createSSLEngine(Ljava/lang/String;I)Ljavax/net/ssl/SSLEngine;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.engineCreateSSLEngine()Ljavax/net/ssl/SSLEngine;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.engineCreateSSLEngine(Ljava/lang/String;I)Ljavax/net/ssl/SSLEngine;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.engineGetClientSessionContext()Ljavax/net/ssl/SSLSessionContext;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.engineGetServerSessionContext()Ljavax/net/ssl/SSLSessionContext;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.engineGetSocketFactory()Ljavax/net/ssl/SSLSocketFactory;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.engineGetServerSocketFactory()Ljavax/net/ssl/SSLServerSocketFactory;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}
	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.getClientSessionContext()Ljavax/net/ssl/SSLSessionContext;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.getDefault()Ljavax/net/ssl/SSLContext;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.getDefaultSSLParameters()Ljavax/net/ssl/SSLParameters;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.getInstance(Ljava/lang/String;)Ljavax/net/ssl/SSLContext;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.getInstance(Ljava/lang/String;Ljava/lang/String;)Ljavax/net/ssl/SSLContext;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.getInstance(Ljava/lang/String;Ljava/security/Provider;)Ljavax/net/ssl/SSLContext;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.getProtocol()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.getProvider()Ljava/security/Provider;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.getServerSessionContext()Ljavax/net/ssl/SSLSessionContext;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.getServerSocketFactory()Ljavax/net/ssl/SSLServerSocketFactory;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.getSocketFactory()Ljavax/net/ssl/SSLSocketFactory;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.getSupportedSSLParameters()Ljavax/net/ssl/SSLParameters;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.init([Ljavax/net/ssl/KeyManager;[Ljavax/net/ssl/TrustManager;Ljava/security/SecureRandom;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["javax/net/ssl/SSLContext.setDefault(Ljavax/net/ssl/SSLContext;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}
}
