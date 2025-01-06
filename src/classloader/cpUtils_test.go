/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package classloader

import (
	"jacobin/frames"
	"jacobin/globals"
	"jacobin/types"
	"math"
	"testing"
)

func TestMeInfoFromMethRefValid(t *testing.T) {
	globals.InitGlobals("test")

	// Initialize classloaders and method area
	err := Init()
	if err != nil {
		t.Errorf("Failure to load classes in TestMeInfoFromMethRefValid")
	}
	LoadBaseClasses() // must follow classloader.Init()

	f := frames.CreateFrame(4)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := CPool{}
	CP.CpIndex = make([]CpEntry, 10)
	CP.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = CpEntry{Type: MethodRef, Slot: 0}

	CP.MethodRefs = make([]MethodRefEntry, 1)
	CP.MethodRefs[0] = MethodRefEntry{ClassIndex: 2, NameAndType: 3}

	CP.CpIndex[2] = CpEntry{Type: ClassRef, Slot: 0}
	CP.ClassRefs = make([]uint32, 4)
	CP.ClassRefs[0] = types.ObjectPoolStringIndex

	CP.CpIndex[3] = CpEntry{Type: NameAndType, Slot: 0}
	CP.NameAndTypes = make([]NameAndTypeEntry, 4)
	CP.NameAndTypes[0] = NameAndTypeEntry{
		NameIndex: 4,
		DescIndex: 5,
	}
	CP.CpIndex[4] = CpEntry{Type: UTF8, Slot: 0} // method name
	CP.Utf8Refs = make([]string, 4)
	CP.Utf8Refs[0] = "<init>"

	CP.CpIndex[5] = CpEntry{Type: UTF8, Slot: 1} // method name
	CP.Utf8Refs[1] = "()V"

	f.CP = &CP

	_, s2, s3 := GetMethInfoFromCPmethref(&CP, 1)
	if s2 != "<init>" && s3 != "()V" {
		t.Errorf("Expect to get a method: <init>()V, got %s%s", s2, s2)
	}

}

func TestGetClassNameFromCPclassref(t *testing.T) {
	globals.InitGlobals("test")

	// Initialize classloaders and method area
	err := Init()
	if err != nil {
		t.Errorf("Failure to load classes in TestGetClassNameFromCPclassref")
	}
	LoadBaseClasses() // must follow classloader.Init()

	f := frames.CreateFrame(4)
	f.Meth = append(f.Meth, 0x00)
	f.Meth = append(f.Meth, 0x01) // Go to slot 0x0001 in the CP

	CP := CPool{}
	CP.CpIndex = make([]CpEntry, 10)
	CP.CpIndex[0] = CpEntry{Type: 0, Slot: 0}
	CP.CpIndex[1] = CpEntry{Type: ClassRef, Slot: 0}
	CP.CpIndex[2] = CpEntry{Type: ClassRef, Slot: 1}

	CP.ClassRefs = make([]uint32, 4)
	CP.ClassRefs[0] = types.ObjectPoolStringIndex
	CP.ClassRefs[1] = types.InvalidStringIndex

	f.CP = &CP

	s1 := GetClassNameFromCPclassref(&CP, 1)
	if s1 != "java/lang/Object" {
		t.Errorf("Expect class name of 'java/lang/Object', got %s", s1)
	}

	s2 := GetClassNameFromCPclassref(&CP, 0)
	if s2 != "" {
		t.Errorf("Expected empty class name, got %s", s2)
	}
}

