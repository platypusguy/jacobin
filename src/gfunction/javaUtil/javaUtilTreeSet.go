/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"container/list"
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"math"
)

const (
	treeSetClassName        = "java/util/TreeSet"
	treeSetStateField       = "$treeSetState"
	treeSetIteratorValues   = "$treeSetIteratorValues"
	treeSetExpectedModCount = "$treeSetExpectedModCount"
	treeSetIteratorLast     = "$treeSetIteratorLast"
	treeSetReverseField     = "$treeSetReverse"
	treeSetDelegateField    = "$treeSetDelegate"
)

type treeSetData struct {
	elements   []*object.Object
	comparator *object.Object
	modCount   int64
}

// Bounds are expressed in the order visible from the view. This keeps nested
// ranges and descending views using the same membership rules.
type treeSetView struct {
	data           *treeSetData
	lower          *object.Object
	upper          *object.Object
	hasLower       bool
	hasUpper       bool
	lowerInclusive bool
	upperInclusive bool
	descending     bool
}

func Load_Util_Tree_Set() {
	reg := func(sig string, slots int, fn func([]any) any, needsContext bool) {
		ghelpers.MethodSignatures["java/util/TreeSet."+sig] = ghelpers.GMeth{ParamSlots: slots, GFunction: fn, NeedsContext: needsContext}
	}
	reg("<clinit>()V", 0, ghelpers.ClinitGeneric, false)
	reg("<init>()V", 0, treeSetInit, false)
	reg("<init>(Ljava/util/Comparator;)V", 1, treeSetInit, false)
	reg("<init>(Ljava/util/Collection;)V", 1, treeSetInitCollection, true)
	reg("<init>(Ljava/util/SortedSet;)V", 1, treeSetInitCollection, true)
	reg("add(Ljava/lang/Object;)Z", 1, treeSetAdd, true)
	reg("addAll(Ljava/util/Collection;)Z", 1, treeSetAddAll, true)
	reg("clear()V", 0, treeSetClear, false)
	reg("clone()Ljava/lang/Object;", 0, treeSetClone, true)
	reg("comparator()Ljava/util/Comparator;", 0, treeSetComparator, false)
	reg("contains(Ljava/lang/Object;)Z", 1, treeSetContains, true)
	reg("containsAll(Ljava/util/Collection;)Z", 1, treeSetContainsAll, true)
	reg("descendingIterator()Ljava/util/Iterator;", 0, treeSetDescendingIterator, true)
	reg("descendingSet()Ljava/util/NavigableSet;", 0, treeSetDescendingSet, false)
	reg("first()Ljava/lang/Object;", 0, treeSetFirst, true)
	reg("floor(Ljava/lang/Object;)Ljava/lang/Object;", 1, treeSetFloor, true)
	reg("headSet(Ljava/lang/Object;)Ljava/util/SortedSet;", 1, treeSetHeadSetInclusive, true)
	reg("headSet(Ljava/lang/Object;Z)Ljava/util/NavigableSet;", 2, treeSetHeadSet, true)
	reg("higher(Ljava/lang/Object;)Ljava/lang/Object;", 1, treeSetHigher, true)
	reg("isEmpty()Z", 0, treeSetIsEmpty, true)
	reg("iterator()Ljava/util/Iterator;", 0, treeSetIterator, true)
	reg("last()Ljava/lang/Object;", 0, treeSetLast, true)
	reg("lower(Ljava/lang/Object;)Ljava/lang/Object;", 1, treeSetLower, true)
	reg("pollFirst()Ljava/lang/Object;", 0, treeSetPollFirst, true)
	reg("pollLast()Ljava/lang/Object;", 0, treeSetPollLast, true)
	reg("remove(Ljava/lang/Object;)Z", 1, treeSetRemove, true)
	reg("removeAll(Ljava/util/Collection;)Z", 1, treeSetRemoveAll, true)
	reg("retainAll(Ljava/util/Collection;)Z", 1, treeSetRetainAll, true)
	reg("size()I", 0, treeSetSize, true)
	reg("subSet(Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/SortedSet;", 2, treeSetSubSetInclusive, true)
	reg("subSet(Ljava/lang/Object;ZLjava/lang/Object;Z)Ljava/util/NavigableSet;", 4, treeSetSubSet, true)
	reg("tailSet(Ljava/lang/Object;)Ljava/util/SortedSet;", 1, treeSetTailSetInclusive, true)
	reg("tailSet(Ljava/lang/Object;Z)Ljava/util/NavigableSet;", 2, treeSetTailSet, true)
	reg("toArray()[Ljava/lang/Object;", 0, treeSetToArray, true)
	reg("toArray([Ljava/lang/Object;)[Ljava/lang/Object;", 1, treeSetToArray, true)
	reg("getFirst()Ljava/lang/Object;", 0, treeSetFirst, true)
	reg("getLast()Ljava/lang/Object;", 0, treeSetLast, true)
	reg("reversed()Ljava/util/SequencedSet;", 0, treeSetDescendingSet, false)
	reg("addFirst(Ljava/lang/Object;)V", 1, ghelpers.TrapFunction, false)
	reg("addLast(Ljava/lang/Object;)V", 1, ghelpers.TrapFunction, false)
	ghelpers.MethodSignatures["java/util/TreeSet$DescendingComparator.compare(Ljava/lang/Object;Ljava/lang/Object;)I"] = ghelpers.GMeth{ParamSlots: 2, GFunction: treeSetDescendingComparatorCompare, NeedsContext: true}
	reg("removeFirst()Ljava/lang/Object;", 0, treeSetRemoveFirst, true)
	reg("removeLast()Ljava/lang/Object;", 0, treeSetRemoveLast, true)
}

