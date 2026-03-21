/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"jacobin/src/object"
	"jacobin/src/types"
	"sync"
)

// instances of java/lang/Class as stored in global.JLCmap
type Jlc struct {
	Lock        sync.RWMutex
	Statics     []string // list of all static fields
	Name        string
	IsPrimitive bool
	KlassPtr    *ClData // points back to the class's data in the method area
}

// JLCmap is a map of java/lang/Class instances for statics and introspection.
// Key is the FQN class name (e.g. "java/lang/String").
// Value is a pointer to the internal Jlc struct.
var JLCmap map[string]*Jlc
var JlcMapLock sync.RWMutex

// InitJlcMap initializes the JLCmap.
// This should be called during classloader initialization.
func InitJlcMap() {
	JlcMapLock.Lock()
	defer JlcMapLock.Unlock()
	JLCmap = make(map[string]*Jlc, 2000)
}

// GetJlcEntry returns the JLC entry for a class.
func GetJlcEntry(className string) *Jlc {
	JlcMapLock.RLock()
	defer JlcMapLock.RUnlock()
	return JLCmap[className]
}

// GetJlcObject fetches the JLC entry for a class and converts
// it into a java/lang/Class object
func GetJlcObject(className string) *object.Object {
	JlcMapLock.RLock()
	defer JlcMapLock.RUnlock()
	jlc := JLCmap[className]
	if jlc == nil {
		return nil
	}

	o := object.MakeEmptyObject()
	o.KlassName = types.StringPoolJavaLangClassIndex
	o.FieldTable["name"] = object.Field{Ftype: types.GolangString,
		Fvalue: object.StringObjectFromGoString(jlc.Name)}
	o.FieldTable["$klass"] = object.Field{Ftype: types.RawGoPointer,
		Fvalue: jlc.KlassPtr} // points to the Klass object in metadata
	o.FieldTable["$statics"] = object.Field{Ftype: types.Array,
		Fvalue: jlc.Statics} // array of static field names for this class
	return o
}

// MakeJlcEntry creates a new JLC entry for a class.
func MakeJlcEntry(className string, primitive bool) *Jlc {
	jlc := Jlc{}
	jlc.Name = className
	klass := MethAreaFetch(className)
	if klass != nil {
		jlc.KlassPtr = klass.Data
	} else {
		jlc.KlassPtr = nil
	}
	jlc.IsPrimitive = primitive
	jlc.Statics = make([]string, 0)
	return &jlc
}
