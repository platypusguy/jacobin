/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package object

import (
	"fmt"
	"jacobin/types"
	"os"
	"sort"
	"strings"
	"sync"
	"unsafe"
)

/*
------------------------------------------
The string repo mid-level functions follow
------------------------------------------
*/

// MakeEmptyStringObject() creates an empty object.Object.
// It is expected that the caller will fill in the FieldTable.
func MakeEmptyStringObject() *Object {
	object := Object{}
	ptrObject := uintptr(unsafe.Pointer(&object))
	object.Mark.Hash = uint32(ptrObject)

	// TODO: Change object.Klass to be type uint32
	object.Klass = &StringClassName // java/lang/String

	// initialize the map of this object's fields
	object.FieldTable = make(map[string]Field)
	return &object
}

func NewRepoStringFromGoString(str string) *Object {
	objPtr := MakeEmptyStringObject()
	/* TODO - Is ignoring the COMPACT_STRINGS flag valid?
	if statics.GetStaticValue("java/lang/String", "COMPACT_STRINGS") == types.JavaBoolFalse {
		objPtr.FieldTable["value"] = Field{types.RuneArray, in}
	} else {
		objPtr.FieldTable["value"] = Field{types.StringIndex, GetStringIndex(&in)}
	}
	*/
	objPtr.FieldTable["value"] = Field{types.StringIndex, GetStringIndex(&str)}
	return objPtr
}

// convenience method to extract a Go string from a repository string
func GetGoStringFromObject(strPtr *Object) string {
	obj := *strPtr
	index := obj.FieldTable["value"].Fvalue.(uint32)
	return *GetStringPointer(index)
}

/*
------------------------------------------
The string repo primitive functions follow
------------------------------------------
*/

var stringTable = make(map[string]uint32)
var stringList []string = nil
var stringNext = uint32(0)
var stringLock sync.Mutex

func GetStringIndex(arg *string) uint32 {
	index, ok := stringTable[*arg]
	if ok {
		return index
	}
	stringLock.Lock()
	index = stringNext
	stringTable[*arg] = index
	stringList = append(stringList, *arg)
	stringNext++
	stringLock.Unlock()
	return index
}

func GetStringPointer(index uint32) *string {
	return &stringList[index]
}

func GetStringRepoSize() uint32 {
	return stringNext
}

func EmptyStringRepo() {
	stringLock.Lock()
	stringTable = make(map[string]uint32)
	stringNext = 0
	stringList = nil
	stringLock.Unlock()
}

func DumpStringRepo(context string) {
	stringLock.Lock()
	if len(context) > 0 {
		_, _ = fmt.Fprintf(os.Stdout, "\n===== DumpStringRepo BEGIN context: %s\n", context)
	} else {
		_, _ = fmt.Fprintln(os.Stdout, "\n===== DumpStringRepo BEGIN")
	}
	// Create an array of keys.
	keys := make([]string, 0, len(stringTable))
	for key := range stringTable {
		keys = append(keys, key)
	}
	// Sort the keys.
	// All the upper case entries precede all the lower case entries.
	sort.Strings(keys)
	// In key sequence order, display the key and its value.
	for _, key := range keys {
		if !strings.HasPrefix(key, "java/") && !strings.HasPrefix(key, "jdk/") &&
			!strings.HasPrefix(key, "javax/") && !strings.HasPrefix(key, "sun") {
			_, _ = fmt.Fprintf(os.Stdout, "%d\t%s\n", stringTable[key], key)
		}
	}
	_, _ = fmt.Fprintln(os.Stdout, "===== DumpStringRepo END")
	stringLock.Unlock()
}
