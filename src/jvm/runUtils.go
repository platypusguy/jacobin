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
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/exceptions"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/opcodes"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"jacobin/src/types"
	"jacobin/src/util"
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

// converts a golang interface{} value to int8. Used for BASTORE, among others.
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

// converts a golang interface{} value on the op stack into a uint64
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

// converts a golang interface{} value into int64
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
	// if globals.TraceVerbose {
	var value interface{}

	if f.TOS == -1 {
		errMsg := fmt.Sprintf("stack underflow in pop() in %s.%s%s",
			util.ConvertInternalClassNameToUserFormat(f.ClName), f.MethName, f.MethType)
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
					TraceObject(f, "POP", obj)
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
		LogTraceStack(f)
	} // trace the resultant stack

	// Return value to caller.
	return value
	// } else {
	// 	value := f.OpStack[f.TOS]
	// 	f.TOS -= 1
	// 	return value
	// }

}

// returns the value at the top of the stack without popping it off.
func peek(f *frames.Frame) interface{} {
	if f.TOS == -1 {
		errMsg := fmt.Sprintf("stack underflow in peek() in %s.%s%s",
			util.ConvertInternalClassNameToUserFormat(f.ClName), f.MethName, f.MethType)
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
			TraceObject(f, "PEEK", obj)
		default:
			traceInfo = fmt.Sprintf("                                                  "+
				"PEEK          TOS:%3d %T %v", f.TOS, value, value)
			trace.Trace(traceInfo)
		}
		// Trace the stack
		LogTraceStack(f)
	}
	return f.OpStack[f.TOS]
}

// push onto the operand stack
func push(f *frames.Frame, x interface{}) {
	// if globals.TraceVerbose {
	if f.TOS == len(f.OpStack)-1 {
		errMsg := fmt.Sprintf("in %s.%s%s, exceeded op stack size of %d",
			util.ConvertInternalClassNameToUserFormat(f.ClName), f.MethName, f.MethType, len(f.OpStack))
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
						TraceObject(f, "PUSH", obj)
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
		LogTraceStack(f)
	} // trace the resultant stack
	// } else { // no tracing and no checking of stack size -- fastest (and the default)
	// 	f.TOS += 1
	// 	f.OpStack[f.TOS] = x
	// }
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

	if classNameIndex == types.StringPoolObjectIndex { // if the object is java/lang/Object, it has no superclasses
		return retval
	}

	thisClassName := stringPool.GetStringPointer(classNameIndex)
	thisClass := classloader.MethAreaFetch(*thisClassName)
	if thisClass == nil {
		return retval
	}
	thisClassSuper := thisClass.Data.SuperclassIndex

	retval = append(retval, thisClassSuper)

	if thisClassSuper == types.StringPoolObjectIndex { // is the immediate superclass java/lang/Object? most cases = yes
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

		if thisClassSuper == types.StringPoolObjectIndex { // is the superclass java/lang/Object? If so, this is the
			break // loop's exit condition as all objects have java/lang/Object at the top of their superclass hierarchy
		} else {
			idx = thisClassSuper
		}
	}
	return retval
}

func checkcastNonArrayObject(srcObj *object.Object, targetClassName string) bool {
	// the object being checked is a class
	// glob := globals.GetGlobalRef()
	classPtr := classloader.MethAreaFetch(targetClassName)
	if classPtr == nil { // class wasn't loaded, so load it now
		if classloader.LoadClassFromNameOnly(targetClassName) != nil {
			// glob.ErrorGoStack = string(debug.Stack())
			// return errors.New("CHECKCAST: Could not load class: "
			// + className)
			return false
		}
		classPtr = classloader.MethAreaFetch(targetClassName)
	}

	// if classPtr does not point to the entry for the same class, then examine superclasses
	if classPtr == classloader.MethAreaFetch(*(stringPool.GetStringPointer(srcObj.KlassName))) {
		return true
	} else if isClassAaSublclassOfB(srcObj.KlassName, stringPool.GetStringIndex(&targetClassName)) {
		return true
	} else {
		// Casting from a java/lang/Object containing a Java byte array to a java/lang/String?
		if isObjectBytesToString(srcObj, targetClassName) {
			return true
		}
	}
	// None of the above
	return false
}

