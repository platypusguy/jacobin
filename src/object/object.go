/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-24 by Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package object

import (
	"jacobin/stringPool"
	"jacobin/types"
	"unsafe"
)

// This file contains basic functions of object creation. (Array objects
// are created in object\arrays.go.)

// With regard to the layout of a created object in Jacobin, note that
// on some architectures, but not Jacobin, there is an additional field
// that insures that the fields that follow the oops (the mark word and
// the class pointer) are aligned in memory for maximal performance.
type Object struct {
	Mark       MarkWord
	KlassName  uint32           // the index of the class name in the string pool
	FieldTable map[string]Field // map mapping field name to field
}

// These mark word contains values for different purposes. Here,
// we use the first four bytes for a hash value, which is taken
// from the address of the object. The 'misc' field will eventually
// contain other values, such as locking and monitoring items.
type MarkWord struct {
	Hash uint32 // contains hash code which is the lower 32 bits of the address
	Misc uint32 // at present unused
}

// We need to know the type of the field only to tell whether
// it occupies one or two slots on the stack when getfield and
// putfield bytecodes are executed. The type also flags static
// fields (with a leading X in the field type, which tells us
// to locate the value in the statics table.
type Field struct {
	Ftype  string // what type of value is stored in the field
	Fvalue any    // the actual value or a pointer to the value (ftype="[something)
}

// MakeEmptyObject() creates an empty basis Object. It is expected that other
// code will fill in the Klass header field and the data fields.
func MakeEmptyObject() *Object {
	o := Object{}
	h := uintptr(unsafe.Pointer(&o))
	o.Mark.Hash = uint32(h)
	o.KlassName = types.InvalidStringIndex // s/be filled in later, when class is filled in.

	// initialize the map of this object's fields
	o.FieldTable = make(map[string]Field)
	return &o
}

// MakeEmptyObjectWithClassName() creates an empty Object using the passed-in class name
func MakeEmptyObjectWithClassName(className *string) *Object {
	o := Object{}
	h := uintptr(unsafe.Pointer(&o))
	o.Mark.Hash = uint32(h)
	o.KlassName = stringPool.GetStringIndex(className)

	// initialize the map of this object's fields
	o.FieldTable = make(map[string]Field)
	return &o
}

// Make an object for a Java primitive field (byte, int, etc.), given the class and field type.
func MakePrimitiveObject(classString string, ftype string, arg any) *Object {
	objPtr := MakeEmptyObject()
	(*objPtr).KlassName = stringPool.GetStringIndex(&classString)
	field := Field{ftype, arg}
	(*objPtr).FieldTable["value"] = field
	return objPtr
}

// As of JACOBIN-613 this function is not used. It can be removed, once all string conversions to JavaBytes are done.
// // UpdateValueFieldFromBytes: Set the value field of the given object to the given byte array
// func UpdateValueFieldFromBytes(objPtr *Object, argBytes []byte) {
// 	fld := Field{Ftype: types.ByteArray, Fvalue: argBytes}
// 	objPtr.FieldTable["value"] = fld
// }

// UpdateValueFieldFromJavaBytes: Set the value field of the given object to the given JavaByte array
func UpdateValueFieldFromJavaBytes(objPtr *Object, argBytes []types.JavaByte) {
	fld := Field{Ftype: types.ByteArray, Fvalue: argBytes}
	objPtr.FieldTable["value"] = fld
}

// Null is the Jacobin implementation of Java's null
// var zero64 = uint64(0)
// var Null *Object = (*Object)(unsafe.Pointer(&zero64))
// JACOBIN-618 changed definition of null to this.
var Null = (*Object)(nil)

// determines whether a value is null or not
func IsNull(value any) bool {
	switch value.(type) {
	case *Object:
		obj := value.(*Object)
		return obj == nil || obj == Null
	}
	return value == nil
}

// Make a replica of an existing object.
func CloneObject(oldObject *Object) *Object {
	// Create new empty object.
	newObject := MakeEmptyObject()
	// Mimic the class.
	newObject.KlassName = oldObject.KlassName
	// Get a slice of keys from the old FieldTable.
	keys := make([]string, 0, len(oldObject.FieldTable))
	for key := range oldObject.FieldTable {
		keys = append(keys, key)
	}
	// For each key in the old FieldTable, copy that entry into the new FieldTable.
	for _, key := range keys {
		newObject.FieldTable[key] = oldObject.FieldTable[key]
	}
	return newObject
}

// Merge all the fields from a source object to a destination object.
func MergeFields(srcObj, dstObj *Object) {
	// Get a slice of keys from the old FieldTable.
	keys := make([]string, 0, len(srcObj.FieldTable))
	for key := range srcObj.FieldTable {
		keys = append(keys, key)
	}
	// For each key in the old FieldTable, copy that entry into the new FieldTable.
	for _, key := range keys {
		dstObj.FieldTable[key] = srcObj.FieldTable[key]
	}
}
