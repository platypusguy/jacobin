/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaNio

import (
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"testing"
)

func Test_FileVisitor(t *testing.T) {
	globals.InitGlobals("test")
	Load_Nio_File_FileVisitor()
	Load_Nio_File_FileVisitResult()

	fs := frames.CreateFrameStack()
	// CONTINUE is fvResultInstances[0]

	// params: [fs, this, arg...]
	p := []interface{}{fs, object.Null, object.Null, object.Null}

	if res := fileVisitorPreVisitDirectory(p); res != fvResultInstances[0] {
		t.Errorf("preVisitDirectory should return CONTINUE")
	}
	if res := fileVisitorVisitFile(p); res != fvResultInstances[0] {
		t.Errorf("visitFile should return CONTINUE")
	}
	if res := fileVisitorVisitFileFailed(p); res != fvResultInstances[0] {
		t.Errorf("visitFileFailed should return CONTINUE")
	}
	if res := fileVisitorPostVisitDirectory(p); res != fvResultInstances[0] {
		t.Errorf("postVisitDirectory should return CONTINUE")
	}
}

func Test_SimpleFileVisitor(t *testing.T) {
	globals.InitGlobals("test")
	Load_Nio_File_SimpleFileVisitor()
	Load_Nio_File_FileVisitResult()

	fs := frames.CreateFrameStack()
	p := []interface{}{fs, object.Null, object.Null, object.Null}

	if res := simpleFileVisitorPreVisitDirectory(p); res != fvResultInstances[0] {
		t.Errorf("preVisitDirectory should return CONTINUE")
	}
	if res := simpleFileVisitorVisitFile(p); res != fvResultInstances[0] {
		t.Errorf("visitFile should return CONTINUE")
	}

	// visitFileFailed: if exc is present, return it
	exc := object.MakeEmptyObjectWithClassName(new(string))
	pFailed := []interface{}{fs, object.Null, object.Null, exc}
	if res := simpleFileVisitorVisitFileFailed(pFailed); res != exc {
		t.Errorf("visitFileFailed should return exception if present")
	}
	if res := simpleFileVisitorVisitFileFailed(p); res != fvResultInstances[0] {
		t.Errorf("visitFileFailed should return CONTINUE if no exception")
	}

	// postVisitDirectory: if exc is present, return it
	if res := simpleFileVisitorPostVisitDirectory(pFailed); res != exc {
		t.Errorf("postVisitDirectory should return exception if present")
	}
	if res := simpleFileVisitorPostVisitDirectory(p); res != fvResultInstances[0] {
		t.Errorf("postVisitDirectory should return CONTINUE if no exception")
	}
}
