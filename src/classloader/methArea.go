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
// var MethArea = make(map[string]Klass)
var MethArea *sync.Map
var methAreaSize = 0
var MethAreaMutex sync.RWMutex // All additions or updates to MethArea map come through this mutex

func MethAreaFetch(key string) *Klass {
    v, _ := MethArea.Load(key)
    if v == nil {
        return nil
    }
    return v.(*Klass)
}

func MethAreaInsert(name string, klass *Klass) {
    MethArea.Store(name, klass)
    MethAreaMutex.Lock()
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

func initMethodArea() {
    ma := sync.Map{}
    MethArea = &ma
    methAreaSize = 0
}
