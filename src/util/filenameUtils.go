/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package util

import "strings"

// accepts a class name with the JVM's internal format and converts
// it to a filename (with backslashes). Returns "" on error.
func ConvertInternalClassNameToFilename(clName string) string {
	name := strings.ReplaceAll(clName, "/", "\\")
	name = strings.ReplaceAll(name, ".", "\\") + ".class"

	return name
}

func ConvertClassFilenameToInternalFormat(fName string) string {
	name := strings.TrimSuffix(fName, ".class")
	name = strings.ReplaceAll(name, ".", "/")
	return name
}