func treeSetContext(params []any) (*list.List, int) {
	if len(params) > 0 {
		if fs, ok := params[0].(*list.List); ok {
			return fs, 1
		}
	}
	return nil, 0
}

func treeSetViewFromObject(obj *object.Object) (*treeSetView, *ghelpers.GErrBlk) {
	if obj == nil || object.IsNull(obj) {
		return nil, ghelpers.GetGErrBlk(excNames.NullPointerException, "TreeSet is null")
	}
	fld, ok := obj.FieldTable[treeSetStateField]
	if ok {
		if view, ok := fld.Fvalue.(*treeSetView); ok && view != nil && view.data != nil {
			return view, nil
		}
	}
	return nil, ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "object is not initialized as a TreeSet")
}

func treeSetSelf(params []any) (*object.Object, *treeSetView, int, *ghelpers.GErrBlk) {
	_, offset := treeSetContext(params)
	if len(params) <= offset {
		return nil, nil, offset, ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TreeSet method is missing its receiver")
	}
	obj, ok := params[offset].(*object.Object)
	if !ok || obj == nil {
		return nil, nil, offset, ghelpers.GetGErrBlk(excNames.NullPointerException, "TreeSet receiver is null or invalid")
	}
	view, err := treeSetViewFromObject(obj)
	return obj, view, offset, err
}

func treeSetWrap(view *treeSetView) *object.Object {
	className := treeSetClassName
	obj := object.MakeEmptyObjectWithClassName(&className)
	obj.FieldTable[treeSetStateField] = object.Field{Ftype: types.RawGoPointer, Fvalue: view}
	return obj
}

func treeSetInit(params []any) any {
	_, offset := treeSetContext(params)
	if len(params) <= offset {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TreeSet constructor is missing its receiver")
	}
	self, ok := params[offset].(*object.Object)
	if !ok || self == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TreeSet constructor receiver is null or invalid")
	}
	view := &treeSetView{data: &treeSetData{elements: make([]*object.Object, 0)}}
	if len(params) > offset+1 {
		arg, ok := params[offset+1].(*object.Object)
		if !ok {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TreeSet constructor argument is invalid")
		}
		if object.IsNull(arg) {
			// A null comparator selects natural ordering.
		} else if klass := object.GoStringFromStringPoolIndex(arg.KlassName); klass == treeSetClassName {
			if source, err := treeSetViewFromObject(arg); err == nil {
				view.data.comparator = source.data.comparator
				for _, elem := range treeSetVisibleElements(source) {
					added := treeSetInsert(view, elem, nil)
					if added.err != nil {
						return added.err
					}
				}
			} else {
				return err
			}
		} else {
			view.data.comparator = arg
		}
	}
	self.FieldTable[treeSetStateField] = object.Field{Ftype: types.RawGoPointer, Fvalue: view}
	return nil
}

func treeSetInitCollection(params []any) any {
	fs, offset := treeSetContext(params)
	if len(params) <= offset+1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TreeSet collection constructor is missing its collection")
	}
	collection, ok := params[offset+1].(*object.Object)
	if !ok || collection == nil || object.IsNull(collection) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TreeSet collection is null")
	}
	var values []*object.Object
	var sourceComparator *object.Object
	if view, err := treeSetViewFromObject(collection); err == nil {
		var gerr *ghelpers.GErrBlk
		values, gerr = treeSetVisibleElementsWithContext(view, fs)
		if gerr != nil {
			return gerr
		}
		sourceComparator = view.data.comparator
		if view.descending {
			sourceComparator = treeSetReverseComparator(sourceComparator)
		}
	} else {
		var gerr *ghelpers.GErrBlk
		values, gerr = treeSetCollectionElements(collection, fs)
		if gerr != nil {
			return gerr
		}
	}
	self, ok := params[offset].(*object.Object)
	if !ok || self == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "TreeSet constructor receiver is null or invalid")
	}
	view := &treeSetView{data: &treeSetData{elements: make([]*object.Object, 0), comparator: sourceComparator}}
	for _, elem := range values {
		result := treeSetInsert(view, elem, fs)
		if result.err != nil {
			return result.err
		}
	}
	self.FieldTable[treeSetStateField] = object.Field{Ftype: types.RawGoPointer, Fvalue: view}
	return nil
}

func treeSetCollectionElements(collection *object.Object, contexts ...*list.List) ([]*object.Object, *ghelpers.GErrBlk) {
	if collection == nil || object.IsNull(collection) {
		return nil, ghelpers.GetGErrBlk(excNames.NullPointerException, "TreeSet collection is null")
	}
	if view, err := treeSetViewFromObject(collection); err == nil {
		var fs *list.List
		if len(contexts) > 0 {
			fs = contexts[0]
		}
		return treeSetVisibleElementsWithContext(view, fs)
	}
	if fld, ok := collection.FieldTable["value"]; ok {
		switch values := fld.Fvalue.(type) {
		case []interface{}:
			out := make([]*object.Object, 0, len(values))
			for _, v := range values {
				obj, ok := v.(*object.Object)
				if !ok {
					return nil, ghelpers.GetGErrBlk(excNames.ClassCastException, "TreeSet collection contains a non-object")
				}
				out = append(out, obj)
			}
			return out, nil
		case []*object.Object:
			return append([]*object.Object(nil), values...), nil
		case *list.List:
			out := make([]*object.Object, 0, values.Len())
			for el := values.Front(); el != nil; el = el.Next() {
				obj, ok := el.Value.(*object.Object)
				if !ok {
					return nil, ghelpers.GetGErrBlk(excNames.ClassCastException, "TreeSet collection contains a non-object")
				}
				out = append(out, obj)
			}
			return out, nil
		}
	}
	if fld, ok := collection.FieldTable["map"]; ok {
		if hm, ok := fld.Fvalue.(types.DefHashMap); ok {
			out := make([]*object.Object, 0, len(hm))
			for _, value := range hm {
				obj, ok := value.(*object.Object)
				if !ok {
					return nil, ghelpers.GetGErrBlk(excNames.ClassCastException, "TreeSet collection contains a non-object")
				}
				out = append(out, obj)
			}
			return out, nil
		}
	}
	return nil, ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "unsupported TreeSet collection representation")
}

