/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2021-5 by the Jacobin authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package globals

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"jacobin/src/config"
	"jacobin/src/types"
	"os"
	"os/exec"
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
	JavaHome        string
	JavaVersion     string
	JDKmajorVersion int // the major version of the JDK, e.g. 11 for Java 11

	// ---- Jacobin Home ----
	JacobinHome string

	// ---- thread management ----
	// Threads ThreadList // list of all app execution threads
	ThreadLock sync.RWMutex
	//    the interface{} values in these maps are object references, but
	//    due to circularity, they need to be described this way here.
	Threads      map[int]interface{}
	ThreadGroups map[string]interface{}
	TGLock       sync.RWMutex

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
	FuncGetStaticValue   func(string, string) any
	FuncInstantiateClass func(string, *list.List) (any, error)
	FuncInvokeGFunction  func(string, []any) any
	FuncMinimalAbort     func(int, string)
	FuncRunThread        func([]any)
	FuncThrowException   func(int, string) bool
	FuncFillInStackTrace func([]any) any
}

// ---- JJ options
var Galt bool // gfunction alternative processing flag -- used strictly for testing

// ---- trace categories
var TraceInit bool
var TraceCloadi bool
var TraceInst bool
var TraceClass bool
var TraceVerbose bool

// ----- String Pool
type StringPoolTable_t map[string]uint32

var StringPoolTable StringPoolTable_t
var StringPoolList []string
var StringPoolNext uint32
var StringPoolLock sync.RWMutex

// ----- map of java/lang/Class instances for statics and instrospection
var JLCmap map[string]any // map of FQN class names to their java/lang/Class instance
var JlcMapLock sync.RWMutex

// LoaderWg is a wait group for various channels used for parallel loading of classes.
var LoaderWg sync.WaitGroup

// Standard Sleep amount in milliseconds used in various places.
var SleepMsecs time.Duration = 5

// the Globals struct.
var global Globals

// Running on Windows?
var OnWindows = runtime.GOOS == "windows"

// InitGlobals initializes the global values that are known at start-up
func InitGlobals(progName string) Globals {

	global = Globals{ // in alpha order
		ArrayAddressList:     InitArrayAddressList(),
		Classpath:            make([]string, 1), // at least one element, the current directory
		ClasspathRaw:         "",
		ErrorGoStack:         "",
		ExitNow:              false,
		FileEncoding:         "UTF-8", // default encoding for file contents
		FileNameEncoding:     "UTF-8", // default encoding for file names
		FuncInstantiateClass: fakeInstantiateClass,
		FuncInvokeGFunction:  fakeInvokeGFunction,
		FuncMinimalAbort:     fakeMinimalAbort,
		FuncRunThread:        fakeRunThread,
		FuncThrowException:   fakeThrowEx,
		GoStackShown:         false,
		JacobinBuildData:     nil,
		JacobinHome:          "",
		JacobinName:          progName,
		JavaHome:             "",
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
		Version:              config.GetJacobinVersion(), // gets version and build #
		VmModel:              "server",
	}

	// ----- map of java/lang/Class instances for statics and instrospection
	JLCmap = make(map[string]any, 2000) // map of FQN class name to its java/lang/Class instance

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
		if progName != "test" {
			os.Exit(1)
		}
	}

	verNumber, verString := GetJDKmajorVersion() // JDK version depends on JAVA_HOME, so place here
	global.JavaVersion = verString
	global.JDKmajorVersion = verNumber

	InitJacobinHome()
	if global.JacobinHome == "" {
		if progName != "test" {
			os.Exit(1)
		}
	}
	InitArrayAddressList()

	if OnWindows {
		global.FileEncoding = "windows-1252"
	} else {
		global.FileEncoding = "UTF-8"
	}

	// Make the encoding for filesystem names be the same as for file contents.
	global.FileNameEncoding = global.FileEncoding

	// Set up headlass mode for AWT from enviromment varialbe. May eventually be removed as we don't support AWT.
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
	global.ThreadGroups = make(map[string]interface{})

	return global
}

