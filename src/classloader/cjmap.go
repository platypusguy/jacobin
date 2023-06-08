/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package classloader

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/shutdown"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// CJMAP contains class-to-Jmod-File relationships for all installed jmod files.
// No class information is stored.
// The key to the map is the class name in String format.
// The value associated with the key is the file name (not the full path) of the jmod file where the class is stored.
var CJMAP map[string]string

// Counting the map size (# of entries) since Go map has no such facility
// When CJMapInit is done (detected by CJMapFetch), a value of zero means that a severe error has occurred.
// In this case, a jacobin Shutdown should be executed.
var cjMapSize = 0

// Did we find a matching gob or did we build the map and gob?
var cjMapFoundGob bool

// Mutex for blocking CJMapFetch until the map is constructed
var cjMapMutex sync.Mutex

// Magic number of a jmod file:
const expectedMagicNumber = 0x4A4D

// Counter element name, needed when restoring a map from a gob file.
const counterElementName = "$COUNT"

// Log level for debugging:
const logLevel = log.FINE

// CJMapFetch retrieves the jmod file name associated with key = the class name.
// In the event that the class is not present there, nil is returned.
func CJMapFetch(key string) string {
	cjMapMutex.Lock()   // Is the map still being built by initialisation?
	cjMapMutex.Unlock() // Immediately unlock.
	if cjMapSize == 0 {
		msg := fmt.Sprintf("CJMapFetch: CJMAP size = 0 detected when key=%s", key)
		_ = log.Log(msg, log.SEVERE)
		shutdown.Exit(shutdown.JVM_EXCEPTION)
	}
	//fmt.Printf("DEBUG key=%s, jmod={%s}\n", key, CJMAP[key])
	return CJMAP[key]
}

// This function returns the number of entries in CJMAP.
func CJMapSize() int {
	return cjMapSize
}

// This function returns the number of entries in CJMAP.
func CJMapFoundGob() bool {
	return cjMapFoundGob
}

// This function initializes CJMAP and cjMapSize.
// Look for an existing gob file that matches global.JavaVersion value.
// If found, load the map from the gob file using buildMapFromGob.
// Otherwise,
//   - Create a new map from that installation's jmod files using buildMapFromJmods.
//   - Save the map to a gob file using saveMapToGob.
func CJMapInit() {

	global := globals.GetGlobalRef()

	// Open JacobinHome directory
	dirOpened, err := os.Open(global.JacobinHome)
	if err != nil {
		msg := fmt.Sprintf("CJMapInit: os.Open(%s) failed", global.JacobinHome)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		cjMapSize = 0
		return
	}
	msg := fmt.Sprintf("CJMapInit: JacobinHome is %s", global.JacobinHome)
	_ = log.Log(msg, logLevel)

	// Get all the file entries in the JacobinHome directory
	names, err := dirOpened.Readdirnames(0) // get all entries
	if err != nil {
		msg := fmt.Sprintf("CJMapInit: Readdirnames(%s) failed", global.JacobinHome)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		cjMapSize = 0
		return
	}

	// For each JacobinHome file, try to find a matching gob file
	for ix := range names {
		name := names[ix]
		//fmt.Printf("DEBUG name = %s\n", name)
		if strings.HasSuffix(name, ".gob") { // Gob file?
			version := strings.TrimSuffix(name, ".gob") // get rid of trailing .gom
			if version == global.JavaVersion {
				// Got a match!  Build map from it.
				gobFullPath := global.JacobinHome + string(os.PathSeparator) + name
				msg := fmt.Sprintf("CJMapInit: Gob file %s selected", gobFullPath)
				_ = log.Log(msg, logLevel)
				buildMapFromGob(gobFullPath)
				// If cjMapSize = 0, buildMapFrom Gob failed.
				cjMapFoundGob = true
				return
			}
		}
	}

	// No matching gob file
	msg = fmt.Sprintf("CJMapInit: No gob files matched Java version %s", global.JavaVersion)
	_ = log.Log(msg, logLevel)
	cjMapFoundGob = false
	buildMapFromJmods()
	if cjMapSize == 0 {
		return
	}
	saveMapToGob()

}

// This is the case where the map must be built from a gob file in global.JacobinHome.
// Lock the mutex and schedule (defer) an unlock upon return or crash.
// gobFile is the full path of the gob file in global.JacobinHome.
func buildMapFromGob(gobFilePath string) {

	// Initialise a new map.
	cjMapMutex.Lock()
	defer cjMapMutex.Unlock()
	CJMAP = make(map[string]string)
	cjMapSize = 0

	// Open input file
	inFile, err := os.Open(gobFilePath)
	if err != nil {
		msg := fmt.Sprintf("buildMapFromGob: os.Open(%s) failed", gobFilePath)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		cjMapSize = 0
		return
	}
	defer inFile.Close()

	// Create a decoder and receive a value.
	dinky := gob.NewDecoder(inFile)
	err = dinky.Decode(&CJMAP)
	if err != nil {
		msg := fmt.Sprintf("buildMapFromGob: gob Decode(%s) failed", gobFilePath)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		cjMapSize = 0
		return
	}

	gobSize := CJMAP[counterElementName]
	cjMapSize, err = strconv.Atoi(gobSize)
	if err != nil {
		msg := fmt.Sprintf("buildMapFromGob: Element (%s) is missing or misformatted", counterElementName)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		cjMapSize = 0
		return
	}
	msg := fmt.Sprintf("buildMapFromGob: Map size from gob file = %d", cjMapSize)
	_ = log.Log(msg, logLevel)

}

