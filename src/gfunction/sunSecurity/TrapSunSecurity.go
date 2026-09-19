/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package sunSecurity

import "jacobin/src/gfunction/ghelpers"

func Load_Traps_Sun_Security() {

	ghelpers.MethodSignatures["sun/security/util/Debug.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["sun/security/util/Debug.getInstance(Ljava/lang/String;)Lsun/security/util/Debug;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.ReturnNull,
		}

	ghelpers.MethodSignatures["sun/security/util/Debug.getInstance(Ljava/lang/String;Ljava/lang/String;)Lsun/security/util/Debug;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.ReturnNull,
		}

	ghelpers.MethodSignatures["sun/security/util/Debug.isOn(Ljava/lang/String;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.ReturnFalse,
		}

}
