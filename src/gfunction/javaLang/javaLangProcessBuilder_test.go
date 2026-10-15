/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"reflect"
	"testing"
)

func TestProcessBuilderCommandArray(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()
	Load_Lang_Process_Builder()

	sig := "java/lang/ProcessBuilder.command([Ljava/lang/String;)Ljava/lang/ProcessBuilder;"
	gmeth, exists := ghelpers.MethodSignatures[sig]
	if !exists {
		t.Fatalf("Method signature %s not registered", sig)
	}
	if gmeth.ParamSlots != 1 {
		t.Errorf("ParamSlots expected 1, got %d", gmeth.ParamSlots)
	}
	if reflect.ValueOf(gmeth.GFunction).Pointer() != reflect.ValueOf(processBuilderCommandArray).Pointer() {
		t.Errorf("GFunction does not point to processBuilderCommandArray")
	}

	// Test successful call with array object
	pbClassName := "java/lang/ProcessBuilder"
	pbObj := object.MakeEmptyObjectWithClassName(&pbClassName)

	str1 := object.StringObjectFromGoString("echo")
	str2 := object.StringObjectFromGoString("hello")
	arrObj := object.MakeArrayFromRawArray([]*object.Object{str1, str2})

	res := processBuilderCommandArray([]any{pbObj, arrObj})
	if res != pbObj {
		t.Fatalf("processBuilderCommandArray did not return self; got %v", res)
	}

	cmdField, ok := pbObj.FieldTable["command"]
	if !ok {
		t.Fatalf("pbObj missing 'command' field")
	}
	if cmdField.Ftype != "java/util/List" {
		t.Errorf("expected Ftype 'java/util/List', got %q", cmdField.Ftype)
	}
	listObj, ok := cmdField.Fvalue.(*object.Object)
	if !ok {
		t.Fatalf("command field value is not *object.Object, got %T", cmdField.Fvalue)
	}
	if className := object.GoStringFromStringPoolIndex(listObj.KlassName); className != "java/util/ArrayList" {
		t.Errorf("expected list class 'java/util/ArrayList', got %q", className)
	}

	valField, ok := listObj.FieldTable["value"]
	if !ok {
		t.Fatalf("listObj missing 'value' field")
	}
	elements, ok := valField.Fvalue.([]any)
	if !ok {
		t.Fatalf("listObj value is not []any, got %T", valField.Fvalue)
	}
	if len(elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(elements))
	}
	elem0, ok := elements[0].(*object.Object)
	if !ok || object.GoStringFromStringObject(elem0) != "echo" {
		t.Errorf("expected elements[0] == 'echo', got %v", elem0)
	}
	elem1, ok := elements[1].(*object.Object)
	if !ok || object.GoStringFromStringObject(elem1) != "hello" {
		t.Errorf("expected elements[1] == 'hello', got %v", elem1)
	}

	// Test with empty array
	emptyArrObj := object.MakeArrayFromRawArray([]*object.Object{})
	resEmpty := processBuilderCommandArray([]any{pbObj, emptyArrObj})
	if resEmpty != pbObj {
		t.Fatalf("processBuilderCommandArray failed on empty array")
	}
	valFieldEmpty := pbObj.FieldTable["command"].Fvalue.(*object.Object).FieldTable["value"]
	if len(valFieldEmpty.Fvalue.([]any)) != 0 {
		t.Errorf("expected 0 elements for empty array, got %d", len(valFieldEmpty.Fvalue.([]any)))
	}

	// Test with Go slice of strings
	resStrings := processBuilderCommandArray([]any{pbObj, []string{"ls", "-la"}})
	if resStrings != pbObj {
		t.Fatalf("processBuilderCommandArray failed on string slice")
	}
	valFieldStrings := pbObj.FieldTable["command"].Fvalue.(*object.Object).FieldTable["value"]
	if len(valFieldStrings.Fvalue.([]any)) != 2 {
		t.Errorf("expected 2 elements for string slice, got %d", len(valFieldStrings.Fvalue.([]any)))
	}
}

