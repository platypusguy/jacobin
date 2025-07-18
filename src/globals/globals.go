/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-2 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package globals

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"jacobin/config"
	"jacobin/types"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

var StringEnvVarHeadless = "java.awt.headless"

// Globals contains variables that need to be globally accessible,
// such as VM and program args, etc.
//
// Note: to avoid circularity, globals cannot depend on exec package.
// As a result, exec contains its own globals
type Globals struct {
	// ---- jacobin version number ----
	// note: all references to version number must come from this literal
	Version string
	VmModel string // "client" or "server" (both the same acc. to JVM docs)

	// ---- processing stoppage? ----
	ExitNow bool

	// ---- command-line items ----
	JacobinName string // name of the executing Jacobin executable
	Args        []string
	CommandLine string

	StartingClass string
	StartingJar   string
	AppArgs       []string
	Options       map[string]Option

	// ---- classloading items ----
	MaxJavaVersion    int // the Java version as commonly known, i.e. Java 11
	MaxJavaVersionRaw int // the Java version as it appears in bytecode i.e., 55 (= Java 11)
	VerifyLevel       int
	ClasspathRaw      string   // the raw classpath as passed in by the user
	Classpath         []string // the classpath as a list of directories and JARs

	// ---- Java Home and Version ----
	JavaHome    string
	JavaVersion string

	// ---- Jacobin Home ----
	JacobinHome string

	// ---- thread management ----
	// Threads ThreadList // list of all app execution threads
	ThreadLock sync.Mutex
	Threads    map[int]interface{} // in reality the interface is a threads.ExecThread, but
	// due to circularity has to be described this way here.
	ThreadNumber int

	// ---- execution context ----
	JacobinBuildData map[string]string

	// ---- special switches ----
	StrictJDK bool // hew closely to actions and error messages of the JDK

	// ---- list of addresses of arrays, see jvm/arrays.go for info ----
	ArrayAddressList *list.List

	// ----- Byte cache for java.base.jmod
	JmodBaseBytes []byte

	// ----- Error handling
	ErrorGoStack       string
	JVMframeStack      *[]string
	PanicCauseShown    bool
	JvmFrameStackShown bool
	GoStackShown       bool

	// Random object mutex
	RandomLock sync.Mutex

	// AtomicInteger mutex
	AtomicIntegerLock sync.Mutex

	// ---- misc properties
	FileEncoding     string // what file encoding are we using?
	FileNameEncoding string // System.getProperty("sun.jnu.encoding")
	Headless         bool   // Headless?

	// Get around the golang circular dependency. To be set up in jvmStart.go
	// Enables gfunctions to call these functions through a global variable.
	FuncInstantiateClass func(string, *list.List) (any, error)
	FuncMinimalAbort     func(int, string)
	FuncThrowException   func(int, string) bool
	FuncFillInStackTrace func([]any) any
}

// ---- JJ options
var Galt bool

// ---- trace categories
var TraceInit bool
var TraceCloadi bool
var TraceInst bool
var TraceClass bool
var TraceVerbose bool

// ----- String Pool
var StringPoolTable map[string]uint32
var StringPoolList []string
var StringPoolNext uint32
var StringPoolLock sync.Mutex
var StringIndexString uint32

// LoaderWg is a wait group for various channels used for parallel loading of classes.
var LoaderWg sync.WaitGroup

// Standard Sleep amount in milliseconds used in various places.
var SleepMsecs time.Duration = 5

// Instantiate the Globals struct.
var global Globals

// InitGlobals initializes the global values that are known at start-up
func InitGlobals(progName string) Globals {
	global = Globals{ // in alpha order
		ArrayAddressList:     InitArrayAddressList(),
		Classpath:            make([]string, 1), // at least one element, the current directory
		ClasspathRaw:         "",
		ErrorGoStack:         "",
		ExitNow:              false,
		FuncInstantiateClass: fakeInstantiateClass,
		FuncMinimalAbort:     fakeMinimalAbort,
		FuncThrowException:   fakeThrowEx,
		GoStackShown:         false,
		JacobinBuildData:     nil,
		JacobinHome:          "",
		JacobinName:          progName,
		JavaHome:             "",
		JavaVersion:          "",
		JmodBaseBytes:        nil,
		JVMframeStack:        nil,
		JvmFrameStackShown:   false,
		MaxJavaVersion:       21, // this value and MaxJavaVersionRaw must *always* be in sync
		MaxJavaVersionRaw:    65, // this value and MaxJavaVersion must *always* be in sync
		Options:              make(map[string]Option),
		PanicCauseShown:      false,
		StartingClass:        "",
		StartingJar:          "",
		StrictJDK:            false,
		ThreadNumber:         0,                          // first thread will be numbered 1, as increment occurs prior
		Version:              config.GetJacobinVersion(), // gets version and build #
		VmModel:              "server",
	}

	// ----- G function alternative processing flag
	Galt = false

	// ----- Tracing flags
	TraceInit = false
	TraceCloadi = false
	TraceInst = false
	TraceClass = false
	TraceVerbose = false

	// ----- String Pool and other values
	InitStringPool()

	InitJavaHome()
	if global.JavaHome == "" || global.JavaVersion == "" {
		os.Exit(1)
	}
	InitJacobinHome()
	if global.JacobinHome == "" {
		os.Exit(1)
	}
	InitArrayAddressList()

	if runtime.GOOS == "windows" {
		global.FileEncoding = "windows-1252"
	} else {
		global.FileEncoding = "UTF-8"
	}

	// Make the encoding for filesystem names be the same as for file contents.
	global.FileNameEncoding = global.FileEncoding

	// Set up headlass boolean.
	strHeadless := os.Getenv(StringEnvVarHeadless)
	global.Headless = false
	if strHeadless != "" {
		if strHeadless == "true" {
			global.Headless = true
		}
	}

	InitClasspath()

	// Capture system properties from the OS and its environment.
	buildGlobalProperties()

	global.Threads = make(map[int]interface{})

	return global
}

