/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 *
 * ================================================
 * THIS IS AN EXPERIMENTAL ALTERNATIVE TO run.go
 * The chages it makes:
 *  - Uses an array of functions rather than a switch for each bytecode
 *  - Does only one push and pull for 64-bit values (longs and doubles)
 *  - All severe errors use ThrowEx() to throw an exception. No errors based on return values.
 */

package jvm

import (
	"container/list"
	"errors"
	"fmt"
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/exceptions"
	"jacobin/frames"
	"jacobin/gfunction"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/object"
	"jacobin/statics"
	"jacobin/stringPool"
	"jacobin/types"
	"jacobin/util"
	"runtime/debug"
)

// set up a DispatchTable with 203 slots that correspond to the bytecodes
// each slot being a pointer to a function that accepts a pointer to the
// current frame and an int parameter. It returns an int that indicates
// how much to increase that frame's PC (program counter) by.
type BytecodeFunc func(*frames.Frame, int64) int

var DispatchTable = [203]BytecodeFunc{
	doNop,           // NOP         0x00
	doAconstNull,    // ACONST_NULL 0x01
	doIconstM1,      // ICONST_M1   0x02
	doIconst0,       // ICONST_0    0x03
	doIconst1,       // ICONST_1    0x04
	doIconst2,       // ICONST_2    0x05
	doIconst3,       // ICONST_3    0x06
	doIconst4,       // ICONST_4    0x07
	doIconst5,       // ICONST_5    0x08
	doLconst0,       // LCONST_0    0x09
	doLconst1,       // LCONST_1    0x0A
	doFconst0,       // FCONST_0    0x0B
	doFconst1,       // FCONST_1    0x0C
	doFconst2,       // FCONST_2    0x0D
	doDconst0,       // DCONST_0    0x0E
	doDconst1,       // DCONST_1    0x0F
	doBiPush,        // BIPUSH      0x10
	notImplemented,  // SIPUSH      0x11
	doLdc,           // LDC             0x12
	doLdcw,          // LDC_W           0x13
	notImplemented,  // LDC2_W          0x14
	notImplemented,  // ILOAD           0x15
	notImplemented,  // LLOAD           0x16
	notImplemented,  // FLOAD           0x17
	notImplemented,  // DLOAD           0x18
	notImplemented,  // ALOAD           0x19
	doIload0,        // ILOAD_0         0x1A
	doIload1,        // ILOAD_1         0x1B
	doIload2,        // ILOAD_2         0x1C
	doIload3,        // ILOAD_3         0x1D
	notImplemented,  // LLOAD_0         0x1E
	notImplemented,  // LLOAD_1         0x1F
	notImplemented,  // LLOAD_2         0x20
	notImplemented,  // LLOAD_3         0x21
	notImplemented,  // FLOAD_0         0x22
	notImplemented,  // FLOAD_1         0x23
	notImplemented,  // FLOAD_2         0x24
	notImplemented,  // FLOAD_3         0x25
	notImplemented,  // DLOAD_0         0x26
	notImplemented,  // DLOAD_1         0x27
	notImplemented,  // DLOAD_2         0x28
	notImplemented,  // DLOAD_3         0x29
	notImplemented,  // ALOAD_0         0x2A
	notImplemented,  // ALOAD_1         0x2B
	notImplemented,  // ALOAD_2         0x2C
	notImplemented,  // ALOAD_3         0x2D
	notImplemented,  // IALOAD          0x2E
	notImplemented,  // LALOAD          0x2F
	notImplemented,  // FALOAD          0x30
	notImplemented,  // DALOAD          0x31
	notImplemented,  // AALOAD          0x32
	notImplemented,  // BALOAD          0x33
	notImplemented,  // CALOAD          0x34
	notImplemented,  // SALOAD          0x35
	notImplemented,  // ISTORE          0x36
	notImplemented,  // LSTORE          0x37
	notImplemented,  // FSTORE          0x38
	notImplemented,  // DSTORE          0x39
	notImplemented,  // ASTORE          0x3A
	doIstore0,       // ISTORE_0        0x3B
	doIstore1,       // ISTORE_1        0x3C
	doIstore2,       // ISTORE_2        0x3D
	doIstore3,       // ISTORE_3        0x3E
	notImplemented,  // LSTORE_0        0x3F
	notImplemented,  // LSTORE_1        0x40
	notImplemented,  // LSTORE_2        0x41
	notImplemented,  // LSTORE_3        0x42
	notImplemented,  // FSTORE_0        0x43
	notImplemented,  // FSTORE_1        0x44
	notImplemented,  // FSTORE_2        0x45
	notImplemented,  // FSTORE_3        0x46
	notImplemented,  // DSTORE_0        0x47
	notImplemented,  // DSTORE_1        0x48
	notImplemented,  // DSTORE_2        0x49
	notImplemented,  // DSTORE_3        0x4A
	notImplemented,  // ASTORE_0        0x4B
	notImplemented,  // ASTORE_1        0x4C
	notImplemented,  // ASTORE_2        0x4D
	notImplemented,  // ASTORE_3        0x4E
	notImplemented,  // IASTORE         0x4F
	notImplemented,  // LASTORE         0x50
	notImplemented,  // FASTORE         0x51
	notImplemented,  // DASTORE         0x52
	notImplemented,  // AASTORE         0x53
	notImplemented,  // BASTORE         0x54
	notImplemented,  // CASTORE         0x55
	notImplemented,  // SASTORE         0x56
	notImplemented,  // POP             0x57
	notImplemented,  // POP2            0x58
	notImplemented,  // DUP             0x59
	notImplemented,  // DUP_X1          0x5A
	notImplemented,  // DUP_X2          0x5B
	notImplemented,  // DUP2            0x5C
	notImplemented,  // DUP2_X1         0x5D
	notImplemented,  // DUP2_X2         0x5E
	notImplemented,  // SWAP            0x5F
	notImplemented,  // IADD            0x60
	notImplemented,  // LADD            0x61
	notImplemented,  // FADD            0x62
	notImplemented,  // DADD            0x63
	notImplemented,  // ISUB            0x64
	notImplemented,  // LSUB            0x65
	notImplemented,  // FSUB            0x66
	notImplemented,  // DSUB            0x67
	notImplemented,  // IMUL            0x68
	notImplemented,  // LMUL            0x69
	notImplemented,  // FMUL            0x6A
	notImplemented,  // DMUL            0x6B
	notImplemented,  // IDIV            0x6C
	notImplemented,  // LDIV            0x6D
	notImplemented,  // FDIV            0x6E
	notImplemented,  // DDIV            0x6F
	notImplemented,  // IREM            0x70
	notImplemented,  // LREM            0x71
	notImplemented,  // FREM            0x72
	notImplemented,  // DREM            0x73
	notImplemented,  // INEG            0x74
	notImplemented,  // LNEG            0x75
	notImplemented,  // FNEG            0x76
	notImplemented,  // DNEG            0x77
	notImplemented,  // ISHL            0x78
	notImplemented,  // LSHL            0x79
	notImplemented,  // ISHR            0x7A
	notImplemented,  // LSHR            0x7B
	notImplemented,  // IUSHR           0x7C
	notImplemented,  // LUSHR           0x7D
	notImplemented,  // IAND            0x7E
	notImplemented,  // LAND            0x7F
	notImplemented,  // IOR             0x80
	notImplemented,  // LOR             0x81
	notImplemented,  // IXOR            0x82
	notImplemented,  // LXOR            0x83
	doIinc,          // IINC            0x84
	notImplemented,  // I2L             0x85
	notImplemented,  // I2F             0x86
	notImplemented,  // I2D             0x87
	notImplemented,  // L2I             0x88
	notImplemented,  // L2F             0x89
	notImplemented,  // L2D             0x8A
	notImplemented,  // F2I             0x8B
	notImplemented,  // F2L             0x8C
	notImplemented,  // F2D             0x8D
	notImplemented,  // D2I             0x8E
	notImplemented,  // D2L             0x8F
	notImplemented,  // D2F             0x90
	notImplemented,  // I2B             0x91
	notImplemented,  // I2C             0x92
	notImplemented,  // I2S             0x93
	notImplemented,  // LCMP            0x94
	notImplemented,  // FCMPL           0x95
	notImplemented,  // FCMPG           0x96
	notImplemented,  // DCMPL           0x97
	notImplemented,  // DCMPG           0x98
	notImplemented,  // IFEQ            0x99
	notImplemented,  // IFNE            0x9A
	notImplemented,  // IFLT            0x9B
	notImplemented,  // IFGE            0x9C
	notImplemented,  // IFGT            0x9D
	notImplemented,  // IFLE            0x9E
	notImplemented,  // IF_ICMPEQ       0x9F
	notImplemented,  // IF_ICMPNE       0xA0
	notImplemented,  // IF_ICMPLT       0xA1
	doIfIcmpge,      // IF_ICMPGE       0xA2
	notImplemented,  // IF_ICMPGT       0xA3
	notImplemented,  // IF_ICMPLE       0xA4
	notImplemented,  // IF_ACMPEQ       0xA5
	notImplemented,  // IF_ACMPNE       0xA6
	doGoto,          // GOTO            0xA7
	notImplemented,  // JSR             0xA8
	notImplemented,  // RET             0xA9
	notImplemented,  // TABLESWITCH     0xAA
	notImplemented,  // LOOKUPSWITCH    0xAB
	notImplemented,  // IRETURN         0xAC
	notImplemented,  // LRETURN         0xAD
	notImplemented,  // FRETURN         0xAE
	notImplemented,  // DRETURN         0xAF
	notImplemented,  // ARETURN         0xB0
	notImplemented,  // RETURN          0xB1
	doGetStatic,     // GETSTATIC       0xB2
	notImplemented,  // PUTSTATIC       0xB3
	notImplemented,  // GETFIELD        0xB4
	notImplemented,  // PUTFIELD        0xB5
	doInvokeVirtual, // INVOKEVIRTUAL   0xB6
	notImplemented,  // INVOKESPECIAL   0xB7
	notImplemented,  // INVOKESTATIC    0xB8
	notImplemented,  // INVOKEINTERFACE 0xB9
	notImplemented,  // INVOKEDYNAMIC   0xBA
	notImplemented,  // NEW             0xBB
	notImplemented,  // NEWARRAY        0xBC
	notImplemented,  // ANEWARRAY       0xBD
	notImplemented,  // ARRAYLENGTH     0xBE
	notImplemented,  // ATHROW          0xBF
	notImplemented,  // CHECKCAST       0xC0
	notImplemented,  // INSTANCEOF      0xC1
	notImplemented,  // MONITORENTER    0xC2
	notImplemented,  // MONITOREXIT     0xC3
	notImplemented,  // WIDE            0xC4
	notImplemented,  // MULTIANEWARRAY  0xC5
	notImplemented,  // IFNULL          0xC6
	notImplemented,  // IFNONNULL       0xC7
	notImplemented,  // GOTO_W          0xC8
	notImplemented,  // JSR_W           0xC9
	notImplemented,  // BREAKPOINT      0xCA
}

