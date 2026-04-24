package javaxCrypto

import (
	"crypto/cipher"
	"encoding/binary"
	"fmt"
)

// RC2 is a block cipher designed by Ron Rivest.
// Implementation based on RFC 2268.

type rc2Cipher struct {
	xkey [64]uint16
}

const blockSize = 8

func newRC2Cipher(key []byte, bits int) (cipher.Block, error) {
	kLen := len(key)
	if kLen < 1 || kLen > 128 {
		return nil, fmt.Errorf("rc2: invalid key length %d", kLen)
	}
	if bits <= 0 {
		bits = kLen * 8
	}

	c := new(rc2Cipher)
	var l [128]byte
	copy(l[:], key)

	// Phase 1: Expand input key to 128 bytes
	for i := kLen; i < 128; i++ {
		l[i] = piTable[int(l[i-1]+l[i-kLen])&0xff]
	}

	// Phase 2: Reduce effective key length
	t8 := (bits + 7) / 8
	tm := 0xff >> (uint(8*t8) - uint(bits))
	i := 128 - t8
	l[i] = piTable[l[i]&byte(tm)]
	for i > 0 {
		i--
		l[i] = piTable[l[i]^l[i+t8]]
	}

	// Phase 3: Copy to xkey
	for i := 0; i < 64; i++ {
		c.xkey[i] = binary.LittleEndian.Uint16(l[i*2 : i*2+2])
	}

	return c, nil
}

func (c *rc2Cipher) BlockSize() int { return blockSize }

func (c *rc2Cipher) Encrypt(dst, src []byte) {
	if len(src) < blockSize {
		panic("rc2: input not full block")
	}
	if len(dst) < blockSize {
		panic("rc2: output not full block")
	}

	r0 := binary.LittleEndian.Uint16(src[0:2])
	r1 := binary.LittleEndian.Uint16(src[2:4])
	r2 := binary.LittleEndian.Uint16(src[4:6])
	r3 := binary.LittleEndian.Uint16(src[6:8])

	// Rounds 1-5
	for i := 0; i < 5; i++ {
		r0 = rotl(r0+c.xkey[i*4+0]+(r3&r2)+(^r3&r1), 1)
		r1 = rotl(r1+c.xkey[i*4+1]+(r0&r3)+(^r0&r2), 2)
		r2 = rotl(r2+c.xkey[i*4+2]+(r1&r0)+(^r1&r3), 3)
		r3 = rotl(r3+c.xkey[i*4+3]+(r2&r1)+(^r2&r0), 5)
	}

	// Mashing
	r0 += c.xkey[r3&63]
	r1 += c.xkey[r0&63]
	r2 += c.xkey[r1&63]
	r3 += c.xkey[r2&63]

	// Rounds 6-11
	for i := 5; i < 11; i++ {
		r0 = rotl(r0+c.xkey[i*4+0]+(r3&r2)+(^r3&r1), 1)
		r1 = rotl(r1+c.xkey[i*4+1]+(r0&r3)+(^r0&r2), 2)
		r2 = rotl(r2+c.xkey[i*4+2]+(r1&r0)+(^r1&r3), 3)
		r3 = rotl(r3+c.xkey[i*4+3]+(r2&r1)+(^r2&r0), 5)
	}

	// Mashing
	r0 += c.xkey[r3&63]
	r1 += c.xkey[r0&63]
	r2 += c.xkey[r1&63]
	r3 += c.xkey[r2&63]

	// Rounds 12-16
	for i := 11; i < 16; i++ {
		r0 = rotl(r0+c.xkey[i*4+0]+(r3&r2)+(^r3&r1), 1)
		r1 = rotl(r1+c.xkey[i*4+1]+(r0&r3)+(^r0&r2), 2)
		r2 = rotl(r2+c.xkey[i*4+2]+(r1&r0)+(^r1&r3), 3)
		r3 = rotl(r3+c.xkey[i*4+3]+(r2&r1)+(^r2&r0), 5)
	}

	binary.LittleEndian.PutUint16(dst[0:2], r0)
	binary.LittleEndian.PutUint16(dst[2:4], r1)
	binary.LittleEndian.PutUint16(dst[4:6], r2)
	binary.LittleEndian.PutUint16(dst[6:8], r3)
}

