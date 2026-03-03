/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
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
