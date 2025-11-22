/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022-3 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package thread

import (
	"container/list"
	"jacobin/src/globals"
	"math/rand"
	"time"
)

// Creates a JVM program execution thread. These threads are extremely limited.
// They basically hold a Stack of frames. They push and popFrame frames as required.
// They begin execution; they exit when execution ends.

type ExecThread struct {
	ID    int        // the thread ID
	Stack *list.List // the JVM Stack (frame stack, that is) for this thread
	Trace bool       // do we trace instructions?
}

// ThreadNumber is a monotonic counter to number threads at creation
// is set only in the threadNumbering functions in javalang/thread.go
var ThreadNumber int64 = 0

// CreateThread creates an execution thread and initializes it with default values
// All Jacobin execution threads *must* use this function to create a thread
func CreateThread() ExecThread {
	gl := globals.GetGlobalRef()
	t := ExecThread{}
	if gl.JacobinName == "test" || gl.JacobinName == "testWithoutShutdown" {
		t.ID = int(time.Now().UnixNano()) + rand.Int()
	} else {
		ID := gl.FuncInvokeGFunction("java/lang/Thread.ThreadNumberingNext()J", nil).(int64)
		t.ID = int(ID)
	}
	t.Stack = nil
	t.Trace = false
	return t
}

//
// // threads are assigned a monotonically incrementing integer ID. This function
// // increments the counter and returns its value as the integer ID to use
// func IncrementThreadNumber() int {
// 	glob := globals.GetGlobalRef()
//
// 	glob.ThreadLock.Lock()
// 	forCaller := glob.ThreadNumber + 1 // ensure that caller sees this one
// 	glob.ThreadNumber = forCaller
// 	glob.ThreadLock.Unlock() // I don't care if glob.ThreadNumber races ahead
// 	return forCaller
// }

// ======= Items for Java threads ======
// Thread state constants matching Java's Thread.State enum
type State int64

const (
	NEW           State = 0
	RUNNABLE            = 1
	BLOCKED             = 2
	WAITING             = 3
	TIMED_WAITING       = 4
	TERMINATED          = 5
)

const (
	MIN_PRIORITY  = 1
	NORM_PRIORITY = 5
	MAX_PRIORITY  = 10
)

/* functions that have been moved to gfunction/javaLangThread.go

func CreateMainThread() *object.Object {
	gl := globals.GetGlobalRef()
	globals.InitGlobals("test")
	main := object.StringObjectFromGoString("main")
	params := []any{main}
	t := gl.FuncInvokeGFunction("java/lang/Thread.<init>(Ljava/lang/String;)V",
		params)
	return t.(*object.Object)
}

// Adds a thread to the global thread table using the ID as the key,
// and a pointer to the thread itself as the value
func (t *ExecThread) AddThreadToTable(glob *globals.Globals) {
	glob.ThreadLock.Lock()
	glob.Threads[t.ID] = t
	glob.ThreadLock.Unlock()
}

// Runs a thread. This is the function that is called by the JVM when a thread is started.
// It calls the run method in the java/lang/Thread class and executes the Runnable that should be there.
func Run(t *object.Object) {
	glob := globals.GetGlobalRef()
	params := []any{t}
	glob.FuncInvokeGFunction("java/lang/Thread.run()V", params)
}

func RegisterThread(t *object.Object) {
	glob := globals.GetGlobalRef()
	ID := int(t.FieldTable["ID"].Fvalue.(int64))
	glob.ThreadLock.Lock()
	glob.Threads[ID] = t
	glob.ThreadLock.Unlock()
}

*/
