package jvm

import (
	"io"
	"jacobin/src/classloader"
	"jacobin/src/exceptions"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/opcodes"
	"jacobin/src/statics"
	"jacobin/src/types"
	"os"
	"strings"
	"testing"
)

// _testPutStaticTestClass initializes globals, disables instruction tracing,
// resets the method area, and registers the "test" class as already
// initialized, so callers no longer have to repeat this boilerplate.
func _testPutStaticTestClass() {
	globals.InitGlobals("test")
	globals.TraceInst = false
	classloader.InitMethodArea()
	_testPutStaticRegisterClass("test")
}

// _testPutStaticRegisterClass registers className in the method area and
// marks it as already initialized, so InitializeClass (called internally by
// doPutStatic) treats it as a known, ready-to-use class instead of trying to
// load it from the classpath.
func _testPutStaticRegisterClass(className string) {
	classloader.MethAreaInsert(className, &classloader.Klass{
		Status: 'N',
		Loader: "bootstrap",
		Data:   &classloader.ClData{Name: className},
	})
	getClassInit(className).done.Store(true)
}

// PUTSTATIC: Update a static field -- invalid b/c does not point to a field ref in the CP
// This test is now performed in codeCheck.go (and so, deleted here)

// Helper: sets up a frame and constant pool so bytecodes 0..2 form: PUTSTATIC, 0x00, 0x01
// CP[1] -> FieldRef slot 0 with (className,fldName,fldSig)
func _setupPutStaticFrame(className, fldName, fldSig string) (*frames.Frame, *classloader.CPool) {
	f := newFrame(opcodes.PUTSTATIC)
	f.Meth = append(f.Meth, 0x00, 0x01)

	cp := &classloader.CPool{}
	cp.CpIndex = make([]classloader.CpEntry, 2)
	cp.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0}
	cp.FieldRefs = make([]classloader.ResolvedFieldEntry, 1)
	cp.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      className,
		FldName:     fldName,
		FldType:     fldSig,
	}
	f.CP = cp
	return f, cp
}

// PUTSTATIC: Update a static field, an int, successfully
func TestPutStaticInt(t *testing.T) {
	_testPutStaticTestClass()

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP
	push(f, int64(420))

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0} // should be a field ref

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      "test",
		FldName:     "field1",
		FldType:     "I",
	}
	f.CP = &CP

	statics.LoadProgramStatics()
	statics.AddStatic("test.field1", statics.Static{
		Type:  "I",
		Value: 42,
	})

	fs := frames.CreateFrameStack()
	fs.PushFront(f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg != "" {
		t.Errorf("PUTSTATIC: Got unexpected error msg: \n%s", errMsg)
	}

	val := statics.GetStaticValue("test", "field1").(int64)
	if val != 420 {
		t.Errorf("PUTSTATIC: Expected static value to be 420, got: %d", val)
	}
}

// PUTSTATIC: Update a static field, an int, successfully (same as previous test, with tracing on)
func TestPutStaticIntWithTrace(t *testing.T) {
	_testPutStaticTestClass()
	globals.TraceInst = true

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP
	push(f, int64(420))

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0} // should be a field ref

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      "test",
		FldName:     "field1",
		FldType:     "I",
	}
	f.CP = &CP

	statics.LoadProgramStatics()
	statics.AddStatic("test.field1", statics.Static{
		Type:  "I",
		Value: 42,
	})

	fs := frames.CreateFrameStack()
	fs.PushFront(f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if !strings.Contains(errMsg, "PUTSTATIC") || !strings.Contains(errMsg, "field1") {
		t.Errorf("PUTSTATIC: Got unexpected error msg: \n%s", errMsg)
	}

	val := statics.GetStaticValue("test", "field1").(int64)
	if val != 420 {
		t.Errorf("PUTSTATIC: Expected static value to be 420, got: %d", val)
	}
}

