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
	"jacobin/globals"
	"jacobin/log"
	"jacobin/shutdown"
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

	// Stat the base jmod file
	//jmodStat, err := os.Stat(jmodBasePath)
	//if err != nil {
	//	msg := fmt.Sprintf("GetBaseJmodBytes: os.Stat(%s) failed", jmodBasePath)
	//	_ = log.Log(msg, log.SEVERE)
	//	_ = log.Log(err.Error(), log.SEVERE)
	//	shutdown.Exit(shutdown.JVM_EXCEPTION)
	//}

	// Allocate byte array, JmodBaseBytes
	//global.JmodBaseBytes = make([]byte, jmodStat.Size())

	// Read the entire base jmod file contents (huge!)
	global.JmodBaseBytes, err = os.ReadFile(jmodBasePath)
	if err != nil {
		msg := fmt.Sprintf("GetBaseJmodBytes: os.ReadFile(%s) failed", jmodBasePath)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		shutdown.Exit(shutdown.JVM_EXCEPTION)
	}

	// Validate the file's magic number
	fileMagicNumber := binary.BigEndian.Uint16(global.JmodBaseBytes[:2])
	if fileMagicNumber != ExpectedMagicNumber {
		msg := fmt.Sprintf("GetBaseJmodBytes: fileMagicNumber != ExpectedMagicNumber in jmod file %s", jmodBasePath)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		shutdown.Exit(shutdown.JVM_EXCEPTION)
	}

	msg := fmt.Sprintf("GetBaseJmodBytes: jmodPath %s is loaded, %d bytes", jmodBasePath, len(global.JmodBaseBytes))
	_ = log.Log(msg, log.CLASS)

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
			msg := fmt.Sprintf("GetClassBytes: os.ReadFile(%s) failed", jmodPath)
			_ = log.Log(msg, log.SEVERE)
			_ = log.Log(err.Error(), log.SEVERE)
			return nil, err
		}

		// Validate the file's magic number
		fileMagicNumber := binary.BigEndian.Uint16(jmodBytes[:2])
		if fileMagicNumber != ExpectedMagicNumber {
			msg := fmt.Sprintf("GetClassBytes: fileMagicNumber != ExpectedMagicNumber in jmod file %s", jmodPath)
			_ = log.Log(msg, log.SEVERE)
			_ = log.Log(err.Error(), log.SEVERE)
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
		msg := fmt.Sprintf("GetClassBytes: zip.NewReader(%s) failed", jmodPath)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		return nil, err
	}

	// Open the file within the zip archive
	fileHandle, err := zipReader.Open(classFileName)
	if err != nil {
		msg := fmt.Sprintf("GetClassBytes: zipReader.Open(class file %s in jmod file %s) failed", classFileName, jmodPath)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		return nil, err
	}

	// Read entire class file contents
	classBytes, err := io.ReadAll(fileHandle)
	if err != nil {
		msg := fmt.Sprintf("GetClassBytes: os.ReadAll(class file %s in jmod file %s) failed", classFileName, jmodPath)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		return nil, err
	}

	// Success!
	msg := fmt.Sprintf("GetClassBytes: jmodPath %s, className %s was loaded", jmodPath, className)
	_ = log.Log(msg, log.CLASS)
	return classBytes, nil

}
