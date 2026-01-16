/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaNio

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
)

// Load_Nio_File_Files registers implementations for java.nio.file.Files public methods (Java 21).
// Policy agreed with user:
// - Implement core filesystem ops.
// - Stream-returning and complex attribute-view methods remain trapped.
// - Where platform limitations apply, return UnsupportedOperationException.
func Load_Nio_File_Files() {
	ghelpers.MethodSignatures["java/nio/file/Files.<clinit>()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.ClinitGeneric}
	// Core boolean checks
	ghelpers.MethodSignatures["java/nio/file/Files.exists(Ljava/nio/file/Path;[Ljava/nio/file/LinkOption;)Z"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: filesExists}
	ghelpers.MethodSignatures["java/nio/file/Files.notExists(Ljava/nio/file/Path;[Ljava/nio/file/LinkOption;)Z"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: filesNotExists}
	ghelpers.MethodSignatures["java/nio/file/Files.isDirectory(Ljava/nio/file/Path;[Ljava/nio/file/LinkOption;)Z"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: filesIsDirectory}
	ghelpers.MethodSignatures["java/nio/file/Files.isRegularFile(Ljava/nio/file/Path;[Ljava/nio/file/LinkOption;)Z"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: filesIsRegularFile}
	ghelpers.MethodSignatures["java/nio/file/Files.isSameFile(Ljava/nio/file/Path;Ljava/nio/file/Path;)Z"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: filesIsSameFile}
	ghelpers.MethodSignatures["java/nio/file/Files.isSymbolicLink(Ljava/nio/file/Path;)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: filesIsSymbolicLink}

	// Size
	ghelpers.MethodSignatures["java/nio/file/Files.size(Ljava/nio/file/Path;)J"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: filesSize}

	// Create/Delete
	ghelpers.MethodSignatures["java/nio/file/Files.createFile(Ljava/nio/file/Path;[Ljava/nio/file/attribute/FileAttribute;)Ljava/nio/file/Path;"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: filesCreateFile}
	ghelpers.MethodSignatures["java/nio/file/Files.createDirectory(Ljava/nio/file/Path;[Ljava/nio/file/attribute/FileAttribute;)Ljava/nio/file/Path;"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: filesCreateDirectory}
	ghelpers.MethodSignatures["java/nio/file/Files.createDirectories(Ljava/nio/file/Path;[Ljava/nio/file/attribute/FileAttribute;)Ljava/nio/file/Path;"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: filesCreateDirectories}
	ghelpers.MethodSignatures["java/nio/file/Files.delete(Ljava/nio/file/Path;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: filesDelete}
	ghelpers.MethodSignatures["java/nio/file/Files.deleteIfExists(Ljava/nio/file/Path;)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: filesDeleteIfExists}

	// Copy/Move (Path, Path variants)
	ghelpers.MethodSignatures["java/nio/file/Files.copy(Ljava/nio/file/Path;Ljava/nio/file/Path;[Ljava/nio/file/CopyOption;)Ljava/nio/file/Path;"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: filesCopyPath}
	ghelpers.MethodSignatures["java/nio/file/Files.move(Ljava/nio/file/Path;Ljava/nio/file/Path;[Ljava/nio/file/CopyOption;)Ljava/nio/file/Path;"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: filesMove}

	// I/O streams
	ghelpers.MethodSignatures["java/nio/file/Files.newInputStream(Ljava/nio/file/Path;[Ljava/nio/file/OpenOption;)Ljava/io/InputStream;"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: filesNewInputStream}
	ghelpers.MethodSignatures["java/nio/file/Files.newOutputStream(Ljava/nio/file/Path;[Ljava/nio/file/OpenOption;)Ljava/io/OutputStream;"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: filesNewOutputStream}

	// Bulk read/write utilities
	ghelpers.MethodSignatures["java/nio/file/Files.readAllBytes(Ljava/nio/file/Path;)[B"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: filesReadAllBytes}
	ghelpers.MethodSignatures["java/nio/file/Files.write(Ljava/nio/file/Path;[B[Ljava/nio/file/OpenOption;)Ljava/nio/file/Path;"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: filesWriteBytes}
	ghelpers.MethodSignatures["java/nio/file/Files.readString(Ljava/nio/file/Path;)Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: filesReadString}
	ghelpers.MethodSignatures["java/nio/file/Files.writeString(Ljava/nio/file/Path;Ljava/lang/CharSequence;[Ljava/nio/file/OpenOption;)Ljava/nio/file/Path;"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: filesWriteString}
	ghelpers.MethodSignatures["java/nio/file/Files.readAllLines(Ljava/nio/file/Path;)Ljava/util/List;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: filesReadAllLines}

	// Links
	ghelpers.MethodSignatures["java/nio/file/Files.readSymbolicLink(Ljava/nio/file/Path;)Ljava/nio/file/Path;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: filesReadSymbolicLink}
	ghelpers.MethodSignatures["java/nio/file/Files.createSymbolicLink(Ljava/nio/file/Path;Ljava/nio/file/Path;[Ljava/nio/file/attribute/FileAttribute;)Ljava/nio/file/Path;"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: filesCreateSymbolicLink}
	ghelpers.MethodSignatures["java/nio/file/Files.createLink(Ljava/nio/file/Path;Ljava/nio/file/Path;)Ljava/nio/file/Path;"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: filesCreateLink}

	// Temp utilities
	ghelpers.MethodSignatures["java/nio/file/Files.createTempFile(Ljava/nio/file/Path;Ljava/lang/String;Ljava/lang/String;[Ljava/nio/file/attribute/FileAttribute;)Ljava/nio/file/Path;"] =
		ghelpers.GMeth{ParamSlots: 4, GFunction: filesCreateTempFileInDir}
	ghelpers.MethodSignatures["java/nio/file/Files.createTempFile(Ljava/lang/String;Ljava/lang/String;[Ljava/nio/file/attribute/FileAttribute;)Ljava/nio/file/Path;"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: filesCreateTempFile}
	ghelpers.MethodSignatures["java/nio/file/Files.createTempDirectory(Ljava/nio/file/Path;Ljava/lang/String;[Ljava/nio/file/attribute/FileAttribute;)Ljava/nio/file/Path;"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: filesCreateTempDirectoryInDir}
	ghelpers.MethodSignatures["java/nio/file/Files.createTempDirectory(Ljava/lang/String;[Ljava/nio/file/attribute/FileAttribute;)Ljava/nio/file/Path;"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: filesCreateTempDirectory}
}

// --- Helpers ---
func pathToGoString(p interface{}) (string, *ghelpers.GErrBlk) {
	if p == nil || p == object.Null {
		return "", ghelpers.GetGErrBlk(excNames.NullPointerException, "Path is null")
	}
	obj := p.(*object.Object)
	sval, ok := obj.FieldTable["value"].Fvalue.(*object.Object)
	if !ok || sval == nil {
		return "", ghelpers.GetGErrBlk(excNames.IOException, "Path.value missing")
	}
	return object.GoStringFromStringObject(sval), nil
}

func boolToJava(b bool) types.JavaBool {
	if b {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// Build a java.io.FileInputStream object pre-initialized with FilePath/FileHandle
func newFileInputStreamObj(path string, fh *os.File) *object.Object {
	className := "java/io/FileInputStream"
	obj := object.MakeEmptyObjectWithClassName(&className)
	obj.FieldTable[ghelpers.FilePath] = object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(path)}
	obj.FieldTable[ghelpers.FileHandle] = object.Field{Ftype: ghelpers.FileHandle, Fvalue: fh}
	return obj
}

// Build a java.io.FileOutputStream object pre-initialized with FilePath/FileHandle
func newFileOutputStreamObj(path string, fh *os.File) *object.Object {
	className := "java/io/FileOutputStream"
	obj := object.MakeEmptyObjectWithClassName(&className)
	obj.FieldTable[ghelpers.FilePath] = object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(path)}
	obj.FieldTable[ghelpers.FileHandle] = object.Field{Ftype: ghelpers.FileHandle, Fvalue: fh}
	return obj
}

// --- Implementations ---

func filesExists(params []interface{}) interface{} {
	p, gerr := pathToGoString(params[0])
	if gerr != nil {
		return gerr
	}
	_, err := os.Stat(p)
	if err == nil {
		return types.JavaBoolTrue
	}
	if os.IsNotExist(err) {
		return types.JavaBoolFalse
	}
	// Any other error -> false per JDK spec guidance
	return types.JavaBoolFalse
}

func filesNotExists(params []interface{}) interface{} {
	res := filesExists([]interface{}{params[0]})
	if geb, ok := res.(*ghelpers.GErrBlk); ok {
		return geb
	}
	return boolToJava(res.(types.JavaBool) == types.JavaBoolFalse)
}

func filesIsDirectory(params []interface{}) interface{} {
	p, gerr := pathToGoString(params[0])
	if gerr != nil {
		return gerr
	}
	fi, err := os.Stat(p)
	if err != nil {
		return types.JavaBoolFalse
	}
	return boolToJava(fi.IsDir())
}

func filesIsRegularFile(params []interface{}) interface{} {
	p, gerr := pathToGoString(params[0])
	if gerr != nil {
		return gerr
	}
	fi, err := os.Stat(p)
	if err != nil {
		return types.JavaBoolFalse
	}
	return boolToJava(fi.Mode().IsRegular())
}

func filesIsSameFile(params []interface{}) interface{} {
	p1, g1 := pathToGoString(params[0])
	if g1 != nil {
		return g1
	}
	p2, g2 := pathToGoString(params[1])
	if g2 != nil {
		return g2
	}
	// Resolve to absolute real paths best-effort
	a1, _ := filepath.Abs(p1)
	a2, _ := filepath.Abs(p2)
	return boolToJava(a1 == a2)
}

func filesIsSymbolicLink(params []interface{}) interface{} {
	p, gerr := pathToGoString(params[0])
	if gerr != nil {
		return gerr
	}
	fi, err := os.Lstat(p)
	if err != nil {
		return types.JavaBoolFalse
	}
	return boolToJava((fi.Mode() & fs.ModeSymlink) != 0)
}

func filesSize(params []interface{}) interface{} {
	p, gerr := pathToGoString(params[0])
	if gerr != nil {
		return gerr
	}
	fi, err := os.Stat(p)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.size: %s", err.Error()))
	}
	return int64(fi.Size())
}

func filesCreateFile(params []interface{}) interface{} {
	p, gerr := pathToGoString(params[0])
	if gerr != nil {
		return gerr
	}
	f, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0o666)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.createFile: %s", err.Error()))
	}
	_ = f.Close()
	return newPath(p)
}

func filesCreateDirectory(params []interface{}) interface{} {
	p, gerr := pathToGoString(params[0])
	if gerr != nil {
		return gerr
	}
	if err := os.Mkdir(p, 0o777); err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.createDirectory: %s", err.Error()))
	}
	return newPath(p)
}

func filesCreateDirectories(params []interface{}) interface{} {
	p, gerr := pathToGoString(params[0])
	if gerr != nil {
		return gerr
	}
	if err := os.MkdirAll(p, 0o777); err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.createDirectories: %s", err.Error()))
	}
	return newPath(p)
}

func filesDelete(params []interface{}) interface{} {
	p, gerr := pathToGoString(params[0])
	if gerr != nil {
		return gerr
	}
	if err := os.Remove(p); err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.delete: %s", err.Error()))
	}
	return nil
}

func filesDeleteIfExists(params []interface{}) interface{} {
	p, gerr := pathToGoString(params[0])
	if gerr != nil {
		return gerr
	}
	err := os.Remove(p)
	if err == nil {
		return types.JavaBoolTrue
	}
	if os.IsNotExist(err) {
		return types.JavaBoolFalse
	}
	return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.deleteIfExists: %s", err.Error()))
}

func filesCopyPath(params []interface{}) interface{} {
	src, g1 := pathToGoString(params[0])
	if g1 != nil {
		return g1
	}
	dst, g2 := pathToGoString(params[1])
	if g2 != nil {
		return g2
	}
	// Only support file-to-file regular copy for now.
	sfi, err := os.Stat(src)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.copy: %s", err.Error()))
	}
	if sfi.IsDir() {
		return ghelpers.GetGErrBlk(excNames.UnsupportedOperationException, "Files.copy: directory copy not supported")
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.copy: %s", err.Error()))
	}
	if err := os.WriteFile(dst, data, 0o666); err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.copy: %s", err.Error()))
	}
	return newPath(dst)
}