// PUTSTATIC: Update a static field, a boolean, successfully
func TestPutStaticBool(t *testing.T) {
	_testPutStaticTestClass()

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP
	push(f, types.JavaBoolTrue)

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0} // should be a field ref

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      "test",
		FldName:     "field1",
		FldType:     types.Bool,
	}
	f.CP = &CP

	statics.LoadProgramStatics()
	statics.AddStatic("test.field1", statics.Static{
		Type:  "Z",
		Value: types.JavaBoolFalse,
	})

	fs := frames.CreateFrameStack()
	fs.PushFront(f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg != "" {
		t.Errorf("PUTSTATIC: Got unexpected error msg: \n%s", errMsg)
	}

	val := statics.GetStaticValue("test", "field1").(int64)
	if val != types.JavaBoolTrue {
		t.Errorf("PUTSTATIC: Expected static value to be true (1), got: %d", val)
	}
}

// PUTSTATIC: Update a static field, a byte, successfully
func TestPutStaticByte(t *testing.T) {
	_testPutStaticTestClass()

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP
	push(f, byte('A'))

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0} // should be a field ref

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      "test",
		FldName:     "field1",
		FldType:     types.Byte,
	}
	f.CP = &CP

	statics.LoadProgramStatics()
	statics.AddStatic("test.field1", statics.Static{
		Type:  types.Byte,
		Value: byte('B'),
	})

	fs := frames.CreateFrameStack()
	fs.PushFront(f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg != "" {
		t.Errorf("PUTSTATIC: Got unexpected error msg: \n%s", errMsg)
	}

	val := statics.GetStaticValue("test", "field1").(int64) // GeStaticValue converts bytes to int64s
	if rune(val) != 'A' {
		t.Errorf("PUTSTATIC: Expected static value to be 'A', got: %c", rune(val))
	}
}

// PUTSTATIC: Update a static field, a byte, successfully
func TestPutStaticJavaByte(t *testing.T) {
	_testPutStaticTestClass()

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP
	push(f, types.JavaByte('A'))

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0} // should be a field ref

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      "test",
		FldName:     "field1",
		FldType:     types.Byte,
	}
	f.CP = &CP

	statics.LoadProgramStatics()
	statics.AddStatic("test.field1", statics.Static{
		Type:  types.Byte,
		Value: types.JavaByte('B'),
	})

	fs := frames.CreateFrameStack()
	fs.PushFront(f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg != "" {
		t.Errorf("PUTSTATIC: Got unexpected error msg: \n%s", errMsg)
	}

	val := statics.GetStaticValue("test", "field1").(int64) // GeStaticValue converts bytes to int64s
	if rune(val) != 'A' {
		t.Errorf("PUTSTATIC: Expected static value to be 'A', got: %c", rune(val))
	}
}

// PUTSTATIC: Update a static field, a byte, successfully. This byte value is passed in as int64
func TestPutStaticByteAsInt64(t *testing.T) {
	_testPutStaticTestClass()

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP
	push(f, int64('A'))

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0} // should be a field ref

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      "test",
		FldName:     "field1",
		FldType:     types.Byte,
	}
	f.CP = &CP

	statics.LoadProgramStatics()
	statics.AddStatic("test.field1", statics.Static{
		Type:  types.Byte,
		Value: int64('B'),
	})

	fs := frames.CreateFrameStack()
	fs.PushFront(f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg != "" {
		t.Errorf("PUTSTATIC: Got unexpected error msg: \n%s", errMsg)
	}

	val := statics.GetStaticValue("test", "field1").(int64) // GeStaticValue converts bytes to int64s
	if rune(val) != 'A' {
		t.Errorf("PUTSTATIC: Expected static value to be 'A', got: %c", rune(val))
	}
}

// PUTSTATIC: Update a static field, an float/double, successfully
func TestPutStaticFloat(t *testing.T) {
	_testPutStaticTestClass()

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP
	push(f, float64(420.1))

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0} // should be a field ref

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      "test",
		FldName:     "field1",
		FldType:     "F",
	}
	f.CP = &CP

	statics.LoadProgramStatics()
	statics.AddStatic("test.field1", statics.Static{
		Type:  "F",
		Value: 42.0,
	})

	fs := frames.CreateFrameStack()
	fs.PushFront(f) // push the new frame
	interpret(fs)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	errMsg := string(msg)
	if errMsg != "" {
		t.Errorf("PUTSTATIC: Got unexpected error msg: \n%s", errMsg)
	}

	val := statics.GetStaticValue("test", "field1").(float64)
	if val != 420.1 {
		t.Errorf("PUTSTATIC: Expected static value to be 420.9, got: %f", val)
	}
}

