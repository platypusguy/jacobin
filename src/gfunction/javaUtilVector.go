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

func Load_Util_Vector() {

	MethodSignatures["java/util/Vector.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/util/Vector.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vectorInit,
		}

	MethodSignatures["java/util/Vector.<init>(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  vectorInitWithCapacity,
		}

	MethodSignatures["java/util/Vector.<init>(II)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  vectorInitWithCapacityAndIncrement,
		}

	MethodSignatures["java/util/Vector.<init>(Ljava/util/Collection;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Vector.add(ILjava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  vectorAddAtIndex,
		}

	MethodSignatures["java/util/Vector.add(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  vectorAdd,
		}

	MethodSignatures["java/util/Vector.addAll(Ljava/util/Collection;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Vector.addAll(ILjava/util/Collection;)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Vector.addElement(Ljava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  vectorAddElement,
		}

	MethodSignatures["java/util/Vector.capacity()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vectorCapacity,
		}

	MethodSignatures["java/util/Vector.clear()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vectorClear,
		}

	MethodSignatures["java/util/Vector.clone()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vectorClone,
		}

	MethodSignatures["java/util/Vector.contains(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  vectorContains,
		}

	MethodSignatures["java/util/Vector.containsAll(Ljava/util/Collection;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Vector.copyInto([Ljava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  vectorCopyInto,
		}

	MethodSignatures["java/util/Vector.elementAt(I)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  vectorElementAt,
		}

	MethodSignatures["java/util/Vector.elements()Ljava/util/Enumeration;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vectorElements,
		}

	MethodSignatures["java/util/Vector.ensureCapacity(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  vectorEnsureCapacity,
		}

	MethodSignatures["java/util/Vector.equals(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  vectorEquals,
		}

	MethodSignatures["java/util/Vector.firstElement()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vectorFirstElement,
		}

	MethodSignatures["java/util/Vector.forEach(Ljava/util/function/Consumer;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Vector.get(I)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  vectorGet,
		}

	MethodSignatures["java/util/Vector.hashCode()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vectorHashCode,
		}

	MethodSignatures["java/util/Vector.indexOf(Ljava/lang/Object;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  vectorIndexOf,
		}

	MethodSignatures["java/util/Vector.indexOf(Ljava/lang/Object;I)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  vectorIndexOfWithIndex,
		}

	MethodSignatures["java/util/Vector.insertElementAt(Ljava/lang/Object;I)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  vectorInsertElementAt,
		}

	MethodSignatures["java/util/Vector.isEmpty()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vectorIsEmpty,
		}

	MethodSignatures["java/util/Vector.iterator()Ljava/util/Iterator;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vectorIterator,
		}

	MethodSignatures["java/util/Vector.lastElement()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vectorLastElement,
		}

	MethodSignatures["java/util/Vector.lastIndexOf(Ljava/lang/Object;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  vectorLastIndexOf,
		}

	MethodSignatures["java/util/Vector.lastIndexOf(Ljava/lang/Object;I)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  vectorLastIndexOfWithIndex,
		}

	MethodSignatures["java/util/Vector.listIterator()Ljava/util/ListIterator;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vectorListIterator,
		}

	MethodSignatures["java/util/Vector.listIterator(I)Ljava/util/ListIterator;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  vectorListIteratorAtIndex,
		}

	MethodSignatures["java/util/Vector.remove(I)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  vectorRemoveAtIndex,
		}

	MethodSignatures["java/util/Vector.remove(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  vectorRemoveObject,
		}

	MethodSignatures["java/util/Vector.removeAll(Ljava/util/Collection;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Vector.removeAllElements()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vectorClear,
		}

	MethodSignatures["java/util/Vector.removeElement(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  vectorRemoveObject,
		}

	MethodSignatures["java/util/Vector.removeElementAt(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  vectorRemoveElementAt,
		}

	MethodSignatures["java/util/Vector.removeIf(Ljava/util/function/Predicate;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Vector.replaceAll(Ljava/util/function/UnaryOperator;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Vector.retainAll(Ljava/util/Collection;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Vector.set(ILjava/lang/Object;)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  vectorSet,
		}

	MethodSignatures["java/util/Vector.setElementAt(Ljava/lang/Object;I)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  vectorSetElementAt,
		}

	MethodSignatures["java/util/Vector.setSize(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  vectorSetSize,
		}

	MethodSignatures["java/util/Vector.size()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vectorSize,
		}

	MethodSignatures["java/util/Vector.sort(Ljava/util/Comparator;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Vector.spliterator()Ljava/util/Spliterator;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Vector.subList(II)Ljava/util/List;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Vector.toArray()[Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vectorToArray,
		}

	MethodSignatures["java/util/Vector.toArray([Ljava/lang/Object;)[Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  vectorToArrayTyped,
		}

	MethodSignatures["java/util/Vector.trimToSize()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vectorTrimToSize,
		}

	MethodSignatures["java/util/Vector.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  vectorToString,
		}
}

