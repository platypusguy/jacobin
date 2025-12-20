package gfunction

import (
	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"path/filepath"
	"strings"
)

func getSep() string {
	if globals.OnWindows {
		return `\`
	}
	return `/`
}

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

	if globals.OnWindows {
		thisVal = strings.ToLower(thisVal)
		otherVal = strings.ToLower(otherVal)
	}

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

	if globals.OnWindows {
		if strings.EqualFold(thisStr, otherStr) {
			return types.JavaBoolTrue
		}
	} else {
		if thisStr == otherStr {
			return types.JavaBoolTrue
		}
	}
	return types.JavaBoolFalse
}

func filePathGetFileName(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))

	parts := getPathParts(thisStr)
	if len(parts) == 0 {
		return object.Null
	}
	fileName := parts[len(parts)-1]
	return object.StringObjectFromGoString(fileName)
}

func filePathGetName(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	i := params[1].(int64)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))

	parts := getPathParts(thisStr)
	if i < 0 || int(i) >= len(parts) {
		return getGErrBlk(excNames.IllegalArgumentException, "Path.getName: index out of bounds")
	}
	return object.StringObjectFromGoString(parts[i])
}

func filePathGetNameCount(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	parts := getPathParts(thisStr)
	return int64(len(parts))
}

func filePathGetParent(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))

	// If it's a root, return null
	root := filePathGetRoot([]interface{}{thisObj})
	if !object.IsNull(root) {
		rootStr := object.GoStringFromStringObject(root.(*object.Object).FieldTable["value"].Fvalue.(*object.Object))
		if thisStr == rootStr {
			return object.Null
		}
	}

	idx := strings.LastIndex(thisStr, getSep())
	if idx < 0 {
		return object.Null
	}
	if idx == 0 {
		if len(thisStr) > 1 {
			return newPath(getSep())
		}
		return object.Null
	}

	// Windows-specific: C:\foo -> idx is 2. parent is C:\
	if globals.OnWindows && idx == 2 && len(thisStr) >= 3 && thisStr[1] == ':' {
		return newPath(thisStr[:3])
	}

	parent := thisStr[:idx]
	if parent == "" {
		return object.Null
	}
	return newPath(parent)
}

func filePathGetRoot(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(
		thisObj.FieldTable["value"].Fvalue.(*object.Object),
	)

	if len(thisStr) == 0 {
		return object.Null
	}

	// Windows
	if globals.OnWindows {
		// UNC path: \\server\share\...
		if strings.HasPrefix(thisStr, `\\`) {
			// Find \\server\share\
			i := strings.Index(thisStr[2:], `\`)
			if i < 0 {
				return object.Null
			}
			j := strings.Index(thisStr[2+i+1:], `\`)
			root := ""
			if j < 0 {
				root = thisStr + `\`
			} else {
				root = thisStr[:2+i+1+j+1]
			}
			return newPath(root)
		}

		// Drive-letter absolute path: C:\...
		if len(thisStr) >= 3 &&
			((thisStr[0] >= 'A' && thisStr[0] <= 'Z') ||
				(thisStr[0] >= 'a' && thisStr[0] <= 'z')) &&
			thisStr[1] == ':' &&
			thisStr[2] == '\\' {
			return newPath(thisStr[:3])
		}

		// Rooted but not absolute (\foo) → no root
		return object.Null
	}

	// Unix
	if strings.HasPrefix(thisStr, "/") {
		return newPath("/")
	}

	return object.Null
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

func isAbsolute(path string) bool {
	if path == "" {
		return false
	}

	if globals.OnWindows {
		// UNC path: \\server\share\...
		if len(path) >= 2 && path[0] == '\\' && path[1] == '\\' {
			return true
		}

		// Drive-letter absolute path: C:\...
		if len(path) >= 3 &&
			((path[0] >= 'A' && path[0] <= 'Z') ||
				(path[0] >= 'a' && path[0] <= 'z')) &&
			path[1] == ':' &&
			(path[2] == '\\' || path[2] == '/') {
			return true
		}

		// Everything else (including "\foo") is not fully absolute
		return false
	}

	// Unix-like systems
	return path[0] == '/'
}

func filePathIsAbsolute(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	if isAbsolute(thisStr) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func filePathIterator(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	parts := strings.Split(thisStr, getSep())
	return newStringIterator(parts)
}

func filePathNormalize(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))

	var normalized string
	if globals.OnWindows {
		// Basic normalization for Windows paths when running on any OS
		// 1. Replace all / with \
		res := strings.ReplaceAll(thisStr, "/", `\`)

		rootObj := filePathGetRoot([]interface{}{thisObj})
		var rootStr string
		pathWithoutRoot := res
		if !object.IsNull(rootObj) {
			rootStr = object.GoStringFromStringObject(rootObj.(*object.Object).FieldTable["value"].Fvalue.(*object.Object))
			if strings.HasPrefix(res, rootStr) {
				pathWithoutRoot = res[len(rootStr):]
			}
		} else if strings.HasPrefix(res, `\`) {
			rootStr = `\`
			pathWithoutRoot = res[1:]
		}

		parts := getPathParts(pathWithoutRoot)

		// Handle .. and .
		var stack []string
		for _, p := range parts {
			if p == "." {
				continue
			}
			if p == ".." {
				if len(stack) > 0 {
					stack = stack[:len(stack)-1]
				}
				continue
			}
			stack = append(stack, p)
		}
		normalized = rootStr + strings.Join(stack, `\`)
	} else {
		normalized = filepath.Clean(thisStr)
	}
	return newPath(normalized)
}

func filePathRelativize(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	otherObj := params[1].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	otherStr := object.GoStringFromStringObject(otherObj.FieldTable["value"].Fvalue.(*object.Object))

	if isAbsolute(thisStr) != isAbsolute(otherStr) {
		return getGErrBlk(excNames.IllegalArgumentException, "Path.relativize: both paths must be either absolute or relative")
	}

	if globals.OnWindows {
		// Basic relativize for Windows paths
		// 1. Must have same root
		thisRootObj := filePathGetRoot([]interface{}{thisObj})
		otherRootObj := filePathGetRoot([]interface{}{otherObj})

		var thisRoot, otherRoot string
		if !object.IsNull(thisRootObj) {
			thisRoot = object.GoStringFromStringObject(thisRootObj.(*object.Object).FieldTable["value"].Fvalue.(*object.Object))
		}
		if !object.IsNull(otherRootObj) {
			otherRoot = object.GoStringFromStringObject(otherRootObj.(*object.Object).FieldTable["value"].Fvalue.(*object.Object))
		}

		if !strings.EqualFold(thisRoot, otherRoot) {
			return getGErrBlk(excNames.IllegalArgumentException, "Path.relativize: both paths must have the same root")
		}

		thisWithoutRoot := thisStr[len(thisRoot):]
		otherWithoutRoot := otherStr[len(otherRoot):]

		thisParts := getPathParts(thisWithoutRoot)
		otherParts := getPathParts(otherWithoutRoot)

		i := 0
		for i < len(thisParts) && i < len(otherParts) && strings.EqualFold(thisParts[i], otherParts[i]) {
			i++
		}

		var relParts []string
		for j := i; j < len(thisParts); j++ {
			relParts = append(relParts, "..")
		}
		relParts = append(relParts, otherParts[i:]...)

		if len(relParts) == 0 {
			return newPath("")
		}
		return newPath(strings.Join(relParts, `\`))
	}

	rel, err := filepath.Rel(thisStr, otherStr)
	if err != nil {
		return getGErrBlk(excNames.IllegalArgumentException, "Path.relativize: "+err.Error())
	}
	return newPath(rel)
}

func filePathResolve(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	otherStr := object.GoStringFromStringObject(params[1].(*object.Object))
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))

	if isAbsolute(otherStr) {
		return newPath(otherStr)
	}
	if otherStr == "" {
		return thisObj
	}

	if globals.OnWindows {
		// If otherStr starts with \ or /, it's rooted.
		// Resolve it against the root of thisStr if thisStr has one.
		if strings.HasPrefix(otherStr, `\`) || strings.HasPrefix(otherStr, `/`) {
			rootObj := filePathGetRoot([]interface{}{thisObj})
			if !object.IsNull(rootObj) {
				rootStr := object.GoStringFromStringObject(rootObj.(*object.Object).FieldTable["value"].Fvalue.(*object.Object))
				// If root is drive letter like C:\, we want C: + otherStr
				if len(rootStr) >= 2 && rootStr[1] == ':' {
					return newPath(rootStr[:2] + otherStr)
				}
			}
			return newPath(otherStr)
		}
	}

	res := thisStr
	if !strings.HasSuffix(res, getSep()) {
		res += getSep()
	}
	return newPath(res + otherStr)
}

func filePathResolvePath(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	otherObj := params[1].(*object.Object)
	otherStr := object.GoStringFromStringObject(otherObj.FieldTable["value"].Fvalue.(*object.Object))
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))

	if isAbsolute(otherStr) {
		return otherObj
	}
	if otherStr == "" {
		return thisObj
	}

	if globals.OnWindows {
		if strings.HasPrefix(otherStr, `\`) || strings.HasPrefix(otherStr, `/`) {
			rootObj := filePathGetRoot([]interface{}{thisObj})
			if !object.IsNull(rootObj) {
				rootStr := object.GoStringFromStringObject(rootObj.(*object.Object).FieldTable["value"].Fvalue.(*object.Object))
				if len(rootStr) >= 2 && rootStr[1] == ':' {
					return newPath(rootStr[:2] + otherStr)
				}
			}
			return otherObj
		}
	}

	res := thisStr
	if !strings.HasSuffix(res, getSep()) {
		res += getSep()
	}
	return newPath(res + otherStr)
}

func filePathResolveSibling(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	otherStr := object.GoStringFromStringObject(params[1].(*object.Object))

	parent := filePathGetParent([]interface{}{thisObj})
	if object.IsNull(parent) || isAbsolute(otherStr) {
		return newPath(otherStr)
	}
	parentPath := parent.(*object.Object)
	return filePathResolve([]interface{}{parentPath, params[1]})
}

func filePathResolveSiblingPath(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	otherObj := params[1].(*object.Object)
	otherStr := object.GoStringFromStringObject(otherObj.FieldTable["value"].Fvalue.(*object.Object))

	parent := filePathGetParent([]interface{}{thisObj})
	if object.IsNull(parent) || isAbsolute(otherStr) {
		return otherObj
	}
	parentPath := parent.(*object.Object)
	return filePathResolvePath([]interface{}{parentPath, otherObj})
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
	parts := getPathParts(thisStr)
	if start < 0 || end > len(parts) || start >= end {
		return getGErrBlk(excNames.IllegalArgumentException, "Path.subpath: invalid range")
	}
	sub := strings.Join(parts[start:end], getSep())
	return newPath(sub)
}

func getPathParts(path string) []string {
	var parts []string
	for _, p := range strings.Split(path, getSep()) {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return parts
}

func filePathToAbsolutePath(params []interface{}) interface{} {
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	if isAbsolute(thisStr) {
		return thisObj
	}
	cwd := globals.GetSystemProperty("user.dir")
	if cwd == "" {
		cwd, _ = os.Getwd()
	}

	if globals.OnWindows {
		// If it starts with \, it's rooted on current drive
		if strings.HasPrefix(thisStr, `\`) || strings.HasPrefix(thisStr, `/`) {
			// Find drive letter in cwd
			if len(cwd) >= 2 && cwd[1] == ':' {
				return newPath(cwd[:2] + thisStr)
			}
			// If no drive letter in cwd, just prepend \ if not present (unlikely for valid cwd)
			return newPath(thisStr)
		}
	}

	if !strings.HasSuffix(cwd, getSep()) {
		cwd += getSep()
	}
	return newPath(cwd + thisStr)
}

func filePathToRealPath(params []interface{}) interface{} {
	// For simplicity, just normalize
	thisObj := params[0].(*object.Object)
	thisStr := object.GoStringFromStringObject(thisObj.FieldTable["value"].Fvalue.(*object.Object))
	doubleSep := getSep() + getSep()
	normalized := strings.ReplaceAll(thisStr, doubleSep, getSep())
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
