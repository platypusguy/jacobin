package gfunction

import (
	"container/list"
	"fmt"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/types"
	"sort"
)

// linkedlistRemove removes the specified element from the list
func linkedlistRemove(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.ClassCastException, "linkedlistRemove: Invalid self argument")
	}
	elementObj := params[1]

	llst, err := getLinkedListFromObject(self)
	if err != nil {
		return err
	}

	// Search for the element and remove it
	for e := llst.Front(); e != nil; e = e.Next() {
		equal, gerr := equalLinkedListElements(elementObj, e.Value)
		if gerr != nil {
			return gerr
		}
		if equal {
			llst.Remove(e)
			return elementObj
		}
	}

	// If the element is not found, return null
	return nil
}

// linkedlistRemoveAtIndex removes the element at the specified position in this list.
func linkedlistRemoveAtIndex(args []interface{}) interface{} {
	// The argument should be the index of the element to remove
	index, ok := args[0].(int64)
	if !ok {
		errMsg := "linkedlistRemoveAtIndex: argument is not an int64 index"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Retrieve the linked list from the "value" field
	listObj, ok := args[0].(*object.Object)
	if !ok || listObj == nil {
		errMsg := "linkedlistRemoveAtIndex: argument is not a valid object"
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}

	// Extract the *list.List
	field, exists := listObj.FieldTable["value"]
	if !exists || field.Ftype != types.LinkedList {
		errMsg := "linkedlistRemoveAtIndex: no linked list found in object"
		return getGErrBlk(excNames.IllegalStateException, errMsg)
	}
	linkedList, ok := field.Fvalue.(*list.List)
	if !ok {
		errMsg := "linkedlistRemoveAtIndex: value field is not a *list.List"
		return getGErrBlk(excNames.IllegalStateException, errMsg)
	}

	// Ensure the index is within bounds
	if index < 0 || index >= int64(linkedList.Len()) {
		errMsg := fmt.Sprintf("linkedlistRemoveAtIndex: index %d out of bounds for list of size %d", index, linkedList.Len())
		return getGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
	}

	// Iterate to the specified index
	elem := linkedList.Front()
	for i := int64(0); i < index; i++ {
		elem = elem.Next()
	}

	// Remove the element at the specified index and return its value
	linkedList.Remove(elem)

	// Return the value of the removed element
	return elem.Value
}

// linkedlistRemoveFirst removes and returns the first element.
func linkedlistRemoveFirst(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.ClassCastException, "linkedlistRemoveFirst: Invalid self argument")
	}

	llst, err := getLinkedListFromObject(self)
	if err != nil {
		return err
	}

	if llst.Len() == 0 {
		return getGErrBlk(excNames.NoSuchElementException, "linkedlistRemoveFirst: LinkedList is empty")
	}

	element := llst.Remove(llst.Front())
	return element
}

// linkedlistRemoveLast removes and returns the last element.
func linkedlistRemoveLast(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.ClassCastException, "linkedlistRemoveLast: Invalid self argument")
	}

	llst, err := getLinkedListFromObject(self)
	if err != nil {
		return err
	}

	if llst.Len() == 0 {
		return getGErrBlk(excNames.NoSuchElementException, "linkedlistRemoveLast: LinkedList is empty")
	}

	element := llst.Remove(llst.Back())
	return element
}