func treeSetIsNull(v *object.Object) bool { return v == nil || object.IsNull(v) }

func treeSetObjectsEqual(left, right *object.Object) bool {
	if treeSetIsNull(left) || treeSetIsNull(right) {
		return treeSetIsNull(left) && treeSetIsNull(right)
	}
	if left == right {
		return true
	}
	leftClass := object.GoStringFromStringPoolIndex(left.KlassName)
	if leftClass != object.GoStringFromStringPoolIndex(right.KlassName) {
		return false
	}
	if object.IsStringObject(left) {
		return object.GoStringFromStringObject(left) == object.GoStringFromStringObject(right)
	}
	leftValue, hasLeft := left.FieldTable["value"]
	rightValue, hasRight := right.FieldTable["value"]
	if !hasLeft || !hasRight {
		return false
	}
	switch x := leftValue.Fvalue.(type) {
	case int64:
		y, ok := rightValue.Fvalue.(int64)
		return ok && x == y
	case bool:
		y, ok := rightValue.Fvalue.(bool)
		return ok && x == y
	case float64:
		y, ok := rightValue.Fvalue.(float64)
		if !ok {
			return false
		}
		if math.IsNaN(x) && math.IsNaN(y) {
			return true
		}
		return x == y && (x != 0 || math.Signbit(x) == math.Signbit(y))
	default:
		return false
	}
}

func treeSetCompareJava(comparator, left, right *object.Object, fs *list.List) (int, *ghelpers.GErrBlk) {
	if comparator != nil && !object.IsNull(comparator) {
		if reversed, ok := comparator.FieldTable[treeSetReverseField]; ok && reversed.Fvalue == true {
			delegate, _ := comparator.FieldTable[treeSetDelegateField].Fvalue.(*object.Object)
			cmp, err := treeSetCompareJava(delegate, left, right, fs)
			return -cmp, err
		}
		global := globals.GetGlobalRef()
		className := object.GoStringFromStringPoolIndex(comparator.KlassName)
		if fs != nil && global != nil && global.FuncRunJavaFromG != nil {
			if result, invoked := loggingInvokeUserJavaMethod(fs, comparator, "compare", "(Ljava/lang/Object;Ljava/lang/Object;)I", left, right); invoked {
				return treeSetCompareResult(result, "comparator.compare")
			}
		}
		if global == nil || global.FuncInvokeGFunction == nil {
			return 0, ghelpers.GetGErrBlk(excNames.ClassCastException, "comparator callback is unavailable")
		}
		result := treeSetInvokeCallback(className+".compare(Ljava/lang/Object;Ljava/lang/Object;)I", []any{comparator, left, right}, fs)
		return treeSetCompareResult(result, "comparator.compare")
	}
	if treeSetIsNull(left) || treeSetIsNull(right) {
		return 0, ghelpers.GetGErrBlk(excNames.NullPointerException, "cannot compare null elements")
	}
	leftClass := object.GoStringFromStringPoolIndex(left.KlassName)
	rightClass := object.GoStringFromStringPoolIndex(right.KlassName)
	if leftClass != rightClass {
		return 0, ghelpers.GetGErrBlk(excNames.ClassCastException, fmt.Sprintf("cannot compare %s with %s", leftClass, rightClass))
	}
	if object.IsStringObject(left) {
		switch {
		case object.GoStringFromStringObject(left) < object.GoStringFromStringObject(right):
			return -1, nil
		case object.GoStringFromStringObject(left) > object.GoStringFromStringObject(right):
			return 1, nil
		default:
			return 0, nil
		}
	}
	if fld, ok := left.FieldTable["value"]; ok {
		other, exists := right.FieldTable["value"]
		if !exists {
			return 0, ghelpers.GetGErrBlk(excNames.ClassCastException, "Comparable element has no value")
		}
		switch x := fld.Fvalue.(type) {
		case int64:
			y, ok := other.Fvalue.(int64)
			if !ok {
				break
			}
			if x < y {
				return -1, nil
			}
			if x > y {
				return 1, nil
			}
			return 0, nil
		case float64:
			y, ok := other.Fvalue.(float64)
			if !ok {
				break
			}
			if math.IsNaN(x) {
				if math.IsNaN(y) {
					return 0, nil
				}
				return 1, nil
			}
			if math.IsNaN(y) {
				return -1, nil
			}
			if x == y {
				if x == 0 && math.Signbit(x) != math.Signbit(y) {
					if math.Signbit(x) {
						return -1, nil
					}
					return 1, nil
				}
				return 0, nil
			}
			if x < y {
				return -1, nil
			}
			return 1, nil
		case bool:
			y, ok := other.Fvalue.(bool)
			if ok {
				if x == y {
					return 0, nil
				}
				if !x {
					return -1, nil
				}
				return 1, nil
			}
		}
	}
	// Non-native Comparable classes are dispatched through the VM callback.
	global := globals.GetGlobalRef()
	if fs != nil && global != nil && global.FuncRunJavaFromG != nil {
		if result, invoked := loggingInvokeUserJavaMethod(fs, left, "compareTo", "(Ljava/lang/Object;)I", right); invoked {
			return treeSetCompareResult(result, "Comparable.compareTo")
		}
	}
	if global != nil && global.FuncInvokeGFunction != nil {
		result := treeSetInvokeCallback(leftClass+".compareTo(Ljava/lang/Object;)I", []any{left, right}, fs)
		if value, err := treeSetCompareResult(result, "Comparable.compareTo"); err == nil {
			return value, nil
		}
	}
	return 0, ghelpers.GetGErrBlk(excNames.ClassCastException, fmt.Sprintf("%s does not implement a supported natural ordering", leftClass))
}

