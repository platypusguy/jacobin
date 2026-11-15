/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024-6 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"encoding/binary"
	"errors"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/trace"
	"jacobin/src/types"
	"math"
	"os"
)

// Here we check the bytecodes of a method. The method is passed as a slice of bytes.
// The instructions are checked via a function lookup table that calls the verifier
// function. That function returns either an error code or the number of bytes to skip
// to reach the next bytecode. If an error occurs, a ClassFormatException is thrown.

var BytecodeNames = [203]string{
	"NOP",             // 0x00
	"ACONST_NULL",     // 0x01
	"ICONST_M1",       // 0x02
	"ICONST_0",        // 0x03
	"ICONST_1",        // 0x04
	"ICONST_2",        // 0x05
	"ICONST_3",        // 0x06
	"ICONST_4",        // 0x07
	"ICONST_5",        // 0x08
	"LCONST_0",        // 0x09
	"LCONST_1",        // 0x0A
	"FCONST_0",        // 0x0B
	"FCONST_1",        // 0x0C
	"FCONST_2",        // 0x0D
	"DCONST_0",        // 0x0E
	"DCONST_1",        // 0x0F
	"BIPUSH",          // 0x10
	"SIPUSH",          // 0x11
	"LDC",             // 0x12
	"LDC_W",           // 0x13
	"LDC2_W",          // 0x14
	"ILOAD",           // 0x15
	"LLOAD",           // 0x16
	"FLOAD",           // 0x17
	"DLOAD",           // 0x18
	"ALOAD",           // 0x19
	"ILOAD_0",         // 0x1A
	"ILOAD_1",         // 0x1B
	"ILOAD_2",         // 0x1C
	"ILOAD_3",         // 0x1D
	"LLOAD_0",         // 0x1E
	"LLOAD_1",         // 0x1F
	"LLOAD_2",         // 0x20
	"LLOAD_3",         // 0x21
	"FLOAD_0",         // 0x22
	"FLOAD_1",         // 0x23
	"FLOAD_2",         // 0x24
	"FLOAD_3",         // 0x25
	"DLOAD_0",         // 0x26
	"DLOAD_1",         // 0x27
	"DLOAD_2",         // 0x28
	"DLOAD_3",         // 0x29
	"ALOAD_0",         // 0x2A
	"ALOAD_1",         // 0x2B
	"ALOAD_2",         // 0x2C
	"ALOAD_3",         // 0x2D
	"IALOAD",          // 0x2E
	"LALOAD",          // 0x2F
	"FALOAD",          // 0x30
	"DALOAD",          // 0x31
	"AALOAD",          // 0x32
	"BALOAD",          // 0x33
	"CALOAD",          // 0x34
	"SALOAD",          // 0x35
	"ISTORE",          // 0x36
	"LSTORE",          // 0x37
	"FSTORE",          // 0x38
	"DSTORE",          // 0x39
	"ASTORE",          // 0x3A
	"ISTORE_0",        // 0x3B
	"ISTORE_1",        // 0x3C
	"ISTORE_2",        // 0x3D
	"ISTORE_3",        // 0x3E
	"LSTORE_0",        // 0x3F
	"LSTORE_1",        // 0x40
	"LSTORE_2",        // 0x41
	"LSTORE_3",        // 0x42
	"FSTORE_0",        // 0x43
	"FSTORE_1",        // 0x44
	"FSTORE_2",        // 0x45
	"FSTORE_3",        // 0x46
	"DSTORE_0",        // 0x47
	"DSTORE_1",        // 0x48
	"DSTORE_2",        // 0x49
	"DSTORE_3",        // 0x4A
	"ASTORE_0",        // 0x4B
	"ASTORE_1",        // 0x4C
	"ASTORE_2",        // 0x4D
	"ASTORE_3",        // 0x4E
	"IASTORE",         // 0x4F
	"LASTORE",         // 0x50
	"FASTORE",         // 0x51
	"DASTORE",         // 0x52
	"AASTORE",         // 0x53
	"BASTORE",         // 0x54
	"CASTORE",         // 0x55
	"SASTORE",         // 0x56
	"POP",             // 0x57
	"POP2",            // 0x58
	"DUP",             // 0x59
	"DUP_X1",          // 0x5A
	"DUP_X2",          // 0x5B
	"DUP2",            // 0x5C
	"DUP2_X1",         // 0x5D
	"DUP2_X2",         // 0x5E
	"SWAP",            // 0x5F
	"IADD",            // 0x60
	"LADD",            // 0x61
	"FADD",            // 0x62
	"DADD",            // 0x63
	"ISUB",            // 0x64
	"LSUB",            // 0x65
	"FSUB",            // 0x66
	"DSUB",            // 0x67
	"IMUL",            // 0x68
	"LMUL",            // 0x69
	"FMUL",            // 0x6A
	"DMUL",            // 0x6B
	"IDIV",            // 0x6C
	"LDIV",            // 0x6D
	"FDIV",            // 0x6E
	"DDIV",            // 0x6F
	"IREM",            // 0x70
	"LREM",            // 0x71
	"FREM",            // 0x72
	"DREM",            // 0x73
	"INEG",            // 0x74
	"LNEG",            // 0x75
	"FNEG",            // 0x76
	"DNEG",            // 0x77
	"ISHL",            // 0x78
	"LSHL",            // 0x79
	"ISHR",            // 0x7A
	"LSHR",            // 0x7B
	"IUSHR",           // 0x7C
	"LUSHR",           // 0x7D
	"IAND",            // 0x7E
	"LAND",            // 0x7F
	"IOR",             // 0x80
	"LOR",             // 0x81
	"IXOR",            // 0x82
	"LXOR",            // 0x83
	"IINC",            // 0x84
	"I2L",             // 0x85
	"I2F",             // 0x86
	"I2D",             // 0x87
	"L2I",             // 0x88
	"L2F",             // 0x89
	"L2D",             // 0x8A
	"F2I",             // 0x8B
	"F2L",             // 0x8C
	"F2D",             // 0x8D
	"D2I",             // 0x8E
	"D2L",             // 0x8F
	"D2F",             // 0x90
	"I2B",             // 0x91
	"I2C",             // 0x92
	"I2S",             // 0x93
	"LCMP",            // 0x94
	"FCMPL",           // 0x95
	"FCMPG",           // 0x96
	"DCMPL",           // 0x97
	"DCMPG",           // 0x98
	"IFEQ",            // 0x99
	"IFNE",            // 0x9A
	"IFLT",            // 0x9B
	"IFGE",            // 0x9C
	"IFGT",            // 0x9D
	"IFLE",            // 0x9E
	"IF_ICMPEQ",       // 0x9F
	"IF_ICMPNE",       // 0xA0
	"IF_ICMPLT",       // 0xA1
	"IF_ICMPGE",       // 0xA2
	"IF_ICMPGT",       // 0xA3
	"IF_ICMPLE",       // 0xA4
	"IF_ACMPEQ",       // 0xA5
	"IF_ACMPNE",       // 0xA6
	"GOTO",            // 0xA7
	"JSR",             // 0xA8
	"RET",             // 0xA9
	"TABLESWITCH",     // 0xAA
	"LOOKUPSWITCH",    // 0xAB
	"IRETURN",         // 0xAC
	"LRETURN",         // 0xAD
	"FRETURN",         // 0xAE
	"DRETURN",         // 0xAF
	"ARETURN",         // 0xB0
	"RETURN",          // 0xB1
	"GETSTATIC",       // 0xB2
	"PUTSTATIC",       // 0xB3
	"GETFIELD",        // 0xB4
	"PUTFIELD",        // 0xB5
	"INVOKEVIRTUAL",   // 0xB6
	"INVOKESPECIAL",   // 0xB7
	"INVOKESTATIC",    // 0xB8
	"INVOKEINTERFACE", // 0xB9
	"INVOKEDYNAMIC",   // 0xBA
	"NEW",             // 0xBB
	"NEWARRAY",        // 0xBC
	"ANEWARRAY",       // 0xBD
	"ARRAYLENGTH",     // 0xBE
	"ATHROW",          // 0xBF
	"CHECKCAST",       // 0xC0
	"INSTANCEOF",      // 0xC1
	"MONITORENTER",    // 0xC2
	"MONITOREXIT",     // 0xC3
	"WIDE",            // 0xC4
	"MULTIANEWARRAY",  // 0xC5
	"IFNULL",          // 0xC6
	"IFNONNULL",       // 0xC7
	"GOTO_W",          // 0xC8
	"JSR_W",           // 0xC9
	"BREAKPOINT",      // 0xCA
}

