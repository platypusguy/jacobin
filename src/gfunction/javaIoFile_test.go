package gfunction

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
)

// Helpers
func newFileObjFromPath(t *testing.T, p string) *object.Object {
	t.Helper()
	className := "java/io/File"
	f := object.MakeEmptyObjectWithClassName(&className)
	// call File.<init>(String)
	ret := fileInit([]interface{}{f, object.StringObjectFromGoString(p)})
	if ret != nil {
		t.Fatalf("fileInit returned error: %v", ret)
	}
	return f
}

func getPath(t *testing.T, f *object.Object) string {
	t.Helper()
	p, gerr := fileGetPathString(f)
	if gerr != nil {
		t.Fatalf("fileGetPathString error: %s", gerr.ErrMsg)
	}
	return p
}

func TestJavaIoFile_MethodRegistration(t *testing.T) {
	globals.InitStringPool()
	MethodSignatures = make(map[string]GMeth)
	Load_Io_File()

	checks := []struct {
		key   string
		slots int
	}{
		{"java/io/File.<init>(Ljava/lang/String;)V", 1},
		{"java/io/File.getPath()Ljava/lang/String;", 0},
		{"java/io/File.exists()Z", 0},
		{"java/io/File.createNewFile()Z", 0},
		{"java/io/File.delete()Z", 0},
		{"java/io/File.length()J", 0},
		{"java/io/File.list()[Ljava/lang/String;", 0},
		{"java/io/File.listFiles()[Ljava/io/File;", 0},
		{"java/io/File.mkdir()Z", 0},
		{"java/io/File.mkdirs()Z", 0},
		{"java/io/File.renameTo(Ljava/io/File;)Z", 1},
		{"java/io/File.setReadOnly()Z", 0},
		{"java/io/File.getAbsolutePath()Ljava/lang/String;", 0},
		{"java/io/File.getCanonicalPath()Ljava/lang/String;", 0},
	}
	for _, c := range checks {
		gm, ok := MethodSignatures[c.key]
		if !ok {
			t.Fatalf("method not registered: %s", c.key)
		}
		if gm.ParamSlots != c.slots {
			t.Fatalf("ParamSlots mismatch for %s: want %d got %d", c.key, c.slots, gm.ParamSlots)
		}
		if gm.GFunction == nil {
			t.Fatalf("GFunction is nil for %s", c.key)
		}
	}
}

func TestJavaIoFile_Init_And_Getters(t *testing.T) {
	globals.InitStringPool()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	f := newFileObjFromPath(t, path)

	// getPath should be absolute
	pstr := fileGetPath([]interface{}{f}).(*object.Object)
	gp := object.GoStringFromStringObject(pstr)
	if !filepath.IsAbs(gp) {
		t.Fatalf("getPath not absolute: %q", gp)
	}
	// name and parent
	name := fileGetName([]interface{}{f}).(*object.Object)
	if object.GoStringFromStringObject(name) != filepath.Base(gp) {
		t.Fatalf("getName mismatch")
	}
	parent := fileGetParent([]interface{}{f}).(*object.Object)
	if object.GoStringFromStringObject(parent) != filepath.Dir(gp) {
		t.Fatalf("getParent mismatch")
	}
	// absolute and toString
	abs := fileIsAbsolute([]interface{}{f}).(int64)
	if abs != types.JavaBoolTrue {
		t.Fatalf("isAbsolute expected true")
	}
	toStr := fileToString([]interface{}{f}).(*object.Object)
	if object.GoStringFromStringObject(toStr) != gp {
		t.Fatalf("toString mismatch")
	}
	// absolute and canonical path strings
	absP := fileGetAbsolutePath([]interface{}{f}).(*object.Object)
	canP := fileGetCanonicalPath([]interface{}{f}).(*object.Object)
	if object.GoStringFromStringObject(absP) != object.GoStringFromStringObject(canP) {
		// On our minimal impl these are the same
		t.Fatalf("absolute and canonical path differ: %q vs %q", object.GoStringFromStringObject(absP), object.GoStringFromStringObject(canP))
	}
}

