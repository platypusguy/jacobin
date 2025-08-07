/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024-5 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"encoding/binary"
	"errors"
	"fmt"
	"jacobin/excNames"
	"jacobin/globals"
	"jacobin/trace"
	"jacobin/types"
	"math"
)

// Here we check the bytecodes of a method. The method is passed as a byte slice.
// The bytecodes are checked against a function lookup table which calls the function
// to performs the check. It then uses the skip table to determine the number of bytes
// to skip to the next bytecode. If an error occurs, a ClassFormatException is thrown.

// NOTE: The unit tests for these functions are in codeCheck_test.go in the jvm directory.
// Placed there to avoid circular dependencies.

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
	return1,              // NOP             0x00
	checkAconstnull,      // ACONST_NULL     0x01
	pushInt,              // ICONST_M1       0x02
	pushInt,              // ICONST_0        0x03
	pushInt,              // ICONST_1        0x04
	pushInt,              // ICONST_2        0x05
	pushInt,              // ICONST_3        0x06
	pushInt,              // ICONST_4        0x07
	pushInt,              // ICONST_5        0x08
	pushInt,              // LCONST_0        0x09
	pushInt,              // LCONST_1        0x0A
	pushFloat,            // FCONST_0        0x0B
	pushFloat,            // FCONST_1        0x0C
	pushFloat,            // FCONST_2        0x0D
	pushFloat,            // DCONST_0        0x0E
	pushFloat,            // DCONST_1        0x0F
	checkBipush,          // BIPUSH          0x10
	checkSipush,          // SIPUSH          0x11
	pushIntRet2,          // LDC             0x12
	pushIntRet3,          // LDC_W           0x13
	pushIntRet3,          // LDC2_W          0x14
	pushIntRet2,          // ILOAD           0x15
	pushIntRet2,          // LLOAD           0x16
	pushFloatRet2,        // FLOAD           0x17
	pushFloatRet2,        // DLOAD           0x18
	pushIntRet2,          // ALOAD           0x19
	pushInt,              // ILOAD_0         0x1A
	pushInt,              // ILOAD_1         0x1B
	pushInt,              // ILOAD_2         0x1C
	pushInt,              // ILOAD_3         0x1D
	pushInt,              // LLOAD_0         0x1E
	pushInt,              // LLOAD_1         0x1F
	pushInt,              // LLOAD_2         0x20
	pushInt,              // LLOAD_3         0x21
	pushFloat,            // FLOAD_0         0x22
	pushFloat,            // FLOAD_1         0x23
	pushFloat,            // FLOAD_2         0x24
	pushFloat,            // FLOAD_3         0x25
	pushFloat,            // DLOAD_0         0x26
	pushFloat,            // DLOAD_1         0x27
	pushFloat,            // DLOAD_2         0x28
	pushFloat,            // DLOAD_3         0x29
	pushInt,              // ALOAD_0         0x2A
	pushInt,              // ALOAD_1         0x2B
	pushInt,              // ALOAD_2         0x2C
	pushInt,              // ALOAD_3         0x2D
	pushInt,              // IALOAD          0x2E
	pushInt,              // LALOAD          0x2F
	pushInt,              // FALOAD          0x30
	pushInt,              // DALOAD          0x31
	pushInt,              // AALOAD          0x32
	pushInt,              // BALOAD          0x33
	pushInt,              // CALOAD          0x34
	pushInt,              // SALOAD          0x35
	storeIntRet2,         // ISTORE          0x36
	storeIntRet2,         // LSTORE          0x37
	storeFloatRet2,       // FSTORE          0x38
	storeFloatRet2,       // DSTORE          0x39
	storeIntRet2,         // ASTORE          0x3A
	storeInt,             // ISTORE_0        0x3B
	storeInt,             // ISTORE_1        0x3C
	storeInt,             // ISTORE_2        0x3D
	storeInt,             // ISTORE_3        0x3E
	storeInt,             // LSTORE_0        0x3F
	storeInt,             // LSTORE_1        0x40
	storeInt,             // LSTORE_2        0x41
	storeInt,             // LSTORE_3        0x42
	storeFloat,           // FSTORE_0        0x43
	storeFloat,           // FSTORE_1        0x44
	storeFloat,           // FSTORE_2        0x45
	storeFloat,           // FSTORE_3        0x46
	storeFloat,           // DSTORE_0        0x47
	storeFloat,           // DSTORE_1        0x48
	storeFloat,           // DSTORE_2        0x49
	storeFloat,           // DSTORE_3        0x4A
	storeInt,             // ASTORE_0        0x4B
	storeInt,             // ASTORE_1        0x4C
	storeInt,             // ASTORE_2        0x4D
	storeInt,             // ASTORE_3        0x4E
	storeInt,             // IASTORE         0x4F
	storeInt,             // LASTORE         0x50
	storeInt,             // FASTORE         0x51
	storeInt,             // DASTORE         0x52
	storeInt,             // AASTORE         0x53
	storeInt,             // BASTORE         0x54
	storeInt,             // CASTORE         0x55
	storeInt,             // SASTORE         0x56
	checkPop,             // POP             0x57
	checkPop2,            // POP2            0x58
	dup1,                 // DUP             0x59
	dup1,                 // DUP_X1          0x5A
	dup1,                 // DUP_X2          0x5B
	dup2,                 // DUP2            0x5C
	dup2,                 // DUP2_X1         0x5D
	dup2,                 // DUP2_X2         0x5E
	return1,              // SWAP            0x5F
	arith,                // IADD            0x60
	arith,                // LADD            0x61
	arith,                // FADD            0x62
	arith,                // DADD            0x63
	arith,                // ISUB            0x64
	arith,                // LSUB            0x65
	arith,                // FSUB            0x66
	arith,                // DSUB            0x67
	arith,                // IMUL            0x68
	arith,                // LMUL            0x69
	arith,                // FMUL            0x6A
	arith,                // DMUL            0x6B
	arith,                // IDIV            0x6C
	arith,                // LDIV            0x6D
	arith,                // FDIV            0x6E
	arith,                // DDIV            0x6F
	arith,                // IREM            0x70
	arith,                // LREM            0x71
	arith,                // FREM            0x72
	arith,                // DREM            0x73
	arith,                // INEG            0x74
	arith,                // LNEG            0x75
	arith,                // FNEG            0x76
	arith,                // DNEG            0x77
	arith,                // ISHL            0x78
	arith,                // LSHL            0x79
	arith,                // ISHR            0x7A
	arith,                // LSHR            0x7B
	arith,                // IUSHR           0x7C
	arith,                // LUSHR           0x7D
	arith,                // IAND            0x7E
	arith,                // LAND            0x7F
	arith,                // IOR             0x80
	arith,                // LOR             0x81
	arith,                // IXOR            0x82
	arith,                // LXOR            0x83
	return3,              // IINC            0x84
	return1,              // I2L             0x85
	return1,              // I2F             0x86
	return1,              // I2D             0x87
	return1,              // L2I             0x88
	return1,              // L2F             0x89
	return1,              // L2D             0x8A
	return1,              // F2I             0x8B
	return1,              // F2L             0x8C
	return1,              // F2D             0x8D
	return1,              // D2I             0x8E
	return1,              // D2L             0x8F
	return1,              // D2F             0x90
	return1,              // I2B             0x91
	return1,              // I2C             0x92
	return1,              // I2S             0x93
	arith,                // LCMP            0x94
	arith,                // FCMPL           0x95
	arith,                // FCMPG           0x96
	arith,                // DCMPL           0x97
	arith,                // DCMPG           0x98
	checkIfZero,          // IFEQ            0x99
	checkIfZero,          // IFNE            0x9A
	checkIfZero,          // IFLT            0x9B
	checkIfZero,          // IFGE            0x9C
	checkIfZero,          // IFGT            0x9D
	checkIfZero,          // IFLE            0x9E
	checkIf,              // IF_ICMPEQ       0x9F
	checkIf,              // IF_ICMPNE       0xA0
	checkIf,              // IF_ICMPLT       0xA1
	checkIf,              // IF_ICMPGE       0xA2
	checkIf,              // IF_ICMPGT       0xA3
	checkIf,              // IF_ICMPLE       0xA4
	checkIf,              // IF_ACMPEQ       0xA5
	checkIf,              // IF_ACMPNE       0xA6 // stack-checking code got this far
	checkGoto,            // GOTO            0xA7
	checkGoto,            // JSR             0xA8
	return2,              // RET             0xA9
	checkTableswitch,     // TABLESWITCH     0xAA
	checkLookupswitch,    // LOOKUPSWITCH    0xAB
	return1,              // IRETURN         0xAC
	return1,              // LRETURN         0xAD
	return1,              // FRETURN         0xAE
	return1,              // DRETURN         0xAF
	return1,              // ARETURN         0xB0
	return1,              // RETURN          0xB1
	return3,              // GETSTATIC       0xB2
	return3,              // PUTSTATIC       0xB3
	checkGetfield,        // GETFIELD        0xB4
	return3,              // PUTFIELD        0xB5
	checkInvokevirtual,   // INVOKEVIRTUAL   0xB6
	checkInvokespecial,   // INVOKESPECIAL   0xB7
	checkInvokestatic,    // INVOKESTATIC    0xB8
	checkInvokeinterface, // INVOKEINTERFACE 0xB9
	return5,              // INVOKEDYNAMIC   0xBA
	return3,              // NEW             0xBB
	return2,              // NEWARRAY        0xBC
	return3,              // ANEWARRAY       0xBD
	return1,              // ARRAYLENGTH     0xBE
	return1,              // ATHROW          0xBF
	return3,              // CHECKCAST       0xC0
	return3,              // INSTANCEOF      0xC1
	return1,              // MONITORENTER    0xC2
	return1,              // MONITOREXIT     0xC3
	return1,              // WIDE            0xC4
	checkMultianewarray,  // MULTIANEWARRAY  0xC5
	return3,              // IFNULL          0xC6
	return3,              // IFNONNULL       0xC7
	checkGotow,           // GOTO_W          0xC8
	return5,              // JSR_W           0xC9
	return1,              // BREAKPOINT      0xCA
}

