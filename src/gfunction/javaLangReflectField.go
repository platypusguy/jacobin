/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

// The implementation of a field in the Java reflection API.

import (
	"jacobin/src/object"
	"jacobin/src/util"
	"sync"
)

type Field struct {
	Class        *object.Object
	Slot         int
	Name         string
	Type         *object.Object
	Modifiers    int
	TrustedFinal bool
	Signature    string
	Annotations  []byte
	// GenericInfo        *FieldRepository
	// FieldAccessor      FieldAccessor
	// OverrideFieldAccessor FieldAccessor
	Root                    *Field
	DeclaredAnnotations     map[string]interface{}
	DeclaredAnnotationsLock sync.Mutex
}

func NewField(obj *object.Object, fieldName string) *Field {
	// f, ok := obj.FieldTable[fieldName]
	// if !ok {
	// 	return nil
	// }

	return &Field{
		Class: obj,
		Name:  fieldName,
		// Type:  f.Ftype,
		// Modifiers:    modifiers,
		// TrustedFinal: trustedFinal,
		// Slot:         slot,
		// Signature:    signature,
		// Annotations:  annotations,
	}
}

func (f *Field) GetDeclaringClass() *object.Object {
	return f.Class
}

func (f *Field) GetName() string {
	return f.Name
}

func (f *Field) GetModifiers() int {
	return f.Modifiers
}

// func (f *Field) IsEnumConstant() bool {
// 	return (f.Modifiers & ModifierEnum) != 0
// }
//
// func (f *Field) IsSynthetic() bool {
// 	return ModifierIsSynthetic(f.Modifiers)
// }

func (f *Field) GetType() *object.Object {
	return f.Type
}

func (f *Field) Equals(other *Field) bool {
	return f.Class == other.Class && f.Name == other.Name && f.Type == other.Type
}

func (f *Field) HashCode() uint64 {
	retVal, _ := util.HashAnything(f.Class)
	return retVal
}

/*
func (f *Field) ToString() string {
	mod := ModifierToString(f.Modifiers)
	return mod + " " + f.Type.GetTypeName() + " " + f.Clazz.GetTypeName() + "." + f.Name
}

func (f *Field) Get(obj *object.Object) (interface{}, error) {
	if f.FieldAccessor == nil {
		return nil, errors.New("field accessor not initialized")
	}
	return f.FieldAccessor.Get(obj)
}

func (f *Field) Set(obj *object.Object, value interface{}) error {
	if f.FieldAccessor == nil {
		return errors.New("field accessor not initialized")
	}
	return f.FieldAccessor.Set(obj, value)
}
*/

// func (f *Field) GetDeclaredAnnotations() map[string]interface{} {
// 	f.DeclaredAnnotationsLock.Lock()
// 	defer f.DeclaredAnnotationsLock.Unlock()
//
// 	if f.DeclaredAnnotations == nil {
// 		f.DeclaredAnnotations = ParseAnnotations(f.Annotations)
// 	}
// 	return f.DeclaredAnnotations
// }
