/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"fmt"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/trace"
	"strings"
)

// ResolveMethodHandle resolves a MethodHandle constant pool entry into a runtime representation.
// This corresponds to the JVM's resolution of a CONSTANT_MethodHandle_info structure.
//
// The resolution process depends on the reference_kind (1-9), which determines whether
// the handle points to a field, a method, a constructor, or an interface method.
//
// Returns:
//   - A pointer to a java.lang.invoke.MethodHandle object (as *object.Object)
//   - An error if resolution fails
func ResolveMethodHandle(cp *CPool, index int, fr *frames.Frame) (*object.Object, error) {
	if index < 1 || index >= len(cp.CpIndex) {
		return nil, fmt.Errorf("ResolveMethodHandle: invalid CP index %d", index)
	}

	// Fetch the MethodHandle entry from the constant pool
	// Note: FetchCPentry returns a CpType struct. For MethodHandle, RetType is IS_STRUCT_ADDR
	// and AddrVal points to a CPuint16s struct where entry1 is RefKind and entry2 is RefIndex.
	mhEntry := FetchCPentry(cp, index)
	if mhEntry.EntryType != MethodHandle {
		return nil, fmt.Errorf("ResolveMethodHandle: CP entry at %d is not a MethodHandle (type %d)", index, mhEntry.EntryType)
	}

	refKind := uint8(mhEntry.AddrVal.entry1) // uint8
	refIndex := mhEntry.AddrVal.entry2

	if globals.TraceClass {
		trace.Trace(fmt.Sprintf("ResolveMethodHandle: Resolving MH at index %d, kind=%d, refIndex=%d", index, refKind, refIndex))
	}

	// The resolution logic varies significantly based on the reference kind.
	// See JVM Spec 5.4.3.5. Method Type and Method Handle Resolution
	switch refKind {
	case 1: // REF_getField
		return resolveFieldHandle(cp, int(refIndex), false, false, fr, refKind)
	case 2: // REF_getStatic
		return resolveFieldHandle(cp, int(refIndex), true, false, fr, refKind)
	case 3: // REF_putField
		return resolveFieldHandle(cp, int(refIndex), false, true, fr, refKind)
	case 4: // REF_putStatic
		return resolveFieldHandle(cp, int(refIndex), true, true, fr, refKind)
	case 5: // REF_invokeVirtual
		return resolveMethodHandleEntry(cp, int(refIndex), false, false, fr, refKind)
	case 6: // REF_invokeStatic
		return resolveMethodHandleEntry(cp, int(refIndex), true, false, fr, refKind)
	case 7: // REF_invokeSpecial
		// TODO: Special handling for <init> vs other methods?
		return resolveMethodHandleEntry(cp, int(refIndex), false, true, fr, refKind)
	case 8: // REF_newInvokeSpecial
		return resolveMethodHandleEntry(cp, int(refIndex), false, true, fr, refKind) // Constructor
	case 9: // REF_invokeInterface
		return resolveMethodHandleEntry(cp, int(refIndex), false, false, fr, refKind)
	default:
		return nil, fmt.Errorf("ResolveMethodHandle: invalid reference kind %d", refKind)
	}
}

