/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"container/list"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/frames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/trace"
	"testing"
)

// setupHelper initializes globals and base class definitions for tests
func setupConcatTest(t *testing.T) {
	t.Helper()
	globals.InitGlobals("test")
	trace.Init()
	classloader.Init()
	classloader.LoadBaseClasses()
}

// helper to build a minimal valid CPool and Frame for ConcatStringElements
func buildValidConcatFrame(callerClass, invdynType, recipe string, opStackArgs []any) (*frames.Frame, *list.List) {
	cp := &classloader.CPool{}
	// CpIndex layout:
	// 0: Dummy
	// 1: InvokeDynamic (slot 0)
	// 2: NameAndType (slot 0)
	// 3: NAT type UTF8 reference (slot 0 in Utf8Refs)
	// 4: Recipe UTF8 (slot 1 in Utf8Refs)
	cp.CpIndex = make([]classloader.CpEntry, 5)
	cp.CpIndex[0] = classloader.CpEntry{Type: classloader.Dummy, Slot: 0}
	cp.CpIndex[1] = classloader.CpEntry{Type: classloader.InvokeDynamic, Slot: 0}
	cp.CpIndex[2] = classloader.CpEntry{Type: classloader.NameAndType, Slot: 0}
	cp.CpIndex[3] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
	cp.CpIndex[4] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1}

	cp.InvokeDynamics = []classloader.InvokeDynamicEntry{
		{BootstrapIndex: 0, NameAndType: 2},
	}
	cp.NameAndTypes = []classloader.NameAndTypeEntry{
		{NameIndex: 0, DescIndex: 3},
	}
	cp.Utf8Refs = []string{invdynType, recipe}
	cp.Bootstraps = []classloader.BootstrapMethod{
		{
			MethodRef: 1,
			Args:      []uint16{4}, // points to CP index 4 (recipe)
		},
	}

	fr := frames.CreateFrame(10)
	fr.ClName = callerClass
	fr.CP = cp
	// INVOKEDYNAMIC opcode (0xba), followed by CP slot 1 (0x00, 0x01), and 2 trailing bytes
	fr.Meth = []byte{0xBA, 0x00, 0x01, 0x00, 0x00}
	fr.PC = 0

	// Push operands in stack order
	for i, arg := range opStackArgs {
		fr.OpStack[i] = arg
	}
	if len(opStackArgs) > 0 {
		fr.TOS = len(opStackArgs) - 1
	} else {
		fr.TOS = -1
	}

	fs := list.New()
	fs.PushFront(fr)
	return fr, fs
}

func TestLoad_Lang_Invoke_ConcatStringElements(t *testing.T) {
	setupConcatTest(t)
	Load_Lang_Invoke_ConcatStringElements()

	sig := "java/lang/invoke.ConcatStringElements()Ljava/lang/String;"
	gMeth, exists := ghelpers.MethodSignatures[sig]
	if !exists {
		t.Fatalf("Method signature %s not registered", sig)
	}
	if gMeth.ParamSlots != 0 {
		t.Errorf("Expected ParamSlots 0, got %d", gMeth.ParamSlots)
	}
	if !gMeth.NeedsContext {
		t.Errorf("Expected NeedsContext true, got false")
	}
	if gMeth.GFunction == nil {
		t.Errorf("Expected non-nil GFunction")
	}
}

