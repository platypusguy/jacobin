/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-3 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"container/list"
	"encoding/binary"
	"errors"
	"fmt"
	"jacobin/classloader"
	"jacobin/exceptions"
	"jacobin/frames"
	"jacobin/gfunction"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/object"
	"jacobin/opcodes"
	"jacobin/shutdown"
	"jacobin/statics"
	"jacobin/thread"
	"jacobin/types"
	"jacobin/util"
	"math"
	"runtime/debug"
	"strconv"
	"strings"
	"unsafe"
)

var MainThread thread.ExecThread

// StartExec is where execution begins. It initializes various structures, such as
// the MTable, then using the passed-in name of the starting class, finds its main() method
// in the method area (it's guaranteed to already be loaded), grabs the executable
// bytes, creates a thread of execution, pushes the main() frame onto the JVM stack
// and begins execution.
func StartExec(className string, mainThread *thread.ExecThread, globals *globals.Globals) error {

	MainThread = *mainThread
	// set tracing, if any
	tracing := false
	trace, exists := globals.Options["-trace"]
	if exists {
		tracing = trace.Set
	}
	MainThread.Trace = tracing

	me, err := classloader.FetchMethodAndCP(className, "main", "([Ljava/lang/String;)V")
	if err != nil {
		return errors.New("Class not found: " + className + ".main()")
	}

	m := me.Meth.(classloader.JmEntry)
	f := frames.CreateFrame(m.MaxStack + 2) // create a new frame (the +2 is arbitrary, but needed)
	f.Thread = MainThread.ID
	f.MethName = "main"
	f.MethType = "([Ljava/lang/String;)V"
	f.ClName = className
	f.CP = m.Cp                        // add its pointer to the class CP
	f.Meth = append(f.Meth, m.Code...) // copy the bytecodes over

	// allocate the local variables
	for k := 0; k < m.MaxLocals; k++ {
		f.Locals = append(f.Locals, 0)
	}

	// create the first thread and place its first frame on it
	// MainThread = *mainThread
	MainThread.Stack = frames.CreateFrameStack()
	// MainThread.ID = thread.AddThreadToTable(&MainThread, &globals.Threads)
	MainThread.Trace = tracing

	// must first instantiate the class, so that any static initializers are run
	_, instantiateError := InstantiateClass(className, MainThread.Stack)
	if instantiateError != nil {
		return errors.New("Error instantiating: " + className + ".main()")
	}

	if frames.PushFrame(MainThread.Stack, f) != nil {
		errMsg := "Memory error allocating frame on thread: " + strconv.Itoa(MainThread.ID)
		_ = log.Log(errMsg, log.SEVERE)
		return errors.New(errMsg)
	}

	if MainThread.Trace {
		traceInfo := fmt.Sprintf("StartExec: class=%s, meth=%s, maxStack=%d, maxLocals=%d, code size=%d",
			f.ClName, f.MethName, m.MaxStack, m.MaxLocals, len(m.Code))
		_ = log.Log(traceInfo, log.TRACE_INST)
	}

	err = runThread(&MainThread)
	if err != nil {
		statics.DumpStatics()
		return err
	}

	if MainThread.Trace {
		statics.DumpStatics()
	}

	return nil
}

// Point the thread to the top of the frame stack and tell it to run from there.
func runThread(t *thread.ExecThread) error {

	defer func() int {
		// only an untrapped panic gets us here
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			glob := globals.GetGlobalRef()
			glob.ErrorGoStack = stack
			exceptions.ShowPanicCause(r)
			exceptions.ShowFrameStack(t)
			exceptions.ShowGoStackTrace(nil)
			statics.DumpStatics()
			return shutdown.Exit(shutdown.APP_EXCEPTION)
		}
		return shutdown.OK
	}()

	for t.Stack.Len() > 0 {
		err := runFrame(t.Stack)
		if err != nil {
			exceptions.ShowFrameStack(t)
			if globals.GetGlobalRef().GoStackShown == false {
				exceptions.ShowGoStackTrace(nil)
				globals.GetGlobalRef().GoStackShown = true
			}
			return err
		}

		if t.Stack.Len() == 1 { // true when the last executed frame was main()
			return nil
		} else {
			t.Stack.Remove(t.Stack.Front()) // pop the frame off
		}
	}
	return nil
}

