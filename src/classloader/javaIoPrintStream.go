/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"fmt"
	"jacobin/exceptions"
	"jacobin/object"
	"math"
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
	MethodSignatures["java/io/PrintStream.println()V"] = // println string
		GMeth{
			ParamSlots: 1, // [0] = PrintStream.out object,
			GFunction:  PrintlnV,
		}
	MethodSignatures["java/io/PrintStream.println(Ljava/lang/String;)V"] = // println string
		GMeth{
			ParamSlots: 2, // [0] = PrintStream.out object,
			// [1] = index to StringConst to print
			GFunction: Println,
		}
	MethodSignatures["java/io/PrintStream.println(I)V"] = // println int
		GMeth{
			ParamSlots: 2,
			GFunction:  PrintlnI,
		}
	MethodSignatures["java/io/PrintStream.println(Z)V"] = // println boolean
		GMeth{
			ParamSlots: 2,
			GFunction:  PrintlnBoolean,
		}
	MethodSignatures["java/io/PrintStream.println(J)V"] = // println long
		GMeth{
			ParamSlots: 3, // PrintStream.out object + 2 slots for the long
			GFunction:  PrintlnLong,
		}

	MethodSignatures["java/io/PrintStream.println(D)V"] = // println double
		GMeth{
			ParamSlots: 3, // PrintStream.out object + 2 slots for the double
			GFunction:  PrintlnDouble,
		}

	MethodSignatures["java/io/PrintStream.println(F)V"] = // println float
		GMeth{
			ParamSlots: 2, // PrintStream.out object + 1 slot for the float
			GFunction:  PrintlnDouble,
		}

	MethodSignatures["java/io/PrintStream.print(Ljava/lang/String;)V"] = // print string
		GMeth{
			ParamSlots: 2, // [0] = PrintStream.out object,
			// [1] = index to StringConst to print
			GFunction: PrintS,
		}
	MethodSignatures["java/io/PrintStream.print(I)V"] = // print int
		GMeth{
			ParamSlots: 2,
			GFunction:  PrintI,
		}
	MethodSignatures["java/io/PrintStream.print(Z)V"] = // print boolean
		GMeth{
			ParamSlots: 2,
			GFunction:  PrintBoolean,
		}
	MethodSignatures["java/io/PrintStream.print(J)V"] = // print long
		GMeth{
			ParamSlots: 3, // PrintStream.out object + 2 slots for the long
			GFunction:  PrintLong,
		}

	MethodSignatures["java/io/PrintStream.print(D)V"] = // print double
		GMeth{
			ParamSlots: 3, // PrintStream.out object + 2 slots for the double
			GFunction:  PrintDouble,
		}

	MethodSignatures["java/io/PrintStream.print(F)V"] = // print float
		GMeth{
			ParamSlots: 2, // PrintStream.out object + 1 slot for the float
			GFunction:  PrintFloat,
		}

	MethodSignatures["java/io/PrintStream.printf(Ljava/lang/String;[Ljava/lang/Object;)Ljava/io/PrintStream;"] =
		GMeth{
			ParamSlots: 3, // the Printstream object, the format string, the parameters (if any)
			GFunction:  Printf,
		}

	return MethodSignatures
}

// Println is the go equivalent of System.out.println(). It accepts two args,
// which are passed in a two-entry slice of type interface{}. The first arg is an
// index in the CP to a StringConst entry; the second arg is an index into the
// array of static fields, Statics. The entry there includes a pointer to the CP
// for this class. The first arg then gets the StringConst ref, which is an index
// into the UTF8 entries of the CP. This string is then printed to stdout. There
// is no return value.
func Println(i []interface{}) interface{} {
	strAddr := i[1].(*object.Object)
	// t := (strAddr.Fields[0].Fvalue).(*[]types.JavaByte) // changed due to JAcOBIN-282
	t := (strAddr.Fields[0].Fvalue).(*[]byte)

	// goChars := make([]byte, len(*t), len(*t))
	// for i, c := range *t {
	// 	goChars[i] = byte(c)
	// }

	fmt.Println(string(*t))
	return nil
}

// PrintlnV = java/io/Prinstream.println() -- println() prints a newline (V = void)
func PrintlnV(i []interface{}) interface{} {
	fmt.Println("")
	return nil
}

// PrintlnI = java/io/Prinstream.println(int) TODO: equivalent (verify that this grabs the right param to print)
func PrintlnI(i []interface{}) interface{} {
	intToPrint := i[1].(int64) // contains an int
	fmt.Println(intToPrint)
	return nil
}

// PrintlnBoolean = java/io/Prinstream.println(boolean) TODO: equivalent (verify that this grabs the right param to print)
func PrintlnBoolean(i []interface{}) interface{} {
	var boolToPrint bool
	boolAsInt64 := i[1].(int64) // contains an int64
	if boolAsInt64 > 0 {
		boolToPrint = true
	} else {
		boolToPrint = false
	}
	fmt.Println(boolToPrint)
	return nil
}

// PrintlnLong = java/io/Prinstream.println(long)
// Long in Java are 64-bit ints, so we just duplicated the logic for println(int)
func PrintlnLong(l []interface{}) interface{} {
	longToPrint := l[1].(int64) // contains to an int64--the equivalent of a Java long
	fmt.Println(longToPrint)
	return nil
}

// PrintlnDouble = java/io/Prinstream.println(double)
// Doubles in Java are 64-bit FP. Like Hotspot, we print at least one decimal place of data.
func PrintlnDouble(l []interface{}) interface{} {
	doubleToPrint := l[1].(float64) // contains to a float64--the equivalent of a Java double
	fmt.Printf(formatDouble(doubleToPrint)+"\n", doubleToPrint)
	return nil
}