type ccvBlock_type struct {
	PC           int
	PrevPC       int // allows us to view the preceding opcode if we need it for analysis
	CP           *CPool
	Code         []byte
	StackEntries int
	MaxStack     int
	MaxLocals    int
	wideInEffect bool
}

type BytecodeFunc func(*ccvBlock_type) int

const ERROR_OCCURRED = math.MaxInt32

var CheckTable = [256]BytecodeFunc{
	return1,               // NOP             0x00
	incrementStackEntries, // ACONST_NULL     0x01
	pushInt,               // ICONST_M1       0x02
	pushInt,               // ICONST_0        0x03
	pushInt,               // ICONST_1        0x04
	pushInt,               // ICONST_2        0x05
	pushInt,               // ICONST_3        0x06
	pushInt,               // ICONST_4        0x07
	pushInt,               // ICONST_5        0x08
	pushInt,               // LCONST_0        0x09
	pushInt,               // LCONST_1        0x0A
	pushFloat,             // FCONST_0        0x0B
	pushFloat,             // FCONST_1        0x0C
	pushFloat,             // FCONST_2        0x0D
	pushFloat,             // DCONST_0        0x0E
	pushFloat,             // DCONST_1        0x0F
	checkBipush,           // BIPUSH          0x10
	checkSipush,           // SIPUSH          0x11
	pushIntRet2,           // LDC             0x12
	pushIntRet3,           // LDC_W           0x13
	pushIntRet3,           // LDC2_W          0x14
	checkIload,            // ILOAD           0x15
	checkIload,            // LLOAD           0x16
	pushFloatRet2,         // FLOAD           0x17
	pushFloatRet2,         // DLOAD           0x18
	pushIntRet2,           // ALOAD           0x19
	checkIload0,           // ILOAD_0         0x1A
	checkIload1,           // ILOAD_1         0x1B
	checkIload2,           // ILOAD_2         0x1C
	checkIload3,           // ILOAD_3         0x1D
	checkLload0,           // LLOAD_0         0x1E
	checkLload1,           // LLOAD_1         0x1F
	checkLload2,           // LLOAD_2         0x20
	checkLload3,           // LLOAD_3         0x21
	checkIload0,           // FLOAD_0         0x22
	checkIload1,           // FLOAD_1         0x23
	checkIload2,           // FLOAD_2         0x24
	checkIload3,           // FLOAD_3         0x25
	checkLload0,           // DLOAD_0         0x26
	checkLload1,           // DLOAD_1         0x27
	checkLload2,           // DLOAD_2         0x28
	checkLload3,           // DLOAD_3         0x29
	pushInt,               // ALOAD_0         0x2A
	pushInt,               // ALOAD_1         0x2B
	pushInt,               // ALOAD_2         0x2C
	pushInt,               // ALOAD_3         0x2D
	pushInt,               // IALOAD          0x2E
	pushInt,               // LALOAD          0x2F
	pushInt,               // FALOAD          0x30
	pushInt,               // DALOAD          0x31
	pushInt,               // AALOAD          0x32
	pushInt,               // BALOAD          0x33
	pushInt,               // CALOAD          0x34
	pushInt,               // SALOAD          0x35
	storeIntRet2,          // ISTORE          0x36
	storeIntRet2,          // LSTORE          0x37
	storeFloatRet2,        // FSTORE          0x38
	storeFloatRet2,        // DSTORE          0x39
	storeIntRet2,          // ASTORE          0x3A
	storeInt,              // ISTORE_0        0x3B
	storeInt,              // ISTORE_1        0x3C
	storeInt,              // ISTORE_2        0x3D
	storeInt,              // ISTORE_3        0x3E
	storeInt,              // LSTORE_0        0x3F
	storeInt,              // LSTORE_1        0x40
	storeInt,              // LSTORE_2        0x41
	storeInt,              // LSTORE_3        0x42
	storeFloat,            // FSTORE_0        0x43
	storeFloat,            // FSTORE_1        0x44
	storeFloat,            // FSTORE_2        0x45
	storeFloat,            // FSTORE_3        0x46
	storeFloat,            // DSTORE_0        0x47
	storeFloat,            // DSTORE_1        0x48
	storeFloat,            // DSTORE_2        0x49
	storeFloat,            // DSTORE_3        0x4A
	storeInt,              // ASTORE_0        0x4B
	storeInt,              // ASTORE_1        0x4C
	storeInt,              // ASTORE_2        0x4D
	storeInt,              // ASTORE_3        0x4E
	storeInt,              // IASTORE         0x4F
	storeInt,              // LASTORE         0x50
	storeInt,              // FASTORE         0x51
	storeInt,              // DASTORE         0x52
	storeInt,              // AASTORE         0x53
	storeInt,              // BASTORE         0x54
	storeInt,              // CASTORE         0x55
	storeInt,              // SASTORE         0x56
	checkPop,              // POP             0x57
	checkPop2,             // POP2            0x58
	checkDup1,             // DUP             0x59
	checkDup1,             // DUP_X1          0x5A
	checkDup1,             // DUP_X2          0x5B
	checkDup2,             // DUP2            0x5C
	checkDup2x1,           // DUP2_X1         0x5D
	checkDup2x2,           // DUP2_X2         0x5E
	return1,               // SWAP            0x5F
	decrementStackEntries, // IADD            0x60
	decrementStackEntries, // LADD            0x61
	decrementStackEntries, // FADD            0x62
	decrementStackEntries, // DADD            0x63
	decrementStackEntries, // ISUB            0x64
	decrementStackEntries, // LSUB            0x65
	decrementStackEntries, // FSUB            0x66
	decrementStackEntries, // DSUB            0x67
	decrementStackEntries, // IMUL            0x68
	decrementStackEntries, // LMUL            0x69
	decrementStackEntries, // FMUL            0x6A
	decrementStackEntries, // DMUL            0x6B
	decrementStackEntries, // IDIV            0x6C
	decrementStackEntries, // LDIV            0x6D
	decrementStackEntries, // FDIV            0x6E
	decrementStackEntries, // DDIV            0x6F
	decrementStackEntries, // IREM            0x70
	decrementStackEntries, // LREM            0x71
	decrementStackEntries, // FREM            0x72
	decrementStackEntries, // DREM            0x73
	decrementStackEntries, // INEG            0x74
	decrementStackEntries, // LNEG            0x75
	decrementStackEntries, // FNEG            0x76
	decrementStackEntries, // DNEG            0x77
	decrementStackEntries, // ISHL            0x78
	decrementStackEntries, // LSHL            0x79
	decrementStackEntries, // ISHR            0x7A
	decrementStackEntries, // LSHR            0x7B
	decrementStackEntries, // IUSHR           0x7C
	decrementStackEntries, // LUSHR           0x7D
	decrementStackEntries, // IAND            0x7E
	decrementStackEntries, // LAND            0x7F
	decrementStackEntries, // IOR             0x80
	decrementStackEntries, // LOR             0x81
	decrementStackEntries, // IXOR            0x82
	decrementStackEntries, // LXOR            0x83
	return3,               // IINC            0x84
	return1,               // I2L             0x85
	return1,               // I2F             0x86
	return1,               // I2D             0x87
	return1,               // L2I             0x88
	return1,               // L2F             0x89
	return1,               // L2D             0x8A
	return1,               // F2I             0x8B
	return1,               // F2L             0x8C
	return1,               // F2D             0x8D
	return1,               // D2I             0x8E
	return1,               // D2L             0x8F
	return1,               // D2F             0x90
	return1,               // I2B             0x91
	return1,               // I2C             0x92
	return1,               // I2S             0x93
	decrementStackEntries, // LCMP            0x94
	decrementStackEntries, // FCMPL           0x95
	decrementStackEntries, // FCMPG           0x96
	decrementStackEntries, // DCMPL           0x97
	decrementStackEntries, // DCMPG           0x98
	checkIfzero,           // IFEQ            0x99
	checkIfzero,           // IFNE            0x9A
	checkIfzero,           // IFLT            0x9B
	checkIfzero,           // IFGE            0x9C
	checkIfzero,           // IFGT            0x9D
	checkIfzero,           // IFLE            0x9E
	checkIf,               // IF_ICMPEQ       0x9F
	checkIf,               // IF_ICMPNE       0xA0
	checkIf,               // IF_ICMPLT       0xA1
	checkIf,               // IF_ICMPGE       0xA2
	checkIf,               // IF_ICMPGT       0xA3
	checkIf,               // IF_ICMPLE       0xA4
	checkIf,               // IF_ACMPEQ       0xA5
	checkIf,               // IF_ACMPNE       0xA6
	checkGoto,             // GOTO            0xA7
	checkGoto,             // JSR             0xA8
	return2,               // RET             0xA9
	checkTableSwitch,      // TABLESWITCH     0xAA
	checkLookupswitch,     // LOOKUPSWITCH    0xAB
	return1,               // IRETURN         0xAC
	return1,               // LRETURN         0xAD
	return1,               // FRETURN         0xAE
	return1,               // DRETURN         0xAF
	return1,               // ARETURN         0xB0
	return1,               // RETURN          0xB1
	checkGetstatic,        // GETSTATIC       0xB2
	checkPutstatic,        // PUTSTATIC       0xB3
	checkGetfield,         // GETFIELD        0xB4
	checkPutfield,         // PUTFIELD        0xB5
	checkInvokevirtual,    // INVOKEVIRTUAL   0xB6
	checkInvokespecial,    // INVOKESPECIAL   0xB7
	checkInvokestatic,     // INVOKESTATIC    0xB8
	checkInvokeinterface,  // INVOKEINTERFACE 0xB9
	checkInvokedynamic,    // INVOKEDYNAMIC   0xBA
	checkNew,              // NEW             0xBB
	checkNewArray,         // NEWARRAY        0xBC
	checkAnewarray,        // ANEWARRAY       0xBD
	return1,               // ARRAYLENGTH     0xBE
	return1,               // ATHROW          0xBF
	return3,               // CHECKCAST       0xC0
	return3,               // INSTANCEOF      0xC1
	return1,               // MONITORENTER    0xC2
	return1,               // MONITOREXIT     0xC3
	checkWide,             // WIDE            0xC4
	checkMultianewarray,   // MULTIANEWARRAY  0xC5
	return3,               // IFNULL          0xC6
	return3,               // IFNONNULL       0xC7
	checkGotow,            // GOTO_W          0xC8
	return5,               // JSR_W           0xC9
	return1,               // BREAKPOINT      0xCA
	invalidBytecode,       // 0xCB All following bytecodes (thru 0xFF) are invalid in a class file
	invalidBytecode,       // 0xCC
	invalidBytecode,       // 0xCD
	invalidBytecode,       // 0xCE
	invalidBytecode,       // 0xCF
	invalidBytecode,       // 0xD0
	invalidBytecode,       // 0xD1
	invalidBytecode,       // 0xD2
	invalidBytecode,       // 0xD3
	invalidBytecode,       // 0xD4
	invalidBytecode,       // 0xD5
	invalidBytecode,       // 0xD6
	invalidBytecode,       // 0xD7
	invalidBytecode,       // 0xD8
	invalidBytecode,       // 0xD9
	invalidBytecode,       // 0xDA
	invalidBytecode,       // 0xDB
	invalidBytecode,       // 0xDC
	invalidBytecode,       // 0xDD
	invalidBytecode,       // 0xDE
	invalidBytecode,       // 0xDF
	invalidBytecode,       // 0xE0
	invalidBytecode,       // 0xE1
	invalidBytecode,       // 0xE2
	invalidBytecode,       // 0xE3
	invalidBytecode,       // 0xE4
	invalidBytecode,       // 0xE5
	invalidBytecode,       // 0xE6
	invalidBytecode,       // 0xE7
	invalidBytecode,       // 0xE8
	invalidBytecode,       // 0xE9
	invalidBytecode,       // 0xEA
	invalidBytecode,       // 0xEB
	invalidBytecode,       // 0xEC
	invalidBytecode,       // 0xED
	invalidBytecode,       // 0xEE
	invalidBytecode,       // 0xEF
	invalidBytecode,       // 0xF0
	invalidBytecode,       // 0xF1
	invalidBytecode,       // 0xF2
	invalidBytecode,       // 0xF3
	invalidBytecode,       // 0xF4
	invalidBytecode,       // 0xF5
	invalidBytecode,       // 0xF6
	invalidBytecode,       // 0xF7
	invalidBytecode,       // 0xF8
	invalidBytecode,       // 0xF9
	invalidBytecode,       // 0xFA
	invalidBytecode,       // 0xFB
	invalidBytecode,       // 0xFC
	invalidBytecode,       // 0xFD
	invalidBytecode,       // 0xFE
	invalidBytecode,       // 0xFF
}

