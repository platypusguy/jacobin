/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"errors"
)

// Here we check the bytecodes of a method. The method is passed as a byte slice.
// The bytecodes are checked against a function lookup table which calls the function
// to performs the check. It then uses the skip table to determine the number of bytes
// to skip to the next bytecode. If an error occurs, a ClassFormatException is thrown.

var bytecodeSkipTable = map[byte]int{
	0x00: 1, // NOP
	0x01: 1, // ACONST_NULL
	0x02: 1, // ICONST_M1
	0x03: 1, // ICONST_0
	0x04: 1, // ICONST_1
	0x05: 1, // ICONST_2
	0x06: 1, // ICONST_3
	0x07: 1, // ICONST_4
	0x08: 1, // ICONST_5
	0x09: 1, // LCONST_0
	0x0A: 1, // LCONST_1
	0x0B: 1, // FCONST_0
	0x0C: 1, // FCONST_1
	0x0D: 1, // FCONST_2
	0x0E: 1, // DCONST_0
	0x0F: 1, // DCONST_1
	0x10: 2, // BIPUSH
	0x11: 3, // SIPUSH
	0x12: 2, // LDC
	0x13: 3, // LDC_W
	0x14: 3, // LDC2_W
	0x15: 2, // ILOAD
	0x16: 2, // LLOAD
	0x17: 2, // FLOAD
	0x18: 2, // DLOAD
	0x19: 2, // ALOAD
	0x1A: 1, // ILOAD_0
	0x1B: 1, // ILOAD_1
	0x1C: 1, // ILOAD_2
	0x1D: 1, // ILOAD_3
	0x1E: 1, // LLOAD_0
	0x1F: 1, // LLOAD_1
	0x20: 1, // LLOAD_2
	0x21: 1, // LLOAD_3
	0x22: 1, // FLOAD_0
	0x23: 1, // FLOAD_1
	0x24: 1, // FLOAD_2
	0x25: 1, // FLOAD_3
	0x26: 1, // DLOAD_0
	0x27: 1, // DLOAD_1
	0x28: 1, // DLOAD_2
	0x29: 1, // DLOAD_3
	0x2A: 1, // ALOAD_0
	0x2B: 1, // ALOAD_1
	0x2C: 1, // ALOAD_2
	0x2D: 1, // ALOAD_3
	0x2E: 1, // IALOAD
	0x2F: 1, // LALOAD
	0x30: 1, // FALOAD
	0x31: 1, // DALOAD
	0x32: 1, // AALOAD
	0x33: 1, // BALOAD
	0x34: 1, // CALOAD
	0x35: 1, // SALOAD
	0x36: 2, // ISTORE
	0x37: 2, // LSTORE
	0x38: 2, // FSTORE
	0x39: 2, // DSTORE
	0x3A: 2, // ASTORE
	0x3B: 1, // ISTORE_0
	0x3C: 1, // ISTORE_1
	0x3D: 1, // ISTORE_2
	0x3E: 1, // ISTORE_3
	0x3F: 1, // LSTORE_0
	0x40: 1, // LSTORE_1
	0x41: 1, // LSTORE_2
	0x42: 1, // LSTORE_3
	0x43: 1, // FSTORE_0
	0x44: 1, // FSTORE_1
	0x45: 1, // FSTORE_2
	0x46: 1, // FSTORE_3
	0x47: 1, // DSTORE_0
	0x48: 1, // DSTORE_1
	0x49: 1, // DSTORE_2
	0x4A: 1, // DSTORE_3
	0x4B: 1, // ASTORE_0
	0x4C: 1, // ASTORE_1
	0x4D: 1, // ASTORE_2
	0x4E: 1, // ASTORE_3
	0x4F: 1, // IASTORE
	0x50: 1, // LASTORE
	0x51: 1, // FASTORE
	0x52: 1, // DASTORE
	0x53: 1, // AASTORE
	0x54: 1, // BASTORE
	0x55: 1, // CASTORE
	0x56: 1, // SASTORE
	0x57: 1, // POP
	0x58: 1, // POP2
	0x59: 1, // DUP
	0x5A: 1, // DUP_X1
	0x5B: 1, // DUP_X2
	0x5C: 1, // DUP2
	0x5D: 1, // DUP2_X1
	0x5E: 1, // DUP2_X2
	0x5F: 1, // SWAP
	0x60: 1, // IADD
	0x61: 1, // LADD
	0x62: 1, // FADD
	0x63: 1, // DADD
	0x64: 1, // ISUB
	0x65: 1, // LSUB
	0x66: 1, // FSUB
	0x67: 1, // DSUB
	0x68: 1, // IMUL
	0x69: 1, // LMUL
	0x6A: 1, // FMUL
	0x6B: 1, // DMUL
	0x6C: 1, // IDIV
	0x6D: 1, // LDIV
	0x6E: 1, // FDIV
	0x6F: 1, // DDIV
	0x70: 1, // IREM
	0x71: 1, // LREM
	0x72: 1, // FREM
	0x73: 1, // DREM
	0x74: 1, // INEG
	0x75: 1, // LNEG
	0x76: 1, // FNEG
	0x77: 1, // DNEG
	0x78: 1, // ISHL
	0x79: 1, // LSHL
	0x7A: 1, // ISHR
	0x7B: 1, // LSHR
	0x7C: 1, // IUSHR
	0x7D: 1, // LUSHR
	0x7E: 1, // IAND
	0x7F: 1, // LAND
	0x80: 1, // IOR
	0x81: 1, // LOR
	0x82: 1, // IXOR
	0x83: 1, // LXOR
	0x84: 3, // IINC
	0x85: 1, // I2L
	0x86: 1, // I2F
	0x87: 1, // I2D
	0x88: 1, // L2I
	0x89: 1, // L2F
	0x8A: 1, // L2D
	0x8B: 1, // F2I
	0x8C: 1, // F2L
	0x8D: 1, // F2D
	0x8E: 1, // D2I
	0x8F: 1, // D2L
	0x90: 1, // D2F
	0x91: 1, // I2B
	0x92: 1, // I2C
	0x93: 1, // I2S
	0x94: 1, // LCMP
	0x95: 1, // FCMPL
	0x96: 1, // FCMPG
	0x97: 1, // DCMPL
	0x98: 1, // DCMPG
	0x99: 3, // IFEQ
	0x9A: 3, // IFNE
	0x9B: 3, // IFLT
	0x9C: 3, // IFGE
	0x9D: 3, // IFGT
	0x9E: 3, // IFLE
	0x9F: 3, // IF_ICMPEQ
	0xA0: 3, // IF_ICMPNE
	0xA1: 3, // IF_ICMPLT
	0xA2: 3, // IF_ICMPGE
	0xA3: 3, // IF_ICMPGT
	0xA4: 3, // IF_ICMPLE
	0xA5: 3, // IF_ACMPEQ
	0xA6: 3, // IF_ACMPNE
	0xA7: 3, // GOTO
	0xA8: 3, // JSR
	0xA9: 2, // RET
	0xAA: 0, // TABLESWITCH
	0xAB: 0, // LOOKUPSWITCH
	0xAC: 1, // IRETURN
	0xAD: 1, // LRETURN
	0xAE: 1, // FRETURN
	0xAF: 1, // DRETURN
	0xB0: 1, // ARETURN
	0xB1: 1, // RETURN
	0xB2: 3, // GETSTATIC
	0xB3: 3, // PUTSTATIC
	0xB4: 3, // GETFIELD
	0xB5: 3, // PUTFIELD
	0xB6: 3, // INVOKEVIRTUAL
	0xB7: 3, // INVOKESPECIAL
	0xB8: 3, // INVOKESTATIC
	0xB9: 5, // INVOKEINTERFACE
	0xBA: 5, // INVOKEDYNAMIC
	0xBB: 3, // NEW
	0xBC: 2, // NEWARRAY
	0xBD: 3, // ANEWARRAY
	0xBE: 1, // ARRAYLENGTH
	0xBF: 1, // ATHROW
	0xC0: 3, // CHECKCAST
	0xC1: 3, // INSTANCEOF
	0xC2: 1, // MONITORENTER
	0xC3: 1, // MONITOREXIT
	0xC4: 0, // WIDE
	0xC5: 4, // MULTIANEWARRAY
	0xC6: 3, // IFNULL
	0xC7: 3, // IFNONNULL
	0xC8: 5, // GOTO_W
	0xC9: 5, // JSR_W
	0xCA: 1, // BREAKPOINT

	//
	// // Example usage
	// bytecode := byte(0x10) // bipush
	// skip := bytecodeSkipTable[bytecode]
	// fmt.Printf("Bytecode: 0x%X, Skip: %d\n", bytecode, skip)
}

