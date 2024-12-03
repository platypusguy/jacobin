/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"encoding/binary"
	"errors"
	"fmt"
	"jacobin/excNames"
	"jacobin/globals"
	"jacobin/types"
	"math"
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
}

type BytecodeFunc func() int

var ERROR_OCCURRED = math.MaxInt32
var WideInEffect = false

var CheckTable = [203]BytecodeFunc{
	return1,           // NOP             0x00
	return1,           // ACONST_NULL     0x01
	return1,           // ICONST_M1       0x02
	return1,           // ICONST_0        0x03
	return1,           // ICONST_1        0x04
	return1,           // ICONST_2        0x05
	return1,           // ICONST_3        0x06
	return1,           // ICONST_4        0x07
	return1,           // ICONST_5        0x08
	return1,           // LCONST_0        0x09
	return1,           // LCONST_1        0x0A
	return1,           // FCONST_0        0x0B
	return1,           // FCONST_1        0x0C
	return1,           // FCONST_2        0x0D
	return1,           // DCONST_0        0x0E
	return1,           // DCONST_1        0x0F
	return2,           // BIPUSH          0x10
	return3,           // SIPUSH          0x11
	return2,           // LDC             0x12
	return3,           // LDC_W           0x13
	return3,           // LDC2_W          0x14
	return2,           // ILOAD           0x15
	return2,           // LLOAD           0x16
	return2,           // FLOAD           0x17
	return2,           // DLOAD           0x18
	return2,           // ALOAD           0x19
	return1,           // ILOAD_0         0x1A
	return1,           // ILOAD_1         0x1B
	return1,           // ILOAD_2         0x1C
	return1,           // ILOAD_3         0x1D
	return1,           // LLOAD_0         0x1E
	return1,           // LLOAD_1         0x1F
	return1,           // LLOAD_2         0x20
	return1,           // LLOAD_3         0x21
	return1,           // FLOAD_0         0x22
	return1,           // FLOAD_1         0x23
	return1,           // FLOAD_2         0x24
	return1,           // FLOAD_3         0x25
	return1,           // DLOAD_0         0x26
	return1,           // DLOAD_1         0x27
	return1,           // DLOAD_2         0x28
	return1,           // DLOAD_3         0x29
	return1,           // ALOAD_0         0x2A
	return1,           // ALOAD_1         0x2B
	return1,           // ALOAD_2         0x2C
	return1,           // ALOAD_3         0x2D
	return1,           // IALOAD          0x2E
	return1,           // LALOAD          0x2F
	return1,           // FALOAD          0x30
	return1,           // DALOAD          0x31
	return1,           // AALOAD          0x32
	return1,           // BALOAD          0x33
	return1,           // CALOAD          0x34
	return1,           // SALOAD          0x35
	return2,           // ISTORE          0x36
	return2,           // LSTORE          0x37
	return2,           // FSTORE          0x38
	return2,           // DSTORE          0x39
	return2,           // ASTORE          0x3A
	return1,           // ISTORE_0        0x3B
	return1,           // ISTORE_1        0x3C
	return1,           // ISTORE_2        0x3D
	return1,           // ISTORE_3        0x3E
	return1,           // LSTORE_0        0x3F
	return1,           // LSTORE_1        0x40
	return1,           // LSTORE_2        0x41
	return1,           // LSTORE_3        0x42
	return1,           // FSTORE_0        0x43
	return1,           // FSTORE_1        0x44
	return1,           // FSTORE_2        0x45
	return1,           // FSTORE_3        0x46
	return1,           // DSTORE_0        0x47
	return1,           // DSTORE_1        0x48
	return1,           // DSTORE_2        0x49
	return1,           // DSTORE_3        0x4A
	return1,           // ASTORE_0        0x4B
	return1,           // ASTORE_1        0x4C
	return1,           // ASTORE_2        0x4D
	return1,           // ASTORE_3        0x4E
	return1,           // IASTORE         0x4F
	return1,           // LASTORE         0x50
	return1,           // FASTORE         0x51
	return1,           // DASTORE         0x52
	return1,           // AASTORE         0x53
	return1,           // BASTORE         0x54
	return1,           // CASTORE         0x55
	return1,           // SASTORE         0x56
	return1,           // POP             0x57
	return1,           // POP2            0x58
	return1,           // DUP             0x59
	return1,           // DUP_X1          0x5A
	return1,           // DUP_X2          0x5B
	return1,           // DUP2            0x5C
	return1,           // DUP2_X1         0x5D
	return1,           // DUP2_X2         0x5E
	return1,           // SWAP            0x5F
	return1,           // IADD            0x60
	return1,           // LADD            0x61
	return1,           // FADD            0x62
	return1,           // DADD            0x63
	return1,           // ISUB            0x64
	return1,           // LSUB            0x65
	return1,           // FSUB            0x66
	return1,           // DSUB            0x67
	return1,           // IMUL            0x68
	return1,           // LMUL            0x69
	return1,           // FMUL            0x6A
	return1,           // DMUL            0x6B
	return1,           // IDIV            0x6C
	return1,           // LDIV            0x6D
	return1,           // FDIV            0x6E
	return1,           // DDIV            0x6F
	return1,           // IREM            0x70
	return1,           // LREM            0x71
	return1,           // FREM            0x72
	return1,           // DREM            0x73
	return1,           // INEG            0x74
	return1,           // LNEG            0x75
	return1,           // FNEG            0x76
	return1,           // DNEG            0x77
	return1,           // ISHL            0x78
	return1,           // LSHL            0x79
	return1,           // ISHR            0x7A
	return1,           // LSHR            0x7B
	return1,           // IUSHR           0x7C
	return1,           // LUSHR           0x7D
	return1,           // IAND            0x7E
	return1,           // LAND            0x7F
	return1,           // IOR             0x80
	return1,           // LOR             0x81
	return1,           // IXOR            0x82
	return1,           // LXOR            0x83
	return3,           // IINC            0x84
	return1,           // I2L             0x85
	return1,           // I2F             0x86
	return1,           // I2D             0x87
	return1,           // L2I             0x88
	return1,           // L2F             0x89
	return1,           // L2D             0x8A
	return1,           // F2I             0x8B
	return1,           // F2L             0x8C
	return1,           // F2D             0x8D
	return1,           // D2I             0x8E
	return1,           // D2L             0x8F
	return1,           // D2F             0x90
	return1,           // I2B             0x91
	return1,           // I2C             0x92
	return1,           // I2S             0x93
	return1,           // LCMP            0x94
	return1,           // FCMPL           0x95
	return1,           // FCMPG           0x96
	return1,           // DCMPL           0x97
	return1,           // DCMPG           0x98
	checkIf,           // IFEQ            0x99
	checkIf,           // IFNE            0x9A
	checkIf,           // IFLT            0x9B
	checkIf,           // IFGE            0x9C
	checkIf,           // IFGT            0x9D
	checkIf,           // IFLE            0x9E
	checkIf,           // IF_ICMPEQ       0x9F
	checkIf,           // IF_ICMPNE       0xA0
	checkIf,           // IF_ICMPLT       0xA1
	checkIf,           // IF_ICMPGE       0xA2
	checkIf,           // IF_ICMPGT       0xA3
	checkIf,           // IF_ICMPLE       0xA4
	checkIf,           // IF_ACMPEQ       0xA5
	checkIf,           // IF_ACMPNE       0xA6
	checkGoto,         // GOTO            0xA7
	checkGoto,         // JSR             0xA8
	return2,           // RET             0xA9
	checkTableswitch,  // TABLESWITCH     0xAA
	checkLookupswitch, // LOOKUPSWITCH    0xAB
	return1,           // IRETURN         0xAC
	return1,           // LRETURN         0xAD
	return1,           // FRETURN         0xAE
	return1,           // DRETURN         0xAF
	return1,           // ARETURN         0xB0
	return1,           // RETURN          0xB1
	return3,           // GETSTATIC       0xB2
	return3,           // PUTSTATIC       0xB3
	return3,           // GETFIELD        0xB4
	return3,           // PUTFIELD        0xB5
	return3,           // INVOKEVIRTUAL   0xB6
	return3,           // INVOKESPECIAL   0xB7
	return3,           // INVOKESTATIC    0xB8
	return5,           // INVOKEINTERFACE 0xB9
	return5,           // INVOKEDYNAMIC   0xBA
	return3,           // NEW             0xBB
	return2,           // NEWARRAY        0xBC
	return3,           // ANEWARRAY       0xBD
	return1,           // ARRAYLENGTH     0xBE
	return1,           // ATHROW          0xBF
	return3,           // CHECKCAST       0xC0
	return3,           // INSTANCEOF      0xC1
	return1,           // MONITORENTER    0xC2
	return1,           // MONITOREXIT     0xC3
	return1,           // WIDE            0xC4
	return4,           // MULTIANEWARRAY  0xC5
	return3,           // IFNULL          0xC6
	return3,           // IFNONNULL       0xC7
	checkGotow,        // GOTO_W          0xC8
	return5,           // JSR_W           0xC9
	return1,           // BREAKPOINT      0xCA
}

