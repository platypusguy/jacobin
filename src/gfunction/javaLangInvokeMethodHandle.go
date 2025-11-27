/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import "jacobin/src/classloader"

// Internal representation of java.lang.invoke.MethodHandle
type MethodHandle struct {
	Kind          MethodHandleKind    // REF_getField, REF_invokeVirtual, etc.
	RefClass      string              // Declaring class
	RefName       string              // Method/field name
	RefDescriptor string              // Method/field descriptor
	DirectMethod  *classloader.Method // For direct method invocations
	IsVarArgs     bool
}

// Method handle reference kinds (JVM spec §5.4.3.5)
type MethodHandleKind uint16

const (
	REF_getField         MethodHandleKind = 1
	REF_getStatic        MethodHandleKind = 2
	REF_putField         MethodHandleKind = 3
	REF_putStatic        MethodHandleKind = 4
	REF_invokeVirtual    MethodHandleKind = 5
	REF_invokeStatic     MethodHandleKind = 6
	REF_invokeSpecial    MethodHandleKind = 7
	REF_newInvokeSpecial MethodHandleKind = 8
	REF_invokeInterface  MethodHandleKind = 9
)

// CallSite represents a resolved invokedynamic call site
type CallSite struct {
	Target     *MethodHandle // The method handle to invoke
	Type       *MethodType   // Expected signature
	IsVolatile bool          // MutableCallSite vs ConstantCallSite
}

// MethodType represents a method signature
type MethodType struct {
	ReturnType string
	ParamTypes []string
}
