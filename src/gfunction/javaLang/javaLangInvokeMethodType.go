/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/statics"
	"strings"
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
	if object.IsNull(descriptorObj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "descriptorObj is null")
	}

	// classLoaderObj := params[1].(*object.Object) // TODO: Might need later if we support custom class loaders

	descriptor := object.GoStringFromStringObject(descriptorObj)

	// Parse the descriptor to get Class objects for return and parameter types
	returnType, paramTypes, err := parseDescriptorToClasses(descriptor)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, err.Error())
	}

	// Now, construct the java.lang.invoke.MethodType object
	mtObj := object.MakeEmptyObject()
	mtObj.KlassName = object.StringPoolIndexFromGoString(methodTypeClassName)

	// Create a Java array of Class objects for the parameters
	paramArray := object.Make1DimRefArray("java/lang/Class", int64(len(paramTypes)))
	rawPtypeArray := paramArray.FieldTable["value"].Fvalue.([]*object.Object)
	copy(rawPtypeArray, paramTypes)

	// Set the fields of the MethodType object.
	// Based on OpenJDK, the fields are 'rtype' and 'ptypes'.
	mtObj.FieldTable["rtype"] = object.Field{Ftype: "Ljava/lang/Class;", Fvalue: returnType}
	mtObj.FieldTable["ptypes"] = object.Field{Ftype: "[Ljava/lang/Class;", Fvalue: paramArray}

	return mtObj
}

// parseDescriptorToClasses parses a method descriptor string and resolves each type
// to its corresponding java.lang.Class object. Returns the return type of the method and the parameter types
// as pointers to java.lang.Class objects. The format for descriptors is specified here:
// https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-4.html#jvms-4.3.3
// Ex: Object m(int i, double d, Thread t) {...} has descriptor: (IDLjava/lang/Thread;)Ljava/lang/Object;
func parseDescriptorToClasses(descriptor string) (returnType *object.Object,
	paramTypes []*object.Object, err error) {
	if len(descriptor) == 0 || descriptor[0] != '(' {
		return nil, nil, fmt.Errorf("invalid method descriptor: %s", descriptor)
	}

	// Find the end of the parameter list
	endParen := strings.IndexRune(descriptor, ')')
	if endParen == -1 {
		return nil, nil, fmt.Errorf("invalid method descriptor: missing ')' in %s", descriptor)
	}

	paramStr := descriptor[1:endParen]   // everything inside the parentheses
	returnStr := descriptor[endParen+1:] // the bytes after the closing parenthesis

	// Parse parameter types
	paramTypes = make([]*object.Object, 0)
	for i := 0; i < len(paramStr); {
		typeStr, width := getNextTypeDescriptor(paramStr[i:])
		if width == 0 {
			return nil, nil, fmt.Errorf("invalid parameter descriptor in %s", descriptor)
		}
		classObj, err := resolveTypeDescriptor(typeStr)
		if err != nil {
			return nil, nil, err
		}
		paramTypes = append(paramTypes, classObj)
		i += width
	}

	// Parse return type
	typeStr, width := getNextTypeDescriptor(returnStr)
	if width == 0 || width != len(returnStr) {
		return nil, nil, fmt.Errorf("invalid return descriptor in %s", descriptor)
	}
	returnType, err = resolveTypeDescriptor(typeStr)
	if err != nil {
		return nil, nil, err
	}

	return returnType, paramTypes, nil
}

// getNextTypeDescriptor extracts the next full type descriptor from a string.
// It returns the descriptor and the number of characters it consumed.
func getNextTypeDescriptor(d string) (string, int) {
	if len(d) == 0 {
		return "", 0
	}
	switch d[0] {
	case 'B', 'C', 'D', 'F', 'I', 'J', 'S', 'Z', 'V':
		return d[0:1], 1
	case 'L':
		end := strings.IndexRune(d, ';')
		if end == -1 {
			return "", 0 // Malformed
		}
		return d[0 : end+1], end + 1
	case '[':
		// It's an array, find the underlying type
		i := 1
		for i < len(d) && d[i] == '[' {
			i++
		}
		_, width := getNextTypeDescriptor(d[i:])
		if width == 0 {
			return "", 0 // Malformed
		}
		return d[0 : i+width], i + width
	default:
		return "", 0 // Invalid character
	}
}

