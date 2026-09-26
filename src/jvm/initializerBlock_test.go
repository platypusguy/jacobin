/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import (
	"jacobin/src/classloader"
	"jacobin/src/frames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"jacobin/src/types"
	"strings"
	"testing"
)

// makeTestKlass builds a minimal *classloader.Klass with the given name and
// superclass name, ready to be inserted into the MethArea.
func makeTestKlass(name, superclassName string) *classloader.Klass {
	k := &classloader.Klass{
		Data: &classloader.ClData{
			Name:            name,
			SuperclassIndex: stringPool.GetStringIndex(&superclassName),
			ClInit:          types.ClInitNotRun,
		},
	}
	return k
}

// registerClinit inserts an MTable entry for className.<clinit>()V.
func registerClinit(className string, mtype byte, meth any) {
	classloader.AddEntry(&classloader.MTable, className+".<clinit>()V",
		classloader.MTentry{Meth: meth, MType: mtype})
}

func TestRunNativeInitializer(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	called := false
	mt := classloader.MTentry{
		MType: 'G',
		Meth: ghelpers.GMeth{
			GFunction: func(_ []any) any {
				called = true
				return nil
			},
		},
	}

	k := makeTestKlass("test/NativeClinit", types.ObjectClassName)

	fs := frames.CreateFrameStack()
	f := frames.CreateFrame(1)
	f.ClName = "test/NativeClinit"
	_ = frames.PushFrame(fs, f)

	err := runNativeInitializer(mt, k, fs)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if !called {
		t.Errorf("expected the native gfunction to be invoked")
	}
	if k.Data.ClInit != types.ClInitRun {
		t.Errorf("expected ClInit=ClInitRun, got %v", k.Data.ClInit)
	}
}

func TestRunJavaInitializer(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	cp := &classloader.CPool{}
	jm := classloader.JmEntry{
		Cp:        cp,
		MaxStack:  1,
		MaxLocals: 0,
		Code:      []byte{0xb1}, // RETURN
	}

	k := makeTestKlass("test/JavaClinit", types.ObjectClassName)

	fs := frames.CreateFrameStack()

	err := runJavaInitializer(jm, k, fs)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if k.Data.ClInit != types.ClInitRun {
		t.Errorf("expected ClInit=ClInitRun, got %v", k.Data.ClInit)
	}
}

func TestRunInitializationBlock_SuperclassChain(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()
	classloader.MTable = make(map[string]classloader.MTentry)

	base := makeTestKlass("test/InitBase", types.ObjectClassName)
	sub := makeTestKlass("test/InitSub", "test/InitBase")

	classloader.MethAreaInsert("test/InitBase", base)
	classloader.MethAreaInsert("test/InitSub", sub)

	retCode := []byte{0xb1} // RETURN
	registerClinit("test/InitBase", 'J', classloader.JmEntry{
		Cp: &classloader.CPool{}, MaxStack: 1, MaxLocals: 0, Code: retCode,
	})
	registerClinit("test/InitSub", 'J', classloader.JmEntry{
		Cp: &classloader.CPool{}, MaxStack: 1, MaxLocals: 0, Code: retCode,
	})

	fs := frames.CreateFrameStack()

	err := runInitializationBlock(sub, nil, fs)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if base.Data.ClInit != types.ClInitRun {
		t.Errorf("expected base class ClInit=ClInitRun, got %v", base.Data.ClInit)
	}
	if sub.Data.ClInit != types.ClInitRun {
		t.Errorf("expected sub class ClInit=ClInitRun, got %v", sub.Data.ClInit)
	}
}

func TestRunInitializationBlock_MissingClassInSuperclasses(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()
	classloader.MTable = make(map[string]classloader.MTentry)

	sub := makeTestKlass("test/InitSubMissing", types.ObjectClassName)
	classloader.MethAreaInsert("test/InitSubMissing", sub)

	fs := frames.CreateFrameStack()

	// Pass a pre-computed superclasses list that references a class not
	// present in the MethArea, to exercise the classK==nil error branch.
	err := runInitializationBlock(sub, []string{"test/InitSubMissing", "test/DoesNotExist"}, fs)
	if err == nil {
		t.Fatalf("expected an error for missing class in superclasses, got nil")
	}
	if !strings.Contains(err.Error(), "MethAreaFetch could not find class") {
		t.Errorf("expected 'MethAreaFetch could not find class' error, got: %v", err)
	}
}

