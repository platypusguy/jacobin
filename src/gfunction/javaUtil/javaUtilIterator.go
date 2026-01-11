/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"container/list"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
)

func Load_Util_Iterator() {

	ghelpers.MethodSignatures["java/util/Iterator.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/Iterator.hasNext()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  iteratorHasNext,
		}

	ghelpers.MethodSignatures["java/util/Iterator.next()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  iteratorNext,
		}

	ghelpers.MethodSignatures["java/util/Iterator.remove()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  iteratorRemove,
		}

	ghelpers.MethodSignatures["java/util/Iterator.forEachRemaining(Ljava/util/function/Consumer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
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
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "iteratorHasNext: Invalid self argument")
	}

	colObj, ok := self.FieldTable[iteratorCollectionField].Fvalue.(*object.Object)
	if !ok || colObj == nil {
		return types.JavaBoolFalse
	}

	className := object.GoStringFromStringPoolIndex(colObj.KlassName)

	switch className {
	case "java/util/ArrayList", "java/util/Vector":
		var wlist []interface{}
		var err interface{}
		if className == "java/util/ArrayList" {
			wlist, err = GetArrayListFromObject(colObj)
		} else {
			wlist, err = GetVectorFromObject(colObj)
		}
		if err != nil {
			return types.JavaBoolFalse
		}
		index := self.FieldTable[iteratorIndexField].Fvalue.(int64)
		if index < int64(len(wlist)) {
			return types.JavaBoolTrue
		}
	case "java/util/LinkedList":
		nextNode := self.FieldTable[iteratorNextNodeField].Fvalue
		if nextNode != nil && nextNode != (*list.Element)(nil) {
			return types.JavaBoolTrue
		}
	default:
		return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, fmt.Sprintf("iteratorHasNext: Unsupported collection type %s", className))
	}

	return types.JavaBoolFalse
}

func iteratorNext(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "iteratorNext: Invalid self argument")
	}

	colObj, ok := self.FieldTable[iteratorCollectionField].Fvalue.(*object.Object)
	if !ok || colObj == nil {
		return ghelpers.GetGErrBlk(excNames.NoSuchElementException, "iteratorNext: No collection")
	}

	className := object.GoStringFromStringPoolIndex(colObj.KlassName)

	switch className {
	case "java/util/ArrayList", "java/util/Vector":
		var wlist []interface{}
		var err interface{}
		if className == "java/util/ArrayList" {
			wlist, err = GetArrayListFromObject(colObj)
		} else {
			wlist, err = GetVectorFromObject(colObj)
		}
		if err != nil {
			return err
		}
		index := self.FieldTable[iteratorIndexField].Fvalue.(int64)
		if index >= int64(len(wlist)) {
			return ghelpers.GetGErrBlk(excNames.NoSuchElementException, "iteratorNext: Index out of bounds")
		}
		val := wlist[index]
		self.FieldTable[iteratorIndexField] = object.Field{Ftype: types.Int, Fvalue: index + 1}
		self.FieldTable[iteratorLastReturnedIndexField] = object.Field{Ftype: types.Int, Fvalue: index}
		return val

	case "java/util/LinkedList":
		nextNode, ok := self.FieldTable[iteratorNextNodeField].Fvalue.(*list.Element)
		if !ok || nextNode == nil {
			return ghelpers.GetGErrBlk(excNames.NoSuchElementException, "iteratorNext: No more elements")
		}
		val := nextNode.Value
		self.FieldTable[iteratorNextNodeField] = object.Field{Ftype: types.NonArrayObject, Fvalue: nextNode.Next()}
		self.FieldTable[iteratorLastReturnedNodeField] = object.Field{Ftype: types.NonArrayObject, Fvalue: nextNode}
		return val
	}

	return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, fmt.Sprintf("iteratorNext: Unsupported collection type %s", className))
}

func iteratorRemove(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "iteratorRemove: Invalid self argument")
	}

	colObj, ok := self.FieldTable[iteratorCollectionField].Fvalue.(*object.Object)
	if !ok || colObj == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalStateException, "iteratorRemove: No collection")
	}

	className := object.GoStringFromStringPoolIndex(colObj.KlassName)

	switch className {
	case "java/util/ArrayList", "java/util/Vector":
		lastIdx := self.FieldTable[iteratorLastReturnedIndexField].Fvalue.(int64)
		if lastIdx < 0 {
			return ghelpers.GetGErrBlk(excNames.IllegalStateException, "iteratorRemove: next() has not been called, or remove() has already been called")
		}

		var wlist []interface{}
		var err interface{}
		if className == "java/util/ArrayList" {
			wlist, err = GetArrayListFromObject(colObj)
		} else {
			wlist, err = GetVectorFromObject(colObj)
		}
		if err != nil {
			return err
		}

		// Remove element at lastIdx
		if lastIdx >= int64(len(wlist)) {
			return ghelpers.GetGErrBlk(excNames.ConcurrentModificationException, "iteratorRemove: Collection changed")
		}

		wlist = append(wlist[:lastIdx], wlist[lastIdx+1:]...)
		if className == "java/util/ArrayList" {
			colObj.FieldTable["value"] = object.Field{Ftype: types.ArrayList, Fvalue: wlist}
		} else {
			colObj.FieldTable["value"] = object.Field{Ftype: types.Vector, Fvalue: wlist}
		}

		// Update iterator state
		self.FieldTable[iteratorIndexField] = object.Field{Ftype: types.Int, Fvalue: lastIdx}
		self.FieldTable[iteratorLastReturnedIndexField] = object.Field{Ftype: types.Int, Fvalue: int64(-1)}

	case "java/util/LinkedList":
		lastNode, ok := self.FieldTable[iteratorLastReturnedNodeField].Fvalue.(*list.Element)
		if !ok || lastNode == nil {
			return ghelpers.GetGErrBlk(excNames.IllegalStateException, "iteratorRemove: next() has not been called, or remove() has already been called")
		}

		llst, err := ghelpers.GetLinkedListFromObject(colObj)
		if err != nil {
			return err
		}

		llst.Remove(lastNode)
		self.FieldTable[iteratorLastReturnedNodeField] = object.Field{Ftype: types.NonArrayObject, Fvalue: nil}

	default:
		return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, fmt.Sprintf("iteratorRemove: Unsupported collection type %s", className))
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
		l, err := ghelpers.GetLinkedListFromObject(collection)
		if err == nil && l != nil {
			o.FieldTable[iteratorNextNodeField] = object.Field{Ftype: types.NonArrayObject, Fvalue: l.Front()}
		} else {
			o.FieldTable[iteratorNextNodeField] = object.Field{Ftype: types.NonArrayObject, Fvalue: nil}
		}
	}
	return o
}
