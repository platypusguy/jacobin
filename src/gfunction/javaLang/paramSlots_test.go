/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/util"
	"strings"
	"testing"
)

func loadAllJavaLang() {
	globals.InitGlobals("test")
	ClassClinitIsh()
	Load_Lang_Boolean()
	Load_Lang_Byte()
	Load_Lang_Character()
	Load_Lang_CharSequence()
	Load_Lang_Class()
	Load_Lang_Classloader()
	Load_Lang_Double()
	Load_Lang_Enum()
	Load_Lang_Float()
	Load_Lang_Integer()
	Load_Lang_Invoke_MethodHandle()
	Load_Lang_Invoke_MethodType()
	Load_Lang_Invoke_StringConcatFactory()
	Load_Lang_Long()
	Load_Lang_Math()
	Load_Lang_Object()
	Load_Lang_Process()
	Load_Lang_Process_Builder()
	Load_Lang_Process_Handle_Impl()
	Load_Lang_Record()
	Load_Lang_Reflect_Modifier()
	Load_Lang_Runtime()
	Load_Lang_Runtime_SwitchBootstraps()
	Load_Lang_SecurityManager()
	Load_Lang_Short()
	Load_Lang_StackTraceELement()
	Load_Lang_String()
	Load_Lang_StringBuffer()
	Load_Lang_StringBuilder()
	Load_Lang_UTF16()
	Load_Lang_System()
	Load_Lang_Thread()
	Load_Lang_Thread_Group()
	Load_Lang_Thread_State()
	Load_Lang_Throwable()
	Load_Lang_Void()

	// java/lang exception and error classes
	Load_Java_Lang_AssertionError()
	Load_Java_Lang_Exception()
	Load_Java_Lang_Invoke_StringConcatException()
	Load_Java_Lang_NullPointerException()
	Load_Java_Lang_RuntimeException()
	Load_Java_Lang_UnsupportedOperationException()
}

const classPrefix = "java/lang/"

func TestJavaLangParamSlotsConsistency(t *testing.T) {
	// Save and restore the global ghelpers.MethodSignatures map to avoid test pollution
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	loadAllJavaLang()

	var mismatches []string
	for k, v := range ghelpers.MethodSignatures {
		if !strings.HasPrefix(k, classPrefix) {
			continue
		}
		_, _, methType := util.ParseFQN(k)
		params := util.ParseIncomingParamsFromMethTypeString(methType)
		expectedSlots := len(params)
		if v.ParamSlots != expectedSlots {
			t.Errorf("Signature mismatch for %s: ParamSlots is %d, expected %d (params: %v)", k, v.ParamSlots, expectedSlots, params)
			mismatches = append(mismatches, k)
		}
	}
	t.Logf("Total %s signatures checked: %d, mismatches: %d", classPrefix, len(ghelpers.MethodSignatures), len(mismatches))
	if len(mismatches) > 0 {
		t.Fatalf("Found %d ParamSlots mismatches", len(mismatches))
	}
}