func treeSetInvokeCallback(signature string, args []any, fs *list.List) any {
	global := globals.GetGlobalRef()
	if global == nil || global.FuncInvokeGFunction == nil {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "comparison callback is unavailable")
	}
	if method, ok := ghelpers.MethodSignatures[signature]; ok && method.NeedsContext {
		args = append([]any{fs}, args...)
	}
	return global.FuncInvokeGFunction(signature, args)
}

func treeSetCompareResult(result any, method string) (int, *ghelpers.GErrBlk) {
	if err, ok := result.(*ghelpers.GErrBlk); ok {
		return 0, err
	}
	if value, ok := result.(int64); ok {
		if value < 0 {
			return -1, nil
		}
		if value > 0 {
			return 1, nil
		}
		return 0, nil
	}
	if value, ok := result.(int); ok {
		if value < 0 {
			return -1, nil
		}
		if value > 0 {
			return 1, nil
		}
		return 0, nil
	}
	return 0, ghelpers.GetGErrBlk(excNames.ClassCastException, method+" returned an invalid comparison result")
}

func treeSetCompare(view *treeSetView, left, right *object.Object, fs *list.List) (int, *ghelpers.GErrBlk) {
	cmp, err := treeSetCompareJava(view.data.comparator, left, right, fs)
	if err == nil && view.descending {
		cmp = -cmp
	}
	return cmp, err
}

func treeSetVisible(view *treeSetView, value *object.Object, fs *list.List) (bool, *ghelpers.GErrBlk) {
	if view.hasLower {
		cmp, err := treeSetCompare(view, value, view.lower, fs)
		if err != nil {
			return false, err
		}
		if cmp < 0 || (cmp == 0 && !view.lowerInclusive) {
			return false, nil
		}
	}
	if view.hasUpper {
		cmp, err := treeSetCompare(view, value, view.upper, fs)
		if err != nil {
			return false, err
		}
		if cmp > 0 || (cmp == 0 && !view.upperInclusive) {
			return false, nil
		}
	}
	return true, nil
}

func treeSetVisibleElements(view *treeSetView) []*object.Object {
	values, _ := treeSetVisibleElementsWithContext(view, nil)
	return values
}

func treeSetVisibleElementsWithContext(view *treeSetView, fs *list.List) ([]*object.Object, *ghelpers.GErrBlk) {
	if view == nil || view.data == nil {
		return nil, nil
	}
	out := make([]*object.Object, 0, len(view.data.elements))
	for _, elem := range view.data.elements {
		visible, err := treeSetVisible(view, elem, fs)
		if err != nil {
			return nil, err
		}
		if visible {
			out = append(out, elem)
		}
	}
	if view.descending {
		for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
			out[i], out[j] = out[j], out[i]
		}
	}
	return out, nil
}

type treeSetInsertResult struct {
	added bool
	err   *ghelpers.GErrBlk
}

func treeSetInsert(view *treeSetView, elem *object.Object, fs *list.List) treeSetInsertResult {
	if treeSetIsNull(elem) {
		if view.data.comparator == nil || object.IsNull(view.data.comparator) {
			return treeSetInsertResult{err: ghelpers.GetGErrBlk(excNames.NullPointerException, "natural-order TreeSet does not permit null")}
		}
	}
	inside, err := treeSetVisible(view, elem, fs)
	if err != nil {
		return treeSetInsertResult{err: err}
	}
	if !inside {
		return treeSetInsertResult{err: ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "element is outside the TreeSet view range")}
	}
	if len(view.data.elements) == 0 {
		if _, err := treeSetCompareJava(view.data.comparator, elem, elem, fs); err != nil {
			return treeSetInsertResult{err: err}
		}
	}
	for i, existing := range view.data.elements {
		cmp, err := treeSetCompareJava(view.data.comparator, existing, elem, fs)
		if err != nil {
			return treeSetInsertResult{err: err}
		}
		if cmp == 0 {
			return treeSetInsertResult{}
		}
		if cmp > 0 {
			view.data.elements = append(view.data.elements, nil)
			copy(view.data.elements[i+1:], view.data.elements[i:])
			view.data.elements[i] = elem
			view.data.modCount++
			return treeSetInsertResult{added: true}
		}
	}
	view.data.elements = append(view.data.elements, elem)
	view.data.modCount++
	return treeSetInsertResult{added: true}
}

