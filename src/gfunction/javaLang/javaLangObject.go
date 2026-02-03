/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaLang

import (
	"container/list"
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/frames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaUtil"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"strings"
	"sync/atomic"
	"unsafe"
)

// Implementation of some of the functions in Java/lang/Class.

func Load_Lang_Object() {

	// --- Already implemented ---
	ghelpers.MethodSignatures["java/lang/Object.<clinit>()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: objectClinitInit}

	ghelpers.MethodSignatures["java/lang/Object.<init>()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: objectClinitInit}

	ghelpers.MethodSignatures["java/lang/Object.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: objectEquals}

	ghelpers.MethodSignatures["java/lang/Object.finalize()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapDeprecated}

	ghelpers.MethodSignatures["java/lang/Object.getClass()Ljava/lang/Class;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ObjectGetClass} // TODO: finish implementing objectGetClass

	ghelpers.MethodSignatures["java/lang/Object.getResourceAsStream(Ljava/lang/String;)Ljava/io/InputStream;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}

	ghelpers.MethodSignatures["java/lang/Object.hashCode()I"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: objectHashCode}

	ghelpers.MethodSignatures["java/lang/Object.notify()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: objectNotify, NeedsContext: true}

	ghelpers.MethodSignatures["java/lang/Object.notifyAll()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: objectNotifyAll, NeedsContext: true}

	ghelpers.MethodSignatures["java/lang/Object.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: objectToString}

	ghelpers.MethodSignatures["java/lang/Object.wait()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: objectWait, NeedsContext: true}

	ghelpers.MethodSignatures["java/lang/Object.wait(J)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: objectWait, NeedsContext: true}

	ghelpers.MethodSignatures["java/lang/Object.wait(JI)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: objectWait, NeedsContext: true}

	// --- All other Object methods as traps (alphabetical) ---
	addTrap := func(signature string, slots int) {
		if _, exists := ghelpers.MethodSignatures[signature]; !exists {
			ghelpers.MethodSignatures[signature] =
				ghelpers.GMeth{ParamSlots: slots, GFunction: ghelpers.TrapFunction}
		}
	}

	// Alphabetically sorted
	trapMethods := map[string]int{
		"java/lang/Object.clone()Ljava/lang/Object;": 0, // protected
	}

	for m, slots := range trapMethods {
		addTrap(m, slots)
	}
}

// === the internal representation of a java.lang.Class() instance ===
// this is not a faithful reproduction of the OpenJDK version, but rather
// the one we use in Jacobin
type javaLangClass struct {
	accessFlags    classloader.AccessFlags
	name           string
	superClassName string
	interfaceNames []string
	constantPool   classloader.CPool
	fields         []classloader.Field
	methods        map[string]*classloader.Method
	loader         string
	superClass     string
	interfaces     []uint16 // indices into UTF8Refs
	// instanceSlotCount uint
	// staticSlotCount   uint
	// staticVars        Slots
}

func objectClinitInit(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	if obj == nil {
		errMsg := fmt.Sprintf("objectClinitInit: Invalid or missing object: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj = object.MakeEmptyObjectWithClassName(&types.ObjectClassName)

	return nil
}

// objectGetClass implements "java/lang/Object.getClass()Ljava/lang/Class;"
// It returns a pointer to the skeletal Class object for this object,
// which is located in the global JLC table (JLC = java/lang/Class).
func ObjectGetClass(params []interface{}) interface{} {
	objPtr := params[0].(*object.Object)
	if objPtr == nil || objPtr.KlassName == types.InvalidStringIndex {
		errMsg := fmt.Sprintf("java/lang/Object.getClass: Invalid object: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	jlc := globals.JLCmap[*stringPool.GetStringPointer(objPtr.KlassName)]
	if jlc == nil {
		errMsg := fmt.Sprintf("java/lang/Object.getClass: Class %s not loaded",
			object.GoStringFromStringPoolIndex(objPtr.KlassName))
		return ghelpers.GetGErrBlk(excNames.ClassNotLoadedException, errMsg)
	} else {
		return jlc
	}
	/*
		name := object.GoStringFromStringPoolIndex(objPtr.KlassName)

		if strings.HasPrefix(name, types.Array) { // arrays are handled differently
			arrClass := arrayGetClass(objPtr, name)
			return arrClass
		}

		// get a pointer to the class contents from the method area
		content := classloader.MethAreaFetch(name)
		if content == nil {
			errMsg := fmt.Sprintf("java/lang/Object.getClass: Class %s not loaded", name)
			return ghelpers.GetGErrBlk(excNames.ClassNotLoadedException, errMsg)
		}

		// syntactic sugar
		obj := *content

		// if we've previously created the Class object, return it
		if obj.Data.ClassObject != nil {
			return obj.Data.ClassObject
		}

		// create the empty java.lang.Class structure
		jlc := object.MakeEmptyObject()

		// points to the internal metaspace representation of the class (in methArea)
		// HotSpot uses a hidden field named _klass for this. So do we.
		jlc.FieldTable = make(map[string]object.Field)
		jlc.FieldTable["_klass"] = object.Field{
			Ftype:  types.RawGoPointer,
			Fvalue: content,
		}

		className := util.ConvertInternalClassNameToUserFormat(name) // FQN uses . not /
		jlc.FieldTable["name"] = object.Field{
			Ftype:  types.Ref,
			Fvalue: object.StringObjectFromGoString(className),
		}

		jlc.FieldTable["classLoader"] = object.Field{
			Ftype:  types.Ref,
			Fvalue: object.StringObjectFromGoString(obj.Loader),
		}

		objData := *obj.Data
		jlc.FieldTable["constantPool"] = object.Field{
			Ftype:  types.Struct,
			Fvalue: objData.CP,
		}

		jlc.FieldTable["superClass"] = object.Field{
			Ftype:  types.GolangString,
			Fvalue: object.GoStringFromStringPoolIndex(objData.SuperclassIndex),
		}

		jlc.FieldTable["fields"] = object.Field{
			Ftype:  types.Struct,
			Fvalue: objData.Fields,
		}

		jlc.FieldTable["interfaces"] = object.Field{
			Ftype:  types.Struct,
			Fvalue: objData.Interfaces,
		}

		jlc.FieldTable["methods"] = object.Field{
			Ftype:  types.Struct,
			Fvalue: objData.MethodTable,
		}

		jlc.FieldTable["modifiers"] = object.Field{
			Ftype:  types.Struct,
			Fvalue: objData.Access,
		}

		return jlc

	*/
}

// "java/lang/Object.toString()Ljava/lang/String;"
func objectToString(params []interface{}) interface{} {
	// params[0]: input Object

	switch params[0].(type) {
	case *object.Object:
		inObj := params[0].(*object.Object)
		classNameSuffix := object.GetClassNameSuffix(inObj, false)
		if classNameSuffix == "LinkedList" {
			return javaUtil.LinkedlistToString(params)
		}
		return object.StringifyAnythingJava(inObj)
	}

	errMsg := fmt.Sprintf("objectToString: Unsupported parameter type: %T", params[0])
	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
}

// "java/lang/Object.hashCode()I"
func objectHashCode(params []interface{}) interface{} {
	// params[0]: input Object
	switch params[0].(type) {
	case *object.Object:
		ptr := uintptr(unsafe.Pointer(params[0].(*object.Object)))
		hashCode := int64(ptr ^ (ptr >> 32))
		return hashCode
	}

	errMsg := fmt.Sprintf("objectHashCode: Unsupported parameter type: %T", params[0])
	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
}

func objectEquals(params []interface{}) interface{} {
	this, ok := params[0].(*object.Object)
	if !ok {
		return types.JavaBoolFalse
	}
	that, ok := params[1].(*object.Object)
	if !ok {
		return types.JavaBoolFalse
	}

	// If they are the same object, even if null, return true.
	if this == that {
		return types.JavaBoolTrue
	}

	// Not the same object.
	return types.JavaBoolFalse
}

func objectWait(params []interface{}) interface{} {

	// Get frame stack.
	fs, ok := params[0].(*list.List)
	if !ok {
		errMsg := fmt.Sprintf("objectWait: params[0] must be the frame stack, saw: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// Get thread ID.
	frame := *fs.Front().Value.(*frames.Frame)
	thID := int32(frame.Thread)

	// Get the object of the synchronized method.
	obj, ok := params[1].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("objectWait: params[1] must be an Object, saw: %T", params[1])
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// Set millis = sleep time in milliseconds.
	// TODO: Sub-millisecond precision is not done well.
	millis := int64(0) // 0 means wait indefinitely
	if len(params) > 2 {
		millis = params[2].(int64)
		nanos := int64(0)
		if len(params) > 3 {
			nanos = params[3].(int64)
			if nanos > 0 {
				millis += 1 // not precise
			}
		}
	}

	err := obj.ObjectWait(thID, millis)
	if err != nil {
		monitor := obj.GetMonitor()
		errMsg := fmt.Sprintf("objectWait: thID=%d, wait-obj-class=%s, obj-monitor.Owner=%d\n%s",
			thID, object.GoStringFromStringPoolIndex(obj.KlassName), atomic.LoadInt32(&monitor.Owner), err.Error())

		// Check for wrong owner.
		if strings.Contains(err.Error(), "does not own lock") {
			return ghelpers.GetGErrBlk(excNames.IllegalMonitorStateException, errMsg)
		}
		// Interrupted?
		if strings.Contains(err.Error(), "interrupted") {
			return ghelpers.GetGErrBlk(excNames.InterruptedException, errMsg)
		}
		// Other errors.
		return ghelpers.GetGErrBlk(excNames.IllegalMonitorStateException, errMsg)
	}

	return nil
}

func objectNotify(params []interface{}) interface{} {
	// Get frame stack.
	fs, ok := params[0].(*list.List)
	if !ok {
		errMsg := fmt.Sprintf("objectNotify: params[0] must be the frame stack, saw: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// Get thread ID.
	frame := *fs.Front().Value.(*frames.Frame)
	thID := int32(frame.Thread)

	// Get the object of the synchronized method.
	obj, ok := params[1].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("objectNotify: params[1] must be an Object, saw: %T", params[1])
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	err := obj.ObjectNotify(thID)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IllegalMonitorStateException, err.Error())
	}

	return nil
}

func objectNotifyAll(params []interface{}) interface{} {
	// Get frame stack.
	fs, ok := params[0].(*list.List)
	if !ok {
		errMsg := fmt.Sprintf("objectNotifyAll: params[0] must be the frame stack, saw: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// Get thread ID.
	frame := *fs.Front().Value.(*frames.Frame)
	thID := int32(frame.Thread)

	// Get the object of the synchronized method.
	obj, ok := params[1].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("objectNotifyAll: params[1] must be an Object, saw: %T", params[1])
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	err := obj.ObjectNotifyAll(thID)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IllegalMonitorStateException, err.Error())
	}

	return nil
}

// arrayGetClass creates a Class object for array types
// Arrays have special handling because they're not loaded from .class files
// Per JVM spec, all arrays have Object as their superclass
func arrayGetClass(objPtr *object.Object, arrayName string) *object.Object {
	jlc := object.MakeEmptyObject()
	jlc.FieldTable = make(map[string]object.Field)

	// Set the name field to the array type descriptor (e.g., "[Ljava/lang/String;" or "[I")
	jlc.FieldTable["name"] = object.Field{
		Ftype:  types.GolangString,
		Fvalue: arrayName,
	}

	// Determine the component type (the type of elements in the array)
	// For example: "[Ljava/lang/String;" -> "java/lang/String"
	//              "[I" -> "int"
	//              "[[I" -> "[I"
	componentType := ""
	if len(arrayName) > 1 {
		componentType = arrayName[1:] // Remove the leading '['

		// Convert internal format to readable format for object arrays
		// e.g., "Ljava/lang/String;" -> "java/lang/String"
		if strings.HasPrefix(componentType, "L") && strings.HasSuffix(componentType, ";") {
			componentType = componentType[1 : len(componentType)-1]
		}

		// Handle primitive types
		switch componentType {
		case "Z":
			componentType = "boolean"
		case "B":
			componentType = "byte"
		case "C":
			componentType = "char"
		case "D":
			componentType = "double"
		case "F":
			componentType = "float"
		case "I":
			componentType = "int"
		case "J":
			componentType = "long"
		case "S":
			componentType = "short"
		}
	}

	jlc.FieldTable["componentType"] = object.Field{
		Ftype:  types.GolangString,
		Fvalue: componentType,
	}

	// Arrays always have Object as their superclass
	jlc.FieldTable["superClass"] = object.Field{
		Ftype:  types.GolangString,
		Fvalue: "java/lang/Object",
	}

	// Arrays don't have fields (other than length, which is implicit)
	jlc.FieldTable["fields"] = object.Field{
		Ftype:  types.Struct,
		Fvalue: []classloader.Field{},
	}

	// Arrays don't have methods
	jlc.FieldTable["methods"] = object.Field{
		Ftype:  types.Struct,
		Fvalue: map[string]*classloader.Method{},
	}

	// Arrays don't have interfaces
	jlc.FieldTable["interfaces"] = object.Field{
		Ftype:  types.Struct,
		Fvalue: []uint16{},
	}

	// Set modifiers - arrays are always public and final
	accessFlags := classloader.AccessFlags{
		ClassIsPublic: true,
		ClassIsFinal:  true,
	}
	jlc.FieldTable["modifiers"] = object.Field{
		Ftype:  types.Struct,
		Fvalue: accessFlags,
	}

	// Arrays use the bootstrap classloader
	jlc.FieldTable["classLoader"] = object.Field{
		Ftype:  types.GolangString,
		Fvalue: "bootstrap",
	}

	return jlc
}
