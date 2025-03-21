/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/trace"
	"jacobin/types"
)

// TODO: Delete this file and the reference to Load_Lang_SecurityManager() in gfunction.go
// TODO: and javaLangSystem.go when the Java library stops using the SecurityManager class.

func Load_Lang_SecurityManager() {

	MethodSignatures["java/lang/SecurityManager.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/lang/SecurityManager.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  secmgrInit,
		}

	MethodSignatures["java/lang/SecurityManager.checkAccept(Ljava/lang/String;I)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkAccess(Ljava/lang/Thread)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkAccess(Ljava/lang/ThreadGroup)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkConnect(Ljava/lang/String;I)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkConnect(Ljava/lang/String;ILjava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkCreateClassLoader()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkDelete(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkExec(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkExit(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkLink(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkListen(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkMulticast(Ljava/net/InetAddress;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkPackageAccess(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkPackageDefinition(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkPermission(Ljava/security/Permission;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkPrintJobAccess()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkPropertiesAccess()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkPropertyAccess(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkRead(Ljava/io/FileDescriptor;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkRead(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkRead(Ljava/lang/String;Ljava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkWrite(Ljava/io/FileDescriptor;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/SecurityManager.checkWrite(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  justReturn,
		}

}

func secmgrInit(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "secmgrInit: SecurityManager parameter is not an object"
		trace.Error(errMsg)
		return getGErrBlk(excNames.VirtualMachineError, errMsg)
	}
	fld := object.Field{Ftype: types.Int, Fvalue: 42}
	obj.FieldTable["value"] = fld
	return nil
}
