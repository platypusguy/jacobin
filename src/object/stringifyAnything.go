/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024-5 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package object

import (
	"container/list"
	"encoding/hex"
	"fmt"
	"jacobin/src/types"
	"reflect"
	"strconv"
	"strings"
)

// StringifyAnythingGo: Stringify anything and return a go string to caller.
// arg: either an *Object or a Field.
func StringifyAnythingGo(arg interface{}) string {

	// Watch out for nil arguments.
	if arg == nil {
		return types.NullString
	}

	// Begin outer switch on arg type.
	switch arg.(type) {
	case *Object:
		obj := arg.(*Object)
		if IsNull(obj) {
			return types.NullString
		}

		className := GoStringFromStringPoolIndex(obj.KlassName)
		if className == "[C" || className == types.CharArray {
			if fld, ok := obj.FieldTable["value"]; ok {
				if charArr, ok := fld.Fvalue.([]int64); ok {
					return GoStringFromJavaCharArray(charArr)
				}
			}
		}

		if strings.HasPrefix(className, "[L") || strings.HasPrefix(className, "[[") {
			if fld, ok := obj.FieldTable["value"]; ok {
				return StringifyAnythingGo(fld)
			}
		}

		classNameSuffix := GetClassNameSuffix(obj, true)
		switch classNameSuffix {
		case "String":
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "*** ERROR, StringifyAnythingGo: corrupted String object, missing \"value\" field"
			}
			switch val := fld.Fvalue.(type) {
			case []types.JavaByte:
				return GoStringFromJavaByteArray(val)
			case []byte:
				return string(val)
			case []int64:
				return GoStringFromJavaCharArray(val)
			default:
				return fmt.Sprintf("*** ERROR, StringifyAnythingGo: String/CharArray value field is %T, not byte or char array", fld.Fvalue)
			}
		case "Boolean":
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "*** ERROR, StringifyAnythingGo: corrupted Boolean object, missing \"value\" field"
			}
			if fld.Fvalue == types.JavaBoolTrue {
				return "true"
			}
			return "false"
		case "Byte":
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "*** ERROR, StringifyAnythingGo: corrupted Byte object, missing \"value\" field"
			}
			return fmt.Sprintf("0x%02x", fld.Fvalue.(int64)&0xff)
		case "Character":
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "*** ERROR, StringifyAnythingGo: corrupted Character object, missing \"value\" field"
			}
			return fmt.Sprintf("%c", fld.Fvalue.(int64))
		case "Double":
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "*** ERROR, StringifyAnythingGo: corrupted Double object, missing \"value\" field"
			}
			return strconv.FormatFloat(fld.Fvalue.(float64), 'g', -1, 64)
		case "Float":
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "*** ERROR, StringifyAnythingGo: corrupted Float object, missing \"value\" field"
			}
			return strconv.FormatFloat(fld.Fvalue.(float64), 'g', -1, 32)
		case "Integer", "Long", "Short":
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "*** ERROR, StringifyAnythingGo: corrupted Integer/Long/Short object, missing \"value\" field"
			}
			return strconv.FormatInt(fld.Fvalue.(int64), 10)
		case types.CharArray:
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "*** ERROR, StringifyAnythingGo: char array missing \"value\" field"
			}
			if charArr, ok := fld.Fvalue.([]int64); ok {
				return GoStringFromJavaCharArray(charArr)
			}
			return fmt.Sprintf("*** ERROR, StringifyAnythingGo: corrupted char array \"value\" field, type %T", fld.Fvalue)
		case types.ByteArray:
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "*** ERROR, StringifyAnythingGo: corrupted byte array, missing \"value\" field"
			}
			fvalue, ok := fld.Fvalue.([]byte)
			if !ok {
				jb, ok := fld.Fvalue.([]types.JavaByte)
				if !ok {
					return "*** ERROR, StringifyAnythingGo: corrupted byte array \"value\" field"
				}
				fvalue = GoByteArrayFromJavaByteArray(jb)
			}
			return "0x" + hex.EncodeToString(fvalue)
		case types.BoolArray:
			strBuffer := "["
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "*** ERROR, StringifyAnythingGo: boolean object missing \"value\" field"
			}
			boolArray, ok := fld.Fvalue.([]int64)
			if !ok {
				return "*** ERROR, StringifyAnythingGo: corrupted boolean array \"value\" field"
			}
			for _, elem := range boolArray {
				if elem > 0 {
					strBuffer += "true"
				} else {
					strBuffer += "false"
				}
				strBuffer += ", "
			}
			strBuffer = strBuffer[:len(strBuffer)-2] + "]"
			return strBuffer
		case types.DoubleArray:
			array, ok := obj.FieldTable["value"].Fvalue.([]float64)
			if !ok {
				return "*** ERROR, StringifyAnythingGo: double array missing \"value\" field or array value corrupted"
			}
			strBuffer := "["
			for ix := 0; ix < len(array); ix++ {
				strBuffer += strconv.FormatFloat(array[ix], 'g', -1, 64) + ", "
			}
			return strBuffer[:len(strBuffer)-2] + "]"
		case types.FloatArray:
			array, ok := obj.FieldTable["value"].Fvalue.([]float64)
			if !ok {
				return "*** ERROR, StringifyAnythingGo: float array missing \"value\" field or array value corrupted"
			}
			strBuffer := "["
			for ix := 0; ix < len(array); ix++ {
				strBuffer += strconv.FormatFloat(array[ix], 'g', -1, 32) + ", "
			}
			return strBuffer[:len(strBuffer)-2] + "]"
		case types.IntArray, types.LongArray, types.ShortArray:
			array, ok := obj.FieldTable["value"].Fvalue.([]int64)
			if !ok {
				return "*** ERROR, StringifyAnythingGo: int/long/short array missing \"value\" field or array value corrupted"
			}
			strBuffer := "["
			for ix := 0; ix < len(array); ix++ {
				strBuffer += strconv.FormatInt(array[ix], 10) + ", "
			}
			return strBuffer[:len(strBuffer)-2] + "]"
		case types.MultiArray, types.Array:
			strBuffer := "["
			fvalue := obj.FieldTable["value"].Fvalue
			if innerObj, ok := fvalue.(*Object); ok {
				if val, ok := innerObj.FieldTable["value"]; ok {
					fvalue = val.Fvalue
				}
			}

			anArray := reflect.ValueOf(fvalue)
			if anArray.Kind() != reflect.Slice && anArray.Kind() != reflect.Array {
				return fmt.Sprintf("*** ERROR, StringifyAnythingGo Field: expected array/slice but saw %T", fvalue)
			}

			for ix := 0; ix < anArray.Len(); ix++ {
				if ix > 0 {
					strBuffer += ", "
				}
				element := anArray.Index(ix).Interface()
				strBuffer += StringifyAnythingGo(element)
			}
			return strBuffer + "]"
		default:
			// Format a small report of the class name and the FieldTable.
			// Concoct a string buffer formatted as: class{name1=value1, name2=value2, ...}.
			strBuffer := classNameSuffix + "{"
			for name := range obj.FieldTable {
				strBuffer += name + "=" + ObjectFieldToString(obj, name) + ", "
			}
			return strBuffer[:len(strBuffer)-2] + "}"
		}
		/* end of case for *Object */
	case Field:
		fld := arg.(Field)
		switch fld.Fvalue.(type) {
		case *Object:
			return StringifyAnythingGo(fld.Fvalue.(*Object))
		}
		switch fld.Ftype {
		case types.StringClassRef:
			if IsNull(fld.Fvalue) {
				return types.NullString
			}
			switch fld.Fvalue.(type) {
			case []types.JavaByte:
				return GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
			case []byte:
				return string(fld.Fvalue.([]byte))
			case *Object:
				return StringifyAnythingGo(fld.Fvalue.(*Object))
			default:
				errMsg := fmt.Sprintf("*** ERROR, StringifyAnythingGo Field types.StringClassRef: expected Object but saw type %T",
					fld.Fvalue)
				return errMsg
			}
		case types.Byte:
			var ba1 []byte
			switch fld.Fvalue.(type) {
			case int64:
				ba1 = []byte{byte(fld.Fvalue.(int64))}
			case byte:
				ba1 = []byte{byte(fld.Fvalue.(byte))}
			default:
				return "*** ERROR, StringifyAnythingGo Field types.Byte: corrupted byte \"value\" field"
			}
			return "0x" + hex.EncodeToString(ba1)
		case types.ByteArray:
			var bytes []byte
			switch fld.Fvalue.(type) {
			case []types.JavaByte:
				bytes = GoByteArrayFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
			case []byte:
				bytes = fld.Fvalue.([]byte)
			default:
				return "*** ERROR, StringifyAnythingGo Field types.ByteArray: corrupted byte array \"value\" field"
			}
			return "0x" + hex.EncodeToString(bytes)
		case types.Bool:
			if fld.Fvalue == types.JavaBoolTrue {
				return "true"
			}
			return "false"
		case types.Char:
			return fmt.Sprintf("%c", fld.Fvalue.(int64))
		case types.Int, types.Short, types.Long:
			var strValue string
			switch fld.Fvalue.(type) {
			case int64:
				strValue = strconv.FormatInt(fld.Fvalue.(int64), 10)
			case uint64:
				strValue = strconv.FormatInt(int64(fld.Fvalue.(uint64)), 10)
			case int32:
				strValue = strconv.FormatInt(int64(fld.Fvalue.(int32)), 10)
			case uint32:
				strValue = strconv.FormatInt(int64(fld.Fvalue.(uint32)), 10)
			case int16:
				strValue = strconv.FormatInt(int64(fld.Fvalue.(int16)), 10)
			case uint16:
				strValue = strconv.FormatInt(int64(fld.Fvalue.(uint16)), 10)
			default:
				errMsg := fmt.Sprintf("*** ERROR, StringifyAnythingGo  Field types.Int: unrecognized field value type, value: %T, %v", arg, arg)
				return errMsg
			}
			return strValue
		case types.Double:
			return strconv.FormatFloat(fld.Fvalue.(float64), 'g', -1, 64)
		case types.Float:
			return strconv.FormatFloat(fld.Fvalue.(float64), 'g', -1, 32)
		case types.BigInteger:
			return fmt.Sprint(fld.Fvalue)
		case types.LinkedList: // LinkedList must contain objects, not primitives due to recursive call to this function
			strBuffer := "["
			llst := fld.Fvalue.(*list.List)
			if llst.Len() > 0 {
				element := llst.Front()
				for ix := 0; ix < llst.Len(); ix++ {
					strBuffer += StringifyAnythingGo(element.Value)
					strBuffer += ", "
					element = element.Next()
				}
				return strBuffer[:len(strBuffer)-2] + "]"
			} else {
				return "[]"
			}
		case types.MultiArray, types.Array:
			strBuffer := "["
			fvalue := fld.Fvalue
			if obj, ok := fvalue.(*Object); ok {
				if val, ok := obj.FieldTable["value"]; ok {
					fvalue = val.Fvalue
				}
			}

			anArray := reflect.ValueOf(fvalue)
			if anArray.Kind() != reflect.Slice && anArray.Kind() != reflect.Array {
				return fmt.Sprintf("*** ERROR, StringifyAnythingGo Field: expected array/slice but saw %T", fvalue)
			}

			for ix := 0; ix < anArray.Len(); ix++ {
				if ix > 0 {
					strBuffer += ", "
				}
				element := anArray.Index(ix).Interface()
				strBuffer += StringifyAnythingGo(element)
			}
			return strBuffer + "]"
		case types.IntArray, types.LongArray, types.ShortArray, types.BoolArray, types.DoubleArray, types.FloatArray:
			return StringifyAnythingGo(&Object{
				KlassName:  types.StringPoolStringIndex, // dummy
				FieldTable: map[string]Field{"value": fld},
			})
		case types.CharArray:
			if charArr, ok := fld.Fvalue.([]int64); ok {
				return GoStringFromJavaCharArray(charArr)
			}
			return fmt.Sprintf("%v", fld.Fvalue)
		case types.Ref:
			if obj, ok := fld.Fvalue.(*Object); ok {
				return StringifyAnythingGo(obj)
			}
			return fmt.Sprintf("%v", fld.Fvalue)
		default:
			if obj, ok := fld.Fvalue.(*Object); ok {
				return StringifyAnythingGo(obj)
			}
			// If it's not an object, it might be a raw slice or primitive.
			// Try the fallback logic.
			res := StringifyAnythingGo(fld.Fvalue)
			if !strings.HasPrefix(res, "*** ERROR, StringifyAnythingGo: neither *Object nor Field") {
				return res
			}
			return fmt.Sprintf("%v", fld.Fvalue)
		}
		/* end of case Field */
	}

	// If we got here, then the argument was neither an *Object nor a Field.
	// But it might be a primitive or something else we can handle with fmt.Sprintf.
	val := reflect.ValueOf(arg)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(val.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'g', -1, 64)
	case reflect.Bool:
		if val.Bool() {
			return "true"
		}
		return "false"
	case reflect.String:
		return val.String()
	case reflect.Slice, reflect.Array:
		strBuffer := "["
		for ix := 0; ix < val.Len(); ix++ {
			if ix > 0 {
				strBuffer += ", "
			}
			strBuffer += StringifyAnythingGo(val.Index(ix).Interface())
		}
		return strBuffer + "]"
	case reflect.Interface, reflect.Ptr:
		if val.IsNil() {
			return types.NullString
		}
		return StringifyAnythingGo(val.Elem().Interface())
	}

	errMsg := fmt.Sprintf("*** ERROR, StringifyAnythingGo: neither *Object nor Field, value: %T, %v", arg, arg)
	return errMsg
}

// StringifyAnythingJava: Stringify anything and return the Java String object for that string.
func StringifyAnythingJava(arg interface{}) *Object {
	return StringObjectFromGoString(StringifyAnythingGo(arg))
}