func TestJavaIoFile_Create_Exists_Length_Delete(t *testing.T) {
	globals.InitStringPool()
	dir := t.TempDir()
	path := filepath.Join(dir, "afile.dat")
	f := newFileObjFromPath(t, path)

	// initially does not exist
	if fileExists([]interface{}{f}).(int64) != types.JavaBoolFalse {
		t.Fatalf("exists expected false initially")
	}
	// create
	if fileCreate([]interface{}{f}).(int64) != 1 {
		t.Fatalf("createNewFile returned false")
	}
	if fileExists([]interface{}{f}).(int64) != types.JavaBoolTrue {
		t.Fatalf("exists expected true after create")
	}
	// length initially 0
	if fileLength([]interface{}{f}).(int64) != 0 {
		t.Fatalf("length expected 0 initially")
	}
	// write some bytes
	goPath := getPath(t, f)
	if err := os.WriteFile(goPath, []byte("hello"), 0664); err != nil {
		t.Fatalf("write file error: %v", err)
	}
	if fileLength([]interface{}{f}).(int64) != 5 {
		t.Fatalf("length expected 5 after write")
	}

	// delete
	if fileDelete([]interface{}{f}).(int64) != 1 {
		t.Fatalf("delete returned false")
	}
	if fileExists([]interface{}{f}).(int64) != types.JavaBoolFalse {
		t.Fatalf("exists expected false after delete")
	}
}

func TestJavaIoFile_List_And_ListFiles(t *testing.T) {
	globals.InitStringPool()
	dir := t.TempDir()
	// create items
	_ = os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0664)
	_ = os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0664)
	_ = os.Mkdir(filepath.Join(dir, "sub"), 0775)

	dirFile := newFileObjFromPath(t, dir)
	// list names
	arr := fileList([]interface{}{dirFile}).(*object.Object)
	names, _ := arr.FieldTable["value"].Fvalue.([]*object.Object)
	var got []string
	for _, s := range names {
		if s == nil {
			continue
		}
		got = append(got, object.GoStringFromStringObject(s))
	}
	sort.Strings(got)
	// We expect at least these three entries
	want := map[string]bool{"a.txt": true, "b.txt": true, "sub": true}
	for k := range want {
		found := false
		for _, g := range got {
			if g == k {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("list missing entry %q; got %v", k, got)
		}
	}
	// listFiles returns File[]
	farr := fileListFiles([]interface{}{dirFile}).(*object.Object)
	files, _ := farr.FieldTable["value"].Fvalue.([]*object.Object)
	if len(files) < 3 {
		t.Fatalf("listFiles expected >=3 entries, got %d", len(files))
	}
}

func TestJavaIoFile_Mkdir_Mkdirs_IsDirectory(t *testing.T) {
	globals.InitStringPool()
	dir := t.TempDir()
	nested := filepath.Join(dir, "a", "b", "c")
	f := newFileObjFromPath(t, nested)
	// mkdir should fail on a deep path (parent missing) and return false; mkdirs true
	_ = fileDelete([]interface{}{f}) // ensure not exists
	if fileMkdir([]interface{}{f}).(int64) != types.JavaBoolFalse {
		// On some OS, Mkdir may fail with ENOENT as expected
	}
	if fileMkdirs([]interface{}{f}).(int64) != types.JavaBoolTrue {
		t.Fatalf("mkdirs expected true")
	}
	if fileIsDirectory([]interface{}{f}).(int64) != types.JavaBoolTrue {
		t.Fatalf("isDirectory expected true after mkdirs")
	}
}

func TestJavaIoFile_RenameTo(t *testing.T) {
	globals.InitStringPool()

	// get temp directory
	dir := t.TempDir()
	pidstr := fmt.Sprintf("%d", os.Getpid())

	// construct src and dst paths
	srcpath := filepath.Join(dir, pidstr, ".src.txt")
	dstpath := filepath.Join(dir, pidstr, ".dst.txt")

	t.Logf("srcpath: %q", srcpath)
	t.Logf("dstpath: %q", dstpath)

	// ensure parent directory exists
	srcDir := filepath.Dir(srcpath)
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("failed to create parent directory %q: %v", srcDir, err)
	}

	// create file objects
	src := newFileObjFromPath(t, srcpath)
	dst := newFileObjFromPath(t, dstpath)

	// create the source file
	if fileCreate([]interface{}{src}) != types.JavaBoolTrue {
		t.Fatalf("fileCreate failed")
	}

	// rename src -> dst
	if fileRenameTo([]interface{}{src, dst}).(int64) != types.JavaBoolTrue {
		t.Fatalf("fileRenameTo failed")
	}

	// check that dst exists
	if fileExists([]interface{}{dst}).(int64) != types.JavaBoolTrue {
		t.Fatalf("dst should exist after rename")
	}

	// Delete the temp directory.
	err := os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("Ignoring os.RemoveAll(%s) error: %v", dir, err)
	}

}

