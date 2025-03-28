/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package jvm

import (
	"container/list"
	"errors"
	"fmt"
	"jacobin/classloader"
	"jacobin/frames"
	"jacobin/gfunction"
	"jacobin/globals"
	"jacobin/stringPool"
	"jacobin/trace"
	"jacobin/types"
)

// Initialization blocks are code blocks that for all intents are methods. They're gathered up by the
// Java compiler into a method called <clinit>, which must be run at class instantiation--that is,
// before any constructor. Because that code might well call other methods, it will need to be run
// just like a regular method with stack frames and depending on the interpreter in run.go
// In addition, we have to make sure that the initialization blocks of superclasses have been
// previously executed.
func runInitializationBlock(k *classloader.Klass, superClasses []string, fs *list.List) error {
	if superClasses == nil || len(superClasses) == 0 {
		// show we're running <clinit>. This prevents circularity errors.
		k.Data.ClInit = types.ClInitInProgress

		// if no superclasses were previously looked up
		// get list of the superclasses up to but not including java.lang.Object
		var superclasses []string

		// put the present class at the bottom of the list of superclasses,
		// because we'll need to run its clinit() code, if any
		superclasses = append(superclasses, k.Data.Name)

		superclass := *stringPool.GetStringPointer(k.Data.SuperclassIndex)
		for {
			if superclass == types.ObjectClassName {
				break
			}

			err := loadThisClass(superclass) // load the superclass
			if err != nil {                  // error message will have been displayed
				return err
			}

			// load only superclasses that have a clInit block that has not been run
			loadedSuperclass := classloader.MethAreaFetch(superclass)
			if loadedSuperclass.Data.ClInit == types.ClInitNotRun {
				superclasses = append(superclasses, superclass)
			}

			// now loop to see whether this superclass has a superclass
			superclass = *stringPool.GetStringPointer(loadedSuperclass.Data.SuperclassIndex)
		}
		superClasses = superclasses
	}

	// now execute any encountered <clinit> code in this class
	for i := len(superClasses) - 1; i >= 0; i-- {
		className := superClasses[i]
		me, err := classloader.FetchMethodAndCP(className, "<clinit>", "()V")
		if err == nil {
			switch me.MType {
			case 'J': // it's a Java initializer (the most common case)
				err = runJavaInitializer(me.Meth, k, fs)
			case 'G': // it's a golang implementation of the initializer
				err = runNativeInitializer(me, k, fs)
			}
			if err != nil {
				return err
			}
		} // if no <clinit> method, then skip that superclass
	}
	return nil
}

// Run the <clinit>() initializer code as a Java method. This effectively duplicates
// the code in run.go that creates a new frame and runs the method. Note that this
// code creates its own frame stack, which is distinct from the applications frame
// stack. The reason is that this is computing that's in most ways apart from the
// bytecode of the app. (This design might be revised at a later point and the two
// frame stacks combined into one.)
func runJavaInitializer(m classloader.MData, k *classloader.Klass, fs *list.List) error {
	meth := m.(classloader.JmEntry)
	f := frames.CreateFrame(meth.MaxStack + types.StackInflator) // Experimental expansion, see JACOBIN-494
	if fs.Front() != nil {
		parentFrame := *(fs.Front().Value.(*frames.Frame))
		f.Thread = parentFrame.Thread
	}
	f.MethName = "<clinit>"
	f.MethType = "()V"
	f.ClName = k.Data.Name
	f.CP = meth.Cp                        // add its pointer to the class CP
	f.Meth = append(f.Meth, meth.Code...) // copy the bytecodes over

	// allocate the local variables
	for j := 0; j < meth.MaxLocals; j++ {
		f.Locals = append(f.Locals, 0)
	}

	k.Data.ClInit = types.ClInitInProgress

	currJvmStackSize := fs.Len()
	if frames.PushFrame(fs, f) != nil {
		errMsg := "memory exception allocating frame in runJavaInitializer()"
		trace.Error(errMsg)
		return errors.New(errMsg)
	}

	if globals.TraceVerbose {
		infoMsg := fmt.Sprintf("Start init: class=%s, meth=%s%s, maxStack=%d, maxLocals=%d, code size=%d",
			f.ClName, f.MethName, f.MethType, meth.MaxStack, meth.MaxLocals, len(meth.Code))
		trace.Trace(infoMsg)
	}

	// the <clinit> method might call other methods, so we can't just determine that
	// it's completed by the return from interpret(). See JACOBIN-665
	for fs.Len() > currJvmStackSize { // loop until the frame stack is back to its pre-<clinit>() size
		interpret(fs) // if an error occurs, ThrowEx() will break us out of here
	}
	k.Data.ClInit = types.ClInitRun // flag showing we've run this class's <clinit>

	// frames.PopFrame(fs)
	return nil
}

func runNativeInitializer(mt classloader.MTentry, k *classloader.Klass, fs *list.List) error {
	_ = gfunction.RunGfunction(mt, fs, k.Data.Name, "<clinit>", "()V", nil, false, false)
	k.Data.ClInit = types.ClInitRun // flag showing we've run this class's <clinit>
	return nil
}