// linkedlistRemoveObject removes the first occurrence of the specified element in the list.
// linkedlistRemoveFirstOccurrence
// If the element is found and removed, it returns true. If the element is not found, it returns false.
func linkedlistRemoveFirstOccurrence(params []interface{}) interface{} {
	// Self should be a Java object, i.e., *object.Object
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		errMsg := "linkedlistRemoveObject: argument is not a valid object"
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}

	// Retrieve the linked list from the "value" field
	llstField, ok := self.FieldTable["value"]
	if !ok || llstField.Ftype != types.LinkedList {
		errMsg := "linkedlistRemoveObject: no linked list found in object"
		return getGErrBlk(excNames.IllegalStateException, errMsg)
	}

	// Extract the *list.List
	llst, ok := llstField.Fvalue.(*list.List)
	if !ok {
		errMsg := "linkedlistRemoveObject: value field is not a *list.List"
		return getGErrBlk(excNames.IllegalStateException, errMsg)
	}

	// The argument should be a Java object, i.e., *object.Object
	elementObj, ok := params[1].(*object.Object)
	if !ok || elementObj == nil {
		errMsg := "linkedlistRemoveObject: argument is not a valid object"
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}

	// Iterate through the linked list to find the element
	for elem := llst.Front(); elem != nil; elem = elem.Next() {
		// If the element matches, remove it and return true
		equal, gerr := equalLinkedListElements(elementObj, elem.Value)
		if gerr != nil {
			return gerr
		}
		if equal {
			llst.Remove(elem)
			return types.JavaBoolTrue // Return true to indicate successful removal
		}
	}

	// If the element was not found, return false
	return types.JavaBoolFalse
}

// linkedlistRemoveLastOccurrence removes the last occurrence of the specified element in the list.
// If the element is found and removed, it returns true. If the element is not found, it returns false.
func linkedlistRemoveLastOccurrence(args []interface{}) interface{} {
	self, ok := args[0].(*object.Object)
	if !ok || self == nil {
		errMsg := "linkedlistRemoveLastOccurrence: self is not a valid object"
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}

	// Retrieve the linked list from the "value" field
	llstField, ok := self.FieldTable["value"]
	if !ok || llstField.Ftype != types.LinkedList {
		errMsg := "linkedlistRemoveLastOccurrence: no linked list found in object"
		return getGErrBlk(excNames.IllegalStateException, errMsg)
	}

	// Extract the *list.List
	llst, ok := llstField.Fvalue.(*list.List)
	if !ok {
		errMsg := "linkedlistRemoveLastOccurrence: value field is not a *list.List"
		return getGErrBlk(excNames.IllegalStateException, errMsg)
	}

	// The argument should be a Java object, i.e., *object.Object
	elementObj, ok := args[1].(*object.Object)
	if !ok || elementObj == nil {
		errMsg := "linkedlistRemoveLastOccurrence: argument is not a valid object"
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}

	// Variable to store the last occurrence element
	var lastElem *list.Element

	// Iterate through the linked list to find the last occurrence
	for elem := llst.Front(); elem != nil; elem = elem.Next() {
		equal, gerr := equalLinkedListElements(elementObj, elem.Value)
		if gerr != nil {
			return gerr
		}
		if equal {
			lastElem = elem
		}
	}

	// If the last occurrence was found, remove it
	if lastElem != nil {
		llst.Remove(lastElem)
		return types.JavaBoolTrue // Return true to indicate successful removal
	}

	// If the element was not found, return false
	return types.JavaBoolFalse
}

