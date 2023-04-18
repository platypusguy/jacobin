/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"fmt"
	"unsafe"
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

type GMeth struct {
	ParamSlots int
	GFunction  function
}

type function func([]interface{}) interface{}

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
			GFunction:  PrintDouble,
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
	// sIndex := i[1].(int64) // points to a String constant entry in the CP
	// // cpi := i[0].(int64)    // int64 which is an index into Statics array
	// // cp := StaticsArray[cpi].CP
	//
	// usIndex := uint64(sIndex)
	// upsIndex := uintptr(usIndex)
	// strAddr := unsafe.Pointer(upsIndex)
	strAddr := i[1].(unsafe.Pointer)
	s := *(*string)(strAddr)
	// s := FetchUTF8stringFromCPEntryNumber(cp, uint16(sIndex))
	fmt.Println(s)
	return nil
}

// PrintlnV = java/io/Prinstream.println() -- println() prints a newline
func PrintlnV(i []interface{}) interface{} {
	// intToPrint := i[1].(int64) // contains an int
	fmt.Println("")
	return nil
}

// PrintlnI = java/io/Prinstream.println(int) TODO: equivalent (verify that this grabs the right param to print)
func PrintlnI(i []interface{}) interface{} {
	intToPrint := i[1].(int64) // contains an int
	fmt.Println(intToPrint)
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
// Doubles in Java are 64-bit FP
func PrintlnDouble(l []interface{}) interface{} {
	doubleToPrint := l[1].(float64) // contains to a float64--the equivalent of a Java double
	fmt.Println(doubleToPrint)
	return nil
}

// PrintI = java/io/Prinstream.print(int) TODO: equivalent (verify that this grabs the right param to print)
func PrintI(i []interface{}) interface{} {
	intToPrint := i[1].(int64) // contains an int
	fmt.Print(intToPrint)
	return nil
}

// PrintLong = java/io/Prinstream.print(long)
// Long in Java are 64-bit ints, so we just duplicated the logic for println(int)
func PrintLong(l []interface{}) interface{} {
	longToPrint := l[1].(int64) // contains to an int64--the equivalent of a Java long
	fmt.Print(longToPrint)
	return nil
}

// PrintDouble = java/io/Prinstream.print(double)
// Doubles in Java are 64-bit FP
func PrintDouble(l []interface{}) interface{} {
	doubleToPrint := l[1].(float64) // contains to a float64--the equivalent of a Java double
	fmt.Print(doubleToPrint)
	return nil
}

// Print string
func PrintS(i []interface{}) interface{} {
	// sIndex := i[1].(int64) // points to a String constant entry in the CP
	//
	// usIndex := uint64(sIndex)
	// upsIndex := uintptr(usIndex)
	// strAddr := unsafe.Pointer(upsIndex)
	strAddr := i[1].(unsafe.Pointer)
	s := *(*string)(strAddr)
	fmt.Print(s)
	return nil
}
