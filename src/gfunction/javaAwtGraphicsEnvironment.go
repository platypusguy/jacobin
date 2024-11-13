/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/globals"
)

func Load_Awt_Graphics_Environment() {

	MethodSignatures["java/awt/GraphicsEnvironment.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/awt/GraphicsEnvironment.isHeadless()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  awtgeIsHeadless,
		}

	MethodSignatures["java/awt/GraphicsEnvironment.isHeadlessInstance()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  awtgeIsHeadless,
		}

}

// "java/awt/GraphicsEnvironment.isHeadless()Z"
func awtgeIsHeadless(params []interface{}) interface{} {
	glob := globals.GetGlobalRef()
	return glob.Headless
}