// the main interpreter loop
func interpret(fs *list.List) {
	fr := fs.Front().Value.(*frames.Frame)
	if fr.FrameStack == nil { // make sure the can reference the frame stack
		fr.FrameStack = fs
	}

	for fr.PC < len(fr.Meth) {
		if MainThread.Trace {
			traceInfo := emitTraceData(fr)
			_ = log.Log(traceInfo, log.TRACE_INST)
		}

		opcode := fr.Meth[fr.PC]
		ret := DispatchTable[opcode](fr, 0)
		switch ret {
		case 0: // shows that a new frame has been pushed, so execute it
			fr = fs.Front().Value.(*frames.Frame)
			fr.PC = 0
		case exceptions.ERROR_OCCURRED: // occurs only in tests
			break
		default:
			fr.PC += ret
		}
	}
}

// the functions, listed here in numerical order of the bytecode
func doNop(_ *frames.Frame, _ int64) int { return 1 } // 0x00

func doAconstNull(fr *frames.Frame, _ int64) int { // 0x01 ACONST_NULL push null onto stack
	push(fr, object.Null)
	return 1
}

// 0x02 - 0x0A ICONST and LCONST, push int or long onto stack
func doIconstM1(fr *frames.Frame, _ int64) int { return pushInt(fr, int64(-1)) }
func doIconst0(fr *frames.Frame, _ int64) int  { return pushInt(fr, int64(0)) }
func doIconst1(fr *frames.Frame, _ int64) int  { return pushInt(fr, int64(1)) }
func doIconst2(fr *frames.Frame, _ int64) int  { return pushInt(fr, int64(2)) }
func doIconst3(fr *frames.Frame, _ int64) int  { return pushInt(fr, int64(3)) }
func doIconst4(fr *frames.Frame, _ int64) int  { return pushInt(fr, int64(4)) }
func doIconst5(fr *frames.Frame, _ int64) int  { return pushInt(fr, int64(5)) }
func doLconst0(fr *frames.Frame, _ int64) int  { return pushInt(fr, int64(0)) }
func doLconst1(fr *frames.Frame, _ int64) int  { return pushInt(fr, int64(1)) }
func doFconst0(fr *frames.Frame, _ int64) int  { return pushFloat(fr, int64(0)) }
func doFconst1(fr *frames.Frame, _ int64) int  { return pushFloat(fr, int64(1)) }
func doFconst2(fr *frames.Frame, _ int64) int  { return pushFloat(fr, int64(2)) }
func doDconst0(fr *frames.Frame, _ int64) int  { return pushFloat(fr, int64(0)) }
func doDconst1(fr *frames.Frame, _ int64) int  { return pushFloat(fr, int64(1)) }

