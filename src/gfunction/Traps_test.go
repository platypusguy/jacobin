/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
    "reflect"
    "strings"
    "testing"

    "jacobin/excNames"
)

func TestLoad_Traps_RegistersSomeMethods(t *testing.T) {
    saved := MethodSignatures
    defer func() { MethodSignatures = saved }()
    MethodSignatures = make(map[string]GMeth)

    Load_Traps()

    // Representative subset across class, function, deprecated
    checks := []struct{
        key   string
        slots int
        fn    func([]interface{}) interface{}
    }{
        {"java/io/BufferedOutputStream.<clinit>()V", 0, trapClass},
        {"java/io/DefaultFileSystem.getFileSystem()Ljava/io/FileSystem;", 0, trapFunction},
        {"java/io/FilterOutputStream.<init>(Ljava/io/OutputStream;)V", 1, trapFunction},
        {"java/rmi/RMISecurityManager.<clinit>()V", 0, trapDeprecated},
        {"java/rmi/RMISecurityManager.<init>()V", 0, trapDeprecated},
    }

    for _, c := range checks {
        got, ok := MethodSignatures[c.key]
        if !ok {
            t.Fatalf("missing MethodSignatures entry for %s", c.key)
        }
        if got.ParamSlots != c.slots {
            t.Fatalf("%s ParamSlots expected %d, got %d", c.key, c.slots, got.ParamSlots)
        }
        if got.GFunction == nil {
            t.Fatalf("%s GFunction expected non-nil", c.key)
        }
        if reflect.ValueOf(got.GFunction).Pointer() != reflect.ValueOf(c.fn).Pointer() {
            t.Fatalf("%s GFunction mismatch", c.key)
        }
    }
}

func TestTrapFunctions_ReturnUnsupported(t *testing.T) {
    // trapClass
    if blk, ok := trapClass(nil).(*GErrBlk); !ok || blk.ExceptionType != excNames.UnsupportedOperationException || !strings.Contains(blk.ErrMsg, "TRAP:") {
        t.Fatalf("trapClass expected UnsupportedOperationException with TRAP: message, got %+v", blk)
    }
    // trapFunction
    if blk, ok := trapFunction(nil).(*GErrBlk); !ok || blk.ExceptionType != excNames.UnsupportedOperationException || !strings.Contains(blk.ErrMsg, "TRAP:") {
        t.Fatalf("trapFunction expected UnsupportedOperationException with TRAP: message, got %+v", blk)
    }
    // trapDeprecated
    if blk, ok := trapDeprecated(nil).(*GErrBlk); !ok || blk.ExceptionType != excNames.UnsupportedOperationException || !strings.Contains(blk.ErrMsg, "TRAP:") {
        t.Fatalf("trapDeprecated expected UnsupportedOperationException with TRAP: message, got %+v", blk)
    }
    // trapUndocumented
    if blk, ok := trapUndocumented(nil).(*GErrBlk); !ok || blk.ExceptionType != excNames.UnsupportedOperationException || !strings.Contains(blk.ErrMsg, "TRAP:") {
        t.Fatalf("trapUndocumented expected UnsupportedOperationException with TRAP: message, got %+v", blk)
    }
    // trapProtected
    if blk, ok := trapProtected(nil).(*GErrBlk); !ok || blk.ExceptionType != excNames.UnsupportedOperationException || !strings.Contains(blk.ErrMsg, "TRAP:") {
        t.Fatalf("trapProtected expected UnsupportedOperationException with TRAP: message, got %+v", blk)
    }
}
