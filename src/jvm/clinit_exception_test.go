/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package jvm

import (
	"container/list"
	"jacobin/src/classloader"
	"jacobin/src/frames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/opcodes"
	"jacobin/src/stringPool"
	"sync"
	"testing"
)

func TestGetStaticClinitTrapExceptionCaught(t *testing.T) {
	globals.InitGlobals("testWithoutShutdown")
	globals.TraceInst = false
	classloader.InitMethodArea()
	initTable = sync.Map{}

	className := "test/TrapClass"
	fieldName := "UTF_8"

	// Register a method signature with ghelpers.TrapClass for <clinit>
	clinitSig := className + ".<clinit>()V"
	gm := ghelpers.GMeth{
		ParamSlots: 0,
		GFunction:  ghelpers.TrapClass,
	}
	ghelpers.MethodSignatures[clinitSig] = gm
	classloader.AddEntry(&classloader.MTable, clinitSig, classloader.MTentry{Meth: gm, MType: 'G'})

	// Register the class in MethodArea with a field
	objClass := "java/lang/Object"
	objIdx := stringPool.GetStringIndex(&objClass)
	classloader.MethAreaInsert(className, &classloader.Klass{
		Status: 'N',
		Loader: "bootstrap",
		Data: &classloader.ClData{
			Name:            className,
			SuperclassIndex: objIdx,
			Fields: []classloader.Field{
				{
					Name:     1,
					Desc:     2,
					IsStatic: true,
				},
			},
			CP: classloader.CPool{
				Utf8Refs: []string{"", fieldName, "Ljava/lang/Object;"},
			},
		},
	})

	callerClass := "test/Caller"
	callerMeth := "run"
	callerType := "()V"
	fullCallerName := callerClass + "." + callerMeth + callerType

	// Setup caller frame executing GETSTATIC with an exception handler covering PC 0..10 targeting PC 20
	f := frames.CreateFrame(5)
	f.Thread = 1
	f.ClName = callerClass
	f.MethName = callerMeth
	f.MethType = callerType
	f.Meth = []byte{opcodes.GETSTATIC, 0x00, 0x01} // PC 0..2: GETSTATIC CP[1]
	for len(f.Meth) < 25 {
		f.Meth = append(f.Meth, opcodes.NOP)
	}
	f.Meth[20] = opcodes.ASTORE_0 // catch block starts at PC 20
	f.Meth[21] = opcodes.RETURN   // PC 21

	cp := &classloader.CPool{}
	cp.CpIndex = make([]classloader.CpEntry, 3)
	cp.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0}
	cp.FieldRefs = make([]classloader.ResolvedFieldEntry, 1)
	cp.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      className,
		FldName:     fieldName,
		FldType:     "Ljava/lang/Object;",
	}

	// CP[2] = ClassRef for java/lang/Exception
	excClassName := "java/lang/Exception"
	excClassIdx := stringPool.GetStringIndex(&excClassName)
	cp.ClassRefs = []uint32{excClassIdx}
	cp.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}

	f.CP = cp

	// Register caller Java method in MTable with exception table
	jm := classloader.JmEntry{
		Exceptions: []classloader.CodeException{
			{
				StartPc:   0,
				EndPc:     10,
				HandlerPc: 20,
				CatchType: 2, // points to CP entry 2
			},
		},
		Cp: cp,
	}
	classloader.AddEntry(&classloader.MTable, fullCallerName, classloader.MTentry{Meth: jm, MType: 'J'})

	fs := frames.CreateFrameStack()
	f.FrameStack = fs
	fs.PushFront(f)

	glob := globals.GetGlobalRef()
	glob.FuncInstantiateClass = func(className string, fs *list.List) (any, error) {
		obj := object.MakeEmptyObject()
		return obj, nil
	}
	glob.FuncFillInStackTrace = func(params []any) any {
		return nil
	}
	glob.Threads = make(map[int]interface{})
	th := object.MakeEmptyObject()
	th.FieldTable["framestack"] = object.Field{Fvalue: fs}
	glob.Threads[1] = th

	ret := doGetStatic(f, 0)
	if ret != RESUME_HERE {
		t.Fatalf("Expected doGetStatic to return RESUME_HERE (%d), got %d", RESUME_HERE, ret)
	}

	if f.PC != 20 {
		t.Fatalf("Expected caller frame PC to be set to catch handler PC (20), got %d", f.PC)
	}

	if f.TOS != 0 {
		t.Fatalf("Expected TOS to be 0 (exception object on stack), got %d", f.TOS)
	}
}

