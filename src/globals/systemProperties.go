/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package globals

import (
	"fmt"
	"jacobin/types"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// the System Properties Map: JVM System Properties
//
// Jacobin uses information from the operating system, startup arguments, and the host environment to
// set up its initial system properties.
//
// Properties are stored in the globalPropertiesMap, fetched with System.getProperties() or
// System.getProperty(key), and include items such as:
// * os.name, os.arch, os.version
// * user.name, user.home, user.dir
// * java.home, java.class.path
// * file.separator, line.separator
//
// These values are derived from:
// * Environment variables (HOME, PATH, etc.)
// * The current working directory
// * Command-line -D options passed when launching the JVM
// * Other means

var systemPropertiesMap types.DefProperties
var systemPropertiesMutex = sync.RWMutex{}

func getOsProperty(arg string) string {
	var value string
	operSys := runtime.GOOS

	switch arg {
	case "file.encoding":
		value = global.FileEncoding
	case "file.separator":
		value = string(os.PathSeparator)
	case "java.class.path":
		value = global.ClasspathRaw
	case "java.compiler": // the name of the JIT compiler (we don't have a JIT)
		value = "no JIT"
	case "java.home":
		value = global.JavaHome
	case "java.io.tmpdir":
		value = os.TempDir()
	case "java.library.path":
		value = global.JavaHome
	case "java.vendor":
		value = "Jacobin"
	case "java.vendor.url":
		value = "https://jacobin.org"
	case "java.vendor.version":
		value = global.Version
	case "java.version":
		_, versionString := GetJDKmajorVersion()
		value = versionString
	// case "java.version.date":
	// 	need to get this
	case "java.vm.name":
		value = fmt.Sprintf(
			"Jacobin VM v. %s (Java %d) 64-bit VM", global.Version, global.MaxJavaVersion)
	case "java.vm.specification.name":
		value = "Java Virtual Machine Specification"
	case "java.vm.specification.vendor":
		value = "Oracle and Jacobin"
	case "java.vm.specification.version":
		value = strconv.Itoa(global.MaxJavaVersion)
	case "java.vm.vendor":
		value = "Jacobin"
	case "java.vm.version":
		value = strconv.Itoa(global.MaxJavaVersion)
	case "jdk.major.version":
		_, ver := GetJDKmajorVersion() // "" if not found
		value = ver
	case "line.separator":
		if operSys == "windows" {
			value = "\\r\\n"
		} else {
			value = "\\n"
		}
	case "native.encoding", "stdout.encoding", "stderr.encoding":
		value = GetCharsetName()
	case "os.arch":
		value = runtime.GOARCH
	case "os.name":
		value = operSys
	case "os.version":
		value = getOSVersion()
	case "path.separator":
		value = string(os.PathSeparator)
	case "sun.jnu.encoding":
		value = "UTF-8" // this is the default encoding for file names in Java
	case "user.dir": // present working directory
		value, _ = os.Getwd()
	case "user.home":
		currentUser, _ := user.Current()
		value = currentUser.HomeDir
	case "user.name":
		currentUser, _ := user.Current()
		value = currentUser.Name
	case "user.timezone":
		now := time.Now()
		value, _ = now.Zone()
	default:
		value = ""
	}

	return value
}

// Build the Global Properties Map.
func buildGlobalProperties() {
	systemPropertiesMap = make(types.DefProperties)
	systemPropertiesMutex.Lock()
	defer systemPropertiesMutex.Unlock()

	systemPropertiesMap["file.encoding"] = getOsProperty("file.encoding")
	systemPropertiesMap["file.separator"] = getOsProperty("file.separator")
	systemPropertiesMap["java.class.path"] = "."
	systemPropertiesMap["java.compiler"] = getOsProperty("java.compiler")
	systemPropertiesMap["java.home"] = getOsProperty("java.home")
	systemPropertiesMap["java.io.tmpdir"] = getOsProperty("java.io.tmpdir")
	systemPropertiesMap["java.library.path"] = getOsProperty("java.library.path")
	systemPropertiesMap["java.vendor"] = getOsProperty("java.vendor")
	systemPropertiesMap["java.vendor.url"] = getOsProperty("java.vendor.url")
	systemPropertiesMap["java.vendor.version"] = getOsProperty("java.vendor.version")
	systemPropertiesMap["java.version"] = getOsProperty("java.version")
	systemPropertiesMap["java.vm.name"] = getOsProperty("java.vm.name")
	systemPropertiesMap["java.vm.specification.name"] = getOsProperty("java.vm.specification.name")
	systemPropertiesMap["java.vm.specification.vendor"] = getOsProperty("java.vm.specification.vendor")
	systemPropertiesMap["java.vm.specification.version"] = getOsProperty("java.vm.specification.version")
	systemPropertiesMap["java.vm.vendor"] = getOsProperty("java.vm.vendor")
	systemPropertiesMap["java.vm.version"] = getOsProperty("java.vm.version")
	systemPropertiesMap["jdk.major.version"] = getOsProperty("java.version")
	systemPropertiesMap["line.separator"] = getOsProperty("line.separator")
	systemPropertiesMap["native.encoding"] = getOsProperty("native.encoding")
	systemPropertiesMap["os.arch"] = getOsProperty("os.arch")
	systemPropertiesMap["os.name"] = getOsProperty("os.name")
	systemPropertiesMap["os.version"] = getOsProperty("os.version")
	systemPropertiesMap["path.separator"] = getOsProperty("path.separator")
	systemPropertiesMap["stdout.encoding"] = getOsProperty("stdout.encoding")
	systemPropertiesMap["stderr.encoding"] = getOsProperty("stderr.encoding")
	systemPropertiesMap["sun.jnu.encoding"] = "UTF-8"
	systemPropertiesMap["user.dir"] = getOsProperty("user.dir")
	systemPropertiesMap["user.home"] = getOsProperty("user.home")
	systemPropertiesMap["user.name"] = getOsProperty("user.name")
	systemPropertiesMap["user.timezone"] = getOsProperty("user.timezone")
}

// GetSystemProperty: get a system property.
func GetSystemProperty(key string) string {
	return systemPropertiesMap[key]
}

// SetSystemProperty: add or update a system property.
func SetSystemProperty(key, value string) {
	systemPropertiesMutex.Lock()
	defer systemPropertiesMutex.Unlock()
	systemPropertiesMap[key] = value
}

// RemoveSystemProperty: remove a system property.
func RemoveSystemProperty(key string) {
	systemPropertiesMutex.Lock()
	defer systemPropertiesMutex.Unlock()
	delete(systemPropertiesMap, key)
}

// ReplaceSystemProperties: replace the current map with a new one.
func ReplaceSystemProperties(newMap types.DefProperties) {
	systemPropertiesMutex.Lock()
	defer systemPropertiesMutex.Unlock()
	systemPropertiesMap = newMap
}
