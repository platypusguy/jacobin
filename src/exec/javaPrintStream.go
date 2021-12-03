/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package exec

import (
	"fmt"
)

/*
 Each object or library that has Go methods contains a reference to MethodSignatures,
 which contain data needed to insert the go method into the vtable of the currently
 executing JVM. MethodSignatures is a map whose key is the fully qualified name and
 type of the method (that is, the method's full signature) and a value consisting of
 a struct of an int (the number of slots to pop off the caller's operand stack when
 creating the new frame and a function. All methods have the same signature, regardless
 of the signature of their Java counterparts. That signature is that it accepts a slice
 of interface{} and returns nothing.

 The slice contains one entry for every parameter passed to the method (which could
 mean an empty slice). There is no return value, because the method will place any
 return value on the operand stack of the calling function.
*/

var MethodSignatures = make(map[string]GMeth)

type GMeth struct {
	ParamSlots int
	GFunction  function
}

type function func([]interface{})

func Load_System_PrintStream() map[string]GMeth {
	MethodSignatures["java/io/PrintStream.println(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2, // [0] = PrintStream.out object,
			// [1] = index to StringConst to print
			GFunction: Println,
		}
	MethodSignatures["java/io/PrintStream.println(I)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  PrintlnI,
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
func Println(i []interface{}) {
	sIndex := i[1].(int64) // points to a String constant entry in the CP
	cpi := i[0].(int64)    // int64 which is an index into Statics array
	cp := StaticsArray[cpi].CP
	s := FetchUTF8stringFromCPEntryNumber(cp, uint16(sIndex))
	fmt.Println(s)
}

// java/io/Prinstream(int) TODO: equivalent (verify that this grabs the right param to print)
func PrintlnI(i []interface{}) {
	intToPrint := i[1].(int64) // points to an int
	// cpi := i[0].(int64)    // int64 which is an index into Statics array
	// cp := StaticsArray[cpi].CP
	// s := FetchUTF8stringFromCPEntryNumber(cp, uint16(sIndex))
	fmt.Println(intToPrint)
}
