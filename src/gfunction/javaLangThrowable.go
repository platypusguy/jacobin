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

	return MethodSignatures
}

// clinit needs to be implmented:
// Code:
// stack=1, locals=0, args_size=0
// 0: ldc           #8                  // class java/lang/Throwable
// 2: invokevirtual #302                // Method java/lang/Class.desiredAssertionStatus:()Z
// 5: ifne          12
// 8: iconst_1
// 9: goto          13
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
	if len(params) != 1 {
		_ = log.Log(fmt.Sprintf("Throwable.clinit() expected one params, got: %d", len(params)), log.SEVERE)
		shutdown.Exit(shutdown.JVM_EXCEPTION)
	}
	// get the throwable object we're doing the clinit on
	// throwable := params[0]
	return nil
}

// this function is called by Throwable.<init>()
func fillInStackTrace(params []interface{}) interface{} {
	if len(params) != 2 {
		_ = log.Log(fmt.Sprintf("fillInsStackTrace() expected two params, got: %d", len(params)), log.SEVERE)
		shutdown.Exit(shutdown.JVM_EXCEPTION)
	}
	frameStack := params[0].(*list.List)
	objRef := params[1].(*object.Object)
	fmt.Printf("Throwable object contains: %v", objRef.FieldTable)

	// thisFrame := frameStack.Front().Next()
	for thisFrame := frameStack.Front().Next(); thisFrame != nil; thisFrame = thisFrame.Next() {
		global := *globals.GetGlobalRef()
		ste, err := global.FuncInstantiateClass("java/lang/StackTraceElement", nil)
		if err != nil {
			_ = log.Log("Throwable.fillInStackTrace: error creating 'java/lang/StackTraceElement", log.SEVERE)
			return ste.(*object.Object)
		}
		fmt.Println(thisFrame.Value)
	}

	// This might require that we add the logic to the class parse showing the Java code source line number.
	// JACOBIN-224 refers to this.
	return objRef
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
