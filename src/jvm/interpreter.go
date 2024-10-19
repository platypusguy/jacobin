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
 */

package jvm

import (
	"jacobin/frames"
	"jacobin/object"
)

// set up a DispatchTable with 203 slots that correspond to the bytecodes
// each slot being a pointer to a function that accepts a pointer to the
// current frame and an int parameter. It returns an int that indicates
// how much to increase that frame's PC (program counter) by.
type BytecodeFunc func(*frames.Frame, int64) int

var DispatchTable = [203]BytecodeFunc{
	doNop,          // NOP         0x00
	doAconstNull,   // ACONST_NULL 0x01
	doIconstM1,     // ICONST_M1   0x02
	doIconst0,      // ICONST_0    0x03
	doIconst1,      // ICONST_1    0x04
	doIconst2,      // ICONST_2    0x05
	doIconst3,      // ICONST_3    0x06
	doIconst4,      // ICONST_4    0x07
	doIconst5,      // ICONST_5    0x08
	doLconst0,      // LCONST_0    0x09
	doLconst1,      // LCONST_1    0x0A
	doFconst0,      // FCONST_0    0x0B
	doFconst1,      // FCONST_1    0x0C
	doFconst2,      // FCONST_2    0x0D
	doDconst0,      // DCONST_0    0x0E
	doDconst1,      // DCONST_1    0x0F
	doBiPush,       // BIPUSH      0x10
	notImplemented, // SIPUSH      0x11
	notImplemented, // LDC         0x12
	notImplemented, // LDC_W       0x13
	notImplemented, // LDC2_W      0x14
	notImplemented, // ILOAD           0x15
	notImplemented, // LLOAD           0x16
	notImplemented, // FLOAD           0x17
	notImplemented, // DLOAD           0x18
	notImplemented, // ALOAD           0x19
	notImplemented, // ILOAD_0         0x1A
	notImplemented, // ILOAD_1         0x1B
	notImplemented, // ILOAD_2         0x1C
	notImplemented, // ILOAD_3         0x1D
	notImplemented, // LLOAD_0         0x1E
	notImplemented, // LLOAD_1         0x1F
	notImplemented, // LLOAD_2         0x20
	notImplemented, // LLOAD_3         0x21
	notImplemented, // FLOAD_0         0x22
	notImplemented, // FLOAD_1         0x23
	notImplemented, // FLOAD_2         0x24
	notImplemented, // FLOAD_3         0x25
	notImplemented, // DLOAD_0         0x26
	notImplemented, // DLOAD_1         0x27
	notImplemented, // DLOAD_2         0x28
	notImplemented, // DLOAD_3         0x29
	notImplemented, // ALOAD_0         0x2A
	notImplemented, // ALOAD_1         0x2B
	notImplemented, // ALOAD_2         0x2C
	notImplemented, // ALOAD_3         0x2D
	notImplemented, // IALOAD          0x2E
	notImplemented, // LALOAD          0x2F
	notImplemented, // FALOAD          0x30
	notImplemented, // DALOAD          0x31
	notImplemented, // AALOAD          0x32
	notImplemented, // BALOAD          0x33
	notImplemented, // CALOAD          0x34
	notImplemented, // SALOAD          0x35
	notImplemented, // ISTORE          0x36
	notImplemented, // LSTORE          0x37
	notImplemented, // FSTORE          0x38
	notImplemented, // DSTORE          0x39
	notImplemented, // ASTORE          0x3A
	notImplemented, // ISTORE_0        0x3B
	doIstore1,      // ISTORE_1        0x3C
	notImplemented, // ISTORE_2        0x3D
	notImplemented, // ISTORE_3        0x3E
	notImplemented, // LSTORE_0        0x3F
	notImplemented, // LSTORE_1        0x40
	notImplemented, // LSTORE_2        0x41
	notImplemented, // LSTORE_3        0x42
	notImplemented, // FSTORE_0        0x43
	notImplemented, // FSTORE_1        0x44
	notImplemented, // FSTORE_2        0x45
	notImplemented, // FSTORE_3        0x46
	notImplemented, // DSTORE_0        0x47
	notImplemented, // DSTORE_1        0x48
	notImplemented, // DSTORE_2        0x49
	notImplemented, // DSTORE_3        0x4A
	notImplemented, // ASTORE_0        0x4B
	notImplemented, // ASTORE_1        0x4C
	notImplemented, // ASTORE_2        0x4D
	notImplemented, // ASTORE_3        0x4E
	notImplemented, // IASTORE         0x4F
	notImplemented, // LASTORE         0x50
	notImplemented, // FASTORE         0x51
	notImplemented, // DASTORE         0x52
	notImplemented, // AASTORE         0x53
	notImplemented, // BASTORE         0x54
	notImplemented, // CASTORE         0x55
	notImplemented, // SASTORE         0x56
	notImplemented, // POP             0x57
	notImplemented, // POP2            0x58
	notImplemented, // DUP             0x59
	notImplemented, // DUP_X1          0x5A
	notImplemented, // DUP_X2          0x5B
	notImplemented, // DUP2            0x5C
	notImplemented, // DUP2_X1         0x5D
	notImplemented, // DUP2_X2         0x5E
	notImplemented, // SWAP            0x5F
	notImplemented, // IADD            0x60
	notImplemented, // LADD            0x61
	notImplemented, // FADD            0x62
	notImplemented, // DADD            0x63
	notImplemented, // ISUB            0x64
	notImplemented, // LSUB            0x65
	notImplemented, // FSUB            0x66
	notImplemented, // DSUB            0x67
	notImplemented, // IMUL            0x68
	notImplemented, // LMUL            0x69
	notImplemented, // FMUL            0x6A
	notImplemented, // DMUL            0x6B
	notImplemented, // IDIV            0x6C
	notImplemented, // LDIV            0x6D
	notImplemented, // FDIV            0x6E
	notImplemented, // DDIV            0x6F
	notImplemented, // IREM            0x70
	notImplemented, // LREM            0x71
	notImplemented, // FREM            0x72
	notImplemented, // DREM            0x73
	notImplemented, // INEG            0x74
	notImplemented, // LNEG            0x75
	notImplemented, // FNEG            0x76
	notImplemented, // DNEG            0x77
	notImplemented, // ISHL            0x78
	notImplemented, // LSHL            0x79
	notImplemented, // ISHR            0x7A
	notImplemented, // LSHR            0x7B
	notImplemented, // IUSHR           0x7C
	notImplemented, // LUSHR           0x7D
	notImplemented, // IAND            0x7E
	notImplemented, // LAND            0x7F
	notImplemented, // IOR             0x80
	notImplemented, // LOR             0x81
	notImplemented, // IXOR            0x82
	notImplemented, // LXOR            0x83
	notImplemented, // IINC            0x84
	notImplemented, // I2L             0x85
	notImplemented, // I2F             0x86
	notImplemented, // I2D             0x87
	notImplemented, // L2I             0x88
	notImplemented, // L2F             0x89
	notImplemented, // L2D             0x8A
	notImplemented, // F2I             0x8B
	notImplemented, // F2L             0x8C
	notImplemented, // F2D             0x8D
	notImplemented, // D2I             0x8E
	notImplemented, // D2L             0x8F
	notImplemented, // D2F             0x90
	notImplemented, // I2B             0x91
	notImplemented, // I2C             0x92
	notImplemented, // I2S             0x93
	notImplemented, // LCMP            0x94
	notImplemented, // FCMPL           0x95
	notImplemented, // FCMPG           0x96
	notImplemented, // DCMPL           0x97
	notImplemented, // DCMPG           0x98
	notImplemented, // IFEQ            0x99
	notImplemented, // IFNE            0x9A
	notImplemented, // IFLT            0x9B
	notImplemented, // IFGE            0x9C
	notImplemented, // IFGT            0x9D
	notImplemented, // IFLE            0x9E
	notImplemented, // IF_ICMPEQ       0x9F
	notImplemented, // IF_ICMPNE       0xA0
	notImplemented, // IF_ICMPLT       0xA1
	notImplemented, // IF_ICMPGE       0xA2
	notImplemented, // IF_ICMPGT       0xA3
	notImplemented, // IF_ICMPLE       0xA4
	notImplemented, // IF_ACMPEQ       0xA5
	notImplemented, // IF_ACMPNE       0xA6
	notImplemented, // GOTO            0xA7
	notImplemented, // JSR             0xA8
	notImplemented, // RET             0xA9
	notImplemented, // TABLESWITCH     0xAA
	notImplemented, // LOOKUPSWITCH    0xAB
	notImplemented, // IRETURN         0xAC
	notImplemented, // LRETURN         0xAD
	notImplemented, // FRETURN         0xAE
	notImplemented, // DRETURN         0xAF
	notImplemented, // ARETURN         0xB0
	notImplemented, // RETURN          0xB1
	notImplemented, // GETSTATIC       0xB2
	notImplemented, // PUTSTATIC       0xB3
	notImplemented, // GETFIELD        0xB4
	notImplemented, // PUTFIELD        0xB5
	notImplemented, // INVOKEVIRTUAL   0xB6
	notImplemented, // INVOKESPECIAL   0xB7
	notImplemented, // INVOKESTATIC    0xB8
	notImplemented, // INVOKEINTERFACE 0xB9
	notImplemented, // INVOKEDYNAMIC   0xBA
	notImplemented, // NEW             0xBB
	notImplemented, // NEWARRAY        0xBC
	notImplemented, // ANEWARRAY       0xBD
	notImplemented, // ARRAYLENGTH     0xBE
	notImplemented, // ATHROW          0xBF
	notImplemented, // CHECKCAST       0xC0
	notImplemented, // INSTANCEOF      0xC1
	notImplemented, // MONITORENTER    0xC2
	notImplemented, // MONITOREXIT     0xC3
	notImplemented, // WIDE            0xC4
	notImplemented, // MULTIANEWARRAY  0xC5
	notImplemented, // IFNULL          0xC6
	notImplemented, // IFNONNULL       0xC7
	notImplemented, // GOTO_W          0xC8
	notImplemented, // JSR_W           0xC9
	notImplemented, // BREAKPOINT      0xCA
}

// the functions, listed here in numerical order of the bytecode
func doNop(_ *frames.Frame, _ int64) int { return 1 }

func doAconstNull(fr *frames.Frame, _ int64) int {
	push(fr, object.Null)
	return 1
}

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

func doBiPush(fr *frames.Frame, _ int64) int {
	wbyte := fr.Meth[fr.PC+1]
	wint64 := byteToInt64(wbyte)
	push(fr, wint64)
	return 2
}

func doIstore1(fr *frames.Frame, _ int64) int { return storeInt(fr, int64(1)) }

func notImplemented(_ *frames.Frame, _ int64) int {
	return 1
}

// the functions call by the dispatched functions
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

func interpretBytecodes(bytecode int, f *frames.Frame) int {
	PC := DispatchTable[bytecode](f, 0)
	println("PC after call to DispatchTable[", bytecode, "] = ", PC)
	return PC
}
