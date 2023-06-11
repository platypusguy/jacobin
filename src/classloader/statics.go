/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
    "errors"
    "sync"
)

// Statics is a fast-lookup map of static variables and functions. The int64 value
// contains the index into the statics array where the entry is stored.
// Statics are placed into this map only when they are first referenced and resolved.
var Statics = make(map[string]Static)

// var StaticsArray []Static

// Static contains all the various items needed for a static variable or function.
type Static struct {
    Type string // the kind of entity we're dealing with. Can be:
    /*
    	B	byte signed byte (includes booleans)
    	C	char	Unicode character code point (UTF-16)
    	D	double
    	F	float
    	I	int	integer
    	J	long integer
    	L ClassName ;	reference	an instance of class ClassName
    	S	signed short int
    	Z	boolean
    	[x   array of x, where x is any primitive or a reference
    	plus (Jacobin implementation-specific):
    	G   native method (that is, one written in Go)
    	T	string (ptr to an object, but facilitates processing knowing it's a string)
    */
    Value any
    // Type      string         // Type data used for reference variables (i.e., objects, etc.)
    // ValueRef  *object.Object // pointer--might need to change this
    // ValueInt  int64          // holds longs, ints, shorts, chars, booleans, byte
    // ValueFP   float64        // holds doubles and floats
    // ValueStr  string         // string
    // ValueFunc func()         // function pointer
    // CP        *CPool         // the constant pool for the class
}

var staticsMutex = sync.RWMutex{}

// AddStatic adds a static field to the Statics table using a mutex
func AddStatic(name string, s Static) error {
    if name == "" {
        return errors.New("AddStatic: Attempting to add invalid static entry")
    }
    staticsMutex.RLock()
    Statics[name] = s
    staticsMutex.RUnlock()
    return nil
}
