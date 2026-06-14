package jvm

import (
	"container/list"
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/exceptions"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/shutdown"
	"jacobin/src/trace"
	"runtime/debug"
)

/*
Run an FQN Java function, called from a G function (JACOBIN-923)
================================================================
* Create a new frame and push it onto the current frame stack.
* In the new frame, set the class name, method name, and method type.
* Handle the method's parameters by pushing them in reverse order onto the frame's OpStack.
* Run the frame by calling interpret, passing the current frame stack.
* Return to the calling G function.
*/
func RunJavaFromG(fs *list.List, clName, methName, methType string, args ...any) {

	// Set up thread ID.
	prevF := fs.Front().Value.(*frames.Frame)
	threadID := prevF.Thread

	// Set up panic recovery.
	defer func() int {
		// only an untrapped panic gets us here
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			glob := globals.GetGlobalRef()
			glob.ErrorGoStack = stack
			exceptions.ShowPanicCause(r)
			exceptions.ShowFrameStack(threadID)
			exceptions.ShowGoStackTrace(nil)
			return shutdown.Exit(shutdown.APP_EXCEPTION)
		}
		return shutdown.OK
	}()

	// Set up the method and CPool for the method's class.
	mte, err := classloader.FetchMethodAndCP(clName, methName, methType)
	if err != nil {
		errMsg := fmt.Sprintf("RunJavaFromG: Could not find run method (%s.%s%s): %v", clName, methName, methType, err)
		exceptions.ThrowEx(excNames.NoSuchMethodError, errMsg, nil)
		return
	}

	// trace.Trace(fmt.Sprintf("DEBUG RunJavaFromG: Found run method: %s.%s%s", clName, methName, methType))

	// Get Mtable entry for the method.
	meth, ok := mte.Meth.(classloader.JmEntry)
	if !ok {
		errMsg := fmt.Sprintf("RunJavaFromG: Method found but it is not a Java method (MType: %c)", mte.MType)
		exceptions.ThrowEx(excNames.NoSuchMethodError, errMsg, nil)
		return
	}

	// Create the frame.
	f := frames.CreateFrame(meth.MaxStack)
	f.Thread = threadID
	f.ClName = clName
	f.MethName = methName
	f.MethType = methType
	f.AccessFlags = meth.AccessFlags
	f.Locals = make([]any, 0, meth.MaxLocals)

	// Add the Mtable pointer to the class CP.
	f.CP = meth.Cp
	// Copy the bytecodes into the frame.
	f.Meth = append(f.Meth, meth.Code...)

	// Populate the method's local variables for this frame with args.
	for k := 0; k < len(args); k++ {
		f.Locals = append(f.Locals, args[k])
	}

	// Allocate the remaining method's local variables.
	for k := len(args); k < meth.MaxLocals; k++ {
		f.Locals = append(f.Locals, int64(0))
	}

	// Push the frame onto the frame stack.
	if frames.PushFrame(fs, f) != nil {
		errMsg := fmt.Sprintf("RunJavaFromG: frames.PushFrame failed on thread: %d", f.Thread)
		exceptions.ThrowEx(excNames.OutOfMemoryError, errMsg, nil)
		return
	}

	if globals.TraceInst {
		traceInfo := fmt.Sprintf("RunJavaFromG: class=%s, meth=%s%s, maxStack=%d, maxLocals=%d, code size=%d",
			f.ClName, f.MethName, f.MethType, meth.MaxStack, meth.MaxLocals, len(meth.Code))
		trace.Trace(traceInfo)
	}

	// Execute the frame until it's popped.
	originalLen := fs.Len()
	for fs.Len() >= originalLen {
		interpret(fs)
	}

	// Return to the G function.

}
