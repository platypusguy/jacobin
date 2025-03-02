package testutil

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/frames"
	"jacobin/gfunction"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/statics"
	"os"
	"testing"
)

// ***** NOT THREAD-SAFE *****
var originalStderr *os.File
var originalStdout *os.File
var testStderr *os.File
var testStdout *os.File
var flagInit = false

// Standard unit test initialisation without processing stderr and stdout.
func UTinit(t *testing.T) {
	if flagInit {
		return
	}
	globals.InitGlobals("test")
	statics.Statics = make(map[string]statics.Static)
	statics.PreloadStatics()
	err := classloader.Init()
	if err != nil {
		t.Fatalf("stdInit: classloader.Init() failed, err: %s", err.Error())
	}
	classloader.LoadBaseClasses()
	classloader.MTable = make(map[string]classloader.MTentry)
	gfunction.MTableLoadGFunctions(&classloader.MTable)
	flagInit = true
}

// Save current stdout and stderr.
// Substitute 2 ones while executing Jacobin functions.
func UTnewConsole(t *testing.T) {
	var err error
	var dummy *os.File

	// Swap in a new stdout.
	originalStdout = os.Stdout
	dummy, testStdout, err = os.Pipe()
	_ = dummy.Close()
	if err != nil {
		os.Stdout = originalStdout // Restore original stdout.
		t.Fatalf("stdTestBegin: os.Pipe()-->stdout failed, err: %s", err.Error())
	}
	os.Stdout = testStdout

	// Swap in a new stderr.
	originalStderr = os.Stderr
	dummy, testStderr, err = os.Pipe()
	_ = dummy.Close()
	if err != nil {
		_ = testStdout.Close()
		os.Stdout = originalStdout // Restore original stdout.
		os.Stderr = originalStderr // Restore original stderr.
		t.Fatalf("stdTestBegin: os.Pipe()-->stderr failed, err: %s", err.Error())
	}
	os.Stderr = testStderr
}

// Close testing stderr & stdout; restore the original stderr & stdout.
func UTrestoreConsole(t *testing.T) {
	_ = testStderr.Close()
	_ = testStdout.Close()
	os.Stdout = originalStdout
	os.Stderr = originalStderr
}

// Execute a G function.
// Return result to caller.
func UTgfunc(t *testing.T, className, methodName, methodType string, obj *object.Object, args []interface{}) interface{} {

	// String form of FQN.
	fqn := fmt.Sprintf("%s.%s%s", className, methodName, methodType)

	// Initialize Jacobin infrastructure and set up new stdout and stderr.
	UTinit(t)
	UTnewConsole(t)

	// Create empty frame stack (fs).
	fs := frames.CreateFrameStack()

	// Create frame (fr).
	fr := frames.CreateFrame(3)
	fr.Thread = 0 // Mainthread
	fr.FrameStack = fs
	fr.ClName = className
	fr.MethName = methodName
	fr.MethType = methodType

	// Push fr to front of fs.
	_ = frames.PushFrame(fs, fr)

	// Create mtEntry.
	mtEntry := classloader.MTable[fqn]
	if mtEntry.Meth == nil {
		t.Fatalf("UTgfunc: classloader.MTable[%s] not found", fqn)
	}

	paramCount := len(args)

	// params = args in reverse order (expected by RunGfunction).
	params := make([]interface{}, paramCount)
	for ix := 0; ix < paramCount; ix++ {
		params[ix] = args[paramCount-1-ix]
	}

	// Add the object reference (Java class or file I/O).
	if obj != nil {
		params = append(params, obj)
	}

	// Run the G function.
	result := gfunction.RunGfunction(mtEntry, fs, className, methodName, methodType, &params, true, false)

	// Restore previous stderr and stdout.
	UTrestoreConsole(t)

	// Return the result to caller.
	return result

}