func treeSetFind(view *treeSetView, elem *object.Object, fs *list.List) (int, *ghelpers.GErrBlk) {
	for i, existing := range view.data.elements {
		cmp, err := treeSetCompareJava(view.data.comparator, existing, elem, fs)
		if err != nil {
			return -1, err
		}
		if cmp == 0 {
			return i, nil
		}
		if cmp > 0 {
			return -1, nil
		}
	}
	return -1, nil
}

func treeSetArg(params []any, offset, index int) (*object.Object, *ghelpers.GErrBlk) {
	if len(params) <= offset+index {
		return nil, ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "TreeSet method is missing an argument")
	}
	value, ok := params[offset+index].(*object.Object)
	if !ok {
		return nil, ghelpers.GetGErrBlk(excNames.ClassCastException, "TreeSet argument is not an object")
	}
	return value, nil
}

func treeSetAdd(params []any) any {
	fs, _ := treeSetContext(params)
	_, view, offset, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	elem, err := treeSetArg(params, offset, 1)
	if err != nil {
		return err
	}
	result := treeSetInsert(view, elem, fs)
	if result.err != nil {
		return result.err
	}
	if result.added {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func treeSetRemoveValue(view *treeSetView, elem *object.Object, fs *list.List) (bool, *ghelpers.GErrBlk) {
	if treeSetIsNull(elem) && (view.data.comparator == nil || object.IsNull(view.data.comparator)) {
		return false, nil
	}
	inside, err := treeSetVisible(view, elem, fs)
	if err != nil {
		return false, err
	}
	if !inside {
		return false, nil
	}
	index, err := treeSetFind(view, elem, fs)
	if err != nil {
		return false, err
	}
	if index < 0 {
		return false, nil
	}
	view.data.elements = append(view.data.elements[:index], view.data.elements[index+1:]...)
	view.data.modCount++
	return true, nil
}

func treeSetRemove(params []any) any {
	fs, _ := treeSetContext(params)
	_, view, offset, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	elem, err := treeSetArg(params, offset, 1)
	if err != nil {
		return err
	}
	removed, err := treeSetRemoveValue(view, elem, fs)
	if err != nil {
		return err
	}
	if removed {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func treeSetContains(params []any) any {
	fs, _ := treeSetContext(params)
	_, view, offset, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	elem, err := treeSetArg(params, offset, 1)
	if err != nil {
		return err
	}
	if treeSetIsNull(elem) && (view.data.comparator == nil || object.IsNull(view.data.comparator)) {
		return types.JavaBoolFalse
	}
	inside, err := treeSetVisible(view, elem, fs)
	if err != nil {
		return err
	}
	if !inside {
		return types.JavaBoolFalse
	}
	index, err := treeSetFind(view, elem, fs)
	if err != nil {
		return err
	}
	if index >= 0 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func treeSetAddAll(params []any) any    { return treeSetBulk(params, "add") }
func treeSetRemoveAll(params []any) any { return treeSetBulk(params, "remove") }
func treeSetRetainAll(params []any) any { return treeSetBulk(params, "retain") }

func treeSetBulk(params []any, mode string) any {
	fs, _ := treeSetContext(params)
	_, view, offset, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	collection, err := treeSetArg(params, offset, 1)
	if err != nil {
		return err
	}
	values, gerr := treeSetCollectionElements(collection, fs)
	if gerr != nil {
		return gerr
	}
	changed := false
	if mode == "add" {
		for _, elem := range values {
			result := treeSetInsert(view, elem, fs)
			if result.err != nil {
				return result.err
			}
			changed = changed || result.added
		}
		if changed {
			return types.JavaBoolTrue
		}
		return types.JavaBoolFalse
	}
	for i := len(view.data.elements) - 1; i >= 0; i-- {
		elem := view.data.elements[i]
		visible, gerr := treeSetVisible(view, elem, fs)
		if gerr != nil {
			return gerr
		}
		if !visible {
			continue
		}
		present := false
		for _, candidate := range values {
			if treeSetObjectsEqual(elem, candidate) {
				present = true
				break
			}
		}
		if mode == "remove" && present || mode == "retain" && !present {
			view.data.elements = append(view.data.elements[:i], view.data.elements[i+1:]...)
			view.data.modCount++
			changed = true
		}
	}
	if changed {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func treeSetContainsAll(params []any) any {
	fs, _ := treeSetContext(params)
	_, view, offset, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	collection, err := treeSetArg(params, offset, 1)
	if err != nil {
		return err
	}
	values, gerr := treeSetCollectionElements(collection, fs)
	if gerr != nil {
		return gerr
	}
	for _, elem := range values {
		inside, gerr := treeSetVisible(view, elem, fs)
		if gerr != nil {
			return gerr
		}
		if !inside {
			return types.JavaBoolFalse
		}
		idx, gerr := treeSetFind(view, elem, fs)
		if gerr != nil {
			return gerr
		}
		if idx < 0 {
			return types.JavaBoolFalse
		}
	}
	return types.JavaBoolTrue
}

func treeSetClear(params []any) any {
	fs, _ := treeSetContext(params)
	_, view, _, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	kept := make([]*object.Object, 0, len(view.data.elements))
	for _, elem := range view.data.elements {
		visible, gerr := treeSetVisible(view, elem, fs)
		if gerr != nil {
			return gerr
		}
		if visible {
			continue
		}
		kept = append(kept, elem)
	}
	if len(kept) != len(view.data.elements) {
		view.data.elements = kept
		view.data.modCount++
	}
	return nil
}

func treeSetSize(params []any) any {
	fs, _ := treeSetContext(params)
	_, view, _, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	values, gerr := treeSetVisibleElementsWithContext(view, fs)
	if gerr != nil {
		return gerr
	}
	return int64(len(values))
}

func treeSetIsEmpty(params []any) any {
	size := treeSetSize(params)
	if _, ok := size.(*ghelpers.GErrBlk); ok {
		return size
	}
	if size.(int64) == 0 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func treeSetToArray(params []any) any {
	fs, _ := treeSetContext(params)
	_, view, _, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	items, gerr := treeSetVisibleElementsWithContext(view, fs)
	if gerr != nil {
		return gerr
	}
	if len(params) > 1 {
		idx := 1
		if _, offset := treeSetContext(params); offset != 0 {
			idx = 2
		}
		if len(params) > idx {
			if arr, ok := params[idx].(*object.Object); ok && arr != nil && !object.IsNull(arr) {
				name := object.GoStringFromStringPoolIndex(arr.KlassName)
				return object.MakePrimitiveObject(name, name, items)
			}
		}
	}
	return object.MakePrimitiveObject(types.ObjectArrayClassName, types.ObjectArrayClassName, items)
}

func treeSetComparator(params []any) any {
	_, view, _, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	if !view.descending {
		if view.data.comparator == nil {
			return object.Null
		}
		return view.data.comparator
	}
	return treeSetReverseComparator(view.data.comparator)
}

func treeSetDescendingComparatorCompare(params []any) any {
	fs, offset := treeSetContext(params)
	if len(params) < offset+3 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "descending comparator is missing an argument")
	}
	comparator, ok := params[offset].(*object.Object)
	if !ok || comparator == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "descending comparator receiver is null")
	}
	left, leftOK := params[offset+1].(*object.Object)
	right, rightOK := params[offset+2].(*object.Object)
	if !leftOK || !rightOK {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "descending comparator arguments are invalid")
	}
	cmp, err := treeSetCompareJava(comparator, left, right, fs)
	if err != nil {
		return err
	}
	return int64(cmp)
}

func treeSetReverseComparator(comparator *object.Object) *object.Object {
	className := "java/util/TreeSet$DescendingComparator"
	obj := object.MakeEmptyObjectWithClassName(&className)
	obj.FieldTable[treeSetReverseField] = object.Field{Ftype: types.Bool, Fvalue: true}
	obj.FieldTable[treeSetDelegateField] = object.Field{Ftype: types.NonArrayObject, Fvalue: comparator}
	return obj
}

func treeSetBoundary(view *treeSetView, bound *object.Object, lower, inclusive bool, fs *list.List) (*treeSetView, *ghelpers.GErrBlk) {
	if treeSetIsNull(bound) && (view.data.comparator == nil || object.IsNull(view.data.comparator)) {
		return nil, ghelpers.GetGErrBlk(excNames.NullPointerException, "null range endpoint")
	}
	next := *view
	if lower {
		if next.hasLower {
			cmp, err := treeSetCompare(view, bound, next.lower, fs)
			if err != nil {
				return nil, err
			}
			if cmp < 0 || (cmp == 0 && inclusive && !next.lowerInclusive) {
				return nil, ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "range endpoint below view bounds")
			}
			if cmp > 0 {
				next.lower, next.lowerInclusive = bound, inclusive
			} else if cmp == 0 {
				next.lowerInclusive = next.lowerInclusive && inclusive
			}
		} else {
			next.lower, next.hasLower, next.lowerInclusive = bound, true, inclusive
		}
	} else {
		if next.hasUpper {
			cmp, err := treeSetCompare(view, bound, next.upper, fs)
			if err != nil {
				return nil, err
			}
			if cmp > 0 || (cmp == 0 && inclusive && !next.upperInclusive) {
				return nil, ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "range endpoint above view bounds")
			}
			if cmp < 0 {
				next.upper, next.upperInclusive = bound, inclusive
			} else if cmp == 0 {
				next.upperInclusive = next.upperInclusive && inclusive
			}
		} else {
			next.upper, next.hasUpper, next.upperInclusive = bound, true, inclusive
		}
	}
	if next.hasLower && next.hasUpper {
		cmp, err := treeSetCompare(view, next.lower, next.upper, fs)
		if err != nil {
			return nil, err
		}
		if cmp > 0 {
			return nil, ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "range lower endpoint exceeds upper endpoint")
		}
	}
	return &next, nil
}

