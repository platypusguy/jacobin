/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classArea

import (
	"jacobin/classloader"
	"jacobin/log"
	"sync"
)

// ClassArea contains all the loaded classes. Key is the class name in java/lang/Object format.
// var ClassArea = make(map[string]Klass)
var ClassArea *sync.Map

func Fetch(key string) *classloader.Klass {
	v, _ := ClassArea.Load(key)
	if v == nil {
		return nil
	}
	return v.(*classloader.Klass)
}

func Insert(name string, klass *classloader.Klass) {
	ClassArea.Store(name, klass)

	if klass.Status == 'F' || klass.Status == 'V' || klass.Status == 'L' {
		_ = log.Log("Class: "+klass.Data.Name+", loader: "+klass.Loader, log.CLASS)
	}
}
