/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024-5 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package jvm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/exceptions"
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/stringPool"
	"jacobin/trace"
	"jacobin/types"
	"jacobin/util"
	"math"
	"runtime/debug"
	"strings"
	"unsafe"
)

// This file contains many support functions for the interpreter in run.go.
// These notably include push, pop, and peek operations on the operand stack,
// as well as some formatting functions for tracing, and utility functions for
// conversions of interfaces and data types.

// Convert a byte to an int64 by extending the sign-bit
func byteToInt64(bite byte) int64 {
	if (bite & 0x80) == 0x80 { // Negative bite value (left-most bit on)?
		// Negative byte - need to extend the sign (left-most) bit
		var wbytes = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00}
		wbytes[7] = bite
		// Form an int64 from the wbytes array
		// If you know C, this is equivalent to memcpy(&wint64, &wbytes, 8)
		return int64(binary.BigEndian.Uint64(wbytes))
	}

	// Not negative (left-most bit off) : just cast bite as an int64
	return int64(bite)
}

// converts an interface{} value to int8. Used for BASTORE
func convertInterfaceToByte(val interface{}) types.JavaByte {
	switch t := val.(type) {
	case byte:
		return types.JavaByte(t)
	case int:
		return types.JavaByte(t)
	case int8:
		return t
	case int64:
		return types.JavaByte(t)
	}
	return 0
}

// converts an interface{} value on the op stack into a uint64
func convertInterfaceToUint64(val interface{}) uint64 {
	// in theory, the only types passed to this function are those
	// found on the operand stack: ints, floats, pointers
	switch t := val.(type) {
	case int64:
		return uint64(t)
	case float64:
		return uint64(math.Round(t))
	case unsafe.Pointer:
		intVal := uintptr(t)
		return uint64(intVal)
	}
	return 0
}

// converts an interface{} value into int64
// notes:
//
//   - an exception is thrown if function is passed a uint64
//
//   - there exists a similar uint function: convertInterfaceToUint64()
//
//   - for converting a byte as a numeric value and so to propagate negative values,
//     use byteToInt64(). This would be done only in numeric conversions of binary data.
//
//     It might appear that you could put most of the case statements into a single case statement,
//     but golang does not allow this. For interface conversion, it needs to be done type-by-type.
func convertInterfaceToInt64(arg interface{}) int64 {
	switch t := arg.(type) {
	case int8:
		return int64(t)
	case uint8:
		return int64(t)
	case int16:
		return int64(t)
	case uint16:
		return int64(t)
	case int:
		return int64(t)
	case int32:
		return int64(t)
	case uint32:
		return int64(t)
	case int64:
		return t
	case bool:
		return types.ConvertGoBoolToJavaBool(t)
	default:
		gl := globals.GetGlobalRef()
		if gl.JacobinName != "test" {
			errMsg := fmt.Sprintf("convertInterfaceToInt64: Invalid argument type: %T", arg)
			exceptions.ThrowEx(excNames.InvalidTypeException, errMsg, nil)
		}
	}
	return 0
}

