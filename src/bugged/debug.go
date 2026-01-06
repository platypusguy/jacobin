package bugged

import (
	"fmt"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"strings"
)

type DebugObject struct {
	*object.Object
}

/*
This odd-looking code is there to force the Go compiler to include object.TVO() and object.STR() in the executable.
The compiler excludes any functions that have no references (calls).
If they are not included, then one cannot use them in the GoLand debugger.
*/
func DebugInit() {
	if os.Getenv("Revenge_is_best_served_cold!") == "Well, maybe?" {
		dummyObj := object.MakeEmptyObject()
		_ = TVO(dummyObj)
		_ = STR(dummyObj)
	}
}

// Tree View Object (TVO) provides a comprehensive debug view of an Object
// for the GoLand debugger.
//
// Call from "Evaluate Expression": TVO(object)
// Displays:
// - Mark.Misc value
// - Class name from string pool
// - String values for []int8 fields
// - Integer values for integer fields
func TVO(obj *object.Object) string {
	if object.IsNull(obj) {
		return "Object: nil"
	}

	var sb strings.Builder
	sb.WriteString("=== Tree View Object ===\n")

	// Display Mark.Misc value
	sb.WriteString(fmt.Sprintf("Mark.Misc: %d (0x%08X)\n", obj.Mark.Misc, obj.Mark.Misc))

	// Display the class name from the string pool.
	className := object.GoStringFromStringPoolIndex(obj.KlassName)
	sb.WriteString(fmt.Sprintf("Class: %s\n", className))

	// Display fields
	if len(obj.FieldTable) > 0 {

		// Create a slice of keys.
		keys := make([]string, 0, len(obj.FieldTable))
		for key := range obj.FieldTable {
			keys = append(keys, key)
		}

		// Sort the keys, case-insensitive.
		globals.SortCaseInsensitive(&keys)

		sb.WriteString("Fields:\n")

		// For each field .....
		for _, fieldName := range keys {

			field := obj.FieldTable[fieldName]
			value := field.Fvalue
			// Check for integer types
			switch value.(type) {

			case int64:
				if field.Ftype == types.Bool {
					var str string
					if field.Fvalue == types.JavaBoolTrue {
						str = "true"
					} else {
						str = "false"
					}
					sb.WriteString(fmt.Sprintf("  %s [%s]: %s\n", fieldName, field.Ftype, str))
				} else {
					sb.WriteString(fmt.Sprintf("  %s [%s]: %d\n", fieldName, field.Ftype, value))
				}

			case []types.JavaByte:
				if field.Ftype == types.ByteArray || field.Ftype == types.StringClassName || field.Ftype == types.StringClassRef {
					str := object.GoStringFromJavaByteArray(value.([]types.JavaByte))
					sb.WriteString(fmt.Sprintf("  %s [%s]: %s\n", fieldName, field.Ftype, str))
				} else {
					sb.WriteString(fmt.Sprintf("  %s [%s]: %v\n", fieldName, field.Ftype, field.Fvalue))
				}

			case *object.Object:
				clname := object.GoStringFromStringPoolIndex(value.(*object.Object).KlassName)
				if clname == types.StringClassName {
					str := object.GoStringFromStringObject(value.(*object.Object))
					sb.WriteString(fmt.Sprintf("  %s [object %s]: %s\n", fieldName, field.Ftype, str))
				} else {
					sb.WriteString(fmt.Sprintf("  %s [%s]: class %s\n", fieldName, field.Ftype, clname))
				}

			case int8, int16, int32, uint8, uint16, uint32, uint64:
				sb.WriteString(fmt.Sprintf("  %s [%s]: %d\n", fieldName, field.Ftype, value))

			default:
				// For other types, show type and value
				sb.WriteString(fmt.Sprintf("  %s [%s]: %v\n", fieldName, field.Ftype, field.Fvalue))
			}
		}
	} else {
		sb.WriteString("Fields: (none)\n")
	}

	return sb.String()
}

// STR provides a string view of a Java byte array ([]types.JavaByte)
// for the GoLand debugger.
//
// Call from "Evaluate Expression": object.STR(array)
func STR(obj *object.Object) string {

	if object.IsNull(obj) {
		return "Object: nil"
	}

	className := object.GoStringFromStringPoolIndex(obj.KlassName)
	if className != types.StringClassName {
		return "Not a Java String"
	}

	field, ok := obj.FieldTable["value"]
	if !ok {
		return "Missing the \"value\" field"
	}

	valueJBA, ok := field.Fvalue.([]types.JavaByte)
	if !ok {
		return "Value field is not a []types.JavaByte"
	}

	return object.GoStringFromJavaByteArray(valueJBA)
}
