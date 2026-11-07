/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"math"
	"strings"
)

var classNameConcurrentHashMap = "java/util/concurrent/ConcurrentHashMap"
var classNameKeySetView = "java/util/concurrent/ConcurrentHashMap$KeySetView"

func Load_Util_Concurrent_Hash_Map() {

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	// -----------------------
	// Instantiation functions
	// -----------------------

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  concurrentHashMapInit,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.<init>(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  concurrentHashMapInitCapacity,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.<init>(IF)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  concurrentHashMapInitCapacityLoadFactor,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.<init>(IFI)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  concurrentHashMapInitCapacityLoadFactorConcurrencyLevel,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.<init>(Ljava/util/Map;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  concurrentHashMapInitMap,
		}

	// ----------------
	// Static factories
	// ----------------

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.newKeySet()Ljava/util/concurrent/ConcurrentHashMap$KeySetView;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  concurrentHashMapNewKeySet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.newKeySet(I)Ljava/util/concurrent/ConcurrentHashMap$KeySetView;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  concurrentHashMapNewKeySetCapacity,
		}

	// ----------------
	// Member functions
	// ----------------

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.clear()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  concurrentHashMapClear,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.compute(Ljava/lang/Object;Ljava/util/function/BiFunction;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.computeIfAbsent(Ljava/lang/Object;Ljava/util/function/Function;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.computeIfPresent(Ljava/lang/Object;Ljava/util/function/BiFunction;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.contains(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  concurrentHashMapContains,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.containsKey(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  concurrentHashMapContainsKey,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.containsValue(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  concurrentHashMapContainsValue,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.elements()Ljava/util/Enumeration;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  concurrentHashMapElements,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.entrySet()Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  concurrentHashMapEntrySet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  concurrentHashMapEquals,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.forEach(Ljava/util/function/BiConsumer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.forEach(JLjava/util/function/BiConsumer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.forEach(JLjava/util/function/BiFunction;Ljava/util/function/Consumer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.forEachEntry(JLjava/util/function/Consumer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.forEachEntry(JLjava/util/function/Function;Ljava/util/function/Consumer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.forEachKey(JLjava/util/function/Consumer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.forEachKey(JLjava/util/function/Function;Ljava/util/function/Consumer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.forEachValue(JLjava/util/function/Consumer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.forEachValue(JLjava/util/function/Function;Ljava/util/function/Consumer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.get(Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  concurrentHashMapGet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.getOrDefault(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  concurrentHashMapGetOrDefault,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  concurrentHashMapHashCode,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.isEmpty()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  concurrentHashMapIsEmpty,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.keySet()Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  concurrentHashMapKeySet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.keySet()Ljava/util/concurrent/ConcurrentHashMap$KeySetView;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  concurrentHashMapKeySet,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.keySet(Ljava/lang/Object;)Ljava/util/concurrent/ConcurrentHashMap$KeySetView;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  concurrentHashMapKeySetWithDefaultValue,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.keys()Ljava/util/Enumeration;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  concurrentHashMapKeys,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.mappingCount()J"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  concurrentHashMapMappingCount,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.merge(Ljava/lang/Object;Ljava/lang/Object;Ljava/util/function/BiFunction;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.put(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  concurrentHashMapPut,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.putAll(Ljava/util/Map;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  concurrentHashMapPutAll,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.putIfAbsent(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  concurrentHashMapPutIfAbsent,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduce(JLjava/util/function/BiFunction;Ljava/util/function/BiFunction;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduceEntries(JLjava/util/function/BiFunction;)Ljava/util/Map$Entry;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduceEntries(JLjava/util/function/Function;Ljava/util/function/BiFunction;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduceEntriesToDouble(JLjava/util/function/ToDoubleFunction;DLjava/util/function/DoubleBinaryOperator;)D"] =
		ghelpers.GMeth{
			ParamSlots: 6,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduceEntriesToInt(JLjava/util/function/ToIntFunction;ILjava/util/function/IntBinaryOperator;)I"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduceEntriesToLong(JLjava/util/function/ToLongFunction;JLjava/util/function/LongBinaryOperator;)J"] =
		ghelpers.GMeth{
			ParamSlots: 6,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduceKeys(JLjava/util/function/BiFunction;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduceKeys(JLjava/util/function/Function;Ljava/util/function/BiFunction;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduceKeysToDouble(JLjava/util/function/ToDoubleFunction;DLjava/util/function/DoubleBinaryOperator;)D"] =
		ghelpers.GMeth{
			ParamSlots: 6,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduceKeysToInt(JLjava/util/function/ToIntFunction;ILjava/util/function/IntBinaryOperator;)I"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduceKeysToLong(JLjava/util/function/ToLongFunction;JLjava/util/function/LongBinaryOperator;)J"] =
		ghelpers.GMeth{
			ParamSlots: 6,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduceToDouble(JLjava/util/function/ToDoubleBiFunction;DLjava/util/function/DoubleBinaryOperator;)D"] =
		ghelpers.GMeth{
			ParamSlots: 6,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduceToInt(JLjava/util/function/ToIntBiFunction;ILjava/util/function/IntBinaryOperator;)I"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduceToLong(JLjava/util/function/ToLongBiFunction;JLjava/util/function/LongBinaryOperator;)J"] =
		ghelpers.GMeth{
			ParamSlots: 6,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduceValues(JLjava/util/function/BiFunction;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduceValues(JLjava/util/function/Function;Ljava/util/function/BiFunction;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduceValuesToDouble(JLjava/util/function/ToDoubleFunction;DLjava/util/function/DoubleBinaryOperator;)D"] =
		ghelpers.GMeth{
			ParamSlots: 6,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduceValuesToInt(JLjava/util/function/ToIntFunction;ILjava/util/function/IntBinaryOperator;)I"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.reduceValuesToLong(JLjava/util/function/ToLongFunction;JLjava/util/function/LongBinaryOperator;)J"] =
		ghelpers.GMeth{
			ParamSlots: 6,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.remove(Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  concurrentHashMapRemove,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.remove(Ljava/lang/Object;Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  concurrentHashMapRemoveEntry,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.replace(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  concurrentHashMapReplace,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.replace(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  concurrentHashMapReplaceOldNew,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.replaceAll(Ljava/util/function/BiFunction;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.search(JLjava/util/function/BiFunction;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.searchEntries(JLjava/util/function/Function;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.searchKeys(JLjava/util/function/Function;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.searchValues(JLjava/util/function/Function;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.size()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  concurrentHashMapSize,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  concurrentHashMapToString,
		}

	ghelpers.MethodSignatures["java/util/concurrent/ConcurrentHashMap.values()Ljava/util/Collection;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  concurrentHashMapValues,
		}

}

func _getConcurrentHashMap(params []interface{}, funcName string) (*object.Object, types.DefHashMap, *ghelpers.GErrBlk) {
	if len(params) < 1 {
		return nil, nil, ghelpers.GetGErrBlk(excNames.IllegalArgumentException, funcName+": missing 'this' parameter")
	}
	if params[0] == nil || params[0] == object.Null {
		return nil, nil, ghelpers.GetGErrBlk(excNames.NullPointerException, funcName+": 'this' is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok {
		return nil, nil, ghelpers.GetGErrBlk(excNames.ClassCastException, funcName+": 'this' is not an object")
	}
	if object.IsNull(this) {
		return nil, nil, ghelpers.GetGErrBlk(excNames.NullPointerException, funcName+": 'this' is null")
	}

	actualKlassName := *stringPool.GetStringPointer(this.KlassName)
	if actualKlassName != classNameConcurrentHashMap {
		errMsg := fmt.Sprintf("%s: expected %s, got %s", funcName, classNameConcurrentHashMap, actualKlassName)
		return nil, nil, ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	this.ThMutex.RLock()
	fld, ok := this.FieldTable[fieldNameMap]
	this.ThMutex.RUnlock()
	if !ok {
		return nil, nil, ghelpers.GetGErrBlk(excNames.IllegalArgumentException, funcName+": missing 'map' field")
	}

	hm, ok := fld.Fvalue.(types.DefHashMap)
	if !ok {
		return nil, nil, ghelpers.GetGErrBlk(excNames.IllegalArgumentException, funcName+": 'map' field is invalid")
	}

	return this, hm, nil
}

func _getKeyConcurrent(param interface{}, funcName string) (interface{}, *ghelpers.GErrBlk) {
	if param == nil || param == object.Null {
		return nil, ghelpers.GetGErrBlk(excNames.NullPointerException, funcName+": key cannot be null")
	}
	keyObj, ok := param.(*object.Object)
	if !ok {
		return param, nil
	}
	if object.IsNull(keyObj) {
		return nil, ghelpers.GetGErrBlk(excNames.NullPointerException, funcName+": key cannot be null")
	}
	if object.IsStringObject(keyObj) {
		return object.GoStringFromStringObject(keyObj), nil
	}
	if fld, ok := keyObj.FieldTable[fieldNameValue]; ok {
		switch v := fld.Fvalue.(type) {
		case int64, float64, int32, int16, int8, uint8, uint16, bool, string:
			return v, nil
		}
	}
	return keyObj, nil
}

func _getValueConcurrent(param interface{}, funcName string) (interface{}, *ghelpers.GErrBlk) {
	if param == nil || param == object.Null {
		return nil, ghelpers.GetGErrBlk(excNames.NullPointerException, funcName+": value cannot be null")
	}
	valObj, ok := param.(*object.Object)
	if ok && (valObj == nil || object.IsNull(valObj)) {
		return nil, ghelpers.GetGErrBlk(excNames.NullPointerException, funcName+": value cannot be null")
	}
	return param, nil
}

func concurrentHashMapInit(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ConcurrentHashMap.<init>: missing 'this' parameter")
	}
	if params[0] == nil || params[0] == object.Null {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ConcurrentHashMap.<init>: 'this' is null")
	}
	this, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "ConcurrentHashMap.<init>: 'this' is not an object")
	}
	if object.IsNull(this) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ConcurrentHashMap.<init>: 'this' is null")
	}

	this.ThMutex.Lock()
	defer this.ThMutex.Unlock()

	if this.KlassName == 0 || this.KlassName == types.InvalidStringIndex {
		this.KlassName = object.StringPoolIndexFromGoString(classNameConcurrentHashMap)
	}
	this.FieldTable[fieldNameMap] = object.Field{Ftype: types.Ref, Fvalue: make(types.DefHashMap)}
	return nil
}

func concurrentHashMapInitCapacity(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ConcurrentHashMap.<init>: missing capacity parameter")
	}
	initialCapacity, ok := params[1].(int64)
	if !ok || initialCapacity < 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ConcurrentHashMap.<init>: initial capacity < 0")
	}
	return concurrentHashMapInit(params)
}

func concurrentHashMapInitCapacityLoadFactor(params []interface{}) interface{} {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ConcurrentHashMap.<init>: missing parameters")
	}
	initialCapacity, ok := params[1].(int64)
	if !ok || initialCapacity < 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ConcurrentHashMap.<init>: initial capacity < 0")
	}
	var lf float64
	switch v := params[2].(type) {
	case float64:
		lf = v
	case float32:
		lf = float64(v)
	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ConcurrentHashMap.<init>: invalid load factor")
	}
	if lf <= 0 || math.IsNaN(lf) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ConcurrentHashMap.<init>: load factor <= 0")
	}
	return concurrentHashMapInit(params)
}

func concurrentHashMapInitCapacityLoadFactorConcurrencyLevel(params []interface{}) interface{} {
	if len(params) < 4 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ConcurrentHashMap.<init>: missing parameters")
	}
	initialCapacity, ok := params[1].(int64)
	if !ok || initialCapacity < 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ConcurrentHashMap.<init>: initial capacity < 0")
	}
	var lf float64
	switch v := params[2].(type) {
	case float64:
		lf = v
	case float32:
		lf = float64(v)
	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ConcurrentHashMap.<init>: invalid load factor")
	}
	if lf <= 0 || math.IsNaN(lf) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ConcurrentHashMap.<init>: load factor <= 0")
	}
	concurrencyLevel, ok := params[3].(int64)
	if !ok || concurrencyLevel <= 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ConcurrentHashMap.<init>: concurrency level <= 0")
	}
	return concurrentHashMapInit(params)
}

func concurrentHashMapInitMap(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ConcurrentHashMap.<init>: missing map parameter")
	}
	if params[1] == nil || params[1] == object.Null || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ConcurrentHashMap.<init>: map is null")
	}
	if err := concurrentHashMapInit(params); err != nil {
		return err
	}
	res := concurrentHashMapPutAll([]interface{}{params[0], params[1]})
	if gerr, ok := res.(*ghelpers.GErrBlk); ok {
		return gerr
	}
	return nil
}

func concurrentHashMapNewKeySet(params []interface{}) interface{} {
	mapObj := object.MakeEmptyObjectWithClassName(&classNameConcurrentHashMap)
	concurrentHashMapInit([]interface{}{mapObj})
	keySetView := object.MakeEmptyObjectWithClassName(&classNameKeySetView)
	keySetView.FieldTable[fieldNameMap] = mapObj.FieldTable[fieldNameMap]
	return keySetView
}

func concurrentHashMapNewKeySetCapacity(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "newKeySet: missing capacity parameter")
	}
	initialCapacity, ok := params[0].(int64)
	if !ok || initialCapacity < 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "newKeySet: initial capacity < 0")
	}
	mapObj := object.MakeEmptyObjectWithClassName(&classNameConcurrentHashMap)
	concurrentHashMapInit([]interface{}{mapObj})
	keySetView := object.MakeEmptyObjectWithClassName(&classNameKeySetView)
	keySetView.FieldTable[fieldNameMap] = mapObj.FieldTable[fieldNameMap]
	return keySetView
}

func concurrentHashMapClear(params []interface{}) interface{} {
	this, _, gerr := _getConcurrentHashMap(params, "concurrentHashMapClear")
	if gerr != nil {
		return gerr
	}
	this.ThMutex.Lock()
	defer this.ThMutex.Unlock()
	this.FieldTable[fieldNameMap] = object.Field{Ftype: types.Ref, Fvalue: make(types.DefHashMap)}
	return nil
}

func concurrentHashMapPut(params []interface{}) interface{} {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "concurrentHashMapPut: requires 3 parameters")
	}
	this, hm, gerr := _getConcurrentHashMap(params, "concurrentHashMapPut")
	if gerr != nil {
		return gerr
	}
	key, kerr := _getKeyConcurrent(params[1], "concurrentHashMapPut")
	if kerr != nil {
		return kerr
	}
	val, verr := _getValueConcurrent(params[2], "concurrentHashMapPut")
	if verr != nil {
		return verr
	}

	this.ThMutex.Lock()
	defer this.ThMutex.Unlock()

	prev, exists := hm[key]
	if !exists {
		prev = object.Null
	}
	hm[key] = val
	return prev
}

func concurrentHashMapPutIfAbsent(params []interface{}) interface{} {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "concurrentHashMapPutIfAbsent: requires 3 parameters")
	}
	this, hm, gerr := _getConcurrentHashMap(params, "concurrentHashMapPutIfAbsent")
	if gerr != nil {
		return gerr
	}
	key, kerr := _getKeyConcurrent(params[1], "concurrentHashMapPutIfAbsent")
	if kerr != nil {
		return kerr
	}
	val, verr := _getValueConcurrent(params[2], "concurrentHashMapPutIfAbsent")
	if verr != nil {
		return verr
	}

	this.ThMutex.Lock()
	defer this.ThMutex.Unlock()

	if prev, exists := hm[key]; exists {
		return prev
	}
	hm[key] = val
	return object.Null
}

func concurrentHashMapGet(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "concurrentHashMapGet: requires 2 parameters")
	}
	this, hm, gerr := _getConcurrentHashMap(params, "concurrentHashMapGet")
	if gerr != nil {
		return gerr
	}
	key, kerr := _getKeyConcurrent(params[1], "concurrentHashMapGet")
	if kerr != nil {
		return kerr
	}

	this.ThMutex.RLock()
	defer this.ThMutex.RUnlock()

	if val, exists := hm[key]; exists {
		return val
	}
	return object.Null
}

func concurrentHashMapGetOrDefault(params []interface{}) interface{} {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "concurrentHashMapGetOrDefault: requires 3 parameters")
	}
	this, hm, gerr := _getConcurrentHashMap(params, "concurrentHashMapGetOrDefault")
	if gerr != nil {
		return gerr
	}
	key, kerr := _getKeyConcurrent(params[1], "concurrentHashMapGetOrDefault")
	if kerr != nil {
		return kerr
	}

	this.ThMutex.RLock()
	defer this.ThMutex.RUnlock()

	if val, exists := hm[key]; exists {
		return val
	}
	return params[2]
}

func concurrentHashMapContainsKey(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "concurrentHashMapContainsKey: requires 2 parameters")
	}
	this, hm, gerr := _getConcurrentHashMap(params, "concurrentHashMapContainsKey")
	if gerr != nil {
		return gerr
	}
	key, kerr := _getKeyConcurrent(params[1], "concurrentHashMapContainsKey")
	if kerr != nil {
		return kerr
	}

	this.ThMutex.RLock()
	defer this.ThMutex.RUnlock()

	if _, exists := hm[key]; exists {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func concurrentHashMapContainsValue(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "concurrentHashMapContainsValue: requires 2 parameters")
	}
	this, hm, gerr := _getConcurrentHashMap(params, "concurrentHashMapContainsValue")
	if gerr != nil {
		return gerr
	}
	val, verr := _getValueConcurrent(params[1], "concurrentHashMapContainsValue")
	if verr != nil {
		return verr
	}

	this.ThMutex.RLock()
	defer this.ThMutex.RUnlock()

	for _, v := range hm {
		if eq, _ := EqualArrayListElements(v, val); eq {
			return types.JavaBoolTrue
		}
	}
	return types.JavaBoolFalse
}

func concurrentHashMapContains(params []interface{}) interface{} {
	return concurrentHashMapContainsValue(params)
}

func concurrentHashMapRemove(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "concurrentHashMapRemove: requires 2 parameters")
	}
	this, hm, gerr := _getConcurrentHashMap(params, "concurrentHashMapRemove")
	if gerr != nil {
		return gerr
	}
	key, kerr := _getKeyConcurrent(params[1], "concurrentHashMapRemove")
	if kerr != nil {
		return kerr
	}

	this.ThMutex.Lock()
	defer this.ThMutex.Unlock()

	if val, exists := hm[key]; exists {
		delete(hm, key)
		return val
	}
	return object.Null
}

func concurrentHashMapRemoveEntry(params []interface{}) interface{} {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "concurrentHashMapRemove: requires 3 parameters")
	}
	this, hm, gerr := _getConcurrentHashMap(params, "concurrentHashMapRemove")
	if gerr != nil {
		return gerr
	}
	key, kerr := _getKeyConcurrent(params[1], "concurrentHashMapRemove")
	if kerr != nil {
		return kerr
	}
	val, verr := _getValueConcurrent(params[2], "concurrentHashMapRemove")
	if verr != nil {
		return verr
	}

	this.ThMutex.Lock()
	defer this.ThMutex.Unlock()

	if existingVal, exists := hm[key]; exists {
		if eq, _ := EqualArrayListElements(existingVal, val); eq {
			delete(hm, key)
			return types.JavaBoolTrue
		}
	}
	return types.JavaBoolFalse
}

func concurrentHashMapReplace(params []interface{}) interface{} {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "concurrentHashMapReplace: requires 3 parameters")
	}
	this, hm, gerr := _getConcurrentHashMap(params, "concurrentHashMapReplace")
	if gerr != nil {
		return gerr
	}
	key, kerr := _getKeyConcurrent(params[1], "concurrentHashMapReplace")
	if kerr != nil {
		return kerr
	}
	val, verr := _getValueConcurrent(params[2], "concurrentHashMapReplace")
	if verr != nil {
		return verr
	}

	this.ThMutex.Lock()
	defer this.ThMutex.Unlock()

	if prev, exists := hm[key]; exists {
		hm[key] = val
		return prev
	}
	return object.Null
}

