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

// ClassArea contains all the loaded classes. Key is the class name in java/lang/Object format.
// var ClassArea = make(map[string]Klass)
var ClassArea *sync.Map
var classAreaSize = 0

func ClassAreaFetch(key string) *Klass {
    v, _ := ClassArea.Load(key)
    if v == nil {
        return nil
    }
    return v.(*Klass)
}

func ClassAreaInsert(name string, klass *Klass) {
    ClassArea.Store(name, klass)
    classAreaSize++

    if klass.Status == 'F' || klass.Status == 'V' || klass.Status == 'L' {
        _ = log.Log("Class: "+klass.Data.Name+", loader: "+klass.Loader, log.CLASS)
    }
}

// Size returns the number of entries in ClassArea.
// Because the golang's sync.Map does not have a len() function
// we have to track our additions with a counter, which is
// returned here.
func ClassAreaSize() int {
    return classAreaSize
}

func initMethodArea() {
    ma := sync.Map{}
    ClassArea = &ma
    classAreaSize = 0
}