func TestProcessBuilderCommandArrayErrors(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()

	pbClassName := "java/lang/ProcessBuilder"
	pbObj := object.MakeEmptyObjectWithClassName(&pbClassName)

	// Self is nil
	res := processBuilderCommandArray([]any{nil, object.MakeArrayFromRawArray([]*object.Object{})})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for nil self, got %v", res)
	}

	// Self is object.Null
	res = processBuilderCommandArray([]any{object.Null, object.MakeArrayFromRawArray([]*object.Object{})})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for object.Null self, got %v", res)
	}

	// Command array is nil
	res = processBuilderCommandArray([]any{pbObj, nil})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for nil command array, got %v", res)
	}

	// Command array is object.Null
	res = processBuilderCommandArray([]any{pbObj, object.Null})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for object.Null command array, got %v", res)
	}

	// Command array containing nil element
	arrWithNil := object.MakeArrayFromRawArray([]*object.Object{object.StringObjectFromGoString("ls"), nil})
	res = processBuilderCommandArray([]any{pbObj, arrWithNil})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for array containing nil, got %v", res)
	}

	// Command array containing object.Null element
	arrWithNull := object.MakeArrayFromRawArray([]*object.Object{object.Null})
	res = processBuilderCommandArray([]any{pbObj, arrWithNull})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for array containing object.Null, got %v", res)
	}

	// Invalid self type
	res = processBuilderCommandArray([]any{"not-an-object", object.MakeArrayFromRawArray([]*object.Object{})})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException for invalid self, got %v", res)
	}

	// Invalid command argument type
	res = processBuilderCommandArray([]any{pbObj, 12345})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException for invalid command arg, got %v", res)
	}
}

func TestProcessBuilderStart(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()
	Load_Lang_Process_Builder()
	Load_Lang_Process()

	sig := "java/lang/ProcessBuilder.start()Ljava/lang/Process;"
	gmeth, exists := ghelpers.MethodSignatures[sig]
	if !exists {
		t.Fatalf("Method signature %s not registered", sig)
	}
	if gmeth.ParamSlots != 0 {
		t.Errorf("ParamSlots expected 0, got %d", gmeth.ParamSlots)
	}
	if reflect.ValueOf(gmeth.GFunction).Pointer() != reflect.ValueOf(processBuilderStart).Pointer() {
		t.Errorf("GFunction does not point to processBuilderStart")
	}

	pbClassName := "java/lang/ProcessBuilder"
	pbObj := object.MakeEmptyObjectWithClassName(&pbClassName)

	// Set command: ["echo", "test"]
	processBuilderCommandArray([]any{pbObj, []string{"echo", "test"}})

	res := processBuilderStart([]any{pbObj})
	procObj, ok := res.(*object.Object)
	if !ok {
		t.Fatalf("processBuilderStart did not return *object.Object; got %v", res)
	}
	if className := object.GoStringFromStringPoolIndex(procObj.KlassName); className != "java/lang/Process" {
		t.Errorf("expected class 'java/lang/Process', got %q", className)
	}

	pidVal, ok := procObj.FieldTable["pid"]
	if !ok {
		t.Fatalf("Process object missing 'pid' field")
	}
	pid, ok := pidVal.Fvalue.(int64)
	if !ok || pid <= 0 {
		t.Errorf("expected positive int64 pid, got %v", pidVal.Fvalue)
	}

	// Test waitFor on the started process
	waitRes := processWaitFor([]any{procObj})
	if exitCode, ok := waitRes.(int64); !ok || exitCode != 0 {
		t.Errorf("expected process exit code 0, got %v", waitRes)
	}

	// Test with working directory set
	pbObjWithDir := object.MakeEmptyObjectWithClassName(&pbClassName)
	processBuilderCommandArray([]any{pbObjWithDir, []string{"pwd"}})
	tempDir := t.TempDir()
	fileClassName := "java/io/File"
	fileObj := object.MakeEmptyObjectWithClassName(&fileClassName)
	fileObj.FieldTable[ghelpers.FilePath] = object.Field{
		Ftype:  types.JavaByteArray,
		Fvalue: object.JavaByteArrayFromGoString(tempDir),
	}
	pbObjWithDir.FieldTable["directory"] = object.Field{
		Ftype:  "java/io/File",
		Fvalue: fileObj,
	}

	resDir := processBuilderStart([]any{pbObjWithDir})
	procDirObj, ok := resDir.(*object.Object)
	if !ok {
		t.Fatalf("processBuilderStart with directory failed: %v", resDir)
	}
	waitResDir := processWaitFor([]any{procDirObj})
	if exitCode, ok := waitResDir.(int64); !ok || exitCode != 0 {
		t.Errorf("expected process with directory exit code 0, got %v", waitResDir)
	}
}

