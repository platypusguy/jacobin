package gfunction

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/trace"
	"jacobin/src/types"
	"os"
	"path/filepath"
)

func Load_Io_File() {
	MethodSignatures["java/io/File.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	// Constructors
	MethodSignatures["java/io/File.<init>(Ljava/io/File;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  fileInit,
		}
	MethodSignatures["java/io/File.<init>(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  fileInit,
		}
	MethodSignatures["java/io/File.<init>(Ljava/lang/String;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  fileInit,
		}
	MethodSignatures["java/io/File.<init>(Ljava/net/URI;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  fileInit,
		}

	// Instance Methods (alphabetical)
	MethodSignatures["java/io/File.canExecute()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileCanExecute,
		}
	MethodSignatures["java/io/File.canRead()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileCanRead,
		}
	MethodSignatures["java/io/File.canWrite()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileCanWrite,
		}
	MethodSignatures["java/io/File.compareTo(Ljava/io/File;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  fileCompareTo,
		}
	MethodSignatures["java/io/File.createNewFile()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileCreate,
		}
	MethodSignatures["java/io/File.createTempFile(Ljava/lang/String;Ljava/lang/String;)Ljava/io/File;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  fileCreateTemp,
		}
	MethodSignatures["java/io/File.delete()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileDelete,
		}
	MethodSignatures["java/io/File.equals(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  fileEquals,
		}
	MethodSignatures["java/io/File.exists()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileExists,
		}
	MethodSignatures["java/io/File.getAbsolutePath()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileGetAbsolutePath,
		}
	MethodSignatures["java/io/File.getCanonicalPath()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileGetCanonicalPath,
		}
	MethodSignatures["java/io/File.getName()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileGetName,
		}
	MethodSignatures["java/io/File.getParent()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileGetParent,
		}
	MethodSignatures["java/io/File.getParentFile()Ljava/io/File;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileGetParentFile,
		}
	MethodSignatures["java/io/File.getPath()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileGetPath,
		}
	MethodSignatures["java/io/File.getFreeSpace()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileGetFreeSpace,
		}
	MethodSignatures["java/io/File.hashCode()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileHashCode,
		}
	MethodSignatures["java/io/File.isAbsolute()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileIsAbsolute,
		}
	MethodSignatures["java/io/File.isDirectory()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileIsDirectory,
		}
	MethodSignatures["java/io/File.isFile()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileIsFile,
		}
	MethodSignatures["java/io/File.isHidden()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileIsHidden,
		}
	MethodSignatures["java/io/File.lastModified()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileLastModified,
		}
	MethodSignatures["java/io/File.length()J"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileLength,
		}
	MethodSignatures["java/io/File.list()[Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileList,
		}
	MethodSignatures["java/io/File.list(Ljava/io/FilenameFilter;)[Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  fileListFiltered,
		}
	MethodSignatures["java/io/File.listFiles()[Ljava/io/File;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileListFiles,
		}
	MethodSignatures["java/io/File.listFiles(Ljava/io/FileFilter;)[Ljava/io/File;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  fileListFilesWithFileFilter,
		}
	MethodSignatures["java/io/File.listFiles(Ljava/io/FilenameFilter;)[Ljava/io/File;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  fileListFilesWithFilenameFilter,
		}
	MethodSignatures["java/io/File.mkdir()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileMkdir,
		}
	MethodSignatures["java/io/File.mkdirs()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileMkdirs,
		}
	MethodSignatures["java/io/File.renameTo(Ljava/io/File;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  fileRenameTo,
		}
	MethodSignatures["java/io/File.setExecutable(Z)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  fileSetExecutable,
		}
	MethodSignatures["java/io/File.setExecutable(ZZ)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  fileSetExecutable2,
		}
	MethodSignatures["java/io/File.setReadable(Z)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  fileSetReadable,
		}
	MethodSignatures["java/io/File.setReadable(ZZ)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  fileSetReadable2,
		}
	MethodSignatures["java/io/File.setReadOnly()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileSetReadOnly,
		}
	MethodSignatures["java/io/File.setWritable(Z)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  fileSetWritable,
		}
	MethodSignatures["java/io/File.setWritable(ZZ)Z"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  fileSetWritable2,
		}
	MethodSignatures["java/io/File.toURI()Ljava/net/URI;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/io/File.toURL()Ljava/net/URL;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/io/File.toPath()Ljava/nio/file/Path;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}
	MethodSignatures["java/io/File.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileToString,
		}

	// Static Methods
	MethodSignatures["java/io/File.createTempFile(Ljava/lang/String;Ljava/lang/String;Ljava/io/File;)Ljava/io/File;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  fileCreateTempWithDir,
		}
	MethodSignatures["java/io/File.listRoots()[Ljava/io/File;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  fileListRoots,
		}
}

// "java/io/File.<init>(Ljava/lang/String;)V"
// File file = new File(path);
func fileInit(params []interface{}) interface{} {

	// Get File object. Initialise the field map if required.
	objFile := params[0].(*object.Object)
	if objFile.FieldTable == nil {
		objFile.FieldTable = make(map[string]object.Field)
	}

	// Initialise the file status as "invalid" (=0).
	fld := object.Field{Ftype: types.Int, Fvalue: int64(0)}
	objFile.FieldTable[FileStatus] = fld

	// Get the argument path string object.
	objPath := params[1]
	if object.IsNull(objPath) {
		errMsg := "fileInit: Path object is null"
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}
	argPathStr := object.GoStringFromStringObject(objPath.(*object.Object))
	if argPathStr == "" {
		errMsg := "fileInit: String argument for path is empty"
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}

	// Create an absolute path string.
	absPathStr, err := filepath.Abs(argPathStr)
	if err != nil {
		errMsg := fmt.Sprintf("fileInit: filepath.Abs(%s) failed, reason: %s", argPathStr, err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}

	// Fill in File attributes that might get accessed by OpenJDK library member functions.

	fld = object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(absPathStr)}
	objFile.FieldTable[FilePath] = fld

	fld = object.Field{Ftype: types.Int, Fvalue: os.PathSeparator}
	objFile.FieldTable["separatorChar"] = fld

	fld = object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoByteArray([]byte{os.PathSeparator})}
	objFile.FieldTable["separator"] = fld

	fld = object.Field{Ftype: types.Int, Fvalue: os.PathListSeparator}
	objFile.FieldTable["pathSeparatorChar"] = fld

	fld = object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoByteArray([]byte{os.PathListSeparator})}
	objFile.FieldTable["pathSeparator"] = fld

	// Set status to "checked" (=1).
	fld = object.Field{Ftype: types.Int, Fvalue: int64(1)}
	objFile.FieldTable[FileStatus] = fld

	return nil
}