// pop from the operand stack.
func pop(f *frames.Frame) interface{} {
	var value interface{}

	if f.TOS == -1 {
		errMsg := fmt.Sprintf("stack underflow in pop() in %s.%s",
			util.ConvertInternalClassNameToUserFormat(f.ClName), f.MethName)
		status := exceptions.ThrowEx(excNames.InternalException, errMsg, f)
		if status != exceptions.Caught {
			return nil // applies only if in test
		}
	} else {
		value = f.OpStack[f.TOS]
	}

	// we show trace info of the TOS *before* we change its value--
	// all traces show TOS before the instruction is executed.
	if globals.TraceVerbose {
		var traceInfo string
		if f.TOS == -1 {
			traceInfo = fmt.Sprintf("%74s", "POP           TOS:  -")
			trace.Trace(traceInfo)
		} else {
			if value == nil {
				traceInfo = fmt.Sprintf("%74s", "POP           TOS:") +
					fmt.Sprintf("%3d <nil>", f.TOS)
				trace.Trace(traceInfo)
			} else {
				switch value.(type) {
				case *object.Object:
					obj := value.(*object.Object)
					trace.TraceObject(f, "POP", obj)
				case *[]uint8:
					strPtr := value.(*[]byte)
					str := string(*strPtr)
					traceInfo = fmt.Sprintf("%74s", "POP           TOS:") +
						fmt.Sprintf("%3d *[]byte: %-10s", f.TOS, str)
					trace.Trace(traceInfo)
				case []uint8:
					bytes := value.([]byte)
					str := string(bytes)
					traceInfo = fmt.Sprintf("%74s", "POP           TOS:") +
						fmt.Sprintf("%3d []byte: %-10s", f.TOS, str)
					trace.Trace(traceInfo)
				case []types.JavaByte:
					bytes := value.([]types.JavaByte)
					str := object.GoStringFromJavaByteArray(bytes)
					traceInfo = fmt.Sprintf("%74s", "POP           TOS:") +
						fmt.Sprintf("%3d []javaByte: %-10s", f.TOS, str)
					trace.Trace(traceInfo)
				default:
					traceInfo = fmt.Sprintf("%74s", "POP           TOS:") +
						fmt.Sprintf("%3d %T %v", f.TOS, value, value)
					trace.Trace(traceInfo)
				}
			}
		}
	}

	f.TOS -= 1 // adjust TOS
	if globals.TraceVerbose {
		trace.LogTraceStack(f)
	} // trace the resultant stack
	return value
}

// returns the value at the top of the stack without popping it off.
func peek(f *frames.Frame) interface{} {
	if f.TOS == -1 {
		errMsg := fmt.Sprintf("stack underflow in peek() in %s.%s",
			util.ConvertInternalClassNameToUserFormat(f.ClName), f.MethName)
		status := exceptions.ThrowEx(excNames.InternalException, errMsg, f)
		if status != exceptions.Caught {
			return nil // applies only if in test
		}
	}

	if globals.TraceVerbose {
		var traceInfo string
		value := f.OpStack[f.TOS]
		switch value.(type) {
		case *object.Object:
			obj := value.(*object.Object)
			trace.TraceObject(f, "PEEK", obj)
		default:
			traceInfo = fmt.Sprintf("                                                  "+
				"PEEK          TOS:%3d %T %v", f.TOS, value, value)
			trace.Trace(traceInfo)
		}
		// Trace the stack
		trace.LogTraceStack(f)
	}
	return f.OpStack[f.TOS]
}

// push onto the operand stack
func push(f *frames.Frame, x interface{}) {
	if f.TOS == len(f.OpStack)-1 {
		errMsg := fmt.Sprintf("in %s.%s, exceeded op stack size of %d",
			util.ConvertInternalClassNameToUserFormat(f.ClName), f.MethName, len(f.OpStack))
		status := exceptions.ThrowEx(excNames.StackOverflowError, errMsg, f)
		if status != exceptions.Caught {
			return // applies only if in test
		}
	}

	// we show trace info of the TOS *before* we change its value--
	// all traces show TOS before the instruction is executed.
	if globals.TraceVerbose {
		var traceInfo string

		if f.TOS == -1 {
			traceInfo = fmt.Sprintf("%77s", "PUSH          TOS:  -")
			trace.Trace(traceInfo)
		} else {
			if x == nil {
				traceInfo = fmt.Sprintf("%74s", "PUSH          TOS:") +
					fmt.Sprintf("%3d <nil>", f.TOS)
				trace.Trace(traceInfo)
			} else {
				if x == object.Null {
					traceInfo = fmt.Sprintf("%74s", "PUSH          TOS:") +
						fmt.Sprintf("%3d null", f.TOS)
					trace.Trace(traceInfo)
				} else {
					switch x.(type) {
					case *object.Object:
						obj := x.(*object.Object)
						trace.TraceObject(f, "PUSH", obj)
					case *[]uint8:
						strPtr := x.(*[]byte)
						str := string(*strPtr)
						traceInfo = fmt.Sprintf("%74s", "PUSH          TOS:") +
							fmt.Sprintf("%3d *[]byte: %-10s", f.TOS, str)
						trace.Trace(traceInfo)
					case []uint8:
						bytes := x.([]byte)
						str := string(bytes)
						traceInfo = fmt.Sprintf("%74s", "PUSH          TOS:") +
							fmt.Sprintf("%3d []byte: %-10s", f.TOS, str)
						trace.Trace(traceInfo)
					case []types.JavaByte:
						bytes := x.([]types.JavaByte)
						str := object.GoStringFromJavaByteArray(bytes)
						traceInfo = fmt.Sprintf("%74s", "PUSH          TOS:") +
							fmt.Sprintf("%3d []javaByte: %-10s", f.TOS, str)
						trace.Trace(traceInfo)
					default:
						traceInfo = fmt.Sprintf("%56s", " ") +
							fmt.Sprintf("PUSH          TOS:%3d %T %v", f.TOS, x, x)
						trace.Trace(traceInfo)
					}
				}
			}
		}
	}

	// the actual push
	f.TOS += 1
	f.OpStack[f.TOS] = x
	if globals.TraceVerbose {
		trace.LogTraceStack(f)
	} // trace the resultant stack
}

