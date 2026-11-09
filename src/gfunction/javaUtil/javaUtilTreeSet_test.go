package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func newTreeSetForTest(t *testing.T, comparator *object.Object) *object.Object {
	t.Helper()
	set := makeTreeSetObject()
	var result any
	if comparator == nil {
		result = treeSetInit([]any{set})
	} else {
		result = treeSetInit([]any{set, comparator})
	}
	if result != nil {
		t.Fatalf("TreeSet constructor failed: %v", result)
	}
	return set
}

func makeTreeSetObject() *object.Object {
	className := treeSetClassName
	return object.MakeEmptyObjectWithClassName(&className)
}

func treeSetInt(value int64) *object.Object {
	return object.MakePrimitiveObject("java/lang/Integer", types.Int, value)
}

func treeSetValues(t *testing.T, set *object.Object) []int64 {
	t.Helper()
	array := treeSetToArray([]any{set}).(*object.Object)
	items, ok := array.FieldTable["value"].Fvalue.([]*object.Object)
	if !ok {
		t.Fatalf("TreeSet.toArray returned %T", array.FieldTable["value"].Fvalue)
	}
	values := make([]int64, len(items))
	for i, item := range items {
		values[i] = item.FieldTable["value"].Fvalue.(int64)
	}
	return values
}

func requireTreeSetError(t *testing.T, result any, exception int) {
	t.Helper()
	block, ok := result.(*ghelpers.GErrBlk)
	if !ok {
		t.Fatalf("expected GErrBlk with exception %d, got %T (%v)", exception, result, result)
	}
	if block.ExceptionType != exception {
		t.Fatalf("expected exception %d, got %d: %s", exception, block.ExceptionType, block.ErrMsg)
	}
}

func TestTreeSetNaturalOrderAndNavigation(t *testing.T) {
	globals.InitStringPool()
	set := newTreeSetForTest(t, nil)
	for _, value := range []int64{7, 2, 5, 2, -1} {
		if got := treeSetAdd([]any{set, treeSetInt(value)}); value == 2 && got == types.JavaBoolTrue {
			// The first 2 is new; the duplicate is checked independently below.
		}
	}
	if got := treeSetAdd([]any{set, treeSetInt(2)}); got != types.JavaBoolFalse {
		t.Fatalf("adding comparator-equal element should return false, got %v", got)
	}
	want := []int64{-1, 2, 5, 7}
	got := treeSetValues(t, set)
	if len(got) != len(want) {
		t.Fatalf("expected values %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected sorted values %v, got %v", want, got)
		}
	}
	checks := []struct {
		name string
		fn   func([]any) any
		key  int64
		want int64
	}{
		{"lower", treeSetLower, 5, 2},
		{"floor", treeSetFloor, 5, 5},
		{"ceiling", treeSetCeiling, 4, 5},
		{"higher", treeSetHigher, 5, 7},
	}
	for _, check := range checks {
		result := check.fn([]any{set, treeSetInt(check.key)})
		if value := result.(*object.Object).FieldTable["value"].Fvalue.(int64); value != check.want {
			t.Errorf("%s(%d) = %d, want %d", check.name, check.key, value, check.want)
		}
	}
	if got := treeSetSize([]any{set}); got != int64(4) {
		t.Fatalf("size = %v, want 4", got)
	}
}

