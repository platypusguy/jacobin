/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-4 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) Consult jacobin.org.
 */

package jvm

import (
	"container/list"
	"errors"
	"fmt"
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/exceptions"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/shutdown"
	"jacobin/statics"
	"jacobin/stringPool"
	"jacobin/trace"
	"jacobin/types"
	"strings"
	"unsafe"
)

// instantiating an object is a two-part process (except for arrays, which are handled
// by special bytecodes):
//
//  1. the class needs to be loaded, so that its details and its methods are knowable
//
//  2. the class fields (if static) and instance fields (if non-static) are allocated.
//     Details for this second step appear in the loop that drives createField().
//
//     NOTE: The "any" type returned is always *object.Object.
//     This is being done to avoid a golang circularity error when the caller
//     is one of the native 'G' functions.
func InstantiateClass(classname string, frameStack *list.List) (any, error) {

	if !strings.HasPrefix(classname, "[") { // do this only for classes, not arrays
		err := loadThisClass(classname)
		if err != nil { // error message will have been displayed
			return nil, err
		}
	}

	// strings are handled separately
	if classname == types.StringClassName {
		return object.NewStringObject(), nil
	}

	// At this point, classname is ready
	k := classloader.MethAreaFetch(classname)
	obj := object.Object{
		KlassName: stringPool.GetStringIndex(&classname),
	}

	if k == nil {
		errMsg := "InstantiateClass: Class is nil after loading, class: " + classname
		trace.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	if k.Data == nil {
		errMsg := "InstantiateClass: class.Data is nil, class: " + classname
		trace.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	// go up the chain of superclasses until we hit java/lang/Object
	superclasses := []string{}
	superclassNamePtr := stringPool.GetStringPointer(k.Data.SuperclassIndex)
	for {
		// if the present class is Object, it has no superclass. If the present
		// class's superclass is Object, we've reached the top of the superclass
		// hierarchy. Otherwise, keep looping up the superclasses.
		if classname == types.ObjectClassName || *superclassNamePtr == types.ObjectClassName {
			break
		}

		err := loadThisClass(*superclassNamePtr) // load the superclass
		if err != nil {                          // error message will have been displayed
			return nil, err
		} else {
			superclasses = append(superclasses, *superclassNamePtr)
		}

		loadedSuperclass := classloader.MethAreaFetch(*superclassNamePtr)
		// now loop to see whether this superclass has a superclass
		superclassNamePtr = stringPool.GetStringPointer(loadedSuperclass.Data.SuperclassIndex)
	}

	// the object's mark field contains the lower 32-bits of the object's
	// address, which serves as the hash code for the object
	uintp := uintptr(unsafe.Pointer(&obj))
	obj.Mark.Hash = uint32(uintp)

	// handle the fields. If the object has no superclass other than Object,
	// the fields are in an array in the order they're declared in the CP.
	// If the object has a non-Object superclass, then the superclasses' fields
	// and the present object's field are stored in a map--indexed by the
	// field name. Eventually, we might coalesce on a single approach for
	// both kinds of objects.
	if len(superclasses) == 0 && len(k.Data.Fields) == 0 {
		goto runInitializer // check to see if any static initializers
	}

	// initialize the map of this object's fields
	obj.FieldTable = make(map[string]object.Field)

	if len(superclasses) == 0 {
		for i := 0; i < len(k.Data.Fields); i++ {
			fld := k.Data.Fields[i]
			fldName := k.Data.CP.Utf8Refs[fld.Name]

			fieldToAdd, err := createField(fld, k, classname)
			if err != nil {
				return nil, err
			}
			obj.FieldTable[fldName] = *fieldToAdd

			// prepare the static fields, by inserting them w/ default values in Statics table
			// See (https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-5.html#jvms-5.4.2)
			if fld.IsStatic {
				var fldValue any
				fldType := []byte(k.Data.CP.Utf8Refs[fld.Desc])
				switch fldType[0] {
				case 'B', 'C', 'S', 'I', 'J', 'Z':
					fldValue = int64(0)
				case 'F', 'D':
					fldValue = float64(0.00)
				case 'L', '[':
					fldValue = object.Null
				}
				statics.AddStatic(classname+"."+fldName,
					statics.Static{Type: string(fldType[0]), Value: fldValue}) // CURR
			}
		} // loop through the fields if any
		goto runInitializer
	} // end of handling fields for objects w/ no superclasses

	// in the case of superclasses, we start at the topmost superclass
	// and work our way down to the present class, adding fields to FieldTable.
	// so we add the present class into position[0] and then loop through
	// the slice of class names
	superclasses = append([]string{classname}, superclasses...)
	for j := len(superclasses) - 1; j >= 0; j-- {
		superclassName := superclasses[j]
		c := classloader.MethAreaFetch(superclassName)
		if c == nil {
			errMsg := fmt.Sprintf("InstantiateClass: MethAreaFetch(superclass: %s) failed", superclassName)
			trace.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		for i := 0; i < len(c.Data.Fields); i++ {
			f := c.Data.Fields[i]
			name := c.Data.CP.Utf8Refs[f.Name]

			fieldToAdd, err := createField(f, c, classname)
			if err != nil {
				return nil, err
			}

			// add the field to the field table for this object
			obj.FieldTable[name] = *fieldToAdd
		} // end of handling fields for one  class or superclass
	} // end of handling fields for classes with superclasses other than Object

runInitializer:
	// check code validity in methods
	for _, m := range k.Data.MethodTable {
		code := m.CodeAttr.Code
		err := classloader.CheckCodeValidity(code, &k.Data.CP)
		if err != nil {
			errMsg := fmt.Sprintf("InstantiateClass: CheckCodeValidity failed with %s.%s", classname, m.Name)
			status := exceptions.ThrowEx(excNames.ClassFormatError, errMsg, nil)
			if status != exceptions.Caught {
				return nil, errors.New(errMsg) // applies only if in test
			}
		}
	}

	// run intialization blocks
	_, ok := k.Data.MethodTable["<clinit>()V"]
	if ok && k.Data.ClInit == types.ClInitNotRun {
		err := runInitializationBlock(k, superclasses, frameStack)
		if err != nil {
			errMsg := fmt.Sprintf("InstantiateClass: runInitializationBlock failed with %s.<clinit>()V", classname)
			trace.Error(errMsg)
			return nil, err
		}
	}

	return &obj, nil
}

// creates a field for insertion into the object representation
func createField(f classloader.Field, k *classloader.Klass, classname string) (*object.Field, error) {
	desc := k.Data.CP.Utf8Refs[f.Desc]

	fieldToAdd := new(object.Field)
	fieldToAdd.Ftype = desc
	switch string(fieldToAdd.Ftype[0]) {
	case types.Ref, types.Array: // it's a reference
		fieldToAdd.Fvalue = nil
	case types.Byte, types.Char, types.Int, types.Long, types.Short, types.Bool:
		fieldToAdd.Fvalue = int64(0)
	case types.Double, types.Float:
		fieldToAdd.Fvalue = 0.0
	default:
		errMsg := fmt.Sprintf("createField: error creating field in: %s,  Invalid type: %s",
			classname, fieldToAdd.Ftype)
		trace.Error(errMsg)
		return nil, classloader.CFE(errMsg)
	}

	presentType := fieldToAdd.Ftype
	if f.IsStatic {
		// in the instantiated class, add a types.Static before the
		// type, which notifies future users that the field
		// is static and should be fetched from the Statics
		// table.
		fieldToAdd.Ftype = types.Static + presentType
	}

	// static fields can have ConstantValue attributes,
	// which specify their initial value.
	if len(f.Attributes) > 0 {
		for j := 0; j < len(f.Attributes); j++ {
			attr := k.Data.CP.Utf8Refs[int(f.Attributes[j].AttrName)]
			if attr == "ConstantValue" && f.IsStatic { // only statics can have ConstantValue attribute
				valueIndex := int(f.Attributes[j].AttrContent[0])*256 +
					int(f.Attributes[j].AttrContent[1])
				valueType := k.Data.CP.CpIndex[valueIndex].Type
				valueSlot := k.Data.CP.CpIndex[valueIndex].Slot
				switch valueType {
				case classloader.IntConst:
					fieldToAdd.Fvalue = int64(k.Data.CP.IntConsts[valueSlot])
				case classloader.LongConst:
					fieldToAdd.Fvalue = k.Data.CP.LongConsts[valueSlot]
				case classloader.FloatConst:
					fieldToAdd.Fvalue = float64(k.Data.CP.Floats[valueSlot])
				case classloader.DoubleConst:
					fieldToAdd.Fvalue = k.Data.CP.Doubles[valueSlot]
				case classloader.StringConst:
					str := k.Data.CP.Utf8Refs[valueSlot]
					fieldToAdd.Fvalue = object.StringObjectFromGoString(str)
				default:
					errMsg := fmt.Sprintf(
						"createField: Unexpected ConstantValue type in instantiate: %d", valueType)
					trace.Error(errMsg)
					return nil, errors.New(errMsg)
				} // end of ConstantValue type switch
			} // end of ConstantValue attribute processing
		} // end of processing attributes
	} // end of search through attributes

	if f.IsStatic {
		s := statics.Static{
			Type:  presentType, // we use the type without the 'X' prefix in the statics table.
			Value: fieldToAdd.Fvalue,
		}
		// add the field to the Statics table
		fieldName := k.Data.CP.Utf8Refs[f.Name]
		fullFieldName := classname + "." + fieldName

		_, alreadyPresent := statics.Statics[fullFieldName]
		if !alreadyPresent { // add only if field has not been pre-loaded
			_ = statics.AddStatic(fullFieldName, s)
		}
	}
	return fieldToAdd, nil
}

// Loads the class (if it's not already loaded) and makes sure it's accessible in the method area
func loadThisClass(className string) error {
	alreadyLoaded := classloader.MethAreaFetch(className)
	if alreadyLoaded != nil { // if the class is already loaded, skip the rest of this
		return nil
	}
	// Try to load class by name
	err := classloader.LoadClassFromNameOnly(className)
	if err != nil {
		shutdown.Exit(shutdown.APP_EXCEPTION)
		return errors.New(err.Error()) // needed for testing, which does not shutdown on failure
	}
	// Success in loaded by name
	if globals.TraceCloadi {
		trace.Trace("loadThisClass: Success in LoadClassFromNameOnly(" + className + ")")
	}

	// at this point the class has been loaded into the method area (MethArea). Wait for it to be ready.
	err = classloader.WaitForClassStatus(className)
	if err != nil {
		errMsg := fmt.Sprintf("loadThisClass: WaitForClassStatus(%s) failed, err: %v", className, err)
		trace.Error(errMsg)
		return errors.New(errMsg) // needed for testing, which does not shutdown on failure
	}
	return nil
}
