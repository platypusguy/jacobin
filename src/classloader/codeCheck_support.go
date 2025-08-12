/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

// Test wrapper functions - these expose private functions for testing
// func TestArith() int                                   { return arith() }
func TestCheckAconstnull() int { return checkAconstnull() }
func TestPushFloat() int       { return pushFloat() }
func TestPushInt() int         { return pushInt() }
func TestCheckBipush() int     { return checkBipush() }
func TestCheckSipush() int     { return checkSipush() }
func TestDup1() int            { return dup1() }
func TestDup2() int            { return dup2() }
func TestCheckPop() int        { return checkPop() }
func TestCheckPop2() int       { return checkPop2() }
func TestCheckGetfield() int   { return checkGetfield() }

// func TestCheckGoto() int       { return checkGoto() }

// func TestCheckIf() int         { return checkIf() }
// func TestCheckIfZero() int     { return checkIfZero() }
// func TestCheckInvokeinterface() int { return checkInvokeinterface() }
// func TestCheckInvokevirtual() int   { return checkInvokevirtual() }
// func TestReturn1() int                                 { return return1() }
// func TestReturn2() int                                 { return return2() }
// func TestReturn3() int                                 { return return3() }
// func TestReturn4() int                                 { return return4() }
// func TestReturn5() int                                 { return return5() }
// func TestCheckTableswitch() int                        { return checkTableswitch() }
// func TestCheckMultianewarray() int                     { return checkMultianewarray() }
// func TestByteCodeIsForLongOrDouble(bytecode byte) bool { return byteCodeIsForLongOrDouble(bytecode) }
