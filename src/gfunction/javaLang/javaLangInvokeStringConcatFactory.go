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

// ConcatStringElements concatenates multiple fragments into a single string.
// This is the universal call site method for string concatenation as handled
// in INVOKEDYNAMIC.
//
// It accepts only one parameter, which is the context data provided in the
// executing frame stack from which the frame of the calling method is located.
// That frame is the source of two different types of elements to be concatenated.
// 1) Dynamic elements are pushed onto the operand stack immediately prior to the
// INVOKEDYNAMIC instruction;
// 2) Static elements are found in the constant pool, pointed to by the INVOKEDYNAMIC
// instruction.
//
// Both types of elements can be retrieved by accessing the frame.
// One of the static arguments that is always present is the "recipe," which is
// a string-like series of bytes that has string data and embedded 0x01 and 0x02 bytes.
// The 0x01 byte specifies the location of a dynamic element, while the 0x02
// specifies the location of a static element.
func ConcatStringElements(params []any) any {
	if len(params) != 2 {
		errMsg := fmt.Sprintf("ConcatStringElements expected 2 args, got %d", len(params))
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.MalformedParametersException, errMsg)
	}

	frameList, ok := params[1].(*list.List)
	if !ok || frameList == nil || frameList.Len() == 0 {
		errMsg := "ConcatStringElements: invalid or empty frame stack in params[1]"
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	elem := frameList.Front()
	if elem == nil || elem.Value == nil {
		errMsg := "ConcatStringElements: no active frame at front of frame stack"
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	fr, ok := elem.Value.(*frames.Frame)
	if !ok || fr == nil {
		errMsg := "ConcatStringElements: front of frame stack is not a valid *frames.Frame"
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	if fr.Meth == nil || fr.PC < 0 || fr.PC+2 >= len(fr.Meth) {
		errMsg := fmt.Sprintf("ConcatStringElements: PC %d out of bounds for method bytecode", fr.PC)
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.ClassFormatError, errMsg)
	}

	// The two bytes following the INVOKEDYNAMIC opcode point to the CP entry
	cpSlot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2])

	cp, ok := fr.CP.(*classloader.CPool)
	if !ok || cp == nil {
		errMsg := fmt.Sprintf("ConcatStringElements: invalid constant pool in frame for %s", fr.ClName)
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.ClassFormatError, errMsg)
	}

	idEntry := classloader.FetchCPentry(cp, cpSlot)
	if idEntry.AddrVal == nil {
		errMsg := fmt.Sprintf("ConcatStringElements: invalid InvokeDynamic CP entry at index %d", cpSlot)
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.ClassFormatError, errMsg)
	}

	bsmIndex := int(idEntry.AddrVal.Entry1)
	if bsmIndex < 0 || bsmIndex >= len(cp.Bootstraps) {
		errMsg := fmt.Sprintf("ConcatStringElements: bsmIndex %d out of range", bsmIndex)
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.ClassFormatError, errMsg)
	}

	bsm := cp.Bootstraps[bsmIndex]
	if len(bsm.Args) == 0 {
		errMsg := fmt.Sprintf("ConcatStringElements: no bootstrap arguments for BSM index %d", bsmIndex)
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.ClassFormatError, errMsg)
	}

	recipeArgIndex := int(bsm.Args[0])
	recipeObj, err := classloader.ResolveStaticArgument(cp, recipeArgIndex, fr)
	if err != nil {
		errMsg := fmt.Sprintf("ConcatStringElements: error resolving recipe static argument at index %d: %s", recipeArgIndex, err.Error())
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.ClassFormatError, errMsg)
	}

	recipe := object.GoStringFromStringObject(recipeObj)
	_ = recipe
	_ = recipeObj

	return nil
}