func doBiPush(fr *frames.Frame, _ int64) int { // 0x10 BIPUSH push following byte onto stack
	wbyte := fr.Meth[fr.PC+1]
	wint64 := byteToInt64(wbyte)
	push(fr, wint64)
	return 2
}

// 0x12, 0x13 LDC functions
func doLdc(fr *frames.Frame, _ int64) int  { return ldc(fr, 1) }
func doLdcw(fr *frames.Frame, _ int64) int { return ldc(fr, 2) }

// 0x3B - 0x3E ISTORE_0 thru _3: Store popped TOS into locals specified as 0-3 in bytecode name
func doIstore0(fr *frames.Frame, _ int64) int { return storeInt(fr, int64(0)) }
func doIstore1(fr *frames.Frame, _ int64) int { return storeInt(fr, int64(1)) }
func doIstore2(fr *frames.Frame, _ int64) int { return storeInt(fr, int64(2)) }
func doIstore3(fr *frames.Frame, _ int64) int { return storeInt(fr, int64(3)) }

func doIload0(fr *frames.Frame, _ int64) int { return loadInt(fr, int64(0)) }
func doIload1(fr *frames.Frame, _ int64) int { return loadInt(fr, int64(1)) }
func doIload2(fr *frames.Frame, _ int64) int { return loadInt(fr, int64(2)) }
func doIload3(fr *frames.Frame, _ int64) int { return loadInt(fr, int64(3)) }

