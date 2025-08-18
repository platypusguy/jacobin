package gfunction

import (
    "jacobin/src/globals"
    "jacobin/src/types"
    "runtime"
    "testing"
)

func TestAtomicLong_VMSupportsCS8_ReflectsArchitecture(t *testing.T) {
    globals.InitStringPool()

    // Determine expected support based on the same architecture map used in the implementation
    arch := runtime.GOARCH
    supportedArchitectures := map[string]bool{
        "amd64":    true,
        "arm64":    true,
        "ppc64":    true,
        "ppc64le":  true,
        "s390x":    true,
        "sparc64":  true,
        "mips64":   true,
        "mips64le": true,
    }
    expected := supportedArchitectures[arch]

    // Invoke the gfunction
    res := atomicLongVMSupportsCS8([]interface{}{})

    // Ensure it returns a Java boolean in int64 form
    val, ok := res.(int64)
    if !ok {
        t.Fatalf("VMSupportsCS8 did not return int64 (Java boolean), got %T", res)
    }

    // Compare to expected mapping
    if expected {
        if val != types.JavaBoolTrue {
            t.Fatalf("VMSupportsCS8 expected true for arch %s, got %d", arch, val)
        }
    } else {
        if val != types.JavaBoolFalse {
            t.Fatalf("VMSupportsCS8 expected false for arch %s, got %d", arch, val)
        }
    }
}
