/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package types

import "encoding/binary"

// converts four bytes into a signed 64-bit integer
func FourBytesToInt64(b1, b2, b3, b4 byte) int64 {
	wbytes := make([]byte, 8)
	wbytes[4] = b1
	wbytes[5] = b2
	wbytes[6] = b3
	wbytes[7] = b4

	if (b1 & 0x80) == 0x80 { // Negative bite value (left-most bit on)?
		// Negative byte - need to extend the sign (left-most) bit
		wbytes[0] = 0xff
		wbytes[1] = 0xff
		wbytes[2] = 0xff
		wbytes[3] = 0xff
	}
	return int64(binary.BigEndian.Uint64(wbytes))
}
