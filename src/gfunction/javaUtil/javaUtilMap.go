package javaUtil

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"strconv"
)

func Load_Util_Map() {

	ghelpers.MethodSignatures["java/util/Map.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/Map.clear()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  mapClear,
		}

	ghelpers.MethodSignatures["java/util/Map.compute(Ljava/lang/Object;Ljava/util/function/BiFunction;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Map.computeIfAbsent(Ljava/lang/Object;Ljava/util/function/Function;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Map.computeIfPresent(Ljava/lang/Object;Ljava/util/function/BiFunction;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Map.containsKey(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  mapContainsKey,
		}

	ghelpers.MethodSignatures["java/util/Map.containsValue(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Map.entrySet()Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  mapEntrySet,
		}

	ghelpers.MethodSignatures["java/util/Map$Entry.comparingByKey()Ljava/util/Comparator;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Map$Entry.comparingByKey(Ljava/util/Comparator;)Ljava/util/Comparator;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Map$Entry.comparingByValue()Ljava/util/Comparator;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Map$Entry.comparingByValue(Ljava/util/Comparator;)Ljava/util/Comparator;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Map$Entry.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Map$Entry.getKey()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  mapEntryGetKey,
		}

	ghelpers.MethodSignatures["java/util/Map$Entry.getValue()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  mapEntryGetValue,
		}

	ghelpers.MethodSignatures["java/util/Map$Entry.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Map$Entry.setValue(Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  mapEntrySetValue,
		}

	ghelpers.MethodSignatures["java/util/AbstractMap$SimpleImmutableEntry.getKey()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  mapEntryGetKey,
		}

	ghelpers.MethodSignatures["java/util/AbstractMap$SimpleImmutableEntry.getValue()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  mapEntryGetValue,
		}

	ghelpers.MethodSignatures["java/util/AbstractMap$SimpleImmutableEntry.setValue(Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  mapEntrySetValue,
		}

	ghelpers.MethodSignatures["java/util/Map.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Map.forEach(Ljava/util/function/BiConsumer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Map.get(Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  mapGet,
		}

	ghelpers.MethodSignatures["java/util/Map.getOrDefault(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  mapGetOrDefault,
		}

	ghelpers.MethodSignatures["java/util/Map.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Map.isEmpty()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  mapIsEmpty,
		}

	ghelpers.MethodSignatures["java/util/Map.keySet()Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  mapKeySet,
		}

	ghelpers.MethodSignatures["java/util/Map.merge(Ljava/lang/Object;Ljava/lang/Object;Ljava/util/function/BiFunction;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Map.put(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  mapPut,
		}

	ghelpers.MethodSignatures["java/util/Map.putAll(Ljava/util/Map;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  mapPutAll,
		}

	ghelpers.MethodSignatures["java/util/Map.putIfAbsent(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  mapPutIfAbsent,
		}

	ghelpers.MethodSignatures["java/util/Map.remove(Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  mapRemove,
		}

	ghelpers.MethodSignatures["java/util/Map.remove(Ljava/lang/Object;Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Map.replace(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Map.replace(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Map.replaceAll(Ljava/util/function/BiFunction;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Map.size()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  mapSize,
		}

	ghelpers.MethodSignatures["java/util/Map.values()Ljava/util/Collection;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}
}

func mapClear(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapClear: missing 'this' parameter")
	}
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "mapClear: 'this' is not an object")
	}

	className := *stringPool.GetStringPointer(this.KlassName)
	switch className {
	case "java/util/HashMap":
		return hashmapInit(params)
	default:
		return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, fmt.Sprintf("mapClear not supported for class %s", className))
	}
}

func mapGet(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapGet: missing parameters")
	}
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "mapGet: 'this' is not an object")
	}

	className := *stringPool.GetStringPointer(this.KlassName)
	switch className {
	case "java/util/HashMap":
		return hashmapGet(params)
	default:
		return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, fmt.Sprintf("mapGet not supported for class %s", className))
	}
}

func mapGetOrDefault(params []interface{}) interface{} {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapGetOrDefault: missing parameters")
	}
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "mapGetOrDefault: 'this' is not an object")
	}

	v := mapGet([]any{this, params[1]})
	if v == object.Null {
		return params[2]
	}
	return v
}

func mapContainsKey(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapContainsKey: missing parameters")
	}
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "mapContainsKey: 'this' is not an object")
	}

	className := *stringPool.GetStringPointer(this.KlassName)
	switch className {
	case "java/util/HashMap":
		return hashmapContainsKey(params)
	default:
		return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, fmt.Sprintf("mapContainsKey not supported for class %s", className))
	}
}

func mapIsEmpty(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapIsEmpty: missing 'this' parameter")
	}
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "mapIsEmpty: 'this' is not an object")
	}

	className := *stringPool.GetStringPointer(this.KlassName)
	switch className {
	case "java/util/HashMap":
		return hashmapIsEmpty(params)
	default:
		return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, fmt.Sprintf("mapIsEmpty not supported for class %s", className))
	}
}

func mapPut(params []interface{}) interface{} {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapPut: missing parameters")
	}
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "mapPut: 'this' is not an object")
	}

	className := *stringPool.GetStringPointer(this.KlassName)
	switch className {
	case "java/util/HashMap":
		return hashmapPut(params)
	default:
		return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, fmt.Sprintf("mapPut not supported for class %s", className))
	}
}

func mapRemove(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapRemove: missing parameters")
	}
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "mapRemove: 'this' is not an object")
	}

	className := *stringPool.GetStringPointer(this.KlassName)
	switch className {
	case "java/util/HashMap":
		return hashmapRemove(params)
	default:
		return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, fmt.Sprintf("mapRemove not supported for class %s", className))
	}
}