func vectorInit(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	self.FieldTable["value"] = object.Field{Ftype: types.Vector, Fvalue: make([]interface{}, 0)}
	return nil
}

func vectorInitWithCapacity(params []interface{}) interface{} {
	if len(params) < 2 {
		return getGErrBlk(excNames.IllegalArgumentException, "vectorInitWithCapacity: too few arguments")
	}
	capacity := int(params[1].(int64))
	if capacity < 0 {
		return getGErrBlk(excNames.IllegalArgumentException, "vectorInitWithCapacity: negative capacity")
	}
	self := params[0].(*object.Object)
	self.FieldTable["value"] = object.Field{Ftype: types.Vector, Fvalue: make([]interface{}, 0, capacity)}
	return nil
}

func vectorInitWithCapacityAndIncrement(params []interface{}) interface{} {
	return vectorInitWithCapacity(params)
}

func getVectorFromObject(self *object.Object) ([]interface{}, interface{}) {
	field, ok := self.FieldTable["value"]
	if !ok {
		return nil, getGErrBlk(excNames.IllegalStateException, "Vector value field missing")
	}
	v, ok := field.Fvalue.([]interface{})
	if !ok {
		return nil, getGErrBlk(excNames.IllegalStateException, "Vector value field is not []interface{}")
	}
	return v, nil
}

func vectorAdd(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	v = append(v, params[1])
	self.FieldTable["value"] = object.Field{Ftype: types.Vector, Fvalue: v}
	return types.JavaBoolTrue
}

func vectorAddAtIndex(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	index := int(params[1].(int64))
	obj := params[2]
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	if index < 0 || index > len(v) {
		return getGErrBlk(excNames.ArrayIndexOutOfBoundsException, fmt.Sprintf("Index: %d, Size: %d", index, len(v)))
	}
	v = append(v, nil)
	copy(v[index+1:], v[index:])
	v[index] = obj
	self.FieldTable["value"] = object.Field{Ftype: types.Vector, Fvalue: v}
	return nil
}

func vectorAddElement(params []interface{}) interface{} {
	vectorAdd(params)
	return nil
}

func vectorGet(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	index := int(params[1].(int64))
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(v) {
		return getGErrBlk(excNames.ArrayIndexOutOfBoundsException, fmt.Sprintf("Index: %d, Size: %d", index, len(v)))
	}
	return v[index]
}

func vectorElementAt(params []interface{}) interface{} {
	return vectorGet(params)
}

func vectorSet(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	index := int(params[1].(int64))
	obj := params[2]
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(v) {
		return getGErrBlk(excNames.ArrayIndexOutOfBoundsException, fmt.Sprintf("Index: %d, Size: %d", index, len(v)))
	}
	old := v[index]
	v[index] = obj
	self.FieldTable["value"] = object.Field{Ftype: types.Vector, Fvalue: v}
	return old
}

func vectorSetElementAt(params []interface{}) interface{} {
	vectorSet(params)
	return nil
}

func vectorSize(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	return int64(len(v))
}

func vectorCapacity(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	return int64(cap(v))
}

func vectorIsEmpty(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	return types.ConvertGoBoolToJavaBool(len(v) == 0)
}

func vectorClear(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	v = v[:0]
	self.FieldTable["value"] = object.Field{Ftype: types.Vector, Fvalue: v}
	return nil
}