func concurrentHashMapReplaceOldNew(params []interface{}) interface{} {
	if len(params) < 4 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "concurrentHashMapReplace: requires 4 parameters")
	}
	this, hm, gerr := _getConcurrentHashMap(params, "concurrentHashMapReplace")
	if gerr != nil {
		return gerr
	}
	key, kerr := _getKeyConcurrent(params[1], "concurrentHashMapReplace")
	if kerr != nil {
		return kerr
	}
	oldVal, oerr := _getValueConcurrent(params[2], "concurrentHashMapReplace")
	if oerr != nil {
		return oerr
	}
	newVal, nerr := _getValueConcurrent(params[3], "concurrentHashMapReplace")
	if nerr != nil {
		return nerr
	}

	this.ThMutex.Lock()
	defer this.ThMutex.Unlock()

	if currentVal, exists := hm[key]; exists {
		if eq, _ := EqualArrayListElements(currentVal, oldVal); eq {
			hm[key] = newVal
			return types.JavaBoolTrue
		}
	}
	return types.JavaBoolFalse
}

func concurrentHashMapSize(params []interface{}) interface{} {
	this, hm, gerr := _getConcurrentHashMap(params, "concurrentHashMapSize")
	if gerr != nil {
		return gerr
	}
	this.ThMutex.RLock()
	defer this.ThMutex.RUnlock()
	sz := len(hm)
	if sz > 2147483647 {
		return int64(2147483647)
	}
	return int64(sz)
}