func treeSetRangeArgs(params []any, lowIndex, highIndex int) (*object.Object, *treeSetView, int, *list.List, *ghelpers.GErrBlk) {
	fs, offset := treeSetContext(params)
	self, ok := params[offset].(*object.Object)
	if !ok {
		return nil, nil, offset, fs, ghelpers.GetGErrBlk(excNames.NullPointerException, "TreeSet receiver is invalid")
	}
	view, err := treeSetViewFromObject(self)
	if err != nil {
		return nil, nil, offset, fs, err
	}
	_ = lowIndex
	_ = highIndex
	return self, view, offset, fs, nil
}

func treeSetSubSet(params []any) any {
	_, view, offset, fs, err := treeSetRangeArgs(params, 1, 3)
	if err != nil {
		return err
	}
	from, err := treeSetArg(params, offset, 1)
	if err != nil {
		return err
	}
	fromInclusive, ok := params[offset+2].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "subSet lower inclusivity must be boolean")
	}
	to, err := treeSetArg(params, offset, 3)
	if err != nil {
		return err
	}
	toInclusive, ok := params[offset+4].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "subSet upper inclusivity must be boolean")
	}
	cmp, gerr := treeSetCompare(view, from, to, fs)
	if gerr != nil {
		return gerr
	}
	if cmp > 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "subSet lower endpoint exceeds upper endpoint")
	}
	low, gerr := treeSetBoundary(view, from, true, fromInclusive == types.JavaBoolTrue, fs)
	if gerr != nil {
		return gerr
	}
	bounded, gerr := treeSetBoundary(low, to, false, toInclusive == types.JavaBoolTrue, fs)
	if gerr != nil {
		return gerr
	}
	return treeSetWrap(bounded)
}

