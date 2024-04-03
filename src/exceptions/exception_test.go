package exceptions

import "testing"

func details(t *testing.T, index int, expected string) {
	observed := JVMexceptionNames[index]
	if observed != expected {
		t.Errorf("JVMexceptionNames[%d]: expected %s, but observed %s\n", index, expected, observed)
	}
}

func TestExceptionTableAlignment(t *testing.T) {
	details(t, IllegalArgumentException, "java.lang.IllegalArgumentException")
	details(t, NoSuchDynamicMethodException, "java.lang.NoSuchDynamicMethodException")
	details(t, WrongMethodTypeException, "java.lang.WrongMethodTypeException")
	details(t, ClassNotLoadedException, "java.lang.ClassNotLoadedException")
	details(t, InvalidTypeException, "java.lang.InvalidTypeException")
	details(t, PrintException, "java.lang.PrintException")
	details(t, UnmodifiableClassException, "java.lang.UnmodifiableClassException")
	details(t, XMLParseException, "java.lang.XMLParseException")
	details(t, VirtualMachineError, "java.lang.VirtualMachineError")
	details(t, UTFDataFormatException, "java.lang.UTFDataFormatException")
}
