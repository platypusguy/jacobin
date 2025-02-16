/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-4 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package util

import (
	"os"
	"path/filepath"
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

func ConvertInternalClassNameToUserFormat(fName string) string {
	name := strings.ReplaceAll(fName, "/", ".")
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
	str := ConvertInternalClassNameToUserFormat(*filename)
	return strings.HasPrefix(str, "java.") ||
		strings.HasPrefix(str, "jdk.") ||
		strings.HasPrefix(str, "com.sun") ||
		strings.HasPrefix(str, "sun.")
}

// SearchDirByFileExtension searches a directory and its subdirectories for
// files with the extension. It returns a pointer to a slice of strings
// containing the path for every qualifying file, or nil if an error has
// occurred. If no qualifying file is found, the returned pointer points
// to a slice of length 0.
func SearchDirByFileExtension(dir, extension string) *[]string {
	var filenames []string

	_, direrr := os.Stat(dir)
	if os.IsNotExist(direrr) {
		return nil
	}

	// Walk through the directory tree
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		// Check if the file has the searched-for extension
		if !info.IsDir() && filepath.Ext(info.Name()) == "."+extension {
			filenames = append(filenames, path)
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return &filenames
}
