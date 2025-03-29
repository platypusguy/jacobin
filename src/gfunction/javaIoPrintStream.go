/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/types"
	"os"
	"strconv"
)

/*
 Each object or library that has Go methods contains a reference to MethodSignatures,
 which contain data needed to insert the go method into the MTable of the currently
 executing JVM. MethodSignatures is a map whose key is the fully qualified name and
 type of the method (that is, the method's full signature) and a value consisting of
 a struct of an int (the number of slots to pop off the caller's operand stack when
 creating the new frame and a function. All methods have the same signature, regardless
 of the signature of their Java counterparts. That signature is that it accepts a slice
 of interface{} and returns an interface{}. The accepted slice can be empty and the
 return interface can be nil. This covers all Java functions. (Objects are returned
 as a 64-bit address in this scheme (as they are in the JVM).

 The slice contains one entry for every parameter passed to the method (which could
 mean an empty slice). There is no return value, because the method will place any
 return value on the operand stack of the calling function.
*/

func Load_Io_PrintStream() {
	MethodSignatures["java/io/PrintStream.println()V"] = // println void
		GMeth{
			ParamSlots: 0,
			GFunction:  PrintlnV,
		}
	MethodSignatures["java/io/PrintStream.println(Ljava/lang/String;)V"] = // println string
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnString,
		}
	MethodSignatures["java/io/PrintStream.println(B)V"] = // println byte
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnBIS,
		}
	MethodSignatures["java/io/PrintStream.println(C)V"] = // println char
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnChar,
		}
	MethodSignatures["java/io/PrintStream.println(I)V"] = // println int
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnBIS,
		}
	MethodSignatures["java/io/PrintStream.println(S)V"] = // println short
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnBIS,
		}
	MethodSignatures["java/io/PrintStream.println(Z)V"] = // println boolean
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnBoolean,
		}
	MethodSignatures["java/io/PrintStream.println(J)V"] = // println long
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnLong,
		}

	MethodSignatures["java/io/PrintStream.println(D)V"] = // println double
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnDouble,
		}

	MethodSignatures["java/io/PrintStream.println(F)V"] = // println float
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnFloat,
		}

	MethodSignatures["java/io/PrintStream.println(Ljava/lang/Object;)V"] = // println object
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnObject,
		}

	MethodSignatures["java/io/PrintStream.println([B)V"] = // println byte array
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnObject,
		}

	MethodSignatures["java/io/PrintStream.println([C)V"] = // println char array
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnObject,
		}

	MethodSignatures["java/io/PrintStream.println([D)V"] = // println double array
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnObject,
		}

	MethodSignatures["java/io/PrintStream.println([F)V"] = // println float array
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnObject,
		}

	MethodSignatures["java/io/PrintStream.println([I)V"] = // println int array
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnObject,
		}

	MethodSignatures["java/io/PrintStream.println([J)V"] = // println long array
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnObject,
		}

	MethodSignatures["java/io/PrintStream.println([S)V"] = // println int array
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnObject,
		}

	MethodSignatures["java/io/PrintStream.println([Z)V"] = // println boolean array
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnObject,
		}

	MethodSignatures["java/io/PrintStream.print(Ljava/lang/String;)V"] = // print string
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintString,
		}
	MethodSignatures["java/io/PrintStream.print(B)V"] = // print byte
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintBIS,
		}
	MethodSignatures["java/io/PrintStream.print(C)V"] = // print char
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintChar,
		}
	MethodSignatures["java/io/PrintStream.print(I)V"] = // print int
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintBIS,
		}
	MethodSignatures["java/io/PrintStream.print(S)V"] = // print short
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintBIS,
		}
	MethodSignatures["java/io/PrintStream.print(Z)V"] = // print boolean
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintBoolean,
		}
	MethodSignatures["java/io/PrintStream.print(J)V"] = // print long
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintLong,
		}

	MethodSignatures["java/io/PrintStream.print(D)V"] = // print double
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintDouble,
		}

	MethodSignatures["java/io/PrintStream.print(F)V"] = // print float
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintFloat,
		}

	MethodSignatures["java/io/PrintStream.print(Ljava/lang/Object;)V"] = // print object
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintObject,
		}

	MethodSignatures["java/io/PrintStream.print([B)V"] = // print byte array
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintObject,
		}

	MethodSignatures["java/io/PrintStream.print([C)V"] = // print char array
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintObject,
		}

	MethodSignatures["java/io/PrintStream.print([D)V"] = // print double array
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintObject,
		}

	MethodSignatures["java/io/PrintStream.print([F)V"] = // print float array
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintObject,
		}

	MethodSignatures["java/io/PrintStream.print([I)V"] = // print int array
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintObject,
		}

	MethodSignatures["java/io/PrintStream.print([J)V"] = // print long array
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintObject,
		}

	MethodSignatures["java/io/PrintStream.print([S)V"] = // print int array
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintObject,
		}

	MethodSignatures["java/io/PrintStream.print([Z)V"] = // print boolean array
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintObject,
		}

	MethodSignatures["java/io/PrintStream.printf(Ljava/lang/String;[Ljava/lang/Object;)Ljava/io/PrintStream;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  Printf,
		}

}

