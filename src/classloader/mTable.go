/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"sync"
)

// MTable value consists of a byte identifying whether the method is a Java method
// ('J'), that is, a method that is executed by executing bytecodes, or a go
// method ('G'), which is a Go funciton that is being used as a stand-in for
// the named Java method. In most contexts, this would be called a native method,
// but that term is used in a different context in Java (see JNI), so avoided here.
//
// The second field in the value is an empty interface, which is Go's way of
// implementing generics. Ultimately, this mechanism supports two types of entries--
// one for each kind of method.
//
// When a function is invoked, the lookup mechanism first checks the MTable, and
// if an entry is found, that entry is what is executed. If no entry is found,
// the search goes to the class and faiing that to the superclass, etc. Once the
// method is located it's added to the MTable so that all future invocations will
// result in fast look-ups in the MTable.
var MTable = make(MT)

// MT is a type alias for the MTable. It's simply syntactic sugar in context.
type MT = map[string]MTentry

// MTentry is described in detail in the comments to MTable
type MTentry struct {
	Meth  mData // the method data
	MType byte  // method type, G = Go method, J = Java method
}

// mData can be a GmEntry or a JmEntry (method in Go or Java, respectively)
type mData interface{}

// GmEntry is the entry in the MTable for Go functions. See MTable comments for details.
// Fu is a go function. All go functions accept a possibly empty slice of interface{} and
// return a possibly nil interface{}
type GmEntry struct {
	ParamSlots int
	Fu         func([]interface{}) interface{}
}

// JmEntry is the entry in the Mtable for Java methods.
type JmEntry struct {
	accessFlags int
	MaxStack    int
	MaxLocals   int
	Code        []byte
	exceptions  []CodeException
	attribs     []Attr
	params      []ParamAttrib
	deprecated  bool
	Cp          *CPool
}

// Function is the generic-style function used for Go entries: a function that accepts a
// slice of empty interfaces and returns nothing (b/c all returns are pushed onto the
// stack rather than actually returned to a caller).
type Function func([]interface{}) interface{}

// MTmutex is used for updates to the MTable because multiple threads could be
// updating it simultaneously.
var MTmutex sync.Mutex

// MTableLoadNatives loads the Go methods from files that contain them. It does this
// by calling the Load_* function in each of those files to load whatever Go functions
// they make available.
func MTableLoadNatives() {
	loadlib(&MTable, Load_System_Io_PrintStream()) // load the java.io.prinstream golang functions
	loadlib(&MTable, Load_Lang_System())           // load the java.lang.system golang functions
}

func loadlib(tbl *MT, libMeths map[string]GMeth) {
	for key, val := range libMeths {
		gme := GmEntry{}
		gme.ParamSlots = val.ParamSlots
		gme.Fu = val.GFunction

		tableEntry := MTentry{
			MType: 'G',
			Meth:  gme,
		}

		addEntry(tbl, key, tableEntry)
	}
}

// adds an entry to the MTable, using a mutex
func addEntry(tbl *MT, key string, mte MTentry) {
	mt := *tbl

	MTmutex.Lock()
	mt[key] = mte
	MTmutex.Unlock()
}