func treeSetSubSetInclusive(params []any) any {
	_, offset := treeSetContext(params)
	if len(params) <= offset+2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "subSet is missing endpoints")
	}
	return treeSetSubSet(append(append([]any(nil), params[:offset+2]...), int64(types.JavaBoolTrue), params[offset+2], int64(types.JavaBoolFalse)))
}

func treeSetHeadSet(params []any) any {
	fs, offset := treeSetContext(params)
	_, view, _, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	to, err := treeSetArg(params, offset, 1)
	if err != nil {
		return err
	}
	inclusive, ok := params[offset+2].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "headSet inclusivity must be boolean")
	}
	bounded, gerr := treeSetBoundary(view, to, false, inclusive == types.JavaBoolTrue, fs)
	if gerr != nil {
		return gerr
	}
	return treeSetWrap(bounded)
}

func treeSetHeadSetInclusive(params []any) any {
	return treeSetHeadSet(append(append([]any(nil), params...), int64(types.JavaBoolFalse)))
}

func treeSetTailSet(params []any) any {
	fs, offset := treeSetContext(params)
	_, view, _, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	from, err := treeSetArg(params, offset, 1)
	if err != nil {
		return err
	}
	inclusive, ok := params[offset+2].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "tailSet inclusivity must be boolean")
	}
	bounded, gerr := treeSetBoundary(view, from, true, inclusive == types.JavaBoolTrue, fs)
	if gerr != nil {
		return gerr
	}
	return treeSetWrap(bounded)
}

func treeSetTailSetInclusive(params []any) any {
	return treeSetTailSet(append(append([]any(nil), params...), int64(types.JavaBoolTrue)))
}

func treeSetDescendingSet(params []any) any {
	_, view, _, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	next := *view
	next.descending = !next.descending
	next.lower, next.upper = view.upper, view.lower
	next.hasLower, next.hasUpper = view.hasUpper, view.hasLower
	next.lowerInclusive, next.upperInclusive = view.upperInclusive, view.lowerInclusive
	return treeSetWrap(&next)
}

func treeSetExtreme(view *treeSetView, first bool, fs *list.List) (*object.Object, *ghelpers.GErrBlk) {
	values, err := treeSetVisibleElementsWithContext(view, fs)
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, ghelpers.GetGErrBlk(excNames.NoSuchElementException, "TreeSet is empty")
	}
	if !first {
		return values[len(values)-1], nil
	}
	return values[0], nil
}

func treeSetFirst(params []any) any {
	fs, _ := treeSetContext(params)
	_, view, _, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	value, gerr := treeSetExtreme(view, true, fs)
	if gerr != nil {
		return gerr
	}
	return value
}
func treeSetLast(params []any) any {
	fs, _ := treeSetContext(params)
	_, view, _, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	value, gerr := treeSetExtreme(view, false, fs)
	if gerr != nil {
		return gerr
	}
	return value
}

func treeSetNavigate(view *treeSetView, target *object.Object, relation int, fs *list.List) (*object.Object, *ghelpers.GErrBlk) {
	values, err := treeSetVisibleElementsWithContext(view, fs)
	if err != nil {
		return nil, err
	}
	var found *object.Object
	for _, elem := range values {
		cmp, err := treeSetCompare(view, elem, target, fs)
		if err != nil {
			return nil, err
		}
		matches := relation == -2 && cmp < 0 || relation == -1 && cmp <= 0 || relation == 1 && cmp >= 0 || relation == 2 && cmp > 0
		if matches {
			found = elem
			if relation < 0 {
				continue
			}
			return found, nil
		}
	}
	return found, nil
}

func treeSetNavigation(params []any, relation int) any {
	fs, _ := treeSetContext(params)
	_, view, offset, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	target, err := treeSetArg(params, offset, 1)
	if err != nil {
		return err
	}
	if treeSetIsNull(target) && (view.data.comparator == nil || object.IsNull(view.data.comparator)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "null navigation key")
	}
	value, gerr := treeSetNavigate(view, target, relation, fs)
	if gerr != nil {
		return gerr
	}
	if value == nil {
		return object.Null
	}
	return value
}
func treeSetLower(p []any) any   { return treeSetNavigation(p, -2) }
func treeSetFloor(p []any) any   { return treeSetNavigation(p, -1) }
func treeSetCeiling(p []any) any { return treeSetNavigation(p, 1) }
func treeSetHigher(p []any) any  { return treeSetNavigation(p, 2) }

func treeSetPoll(params []any, first bool) any {
	fs, _ := treeSetContext(params)
	_, view, _, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	value, gerr := treeSetExtreme(view, first, fs)
	if gerr != nil {
		if gerr.ExceptionType == excNames.NoSuchElementException {
			return object.Null
		}
		return gerr
	}
	removed, gerr := treeSetRemoveValue(view, value, fs)
	if gerr != nil {
		return gerr
	}
	if !removed {
		return object.Null
	}
	return value
}
func treeSetPollFirst(p []any) any { return treeSetPoll(p, true) }
func treeSetPollLast(p []any) any  { return treeSetPoll(p, false) }

