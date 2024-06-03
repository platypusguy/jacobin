/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package excNames

import "testing"

func details(t *testing.T, index int, expected string) {
	observed := JVMexceptionNames[index]
	if observed != expected {
		t.Errorf("JVMexceptionNames[%d]: expected %s, but observed %s\n", index, expected, observed)
	}
}

func TestExceptionTableAlignment(t *testing.T) {
	details(t, IllegalArgumentException, "java.lang.IllegalArgumentException")
	details(t, NoSuchDynamicMethodException, "jdk.dynalink.NoSuchDynamicMethodException")
	details(t, WrongMethodTypeException, "java.lang.invoke.WrongMethodTypeException")
	details(t, ClassNotLoadedException, "com.sun.jdi.ClassNotLoadedException")
	details(t, InvalidTypeException, "com.sun.jdi.InvalidTypeException")
	details(t, PrintException, "javax.print.PrintException")
	details(t, UnmodifiableClassException, "java.lang.instrument.UnmodifiableClassException")
	details(t, XMLParseException, "javax.management.modelmbean.XMLParseException")
	details(t, VirtualMachineError, "java.lang.VirtualMachineError")
	details(t, UTFDataFormatException, "java.io.UTFDataFormatException")
}
