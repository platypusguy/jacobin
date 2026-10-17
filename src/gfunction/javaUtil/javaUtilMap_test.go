package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func TestMap_Delegation(t *testing.T) {
	globals.InitGlobals("test")
	Load_Util_Map()
	Load_Util_Hash_Map()

	// Create a HashMap object
	hm := object.MakeEmptyObject()
	hashmapInit([]any{hm})

	k := object.StringObjectFromGoString("key")
	v := object.StringObjectFromGoString("value")

	// Test mapPut via Map.put delegation
	ret := mapPut([]any{hm, k, v})
	if ret != object.Null {
		t.Errorf("Expected Null for new put, got %v", ret)
	}

	// Test mapGet
	ret = mapGet([]any{hm, k})
	if !object.IsStringObject(ret.(*object.Object)) || object.GoStringFromStringObject(ret.(*object.Object)) != "value" {
		t.Errorf("Expected 'value', got %v", ret)
	}

	// Test mapSize
	ret = mapSize([]any{hm})
	if ret.(int64) != 1 {
		t.Errorf("Expected size 1, got %v", ret)
	}

	// Test mapContainsKey
	ret = mapContainsKey([]any{hm, k})
	if ret.(int64) != types.JavaBoolTrue {
		t.Errorf("Expected true for containsKey")
	}

	// Test mapIsEmpty
	ret = mapIsEmpty([]any{hm})
	if ret.(int64) != types.JavaBoolFalse {
		t.Errorf("Expected false for isEmpty")
	}

	// Test mapRemove
	ret = mapRemove([]any{hm, k})
	if !object.IsStringObject(ret.(*object.Object)) || object.GoStringFromStringObject(ret.(*object.Object)) != "value" {
		t.Errorf("Expected 'value' for removed object")
	}

	// Test mapIsEmpty after remove
	ret = mapIsEmpty([]any{hm})
	if ret.(int64) != types.JavaBoolTrue {
		t.Errorf("Expected true for isEmpty after remove")
	}

	// Test mapClear
	mapPut([]any{hm, k, v})
	mapClear([]any{hm})
	if mapSize([]any{hm}).(int64) != 0 {
		t.Errorf("Expected size 0 after clear")
	}

	// Test mapPutIfAbsent
	mapPutIfAbsent([]any{hm, k, v})
	if mapSize([]any{hm}).(int64) != 1 {
		t.Errorf("Expected size 1 after putIfAbsent on empty map")
	}
	v2 := object.StringObjectFromGoString("value2")
	ret = mapPutIfAbsent([]any{hm, k, v2})
	if object.GoStringFromStringObject(ret.(*object.Object)) != "value" {
		t.Errorf("Expected original 'value' to be returned by putIfAbsent")
	}
	if object.GoStringFromStringObject(mapGet([]any{hm, k}).(*object.Object)) != "value" {
		t.Errorf("Expected 'value' to remain after putIfAbsent with existing key")
	}

	// Test mapGetOrDefault
	ret = mapGetOrDefault([]any{hm, k, v2})
	if object.GoStringFromStringObject(ret.(*object.Object)) != "value" {
		t.Errorf("Expected 'value' for existing key in getOrDefault")
	}
	k2 := object.StringObjectFromGoString("key2")
	ret = mapGetOrDefault([]any{hm, k2, v2})
	if object.GoStringFromStringObject(ret.(*object.Object)) != "value2" {
		t.Errorf("Expected 'value2' (default) for missing key in getOrDefault")
	}
}

