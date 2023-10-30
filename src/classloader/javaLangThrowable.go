/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package classloader

import (
	"container/list"
	"fmt"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/object"
	"jacobin/thread"
)

func Load_Lang_Throwable() map[string]GMeth {

	MethodSignatures["java/lang/Throwable.fillInStackTrace()Ljava/lang/Throwable;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fillInStackTrace,
		}
	return MethodSignatures
}

func fillInStackTrace(params []interface{}) interface{} {
	glob := globals.GetGlobalRef()
	if glob.JVMframeStack == nil { // if we haven't captured the JVM stack before now, we're hosed.
		_ = log.Log("No stack data available for this error. Incomplete data will be shown.", log.SEVERE)
		return nil
	}

	thisThread := params[0].(*thread.ExecThread)
	thisFrameStack := thisThread.Stack
	stackListing := GetStackTraces(thisFrameStack)
	listing := stackListing.FieldTable["stackTrace"].Fvalue.([]*object.Object)
	fmt.Printf("Stack trace contains %d elements", len(listing))

	// thisFrame := thisFrameStack.Front().Next()

	// CURR: next steps
	// This might require that we add the logic to the class parse showing the Java code source line number.
	// JACOBIN-224 refers to this.
	return nil
}

// gets the full JVM stack trace using java.lang.StackTraceElement slice to hold the data
// in case of error, nil is returned
func GetStackTraces(fs *list.List) *object.Object {
	var stackListing []*object.Object
	// stackListing =

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
		k := MethAreaFetch(classname)
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
		fld := object.Field{}
		fld.Fvalue = frame.ClName
		stackTrace.FieldTable["declaringClass"] = &fld

		fld = object.Field{}
		fld.Fvalue = frame.MethName
		stackTrace.FieldTable["methodName"] = &fld

		methClass := MethAreaFetch(frame.ClName)
		if methClass == nil {
			return nil
		}

		fld = object.Field{}
		fld.Fvalue = methClass.Loader
		stackTrace.FieldTable["classLoaderName"] = &fld

		fld = object.Field{}
		fld.Fvalue = methClass.Data.SourceFile
		stackTrace.FieldTable["fileName"] = &fld

		fld = object.Field{}
		fld.Fvalue = methClass.Data.Module
		stackTrace.FieldTable["moduleName"] = &fld

		stackListing = append(stackListing, stackTrace)
	}

	// now that we have our data items loaded into the StackTraceElement
	// put the elments into an array, which is converted into an object
	obj := object.MakeEmptyObject()
	klassName := "java/lang/StackTraceElement"
	obj.Klass = &klassName

	// add array to the object we're returning
	fieldToAdd := new(object.Field)
	fieldToAdd.Ftype = "[Ljava/lang/StackTraceElement;"
	fieldToAdd.Fvalue = stackListing

	// add the field to the field table for this object
	obj.FieldTable["stackTrace"] = fieldToAdd

	// 	methName := fmt.Sprintf("%s.%s", frame.ClName, frame.MethName)
	// 	// stackTrace.FieldTable[]
	// 	entry := fmt.Sprintf("Method: %-40s PC: %03d", methName, frame.PC)
	// 	stackListing = append(stackListing, object.NewStringFromGoString(entry))
	// 	return stackListing
	// }
	// // return *stackListing

	return obj
}