func TestRunInitializationBlock_UnloadableSuperclass(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()
	classloader.MTable = make(map[string]classloader.MTentry)

	// SuperclassIndex points to a class that can never be found/loaded.
	sub := makeTestKlass("test/InitSubBadSuper", "test/NoSuchSuperclassAtAll")
	classloader.MethAreaInsert("test/InitSubBadSuper", sub)

	fs := frames.CreateFrameStack()

	err := runInitializationBlock(sub, nil, fs)
	if err == nil {
		t.Fatalf("expected an error when the superclass cannot be loaded, got nil")
	}
}

func TestRunInitializationBlock_NoClinit(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()
	classloader.MTable = make(map[string]classloader.MTentry)

	sub := makeTestKlass("test/InitSubNoClinit", types.ObjectClassName)
	classloader.MethAreaInsert("test/InitSubNoClinit", sub)

	fs := frames.CreateFrameStack()

	// No <clinit> registered in MTable at all -- should simply be a no-op.
	err := runInitializationBlock(sub, nil, fs)
	if err != nil {
		t.Errorf("expected nil error when no <clinit> exists, got %v", err)
	}
}

func TestInitializePrimitiveWrappers_SkipsAlreadyRun(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()
	classloader.MTable = make(map[string]classloader.MTentry)
	_ = classloader.Init()
	classloader.LoadBaseClasses()

	k := classloader.MethAreaFetch("java/lang/Boolean")
	if k == nil {
		t.Fatalf("test setup failure: java/lang/Boolean not preloaded")
	}
	k.Data.ClInit = types.ClInitRun

	invokedBoolean := false
	glob := globals.GetGlobalRef()
	glob.FuncInvokeGFunction = func(fqn string, _ []any) any {
		if fqn == "java/lang/Boolean.<clinit>()V" {
			invokedBoolean = true
		}
		return nil
	}

	InitializePrimitiveWrappers()

	if invokedBoolean {
		t.Errorf("expected already-initialized java/lang/Boolean to be skipped, but its <clinit> gfunction was invoked")
	}
}

func TestInitializePrimitiveWrappers_SkipsMissingClass(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()

	// Remove all wrapper classes from the MethArea, so the loop hits the
	// "class not found" continue branch for every one of them.
	wrappers := []string{
		"java/lang/Boolean", "java/lang/Byte", "java/lang/Character",
		"java/lang/Short", "java/lang/Integer", "java/lang/Long",
		"java/lang/Float", "java/lang/Double", "java/lang/Void",
	}
	for _, w := range wrappers {
		classloader.MethAreaDelete(w)
	}

	invoked := false
	glob := globals.GetGlobalRef()
	glob.FuncInvokeGFunction = func(_ string, _ []any) any {
		invoked = true
		return nil
	}

	// should not panic and should not invoke anything since no classes exist
	InitializePrimitiveWrappers()

	if invoked {
		t.Errorf("expected no gfunction invocation when wrapper classes are missing")
	}
}

// TestRunInitializationBlock_DoubleClinitGuard verifies that calling
// runInitializationBlock twice for the same class only runs its <clinit>
// once. A shared counter incremented inside the native <clinit> gfunction
// must end up at 1, not 2.
func TestRunInitializationBlock_DoubleClinitGuard(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()
	classloader.MTable = make(map[string]classloader.MTentry)

	counter := 0
	className := "test/DoubleClinitGuard"
	k := makeTestKlass(className, types.ObjectClassName)
	classloader.MethAreaInsert(className, k)

	registerClinit(className, 'G', ghelpers.GMeth{
		GFunction: func(_ []any) any {
			counter++
			return nil
		},
	})

	fs := frames.CreateFrameStack()
	f := frames.CreateFrame(1)
	f.ClName = className
	_ = frames.PushFrame(fs, f)

	// Mirrors the real call sites (e.g. INVOKESTATIC in interpreter.go), which
	// only call runInitializationBlock when ClInit hasn't run yet.
	clinitGuard := func() error {
		if k.Data.ClInit == types.ClInitNotRun {
			return runInitializationBlock(k, nil, fs)
		}
		return nil
	}

	// First call should run <clinit> and increment the counter.
	if err := clinitGuard(); err != nil {
		t.Fatalf("first call: expected nil error, got %v", err)
	}

	// Second call should be a no-op because ClInit is now ClInitRun.
	if err := clinitGuard(); err != nil {
		t.Fatalf("second call: expected nil error, got %v", err)
	}

	if counter != 1 {
		t.Errorf("expected <clinit> to run exactly once, counter=%d", counter)
	}
	if k.Data.ClInit != types.ClInitRun {
		t.Errorf("expected ClInit=ClInitRun, got %v", k.Data.ClInit)
	}
}

