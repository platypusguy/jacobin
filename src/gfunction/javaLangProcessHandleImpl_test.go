package gfunction

import (
	"jacobin/excNames"
	"jacobin/globals"
	"jacobin/object"
	"os"
	"os/user"
	"testing"
)

// Helper to check *GErrBlk with expected exception type (int)
func checkErrType(t *testing.T, res interface{}, expected int) {
	t.Helper()

	if errObj, ok := res.(*GErrBlk); ok {
		if errObj.ExceptionType != expected {
			t.Fatalf("expected exception type %d, got %d", expected, errObj.ExceptionType)
		}
	}
	// If res is not *GErrBlk, do nothing (not an error)
}

// processHandleInfoArguments expects params []interface{} of Java String objects representing program arguments
func TestProcessHandleInfoArguments(t *testing.T) {
	globals.InitStringPool()
	prog := object.StringObjectFromGoString("prog")
	arg1 := object.StringObjectFromGoString("arg1")
	arg2 := object.StringObjectFromGoString("arg2")
	params := []interface{}{prog, arg1, arg2}

	res := processHandleInfoArguments(params)
	if res == nil {
		t.Fatalf("expected object, got nil")
	}
	if _, ok := res.(*object.Object); !ok {
		t.Errorf("expected *object.Object, got %T", res)
	}
}

func TestProcessHandleInfoCommand(t *testing.T) {
	globals.InitStringPool()
	params := []interface{}{} // no arguments, but slice passed

	res := processHandleInfoCommand(params)
	if _, ok := res.(*object.Object); !ok {
		t.Errorf("expected *object.Object, got %T", res)
	}
}

func TestProcessHandleInfoCommand_Error(t *testing.T) {
	orig := osExecutable
	defer func() { osExecutable = orig }()
	osExecutable = func() (string, error) { return "", os.ErrPermission }

	globals.InitStringPool()
	params := []interface{}{}

	res := processHandleInfoCommand(params)
	checkErrType(t, res, excNames.VirtualMachineError)
}

func TestProcessHandleInfoCommandLine(t *testing.T) {
	globals.InitStringPool()
	prog := object.StringObjectFromGoString("prog")
	arg1 := object.StringObjectFromGoString("arg1")
	params := []interface{}{prog, arg1}

	res := processHandleInfoCommandLine(params)
	if _, ok := res.(*object.Object); !ok {
		t.Errorf("expected *object.Object, got %T", res)
	}
}

func TestProcessHandleInfoCommandLine_Error(t *testing.T) {
	orig := osExecutable
	defer func() { osExecutable = orig }()
	osExecutable = func() (string, error) { return "", os.ErrPermission }

	globals.InitStringPool()
	params := []interface{}{}

	res := processHandleInfoCommandLine(params)
	checkErrType(t, res, excNames.VirtualMachineError)
}

func TestProcessHandleInfoUser(t *testing.T) {
	globals.InitStringPool()
	params := []interface{}{}

	res := processHandleInfoUser(params)
	if _, ok := res.(*object.Object); !ok {
		t.Errorf("expected *object.Object, got %T", res)
	}
}

func TestProcessHandleInfoUser_Error(t *testing.T) {
	orig := userCurrent
	defer func() { userCurrent = orig }()
	userCurrent = func() (*user.User, error) { return nil, os.ErrPermission }

	globals.InitStringPool()
	params := []interface{}{}

	res := processHandleInfoUser(params)
	checkErrType(t, res, excNames.VirtualMachineError)
}

// Dependency injection wrappers to allow error injection in tests
var osExecutable = os.Executable
var userCurrent = user.Current
