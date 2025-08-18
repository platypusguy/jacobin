/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-3 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package util

import (
	"jacobin/src/types"
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
			// validate reference descriptor terminates with ';' before appending
			var j int
			for j = i + 1; j < paramLen && paramChars[j] != ';'; j++ {
				// scan to terminating ';'
			}
			if j >= paramLen {
				// invalid reference descriptor: missing terminating ';'
				return make([]string, 0)
			}
			params = append(params, types.Ref)
			i = j
		case '[': // arrays
			elements := make([]byte, 0)
			for i < paramLen && paramChars[i] == '[' {
				elements = append(elements, '[')
				i++
			}
			if i >= paramLen {
				return make([]string, 0)
			}
			if paramChars[i] == ')' {
				return make([]string, 0)
			}
			// validate reference arrays fully, then append
			var j int
			if paramChars[i] == 'L' {
				for j = i + 1; j < paramLen && paramChars[j] != ';'; j++ {
				}
				if j >= paramLen {
					return make([]string, 0)
				}
			}
			elements = append(elements, paramChars[i])
			params = append(params, string(elements))
			if paramChars[i] == 'L' {
				i = j
			}
		default:
			// illegal or unexpected character encountered
			return make([]string, 0)
		}
	}
	return params
}
