/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package sunSecurity

import (
	"jacobin/src/gfunction/ghelpers"
)

func Load_Sun_Security_Action_GetLongAction() {

	ghelpers.MethodSignatures["sun/security/action/GetLongAction.privilegedGetProperty(Ljava/lang/String;)Ljava/lang/Long;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["sun/security/action/GetLongAction.privilegedGetProperty(Ljava/lang/String;J)Ljava/lang/Long;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["sun/security/action/GetLongAction.run()Ljava/lang/Long;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["sun/security/action/GetLongAction.run()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}
}
