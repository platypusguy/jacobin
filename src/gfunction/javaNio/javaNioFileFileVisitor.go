/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaNio

import (
	"jacobin/src/gfunction/ghelpers"
)

func Load_Nio_File_FileVisitor() {
	// preVisitDirectory(Object dir, BasicFileAttributes attrs)
	ghelpers.MethodSignatures["java/nio/file/FileVisitor.preVisitDirectory(Ljava/lang/Object;Ljava/nio/file/attribute/BasicFileAttributes;)Ljava/nio/file/FileVisitResult;"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: fileVisitorPreVisitDirectory, NeedsContext: true}

	// visitFile(Object file, BasicFileAttributes attrs)
	ghelpers.MethodSignatures["java/nio/file/FileVisitor.visitFile(Ljava/lang/Object;Ljava/nio/file/attribute/BasicFileAttributes;)Ljava/nio/file/FileVisitResult;"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: fileVisitorVisitFile, NeedsContext: true}

	// visitFileFailed(Object file, IOException exc)
	ghelpers.MethodSignatures["java/nio/file/FileVisitor.visitFileFailed(Ljava/lang/Object;Ljava/io/IOException;)Ljava/nio/file/FileVisitResult;"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: fileVisitorVisitFileFailed, NeedsContext: true}

	// postVisitDirectory(Object dir, IOException exc)
	ghelpers.MethodSignatures["java/nio/file/FileVisitor.postVisitDirectory(Ljava/lang/Object;Ljava/io/IOException;)Ljava/nio/file/FileVisitResult;"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: fileVisitorPostVisitDirectory, NeedsContext: true}
}

func fileVisitorPreVisitDirectory(params []interface{}) interface{} {
	ensureFvResultInited()
	// Default behavior: CONTINUE
	return fvResultInstances[0] // CONTINUE
}

func fileVisitorVisitFile(params []interface{}) interface{} {
	ensureFvResultInited()
	// Default behavior: CONTINUE
	return fvResultInstances[0] // CONTINUE
}

func fileVisitorVisitFileFailed(params []interface{}) interface{} {
	ensureFvResultInited()
	// Default behavior: CONTINUE
	return fvResultInstances[0] // CONTINUE
}

func fileVisitorPostVisitDirectory(params []interface{}) interface{} {
	ensureFvResultInited()
	// Default behavior: CONTINUE
	return fvResultInstances[0] // CONTINUE
}
