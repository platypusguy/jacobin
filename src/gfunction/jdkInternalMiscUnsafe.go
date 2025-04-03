/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/trace"
	"reflect"
	"strings"
	"unsafe"
)

/*
 Each object or library that has Go methods contains a reference to MethodSignatures,
 which contain data needed to insert the go method into the MTable of the currently
 executing JVM. MethodSignatures is a map whose key is the fully qualified name and
 type of the method (that is, the method's full signature) and a value consisting of
 a struct of an int (the number of slots to pop off the caller's operand stack when
 creating the new frame and a function. All methods have the same signature, regardless
 of the signature of their Java counterparts. That signature is that it accepts a slice
 of interface{} and returns an interface{}. The accepted slice can be empty and the
 return interface can be nil. This covers all Java functions. (Objects are returned
 as a 64-bit address in this scheme (as they are in the JVM).

 The passed-in slice contains one entry for every parameter passed to the method (which
 could mean an empty slice).
*/

func Load_Jdk_Internal_Misc_Unsafe() {

	MethodSignatures["jdk/internal/misc/Unsafe.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["jdk/internal/misc/Unsafe.unsafeArrayBaseOffset(Ljava/lang/Class;)I"] = // offset to start of first item in an array
		GMeth{
			ParamSlots: 1,
			GFunction:  unsafeArrayBaseOffset,
		}

	MethodSignatures["jdk/internal/misc/Unsafe.unsafeArrayIndexScale(Ljava/lang/Class;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  unsafeArrayIndexScale,
		}

	MethodSignatures["jdk/internal/misc/Unsafe.unsafeArrayIndexScale0(Ljava/lang/Class;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  unsafeArrayIndexScale0,
		}

	MethodSignatures["jdk/internal/misc/Unsafe.compareAndSetInt(Ljava/lang/Object;JII)Z"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  unsafeCompareAndSetInt,
		}

	MethodSignatures["jdk/internal/misc/Unsafe.getAndAddInt(Ljava/lang/Object;JI)I"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  unsafeCompareAndSetInt,
		}

	MethodSignatures["jdk/internal/misc/Unsafe.getIntVolatile(Ljava/lang/Object;J)I"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  unsafeGetIntVolatile,
		}

	MethodSignatures["jdk/internal/misc/Unsafe.getLong(Ljava/lang/Object;J)J"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  unsafeGetLong,
		}

	MethodSignatures["jdk/internal/misc/Unsafe.getUnsafe()Ljdk/internal/misc/Unsafe;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  unsafeGetUnsafe,
		}

	MethodSignatures["jdk/internal/misc/Unsafe.objectFieldOffset1(Ljava/lang/Class;Ljava/lang/String;)J"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  unsafeObjectFieldOffset1,
		}

}

var classUnsafeName = "jdk/internal/misc/Unsafe"

// Return the number of bytes between the beginning of the object and the first element.
// This is used in computing the pointer to a given element
// "jdk/internal/misc/Unsafe.unsafeArrayBaseOffset(Ljava/lang/Class;)I"
func unsafeArrayBaseOffset(params []interface{}) interface{} {
	p := params[0]
	if p == nil || p == object.Null {
		errMsg := "unsafeArrayBaseOffset: Object is a null pointer"
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}
	return int64(0) // this should work...
}

// Return the size of the elements of an array
func unsafeArrayIndexScale(params []interface{}) interface{} {
	arrObj := params[0] // array class whose scale factor is to be returned
	if arrObj == object.Null {
		errMsg := "unsafeArrayIndexScale: Object is a null pointer"
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}

	return unsafeArrayIndexScale0(params)
}

// Utility function that does the work of Unsafe.unsafeArrayIndexScale()
func unsafeArrayIndexScale0(params []interface{}) interface{} {
	// The array class is passed in as a string, so we need to convert it to an object
	// to get the class name.
	arrClass := params[0].(*object.Object).FieldTable["value"].Ftype
	if strings.HasPrefix(arrClass, "[[") { // multi-dimensional array, the first dimension is always pointers
		return int64(8)
	}

	switch arrClass {
	case "[Z", "[B":
		return int64(1)
	default:
		return int64(8)
	}
}

// SWAG
// "jdk/internal/misc/Unsafe.getIntVolatile(Ljava/lang/Object;J)I"
func unsafeGetIntVolatile(params []interface{}) interface{} {
	var hash int64
	switch params[1].(type) {
	case nil:
		hash = 0
	case *object.Object:
		obj := params[0].(*object.Object)
		hash = int64(obj.Mark.Hash)
	}
	offset := params[2].(int64)
	wint := hash + offset
	return wint
}

// SWAG
// "jdk/internal/misc/Unsafe.compareAndSetInt(Ljava/lang/Object;JII)Z"
func unsafeCompareAndSetInt(params []interface{}) interface{} {
	return int64(1) // SWAG
}

// SWAG
// "jdk/internal/misc/Unsafe.getAndAddInt(Ljava/lang/Object;JI)I"
func unsafeGetAndAddInt(params []interface{}) interface{} {
	var hash int64
	switch params[1].(type) {
	case nil:
		hash = 0
	case *object.Object:
		obj := params[0].(*object.Object)
		hash = int64(obj.Mark.Hash)
	}
	offset := params[2].(int64)
	delta := params[3].(int64)
	wint := hash + offset + delta
	return wint
}

func unsafeGetUnsafe([]interface{}) interface{} {
	obj := object.MakeEmptyObjectWithClassName(&classUnsafeName)
	return obj
}

func unsafeObjectFieldOffset1([]interface{}) interface{} {
	return int64(0)
}

func unsafeGetLong(params []interface{}) interface{} {
	obj, ok := params[1].(*object.Object)
	if !ok {
		trace.Warning("unsafeGetLong: Not an object, returning 0")
		return int64(0)
	}
	offset, ok := params[2].(int64)
	if !ok {
		trace.Warning("unsafeGetLong: Invalid offset, returning 0")
		return int64(0)
	}

	// Get the reflect.Value of obj.
	value := reflect.ValueOf(obj)

	// Ensure that the reflect value is addressable.
	if value.Kind() != reflect.Ptr {
		trace.Warning("unsafeGetLong: Object must be a pointer, returning 0")
		return int64(0)
	}

	// Get the unsafe pointer to the object.
	ptr := unsafe.Pointer(value.Pointer())

	// Compute the target memory location.
	target := unsafe.Pointer(uintptr(ptr) + uintptr(offset))

	// Read the int64 value at the computed address and hope for no ka-boom!
	return *(*int64)(target)
}