var PC int
var prevPC int // allows us to view the preceding opcode if we need it for analysis
var CP *CPool
var Code []byte
var StackEntries int
var MaxStack int

func CheckCodeValidity(codePtr *[]byte, cp *CPool, maxStack int, access AccessFlags) error {
	if codePtr == nil {
		errMsg := "CheckCodeValidity: ptr to code segment is nil"
		return errors.New(errMsg)
	}
	code := *codePtr
	// check that the code is valid
	if code == nil {
		if access.ClassIsAbstract {
			return nil
		} else {
			errMsg := "CheckCodeValidity: Empty code segment"
			return errors.New(errMsg)
		}
	}

	if cp == nil {
		errMsg := "CheckCodeValidity: ptr to constant pool is nil"
		return errors.New(errMsg)
	}

	CP = cp
	if len(CP.CpIndex) == 0 {
		errMsg := "CheckCodeValidity: empty constant pool"
		return errors.New(errMsg)
	}

	Code = code
	PC = 0
	prevPC = -1 // -1 means no previous PC
	MaxStack = maxStack
	StackEntries = 0

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
			prevPC = PC
			PC += ret
		}
	}
	return nil
}

// === check functions in alpha order by name of bytecode ===

// *ADD (pop 2 values, add/sub/multiply/divide/modulo them, push result)
func arith() int {
	StackEntries -= 1
	return 1
}