func filesMove(params []interface{}) interface{} {
	src, g1 := pathToGoString(params[0])
	if g1 != nil {
		return g1
	}
	dst, g2 := pathToGoString(params[1])
	if g2 != nil {
		return g2
	}
	if err := os.Rename(src, dst); err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.move: %s", err.Error()))
	}
	return newPath(dst)
}

func filesNewInputStream(params []interface{}) interface{} {
	p, gerr := pathToGoString(params[0])
	if gerr != nil {
		return gerr
	}
	fh, err := os.Open(p)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.newInputStream: %s", err.Error()))
	}
	return newFileInputStreamObj(p, fh)
}

func filesNewOutputStream(params []interface{}) interface{} {
	p, gerr := pathToGoString(params[0])
	if gerr != nil {
		return gerr
	}
	// Simplified: create/truncate
	fh, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o666)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.newOutputStream: %s", err.Error()))
	}
	return newFileOutputStreamObj(p, fh)
}

func filesReadAllBytes(params []interface{}) interface{} {
	p, gerr := pathToGoString(params[0])
	if gerr != nil {
		return gerr
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.readAllBytes: %s", err.Error()))
	}
	return object.JavaByteArrayFromGoByteArray(data)
}

func filesWriteBytes(params []interface{}) interface{} {
	p, gerr := pathToGoString(params[0])
	if gerr != nil {
		return gerr
	}
	// params[1] is a Java byte[] wrapped in object.Object? In this codebase, a Java byte[] is passed directly as []types.JavaByte
	jb, ok := params[1].([]types.JavaByte)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IOException, "Files.write: expected byte[] argument")
	}
	data := object.GoByteArrayFromJavaByteArray(jb)
	if err := os.WriteFile(p, data, 0o666); err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.write: %s", err.Error()))
	}
	return newPath(p)
}