// PUTSTATIC: this should bonk because the class of the static cannot be found/loaded
func TestPutStaticInvalidNoSuchClass(t *testing.T) {
	_testPutStaticTestClass()
	globals.TraceInst = true

	classloader.InitMethodArea()
	statics.Statics = make(map[string]statics.Static)

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f := newFrame(opcodes.PUTSTATIC)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP
	push(f, float64(420.1))

	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.FieldRef, Slot: 0} // should be a field ref

	// now create the pointed-to FieldRef
	CP.FieldRefs = make([]classloader.ResolvedFieldEntry, 1)
	CP.FieldRefs[0] = classloader.ResolvedFieldEntry{
		AccessFlags: 0,
		IsStatic:    true,
		IsFinal:     false,
		ClName:      "test",
		FldName:     "field1",
		FldType:     "F",
	}
	f.CP = &CP

	ret := doPutStatic(f, 0)

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	_, _ = io.ReadAll(r)
	os.Stderr = normalStderr

	if ret != exceptions.ERROR_OCCURRED {
		t.Errorf("TestPutStaticInvalidNoSuchClass: Expected ret=exceptions.ERROR_OCCURRED, observed: %d", ret)
		t.Log(string(msg))
	}
}

// PUTSTATIC: store boolean (normalize to int64 0/1)
func TestPutStaticStoreBoolean_NoTable_NoTypesType(t *testing.T) {
	_testPutStaticTestClass()
	statics.Statics = make(map[string]statics.Static)

	cls, fld := "Cbool", "b"
	_testPutStaticRegisterClass(cls)
	key := cls + "." + fld
	statics.AddStatic(key, statics.Static{Type: types.Bool, Value: int64(0)})

	f, _ := _setupPutStaticFrame(cls, fld, types.Bool)
	push(f, int64(3)) // doPutStatic will mask with & 0x01
	ret := doPutStatic(f, 0)
	if ret != 3 {
		t.Fatalf("expected ret=3, got %d", ret)
	}
	st := statics.Statics[key]
	if st.Type != types.Bool {
		t.Fatalf("type mismatch: %v", st.Type)
	}
	if v := st.Value.(int64); v != 1 {
		t.Fatalf("expected 1, got %v", v)
	}
}

// PUTSTATIC: store Char as int64
func TestPutStaticStoreChar_NoTable_NoTypesType(t *testing.T) {
	_testPutStaticTestClass()
	statics.Statics = make(map[string]statics.Static)

	cls, fld := "Cchar", "x"
	_testPutStaticRegisterClass(cls)
	key := cls + "." + fld
	statics.AddStatic(key, statics.Static{Type: types.Char, Value: int64(0)})

	f, _ := _setupPutStaticFrame(cls, fld, types.Char)
	push(f, int64(65))
	ret := doPutStatic(f, 0)
	if ret != 3 {
		t.Fatalf("expected ret=3, got %d", ret)
	}
	st := statics.Statics[key]
	if st.Type != types.Char {
		t.Fatalf("type mismatch: %v", st.Type)
	}
	if v := st.Value.(int64); v != 65 {
		t.Fatalf("expected 65, got %d", v)
	}
}

// PUTSTATIC: store Short as int64
func TestPutStaticStoreShort_NoTable_NoTypesType(t *testing.T) {
	_testPutStaticTestClass()
	statics.Statics = make(map[string]statics.Static)

	cls, fld := "Cshort", "x"
	_testPutStaticRegisterClass(cls)
	key := cls + "." + fld
	statics.AddStatic(key, statics.Static{Type: types.Short, Value: int64(0)})

	f, _ := _setupPutStaticFrame(cls, fld, types.Short)
	push(f, int64(-123))
	ret := doPutStatic(f, 0)
	if ret != 3 {
		t.Fatalf("expected ret=3, got %d", ret)
	}
	st := statics.Statics[key]
	if st.Type != types.Short {
		t.Fatalf("type mismatch: %v", st.Type)
	}
	if v := st.Value.(int64); v != -123 {
		t.Fatalf("expected -123, got %d", v)
	}
}

