/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package exceptions

import (
	"container/list"
	"fmt"
	"jacobin/classloader"
	"jacobin/frames"
	"jacobin/log"
)

// This routine looks for a handler for the given exception (excName) in the
// current frame stack working its way up the frame stack (fs). If one is found,
// it returns a pointer to that frame, otherwise it returns nil. Param pc is the
// program counter in the current frame where the execption was thrown.
func FindCatchFrame(fs *list.List, excName string, pc int) *frames.Frame {
	// presentPC := pc

	for fr := fs.Front(); fr != nil; fr = fr.Next() {
		var f = fr.Value.(*frames.Frame)

		// get the method and check for an exception catch table
		// get the full method nameclassloader.MTable = {map[string]classloader.MTentry}
		fullMethName := f.ClName + "." + f.MethName + f.MethType
		methEntry, found := classloader.MTable[fullMethName]
		if !found {
			errMsg := fmt.Sprintf("ATHROW: Method %s not found in MTable", fullMethName)
			_ = log.Log(errMsg, log.SEVERE)
			return nil
		}

		if methEntry.MType != 'J' {
			errMsg := fmt.Sprintf("ATHROW: Method %s is a native method", fullMethName)
			_ = log.Log(errMsg, log.SEVERE)
			return nil
		}

		method := methEntry.Meth.(classloader.JmEntry)
		if method.Exceptions == nil {
			errMsg := fmt.Sprintf("ATHROW: Method %s has no exception table", fullMethName)
			_ = log.Log(errMsg, log.INFO)
			continue // loop to the next frame
		}

		return f
	}
	return nil
}
