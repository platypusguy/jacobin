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

// JMODMAP contains class-to-Jmod-File relationships for all installed jmod files.
// No class information is stored.
// The key to the map is the class name in String format.
// The value associated with the key is the file name (not the full path) of the jmod file where the class is stored.
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

// Log level for debugging:
const logLevel = log.FINE

// JmodMapFetch retrieves the jmod file name associated with key = the class name.
// The input class name is suffixed with ".class" before accessing the map.
// In the event that the class is not present there, nil is returned.
func JmodMapFetch(className string) string {
    jmodMapMutex.Lock()   // Wait if the map still being built by initialisation.
    jmodMapMutex.Unlock() // Immediately unlock.
    if jmodMapSize == 0 {
        msg := fmt.Sprintf("JmodMapFetch: JMODMAP size = 0 detected when key=%s", className)
        _ = log.Log(msg, log.SEVERE)
        shutdown.Exit(shutdown.JVM_EXCEPTION)
    }
    jmodFile := JMODMAP[className+".class"]
    // fmt.Printf("$$$$$$$$$$$$$$$$$$$$$$$$$$$$ DEBUG key={%s}, jmod={%s}\n", className, jmodFile)
    return jmodFile
}

// This function returns the number of entries in JMODMAP.
func JmodMapSize() int {
    return jmodMapSize
}

// This function returns the number of entries in JMODMAP.
func JmodMapFoundGob() bool {
    return jmodMapFoundGob
}

// This function initializes JMODMAP and jmodMapSize.
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
        msg := fmt.Sprintf("JmodMapInit: os.Open(%s) failed", global.JacobinHome)
        _ = log.Log(msg, log.SEVERE)
        _ = log.Log(err.Error(), log.SEVERE)
        jmodMapSize = 0
        return
    }
    msg := fmt.Sprintf("JmodMapInit: JacobinHome is %s", global.JacobinHome)
    _ = log.Log(msg, logLevel)

    // Get all the file entries in the JacobinHome directory
    names, err := dirOpened.Readdirnames(0) // get all entries
    if err != nil {
        msg := fmt.Sprintf("JmodMapInit: Readdirnames(%s) failed", global.JacobinHome)
        _ = log.Log(msg, log.SEVERE)
        _ = log.Log(err.Error(), log.SEVERE)
        jmodMapSize = 0
        return
    }

    // For each JacobinHome file, try to find a matching gob file
    for ix := range names {
        name := names[ix]
        // fmt.Printf("DEBUG name = %s\n", name)
        if strings.HasSuffix(name, ".gob") { // Gob file?
            version := strings.TrimSuffix(name, ".gob") // get rid of trailing .gom
            if version == global.JavaVersion {
                // Got a match!  Build map from it.
                gobFullPath := global.JacobinHome + string(os.PathSeparator) + name
                msg := fmt.Sprintf("JmodMapInit: Gob file %s selected", gobFullPath)
                _ = log.Log(msg, logLevel)
                if !buildMapFromGob(gobFullPath) {
                    // Gob file trouble
                    // Force re-creation
                    break
                }

                // Map built form gob file succeeded
                jmodMapFoundGob = true
                return
            }
        }
    }

    // No matching gob file
    msg = fmt.Sprintf("JmodMapInit: No gob files matched Java version %s", global.JavaVersion)
    _ = log.Log(msg, logLevel)
    jmodMapFoundGob = false
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
        msg := fmt.Sprintf("buildMapFromGob: os.Open(%s) failed", gobFilePath)
        _ = log.Log(msg, log.WARNING)
        _ = log.Log(err.Error(), log.WARNING)
        jmodMapSize = 0
        return false
    }
    defer inFile.Close()

    // Create a decoder and receive a value.
    decoder := gob.NewDecoder(inFile)
    err = decoder.Decode(&JMODMAP)
    if err != nil {
        msg := fmt.Sprintf("buildMapFromGob: gob Decode(%s) failed", gobFilePath)
        _ = log.Log(msg, log.WARNING)
        _ = log.Log(err.Error(), log.WARNING)
        jmodMapSize = 0
        return false
    }

    gobSize := JMODMAP[counterElementName]
    jmodMapSize, err = strconv.Atoi(gobSize)
    if err != nil {
        msg := fmt.Sprintf("buildMapFromGob: Element (%s) is missing or misformatted", counterElementName)
        _ = log.Log(msg, log.WARNING)
        _ = log.Log(err.Error(), log.WARNING)
        jmodMapSize = 0
        return false
    }
    msg := fmt.Sprintf("buildMapFromGob: Map size from gob file = %d", jmodMapSize)
    _ = log.Log(msg, logLevel)

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
        msg := fmt.Sprintf("buildMapFromJmods: os.Open(%s) failed", dirPath)
        _ = log.Log(msg, log.SEVERE)
        _ = log.Log(err.Error(), log.SEVERE)
        jmodMapSize = 0
        return
    }

    // Get all the file entries in the jmods directory
    names, err := dirOpened.Readdirnames(0) // get all entries
    if err != nil {
        _ = log.Log("buildMapFromJmods: Readdirnames(jmods directory) failed", log.SEVERE)
        _ = log.Log(err.Error(), log.SEVERE)
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
        classFileName := strings.Replace(fileEntry.Name, "classes/", "", 1)

        // Add to map
        JMODMAP[classFileName] = jmodFileName
        // fmt.Printf("DEBUG processJmodFile: classFileName=%s, jmodFileName=%s\n", classFileName, jmodFileName)

        // Add to count of classes
        countClasses++

        // Add to size of JMODMAP
        jmodMapSize++
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
    _ = os.Remove(gobFile)
    outFile, err := os.Create(gobFile)
    if err != nil {
        msg := fmt.Sprintf("saveMapToGob: os.Create(%s) failed", gobFile)
        _ = log.Log(msg, log.SEVERE)
        _ = log.Log(err.Error(), log.SEVERE)
        jmodMapSize = 0
        return
    }

    // Create a gob encoder and encode the cross-reference map
    encoder := gob.NewEncoder(outFile)
    err = encoder.Encode(JMODMAP)
    if err != nil {
        msg := fmt.Sprintf("saveMapToGob: gob Encode(%s) failed", gobFile)
        _ = log.Log(msg, log.SEVERE)
        _ = log.Log(err.Error(), log.SEVERE)
        jmodMapSize = 0
        return
    }

    // Close the output file
    err = outFile.Close()
    if err != nil {
        msg := fmt.Sprintf("saveMapToGob: close(%s) failed", gobFile)
        _ = log.Log(msg, log.SEVERE)
        _ = log.Log(err.Error(), log.SEVERE)
        jmodMapSize = 0
    }

}