func TestProcessBuilderStartErrors(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()

	pbClassName := "java/lang/ProcessBuilder"
	pbObj := object.MakeEmptyObjectWithClassName(&pbClassName)

	// Null self
	res := processBuilderStart([]any{nil})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for nil self, got %v", res)
	}

	// object.Null self
	res = processBuilderStart([]any{object.Null})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for object.Null self, got %v", res)
	}

	// Invalid self type
	res = processBuilderStart([]any{"not-an-object"})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException for non-object self, got %v", res)
	}

	// Missing command field
	res = processBuilderStart([]any{pbObj})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for missing command, got %v", res)
	}

	// Empty command list -> IndexOutOfBoundsException
	pbObjEmpty := object.MakeEmptyObjectWithClassName(&pbClassName)
	processBuilderCommandArray([]any{pbObjEmpty, []string{}})
	res = processBuilderStart([]any{pbObjEmpty})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.IndexOutOfBoundsException {
		t.Errorf("expected IndexOutOfBoundsException for empty command list, got %v", res)
	}

	// Command list containing null
	pbObjNullElem := object.MakeEmptyObjectWithClassName(&pbClassName)
	listWithNull := object.MakePrimitiveObject("java/util/ArrayList", types.Ref, []any{object.Null})
	pbObjNullElem.FieldTable["command"] = object.Field{
		Ftype:  "java/util/List",
		Fvalue: listWithNull,
	}
	res = processBuilderStart([]any{pbObjNullElem})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for command list with null element, got %v", res)
	}

	// Non-existent executable -> IOException
	pbObjNonExistent := object.MakeEmptyObjectWithClassName(&pbClassName)
	processBuilderCommandArray([]any{pbObjNonExistent, []string{"/nonexistent/executable/binary/path"}})
	res = processBuilderStart([]any{pbObjNonExistent})
	if err, ok := res.(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.IOException {
		t.Errorf("expected IOException for non-existent executable, got %v", res)
	}
}

func TestProcessBuilderInitStringArray(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()
	Load_Lang_Process_Builder()

	sig := "java/lang/ProcessBuilder.<init>([Ljava/lang/String;)V"
	gmeth, exists := ghelpers.MethodSignatures[sig]
	if !exists {
		t.Fatalf("Method signature %s not registered", sig)
	}
	if gmeth.ParamSlots != 1 {
		t.Errorf("ParamSlots expected 1, got %d", gmeth.ParamSlots)
	}
	if reflect.ValueOf(gmeth.GFunction).Pointer() != reflect.ValueOf(processBuilderInitStringArray).Pointer() {
		t.Errorf("GFunction does not point to processBuilderInitStringArray")
	}

	pbClassName := "java/lang/ProcessBuilder"
	pbObj := object.MakeEmptyObjectWithClassName(&pbClassName)

	str1 := object.StringObjectFromGoString("ls")
	str2 := object.StringObjectFromGoString("-la")
	arrObj := object.MakeArrayFromRawArray([]*object.Object{str1, str2})

	res := processBuilderInitStringArray([]any{pbObj, arrObj})
	if res != nil {
		t.Fatalf("processBuilderInitStringArray expected nil, got %v", res)
	}

	cmdField, ok := pbObj.FieldTable["command"]
	if !ok {
		t.Fatalf("pbObj missing 'command' field after init")
	}
	if cmdField.Ftype != "java/util/List" {
		t.Errorf("expected Ftype 'java/util/List', got %q", cmdField.Ftype)
	}
	listObj, ok := cmdField.Fvalue.(*object.Object)
	if !ok {
		t.Fatalf("command field value is not *object.Object, got %T", cmdField.Fvalue)
	}
	valField, ok := listObj.FieldTable["value"]
	if !ok {
		t.Fatalf("listObj missing 'value' field")
	}
	elements, ok := valField.Fvalue.([]any)
	if !ok || len(elements) != 2 {
		t.Fatalf("expected 2 elements, got %v", valField.Fvalue)
	}

	// Test with Go slice of strings
	pbObj2 := object.MakeEmptyObjectWithClassName(&pbClassName)
	res2 := processBuilderInitStringArray([]any{pbObj2, []string{"echo", "hello"}})
	if res2 != nil {
		t.Fatalf("processBuilderInitStringArray expected nil, got %v", res2)
	}

	// Test error cases
	// Null self
	if err, ok := processBuilderInitStringArray([]any{nil, []string{"echo"}}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for nil self")
	}

	// Null array
	if err, ok := processBuilderInitStringArray([]any{pbObj, nil}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for nil command array")
	}

	// Null element in array
	arrWithNull := object.MakeArrayFromRawArray([]*object.Object{str1, nil})
	if err, ok := processBuilderInitStringArray([]any{pbObj, arrWithNull}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for null element in command array")
	}
}

