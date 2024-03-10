/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/exceptions"
	"jacobin/object"
	"jacobin/types"
	"math"
	"os"
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

var MethodSignatures = make(map[string]GMeth)

func Load_Io_PrintStream() map[string]GMeth {
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
			GFunction:  PrintlnI,
		}
	MethodSignatures["java/io/PrintStream.println(C)V"] = // println char
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnC,
		}
	MethodSignatures["java/io/PrintStream.println(I)V"] = // println int
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnI,
		}
	MethodSignatures["java/io/PrintStream.println(S)V"] = // println short
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnI,
		}
	MethodSignatures["java/io/PrintStream.println(Z)V"] = // println boolean
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintlnBoolean,
		}
	MethodSignatures["java/io/PrintStream.println(J)V"] = // println long
		GMeth{
			ParamSlots: 2, // 2 slots for the long
			GFunction:  PrintlnLong,
		}

	MethodSignatures["java/io/PrintStream.println(D)V"] = // println double
		GMeth{
			ParamSlots: 2, // 2 slots for the double
			GFunction:  PrintlnDouble,
		}

	MethodSignatures["java/io/PrintStream.println(F)V"] = // println float
		GMeth{
			ParamSlots: 1, // 1 slot for the float
			GFunction:  PrintlnDouble,
		}

	MethodSignatures["java/io/PrintStream.println(Ljava/lang/Object;)V"] = // println object
		GMeth{
			ParamSlots: 1, // 1 slot for the Object
			GFunction:  PrintlnObject,
		}

	MethodSignatures["java/io/PrintStream.print(Ljava/lang/String;)V"] = // print string
		GMeth{
			ParamSlots: 1, // [0] =  StringConst to print
			GFunction:  PrintString,
		}
	MethodSignatures["java/io/PrintStream.print(B)V"] = // print byte
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintI,
		}
	MethodSignatures["java/io/PrintStream.print(C)V"] = // print char
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintC,
		}
	MethodSignatures["java/io/PrintStream.print(I)V"] = // print int
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintI,
		}
	MethodSignatures["java/io/PrintStream.print(S)V"] = // print short
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintI,
		}
	MethodSignatures["java/io/PrintStream.print(Z)V"] = // print boolean
		GMeth{
			ParamSlots: 1,
			GFunction:  PrintBoolean,
		}
	MethodSignatures["java/io/PrintStream.print(J)V"] = // print long
		GMeth{
			ParamSlots: 2, // 2 slots for the long
			GFunction:  PrintLong,
		}

	MethodSignatures["java/io/PrintStream.print(D)V"] = // print double
		GMeth{
			ParamSlots: 2, // 2 slots for the double
			GFunction:  PrintDouble,
		}

	MethodSignatures["java/io/PrintStream.print(F)V"] = // print float
		GMeth{
			ParamSlots: 1, // 1 slot for the float
			GFunction:  PrintFloat,
		}

	MethodSignatures["java/io/PrintStream.print(Ljava/lang/Object;)V"] = // print object
		GMeth{
			ParamSlots: 1, // 1 slot for the Object
			GFunction:  PrintObject,
		}

	MethodSignatures["java/io/PrintStream.printf(Ljava/lang/String;[Ljava/lang/Object;)Ljava/io/PrintStream;"] =
		GMeth{
			ParamSlots: 2, // the format string, the parameters (if any)
			GFunction:  Printf,
		}

	return MethodSignatures
}

func PrintlnString(params []interface{}) interface{} {
	var str string
	switch params[1].(type) {
	case *object.Object:
		obj := params[1].(*object.Object)
		switch obj.FieldTable["value"].Ftype {
		case types.StringIndex:
			str = object.GetGoStringFromObject(obj)
		case types.ByteArray:
			str = string(obj.FieldTable["value"].Fvalue.([]byte))
		}
	default:
		errMsg := fmt.Sprintf("PrintlnString: expected params[1] of type *object.Object but observed type %T\n", params[1])
		return getGErrBlk(exceptions.IllegalArgumentException, errMsg)
	}
	fmt.Fprintln(params[0].(*os.File), str)
	return nil
}

// PrintlnV = java/io/Prinstream.println() -- println() prints a newline (V = void)
func PrintlnV(params []interface{}) interface{} {
	fmt.Fprintln(params[0].(*os.File), "")
	return nil
}

// PrintlnC = java/io/Prinstream.println(char)
func PrintlnC(params []interface{}) interface{} {
	cc := fmt.Sprint(params[1].(int64))
	fmt.Fprintln(params[0].(*os.File), cc)
	return nil
}

// PrintlnI = java/io/Prinstream.println(int)
func PrintlnI(params []interface{}) interface{} {
	intToPrint := params[1].(int64) // contains an int
	fmt.Fprintln(params[0].(*os.File), intToPrint)
	return nil
}

// PrintlnBoolean = java/io/Prinstream.println(boolean)
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

