/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) Consult jacobin.org.
 */

package jvm

import (
	"errors"
	"fmt"
	"io"
	"jacobin/classloader"
	"jacobin/frames"
	"jacobin/gfunction"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/object"
	"jacobin/opcodes"
	"jacobin/stringPool"
	"os"
	"strings"
	"testing"
)

type TestVars struct {
	objName, // the name of the object whose method we're invoking
	methName, // the name of the method, without the class name or the type
	methType, // the type of the method
	stringParam, // an optional string parameter to pass to the method under test
	stderrText, // the expected output text to stderr, if any. "" = no output expected
	stdoutText string // the expected text to stdout, if any. "" = no output expected
}

// TestPrintln() is a demonstration of how to use RunTestFuncWith0or1Strings()
func TestPrintln(t *testing.T) {
	testName := "TestPrintln"

	var params TestVars
	params.objName = "java/io/PrintStream"
	params.methName = "println"
	params.methType = "(Ljava/lang/String;)V"
	params.stringParam = "hello from table test"
	params.stderrText = ""
	params.stdoutText = "hello"
	err := RunTestFuncWith0or1Strings(params)
	if err != nil {
		t.Errorf("Error in %s: %s", testName, err.Error())
	}
}

// RunTestFuncWith0or1Strings runs the test and reports any errors via its return value.
// Unit tests should test this value to determine whether the unit test has passed/failed.
func RunTestFuncWith0or1Strings(params TestVars) error {

	var objClassName = params.objName
	var methName = params.methName
	var methType = params.methType
	var stringParam = params.stringParam
	var stderrExpected = params.stderrText
	var stdoutExpected = params.stdoutText

	// set up the environment

	globals.InitGlobals("test")
	log.Init()

	// redirect stderr and stdout
	normalStderr := os.Stderr
	rerr, werr, _ := os.Pipe()
	os.Stderr = werr

	normalStdout := os.Stdout
	rout, wout, _ := os.Pipe()
	os.Stdout = wout

	// load the classes into MTable
	classloader.MTable = make(map[string]classloader.MTentry)
	gfunction.MTableLoadGFunctions(&classloader.MTable)

	// build up the CP and the various frame fields
	CP := classloader.CPool{}
	CP.CpIndex = make([]classloader.CpEntry, 10)
	CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = classloader.CpEntry{Type: classloader.MethodRef, Slot: 0}
	CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
	CP.CpIndex[3] = classloader.CpEntry{Type: classloader.NameAndType, Slot: 0}
	CP.CpIndex[4] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
	CP.CpIndex[5] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1}
	CP.CpIndex[6] = classloader.CpEntry{Type: classloader.StringConst, Slot: 7}
	CP.CpIndex[7] = classloader.CpEntry{Type: classloader.UTF8, Slot: 2} // point to UTF8[2]

	CP.MethodRefs = make([]classloader.MethodRefEntry, 1)
	methRef := classloader.MethodRefEntry{
		ClassIndex:  2, // these are CP entries
		NameAndType: 3,
	}
	CP.MethodRefs[0] = methRef

	CP.ClassRefs = append(CP.ClassRefs, stringPool.GetStringIndex(&objClassName))

	CP.Utf8Refs = append(CP.Utf8Refs, methName)
	CP.Utf8Refs = append(CP.Utf8Refs, methType)
	CP.Utf8Refs = append(CP.Utf8Refs, stringParam)
	nAndT := classloader.NameAndTypeEntry{
		NameIndex: uint16(4),
		DescIndex: uint16(5),
	}
	CP.NameAndTypes = append(CP.NameAndTypes, nAndT)

	var f frames.Frame
	if stringParam != "" { // if there is a string param, load it with LDC before calling INVOKEVIRTUAL
		f = newFrame(opcodes.LDC)
		f.Meth = append(f.Meth, uint8(6)) // point to the string constant parameter indexed by CPindex[6]
		f.Meth = append(f.Meth, opcodes.INVOKEVIRTUAL)
	} else {
		f = newFrame(opcodes.INVOKEVIRTUAL)
	}
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to method referred to in 0x0001 of the CP

	f.CP = &CP

	// create the opStack
	for j := 0; j < 10; j++ {
		f.OpStack = append(f.OpStack, 0)
	}

	fs := frames.CreateFrameStack()

	// now push a reference to the object whose method we're calling. In the event, it's a prinstream,
	// we force it to be stdout. Otherwise, we instantiate the class.
	if objClassName == "java/io/PrintStream" { // if we're working with a printstream, force-set it to stdout
		f.OpStack[0] = os.Stdout
	} else {
		objPtr, err := InstantiateClass(objClassName, fs)
		if err != nil {
			errMsg := fmt.Sprintf("could not instantiate class object: %s  %v",
				objClassName, err)
			return errors.New(errMsg)
		} else {
			f.OpStack[0] = objPtr.(*object.Object)
		}
	}
	f.TOS = 0

	fs.PushFront(&f) // push the new frame

	// run the method (by running the frame)
	_ = runFrame(fs)

	// get contents written by stderr and stdout, then restore
	// stderr and stdout to what they were before the test
	_ = werr.Close()
	rawStderrMsg, _ := io.ReadAll(rerr)
	os.Stderr = normalStderr

	_ = wout.Close()
	rawStdoutMsg, _ := io.ReadAll(rout)
	os.Stdout = normalStdout

	// convert the contents written to stderr and stdout into strings
	// and run tests on those strings

	errMsg := string(rawStderrMsg[:])
	if len(errMsg) != 0 && !strings.Contains(errMsg, params.stderrText) {
		failureMsg := fmt.Sprintf("expected error msg: %s, got: %s",
			stderrExpected, errMsg)
		return errors.New(failureMsg)
	}

	outMsg := string(rawStdoutMsg[:])
	if len(outMsg) != 0 && !strings.Contains(outMsg, params.stdoutText) {
		failureMsg := fmt.Sprintf("expected output: %s, got: %s",
			stdoutExpected, outMsg)
		return errors.New(failureMsg)
	}
	return nil
}