// PrintlnV = java/io/Prinstream.println() -- println() prints a newline (V = void)
// "java/io/PrintStream.println()V"
func PrintlnV(params []interface{}) interface{} {
	fmt.Fprintln(params[0].(*os.File), "")
	return nil
}

// "java/io/PrintStream.println(C)V"
func PrintlnChar(params []interface{}) interface{} {
	bb := byte(params[1].(int64))
	fmt.Fprintln(params[0].(*os.File), string(bb))
	return nil
}

// "java/io/PrintStream.println(B)V"
// "java/io/PrintStream.println(I)V"
// "java/io/PrintStream.println(S)V"
func PrintlnBIS(params []interface{}) interface{} {
	intToPrint := params[1].(int64) // contains an int
	fmt.Fprintln(params[0].(*os.File), intToPrint)
	return nil
}

// "java/io/PrintStream.println(Z)V"
func PrintlnBoolean(params []interface{}) interface{} {
	var boolToPrint bool
	boolAsInt64 := params[1].(int64) // contains an int64
	if boolAsInt64 > 0 {
		boolToPrint = true
	} else {
		boolToPrint = false
	}
	fmt.Fprintln(params[0].(*os.File), boolToPrint)
	return nil
}

// "java/io/PrintStream.println(J)V"
func PrintlnLong(params []interface{}) interface{} {
	longToPrint := params[1].(int64) // contains to an int64--the equivalent of a Java long
	fmt.Fprintln(params[0].(*os.File), longToPrint)
	return nil
}

// PrintlnDouble = java/io/Prinstream.print(double)
func PrintlnDouble(params []interface{}) interface{} {
	xx := params[1].(float64) // contains to a float64--the equivalent of a Java double
	fmt.Fprintln(params[0].(*os.File), strconv.FormatFloat(xx, 'g', -1, 64))
	return nil
}

// PrintlnFloat = java/io/Prinstream.print(float)
func PrintlnFloat(params []interface{}) interface{} {
	xx := params[1].(float64) // contains to a float64--the equivalent of a Java double
	fmt.Fprintln(params[0].(*os.File), strconv.FormatFloat(xx, 'g', -1, 32))
	return nil
}

// "java/io/PrintStream.print(C)V"
func PrintChar(params []interface{}) interface{} {
	bb := byte(params[1].(int64))
	fmt.Fprint(params[0].(*os.File), string(bb))
	return nil
}

// "java/io/PrintStream.print(B)V"
// "java/io/PrintStream.print(I)V"
// "java/io/PrintStream.print(S)V"
func PrintBIS(params []interface{}) interface{} {
	intToPrint := params[1].(int64) // contains an int
	fmt.Fprint(params[0].(*os.File), intToPrint)
	return nil
}

// PrintBoolean = java/io/Prinstream.print(boolean)
// "java/io/PrintStream.print(Z)V"
func PrintBoolean(params []interface{}) interface{} {
	var boolToPrint bool
	boolAsInt64 := params[1].(int64) // contains an int64
	if boolAsInt64 > 0 {
		boolToPrint = true
	} else {
		boolToPrint = false
	}
	fmt.Fprint(params[0].(*os.File), boolToPrint)
	return nil
}

// PrintLong = java/io/Prinstream.print(long)
// Long in Java are 64-bit ints, so we just duplicated the logic for println(int)
// "java/io/PrintStream.print(J)V"
func PrintLong(params []interface{}) interface{} {
	longToPrint := params[1].(int64) // contains to an int64--the equivalent of a Java long
	fmt.Fprint(params[0].(*os.File), longToPrint)
	return nil
}

// PrintDouble = java/io/Prinstream.print(double)
func PrintDouble(params []interface{}) interface{} {
	xx := params[1].(float64) // contains to a float64--the equivalent of a Java double
	fmt.Fprint(params[0].(*os.File), strconv.FormatFloat(xx, 'g', -1, 64))
	return nil
}

