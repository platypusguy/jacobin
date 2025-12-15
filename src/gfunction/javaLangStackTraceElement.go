/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"container/list"
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/exceptions"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/shutdown"
	"jacobin/src/trace"
	"jacobin/src/util"
	"sort"
)

// StackTraceElement is a class that is primarily used by Throwable to gather data about the
// entries in the JVM stack. Because this data is so tightly bound to the specific implementation
// of the JVM, the methods of this class are fairly faithfully reproduced in the golang-native
// methods in this file. Consult:
// https://docs.oracle.com/en/java/javase/17/docs/api/java.base/java/lang/StackTraceElement.html

func Load_Lang_StackTraceELement() {

	MethodSignatures["java/lang/StackTraceElement.of(Ljava/lang/Throwable;I)[Ljava/lang/StackTraceElement;"] =
		GMeth{
			ParamSlots:   2,
			GFunction:    of,
			NeedsContext: true,
		}

	MethodSignatures["java/lang/StackTraceElement.initStackTraceElements([Ljava/lang/StackTraceElement;Ljava/lang/Throwable;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  initStackTraceElements,
		}

}

/*
	 From the Java code for this method:

	 Returns an array of StackTraceElements of the given depth
	 filled from the backtrace of a given Throwable.

	    static StackTraceElement[] of(Throwable x, int depth) {
			StackTraceElement[] stackTrace = new StackTraceElement[depth];
			for (int i = 0; i < depth; i++) {
				stackTrace[i] = new StackTraceElement();
		  }

		// VM to fill in StackTraceElement
		initStackTraceElements(stackTrace, x);

		// ensure the proper StackTraceElement initialization
		for (StackTraceElement ste : stackTrace) {
			ste.computeFormat();
		}
		return stackTrace;
	}
*/
func of(params []interface{}) interface{} {

	throwable := params[0].(*object.Object)
	depth := params[1].(int64)

	// get a pointer to the JVM stack
	jvmStackRef := throwable.FieldTable["frameStackRef"].Fvalue.(*list.List)
	if jvmStackRef == nil {
		errMsg := "java/lang/StackTraceElement.of: Nil parameter for 'frameStackRef' in Throwable, found in StackTraceElement.of()"
		trace.Error(errMsg)
		shutdown.Exit(shutdown.JVM_EXCEPTION)
	}

	// create the 1-dimensional array of stackTraceElements
	stackTrace := object.Make1DimRefArray("java/lang/StackTraceElement", depth)

	// insert empty stackTraceElements into the array.
	rawArray := stackTrace.FieldTable["value"].Fvalue.([]*object.Object)
	global := globals.GetGlobalRef()
	for i := int64(0); i < depth; i++ {
		ste, err := global.FuncInstantiateClass("java/lang/StackTraceElement", nil)
		if err == nil {
			rawArray[i] = ste.(*object.Object)
		}
	}

	argsToPass := []interface{}{stackTrace, throwable}
	initStackTraceElements(argsToPass)

	return stackTrace
}

// This is a native function in HotSpot that accepts an array of empty
// stackTraceElements and a Throwable and fills in the values in the array
// by repeated calls to initStackTraceElement() below.
// Returns nothing.
func initStackTraceElements(params []interface{}) interface{} {
	arrayObjPtr := params[0].(*object.Object) // the array of stackTraceElements we'll fill in
	arrayObj := *arrayObjPtr
	rawSteArray := arrayObj.FieldTable["value"].Fvalue.([]*object.Object)
	// rawSteArray := *rawSteArrayPtr

	throwable := params[1].(*object.Object) // pointer to the Throwable object
	jvmStack := throwable.FieldTable["frameStackRef"].Fvalue.(*list.List)

	var i = 0
	isFirstFrame := true
	for e := jvmStack.Front(); e != nil; e = e.Next() {
		frame := e.Value.(*frames.Frame)

		// Skip printStackTrace and related printing methods
		if frame.MethName == "printStackTrace" ||
			frame.MethName == "printStackTraceToPrintStream" ||
			frame.MethName == "printStackTraceToPrintWriter" {
			continue
		}

		ste := rawSteArray[i]
		i += 1
		initStackTraceElement(ste, frame, isFirstFrame)
		isFirstFrame = false
	}

	return nil
}

