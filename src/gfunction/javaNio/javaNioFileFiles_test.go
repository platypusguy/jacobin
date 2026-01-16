/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaNio

import (
	"os"
	"path/filepath"
	"testing"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
)

func Test_Files_Exists_And_NotExists(t *testing.T) {
	// Ensure string pool and related globals are initialized for object/string creation
	globals.InitStringPool()
	dir := t.TempDir()
	f := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatalf("prep: %v", err)
	}
	p := newPath(f)

	if v := filesExists([]interface{}{p, object.Null}); v != types.JavaBoolTrue {
		t.Fatalf("exists should be true, got %v", v)
	}
	if v := filesNotExists([]interface{}{p, object.Null}); v != types.JavaBoolFalse {
		t.Fatalf("notExists should be false, got %v", v)
	}

	pn := newPath(filepath.Join(dir, "nope.txt"))
	if v := filesExists([]interface{}{pn, object.Null}); v != types.JavaBoolFalse {
		t.Fatalf("exists (no file) should be false, got %v", v)
	}
	if v := filesNotExists([]interface{}{pn, object.Null}); v != types.JavaBoolTrue {
		t.Fatalf("notExists (no file) should be true, got %v", v)
	}
}

func Test_Files_IsDirectory_IsRegularFile(t *testing.T) {
	dir := t.TempDir()
	dpath := newPath(dir)
	f := filepath.Join(dir, "b.bin")
	if err := os.WriteFile(f, []byte{1, 2, 3}, 0o644); err != nil {
		t.Fatalf("prep: %v", err)
	}
	fpath := newPath(f)

	if filesIsDirectory([]interface{}{dpath, object.Null}) != types.JavaBoolTrue {
		t.Fatalf("dir should be directory")
	}
	if filesIsRegularFile([]interface{}{dpath, object.Null}) != types.JavaBoolFalse {
		t.Fatalf("dir should not be regular file")
	}
	if filesIsRegularFile([]interface{}{fpath, object.Null}) != types.JavaBoolTrue {
		t.Fatalf("file should be regular file")
	}
	if filesIsDirectory([]interface{}{fpath, object.Null}) != types.JavaBoolFalse {
		t.Fatalf("file should not be directory")
	}
}

func Test_Files_Size(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "c.txt")
	data := []byte("hello")
	if err := os.WriteFile(f, data, 0o644); err != nil {
		t.Fatalf("prep: %v", err)
	}
	v := filesSize([]interface{}{newPath(f)})
	if n, ok := v.(int64); !ok || n != int64(len(data)) {
		t.Fatalf("size got %T %v, want %d", v, v, len(data))
	}
	// error path
	v2 := filesSize([]interface{}{newPath(filepath.Join(dir, "none"))})
	if _, ok := v2.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected error for non-existent size, got %T", v2)
	}
}

func Test_Files_CreateFile_Directory_Delete_DeleteIfExists(t *testing.T) {
	dir := t.TempDir()
	f := newPath(filepath.Join(dir, "d.txt"))
	res := filesCreateFile([]interface{}{f, object.Null})
	if _, ok := res.(*object.Object); !ok {
		t.Fatalf("createFile did not return Path, got %T", res)
	}
	// second create -> error
	res2 := filesCreateFile([]interface{}{f, object.Null})
	if _, ok := res2.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected error on duplicate createFile, got %T", res2)
	}

	// createDirectory
	d := newPath(filepath.Join(dir, "adir"))
	dr := filesCreateDirectory([]interface{}{d, object.Null})
	if _, ok := dr.(*object.Object); !ok {
		t.Fatalf("createDirectory did not return Path")
	}
	// duplicate -> error
	dr2 := filesCreateDirectory([]interface{}{d, object.Null})
	if _, ok := dr2.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected error on duplicate createDirectory")
	}

	// deleteIfExists existing
	if v := filesDeleteIfExists([]interface{}{f}); v != types.JavaBoolTrue {
		t.Fatalf("deleteIfExists should return true on existing file")
	}
	// deleteIfExists non-existing
	if v := filesDeleteIfExists([]interface{}{f}); v != types.JavaBoolFalse {
		t.Fatalf("deleteIfExists should return false on non-existing file")
	}
	// delete non-existing -> error
	del := filesDelete([]interface{}{f})
	if _, ok := del.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected error for delete non-existing, got %T", del)
	}
}

func Test_Files_Copy_And_Move(t *testing.T) {
	dir := t.TempDir()
	s := filepath.Join(dir, "src.txt")
	if err := os.WriteFile(s, []byte("abc"), 0o644); err != nil {
		t.Fatalf("prep: %v", err)
	}
	d := filepath.Join(dir, "dst.txt")
	r := filesCopyPath([]interface{}{newPath(s), newPath(d), object.Null})
	if _, ok := r.(*object.Object); !ok {
		t.Fatalf("copy should return Path, got %T", r)
	}
	b, _ := os.ReadFile(d)
	if string(b) != "abc" {
		t.Fatalf("copy content mismatch: %q", string(b))
	}

	// copy directory -> unsupported
	adir := filepath.Join(dir, "dd")
	if err := os.Mkdir(adir, 0o755); err != nil {
		t.Fatalf("prep: %v", err)
	}
	r2 := filesCopyPath([]interface{}{newPath(adir), newPath(filepath.Join(dir, "dd2")), object.Null})
	if geb, ok := r2.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.UnsupportedOperationException {
		t.Fatalf("expected UnsupportedOperationException for dir copy, got %T %+v", r2, r2)
	}

	// move
	m := filepath.Join(dir, "moved.txt")
	r3 := filesMove([]interface{}{newPath(d), newPath(m), object.Null})
	if _, ok := r3.(*object.Object); !ok {
		t.Fatalf("move should return Path")
	}
	if _, err := os.Stat(d); !os.IsNotExist(err) {
		t.Fatalf("old file should not exist after move")
	}
}