var PC int
var CP *CPool
var Code []byte

func CheckCodeValidity(code []byte, cp *CPool) error {
	// check that the code is valid
	if code == nil || cp == nil {
		errMsg := "CheckCodeValidity: nil code or constant pool"
		return errors.New(errMsg)
	}

	CP := cp
	if len(CP.CpIndex) == 0 {
		errMsg := "CheckCodeValidity: empty constant pool"
		return errors.New(errMsg)
	}

	Code = code
	PC = 0
	for PC < len(code) {
		opcode := code[PC]
		ret := CheckTable[opcode]()
		if ret == ERROR_OCCURRED {
			errMsg := fmt.Sprintf("Invalid bytecode or argument at location %d", PC)
			status := globals.GetGlobalRef().FuncThrowException(excNames.ClassFormatError, errMsg)
			if status != true { // will only happen in test
				globals.InitGlobals("test")
				return errors.New(errMsg)
			}
		} else {
			if ret+PC > len(code) {
				errMsg := fmt.Sprintf("Invalid bytecode or argument at location %d", PC)
				status := globals.GetGlobalRef().FuncThrowException(excNames.ClassFormatError, errMsg)
				if status != true { // will only happen in test
					globals.InitGlobals("test")
					return errors.New(errMsg)
				}
			}
			PC += ret
		}
	}
	return nil
}