func treeSetRemoveFirst(params []any) any { return treeSetRemoveEdge(params, true) }
func treeSetRemoveLast(params []any) any  { return treeSetRemoveEdge(params, false) }

func treeSetRemoveEdge(params []any, first bool) any {
	fs, _ := treeSetContext(params)
	_, view, _, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	value, gerr := treeSetExtreme(view, first, fs)
	if gerr != nil {
		return gerr
	}
	if _, gerr = treeSetRemoveValue(view, value, fs); gerr != nil {
		return gerr
	}
	return value
}

func treeSetIterator(params []any) any           { return treeSetMakeIterator(params, false) }
func treeSetDescendingIterator(params []any) any { return treeSetMakeIterator(params, true) }

func treeSetMakeIterator(params []any, reverse bool) any {
	fs, _ := treeSetContext(params)
	self, view, _, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	values, gerr := treeSetVisibleElementsWithContext(view, fs)
	if gerr != nil {
		return gerr
	}
	if reverse {
		for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
			values[i], values[j] = values[j], values[i]
		}
	}
	className := "java/util/Iterator"
	iterator := object.MakeEmptyObjectWithClassName(&className)
	iterator.FieldTable[iteratorCollectionField] = object.Field{Ftype: types.NonArrayObject, Fvalue: self}
	iterator.FieldTable[treeSetIteratorValues] = object.Field{Ftype: types.Ref, Fvalue: values}
	iterator.FieldTable[treeSetExpectedModCount] = object.Field{Ftype: types.Long, Fvalue: view.data.modCount}
	iterator.FieldTable[iteratorIndexField] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
	iterator.FieldTable[treeSetIteratorLast] = object.Field{Ftype: types.Int, Fvalue: int64(-1)}
	return iterator
}

func treeSetIteratorView(iterator *object.Object) (*treeSetView, bool) {
	collection, ok := iterator.FieldTable[iteratorCollectionField].Fvalue.(*object.Object)
	if !ok {
		return nil, false
	}
	view, err := treeSetViewFromObject(collection)
	return view, err == nil
}

func treeSetIteratorHasNext(iterator *object.Object) (bool, *ghelpers.GErrBlk) {
	view, ok := treeSetIteratorView(iterator)
	if !ok {
		return false, ghelpers.GetGErrBlk(excNames.IllegalStateException, "TreeSet iterator has no backing set")
	}
	expected, _ := iterator.FieldTable[treeSetExpectedModCount].Fvalue.(int64)
	if expected != view.data.modCount {
		return false, ghelpers.GetGErrBlk(excNames.ConcurrentModificationException, "TreeSet modified during iteration")
	}
	values, _ := iterator.FieldTable[treeSetIteratorValues].Fvalue.([]*object.Object)
	index, _ := iterator.FieldTable[iteratorIndexField].Fvalue.(int64)
	return index < int64(len(values)), nil
}

func treeSetIteratorNext(iterator *object.Object) any {
	hasNext, err := treeSetIteratorHasNext(iterator)
	if err != nil {
		return err
	}
	if !hasNext {
		return ghelpers.GetGErrBlk(excNames.NoSuchElementException, "TreeSet iterator has no next element")
	}
	values := iterator.FieldTable[treeSetIteratorValues].Fvalue.([]*object.Object)
	index := iterator.FieldTable[iteratorIndexField].Fvalue.(int64)
	iterator.FieldTable[iteratorIndexField] = object.Field{Ftype: types.Int, Fvalue: index + 1}
	iterator.FieldTable[treeSetIteratorLast] = object.Field{Ftype: types.Int, Fvalue: index}
	return values[index]
}

func treeSetIteratorRemove(iterator *object.Object) any {
	view, ok := treeSetIteratorView(iterator)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalStateException, "TreeSet iterator has no backing set")
	}
	expected, _ := iterator.FieldTable[treeSetExpectedModCount].Fvalue.(int64)
	if expected != view.data.modCount {
		return ghelpers.GetGErrBlk(excNames.ConcurrentModificationException, "TreeSet modified during iteration")
	}
	last := iterator.FieldTable[treeSetIteratorLast].Fvalue.(int64)
	if last < 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalStateException, "next() has not been called, or remove() has already been called")
	}
	values := iterator.FieldTable[treeSetIteratorValues].Fvalue.([]*object.Object)
	value := values[last]
	if _, err := treeSetRemoveValue(view, value, nil); err != nil {
		return err
	}
	iterator.FieldTable[treeSetExpectedModCount] = object.Field{Ftype: types.Long, Fvalue: view.data.modCount}
	iterator.FieldTable[treeSetIteratorLast] = object.Field{Ftype: types.Int, Fvalue: int64(-1)}
	return nil
}

func treeSetClone(params []any) any {
	fs, _ := treeSetContext(params)
	_, view, _, err := treeSetSelf(params)
	if err != nil {
		return err
	}
	comparator := view.data.comparator
	if view.descending {
		comparator = treeSetReverseComparator(comparator)
	}
	clone := &treeSetView{data: &treeSetData{elements: make([]*object.Object, 0), comparator: comparator}}
	values, gerr := treeSetVisibleElementsWithContext(view, fs)
	if gerr != nil {
		return gerr
	}
	for _, elem := range values {
		result := treeSetInsert(clone, elem, fs)
		if result.err != nil {
			return result.err
		}
	}
	return treeSetWrap(clone)
}
