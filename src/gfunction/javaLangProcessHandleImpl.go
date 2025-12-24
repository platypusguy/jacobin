/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"os"
	"os/user"
)

func Load_Lang_Process_Handle_Impl() {

	MethodSignatures["java/lang/ProcessHandle.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processHandleImplClinit,
		}

	MethodSignatures["java/lang/ProcessHandle.allProcesses()Ljava/util/Stream;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessHandle.children()Ljava/util/Stream;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessHandle.compareTo(Ljava/lang/ProcessHandle;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessHandle.current()Ljava/lang/ProcessHandle;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processHandleImplCurrent,
		}

	MethodSignatures["java/lang/ProcessHandle.descendents()Ljava/util/Stream;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessHandle.destroy()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction, // O/S dependent
		}

	MethodSignatures["java/lang/ProcessHandle.destroyForcibly()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction, // O/S dependent
		}

	MethodSignatures["java/lang/ProcessHandle.equals(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  processHandleImplEquals,
		}

	MethodSignatures["java/lang/ProcessHandle.getCurrentPid()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processHandleImplCurrentPid,
		}

	MethodSignatures["java/lang/ProcessHandle.getCurrentPid0()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processHandleImplCurrentPid,
		}

	MethodSignatures["java/lang/ProcessHandle.hashCode()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessHandle.info()Ljava/lang/ProcessHandle$Info;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processHandleImplInfo,
		}

	MethodSignatures["java/lang/ProcessHandle.initNative()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/ProcessHandle.isAlive()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processHandleImplIsAlive,
		}

	MethodSignatures["java/lang/ProcessHandle.isAlive0(J)J"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  processHandleImplIsAlive0,
		}

	MethodSignatures["java/lang/ProcessHandle.of(J)Ljava/util/Optional;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessHandle.onExit()Ljava/util.concurrent/CompletableFuture;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessHandle.parent()Ljava/util/Optional;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessHandle.pid()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processHandleImplCurrentPid,
		}

	MethodSignatures["java/lang/ProcessHandle.supportsNormalTermination()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessHandle$Info.arguments()Ljava/util/Optional;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processHandleInfoArguments,
		}

	MethodSignatures["java/lang/ProcessHandle$Info.command()Ljava/util/Optional;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processHandleInfoCommand,
		}

	MethodSignatures["java/lang/ProcessHandle$Info.commandLine()Ljava/util/Optional;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processHandleInfoCommandLine,
		}

	MethodSignatures["java/lang/ProcessHandle$Info.startInstant()Ljava/util/Optional;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessHandle$Info.totalCpuDuration()Ljava/util/Optional;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/ProcessHandle$Info.user()Ljava/util/Optional;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  processHandleInfoUser,
		}

}

var classNameProcessHandle string = "java/lang/ProcessHandle"
var classNameProcessHandleImpl string = "java/lang/ProcessHandleImpl"
var classNameProcessHandleInfo string = "java/lang/ProcessHandle$Info"
var pidFieldName string = "pid"
var pid int64

func processHandleImplClinit(params []interface{}) interface{} {
	statics.AddStatic(classNameProcessHandleImpl+".NOT_A_CHILD", statics.Static{
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

// processHandleImplIsAlive checks to see whether or not the ProcessHandleImpl object contains a pid that represents
// a process that is still alive.
func processHandleImplIsAlive(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "processHandleImplIsAlive: ProcessHandleImpl parameter is not an object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	processId := obj.FieldTable[pidFieldName].Fvalue.(int64)
	_, err := os.FindProcess(int(processId))
	if err != nil {
		return types.JavaBoolFalse
	}
	// Note: process.Signal is not implemented on Windows
	// err = process.Signal(syscall.Signal(0))
	// return err == nil
	return types.JavaBoolTrue
}

// processHandleImplIsAlive0 is a low-level function that checks if a process is alive and
// returns the start time in milliseconds since 1970, 0 if the start time cannot be determined,
// and -1 if the process is not alive.
func processHandleImplIsAlive0(params []interface{}) interface{} {
	return int64(0) // golang has no O/S-independent way to get the start time of a process
}

// processHandleImplCurrent creates a java/lang/ProcessHandle object with the pid field set to this process ID.
func processHandleImplCurrent(params []interface{}) interface{} {
	obj := object.MakeOneFieldObject(classNameProcessHandle, pidFieldName, types.Int, pid)
	return obj
}

// processHandleImplEquals determines whether the argument ProcessHandleImpl holds a pid
// that is the same as the current process.
func processHandleImplEquals(params []interface{}) interface{} {
	this, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "processHandleImplEquals: ProcessHandleImpl parameter is not an object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	that, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "processHandleImplEquals: Parameter is not an object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	if *stringPool.GetStringPointer(this.KlassName) != classNameProcessHandle {
		return types.JavaBoolFalse
	}
	thatPid, ok := that.FieldTable[pidFieldName].Fvalue.(int64)
	if !ok {
		return types.JavaBoolFalse
	}
	if thatPid != pid {
		return types.JavaBoolFalse
	}

	return types.JavaBoolTrue
}

func processHandleImplInfo(params []interface{}) interface{} {
	return object.MakeEmptyObjectWithClassName(&classNameProcessHandleInfo)
}

func processHandleInfoArguments(params []interface{}) interface{} {
	args := os.Args[1:]
	strObjArray := object.StringObjectArrayFromGoStringArray(args)
	outObj := object.MakePrimitiveObject(classNameOptional, types.RefArray, strObjArray)
	return outObj
}

func processHandleInfoCommand(params []interface{}) interface{} {
	exePath, err := os.Executable()
	if err != nil {
		errMsg := "processHandleInfoCommand: os.Executable() failed, err: " + err.Error()
		return getGErrBlk(excNames.VirtualMachineError, errMsg)
	}
	strObj := object.StringObjectFromGoString(exePath)
	outObj := object.MakePrimitiveObject(classNameOptional, types.StringClassName, strObj)
	return outObj
}

func processHandleInfoCommandLine(params []interface{}) interface{} {
	args := os.Args[1:]
	cmdLine, err := os.Executable()
	if err != nil {
		errMsg := "processHandleInfoCommandLine: os.Executable() failed, err: " + err.Error()
		return getGErrBlk(excNames.VirtualMachineError, errMsg)
	}
	for _, arg := range args {
		cmdLine += " " + arg
	}

	strObj := object.StringObjectFromGoString(cmdLine)
	outObj := object.MakePrimitiveObject(classNameOptional, types.StringClassName, strObj)
	return outObj
}

func processHandleInfoUser(params []interface{}) interface{} {
	user, err := user.Current()
	if err != nil {
		errMsg := "processHandleInfoUser: user.Current() failed, err: " + err.Error()
		return getGErrBlk(excNames.VirtualMachineError, errMsg)
	}
	strObj := object.StringObjectFromGoString(user.Username)
	outObj := object.MakePrimitiveObject(classNameOptional, types.StringClassName, strObj)
	return outObj
}