func TestMap_EntrySet(t *testing.T) {
	globals.InitGlobals("test")
	Load_Util_Map()
	Load_Util_Hash_Map()

	// Create a HashMap object
	hm := object.MakeEmptyObjectWithClassName(&classNameHashMap)
	hashmapInit([]any{hm})

	k1 := object.StringObjectFromGoString("k1")
	v1 := object.StringObjectFromGoString("v1")
	k2 := object.StringObjectFromGoString("k2")
	v2 := object.StringObjectFromGoString("v2")

	mapPut([]any{hm, k1, v1})
	mapPut([]any{hm, k2, v2})

	// Get entrySet
	ret := mapEntrySet([]any{hm})
	if err, ok := ret.(*ghelpers.GErrBlk); ok {
		t.Fatalf("entrySet returned error: %s", err.ErrMsg)
	}
	es := ret.(*object.Object)
	if es == nil {
		t.Fatal("entrySet returned nil")
	}

	// entrySet is a HashSet (HashMap object in Jacobin)
	sz := hashmapSize([]any{es}).(int64)
	if sz != 2 {
		t.Errorf("Expected entrySet size 2, got %d", sz)
	}

	// The entrySet contains Map.Entry objects.
	// In our implementation, they are stored in the internal types.DefHashMap of the entrySet.
	fld := es.FieldTable[fieldNameMap]
	innerMap := fld.Fvalue.(types.DefHashMap)

	foundK1 := false
	foundK2 := false

	for _, entryVal := range innerMap {
		entryObj, ok := entryVal.(*object.Object)
		if !ok {
			t.Errorf("Expected *object.Object in entrySet, got %T", entryVal)
			continue
		}
		ek_any := mapEntryGetKey([]any{entryObj})
		ev_any := mapEntryGetValue([]any{entryObj})

		var ekStr, evStr string

		if s, ok := ek_any.(string); ok {
			ekStr = s
		} else if obj, ok := ek_any.(*object.Object); ok {
			ekStr = object.GoStringFromStringObject(obj)
		} else {
			t.Errorf("Unexpected key type: %T", ek_any)
		}

		if s, ok := ev_any.(string); ok {
			evStr = s
		} else if obj, ok := ev_any.(*object.Object); ok {
			evStr = object.GoStringFromStringObject(obj)
		} else {
			t.Errorf("Unexpected value type: %T", ev_any)
		}

		if ekStr == "k1" && evStr == "v1" {
			foundK1 = true
		} else if ekStr == "k2" && evStr == "v2" {
			foundK2 = true
		}
	}

	if !foundK1 {
		t.Errorf("Entry for k1 not found in entrySet")
	}
	if !foundK2 {
		t.Errorf("Entry for k2 not found in entrySet")
	}
}

func TestMapEntry_SetValue(t *testing.T) {
	globals.InitGlobals("test")
	Load_Util_Map()

	// Create a SimpleImmutableEntry object
	entryObj := object.MakeEmptyObjectWithClassName(new("java/util/AbstractMap$SimpleImmutableEntry"))
	entryObj.ThMutex.Lock()
	entryObj.FieldTable["key"] = object.Field{Ftype: "Ljava/lang/Object;", Fvalue: object.StringObjectFromGoString("key")}
	entryObj.FieldTable["value"] = object.Field{Ftype: "Ljava/lang/Object;", Fvalue: object.StringObjectFromGoString("oldValue")}
	entryObj.ThMutex.Unlock()

	newValue := object.StringObjectFromGoString("newValue")

	// Call setValue
	ret := mapEntrySetValue([]any{entryObj, newValue})

	// Verify old value returned
	oldValObj, ok := ret.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", ret)
	}
	if object.GoStringFromStringObject(oldValObj) != "oldValue" {
		t.Errorf("Expected 'oldValue', got '%s'", object.GoStringFromStringObject(oldValObj))
	}

	// Verify new value set
	val := mapEntryGetValue([]any{entryObj})
	valObj, ok := val.(*object.Object)
	if !ok {
		t.Fatalf("Expected *object.Object, got %T", val)
	}
	if object.GoStringFromStringObject(valObj) != "newValue" {
		t.Errorf("Expected 'newValue', got '%s'", object.GoStringFromStringObject(valObj))
	}
}

func TestMapEntry_Registration(t *testing.T) {
	globals.InitGlobals("test")
	Load_Util_Map()

	methods := []string{
		"java/util/Map.entry(Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/Map$Entry;",
		"java/util/Map.ofEntries([Ljava/util/Map$Entry;)Ljava/util/Map;",
		"java/util/Map$Entry.getKey()Ljava/lang/Object;",
		"java/util/Map$Entry.getValue()Ljava/lang/Object;",
		"java/util/Map$Entry.setValue(Ljava/lang/Object;)Ljava/lang/Object;",
		"java/util/Map$Entry.equals(Ljava/lang/Object;)Z",
		"java/util/Map$Entry.hashCode()I",
		"java/util/Map$Entry.comparingByKey()Ljava/util/Comparator;",
		"java/util/Map$Entry.comparingByValue()Ljava/util/Comparator;",
	}

	for _, m := range methods {
		if _, ok := ghelpers.MethodSignatures[m]; !ok {
			t.Errorf("Method %s not registered", m)
		}
	}
}