func TestTreeSetRangeAndDescendingViewsAreBacked(t *testing.T) {
	globals.InitStringPool()
	set := newTreeSetForTest(t, nil)
	for _, value := range []int64{1, 2, 3, 4, 5} {
		treeSetAdd([]any{set, treeSetInt(value)})
	}
	view := treeSetSubSet([]any{set, treeSetInt(2), types.JavaBoolTrue, treeSetInt(5), types.JavaBoolFalse})
	if _, ok := view.(*ghelpers.GErrBlk); ok {
		t.Fatalf("subSet failed: %v", view)
	}
	viewObj := view.(*object.Object)
	if got, want := treeSetValues(t, viewObj), []int64{2, 3, 4}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Fatalf("subSet values = %v, want %v", got, want)
	}
	if got := treeSetAdd([]any{viewObj, treeSetInt(3)}); got != types.JavaBoolFalse {
		t.Fatalf("duplicate add through range view = %v, want false", got)
	}
	if got := treeSetAdd([]any{viewObj, treeSetInt(6)}); got == types.JavaBoolTrue {
		t.Fatal("range view accepted an element above its upper bound")
	}
	if got := treeSetAdd([]any{viewObj, treeSetInt(4)}); got != types.JavaBoolFalse {
		t.Fatalf("expected existing range element to remain present, got %v", got)
	}
	sortedView := treeSetSubSetInclusive([]any{set, treeSetInt(2), treeSetInt(5)}).(*object.Object)
	if got, want := treeSetValues(t, sortedView), []int64{2, 3, 4}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Fatalf("SortedSet subSet values = %v, want %v", got, want)
	}
	headView := treeSetHeadSetInclusive([]any{set, treeSetInt(4)}).(*object.Object)
	if got, want := treeSetValues(t, headView), []int64{1, 2, 3}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Fatalf("SortedSet headSet values = %v, want %v", got, want)
	}
	if got := treeSetRemove([]any{viewObj, treeSetInt(3)}); got != types.JavaBoolTrue {
		t.Fatalf("remove through view = %v, want true", got)
	}
	if got := treeSetContains([]any{set, treeSetInt(3)}); got != types.JavaBoolFalse {
		t.Fatal("removing from range view did not modify the backing set")
	}

	descending := treeSetDescendingSet([]any{set}).(*object.Object)
	if got, want := treeSetValues(t, descending), []int64{5, 4, 2, 1}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] || got[3] != want[3] {
		t.Fatalf("descending values = %v, want %v", got, want)
	}
	if got := treeSetFirst([]any{descending}).(*object.Object).FieldTable["value"].Fvalue.(int64); got != 5 {
		t.Fatalf("descending first = %d, want 5", got)
	}
	if got := treeSetAdd([]any{descending, treeSetInt(6)}); got != types.JavaBoolTrue {
		t.Fatalf("descending view add = %v, want true", got)
	}
	if got := treeSetContains([]any{set, treeSetInt(6)}); got != types.JavaBoolTrue {
		t.Fatal("descending view insertion did not update backing set")
	}
}

func TestTreeSetIteratorRemoveAndFailFast(t *testing.T) {
	globals.InitStringPool()
	set := newTreeSetForTest(t, nil)
	for _, value := range []int64{1, 2, 3} {
		treeSetAdd([]any{set, treeSetInt(value)})
	}
	iterator := treeSetIterator([]any{set}).(*object.Object)
	requireTreeSetError(t, iteratorRemove([]any{iterator}), excNames.IllegalStateException)
	if got := iteratorNext([]any{iterator}).(*object.Object).FieldTable["value"].Fvalue.(int64); got != 1 {
		t.Fatalf("first iterator value = %d, want 1", got)
	}
	if err := iteratorRemove([]any{iterator}); err != nil {
		t.Fatalf("iterator remove failed: %v", err)
	}
	if got := treeSetSize([]any{set}); got != int64(2) {
		t.Fatalf("size after iterator remove = %v, want 2", got)
	}
	if got := iteratorNext([]any{iterator}).(*object.Object).FieldTable["value"].Fvalue.(int64); got != 2 {
		t.Fatalf("iterator value after remove = %d, want 2", got)
	}

	other := treeSetIterator([]any{set}).(*object.Object)
	treeSetAdd([]any{set, treeSetInt(4)})
	requireTreeSetError(t, iteratorHasNext([]any{other}), excNames.ConcurrentModificationException)
	requireTreeSetError(t, iteratorNext([]any{other}), excNames.ConcurrentModificationException)
}

