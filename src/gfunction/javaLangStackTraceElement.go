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
	"jacobin/statics"
)

// StackTraceElement is a class that is primarily used by Throwable to gather data about the
// entries in the JVM stack. Because this data is so tightly bound to the specific implementation
// of the JVM, the methods of this class are fairly faithfully reproduced in the golang-native
// methods in this file.

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

func of(params []interface{}) interface{} {
	stackTraceElementClassName := "java/lang/StackTraceElement"
	emptyStackTraceElementArray := object.Make1DimRefArray(&stackTraceElementClassName, 0)
	statics.AddStatic("Throwable.UNASSIGNED_STACK", statics.Static{
		Type:  "[Ljava/lang/StackTraceElement",
		Value: emptyStackTraceElementArray,
	})

	// for the time being, SUPPRESSED SENTINEL is set to nil.
	// We might later need to set it to an empty List.
	statics.AddStatic("Throwable.SUPPRESSED_SENTINEL", statics.Static{
		Type: "Ljava/util/List", Value: nil})

	emptyThrowableClassName := "java/lang/Throwable"
	emptyThrowableArray := object.Make1DimRefArray(&emptyThrowableClassName, 0)
	statics.AddStatic("Throwable.EMPTY_THROWABLE_ARRAY", statics.Static{
		Type:  "[Ljava/lang/Throwable",
		Value: emptyThrowableArray,
	})
	return nil
}

// This function is called by Throwable.<init>(). In Throwable.java, it consists of one line:
//      getOurStackTrace().clone();
// In turn, getOurStackTrace() calls
//      StackTraceElement.of(this, depth);
// In turn, this method calls
//      StackTraceElement.initStackTraceElements:([Ljava/lang/StackTraceElement;Ljava/lang/Throwable;)V
// which actually fills in the fields of the StackTraceElement (done as a native function)
//
// Despite this simple function chaining, there is value in reading the
// Javadoc for this function from Throwable.java (copyright Oracle Corp.):
/*
 * Provides programmatic access to the stack trace information printed by
 * {@link #printStackTrace()}. Returns an array of stack trace elements,
 * each representing one stack frame. The zeroth element of the array
 * (assuming the array's length is non-zero) represents the top of the
 * stack, which is the last method invocation in the sequence.  Typically,
 * this is the point at which this throwable was created and thrown.
 * The last element of the array (assuming the array's length is non-zero)
 * represents the bottom of the stack, which is the first method invocation
 * in the sequence.
 *
 * <p> [...] Generally speaking, the array returned by this method will
 * contain one element for every frame that would be printed by
 * {@code printStackTrace}.  Writes to the returned array do not
 * affect future calls to this method.
 *
 * @return an array of stack trace elements representing the stack trace
 *         pertaining to this throwable.
 */
func initStackTraceElements(params []interface{}) interface{} {
	if len(params) != 2 {
		_ = log.Log(fmt.Sprintf("fillInsStackTrace() expected two params, got: %d", len(params)), log.SEVERE)
		shutdown.Exit(shutdown.JVM_EXCEPTION)
	}
	frameStack := params[0].(*list.List)
	objRef := params[1].(*object.Object)
	fmt.Printf("Throwable object contains: %v", objRef.FieldTable)

	global := *globals.GetGlobalRef()
	// step through the JVM stack frame and fill in a StackTraceElement for each frame
	for thisFrame := frameStack.Front().Next(); thisFrame != nil; thisFrame = thisFrame.Next() {
		ste, err := global.FuncInstantiateClass("java/lang/StackTraceElement", nil)
		if err != nil {
			_ = log.Log("Throwable.fillInStackTrace: error creating 'java/lang/StackTraceElement", log.SEVERE)
			// return ste.(*object.Object)
			ste = nil
			return ste
		}

		fmt.Println(thisFrame.Value)
	}

	// This might require that we add the logic to the class parse showing the Java code source line number.
	// JACOBIN-224 refers to this.
	return objRef
}

// GetStackTraces gets the full JVM stack trace using java.lang.StackTraceElement
// slice to hold the data. In case of error, nil is returned.
func initStackTraceElement(fs *list.List) *object.Object {
	var stackListing []*object.Object

	frameStack := fs.Front()
	if frameStack == nil {
		// return an empty stack listing
		return nil
	}

	// ...will eventually go into java/lang/Throwable.stackTrace
	// ...Type will be: [Ljava/lang/StackTraceElement;
	// ...other fields to be sure to capture: cause, detailMessage,
	// ....not sure about backtrace

	// step through the list-based stack of called methods and print contents

	var frame *frames.Frame

	for e := frameStack; e != nil; e = e.Next() {
		classname := "java/lang/StackTraceElement"
		stackTrace := object.MakeEmptyObject()
		k := classloader.MethAreaFetch(classname)
		stackTrace.Klass = &classname

		if k == nil {
			errMsg := "Class is nil after loading, class: " + classname
			_ = log.Log(errMsg, log.SEVERE)
			return nil
		}

		if k.Data == nil {
			errMsg := "class.Data is nil, class: " + classname
			_ = log.Log(errMsg, log.SEVERE)
			return nil
		}

		frame = e.Value.(*frames.Frame)

		// helper function to facilitate subsequent field updates
		// thanks to JetBrains' AI Assistant for this suggestion
		addField := func(name, value string) {
			fld := object.Field{}
			fld.Fvalue = value
			stackTrace.FieldTable[name] = &fld
		}

		addField("declaringClass", frame.ClName)
		addField("methodName", frame.MethName)

		methClass := classloader.MethAreaFetch(frame.ClName)
		if methClass == nil {
			return nil
		}
		addField("classLoaderName", methClass.Loader)
		addField("fileName", methClass.Data.SourceFile)
		addField("moduleName", methClass.Data.Module)

		stackListing = append(stackListing, stackTrace)
	}

	// now that we have our data items loaded into the StackTraceElement
	// put the elements into an array, which is converted into an object
	obj := object.MakeEmptyObject()
	klassName := "java/lang/StackTraceElement"
	obj.Klass = &klassName

	// add array to the object we're returning
	fieldToAdd := new(object.Field)
	fieldToAdd.Ftype = "[Ljava/lang/StackTraceElement;"
	fieldToAdd.Fvalue = stackListing

	// add the field to the field table for this object
	obj.FieldTable["stackTrace"] = fieldToAdd

	return obj
}
