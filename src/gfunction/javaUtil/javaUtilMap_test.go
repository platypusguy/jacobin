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

func TestMap_InvalidClass(t *testing.T) {
	globals.InitGlobals("test")

	// Create a non-Map object (e.g., a String)
	obj := object.StringObjectFromGoString("not a map")

	ret := mapPut([]any{obj, obj, obj})
	if _, ok := ret.(*ghelpers.GErrBlk); !ok {
		t.Errorf("Expected GErrBlk for invalid class")
	}
}