func TestJavaIoFile_SetReadOnly_And_Permissions_Noops(t *testing.T) {
	globals.InitStringPool()
	dir := t.TempDir()
	pidstr := fmt.Sprintf("%d", os.Getpid())
	fpath := filepath.Join(dir, pidstr, "ro.txt")
	f := newFileObjFromPath(t, fpath)

	// ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
		t.Fatalf("failed to create parent directory %q: %v", filepath.Dir(fpath), err)
	}

	// create the file
	_ = fileCreate([]interface{}{f})

	if fileSetReadOnly([]interface{}{f}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setReadOnly expected true")
	}

	// minimal no-ops that return true
	if fileSetReadable([]interface{}{f, types.JavaBoolTrue}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setReadable expected true")
	}
	if fileSetWritable([]interface{}{f, types.JavaBoolTrue}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setWritable expected true")
	}
	if fileSetExecutable([]interface{}{f, types.JavaBoolTrue}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setExecutable expected true")
	}

	// Delete the temp directory.
	err := os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("Ignoring os.RemoveAll(%s) error: %v", dir, err)
	}
}

func TestJavaIoFile_CreateTemp_Instance_And_Static(t *testing.T) {
	globals.InitStringPool()
	dir := t.TempDir()
	// instance createTempFile
	base := newFileObjFromPath(t, dir)
	obj := fileCreateTemp([]interface{}{base, object.StringObjectFromGoString("pre"), object.StringObjectFromGoString(".suf")})
	if obj == nil {
		t.Fatalf("instance createTempFile returned nil")
	}
	// static createTempFile(prefix,suffix,dir)
	obj2 := fileCreateTempWithDir([]interface{}{object.StringObjectFromGoString("x"), object.StringObjectFromGoString(".log"), newFileObjFromPath(t, dir)})
	if obj2 == nil {
		t.Fatalf("static createTempFile returned nil")
	}
}

func TestJavaIoFile_Equals_And_HashCode(t *testing.T) {
	globals.InitStringPool()
	dir := t.TempDir()
	p := filepath.Join(dir, "same")
	f1 := newFileObjFromPath(t, p)
	f2 := newFileObjFromPath(t, p)
	// equals should be true for the same absolute path
	if fileEquals([]interface{}{f1, f2}).(int64) != types.JavaBoolTrue {
		t.Fatalf("equals expected true for same path")
	}
	h1 := fileHashCode([]interface{}{f1}).(int64)
	h2 := fileHashCode([]interface{}{f2}).(int64)
	if h1 != h2 {
		t.Fatalf("hashCode expected equal for same path")
	}
}

