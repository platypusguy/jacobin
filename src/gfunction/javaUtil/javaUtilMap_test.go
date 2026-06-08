package javaUtil

import (
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
