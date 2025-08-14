package gfunction

import (
	"fmt"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/types"
	"strings"
)

// String formatting given a format string and a slice of arguments.
// Called by sprintf, javaIoConsole.go, and javaIoPrintStream.go.
func StringFormatter(params []interface{}) interface{} {
	// params[0]: format string
	// params[1]: argument slice (array of object pointers)

	// Check the parameter length. It should be 2.
	lenParams := len(params)
	if lenParams < 1 || lenParams > 2 {
		errMsg := fmt.Sprintf("StringFormatter: Invalid parameter count: %d", lenParams)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	if lenParams == 1 { // No parameters beyond the format string
		formatStringObj := params[0].(*object.Object)
		return formatStringObj
	}

	// Check the format string.
	var formatString string
	switch params[0].(type) {
	case *object.Object:
		formatStringObj := params[0].(*object.Object) // the format string is passed as a pointer to a string object
		switch formatStringObj.FieldTable["value"].Fvalue.(type) {
		case []byte:
			formatString = object.GoStringFromStringObject(formatStringObj)
		case []types.JavaByte:
			formatString =
				object.GoStringFromJavaByteArray(formatStringObj.FieldTable["value"].Fvalue.([]types.JavaByte))
		default:
			errMsg := fmt.Sprintf("StringFormatter: In the format string object, expected Ftype=%s but observed: %s",
				types.ByteArray, formatStringObj.FieldTable["value"].Ftype)
			return getGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	default:
		errMsg := fmt.Sprintf("StringFormatter: Expected a string object for the format string but observed: %T", params[0])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Make sure that the argument slice is a reference array.
	valuesOut := []any{}
	field := params[1].(*object.Object).FieldTable["value"]
	if !strings.HasPrefix(field.Ftype, types.RefArray) {
		errMsg := fmt.Sprintf("StringFormatter: Expected Ftype=%s for params[1]: fld.Ftype=%s, fld.Fvalue=%v",
			types.RefArray, field.Ftype, field.Fvalue)
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// valuesIn = the reference array
	valuesIn := field.Fvalue.([]*object.Object)

	// Main loop for reference array.
	for ii := 0; ii < len(valuesIn); ii++ {

		// Get the current object's value field.
		fld := valuesIn[ii].FieldTable["value"]

		// If type is string object, process it.
		if fld.Ftype == types.StringClassRef {
			var str string
			switch fld.Fvalue.(type) {
			case []byte:
				str = string(fld.Fvalue.([]byte))
			case []types.JavaByte:
				str = object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
			}
			valuesOut = append(valuesOut, str)
		} else {
			// Not a string object.
			switch fld.Ftype {
			case types.ByteArray:
				var str string
				switch fld.Fvalue.(type) {
				case []byte:
					str = string(fld.Fvalue.([]byte))
				case []types.JavaByte:
					str = object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
				}
				valuesOut = append(valuesOut, str)
			case types.Byte:
				valuesOut = append(valuesOut, uint8(fld.Fvalue.(int64)))
			case types.Bool:
				var zz bool
				if fld.Fvalue.(int64) == 0 {
					zz = false
				} else {
					zz = true
				}
				valuesOut = append(valuesOut, zz)
			case types.Char:
				rooney := rune(fld.Fvalue.(int64))
				valuesOut = append(valuesOut, rooney)
			case types.Double:
				valuesOut = append(valuesOut, fld.Fvalue.(float64))
			case types.Float:
				valuesOut = append(valuesOut, fld.Fvalue.(float64))
			case types.Int:
				valuesOut = append(valuesOut, fld.Fvalue.(int64))
			case types.Long:
				valuesOut = append(valuesOut, fld.Fvalue.(int64))
			case types.Short:
				valuesOut = append(valuesOut, fld.Fvalue.(int64))
			default:
				errMsg := fmt.Sprintf("StringFormatter: Invalid parameter %d is of type %s", ii+1, fld.Ftype)
				return getGErrBlk(excNames.IllegalArgumentException, errMsg)
			}
		}
	}

	// Use golang fmt.Sprintf to do the heavy lifting.
	str := fmt.Sprintf(formatString, valuesOut...)

	// Return a pointer to an object.Object that wraps the string byte array.
	return object.StringObjectFromGoString(str)
}