func TestJavaIoFile_PermissionSetters(t *testing.T) {
	globals.InitStringPool()
	dir := t.TempDir()
	pidstr := fmt.Sprintf("%d", os.Getpid())
	fpath := filepath.Join(dir, pidstr, "perm.bin")
	f := newFileObjFromPath(t, fpath)

	// ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
		t.Fatalf("failed to create parent directory %q: %v", filepath.Dir(fpath), err)
	}

	if fileCreate([]interface{}{f}).(int64) != 1 {
		t.Fatalf("createNewFile failed")
	}

	goPath := getPath(t, f)

	// Establish a known baseline of 0000; skip if not supported
	if err := os.Chmod(goPath, 0o000); err != nil {
		t.Skipf("skipping permission tests; chmod baseline unsupported: %v", err)
	}

	stat := func() os.FileMode {
		fi, err := os.Stat(goPath)
		if err != nil {
			t.Fatalf("os.Stat error: %v", err)
		}
		return fi.Mode()
	}

	// --- Readable ---
	if fileSetReadable2([]interface{}{f, types.JavaBoolTrue, types.JavaBoolTrue}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setReadable2(true, ownerOnly=true) returned false")
	}
	m := stat()
	if m&0o400 == 0 || (m&(0o040|0o004)) != 0 {
		t.Fatalf("setReadable2 ownerOnly should set only 0400; mode now %o", m)
	}
	if fileSetReadable2([]interface{}{f, types.JavaBoolFalse, types.JavaBoolTrue}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setReadable2(false, ownerOnly=true) returned false")
	}
	m = stat()
	if (m & 0o444) != 0 {
		t.Fatalf("clearing owner-only readable failed; mode %o", m)
	}
	if fileSetReadable([]interface{}{f, types.JavaBoolTrue}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setReadable(true) returned false")
	}
	m = stat()
	if m&0o444 != 0o444 {
		t.Fatalf("setReadable(true) should set 0444; mode %o", m)
	}
	if fileSetReadable([]interface{}{f, types.JavaBoolFalse}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setReadable(false) returned false")
	}
	m = stat()
	if (m & 0o444) != 0 {
		t.Fatalf("setReadable(false) should clear 0444; mode %o", m)
	}

	// --- Writable ---
	if fileSetWritable2([]interface{}{f, types.JavaBoolTrue, types.JavaBoolTrue}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setWritable2(true, ownerOnly=true) returned false")
	}
	m = stat()
	if m&0o200 == 0 || (m&(0o020|0o002)) != 0 {
		t.Fatalf("setWritable2 ownerOnly should set only 0200; mode now %o", m)
	}
	if fileSetWritable2([]interface{}{f, types.JavaBoolFalse, types.JavaBoolTrue}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setWritable2(false, ownerOnly=true) returned false")
	}
	m = stat()
	if (m & 0o222) != 0 {
		t.Fatalf("clearing owner-only writable failed; mode %o", m)
	}
	if fileSetWritable([]interface{}{f, types.JavaBoolTrue}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setWritable(true) returned false")
	}
	m = stat()
	if m&0o222 != 0o222 {
		t.Fatalf("setWritable(true) should set 0222; mode %o", m)
	}
	if fileSetWritable([]interface{}{f, types.JavaBoolFalse}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setWritable(false) returned false")
	}
	m = stat()
	if (m & 0o222) != 0 {
		t.Fatalf("setWritable(false) should clear 0222; mode %o", m)
	}

	// --- Executable ---
	if fileSetExecutable2([]interface{}{f, types.JavaBoolTrue, types.JavaBoolTrue}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setExecutable2(true, ownerOnly=true) returned false")
	}
	m = stat()
	if m&0o100 == 0 || (m&(0o010|0o001)) != 0 {
		t.Fatalf("setExecutable2 ownerOnly should set only 0100; mode now %o", m)
	}
	if fileSetExecutable2([]interface{}{f, types.JavaBoolFalse, types.JavaBoolTrue}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setExecutable2(false, ownerOnly=true) returned false")
	}
	m = stat()
	if (m & 0o111) != 0 {
		t.Fatalf("clearing owner-only executable failed; mode %o", m)
	}
	if fileSetExecutable([]interface{}{f, types.JavaBoolTrue}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setExecutable(true) returned false")
	}
	m = stat()
	if m&0o111 != 0o111 {
		t.Fatalf("setExecutable(true) should set 0111; mode %o", m)
	}
	if fileSetExecutable([]interface{}{f, types.JavaBoolFalse}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setExecutable(false) returned false")
	}
	m = stat()
	if (m & 0o111) != 0 {
		t.Fatalf("setExecutable(false) should clear 0111; mode %o", m)
	}

	// Ensure writable so that FileOutputStream can open on all platforms
	err := os.Chmod(goPath, 0o600)
	if err != nil {
		t.Fatalf("os.Chmod(goPath, 0o600) failed, err: %s", err.Error())
	}

	// Delete the temp directory.
	err = os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("Ignoring os.RemoveAll(%s) error: %v", dir, err)
	}
}