func TestTreeSetConstructorsAndErrors(t *testing.T) {
	globals.InitStringPool()
	set := newTreeSetForTest(t, nil)
	for _, value := range []int64{3, 1, 2} {
		treeSetAdd([]any{set, treeSetInt(value)})
	}
	collection := object.MakeEmptyObjectWithClassName(&classNameArrayList)
	collection.FieldTable["value"] = object.Field{Ftype: types.Ref, Fvalue: []interface{}{treeSetInt(8), treeSetInt(6), treeSetInt(8)}}
	copySet := makeTreeSetObject()
	if result := treeSetInitCollection([]any{copySet, collection}); result != nil {
		t.Fatalf("collection constructor failed: %v", result)
	}
	if got, want := treeSetValues(t, copySet), []int64{6, 8}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("collection constructor values = %v, want %v", got, want)
	}
	if result := treeSetInitCollection([]any{makeTreeSetObject(), object.Null}); result == nil {
		t.Fatal("null collection constructor should fail")
	} else {
		requireTreeSetError(t, result, excNames.NullPointerException)
	}
	requireTreeSetError(t, treeSetAdd([]any{set, object.Null}), excNames.NullPointerException)
	requireTreeSetError(t, treeSetAdd([]any{set, object.MakeEmptyObjectWithClassName(&classNameObject)}), excNames.ClassCastException)
	nonComparableSet := newTreeSetForTest(t, nil)
	requireTreeSetError(t, treeSetAdd([]any{nonComparableSet, object.MakeEmptyObjectWithClassName(&classNameObject)}), excNames.ClassCastException)
	if got := treeSetComparator([]any{set}); got != object.Null {
		t.Fatalf("natural-order comparator = %v, want null", got)
	}
	values := object.MakeEmptyObjectWithClassName(&classNameArrayList)
	values.FieldTable["value"] = object.Field{Ftype: types.Ref, Fvalue: []interface{}{treeSetInt(2)}}
	if got := treeSetRemoveAll([]any{set, values}); got != types.JavaBoolTrue {
		t.Fatalf("removeAll with an equal but distinct Integer returned %v", got)
	}
	if got := treeSetContains([]any{set, treeSetInt(2)}); got != types.JavaBoolFalse {
		t.Fatal("removeAll did not use Java object equality")
	}
}

func TestTreeSetCustomComparatorAndNullOrdering(t *testing.T) {
	globals.InitStringPool()
	global := globals.GetGlobalRef()
	oldInvoker := global.FuncInvokeGFunction
	t.Cleanup(func() { global.FuncInvokeGFunction = oldInvoker })
	global.FuncInvokeGFunction = func(signature string, params []any) any {
		if signature != "test/ReverseComparator.compare(Ljava/lang/Object;Ljava/lang/Object;)I" {
			t.Errorf("unexpected comparator method signature: %s", signature)
			return int64(0)
		}
		left, _ := params[1].(*object.Object)
		right, _ := params[2].(*object.Object)
		leftNull, rightNull := treeSetIsNull(left), treeSetIsNull(right)
		if leftNull || rightNull {
			if leftNull && rightNull {
				return int64(0)
			}
			if leftNull {
				return int64(-1)
			}
			return int64(1)
		}
		x := left.FieldTable["value"].Fvalue.(int64)
		y := right.FieldTable["value"].Fvalue.(int64)
		if x > y {
			return int64(-1)
		}
		if x < y {
			return int64(1)
		}
		return int64(0)
	}
	comparatorClass := "test/ReverseComparator"
	comparator := object.MakeEmptyObjectWithClassName(&comparatorClass)
	set := newTreeSetForTest(t, comparator)
	for _, value := range []int64{1, 3, 2} {
		if result := treeSetAdd([]any{set, treeSetInt(value)}); result != types.JavaBoolTrue {
			t.Fatalf("adding %d with custom comparator returned %v", value, result)
		}
	}
	if result := treeSetAdd([]any{set, object.Null}); result != types.JavaBoolTrue {
		t.Fatalf("custom comparator should be able to order null, got %v", result)
	}
	items := treeSetVisibleElements(treeSetViewFromObjectForTest(t, set))
	if len(items) != 4 || !object.IsNull(items[0]) {
		t.Fatalf("custom comparator did not order null first: %v", items)
	}
	for i, want := range []int64{3, 2, 1} {
		if got := items[i+1].FieldTable["value"].Fvalue.(int64); got != want {
			t.Fatalf("custom comparator order = %v, want [null 3 2 1]", items)
		}
	}
	reversed := treeSetDescendingSet([]any{set}).(*object.Object)
	reverseComparator := treeSetComparator([]any{reversed}).(*object.Object)
	if got := treeSetDescendingComparatorCompare([]any{reverseComparator, treeSetInt(1), treeSetInt(3)}); got != int64(-1) {
		t.Fatalf("descending comparator compare(1, 3) = %v, want -1", got)
	}
	if got := treeSetRemoveFirst([]any{set}); !object.IsNull(got.(*object.Object)) {
		t.Fatalf("removeFirst should remove null, got %v", got)
	}
	empty := newTreeSetForTest(t, comparator)
	requireTreeSetError(t, treeSetRemoveFirst([]any{empty}), excNames.NoSuchElementException)
}

func treeSetViewFromObjectForTest(t *testing.T, set *object.Object) *treeSetView {
	t.Helper()
	view, err := treeSetViewFromObject(set)
	if err != nil {
		t.Fatalf("failed to get TreeSet view: %v", err)
	}
	return view
}
