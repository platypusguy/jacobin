/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
)

func Load_Util_ListIterator() {

	MethodSignatures["java/util/ListIterator.add(Ljava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  listiteratorAdd,
		}

	MethodSignatures["java/util/ListIterator.hasNext()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  listiteratorHasNext,
		}

	MethodSignatures["java/util/ListIterator.hasPrevious()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  listiteratorHasPrevious,
		}

	MethodSignatures["java/util/ListIterator.next()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  listiteratorNext,
		}

	MethodSignatures["java/util/ListIterator.nextIndex()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  listiteratorNextIndex,
		}

	MethodSignatures["java/util/ListIterator.previous()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  listiteratorPrevious,
		}

	MethodSignatures["java/util/ListIterator.previousIndex()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  listiteratorPreviousIndex,
		}

	MethodSignatures["java/util/ListIterator.remove()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  listiteratorRemove,
		}

	MethodSignatures["java/util/ListIterator.set(Ljava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  listiteratorSet,
		}

	MethodSignatures["java/util/ListIterator.forEachRemaining(Ljava/util/function/Consumer;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}
}

type ListIteratorState struct {
	collection   *object.Object
	cursor       int
	lastReturned int
}

func NewListIterator(collection *object.Object, index int) *object.Object {
	state := &ListIteratorState{
		collection:   collection,
		cursor:       index,
		lastReturned: -1,
	}
	li := object.MakePrimitiveObject("java/util/ListIterator", types.Interface, state)
	return li
}

func getListIteratorState(params []interface{}) (*ListIteratorState, interface{}) {
	self := params[0].(*object.Object)
	state, ok := self.FieldTable["value"].Fvalue.(*ListIteratorState)
	if !ok {
		return nil, getGErrBlk(excNames.IllegalStateException, "ListIterator state missing")
	}
	return state, nil
}

func listiteratorHasNext(params []interface{}) interface{} {
	state, err := getListIteratorState(params)
	if err != nil {
		return err
	}
	var size int
	className := object.GoStringFromStringPoolIndex(state.collection.KlassName)
	if className == "java/util/ArrayList" {
		list, _ := getArrayListFromObject(state.collection)
		size = len(list)
	} else if className == "java/util/Vector" {
		list, _ := getVectorFromObject(state.collection)
		size = len(list)
	}
	return types.ConvertGoBoolToJavaBool(state.cursor < size)
}

func listiteratorNext(params []interface{}) interface{} {
	state, err := getListIteratorState(params)
	if err != nil {
		return err
	}
	var list []interface{}
	className := object.GoStringFromStringPoolIndex(state.collection.KlassName)
	if className == "java/util/ArrayList" {
		list, _ = getArrayListFromObject(state.collection)
	} else if className == "java/util/Vector" {
		list, _ = getVectorFromObject(state.collection)
	}

	if state.cursor >= len(list) {
		return getGErrBlk(excNames.NoSuchElementException, "ListIterator.next: no more elements")
	}

	state.lastReturned = state.cursor
	res := list[state.cursor]
	state.cursor++
	return res
}

func listiteratorHasPrevious(params []interface{}) interface{} {
	state, err := getListIteratorState(params)
	if err != nil {
		return err
	}
	return types.ConvertGoBoolToJavaBool(state.cursor > 0)
}

func listiteratorPrevious(params []interface{}) interface{} {
	state, err := getListIteratorState(params)
	if err != nil {
		return err
	}
	var list []interface{}
	className := object.GoStringFromStringPoolIndex(state.collection.KlassName)
	if className == "java/util/ArrayList" {
		list, _ = getArrayListFromObject(state.collection)
	} else if className == "java/util/Vector" {
		list, _ = getVectorFromObject(state.collection)
	}

	if state.cursor <= 0 {
		return getGErrBlk(excNames.NoSuchElementException, "ListIterator.previous: at beginning")
	}

	state.cursor--
	state.lastReturned = state.cursor
	return list[state.cursor]
}

func listiteratorNextIndex(params []interface{}) interface{} {
	state, err := getListIteratorState(params)
	if err != nil {
		return err
	}
	return int64(state.cursor)
}

func listiteratorPreviousIndex(params []interface{}) interface{} {
	state, err := getListIteratorState(params)
	if err != nil {
		return err
	}
	return int64(state.cursor - 1)
}

func listiteratorRemove(params []interface{}) interface{} {
	state, err := getListIteratorState(params)
	if err != nil {
		return err
	}
	if state.lastReturned == -1 {
		return getGErrBlk(excNames.IllegalStateException, "ListIterator.remove: next/previous not called or remove/add already called")
	}

	className := object.GoStringFromStringPoolIndex(state.collection.KlassName)
	if className == "java/util/ArrayList" {
		list, _ := getArrayListFromObject(state.collection)
		list = append(list[:state.lastReturned], list[state.lastReturned+1:]...)
		state.collection.FieldTable["value"] = object.Field{Ftype: types.ArrayList, Fvalue: list}
	} else if className == "java/util/Vector" {
		list, _ := getVectorFromObject(state.collection)
		list = append(list[:state.lastReturned], list[state.lastReturned+1:]...)
		state.collection.FieldTable["value"] = object.Field{Ftype: types.Vector, Fvalue: list}
	}

	if state.lastReturned < state.cursor {
		state.cursor--
	}
	state.lastReturned = -1
	return nil
}

func listiteratorSet(params []interface{}) interface{} {
	state, err := getListIteratorState(params)
	if err != nil {
		return err
	}
	if state.lastReturned == -1 {
		return getGErrBlk(excNames.IllegalStateException, "ListIterator.set: next/previous not called or remove/add already called")
	}

	obj := params[1]
	className := object.GoStringFromStringPoolIndex(state.collection.KlassName)
	if className == "java/util/ArrayList" {
		list, _ := getArrayListFromObject(state.collection)
		list[state.lastReturned] = obj
		state.collection.FieldTable["value"] = object.Field{Ftype: types.ArrayList, Fvalue: list}
	} else if className == "java/util/Vector" {
		list, _ := getVectorFromObject(state.collection)
		list[state.lastReturned] = obj
		state.collection.FieldTable["value"] = object.Field{Ftype: types.Vector, Fvalue: list}
	}
	return nil
}

func listiteratorAdd(params []interface{}) interface{} {
	state, err := getListIteratorState(params)
	if err != nil {
		return err
	}
	obj := params[1]
	className := object.GoStringFromStringPoolIndex(state.collection.KlassName)
	if className == "java/util/ArrayList" {
		list, _ := getArrayListFromObject(state.collection)
		list = append(list, nil)
		copy(list[state.cursor+1:], list[state.cursor:])
		list[state.cursor] = obj
		state.collection.FieldTable["value"] = object.Field{Ftype: types.ArrayList, Fvalue: list}
	} else if className == "java/util/Vector" {
		list, _ := getVectorFromObject(state.collection)
		list = append(list, nil)
		copy(list[state.cursor+1:], list[state.cursor:])
		list[state.cursor] = obj
		state.collection.FieldTable["value"] = object.Field{Ftype: types.Vector, Fvalue: list}
	}
	state.cursor++
	state.lastReturned = -1
	return nil
}
