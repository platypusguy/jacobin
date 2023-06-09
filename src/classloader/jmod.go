/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"jacobin/exceptions"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/shutdown"
	"os"
	"strings"
)

type WalkEntryFunc func(bytes []byte, filename string) error

// MagicNumber JMOD Magic Number
const MagicNumber = 0x4A4D

// Jmod Holds the file referring to a Java Module (JMOD)
// Allows walking a Java Module (JMOD). The `Walk` method will walk the module and invoke the `walk` parameter for all
// classes found. If there is a classlist file in lib\classlist (in the module), it will filter out any classes not
// contained in the classlist file; otherwise, all classes found in classes/ in the module.
type Jmod struct {
	File os.File
}

// Walk a Jmod file and invoke the indicated WalkEntryFunc for each class found in the classlist
// Only called in one place: LoadJmodClasses.
func (jmodFile *Jmod) Walk(walk WalkEntryFunc) error {
	b, err := os.ReadFile(jmodFile.File.Name())
	if err != nil {
		return err
	}

	fileMagic := binary.BigEndian.Uint16(b[:2])

	if fileMagic != MagicNumber {

		if !globals.GetGlobalRef().StrictJDK {
			msg := fmt.Sprintf("An IOException occurred reading %s: the magic number is invalid. Expected: %x, Got: %x", jmodFile.File.Name(), MagicNumber, fileMagic)
			_ = log.Log(msg, log.SEVERE)
		}

		exceptions.JVMexception(exceptions.IOException, fmt.Sprintf("Invalid JMOD file: %s", jmodFile.File.Name()))
		shutdown.Exit(shutdown.JVM_EXCEPTION)
	}

	// Skip over the JMOD header so that it is recognized as a ZIP file
	offsetReader := bytes.NewReader(b[4:])

	r, err := zip.NewReader(offsetReader, int64(len(b)-4))
	if err != nil {
		_ = log.Log(err.Error(), log.WARNING)
		return err
	}

	classSet := getClasslist(*r)

	useClassSet := len(classSet) > 0

	for _, f := range r.File {
		if !strings.HasPrefix(f.Name, "classes") {
			continue
		}

		classFileName := strings.Replace(f.Name, "classes/", "", 1)

		if useClassSet {
			_, ok := classSet[classFileName]
			if !ok {
				continue
			}
		} else {
			if !strings.HasSuffix(f.Name, ".class") {
				continue
			}
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		b, err := io.ReadAll(rc)
		if err != nil {
			return err
		}

		_ = walk(b, jmodFile.File.Name()+"+"+f.Name)

		_ = rc.Close()
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
		_ = log.Log(err.Error(), log.CLASS)
		_ = log.Log("Unable to read lib/classlist from jmod file. Loading all classes in jmod file.", log.CLASS)
		return classSet
	}

	classlistContent, err := io.ReadAll(classlist)
	if err != nil {
		_ = log.Log(err.Error(), log.CLASS)
		_ = log.Log("Unable to read lib/classlist from jmod file. Loading all classes in jmod file.", log.CLASS)
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
