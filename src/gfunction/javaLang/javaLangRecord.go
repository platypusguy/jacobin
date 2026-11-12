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
	"math"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

func Load_Lang_Record() {

	ghelpers.MethodSignatures["java/lang/Record.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/lang/Record.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  recordInit,
		}

	ghelpers.MethodSignatures["java/lang/Record.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  recordEquals,
		}

	ghelpers.MethodSignatures["java/lang/Record.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  recordHashCode,
		}

	ghelpers.MethodSignatures["java/lang/Record.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  recordToString,
		}

}

// recordInit implements "java/lang/Record.<init>()V"
// Protected constructor for java.lang.Record.
func recordInit(params []any) any {
	if len(params) == 0 {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "recordInit: missing receiver")
	}
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		errMsg := fmt.Sprintf("recordInit: Invalid or null receiver: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}
	return nil
}

// recordEquals implements "java/lang/Record.equals(Ljava/lang/Object;)Z"
// Indicates whether some other object is "equal to" this one.
func recordEquals(params []any) any {
	if len(params) == 0 {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "recordEquals: missing receiver")
	}
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		errMsg := fmt.Sprintf("recordEquals: Invalid or null receiver: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	if len(params) < 2 {
		return types.JavaBoolFalse
	}

	other, ok := params[1].(*object.Object)
	if !ok || object.IsNull(other) {
		return types.JavaBoolFalse
	}

	if self == other {
		return types.JavaBoolTrue
	}

	if self.KlassName != other.KlassName {
		return types.JavaBoolFalse
	}

	if len(self.FieldTable) != len(other.FieldTable) {
		return types.JavaBoolFalse
	}

	for k, selfFld := range self.FieldTable {
		otherFld, exists := other.FieldTable[k]
		if !exists {
			return types.JavaBoolFalse
		}
		if !fieldValuesEqual(selfFld.Fvalue, otherFld.Fvalue) {
			return types.JavaBoolFalse
		}
	}

	return types.JavaBoolTrue
}

func fieldValuesEqual(v1, v2 any) bool {
	if v1 == nil && v2 == nil {
		return true
	}
	if v1 == nil || v2 == nil {
		return false
	}

	switch val1 := v1.(type) {
	case *object.Object:
		val2, ok := v2.(*object.Object)
		if !ok {
			return false
		}
		if object.IsNull(val1) && object.IsNull(val2) {
			return true
		}
		if object.IsNull(val1) || object.IsNull(val2) {
			return false
		}
		if val1 == val2 {
			return true
		}
		if object.IsStringObject(val1) && object.IsStringObject(val2) {
			return object.GoStringFromStringObject(val1) == object.GoStringFromStringObject(val2)
		}
		return val1 == val2
	case int64:
		val2, ok := v2.(int64)
		return ok && val1 == val2
	case float64:
		val2, ok := v2.(float64)
		if !ok {
			return false
		}
		return math.Float64bits(val1) == math.Float64bits(val2)
	case float32:
		val2, ok := v2.(float32)
		if !ok {
			return false
		}
		return math.Float32bits(val1) == math.Float32bits(val2)
	case types.JavaByte:
		val2, ok := v2.(types.JavaByte)
		return ok && val1 == val2
	case string:
		val2, ok := v2.(string)
		return ok && val1 == val2
	case bool:
		val2, ok := v2.(bool)
		return ok && val1 == val2
	case []int64:
		val2, ok := v2.([]int64)
		if !ok || len(val1) != len(val2) {
			return false
		}
		for i := range val1 {
			if val1[i] != val2[i] {
				return false
			}
		}
		return true
	case []types.JavaByte:
		val2, ok := v2.([]types.JavaByte)
		if !ok || len(val1) != len(val2) {
			return false
		}
		for i := range val1 {
			if val1[i] != val2[i] {
				return false
			}
		}
		return true
	case []byte:
		val2, ok := v2.([]byte)
		if !ok || len(val1) != len(val2) {
			return false
		}
		for i := range val1 {
			if val1[i] != val2[i] {
				return false
			}
		}
		return true
	case []*object.Object:
		val2, ok := v2.([]*object.Object)
		if !ok || len(val1) != len(val2) {
			return false
		}
		for i := range val1 {
			if !fieldValuesEqual(val1[i], val2[i]) {
				return false
			}
		}
		return true
	default:
		return v1 == v2
	}
}

