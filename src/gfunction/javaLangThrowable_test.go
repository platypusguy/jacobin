/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/globals"
	"jacobin/log"
	"jacobin/statics"
	"strings"
	"testing"
)

func TestJavaLangThrowableClinit(t *testing.T) {
	statics.Statics = make(map[string]statics.Static)

	throwableClinit(nil)
	_, ok := statics.Statics["Throwable.UNASSIGNED_STACK"]
	if !ok {
		t.Error("JavaLangThrowableClinit: Throwable.UNASSIGNED_STACK not found")
	}

	_, ok = statics.Statics["Throwable.SUPPRESSED_SENTINEL"]
	if !ok {
		t.Error("JavaLangThrowableClinit: Throwable.SUPPRESSED_SENTINEL not found")
	}

	_, ok = statics.Statics["Throwable.EMPTY_THROWABLE_ARRAY"]
	if !ok {
		t.Error("Throwable.EMPTY_THROWABLE_ARRAY not found")
	}
}

func TestJavaLangThrowableFillInStackTraceWrongParmCount(t *testing.T) {
	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.SEVERE)

	params := []interface{}{1}
	err := fillInStackTrace(params)

	var retVal error
	switch err.(type) {
	case error:
		retVal = err.(error)
	default:
		t.Error("JavaLangThrowableFillInStack should have returned an error, but did not")
	}

	errMsg := retVal.Error()
	if !strings.HasPrefix(errMsg, "fillInStackTrace() expected two parameters") {
		t.Errorf("did not get expected error message, got: %s", errMsg)
	}
}