// This is the case where the map must be built from the jmod files of the Java installation
// indicated by JAVA_HOME.
// Lock the mutex and schedule (defer) an unlock upon return or crash.
func buildMapFromJmods() {

	global := globals.GetGlobalRef()

	// Initialise a new map.
	cjMapMutex.Lock()
	defer cjMapMutex.Unlock()
	CJMAP = make(map[string]string)
	cjMapSize = 0

	// Get path of jmods directory
	dirPath := global.JavaHome + string(os.PathSeparator) + "jmods"

	// Open jmods directory
	dirOpened, err := os.Open(dirPath)
	if err != nil {
		msg := fmt.Sprintf("buildMapFromJmods: os.Open(%s) failed", dirPath)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		cjMapSize = 0
		return
	}

	// Get all the file entries in the jmods directory
	names, err := dirOpened.Readdirnames(0) // get all entries
	if err != nil {
		_ = log.Log("buildMapFromJmods: Readdirnames(jmods directory) failed", log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		cjMapSize = 0
		return
	}

	// For each jmod file, process it
	count := 0
	for index := range names {
		count++
		name := names[index]
		jmodFullPath := filepath.Join(dirPath, name)
		if !processJmodFile(name, jmodFullPath) {
			cjMapSize = 0
			return
		}
	}

	CJMAP[counterElementName] = fmt.Sprint(cjMapSize)
	msg := fmt.Sprintf("buildMapFromJmods: Map built from %d jmod files", count)
	_ = log.Log(msg, logLevel)

}

// Given a jmod file, process all of the embedded class files.
// Called by buildMapFromJmods
// jmodFullPath: Full path of the jmod file under the Java jmods subdirectory
// jmodFileName: Just the jmod file name
func processJmodFile(jmodFileName string, jmodFullPath string) bool {

	// Open the jmods file
	_, err := os.Open(jmodFullPath)
	if err != nil {
		msg := fmt.Sprintf("processJmodFile: os.Open(%s) failed", jmodFullPath)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		return false
	}

	// Read entire file contents
	jmodBytes, err := os.ReadFile(jmodFullPath)
	if err != nil {
		msg := fmt.Sprintf("processJmodFile: os.ReadFile(%s) failed", jmodFullPath)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		return false
	}

	// Validate the file's magic number
	fileMagicNumber := binary.BigEndian.Uint16(jmodBytes[:2])
	if fileMagicNumber != expectedMagicNumber {
		msg := fmt.Sprintf("processJmodFile: fileMagicNumber != ExpectedMagicNumber in %s", jmodFullPath)
		_ = log.Log(msg, log.SEVERE)
		return false
	}

	// Skip over the jmod header so that it is recognized as a ZIP file
	offsetReader := bytes.NewReader(jmodBytes[4:])

	// Prepare the reader for the zip archive
	zipReader, err := zip.NewReader(offsetReader, int64(len(jmodBytes)-4))
	if err != nil {
		msg := fmt.Sprintf("processJmodFile: zip.NewReader failed(%s) failed", jmodFullPath)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
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
		classFileName := strings.Replace(fileEntry.Name, "classes"+string(os.PathSeparator), "", 1)

		// Add to map
		CJMAP[classFileName] = jmodFileName
		//fmt.Printf("DEBUG processJmodFile: classFileName=%s, jmodFileName=%s\n", classFileName, jmodFileName)

		// Add to count of classes
		countClasses++

		// Add to size of CJMAP
		cjMapSize++
	}

	msg := fmt.Sprintf("processJmodFile: Total classes added for %s = %d", jmodFileName, countClasses)
	_ = log.Log(msg, logLevel)

	return true

}

// Save the map to a gob file.
// No map locking is necessary.
func saveMapToGob() {

	global := globals.GetGlobalRef()
	gobFile := global.JacobinHome + string(os.PathSeparator) + global.JavaVersion + ".gob"
	// Open output gob file
	outFile, err := os.Create(gobFile)
	if err != nil {
		msg := fmt.Sprintf("saveMapToGob: os.Create(%s) failed", gobFile)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		cjMapSize = 0
		return
	}

	// Create a gob encoder and encode the cross-reference map
	inky := gob.NewEncoder(outFile)
	err = inky.Encode(&CJMAP)
	if err != nil {
		msg := fmt.Sprintf("saveMapToGob: gob Encode(%s) failed", gobFile)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		cjMapSize = 0
		return
	}

	// Close the output file
	err = outFile.Close()
	if err != nil {
		msg := fmt.Sprintf("saveMapToGob: close(%s) failed", gobFile)
		_ = log.Log(msg, log.SEVERE)
		_ = log.Log(err.Error(), log.SEVERE)
		cjMapSize = 0
	}

}
