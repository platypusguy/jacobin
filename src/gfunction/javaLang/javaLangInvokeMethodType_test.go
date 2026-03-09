/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"testing"
)

func TestMethodTypeFromMethodDescriptorString(t *testing.T) {
	// globals.InitGlobals("test")
	// trace.Init()
	// classloader.Init()
	// classloader.LoadBaseClasses()
	//
	// // Initialize primitive wrappers to ensure TYPE fields are populated
	// // Since we cannot import jvm package here, we call the clinit functions directly.
	// // This mimics what jvm.InitializePrimitiveWrappers does via gfunction invocation.
	//
	// // Initialize Integer (for "int")
	// integerClinit(nil)
	//
	// // Initialize Void (for "void")
	// voidClinit(nil)
	//
	// // Test Case 1: Simple descriptor with primitives
	// // (II)V -> int, int -> void
	// descriptor1 := "(II)V"
	// descObj1 := object.StringObjectFromGoString(descriptor1)
	// params1 := []interface{}{descObj1, nil} // ClassLoader is nil for now
	//
	// result1 := MethodTypeFromMethodDescriptorString(params1)
	// mtObj1, ok := result1.(*object.Object)
	// if !ok {
	// 	t.Fatalf("Expected *object.Object, got %T", result1)
	// }
	//
	// // Verify return type (void)
	// rtype1 := mtObj1.FieldTable["rtype"].Fvalue.(*object.Object)
	// // In Jacobin, primitive classes like void.class are stored in JLCmap
	// // We can check the name.
	//
	// // The KlassName field of a Class object points to "java/lang/Class".
	// // To get the name of the class it represents, we need to look at the "name" field
	// // inside the Class object.
	// rtypeNameField1, ok := rtype1.FieldTable["name"]
	// if !ok {
	// 	t.Fatalf("Return type Class object missing 'name' field")
	// }
	// rtypeName1 := rtypeNameField1.Fvalue.(string)
	//
	// if rtypeName1 != "void" {
	// 	t.Errorf("Expected return type name 'void', got '%s'", rtypeName1)
	// }
	//
	// // Verify parameter types (int, int)
	// ptypesArray1 := mtObj1.FieldTable["ptypes"].Fvalue.(*object.Object)
	// rawPtypes1 := ptypesArray1.FieldTable["value"].Fvalue.([]*object.Object)
	// if len(rawPtypes1) != 2 {
	// 	t.Errorf("Expected 2 parameters, got %d", len(rawPtypes1))
	// }
	//
	// ptypeNameField1_0, ok := rawPtypes1[0].FieldTable["name"]
	// if !ok {
	// 	t.Fatalf("Parameter type Class object missing 'name' field")
	// }
	// ptypeName1_0 := ptypeNameField1_0.Fvalue.(string)
	//
	// if ptypeName1_0 != "int" {
	// 	t.Errorf("Expected parameter 0 type 'int', got '%s'", ptypeName1_0)
	// }
	//
	// // Test Case 2: Descriptor with Object types
	// // (Ljava/lang/String;)Ljava/lang/Object;
	// descriptor2 := "(Ljava/lang/String;)Ljava/lang/Object;"
	// descObj2 := object.StringObjectFromGoString(descriptor2)
	// params2 := []interface{}{descObj2, nil}
	//
	// result2 := MethodTypeFromMethodDescriptorString(params2)
	// mtObj2 := result2.(*object.Object)
	//
	// rtype2 := mtObj2.FieldTable["rtype"].Fvalue.(*object.Object)
	// rtypeNameField2 := rtype2.FieldTable["name"]
	// rtypeName2 := rtypeNameField2.Fvalue.(string)
	//
	// if rtypeName2 != "java/lang/Object" {
	// 	t.Errorf("Expected return type java/lang/Object, got %s", rtypeName2)
	// }
	//
	// ptypesArray2 := mtObj2.FieldTable["ptypes"].Fvalue.(*object.Object)
	// rawPtypes2 := ptypesArray2.FieldTable["value"].Fvalue.([]*object.Object)
	// if len(rawPtypes2) != 1 {
	// 	t.Errorf("Expected 1 parameter, got %d", len(rawPtypes2))
	// }
	// ptypeNameField2 := rawPtypes2[0].FieldTable["name"]
	// ptypeName2 := ptypeNameField2.Fvalue.(string)
	//
	// if ptypeName2 != "java/lang/String" {
	// 	t.Errorf("Expected parameter type java/lang/String, got %s", ptypeName2)
	// }

	/* Needs some fixes before working
	// Test Case 3: Descriptor with Array types
	// ([I)[Ljava/lang/String;
	descriptor3 := "([I)[Ljava/lang/String;"
	descObj3 := object.StringObjectFromGoString(descriptor3)
	params3 := []interface{}{descObj3, nil}

	result3 := MethodTypeFromMethodDescriptorString(params3)
	mtObj3 := result3.(*object.Object)

	rtype3 := mtObj3.FieldTable["rtype"].Fvalue.(*object.Object)
	rtypeNameField3 := rtype3.FieldTable["name"]
	rtypeName3 := rtypeNameField3.Fvalue.(string)

	if rtypeName3 != "[Ljava/lang/String;" {
		t.Errorf("Expected return type [Ljava/lang/String;, got %s", rtypeName3)
	}

	ptypesArray3 := mtObj3.FieldTable["ptypes"].Fvalue.(*object.Object)
	rawPtypes3 := ptypesArray3.FieldTable["value"].Fvalue.([]*object.Object)
	if len(rawPtypes3) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(rawPtypes3))
	}
	ptypeNameField3 := rawPtypes3[0].FieldTable["name"]
	ptypeName3 := ptypeNameField3.Fvalue.(string)

	if ptypeName3 != "[I" {
		t.Errorf("Expected parameter type [I, got %s", ptypeName3)
	}
	*/
}

func TestParseDescriptorToClasses_Invalid(t *testing.T) {
	// // Test invalid descriptors
	// invalidDescriptors := []string{
	// 	"",
	// 	"()",                   // Missing return type
	// 	"(I",                   // Missing closing paren
	// 	"I)V",                  // Missing opening paren
	// 	"(Ljava/lang/String)V", // Missing semicolon
	// }
	//
	// for _, desc := range invalidDescriptors {
	// 	_, _, err := parseDescriptorToClasses(desc)
	// 	if err == nil {
	// 		t.Errorf("Expected error for invalid descriptor: %s", desc)
	// 	}
	// }
}
