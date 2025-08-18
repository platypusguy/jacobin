package gfunction

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
)

var atomicIntegerClassName = "java/util/concurrent/atomic/AtomicInteger"

func Load_Util_Concurrent_Atomic_AtomicInteger() {

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerClinit,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerInitVoid,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.<init>(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  atomicIntegerInitInt,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.accumulateAndGet(ILjava/util/function.IntBinaryOperator;)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.addAndGet(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  atomicIntegerAddAndGet,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.compareAndExchange(II)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.compareAndExchangeAcquire(II)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.compareAndExchangeRelease(II)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.compareAndSet(II)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  atomicIntegerCompareAndSet,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.decrementAndGet()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerDecrementAndGet,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.doubleValue()D"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerToFloat,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.floatValue()F"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerToFloat,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.get()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerGet,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.getAcquire()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerGet,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.getAndAccumulate(ILjava/util/function.IntBinaryOperator;)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.getAndAdd(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  atomicIntegerGetAndAdd,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.getAndDecrement()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerGetAndDecrement,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.getAndIncrement()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerGetAndIncrement,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.getAndSet(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  atomicIntegerGetAndSet,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.getAndUpdate(Ljava/util/function.IntUnaryOperator;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.getOpaque()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.getPlain()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerGet,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.incrementAndGet()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerIncrementAndGet,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.intValue()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerGet,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.lazySet(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  atomicIntegerInitInt, // same as <init>(I)V
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.longValue()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerGet,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.set(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  atomicIntegerSet, // same as <init>(I)V
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.setOpaque(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  atomicIntegerSet,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.setPlain(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  atomicIntegerSet,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.setRelease(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  atomicIntegerSet,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.toString()Ljava/base/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerToString,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.updateAndGet(Ljava/util/function.IntUnaryOperator;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.weakCompareAndSet(II)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.weakCompareAndSetAcquire(II)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.weakCompareAndSetPlain(II)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.weakCompareAndSetRelease(II)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

}

// "java/util/concurrent/atomic/AtomicInteger.<clinit>()V"
func atomicIntegerClinit(params []interface{}) interface{} {
	className := "java/util/concurrent/atomic/AtomicInteger"
	obj := object.MakeEmptyObjectWithClassName(&className)
	initialField := object.Field{Ftype: types.Int, Fvalue: int64(0)}
	obj.FieldTable["value"] = initialField
	return nil
}

// "java/util/concurrent/atomic/AtomicInteger.<init>()V"
func atomicIntegerInitVoid(params []interface{}) interface{} {
	initialField := object.Field{Ftype: types.Int, Fvalue: int64(0)}
	obj := params[0].(*object.Object)
	obj.FieldTable["value"] = initialField
	return nil
}

// "java/util/concurrent/atomic/AtomicInteger.<init>(I)V"
func atomicIntegerInitInt(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	initialValue := params[1].(int64)
	initialField := object.Field{Ftype: types.Int, Fvalue: initialValue}
	obj.FieldTable["value"] = initialField
	return nil
}

// "java/util/concurrent/atomic/AtomicInteger.Set(I)V"
func atomicIntegerSet(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	initialValue := params[1].(int64)
	initialField := object.Field{Ftype: types.Int, Fvalue: initialValue}
	obj.FieldTable["value"] = initialField
	return nil
}

// "java/util/concurrent/atomic/AtomicInteger.get()I"
func atomicIntegerGet(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	wint := obj.FieldTable["value"].Fvalue.(int64)
	return wint
}

// func atomicIntegerSet = atomicIntegerInitInt

// func atomicIntegerLazySet = atomicIntegerInitInt

func atomicIntegerGetAndSet(params []interface{}) interface{} {
	global := globals.GetGlobalRef()
	global.AtomicIntegerLock.Lock()
	defer global.AtomicIntegerLock.Unlock()
	obj := params[0].(*object.Object)
	oldValue := obj.FieldTable["value"].Fvalue.(int64)
	newValue := params[1].(int64)
	newField := object.Field{Ftype: types.Int, Fvalue: newValue}
	obj.FieldTable["value"] = newField
	return oldValue
}

func atomicIntegerCompareAndSet(params []interface{}) interface{} {
	global := globals.GetGlobalRef()
	global.AtomicIntegerLock.Lock()
	defer global.AtomicIntegerLock.Unlock()
	obj := params[0].(*object.Object)
	oldValue := obj.FieldTable["value"].Fvalue.(int64)
	expectedValue := params[1].(int64)
	if oldValue != expectedValue {
		return int64(0)
	}
	newValue := params[2].(int64)
	newField := object.Field{Ftype: types.Int, Fvalue: newValue}
	obj.FieldTable["value"] = newField
	return int64(1)
}

// func WeakCompareAndSet = trapDeprecated

func atomicIntegerGetAndIncrement(params []interface{}) interface{} {
	var fnParams []interface{}
	fnParams = append(fnParams, params[0])
	fnParams = append(fnParams, int64(1))
	ret := fnAtomicIntegerAdd(fnParams, false)
	return ret
}

func atomicIntegerGetAndDecrement(params []interface{}) interface{} {
	var fnParams []interface{}
	fnParams = append(fnParams, params[0])
	fnParams = append(fnParams, int64(-1))
	ret := fnAtomicIntegerAdd(fnParams, false)
	return ret
}

func atomicIntegerGetAndAdd(params []interface{}) interface{} {
	ret := fnAtomicIntegerAdd(params, false)
	return ret
}

func atomicIntegerIncrementAndGet(params []interface{}) interface{} {
	var fnParams []interface{}
	fnParams = append(fnParams, params[0])
	fnParams = append(fnParams, int64(1))
	ret := fnAtomicIntegerAdd(fnParams, true)
	return ret
}

func atomicIntegerDecrementAndGet(params []interface{}) interface{} {
	var fnParams []interface{}
	fnParams = append(fnParams, params[0])
	fnParams = append(fnParams, int64(-1))
	ret := fnAtomicIntegerAdd(fnParams, true)
	return ret
}

func atomicIntegerAddAndGet(params []interface{}) interface{} {
	ret := fnAtomicIntegerAdd(params, true)
	return ret
}

func atomicIntegerToString(params []interface{}) interface{} {
	global := globals.GetGlobalRef()
	global.AtomicIntegerLock.Lock()
	defer global.AtomicIntegerLock.Unlock()
	obj := params[0].(*object.Object)
	intValue := obj.FieldTable["value"].Fvalue.(int64)
	str := fmt.Sprintf("%d", intValue)
	return object.StringObjectFromGoString(str)
}

func atomicIntegerToFloat(params []interface{}) interface{} {
	global := globals.GetGlobalRef()
	global.AtomicIntegerLock.Lock()
	defer global.AtomicIntegerLock.Unlock()
	obj := params[0].(*object.Object)
	intValue := obj.FieldTable["value"].Fvalue.(int64)
	return float64(intValue)
}

// Internal function to add/subtract or increment/decrement
// and return either the old value or the new value, depending on newFlag.
func fnAtomicIntegerAdd(params []interface{}, newFlag bool) interface{} {

	// Validate the number of parameters.
	if len(params) != 2 {
		errMsg := fmt.Sprintf("fnAtomicIntegerAdd: Expected 2 parameters, observed %d", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Validate the first parameter (AtomicInteger object).
	obj, ok := params[0].(*object.Object)
	if !ok {
		var errMsg string
		if object.IsNull(params[0]) {
			errMsg = fmt.Sprintf("fnAtomicIntegerAdd: First parameter is null")
		} else {
			errMsg = fmt.Sprintf("fnAtomicIntegerAdd: First parameter is not an AtomicInteger object, observed %T", params[0])
		}
		return getGErrBlk(excNames.ClassCastException, errMsg)
	}
	// Guard against a typed nil object (null in Java terms)
	if obj == nil || object.IsNull(obj) {
		errMsg := "fnAtomicIntegerAdd: First parameter is null"
		return getGErrBlk(excNames.ClassCastException, errMsg)
	}

	// Validate the second parameter (int64 value to add)
	addend, ok := params[1].(int64)
	if !ok {
		errMsg := "fnAtomicIntegerAdd: Second parameter is not a valid int64"
		return getGErrBlk(excNames.ClassCastException, errMsg)
	}

	// Set up for lock and deferred unlock.
	global := globals.GetGlobalRef()
	global.AtomicIntegerLock.Lock()
	defer global.AtomicIntegerLock.Unlock()

	// Retrieve the current value field from the AtomicInteger object.
	valueField, exists := obj.FieldTable["value"]
	if !exists {
		errMsg := "fnAtomicIntegerAdd: AtomicInteger object does not have a 'value' field"
		return getGErrBlk(excNames.NoSuchFieldException, errMsg)
	}
	if valueField.Ftype != types.Int {
		errMsg := fmt.Sprintf("fnAtomicIntegerAdd: Expected 'value' field to be of type integer, observed %s", valueField.Ftype)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get the current value.
	formerValue, ok := valueField.Fvalue.(int64)
	if !ok {
		errMsg := "fnAtomicIntegerAdd: The 'value' field does not contain a valid int64"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Perform addition and update the AtomicInteger value field.
	newValue := formerValue + addend
	obj.FieldTable["value"] = object.Field{
		Ftype:  types.Int,
		Fvalue: newValue,
	}

	// Return the former value.
	if newFlag {
		return newValue
	}
	return formerValue
}
