/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package types

import "testing"

func TestFourBytesToInt64_Table(t *testing.T) {
    // Table-driven tests: b1..b4 -> expected int64
    tests := []struct {
        b1, b2, b3, b4 byte
        want           int64
        name           string
    }{
        {0x00, 0x00, 0x00, 0x00, 0, "zero"},
        {0x00, 0x00, 0x00, 0x01, 1, "one"},
        {0x12, 0x34, 0x56, 0x78, 0x12345678, "0x12345678"},
        {0x7f, 0xff, 0xff, 0xff, 2147483647, "max_int32"},
        {0x80, 0x00, 0x00, 0x00, -2147483648, "min_int32"},
        {0xff, 0xff, 0xff, 0xff, -1, "minus_one"},
    }

    for _, tc := range tests {
        got := FourBytesToInt64(tc.b1, tc.b2, tc.b3, tc.b4)
        if got != tc.want {
            t.Fatalf("%s: FourBytesToInt64(%#02x,%#02x,%#02x,%#02x) = %d; want %d",
                tc.name, tc.b1, tc.b2, tc.b3, tc.b4, got, tc.want)
        }
    }
}
