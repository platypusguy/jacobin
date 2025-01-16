/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"fmt"
	"os"
	"sort"
	"sync"
)

// GMT = global method table. It's a map of fully qualified method (FQN) name to
// method data, called gmtEntry. The gmtEntry consists of a byte identifying whether
// the method is a Java method ('J') or a golang-implemented native method ('G').
// he other field is a pointer to the method data.
//
// When a function is invoked, the lookup mechanism checks the GMT for a pointer
// to the method data. The method in the class is never called directly--only ever
// through a GMT look-up. This allows for the method to be overridden in a subclass,
// and for us to insert go methods in lieu of JDK native methods
//
// Entries are added to the GMT when a class is instantiated. In other words, it's
// lazy loading.

var GMT = make(gmt)

type gmt map[string]GmtEntry
type GmtEntry struct {
	MethData interface{} // pointer to the method data
	MType    byte        // method type, G = Go method, J = Java method
}

// // JmEntry is the entry in the Mtable for Java methods.
// type JmEntry struct {
// 	AccessFlags int
// 	MaxStack    int
// 	MaxLocals   int
// 	Code        []byte
// 	// Exceptions  []uint16 prior to JACOBIN-575
// 	Exceptions []CodeException // the exception data stored in the code attribute.
// 	Attribs    []Attr
// 	params     []ParamAttrib
// 	CodeAttr   CodeAttrib
// 	deprecated bool
// 	Cp         *CPool
// }

// Function is the generic-style function used for Go entries: a function that accepts a
// slice of empty interfaces and returns an empty interface
// type Function func([]interface{}) interface{}

// MTmutex is used for updates to the MTable because multiple threads could be
// updating it simultaneously.
var GMTmutex sync.Mutex

// adds an entry to the GMT, using a mutex
func GmtAddEntry(key string, mte GmtEntry) {
	// fmt.Printf("DEBUG gmt.go AddEntry key=%s, gmtEntry=%s\n", key, string(mte.MType))
	GMTmutex.Lock()
	GMT[key] = mte
	GMTmutex.Unlock()
}

// DumpMTable dumps the contents of GMT in sorted order to stderr
func DumpGmt() {
	_, _ = fmt.Fprintln(os.Stderr, "\n===== DumpGMT BEGIN")
	// Create an array of keys.
	keys := make([]string, 0, len(GMT))
	for key := range MTable {
		keys = append(keys, key)
	}

	// Sort the keys (FQNs).
	// All the upper case entries precede all the lower case entries.
	sort.Strings(keys)

	// In key sequence order, display the key and its value.
	for _, key := range keys {
		entry := MTable[key]
		_, _ = fmt.Fprintf(os.Stderr, "%s   %s\n", string(entry.MType), key)
	}
	_, _ = fmt.Fprintln(os.Stderr, "===== DumpGMT END")
}
