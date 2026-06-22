/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package sunSecurity

import (
	"jacobin/src/gfunction/ghelpers"
)

func Load_Sun_Security_Jca_Providers() {

	ghelpers.MethodSignatures["sun/security/jca/Providers.getProviderList()Lsun/security/jca/ProviderList;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["sun/security/jca/Providers.setProviderList(Lsun/security/jca/ProviderList;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["sun/security/jca/Providers.getFullProviderList()Lsun/security/jca/ProviderList;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["sun/security/jca/Providers.getThreadProviderList()Lsun/security/jca/ProviderList;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["sun/security/jca/Providers.beginThreadProviderList(Lsun/security/jca/ProviderList;)Lsun/security/jca/ProviderList;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["sun/security/jca/Providers.endThreadProviderList(Lsun/security/jca/ProviderList;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}
}
