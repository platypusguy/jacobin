/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"fmt"
	"jacobin/src/frames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/trace"
	"jacobin/src/types"
	"jacobin/src/util"
	"strings"
)

// ResolveCallSite is the high-level function called by the INVOKEDYNAMIC instruction.
// It coordinates the resolution of the bootstrap method and the creation of the CallSite.
func ResolveCallSite(cp *CPool, index int, fr *frames.Frame) (*object.Object, error) {
	// index is the index into the constant pool for the CONSTANT_InvokeDynamic_info entry

	// Fetch the InvokeDynamic entry (it was previously validated in codeCheck.go)
	idEntry := FetchCPentry(cp, index)

	// idEntry.AddrVal.entry1 is the bootstrap_method_attr_index
	// idEntry.AddrVal.entry2 is the name_and_type_index for the CallSite
	bsmIndex := int(idEntry.AddrVal.entry1)
	natIndex := int(idEntry.AddrVal.entry2)

	// Get the Bootstrap Method info from the class attributes
	klass := MethAreaFetch(fr.ClName)
	if klass == nil {
		return nil, fmt.Errorf("ResolveCallSite: could not find class %s", fr.ClName)
	}

	// Resolve the NameAndType (method name and type for the CallSite)
	// natIndex points to NameAndType entry
	methName, methType := GetNATfieldsFromCPindex(cp, natIndex)
	if methName == "" || methType == "" {
		return nil, fmt.Errorf("ResolveCallSite: could not resolve NameAndType entry at index %d", natIndex)
	}

	// Resolve the Bootstrap Method Handle
	// bsm.MethodRef is an index into the Constant Pool (MethodHandle)
	bsm := cp.Bootstraps[bsmIndex]
	if globals.TraceClass {
		trace.Trace(fmt.Sprintf("ResolveCallSite: BSM index=%d, MethodRef=%d, ArgCount=%d",
			bsmIndex, bsm.MethodRef, len(bsm.Args)))
	}

	bsmHandle, err := ResolveMethodHandle(cp, int(bsm.MethodRef), fr)
	if err != nil {
		return nil, err
	}

	// Retrieve BSM signature to determine argument layout
	mhCpe := FetchCPentry(cp, int(bsm.MethodRef))
	if mhCpe.EntryType != MethodHandle {
		return nil, fmt.Errorf("ResolveCallSite: CP entry at %d is not a MethodHandle", bsm.MethodRef)
	}
	refIndex := int(mhCpe.AddrVal.entry2)
	var bsmSig string
	if refIndex >= 1 && refIndex < len(cp.CpIndex) && cp.CpIndex[refIndex].Type == Interface {
		_, _, bsmSig = GetMethInfoFromCPinterfaceRef(cp, refIndex)
	} else {
		_, _, bsmSig, _ = GetMethInfoFromCPmethref(cp, refIndex)
	}

	// Resolve Static Arguments
	// bsm.Args is a list of indices into the Constant Pool.
	// These must be resolved to Java objects
	staticArgs := make([]*object.Object, len(bsm.Args))
	for i, argIndex := range bsm.Args {
		staticArg, err := ResolveStaticArgument(cp, int(argIndex), fr)
		if err != nil {
			return nil, fmt.Errorf("ResolveCallSite: error resolving static argument at index %d: %w", argIndex, err)
		}
		staticArgs[i] = staticArg
	}

	// Marshal the arguments for the bootstrap method.
	// Per JVMS 5.4.3.6, the BSM is invoked with a Lookup object for the caller class, the call site name,
	// the call site MethodType, then the resolved static arguments in order.
	// The resolves arguments to the BSM are always static and if primitives, always boxed.
	// When the BSM expects varargs (e.g. Object... constants), trailing static arguments are packed into an array.
	bsmArgs, err := marshalBSMargs(fr, methName, methType, bsmSig, staticArgs)
	if err != nil {
		return nil, err
	}

	// Invoke the Bootstrap Method to obtain the CallSite object.
	return invokeBSM(cp, bsm, bsmHandle, bsmArgs)
}

