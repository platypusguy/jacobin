/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/types"
)

// Load_Util_Logging_Filter loads the method signatures for the
// java/util/logging/Filter interface. Filter is a functional interface
// whose sole abstract method is isLoggable(LogRecord), which is not
// deprecated, so it uses ghelpers.TrapFunction.
func Load_Util_Logging_Filter() {

	ghelpers.MethodSignatures["java/util/logging/Filter.isLoggable(Ljava/util/logging/LogRecord;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingFilterIsLoggable,
		}
}

/*
 * A Filter can be used to provide fine grain control over what is logged, beyond the control provided by log levels.
 * Each Logger and each Handler can have a filter associated with it.
 * The Logger or Handler will call the isLoggable method to check if a given LogRecord should be published.
 * If isLoggable returns false, the LogRecord will be discarded.
 */
func loggingFilterIsLoggable([]interface{}) interface{} {
	// Allow all LogRecords to be published.
	// Normally an object implementing Filter would return true or false
	// based on the LogRecord level or some other criteria.
	return types.JavaBoolTrue
}