func TestPutStaticClinitTrapExceptionCaught(t *testing.T) {
	globals.InitGlobals("testWithoutShutdown")
	globals.TraceInst = false
	classloader.InitMethodArea()
	initTable = sync.Map{}

	className := "test/TrapClassPut"
	fieldName := "field1"

	clinitSig := className + ".<clinit>()V"
	gm := ghelpers.GMeth{
		ParamSlots: 0,
		GFunction:  ghelpers.TrapClass,
	}
	ghelpers.MethodSignatures[clinitSig] = gm
	classloader.AddEntry(&classloader.MTable, clinitSig, classloader.MTentry{Meth: gm, MType: 'G'})

	objClass := "java/lang/Object"
	objIdx := stringPool.GetStringIndex(&objClass)
	classloader.MethAreaInsert(className, &classloader.Klass{
		Status: 'N',
		Loader: "bootstrap",
		Data: &classloader.ClData{
			Name:            className,
			SuperclassIndex: objIdx,
			Fields: []classloader.Field{
				{
					Name:     1,
					Desc:     2,
					IsStatic: true,
				},
			},
			CP: classloader.CPool{
				Utf8Refs: []string{"", fieldName, "I"},
			},
		},
	})

	callerClass := "test/CallerPut"
	callerMeth := "run"
	callerType := "()V"
	fullCallerName := callerClass + "." + callerMeth + callerType

	f := frames.CreateFrame(5)
	f.Thread = 1
	f.ClName = callerClass
	f.MethName = callerMeth
	f.MethType = callerType
	f.Meth = []byte{opcodes.PUTSTATIC, 0x00, 0x01}
	for len(f.Meth) < 25 {
		f.Meth = append(f.Meth, opcodes.NOP)
	}
	f.Meth[20] = opcodes.ASTORE_0
	f.Meth[21] = opcodes.RETURN

	cp := &classloader.CPool{}
	cp.CpIndex = make([]classloader.CpEntry, 3)
	cp.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0}
	cp.FieldRefs = make([]classloader.ResolvedFieldEntry, 1)
	cp.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      className,
		FldName:     fieldName,
		FldType:     "I",
	}

	excClassName := "java/lang/Exception"
	excClassIdx := stringPool.GetStringIndex(&excClassName)
	cp.ClassRefs = []uint32{excClassIdx}
	cp.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}

	f.CP = cp

	jm := classloader.JmEntry{
		Exceptions: []classloader.CodeException{
			{
				StartPc:   0,
				EndPc:     10,
				HandlerPc: 20,
				CatchType: 2,
			},
		},
		Cp: cp,
	}
	classloader.AddEntry(&classloader.MTable, fullCallerName, classloader.MTentry{Meth: jm, MType: 'J'})

	fs := frames.CreateFrameStack()
	f.FrameStack = fs
	fs.PushFront(f)

	glob := globals.GetGlobalRef()
	glob.FuncInstantiateClass = func(className string, fs *list.List) (any, error) {
		obj := object.MakeEmptyObject()
		return obj, nil
	}
	glob.FuncFillInStackTrace = func(params []any) any {
		return nil
	}
	glob.Threads = make(map[int]interface{})
	th := object.MakeEmptyObject()
	th.FieldTable["framestack"] = object.Field{Fvalue: fs}
	glob.Threads[1] = th

	ret := doPutStatic(f, 0)
	if ret != RESUME_HERE {
		t.Fatalf("Expected doPutStatic to return RESUME_HERE (%d), got %d", RESUME_HERE, ret)
	}

	if f.PC != 20 {
		t.Fatalf("Expected caller frame PC to be set to catch handler PC (20), got %d", f.PC)
	}
}

