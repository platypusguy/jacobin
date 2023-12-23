/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-23 by Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"errors"
	"fmt"
	"jacobin/log"
	"jacobin/types"
	"sync"
	"time"
)

// MethArea contains all the loaded classes. Key is the class name in java/lang/Object format.
// var MethArea = make(map[string]Klass)
var MethArea *sync.Map
var methAreaSize = 0
var MethAreaMutex sync.RWMutex // All additions or updates to MethArea map come through this mutex

// MethAreaFetch retrieves a pointer to a loaded class from the
// method area. In the event the class is not present there, the
// function returns nil.
func MethAreaFetch(key string) *Klass {
	MethAreaMutex.RLock()
	v, _ := MethArea.Load(key)
	MethAreaMutex.RUnlock()
	if v == nil {
		_ = log.Log("MethAreaFetch: key("+key+") --> nil", log.CLASS)
		return nil
	}
	_ = log.Log("MethAreaFetch: key("+key+") --> not nil", log.CLASS)
	return v.(*Klass)
}

// MethAreaInsert adds a class to the method area, using a pointer
// to the parsed class.
func MethAreaInsert(name string, klass *Klass) {
	_ = log.Log("MethAreaInsert: key("+name+")", log.CLASS)
	MethAreaMutex.Lock()
	MethArea.Store(name, klass)
	methAreaSize++
	MethAreaMutex.Unlock()

	if klass.Status == 'F' || klass.Status == 'V' || klass.Status == 'L' {
		_ = log.Log("Method area insert: "+klass.Data.Name+", loader: "+klass.Loader, log.CLASS)
	}
}

// Size returns the number of entries in MethArea.
// Because the golang's sync.Map does not have a len() function
// we have to track our additions with a counter, which is
// returned here.
func MethAreaSize() int {
	MethAreaMutex.RLock()
	size := methAreaSize
	MethAreaMutex.RUnlock()
	return size
}

// this function deletes an entry in the method area
// at present, it is used only in testing
func MethAreaDelete(key string) {
	if MethAreaFetch(key) != nil {
		MethAreaMutex.Lock()
		MethArea.Delete(key)
		methAreaSize--
		MethAreaMutex.Unlock()
	}
}

// Wait for klass.Status to no longer be "I"
// TODO: must be a better way to do this!
func WaitForClassStatus(className string) error {
	_ = log.Log("WaitForClassStatus: class name: "+className, log.CLASS)
	klass := MethAreaFetch(className)
	if klass == nil { // class not there yet
		time.Sleep(100 * time.Millisecond) // sleep 100 milliseconds
		klass = MethAreaFetch(className)
		if klass == nil {
			msg := fmt.Sprintf("WaitClassStatus: Timeout waiting for class %s to load", className)
			return errors.New(msg)
		}
	}
	if klass.Status == 'I' { // class is being initialized by a loader, so wait
		time.Sleep(100 * time.Millisecond) // sleep 100 milliseconds
		klass = MethAreaFetch(className)
		if klass.Status == 'I' {
			msg := fmt.Sprintf("WaitClassStatus: Timeout waiting for class {%s} status", className)
			return errors.New(msg)
		}
	}
	return nil
}

// InitMethodArea simply initializes MethArea (the method area
// table of loaded classes) and initializes the counter of classes.
func InitMethodArea() {
	MethAreaMutex.Lock()
	ma := sync.Map{}
	MethArea = &ma
	methAreaSize = 0
	MethAreaMutex.Unlock()

	// preload the synthetic classes for arrays
	MethAreaPreload()
}

// MethAreaPreload preloads the synthetic entries for array types into
// the method area.
func MethAreaPreload() {
	emptyKlass := Klass{
		Status: 'N', // N = instantiated
		Loader: "bootstrap",
		Data:   &ClData{Superclass: "java/lang/Object"}, // empty class info
	}
	classesToPreload := []string{
		types.ByteArray, types.FloatArray, types.IntArray,
		types.RefArray, types.RuneArray,
	}

	for _, x := range classesToPreload {
		k := emptyKlass
		k.Data.Name = x
		MethAreaInsert(x, &emptyKlass)
	}
}
