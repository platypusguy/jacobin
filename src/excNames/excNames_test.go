/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package excNames

import (
	"testing"
)

// Make sure that the enum for the exception correctly matches the string. Exceptions tested here
// are spaced more or less evenly through the list of exceptions
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
	details(t, NoSuchAlgorithmException, "java.security.NoSuchAlgorithmException")
}

// Make sure that the enum for the exception correctly matches the string. Exceptions tested here
// are spaced more or less evenly through the list of exceptions. This tests the Jacobin-specific
// exceptions. The test is similar to the one above, but it uses the Jacobin-default exception names
func TestJacobinExceptionTableAlignment(t *testing.T) {
	detailsJacobin(t, IllegalArgumentException, "java.lang.IllegalArgumentException")
	detailsJacobin(t, NoSuchDynamicMethodException, "jdk.dynalink.NoSuchDynamicMethodException")
	detailsJacobin(t, WrongMethodTypeException, "java.lang.invoke.WrongMethodTypeException")
	detailsJacobin(t, ClassNotLoadedException, "org.jacobin.ClassNotLoadedException")
	detailsJacobin(t, InvalidTypeException, "org.jacobin.InvalidTypeException")
	detailsJacobin(t, PrintException, "javax.print.PrintException")
	detailsJacobin(t, UnmodifiableClassException, "java.lang.instrument.UnmodifiableClassException")
	detailsJacobin(t, XMLParseException, "javax.management.modelmbean.XMLParseException")
	detailsJacobin(t, VirtualMachineError, "java.lang.VirtualMachineError")
	detailsJacobin(t, VMStartException, "org.jacobin.connect.VMStartException")
}

// utility function for previous unit test
func details(t *testing.T, index int, expected string) {
	observed := JVMexceptionNames[index]
	if observed != expected {
		t.Errorf("JVMexceptionNames[%d]: expected %s, but observed %s\n", index, expected, observed)
	}
}

func detailsJacobin(t *testing.T, index int, expected string) {
	observed := JVMexceptionNamesJacobin[index]
	if observed != expected {
		t.Errorf("JVMexceptionNamesJacobin[%d]: expected %s, but observed %s\n", index, expected, observed)
	}
}
