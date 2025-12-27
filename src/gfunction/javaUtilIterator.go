/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"container/list"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
)

func Load_Util_Iterator() {

	MethodSignatures["java/util/Iterator.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/util/Iterator.hasNext()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  iteratorHasNext,
		}

	MethodSignatures["java/util/Iterator.next()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  iteratorNext,
		}

	MethodSignatures["java/util/Iterator.remove()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  iteratorRemove,
		}

	MethodSignatures["java/util/Iterator.forEachRemaining(Ljava/util/function/Consumer;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}
}

const (
	iteratorCollectionField        = "collection"
	iteratorIndexField             = "index"
	iteratorLastReturnedIndexField = "lastReturnedIndex" // for ArrayList
	iteratorNextNodeField          = "nextNode"          // for LinkedList
	iteratorLastReturnedNodeField  = "lastReturnedNode"  // for LinkedList
)

func iteratorHasNext(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "iteratorHasNext: Invalid self argument")
	}

	colObj, ok := self.FieldTable[iteratorCollectionField].Fvalue.(*object.Object)
	if !ok || colObj == nil {
		return types.JavaBoolFalse
	}

	className := object.GoStringFromStringPoolIndex(colObj.KlassName)

	switch className {
	case "java/util/ArrayList", "java/util/Vector":
		var list []interface{}
		var err interface{}
		if className == "java/util/ArrayList" {
			list, err = getArrayListFromObject(colObj)
		} else {
			list, err = getVectorFromObject(colObj)
		}
		if err != nil {
			return types.JavaBoolFalse
		}
		index := self.FieldTable[iteratorIndexField].Fvalue.(int64)
		if index < int64(len(list)) {
			return types.JavaBoolTrue
		}
	case "java/util/LinkedList":
		nextNode := self.FieldTable[iteratorNextNodeField].Fvalue
		if nextNode != nil && nextNode != (*list.Element)(nil) {
			return types.JavaBoolTrue
		}
	default:
		return getGErrBlk(excNames.UnsupportedOperationException, fmt.Sprintf("iteratorHasNext: Unsupported collection type %s", className))
	}

	return types.JavaBoolFalse
}

func iteratorNext(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "iteratorNext: Invalid self argument")
	}

	colObj, ok := self.FieldTable[iteratorCollectionField].Fvalue.(*object.Object)
	if !ok || colObj == nil {
		return getGErrBlk(excNames.NoSuchElementException, "iteratorNext: No collection")
	}

	className := object.GoStringFromStringPoolIndex(colObj.KlassName)

	switch className {
	case "java/util/ArrayList", "java/util/Vector":
		var list []interface{}
		var err interface{}
		if className == "java/util/ArrayList" {
			list, err = getArrayListFromObject(colObj)
		} else {
			list, err = getVectorFromObject(colObj)
		}
		if err != nil {
			return err
		}
		index := self.FieldTable[iteratorIndexField].Fvalue.(int64)
		if index >= int64(len(list)) {
			return getGErrBlk(excNames.NoSuchElementException, "iteratorNext: Index out of bounds")
		}
		val := list[index]
		self.FieldTable[iteratorIndexField] = object.Field{Ftype: types.Int, Fvalue: index + 1}
		self.FieldTable[iteratorLastReturnedIndexField] = object.Field{Ftype: types.Int, Fvalue: index}
		return val

	case "java/util/LinkedList":
		nextNode, ok := self.FieldTable[iteratorNextNodeField].Fvalue.(*list.Element)
		if !ok || nextNode == nil {
			return getGErrBlk(excNames.NoSuchElementException, "iteratorNext: No more elements")
		}
		val := nextNode.Value
		self.FieldTable[iteratorNextNodeField] = object.Field{Ftype: types.NonArrayObject, Fvalue: nextNode.Next()}
		self.FieldTable[iteratorLastReturnedNodeField] = object.Field{Ftype: types.NonArrayObject, Fvalue: nextNode}
		return val
	}

	return getGErrBlk(excNames.UnsupportedOperationException, fmt.Sprintf("iteratorNext: Unsupported collection type %s", className))
}

