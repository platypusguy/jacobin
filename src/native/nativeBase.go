/*

During initialization,
* The xrefTable is built by either a POSIX loader or a Windows loader. Note that both the library path and handle are populated.
* The methodsTable remains nil.

At run-time, RunNativeFunction will do the following in order to get a native function handle:
* Form the functionKey = methodName concatenated with methodType.
* Look up the functionKey in the nfuncTable.
* If not found,
     - Look up functionKey in xrefTable.
     - Not found ---> error.
     - Use the libHandle to get the function handle.
     - Failure (E.g. not found) ---> error.
     - Store the function handle in nfuncTable.
* Use the function handle for the function call.

*/

package native

import (
	"fmt"
	"github.com/ebitengine/purego"
	"jacobin/excNames"
)

var nfuncTable = map[string]uintptr{} // Functions encountered and therefore have a handle

type typeNxref struct {
	LibPath   string
	LibHandle uintptr
}

var xrefTable = map[string]typeNxref{} // Function-to-library cross reference table

// Native function error block.
type NativeErrBlk struct {
	ExceptionType int
	ErrMsg        string
}

func getFuncHandle(methodName, methodType string) interface{} {
	functionKey := methodName + methodType
	funcHandle, ok := nfuncTable[functionKey]
	if !ok {
		libHandle := xrefTable[functionKey].LibHandle
		funcHandle, err := purego.Dlsym(libHandle, methodName)
		if err != nil {
			libPath := xrefTable[functionKey].LibPath
			errMsg := fmt.Sprintf("purego.Dlsym(%s : %s) failed, reason: %s", libPath, functionKey, err.Error())
			return NativeErrBlk{ExceptionType: excNames.VirtualMachineError, ErrMsg: errMsg}
		}
		if funcHandle == 0 {
			libPath := xrefTable[functionKey].LibPath
			errMsg := fmt.Sprintf("purego.Dlsym(%s : %s) function not found", libPath, functionKey)
			return NativeErrBlk{ExceptionType: excNames.VirtualMachineError, ErrMsg: errMsg}
		}
	}
	return funcHandle
}