func concurrentHashMapMappingCount(params []interface{}) interface{} {
	this, hm, gerr := _getConcurrentHashMap(params, "concurrentHashMapMappingCount")
	if gerr != nil {
		return gerr
	}
	this.ThMutex.RLock()
	defer this.ThMutex.RUnlock()
	return int64(len(hm))
}

func concurrentHashMapIsEmpty(params []interface{}) interface{} {
	this, hm, gerr := _getConcurrentHashMap(params, "concurrentHashMapIsEmpty")
	if gerr != nil {
		return gerr
	}
	this.ThMutex.RLock()
	defer this.ThMutex.RUnlock()
	if len(hm) == 0 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func concurrentHashMapPutAll(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "concurrentHashMapPutAll: missing parameters")
	}
	this, thisHm, gerr := _getConcurrentHashMap(params, "concurrentHashMapPutAll")
	if gerr != nil {
		return gerr
	}
	if params[1] == nil || params[1] == object.Null || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "concurrentHashMapPutAll: source map is null")
	}
	otherObj, ok := params[1].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "concurrentHashMapPutAll: source map is not an object")
	}
	otherObj.ThMutex.RLock()
	otherFld, ok := otherObj.FieldTable[fieldNameMap]
	otherObj.ThMutex.RUnlock()
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "concurrentHashMapPutAll: source map missing map field")
	}
	otherHm, ok := otherFld.Fvalue.(types.DefHashMap)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "concurrentHashMapPutAll: source map field is invalid")
	}

	otherObj.ThMutex.RLock()
	for k, v := range otherHm {
		if k == nil || k == object.Null || object.IsNull(k) {
			otherObj.ThMutex.RUnlock()
			return ghelpers.GetGErrBlk(excNames.NullPointerException, "concurrentHashMapPutAll: null key in map")
		}
		if v == nil || v == object.Null || object.IsNull(v) {
			otherObj.ThMutex.RUnlock()
			return ghelpers.GetGErrBlk(excNames.NullPointerException, "concurrentHashMapPutAll: null value in map")
		}
	}
	this.ThMutex.Lock()
	for k, v := range otherHm {
		thisHm[k] = v
	}
	this.ThMutex.Unlock()
	otherObj.ThMutex.RUnlock()
	return nil
}