// PrintlnLong = java/io/Prinstream.println(long)
// Long in Java are 64-bit ints, so we just duplicated the logic for println(int)
func PrintlnLong(params []interface{}) interface{} {
	longToPrint := params[1].(int64) // contains to an int64--the equivalent of a Java long
	fmt.Fprintln(params[0].(*os.File), longToPrint)
	return nil
}

// PrintlnDouble = java/io/Prinstream.println(double) or java/io/Prinstream.println(float)
// Doubles in Java are 64-bit FP. Like Hotspot, we print at least one decimal place of data.
func PrintlnDouble(params []interface{}) interface{} {
	doubleToPrint := params[1].(float64) // contains to a float64--the equivalent of a Java double
	fmt.Fprintf(params[0].(*os.File), getDoubleFormat(doubleToPrint)+"\n", doubleToPrint)
	return nil
}

// Println an Object's contents
func PrintlnObject(params []interface{}) interface{} {
	objPtr := params[1].(*object.Object)
	fld := objPtr.FieldTable["value"]
	if fld.Ftype == types.ByteArray {
		fmt.Fprintln(params[0].(*os.File), string(fld.Fvalue.([]byte)))
		return nil
	}
	fmt.Fprintln(params[0].(*os.File), fld.Fvalue)
	return nil
}

// PrintC = java/io/Prinstream.print(char)
func PrintC(params []interface{}) interface{} {
	cc := fmt.Sprint(params[1].(int64))
	fmt.Fprint(params[0].(*os.File), cc)
	return nil
}

// PrintI = java/io/Prinstream.print(int)
func PrintI(params []interface{}) interface{} {
	intToPrint := params[1].(int64) // contains an int
	fmt.Fprint(params[0].(*os.File), intToPrint)
	return nil
}

// PrintBoolean = java/io/Prinstream.print(boolean)
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
func PrintLong(params []interface{}) interface{} {
	longToPrint := params[1].(int64) // contains to an int64--the equivalent of a Java long
	fmt.Fprint(params[0].(*os.File), longToPrint)
	return nil
}

// Printfloat = java/io/Prinstream.print(float)
// Doubles in Java are 64-bit FP
func PrintFloat(params []interface{}) interface{} {
	floatToPrint := params[1].(float64) // contains to a float64--the equivalent of a Java double
	fmt.Fprintf(params[0].(*os.File), getDoubleFormat(floatToPrint), floatToPrint)
	return nil
}

// PrintDouble = java/io/Prinstream.print(double)
// Doubles in Java are 64-bit FP
func PrintDouble(params []interface{}) interface{} {
	doubleToPrint := params[1].(float64) // contains to a float64--the equivalent of a Java double
	fmt.Fprintf(params[0].(*os.File), getDoubleFormat(doubleToPrint), doubleToPrint)
	return nil
}

// Print string
func PrintString(params []interface{}) interface{} {
	var str string
	switch params[1].(type) {
	case *object.Object:
		obj := params[1].(*object.Object)
		switch obj.FieldTable["value"].Ftype {
		case types.StringIndex:
			str = object.GetGoStringFromObject(obj)
		case types.ByteArray:
			str = string(obj.FieldTable["value"].Fvalue.([]byte))
		}
	default:
		errMsg := fmt.Sprintf("PrintString: expected params[1] of type *object.Object but observed type %T\n", params[1])
		return getGErrBlk(exceptions.IllegalArgumentException, errMsg)
	}
	fmt.Fprint(params[0].(*os.File), str)
	return nil
}

// Print an Object's contents
func PrintObject(params []interface{}) interface{} {
	objPtr := params[1].(*object.Object)
	fld := objPtr.FieldTable["value"]
	switch fld.Ftype {
	case types.StringIndex:
		str := object.GetGoStringFromObject(objPtr)
		fmt.Fprintln(params[0].(*os.File), str)
	case types.ByteArray:
		fmt.Fprint(params[0].(*os.File), string(fld.Fvalue.([]byte)))
	default:
		fmt.Fprint(params[0].(*os.File), fld.Fvalue)
	}
	return nil
}

// Printf -- handle the variable args and then call golang's own printf function
func Printf(params []interface{}) interface{} {
	var intfSprintf = new([]interface{})
	*intfSprintf = append(*intfSprintf, params[1])
	*intfSprintf = append(*intfSprintf, params[2])
	retval := StringFormatter(*intfSprintf)
	switch retval.(type) {
	case *object.Object:
	default:
		return retval
	}
	objPtr := retval.(*object.Object)
	str := object.GetGoStringFromObject(objPtr)
	fmt.Fprint(params[0].(*os.File), str)
	return params[0] // Return the PrintStream object

}

// Trying to approximate the exact formatting used in HotSpot JVM
// TODO: look at the JDK source code to map this formatting exactly.
func getDoubleFormat(d float64) string {
	if d < 0.0000001 || d > 10_000_000 {
		return "%E"
	} else {
		if d == math.Floor(d) { // if the fractional part is 0, print a trailing .0
			return "%.01f"
		} else {
			return "%f"
		}
	}
}
