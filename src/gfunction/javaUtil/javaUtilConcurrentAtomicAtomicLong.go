/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"runtime"
)

var atomicLongClassName = "java/util/concurrent/atomic/AtomicLong"

func Load_Util_Concurrent_Atomic_Atomic_Long() {

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  atomicLongClinit,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  atomicLongInitVoid,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.<init>(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  atomicLongInitLong,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.accumulateAndGet(JLjava/util/function.LongBinaryOperator;)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.addAndGet(J)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  atomicLongAddAndGet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.compareAndExchange(JJ)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.compareAndExchangeAcquire(JJ)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.compareAndExchangeRelease(JJ)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.compareAndSet(JJ)Z"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  atomicLongCompareAndSet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.decrementAndGet()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  atomicLongDecrementAndGet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.doubleValue()D"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  atomicLongToFloat,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.floatValue()F"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  atomicLongToFloat,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.get()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  atomicLongGet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.getAcquire()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  atomicLongGet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.getAndAccumulate(JLjava/util/function.LongBinaryOperator;)J"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.getAndAdd(J)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  atomicLongGetAndAdd,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.getAndDecrement()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  atomicLongGetAndDecrement,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.getAndIncrement()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  atomicLongGetAndIncrement,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.getAndSet(J)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  atomicLongGetAndSet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.getAndUpdate(Ljava/util/function.LongUnaryOperator;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.getOpaque()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.getPlain()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  atomicLongGet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.incrementAndGet()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  atomicLongIncrementAndGet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.intValue()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  atomicLongGet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.lazySet(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  atomicLongInitLong,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.longValue()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  atomicLongGet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.set(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  atomicLongSet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.setOpaque(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  atomicLongSet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.setPlain(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  atomicLongSet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.setRelease(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  atomicLongSet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  atomicLongToString,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.updateAndGet(Ljava/util/function.LongUnaryOperator;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.weakCompareAndSet(JJ)Z"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapDeprecated,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.weakCompareAndSetAcquire(JJ)Z"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.weakCompareAndSetPlain(JJ)Z"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.weakCompareAndSetRelease(JJ)Z"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/atomic/AtomicLong.VMSupportsCS8()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  atomicLongVMSupportsCS8,
		}

}

func atomicLongClinit([]interface{}) interface{} {
	className := "java/util/concurrent/atomic/AtomicLong"
	obj := object.MakeEmptyObjectWithClassName(&className)
	initialField := object.Field{Ftype: types.Long, Fvalue: int64(0)}
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable["value"] = initialField
	return nil
}

func atomicLongInitVoid(params []interface{}) interface{} {
	initialField := object.Field{Ftype: types.Long, Fvalue: int64(0)}
	obj := params[0].(*object.Object)
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable["value"] = initialField
	return nil
}

func atomicLongInitLong(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	initialValue := params[1].(int64)
	initialField := object.Field{Ftype: types.Long, Fvalue: initialValue}
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable["value"] = initialField
	return nil
}

func atomicLongSet(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	initialValue := params[1].(int64)
	initialField := object.Field{Ftype: types.Long, Fvalue: initialValue}
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable["value"] = initialField
	return nil
}

func atomicLongGet(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	wlong := obj.FieldTable["value"].Fvalue.(int64)
	return wlong
}

func atomicLongGetAndSet(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	oldValue := obj.FieldTable["value"].Fvalue.(int64)
	newValue := params[1].(int64)
	newField := object.Field{Ftype: types.Long, Fvalue: newValue}
	obj.FieldTable["value"] = newField
	return oldValue
}

func atomicLongCompareAndSet(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	oldValue := obj.FieldTable["value"].Fvalue.(int64)
	expectedValue := params[1].(int64)
	if oldValue != expectedValue {
		return types.JavaBoolFalse
	}
	newValue := params[2].(int64)
	newField := object.Field{Ftype: types.Long, Fvalue: newValue}
	obj.FieldTable["value"] = newField
	return types.JavaBoolTrue
}

func atomicLongGetAndIncrement(params []interface{}) interface{} {
	var fnParams []interface{}
	fnParams = append(fnParams, params[0])
	fnParams = append(fnParams, int64(1))
	ret := fnAtomicLongAdd(fnParams, false)
	return ret
}

func atomicLongGetAndDecrement(params []interface{}) interface{} {
	var fnParams []interface{}
	fnParams = append(fnParams, params[0])
	fnParams = append(fnParams, int64(-1))
	ret := fnAtomicLongAdd(fnParams, false)
	return ret
}

func atomicLongGetAndAdd(params []interface{}) interface{} {
	ret := fnAtomicLongAdd(params, false)
	return ret
}

func atomicLongIncrementAndGet(params []interface{}) interface{} {
	var fnParams []interface{}
	fnParams = append(fnParams, params[0])
	fnParams = append(fnParams, int64(1))
	ret := fnAtomicLongAdd(fnParams, true)
	return ret
}

func atomicLongDecrementAndGet(params []interface{}) interface{} {
	var fnParams []interface{}
	fnParams = append(fnParams, params[0])
	fnParams = append(fnParams, int64(-1))
	ret := fnAtomicLongAdd(fnParams, true)
	return ret
}

func atomicLongAddAndGet(params []interface{}) interface{} {
	ret := fnAtomicLongAdd(params, true)
	return ret
}

func atomicLongToString(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	longValue := obj.FieldTable["value"].Fvalue.(int64)
	str := fmt.Sprintf("%d", longValue)
	return object.StringObjectFromGoString(str)
}

func atomicLongToFloat(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	longValue := obj.FieldTable["value"].Fvalue.(int64)
	return float64(longValue)
}

func fnAtomicLongAdd(params []interface{}, newFlag bool) interface{} {
	if len(params) != 2 {
		errMsg := fmt.Sprintf("fnAtomicLongAdd: Expected 2 parameters, observed %d", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "fnAtomicLongAdd: First parameter is not a valid object")
	}

	addend, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "fnAtomicLongAdd: Second parameter is not a valid int64")
	}

	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()

	valueField, exists := obj.FieldTable["value"]
	if !exists {
		return ghelpers.GetGErrBlk(excNames.NoSuchFieldException, "fnAtomicLongAdd: AtomicLong object does not have a 'value' field")
	}
	if valueField.Ftype != types.Long {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "fnAtomicLongAdd: Expected 'value' field to be of type long")
	}

	formerValue := valueField.Fvalue.(int64)
	newValue := formerValue + addend
	obj.FieldTable["value"] = object.Field{
		Ftype:  types.Long,
		Fvalue: newValue,
	}

	if newFlag {
		return newValue
	}
	return formerValue
}

func atomicLongVMSupportsCS8([]interface{}) interface{} {
	arch := runtime.GOARCH
	supportedArchitectures := map[string]bool{
		"amd64":    true,
		"arm64":    true,
		"ppc64":    true,
		"ppc64le":  true,
		"s390x":    true,
		"sparc64":  true,
		"mips64":   true,
		"mips64le": true,
	}

	if supported, ok := supportedArchitectures[arch]; ok {
		return object.JavaBooleanFromGoBoolean(supported)
	}
	return object.JavaBooleanFromGoBoolean(false)
}