// linkedlistReversed - Returns a reverse-ordered view of this linked list without performing a deep copy.
// TODO - How should this really work? Trapped for the moment.
// NOTE: The following code is a deep copy, not what is asked for.
func linkedlistReversed(args []interface{}) interface{} {
	if len(args) != 1 {
		errMsg := "linkedlistReversed: expected 1 argument (the linked list object)"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	origObj, ok := args[0].(*object.Object)
	if !ok || origObj == nil {
		errMsg := "linkedlistReversed: argument is not a valid linked list object"
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}

	field, ok := origObj.FieldTable["value"]
	if !ok || field.Ftype != types.LinkedList {
		errMsg := "linkedlistReversed: no linked list found in object"
		return getGErrBlk(excNames.IllegalStateException, errMsg)
	}

	origList, ok := field.Fvalue.(*list.List)
	if !ok {
		errMsg := "linkedlistReversed: value field is not a *list.List"
		return getGErrBlk(excNames.IllegalStateException, errMsg)
	}

	// Create a new list and add elements in reverse order
	newList := list.New()
	for e := origList.Back(); e != nil; e = e.Prev() {
		newList.PushBack(e.Value)
	}

	// Wrap in a new object
	reversedObj := &object.Object{
		KlassName: origObj.KlassName,
		FieldTable: map[string]object.Field{
			"value": {
				Ftype:  types.LinkedList,
				Fvalue: newList,
			},
		},
	}

	return reversedObj
}

// linkedlistSet replaces the element at the specified position in the linked list with the specified element.
func linkedlistSet(params []interface{}) interface{} {
	// Get the linked list.
	self, ok := params[0].(*object.Object)
	if !ok || self == nil {
		errMsg := "linkedlistSet: self is not a valid object"
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}
	field, exists := self.FieldTable["value"]
	if !exists || field.Ftype != types.LinkedList {
		errMsg := "linkedlistSet: no linked list found in object"
		return getGErrBlk(excNames.IllegalStateException, errMsg)
	}
	llst, ok := field.Fvalue.(*list.List)
	if !ok {
		errMsg := "linkedlistSet: value field is not a *list.List"
		return getGErrBlk(excNames.IllegalStateException, errMsg)
	}

	// Get the index.
	index, ok := params[1].(int64)
	if !ok {
		errMsg := "linkedlistSet: first argument is not an int64 index"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Ensure the index is within bounds.
	if index < 0 || index >= int64(llst.Len()) {
		errMsg := fmt.Sprintf("linkedlistSet: index %d out of bounds for list of size %d", index, llst.Len())
		return getGErrBlk(excNames.IndexOutOfBoundsException, errMsg)
	}

	// Get the argument object.
	listObj, ok := params[2].(*object.Object)
	if !ok || listObj == nil {
		errMsg := "linkedlistSet: second argument is not a valid object"
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}

	// Find the element at the specified index
	elem := llst.Front()
	for i := int64(0); i < index; i++ {
		elem = elem.Next()
	}

	// Store the old value to return it
	oldValue := elem.Value
	// Replace the value at the current position
	elem.Value = listObj

	return oldValue
}

// linkedlistSize returns the number of elements in the list
func linkedlistSize(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.ClassCastException, "linkedlistSize: Invalid self argument")
	}
	llst, err := getLinkedListFromObject(self)
	if err != nil {
		return err
	}
	return int64(llst.Len())
}

// linkedlistSort sorts the elements of the LinkedList in ascending order using a provided comparator.
// If the list cannot be sorted due to invalid elements or comparator issues, it returns an exception.
// TODO: Implement actual sorting.
func linkedlistSort(args []interface{}) interface{} {
	// The argument should be a Java object, i.e., *object.Object, which is the comparator function
	comparator, ok := args[0].(*object.Object)
	if !ok || comparator == nil {
		errMsg := "linkedlistSort: argument is not a valid comparator object"
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}

	// Retrieve the linked list from the "value" field
	listObj, ok := comparator.FieldTable["value"]
	if !ok || listObj.Ftype != types.LinkedList {
		errMsg := "linkedlistSort: no linked list found in object"
		return getGErrBlk(excNames.IllegalStateException, errMsg)
	}

	// Extract the *list.List
	linkedList, ok := listObj.Fvalue.(*list.List)
	if !ok {
		errMsg := "linkedlistSort: value field is not a *list.List"
		return getGErrBlk(excNames.IllegalStateException, errMsg)
	}

	// Convert the linked list to a slice for sorting
	var slice []interface{}
	for elem := linkedList.Front(); elem != nil; elem = elem.Next() {
		slice = append(slice, elem.Value)
	}

	// Use the Go sort function with the provided comparator to sort the slice
	sort.Slice(slice, func(i, j int) bool {
		// Assuming the comparator is a function that compares elements i and j
		// You can implement the logic for invoking the comparator function here
		// For example, if it's a Java method, you can call the corresponding Go function
		// that implements the comparison logic for the two elements
		// Here we assume the elements are directly comparable, replace with actual comparison logic

		// This is a placeholder for actual comparator invocation
		// return someComparator(slice[i], slice[j])

		// If the comparator is not implemented or the elements are not comparable, return false
		return false
	})

	// After sorting, rebuild the linked list
	linkedList.Init() // Clear the existing list
	for _, value := range slice {
		linkedList.PushBack(value)
	}

	return nil
}
