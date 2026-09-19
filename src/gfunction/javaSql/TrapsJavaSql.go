/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSql

import "jacobin/src/gfunction/ghelpers"

func Load_Traps_Java_Sql() {

	ghelpers.MethodSignatures["java/sql/Driver.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapClass,
		}

	ghelpers.MethodSignatures["java/sql/DriverAction.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapClass,
		}

	ghelpers.MethodSignatures["java/sql/DriverPropertyInfo.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapClass,
		}

	ghelpers.MethodSignatures["java/sql/DriverPropertyInfo.<init>(Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapClass,
		}

	ghelpers.MethodSignatures["java/sql/DriverManager.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapClass,
		}

}
