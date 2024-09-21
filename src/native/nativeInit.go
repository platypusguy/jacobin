package native

import (
	"fmt"
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
	infoMsg := fmt.Sprintf("nativeInit: End, connect to %s ok", PathLibjvm)
	_ = log.Log(infoMsg, log.TRACE_INST)

	return true

}
