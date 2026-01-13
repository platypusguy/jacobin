/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import "jacobin/src/gfunction/ghelpers"

func Load_Util_Logging_Logger() {

	ghelpers.MethodSignatures["java/util/logging/Logger.<clinit>()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapClass}

	ghelpers.MethodSignatures["java/util/logging/Logger.config(Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/util/logging/Logger.entering(Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/util/logging/Logger.entering(Ljava/lang/String;Ljava/lang/String;Ljava/lang/Object;)V"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/util/logging/Logger.exiting(Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/util/logging/Logger.exiting(Ljava/lang/String;Ljava/lang/String;Ljava/lang/Object;)V"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/util/logging/Logger.finest(Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/util/logging/Logger.finer(Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/util/logging/Logger.info(Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/util/logging/Logger.log(Ljava/util/logging/Level;Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/util/logging/Logger.log(Ljava/util/logging/Level;Ljava/lang/String;[Ljava/lang/Object;)V"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/util/logging/Logger.logp(Ljava/util/logging/Level;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 4, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/util/logging/Logger.logp(Ljava/util/logging/Level;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;[Ljava/lang/Object;)V"] =
		ghelpers.GMeth{ParamSlots: 5, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/util/logging/Logger.logrb(Ljava/util/logging/Level;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 5, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/util/logging/Logger.logrb(Ljava/util/logging/Level;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;[Ljava/lang/Object;)V"] =
		ghelpers.GMeth{ParamSlots: 6, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/util/logging/Logger.severe(Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/util/logging/Logger.throwing(Ljava/lang/String;Ljava/lang/String;Ljava/lang/Throwable;)V"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/util/logging/Logger.warning(Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/util/logging/Logger.warning(Ljava/lang/String;[Ljava/lang/Object;)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
}
