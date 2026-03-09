/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
)

func Load_Lang_Invoke_MethodType() {
	ghelpers.MethodSignatures["java/lang/invoke/MethodType.fromMethodDescriptorString(Ljava/lang/String;Ljava/lang/ClassLoader;)Ljava/lang/invoke/MethodType;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  MethodTypeFromMethodDescriptorString,
		}
}

const methodTypeClassName = "java/lang/invoke/MethodType"

// "java/lang/invoke/MethodType.fromMethodDescriptorString(Ljava/lang/String;Ljava/lang/ClassLoader;)Ljava/lang/invoke/MethodType;"
// Create a method type, suitable for use with method handles, from a descriptor string
//  1. Parses the descriptor: It splits the descriptor into parameter types and return type.
//  2. Resolves types: It resolves primitive types (e.g., "I" -> int.class) and object types (e.g.,
//     "Ljava/lang/String;" -> String.class) to their corresponding java.lang.Class objects.
//  3. Constructs the object: It creates a new java.lang.invoke.MethodType object.
//  4. Populates fields: It sets the rtype (return type) and ptypes (parameter types array) fields
//     of the MethodType object, which matches the internal structure of the JDK class.
func MethodTypeFromMethodDescriptorString(params []interface{}) interface{} {

	descriptorObj := params[0].(*object.Object)
	// classLoaderObj := params[1].(*object.Object) // TODO: Might need later if we support custom class loaders

	descriptor := object.GoStringFromStringObject(descriptorObj)
	return descriptor
}