// ACONST_NULL 0x01 Push null onto op stack
func checkAconstnull() int {
	StackEntries += 1
	return 1
}

// BIPUSH 0x10 Push following byte onto op stack
func checkBipush() int {
	if len(Code) >= PC+1 {
		StackEntries += 1
		return 2
	} else {
		return ERROR_OCCURRED
	}
}

func dup1() int {
	StackEntries += 1
	return 1
}

// DUP2 is like DUP, but duplicates the top 2 stack values-- frequenly generated for longs and doubles,
// which take up 2 stack entries on HotSpot and other OpenJDK JVMs. On Jacobin, doubles and longs
// take up 1 stack entry, so we need to be check whether the operation is on a double or long. If
// it is, then we convert DUP2 to DUP, which duplicates only the top stack entry.
func dup2() int {
	if byteCodeIsForLongOrDouble(Code[PC+1]) { // check if the next bytecode is for a long or double
		Code[PC] = 0x59 // change DUP2 to DUP
		StackEntries += 1
		return 1
	}
	StackEntries += 2
	return 1
}

// FCONST and DCONST Push a float onto the op stack
func pushFloat() int {
	StackEntries += 1
	return 1
}

// FLOAD* and DLOAD* Push an int or long from local onto op stack
func pushFloatRet2() int {
	StackEntries += 1
	return 2
}

// FSORE* and DSTORE*
func storeFloat() int {
	StackEntries -= 1
	return 1
}

func storeFloatRet2() int {
	StackEntries -= 1
	return 2
}

// ICONST* and LCONST Push an int or long onto op stack
func pushInt() int {
	StackEntries += 1
	return 1
}

// ILOAD* and LLOAD* Push an int or long from local onto op stack
func pushIntRet2() int {
	StackEntries += 1
	return 2
}

// for LDC variants (but not LDC itself)
func pushIntRet3() int {
	StackEntries += 1
	return 3
}