// determines whether classA is a subset of classB, using the stringpool indices that point to the class names
func isClassAaSublclassOfB(classA uint32, classB uint32) bool {
	if classA == classB {
		return true
	}

	superclasses := getSuperclasses(classA)
	if len(superclasses) > 0 {
		for _, superclass := range superclasses {
			if superclass == classB {
				return true
			}
		}
	}
	return false
}

// accepts a stringpool index to the classname and returns an array of names of superclasses.
// These names are returned in the form of stringPool indexes, that is, uint32 values.
func getSuperclasses(classNameIndex uint32) []uint32 {
	retval := []uint32{}
	if classNameIndex == types.InvalidStringIndex {
		return retval
	}

	if classNameIndex == types.ObjectPoolStringIndex { // if the object is java/lang/Object, it has no superclasses
		return retval
	}

	thisClassName := stringPool.GetStringPointer(classNameIndex)
	thisClass := classloader.MethAreaFetch(*thisClassName)
	thisClassSuper := thisClass.Data.SuperclassIndex

	retval = append(retval, thisClassSuper)

	if thisClassSuper == types.ObjectPoolStringIndex { // is the immediate superclass java/lang/Object? most cases = yes
		return retval
	}

	for {
		idx := thisClassSuper
		thisClassName = stringPool.GetStringPointer(idx)
		thisClass = classloader.MethAreaFetch(*thisClassName)
		if thisClass == nil {
			_ = classloader.LoadClassFromNameOnly(*thisClassName)
			thisClass = classloader.MethAreaFetch(*thisClassName)
		}

		thisClassSuper = thisClass.Data.SuperclassIndex
		retval = append(retval, thisClassSuper)

		if thisClassSuper == types.ObjectPoolStringIndex { // is the superclass java/lang/Object? If so, this is the
			break // loop's exit condition as all objects have java/lang/Object at the top of their superclass hierarchy
		} else {
			idx = thisClassSuper
		}
	}
	return retval
}

func checkcastNonArrayObject(obj *object.Object, className string) bool {
	// the object being checked is a class
	// glob := globals.GetGlobalRef()
	classPtr := classloader.MethAreaFetch(className)
	if classPtr == nil { // class wasn't loaded, so load it now
		if classloader.LoadClassFromNameOnly(className) != nil {
			// glob.ErrorGoStack = string(debug.Stack())
			// return errors.New("CHECKCAST: Could not load class: "
			// + className)
			return false
		}
		classPtr = classloader.MethAreaFetch(className)
	}

	// if classPtr does not point to the entry for the same class, then examine superclasses
	if classPtr == classloader.MethAreaFetch(*(stringPool.GetStringPointer(obj.KlassName))) {
		return true
	} else if isClassAaSublclassOfB(obj.KlassName, stringPool.GetStringIndex(&className)) {
		return true
	}
	return false
}

