/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package jvm

import (
	"jacobin/src/classloader"
	"jacobin/src/frames"
	"jacobin/src/gfunction"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"jacobin/src/types"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Arrays are preloaded, so this should only confirm the presence of the class
// in the method area--and make sure it has no fields.
func TestInstantiateArray(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()

	anything, err := InstantiateClass(types.JavaByteArray, nil)
	if err != nil {
		t.Errorf("Got unexpected error from instantiating array: %s", err.Error())
	}
	obj := anything.(*object.Object)
	if len(obj.FieldTable) != 0 {
		t.Errorf("Expected 0 fields in array class, got %d fields", len(obj.FieldTable))
	}
}

func TestInstantiateString1(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()

	// initialize the MTable and other class entries
	classloader.MTable = make(map[string]classloader.MTentry)

	// Init classloader and load base classes
	err := classloader.Init() // must precede classloader.LoadBaseClasses
	if err != nil {
		t.Errorf("Got unexpected error from classloader.Init: %s", err.Error())
	}
	classloader.LoadBaseClasses()

	myobj, err := InstantiateClass(types.StringClassName, nil)
	if err != nil {
		t.Errorf("Got unexpected error from instantiating string: %s", err.Error())
	}

	obj := myobj.(*object.Object)
	klassType := stringPool.GetStringPointer(obj.KlassName)
	if *klassType != types.StringClassName {
		t.Errorf("Expected 'java/lang/String', got %s", *klassType)
	}

	if len(obj.FieldTable) < 2 {
		t.Errorf("Expected more than 1 field in String object, got %d fields", len(obj.FieldTable))
	}
}

func TestInstantiateNonExistentClass(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// redirect stderr, to avoid all the error msgs for a non-existent class
	normalStderr := os.Stderr
	_, werr, err := os.Pipe()

	os.Stderr = werr

	classloader.InitMethodArea()

	// initialize the MTable and other class entries
	classloader.MTable = make(map[string]classloader.MTentry)

	// Init classloader and load base classes
	err = classloader.Init() // must precede classloader.LoadBaseClasses
	if err != nil {
		t.Errorf("Got unexpected error from classloader.Init: %s", err.Error())
	}
	classloader.LoadBaseClasses()
	gfunction.MTableLoadGFunctions(&classloader.MTable)
	statics.PreloadStatics()

	myobj, err := InstantiateClass("$nosuchclass", nil)

	// restore stderr
	_ = werr.Close()
	os.Stderr = normalStderr

	if err == nil {
		t.Errorf("Expected error message for nonexistent class, but got none")
	}

	if myobj != nil {
		t.Errorf("Expected nil object, got %v", myobj)
	}
}

func TestLoadValidClass(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// redirect stderr, to avoid all the error msgs for a non-existent class
	normalStderr := os.Stderr
	_, werr, err := os.Pipe()
	if err != nil {
		t.Error("cannot create pipe for stderr")
	}
	os.Stderr = werr

	classloader.InitMethodArea()

	// initialize the MTable and other class entries
	classloader.MTable = make(map[string]classloader.MTentry)

	// Init classloader and load base classes
	err = classloader.Init() // must precede classloader.LoadBaseClasses
	if err != nil {
		t.Errorf("Got unexpected error from classloader.Init: %s", err.Error())
	}
	classloader.LoadBaseClasses()

	// we'll check that the class is loaded, then delete it, then load it and check again

	class := classloader.MethAreaFetch("java/lang/Integer")
	if class == nil {
		t.Errorf("Expected java.lang.Integer to be loaded in method area, but it wasn't")
	}

	classloader.MethAreaDelete("java/lang/Integer")
	class = classloader.MethAreaFetch("java/lang/Integer")
	if class != nil {
		t.Errorf("Expected java.lang.Integer to be absent from method area, but it wasn't")
	}

	// now load the class
	err = loadThisClass("java/lang/Integer")
	if err != nil {
		t.Errorf("Got unexpected error from loadThisClass(\"java/lang/Integer\"): %s", err.Error())
	}
	class = classloader.MethAreaFetch("java/lang/Integer")
	if class == nil {
		t.Errorf("Expected java.lang.Integer to be loaded in method area, but it wasn't")
	}

	// restore stderr
	_ = werr.Close()
	os.Stderr = normalStderr
}

// createField() should correctly set the default value and type for each of
// the primitive/reference field-descriptor kinds it supports.
func TestCreateFieldPrimitiveAndRefTypes(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()
	statics.PreloadStatics()

	tests := []struct {
		desc     string
		wantType string
		wantVal  any
	}{
		{types.Ref, types.Ref, nil},
		{types.Array, types.Array, nil},
		{types.Byte, types.Byte, int8(0)},
		{types.Char, types.Char, int64(0)},
		{types.Int, types.Int, int64(0)},
		{types.Long, types.Long, int64(0)},
		{types.Short, types.Short, int64(0)},
		{types.Bool, types.Bool, int64(0)},
		{types.Double, types.Double, 0.0},
		{types.Float, types.Float, 0.0},
	}

	for _, tt := range tests {
		k := &classloader.Klass{Data: &classloader.ClData{}}
		k.Data.CP.Utf8Refs = []string{tt.desc}
		f := classloader.Field{Desc: 0}

		fld, err := createField(f, k, "test/Class")
		if err != nil {
			t.Errorf("createField(%s) returned unexpected error: %s", tt.desc, err.Error())
			continue
		}
		if fld.Ftype != tt.wantType {
			t.Errorf("createField(%s): expected Ftype %s, got %s", tt.desc, tt.wantType, fld.Ftype)
		}
		if fld.Fvalue != tt.wantVal {
			t.Errorf("createField(%s): expected Fvalue %v, got %v", tt.desc, tt.wantVal, fld.Fvalue)
		}
	}
}

// createField() should return an error for a field descriptor it doesn't recognize.
func TestCreateFieldInvalidType(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()
	statics.PreloadStatics()

	k := &classloader.Klass{Data: &classloader.ClData{}}
	k.Data.CP.Utf8Refs = []string{"Q"} // not a valid field descriptor
	f := classloader.Field{Desc: 0}

	fld, err := createField(f, k, "test/Class")
	if err == nil {
		t.Errorf("Expected error for invalid field type, but got none")
	}
	if fld != nil {
		t.Errorf("Expected nil field on error, got %v", fld)
	}
}

// createField() should mark a static field's Ftype with the types.Static prefix
// and register the field in the Statics table.
func TestCreateFieldStaticField(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()
	statics.PreloadStatics()

	k := &classloader.Klass{Data: &classloader.ClData{}}
	k.Data.CP.Utf8Refs = []string{types.Int, "myStaticField"}
	f := classloader.Field{Desc: 0, Name: 1, IsStatic: true}

	classname := "test/StaticFieldClass"
	fld, err := createField(f, k, classname)
	if err != nil {
		t.Fatalf("createField() returned unexpected error: %s", err.Error())
	}

	if fld.Ftype != types.Static+types.Int {
		t.Errorf("Expected Ftype %s, got %s", types.Static+types.Int, fld.Ftype)
	}

	statVal, present := statics.QueryStatic(classname, "myStaticField")
	if !present {
		t.Errorf("Expected static field %s.myStaticField to be present in Statics table", classname)
	}
	if statVal.Type != types.Int {
		t.Errorf("Expected static field type %s, got %s", types.Int, statVal.Type)
	}
}

// doStaticDefaults() should populate the Statics table with default values
// for all the static fields of the class, and should skip instance fields.
func TestDoStaticDefaults(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()
	statics.PreloadStatics()

	classname := "test/DoStaticDefaultsClass"
	k := &classloader.Klass{Data: &classloader.ClData{}}
	k.Data.CP.Utf8Refs = []string{types.Int, "staticIntField", types.Ref, "instanceField"}
	k.Data.Fields = []classloader.Field{
		{Desc: 0, Name: 1, IsStatic: true},
		{Desc: 2, Name: 3, IsStatic: false},
	}

	doStaticDefaults(k, classname)

	statVal, present := statics.QueryStatic(classname, "staticIntField")
	if !present {
		t.Errorf("Expected static field %s.staticIntField to be present in Statics table", classname)
	} else if statVal.Value != int64(0) {
		t.Errorf("Expected default value 0, got %v", statVal.Value)
	}

	_, present = statics.QueryStatic(classname, "instanceField")
	if present {
		t.Errorf("Instance field %s.instanceField should not have been added to Statics table", classname)
	}
}

// superclassChain() should return the chain of superclasses (nearest first),
// stopping before java/lang/Object, for a class with a non-Object superclass.
func TestSuperclassChain(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()

	classloader.MTable = make(map[string]classloader.MTentry)
	err := classloader.Init()
	if err != nil {
		t.Fatalf("Got unexpected error from classloader.Init: %s", err.Error())
	}
	classloader.LoadBaseClasses()

	err = loadThisClass("java/lang/Exception")
	if err != nil {
		t.Fatalf("Got unexpected error from loadThisClass(java/lang/Exception): %s", err.Error())
	}
	k := classloader.MethAreaFetch("java/lang/Exception")
	if k == nil {
		t.Fatalf("Expected java/lang/Exception to be loaded, but it wasn't")
	}

	superclasses, err := superclassChain(k, "java/lang/Exception")
	if err != nil {
		t.Fatalf("superclassChain() returned unexpected error: %s", err.Error())
	}

	found := false
	for _, sc := range superclasses {
		if sc == "java/lang/Throwable" {
			found = true
		}
		if sc == types.ObjectClassName {
			t.Errorf("Expected superclassChain() to stop before %s, but it was included", types.ObjectClassName)
		}
	}
	if !found {
		t.Errorf("Expected java/lang/Throwable in superclass chain of java/lang/Exception, got %v", superclasses)
	}
}

// superclassChain() for java/lang/Object itself should return an empty chain,
// since java/lang/Object has no superclass.
func TestSuperclassChainForObject(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()

	classloader.MTable = make(map[string]classloader.MTentry)
	err := classloader.Init()
	if err != nil {
		t.Fatalf("Got unexpected error from classloader.Init: %s", err.Error())
	}
	classloader.LoadBaseClasses()

	err = loadThisClass(types.ObjectClassName)
	if err != nil {
		t.Fatalf("Got unexpected error from loadThisClass(%s): %s", types.ObjectClassName, err.Error())
	}
	k := classloader.MethAreaFetch(types.ObjectClassName)
	if k == nil {
		t.Fatalf("Expected %s to be loaded, but it wasn't", types.ObjectClassName)
	}

	superclasses, err := superclassChain(k, types.ObjectClassName)
	if err != nil {
		t.Fatalf("superclassChain() returned unexpected error: %s", err.Error())
	}
	if len(superclasses) != 0 {
		t.Errorf("Expected empty superclass chain for %s, got %v", types.ObjectClassName, superclasses)
	}
}

// InitializeClass() should be idempotent: calling it a second time on the
// same class should succeed via the fast path, without error.
func TestInitializeClassIdempotent(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()

	classloader.MTable = make(map[string]classloader.MTentry)
	err := classloader.Init()
	if err != nil {
		t.Fatalf("Got unexpected error from classloader.Init: %s", err.Error())
	}
	classloader.LoadBaseClasses()
	gfunction.MTableLoadGFunctions(&classloader.MTable)
	statics.PreloadStatics()

	fs := frames.CreateFrameStack()
	fs.PushFront(frames.CreateFrame(0))

	err = InitializeClass("java/lang/Integer", fs)
	if err != nil {
		t.Fatalf("First InitializeClass() call returned unexpected error: %s", err.Error())
	}

	err = InitializeClass("java/lang/Integer", fs)
	if err != nil {
		t.Errorf("Second InitializeClass() call (fast path) returned unexpected error: %s", err.Error())
	}
}

// If the method area is reset (as classloader.Init() does between unit tests)
// after a class was initialized, InitializeClass() should detect the stale
// cache entry and successfully reload/reinitialize the class rather than
// skipping initialization and later failing.
func TestInitializeClassStaleCacheIsReloaded(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()

	classloader.MTable = make(map[string]classloader.MTentry)
	err := classloader.Init()
	if err != nil {
		t.Fatalf("Got unexpected error from classloader.Init: %s", err.Error())
	}
	classloader.LoadBaseClasses()
	gfunction.MTableLoadGFunctions(&classloader.MTable)
	statics.PreloadStatics()

	fs := frames.CreateFrameStack()
	fs.PushFront(frames.CreateFrame(0))

	err = InitializeClass("java/lang/Integer", fs)
	if err != nil {
		t.Fatalf("First InitializeClass() call returned unexpected error: %s", err.Error())
	}

	// simulate the method area being reset, as happens between unit tests
	classloader.InitMethodArea()
	err = classloader.Init()
	if err != nil {
		t.Fatalf("Got unexpected error from second classloader.Init: %s", err.Error())
	}
	classloader.LoadBaseClasses()

	fs2 := frames.CreateFrameStack()
	fs2.PushFront(frames.CreateFrame(0))
	err = InitializeClass("java/lang/Integer", fs2)
	if err != nil {
		t.Errorf("InitializeClass() after stale cache reset returned unexpected error: %s", err.Error())
	}
	if classloader.MethAreaFetch("java/lang/Integer") == nil {
		t.Errorf("Expected java/lang/Integer to be reloaded into the method area after stale cache reset")
	}
}

// This should always work. java/lang/Object contains no instance or static fields,
// so this is about as simple a class instantiation as possible
func TestLoadClassJavaLangObject(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// redirect stderr, to avoid all the error msgs for a non-existent class
	normalStderr := os.Stderr
	_, werr, err := os.Pipe()
	os.Stderr = werr

	classloader.InitMethodArea()

	// initialize the MTable and other class entries
	classloader.MTable = make(map[string]classloader.MTentry)

	// Init classloader and load base classes
	err = classloader.Init() // must precede classloader.LoadBaseClasses
	if err != nil {
		t.Errorf("Got unexpected error from classloader.Init: %s", err.Error())
	}
	classloader.LoadBaseClasses()

	err = loadThisClass(types.ObjectClassName)

	// this should always work. java/lang/Object contains no instance or static fields,
	// so this is about as simple a class instantiation as possible

	// restore stderr
	_ = werr.Close()
	os.Stderr = normalStderr

	if err != nil {
		t.Errorf("Got unexpected error from loadThisClass: %s", err.Error())
	}
}

func registerDummyKlass(t *testing.T, classname string) {
	t.Helper()
	classloader.MTable = make(map[string]classloader.MTentry)
	if err := classloader.Init(); err != nil {
		t.Fatalf("classloader.Init() failed: %s", err.Error())
	}
	classloader.LoadBaseClasses()
	classloader.MethAreaInsert(classname, &classloader.Klass{
		Data: &classloader.ClData{Name: classname},
	})
}

// doStaticDefaults() must not overwrite a final static's compile-time
// constant value (from the ConstantValue attribute) with a generic
// type-default. javac typically emits no <clinit> PUTSTATIC for these
// fields, so this is the only place the real value is ever installed.
func TestDoStaticDefaultsHonorsConstValue(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()
	classloader.InitMethodArea()
	statics.PreloadStatics()

	classname := "test/ConstValueClass"
	k := &classloader.Klass{Data: &classloader.ClData{}}
	k.Data.CP.Utf8Refs = []string{types.Int, "constIntField", types.Int, "plainIntField"}
	k.Data.Fields = []classloader.Field{
		{Desc: 0, Name: 1, IsStatic: true, ConstValue: int64(42)},
		{Desc: 2, Name: 3, IsStatic: true}, // no ConstValue: should get the type default
	}

	doStaticDefaults(k, classname)

	constVal, present := statics.QueryStatic(classname, "constIntField")
	if !present {
		t.Fatalf("Expected static field %s.constIntField to be present", classname)
	}
	if constVal.Value != int64(42) {
		t.Errorf("Expected ConstValue 42 to be preserved, got %v", constVal.Value)
	}

	plainVal, present := statics.QueryStatic(classname, "plainIntField")
	if !present {
		t.Fatalf("Expected static field %s.plainIntField to be present", classname)
	}
	if plainVal.Value != int64(0) {
		t.Errorf("Expected type-default 0 for field with no ConstValue, got %v", plainVal.Value)
	}
}

// ensureInitialized must run initFn exactly once even when many goroutines
// call it concurrently for the same class.
func TestEnsureInitializedRunsOnce(t *testing.T) {
	globals.InitGlobals("test")

	classname := "test/ConcurrentInitClass"
	registerDummyKlass(t, "test/ConcurrentInitClass")
	fs := frames.CreateFrameStack()

	var runCount int32
	const goroutines = 50

	var wg sync.WaitGroup
	errs := make([]error, goroutines)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			errs[idx] = ensureInitialized(classname, fs, func() error {
				atomic.AddInt32(&runCount, 1)
				time.Sleep(10 * time.Millisecond) // widen the race window
				return nil
			})
		}(i)
	}
	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Errorf("goroutine %d: unexpected error: %s", i, err.Error())
		}
	}
	if runCount != 1 {
		t.Errorf("Expected initFn to run exactly once, ran %d times", runCount)
	}
}

