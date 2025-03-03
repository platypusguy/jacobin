/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package statics

import (
	"errors"
	"fmt"
	"jacobin/excNames"
	"jacobin/globals"
	"jacobin/types"
	"jacobin/util"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"testing"
)

const flagTraceStatics = false

// Statics is a fast-lookup map of static variables and functions. The int64 value
// contains the index into the statics array where the entry is stored.
// Statics are placed into this map only when they are first referenced and resolved.
var Statics = make(map[string]Static)

// Static contains all the various items needed for a static variable or function.
type Static struct {
	Type string // see the possible returns in types/javatypes.go
	// the kind of entity we're dealing with. Can be:
	// B	byte signed byte (includes booleans)
	// C	char	Unicode character code point (UTF-16)
	// D	double
	// F	float
	// I	int	integer
	// J	long integer
	// LClassName;	reference to an instance of an object of class ClassName
	// S	signed short int
	// Z	boolean
	// [x   array of x, where x is any primitive or a reference
	//
	// plus (Jacobin implementation-specific):
	// G   native method (that is, one written in Go)
	// T   string (ptr to an object, but facilitates processing knowing it's a string)
	// GS  Go I/O stream (os.Stdin, os.Stdout, os.Stderr)

	Value any
}

var staticsMutex = sync.RWMutex{}

// AddStatic adds a static field to the Statics table using a mutex
// name: className.fieldName
func AddStatic(name string, s Static) error {
	if name == "" {
		errMsg := fmt.Sprintf("AddStatic: Attempting to add static entry with a nil name, type=%s, value=%v", s.Type, s.Value)
		globals.GetGlobalRef().FuncThrowException(excNames.InvalidTypeException, errMsg)
		return errors.New(errMsg)
	}
	staticsMutex.RLock()
	Statics[name] = s
	staticsMutex.RUnlock()
	if flagTraceStatics && !util.IsFilePartOfJDK(&name) {
		if !testing.Testing() {
			_, _ = fmt.Fprintf(os.Stderr, ">>>trace>>>AddStatic: Adding static entry with name=%s, value=%v\n", name, s.Value)
		}
	}
	return nil
}

// PreloadStatics preloads static fields from java.lang.String and other
// immediately necessary statics. It's called in jvmStart.go
func PreloadStatics() {
	LoadProgramStatics()
	// java.lang.*
	LoadStaticsByte()
	LoadStaticsCharacter()
	LoadStaticsDouble()
	LoadStaticsFloat()
	LoadStaticsInteger()
	LoadStaticsLong()
	LoadStaticsMath()
	LoadStaticsShort()
	LoadStaticsStrictMath()
	LoadStaticsString()
}

// LoadProgramStatics loads static fields that the JVM expects to have
// loaded as execution begins.
func LoadProgramStatics() {
	_ = AddStatic("main.$assertionsDisabled",
		Static{Type: types.Int, Value: types.JavaBoolTrue})
}

// LoadStaticsString loads the statics from java/lang/String directly
// into the Statics table as part of the setup operations of Jacobin.
// This is done primarily for speed.
func LoadStaticsString() {
	_ = AddStatic("java/lang/String.COMPACT_STRINGS",
		Static{Type: types.Bool, Value: true})
	_ = AddStatic("java/lang/String.UTF16",
		Static{Type: types.Byte, Value: int64(1)})
	_ = AddStatic("java/lang/String.LATIN1",
		Static{Type: types.Byte, Value: int64(0)})
	_ = AddStatic("sun/nio/cs/UTF_8.INSTANCE",
		Static{Type: types.Ref, Value: nil})
	_ = AddStatic("sun/nio/cs/.ISO_8859_1.INSTANCE",
		Static{Type: types.Ref, Value: nil})
	_ = AddStatic("sun/nio/cs/.US_ASCII.INSTANCE",
		Static{Type: types.Ref, Value: nil})
	_ = AddStatic("java/nio/charset/CodingErrorAction.REPLACE",
		Static{Type: types.Ref, Value: nil})
	_ = AddStatic("java/lang/String.serialPersistentFields",
		Static{Type: types.Ref, Value: nil})
	// next entry points to a comparator. Might be useful to fill in later
	_ = AddStatic("java/lang/String.CASE_INSENSITIVE_ORDER",
		Static{Type: types.Ref, Value: nil})

}