func concurrentHashMapEquals(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "concurrentHashMapEquals: missing parameters")
	}
	this, thisHm, gerr := _getConcurrentHashMap(params, "concurrentHashMapEquals")
	if gerr != nil {
		return gerr
	}
	if params[1] == this {
		return types.JavaBoolTrue
	}
	if params[1] == nil || params[1] == object.Null || object.IsNull(params[1]) {
		return types.JavaBoolFalse
	}
	otherObj, ok := params[1].(*object.Object)
	if !ok {
		return types.JavaBoolFalse
	}
	otherObj.ThMutex.RLock()
	otherFld, ok := otherObj.FieldTable[fieldNameMap]
	otherObj.ThMutex.RUnlock()
	if !ok {
		return types.JavaBoolFalse
	}
	otherHm, ok := otherFld.Fvalue.(types.DefHashMap)
	if !ok {
		return types.JavaBoolFalse
	}
	this.ThMutex.RLock()
	defer this.ThMutex.RUnlock()
	otherObj.ThMutex.RLock()
	defer otherObj.ThMutex.RUnlock()
	if len(thisHm) != len(otherHm) {
		return types.JavaBoolFalse
	}
	for k, v := range thisHm {
		otherVal, exists := otherHm[k]
		if !exists {
			return types.JavaBoolFalse
		}
		if eq, _ := EqualArrayListElements(v, otherVal); !eq {
			return types.JavaBoolFalse
		}
	}
	return types.JavaBoolTrue
}

