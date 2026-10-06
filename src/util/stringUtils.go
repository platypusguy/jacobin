/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package util

// TEST
// DecodeModifiedUTF8 decodes Modified UTF-8 null: 0xC0 0x80 → 0x00
func DecodeModifiedUTF8(baIn []byte) []byte {
	var baOut []byte
	for i := 0; i < len(baIn); {
		// Modified UTF-8 null: 0xC0 0x80 → 0x00
		if baIn[i] == 0xC0 && i+1 < len(baIn) && baIn[i+1] == 0x80 {
			baOut = append(baOut, 0x00)
			i += 2
		} else {
			baOut = append(baOut, baIn[i])
			i++
		}
	}
	return baOut
}