func filesReadString(params []interface{}) interface{} {
	p, gerr := pathToGoString(params[0])
	if gerr != nil {
		return gerr
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.readString: %s", err.Error()))
	}
	return object.StringObjectFromGoString(string(data))
}

func filesWriteString(params []interface{}) interface{} {
	p, gerr := pathToGoString(params[0])
	if gerr != nil {
		return gerr
	}
	// CharSequence is a java.lang.String object in our use
	sObj := params[1].(*object.Object)
	text := object.GoStringFromStringObject(sObj)
	if err := os.WriteFile(p, []byte(text), 0o666); err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.writeString: %s", err.Error()))
	}
	return newPath(p)
}

func filesReadAllLines(params []interface{}) interface{} {
	p, gerr := pathToGoString(params[0])
	if gerr != nil {
		return gerr
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.readAllLines: %s", err.Error()))
	}
	// Split on \n, drop trailing \r if present (CRLF)
	var lines []string
	start := 0
	for i, b := range data {
		if b == '\n' {
			line := string(data[start:i])
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, line)
			start = i + 1
		}
	}
	if start <= len(data) {
		line := string(data[start:])
		if len(line) > 0 && len(data) > 0 && data[len(data)-1] != '\n' {
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, line)
		} else if start < len(data) {
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, line)
		}
	}
	// Represent as java.util.LinkedList of String[] as used elsewhere (see newStringIterator)
	arr := object.StringObjectArrayFromGoStringArray(lines)
	listObj := object.MakePrimitiveObject("java/util/LinkedList", types.Array+types.StringClassName, arr)
	return listObj
}