func javaStringHashCode(s string) int32 {
	var h int32
	for i := range len(s) {
		h = 31*h + int32(s[i])
	}
	return h
}

func concurrentHashMapHashCode(params []interface{}) interface{} {
	_, hm, gerr := _getConcurrentHashMap(params, "concurrentHashMapHashCode")
	if gerr != nil {
		return gerr
	}
	var totalHash int32 = 0
	for k, v := range hm {
		var kHash int32 = 0
		switch key := k.(type) {
		case string:
			kHash = javaStringHashCode(key)
		case int64:
			kHash = int32(key ^ (key >> 32))
		case float64:
			bits := math.Float64bits(key)
			kHash = int32(bits ^ (bits >> 32))
		case *object.Object:
			if object.IsStringObject(key) {
				kHash = javaStringHashCode(object.GoStringFromStringObject(key))
			} else {
				kHash = int32(key.Mark.Hash)
			}
		}
		var vHash int32 = 0
		switch val := v.(type) {
		case string:
			vHash = javaStringHashCode(val)
		case int64:
			vHash = int32(val ^ (val >> 32))
		case float64:
			bits := math.Float64bits(val)
			vHash = int32(bits ^ (bits >> 32))
		case *object.Object:
			if object.IsStringObject(val) {
				vHash = javaStringHashCode(object.GoStringFromStringObject(val))
			} else {
				vHash = int32(val.Mark.Hash)
			}
		}
		totalHash += (kHash ^ vHash)
	}
	return int64(totalHash)
}