func TestProcessBuilderInitList(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()
	Load_Lang_Process_Builder()

	sig := "java/lang/ProcessBuilder.<init>(Ljava/util/List;)V"
	gmeth, exists := ghelpers.MethodSignatures[sig]
	if !exists {
		t.Fatalf("Method signature %s not registered", sig)
	}
	if gmeth.ParamSlots != 1 {
		t.Errorf("ParamSlots expected 1, got %d", gmeth.ParamSlots)
	}
	if reflect.ValueOf(gmeth.GFunction).Pointer() != reflect.ValueOf(processBuilderInitList).Pointer() {
		t.Errorf("GFunction does not point to processBuilderInitList")
	}

	pbClassName := "java/lang/ProcessBuilder"
	pbObj := object.MakeEmptyObjectWithClassName(&pbClassName)

	str1 := object.StringObjectFromGoString("git")
	str2 := object.StringObjectFromGoString("status")
	listArg := object.MakePrimitiveObject("java/util/ArrayList", types.Ref, []any{str1, str2})

	res := processBuilderInitList([]any{pbObj, listArg})
	if res != nil {
		t.Fatalf("processBuilderInitList expected nil, got %v", res)
	}

	cmdField, ok := pbObj.FieldTable["command"]
	if !ok {
		t.Fatalf("pbObj missing 'command' field after init")
	}
	if cmdField.Ftype != "java/util/List" {
		t.Errorf("expected Ftype 'java/util/List', got %q", cmdField.Ftype)
	}
	listObj, ok := cmdField.Fvalue.(*object.Object)
	if !ok {
		t.Fatalf("command field value is not *object.Object, got %T", cmdField.Fvalue)
	}
	valField, ok := listObj.FieldTable["value"]
	if !ok {
		t.Fatalf("listObj missing 'value' field")
	}
	elements, ok := valField.Fvalue.([]any)
	if !ok || len(elements) != 2 {
		t.Fatalf("expected 2 elements, got %v", valField.Fvalue)
	}

	// Test with Go slice of strings
	pbObj2 := object.MakeEmptyObjectWithClassName(&pbClassName)
	res2 := processBuilderInitList([]any{pbObj2, []string{"echo", "hello"}})
	if res2 != nil {
		t.Fatalf("processBuilderInitList expected nil, got %v", res2)
	}

	// Test error cases
	// Null self
	if err, ok := processBuilderInitList([]any{nil, listArg}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for nil self")
	}

	// Null list
	if err, ok := processBuilderInitList([]any{pbObj, nil}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for nil command list")
	}

	// Null element in list
	listWithNull := object.MakePrimitiveObject("java/util/ArrayList", types.Ref, []any{str1, nil})
	if err, ok := processBuilderInitList([]any{pbObj, listWithNull}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for null element in command list")
	}
}

