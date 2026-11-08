/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"sync"
)

var classNameTemporalAdjusters = "java/time/temporal/TemporalAdjusters"
var classNameTemporalAdjusterImpl = "java/time/temporal/TemporalAdjuster"

var temporalAdjustersMutex sync.Mutex
var temporalAdjustersInitialized bool = false

var (
	adjusterFirstDayOfMonth     *object.Object
	adjusterLastDayOfMonth      *object.Object
	adjusterFirstDayOfNextMonth *object.Object
	adjusterFirstDayOfYear      *object.Object
	adjusterLastDayOfYear       *object.Object
	adjusterFirstDayOfNextYear  *object.Object
)

func makeAdjusterObject(adjusterType string) *object.Object {
	adj := object.MakeEmptyObjectWithClassName(&classNameTemporalAdjusterImpl)
	adj.FieldTable["adjusterType"] = object.Field{
		Ftype:  types.StringClassRef,
		Fvalue: object.StringObjectFromGoString(adjusterType),
	}
	return adj
}

func ensureTemporalAdjustersInited() {
	temporalAdjustersMutex.Lock()
	defer temporalAdjustersMutex.Unlock()
	if temporalAdjustersInitialized {
		return
	}

	adjusterFirstDayOfMonth = makeAdjusterObject("firstDayOfMonth")
	adjusterLastDayOfMonth = makeAdjusterObject("lastDayOfMonth")
	adjusterFirstDayOfNextMonth = makeAdjusterObject("firstDayOfNextMonth")
	adjusterFirstDayOfYear = makeAdjusterObject("firstDayOfYear")
	adjusterLastDayOfYear = makeAdjusterObject("lastDayOfYear")
	adjusterFirstDayOfNextYear = makeAdjusterObject("firstDayOfNextYear")

	temporalAdjustersInitialized = true
}

// Load_Time_Temporal_TemporalAdjusters registers all java/time/temporal/TemporalAdjusters methods.
func Load_Time_Temporal_TemporalAdjusters() {
	ghelpers.MethodSignatures["java/time/temporal/TemporalAdjusters.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalAdjustersClinit,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAdjusters.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalAdjustersInit,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAdjusters.dayOfWeekInMonth(ILjava/time/DayOfWeek;)Ljava/time/temporal/TemporalAdjuster;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  temporalAdjustersDayOfWeekInMonth,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAdjusters.firstDayOfMonth()Ljava/time/temporal/TemporalAdjuster;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalAdjustersFirstDayOfMonth,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAdjusters.firstDayOfNextMonth()Ljava/time/temporal/TemporalAdjuster;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalAdjustersFirstDayOfNextMonth,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAdjusters.firstDayOfNextYear()Ljava/time/temporal/TemporalAdjuster;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalAdjustersFirstDayOfNextYear,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAdjusters.firstDayOfYear()Ljava/time/temporal/TemporalAdjuster;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalAdjustersFirstDayOfYear,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAdjusters.firstInMonth(Ljava/time/DayOfWeek;)Ljava/time/temporal/TemporalAdjuster;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalAdjustersFirstInMonth,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAdjusters.lastDayOfMonth()Ljava/time/temporal/TemporalAdjuster;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalAdjustersLastDayOfMonth,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAdjusters.lastDayOfYear()Ljava/time/temporal/TemporalAdjuster;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  temporalAdjustersLastDayOfYear,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAdjusters.lastInMonth(Ljava/time/DayOfWeek;)Ljava/time/temporal/TemporalAdjuster;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalAdjustersLastInMonth,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAdjusters.next(Ljava/time/DayOfWeek;)Ljava/time/temporal/TemporalAdjuster;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalAdjustersNext,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAdjusters.nextOrSame(Ljava/time/DayOfWeek;)Ljava/time/temporal/TemporalAdjuster;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalAdjustersNextOrSame,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAdjusters.ofDateAdjuster(Ljava/util/function/UnaryOperator;)Ljava/time/temporal/TemporalAdjuster;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalAdjustersOfDateAdjuster,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAdjusters.previous(Ljava/time/DayOfWeek;)Ljava/time/temporal/TemporalAdjuster;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalAdjustersPrevious,
		}

	ghelpers.MethodSignatures["java/time/temporal/TemporalAdjusters.previousOrSame(Ljava/time/DayOfWeek;)Ljava/time/temporal/TemporalAdjuster;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  temporalAdjustersPreviousOrSame,
		}
}

func temporalAdjustersClinit(_ []any) any {
	temporalAdjustersMutex.Lock()
	temporalAdjustersInitialized = false
	temporalAdjustersMutex.Unlock()
	ensureTemporalAdjustersInited()
	return nil
}

func temporalAdjustersInit(_ []any) any {
	return nil
}

// dayOfWeekInMonth(int ordinal, DayOfWeek dayOfWeek)
func temporalAdjustersDayOfWeekInMonth(params []any) any {
	ensureTemporalAdjustersInited()
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TemporalAdjusters.dayOfWeekInMonth requires 2 arguments")
	}
	ordinal, okOrd := getInt64Val(params[0])
	if !okOrd {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TemporalAdjusters.dayOfWeekInMonth: invalid ordinal")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAdjusters.dayOfWeekInMonth: dayOfWeek is null")
	}
	dowObj := params[1].(*object.Object)
	dowVal, okVal := getDayOfWeekValueFromObject(dowObj)
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TemporalAdjusters.dayOfWeekInMonth: invalid dayOfWeek")
	}

	adj := makeAdjusterObject("dayOfWeekInMonth")
	adj.FieldTable["ordinal"] = object.Field{Ftype: types.Int, Fvalue: ordinal}
	adj.FieldTable["dowValue"] = object.Field{Ftype: types.Int, Fvalue: dowVal}
	adj.FieldTable["dayOfWeek"] = object.Field{Ftype: types.Ref, Fvalue: dowObj}
	return adj
}

