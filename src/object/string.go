/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package object

// This code attempts to codify the various kinds of strings and provide
// convenience methods to create and manipulate them.
// There exist only the following kinds of strings:
// 1) golang string -- this is commonly used inside the JVM
// 2) byte array -- less frequently used. Can represent either a string of bytes or a string of chars
// 3) String object -- used only when passing args to/from Java methods and gfunctions. Important note:
//    string objects are the *only* form of strings passed to/from Java methods and gfunctions.
//
// Implementation details:
// * the string pool stores only golang strings. This is done for performance reasons.
// * string objects' "value" field contains a byte array, which is required by Java methods and gfunctions.

import (
	"fmt"
	"jacobin/stringPool"
	"jacobin/trace"
	"jacobin/types"
	"strconv"
	"strings"
	"unsafe"
)

// NewStringObject creates an empty string object (aka Java String)
func NewStringObject() *Object {
	s := new(Object)
	s.Mark.Hash = 0
	s.KlassName = types.StringPoolStringIndex // =  java/lang/String
	s.FieldTable = make(map[string]Field)

	// ==== now the fields ====

	// value: the content of the string as array of runes or JavaBytes (int8)
	// Note: Post JDK9, this field is an array of bytes, so as to
	// enable compact strings.

	value := make([]types.JavaByte, 0)
	valueField := Field{Ftype: types.ByteArray, Fvalue: value} // empty string
	s.FieldTable["value"] = valueField

	// coder has two possible values:
	// LATIN(e.g., 0 = bytes for compact strings) or UTF16(e.g., 1 = UTF16)
	coderField := Field{Ftype: types.Byte, Fvalue: byte(0)}
	s.FieldTable["coder"] = coderField

	// the hash code, which is initialized to 0
	hash := Field{Ftype: types.Int, Fvalue: uint32(0)}
	s.FieldTable["hash"] = hash

	// hashIsZero: only true in rare case where compute hash is 0
	hashIsZero := Field{Ftype: types.Byte, Fvalue: byte(0)}
	s.FieldTable["hashIsZero"] = hashIsZero

	// The following static fields are preloaded in statics/LoadStaticsString()
	//   COMPACT_STRINGS (always true for JDK >= 9)
	//   UTF_8.INSTANCE ptr to encoder
	//   ISO_8859_1.INSTANCE ptr to encoder
	//   US_ASCII.INSTANCE ptr to encoder
	//   CodingErrorAction.REPLACE
	//   CASE_INSENSITIVE_ORDER
	//   serialPersistentFields

	return s
}

// StringObjectFromGoString: convenience method to create a string object from a Golang string
func StringObjectFromGoString(str string) *Object {
	newStr := NewStringObject()
	newStr.FieldTable["value"] = Field{Ftype: types.ByteArray, Fvalue: []byte(str)}
	return newStr
}

// GoStringFromStringObject: convenience method to extract a Go string from a String object (Java string)
func GoStringFromStringObject(obj *Object) string {
	if IsNull(obj) {
		return ""
	}
	fld, ok := obj.FieldTable["value"]
	if !ok {
		return ""
	}

	bytes := fld.Fvalue
	switch bytes.(type) {
	case []byte:
		return string(bytes.([]byte))
	case []types.JavaByte:
		return GoStringFromJavaByteArray(bytes.([]types.JavaByte))
	}

	return ""
}

// JavaByteFromStringObject: convenience method to extract a Java byte array from a String object (Java string)
func JavaByteFromStringObject(obj *Object) []types.JavaByte {
	if obj != nil && obj.KlassName == types.StringPoolStringIndex {
		return obj.FieldTable["value"].Fvalue.([]types.JavaByte)
	} else {
		return nil
	}
}

// ByteArrayFromStringObject: convenience method to extract a byte array from a String object (Java string)
func ByteArrayFromStringObject(obj *Object) []byte {
	if obj != nil && obj.KlassName == types.StringPoolStringIndex {
		return obj.FieldTable["value"].Fvalue.([]byte)
	} else {
		return nil
	}
}

// JavaByteArrayFromStringObject: convenience method to extract a JavaByte array from a String object (Java string)
func JavaByteArrayFromStringObject(obj *Object) []types.JavaByte {
	if obj != nil && obj.KlassName == types.StringPoolStringIndex {
		switch obj.FieldTable["value"].Fvalue.(type) {
		case []types.JavaByte:
			return obj.FieldTable["value"].Fvalue.([]types.JavaByte)
		case []byte:
			return JavaByteArrayFromGoByteArray(obj.FieldTable["value"].Fvalue.([]byte))
		}
	}
	return nil
}

// StringObjectFromByteArray: convenience method to create a string object from a byte array
func StringObjectFromByteArray(bytes []byte) *Object {
	newStr := NewStringObject()
	newStr.FieldTable["value"] = Field{Ftype: types.ByteArray, Fvalue: bytes}
	return newStr
}

// StringObjectFromJavaByteArray: convenience method to create a string object from a JavaByte array
func StringObjectFromJavaByteArray(bytes []types.JavaByte) *Object {
	newStr := NewStringObject()
	newStr.FieldTable["value"] = Field{Ftype: types.ByteArray, Fvalue: bytes}
	return newStr
}

