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
	"strings"
)

// Here we check the bytecodes of a method. The method is passed as a byte slice.
// The bytecodes are checked against a function lookup table which calls the function
// to performs the check. It then uses the skip table to determine the number of bytes
// to skip to the next bytecode. If an error occurs, a ClassFormatException is thrown.

// NOTE: The unit tests for these functions are in codeCheck_test.go in the jvm directory.
// Placed there to avoid circular dependencies.

type BytecodeFunc func() int

var ERROR_OCCURRED = math.MaxInt32
var WideInEffect = false

var CheckTable = [203]BytecodeFunc{
	return1,              // NOP             0x00
	checkAconstnull,      // ACONST_NULL     0x01
	checkIconst,          // ICONST_M1       0x02
	checkIconst,          // ICONST_0        0x03
	checkIconst,          // ICONST_1        0x04
	checkIconst,          // ICONST_2        0x05
	checkIconst,          // ICONST_3        0x06
	checkIconst,          // ICONST_4        0x07
	checkIconst,          // ICONST_5        0x08
	checkIconst,          // LCONST_0        0x09
	checkIconst,          // LCONST_1        0x0A
	return1,              // FCONST_0        0x0B
	return1,              // FCONST_1        0x0C
	return1,              // FCONST_2        0x0D
	return1,              // DCONST_0        0x0E
	return1,              // DCONST_1        0x0F
	checkBipush,          // BIPUSH          0x10
	return3,              // SIPUSH          0x11
	checkLdc,             // LDC             0x12
	return3,              // LDC_W           0x13
	return3,              // LDC2_W          0x14
	return2,              // ILOAD           0x15
	return2,              // LLOAD           0x16
	return2,              // FLOAD           0x17
	return2,              // DLOAD           0x18
	checkAload,           // ALOAD           0x19
	checkIload0,          // ILOAD_0         0x1A
	checkIload1,          // ILOAD_1         0x1B
	return1,              // ILOAD_2         0x1C
	return1,              // ILOAD_3         0x1D
	return1,              // LLOAD_0         0x1E
	return1,              // LLOAD_1         0x1F
	return1,              // LLOAD_2         0x20
	return1,              // LLOAD_3         0x21
	return1,              // FLOAD_0         0x22
	return1,              // FLOAD_1         0x23
	return1,              // FLOAD_2         0x24
	return1,              // FLOAD_3         0x25
	return1,              // DLOAD_0         0x26
	return1,              // DLOAD_1         0x27
	return1,              // DLOAD_2         0x28
	return1,              // DLOAD_3         0x29
	return1,              // ALOAD_0         0x2A
	return1,              // ALOAD_1         0x2B
	return1,              // ALOAD_2         0x2C
	return1,              // ALOAD_3         0x2D
	return1,              // IALOAD          0x2E
	return1,              // LALOAD          0x2F
	return1,              // FALOAD          0x30
	return1,              // DALOAD          0x31
	return1,              // AALOAD          0x32
	return1,              // BALOAD          0x33
	return1,              // CALOAD          0x34
	return1,              // SALOAD          0x35
	return2,              // ISTORE          0x36
	return2,              // LSTORE          0x37
	return2,              // FSTORE          0x38
	return2,              // DSTORE          0x39
	return2,              // ASTORE          0x3A
	checkIstore0,         // ISTORE_0        0x3B
	checkIstore1,         // ISTORE_1        0x3C
	return1,              // ISTORE_2        0x3D
	return1,              // ISTORE_3        0x3E
	return1,              // LSTORE_0        0x3F
	return1,              // LSTORE_1        0x40
	return1,              // LSTORE_2        0x41
	return1,              // LSTORE_3        0x42
	return1,              // FSTORE_0        0x43
	return1,              // FSTORE_1        0x44
	return1,              // FSTORE_2        0x45
	return1,              // FSTORE_3        0x46
	return1,              // DSTORE_0        0x47
	return1,              // DSTORE_1        0x48
	return1,              // DSTORE_2        0x49
	return1,              // DSTORE_3        0x4A
	return1,              // ASTORE_0        0x4B
	return1,              // ASTORE_1        0x4C
	return1,              // ASTORE_2        0x4D
	return1,              // ASTORE_3        0x4E
	return1,              // IASTORE         0x4F
	return1,              // LASTORE         0x50
	return1,              // FASTORE         0x51
	return1,              // DASTORE         0x52
	return1,              // AASTORE         0x53
	return1,              // BASTORE         0x54
	return1,              // CASTORE         0x55
	return1,              // SASTORE         0x56
	checkPop,             // POP             0x57
	checkPop2,            // POP2            0x58
	checkDup,             // DUP             0x59
	checkDupx1,           // DUP_X1          0x5A
	checkDupx2,           // DUP_X2          0x5B
	checkDup2,            // DUP2            0x5C
	checkDup2x1,          // DUP2_X1         0x5D
	checkDup2x2,          // DUP2_X2         0x5E
	checkSwap,            // SWAP            0x5F
	return1,              // IADD            0x60
	return1,              // LADD            0x61
	return1,              // FADD            0x62
	return1,              // DADD            0x63
	return1,              // ISUB            0x64
	return1,              // LSUB            0x65
	return1,              // FSUB            0x66
	return1,              // DSUB            0x67
	return1,              // IMUL            0x68
	return1,              // LMUL            0x69
	return1,              // FMUL            0x6A
	return1,              // DMUL            0x6B
	return1,              // IDIV            0x6C
	return1,              // LDIV            0x6D
	return1,              // FDIV            0x6E
	return1,              // DDIV            0x6F
	return1,              // IREM            0x70
	return1,              // LREM            0x71
	return1,              // FREM            0x72
	return1,              // DREM            0x73
	return1,              // INEG            0x74
	return1,              // LNEG            0x75
	return1,              // FNEG            0x76
	return1,              // DNEG            0x77
	return1,              // ISHL            0x78
	return1,              // LSHL            0x79
	return1,              // ISHR            0x7A
	return1,              // LSHR            0x7B
	return1,              // IUSHR           0x7C
	return1,              // LUSHR           0x7D
	return1,              // IAND            0x7E
	return1,              // LAND            0x7F
	return1,              // IOR             0x80
	return1,              // LOR             0x81
	return1,              // IXOR            0x82
	return1,              // LXOR            0x83
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
	return1,              // LCMP            0x94
	return1,              // FCMPL           0x95
	return1,              // FCMPG           0x96
	return1,              // DCMPL           0x97
	return1,              // DCMPG           0x98
	checkIfwithint,       // IFEQ            0x99
	checkIfwithint,       // IFNE            0x9A
	checkIfwithint,       // IFLT            0x9B
	checkIfwithint,       // IFGE            0x9C
	checkIfwithint,       // IFGT            0x9D
	checkIfwithint,       // IFLE            0x9E
	checkIfwith2ints,     // IF_ICMPEQ       0x9F
	checkIfwith2ints,     // IF_ICMPNE       0xA0
	checkIfwith2ints,     // IF_ICMPLT       0xA1
	checkIfwith2ints,     // IF_ICMPGE       0xA2
	checkIfwith2ints,     // IF_ICMPGT       0xA3
	checkIfwith2ints,     // IF_ICMPLE       0xA4
	checkIfwith2refs,     // IF_ACMPEQ       0xA5
	checkIfwith2refs,     // IF_ACMPNE       0xA6
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
	checkReturn,          // RETURN          0xB1
	checkGetstatic,       // GETSTATIC       0xB2
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
var CP *CPool
var Code []byte
var OpStack []byte // values are: N = nil, I = int, L = long, F = float, R = reference, U = unknown
var TOS int        // index to top of stack, 0 = empty (note: stack is 0-based, not -1 based as in the interpreter)
var LocalsCount int
var Locals []byte // uses samve values as OpStack

func CheckCodeValidity(code []byte, cp *CPool, stackSize int, locals int) error {
	// check that the code is valid
	if code == nil || cp == nil {
		errMsg := "CheckCodeValidity: nil code or constant pool"
		return errors.New(errMsg)
	}

	CP = cp
	if len(CP.CpIndex) == 0 {
		errMsg := "CheckCodeValidity: empty constant pool"
		return errors.New(errMsg)
	}

	// set up the simulated operand stack
	OpStack = make([]byte, stackSize+1) // +1 for 0-based stack
	for i := 0; i < stackSize+1; i++ {
		OpStack[i] = 'N'
	}
	TOS = -1

	// set up the simulated local variables
	LocalsCount = locals
	Locals = make([]byte, LocalsCount)
	for i := 0; i < locals; i++ {
		Locals[i] = 'N'
	}

	Code = code
	PC = 0
	for PC < len(code) {
		opcode := code[PC]
		ret := CheckTable[opcode]()
		if ret == ERROR_OCCURRED {
			errMsg := fmt.Sprintf("Invalid bytecode or argument at location %d", PC)
			return errors.New(errMsg)
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

// ACONST_NULL 0x01 // Push null onto the stack
func checkAconstnull() int {
	TOS += 1
	if TOS > len(OpStack) {
		return ERROR_OCCURRED
	}
	OpStack[TOS] = 'R'
	return 1
}

// ALOAD 0x19 Load reference onto the stack from local variable specified by the following index byte
func checkAload() int {
	index := int(Code[PC+1])
	if index >= LocalsCount {
		return ERROR_OCCURRED
	}

	if Locals[index] != 'R' && Locals[index] != 'U' {
		return ERROR_OCCURRED
	}

	TOS += 1
	if TOS > len(OpStack) {
		return ERROR_OCCURRED
	}

	OpStack[TOS] = 'R'
	return 2
}

// BIPUSH 0x10 Push byte onto the stack
func checkBipush() int {
	TOS += 1
	if TOS > len(OpStack) {
		return ERROR_OCCURRED
	}

	OpStack[TOS] = 'I'
	return 2
}

// DUP 0x59 Duplicate the top value on the stack
func checkDup() int {
	if TOS < 1 {
		return ERROR_OCCURRED
	}

	TOS += 1
	if TOS > len(OpStack) {
		return ERROR_OCCURRED
	}

	OpStack[TOS] = OpStack[TOS-1]
	return 1
}

// DUP_X1 0x5A Duplicate the top value on the stack and insert two down
func checkDupx1() int {
	if TOS < 2 {
		return ERROR_OCCURRED
	}

	TOS += 1
	if TOS > len(OpStack) {
		return ERROR_OCCURRED
	}

	initialTOSvalue := OpStack[TOS]

	OpStack[TOS] = OpStack[TOS-1]
	OpStack[TOS-1] = OpStack[TOS-2]
	OpStack[TOS-2] = initialTOSvalue
	return 1
}

// DUP_X2 0x5B Duplicate the top value on the stack and insert three down
func checkDupx2() int {
	if TOS < 3 {
		return ERROR_OCCURRED
	}

	TOS += 1
	if TOS > len(OpStack) {
		return ERROR_OCCURRED
	}

	initialTOSvalue := OpStack[TOS]

	OpStack[TOS] = OpStack[TOS-1]
	OpStack[TOS-1] = OpStack[TOS-2]
	OpStack[TOS-2] = OpStack[TOS-3]
	OpStack[TOS-3] = initialTOSvalue
	return 1
}

// DUP2 0x5C Duplicate the top two values on the stack
func checkDup2() int {
	if TOS < 2 {
		return ERROR_OCCURRED
	}

	TOS += 2
	if TOS > len(OpStack) {
		return ERROR_OCCURRED
	}

	OpStack[TOS] = OpStack[TOS-2]
	OpStack[TOS-1] = OpStack[TOS-3]
	return 1
}

// DUP2_X1 0x5D Duplicate the top two values on the stack and insert them two down. So,
// ..., value3, value2, value1 <-TOS
//
//	becomes:
//
// ..., value2, value1, value3, value2, value1
func checkDup2x1() int {
	if TOS < 3 {
		return ERROR_OCCURRED
	}

	initialTOSvalue := OpStack[TOS]
	initialTOSplus1value := OpStack[TOS-1]

	TOS += 2
	if TOS > len(OpStack) {
		return ERROR_OCCURRED
	}

	OpStack[TOS] = OpStack[TOS-2]
	OpStack[TOS-1] = OpStack[TOS-3]
	OpStack[TOS-2] = OpStack[TOS-4]
	OpStack[TOS-3] = initialTOSvalue
	OpStack[TOS-4] = initialTOSplus1value
	return 1
}

// DUP2_X2 0x5E Duplicate the top two values on the stack and insert them three down. So,
// ..., value4, value3, value2, value1 <-TOS
//
//	becomes:
//
// ..., value2, value1, value4, value3, value2, value1
func checkDup2x2() int {
	if TOS < 4 {
		return ERROR_OCCURRED
	}

	initialTOSvalue := OpStack[TOS]
	initialTOSplus1value := OpStack[TOS-1]

	TOS += 2
	if TOS > len(OpStack) {
		return ERROR_OCCURRED
	}

	OpStack[TOS] = OpStack[TOS-2]
	OpStack[TOS-1] = OpStack[TOS-3]
	OpStack[TOS-2] = OpStack[TOS-4]
	OpStack[TOS-3] = OpStack[TOS-5]
	OpStack[TOS-4] = initialTOSvalue
	OpStack[TOS-5] = initialTOSplus1value
	return 1
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

	if TOS+1 > len(OpStack) {
		return ERROR_OCCURRED
	}

	TOS += 1
	OpStack[TOS] = 'U' // unknown type, as we don't know the type of the field.
	return 3
}

// GETSTATIC 0xB2 Get static field and push it onto the stack
func checkGetstatic() int {
	// check that the index points to a field reference in the CP
	CPslot := (int(Code[PC+1]) * 256) + int(Code[PC+2]) // next 2 bytes point to CP entry
	if CPslot < 1 || CPslot >= len(CP.CpIndex) {
		return ERROR_OCCURRED
	}

	CPentry := CP.CpIndex[CPslot]
	if CPentry.Type != FieldRef {
		errMsg := fmt.Sprintf("%s:\n GETSTATIC at %d: CP entry (%d) is not a field reference",
			excNames.JVMexceptionNames[excNames.VerifyError], PC, CPentry.Type)
		trace.Error(errMsg)
		return ERROR_OCCURRED
	}

	if TOS+1 > len(OpStack) {
		return ERROR_OCCURRED
	}

	TOS += 1
	OpStack[TOS] = 'U' // unknown type, as we don't know the type of the field.
	return 3
}

// GOTO 0xA7
func checkGoto() int {
	jumpTo := int(int16(Code[PC+1])*256 + int16(Code[PC+2]))
	if PC+jumpTo < 0 || PC+jumpTo >= len(Code) {
		return ERROR_OCCURRED
	}

	// TODO handle saving state for jump. Don't jump backwards.
	return 3
}

// GOTO_W 0xC8
func checkGotow() int {
	jumpTo := int(types.FourBytesToInt64(Code[PC+1], Code[PC+2], Code[PC+3], Code[PC+4]))
	if PC+jumpTo < 0 || PC+jumpTo >= len(Code) {
		return ERROR_OCCURRED
	}
	return 5
}

// ICONST_M1 0x02 Push int constant -1 onto the stack
// ICONST_0 0x03 Push int constant 0 onto the stack
// ICONST_1 0x04 Push int constant 1 onto the stack
// ICONST_2 0x05 Push int constant 2 onto the stack
// ICONST_3 0x06 Push int constant 3 onto the stack
// ICONST_4 0x07 Push int constant 4 onto the stack
// ICONST_5 0x08 Push int constant 5 onto the stack
// LCONST_0 0x09 Push long constant 0 onto the stack
// LCONST_1 0x0A Push long constant 1 onto the stack
func checkIconst() int {
	TOS += 1
	if TOS > len(OpStack) {
		return ERROR_OCCURRED
	}

	OpStack[TOS] = 'I'
	return 1
}

// IF_ACMPEQ 0xA5 Pop two references off the stack and jump if they are equal
// IF_ACMPNE 0xA6 Pop two references off the stack and jump if they are not equal
func checkIfwith2refs() int {
	jumpSize := int(int16(Code[PC+1])*256 + int16(Code[PC+2]))
	if PC+jumpSize < 0 || PC+jumpSize >= len(Code) {
		return ERROR_OCCURRED
	}

	if TOS < 2 {
		return ERROR_OCCURRED
	} else {
		if (OpStack[TOS] != 'R' && OpStack[TOS] != 'U') ||
			(OpStack[TOS-1] != 'R' && OpStack[TOS-1] != 'U') {
			return ERROR_OCCURRED
		}
	}

	TOS -= 2
	return 3
}

// IF_ICMPEQ       0x9F pop two ints off the stack and jump if comparison succeeds
// IF_ICMPNE       0xA0
// IF_ICMPLT       0xA1
// IF_ICMPGE       0xA2
// IF_ICMPGT       0xA3
// IF_ICMPLE       0xA4
func checkIfwith2ints() int {
	jumpSize := int(int16(Code[PC+1])*256 + int16(Code[PC+2]))
	if PC+jumpSize < 0 || PC+jumpSize >= len(Code) {
		return ERROR_OCCURRED
	}

	if TOS < 2 {
		return ERROR_OCCURRED
	} else {
		if (OpStack[TOS] != 'I' && OpStack[TOS] != 'U') ||
			(OpStack[TOS-1] != 'I' && OpStack[TOS-1] != 'U') {
			return ERROR_OCCURRED
		}
	}

	TOS -= 2
	return 3
}

// IFEQ 0x99 pop int off the stack and jump if comparison with zero succeeds
// IFNE 0x9A
// IFLT 0x9B
// IFGE 0x9C
// IFGT 0x9D
// IFLE 0x9E
func checkIfwithint() int { //
	jumpSize := int(int16(Code[PC+1])*256 + int16(Code[PC+2]))
	if PC+jumpSize < 0 || PC+jumpSize >= len(Code) {
		return ERROR_OCCURRED
	}

	if TOS < 1 {
		return ERROR_OCCURRED
	} else {
		if OpStack[TOS] == 'F' || OpStack[TOS] == 'R' {
			return ERROR_OCCURRED
		}
	}
	TOS -= 1
	return 3
}

// // IINC 0x84 Increment local variable by constant
// func checkIinc() int {
// 	index := int(Code[PC+1])
// 	if index >= LocalsCount {
// 		return ERROR_OCCURRED
// 	}
//
// 	if Locals[index] != 'F' || Locals[index] != 'R' {
// 		return ERROR_OCCURRED
// 	}
//
// 	if Locals[index] != 'U' {
// 		Locals[index] = 'I'
// 	}
//
// 	return 3
// }

// ILOAD_0 0x1A Load int from local variable 0
func checkIload0() int {
	if LocalsCount < 1 {
		return ERROR_OCCURRED
	}

	if Locals[0] == 'F' || Locals[0] != 'R' {
		return ERROR_OCCURRED
	}

	TOS += 1
	if TOS > len(OpStack) {
		return ERROR_OCCURRED
	}

	OpStack[TOS] = 'I'
	return 1
}

// ILOAD_1 0x1B Load int from local variable 1
func checkIload1() int {
	if LocalsCount < 2 {
		return ERROR_OCCURRED
	}

	if Locals[1] == 'F' || Locals[1] != 'R' {
		return ERROR_OCCURRED
	}

	TOS += 1
	if TOS > len(OpStack) {
		return ERROR_OCCURRED
	}

	OpStack[TOS] = 'I'
	return 1
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

	_, _, methodType := GetMethInfoFromCPmethref(CP, CPslot)
	if !strings.HasSuffix(methodType, "V") { // if the return is not void
		// so, we need to push the return value onto the stack
		TOS += 1
		if TOS > len(OpStack) {
			return ERROR_OCCURRED
		}
		OpStack[TOS] = 'U' // unknown type. TODO: use the type data to determine the type of the return value
	}
	return 3
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

// ISTORE_0 0x3B Store int from stack into local variable 0
func checkIstore0() int {
	if LocalsCount < 1 {
		return ERROR_OCCURRED
	}

	if TOS < 1 {
		return ERROR_OCCURRED
	}

	Locals[0] = OpStack[TOS]
	TOS -= 1
	return 1
}

// ISTORE_1 0x3C Store int from stack into local variable 1
func checkIstore1() int {
	if LocalsCount < 2 {
		return ERROR_OCCURRED
	}

	if TOS < 1 {
		return ERROR_OCCURRED
	}

	Locals[1] = OpStack[TOS]
	TOS -= 1
	return 1
}

// LDC 0x12 Push item from constant pool onto the stack
func checkLdc() int {
	CPslot := int(Code[PC+1])
	if CPslot < 1 || CPslot >= len(CP.CpIndex) {
		return ERROR_OCCURRED
	}

	CPentry := CP.CpIndex[CPslot]
	if CPentry.Type != StringConst && CPentry.Type != IntConst && CPentry.Type != FloatConst {
		return ERROR_OCCURRED
	}

	TOS += 1
	if TOS > len(OpStack) {
		return ERROR_OCCURRED
	}

	switch CPentry.Type {
	case IntConst:
		OpStack[TOS] = 'I'
	case FloatConst:
		OpStack[TOS] = 'F'
	default:
		OpStack[TOS] = 'R'
	}

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

// POP 0x57 Pop the top value off the stack
func checkPop() int {
	if TOS < 1 {
		return ERROR_OCCURRED
	}
	TOS -= 1
	return 1
}

// POP2 0x58 Pop the top two values off the stack
func checkPop2() int {
	if TOS < 2 {
		return ERROR_OCCURRED
	}
	TOS -= 2
	return 1
}

// RETURN 0xB1 Return void from method
func checkReturn() int {
	return 1
}

// SWAP 0x5F Swap the top two values on the stack
func checkSwap() int {
	if TOS < 2 {
		return ERROR_OCCURRED
	}

	temp := OpStack[TOS]
	OpStack[TOS] = OpStack[TOS-1]
	OpStack[TOS-1] = temp
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
