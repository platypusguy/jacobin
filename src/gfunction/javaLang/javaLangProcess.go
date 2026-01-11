/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"runtime"
	"syscall"
)

func Load_Lang_Process() {

	ghelpers.MethodSignatures["java/lang/Process.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/lang/Process.destroy()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  processDestroy,
		}

	ghelpers.MethodSignatures["java/lang/Process.destroyForcibly()Ljava/lang/Process;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  processDestroyForcibly,
		}

	ghelpers.MethodSignatures["java/lang/Process.exitValue()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  processExitValue,
		}

	ghelpers.MethodSignatures["java/lang/Process.getErrorStream()Ljava/io/InputStream;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Process.getInputStream()Ljava/io/InputStream;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Process.getOutputStream()Ljava/io/OutputStream;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/lang/Process.isAlive()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  processIsAlive,
		}

	ghelpers.MethodSignatures["java/lang/Process.waitFor()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  processWaitFor,
		}

	ghelpers.MethodSignatures["java/lang/Process.waitFor(JLjava/util/concurrent/TimeUnit;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  processWaitForTimeout,
		}

	ghelpers.MethodSignatures["java/lang/Process.pid()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  processPid,
		}

	ghelpers.MethodSignatures["java/lang/Process.info()Ljava/lang/ProcessHandle$Info;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  processInfo,
		}

	ghelpers.MethodSignatures["java/lang/Process.toHandle()Ljava/lang/ProcessHandle;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  processToHandle,
		}

	ghelpers.MethodSignatures["java/lang/Process.onExit()Ljava/util/concurrent/CompletableFuture;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}
}

func getPidFromProcessObject(obj *object.Object) (int, *ghelpers.GErrBlk) {
	val, ok := obj.FieldTable["pid"]
	if !ok {
		return 0, ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Process object missing 'pid' field")
	}
	pid, ok := val.Fvalue.(int64)
	if !ok {
		return 0, ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Process object 'pid' field is not an int64")
	}
	return int(pid), nil
}

func processDestroy(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	pid, gerr := getPidFromProcessObject(obj)
	if gerr != nil {
		return gerr
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil
	}
	_ = proc.Kill()
	return nil
}

func processDestroyForcibly(params []interface{}) interface{} {
	_ = processDestroy(params)
	return params[0]
}

func processExitValue(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	pid, gerr := getPidFromProcessObject(obj)
	if gerr != nil {
		return gerr
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IllegalThreadStateException, "Process not found")
	}
	// Note: non-blocking check is platform dependent.
	// For now we just attempt to Wait.
	state, err := proc.Wait()
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IllegalThreadStateException, "Process has not exited or is not a child")
	}
	return int64(state.ExitCode())
}

func processIsAlive(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	pid, gerr := getPidFromProcessObject(obj)
	if gerr != nil {
		return gerr
	}
	if isProcessAlive(pid) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func processWaitFor(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	pid, gerr := getPidFromProcessObject(obj)
	if gerr != nil {
		return gerr
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return int64(-1)
	}
	state, err := proc.Wait()
	if err != nil {
		return int64(-1)
	}
	return int64(state.ExitCode())
}

func processWaitForTimeout(params []interface{}) interface{} {
	// Placeholder implementation for timeout
	return processWaitFor(params[:1])
}

func processPid(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	pid, gerr := getPidFromProcessObject(obj)
	if gerr != nil {
		return gerr
	}
	return int64(pid)
}

func processInfo(params []interface{}) interface{} {
	return processHandleImplInfo(params)
}

func processToHandle(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	pid, gerr := getPidFromProcessObject(obj)
	if gerr != nil {
		return gerr
	}
	return object.MakeOneFieldObject(classNameProcessHandle, pidFieldName, types.Int, int64(pid))
}

func isProcessAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	if runtime.GOOS == "windows" {
		// On Windows, FindProcess fails if the process does not exist.
		// If it succeeded, the process exists.
		return true
	}

	// On Unix, FindProcess always succeeds, so we use Signal(0) to check if the process exists.
	err = proc.Signal(syscall.Signal(0))
	return err == nil
}
