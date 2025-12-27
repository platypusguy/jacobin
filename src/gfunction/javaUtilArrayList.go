/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
)

func Load_Util_ArrayList() {

	MethodSignatures["java/util/ArrayList.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/util/ArrayList.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  arraylistInit,
		}

	MethodSignatures["java/util/ArrayList.<init>(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  arraylistInitWithCapacity,
		}

	MethodSignatures["java/util/ArrayList.<init>(Ljava/util/Collection;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/ArrayList.add(ILjava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  arraylistAddAtIndex,
		}

	MethodSignatures["java/util/ArrayList.add(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  arraylistAdd,
		}

	MethodSignatures["java/util/ArrayList.addAll(Ljava/util/Collection;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/ArrayList.addAll(ILjava/util/Collection;)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/ArrayList.clear()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  arraylistClear,
		}

	MethodSignatures["java/util/ArrayList.clone()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  arraylistClone,
		}

	MethodSignatures["java/util/ArrayList.contains(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  arraylistContains,
		}

	MethodSignatures["java/util/ArrayList.containsAll(Ljava/util/Collection;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/ArrayList.ensureCapacity(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  arraylistEnsureCapacity,
		}

	MethodSignatures["java/util/ArrayList.forEach(Ljava/util/function/Consumer;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/ArrayList.get(I)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  arraylistGet,
		}

	MethodSignatures["java/util/ArrayList.indexOf(Ljava/lang/Object;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  arraylistIndexOf,
		}

	MethodSignatures["java/util/ArrayList.isEmpty()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  arraylistIsEmpty,
		}

	MethodSignatures["java/util/ArrayList.iterator()Ljava/util/Iterator;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  arraylistIterator,
		}

	MethodSignatures["java/util/ArrayList.lastIndexOf(Ljava/lang/Object;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  arraylistLastIndexOf,
		}

	MethodSignatures["java/util/ArrayList.listIterator()Ljava/util/ListIterator;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  arraylistListIterator,
		}

	MethodSignatures["java/util/ArrayList.listIterator(I)Ljava/util/ListIterator;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  arraylistListIteratorAtIndex,
		}

	MethodSignatures["java/util/ArrayList.remove(I)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  arraylistRemoveAtIndex,
		}

	MethodSignatures["java/util/ArrayList.remove(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  arraylistRemoveObject,
		}

	MethodSignatures["java/util/ArrayList.removeAll(Ljava/util/Collection;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/ArrayList.removeIf(Ljava/util/function/Predicate;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/ArrayList.replaceAll(Ljava/util/function/UnaryOperator;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/ArrayList.retainAll(Ljava/util/Collection;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/ArrayList.set(ILjava/lang/Object;)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  arraylistSet,
		}

	MethodSignatures["java/util/ArrayList.size()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  arraylistSize,
		}

	MethodSignatures["java/util/ArrayList.spliterator()Ljava/util/Spliterator;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/ArrayList.subList(II)Ljava/util/List;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/ArrayList.sort(Ljava/util/Comparator;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/ArrayList.toArray()[Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  arraylistToArray,
		}

	MethodSignatures["java/util/ArrayList.toArray([Ljava/lang/Object;)[Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  arraylistToArrayTyped,
		}

	MethodSignatures["java/util/ArrayList.trimToSize()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  arraylistTrimToSize,
		}
}

var classNameArrayList = "java/util/ArrayList"

func arraylistInit(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "arraylistInit: Invalid self argument")
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
		return getGErrBlk(excNames.IllegalArgumentException, "arraylistInitWithCapacity: Invalid self argument")
	}

	capacity, ok := params[1].(int64)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "arraylistInitWithCapacity: Invalid capacity argument")
	}

	if capacity < 0 {
		return getGErrBlk(excNames.IllegalArgumentException, fmt.Sprintf("Illegal Capacity: %d", capacity))
	}

	object.ClearFieldTable(self)
	self.FieldTable["value"] = object.Field{
		Ftype:  types.ArrayList,
		Fvalue: make([]interface{}, 0, int(capacity)),
	}
	return nil
}

