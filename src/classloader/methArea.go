/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
    "jacobin/log"
    "sync"
)

// MethArea contains all the loaded classes. Key is the class name in java/lang/Object format.
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
        return nil
    }
    return v.(*Klass)
}

// MethAreaInsert adds a class to the method area, using a pointer
// to the parsed class.
func MethAreaInsert(name string, klass *Klass) {
    MethAreaMutex.Lock()
    MethArea.Store(name, klass)
    methAreaSize++
    MethAreaMutex.Unlock()

    if klass.Status == 'F' || klass.Status == 'V' || klass.Status == 'L' {
        _ = log.Log("Class: "+klass.Data.Name+", loader: "+klass.Loader, log.CLASS)
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

// initMethodArea simply initializes MethArea (the method area
// table of loaded classes) and initializes the counter of classes.
func initMethodArea() {
    MethAreaMutex.Lock()
    ma := sync.Map{}
    MethArea = &ma
    methAreaSize = 0
    MethAreaMutex.Unlock()
}
