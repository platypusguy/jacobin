/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

// This file contains the definition/creation of a java/lang/Class object,
// which is a representation of a Java class used for reflection. It is not
// how loaded classes are represented in Jacobin.

package classloader

import (
	"fmt"
	"jacobin/excNames"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/stringPool"
	"jacobin/types"
)

type JavaLangClass struct {
	Name        *string                 // Fully qualified name of the class
	Modifiers   int                     // Access modifiers (e.g., public, private, etc.)
	Fields      map[string]*ClassField  // Map of field names to Field objects
	Methods     map[string]*ClassMethod // Map of method names to Method objects
	Interfaces  []string                // List of implemented interfaces
	SuperClass  *Klass                  // Reference to the superclass as it appears in the method area
	Annotations map[string]*Annotation  // Map of annotations
	IsPrimitive bool                    // Indicates if the class represents a primitive type
	IsArray     bool                    // Indicates if the class represents an array type
	IsEnum      bool                    // Indicates if the class represents an enum type
	IsInterface bool                    // Indicates if the class represents an interface
	ClassLoader *object.Object          // Reference to the class loader
	PackageName string                  // Package name of the class
}

// == Componnents of a Class ==
type Annotation struct {
	Name       string            // Fully qualified name of the annotation
	Attributes map[string]string // Map of attribute names to their values
}

type ClassMethod struct {
	Name           string                 // Name of the method
	ReturnType     string                 // Fully qualified name of the return type
	ParameterTypes []string               // List of fully qualified names of parameter types
	Modifiers      int                    // Access modifiers (e.g., public, private, etc.)
	Exceptions     []string               // List of fully qualified names of exceptions thrown
	Annotations    map[string]*Annotation // Map of annotations
	IsStatic       bool                   // Indicates if the method is static
	IsAbstract     bool                   // Indicates if the method is abstract
	IsFinal        bool                   // Indicates if the method is final
	IsSynchronized bool                   // Indicates if the method is synchronized
}

type ClassField struct {
	Name        string                 // Name of the field
	Type        string                 // Fully qualified name of the field's type
	Modifiers   int                    // Access modifiers (e.g., public, private, etc.)
	Annotations map[string]*Annotation // Map of annotations
	IsStatic    bool                   // Indicates if the field is static
	IsFinal     bool                   // Indicates if the field is final
	IsVolatile  bool                   // Indicates if the field is volatile
	IsTransient bool                   // Indicates if the field is transient
}

func ClassFromInstance(obj *object.Object) *JavaLangClass {
	// Create a new Class object from an instance of a class.
	if obj == nil {
		return nil
	}
	objName := stringPool.GetStringPointer(obj.KlassName)

	cl := JavaLangClass{
		Name:        objName,
		Modifiers:   0, // Default to 0, can be set later
		Fields:      make(map[string]*ClassField),
		Methods:     make(map[string]*ClassMethod),
		Interfaces:  []string{},
		SuperClass:  nil, // Default to nil, can be set later
		Annotations: make(map[string]*Annotation),
		IsPrimitive: false, // Default to false, can be set later
		IsArray:     false, // Default to false, can be set later
		IsEnum:      false, // Default to false, can be set later
		IsInterface: false, // Default to false, can be set later
		ClassLoader: nil,   // Default to nil, can be set later
	}

	klass := MethAreaFetch(*objName)
	if klass == nil {
		// this should never happen as getClass is alwasy called on an existing object
		errMsg := fmt.Sprintf(
			"%s\n\tClassFromInstance(): Class %s not found in method area",
			excNames.JVMexceptionNamesJacobin[excNames.InternalException], *objName)
		globals.GetGlobalRef().FuncThrowException(excNames.InternalException, errMsg)
		return nil
	}

	superclassName := stringPool.GetStringPointer(klass.Data.SuperclassIndex)
	cl.SuperClass = MethAreaFetch(*superclassName)

	value, ok := obj.FieldTable["value"]
	if ok {
		valueType := value.Ftype
		if types.IsArray(valueType) {
			cl.IsArray = true
		}
	}

	for _, fld := range klass.Data.Fields {
		fldName := fld.NameStr
		cl.Fields[fldName] = &ClassField{
			Name:        fldName,
			Type:        fld.DescStr,
			Modifiers:   0, // Default to 0, can be set later
			Annotations: make(map[string]*Annotation),
			IsStatic:    fld.IsStatic,
			IsFinal:     false, // Default to false, can be set later
			IsVolatile:  false, // Default to false, can be set later
			IsTransient: false, // Default to false, can be set later
		}
	}

	for key, _ := range klass.Data.MethodTable {
		cl.Methods[key] = &ClassMethod{
			Name:           key,
			IsStatic:       false, // Default to false, can be set later
			IsAbstract:     false, // Default to false, can be set later
			IsFinal:        false, // Default to false, can be set later
			IsSynchronized: false, // Default to false, can be set later
		}
	}

	return &cl
}
