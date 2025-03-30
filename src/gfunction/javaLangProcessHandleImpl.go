/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import "os"

func Load_Lang_Process_Handle_Impl() {

	MethodSignatures["java/lang/ProcessHandleImpl.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/lang/ProcessHandleImpl.initNative()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/ProcessHandleImpl.getCurrentPid0()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  phimplCurrentPid,
		}

}

func phimplCurrentPid(params []interface{}) interface{} {
	return int64(os.Getpid())
}