func TestProcessBuilderEnvironment(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()
	Load_Lang_Process_Builder()

	sig := "java/lang/ProcessBuilder.environment()Ljava/util/Map;"
	gmeth, exists := ghelpers.MethodSignatures[sig]
	if !exists {
		t.Fatalf("Method signature %s not registered", sig)
	}
	if gmeth.ParamSlots != 0 {
		t.Errorf("ParamSlots expected 0, got %d", gmeth.ParamSlots)
	}
	if reflect.ValueOf(gmeth.GFunction).Pointer() != reflect.ValueOf(processBuilderEnvironment).Pointer() {
		t.Errorf("GFunction does not point to processBuilderEnvironment")
	}

	// Null receiver
	if err, ok := processBuilderEnvironment([]any{nil}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for nil receiver")
	}

	pbClassName := "java/lang/ProcessBuilder"
	pbObj := object.MakeEmptyObjectWithClassName(&pbClassName)

	os.Setenv("JACOBIN_TEST_ENV_VAR", "JACOBIN_TEST_VAL")
	defer os.Unsetenv("JACOBIN_TEST_ENV_VAR")

	res := processBuilderEnvironment([]any{pbObj})
	mapObj, ok := res.(*object.Object)
	if !ok || mapObj == nil {
		t.Fatalf("expected *object.Object from environment(), got %T (%v)", res, res)
	}

	mapFld, ok := mapObj.FieldTable["map"]
	if !ok {
		t.Fatalf("expected 'map' field in HashMap object")
	}
	defMap, ok := mapFld.Fvalue.(types.DefHashMap)
	if !ok {
		t.Fatalf("expected types.DefHashMap in 'map' field, got %T", mapFld.Fvalue)
	}

	valObj, found := defMap["JACOBIN_TEST_ENV_VAR"]
	if !found {
		t.Fatalf("expected key 'JACOBIN_TEST_ENV_VAR' in environment map")
	}
	valStrObj, ok := valObj.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object value in environment map, got %T", valObj)
	}
	if val := object.GoStringFromStringObject(valStrObj); val != "JACOBIN_TEST_VAL" {
		t.Errorf("expected value 'JACOBIN_TEST_VAL', got %q", val)
	}

	// Calling environment() again returns the same map object
	res2 := processBuilderEnvironment([]any{pbObj})
	if res2 != mapObj {
		t.Errorf("expected environment() to return the same map instance, got %v vs %v", res2, mapObj)
	}

	// Modify the environment map and test processBuilderStart with it
	defMap["CUSTOM_VAR"] = object.StringObjectFromGoString("custom_value")
	processBuilderInitStringArray([]any{pbObj, []string{"sh", "-c", "echo $CUSTOM_VAR"}})

	startRes := processBuilderStart([]any{pbObj})
	procObj, ok := startRes.(*object.Object)
	if !ok || procObj == nil {
		t.Fatalf("processBuilderStart expected *object.Object, got %T (%v)", startRes, startRes)
	}
}

