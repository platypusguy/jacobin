/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024-6 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package stringPool

import (
	"fmt"
	"jacobin/src/globals"
	"jacobin/src/types"
	"os"
	"sort"
)

// Overview of the String Pool Functions
//
// The purpose of the String Pool is to provide a mechanism for Jacobin whereby string values are not
// duplicated. When a request is made to add a second instance of a string which is already in the Pool
// (via GetStringIndex), the existing instance of the string will be used to satisfy the request.
// It is transparent to the caller of String Pool functions whether the presented string argument is
// pre-existing or is new to the Pool.
//
// String Pool components, in the globals package, common across all frames and threads:
// -------------------------------------------------------------------------------------
//
// stringTable = map[string] --> uint32, an index to stringList (initially empty)
// stringList []string - the array of unique strings (initially empty)
// stringNext uint32 - the index of the next available stringList entry (initially 0)
// stringLock sync.Mutex - control modifications to this Pool (initially unlocked)

// Functions are in alpha order

// DumpStringPool(dumpContext string) -
// - Dump the contents of the string Pool.
// - If the context parameter is not "", the context string will be shown at the beginning of the dump.
func DumpStringPool(dumpContext string) {
	globals.StringPoolLock.RLock()
	defer globals.StringPoolLock.RUnlock()
	if len(dumpContext) > 0 {
		_, _ = fmt.Fprintf(os.Stdout, "\n===== DumpStringPool BEGIN context: %s\n", dumpContext)
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
}

// EmptyStringPool is used exclusively for testing. If used in production, remove this comment.
func EmptyStringPool() {
	globals.InitStringPool()
}

// GetStringIndex(arg *string) uint32 returns the string index for the string pointed to by arg.
// If the string is not already present in the pool, it is added. In fact, this is the primary
// function for adding strings to the pool.
func GetStringIndex(arg *string) uint32 {
	if arg == nil {
		nilString := ""
		arg = &nilString
	}

	// Try to get the index from the map.
	globals.StringPoolLock.RLock()
	index, ok := globals.StringPoolTable[*arg]
	globals.StringPoolLock.RUnlock()
	if ok {
		return index // Found it!
	}

	// Not found. Lock and defer unlock.
	globals.StringPoolLock.Lock()
	defer globals.StringPoolLock.Unlock()

	// Add it to the map and the list.
	index = globals.StringPoolNext
	globals.StringPoolTable[*arg] = index
	globals.StringPoolList = append(globals.StringPoolList, *arg)

	// Increment the next available index.
	globals.StringPoolNext++

	return index
}

// GetStringPointer retrieves a pointer to the string at the index into the string pool slice
// Returns nil on index out of range (which is the only possible error)
func GetStringPointer(index uint32) *string {
	if index < globals.StringPoolNext {
		return &globals.StringPoolList[index]
	} else {
		return nil
	}
}

// GetStringPoolSize() uint32: Get the current string Pool size.
func GetStringPoolSize() uint32 {
	// glob := globals.GetGlobalRef()
	return globals.StringPoolNext
}

// PreloadArrayClassesToStringPool() adds the names of the array classes for primitives
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