// InitClasspath initializes the classpath from the CLASSPATH environment variable.
// If CLASSPATH is not set, it uses the current working directory as the classpath.
// This will be overriden by the -cp or -classpath command-line options, if any.
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
		_, _ = fmt.Fprintf(os.Stderr, "InitJavaHome: Cannot find the specified path: %s. Exiting.\n", releasePath)
		return
	}
	defer handle.Close()

	// Get the JDK major version and the full version string from the release file.
	major, versionString := GetJDKmajorVersion()
	if major == 0 {
		global.JavaVersion = "" // Set to empty string to signal failure
		global.JDKmajorVersion = 0
		_, _ = fmt.Fprintf(os.Stderr, "InitJavaHome: Did not find the JAVA_VERSION record in %s. Exiting.\n",
			releasePath)
		return
	} else {
		global.JavaVersion = versionString // Set the Java version from the release file
		global.JDKmajorVersion = major
	}
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

// Fake RunThread() in run.go
func fakeRunThread(_ []interface{}) {
	errMsg := fmt.Sprintf("\n*Attempt to access uninitialized RunThread pointer func\n")
	fmt.Fprintf(os.Stderr, "%s", errMsg)
}

// Fake ThrowEx() in exceptions.go
func fakeThrowEx(whichEx int, msg string) bool {
	errMsg := fmt.Sprintf("\n*Attempt to access uninitialized ThrowEx pointer func\n")
	fmt.Fprintf(os.Stderr, "%s", errMsg)
	return false
}

// Fake Invoke in gfunctions.go
func fakeInvokeGFunction(name string, args []any) any {
	errMsg := fmt.Sprintf("\n*Attempt to access uninitialized InvokeGFunction pointer func: name=%s\n", name)
	fmt.Fprintf(os.Stderr, "%s", errMsg)
	return nil
}

func InitStringPool() {

	StringPoolLock.Lock()
	defer StringPoolLock.Unlock()
	StringPoolTable = make(StringPoolTable_t)
	// Prestored values: 0 = nil, 1 = java/lang/String, 2 = java/lang/Object,
	// 3 = java/lang/Thread, 4 = java/lang/Class
	StringPoolTable[""] = 0
	StringPoolTable["java/lang/String"] = types.StringPoolStringIndex
	StringPoolTable["java/lang/Object"] = types.StringPoolObjectIndex
	StringPoolTable["java/lang/Thread"] = types.StringPoolThreadIndex
	StringPoolTable["java/lang/Class"] = types.StringPoolJavaLangClassIndex

	// Pre-stored string list. This list must exactly match the preceding list.
	StringPoolList = nil
	StringPoolList = append(StringPoolList, types.EmptyString)
	StringPoolList = append(StringPoolList, types.StringClassName)
	StringPoolList = append(StringPoolList, types.ObjectClassName)
	StringPoolList = append(StringPoolList, types.ClassNameThread)
	StringPoolList = append(StringPoolList, types.ClassNameJavaLangClass)

	// Set up the next available index.
	StringPoolNext = uint32(len(StringPoolList))

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

// Reads the JDK release file and returns the major version number and the full version string.
func GetJDKmajorVersion() (int, string) {
	releaseFilePath := global.JavaHome + string(os.PathSeparator) + "release"
	file, err := os.Open(releaseFilePath)
	if err != nil {
		if TraceVerbose {
			fmt.Fprintf(os.Stderr, "GetJDKversion(): open release file failed: %v\n", err)
		}
		return 0, ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "JAVA_VERSION=") {
			// Extract the value after "JAVA_VERSION="
			versionString := strings.Trim(line[len("JAVA_VERSION="):], "\"")
			parts := strings.Split(versionString, ".")
			majorVersion, _ := strconv.Atoi(parts[0])
			return majorVersion, versionString // Return major version number as an int and version string
		}
	}

	if err = scanner.Err(); err != nil {
		if TraceVerbose {
			fmt.Fprintf(os.Stderr, "error reading release file in getJDKversion(): %v", err)
		}
		return 0, ""
	}

	if TraceVerbose {
		fmt.Fprintf(os.Stderr, "getJDKversion(): JAVA_VERSION not found in release file")
	}
	return 0, ""
}