// StringPoolIndexFromStringObject: convenience method to extract a string pool index from a String object
func StringPoolIndexFromStringObject(obj *Object) uint32 {
	if obj != nil && obj.KlassName == types.StringPoolStringIndex {
		str := string(obj.FieldTable["value"].Fvalue.([]byte))
		index := stringPool.GetStringIndex(&str)
		return index
	} else {
		return types.InvalidStringIndex
	}
}

// GoStringFromStringPoolIndex: convenience method to extract a Go string from a string pool index
func GoStringFromStringPoolIndex(index uint32) string {
	if index < stringPool.GetStringPoolSize() {
		return *stringPool.GetStringPointer(index)
	} else {
		return ""
	}
}

// StringObjectFromStringPoolIndex: convenience method to create a string object using a string pool index
func StringObjectFromPoolIndex(index uint32) *Object {
	if index < stringPool.GetStringPoolSize() {
		return StringObjectFromGoString(*stringPool.GetStringPointer(index))
	} else {
		return nil
	}
}

// ByteArrayFromStringPoolIndex: convenience method to get a byte array using a string pool index
func ByteArrayFromStringPoolIndex(index uint32) []byte {
	if index < stringPool.GetStringPoolSize() {
		return []byte(*stringPool.GetStringPointer(index))
	} else {
		return nil
	}
}

// IsStringObject determines whether an object is a string object
// (i.e., a Java string). It assumes that any object whose
// KlassName refers to java/lang/String is an instance of a Java string
func IsStringObject(unknown any) bool {
	if unknown == nil {
		return false
	}

	o, ok := unknown.(*Object)
	if !ok {
		return false
	}

	return o.KlassName == types.StringPoolStringIndex
}

func JavaByteSliceToGoString(vector types.JavaByteSlice) string {
	uint8Slice := *(*[]uint8)(unsafe.Pointer(&vector))
	return string(uint8Slice)
}

func GoStringToJavaByteSlice(arg string) types.JavaByteSlice {
	byteSlice := []byte(arg)
	javaByteSlice := make(types.JavaByteSlice, len(byteSlice))
	for ix, bb := range byteSlice {
		javaByteSlice[ix] = types.JavaByte(bb)
	}
	return javaByteSlice
}

// With the specified object and field, return a string representing the field value.
func ObjectFieldToString(obj *Object, fieldName string) string {
	// If null, return "null".
	if IsNull(obj) {
		return "null"
	}

	// If the field is missing, return "null".
	fld, ok := obj.FieldTable[fieldName]
	if !ok {
		return "null"
	}

	// If a static, remove the leading types.Static.
	if strings.HasPrefix(fld.Ftype, types.Static) {
		bytes := []byte(fld.Ftype)
		fld.Ftype = string(bytes[1:])
	}

	// What type is the field?
	switch fld.Ftype {
	case types.BigInteger:
		return fmt.Sprint(fld.Fvalue)
	case types.Bool:
		boolAsInt64 := fld.Fvalue.(int64)
		if boolAsInt64 > 0 {
			return "true"
		}
		return "false"
	case types.BoolArray:
		var str string
		for _, elem := range fld.Fvalue.([]int64) {
			if elem > 0 {
				str += "true"
			} else {
				str += "false"
			}
			str += " "
		}
		str = strings.TrimSuffix(str, " ")
		return str
	case types.Byte, types.Char, types.Int, types.Long, types.Rune, types.Short:
		return fmt.Sprintf("%d", fld.Fvalue.(int64))
	case types.ByteArray:
		switch fld.Fvalue.(type) {
		case []byte:
			return string(fld.Fvalue.([]byte))
		default:
			return JavaByteSliceToGoString(fld.Fvalue.(types.JavaByteSlice))
		}
	case types.CharArray, types.IntArray, types.LongArray, types.ShortArray:
		var str string
		for _, elem := range fld.Fvalue.([]int64) {
			str += fmt.Sprint(elem)
			str += " "
		}
		str = strings.TrimSuffix(str, " ")
		return str
	case types.Double, types.Float:
		return strconv.FormatFloat(fld.Fvalue.(float64), 'f', -1, 64)
	case types.DoubleArray, types.FloatArray:
		var str string
		for _, elem := range fld.Fvalue.([]float64) {
			str += strconv.FormatFloat(elem, 'f', -1, 64)
			str += " "
		}
		str = strings.TrimSuffix(str, " ")
		return str
	case types.FileHandle:
		return "FileHandle"
	case types.Ref, types.RefArray:
		return GoStringFromStringPoolIndex(obj.KlassName)
	}

	// None of the above.
	errMsg := fmt.Sprintf("ObjectFieldToString: field \"%s\" Ftype \"%s\" not yet supported. Returning the class name",
		fieldName, fld.Ftype)
	trace.Error(errMsg)
	return GoStringFromStringPoolIndex(obj.KlassName)
}

func GoStringFromJavaByteArray(jbarr []types.JavaByte) string {
	var sb strings.Builder
	for _, b := range jbarr {
		sb.WriteByte(byte(b))
	}
	return sb.String()
}

func JavaByteArrayFromGoString(str string) []types.JavaByte {
	jbarr := make([]types.JavaByte, len(str))
	for i, b := range str {
		jbarr[i] = types.JavaByte(b)
	}
	return jbarr
}

func JavaByteArrayFromGoByteArray(gbarr []byte) []types.JavaByte {
	jbarr := make([]types.JavaByte, len(gbarr))
	for i, b := range gbarr {
		jbarr[i] = types.JavaByte(b)
	}
	return jbarr
}
