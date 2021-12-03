/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package exec

// ParseIncomingParamsFromMethTypeString takes a type string from a CP
// and parses its passed-in parameters, returning them in reduced form
// as a slice. By reduced, we mean, for example, ints, shorts, chars, etc.
// are all marked as ints.
func ParseIncomingParamsFromMethTypeString(s string) []byte {
	params := make([]byte, 0)
	if s == "" {
		return params
	}

	paramChars := []byte(s)
	for i := 0; i < len(paramChars); i++ {
		switch paramChars[i] {
		case '(':
			continue
		case ')':
			return params
		case 'I', 'S', 'C', 'B', 'Z': // int, short, char, byte, bool -> int
			params = append(params, 'I')
		case 'F':
			params = append(params, 'F')
		case 'J':
			params = append(params, 'J')
		case 'D':
			params = append(params, 'D')
		case 'L', '[': // objects and arrays -> object references
			params = append(params, 'L')
		}
	}
	return params
}
