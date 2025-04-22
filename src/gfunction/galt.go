/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

func Load_Experiment() {

	MethodSignatures["java/io/FileInputStream.initIDs()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  returnTrue,
		}

	MethodSignatures["java/io/UnixFileSystem.initIDs()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  returnTrue,
		}

	MethodSignatures["java/lang/Class.desiredAssertionStatus()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  returnTrue,
		}

	MethodSignatures["java/lang/Class.desiredAssertionStatus0()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  returnTrue,
		}

	MethodSignatures["java/lang/Class.registerNatives()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/Runtime.availableProcessors()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  runtimeAvailableProcessors,
		}

	MethodSignatures["java/lang/System.registerNatives()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/System.currentTimeMillis()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  systemCurrentTimeMillis,
		}

	MethodSignatures["java/lang/System.nanoTime()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  systemNanoTime,
		}

	MethodSignatures["java/lang/Thread.registerNatives()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicLong.VMSupportsCS8()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  returnTrue,
		}

	MethodSignatures["jdk/internal/misc/CDS.isDumpingArchive0()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  returnFalse,
		}

	MethodSignatures["jdk/internal/misc/CDS.isDumpingClassList0()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  returnFalse,
		}

	MethodSignatures["jdk/internal/misc/CDS.initializeFromArchive(Ljava/lang/Class;)V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["jdk/internal/misc/CDS.isSharingEnabled0()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  returnFalse,
		}

	MethodSignatures["jdk/internal/misc/Unsafe.registerNatives()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["jdk/internal/misc/VM.initialize()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

}