func TestMap_InvalidClass(t *testing.T) {
	globals.InitGlobals("test")

	// Create a non-Map object (e.g., a String)
	obj := object.StringObjectFromGoString("not a map")

	ret := mapPut([]any{obj, obj, obj})
	if _, ok := ret.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk for invalid class")
	}
}

func TestMap_KeySet(t *testing.T) {
	globals.InitGlobals("test")
	Load_Util_Map()
	Load_Util_Hash_Map()

	// Create a HashMap object
	hm := object.MakeEmptyObjectWithClassName(&classNameHashMap)
	hashmapInit([]any{hm})

	k1 := object.StringObjectFromGoString("k1")
	v1 := object.StringObjectFromGoString("v1")
	k2 := object.StringObjectFromGoString("k2")
	v2 := object.StringObjectFromGoString("v2")

	mapPut([]any{hm, k1, v1})
	mapPut([]any{hm, k2, v2})

	// Get keySet
	ret := mapKeySet([]any{hm})
	if err, ok := ret.(*ghelpers.GErrBlk); ok {
		t.Fatalf("keySet returned error: %s", err.ErrMsg)
	}
	ks := ret.(*object.Object)
	if ks == nil {
		t.Fatal("keySet returned nil")
	}

	// keySet is a HashSet (HashMap object in Jacobin)
	sz := hashmapSize([]any{ks}).(int64)
	if sz != 2 {
		t.Errorf("Expected keySet size 2, got %d", sz)
	}

	// Verify keys are present in the keySet
	if hashsetContains([]any{ks, k1}).(int64) != types.JavaBoolTrue {
		t.Errorf("k1 not found in keySet")
	}
	if hashsetContains([]any{ks, k2}).(int64) != types.JavaBoolTrue {
		t.Errorf("k2 not found in keySet")
	}
}

func TestMap_EntrySetInvoke(t *testing.T) {
	globals.InitGlobals("test")
	Load_Util_Map()
	Load_Util_Hash_Map()

	// Create a HashMap object
	hm := object.MakeEmptyObjectWithClassName(&classNameHashMap)
	hashmapInit([]any{hm})

	k1 := object.StringObjectFromGoString("k1")
	v1 := object.StringObjectFromGoString("v1")
	mapPut([]any{hm, k1, v1})

	// Invoke entrySet using ghelpers.Invoke
	ret := ghelpers.Invoke("java/util/Map.entrySet()Ljava/util/Set;", []interface{}{hm})
	if err, ok := ret.(*ghelpers.GErrBlk); ok {
		t.Fatalf("Invoke entrySet returned error: %s", err.ErrMsg)
	}

	es, ok := ret.(*object.Object)
	if !ok || es == nil {
		t.Fatal("Invoke entrySet did not return a valid object")
	}

	// Verify size of entrySet
	sz := hashmapSize([]any{es}).(int64)
	if sz != 1 {
		t.Errorf("Expected entrySet size 1, got %d", sz)
	}
}

func TestHashMap_EntrySetKeySetInvoke(t *testing.T) {
	globals.InitGlobals("test")
	Load_Util_Map()
	Load_Util_Hash_Map()

	// Create a HashMap object
	hm := object.MakeEmptyObjectWithClassName(&classNameHashMap)
	hashmapInit([]any{hm})

	k1 := object.StringObjectFromGoString("k1")
	v1 := object.StringObjectFromGoString("v1")
	mapPut([]any{hm, k1, v1})

	// Invoke HashMap.entrySet using ghelpers.Invoke
	ret := ghelpers.Invoke("java/util/HashMap.entrySet()Ljava/util/Set;", []interface{}{hm})
	if err, ok := ret.(*ghelpers.GErrBlk); ok {
		t.Fatalf("Invoke HashMap.entrySet returned error: %s", err.ErrMsg)
	}

	es, ok := ret.(*object.Object)
	if !ok || es == nil {
		t.Fatal("Invoke HashMap.entrySet did not return a valid object")
	}

	// Verify size of entrySet
	sz := hashmapSize([]any{es}).(int64)
	if sz != 1 {
		t.Errorf("Expected entrySet size 1, got %d", sz)
	}

	// Invoke HashMap.keySet using ghelpers.Invoke
	ret = ghelpers.Invoke("java/util/HashMap.keySet()Ljava/util/Set;", []interface{}{hm})
	if err, ok := ret.(*ghelpers.GErrBlk); ok {
		t.Fatalf("Invoke HashMap.keySet returned error: %s", err.ErrMsg)
	}

	ks, ok := ret.(*object.Object)
	if !ok || ks == nil {
		t.Fatal("Invoke HashMap.keySet did not return a valid object")
	}

	// Verify size of keySet
	sz = hashmapSize([]any{ks}).(int64)
	if sz != 1 {
		t.Errorf("Expected keySet size 1, got %d", sz)
	}
}

