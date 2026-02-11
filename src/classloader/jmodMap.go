/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-4 by Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"jacobin/src/globals"
	"jacobin/src/shutdown"
	"jacobin/src/trace"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// JMODMAP contains class-to-Jmod-File relationships for all installed jmod files.
// No class information is stored.
// The key to the map is the class name in String format.
// The value is the file name (not the full path) of the jmod file where the class is stored.
var JMODMAP map[string]string

// Counting the map size (# of entries) since Go map has no such facility
// When JmodMapInit is done (detected by JmodMapFetch), a value of zero means that a severe error has occurred.
// In this case, a jacobin Shutdown should be executed.
var jmodMapSize = 0

// Did we find a matching gob or did we build the map and gob?
var jmodMapFoundGob bool

// Mutex for blocking JmodMapFetch until the map is constructed
var jmodMapMutex sync.Mutex

// Magic number of a jmod file:
const expectedMagicNumber = 0x4A4D

// Counter element name, needed when restoring a map from a gob file.
const counterElementName = "$COUNT"

// JmodMapFetch retrieves the jmod file name associated with key = the class name.
// The input class name is suffixed with ".class" before accessing the map.
// In the event that the class is not present there, nil is returned.
func JmodMapFetch(className string) string {
	jmodMapMutex.Lock()   // Wait if the map still being built by initialisation.
	jmodMapMutex.Unlock() // Immediately unlock.
	if jmodMapSize == 0 {
		errMsg := fmt.Sprintf("JmodMapFetch: JMODMAP size = 0 detected when key=%s", className)
		trace.Error(errMsg)
		shutdown.Exit(shutdown.JVM_EXCEPTION)
	}
	jmodFile := JMODMAP[className+".class"]
	return jmodFile
}

// JmodMapSize returns the number of entries in JMODMAP.
func JmodMapSize() int {
	return jmodMapSize
}

// JmodMapFoundGob returns the number of entries in JMODMAP.
func JmodMapFoundGob() bool {
	return jmodMapFoundGob
}

// JmodMapInit initializes JMODMAP and jmodMapSize.
// Look for an existing gob file that matches global.JavaVersion value.
// If found, load the map from the gob file using buildMapFromGob.
// Otherwise,
//   - Create a new map from that installation's jmod files using buildMapFromJmods.
//   - Save the map to a gob file using saveMapToGob.
func JmodMapInit() {

	global := globals.GetGlobalRef()

	// Open JacobinHome directory
	dirOpened, err := os.Open(global.JacobinHome)
	if err != nil {
		errMsg := fmt.Sprintf("JmodMapInit: os.Open(%s) failed, err: %v", global.JacobinHome, err)
		trace.Error(errMsg)
		jmodMapSize = 0
		return
	}

	// Get all the file entries in the JacobinHome directory
	names, err := dirOpened.Readdirnames(0) // get all entries
	if err != nil {
		errMsg := fmt.Sprintf("JmodMapInit: Readdirnames(%s) failed, err: %v", global.JacobinHome, err)
		trace.Error(errMsg)
		jmodMapSize = 0
		return
	}

	// For each JacobinHome file, try to find a matching gob file
	jmodMapFoundGob = false
	for ix := range names {
		name := names[ix]
		if strings.HasSuffix(name, ".gob") { // Gob file?
			version := strings.TrimSuffix(name, ".gob") // get rid of trailing .gob
			if version == global.JavaVersion {
				// Got a match!  Build map from it.
				gobFullPath := global.JacobinHome + string(os.PathSeparator) + name
				if !buildMapFromGob(gobFullPath) {
					// Gob file trouble
					// Force re-creation
					break
				}

				// Map built from gob file succeeded
				jmodMapFoundGob = true
				return
			}
		}
	}

	// No matching gob file or had gob file trouble (force re-creation).
	buildMapFromJmods()
	if jmodMapSize == 0 {
		return
	}
	saveMapToGob()
}

// This is the case where the map must be built from a gob file in global.JacobinHome.
// Lock the mutex and schedule (defer) an unlock upon return or crash.
// gobFile is the full path of the gob file in global.JacobinHome.
func buildMapFromGob(gobFilePath string) bool {

	// Initialise a new map.
	jmodMapMutex.Lock()
	defer jmodMapMutex.Unlock()
	JMODMAP = make(map[string]string)
	jmodMapSize = 0

	// Open input file
	inFile, err := os.Open(gobFilePath)
	if err != nil {
		errMsg := fmt.Sprintf("buildMapFromGob: os.Open(%s) failed, err: %v", gobFilePath, err)
		trace.Error(errMsg)
		jmodMapSize = 0
		return false
	}
	defer inFile.Close()

	// Create a decoder and receive a value.
	decoder := gob.NewDecoder(inFile)
	err = decoder.Decode(&JMODMAP)
	if err != nil {
		errMsg := fmt.Sprintf("buildMapFromGob: gob Decode(%s) failed, err: %v", gobFilePath, err)
		trace.Error(errMsg)
		jmodMapSize = 0
		return false
	}

	gobSize := JMODMAP[counterElementName]
	jmodMapSize, err = strconv.Atoi(gobSize)
	if err != nil {
		errMsg := fmt.Sprintf("buildMapFromGob: Element (%s) is missing or misformatted, err: %v", counterElementName, err)
		trace.Error(errMsg)
		jmodMapSize = 0
		return false
	}

	// Success!
	return true
}

