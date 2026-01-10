/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
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
)

func Load_Util_ArrayList() {

	ghelpers.MethodSignatures["java/util/ArrayList.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  arraylistInit,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.<init>(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  arraylistInitWithCapacity,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.<init>(Ljava/util/Collection;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.add(ILjava/lang/Object;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  arraylistAddAtIndex,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.add(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  arraylistAdd,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.addAll(Ljava/util/Collection;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.addAll(ILjava/util/Collection;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.clear()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  arraylistClear,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.clone()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  arraylistClone,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.contains(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  arraylistContains,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.containsAll(Ljava/util/Collection;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.ensureCapacity(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  arraylistEnsureCapacity,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.forEach(Ljava/util/function/Consumer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.get(I)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  arraylistGet,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.indexOf(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  arraylistIndexOf,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.isEmpty()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  arraylistIsEmpty,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.iterator()Ljava/util/Iterator;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  arraylistIterator,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.lastIndexOf(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  arraylistLastIndexOf,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.listIterator()Ljava/util/ListIterator;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  arraylistListIterator,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.listIterator(I)Ljava/util/ListIterator;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  arraylistListIteratorAtIndex,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.remove(I)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  arraylistRemoveAtIndex,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.remove(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  arraylistRemoveObject,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.removeAll(Ljava/util/Collection;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.removeIf(Ljava/util/function/Predicate;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.replaceAll(Ljava/util/function/UnaryOperator;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.retainAll(Ljava/util/Collection;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.set(ILjava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  arraylistSet,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.size()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  arraylistSize,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.spliterator()Ljava/util/Spliterator;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.subList(II)Ljava/util/List;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.sort(Ljava/util/Comparator;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.toArray()[Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  arraylistToArray,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.toArray([Ljava/lang/Object;)[Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  arraylistToArrayTyped,
		}

	ghelpers.MethodSignatures["java/util/ArrayList.trimToSize()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  arraylistTrimToSize,
		}
}

var classNameArrayList = "java/util/ArrayList"

func arraylistInit(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "arraylistInit: Invalid self argument")
	}

	object.ClearFieldTable(self)
	self.FieldTable["value"] = object.Field{
		Ftype:  types.ArrayList,
		Fvalue: make([]interface{}, 0, 10),
	}
	return nil
}

func arraylistInitWithCapacity(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "arraylistInitWithCapacity: Invalid self argument")
	}

	capacity, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "arraylistInitWithCapacity: Invalid capacity argument")
	}

	if capacity < 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("Illegal Capacity: %d", capacity))
	}

	object.ClearFieldTable(self)
	self.FieldTable["value"] = object.Field{
		Ftype:  types.ArrayList,
		Fvalue: make([]interface{}, 0, int(capacity)),
	}
	return nil
}

func GetArrayListFromObject(self *object.Object) ([]interface{}, interface{}) {
	field, exists := self.FieldTable["value"]
	if !exists {
		return nil, ghelpers.GetGErrBlk(excNames.NullPointerException, "GetArrayListFromObject: ArrayList not initialized")
	}
	list, ok := field.Fvalue.([]interface{})
	if !ok {
		return nil, ghelpers.GetGErrBlk(excNames.VirtualMachineError, "GetArrayListFromObject: Invalid ArrayList storage")
	}
	return list, nil
}

func arraylistAdd(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	element := params[1]

	list, err := GetArrayListFromObject(self)
	if err != nil {
		return err
	}
	list = append(list, element)
	self.FieldTable["value"] = object.Field{Ftype: types.ArrayList, Fvalue: list}

	return types.JavaBoolTrue
}

func arraylistAddAtIndex(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	index, ok1 := params[1].(int64)
	element := params[2]

	if !ok1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "arraylistAddAtIndex: Invalid index argument")
	}

	list, err := GetArrayListFromObject(self)
	if err != nil {
		return err
	}

	if index < 0 || index > int64(len(list)) {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("Index: %d, Size: %d", index, len(list)))
	}

	// Insert element at index
	list = append(list, nil)
	copy(list[index+1:], list[index:])
	list[index] = element

	self.FieldTable["value"] = object.Field{Ftype: types.ArrayList, Fvalue: list}

	return nil
}

func arraylistGet(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	index, ok := params[1].(int64)

	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "arraylistGet: Invalid index argument")
	}

	list, err := GetArrayListFromObject(self)
	if err != nil {
		return err
	}

	if index < 0 || index >= int64(len(list)) {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("Index: %d, Size: %d", index, len(list)))
	}

	return list[index]
}

func arraylistSet(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	index, ok := params[1].(int64)
	element := params[2]

	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "arraylistSet: Invalid index argument")
	}

	list, err := GetArrayListFromObject(self)
	if err != nil {
		return err
	}

	if index < 0 || index >= int64(len(list)) {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("Index: %d, Size: %d", index, len(list)))
	}

	oldValue := list[index]
	list[index] = element

	return oldValue
}

func arraylistSize(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	list, err := GetArrayListFromObject(self)
	if err != nil {
		return err
	}
	return int64(len(list))
}

func arraylistIterator(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "arraylistIterator: Invalid self argument")
	}
	return NewIterator(self)
}

