/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
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
		NANOSECONDS:  1000,
		MICROSECONDS: 1,
		MILLISECONDS: 1000,
		SECONDS:      1000000,
		MINUTES:      60000000,
		HOURS:        3600000000,
		DAYS:         86400000000,
	},
	MILLISECONDS: {
		NANOSECONDS:  1000000,
		MICROSECONDS: 1000,
		MILLISECONDS: 1,
		SECONDS:      1000,
		MINUTES:      60000,
		HOURS:        3600000,
		DAYS:         86400000,
	},
	SECONDS: {
		NANOSECONDS:  1000000000,
		MICROSECONDS: 1000000,
		MILLISECONDS: 1000,
		SECONDS:      1,
		MINUTES:      60,
		HOURS:        3600,
		DAYS:         86400,
	},
	MINUTES: {
		NANOSECONDS:  60000000000,
		MICROSECONDS: 60000000,
		MILLISECONDS: 60000,
		SECONDS:      60,
		MINUTES:      1,
		HOURS:        60,
		DAYS:         1440,
	},
	HOURS: {
		NANOSECONDS:  3600000000000,
		MICROSECONDS: 3600000000,
		MILLISECONDS: 3600000,
		SECONDS:      3600,
		MINUTES:      60,
		HOURS:        1,
		DAYS:         24,
	},
	DAYS: {
		NANOSECONDS:  86400000000000,
		MICROSECONDS: 86400000000,
		MILLISECONDS: 86400000,
		SECONDS:      86400,
		MINUTES:      1440,
		HOURS:        24,
		DAYS:         1,
	},
}

// toMillis converts the given time to milliseconds
func toMillis(params []interface{}) interface{} {
	unit := params[0].(*object.Object)
	duration := params[1].(int64)

	unitName := object.GoStringFromStringObject(unit)
	if isLarger(unitName, MILLISECONDS) {
		conversionFactor, ok := timeUnitConversion[unitName][MILLISECONDS]
		if !ok {
			errMsg := "toMillis: invalid TimeUnit"
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
		return duration * conversionFactor
	}
	conversionFactor, ok := timeUnitConversion[MILLISECONDS][unitName]
	if !ok {
		errMsg := "toMillis: invalid TimeUnit"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	return duration / conversionFactor
}

// toSeconds converts the given time to seconds
func toSeconds(params []interface{}) interface{} {
	unit := params[0].(*object.Object)
	duration := params[1].(int64)

	unitName := object.GoStringFromStringObject(unit)
	if isLarger(unitName, SECONDS) {
		conversionFactor, ok := timeUnitConversion[unitName][SECONDS]
		if !ok {
			errMsg := "toSeconds: invalid TimeUnit"
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
		return duration * conversionFactor
	}
	conversionFactor, ok := timeUnitConversion[SECONDS][unitName]
	if !ok {
		errMsg := "toSeconds: invalid TimeUnit"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	return duration / conversionFactor
}

// toMinutes converts the given time to minutes
func toMinutes(params []interface{}) interface{} {
	unit := params[0].(*object.Object)
	duration := params[1].(int64)

	unitName := object.GoStringFromStringObject(unit)
	if isLarger(unitName, MINUTES) {
		conversionFactor, ok := timeUnitConversion[unitName][MINUTES]
		if !ok {
			errMsg := "toMinutes: invalid TimeUnit"
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
		return duration * conversionFactor
	}
	conversionFactor, ok := timeUnitConversion[MINUTES][unitName]
	if !ok {
		errMsg := "toMinutes: invalid TimeUnit"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	return duration / conversionFactor
}

// toHours converts the given time to hours
func toHours(params []interface{}) interface{} {
	unit := params[0].(*object.Object)
	duration := params[1].(int64)

	unitName := object.GoStringFromStringObject(unit)
	if isLarger(unitName, HOURS) {
		conversionFactor, ok := timeUnitConversion[unitName][HOURS]
		if !ok {
			errMsg := "toHours: invalid TimeUnit"
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
		return duration * conversionFactor
	}
	conversionFactor, ok := timeUnitConversion[HOURS][unitName]
	if !ok {
		errMsg := "toHours: invalid TimeUnit"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	return duration / conversionFactor
}

// toDays converts the given time to days
func toDays(params []interface{}) interface{} {
	unit := params[0].(*object.Object)
	duration := params[1].(int64)

	unitName := object.GoStringFromStringObject(unit)
	if isLarger(unitName, DAYS) {
		conversionFactor, ok := timeUnitConversion[unitName][DAYS]
		if !ok {
			errMsg := "toDays: invalid TimeUnit"
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
		return duration * conversionFactor
	}
	conversionFactor, ok := timeUnitConversion[DAYS][unitName]
	if !ok {
		errMsg := "toDays: invalid TimeUnit"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	return duration / conversionFactor
}

var unitOrder = map[string]int{
	NANOSECONDS:  0,
	MICROSECONDS: 1,
	MILLISECONDS: 2,
	SECONDS:      3,
	MINUTES:      4,
	HOURS:        5,
	DAYS:         6,
}

func isLarger(unit1, unit2 string) bool {
	return unitOrder[unit1] >= unitOrder[unit2]
}
