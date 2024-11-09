/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"container/list"
	"errors"
	"fmt"
	"jacobin/object"
	"jacobin/shutdown"
	"jacobin/statics"
	"jacobin/trace"
	"jacobin/types"
)

func Load_Lang_Throwable() {

	MethodSignatures["java/lang/Throwable.fillInStackTrace()Ljava/lang/Throwable;"] =
		GMeth{
			ParamSlots:   0,
			GFunction:    FillInStackTrace,
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
func throwableClinit([]interface{}) interface{} {
	stackTraceElementClassName := "java/lang/StackTraceElement"
	emptyStackTraceElementArray := object.Make1DimRefArray(&stackTraceElementClassName, 0)
	_ = statics.AddStatic("Throwable.UNASSIGNED_STACK", statics.Static{
		Type:  "[Ljava/lang/StackTraceElement",
		Value: emptyStackTraceElementArray,
	})

	// for the time being, SUPPRESSED SENTINEL is set to nil.
	// We might later need to set it to an empty List.
	_ = statics.AddStatic("Throwable.SUPPRESSED_SENTINEL", statics.Static{
		Type: "Ljava/util/List", Value: nil})

	emptyThrowableClassName := "java/lang/Throwable"
	emptyThrowableArray := object.Make1DimRefArray(&emptyThrowableClassName, 0)
	_ = statics.AddStatic("Throwable.EMPTY_THROWABLE_ARRAY", statics.Static{
		Type:  "[Ljava/lang/Throwable",
		Value: emptyThrowableArray,
	})
	return nil
}

// This function is called by Throwable.<init>().
// In Throwable.java, it consists of one line:
//      return getOurStackTrace().clone(); // public, returns a StackTraceElement[]
// In turn, getOurStackTrace() calls
//      StackTraceElement.of(this, depth); // private, returns StackTraceElement[]
// In turn, this method calls
//      StackTraceElement.initStackTraceElements:([Ljava/lang/StackTraceElement;Ljava/lang/Throwable;)V
// which actually fills in the fields of the StackTraceElement (done as a native function)
//
// Despite this simple function chaining, there is value in reading the
// Javadoc for this function from Throwable.java (Copyright Oracle Corp.):
/*
 * Provides programmatic access to the stack trace information printed by
 * printStackTrace(). Returns an array of stack trace elements,
 * each representing one stack frame. The zeroth element of the array
 * (assuming the array's length is non-zero) represents the top of the
 * stack, which is the last method invocation in the sequence.  Typically,
 * this is the point at which this throwable was created and thrown.
 * The last element of the array (assuming the array's length is non-zero)
 * represents the bottom of the stack, which is the first method invocation
 * in the sequence.
 *
 * [...] Generally speaking, the array returned by this method will
 * contain one element for every frame that would be printed by
 * {@code printStackTrace}.  Writes to the returned array do not
 * affect future calls to this method.
 *
 * @return an array of stack trace elements representing the stack trace
 *         pertaining to this throwable.
 */
func FillInStackTrace(params []interface{}) interface{} {
	// get our parameters vetted and ready for use, then call getOurStackTrace()
	if len(params) != 2 {
		errMsg := fmt.Sprintf("FillInStackTrace: expected two parameters, got: %d", len(params))
		trace.Error(errMsg)
		shutdown.Exit(shutdown.JVM_EXCEPTION)
		return errors.New(errMsg) // needed only for testing b/c shutdown.Exit() doesn't exit in tests
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
	objRef.FieldTable["frameStackRef"] = jacobinSpecificField
	// fmt.Printf("Throwable object contains: %v\n", objRef.FieldTable)

	args := []interface{}{objRef}
	stackData := getOurStackTrace(args)
	throwable := *objRef

	stackTraceField := object.Field{
		Ftype:  types.Ref,
		Fvalue: stackData,
	}
	throwable.FieldTable["stackTrace"] = stackTraceField
	return &stackTraceField
}

// as described above, this function simply chains to GetStackTraces
func getOurStackTrace(params []interface{}) interface{} {
	args := []interface{}{params[0].(*object.Object)}
	return GetStackTraces(args)
}

// Calls stackTraceElement.of() which populates the entries in a
// slice of entries representing each frame in the JVM stack
func GetStackTraces(params []interface{}) *object.Object {
	throwable := params[0].(*object.Object)
	stack := throwable.FieldTable["frameStackRef"].Fvalue.(*list.List)
	depth := stack.Len()
	args := []interface{}{throwable, int64(depth)}
	retVal := of(args) // this is javaLangStackTraceElement.of()
	return retVal.(*object.Object)
}
