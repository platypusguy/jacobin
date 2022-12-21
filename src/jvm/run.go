/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"container/list"
	"errors"
	"fmt"
	"jacobin/classloader"
	"jacobin/exceptions"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/shutdown"
	"jacobin/thread"
	"jacobin/util"
	"strconv"
	"unsafe"
)

var MainThread thread.ExecThread

// StartExec is where execution begins. It initializes various structures, such as
// the MTable, then using the passed-in name of the starting class, finds its main() method
// in the method area (it's guaranteed to already be loaded), grabs the executable
// bytes, creates a thread of execution, pushes the main() frame onto the JVM stack
// and begins execution.
func StartExec(className string, globals *globals.Globals) error {
	// initialize the MTable
	classloader.MTable = make(map[string]classloader.MTentry)
	classloader.MTableLoadNatives()

	me, err := classloader.FetchMethodAndCP(className, "main", "([Ljava/lang/String;)V")
	if err != nil {
		return errors.New("Class not found: " + className + ".main()")
	}

	m := me.Meth.(classloader.JmEntry)
	f := frames.CreateFrame(m.MaxStack) // create a new frame
	f.MethName = "main"
	f.ClName = className
	f.CP = m.Cp                        // add its pointer to the class CP
	for i := 0; i < len(m.Code); i++ { // copy the bytecodes over
		f.Meth = append(f.Meth, m.Code[i])
	}

	// allocate the local variables
	for k := 0; k < m.MaxLocals; k++ {
		f.Locals = append(f.Locals, 0)
	}

	// create the first thread and place its first frame on it
	MainThread = thread.CreateThread()
	MainThread.Stack = frames.CreateFrameStack()
	MainThread.ID = thread.AddThreadToTable(&MainThread, &globals.Threads)

	tracing := false
	trace, exists := globals.Options["-trace"]
	if exists {
		tracing = trace.Set
	}
	MainThread.Trace = tracing
	f.Thread = MainThread.ID

	if frames.PushFrame(MainThread.Stack, f) != nil {
		_ = log.Log("Memory exceptions allocating frame on thread: "+strconv.Itoa(MainThread.ID), log.SEVERE)
		return errors.New("outOfMemory Exception")
	}

	err = runThread(&MainThread)
	if err != nil {
		return err
	}
	return nil
}

// Point the thread to the top of the frame stack and tell it to run from there.
func runThread(t *thread.ExecThread) error {
	for t.Stack.Len() > 0 {
		err := runFrame(t.Stack)
		if err != nil {
			return err
		}

		if t.Stack.Len() == 1 { // true when the last executed frame was main()
			return nil
		}
	}
	return nil
}

