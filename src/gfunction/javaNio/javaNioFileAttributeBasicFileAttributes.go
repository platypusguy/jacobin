/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaNio

import (
	"io/fs"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"time"
)

func Load_Nio_File_Attribute_BasicFileAttributes() {
	// isRegularFile()Z
	ghelpers.MethodSignatures["java/nio/file/attribute/BasicFileAttributes.isRegularFile()Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: bfaIsRegularFile}

	// isDirectory()Z
	ghelpers.MethodSignatures["java/nio/file/attribute/BasicFileAttributes.isDirectory()Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: bfaIsDirectory}

	// isSymbolicLink()Z
	ghelpers.MethodSignatures["java/nio/file/attribute/BasicFileAttributes.isSymbolicLink()Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: bfaIsSymbolicLink}

	// isOther()Z
	ghelpers.MethodSignatures["java/nio/file/attribute/BasicFileAttributes.isOther()Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: bfaIsOther}

	// size()J
	ghelpers.MethodSignatures["java/nio/file/attribute/BasicFileAttributes.size()J"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: bfaSize}

	// lastModifiedTime()Ljava/nio/file/attribute/FileTime;
	ghelpers.MethodSignatures["java/nio/file/attribute/BasicFileAttributes.lastModifiedTime()Ljava/nio/file/attribute/FileTime;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: bfaLastModifiedTime}
}

func bfaIsRegularFile(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	info := obj.FieldTable["info"].Fvalue.(fs.FileInfo)
	return boolToJava(info.Mode().IsRegular())
}

func bfaIsDirectory(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	info := obj.FieldTable["info"].Fvalue.(fs.FileInfo)
	return boolToJava(info.IsDir())
}

func bfaIsSymbolicLink(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	info := obj.FieldTable["info"].Fvalue.(fs.FileInfo)
	return boolToJava(info.Mode()&fs.ModeSymlink != 0)
}

func bfaIsOther(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	info := obj.FieldTable["info"].Fvalue.(fs.FileInfo)
	mode := info.Mode()
	return boolToJava(!mode.IsRegular() && !info.IsDir() && mode&fs.ModeSymlink == 0)
}

func bfaSize(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	info := obj.FieldTable["info"].Fvalue.(fs.FileInfo)
	return info.Size()
}

func bfaLastModifiedTime(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	info := obj.FieldTable["info"].Fvalue.(fs.FileInfo)
	return newFileTime(info.ModTime())
}

func newBasicFileAttributes(info fs.FileInfo) *object.Object {
	className := "java/nio/file/attribute/BasicFileAttributes"
	obj := object.MakeEmptyObjectWithClassName(&className)
	obj.FieldTable["info"] = object.Field{Ftype: "any", Fvalue: info}
	return obj
}

func newFileTime(t time.Time) *object.Object {
	className := "java/nio/file/attribute/FileTime"
	obj := object.MakeEmptyObjectWithClassName(&className)
	obj.FieldTable["value"] = object.Field{Ftype: types.Long, Fvalue: t.UnixMilli()}
	return obj
}

func Load_Nio_File_Attribute_FileTime() {
	// toMillis()J
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.toMillis()J"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: fileTimeToMillis}
}

func fileTimeToMillis(params []interface{}) interface{} {
	obj := params[0].(*object.Object)
	return obj.FieldTable["value"].Fvalue.(int64)
}
