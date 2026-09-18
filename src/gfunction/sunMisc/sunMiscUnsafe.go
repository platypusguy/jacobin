/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package sunMisc

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/jdkInternal"
)

func Load_Sun_Misc_Unsafe() {

	ghelpers.MethodSignatures["sun/misc/Unsafe.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.getUnsafe()Lsun/misc/Unsafe;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  jdkInternal.UnsafeGetUnsafe,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.arrayBaseOffset(Ljava/lang/Class;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  jdkInternal.UnsafeArrayBaseOffset,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.arrayIndexScale(Ljava/lang/Class;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  jdkInternal.UnsafeArrayIndexScale,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.objectFieldOffset(Ljava/lang/reflect/Field;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  jdkInternal.UnsafeObjectFieldOffset1,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.compareAndSwapInt(Ljava/lang/Object;JII)Z"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  jdkInternal.UnsafeCompareAndSetInt,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.compareAndSwapLong(Ljava/lang/Object;JJJ)Z"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  jdkInternal.UnsafeCompareAndSetLong,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.compareAndSwapObject(Ljava/lang/Object;JLjava/lang/Object;Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  jdkInternal.UnsafeCompareAndSetReference,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.getInt(Ljava/lang/Object;J)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  jdkInternal.UnsafeGetIntVolatile,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.putInt(Ljava/lang/Object;JI)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.getLong(Ljava/lang/Object;J)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  jdkInternal.UnsafeGetLong,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.putLong(Ljava/lang/Object;JJ)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.getObject(Ljava/lang/Object;J)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.putObject(Ljava/lang/Object;JLjava/lang/Object;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.getIntVolatile(Ljava/lang/Object;J)I"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  jdkInternal.UnsafeGetIntVolatile,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.getLongVolatile(Ljava/lang/Object;J)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  jdkInternal.UnsafeGetLong,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.getObjectVolatile(Ljava/lang/Object;J)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.allocateMemory(J)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.freeMemory(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.setMemory(Ljava/lang/Object;JJB)V"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.copyMemory(Ljava/lang/Object;JLjava/lang/Object;JJ)V"] =
		ghelpers.GMeth{
			ParamSlots: 6,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.staticFieldOffset(Ljava/lang/reflect/Field;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  jdkInternal.UnsafeObjectFieldOffset1,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.staticFieldBase(Ljava/lang/reflect/Field;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.shouldBeInitialized(Ljava/lang/Class;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.ReturnTrue,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.ensureClassInitialized(Ljava/lang/Class;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.fullFence()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.loadFence()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["sun/misc/Unsafe.storeFence()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}
}