// marshalBSMargs assembles the argument list for bootstrap method invocation.
// When the target bootstrap method signature expects varargs (e.g. [Ljava/lang/Object;),
// any trailing static arguments beyond the fixed parameters are packed into an array
// (or an empty array/slice if no extra constants exist).
func marshalBSMargs(fr *frames.Frame, methName, methType, bsmSig string, staticArgs []*object.Object) ([]any, error) {
	// arg 0: a Lookup object bound to the caller class (for caller-sensitive BSMs).
	callerClassObj, err := getClassObj("L"+fr.ClName+";", fr)
	if err != nil {
		return nil, fmt.Errorf("marshalBSMargs: could not get Class object for caller %s: %w", fr.ClName, err)
	}

	// arg 1: the call site name.
	nameObj := object.StringObjectFromGoString(methName)

	// arg 2: the call site type as a MethodType object.
	typeObj, err := getMethodTypeObject(methType, fr)
	if err != nil {
		return nil, fmt.Errorf("marshalBSMargs: %w", err)
	}

	// Check if the BSM signature expects 5 parameters ending with [Ljava/lang/Object;
	// e.g. StringConcatFactory.makeConcatWithConstants(Lookup, String, MethodType, String, Object...)
	// which has descriptor (
	//   Ljava/lang/invoke/MethodHandles$Lookup;
	//   Ljava/lang/String;Ljava/lang/invoke/MethodType;
	//   Ljava/lang/String;[Ljava/lang/Object;)
	//      Ljava/lang/invoke/CallSite;
	if strings.Contains(bsmSig, "[Ljava/lang/Object;)") {
		params := make([]any, 5)
		params[0] = makeLookupObject(callerClassObj)
		params[1] = nameObj
		params[2] = typeObj

		// Parameter 3 is the recipe string
		if len(staticArgs) > 0 {
			params[3] = staticArgs[0]
		} else {
			params[3] = object.StringObjectFromGoString("")
		}

		// Parameter 4 is Object... constants (pack staticArgs[1:] into an array/slice)
		var extraConstants []*object.Object
		if len(staticArgs) > 1 {
			extraConstants = staticArgs[1:]
		} else {
			extraConstants = make([]*object.Object, 0)
		}
		params[4] = extraConstants

		return params, nil
	}

	// Standard fallback for non-varargs BSMs
	params := make([]any, 0, 3+len(staticArgs))
	params = append(params, makeLookupObject(callerClassObj), nameObj, typeObj)
	for _, sa := range staticArgs {
		params = append(params, sa)
	}
	return params, nil
}

// makeLookupObject creates a minimal java.lang.invoke.MethodHandles$Lookup object
// bound to the caller class. Jacobin does not yet implement the Lookup API; this
// object exists so bootstrap methods receive the caller context JVMS 5.4.3.6 requires.
func makeLookupObject(callerClassObj *object.Object) *object.Object {
	lookupClassName := "java/lang/invoke/MethodHandles$Lookup"
	lookupObj := object.MakeEmptyObjectWithClassName(&lookupClassName)
	lookupObj.FieldTable["lookupClass"] = object.Field{Ftype: types.Ref, Fvalue: callerClassObj}
	return lookupObj
}

// invokeBSM invokes the bootstrap method handle with the marshaled arguments and
// returns the resulting CallSite object.
//
// The handle's $target field holds the executable MTentry. Only BSMs backed by
// gfunctions (MType == 'G') are supported: a BSM implemented as a Java method would
// require re-entering the interpreter, which is not available from the classloader.
// (All of Jacobin's current BSMs -- StringConcatFactory -- are gfunctions.)
func invokeBSM(cp *CPool, bsm BootstrapMethod, bsmHandle *object.Object, params []any) (*object.Object, error) {
	// Per JVMS 4.7.23, the BSM handle should be kind 6 (invokeStatic) or 8 (newInvokeSpecial).
	kindField, ok := bsmHandle.FieldTable["Kind"]
	kind, kindOk := kindField.Fvalue.(int64)
	if !ok || !kindOk || (kind != 6 && kind != 8) {
		return nil, fmt.Errorf("invokeBSM: bootstrap method handle has unexpected kind %v", kindField.Fvalue)
	}

	// Reconstruct the BSM's fully-qualified name from the constant pool for the
	// gfunction invocation: MethodHandle entry -> refIndex -> MethodRef entry.
	mhCpe := FetchCPentry(cp, int(bsm.MethodRef))
	if mhCpe.EntryType != MethodHandle {
		return nil, fmt.Errorf("invokeBSM: CP entry at %d is not a MethodHandle", bsm.MethodRef)
	}
	refIndex := int(mhCpe.AddrVal.entry2)
	var bsmClass, bsmName, bsmSig string
	if refIndex >= 1 && refIndex < len(cp.CpIndex) && cp.CpIndex[refIndex].Type == Interface {
		bsmClass, bsmName, bsmSig = GetMethInfoFromCPinterfaceRef(cp, refIndex)
	} else {
		bsmClass, bsmName, bsmSig, _ = GetMethInfoFromCPmethref(cp, refIndex)
	}
	bsmFQN := bsmClass + "." + bsmName + bsmSig

	// The executable payload must be a gfunction (see doc comment above).
	targetField, ok := bsmHandle.FieldTable["$target"]
	if !ok {
		return nil, fmt.Errorf("invokeBSM: method handle has no $target")
	}
	mtEntry, ok := targetField.Fvalue.(MTentry)
	if !ok {
		return nil, fmt.Errorf("invokeBSM: method handle $target is not an MTentry")
	}
	if mtEntry.MType != 'G' {
		return nil, fmt.Errorf("invokeBSM: bootstrap methods implemented in Java are not yet supported (MType=%c)", mtEntry.MType)
	}

	result := globals.GetGlobalRef().FuncInvokeGFunction(bsmFQN, params)
	switch r := result.(type) {
	case *object.Object:
		// Per the java.lang.invoke specification, the BSM must return a non-null CallSite object.
		if r == nil || object.IsNull(r) {
			return nil, fmt.Errorf("invokeBSM: bootstrap method %s returned null", bsmFQN)
		}
		if className := stringPool.GetStringPointer(r.KlassName); className == nil ||
			*className != "java/lang/invoke/CallSite" {
			return nil, fmt.Errorf("invokeBSM: bootstrap method %s returned %v, not a CallSite",
				bsmFQN, klassNameOf(r))
		}
		return r, nil
	case *ghelpers.GErrBlk:
		return nil, fmt.Errorf("invokeBSM: bootstrap method %s returned error %s", bsmFQN, result.(*ghelpers.GErrBlk).ErrMsg)
	case nil:
		return nil, fmt.Errorf("invokeBSM: bootstrap method %s returned nil", bsmFQN)
	default:
		return nil, fmt.Errorf("invokeBSM: bootstrap method %s returned unexpected type %T", bsmFQN, result)
	}
}