func vectorRemoveAtIndex(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	index := int(params[1].(int64))
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(v) {
		return getGErrBlk(excNames.ArrayIndexOutOfBoundsException, fmt.Sprintf("Index: %d, Size: %d", index, len(v)))
	}
	old := v[index]
	v = append(v[:index], v[index+1:]...)
	self.FieldTable["value"] = object.Field{Ftype: types.Vector, Fvalue: v}
	return old
}

func vectorRemoveElementAt(params []interface{}) interface{} {
	vectorRemoveAtIndex(params)
	return nil
}

func vectorRemoveObject(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	obj := params[1]
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	for i, e := range v {
		eq, gerr := equalArrayListElements(e, obj)
		if gerr != nil {
			return gerr
		}
		if eq {
			v = append(v[:i], v[i+1:]...)
			self.FieldTable["value"] = object.Field{Ftype: types.Vector, Fvalue: v}
			return types.JavaBoolTrue
		}
	}
	return types.JavaBoolFalse
}

func vectorIndexOf(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	obj := params[1]
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	for i, e := range v {
		eq, gerr := equalArrayListElements(e, obj)
		if gerr != nil {
			return gerr
		}
		if eq {
			return int64(i)
		}
	}
	return int64(-1)
}

func vectorIndexOfWithIndex(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	obj := params[1]
	index := int(params[2].(int64))
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	if index < 0 {
		return getGErrBlk(excNames.ArrayIndexOutOfBoundsException, fmt.Sprintf("Index: %d", index))
	}
	for i := index; i < len(v); i++ {
		eq, gerr := equalArrayListElements(v[i], obj)
		if gerr != nil {
			return gerr
		}
		if eq {
			return int64(i)
		}
	}
	return int64(-1)
}

func vectorLastIndexOf(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	obj := params[1]
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	for i := len(v) - 1; i >= 0; i-- {
		eq, gerr := equalArrayListElements(v[i], obj)
		if gerr != nil {
			return gerr
		}
		if eq {
			return int64(i)
		}
	}
	return int64(-1)
}

func vectorLastIndexOfWithIndex(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	obj := params[1]
	index := int(params[2].(int64))
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	if index >= len(v) {
		return getGErrBlk(excNames.ArrayIndexOutOfBoundsException, fmt.Sprintf("Index: %d, Size: %d", index, len(v)))
	}
	for i := index; i >= 0; i-- {
		eq, gerr := equalArrayListElements(v[i], obj)
		if gerr != nil {
			return gerr
		}
		if eq {
			return int64(i)
		}
	}
	return int64(-1)
}

func vectorContains(params []interface{}) interface{} {
	res := vectorIndexOf(params)
	if i, ok := res.(int64); ok {
		return types.ConvertGoBoolToJavaBool(i >= 0)
	}
	return res
}

func vectorFirstElement(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	if len(v) == 0 {
		return getGErrBlk(excNames.NoSuchElementException, "Vector is empty")
	}
	return v[0]
}

func vectorLastElement(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	if len(v) == 0 {
		return getGErrBlk(excNames.NoSuchElementException, "Vector is empty")
	}
	return v[len(v)-1]
}

func vectorInsertElementAt(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	obj := params[1]
	index := int(params[2].(int64))
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	if index < 0 || index > len(v) {
		return getGErrBlk(excNames.ArrayIndexOutOfBoundsException, fmt.Sprintf("Index: %d, Size: %d", index, len(v)))
	}
	v = append(v, nil)
	copy(v[index+1:], v[index:])
	v[index] = obj
	self.FieldTable["value"] = object.Field{Ftype: types.Vector, Fvalue: v}
	return nil
}

func vectorSetSize(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	newSize := int(params[1].(int64))
	if newSize < 0 {
		return getGErrBlk(excNames.ArrayIndexOutOfBoundsException, fmt.Sprintf("New size: %d", newSize))
	}
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	if newSize <= len(v) {
		v = v[:newSize]
	} else {
		for i := len(v); i < newSize; i++ {
			v = append(v, object.Null)
		}
	}
	self.FieldTable["value"] = object.Field{Ftype: types.Vector, Fvalue: v}
	return nil
}