func concurrentHashMapToString(params []interface{}) interface{} {
	this, hm, gerr := _getConcurrentHashMap(params, "concurrentHashMapToString")
	if gerr != nil {
		return gerr
	}
	this.ThMutex.RLock()
	defer this.ThMutex.RUnlock()
	if len(hm) == 0 {
		return object.StringObjectFromGoString("{}")
	}
	keys := make([]string, 0, len(hm))
	keyToEntry := make(map[string]any, len(hm))
	for k, v := range hm {
		kStr := object.StringifyAnythingGo(k)
		keys = append(keys, kStr)
		keyToEntry[kStr] = v
	}
	globals.SortCaseInsensitive(&keys)
	var sb strings.Builder
	sb.WriteString("{")
	for i, kStr := range keys {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(kStr)
		sb.WriteString("=")
		sb.WriteString(object.StringifyAnythingGo(keyToEntry[kStr]))
	}
	sb.WriteString("}")
	return object.StringObjectFromGoString(sb.String())
}

func concurrentHashMapKeySet(params []interface{}) interface{} {
	return mapKeySet(params)
}

func concurrentHashMapKeySetWithDefaultValue(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "keySet: missing default value parameter")
	}
	this, _, gerr := _getConcurrentHashMap(params, "concurrentHashMapKeySet")
	if gerr != nil {
		return gerr
	}
	if params[1] == nil || params[1] == object.Null || object.IsNull(params[1]) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "keySet: default value is null")
	}
	keySetView := object.MakeEmptyObjectWithClassName(&classNameKeySetView)
	this.ThMutex.RLock()
	keySetView.FieldTable[fieldNameMap] = this.FieldTable[fieldNameMap]
	keySetView.FieldTable["value"] = object.Field{Ftype: "Ljava/lang/Object;", Fvalue: params[1]}
	this.ThMutex.RUnlock()
	return keySetView
}

