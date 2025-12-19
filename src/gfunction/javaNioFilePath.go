package gfunction

import (
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"strings"
)

var sep = string(os.PathSeparator)

// Load_Nio_File_Path loads MethodSignatures entries for java.nio.file.Path
func Load_Nio_File_Path() {

	MethodSignatures["java/nio/file/Path.<clinit>()V"] =
		GMeth{ParamSlots: 0, GFunction: clinitGeneric}
	MethodSignatures["java/nio/file/Path.<init>()V"] =
		GMeth{ParamSlots: 0, GFunction: trapProtected}

	MethodSignatures["java/nio/file/Path.compareTo(Ljava/nio/file/Path;)I"] =
		GMeth{ParamSlots: 1, GFunction: filePathCompareTo}

	MethodSignatures["java/nio/file/Path.endsWith(Ljava/lang/String;)Z"] =
		GMeth{ParamSlots: 1, GFunction: filePathEndsWith}
	MethodSignatures["java/nio/file/Path.endsWith(Ljava/nio/file/Path;)Z"] =
		GMeth{ParamSlots: 1, GFunction: filePathEndsWithPath}

	MethodSignatures["java/nio/file/Path.equals(Ljava/lang/Object;)Z"] =
		GMeth{ParamSlots: 1, GFunction: filePathEquals}

	// leave getFileSystem trapped
	MethodSignatures["java/nio/file/Path.getFileSystem()Ljava/nio/file/FileSystem;"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}

	MethodSignatures["java/nio/file/Path.getFileName()Ljava/nio/file/Path;"] =
		GMeth{ParamSlots: 0, GFunction: filePathGetFileName}
	MethodSignatures["java/nio/file/Path.getName(I)Ljava/nio/file/Path;"] =
		GMeth{ParamSlots: 1, GFunction: filePathGetName}
	MethodSignatures["java/nio/file/Path.getNameCount()I"] =
		GMeth{ParamSlots: 0, GFunction: filePathGetNameCount}
	MethodSignatures["java/nio/file/Path.getParent()Ljava/nio/file/Path;"] =
		GMeth{ParamSlots: 0, GFunction: filePathGetParent}
	MethodSignatures["java/nio/file/Path.getRoot()Ljava/nio/file/Path;"] =
		GMeth{ParamSlots: 0, GFunction: filePathGetRoot}

	MethodSignatures["java/nio/file/Path.hashCode()I"] =
		GMeth{ParamSlots: 0, GFunction: filePathHashCode}

	MethodSignatures["java/nio/file/Path.isAbsolute()Z"] =
		GMeth{ParamSlots: 0, GFunction: filePathIsAbsolute}

	MethodSignatures["java/nio/file/Path.iterator()Ljava/util/Iterator;"] =
		GMeth{ParamSlots: 0, GFunction: filePathIterator}

	MethodSignatures["java/nio/file/Path.normalize()Ljava/nio/file/Path;"] =
		GMeth{ParamSlots: 0, GFunction: filePathNormalize}

	MethodSignatures["java/nio/file/Path.relativize(Ljava/nio/file/Path;)Ljava/nio/file/Path;"] =
		GMeth{ParamSlots: 1, GFunction: filePathRelativize}

	MethodSignatures["java/nio/file/Path.resolve(Ljava/lang/String;)Ljava/nio/file/Path;"] =
		GMeth{ParamSlots: 1, GFunction: filePathResolve}
	MethodSignatures["java/nio/file/Path.resolve(Ljava/nio/file/Path;)Ljava/nio/file/Path;"] =
		GMeth{ParamSlots: 1, GFunction: filePathResolvePath}

	MethodSignatures["java/nio/file/Path.resolveSibling(Ljava/lang/String;)Ljava/nio/file/Path;"] =
		GMeth{ParamSlots: 1, GFunction: filePathResolveSibling}
	MethodSignatures["java/nio/file/Path.resolveSibling(Ljava/nio/file/Path;)Ljava/nio/file/Path;"] =
		GMeth{ParamSlots: 1, GFunction: filePathResolveSiblingPath}

	MethodSignatures["java/nio/file/Path.startsWith(Ljava/lang/String;)Z"] =
		GMeth{ParamSlots: 1, GFunction: filePathStartsWith}
	MethodSignatures["java/nio/file/Path.startsWith(Ljava/nio/file/Path;)Z"] =
		GMeth{ParamSlots: 1, GFunction: filePathStartsWithPath}

	MethodSignatures["java/nio/file/Path.subpath(II)Ljava/nio/file/Path;"] =
		GMeth{ParamSlots: 2, GFunction: filePathSubpath}

	MethodSignatures["java/nio/file/Path.toAbsolutePath()Ljava/nio/file/Path;"] =
		GMeth{ParamSlots: 0, GFunction: filePathToAbsolutePath}
	MethodSignatures["java/nio/file/Path.toFile()Ljava/io/File;"] =
		GMeth{ParamSlots: 0, GFunction: trapFunction}
	MethodSignatures["java/nio/file/Path.toRealPath([Ljava/nio/file/LinkOption;)Ljava/nio/file/Path;"] =
		GMeth{ParamSlots: 1, GFunction: filePathToRealPath}

	MethodSignatures["java/nio/file/Path.toString()Ljava/lang/String;"] =
		GMeth{ParamSlots: 0, GFunction: filePathToString}

	// leave register functions trapped
	MethodSignatures["java/nio/file/Path.register(Ljava/nio/file/WatchService;[Ljava/nio/file/WatchEvent$Kind;[Ljava/nio/file/WatchEvent$Modifier;)Ljava/nio/file/WatchKey;"] =
		GMeth{ParamSlots: 3, GFunction: trapFunction}
	MethodSignatures["java/nio/file/Path.register(Ljava/nio/file/WatchService;[Ljava/nio/file/WatchEvent$Kind;)Ljava/nio/file/WatchKey;"] =
		GMeth{ParamSlots: 2, GFunction: trapFunction}
}