// test al CP lookups that do not return structs consisting of two CP entries
func TestFetchCPentry(t *testing.T) {
	globals.InitGlobals("test")

	cp := FetchCPentry(nil, 6)
	if cp.RetType != IS_ERROR {
		t.Errorf("Expect IS_ERROR, got %d", cp.RetType)
	}

	CP := CPool{
		CpIndex:        []CpEntry{},
		ClassRefs:      []uint32{},
		Doubles:        []float64{},
		Dynamics:       []DynamicEntry{},
		Floats:         []float32{},
		IntConsts:      []int32{},
		InterfaceRefs:  []InterfaceRefEntry{},
		InvokeDynamics: []InvokeDynamicEntry{},
		LongConsts:     []int64{},
		MethodHandles:  []MethodHandleEntry{},
		MethodRefs:     []MethodRefEntry{},
		MethodTypes:    []uint16{},
		NameAndTypes:   []NameAndTypeEntry{},
		Utf8Refs:       []string{},
	}

	CP.CpIndex = make([]CpEntry, 20)
	CP.CpIndex[0] = CpEntry{Type: 0, Slot: 0} // mandatory dummy entry
	CP.CpIndex[1] = CpEntry{Type: IntConst, Slot: 0}
	CP.CpIndex[2] = CpEntry{Type: LongConst, Slot: 0}
	CP.CpIndex[3] = CpEntry{Type: StringConst, Slot: 4}
	CP.CpIndex[4] = CpEntry{Type: UTF8, Slot: 0}
	CP.CpIndex[5] = CpEntry{Type: MethodType, Slot: 0}
	CP.CpIndex[6] = CpEntry{Type: FloatConst, Slot: 0}
	CP.CpIndex[7] = CpEntry{Type: DoubleConst, Slot: 0}
	CP.CpIndex[8] = CpEntry{Type: ClassRef, Slot: 0}
	CP.CpIndex[9] = CpEntry{Type: StringConst, Slot: 8} // causes an error, s/point to at UTF8 entry
	CP.CpIndex[10] = CpEntry{Type: Module, Slot: 0}     // should cause an error
	CP.CpIndex[11] = CpEntry{Type: Package, Slot: 0}    // should cause an error

	CP.IntConsts = []int32{25}
	cp = FetchCPentry(&CP, 1)
	if cp.RetType != IS_INT64 || cp.IntVal != int64(25) {
		t.Errorf("Expect IS_INT64 with value of 25, got %d with value of %d", cp.RetType, cp.IntVal)
	}

	CP.LongConsts = []int64{250}
	cp = FetchCPentry(&CP, 2)
	if cp.RetType != IS_INT64 || cp.IntVal != int64(250) {
		t.Errorf("Expect IS_INT64 with value of 250, got %d with value of %d", cp.RetType, cp.IntVal)
	}

	CP.Utf8Refs = []string{"Hello from Jacobin JVM!"}
	cp = FetchCPentry(&CP, 3) // get the string const that points to the string in Utf8Refs[0]
	if cp.RetType != IS_STRING_ADDR || *cp.StringVal != "Hello from Jacobin JVM!" {
		t.Errorf("Expect IS_STRING_ADDR pointing to 'Hello from Jacobin JVM!', got %s", *cp.StringVal)
	}

	cp = FetchCPentry(&CP, 9) // get an invalid string const
	if cp.RetType != IS_ERROR || cp.EntryType != 0 {
		t.Errorf("Expect IS_ERROR with value 0, got %d with value %d", cp.RetType, cp.EntryType)
	}

	CP.MethodTypes = []uint16{24}
	cp = FetchCPentry(&CP, 5)
	if cp.RetType != IS_INT64 || cp.IntVal != int64(24) {
		t.Errorf("Expect IS_INT64 with value of 24, got %d with value of %d", cp.RetType, cp.IntVal)
	}

	CP.Floats = []float32{float32(24.100000)}
	cp = FetchCPentry(&CP, 6)
	if cp.RetType != IS_FLOAT64 || math.Abs(24.1-cp.FloatVal) > 0.001 {
		t.Errorf("Expected IS_FLOAT64 with value of 24.1, got %d with value of %f", cp.RetType, cp.FloatVal)
	}

	CP.Doubles = []float64{24.20}
	cp = FetchCPentry(&CP, 7)
	if cp.RetType != IS_FLOAT64 || math.Abs(24.20-cp.FloatVal) > 0.001 {
		t.Errorf("Expected IS_FLOAT64 with value of 24.2, got %d with value of %f", cp.RetType, cp.FloatVal)
	}

	err := Init()
	if err != nil {
		t.Errorf("Failure to load classes in TestFetchCPentry")
	}
	LoadBaseClasses()
	CP.ClassRefs = []uint32{types.StringPoolStringIndex}
	cp = FetchCPentry(&CP, 8)
	if cp.RetType != IS_STRING_ADDR || *cp.StringVal != "java/lang/String" {
		t.Errorf("Expected IS_STRING_ADDR pointing to 'java/lang/String', got %d pointing to %s",
			cp.RetType, *cp.StringVal)
	}

	cp = FetchCPentry(&CP, 4) // UTF-8
	if cp.RetType != IS_STRING_ADDR || *cp.StringVal != "Hello from Jacobin JVM!" {
		t.Errorf("Expected IS_STRING_ADDR pointing to 'Hello from Jacobin JVM!', got %s", *cp.StringVal)
	}

	cp = FetchCPentry(&CP, 10) // Module, should cause an error
	if cp.RetType != IS_ERROR {
		t.Errorf("Expect IS_ERRROR for Module lookup, got %d", cp.RetType)
	}

	cp = FetchCPentry(&CP, 11) // Package, should cause an error
	if cp.RetType != IS_ERROR {
		t.Errorf("Expect IS_ERRROR for Package lookup, got %d", cp.RetType)
	}
}