// GetStaticValue: Given the class name and field name,
// return the field contents.
// If successful, return the field value and a nil error;
// Else (error), return errors.New(errMsg).
func GetStaticValue(className string, fieldName string) any {
	var retValue any

	staticName := className + "." + fieldName

	// was this static field previously loaded? Is so, get its location and move on.
	prevLoaded, ok := Statics[staticName]
	if !ok {
		glob := globals.GetGlobalRef()
		glob.ErrorGoStack = string(debug.Stack())
		errMsg := fmt.Sprintf("GetStaticValue: could not find static: %s", staticName)
		glob.FuncThrowException(excNames.InvalidTypeException, errMsg)
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
	case types.JavaByte:
		retValue = int64(prevLoaded.Value.(types.JavaByte))
	case int32:
		retValue = int64(prevLoaded.Value.(int32))
	case int:
		retValue = int64(prevLoaded.Value.(int))
	default:
		retValue = prevLoaded.Value
	}

	return retValue
}

const SelectAll = int64(1)
const SelectClass = int64(2)
const SelectUser = int64(3)

// DumpStatics dumps the contents of the statics table in sorted order to stderr
func DumpStatics(from string, selection int64, className string) {
	_, _ = fmt.Fprintf(os.Stderr, "\n===== DumpStatics BEGIN, from=\"%s\", selection=%d, className=\"%s\"\n",
		from, selection, className)

	if selection == SelectClass && len(className) < 1 {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR, no class name specified!\n===== DumpStatics END")
		return
	}

	// Create a slice of keys.
	keys := make([]string, 0, len(Statics))
	for key := range Statics {
		keys = append(keys, key)
	}

	// Sort the keys, case-insensitive.
	globals.SortCaseInsensitive(&keys)

	// Process the key slice, depending on selection value.
	var value string
	for _, key := range keys {
		st := Statics[key]

		// Filter switch.
		switch selection {
		case SelectClass:
			left := strings.Split(key, ".")
			if left[0] != className {
				continue
			}
		case SelectUser:
			if strings.HasPrefix(key, "java/") || strings.HasPrefix(key, "jdk/") ||
				strings.HasPrefix(key, "javax/") || strings.HasPrefix(key, "sun") {
				continue
			}
		case SelectAll: // passthrough: nothing here to filter
		default:
			_, _ = fmt.Fprintf(os.Stderr, "ERROR, illegal selection specified: %d!\n===== DumpStatics END", selection)
			return
		}

		// due to circular dependence on object, we can't test directly for object.Null, so we do this.
		if (strings.HasPrefix(st.Type, "L") || strings.HasPrefix(st.Type, "[")) && st.Value == nil {
			value = "<null>"
		} else {
			switch st.Type {
			case types.Bool:
				switch st.Value.(type) {
				case bool:
					if st.Value.(bool) {
						value = "true"
					} else {
						value = "false"
					}
				default:
					if st.Value.(int64) == 1 {
						value = "true"
					} else {
						value = "false"
					}
				}
			case types.Byte:
				value = fmt.Sprintf("0x%02x", st.Value)
			case types.Char, types.Rune:
				value = fmt.Sprintf("'%c'", st.Value)
			case "Ljava/lang/String;":
				// TODO: Avoiding a circularity issue between packages statics and object. What a pity!
				value = fmt.Sprintf("%v", st.Value)
			default:
				value = fmt.Sprintf("%v", st.Value)
			}

		}

		// Prefix name with statics designation (X).
		if strings.HasPrefix(st.Type, "X") {
			st.Type = st.Type[1:] // remove X type prefix, which says field is static
		}

		// Print it.
		_, _ = fmt.Fprintf(os.Stderr, "%-40s   %s %s\n", key, st.Type, value)
	}
	_, _ = fmt.Fprintln(os.Stderr, "===== DumpStatics END")
}