// runFrame() is the principal execution function in Jacobin. It first tests for a
// golang function in the present frame. If it is a golang function, it's sent to
// a different function for execution. Otherwise, bytecode interpretation takes
// place through a giant switch statement.
func runFrame(fs *list.List) error {
	glob := globals.GetGlobalRef()

frameInterpreter:
	// the current frame is always the head of the linked list of frames.
	// the next statement converts the address of that frame to the more readable 'f'
	f := fs.Front().Value.(*frames.Frame)

	// if the frame contains a golang method, execute it using runGframe(),
	// which returns a value (possibly nil) and an exceptions code. Presuming no exceptions,
	// if the return value (here, retval) is not nil, it is placed on the stack
	// of the calling frame.
	if f.Ftype == 'G' {
		retval, slotCount, err := runGframe(f)

		if retval != nil {
			f = fs.Front().Next().Value.(*frames.Frame)
			push(f, retval) // if slotCount = 1

			if slotCount == 2 {
				push(f, retval) // push a second time, if a long, double, etc.
			}
		}
		return err
	}

	// the frame's method is not a golang method, so it's Java bytecode, which
	// is interpreted in the rest of this function.
	for f.PC < len(f.Meth) {
		if MainThread.Trace && f.Meth[f.PC] != opcodes.IMPDEP2 {
			traceInfo := emitTraceData(f)
			_ = log.Log(traceInfo, log.TRACE_INST)
		}

		switch f.Meth[f.PC] { // cases listed in numerical value of opcode
		case opcodes.NOP:
			break
		case opcodes.ACONST_NULL: // 0x01   (push null onto opStack)
			// push(f, int64(0)) // replaced in JACOBIN-286
			push(f, object.Null)
		case opcodes.ICONST_M1: //	x02	(push -1 onto opStack)
			push(f, int64(-1))
		case opcodes.ICONST_0: // 	0x03	(push int 0 onto opStack)
			push(f, int64(0))
		case opcodes.ICONST_1: //  	0x04	(push int 1 onto opStack)
			push(f, int64(1))
		case opcodes.ICONST_2: //   0x05	(push 2 onto opStack)
			push(f, int64(2))
		case opcodes.ICONST_3: //   0x06	(push 3 onto opStack)
			push(f, int64(3))
		case opcodes.ICONST_4: //   0x07	(push 4 onto opStack)
			push(f, int64(4))
		case opcodes.ICONST_5: //   0x08	(push 5 onto opStack)
			push(f, int64(5))
		case opcodes.LCONST_0: //   0x09    (push long 0 onto opStack)
			push(f, int64(0)) // b/c longs take two slots on the stack, it's pushed twice
			push(f, int64(0))
		case opcodes.LCONST_1: //   0x0A    (push long 1 on to opStack)
			push(f, int64(1)) // b/c longs take two slots on the stack, it's pushed twice
			push(f, int64(1))
		case opcodes.FCONST_0: // 0x0B
			push(f, 0.0)
		case opcodes.FCONST_1: // 0x0C
			push(f, 1.0)
		case opcodes.FCONST_2: // 0x0D
			push(f, 2.0)
		case opcodes.DCONST_0: // 0x0E
			push(f, 0.0)
			push(f, 0.0)
		case opcodes.DCONST_1: // 0xoF
			push(f, 1.0)
			push(f, 1.0)
		case opcodes.BIPUSH: //	0x10	(push the following byte as an int onto the stack)
			wbyte := f.Meth[f.PC+1]
			wint64 := byteToInt64(wbyte)
			f.PC += 1
			push(f, wint64)
		case opcodes.SIPUSH: //	0x11	(create int from next two bytes and push the int)
			wbyte1 := f.Meth[f.PC+1]
			wbyte2 := f.Meth[f.PC+2]
			var wint64 int64
			if (wbyte1 & 0x80) == 0x80 { // Negative wbyte1 (left-most bit on)?
				// Negative wbyte1 : form wbytes = 6 0xFFs concatenated with the wbyte1 and wbyte2
				var wbytes = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00}
				wbytes[6] = wbyte1
				wbytes[7] = wbyte2
				// Form an int64 from the wbytes array
				// If you know C, this is equivalent to memcpy(&wint64, &wbytes, 8)
				wint64 = int64(binary.BigEndian.Uint64(wbytes))
			} else {
				// Not negative (left-most bit off) : just cast wbyte as an int64
				wint64 = (int64(wbyte1) * 256) + int64(wbyte2)
			}
			f.PC += 2
			push(f, wint64)
		case opcodes.LDC: // 	0x12   	(push constant from CP indexed by next byte)
			idx := f.Meth[f.PC+1]
			f.PC += 1

			CPe := classloader.FetchCPentry(f.CP.(*classloader.CPool), int(idx))
			if CPe.EntryType != 0 && // 0 = error
				// Note: an invalid CP entry causes a java.lang.Verify error and
				//       is caught before execution of the program begins.
				// This instruction does not load longs or doubles
				CPe.EntryType != classloader.DoubleConst &&
				CPe.EntryType != classloader.LongConst { // if no error
				if CPe.RetType == classloader.IS_INT64 {
					push(f, CPe.IntVal)
				} else if CPe.RetType == classloader.IS_FLOAT64 {
					push(f, CPe.FloatVal)
				} else if CPe.RetType == classloader.IS_STRUCT_ADDR {
					push(f, (*object.Object)(unsafe.Pointer(CPe.AddrVal)))
				} else if CPe.RetType == classloader.IS_STRING_ADDR {
					stringAddr :=
						object.CreateCompactStringFromGoString(CPe.StringVal)
					stringAddr.Klass = &object.StringClassName
					if classloader.MethAreaFetch(*stringAddr.Klass) == nil {
						// glob := globals.GetGlobalRef()
						glob.ErrorGoStack = string(debug.Stack())
						errMsg := fmt.Sprintf("LDC: MethAreaFetch could not find class java/lang/String")
						_ = log.Log(errMsg, log.SEVERE)
						return errors.New("LDC: MethAreaFetch could not find class java/lang/String")
					}
					push(f, stringAddr)
				}
			} else { // TODO: Determine what exception to throw
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "LDC: Invalid type for instruction"
				exceptions.Throw(exceptions.InaccessibleObjectException, errMsg)
				return errors.New(errMsg)
			}
		case opcodes.LDC_W: // 	0x13	(push constant from CP indexed by next two bytes)
			idx := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2])
			f.PC += 2

			CPe := classloader.FetchCPentry(f.CP.(*classloader.CPool), idx)
			if CPe.EntryType != 0 && // this instruction does not load longs or doubles
				CPe.EntryType != classloader.DoubleConst &&
				CPe.EntryType != classloader.LongConst { // if no error
				if CPe.RetType == classloader.IS_INT64 {
					push(f, CPe.IntVal)
				} else if CPe.RetType == classloader.IS_FLOAT64 {
					push(f, CPe.FloatVal)
				} else if CPe.RetType == classloader.IS_STRUCT_ADDR {
					push(f, (*object.Object)(unsafe.Pointer(CPe.AddrVal)))
				} else if CPe.RetType == classloader.IS_STRING_ADDR {
					stringAddr :=
						object.CreateCompactStringFromGoString(CPe.StringVal)
					stringAddr.Klass = &object.StringClassName
					if classloader.MethAreaFetch(*stringAddr.Klass) == nil {
						// glob := globals.GetGlobalRef()
						glob.ErrorGoStack = string(debug.Stack())
						errMsg := fmt.Sprintf("LDC_W: MethAreaFetch could not find class java/lang/String")
						_ = log.Log(errMsg, log.SEVERE)
						return errors.New("LDC_W: MethAreaFetch could not find class java/lang/String")
					}
					push(f, stringAddr)
				}
			} else { // TODO: Determine what exception to throw
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "LDC_W: Invalid type for instruction"
				exceptions.Throw(exceptions.InaccessibleObjectException, errMsg)
				return errors.New(errMsg)
			}
		case opcodes.LDC2_W: // 0x14 	(push long or double from CP indexed by next two bytes)
			idx := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2])
			f.PC += 2

			CPe := classloader.FetchCPentry(f.CP.(*classloader.CPool), idx)
			if CPe.RetType == classloader.IS_INT64 { // push value twice (due to 64-bit width)
				push(f, CPe.IntVal)
				push(f, CPe.IntVal)
			} else if CPe.RetType == classloader.IS_FLOAT64 {
				push(f, CPe.FloatVal)
				push(f, CPe.FloatVal)
			} else { // TODO: Determine what exception to throw
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "LDC2_W: Invalid type for LDC2_W instruction"
				exceptions.Throw(exceptions.InaccessibleObjectException, errMsg)
				return errors.New(errMsg)
			}
		case opcodes.ILOAD, // 0x15	(push int from local var, using next byte as index)
			opcodes.FLOAD, //  0x17 (push float from local var, using next byte as index)
			opcodes.ALOAD: //  0x19 (push ref from local var, using next byte as index)
			index := int(f.Meth[f.PC+1])
			f.PC += 1
			push(f, f.Locals[index])
		case opcodes.LLOAD: // 0x16 (push long from local var, using next byte as index)
			index := int(f.Meth[f.PC+1])
			f.PC += 1
			val := f.Locals[index].(int64)
			push(f, val)
			push(f, val) // push twice due to item being 64 bits wide
		case opcodes.DLOAD: // 0x18 (push double from local var, using next byte as index)
			index := int(f.Meth[f.PC+1])
			f.PC += 1
			val := f.Locals[index].(float64)
			push(f, val)
			push(f, val) // push twice due to item being 64 bits wide
		case opcodes.ILOAD_0: // 	0x1A    (push local variable 0)
			push(f, f.Locals[0].(int64))
		case opcodes.ILOAD_1: //    OX1B    (push local variable 1)
			push(f, f.Locals[1].(int64))
		case opcodes.ILOAD_2: //    0X1C    (push local variable 2)
			push(f, f.Locals[2].(int64))
		case opcodes.ILOAD_3: //  	0x1D   	(push local variable 3)
			push(f, f.Locals[3].(int64))
		// LLOAD use two slots, so the same value is pushed twice
		case opcodes.LLOAD_0: //	0x1E	(push local variable 0, as long)
			push(f, f.Locals[0].(int64))
			push(f, f.Locals[0].(int64))
		case opcodes.LLOAD_1: //	0x1F	(push local variable 1, as long)
			push(f, f.Locals[1].(int64))
			push(f, f.Locals[1].(int64))
		case opcodes.LLOAD_2: //	0x20	(push local variable 2, as long)
			push(f, f.Locals[2].(int64))
			push(f, f.Locals[2].(int64))
		case opcodes.LLOAD_3: //	0x21	(push local variable 3, as long)
			push(f, f.Locals[3].(int64))
			push(f, f.Locals[3].(int64))
		case opcodes.FLOAD_0: // 0x22
			push(f, f.Locals[0])
		case opcodes.FLOAD_1: // 0x23
			push(f, f.Locals[1])
		case opcodes.FLOAD_2: // 0x24
			push(f, f.Locals[2])
		case opcodes.FLOAD_3: // 0x25
			push(f, f.Locals[3])
		case opcodes.DLOAD_0: //	0x26	(push local variable 0, as double)
			push(f, f.Locals[0])
			push(f, f.Locals[0])
		case opcodes.DLOAD_1: //	0x27	(push local variable 1, as double)
			push(f, f.Locals[1])
			push(f, f.Locals[1])
		case opcodes.DLOAD_2: //	0x28	(push local variable 2, as double)
			push(f, f.Locals[2])
			push(f, f.Locals[2])
		case opcodes.DLOAD_3: //	0x29	(push local variable 3, as double)
			push(f, f.Locals[3])
			push(f, f.Locals[3])
		case opcodes.ALOAD_0: //	0x2A	(push reference stored in local variable 0)
			push(f, f.Locals[0])
		case opcodes.ALOAD_1: //	0x2B	(push reference stored in local variable 1)
			push(f, f.Locals[1])
		case opcodes.ALOAD_2: //	0x2C    (push reference stored in local variable 2)
			push(f, f.Locals[2])
		case opcodes.ALOAD_3: //	0x2D	(push reference stored in local variable 3)
			push(f, f.Locals[3])
		case opcodes.IALOAD, //		0x2E	(push contents of an int array element)
			opcodes.CALOAD, //		0x34	(push contents of a (two-byte) char array element)
			opcodes.SALOAD: //		0x35    (push contents of a short array element)
			index := pop(f).(int64)
			iAref := pop(f).(*object.Object) // ptr to array object
			if iAref == object.Null {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "I/C/SALOAD: Invalid (null) reference to an array"
				exceptions.Throw(exceptions.NullPointerException, errMsg)
				return errors.New(errMsg)
			}

			array := *(iAref.Fields[0].Fvalue).(*[]int64)

			if index >= int64(len(array)) {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "IALOAD: Invalid array subscript"
				exceptions.Throw(exceptions.ArrayIndexOutOfBoundsException, errMsg)
				return errors.New(errMsg)
			}
			var value = array[index]
			push(f, value)
		case opcodes.LALOAD: //		0x2F	(push contents of a long array element)
			index := pop(f).(int64)
			iAref := pop(f).(*object.Object) // ptr to array object
			if iAref == nil {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "LALOAD: Invalid (null) reference to an array"
				exceptions.Throw(exceptions.NullPointerException, errMsg)
				return errors.New(errMsg)
			}

			array := *(iAref.Fields[0].Fvalue).(*[]int64)
			if index >= int64(len(array)) {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				exceptions.Throw(exceptions.ArrayIndexOutOfBoundsException,
					"LALOAD: Invalid array subscript")
				return errors.New("LALOAD error")
			}
			var value = array[index]
			push(f, value)
			push(f, value) // pushed twice due to JDK longs being 64 bits wide

		case opcodes.FALOAD: //		0x30	(push contents of an float array element)
			index := pop(f).(int64)
			ref := pop(f) // ptr to array object
			// fAref := (*object.JacobinFloatArray)(ref)
			if ref == nil || ref == object.Null {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "FALOAD: Invalid (null) reference to an array"
				exceptions.Throw(exceptions.NullPointerException, errMsg)
				return errors.New(errMsg)
			}

			fAref := ref.(*object.Object)
			array := *(fAref.Fields[0].Fvalue).(*[]float64)
			if index >= int64(len(array)) {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "FALOAD: Invalid array subscript"
				exceptions.Throw(exceptions.ArrayIndexOutOfBoundsException, errMsg)
				return errors.New(errMsg)
			}
			var value = array[index]
			push(f, value)

		case opcodes.DALOAD: //		0x31	(push contents of a double array element)
			index := pop(f).(int64)
			fAref := pop(f).(*object.Object) // ptr to array object
			if fAref == nil {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "DALOAD: Invalid (null) reference to an array"
				exceptions.Throw(exceptions.NullPointerException, errMsg)
				return errors.New(errMsg)
			}
			array := *(fAref.Fields[0].Fvalue).(*[]float64)

			if index >= int64(len(array)) {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "DALOAD: Invalid array subscript"
				exceptions.Throw(exceptions.ArrayIndexOutOfBoundsException, errMsg)
				return errors.New(errMsg)
			}
			var value = array[index]
			push(f, value)
			push(f, value)
		case opcodes.AALOAD: // 0x32    (push contents of a reference array element)
			index := pop(f).(int64)
			rAref := pop(f) // the array object. Can't be cast to *Object b/c might be nil
			if rAref == nil {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "AALOAD: Invalid (null) reference to an array"
				exceptions.Throw(exceptions.NullPointerException, errMsg)
				return errors.New(errMsg)
			}

			arrayPtr := (rAref.(*object.Object)).Fields[0].Fvalue.(*[]*object.Object)
			size := int64(len(*arrayPtr))
			if index >= size {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "AALOAD: Invalid array subscript"
				exceptions.Throw(exceptions.ArrayIndexOutOfBoundsException, errMsg)
				return errors.New(errMsg)
			}
			array := *(arrayPtr)
			var value = array[index]
			push(f, value)

		case opcodes.BALOAD: // 0x33	(push contents of a byte/boolean array element)
			index := pop(f).(int64)
			ref := pop(f) // the array object
			if ref == nil || ref == object.Null {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "BALOAD: Invalid (null) reference to an array"
				exceptions.Throw(exceptions.InvalidTypeException, errMsg)
				return errors.New(errMsg)
			}

			var bAref *object.Object
			var arrayPtr *[]byte
			switch ref.(type) {
			case *object.Object:
				bAref = ref.(*object.Object)
				arrayPtr = bAref.Fields[0].Fvalue.(*[]byte)
			case *[]uint8:
				arrayPtr = ref.(*[]uint8)
			default:
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("BALOAD: Invalid type of object ref: %T", ref)
				exceptions.Throw(exceptions.InvalidTypeException, errMsg)
				return errors.New(errMsg)
			}
			size := int64(len(*arrayPtr))

			if index >= size {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "BALOAD: Invalid array subscript"
				exceptions.Throw(exceptions.ArrayIndexOutOfBoundsException, errMsg)
				return errors.New(errMsg)
			}
			array := *(arrayPtr)
			var value = array[index]
			push(f, int64(value))

		case opcodes.ISTORE, //  0x36 	(store popped top of stack int into local[index])
			opcodes.LSTORE: //  0x37 (store popped top of stack long into local[index])
			bytecode := f.Meth[f.PC]
			index := int(f.Meth[f.PC+1])
			f.PC += 1
			f.Locals[index] = pop(f).(int64)
			// longs and doubles are stored in localvar[x] and again in localvar[x+1]
			if bytecode == opcodes.LSTORE {
				f.Locals[index+1] = pop(f).(int64)
			}
		case opcodes.FSTORE: //  0x38 (store popped top of stack float into local[index])
			index := int(f.Meth[f.PC+1])
			f.PC += 1
			f.Locals[index] = pop(f).(float64)
		case opcodes.DSTORE: //  0x39 (store popped top of stack double into local[index])
			index := int(f.Meth[f.PC+1])
			f.PC += 1
			f.Locals[index] = pop(f).(float64)
			// longs and doubles are stored in localvar[x] and again in localvar[x+1]
			f.Locals[index+1] = pop(f).(float64)
		case opcodes.ASTORE: //  0x3A (store popped top of stack ref into localc[index])
			index := int(f.Meth[f.PC+1])
			f.PC += 1
			f.Locals[index] = pop(f)
		case opcodes.ISTORE_0: //   0x3B    (store popped top of stack int into local 0)
			f.Locals[0] = pop(f).(int64)
		case opcodes.ISTORE_1: //   0x3C   	(store popped top of stack int into local 1)
			f.Locals[1] = pop(f).(int64)
		case opcodes.ISTORE_2: //   0x3D   	(store popped top of stack int into local 2)
			f.Locals[2] = pop(f).(int64)
		case opcodes.ISTORE_3: //   0x3E    (store popped top of stack int into local 3)
			f.Locals[3] = pop(f).(int64)
		case opcodes.LSTORE_0: //   0x3F    (store long from top of stack into locals 0 and 1)
			var v = pop(f).(int64)
			f.Locals[0] = v
			f.Locals[1] = v
			pop(f)
		case opcodes.LSTORE_1: //   0x40    (store long from top of stack into locals 1 and 2)
			var v = pop(f).(int64)
			f.Locals[1] = v
			f.Locals[2] = v
			pop(f)
		case opcodes.LSTORE_2: //   0x41    (store long from top of stack into locals 2 and 3)
			var v = pop(f).(int64)
			f.Locals[2] = v
			f.Locals[3] = v
			pop(f)
		case opcodes.LSTORE_3: //   0x42    (store long from top of stack into locals 3 and 4)
			var v = pop(f).(int64)
			f.Locals[3] = v
			f.Locals[4] = v
			pop(f)
		case opcodes.FSTORE_0: // 0x43
			f.Locals[0] = pop(f).(float64)
		case opcodes.FSTORE_1: // 0x44
			f.Locals[1] = pop(f).(float64)
		case opcodes.FSTORE_2: // 0x45
			f.Locals[2] = pop(f).(float64)
		case opcodes.FSTORE_3: // 0x46
			f.Locals[3] = pop(f).(float64)
		case opcodes.DSTORE_0: // 0x47
			f.Locals[0] = pop(f).(float64)
			f.Locals[1] = pop(f).(float64)
		case opcodes.DSTORE_1: // 0x48
			f.Locals[1] = pop(f).(float64)
			f.Locals[2] = pop(f).(float64)
		case opcodes.DSTORE_2: // 0x49
			f.Locals[2] = pop(f).(float64)
			f.Locals[3] = pop(f).(float64)
		case opcodes.DSTORE_3: // 0x4A
			f.Locals[3] = pop(f).(float64)
			f.Locals[4] = pop(f).(float64)
		case opcodes.ASTORE_0: //	0x4B	(pop reference into local variable 0)
			f.Locals[0] = pop(f)
		case opcodes.ASTORE_1: //   0x4C	(pop reference into local variable 1)
			f.Locals[1] = pop(f)
		case opcodes.ASTORE_2: // 	0x4D	(pop reference into local variable 2)
			f.Locals[2] = pop(f)
		case opcodes.ASTORE_3: //	0x4E	(pop reference into local variable 3)
			f.Locals[3] = pop(f)
		case opcodes.IASTORE, //	0x4F	(store int in an array)
			opcodes.CASTORE, //		0x55 	(store char (2 bytes) in an array)
			opcodes.SASTORE: //    	0x56	(store a short in an array)
			value := pop(f).(int64)
			index := pop(f).(int64)
			arrObj := pop(f).(*object.Object) // the array object
			if arrObj == nil {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				exceptions.Throw(exceptions.NullPointerException,
					"IA/CA/SASTORE: Invalid (null) reference to an array")
				return errors.New("IA/CA/SASTORE: Invalid array address")
			}

			if arrObj.Fields[0].Ftype != "[I" {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("IA/CA/SASTORE: field type expected=[I, observed=%s", arrObj.Fields[0].Ftype)
				_ = log.Log(errMsg, log.SEVERE)
				exceptions.Throw(exceptions.ArrayStoreException, errMsg)
				return errors.New(errMsg)
			}

			array := *(arrObj.Fields[0].Fvalue).(*[]int64)
			size := int64(len(array))
			if index >= size {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("IA/CA/SASTORE: array size= %d but array index= %d (too large)", size, index)
				_ = log.Log(errMsg, log.SEVERE)
				exceptions.Throw(exceptions.ArrayIndexOutOfBoundsException, errMsg)
				return errors.New(errMsg)
			}
			array[index] = value

		case opcodes.LASTORE: // 0x50	(store a long in a long array)
			value := pop(f).(int64)
			pop(f) // second pop b/c longs use two slots
			index := pop(f).(int64)
			lAref := pop(f).(*object.Object) // ptr to array object
			if lAref == nil {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("LASTORE: Invalid (null) reference to an array")
				exceptions.Throw(exceptions.NullPointerException, errMsg)
				return errors.New(errMsg)
			}

			arrType := lAref.Fields[0].Ftype

			if arrType != "[I" {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("LASTORE: field type expected=[I, observed=%s", arrType)
				_ = log.Log(errMsg, log.SEVERE)
				exceptions.Throw(exceptions.ArrayStoreException,
					"LASTORE: Attempt to access array of incorrect type")
				return errors.New("LASTORE: Invalid array type")
			}

			array := *(lAref.Fields[0].Fvalue).(*[]int64)
			size := int64(len(array))
			if index >= size {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("LASTORE: array size=%d but index=%d (too large)", size, index)
				_ = log.Log(errMsg, log.SEVERE)
				exceptions.Throw(exceptions.ArrayIndexOutOfBoundsException,
					"LASTORE: Invalid array subscript")
				return errors.New("LASTORE: Invalid array index")
			}
			array[index] = value

		case opcodes.FASTORE: // 0x51	(store a float in a float array)
			value := pop(f).(float64)
			index := pop(f).(int64)
			fAref := pop(f).(*object.Object) // ptr to array object
			if fAref == nil {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				exceptions.Throw(exceptions.NullPointerException,
					"FASTORE: Invalid (null) reference to an array")
				return errors.New("FASTORE: Invalid array address")
			}

			if fAref.Fields[0].Ftype != "[F" {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("FASTORE: field type expected=[F, observed=%s", fAref.Fields[0].Ftype)
				_ = log.Log(errMsg, log.SEVERE)
				exceptions.Throw(exceptions.ArrayStoreException,
					"FASTORE: Attempt to access array of incorrect type")
				return errors.New("FASTORE: Invalid array type")
			}

			array := *(fAref.Fields[0].Fvalue).(*[]float64)
			size := int64(len(array))
			if index >= size {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("FASTORE: array size=%d but index=%d (too large)", size, index)
				_ = log.Log(errMsg, log.SEVERE)
				exceptions.Throw(exceptions.ArrayIndexOutOfBoundsException,
					"FASTORE: Invalid array subscript")
				return errors.New("FASTORE: Invalid array index")
			}
			array[index] = value

		case opcodes.DASTORE: // 0x52	(store a double in a doubles array)
			value := pop(f).(float64)
			pop(f) // second pop b/c doubles take two slots on the operand stack
			index := pop(f).(int64)
			dAref := pop(f).(*object.Object)
			if dAref == nil {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				exceptions.Throw(exceptions.NullPointerException,
					"DASTORE: Invalid (null) reference to an array")
				return errors.New("DASTORE: Invalid array reference")
			}

			if dAref.Fields[0].Ftype != "[F" {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("DASTORE: field type expected=[F, observed=%s", dAref.Fields[0].Ftype)
				_ = log.Log(errMsg, log.SEVERE)
				exceptions.Throw(exceptions.ArrayStoreException,
					"DASTORE: Attempt to access array of incorrect type")
				return errors.New("DASTORE: Invalid array type")
			}

			array := *(dAref.Fields[0].Fvalue).(*[]float64)
			size := int64(len(array))
			if index >= size {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("DASTORE: array size=%d but index=%d (too large)", size, index)
				_ = log.Log(errMsg, log.SEVERE)
				exceptions.Throw(exceptions.ArrayIndexOutOfBoundsException,
					"DASTORE: Invalid array subscript")
				return errors.New("DASTORE: Invalid array index")
			}

			array[index] = value

		case opcodes.AASTORE: // 0x53   (store a reference in a reference array)
			value := pop(f).(*object.Object)  // reference we're inserting
			index := pop(f).(int64)           // index into the array
			ptrObj := pop(f).(*object.Object) // ptr to the array object

			if ptrObj == nil {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				exceptions.Throw(exceptions.NullPointerException,
					"AASTORE: Invalid (null) reference to an array")
				return errors.New("AASTORE: Invalid array address")
			}

			if !strings.HasPrefix(ptrObj.Fields[0].Ftype, types.RefArray) {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("AASTORE: field type must start with '[L', got %s", ptrObj.Fields[0].Ftype)
				_ = log.Log(errMsg, log.SEVERE)
				exceptions.Throw(exceptions.ArrayStoreException,
					"AASTORE: Attempt to access array of incorrect type")
				return errors.New("AASTORE: Invalid array type")
			}

			// get pointer to the actual array
			arrayPtr := ptrObj.Fields[0].Fvalue.(*[]*object.Object)
			size := int64(len(*arrayPtr))
			if index >= size {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("AASTORE: array size=%d but index=%d (too large)", size, index)
				_ = log.Log(errMsg, log.SEVERE)
				exceptions.Throw(exceptions.ArrayIndexOutOfBoundsException,
					"AASTORE: Invalid array subscript")
				return errors.New("AASTORE: Invalid array index")
			}

			array := *arrayPtr
			array[index] = value

		case opcodes.BASTORE: // 0x54 	(store a boolean or byte in byte array)
			var value byte = 0
			rawValue := pop(f)
			value = convertInterfaceToByte(rawValue)
			index := pop(f).(int64)
			ptrObj := pop(f).(*object.Object) // ptr to array object
			if ptrObj == nil {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "BASTORE: Invalid (null) reference to an array"
				exceptions.Throw(exceptions.NullPointerException, errMsg)
				return errors.New(errMsg)
			}

			if ptrObj.Fields[0].Ftype != "[B" {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("BASTORE: Attempt to access array of incorrect type, expected=[B, observed=%s",
					ptrObj.Fields[0].Ftype)
				_ = log.Log(errMsg, log.SEVERE)
				exceptions.Throw(exceptions.ArrayStoreException, errMsg)
				return errors.New(errMsg)
			}

			// array := *(ptrObj.Fields[0].Fvalue.(*[]types.JavaByte)) // changed w/ JACOBIN-282
			array := *(ptrObj.Fields[0].Fvalue.(*[]byte))
			size := int64(len(array))
			if index >= size {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("BASTORE: Invalid array subscript: %d (size=%d) ", index, size)
				_ = log.Log(errMsg, log.SEVERE)
				exceptions.Throw(exceptions.ArrayIndexOutOfBoundsException, errMsg)
				return errors.New(errMsg)
			}
			array[index] = value

		case opcodes.POP: // 0x57 	(pop an item off the stack and discard it)
			if f.TOS < 0 {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				exceptions.FormatStackUnderflowError(f)
				break // the error will be picked up on the next instruction
			}
			f.TOS -= 1

		case opcodes.POP2: // 0x58	(pop 2 items from stack and discard them)
			if f.TOS < 1 {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				exceptions.FormatStackUnderflowError(f)
				break // the error will be picked up on the next instruction
			}
			f.TOS -= 2

		case opcodes.DUP: // 0x59 			(push an item equal to the current top of the stack
			tosItem := peek(f)
			if len(f.Meth) > 1 && f.Meth[f.PC+1] == opcodes.IMPDEP2 {
				break
			} // if invalid peek break now
			push(f, tosItem)
		case opcodes.DUP_X1: // 0x5A		(Duplicate the top stack value and insert two values down)
			top := pop(f)
			next := pop(f)
			push(f, top)
			push(f, next)
			push(f, top)
		case opcodes.DUP_X2: // 0x5B		(Duplicate top stack value and insert it three slots earlier)
			top := pop(f)
			next := pop(f)
			third := pop(f)
			push(f, top)
			push(f, third)
			push(f, next)
			push(f, top)
		case opcodes.DUP2: // 0x5C			(Duplicate the top two stack values)
			top := pop(f)
			next := peek(f)
			push(f, top)
			push(f, next)
			push(f, top)
		case opcodes.DUP2_X1: // 0x5D		(Duplicate the top two values, three slots down)
			top := pop(f)
			next := pop(f)
			third := pop(f)
			push(f, next) // so: top-next-third -> top-next-third->top->next
			push(f, top)
			push(f, third)
			push(f, next)
			push(f, top)
		case opcodes.DUP2_X2: // 0x5E		(Duplicate the top two values, four slots down)
			top := pop(f)
			next := pop(f)
			third := pop(f)
			fourth := pop(f)
			push(f, next) // so: top-next-third-fourth -> top-next-third-fourth-top-next
			push(f, top)
			push(f, fourth)
			push(f, third)
			push(f, next)
			push(f, top)
		case opcodes.SWAP: // 0x5F 	(swap top two items on stack)
			top := pop(f)
			next := pop(f)
			push(f, top)
			push(f, next)
		case opcodes.IADD: //  0x60		(add top 2 integers on operand stack, push result)
			i2 := pop(f).(int64)
			i1 := pop(f).(int64)
			sum := add(i1, i2)
			push(f, sum)
		case opcodes.LADD: //  0x61     (add top 2 longs on operand stack, push result)
			l2 := pop(f).(int64) //    longs occupy two slots, hence double pushes and pops
			pop(f)
			l1 := pop(f).(int64)
			pop(f)
			sum := add(l1, l2)
			push(f, sum)
			push(f, sum)
		case opcodes.FADD: // 0x62
			lhs := float32(pop(f).(float64))
			rhs := float32(pop(f).(float64))
			push(f, float64(lhs+rhs))
		case opcodes.DADD: // 0x63
			lhs := pop(f).(float64)
			pop(f)
			rhs := pop(f).(float64)
			pop(f)
			res := add(lhs, rhs)
			push(f, res)
			push(f, res)
		case opcodes.ISUB: //  0x64	(subtract top 2 integers on operand stack, push result)
			i2 := pop(f).(int64)
			i1 := pop(f).(int64)
			diff := subtract(i1, i2)
			push(f, diff)
		case opcodes.LSUB: //  0x65 (subtract top 2 longs on operand stack, push result)
			i2 := pop(f).(int64) //    longs occupy two slots, hence double pushes and pops
			pop(f)
			i1 := pop(f).(int64)
			pop(f)
			diff := subtract(i1, i2)

			push(f, diff)
			push(f, diff)
		case opcodes.FSUB: // 0x66
			i2 := float32(pop(f).(float64))
			i1 := float32(pop(f).(float64))
			push(f, float64(i1-i2))
		case opcodes.DSUB: // 0x67
			val2 := pop(f).(float64)
			pop(f)
			val1 := pop(f).(float64)
			pop(f)
			res := val1 - val2
			push(f, res)
			push(f, res)
		case opcodes.IMUL: //  0x68  	(multiply 2 integers on operand stack, push result)
			i2 := pop(f).(int64)
			i1 := pop(f).(int64)
			product := multiply(i1, i2)
			push(f, product)
		case opcodes.LMUL: //  0x69     (multiply 2 longs on operand stack, push result)
			l2 := pop(f).(int64) //    longs occupy two slots, hence double pushes and pops
			pop(f)
			l1 := pop(f).(int64)
			pop(f)
			product := multiply(l1, l2)
			push(f, product)
			push(f, product)
		case opcodes.FMUL: // 0x6A
			val1 := float32(pop(f).(float64))
			val2 := float32(pop(f).(float64))
			push(f, float64(val1*val2))
		case opcodes.DMUL: // 0x6B
			val1 := pop(f).(float64)
			pop(f)
			val2 := pop(f).(float64)
			pop(f)
			res := multiply(val1, val2)
			push(f, res)
			push(f, res)
		case opcodes.IDIV: //  0x6C (integer divide tos-1 by tos)
			val1 := pop(f).(int64)
			if val1 == 0 {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				exceptions.Throw(exceptions.ArithmeticException, ""+
					"IDIV: Arithmetic Exception: divide by zero")
				return errors.New("IDIV: Arithmetic Exception: divide by zero")
			} else {
				val2 := pop(f).(int64)
				push(f, val2/val1)
			}
		case opcodes.LDIV: //  0x6D   (long divide tos-2 by tos)
			val2 := pop(f).(int64)
			pop(f) //    longs occupy two slots, hence double pushes and pops
			if val2 == 0 {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "LDIV: Arithmetic Exception: divide by zero"
				exceptions.Throw(exceptions.ArithmeticException, errMsg)
				return errors.New(errMsg)
			} else {
				val1 := pop(f).(int64)
				pop(f)
				res := val1 / val2
				push(f, res)
				push(f, res)
			}

		case opcodes.FDIV: // 0x6E
			val1 := pop(f).(float64)
			val2 := pop(f).(float64)
			if val1 == 0.0 {
				if val2 == 0.0 {
					push(f, math.NaN())
				} else if math.Signbit(val1) { // this test for negative zero
					push(f, math.Inf(-1)) // but alas there is no -0 in golang (as of 1.20)
				} else {
					push(f, math.Inf(1))
				}
			} else {
				push(f, float64(float32(val2)/float32(val1)))
			}

		case opcodes.DDIV: // 0x6F
			val1 := pop(f).(float64)
			pop(f)
			val2 := pop(f).(float64)
			pop(f)
			if val1 == 0.0 {
				if val2 == 0.0 {
					push(f, math.NaN())
				} else if math.Signbit(val1) { // this tests for negative zero
					push(f, math.Inf(-1)) // but golang has no -0 as of v. 1.20
				} else {
					push(f, math.Inf(1))
				}
			} else {
				res := val2 / val1
				push(f, res)
				push(f, res)
			}
		case opcodes.IREM: // 	0x70	(remainder after int division, modulo)
			val2 := pop(f).(int64)
			if val2 == 0 {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "IREM: Arithmetic Exception: divide by zero"
				exceptions.Throw(exceptions.ArithmeticException, errMsg)
				return errors.New(errMsg)
			} else {
				val1 := pop(f).(int64)
				res := val1 % val2
				push(f, res)
			}
		case opcodes.LREM: // 	0x71	(remainder after long division)
			val2 := pop(f).(int64)
			pop(f) //    longs occupy two slots, hence double pushes and pops
			if val2 == 0 {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "LREM: Arithmetic Exception: divide by zero"
				exceptions.Throw(exceptions.ArithmeticException, errMsg)
				return errors.New(errMsg)
			} else {
				val1 := pop(f).(int64)
				pop(f)
				res := val1 % val2
				push(f, res)
				push(f, res)
			}
		case opcodes.FREM: // 0x72
			val2 := pop(f).(float64)
			val1 := pop(f).(float64)
			push(f, float64(float32(math.Remainder(val1, val2))))
		case opcodes.DREM: // 0x73
			val2 := pop(f).(float64)
			pop(f)
			val1 := pop(f).(float64)
			pop(f)
			drem := math.Remainder(val1, val2)
			push(f, drem)
			push(f, drem)
		case opcodes.INEG: //	0x74 	(negate an int)
			val := pop(f).(int64)
			push(f, -val)
		case opcodes.LNEG: //   0x75	(negate a long)
			val := pop(f).(int64)
			pop(f) // pop a second time because it's a long, which occupies 2 slots
			val = val * (-1)
			push(f, val)
			push(f, val)
		case opcodes.FNEG: //	0x76	(negate a float)
			val := pop(f).(float64)
			push(f, -val)
		case opcodes.DNEG: // 0x77
			pop(f)
			val := pop(f).(float64)
			push(f, -val)
			push(f, -val)
		case opcodes.ISHL: //	0x78 	(shift int left)
			shiftBy := pop(f).(int64)
			val1 := pop(f).(int64)
			var val2 int64
			if val1 < 0 { // if neg, shift as pos, then make neg
				val2 = (-val1) << (shiftBy & 0x1F) // only the bottom five bits are used
				push(f, -val2)
			} else {
				push(f, val1<<(shiftBy&0x1F))
			}
		case opcodes.LSHL: // 	0x79	(shift value1 (long) left by value2 (int) bits)
			shiftBy := pop(f).(int64)
			ushiftBy := uint64(shiftBy) & 0x3f // must be unsigned in golang; 0-63 bits per JVM
			val1 := pop(f).(int64)
			pop(f)
			val3 := val1 << ushiftBy
			push(f, val3)
			push(f, val3)
		case opcodes.ISHR: //  0x7A	(shift int value right)
			shiftBy := pop(f).(int64)
			val1 := pop(f).(int64)
			var val2 int64
			if val1 < 0 { // if neg, shift as pos, then make neg
				val2 = (-val1) >> (shiftBy & 0x1F) // only the bottom five bits are used
				push(f, -val2)
			} else {
				push(f, val1>>(shiftBy&0x1F))
			}
		case opcodes.LSHR, // 	0x7B	(shift value1 (long) right by value2 (int) bits)
			opcodes.LUSHR: // 	0x70
			shiftBy := pop(f).(int64)
			ushiftBy := uint64(shiftBy) & 0x3f // must be unsigned in golang; 0-63 bits per JVM
			val1 := pop(f).(int64)
			pop(f)
			val3 := val1 >> ushiftBy
			push(f, val3)
			push(f, val3)
		case opcodes.IUSHR: // 0x7C (unsigned shift right of int)
			shiftBy := pop(f).(int64) // TODO: verify the result against JDK
			val1 := pop(f).(int64)
			if val1 < 0 {
				val1 = -val1
			}
			push(f, val1>>(shiftBy&0x1F)) // only the bottom five bits are used
		case opcodes.IAND: //	0x7E	(logical and of two ints, push result)
			val1 := pop(f).(int64)
			val2 := pop(f).(int64)
			push(f, val1&val2)
		case opcodes.LAND: //   0x7F    (logical and of two longs, push result)
			val1 := pop(f).(int64)
			pop(f)
			val2 := pop(f).(int64)
			pop(f)
			val3 := val1 & val2
			push(f, val3)
			push(f, val3)
		case opcodes.IOR: // 0x 80 (logical OR of two ints, push result)
			val1 := pop(f).(int64)
			val2 := pop(f).(int64)
			push(f, val1|val2)
		case opcodes.LOR: // 0x81  (logical OR of two longs, push result)
			val1 := pop(f).(int64)
			pop(f)
			val2 := pop(f).(int64)
			pop(f)
			val3 := val1 | val2
			push(f, val3)
			push(f, val3)
		case opcodes.IXOR: // 	0x82	(logical XOR of two ints, push result)
			val1 := pop(f).(int64)
			val2 := pop(f).(int64)
			push(f, val1^val2)
		case opcodes.LXOR: // 	0x83  	(logical XOR of two longs, push result)
			val1 := pop(f).(int64)
			pop(f)
			val2 := pop(f).(int64)
			pop(f)
			val3 := val1 ^ val2
			push(f, val3)
			push(f, val3)
		case opcodes.IINC: // 	0x84    (increment local variable by a signed byte constant)
			localVarIndex := int64(f.Meth[f.PC+1])
			wbyte := f.Meth[f.PC+2]
			increment := byteToInt64(wbyte)
			orig := f.Locals[localVarIndex].(int64)
			f.Locals[localVarIndex] = orig + increment
			f.PC += 2
		case opcodes.I2F: //	0x86 	( convert int to float)
			intVal := pop(f).(int64)
			push(f, float64(intVal))
		case opcodes.I2L: // 	0x85     (convert int to long)
			// 	ints are already 64-bits, so this just pushes a second instance
			val := peek(f).(int64) // look without popping
			push(f, val)           // push the int a second time
		case opcodes.I2D: // 	0x87	(convert int to double)
			intVal := pop(f).(int64)
			dval := float64(intVal)
			push(f, dval) // doubles use two slots, hence two pushes
			push(f, dval)
		case opcodes.L2I: // 	0x88 	(convert long to int)
			longVal := pop(f).(int64)
			pop(f)
			intVal := longVal << 32 // remove high-end 4 bytes. this maintains the sign
			intVal >>= 32
			push(f, intVal)
		case opcodes.L2F: // 	0x89 	(convert long to float)
			longVal := pop(f).(int64)
			pop(f)
			float32Val := float32(longVal) //
			float64Val := float64(float32Val)
			push(f, float64Val) // floats tke up only 1 slot in the JVM
		case opcodes.L2D: // 	0x8A (convert long to double)
			longVal := pop(f).(int64)
			pop(f)
			dblVal := float64(longVal)
			push(f, dblVal)
			push(f, dblVal)
		case opcodes.D2I: // 0xBE
			pop(f)
			fallthrough
		case opcodes.F2I: // 0x8B
			floatVal := pop(f).(float64)
			push(f, int64(math.Trunc(floatVal)))
		case opcodes.F2D: // 0x8D
			floatVal := pop(f).(float64)
			push(f, floatVal)
			push(f, floatVal)
		case opcodes.D2L: // 	0x8F convert double to long
			pop(f)
			fallthrough
		case opcodes.F2L: // 	0x8C convert float to long
			floatVal := pop(f).(float64)
			truncated := int64(math.Trunc(floatVal))
			push(f, truncated)
			push(f, truncated)
		case opcodes.D2F: // 	0x90 Double to float
			floatVal := float32(pop(f).(float64))
			pop(f)
			push(f, float64(floatVal))
		case opcodes.I2B: //	0x91 convert into to byte preserving sign
			intVal := pop(f).(int64)
			byteVal := intVal & 0xFF
			if !(intVal > 0 && byteVal > 0) &&
				!(intVal < 0 && byteVal < 0) {
				byteVal = -byteVal
			}
			push(f, byteVal)
		case opcodes.I2C: //	0x92 convert to 16-bit char
			// determine what happens in Java if the int is negative
			intVal := pop(f).(int64)
			charVal := uint16(intVal) // Java chars are 16-bit unsigned values
			push(f, int64(charVal))
		case opcodes.I2S: //	0x93 convert int to short
			intVal := pop(f).(int64)
			shortVal := int16(intVal) // Java shorts are 16-bit signed values
			push(f, int64(shortVal))
		case opcodes.LCMP: // 	0x94 (compare two longs, push int -1, 0, or 1, depending on result)
			value2 := pop(f).(int64)
			pop(f)
			value1 := pop(f).(int64)
			pop(f)
			if value1 == value2 {
				push(f, int64(0))
			} else if value1 > value2 {
				push(f, int64(1))
			} else {
				push(f, int64(-1))
			}
		case opcodes.FCMPL, opcodes.FCMPG: // Ox95, 0x96 - float comparison - they differ only in NaN treatment
			value2 := pop(f).(float64)
			value1 := pop(f).(float64)
			if math.IsNaN(value1) || math.IsNaN(value2) {
				if f.Meth[f.PC] == opcodes.FCMPG {
					push(f, int64(1))
				} else {
					push(f, int64(-1))
				}
			} else if value1 > value2 {
				push(f, int64(1))
			} else if value1 < value2 {
				push(f, int64(-1))
			} else {
				push(f, int64(0))
			}
		case opcodes.DCMPL, opcodes.DCMPG: // 0x98, 0x97 - double comparison - they only differ in NaN treatment
			value2 := pop(f).(float64)
			pop(f)
			value1 := pop(f).(float64)
			pop(f)

			if math.IsNaN(value1) || math.IsNaN(value2) {
				if f.Meth[f.PC] == opcodes.DCMPG {
					push(f, int64(1))
				} else {
					push(f, int64(-1))
				}
			} else if value1 > value2 {
				push(f, int64(1))
			} else if value1 < value2 {
				push(f, int64(-1))
			} else {
				push(f, int64(0))
			}
		case opcodes.IFEQ: // 0x99 pop int, if it's == 0, go to the jump location
			// specified in the next two bytes
			popValue := pop(f)
			// bools are treated in the JVM as ints, so convert here if bool;
			// otherwise, values should be int64's
			var value int64
			switch popValue.(type) {
			case bool:
				if popValue == true {
					value = int64(1)
				} else {
					value = int64(0)
				}
			default:
				value = popValue.(int64)
			}
			if value == 0 {
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1
			} else {
				f.PC += 2
			}
		case opcodes.IFNE: // 0x9A pop int, it it's !=0, go to the jump location
			// specified in the next two bytes
			popValue := pop(f)
			// bools are treated in the JVM as ints, so convert here if bool;
			// otherwise, values should be int64's
			var value int64
			switch popValue.(type) {
			case bool:
				if popValue == true {
					value = int64(1)
				} else {
					value = int64(0)
				}
			default:
				value = popValue.(int64)
			}
			if value != 0 {
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1
			} else {
				f.PC += 2
			}
		case opcodes.IFLT: // 0x9B pop int, if it's < 0, go to the jump location
			// specified in the next two bytes
			value := pop(f).(int64)
			if value < 0 {
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1
			} else {
				f.PC += 2
			}
		case opcodes.IFGE: // 0x9C pop int, if it's >= 0, go to the jump location
			// specified in the next two bytes
			value := pop(f).(int64)
			if value >= 0 {
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1
			} else {
				f.PC += 2
			}
		case opcodes.IFGT: // 0x9D pop int, if it's > 0, go to the jump location
			// specified in the next two bytes
			value := pop(f).(int64)
			if value > 0 {
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1
			} else {
				f.PC += 2
			}
		case opcodes.IFLE: // 0x9E pop int, if it's <= 0, go to the jump location
			// specified in the next two bytes
			value := pop(f).(int64)
			if value <= 0 {
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1
			} else {
				f.PC += 2
			}
		case opcodes.IF_ICMPEQ: //  0x9F 	(jump if top two ints are equal)
			val2 := pop(f).(int64)
			val1 := pop(f).(int64)
			if int32(val1) == int32(val2) { // if comp succeeds, next 2 bytes hold instruction index
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1 // -1 b/c on the next iteration, pc is bumped by 1
			} else {
				f.PC += 2
			}
		case opcodes.IF_ICMPNE: //  0xA0    (jump if top two ints are not equal)
			val2 := pop(f).(int64)
			val1 := pop(f).(int64)
			if int32(val1) != int32(val2) { // if comp succeeds, next 2 bytes hold instruction index
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1 // -1 b/c on the next iteration, pc is bumped by 1
			} else {
				f.PC += 2
			}
		case opcodes.IF_ICMPLT: //  0xA1    (jump if popped val1 < popped val2)
			val2 := pop(f).(int64)
			val1 := pop(f).(int64)
			val1a := val1
			val2a := val2
			if val1a < val2a { // if comp succeeds, next 2 bytes hold instruction index
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1 // -1 b/c on the next iteration, pc is bumped by 1
			} else {
				f.PC += 2
			}
		case opcodes.IF_ICMPGE: //  0xA2    (jump if popped val1 >= popped val2)
			val2 := pop(f).(int64)
			val1 := pop(f).(int64)
			if val1 >= val2 { // if comp succeeds, next 2 bytes hold instruction index
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1 // -1 b/c on the next iteration, pc is bumped by 1
			} else {
				f.PC += 2
			}
		case opcodes.IF_ICMPGT: //  0xA3    (jump if popped val1 > popped val2)
			val2 := pop(f).(int64)
			val1 := pop(f).(int64)
			if int32(val1) > int32(val2) { // if comp succeeds, next 2 bytes hold instruction index
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1 // -1 b/c on the next iteration, pc is bumped by 1
			} else {
				f.PC += 2
			}
		case opcodes.IF_ICMPLE: //	0xA4	(jump if popped val1 <= popped val2)
			val2 := pop(f).(int64)
			val1 := pop(f).(int64)
			if val1 <= val2 { // if comp succeeds, next 2 bytes hold instruction index
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1 // -1 b/c on the next iteration, pc is bumped by 1
			} else {
				f.PC += 2
			}
		case opcodes.IF_ACMPEQ: // 0xA5		(jump if two addresses are equal)
			val2 := pop(f)
			val1 := pop(f)
			if val1 == val2 { // if comp succeeds, next 2 bytes hold instruction index
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1 // -1 b/c on the next iteration, pc is bumped by 1
			} else {
				f.PC += 2
			}
		case opcodes.IF_ACMPNE: // 0xA6		(jump if two addresses are note equal)
			val2 := pop(f)
			val1 := pop(f)
			if val1 != val2 { // if comp succeeds, next 2 bytes hold instruction index
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1 // -1 b/c on the next iteration, pc is bumped by 1
			} else {
				f.PC += 2
			}
		case opcodes.GOTO: // 0xA7     (goto an instruction)
			jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
			f.PC = f.PC + int(jumpTo) - 1 // -1 because this loop will increment f.PC by 1
		case opcodes.IRETURN: // 0xAC (return an int and exit current frame)
			valToReturn := pop(f)
			f = fs.Front().Next().Value.(*frames.Frame)
			push(f, valToReturn) // TODO: check what happens when main() ends on IRETURN
			return nil
		case opcodes.LRETURN: // 0xAD (return a long and exit current frame)
			valToReturn := pop(f).(int64)
			f = fs.Front().Next().Value.(*frames.Frame)
			push(f, valToReturn) // pushed twice b/c a long uses two slots
			push(f, valToReturn)
			return nil
		case opcodes.FRETURN: // 0xAE
			valToReturn := pop(f).(float64)
			f = fs.Front().Next().Value.(*frames.Frame)
			push(f, valToReturn)
			return nil
		case opcodes.DRETURN: // 0xAF (return a double and exit current frame)
			valToReturn := pop(f).(float64)
			f = fs.Front().Next().Value.(*frames.Frame)
			push(f, valToReturn) // pushed twice b/c a float uses two slots
			push(f, valToReturn)
			return nil
		case opcodes.ARETURN: // 0xB0	(return a reference)
			valToReturn := pop(f)
			// prevFrame := f
			f = fs.Front().Next().Value.(*frames.Frame)
			push(f, valToReturn)
			// fs.PushFront(prevFrame) //
			return nil
		case opcodes.RETURN: // 0xB1    (return from void function)
			f.TOS = -1 // empty the stack
			return nil
		case opcodes.GETSTATIC: // 0xB2		(get static field)
			CPslot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2]) // next 2 bytes point to CP entry
			f.PC += 2
			CP := f.CP.(*classloader.CPool)
			CPentry := CP.CpIndex[CPslot]
			if CPentry.Type != classloader.FieldRef { // the pointed-to CP entry must be a field reference
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("GETSTATIC: Expected a field ref, but got %d in"+
					"location %d in method %s of class %s\n",
					CPentry.Type, f.PC, f.MethName, f.ClName)
				_ = log.Log(errMsg, log.SEVERE)
				return errors.New(errMsg)
			}

			// get the field entry
			field := CP.FieldRefs[CPentry.Slot]

			// get the class entry from the field entry for this field. It's the class name.
			classRef := field.ClassIndex
			classNameIndex := CP.ClassRefs[CP.CpIndex[classRef].Slot]
			classNameEntry := CP.CpIndex[classNameIndex]
			className := CP.Utf8Refs[classNameEntry.Slot]

			// process the name and type entry for this field
			nAndTindex := field.NameAndType
			nAndTentry := CP.CpIndex[nAndTindex]
			nAndTslot := nAndTentry.Slot
			nAndT := CP.NameAndTypes[nAndTslot]
			fieldNameIndex := nAndT.NameIndex
			fieldName := classloader.FetchUTF8stringFromCPEntryNumber(CP, fieldNameIndex)
			fieldName = className + "." + fieldName

			// was this static field previously loaded? Is so, get its location and move on.
			prevLoaded, ok := statics.Statics[fieldName]
			if !ok { // if field is not already loaded, then
				// the class has not been instantiated, so
				// instantiate the class
				_, err := InstantiateClass(className, fs)
				if err == nil {
					prevLoaded, ok = statics.Statics[fieldName]
				} else {
					// glob := globals.GetGlobalRef()
					glob.ErrorGoStack = string(debug.Stack())
					errMsg := fmt.Sprintf("GETSTATIC: could not load class %s", className)
					_ = log.Log(errMsg, log.SEVERE)
					return errors.New(errMsg)
				}
			}

			// if the field can't be found even after instantiating the
			// containing class, something is wrong so get out of here.
			if !ok {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("GETSTATIC: could not find static field %s in class %s"+
					"\n", fieldName, className)
				_ = log.Log(errMsg, log.SEVERE)
				return errors.New(errMsg)
			}

			switch prevLoaded.Value.(type) {
			case bool:
				// a boolean, which might
				// be stored as a boolean, a byte (in an array), or int64
				// We want all forms normalized to int64
				value := prevLoaded.Value.(bool)
				prevLoaded.Value =
					types.ConvertGoBoolToJavaBool(value)
				push(f, prevLoaded.Value)
			case byte:
				value := prevLoaded.Value.(byte)
				prevLoaded.Value = int64(value)
				push(f, prevLoaded.Value)
			case int:
				value := prevLoaded.Value.(int)
				push(f, int64(value))
			default:
				push(f, prevLoaded.Value)
			}

			// doubles and longs consume two slots on the op stack
			// so push a second time
			if types.UsesTwoSlots(prevLoaded.Type) {
				push(f, prevLoaded.Value)
			}

		case opcodes.PUTSTATIC: // 0xB2		(get static field)
			CPslot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2]) // next 2 bytes point to CP entry
			f.PC += 2
			CP := f.CP.(*classloader.CPool)
			CPentry := CP.CpIndex[CPslot]
			if CPentry.Type != classloader.FieldRef { // the pointed-to CP entry must be a field reference
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("PUTSTATIC: Expected a field ref, but got %d in"+
					"location %d in method %s of class %s\n",
					CPentry.Type, f.PC, f.MethName, f.ClName)
				_ = log.Log(errMsg, log.SEVERE)
				return errors.New(errMsg)
			}

			// get the field entry
			field := CP.FieldRefs[CPentry.Slot]

			// get the class entry from the field entry for this field. It's the class name.
			classRef := field.ClassIndex
			classNameIndex := CP.ClassRefs[CP.CpIndex[classRef].Slot]
			classNameEntry := CP.CpIndex[classNameIndex]
			className := CP.Utf8Refs[classNameEntry.Slot]

			// process the name and type entry for this field
			nAndTindex := field.NameAndType
			nAndTentry := CP.CpIndex[nAndTindex]
			nAndTslot := nAndTentry.Slot
			nAndT := CP.NameAndTypes[nAndTslot]
			fieldNameIndex := nAndT.NameIndex
			fieldName := classloader.FetchUTF8stringFromCPEntryNumber(CP, fieldNameIndex)
			fieldName = className + "." + fieldName

			// was this static field previously loaded? Is so, get its location and move on.
			prevLoaded, ok := statics.Statics[fieldName]
			if !ok { // if field is not already loaded, then
				// the class has not been instantiated, so
				// instantiate the class
				_, err := InstantiateClass(className, fs)
				if err == nil {
					prevLoaded, ok = statics.Statics[fieldName]
				} else {
					// glob := globals.GetGlobalRef()
					glob.ErrorGoStack = string(debug.Stack())
					errMsg := fmt.Sprintf("PUTSTATIC: could not load class %s", className)
					_ = log.Log(errMsg, log.SEVERE)
					return errors.New(errMsg)
				}
			}

			// if the field can't be found even after instantiating the
			// containing class, something is wrong so get out of here.
			if !ok {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("PUTSTATIC: could not find static field %s", fieldName)
				_ = log.Log(errMsg, log.SEVERE)
				return errors.New(errMsg)
			}

			var value interface{}
			switch prevLoaded.Type {
			case types.Bool:
				// a boolean, which might
				// be stored as a boolean, a byte (in an array), or int64
				// We want all forms normalized to int64
				value = pop(f).(int64) & 0x01
				statics.Statics[fieldName] = statics.Static{
					Type:  prevLoaded.Type,
					Value: value,
				}
			case types.Char, types.Short, types.Int, types.Long:
				value = pop(f).(int64)
				statics.Statics[fieldName] = statics.Static{
					Type:  prevLoaded.Type,
					Value: value,
				}
			case types.Byte:
				var val byte
				v := pop(f)
				switch v.(type) { // could be passed a byte or an integral type for a value
				case int64:
					newVal := v.(int64)
					val = byte(newVal)
				case byte:
					val = v.(byte)
				}
				statics.Statics[fieldName] = statics.Static{
					Type:  prevLoaded.Type,
					Value: val,
				}
			case types.Float, types.Double:
				value = pop(f).(float64)
				statics.Statics[fieldName] = statics.Static{
					Type:  prevLoaded.Type,
					Value: value,
				}

			default:
				// if it's not a primitive or a pointer to a class,
				// then it should be a pointer to an object or to
				// a loaded class
				value = pop(f)
				switch value.(type) {
				case *object.Object:
					statics.Statics[fieldName] = statics.Static{
						Type:  prevLoaded.Type,
						Value: value,
					}
				case *classloader.Klass:
					// convert to an *object.Object
					kPtr := value.(*classloader.Klass)
					obj := object.MakeEmptyObject()
					obj.Klass = &kPtr.Data.Name
					objField := object.Field{
						Ftype:  "L" + kPtr.Data.Name + ";",
						Fvalue: kPtr,
					}
					obj.Fields = append(obj.Fields, objField)
					obj.FieldTable = nil

					statics.Statics[fieldName] = statics.Static{
						Type:  objField.Ftype,
						Value: value,
					}
				default:
					// glob := globals.GetGlobalRef()
					glob.ErrorGoStack = string(debug.Stack())
					errMsg := fmt.Sprintf("PUTSTATIC: field %s, type unrecognized: %v", fieldName, value)
					_ = log.Log(errMsg, log.SEVERE)
					return errors.New(errMsg)
				}
			}

			// doubles and longs consume two slots on the op stack
			// so push a second time
			if types.UsesTwoSlots(prevLoaded.Type) {
				pop(f)
			}

		case opcodes.GETFIELD: // 0xB4 get field in pointed-to-object
			CPslot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2]) // next 2 bytes point to CP entry
			f.PC += 2
			CP := f.CP.(*classloader.CPool)
			fieldEntry := CP.CpIndex[CPslot]
			if fieldEntry.Type != classloader.FieldRef { // the pointed-to CP entry must be a field reference
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("GETFIELD: Expected a field ref, but got %d in"+
					"location %d in method %s of class %s\n",
					fieldEntry.Type, f.PC, f.MethName, f.ClName)
				_ = log.Log(errMsg, log.SEVERE)
				return errors.New(errMsg)
			}

			// Get object reference from stack.
			ref := pop(f)
			switch ref.(type) {
			case *object.Object:
				break
			default:
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("GETFIELD: Invalid type of object ref: %T", ref)
				_ = log.Log(errMsg, log.SEVERE)
				return errors.New(errMsg)
			}

			// Extract field.
			obj := *ref.(*object.Object)
			var fieldType string
			var fieldValue interface{}

			if len(obj.FieldTable) < 1 {
				// Extract field from FieldTable map.
				fieldType = obj.Fields[fieldEntry.Slot].Ftype
				fieldValue = obj.Fields[fieldEntry.Slot].Fvalue
			} else {
				// Extract field from Fields slice.
				fullFieldEntry := CP.FieldRefs[fieldEntry.Slot]
				nameAndTypeCPIndex := fullFieldEntry.NameAndType
				nameAndTypeIndex := CP.CpIndex[nameAndTypeCPIndex]
				nameAndType := CP.NameAndTypes[nameAndTypeIndex.Slot]
				nameCPIndex := nameAndType.NameIndex
				nameCPentry := CP.CpIndex[nameCPIndex]
				fieldName := CP.Utf8Refs[nameCPentry.Slot]
				objField := obj.FieldTable[fieldName]
				fieldType = objField.Ftype
				fieldValue = objField.Fvalue // <<<< test for string and return pointer to String object
			}

			push(f, fieldValue)
			// fmt.Printf("DEBUG GETFIELD pushed type %s, value %v\n", fieldType, fieldValue)

			// doubles and longs consume two slots on the op stack
			// so push a second time
			if types.UsesTwoSlots(fieldType) {
				push(f, fieldValue)
			}

		case opcodes.PUTFIELD: // 0xB5 place value into an object's field
			CPslot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2]) // next 2 bytes point to CP entry
			f.PC += 2
			CP := f.CP.(*classloader.CPool)
			fieldEntry := CP.CpIndex[CPslot]
			if fieldEntry.Type != classloader.FieldRef { // the pointed-to CP entry must be a method reference
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("PUTFIELD: Expected a field ref, but got %d in"+
					"location %d in method %s of class %s\n",
					fieldEntry.Type, f.PC, f.MethName, f.ClName)
				_ = log.Log(errMsg, log.SEVERE)
				logTraceStack(f)
				return errors.New(errMsg)
			}

			var ref interface{} // pointer to object we're updating
			value := pop(f)     // the value we're placing in the field
			ref = pop(f)        // on non-long, non-double values, this will be a
			// reference to the object. On longs and doubles
			// it will be the second pop of the value field,
			// so we check for this.

			switch ref.(type) {
			case int64, float64: // if it is a float or double, then pop
				// once more to get the pointer to object. If it's an int64,
				// we know it's a long (likewise a float64 shows a double)
				// because that's the only reason a second pop would find
				// identical value types pushed twice. So pop once more to
				// get the object reference.
				ref = pop(f).(*object.Object)
			}

			obj := *(ref.(*object.Object))

			// if the value we're inserting is a reference to an
			// array object, we have to modify it to point directly
			// to the array of primitives, rather than to the array
			// object
			switch value.(type) {
			case *object.Object:
				if value != object.Null {
					v := *(value.(*object.Object))
					if obj.Fields != nil {
						if strings.HasPrefix(v.Fields[0].Ftype, types.Array) {
							value = v.Fields[0].Fvalue
						}
					}
				}
			}

			if obj.Fields != nil {
				// If it's a simple object w/out superclasses other than Object,
				// the fields in the object are numbered in the same
				// order they are declared in the constant pool. So,
				// to get to the right field, we only need to know
				// the slot number in CP.Fields. It will be the same
				// index into the object's fields.
				if strings.HasPrefix(obj.Fields[fieldEntry.Slot].Ftype, types.Static) {
					// glob := globals.GetGlobalRef()
					glob.ErrorGoStack = string(debug.Stack())
					errMsg := fmt.Sprintf("PUTFIELD: invalid attempt to update a static variable in %s.%s",
						f.ClName, f.MethName)
					_ = log.Log(errMsg, log.SEVERE)
					logTraceStack(f)
					return errors.New(errMsg)
				} else {
					obj.Fields[fieldEntry.Slot].Fvalue = value
				}
			} else {
				// otherwise, it's an object that contains superclass fields and
				// we need to access the field via the field table using the field name
				fullFieldEntry := CP.FieldRefs[fieldEntry.Slot]
				nameAndTypeCPIndex := fullFieldEntry.NameAndType
				nameAndTypeIndex := CP.CpIndex[nameAndTypeCPIndex]
				nameAndType := CP.NameAndTypes[nameAndTypeIndex.Slot]
				nameCPIndex := nameAndType.NameIndex
				nameCPentry := CP.CpIndex[nameCPIndex]
				fieldName := CP.Utf8Refs[nameCPentry.Slot]

				objField, ok := obj.FieldTable[fieldName]
				if !ok {
					errMsg := fmt.Sprintf("PUTFIELD: In trying for a superclass field, %s referenced by %s.%s is not present",
						fieldName, f.ClName, f.MethName)
					_ = log.Log(errMsg, log.SEVERE)
					logTraceStack(f)
					return errors.New(errMsg)
				}
				objField.Fvalue = value
				obj.FieldTable[fieldName] = objField
			}

		case opcodes.INVOKEVIRTUAL: // 	0xB6 invokevirtual (create new frame, invoke function)
			var err error
			CPslot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2]) // next 2 bytes point to CP entry
			f.PC += 2
			CP := f.CP.(*classloader.CPool)
			CPentry := CP.CpIndex[CPslot]
			if CPentry.Type != classloader.MethodRef { // the pointed-to CP entry must be a method reference
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("INVOKEVIRTUAL: Expected a method ref, but got %d in"+
					"location %d in method %s of class %s\n",
					CPentry.Type, f.PC, f.MethName, f.ClName)
				_ = log.Log(errMsg, log.SEVERE)
				return errors.New(errMsg)
			}

			// get the methodRef entry
			method := CP.MethodRefs[CPentry.Slot]

			// get the class entry from this method
			classRef := method.ClassIndex
			classNameIndex := CP.ClassRefs[CP.CpIndex[classRef].Slot]
			classNameEntry := CP.CpIndex[classNameIndex]
			className := CP.Utf8Refs[classNameEntry.Slot]

			// get the method name for this method
			nAndTindex := method.NameAndType
			nAndTentry := CP.CpIndex[nAndTindex]
			nAndTslot := nAndTentry.Slot
			nAndT := CP.NameAndTypes[nAndTslot]
			methodNameIndex := nAndT.NameIndex
			methodName := classloader.FetchUTF8stringFromCPEntryNumber(CP, methodNameIndex)

			// get the signature for this method
			methodSigIndex := nAndT.DescIndex
			methodType := classloader.FetchUTF8stringFromCPEntryNumber(CP, methodSigIndex)

			mtEntry := classloader.MTable[className+"."+methodName+methodType]
			if mtEntry.Meth == nil { // if the method is not in the method table, find it
				mtEntry, err = classloader.FetchMethodAndCP(className, methodName, methodType)
				if err != nil || mtEntry.Meth == nil {
					// TODO: search the classpath and retry
					// glob := globals.GetGlobalRef()
					glob.ErrorGoStack = string(debug.Stack())
					errMsg := "INVOKEVIRTUAL: Class method not found: " + className + "." + methodName
					_ = log.Log(errMsg, log.SEVERE)
					return errors.New(errMsg)
				}
			}

			// if we have a native function (here, one implemented in golang, rather than Java),
			// then follow the JVM spec and push the objectRef and the parameters to the function
			// as parameters. Consult:
			// https://docs.oracle.com/javase/specs/jvms/se17/html/jvms-6.html#jvms-6.5.invokevirtual
			if mtEntry.MType == 'G' { // so we have a native golang function
				// get the parameters/args off the stack
				gmethData := mtEntry.Meth.(gfunction.GMeth)
				paramCount := gmethData.ParamSlots
				var params []interface{}
				for i := 0; i < paramCount; i++ {
					params = append(params, pop(f))
				}

				// now get the objectRef (the object whose method we're invoking)
				objRef := pop(f).(*object.Object)
				params = append(params, objRef)

				_, err = runGmethod(mtEntry, fs, className, methodName, methodType, &params, true)
				if err != nil {
					// any exception message will already have been displayed to the user
					// glob := globals.GetGlobalRef()
					glob.ErrorGoStack = string(debug.Stack())
					errMsg := fmt.Sprintf("INVOKEVIRTUAL: Error encountered in: %s.%s"+
						className, methodName)
					return errors.New(errMsg)
				}
				break
			}

			if mtEntry.MType == 'J' { // it's a Java or Native function
				m := mtEntry.Meth.(classloader.JmEntry)
				if m.AccessFlags&0x0100 > 0 {
					// Native code
					// glob := globals.GetGlobalRef()
					glob.ErrorGoStack = string(debug.Stack())
					errMsg := "INVOKEVIRTUAL: Native method requested: " + className + "." + methodName
					_ = log.Log(errMsg, log.SEVERE)
					return errors.New(errMsg)
				}
				fram, err := createAndInitNewFrame(
					className, methodName, methodType, &m, true, f)
				if err != nil {
					// glob := globals.GetGlobalRef()
					glob.ErrorGoStack = string(debug.Stack())
					errMsg := "INVOKEVIRTUAL: Error creating frame in: " + className + "." + methodName
					return errors.New(errMsg)
				}
				f.PC += 1                            // move to next bytecode before exiting
				fs.PushFront(fram)                   // push the new frame
				f = fs.Front().Value.(*frames.Frame) // point f to the new head
				return runFrame(fs)
			}
		case opcodes.INVOKESPECIAL: //	0xB7 invokespecial (invoke constructors, private methods, etc.)
			CPslot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2]) // next 2 bytes point to CP entry
			f.PC += 2
			CP := f.CP.(*classloader.CPool)
			className, methName, methSig := classloader.GetMethInfoFromCPmethref(CP, CPslot)

			// if it's a call to java/lang/Object."<init>":()V, which happens frequently,
			// that function simply returns. So test for it here and if it is, skip the rest
			fullConstructorName := className + "." + methName + methSig
			if fullConstructorName == "java/lang/Object.<init>()V" {
				break
			}

			mtEntry, err := classloader.FetchMethodAndCP(className, methName, methSig)
			if err != nil || mtEntry.Meth == nil {
				// TODO: search the classpath and retry
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "INVOKESPECIAL: Class method not found: " + className + "." + methName
				_ = log.Log(errMsg, log.SEVERE)
				return errors.New(errMsg)
			}

			if mtEntry.MType == 'G' { // it's a golang method

				// get the parameters/args, if any, off the stack
				gmethData := mtEntry.Meth.(gfunction.GMeth)
				paramCount := gmethData.ParamSlots
				var params []interface{}
				for i := 0; i < paramCount; i++ {
					params = append(params, pop(f))
				}

				// now get the objectRef (the object whose method we're invoking)
				objRef := pop(f).(*object.Object)
				params = append(params, objRef)

				_, err = runGmethod(mtEntry, fs, className, methName, methSig, &params, true)
				if err != nil {
					// any exception message will already have been displayed to the user
					// glob := globals.GetGlobalRef()
					glob.ErrorGoStack = string(debug.Stack())
					errMsg := fmt.Sprintf("INVOKESPECIAL: Error encountered in: %s.%s", className, methName)
					return errors.New(errMsg)
				}
				break
			} else if mtEntry.MType == 'J' {
				// TODO: handle arguments to method, if any
				m := mtEntry.Meth.(classloader.JmEntry)
				fram, err := createAndInitNewFrame(className, methName, methSig, &m, true, f)
				if err != nil {
					// glob := globals.GetGlobalRef()
					glob.ErrorGoStack = string(debug.Stack())
					errMsg := "INVOKESPECIAL: Error creating frame in: " + className + "." + methName
					return errors.New(errMsg)
				}

				f.PC += 1
				fs.PushFront(fram)                   // push the new frame
				f = fs.Front().Value.(*frames.Frame) // point f to the new head
				return runFrame(fs)
			}
		case opcodes.INVOKESTATIC: // 	0xB8 invokestatic (create new frame, invoke static function)
			CPslot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2]) // next 2 bytes point to CP entry
			f.PC += 2
			CP := f.CP.(*classloader.CPool)
			CPentry := CP.CpIndex[CPslot]
			// get the methodRef entry
			method := CP.MethodRefs[CPentry.Slot]

			// get the class entry from this method
			classRef := method.ClassIndex
			classNameIndex := CP.ClassRefs[CP.CpIndex[classRef].Slot]
			classNameEntry := CP.CpIndex[classNameIndex]
			className := CP.Utf8Refs[classNameEntry.Slot]

			// get the method name for this method
			nAndTindex := method.NameAndType
			nAndTentry := CP.CpIndex[nAndTindex]
			nAndTslot := nAndTentry.Slot
			nAndT := CP.NameAndTypes[nAndTslot]
			methodNameIndex := nAndT.NameIndex
			methodName := classloader.FetchUTF8stringFromCPEntryNumber(CP, methodNameIndex)

			// get the signature for this method
			methodSigIndex := nAndT.DescIndex
			methodType := classloader.FetchUTF8stringFromCPEntryNumber(
				CP, methodSigIndex)

			mtEntry, err := classloader.FetchMethodAndCP(className, methodName, methodType)
			if err != nil || mtEntry.Meth == nil {
				// TODO: search the classpath and retry
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "INVOKESTATIC: Class method not found: " + className + "." + methodName
				_ = log.Log(errMsg, log.SEVERE)
				return errors.New(errMsg)
			}

			// before we can run the method, we need to either instantiate the class and/or
			// make sure that its static intializer block (if any) has been run. At this point,
			// all we know the class exists and has been loaded.
			k := classloader.MethAreaFetch(className)
			if k.Data.ClInit == types.ClInitNotRun {
				err = runInitializationBlock(k, nil, fs)
				if err != nil {
					// glob := globals.GetGlobalRef()
					glob.ErrorGoStack = string(debug.Stack())
					errMsg := fmt.Sprintf("INVOKESTATIC: error running initializer block in %s",
						className)
					_ = log.Log(errMsg, log.SEVERE)
					return errors.New(errMsg)
				}
			}

			if mtEntry.MType == 'G' {
				gmethData := mtEntry.Meth.(gfunction.GMeth)
				paramCount := gmethData.ParamSlots
				var params []interface{}
				for i := 0; i < paramCount; i++ {
					params = append(params, pop(f))
				}

				f, err = runGmethod(mtEntry, fs, className, methodName, methodType, &params, false)

				if err != nil {
					// any exceptions message will already have been displayed to the user
					// glob := globals.GetGlobalRef()
					glob.ErrorGoStack = string(debug.Stack())
					return errors.New("INVOKESTATIC: Error encountered in: " +
						className + "." + methodName)
				}
			} else if mtEntry.MType == 'J' {
				m := mtEntry.Meth.(classloader.JmEntry)
				fram, err := createAndInitNewFrame(
					className, methodName, methodType, &m, false, f)
				if err != nil {
					// glob := globals.GetGlobalRef()
					glob.ErrorGoStack = string(debug.Stack())
					return errors.New("INVOKESTATIC: Error creating frame in: " +
						className + "." + methodName)
				}

				f.PC += 1                            // point to the next bytecode before exiting
				fs.PushFront(fram)                   // push the new frame
				f = fs.Front().Value.(*frames.Frame) // point f to the new head
				return runFrame(fs)
			}

		case opcodes.NEW: // 0xBB 	new: create and instantiate a new object
			CPslot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2]) // next 2 bytes point to CP entry
			f.PC += 2
			CP := f.CP.(*classloader.CPool)
			CPentry := CP.CpIndex[CPslot]
			if CPentry.Type != classloader.ClassRef && CPentry.Type != classloader.Interface {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("NEW: Invalid type for new object")
				_ = log.Log(errMsg, log.SEVERE)
				return errors.New(errMsg)
			}

			// the classref points to a UTF8 record with the name of the class to instantiate
			var className string
			if CPentry.Type == classloader.ClassRef {
				utf8Index := CP.ClassRefs[CPentry.Slot]
				className = classloader.FetchUTF8stringFromCPEntryNumber(CP, utf8Index)
			}

			ref, err := InstantiateClass(className, fs)
			if err != nil {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("NEW: could not load class %s", className)
				_ = log.Log(errMsg, log.SEVERE)
				return errors.New(errMsg)
			}
			push(f, ref.(*object.Object))

		case opcodes.NEWARRAY: // 0xBC create a new array of primitives
			size := pop(f).(int64)
			if size < 0 {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "NEWARRAY: Invalid size for array"
				exceptions.Throw(exceptions.NegativeArraySizeException, errMsg)
				return errors.New(errMsg)
			}

			arrayType := int(f.Meth[f.PC+1])
			f.PC += 1

			actualType := object.JdkArrayTypeToJacobinType(arrayType)
			if actualType == object.ERROR {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "NEWARRAY: Invalid array type specified"
				_ = log.Log(errMsg, log.SEVERE)
				return errors.New(errMsg)
			}

			arrayPtr := object.Make1DimArray(uint8(actualType), size)
			g := globals.GetGlobalRef()
			g.ArrayAddressList.PushFront(arrayPtr)
			push(f, arrayPtr)

		case opcodes.ANEWARRAY: // 0xBD create array of references
			size := pop(f).(int64)
			if size < 0 {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "ANEWARRAY: Invalid size for array"
				exceptions.Throw(exceptions.NegativeArraySizeException, errMsg)
				return errors.New(errMsg)
			}

			refTypeSlot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2]) // next 2 bytes point to CP entry
			f.PC += 2
			CP := f.CP.(*classloader.CPool)
			refType := CP.CpIndex[refTypeSlot]
			if refType.Type != classloader.ClassRef && refType.Type != classloader.Interface {
				errMsg := fmt.Sprintf("ANEWARRAY: Presently works only with classes and interfaces")
				_ = log.Log(errMsg, log.SEVERE)
				return errors.New(errMsg)
			}

			var refTypeName = ""
			if refType.Type == classloader.ClassRef {
				utf8Index := CP.ClassRefs[refType.Slot]
				refTypeName = classloader.FetchUTF8stringFromCPEntryNumber(CP, utf8Index)
			}

			arrayPtr := object.Make1DimRefArray(&refTypeName, size)
			g := globals.GetGlobalRef()
			g.ArrayAddressList.PushFront(arrayPtr)
			push(f, arrayPtr)

			// The bytecode is followed by a two-byte index into the CP
			// which indicates what type the reference points to. We
			// don't presently check the type, so we skip over these
			// two bytes.
			// f.PC += 2

		case opcodes.ARRAYLENGTH: // OxBE get size of array
			// expects a pointer to an array
			ref := pop(f)
			if ref == nil {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "ARRAYLENGTH: Invalid (null) reference to an array"
				exceptions.Throw(exceptions.NullPointerException, errMsg)
				return errors.New(errMsg)
			}

			var size int64
			switch ref.(type) {
			// the type of array reference can vary. For many instances,
			// it will be a pointer to an array object. In other cases,
			// such as inside Java String class, the actual primitive
			// array of bytes will be extracted as a field and passed
			// to this function, so we need to accommodate all types--
			// hence, the switch on type.
			case *[]int8:
				array := *ref.(*[]int8)
				size = int64(len(array))
			case *[]uint8: // = go byte
				array := *ref.(*[]uint8)
				size = int64(len(array))
			case *object.Object:
				r := ref.(*object.Object)
				arrayType := r.Fields[0].Ftype
				switch arrayType {
				case types.ByteArray:
					arrayPtr := r.Fields[0].Fvalue.(*[]byte)
					size = int64(len(*arrayPtr))
				case types.RefArray:
					arrayPtr := r.Fields[0].Fvalue.(*[]*object.Object)
					size = int64(len(*arrayPtr))
				case types.FloatArray:
					arrayPtr := r.Fields[0].Fvalue.(*[]float64)
					size = int64(len(*arrayPtr))
				case types.IntArray:
					arrayPtr := r.Fields[0].Fvalue.(*[]int64)
					size = int64(len(*arrayPtr))
				default:
					arrayPtr := r.Fields[0].Fvalue.(*[]*object.Object)
					size = int64(len(*arrayPtr))
				}
			}
			push(f, size)
		case opcodes.ATHROW: // 0xBF throw an exception
			// objRef points to an instance of the error/exception class that's being thrown
			objectRef := pop(f).(*object.Object)
			if object.IsNull(objectRef) {
				errMsg := "ATHROW: Invalid (null) reference to a exception/error class to throw"
				exceptions.Throw(exceptions.NullPointerException, errMsg)
				return errors.New(errMsg)
			}

			// capture the golang stack
			stack := string(debug.Stack())
			// glob := globals.GetGlobalRef()
			glob.ErrorGoStack = stack

			// capture the JVM frame stack
			glob.JVMframeStack = exceptions.GrabFrameStack(fs)

			// get the name of the exception in the format used by HotSpot
			exceptionClass := *objectRef.Klass
			exceptionName := strings.Replace(exceptionClass, "/", ".", -1)

			// get the PC of the exception and check for any catch blocks
			exceptionPC := f.PC

			// find the frame with a valid catch block for this exception, if any
			catchFrame, handlerBytecode := exceptions.FindCatchFrame(fs, exceptionName, exceptionPC)
			// if there is no catch block, then print out the data we have (conforming
			// with whether we want the standard JDK info as elected with the -strictJDK
			// command-line option)
			if catchFrame == nil {
				// if the exception is not caught, then print the data from the stackTraceElements (STEs)
				// in the Throwable object or subclass (which is generally the specific exception class).

				// start by printing out the name of the exception/error and the thread it occurred on
				msg := ""
				if f.Thread == 1 { // if it's thread #1, use its name, "main"
					msg = fmt.Sprintf("Exception in thread \"main\" %s", exceptionName)
				} else {
					msg = fmt.Sprintf("Exception in thread %d %s", f.Thread, exceptionName)
				}
				_ = log.Log(msg, log.SEVERE)

				steArrayPtr := objectRef.FieldTable["stackTrace"].Fvalue.(*object.Object)
				rawSteArrayPtr := steArrayPtr.Fields[0].Fvalue.(*[]*object.Object) // *[]*object.Object (each of which is an STE)
				rawSteArray := *rawSteArrayPtr
				for i := 0; i < len(rawSteArray); i++ {
					ste := rawSteArray[i]
					methodName := ste.FieldTable["methodName"].Fvalue.(string)
					if methodName == "<init>" { // don't show constructors
						continue
					}
					rawClassName := ste.FieldTable["declaringClass"].Fvalue.(string)
					if rawClassName == "java/lang/Throwable" { // don't show Throwable methods
						continue
					}
					className := strings.Replace(rawClassName, "/", ".", -1)

					sourceLine := ste.FieldTable["sourceLine"].Fvalue.(string)

					var s string
					if sourceLine != "" {
						s = fmt.Sprintf("\tat %s.%s(%s:%s)", className,
							methodName, ste.FieldTable["fileName"].Fvalue, sourceLine)
					} else {
						s = fmt.Sprintf("\tat %s.%s(%s)", className,
							methodName, ste.FieldTable["fileName"].Fvalue)
					}
					_ = log.Log(s, log.SEVERE)
				}

				// show Jacobin's JVM stack info if -strictJDK is not set
				if glob.StrictJDK == false {
					_ = log.Log(" ", log.SEVERE)
					for _, frameData := range *glob.JVMframeStack {
						colon := strings.Index(frameData, ":")
						shortenedFrameData := frameData[colon+1:]
						_ = log.Log("\tat"+shortenedFrameData, log.SEVERE)
					}
				}

				// all exceptions that got this far are untrapped, so shutdown with an error code
				shutdown.Exit(shutdown.APP_EXCEPTION)

			} else { // perform the catch operation. We know the frame and the starting bytecode for the handler
				for fr := fs.Front(); fr != nil; fr = fr.Next() {
					var frm = fr.Value.(*frames.Frame)
					// f.ExceptionTable = &m.Exceptions
					if frm == catchFrame {
						frm.Meth = f.Meth[handlerBytecode:]
						frm.TOS = -1
						push(frm, objectRef)
						frm.PC = 0
						// make the frame with the catch block active
						fs.Front().Value = frm
						goto frameInterpreter
					}
				}
			}
		case opcodes.CHECKCAST: // 0xC0 same as INSTANCEOF but throws exception on null
			// because this uses the same logic as INSTANCEOF, any change here should
			// be made to INSTANCEOF
			ref := peek(f)
			if ref == nil { // if ref is nil, just carry on
				f.PC += 2 // move past two bytes pointing to comp object
				f.PC += 1
				continue
			}

			var obj *object.Object
			switch ref.(type) {
			case *object.Object:
				if ref == object.Null { // if ref is null, just carry on
					f.PC += 2 // move past two bytes pointing to comp object
					f.PC += 1
					continue
				} else {
					obj = (ref).(*object.Object)
				}
			default:
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "CHECKCAST: Invalid class reference"
				exceptions.Throw(exceptions.ClassCastException, errMsg)
				return errors.New(errMsg)
			}

			// at this point, we know we have a valid non-nil, non-null pointer to an object
			CPslot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2])
			f.PC += 2
			CP := f.CP.(*classloader.CPool)
			CPentry := CP.CpIndex[CPslot]
			if CPentry.Type == classloader.ClassRef { // slot of ClassRef points to
				// a CP entry for a UTF8 record w/ name of class
				var className string
				classNamePtr := classloader.FetchCPentry(CP, CPslot)
				if classNamePtr.RetType != classloader.IS_STRING_ADDR {
					// glob := globals.GetGlobalRef()
					glob.ErrorGoStack = string(debug.Stack())
					errMsg := "CHECKCAST: Invalid classRef found"
					_ = log.Log(errMsg, log.SEVERE)
					return errors.New(errMsg)
				}

				className = *(classNamePtr.StringVal)
				if MainThread.Trace {
					var traceInfo string
					if strings.HasPrefix(className, "[") {
						traceInfo = fmt.Sprintf("CHECKCAST: class is an array = %s", className)
					} else {
						traceInfo = fmt.Sprintf("CHECKCAST: className = %s", className)
					}
					_ = log.Log(traceInfo, log.TRACE_INST)
				}

				if strings.HasPrefix(className, "[") { // the object being checked is an array
					if obj.Klass != nil {
						sptr := obj.Klass
						// for the nonce if they're both the same type of arrays, we're good
						// TODO: if both are arrays of reference, check the leaf types
						if *sptr == className || strings.HasPrefix(className, *sptr) {
							break // exit this bytecode processing
						} else {
							/*** TODO: bypass this Throw action. Right thing to do?
							glob := globals.GetGlobalRef()
							glob.ErrorGoStack = string(debug.Stack())
							errMsg := fmt.Sprintf("CHECKCAST: %s is not castable with respect to %s", className, *sptr)
							exceptions.Throw(exceptions.ClassCastException, errMsg)
							return errors.New(errMsg)
							***/
							warnMsg := fmt.Sprintf("CHECKCAST: casting %s to %s might be unpleasant!", className, *sptr)
							_ = log.Log(warnMsg, log.WARNING)
						}
					} else {
						// glob := globals.GetGlobalRef()
						glob.ErrorGoStack = string(debug.Stack())
						errMsg := fmt.Sprintf("CHECKCAST: Klass field for object is nil")
						exceptions.Throw(exceptions.ClassCastException, errMsg)
						return errors.New(errMsg)
					}
				} else { // the object being checked is a class
					classPtr := classloader.MethAreaFetch(className)
					if classPtr == nil { // class wasn't loaded, so load it now
						if classloader.LoadClassFromNameOnly(className) != nil {
							// glob := globals.GetGlobalRef()
							glob.ErrorGoStack = string(debug.Stack())
							return errors.New("CHECKCAST: Could not load class: " + className)
						}
						classPtr = classloader.MethAreaFetch(className)
					}

					if classPtr != classloader.MethAreaFetch(*obj.Klass) {
						// glob := globals.GetGlobalRef()
						glob.ErrorGoStack = string(debug.Stack())
						errMsg := fmt.Sprintf("CHECKCAST: %s is not castable with respect to %s", className, classPtr.Data.Name)
						exceptions.Throw(exceptions.ClassCastException, errMsg)
						return errors.New(errMsg)
					}
					// note that if the classPtr == obj.Klass, which is the desired outcome,
					// do nothing. That is, the incoming stack should remain the same.
				}
			}

		case opcodes.INSTANCEOF: // 0xC1 validate the type of object (if not nil or null)
			// because this uses similar logic to CHECKCAST, any change here should
			// likely be made to CHECKCAST as well
			ref := pop(f)
			if ref == nil || ref == object.Null {
				push(f, int64(0))
				f.PC += 2 // move past index bytes to comp object
				break
			}

			switch ref.(type) {
			case *object.Object:
				if ref == object.Null {
					push(f, int64(0))
					f.PC += 2 // move past two bytes pointing to comp object
					break
				} else {
					obj := *ref.(*object.Object)
					CPslot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2])
					f.PC += 2
					CP := f.CP.(*classloader.CPool)
					CPentry := CP.CpIndex[CPslot]
					if CPentry.Type == classloader.ClassRef { // slot of ClassRef points to
						// a CP entry for a UTF8 record w/ name of class
						var className string
						classNamePtr := classloader.FetchCPentry(CP, CPslot)
						if classNamePtr.RetType != classloader.IS_STRING_ADDR {
							// glob := globals.GetGlobalRef()
							glob.ErrorGoStack = string(debug.Stack())
							errMsg := "INSTANCEOF: Invalid classRef found"
							_ = log.Log(errMsg, log.SEVERE)
							return errors.New(errMsg)
						} else {
							className = *(classNamePtr.StringVal)
							if MainThread.Trace {
								traceInfo := fmt.Sprintf("INSTANCEOF: className = %s", className)
								_ = log.Log(traceInfo, log.TRACE_INST)
							}
						}
						classPtr := classloader.MethAreaFetch(className)
						if classPtr == nil { // class wasn't loaded, so load it now
							if classloader.LoadClassFromNameOnly(className) != nil {
								// glob := globals.GetGlobalRef()
								glob.ErrorGoStack = string(debug.Stack())
								errMsg := "INSTANCEOF: Could not load class: " + className
								_ = log.Log(errMsg, log.SEVERE)
								return errors.New(errMsg)
							}
							classPtr = classloader.MethAreaFetch(className)
						}
						if classPtr == classloader.MethAreaFetch(*obj.Klass) {
							push(f, int64(1))
						} else {
							push(f, int64(0))
						}
					}
				}
			}

		case opcodes.MONITORENTER, opcodes.MONITOREXIT: // OxC2 and OxC3. These  are not implemented in the JDK JVM
			_ = pop(f) // so just pop off the reference on the stack

		case opcodes.MULTIANEWARRAY: // 0xC5 create multi-dimensional array
			var arrayDesc string
			var arrayType uint8

			// The first two chars after the bytecode point to a
			// classref entry in the CP. In turn, it points to a
			// string describing the array. Of the form [[L or
			// similar, in which one [ is present for every dimension
			// followed by a single letter describing the type of
			// entry in the leaf dimension of the array. The letters
			// are the usual ones used in the JVM for primitives, etc.
			// as in: https://docs.oracle.com/javase/specs/jvms/se17/html/jvms-4.html#jvms-4.3.2-200
			CPslot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2]) // next 2 bytes point to CP entry
			f.PC += 2
			CP := f.CP.(*classloader.CPool)
			CPentry := CP.CpIndex[CPslot]
			if CPentry.Type != classloader.ClassRef {
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				return errors.New("MULTIANEWARRAY: multi-dimensional array presently supports classes only")
			} else {
				utf8Index := CP.ClassRefs[CPentry.Slot]
				arrayDesc = classloader.FetchUTF8stringFromCPEntryNumber(CP, utf8Index)
			}

			var rawArrayType uint8
			for i := 0; i < len(arrayDesc); i++ {
				if arrayDesc[i] != '[' {
					rawArrayType = arrayDesc[i]
					break
				}
			}

			switch rawArrayType {
			case 'B', 'Z':
				arrayType = object.BYTE
			case 'F', 'D':
				arrayType = object.FLOAT
			case 'L':
				arrayType = object.REF
			default:
				arrayType = object.INT
			}

			// get the number of dimensions, then pop off the operand
			// stack an int for every dimension, giving the size of that
			// dimension and put them into a slice that starts with
			// the highest dimension first. So a two-dimensional array
			// such as x[4][3], would have entries of 4 and 3 respectively
			// in the dimsizes slice.
			dimensionCount := int(f.Meth[f.PC+1])
			f.PC += 1

			if dimensionCount > 3 { // TODO: explore arrays of > 5-255 dimensions
				// glob := globals.GetGlobalRef()
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := "MULTIANEWARRAY: Jacobin supports arrays only up to three dimensions"
				_ = log.Log(errMsg, log.SEVERE)
				return errors.New(errMsg)
			}

			dimSizes := make([]int64, dimensionCount)

			// the values on the operand stack give the last dimension
			// first when popped off the stack, so, they're stored here
			// in reverse order, so that dimSizes[0] will hold the first
			// dimenion.
			for i := dimensionCount - 1; i >= 0; i-- {
				dimSizes[i] = pop(f).(int64)
			}

			// A dimension of zero ends the dimensions, so we check
			// and cut off the dimensions below and includingthe 0-sized
			// one. Because this is almost certainly an error, we also
			// issue a warning.
			for i := range dimSizes {
				if dimSizes[i] == 0 {
					dimSizes = dimSizes[i+1:] // lop off the prev dims
					_ = log.Log("MULTIANEWARRAY: Multidimensional array with one dimension of size 0 encountered.",
						log.WARNING)
					break
				}
			}

			// Because of the possibility of a zero-sized dimension
			// affecting the valid number of dimensions, dimensionCount
			// can no longer be considered reliable. Use len(dimSizes).
			if len(dimSizes) == 3 {
				multiArr := object.Make1DimArray(object.REF, dimSizes[0])
				actualArray := *multiArr.Fields[0].Fvalue.(*[]*object.Object)
				for i := 0; i < len(actualArray); i++ {
					actualArray[i], _ = object.Make2DimArray(dimSizes[1],
						dimSizes[2], arrayType)
				}
				push(f, multiArr)
				break
			} else if len(dimSizes) == 2 { // 2-dim array is a special, trivial case
				multiArr, _ := object.Make2DimArray(dimSizes[0], dimSizes[1], arrayType)
				push(f, multiArr)
				break
				// It's possible due to a zero-length dimension, that we
				// need to create a single-dimension array.
			} else if len(dimSizes) == 1 {
				oneDimArr := object.Make1DimArray(arrayType, dimSizes[0])
				push(f, oneDimArr)
				break
			}

		case opcodes.IFNULL: // 0xC6 jump if TOS holds a null address
			// null = 0, so we duplicate logic of IFEQ instruction
			value := pop(f)
			if value == nil || value == object.Null {
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1
			} else {
				f.PC += 2
			}

		case opcodes.IFNONNULL: // 0xC7 jump if TOS does not hold a null address, where null = nil or object.Null
			value := pop(f)
			if value != nil { // it's not nil, but is it a null pointer?
				checkForPtr := value.(*object.Object)
				if checkForPtr == nil || checkForPtr == object.Null { // it really is a null pointer, so just move on
					f.PC += 2
				} else { // no, it's not nil nor a null pointer--so do the jump
					jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
					f.PC = f.PC + int(jumpTo) - 1
				}
			} else { // value is nil, so just move along
				f.PC += 2
			}

		case opcodes.IMPDEP2: // 0xFF private bytecode to flag an error. Next byte shows error type.
			// glob := globals.GetGlobalRef()
			glob.ErrorGoStack = string(debug.Stack())
			errCode := f.Meth[2]
			switch errCode {
			case 0x01: // stack overflow
				bytes := make([]byte, 2)
				bytes[0] = f.Meth[3]
				bytes[1] = f.Meth[4]
				location := int16(binary.BigEndian.Uint16(bytes))

				methName := fmt.Sprintf("%s.%s", f.ClName, f.MethName)
				rootCause := "stack overflow"
				exceptions.ShowPanicCause(rootCause)
				errMsg := fmt.Sprintf("Method: %-40s PC: %03d", methName, location)
				_ = log.Log(errMsg, log.SEVERE)

				fs.Remove(fs.Front()) // having reported on this frame's error, pop the frame off
				// return errors.New(rootCause)
				return errors.New(string(debug.Stack()))

			case 0x02: // stack underflow
				bytes := make([]byte, 2)
				bytes[0] = f.Meth[3]
				bytes[1] = f.Meth[4]
				location := int16(binary.BigEndian.Uint16(bytes))
				methName := fmt.Sprintf("%s.%s", f.ClName, f.MethName)
				rootCause := "stack underflow"
				exceptions.ShowPanicCause(rootCause)
				errMsg := fmt.Sprintf("Method: %-40s PC: %03d", methName, location)
				_ = log.Log(errMsg, log.SEVERE)

				fs.Remove(fs.Front()) // having reported on this frame's error, pop the frame off
				return errors.New(string(debug.Stack()))

			default:
				return errors.New("unknown error encountered")
			}

		default:
			missingOpCode := fmt.Sprintf("%d (0x%X)", f.Meth[f.PC], f.Meth[f.PC])

			if int(f.Meth[f.PC]) < len(opcodes.BytecodeNames) && int(f.Meth[f.PC]) > 0 {
				missingOpCode += fmt.Sprintf("(%s)", opcodes.BytecodeNames[f.Meth[f.PC]])
			}

			// glob := globals.GetGlobalRef()
			glob.ErrorGoStack = string(debug.Stack())
			errMsg := fmt.Sprintf("Invalid bytecode found: %s at location %d in method %s() of class %s\n",
				missingOpCode, f.PC, f.MethName, f.ClName)
			_ = log.Log(errMsg, log.SEVERE)
			return errors.New("invalid bytecode encountered")
		}
		f.PC += 1
	}
	return nil
}