func TestMap_Entry(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()
	Load_Util_Map()

	// 1. Check registration
	gm, ok := ghelpers.MethodSignatures["java/util/Map.entry(Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/Map$Entry;"]
	if !ok {
		t.Fatalf("Method java/util/Map.entry not registered")
	}
	if gm.ParamSlots != 2 {
		t.Errorf("Expected ParamSlots 2, got %d", gm.ParamSlots)
	}
	if gm.GFunction == nil {
		t.Fatalf("Expected non-nil GFunction")
	}

	// 2. Test valid entry creation via direct function
	k := object.StringObjectFromGoString("testKey")
	v := object.StringObjectFromGoString("testVal")

	res := mapEntry([]any{k, v})
	entryObj, ok := res.(*object.Object)
	if !ok || entryObj == nil {
		t.Fatalf("mapEntry returned %T, want *object.Object", res)
	}

	className := object.GoStringFromStringPoolIndex(entryObj.KlassName)
	if className != "java/util/AbstractMap$SimpleImmutableEntry" {
		t.Errorf("Expected class java/util/AbstractMap$SimpleImmutableEntry, got %s", className)
	}

	// Check getKey and getValue
	gotK := mapEntryGetKey([]any{entryObj})
	if gotK != k {
		t.Errorf("getKey returned %v, want %v", gotK, k)
	}
	gotV := mapEntryGetValue([]any{entryObj})
	if gotV != v {
		t.Errorf("getValue returned %v, want %v", gotV, v)
	}

	// 3. Test invocation via ghelpers.Invoke
	invRes := ghelpers.Invoke("java/util/Map.entry(Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/Map$Entry;", []any{k, v})
	invObj, ok := invRes.(*object.Object)
	if !ok || invObj == nil {
		t.Fatalf("Invoke mapEntry returned %T, want *object.Object", invRes)
	}
	if mapEntryGetKey([]any{invObj}) != k || mapEntryGetValue([]any{invObj}) != v {
		t.Errorf("Invoked entry has incorrect key or value")
	}

	// 4. Test error handling: missing parameters (< 2)
	errRes := mapEntry([]any{})
	if errBlk, ok := errRes.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for 0 params, got %v", errRes)
	}
	errRes = mapEntry([]any{k})
	if errBlk, ok := errRes.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for 1 param, got %v", errRes)
	}

	// 5. Test error handling: null key or null value
	errRes = mapEntry([]any{nil, v})
	if errBlk, ok := errRes.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for nil key, got %v", errRes)
	}
	errRes = mapEntry([]any{k, nil})
	if errBlk, ok := errRes.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for nil value, got %v", errRes)
	}
	errRes = mapEntry([]any{object.Null, v})
	if errBlk, ok := errRes.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for Null key, got %v", errRes)
	}
	errRes = mapEntry([]any{k, object.Null})
	if errBlk, ok := errRes.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for Null value, got %v", errRes)
	}
}

