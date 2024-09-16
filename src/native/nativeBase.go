/*

During initialization,
* The NfLibXrefTable is built by either a POSIX loader or a Windows loader. Note that both the library path and handle are populated.
* The nfToTmplTable remains nil.

At run-time, RunNativeFunction will do the following in order to get (1) a native function handle
and (2) the corresponding template function address:
* Look up the funcName in the nfToTmplTable.
* If not found,
     - Look up funcName in nfToLibTable. Not found ---> error.
     - Derive the template function to use for this methodName based on the methodType.
     - Store the template function handle in nfToTmplTable.
* Call the template function (by address) with arguments: library handle and the function name.

*/

package native

import "errors"

var CaughtNativeFunctionException = errors.New("caught native function exception")

// Native function to library handle table
// Input: Native function name.
// Output: Library file handle.
var nfToLibTable = map[string]uintptr{}

// Native function to template function handle table
// Input: Native function name.
// Output: Template function handle.
var nfToTmplTable = map[string]typeTemplateFunction{}

// Native function error block.
type NativeErrBlk struct {
	ExceptionType int
	ErrMsg        string
}

// Type definition for all the template functions
type typeTemplateFunction func(libHandle uintptr, functionName string, params []interface{}) interface{}

// Argument types for template functions.
type NFboolean uint8
type NFbyte uint8
type NFchar uint16
type NFshort int16
type NFint int32
type NFlong int64
type NFfloat float32
type NFdouble float64
