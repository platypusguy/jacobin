/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024-5 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

// This is the main interpreter loop. Bytecodes are executed by using the bytecode
// value as an index into a dispatch table. The dispatch table contains pointers
// to functions that implement the bytecode.

package jvm

import (
	"container/list"
	"encoding/binary"
	"errors"
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/exceptions"
	"jacobin/src/frames"
	"jacobin/src/gfunction"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/opcodes"
	"jacobin/src/shutdown"
	"jacobin/src/statics"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"jacobin/src/types"
	"jacobin/src/util"
	"math"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
)

var mtrdebug = false // set to true to enable debug output for MONITORENTER and EXIT

// set up a DispatchTable with 203 slots that correspond to the bytecodes
// each slot being a pointer to a function that accepts a pointer to the
// current frame and an int parameter. It returns an int that indicates
// how much to increase that frame's PC (program counter) by.

type BytecodeFunc func(*frames.Frame, int64) int

var DispatchTable = [203]BytecodeFunc{
	doNothing,         // NOP             0x00
	doAconstNull,      // ACONST_NULL     0x01
	doIconstM1,        // ICONST_M1       0x02
	doIconst0,         // ICONST_0        0x03
	doIconst1,         // ICONST_1        0x04
	doIconst2,         // ICONST_2        0x05
	doIconst3,         // ICONST_3        0x06
	doIconst4,         // ICONST_4        0x07
	doIconst5,         // ICONST_5        0x08
	doLconst0,         // LCONST_0        0x09
	doLconst1,         // LCONST_1        0x0A
	doFconst0,         // FCONST_0        0x0B
	doFconst1,         // FCONST_1        0x0C
	doFconst2,         // FCONST_2        0x0D
	doDconst0,         // DCONST_0        0x0E
	doDconst1,         // DCONST_1        0x0F
	doBipush,          // BIPUSH          0x10
	doSipush,          // SIPUSH          0x11
	doLdc,             // LDC             0x12
	doLdcw,            // LDC_W           0x13
	doLdc2w,           // LDC2_W          0x14
	doLoad,            // ILOAD           0x15
	doLoad,            // LLOAD           0x16
	doLoad,            // FLOAD           0x17
	doLoad,            // DLOAD           0x18
	doLoad,            // ALOAD           0x19
	doIload0,          // ILOAD_0         0x1A
	doIload1,          // ILOAD_1         0x1B
	doIload2,          // ILOAD_2         0x1C
	doIload3,          // ILOAD_3         0x1D
	doIload0,          // LLOAD_0         0x1E
	doIload1,          // LLOAD_1         0x1F
	doIload2,          // LLOAD_2         0x20
	doIload3,          // LLOAD_3         0x21
	doFload0,          // FLOAD_0         0x22
	doFload1,          // FLOAD_1         0x23
	doFload2,          // FLOAD_2         0x24
	doFload3,          // FLOAD_3         0x25
	doFload0,          // DLOAD_0         0x26
	doFload1,          // DLOAD_1         0x27
	doFload2,          // DLOAD_2         0x28
	doFload3,          // DLOAD_3         0x29
	doAload0,          // ALOAD_0         0x2A
	doAload1,          // ALOAD_1         0x2B
	doAload2,          // ALOAD_2         0x2C
	doAload3,          // ALOAD_3         0x2D
	doIaload,          // IALOAD          0x2E
	doIaload,          // LALOAD          0x2F
	doFaload,          // FALOAD          0x30
	doFaload,          // DALOAD          0x31
	doAaload,          // AALOAD          0x32
	doBaload,          // BALOAD          0x33
	doIaload,          // CALOAD          0x34
	doIaload,          // SALOAD          0x35
	doIstore,          // ISTORE          0x36
	doIstore,          // LSTORE          0x37
	doFstore,          // FSTORE          0x38
	doFstore,          // DSTORE          0x39
	doAstore,          // ASTORE          0x3A
	doIstore0,         // ISTORE_0        0x3B
	doIstore1,         // ISTORE_1        0x3C
	doIstore2,         // ISTORE_2        0x3D
	doIstore3,         // ISTORE_3        0x3E
	doIstore0,         // LSTORE_0        0x3F
	doIstore1,         // LSTORE_1        0x40
	doIstore2,         // LSTORE_2        0x41
	doIstore3,         // LSTORE_3        0x42
	doFstore0,         // FSTORE_0        0x43
	doFstore1,         // FSTORE_1        0x44
	doFstore2,         // FSTORE_2        0x45
	doFstore3,         // FSTORE_3        0x46
	doFstore0,         // DSTORE_0        0x47
	doFstore1,         // DSTORE_1        0x48
	doFstore2,         // DSTORE_2        0x49
	doFstore3,         // DSTORE_3        0x4A
	doAstore0,         // ASTORE_0        0x4B
	doAstore1,         // ASTORE_1        0x4C
	doAstore2,         // ASTORE_2        0x4D
	doAstore3,         // ASTORE_3        0x4E
	doIastore,         // IASTORE         0x4F
	doIastore,         // LASTORE         0x50
	doFastore,         // FASTORE         0x51
	doFastore,         // DASTORE         0x52
	doAastore,         // AASTORE         0x53
	doBastore,         // BASTORE         0x54
	doIastore,         // CASTORE         0x55
	doIastore,         // SASTORE         0x56
	doPop,             // POP             0x57
	doPop,             // POP2            0x58
	doDup,             // DUP             0x59
	doDupx1,           // DUP_X1          0x5A
	doDupx2,           // DUP_X2          0x5B
	doDup2,            // DUP2            0x5C
	doDup2x1,          // DUP2_X1         0x5D
	doDup2x2,          // DUP2_X2         0x5E
	doSwap,            // SWAP            0x5F
	doIadd,            // IADD            0x60
	doLadd,            // LADD            0x61
	doFadd,            // FADD            0x62
	doDadd,            // DADD            0x63
	doIsub,            // ISUB            0x64
	doLsub,            // LSUB            0x65
	doFsub,            // FSUB            0x66
	doDsub,            // DSUB            0x67
	doImul,            // IMUL            0x68
	doLmul,            // LMUL            0x69
	doFmul,            // FMUL            0x6A
	doDmul,            // DMUL            0x6B
	doIdiv,            // IDIV            0x6C
	doLdiv,            // LDIV            0x6D
	doFdiv,            // FDIV            0x6E
	doDdiv,            // DDIV            0x6F
	doIrem,            // IREM            0x70
	doLrem,            // LREM            0x71
	doFrem,            // FREM            0x72
	doDrem,            // DREM            0x73
	doIneg,            // INEG            0x74
	doLneg,            // LNEG            0x75
	doFneg,            // FNEG            0x76
	doDneg,            // DNEG            0x77
	doIshl,            // ISHL            0x78
	doLshl,            // LSHL            0x79
	doIshr,            // ISHR            0x7A
	doLshr,            // LSHR            0x7B
	doIushr,           // IUSHR           0x7C
	doLushr,           // LUSHR           0x7D
	doIand,            // IAND            0x7E
	doLand,            // LAND            0x7F
	doIor,             // IOR             0x80
	doLor,             // LOR             0x81
	doIxor,            // IXOR            0x82
	doLxor,            // LXOR            0x83
	doIinc,            // IINC            0x84
	doNothing,         // I2L             0x85
	doI2f,             // I2F             0x86
	doI2d,             // I2D             0x87
	doNothing,         // L2I             0x88
	doL2f,             // L2F             0x89
	doL2d,             // L2D             0x8A
	doF2i,             // F2I             0x8B
	doF2l,             // F2L             0x8C
	doF2d,             // F2D             0x8D
	doD2i,             // D2I             0x8E
	doD2l,             // D2L             0x8F
	doD2f,             // D2F             0x90
	doI2b,             // I2B             0x91
	doI2c,             // I2C             0x92
	doI2s,             // I2S             0x93
	doLcmp,            // LCMP            0x94
	doFcmpl,           // FCMPL           0x95
	doFcmpg,           // FCMPG           0x96
	doDcmpl,           // DCMPL           0x97
	doDcmpg,           // DCMPG           0x98
	doIfeq,            // IFEQ            0x99
	doIfne,            // IFNE            0x9A
	doIflt,            // IFLT            0x9B
	doIfge,            // IFGE            0x9C
	doIfgt,            // IFGT            0x9D
	doIfle,            // IFLE            0x9E
	doIficmpeq,        // IF_ICMPEQ       0x9F
	doIficmpne,        // IF_ICMPNE       0xA0
	doIficmplt,        // IF_ICMPLT       0xA1
	doIficmpge,        // IF_ICMPGE       0xA2
	doIficmpgt,        // IF_ICMPGT       0xA3
	doIficmple,        // IF_ICMPLE       0xA4
	doIfacmpeq,        // IF_ACMPEQ       0xA5
	doIfacmpne,        // IF_ACMPNE       0xA6
	doGoto,            // GOTO            0xA7
	doJsr,             // JSR             0xA8
	doRet,             // RET             0xA9
	doTableswitch,     // TABLESWITCH     0xAA
	doLookupswitch,    // LOOKUPSWITCH    0xAB
	doIreturn,         // IRETURN         0xAC
	doIreturn,         // LRETURN         0xAD
	doIreturn,         // FRETURN         0xAE
	doIreturn,         // DRETURN         0xAF
	doIreturn,         // ARETURN         0xB0
	doReturn,          // RETURN          0xB1
	nil,               // GETSTATIC       0xB2 initialized in initializeDispatchTable()
	nil,               // PUTSTATIC       0xB3 initialized in initializeDispatchTable()
	doGetfield,        // GETFIELD        0xB4
	doPutfield,        // PUTFIELD        0xB5
	doInvokeVirtual,   // INVOKEVIRTUAL   0xB6
	doInvokespecial,   // INVOKESPECIAL   0xB7
	nil,               // INVOKESTATIC    0xB8 initialized in initializeDispatchTable()
	doInvokeinterface, // INVOKEINTERFACE 0xB9
	doInvokedynamic,   // INVOKEDYNAMIC   0xBA
	nil,               // NEW             0xBB initialized in initializeDispatchTable()
	doNewarray,        // NEWARRAY        0xBC
	doAnewarray,       // ANEWARRAY       0xBD
	doArraylength,     // ARRAYLENGTH     0xBE
	doAthrow,          // ATHROW          0xBF
	doCheckcast,       // CHECKCAST       0xC0
	doInstanceof,      // INSTANCEOF      0xC1
	doMonitorenter,    // MONITORENTER    0xC2
	doMonitorexit,     // MONITOREXIT     0xC3
	doWide,            // WIDE            0xC4
	doMultinewarray,   // MULTIANEWARRAY  0xC5
	doIfnull,          // IFNULL          0xC6
	doIfnonnull,       // IFNONNULL       0xC7
	doGotow,           // GOTO_W          0xC8
	doJsrw,            // JSR_W           0xC9
	doWarninvalid,     // BREAKPOINT      0xCA not implemented, generates warning, not exception
}

// initializeDispatchTable initializes a few bytecodes that call interpret(). If they were
// initialized to their respective functions directly in the table above, golang gives a
// circularity error, By initializing those bytecodes with their methods here, the circularity
// issue goes away. Golang can't tell in the original setup that circularity will never occur.
func initializeDispatchTable() {
	DispatchTable[opcodes.GETSTATIC] = doGetStatic
	DispatchTable[opcodes.PUTSTATIC] = doPutStatic
	DispatchTable[opcodes.INVOKESTATIC] = doInvokestatic
	DispatchTable[opcodes.NEW] = doNew
}

const ( // result values from bytecode interpretation
	ERROR_OCCURRED = math.MaxInt32
	RESUME_HERE    = math.MaxInt32 - 1
	// all special result values must be greater than SPECIAL_CASE,
	// which is the value tested against in interpret()'s principal
	// loop to identify special cases
	SPECIAL_CASE = math.MaxInt32 - 10
)

// the main interpreter loop. This loop takes responsibility for
// pushing a new frame for a called method onto the stack, and for
// popping the current frame when a bytecode of the RETURN family
// is encountered. In both cases, interpret() returns and the
// RunJavaThread() loop goes to the top of the frame stack and calls
// interpret() on the frame found there, if any.
func interpret(fs *list.List) {
	const maxBytecode = byte(len(DispatchTable) - 1)
	if DispatchTable[opcodes.NEW] == nil { // test whether the table is fully initialized
		initializeDispatchTable()
	}

	fr := fs.Front().Value.(*frames.Frame)
	if fr.FrameStack == nil { // make sure we can reference the frame stack
		fr.FrameStack = fs
	}

	// Don't allow a nil code segment (E.g. mishandled abstract).
	if len(fr.Meth) < 1 {
		errMsg := "Empty code segment"
		status := exceptions.ThrowEx(excNames.VirtualMachineError, errMsg, fr)
		if status != exceptions.Caught { // will only happen in test
			globals.InitGlobals("test")
			return
		}
	}

	defer func() int {
		// only an untrapped panic gets us here
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			glob := globals.GetGlobalRef()
			glob.ErrorGoStack = stack
			exceptions.ShowPanicCause(r)
			exceptions.ShowFrameStack(fs)
			exceptions.ShowGoStackTrace(nil)
			return shutdown.Exit(shutdown.APP_EXCEPTION)
		}
		return shutdown.OK
	}()

	// the main bytecode interpreter loop
	for fr.PC < len(fr.Meth) {
		if globals.TraceInst {
			traceInfo := EmitTraceData(fr)
			trace.Trace(traceInfo)
		}

		opcode := fr.Meth[fr.PC]
		if opcode <= maxBytecode {
			ret := DispatchTable[opcode](fr, 0)
			if ret < SPECIAL_CASE && ret != 0 {
				fr.PC += ret
			} else {
				switch ret {
				case 0:
					// exiting will either end program or call this function
					// again for the frame at the top of the frame stack
					return
				case ERROR_OCCURRED: // occurs only in tests
					fs.Remove(fs.Front()) // pop the frame off, else we loop endlessly
					return
				case RESUME_HERE: // continue processing from the present fr.PC
					// This primarily occurs when an exception is caught. The catch resets
					// the PC to the catch code to execute. So, we don't need any update to
					// the PC. However, we have to refresh the current frame b/c the
					// exception will refresh the topmost frame with any exception handling
					fr = fs.Front().Value.(*frames.Frame)
				}
			}
		} else {
			errMsg := fmt.Sprintf("Invalid bytecode: %d", opcode)
			status := exceptions.ThrowEx(excNames.ClassFormatError, errMsg, fr)
			if status != exceptions.Caught { // will only happen in test
				globals.InitGlobals("test")
				return
			}
		}
	}
}

// the functions, listed here in numerical order of the bytecode
func doNothing(_ *frames.Frame, _ int64) int { return 1 } // 0x00

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

// 0x0B - 0x0F push float
func doFconst0(fr *frames.Frame, _ int64) int { return pushFloat(fr, int64(0)) }
func doFconst1(fr *frames.Frame, _ int64) int { return pushFloat(fr, int64(1)) }
func doFconst2(fr *frames.Frame, _ int64) int { return pushFloat(fr, int64(2)) }
func doDconst0(fr *frames.Frame, _ int64) int { return pushFloat(fr, int64(0)) }
func doDconst1(fr *frames.Frame, _ int64) int { return pushFloat(fr, int64(1)) }

// 0x10 BIPUSH push following byte onto stack
func doBipush(fr *frames.Frame, _ int64) int {
	wbyte := fr.Meth[fr.PC+1]
	wint64 := byteToInt64(wbyte)
	push(fr, wint64)
	return 2
}

// 0x11 SIPUSH create int from next 2 bytes and push it
func doSipush(fr *frames.Frame, _ int64) int {
	wbyte1 := fr.Meth[fr.PC+1]
	wbyte2 := fr.Meth[fr.PC+2]
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
	push(fr, wint64)
	return 3
}

// 0x12, 0x13 LDC, LDC_W load constants
func doLdc(fr *frames.Frame, _ int64) int  { return ldc(fr, 1) }
func doLdcw(fr *frames.Frame, _ int64) int { return ldc(fr, 2) }

// 0x14 LDC2_W (push long or double from CP indexed by next two bytes)
func doLdc2w(fr *frames.Frame, _ int64) int {
	idx := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2])

	CPe := classloader.FetchCPentry(fr.CP.(*classloader.CPool), idx)
	if CPe.RetType == classloader.IS_INT64 {
		push(fr, CPe.IntVal)
	} else if CPe.RetType == classloader.IS_FLOAT64 {
		push(fr, CPe.FloatVal)
	} else {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, LDC2_W: Invalid type for bytecode operand",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName)
		status := exceptions.ThrowEx(excNames.ClassFormatError, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}
	return 3 // 2 for idx + 1 for next bytecode
}

// 0x15 - 0x19: ILOAD, LLOAD, FLOAD, ALOAD
func doLoad(fr *frames.Frame, _ int64) int {
	var index int
	var PCadvance int    // how much to advance fr.PC, the program counter
	if fr.WideInEffect { // if wide is in effect, index is two bytes wide, otherwise one byte
		index = (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2])
		PCadvance = 2
		fr.WideInEffect = false
	} else {
		index = int(fr.Meth[fr.PC+1])
		PCadvance = 1
	}
	push(fr, fr.Locals[index])
	return PCadvance + 1
}

// 0x1A - 0x1D ILOAD_x push int from local x
// 0x1E - 0x2B LLOAD_x push long from local x
func doIload0(fr *frames.Frame, _ int64) int { return load(fr, int64(0)) }
func doIload1(fr *frames.Frame, _ int64) int { return load(fr, int64(1)) }
func doIload2(fr *frames.Frame, _ int64) int { return load(fr, int64(2)) }
func doIload3(fr *frames.Frame, _ int64) int { return load(fr, int64(3)) }

