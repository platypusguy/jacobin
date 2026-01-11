/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
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

func Load_Util_LinkedList() {

	ghelpers.MethodSignatures["java/util/LinkedList.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistInit,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.<init>(Ljava/util/Collection;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.add(ILjava/lang/Object;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  linkedlistAddAtIndex,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.add(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistAddLastRetBool,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.addAll(Ljava/util/Collection;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.addAll(ILjava/util/Collection;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.addFirst(Ljava/lang/Object;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistAddFirst,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.addLast(Ljava/lang/Object;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistAddLast,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.clear()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistClear,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.clone()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistClone,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.contains(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistContains,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.descendingIterator()Ljava/util/Iterator;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.element()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistGetFirst,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.get(I)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistGet,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.getFirst()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistGetFirst,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.getLast()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistGetLast,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.indexOf(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistIndexOf,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.isEmpty()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistIsEmpty,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.iterator()Ljava/util/Iterator;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistIterator,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.lastIndexOf(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistLastIndexOf,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.listIterator(I)Ljava/util/ListIterator;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.offer(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistAddLast,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.offerFirst(Ljava/lang/Object;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistAddFirst,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.offerLast(Ljava/lang/Object;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistAddLast,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.peek()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistGetFirst,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.peekFirst()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistGetFirst,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.peekLast()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistGetLast,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.poll()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistRemoveFirst,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.pollFirst()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistRemoveFirst,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.pollLast()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistRemoveLast,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.pop()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistRemoveFirst,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.push(Ljava/lang/Object;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistAddFirst,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.remove()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistRemove,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.remove(I)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistRemoveAtIndex,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.remove(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistRemoveLastOccurrence,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.removeFirst()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistRemoveFirst,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.removeFirstOccurrence(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistRemoveFirstOccurrence,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.removeLast()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistRemoveLast,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.removeLastOccurrence(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistRemoveLastOccurrence,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.reversed()Ljava/util/LinkedList;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.set(ILjava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  linkedlistSet,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.size()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistSize,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.sort(Ljava/util/Comparator;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.spliterator()Ljava/util/Spliterator;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.toArray()[Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistToArray,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.toArray([Ljava/lang/Object;)[Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistToArrayTyped,
		}

	ghelpers.MethodSignatures["java/util/LinkedList.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  LinkedlistToString,
		}
}

// linkedlistInit (<init>) initializes a new LinkedList object.
func linkedlistIterator(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "linkedlistIterator: Invalid self argument")
	}
	return NewIterator(self)
}

func linkedlistInit(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "linkedlistInit: Invalid self argument")
	}

	object.ClearFieldTable(self)
	self.FieldTable["value"] = object.Field{
		Ftype:  types.LinkedList,
		Fvalue: list.New(),
	}
	return nil
}

// newLinkedListObject (internal function) creates a new *object.Object that contains an empty *list.List in its "value" field.
func newLinkedListObject() *object.Object {
	return object.MakePrimitiveObject(types.ClassNameLinkedList, types.LinkedList, list.New())
}

// getLinkedListFromObject (internal function) extracts the *list.List from the object
func equalLinkedListElements(argA any, argB any) (bool, *ghelpers.GErrBlk) {
	// Compare based on the actual type of the searched element (argA).
	switch a := argA.(type) {
	case *object.Object:
		// Only String objects are supported for object comparisons per current implementation.
		if !object.IsStringObject(a) {
			gerr := ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, "linkedlistContains: Cannot yet suport non-String objects")
			return false, gerr
		}
		// If the list element is not an object, it's definitely not equal to a String object.
		bObj, ok := argB.(*object.Object)
		if !ok || bObj == nil {
			return false, nil
		}
		// If the list element is not a String object, not equal.
		if !object.IsStringObject(bObj) {
			return false, nil
		}
		if object.EqualStringObjects(a, bObj) {
			return true, nil
		}
		return false, nil
	case int64:
		// Only equal if the list element is also an int64 with same value.
		bInt, ok := argB.(int64)
		if !ok {
			return false, nil
		}
		if a == bInt {
			return true, nil
		}
		return false, nil
	case float64:
		// Only equal if the list element is also a float64 with same value.
		bFlt, ok := argB.(float64)
		if !ok {
			return false, nil
		}
		if a == bFlt {
			return true, nil
		}
		return false, nil
	}

	// Unsupported search element type.
	errMsg := fmt.Sprintf("linkedlistContains: Cannot yet suport element type %T", argA)
	gerr := ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, errMsg)
	return false, gerr
}
