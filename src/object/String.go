/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package object

import (
	"fmt"
	"jacobin/statics"
	"jacobin/types"
	"os"
	"sort"
	"strings"
	"sync"
)

// Strings are so commonly used in Java, that it makes sense
// to have a means of creating them quickly, rather than building
// them from scratch each time by walking through the constant pool.

var StringClassName = "java/lang/String"
var StringClassRef = "Ljava/lang/String;"
var EmptyString = ""

// NewString creates an empty string. However, it lacks an updated
// Klass field, which due to circularity issues, is updated in
// classloader.MakeString(). DO NOT CALL THIS FUNCTION DIRECTLYy.
// It should be called ONLY by classloader.MakeString
func NewString() *Object {
	s := new(Object)
	s.Mark.Hash = 0
	s.Klass = &StringClassName // java/lang/String
	s.FieldTable = make(map[string]Field)

	// ==== now the fields ====

	// value: the content of the string as array of runes or bytes
	// Note: Post JDK9, this field is an array of bytes, so as to
	// enable compact strings.
	value := make([]byte, 0) // presently empty
	valueField := Field{Ftype: types.ByteArray, Fvalue: value}
	s.FieldTable["value"] = valueField

	// coder has two possible values:
	// LATIN(e.g., 0 = bytes for compact strings) or UTF16(e.g., 1 = UTF16)
	coderField := Field{Ftype: types.Byte, Fvalue: int64(0)}
	s.FieldTable["coder"] = coderField

	// the hash code, which is initialized to 0
	hash := Field{Ftype: types.Int, Fvalue: int64(0)}
	s.FieldTable["hash"] = hash

	// hashIsZero: only true in rare case where compute hash is 0
	hashIsZero := Field{Ftype: types.Bool, Fvalue: types.JavaBoolFalse}
	s.FieldTable["hashIsZero"] = hashIsZero

	// The following static fields are preloaded in statics/LoadStaticsString()
	//   COMPACT_STRINGS (always true for JDK >= 9)
	//   UTF_8.INSTANCE ptr to encoder
	//   ISO_8859_1.INSTANCE ptr to encoder
	//   US_ASCII.INSTANCE ptr to encoder
	//   CodingErrorAction.REPLACE
	//   CASE_INSENSITIVE_ORDER
	//   serialPersistentFields

	return s
}

// NewStringFromGoString converts a go string to a Java string-like
// entity, in which the chars are stored as runes, rather than chars
// or as bytes, depending on the status of COMPACT_STRINGS (which defaults
// to true for JDK >= 9).

func NewStringFromGoString(in string) *Object {
	s := NewString()
	if statics.GetStaticValue("java/lang/String", "COMPACT_STRINGS") == types.JavaBoolFalse {
		s.FieldTable["value"] = Field{types.RuneArray, in}
	} else {
		s.FieldTable["value"] = Field{types.ByteArray, []byte(in)}
	}
	return s
}

// CreateCompactStringFromGoString creates a string in which the chars
// are stored as bytes--that is, a compact string.
func CreateCompactStringFromGoString(in *string) *Object {
	s := NewString()
	s.FieldTable["value"] = Field{types.ByteArray, []byte(*in)}
	return s
}

// convenience method to extract a Go string from a Java string
func GetGoStringFromJavaStringPtr(strPtr *Object) string {
	s := *strPtr
	bytes := s.FieldTable["value"].Fvalue.([]byte)
	return string(bytes)
}

// determine whether an object is a Java string
// assumes that any object whose Klass pointer points to java/lang/String
// is an instance of a Java string
func IsJavaString(unknown any) bool {
	var objPtr *Object

	if unknown == nil {
		return false
	}

	switch unknown.(type) {
	case *Object:
		objPtr = unknown.(*Object)
		break
	default:
		return false
	}

	return *objPtr.Klass == "java/lang/String"
}

/*
--------------------------------
The new string primitives follow
--------------------------------
*/

var stringTable = make(map[string]uint32)
var stringList []string
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

func DumpStringRepo() {
	stringLock.Lock()
	_, _ = fmt.Fprintln(os.Stderr, "\n===== DumpStringRepo BEGIN")
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
			_, _ = fmt.Fprintf(os.Stderr, "%d\t%s\n", stringTable[key], key)
		}
	}
	_, _ = fmt.Fprintln(os.Stderr, "===== DumpStringRepo END")
	stringLock.Unlock()
}
