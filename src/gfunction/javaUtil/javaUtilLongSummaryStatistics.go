/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaUtil

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"math"
)

func Load_Util_LongSummaryStatistics() {

	ghelpers.MethodSignatures["java/util/LongSummaryStatistics.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/LongSummaryStatistics.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  longSummaryStatisticsInit,
		}

	ghelpers.MethodSignatures["java/util/LongSummaryStatistics.<init>(JJJJ)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  longSummaryStatisticsInitJJJJ,
		}

	ghelpers.MethodSignatures["java/util/LongSummaryStatistics.accept(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  longSummaryStatisticsAccept,
		}

	ghelpers.MethodSignatures["java/util/LongSummaryStatistics.combine(Ljava/util/LongSummaryStatistics;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  longSummaryStatisticsCombine,
		}

	ghelpers.MethodSignatures["java/util/LongSummaryStatistics.getCount()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  longSummaryStatisticsGetCount,
		}

	ghelpers.MethodSignatures["java/util/LongSummaryStatistics.getSum()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  longSummaryStatisticsGetSum,
		}

	ghelpers.MethodSignatures["java/util/LongSummaryStatistics.getMin()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  longSummaryStatisticsGetMin,
		}

	ghelpers.MethodSignatures["java/util/LongSummaryStatistics.getMax()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  longSummaryStatisticsGetMax,
		}

	ghelpers.MethodSignatures["java/util/LongSummaryStatistics.getAverage()D"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  longSummaryStatisticsGetAverage,
		}

	ghelpers.MethodSignatures["java/util/LongSummaryStatistics.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  longSummaryStatisticsToString,
		}

}

// Initialize a LongSummaryStatistics object with count=0, sum=0, min=Long.MAX_VALUE, max=Long.MIN_VALUE,
// matching the default state of java.util.LongSummaryStatistics.
func longSummaryStatisticsInit(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable["count"] = object.Field{Ftype: types.Long, Fvalue: int64(0)}
	obj.FieldTable["sum"] = object.Field{Ftype: types.Long, Fvalue: int64(0)}
	obj.FieldTable["min"] = object.Field{Ftype: types.Long, Fvalue: int64(math.MaxInt64)}
	obj.FieldTable["max"] = object.Field{Ftype: types.Long, Fvalue: int64(math.MinInt64)}
	return nil
}

// Initialize a LongSummaryStatistics object from explicit count, min, max, sum values.
// Matches java.util.LongSummaryStatistics(long count, long min, long max, long sum).
// If count is zero, the remaining arguments are ignored and an empty instance is constructed.
// If the arguments are inconsistent (count < 0 or min > max), an IllegalArgumentException is thrown.
func longSummaryStatisticsInitJJJJ(params []interface{}) interface{} {
	if len(params) < 5 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "longSummaryStatisticsInitJJJJ: too few arguments")
	}

	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "longSummaryStatisticsInitJJJJ: receiver is not a valid object")
	}

	count, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "longSummaryStatisticsInitJJJJ: count argument is not a valid long")
	}
	min, ok := params[2].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "longSummaryStatisticsInitJJJJ: min argument is not a valid long")
	}
	max, ok := params[3].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "longSummaryStatisticsInitJJJJ: max argument is not a valid long")
	}
	sum, ok := params[4].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "longSummaryStatisticsInitJJJJ: sum argument is not a valid long")
	}

	if count < 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Negative count value")
	} else if count > 0 {
		if min > max {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Minimum greater than maximum")
		}
	} else {
		// Zero count should only have the default min/max values.
		min = math.MaxInt64
		max = math.MinInt64
		sum = 0
	}

	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable["count"] = object.Field{Ftype: types.Long, Fvalue: count}
	obj.FieldTable["sum"] = object.Field{Ftype: types.Long, Fvalue: sum}
	obj.FieldTable["min"] = object.Field{Ftype: types.Long, Fvalue: min}
	obj.FieldTable["max"] = object.Field{Ftype: types.Long, Fvalue: max}

	return nil
}