func doIinc(fr *frames.Frame, _ int64) int {
	var index int
	var increment int64
	var PCtoSkip int
	if fr.WideInEffect { // if wide is in effect, index  and increment are two bytes wide, otherwise one byte each
		index = (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2])
		increment = int64(fr.Meth[fr.PC+1])*256 + int64(fr.Meth[fr.PC+2])
		PCtoSkip = 4
		fr.WideInEffect = false
	} else {
		index = int(fr.Meth[fr.PC+1])
		increment = byteToInt64(fr.Meth[fr.PC+2])
		PCtoSkip = 2
	}
	orig := fr.Locals[index].(int64)
	fr.Locals[index] = orig + increment
	return PCtoSkip + 1
}
func doIfIcmpge(fr *frames.Frame, _ int64) int { // 0xA2 IF_ICMPGE Compare ints for >=
	popValue := pop(fr)
	val2 := convertInterfaceToInt64(popValue)
	popValue = pop(fr)
	val1 := convertInterfaceToInt64(popValue)
	if val1 >= val2 { // if comp succeeds, next 2 bytes hold instruction index
		jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
		return int(jumpTo)
	} else {
		return 3 // the 2 bytes forming the unused jumpTo + 1 byte to next bytecode
	}
}

func doGoto(fr *frames.Frame, _ int64) int { // 0xA7 GOTO unconditional jump within method
	jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
	return int(jumpTo) // note the value can be negative to jump to earlier bytecode
}

