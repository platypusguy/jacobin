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
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"sync"
)

var classNameHashMap = "java/util/HashMap"
var hashmapMutex = sync.RWMutex{}
var fieldNameMap = "map"

func Load_Util_Hash_Map() {

	ghelpers.MethodSignatures["java/util/HashMap.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/HashMap.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  hashmapInit,
		}

	ghelpers.MethodSignatures["java/util/HashMap.<init>(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  hashmapInit,
		}

	ghelpers.MethodSignatures["java/util/HashMap.<init>(IF)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  hashmapInit,
		}

	ghelpers.MethodSignatures["java/util/HashMap.clear()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  hashmapInit,
		}

	ghelpers.MethodSignatures["java/util/HashMap.clone()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  hashmapInit,
		}

	ghelpers.MethodSignatures["java/util/HashMap.compute(Ljava/lang/Object;Ljava/util/function/BiFunction;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/HashMap.computeIfAbsent(Ljava/lang/Object;Ljava/util/function/BiFunction;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/HashMap.computeIfPresent(Ljava/lang/Object;Ljava/util/function/BiFunction;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/HashMap.containsKey(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  hashmapContainsKey,
		}

	ghelpers.MethodSignatures["java/util/HashMap.containsValue(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/HashMap.entrySet()Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/HashMap.get(Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  hashmapGet,
		}

	ghelpers.MethodSignatures["java/util/HashMap.isEmpty()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  hashmapIsEmpty,
		}

	ghelpers.MethodSignatures["java/util/HashMap.keySet()Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/HashMap.merge(Ljava/lang/Object;Ljava/lang/Object;Ljava/util/function/BiFunction;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/HashMap.newHashMap(I)Ljava/util/HashMap;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  hashmapInit,
		}

	ghelpers.MethodSignatures["java/util/HashMap.put(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  hashmapPut,
		}

	ghelpers.MethodSignatures["java/util/HashMap.putAll(Ljava/util/Map;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  hashmapPutAll,
		}

	ghelpers.MethodSignatures["java/util/HashMap.remove(Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  hashmapRemove,
		}

	ghelpers.MethodSignatures["java/util/HashMap.size()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  hashmapSize,
		}

	ghelpers.MethodSignatures["java/util/HashMap.values()Ljava/util/Collection;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

}

// Initialise a hash map object to an empty state.
func hashmapInit(params []interface{}) interface{} {
	hashmapMutex.Lock()
	defer hashmapMutex.Unlock()
	nilMap := make(types.DefHashMap)
	obj := params[0].(*object.Object)
	obj.KlassName = object.StringPoolIndexFromGoString(classNameHashMap)
	fld := obj.FieldTable[fieldNameMap]
	fld.Ftype = types.HashMap
	fld.Fvalue = nilMap
	obj.FieldTable[fieldNameMap] = fld
	return nil
}

// An internal function to extract the HashMap key field value.
func _getKey(param interface{}) (interface{}, bool) {
	keyObj, ok := param.(*object.Object)
	if !ok || keyObj == nil {
		errMsg := "HashMap:_getKey: Key parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg), false
	}
	if object.IsStringObject(keyObj) {
		return object.GoStringFromStringObject(keyObj), true
	}
	fvalue := keyObj.FieldTable[fieldNameValue].Fvalue
	switch fvalue.(type) {
	case int64, float64:
	default:
		ftype := keyObj.FieldTable[fieldNameValue].Ftype
		errMsg := fmt.Sprintf("HashMap:_getKey: Unsupported key type: {Ftype: %s, Fvalue: %T}", ftype, fvalue)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg), false
	}
	return fvalue, true
}

