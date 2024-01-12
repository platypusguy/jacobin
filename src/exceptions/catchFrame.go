/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package exceptions

import (
	"container/list"
	"jacobin/frames"
)

// This routine looks for a handler for the given exception (excName) in the
// current frame stack working its way up the frame stack (fs). If one is found,
// it returns a pointer to that frame, otherwise it returns nil. pc is the
// program counter in the current frame.
func FindCatchFrame(fs *list.List, excName string, pc int) *frames.Frame {
	// presentPC := pc
	f := fs.Front()
	for e := fs.Front(); e != nil; e = e.Next() {

		return f.Value.(*frames.Frame)
	}
	return nil
}
