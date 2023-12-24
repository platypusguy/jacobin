/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-3 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package util

import (
	"fmt"
	"jacobin/log"
	"jacobin/types"
)

// ParseIncomingParamsFromMethTypeString takes a type string from a CP
// and parses its passed-in parameters, returning them in reduced form
// as a slice. By reduced, we mean, for example, ints, shorts, chars, etc.
// are all marked as ints.
func ParseIncomingParamsFromMethTypeString(s string) []string {
	params := make([]string, 0)
	if s == "" {
		return params
	}

	paramChars := []byte(s)
	paramLen := len(paramChars)
	for i := 0; i < len(paramChars); i++ {
		switch paramChars[i] {
		case '(':
			continue
		case ')':
			return params
		case 'I', 'S', 'C', 'B', 'Z': // int, short, char, byte, bool -> int
			params = append(params, types.Int)
		case 'F':
			params = append(params, types.Float)
		case 'J':
			params = append(params, types.Long)
		case 'D':
			params = append(params, types.Double)
		case 'L':
			params = append(params, types.Ref)
			var j int
			for j = i + 1; j < paramLen; j++ {
				if paramChars[j] != ';' { // the end of the link is a ;
					continue
				} else {
					i = j // j now points to the ;, continue will add 1
					break
				}
			}
			if j >= paramLen {
				errMsg := fmt.Sprintf("ParseIncomingParamsFromMethTypeString case 'L': failed to find final ';'")
				_ = log.Log(errMsg, log.SEVERE)
				return make([]string, 0)
			}
		case '[': // arrays
			elements := make([]byte, 0)
			for paramChars[i] == '[' && i < paramLen {
				elements = append(elements, '[')
				i += 1
			}
			if i >= paramLen {
				errMsg := fmt.Sprintf("ParseIncomingParamsFromMethTypeString case '[': unending '[' repetitions")
				_ = log.Log(errMsg, log.SEVERE)
				return make([]string, 0)
			}
			if paramChars[i] == ')' {
				errMsg := fmt.Sprintf("ParseIncomingParamsFromMethTypeString case '[': no array type specified")
				_ = log.Log(errMsg, log.SEVERE)
				return make([]string, 0)
			}
			// i is now pointing to the type-character of the array
			elements = append(elements, paramChars[i])
			params = append(params, string(elements))
			if paramChars[i] == 'L' {
				var j int
				for j = i + 1; j < paramLen; j++ {
					if paramChars[j] != ';' { // the end of the link is a ;
						continue
					} else {
						i = j // j now points to the ;, continue will add 1
						break
					}
				}
				if j >= paramLen {
					errMsg := fmt.Sprintf("ParseIncomingParamsFromMethTypeString case '[': failed to find final ';'")
					_ = log.Log(errMsg, log.SEVERE)
					return make([]string, 0)
				}
			}
		default:
			errMsg := fmt.Sprintf("ParseIncomingParamsFromMethTypeString default: illegal character '%c'", paramChars[i])
			_ = log.Log(errMsg, log.SEVERE)
			return make([]string, 0)
		}
	}
	return params
}