// This is the case where the map must be built from the jmod files of the Java installation
// indicated by JAVA_HOME.
// Lock the mutex and schedule (defer) an unlock upon return or crash.
func buildMapFromJmods() {

	global := globals.GetGlobalRef()

	// Initialise a new map.
	jmodMapMutex.Lock()
	defer jmodMapMutex.Unlock()
	JMODMAP = make(map[string]string)
	jmodMapSize = 0

	// Get path of jmods directory
	dirPath := global.JavaHome + string(os.PathSeparator) + "jmods"

	// Open jmods directory
	dirOpened, err := os.Open(dirPath)
	if err != nil {
		errMsg := fmt.Sprintf("buildMapFromJmods: os.Open(%s) failed, err: %v", dirPath, err)
		trace.Error(errMsg)
		jmodMapSize = 0
		return
	}

	// Get all the file entries in the jmods directory
	names, err := dirOpened.Readdirnames(0) // get all entries
	if err != nil {
		errMsg := fmt.Sprintf("buildMapFromJmods: Readdirnames(jmods directory) failed, err: %v", err)
		trace.Error(errMsg)
		jmodMapSize = 0
		return
	}

	// For each jmod file, process it
	count := 0
	for index := range names {
		count++
		name := names[index]
		jmodFullPath := filepath.Join(dirPath, name)
		if !processJmodFile(name, jmodFullPath) {
			jmodMapSize = 0
			return
		}
	}

	JMODMAP[counterElementName] = fmt.Sprint(jmodMapSize)
}

// Given a jmod file, process all of the embedded class files.
// Called by buildMapFromJmods
// jmodFullPath: Full path of the jmod file under the Java jmods subdirectory
// jmodFileName: Just the jmod file name
func processJmodFile(jmodFileName string, jmodFullPath string) bool {

	// Open the jmods file
	_, err := os.Open(jmodFullPath)
	if err != nil {
		errMsg := fmt.Sprintf("processJmodFile: os.Open(%s) failed, err: %v", jmodFullPath, err)
		trace.Error(errMsg)
		jmodMapSize = 0
		return false
	}

	// Read entire file contents
	jmodBytes, err := os.ReadFile(jmodFullPath)
	if err != nil {
		errMsg := fmt.Sprintf("processJmodFile: os.ReadFile(%s) failed, err: %v", jmodFullPath, err)
		trace.Error(errMsg)
		jmodMapSize = 0
		return false
	}

	// Validate the file's magic number
	fileMagicNumber := binary.BigEndian.Uint16(jmodBytes[:2])
	if fileMagicNumber != expectedMagicNumber {
		errMsg := fmt.Sprintf("processJmodFile: fileMagicNumber != ExpectedMagicNumber in %s", jmodFullPath)
		trace.Error(errMsg)
		jmodMapSize = 0
		return false
	}

	// Skip over the jmod header so that it is recognized as a ZIP file
	offsetReader := bytes.NewReader(jmodBytes[4:])

	// Prepare the reader for the zip archive
	zipReader, err := zip.NewReader(offsetReader, int64(len(jmodBytes)-4))
	if err != nil {
		errMsg := fmt.Sprintf("processJmodFile: zip.NewReader failed(%s) failed, err: %v", jmodFullPath, err)
		trace.Error(errMsg)
		jmodMapSize = 0
		return false
	}

	// For each file entry within the zip reader, process it
	countClasses := 0
	for _, fileEntry := range zipReader.File {

		// Has the right prefix and suffix?
		if !strings.HasPrefix(fileEntry.Name, "classes") {
			continue
		}
		if !strings.HasSuffix(fileEntry.Name, ".class") {
			continue
		}

		// Remove the "classes" + string(os.PathSeparator) prefix.
		classFileName := strings.Replace(fileEntry.Name, "classes/", "", 1)

		// Add to map
		JMODMAP[classFileName] = jmodFileName
		if classFileName == "java/security/interfaces/DHPublicKey.class" {
			fmt.Printf("DEBUG processJmodFile: stored classFileName=%s, jmodFileName=%s\n", classFileName, jmodFileName)
		}

		// Add to count of classes
		countClasses++

		// Add to size of JMODMAP
		jmodMapSize++
	}

	return true
}

// Save the map to a gob (that is, a binary go object file).
// No map locking is necessary.
func saveMapToGob() {

	global := globals.GetGlobalRef()
	gobFile := global.JacobinHome + string(os.PathSeparator) + global.JavaVersion + ".gob"
	// Open output gob file
	_ = os.Remove(gobFile)
	outFile, err := os.Create(gobFile)
	if err != nil {
		errMsg := fmt.Sprintf("saveMapToGob: os.Create(%s) failed, err: %v", gobFile, err)
		trace.Error(errMsg)
		jmodMapSize = 0
		return
	}

	// Create a gob encoder and encode the cross-reference map
	encoder := gob.NewEncoder(outFile)
	err = encoder.Encode(JMODMAP)
	if err != nil {
		errMsg := fmt.Sprintf("saveMapToGob: gob Encode(%s) failed, err: %v", gobFile, err)
		trace.Error(errMsg)
		jmodMapSize = 0
		return
	}

	// Close the output file
	err = outFile.Close()
	if err != nil {
		errMsg := fmt.Sprintf("saveMapToGob: close(%s) failed, err: %v", gobFile, err)
		trace.Error(errMsg)
		jmodMapSize = 0
	}
}