// "java/io/File.getPath()Ljava/lang/String;"
func fileGetPath(params []interface{}) interface{} {
	fld, ok := params[0].(*object.Object).FieldTable[FilePath]
	if !ok {
		errMsg := "fileGetPath: File object lacks a FilePath field"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	return object.StringObjectFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
}

// "java/io/File.isInvalid()Z"
func fileIsInvalid(params []interface{}) interface{} {
	status, ok := params[0].(*object.Object).FieldTable[FileStatus].Fvalue.(int64)
	if !ok {
		errMsg := "fileIsInvalid: File object lacks a FileStatus field"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	if status == 0 {
		return int64(1)
	} else {
		return int64(0)
	}
}

// "java/io/File.delete()Z"
func fileDelete(params []interface{}) interface{} {
	// Close the file if it is open (Windows).
	osFile, ok := params[0].(*object.Object).FieldTable[FileHandle].Fvalue.(*os.File)
	if ok {
		_ = osFile.Close()
	}

	// Get file path string.
	fld, ok := params[0].(*object.Object).FieldTable[FilePath]
	if !ok {
		errMsg := "fileDelete: File object lacks a FilePath field"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	pathStr := object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))

	err := os.Remove(pathStr)
	if err != nil {
		errMsg := fmt.Sprintf("fileDelete: Failed to remove file %s, reason: %s", pathStr, err.Error())
		trace.Error(errMsg)
		return int64(0)
	}
	return int64(1)
}

// "java/io/File.createNewFile()Z"
func fileCreate(params []interface{}) interface{} {
	// Get file path string.
	fld, ok := params[0].(*object.Object).FieldTable[FilePath]
	if !ok {
		errMsg := "fileCreate: File object lacks a FilePath field"
		return getGErrBlk(excNames.IOException, errMsg)
	}
	pathStr := object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))

	// Create the file and keep it open, storing the handle in the File object.
	osFile, err := os.Create(pathStr)
	if err != nil {
		errMsg := fmt.Sprintf("fileCreate: Failed to create file %s, reason: %s", pathStr, err.Error())
		trace.Error(errMsg)
		return types.JavaBoolFalse
	}

	// Copy the file handle into the File object; the test is responsible for closing it when finished.
	fld = object.Field{Ftype: types.FileHandle, Fvalue: osFile}
	params[0].(*object.Object).FieldTable[FileHandle] = fld

	return types.JavaBoolTrue
}