// do the checkcast logic for an array. The rules are:
// S = obj
// T = className
//
// If S is the type of the object referred to by objectref, and T is the resolved class, array, or
// interface type, then checkcast determines whether objectref can be cast to type T as follows:
//
// If S is an array type SC[], that is, an array of components of type SC, then:
// * If T is a class type, then T must be Object.
// * If T is an interface type, then T must be one of the interfaces implemented by arrays (JLS §4.10.3).
// * If T is an array type TC[], that is, an array of components of type TC,
// then one of the following must be true:
// >          TC and SC are the same primitive type.
// >          TC and SC are reference types, and type SC can be cast to TC by
// >             recursive application of these rules.
func checkcastArray(obj *object.Object, className string) bool {
	if obj.KlassName == types.InvalidStringIndex {
		errMsg := "CHECKCAST: expected to verify class or interface, but got none"
		status := exceptions.ThrowExNil(excNames.InvalidTypeException, errMsg)
		if status != exceptions.Caught {
			return false // applies only if in test
		}
	}

	sptr := stringPool.GetStringPointer(obj.KlassName)
	// if they're both the same type of arrays, we're good
	if *sptr == className || strings.HasPrefix(className, *sptr) {
		return true
	}

	// If S (obj) is an array type SC[], that is, an array of components of type SC,
	// then: If T (className) is a class type, then T must be Object.
	if !strings.HasPrefix(className, types.Array) {
		return className == "java/lang/Object"
	}

	// If S (obj) is an array type SC[], that is, an array of components of type SC,
	// if T is an array type TC[], that is, an array of components of type TC,
	// then one of the following must be true:
	// >          TC and SC are the same primitive type.
	objArrayType := object.GetArrayType(*sptr)
	classArrayType := object.GetArrayType(className)
	if !strings.HasPrefix(objArrayType, "L") && // if both array types are primitives
		!strings.HasPrefix(classArrayType, "L") {
		if objArrayType == classArrayType { // are they the same primitive?
			return true
		}
	}

	// we now know both types are arrays of references, so we test to see whether
	// the reference object is castable, using this guideline:
	// If TC and SC are reference types, and type SC can be cast to TC by
	//    recursive application of these rules.
	rawObjArrayType, _ := strings.CutPrefix(objArrayType, "L")
	rawObjArrayType = strings.TrimSuffix(rawObjArrayType, ";")
	rawClassArrayType, _ := strings.CutPrefix(classArrayType, "L")
	rawClassArrayType = strings.TrimSuffix(rawClassArrayType, ";")
	if rawObjArrayType == classArrayType || rawClassArrayType == "java/lang/Object" {
		return true
	} else {
		return isClassAaSublclassOfB(
			stringPool.GetStringIndex(&rawObjArrayType),
			stringPool.GetStringIndex(&rawClassArrayType))
	}
}

func checkcastInterface(obj *object.Object, className string) bool {
	return false // TODO: fill this in
}

