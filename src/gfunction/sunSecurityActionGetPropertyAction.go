/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

func Load_Sun_Security_Action_GetPropertyAction() {

	MethodSignatures["sun/security/action/GetPropertyAction.privilegedGetProperties()Ljava/util/Properties;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  systemGetProperties,
		}

	MethodSignatures["sun/security/action/GetPropertyAction.privilegedGetProperty(Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  systemGetProperty,
		}

	MethodSignatures["sun/security/action/GetPropertyAction.privilegedGetProperty(Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  systemGetProperty,
		}

	MethodSignatures["sun/security/action/GetPropertyAction.privilegedGetTimeoutProp(Ljava/lang/String;ILsun/security/util/Debug;)I"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

}
