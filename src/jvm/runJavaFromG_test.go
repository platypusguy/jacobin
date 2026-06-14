package jvm

import (
	"container/list"
	"jacobin/src/classloader"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/opcodes"
	"testing"
)

func TestRunJavaFromG_ArgsCount(t *testing.T) {
	globals.InitGlobals("test")
	classloader.InitMethodArea()
	fs := list.New()

	baseFrame := frames.CreateFrame(1)
	baseFrame.Thread = 1
	fs.PushFront(baseFrame)

	// We removed the arg count check, so this should try to execute and fail
	// because "Dummy" class is not found, but it shouldn't panic on arg count.
	RunJavaFromG(fs, "Dummy", "meth", "()V", 1, 2, 3)
}

func TestRunJavaFromG_Success_4Args(t *testing.T) {
	globals.InitGlobals("test")
	classloader.InitMethodArea()
	fs := list.New()

	baseFrame := frames.CreateFrame(1)
	baseFrame.Thread = 1
	fs.PushFront(baseFrame)

	clName := "com/test/Dummy4"
	methName := "testMeth"
	methType := "()V"
	methFQN := clName + "." + methName + methType

	// Create Klass and insert into MethArea to bypass LoadClassFromNameOnly
	k := &classloader.Klass{
		Status: 'I', // Initialized
		Data: &classloader.ClData{
			Name: clName,
		},
	}
	classloader.MethAreaInsert(clName, k)

	jme := classloader.JmEntry{
		MaxStack:  4,
		MaxLocals: 4,
		Code:      []byte{opcodes.RETURN},
		Cp:        &classloader.CPool{},
	}

	classloader.AddEntry(&classloader.MTable, methFQN, classloader.MTentry{
		Meth:  jme,
		MType: 'J',
	})

	args := []any{int64(1), int64(2), int64(3), int64(4)}
	initializeDispatchTable()
	RunJavaFromG(fs, clName, methName, methType, args...)

	if fs.Len() != 1 {
		t.Errorf("Expected base frame to remain on stack, got length %d", fs.Len())
	}
}

func TestRunJavaFromG_Success_5Args(t *testing.T) {
	globals.InitGlobals("test")
	classloader.InitMethodArea()
	fs := list.New()

	// Create a base frame for thread ID
	baseFrame := frames.CreateFrame(1)
	baseFrame.Thread = 1
	fs.PushFront(baseFrame)

	clName := "com/test/Dummy5"
	methName := "testMeth"
	methType := "()V"
	methFQN := clName + "." + methName + methType

	// Create Klass and insert into MethArea
	k := &classloader.Klass{
		Status: 'I',
		Data: &classloader.ClData{
			Name: clName,
		},
	}
	classloader.MethAreaInsert(clName, k)

	// Mock JmEntry with a RETURN instruction
	jme := classloader.JmEntry{
		MaxStack:  5,
		MaxLocals: 5,
		Code:      []byte{opcodes.RETURN},
		Cp:        &classloader.CPool{},
	}

	classloader.AddEntry(&classloader.MTable, methFQN, classloader.MTentry{
		Meth:  jme,
		MType: 'J',
	})

	// RunJavaFromG expects at least 4 arguments
	args := []any{int64(10), int64(20), int64(30), int64(40), int64(50)}

	// We need to initialize the dispatch table for interpret to work
	initializeDispatchTable()

	// This should run and complete because of the RETURN instruction
	RunJavaFromG(fs, clName, methName, methType, args...)

	if fs.Len() != 1 {
		t.Errorf("Expected base frame to remain on stack, got length %d", fs.Len())
	}
}

func TestRunJavaFromG_NotFound(t *testing.T) {
	globals.InitGlobals("test")
	classloader.InitMethodArea()
	fs := list.New()

	baseFrame := frames.CreateFrame(1)
	baseFrame.Thread = 1
	fs.PushFront(baseFrame)

	// In test mode, ThrowEx should return to RunJavaFromG, which then returns to us.
	// No panic should occur if we don't mock it to panic.

	// FetchMethodAndCP will fail to find "MissingClass"
	RunJavaFromG(fs, "MissingClass", "meth", "()V", 1, 2, 3, 4, 5)

	// If it reached here without panicking, it means RunJavaFromG handled the error via ThrowEx
	// which in test mode just prints to stderr and returns.
}
