/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/stringPool"
	"jacobin/types"
	"jacobin/util"
	"strconv"
)

var classNameObject = "java/lang/Object"

func Load_Util_Hash_Set() {

	MethodSignatures["java/util/HashSet.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/util/HashSet.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  hashmapInit,
		}

	MethodSignatures["java/util/HashSet.<init>(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hashmapInit,
		}

	MethodSignatures["java/util/HashSet.<init>(IF)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  hashmapInit,
		}

	MethodSignatures["java/util/HashSet.add(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hashsetAdd,
		}

	MethodSignatures["java/util/HashSet.addAll(Ljava/util/Collection;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/HashSet.clear()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  hashmapInit,
		}

	MethodSignatures["java/util/HashSet.clone()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  hashmapInit,
		}

	MethodSignatures["java/util/HashSet.contains(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hashsetContains,
		}

	MethodSignatures["java/util/HashSet.containsAll(Ljava/util/Collection;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/HashSet.isEmpty()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  hashsetIsEmpty,
		}

	MethodSignatures["java/util/HashSet.iterator()Ljava/util/Iterator;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/HashSet.newHashSet(I)Ljava/util/HashSet;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hashmapInit,
		}

	MethodSignatures["java/util/HashSet.remove(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  hashsetRemove,
		}

	MethodSignatures["java/util/HashSet.removeAll(Ljava/util/Collection;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/HashSet.retainAll(Ljava/util/Collection;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/HashSet.size()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  hashmapSize,
		}

	MethodSignatures["java/util/HashSet.spliterator()Ljava/util/Spliterator;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/HashSet.toArray()[Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  hashsetToArray,
		}

	MethodSignatures["java/util/HashSet.toArray([Ljava/lang/Object;)[Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

}

// Compute the hash of the object being added to the HashSet.
// Is it already present? Remember for later.
// Use HashMap.put to add key=hash, value=parameter.
// Return true if this entry did not previously exist; else return false.
func hashsetAdd(params []interface{}) interface{} {

	// Validate HashSet object parameter.
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		errMsg := "hashsetAdd: HashSet parameter is nil or not an object"
		return getGErrBlk(excNames.ClassCastException, errMsg)
	}

	// It must be a HashMap object.
	if *stringPool.GetStringPointer(this.KlassName) != classNameHashMap {
		errMsg := "hashsetAdd: HashSet parameter is not a HashMap object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get argument to add.
	that, ok := params[1].(*object.Object)
	if !ok || that == nil {
		errMsg := "hashsetAdd: Argument parameter is nil or not an object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	fld, ok := that.FieldTable["value"]
	if !ok || that == nil {
		errMsg := "hashsetAdd: Argument parameter is missing the \"value\" field"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Hash the argument and build a composite key including the element's class name
	hashedUint64, err := util.HashAnything(fld.Fvalue)
	if err != nil {
		errMsg := fmt.Sprintf("hashsetAdd: util.HashAnything failed, err: %v", err)
		return getGErrBlk(excNames.VirtualMachineError, errMsg)
	}
	className := *stringPool.GetStringPointer(that.KlassName)
	keyString := className + ":" + strconv.FormatUint(hashedUint64, 10)

	// Remember whether the key already exists.
	flagReturn := types.JavaBoolTrue // Assume not existing yet.
	fld = this.FieldTable[fieldNameMap]
	hm, ok := fld.Fvalue.(types.DefHashMap)
	if !ok {
		errMsg := "hashsetAdd: HashMap is not present"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	_, exists := hm[keyString]
	if exists {
		flagReturn = types.JavaBoolFalse
	}

	// Add this argument.
	var extparams = new([]interface{})
	key := object.StringObjectFromGoString(keyString)
	*extparams = append(*extparams, this)
	*extparams = append(*extparams, key)
	*extparams = append(*extparams, that)
	result := hashmapPut(*extparams)
	switch result.(type) {
	case *GErrBlk:
		return result
	}

	return flagReturn
}

// Remove a hash set entry. Return true if something was actually removed; else return false.
func hashsetRemove(params []interface{}) interface{} {
	// Validate HashSet object parameter.
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		errMsg := "hashsetRemove: HashSet parameter is nil or not an object"
		return getGErrBlk(excNames.ClassCastException, errMsg)
	}

	// It must be a HashMap object.
	if *stringPool.GetStringPointer(this.KlassName) != classNameHashMap {
		errMsg := "hashsetRemove: HashSet parameter is not a HashMap object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get argument to remove.
	that, ok := params[1].(*object.Object)
	if !ok || that == nil {
		errMsg := "hashsetRemove: Argument parameter is nil or not an object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	fld, ok := that.FieldTable["value"]
	if !ok || that == nil {
		errMsg := "hashsetRemove: Argument parameter is missing the \"value\" field"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Hash the argument and build a composite key including the element's class name
	hashedUint64, err := util.HashAnything(fld.Fvalue)
	if err != nil {
		errMsg := fmt.Sprintf("hashsetRemove: util.HashAnything failed, err: %v", err)
		return getGErrBlk(excNames.VirtualMachineError, errMsg)
	}
	className := *stringPool.GetStringPointer(that.KlassName)
	hashedString := className + ":" + strconv.FormatUint(hashedUint64, 10)

	// Remove this argument.
	var extparams = new([]interface{})
	key := object.StringObjectFromGoString(hashedString)
	*extparams = append(*extparams, this)
	*extparams = append(*extparams, key)
	result := hashmapRemove(*extparams)
	switch result.(type) {
	case *GErrBlk:
		return result
	}

	if result == object.Null {
		return types.JavaBoolFalse
	}
	return types.JavaBoolTrue
}

// Does the hash map Have the given key?
func hashsetContains(params []interface{}) interface{} {
	// Validate HashSet object parameter.
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		errMsg := "hashsetContains: HashSet parameter is nil or not an object"
		return getGErrBlk(excNames.ClassCastException, errMsg)
	}

	// It must be a HashMap object.
	if *stringPool.GetStringPointer(this.KlassName) != classNameHashMap {
		errMsg := "hashsetContains: HashSet parameter is not a HashMap object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get argument to look for.
	that, ok := params[1].(*object.Object)
	if !ok || that == nil {
		errMsg := "hashsetContains: Argument parameter is nil or not an object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	fld, ok := that.FieldTable["value"]
	if !ok || that == nil {
		errMsg := "hashsetContains: Argument parameter is missing the \"value\" field"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Hash the argument and build a composite key including the element's class name
	hashedUint64, err := util.HashAnything(fld.Fvalue)
	if err != nil {
		errMsg := fmt.Sprintf("hashsetContains: util.HashAnything failed, err: %v", err)
		return getGErrBlk(excNames.VirtualMachineError, errMsg)
	}
	className := *stringPool.GetStringPointer(that.KlassName)
	hashedString := className + ":" + strconv.FormatUint(hashedUint64, 10)

	// Do the search.
	var extparams = new([]interface{})
	key := object.StringObjectFromGoString(hashedString)
	*extparams = append(*extparams, this)
	*extparams = append(*extparams, key)
	return hashmapContainsKey(*extparams)
}

func hashsetIsEmpty(params []interface{}) interface{} {

	result := hashmapSize(params)
	switch result.(type) {
	case *GErrBlk:
		return result
	}

	if result.(int64) == 0 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func hashsetToArray(params []interface{}) interface{} {

	// Validate HashSet object parameter.
	this, ok := params[0].(*object.Object)
	if !ok || this == nil {
		errMsg := "hashsetToArray: HashSet parameter is nil or not an object"
		return getGErrBlk(excNames.ClassCastException, errMsg)
	}

	// It must be a HashMap object.
	if *stringPool.GetStringPointer(this.KlassName) != classNameHashMap {
		errMsg := "hashsetToArray: HashSet parameter is not a HashMap object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get the current hash map.
	fld := this.FieldTable[fieldNameMap]
	hm, ok := fld.Fvalue.(types.DefHashMap)
	if !ok {
		errMsg := "hashsetToArray: HashMap field is not present"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Create an array of objects.
	objArray := make([]*object.Object, 0, len(hm))
	for _, value := range hm {
		objArray = append(objArray, value.(*object.Object))
	}

	return object.MakePrimitiveObject(classNameObject, types.RefArray, objArray)

}