// resolveFieldHandle resolves a field access handle (kinds 1-4)
func resolveFieldHandle(cp *CPool, refIndex int, isStatic bool, isSetter bool, fr *frames.Frame, refKind uint8) (*object.Object, error) {
	// 1. Resolve the field reference in the constant pool
	// The refIndex points to a CONSTANT_Fieldref_info structure
	if refIndex < 1 || refIndex >= len(cp.CpIndex) {
		return nil, fmt.Errorf("resolveFieldHandle: invalid field ref index %d", refIndex)
	}

	// In Jacobin's CPool, FieldRefs are stored in a separate slice, indexed by the slot in CpIndex
	cpEntry := cp.CpIndex[refIndex]
	if cpEntry.Type != FieldRef {
		return nil, fmt.Errorf("resolveFieldHandle: expected FieldRef at index %d, got %d", refIndex, cpEntry.Type)
	}

	fieldRef := cp.FieldRefs[cpEntry.Slot]
	className := fieldRef.ClName
	fieldName := fieldRef.FldName
	fieldType := fieldRef.FldType

	if globals.TraceClass {
		trace.Trace(fmt.Sprintf("resolveFieldHandle: Class=%s, Field=%s, Type=%s, Static=%v, Setter=%v",
			className, fieldName, fieldType, isStatic, isSetter))
	}

	// 2. Get java.lang.Class object for the defining class.
	// getClassObj expects a descriptor, and className is an internal name (e.g. "java/lang/Object").
	defClassObj, err := getClassObj("L"+className+";", fr)
	if err != nil {
		return nil, fmt.Errorf("resolveFieldHandle: could not get Class object for %s: %w", className, err)
	}

	// 3. Get java.lang.Class object for the field's type. fieldType is already a descriptor.
	fieldTypeObj, err := getClassObj(fieldType, fr)
	if err != nil {
		return nil, fmt.Errorf("resolveFieldHandle: could not get Class object for field type %s: %w", fieldType, err)
	}

	// 4. Get java.lang.Class object for the caller class (for access checks).
	callerClassObj, err := getClassObj("L"+fr.ClName+";", fr)
	if err != nil {
		return nil, fmt.Errorf("resolveFieldHandle: could not get Class object for caller %s: %w", fr.ClName, err)
	}

	// 5. Create Java String for field name
	fieldNameObj := object.StringObjectFromGoString(fieldName)

	// 6. Invoke an internal gfunction to create the MethodHandle.
	// This is a hypothetical internal API for the VM to create method handles
	// without going through the full MethodHandles.Lookup security checks,
	// as is permitted for 'ldc' resolution. The Java-side implementation
	// of this gfunction would perform the necessary resolution and access checks.
	params := []interface{}{
		defClassObj,
		fieldNameObj,
		fieldTypeObj,
		int64(refKind),
		callerClassObj,
	}

	gfuncName := "jacobin/internal/VM.resolveFieldHandle(Ljava/lang/Class;Ljava/lang/String;Ljava/lang/Class;ILjava/lang/Class;)Ljava/lang/invoke/MethodHandle;"
	result := globals.GetGlobalRef().FuncInvokeGFunction(gfuncName, params)

	if result == nil { // TODO: Or check for error block
		return nil, fmt.Errorf("resolveFieldHandle: gfunction call to create MethodHandle failed for field %s.%s", className, fieldName)
	}

	return result.(*object.Object), nil
}

// GetPrimitiveClass returns the java.lang.Class object representing the
// primitive type specified by the descriptor (e.g., "I" for int).
// It retrieves this by looking up the static TYPE field in the corresponding
// wrapper class (e.g., java/lang/Integer).
func GetPrimitiveClass(descriptor string) *object.Object {
	var wrapperClass string
	switch descriptor {
	case "B":
		wrapperClass = "java/lang/Byte"
	case "C":
		wrapperClass = "java/lang/Character"
	case "D":
		wrapperClass = "java/lang/Double"
	case "F":
		wrapperClass = "java/lang/Float"
	case "I":
		wrapperClass = "java/lang/Integer"
	case "J":
		wrapperClass = "java/lang/Long"
	case "S":
		wrapperClass = "java/lang/Short"
	case "Z":
		wrapperClass = "java/lang/Boolean"
	case "V":
		wrapperClass = "java/lang/Void"
	default:
		return nil
	}

	// The primitive Class objects are stored in the TYPE static field of their wrapper classes.
	staticField, ok := statics.QueryStatic(wrapperClass, "TYPE")
	if !ok || staticField.Value == nil {
		return nil
	}

	primClassObj, ok := staticField.Value.(*object.Object)
	if !ok {
		return nil
	}

	return primClassObj
}