func TestConcatStringElements_ValidationErrors_Table(t *testing.T) {
	setupConcatTest(t)

	callerClass := "caller/MyClass"
	classloader.MethAreaInsert(callerClass, &classloader.Klass{Status: 'N', Data: &classloader.ClData{Name: callerClass}})

	tests := []struct {
		name         string
		setupParams  func() []any
		expectedExcp int
	}{
		{
			name: "wrong_param_count_0",
			setupParams: func() []any {
				return []any{}
			},
			expectedExcp: excNames.MalformedParametersException,
		},
		{
			name: "wrong_param_count_2",
			setupParams: func() []any {
				return []any{"arg0", "arg1"}
			},
			expectedExcp: excNames.MalformedParametersException,
		},
		{
			name: "param_not_list",
			setupParams: func() []any {
				return []any{"not-a-list"}
			},
			expectedExcp: excNames.NullPointerException,
		},
		{
			name: "param_nil_list",
			setupParams: func() []any {
				return []any{(*list.List)(nil)}
			},
			expectedExcp: excNames.NullPointerException,
		},
		{
			name: "empty_frame_stack",
			setupParams: func() []any {
				return []any{list.New()}
			},
			expectedExcp: excNames.NullPointerException,
		},
		{
			name: "frame_front_value_nil",
			setupParams: func() []any {
				fs := list.New()
				fs.PushFront(nil)
				return []any{fs}
			},
			expectedExcp: excNames.NullPointerException,
		},
		{
			name: "frame_front_not_frame",
			setupParams: func() []any {
				fs := list.New()
				fs.PushFront("not-a-frame")
				return []any{fs}
			},
			expectedExcp: excNames.NullPointerException,
		},
		{
			name: "frame_front_nil_frame",
			setupParams: func() []any {
				fs := list.New()
				fs.PushFront((*frames.Frame)(nil))
				return []any{fs}
			},
			expectedExcp: excNames.NullPointerException,
		},
		{
			name: "frame_meth_nil",
			setupParams: func() []any {
				fr, fs := buildValidConcatFrame(callerClass, "()Ljava/lang/String;", "test", nil)
				fr.Meth = nil
				return []any{fs}
			},
			expectedExcp: excNames.ClassFormatError,
		},
		{
			name: "frame_pc_negative",
			setupParams: func() []any {
				fr, fs := buildValidConcatFrame(callerClass, "()Ljava/lang/String;", "test", nil)
				fr.PC = -1
				return []any{fs}
			},
			expectedExcp: excNames.ClassFormatError,
		},
		{
			name: "frame_pc_out_of_bounds",
			setupParams: func() []any {
				fr, fs := buildValidConcatFrame(callerClass, "()Ljava/lang/String;", "test", nil)
				fr.PC = len(fr.Meth) - 2 // fr.PC + 2 == len(fr.Meth), out of bounds
				return []any{fs}
			},
			expectedExcp: excNames.ClassFormatError,
		},
		{
			name: "frame_cp_not_cpool",
			setupParams: func() []any {
				fr, fs := buildValidConcatFrame(callerClass, "()Ljava/lang/String;", "test", nil)
				fr.CP = "invalid-cp"
				return []any{fs}
			},
			expectedExcp: excNames.ClassFormatError,
		},
		{
			name: "frame_cp_nil",
			setupParams: func() []any {
				fr, fs := buildValidConcatFrame(callerClass, "()Ljava/lang/String;", "test", nil)
				fr.CP = (*classloader.CPool)(nil)
				return []any{fs}
			},
			expectedExcp: excNames.ClassFormatError,
		},
		{
			name: "cp_id_entry_addrval_nil",
			setupParams: func() []any {
				fr, fs := buildValidConcatFrame(callerClass, "()Ljava/lang/String;", "test", nil)
				cp := fr.CP.(*classloader.CPool)
				// Overwrite CP slot 1 with a UTF8 entry so FetchCPentry returns string (AddrVal == nil)
				cp.CpIndex[1] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
				return []any{fs}
			},
			expectedExcp: excNames.ClassFormatError,
		},
		{
			name: "bsm_index_out_of_range",
			setupParams: func() []any {
				fr, fs := buildValidConcatFrame(callerClass, "()Ljava/lang/String;", "test", nil)
				cp := fr.CP.(*classloader.CPool)
				cp.InvokeDynamics[0].BootstrapIndex = 99
				return []any{fs}
			},
			expectedExcp: excNames.ClassFormatError,
		},
		{
			name: "bsm_no_args",
			setupParams: func() []any {
				fr, fs := buildValidConcatFrame(callerClass, "()Ljava/lang/String;", "test", nil)
				cp := fr.CP.(*classloader.CPool)
				cp.Bootstraps[0].Args = []uint16{}
				return []any{fs}
			},
			expectedExcp: excNames.ClassFormatError,
		},
		{
			name: "recipe_resolve_error_invalid_cp_index",
			setupParams: func() []any {
				fr, fs := buildValidConcatFrame(callerClass, "()Ljava/lang/String;", "test", nil)
				cp := fr.CP.(*classloader.CPool)
				// point recipe to out-of-range CP index 99
				cp.Bootstraps[0].Args = []uint16{99}
				return []any{fs}
			},
			expectedExcp: excNames.ClassFormatError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			params := tc.setupParams()
			res := ConcatStringElements(params)
			errBlk, ok := res.(*ghelpers.GErrBlk)
			if !ok {
				t.Fatalf("%s: expected *ghelpers.GErrBlk, got %T (%v)", tc.name, res, res)
			}
			if errBlk.ExceptionType != tc.expectedExcp {
				t.Errorf("%s: expected exception %d, got %d (msg: %s)", tc.name, tc.expectedExcp, errBlk.ExceptionType, errBlk.ErrMsg)
			}
		})
	}
}