// PUTSTATIC: store Int as int64
func TestPutStaticStoreInt_NoTable_NoTypesType(t *testing.T) {
	_testPutStaticTestClass()
	statics.Statics = make(map[string]statics.Static)

	cls, fld := "Cint", "x"
	_testPutStaticRegisterClass(cls)
	key := cls + "." + fld
	statics.AddStatic(key, statics.Static{Type: types.Int, Value: int64(0)})

	f, _ := _setupPutStaticFrame(cls, fld, types.Int)
	push(f, int64(42))
	ret := doPutStatic(f, 0)
	if ret != 3 {
		t.Fatalf("expected ret=3, got %d", ret)
	}
	st := statics.Statics[key]
	if st.Type != types.Int {
		t.Fatalf("type mismatch: %v", st.Type)
	}
	if v := st.Value.(int64); v != 42 {
		t.Fatalf("expected 42, got %d", v)
	}
}

// PUTSTATIC: store Long as int64
func TestPutStaticStoreLong_NoTable_NoTypesType(t *testing.T) {
	_testPutStaticTestClass()
	statics.Statics = make(map[string]statics.Static)

	cls, fld := "Clong", "x"
	_testPutStaticRegisterClass(cls)
	key := cls + "." + fld
	statics.AddStatic(key, statics.Static{Type: types.Long, Value: int64(0)})

	f, _ := _setupPutStaticFrame(cls, fld, types.Long)
	push(f, int64(9876543210))
	ret := doPutStatic(f, 0)
	if ret != 3 {
		t.Fatalf("expected ret=3, got %d", ret)
	}
	st := statics.Statics[key]
	if st.Type != types.Long {
		t.Fatalf("type mismatch: %v", st.Type)
	}
	if v := st.Value.(int64); v != 9876543210 {
		t.Fatalf("expected 9876543210, got %d", v)
	}
}

// PUTSTATIC: store byte from int64
func TestPutStaticStoreByteFromInt64_NoTable_NoTypesType(t *testing.T) {
	_testPutStaticTestClass()
	statics.Statics = make(map[string]statics.Static)

	cls, fld := "Cbyte", "b"
	_testPutStaticRegisterClass(cls)
	key := cls + "." + fld
	statics.AddStatic(key, statics.Static{Type: types.Byte, Value: int64(0)})

	f, _ := _setupPutStaticFrame(cls, fld, types.Byte)
	push(f, int64(-7))
	_ = doPutStatic(f, 0)
	if v := statics.Statics[key].Value.(int64); v != -7 {
		t.Fatalf("expected -7, got %d", v)
	}
}

// PUTSTATIC: store byte from uint8
func TestPutStaticStoreByteFromUint8_NoTable_NoTypesType(t *testing.T) {
	_testPutStaticTestClass()
	statics.Statics = make(map[string]statics.Static)

	cls, fld := "Cbyte", "b"
	_testPutStaticRegisterClass(cls)
	key := cls + "." + fld
	statics.AddStatic(key, statics.Static{Type: types.Byte, Value: int64(0)})

	f, _ := _setupPutStaticFrame(cls, fld, types.Byte)
	push(f, uint8(250))
	_ = doPutStatic(f, 0)
	if v := statics.Statics[key].Value.(int64); v != 250 {
		t.Fatalf("expected 250, got %d", v)
	}
}

// PUTSTATIC: store byte from types.JavaByte
func TestPutStaticStoreByteFromJavaByte_NoTable_NoTypesType(t *testing.T) {
	_testPutStaticTestClass()
	statics.Statics = make(map[string]statics.Static)

	cls, fld := "Cbyte", "b"
	_testPutStaticRegisterClass(cls)
	key := cls + "." + fld
	statics.AddStatic(key, statics.Static{Type: types.Byte, Value: int64(0)})

	f, _ := _setupPutStaticFrame(cls, fld, types.Byte)
	push(f, types.JavaByte(-120))
	_ = doPutStatic(f, 0)
	if v := statics.Statics[key].Value.(int64); v != -120 {
		t.Fatalf("expected -120, got %d", v)
	}
}

