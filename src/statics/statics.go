/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package statics

import (
	"errors"
	"fmt"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/types"
	"os"
	"runtime/debug"
	"sort"
	"sync"
)

// Statics is a fast-lookup map of static variables and functions. The int64 value
// contains the index into the statics array where the entry is stored.
// Statics are placed into this map only when they are first referenced and resolved.
var Statics = make(map[string]Static)

// Static contains all the various items needed for a static variable or function.
type Static struct {
	Type string // see the possible returns in types/javatypes.go
	// the kind of entity we're dealing with. Can be:
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
}

var staticsMutex = sync.RWMutex{}

// AddStatic adds a static field to the Statics table using a mutex
func AddStatic(name string, s Static) error {
	if name == "" {
		errMsg := fmt.Sprintf("AddStatic: Attempting to add static entry with a nil name, type=%s, value=%v", s.Type, s.Value)
		_ = log.Log(errMsg, log.SEVERE)
		return errors.New(errMsg)
	}
	staticsMutex.RLock()
	Statics[name] = s
	staticsMutex.RUnlock()
	return nil
}

// PreloadStatics preloads static fields from java.lang.String and other
// immediately necessary statics. It's called in jvmStart.go
func PreloadStatics() {
	LoadProgramStatics()
	LoadStringStatics()
	LoadStaticsInteger()
}

// LoadProgramStatics loads static fields that the JVM expects to have
// loaded as execution begins.
func LoadProgramStatics() {
	_ = AddStatic("main.$assertionsDisabled",
		Static{Type: types.Int, Value: types.JavaBoolTrue})
}

// LoadStringStatics loads the statics from java/lang/String directly
// into the Statics table as part of the setup operations of Jacobin.
// This is done primarily for speed.
func LoadStringStatics() {
	_ = AddStatic("java/lang/String.COMPACT_STRINGS",
		Static{Type: types.Bool, Value: true})
	_ = AddStatic("java/lang/String.UTF16",
		Static{Type: types.Byte, Value: int64(1)})
	_ = AddStatic("java/lang/String.LATIN1",
		Static{Type: types.Byte, Value: int64(0)})
}

// GetStaticValue: Given the frame, frame stack, and field name,
// return the field contents.
// If successful, return the field value and a nil error;
// Else (error), return a nil field value and the non-nil error.
func GetStaticValue(className string, fieldName string) any {
	var retValue any

	staticName := className + "." + fieldName

	// was this static field previously loaded? Is so, get its location and move on.
	prevLoaded, ok := Statics[staticName]
	if !ok {
		glob := globals.GetGlobalRef()
		glob.ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("GetStaticValue: could not find static: %s", staticName)
		_ = log.Log(errMsg, log.SEVERE)
		return errors.New(errMsg)
	}

	// Field types bool, byte, and int need conversion to int64.
	// The other types are OK as is.
	switch prevLoaded.Value.(type) {
	case bool:
		value := prevLoaded.Value.(bool)
		retValue = types.ConvertGoBoolToJavaBool(value)
	case byte:
		retValue = int64(prevLoaded.Value.(byte))
	case int32:
		retValue = int64(prevLoaded.Value.(int32))
	case int:
		retValue = int64(prevLoaded.Value.(int))
	default:
		retValue = prevLoaded.Value
	}

	return retValue
}

// DumpStatics dumps the contents of the statics table in sorted order to stderr
func DumpStatics() {
	_, _ = fmt.Fprintln(os.Stderr, "\n===== DumpStatics BEGIN")
	// Create an array of keys.
	keys := make([]string, 0, len(Statics))
	for key := range Statics {
		keys = append(keys, key)
	}
	// Sort the keys.
	// All the upper case entries precede all the lower case entries.
	sort.Strings(keys)
	// In key sequence order, display the key and its value.
	for _, key := range keys {
		_, _ = fmt.Fprintf(os.Stderr, "%s     %v\n", key, Statics[key])
	}
	_, _ = fmt.Fprintln(os.Stderr, "===== DumpStatics END")
}
