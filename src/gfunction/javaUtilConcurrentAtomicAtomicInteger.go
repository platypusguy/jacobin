package gfunction

import (
	"jacobin/globals"
	"jacobin/object"
	"jacobin/types"
)

var atomicIntegerClassName = "java/util/concurrent/atomic/AtomicInteger"

func Load_Util_Concurrent_Atomic_AtomicInteger() map[string]GMeth {

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.<clinit>V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
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

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.get()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerGet,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.set(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  atomicIntegerInitInt,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.lazySet(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  atomicIntegerInitInt,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.getAndSet(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  atomicIntegerInitInt,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.compareAndSet(II)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  atomicIntegerInitInt,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.weakCompareAndSet(II)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.getAndIncrement()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerGetAndIncrement,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.getAndDecrement()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerGetAndDecrement,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.getAndAdd(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  atomicIntegerGetAndDecrement,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.incrementAndGet()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerGetAndIncrement,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.decrementAndGet()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  atomicIntegerGetAndDecrement,
		}

	MethodSignatures["java/util/concurrent/atomic/AtomicInteger.addAndGet(I)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  atomicIntegerGetAndDecrement,
		}

	return MethodSignatures
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
	initialValue := params[1].(int64)
	initialField := object.Field{Ftype: types.Int, Fvalue: initialValue}
	obj := params[0].(*object.Object)
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
	global.AtomicIntegerLock.Lock() // <-------------------
	obj := params[0].(*object.Object)
	oldValue := obj.FieldTable["value"].Fvalue.(int64)
	newValue := params[1].(int64)
	newField := object.Field{Ftype: types.Int, Fvalue: newValue}
	obj.FieldTable["value"] = newField
	global.AtomicIntegerLock.Unlock() // <-------------------
	return oldValue
}

func atomicIntegerCompareAndSet(params []interface{}) interface{} {
	global := globals.GetGlobalRef()
	global.AtomicIntegerLock.Lock() // <-------------------
	obj := params[0].(*object.Object)
	oldValue := obj.FieldTable["value"].Fvalue.(int64)
	expectedValue := params[1].(int64)
	if oldValue != expectedValue {
		global.AtomicIntegerLock.Unlock() // <-------------------
		return int64(0)
	}
	newValue := params[2].(int64)
	newField := object.Field{Ftype: types.Int, Fvalue: newValue}
	obj.FieldTable["value"] = newField
	global.AtomicIntegerLock.Unlock() // <-------------------
	return int64(1)
}

// func WeakCompareAndSet = trapDeprecated

func atomicIntegerGetAndIncrement(params []interface{}) interface{} {
	global := globals.GetGlobalRef()
	global.AtomicIntegerLock.Lock() // <-------------------
	obj := params[0].(*object.Object)
	oldValue := obj.FieldTable["value"].Fvalue.(int64)
	newValue := oldValue + 1
	newField := object.Field{Ftype: types.Int, Fvalue: newValue}
	obj.FieldTable["value"] = newField
	global.AtomicIntegerLock.Unlock() // <-------------------
	return oldValue                   // previous value
}

func atomicIntegerGetAndDecrement(params []interface{}) interface{} {
	global := globals.GetGlobalRef()
	global.AtomicIntegerLock.Lock() // <-------------------
	obj := params[0].(*object.Object)
	oldValue := obj.FieldTable["value"].Fvalue.(int64)
	newValue := oldValue - 1
	newField := object.Field{Ftype: types.Int, Fvalue: newValue}
	obj.FieldTable["value"] = newField
	global.AtomicIntegerLock.Unlock() // <-------------------
	return oldValue                   // previous value
}

func atomicIntegerGetAndAdd(params []interface{}) interface{} {
	global := globals.GetGlobalRef()
	global.AtomicIntegerLock.Lock() // <-------------------
	obj := params[0].(*object.Object)
	delta := params[1].(int64)
	oldValue := obj.FieldTable["value"].Fvalue.(int64)
	newValue := oldValue + delta
	newField := object.Field{Ftype: types.Int, Fvalue: newValue}
	obj.FieldTable["value"] = newField
	global.AtomicIntegerLock.Unlock() // <-------------------
	return oldValue                   // previous value
}

func atomicIntegerIncrementAndGet(params []interface{}) interface{} {
	global := globals.GetGlobalRef()
	global.AtomicIntegerLock.Lock() // <-------------------
	obj := params[0].(*object.Object)
	oldValue := obj.FieldTable["value"].Fvalue.(int64)
	newValue := oldValue + 1
	newField := object.Field{Ftype: types.Int, Fvalue: newValue}
	obj.FieldTable["value"] = newField
	global.AtomicIntegerLock.Unlock() // <-------------------
	return newValue                   // previous value
}

func atomicIntegerDecrementAndGet(params []interface{}) interface{} {
	global := globals.GetGlobalRef()
	global.AtomicIntegerLock.Lock() // <-------------------
	obj := params[0].(*object.Object)
	oldValue := obj.FieldTable["value"].Fvalue.(int64)
	newValue := oldValue - 1
	newField := object.Field{Ftype: types.Int, Fvalue: newValue}
	obj.FieldTable["value"] = newField
	global.AtomicIntegerLock.Unlock() // <-------------------
	return newValue                   // previous value
}

func atomicIntegerAddAndGet(params []interface{}) interface{} {
	global := globals.GetGlobalRef()
	global.AtomicIntegerLock.Lock() // <-------------------
	obj := params[0].(*object.Object)
	delta := params[1].(int64)
	oldValue := obj.FieldTable["value"].Fvalue.(int64)
	newValue := oldValue + delta
	newField := object.Field{Ftype: types.Int, Fvalue: newValue}
	obj.FieldTable["value"] = newField
	global.AtomicIntegerLock.Unlock() // <-------------------
	return newValue                   // previous value
}
