/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package exec

import (
	// "errors"
	// "fmt"
	// "jacobin/log"
	// "strconv"
	"sync"
)

// VTable is the table in which virtual method data is stored for quick reference at
// method invocation. It consists of a map whose key is a string consisting of a
// concatenation of the class name, method name, and method type. The value consists
// of fields used in the execution of the method: ParamSlots indicates how many slots
// on the calling method's operand stack are items for the called method, a pointer
// to a generic function, and a MethType byte, which indicates what kind of method is
// pointed to by the previous field: 'J' = Java method, 'G' = golang method, and 'N'
// which is a native method in the JNI sense of the term.

var VTable map[string]Ventry

type Ventry struct {
	ParamSlots int
	Fu         Function
	MethType   byte
}

type Function func([]interface{})

// VTmutex is used for updates to the VTable because multiple threads could be
// updating it simultaneously.
var VTmutex sync.Mutex
