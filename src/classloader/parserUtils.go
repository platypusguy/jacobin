/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

// various utilities frequently used in parsing classfiles

// read two bytes in big endian order and convert to an int
func intFrom2Bytes(bytes []byte, pos int) (int, error) {
	if len(bytes) < pos+2 {
		return 0, cfe("invalid offset into file")
	}

	value := (uint16(bytes[pos]) << 8) + uint16(bytes[pos+1])
	return int(value), nil
}

// read four bytes in big endian order and convert to an int
func intFrom4Bytes(bytes []byte, pos int) (int, error) {
	if len(bytes) < pos+4 {
		return 0, cfe("invalid offset into file")
	}

	value1 := (uint32(bytes[pos]) << 8) + uint32(bytes[pos+1])
	value2 := (uint32(bytes[pos+2]) << 8) + uint32(bytes[pos+3])
	retVal := int(value1<<16) + int(value2)
	return retVal, nil
}
