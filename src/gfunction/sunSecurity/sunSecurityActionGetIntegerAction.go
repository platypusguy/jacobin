/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package sunSecurity

import (
	"jacobin/src/gfunction/ghelpers"
)

func Load_Sun_Security_Action_GetIntegerAction() {

	ghelpers.MethodSignatures["sun/security/action/GetIntegerAction.privilegedGetProperty(Ljava/lang/String;)Ljava/lang/Integer;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["sun/security/action/GetIntegerAction.privilegedGetProperty(Ljava/lang/String;I)Ljava/lang/Integer;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["sun/security/action/GetIntegerAction.run()Ljava/lang/Integer;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["sun/security/action/GetIntegerAction.run()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}
}
