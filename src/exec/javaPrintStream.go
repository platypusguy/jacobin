/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package exec

import (
	"fmt"
	"os"
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
			ParamSlots: 2, // [0] = PrintStream.out object, [1] = string to print
			GFunction:  Println,
		}
	return MethodSignatures
}

// a temporary stand-in for java\io\PrintStream
type stream *os.File

var Out stream

func PrintStream(out stream) {
	Out = out
}

func init() {
	Out = os.Stdout
}

func Println(i []interface{}) {
	sIndex := i[1].(int64) // points to a String constant entry in the CP
	cpi := i[0].(int64)    // int64 which is an index into Statics array
	cp := StaticsArray[cpi].CP
	stringRef := cp.CpIndex[sIndex]
	utf8index := cp.StringRefs[int(stringRef.Slot)]
	s := FetchUTF8stringFromCPEntryNumber(cp, uint16(utf8index))
	fmt.Fprintf(os.Stdout, "%s\n", s)
}
