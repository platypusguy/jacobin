/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"unsafe"
)

// Implementation of java.lang.Enum, the common base class of all Java enum types.

func Load_Lang_Enum() {

	ghelpers.MethodSignatures["java/lang/Enum.<clinit>()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.ClinitGeneric}

	ghelpers.MethodSignatures["java/lang/Enum.<init>(Ljava/lang/String;I)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: enumInit}

	ghelpers.MethodSignatures["java/lang/Enum.clone()Ljava/lang/Object;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: enumClone}

	ghelpers.MethodSignatures["java/lang/Enum.compareTo(Ljava/lang/Enum;)I"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: enumCompareTo}

	ghelpers.MethodSignatures["java/lang/Enum.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: enumCompareTo}

	ghelpers.MethodSignatures["java/lang/Enum.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: enumEquals}

	ghelpers.MethodSignatures["java/lang/Enum.finalize()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapDeprecated}

	ghelpers.MethodSignatures["java/lang/Enum.getDeclaringClass()Ljava/lang/Class;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Enum.hashCode()I"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: enumHashCode}

	ghelpers.MethodSignatures["java/lang/Enum.name()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: enumName}

	ghelpers.MethodSignatures["java/lang/Enum.ordinal()I"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: enumOrdinal}

	ghelpers.MethodSignatures["java/lang/Enum.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: enumToString}

	ghelpers.MethodSignatures["java/lang/Enum.valueOf(Ljava/lang/Class;Ljava/lang/String;)Ljava/lang/Enum;"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
}

// "java/lang/Enum.<init>(Ljava/lang/String;I)V"
// Sets the name and ordinal fields on the enum instance, as done by the
// protected Enum(String name, int ordinal) constructor.
func enumInit(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		errMsg := fmt.Sprintf("enumInit: Invalid or null receiver: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	nameObj, ok := params[1].(*object.Object)
	if !ok || object.IsNull(nameObj) || !object.IsStringObject(nameObj) {
		errMsg := fmt.Sprintf("enumInit: Invalid name argument: %T", params[1])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	ordinal, ok := params[2].(int64)
	if !ok {
		errMsg := fmt.Sprintf("enumInit: Invalid ordinal argument: %T", params[2])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	self.FieldTable["name"] = object.Field{Ftype: types.StringClassRef, Fvalue: nameObj}
	self.FieldTable["ordinal"] = object.Field{Ftype: types.Int, Fvalue: ordinal}
	return nil
}

// "java/lang/Enum.name()Ljava/lang/String;"
func enumName(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		errMsg := fmt.Sprintf("enumName: Invalid or null receiver: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}
	fld, exists := self.FieldTable["name"]
	if !exists {
		errMsg := "enumName: missing name field"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	strObj, ok := fld.Fvalue.(*object.Object)
	if !ok || !object.IsStringObject(strObj) {
		errMsg := "enumName: name field is not a String"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	return strObj
}

// "java/lang/Enum.ordinal()I"
func enumOrdinal(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		errMsg := fmt.Sprintf("enumOrdinal: Invalid or null receiver: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}
	fld, exists := self.FieldTable["ordinal"]
	if !exists {
		errMsg := "enumOrdinal: missing ordinal field"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	ordinal, ok := fld.Fvalue.(int64)
	if !ok {
		errMsg := "enumOrdinal: ordinal field is not an int"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	return ordinal
}

// "java/lang/Enum.toString()Ljava/lang/String;"
// By default, returns the name of the enum constant, as declared in the enum declaration.
func enumToString(params []interface{}) interface{} {
	return enumName(params)
}

// "java/lang/Enum.equals(Ljava/lang/Object;)Z"
// Enum equality is identity-based: two enum constants are equal iff they're the same object.
func enumEquals(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		errMsg := fmt.Sprintf("enumEquals: Invalid or null receiver: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}
	other := params[1]
	if object.IsNull(other) {
		return types.JavaBoolFalse
	}
	otherObj, ok := other.(*object.Object)
	if !ok {
		return types.JavaBoolFalse
	}
	if self == otherObj {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// "java/lang/Enum.hashCode()I"
// Returns the identity hash code, as in Object.hashCode().
func enumHashCode(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		errMsg := fmt.Sprintf("enumHashCode: Invalid or null receiver: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}
	ptr := uintptr(unsafe.Pointer(self))
	hashCode := int64(ptr ^ (ptr >> 32))
	return hashCode
}

// "java/lang/Enum.compareTo(Ljava/lang/Enum;)I" and "compareTo(Ljava/lang/Object;)I"
// Compares this enum with the specified object for order, based on their ordinals.
// The two constants must belong to the same enum type.
func enumCompareTo(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		errMsg := fmt.Sprintf("enumCompareTo: Invalid or null receiver: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}
	other, ok := params[1].(*object.Object)
	if !ok || object.IsNull(other) {
		errMsg := fmt.Sprintf("enumCompareTo: Invalid or null argument: %T", params[1])
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	if self.KlassName != other.KlassName {
		errMsg := fmt.Sprintf("enumCompareTo: %s cannot be cast to %s",
			object.GoStringFromStringPoolIndex(other.KlassName),
			object.GoStringFromStringPoolIndex(self.KlassName))
		return ghelpers.GetGErrBlk(excNames.ClassCastException, errMsg)
	}

	selfOrdinalRes := enumOrdinal([]interface{}{self})
	if _, isErr := selfOrdinalRes.(*ghelpers.GErrBlk); isErr {
		return selfOrdinalRes
	}
	otherOrdinalRes := enumOrdinal([]interface{}{other})
	if _, isErr := otherOrdinalRes.(*ghelpers.GErrBlk); isErr {
		return otherOrdinalRes
	}

	selfOrdinal := selfOrdinalRes.(int64)
	otherOrdinal := otherOrdinalRes.(int64)
	return selfOrdinal - otherOrdinal
}

// "java/lang/Enum.clone()Ljava/lang/Object;"
// Enum.clone() is final and always throws CloneNotSupportedException,
// to guarantee that constants can never be cloned, preserving singleton semantics.
func enumClone([]interface{}) interface{} {
	errMsg := "enumClone: enum constants cannot be cloned"
	return ghelpers.GetGErrBlk(excNames.CloneNotSupportedException, errMsg)
}