func TestConcatStringElements_ExecutionPaths_Table(t *testing.T) {
	setupConcatTest(t)

	callerClass := "caller/MyClass"
	classloader.MethAreaInsert(callerClass, &classloader.Klass{Status: 'N', Data: &classloader.ClData{Name: callerClass}})

	tests := []struct {
		name         string
		invdynType   string
		recipe       string
		opStackArgs  []any
		expectedTOS  int
		expectedText string
	}{
		{
			name:         "literal_only_recipe",
			invdynType:   "()Ljava/lang/String;",
			recipe:       "Hello World",
			opStackArgs:  nil,
			expectedTOS:  -1,
			expectedText: "Hello World",
		},
		{
			name:         "empty_recipe",
			invdynType:   "()Ljava/lang/String;",
			recipe:       "",
			opStackArgs:  nil,
			expectedTOS:  -1,
			expectedText: "",
		},
		{
			name:         "dynamic_byte_argument_ascii",
			invdynType:   "(B)Ljava/lang/String;",
			recipe:       "Char: \x01",
			opStackArgs:  []any{int64(42)}, // ASCII 42 = '*'
			expectedTOS:  -1,
			expectedText: "Char: *",
		},
		{
			name:         "dynamic_byte_argument_letter",
			invdynType:   "(B)Ljava/lang/String;",
			recipe:       "ByteLetter: \x01!",
			opStackArgs:  []any{int64(65)}, // ASCII 65 = 'A'
			expectedTOS:  -1,
			expectedText: "ByteLetter: A!",
		},
		{
			name:         "dynamic_int_argument",
			invdynType:   "(I)Ljava/lang/String;",
			recipe:       "Count: \x01",
			opStackArgs:  []any{int64(12345)},
			expectedTOS:  -1,
			expectedText: "Count: 12345",
		},
		{
			name:         "dynamic_string_reference_argument",
			invdynType:   "(Ljava/lang/String;)Ljava/lang/String;",
			recipe:       "Hello \x01!",
			opStackArgs:  []any{object.StringObjectFromGoString("Jacobin")},
			expectedTOS:  -1,
			expectedText: "Hello Jacobin!",
		},
		{
			name:         "multiple_dynamic_arguments_popping_order",
			invdynType:   "(Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;",
			recipe:       "[\x01, \x01]",
			opStackArgs:  []any{object.StringObjectFromGoString("First"), object.StringObjectFromGoString("Second")},
			expectedTOS:  -1,
			expectedText: "[First, Second]",
		},
		{
			name:         "multiple_mixed_dynamic_arguments",
			invdynType:   "(Ljava/lang/String;IB)Ljava/lang/String;",
			recipe:       "\x01 has ID \x01 and grade \x01",
			opStackArgs:  []any{object.StringObjectFromGoString("Alice"), int64(1001), int64(65)},
			expectedTOS:  -1,
			expectedText: "Alice has ID 1001 and grade A",
		},
		{
			name:         "surplus_items_on_opstack_remain_intact",
			invdynType:   "(Ljava/lang/String;)Ljava/lang/String;",
			recipe:       "Result: \x01",
			opStackArgs:  []any{"base-item", object.StringObjectFromGoString("popped-item")},
			expectedTOS:  0, // "base-item" remains at index 0, TOS = 0
			expectedText: "Result: popped-item",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fr, fs := buildValidConcatFrame(callerClass, tc.invdynType, tc.recipe, tc.opStackArgs)
			res := ConcatStringElements([]any{fs})
			obj, ok := res.(*object.Object)
			if !ok || obj == nil {
				t.Fatalf("%s: expected *object.Object, got %T (%v)", tc.name, res, res)
			}
			gotStr := object.GoStringFromStringObject(obj)
			if gotStr != tc.expectedText {
				t.Errorf("%s: expected %q, got %q", tc.name, tc.expectedText, gotStr)
			}
			if fr.TOS != tc.expectedTOS {
				t.Errorf("%s: expected fr.TOS %d, got %d", tc.name, tc.expectedTOS, fr.TOS)
			}
		})
	}
}