// the function that finds the interface method to execute (and returns it).
func locateInterfaceMeth(
	class *classloader.Klass, // the objRef class
	f *frames.Frame,
	objRefClassName string,
	interfaceName string,
	interfaceMethodName string,
	interfaceMethodType string) (classloader.MTentry, error) {

	glob := globals.GetGlobalRef()

	// Find the interface method. Section 5.4.3.4 of the JVM spec lists the order in which
	// the steps are taken, where C is the interface:
	//
	// 1) If C is not an interface, interface method resolution throws an IncompatibleClassChangeError.
	//
	// 2) Otherwise, if C declares a method with the name and descriptor specified by the
	// interface method reference, method lookup succeeds.
	//
	// 3) Otherwise, if the class Object declares a method with the name and descriptor specified by the
	// interface method reference, which has its ACC_PUBLIC flag set and does not have its ACC_STATIC flag set,
	// method lookup succeeds.
	//
	// 4) Otherwise, if the maximally-specific superinterface methods (§5.4.3.3) of C for the name and descriptor
	// specified by the method reference include exactly one method that does not have its ACC_ABSTRACT flag set,
	// then this method is chosen and method lookup succeeds.
	//
	// 5) Otherwise, if any superinterface of C declares a method with the name and descriptor specified by the
	// method reference that has neither its ACC_PRIVATE flag nor its ACC_STATIC flag set, one of these is
	// arbitrarily chosen and method lookup succeeds.
	//
	// 6) Otherwise, method lookup fails.
	//
	// For more info: https://docs.oracle.com/javase/specs/jvms/se17/html/jvms-5.html#jvms-5.4.3.4

	clData := *class.Data
	if len(clData.Interfaces) == 0 { // TODO: Determine whether this is correct behavior. See Jacotest results.
		errMsg := fmt.Sprintf("INVOKEINTERFACE: class %s does not implement interface %s",
			objRefClassName, interfaceName)
		status := exceptions.ThrowEx(excNames.IncompatibleClassChangeError, errMsg, f)
		if status != exceptions.Caught {
			return classloader.MTentry{}, errors.New(errMsg) // applies only if in test
		}
	}

	var foundIntfaceName = ""
	var mtEntry classloader.MTentry
	var meth *classloader.Method
	var ok bool
	for i := 0; i < len(clData.Interfaces); i++ {
		index := uint32(clData.Interfaces[i])
		foundIntfaceName = *stringPool.GetStringPointer(index)
		if foundIntfaceName == interfaceName {
			// at this point we know that clData's class implements the required interface.
			// Now, check whether clData contains the desired method.
			meth, ok = clData.MethodTable[interfaceMethodName+interfaceMethodType]
			if ok {
				mtEntry, _ = classloader.FetchMethodAndCP(
					clData.Name, interfaceMethodName, interfaceMethodType)
				goto verifyInterfaceMethod
			}

			if err := classloader.LoadClassFromNameOnly(interfaceName); err != nil {
				// in this case, LoadClassFromNameOnly() will have already thrown the exception
				if globals.JacobinHome() == "test" {
					return classloader.MTentry{}, err // applies only if in test
				}
			}
			mtEntry, _ = classloader.FetchMethodAndCP(
				interfaceName, interfaceMethodName, interfaceMethodType)
			if mtEntry.Meth == nil {
				glob.ErrorGoStack = string(debug.Stack())
				errMsg := fmt.Sprintf("INVOKEINTERFACE: Interface method not found: %s.%s%s"+
					interfaceName, interfaceMethodName, interfaceMethodType)
				status := exceptions.ThrowEx(excNames.ClassNotLoadedException, errMsg, f)
				if status != exceptions.Caught {
					return classloader.MTentry{}, errors.New(errMsg) // applies only if in test
				}
			}
			goto verifyInterfaceMethod // method found, move on to execution
		} else { // CURR: check for superclasses, after checking Object
			foundIntfaceName = ""
		}
	}

	if foundIntfaceName == "" { // no interface was found, check java.lang.Object()
		errMsg := fmt.Sprintf("INVOKEINTERFACE: class %s does not implement interface %s",
			objRefClassName, interfaceName)
		status := exceptions.ThrowEx(excNames.IncompatibleClassChangeError, errMsg, f)
		if status != exceptions.Caught {
			return classloader.MTentry{}, errors.New(errMsg) // applies only if in test
		}
	}

verifyInterfaceMethod:
	if mtEntry.MType == 'J' && meth.AccessFlags&0x0100 > 0 { // if a J method calls native code, JVM spec throws exception
		glob.ErrorGoStack = string(debug.Stack())
		errMsg := "INVOKEINTERFACE: Native method requested: " +
			clData.Name + "." + interfaceMethodName + interfaceMethodType
		status := exceptions.ThrowEx(excNames.UnsupportedOperationException, errMsg, f)
		if status != exceptions.Caught {
			return classloader.MTentry{}, errors.New(errMsg) // applies only if in test
		}
	}

	return mtEntry, nil
}
