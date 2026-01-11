package javaUtil

import (
	"container/list"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
)

// linkedlistAdd inserts the specified element at the specified index in this list.
// The element currently at that position (if any) and any subsequent elements
// are shifted to the right.
func linkedlistAddAtIndex(params []interface{}) interface{} {
	// Get the linked list.
	self := params[0].(*object.Object)
	valueField, found := self.FieldTable["value"]
	if !found {
		errMsg := "linkedlistAddAtIndex: Field 'value' not found in self object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	llst, ok := valueField.Fvalue.(*list.List)
	if !ok {
		errMsg := fmt.Sprintf("linkedlistAddAtIndex: Expected a linked list field, got %T", valueField.Fvalue)
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get index.
	index, ok := params[1].(int64)
	if !ok {
		errMsg := fmt.Sprintf("linkedlistAddAtIndex: Expected integer for index, got %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get element to add to the linked list.
	newElement, ok := params[2].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("linkedlistAddAtIndex: Expected object for element, got %T", params[1])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Add element at the specified index.
	if index < 0 || index > int64(llst.Len()) {
		errMsg := fmt.Sprintf("linkedlistAddAtIndex: Index out of bounds: %d", index)
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
	}

	targetE := llst.Front()
	for ix := int64(0); ix < index && targetE != nil; ix++ {
		targetE = targetE.Next()
	}

	// If reached the end of the linked list before the index, then append element at end.
	// Else insert new element before the indexed list element.
	if targetE == nil {
		llst.PushBack(newElement)
	} else {
		llst.InsertBefore(newElement, targetE)
	}

	return nil
}

// linkedlistAddFirst inserts an element at the beginning of the list.
func linkedlistAddFirst(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "linkedlistAddFirst: Invalid self argument")
	}
	element := params[1]

	llst, err := ghelpers.GetLinkedListFromObject(self)
	if err != nil {
		return err
	}

	llst.PushFront(element)
	return nil
}

// linkedlistAddLast appends an element at the end of the list.
func linkedlistAddLast(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "linkedlistAddLast: Invalid self argument")
	}
	element := params[1]

	llst, gerr := ghelpers.GetLinkedListFromObject(self)
	if gerr != nil {
		return gerr
	}

	llst.PushBack(element)
	return nil
}

// linkedlistAddLastRetBool appends an element at the end of the linked list. Return true.
func linkedlistAddLastRetBool(params []interface{}) interface{} {
	if len(params) != 2 {
		errMsg := fmt.Sprintf("linkedlistAddLastRetBool: Expected 2 arguments, got %d", len(params))
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "linkedlistAddLastRetBool: Invalid self argument")
	}
	element := params[1]

	llst, gerr := ghelpers.GetLinkedListFromObject(self)
	if gerr != nil {
		return gerr
	}

	llst.PushBack(element)
	return types.JavaBoolTrue
}

// linkedlistClear removes all elements from the list
func linkedlistClear(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "linkedlistClear: Invalid self argument")
	}
	llst, err := ghelpers.GetLinkedListFromObject(self)
	if err != nil {
		return err
	}
	for element := llst.Front(); element != nil; {
		next := element.Next()
		llst.Remove(element)
		element = next
	}
	return nil
}

// linkedlistClone returns a shallow copy of the linked list
func linkedlistClone([]interface{}) interface{} {
	return newLinkedListObject()
}

// linkedlistContains checks if the list contains the given element
func linkedlistContains(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "linkedlistContains: Invalid self argument")
	}
	element := params[1]

	llst, err := ghelpers.GetLinkedListFromObject(self)
	if err != nil {
		return err
	}

	// Iterate through the list and check if element is present
	for e := llst.Front(); e != nil; e = e.Next() {
		equal, gerr := equalLinkedListElements(element, e.Value)
		if gerr != nil {
			return gerr
		}
		if equal {
			return types.JavaBoolTrue
		}
	}

	// There isn't a match.
	return types.JavaBoolFalse
}

// linkedlistGet returns the element at the specified index
func linkedlistGet(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "linkedlistGet: Invalid self argument")
	}

	idx, ok := params[1].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "linkedlistGet: Invalid index argument")
	}

	llst, err := ghelpers.GetLinkedListFromObject(self)
	if err != nil {
		return err
	}

	if idx < 0 || int(idx) >= llst.Len() {
		errMsg := fmt.Sprintf("linkedlistGet: Index %d out of bounds for linked list length %d", idx, llst.Len())
		return ghelpers.GetGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
	}

	element := llst.Front()
	for i := int64(0); i < idx; i++ {
		element = element.Next()
	}

	return element.Value
}

