/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaIo

import (
	"fmt"
	"io"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaLang"
	"jacobin/src/gfunction/javaMath"
	"jacobin/src/gfunction/misc"
	"jacobin/src/object"
	"jacobin/src/types"
	"strconv"
)

/*
 Each object or library that has Go methods contains a reference to ghelpers.MethodSignatures,
 which contain data needed to insert the go method into the MTable of the currently
 executing JVM. ghelpers.MethodSignatures is a map whose key is the fully qualified name and
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
	ghelpers.MethodSignatures["java/io/PrintStream.println()V"] = // println void
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  PrintlnV,
		}
	ghelpers.MethodSignatures["java/io/PrintStream.println(Ljava/lang/String;)V"] = // println string
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnString,
		}
	ghelpers.MethodSignatures["java/io/PrintStream.println(B)V"] = // println byte
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnBIS,
		}
	ghelpers.MethodSignatures["java/io/PrintStream.println(C)V"] = // println char
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnChar,
		}
	ghelpers.MethodSignatures["java/io/PrintStream.println(I)V"] = // println int
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnBIS,
		}
	ghelpers.MethodSignatures["java/io/PrintStream.println(S)V"] = // println short
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnBIS,
		}
	ghelpers.MethodSignatures["java/io/PrintStream.println(Z)V"] = // println boolean
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnBoolean,
		}
	ghelpers.MethodSignatures["java/io/PrintStream.println(J)V"] = // println long
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnLong,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.println(D)V"] = // println double
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnDouble,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.println(F)V"] = // println float
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnFloat,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.println(Ljava/lang/Object;)V"] = // println object
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnObject,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.println([B)V"] = // println byte array
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnObject,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.println([C)V"] = // println char array
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnObject,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.println([D)V"] = // println double array
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnObject,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.println([F)V"] = // println float array
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnObject,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.println([I)V"] = // println int array
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnObject,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.println([J)V"] = // println long array
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnObject,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.println([S)V"] = // println int array
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnObject,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.println([Z)V"] = // println boolean array
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnObject,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.print(Ljava/lang/String;)V"] = // print string
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintString,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.print(B)V"] = // print byte
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintBIS,
		}
	ghelpers.MethodSignatures["java/io/PrintStream.print(C)V"] = // print char
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintChar,
		}
	ghelpers.MethodSignatures["java/io/PrintStream.print(I)V"] = // print int
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintBIS,
		}
	ghelpers.MethodSignatures["java/io/PrintStream.print(S)V"] = // print short
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintBIS,
		}
	ghelpers.MethodSignatures["java/io/PrintStream.print(Z)V"] = // print boolean
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintBoolean,
		}
	ghelpers.MethodSignatures["java/io/PrintStream.print(J)V"] = // print long
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintLong,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.print(D)V"] = // print double
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintDouble,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.print(F)V"] = // print float
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintFloat,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.print(Ljava/lang/Object;)V"] = // print object
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintObject,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.print([B)V"] = // print byte array
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintObject,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.print([C)V"] = // print char array
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintObject,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.print([D)V"] = // print double array
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintObject,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.print([F)V"] = // print float array
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintObject,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.print([I)V"] = // print int array
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintObject,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.print([J)V"] = // print long array
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintObject,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.print([S)V"] = // print int array
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintObject,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.print([Z)V"] = // print boolean array
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  PrintObject,
		}

	ghelpers.MethodSignatures["java/io/PrintStream.printf(Ljava/lang/String;[Ljava/lang/Object;)Ljava/io/PrintStream;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  Printf,
		}

}

// PrintlnV = java/io/Prinstream.println() -- println() prints a newline (V = void)
// "java/io/PrintStream.println()V"
func PrintlnV(params []interface{}) interface{} {
	writer, ok := params[0].(io.Writer)
	if !ok {
		errMsg := fmt.Sprintf("PrintlnV: Expected io.Writer, observed %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	fmt.Fprintln(writer, "")
	return nil
}

// "java/io/PrintStream.println(C)V"
func PrintlnChar(params []interface{}) interface{} {
	writer, ok := params[0].(io.Writer)
	if !ok {
		errMsg := fmt.Sprintf("PrintlnChar: Expected io.Writer, observed %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	bb := byte(params[1].(int64))
	fmt.Fprintln(writer, string(bb))
	return nil
}

// "java/io/PrintStream.println(B)V"
// "java/io/PrintStream.println(I)V"
// "java/io/PrintStream.println(S)V"
func PrintlnBIS(params []interface{}) interface{} {
	writer, ok := params[0].(io.Writer)
	if !ok {
		errMsg := fmt.Sprintf("PrintlnBIS: Expected io.Writer, observed %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	intToPrint, ok := params[1].(int64) // contains an int
	if !ok {
		intToPrint = int64(params[1].(int8))
	}
	fmt.Fprintln(writer, intToPrint)
	return nil
}

// "java/io/PrintStream.println(Z)V"
func PrintlnBoolean(params []interface{}) interface{} {
	writer, ok := params[0].(io.Writer)
	if !ok {
		errMsg := fmt.Sprintf("PrintlnBoolean: Expected io.Writer, observed %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	var boolToPrint bool
	boolAsInt64 := params[1].(int64) // contains an int64
	if boolAsInt64 > 0 {
		boolToPrint = true
	} else {
		boolToPrint = false
	}
	fmt.Fprintln(writer, boolToPrint)
	return nil
}

// "java/io/PrintStream.println(J)V"
func PrintlnLong(params []interface{}) interface{} {
	writer, ok := params[0].(io.Writer)
	if !ok {
		errMsg := fmt.Sprintf("PrintlnLong: Expected io.Writer, observed %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	longToPrint := params[1].(int64) // contains to an int64--the equivalent of a Java long
	fmt.Fprintln(writer, longToPrint)
	return nil
}

// PrintlnDouble = java/io/Prinstream.print(double)
func PrintlnDouble(params []interface{}) interface{} {
	writer, ok := params[0].(io.Writer)
	if !ok {
		errMsg := fmt.Sprintf("PrintlnDouble: Expected io.Writer, observed %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	xx := params[1].(float64) // contains to a float64--the equivalent of a Java double
	fmt.Fprintln(writer, strconv.FormatFloat(xx, 'g', -1, 64))
	return nil
}

// PrintlnFloat = java/io/Prinstream.print(float)
func PrintlnFloat(params []interface{}) interface{} {
	writer, ok := params[0].(io.Writer)
	if !ok {
		errMsg := fmt.Sprintf("PrintlnFloat: Expected io.Writer, observed %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	xx := params[1].(float64) // contains to a float64--the equivalent of a Java double
	fmt.Fprintln(writer, strconv.FormatFloat(xx, 'g', -1, 32))
	return nil
}

// "java/io/PrintStream.print(C)V"
func PrintChar(params []interface{}) interface{} {
	writer, ok := params[0].(io.Writer)
	if !ok {
		errMsg := fmt.Sprintf("PrintChar: Expected io.Writer, observed %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	bb := byte(params[1].(int64))
	fmt.Fprint(writer, string(bb))
	return nil
}

// "java/io/PrintStream.print(B)V"
// "java/io/PrintStream.print(I)V"
// "java/io/PrintStream.print(S)V"
func PrintBIS(params []interface{}) interface{} {
	writer, ok := params[0].(io.Writer)
	if !ok {
		errMsg := fmt.Sprintf("PrintBIS: Expected io.Writer, observed %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	intToPrint, ok := params[1].(int64) // contains an int
	if !ok {
		intToPrint = int64(params[1].(int8))
	}
	fmt.Fprint(writer, intToPrint)
	return nil
}

// PrintBoolean = java/io/Prinstream.print(boolean)
// "java/io/PrintStream.print(Z)V"
func PrintBoolean(params []interface{}) interface{} {
	writer, ok := params[0].(io.Writer)
	if !ok {
		errMsg := fmt.Sprintf("PrintBoolean: Expected io.Writer, observed %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	var boolToPrint bool
	boolAsInt64 := params[1].(int64) // contains an int64
	if boolAsInt64 > 0 {
		boolToPrint = true
	} else {
		boolToPrint = false
	}
	fmt.Fprint(writer, boolToPrint)
	return nil
}

// PrintLong = java/io/Prinstream.print(long)
// Long in Java are 64-bit ints, so we just duplicated the logic for println(int)
// "java/io/PrintStream.print(J)V"
func PrintLong(params []interface{}) interface{} {
	writer, ok := params[0].(io.Writer)
	if !ok {
		errMsg := fmt.Sprintf("PrintLong: Expected io.Writer, observed %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	longToPrint := params[1].(int64) // contains to an int64--the equivalent of a Java long
	fmt.Fprint(writer, longToPrint)
	return nil
}

// PrintDouble = java/io/Prinstream.print(double)
func PrintDouble(params []interface{}) interface{} {
	writer, ok := params[0].(io.Writer)
	if !ok {
		errMsg := fmt.Sprintf("PrintDouble: Expected io.Writer, observed %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	xx := params[1].(float64) // contains to a float64--the equivalent of a Java double
	fmt.Fprint(writer, strconv.FormatFloat(xx, 'g', -1, 64))
	return nil
}

// PrintFloat = java/io/Prinstream.print(float)
func PrintFloat(params []interface{}) interface{} {
	writer, ok := params[0].(io.Writer)
	if !ok {
		errMsg := fmt.Sprintf("PrintFloat: Expected io.Writer, observed %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	xx := params[1].(float64) // contains to a float64--the equivalent of a Java double
	fmt.Fprint(writer, strconv.FormatFloat(xx, 'g', -1, 32))
	return nil
}

// Printf -- handle the variable args and then call golang's own printf function
// "java/io/PrintStream.printf(Ljava/lang/String;[Ljava/lang/Object;)Ljava/io/PrintStream;"
func Printf(params []interface{}) interface{} {
	writer, ok := params[0].(io.Writer)
	if !ok {
		errMsg := fmt.Sprintf("Printf: Expected io.Writer, observed %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	var intfSprintf = new([]interface{})
	*intfSprintf = append(*intfSprintf, params[1]) // The format string
	*intfSprintf = append(*intfSprintf, params[2]) // The object array
	retval := misc.StringFormatter(*intfSprintf)
	switch retval.(type) {
	case *object.Object:
	default:
		return retval
	}
	objPtr := retval.(*object.Object)
	str := object.GoStringFromStringObject(objPtr)

	fmt.Fprint(writer, str)

	return params[0] // Return the PrintStream object
}

// "java/io/PrintStream.println(Ljava/lang/String;)V"
func _printString(params []interface{}, newLine bool) interface{} {
	writer, ok := params[0].(io.Writer)
	if !ok {
		errMsg := fmt.Sprintf("_printString: Expected io.Writer, observed %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	var str string
	param1, ok := params[1].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("_printString: Expected params[1] of type *object.Object but observed type %T\n", params[1])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Handle null strings.
	if param1 == nil || object.IsNull(param1) {
		str = types.NullString
	} else {
		fld, ok := param1.FieldTable["value"]
		if !ok {
			className := object.GoStringFromStringPoolIndex(param1.KlassName)
			errMsg := fmt.Sprintf("_printString: Class %s (String?), \"value\" field is missing", className)
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
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
				return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
			}
		}
	}

	if newLine {
		fmt.Fprintln(writer, str)
	} else {
		fmt.Fprint(writer, str)
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

// Called by PrintObject and PrintlnObject
func _printObject(params []interface{}, newLine bool) interface{} {
	writer, ok := params[0].(io.Writer)
	if !ok {
		errMsg := fmt.Sprintf("_printObject: Expected io.Writer, observed %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	var strBuffer string

	// Watch out for a null object.
	if params[1] == nil || object.IsNull(params[1]) {
		strBuffer = types.NullString
	} else {
		switch params[1].(type) {
		case *object.Object:
			inObj := params[1].(*object.Object)
			if object.IsStringObject(inObj) {
				strBuffer = object.GoStringFromStringObject(inObj)
				break
			}
			if object.IsObjectClass(inObj, types.ClassNameBigDecimal) {
				strObj := javaMath.BigdecimalToString([]interface{}{inObj}).(*object.Object)
				strBuffer = object.GoStringFromStringObject(strObj)
				break
			}
			if object.IsObjectClass(inObj, types.ClassNameBigInteger) {
				strObj := javaMath.BigIntegerToString([]interface{}{inObj}).(*object.Object)
				strBuffer = object.GoStringFromStringObject(strObj)
				break
			}
			if object.IsObjectClass(inObj, types.ClassNameThread) {
				strObj := javaLang.ThreadToString([]interface{}{inObj}).(*object.Object)
				strBuffer = object.GoStringFromStringObject(strObj)
				break
			}
			if object.IsObjectClass(inObj, types.ClassNameThreadState) {
				strObj := javaLang.ThreadStateToString([]interface{}{inObj}).(*object.Object)
				strBuffer = object.GoStringFromStringObject(strObj)
				break
			}
			classNameSuffix := object.GetClassNameSuffix(inObj, true)
			strBuffer = classNameSuffix + "{"
			for name, field := range inObj.FieldTable {
				strBuffer += fmt.Sprintf("%s=%s, ", name, object.StringifyAnythingGo(field))
			}
			strBuffer = strBuffer[:len(strBuffer)-2] + "}"
			if newLine {
				fmt.Fprintln(writer, strBuffer)
				return nil
			} else {
				fmt.Fprint(writer, strBuffer)
				return nil
			}
		default:
			errMsg := fmt.Sprintf("_printObject: Unsupported parameter type: %T", params[1])
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
	}

	if newLine {
		fmt.Fprintln(writer, strBuffer)
	} else {
		fmt.Fprint(writer, strBuffer)
	}

	return nil
}

// Print an Object's contents
// "java/io/PrintStream.print(Ljava/lang/Object;)V"
func PrintObject(params []interface{}) interface{} {
	// Check for null object.
	if params[1] == nil || object.IsNull(params[1]) {
		writer, ok := params[0].(io.Writer)
		if !ok {
			errMsg := fmt.Sprintf("PrintObject: Expected io.Writer, observed %T", params[0])
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
		fmt.Fprint(writer, types.NullString)
		return nil
	}

	// Check for linked list object.
	if object.GoStringFromStringPoolIndex(params[1].(*object.Object).KlassName) == types.ClassNameLinkedList {
		return _printLinkedList(params, false)
	}

	// It's some other object.
	return _printObject(params, false)
}

// Println an Object's contents
// "java/io/PrintStream.println(Ljava/lang/Object;)V"
func PrintlnObject(params []interface{}) interface{} {
	// Check for null object.
	if params[1] == nil || object.IsNull(params[1]) {
		writer, ok := params[0].(io.Writer)
		if !ok {
			errMsg := fmt.Sprintf("PrintlnObject: Expected io.Writer, observed %T", params[0])
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
		}
		fmt.Fprintln(writer, types.NullString)
		return nil
	}

	// Check for linked list object.
	if object.GoStringFromStringPoolIndex(params[1].(*object.Object).KlassName) == types.ClassNameLinkedList {
		return _printLinkedList(params, true)
	}

	// It's some other object.
	return _printObject(params, true)
}

// Print a linked list like this: [A, B, C]
func _printLinkedList(params []interface{}, newLine bool) interface{} {
	writer, ok := params[0].(io.Writer)
	if !ok {
		errMsg := fmt.Sprintf("_printLinkedList: Expected io.Writer, observed %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	var strBuffer string

	// Get linked list object.
	param1, ok := params[1].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("_printLinkedList: Expected params[1] of type *object.Object but observed type %T\n", params[1])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Handle null LinkedList objects.
	if param1 == nil || object.IsNull(param1) {
		strBuffer = types.NullString
	} else {
		// Get value field, holding the linked list reference.
		fld, ok := param1.FieldTable["value"]
		if !ok || fld.Ftype != types.LinkedList {
			className := object.GoStringFromStringPoolIndex(param1.KlassName)
			errMsg := fmt.Sprintf("_printLinkedList: Class %s (LinkedList?), \"value\" field is missing or is not a LinkedList", className)
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
		}

		// Get the linked list.
		llst, gerr := ghelpers.GetLinkedListFromObject(param1)
		if gerr != nil {
			return gerr
		}

		// Start with the front element.
		// Continue to the end.
		element := llst.Front()
		fmt.Fprint(writer, "[")
		for ix := 0; ix < llst.Len(); ix++ {
			strBuffer += object.StringifyAnythingGo(element.Value)
			strBuffer += ", "
			element = element.Next()
		}
	}

	// Remove the final ", " and add a closing right bracket.
	strBuffer = strBuffer[:len(strBuffer)-2] + "]"

	if newLine {
		fmt.Fprintln(writer, strBuffer)
	} else {
		fmt.Fprint(writer, strBuffer)
	}

	return nil
}
