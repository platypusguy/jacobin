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
	"jacobin/src/trace"
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

	refKind := mhEntry.AddrVal.entry1
	refIndex := mhEntry.AddrVal.entry2

	if globals.TraceClass {
		trace.Trace(fmt.Sprintf("ResolveMethodHandle: Resolving MH at index %d, kind=%d, refIndex=%d", index, refKind, refIndex))
	}

	// The resolution logic varies significantly based on the reference kind.
	// See JVM Spec 5.4.3.5. Method Type and Method Handle Resolution
	switch refKind {
	case 1: // REF_getField
		return resolveFieldHandle(cp, int(refIndex), false, false, fr)
	case 2: // REF_getStatic
		return resolveFieldHandle(cp, int(refIndex), true, false, fr)
	case 3: // REF_putField
		return resolveFieldHandle(cp, int(refIndex), false, true, fr)
	case 4: // REF_putStatic
		return resolveFieldHandle(cp, int(refIndex), true, true, fr)
	case 5: // REF_invokeVirtual
		return resolveMethodHandleEntry(cp, int(refIndex), false, false, fr)
	case 6: // REF_invokeStatic
		return resolveMethodHandleEntry(cp, int(refIndex), true, false, fr)
	case 7: // REF_invokeSpecial
		// TODO: Special handling for <init> vs other methods?
		return resolveMethodHandleEntry(cp, int(refIndex), false, true, fr)
	case 8: // REF_newInvokeSpecial
		return resolveMethodHandleEntry(cp, int(refIndex), false, true, fr) // Constructor
	case 9: // REF_invokeInterface
		return resolveMethodHandleEntry(cp, int(refIndex), false, false, fr)
	default:
		return nil, fmt.Errorf("ResolveMethodHandle: invalid reference kind %d", refKind)
	}
}

// resolveFieldHandle resolves a field access handle (kinds 1-4)
func resolveFieldHandle(cp *CPool, refIndex int, isStatic bool, isSetter bool, fr *frames.Frame) (*object.Object, error) {
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

	// 2. Load the class containing the field
	// This ensures the class is loaded and we can access its metadata
	if err := LoadClassFromNameOnly(className); err != nil {
		return nil, err
	}

	// 3. Create a java.lang.invoke.MethodHandle object representing this field access
	// This involves creating a DirectMethodHandle (or similar internal subclass)
	// that knows how to get/put the field.

	// TODO: Implement the actual creation of the MethodHandle object.
	// For now, we will return a placeholder or throw an error if the MH class isn't ready.
	// We need to instantiate java.lang.invoke.DirectMethodHandle (or similar).

	// For the purpose of this step, we'll assume we need to return a valid Object pointer.
	// In a full implementation, this would be a fully initialized MethodHandle.
	// Since we are building this incrementally, we might need to stub this out.

	// Placeholder: Return null for now until we have the MH classes loaded and ready to instantiate
	return nil, fmt.Errorf("resolveFieldHandle: implementation pending for field handles")
}

// resolveMethodHandleEntry resolves a method invocation handle (kinds 5-9)
func resolveMethodHandleEntry(cp *CPool, refIndex int, isStatic bool, isSpecial bool, fr *frames.Frame) (*object.Object, error) {
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

	// 2. Load the class
	if err := LoadClassFromNameOnly(className); err != nil {
		return nil, err
	}

	// 3. Create the MethodHandle object
	// This requires mapping the method info to a MemberName and then to a MethodHandle.
	// This is a complex interaction with the JDK's java.lang.invoke code.

	// TODO: Implement creation of MethodHandle object for methods.
	return nil, fmt.Errorf("resolveMethodHandleEntry: implementation pending for method handles")
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

	// 2. Call java.lang.invoke.MethodType.fromMethodDescriptorString()
	// This is the standard way to create a MethodType from a string.
	// Create the descriptor string object
	descriptorObj := object.StringObjectFromGoString(descriptor)

	// For now, we can pass nil for the class loader  b/c we don't presently support custom class loaders.
	params := []interface{}{descriptorObj, nil}

	// We invoke the gfunction logic directly.
	result := globals.GetGlobalRef().FuncInvokeGFunction(
		"java/lang/invoke/MethodType.fromMethodDescriptorString(Ljava/lang/String;Ljava/lang/ClassLoader;)Ljava/lang/invoke/MethodType;",
		params,
	)

	// TODO: check for errBlk and proceed accordingly.
	if result == nil {
		return nil, fmt.Errorf("ResolveMethodType: failed to create MethodType")
	}

	if errBlock, ok := result.(*object.Object); ok && errBlock.KlassName == 0 { // Check if it's an error block?
		// Jacobin error blocks usually have a specific structure.
		// For now, assume if it returns an object it's the MethodType.
		return result.(*object.Object), nil
	}

	// If result is an error block (which is an *object.Object in Jacobin gfunctions usually)
	// we might need to check. But let's assume success for now.
	return result.(*object.Object), nil
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
