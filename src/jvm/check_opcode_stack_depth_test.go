package jvm

// Adjust the package name and import paths to match the rest of your test files.
import (
	"jacobin/src/frames"
	"jacobin/src/opcodes"
	"testing"
)

// These tests check how many operand-stack entries a long or double occupies.
// Expected: exactly one. So after the first long/double lands on an empty stack, TOS must be 0.

// LCONST_1: push long constant 1 onto an empty stack
func TestLconst1StackDepth(t *testing.T) {
	f := newFrame(opcodes.LCONST_1)
	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	if f.TOS != 0 {
		t.Errorf("TestLconst1StackDepth: expected TOS == 0 (one entry), got: %d", f.TOS)
	}
	if f.TOS < 0 {
		return
	}

	value, ok := pop(f).(int64)
	if !ok {
		t.Errorf("TestLconst1StackDepth: expected an int64 on top of the stack")
	} else if value != 1 {
		t.Errorf("TestLconst1StackDepth: expected value 1, got: %d", value)
	}
}

// LLOAD_0: load a long from local variable 0 onto an empty stack
func TestLload0StackDepth(t *testing.T) {
	f := newFrame(opcodes.LLOAD_0)
	if len(f.Locals) < 2 { // a long occupies local indexes 0 and 1 in the bytecode's numbering
		f.Locals = make([]any, 2)
	}
	f.Locals[0] = int64(0x123456789A) // wider than 32 bits on purpose
	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	if f.TOS != 0 {
		t.Errorf("TestLload0StackDepth: expected TOS == 0 (one entry), got: %d", f.TOS)
	}
	if f.TOS < 0 {
		return
	}

	value, ok := pop(f).(int64)
	if !ok {
		t.Errorf("TestLload0StackDepth: expected an int64 on top of the stack")
	} else if value != 0x123456789A {
		t.Errorf("TestLload0StackDepth: expected value 0x123456789A, got: %X", value)
	}
}

// DCONST_1: push double constant 1.0 onto an empty stack
func TestDconst1StackDepth(t *testing.T) {
	f := newFrame(opcodes.DCONST_1)
	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	if f.TOS != 0 {
		t.Errorf("TestDconst1StackDepth: expected TOS == 0 (one entry), got: %d", f.TOS)
	}
	if f.TOS < 0 {
		return
	}

	value, ok := pop(f).(float64)
	if !ok {
		t.Errorf("TestDconst1StackDepth: expected a float64 on top of the stack")
	} else if value != 1.0 {
		t.Errorf("TestDconst1StackDepth: expected value 1.0, got: %f", value)
	}
}

// LADD: two longs in (one entry each), one long out
func TestLaddStackDepth(t *testing.T) {
	f := newFrame(opcodes.LADD)
	push(f, int64(0x100000000)) // wider than 32 bits on purpose
	push(f, int64(2))
	fs := frames.CreateFrameStack()
	fs.PushFront(f)
	interpret(fs)

	if f.TOS != 0 {
		t.Errorf("TestLaddStackDepth: expected TOS == 0 (one entry), got: %d", f.TOS)
	}
	if f.TOS < 0 {
		return
	}

	value, ok := pop(f).(int64)
	if !ok {
		t.Errorf("TestLaddStackDepth: expected an int64 on top of the stack")
	} else if value != 0x100000002 {
		t.Errorf("TestLaddStackDepth: expected value 0x100000002, got: %X", value)
	}
}

// DUP2: two entries only
func TestDup2OnTwoEntries(t *testing.T) {
	fr := newFrame(opcodes.DUP2)
	push(fr, int64(1)) // one entry only
	push(fr, int64(2)) // one entry only
	fs := frames.CreateFrameStack()
	fs.PushFront(fr)
	interpret(fs)

	if fr.TOS != 3 {
		t.Errorf("TestDup2OnTwoEntries: expected TOS == 3 (one entry), got: %d", fr.TOS)
	}

	value, ok := pop(fr).(int64)
	if !ok {
		t.Errorf("TestDup2OnTwoEntries: expected an int64 on top of the stack")
	} else if value != 2 {
		t.Errorf("TestDup2OnTwoEntries: expected value 2, got: %X", value)
	}
}
