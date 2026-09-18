/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package exceptions

import (
	"io"
	"jacobin/src/classloader"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"os"
	"testing"
)

// buildCPwithClassRef builds a minimal *classloader.CPool whose CP index 1
// is a ClassRef entry pointing to className, suitable for use as the
// CatchType of a CodeException entry.
func buildCPwithClassRef(className string) *classloader.CPool {
	cp := &classloader.CPool{}
	nameCopy := className
	idx := stringPool.GetStringIndex(&nameCopy)
	cp.ClassRefs = []uint32{idx}
	cp.CpIndex = []classloader.CpEntry{
		{Type: 0, Slot: 0}, // index 0 is unused/dummy
		{Type: classloader.ClassRef, Slot: 0},
	}
	return cp
}

// registerJavaMethod inserts a Java ('J') MTable entry for fullMethName with
// the given exception table and constant pool.
func registerJavaMethod(fullMethName string, cp *classloader.CPool, excs []classloader.CodeException) {
	jm := classloader.JmEntry{
		Exceptions: excs,
		Cp:         cp,
	}
	classloader.AddEntry(&classloader.MTable, fullMethName, classloader.MTentry{Meth: jm, MType: 'J'})
}

// captureStderr runs fn while redirecting os.Stderr and returns what was written.
func captureStderr(fn func()) string {
	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	fn()
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stderr = normalStderr
	return string(out)
}

func makeFrame(clName, methName, methType string, cp *classloader.CPool) *frames.Frame {
	f := frames.CreateFrame(2)
	f.ClName = clName
	f.MethName = methName
	f.MethType = methType
	f.CP = cp
	return f
}

// --- locateExceptionFrame ---

func TestLocateExceptionFrame_MethodNotFound(t *testing.T) {
	globals.InitGlobals("test")
	f := makeFrame("com/test/NoSuchClass", "noSuchMethod", "()V", nil)

	out := captureStderr(func() {
		fr, pc := locateExceptionFrame(f, "java/lang/Exception", 0)
		if fr != nil || pc != -1 {
			t.Errorf("expected nil frame and pc=-1, got frame=%v pc=%d", fr, pc)
		}
	})
	if out == "" {
		t.Errorf("expected an error message to be logged for missing method entry")
	}
}

func TestLocateExceptionFrame_GMethodHasNoHandler(t *testing.T) {
	globals.InitGlobals("test")
	fqn := "com/test/GClass.gMethod()V"
	classloader.AddEntry(&classloader.MTable, fqn, classloader.MTentry{Meth: "dummy-g-method", MType: 'G'})

	f := makeFrame("com/test/GClass", "gMethod", "()V", nil)
	fr, pc := locateExceptionFrame(f, "java/lang/Exception", 0)
	if fr != nil || pc != -1 {
		t.Errorf("expected nil frame and pc=-1 for G-method, got frame=%v pc=%d", fr, pc)
	}
}

func TestLocateExceptionFrame_NoExceptionTable(t *testing.T) {
	globals.InitGlobals("test")
	fqn := "com/test/NoExcTable.method()V"
	registerJavaMethod(fqn, nil, nil)

	f := makeFrame("com/test/NoExcTable", "method", "()V", nil)
	fr, pc := locateExceptionFrame(f, "java/lang/Exception", 0)
	if fr != nil || pc != -1 {
		t.Errorf("expected nil frame and pc=-1 when no exception table, got frame=%v pc=%d", fr, pc)
	}
}

func TestLocateExceptionFrame_PcOutsideRange(t *testing.T) {
	globals.InitGlobals("test")
	cp := buildCPwithClassRef("java/lang/Exception")
	excs := []classloader.CodeException{
		{StartPc: 5, EndPc: 10, HandlerPc: 20, CatchType: 1},
	}
	fqn := "com/test/OutOfRange.method()V"
	registerJavaMethod(fqn, cp, excs)

	f := makeFrame("com/test/OutOfRange", "method", "()V", cp)
	fr, pc := locateExceptionFrame(f, "java/lang/Exception", 1) // pc=1 not in [5,10)
	if fr != nil || pc != -1 {
		t.Errorf("expected nil frame and pc=-1 when pc is outside handler range, got frame=%v pc=%d", fr, pc)
	}
}

