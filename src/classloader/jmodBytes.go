/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"jacobin/src/globals"
	"jacobin/src/shutdown"
	"jacobin/src/trace"
	"os"
)

const ExpectedMagicNumber = 0x4A4D
const BaseJmodFileName = "java.base.jmod"

// Load the entirety of the base jmod file into the byte cache: JmodBaseBytes
// Called during classloader initialisation
// Any error --> shutdown
func GetBaseJmodBytes() {

	var err error
	global := globals.GetGlobalRef()
	jmodBasePath := global.JavaHome + string(os.PathSeparator) + "jmods" + string(os.PathSeparator) + BaseJmodFileName

	// Read the entire base jmod file contents (huge!)
	global.JmodBaseBytes, err = os.ReadFile(jmodBasePath)
	if err != nil {
		errMsg := fmt.Sprintf("GetBaseJmodBytes: os.ReadFile(%s) failed, err: %v", jmodBasePath, err)
		trace.Error(errMsg)
		shutdown.Exit(shutdown.JVM_EXCEPTION)
	}

	// Validate the file's magic number
	fileMagicNumber := binary.BigEndian.Uint16(global.JmodBaseBytes[:2])
	if fileMagicNumber != ExpectedMagicNumber {
		errMsg := fmt.Sprintf("GetBaseJmodBytes: fileMagicNumber != ExpectedMagicNumber in jmod file %s, err: %v", jmodBasePath, err)
		trace.Error(errMsg)
		shutdown.Exit(shutdown.JVM_EXCEPTION)
	}

	if globals.TraceCloadi {
		infoMsg := fmt.Sprintf("GetBaseJmodBytes: jmodPath %s is loaded, %d bytes", jmodBasePath, len(global.JmodBaseBytes))
		trace.Trace(infoMsg)
	}

}

// For the given jmod and class name, return the class byte array to caller
func GetClassBytes(jmodFileName string, className string) ([]byte, error) {

	var jmodBytes []byte // <-- used if jmod file is not java.base.jmod
	var err error
	var ioReader *bytes.Reader
	var newReaderLength int64

	global := globals.GetGlobalRef()
	jmodPath := global.JavaHome + string(os.PathSeparator) + "jmods" + string(os.PathSeparator) + jmodFileName
	classFileName := "classes/" + className + ".class"

	//fmt.Printf("DEBUG GetClassBytes: jmod=%s, class=%s\n", jmodFileName, className)
	if jmodFileName == BaseJmodFileName {
		// Already loaded in JmodBaseBytes during classloader initialisation
		// Skip over the jmod header so that it is recognized as a ZIP file
		ioReader = bytes.NewReader(global.JmodBaseBytes[4:])
		newReaderLength = int64(len(global.JmodBaseBytes) - 4)
	} else {
		// Not the base jmod
		// Read entire jmod file contents
		jmodBytes, err = os.ReadFile(jmodPath)
		if err != nil {
			errMsg := fmt.Sprintf("GetClassBytes: os.ReadFile(%s) failed, err: %v", jmodPath, err)
			trace.Error(errMsg)
			return nil, err
		}

		// Validate the file's magic number
		fileMagicNumber := binary.BigEndian.Uint16(jmodBytes[:2])
		if fileMagicNumber != ExpectedMagicNumber {
			errMsg := fmt.Sprintf("GetClassBytes: fileMagicNumber != ExpectedMagicNumber in jmod file %s", jmodPath)
			trace.Error(errMsg)
			return nil, err
		}

		// Skip over the jmod header so that it is recognized as a ZIP file
		ioReader = bytes.NewReader(jmodBytes[4:])
		newReaderLength = int64(len(jmodBytes) - 4)

	}

	// Prepare the reader for the zip archive
	//fmt.Printf("DEBUG GetClassBytes: zip.NewReader newReaderLength=%d\n", newReaderLength)
	zipReader, err := zip.NewReader(ioReader, newReaderLength)
	if err != nil {
		errMsg := fmt.Sprintf("GetClassBytes: zip.NewReader(%s) failed, err: %v", jmodPath, err)
		trace.Error(errMsg)
		return nil, err
	}

	// Open the file within the zip archive
	fileHandle, err := zipReader.Open(classFileName)
	if err != nil {
		errMsg := fmt.Sprintf("GetClassBytes: zipReader.Open(class file %s in jmod file %s) failed, err: %v", classFileName, jmodPath, err)
		trace.Error(errMsg)
		return nil, err
	}

	// Read entire class file contents
	classBytes, err := io.ReadAll(fileHandle)
	if err != nil {
		errMsg := fmt.Sprintf("GetClassBytes: os.ReadAll(class file %s in jmod file %s) failed, err: %v", classFileName, jmodPath, err)
		trace.Error(errMsg)
		return nil, err
	}

	// Success!
	return classBytes, nil

}