// GETFIELD 0xB4 Get field from object and push it onto the stack
func checkGetfield() int {
	// check that the index points to a field reference in the CP
	CPslot := (int(Code[PC+1]) * 256) + int(Code[PC+2]) // next 2 bytes point to CP entry
	if CPslot < 1 || CPslot >= len(CP.CpIndex) {
		return ERROR_OCCURRED
	}

	CPentry := CP.CpIndex[CPslot]
	if CPentry.Type != FieldRef {
		errMsg := fmt.Sprintf("%s:\n GETFIELD at %d: CP entry (%d) is not a field reference",
			excNames.JVMexceptionNames[excNames.VerifyError], PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	return 3
}

// GOTO 0xA7
func checkGoto() int {
	jumpTo := int(int16(Code[PC+1])*256 + int16(Code[PC+2]))
	if PC+jumpTo < 0 || PC+jumpTo >= len(Code) {
		errMsg := fmt.Sprintf("%s:\n GOTO at %d: illegal jump to %d",
			excNames.JVMexceptionNames[excNames.VerifyError], PC, PC+jumpTo)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	return 3
}

// GOTO_W 0xC8
func checkGotow() int {
	jumpTo := int(types.FourBytesToInt64(Code[PC+1], Code[PC+2], Code[PC+3], Code[PC+4]))
	if PC+jumpTo < 0 || PC+jumpTo >= len(Code) {
		errMsg := fmt.Sprintf("%s:\n GOTO_W at %d: illegal jump to %d",
			excNames.JVMexceptionNames[excNames.VerifyError], PC, PC+jumpTo)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	return 5
}

// IF_ACMPEQ 0xA5 (and the many other IF* bytecodes)
func checkIf() int { // most IF* bytecodes come here. Jump if condition is met
	jumpSize := int(int16(Code[PC+1])*256 + int16(Code[PC+2]))
	if PC+jumpSize < 0 || PC+jumpSize >= len(Code) {
		errMsg := fmt.Sprintf("%s:\n IF_ACMPEQ at %d: illegal jump to %d",
			excNames.JVMexceptionNames[excNames.VerifyError], PC, PC+jumpSize)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	return 3
}

func checkIfZero() int { // Jump if condition w.r.t 0 is met
	jumpSize := int(int16(Code[PC+1])*256 + int16(Code[PC+2]))
	if PC+jumpSize < 0 || PC+jumpSize >= len(Code) {
		errMsg := fmt.Sprintf("%s:\n IF* test at %d: illegal jump to %d",
			excNames.JVMexceptionNames[excNames.VerifyError], PC, PC+jumpSize)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	StackEntries -= 1
	return 3
}

// INVOKEINTERFACE 0xB9
func checkInvokeinterface() int {
	CPslot := (int(Code[PC+1]) * 256) + int(Code[PC+2]) // next 2 bytes point to CP entry
	if CPslot < 1 || CPslot >= len(CP.CpIndex) {
		return ERROR_OCCURRED
	}

	countByte := Code[PC+3]
	if countByte == 0 {
		return ERROR_OCCURRED
	}

	zeroByte := Code[PC+4]
	if zeroByte != 0 {
		return ERROR_OCCURRED
	}

	CPentry := CP.CpIndex[CPslot]
	if CPentry.Type != Interface {
		// because this is not a ClassFormatError, we output a trace error message here
		errMsg := fmt.Sprintf("%s:\n INVOKEINTERFACE at %d: CP entry (%d) is not an interface reference",
			excNames.JVMexceptionNames[excNames.VerifyError], PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	return 4
}

// INVOKESPECIAL 0xB7
func checkInvokespecial() int {
	// check that the index points to a method or interface reference in the CP
	CPslot := (int(Code[PC+1]) * 256) + int(Code[PC+2]) // next 2 bytes point to CP entry
	if CPslot < 1 || CPslot >= len(CP.CpIndex) {
		return ERROR_OCCURRED
	}

	CPentry := CP.CpIndex[CPslot]
	if CPentry.Type != MethodRef && CPentry.Type != Interface {
		// because this is not a ClassFormatError, we output a trace error message here
		errMsg := fmt.Sprintf("%s:\n INVOKESPECIAL at %d: CP entry (%d) is not a method or interface reference",
			excNames.JVMexceptionNames[excNames.VerifyError], PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	return 3
}

// INVOKESTATIC 0xB8
func checkInvokestatic() int {
	// check that the index points to a method or interface reference in the CP
	CPslot := (int(Code[PC+1]) * 256) + int(Code[PC+2]) // next 2 bytes point to CP entry
	if CPslot < 1 || CPslot >= len(CP.CpIndex) {
		return ERROR_OCCURRED
	}

	CPentry := CP.CpIndex[CPslot]
	if CPentry.Type != MethodRef && CPentry.Type != Interface {
		// because this is not a ClassFormatError, we output a trace message here
		errMsg := fmt.Sprintf("%s:\n INVOKESTATIC at %d: CP entry (%d) is not a method or interface reference",
			excNames.JVMexceptionNames[excNames.VerifyError], PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	return 3
}

// INVOKEVIRTUAL 0xB6
func checkInvokevirtual() int {
	// check that the index points to a method reference in the CP
	CPslot := (int(Code[PC+1]) * 256) + int(Code[PC+2]) // next 2 bytes point to CP entry
	if CPslot < 1 || CPslot >= len(CP.CpIndex) {
		return ERROR_OCCURRED
	}

	CPentry := CP.CpIndex[CPslot]
	if CPentry.Type != MethodRef {
		// because this is not a ClassFormatError, we emit a trace message here
		errMsg := fmt.Sprintf("%s:\n INVOKEVIRTUAL at %d: CP entry (%d) is not a method reference",
			excNames.JVMexceptionNames[excNames.VerifyError], PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	return 3
}

// ISTORE
func storeInt() int {
	StackEntries -= 1
	return 1
}

func storeIntRet2() int {
	StackEntries -= 1
	return 2
}

// LOOKUPSWITCH 0xAB
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

func checkPop() int {
	StackEntries -= 1
	return 1
}

func checkPop2() int {
	StackEntries -= 2
	return 1
}

// TABLESWITCH 0xAA
func checkTableswitch() int {
	basePC := PC

	paddingBytes := 4 - ((PC + 1) % 4)
	if paddingBytes == 4 {
		paddingBytes = 0
	}
	basePC += paddingBytes
	basePC += 4 // jump size for default
	low := types.FourBytesToInt64(Code[basePC+1], Code[basePC+2], Code[basePC+3], Code[basePC+4])
	high := types.FourBytesToInt64(Code[basePC+5], Code[basePC+6], Code[basePC+7], Code[basePC+8])
	basePC += 8 // 4 bytes for low, 4 bytes for high

	if !(low <= high) {
		return ERROR_OCCURRED
	}

	offsetsCount := high - low + 1
	basePC += int(offsetsCount) * 4
	return (basePC - PC) + 1
}

// MULTIANEWARRAY 0xC5 (create a multidimensional array)
func checkMultianewarray() int {
	CPslot := (int(Code[PC+1]) * 256) + int(Code[PC+2]) // next 2 bytes point to CP entry
	if CPslot < 1 || CPslot >= len(CP.CpIndex) {
		return ERROR_OCCURRED
	}

	CPentry := CP.CpIndex[CPslot]
	if CPentry.Type != ClassRef {
		// because this is not a ClassFormatError, we emit a trace message here
		errMsg := fmt.Sprintf("%s:\n MULTIANEWARRAY at %d: CP entry (%d) is not a class reference",
			excNames.JVMexceptionNames[excNames.VerifyError], PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	dimensions := Code[PC+3]
	if dimensions == 0 {
		return ERROR_OCCURRED
	}

	return 4
}

// SIPUSH 0x11 create int from next 2 bytes and push it
func checkSipush() int {
	if len(Code) >= PC+2 {
		StackEntries += 1
		return 3
	} else {
		return ERROR_OCCURRED
	}
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

func byteCodeIsForLongOrDouble(bytecode byte) bool {
	switch bytecode {
	case
		0x09, 0x0A, 0x0E, 0x0F, 0x14, 0x16, 0x1E, 0x1F, // LCONST_0, LCONST_1, DCONST_0, DCONST_1, LDC2_W, LLOAD, LLOAD_0, LLOAD_1
		0x20, 0x21, 0x26, 0x27, 0x28, 0x29, 0x3F, 0x40, // LLOAD_2, LLOAD_3, DLOAD_0, DLOAD_1, DLOAD_2, DLOAD_3, LSTORE_0, LSTORE_1
		0x41, 0x42, 0x47, 0x48, 0x49, 0x4A, 0x62, 0x63, // LSTORE_2, LSTORE_3, DSTORE_0, DSTORE_1, DSTORE_2, DSTORE_3, LADD, DADD
		0x65, 0x66, 0x69, 0x6A, 0x6D, 0x6E, 0x71, 0x72, // LSUB, DSUB, LMUL, DMUL, LDIV, DDIV, LREM, DREM
		0x75, 0x76, 0x79, 0x7A, 0x7D, 0x7E, 0x81, 0x83, // LNEG, DNEG, LSHL, ISHR, LUSHR, LAND, LOR, LXOR
		0x94, 0x98: // LCMP, DCMPG
		return true
	default:
		return false
	}
}
