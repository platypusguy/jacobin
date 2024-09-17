package native

import (
	"fmt"
	"github.com/ebitengine/purego"
	"jacobin/globals"
	"jacobin/log"
	"runtime"
)

func nativeInit() bool {

	_ = log.Log("nativeInit: Begin", log.TRACE_INST)

	// Set up library file extension and library path string as a function of O/S.
	OperSys = runtime.GOOS
	switch OperSys {
	case "darwin":
		FileExt = "dylib"
	case "linux":
		FileExt = "so"
	case "windows":
		FileExt = "dll"
		WindowsOS = true
	default:
		errMsg := fmt.Sprintf("nativeInit: Unsupported O/S: %s", OperSys)
		_ = log.Log(errMsg, log.SEVERE)
		return false
	}

	// Calculate some needed paths.
	if WindowsOS {
		PathDirLibs = globals.JavaHome() + SepPathString + "bin"
		PathLibjvm = PathDirLibs + SepPathString + "server" + SepPathString + "jvm." + FileExt
		PathLibjava = PathDirLibs + SepPathString + "java." + FileExt
	} else {
		PathDirLibs = globals.JavaHome() + SepPathString + "lib"
		PathLibjvm = PathDirLibs + SepPathString + "server" + SepPathString + "libjvm." + FileExt
		PathLibjava = PathDirLibs + SepPathString + "libjava." + FileExt
	}

	// Connect to libjvm.
	HandleLibjvm = ConnectLibrary(PathLibjvm)
	if HandleLibjvm == 0 {
		return false
	}
	infoMsg := fmt.Sprintf("nativeInit: connect to %s ok", PathLibjvm)
	_ = log.Log(infoMsg, log.TRACE_INST)

	// Connect to libjava.
	HandleLibjava = ConnectLibrary(PathLibjava)
	if HandleLibjvm == 0 {
		return false
	}
	infoMsg = fmt.Sprintf("nativeInit: connect to %s ok", PathLibjava)
	_ = log.Log(infoMsg, log.TRACE_INST)

	// Register the JVM creator library function.
	funcName := "JNI_CreateJavaVM"
	var JvmEnv uintptr
	var createJvm func(*uintptr, *uintptr, *t_JavaVMInitArgs) NFint // (& ptr to JVM, & ptr to env, & arguments) returns JNIint
	purego.RegisterLibFunc(&createJvm, HandleLibjvm, funcName)
	infoMsg = fmt.Sprintf("nativeInit: purego.RegisterLibFunc (%s) ok", funcName)
	_ = log.Log(infoMsg, log.TRACE_INST)

	// Create the JVM.

	ret := createJvm(&HandleJVM, &JvmEnv, &JavaVMInitArgs)
	if ret < 0 {
		_ = log.Log("nativeInit: Cannot create a JVM. Exiting.", log.SEVERE)
		return false
	}
	_ = log.Log("nativeInit: createJvm ok", log.TRACE_INST)

	// Register the GetEnv library function.
	funcName = "JNU_GetEnv"
	var getEnv func(uintptr, *uintptr, NFint) NFint // (ptr to JVM, & ptr to env,JNI version) returns JNIint
	purego.RegisterLibFunc(&getEnv, HandleLibjava, funcName)
	infoMsg = fmt.Sprintf("nativeInit: purego.RegisterLibFunc (%s) ok", funcName)
	_ = log.Log(infoMsg, log.TRACE_INST)

	// Get the JNI environment pointer for the current thread.
	ret = getEnv(HandleJVM, &HandleENV, JavaVMInitArgs.version)
	if ret < 0 {
		_ = log.Log("nativeInit: Cannot get the JNI environment pointer. Exiting.", log.SEVERE)
		return false
	}
	_ = log.Log("nativeInit: End, got JNI env handle", log.TRACE_INST)

	return true

}
