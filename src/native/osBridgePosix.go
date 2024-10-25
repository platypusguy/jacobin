//go:build !windows

package native

import (
	"fmt"
	"github.com/ebitengine/purego"
	"jacobin/trace"
)

func ConnectLibrary(libPath string) uintptr {
	var handle uintptr
	var err error
	handle, err = purego.Dlopen(libPath, purego.RTLD_LAZY)
	if err != nil {
		errMsg := fmt.Sprintf("ConnectLibrary: purego.Dlopen for [%s] failed, reason: [%s]",
			libPath, err.Error())
		trace.ErrorMsg(errMsg)
		handle = 0
	}
	return handle
}