func TestLocateExceptionFrame_DirectMatch(t *testing.T) {
	globals.InitGlobals("test")
	cp := buildCPwithClassRef("java/lang/ArithmeticException")
	excs := []classloader.CodeException{
		{StartPc: 0, EndPc: 10, HandlerPc: 42, CatchType: 1},
	}
	fqn := "com/test/DirectMatch.method()V"
	registerJavaMethod(fqn, cp, excs)

	f := makeFrame("com/test/DirectMatch", "method", "()V", cp)
	fr, pc := locateExceptionFrame(f, "java/lang/ArithmeticException", 3)
	if fr != f {
		t.Errorf("expected returned frame to be the input frame, got %v", fr)
	}
	if pc != 42 {
		t.Errorf("expected handler pc 42, got %d", pc)
	}
}

func TestLocateExceptionFrame_MatchOnThrowableLiteral(t *testing.T) {
	globals.InitGlobals("test")
	cp := buildCPwithClassRef("java/lang/Throwable")
	excs := []classloader.CodeException{
		{StartPc: 0, EndPc: 10, HandlerPc: 7, CatchType: 1},
	}
	fqn := "com/test/ThrowableLiteral.method()V"
	registerJavaMethod(fqn, cp, excs)

	f := makeFrame("com/test/ThrowableLiteral", "method", "()V", cp)
	fr, pc := locateExceptionFrame(f, "com/test/SomeCustomException", 2)
	if fr != f || pc != 7 {
		t.Errorf("expected match via literal 'java/lang/Throwable' catch type, got frame=%v pc=%d", fr, pc)
	}
}

func TestLocateExceptionFrame_MatchViaSuperclassInMethArea(t *testing.T) {
	globals.InitGlobals("test")
	if classloader.MethArea == nil {
		classloader.InitMethodArea()
	}

	// register a "custom/MyException" class whose superclass is java/lang/Exception
	customExcName := "custom/MyException"
	cp := buildCPwithClassRef(customExcName)
	excs := []classloader.CodeException{
		{StartPc: 0, EndPc: 10, HandlerPc: 99, CatchType: 1},
	}
	fqn := "com/test/SuperclassMatch.method()V"
	registerJavaMethod(fqn, cp, excs)

	superName := "java/lang/Exception"
	klass := &classloader.Klass{
		Data: &classloader.ClData{
			Name:            customExcName,
			SuperclassIndex: stringPool.GetStringIndex(&superName),
		},
	}
	classloader.MethAreaInsert(customExcName, klass)

	f := makeFrame("com/test/SuperclassMatch", "method", "()V", cp)
	fr, pc := locateExceptionFrame(f, "com/test/DifferentException", 3)
	if fr != f || pc != 99 {
		t.Errorf("expected match via superclass lookup in method area, got frame=%v pc=%d", fr, pc)
	}
}

func TestLocateExceptionFrame_CatchClassNotInMethAreaIsSkipped(t *testing.T) {
	globals.InitGlobals("test")
	if classloader.MethArea == nil {
		classloader.InitMethodArea()
	}

	unknownExcName := "custom/UnknownException"
	classloader.MethAreaDelete(unknownExcName) // ensure it's not present
	cp := buildCPwithClassRef(unknownExcName)
	excs := []classloader.CodeException{
		{StartPc: 0, EndPc: 10, HandlerPc: 55, CatchType: 1},
	}
	fqn := "com/test/UnknownCatchClass.method()V"
	registerJavaMethod(fqn, cp, excs)

	f := makeFrame("com/test/UnknownCatchClass", "method", "()V", cp)
	fr, pc := locateExceptionFrame(f, "com/test/DifferentException", 3)
	if fr != nil || pc != -1 {
		t.Errorf("expected no handler found when catch class is absent from method area, got frame=%v pc=%d", fr, pc)
	}
}