// Record a new long value into the statistics.
func longSummaryStatisticsAccept(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "longSummaryStatisticsAccept: too few arguments")
	}

	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "longSummaryStatisticsAccept: receiver is not a valid object")
	}

	value, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "longSummaryStatisticsAccept: argument is not a valid long")
	}

	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()

	count := obj.FieldTable["count"].Fvalue.(int64)
	sum := obj.FieldTable["sum"].Fvalue.(int64)
	min := obj.FieldTable["min"].Fvalue.(int64)
	max := obj.FieldTable["max"].Fvalue.(int64)

	count++
	sum += value
	if value < min {
		min = value
	}
	if value > max {
		max = value
	}

	obj.FieldTable["count"] = object.Field{Ftype: types.Long, Fvalue: count}
	obj.FieldTable["sum"] = object.Field{Ftype: types.Long, Fvalue: sum}
	obj.FieldTable["min"] = object.Field{Ftype: types.Long, Fvalue: min}
	obj.FieldTable["max"] = object.Field{Ftype: types.Long, Fvalue: max}

	return nil
}

// Combine the state of another LongSummaryStatistics into this one.
func longSummaryStatisticsCombine(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "longSummaryStatisticsCombine: too few arguments")
	}

	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "longSummaryStatisticsCombine: receiver is not a valid object")
	}

	other, ok := params[1].(*object.Object)
	if !ok || other == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "longSummaryStatisticsCombine: argument is not a valid object")
	}

	other.ThMutex.RLock()
	otherCount := other.FieldTable["count"].Fvalue.(int64)
	otherSum := other.FieldTable["sum"].Fvalue.(int64)
	otherMin := other.FieldTable["min"].Fvalue.(int64)
	otherMax := other.FieldTable["max"].Fvalue.(int64)
	other.ThMutex.RUnlock()

	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()

	count := obj.FieldTable["count"].Fvalue.(int64)
	sum := obj.FieldTable["sum"].Fvalue.(int64)
	min := obj.FieldTable["min"].Fvalue.(int64)
	max := obj.FieldTable["max"].Fvalue.(int64)

	count += otherCount
	sum += otherSum
	if otherMin < min {
		min = otherMin
	}
	if otherMax > max {
		max = otherMax
	}

	obj.FieldTable["count"] = object.Field{Ftype: types.Long, Fvalue: count}
	obj.FieldTable["sum"] = object.Field{Ftype: types.Long, Fvalue: sum}
	obj.FieldTable["min"] = object.Field{Ftype: types.Long, Fvalue: min}
	obj.FieldTable["max"] = object.Field{Ftype: types.Long, Fvalue: max}

	return nil
}

func longSummaryStatisticsGetCount(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	return obj.FieldTable["count"].Fvalue.(int64)
}

func longSummaryStatisticsGetSum(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	return obj.FieldTable["sum"].Fvalue.(int64)
}

func longSummaryStatisticsGetMin(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	return obj.FieldTable["min"].Fvalue.(int64)
}

func longSummaryStatisticsGetMax(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	return obj.FieldTable["max"].Fvalue.(int64)
}

func longSummaryStatisticsGetAverage(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	count := obj.FieldTable["count"].Fvalue.(int64)
	sum := obj.FieldTable["sum"].Fvalue.(int64)
	if count == 0 {
		return float64(0)
	}
	return float64(sum) / float64(count)
}

func longSummaryStatisticsToString(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	obj.ThMutex.RLock()
	count := obj.FieldTable["count"].Fvalue.(int64)
	sum := obj.FieldTable["sum"].Fvalue.(int64)
	min := obj.FieldTable["min"].Fvalue.(int64)
	max := obj.FieldTable["max"].Fvalue.(int64)
	obj.ThMutex.RUnlock()

	var average float64
	if count == 0 {
		average = 0
	} else {
		average = float64(sum) / float64(count)
	}

	str := fmt.Sprintf("LongSummaryStatistics{count=%d, sum=%d, min=%d, average=%f, max=%d}",
		count, sum, min, average, max)
	return object.StringObjectFromGoString(str)
}