func add[N frames.Number](num1, num2 N) N {
	return num1 + num2
}

// multiply two numbers
func multiply[N frames.Number](num1, num2 N) N {
	return num1 * num2
}

func subtract[N frames.Number](num1, num2 N) N {
	return num1 - num2
}

// create a new frame and load up the local variables with the passed
// arguments, set up the stack, and all the remaining items to begin execution
// Note: the includeObjectRef parameter is a boolean. When true, it indicates
// that in addition to the method parameter, an object reference is also on
// the stack and needs to be popped off the caller's opStack and passed in.
// (This would be the case for invokevirtual, among others.) When false, no
// object pointer is needed (for invokestatic, among others).
func createAndInitNewFrame(
	className string, methodName string, methodType string,
	m *classloader.JmEntry,
	includeObjectRef bool,
	currFrame *frames.Frame) (*frames.Frame, error) {

	if MainThread.Trace {
		traceInfo := fmt.Sprintf("\tcreateAndInitNewFrame: class=%s, meth=%s%s, includeObjectRef=%v, maxStack=%d, maxLocals=%d",
			className, methodName, methodType, includeObjectRef, m.MaxStack, m.MaxLocals)
		_ = log.Log(traceInfo, log.TRACE_INST)
	}

	f := currFrame

	stackSize := m.MaxStack
	if stackSize < 1 {
		stackSize = 2
	}
	// we increase the stack size by 2 because in many methods, the stack size must needs be
	// larger than what the class file specifies. For example, see the code in subbie() in main.class
	// of Jacobin test 389 (among many other examples). The value of 2 is chosen arbitrarily. It
	// might well be revised in future releases.
	stackSize += 2
	fram := frames.CreateFrame(stackSize)
	fram.Thread = currFrame.Thread
	fram.ClName = className
	fram.MethName = methodName
	fram.MethType = methodType
	fram.CP = m.Cp                           // add its pointer to the class CP
	fram.Meth = append(fram.Meth, m.Code...) // copy the method's bytecodes over

	// pop the parameters off the present stack and put them in
	// the new frame's locals. This is done in reverse order so
	// that the parameters are pushed in the right order to be
	// popped off by the receiving function
	var argList []interface{}
	paramsToPass :=
		util.ParseIncomingParamsFromMethTypeString(methodType)

	// primitives use a single byte/letter, but arrays can be many bytes:
	// a minimum of two (e.g., [I for array of ints). If the array
	// is multidimensional, the bytes will be [[I with one instance
	// of [ for every dimension. In the case of multidimensional
	// arrays, the arrays are always pushed as arrays of references,
	// and we simply mark off the number of [. For single-dimensional
	// arrays, we pass the kind of pointer that applies and mark off
	// a single instance of [
	for j := len(paramsToPass) - 1; j > -1; j-- {
		param := paramsToPass[j]
		primitive := param[0]

		arrayDimensions := 0
		if primitive == '[' {
			i := 0
			for i = 0; i < len(param); i++ {
				if param[i] == '[' {
					arrayDimensions += 1
				} else {
					break
				}
			}
			// param[i] now holds the primitive of the array
			primitive = param[i]
		}

		if arrayDimensions > 1 { // a multidimensional array
			// if the array is multidimensional, then we are
			// passing in an pointer to an array of references
			// to objects (lower arrays) regardless of the
			// lowest level of primitive in the array
			arg := pop(f).(*object.Object)
			argList = append(argList, arg)
			continue
		}

		if arrayDimensions == 1 { // a single-dimension array
			// a bunch of Java functions return raw arrays (like String.toCharArray()), which
			// are not really viewed by the JVM as objects in the full sense of the term. These
			// almost invariably are single-dimension arrays. So we test for these here and
			// return the corresponding object entity.
			value := pop(f)
			arg := object.MakeArrayFromRawArray(value)
			// arg := pop(f).(*object.Object)
			argList = append(argList, arg)
			continue
		}

		switch primitive { // it's not an array
		case 'D': // double
			arg := pop(f).(float64)
			argList = append(argList, arg)
			argList = append(argList, arg)
			pop(f)
		case 'F': // float
			arg := pop(f).(float64)
			argList = append(argList, arg)
		case 'B', 'C', 'I', 'S': // byte, char, integer, short
			arg := pop(f)
			switch arg.(type) {
			case int: // the arg should be int64, but is occasionally int. Tracking this down.
				arg = int64(arg.(int))
			}
			argList = append(argList, arg)
		case 'J': // long
			arg := pop(f).(int64)
			argList = append(argList, arg)
			argList = append(argList, arg)
			pop(f)
		case 'L': // pointer/reference
			// arg := pop(f).(*object.Object)
			arg := pop(f) // can't be case to *Object b/c it could be nil, which would panic
			argList = append(argList, arg)
		default:
			arg := pop(f)
			argList = append(argList, arg)
		}
	}

	// Initialize lenLocals = max (m.MaxLocals, len(argList)) but at least 1
	lenArgList := len(argList)
	lenLocals := m.MaxLocals
	if lenArgList > m.MaxLocals {
		lenLocals = lenArgList
	}
	if lenLocals < 1 {
		lenLocals = 1
	}

	// allocate the local variables
	for k := 0; k < lenLocals; k++ {
		fram.Locals = append(fram.Locals, int64(0))
	}

	// if includeObjectRef is true then objectRef != nil.
	// Insert it in the local[0]
	// This is used in invokevirtual, invokespecial, and invokeinterface.
	destLocal := 0
	if includeObjectRef {
		fram.Locals[0] = pop(f)
		fram.Locals = append(fram.Locals, int64(0)) // add the slot taken up by objectRef
		destLocal = 1                               // The first parameter starts at index 1
		lenLocals++                                 // There is 1 more local needed
	}

	if MainThread.Trace {
		traceInfo := fmt.Sprintf("\tcreateAndInitNewFrame: lenArgList=%d, lenLocals=%d, stackSize=%d",
			lenArgList, lenLocals, stackSize)
		_ = log.Log(traceInfo, log.TRACE_INST)
	}

	for j := lenArgList - 1; j >= 0; j-- {
		fram.Locals[destLocal] = argList[j]
		destLocal += 1
	}

	fram.TOS = -1

	return fram, nil
}