func vectorToArray(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	elements := make([]*object.Object, len(v))
	for i, e := range v {
		if obj, ok := e.(*object.Object); ok {
			elements[i] = obj
		} else {
			elements[i] = object.Null
		}
	}
	return object.MakeArrayFromRawArray(elements)
}

func vectorToArrayTyped(params []interface{}) interface{} {
	return vectorToArray(params)
}

func vectorTrimToSize(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	newV := make([]interface{}, len(v))
	copy(newV, v)
	self.FieldTable["value"] = object.Field{Ftype: types.Vector, Fvalue: newV}
	return nil
}

func vectorEnsureCapacity(params []interface{}) interface{} {
	return nil
}

func vectorCopyInto(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	destArr := params[1].(*object.Object)
	field, ok := destArr.FieldTable["value"]
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "Destination array missing value field")
	}

	switch dest := field.Fvalue.(type) {
	case []*object.Object:
		for i, e := range v {
			if i < len(dest) {
				if obj, ok := e.(*object.Object); ok {
					dest[i] = obj
				} else {
					dest[i] = object.Null
				}
			}
		}
	case []interface{}:
		for i, e := range v {
			if i < len(dest) {
				dest[i] = e
			}
		}
	}
	return nil
}

func vectorClone(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	newV := make([]interface{}, len(v))
	copy(newV, v)
	clone := object.MakePrimitiveObject("java/util/Vector", types.Vector, newV)
	return clone
}

func vectorEquals(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	otherObj, ok := params[1].(*object.Object)
	if !ok || otherObj == object.Null {
		return types.JavaBoolFalse
	}
	if object.GoStringFromStringPoolIndex(otherObj.KlassName) != "java/util/Vector" && object.GoStringFromStringPoolIndex(otherObj.KlassName) != "java/util/ArrayList" {
		// Strictly it should be any List, but let's stick to what we have
		return types.JavaBoolFalse
	}
	v1, err1 := getVectorFromObject(self)
	if err1 != nil {
		return err1
	}

	var v2 []interface{}
	if object.GoStringFromStringPoolIndex(otherObj.KlassName) == "java/util/Vector" {
		var err2 interface{}
		v2, err2 = getVectorFromObject(otherObj)
		if err2 != nil {
			return err2
		}
	} else {
		var err2 interface{}
		v2, err2 = getArrayListFromObject(otherObj)
		if err2 != nil {
			return err2
		}
	}

	if len(v1) != len(v2) {
		return types.JavaBoolFalse
	}
	for i := range v1 {
		eq, gerr := equalArrayListElements(v1[i], v2[i])
		if gerr != nil {
			return gerr
		}
		if !eq {
			return types.JavaBoolFalse
		}
	}
	return types.JavaBoolTrue
}

func vectorHashCode(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}
	var hashCode int32 = 1
	for _, e := range v {
		var elementHash int32 = 0
		if e == nil || e == object.Null {
			elementHash = 0
		} else if _, ok := e.(*object.Object); ok {
			// This is a simplification
			elementHash = 0
		} else {
			// Primitive or other
			elementHash = 0
		}
		hashCode = 31*hashCode + elementHash
	}
	return int64(hashCode)
}

func vectorIterator(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	return NewIterator(self)
}

func vectorElements(params []interface{}) interface{} {
	// Vector.elements() returns an Enumeration.
	// For simplicity, we can use the same Iterator mechanism if we wrap it.
	// But let's see if we need java/util/Enumeration.
	self := params[0].(*object.Object)
	return NewIterator(self) // TRAP or implement Enumeration?
}

func vectorListIterator(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	return NewListIterator(self, 0)
}

func vectorListIteratorAtIndex(params []interface{}) interface{} {
	self := params[0].(*object.Object)
	index := int(params[1].(int64))
	return NewListIterator(self, index)
}

func vectorToString(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		return getGErrBlk(excNames.NullPointerException, "vectorToString: self is null")
	}
	v, err := getVectorFromObject(self)
	if err != nil {
		return err
	}

	strBuffer := "["
	for i, element := range v {
		strBuffer += object.StringifyAnythingGo(element)
		if i < len(v)-1 {
			strBuffer += ", "
		}
	}
	strBuffer += "]"

	return object.StringObjectFromGoString(strBuffer)
}