func (c *rc2Cipher) Decrypt(dst, src []byte) {
	if len(src) < blockSize {
		panic("rc2: input not full block")
	}
	if len(dst) < blockSize {
		panic("rc2: output not full block")
	}

	r0 := binary.LittleEndian.Uint16(src[0:2])
	r1 := binary.LittleEndian.Uint16(src[2:4])
	r2 := binary.LittleEndian.Uint16(src[4:6])
	r3 := binary.LittleEndian.Uint16(src[6:8])

	// Inverse Rounds 16-12
	for i := 15; i >= 11; i-- {
		r3 = rotr(r3, 5) - (c.xkey[i*4+3] + (r2 & r1) + (^r2 & r0))
		r2 = rotr(r2, 3) - (c.xkey[i*4+2] + (r1 & r0) + (^r1 & r3))
		r1 = rotr(r1, 2) - (c.xkey[i*4+1] + (r0 & r3) + (^r0 & r2))
		r0 = rotr(r0, 1) - (c.xkey[i*4+0] + (r3 & r2) + (^r3 & r1))
	}

	// Inverse Mashing
	r3 -= c.xkey[r2&63]
	r2 -= c.xkey[r1&63]
	r1 -= c.xkey[r0&63]
	r0 -= c.xkey[r3&63]

	// Inverse Rounds 11-6
	for i := 10; i >= 5; i-- {
		r3 = rotr(r3, 5) - (c.xkey[i*4+3] + (r2 & r1) + (^r2 & r0))
		r2 = rotr(r2, 3) - (c.xkey[i*4+2] + (r1 & r0) + (^r1 & r3))
		r1 = rotr(r1, 2) - (c.xkey[i*4+1] + (r0 & r3) + (^r0 & r2))
		r0 = rotr(r0, 1) - (c.xkey[i*4+0] + (r3 & r2) + (^r3 & r1))
	}

	// Inverse Mashing
	r3 -= c.xkey[r2&63]
	r2 -= c.xkey[r1&63]
	r1 -= c.xkey[r0&63]
	r0 -= c.xkey[r3&63]

	// Inverse Rounds 5-1
	for i := 4; i >= 0; i-- {
		r3 = rotr(r3, 5) - (c.xkey[i*4+3] + (r2 & r1) + (^r2 & r0))
		r2 = rotr(r2, 3) - (c.xkey[i*4+2] + (r1 & r0) + (^r1 & r3))
		r1 = rotr(r1, 2) - (c.xkey[i*4+1] + (r0 & r3) + (^r0 & r2))
		r0 = rotr(r0, 1) - (c.xkey[i*4+0] + (r3 & r2) + (^r3 & r1))
	}

	binary.LittleEndian.PutUint16(dst[0:2], r0)
	binary.LittleEndian.PutUint16(dst[2:4], r1)
	binary.LittleEndian.PutUint16(dst[4:6], r2)
	binary.LittleEndian.PutUint16(dst[6:8], r3)
}

func rotl(v uint16, n uint) uint16 {
	return (v << n) | (v >> (16 - n))
}

func rotr(v uint16, n uint) uint16 {
	return (v >> n) | (v << (16 - n))
}

var piTable = [256]byte{
	0xd9, 0x78, 0xf9, 0xc4, 0x19, 0xdd, 0xb5, 0xed, 0x28, 0xe9, 0x82, 0x79, 0x4a, 0x1a, 0x07, 0xad,
	0x64, 0x51, 0x04, 0x2e, 0xf1, 0x27, 0xb6, 0x71, 0x52, 0x83, 0x5a, 0x47, 0x13, 0x2d, 0x22, 0xeb,
	0x3b, 0x44, 0x1e, 0xc1, 0x0f, 0xcc, 0xe6, 0x81, 0x53, 0xaf, 0xa2, 0x24, 0xd0, 0x06, 0x3a, 0x93,
	0x86, 0x9e, 0x5d, 0x90, 0x41, 0x40, 0x75, 0xdf, 0x2f, 0xcb, 0x7b, 0xa0, 0x74, 0x0c, 0xc2, 0x58,
	0xcf, 0x3e, 0xb0, 0x18, 0x0a, 0x25, 0x23, 0x45, 0xd7, 0x05, 0x42, 0xed, 0x38, 0xd2, 0x16, 0xc8,
	0x46, 0x5f, 0x48, 0x72, 0x2a, 0x92, 0x01, 0x4b, 0x43, 0x56, 0xda, 0x14, 0x54, 0xfd, 0x68, 0x5c,
	0x33, 0xbb, 0x59, 0x09, 0x21, 0xfe, 0x69, 0x4c, 0xec, 0xe1, 0xbb, 0xd6, 0x29, 0x34, 0xa3, 0xf0,
	0xb3, 0x73, 0x0c, 0x1c, 0xd4, 0x02, 0x88, 0x61, 0x00, 0xa1, 0x1b, 0x60, 0x4d, 0xfb, 0x26, 0x0e,
	0x61, 0x67, 0x6b, 0x4e, 0x16, 0x7a, 0x6a, 0xc0, 0x35, 0xf7, 0xbc, 0xaf, 0xf2, 0x50, 0xac, 0xd3,
	0x7c, 0x2b, 0x21, 0xc3, 0x1a, 0x36, 0x08, 0xa2, 0xed, 0x62, 0xab, 0x31, 0x28, 0x47, 0x76, 0x24,
	0xe4, 0x70, 0xa0, 0xf4, 0xb1, 0x27, 0x22, 0x75, 0x0b, 0x2c, 0xa4, 0x7d, 0xb9, 0x1d, 0x11, 0x57,
	0xb4, 0x49, 0x21, 0x10, 0x00, 0xea, 0xfb, 0xbf, 0x6e, 0x2e, 0x7a, 0xd2, 0x6d, 0x3c, 0x51, 0xea,
	0xce, 0x0e, 0x5b, 0x52, 0x23, 0x3e, 0x20, 0x10, 0x31, 0xfc, 0x1d, 0x5e, 0xcd, 0x24, 0x3a, 0xb1,
	0x22, 0x2b, 0x6b, 0x6a, 0x11, 0xcf, 0x30, 0x67, 0x62, 0xa5, 0xf3, 0x3d, 0xa2, 0x70, 0x31, 0x11,
	0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
	0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
}