// ensureInitialized must let the initializing goroutine re-enter (e.g. from
// its own <clinit> constructing an instance of the class being initialized)
// without deadlocking, and must not run initFn a second time for that thread.
func TestEnsureInitializedReentrant(t *testing.T) {
	globals.InitGlobals("test")

	classname := "test/ReentrantInitClass"
	registerDummyKlass(t, "test/ConcurrentInitClass")
	fs := frames.CreateFrameStack()

	var outerRuns, innerRuns int32
	done := make(chan error, 1)
	go func() {
		done <- ensureInitialized(classname, fs, func() error {
			atomic.AddInt32(&outerRuns, 1)
			// Simulate <clinit> re-entering, e.g. constructing an instance
			// of the class it's still initializing.
			return ensureInitialized(classname, fs, func() error {
				atomic.AddInt32(&innerRuns, 1)
				return nil
			})
		})
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}
	case <-time.After(2 * time.Second):
		t.Fatal("ensureInitialized deadlocked on re-entrant call")
	}

	if outerRuns != 1 {
		t.Errorf("Expected outer initFn to run once, ran %d times", outerRuns)
	}
	if innerRuns != 0 {
		t.Errorf("Expected inner (re-entrant) initFn to be skipped, ran %d times", innerRuns)
	}
}