func CheckCodeValidity(codePtr *[]byte, cp *CPool, maxStack int, maxLocals int, access AccessFlags, methNamePtr *string) error {

	ccvb := ccvBlock_type{}

	if codePtr == nil {
		errMsg := "CheckCodeValidity: ptr to code segment is nil"
		return errors.New(errMsg)
	}
	code := *codePtr
	// check that the code is valid
	if code == nil || len(code) == 0 {
		if access.ClassIsAbstract {
			return nil
		}
		errMsg := "CheckCodeValidity: Empty code segment"
		return errors.New(errMsg)
	}

	if cp == nil {
		errMsg := "CheckCodeValidity: ptr to constant pool is nil"
		return errors.New(errMsg)
	}

	ccvb.CP = cp
	if len(ccvb.CP.CpIndex) == 0 {
		errMsg := "CheckCodeValidity: empty constant pool"
		return errors.New(errMsg)
	}

	var methName string
	if methNamePtr == nil {
		methName = "not specified"
	} else {
		methName = *methNamePtr
	}

	ccvb.Code = code
	ccvb.PC = 0
	ccvb.PrevPC = -1 // -1 means no previous PC
	ccvb.MaxStack = maxStack
	ccvb.MaxLocals = maxLocals
	ccvb.StackEntries = 0
	ccvb.wideInEffect = false

	if globals.TraceCodeCheck {
		_, _ = fmt.Fprintf(os.Stderr, "\nCode check for method: %s\n", methName)
	}

	for ccvb.PC < len(code) {
		opcode := code[ccvb.PC]
		bcName := "UNKNOWN"
		if int(opcode) < len(BytecodeNames) {
			bcName = BytecodeNames[opcode]
		}
		if globals.TraceCodeCheck {
			_, _ = fmt.Fprintf(os.Stderr, "PC: %03d %-14s (%02X)\n", ccvb.PC, bcName, opcode)
		}

		ret := CheckTable[opcode](&ccvb)
		if ret == ERROR_OCCURRED {
			errMsg := fmt.Sprintf("Check of %s instruction: Invalid bytecode or argument at location %d",
				bcName, ccvb.PC)
			return errors.New(errMsg)
		}

		if ret+ccvb.PC > len(code) {
			errMsg := fmt.Sprintf("Invalid bytecode or argument at location %d", ccvb.PC)
			return errors.New(errMsg)
		}

		ccvb.PrevPC = ccvb.PC
		ccvb.PC += ret
	}
	return nil
}

