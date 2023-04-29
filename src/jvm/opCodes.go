/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

const AALOAD = 0x32
const AASTORE = 0x53
const ACONST_NULL = 0x01
const ALOAD = 0x19
const ALOAD_0 = 0x2A
const ALOAD_1 = 0x2B
const ALOAD_2 = 0x2C
const ALOAD_3 = 0x2D
const ANEWARRAY = 0xBD
const ARETURN = 0xB0
const ARRAYLENGTH = 0xBE
const ASTORE = 0x3A
const ASTORE_0 = 0x4B
const ASTORE_1 = 0x4C
const ASTORE_2 = 0x4D
const ASTORE_3 = 0x4E
const ATHROW = 0xBF
const BALOAD = 0x33
const BASTORE = 0x54
const BIPUSH = 0x10
const BREAKPOINT = 0xCA
const CALOAD = 0x34
const CASTORE = 0x55
const CHECKCAST = 0xC0
const D2F = 0x90
const D2I = 0x8E
const D2L = 0x8F
const DADD = 0x63
const DALOAD = 0x31
const DASTORE = 0x52
const DCMPG = 0x98
const DCMPL = 0x97
const DCONST_0 = 0x0E
const DCONST_1 = 0x0F
const DDIV = 0x6F
const DLOAD = 0x18
const DLOAD_0 = 0x26
const DLOAD_1 = 0x27
const DLOAD_2 = 0x28
const DLOAD_3 = 0x29
const DMUL = 0x6B
const DNEG = 0x77
const DREM = 0x73
const DRETURN = 0xAF
const DSTORE = 0x39
const DSTORE_0 = 0x47
const DSTORE_1 = 0x48
const DSTORE_2 = 0x49
const DSTORE_3 = 0x4A
const DSUB = 0x67
const DUP = 0x59
const DUP_X1 = 0x5A
const DUP_X2 = 0x5B
const DUP2 = 0x5C
const DUP2_X1 = 0x5D
const DUP2_X2 = 0x5E
const F2D = 0x8D
const F2I = 0x8B
const F2L = 0x8C
const FADD = 0x62
const FALOAD = 0x30
const FASTORE = 0x51
const FCMPG = 0x96
const FCMPL = 0x95
const FCONST_0 = 0x0B
const FCONST_1 = 0x0C
const FCONST_2 = 0x0D
const FDIV = 0x6E
const FLOAD = 0x17
const FLOAD_0 = 0x22
const FLOAD_1 = 0x23
const FLOAD_2 = 0x24
const FLOAD_3 = 0x25
const FMUL = 0x6A
const FNEG = 0x76
const FREM = 0x72
const FRETURN = 0xAE
const FSTORE = 0x38
const FSTORE_0 = 0x43
const FSTORE_1 = 0x44
const FSTORE_2 = 0x45
const FSTORE_3 = 0x46
const FSUB = 0x66
const GETFIELD = 0xB4
const GETSTATIC = 0xB2
const GOTO = 0xA7
const GOTO_W = 0xC8
const I2B = 0x91
const I2C = 0x92
const I2D = 0x87
const I2F = 0x86
const I2L = 0x85
const I2S = 0x93
const IADD = 0x60
const IALOAD = 0x2E
const IAND = 0x7E
const IASTORE = 0x4F
const ICONST_N1 = 0x02
const ICONST_0 = 0x03
const ICONST_1 = 0x04
const ICONST_2 = 0x05
const ICONST_3 = 0x06
const ICONST_4 = 0x07
const ICONST_5 = 0x08
const IDIV = 0x6C
const IF_ACMPEQ = 0xA5
const IF_ACMPNE = 0xA6
const IF_ICMPEQ = 0x9F
const IF_ICMPGE = 0xA2
const IF_ICMPGT = 0xA3
const IF_ICMPLE = 0xA4
const IF_ICMPLT = 0xA1
const IF_ICMPNE = 0xA0
const IFEQ = 0x99
const IFGE = 0x9C
const IFGT = 0x9D
const IFLE = 0x9E
const IFLT = 0x9B
const IFNE = 0x9A
const IFNONNULL = 0xC7
const IFNULL = 0xC6
const IINC = 0x84
const ILOAD = 0x15
const ILOAD_0 = 0x1A
const ILOAD_1 = 0x1B
const ILOAD_2 = 0x1C
const ILOAD_3 = 0x1D
const IMPDEP1 = 0xFE
const IMPDEP2 = 0xFF
const IMUL = 0x68
const INEG = 0x74
const INSTANCEOF = 0xC1
const INVOKEDYNAMIC = 0xBA
const INVOKEINTERFACE = 0xB9
const INVOKESPECIAL = 0xB7
const INVOKESTATIC = 0xB8
const INVOKEVIRTUAL = 0xB6
const IOR = 0x80
const IREM = 0x70
const IRETURN = 0xAC
const ISHL = 0x78
const ISHR = 0x7A
const ISTORE = 0x36
const ISTORE_0 = 0x3B
const ISTORE_1 = 0x3C
const ISTORE_2 = 0x3D
const ISTORE_3 = 0x3E
const ISUB = 0x64
const IUSHR = 0x7C
const IXOR = 0x82
const JSR = 0xA8
const JSR_W = 0xC9
const L2D = 0x8A
const L2F = 0x89
const L2I = 0x88
const LADD = 0x61
const LALOAD = 0x2F
const LAND = 0x7F
const LASTORE = 0x50
const LCMP = 0x94
const LCONST_0 = 0x09
const LCONST_1 = 0x0A
const LDC = 0x12
const LDC_W = 0x13
const LDC2_W = 0x14
const LDIV = 0x6D
const LLOAD = 0x16
const LLOAD_0 = 0x1E
const LLOAD_1 = 0x1F
const LLOAD_2 = 0x20
const LLOAD_3 = 0x21
const LMUL = 0x69
const LNEG = 0x75
const LOOKUPSWITCH = 0xAB
const LOR = 0x81
const LREM = 0x71
const LRETURN = 0xAD
const LSHL = 0x79
const LSHR = 0x7B
const LSTORE = 0x37
const LSTORE_0 = 0x3F
const LSTORE_1 = 0x40
const LSTORE_2 = 0x41
const LSTORE_3 = 0x42
const LSUB = 0x65
const LUSHR = 0x7D
const LXOR = 0x83
const MONITORENTER = 0xC2
const MONITOREXIT = 0xC3
const MULTIANEWARRAY = 0xC5
const NEW = 0xBB
const NEWARRAY = 0xBC
const NOP = 0x00
const POP = 0x57
const POP2 = 0x58
const PUTFIELD = 0xB5
const PUTSTATIC = 0xB3
const RET = 0xA9
const RETURN = 0xB1
const SALOAD = 0x35
const SASTORE = 0x56
const SIPUSH = 0x11
const SWAP = 0x5F
const TABLESWITCH = 0xAA
const WIDE = 0xC4

var BytecodeNames = []string{
    "NOP",             // 0x00
    "ACONST_NULL",     // 0x01
    "ICONST_N1",       // 0x02
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
    "FCONST_2",        // OxOD
    "DCONST_0",        // 0x0E
    "DCONST_1",        // 0x0F
    "BIPUSH",          // 0X10
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