func Test_Files_NewInputStream_NewOutputStream(t *testing.T) {
	dir := t.TempDir()
	// InputStream error path (no file)
	bad := filesNewInputStream([]interface{}{newPath(filepath.Join(dir, "nope")), object.Null})
	if _, ok := bad.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected error for missing file input stream")
	}

	// OutputStream success; we will write using the handle field
	op := filesNewOutputStream([]interface{}{newPath(filepath.Join(dir, "out.txt")), object.Null})
	outObj, ok := op.(*object.Object)
	if !ok {
		t.Fatalf("expected FileOutputStream object, got %T", op)
	}
	fld := outObj.FieldTable[ghelpers.FileHandle]
	fh, ok := fld.Fvalue.(*os.File)
	if !ok {
		t.Fatalf("missing FileHandle in output stream object")
	}
	if _, err := fh.Write([]byte("Q")); err != nil {
		t.Fatalf("write via handle: %v", err)
	}
	_ = fh.Close()

	b, err := os.ReadFile(filepath.Join(dir, "out.txt"))
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if string(b) != "Q" {
		t.Fatalf("content mismatch: %q", string(b))
	}
}

func Test_Files_ReadAllBytes_WriteBytes(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "rw.bin")
	jb := object.JavaByteArrayFromGoByteArray([]byte{9, 8, 7})
	wr := filesWriteBytes([]interface{}{newPath(f), jb, object.Null})
	if _, ok := wr.(*object.Object); !ok {
		t.Fatalf("write should return Path")
	}
	rd := filesReadAllBytes([]interface{}{newPath(f)})
	arr, ok := rd.([]types.JavaByte)
	if !ok {
		t.Fatalf("readAllBytes did not return byte[]: %T", rd)
	}
	gb := object.GoByteArrayFromJavaByteArray(arr)
	if len(gb) != 3 || gb[0] != 9 || gb[1] != 8 || gb[2] != 7 {
		t.Fatalf("bytes mismatch: %v", gb)
	}

	// error path
	rd2 := filesReadAllBytes([]interface{}{newPath(filepath.Join(dir, "nope"))})
	if _, ok := rd2.(*ghelpers.GErrBlk); !ok {
		t.Fatalf("expected error reading missing file")
	}
}

func Test_Files_ReadString_WriteString_ReadAllLines(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "rw.txt")
	s := object.StringObjectFromGoString("line1\nline2")
	r := filesWriteString([]interface{}{newPath(f), s, object.Null})
	if _, ok := r.(*object.Object); !ok {
		t.Fatalf("writeString should return Path")
	}

	rds := filesReadString([]interface{}{newPath(f)})
	so, ok := rds.(*object.Object)
	if !ok {
		t.Fatalf("readString did not return String object: %T", rds)
	}
	if txt := object.GoStringFromStringObject(so); txt != "line1\nline2" {
		t.Fatalf("readString mismatch: %q", txt)
	}

	lst := filesReadAllLines([]interface{}{newPath(f)})
	if _, ok := lst.(*object.Object); !ok {
		t.Fatalf("readAllLines should return a List object")
	}
}

func Test_Files_IsSameFile_And_Temps(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "x.txt")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatalf("prep: %v", err)
	}
	p1 := newPath(f)
	p2 := newPath(f)
	same := filesIsSameFile([]interface{}{p1, p2})
	if same != types.JavaBoolTrue {
		t.Fatalf("same file should be true")
	}
	other := filesIsSameFile([]interface{}{p1, newPath(filepath.Join(dir, "y.txt"))})
	if other != types.JavaBoolFalse {
		t.Fatalf("different file should be false")
	}

	// temp file/dir
	tf := filesCreateTempFile([]interface{}{object.StringObjectFromGoString("pre"), object.StringObjectFromGoString(".suf"), object.Null})
	if _, ok := tf.(*object.Object); !ok {
		t.Fatalf("createTempFile should return Path")
	}
	td := filesCreateTempDirectory([]interface{}{object.StringObjectFromGoString("pfx"), object.Null})
	if _, ok := td.(*object.Object); !ok {
		t.Fatalf("createTempDirectory should return Path")
	}
}

func Test_Files_Symlink_Paths(t *testing.T) {
	dir := t.TempDir()
	tgt := filepath.Join(dir, "t.txt")
	if err := os.WriteFile(tgt, []byte("z"), 0o644); err != nil {
		t.Fatalf("prep: %v", err)
	}
	link := filepath.Join(dir, "lnk")

	// createSymbolicLink: may fail on platform without privileges; both paths are valid executable branches
	r := filesCreateSymbolicLink([]interface{}{newPath(link), newPath(tgt), object.Null})
	if _, ok := r.(*ghelpers.GErrBlk); ok {
		// error path covered; now readSymbolicLink should also error for non-link
		rl := filesReadSymbolicLink([]interface{}{newPath(tgt)})
		if _, ok := rl.(*ghelpers.GErrBlk); !ok {
			t.Fatalf("expected error reading non-link")
		}
		// isSymbolicLink on regular file false
		if filesIsSymbolicLink([]interface{}{newPath(tgt)}) != types.JavaBoolFalse {
			t.Fatalf("regular file is not symlink")
		}
		return
	}
	// success path: link exists
	if filesIsSymbolicLink([]interface{}{newPath(link)}) != types.JavaBoolTrue {
		t.Fatalf("link should be symlink")
	}
	rl := filesReadSymbolicLink([]interface{}{newPath(link)})
	if _, ok := rl.(*object.Object); !ok {
		t.Fatalf("readSymbolicLink should return Path")
	}
}