func doGetStatic(fr *frames.Frame, _ int64) int { // 0xB2 GETSTATIC
	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2]) // next 2 bytes point to CP entry
	// f.PC += 2
	CP := fr.CP.(*classloader.CPool)
	CPentry := CP.CpIndex[CPslot]
	if CPentry.Type != classloader.FieldRef { // the pointed-to CP entry must be a field reference
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("GETSTATIC: Expected a field ref, but got %d in"+
			"location %d in method %s of class %s\n",
			CPentry.Type, fr.PC, fr.MethName, fr.ClName)
		exceptions.ThrowEx(excNames.NoSuchFieldException, errMsg, fr)
	}

	// get the field entry
	field := CP.FieldRefs[CPentry.Slot]

	// get the class entry from the field entry for this field. It's the class name.
	classRef := field.ClassIndex
	classNameIndex := CP.ClassRefs[CP.CpIndex[classRef].Slot]
	classNamePtr := stringPool.GetStringPointer(classNameIndex)
	className := *classNamePtr

	// process the name and type entry for this field
	nAndTindex := field.NameAndType
	nAndTentry := CP.CpIndex[nAndTindex]
	nAndTslot := nAndTentry.Slot
	nAndT := CP.NameAndTypes[nAndTslot]
	fieldNameIndex := nAndT.NameIndex
	fieldName := classloader.FetchUTF8stringFromCPEntryNumber(CP, fieldNameIndex)
	fieldName = className + "." + fieldName
	if MainThread.Trace {
		emitTraceFieldID("GETSTATIC", fieldName)
	}

	// was this static field previously loaded? Is so, get its location and move on.
	prevLoaded, ok := statics.Statics[fieldName]
	if !ok { // if field is not already loaded, then
		// the class has not been instantiated, so
		// instantiate the class
		_, err := InstantiateClass(className, fr.FrameStack)
		if err == nil {
			prevLoaded, ok = statics.Statics[fieldName]
		} else {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := fmt.Sprintf("GETSTATIC: could not load class %s", className)
			_ = log.Log(errMsg, log.SEVERE)
			exceptions.ThrowEx(excNames.ClassNotFoundException, errMsg, fr)
		}
	}

	// if the field can't be found even after instantiating the
	// containing class, something is wrong so get out of here.
	if !ok {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("GETSTATIC: could not find static field %s in class %s"+
			"\n", fieldName, className)
		exceptions.ThrowEx(excNames.NoSuchFieldException, errMsg, fr)
	}

	switch prevLoaded.Value.(type) {
	case bool:
		// a boolean, which might
		// be stored as a boolean, a byte (in an array), or int64
		// We want all forms normalized to int64
		value := prevLoaded.Value.(bool)
		prevLoaded.Value =
			types.ConvertGoBoolToJavaBool(value)
		push(fr, prevLoaded.Value)
	case byte:
		value := prevLoaded.Value.(byte)
		prevLoaded.Value = int64(value)
		push(fr, prevLoaded.Value)
	case int:
		value := prevLoaded.Value.(int)
		push(fr, int64(value))
	default:
		push(fr, prevLoaded.Value)
	}
	return 3 // 2 for the CP slot + 1 for the next bytecode
}