func TestFetchCPentriesThatAreStructAddresses(t *testing.T) {
	globals.InitGlobals("test")

	CP := CPool{
		CpIndex:        []CpEntry{},
		ClassRefs:      []uint32{},
		Doubles:        []float64{},
		Dynamics:       []DynamicEntry{},
		Floats:         []float32{},
		IntConsts:      []int32{},
		InterfaceRefs:  []InterfaceRefEntry{},
		InvokeDynamics: []InvokeDynamicEntry{},
		LongConsts:     []int64{},
		MethodHandles:  []MethodHandleEntry{},
		MethodRefs:     []MethodRefEntry{},
		MethodTypes:    []uint16{},
		NameAndTypes:   []NameAndTypeEntry{},
		Utf8Refs:       []string{},
	}

	CP.CpIndex = make([]CpEntry, 10)
	CP.CpIndex[0] = CpEntry{Type: 0, Slot: 0} // mandatory dummy entry
	CP.CpIndex[1] = CpEntry{Type: Interface, Slot: 0}
	CP.CpIndex[2] = CpEntry{Type: Dynamic, Slot: 0}
	CP.CpIndex[3] = CpEntry{Type: InvokeDynamic, Slot: 0}
	CP.CpIndex[4] = CpEntry{Type: MethodHandle, Slot: 0}
	CP.CpIndex[5] = CpEntry{Type: MethodRef, Slot: 0}
	CP.CpIndex[6] = CpEntry{Type: NameAndType, Slot: 0}
	CP.CpIndex[7] = CpEntry{Type: InvokeDynamic, Slot: 0}

	// Interface
	irf := InterfaceRefEntry{
		ClassIndex:  4,
		NameAndType: 2,
	}
	CP.InterfaceRefs = []InterfaceRefEntry{irf}

	cp := FetchCPentry(&CP, 1)
	if cp.RetType != IS_STRUCT_ADDR {
		t.Errorf("Expected IS_STRUCT_ADDR, got %d", cp.RetType)
	}

	struc := *cp.AddrVal
	if struc.entry1 != uint16(4) || struc.entry2 != 2 {
		t.Errorf("Expected returned struc to contain 4 and 2, got %d and %d",
			struc.entry1, struc.entry2)
	}

	// Dynamic
	de := DynamicEntry{
		BootstrapIndex: 5,
		NameAndType:    3,
	}
	CP.Dynamics = []DynamicEntry{de}

	cp = FetchCPentry(&CP, 2)
	if cp.RetType != IS_STRUCT_ADDR {
		t.Errorf("Expected IS_STRUCT_ADDR, got %d", cp.RetType)
	}

	struc = *cp.AddrVal
	if struc.entry1 != uint16(5) || struc.entry2 != uint16(3) {
		t.Errorf("Expected returned struc to contain 5 and 3, got %d and %d",
			struc.entry1, struc.entry2)
	}

	// InvokeDynamic
	id := InvokeDynamicEntry{
		BootstrapIndex: 20,
		NameAndType:    21,
	}
	CP.InvokeDynamics = []InvokeDynamicEntry{id}

	cp = FetchCPentry(&CP, 7)
	if cp.RetType != IS_STRUCT_ADDR {
		t.Errorf("Expected IS_STRUCT_ADDR, got %d", cp.RetType)
	}

	struc = *cp.AddrVal
	if struc.entry1 != uint16(20) || struc.entry2 != uint16(21) {
		t.Errorf("Expected returned struc to contain 20 and 21, got %d and %d",
			struc.entry1, struc.entry2)
	}

	// Method Handle
	mh := MethodHandleEntry{
		RefKind:  8,
		RefIndex: 9,
	}
	CP.MethodHandles = []MethodHandleEntry{mh}

	cp = FetchCPentry(&CP, 4)
	if cp.RetType != IS_STRUCT_ADDR {
		t.Errorf("Expected IS_STRUCT_ADDR, got %d", cp.RetType)
	}

	struc = *cp.AddrVal
	if struc.entry1 != uint16(8) || struc.entry2 != uint16(9) {
		t.Errorf("Expected returned struc to contain 8 and 9, got %d and %d",
			struc.entry1, struc.entry2)
	}

	// Method Ref
	mr := MethodRefEntry{
		ClassIndex:  10,
		NameAndType: 11,
	}
	CP.MethodRefs = []MethodRefEntry{mr}

	cp = FetchCPentry(&CP, 5)
	if cp.RetType != IS_STRUCT_ADDR {
		t.Errorf("Expected IS_STRUCT_ADDR, got %d", cp.RetType)
	}

	struc = *cp.AddrVal
	if struc.entry1 != uint16(10) || struc.entry2 != uint16(11) {
		t.Errorf("Expected returned struc to contain 10 and 11, got %d and %d",
			struc.entry1, struc.entry2)
	}

	// NameAndType
	nt := NameAndTypeEntry{
		NameIndex: 12,
		DescIndex: 13,
	}
	CP.NameAndTypes = []NameAndTypeEntry{nt}

	cp = FetchCPentry(&CP, 6)
	if cp.RetType != IS_STRUCT_ADDR {
		t.Errorf("Expected IS_STRUCT_ADDR, got %d", cp.RetType)
	}

	struc = *cp.AddrVal
	if struc.entry1 != uint16(12) || struc.entry2 != uint16(13) {
		t.Errorf("Expected returned struc to contain 12 and 13, got %d and %d",
			struc.entry1, struc.entry2)
	}
}