func TestProcessBuilderRedirectErrorStream(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()
	Load_Lang_Process_Builder()

	// Check signature redirectErrorStream(Z)Ljava/lang/ProcessBuilder;
	setSig := "java/lang/ProcessBuilder.redirectErrorStream(Z)Ljava/lang/ProcessBuilder;"
	setMeth, exists := ghelpers.MethodSignatures[setSig]
	if !exists {
		t.Fatalf("Method signature %s not registered", setSig)
	}
	if setMeth.ParamSlots != 1 {
		t.Errorf("ParamSlots expected 1, got %d", setMeth.ParamSlots)
	}
	if reflect.ValueOf(setMeth.GFunction).Pointer() != reflect.ValueOf(processBuilderRedirectErrorStreamSet).Pointer() {
		t.Errorf("GFunction does not point to processBuilderRedirectErrorStreamSet")
	}

	// Check signature redirectErrorStream()Z
	getSig := "java/lang/ProcessBuilder.redirectErrorStream()Z"
	getMeth, exists := ghelpers.MethodSignatures[getSig]
	if !exists {
		t.Fatalf("Method signature %s not registered", getSig)
	}
	if getMeth.ParamSlots != 0 {
		t.Errorf("ParamSlots expected 0, got %d", getMeth.ParamSlots)
	}
	if reflect.ValueOf(getMeth.GFunction).Pointer() != reflect.ValueOf(processBuilderRedirectErrorStreamGet).Pointer() {
		t.Errorf("GFunction does not point to processBuilderRedirectErrorStreamGet")
	}

	pbClassName := "java/lang/ProcessBuilder"
	pbObj := object.MakeEmptyObjectWithClassName(&pbClassName)

	// Default redirectErrorStream() should return false (0)
	getRes := processBuilderRedirectErrorStreamGet([]any{pbObj})
	if getRes != types.JavaBoolFalse {
		t.Errorf("expected default redirectErrorStream() to be JavaBoolFalse (0), got %v", getRes)
	}

	// Set redirectErrorStream(true) using Go bool
	setRes := processBuilderRedirectErrorStreamSet([]any{pbObj, true})
	if setRes != pbObj {
		t.Errorf("expected redirectErrorStream(true) to return self, got %v", setRes)
	}

	getRes = processBuilderRedirectErrorStreamGet([]any{pbObj})
	if getRes != types.JavaBoolTrue {
		t.Errorf("expected redirectErrorStream() to be JavaBoolTrue (1), got %v", getRes)
	}

	// Set redirectErrorStream(false) using int64(0)
	setRes = processBuilderRedirectErrorStreamSet([]any{pbObj, int64(0)})
	if setRes != pbObj {
		t.Errorf("expected redirectErrorStream(0) to return self, got %v", setRes)
	}

	getRes = processBuilderRedirectErrorStreamGet([]any{pbObj})
	if getRes != types.JavaBoolFalse {
		t.Errorf("expected redirectErrorStream() to be JavaBoolFalse (0), got %v", getRes)
	}

	// Set redirectErrorStream(true) using int64(1)
	setRes = processBuilderRedirectErrorStreamSet([]any{pbObj, int64(1)})
	if setRes != pbObj {
		t.Errorf("expected redirectErrorStream(1) to return self, got %v", setRes)
	}

	getRes = processBuilderRedirectErrorStreamGet([]any{pbObj})
	if getRes != types.JavaBoolTrue {
		t.Errorf("expected redirectErrorStream() to be JavaBoolTrue (1), got %v", getRes)
	}

	// Set redirectErrorStream using java/lang/Boolean object
	boolTrueObj := object.MakePrimitiveObject("java/lang/Boolean", types.Bool, types.JavaBoolTrue)
	setRes = processBuilderRedirectErrorStreamSet([]any{pbObj, boolTrueObj})
	if setRes != pbObj {
		t.Errorf("expected redirectErrorStream(Boolean(true)) to return self, got %v", setRes)
	}
	if getRes = processBuilderRedirectErrorStreamGet([]any{pbObj}); getRes != types.JavaBoolTrue {
		t.Errorf("expected redirectErrorStream() to be JavaBoolTrue (1), got %v", getRes)
	}

	boolFalseObj := object.MakePrimitiveObject("java/lang/Boolean", types.Bool, types.JavaBoolFalse)
	setRes = processBuilderRedirectErrorStreamSet([]any{pbObj, boolFalseObj})
	if setRes != pbObj {
		t.Errorf("expected redirectErrorStream(Boolean(false)) to return self, got %v", setRes)
	}
	if getRes = processBuilderRedirectErrorStreamGet([]any{pbObj}); getRes != types.JavaBoolFalse {
		t.Errorf("expected redirectErrorStream() to be JavaBoolFalse (0), got %v", getRes)
	}

	// Error test cases
	// Null self on get
	if err, ok := processBuilderRedirectErrorStreamGet([]any{nil}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for nil self on get")
	}

	// Null self on set
	if err, ok := processBuilderRedirectErrorStreamSet([]any{nil, true}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for nil self on set")
	}

	// Missing / null boolean param on set
	if err, ok := processBuilderRedirectErrorStreamSet([]any{pbObj}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for missing parameter on set")
	}
	if err, ok := processBuilderRedirectErrorStreamSet([]any{pbObj, nil}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for nil parameter on set")
	}

	// Invalid parameter type on set
	if err, ok := processBuilderRedirectErrorStreamSet([]any{pbObj, "invalid"}).(*ghelpers.GErrBlk); !ok || err.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException for invalid parameter type on set")
	}
}