func filesReadSymbolicLink(params []interface{}) interface{} {
	p, gerr := pathToGoString(params[0])
	if gerr != nil {
		return gerr
	}
	target, err := os.Readlink(p)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.readSymbolicLink: %s", err.Error()))
	}
	return newPath(target)
}

func filesCreateSymbolicLink(params []interface{}) interface{} {
	link, g1 := pathToGoString(params[0])
	if g1 != nil {
		return g1
	}
	target, g2 := pathToGoString(params[1])
	if g2 != nil {
		return g2
	}
	if err := os.Symlink(target, link); err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.createSymbolicLink: %s", err.Error()))
	}
	return newPath(link)
}

func filesCreateLink(params []interface{}) interface{} {
	link, g1 := pathToGoString(params[0])
	if g1 != nil {
		return g1
	}
	existing, g2 := pathToGoString(params[1])
	if g2 != nil {
		return g2
	}
	if err := os.Link(existing, link); err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.createLink: %s", err.Error()))
	}
	return newPath(link)
}

func filesCreateTempFileInDir(params []interface{}) interface{} {
	dir, g1 := pathToGoString(params[0])
	if g1 != nil {
		return g1
	}
	prefix := object.GoStringFromStringObject(params[1].(*object.Object))
	suffix := object.GoStringFromStringObject(params[2].(*object.Object))
	pattern := prefix + "*" + suffix
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.createTempFile: %s", err.Error()))
	}
	_ = f.Close()
	return newPath(f.Name())
}

func filesCreateTempFile(params []interface{}) interface{} {
	prefix := object.GoStringFromStringObject(params[0].(*object.Object))
	suffix := object.GoStringFromStringObject(params[1].(*object.Object))
	pattern := prefix + "*" + suffix
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.createTempFile: %s", err.Error()))
	}
	_ = f.Close()
	return newPath(f.Name())
}

func filesCreateTempDirectoryInDir(params []interface{}) interface{} {
	dir, g1 := pathToGoString(params[0])
	if g1 != nil {
		return g1
	}
	prefix := object.GoStringFromStringObject(params[1].(*object.Object))
	path, err := os.MkdirTemp(dir, prefix+"*")
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.createTempDirectory: %s", err.Error()))
	}
	return newPath(path)
}

func filesCreateTempDirectory(params []interface{}) interface{} {
	prefix := object.GoStringFromStringObject(params[0].(*object.Object))
	path, err := os.MkdirTemp("", prefix+"*")
	if err != nil {
		return ghelpers.GetGErrBlk(excNames.IOException, fmt.Sprintf("Files.createTempDirectory: %s", err.Error()))
	}
	return newPath(path)
}