// InitClasspath initializes the classpath from the CLASSPATH environment variable.
// If CLASSPATH is not set, it uses the current working directory as the classpath.
// This will be overriden by the -cp or -classpath command-line options.
func InitClasspath() {
	cp := os.Getenv("CLASSPATH")
	if cp != "" {
		cp = strings.TrimSpace(cp)
		cp = cleanupPath(cp) // convert slashes to current platform's path separator
		global.ClasspathRaw = cp
		global.Classpath = strings.Split(cp, string(os.PathListSeparator))
	} else {
		global.ClasspathRaw, _ = os.Getwd()
		global.Classpath[0] = global.ClasspathRaw
	}
}

// GetGlobalRef returns a pointer to the singleton instance of Globals
func GetGlobalRef() *Globals {
	return &global
}

// Option is the value portion of the globals.options table. This table is described in
// more detail in option_table_loader.go introductory comments
type Option struct {
	Supported bool
	Set       bool
	ArgStyle  int16
	Action    func(position int, name string, gl *Globals) (int, error)
}

// InitJacobinHome gets JACOBIN_HOME and formats it as expected
// Note: any trailing separator is removed from the retrieved string per JACOBIN-184
func InitJacobinHome() {
	jacobinHome := os.Getenv("JACOBIN_HOME")
	if jacobinHome != "" {
		jacobinHome = strings.TrimRight(jacobinHome, "\\/") // remove any trailing separator
		jacobinHome = cleanupPath(jacobinHome)
	} else {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "InitJacobinHome: os.UserHomeDir() failed. Exiting.\n")
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return
		}
		jacobinHome = userHomeDir + string(os.PathSeparator) + "jacobin_data"
	}
	// 0755 (Unix octal): user(owner) can do anything, group and other can read and visit directory ("execute").
	// Ref: https://opensource.com/article/19/8/linux-permissions-101
	err := os.MkdirAll(jacobinHome, 0755)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "InitJacobinHome: os.MkDirAll(%s) failed. Exiting.\n", jacobinHome)
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	// Success!
	global.JacobinHome = jacobinHome
}

func JacobinHome() string { return global.JacobinHome }

// InitJavaHome gets JAVA_HOME from the environment and formats it as expected. (It
// also checks that the directory is valid by looking for the release file.)
// Note: any trailing separator is removed from the retrieved string per JACOBIN-184
func InitJavaHome() {

	javaHome := os.Getenv("JAVA_HOME")
	if javaHome == "" {
		_, _ = fmt.Fprintf(os.Stderr, "InitJavaHome: Environment variable JAVA_HOME missing but is required. Exiting.\n")
		return
	}
	javaHome = strings.TrimRight(javaHome, "\\/") // remove any trailing separator
	javaHome = cleanupPath(javaHome)
	global.JavaHome = javaHome

	// Check if JAVA_HOME is a valid directory by looking for the release file.
	releasePath := javaHome + string(os.PathSeparator) + "release"
	handle, err := os.Open(releasePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "InitJavaHome: os.Open(%s) failed. Exiting.\n", releasePath)
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
	defer handle.Close()
	scanner := bufio.NewScanner(handle)

	// Scan the release file for a mandatory JAVA_VERSION record.
	for scanner.Scan() {
		// do something with a line
		line := scanner.Text()
		tokens := strings.Split(line, "=")
		if len(tokens) != 2 {
			_, _ = fmt.Fprintf(os.Stderr, "InitJavaHome: File format error in %s. Exiting.\n", releasePath)
			return
		}
		if tokens[0] == "JAVA_VERSION" {
			global.JavaVersion = strings.Trim(tokens[1], "\"")
			return
		}
	}

	// At this pint, we did not find a Java version record,
	// so either the JAVA_HOME is not a valid JDK or the release file is corrupted.
	// Either way, we cannot proceed.
	_, _ = fmt.Fprintf(os.Stderr, "InitJavaHome: Did not find the JAVA_VERSION record in %s. Exiting.\n", releasePath)
	os.Exit(1)

}

func JavaHome() string    { return global.JavaHome }
func JavaVersion() string { return global.JavaVersion }

