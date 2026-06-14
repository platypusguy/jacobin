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

func Test_SimpleFileVisitorDetailed(t *testing.T) {
	globals.InitGlobals("test")
	Load_Nio_File_SimpleFileVisitor()
	Load_Nio_File_FileVisitResult()

	fs := frames.CreateFrameStack()

	t.Run("preVisitDirectory", func(t *testing.T) {
		p := []interface{}{fs, object.Null, object.Null, object.Null}
		res := simpleFileVisitorPreVisitDirectory(p)
		if res != fvResultInstances[0] {
			t.Errorf("expected CONTINUE, got %v", res)
		}
	})

	t.Run("visitFile", func(t *testing.T) {
		p := []interface{}{fs, object.Null, object.Null, object.Null}
		res := simpleFileVisitorVisitFile(p)
		if res != fvResultInstances[0] {
			t.Errorf("expected CONTINUE, got %v", res)
		}
	})

	t.Run("visitFileFailed", func(t *testing.T) {
		// Case 1: No exception
		p1 := []interface{}{fs, object.Null, object.Null, object.Null}
		res1 := simpleFileVisitorVisitFileFailed(p1)
		if res1 != fvResultInstances[0] {
			t.Errorf("expected CONTINUE when no exception, got %v", res1)
		}

		// Case 2: Exception present
		exc := object.MakeEmptyObjectWithClassName(new(string))
		p2 := []interface{}{fs, object.Null, object.Null, exc}
		res2 := simpleFileVisitorVisitFileFailed(p2)
		if res2 != exc {
			t.Errorf("expected exception to be returned, got %v", res2)
		}

		// Case 3: Too few params
		p3 := []interface{}{fs, object.Null, object.Null}
		res3 := simpleFileVisitorVisitFileFailed(p3)
		if res3 != fvResultInstances[0] {
			t.Errorf("expected CONTINUE with too few params, got %v", res3)
		}
	})

	t.Run("postVisitDirectory", func(t *testing.T) {
		// Case 1: No exception
		p1 := []interface{}{fs, object.Null, object.Null, object.Null}
		res1 := simpleFileVisitorPostVisitDirectory(p1)
		if res1 != fvResultInstances[0] {
			t.Errorf("expected CONTINUE when no exception, got %v", res1)
		}

		// Case 2: Exception present
		exc := object.MakeEmptyObjectWithClassName(new(string))
		p2 := []interface{}{fs, object.Null, object.Null, exc}
		res2 := simpleFileVisitorPostVisitDirectory(p2)
		if res2 != exc {
			t.Errorf("expected exception to be returned, got %v", res2)
		}

		// Case 3: Too few params
		p3 := []interface{}{fs, object.Null, object.Null}
		res3 := simpleFileVisitorPostVisitDirectory(p3)
		if res3 != fvResultInstances[0] {
			t.Errorf("expected CONTINUE with too few params, got %v", res3)
		}
	})
}