// === check functions in alpha order by name of bytecode ===
func checkGoto() int {
	jumpTo := int(int16(Code[PC+1])*256 + int16(Code[PC+2]))
	if PC+jumpTo < 0 || PC+jumpTo >= len(Code) {
		return ERROR_OCCURRED
	}
	return 3
}

func checkGotow() int {
	jumpTo := int(types.FourBytesToInt64(Code[PC+1], Code[PC+2], Code[PC+3], Code[PC+4]))
	if PC+jumpTo < 0 || PC+jumpTo >= len(Code) {
		return ERROR_OCCURRED
	}
	return 5
}

func checkLookupswitch() int { // need to check this
	basePC := PC

	paddingBytes := 4 - ((PC + 1) % 4)
	if paddingBytes == 4 {
		paddingBytes = 0
	}
	basePC += paddingBytes
	basePC += 4 // jump size for default

	npairs := binary.BigEndian.Uint32(
		[]byte{Code[basePC+1], Code[basePC+2], Code[basePC+3], Code[basePC+4]})

	basePC += 4
	basePC += int(npairs) * 8

	return (basePC - PC) + 1
}

func checkTableswitch() int {
	basePC := PC

	paddingBytes := 4 - ((PC + 1) % 4)
	if paddingBytes == 4 {
		paddingBytes = 0
	}
	basePC += paddingBytes
	basePC += 12 // 4 bytes for default, 4 bytes for low, 4 bytes for high
	return (basePC - PC) + 1
}

func checkIf() int { // most IF* bytecodes come here. Jump if condition is met
	jumpSize := int(int16(Code[PC+1])*256 + int16(Code[PC+2]))
	if PC+jumpSize < 0 || PC+jumpSize >= len(Code) {
		return ERROR_OCCURRED
	}
	return 3
}

// === utility functions ===

// a one-byte opcode that has nothing that can be checked
func return1() int {
	return 1
}

func return2() int {
	return 2
}

func return3() int {
	return 3
}

func return4() int {
	return 4
}

func return5() int {
	return 5
}