// Put inserts a key-value pair into the HashMap and returns the previous value or null.
func hashmapPut(params []interface{}) interface{} {
	hashmapMutex.Lock()
	defer hashmapMutex.Unlock()

	if len(params) < 3 {
		errMsg := "hashmapPut: requires 3 parameters: HashMap, key, and value"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		errMsg := "hashmapPut: HashMap parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.ClassCastException, errMsg)
	}

	if *stringPool.GetStringPointer(this.KlassName) != classNameHashMap {
		errMsg := "hashmapPut: HashMap parameter is not a HashMap"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Extract the key.
	key, ok := _getKey(params[1])
	if !ok {
		return key
	}

	// Use the value object as-is.
	value := params[2]

	// Get the current hash map.
	fld := this.FieldTable[fieldNameMap]
	hm, ok := fld.Fvalue.(types.DefHashMap)
	if !ok {
		errMsg := "hashmapPut: HashMap parameter is missing its \"value\" field"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Check for the previous value associated with the key.
	prevValue, exists := hm[key]
	if !exists {
		prevValue = object.Null
	}

	// Store the new key-value pair in the hash map.
	hm[key] = value
	fld.Fvalue = hm
	this.FieldTable[fieldNameMap] = fld

	// Return the previous value for that key.
	return prevValue
}

// Get a hash map entry. Return nil if there is not one that matches the key.
func hashmapGet(params []interface{}) interface{} {
	if len(params) < 2 {
		errMsg := "hashmapGet: Requires 2 parameters: HashMap and key"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		errMsg := "hashmapGet: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.ClassCastException, errMsg)
	}

	if *stringPool.GetStringPointer(this.KlassName) != classNameHashMap {
		errMsg := "hashmapGet: The object is not a HashMap"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Extract the key.
	key, ok := _getKey(params[1])
	if !ok {
		return key
	}

	// Get the current hash map.
	fld := this.FieldTable[fieldNameMap]
	hm, ok := fld.Fvalue.(types.DefHashMap)
	if !ok {
		errMsg := "hashmapGet: The HashMap is not present"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Retrieve the field associated with the key
	value, exists := hm[key]
	if !exists {
		return object.Null
	}

	return value
}

// Remove a hash map entry. Return the removed value or nil if there is not one that matches the key.
func hashmapRemove(params []interface{}) interface{} {
	hashmapMutex.Lock()
	defer hashmapMutex.Unlock()

	if len(params) < 2 {
		errMsg := "hashmapRemove: Requires 2 parameters: HashMap and key"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		errMsg := "hashmapRemove: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.ClassCastException, errMsg)
	}

	if *stringPool.GetStringPointer(this.KlassName) != classNameHashMap {
		errMsg := "hashmapRemove: The object is not a HashMap"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Extract the key.
	key, ok := _getKey(params[1])
	if !ok {
		return key
	}

	// Get the current hash map.
	fld := this.FieldTable[fieldNameMap]
	hm, ok := fld.Fvalue.(types.DefHashMap)
	if !ok {
		errMsg := "hashmapRemove: The HashMap is not present"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Retrieve the field associated with the key
	value, exists := hm[key]
	if !exists {
		return object.Null
	}

	// Delete key-value entry.
	delete(hm, key)

	// Return the deleted value object.
	return value
}

// Get the size of the hash map.
func hashmapSize(params []interface{}) interface{} {

	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		errMsg := "hashmapSize: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.ClassCastException, errMsg)
	}

	if *stringPool.GetStringPointer(this.KlassName) != classNameHashMap {
		errMsg := "hashmapSize: The object is not a HashMap"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get the current hash map.
	fld := this.FieldTable[fieldNameMap]
	hm, ok := fld.Fvalue.(types.DefHashMap)
	if !ok {
		errMsg := "hashmapSize: The HashMap is not present"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Return the current size.
	return int64(len(hm))
}

func hashmapIsEmpty(params []interface{}) interface{} {

	result := hashmapSize(params)
	switch result.(type) {
	case *ghelpers.GErrBlk:
		return result
	}

	if result.(int64) == 0 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func hashmapPutAll(params []interface{}) interface{} {
	hashmapMutex.Lock()
	defer hashmapMutex.Unlock()

	if len(params) < 2 {
		errMsg := "hashmapPutAll: requires 2 parameters: this HashMap and that HashMap"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		errMsg := "hashmapPutAll: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.ClassCastException, errMsg)
	}

	if *stringPool.GetStringPointer(this.KlassName) != classNameHashMap {
		errMsg := "hashmapPutAll: The object is not a HashMap"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get the hash map from this.
	thisFld := this.FieldTable[fieldNameMap]
	thisHmap, ok := thisFld.Fvalue.(types.DefHashMap)
	if !ok {
		errMsg := "hashmapPutAll: The HashMap is not present in the 1st parameter"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	that, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "hashmapPutAll: The 2nd parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.ClassCastException, errMsg)
	}

	if *stringPool.GetStringPointer(this.KlassName) != classNameHashMap {
		errMsg := "hashmapPutAll: The 2nd parameter is not a HashMap"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get the hash map from that.
	thatFld := that.FieldTable[fieldNameMap]
	thatHmap, ok := thatFld.Fvalue.(types.DefHashMap)
	if !ok {
		errMsg := "hashmapPutAll: The HashMap is not present in the 2nd parameter"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	for key, value := range thatHmap {
		thisHmap[key] = value
	}

	return nil
}

// Does the hash map Have the given key?
func hashmapContainsKey(params []interface{}) interface{} {
	if len(params) < 2 {
		errMsg := "hashmapContainsKey: Requires 2 parameters: HashMap and key"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		errMsg := "hashmapContainsKey: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.ClassCastException, errMsg)
	}

	if *stringPool.GetStringPointer(this.KlassName) != classNameHashMap {
		errMsg := "hashmapContainsKey: The object is not a HashMap"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Extract the key.
	key, ok := _getKey(params[1])
	if !ok {
		return key
	}

	// Get the current hash map.
	fld := this.FieldTable[fieldNameMap]
	hm, ok := fld.Fvalue.(types.DefHashMap)
	if !ok {
		errMsg := "hashmapContainsKey: The HashMap is not present"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Retrieve the field associated with the key
	_, exists := hm[key]
	if exists {
		return types.JavaBoolTrue
	}

	return types.JavaBoolFalse
}