func TestMap_OfEntries(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()
	Load_Util_Map()
	Load_Util_Hash_Map()

	// 1. Check registration
	gm, ok := ghelpers.MethodSignatures["java/util/Map.ofEntries([Ljava/util/Map$Entry;)Ljava/util/Map;"]
	if !ok {
		t.Fatalf("Method java/util/Map.ofEntries not registered")
	}
	if gm.ParamSlots != 1 {
		t.Errorf("Expected ParamSlots 1, got %d", gm.ParamSlots)
	}
	if gm.GFunction == nil {
		t.Fatalf("Expected non-nil GFunction")
	}

	// 2. Test empty entries array
	emptyArr := object.MakeArrayFromRawArray([]*object.Object{})
	emptyRes := mapOfEntries([]any{emptyArr})
	emptyMapObj, ok := emptyRes.(*object.Object)
	if !ok || emptyMapObj == nil {
		t.Fatalf("mapOfEntries returned %T for empty array, want *object.Object", emptyRes)
	}
	sz := mapSize([]any{emptyMapObj}).(int64)
	if sz != 0 {
		t.Errorf("Expected empty map size 0, got %d", sz)
	}

	// 3. Test non-empty entries array
	k1 := object.StringObjectFromGoString("k1")
	v1 := object.StringObjectFromGoString("v1")
	k2 := object.StringObjectFromGoString("k2")
	v2 := object.StringObjectFromGoString("v2")

	e1 := mapEntry([]any{k1, v1}).(*object.Object)
	e2 := mapEntry([]any{k2, v2}).(*object.Object)

	arr := object.MakeArrayFromRawArray([]*object.Object{e1, e2})
	res := mapOfEntries([]any{arr})
	mapObj, ok := res.(*object.Object)
	if !ok || mapObj == nil {
		t.Fatalf("mapOfEntries returned %T, want *object.Object", res)
	}

	sz = mapSize([]any{mapObj}).(int64)
	if sz != 2 {
		t.Errorf("Expected map size 2, got %d", sz)
	}

	gotV1 := mapGet([]any{mapObj, k1})
	if gotV1 != v1 {
		t.Errorf("mapGet for k1 returned %v, want %v", gotV1, v1)
	}
	gotV2 := mapGet([]any{mapObj, k2})
	if gotV2 != v2 {
		t.Errorf("mapGet for k2 returned %v, want %v", gotV2, v2)
	}

	// 4. Test invocation via ghelpers.Invoke
	invRes := ghelpers.Invoke("java/util/Map.ofEntries([Ljava/util/Map$Entry;)Ljava/util/Map;", []any{arr})
	invMapObj, ok := invRes.(*object.Object)
	if !ok || invMapObj == nil {
		t.Fatalf("Invoke mapOfEntries returned %T, want *object.Object", invRes)
	}
	if mapSize([]any{invMapObj}).(int64) != 2 {
		t.Errorf("Invoked map size != 2")
	}

	// 5. Test duplicate keys -> IllegalArgumentException
	e2Dup := mapEntry([]any{k1, v2}).(*object.Object)
	dupArr := object.MakeArrayFromRawArray([]*object.Object{e1, e2Dup})
	dupRes := mapOfEntries([]any{dupArr})
	if errBlk, ok := dupRes.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for duplicate keys, got %v", dupRes)
	}

	// 6. Test null array / missing arguments -> NullPointerException
	errRes := mapOfEntries([]any{})
	if errBlk, ok := errRes.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for 0 params, got %v", errRes)
	}
	errRes = mapOfEntries([]any{nil})
	if errBlk, ok := errRes.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for nil param, got %v", errRes)
	}
	errRes = mapOfEntries([]any{object.Null})
	if errBlk, ok := errRes.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for Null param, got %v", errRes)
	}

	// 7. Test null entry in array -> NullPointerException
	nullEntryArr := object.MakeArrayFromRawArray([]*object.Object{e1, nil})
	nullRes := mapOfEntries([]any{nullEntryArr})
	if errBlk, ok := nullRes.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for null entry in array, got %v", nullRes)
	}

	nullEntryArr2 := object.MakeArrayFromRawArray([]*object.Object{e1, object.Null})
	nullRes2 := mapOfEntries([]any{nullEntryArr2})
	if errBlk, ok := nullRes2.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for Null entry in array, got %v", nullRes2)
	}

	// 8. Test entry with null key or value -> NullPointerException
	badEntry := object.MakeEmptyObjectWithClassName(new("java/util/AbstractMap$SimpleImmutableEntry"))
	badEntry.FieldTable["key"] = object.Field{Ftype: "Ljava/lang/Object;", Fvalue: nil}
	badEntry.FieldTable["value"] = object.Field{Ftype: "Ljava/lang/Object;", Fvalue: v1}
	badArr := object.MakeArrayFromRawArray([]*object.Object{badEntry})
	badRes := mapOfEntries([]any{badArr})
	if errBlk, ok := badRes.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NullPointerException for entry with nil key, got %v", badRes)
	}
}
