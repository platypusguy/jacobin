/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
)

// Load_Time_Temporal_TemporalQuery loads method signatures for java/time/temporal/TemporalQuery.
func Load_Time_Temporal_TemporalQuery() {
	ghelpers.MethodSignatures["java/time/temporal/TemporalQuery.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalQuery.queryFrom(Ljava/time/temporal/TemporalAccessor;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalQueryQueryFrom,
		}
}

// temporalQueryQueryFrom queries the specified temporal object.
// R queryFrom(TemporalAccessor temporal)
func temporalQueryQueryFrom(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalQuery.queryFrom: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalQuery.queryFrom: temporal is null")
	}
	selfObject := params[0].(*object.Object)
	temporalObject := params[1].(*object.Object)

	return executeTemporalQuery(selfObject, temporalObject)
}
