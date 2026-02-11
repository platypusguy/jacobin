package javaUtil

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/stringPool"
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
			GFunction:  ghelpers.TrapFunction,
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
			GFunction:  ghelpers.TrapFunction,
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
