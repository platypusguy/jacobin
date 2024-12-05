/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"fmt"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/statics"
	"jacobin/types"
	"strconv"
)

// jj (Jacobin JVM) functions are functions that can be inserted inside Java programs
// for diagnostic purposes. They simply return when run in the JDK, but do what they're
// supposed to do when run under Jacobin.
//
// Note this is a rough first design that will surely be refined. (JACOBIN-624)

func Load_jj() {

	MethodSignatures["jj._dumpStatics(Ljava/lang/String;ILjava/lang/String;)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  jjDumpStatics,
		}

	MethodSignatures["jj._dumpObject(Ljava/lang/Object;Ljava/lang/String;I)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  jjDumpObject,
		}

	MethodSignatures["jj._getStaticString(Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  jjGetStaticString,
		}

	MethodSignatures["jj._getFieldString(Ljava/lang/Object;Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  jjGetFieldString,
		}
}

func jjStringify(ftype string, fvalue any) *object.Object {
	var str string
	switch ftype {
	case types.Bool:
		if fvalue.(int64) == 1 {
			str = "true"
		} else {
			str = "false"
		}
	case types.Byte: // uint8
		str = fmt.Sprintf("%02x", fvalue.(uint8))
	case types.Char, types.Rune:
		str = fmt.Sprintf("%c", fvalue.(int64))
	case types.Double:
		str = strconv.FormatFloat(fvalue.(float64), 'g', -1, 64)
	case types.Float:
		str = strconv.FormatFloat(float64(fvalue.(float64)), 'g', -1, 64)
	case types.Int:
		str = fmt.Sprintf("%d", fvalue.(int64))
	case types.Long:
		str = fmt.Sprintf("%d", fvalue.(int64))
	case "Ljava/lang/String;":
		str = object.GoStringFromStringObject(fvalue.(*object.Object))
	case types.Ref, types.ByteArray:
		if object.IsNull(fvalue.(*object.Object)) {
			str = "null"
		} else {
			obj := fvalue.(*object.Object)
			if obj.KlassName == globals.StringIndexString {
				// It is a Java String object. Return it as-is.
				return obj
			}
			// Not a Java String object.
			str = fmt.Sprintf("%v", fvalue)
		}
	case types.Short:
		str = fmt.Sprintf("%d", fvalue.(int64))
	default:
		str = fmt.Sprintf("%v", fvalue)
	}
	return object.StringObjectFromGoString(str)
}

func jjDumpStatics(params []interface{}) interface{} {
	fromObj := params[0].(*object.Object)
	if fromObj == nil || fromObj.KlassName == types.InvalidStringIndex {
		errMsg := fmt.Sprintf("jjDumpStatics: Invalid from object: %T", params[0])
		return object.StringObjectFromGoString(errMsg)
	}
	from := object.ObjectFieldToString(fromObj, "value")
	selection := params[1].(int64)
	classNameObj := params[2].(*object.Object)
	className := object.ObjectFieldToString(classNameObj, "value")

	statics.DumpStatics(from, selection, className)
	return nil
}

func jjDumpObject(params []interface{}) interface{} {
	this := params[0].(*object.Object)
	objTitle := params[1].(*object.Object)
	title := object.ObjectFieldToString(objTitle, "value")
	indent := params[2].(int64)
	this.DumpObject(title, int(indent))
	return nil
}

func jjGetStaticString(params []interface{}) interface{} {

	// Get class name.
	classObj := params[0].(*object.Object)
	if classObj == nil || classObj.KlassName == types.InvalidStringIndex {
		errMsg := fmt.Sprintf("jjGetStaticString: Invalid class object: %T", params[0])
		return object.StringObjectFromGoString(errMsg)
	}
	className := object.ObjectFieldToString(classObj, "value")

	// Get field name.
	fieldObj := params[1].(*object.Object)
	if fieldObj == nil || fieldObj.KlassName == types.InvalidStringIndex {
		errMsg := fmt.Sprintf("jjGetStaticString: Invalid field object: %T", params[1])
		return object.StringObjectFromGoString(errMsg)
	}
	fieldName := object.ObjectFieldToString(fieldObj, "value")

	// Convert statics entry to a string object.
	static := statics.Statics[className+"."+fieldName]
	return jjStringify(static.Type, static.Value)
}

func jjGetFieldString(params []interface{}) interface{} {

	// Get this object.
	thisObj := params[0].(*object.Object)

	// Get field name.
	fieldObj := params[1].(*object.Object)
	if fieldObj == nil || fieldObj.KlassName == types.InvalidStringIndex {
		errMsg := fmt.Sprintf("jjGetFieldString: Invalid field object: %T", params[1])
		return object.StringObjectFromGoString(errMsg)
	}
	fieldName := object.ObjectFieldToString(fieldObj, "value")

	// Convert field entry to a string object.
	fld, ok := thisObj.FieldTable[fieldName]
	if !ok {
		errMsg := fmt.Sprintf("jjGetFieldString: No such field name: %s", fieldName)
		return object.StringObjectFromGoString(errMsg)
	}
	if fld.Ftype == "Ljava/lang/String;" {
		return object.StringObjectFromByteArray(fld.Fvalue.([]byte))
	}
	return jjStringify(fld.Ftype, fld.Fvalue)
}