// --- FindCatchFrame ---

func TestFindCatchFrame_FoundInFirstFrame(t *testing.T) {
	globals.InitGlobals("test")
	cp := buildCPwithClassRef("java/lang/ArithmeticException")
	excs := []classloader.CodeException{
		{StartPc: 0, EndPc: 10, HandlerPc: 42, CatchType: 1},
	}
	fqn := "com/test/FirstFrame.method()V"
	registerJavaMethod(fqn, cp, excs)

	f := makeFrame("com/test/FirstFrame", "method", "()V", cp)
	f.PC = 3
	f.ExceptionPC = -1

	fs := frames.CreateFrameStack()
	_ = frames.PushFrame(fs, f)

	fr, pc := FindCatchFrame(fs, "java.lang.ArithmeticException", 3)
	if fr != f || pc != 42 {
		t.Errorf("expected to find handler in first frame, got frame=%v pc=%d", fr, pc)
	}
}

func TestFindCatchFrame_UnlocksAndMovesToNextFrame(t *testing.T) {
	globals.InitGlobals("test")

	// first frame: no exception table, but has a synchronized object to unlock
	fqn1 := "com/test/Outer.method()V"
	registerJavaMethod(fqn1, nil, nil)
	f1 := makeFrame("com/test/Outer", "method", "()V", nil)
	f1.PC = 5
	f1.ExceptionPC = -1
	obj := object.MakeEmptyObject()
	if err := obj.ObjLock(int32(f1.Thread)); err != nil {
		t.Fatalf("failed to lock test object: %v", err)
	}
	f1.ObjSync = obj

	// second frame: has a matching handler
	cp2 := buildCPwithClassRef("java/lang/Exception")
	excs2 := []classloader.CodeException{
		{StartPc: 0, EndPc: 10, HandlerPc: 17, CatchType: 1},
	}
	fqn2 := "com/test/Inner.method()V"
	registerJavaMethod(fqn2, cp2, excs2)
	f2 := makeFrame("com/test/Inner", "method", "()V", cp2)
	f2.PC = 4
	f2.ExceptionPC = -1

	fs := frames.CreateFrameStack()
	_ = frames.PushFrame(fs, f1)
	// f1 is now at front; append f2 after it
	fs.PushBack(f2)

	fr, pc := FindCatchFrame(fs, "java.lang.Exception", 5)
	if fr != f2 {
		t.Errorf("expected handler to be found in second frame, got %v", fr)
	}
	if pc != 17 {
		t.Errorf("expected handler pc 17, got %d", pc)
	}
}

func TestFindCatchFrame_NotFoundReturnsNil(t *testing.T) {
	globals.InitGlobals("test")

	fqn := "com/test/NoHandlerAnywhere.method()V"
	registerJavaMethod(fqn, nil, nil)
	f := makeFrame("com/test/NoHandlerAnywhere", "method", "()V", nil)
	f.PC = 2
	f.ExceptionPC = -1

	fs := frames.CreateFrameStack()
	_ = frames.PushFrame(fs, f)

	fr, pc := FindCatchFrame(fs, "java.lang.Exception", 2)
	if fr != nil || pc != -1 {
		t.Errorf("expected no handler found across the whole frame stack, got frame=%v pc=%d", fr, pc)
	}
}

func TestFindCatchFrame_EmptyStackReturnsNil(t *testing.T) {
	globals.InitGlobals("test")

	fs := frames.CreateFrameStack()
	// With an empty stack, the loop body never runs, so the zero-value
	// results (nil frame, pc=0) are returned rather than the "not found"
	// sentinel (nil, -1) used elsewhere in the function.
	fr, pc := FindCatchFrame(fs, "java.lang.Exception", 0)
	if fr != nil || pc != 0 {
		t.Errorf("expected nil frame and pc=0 for empty frame stack, got frame=%v pc=%d", fr, pc)
	}
}