// --- Helpers ---
func fileGetPathString(obj *object.Object) (string, *GErrBlk) {
	fld, ok := obj.FieldTable[FilePath]
	if !ok {
		return "", getGErrBlk(excNames.IOException, "File object lacks a FilePath field")
	}
	jb, ok := fld.Fvalue.([]types.JavaByte)
	if !ok {
		return "", getGErrBlk(excNames.IOException, "FilePath field has invalid type")
	}
	return object.GoStringFromJavaByteArray(jb), nil
}

func newFileObjectFromPath(pathStr string) *object.Object {
	className := "java/io/File"
	obj := object.MakeEmptyObjectWithClassName(&className)
	// initialize fields similar to fileInit
	if obj.FieldTable == nil {
		obj.FieldTable = make(map[string]object.Field)
	}
	obj.FieldTable[FileStatus] = object.Field{Ftype: types.Int, Fvalue: int64(1)}
	obj.FieldTable[FilePath] = object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(pathStr)}
	obj.FieldTable["separatorChar"] = object.Field{Ftype: types.Int, Fvalue: os.PathSeparator}
	obj.FieldTable["separator"] = object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoByteArray([]byte{os.PathSeparator})}
	obj.FieldTable["pathSeparatorChar"] = object.Field{Ftype: types.Int, Fvalue: os.PathListSeparator}
	obj.FieldTable["pathSeparator"] = object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoByteArray([]byte{os.PathListSeparator})}
	return obj
}