// initStackTraceElement accepts a single stackTraceElement and JVM stack
// info and fills in the former with the latter. It's a private method and
// called only from initStackTraceElements(), so we don't need it to strictly
// follow the HotSpot way of implementing it. Official definition:
// initStackTraceElement(Ljava/lang/StackTraceElement;Ljava/lang/StackFrameInfo;)V
func initStackTraceElement(ste *object.Object, frm *frames.Frame, isFirstFrame bool) {
	frame := *frm
	stackTrace := *ste

	// helper function to facilitate subsequent field updates
	// (Thanks to JetBrains' AI Assistant for this suggestion)
	addField := func(name, value string) {
		fld := object.Field{}
		fld.Fvalue = value
		stackTrace.FieldTable[name] = fld
	}

	addField("declaringClass", frame.ClName)
	addField("methodName", frame.MethName)

	methClass := classloader.MethAreaFetch(frame.ClName)

	addField("classLoaderName", methClass.Loader)
	addField("fileName", methClass.Data.SourceFile)
	addField("moduleName", methClass.Data.Module)

	// now get the source line number for any non-JDK classes and non-constructors
	// Unsure why this limitation. It's commented out for the moment (JACOBIN-781)

	addField("sourceLine", "") // the default if no source line data is available
	// if !util.IsFilePartOfJDK(&frame.MethName) && !strings.HasPrefix(frame.MethName, "<init>") {
	rawMethod, _ := classloader.FetchMethodAndCP(frame.ClName, frame.MethName, frame.MethType)
	if rawMethod.MType == 'G' { // nothing more to do if it's a native method
		return
	}
	method, ok := rawMethod.Meth.(classloader.JmEntry)
	if !ok {
		errMsg := fmt.Sprintf("initStackTraceElement: %s.%s, Invalid operand type for rawMethod.Meth: %T",
			util.ConvertInternalClassNameToUserFormat(frame.ClName), frame.MethName, rawMethod.Meth)
		_ = exceptions.ThrowEx(excNames.InternalException, errMsg, &frame)
	}
	for i := 0; i < len(method.Attribs); i++ {
		index := method.Attribs[i].AttrName
		if method.Cp.Utf8Refs[index] == "LineNumberTable" {
			// Use ExceptionPC if it's set (not -1), otherwise fall back to PC
			pcToUse := frame.PC
			if frame.ExceptionPC != -1 {
				pcToUse = frame.ExceptionPC
			}

			// For non-first frames (caller frames), the PC points to the instruction
			// after the method call. We need to look back to find the actual call instruction.
			// Most invoke instructions are 3 bytes (opcode + 2-byte CP index).
			if !isFirstFrame && pcToUse > 0 {
				// Adjust PC backward to point to the call instruction
				// We subtract 1 to get into the range of the previous instruction
				pcToUse = pcToUse - 1
			}

			line := searchLineNumberTable(method.Attribs[i].AttrContent, pcToUse)
			if line != -1 { // -1 means not found
				addField("sourceLine", fmt.Sprintf("%d", line))
			}
		}
	}
	// }
}

// get the source line number from the location of the bytecode where exception occurred
//
// We first create a table of entries consisting of bytecode number and source line number
// then we sort the table, then we traverse the table to find the matching line number
func searchLineNumberTable(attrContent []byte, PC int) int {
	entryCount := uint(attrContent[0])*256 + uint(attrContent[1])
	loc := 2 // 2 bytes into attrContent for entry count uint16 entryCount
	if entryCount < 1 {
		return -1
	}

	var table b2sTable
	var i uint

	// build the table
	for i = 0; i < entryCount; i++ {
		bytecodeNumber := uint16(attrContent[loc])*256 + uint16(attrContent[loc+1])
		sourceLineNumber := uint16(attrContent[loc+2])*256 + uint16(attrContent[loc+3])
		loc += 4

		tableEntry := BytecodeToSourceLine{bytecodeNumber, sourceLineNumber}
		table = append(table, tableEntry)
	}

	// sort the table
	if len(table) > 1 {
		sort.Sort(b2sTable(table))
	}

	// traverse the table
	var prev uint16 = 0
	for i = 0; i < entryCount; i++ {
		bytecodeNumber := table[i].BytecodePos
		sourceLineNumber := table[i].SourceLine
		if bytecodeNumber > uint16(PC) {
			break
		} else if bytecodeNumber == uint16(PC) {
			prev = sourceLineNumber
			break
		} else {
			prev = sourceLineNumber
		}
	}
	return int(prev)
}

// the following four lines are all needed for the call to Sort()
type b2sTable []BytecodeToSourceLine

func (t b2sTable) Len() int           { return len(t) }
func (t b2sTable) Swap(k, j int)      { (t)[k], (t)[j] = (t)[j], (t)[k] }
func (t b2sTable) Less(k, j int) bool { return (t)[k].BytecodePos < (t)[j].BytecodePos }

// BytecodeToSourceLine maps the PC in a method to the
// corresponding source line in the original source file.
type BytecodeToSourceLine struct {
	BytecodePos uint16
	SourceLine  uint16
}
