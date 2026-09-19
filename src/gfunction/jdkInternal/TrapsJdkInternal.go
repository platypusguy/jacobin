/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package jdkInternal

import "jacobin/src/gfunction/ghelpers"

func Load_Traps_Jdk_Internal() {

	ghelpers.MethodSignatures["jdk/internal/access/SharedSecrets.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/VM.initialize()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/CDS.getRandomSeedForDumping()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnRandomLong,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/CDS.initializeFromArchive(Ljava/lang/Class;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.JustReturn,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/CDS.isDumpingArchive0()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnFalse,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/CDS.isDumpingClassList0()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnFalse,
		}

	ghelpers.MethodSignatures["jdk/internal/misc/CDS.isSharingEnabled0()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ReturnFalse,
		}

	ghelpers.MethodSignatures["jdk/internal/util/ArraysSupport.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

}
