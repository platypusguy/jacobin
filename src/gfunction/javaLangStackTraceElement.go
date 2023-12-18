/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"container/list"
	"fmt"
	"jacobin/classloader"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/object"
	"jacobin/shutdown"
	"jacobin/util"
	// "sort"
	"strings"
)

// StackTraceElement is a class that is primarily used by Throwable to gather data about the
// entries in the JVM stack. Because this data is so tightly bound to the specific implementation
// of the JVM, the methods of this class are fairly faithfully reproduced in the golang-native
// methods in this file. Consult:
// https://docs.oracle.com/en/java/javase/17/docs/api/java.base/java/lang/StackTraceElement.html

func Load_Lang_StackTraceELement() map[string]GMeth {

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

	return MethodSignatures
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
		errMsg := "nil parameter for 'frameStackRef' in Throwable, found in StackTraceElement.of()"
		_ = log.Log(errMsg, log.SEVERE)
		shutdown.Exit(shutdown.JVM_EXCEPTION)
	}

	// create the 1-dimensional array of stackTraceElements
	stackTraceElementClassName := "java/lang/StackTraceElement"
	stackTrace := object.Make1DimRefArray(&stackTraceElementClassName, depth)

	// insert empty stackTraceElements into the array.
	rawArrayPtr := stackTrace.Fields[0].Fvalue.(*[]*object.Object)
	rawArray := *rawArrayPtr
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
	rawSteArrayPtr := arrayObj.Fields[0].Fvalue.(*[]*object.Object)
	rawSteArray := *rawSteArrayPtr

	throwable := params[1].(*object.Object) // pointer to the Throwable object
	jvmStack := throwable.FieldTable["frameStackRef"].Fvalue.(*list.List)

	var i = 0
	for e := jvmStack.Front(); e != nil; e = e.Next() {
		frame := e.Value.(*frames.Frame)
		ste := rawSteArray[i]
		i += 1
		initStackTraceElement(ste, frame)
	}

	return nil

	// Note: an improvement is to add the logic to the class parse showing
	// the Java code source line number. JACOBIN-224 refers to this.
}

// initStackTraceElement accepts a single stackTraceElement and JVM stack
// info and fills in the former with the latter. It's a private method and
// called only from initStackTraceElements(), so we don't need it to strictly
// follow the HotSpot way of implementing it. Official definition:
// initStackTraceElement(Ljava/lang/StackTraceElement;Ljava/lang/StackFrameInfo;)V
func initStackTraceElement(ste *object.Object, frm *frames.Frame) {
	frame := *frm
	stackTrace := *ste

	// helper function to facilitate subsequent field updates
	// (Thanks to JetBrains' AI Assistant for this suggestion)
	addField := func(name, value string) {
		fld := object.Field{}
		fld.Fvalue = value
		stackTrace.FieldTable[name] = &fld
	}

	addField("declaringClass", frame.ClName)
	addField("methodName", frame.MethName)

	methClass := classloader.MethAreaFetch(frame.ClName)

	addField("classLoaderName", methClass.Loader)
	addField("fileName", methClass.Data.SourceFile)
	addField("moduleName", methClass.Data.Module)

	// now get the source line number for any non-JDK files
	if util.IsFilePartOfJDK(&frame.MethName) || strings.HasPrefix(frame.MethName, "<init>") {
		addField("sourceLine", "")
	} else {
		rawMethod, _ := classloader.FetchMethodAndCP(frame.ClName, frame.MethName, frame.MethType)
		if rawMethod.MType == 'G' { // nothing more to do if it's a native method
			return
		}
		method := rawMethod.Meth.(classloader.JmEntry)
		for i := 0; i < len(method.Attribs); i++ {
			index := method.Attribs[i].AttrName
			if method.Cp.Utf8Refs[index] == "LineNumberTable" {
				line := searchLineNumberTable(method.Attribs[i].AttrContent, frame.PC)
				if line != -1 { // -1 means not found
					addField("sourceLine", fmt.Sprintf("%d", line))
				}
				// fmt.Fprintf(os.Stderr, "line: %d\n", line)
			}
		}
	}
}

// get the source line number from the location of the bytecode where exception occurred
func searchLineNumberTable(attrContent []byte, PC int) int {
	entryCount := uint(attrContent[0])*256 + uint(attrContent[1])
	loc := 2 // we're two bytes into the attr.Content byte array
	if entryCount < 1 {
		return -1
	}

	var i uint
	var prev uint16 = 0
	for i = 0; i < entryCount; i++ {
		bytecodeNumber := uint16(attrContent[loc])*256 + uint16(attrContent[loc+1])
		sourceLineNumber := uint16(attrContent[loc+2])*256 + uint16(attrContent[loc+3])
		loc += 4

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

// sourceMap := class.Data.Methods[0].CodeAttr.BytecodeSourceMap
// prev := uint16(0)
// for i := 0; i < len(sourceMap); i++ {
// 	entry := sourceMap[i]
// 	if entry.BytecodePos > uint16(frame.PC) {
// 		break
// 	} else if entry.BytecodePos == uint16(frame.PC) {
// 		prev = entry.SourceLine
// 		break
// 	} else {
// 		prev = entry.SourceLine
// 	}
// }
// sourceLineNumber = fmt.Sprintf("%d", prev)
// addField("sourceLine", sourceLineNumber)

// build the table of line numbers (that map bytecode location to source line #)
// consult https://docs.oracle.com/javase/specs/jvms/se17/html/jvms-4.html#jvms-4.7.12
// func buildLineNumberTable(codeAttr *codeAttrib, thisAttr *attr, methodName string) {
// 	entryCount := uint(thisAttr.attrContent[0])*256 + uint(thisAttr.attrContent[1])
// 	loc := 2 // we're two bytes into the attr.Content byte array
// 	if entryCount < 1 {
// 		(*codeAttr).sourceLineTable = nil
// 		return
// 	}
//
// 	var table []BytecodeToSourceLine
// 	if (*codeAttr).sourceLineTable != nil { // we could be adding to the table
// 		table = []BytecodeToSourceLine{}
// 		(*codeAttr).sourceLineTable = &table
// 	}
// 	var i uint
// 	for i = 0; i < entryCount; i++ {
// 		bytecodeNumber := uint16(thisAttr.attrContent[loc])*256 + uint16(thisAttr.attrContent[loc+1])
// 		sourceLineNumber := uint16(thisAttr.attrContent[loc+2])*256 + uint16(thisAttr.attrContent[loc+3])
// 		loc += 4
//
// 		tableEntry := BytecodeToSourceLine{bytecodeNumber, sourceLineNumber}
// 		table = append(table, tableEntry)
// 	}
//
// 	// now sort the table
// 	if len(table) > 1 {
// 		sort.Sort(b2sTable(table))
// 	}
//
// 	(*codeAttr).sourceLineTable = &table
//
// 	// if methodName == "main" {
// 	// 	fmt.Fprintf(os.Stderr, "%v\n", table)
// 	// }
// }
//
// // the following four lines are all needed for the call to Sort()
// type b2sTable []BytecodeToSourceLine
//
// func (t b2sTable) Len() int           { return len(t) }
// func (t b2sTable) Swap(k, j int)      { (t)[k], (t)[j] = (t)[j], (t)[k] }
// func (t b2sTable) Less(k, j int) bool { return (t)[k].BytecodePos < (t)[j].BytecodePos }
//
// // BytecodeToSourceLine maps the PC in a method to the
// // corresponding source line in the original source file.
// // This data is captured in the method's attributes
// type BytecodeToSourceLine struct {
// 	BytecodePos uint16
// 	SourceLine  uint16
// }