// ---- GFunction implementation attempt
// filePathCompareTo compares two Path objects by their string representations
// This mirrors Java's Comparable<Path> contract at a minimal level.
func filePathCompareTo(params []interface{}) interface{} {
	thisObj, ok := params[0].(*object.Object)
	if !ok || thisObj == nil {
		return getGErrBlk(
			excNames.NullPointerException,
			"Path.compareTo: this is null",
		)
	}

	otherObj, ok := params[1].(*object.Object)
	if !ok || otherObj == nil {
		return getGErrBlk(
			excNames.NullPointerException,
			"Path.compareTo: other is null",
		)
	}

	// ---- same-provider check
	if thisObj.FieldTable["provider"] != otherObj.FieldTable["provider"] {
		return getGErrBlk(
			excNames.ClassCastException,
			"Path.compareTo: incompatible Path providers",
		)
	}

	// ---- extract Java Strings
	thisStrObj, ok1 := thisObj.FieldTable["value"].Fvalue.(*object.Object)
	otherStrObj, ok2 := otherObj.FieldTable["value"].Fvalue.(*object.Object)

	if !ok1 || !ok2 ||
		!object.IsStringObject(thisStrObj) ||
		!object.IsStringObject(otherStrObj) {
		return getGErrBlk(
			excNames.ClassCastException,
			"Path.compareTo: backing value is not a String",
		)
	}

	thisVal := object.GoStringFromStringObject(thisStrObj)
	otherVal := object.GoStringFromStringObject(otherStrObj)

	if thisVal < otherVal {
		return int64(-1)
	}
	if thisVal > otherVal {
		return int64(1)
	}
	return int64(0)
}

func filePathEndsWith(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	suffixObj := params[1].(*object.Object)

	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	suffixStr := object.GoStringFromStringObject(suffixObj)

	if strings.HasSuffix(thisStr, suffixStr) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func filePathEndsWithPath(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	otherObj := params[1].(*object.Object)

	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	otherStr := object.GoStringFromStringObject(otherObj.FieldTable["value"].Fvalue.(*object.Object))

	if strings.HasSuffix(thisStr, otherStr) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func filePathEquals(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	otherObj := params[1].(*object.Object)

	if thisObj.FieldTable["provider"] != otherObj.FieldTable["provider"] {
		return types.JavaBoolFalse
	}

	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	otherStr := object.GoStringFromStringObject(otherObj.FieldTable["value"].Fvalue.(*object.Object))

	if thisStr == otherStr {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func filePathGetFileName(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))

	parts := strings.Split(thisStr, sep)
	fileName := parts[len(parts)-1]
	return object.StringObjectFromGoString(fileName)
}

func filePathGetName(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	i := params[1].(int64)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))

	parts := strings.Split(thisStr, sep)
	if i < 0 || int(i) >= len(parts) {
		return getGErrBlk(excNames.IllegalArgumentException, "Path.getName: index out of bounds")
	}
	return object.StringObjectFromGoString(parts[i])
}

func filePathGetNameCount(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	parts := strings.Split(thisStr, sep)
	return int64(len(parts))
}

func filePathGetParent(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	idx := strings.LastIndex(thisStr, sep)
	if idx <= 0 {
		return nil
	}
	parent := thisStr[:idx]
	return newPath(parent)
}

func filePathGetRoot(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	if len(thisStr) > 0 && strings.HasPrefix(thisStr, sep) {
		return newPath(sep)
	}
	return nil
}