// resolveMethodHandleEntry resolves a method invocation handle (kinds 5-9)
func resolveMethodHandleEntry(cp *CPool, refIndex int, isStatic bool, isSpecial bool, fr *frames.Frame, refKind uint8) (*object.Object, error) {
	// 1. Resolve the method reference
	// refIndex points to MethodRef (10) or InterfaceMethodRef (11)
	if refIndex < 1 || refIndex >= len(cp.CpIndex) {
		return nil, fmt.Errorf("resolveMethodHandleEntry: invalid method ref index %d", refIndex)
	}

	cpEntry := cp.CpIndex[refIndex]
	var className, methodName, methodSig string

	if cpEntry.Type == MethodRef {
		// Use the resolved method refs if available, or look them up
		// In Jacobin, cp.MethodRefs holds the raw indices, cp.ResolvedMethodRefs holds resolved strings
		// We can use the helper function from cpUtils.go
		className, methodName, methodSig, _ = GetMethInfoFromCPmethref(cp, refIndex)
	} else if cpEntry.Type == Interface {
		className, methodName, methodSig = GetMethInfoFromCPinterfaceRef(cp, refIndex)
	} else {
		return nil, fmt.Errorf("resolveMethodHandleEntry: expected MethodRef or InterfaceMethodRef at index %d, got %d", refIndex, cpEntry.Type)
	}

	if globals.TraceClass {
		trace.Trace(fmt.Sprintf("resolveMethodHandleEntry: Class=%s, Method=%s, Sig=%s, Static=%v",
			className, methodName, methodSig, isStatic))
	}

	// 2. Get java.lang.Class object for the defining class.
	defClassObj, err := getClassObj("L"+className+";", fr)
	if err != nil {
		return nil, fmt.Errorf("resolveMethodHandleEntry: could not get Class object for %s: %w", className, err)
	}

	// 3. Get java.lang.invoke.MethodType object for the method signature.
	methodTypeObj, err := getMethodTypeObject(methodSig, fr)
	if err != nil {
		return nil, fmt.Errorf("resolveMethodHandleEntry: could not create MethodType for %s: %w", methodSig, err)
	}

	// 4. Get java.lang.Class object for the caller class (for access checks).
	callerClassObj, err := getClassObj("L"+fr.ClName+";", fr)
	if err != nil {
		return nil, fmt.Errorf("resolveMethodHandleEntry: could not get Class object for caller %s: %w", fr.ClName, err)
	}

	// 5. Create Java String for method name
	methodNameObj := object.StringObjectFromGoString(methodName)

	// 6. Invoke an internal gfunction to create the MethodHandle.
	// This is a hypothetical internal API for the VM to create method handles
	// without going through the full MethodHandles.Lookup security checks.
	params := []interface{}{
		defClassObj,
		methodNameObj,
		methodTypeObj,
		int64(refKind),
		callerClassObj,
	}

	gfuncName := "jacobin/internal/VM.resolveMethodHandle(Ljava/lang/Class;Ljava/lang/String;Ljava/lang/invoke/MethodType;ILjava/lang/Class;)Ljava/lang/invoke/MethodHandle;"
	result := globals.GetGlobalRef().FuncInvokeGFunction(gfuncName, params)

	if result == nil { // TODO: Or check for error block
		return nil, fmt.Errorf("resolveMethodHandleEntry: gfunction call to create MethodHandle failed for method %s.%s%s", className, methodName, methodSig)
	}

	return result.(*object.Object), nil
}

// ResolveMethodType resolves a MethodType constant pool entry.
// It parses the descriptor string (e.g. "(Ljava/lang/String;)V") and creates
// a java.lang.invoke.MethodType object.
func ResolveMethodType(cp *CPool, index int, fr *frames.Frame) (*object.Object, error) {
	// 1. Get the descriptor string from the CP
	// MethodType entry contains an index to a UTF8 string
	mtEntry := FetchCPentry(cp, index)
	if mtEntry.EntryType != MethodType {
		return nil, fmt.Errorf("ResolveMethodType: CP entry at %d is not a MethodType", index)
	}

	// FetchCPentry for MethodType returns the descriptor index as IntVal (cast to int64)
	descIndex := int(mtEntry.IntVal)
	descriptor := FetchUTF8stringFromCPEntryNumber(cp, uint16(descIndex))

	if globals.TraceClass {
		trace.Trace(fmt.Sprintf("ResolveMethodType: Resolving MT at index %d, descriptor=%s", index, descriptor))
	}

	// 2. Create the MethodType object using the helper
	return getMethodTypeObject(descriptor, fr)
}

// ResolveCallSite is the high-level function called by the INVOKEDYNAMIC instruction.
// It coordinates the resolution of the bootstrap method and the creation of the CallSite.
func ResolveCallSite(cp *CPool, index int, fr *frames.Frame) (*object.Object, error) {
	// index is the index into the constant pool for the CONSTANT_InvokeDynamic_info entry

	// 1. Fetch the InvokeDynamic entry (it was previously validated in codeCheck.go)
	idEntry := FetchCPentry(cp, index)

	// idEntry.AddrVal.entry1 is the bootstrap_method_attr_index
	// idEntry.AddrVal.entry2 is the name_and_type_index
	bsmIndex := int(idEntry.AddrVal.entry1)
	natIndex := int(idEntry.AddrVal.entry2)

	// 2. Get the Bootstrap Method info from the class attributes

	// In interpreter.go, fr.CP is an interface{}, usually *classloader.CPool.
	// The frame also has the class name fr.ClName. We can look up the class in MethArea.
	klass := MethAreaFetch(fr.ClName)
	if klass == nil {
		return nil, fmt.Errorf("ResolveCallSite: could not find class %s", fr.ClName)
	}

	bsm := cp.Bootstraps[bsmIndex]

	if globals.TraceClass {
		trace.Trace(fmt.Sprintf("ResolveCallSite: BSM index=%d, MethodRef=%d, ArgCount=%d",
			bsmIndex, bsm.MethodRef, len(bsm.Args)))
	}

	// 3. Resolve the Bootstrap Method Handle
	// bsm.MethodRef is an index into the Constant Pool (MethodHandle)
	bsmHandle, err := ResolveMethodHandle(cp, int(bsm.MethodRef), fr)
	if err != nil {
		return nil, err
	}

	// 4. Resolve the NameAndType (method name and type for the CallSite)
	// natIndex points to NameAndType entry
	// We need to create a String for the name and a MethodType for the type.
	// ...

	// 5. Resolve Static Arguments
	// bsm.Args is a list of indices into the Constant Pool.
	// These must be resolved to Java objects (String, Class, MethodType, MethodHandle, int, long, etc.)
	// ...

	// 6. Invoke the Bootstrap Method
	// This is the critical step: executing the BSM to get the CallSite object.
	// ...

	_ = bsmHandle // suppress unused var error for now
	_ = natIndex

	return nil, fmt.Errorf("ResolveCallSite: implementation pending")
}

