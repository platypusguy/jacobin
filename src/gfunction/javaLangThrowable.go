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
	"jacobin/log"
	"jacobin/object"
	"jacobin/shutdown"
	"jacobin/statics"
	"jacobin/types"
)

func Load_Lang_Throwable() map[string]GMeth {

	MethodSignatures["java/lang/Throwable.fillInStackTrace()Ljava/lang/Throwable;"] =
		GMeth{
			ParamSlots:   0,
			GFunction:    fillInStackTrace,
			NeedsContext: true,
		}

	MethodSignatures["java/lang/Throwable.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  throwableClinit,
		}

	MethodSignatures["java/lang/Throwable.getOurStackTrace:()[Ljava/lang/StackTraceElement;"] =
		GMeth{
			ParamSlots:   0,
			GFunction:    getOurStackTrace,
			NeedsContext: true,
		}

	return MethodSignatures
}

// This method duplicates the following bytecode, with these exceptions:
//  1. we don't check for assertion status, which is determined already at start-up
//  2. for the nonce, Throwable.SUPPRESSED_SENTINEL is set to nil. It's unlikely we'll
//     ever need it, but if we do, we'll implement it then.
//
// So, essentially, we're just initializing several static fields (as expected in clinit())
//
// 0:  ldc           #8                  // class java/lang/Throwable
// 2:  invokevirtual #302                // Method java/lang/Class.desiredAssertionStatus:()Z
// 5:  ifne          12
// 8:  iconst_1
// 9:  goto          13
// 12: iconst_0
// 13: putstatic     #153                // Field $assertionsDisabled:Z
// 16: iconst_0
// 17: anewarray     #173                // class java/lang/StackTraceElement
// 20: putstatic     #13                 // Field UNASSIGNED_STACK:[Ljava/lang/StackTraceElement;
// 23: invokestatic  #305                // Method java/util/Collections.emptyList:()Ljava/util/List;
// 26: putstatic     #20                 // Field SUPPRESSED_SENTINEL:Ljava/util/List;
// 29: iconst_0
// 30: anewarray     #8                  // class java/lang/Throwable
// 33: putstatic     #293                // Field EMPTY_THROWABLE_ARRAY:[Ljava/lang/Throwable;
// java of previous: private static final Throwable[] EMPTY_THROWABLE_ARRAY = new Throwable[0];
// 36: return
func throwableClinit(params []interface{}) interface{} {
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
func fillInStackTrace(params []interface{}) interface{} {

	if len(params) != 2 {
		_ = log.Log(fmt.Sprintf("fillInsStackTrace() expected two parameterss, got: %d",
			len(params)), log.SEVERE)
		shutdown.Exit(shutdown.JVM_EXCEPTION)
	}
	frameStack := params[0].(*list.List)
	objRef := params[1].(*object.Object)

	// we're adding the frame stack reference as a field to Throwable. This is
	// unique to Jacobin. (HotSpot accesses the frame stack through a completely
	// different mechanism that has no direct counterpart in Jacobin). This
	// step allows any Throwable to access the JVM frame stack, which is
	// necessary in stackTraceElement methods.
	jacobinSpecificField := object.Field{
		Ftype:  types.Ref,
		Fvalue: frameStack,
	}
	objRef.FieldTable["frameStackRef"] = &jacobinSpecificField
	fmt.Printf("Throwable object contains: %v\n", objRef.FieldTable)

	args := []interface{}{frameStack}
	return getOurStackTrace(args) // <<<<<<<<<<< we get here currently <<<<<<<<<<
}

func getOurStackTrace(params []interface{}) interface{} {
	return GetStackTraces(params[0].(*list.List))
}

// GetStackTraces gets the full JVM stack trace using java.lang.StackTraceElement
// slice to hold the data. In case of error, nil is returned.
func GetStackTraces(fs *list.List) *object.Object {
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
