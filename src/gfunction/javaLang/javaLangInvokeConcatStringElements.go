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
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/trace"
	"jacobin/src/util"
	"slices"
)

func Load_Lang_Invoke_ConcatStringElements() {
	ghelpers.MethodSignatures["java/lang/invoke.ConcatStringElements()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots:   0,
			GFunction:    ConcatStringElements,
			NeedsContext: true,
		}
}

// ConcatStringElements concatenates multiple fragments into a single string.
// This is the universal call site method for string concatenation as handled
// in INVOKEDYNAMIC.
//
// It accepts only one parameter, which is the context data provided in the
// frame stack from which the frame of the calling method is located.
// That frame is the source of two different types of elements to be concatenated.
// 1) Dynamic elements previously pushed onto the operand stack immediately prior
// to the INVOKEDYNAMIC instruction;
// 2) Static elements that are found in the constant pool, pointed to by the
// INVOKEDYNAMIC instruction.
//
// Both types of elements can be retrieved by accessing the frame.
// One of the static arguments that is always present is the "recipe," which is
// a string-like series of bytes that has string data and embedded 0x01 and 0x02 bytes.
// The 0x01 byte specifies the location of a dynamic element, while the 0x02
// specifies the location of a static element.
func ConcatStringElements(params []any) any {
	if len(params) != 1 {
		errMsg := fmt.Sprintf("ConcatStringElements expected 1 argument, got %d", len(params))
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.MalformedParametersException, errMsg)
	}

	frameList, ok := params[0].(*list.List)
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

	// The two bytes following the INVOKEDYNAMIC opcode are an index into the CP index pointing to
	// the INVOKEDYNAMIC entry
	cpSlot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2])

	cp, ok := fr.CP.(*classloader.CPool)
	if !ok || cp == nil {
		errMsg := fmt.Sprintf("ConcatStringElements: invalid constant pool in frame for %s", fr.ClName)
		trace.Error(errMsg)
		return ghelpers.GetGErrBlk(excNames.ClassFormatError, errMsg)
	}

	invdynIndex := cp.CpIndex[cpSlot]
	invdyn := cp.InvokeDynamics[invdynIndex.Slot]
	invdynNAT := invdyn.NameAndType
	invdynTypeIndex := cp.CpIndex[invdynNAT]
	invdynType := cp.Utf8Refs[invdynTypeIndex.Type]
	invdynArgs := util.ParseIncomingParamsFromMethTypeStringRaw(invdynType) // the types of the dynamic args
	if globals.TraceInst {
		trace.Raw(fmt.Sprintf("invdynArgs: %v", invdynArgs))
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
	if len(bsm.Args) == 0 { // these should be the static args
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
	recipeBytes := []byte(recipe)
	_ = recipeObj

	const (
		tagArg   byte = 0x01 // '0x01' placeholder for dynamic arguments
		tagConst byte = 0x02 // '0x02' placeholder for static constant arguments
	)
	var j int = 0
	var argIndex int = 0
	concatBytes := make([]byte, len(recipeBytes))
	for i := 0; i < len(recipeBytes); i++ {
		switch recipeBytes[i] {
		case tagArg: // these arguments are pulled off op stack of the present frame
			arg := fr.OpStack[fr.TOS]
			fr.TOS -= 1
			argType := invdynArgs[argIndex]
			argIndex++
			if argType == "B" {
				concatBytes[j] = byte(arg.(int64))
			} else {
				s := object.StringifyAnythingGo(arg)
				concatBytes = slices.Concat(concatBytes, []byte(s))
				j = len(concatBytes) // and point j to the end of the new string
			}
		case tagConst:
			// Consume next constant from static constants array
		default:
			concatBytes[j] = recipeBytes[i]
			j++
		}
	}

	return object.StringObjectFromGoString(string(concatBytes))
}