// PUTSTATIC: store Float as float64
func TestPutStaticStoreFloat_NoTable_NoTypesType(t *testing.T) {
	_testPutStaticTestClass()
	statics.Statics = make(map[string]statics.Static)

	cls, fld := "Cfp", "y"
	_testPutStaticRegisterClass(cls)
	key := cls + "." + fld
	statics.AddStatic(key, statics.Static{Type: types.Float, Value: float64(0)})

	f, _ := _setupPutStaticFrame(cls, fld, types.Float)
	push(f, float64(3.5))
	ret := doPutStatic(f, 0)
	if ret != 3 {
		t.Fatalf("expected ret=3, got %d", ret)
	}
	st := statics.Statics[key]
	if st.Type != types.Float {
		t.Fatalf("type mismatch: %v", st.Type)
	}
	if v := st.Value.(float64); v != 3.5 {
		t.Fatalf("expected 3.5, got %v", v)
	}
}

// PUTSTATIC: store Double as float64
func TestPutStaticStoreDouble_NoTable_NoTypesType(t *testing.T) {
	_testPutStaticTestClass()
	statics.Statics = make(map[string]statics.Static)

	cls, fld := "Cdp", "y"
	_testPutStaticRegisterClass(cls)
	key := cls + "." + fld
	statics.AddStatic(key, statics.Static{Type: types.Double, Value: float64(0)})

	f, _ := _setupPutStaticFrame(cls, fld, types.Double)
	push(f, float64(-42.25))
	ret := doPutStatic(f, 0)
	if ret != 3 {
		t.Fatalf("expected ret=3, got %d", ret)
	}
	st := statics.Statics[key]
	if st.Type != types.Double {
		t.Fatalf("type mismatch: %v", st.Type)
	}
	if v := st.Value.(float64); v != -42.25 {
		t.Fatalf("expected -42.25, got %v", v)
	}
}

// PUTSTATIC: default branch (references) — nil coerces to object.Null
func TestPutStaticStoreRefNil_NoTable_NoTypesType(t *testing.T) {
	_testPutStaticTestClass()
	statics.Statics = make(map[string]statics.Static)

	cls, fld := "CrefN", "r"
	_testPutStaticRegisterClass(cls)
	key := cls + "." + fld
	statics.AddStatic(key, statics.Static{Type: types.Ref, Value: object.Null})

	f, _ := _setupPutStaticFrame(cls, fld, "Ljava/lang/Object;")
	push(f, nil)
	_ = doPutStatic(f, 0)
	st := statics.Statics[key]
	if st.Type != types.Ref {
		t.Fatalf("type mismatch: %v", st.Type)
	}
	if st.Value != object.Null {
		t.Fatalf("expected object.Null, got %v", st.Value)
	}
}

// PUTSTATIC: default branch (references) — *object.Object stored as-is
func TestPutStaticStoreRefObject_NoTable_NoTypesType(t *testing.T) {
	_testPutStaticTestClass()
	statics.Statics = make(map[string]statics.Static)

	cls, fld := "CrefO", "o"
	_testPutStaticRegisterClass(cls)
	key := cls + "." + fld
	statics.AddStatic(key, statics.Static{Type: types.Ref, Value: object.Null})

	f, _ := _setupPutStaticFrame(cls, fld, "Ljava/lang/Object;")
	obj := object.MakeEmptyObject()
	push(f, obj)
	_ = doPutStatic(f, 0)
	st := statics.Statics[key]
	if st.Type != types.Ref {
		t.Fatalf("type mismatch: %v", st.Type)
	}
	if st.Value != obj {
		t.Fatalf("expected same object pointer, got %v", st.Value)
	}
}

// PUTSTATIC: instantiate fails (ClassNotFound) path — direct doPutStatic call
func TestPutStaticClassNotFound_NoTable_NoTypesType(t *testing.T) {
	_testPutStaticTestClass()
	statics.Statics = make(map[string]statics.Static)

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f, _ := _setupPutStaticFrame("no/such/Class", "x", types.Int)
	push(f, int64(1))
	ret := doPutStatic(f, 0)

	_ = w.Close()
	_, _ = io.ReadAll(r) // drain
	os.Stderr = normalStderr

	if ret != exceptions.ERROR_OCCURRED {
		t.Fatalf("expected ERROR_OCCURRED, got %d", ret)
	}
}