// --- Implementations for File methods ---
func fileCanExecute(params []interface{}) interface{} {
	// Minimal: executable if file exists (we don't track execute bit portably)
	p, err := fileGetPathString(params[0].(*object.Object))
	if err != nil {
		return types.JavaBoolFalse
	}
	if _, e := os.Stat(p); e == nil {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func fileCanRead(params []interface{}) interface{} {
	p, err := fileGetPathString(params[0].(*object.Object))
	if err != nil {
		return types.JavaBoolFalse
	}
	if _, e := os.Stat(p); e == nil {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func fileCanWrite(params []interface{}) interface{} {
	p, err := fileGetPathString(params[0].(*object.Object))
	if err != nil {
		return types.JavaBoolFalse
	}
	if fi, e := os.Stat(p); e == nil {
		// On Unix, writable if not read-only for owner; we approximate by attempting to open write-only
		f, e2 := os.OpenFile(p, os.O_WRONLY, fi.Mode())
		if e2 == nil {
			_ = f.Close()
			return types.JavaBoolTrue
		}
	}
	return types.JavaBoolFalse
}

func fileCompareTo(params []interface{}) interface{} {
	p1, _ := fileGetPathString(params[0].(*object.Object))
	p2 := ""
	if other, ok := params[1].(*object.Object); ok && other != nil {
		p2, _ = fileGetPathString(other)
	}
	if p1 == p2 {
		return int64(0)
	}
	if p1 < p2 {
		return int64(-1)
	}
	return int64(1)
}

func fileCreateTemp(params []interface{}) interface{} {
	// instance createTempFile(prefix,suffix) -> File in default temp dir
	prefixObj := params[0].(*object.Object)
	suffixObj := params[1].(*object.Object)
	prefix := object.GoStringFromStringObject(prefixObj)
	suffix := object.GoStringFromStringObject(suffixObj)
	if suffix == "" {
		suffix = ".tmp"
	}
	f, err := os.CreateTemp("", prefix+"*")
	if err != nil {
		return object.Null
	}
	_ = f.Close()
	path := f.Name()
	if suffix != ".tmp" {
		// rename to ensure suffix
		newPath := path + suffix
		_ = os.Rename(path, newPath)
		path = newPath
	}
	return newFileObjectFromPath(path)
}

func fileEquals(params []interface{}) interface{} {
	p1, err1 := fileGetPathString(params[0].(*object.Object))
	if err1 != nil {
		return types.JavaBoolFalse
	}
	p2 := ""
	if other, ok := params[1].(*object.Object); ok && other != nil {
		p2, _ = fileGetPathString(other)
	}
	if p1 == p2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func fileExists(params []interface{}) interface{} {
	p, err := fileGetPathString(params[0].(*object.Object))
	if err != nil {
		return types.JavaBoolFalse
	}
	_, e := os.Stat(p)
	if e == nil {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func fileGetAbsolutePath(params []interface{}) interface{} {
	p, err := fileGetPathString(params[0].(*object.Object))
	if err != nil {
		return object.Null
	}
	abs, e := filepath.Abs(p)
	if e != nil {
		abs = p
	}
	return object.StringObjectFromGoString(abs)
}

func fileGetCanonicalPath(params []interface{}) interface{} {
	// Minimal: same as absolute path with filepath.Clean
	p, err := fileGetPathString(params[0].(*object.Object))
	if err != nil {
		return object.Null
	}
	abs, e := filepath.Abs(p)
	if e != nil {
		abs = p
	}
	can := filepath.Clean(abs)
	return object.StringObjectFromGoString(can)
}

func fileGetName(params []interface{}) interface{} {
	p, err := fileGetPathString(params[0].(*object.Object))
	if err != nil {
		return object.Null
	}
	name := filepath.Base(p)
	return object.StringObjectFromGoString(name)
}

func fileGetParent(params []interface{}) interface{} {
	p, err := fileGetPathString(params[0].(*object.Object))
	if err != nil {
		return object.Null
	}
	parent := filepath.Dir(p)
	return object.StringObjectFromGoString(parent)
}

func fileGetParentFile(params []interface{}) interface{} {
	p, err := fileGetPathString(params[0].(*object.Object))
	if err != nil {
		return object.Null
	}
	parent := filepath.Dir(p)
	return newFileObjectFromPath(parent)
}

func fileGetFreeSpace(params []interface{}) interface{} {
	// Minimal: return 0 (we do not query filesystem stats portable here)
	return int64(0)
}

func fileHashCode(params []interface{}) interface{} {
	p, _ := fileGetPathString(params[0].(*object.Object))
	// Java String hashCode algorithm
	h := int32(0)
	for i := 0; i < len(p); i++ {
		h = h*31 + int32(p[i])
	}
	return int64(h)
}

func fileIsAbsolute(params []interface{}) interface{} {
	p, err := fileGetPathString(params[0].(*object.Object))
	if err != nil {
		return types.JavaBoolFalse
	}
	if filepath.IsAbs(p) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func fileIsDirectory(params []interface{}) interface{} {
	p, err := fileGetPathString(params[0].(*object.Object))
	if err != nil {
		return types.JavaBoolFalse
	}
	fi, e := os.Stat(p)
	if e == nil && fi.IsDir() {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func fileIsFile(params []interface{}) interface{} {
	p, err := fileGetPathString(params[0].(*object.Object))
	if err != nil {
		return types.JavaBoolFalse
	}
	fi, e := os.Stat(p)
	if e == nil && !fi.IsDir() {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func fileIsHidden(params []interface{}) interface{} {
	p, _ := fileGetPathString(params[0].(*object.Object))
	name := filepath.Base(p)
	if len(name) > 0 && name[0] == '.' {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func fileLastModified(params []interface{}) interface{} {
	p, err := fileGetPathString(params[0].(*object.Object))
	if err != nil {
		return int64(0)
	}
	fi, e := os.Stat(p)
	if e != nil {
		return int64(0)
	}
	return fi.ModTime().UnixMilli()
}

func fileLength(params []interface{}) interface{} {
	p, err := fileGetPathString(params[0].(*object.Object))
	if err != nil {
		return int64(0)
	}
	fi, e := os.Stat(p)
	if e != nil || fi.IsDir() {
		return int64(0)
	}
	return fi.Size()
}

func fileList(params []interface{}) interface{} {
	p, err := fileGetPathString(params[0].(*object.Object))
	if err != nil {
		return object.MakeArrayFromRawArray([]*object.Object{})
	}
	entries, e := os.ReadDir(p)
	if e != nil {
		return object.MakeArrayFromRawArray([]*object.Object{})
	}
	var names []*object.Object
	for _, ent := range entries {
		names = append(names, object.StringObjectFromGoString(ent.Name()))
	}
	return object.MakeArrayFromRawArray(names)
}

func fileListFiltered(params []interface{}) interface{} {
	// Minimal: ignore filter and delegate to list()
	return fileList([]interface{}{params[0]})
}

func fileListFiles(params []interface{}) interface{} {
	p, err := fileGetPathString(params[0].(*object.Object))
	if err != nil {
		return object.MakeArrayFromRawArray([]*object.Object{})
	}
	entries, e := os.ReadDir(p)
	if e != nil {
		return object.MakeArrayFromRawArray([]*object.Object{})
	}
	var files []*object.Object
	for _, ent := range entries {
		child := filepath.Join(p, ent.Name())
		files = append(files, newFileObjectFromPath(child))
	}
	return object.MakeArrayFromRawArray(files)
}

func fileListFilesWithFileFilter(params []interface{}) interface{} {
	// Minimal: ignore filter
	return fileListFiles([]interface{}{params[0]})
}

func fileListFilesWithFilenameFilter(params []interface{}) interface{} {
	// Minimal: ignore filter
	return fileListFiles([]interface{}{params[0]})
}

func fileMkdir(params []interface{}) interface{} {
	p, err := fileGetPathString(params[0].(*object.Object))
	if err != nil {
		return types.JavaBoolFalse
	}
	if e := os.Mkdir(p, 0775); e == nil {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func fileMkdirs(params []interface{}) interface{} {
	p, err := fileGetPathString(params[0].(*object.Object))
	if err != nil {
		return types.JavaBoolFalse
	}
	if e := os.MkdirAll(p, 0775); e == nil {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func fileRenameTo(params []interface{}) interface{} {
	p1, err1 := fileGetPathString(params[0].(*object.Object))
	if err1 != nil {
		return types.JavaBoolFalse
	}
	p2 := ""
	if other, ok := params[1].(*object.Object); ok && other != nil {
		p2, _ = fileGetPathString(other)
	}
	if p2 == "" {
		return types.JavaBoolFalse
	}
	if e := os.Rename(p1, p2); e != nil {
		return types.JavaBoolFalse
	}
	// Update this object's path
	params[0].(*object.Object).FieldTable[FilePath] = object.Field{Ftype: types.ByteArray, Fvalue: object.JavaByteArrayFromGoString(p2)}
	return types.JavaBoolTrue
}

// --- permission helpers ---
func fileSetPermGeneric(obj *object.Object, enable bool, ownerOnly bool, allMask os.FileMode, ownerMask os.FileMode) int64 {
	p, _ := fileGetPathString(obj)
	if p == "" {
		return types.JavaBoolFalse
	}
	fi, e := os.Stat(p)
	if e != nil {
		return types.JavaBoolFalse
	}
	mode := fi.Mode()
	mask := allMask
	if ownerOnly {
		mask = ownerMask
	}
	if enable {
		mode = mode | mask
	} else {
		mode = mode &^ mask
	}
	if e2 := os.Chmod(p, mode); e2 != nil {
		return types.JavaBoolFalse
	}
	return types.JavaBoolTrue
}

func fileSetExecutable(params []interface{}) interface{} {
	// setExecutable(boolean)
	obj := params[0].(*object.Object)
	flag := params[1].(int64)
	enable := object.GoBooleanFromJavaBoolean(flag)
	// affect all (ownerOnly=false)
	return fileSetPermGeneric(obj, enable, false, 0o111, 0o100)
}
func fileSetExecutable2(params []interface{}) interface{} {
	// setExecutable(boolean, boolean ownerOnly)
	obj := params[0].(*object.Object)
	flag := params[1].(int64)
	ownerOnlyFlag := params[2].(int64)
	enable := object.GoBooleanFromJavaBoolean(flag)
	ownerOnly := object.GoBooleanFromJavaBoolean(ownerOnlyFlag)
	return fileSetPermGeneric(obj, enable, ownerOnly, 0o111, 0o100)
}
func fileSetReadable(params []interface{}) interface{} {
	// setReadable(boolean)
	obj := params[0].(*object.Object)
	flag := params[1].(int64)
	enable := object.GoBooleanFromJavaBoolean(flag)
	return fileSetPermGeneric(obj, enable, false, 0o444, 0o400)
}
func fileSetReadable2(params []interface{}) interface{} {
	// setReadable(boolean, boolean ownerOnly)
	obj := params[0].(*object.Object)
	flag := params[1].(int64)
	ownerOnlyFlag := params[2].(int64)
	enable := object.GoBooleanFromJavaBoolean(flag)
	ownerOnly := object.GoBooleanFromJavaBoolean(ownerOnlyFlag)
	return fileSetPermGeneric(obj, enable, ownerOnly, 0o444, 0o400)
}

func fileSetReadOnly(params []interface{}) interface{} {
	p, _ := fileGetPathString(params[0].(*object.Object))
	if p == "" {
		return types.JavaBoolFalse
	}
	fi, e := os.Stat(p)
	if e != nil {
		return types.JavaBoolFalse
	}
	mode := fi.Mode() &^ 0222
	if e2 := os.Chmod(p, mode); e2 != nil {
		return types.JavaBoolFalse
	}
	return types.JavaBoolTrue
}

func fileSetWritable(params []interface{}) interface{} {
	// setWritable(boolean)
	obj := params[0].(*object.Object)
	flag := params[1].(int64)
	enable := object.GoBooleanFromJavaBoolean(flag)
	return fileSetPermGeneric(obj, enable, false, 0o222, 0o200)
}
func fileSetWritable2(params []interface{}) interface{} {
	// setWritable(boolean, boolean ownerOnly)
	obj := params[0].(*object.Object)
	flag := params[1].(int64)
	ownerOnlyFlag := params[2].(int64)
	enable := object.GoBooleanFromJavaBoolean(flag)
	ownerOnly := object.GoBooleanFromJavaBoolean(ownerOnlyFlag)
	return fileSetPermGeneric(obj, enable, ownerOnly, 0o222, 0o200)
}

func fileToString(params []interface{}) interface{} {
	// Same as getPath()
	return fileGetPath(params)
}

func fileCreateTempWithDir(params []interface{}) interface{} {
	prefix := object.GoStringFromStringObject(params[0].(*object.Object))
	suffix := object.GoStringFromStringObject(params[1].(*object.Object))
	var dir string
	if d, ok := params[2].(*object.Object); ok && d != nil {
		d, _ := fileGetPathString(d)
		dir = d
	}
	if suffix == "" {
		suffix = ".tmp"
	}
	f, err := os.CreateTemp(dir, prefix+"*")
	if err != nil {
		return object.Null
	}
	_ = f.Close()
	path := f.Name()
	if suffix != ".tmp" {
		newPath := path + suffix
		_ = os.Rename(path, newPath)
		path = newPath
	}
	return newFileObjectFromPath(path)
}

func fileListRoots(params []interface{}) interface{} {
	// Minimal: return root directory of current OS ("/" on Unix)
	root := string(os.PathSeparator)
	return object.MakeArrayFromRawArray([]*object.Object{newFileObjectFromPath(root)})
}
