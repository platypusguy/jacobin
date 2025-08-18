/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"jacobin/src/excNames"
	"jacobin/src/object"
)

// TimeUnit constants
const (
	NANOSECONDS  = "NANOSECONDS"
	MICROSECONDS = "MICROSECONDS"
	MILLISECONDS = "MILLISECONDS"
	SECONDS      = "SECONDS"
	MINUTES      = "MINUTES"
	HOURS        = "HOURS"
	DAYS         = "DAYS"
)

// Conversion factors
var timeUnitConversion = map[string]map[string]int64{
	NANOSECONDS: {
		NANOSECONDS:  1,
		MICROSECONDS: 1000,
		MILLISECONDS: 1000000,
		SECONDS:      1000000000,
		MINUTES:      60000000000,
		HOURS:        3600000000000,
		DAYS:         86400000000000,
	},
	MICROSECONDS: {
		NANOSECONDS:  1 / 1000,
		MICROSECONDS: 1,
		MILLISECONDS: 1000,
		SECONDS:      1000000,
		MINUTES:      60000000,
		HOURS:        3600000000,
		DAYS:         86400000000,
	},
	MILLISECONDS: {
		NANOSECONDS:  1 / 1000000,
		MICROSECONDS: 1 / 1000,
		MILLISECONDS: 1,
		SECONDS:      1000,
		MINUTES:      60000,
		HOURS:        3600000,
		DAYS:         86400000,
	},
	SECONDS: {
		NANOSECONDS:  1 / 1000000000,
		MICROSECONDS: 1 / 1000000,
		MILLISECONDS: 1 / 1000,
		SECONDS:      1,
		MINUTES:      60,
		HOURS:        3600,
		DAYS:         86400,
	},
	MINUTES: {
		NANOSECONDS:  1 / 60000000000,
		MICROSECONDS: 1 / 60000000,
		MILLISECONDS: 1 / 60000,
		SECONDS:      1 / 60,
		MINUTES:      1,
		HOURS:        60,
		DAYS:         1440,
	},
	HOURS: {
		NANOSECONDS:  1 / 3600000000000,
		MICROSECONDS: 1 / 3600000000,
		MILLISECONDS: 1 / 3600000,
		SECONDS:      1 / 3600,
		MINUTES:      1 / 60,
		HOURS:        1,
		DAYS:         24,
	},
	DAYS: {
		NANOSECONDS:  1 / 86400000000000,
		MICROSECONDS: 1 / 86400000000,
		MILLISECONDS: 1 / 86400000,
		SECONDS:      1 / 86400,
		MINUTES:      1 / 1440,
		HOURS:        1 / 24,
		DAYS:         1,
	},
}

// toMillis converts the given time to milliseconds
func toMillis(params []interface{}) interface{} {
	unit := params[0].(*object.Object)
	time := params[1].(int64)

	unitName := object.GoStringFromStringObject(unit)
	conversionFactor, ok := timeUnitConversion[unitName][MILLISECONDS]
	if !ok {
		errMsg := "toMillis: invalid TimeUnit"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	return time * conversionFactor
}

// toSeconds converts the given time to seconds
func toSeconds(params []interface{}) interface{} {
	unit := params[0].(*object.Object)
	time := params[1].(int64)

	unitName := object.GoStringFromStringObject(unit)
	conversionFactor, ok := timeUnitConversion[unitName][SECONDS]
	if !ok {
		errMsg := "toSeconds: invalid TimeUnit"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	return time * conversionFactor
}

// toMinutes converts the given time to minutes
func toMinutes(params []interface{}) interface{} {
	unit := params[0].(*object.Object)
	time := params[1].(int64)

	unitName := object.GoStringFromStringObject(unit)
	conversionFactor, ok := timeUnitConversion[unitName][MINUTES]
	if !ok {
		errMsg := "toMinutes: invalid TimeUnit"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	return time * conversionFactor
}

// toHours converts the given time to hours
func toHours(params []interface{}) interface{} {
	unit := params[0].(*object.Object)
	time := params[1].(int64)

	unitName := object.GoStringFromStringObject(unit)
	conversionFactor, ok := timeUnitConversion[unitName][HOURS]
	if !ok {
		errMsg := "toHours: invalid TimeUnit"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	return time * conversionFactor
}

// toDays converts the given time to days
func toDays(params []interface{}) interface{} {
	unit := params[0].(*object.Object)
	time := params[1].(int64)

	unitName := object.GoStringFromStringObject(unit)
	conversionFactor, ok := timeUnitConversion[unitName][DAYS]
	if !ok {
		errMsg := "toDays: invalid TimeUnit"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	return time * conversionFactor
}