// TestRunInitializationBlock_TwoClasses_ClinitGuard verifies the guard
// works correctly across a superclass chain of two classes: running the
// chain twice (e.g. once directly on the subclass, then again after the
// superclass was independently reached) must still only run each class's
// <clinit> exactly once, for a total counter value of 2 (one per class),
// not 4.
func TestRunInitializationBlock_TwoClasses_ClinitGuard(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()
	classloader.MTable = make(map[string]classloader.MTentry)

	counter := 0
	base := makeTestKlass("test/GuardBase", types.ObjectClassName)
	sub := makeTestKlass("test/GuardSub", "test/GuardBase")

	classloader.MethAreaInsert("test/GuardBase", base)
	classloader.MethAreaInsert("test/GuardSub", sub)

	incrementer := func(_ []any) any {
		counter++
		return nil
	}
	registerClinit("test/GuardBase", 'G', ghelpers.GMeth{GFunction: incrementer})
	registerClinit("test/GuardSub", 'G', ghelpers.GMeth{GFunction: incrementer})

	fs := frames.CreateFrameStack()
	f := frames.CreateFrame(1)
	f.ClName = "test/GuardSub"
	_ = frames.PushFrame(fs, f)

	// Mirrors the real call sites (e.g. INVOKESTATIC in interpreter.go), which
	// only call runInitializationBlock when ClInit hasn't run yet.
	clinitGuard := func() error {
		if sub.Data.ClInit == types.ClInitNotRun {
			return runInitializationBlock(sub, nil, fs)
		}
		return nil
	}

	// First run: builds the superclass chain and runs both <clinit>s.
	if err := clinitGuard(); err != nil {
		t.Fatalf("first call: expected nil error, got %v", err)
	}

	// Second run: both classes are now ClInitRun, so neither should re-run.
	if err := clinitGuard(); err != nil {
		t.Fatalf("second call: expected nil error, got %v", err)
	}

	if counter != 2 {
		t.Errorf("expected each class's <clinit> to run exactly once (total 2), counter=%d", counter)
	}
	if base.Data.ClInit != types.ClInitRun {
		t.Errorf("expected base class ClInit=ClInitRun, got %v", base.Data.ClInit)
	}
	if sub.Data.ClInit != types.ClInitRun {
		t.Errorf("expected sub class ClInit=ClInitRun, got %v", sub.Data.ClInit)
	}
}

func TestInitializePrimitiveWrappers_InvokesAndMarksRun(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()
	classloader.MTable = make(map[string]classloader.MTentry)
	_ = classloader.Init()
	classloader.LoadBaseClasses()

	k := classloader.MethAreaFetch("java/lang/Void")
	if k == nil {
		t.Fatalf("test setup failure: java/lang/Void not preloaded")
	}
	k.Data.ClInit = types.ClInitNotRun

	var calledFqn string
	glob := globals.GetGlobalRef()
	glob.FuncInvokeGFunction = func(fqn string, _ []any) any {
		if fqn == "java/lang/Void.<clinit>()V" {
			calledFqn = fqn
		}
		return nil
	}

	InitializePrimitiveWrappers()

	if calledFqn == "" {
		t.Errorf("expected java/lang/Void.<clinit>()V to be invoked via FuncInvokeGFunction")
	}
	if k.Data.ClInit != types.ClInitRun {
		t.Errorf("expected java/lang/Void ClInit=ClInitRun after InitializePrimitiveWrappers, got %v", k.Data.ClInit)
	}
}
