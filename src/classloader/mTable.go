/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-4 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"fmt"
	"os"
	"sort"
	"sync"
)

// MTable value consists of a byte identifying whether the method is a Java method
// ('J'), that is, a method that is executed by executing bytecodes, or a golang-style
// native method ('G'), i.e., a golang function that is being used as a stand-in for
// the named Java method. In most contexts, this would be called a native method,
// but that term is used in a different context in Java (see JNI), so avoided here.
//
// The second field in the value is an empty interface, which is Go's way of
// implementing generics. Ultimately, this mechanism supports two types of entries:
// one for each kind of method.
//
// When a function is invoked, the lookup mechanism first checks the MTable, and
// if an entry is found, that entry is what is executed. If no entry is found,
// the search goes to the class and failing that to the superclass, etc. Once the
// method is located, it's added to the MTable so that all future invocations will
// result in fast look-ups in the MTable.
var MTable = make(MT)

// MT is a type alias for the MTable. It's simply syntactic sugar in context.
type MT = map[string]MTentry

// MTentry is described in detail in the comments to MTable
type MTentry struct {
	Meth  MData // the method data
	MType byte  // method type, G = Go method, J = Java method
}

// MData can be a GMeth or a JmEntry (method in Go or Java, respectively)
type MData interface{}

// JmEntry is the entry in the Mtable for Java methods.
type JmEntry struct {
	AccessFlags int
	MaxStack    int
	MaxLocals   int
	Code        []byte
	// Exceptions  []uint16 prior to JACOBIN-575
	Exceptions []CodeException // the exception data stored in the code attribute.
	Attribs    []Attr
	params     []ParamAttrib
	CodeAttr   CodeAttrib
	deprecated bool
	Cp         *CPool
}

// Function is the generic-style function used for Go entries: a function that accepts a
// slice of empty interfaces and returns an empty interface
type Function func([]interface{}) interface{}

// MTmutex is used for updates to the MTable because multiple threads could be
// updating it simultaneously.
var MTmutex sync.RWMutex

// adds an entry to the MTable, using a mutex
func AddEntry(tbl *MT, key string, mte MTentry) {
	mt := *tbl

	MTmutex.Lock()
	defer MTmutex.Unlock()
	mt[key] = mte
}

// GetMtableEntry returns the entry for the given key, or nil if it doesn't exist.
func GetMtableEntry(key string) MTentry {
	MTmutex.RLock()
	defer MTmutex.RUnlock()
	entry, ok := MTable[key]
	if ok {
		return entry
	}
	return MTentry{}
}

// DumpMTable dumps the contents of MTable in sorted order to stderr
func DumpMTable() {
	MTmutex.RLock()
	defer MTmutex.RUnlock()
	_, _ = fmt.Fprintln(os.Stderr, "\n===== DumpMTable BEGIN")
	// Create an array of keys.
	keys := make([]string, 0, len(MTable))
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
	_, _ = fmt.Fprintln(os.Stderr, "===== DumpMTable END")
}