func TestNewClinitTrapExceptionCaught(t *testing.T) {
	globals.InitGlobals("testWithoutShutdown")
	globals.TraceInst = false
	classloader.InitMethodArea()
	initTable = sync.Map{}

	className := "test/TrapClassNew"

	clinitSig := className + ".<clinit>()V"
	gm := ghelpers.GMeth{
		ParamSlots: 0,
		GFunction:  ghelpers.TrapClass,
	}
	ghelpers.MethodSignatures[clinitSig] = gm
	classloader.AddEntry(&classloader.MTable, clinitSig, classloader.MTentry{Meth: gm, MType: 'G'})

	objClass := "java/lang/Object"
	objIdx := stringPool.GetStringIndex(&objClass)
	classloader.MethAreaInsert(className, &classloader.Klass{
		Status: 'N',
		Loader: "bootstrap",
		Data: &classloader.ClData{
			Name:            className,
			SuperclassIndex: objIdx,
		},
	})

	callerClass := "test/CallerNew"
	callerMeth := "run"
	callerType := "()V"
	fullCallerName := callerClass + "." + callerMeth + callerType

	f := frames.CreateFrame(5)
	f.Thread = 1
	f.ClName = callerClass
	f.MethName = callerMeth
	f.MethType = callerType
	f.Meth = []byte{opcodes.NEW, 0x00, 0x01}
	for len(f.Meth) < 25 {
		f.Meth = append(f.Meth, opcodes.NOP)
	}
	f.Meth[20] = opcodes.ASTORE_0
	f.Meth[21] = opcodes.RETURN

	cp := &classloader.CPool{}
	cp.CpIndex = make([]classloader.CpEntry, 3)
	classNameCopy := className
	classNameIdx := stringPool.GetStringIndex(&classNameCopy)
	cp.ClassRefs = []uint32{classNameIdx}
	cp.CpIndex[1] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}

	excClassName := "java/lang/Exception"
	excClassIdx := stringPool.GetStringIndex(&excClassName)
	cp.ClassRefs = append(cp.ClassRefs, excClassIdx)
	cp.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 1}

	f.CP = cp

	jm := classloader.JmEntry{
		Exceptions: []classloader.CodeException{
			{
				StartPc:   0,
				EndPc:     10,
				HandlerPc: 20,
				CatchType: 2,
			},
		},
		Cp: cp,
	}
	classloader.AddEntry(&classloader.MTable, fullCallerName, classloader.MTentry{Meth: jm, MType: 'J'})

	fs := frames.CreateFrameStack()
	f.FrameStack = fs
	fs.PushFront(f)

	glob := globals.GetGlobalRef()
	glob.FuncInstantiateClass = func(className string, fs *list.List) (any, error) {
		obj := object.MakeEmptyObject()
		return obj, nil
	}
	glob.FuncFillInStackTrace = func(params []any) any {
		return nil
	}
	glob.Threads = make(map[int]interface{})
	th := object.MakeEmptyObject()
	th.FieldTable["framestack"] = object.Field{Fvalue: fs}
	glob.Threads[1] = th

	ret := doNew(f, 0)
	if ret != RESUME_HERE {
		t.Fatalf("Expected doNew to return RESUME_HERE (%d), got %d", RESUME_HERE, ret)
	}

	if f.PC != 20 {
		t.Fatalf("Expected caller frame PC to be set to catch handler PC (20), got %d", f.PC)
	}
}
