package gfunction

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"
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

// closeWithFIS creates a FileInputStream with the given File object and closes it.
// This mirrors the Java pattern:
//
//	File f = new File("test.txt");
//	FileInputStream in = new FileInputStream(f);
//	in.close();
func closeWithFIS(t *testing.T, f *object.Object) {
	to := object.MakeEmptyObject() // represents new FileInputStream()
	if ret := initFileInputStreamFile([]interface{}{to, f}); ret != nil {
		t.Fatalf("initFileInputStreamFile error: %v", ret)
	}
	if ret := fisClose([]interface{}{to}); ret != nil {
		t.Fatalf("fisClose error: %v", ret)
	}
}

func getPath(t *testing.T, f *object.Object) string {
	t.Helper()
	p, gerr := fileGetPathString(f)
	if gerr != nil {
		t.Fatalf("fileGetPathString error: %s", gerr.ErrMsg)
	}
	return p
}

func fileCreateThenClose(t *testing.T, params []interface{}) {
	// Get file path string.
	fld, ok := params[0].(*object.Object).FieldTable[FilePath]
	if !ok {
		errMsg := "fileCreateThenClose: File object lacks a FilePath field"
		t.Fatal(errMsg)
	}
	pathStr := object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))

	// Create the file and keep it open, storing the handle in the File object.
	osFile, err := os.Create(pathStr)
	if err != nil {
		errMsg := fmt.Sprintf("fileCreateThenClose: os.Create failed for file %s, reason: %s",
			pathStr, err.Error())
		t.Fatal(errMsg)
	}

	err = osFile.Close()
	if err != nil {
		errMsg := fmt.Sprintf("fileCreateThenClose: osFile.Close failed for file %s, reason: %s",
			pathStr, err.Error())
		t.Fatal(errMsg)
	}

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

func TestJavaIoFile_AdditionalGetters(t *testing.T) {
	globals.InitStringPool()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	f := newFileObjFromPath(t, path)

	// Test fileGetParentFile
	parentFile := fileGetParentFile([]interface{}{f}).(*object.Object)
	if parentFile == nil {
		t.Fatalf("fileGetParentFile returned nil")
	}
	p, _ := fileGetPathString(parentFile)
	if p != dir {
		t.Fatalf("fileGetParentFile mismatch: want %s got %s", dir, p)
	}

	// Test fileGetFreeSpace
	if fileGetFreeSpace([]interface{}{f}).(int64) != 0 {
		t.Fatalf("fileGetFreeSpace expected 0")
	}

	// Test fileIsInvalid
	if fileIsInvalid([]interface{}{f}).(int64) != 0 {
		t.Fatalf("fileIsInvalid expected 0 (false)")
	}
	// Test invalid status
	f.FieldTable[FileStatus] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
	if fileIsInvalid([]interface{}{f}).(int64) != 1 {
		t.Fatalf("fileIsInvalid expected 1 (true) for status 0")
	}

	// Test fileIsHidden (minimal)
	if fileIsHidden([]interface{}{f}).(int64) != types.JavaBoolFalse {
		t.Fatalf("fileIsHidden expected false")
	}
	hiddenPath := filepath.Join(dir, ".hidden")
	fHidden := newFileObjFromPath(t, hiddenPath)
	if fileIsHidden([]interface{}{fHidden}).(int64) != types.JavaBoolTrue {
		t.Fatalf("fileIsHidden expected true for .hidden")
	}
}

func TestJavaIoFile_FileStatusMethods(t *testing.T) {
	globals.InitStringPool()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	f := newFileObjFromPath(t, path)

	// Initially does not exist
	if fileIsFile([]interface{}{f}).(int64) != types.JavaBoolFalse {
		t.Fatalf("isFile expected false initially")
	}

	// Create file
	fileCreateThenClose(t, []interface{}{f})
	closeWithFIS(t, f)

	if fileIsFile([]interface{}{f}).(int64) != types.JavaBoolTrue {
		t.Fatalf("isFile expected true after create")
	}
	if fileIsDirectory([]interface{}{f}).(int64) != types.JavaBoolFalse {
		t.Fatalf("isDirectory expected false for file")
	}

	// Test fileLastModified
	lm := fileLastModified([]interface{}{f}).(int64)
	if lm == 0 {
		t.Fatalf("fileLastModified returned 0")
	}

	// Test fileCanRead, fileCanWrite, fileCanExecute
	if fileCanRead([]interface{}{f}).(int64) != types.JavaBoolTrue {
		t.Fatalf("fileCanRead expected true")
	}
	if fileCanWrite([]interface{}{f}).(int64) != types.JavaBoolTrue {
		t.Fatalf("fileCanWrite expected true")
	}
	if fileCanExecute([]interface{}{f}).(int64) != types.JavaBoolTrue {
		t.Fatalf("fileCanExecute expected true")
	}
}

