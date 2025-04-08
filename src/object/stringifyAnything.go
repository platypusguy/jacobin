package object

import (
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
		/*
			Start of objects
		*/
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
		case "Double", "Float":
			fld, ok := obj.FieldTable["value"]
			if !ok {
				return "StringifyAnythingGo: corrupted Double/Float object, missing \"value\" field"
			}
			return strconv.FormatFloat(fld.Fvalue.(float64), 'g', -1, 64)
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
		case types.DoubleArray, types.FloatArray:
			array, ok := obj.FieldTable["value"].Fvalue.([]float64)
			if !ok {
				return "StringifyAnythingGo: double/float array missing \"value\" field or array value corrupted"
			}
			strBuffer := "["
			for ix := 0; ix < len(array); ix++ {
				strBuffer += strconv.FormatFloat(array[ix], 'g', -1, 64) + ", "
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
			for name, _ := range obj.FieldTable {
				strBuffer += name + "=" + ObjectFieldToString(obj, name) + ", "
			}
			return strBuffer[:len(strBuffer)-2] + "}"
		}
		/*
			End of objects
		*/
		// ------------------------------------------------------------------------------------------------------
	case Field:
		/*
			Start of Field
		*/
		fld := arg.(Field)
		switch fld.Ftype {
		case types.StringClassRef:
			return GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
		case types.Byte:
			return "0x" + hex.EncodeToString([]byte{fld.Fvalue.(byte)})
		case types.Bool:
			if fld.Fvalue == types.JavaBoolTrue {
				return "true"
			}
			return "false"
		case types.Int, types.Short, types.Long:
			return strconv.FormatInt(fld.Fvalue.(int64), 10)
		case types.Double, types.Float:
			return strconv.FormatFloat(fld.Fvalue.(float64), 'g', -1, 64)
		case types.BigInteger:
			return fmt.Sprint(fld.Fvalue)
		default:
			errMsg := fmt.Sprintf("StringifyAnythingGo: unrecognized argument type, value: %T, %v", arg, arg)
			return errMsg
		}
		/*
			End of Field
		*/
		// ------------------------------------------------------------------------------------------------------
	}
	/*
		End of Field
	*/

	errMsg := fmt.Sprintf("StringifyAnythingGo: unrecognized argument type, value: %T, %v", arg, arg)
	return errMsg

}

// StringifyAnythingJava: Stringify anything and return the Java String version of that string.
func StringifyAnythingJava(arg interface{}) *Object {
	return StringObjectFromGoString(StringifyAnythingGo(arg))
}