// resolveTypeDescriptor converts a type descriptor string into a java.lang.Class object.
// The descriptor formats are described here: https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-4.html#jvms-4.3.2
// and https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-4.html#jvms-4.3.3
// For primitives, the function maps the single-character descriptor (e.g. 'I') to the wrapper
// class name ('java/lang/Integer') and accesses the TYPE static variable of the wrapper class
// to get the java/lang/Class object for the primitive type.
// If the TYPE field is missing, it means InitializePrimitiveWrappers() hasn't run.
func resolveTypeDescriptor(typeStr string) (*object.Object, error) {
	var className string
	var primitiveName string
	var isPrimitive bool

	switch typeStr {
	case "B":
		className, primitiveName, isPrimitive = "java/lang/Byte", "byte", true
	case "C":
		className, primitiveName, isPrimitive = "java/lang/Character", "char", true
	case "D":
		className, primitiveName, isPrimitive = "java/lang/Double", "double", true
	case "F":
		className, primitiveName, isPrimitive = "java/lang/Float", "float", true
	case "I":
		className, primitiveName, isPrimitive = "java/lang/Integer", "int", true
	case "J":
		className, primitiveName, isPrimitive = "java/lang/Long", "long", true
	case "S":
		className, primitiveName, isPrimitive = "java/lang/Short", "short", true
	case "Z":
		className, primitiveName, isPrimitive = "java/lang/Boolean", "boolean", true
	case "V":
		className, primitiveName, isPrimitive = "java/lang/Void", "void", true
	default:
		// Object or Array type
		className = strings.ReplaceAll(typeStr, ".", "/")
		if strings.HasPrefix(className, "L") && strings.HasSuffix(className, ";") {
			className = className[1 : len(className)-1]
		}
	}

	if primitiveName == "" && isPrimitive { // delete this once we know what primitive name is used for
		return nil, fmt.Errorf("invalid primitive type descriptor: %s", typeStr)
	}

	// For primitive types, we need to return the java.lang.Class object that represents
	// the primitive. In Jacobin, this object is created during the <clinit> of the
	// corresponding wrapper class (e.g., java/lang/Integer.<clinit> creates the Class
	// object for 'int'). That Class object is then stored in the wrapper's TYPE static field.
	//
	// To resolve a primitive descriptor (like "I"):
	// 1. We determine the wrapper class name (className = "java/lang/Integer").
	// 2. We look up the "TYPE" static field for that wrapper class in the statics map.
	// 3. If found, its value is the *object.Object representing the primitive class.
	//
	// If the TYPE field is not found, it means the wrapper class hasn't been initialized
	// (its <clinit> hasn't run). In the Jacobin boot sequence, InitializePrimitiveWrappers()
	// is expected to run these initializers. If it fails here, that initialization was missed.
	if isPrimitive {
		staticField, ok := statics.QueryStatic(className, "TYPE")
		if !ok {
			return nil, fmt.Errorf("primitive TYPE field not found for %s (ensure InitializePrimitiveWrappers ran)", className)
		}
		return staticField.Value.(*object.Object), nil
	}

	// For non-primitive types (including arrays), check if it's already loaded.
	k := classloader.MethAreaFetch(className)
	if k == nil {
		// Not loaded.

		// If it's an array type, we shouldn't try to load it from a file.
		if strings.HasPrefix(className, "[") {
			// It's an array. Jacobin preloads primitive arrays.
			// If it's not found in MethArea, it's an object array that hasn't been created yet.
			// We create a synthetic class for it.

			// Create Klass
			k = &classloader.Klass{
				Status: 'L', // Loaded/Linked
				Loader: "bootstrap", // or app
				Data: &classloader.ClData{
					Name: className,
					NameIndex: object.StringPoolIndexFromGoString(className),
					SuperclassIndex: object.StringPoolIndexFromGoString("java/lang/Object"),
				},
			}

			// Create Class Object
			if err := classloader.LoadClassFromNameOnly("java/lang/Class"); err != nil {
				return nil, err
			}
			classObj := object.MakeEmptyObject()
			classObj.KlassName = object.StringPoolIndexFromGoString("java/lang/Class")
			classObj.ThMutex.Lock()
			classObj.FieldTable["name"] = object.Field{Ftype: "Ljava/lang/String;", Fvalue: className}
			classObj.ThMutex.Unlock()

			k.Data.ClassObject = classObj

			classloader.MethAreaInsert(className, k)

			return classObj, nil
		}

		// For regular classes, load them.
		if err := classloader.LoadClassFromNameOnly(className); err != nil {
			return nil, fmt.Errorf("could not load class for descriptor %s: %v", className, err)
		}
		k = classloader.MethAreaFetch(className)
	}

	if k == nil || k.Data == nil {
		return nil, fmt.Errorf("class %s loaded but not found in MethArea", className)
	}

	if k.Data.ClassObject == nil {
		// Create the Class object if it doesn't exist (lazy creation)
		// This usually happens during class loading, but if it's missing, create it now.
		if err := classloader.LoadClassFromNameOnly("java/lang/Class"); err != nil {
			return nil, err
		}
		classObj := object.MakeEmptyObject()
		classObj.KlassName = object.StringPoolIndexFromGoString("java/lang/Class")

		// Set the "name" field for the Class object so tests can verify it
		// This mirrors what we did for primitives in MakeJlcObject
		classObj.ThMutex.Lock()
		classObj.FieldTable["name"] = object.Field{Ftype: "Ljava/lang/String;", Fvalue: className}
		classObj.ThMutex.Unlock()

		// Link it back
		k.Data.ClassObject = classObj
	}

	return k.Data.ClassObject, nil
}
