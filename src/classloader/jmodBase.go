/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"jacobin/globals"
	"jacobin/trace"
	"strings"
)

// Walk the Base Jmod file and invoke ParseAndPostClass for each class found in the classlist
// Only called in one place: LoadBaseClasses.
func WalkBaseJmod() error {

	// Skip over the JMOD header so that it is recognized as a ZIP file
	global := globals.GetGlobalRef()
	ioReader := bytes.NewReader(global.JmodBaseBytes[4:])
	zipReader, err := zip.NewReader(ioReader, int64(len(global.JmodBaseBytes)-4))
	if err != nil {
		errMsg := fmt.Sprintf("WalkBaseJmod: zip.NewReader failed, err: %v", err)
		trace.Error(errMsg)
		return err
	}

	// Get the lib/classlist (bootstrap set of classes) if it exists
	bootstrapSet := getClasslist(*zipReader)
	useBootstrapSet := len(bootstrapSet) > 0

	// For each class file in the base jmod,
	// if it is in the classlist
	for _, classFile := range zipReader.File {

		// If not prefixed by "classes" or suffixed by ".class", skip this file
		if !strings.HasPrefix(classFile.Name, "classes") {
			continue
		}
		if !strings.HasSuffix(classFile.Name, ".class") {
			continue
		}

		// Remove prefix for bootstrap list check
		strapFileName := strings.Replace(classFile.Name, "classes/", "", 1)

		// Is there a bootstrap list?
		if useBootstrapSet {
			// Yes, make sure that this class is on the list
			_, onList := bootstrapSet[strapFileName]
			if !onList {
				continue
			}
		}

		// Open the class file
		rc, err := classFile.Open()
		if err != nil {
			return err
		}

		// Read all of the bytes
		classBytes, err := io.ReadAll(rc)
		if err != nil {
			return err
		}
		_ = rc.Close()

		// Parse and post class into MethArea
		ParseAndPostClass(&BootstrapCL, classFile.Name, classBytes)

	}

	return nil
}

// getClasslist returns the bootstrap lib/classlist as a Go-language map from the Java installation.
// There is a lib/classlist under the Java installation.
// However, that file only has entries from jmods/java.base.jmod and this classlist is duplicated as a member in that file.
// So, this function uses jmods/java.base.jmod to fetch the bootstrap map.
func getClasslist(reader zip.Reader) map[string]struct{} {
	classSet := make(map[string]struct{})

	classlist, err := reader.Open("lib/classlist")
	if err != nil {
		errMsg := fmt.Sprintf("getClasslist: reader.Open(lib/classlist) failed, err: %v", err)
		trace.Error(errMsg)
		return classSet
	}

	classlistContent, err := io.ReadAll(classlist)
	if err != nil {
		errMsg := fmt.Sprintf("getClasslist: io.ReadAll(classList) failed, err: %v", err)
		trace.Error(errMsg)
		return classSet
	}

	classes := strings.Split(string(classlistContent), "\n")

	var empty struct{}

	for _, c := range classes {
		if strings.HasSuffix(c, "\r") || strings.HasSuffix(c, "\n") {
			c = strings.TrimRight(c, "\r\n")
		}
		classSet[c+".class"] = empty
	}

	return classSet
}