// Normalize a file path. Slashes are converted to the current platform's path separator if necessary.
func cleanupPath(path string) string {
	path = filepath.FromSlash(path)
	return path
}

// Array addresses must be kept in a list to avoid being GC'd.
// This creates that list.
func InitArrayAddressList() *list.List {
	return list.New()
}

// Fake GoStringFromStringObject()
func fakeGoStringFromStringObject(obj interface{}) string {
	errMsg := fmt.Sprintf("\n*Attempt to access uninitialized GoStringFromStringObject pointer func\n")
	fmt.Fprintf(os.Stderr, "%s", errMsg)
	return ""
}

// Fake InstantiateClass
func fakeInstantiateClass(classname string, frameStack *list.List) (any, error) {
	errMsg := fmt.Sprintf("\n*Attempt to access uninitialized InstantiateClass pointer func: classname=%s\n", classname)
	fmt.Fprintf(os.Stderr, "%s", errMsg)
	return nil, errors.New(errMsg)
}

// Fake MinimalAbort() in exceptions.go
func fakeMinimalAbort(whichEx int, msg string) {
	errMsg := fmt.Sprintf("\n*Attempt to access uninitialized MinimalAbort pointer func\n")
	fmt.Fprintf(os.Stderr, "%s", errMsg)
}

// Fake ThrowEx() in exceptions.go
func fakeThrowEx(whichEx int, msg string) bool {
	errMsg := fmt.Sprintf("\n*Attempt to access uninitialized ThrowEx pointer func\n")
	fmt.Fprintf(os.Stderr, "%s", errMsg)
	return false
}

func InitStringPool() {

	StringPoolLock.Lock()

	// create the string pool
	StringPoolTable = make(map[string]uint32)
	StringPoolList = nil

	// Changed on 9-Apr-2024: 0 = nil, 1 = String, 2 = Object
	// Preload two values. java/lang/Object is always 0
	// and java/lang/String is always 1.

	// Add empty string (for when an index field has not been use, and so = 0
	StringPoolTable[""] = 0
	StringPoolList = append(StringPoolList, types.EmptyString)

	// Add "java/lang/String"
	StringPoolTable[types.StringClassName] = types.StringPoolStringIndex
	StringPoolList = append(StringPoolList, types.StringClassName)

	// Add "java/lang/Object"
	StringPoolTable[types.ObjectClassName] = types.ObjectPoolStringIndex
	StringPoolList = append(StringPoolList, types.ObjectClassName)

	// Set up next available index
	StringPoolNext = uint32(3)

	StringPoolLock.Unlock()
}

// Get the character set name.
func GetCharsetName() string {
	return global.FileEncoding
}

// Case-insensitive sort.
// Golang should have provided this!
func SortCaseInsensitive(ptrSlice *[]string) {
	slices.SortFunc(*ptrSlice, func(a, b string) int {
		return strings.Compare(strings.ToLower(a), strings.ToLower(b))
	})
}

// the System Properties Map: JVM System Properties
//
// Jacobin uses information from the operating system, startup arguments, and the host environment to set up its initial system properties.
//
// Properties are stored in the globalPropertiesMap, fetched with System.getProperties() or System.getProperty(key), and include items like:
// * os.name, os.arch, os.version
// * user.name, user.home, user.dir
// * java.home, java.class.path
// * file.separator, line.separator
//
// These values are derived from:
// * Environment variables (HOME, PATH, etc.)
// * The current working directory
// * Command-line -D options passed when launching the JVM

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
		value = "." // OpenJDK JVM default value
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
		value = strconv.Itoa(global.MaxJavaVersion)
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
		value = global.FileNameEncoding
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
	systemPropertiesMap["java.class.path"] = "." // TODO - fix this during CLASSPATH development
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
	systemPropertiesMap["line.separator"] = getOsProperty("line.separator")
	systemPropertiesMap["native.encoding"] = getOsProperty("native.encoding")
	systemPropertiesMap["os.arch"] = getOsProperty("os.arch")
	systemPropertiesMap["os.name"] = getOsProperty("os.name")
	systemPropertiesMap["os.version"] = getOsProperty("os.version")
	systemPropertiesMap["path.separator"] = getOsProperty("path.separator")
	systemPropertiesMap["stdout.encoding"] = getOsProperty("stdout.encoding")
	systemPropertiesMap["stderr.encoding"] = getOsProperty("stderr.encoding")
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

// getOSVersion: Get the O/S version string and return it to caller.
func getOSVersion() string {
	var cmd *exec.Cmd

	operSys := runtime.GOOS
	switch operSys {
	case "windows":
		cmd = exec.Command("cmd", "/C", "ver")
	default:
		cmd = exec.Command("uname", "-r")
	}

	cmdBytes, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("getOSVersion: cmd.CombinedOutput() failed on %s: %v", operSys, err)
		return errMsg
	}

	var cleanBytes []byte
	for ix := 0; ix < len(cmdBytes); ix++ {
		if unicode.IsPrint(rune(cmdBytes[ix])) {
			cleanBytes = append(cleanBytes, cmdBytes[ix])
		}
	}

	return string(cleanBytes)
}
