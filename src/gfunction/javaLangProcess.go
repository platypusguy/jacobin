/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
)

func Load_Lang_Process() {

	MethodSignatures["java/lang/Process.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/lang/Process.destroy()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processDestroy,
		}

	MethodSignatures["java/lang/Process.destroyForcibly()Ljava/lang/Process;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processDestroyForcibly,
		}

	MethodSignatures["java/lang/Process.exitValue()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processExitValue,
		}

	MethodSignatures["java/lang/Process.getErrorStream()Ljava/io/InputStream;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Process.getInputStream()Ljava/io/InputStream;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Process.getOutputStream()Ljava/io/OutputStream;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Process.isAlive()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processIsAlive,
		}

	MethodSignatures["java/lang/Process.waitFor()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processWaitFor,
		}

	MethodSignatures["java/lang/Process.waitFor(JLjava/util/concurrent/TimeUnit;)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  processWaitForTimeout,
		}

	MethodSignatures["java/lang/Process.pid()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processPid,
		}

	MethodSignatures["java/lang/Process.info()Ljava/lang/ProcessHandle$Info;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processInfo,
		}

	MethodSignatures["java/lang/Process.toHandle()Ljava/lang/ProcessHandle;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processToHandle,
		}

	MethodSignatures["java/lang/Process.onExit()Ljava/util/concurrent/CompletableFuture;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}
}

func getPidFromProcessObject(obj *object.Object) (int, *GErrBlk) {
	val, ok := obj.FieldTable["pid"]
	if !ok {
		return 0, getGErrBlk(excNames.IllegalArgumentException, "Process object missing 'pid' field")
	}
	pid, ok := val.Fvalue.(int64)
	if !ok {
		return 0, getGErrBlk(excNames.IllegalArgumentException, "Process object 'pid' field is not an int64")
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
		return getGErrBlk(excNames.IllegalThreadStateException, "Process not found")
	}
	// Note: non-blocking check is platform dependent.
	// For now we just attempt to Wait.
	state, err := proc.Wait()
	if err != nil {
		return getGErrBlk(excNames.IllegalThreadStateException, "Process has not exited or is not a child")
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
