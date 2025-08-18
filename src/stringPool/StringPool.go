/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package stringPool

/*
Overview of the String Pool Functions

The purpose of the String Pool is to provide a mechanism for Jacobin whereby string values are not duplicated.
When a request is made to add a second instance of a string which is already in the Pool (GetStringIndex),
the existing instance of the string will be used to satisfy the request.
It is transparent to the caller of String Pool functions whether the presented string argument is
pre-existing or is new to the Pool.

String Pool components, in the globals package, common across all frames and threads:
-------------------------------------------------------------------------------------

stringTable = map[string] --> uint32, an index to stringList (initially empty)
stringList []string - the array of unique strings (initially empty)
stringNext uint32 - the index of the next available stringList entry (initially 0)
stringLock sync.Mutex - control modifications to this Pool (initially unlocked)

Mid-level Functions:
--------------------

MakeEmptyStringObject() *Object
  - Create an object.Object for a java/lang/string with an empty FieldTable.
  - Return a pointer to the object.

NewPoolStringFromGoString(str string) *Object

	Given a Go string,
	* Store the string in the pool.
	* Create an object.Object containing the index to the pool string.
	* Return a pointer to the object.

GetGoStringFromObject(strPtr *Object) string

	Given a pointer to an object.Object containing an index to a pool string, return a Go string.

Primitive Functions:
--------------------

GetStringIndex(arg *string) uint32 -
  - Given a pointer to a Go string, add the string to the pool if the string is not already present.
  - Whether new or existing, return the index for the string for subsequent direct retrievals using stringList.

GetStringPoolSize() uint32

	Get the current string Pool size.

EmptyStringPool() - Put the string Pool into an initial state. Useful in testing!

DumpStringPool(context string) -
  - Dump the contents of the string Pool.
  - If the context parameter is not "", the context string will be shown at the beginning of the dump.
*/

import (
	"fmt"
	"jacobin/src/globals"
	"jacobin/src/types"
	"os"
	"sort"
)

/*
-----------------------------------
The string pool primitive functions
-----------------------------------
*/

func GetStringIndex(arg *string) uint32 {
	if arg == nil {
		nilString := ""
		arg = &nilString
	}

	index, ok := globals.StringPoolTable[*arg]
	if ok {
		return index
	}
	globals.StringPoolLock.Lock()
	index = globals.StringPoolNext
	globals.StringPoolTable[*arg] = index
	globals.StringPoolList = append(globals.StringPoolList, *arg)
	globals.StringPoolNext++
	globals.StringPoolLock.Unlock()
	return index
}

// GetStringPointer retrieves a pointer to the string at the index into the string pool slice
// Returns nil on index out of range (which is the only possible error)
func GetStringPointer(index uint32) *string {
	if index < uint32(len(globals.StringPoolList)) {
		return &globals.StringPoolList[index]
	} else {
		return nil
	}
}

func GetStringPoolSize() uint32 {
	// glob := globals.GetGlobalRef()
	return globals.StringPoolNext
}

// EmptyStringPool is used exclusively for testing. If used in production, remove this comment.
func EmptyStringPool() {
	globals.InitStringPool()
}

func DumpStringPool(context string) {
	globals.StringPoolLock.Lock()
	if len(context) > 0 {
		_, _ = fmt.Fprintf(os.Stdout, "\n===== DumpStringPool BEGIN context: %s\n", context)
	} else {
		_, _ = fmt.Fprintln(os.Stdout, "\n===== DumpStringPool BEGIN")
	}
	// Create an array of keys.
	keys := make([]string, 0, len(globals.StringPoolTable))
	for key := range globals.StringPoolTable {
		keys = append(keys, key)
	}
	// Sort the keys.
	// All the upper case entries precede all the lower case entries.
	sort.Strings(keys)
	// In key sequence order, display the key and its value.
	for _, key := range keys {
		_, _ = fmt.Fprintf(os.Stdout, "%d\t%s\n", globals.StringPoolTable[key], key)
	}
	_, _ = fmt.Fprintln(os.Stdout, "===== DumpStringPool END")
	globals.StringPoolLock.Unlock()
}

func PreloadArrayClassesToStringPool() {
	arrayClassesToPreload := []string{
		types.BoolArray,
		types.ByteArray,
		types.DoubleArray,
		types.FloatArray,
		types.IntArray,
		types.LongArray,
		types.RefArray,
		types.RuneArray,
	}

	for _, className := range arrayClassesToPreload {
		_ = GetStringIndex(&className)
	}
}
