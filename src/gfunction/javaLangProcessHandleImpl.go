/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/statics"
	"jacobin/types"
	"os"
	"syscall"
)

func Load_Lang_Process_Handle_Impl() {

	MethodSignatures["java/lang/ProcessHandleImpl.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processHandleImplClinit,
		}

	MethodSignatures["java/lang/ProcessHandleImpl.initNative()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/ProcessHandleImpl.getCurrentPid0()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processHandleImplCurrentPid,
		}

	MethodSignatures["java/lang/ProcessHandleImpl.getCurrentPid()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processHandleImplCurrentPid,
		}

	MethodSignatures["java/lang/ProcessHandleImpl.isAlive()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processHandleImplIsAlive,
		}

	MethodSignatures["java/lang/ProcessHandleImpl.isAlive0()J"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  processHandleImplIsAlive0,
		}

	MethodSignatures["java/lang/ProcessHandleImpl.pid()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processHandleImplCurrentPid,
		}
}

var className string = "java/lang/ProcessHandleImpl"
var pid int64

func processHandleImplClinit(params []interface{}) interface{} {
	statics.AddStatic(className+".NOT_A_CHILD", statics.Static{
		Type:  types.Int,
		Value: -2,
	})
	pid = int64(os.Getpid())

	return nil
}

func processHandleImplPid(params []interface{}) interface{} {
	// This function is called by the java.lang.ProcessHandleImpl.getCurrentPid() method.
	// It returns the current process ID.
	// The pid is stored in the static variable pid, which is set in the processHandleImplClinit function.
	// The pid is returned as a long integer.

	return pid
}

func processHandleImplCurrentPid(params []interface{}) interface{} {
	return int64(os.Getpid())
}

func processHandleImplIsAlive(params []interface{}) interface{} {
	processId := params[0].(int64)
	process, err := os.FindProcess(int(processId))
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// processHandleImplIsAlive0 is a low-level function that checks if a process is alive and
// returns the start time in milliseconds since 1970, 0 if the start time cannot be determined,
// and -1 if the process is not alive.
func processHandleImplIsAlive0(params []interface{}) interface{} {
	return 0 // golang has no way to get the start time of a process
}
