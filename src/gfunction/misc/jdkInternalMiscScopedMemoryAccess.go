/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package misc

import "jacobin/src/gfunction/ghelpers"

/*
 Each object or library that has Go methods contains a reference to ghelpers.MethodSignatures,
 which contain data needed to insert the go method into the MTable of the currently
 executing JVM. ghelpers.MethodSignatures is a map whose key is the fully qualified name and
 type of the method (that is, the method's full signature) and a value consisting of
 a struct of an int (the number of slots to pop off the caller's operand stack when
 creating the new frame and a function. All methods have the same signature, regardless
 of the signature of their Java counterparts. That signature is that it accepts a slice
 of interface{} and returns an interface{}. The accepted slice can be empty and the
 return interface can be nil. This covers all Java functions. (Objects are returned
 as a 64-bit address in this scheme (as they are in the JVM).

 The passed-in slice contains one entry for every parameter passed to the method (which
 could mean an empty slice).
*/

func Load_Jdk_Internal_Misc_ScopedMemoryAccess() {

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.registerNatives()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.copyMemory(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;JLjava/lang/Object;JJ)V"] =
		ghelpers.GMeth{
			ParamSlots: 7,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.copySwapMemory(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;JLjava/lang/Object;JJJ)V"] =
		ghelpers.GMeth{
			ParamSlots: 8,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.setMemory(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;JJB)V"] =
		ghelpers.GMeth{
			ParamSlots: 6,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.getBoolean(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;J)Z"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.putBoolean(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;JZ)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.getByte(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;J)B"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.putByte(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;JB)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.getChar(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;JZ)C"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.putChar(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;JCZ)V"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.getShort(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;JZ)S"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.putShort(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;JSZ)V"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.getInt(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;JZ)I"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.putInt(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;JIZ)V"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.getLong(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;JZ)J"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.putLong(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;JJZ)V"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.getFloat(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;JZ)F"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.putFloat(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;JFZ)V"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.getDouble(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;JZ)D"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.putDouble(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;JDZ)V"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.getReference(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;J)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/ScopedMemoryAccess.putReference(Ljdk/internal/misc/ScopedMemoryAccess$Scope;Ljava/lang/Object;JLjava/lang/Object;)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

}
