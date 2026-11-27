/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by \ the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import (
	"errors"
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/exceptions"
	"jacobin/src/frames"
	"jacobin/src/gfunction"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaLang"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"runtime/debug"
	"strings"
)

// 0xBA INVOKEDYNAMIC
func doInvokeDynamic(fr *frames.Frame, _ int64) int {
	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2]) // next 2 bytes point to CP entry

	// the components of the InvokeDynamic entry are validated by codeCheck prior to getting here
	// so we can call the callsite resolution method without further verifying the entries.
	CP := fr.CP.(*classloader.CPool)

	// Fast-path: Check if the BSM is StringConcatFactory.makeConcatWithConstants
	// If so, just call that method directly without going through the BSM processing
	if isStringConcatBSM(CP, CPslot) {
		res := javaLang.ConcatStringElements([]any{fr.FrameStack})
		if errBlk, isErr := res.(*ghelpers.GErrBlk); isErr {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			status := exceptions.ThrowEx(errBlk.ExceptionType, errBlk.ErrMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED
			}
			return RESUME_HERE
		}

		// Push resulting String object onto the caller frame's operand stack
		push(fr, res)
		return 5 // 2 bytes CP slot + 2 zero bytes + 1 next bytecode
	}

	callSite, err := classloader.ResolveCallSite(CP, CPslot, fr)
	if callSite == nil || err != nil {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("INVOKEDYNAMIC: error resolving callsite: %s", err.Error())
		status := exceptions.ThrowEx(excNames.ClassFormatError, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}
	methHandle := callSite.FieldTable["target"].Fvalue.(*object.Object)
	mtEntry := methHandle.FieldTable["$target"].Fvalue.(classloader.MTentry)

	// presently, all callsite method handles point to gfunctions
	gmethData := mtEntry.Meth.(ghelpers.GMeth)
	paramCount := gmethData.ParamSlots
	var params []any
	for range paramCount {
		params = append(params, pop(fr))
	}

	ret := gfunction.RunGfunction(mtEntry, fr.FrameStack, &params, true, globals.TraceInst)
	if ret != nil {
		switch ret.(type) {
		case error:
			if globals.GetGlobalRef().JacobinName == "test" {
				return ERROR_OCCURRED
			} else if errors.Is(ret.(error), gfunction.CaughtGfunctionException) {
				return RESUME_HERE // resume at the present PC, which points to the exception code
			}
		default:
			methodType := *stringPool.GetStringPointer(mtEntry.MethType)
			if !strings.HasSuffix(methodType, "V") { // if it's not an error,
				// then it's a legitimate return value, which we push provided
				// the method does not return void
				push(fr, ret)
			}
		}
	}

	return 5 // the two bytes for the CP slot + 2 bytes with value 0x00 + 1 for next bytecode
}
