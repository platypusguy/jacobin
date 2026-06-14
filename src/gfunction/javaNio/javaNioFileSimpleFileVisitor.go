/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaNio

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
)

func Load_Nio_File_SimpleFileVisitor() {
	// preVisitDirectory(Object dir, BasicFileAttributes attrs)
	ghelpers.MethodSignatures["java/nio/file/SimpleFileVisitor.preVisitDirectory(Ljava/lang/Object;Ljava/nio/file/attribute/BasicFileAttributes;)Ljava/nio/file/FileVisitResult;"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: simpleFileVisitorPreVisitDirectory, NeedsContext: true}

	// visitFile(Object file, BasicFileAttributes attrs)
	ghelpers.MethodSignatures["java/nio/file/SimpleFileVisitor.visitFile(Ljava/lang/Object;Ljava/nio/file/attribute/BasicFileAttributes;)Ljava/nio/file/FileVisitResult;"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: simpleFileVisitorVisitFile, NeedsContext: true}

	// visitFileFailed(Object file, IOException exc)
	ghelpers.MethodSignatures["java/nio/file/SimpleFileVisitor.visitFileFailed(Ljava/lang/Object;Ljava/io/IOException;)Ljava/nio/file/FileVisitResult;"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: simpleFileVisitorVisitFileFailed, NeedsContext: true}

	// postVisitDirectory(Object dir, IOException exc)
	ghelpers.MethodSignatures["java/nio/file/SimpleFileVisitor.postVisitDirectory(Ljava/lang/Object;Ljava/io/IOException;)Ljava/nio/file/FileVisitResult;"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: simpleFileVisitorPostVisitDirectory, NeedsContext: true}
}

func simpleFileVisitorPreVisitDirectory(params []interface{}) interface{} {
	ensureFvResultInited()
	// Default behavior: CONTINUE
	return fvResultInstances[0] // CONTINUE
}

func simpleFileVisitorVisitFile(params []interface{}) interface{} {
	ensureFvResultInited()
	// Default behavior: CONTINUE
	return fvResultInstances[0] // CONTINUE
}

func simpleFileVisitorVisitFileFailed(params []interface{}) interface{} {
	ensureFvResultInited()
	// params[0] is frame stack if NeedsContext is true
	// params[1] is 'this', params[2] is 'file', params[3] is 'exc'
	if len(params) >= 4 {
		exc := params[3]
		if exc != nil && !object.IsNull(exc) {
			// SimpleFileVisitor.visitFileFailed throws the exception it's passed
			return exc
		}
	}
	return fvResultInstances[0] // CONTINUE
}

func simpleFileVisitorPostVisitDirectory(params []interface{}) interface{} {
	ensureFvResultInited()
	// params[0] is frame stack if NeedsContext is true
	// params[1] is 'this', params[2] is 'dir', params[3] is 'exc'
	if len(params) >= 4 {
		exc := params[3]
		if exc != nil && !object.IsNull(exc) {
			// SimpleFileVisitor.postVisitDirectory throws the exception it's passed
			return exc
		}
	}
	return fvResultInstances[0] // CONTINUE
}