// runFrame() is the principal execution function in Jacobin. It first tests for a
// golang function in the present frame. If it is a golang function, it's sent to
// a different function for execution. Otherwise, bytecode interpretation takes
// place through a giant switch statement.
func runFrame(fs *list.List) error {
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
			push(f, retval.(int64)) // if slotCount = 1

			if slotCount == 2 {
				push(f, retval.(int64)) // push a second time, if a long, double, etc.
			}
		}
		return err
	}

	// the frame's method is not a golang method, so it's Java bytecode, which
	// is interpreted in the rest of this function.
	for f.PC < len(f.Meth) {
		if MainThread.Trace {
			_ = log.Log("class: "+f.ClName+
				", meth: "+f.MethName+
				", pc: "+strconv.Itoa(f.PC)+
				", inst: "+BytecodeNames[int(f.Meth[f.PC])]+
				", tos: "+strconv.Itoa(f.TOS),
				log.TRACE_INST)
		}
		switch f.Meth[f.PC] { // cases listed in numerical value of opcode
		case NOP:
			break
		case ICONST_N1: //	0x02	(push -1 onto opStack)
			push(f, int64(-1))
		case ICONST_0: // 	0x03	(push int 0 onto opStack)
			push(f, int64(0))
		case ICONST_1: //  	0x04	(push int 1 onto opStack)
			push(f, int64(1))
		case ICONST_2: //   0x05	(push 2 onto opStack)
			push(f, int64(2))
		case ICONST_3: //   0x06	(push 3 onto opStack)
			push(f, int64(3))
		case ICONST_4: //   0x07	(push 4 onto opStack)
			push(f, int64(4))
		case ICONST_5: //   0x08	(push 5 onto opStack)
			push(f, int64(5))
		case LCONST_0: //   0x09    (push long 0 onto opStack)
			push(f, int64(0)) // b/c longs take two slots on the stack, it's pushed twice
			push(f, int64(0))
		case LCONST_1: //   0x0A    (push long 1 on to opStack)
			push(f, int64(1)) // b/c longs take two slots on the stack, it's pushed twice
			push(f, int64(1))
		case BIPUSH: //	0x10	(push the following byte as an int onto the stack)
			push(f, int64(f.Meth[f.PC+1]))
			f.PC += 1
		case LDC: // 	0x12   	(push constant from CP indexed by next byte)
			idx := f.Meth[f.PC+1]
			f.PC += 1

			CPe := FetchCPentry(f.CP, int(idx))
			if CPe.entryType != 0 && // 0 = error
				// Note: an invalid CP entry causes a java.lang.Verify error and
				//       is caught before execution of the program beings.
				// This instruction does not load longs or doubles
				CPe.entryType != classloader.DoubleConst &&
				CPe.entryType != classloader.LongConst { // if no error
				if CPe.retType == IS_INT64 {
					push(f, CPe.intVal)
				} else if CPe.retType == IS_FLOAT64 {
					push(f, CPe.floatVal)
				} else if CPe.retType == IS_STRUCT_ADDR || CPe.retType == IS_STRING_ADDR {
					push(f, int64(CPe.addrVal))
				}
			} else { // TODO: Determine what exception to throw
				exceptions.Throw(exceptions.InaccessibleObjectException, "Invalid type for LDC2_W instruction")
				shutdown.Exit(shutdown.APP_EXCEPTION)
			}

		case LDC_W: // 	0x13	(push constant from CP indexed by next two bytes)
			idx := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2])
			f.PC += 2

			CPe := FetchCPentry(f.CP, idx)
			if CPe.entryType != 0 && // this instruction does not load longs or doubles
				CPe.entryType != classloader.DoubleConst &&
				CPe.entryType != classloader.LongConst { // if no error
				if CPe.retType == IS_INT64 {
					push(f, CPe.intVal)
				} else if CPe.retType == IS_FLOAT64 {
					push(f, CPe.floatVal)
				} else {
					push(f, CPe.addrVal)
				}
			} else { // TODO: Determine what exception to throw
				exceptions.Throw(exceptions.InaccessibleObjectException, "Invalid type for LDC2_W instruction")
				shutdown.Exit(shutdown.APP_EXCEPTION)
			}
		case LDC2_W: // 0x14 	(push long or double from CP indexed by next two bytes)
			idx := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2])
			f.PC += 2

			CPe := FetchCPentry(f.CP, idx)
			if CPe.retType == IS_INT64 { // push value twice (due to 64-bit width)
				push(f, CPe.intVal)
				push(f, CPe.intVal)
			} else if CPe.retType == IS_FLOAT64 {
				push(f, CPe.floatVal)
				push(f, CPe.floatVal)
			} else { // TODO: Determine what exception to throw
				exceptions.Throw(exceptions.InaccessibleObjectException, "Invalid type for LDC2_W instruction")
				shutdown.Exit(shutdown.APP_EXCEPTION)
			}
		case ILOAD, // 0x15	(push int from local var, using next byte as index)
			FLOAD, //  0x17 (push float from local var, using next byte as index)
			ALOAD: //  0x19 (push ref from local var, using next byte as index)
			index := int(f.Meth[f.PC+1])
			f.PC += 1
			push(f, f.Locals[index])
		case LLOAD: // 0x16 (push long from local var, using next byte as index)
			index := int(f.Meth[f.PC+1])
			f.PC += 1
			val := f.Locals[index].(int64)
			push(f, val)
			push(f, val) // push twice due to item being 64 bits wide
		case DLOAD: // 0x18 (push double from local var, using next byte as index)
			index := int(f.Meth[f.PC+1])
			f.PC += 1
			val := f.Locals[index].(float64)
			push(f, val)
			push(f, val) // push twice due to item being 64 bits wide
		case ILOAD_0: // 	0x1A    (push local variable 0)
			push(f, f.Locals[0].(int64))
		case ILOAD_1: //    OX1B    (push local variable 1)
			push(f, f.Locals[1].(int64))
		case ILOAD_2: //    0X1C    (push local variable 2)
			push(f, f.Locals[2].(int64))
		case ILOAD_3: //  	0x1D   	(push local variable 3)
			push(f, f.Locals[3].(int64))
		// LLOAD use two slots, so the same value is pushed twice
		case LLOAD_0: //	0x1E	(push local variable 0, as long)
			push(f, f.Locals[0].(int64))
			push(f, f.Locals[0].(int64))
		case LLOAD_1: //	0x1F	(push local variable 1, as long)
			push(f, f.Locals[1].(int64))
			push(f, f.Locals[1].(int64))
		case LLOAD_2: //	0x20	(push local variable 2, as long)
			push(f, f.Locals[2].(int64))
			push(f, f.Locals[2].(int64))
		case LLOAD_3: //	0x21	(push local variable 3, as long)
			push(f, f.Locals[3].(int64))
			push(f, f.Locals[3].(int64))
		case DLOAD_0: //	0x26	(push local variable 0, as double)
			push(f, f.Locals[0])
			push(f, f.Locals[0])
		case DLOAD_1: //	0x27	(push local variable 1, as double)
			push(f, f.Locals[1])
			push(f, f.Locals[1])
		case DLOAD_2: //	0x28	(push local variable 2, as double)
			push(f, f.Locals[2])
			push(f, f.Locals[2])
		case DLOAD_3: //	0x29	(push local variable 3, as double)
			push(f, f.Locals[3])
			push(f, f.Locals[3])
		case ALOAD_0: //	0x2A	(push reference stored in local variable 0)
			push(f, f.Locals[0])
		case ALOAD_1: //	0x2B	(push reference stored in local variable 1)
			push(f, f.Locals[1])
		case ALOAD_2: //	0x2C    (push reference stored in local variable 2)
			push(f, f.Locals[2])
		case ALOAD_3: //	0x2D	(push reference stored in local variable 3)
			push(f, f.Locals[3])
		case ISTORE, //  0x36 	(store popped top of stack int into local[index])
			LSTORE, //  0x37 (store popped top of stack long into local[index])
			ASTORE: //  0x3A (store popped top of stack ref into localc[index])
			bytecode := f.Meth[f.PC]
			index := int(f.Meth[f.PC+1])
			f.PC += 1
			f.Locals[index] = pop(f).(int64)
			// longs and doubles are stored in localvar[x] and again in localvar[x+1]
			if bytecode == LSTORE {
				f.Locals[index+1] = pop(f).(int64)
			}
		case FSTORE, //  0x38 (store popped top of stack float into local[index])
			DSTORE: //  0x39 (store popped top of stack double into local[index])
			index := int(f.Meth[f.PC+1])
			f.PC += 1
			f.Locals[index] = pop(f).(float64)
			// longs and doubles are stored in localvar[x] and again in localvar[x+1]
			f.Locals[index+1] = pop(f).(float64)
		case ISTORE_0: //   0x3B    (store popped top of stack int into local 0)
			f.Locals[0] = pop(f).(int64)
		case ISTORE_1: //   0x3C   	(store popped top of stack int into local 1)
			f.Locals[1] = pop(f).(int64)
		case ISTORE_2: //   0x3D   	(store popped top of stack int into local 2)
			f.Locals[2] = pop(f).(int64)
		case ISTORE_3: //   0x3E    (store popped top of stack int into local 3)
			f.Locals[3] = pop(f).(int64)
		case LSTORE_0: //   0x3F    (store long from top of stack into locals 0 and 1)
			var v = pop(f).(int64)
			f.Locals[0] = v
			f.Locals[1] = v
			pop(f)
		case LSTORE_1: //   0x40    (store long from top of stack into locals 1 and 2)
			var v = pop(f).(int64)
			f.Locals[1] = v
			f.Locals[2] = v
			pop(f)
		case LSTORE_2: //   0x41    (store long from top of stack into locals 2 and 3)
			var v = pop(f).(int64)
			f.Locals[2] = v
			f.Locals[3] = v
			pop(f)
		case LSTORE_3: //   0x42    (store long from top of stack into locals 3 and 4)
			var v = pop(f).(int64)
			f.Locals[3] = v
			f.Locals[4] = v
			pop(f)
		case ASTORE_0: //	0x4B	(pop reference into local variable 0)
			f.Locals[0] = pop(f).(int64)
		case ASTORE_1: //   0x4C	(pop reference into local variable 1)
			f.Locals[1] = pop(f).(int64)
		case ASTORE_2: // 	0x4D	(pop reference into local variable 2)
			f.Locals[2] = pop(f).(int64)
		case ASTORE_3: //	0x4E	(pop reference into local variable 3)
			f.Locals[3] = pop(f).(int64)
		case DUP: // 0x59 			(push an item equal to the current top of the stack
			push(f, peek(f))
		case DUP_X1: // 0x5A		(Duplicate the top stack value and insert two values down)
			top := pop(f)
			next := pop(f)
			push(f, top)
			push(f, next)
			push(f, top)
		case SWAP: // 0x5C 	(swap top two items on stack)
			top := pop(f)
			next := pop(f)
			push(f, top)
			push(f, next)
		case IADD: //  0x60		(add top 2 integers on operand stack, push result)
			i2 := pop(f).(int64)
			i1 := pop(f).(int64)
			sum := add(i1, i2)
			push(f, sum)
		case LADD: //  0x61     (add top 2 longs on operand stack, push result)
			l2 := pop(f).(int64) //    longs occupy two slots, hence double pushes and pops
			pop(f)
			l1 := pop(f).(int64)
			pop(f)
			sum := add(l1, l2)
			push(f, sum)
			push(f, sum)
		case ISUB: //  0x64	(subtract top 2 integers on operand stack, push result)
			i2 := pop(f).(int64)
			i1 := pop(f).(int64)
			diff := subtract(i1, i2)
			push(f, diff)
		case LSUB: //  0x65 (subtract top 2 longs on operand stack, push result)
			i2 := pop(f).(int64) //    longs occupy two slots, hence double pushes and pops
			pop(f)
			i1 := pop(f).(int64)
			pop(f)
			diff := subtract(i1, i2)

			push(f, diff)
			push(f, diff)
		case IMUL: //  0x68  	(multiply 2 integers on operand stack, push result)
			i2 := pop(f).(int64)
			i1 := pop(f).(int64)
			product := multiply(i1, i2)

			push(f, product)
		case LMUL: //  0x69     (multiply 2 longs on operand stack, push result)
			l2 := pop(f).(int64) //    longs occupy two slots, hence double pushes and pops
			pop(f)
			l1 := pop(f).(int64)
			pop(f)
			product := multiply(l1, l2)

			push(f, product)
			push(f, product)
		case IDIV: //  0x6C (integer divide tos-1 by tos)
			val1 := pop(f).(int64)
			if val1 == 0 {
				exceptions.Throw(exceptions.ArithmeticException, "Arithmetic Exception: divide by zero")
				shutdown.Exit(shutdown.APP_EXCEPTION)
			} else {
				val2 := pop(f).(int64)
				push(f, val2/val1)
			}
		case LDIV: //  0x6D   (long divide tos-2 by tos)
			val2 := pop(f).(int64)
			pop(f) //    longs occupy two slots, hence double pushes and pops
			if val2 == 0 {
				exceptions.Throw(exceptions.ArithmeticException, "Arithmetic Exception: divide by zero")
				shutdown.Exit(shutdown.APP_EXCEPTION)
			} else {
				val1 := pop(f).(int64)
				pop(f)
				res := val1 / val2
				push(f, res)
				push(f, res)
			}
		case LREM: // 	0x71	(remainder after long division)
			val2 := pop(f).(int64)
			pop(f) //    longs occupy two slots, hence double pushes and pops
			if val2 == 0 {
				exceptions.Throw(exceptions.ArithmeticException, "Arithmetic Exception: divide by zero")
				shutdown.Exit(shutdown.APP_EXCEPTION)
			} else {
				val1 := pop(f).(int64)
				pop(f)
				res := val1 % val2
				push(f, res)
				push(f, res)
			}
		case INEG: //	0x74 	(negate an int)
			val := pop(f).(int64)
			push(f, -val)
		case LNEG: //   0x75	(negate a long)
			val := pop(f).(int64)
			pop(f) // pop a second time because it's a long, which occupies 2 slots
			val = val * (-1)
			push(f, val)
			push(f, val)
		// case FNEG: //	0x76	(negate a float)
		// 	val := float64(pop(f))
		// 	val = val * (-1.0)
		// 	push(f, val)
		case LSHL: // 	0x79	(shift value1 (long) left by value2 (int) bits)
			shiftBy := pop(f).(int64)
			ushiftBy := uint64(shiftBy) & 0x3f // must be unsigned in golang; 0-63 bits per JVM
			val1 := pop(f).(int64)
			pop(f)
			val3 := val1 << ushiftBy
			push(f, val3)
			push(f, val3)
		case LSHR, // 	0x7B	(shift value1 (long) right by value2 (int) bits)
			LUSHR: // 	0x70
			shiftBy := pop(f).(int64)
			ushiftBy := uint64(shiftBy) & 0x3f // must be unsigned in golang; 0-63 bits per JVM
			val1 := pop(f).(int64)
			pop(f)
			val3 := val1 >> ushiftBy
			push(f, val3)
			push(f, val3)
		case LAND: //   0x7F    (logical and of two longs, push result)
			val1 := pop(f).(int64)
			pop(f)
			val2 := pop(f).(int64)
			pop(f)
			val3 := val1 & val2
			push(f, val3)
			push(f, val3)
		case LOR: // 0x81  (logical OR of two longs, push result)
			val1 := pop(f).(int64)
			pop(f)
			val2 := pop(f).(int64)
			pop(f)
			val3 := val1 | val2
			push(f, val3)
			push(f, val3)
		case LXOR: // 0x83  (logical XOR of two longs, push result)
			val1 := pop(f).(int64)
			pop(f)
			val2 := pop(f).(int64)
			pop(f)
			val3 := val1 ^ val2
			push(f, val3)
			push(f, val3)
		case IINC: // 	0x84    (increment local variable by a constant)
			localVarIndex := int64(f.Meth[f.PC+1])
			constAmount := int64(f.Meth[f.PC+2])
			f.PC += 2
			orig := f.Locals[localVarIndex].(int64)
			f.Locals[localVarIndex] = orig + constAmount
		case I2L: // 	0x85     (convert int to long)
			val := pop(f).(int64)
			push(f, val) // all ints are 64-bits wide so already longs,
			push(f, val) // so push the int a second time
		case L2I: // 	0x88 	(convert long to int)
			longVal := pop(f).(int64)
			pop(f)
			intVal := longVal << 32 // remove high-end 4 bytes. this maintains the sign
			intVal >>= 32
			push(f, intVal)
		case L2F: // 	0x89 	(convert long to float)
			longVal := pop(f).(int64)
			pop(f)
			float32Val := float32(longVal) //
			float64Val := float64(float32Val)
			push(f, float64Val) // floats tke up only 1 slot in the JVM
		case L2D: // 	0x8A (convert long to double)
			longVal := pop(f).(int64)
			pop(f)
			dblVal := float64(longVal)
			push(f, dblVal)
			push(f, dblVal)
		case LCMP: // 0x94 (compare two longs, push int -1, 0, or 1, depending on result)
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
		case IFEQ: // 0x99 pop int, if it's == 0, go to the jump location
			// specified in the next two bytes
			value := pop(f).(int64)
			if value == 0 {
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1
			} else {
				f.PC += 2
			}
		case IFNE: // 0x9A pop int, it it's !=0, go to the jump location
			// specified in the next two bytes
			value := pop(f).(int64)
			if value != 0 {
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1
			} else {
				f.PC += 2
			}
		case IFLT: // 0x9B pop int, if it's < 0, go to the jump location
			// specified in the next two bytes
			value := pop(f).(int64)
			if value < 0 {
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1
			} else {
				f.PC += 2
			}
		case IFGE: // 0x9C pop int, if it's >= 0, go to the jump location
			// specified in the next two bytes
			value := pop(f).(int64)
			if value >= 0 {
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1
			} else {
				f.PC += 2
			}
		case IFGT: // 0x9D pop int, if it's > 0, go to the jump location
			// specified in the next two bytes
			value := pop(f).(int64)
			if value > 0 {
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1
			} else {
				f.PC += 2
			}
		case IFLE: // 0x9E pop int, if it's <= 0, go to the jump location
			// specified in the next two bytes
			value := pop(f).(int64)
			if value <= 0 {
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1
			} else {
				f.PC += 2
			}
		case IF_ICMPEQ: //  0x9F 	(jump if top two ints are equal)
			val2 := pop(f).(int64)
			val1 := pop(f).(int64)
			if int32(val1) == int32(val2) { // if comp succeeds, next 2 bytes hold instruction index
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1 // -1 b/c on the next iteration, pc is bumped by 1
			} else {
				f.PC += 2
			}
		case IF_ICMPNE: //  0xA0    (jump if top two ints are not equal)
			val2 := pop(f).(int64)
			val1 := pop(f).(int64)
			if int32(val1) != int32(val2) { // if comp succeeds, next 2 bytes hold instruction index
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1 // -1 b/c on the next iteration, pc is bumped by 1
			} else {
				f.PC += 2
			}
		case IF_ICMPLT: //  0xA1    (jump if popped val1 < popped val2)
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
		case IF_ICMPGE: //  0xA2    (jump if popped val1 >= popped val2)
			val2 := pop(f).(int64)
			val1 := pop(f).(int64)
			if val1 >= val2 { // if comp succeeds, next 2 bytes hold instruction index
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1 // -1 b/c on the next iteration, pc is bumped by 1
			} else {
				f.PC += 2
			}
		case IF_ICMPGT: //  0xA3    (jump if popped val1 > popped val2)
			val2 := pop(f).(int64)
			val1 := pop(f).(int64)
			if int32(val1) > int32(val2) { // if comp succeeds, next 2 bytes hold instruction index
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1 // -1 b/c on the next iteration, pc is bumped by 1
			} else {
				f.PC += 2
			}
		case IF_ICMPLE: //	0xA4	(jump if popped val1 <= popped val2)
			val2 := pop(f).(int64)
			val1 := pop(f).(int64)
			if val1 <= val2 { // if comp succeeds, next 2 bytes hold instruction index
				jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
				f.PC = f.PC + int(jumpTo) - 1 // -1 b/c on the next iteration, pc is bumped by 1
			} else {
				f.PC += 2
			}
		case GOTO: // 0xA7     (goto an instruction)
			jumpTo := (int16(f.Meth[f.PC+1]) * 256) + int16(f.Meth[f.PC+2])
			f.PC = f.PC + int(jumpTo) - 1 // -1 because this loop will increment f.PC by 1
		case IRETURN: // 0xAC (return an int and exit current frame)
			valToReturn := pop(f)
			f = fs.Front().Next().Value.(*frames.Frame)
			push(f, valToReturn) // TODO: check what happens when main() ends on IRETURN
			return nil
		case LRETURN: // 0xAD (return a long and exit current frame)
			valToReturn := pop(f).(int64)
			f = fs.Front().Next().Value.(*frames.Frame)
			push(f, valToReturn) // pushed twice b/c a long uses two slots
			push(f, valToReturn)
			return nil
		case DRETURN: // 0xAF (return a double and exit current frame)
			valToReturn := pop(f).(float64)
			f = fs.Front().Next().Value.(*frames.Frame)
			push(f, valToReturn) // pushed twice b/c a float uses two slots
			push(f, valToReturn)
			return nil
		case RETURN: // 0xB1    (return from void function)
			f.TOS = -1 // empty the stack
			return nil
		case GETSTATIC: // 0xB2		(get static field)
			// TODO: getstatic will instantiate a static class if it's not already instantiated
			// that logic has not yet been implemented and the code here is simply a reasonable
			// placeholder, which consists of creating a struct that holds most of the needed info
			// puts it into a slice of such static fields and pushes the index of this item in the slice
			// onto the stack of the frame.
			CPslot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2]) // next 2 bytes point to CP entry
			f.PC += 2
			CPentry := f.CP.CpIndex[CPslot]
			if CPentry.Type != classloader.FieldRef { // the pointed-to CP entry must be a field reference
				return fmt.Errorf("Expected a field ref on getstatic, but got %d in"+
					"location %d in method %s of class %s\n",
					CPentry.Type, f.PC, f.MethName, f.ClName)
			}

			// get the field entry
			field := f.CP.FieldRefs[CPentry.Slot]

			// get the class entry from the field entry for this field. It's the class name.
			classRef := field.ClassIndex
			classNameIndex := f.CP.ClassRefs[f.CP.CpIndex[classRef].Slot]
			classNameEntry := f.CP.CpIndex[classNameIndex]
			className := f.CP.Utf8Refs[classNameEntry.Slot]
			// println("Field name: " + className)

			// process the name and type entry for this field
			nAndTindex := field.NameAndType
			nAndTentry := f.CP.CpIndex[nAndTindex]
			nAndTslot := nAndTentry.Slot
			nAndT := f.CP.NameAndTypes[nAndTslot]
			fieldNameIndex := nAndT.NameIndex
			fieldName := classloader.FetchUTF8stringFromCPEntryNumber(f.CP, fieldNameIndex)
			fieldName = className + "." + fieldName

			// was this static field previously loaded? Is so, get its location and move on.
			prevLoaded, ok := classloader.Statics[fieldName]
			if ok { // if preloaded, then push the index into the array of constant fields
				push(f, prevLoaded)
				break
			}

			fieldTypeIndex := nAndT.DescIndex
			fieldType := classloader.FetchUTF8stringFromCPEntryNumber(f.CP, fieldTypeIndex)
			// println("full field name: " + fieldName + ", type: " + fieldType)
			newStatic := classloader.Static{
				Class:     'L',
				Type:      fieldType,
				ValueRef:  "",
				ValueInt:  0,
				ValueFP:   0,
				ValueStr:  "",
				ValueFunc: nil,
				CP:        f.CP,
			}
			classloader.StaticsArray = append(classloader.StaticsArray, newStatic)
			classloader.Statics[fieldName] = int64(len(classloader.StaticsArray) - 1)

			// push the pointer to the stack of the frame
			push(f, int64(len(classloader.StaticsArray)-1))

		case INVOKEVIRTUAL: // 	0xB6 invokevirtual (create new frame, invoke function)
			CPslot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2]) // next 2 bytes point to CP entry
			f.PC += 2
			CPentry := f.CP.CpIndex[CPslot]
			if CPentry.Type != classloader.MethodRef { // the pointed-to CP entry must be a method reference
				return fmt.Errorf("Expected a method ref for invokevirtual, but got %d in"+
					"location %d in method %s of class %s\n",
					CPentry.Type, f.PC, f.MethName, f.ClName)
			}

			// get the methodRef entry
			method := f.CP.MethodRefs[CPentry.Slot]

			// get the class entry from this method
			classRef := method.ClassIndex
			classNameIndex := f.CP.ClassRefs[f.CP.CpIndex[classRef].Slot]
			classNameEntry := f.CP.CpIndex[classNameIndex]
			className := f.CP.Utf8Refs[classNameEntry.Slot]

			// get the method name for this method
			nAndTindex := method.NameAndType
			nAndTentry := f.CP.CpIndex[nAndTindex]
			nAndTslot := nAndTentry.Slot
			nAndT := f.CP.NameAndTypes[nAndTslot]
			methodNameIndex := nAndT.NameIndex
			methodName := classloader.FetchUTF8stringFromCPEntryNumber(f.CP, methodNameIndex)
			methodName = className + "." + methodName

			// get the signature for this method
			methodSigIndex := nAndT.DescIndex
			methodType := classloader.FetchUTF8stringFromCPEntryNumber(f.CP, methodSigIndex)
			// println("Method signature for invokevirtual: " + methodName + methodType)

			v := classloader.MTable[methodName+methodType]
			if v.Meth != nil && v.MType == 'G' { // so we have a golang function
				_, err := runGmethod(v, fs, className, methodName, methodType)
				if err != nil {
					shutdown.Exit(shutdown.APP_EXCEPTION) // any exceptions message will already have been displayed to the user
				}
				break
			}
		case INVOKESTATIC: // 	0xB8 invokestatic (create new frame, invoke static function)
			CPslot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2]) // next 2 bytes point to CP entry
			f.PC += 2
			CPentry := f.CP.CpIndex[CPslot]
			// get the methodRef entry
			method := f.CP.MethodRefs[CPentry.Slot]

			// get the class entry from this method
			classRef := method.ClassIndex
			classNameIndex := f.CP.ClassRefs[f.CP.CpIndex[classRef].Slot]
			classNameEntry := f.CP.CpIndex[classNameIndex]
			className := f.CP.Utf8Refs[classNameEntry.Slot]

			// get the method name for this method
			nAndTindex := method.NameAndType
			nAndTentry := f.CP.CpIndex[nAndTindex]
			nAndTslot := nAndTentry.Slot
			nAndT := f.CP.NameAndTypes[nAndTslot]
			methodNameIndex := nAndT.NameIndex
			methodName := classloader.FetchUTF8stringFromCPEntryNumber(f.CP, methodNameIndex)
			// println("Method name for invokestatic: " + className + "." + methodName)

			// get the signature for this method
			methodSigIndex := nAndT.DescIndex
			methodType := classloader.FetchUTF8stringFromCPEntryNumber(f.CP, methodSigIndex)
			// println("Method signature for invokestatic: " + methodName + methodType)

			// m, cpp, err := fetchMethodAndCP(className, methodName, methodType)
			mtEntry, err := classloader.FetchMethodAndCP(className, methodName, methodType)
			if err != nil {
				return errors.New("Class not found: " + className + methodName)
			}

			if mtEntry.MType == 'G' {
				f, err = runGmethod(mtEntry, fs, className, className+"."+methodName, methodType)
				if err != nil {
					shutdown.Exit(shutdown.APP_EXCEPTION) // any exceptions message will already have been displayed to the user
				}
			} else if mtEntry.MType == 'J' {
				m := mtEntry.Meth.(classloader.JmEntry)
				maxStack := m.MaxStack
				fram := frames.CreateFrame(maxStack)

				fram.ClName = className
				fram.MethName = methodName
				fram.CP = m.Cp                     // add its pointer to the class CP
				for i := 0; i < len(m.Code); i++ { // copy the bytecodes over
					fram.Meth = append(fram.Meth, m.Code[i])
				}

				// allocate the local variables
				for k := 0; k < m.MaxLocals; k++ {
					fram.Locals = append(fram.Locals, 0)
				}

				// pop the parameters off the present stack and put them in the new frame's locals
				var argList []interface{}
				paramsToPass := util.ParseIncomingParamsFromMethTypeString(methodType)
				if len(paramsToPass) > 0 {
					for i := 0; i < len(paramsToPass); i++ {
						switch paramsToPass[i] {
						case 'D':
							arg := pop(f).(float64)
							argList = append(argList, arg)
							argList = append(argList, arg)
							pop(f)
						case 'F':
							arg := pop(f).(float64)
							argList = append(argList, arg)
						case 'J': // long
							arg := pop(f).(int64)
							argList = append(argList, arg)
							argList = append(argList, arg)
							pop(f)
						default:
							arg := pop(f).(int64)
							argList = append(argList, arg)
						}
					}
				}

				destLocal := 0
				for j := len(argList) - 1; j >= 0; j-- {
					fram.Locals[destLocal] = argList[j]
					destLocal += 1
				}
				fram.TOS = -1

				fs.PushFront(fram)                   // push the new frame
				f = fs.Front().Value.(*frames.Frame) // point f to the new head
				err = runFrame(fs)
				if err != nil {
					return err
				}

				// if the static method is main(), when we get here the
				// frame stack will be empty to exit from here, otherwise
				// there's still a frame on the stack, pop it off and continue.
				if fs.Len() == 0 {
					return nil
				}
				fs.Remove(fs.Front()) // pop the frame off

				// the previous frame pop might have been main()
				// if so, then we can't reset f to a non-existent frame
				// so we test for this before resetting f.
				if fs.Len() != 0 {
					f = fs.Front().Value.(*frames.Frame)
				} else {
					return nil
				}
			}
		case NEW: // 0xBB 	new: create and instantiate a new object
			CPslot := (int(f.Meth[f.PC+1]) * 256) + int(f.Meth[f.PC+2]) // next 2 bytes point to CP entry
			f.PC += 2
			CPentry := f.CP.CpIndex[CPslot]
			if CPentry.Type != classloader.ClassRef && CPentry.Type != classloader.Interface {
				msg := fmt.Sprintf("Invalid type for new object")
				_ = log.Log(msg, log.SEVERE)
			}

			// the classref points to a UTF8 record with the name of the class to instantiate
			var className string
			if CPentry.Type == classloader.ClassRef {
				utf8Index := f.CP.ClassRefs[CPentry.Slot]
				className = classloader.FetchUTF8stringFromCPEntryNumber(f.CP, utf8Index)
			}

			ref, err := instantiateClass(className)
			if err != nil {
				_ = log.Log("Error instantiating class: "+className, log.SEVERE)
				return errors.New("error instantiating class")
			}

			// to push the object reference as an int64, it must first be converted to an unsafe pointer
			rawRef := uintptr(unsafe.Pointer(ref))
			push(f, int64(rawRef))

		default:
			msg := fmt.Sprintf("Invalid bytecode found: %d at location %d in method %s() of class %s\n",
				f.Meth[f.PC], f.PC, f.MethName, f.ClName)
			_ = log.Log(msg, log.SEVERE)
			return errors.New("invalid bytecode encountered")
		}
		f.PC += 1
	}
	return nil
}

// pop from the operand stack. TODO: need to put in checks for invalid pops
func pop(f *frames.Frame) interface{} {
	value := f.OpStack[f.TOS]
	f.TOS -= 1
	return value
}

// returns the value at the top of the stack without popping it off.
func peek(f *frames.Frame) interface{} {
	return f.OpStack[f.TOS]
}

// push onto the operand stack
func push(f *frames.Frame, x interface{}) {
	f.TOS += 1
	f.OpStack[f.TOS] = x
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