// 0x22 - 0x29 FLOAD_x and DLOAD_x push float from locals[x]
// These are the same as the ILOAD_x functions. However, at some point,
// we might want to verify or handle floats differently from ints.
func doFload0(fr *frames.Frame, _ int64) int { return load(fr, int64(0)) }
func doFload1(fr *frames.Frame, _ int64) int { return load(fr, int64(1)) }
func doFload2(fr *frames.Frame, _ int64) int { return load(fr, int64(2)) }
func doFload3(fr *frames.Frame, _ int64) int { return load(fr, int64(3)) }

// 0x2A - 0x2D ALOAD_x push reference value from locals[x]
func doAload0(fr *frames.Frame, _ int64) int { return load(fr, int64(0)) }
func doAload1(fr *frames.Frame, _ int64) int { return load(fr, int64(1)) }
func doAload2(fr *frames.Frame, _ int64) int { return load(fr, int64(2)) }
func doAload3(fr *frames.Frame, _ int64) int { return load(fr, int64(3)) }

// 0x2E, 0x2F IALOAD, LALOAD push contents of an int/long array element
// 0x34, 0x35 CALOAD, SALOAD push contents of a char/short array element
func doIaload(fr *frames.Frame, _ int64) int {
	var array []int64
	index := pop(fr).(int64)
	ref := pop(fr)
	switch ref.(type) {
	case *object.Object:
		obj := ref.(*object.Object)
		if object.IsNull(obj) {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := fmt.Sprintf("in %s.%s, I/C/S/LALOAD: Invalid null reference to an array",
				util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName)
			status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}
		obj.ThMutex.RLock()
		array = obj.FieldTable["value"].Fvalue.([]int64)
		obj.ThMutex.RUnlock()
	case []int64:
		array = ref.([]int64)
	default:
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, I/C/S/LALOAD: Invalid reference to an array",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName)
		status := exceptions.ThrowEx(excNames.InvalidTypeException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	if index < 0 || index >= int64(len(array)) {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, I/C/S/LALOAD: Invalid array subscript",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName)
		status := exceptions.ThrowEx(excNames.ArrayIndexOutOfBoundsException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	var value = array[index]
	opcode := fr.Meth[fr.PC]
	switch opcode {
	case opcodes.IALOAD:
		value = int64(int32(value))
	case opcodes.LALOAD:
		// value is already int64
	case opcodes.CALOAD:
		value = int64(uint16(value))
	case opcodes.SALOAD:
		value = int64(int16(value))
	}
	push(fr, value)
	return 1
}

// 0x30, 0x31 FALOAD, DALOAD push contents of a float/double array element
func doFaload(fr *frames.Frame, _ int64) int {
	var array []float64
	index := pop(fr).(int64)
	ref := pop(fr)
	switch ref.(type) {
	case []float64:
		array = ref.([]float64)
	case *object.Object:
		obj := ref.(*object.Object)
		if object.IsNull(obj) {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := fmt.Sprintf("in %s.%s, D/FALOAD: Invalid object pointer (nil)",
				util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName)
			status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}
		obj.ThMutex.RLock()
		array = (*obj).FieldTable["value"].Fvalue.([]float64)
		obj.ThMutex.RUnlock()
	default:
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, D/FALOAD: Reference invalid type of array: %T",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, ref)
		status := exceptions.ThrowEx(excNames.InvalidTypeException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	if index < 0 || index >= int64(len(array)) {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, D/FALOAD: Invalid array subscript",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName)
		status := exceptions.ThrowEx(excNames.ArrayIndexOutOfBoundsException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	var value = array[index]
	opcode := fr.Meth[fr.PC]
	if opcode == opcodes.FALOAD {
		push(fr, float64(float32(value)))
	} else {
		push(fr, value)
	}
	return 1
}

// 0x32 AALOAD push contents of a reference array element
func doAaload(fr *frames.Frame, _ int64) int {
	index := pop(fr).(int64)
	rAref := pop(fr) // the array object. Can't be cast to *Object b/c might be nil
	if rAref == nil {
		errMsg := fmt.Sprintf("in %s.%s, AALOAD: Invalid (null) reference to an array",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName)
		status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	obj := rAref.(*object.Object)
	obj.ThMutex.RLock()
	fvalue := obj.FieldTable["value"].Fvalue
	obj.ThMutex.RUnlock()
	array := fvalue.([]*object.Object)

	size := int64(len(array))
	if index >= size {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, AALOAD: Invalid array subscript: %d",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, index)
		status := exceptions.ThrowEx(excNames.ArrayIndexOutOfBoundsException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
	}

	var value = array[index]
	push(fr, value)
	return 1
}

// 0x33 BALOAD push contents of a byte/boolean array element
func doBaload(fr *frames.Frame, _ int64) int {
	index := pop(fr).(int64)
	ref := pop(fr) // the array object
	if object.IsNull(ref) {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, BALOAD: Invalid (null) reference to an array",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName)
		status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	var bAref *object.Object
	var byteArray []types.JavaByte
	var pushValue int64

	switch ref.(type) {
	case *object.Object:
		bAref = ref.(*object.Object)
		if object.IsNull(bAref) {
			byteArray = make([]types.JavaByte, 0)
			size := int64(len(byteArray))
			ret := doBaload_index_vs_size(fr, index, size)
			if ret != 0 {
				return ret
			}
			value := byteArray[index]
			pushValue = int64(value)
		} else {
			bAref.ThMutex.RLock()
			fvalue := bAref.FieldTable["value"].Fvalue
			bAref.ThMutex.RUnlock()
			switch fvalue.(type) {
			case []types.JavaByte, []types.JavaBool:
				// Get byte array from the field.
				// It's the same process for booleans and bytes because
				// the elements of both array types are of type JavaByte.
				// Reference: JVM spec.
				byteArray = fvalue.([]types.JavaByte)
				size := int64(len(byteArray))
				ret := doBaload_index_vs_size(fr, index, size)
				if ret != 0 {
					return ret
				}
				value := byteArray[index]
				pushValue = int64(value)
			case []byte: // TODO: Go byte array? Convert it to a JavaByte array
				byteArray =
					object.JavaByteArrayFromGoByteArray(fvalue.([]byte))
				size := int64(len(byteArray))
				ret := doBaload_index_vs_size(fr, index, size)
				if ret != 0 {
					return ret
				}
				value := byteArray[index]
				pushValue = int64(value)
			}
		}
	case []int8: // TODO: Raw JavaByte arrays?
		arr := ref.([]types.JavaByte)
		size := int64(len(arr))
		ret := doBaload_index_vs_size(fr, index, size)
		if ret != 0 {
			return ret
		}
		val := arr[index]
		pushValue = int64(val)
	default:
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, BALOAD: Invalid  type of object ref: %T",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, ref)
		status := exceptions.ThrowEx(excNames.InvalidTypeException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	push(fr, pushValue)
	return 1
}

func doBaload_index_vs_size(fr *frames.Frame, index, size int64) int {
	if index >= size {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, BALOAD: array size is %d but array index is %d",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, size, index)
		status := exceptions.ThrowEx(excNames.ArrayIndexOutOfBoundsException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}
	return 0
}

// 0x36, 0x37 ISTORE/LSTORE store TOS int into a local
func doIstore(fr *frames.Frame, _ int64) int {
	var index int
	var PCadvance int    // how much to advance fr.PC, the program counter
	if fr.WideInEffect { // if wide is in effect, index is two bytes wide, otherwise one byte
		index = (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2])
		PCadvance = 2
		fr.WideInEffect = false
	} else {
		index = int(fr.Meth[fr.PC+1])
		PCadvance = 1
	}

	popped := pop(fr)
	fr.Locals[index] = convertInterfaceToInt64(popped)
	return PCadvance + 1
}

// 0x38, 0x39 FSTORE and DSTORE Store popped TOS into specified local
func doFstore(fr *frames.Frame, _ int64) int {
	var index int
	var PCadvance int    // how much to advance fr.PC, the program counter
	if fr.WideInEffect { // if wide is in effect, index is two bytes wide, otherwise one byte
		index = (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2])
		PCadvance = 2
		fr.WideInEffect = false
	} else {
		index = int(fr.Meth[fr.PC+1])
		PCadvance = 1
	}
	fr.Locals[index] = pop(fr).(float64)
	return PCadvance + 1
}

// 0x3A ASTORE store popped TOS ref into locals[index]
func doAstore(fr *frames.Frame, _ int64) int {
	var index int
	var PCadvance int    // how much to advance fr.PC, the program counter
	if fr.WideInEffect { // if wide is in effect, index is two bytes wide, otherwise one byte
		index = (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2])
		PCadvance = 2
		fr.WideInEffect = false
	} else {
		index = int(fr.Meth[fr.PC+1])
		PCadvance = 1
	}

	popped := pop(fr)
	fr.Locals[index] = popped
	return PCadvance + 1
}

// 0x3B - 0x3E ISTORE_x: Store popped TOS into locals[x]
// 0x3F - 0x42 LSTORE_x:    "    "     "   "     "
func doIstore0(fr *frames.Frame, _ int64) int { return storeInt(fr, int64(0)) }
func doIstore1(fr *frames.Frame, _ int64) int { return storeInt(fr, int64(1)) }
func doIstore2(fr *frames.Frame, _ int64) int { return storeInt(fr, int64(2)) }
func doIstore3(fr *frames.Frame, _ int64) int { return storeInt(fr, int64(3)) }

// 0x4B - 0x4E ASTORE_x: Store popped address into locals[x]
func doAstore0(fr *frames.Frame, _ int64) int { return store(fr, int64(0)) }
func doAstore1(fr *frames.Frame, _ int64) int { return store(fr, int64(1)) }
func doAstore2(fr *frames.Frame, _ int64) int { return store(fr, int64(2)) }
func doAstore3(fr *frames.Frame, _ int64) int { return store(fr, int64(3)) }

// 0x43 - 0x4A FSTORE_x and DSTORE_x: Store popped TOS into locals[x]
// These are the same as the ISTORE_x functions. However, at some point,
// we might want to verify or handle floats differently from ints.
func doFstore0(fr *frames.Frame, _ int64) int { return store(fr, int64(0)) }
func doFstore1(fr *frames.Frame, _ int64) int { return store(fr, int64(1)) }
func doFstore2(fr *frames.Frame, _ int64) int { return store(fr, int64(2)) }
func doFstore3(fr *frames.Frame, _ int64) int { return store(fr, int64(3)) }

// 0x4F, 0x50 IASTORE, LASTORE store an int, long into an array
// 0x55, 0x56 CASTORE, SASTORE store an char, short into an array
func doIastore(fr *frames.Frame, _ int64) int {
	var array []int64
	value := pop(fr).(int64)
	index := pop(fr).(int64)
	ref := pop(fr)
	switch ref.(type) {
	case *object.Object:
		obj := ref.(*object.Object)
		if object.IsNull(obj) {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := fmt.Sprintf("in %s.%s, I/C/S/LASTORE: Invalid (null) reference to an array",
				util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName)
			status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}
		obj.ThMutex.RLock()
		fld := obj.FieldTable["value"]
		obj.ThMutex.RUnlock()
		if fld.Ftype != types.IntArray && fld.Ftype != types.LongArray && fld.Ftype != types.CharArray && fld.Ftype != types.ShortArray {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := fmt.Sprintf("in %s.%s, I/J/C/S/LASTORE: field type expected=[I|J|C|S, observed=%s",
				util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, fld.Ftype)
			status := exceptions.ThrowEx(excNames.ArrayStoreException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}
		array = fld.Fvalue.([]int64)
	case []int64:
		array = ref.([]int64)
	default:
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, I/C/S/LASTORE: unexpected reference type: %T",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, ref)
		status := exceptions.ThrowEx(excNames.ArrayStoreException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	size := int64(len(array))
	if index < 0 || index >= size {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, I/C/S/LASTORE: array size is %d but array index is %d",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, size, index)
		status := exceptions.ThrowEx(excNames.ArrayIndexOutOfBoundsException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	opcode := fr.Meth[fr.PC]
	switch opcode {
	case opcodes.IASTORE:
		value = int64(int32(value))
	case opcodes.LASTORE:
		// already int64
	case opcodes.CASTORE:
		value = int64(uint16(value))
	case opcodes.SASTORE:
		value = int64(int16(value))
	}
	array[index] = value
	return 1
}

// 0x51, 0x52 FASTORE, DASTORE store a float, double in a float/doubles array
func doFastore(fr *frames.Frame, _ int64) int {
	var array []float64
	value := pop(fr).(float64)
	index := pop(fr).(int64)
	ref := pop(fr)
	switch ref.(type) {
	case *object.Object:
		obj := ref.(*object.Object)
		if object.IsNull(obj) {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := fmt.Sprintf("in %s.%s, F/DASTORE: Invalid (null) reference to an array",
				util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName)
			status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}
		obj.ThMutex.RLock()
		fld := obj.FieldTable["value"]
		obj.ThMutex.RUnlock()
		if fld.Ftype != types.FloatArray && fld.Ftype != types.DoubleArray {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := fmt.Sprintf("in %s.%s, D/FASTORE: field type expected=[F, observed=%s",
				util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, fld.Ftype)
			status := exceptions.ThrowEx(excNames.ArrayStoreException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}
		array = fld.Fvalue.([]float64)
	case []float64:
		array = ref.([]float64)
	default:
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, D/FASTORE: unexpected reference type: %T",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, ref)
		status := exceptions.ThrowEx(excNames.ArrayStoreException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	size := int64(len(array))
	if index < 0 || index >= size {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, D/FASTORE: array size is %d but array index is %d",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, size, index)
		status := exceptions.ThrowEx(excNames.ArrayIndexOutOfBoundsException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	opcode := fr.Meth[fr.PC]
	if opcode == opcodes.FASTORE {
		array[index] = float64(float32(value))
	} else {
		array[index] = value
	}
	return 1
}

// 0x53 AASTORE store a ref in a ref array
func doAastore(fr *frames.Frame, _ int64) int {
	popped := pop(fr) // reference we're inserting
	if popped == nil {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s%s, AASTORE: Invalid (null) interface[any] on stack",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, fr.MethType)
		status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}
	value, ok := popped.(*object.Object) // reference we're inserting
	if !ok {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s%s, AASTORE: Illegal argument type (%T) for inserting into a reference array",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, fr.MethType, popped)
		status := exceptions.ThrowEx(excNames.IllegalArgumentException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	index := pop(fr).(int64) // index into the array

	arrayRef, ok := pop(fr).(*object.Object) // ptr to the array object
	if !ok || object.IsNull(arrayRef) {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s%s, AASTORE: Invalid (null) reference array",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, fr.MethType)
		status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	arrayRef.ThMutex.RLock()
	rawArrayObj, ok := arrayRef.FieldTable["value"]
	arrayRef.ThMutex.RUnlock()
	if !ok {
		// Should probably throw an exception here if "value" is missing
	}

	if !strings.HasPrefix(rawArrayObj.Ftype, types.RefArray) &&
		!strings.HasPrefix(rawArrayObj.Ftype, types.MultiArray) {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, AASTORE: field type must start with '[L', got %s",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, rawArrayObj.Ftype)
		status := exceptions.ThrowEx(excNames.ArrayStoreException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	// get pointer to the actual array
	rawArray := rawArrayObj.Fvalue.([]*object.Object)
	size := int64(len(rawArray))
	if index >= size {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, AASTORE: array size is %d but array index is %d",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, size, index)
		status := exceptions.ThrowEx(excNames.ArrayIndexOutOfBoundsException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	rawArray[index] = value
	return 1
}

// 0x54 BASTORE store a boolean or byte in a bool or byte array
func doBastore(fr *frames.Frame, _ int64) int {
	var rawByteArray []types.JavaByte
	var size int64

	// Pop off value to be stored.
	value := convertInterfaceToByte(pop(fr))

	// Pop off index.
	index := pop(fr).(int64)

	// Pop off array reference.
	arrayRef := pop(fr)

	switch arrayRef.(type) {
	case *object.Object: // Arrays should be wrapped in a Jacobin object.
		obj := arrayRef.(*object.Object)
		if object.IsNull(obj) {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := fmt.Sprintf("in %s.%s, BASTORE: Invalid (null) reference to an array",
				util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName)
			status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}

		// Get array field.
		obj.ThMutex.RLock()
		fld := obj.FieldTable["value"]
		obj.ThMutex.RUnlock()

		if fld.Ftype != types.ByteArray && fld.Ftype != types.BoolArray {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := fmt.Sprintf("in %s.%s, BASTORE: field type expected boolean or byte, observed=%s",
				util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, fld.Ftype)
			status := exceptions.ThrowEx(excNames.ArrayStoreException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}

		// Get byte array from the field.
		// It's the same process for booleans and bytes because
		// the elements of both array types are of type JavaByte.
		// Reference: JVM spec.
		rawByteArray = fld.Fvalue.([]types.JavaByte)
		size = int64(len(rawByteArray))
		ret := doBastore_index_vs_size(fr, index, size)
		if ret != 0 {
			return ret
		}
		rawByteArray[index] = value

	case []types.JavaByte: // Raw JavaByte array: TODO: Do non-wrapped arrays still happen?
		int8Array := arrayRef.([]types.JavaByte)
		size = int64(len(int8Array))
		ret := doBastore_index_vs_size(fr, index, size)
		if ret != 0 {
			return ret
		}
		int8Array[index] = value
		return 1

	default:
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, BASTORE: unexpected reference type: %T",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, arrayRef)
		status := exceptions.ThrowEx(excNames.ArrayStoreException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught

	}

	return 1
}

func doBastore_index_vs_size(fr *frames.Frame, index, size int64) int {
	if index >= size {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, BASTORE: array size is %d but array index is %d",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, size, index)
		status := exceptions.ThrowEx(excNames.ArrayIndexOutOfBoundsException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}
	return 0
}

// 0x57 POP pop 1 item off op stack
// 0x58 POP2 per JACOBIN-710, POP2 is used by HotSpot to pop two 32-bit values
// off the stack (for longs and doubles). However, our longs and doubles are
// ingle 64-bit pops, so POP2 is implemented as a single pop.
func doPop(fr *frames.Frame, _ int64) int {
	if fr.TOS < 0 {
		errMsg := fmt.Sprintf("stack underflow in POP/POP2 in %s.%s",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName)
		status := exceptions.ThrowEx(excNames.InternalException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}
	fr.TOS -= 1
	return 1
}

// 0x59 DUP duplicate item at TOS
func doDup(fr *frames.Frame, _ int64) int {
	tosItem := peek(fr)
	push(fr, tosItem)
	return 1
}

// 0x5A	DUP_X1	Duplicate the top stack value and insert two slots down
func doDupx1(fr *frames.Frame, _ int64) int {
	top := pop(fr)
	next := pop(fr)
	push(fr, top)
	push(fr, next)
	push(fr, top)
	return 1
}

// 0x5B	DUP_X2	Duplicate the top stack value and insert three slots down
func doDupx2(fr *frames.Frame, _ int64) int {
	top := pop(fr)
	next := pop(fr)
	third := pop(fr)
	push(fr, top)
	push(fr, third)
	push(fr, next)
	push(fr, top)
	return 1
}

// 0x5C	DUP2 Duplicate the top two stack values
func doDup2(fr *frames.Frame, _ int64) int {
	top := pop(fr)
	next := peek(fr)
	push(fr, top)
	push(fr, next)
	push(fr, top)
	return 1
}

// 0x5D	DUP2_X1	 Duplicate the top two values, three slots down
func doDup2x1(fr *frames.Frame, _ int64) int {
	top := pop(fr)
	next := pop(fr)
	third := pop(fr)
	push(fr, next) // so: top-next-third -> top-next-third->top->next
	push(fr, top)
	push(fr, third)
	push(fr, next)
	push(fr, top)
	return 1
}

// 0x5E	DUP2_X2	Duplicate the top two values, four slots down
func doDup2x2(fr *frames.Frame, _ int64) int {
	top := pop(fr)
	next := pop(fr)
	third := pop(fr)
	fourth := pop(fr)
	push(fr, next) // so: top-next-third-fourth -> top-next-third-fourth-top-next
	push(fr, top)
	push(fr, fourth)
	push(fr, third)
	push(fr, next)
	push(fr, top)
	return 1
}

// 0x5F SWAP swap top two items on stack
func doSwap(fr *frames.Frame, _ int64) int {
	top := pop(fr)
	next := pop(fr)
	push(fr, top)
	push(fr, next)
	return 1
}

// 0x60 IADD integer addition, push result
func doIadd(fr *frames.Frame, _ int64) int {
	i2 := pop(fr).(int64)
	i1 := pop(fr).(int64)
	sum := int32(i1) + int32(i2)
	push(fr, int64(sum))
	return 1
}

// 0x61 LADD integer addition, push result
func doLadd(fr *frames.Frame, _ int64) int {
	i2 := pop(fr).(int64)
	i1 := pop(fr).(int64)
	push(fr, i1+i2)
	return 1
}

// 0x62 FADD float addition, push result
func doFadd(fr *frames.Frame, _ int64) int {
	rhs := float32(pop(fr).(float64))
	lhs := float32(pop(fr).(float64))
	push(fr, float64(lhs+rhs))
	return 1
}

// 0x63 DADD double addition, push result
func doDadd(fr *frames.Frame, _ int64) int {
	rhs := pop(fr).(float64)
	lhs := pop(fr).(float64)
	push(fr, lhs+rhs)
	return 1
}

// Ox64 ISUB subtract subtract TOS-1 from TOS
func doIsub(fr *frames.Frame, _ int64) int {
	i2 := pop(fr).(int64)
	i1 := pop(fr).(int64)
	diff := int32(i1) - int32(i2)
	push(fr, int64(diff))
	return 1
}

// 0x65 LSUB subtract subtract TOS-1 from TOS
func doLsub(fr *frames.Frame, _ int64) int {
	i2 := pop(fr).(int64)
	i1 := pop(fr).(int64)
	diff := i1 - i2
	push(fr, diff)
	return 1
}

// 0x66 FSUB subtract float at TOS from TOS-1
func doFsub(fr *frames.Frame, _ int64) int {
	rhs := float32(pop(fr).(float64))
	lhs := float32(pop(fr).(float64))
	push(fr, float64(lhs-rhs))
	return 1
}

// 0x67 DSUB subtract double at TOS from TOS-1
func doDsub(fr *frames.Frame, _ int64) int {
	rhs := pop(fr).(float64)
	lhs := pop(fr).(float64)
	push(fr, lhs-rhs)
	return 1
}

// 0x68 IMUL multiply two int32s
func doImul(fr *frames.Frame, _ int64) int {
	i2 := int32(pop(fr).(int64))
	i1 := int32(pop(fr).(int64))
	product := multiply(i1, i2)
	push(fr, int64(product))
	return 1
}

// 0x68 LMUL multiply two int64s, i.e. longs
func doLmul(fr *frames.Frame, _ int64) int {
	i2 := pop(fr).(int64)
	i1 := pop(fr).(int64)
	product := multiply(i1, i2)
	push(fr, product)
	return 1
}

// 0x6A FMUL multiply floats
func doFmul(fr *frames.Frame, _ int64) int {
	rhs := float32(pop(fr).(float64))
	lhs := float32(pop(fr).(float64))
	push(fr, float64(lhs*rhs))
	return 1
}

// 0x6B DMUL multiply doubles
func doDmul(fr *frames.Frame, _ int64) int {
	rhs := pop(fr).(float64)
	lhs := pop(fr).(float64)
	push(fr, lhs*rhs)
	return 1
}

// 0x6C, 0x6D IDIV, LDIV divide TOS into TOS-1
func doIdiv(fr *frames.Frame, _ int64) int {
	val1 := pop(fr).(int64) // divisor
	val2 := pop(fr).(int64) // dividend
	if val1 == 0 {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errInfo := fmt.Sprintf("IDIV or LDIV: division by zero -- %d/0", val2)
		if globals.GetGlobalRef().StrictJDK { // use the HotSpot JDK's error message instead of ours
			errInfo = "/ by zero"
		}
		errMsg := fmt.Sprintf("in %s.%s %s",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, errInfo)
		status := exceptions.ThrowEx(excNames.ArithmeticException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		} else {
			// Make the current frame the caught exception frame.
			fs := fr.FrameStack
			fr = fs.Front().Value.(*frames.Frame)
			return 0 // PC is already set up so indicate that to caller.
		}
	} else {
		res := int32(val2) / int32(val1)
		push(fr, int64(res))
	}
	return 1
}

func doLdiv(fr *frames.Frame, _ int64) int {
	val1 := pop(fr).(int64) // divisor
	val2 := pop(fr).(int64) // dividend
	if val1 == 0 {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errInfo := fmt.Sprintf("IDIV or LDIV: division by zero -- %d/0", val2)
		if globals.GetGlobalRef().StrictJDK { // use the HotSpot JDK's error message instead of ours
			errInfo = "/ by zero"
		}
		errMsg := fmt.Sprintf("in %s.%s %s",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, errInfo)
		status := exceptions.ThrowEx(excNames.ArithmeticException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		} else {
			// Make the current frame the caught exception frame.
			fs := fr.FrameStack
			fr = fs.Front().Value.(*frames.Frame)
			return 0 // PC is already set up so indicate that to caller.
		}
	} else {
		push(fr, val2/val1)
	}
	return 1
}

// 0x6E FDIV floating-point division
func doFdiv(fr *frames.Frame, _ int64) int {
	val1 := float32(pop(fr).(float64)) // divisor
	val2 := float32(pop(fr).(float64)) // dividend
	push(fr, float64(val2/val1))
	return 1
}

// 0x6F DDIV double-precision floating-point division
func doDdiv(fr *frames.Frame, _ int64) int {
	val1 := pop(fr).(float64) // divisor
	val2 := pop(fr).(float64) // dividend
	push(fr, val2/val1)
	return 1
}

// 0x70, 0x71 IREM, LREM get remainder of integer division
func doIrem(fr *frames.Frame, _ int64) int {
	val2 := pop(fr).(int64)
	val1 := pop(fr).(int64)
	if val2 == 0 {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errInfo := fmt.Sprintf("IREM or LREM: division by zero -- %d/0", val2)
		if globals.GetGlobalRef().StrictJDK { // use the HotSpot JDK's error message instead of ours
			errInfo = "/ by zero"
		}
		errMsg := fmt.Sprintf("in %s.%s %s",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, errInfo)
		status := exceptions.ThrowEx(excNames.ArithmeticException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		} else {
			// Make the current frame the caught exception frame.
			fs := fr.FrameStack
			fr = fs.Front().Value.(*frames.Frame)
			return 0 // PC is already set up so indicate that to caller.
		}
	} else {
		res := int32(val1) % int32(val2)
		push(fr, int64(res))
	}
	return 1
}

func doLrem(fr *frames.Frame, _ int64) int {
	val2 := pop(fr).(int64)
	val1 := pop(fr).(int64)
	if val2 == 0 {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errInfo := fmt.Sprintf("IREM or LREM: division by zero -- %d/0", val2)
		if globals.GetGlobalRef().StrictJDK { // use the HotSpot JDK's error message instead of ours
			errInfo = "/ by zero"
		}
		errMsg := fmt.Sprintf("in %s.%s %s",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, errInfo)
		status := exceptions.ThrowEx(excNames.ArithmeticException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		} else {
			// Make the current frame the caught exception frame.
			fs := fr.FrameStack
			fr = fs.Front().Value.(*frames.Frame)
			return 0 // PC is already set up so indicate that to caller.
		}
	} else {
		res := val1 % val2
		push(fr, res)
	}
	return 1
}

// 0x72 FREM get remainder of floating-point division
func doFrem(fr *frames.Frame, _ int64) int {
	val2 := float32(pop(fr).(float64))
	val1 := float32(pop(fr).(float64))
	push(fr, float64(float32(math.Mod(float64(val1), float64(val2)))))
	return 1
}

// 0x73 DREM get remainder of double-precision floating-point division
func doDrem(fr *frames.Frame, _ int64) int {
	val2 := pop(fr).(float64)
	val1 := pop(fr).(float64)
	push(fr, math.Mod(val1, val2))
	return 1
}

// 0x74 INEG negate integer at TOS
func doIneg(fr *frames.Frame, _ int64) int {
	val := pop(fr).(int64)
	push(fr, int64(-int32(val)))
	return 1
}

// 0x75 LNEG negate long at TOS
func doLneg(fr *frames.Frame, _ int64) int {
	val := pop(fr).(int64)
	push(fr, -val)
	return 1
}

// 0x76 FNEG negate float at TOS
func doFneg(fr *frames.Frame, _ int64) int {
	val := float32(pop(fr).(float64))
	push(fr, float64(-val))
	return 1
}

// 0x77 DNEG negate double at TOS
func doDneg(fr *frames.Frame, _ int64) int {
	val := pop(fr).(float64)
	push(fr, -val)
	return 1
}

// 0x78, 0x79 ISHL, LSHL shift int/long to the left
func doIshl(fr *frames.Frame, _ int64) int {
	shiftBy := pop(fr).(int64)
	ushiftBy := uint32(shiftBy) & 0x1F // 0-31 bits for int per JVM
	val1 := pop(fr).(int64)
	val2 := int32(val1) << ushiftBy
	push(fr, int64(val2))
	return 1
}

func doLshl(fr *frames.Frame, _ int64) int {
	shiftBy := pop(fr).(int64)
	ushiftBy := uint64(shiftBy) & 0x3F // 0-63 bits for long per JVM
	val1 := pop(fr).(int64)
	val2 := val1 << ushiftBy
	push(fr, val2)
	return 1
}

// 0x7A, 0x7B ISHR, LSHR shift int/long to the right
func doIshr(fr *frames.Frame, _ int64) int {
	shiftBy := pop(fr).(int64)
	val1 := pop(fr).(int64)
	val2 := int32(val1) >> (shiftBy & 0x1F)
	push(fr, int64(val2))
	return 1
}

func doLshr(fr *frames.Frame, _ int64) int {
	shiftBy := pop(fr).(int64)
	val1 := pop(fr).(int64)
	val2 := val1 >> (shiftBy & 0x3F)
	push(fr, val2)
	return 1
}

// 0x7C IUSHR unsigned shift right of int (32 bits)
func doIushr(fr *frames.Frame, _ int64) int {
	shiftBy := pop(fr).(int64)
	value := pop(fr).(int64)
	shiftedVal := int64(uint32(value) >> (shiftBy & 0x1F))
	push(fr, shiftedVal)
	return 1
}

// 0x7D LUSHR unsigned shift right of long (64 bits)
func doLushr(fr *frames.Frame, _ int64) int {
	shiftBy := pop(fr).(int64)
	value := pop(fr).(int64)
	shiftedVal := int64(uint64(value) >> (shiftBy & 0x3F))
	push(fr, shiftedVal)
	return 1
}

// 0x7E, 0x7F IAND, LAND logical AND of two ints/longs, push result
func doIand(fr *frames.Frame, _ int64) int {
	val1 := pop(fr).(int64)
	val2 := pop(fr).(int64)
	push(fr, int64(int32(val1)&int32(val2)))
	return 1
}

func doLand(fr *frames.Frame, _ int64) int {
	val1 := pop(fr).(int64)
	val2 := pop(fr).(int64)
	push(fr, val1&val2)
	return 1
}

// 0x80, 0x81 IOR, LOR logical OR of two ints/longs, push result
func doIor(fr *frames.Frame, _ int64) int {
	val1 := pop(fr).(int64)
	val2 := pop(fr).(int64)
	push(fr, int64(int32(val1)|int32(val2)))
	return 1
}

func doLor(fr *frames.Frame, _ int64) int {
	val1 := pop(fr).(int64)
	val2 := pop(fr).(int64)
	push(fr, val1|val2)
	return 1
}

// 0x82, 0x83 IXOR, LXOR logical XOR of two ints/longs, push result
func doIxor(fr *frames.Frame, _ int64) int {
	val1 := pop(fr).(int64)
	val2 := pop(fr).(int64)
	push(fr, int64(int32(val1)^int32(val2)))
	return 1
}

func doLxor(fr *frames.Frame, _ int64) int {
	val1 := pop(fr).(int64)
	val2 := pop(fr).(int64)
	push(fr, val1^val2)
	return 1
}

// 0x84 IINC increment int variable
func doIinc(fr *frames.Frame, _ int64) int {
	var index int
	var increment int64
	var PCtoSkip int
	if fr.WideInEffect { // if wide is in effect, index  and increment are two bytes wide, otherwise one byte each
		index = (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2])
		increment = int64(fr.Meth[fr.PC+3])*256 + int64(fr.Meth[fr.PC+4])
		PCtoSkip = 4
		fr.WideInEffect = false
	} else {
		index = int(fr.Meth[fr.PC+1])
		increment = byteToInt64(fr.Meth[fr.PC+2])
		PCtoSkip = 2
	}

	// shoehorn the result into Java's 32-bit int
	orig := fr.Locals[index].(int64)
	fr.Locals[index] = int64(int32(orig) + int32(increment))
	return PCtoSkip + 1
}

// 0x86 I2F convert int to float
func doI2f(fr *frames.Frame, _ int64) int {
	intVal := pop(fr).(int64)
	push(fr, float64(float32(int32(intVal))))
	return 1
}

// 0x87 I2D convert int to double
func doI2d(fr *frames.Frame, _ int64) int {
	intVal := pop(fr).(int64)
	push(fr, float64(int32(intVal)))
	return 1
}

// 0x89 L2F convert long to float
func doL2f(fr *frames.Frame, _ int64) int {
	longVal := pop(fr).(int64)
	push(fr, float64(float32(longVal)))
	return 1
}

// 0x8A L2D convert long to double
func doL2d(fr *frames.Frame, _ int64) int {
	longVal := pop(fr).(int64)
	push(fr, float64(longVal))
	return 1
}

// 0x8B F2I convert float to int
func doF2i(fr *frames.Frame, _ int64) int {
	floatVal := float32(pop(fr).(float64))
	var res int32
	if math.IsNaN(float64(floatVal)) {
		res = 0
	} else if math.IsInf(float64(floatVal), 1) {
		res = math.MaxInt32
	} else if math.IsInf(float64(floatVal), -1) {
		res = math.MinInt32
	} else {
		res = int32(floatVal)
	}
	push(fr, int64(res))
	return 1
}

// 0x8C F2L convert float to long
func doF2l(fr *frames.Frame, _ int64) int {
	floatVal := float32(pop(fr).(float64))
	var res int64
	if math.IsNaN(float64(floatVal)) {
		res = 0
	} else if math.IsInf(float64(floatVal), 1) {
		res = math.MaxInt64
	} else if math.IsInf(float64(floatVal), -1) {
		res = math.MinInt64
	} else {
		res = int64(floatVal)
	}
	push(fr, res)
	return 1
}

// 0x8D F2D convert float to double
func doF2d(fr *frames.Frame, _ int64) int {
	floatVal := float32(pop(fr).(float64))
	push(fr, float64(floatVal))
	return 1
}

// 0x8E D2I convert double to int
func doD2i(fr *frames.Frame, _ int64) int {
	doubleVal := pop(fr).(float64)
	var res int32
	if math.IsNaN(doubleVal) {
		res = 0
	} else if math.IsInf(doubleVal, 1) {
		res = math.MaxInt32
	} else if math.IsInf(doubleVal, -1) {
		res = math.MinInt32
	} else {
		res = int32(doubleVal)
	}
	push(fr, int64(res))
	return 1
}

// 0x8F D2L convert double to long
func doD2l(fr *frames.Frame, _ int64) int {
	doubleVal := pop(fr).(float64)
	var res int64
	if math.IsNaN(doubleVal) {
		res = 0
	} else if math.IsInf(doubleVal, 1) {
		res = math.MaxInt64
	} else if math.IsInf(doubleVal, -1) {
		res = math.MinInt64
	} else {
		res = int64(doubleVal)
	}
	push(fr, res)
	return 1
}

// 0x90 D2F convert double to float
func doD2f(fr *frames.Frame, _ int64) int {
	doubleVal := pop(fr).(float64)
	push(fr, float64(float32(doubleVal)))
	return 1
}

// 0x91 I2B convert int to byte, preserving sign
func doI2b(fr *frames.Frame, _ int64) int {
	intVal := pop(fr).(int64)
	byteVal := int8(intVal)
	push(fr, int64(byteVal))
	return 1
}

// 0x92 I2C convert int to 16-bit char
func doI2c(fr *frames.Frame, _ int64) int {
	intVal := pop(fr).(int64)
	charVal := uint16(intVal) // Java chars are 16-bit unsigned values
	push(fr, int64(charVal))
	return 1
}

// 0x93 I2S convert int to short (16-bits)
func doI2s(fr *frames.Frame, _ int64) int {
	intVal := pop(fr).(int64)
	shortVal := int16(intVal) // Java shorts are 16-bit signed values
	push(fr, int64(shortVal))
	return 1
}

// 0x94 LCMP (compare two longs, push int -1, 0, or 1, depending on result)
func doLcmp(fr *frames.Frame, _ int64) int {
	value2 := pop(fr).(int64)
	value1 := pop(fr).(int64)
	if value1 == value2 {
		push(fr, int64(0))
	} else if value1 > value2 {
		push(fr, int64(1))
	} else {
		push(fr, int64(-1))
	}
	return 1
}

// 0x95, 0x96 FCMPL, FCMPG float comparison differing only in handling NaN
func doFcmpl(fr *frames.Frame, _ int64) int {
	value2 := float32(pop(fr).(float64))
	value1 := float32(pop(fr).(float64))
	if math.IsNaN(float64(value1)) || math.IsNaN(float64(value2)) {
		push(fr, int64(-1))
	} else if value1 > value2 {
		push(fr, int64(1))
	} else if value1 < value2 {
		push(fr, int64(-1))
	} else {
		push(fr, int64(0))
	}
	return 1
}

func doFcmpg(fr *frames.Frame, _ int64) int {
	value2 := float32(pop(fr).(float64))
	value1 := float32(pop(fr).(float64))
	if math.IsNaN(float64(value1)) || math.IsNaN(float64(value2)) {
		push(fr, int64(1))
	} else if value1 > value2 {
		push(fr, int64(1))
	} else if value1 < value2 {
		push(fr, int64(-1))
	} else {
		push(fr, int64(0))
	}
	return 1
}

// 0x97, 0x98 DCMPL, DCMPG double comparison differing only in handling NaN
func doDcmpl(fr *frames.Frame, _ int64) int {
	value2 := pop(fr).(float64)
	value1 := pop(fr).(float64)
	if math.IsNaN(value1) || math.IsNaN(value2) {
		push(fr, int64(-1))
	} else if value1 > value2 {
		push(fr, int64(1))
	} else if value1 < value2 {
		push(fr, int64(-1))
	} else {
		push(fr, int64(0))
	}
	return 1
}

func doDcmpg(fr *frames.Frame, _ int64) int {
	value2 := pop(fr).(float64)
	value1 := pop(fr).(float64)
	if math.IsNaN(value1) || math.IsNaN(value2) {
		push(fr, int64(1))
	} else if value1 > value2 {
		push(fr, int64(1))
	} else if value1 < value2 {
		push(fr, int64(-1))
	} else {
		push(fr, int64(0))
	}
	return 1
}

// 0x99 IFEQ pop int, if it's == 0, go to the jump location
func doIfeq(fr *frames.Frame, _ int64) int {
	// bools are treated in the JVM as ints, so convert here if bool;
	// otherwise, values should be int64's
	popValue := pop(fr)
	value := convertInterfaceToInt64(popValue)
	if value == 0 {
		jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
		return int(jumpTo)
	} else {
		return 3
	}
}

// 0x9A IFNE pop int, if it's != 0, go to the jump location
func doIfne(fr *frames.Frame, _ int64) int {
	// bools are treated in the JVM as ints, so convert here if bool;
	// otherwise, values should be int64's
	popValue := pop(fr)
	value := convertInterfaceToInt64(popValue)
	if value != 0 {
		jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
		return int(jumpTo)
	} else {
		return 3
	}
}

// 0x9B IFLT pop int, if it's < 0, go to the jump location
func doIflt(fr *frames.Frame, _ int64) int {
	// bools are treated in the JVM as ints, so convert here if bool;
	// otherwise, values should be int64's
	popValue := pop(fr)
	value := convertInterfaceToInt64(popValue)
	if value < 0 {
		jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
		return int(jumpTo)
	} else {
		return 3
	}
}

// 0x9C IFGE pop int, if it's >= 0, go to the jump location
func doIfge(fr *frames.Frame, _ int64) int {
	// bools are treated in the JVM as ints, so convert here if bool;
	// otherwise, values should be int64's
	popValue := pop(fr)
	value := convertInterfaceToInt64(popValue)
	if value >= 0 {
		jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
		return int(jumpTo)
	} else {
		return 3
	}
}

// 0x9D IFGT pop int, if it's > 0, go to the jump location
func doIfgt(fr *frames.Frame, _ int64) int {
	// bools are treated in the JVM as ints, so convert here if bool;
	// otherwise, values should be int64's
	popValue := pop(fr)
	value := convertInterfaceToInt64(popValue)
	if value > 0 {
		jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
		return int(jumpTo)
	} else {
		return 3
	}
}

// 0x9E IFLE pop int, if it's <!>= 0, go to the jump location
func doIfle(fr *frames.Frame, _ int64) int {
	// bools are treated in the JVM as ints, so convert here if bool;
	// otherwise, values should be int64's
	popValue := pop(fr)
	value := convertInterfaceToInt64(popValue)
	if value <= 0 {
		jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
		return int(jumpTo)
	} else {
		return 3
	}
}

// 0x9F IF_ICMPEQ  jump if two popped ints are equal
func doIficmpeq(fr *frames.Frame, _ int64) int {
	popValue := pop(fr)
	val2 := convertInterfaceToInt64(popValue)
	popValue = pop(fr)
	val1 := convertInterfaceToInt64(popValue)
	if int32(val1) == int32(val2) { // if comp succeeds, next 2 bytes hold instruction index
		jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
		return int(jumpTo)
	} else {
		return 3 // 2 for the jumpTo + 1 for next bytecode
	}
}

// 0xA0 IF_ICMPNE jump if two popped ints are not equal
func doIficmpne(fr *frames.Frame, _ int64) int {
	popValue := pop(fr)
	val2 := convertInterfaceToInt64(popValue)
	popValue = pop(fr)
	val1 := convertInterfaceToInt64(popValue)
	if int32(val1) != int32(val2) { // if comp fails, next 2 bytes hold instruction index
		jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
		return int(jumpTo)
	} else {
		return 3 // 2 for the jumpTo + 1 for next bytecode
	}
}

// 0xA1 IF_ICMPLT Compare popped ints for <
func doIficmplt(fr *frames.Frame, _ int64) int {
	popValue := pop(fr)
	val2 := convertInterfaceToInt64(popValue)
	popValue = pop(fr)
	val1 := convertInterfaceToInt64(popValue)
	if val1 < val2 { // if comp succeeds, next 2 bytes hold instruction index
		jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
		return int(jumpTo)
	} else {
		return 3 // 2 bytes for the jumpTo + 1 byte to next bytecode
	}
}

// 0xA2 IF_ICMPGE Compare ints for >=
func doIficmpge(fr *frames.Frame, _ int64) int {
	popValue := pop(fr)
	val2 := convertInterfaceToInt64(popValue)
	popValue = pop(fr)
	val1 := convertInterfaceToInt64(popValue)
	if val1 >= val2 { // if comp succeeds, next 2 bytes hold instruction index
		jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
		return int(jumpTo)
	} else {
		return 3 // 2 bytes for the jumpTo + 1 byte to next bytecode
	}
}

// 0xA3 IF_ICMPGT  jump if popped int > int at TOS
func doIficmpgt(fr *frames.Frame, _ int64) int {
	popValue := pop(fr)
	val2 := convertInterfaceToInt64(popValue)
	popValue = pop(fr)
	val1 := convertInterfaceToInt64(popValue)
	if int32(val1) > int32(val2) { // if comp succeeds, next 2 bytes hold instruction index
		jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
		return int(jumpTo)
	} else {
		return 3 // 2 for the jumpTo + 1 for next bytecode
	}
}

// 0xA4 IF_ICMPLE  jump if popped int <= int at TOS
func doIficmple(fr *frames.Frame, _ int64) int {
	popValue := pop(fr)
	val2 := convertInterfaceToInt64(popValue)
	popValue = pop(fr)
	val1 := convertInterfaceToInt64(popValue)
	if int32(val1) <= int32(val2) { // if comp succeeds, next 2 bytes hold instruction index
		jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
		return int(jumpTo)
	} else {
		return 3 // 2 for the jumpTo + 1 for next bytecode
	}
}

// 0xA5 IF_ACMPEQ  jump if two addresses are equal
func doIfacmpeq(fr *frames.Frame, _ int64) int {
	val2 := pop(fr)
	val1 := pop(fr)
	if val1 == val2 { // if comp succeeds, next 2 bytes hold instruction index
		jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
		return int(jumpTo)
	} else {
		return 3 // 2 for the jumpTo + 1 for next bytecode
	}
}

// 0xA6 IF_ACMPNE  jump if two addresses are equal
func doIfacmpne(fr *frames.Frame, _ int64) int {
	val2 := pop(fr)
	val1 := pop(fr)
	if val1 != val2 { // if comp fails, next 2 bytes hold instruction index
		jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
		return int(jumpTo)
	} else {
		return 3 // 2 for the jumpTo + 1 for next bytecode
	}
}

// 0xA7 GOTO unconditional jump within method
func doGoto(fr *frames.Frame, _ int64) int {
	jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
	return int(jumpTo) // note the value can be negative to jump to earlier bytecode
}

// 0xA8 JSR jump to a bytecode in the method at jumpTo bytes
func doJsr(fr *frames.Frame, _ int64) int {
	jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
	return int(jumpTo)
}

// 0xA9 RET return by jumping to a return address stored in a local
func doRet(fr *frames.Frame, _ int64) int {
	var index int64
	if fr.WideInEffect { // if wide is in effect, index is two bytes wide, otherwise one byte
		index = (byteToInt64(fr.Meth[fr.PC+1]) * 256) + byteToInt64(fr.Meth[fr.PC+2])
		fr.WideInEffect = false
	} else {
		index = byteToInt64(fr.Meth[fr.PC+1])
	}
	newPC := fr.Locals[index].(int64)
	return int(newPC)
}

// 0xAA TABLESWITCH switch based on table of offsets
func doTableswitch(fr *frames.Frame, _ int64) int {
	// https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-6.html#jvms-6.5.tableswitch
	basePC := fr.PC // where we are when the processing begins

	paddingBytes := 4 - ((fr.PC + 1) % 4)
	if paddingBytes == 4 {
		paddingBytes = 0
	}
	fr.PC += paddingBytes

	defaultJump := types.FourBytesToInt64( // the jump if the value is not in the table
		fr.Meth[fr.PC+1], fr.Meth[fr.PC+2], fr.Meth[fr.PC+3], fr.Meth[fr.PC+4])
	fr.PC += 4
	lowValue := types.FourBytesToInt64( // the lowest value in the table
		fr.Meth[fr.PC+1], fr.Meth[fr.PC+2], fr.Meth[fr.PC+3], fr.Meth[fr.PC+4])
	fr.PC += 4
	highValue := types.FourBytesToInt64( // the highest value in the table
		fr.Meth[fr.PC+1], fr.Meth[fr.PC+2], fr.Meth[fr.PC+3], fr.Meth[fr.PC+4])
	fr.PC += 4

	index := pop(fr).(int64) // the value we're looking to match

	// Compute PC for jump.
	jumpOffset := 0 //
	for value := lowValue; value <= highValue; value++ {
		if value == index {
			fr.PC += jumpOffset
			jumpPC := types.FourBytesToInt64(
				fr.Meth[fr.PC+1], fr.Meth[fr.PC+2], fr.Meth[fr.PC+3], fr.Meth[fr.PC+4])
			fr.PC = basePC
			return int(jumpPC)
		}
		jumpOffset += 4
	}

	// Default case.
	fr.PC = basePC
	return int(defaultJump)
}

// 0xAB LOOKUPSWITCH switch using lookup table
func doLookupswitch(fr *frames.Frame, _ int64) int {
	// https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-6.html#jvms-6.5.lookupswitch
	basePC := fr.PC // where we are when the processing begins

	paddingBytes := 4 - ((fr.PC + 1) % 4)
	if paddingBytes == 4 {
		paddingBytes = 0
	}
	fr.PC += paddingBytes

	// get the jump size for the default branch
	defaultJump := int64(binary.BigEndian.Uint32(
		[]byte{fr.Meth[fr.PC+1], fr.Meth[fr.PC+2], fr.Meth[fr.PC+3], fr.Meth[fr.PC+4]}))
	fr.PC += 4

	// how many branches in this switch (other than default)
	npairs := binary.BigEndian.Uint32(
		[]byte{fr.Meth[fr.PC+1], fr.Meth[fr.PC+2], fr.Meth[fr.PC+3], fr.Meth[fr.PC+4]})
	fr.PC += 4

	jumpTable := make(map[int64]int)
	for i := 0; i < int(npairs); i++ {
		// get the jump size for each case branch
		caseValue := types.FourBytesToInt64(
			fr.Meth[fr.PC+1], fr.Meth[fr.PC+2], fr.Meth[fr.PC+3], fr.Meth[fr.PC+4])
		fr.PC += 4
		jumpOffset := types.FourBytesToInt64(fr.Meth[fr.PC+1], fr.Meth[fr.PC+2], fr.Meth[fr.PC+3], fr.Meth[fr.PC+4])
		fr.PC += 4
		jumpTable[caseValue] = int(jumpOffset)
	}

	// now get the value we're switching on and find the distance to jump
	fr.PC = basePC
	key := pop(fr).(int64)
	jumpDistance, present := jumpTable[key]
	if present {
		// Validate jumpDistance.
		// Do not return a jump distance that would crash the main line interpreter.
		if fr.PC+jumpDistance < 0 {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			fqn := fr.ClName + "." + fr.MethName + fr.MethType
			errMsg := fmt.Sprintf("LOOKUPSWITCH: In %s, impossible jump, fr.PC=%d, jumpDistance=%d", fqn, fr.PC, jumpDistance)
			status := exceptions.ThrowEx(excNames.VirtualMachineError, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}
		return jumpDistance
	} else {
		return int(defaultJump)
	}
}

// 0xAC - 0xB0 IRETURN, LRETURN, DRETURN, FRETURN, ARETURN
// return a value from method call. Important note:
// This implementation pops off the current frame and tells the
// interpreter loop to resume execution in the previous frame.
func doIreturn(fr *frames.Frame, _ int64) int {
	// If the current method is synchronized, unlock the object.
	if fr.ObjSync != nil {
		_ = fr.ObjSync.ObjUnlock(int32(fr.Thread))
		if globals.TraceInst {
			traceInfo := fmt.Sprintf("\tdoIreturn: Unlocked object %s",
				object.GoStringFromStringPoolIndex(fr.ObjSync.KlassName))
			trace.Trace(traceInfo)
		}
	}
	valToReturn := pop(fr)
	f := fr.FrameStack.Front().Next().Value.(*frames.Frame)
	push(f, valToReturn)
	fr.FrameStack.Remove(fr.FrameStack.Front())
	return 0
}

// 0xB1 RETURN return from void method
func doReturn(fr *frames.Frame, _ int64) int {
	// If the current method is synchronized, unlock the object.
	if fr.ObjSync != nil {
		_ = fr.ObjSync.ObjUnlock(int32(fr.Thread))
		if globals.TraceInst {
			traceInfo := fmt.Sprintf("\tdoReturn: Unlocked object %s",
				object.GoStringFromStringPoolIndex(fr.ObjSync.KlassName))
			trace.Trace(traceInfo)
		}
	}
	fr.FrameStack.Remove(fr.FrameStack.Front())
	return 0
}

// 0xB2 GETSTATIC
// statics are stored in the minimal java/lang/Class mirror of the loaded class
// so we get that mirror from the statics table (using the class name), then
// get its table of static fields, and finally read the value there (using a lock)
func doGetStatic(fr *frames.Frame, _ int64) int {
	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2]) // next 2 bytes point to CP entry
	CP := fr.CP.(*classloader.CPool)
	CPentry := CP.CpIndex[CPslot] // value is checked in codeCheck.go

	// get the field entry
	field := CP.FieldRefs[CPentry.Slot]
	className := field.ClName
	fieldName := field.FldName
	if globals.TraceInst {
		EmitTraceFieldID("GETSTATIC", className+"."+fieldName)
	}

	// was this static field previously loaded? Is so, get its location and move on.
	prevLoaded, ok := statics.QueryStatic(className, field.FldName)
	if !ok { // if the field is not already loaded, then
		// the class has not been instantiated, so instantiate the class
		_, err := InstantiateClass(className, fr.FrameStack)
		if err == nil {
			prevLoaded, ok = statics.QueryStatic(className, field.FldName)
		} else {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := fmt.Sprintf("GETSTATIC: could not load class %s", className)
			status := exceptions.ThrowEx(excNames.ClassNotFoundException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}
	}

	// if the field can't be found even after instantiating the
	// containing class, something is wrong so get out of here.
	if !ok {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("GETSTATIC: could not find static field %s in class %s"+
			"\n", className+"."+fieldName, className)
		status := exceptions.ThrowEx(excNames.NoSuchFieldException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
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

// 0xB3 PUTSTATIC
func doPutStatic(fr *frames.Frame, _ int64) int {
	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2])
	CP := fr.CP.(*classloader.CPool)
	CPentry := CP.CpIndex[CPslot] // value is checked in codeCheck.go

	// get the field entry
	field := CP.FieldRefs[CPentry.Slot]
	className := field.ClName
	fieldName := field.FldName
	fieldName = className + "." + fieldName
	if globals.TraceInst {
		EmitTraceFieldID("PUTSTATIC", fieldName)
	}

	// was this static field previously loaded? Is so, get its location and move on.
	prevLoaded, ok := statics.QueryStatic(className, field.FldName)
	if !ok { // if field is not already loaded, then
		if globals.TraceInst {
			msg := fmt.Sprintf("doPutStatic: Field was not previously loaded: %s", fieldName)
			trace.Trace(msg)
		}
		// the class has not been instantiated, so
		// instantiate the class
		_, err := InstantiateClass(className, fr.FrameStack)
		if err == nil {
			if globals.TraceInst {
				msg := fmt.Sprintf("doPutStatic: Loaded class %s", className)
				trace.Trace(msg)
			}
			prevLoaded, ok = statics.QueryStatic(className, field.FldName)
		} else {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := fmt.Sprintf("PUTSTATIC: could not load class %s", className)
			trace.Error(errMsg)
			return ERROR_OCCURRED
		}
	} else {
		if globals.TraceInst {
			msg := fmt.Sprintf("doPutStatic: Field was previously loaded: %s", fieldName)
			trace.Trace(msg)
		}
	}

	// if the field can't be found even after instantiating the
	// containing class, something is wrong so get out of here.
	if !ok {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("PUTSTATIC: could not find static field %s.%s", className, fieldName)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	var value interface{}
	switch prevLoaded.Type {
	// We want all forms normalized to int64
	case types.Bool:
		value = pop(fr).(int64) & 0x01
		statics.AddStatic(fieldName, statics.Static{
			Type:  prevLoaded.Type,
			Value: value,
		})
	case types.Char, types.Short, types.Int, types.Long:
		value = pop(fr).(int64)
		statics.AddStatic(fieldName, statics.Static{
			Type:  prevLoaded.Type,
			Value: value,
		})
	case types.Byte:
		var value int64
		v := pop(fr)
		switch v.(type) { // could be passed a byte or an integral type for a value
		case int64:
			value = v.(int64)
		case uint8:
			value = int64(v.(uint8))
		case types.JavaByte:
			value = int64(v.(types.JavaByte))
		}
		statics.AddStatic(fieldName, statics.Static{
			Type:  prevLoaded.Type,
			Value: value,
		})
	case types.Float, types.Double:
		value = pop(fr).(float64)
		statics.AddStatic(fieldName, statics.Static{
			Type:  prevLoaded.Type,
			Value: value,
		})

	default:
		// if it's not a primitive or a pointer to a class,
		// then it should be a pointer to an object or to
		// a loaded class
		value = pop(fr)
		if value == nil {
			value = object.Null
		}
		switch value.(type) {
		case *object.Object:
			statics.AddStatic(fieldName, statics.Static{
				Type:  prevLoaded.Type,
				Value: value,
			})

		case *classloader.Klass:
			// convert to an *object.Object
			kPtr := value.(*classloader.Klass)
			obj := object.MakeEmptyObject()
			obj.KlassName = stringPool.GetStringIndex(&kPtr.Data.Name)
			objField := object.Field{
				Ftype:  "L" + kPtr.Data.Name + ";",
				Fvalue: kPtr,
			}

			obj.ThMutex.Lock()
			obj.FieldTable[fieldName] = objField
			obj.ThMutex.Unlock()

			statics.AddStatic(fieldName, statics.Static{
				Type:  prevLoaded.Type,
				Value: value,
			})
		default:
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := fmt.Sprintf("PUTSTATIC: field %s.%s, type unrecognized: %T %v", className, fieldName, value, value)
			trace.Error(errMsg)
			return ERROR_OCCURRED
		}
	}
	return 3 // 2 for the CP slot + 1 for next bytecode
}

// 0xB4 GETFIELD get field in a pointed-to-object
func doGetfield(fr *frames.Frame, _ int64) int {
	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2]) // next 2 bytes point to CP entry
	CP := fr.CP.(*classloader.CPool)
	fieldEntry := CP.CpIndex[CPslot]
	// we check that the pointed-to CP entry is a field reference in codeCheck.go

	// Get field name.
	fullFieldEntry := CP.FieldRefs[fieldEntry.Slot]
	fieldName := fullFieldEntry.FldName
	if globals.TraceVerbose {
		EmitTraceFieldID("GETFIELD", fieldName)
	}

	// Get object reference from stack.
	ref := pop(fr)
	switch ref.(type) {
	case *object.Object:
		break
	default:
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("GETFIELD: Invalid type of object ref: %T, fieldName: %s.%s", ref, fr.ClName, fieldName)
		status := exceptions.ThrowEx(excNames.IllegalArgumentException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	// Check reference for a nil pointer.
	if object.IsNull(ref) {
		errMsg := fmt.Sprintf("GETFIELD: Null object reference, fieldName: %s.%s", fr.ClName, fieldName)
		status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	// Extract field.
	obj := ref.(*object.Object)
	var fieldType string
	var fieldValue interface{}

	obj.ThMutex.RLock()
	objField, ok := obj.FieldTable[fieldName]
	obj.ThMutex.RUnlock()
	if !ok {
		errMsg := fmt.Sprintf("GETFIELD PC=%d: Missing field (%s) in FieldTable", fr.PC, fieldName)
		status := exceptions.ThrowEx(excNames.IllegalArgumentException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
	}

	/* JACOBIN-824: Temporarily commenting out the following lines
	obj := *ref.(*object.Object)
	var fieldType string
	var fieldValue interface{}
	var objField object.Field
	var ok bool

	if obj.KlassName == types.StringPoolThreadIndex {
		runnable := obj.FieldTable["target"].Fvalue.(*object.Object)
		objField, ok = runnable.FieldTable[fieldName]
	} else {
		objField, ok = obj.FieldTable[fieldName]
	}
	if !ok {
		errMsg := fmt.Sprintf("GETFIELD PC=%d: Missing field (%s) in FieldTable", fr.PC, fieldName)
		status := exceptions.ThrowEx(excNames.IllegalArgumentException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
	}
	*/

	fieldType = objField.Ftype
	if fieldType == types.StringIndex {
		fieldValue = stringPool.GetStringPointer(objField.Fvalue.(uint32))
	} else if fieldType == types.StringClassRef {
		// if the field type is String pointer and value is a byte array, convert it to a string
		switch objField.Fvalue.(type) {
		case []byte:
			fieldValue = object.StringObjectFromByteArray(objField.Fvalue.([]byte))
		case []types.JavaByte:
			fieldValue = object.StringObjectFromJavaByteArray(objField.Fvalue.([]types.JavaByte))
		}
	} else if types.IsArray(fieldType) {
		// if the field type is an array, other than a string, convert it to an object
		o := object.MakeEmptyObject()
		of := object.Field{Ftype: fieldType, Fvalue: objField.Fvalue}
		o.ThMutex.Lock()
		o.FieldTable["value"] = of
		o.ThMutex.Unlock()
		o.KlassName = stringPool.GetStringIndex(&of.Ftype)
		fieldValue = o
	} else if fieldType == "Ljava/lang/Object;" {
		// if it's a pointer to an Object and the value field is an array or slice, wrap the array in an Object
		v := reflect.ValueOf(objField.Fvalue)
		if v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
			var newFieldType string
			w := reflect.TypeOf(objField.Fvalue)
			elemType := w.Elem()
			switch elemType.Kind() {
			case reflect.Int8:
				newFieldType = types.ByteArray
			case reflect.Int64:
				newFieldType = types.IntArray
			case reflect.Float64:
				newFieldType = types.FloatArray // types.DoubleArray?
			default:
				newFieldType = "[Ljava/lang/Object;"
			}
			o := object.MakeEmptyObject()
			of := object.Field{Ftype: newFieldType, Fvalue: objField.Fvalue}
			o.ThMutex.Lock()
			o.FieldTable["value"] = of
			o.ThMutex.Unlock()
			klassName := "[Ljava/lang/Object;"
			o.KlassName = stringPool.GetStringIndex(&klassName) // "[Ljava/lang/Object;"
			fieldValue = o
		} else {
			fieldValue = objField.Fvalue
		}
	} else { // not an index to the string pool, nor a String pointer with a byte array
		// The field is another kind of Object.
		fieldValue = objField.Fvalue
	}

	push(fr, fieldValue)
	return 3 // 2 for CPslot + 1 for next bytecode
}

// 0xB5 PUTFIELD place value into an object's field
func doPutfield(fr *frames.Frame, _ int64) int {
	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2]) // next 2 bytes point to CP entry
	CP := fr.CP.(*classloader.CPool)
	fieldEntry := CP.CpIndex[CPslot]
	/* // Checked for in codeCheck.go Delete after Feb 1, 2026
	if fieldEntry.Type != classloader.FieldRef { // the pointed-to CP entry must be a field reference
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("PUTFIELD: Expected a field ref, but got %d in"+
			"location %d in method %s of class %s\n",
			fieldEntry.Type, fr.PC, fr.MethName, fr.ClName)
		status := exceptions.ThrowEx(excNames.NoSuchFieldException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
	} */

	value := pop(fr) // the value we're placing in the field
	ref := pop(fr)   // reference to the object we're updating

	switch ref.(type) {
	case *object.Object:
		// Handle the Object after this switch
	default:
		// *** unexpected type of ref ***
		errMsg := fmt.Sprintf("PUTFIELD: Expected an object ref, but observed type %T in "+
			"location %d in method %s of class %s, previously popped a value(type %T):\n%v\n",
			ref, fr.PC, fr.MethName, fr.ClName, value, value)
		status := exceptions.ThrowEx(excNames.InvalidTypeException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
	}

	// Get Object struct.
	obj := ref.(*object.Object)

	// if the value we're inserting is a reference to an array object, we need to modify it
	// to point directly to the array of primitives, rather than to the array object
	switch value.(type) {
	case *object.Object:
		if !object.IsNull(value.(*object.Object)) {
			v := value.(*object.Object)
			v.ThMutex.RLock()
			o, ok := v.FieldTable["value"]
			v.ThMutex.RUnlock()
			if ok && strings.HasPrefix(o.Ftype, types.Array) {
				v.ThMutex.RLock()
				value = v.FieldTable["value"].Fvalue
				v.ThMutex.RUnlock()
			}
		}
	}

	// otherwise look up the field name in the CP and find it in the FieldTable, then do the update
	obj.ThMutex.RLock()
	fieldTableLen := len(obj.FieldTable)
	obj.ThMutex.RUnlock()
	if fieldTableLen != 0 {
		fullFieldEntry := CP.FieldRefs[fieldEntry.Slot]
		fieldName := fullFieldEntry.FldName
		if globals.TraceVerbose {
			EmitTraceFieldID("PUTFIELD", fieldName)
		}

		obj.ThMutex.RLock()
		objField, ok := obj.FieldTable[fieldName]
		obj.ThMutex.RUnlock()
		if !ok {
			errMsg := fmt.Sprintf("PUTFIELD: In trying for a superclass field, %s is not present in object of class %s",
				fieldName, object.GoStringFromStringPoolIndex(obj.KlassName))
			status := exceptions.ThrowEx(excNames.NoSuchFieldException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}

		// PUTFIELD is not used to update statics. That's for PUTSTATIC to do.
		if strings.HasPrefix(objField.Ftype, types.Static) {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "PUTFIELD: invalid attempt to update a static variable"
			status := exceptions.ThrowEx(excNames.InvalidTypeException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}

		switch objField.Ftype {
		case types.Byte:
			objField.Fvalue = int64(int8(uint8(value.(int64))))
		case types.Int:
			objField.Fvalue = int64(int32(uint32(value.(int64))))
		case types.Long:
			objField.Fvalue = value.(int64)
		case types.Short:
			objField.Fvalue = int64(int16(uint16(value.(int64))))
		default:
			objField.Fvalue = value
		}
		obj.ThMutex.Lock()
		obj.FieldTable[fieldName] = objField
		obj.ThMutex.Unlock()
	}
	return 3 // 2 for CPslot + 1 for next bytecode
}

// 0xB6 INVOKEVIRTUAL
func doInvokeVirtual(fr *frames.Frame, _ int64) int {
	var err error
	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2]) // next 2 bytes point to CP entry
	CP := fr.CP.(*classloader.CPool)                                // codeCheck.go ensures that CPslot is a valid index to a methodRef

	// Get the method table entry for the FQN indicated in CP.
	className, methodName, methodType, fqn :=
		classloader.GetMethInfoFromCPmethref(CP, CPslot)
	mtEntry := classloader.GetMtableEntry(className + "." + methodName + methodType)
	if mtEntry.Meth == nil { // if the method is not in the method table, search classes or superclasses
		mtEntry, err = classloader.FetchMethodAndCP(className, methodName, methodType)
	}

	// Not found after a class-superclass search. Check the interfaces.
	if err != nil || mtEntry.Meth == nil { // the method is not in the superclasses, so check interfaces.
		// When a class implements an interface and inherits default methods (or doesn't override them),
		// the compiler generates INVOKEVIRTUAL
		klass := classloader.MethAreaFetch(className)
		if len(klass.Data.Interfaces) > 0 {
			for i := 0; i < len(klass.Data.Interfaces); i++ {
				index := uint32(klass.Data.Interfaces[i])
				interfaceName := *stringPool.GetStringPointer(index)
				mtEntry, err = locateInterfaceMeth(klass, fr, className, interfaceName, methodName, methodType)
				if mtEntry.Meth != nil {
					// found a match.
					break
				}
			} // end of search of interfaces if method has any

			// Any matches in the interfaces?
			if err != nil || mtEntry.Meth == nil {
				// method was not found in interfaces, so throw an exception
				globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
				errMsg := "INVOKEVIRTUAL: Class method not found: " + fqn
				status := exceptions.ThrowEx(excNames.NoSuchMethodException, errMsg, fr)
				if status != exceptions.Caught {
					return ERROR_OCCURRED // applies only if in test
				}
				return RESUME_HERE // caught
			}
		}
	}

	// if we got here, we have a method to call in mtEntry.Meth

	// if we have a native function (here, one implemented in golang, rather than Java),
	// then follow the JVM spec and push the objectRef and the parameters to the function
	// as parameters. Consult:
	// https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-6.html#jvms-6.5.invokevirtual

	if mtEntry.MType == 'G' { // so we have a golang function
		return invokeVirtualGfunction(fr, mtEntry, className, methodName, methodType)
	}

	// 	To resolve a J method (i.e., a Java method) for invokevirtual:
	//  - If it's a native Java function (written in C/C++), Jacobin does not support it.
	//  - Get the reference object from the stack.
	// 	- Try searching the reference object class and its superclass chain.
	// 	- If the method is not found, try the reference object class interface hierarchy (JVM spec 5.4.3.4).
	if mtEntry.MType == 'J' { // it's a Java function
		m := mtEntry.Meth.(classloader.JmEntry)
		if m.AccessFlags&classloader.ACC_NATIVE > 0 {
			// Native code
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKEVIRTUAL: Native method requested: " + fqn
			status := exceptions.ThrowEx(excNames.UnsupportedOperationException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}

		// The run-time class object is on the stack, below the method arguments.
		// To locate it, get the number of arguments for the method.
		nslots := len(util.ParseIncomingParamsFromMethTypeString(methodType))

		// Extract the reference object from the stack.
		refObj, ok := fr.OpStack[fr.TOS-nslots].(*object.Object)
		if !ok {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKEVIRTUAL: Stack reference object is nil"
			status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}

		// Get the reference object class name.
		clNameIdx := refObj.KlassName
		className = *(stringPool.GetStringPointer(clNameIdx))

		// === Method resolution ===
		// First, try superclass resolution.
		mtEntry, err = classloader.FetchMethodAndCP(className, methodName, methodType)
		if err != nil || mtEntry.Meth == nil {

			// That did not succeed. So, try for an interface default method.
			var ret any
			ret, mtEntry = searchForDefaultInterfaceFunction(className, methodName, methodType)
			if ret == nil {
				globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
				errMsg := "INVOKEVIRTUAL: Concreted class method not found: " + fqn
				status := exceptions.ThrowEx(excNames.NoSuchMethodError, errMsg, fr)
				if status != exceptions.Caught {
					return ERROR_OCCURRED // applies only if in test
				}
				return RESUME_HERE // caught
			}

			// Found an interface default method.
			className = ret.(string)
		}

		// Resolve to a G function?
		if mtEntry.MType == 'G' {
			return invokeVirtualGfunction(fr, mtEntry, className, methodName, methodType)
		}

		// It's a J function. Get its JmEntry.
		m = mtEntry.Meth.(classloader.JmEntry)
		fqn = className + "." + methodName + methodType

		// If an empty code segment, that's an error. Its probably abstract or an interface.
		// In this case, flag it as an AbstractMethodError.
		if len(m.Code) == 0 {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKEVIRTUAL: J class method code is empty: " + fqn
			status := exceptions.ThrowEx(excNames.AbstractMethodError, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}

		// Create the next frame to execute.
		nextFrame, err := createAndInitNewFrame(
			className, methodName, methodType, &m, true, fr)
		if err != nil {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKEVIRTUAL: Error creating frame in: " + fqn
			status := exceptions.ThrowEx(excNames.InvalidStackFrameException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}

		fr.PC += 3                         // 2 for PC slot, move to next bytecode before exiting
		fr.FrameStack.PushFront(nextFrame) // push the new frame
		return 0
	}
	return ERROR_OCCURRED // in theory, unreachable
}

// searchForDefaultInterfaceFunction searches for a default method in all interfaces
// implemented by className (and their superinterfaces), following JVM rules.
//
// JVM rules for default methods:
//  1. Start with all interfaces directly implemented by the class.
//  2. Recursively check superinterfaces.
//  3. If multiple default methods are found with the same signature, throw
//     IncompatibleClassChangeError (ambiguity).
//  4. Otherwise, return the most-specific default method (closest to the class).
//
// Returns the interface name and mtEntry if a default method is found, or nil if not.
//
// Reference: https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-5.html#jvms-5.4.3.4
//
// Note: this function is similar in many aspects to locateInterfaceMeth() in run.go.
// The two functions might eventually be integrated into one.
func searchForDefaultInterfaceFunction(
	className, methodName, methodType string,
) (any, classloader.MTentry) {

	var mtEntry classloader.MTentry

	// Fetch runtime class metadata
	klass := classloader.MethAreaFetch(className)
	if klass == nil || len(klass.Data.Interfaces) == 0 {
		// Class does not implement any interfaces
		return nil, mtEntry
	}

	// Track visited interfaces to prevent infinite loops
	visited := make(map[string]bool)

	// Recursive helper to search an interface and its superinterfaces
	var searchInterface func(intfName string) (string, classloader.MTentry)
	searchInterface = func(intfName string) (string, classloader.MTentry) {
		if visited[intfName] {
			return "", classloader.MTentry{}
		}
		visited[intfName] = true

		intfClass := classloader.MethAreaFetch(intfName)
		if intfClass == nil {
			return "", classloader.MTentry{}
		}

		// Step 3a: Build fully qualified method name and check method table
		fqn := intfName + "." + methodName + methodType
		mtEntry = classloader.GetMtableEntry(fqn)
		if mtEntry.Meth != nil {
			m := mtEntry.Meth.(classloader.JmEntry)
			if len(m.Code) > 0 {
				// Found a non-abstract method → candidate default method
				return intfName, mtEntry
			}
		}

		// Step 3b: Recurse into superinterfaces
		for _, idx := range intfClass.Data.Interfaces {
			superIntfName := *stringPool.GetStringPointer(uint32(idx))
			if resName, resMT := searchInterface(superIntfName); resName != "" {
				return resName, resMT
			}
		}

		return "", classloader.MTentry{}
	}

	// Track all candidates found
	candidates := make([]struct {
		intf string
		mt   classloader.MTentry
	}, 0)

	// Search all interfaces directly implemented by this class
	for _, idx := range klass.Data.Interfaces {
		intfName := *stringPool.GetStringPointer(uint32(idx))
		if resName, resMT := searchInterface(intfName); resName != "" {
			candidates = append(candidates, struct {
				intf string
				mt   classloader.MTentry
			}{resName, resMT})
		}
	}

	// JVM rule: handle multiple candidates
	if len(candidates) > 1 {
		// Conflict detected: multiple default methods with same signature
		errMsg := "INVOKEVIRTUAL: Ambiguous default method for " + className + "." + methodName + methodType
		status := exceptions.ThrowEx(excNames.IncompatibleClassChangeError, errMsg, nil)
		if status != exceptions.Caught {
			// Only relevant in test; return empty
			return nil, classloader.MTentry{}
		}
	}

	// Return the most-specific default method (first found)
	if len(candidates) == 1 {
		return candidates[0].intf, candidates[0].mt
	}

	// No default method found
	return nil, mtEntry
}

// Execute an INVOKEVIRTUAL G function.
func invokeVirtualGfunction(fr *frames.Frame,
	mtEntry classloader.MTentry,
	className, methodName, methodType string) int {

	// Parameter array for G function.
	var params []interface{}

	// Append the parameters/args off the stack to params.
	gmethData := mtEntry.Meth.(ghelpers.GMeth)
	paramCount := gmethData.ParamSlots
	for i := 0; i < paramCount; i++ {
		params = append(params, pop(fr))
	}

	// now get the objectRef (the object whose method we're invoking) or a *os.File (stream I/O)
	popped := pop(fr)
	params = append(params, popped)

	if globals.TraceInst {
		infoMsg := fmt.Sprintf("G-function: class=%s, meth=%s%s", className, methodName, methodType)
		trace.Trace(infoMsg)
	}

	// Execute the G function.
	ret := gfunction.RunGfunction(
		mtEntry, fr.FrameStack, className, methodName, methodType, &params, true, globals.TraceInst)
	if ret != nil {
		switch ret.(type) {
		case error: // only occurs in testing
			if globals.GetGlobalRef().JacobinName == "test" {
				return ERROR_OCCURRED
			}
			if errors.Is(ret.(error), gfunction.CaughtGfunctionException) {
				// return 3 // 2 for CP slot + 1 for next bytecode
				// per JACOBIN-59x, we return RESUME_HERE telling
				// the interpreter that the fr.PC has been set to a new position
				// from which processing should continue. This is used primarily
				// when a frame has caught an exception and we're pointing the
				// interpreter to the first bytecode in the exception handler.
				return RESUME_HERE // caught
			}
		default: // if it's not an error, then it's a legitimate return value, which we simply push
			push(fr, ret)
		}
		// any exception will already have been handled.
	}

	return 3 // 2 for CP slot + 1 for next bytecode

}

// OxB7 INVOKESPECIAL
func doInvokespecial(fr *frames.Frame, _ int64) int {
	var className, methodName, methodType, fqn string

	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2]) // next 2 bytes point to CP entry
	CP := fr.CP.(*classloader.CPool)

	entry := CP.CpIndex[CPslot]
	if entry.Type == classloader.Interface {
		className, methodName, methodType =
			classloader.GetMethInfoFromCPinterfaceRef(CP, CPslot)
	} else {
		className, methodName, methodType, fqn = // fqn is the fully qualified name of the method
			classloader.GetMethInfoFromCPmethref(CP, CPslot)
	}

	// if it's a call to java/lang/Object."<init>"()V, which happens frequently,
	// that function simply returns. So test for it here and if it is, skip the rest
	// fullConstructorName := className + "." + methodName + methodType
	// if fqn == "java/lang/Object.<init>()V" { // the java/lang/Object plain constructor just returns
	//	return 3 // 2 for the CPslot + 1 for next bytecode
	// }

	mtEntry, err := classloader.FetchMethodAndCP(className, methodName, methodType)
	if err != nil || mtEntry.Meth == nil {
		// TODO: search the classpath and retry
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := "INVOKESPECIAL: Class method not found: " + fqn
		status := exceptions.ThrowEx(excNames.NoSuchMethodException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	if mtEntry.MType == 'G' { // it's a golang method
		// get the parameters/args, if any, off the stack
		gmethData := mtEntry.Meth.(ghelpers.GMeth)
		paramCount := gmethData.ParamSlots
		var params []interface{}
		for i := 0; i < paramCount; i++ {
			// This is not problematic because the params count in the gfunction definition
			// counts slots, rather than items, so doubles and longs are listed as two slots.
			params = append(params, pop(fr))
		}

		// now get the objectRef (the object whose method we're invoking)
		objRef := pop(fr).(*object.Object)
		params = append(params, objRef)

		if globals.TraceInst {
			infoMsg := fmt.Sprintf("G-function: class=%s, meth=%s%s", className, methodName, methodType)
			trace.Trace(infoMsg)
		}

		ret := gfunction.RunGfunction(
			mtEntry, fr.FrameStack, className, methodName, methodType, &params, true, globals.TraceInst)

		if ret != nil {
			switch ret.(type) {
			case error:
				if globals.GetGlobalRef().JacobinName == "test" {
					return ERROR_OCCURRED
				}
				if errors.Is(ret.(error), gfunction.CaughtGfunctionException) {
					return RESUME_HERE // resume at the present PC, which points to the exception code
				}
			default: // if it's not an error, then it's a legitimate return value, which we simply push
				push(fr, ret)
			}
			// any exception will already have been handled.
		}
		return 3 // 2 for CP slot + 1 for next bytecode
	}

	if mtEntry.MType == 'J' {
		// The arguments are correctly handled in createAndInitNewFrame()
		m := mtEntry.Meth.(classloader.JmEntry)
		if m.AccessFlags&classloader.ACC_NATIVE > 0 {
			// Native code
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKESPECIAL: Native method requested: " + fqn
			status := exceptions.ThrowEx(excNames.UnsupportedOperationException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}
		fram, err := createAndInitNewFrame(className, methodName, methodType, &m, true, fr)
		if err != nil {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKESPECIAL: Error creating frame in: " + fqn
			status := exceptions.ThrowEx(excNames.InvalidStackFrameException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}

		fr.PC += 3                    // point to the next bytecode for when we return from the invoked method.
		fr.FrameStack.PushFront(fram) // push the new frame
		return 0
	}
	return ERROR_OCCURRED // in theory, unreachable
}

// 0xB8 INVOKESTATIC
func doInvokestatic(fr *frames.Frame, _ int64) int {
	var className, methodName, methodType, fqn string

	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2]) // next 2 bytes point to CP entry
	CP := fr.CP.(*classloader.CPool)

	entry := CP.CpIndex[CPslot]
	if entry.Type == classloader.Interface {
		className, methodName, methodType =
			classloader.GetMethInfoFromCPinterfaceRef(CP, CPslot)
	} else {
		className, methodName, methodType, fqn = // fqn is the fully qualified name of the method
			classloader.GetMethInfoFromCPmethref(CP, CPslot)
	}
	mtEntry, err := classloader.FetchMethodAndCP(className, methodName, methodType)
	if err != nil || mtEntry.Meth == nil {
		// TODO: search the classpath and retry
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := "INVOKESTATIC: Class method not found: " + fqn
		status := exceptions.ThrowEx(excNames.NoSuchMethodException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	// before we can run the method, we need to either instantiate the class and/or
	// make sure that its static intializer block (if any) has been run. At this point,
	// all we know is that the class exists and has been loaded.
	k := classloader.MethAreaFetch(className)
	if k.Data.ClInit == types.ClInitNotRun {
		err = runInitializationBlock(k, nil, fr.FrameStack)
		if err != nil {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := fmt.Sprintf("INVOKESTATIC: error running initializer block in %s", fqn)
			status := exceptions.ThrowEx(excNames.ClassNotLoadedException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}
	}

	if mtEntry.MType == 'G' {
		gmethData := mtEntry.Meth.(ghelpers.GMeth)
		paramCount := gmethData.ParamSlots
		var params []interface{}
		for i := 0; i < paramCount; i++ {
			params = append(params, pop(fr))
		}

		if globals.TraceInst {
			infoMsg := fmt.Sprintf("G-function: class=%s, meth=%s%s", className, methodName, methodType)
			trace.Trace(infoMsg)
		}

		ret := gfunction.RunGfunction(mtEntry, fr.FrameStack, className, methodName, methodType, &params, false, globals.TraceInst)
		if ret != nil {
			switch ret.(type) {
			case error:
				if globals.GetGlobalRef().JacobinName == "test" {
					return ERROR_OCCURRED
				} else if errors.Is(ret.(error), gfunction.CaughtGfunctionException) {
					return RESUME_HERE // resume at the present PC, which points to the exception code
				}
			default: // if it's not an error, then it's a legitimate return value, which we simply push
				push(fr, ret)
			}
		}
		return 3
		// any exception will already have been handled.
	} else if mtEntry.MType == 'J' {
		m := mtEntry.Meth.(classloader.JmEntry)
		if m.AccessFlags&classloader.ACC_NATIVE > 0 {
			// Native code
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKESTATIC: Native method requested: " + className + "." + methodName + methodType
			status := exceptions.ThrowEx(excNames.UnsupportedOperationException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}
		fram, err := createAndInitNewFrame(
			className, methodName, methodType, &m, false, fr)
		if err != nil {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKESTATIC: Error creating frame in: " + className + "." + methodName + methodType
			status := exceptions.ThrowEx(excNames.InvalidStackFrameException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}

		fr.PC += 3                    // 2 == initial PC advance in this bytecode + 1 for next bytecode
		fr.FrameStack.PushFront(fram) // push the new frame
		return 0
	}
	return ERROR_OCCURRED // in theory, unreachable code
}

// 0xB9 INVOKEINTERFACE
func doInvokeinterface(fr *frames.Frame, _ int64) int {
	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2]) // next 2 bytes point to CP entry
	count := fr.Meth[fr.PC+3]
	zeroByte := fr.Meth[fr.PC+4]

	CP := fr.CP.(*classloader.CPool)
	CPentry := CP.CpIndex[CPslot]
	if CPentry.Type != classloader.Interface || zeroByte != 0 { // remove the zeroByte test later
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("INVOKEINTERFACE: CP entry type (%d) did not point to an interface method type (%d)",
			CPentry.Type, classloader.Interface)
		status := exceptions.ThrowEx(excNames.IncompatibleClassChangeError, errMsg, fr) // this is the error thrown by JDK
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	method := CP.InterfaceRefs[CPentry.Slot]

	// get the class entry from this method
	interfaceRef := method.ClassIndex
	interfaceNameIndex := CP.ClassRefs[CP.CpIndex[interfaceRef].Slot]
	interfaceNamePtr := stringPool.GetStringPointer(interfaceNameIndex)
	interfaceName := *interfaceNamePtr

	// get the method name for this method
	nAndTindex := method.NameAndType
	nAndTentry := CP.CpIndex[nAndTindex]
	nAndTslot := nAndTentry.Slot
	nAndT := CP.NameAndTypes[nAndTslot]
	interfaceMethodNameIndex := nAndT.NameIndex
	interfaceMethodName := classloader.FetchUTF8stringFromCPEntryNumber(CP, interfaceMethodNameIndex)

	// get the signature for this method
	interfaceMethodSigIndex := nAndT.DescIndex
	interfaceMethodType := classloader.FetchUTF8stringFromCPEntryNumber(
		CP, interfaceMethodSigIndex)

	// now get the objRef pointing to the class containing the call to the method
	// described just previously. It is located on the f.OpStack below the args to
	// be passed to the method.
	// The objRef object has previously been instantiated and its constructor called.
	objRef := fr.OpStack[fr.TOS-int(count)+1]
	if objRef == nil {
		errMsg := fmt.Sprintf("INVOKEINTERFACE: object whose method, %s, is invoked is null",
			interfaceName+interfaceMethodName+interfaceMethodType)
		status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	// get the name of the objectRef's class, and make sure it's loaded
	objRefClassName := *(stringPool.GetStringPointer(objRef.(*object.Object).KlassName))
	if err := classloader.LoadClassFromNameOnly(objRefClassName); err != nil {
		// in this case, LoadClassFromNameOnly() will have already thrown the exception
		if globals.JacobinHome() == "test" {
			return ERROR_OCCURRED // applies only if in test
		}
	}

	class := classloader.MethAreaFetch(objRefClassName)
	if class == nil {
		// in theory, this can't happen due to immediately previous loading, but making sure
		errMsg := fmt.Sprintf("INVOKEINTERFACE: class %s not found", objRefClassName)
		status := exceptions.ThrowEx(excNames.ClassNotLoadedException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	var mtEntry classloader.MTentry
	var err error
	mtEntry, err = locateInterfaceMeth(class, fr, objRefClassName, interfaceName,
		interfaceMethodName, interfaceMethodType)
	if err != nil { // any error will already have been handled
		return ERROR_OCCURRED
	}

	clData := *class.Data
	if mtEntry.MType == 'J' {
		var className string
		entry := mtEntry.Meth.(classloader.JmEntry)
		if entry.AccessFlags&classloader.ACC_PRIVATE > 0 {
			className = interfaceName
		} else {
			className = clData.Name
		}
		// the args and the objRef are popped off the stack by following call
		fram, err := createAndInitNewFrame(
			className, interfaceMethodName, interfaceMethodType, &entry, true, fr)
		if err != nil {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "INVOKEINTERFACE: Error creating frame in: " + className + "." +
				interfaceMethodName + interfaceMethodType
			status := exceptions.ThrowEx(excNames.InvalidStackFrameException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
		}
		fr.PC += 5                    // 2 for CP slot, 1 for count, 1 for zero byte, 1 for next bytecode
		fr.FrameStack.PushFront(fram) // push the new frame
		return 0                      // forcing execution of the new frame
	} else if mtEntry.MType == 'G' { // it's a gfunction (i.e., a native function implemented in golang)
		gmethData := mtEntry.Meth.(ghelpers.GMeth)
		paramCount := gmethData.ParamSlots
		var params []interface{}
		for i := 0; i < paramCount; i++ {
			params = append(params, pop(fr))
		}

		// now push the objectRef (the object whose method we're invoking)
		params = append(params, objRef)
		_ = pop(fr) // the objRef is still on the stack, so pop it off now

		if globals.TraceInst {
			infoMsg := fmt.Sprintf("G-function: interface=%s, meth=%s%s", interfaceName, interfaceName, interfaceMethodType)
			trace.Trace(infoMsg)
		}
		ret := gfunction.RunGfunction(
			mtEntry, fr.FrameStack, interfaceName, interfaceMethodName, interfaceMethodType, &params, true,
			globals.TraceVerbose)
		if ret != nil {
			switch ret.(type) {
			case error:
				if globals.GetGlobalRef().JacobinName == "test" {
					return ERROR_OCCURRED
				} else if errors.Is(ret.(error), gfunction.CaughtGfunctionException) {
					return ERROR_OCCURRED
				}
			default: // if it's not an error, then it's a legitimate return value, which we simply push
				push(fr, ret)
				return 5 // 2 for CP slot + 1 for count, 1 for zero byte, and 1 for next bytecode
			}
		}
		// Normal nil return.
		return 5 // 2 for CP slot + 1 for count, 1 for zero byte, and 1 for next bytecode
	}

	// Impossible return point unless there is a preceding software error.
	globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
	errMsg := "INVOKEINTERFACE: Logic error creating frame in: " + clData.Name + "." +
		interfaceMethodName + interfaceMethodType
	status := exceptions.ThrowEx(excNames.VirtualMachineError, errMsg, fr)
	if status != exceptions.Caught {
		return ERROR_OCCURRED // applies only if in test
	}
	return notImplemented(fr, 0) // in theory, unreachable code
}

// 0xBA INVOKEDYNAMIC
func doInvokedynamic(fr *frames.Frame, _ int64) int {
	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2]) // next 2 bytes point to CP entry
	CP := fr.CP.(*classloader.CPool)
	CPentry := CP.CpIndex[CPslot]

	if CPentry.Type != classloader.InvokeDynamic {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg :=
			fmt.Sprintf("INVOKEDYNAMIC: constant pool entry is (%d) rather than Constant_InvokeDynamic_info", CPentry.Type)
		status := exceptions.ThrowEx(excNames.ClassFormatError, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		} else {
			return RESUME_HERE // caught
		}
	}

	id := CP.InvokeDynamics[CPentry.Slot]
	bootstrapIndex := id.BootstrapIndex
	bootstrap := CP.Bootstraps[bootstrapIndex]
	if bootstrap.MethodRef < 0 { // dummy test to get a compilation
		return ERROR_OCCURRED
	}

	// bootstrapNAT := id.NameAndType
	
	return 5 // the two bytes for the CP slot + 2 bytes with value 0x00 + 1 for next bytecode

}

// 0xBB NEW create a new object
func doNew(fr *frames.Frame, _ int64) int {
	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2]) // next 2 bytes point to CP entry
	CP := fr.CP.(*classloader.CPool)
	CPentry := CP.CpIndex[CPslot]
	if CPentry.Type != classloader.ClassRef && CPentry.Type != classloader.Interface {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("NEW: Invalid type for new object")
		status := exceptions.ThrowEx(excNames.ClassFormatError, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	// the classref points to a UTF8 record with the name of the class to instantiate
	var className string
	if CPentry.Type == classloader.ClassRef {
		nameStringPoolIndex := CP.ClassRefs[CPentry.Slot]
		className = *stringPool.GetStringPointer(nameStringPoolIndex)
	}

	ref, err := InstantiateClass(className, fr.FrameStack)
	if err != nil {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("NEW: could not load class %s", className)
		status := exceptions.ThrowEx(excNames.ClassNotLoadedException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}
	push(fr, ref.(*object.Object))
	return 3 // 2 for CPslot + 1 for next bytecode
}

// 0xBC NEWARRAY create a new array of primitives
func doNewarray(fr *frames.Frame, _ int64) int {
	size := pop(fr).(int64)
	if size < 0 {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := "NEWARRAY: Invalid size for array"
		status := exceptions.ThrowEx(excNames.NegativeArraySizeException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	arrayType := int(fr.Meth[fr.PC+1])

	arrayPtr := object.Make1DimArray(uint8(arrayType), size)
	g := globals.GetGlobalRef()
	g.ArrayAddressList.PushFront(arrayPtr)
	push(fr, arrayPtr)
	return 2 // 1 for the array type + 1 for next byte
}

// 0xBD ANEWARRAY create an array of pointers
func doAnewarray(fr *frames.Frame, _ int64) int {
	size := pop(fr).(int64)
	if size < 0 {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := "ANEWARRAY: Invalid size for array"
		status := exceptions.ThrowEx(excNames.NegativeArraySizeException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	refTypeSlot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2]) // 2 bytes point to CP entry
	CP := fr.CP.(*classloader.CPool)
	refType := CP.CpIndex[refTypeSlot]
	if refType.Type != classloader.ClassRef && refType.Type != classloader.Interface {
		// TODO: it could also point to an array, per the JVM spec
		errMsg := fmt.Sprintf("ANEWARRAY: Presently works only with classes and interfaces")
		status := exceptions.ThrowEx(excNames.UnsupportedOperationException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	var refTypeName = ""
	if refType.Type == classloader.ClassRef {
		refNameStringPoolIndex := CP.ClassRefs[refType.Slot]
		refTypeName = *stringPool.GetStringPointer(refNameStringPoolIndex)
	}

	arrayPtr := object.Make1DimRefArray(refTypeName, size)
	g := globals.GetGlobalRef()
	g.ArrayAddressList.PushFront(arrayPtr)
	push(fr, arrayPtr)
	return 3 // 2 for RefTypeSlot + 1 for next bytecode
}

// 0xBE ARRAYLENGTH get size of an array
func doArraylength(fr *frames.Frame, _ int64) int {
	ref := pop(fr) // pointer to the array
	if ref == nil {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := "ARRAYLENGTH: Invalid (null) reference to an array"
		status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
	}

	var size int64
	switch ref.(type) {
	// the type of array reference can vary. For many instances,
	// it will be a pointer to an array object. In other cases,
	// such as inside Java String class, the actual primitive
	// array of bytes will be extracted as a field and passed
	// to this function, so we need to accommodate all types--
	// hence, the switch on type.
	case []byte:
		array := ref.([]uint8)
		size = int64(len(array))
	case []types.JavaByte:
		array := ref.([]types.JavaByte)
		size = int64(len(array))
	case []float64:
		array := ref.([]float64)
		size = int64(len(array))
	case []int64:
		array := ref.([]int64)
		size = int64(len(array))
	case *[]byte:
		array := *ref.(*[]uint8)
		size = int64(len(array))
	case *[]types.JavaByte:
		array := *ref.(*[]types.JavaByte)
		size = int64(len(array))
	case []*object.Object:
		array := ref.([]*object.Object)
		size = int64(len(array))
	case *object.Object:
		r := ref.(*object.Object)
		if object.IsNull(r) {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "ARRAYLENGTH: Invalid (null) value for *object.Object"
			status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}
		r.ThMutex.RLock()
		field, ok := r.FieldTable["value"]
		r.ThMutex.RUnlock()
		if !ok {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := "ARRAYLENGTH: Value field missing for *object.Object"
			status := exceptions.ThrowEx(excNames.VirtualMachineError, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}
		if !strings.HasPrefix(field.Ftype, types.Array) {
			globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
			errMsg := fmt.Sprintf("ARRAYLENGTH: Value field for *object.Object is not an array, got: %s", field.Ftype)
			status := exceptions.ThrowEx(excNames.VirtualMachineError, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED // applies only if in test
			}
			return RESUME_HERE // caught
		}
		size = object.ArrayLength(r)
	default:
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("ARRAYLENGTH: Invalid ref.(type): %T", ref)
		status := exceptions.ThrowEx(excNames.IllegalArgumentException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}
	push(fr, size)
	return 1
}

// 0xBF ATHROW throw an exception
func doAthrow(fr *frames.Frame, _ int64) int {
	// objRef points to an instance of the error/exception class that's being thrown
	objectRef := pop(fr).(*object.Object)
	if object.IsNull(objectRef) {
		errMsg := "ATHROW: Invalid (null) reference to an exception/error class to throw"
		status := exceptions.ThrowEx(excNames.NullPointerException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
	}

	// capture the golang stack
	stack := string(debug.Stack())
	globals.GetGlobalRef().ErrorGoStack = stack

	// capture the JVM frame stack
	globals.GetGlobalRef().JVMframeStack = exceptions.GrabFrameStack(fr.FrameStack)

	// get the name of the exception in the format used by HotSpot
	exceptionClass := *(stringPool.GetStringPointer(objectRef.KlassName))
	exceptionName := strings.Replace(exceptionClass, "/", ".", -1)

	// get the PC of the exception and check for any catch blocks
	// if f.ExceptionPC == -1 {
	//	f.ExceptionPC = f.PC
	// }

	// find the frame with a valid catch block for this exception, if any
	catchFrame, handlerBytecode := exceptions.FindCatchFrame(fr.FrameStack, exceptionName, fr.ExceptionPC)
	// if there is no catch block, then print out the data we have (conforming
	// with whether we want the standard JDK info as elected with the -strictJDK
	// command-line option)
	if catchFrame == nil {
		// if the exception is not caught, then print the data from the stackTraceElements (STEs)
		// in the Throwable object or subclass (which is generally the specific exception class).

		// start by printing out the name of the exception/error and the thread it occurred on
		errMsg := ""
		if fr.Thread == 1 { // if it's thread #1, use its name, "main"
			errMsg = fmt.Sprintf("Exception in thread \"main\" %s", exceptionName)
		} else {
			errMsg = fmt.Sprintf("Exception in thread %d %s", fr.Thread, exceptionName)
		}

		objectRef.ThMutex.RLock()
		appMsgFld, ok := objectRef.FieldTable["detailMessage"]
		objectRef.ThMutex.RUnlock()
		if ok {
			appMsg := appMsgFld.Fvalue
			if appMsg != object.Null && appMsg != nil {
				switch appMsg.(type) {
				case []types.JavaByte:
					jbarray := appMsg.([]types.JavaByte)
					errMsg += fmt.Sprintf(": %s", object.GoStringFromJavaByteArray(jbarray))
				case *object.Object:
					var value any
					obj := appMsg.(*object.Object)
					obj.ThMutex.RLock()
					fld, ok := obj.FieldTable["value"]
					obj.ThMutex.RUnlock()
					if !ok {
						value = "<missing>"
					} else {
						value = fld.Fvalue
					}
					switch value.(type) {
					case []byte:
						obj.ThMutex.RLock()
						errMsg += fmt.Sprintf(": %s", string(obj.FieldTable["value"].Fvalue.([]byte)))
						obj.ThMutex.RUnlock()
					case []types.JavaByte:
						obj.ThMutex.RLock()
						errMsg += fmt.Sprintf(": %s", object.GoStringFromJavaByteArray(obj.FieldTable["value"].Fvalue.([]types.JavaByte)))
						obj.ThMutex.RUnlock()
					case uint32:
						str := stringPool.GetStringPointer(value.(uint32))
						errMsg += fmt.Sprintf(": %s", *str)
					default:
						errMsg += fmt.Sprintf(": <default value> %v", value)
					}
				default:
					errMsg += ": objectRef.FieldTable[\"detailMessage\"] is object.Null"
				}
			}
		}
		trace.AsIs(errMsg)

		objectRef.ThMutex.RLock()
		steArrayPtrObj := objectRef.FieldTable["stackTrace"].Fvalue.(*object.Object)
		objectRef.ThMutex.RUnlock()

		steArrayPtrObj.ThMutex.RLock()
		rawSteArray := steArrayPtrObj.FieldTable["value"].Fvalue.([]*object.Object) // []*object.Object (each of which is an STE)
		steArrayPtrObj.ThMutex.RUnlock()

		for i := 0; i < len(rawSteArray); i++ {
			ste := rawSteArray[i]
			ste.ThMutex.RLock()
			rawClassName := ste.FieldTable["declaringClass"].Fvalue.(string)
			ste.ThMutex.RUnlock()
			if rawClassName == "java/lang/Throwable" || // don't show Throwable methods
				rawClassName == "java/lang/Error" || // don't show Error methods
				rawClassName == "java/lang/AssertionError" { // don't show AssertionError methods
				continue
			}
			ste.ThMutex.RLock()
			methodName := ste.FieldTable["methodName"].Fvalue.(string)
			ste.ThMutex.RUnlock()
			if methodName == "<init>" { // don't show constructors
				continue
			}
			className := strings.Replace(rawClassName, "/", ".", -1)

			ste.ThMutex.RLock()
			sourceLine := ste.FieldTable["sourceLine"].Fvalue.(string)
			fileName := ste.FieldTable["fileName"].Fvalue
			ste.ThMutex.RUnlock()

			var errMsg string
			if sourceLine != "" {
				errMsg = fmt.Sprintf("\tat %s.%s(%s:%s)", className,
					methodName, fileName, sourceLine)
			} else {
				errMsg = fmt.Sprintf("\tat %s.%s(%s)", className,
					methodName, fileName)
			}
			trace.AsIs(errMsg)
		}

		// show Jacobin's JVM stack info if -strictJDK is not set
		if globals.GetGlobalRef().StrictJDK == false {
			trace.Trace(" ")
			for _, frameData := range *globals.GetGlobalRef().JVMframeStack {
				colon := strings.Index(frameData, ":")
				shortenedFrameData := frameData[colon+1:]
				trace.Trace("\tat" + shortenedFrameData)
			}
		}

		// all exceptions that got this far are untrapped, so shutdown with an error code
		shutdown.Exit(shutdown.APP_EXCEPTION)

	} else { // perform the catch operation. We know the frame and the starting bytecode for the handler
		for listElem := fr.FrameStack.Front(); listElem != nil; {
			var wframe = listElem.Value.(*frames.Frame)
			if wframe == catchFrame {
				wframe.TOS = -1
				push(wframe, objectRef)
				wframe.PC = handlerBytecode
				return 0
			}
			// if the frame we're about to pop is synchronized, unlock it
			if wframe.ObjSync != nil {
				_ = wframe.ObjSync.ObjUnlock(int32(wframe.Thread))
				if globals.TraceInst {
					traceInfo := fmt.Sprintf("\tdoAthrow: Unlocked object %s",
						object.GoStringFromStringPoolIndex(wframe.ObjSync.KlassName))
					trace.Trace(traceInfo)
				}
			}
			next := listElem.Next()
			fr.FrameStack.Remove(listElem)
			listElem = next
		}
	}
	return 1 // should not be reached, in theory
}

// 0xC0 CHECKCAST
func doCheckcast(fr *frames.Frame, _ int64) int {
	//
	// CHECKCAST and INSTANCEOF are nearly the same. Differences:
	//
	// * CHECKCAST: Peeks at the objectRef on the stack. It returns immediately if objectRef is null (not an error).
	// * INSTANCEOF: Pops the objectRef and pushes a result of types.JavaBoolTrue (if an instance of) or types.JavaBoolFalse (if not).
	// * Both return 3 (2 for CPslot + 1 for next byte) or some sort of error.
	//
	// Because doCheckcast uses the same middle logic as doInstanceof,
	// any change here should probably be made to doInstanceof
	// and vice versa!

	ref := peek(fr) // peek b/c the objectRef is *not* removed from the op stack
	if ref == nil { // if ref is nil, just carry on
		return 3 // move past two bytes pointing to comp object + 1 for next bytecode
	}

	var obj *object.Object
	var objName string
	switch ref.(type) {
	case *object.Object:
		if object.IsNull(ref) { // if ref is null, just carry on
			return 3
		} else {
			obj = (ref).(*object.Object)
			objName = *(stringPool.GetStringPointer(obj.KlassName))
		}
	default: // objectRef must be a reference to an object
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("CHECKCAST: Invalid class reference, type=%T", ref)
		status := exceptions.ThrowEx(excNames.ClassCastException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	// at this point, we know we have a non-nil/non-null pointer to an object;
	// now, get the class we're casting the object to.
	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2])
	CP := fr.CP.(*classloader.CPool)
	// CPentry := CP.CpIndex[CPslot]
	classNamePtr := classloader.FetchCPentry(CP, CPslot)
	targetClassName := *(classNamePtr.StringVal)

	if globals.TraceVerbose {
		trace.Trace(fmt.Sprintf("CHECKCAST: object=%s, target=%s",
			objName, targetClassName))
	}

	// Derive objClassType, the reference object type.
	var objClassType = types.Error
	if strings.HasPrefix(objName, types.Array) {
		// The reference object  is an array.
		objClassType = types.Array
	} else {
		// The reference object  is not an array.
		objData := classloader.MethAreaFetch(objName)
		if objData == nil || objData.Data == nil {
			_ = classloader.LoadClassFromNameOnly(objName)
			objData = classloader.MethAreaFetch(objName)
		}
		// Interface or a non-array object?
		if objData.Data.Access.ClassIsInterface {
			objClassType = types.Interface
		} else {
			objClassType = types.NonArrayObject
		}
	}

	// Depending on objClassType, execute one of the cast checks.
	var checkcastStatus bool
	switch objClassType {
	case types.NonArrayObject:
		checkcastStatus = checkcastNonArrayObject(obj, targetClassName)
	case types.Array:
		checkcastStatus = checkcastArray(obj, targetClassName)
	case types.Interface:
		checkcastStatus = checkcastInterface(obj, targetClassName)
		if checkcastStatus {
			stObj := pop(fr).(*object.Object)
			stObj.KlassName = object.StringPoolIndexFromGoString(targetClassName)
			push(fr, stObj)
		}
	default:
		errMsg := fmt.Sprintf("CHECKCAST: expected to verify class or interface, but got: %s", objClassType)
		status := exceptions.ThrowEx(excNames.InvalidTypeException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	// If false (failure) return, throw an exception.
	if checkcastStatus == false {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("CHECKCAST: %s is not castable with respect to %s",
			*(stringPool.GetStringPointer(obj.KlassName)), targetClassName)
		status := exceptions.ThrowEx(excNames.ClassCastException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	return 3 // 2 for CPslot + 1 for next byte
}

// 0xC1 INSTANCEOF validate the type of object (if not nil or null)
func doInstanceof(fr *frames.Frame, _ int64) int {
	//
	// See detail design documentation in doCheckcast.
	//
	// Because doInstanceof uses the same middle logic as doCheckcast,
	// any change here should probably be made to doCheckcast
	// and vice versa!
	//
	ref := pop(fr)
	if ref == nil || object.IsNull(ref) {
		if globals.TraceVerbose {
			trace.Trace("INSTANCEOF: null reference -> false")
		}
		push(fr, types.JavaBoolFalse)
		return 3 // 2 bytes for CP index + 1 for next bytecode
	}

	var obj *object.Object
	var objName string
	switch ref.(type) {
	case *object.Object:
		if object.IsNull(ref) { // if ref is null, just carry on
			return 3
		} else {
			obj = (ref).(*object.Object)
			objName = *(stringPool.GetStringPointer(obj.KlassName))
		}
	default: // objectRef must be a reference to an object
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("CHECKCAST: Invalid class reference, type=%T", ref)
		status := exceptions.ThrowEx(excNames.ClassCastException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	// at this point, we know we have a non-nil/non-null pointer to an object;
	// now, get the class we're casting the object to.
	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2])
	CP := fr.CP.(*classloader.CPool)
	// CPentry := CP.CpIndex[CPslot]
	classNamePtr := classloader.FetchCPentry(CP, CPslot)
	targetClassName := *(classNamePtr.StringVal)

	if globals.TraceVerbose {
		trace.Trace(fmt.Sprintf("INSTANCEOF: object=%s, target=%s",
			objName, targetClassName))
	}

	var objClassType = types.Error
	if strings.HasPrefix(objName, "[") {
		objClassType = types.Array
	} else {
		objData := classloader.MethAreaFetch(objName)
		if objData == nil || objData.Data == nil {
			_ = classloader.LoadClassFromNameOnly(objName)
			objData = classloader.MethAreaFetch(objName)
		}
		if objData.Data.Access.ClassIsInterface {
			objClassType = types.Interface
		} else {
			objClassType = types.NonArrayObject
		}
	}

	// Perform assignability check using same helpers as CHECKCAST
	var isInstance bool
	switch objClassType {
	case types.NonArrayObject:
		isInstance = checkcastNonArrayObject(obj, targetClassName)
	case types.Array:
		isInstance = checkcastArray(obj, targetClassName)
	case types.Interface:
		isInstance = checkcastInterface(obj, targetClassName)
	default:
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := "INSTANCEOF: expected to verify class or interface, but got none"
		status := exceptions.ThrowEx(excNames.InvalidTypeException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED
		}
		return RESUME_HERE
	}

	// Push result onto operand stack
	if isInstance {
		if globals.TraceVerbose {
			trace.Trace("INSTANCEOF: result = true")
		}
		push(fr, types.JavaBoolTrue)
	} else {
		if globals.TraceVerbose {
			trace.Trace("INSTANCEOF: result = false")
		}
		push(fr, types.JavaBoolFalse)
	}

	return 3 // 2 bytes for CP slot + 1 for next byte
}

// 0xC2 MONITORENTER: Lock an object.
// note: set mtrdebug to true to enable tracing/debugging

func doMonitorenter(fr *frames.Frame, _ int64) int {
	// Check stack size.
	if fr.TOS < 0 {
		errMsg := fmt.Sprintf("MONITORENTRY: stack underflow in %s.%s",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName)
		status := exceptions.ThrowEx(excNames.InternalException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	// Get object reference from stack.
	obj, ok := pop(fr).(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("MONITORENTRY: expected an object but encountered: %v", obj)
		status := exceptions.ThrowEx(excNames.InvalidTypeException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	// Yield.
	if mtrdebug {
		var msg string
		msg = fmt.Sprintf("DEBUG MONITOR-ENTRY YIELD threadID=%d", fr.Thread)
		trace.Trace(msg)
	}
	runtime.Gosched()

	// DEBUG if configured to do so.
	if mtrdebug {
		var msg string
		if obj.Monitor != nil {
			owner := obj.GetMonitorOwner()
			recursion := obj.GetMonitorRecursion()
			msg = fmt.Sprintf("DEBUG MONITOR-ENTRY try for ObjLock-fat threadID=%d, owner=%d count=%d", fr.Thread, owner, recursion)
		} else {
			msg = fmt.Sprintf("DEBUG MONITOR-ENTRY try for ObjLock-thin threadID=%d", fr.Thread)
		}
		trace.Trace(msg)
	}

	// Lock the object to this thread.
	err := obj.ObjLock(int32(fr.Thread))
	if err != nil {
		errMsg := fmt.Sprintf("MONITORENTRY: %s in %s.%s%s", err.Error(),
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, fr.MethType)
		status := exceptions.ThrowEx(excNames.InternalException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	// DEBUG if configured to do so.
	if mtrdebug {
		var msg string
		if obj.Monitor != nil {
			owner := obj.GetMonitorOwner()
			recursion := obj.GetMonitorRecursion()
			msg = fmt.Sprintf("DEBUG MONITOR-ENTRY success ObjLock-fat threadID=%d, owner=%d count=%d", fr.Thread, owner, recursion)
		} else {
			msg = fmt.Sprintf("DEBUG MONITOR-ENTRY success ObjLock-thin threadID=%d", fr.Thread)
		}
		trace.Trace(msg)
	}

	return 1
}

// 0xC3 MONITOREXIT: Release an object.
// note: set mtrdebug to true to enable tracing/debugging
func doMonitorexit(fr *frames.Frame, _ int64) int {
	// Check stack size.
	if fr.TOS < 0 {
		errMsg := fmt.Sprintf("MONITOREXIT: stack underflow in %s.%s",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName)
		status := exceptions.ThrowEx(excNames.InternalException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	// Get object reference from stack.
	obj, ok := pop(fr).(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("MONITOREXIT: expected an object but encountered: %v", obj)
		status := exceptions.ThrowEx(excNames.InvalidTypeException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	// DEBUG if configured to do so.
	if mtrdebug {
		var msg string
		if obj.Monitor != nil {
			owner := obj.GetMonitorOwner()
			recursion := obj.GetMonitorRecursion()
			msg = fmt.Sprintf("DEBUG MONITOR-EXIT ObjUnlock-fat threadID=%d, owner=%d count=%d", fr.Thread, owner, recursion)
		} else {
			msg = fmt.Sprintf("DEBUG MONITOR-EXIT ObjUnlock-thin threadID=%d", fr.Thread)
		}
		trace.Trace(msg)
	}

	// Release the object from this thread.
	err := obj.ObjUnlock(int32(fr.Thread))
	if err != nil {
		// Specific granular checks for compatibility with tests
		if !object.IsObjectLocked(obj) {
			errMsg := fmt.Sprintf("MONITOREXIT: object is already unlocked in %s.%s%s",
				util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, fr.MethType)
			status := exceptions.ThrowEx(excNames.InternalException, errMsg, fr)
			if status != exceptions.Caught {
				return ERROR_OCCURRED
			}
			return RESUME_HERE
		}
		if (obj.Mark.Misc & 0b11) == 2 {
			monitor := obj.GetMonitor()
			if monitor == nil || monitor.Owner == -1 {
				errMsg := fmt.Sprintf("MONITOREXIT: fat lock exists but monitor is nil in %s.%s%s",
					util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, fr.MethType)
				status := exceptions.ThrowEx(excNames.InternalException, errMsg, fr)
				if status != exceptions.Caught {
					return ERROR_OCCURRED
				}
				return RESUME_HERE
			}
		}

		errMsg := fmt.Sprintf("MONITOREXIT: %s in %s.%s%s", err.Error(),
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, fr.MethType)
		status := exceptions.ThrowEx(excNames.InternalException, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	// Yield.
	if mtrdebug {
		var msg string
		msg = fmt.Sprintf("DEBUG MONITOR-EXIT YIELD threadID=%d", fr.Thread)
		trace.Trace(msg)
	}
	runtime.Gosched()

	return 1
}

// 0xC4 WIDE use wide versions of bytecode arguments
func doWide(fr *frames.Frame, _ int64) int {
	fr.WideInEffect = true
	return 1
}

// 0xC5 MULTIANEWARRAY create a multi-dimensional array
func doMultinewarray(fr *frames.Frame, _ int64) int {
	var arrayDesc string
	var arrayType uint8

	// The first two bytes after the bytecode point to a classref entry in the CP.
	// In turn, it points to a string describing the array of the form [[L or
	// similar, in which one [ is present for every array dimension, followed by a
	// single letter describing the type of primitive in the leaf dimension of the array.
	// The letters are the usual ones used in the JVM for primitives, etc.
	// as in: https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-4.html#jvms-4.3.2-200
	CPslot := (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2]) // point to CP entry
	CP := fr.CP.(*classloader.CPool)
	CPentry := CP.CpIndex[CPslot]
	arrayDescStringPoolIndex := CP.ClassRefs[CPentry.Slot]
	arrayDesc = *stringPool.GetStringPointer(arrayDescStringPoolIndex)

	var rawArrayType uint8
	for i := 0; i < len(arrayDesc); i++ {
		if arrayDesc[i] != '[' {
			rawArrayType = arrayDesc[i]
			break
		}
	}

	switch rawArrayType {
	case 'Z':
		arrayType = object.T_BOOLEAN
	case 'B':
		arrayType = object.T_BYTE
	case 'C':
		arrayType = object.T_CHAR
	case 'D':
		arrayType = object.T_DOUBLE
	case 'F':
		arrayType = object.T_FLOAT
	case 'I':
		arrayType = object.T_INT
	case 'J':
		arrayType = object.T_LONG
	case 'S':
		arrayType = object.T_SHORT
	case 'L':
		arrayType = object.T_REF
	default:
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("MULTIANEWARRAY: Unexpected raw array type: %T", rawArrayType)
		status := exceptions.ThrowEx(excNames.VirtualMachineError, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	// get the number of dimensions, then pop off the operand
	// stack an int for every dimension, giving the size of that
	// dimension and put them into a slice that starts with
	// the highest dimension first. So a two-dimensional array
	// such as x[4][3], would have entries of 4 and 3 respectively
	// in the dimsizes slice.
	dimensionCount := int(fr.Meth[fr.PC+3])

	if dimensionCount > 3 { // TODO: explore arrays of > 5-255 dimensions
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := "MULTIANEWARRAY: Jacobin supports arrays only up to three dimensions"
		status := exceptions.ThrowEx(excNames.VirtualMachineError, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
	}

	dimSizes := make([]int64, dimensionCount)

	// the values on the operand stack give the last dimension
	// first when popped off the stack, so, they're stored here
	// in reverse order, so that dimSizes[0] will hold the first
	// dimenion.
	for i := dimensionCount - 1; i >= 0; i-- {
		dimSizes[i] = pop(fr).(int64)
	}

	// A dimension of zero ends the dimensions, so we check
	// and cut off the dimensions below and includingthe 0-sized
	// one. Because this is almost certainly an error, we also
	// issue a warning.
	for i := range dimSizes {
		if dimSizes[i] == 0 {
			dimSizes = dimSizes[i+1:] // lop off the prev dims
			trace.Error("MULTIANEWARRAY: Multidimensional array with one dimension of size 0 encountered.")
			break
		}
	}

	// Because of the possibility of a zero-sized dimension
	// affecting the valid number of dimensions, dimensionCount
	// can no longer be considered reliable. Use len(dimSizes).
	if len(dimSizes) == 3 {
		multiArr := object.Make1DimArray(object.T_REF, dimSizes[0])
		multiArr.ThMutex.RLock()
		actualArray := multiArr.FieldTable["value"].Fvalue.([]*object.Object)
		multiArr.ThMutex.RUnlock()
		for i := 0; i < len(actualArray); i++ {
			actualArray[i], _ = object.Make2DimArray(dimSizes[1],
				dimSizes[2], arrayType)
		}
		push(fr, multiArr)

	} else if len(dimSizes) == 2 { // 2-dim array is a special, trivial case
		multiArr, _ := object.Make2DimArray(dimSizes[0], dimSizes[1], arrayType)
		push(fr, multiArr)
		// It's possible due to a zero-length dimension, that we
		// need to create a single-dimension array.
	} else if len(dimSizes) == 1 {
		oneDimArr := object.Make1DimArray(arrayType, dimSizes[0])
		push(fr, oneDimArr)

	}
	return 4 // 2 for CPslot + 1 for dimensions + 1 for next bytecode
}

// 0xC6 IFNULL jump if TOS holds a null address
func doIfnull(fr *frames.Frame, _ int64) int {
	// null = nil or object.Null (a pointer to nil)
	value := pop(fr)
	if object.IsNull(value) {
		jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
		return int(jumpTo)
	} else {
		return 3
	}
}

// 0xC7 IFNONNULL jump if TOS does not hold a null address
func doIfnonnull(fr *frames.Frame, _ int64) int {
	value := pop(fr)
	if object.IsNull(value) { // if == null, move along
		return 3
	} else { // it's not nil nor a null pointer--so do the jump
		jumpTo := (int16(fr.Meth[fr.PC+1]) * 256) + int16(fr.Meth[fr.PC+2])
		return int(jumpTo)
	}
}

// 0xC8 GOTO_W jump to a four-byte offset
func doGotow(fr *frames.Frame, _ int64) int {
	jumpTo := types.FourBytesToInt64(
		fr.Meth[fr.PC+1], fr.Meth[fr.PC+2], fr.Meth[fr.PC+3], fr.Meth[fr.PC+4])
	return int(jumpTo)
}

// 0xC9 JSR_W jump to a four-byte offset
func doJsrw(fr *frames.Frame, _ int64) int {
	jumpTo := types.FourBytesToInt64(
		fr.Meth[fr.PC+1], fr.Meth[fr.PC+2], fr.Meth[fr.PC+3], fr.Meth[fr.PC+4])
	push(fr, jumpTo) // JSR and JSR_W both push the jump offset and jump to it
	return int(jumpTo)
}

func notImplemented(fr *frames.Frame, _ int64) int {
	opcode := fr.Meth[fr.PC]
	opcodeName := opcodes.BytecodeNames[opcode]
	errMsg := fmt.Sprintf("bytecode %s not implemented at present", opcodeName)
	_ = exceptions.ThrowEx(excNames.UnsupportedOperationException, errMsg, fr)
	return ERROR_OCCURRED
}

func doWarninvalid(fr *frames.Frame, _ int64) int {
	opcode := fr.Meth[fr.PC]
	opcodeName := opcodes.BytecodeNames[opcode]
	errMsg := fmt.Sprintf("bytecode %s not implemented at present", opcodeName)
	trace.Warning(errMsg)
	return 1
}

// === helper methods--that is, functions called by dispatched methods (in alpha order) ===

func load(fr *frames.Frame, local int64) int {
	push(fr, fr.Locals[local])
	return 1
}

// load a constant
// TODO: eventually need to handle CONSTANT_Dynamic_info (tag 17)
func ldc(fr *frames.Frame, width int) int {
	var idx int
	if width == 1 { // LDC uses a 1-byte index into the CP, LDC_W uses a 2-byte index
		idx = int(fr.Meth[fr.PC+1])
	} else {
		idx = (int(fr.Meth[fr.PC+1]) * 256) + int(fr.Meth[fr.PC+2])
	}

	CPe := classloader.FetchCPentry(fr.CP.(*classloader.CPool), idx)
	if CPe.EntryType == classloader.Dummy || // 0 = error
		// Note: an invalid CP entry causes a java.lang.Verify error and
		//       is caught before execution of the program begins.
		// This bytecode does not load longs or doubles
		CPe.EntryType == classloader.DoubleConst ||
		CPe.EntryType == classloader.LongConst {
		globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("in %s.%s, LDC: Invalid type for bytecode operand: %d",
			util.ConvertInternalClassNameToUserFormat(fr.ClName), fr.MethName, CPe.EntryType)
		status := exceptions.ThrowEx(excNames.ClassFormatError, errMsg, fr)
		if status != exceptions.Caught {
			return ERROR_OCCURRED // applies only if in test
		}
		return RESUME_HERE // caught
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
	case classloader.IS_CLASS_REF: // push a class object in support of static synchronized methods
		push(fr, fr.ObjSync)
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

func store(fr *frames.Frame, local int64) int {
	fr.Locals[local] = pop(fr)
	return 1
}

func storeInt(fr *frames.Frame, local int64) int {
	// because we could be storing a byte, boolean, short, etc.
	// we must convert the interface to an int64.
	fr.Locals[local] = convertInterfaceToInt64(pop(fr))
	return 1
}
