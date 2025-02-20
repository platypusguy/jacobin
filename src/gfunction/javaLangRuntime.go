/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/object"
	"jacobin/statics"
	"jacobin/types"
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

	MethodSignatures["java/lang/Runtime.exit(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  exitI, // javaLangSystem.go
		}

	MethodSignatures["java/lang/Runtime.getRuntime()Ljava/lang/Runtime;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  runtimeGetRuntime,
		}

	MethodSignatures["java/lang/Runtime.halt(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  exitI, // javaLangSystem.go
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

	MethodSignatures["java/lang/Runtime.version()Ljava/lang/Runtime/Version;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
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
	return object.StringObjectFromGoString(stringClassnameRuntime)
}

// runtimeGetRuntime: Get the singleton Runtime object.
func runtimeGetRuntime([]interface{}) interface{} {
	return statics.GetStaticValue(stringClassnameRuntime, stringFieldCurrentRuntime)
}

// runtimeAvailableProcessors: Get the number of CPU cores.
func runtimeAvailableProcessors([]interface{}) interface{} {
	return int64(runtime.NumCPU())
}