func concurrentHashMapValues(params []interface{}) interface{} {
	this, hm, gerr := _getConcurrentHashMap(params, "concurrentHashMapValues")
	if gerr != nil {
		return gerr
	}
	this.ThMutex.RLock()
	defer this.ThMutex.RUnlock()
	valuesList := make([]any, 0, len(hm))
	for _, v := range hm {
		valuesList = append(valuesList, v)
	}
	return object.MakePrimitiveObject(classNameArrayList, types.Ref, valuesList)
}

func concurrentHashMapEntrySet(params []interface{}) interface{} {
	return mapEntrySet(params)
}

func concurrentHashMapKeys(params []interface{}) interface{} {
	this, _, gerr := _getConcurrentHashMap(params, "concurrentHashMapKeys")
	if gerr != nil {
		return gerr
	}
	return NewIterator(this)
}

func concurrentHashMapElements(params []interface{}) interface{} {
	this, hm, gerr := _getConcurrentHashMap(params, "concurrentHashMapElements")
	if gerr != nil {
		return gerr
	}
	this.ThMutex.RLock()
	valuesList := make([]any, 0, len(hm))
	for _, v := range hm {
		valuesList = append(valuesList, v)
	}
	this.ThMutex.RUnlock()
	valCol := object.MakePrimitiveObject(classNameArrayList, types.Ref, valuesList)
	return NewIterator(valCol)
}