func linkedlistGetFirst(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "linkedlistGetFirst: Invalid self argument")
	}

	llst, err := ghelpers.GetLinkedListFromObject(self)
	if err != nil {
		return err
	}

	if llst.Len() == 0 {
		return ghelpers.GetGErrBlk(excNames.NoSuchElementException, "linkedlistGetFirst: LinkedList is empty")
	}

	return llst.Front().Value
}

// linkedlistGetLast returns the last element without removing it.
func linkedlistGetLast(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "linkedlistGetLast: Invalid self argument")
	}

	llst, err := ghelpers.GetLinkedListFromObject(self)
	if err != nil {
		return err
	}

	if llst.Len() == 0 {
		return ghelpers.GetGErrBlk(excNames.NoSuchElementException, "linkedlistGetLast: LinkedList is empty")
	}

	return llst.Back().Value
}

// linkedlistIndexOf returns the index of the first occurrence of the specified element in the linked list.
// If the element is not found, returns 1.
func linkedlistIndexOf(params []interface{}) interface{} {
	// Retrieve the linked list.
	field, ok := params[0].(*object.Object).FieldTable["value"]
	if !ok || field.Ftype != types.LinkedList {
		errMsg := "linkedlistIndexOf: No linked list found in object"
		return ghelpers.GetGErrBlk(excNames.IllegalStateException, errMsg)
	}
	llst, ok := field.Fvalue.(*list.List)
	if !ok {
		errMsg := "linkedlistIndexOf: Value field is not a linked list"
		return ghelpers.GetGErrBlk(excNames.IllegalStateException, errMsg)
	}

	// The argument should be a Java object, i.e., *object.Object
	elementObj, ok := params[1].(*object.Object)
	if !ok || elementObj == nil {
		errMsg := "linkedlistIndexOf: argument is not a valid object"
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	// Iterate through the linked list to find the element
	index := int64(0)
	for elem := llst.Front(); elem != nil; elem = elem.Next() {
		// If we find the matching element, return the index
		equal, gerr := equalLinkedListElements(elementObj, elem.Value)
		if gerr != nil {
			return gerr
		}
		if equal {
			return index
		}
		index++
	}

	// The element was not found.
	return int64(-1)
}

// linkedlistIsEmpty checks whether the list is empty
func linkedlistIsEmpty(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "linkedlistIsEmpty: Invalid self argument")
	}
	llst, err := ghelpers.GetLinkedListFromObject(self)
	if err != nil {
		return err
	}
	if llst.Len() == 0 {
		return types.JavaBoolTrue
	}

	// Not empty.
	return types.JavaBoolFalse
}

// linkedlistLastIndexOf returns the index of the last occurrence of the specified element in the list.
// If the element is not found, it returns -1.
func linkedlistLastIndexOf(args []interface{}) interface{} {
	// The argument should be a Java object, i.e., *object.Object
	self, ok := args[0].(*object.Object)
	if !ok || self == nil {
		errMsg := "linkedlistLastIndexOf: argument is not a valid object"
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	// Retrieve the linked list from the "value" field
	field, ok := self.FieldTable["value"]
	if !ok || field.Ftype != types.LinkedList {
		errMsg := "linkedlistLastIndexOf: no linked list found in object"
		return ghelpers.GetGErrBlk(excNames.IllegalStateException, errMsg)
	}

	// Extract the *list.List
	linkedList, ok := field.Fvalue.(*list.List)
	if !ok {
		errMsg := "linkedlistLastIndexOf: value field is not a *list.List"
		return ghelpers.GetGErrBlk(excNames.IllegalStateException, errMsg)
	}

	// Get the argument.
	elementObj, ok := args[1].(*object.Object)
	if !ok || elementObj == nil {
		errMsg := "linkedlistLastIndexOf: argument is not a valid object"
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	// Iterate through the linked list from the back to the front to find the last occurrence
	for ix, elem := int64(0), linkedList.Back(); elem != nil; elem = elem.Prev() {
		// If we find the matching element (closest to the end), return its distance-from-end immediately
		equal, gerr := equalLinkedListElements(elementObj, elem.Value)
		if gerr != nil {
			return gerr
		}
		if equal {
			return ix
		}

		ix++
	}

	// The element was not found; return -1
	return int64(-1)
}
