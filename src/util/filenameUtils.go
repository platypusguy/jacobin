/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package util

import (
	"os"
	"strings"
)

// ConvertInternalClassNameToFilename accepts a class name with
// the JVM's internal format and converts it to a filename with
// OS-specific filepath separator chars.
func ConvertInternalClassNameToFilename(clName string) string {
	name := strings.ReplaceAll(clName, "/", "\\")
	name = strings.ReplaceAll(name, ".", "\\") + ".class"

	return ConvertToPlatformPathSeparators(name)
}

// ConvertClassFilenameToInternalFormat converts a class name
// with embedded . to the internal JVM class name format
func ConvertClassFilenameToInternalFormat(fName string) string {
	name := strings.TrimSuffix(fName, ".class")
	name = strings.ReplaceAll(name, ".", "/")
	return name
}

// ConvertToPlatformPathSeparators accepts a file path and,
// if necessary, converts the filepath separator characters
// to those used on the runtime platform
func ConvertToPlatformPathSeparators(pathIn string) string {
	osps := os.PathSeparator
	if strings.ContainsRune(pathIn, '/') && osps != '/' {
		return strings.ReplaceAll(pathIn, "/", string(osps))
	}

	if strings.ContainsRune(pathIn, '\\') && osps != '\\' {
		return strings.ReplaceAll(pathIn, "\\", string(osps))
	}
	return pathIn
}

// IsFilePartOfJDK accepts a filename and returns true if the filename
// is part of the JDK distribution
func IsFilePartOfJDK(filename *string) bool {
	return strings.HasPrefix(*filename, "java") ||
		strings.HasPrefix(*filename, "jdk") ||
		strings.HasPrefix(*filename, "sun")
}
