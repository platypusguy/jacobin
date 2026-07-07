/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"strings"
)

func Load_Util_StringTokenizer() {
	ghelpers.MethodSignatures["java/util/StringTokenizer.<init>(Ljava/lang/String;Ljava/lang/String;Z)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  stringTokenizerInit,
		}

	ghelpers.MethodSignatures["java/util/StringTokenizer.<init>(Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  stringTokenizerInit2,
		}

	ghelpers.MethodSignatures["java/util/StringTokenizer.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringTokenizerInit1,
		}

	ghelpers.MethodSignatures["java/util/StringTokenizer.hasMoreTokens()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringTokenizerHasMoreTokens,
		}

	ghelpers.MethodSignatures["java/util/StringTokenizer.nextToken()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringTokenizerNextToken,
		}

	ghelpers.MethodSignatures["java/util/StringTokenizer.nextToken(Ljava/lang/String;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  stringTokenizerNextTokenWithDelims,
		}

	ghelpers.MethodSignatures["java/util/StringTokenizer.hasMoreElements()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringTokenizerHasMoreTokens,
		}

	ghelpers.MethodSignatures["java/util/StringTokenizer.nextElement()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringTokenizerNextToken,
		}

	ghelpers.MethodSignatures["java/util/StringTokenizer.countTokens()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  stringTokenizerCountTokens,
		}
}

type stringTokenizerState struct {
	str          string
	delims       string
	returnDelims bool
	position     int
	maxPosition  int
}

func getStringTokenizerState(self *object.Object) (*stringTokenizerState, interface{}) {
	field, exists := self.FieldTable["value"]
	if !exists {
		return nil, ghelpers.GetGErrBlk(excNames.NullPointerException, "StringTokenizer not initialized")
	}
	state, ok := field.Fvalue.(*stringTokenizerState)
	if !ok {
		return nil, ghelpers.GetGErrBlk(excNames.VirtualMachineError, "Invalid StringTokenizer state")
	}
	return state, nil
}

func stringTokenizerInit(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	strObj := params[1].(*object.Object)
	delimObj := params[2].(*object.Object)
	returnDelims := params[3].(int64) != 0

	if strObj == nil || delimObj == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Null string or delimiters")
	}

	str := object.GoStringFromStringObject(strObj)
	delims := object.GoStringFromStringObject(delimObj)

	state := &stringTokenizerState{
		str:          str,
		delims:       delims,
		returnDelims: returnDelims,
		position:     0,
		maxPosition:  len(str),
	}

	object.ClearFieldTable(self)
	self.FieldTable["value"] = object.Field{
		Ftype:  types.StringTokenizer,
		Fvalue: state,
	}

	return nil
}

func stringTokenizerInit2(params []interface{}) interface{} {
	newParams := make([]interface{}, 4)
	newParams[0] = params[0]
	newParams[1] = params[1]
	newParams[2] = params[2]
	newParams[3] = int64(0) // returnDelims = false
	return stringTokenizerInit(newParams)
}

func stringTokenizerInit1(params []interface{}) interface{} {
	newParams := make([]interface{}, 4)
	newParams[0] = params[0]
	newParams[1] = params[1]
	newParams[2] = object.StringObjectFromGoString(" \t\n\r\f")
	newParams[3] = int64(0) // returnDelims = false
	return stringTokenizerInit(newParams)
}

func stringTokenizerHasMoreTokens(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	state, err := getStringTokenizerState(self)
	if err != nil {
		return err
	}

	state.skipDelimiters()
	if state.position < state.maxPosition {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func (s *stringTokenizerState) skipDelimiters() {
	if s.returnDelims {
		return
	}
	for s.position < s.maxPosition && strings.ContainsRune(s.delims, rune(s.str[s.position])) {
		s.position++
	}
}

func stringTokenizerNextToken(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	state, err := getStringTokenizerState(self)
	if err != nil {
		return err
	}

	state.skipDelimiters()
	if state.position >= state.maxPosition {
		return ghelpers.GetGErrBlk(excNames.NoSuchElementException, "")
	}

	start := state.position
	if state.returnDelims && strings.ContainsRune(state.delims, rune(state.str[state.position])) {
		state.position++
	} else {
		for state.position < state.maxPosition && !strings.ContainsRune(state.delims, rune(state.str[state.position])) {
			state.position++
		}
	}

	return object.StringObjectFromGoString(state.str[start:state.position])
}

func stringTokenizerNextTokenWithDelims(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	delimObj := params[1].(*object.Object)

	if delimObj == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Null delimiters")
	}

	state, err := getStringTokenizerState(self)
	if err != nil {
		return err
	}

	state.delims = object.GoStringFromStringObject(delimObj)
	return stringTokenizerNextToken(params)
}

func stringTokenizerCountTokens(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	state, err := getStringTokenizerState(self)
	if err != nil {
		return err
	}

	count := int64(0)
	currPos := state.position

	for currPos < state.maxPosition {
		if !state.returnDelims {
			for currPos < state.maxPosition && strings.ContainsRune(state.delims, rune(state.str[currPos])) {
				currPos++
			}
		}

		if currPos >= state.maxPosition {
			break
		}

		if state.returnDelims && strings.ContainsRune(state.delims, rune(state.str[currPos])) {
			currPos++
		} else {
			for currPos < state.maxPosition && !strings.ContainsRune(state.delims, rune(state.str[currPos])) {
				currPos++
			}
		}
		count++
	}

	return count
}