func filePathHashCode(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	var hash int64
	for _, c := range thisStr {
		hash = 31*hash + int64(c)
	}
	return hash
}

func filePathIsAbsolute(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	if len(thisStr) > 0 && strings.HasPrefix(thisStr, sep) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func filePathIterator(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	parts := strings.Split(thisStr, sep)
	return newStringIterator(parts)
}

func filePathNormalize(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	doubleSep := sep + sep
	normalized := strings.ReplaceAll(thisStr, doubleSep, sep) // simple normalization
	return newPath(normalized)
}

func filePathRelativize(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	otherObj := params[1].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	otherStr := object.GoStringFromStringObject(otherObj.FieldTable["value"].Fvalue.(*object.Object))
	if strings.HasPrefix(otherStr, thisStr) {
		rel := strings.TrimPrefix(otherStr, thisStr)
		if len(rel) > 0 && strings.HasPrefix(rel, sep) {
			rel = rel[1:]
		}
		return newPath(rel)
	}
	return getGErrBlk(excNames.IllegalArgumentException, "Path.relativize: not a subpath")
}

func filePathResolve(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	otherObj := params[1].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	otherStr := object.GoStringFromStringObject(otherObj)
	if strings.HasPrefix(otherStr, sep) {
		return newPath(otherStr)
	}
	return newPath(thisStr + sep + otherStr)
}

func filePathResolvePath(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	otherObj := params[1].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	otherStr := object.GoStringFromStringObject(otherObj.FieldTable["value"].Fvalue.(*object.Object))
	if strings.HasPrefix(otherStr, sep) {
		return newPath(otherStr)
	}
	return newPath(thisStr + sep + otherStr)
}

func filePathResolveSibling(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	otherObj := params[1].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	otherStr := object.GoStringFromStringObject(otherObj)
	idx := strings.LastIndex(thisStr, sep)
	parent := ""
	if idx >= 0 {
		parent = thisStr[:idx]
	}
	return newPath(parent + sep + otherStr)
}

func filePathResolveSiblingPath(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	otherObj := params[1].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	otherStr := object.GoStringFromStringObject(otherObj.FieldTable["value"].Fvalue.(*object.Object))
	idx := strings.LastIndex(thisStr, sep)
	parent := ""
	if idx >= 0 {
		parent = thisStr[:idx]
	}
	return newPath(parent + sep + otherStr)
}

func filePathStartsWith(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	prefixObj := params[1].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	prefixStr := object.GoStringFromStringObject(prefixObj)
	if strings.HasPrefix(thisStr, prefixStr) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func filePathStartsWithPath(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	otherObj := params[1].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	otherStr := object.GoStringFromStringObject(otherObj.FieldTable["value"].Fvalue.(*object.Object))
	if strings.HasPrefix(thisStr, otherStr) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func filePathSubpath(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	start := int(params[1].(int64))
	end := int(params[2].(int64))
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	parts := strings.Split(thisStr, sep)
	if start < 0 || end > len(parts) || start >= end {
		return getGErrBlk(excNames.IllegalArgumentException, "Path.subpath: invalid range")
	}
	sub := strings.Join(parts[start:end], sep)
	return newPath(sub)
}

func filePathToAbsolutePath(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	if len(thisStr) > 0 && strings.HasPrefix(thisStr, sep) {
		return thisObj
	}
	// for simplicity, prepend sep as root
	return newPath(sep + thisStr)
}

func filePathToRealPath(params []interface{}) interface{} {
	// For simplicity, just normalize
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	doubleSep := sep + sep
	normalized := strings.ReplaceAll(thisStr, doubleSep, sep)
	return newPath(normalized)
}

func filePathToString(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	strObj := thisObj.FieldTable["value"].Fvalue.(*object.Object)
	return strObj // Return the Java String object directly
}

// --- Helper to create Path objects ---
func newPath(goPath string) *object.Object {
	className := "java/nio/file/Path"
	obj := object.MakeEmptyObjectWithClassName(&className)
	obj.FieldTable["value"] = object.Field{
		Ftype:  types.StringClassName,
		Fvalue: object.StringObjectFromGoString(goPath),
	}
	obj.FieldTable["provider"] = object.Field{Ftype: types.FileSystemProviderType,
		Fvalue: types.FileSystemProviderValue}
	return obj
}

// --- Helper to create a Java LinkedList object from a slice of strings ---
func newStringIterator(items []string) *object.Object {
	array := object.StringObjectArrayFromGoStringArray(items)
	listObj := object.MakePrimitiveObject("java/util/LinkedList", types.Array+types.StringClassName, array)
	return listObj
}
