/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/statics"
	"jacobin/src/types"
)

func Load_Lang_Invoke_LambdaMetafactory() {
	ghelpers.MethodSignatures["java/lang/invoke/LambdaMetafactory.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  lambdaMetafactoryClinit,
		}

	ghelpers.MethodSignatures["java/lang/invoke/LambdaMetafactory.metafactory(Ljava/lang/invoke/MethodHandles$Lookup;Ljava/lang/String;Ljava/lang/invoke/MethodType;Ljava/lang/invoke/MethodType;Ljava/lang/invoke/MethodHandle;Ljava/lang/invoke/MethodType;)Ljava/lang/invoke/CallSite;"] =
		ghelpers.GMeth{
			ParamSlots: 6,
			GFunction:  lambdaMetafactoryClinit,
		}

	ghelpers.MethodSignatures["java/lang/invoke/LambdaMetafactory.altMetafactory(Ljava/lang/invoke/MethodHandles$Lookup;Ljava/lang/String;Ljava/lang/invoke/MethodType;[Ljava/lang/Object;)Ljava/lang/invoke/CallSite;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  lambdaMetafactoryClinit,
		}
}

func lambdaMetafactoryClinit(_ []any) any {
	_ = statics.AddStatic("java/lang/invoke/LambdaMetafactory.FLAG_BRIDGES",
		statics.Static{
			Type:  types.Int,
			Value: 4,
		})

	_ = statics.AddStatic("java/lang/invoke/LambdaMetafactory.FLAG_MARKERS",
		statics.Static{
			Type:  types.Int,
			Value: 2,
		})

	_ = statics.AddStatic("java/lang/invoke/LambdaMetafactory.FLAG_SERIALIZABLE",
		statics.Static{
			Type:  types.Int,
			Value: 1,
		})

	return nil
}
