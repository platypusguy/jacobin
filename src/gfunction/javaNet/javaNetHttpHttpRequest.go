/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaNet

import (
	"jacobin/src/gfunction/ghelpers"
)

func Load_Net_Http_HttpRequest() {

	// java.net.http.HttpRequest
	ghelpers.MethodSignatures["java/net/http/HttpRequest.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest.bodyPublisher()Ljava/util/Optional;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest.expectContinue()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest.headers()Ljava/net/http/HttpHeaders;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest.method()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest.newBuilder()Ljava/net/http/HttpRequest$Builder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest.newBuilder(Ljava/net/URI;)Ljava/net/http/HttpRequest$Builder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest.newBuilder(Ljava/net/http/HttpRequest;Ljava/util/function/BiPredicate;)Ljava/net/http/HttpRequest$Builder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest.timeout()Ljava/util/Optional;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest.uri()Ljava/net/URI;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest.version()Ljava/util/Optional;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	// java.net.http.HttpRequest$Builder
	ghelpers.MethodSignatures["java/net/http/HttpRequest$Builder.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$Builder.DELETE()Ljava/net/http/HttpRequest$Builder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$Builder.GET()Ljava/net/http/HttpRequest$Builder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$Builder.HEAD()Ljava/net/http/HttpRequest$Builder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$Builder.POST(Ljava/net/http/HttpRequest$BodyPublisher;)Ljava/net/http/HttpRequest$Builder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$Builder.PUT(Ljava/net/http/HttpRequest$BodyPublisher;)Ljava/net/http/HttpRequest$Builder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$Builder.build()Ljava/net/http/HttpRequest;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$Builder.copy()Ljava/net/http/HttpRequest$Builder;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$Builder.expectContinue(Z)Ljava/net/http/HttpRequest$Builder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$Builder.header(Ljava/lang/String;Ljava/lang/String;)Ljava/net/http/HttpRequest$Builder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$Builder.headers([Ljava/lang/String;)Ljava/net/http/HttpRequest$Builder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$Builder.method(Ljava/lang/String;Ljava/net/http/HttpRequest$BodyPublisher;)Ljava/net/http/HttpRequest$Builder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$Builder.setHeader(Ljava/lang/String;Ljava/lang/String;)Ljava/net/http/HttpRequest$Builder;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$Builder.timeout(Ljava/time/Duration;)Ljava/net/http/HttpRequest$Builder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$Builder.uri(Ljava/net/URI;)Ljava/net/http/HttpRequest$Builder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$Builder.version(Ljava/net/http/HttpClient$Version;)Ljava/net/http/HttpRequest$Builder;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	// java.net.http.HttpRequest$BodyPublisher
	ghelpers.MethodSignatures["java/net/http/HttpRequest$BodyPublisher.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$BodyPublisher.contentLength()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	// java.net.http.HttpRequest$BodyPublishers
	ghelpers.MethodSignatures["java/net/http/HttpRequest$BodyPublishers.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$BodyPublishers.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$BodyPublishers.concat([Ljava/net/http/HttpRequest$BodyPublisher;)Ljava/net/http/HttpRequest$BodyPublisher;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$BodyPublishers.fromPublisher(Ljava/util/concurrent/Flow$Publisher;)Ljava/net/http/HttpRequest$BodyPublisher;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$BodyPublishers.fromPublisher(Ljava/util/concurrent/Flow$Publisher;J)Ljava/net/http/HttpRequest$BodyPublisher;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$BodyPublishers.noBody()Ljava/net/http/HttpRequest$BodyPublisher;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$BodyPublishers.ofByteArray([B)Ljava/net/http/HttpRequest$BodyPublisher;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$BodyPublishers.ofByteArray([BII)Ljava/net/http/HttpRequest$BodyPublisher;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$BodyPublishers.ofByteArrays(Ljava/lang/Iterable;)Ljava/net/http/HttpRequest$BodyPublisher;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$BodyPublishers.ofFile(Ljava/nio/file/Path;)Ljava/net/http/HttpRequest$BodyPublisher;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$BodyPublishers.ofInputStream(Ljava/util/function/Supplier;)Ljava/net/http/HttpRequest$BodyPublisher;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$BodyPublishers.ofString(Ljava/lang/String;)Ljava/net/http/HttpRequest$BodyPublisher;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/net/http/HttpRequest$BodyPublishers.ofString(Ljava/lang/String;Ljava/nio/charset/Charset;)Ljava/net/http/HttpRequest$BodyPublisher;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}
}