// getMethodTypeObject creates a java.lang.invoke.MethodType object from a descriptor string.
func getMethodTypeObject(descriptor string, fr *frames.Frame) (*object.Object, error) {
	// Create the descriptor string object
	descriptorObj := object.StringObjectFromGoString(descriptor)

	// The class loader is used to resolve class names in the descriptor.
	// TODO: Pass the correct class loader from the frame/context. For now, we
	// pass nil, which corresponds to the bootstrap class loader.
	params := []interface{}{descriptorObj, nil}

	// We invoke the gfunction for java.lang.invoke.MethodType.fromMethodDescriptorString
	result := globals.GetGlobalRef().FuncInvokeGFunction(
		"java/lang/invoke/MethodType.fromMethodDescriptorString(Ljava/lang/String;Ljava/lang/ClassLoader;)Ljava/lang/invoke/MethodType;",
		params,
	)

	// TODO: check for errBlk and proceed accordingly.
	if result == nil {
		return nil, fmt.Errorf("getMethodTypeObject: failed to create MethodType for descriptor '%s'", descriptor)
	}

	if errBlock, ok := result.(*object.Object); ok && errBlock.KlassName == 0 { // A guess at how error blocks are identified
		return nil, fmt.Errorf("getMethodTypeObject: fromMethodDescriptorString threw an exception for '%s'", descriptor)
	}

	return result.(*object.Object), nil
}

// getClassObj gets a java.lang.Class object for a given class name or descriptor.
// It handles primitive types, array types, and object types by calling the equivalent
// of Class.forName() via a gfunction.
func getClassObj(descriptor string, fr *frames.Frame) (*object.Object, error) {
	// Check for primitive types (single-character descriptors). The VM pre-loads
	// Class objects for primitive types (e.g., Integer.TYPE).
	if len(descriptor) == 1 {
		primClass := GetPrimitiveClass(descriptor)
		if primClass != nil {
			return primClass, nil
		}
	}

	// The java.lang.Class.forName() method expects class names in "binary name"
	// format as defined by JLS §13.1.
	// - "java.lang.String"
	// - "[I" for int[]
	// - "[Ljava.lang.String;" for String[]
	// We must convert the field descriptor format to the binary name format.
	var forNameArg string
	if strings.HasPrefix(descriptor, "L") && strings.HasSuffix(descriptor, ";") {
		// It's an object type, e.g., "Ljava/lang/Object;".
		// forName expects "java.lang.Object".
		internalName := descriptor[1 : len(descriptor)-1]
		forNameArg = strings.ReplaceAll(internalName, "/", ".")
	} else {
		// It's an array type (e.g., "[I", "[Ljava/lang/Object;") or a primitive
		// that wasn't found in the cache. forName handles these formats directly,
		// but we need to convert '/' to '.'.
		forNameArg = strings.ReplaceAll(descriptor, "/", ".")
	}

	// Create a Java String for the class name
	nameObj := object.StringObjectFromGoString(forNameArg)

	// We use the gfunction for Class.forName(String, boolean, ClassLoader).
	// We must initialize the class, as per JVM spec for 'ldc' resolution.
	// TODO: Pass the correct class loader. For now, nil uses the bootstrap loader.
	params := []interface{}{
		nameObj,
		true, // initialize
		nil,  // class loader
	}
	gfuncName := "java/lang/Class.forName(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;"

	result := globals.GetGlobalRef().FuncInvokeGFunction(gfuncName, params)

	// TODO: check for errBlk and handle exceptions like ClassNotFoundException.
	if result == nil {
		return nil, fmt.Errorf("getClassObj: Class.forName failed for '%s'", forNameArg)
	}

	return result.(*object.Object), nil
}