// === check functions in alpha order by name of bytecode ===

// *ADD (pop 2 values, add/sub/multiply/divide/modulo them, push result)
func decrementStackEntries(ccvb *ccvBlock_type) int {
	ccvb.StackEntries -= 1
	return 1
}

// ACONST_NULL 0x01 Push null onto op stack
func incrementStackEntries(ccvb *ccvBlock_type) int {
	ccvb.StackEntries += 1
	return 1
}

// ANEWARRAY 0xBD Create a new array of reference type elements
func checkAnewarray(ccvb *ccvBlock_type) int {
	// check that the index points to an array of reference in the CP
	CPslot := (int(ccvb.Code[ccvb.PC+1]) * 256) + int(ccvb.Code[ccvb.PC+2]) // next 2 bytes point to CP entry
	if CPslot < 1 || CPslot >= len(ccvb.CP.CpIndex) {
		errMsg := fmt.Sprintf("%s:\n ANEWARRAY at %d: CPslot (%d) is invalid",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPslot)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	CPentry := ccvb.CP.CpIndex[CPslot]
	if CPentry.Type != ClassRef && CPentry.Type != Interface {
		errMsg := fmt.Sprintf("%s:\n ANEWARRAY at %d: CP entry (%d) is not a field or interface reference",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	ccvb.StackEntries += 1 // for the dimension of the array

	return 3
}

// BIPUSH 0x10 Push following byte onto op stack
func checkBipush(ccvb *ccvBlock_type) int {
	if len(ccvb.Code) > ccvb.PC+1 {
		ccvb.StackEntries += 1
		return 2
	}
	errMsg := fmt.Sprintf("%s:\n BIPUSH at %d: len(ccvb.Code) (%d) is exceeds ccvb.PC+1 (%d)",
		excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, len(ccvb.Code), ccvb.PC+1)
	trace.Error(errMsg)
	return ERROR_OCCURRED
}

// DUP  Push a duplicate of the top stack value onto the stack.
func checkDup1(ccvb *ccvBlock_type) int {
	ccvb.StackEntries += 1
	return 1
}

// DUP2 is like DUP, but duplicates the top 2 stack values-- frequenly generated for longs and doubles,
// which take up 2 stack entries on HotSpot and other OpenJDK JVMs. On Jacobin, doubles and longs
// take up 1 stack entry, so we need to check whether the operation is on a double or long. If
// it is, then we convert DUP2 to DUP, which duplicates only the top stack entry.
func checkDup2(ccvb *ccvBlock_type) int {
	if ccvb.PrevPC > -1 && bytecodePushes32BitValue(ccvb.Code[ccvb.PrevPC]) { // check if the previous bytecode is a 32-bit load bytecode
		goto dup2 // if so, we can safely use DUP2
	}

	// TODO: should we be testing the previous bytecode instead? That's what we do in Dup2x1
	if ccvb.PC+1 < len(ccvb.Code) && bytecodeIsForLongOrDouble(ccvb.Code[ccvb.PC+1]) { // check if the next bytecode is for a long or double
		ccvb.Code[ccvb.PC] = 0x59 // change DUP2 to DUP
		ccvb.StackEntries += 1
		return 1
	}
dup2:
	ccvb.StackEntries += 2
	return 1
}

// checkDup2x1 converts DUP2_X1 to DUP_X1 when it detects that the bytecode is handling a 64-bit value.
func checkDup2x1(ccvb *ccvBlock_type) int {
	if ccvb.PrevPC > -1 && bytecodePushes32BitValue(ccvb.Code[ccvb.PrevPC]) { // check if the previous bytecode is a 32-bit load bytecode
		goto dup2x1 // if so, we can safely use DUP2
	}

	if ccvb.PrevPC > -1 && bytecodeIsForLongOrDouble(ccvb.Code[ccvb.PrevPC]) { // check if the preceding bytecode is for a long or double
		ccvb.Code[ccvb.PC] = 0x5A // change DUP2_X1 to DUP_X1
		ccvb.StackEntries += 1
		return 1
	}

	if ccvb.PC+1 < len(ccvb.Code) && bytecodeIsForLongOrDouble(ccvb.Code[ccvb.PC+1]) { // check if the next bytecode is for a long or double
		ccvb.Code[ccvb.PC] = 0x5A // change DUP2_X1 to DUP_X1
		ccvb.StackEntries += 1
		return 1
	}
dup2x1:
	ccvb.StackEntries += 2
	return 1
}

// DUP2_X2: *NOTE* checkDup2x2 converts DUP2_X2 to DUP_X2 when it detects that the bytecode is handling a 64-bit value.
func checkDup2x2(ccvb *ccvBlock_type) int {
	if ccvb.PrevPC > -1 && bytecodePushes32BitValue(ccvb.Code[ccvb.PrevPC]) { // check if the previous bytecode is a 32-bit load bytecode
		goto dup2x2 // if so, we can safely use DUP2
	}

	if ccvb.PrevPC > -1 && bytecodeIsForLongOrDouble(ccvb.Code[ccvb.PrevPC]) { // check if the preceding bytecode is for a long or double
		ccvb.Code[ccvb.PC] = 0x5B // change DUP2_X2 to DUP_X2
		ccvb.StackEntries += 1
		return 1
	}

	if ccvb.PC+1 < len(ccvb.Code) && bytecodeIsForLongOrDouble(ccvb.Code[ccvb.PC+1]) { // check if the following bytecode is for a long or double
		ccvb.Code[ccvb.PC] = 0x5B // change DUP2_X2 to DUP_X2
		ccvb.StackEntries += 1
		return 1
	}
dup2x2:
	ccvb.StackEntries += 2
	return 1
}

// FCONST and DCONST Push a float onto the op stack
func pushFloat(ccvb *ccvBlock_type) int {
	ccvb.StackEntries += 1
	return 1
}

// FLOAD* and DLOAD* Push a float or double from local onto op stack
func pushFloatRet2(ccvb *ccvBlock_type) int {
	ccvb.StackEntries += 1
	return 2
}

// FSTORE* and DSTORE*
func storeFloat(ccvb *ccvBlock_type) int {
	ccvb.StackEntries -= 1
	return 1
}

func storeFloatRet2(ccvb *ccvBlock_type) int {
	ccvb.StackEntries -= 1
	return 2
}

// GETFIELD 0xB4 Get non-static field from object and push it onto the stack
func checkGetfield(ccvb *ccvBlock_type) int {
	// check that the index points to a field reference in the CP
	CPslot := (int(ccvb.Code[ccvb.PC+1]) * 256) + int(ccvb.Code[ccvb.PC+2]) // next 2 bytes point to CP entry
	if CPslot < 1 || CPslot >= len(ccvb.CP.CpIndex) {
		errMsg := fmt.Sprintf("%s:\n GETFIELD at %d: invalid CPslot (%d)",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPslot)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	CPentry := ccvb.CP.CpIndex[CPslot]
	if CPentry.Type != FieldRef {
		errMsg := fmt.Sprintf("%s:\n GETFIELD at %d: CP entry (%d) is not a field reference",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	return 3
}

// GETSTATIC 0xB2 Get static field from object and push it onto the stack
func checkGetstatic(ccvb *ccvBlock_type) int {
	// check that the index points to a field reference in the CP
	CPslot := (int(ccvb.Code[ccvb.PC+1]) * 256) + int(ccvb.Code[ccvb.PC+2]) // next 2 bytes point to CP entry
	if CPslot < 1 || CPslot >= len(ccvb.CP.CpIndex) {
		errMsg := fmt.Sprintf("%s:\n GETSTATIC at %d: invalid CPslot (%d)",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPslot)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	CPentry := ccvb.CP.CpIndex[CPslot]
	if CPentry.Type != FieldRef {
		errMsg := fmt.Sprintf(
			"%s:\n GETSTATIC at %d: CP entry (%d) is not a field reference",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	return 3
}

// GOTO 0xA7
func checkGoto(ccvb *ccvBlock_type) int {
	jumpTo := int(int16(ccvb.Code[ccvb.PC+1])*256 + int16(ccvb.Code[ccvb.PC+2]))
	if ccvb.PC+jumpTo < 0 || ccvb.PC+jumpTo >= len(ccvb.Code) {
		errMsg := fmt.Sprintf("%s:\n %s at %d: illegal jump to %d",
			excNames.JVMexceptionNames[excNames.VerifyError], BytecodeNames[ccvb.Code[ccvb.PC]], ccvb.PC, ccvb.PC+jumpTo)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	return 3
}

// GOTO_W 0xC8
func checkGotow(ccvb *ccvBlock_type) int {
	jumpTo := int(types.FourBytesToInt64(ccvb.Code[ccvb.PC+1], ccvb.Code[ccvb.PC+2], ccvb.Code[ccvb.PC+3], ccvb.Code[ccvb.PC+4]))
	if ccvb.PC+jumpTo < 0 || ccvb.PC+jumpTo >= len(ccvb.Code) {
		errMsg := fmt.Sprintf("%s:\n GOTO_W at %d: illegal jump to %d",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, ccvb.PC+jumpTo)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	return 5
}

// ICONST* and LCONST Push an int or long onto op stack
func pushInt(ccvb *ccvBlock_type) int {
	ccvb.StackEntries += 1
	return 1
}

// IF_ACMPEQ 0xA5 (and the many other IF* bytecodes)
func checkIf(ccvb *ccvBlock_type) int { // most IF* bytecodes come here. Jump if condition is met
	jumpSize := int(int16(ccvb.Code[ccvb.PC+1])*256 + int16(ccvb.Code[ccvb.PC+2]))
	if ccvb.PC+jumpSize < 0 || ccvb.PC+jumpSize >= len(ccvb.Code) {
		errMsg := fmt.Sprintf("%s:\n %s at %d: illegal jump to %d",
			excNames.JVMexceptionNames[excNames.VerifyError], BytecodeNames[ccvb.Code[ccvb.PC]], ccvb.PC, ccvb.PC+jumpSize)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	ccvb.StackEntries -= 2
	return 3
}

// IF_ZERO
func checkIfzero(ccvb *ccvBlock_type) int { // Jump if condition w.r.t 0 is met
	jumpSize := int(int16(ccvb.Code[ccvb.PC+1])*256 + int16(ccvb.Code[ccvb.PC+2]))
	if ccvb.PC+jumpSize < 0 || ccvb.PC+jumpSize >= len(ccvb.Code) {
		errMsg := fmt.Sprintf("%s:\n %s at %d: illegal jump to %d",
			excNames.JVMexceptionNames[excNames.VerifyError], BytecodeNames[ccvb.Code[ccvb.PC]], ccvb.PC, ccvb.PC+jumpSize)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	ccvb.StackEntries -= 1
	return 3
}

// ILOAD and LLOAD
func checkIload(ccvb *ccvBlock_type) int {
	ccvb.StackEntries += 1
	if ccvb.wideInEffect {
		ccvb.wideInEffect = false
		return 3
	}
	return 2
}

// ILOAD_0
func checkIload0(ccvb *ccvBlock_type) int {
	if ccvb.MaxLocals < 1 {
		errMsg := fmt.Sprintf("%s:\n %s at %d: index 0 into locals is invalid",
			excNames.JVMexceptionNames[excNames.VerifyError], BytecodeNames[ccvb.Code[ccvb.PC]], ccvb.PC)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	ccvb.StackEntries += 1
	return 1
}

// ILOAD_1
func checkIload1(ccvb *ccvBlock_type) int {
	if ccvb.MaxLocals < 2 {
		errMsg := fmt.Sprintf("%s:\n %s at %d: index 1 into locals is invalid",
			excNames.JVMexceptionNames[excNames.VerifyError], BytecodeNames[ccvb.Code[ccvb.PC]], ccvb.PC)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	ccvb.StackEntries += 1
	return 1
}

// ILOAD_2
func checkIload2(ccvb *ccvBlock_type) int {
	if ccvb.MaxLocals < 3 {
		errMsg := fmt.Sprintf("%s:\n %s at %d: index 2 into locals is invalid",
			excNames.JVMexceptionNames[excNames.VerifyError], BytecodeNames[ccvb.Code[ccvb.PC]], ccvb.PC)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	ccvb.StackEntries += 1
	return 1
}

// ILOAD_3
func checkIload3(ccvb *ccvBlock_type) int {
	if ccvb.MaxLocals < 4 {
		errMsg := fmt.Sprintf("%s:\n %s at %d: index 3 into locals is invalid",
			excNames.JVMexceptionNames[excNames.VerifyError], BytecodeNames[ccvb.Code[ccvb.PC]], ccvb.PC)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	ccvb.StackEntries += 1
	return 1
}

// INVOKEDYNAMIC 0xBA
func checkInvokedynamic(ccvb *ccvBlock_type) int {
	CPslot := (int(ccvb.Code[ccvb.PC+1]) * 256) + int(ccvb.Code[ccvb.PC+2]) // next 2 bytes point to CP entry

	// is CP slot pointing to a valid CP entry?
	if CPslot < 1 || CPslot >= len(ccvb.CP.CpIndex) {
		errMsg := fmt.Sprintf("%s:\n INVOKEDYNAMIC at %d: invalid CP slot %d",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPslot)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	// is the CP entry of type InvokeDynamic?
	CPentry := ccvb.CP.CpIndex[CPslot]
	if CPentry.Type != InvokeDynamic {
		errMsg := fmt.Sprintf("%s:\n INVOKEDYNAMIC at %d: Invalid CP entry type %d, should be 18",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	// is the CP entry slot pointing to a valid InvokeDynamic entry?
	if CPentry.Slot >= uint16(len(ccvb.CP.InvokeDynamics)) {
		errMsg := fmt.Sprintf("%s:\n INVOKEDYNAMIC at %d: invalid InvokeDynamic slot %d",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPentry.Slot)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	// is the bootstrap method index pointed to by the InvokeDynamic entry valid?
	CPdyn := ccvb.CP.InvokeDynamics[CPentry.Slot]
	bootstrapIndex := CPdyn.BootstrapIndex
	if bootstrapIndex >= uint16(len(ccvb.CP.Bootstraps)) {
		errMsg := fmt.Sprintf("%s:\n INVOKEDYNAMIC at %d: invalid BootstrapIndex %d",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, bootstrapIndex)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	// is the bootstrap name and type index pointed to by the InvokeDynamic entry valid?
	bootstrapNatIndex := CPdyn.NameAndType
	if int(bootstrapNatIndex) >= len(ccvb.CP.CpIndex) {
		errMsg := fmt.Sprintf("%s:\n INVOKEDYNAMIC at %d: invalid NameAndType index %d",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, bootstrapNatIndex)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	bootstrapNat := ccvb.CP.CpIndex[bootstrapNatIndex]
	if bootstrapNat.Type != NameAndType {
		errMsg := fmt.Sprintf("%s:\n INVOKEDYNAMIC at %d: invalid NameAndType index %d",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, bootstrapNatIndex)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	// finally, we have the NAT entry that points to two UTF-8 string entries, DescIndex
	// and NameIndex, which are entries in CpIndex that point to UTF8 entries.
	// make sure that the indices into CP.Utf8Refs are valid
	NAT := ccvb.CP.NameAndTypes[bootstrapNat.Slot]
	name := ccvb.CP.CpIndex[NAT.NameIndex]
	desc := ccvb.CP.CpIndex[NAT.DescIndex]
	if name.Type != UTF8 || desc.Type != UTF8 {
		errMsg := fmt.Sprintf("%s:\n INVOKEDYNAMIC at %d: invalid NameAndType name or descriptor field at %d",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, bootstrapNatIndex)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	return 5
}

// INVOKEINTERFACE 0xB9
func checkInvokeinterface(ccvb *ccvBlock_type) int {
	CPslot := (int(ccvb.Code[ccvb.PC+1]) * 256) + int(ccvb.Code[ccvb.PC+2]) // next 2 bytes point to CP entry
	if CPslot < 1 || CPslot >= len(ccvb.CP.CpIndex) {
		errMsg := fmt.Sprintf("%s:\n INVOKEINTERFACE at %d: CPslot (%d) is invalid",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPslot)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	countByte := ccvb.Code[ccvb.PC+3]
	if countByte == 0 {
		errMsg := fmt.Sprintf("%s:\n INVOKEINTERFACE at %d: countByte (%d) is invalid",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, countByte)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	zeroByte := ccvb.Code[ccvb.PC+4]
	if zeroByte != 0 {
		errMsg := fmt.Sprintf("%s:\n INVOKEINTERFACE at %d: zeroByte (%d) is invalid",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, zeroByte)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	CPentry := ccvb.CP.CpIndex[CPslot]
	if CPentry.Type != Interface {
		// because this is not a ClassFormatError, we output a trace error message here
		errMsg := fmt.Sprintf("%s:\n INVOKEINTERFACE at %d: CP entry (%d) is not an interface reference",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	return 5
}

// INVOKESPECIAL 0xB7
func checkInvokespecial(ccvb *ccvBlock_type) int {
	// check that the index points to a method or interface reference in the CP
	CPslot := (int(ccvb.Code[ccvb.PC+1]) * 256) + int(ccvb.Code[ccvb.PC+2]) // next 2 bytes point to CP entry
	if CPslot < 1 || CPslot >= len(ccvb.CP.CpIndex) {
		errMsg := fmt.Sprintf("%s:\n INVOKESPECIAL at %d: invalid CPslot (%d)",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPslot)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	CPentry := ccvb.CP.CpIndex[CPslot]
	if CPentry.Type != MethodRef && CPentry.Type != Interface && CPentry.Type != CachedMeth {
		// because this is not a ClassFormatError, we output a trace error message here
		errMsg := fmt.Sprintf("%s:\n INVOKESPECIAL at %d: CP entry (%d) is not a method or interface reference",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	return 3
}

// INVOKESTATIC 0xB8
func checkInvokestatic(ccvb *ccvBlock_type) int {
	// check that the index points to a method or interface reference in the CP
	CPslot := (int(ccvb.Code[ccvb.PC+1]) * 256) + int(ccvb.Code[ccvb.PC+2]) // next 2 bytes point to CP entry
	if CPslot < 1 || CPslot >= len(ccvb.CP.CpIndex) {
		errMsg := fmt.Sprintf("%s:\n INVOKESTATIC at %d: invalid CPslot (%d)",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPslot)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	CPentry := ccvb.CP.CpIndex[CPslot]
	if CPentry.Type != MethodRef && CPentry.Type != Interface && CPentry.Type != CachedMeth {
		// because this is not a ClassFormatError, we output a trace message here
		errMsg := fmt.Sprintf("%s:\n INVOKESTATIC at %d: CP entry (%d) is not a method or interface reference",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	return 3
}

// INVOKEVIRTUAL 0xB6
func checkInvokevirtual(ccvb *ccvBlock_type) int {
	// check that the index points to a method reference in the CP
	CPslot := (int(ccvb.Code[ccvb.PC+1]) * 256) + int(ccvb.Code[ccvb.PC+2]) // next 2 bytes point to CP entry
	if CPslot < 1 || CPslot >= len(ccvb.CP.CpIndex) {
		errMsg := fmt.Sprintf("%s:\n INVOKEVIRTUAL at %d: invalid CPslot (%d)",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPslot)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	CPentry := ccvb.CP.CpIndex[CPslot]
	if CPentry.Type != MethodRef && CPentry.Type != CachedMeth {
		// because this is not a ClassFormatError, we emit a trace message here
		errMsg := fmt.Sprintf("%s:\n INVOKEVIRTUAL at %d: CP entry (%d) is not a method reference",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	// TODO: figure out how many parameters (+ the this reference) are on
	// the op stack for this function and make sure that the stack is at least that big.
	return 3
}

// ISTORE
func storeInt(ccvb *ccvBlock_type) int {
	ccvb.StackEntries -= 1
	return 1
}

func storeIntRet2(ccvb *ccvBlock_type) int {
	ccvb.StackEntries -= 1
	return 2
}

// LLOAD_0
func checkLload0(ccvb *ccvBlock_type) int {
	if ccvb.MaxLocals < 2 { // LLOAD requires two locals entries
		errMsg := fmt.Sprintf("%s:\n %s at %d: index 0 into locals is invalid",
			excNames.JVMexceptionNames[excNames.VerifyError], BytecodeNames[ccvb.Code[ccvb.PC]], ccvb.PC)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	ccvb.StackEntries += 1
	return 1
}

// LLOAD_1
func checkLload1(ccvb *ccvBlock_type) int {
	if ccvb.MaxLocals < 3 { // locals[0] + two locals for long
		errMsg := fmt.Sprintf("%s:\n %s at %d: index 1 into locals is invalid",
			excNames.JVMexceptionNames[excNames.VerifyError], BytecodeNames[ccvb.Code[ccvb.PC]], ccvb.PC)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	ccvb.StackEntries += 1
	return 1
}

// LLOAD_2
func checkLload2(ccvb *ccvBlock_type) int {
	if ccvb.MaxLocals < 4 { // locals 0, 1 + two locals for long
		errMsg := fmt.Sprintf("%s:\n %s at %d: index 2 into locals is invalid",
			excNames.JVMexceptionNames[excNames.VerifyError], BytecodeNames[ccvb.Code[ccvb.PC]], ccvb.PC)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	ccvb.StackEntries += 1
	return 1
}

// LLOAD_3
func checkLload3(ccvb *ccvBlock_type) int {
	if ccvb.MaxLocals < 5 { // locals 0, 1, 2 + two locals for long
		errMsg := fmt.Sprintf("%s:\n %s at %d: index 3 into locals is invalid",
			excNames.JVMexceptionNames[excNames.VerifyError], BytecodeNames[ccvb.Code[ccvb.PC]], ccvb.PC)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	ccvb.StackEntries += 1
	return 1
}

// LOOKUPSWITCH 0xAB
func checkLookupswitch(ccvb *ccvBlock_type) int { // need to check this
	basePC := ccvb.PC

	paddingBytes := 4 - ((ccvb.PC + 1) % 4)
	if paddingBytes == 4 {
		paddingBytes = 0
	}
	basePC += paddingBytes
	basePC += 4 // jump-size for default

	npairs := binary.BigEndian.Uint32(
		[]byte{ccvb.Code[basePC+1], ccvb.Code[basePC+2], ccvb.Code[basePC+3], ccvb.Code[basePC+4]})

	basePC += 4
	basePC += int(npairs) * 8

	return (basePC - ccvb.PC) + 1
}

// NEW 0xBB
func checkNew(ccvb *ccvBlock_type) int {
	CPslot := (int(ccvb.Code[ccvb.PC+1]) * 256) + int(ccvb.Code[ccvb.PC+2]) // next 2 bytes point to CP entry
	if CPslot < 1 || CPslot >= len(ccvb.CP.CpIndex) {
		errMsg := fmt.Sprintf("%s:\n NEW at %d: invalid CPslot (%d)",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPslot)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	CPentry := ccvb.CP.CpIndex[CPslot]
	if CPentry.Type != ClassRef && CPentry.Type != Interface {
		errMsg := fmt.Sprintf("%s:\n NEW at %d: Invalid reference to object type (%d)",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	ccvb.StackEntries += 1
	return 3
}

// NEWARRAY 0xBC Create a 1-dimensional array of a primitive type
func checkNewArray(ccvb *ccvBlock_type) int {
	atype := ccvb.Code[ccvb.PC+1]
	if atype < 4 || atype > 11 {
		errMsg := fmt.Sprintf("%s:\n NEWARRAY at %d: invalid array type code %d",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, atype)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	// pops the count and pushes the array reference: net 0
	return 2
}

// POP
func checkPop(ccvb *ccvBlock_type) int {
	ccvb.StackEntries -= 1
	return 1
}

// POP2
func checkPop2(ccvb *ccvBlock_type) int {
	ccvb.StackEntries -= 2
	return 1
}

// PUTFIELD 0xB5 Put non-static field
func checkPutfield(ccvb *ccvBlock_type) int {
	// check that the index points to a field reference in the CP
	CPslot := (int(ccvb.Code[ccvb.PC+1]) * 256) + int(ccvb.Code[ccvb.PC+2]) // next 2 bytes point to CP entry
	if CPslot < 1 || CPslot >= len(ccvb.CP.CpIndex) {
		errMsg := fmt.Sprintf("%s:\n PUTFIELD at %d: CPslot (%d) invalid",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPslot)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	CPentry := ccvb.CP.CpIndex[CPslot]
	if CPentry.Type != FieldRef {
		errMsg := fmt.Sprintf("%s:\n PUTFIELD at %d: CP entry (%d) is not a field reference",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	return 3
}

// PUTSTATIC 0xB3 Put static field
func checkPutstatic(ccvb *ccvBlock_type) int {
	// check that the index points to a field reference in the CP
	CPslot := (int(ccvb.Code[ccvb.PC+1]) * 256) + int(ccvb.Code[ccvb.PC+2]) // next 2 bytes point to CP entry
	if CPslot < 1 || CPslot >= len(ccvb.CP.CpIndex) {
		errMsg := fmt.Sprintf("%s:\n PUTSTATIC at %d: CPslot (%d) invalid",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPslot)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	CPentry := ccvb.CP.CpIndex[CPslot]
	if CPentry.Type != FieldRef {
		errMsg := fmt.Sprintf("%s:\n PUTSTATIC at %d: CP entry (%d) is not a field reference",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}
	return 3
}

// TABLESWITCH 0xAA
func checkTableSwitch(ccvb *ccvBlock_type) int {
	basePC := ccvb.PC

	paddingBytes := 4 - ((ccvb.PC + 1) % 4)
	if paddingBytes == 4 {
		paddingBytes = 0
	}
	basePC += paddingBytes
	basePC += 4 // jump-size for default
	low := types.FourBytesToInt64(ccvb.Code[basePC+1], ccvb.Code[basePC+2], ccvb.Code[basePC+3], ccvb.Code[basePC+4])
	high := types.FourBytesToInt64(ccvb.Code[basePC+5], ccvb.Code[basePC+6], ccvb.Code[basePC+7], ccvb.Code[basePC+8])
	basePC += 8 // 4 bytes for low, 4 bytes for high

	if !(low <= high) {
		errMsg := fmt.Sprintf("%s:\n TABLESWITCH at %d: low (%d) exceeds high (%d)",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, low, high)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	offsetsCount := high - low + 1
	basePC += int(offsetsCount) * 4
	return (basePC - ccvb.PC) + 1
}

// MULTIANEWARRAY 0xC5 (create a multidimensional array)
func checkMultianewarray(ccvb *ccvBlock_type) int {
	CPslot := (int(ccvb.Code[ccvb.PC+1]) * 256) + int(ccvb.Code[ccvb.PC+2]) // next 2 bytes point to CP entry
	if CPslot < 1 || CPslot >= len(ccvb.CP.CpIndex) {
		errMsg := fmt.Sprintf("%s:\n MULTIANEWARRAY at %d: CPslot (%d) is invalid",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPslot)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	CPentry := ccvb.CP.CpIndex[CPslot]
	if CPentry.Type != ClassRef {
		// because this is not a ClassFormatError, we emit a trace message here
		errMsg := fmt.Sprintf("%s:\n MULTIANEWARRAY at %d: CP entry (%d) is not a class reference",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	dimensions := ccvb.Code[ccvb.PC+3]
	if dimensions == 0 {
		errMsg := fmt.Sprintf("%s:\n MULTIANEWARRAY at %d: dimensions (%d) is invalid",
			excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, dimensions)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	// we pop off as many ints from the stack as there are dimensions
	ccvb.StackEntries -= int(dimensions)

	return 4
}

// SIPUSH 0x11 create int from next 2 bytes and push it
func checkSipush(ccvb *ccvBlock_type) int {
	if len(ccvb.Code) > ccvb.PC+2 {
		ccvb.StackEntries += 1
		return 3
	}
	errMsg := fmt.Sprintf("%s:\n SIPUSH at %d: len(ccvb.Code) (%d) is invalid",
		excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, len(ccvb.Code))
	trace.Error(errMsg)
	return ERROR_OCCURRED
}

// WIDE
func checkWide(ccvb *ccvBlock_type) int {
	ccvb.wideInEffect = true
	return 1
}

// === utility functions ===

// a one-byte opcode that has nothing that can be checked
func return1(*ccvBlock_type) int {
	return 1
}

func return2(*ccvBlock_type) int {
	return 2
}

func return3(*ccvBlock_type) int {
	return 3
}

func return5(*ccvBlock_type) int {
	return 5
}

func pushIntRet2(ccvb *ccvBlock_type) int {
	ccvb.StackEntries += 1
	return 2
}

func pushIntRet3(ccvb *ccvBlock_type) int {
	ccvb.StackEntries += 1
	return 3
}

func bytecodeIsForLongOrDouble(bytecode byte) bool {
	switch bytecode {
	case 0x09, 0x0A, 0x0E, 0x0F, 0x14, 0x16, 0x18, 0x1E,
		0x1F, 0x20, 0x21, 0x26, 0x27, 0x28, 0x29, 0x2F,
		0x31, 0x37, 0x39, 0x3F, 0x40, 0x41, 0x42, 0x47,
		0x48, 0x49, 0x4A, 0x50, 0x52, 0x61, 0x63, 0x65,
		0x67, 0x69, 0x6B, 0x6D, 0x6F, 0x71, 0x73, 0x75,
		0x77, 0x79, 0x7B, 0x7D, 0x7F, 0x81, 0x83, 0x85,
		0x87, 0x88, 0x89, 0x8A, 0x8C, 0x8D, 0x8E, 0x8F,
		0x90, 0x94, 0x97, 0x98, 0xAD, 0xAF:
		// handle long/double bytecodes
		return true
	default:
		return false
	}
}

func bytecodePushes32BitValue(bytecode byte) bool {
	switch bytecode {
	case 0x02, 0x03, 0x04, 0x05, // ICONST_M1, ICONST_0, ICONST_1, ICONST_2
		0x06, 0x07, 0x08, 0x0B, // ICONST_3, ICONST_4, ICONST_5, FCONST_0
		0x0C, 0x0D, 0x10, 0x11, // FCONST_1, FCONST_2, BIPUSH, SIPUSH
		0x12, 0x13, 0x15, 0x17, // LDC, LDC_W, ILOAD, FLOAD
		0x1A, 0x1B, 0x1C, 0x1D, // ILOAD_0, ILOAD_1, ILOAD_2, ILOAD_3
		0x22, 0x23, 0x24, 0x25, // FLOAD_0, FLOAD_1, FLOAD_2, FLOAD_3
		0x2E, 0x30, 0x33, 0x34, // IALOAD, FALOAD, BALOAD, CALOAD
		0x35, 0x60, 0x62, 0x64, // SALOAD, IADD, FADD, ISUB
		0x66, 0x68, 0x6A, 0x6C, // FSUB, IMUL, FMUL, IDIV
		0x6E, 0x70, 0x72, 0x74, // FDIV, IREM, FREM, INEG
		0x76, 0x78, 0x7A, 0x7C, // FNEG, ISHL, ISHR, IUSHR
		0x7E, 0x80, 0x82, 0x86, // IAND, IOR, IXOR, I2F
		0x88, 0x89, 0x8B, 0x8E, // L2I, L2F, F2I, D2I
		0x90, 0x91, 0x92, 0x93, // D2F, I2B, I2C, I2S
		0x95, 0x96: // FCMPL, FCMPG
		return true
	default:
		return false
	}
}

// bytecode instructions 204-255 (0xCB-0xFF) are disallowed in class files
func invalidBytecode(ccvb *ccvBlock_type) int {
	errMsg := fmt.Sprintf("%s:\n Invalid bytecode instruction at %d: %X",
		excNames.JVMexceptionNames[excNames.VerifyError], ccvb.PC, ccvb.Code[ccvb.PC])
	trace.Error(errMsg)
	return ERROR_OCCURRED
}
