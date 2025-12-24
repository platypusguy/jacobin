/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
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

func Load_Util_LinkedList() {

	MethodSignatures["java/util/LinkedList.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/util/LinkedList.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistInit,
		}

	MethodSignatures["java/util/LinkedList.<init>(Ljava/util/Collection;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/LinkedList.add(ILjava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  linkedlistAddAtIndex,
		}

	MethodSignatures["java/util/LinkedList.add(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistAddLastRetBool,
		}

	MethodSignatures["java/util/LinkedList.addAll(Ljava/util/Collection;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/LinkedList.addAll(ILjava/util/Collection;)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/LinkedList.addFirst(Ljava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistAddFirst,
		}

	MethodSignatures["java/util/LinkedList.addLast(Ljava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistAddLast,
		}

	MethodSignatures["java/util/LinkedList.clear()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistClear,
		}

	MethodSignatures["java/util/LinkedList.clone()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistClone,
		}

	MethodSignatures["java/util/LinkedList.contains(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistContains,
		}

	MethodSignatures["java/util/LinkedList.descendingIterator()Ljava/util/Iterator;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/LinkedList.element()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistGetFirst,
		}

	MethodSignatures["java/util/LinkedList.get(I)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistGet,
		}

	MethodSignatures["java/util/LinkedList.getFirst()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistGetFirst,
		}

	MethodSignatures["java/util/LinkedList.getLast()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistGetLast,
		}

	MethodSignatures["java/util/LinkedList.indexOf(Ljava/lang/Object;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistIndexOf,
		}

	MethodSignatures["java/util/LinkedList.isEmpty()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistIsEmpty,
		}

	MethodSignatures["java/util/LinkedList.iterator()Ljava/util/Iterator;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistIterator,
		}

	MethodSignatures["java/util/LinkedList.lastIndexOf(Ljava/lang/Object;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistLastIndexOf,
		}

	MethodSignatures["java/util/LinkedList.listIterator(I)Ljava/util/ListIterator;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/LinkedList.offer(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistAddLast,
		}

	MethodSignatures["java/util/LinkedList.offerFirst(Ljava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistAddFirst,
		}

	MethodSignatures["java/util/LinkedList.offerLast(Ljava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistAddLast,
		}

	MethodSignatures["java/util/LinkedList.peek()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistGetFirst,
		}

	MethodSignatures["java/util/LinkedList.peekFirst()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistGetFirst,
		}

	MethodSignatures["java/util/LinkedList.peekLast()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistGetLast,
		}

	MethodSignatures["java/util/LinkedList.poll()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistRemoveFirst,
		}

	MethodSignatures["java/util/LinkedList.pollFirst()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistRemoveFirst,
		}

	MethodSignatures["java/util/LinkedList.pollLast()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistRemoveLast,
		}

	MethodSignatures["java/util/LinkedList.pop()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistRemoveFirst,
		}

	MethodSignatures["java/util/LinkedList.push(Ljava/lang/Object;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistAddFirst,
		}

	MethodSignatures["java/util/LinkedList.remove()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistRemove,
		}

	MethodSignatures["java/util/LinkedList.remove(I)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistRemoveAtIndex,
		}

	MethodSignatures["java/util/LinkedList.remove(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistRemoveLastOccurrence,
		}

	MethodSignatures["java/util/LinkedList.removeFirst()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistRemoveFirst,
		}

	MethodSignatures["java/util/LinkedList.removeFirstOccurrence(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistRemoveFirstOccurrence,
		}

	MethodSignatures["java/util/LinkedList.removeLast()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistRemoveLast,
		}

	MethodSignatures["java/util/LinkedList.removeLastOccurrence(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistRemoveLastOccurrence,
		}

	MethodSignatures["java/util/LinkedList.reversed()Ljava/util/LinkedList;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
			//GFunction:  linkedlistReversed,
		}

	MethodSignatures["java/util/LinkedList.set(ILjava/lang/Object;)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  linkedlistSet,
		}

	MethodSignatures["java/util/LinkedList.size()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistSize,
		}

	MethodSignatures["java/util/LinkedList.sort(Ljava/util/Comparator;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/LinkedList.spliterator()Ljava/util/Spliterator;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/LinkedList.toArray()[Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistToArray,
		}

	MethodSignatures["java/util/LinkedList.toArray([Ljava/lang/Object;)[Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  linkedlistToArrayTyped,
		}

	MethodSignatures["java/util/LinkedList.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  linkedlistToString,
		}
}

var classNameLinkedList = "java/util/LinkedList"

// linkedlistInit (<init>) initializes a new LinkedList object.
func linkedlistIterator(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "linkedlistIterator: Invalid self argument")
	}
	return NewIterator(self)
}

func linkedlistInit(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "linkedlistInit: Invalid self argument")
	}

	object.ClearFieldTable(self)
	self.FieldTable["value"] = object.Field{
		Ftype:  types.LinkedList,
		Fvalue: list.New(),
	}
	return nil
}

// getLinkedListFromObject (internal function) extracts the *list.List from the object
func getLinkedListFromObject(self *object.Object) (*list.List, interface{}) {
	field, exists := self.FieldTable["value"]
	if !exists {
		return nil, getGErrBlk(excNames.NullPointerException, "getLinkedListFromObject: LinkedList not initialized")
	}
	llst, ok := field.Fvalue.(*list.List)
	if !ok {
		return nil, getGErrBlk(excNames.VirtualMachineError, "getLinkedListFromObject: Invalid LinkedList storage")
	}
	return llst, nil
}

// newLinkedListObject (internal function) creates a new *object.Object that contains an empty *list.List in its "value" field.
func newLinkedListObject() *object.Object {
	return object.MakePrimitiveObject(classNameLinkedList, types.LinkedList, list.New())
}

// getLinkedListFromObject (internal function) extracts the *list.List from the object
func equalLinkedListElements(argA any, argB any) (bool, *GErrBlk) {
	// Compare based on the actual type of the searched element (argA).
	switch a := argA.(type) {
	case *object.Object:
		// Only String objects are supported for object comparisons per current implementation.
		if !object.IsStringObject(a) {
			gerr := getGErrBlk(excNames.UnsupportedOperationException, "linkedlistContains: Cannot yet suport non-String objects")
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
	gerr := getGErrBlk(excNames.UnsupportedOperationException, errMsg)
	return false, gerr
}