func doInvokeVirtual(fr *frames.Frame, _ int64) int { // 0xB6 INVOKEVIRTUAL
	var err error
	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2]) // next 2 bytes point to CP entry
	fr.PC += 2
	CP := fr.CP.(*classloader.CPool)
	CPentry := CP.CpIndex[CPslot]
	if CPentry.Type != classloader.MethodRef { // the pointed-to CP entry must be a method reference
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("INVOKEVIRTUAL: Expected a method ref, but got %d in"+
			"location %d in method %s of class %s\n",
			CPentry.Type, fr.PC, fr.MethName, fr.ClName)
		status := exceptions.ThrowEx(excNames.WrongMethodTypeException, errMsg, fr)
		if status != exceptions.Caught {
			return exceptions.ERROR_OCCURRED // applies only if in test
		}
	}

	className, methodName, methodType :=
		classloader.GetMethInfoFromCPmethref(CP, CPslot)

	mtEntry := classloader.MTable[className+"."+methodName+methodType]
	if mtEntry.Meth == nil { // if the method is not in the method table, find it
		mtEntry, err = classloader.FetchMethodAndCP(className, methodName, methodType)
		if err != nil || mtEntry.Meth == nil {
			// TODO: search the superclasses, then the classpath and retry
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKEVIRTUAL: Class method not found: " + className + "." + methodName + methodType
			status := exceptions.ThrowEx(excNames.UnsupportedOperationException, errMsg, fr)
			if status != exceptions.Caught {
				return exceptions.ERROR_OCCURRED // applies only if in test
			}
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
			params = append(params, pop(fr))
		}

		// now get the objectRef (the object whose method we're invoking) or a *os.File (stream I/O)
		popped := pop(fr)
		params = append(params, popped)

		ret := gfunction.RunGfunction(mtEntry, fr.FrameStack, className, methodName, methodType, &params, true, MainThread.Trace)
		// if err != nil {
		if ret != nil {
			switch ret.(type) {
			case error: // only occurs in testing
				if globals.GetGlobalRef().JacobinName == "test" {
					return exceptions.ERROR_OCCURRED
				}
				if errors.Is(ret.(error), gfunction.CaughtGfunctionException) {
					return 3 // 2 for CP slot + 1 for next bytecode
				}
			default: // if it's not an error, then it's a legitimate return value, which we simply push
				push(fr, ret)
			}
			// any exception will already have been handled.
		}
		return 3 // 2 for CP slot + 1 for next bytecode
	}

	if mtEntry.MType == 'J' { // it's a Java or Native function
		m := mtEntry.Meth.(classloader.JmEntry)
		if m.AccessFlags&0x0100 > 0 {
			// Native code
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKEVIRTUAL: Native method requested: " + className + "." + methodName + methodType
			status := exceptions.ThrowEx(excNames.UnsupportedOperationException, errMsg, fr)
			if status != exceptions.Caught {
				return exceptions.ERROR_OCCURRED // applies only if in test
			}
		}
		fram, err := createAndInitNewFrame(
			className, methodName, methodType, &m, true, fr)
		if err != nil {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKEVIRTUAL: Error creating frame in: " + className + "." + methodName + methodType
			status := exceptions.ThrowEx(excNames.InvalidStackFrameException, errMsg, fr)
			if status != exceptions.Caught {
				return exceptions.ERROR_OCCURRED // applies only if in test
			}
		}

		fr.PC += 1                    // move to next bytecode before exiting
		fr.FrameStack.PushFront(fram) // push the new frame
		return 0
	}
	return exceptions.ERROR_OCCURRED
}

func notImplemented(_ *frames.Frame, _ int64) int {
	return 1
}

// === helper methods--that is, functions called by dispatched methods (in alpha order) ===

func loadInt(fr *frames.Frame, local int64) int {
	push(fr, fr.Locals[local])
	return 1
}

func ldc(fr *frames.Frame, width int) int {
	var idx int
	if width == 1 { // LDC uses a 1-byte index into the CP, LDC_W uses a 2-byte index
		idx = int(fr.Meth[fr.PC+1])
	} else {
		idx = (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2])
	}

	CPe := classloader.FetchCPentry(fr.CP.(*classloader.CPool), idx)
	if CPe.EntryType == 0 || // 0 = error
		// Note: an invalid CP entry causes a java.lang.Verify error and
		//       is caught before execution of the program begins.
		// This bytecode does not load longs or doubles
		CPe.EntryType == classloader.DoubleConst ||
		CPe.EntryType == classloader.LongConst {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, LDC: Invalid type for bytecode operand",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName)
		status := exceptions.ThrowEx(excNames.ClassFormatError, errMsg, fr)
		if status != exceptions.Caught {
			return exceptions.ERROR_OCCURRED // applies only if in test
		}
	}
	// if no error
	switch CPe.RetType {
	case classloader.IS_INT64:
		push(fr, CPe.IntVal)
	case classloader.IS_FLOAT64:
		push(fr, CPe.FloatVal)
	case classloader.IS_STRUCT_ADDR:
		push(fr, CPe.AddrVal)
	case classloader.IS_STRING_ADDR: // returns a string object whose "value" field is a byte array
		stringAddr := object.StringObjectFromGoString(*CPe.StringVal)
		push(fr, stringAddr)
	}

	if width == 1 {
		return 2 // 1 for the index + 1 for the next bytecode
	} else {
		return 3 // 2 for the index + 1 for the next bytecode
	}
}

func pushInt(fr *frames.Frame, intToPush int64) int {
	push(fr, intToPush)
	return 1
}

func pushFloat(fr *frames.Frame, intToPush int64) int {
	push(fr, float64(intToPush))
	return 1
}

func storeInt(fr *frames.Frame, local int64) int {
	fr.Locals[local] = pop(fr)
	return 1
}
