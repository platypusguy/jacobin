package gfunction

func Load_Javax_Net_Ssl_SSLContext() {

	MethodSignatures["javax/net/ssl/SSLContext.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["javax/net/ssl/SSLContext.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.<init>(Ljavax/net/ssl/SSLContextSpi;Ljava/security/Provider;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.<init>(Ljavax/net/ssl/SSLContextSpi;Ljava/security/Provider;Ljava/lang/String;Z)V"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.createSSLEngine()Ljavax/net/ssl/SSLEngine;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.createSSLEngine(Ljava/lang/String;I)Ljavax/net/ssl/SSLEngine;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.engineCreateSSLEngine()Ljavax/net/ssl/SSLEngine;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.engineCreateSSLEngine(Ljava/lang/String;I)Ljavax/net/ssl/SSLEngine;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.engineGetClientSessionContext()Ljavax/net/ssl/SSLSessionContext;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.engineGetServerSessionContext()Ljavax/net/ssl/SSLSessionContext;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.engineGetSocketFactory()Ljavax/net/ssl/SSLSocketFactory;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.engineGetServerSocketFactory()Ljavax/net/ssl/SSLServerSocketFactory;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}
	MethodSignatures["javax/net/ssl/SSLContext.getClientSessionContext()Ljavax/net/ssl/SSLSessionContext;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.getDefault()Ljavax/net/ssl/SSLContext;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.getDefaultSSLParameters()Ljavax/net/ssl/SSLParameters;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.getInstance(Ljava/lang/String;)Ljavax/net/ssl/SSLContext;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.getInstance(Ljava/lang/String;Ljava/lang/String;)Ljavax/net/ssl/SSLContext;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.getInstance(Ljava/lang/String;Ljava/security/Provider;)Ljavax/net/ssl/SSLContext;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.getProtocol()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.getProvider()Ljava/security/Provider;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.getServerSessionContext()Ljavax/net/ssl/SSLSessionContext;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.getServerSocketFactory()Ljavax/net/ssl/SSLServerSocketFactory;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.getSocketFactory()Ljavax/net/ssl/SSLSocketFactory;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.getSupportedSSLParameters()Ljavax/net/ssl/SSLParameters;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.init([Ljavax/net/ssl/KeyManager;[Ljavax/net/ssl/TrustManager;Ljava/security/SecureRandom;)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["javax/net/ssl/SSLContext.setDefault(Ljavax/net/ssl/SSLContext;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}
}