// Allow a cast from src:java/lang/Object to dest:java/lang/String if src.value is a Java byte array.
func isObjectBytesToString(srcObj *object.Object, destClassName string) bool {
	// Destination object is of type java/lang/String?
	if destClassName != types.StringClassName {
		return false
	}
	// Source object is of type java/lang/Object?
	srcClassName := *stringPool.GetStringPointer(srcObj.KlassName)
	if srcClassName != types.ObjectClassName {
		return false
	}
	// Source object value field contains an array of Java bytes?
	srcField, ok := srcObj.FieldTable["value"]
	if !ok {
		return false
	}
	switch srcField.Fvalue.(type) {
	case []types.JavaByte:
		return true
	}
	// None of the above.
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
func checkcastArray(obj *object.Object, targetClassName string) bool {
	if obj.KlassName == types.InvalidStringIndex {
		errMsg := "CHECKCAST: expected to verify class or interface, but got none"
		status := exceptions.ThrowExNil(excNames.InvalidTypeException, errMsg)
		if status != exceptions.Caught {
			return false // applies only if in test
		}
	}

	sptr := stringPool.GetStringPointer(obj.KlassName)
	// if they're both the same type of arrays, we're good
	if *sptr == targetClassName {
		return true
	}

	// If S (obj) is an array type SC[], that is, an array of components of type SC,
	// then: If T (className) is a class type, then T must be Object.
	if !strings.HasPrefix(targetClassName, types.Array) {
		return targetClassName == "java/lang/Object"
	}

	// If S (obj) is an array type SC[], that is, an array of components of type SC,
	// if T is an array type TC[], that is, an array of components of type TC,
	// then one of the following must be true:
	// >          TC and SC are the same primitive type.
	objArrayType := object.GetArrayType(*sptr)
	classArrayType := object.GetArrayType(targetClassName)
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

	// String output test.

	// Subclass test
	if rawObjArrayType == rawClassArrayType {
		return true
	}
	if rawClassArrayType == types.ObjectClassName {
		return true
	}
	if isClassAaSublclassOfB(
		stringPool.GetStringIndex(&rawObjArrayType),
		stringPool.GetStringIndex(&rawClassArrayType)) {
		return true
	}

	// None of the above
	return false
}

func checkcastInterface(obj *object.Object, targetClassName string) bool {
	return true // TODO: fill this in     2026-02-05 This was returning false
}

// the function that finds the interface method to execute (and returns it).
// This is a two-part process: first, we verify the signature of the method,
// then we locate the concrete implementation.
//
// Note: this function is similar in many aspects to searchForDefaultInterfaceFunction()
// in interpreter.go. The two functions might eventually be integrated into one.
func locateInterfaceMeth(
	class *classloader.Klass, // the objRef class
	f *frames.Frame,
	objRefClassName string,
	interfaceName string,
	interfaceMethodName string,
	interfaceMethodType string) (classloader.MTentry, error) {

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

	// == Phase 1: Verify the signature of the method (might find an abastract method, that's OK)
	// step 1: check whether the interface is truly an interface
	interfaceKlass := classloader.MethAreaFetch(interfaceName)
	if interfaceKlass == nil {
		if globals.TraceVerbose {
			trace.Trace(fmt.Sprintf("[INVOKEINTERFACE] Interface %s not in method area, loading...", interfaceName))
		}
		err := classloader.LoadClassFromNameOnly(interfaceName)
		if err != nil {
			errMsg := fmt.Sprintf("INVOKEINTERFACE: Failed to load interface %s: %v", interfaceName, err)
			status := exceptions.ThrowEx(excNames.NoClassDefFoundError, errMsg, f)
			if status != exceptions.Caught {
				return classloader.MTentry{}, errors.New(errMsg)
			}
			return classloader.MTentry{}, errors.New(errMsg) // for tests
		}
		interfaceKlass = classloader.MethAreaFetch(interfaceName)
		if interfaceKlass == nil {
			errMsg := fmt.Sprintf("INVOKEINTERFACE: Interface %s not found in method area after loading", interfaceName)
			status := exceptions.ThrowEx(excNames.NoClassDefFoundError, errMsg, f)
			if status != exceptions.Caught {
				return classloader.MTentry{}, errors.New(errMsg)
			}
			return classloader.MTentry{}, errors.New(errMsg) // for tests
		}
	}

	// Now interfaceKlass is guaranteed to be non-nil
	if !interfaceKlass.Data.Access.ClassIsInterface {
		errMsg := fmt.Sprintf("INVOKEINTERFACE: %s is not an interface", interfaceName)
		status := exceptions.ThrowEx(excNames.IncompatibleClassChangeError, errMsg, f)
		if status != exceptions.Caught {
			return classloader.MTentry{}, errors.New(errMsg) // applies only if in test
		}
		return classloader.MTentry{}, errors.New(errMsg) // for tests
	}

	if globals.TraceVerbose {
		trace.Trace(fmt.Sprintf("[INVOKEINTERFACE] Step 1: Verified %s is an interface", interfaceName))
	}

	signatureFound := false

	// step 2: Check if interface C directly declares the method
	// Per spec §5.4.3.4, this succeeds even if the method is abstract.
	// Abstract methods will be caught during invocation.
	var mtEntry classloader.MTentry
	var err error

	mtEntry, err = classloader.FetchMethodAndCP(
		interfaceName, interfaceMethodName, interfaceMethodType)
	if err == nil && mtEntry.Meth != nil {
		signatureFound = true
	}

	// step 3: Check if java/lang/Object declares the method
	// this is already done in FetchMethodAndCP as the methods are loaded into every class,
	// so we don't need to do it again

	// step 4: finc the maximally-specific method, which means: suppose we have:
	// interface A { void m(); }
	// interface B extends A { default void m() { } }
	// interface C extends A { default void m() { } }
	// interface D extends B, C { } // Which m() to use?
	// A "maximally-specific" method is one that:
	// * is NOT abstract (has implementation)
	// * is NOT overridden by any other candidate method (is "most specific")
	if !signatureFound {
		_, found := findMaximallySpecificSuperinterfaceMethods(
			interfaceName, interfaceMethodName, interfaceMethodType, f)
		if found {
			signatureFound = true // continue to part 2
		}
	}

	// clData := *class.Data
	// if len(clData.Interfaces) == 0 { // TODO: Determine whether this is correct behavior. See Jacotest results.
	// 	errMsg := fmt.Sprintf("INVOKEINTERFACE: class %s does not implement interface %s",
	// 		objRefClassName, interfaceName)
	// 	status := exceptions.ThrowEx(excNames.IncompatibleClassChangeError, errMsg, f)
	// 	if status != exceptions.Caught {
	// 		return classloader.MTentry{}, errors.New(errMsg) // applies only if in test
	// 	}
	// }

	// STEP 5: Any superinterface (with filtering)
	if !signatureFound {
		superInterfaces := getSuperInterfaces([]string{interfaceName})
		for _, siface := range superInterfaces {
			mtEntry, err = classloader.FetchMethodAndCP(siface, interfaceMethodName, interfaceMethodType)
			if err == nil && mtEntry.Meth != nil {
				m := mtEntry.Meth.(classloader.JmEntry) // filter out abstract and private methods
				if m.AccessFlags&classloader.ACC_ABSTRACT > 0 ||
					m.AccessFlags&classloader.ACC_PRIVATE > 0 {
					continue
				}
				signatureFound = true
				break
			}
		}
	}

	if !signatureFound {
		errMsg := "INVOKEINTERFACE: Interface method not found: " +
			interfaceName + "." + interfaceMethodName + interfaceMethodType
		status := exceptions.ThrowEx(excNames.NoSuchMethodError, errMsg, f)
		if status == exceptions.Caught {
			return classloader.MTentry{}, errors.New(errMsg) // applies only if in test
		}
	}

	// === Phase 2: Find the concrete implementation of the method
	// check whether the class or its superclasses directly implement the method
	mtEntry, _ = classloader.FetchMethodAndCP(
		objRefClassName, interfaceMethodName, interfaceMethodType)
	if err == nil && mtEntry.Meth != nil {
		// found concrete implementation in the class
		return mtEntry, nil
	}

	// check all the interfaces this class implements, going from left to right
	// in the interface declarations.
	interfaces := getClassInterfaces(class)
	for _, iface := range interfaces {
		mtEntry, err = classloader.FetchMethodAndCP(
			iface, interfaceMethodName, interfaceMethodType)
		if err == nil && mtEntry.Meth != nil {
			// Check if it's NOT abstract (i.e., a default method)
			if mtEntry.MType == 'J' {
				jmEntry := mtEntry.Meth.(classloader.JmEntry)
				if (jmEntry.AccessFlags & classloader.ACC_ABSTRACT) == 0 {
					return mtEntry, nil // Found default method
				}
			}
		}
	}

	// if we got here, the method was not found in the class or its interfaces,
	// so we look at the superinterfaces
	superInterfaces := getSuperInterfaces(interfaces)
	for _, siface := range superInterfaces {
		mtEntry, err = classloader.FetchMethodAndCP(
			siface, interfaceMethodName, interfaceMethodType)
		if err == nil && mtEntry.Meth != nil {
			if mtEntry.MType == 'J' {
				jmEntry := mtEntry.Meth.(classloader.JmEntry)
				if (jmEntry.AccessFlags & classloader.ACC_ABSTRACT) == 0 {
					return mtEntry, nil
				}
			}
		}
	}

	// if we got here, the method was not found in the interface, java/lang/Object, or in superinterfaces
	globals.GetGlobalRef().ErrorGoStack = string(debug.Stack())
	errMsg := fmt.Sprintf("INVOKEINTERFACE: Interface method not found: %s.%s%s",
		interfaceName, interfaceMethodName, interfaceMethodType)
	status := exceptions.ThrowEx(excNames.ClassNotLoadedException, errMsg, f)
	if status != exceptions.Caught {
		return classloader.MTentry{}, errors.New(errMsg) // applies only if in test
	}

	// unreachable due to exception thrown immediately above
	// but golang reports an error if we don't return here
	return classloader.MTentry{}, errors.New(errMsg)
}

// goes through the set of interfaces that a class implements and returns them as an array of strings
func getClassInterfaces(class *classloader.Klass) []string {
	interfaces := []string{}
	clData := *class.Data
	for i := 0; i < len(clData.Interfaces); i++ {
		index := uint32(clData.Interfaces[i])
		interfaces = append(interfaces, *stringPool.GetStringPointer(index))
	}
	return interfaces
}

// accepts an array of interfaces names and returns an array of superinterfaces
func getSuperInterfaces(interfaces []string) []string {
	superinterfaces := []string{}

	// this gets one level of superinterfaces to all the interfaces in retval
	for i := 0; i < len(interfaces); i++ {
		interfaceName := interfaces[i]
		interfaceClass := classloader.MethAreaFetch(interfaceName)
		if interfaceClass == nil {
			if err := classloader.LoadClassFromNameOnly(interfaceName); err != nil {
				// in this case, LoadClassFromNameOnly() will have already thrown the exception
				if globals.JacobinHome() == "test" {
					return superinterfaces // applies only if in test. At this point, superinterfaces is empty
				}
			}
		}

		for j := 0; j < len(interfaceClass.Data.Interfaces); j++ {
			index := uint32(interfaceClass.Data.Interfaces[j])
			superinterfaces = append(superinterfaces, *stringPool.GetStringPointer(index))
		}

		// get any superinterfaces of the superinterfaces
		superSuperInterfaces := []string{}
		for k := 0; k < len(superinterfaces); k++ {
			superInterfaceName := superinterfaces[k]
			if err := classloader.LoadClassFromNameOnly(superInterfaceName); err != nil {
				// in this case, LoadClassFromNameOnly() will have already thrown the exception
				if globals.JacobinHome() == "test" {
					return []string{} // applies only if in test
				}
			}

			superInterfaceClass := classloader.MethAreaFetch(superInterfaceName)
			for m := 0; m < len(superInterfaceClass.Data.Interfaces); m++ {
				index := uint32(superInterfaceClass.Data.Interfaces[m])
				superSuperInterfaces = append(superSuperInterfaces, *stringPool.GetStringPointer(index))
			}
		}
		superinterfaces = append(superinterfaces, superSuperInterfaces...)
	}
	return superinterfaces
}

// findMaximallySpecificSuperinterfaceMethods implements JVM spec §5.4.3.3
// Returns a non-abstract method if exactly one maximally-specific method exists
// in the superinterface hierarchy. Returns (mtEntry, true) if found, (empty, false) otherwise.
//
// Per spec: "if the maximally-specific superinterface methods of C for the name
// and descriptor include exactly one method that does not have its ACC_ABSTRACT
// flag set, then this method is chosen"
func findMaximallySpecificSuperinterfaceMethods(
	interfaceName string,
	methodName string,
	methodType string,
	f *frames.Frame) (classloader.MTentry, bool) {

	// Get the interface C
	interfaceKlass := classloader.MethAreaFetch(interfaceName)
	if interfaceKlass == nil {
		_ = classloader.LoadClassFromNameOnly(interfaceName)
		interfaceKlass = classloader.MethAreaFetch(interfaceName)
	}
	if interfaceKlass == nil {
		return classloader.MTentry{}, false
	}

	// Get all direct superinterfaces of C (not transitive)
	directSuperInterfaces := []string{}
	for _, idx := range interfaceKlass.Data.Interfaces {
		superIfaceName := *stringPool.GetStringPointer(uint32(idx))
		directSuperInterfaces = append(directSuperInterfaces, superIfaceName)
	}

	if len(directSuperInterfaces) == 0 {
		return classloader.MTentry{}, false
	}

	// Find all methods with matching signature in superinterfaces
	// that are NOT abstract (i.e., default methods or static methods)
	type MethodCandidate struct {
		entry      classloader.MTentry
		iface      string
		isAbstract bool
	}

	candidates := []MethodCandidate{}

	// Check each superinterface and its superinterfaces recursively
	checkedInterfaces := make(map[string]bool)

	var checkInterface func(ifaceName string)
	checkInterface = func(ifaceName string) {
		if checkedInterfaces[ifaceName] {
			return
		}
		checkedInterfaces[ifaceName] = true

		// Load interface if needed
		iface := classloader.MethAreaFetch(ifaceName)
		if iface == nil {
			_ = classloader.LoadClassFromNameOnly(ifaceName)
			iface = classloader.MethAreaFetch(ifaceName)
		}
		if iface == nil {
			return
		}

		// Check if this interface has the method
		mtEntry, err := classloader.FetchMethodAndCP(ifaceName, methodName, methodType)
		if err == nil && mtEntry.Meth != nil {
			candidate := MethodCandidate{
				entry: mtEntry,
				iface: ifaceName,
			}

			// Check if abstract
			if mtEntry.MType == 'J' {
				jmEntry := mtEntry.Meth.(classloader.JmEntry)
				candidate.isAbstract = (jmEntry.AccessFlags & 0x0400) != 0 // ACC_ABSTRACT
			} else {
				candidate.isAbstract = false // G-functions are concrete
			}

			candidates = append(candidates, candidate)
		}

		// Recursively check superinterfaces
		for _, idx := range iface.Data.Interfaces {
			superIfaceName := *stringPool.GetStringPointer(uint32(idx))
			checkInterface(superIfaceName)
		}
	}

	// Check all superinterfaces of C
	for _, si := range directSuperInterfaces {
		checkInterface(si)
	}

	if len(candidates) == 0 {
		return classloader.MTentry{}, false
	}

	// Filter to only non-abstract methods
	nonAbstractCandidates := []MethodCandidate{}
	for _, candidate := range candidates {
		if !candidate.isAbstract {
			nonAbstractCandidates = append(nonAbstractCandidates, candidate)
		}
	}

	// Spec requires EXACTLY ONE non-abstract method
	if len(nonAbstractCandidates) == 1 {
		return nonAbstractCandidates[0].entry, true
	}

	// Multiple non-abstract methods found (ambiguous)
	// or only abstract methods found
	// In either case, this step fails
	if len(nonAbstractCandidates) > 1 {
		// This would be an IncompatibleClassChangeError at runtime
		// but for now we just return false to try the next step
		if globals.TraceVerbose {
			errMsg := fmt.Sprintf("[Step 4] Multiple non-abstract methods found for %s.%s%s",
				interfaceName, methodName, methodType)
			trace.Trace(errMsg)
		}
	}

	return classloader.MTentry{}, false
}

// the generation and formatting of trace data for each executed bytecode.
// Returns the formatted data for output to logging, console, or other uses.
func EmitTraceData(f *frames.Frame) string {
	var tos = " -"
	var stackTop = ""
	if f.TOS != -1 {
		tos = fmt.Sprintf("%2d", f.TOS)
		switch f.OpStack[f.TOS].(type) {
		// if the value at TOS is a string, say so and print the first 10 chars of the string
		case *object.Object:
			if object.IsNull(f.OpStack[f.TOS].(*object.Object)) {
				stackTop = fmt.Sprintf("<null>")
			} else {
				objPtr := f.OpStack[f.TOS].(*object.Object)
				if objPtr.KlassName == types.StringPoolStringIndex {
					str := object.GoStringFromStringObject(objPtr)
					stackTop = fmt.Sprintf("String: %-10s", str)
				} else {
					stackTop = objPtr.FormatField("")
				}
			}
		case *[]uint8:
			value := f.OpStack[f.TOS]
			strPtr := value.(*[]byte)
			str := string(*strPtr)
			stackTop = fmt.Sprintf("*[]byte: %-10s", str)
		case []uint8:
			value := f.OpStack[f.TOS]
			bytes := value.([]byte)
			str := string(bytes)
			stackTop = fmt.Sprintf("[]byte: %-10s", str)
		case []types.JavaByte:
			value := f.OpStack[f.TOS]
			bytes := value.([]types.JavaByte)
			str := object.GoStringFromJavaByteArray(bytes)
			stackTop = fmt.Sprintf("[]JavaByte: %-10s", str)
		default:
			stackTop = fmt.Sprintf("%T %v ", f.OpStack[f.TOS], f.OpStack[f.TOS])
		}
	}

	traceInfo := fmt.Sprintf("th: %d, class: %-22s meth:%-10s PC: % 3d, %-13s TOS: %s %s ",
		f.Thread, f.ClName, f.MethName, f.PC, opcodes.BytecodeNames[int(f.Meth[f.PC])], tos, stackTop)
	return traceInfo
}

// Generate a trace of a field ID (static or non-static).
func EmitTraceFieldID(opcode, fld string) {
	traceInfo := fmt.Sprintf("%65s fieldName: %s", opcode, fld)
	trace.Trace(traceInfo)
}

// Log the existing stack
// Could be called for tracing -or- supply info for an error section
func LogTraceStack(f *frames.Frame) {
	var traceInfo, output string
	if f.TOS == -1 {
		traceInfo = fmt.Sprintf("%55s %s.%s stack <empty>", "", f.ClName, f.MethName)
		trace.Trace(traceInfo)
		return
	}
	for ii := 0; ii <= f.TOS; ii++ {
		switch f.OpStack[ii].(type) {
		case *object.Object:
			if object.IsNull(f.OpStack[ii].(*object.Object)) {
				output = fmt.Sprintf("<null>")
			} else {
				objPtr := f.OpStack[ii].(*object.Object)
				output = objPtr.FormatField("")
			}
		case *[]uint8:
			value := f.OpStack[ii]
			strPtr := value.(*[]byte)
			str := string(*strPtr)
			output = fmt.Sprintf("*[]byte: %-10s", str)
		case []uint8:
			value := f.OpStack[ii]
			bytes := value.([]byte)
			str := string(bytes)
			output = fmt.Sprintf("[]byte: %-10s", str)
		case []types.JavaByte:
			value := f.OpStack[ii]
			bytes := value.([]types.JavaByte)
			str := object.GoStringFromJavaByteArray(bytes)
			output = fmt.Sprintf("[]javaByte: %-10s", str)
		default:
			output = fmt.Sprintf("%T %v ", f.OpStack[ii], f.OpStack[ii])
		}
		if f.TOS == ii {
			traceInfo = fmt.Sprintf("%55s %s.%s TOS   [%d] %s", "", f.ClName, f.MethName, ii, output)
		} else {
			traceInfo = fmt.Sprintf("%55s %s.%s stack [%d] %s", "", f.ClName, f.MethName, ii, output)
		}
		trace.Trace(traceInfo)
	}
}

// TraceObject : Used by push, pop, and peek in tracing an object.
func TraceObject(f *frames.Frame, opStr string, obj *object.Object) {
	var traceInfo string
	prefix := fmt.Sprintf(" %4s          TOS:", opStr)

	// Nil pointer to object?
	if obj == nil {
		traceInfo = fmt.Sprintf("%74s%3d null", prefix, f.TOS)
		trace.Trace(traceInfo)
		return
	}

	// The object pointer is not nil.
	klass := object.GoStringFromStringPoolIndex(obj.KlassName)
	traceInfo = fmt.Sprintf("%74s%3d, class: %s", prefix, f.TOS, klass)
	trace.Trace(traceInfo)

	// Trace field table.
	prefix = " "
	if len(obj.FieldTable) > 0 {
		for fieldName := range obj.FieldTable {
			fld := obj.FieldTable[fieldName]
			if klass == types.StringClassName && fieldName == "value" {
				var str string
				switch fld.Fvalue.(type) {
				case []types.JavaByte:
					str = object.GoStringFromJavaByteArray(fld.Fvalue.([]types.JavaByte))
				default:
					str = string(fld.Fvalue.([]byte))
				}

				traceInfo = fmt.Sprintf("%74sfield: %s %s %v \"%s\"", prefix, fieldName, fld.Ftype, fld.Fvalue, str)
			} else {
				traceInfo = fmt.Sprintf("%74sfield: %s %s %v", prefix, fieldName, fld.Ftype, fld.Fvalue)
			}
			trace.Trace(traceInfo)
		}
	} else { // nil FieldTable
		traceInfo = fmt.Sprintf("%74sno fields", prefix)
		trace.Trace(traceInfo)
	}
}