func getArrayListFromObject(self *object.Object) ([]interface{}, interface{}) {
	field, exists := self.FieldTable["value"]
	if !exists {
		return nil, getGErrBlk(excNames.NullPointerException, "getArrayListFromObject: ArrayList not initialized")
	}
	list, ok := field.Fvalue.([]interface{})
	if !ok {
		return nil, getGErrBlk(excNames.VirtualMachineError, "getArrayListFromObject: Invalid ArrayList storage")
	}
	return list, nil
}

func arraylistAdd(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	element := params[1]

	list, err := getArrayListFromObject(self)
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
		return getGErrBlk(excNames.IllegalArgumentException, "arraylistAddAtIndex: Invalid index argument")
	}

	list, err := getArrayListFromObject(self)
	if err != nil {
		return err
	}

	if index < 0 || index > int64(len(list)) {
		return getGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("Index: %d, Size: %d", index, len(list)))
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
		return getGErrBlk(excNames.IllegalArgumentException, "arraylistGet: Invalid index argument")
	}

	list, err := getArrayListFromObject(self)
	if err != nil {
		return err
	}

	if index < 0 || index >= int64(len(list)) {
		return getGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("Index: %d, Size: %d", index, len(list)))
	}

	return list[index]
}

func arraylistSet(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	index, ok := params[1].(int64)
	element := params[2]

	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "arraylistSet: Invalid index argument")
	}

	list, err := getArrayListFromObject(self)
	if err != nil {
		return err
	}

	if index < 0 || index >= int64(len(list)) {
		return getGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("Index: %d, Size: %d", index, len(list)))
	}

	oldValue := list[index]
	list[index] = element

	return oldValue
}

func arraylistSize(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	list, err := getArrayListFromObject(self)
	if err != nil {
		return err
	}
	return int64(len(list))
}

func arraylistIterator(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "arraylistIterator: Invalid self argument")
	}
	return NewIterator(self)
}

func arraylistListIterator(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "arraylistListIterator: Invalid self argument")
	}
	return NewListIterator(self, 0)
}

func arraylistListIteratorAtIndex(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "arraylistListIteratorAtIndex: Invalid self argument")
	}
	index, ok := params[1].(int64)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "arraylistListIteratorAtIndex: Invalid index argument")
	}
	return NewListIterator(self, int(index))
}

func arraylistIsEmpty(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	list, err := getArrayListFromObject(self)
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
	list, err := getArrayListFromObject(self)
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
		return getGErrBlk(excNames.IllegalArgumentException, "arraylistRemoveAtIndex: Invalid index argument")
	}

	list, err := getArrayListFromObject(self)
	if err != nil {
		return err
	}

	if index < 0 || index >= int64(len(list)) {
		return getGErrBlk(excNames.IndexOutOfBoundsException, fmt.Sprintf("Index: %d, Size: %d", index, len(list)))
	}

	oldValue := list[index]
	list = append(list[:index], list[index+1:]...)
	self.FieldTable["value"] = object.Field{Ftype: types.ArrayList, Fvalue: list}

	return oldValue
}

func arraylistRemoveObject(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	target := params[1]

	list, err := getArrayListFromObject(self)
	if err != nil {
		return err
	}

	for i, element := range list {
		eq, gerr := equalArrayListElements(target, element)
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

	list, err := getArrayListFromObject(self)
	if err != nil {
		return err
	}

	for i, element := range list {
		eq, gerr := equalArrayListElements(target, element)
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

	list, err := getArrayListFromObject(self)
	if err != nil {
		return err
	}

	for i := len(list) - 1; i >= 0; i-- {
		eq, gerr := equalArrayListElements(target, list[i])
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
	if _, ok := res.(*GErrBlk); ok {
		return res
	}
	return types.JavaBoolFalse
}

func arraylistToArray(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	list, err := getArrayListFromObject(self)
	if err != nil {
		return err
	}

	return Populator("[Ljava/lang/Object;", types.RefArray, list)
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
	list, err := getArrayListFromObject(self)
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
	list, err := getArrayListFromObject(self)
	if err != nil {
		return err
	}
	newlist := make([]interface{}, len(list))
	copy(newlist, list)

	clone := object.MakePrimitiveObject(classNameArrayList, types.ArrayList, newlist)
	return clone
}

func equalArrayListElements(argA any, argB any) (bool, *GErrBlk) {
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
