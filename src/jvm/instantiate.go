/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"errors"
	"fmt"
	"jacobin/classloader"
	"jacobin/log"
	"jacobin/object"
	"jacobin/shutdown"
	"os"
	"unsafe"
)

// instantiating a class is a two-part process:
//  1. the class needs to be loaded, so that its details and its methods are knowable
//  2. the class fields (if static) and instance fields (if non-static) are allocated. Details
//     for this second step appear in front of the initializeFields() method.
func instantiateClass(classname string) (*object.Object, error) {

	// Try to load class by name
	err := classloader.LoadClassFromNameOnly(classname)
	if err != nil {
		msg := "instantiateClass: Failed to load class " + classname
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		shutdown.Exit(shutdown.APP_EXCEPTION)
	}
	// Success in loaded by name
	_ = log.Log("instantiateClass: Success in LoadClassFromNameOnly("+classname+")", log.TRACE_INST)

	// at this point the class has been loaded into the method area (MethArea). Wait for it to be ready.
	err = classloader.WaitForClassStatus(classname)
	if err != nil {
		msg := fmt.Sprintf("instantiateClass: %s", err.Error())
		_ = log.Log(msg, log.SEVERE)
		return nil, errors.New(msg)
	}

	// At this point, classname is ready
	k := classloader.MethAreaFetch(classname)
	obj := object.Object{
		Klass: &classname,
	}

	// the object's mark field contains the lower 32-bits of the object's
	// address, which serves as the hash code for the object
	uintp := uintptr(unsafe.Pointer(&obj))
	obj.Mark.Hash = uint32(uintp)

	if len(k.Data.Fields) > 0 {
		for i := 0; i < len(k.Data.Fields); i++ {
			f := k.Data.Fields[i]
			desc := k.Data.CP.Utf8Refs[f.Desc]

			fieldToAdd := new(object.Field)
			fieldToAdd.Ftype = desc
			switch string(fieldToAdd.Ftype[0]) {
			case "L", "[": // it's a reference
				fieldToAdd.Fvalue = nil
			case "B", "C", "I", "J", "S", "Z":
				fieldToAdd.Fvalue = 0
			case "D", "F":
				fieldToAdd.Fvalue = 0.0
			default:
				_ = log.Log("error creating field in: "+classname+
					" Invalid type: "+fieldToAdd.Ftype, log.SEVERE)
				return nil, classloader.CFE("invalid field type")
			}

			if f.IsStatic {
				// in the instantiated class, add an 'X' before the
				// type, which notifies future users that the field
				// is static and should be fetched from the Statics
				// table.
				presentType := fieldToAdd.Ftype
				fieldToAdd.Ftype = "X" + presentType

				// static fields can have ConstantValue attributes,
				// which specify their initial value.
				if len(f.Attributes) > 0 {
					for j := 0; j < len(f.Attributes); j++ {
						attr := k.Data.CP.Utf8Refs[int(f.Attributes[j].AttrName)]
						if attr == "ConstantValue" {
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
								fieldToAdd.Fvalue = classloader.NewStringFromGoString(str)
							default:
								errMsg := fmt.Sprintf(
									"Unexpected ConstantValue type in instantiate: %d", valueType)
								_ = log.Log(errMsg, log.SEVERE)
								return nil, errors.New(errMsg)
							} // end of ConstantValue type switch
						} // end of ConstantValue attribute processing
					} // end of processing attributes
				} // end of search through attributes
				s := classloader.Static{
					Type:  fieldToAdd.Ftype,
					Value: fieldToAdd.Fvalue,
				}
				// add the field to the Statics table
				fieldName := k.Data.CP.Utf8Refs[f.Name]
				fullFieldName := classname + "." + fieldName

				_, alreadyPresent := classloader.Statics[fullFieldName]
				if !alreadyPresent { // add only if field has not been pre-loaded
					_ = classloader.AddStatic(fullFieldName, s)
				}
			}
			obj.Fields = append(obj.Fields, *fieldToAdd)
		} // end of processing fields
	}
	return &obj, nil
}

// the only fields allocated during class instantiation are instance fields--
// method-local fields are created on the stack during method execution.
// The allocated fields are in a structure that is defined in object.go
func initializeField(f classloader.Field, cp *classloader.CPool, cn string, obj *object.Object) {
	name := cp.Utf8Refs[int(f.Name)]
	desc := cp.Utf8Refs[int(f.Desc)]
	var attr = ""
	// var value int64
	if len(f.Attributes) > 0 {
		for i := 0; i < len(f.Attributes); i++ {
			attr = cp.Utf8Refs[int(f.Attributes[i].AttrName)]
			if attr == "ConstantValue" {
				// valueIndex := int(f.Attributes[i].AttrContent[0])*256 +
				//     int(f.Attributes[i].AttrContent[1])
				// // valueType := cp.CpIndex[valueIndex].Type
				// // valueSlot := cp.CpIndex[valueIndex].Slot
			}
			// } else {
			// 	value = 0
			// }

		}
	}
	// append field to the object.fields slice TODO: check that this is a good solution.
	// obj.fields = append(obj.fields, Field{
	// 	metadata: classloader.Field{},
	// 	value:    value,
	// })
	_, _ = fmt.Fprintf(os.Stdout, "Class: %s, Field to initialize: %s, type: %s\n", cn, name, desc)
	if attr != "" {
		_, _ = fmt.Fprintf(os.Stdout, "Attribute name: %s\n", attr)
	}
}