type BytecodeFunc func(cp *CPool) error

var DispatchTable = [203]BytecodeFunc{
	// checkNop, // NOP             0x00
	/*	checkAconstNull,      // ACONST_NULL     0x01
		checkIconstM1,        // ICONST_M1       0x02
		checkIconst0,         // ICONST_0        0x03
		checkIconst1,         // ICONST_1        0x04
		checkIconst2,         // ICONST_2        0x05
		checkIconst3,         // ICONST_3        0x06
		checkIconst4,         // ICONST_4        0x07
		checkIconst5,         // ICONST_5        0x08
		checkLconst0,         // LCONST_0        0x09
		checkLconst1,         // LCONST_1        0x0A
		checkFconst0,         // FCONST_0        0x0B
		checkFconst1,         // FCONST_1        0x0C
		checkFconst2,         // FCONST_2        0x0D
		checkDconst0,         // DCONST_0        0x0E
		checkDconst1,         // DCONST_1        0x0F
		checkBipush,          // BIPUSH          0x10
		checkSipush,          // SIPUSH          0x11
		checkLdc,             // LDC             0x12
		checkLdcw,            // LDC_W           0x13
		checkLdc2w,           // LDC2_W          0x14
		checkLoad,            // ILOAD           0x15
		checkLoad,            // LLOAD           0x16
		checkLoad,            // FLOAD           0x17
		checkLoad,            // DLOAD           0x18
		checkLoad,            // ALOAD           0x19
		checkIload0,          // ILOAD_0         0x1A
		checkIload1,          // ILOAD_1         0x1B
		checkIload2,          // ILOAD_2         0x1C
		checkIload3,          // ILOAD_3         0x1D
		checkIload0,          // LLOAD_0         0x1E
		checkIload1,          // LLOAD_1         0x1F
		checkIload2,          // LLOAD_2         0x20
		checkIload3,          // LLOAD_3         0x21
		checkFload0,          // FLOAD_0         0x22
		checkFload1,          // FLOAD_1         0x23
		checkFload2,          // FLOAD_2         0x24
		checkFload3,          // FLOAD_3         0x25
		checkFload0,          // DLOAD_0         0x26
		checkFload1,          // DLOAD_1         0x27
		checkFload2,          // DLOAD_2         0x28
		checkFload3,          // DLOAD_3         0x29
		checkAload0,          // ALOAD_0         0x2A
		checkAload1,          // ALOAD_1         0x2B
		checkAload2,          // ALOAD_2         0x2C
		checkAload3,          // ALOAD_3         0x2D
		checkIaload,          // IALOAD          0x2E
		checkIaload,          // LALOAD          0x2F
		checkFaload,          // FALOAD          0x30
		checkFaload,          // DALOAD          0x31
		checkAaload,          // AALOAD          0x32
		checkBaload,          // BALOAD          0x33
		checkIaload,          // CALOAD          0x34
		checkIaload,          // SALOAD          0x35
		checkIstore,          // ISTORE          0x36
		checkIstore,          // LSTORE          0x37
		checkFstore,          // FSTORE          0x38
		checkFstore,          // DSTORE          0x39
		checkAstore,          // ASTORE          0x3A
		checkIstore0,         // ISTORE_0        0x3B
		checkIstore1,         // ISTORE_1        0x3C
		checkIstore2,         // ISTORE_2        0x3D
		checkIstore3,         // ISTORE_3        0x3E
		checkIstore0,         // LSTORE_0        0x3F
		checkIstore1,         // LSTORE_1        0x40
		checkIstore2,         // LSTORE_2        0x41
		checkIstore3,         // LSTORE_3        0x42
		checkFstore0,         // FSTORE_0        0x43
		checkFstore1,         // FSTORE_1        0x44
		checkFstore2,         // FSTORE_2        0x45
		checkFstore3,         // FSTORE_3        0x46
		checkFstore0,         // DSTORE_0        0x47
		checkFstore1,         // DSTORE_1        0x48
		checkFstore2,         // DSTORE_2        0x49
		checkFstore3,         // DSTORE_3        0x4A
		checkAstore0,         // ASTORE_0        0x4B
		checkAstore1,         // ASTORE_1        0x4C
		checkAstore2,         // ASTORE_2        0x4D
		checkAstore3,         // ASTORE_3        0x4E
		checkIastore,         // IASTORE         0x4F
		checkIastore,         // LASTORE         0x50
		checkFastore,         // FASTORE         0x51
		checkFastore,         // DASTORE         0x52
		checkAastore,         // AASTORE         0x53
		checkBastore,         // BASTORE         0x54
		checkIastore,         // CASTORE         0x55
		checkIastore,         // SASTORE         0x56
		checkPop,             // POP             0x57
		checkPop2,            // POP2            0x58
		checkDup,             // DUP             0x59
		checkDupx1,           // DUP_X1          0x5A
		checkDupx2,           // DUP_X2          0x5B
		checkDup2,            // DUP2            0x5C
		checkDup2x1,          // DUP2_X1         0x5D
		checkDup2x2,          // DUP2_X2         0x5E
		checkSwap,            // SWAP            0x5F
		checkIadd,            // IADD            0x60
		checkIadd,            // LADD            0x61
		checkFadd,            // FADD            0x62
		checkFadd,            // DADD            0x63
		checkIsub,            // ISUB            0x64
		checkIsub,            // LSUB            0x65
		checkFsub,            // FSUB            0x66
		checkFsub,            // DSUB            0x67
		checkImul,            // IMUL            0x68
		checkImul,            // LMUL            0x69
		checkFmul,            // FMUL            0x6A
		checkFmul,            // DMUL            0x6B
		checkIdiv,            // IDIV            0x6C
		checkIdiv,            // LDIV            0x6D
		checkFdiv,            // FDIV            0x6E
		checkFdiv,            // DDIV            0x6F
		checkIrem,            // IREM            0x70
		checkIrem,            // LREM            0x71
		checkFrem,            // FREM            0x72
		checkFrem,            // DREM            0x73
		checkIneg,            // INEG            0x74
		checkIneg,            // LNEG            0x75
		checkFneg,            // FNEG            0x76
		checkFneg,            // DNEG            0x77
		checkIshl,            // ISHL            0x78
		checkIshl,            // LSHL            0x79
		checkIshr,            // ISHR            0x7A
		checkIshr,            // LSHR            0x7B
		checkIushr,           // IUSHR           0x7C
		checkIushr,           // LUSHR           0x7D
		checkIand,            // IAND            0x7E
		checkIand,            // LAND            0x7F
		checkIor,             // IOR             0x80
		checkIor,             // LOR             0x81
		checkIxor,            // IXOR            0x82
		checkIxor,            // LXOR            0x83
		checkIinc,            // IINC            0x84
		checkNothing,         // I2L             0x85
		checkI2f,             // I2F             0x86
		checkI2f,             // I2D             0x87
		checkNothing,         // L2I             0x88
		checkL2f,             // L2F             0x89
		checkL2f,             // L2D             0x8A
		checkF2i,             // F2I             0x8B
		checkF2i,             // F2L             0x8C
		checkNothing,         // F2D             0x8D
		checkD2i,             // D2I             0x8E
		checkD2i,             // D2L             0x8F
		checkNothing,         // D2F             0x90
		checkI2b,             // I2B             0x91
		checkI2c,             // I2C             0x92
		checkI2s,             // I2S             0x93
		checkLcmp,            // LCMP            0x94
		checkFcmpl,           // FCMPL           0x95
		checkFcmpl,           // FCMPG           0x96
		checkFcmpl,           // DCMPL           0x97
		checkFcmpl,           // DCMPG           0x98
		checkIfeq,            // IFEQ            0x99
		checkIfne,            // IFNE            0x9A
		checkIflt,            // IFLT            0x9B
		checkIfge,            // IFGE            0x9C
		checkIfgt,            // IFGT            0x9D
		checkIfle,            // IFLE            0x9E
		checkIficmpeq,        // IF_ICMPEQ       0x9F
		checkIficmpne,        // IF_ICMPNE       0xA0
		checkIficmplt,        // IF_ICMPLT       0xA1
		checkIficmpge,        // IF_ICMPGE       0xA2
		checkIficmpgt,        // IF_ICMPGT       0xA3
		checkIficmple,        // IF_ICMPLE       0xA4
		checkIfacmpeq,        // IF_ACMPEQ       0xA5
		checkIfacmpne,        // IF_ACMPNE       0xA6
		checkGoto,            // GOTO            0xA7
		checkJsr,             // JSR             0xA8
		checkRet,             // RET             0xA9
		checkTableswitch,     // TABLESWITCH     0xAA
		checkLookupswitch,    // LOOKUPSWITCH    0xAB
		checkIreturn,         // IRETURN         0xAC
		checkIreturn,         // LRETURN         0xAD
		checkIreturn,         // FRETURN         0xAE
		checkIreturn,         // DRETURN         0xAF
		checkIreturn,         // ARETURN         0xB0
		checkReturn,          // RETURN          0xB1
		checkGetstatic,       // GETSTATIC       0xB2
		checkPutstatic,       // PUTSTATIC       0xB3
		checkGetfield,        // GETFIELD        0xB4
		checkPutfield,        // PUTFIELD        0xB5
		checkInvokeVirtual,   // INVOKEVIRTUAL   0xB6
		checkInvokespecial,   // INVOKESPECIAL   0xB7
		checkInvokestatic,    // INVOKESTATIC    0xB8
		checkInvokeinterface, // INVOKEINTERFACE 0xB9
		checkInvokedynamic,   // INVOKEDYNAMIC   0xBA
		checkNew,             // NEW             0xBB
		checkNewarray,        // NEWARRAY        0xBC
		checkAnewarray,       // ANEWARRAY       0xBD
		checkArraylength,     // ARRAYLENGTH     0xBE
		checkAthrow,          // ATHROW          0xBF
		checkCheckcast,       // CHECKCAST       0xC0
		checkInstanceof,      // INSTANCEOF      0xC1
		checkPop,             // MONITORENTER    0xC2
		checkPop,             // MONITOREXIT     0xC3
		checkWide,            // WIDE            0xC4
		checkMultinewarray,   // MULTIANEWARRAY  0xC5
		checkIfnull,          // IFNULL          0xC6
		checkIfnonnull,       // IFNONNULL       0xC7
		checkGotow,           // GOTO_W          0xC8
		checkJsrw,            // JSR_W           0xC9
		checkWarninvalid,     // BREAKPOINT      0xCA
	*/
}

func CheckCodeValidity(code []byte, cp *CPool) error {
	// check that the code is valid
	if code == nil || cp == nil {
		errMsg := "CheckCodeValidity: nil code or constant pool"
		return errors.New(errMsg)
	}

	CP := *cp
	if len(CP.CpIndex) == 0 {
		errMsg := "CheckCodeValidity: empty constant pool"
		return errors.New(errMsg)
	}

	return nil
}