// TestFuncWith0or1StringsTable() is a table-driven test that exercises the INVOKEVIRTUAL bytecode. To use it:
// fill in an instance of the testVars structure and add it to the tests map, as illustrated below.
// Several requirements:
//   - The method being invoked must be in the MTable, with the same name and type as stated in testVars;
//     if the method is a gfunction, that gfunction must be loaded in the MTable (via LoadLib).
//   - At present, you can pass only 0 or 1 strings to the invoked method.
//   - The test result must appear either on stdout or stderr (or both). Both are captured here and you can
//     specify what text either or both must contain for the test to pass.
func TestGfunWith0or1StringsTable(t *testing.T) {
	// type testVars struct {
	// 	objName, // the name of the object whose method we're invoking
	// 	methName, // the name of the method, without the class name or the type
	// 	methType, // the type of the method
	// 	stringParam, // an optional string parameter to pass to the method under test
	// 	stderrText, // the expected output text to stderr, if any. "" = no output expected
	// 	stdoutText string // the expected text to stdout, if any. "" = no output expected
	// }

	// the map holding out tests. The key is the name of the test, the value is a struct of testVars, shown next
	tests := make(map[string]TestVars)

	tv := TestVars{
		objName:     "java/io/PrintStream",
		methName:    "println",
		methType:    "(Ljava/lang/String;)V",
		stringParam: "hello from table test",
		stderrText:  "",
		stdoutText:  "hello from table test",
	}

	tests["testPrintlnWith1validString"] = tv

	for name, test := range tests {
		t.Run(name, func(t *testing.T) { // from here to the end, we define the function that is the test

			var objClassName = test.objName
			var methName = test.methName
			var methType = test.methType
			var stringParam = test.stringParam
			var stderrExpected = test.stderrText
			var stdoutExpected = test.stdoutText

			// set up the environment

			globals.InitGlobals("test")
			log.Init()

			// redirect stderr and stdout
			normalStderr := os.Stderr
			rerr, werr, _ := os.Pipe()
			os.Stderr = werr

			normalStdout := os.Stdout
			rout, wout, _ := os.Pipe()
			os.Stdout = wout

			// load the classes into MTable
			classloader.MTable = make(map[string]classloader.MTentry)
			gfunction.MTableLoadGFunctions(&classloader.MTable)

			// build up the CP and the various frame fields
			CP := classloader.CPool{}
			CP.CpIndex = make([]classloader.CpEntry, 10)
			CP.CpIndex[0] = classloader.CpEntry{Type: 0, Slot: 0}
			CP.CpIndex[1] = classloader.CpEntry{Type: classloader.MethodRef, Slot: 0}
			CP.CpIndex[2] = classloader.CpEntry{Type: classloader.ClassRef, Slot: 0}
			CP.CpIndex[3] = classloader.CpEntry{Type: classloader.NameAndType, Slot: 0}
			CP.CpIndex[4] = classloader.CpEntry{Type: classloader.UTF8, Slot: 0}
			CP.CpIndex[5] = classloader.CpEntry{Type: classloader.UTF8, Slot: 1}
			CP.CpIndex[6] = classloader.CpEntry{Type: classloader.StringConst, Slot: 7}
			CP.CpIndex[7] = classloader.CpEntry{Type: classloader.UTF8, Slot: 2} // point to UTF8[2]

			CP.MethodRefs = make([]classloader.MethodRefEntry, 1)
			methRef := classloader.MethodRefEntry{
				ClassIndex:  2, // these are CP entries
				NameAndType: 3,
			}
			CP.MethodRefs[0] = methRef

			CP.ClassRefs = append(CP.ClassRefs, stringPool.GetStringIndex(&objClassName))

			CP.Utf8Refs = append(CP.Utf8Refs, methName)
			CP.Utf8Refs = append(CP.Utf8Refs, methType)
			CP.Utf8Refs = append(CP.Utf8Refs, stringParam)
			nAndT := classloader.NameAndTypeEntry{
				NameIndex: uint16(4),
				DescIndex: uint16(5),
			}
			CP.NameAndTypes = append(CP.NameAndTypes, nAndT)

			var f frames.Frame
			if stringParam != "" { // if there is a string param, load it with LDC before calling INVOKEVIRTUAL
				f = newFrame(opcodes.LDC)
				f.Meth = append(f.Meth, uint8(6)) // point to the string constant parameter indexed by CPindex[6]
				f.Meth = append(f.Meth, opcodes.INVOKEVIRTUAL)
			} else {
				f = newFrame(opcodes.INVOKEVIRTUAL)
			}
			f.Meth = append(f.Meth, 0x00)
			f.Meth = append(f.Meth, 0x01) // Go to method referred to in 0x0001 of the CP

			f.CP = &CP

			// create the opStack
			for j := 0; j < 10; j++ {
				f.OpStack = append(f.OpStack, 0)
			}

			fs := frames.CreateFrameStack()

			// now push a reference to the object whose method we're calling. In the event, it's a prinstream,
			// we force it to be stdout. Otherwise, we instantiate the class.
			if objClassName == "java/io/PrintStream" { // if we're working with a printstream, force-set it to stdout
				f.OpStack[0] = os.Stdout
			} else {
				objPtr, err := InstantiateClass(objClassName, fs)
				if err != nil {
					errMsg := fmt.Sprintf("in test %s, could not instantiate class object: %s  %v",
						name, objClassName, err)
					t.Skip(errMsg)
				} else {
					f.OpStack[0] = objPtr.(*object.Object)
				}
			}
			f.TOS = 0

			fs.PushFront(&f) // push the new frame

			// run the method (by running the frame)
			_ = runFrame(fs)

			// get contents written by stderr and stdout, then restore
			// stderr and stdout to what they were before the test
			_ = werr.Close()
			rawStderrMsg, _ := io.ReadAll(rerr)
			os.Stderr = normalStderr

			_ = wout.Close()
			rawStdoutMsg, _ := io.ReadAll(rout)
			os.Stdout = normalStdout

			// convert the contents written to stderr and stdout into strings
			// and run tests on those strings

			errMsg := string(rawStderrMsg[:])
			if len(errMsg) != 0 && !strings.Contains(errMsg, test.stderrText) {
				t.Errorf("gfunctionExec: Test %s, expected error msg: %s, got: %s",
					name, stderrExpected, errMsg)
			}

			outMsg := string(rawStdoutMsg[:])
			if len(outMsg) != 0 && !strings.Contains(outMsg, test.stdoutText) {
				if !strings.Contains(outMsg, test.stdoutText) {
					t.Errorf("gfunctionExec: Test %s, expected output: %s, got: %s",
						name, stdoutExpected, outMsg)
				}
			}
		})
	}
}
