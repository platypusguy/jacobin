/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-3 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package util

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
    for i := 0; i < len(paramChars); i++ {
        switch paramChars[i] {
        case '(':
            continue
        case ')':
            return params
        case 'I', 'S', 'C', 'B', 'Z': // int, short, char, byte, bool -> int
            params = append(params, "I")
        case 'F':
            params = append(params, "F")
        case 'J':
            params = append(params, "J")
        case 'D':
            params = append(params, "D")
        case 'L':
            params = append(params, "L")
            for j := i + 1; j < len(paramChars); j++ {
                if paramChars[j] != ';' { // the end of the link is a ;
                    continue
                } else {
                    i = j // j now points to the ;, continue will add 1
                    break
                }
            }
        case '[': // arrays
            elements := make([]byte, 0)
            for paramChars[i] == '[' {
                elements = append(elements, '[')
                i += 1
            }
            // i is now pointing to the primitive in the array
            elements = append(elements, paramChars[i])
            params = append(params, string(elements))
        }
    }
    return params
}
