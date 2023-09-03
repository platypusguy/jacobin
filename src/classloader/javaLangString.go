/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"jacobin/log"
	"jacobin/types"
)

/*
   We don't run String's static initializer block because the initialization
   is already handled in String creation
*/

func Load_Lang_String() map[string]GMeth {
	// need to replace eventually by enbling the Java intializer to run
	MethodSignatures["java/lang/String.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  stringClinit,
		}

	return MethodSignatures
}

func stringClinit([]interface{}) interface{} {
	klass := MethAreaFetch("java/lang/String")
	if klass == nil {
		errMsg := "In <clinit>, expected java/lang/String to be in the MethodArea, but it was not"
		_ = log.Log(errMsg, log.SEVERE)
	}
	klass.Data.ClInit = types.ClInitRun // just mark that String.<clinit>() has been run
	return nil
}