// recordHashCode implements "java/lang/Record.hashCode()I"
// Returns a hash code value for the record based on its components.
func recordHashCode(params []any) any {
	if len(params) == 0 {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "recordHashCode: missing receiver")
	}
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		errMsg := fmt.Sprintf("recordHashCode: Invalid or null receiver: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	keys := make([]string, 0, len(self.FieldTable))
	for k := range self.FieldTable {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var result int32 = 1
	for _, k := range keys {
		fld := self.FieldTable[k]
		fHash := computeFieldHashCode(fld)
		result = 31*result + fHash
	}

	return int64(result)
}

func computeFieldHashCode(fld object.Field) int32 {
	if fld.Fvalue == nil {
		return 0
	}
	if fld.Ftype == types.Bool || fld.Ftype == "Z" {
		if v, ok := fld.Fvalue.(int64); ok {
			if v == types.JavaBoolTrue {
				return 1231
			}
			return 1237
		}
	}
	switch v := fld.Fvalue.(type) {
	case *object.Object:
		if object.IsNull(v) {
			return 0
		}
		if object.IsStringObject(v) {
			return computeStringHashCode(object.GoStringFromStringObject(v))
		}
		ptr := uintptr(unsafe.Pointer(v))
		return int32(ptr ^ (ptr >> 32))
	case int64:
		return int32(v ^ (v >> 32))
	case float64:
		bits := math.Float64bits(v)
		return int32(bits ^ (bits >> 32))
	case float32:
		return int32(math.Float32bits(v))
	case types.JavaByte:
		return int32(v)
	case string:
		return computeStringHashCode(v)
	case bool:
		if v {
			return 1231
		}
		return 1237
	case []int64:
		var h int32 = 1
		for _, elem := range v {
			h = 31*h + int32(elem^(elem>>32))
		}
		return h
	case []types.JavaByte:
		var h int32 = 1
		for _, elem := range v {
			h = 31*h + int32(elem)
		}
		return h
	case []byte:
		var h int32 = 1
		for _, elem := range v {
			h = 31*h + int32(elem)
		}
		return h
	default:
		return 0
	}
}

func computeStringHashCode(s string) int32 {
	var h int32
	for i := 0; i < len(s); i++ {
		h = 31*h + int32(s[i])
	}
	return h
}

// recordToString implements "java/lang/Record.toString()Ljava/lang/String;"
// Returns a string representation of the record and its components.
func recordToString(params []any) any {
	if len(params) == 0 {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "recordToString: missing receiver")
	}
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		errMsg := fmt.Sprintf("recordToString: Invalid or null receiver: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.NullPointerException, errMsg)
	}

	className := object.GoStringFromStringPoolIndex(self.KlassName)
	simpleName := getSimpleClassName(className)

	keys := make([]string, 0, len(self.FieldTable))
	for k := range self.FieldTable {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		fld := self.FieldTable[k]
		parts = append(parts, fmt.Sprintf("%s=%s", k, formatFieldValue(fld)))
	}

	str := fmt.Sprintf("%s[%s]", simpleName, strings.Join(parts, ", "))
	return object.StringObjectFromGoString(str)
}

func getSimpleClassName(className string) string {
	if idx := strings.LastIndex(className, "/"); idx >= 0 {
		className = className[idx+1:]
	}
	if idx := strings.LastIndex(className, "."); idx >= 0 {
		className = className[idx+1:]
	}
	return className
}

func formatFieldValue(fld object.Field) string {
	if fld.Fvalue == nil {
		return "null"
	}
	if fld.Ftype == types.Bool || fld.Ftype == "Z" {
		if v, ok := fld.Fvalue.(int64); ok {
			if v == types.JavaBoolTrue {
				return "true"
			}
			return "false"
		}
	}
	if fld.Ftype == types.Char || fld.Ftype == "C" {
		if v, ok := fld.Fvalue.(int64); ok {
			return string(rune(v))
		}
	}
	switch v := fld.Fvalue.(type) {
	case *object.Object:
		if object.IsNull(v) {
			return "null"
		}
		if object.IsStringObject(v) {
			return object.GoStringFromStringObject(v)
		}
		return object.StringifyAnythingGo(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return ghelpers.JavaFloatingPointToString(v, 64)
	case float32:
		return ghelpers.JavaFloatingPointToString(float64(v), 32)
	case types.JavaByte:
		return strconv.Itoa(int(v))
	case string:
		return v
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}
