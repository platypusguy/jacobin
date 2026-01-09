/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"math"
	"runtime"
)

func Load_Lang_Runtime() {

	MethodSignatures["java/lang/Runtime.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  runtimeClinit,
		}

	MethodSignatures["java/lang/Runtime.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Runtime.addShutdownHook(Ljava/lang/Thread;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Runtime.availableProcessors()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  runtimeAvailableProcessors,
		}

	MethodSignatures["java/lang/Runtime.exec(Ljava/lang/String;)Ljava/lang/Process;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/lang/Runtime.exec([Ljava/lang/String;)Ljava/lang/Process;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Runtime.exec([Ljava/lang/String;[Ljava/lang/String;)Ljava/lang/Process;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Runtime.exec([Ljava/lang/String;[Ljava/lang/String;Ljava/io/File;)Ljava/lang/Process;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Runtime.exec(Ljava/lang/String;[Ljava/lang/String;)Ljava/lang/Process;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/lang/Runtime.exec(Ljava/lang/String;[Ljava/lang/String;Ljava/io/File;)Ljava/lang/Process;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/lang/Runtime.exit(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  systemExitI, // javaLangSystem.go
		}

	MethodSignatures["java/lang/Runtime.freeMemory()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  freeMemory,
		}

	MethodSignatures["java/lang/Runtime.gc()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  runtimeGC,
		}

	MethodSignatures["java/lang/Runtime.getRuntime()Ljava/lang/Runtime;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  runtimeGetRuntime,
		}

	MethodSignatures["java/lang/Runtime.halt(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  systemExitI, // javaLangSystem.go
		}

	MethodSignatures["java/lang/Runtime.load(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Runtime.load0(Ljava/lang/Class;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Runtime.loadLibrary(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Runtime.loadLibrary0(Ljava/lang/Class;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Runtime.maxMemory()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  maxMemory,
		}

	MethodSignatures["java/lang/Runtime.removeShutdownHook(Ljava/lang/Thread;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Runtime.runFinalization()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/Runtime.totalMemory()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  totalMemory,
		}

	MethodSignatures["java/lang/Runtime.version()Ljava/lang/Runtime$Version;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  runtimeVersion,
		}

}

const stringClassnameRuntime = "java/lang/Runtime"
const stringFieldCurrentRuntime = "currentRuntime"

func runtimeClinit([]interface{}) interface{} {
	obj := object.MakePrimitiveObject(stringClassnameRuntime, types.ByteArray, nil)
	_ = statics.AddStatic(stringClassnameRuntime+"."+stringFieldCurrentRuntime, statics.Static{
		Type:  types.Ref + stringClassnameRuntime,
		Value: obj,
	})
	return nil
}

// runtimeGetRuntime: Get the singleton Runtime object.
func runtimeGetRuntime([]interface{}) interface{} {
	return statics.GetStaticValue(stringClassnameRuntime, stringFieldCurrentRuntime)
}

// runtimeAvailableProcessors: Get the number of CPU cores.
func runtimeAvailableProcessors([]interface{}) interface{} {
	return int64(runtime.NumCPU())
}

// freeMemory: Returns the amount of free memory in the Java Virtual Machine.
func freeMemory([]interface{}) interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// HeapIdle is memory that is idle and could be used for heap,
	// but it is still allocated from the OS.
	// HeapInuse is memory that is currently being used for heap.
	// This is a rough approximation of Java's freeMemory().
	return int64(m.HeapIdle)
}

// runtimeGC: Runs the garbage collector.
func runtimeGC([]interface{}) interface{} {
	runtime.GC()
	return nil
}

// maxMemory: Get the maximum amount of memory that the max Jacobin will attempt to use. If there is no limit,
// Java return Long.MAX_VALUE, which is what we do here
func maxMemory([]interface{}) interface{} {
	return int64(math.MaxInt64)
}

// totalMemory: Get the maximum amount of memory that the max Jacobin will attempt to use.
func totalMemory([]interface{}) interface{} {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	return int64(memStats.Sys)
}

// runtimeVersion returns the version of the runtime.
func runtimeVersion([]interface{}) interface{} {
	// For now, we return a version object that represents Java 17.
	// In a real implementation, this would be more complex.
	return object.MakePrimitiveObject("java/lang/Runtime$Version", types.Ref, nil)
}
