/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"reflect"
	"testing"
)

func newConcurrentHashMapObj() *object.Object {
	return object.MakeEmptyObjectWithClassName(&classNameConcurrentHashMap)
}

func chmInit(t *testing.T, chm *object.Object) {
	t.Helper()
	if ret := concurrentHashMapInit([]interface{}{chm}); ret != nil {
		t.Fatalf("concurrentHashMapInit returned error: %v", ret)
	}
}

func TestLoad_Util_Concurrent_Hash_Map(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Util_Concurrent_Hash_Map()

	expectedMethods := map[string]struct {
		slots int
		trap  any
	}{
		"java/util/concurrent/ConcurrentHashMap.<clinit>()V":                                                                                              {0, ghelpers.ClinitGeneric},
		"java/util/concurrent/ConcurrentHashMap.<init>()V":                                                                                                {0, concurrentHashMapInit},
		"java/util/concurrent/ConcurrentHashMap.<init>(I)V":                                                                                               {1, concurrentHashMapInitCapacity},
		"java/util/concurrent/ConcurrentHashMap.<init>(IF)V":                                                                                              {2, concurrentHashMapInitCapacityLoadFactor},
		"java/util/concurrent/ConcurrentHashMap.<init>(IFI)V":                                                                                             {3, concurrentHashMapInitCapacityLoadFactorConcurrencyLevel},
		"java/util/concurrent/ConcurrentHashMap.<init>(Ljava/util/Map;)V":                                                                                 {1, concurrentHashMapInitMap},
		"java/util/concurrent/ConcurrentHashMap.newKeySet()Ljava/util/concurrent/ConcurrentHashMap$KeySetView;":                                           {0, concurrentHashMapNewKeySet},
		"java/util/concurrent/ConcurrentHashMap.newKeySet(I)Ljava/util/concurrent/ConcurrentHashMap$KeySetView;":                                          {1, concurrentHashMapNewKeySetCapacity},
		"java/util/concurrent/ConcurrentHashMap.clear()V":                                                                                                 {0, concurrentHashMapClear},
		"java/util/concurrent/ConcurrentHashMap.compute(Ljava/lang/Object;Ljava/util/function/BiFunction;)Ljava/lang/Object;":                             {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.computeIfAbsent(Ljava/lang/Object;Ljava/util/function/Function;)Ljava/lang/Object;":                       {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.computeIfPresent(Ljava/lang/Object;Ljava/util/function/BiFunction;)Ljava/lang/Object;":                    {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.contains(Ljava/lang/Object;)Z":                                                                            {1, concurrentHashMapContains},
		"java/util/concurrent/ConcurrentHashMap.containsKey(Ljava/lang/Object;)Z":                                                                         {1, concurrentHashMapContainsKey},
		"java/util/concurrent/ConcurrentHashMap.containsValue(Ljava/lang/Object;)Z":                                                                       {1, concurrentHashMapContainsValue},
		"java/util/concurrent/ConcurrentHashMap.elements()Ljava/util/Enumeration;":                                                                        {0, concurrentHashMapElements},
		"java/util/concurrent/ConcurrentHashMap.entrySet()Ljava/util/Set;":                                                                                {0, concurrentHashMapEntrySet},
		"java/util/concurrent/ConcurrentHashMap.equals(Ljava/lang/Object;)Z":                                                                              {1, concurrentHashMapEquals},
		"java/util/concurrent/ConcurrentHashMap.forEach(Ljava/util/function/BiConsumer;)V":                                                                {1, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.forEach(JLjava/util/function/BiConsumer;)V":                                                               {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.forEach(JLjava/util/function/BiFunction;Ljava/util/function/Consumer;)V":                                  {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.forEachEntry(JLjava/util/function/Consumer;)V":                                                            {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.forEachEntry(JLjava/util/function/Function;Ljava/util/function/Consumer;)V":                               {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.forEachKey(JLjava/util/function/Consumer;)V":                                                              {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.forEachKey(JLjava/util/function/Function;Ljava/util/function/Consumer;)V":                                 {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.forEachValue(JLjava/util/function/Consumer;)V":                                                            {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.forEachValue(JLjava/util/function/Function;Ljava/util/function/Consumer;)V":                               {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.get(Ljava/lang/Object;)Ljava/lang/Object;":                                                                {1, concurrentHashMapGet},
		"java/util/concurrent/ConcurrentHashMap.getOrDefault(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;":                                     {2, concurrentHashMapGetOrDefault},
		"java/util/concurrent/ConcurrentHashMap.hashCode()I":                                                                                              {0, concurrentHashMapHashCode},
		"java/util/concurrent/ConcurrentHashMap.isEmpty()Z":                                                                                               {0, concurrentHashMapIsEmpty},
		"java/util/concurrent/ConcurrentHashMap.keySet()Ljava/util/Set;":                                                                                  {0, concurrentHashMapKeySet},
		"java/util/concurrent/ConcurrentHashMap.keySet()Ljava/util/concurrent/ConcurrentHashMap$KeySetView;":                                              {0, concurrentHashMapKeySet},
		"java/util/concurrent/ConcurrentHashMap.keySet(Ljava/lang/Object;)Ljava/util/concurrent/ConcurrentHashMap$KeySetView;":                            {1, concurrentHashMapKeySetWithDefaultValue},
		"java/util/concurrent/ConcurrentHashMap.keys()Ljava/util/Enumeration;":                                                                            {0, concurrentHashMapKeys},
		"java/util/concurrent/ConcurrentHashMap.mappingCount()J":                                                                                          {0, concurrentHashMapMappingCount},
		"java/util/concurrent/ConcurrentHashMap.merge(Ljava/lang/Object;Ljava/lang/Object;Ljava/util/function/BiFunction;)Ljava/lang/Object;":             {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.put(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;":                                              {2, concurrentHashMapPut},
		"java/util/concurrent/ConcurrentHashMap.putAll(Ljava/util/Map;)V":                                                                                 {1, concurrentHashMapPutAll},
		"java/util/concurrent/ConcurrentHashMap.putIfAbsent(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;":                                      {2, concurrentHashMapPutIfAbsent},
		"java/util/concurrent/ConcurrentHashMap.reduce(JLjava/util/function/BiFunction;Ljava/util/function/BiFunction;)Ljava/lang/Object;":                {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceEntries(JLjava/util/function/BiFunction;)Ljava/util/Map$Entry;":                                     {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceEntries(JLjava/util/function/Function;Ljava/util/function/BiFunction;)Ljava/lang/Object;":           {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceEntriesToDouble(JLjava/util/function/ToDoubleFunction;DLjava/util/function/DoubleBinaryOperator;)D": {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceEntriesToInt(JLjava/util/function/ToIntFunction;ILjava/util/function/IntBinaryOperator;)I":          {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceEntriesToLong(JLjava/util/function/ToLongFunction;JLjava/util/function/LongBinaryOperator;)J":       {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceKeys(JLjava/util/function/BiFunction;)Ljava/lang/Object;":                                           {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceKeys(JLjava/util/function/Function;Ljava/util/function/BiFunction;)Ljava/lang/Object;":              {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceKeysToDouble(JLjava/util/function/ToDoubleFunction;DLjava/util/function/DoubleBinaryOperator;)D":    {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceKeysToInt(JLjava/util/function/ToIntFunction;ILjava/util/function/IntBinaryOperator;)I":             {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceKeysToLong(JLjava/util/function/ToLongFunction;JLjava/util/function/LongBinaryOperator;)J":          {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceToDouble(JLjava/util/function/ToDoubleBiFunction;DLjava/util/function/DoubleBinaryOperator;)D":      {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceToInt(JLjava/util/function/ToIntBiFunction;ILjava/util/function/IntBinaryOperator;)I":               {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceToLong(JLjava/util/function/ToLongBiFunction;JLjava/util/function/LongBinaryOperator;)J":            {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceValues(JLjava/util/function/BiFunction;)Ljava/lang/Object;":                                         {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceValues(JLjava/util/function/Function;Ljava/util/function/BiFunction;)Ljava/lang/Object;":            {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceValuesToDouble(JLjava/util/function/ToDoubleFunction;DLjava/util/function/DoubleBinaryOperator;)D":  {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceValuesToInt(JLjava/util/function/ToIntFunction;ILjava/util/function/IntBinaryOperator;)I":           {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceValuesToLong(JLjava/util/function/ToLongFunction;JLjava/util/function/LongBinaryOperator;)J":        {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.remove(Ljava/lang/Object;)Ljava/lang/Object;":                                                             {1, concurrentHashMapRemove},
		"java/util/concurrent/ConcurrentHashMap.remove(Ljava/lang/Object;Ljava/lang/Object;)Z":                                                            {2, concurrentHashMapRemoveEntry},
		"java/util/concurrent/ConcurrentHashMap.replace(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;":                                          {2, concurrentHashMapReplace},
		"java/util/concurrent/ConcurrentHashMap.replace(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Z":                                         {3, concurrentHashMapReplaceOldNew},
		"java/util/concurrent/ConcurrentHashMap.replaceAll(Ljava/util/function/BiFunction;)V":                                                             {1, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.search(JLjava/util/function/BiFunction;)Ljava/lang/Object;":                                               {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.searchEntries(JLjava/util/function/Function;)Ljava/lang/Object;":                                          {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.searchKeys(JLjava/util/function/Function;)Ljava/lang/Object;":                                             {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.searchValues(JLjava/util/function/Function;)Ljava/lang/Object;":                                           {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.size()I":                                                                                                  {0, concurrentHashMapSize},
		"java/util/concurrent/ConcurrentHashMap.toString()Ljava/lang/String;":                                                                             {0, concurrentHashMapToString},
		"java/util/concurrent/ConcurrentHashMap.values()Ljava/util/Collection;":                                                                           {0, concurrentHashMapValues},
	}

	if len(ghelpers.MethodSignatures) != len(expectedMethods) {
		t.Errorf("Expected %d registered methods, got %d", len(expectedMethods), len(ghelpers.MethodSignatures))
	}

	for sig, exp := range expectedMethods {
		gm, ok := ghelpers.MethodSignatures[sig]
		if !ok {
			t.Errorf("Missing registration for %s", sig)
			continue
		}
		if gm.ParamSlots != exp.slots {
			t.Errorf("%s: ParamSlots expected %d, got %d", sig, exp.slots, gm.ParamSlots)
		}
		if gm.GFunction == nil {
			t.Errorf("%s: GFunction is nil", sig)
		}
		if reflect.ValueOf(gm.GFunction).Pointer() != reflect.ValueOf(exp.trap).Pointer() {
			t.Errorf("%s: GFunction mismatch", sig)
		}
	}
}

func TestConcurrentHashMap_Constructors(t *testing.T) {
	globals.InitStringPool()

	// Default <init>()V
	chm1 := newConcurrentHashMapObj()
	if ret := concurrentHashMapInit([]interface{}{chm1}); ret != nil {
		t.Fatalf("concurrentHashMapInit returned error: %v", ret)
	}
	if sz := concurrentHashMapSize([]interface{}{chm1}).(int64); sz != 0 {
		t.Fatalf("expected size 0, got %d", sz)
	}

	// <init>(I)V with valid and negative capacity
	chm2 := newConcurrentHashMapObj()
	if ret := concurrentHashMapInitCapacity([]interface{}{chm2, int64(32)}); ret != nil {
		t.Fatalf("concurrentHashMapInitCapacity returned error: %v", ret)
	}
	if err := concurrentHashMapInitCapacity([]interface{}{chm2, int64(-5)}); err == nil {
		t.Fatalf("expected error for negative capacity")
	}

	// <init>(IF)V with valid, negative capacity, and invalid loadFactor
	chm3 := newConcurrentHashMapObj()
	if ret := concurrentHashMapInitCapacityLoadFactor([]interface{}{chm3, int64(16), 0.75}); ret != nil {
		t.Fatalf("concurrentHashMapInitCapacityLoadFactor returned error: %v", ret)
	}
	if err := concurrentHashMapInitCapacityLoadFactor([]interface{}{chm3, int64(-1), 0.75}); err == nil {
		t.Fatalf("expected error for negative capacity in InitCapacityLoadFactor")
	}
	if err := concurrentHashMapInitCapacityLoadFactor([]interface{}{chm3, int64(16), float64(0)}); err == nil {
		t.Fatalf("expected error for zero load factor in InitCapacityLoadFactor")
	}
	if err := concurrentHashMapInitCapacityLoadFactor([]interface{}{chm3, int64(16), -0.5}); err == nil {
		t.Fatalf("expected error for negative load factor in InitCapacityLoadFactor")
	}

	// <init>(IFI)V with valid and invalid concurrencyLevel
	chm4 := newConcurrentHashMapObj()
	if ret := concurrentHashMapInitCapacityLoadFactorConcurrencyLevel([]interface{}{chm4, int64(16), 0.75, int64(16)}); ret != nil {
		t.Fatalf("concurrentHashMapInitCapacityLoadFactorConcurrencyLevel returned error: %v", ret)
	}
	if err := concurrentHashMapInitCapacityLoadFactorConcurrencyLevel([]interface{}{chm4, int64(16), 0.75, int64(0)}); err == nil {
		t.Fatalf("expected error for zero concurrencyLevel")
	}

	// <init>(Ljava/util/Map;)V
	src := newConcurrentHashMapObj()
	chmInit(t, src)
	k := strKey("hello")
	v := object.StringObjectFromGoString("world")
	_ = concurrentHashMapPut([]interface{}{src, k, v})

	dst := newConcurrentHashMapObj()
	if ret := concurrentHashMapInitMap([]interface{}{dst, src}); ret != nil {
		t.Fatalf("concurrentHashMapInitMap returned error: %v", ret)
	}
	if sz := concurrentHashMapSize([]interface{}{dst}).(int64); sz != 1 {
		t.Fatalf("expected dst size 1, got %d", sz)
	}
	if got := concurrentHashMapGet([]interface{}{dst, k}).(*object.Object); object.GoStringFromStringObject(got) != "world" {
		t.Fatalf("expected 'world', got %v", object.GoStringFromStringObject(got))
	}

	// <init>(null Map) -> NullPointerException
	dstNull := newConcurrentHashMapObj()
	if err := concurrentHashMapInitMap([]interface{}{dstNull, object.Null}); err == nil {
		t.Fatalf("expected error for null map in InitMap")
	} else if geb, ok := err.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.NullPointerException {
		t.Fatalf("expected NullPointerException, got %v", err)
	}
}

func TestConcurrentHashMap_Put_Get_PutIfAbsent_Remove_Replace(t *testing.T) {
	globals.InitStringPool()

	chm := newConcurrentHashMapObj()
	chmInit(t, chm)

	k1 := strKey("key1")
	v1 := object.StringObjectFromGoString("val1")
	v2 := object.StringObjectFromGoString("val2")

	// 1. Initial put returns object.Null
	prev := concurrentHashMapPut([]interface{}{chm, k1, v1})
	if prev != object.Null {
		t.Fatalf("expected null previous value on initial put, got %v", prev)
	}
	if sz := concurrentHashMapSize([]interface{}{chm}).(int64); sz != 1 {
		t.Fatalf("expected size 1, got %d", sz)
	}

	// 2. get
	got := concurrentHashMapGet([]interface{}{chm, k1}).(*object.Object)
	if got != v1 {
		t.Fatalf("get returned different object than stored")
	}

	// 3. put overwrite returns previous value
	prev2 := concurrentHashMapPut([]interface{}{chm, k1, v2})
	if prev2 != v1 {
		t.Fatalf("expected v1 returned as previous, got %v", prev2)
	}
	if got2 := concurrentHashMapGet([]interface{}{chm, k1}).(*object.Object); got2 != v2 {
		t.Fatalf("expected v2 after overwrite")
	}

	// 4. putIfAbsent when key exists -> returns existing value, doesn't overwrite
	v3 := object.StringObjectFromGoString("val3")
	pia1 := concurrentHashMapPutIfAbsent([]interface{}{chm, k1, v3})
	if pia1 != v2 {
		t.Fatalf("expected existing value v2 on putIfAbsent, got %v", pia1)
	}
	if gotAfter := concurrentHashMapGet([]interface{}{chm, k1}).(*object.Object); gotAfter != v2 {
		t.Fatalf("putIfAbsent should not have overwritten v2")
	}

	// 5. putIfAbsent when key does not exist -> returns null, inserts
	k2 := strKey("key2")
	pia2 := concurrentHashMapPutIfAbsent([]interface{}{chm, k2, v3})
	if pia2 != object.Null {
		t.Fatalf("expected null on putIfAbsent for new key, got %v", pia2)
	}
	if sz := concurrentHashMapSize([]interface{}{chm}).(int64); sz != 2 {
		t.Fatalf("expected size 2, got %d", sz)
	}

	// 6. getOrDefault
	defVal := object.StringObjectFromGoString("default")
	if god1 := concurrentHashMapGetOrDefault([]interface{}{chm, k1, defVal}).(*object.Object); god1 != v2 {
		t.Fatalf("getOrDefault expected v2, got %v", god1)
	}
	if god2 := concurrentHashMapGetOrDefault([]interface{}{chm, strKey("absent"), defVal}).(*object.Object); god2 != defVal {
		t.Fatalf("getOrDefault expected defVal, got %v", god2)
	}

	// 7. replace(key, value)
	v4 := object.StringObjectFromGoString("val4")
	rep1 := concurrentHashMapReplace([]interface{}{chm, k1, v4})
	if rep1 != v2 {
		t.Fatalf("expected replace to return old value v2, got %v", rep1)
	}
	if gotRep := concurrentHashMapGet([]interface{}{chm, k1}).(*object.Object); gotRep != v4 {
		t.Fatalf("expected value to be updated to v4")
	}
	// replace on non-existent key returns null
	if repAbsent := concurrentHashMapReplace([]interface{}{chm, strKey("absent"), v4}); repAbsent != object.Null {
		t.Fatalf("expected null for replace on absent key, got %v", repAbsent)
	}

	// 8. replace(key, oldValue, newValue)
	v5 := object.StringObjectFromGoString("val5")
	// Wrong oldValue -> false
	if repCondFail := concurrentHashMapReplaceOldNew([]interface{}{chm, k1, object.StringObjectFromGoString("wrong"), v5}).(int64); repCondFail != types.JavaBoolFalse {
		t.Fatalf("expected replaceOldNew to fail with wrong old value")
	}
	// Correct oldValue -> true
	if repCondSuccess := concurrentHashMapReplaceOldNew([]interface{}{chm, k1, v4, v5}).(int64); repCondSuccess != types.JavaBoolTrue {
		t.Fatalf("expected replaceOldNew to succeed with correct old value")
	}
	if gotRep5 := concurrentHashMapGet([]interface{}{chm, k1}).(*object.Object); gotRep5 != v5 {
		t.Fatalf("expected value to be updated to v5")
	}

	// 9. remove(key, value) conditional remove
	// Wrong value -> false, not removed
	if remCondFail := concurrentHashMapRemoveEntry([]interface{}{chm, k1, object.StringObjectFromGoString("wrong")}).(int64); remCondFail != types.JavaBoolFalse {
		t.Fatalf("expected removeEntry to fail with wrong value")
	}
	if sz := concurrentHashMapSize([]interface{}{chm}).(int64); sz != 2 {
		t.Fatalf("expected size still 2, got %d", sz)
	}
	// Correct value -> true, removed
	if remCondSuccess := concurrentHashMapRemoveEntry([]interface{}{chm, k1, v5}).(int64); remCondSuccess != types.JavaBoolTrue {
		t.Fatalf("expected removeEntry to succeed with matching value")
	}
	if sz := concurrentHashMapSize([]interface{}{chm}).(int64); sz != 1 {
		t.Fatalf("expected size 1 after removeEntry, got %d", sz)
	}

	// 10. remove(key)
	remVal := concurrentHashMapRemove([]interface{}{chm, k2})
	if remVal != v3 {
		t.Fatalf("expected remove to return v3, got %v", remVal)
	}
	if sz := concurrentHashMapSize([]interface{}{chm}).(int64); sz != 0 {
		t.Fatalf("expected size 0 after all removes, got %d", sz)
	}
	if remAbsent := concurrentHashMapRemove([]interface{}{chm, k2}); remAbsent != object.Null {
		t.Fatalf("expected null for remove of absent key, got %v", remAbsent)
	}
}

func TestConcurrentHashMap_NullChecks(t *testing.T) {
	globals.InitStringPool()

	chm := newConcurrentHashMapObj()
	chmInit(t, chm)

	k := strKey("k")
	v := object.StringObjectFromGoString("v")

	checkNPE := func(name string, err interface{}) {
		t.Helper()
		if err == nil {
			t.Fatalf("%s: expected NullPointerException, got nil", name)
		}
		if geb, ok := err.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.NullPointerException {
			t.Fatalf("%s: expected NullPointerException, got %v", name, err)
		}
	}

	// put
	checkNPE("put(null, v)", concurrentHashMapPut([]interface{}{chm, object.Null, v}))
	checkNPE("put(k, null)", concurrentHashMapPut([]interface{}{chm, k, object.Null}))
	checkNPE("put(nil, v)", concurrentHashMapPut([]interface{}{chm, nil, v}))

	// putIfAbsent
	checkNPE("putIfAbsent(null, v)", concurrentHashMapPutIfAbsent([]interface{}{chm, object.Null, v}))
	checkNPE("putIfAbsent(k, null)", concurrentHashMapPutIfAbsent([]interface{}{chm, k, object.Null}))

	// get
	checkNPE("get(null)", concurrentHashMapGet([]interface{}{chm, object.Null}))

	// getOrDefault
	checkNPE("getOrDefault(null, default)", concurrentHashMapGetOrDefault([]interface{}{chm, object.Null, v}))

	// containsKey
	checkNPE("containsKey(null)", concurrentHashMapContainsKey([]interface{}{chm, object.Null}))

	// containsValue / contains
	checkNPE("containsValue(null)", concurrentHashMapContainsValue([]interface{}{chm, object.Null}))
	checkNPE("contains(null)", concurrentHashMapContains([]interface{}{chm, object.Null}))

	// remove
	checkNPE("remove(null)", concurrentHashMapRemove([]interface{}{chm, object.Null}))
	checkNPE("remove(null, v)", concurrentHashMapRemoveEntry([]interface{}{chm, object.Null, v}))
	checkNPE("remove(k, null)", concurrentHashMapRemoveEntry([]interface{}{chm, k, object.Null}))

	// replace
	checkNPE("replace(null, v)", concurrentHashMapReplace([]interface{}{chm, object.Null, v}))
	checkNPE("replace(k, null)", concurrentHashMapReplace([]interface{}{chm, k, object.Null}))
	checkNPE("replace(null, old, new)", concurrentHashMapReplaceOldNew([]interface{}{chm, object.Null, v, v}))
	checkNPE("replace(k, null, new)", concurrentHashMapReplaceOldNew([]interface{}{chm, k, object.Null, v}))
	checkNPE("replace(k, old, null)", concurrentHashMapReplaceOldNew([]interface{}{chm, k, v, object.Null}))

	// keySet(defaultVal)
	checkNPE("keySet(null)", concurrentHashMapKeySetWithDefaultValue([]interface{}{chm, object.Null}))

	// putAll(null)
	checkNPE("putAll(null)", concurrentHashMapPutAll([]interface{}{chm, object.Null}))
}

func TestConcurrentHashMap_Size_IsEmpty_Clear_MappingCount(t *testing.T) {
	globals.InitStringPool()

	chm := newConcurrentHashMapObj()
	chmInit(t, chm)

	if empty := concurrentHashMapIsEmpty([]interface{}{chm}).(int64); empty != types.JavaBoolTrue {
		t.Fatalf("expected isEmpty true")
	}
	if sz := concurrentHashMapSize([]interface{}{chm}).(int64); sz != 0 {
		t.Fatalf("expected size 0, got %d", sz)
	}
	if mc := concurrentHashMapMappingCount([]interface{}{chm}).(int64); mc != 0 {
		t.Fatalf("expected mappingCount 0, got %d", mc)
	}

	_ = concurrentHashMapPut([]interface{}{chm, strKey("a"), object.StringObjectFromGoString("1")})
	_ = concurrentHashMapPut([]interface{}{chm, strKey("b"), object.StringObjectFromGoString("2")})

	if empty := concurrentHashMapIsEmpty([]interface{}{chm}).(int64); empty != types.JavaBoolFalse {
		t.Fatalf("expected isEmpty false")
	}
	if sz := concurrentHashMapSize([]interface{}{chm}).(int64); sz != 2 {
		t.Fatalf("expected size 2, got %d", sz)
	}
	if mc := concurrentHashMapMappingCount([]interface{}{chm}).(int64); mc != 2 {
		t.Fatalf("expected mappingCount 2, got %d", mc)
	}

	// clear
	if ret := concurrentHashMapClear([]interface{}{chm}); ret != nil {
		t.Fatalf("clear returned error: %v", ret)
	}
	if empty := concurrentHashMapIsEmpty([]interface{}{chm}).(int64); empty != types.JavaBoolTrue {
		t.Fatalf("expected isEmpty true after clear")
	}
	if sz := concurrentHashMapSize([]interface{}{chm}).(int64); sz != 0 {
		t.Fatalf("expected size 0 after clear, got %d", sz)
	}
}

func TestConcurrentHashMap_Contains_ContainsKey_ContainsValue(t *testing.T) {
	globals.InitStringPool()

	chm := newConcurrentHashMapObj()
	chmInit(t, chm)

	k := strKey("apple")
	v := object.StringObjectFromGoString("red")
	_ = concurrentHashMapPut([]interface{}{chm, k, v})

	if ck := concurrentHashMapContainsKey([]interface{}{chm, k}).(int64); ck != types.JavaBoolTrue {
		t.Fatalf("containsKey expected true")
	}
	if ck := concurrentHashMapContainsKey([]interface{}{chm, strKey("banana")}).(int64); ck != types.JavaBoolFalse {
		t.Fatalf("containsKey expected false for absent key")
	}

	if cv := concurrentHashMapContainsValue([]interface{}{chm, v}).(int64); cv != types.JavaBoolTrue {
		t.Fatalf("containsValue expected true")
	}
	if cv := concurrentHashMapContainsValue([]interface{}{chm, object.StringObjectFromGoString("blue")}).(int64); cv != types.JavaBoolFalse {
		t.Fatalf("containsValue expected false for absent value")
	}

	if c := concurrentHashMapContains([]interface{}{chm, v}).(int64); c != types.JavaBoolTrue {
		t.Fatalf("contains expected true")
	}
}

func TestConcurrentHashMap_PutAll(t *testing.T) {
	globals.InitStringPool()

	src := newConcurrentHashMapObj()
	chmInit(t, src)
	_ = concurrentHashMapPut([]interface{}{src, strKey("k1"), object.StringObjectFromGoString("v1")})
	_ = concurrentHashMapPut([]interface{}{src, strKey("k2"), object.StringObjectFromGoString("v2")})

	dst := newConcurrentHashMapObj()
	chmInit(t, dst)
	_ = concurrentHashMapPut([]interface{}{dst, strKey("k0"), object.StringObjectFromGoString("v0")})

	if ret := concurrentHashMapPutAll([]interface{}{dst, src}); ret != nil {
		t.Fatalf("putAll returned error: %v", ret)
	}

	if sz := concurrentHashMapSize([]interface{}{dst}).(int64); sz != 3 {
		t.Fatalf("expected dst size 3, got %d", sz)
	}
	if got := concurrentHashMapGet([]interface{}{dst, strKey("k1")}).(*object.Object); object.GoStringFromStringObject(got) != "v1" {
		t.Fatalf("expected 'v1', got %v", object.GoStringFromStringObject(got))
	}
}

func TestConcurrentHashMap_Equals_HashCode_ToString(t *testing.T) {
	globals.InitStringPool()

	chm1 := newConcurrentHashMapObj()
	chmInit(t, chm1)
	_ = concurrentHashMapPut([]interface{}{chm1, strKey("a"), object.StringObjectFromGoString("1")})
	_ = concurrentHashMapPut([]interface{}{chm1, strKey("b"), object.StringObjectFromGoString("2")})

	chm2 := newConcurrentHashMapObj()
	chmInit(t, chm2)
	_ = concurrentHashMapPut([]interface{}{chm2, strKey("a"), object.StringObjectFromGoString("1")})
	_ = concurrentHashMapPut([]interface{}{chm2, strKey("b"), object.StringObjectFromGoString("2")})

	chm3 := newConcurrentHashMapObj()
	chmInit(t, chm3)
	_ = concurrentHashMapPut([]interface{}{chm3, strKey("a"), object.StringObjectFromGoString("different")})

	// equals
	if eq := concurrentHashMapEquals([]interface{}{chm1, chm1}).(int64); eq != types.JavaBoolTrue {
		t.Fatalf("expected equals self true")
	}
	if eq := concurrentHashMapEquals([]interface{}{chm1, chm2}).(int64); eq != types.JavaBoolTrue {
		t.Fatalf("expected equals chm2 true")
	}
	if eq := concurrentHashMapEquals([]interface{}{chm1, chm3}).(int64); eq != types.JavaBoolFalse {
		t.Fatalf("expected equals chm3 false")
	}
	if eq := concurrentHashMapEquals([]interface{}{chm1, object.Null}).(int64); eq != types.JavaBoolFalse {
		t.Fatalf("expected equals null false")
	}

	// hashCode
	h1 := concurrentHashMapHashCode([]interface{}{chm1}).(int64)
	h2 := concurrentHashMapHashCode([]interface{}{chm2}).(int64)
	if h1 != h2 {
		t.Fatalf("expected equal hash codes for equal maps, got %d and %d", h1, h2)
	}

	// toString
	emptyChm := newConcurrentHashMapObj()
	chmInit(t, emptyChm)
	emptyStr := concurrentHashMapToString([]interface{}{emptyChm}).(*object.Object)
	if s := object.GoStringFromStringObject(emptyStr); s != "{}" {
		t.Fatalf("expected '{}', got %s", s)
	}

	popStr := concurrentHashMapToString([]interface{}{chm1}).(*object.Object)
	if s := object.GoStringFromStringObject(popStr); s != "{a=1, b=2}" {
		t.Fatalf("expected '{a=1, b=2}', got %s", s)
	}
}

func TestConcurrentHashMap_KeySet_Values_EntrySet_Keys_Elements_NewKeySet(t *testing.T) {
	globals.InitStringPool()

	chm := newConcurrentHashMapObj()
	chmInit(t, chm)
	_ = concurrentHashMapPut([]interface{}{chm, strKey("k1"), object.StringObjectFromGoString("v1")})
	_ = concurrentHashMapPut([]interface{}{chm, strKey("k2"), object.StringObjectFromGoString("v2")})

	// keySet()
	keySetObj := concurrentHashMapKeySet([]interface{}{chm}).(*object.Object)
	if keySetObj == nil {
		t.Fatalf("keySet returned nil")
	}

	// keySet(defaultVal)
	defVal := object.StringObjectFromGoString("def")
	keySetViewObj := concurrentHashMapKeySetWithDefaultValue([]interface{}{chm, defVal}).(*object.Object)
	if keySetViewObj == nil {
		t.Fatalf("keySetWithDefaultValue returned nil")
	}

	// values()
	valuesObj := concurrentHashMapValues([]interface{}{chm}).(*object.Object)
	if valuesObj == nil {
		t.Fatalf("values returned nil")
	}

	// entrySet()
	entrySetObj := concurrentHashMapEntrySet([]interface{}{chm}).(*object.Object)
	if entrySetObj == nil {
		t.Fatalf("entrySet returned nil")
	}

	// keys()
	keysIter := concurrentHashMapKeys([]interface{}{chm}).(*object.Object)
	if keysIter == nil {
		t.Fatalf("keys returned nil")
	}
	if hasNext := iteratorHasNext([]interface{}{keysIter}).(int64); hasNext != types.JavaBoolTrue {
		t.Fatalf("keys iterator expected hasNext true")
	}

	// elements()
	elemIter := concurrentHashMapElements([]interface{}{chm}).(*object.Object)
	if elemIter == nil {
		t.Fatalf("elements returned nil")
	}
	if hasNext := iteratorHasNext([]interface{}{elemIter}).(int64); hasNext != types.JavaBoolTrue {
		t.Fatalf("elements iterator expected hasNext true")
	}

	// newKeySet()
	nks1 := concurrentHashMapNewKeySet([]interface{}{}).(*object.Object)
	if nks1 == nil {
		t.Fatalf("newKeySet returned nil")
	}

	// newKeySet(int)
	nks2 := concurrentHashMapNewKeySetCapacity([]interface{}{int64(10)}).(*object.Object)
	if nks2 == nil {
		t.Fatalf("newKeySet(int) returned nil")
	}
	if err := concurrentHashMapNewKeySetCapacity([]interface{}{int64(-1)}); err == nil {
		t.Fatalf("expected error for negative capacity in newKeySet")
	}
}

func TestConcurrentHashMap_Map_Delegation(t *testing.T) {
	globals.InitStringPool()

	chm := newConcurrentHashMapObj()
	chmInit(t, chm)

	k := strKey("delegationKey")
	v := object.StringObjectFromGoString("delegationVal")

	// mapPut
	ret := mapPut([]interface{}{chm, k, v})
	if ret != object.Null {
		t.Fatalf("expected null on mapPut, got %v", ret)
	}

	// mapGet
	got := mapGet([]interface{}{chm, k}).(*object.Object)
	if got != v {
		t.Fatalf("expected v from mapGet")
	}

	// mapSize
	if sz := mapSize([]interface{}{chm}).(int64); sz != 1 {
		t.Fatalf("expected size 1 from mapSize, got %d", sz)
	}

	// mapContainsKey
	if ck := mapContainsKey([]interface{}{chm, k}).(int64); ck != types.JavaBoolTrue {
		t.Fatalf("expected mapContainsKey true")
	}

	// mapIsEmpty
	if empty := mapIsEmpty([]interface{}{chm}).(int64); empty != types.JavaBoolFalse {
		t.Fatalf("expected mapIsEmpty false")
	}

	// mapRemove
	rem := mapRemove([]interface{}{chm, k}).(*object.Object)
	if rem != v {
		t.Fatalf("expected v from mapRemove")
	}
	if empty := mapIsEmpty([]interface{}{chm}).(int64); empty != types.JavaBoolTrue {
		t.Fatalf("expected mapIsEmpty true after remove")
	}

	// mapClear
	_ = mapPut([]interface{}{chm, k, v})
	_ = mapClear([]interface{}{chm})
	if sz := mapSize([]interface{}{chm}).(int64); sz != 0 {
		t.Fatalf("expected size 0 after mapClear, got %d", sz)
	}
}

func TestConcurrentHashMap_ErrorPaths(t *testing.T) {
	globals.InitStringPool()

	chm := newConcurrentHashMapObj()
	chmInit(t, chm)

	// Missing params
	if err := concurrentHashMapPut([]interface{}{chm, strKey("k")}); err == nil {
		t.Fatalf("expected error for missing param in put")
	}

	// First param not object
	if err := concurrentHashMapGet([]interface{}{int64(42), strKey("k")}); err == nil {
		t.Fatalf("expected error for non-object 'this'")
	}

	// First param wrong class
	wrongObj := object.MakeEmptyObjectWithClassName(&classNameHashMap)
	if err := concurrentHashMapSize([]interface{}{wrongObj}); err == nil {
		t.Fatalf("expected error for wrong class in size")
	}

	// Null this
	if err := concurrentHashMapSize([]interface{}{object.Null}); err == nil {
		t.Fatalf("expected error for null 'this'")
	}
}
