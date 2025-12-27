package gfunction

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"testing"
)

func TestProcess_BasicMethods(t *testing.T) {
	globals.InitStringPool()

	// Use current process ID for testing
	currentPid := int64(os.Getpid())

	// Create a mock Process object
	procObj := object.MakeOneFieldObject("java/lang/Process", "pid", types.Int, currentPid)

	// Test pid()
	res := processPid([]interface{}{procObj})
	if res.(int64) != currentPid {
		t.Errorf("pid() expected %d, got %v", currentPid, res)
	}

	// Test isAlive()
	res = processIsAlive([]interface{}{procObj})
	if res != types.JavaBoolTrue {
		t.Errorf("isAlive() expected true for current process")
	}

	// Test toHandle()
	res = processToHandle([]interface{}{procObj})
	handleObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("toHandle() did not return an object")
	}
	hPid := handleObj.FieldTable["pid"].Fvalue.(int64)
	if hPid != currentPid {
		t.Errorf("ProcessHandle PID expected %d, got %d", currentPid, hPid)
	}
}

func TestProcess_InvalidObject(t *testing.T) {
	globals.InitStringPool()

	// Object without pid field
	procObj := object.MakeEmptyObject()

	res := processPid([]interface{}{procObj})
	if _, ok := res.(*GErrBlk); !ok {
		t.Errorf("pid() should return error for object without pid field")
	}
}