func iteratorRemove(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "iteratorRemove: Invalid self argument")
	}

	colObj, ok := self.FieldTable[iteratorCollectionField].Fvalue.(*object.Object)
	if !ok || colObj == nil {
		return getGErrBlk(excNames.IllegalStateException, "iteratorRemove: No collection")
	}

	className := object.GoStringFromStringPoolIndex(colObj.KlassName)

	switch className {
	case "java/util/ArrayList", "java/util/Vector":
		lastIdx := self.FieldTable[iteratorLastReturnedIndexField].Fvalue.(int64)
		if lastIdx < 0 {
			return getGErrBlk(excNames.IllegalStateException, "iteratorRemove: next() has not been called, or remove() has already been called")
		}

		var list []interface{}
		var err interface{}
		if className == "java/util/ArrayList" {
			list, err = getArrayListFromObject(colObj)
		} else {
			list, err = getVectorFromObject(colObj)
		}
		if err != nil {
			return err
		}

		// Remove element at lastIdx
		if lastIdx >= int64(len(list)) {
			return getGErrBlk(excNames.ConcurrentModificationException, "iteratorRemove: Collection changed")
		}

		list = append(list[:lastIdx], list[lastIdx+1:]...)
		if className == "java/util/ArrayList" {
			colObj.FieldTable["value"] = object.Field{Ftype: types.ArrayList, Fvalue: list}
		} else {
			colObj.FieldTable["value"] = object.Field{Ftype: types.Vector, Fvalue: list}
		}

		// Update iterator state
		self.FieldTable[iteratorIndexField] = object.Field{Ftype: types.Int, Fvalue: lastIdx}
		self.FieldTable[iteratorLastReturnedIndexField] = object.Field{Ftype: types.Int, Fvalue: int64(-1)}

	case "java/util/LinkedList":
		lastNode, ok := self.FieldTable[iteratorLastReturnedNodeField].Fvalue.(*list.Element)
		if !ok || lastNode == nil {
			return getGErrBlk(excNames.IllegalStateException, "iteratorRemove: next() has not been called, or remove() has already been called")
		}

		llst, err := getLinkedListFromObject(colObj)
		if err != nil {
			return err
		}

		llst.Remove(lastNode)
		self.FieldTable[iteratorLastReturnedNodeField] = object.Field{Ftype: types.NonArrayObject, Fvalue: nil}

	default:
		return getGErrBlk(excNames.UnsupportedOperationException, fmt.Sprintf("iteratorRemove: Unsupported collection type %s", className))
	}

	return nil
}

func NewIterator(collection *object.Object) *object.Object {
	iterClassName := "java/util/Iterator"
	o := object.MakeEmptyObjectWithClassName(&iterClassName)
	o.FieldTable[iteratorCollectionField] = object.Field{Ftype: types.NonArrayObject, Fvalue: collection}

	className := object.GoStringFromStringPoolIndex(collection.KlassName)
	switch className {
	case "java/util/ArrayList", "java/util/Vector":
		o.FieldTable[iteratorIndexField] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
		o.FieldTable[iteratorLastReturnedIndexField] = object.Field{Ftype: types.Int, Fvalue: int64(-1)}
	case "java/util/LinkedList":
		o.FieldTable[iteratorLastReturnedNodeField] = object.Field{Ftype: types.NonArrayObject, Fvalue: nil}
		l, err := getLinkedListFromObject(collection)
		if err == nil && l != nil {
			o.FieldTable[iteratorNextNodeField] = object.Field{Ftype: types.NonArrayObject, Fvalue: l.Front()}
		} else {
			o.FieldTable[iteratorNextNodeField] = object.Field{Ftype: types.NonArrayObject, Fvalue: nil}
		}
	}
	return o
}