func arraylistListIterator(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "arraylistListIterator: Invalid self argument")
	}
	return NewListIterator(self, 0)
}

func arraylistListIteratorAtIndex(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "arraylistListIteratorAtIndex: Invalid self argument")
	}
	index, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "arraylistListIteratorAtIndex: Invalid index argument")
	}
	return NewListIterator(self, int(index))
}

func arraylistIsEmpty(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	list, err := GetArrayListFromObject(self)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func arraylistClear(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	list, err := GetArrayListFromObject(self)
	if err != nil {
		return err
	}
	list = list[:0]
	self.FieldTable["value"] = object.Field{Ftype: types.ArrayList, Fvalue: list}
	return nil
}

func arraylistRemoveAtIndex(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	index, ok := params[1].(int64)

	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "arraylistRemoveAtIndex: Invalid index argument")
	}

	list, err := GetArrayListFromObject(self)
	if err != nil {
		return err
	}

	if index < 0 || index >= int64(len(list)) {
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("Index: %d, Size: %d", index, len(list)))
	}

	oldValue := list[index]
	list = append(list[:index], list[index+1:]...)
	self.FieldTable["value"] = object.Field{Ftype: types.ArrayList, Fvalue: list}

	return oldValue
}

func arraylistRemoveObject(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	target := params[1]

	list, err := GetArrayListFromObject(self)
	if err != nil {
		return err
	}

	for i, element := range list {
		eq, gerr := EqualArrayListElements(target, element)
		if gerr != nil {
			return gerr
		}
		if eq {
			list = append(list[:i], list[i+1:]...)
			self.FieldTable["value"] = object.Field{Ftype: types.ArrayList, Fvalue: list}
			return types.JavaBoolTrue
		}
	}

	return types.JavaBoolFalse
}

func arraylistIndexOf(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	target := params[1]

	list, err := GetArrayListFromObject(self)
	if err != nil {
		return err
	}

	for i, element := range list {
		eq, gerr := EqualArrayListElements(target, element)
		if gerr != nil {
			return gerr
		}
		if eq {
			return int64(i)
		}
	}

	return int64(-1)
}

func arraylistLastIndexOf(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	target := params[1]

	list, err := GetArrayListFromObject(self)
	if err != nil {
		return err
	}

	for i := len(list) - 1; i >= 0; i-- {
		eq, gerr := EqualArrayListElements(target, list[i])
		if gerr != nil {
			return gerr
		}
		if eq {
			return int64(i)
		}
	}

	return int64(-1)
}

func arraylistContains(params []interface{}) interface{} {
	res := arraylistIndexOf(params)
	if idx, ok := res.(int64); ok && idx >= 0 {
		return types.JavaBoolTrue
	}
	if _, ok := res.(*ghelpers.GErrBlk); ok {
		return res
	}
	return types.JavaBoolFalse
}

func arraylistToArray(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	list, err := GetArrayListFromObject(self)
	if err != nil {
		return err
	}

	// get the target array from params[1] if supplied (the T[] argument)
	if len(params) > 1 && !object.IsNull(params[1]) {
		targetArray := params[1].(*object.Object)
		targetClassName := *stringPool.GetStringPointer(targetArray.KlassName)
		return object.MakePrimitiveObject(targetClassName, targetClassName, list)
	}

	// fallback for no array passed in, default to Object[]
	return object.MakePrimitiveObject(types.ObjectArrayClassName, types.ObjectArrayClassName, list)
}

func arraylistToArrayTyped(params []interface{}) interface{} {
	// Simple implementation: ignore the input array and return a new one
	return arraylistToArray(params)
}

func arraylistEnsureCapacity(params []interface{}) interface{} {
	// ArrayList implementation in Go's slice handles capacity
	return nil
}

func arraylistTrimToSize(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	list, err := GetArrayListFromObject(self)
	if err != nil {
		return err
	}
	newlist := make([]interface{}, len(list))
	copy(newlist, list)
	self.FieldTable["value"] = object.Field{Ftype: types.ArrayList, Fvalue: newlist}
	return nil
}

func arraylistClone(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	list, err := GetArrayListFromObject(self)
	if err != nil {
		return err
	}
	newlist := make([]interface{}, len(list))
	copy(newlist, list)

	clone := object.MakePrimitiveObject(classNameArrayList, types.ArrayList, newlist)
	return clone
}

func EqualArrayListElements(argA any, argB any) (bool, *ghelpers.GErrBlk) {
	if argA == nil {
		return argB == nil, nil
	}
	if argB == nil {
		return false, nil
	}

	switch a := argA.(type) {
	case *object.Object:
		bObj, ok := argB.(*object.Object)
		if !ok || bObj == nil {
			return false, nil
		}
		if object.IsStringObject(a) && object.IsStringObject(bObj) {
			return object.EqualStringObjects(a, bObj), nil
		}
		// For other objects, we might need a more general equals call,
		// but following LinkedList's lead for now.
		return a == bObj, nil
	case int64:
		bInt, ok := argB.(int64)
		return ok && a == bInt, nil
	case float64:
		bFlt, ok := argB.(float64)
		return ok && a == bFlt, nil
	case bool:
		bBool, ok := argB.(bool)
		return ok && a == bBool, nil
	}

	return argA == argB, nil
}
