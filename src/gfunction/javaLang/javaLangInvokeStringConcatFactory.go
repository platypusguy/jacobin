/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"container/list"
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/frames"
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
	return nil
}

// Concatenates a single parameter to an existing item.
// Requires the frame data so that it can reach into the INVOKEDYNAMIC's
// constant pool to find the first existing item. The second item is found
// on the op stack and the second arg to this function defines its type.
// NeedsContext = true (so last param is the frame stack)
func concat1param(params []any) any {
	if len(params) != 2 {
		errMsg := fmt.Sprintf("concat1param expected 2 args, got %d", len(params))
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.MalformedParametersException, errMsg)
	}

	frameList, ok := params[1].(*list.List)
	if !ok || frameList == nil || frameList.Len() == 0 {
		errMsg := "concat1param: invalid or empty frame stack in params[1]"
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	elem := frameList.Front()
	if elem == nil || elem.Value == nil {
		errMsg := "concat1param: no active frame at front of frame stack"
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	fr, ok := elem.Value.(*frames.Frame)
	if !ok || fr == nil {
		errMsg := "concat1param: front of frame stack is not a valid *frames.Frame"
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	if fr.Meth == nil || fr.PC < 0 || fr.PC+2 >= len(fr.Meth) {
		errMsg := fmt.Sprintf("concat1param: PC %d out of bounds for method bytecode", fr.PC)
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.ClassFormatError, errMsg)
	}

	// The two bytes following the INVOKEDYNAMIC opcode point to the CP entry
	cpSlot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2])

	cp, ok := fr.CP.(*classloader.CPool)
	if !ok || cp == nil {
		errMsg := fmt.Sprintf("concat1param: invalid constant pool in frame for %s", fr.ClName)
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.ClassFormatError, errMsg)
	}

	idEntry := classloader.FetchCPentry(cp, cpSlot)
	if idEntry.AddrVal == nil {
		errMsg := fmt.Sprintf("concat1param: invalid InvokeDynamic CP entry at index %d", cpSlot)
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.ClassFormatError, errMsg)
	}

	bsmIndex := int(idEntry.AddrVal.Entry1)
	if bsmIndex < 0 || bsmIndex >= len(cp.Bootstraps) {
		errMsg := fmt.Sprintf("concat1param: bsmIndex %d out of range", bsmIndex)
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.ClassFormatError, errMsg)
	}

	bsm := cp.Bootstraps[bsmIndex]
	if len(bsm.Args) == 0 {
		errMsg := fmt.Sprintf("concat1param: no bootstrap arguments for BSM index %d", bsmIndex)
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.ClassFormatError, errMsg)
	}

	recipeArgIndex := int(bsm.Args[0])
	recipeObj, err := classloader.ResolveStaticArgument(cp, recipeArgIndex, fr)
	if err != nil {
		errMsg := fmt.Sprintf("concat1param: error resolving recipe static argument at index %d: %s", recipeArgIndex, err.Error())
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.ClassFormatError, errMsg)
	}

	recipe := object.GoStringFromStringObject(recipeObj)
	_ = recipe
	_ = recipeObj

	return nil
}
