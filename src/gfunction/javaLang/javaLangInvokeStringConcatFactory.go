/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/trace"
)

func Load_Lang_Invoke_StringConcatFactory() {
	ghelpers.MethodSignatures["java/lang/invoke/StringConcatFactory.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/lang/invoke/StringConcatFactory.makeConcat(Ljava/lang/invoke/MethodHandles$Lookup;Ljava/lang/String;Ljava/lang/invoke/MethodType;)Ljava/lang/invoke/CallSite;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/invoke/StringConcatFactory.makeConcatWithConstants(Ljava/lang/invoke/MethodHandles$Lookup;Ljava/lang/String;Ljava/lang/invoke/MethodType;Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/invoke/CallSite;"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  stringConcatFactory_makeConcatWithConstants,
		}

	ghelpers.MethodSignatures["java/lang/invoke/StringConcatFactory.makeConcatWithTemplate(Ljava/util/List;Ljava/util/List;)Ljava/lang/invoke/MethodHandle;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapPreview,
		}

	ghelpers.MethodSignatures["java/lang/invoke/StringConcatFactory.makeConcatWithTemplateCluster(Ljava/util/List;Ljava/util/List;I)Ljava/util/List;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapPreview,
		}

	ghelpers.MethodSignatures["java/lang/invoke/StringConcatFactory.makeConcatWithTemplateGetters(Ljava/util/List;Ljava/util/List;I)Ljava/lang/invoke/MethodHandle;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapPreview,
		}
}

func stringConcatFactory_makeConcatWithConstants(params []any) any {
	if len(params) != 5 {
		errMsg := fmt.Sprintf(
			"stringConcatFactory_makeConcatWithConstants: incorrect number of parameters: %d (should be 5)",
			len(params),
		)
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.MalformedParametersException, errMsg)
	}

	trace.Raw(fmt.Sprintf("name of method: %s", object.GoStringFromStringObject(params[1].(*object.Object))))

	return nil
}