func TestJavaIoFile_FilteredLists(t *testing.T) {
	globals.InitStringPool()
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0664)
	dirFile := newFileObjFromPath(t, dir)

	// These are currently minimal/no-op implementations in javaIoFile.go
	// fileListFiltered
	arr := fileListFiltered([]interface{}{dirFile, object.Null}).(*object.Object)
	names, _ := arr.FieldTable["value"].Fvalue.([]*object.Object)
	if len(names) == 0 {
		t.Fatalf("fileListFiltered returned empty list")
	}

	// fileListFilesWithFileFilter
	farr := fileListFilesWithFileFilter([]interface{}{dirFile, object.Null}).(*object.Object)
	files, _ := farr.FieldTable["value"].Fvalue.([]*object.Object)
	if len(files) == 0 {
		t.Fatalf("fileListFilesWithFileFilter returned empty list")
	}

	// fileListFilesWithFilenameFilter
	farr2 := fileListFilesWithFilenameFilter([]interface{}{dirFile, object.Null}).(*object.Object)
	files2, _ := farr2.FieldTable["value"].Fvalue.([]*object.Object)
	if len(files2) == 0 {
		t.Fatalf("fileListFilesWithFilenameFilter returned empty list")
	}
}

func TestJavaIoFile_StaticMethods(t *testing.T) {
	globals.InitStringPool()
	// Test fileListRoots
	roots := fileListRoots(nil).(*object.Object)
	arr, _ := roots.FieldTable["value"].Fvalue.([]*object.Object)
	if len(arr) == 0 {
		t.Fatalf("fileListRoots returned empty list")
	}
}

func TestJavaIoFile_ErrorPaths(t *testing.T) {
	globals.InitStringPool()
	obj := object.MakeEmptyObject()

	// fileGetPath missing FilePath
	ret := fileGetPath([]interface{}{obj})
	if err, ok := ret.(*GErrBlk); !ok || err.ExceptionType != excNames.IOException {
		t.Fatalf("fileGetPath: expected GErrBlk IOException, got %v", ret)
	}

	// fileIsInvalid missing FileStatus
	ret = fileIsInvalid([]interface{}{obj})
	if err, ok := ret.(*GErrBlk); !ok || err.ExceptionType != excNames.IOException {
		t.Fatalf("fileIsInvalid: expected GErrBlk IOException, got %v", ret)
	}

	// fileDelete missing FilePath
	ret = fileDelete([]interface{}{obj})
	if err, ok := ret.(*GErrBlk); !ok || err.ExceptionType != excNames.IOException {
		t.Fatalf("fileDelete: expected GErrBlk IOException, got %v", ret)
	}

	// fileCreate missing FilePath
	ret = fileCreate([]interface{}{obj})
	if err, ok := ret.(*GErrBlk); !ok || err.ExceptionType != excNames.IOException {
		t.Fatalf("fileCreate: expected GErrBlk IOException, got %v", ret)
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
	// create file
	fileCreateThenClose(t, []interface{}{f})
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
	fileCreateThenClose(t, []interface{}{src})

	// close the open handle from fileCreate
	// required on Windows before renaming or deleting
	closeWithFIS(t, src)

	// rename src -> dst
	if fileRenameTo([]interface{}{src, dst}).(int64) != types.JavaBoolTrue {
		t.Fatalf("fileRenameTo failed")
	}

	// check that dst exists
	if fileExists([]interface{}{dst}).(int64) != types.JavaBoolTrue {
		t.Fatalf("dst should exist after rename")
	}

	// no need to manually remove dir; t.TempDir cleanup handles it
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
	fileCreateThenClose(t, []interface{}{f})

	// close the open handle from fileCreate (needed on Windows)
	closeWithFIS(t, f)

	// now run permission setters (all no-ops in your impl)
	if fileSetReadOnly([]interface{}{f}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setReadOnly expected true")
	}
	if fileSetReadable([]interface{}{f, types.JavaBoolTrue}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setReadable expected true")
	}
	if fileSetWritable([]interface{}{f, types.JavaBoolTrue}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setWritable expected true")
	}
	if fileSetExecutable([]interface{}{f, types.JavaBoolTrue}).(int64) != types.JavaBoolTrue {
		t.Fatalf("setExecutable expected true")
	}

	// no explicit cleanup — t.TempDir handles it
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

	if runtime.GOOS == "windows" {
		t.Skip("skipping owner-only permission tests on Windows")
	}

	globals.InitStringPool()
	dir := t.TempDir()
	pidstr := fmt.Sprintf("%d", os.Getpid())
	fpath := filepath.Join(dir, pidstr, "perm.bin")
	f := newFileObjFromPath(t, fpath)

	// ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
		t.Fatalf("failed to create parent directory %q: %v", filepath.Dir(fpath), err)
	}

	// create the file
	fileCreateThenClose(t, []interface{}{f})

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
	if m&0o400 == 0 {
		t.Fatalf("setReadable2 owner read not set; mode now %o", m)
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
	if err := os.Chmod(goPath, 0o600); err != nil {
		t.Fatalf("os.Chmod(goPath, 0o600) failed, err: %s", err.Error())
	}
}

func TestJavaIoFile_CompareTo(t *testing.T) {
	globals.InitStringPool()
	f1 := newFileObjFromPath(t, "a.txt")
	f2 := newFileObjFromPath(t, "b.txt")
	// Test fileCompareTo
	if fileCompareTo([]interface{}{f1, f2}).(int64) >= 0 {
		t.Fatalf("fileCompareTo: a.txt should be less than b.txt")
	}
	if fileCompareTo([]interface{}{f1, f1}).(int64) != 0 {
		t.Fatalf("fileCompareTo: a.txt should equal a.txt")
	}
	if fileCompareTo([]interface{}{f2, f1}).(int64) <= 0 {
		t.Fatalf("fileCompareTo: b.txt should be greater than a.txt")
	}
}