// A goroutine with a different frame stack must wait for the initializing
// goroutine to finish, then see it as done rather than re-running initFn.
func TestEnsureInitializedWaiterSeesCompletion(t *testing.T) {
	globals.InitGlobals("test")
	classloader.InitMethodArea()

	classname := "test/WaiterInitClass"
	registerDummyKlass(t, classname)
	fsOwner := frames.CreateFrameStack()
	fsWaiter := frames.CreateFrameStack()

	var runCount int32
	ownerStarted := make(chan struct{})
	release := make(chan struct{})

	ownerDone := make(chan error, 1)
	go func() {
		ownerDone <- ensureInitialized(classname, fsOwner, func() error {
			atomic.AddInt32(&runCount, 1)
			close(ownerStarted) // signals initInProgress is definitely held by us
			<-release
			return nil
		})
	}()

	<-ownerStarted // guarantees the owner has won initNone->initInProgress

	waiterDone := make(chan error, 1)
	go func() {
		waiterDone <- ensureInitialized(classname, fsWaiter, func() error {
			atomic.AddInt32(&runCount, 1) // should never run
			return nil
		})
	}()

	// Best-effort: give the waiter time to reach initCond.Wait() before we
	// let the owner finish. ensureInitialized has no hook to observe this
	// directly, so this is a generous sleep rather than a hard guarantee.
	time.Sleep(200 * time.Millisecond)
	close(release)

	select {
	case err := <-ownerDone:
		if err != nil {
			t.Errorf("owner: unexpected error: %s", err.Error())
		}
	case <-time.After(3 * time.Second):
		t.Fatal("owner call did not complete")
	}
	select {
	case err := <-waiterDone:
		if err != nil {
			t.Errorf("waiter: unexpected error: %s", err.Error())
		}
	case <-time.After(3 * time.Second):
		t.Fatal("waiter call did not complete — possibly missed the broadcast")
	}

	if runCount != 1 {
		t.Errorf("Expected initFn to run exactly once, ran %d times", runCount)
	}
}