// 	// Parse the descriptor to get Class objects for return and parameter types
// 	returnType, paramTypes, err := parseDescriptorToClasses(descriptor)
// 	if err != nil {
// 		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, err.Error())
// 	}
//
// 	// Now, construct the java.lang.invoke.MethodType object
// 	mtObj := object.MakeEmptyObject()
// 	mtObj.KlassName = object.StringPoolIndexFromGoString(methodTypeClassName)
//
// 	// Create a Java array of Class objects for the parameters
// 	paramArray := object.Make1DimRefArray("java/lang/Class", int64(len(paramTypes)))
// 	rawPtypeArray := paramArray.FieldTable["value"].Fvalue.([]*classloader.Jlc)
// 	copy(rawPtypeArray, paramTypes)
//
// 	// Set the fields of the MethodType object.
// 	// Based on OpenJDK, the fields are 'rtype' and 'ptypes'.
// 	mtObj.FieldTable["rtype"] = object.Field{Ftype: "Ljava/lang/Class;", Fvalue: returnType}
// 	mtObj.FieldTable["ptypes"] = object.Field{Ftype: "[Ljava/lang/Class;", Fvalue: paramArray}
//
// 	return mtObj
// }
//
// // parseDescriptorToClasses parses a method descriptor string and resolves each type
// // to its corresponding java.lang.Class object. Returns the return type of the method and the parameter types
// // as pointers to java.lang.Class instances.
// func parseDescriptorToClasses(descriptor string) (returnType *classloader.Jlc, paramTypes []*classloader.Jlc, err error) {
// 	if len(descriptor) == 0 || descriptor[0] != '(' {
// 		return nil, nil, fmt.Errorf("invalid method descriptor: %s", descriptor)
// 	}
//
// 	// Find the end of the parameter list
// 	endParen := strings.IndexRune(descriptor, ')')
// 	if endParen == -1 {
// 		return nil, nil, fmt.Errorf("invalid method descriptor: missing ')' in %s", descriptor)
// 	}
//
// 	paramStr := descriptor[1:endParen]
// 	returnStr := descriptor[endParen+1:]
//
// 	// Parse parameter types
// 	paramTypes = make([]*classloader.Jlc, 0)
// 	for i := 0; i < len(paramStr); {
// 		typeStr, width := getNextTypeDescriptor(paramStr[i:])
// 		if width == 0 {
// 			return nil, nil, fmt.Errorf("malformed parameter descriptor in %s", descriptor)
// 		}
// 		classObj, err := resolveTypeDescriptor(typeStr)
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		paramTypes = append(paramTypes, classObj)
// 		i += width
// 	}
//
// 	// Parse return type
// 	typeStr, width := getNextTypeDescriptor(returnStr)
// 	if width == 0 || width != len(returnStr) {
// 		return nil, nil, fmt.Errorf("malformed return descriptor in %s", descriptor)
// 	}
// 	returnType, err = resolveTypeDescriptor(typeStr)
// 	if err != nil {
// 		return nil, nil, err
// 	}
//
// 	return returnType, paramTypes, nil
// }
//
// // getNextTypeDescriptor extracts the next full type descriptor from a string.
// // It returns the descriptor and the number of characters it consumed.
// func getNextTypeDescriptor(d string) (string, int) {
// 	if len(d) == 0 {
// 		return "", 0
// 	}
// 	switch d[0] {
// 	case 'B', 'C', 'D', 'F', 'I', 'J', 'S', 'Z', 'V':
// 		return d[0:1], 1
// 	case 'L':
// 		end := strings.IndexRune(d, ';')
// 		if end == -1 {
// 			return "", 0 // Malformed
// 		}
// 		return d[0 : end+1], end + 1
// 	case '[':
// 		// It's an array, find the underlying type
// 		i := 1
// 		for i < len(d) && d[i] == '[' {
// 			i++
// 		}
// 		_, width := getNextTypeDescriptor(d[i:])
// 		if width == 0 {
// 			return "", 0 // Malformed
// 		}
// 		return d[0 : i+width], i + width
// 	default:
// 		return "", 0 // Invalid character
// 	}
// }
//
// // resolveTypeDescriptor converts a type descriptor string into a java.lang.Class object.
// func resolveTypeDescriptor(typeStr string) (*object.Object, error) {
// 	var className string
// 	var isPrimitive bool
//
// 	switch typeStr {
// 	case "B":
// 		className, isPrimitive = "java/lang/Byte", true
// 	case "C":
// 		className, isPrimitive = "java/lang/Character", true
// 	case "D":
// 		className, isPrimitive = "java/lang/Double", true
// 	case "F":
// 		className, isPrimitive = "java/lang/Float", true
// 	case "I":
// 		className, isPrimitive = "java/lang/Integer", true
// 	case "J":
// 		className, isPrimitive = "java/lang/Long", true
// 	case "S":
// 		className, isPrimitive = "java/lang/Short", true
// 	case "Z":
// 		className, isPrimitive = "java/lang/Boolean", true
// 	case "V":
// 		className, isPrimitive = "java/lang/Void", true
// 	default:
// 		// Object or Array type
// 		className = strings.ReplaceAll(typeStr, ".", "/")
// 		if strings.HasPrefix(className, "L") && strings.HasSuffix(className, ";") {
// 			className = className[1 : len(className)-1]
// 		}
// 	}
//
// 	// For primitive types, we need the special TYPE field (e.g., Integer.TYPE)
// 	if isPrimitive {
// 		staticField, ok := statics.QueryStatic(className, "TYPE")
// 		if !ok {
// 			// The wrapper class might not be initialized yet.
// 			if err := classloader.LoadClassFromNameOnly(className); err != nil {
// 				return nil, fmt.Errorf("could not load wrapper class %s: %v", className, err)
// 			}
// 			// Trigger static initialization which should populate the TYPE field.
// 			k := classloader.MethAreaFetch(className)
// 			if k.Data.ClInit == types.ClInitNotRun {
// 				globals.GetGlobalRef().FuncInvokeGFunction(k.Data.Name+".<clinit>()V", nil)
// 			}
// 			staticField, ok = statics.QueryStatic(className, "TYPE")
// 			if !ok {
// 				return nil, fmt.Errorf("primitive TYPE field not found for %s", className)
// 			}
// 		}
// 		return staticField.Value.(*object.Object), nil
// 	}
//
// 	// For non-primitive types, load the class and get its Class object.
// 	if err := classloader.LoadClassFromNameOnly(className); err != nil {
// 		return nil, fmt.Errorf("could not load class for descriptor %s: %v", className, err)
// 	}
//
//
// 	classloader.JlcMapLock.RLock()
// 	jlc, ok := classloader.JLCmap[className]
// 	classloader.JlcMapLock.RUnlock()
//
// 	if !ok {
// 		return nil, fmt.Errorf("Class object not found in JLCmap for %s", className)
// 	}
//
// 	// The JLC object itself is the java.lang.Class instance.
// 	return object.MakePrimitiveObjectFromJlcInstance(className), nil
// }
