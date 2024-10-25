package native

import (
	"fmt"
	"jacobin/globals"
	"jacobin/trace"
	"runtime"
)

func nativeInit() bool {

	if globals.TraceInit {
		trace.Trace("nativeInit: Begin")
	}

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
		trace.ErrorMsg(errMsg)
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
	if globals.TraceInit {
		infoMsg := fmt.Sprintf("nativeInit: End, connected to %s", PathLibjvm)
		trace.Trace(infoMsg)
	}

	return true

}