// klassNameOf is a small helper for error messages.
func klassNameOf(obj *object.Object) string {
	if p := stringPool.GetStringPointer(obj.KlassName); p != nil {
		return *p
	}
	return "<?>"
}

// ResolveStaticArgument resolves a constant pool entry (used as a static argument
// for a bootstrap method) into a Java object.
//
// Per JVMS 4.7.23, a bootstrap argument must be a CONSTANT_String, CONSTANT_Class,
// CONSTANT_Integer, CONSTANT_Long, CONSTANT_Float, CONSTANT_Double,
// CONSTANT_MethodHandle, or CONSTANT_MethodType entry. Per the java.lang.invoke
// specification, each resolves to a live object before the BSM is invoked:
//
//	CP entry type         -> Java object
//	CONSTANT_String       -> java.lang.String
//	CONSTANT_Class        -> java.lang.Class (resolved, as if by ldc)
//	CONSTANT_Integer      -> java.lang.Integer (boxed)
//	CONSTANT_Long         -> java.lang.Long (boxed)
//	CONSTANT_Float        -> java.lang.Float (boxed)
//	CONSTANT_Double       -> java.lang.Double (boxed)
//	CONSTANT_MethodHandle -> java.lang.invoke.MethodHandle
//	CONSTANT_MethodType   -> java.lang.invoke.MethodType
func ResolveStaticArgument(cp *CPool, index int, fr *frames.Frame) (*object.Object, error) {
	CPe := FetchCPentry(cp, index)
	if CPe.EntryType == 0 { // Dummy = 0; FetchCPentry also returns 0 for a bad index
		return nil, fmt.Errorf("ResolveStaticArgument: invalid CP index %d", index)
	}

	switch CPe.EntryType {
	case IntConst: // CONSTANT_Integer -> box to java.lang.Integer
		return object.MakePrimitiveObject("java/lang/Integer", types.Int, CPe.IntVal), nil
	case LongConst: // CONSTANT_Long -> box to java.lang.Long
		return object.MakePrimitiveObject("java/lang/Long", types.Long, CPe.IntVal), nil
	case FloatConst: // CONSTANT_Float -> box to java.lang.Float
		return object.MakePrimitiveObject("java/lang/Float", types.Float, CPe.FloatVal), nil
	case DoubleConst: // CONSTANT_Double -> box to java.lang.Double
		return object.MakePrimitiveObject("java/lang/Double", types.Double, CPe.FloatVal), nil
	case StringConst, UTF8: // CONSTANT_String -> java.lang.String
		// Utf8Refs holds raw Modified UTF-8 bytes (see cpParser.go), so it
		// must be decoded here. Same logic as the IS_STRING_ADDR case of ldc.
		mutf8 := util.DecodeModifiedUTF8([]byte(*CPe.StringVal))
		return object.StringObjectFromGoString(string(mutf8)), nil
	case ClassRef: // CONSTANT_Class -> java.lang.Class, resolved as if by ldc
		// This matches ldc semantics in the interpreter: the class is loaded
		// (resolved) but NOT initialized -- a class literal does not run
		// <clinit> (JVMS 5.5).
		className := *CPe.StringVal
		if strings.HasPrefix(className, "[") {
			classObj, err := getClassObj(className, fr)
			if err != nil {
				return nil, fmt.Errorf("ResolveStaticArgument: could not resolve array class %s: %w", className, err)
			}
			return classObj, nil
		}
		if err := LoadClassFromNameOnly(className); err != nil {
			return nil, fmt.Errorf("ResolveStaticArgument: could not load class %s: %w", className, err)
		}
		cl := MethAreaFetch(className)
		if cl == nil {
			return nil, fmt.Errorf("ResolveStaticArgument: class %s loaded but not present in method area", className)
		}
		return cl.Data.ClassObject, nil
	case MethodType: // CONSTANT_MethodType -> java.lang.invoke.MethodType
		return ResolveMethodType(cp, index, fr)
	case MethodHandle: // CONSTANT_MethodHandle -> java.lang.invoke.MethodHandle
		return ResolveMethodHandle(cp, index, fr)
	default:
		return nil, fmt.Errorf(
			"ResolveStaticArgument: CP entry at index %d has type %d, which is not a valid bootstrap static argument (JVMS 4.7.23)",
			index, CPe.EntryType)
	}
}
