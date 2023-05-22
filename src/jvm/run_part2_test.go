/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

// Additional tests for run.go. The previous 4K LOC of tests
// appear in run_test.go

package jvm

import (
    "testing"
)

func TestConvertInterfaceToUint64(t *testing.T) {
    var i64 int64 = 200
    // var f64 float64 = 345.0
    // var ref unsafe.Pointer = unsafe.Pointer(&f64)

    ret := convertInterfaceToUint64(i64)
    if ret != 200 {
        t.Errorf("Expected TestConvertInterfaceToUint64() to retun 200, got %d\n",
            ret)
    }

    // ret = convertInterfaceToUint64(f64)
    // if ret != 345 {
    // 	t.Errorf("Expected TestConvertInterfaceToUint64() to retun 345, got %d\n",
    // 		ret)
    // }
}
