/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package object

import (
	"container/list"
	"encoding/hex"
	"fmt"
	"jacobin/types"
	"strconv"
)

// StringifyAnythingGo: Stringify anything and return string to caller.
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
		classNameSuffix := GetClassNameSuffix(obj, true)
		switch classNameSuffix {
		case "String":
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "StringifyAnythingGo: corrupted String object, missing \"value\" field"
			}
			return GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
		case "Boolean":
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "StringifyAnythingGo: corrupted Boolean object, missing \"value\" field"
			}
			if fld.Fvalue == types.JavaBoolTrue {
				return "true"
			}
			return "false"
		case "Byte":
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "StringifyAnythingGo: corrupted Byte object, missing \"value\" field"
			}
			return fmt.Sprintf("0x%02x", fld.Fvalue.(int64)&0xff)
		case "Character":
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "StringifyAnythingGo: corrupted Character object, missing \"value\" field"
			}
			return fmt.Sprintf("%d", fld.Fvalue.(int64)&0xff)
		case "Double":
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "StringifyAnythingGo: corrupted Double object, missing \"value\" field"
			}
			return strconv.FormatFloat(fld.Fvalue.(float64), 'g', -1, 64)
		case "Float":
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "StringifyAnythingGo: corrupted Float object, missing \"value\" field"
			}
			return strconv.FormatFloat(fld.Fvalue.(float64), 'g', -1, 32)
		case "Integer", "Long", "Short":
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "StringifyAnythingGo: corrupted Integer/Long/Short object, missing \"value\" field"
			}
			return strconv.FormatInt(fld.Fvalue.(int64), 10)
		case types.ByteArray:
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "StringifyAnythingGo: corrupted byte array, missing \"value\" field"
			}
			fvalue, ok := fld.Fvalue.([]byte)
			if !ok {
				jb, ok := fld.Fvalue.([]types.JavaByte)
				if !ok {
					return "StringifyAnythingGo: corrupted byte array \"value\" field"
				}
				fvalue = GoByteArrayFromJavaByteArray(jb)
			}
			return "0x" + hex.EncodeToString(fvalue)
		case types.BoolArray:
			strBuffer := "["
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "StringifyAnythingGo: boolean object missing \"value\" field"
			}
			boolArray, ok := fld.Fvalue.([]int64)
			if !ok {
				return "StringifyAnythingGo: corrupted boolean array \"value\" field"
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
				return "StringifyAnythingGo: double array missing \"value\" field or array value corrupted"
			}
			strBuffer := "["
			for ix := 0; ix < len(array); ix++ {
				strBuffer += strconv.FormatFloat(array[ix], 'g', -1, 64) + ", "
			}
			return strBuffer[:len(strBuffer)-2] + "]"
		case types.FloatArray:
			array, ok := obj.FieldTable["value"].Fvalue.([]float64)
			if !ok {
				return "StringifyAnythingGo: float array missing \"value\" field or array value corrupted"
			}
			strBuffer := "["
			for ix := 0; ix < len(array); ix++ {
				strBuffer += strconv.FormatFloat(array[ix], 'g', -1, 32) + ", "
			}
			return strBuffer[:len(strBuffer)-2] + "]"
		case types.IntArray, types.LongArray, types.ShortArray:
			array, ok := obj.FieldTable["value"].Fvalue.([]int64)
			if !ok {
				return "StringifyAnythingGo: int/long/short array missing \"value\" field or array value corrupted"
			}
			strBuffer := "["
			for ix := 0; ix < len(array); ix++ {
				strBuffer += strconv.FormatInt(array[ix], 10) + ", "
			}
			return strBuffer[:len(strBuffer)-2] + "]"
		case types.MultiArray:
			GoStringFromStringPoolIndex(arg.(*Object).KlassName)
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
		switch fld.Ftype {
		case types.StringClassRef:
			return GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
		case types.Byte:
			var ba1 []byte
			switch fld.Fvalue.(type) {
			case int64:
				ba1 = []byte{byte(fld.Fvalue.(int64))}
			case byte:
				ba1 = []byte{byte(fld.Fvalue.(byte))}
			default:
				return "StringifyAnythingGo Field types.Byte: corrupted byte \"value\" field"
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
				return "StringifyAnythingGo Field types.ByteArray: corrupted byte array \"value\" field"
			}
			return "0x" + hex.EncodeToString(bytes)
		case types.Bool:
			if fld.Fvalue == types.JavaBoolTrue {
				return "true"
			}
			return "false"
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
				errMsg := fmt.Sprintf("StringifyAnythingGo  Field types.Int: unrecognized field value type, value: %T, %v", arg, arg)
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
		default:
			errMsg := fmt.Sprintf("StringifyAnythingGo Field default: unrecognized argument type, value: %T, %v", arg, arg)
			return errMsg
		}
		/* end of case Field */
	}

	// If we got here, then the argument was neither an *Object nor a Field.
	errMsg := fmt.Sprintf("StringifyAnythingGo: neither *Object nor Field, value: %T, %v", arg, arg)
	return errMsg
}

// StringifyAnythingJava: Stringify anything and return the Java String version of that string.
func StringifyAnythingJava(arg interface{}) *Object {
	return StringObjectFromGoString(StringifyAnythingGo(arg))
}