// PrintI = java/io/Prinstream.print(int) TODO: equivalent (verify that this grabs the right param to print)
func PrintI(i []interface{}) interface{} {
	intToPrint := i[1].(int64) // contains an int
	fmt.Print(intToPrint)
	return nil
}

// PrintBoolean = java/io/Prinstream.print(boolean) TODO: equivalent (verify that this grabs the right param to print)
func PrintBoolean(i []interface{}) interface{} {
	var boolToPrint bool
	boolAsInt64 := i[1].(int64) // contains an int64
	if boolAsInt64 > 0 {
		boolToPrint = true
	} else {
		boolToPrint = false
	}
	fmt.Print(boolToPrint)
	return nil
}

// PrintLong = java/io/Prinstream.print(long)
// Long in Java are 64-bit ints, so we just duplicated the logic for println(int)
func PrintLong(l []interface{}) interface{} {
	longToPrint := l[1].(int64) // contains to an int64--the equivalent of a Java long
	fmt.Print(longToPrint)
	return nil
}

// Printfloat = java/io/Prinstream.print(float)
// Doubles in Java are 64-bit FP
func PrintFloat(l []interface{}) interface{} {
	floatToPrint := l[1].(float64) // contains to a float64--the equivalent of a Java double
	fmt.Printf(formatDouble(floatToPrint), floatToPrint)
	return nil
}

// PrintDouble = java/io/Prinstream.print(double)
// Doubles in Java are 64-bit FP
func PrintDouble(l []interface{}) interface{} {
	doubleToPrint := l[1].(float64) // contains to a float64--the equivalent of a Java double
	fmt.Printf(formatDouble(doubleToPrint), doubleToPrint)
	return nil
}

// Print string
func PrintS(i []interface{}) interface{} {
	strAddr := i[1].(*object.Object)
	// eventually will need to check wherther it's a compact string.
	// Presently, we assume it is.
	t := (strAddr.Fields[0].Fvalue).(*[]byte)
	fmt.Print(string(*t))
	return nil
}

// Printf -- handle the variable args and then call golang's own printf function
func Printf(params []interface{}) interface{} {
	ps := params[0]
	formatStringObj := params[1].(*object.Object) // the format string is passed as a pointer to a string object
	formatString := object.GetGoStringFromJavaStringPtr(formatStringObj)

	// now peel off the parameters beyond the format string and pass them to golang's printf function
	switch len(params) {
	case 0, 1:
		errMsg := "printf(): Invalid parameter count"
		exceptions.Throw(exceptions.IllegalClassFormatException, errMsg)
	case 2: // 0 parameters beyond the format string
		fmt.Printf(formatString)
	case 3: // 1 parameter beyond the format string, which will be an array of pointers to objects
		valuesIn := (params[2].(*object.Object).Fields[0].Fvalue).([]*object.Object) // array of pointers to 1 or more objects
		valuesOut := []any{}
		for i := 0; i < len(valuesIn); i++ {
			value := getRawParameter(valuesIn[i])
			valuesOut = append(valuesOut, value)
		}

		switch len(valuesOut) {
		case 1:
			fmt.Printf(formatString, valuesOut[0])
		case 2:
			fmt.Printf(formatString, valuesOut[0], valuesOut[1])
		case 3:
			fmt.Printf(formatString, valuesOut[0], valuesOut[1], valuesOut[2])
		case 4:
			fmt.Printf(formatString, valuesOut[0], valuesOut[1], valuesOut[2], valuesOut[3])
		}
	}
	//
	//
	// 	var param1 any
	// 	if object.IsJavaString(params[2]) {
	// 		param1 = object.GetGoStringFromJavaStringPtr(params[2].(*object.Object))
	// 	} else {
	// 		param1 = params[2]
	// 	}
	// 	fmt.Printf(formatString, param1)
	// case 4: // 2 parameters beyond the format string
	// 	var param1, param2 any
	// 	if object.IsJavaString(params[2]) {
	// 		param1 = object.GetGoStringFromJavaStringPtr(params[2].(*object.Object))
	// 	} else {
	// 		param1 = params[2]
	// 	}
	// 	if object.IsJavaString(params[3]) {
	// 		param2 = object.GetGoStringFromJavaStringPtr(params[3].(*object.Object))
	// 	} else {
	// 		param2 = params[3]
	// 	}
	// 	fmt.Printf(formatString, param1, param2)
	// case 5: // 3 parameters beyond the format string
	// 	var param1, param2, param3 any
	// 	if object.IsJavaString(params[2]) {
	// 		param1 = object.GetGoStringFromJavaStringPtr(params[2].(*object.Object))
	// 	} else {
	// 		param1 = params[2]
	// 	}
	// 	if object.IsJavaString(params[3]) {
	// 		param2 = object.GetGoStringFromJavaStringPtr(params[3].(*object.Object))
	// 	} else {
	// 		param2 = params[3]
	// 	}
	// 	if object.IsJavaString(params[4]) {
	// 		param3 = object.GetGoStringFromJavaStringPtr(params[4].(*object.Object))
	// 	} else {
	// 		param3 = params[4]
	// 	}
	// 	fmt.Printf(formatString, param1, param2, param3)
	// }

	return ps // return the printStream (even though we don't use it here)
}

func getRawParameter(param any) any {
	if object.IsJavaString(param) {
		return object.GetGoStringFromJavaStringPtr(param.(*object.Object))
	} else {
		return param
	}
}

// Trying to approximate the exact formatting used in HotSpot JVM
// TODO: look at the JDK source code to map this formatting exactly.
func formatDouble(d float64) string {
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