func mapSize(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapSize: missing 'this' parameter")
	}
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "mapSize: 'this' is not an object")
	}

	className := *stringPool.GetStringPointer(this.KlassName)
	switch className {
	case "java/util/HashMap":
		return hashmapSize(params)
	default:
		return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, fmt.Sprintf("mapSize not supported for class %s", className))
	}
}

func mapPutAll(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapPutAll: missing parameters")
	}
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "mapPutAll: 'this' is not an object")
	}

	className := *stringPool.GetStringPointer(this.KlassName)
	switch className {
	case "java/util/HashMap":
		return hashmapPutAll(params)
	default:
		return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, fmt.Sprintf("mapPutAll not supported for class %s", className))
	}
}

func mapPutIfAbsent(params []interface{}) interface{} {
	if len(params) < 3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapPutIfAbsent: missing parameters")
	}
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "mapPutIfAbsent: 'this' is not an object")
	}

	// Default implementation for putIfAbsent
	v := mapGet([]any{this, params[1]})
	if v == object.Null {
		return mapPut(params)
	}
	return v
}

func mapEntrySet(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapEntrySet: missing 'this' parameter")
	}
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "mapEntrySet: 'this' is not an object")
	}

	// Get the current hash map.
	this.ThMutex.RLock()
	fld, ok := this.FieldTable[fieldNameMap]
	this.ThMutex.RUnlock()
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapEntrySet: 'this' does not have a 'map' field")
	}
	hm, ok := fld.Fvalue.(types.DefHashMap)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapEntrySet: 'map' field is not a HashMap")
	}

	// Create a new HashSet (which is a HashMap object in Jacobin)
	entrySet := object.MakeEmptyObjectWithClassName(&classNameHashSet)
	if ret := hashmapInit([]any{entrySet}); ret != nil {
		return ret
	}

	// Iterate over the source map
	for k, v := range hm {
		// Represent each entry as a SimpleImmutableEntry
		entryObj := object.MakeEmptyObjectWithClassName(new("java/util/AbstractMap$SimpleImmutableEntry"))
		entryObj.ThMutex.Lock()
		entryObj.FieldTable["key"] = object.Field{Ftype: "Ljava/lang/Object;", Fvalue: k}
		entryObj.FieldTable["value"] = object.Field{Ftype: "Ljava/lang/Object;", Fvalue: v}
		entryObj.ThMutex.Unlock()

		// Add entry to entrySet (which is a HashSet, so it needs hashed key)
		// We use the object's hash code (from its address) to avoid JSON marshal issues with unsafe.Pointer
		keyString := "java/util/AbstractMap$SimpleImmutableEntry:" + strconv.FormatUint(uint64(entryObj.Mark.Hash), 10)
		keyObj := object.StringObjectFromGoString(keyString)

		hashmapPut([]any{entrySet, keyObj, entryObj})
	}

	return entrySet
}

func mapEntryGetKey(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapEntryGetKey: missing 'this' parameter")
	}
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "mapEntryGetKey: 'this' is not an object")
	}
	this.ThMutex.RLock()
	defer this.ThMutex.RUnlock()
	fld, ok := this.FieldTable["key"]
	if !ok {
		return object.Null
	}
	return fld.Fvalue
}

func mapEntryGetValue(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapEntryGetValue: missing 'this' parameter")
	}
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "mapEntryGetValue: 'this' is not an object")
	}
	this.ThMutex.RLock()
	defer this.ThMutex.RUnlock()
	fld, ok := this.FieldTable["value"]
	if !ok {
		return object.Null
	}
	return fld.Fvalue
}

func mapEntrySetValue(params []interface{}) interface{} {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapEntrySetValue: missing parameters")
	}
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "mapEntrySetValue: 'this' is not an object")
	}

	this.ThMutex.Lock()
	defer this.ThMutex.Unlock()

	oldFld, ok := this.FieldTable["value"]
	var oldValue interface{} = object.Null
	if ok {
		oldValue = oldFld.Fvalue
	}

	this.FieldTable["value"] = object.Field{Ftype: "Ljava/lang/Object;", Fvalue: params[1]}
	return oldValue
}

func mapKeySet(params []interface{}) interface{} {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapKeySet: missing 'this' parameter")
	}
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "mapKeySet: 'this' is not an object")
	}

	// Get the current hash map.
	this.ThMutex.RLock()
	fld, ok := this.FieldTable[fieldNameMap]
	this.ThMutex.RUnlock()
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapKeySet: 'this' does not have a 'map' field")
	}
	hm, ok := fld.Fvalue.(types.DefHashMap)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "mapKeySet: 'map' field is not a HashMap")
	}

	// Create a new HashSet (which is a HashMap object in Jacobin)
	keySet := object.MakeEmptyObjectWithClassName(&classNameHashSet)
	if ret := hashmapInit([]any{keySet}); ret != nil {
		return ret
	}

	// Iterate over the source map
	for k := range hm {
		// Add key to keySet (which is a HashSet)
		var keyObj *object.Object
		switch v := k.(type) {
		case string:
			keyObj = object.StringObjectFromGoString(v)
		case int64:
			keyObj = object.MakePrimitiveObject("java/lang/Integer", types.Int, v)
		case float64:
			keyObj = object.MakePrimitiveObject("java/lang/Double", types.Double, v)
		case *object.Object:
			keyObj = v
		default:
			// Fallback for other types if any
			keyObj = object.MakeOneFieldObject("java/lang/Object", "value", "Ljava/lang/Object;", v)
		}

		if ret := hashsetAdd([]any{keySet, keyObj}); ret != nil {
			if _, ok := ret.(*ghelpers.GErrBlk); ok {
				return ret
			}
		}
	}

	return keySet
}