// firstDayOfMonth()
func temporalAdjustersFirstDayOfMonth(_ []any) any {
	ensureTemporalAdjustersInited()
	return adjusterFirstDayOfMonth
}

// firstDayOfNextMonth()
func temporalAdjustersFirstDayOfNextMonth(_ []any) any {
	ensureTemporalAdjustersInited()
	return adjusterFirstDayOfNextMonth
}

// firstDayOfNextYear()
func temporalAdjustersFirstDayOfNextYear(_ []any) any {
	ensureTemporalAdjustersInited()
	return adjusterFirstDayOfNextYear
}

// firstDayOfYear()
func temporalAdjustersFirstDayOfYear(_ []any) any {
	ensureTemporalAdjustersInited()
	return adjusterFirstDayOfYear
}

// firstInMonth(DayOfWeek dayOfWeek)
func temporalAdjustersFirstInMonth(params []any) any {
	ensureTemporalAdjustersInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAdjusters.firstInMonth: dayOfWeek is null")
	}
	return temporalAdjustersDayOfWeekInMonth([]any{int64(1), params[0]})
}

// lastDayOfMonth()
func temporalAdjustersLastDayOfMonth(_ []any) any {
	ensureTemporalAdjustersInited()
	return adjusterLastDayOfMonth
}

// lastDayOfYear()
func temporalAdjustersLastDayOfYear(_ []any) any {
	ensureTemporalAdjustersInited()
	return adjusterLastDayOfYear
}

// lastInMonth(DayOfWeek dayOfWeek)
func temporalAdjustersLastInMonth(params []any) any {
	ensureTemporalAdjustersInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAdjusters.lastInMonth: dayOfWeek is null")
	}
	return temporalAdjustersDayOfWeekInMonth([]any{int64(-1), params[0]})
}

// next(DayOfWeek dayOfWeek)
func temporalAdjustersNext(params []any) any {
	ensureTemporalAdjustersInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAdjusters.next: dayOfWeek is null")
	}
	dowObj := params[0].(*object.Object)
	dowVal, okVal := getDayOfWeekValueFromObject(dowObj)
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TemporalAdjusters.next: invalid dayOfWeek")
	}

	adj := makeAdjusterObject("next")
	adj.FieldTable["dowValue"] = object.Field{Ftype: types.Int, Fvalue: dowVal}
	adj.FieldTable["dayOfWeek"] = object.Field{Ftype: types.Ref, Fvalue: dowObj}
	return adj
}

// nextOrSame(DayOfWeek dayOfWeek)
func temporalAdjustersNextOrSame(params []any) any {
	ensureTemporalAdjustersInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAdjusters.nextOrSame: dayOfWeek is null")
	}
	dowObj := params[0].(*object.Object)
	dowVal, okVal := getDayOfWeekValueFromObject(dowObj)
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TemporalAdjusters.nextOrSame: invalid dayOfWeek")
	}

	adj := makeAdjusterObject("nextOrSame")
	adj.FieldTable["dowValue"] = object.Field{Ftype: types.Int, Fvalue: dowVal}
	adj.FieldTable["dayOfWeek"] = object.Field{Ftype: types.Ref, Fvalue: dowObj}
	return adj
}

// ofDateAdjuster(UnaryOperator dateBasedAdjuster)
func temporalAdjustersOfDateAdjuster(params []any) any {
	ensureTemporalAdjustersInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAdjusters.ofDateAdjuster: dateBasedAdjuster is null")
	}
	fnObj := params[0].(*object.Object)

	adj := makeAdjusterObject("ofDateAdjuster")
	adj.FieldTable["dateBasedAdjuster"] = object.Field{Ftype: types.Ref, Fvalue: fnObj}
	return adj
}

// previous(DayOfWeek dayOfWeek)
func temporalAdjustersPrevious(params []any) any {
	ensureTemporalAdjustersInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAdjusters.previous: dayOfWeek is null")
	}
	dowObj := params[0].(*object.Object)
	dowVal, okVal := getDayOfWeekValueFromObject(dowObj)
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TemporalAdjusters.previous: invalid dayOfWeek")
	}

	adj := makeAdjusterObject("previous")
	adj.FieldTable["dowValue"] = object.Field{Ftype: types.Int, Fvalue: dowVal}
	adj.FieldTable["dayOfWeek"] = object.Field{Ftype: types.Ref, Fvalue: dowObj}
	return adj
}

// previousOrSame(DayOfWeek dayOfWeek)
func temporalAdjustersPreviousOrSame(params []any) any {
	ensureTemporalAdjustersInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TemporalAdjusters.previousOrSame: dayOfWeek is null")
	}
	dowObj := params[0].(*object.Object)
	dowVal, okVal := getDayOfWeekValueFromObject(dowObj)
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TemporalAdjusters.previousOrSame: invalid dayOfWeek")
	}

	adj := makeAdjusterObject("previousOrSame")
	adj.FieldTable["dowValue"] = object.Field{Ftype: types.Int, Fvalue: dowVal}
	adj.FieldTable["dayOfWeek"] = object.Field{Ftype: types.Ref, Fvalue: dowObj}
	return adj
}
