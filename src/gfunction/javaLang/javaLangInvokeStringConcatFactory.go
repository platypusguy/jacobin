/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/trace"
	"jacobin/src/types"
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

// stringConcatFactory_makeConcatWithConstants returns a call site that includes a method handle for
// a method that concatenates strings with constants
func stringConcatFactory_makeConcatWithConstants(params []any) any {
	if len(params) != 5 {
		errMsg := fmt.Sprintf(
			"stringConcatFactory_makeConcatWithConstants: incorrect number of parameters: %d (should be 5)",
			len(params),
		)
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.MalformedParametersException, errMsg)
	}

	trace.Raw(fmt.Sprintf("calling class: %s", object.GoStringFromStringObject(params[0].(*object.Object))))
	trace.Raw(fmt.Sprintf("name of callsite method: %s", object.GoStringFromStringObject(params[1].(*object.Object))))
	// the method type is a java/lang/invoke/MethodType object with ptypes (parameter types) and
	// rtype (return type), which eventually must be extracted.
	ptypesArray := params[2].(*object.Object).FieldTable["ptypes"].Fvalue.(*object.Object)
	ptypes := ptypesArray.FieldTable["value"].Fvalue.([]*object.Object)
	for i, ptype := range ptypes {
		klassVal := ptype.FieldTable["$klass"].Fvalue
		if clData, ok := klassVal.(*classloader.ClData); ok && clData != nil {
			trace.Raw(fmt.Sprintf("ptype %d $klass: %s", i, clData.Name))
		} else {
			trace.Raw(fmt.Sprintf("ptype %d $klass: %v", i, klassVal))
		}
	}

	trace.Raw(fmt.Sprintf("recipe: %s", object.GoStringFromStringObject(params[3].(*object.Object))))
	constants := params[4].([]*object.Object)
	if len(constants) > 0 {
		trace.Raw("constants/args:")
		for _, constant := range constants {
			trace.Raw("\t" + object.StringifyAnythingGo(constant))
		}
	}

	// create a method handle for ConcatStringElements
	mh := createRawMethodHandleObject()
	mh.FieldTable["Kind"] = object.Field{Ftype: types.Int, Fvalue: REF_invokeStatic}
	mh.FieldTable["$target"] = object.Field{Ftype: types.Ref, Fvalue: classloader.MTentry{
		Meth:      ConcatStringElements,
		MType:     'G',
		MethClass: 0,
		MethName:  0,
		MethType:  0,
	}}
	// wrap the method handle in a call site
	// return the call site
	return nil
}
