/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package sunSecurity

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaLang"
)

func Load_Sun_Security_Action_GetPropertyAction() {

	ghelpers.MethodSignatures["sun/security/action/GetPropertyAction.privilegedGetProperties()Ljava/util/Properties;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  javaLang.SystemGetProperties,
		}

	ghelpers.MethodSignatures["sun/security/action/GetPropertyAction.privilegedGetProperty(Ljava/lang/String;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  javaLang.SystemGetProperty,
		}

	ghelpers.MethodSignatures["sun/security/action/GetPropertyAction.privilegedGetProperty(Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  javaLang.SystemGetProperty,
		}

	ghelpers.MethodSignatures["sun/security/action/GetPropertyAction.privilegedGetTimeoutProp(Ljava/lang/String;ILsun/security/util/Debug;)I"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

}