// PrintFloat = java/io/Prinstream.print(float)
func PrintFloat(params []interface{}) interface{} {
	xx := params[1].(float64) // contains to a float64--the equivalent of a Java double
	fmt.Fprint(params[0].(*os.File), strconv.FormatFloat(xx, 'g', -1, 32))
	return nil
}

// Printf -- handle the variable args and then call golang's own printf function
// "java/io/PrintStream.printf(Ljava/lang/String;[Ljava/lang/Object;)Ljava/io/PrintStream;"
func Printf(params []interface{}) interface{} {
	var intfSprintf = new([]interface{})
	*intfSprintf = append(*intfSprintf, params[1]) // The format string
	*intfSprintf = append(*intfSprintf, params[2]) // The object array
	retval := StringFormatter(*intfSprintf)
	switch retval.(type) {
	case *object.Object:
	default:
		return retval
	}
	objPtr := retval.(*object.Object)
	str := object.GoStringFromStringObject(objPtr)
	switch params[0].(type) {
	case *os.File:
		break
	default:
		errMsg := fmt.Sprintf("Printf: Expected parameter type *os.File, observed: %T", params[0])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	fmt.Fprint(params[0].(*os.File), str)
	return params[0] // Return the PrintStream object

}

// Called by PrintObject and PrintlnObject
func _printObject(params []interface{}, newLine bool) interface{} {
	var str string
	switch params[1].(type) {
	case *object.Object:
		inObj := params[1].(*object.Object)
		str = object.ObjectFieldToString(inObj, "FilePath")
		if str == types.NullString {
			str = object.ObjectFieldToString(inObj, "value")
			if str == types.NullString {
				className := object.GoStringFromStringPoolIndex(inObj.KlassName)
				if newLine {
					fmt.Fprintf(params[0].(*os.File), "class: %s, fields:\n", className)
				} else {
					fmt.Fprintf(params[0].(*os.File), "class: %s, fields: ", className)
				}
				str = ""
				for name, _ := range inObj.FieldTable {
					if newLine {
						str += fmt.Sprintf("%s=%s\n", name, object.ObjectFieldToString(inObj, name))
					} else {
						str += fmt.Sprintf("%s=%s, ", name, object.ObjectFieldToString(inObj, name))
					}
				}
				if newLine {
					fmt.Fprint(params[0].(*os.File), str)
					return nil
				} else {
					str = str[:len(str)-2]
					fmt.Fprint(params[0].(*os.File), str)
					return nil
				}
			}

		}
	case nil:
		str = types.NullString
	default:
		errMsg := fmt.Sprintf("_printObject: Unsupported parameter type: %T", params[1])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	if newLine {
		fmt.Fprintln(params[0].(*os.File), str)
	} else {
		fmt.Fprint(params[0].(*os.File), str)
	}

	return nil
}

// Print an Object's contents
// "java/io/PrintStream.print(Ljava/lang/Object;)V"
func PrintObject(params []interface{}) interface{} {
	return _printObject(params, false)
}

// Println an Object's contents
// "java/io/PrintStream.println(Ljava/lang/Object;)V"
func PrintlnObject(params []interface{}) interface{} {
	return _printObject(params, true)
}

// "java/io/PrintStream.println(Ljava/lang/String;)V"
func _printString(params []interface{}, newLine bool) interface{} {
	var str string
	param1, ok := params[1].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("_printString: Expected params[1] of type *object.Object but observed type %T\n", params[1])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Handle null strings as well as []byte.
	fld := param1.FieldTable["value"]
	if fld.Fvalue == nil {
		str = ""
	} else {
		switch fld.Fvalue.(type) {
		case []byte:
			str = string(fld.Fvalue.([]byte))
		case []types.JavaByte:
			str = object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
		default:
			errMsg := fmt.Sprintf("_printString: Expected value field to be type byte but observed type %T\n", fld.Fvalue)
			return getGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	}

	if newLine {
		fmt.Fprintln(params[0].(*os.File), str)
	} else {
		fmt.Fprint(params[0].(*os.File), str)
	}

	return nil
}

// Print string
// "java/io/PrintStream.print(Ljava/lang/String;)V"
func PrintString(params []interface{}) interface{} {
	return _printString(params, false)
}

// "java/io/PrintStream.println(Ljava/lang/String;)V"
func PrintlnString(params []interface{}) interface{} {
	return _printString(params, true)
}
