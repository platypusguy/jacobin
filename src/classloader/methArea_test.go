/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"jacobin/globals"
	"jacobin/trace"
	"jacobin/types"
	"sync"
	"testing"
)

// Note: many MethArea functions are tested in classes_test,go
// These tests simply fill in untested functions, rather than duplicate those tests

func TestMethAreadDelete(t *testing.T) {
	MethArea = &sync.Map{}
	methAreaSize = 0
	currLen := MethAreaSize()
	if currLen != 0 {
		t.Errorf("Expecting MethArea size of 0, got: %d", currLen)
	}

	k := Klass{
		Status: 0,
		Loader: "",
		Data:   &ClData{},
	}
	k.Data.Name = "testClass1"
	k.Data.SuperclassIndex = types.ObjectPoolStringIndex
	k.Loader = "testloader"
	k.Status = 'F'
	MethAreaInsert("TestEntry1", &k)
	MethAreaInsert("TestEntry2", &k)
	currLen = MethAreaSize()
	if MethAreaSize() != 2 {
		t.Errorf("Expecting MethArea size of 2, got: %d", currLen)
	}

	MethAreaDelete("TestEntry1")
	newLen := MethAreaSize()
	if newLen != 1 {
		t.Errorf("Expected post-deletion MethArea[] to have length of 1, got: %d",
			newLen)
	}
}

func TestMethAreadDeleteNonExistentEntry(t *testing.T) {
	MethArea = &sync.Map{}
	methAreaSize = 0
	currLen := MethAreaSize()
	if currLen != 0 {
		t.Errorf("Expecting MethArea size of 0, got: %d", currLen)
	}

	k := Klass{
		Status: 0,
		Loader: "",
		Data:   &ClData{},
	}
	k.Data.Name = "testClass1"
	k.Data.SuperclassIndex = types.ObjectPoolStringIndex
	k.Loader = "testloader"
	k.Status = 'F'
	MethAreaInsert("TestEntry", &k)
	currLen = MethAreaSize()
	if MethAreaSize() != 1 {
		t.Errorf("Expecting MethArea size of 1, got: %d", currLen)
	}

	// deleting a non-entry should not cause an error or reduce MethArea size
	MethAreaDelete("NoSuchEntry")
	newLen := MethAreaSize()
	if newLen != 1 {
		t.Errorf("Expected post-deletion MethArea[] to have length of 1, got: %d",
			newLen)
	}
}

func tryMethod(t *testing.T, class string, methodName string, methodType string) {
	t.Logf("tryMethod: Input class=%s, methodName=%s, methodType=%s\n", class, methodName, methodType)
	mte, err := FetchMethodAndCP(class, methodName, methodType)
	if err != nil {
		t.Errorf("tryMethod: FetchMethodAndCP failed: %s\n", err.Error())
		return
	}
	t.Logf("tryMethod: FetchMethodAndCP returned MType: %s\n", string(mte.MType))
}
func TestMethArea42(t *testing.T) {
	globals.InitGlobals("test")
	trace.Init()

	// Initialise JMODMAP
	JmodMapInit()
	t.Logf("JmodMapInit ok\n")
	mapSize := JmodMapSize()
	if mapSize < 1 {
		t.Errorf("TestMethArea42: JMODMAP size < 1\n")
		return
	}
	t.Logf("JMODMAP size is %d\n", mapSize)

	// Initialise classloader
	Init()
	t.Logf("classloader.Init ok\n")

	// Load base classes.
	LoadBaseClasses()

	// Find some specific class-methods.
	tryMethod(t, "java/io/PrintStream", "println", "(Ljava/lang/String;)V")
	tryMethod(t, "java/io/BufferedOutputStream", "<init>", "(Ljava/io/OutputStream;)V")
	tryMethod(t, "java/io/BufferedOutputStream", "<init>", "(Ljava/io/OutputStream;I)V")
	tryMethod(t, "java/io/InputStream", "<init>", "()V")
}
