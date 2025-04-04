/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/statics"
	"jacobin/stringPool"
	"jacobin/types"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestSystemClinit(t *testing.T) {
	globals.InitGlobals("test")
	classloader.InitMethodArea()
	classloader.MethAreaInsert("java/lang/System", &classloader.Klass{Data: &classloader.ClData{ClInit: types.ClInitRun}})
	ret := systemClinit(nil)
	if ret != nil {
		gErr := ret.(*GErrBlk)
		t.Errorf("TestSystemClinit: Unexpected error message. got %s", gErr.ErrMsg)
	}
	t.Log("TestSystemClinit: stringClinit() returned nil as expected")
}

func TestArrayCopyNonOverlapping(t *testing.T) {
	globals.InitGlobals("test")

	src := object.Make1DimArray(object.INT, 10)
	dest := object.Make1DimArray(object.INT, 10)

	rawSrcArray := src.FieldTable["value"].Fvalue.([]int64)
	for i := 0; i < 10; i++ {
		rawSrcArray[i] = int64(1)
	}

	params := make([]interface{}, 5)
	params[0] = src
	params[1] = int64(2)
	params[2] = dest
	params[3] = int64(0)
	params[4] = int64(5)

	err := arrayCopy(params)

	if err != nil {
		e := err.(error)
		t.Errorf("Unexpected error in test of arrayCopy(): %s", error.Error(e))
	}

	rawDestArray := dest.FieldTable["value"].Fvalue.([]int64)
	j := int64(0)
	for i := 0; i < 10; i++ {
		j += rawDestArray[i]
	}

	if j != 5 {
		t.Errorf("Expected total to be 5, got %d", j)
	}

	if rawDestArray[0] != 1 || rawDestArray[5] != 0 {
		t.Errorf("Expedting [0] to be 1, [5] to be 0, got %d, %d",
			rawDestArray[0], rawDestArray[5])
	}
}

func TestArrayCopyOverlappingSameArray(t *testing.T) {
	globals.InitGlobals("test")

	src := object.Make1DimArray(object.BYTE, 10)
	// dest := object.Make1DimArray(object.BYTE, 10)

	rawSrcArray := src.FieldTable["value"].Fvalue.([]types.JavaByte)
	for i := 0; i < 10; i++ {
		rawSrcArray[i] = types.JavaByte(i)
	}

	params := make([]interface{}, 5)
	params[0] = src
	params[1] = int64(2)
	params[2] = src
	params[3] = int64(0)
	params[4] = int64(5)

	// result should be 2,3,4,5,6,5,6,7,8,9 (which totals 55)

	err := arrayCopy(params)

	if err != nil {
		e := err.(error)
		t.Errorf("Unexpected error in test of arrayCopy(): %s", error.Error(e))
	}

	j := types.JavaByte(0)
	for i := 0; i < 10; i++ {
		j += rawSrcArray[i]
	}

	if j != 55 {
		t.Errorf("Expected total to be 55, got %d", j)
	}
}

func TestArrayInvalidParmCount(t *testing.T) {
	globals.InitGlobals("test")

	src := object.Make1DimArray(object.INT, 10)
	dest := object.Make1DimArray(object.INT, 10)

	rawSrcArray := src.FieldTable["value"].Fvalue.([]int64)
	for i := 0; i < 10; i++ {
		rawSrcArray[i] = int64(1)
	}

	params := make([]interface{}, 4)
	params[0] = src
	params[1] = int64(2)
	params[2] = dest
	params[3] = int64(0)
	// params[4] = int64(5)

	err := arrayCopy(params)

	if err == nil {
		t.Errorf("Expecting error, but got none")
	}

	errMsg := err.(*GErrBlk).ErrMsg
	if !strings.Contains(errMsg, "Expected 5 parameters") {
		t.Errorf("Expected error re 5 parameters, got %s", errMsg)
	}
}

func TestArrayCopyInvalidPos(t *testing.T) {
	globals.InitGlobals("test")

	src := object.Make1DimArray(object.INT, 10)
	dest := object.Make1DimArray(object.INT, 10)

	rawSrcArray := src.FieldTable["value"].Fvalue.([]int64)
	for i := 0; i < 10; i++ {
		rawSrcArray[i] = int64(1)
	}

	params := make([]interface{}, 5)
	params[0] = src
	params[1] = int64(-1) // this is an invalid position in the array
	params[2] = dest
	params[3] = int64(0)
	params[4] = int64(5)

	err := arrayCopy(params)

	if err == nil {
		t.Errorf("Exoected an error message, but got none")
	}

	errMsg := err.(*GErrBlk).ErrMsg
	if !strings.Contains(errMsg, "Negative position") {
		t.Errorf("Expected error re invalid position, got %s", errMsg)
	}
}

func TestArrayCopyNullArray(t *testing.T) {
	globals.InitGlobals("test")

	src := object.Make1DimArray(object.INT, 10)
	dest := object.Make1DimArray(object.INT, 10)

	rawSrcArray := src.FieldTable["value"].Fvalue.([]int64)
	for i := 0; i < 10; i++ {
		rawSrcArray[i] = int64(1)
	}

	params := make([]interface{}, 5)
	params[0] = object.Null // clearly invalid
	params[1] = int64(2)
	params[2] = dest
	params[3] = int64(0)
	params[4] = int64(5)

	err := arrayCopy(params)

	if err == nil {
		t.Errorf("Exoected an error message, but got none")
	}

	errMsg := err.(*GErrBlk).ErrMsg
	if !strings.Contains(errMsg, "null src or dest") {
		t.Errorf("Expected error re null array, got %s", errMsg)
	}
}

func TestArrayCopyInvalidObject(t *testing.T) {
	globals.InitGlobals("test")

	src := object.Make1DimArray(object.INT, 10)
	dest := object.Make1DimArray(object.INT, 10)

	rawSrcArray := src.FieldTable["value"].Fvalue.([]int64)
	for i := 0; i < 10; i++ {
		rawSrcArray[i] = int64(1)
	}

	params := make([]interface{}, 5)
	o := object.MakeEmptyObject()
	objType := "invalid object"
	o.KlassName = stringPool.GetStringIndex(&objType)
	params[0] = o
	params[1] = int64(2)
	params[2] = dest
	params[3] = int64(0)
	params[4] = int64(5)

	err := arrayCopy(params)

	if err == nil {
		t.Errorf("Exoected an error message, but got none")
	}

	errMsg := err.(*GErrBlk).ErrMsg
	if !strings.Contains(errMsg, "invalid src or dest array") {
		t.Errorf("Expected error re invalid array type, got %s", errMsg)
	}
}

func TestArrayCopyInvalidLength(t *testing.T) {
	globals.InitGlobals("test")

	src := object.Make1DimArray(object.INT, 10)
	dest := object.Make1DimArray(object.INT, 10)

	rawSrcArray := src.FieldTable["value"].Fvalue.([]int64)
	for i := 0; i < 10; i++ {
		rawSrcArray[i] = int64(1)
	}

	params := make([]interface{}, 5)
	params[0] = src
	params[1] = int64(2)
	params[2] = dest
	params[3] = int64(0)
	params[4] = int64(200) // the invalid length

	err := arrayCopy(params)

	if err == nil {
		t.Errorf("Exoected an error message, but got none")
	}

	errMsg := err.(*GErrBlk).ErrMsg
	if !strings.Contains(errMsg, "Array position + length exceeds array size") {
		t.Errorf("Expected error re invalid length, got %s", errMsg)
	}
}

func TestGetMilliTime(t *testing.T) {
	globals.InitGlobals("test")
	ret := currentTimeMillis(nil).(int64)
	if ret < 1739512706877 { // milli time on 13 Feb 2025 at roughtly 10PM PST
		t.Errorf("Expected a greater value from nanoTime(), got %d", ret)
	}
}

func TestGetNanoTime(t *testing.T) {
	globals.InitGlobals("test")
	ret := nanoTime(nil).(int64)
	if ret < 1739512706877498200 { // nanotime on 13 Feb 2025 at roughtly 10PM PST
		t.Errorf("Expected a greater value from nanoTime(), got %d", ret)
	}
}

func TestExitI(t *testing.T) {
	globals.InitGlobals("test")
	ret := exitI([]interface{}{int64(17)})
	if ret.(int64) != 17 {
		t.Errorf("Expected exit code of 17, got %d", ret.(int64))
	}
}

func TestGetConsole(t *testing.T) {
	globals.InitGlobals("test")

	// systemClnit() will initialize stdin
	classloader.InitMethodArea()
	classloader.MethAreaInsert("java/lang/System", &classloader.Klass{Data: &classloader.ClData{ClInit: types.ClInitNotRun}})
	statics.PreloadStatics()
	_ = systemClinit(nil)

	ret := getConsole(nil)
	if ret.(*os.File) != os.Stdin {
		t.Errorf("Expected getConsole() to return stdin, got %v", ret)
	}
}

// the various property retrievals tested next

func TestGetProperty_FileEncoding(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("file.encoding")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString(globals.GetGlobalRef().FileEncoding)
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_FileSeparator(t *testing.T) {
	propObj := object.StringObjectFromGoString("file.separator")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString(string(os.PathSeparator))
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaClassPath(t *testing.T) {
	propObj := object.StringObjectFromGoString("java.class.path")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString(".")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaCompiler(t *testing.T) {
	propObj := object.StringObjectFromGoString("java.compiler")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString("no JIT")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaHome(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.home")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString(globals.GetGlobalRef().JavaHome)
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaIoTmpdir(t *testing.T) {
	propObj := object.StringObjectFromGoString("java.io.tmpdir")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString(os.TempDir())
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaLibraryPath(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.library.path")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString(globals.GetGlobalRef().JavaHome)
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVendor(t *testing.T) {
	propObj := object.StringObjectFromGoString("java.vendor")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString("Jacobin")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVendorUrl(t *testing.T) {
	propObj := object.StringObjectFromGoString("java.vendor.url")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString("https://jacobin.org")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVendorVersion(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.vendor.version")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString(globals.GetGlobalRef().Version)
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVersion(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.version")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString(strconv.Itoa(globals.GetGlobalRef().MaxJavaVersion))
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVmName(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.vm.name")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString(
		fmt.Sprintf("Jacobin VM v. %s (Java %d) 64-bit VM", globals.GetGlobalRef().Version, globals.GetGlobalRef().MaxJavaVersion))
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVmSpecificationName(t *testing.T) {
	propObj := object.StringObjectFromGoString("java.vm.specification.name")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString("Java Virtual Machine Specification")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVmSpecificationVendor(t *testing.T) {
	propObj := object.StringObjectFromGoString("java.vm.specification.vendor")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString("Oracle and Jacobin")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVmSpecificationVersion(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.vm.specification.version")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString(strconv.Itoa(globals.GetGlobalRef().MaxJavaVersion))
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVmVendor(t *testing.T) {
	propObj := object.StringObjectFromGoString("java.vm.vendor")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString("Jacobin")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVmVersion(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.vm.version")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString(strconv.Itoa(globals.GetGlobalRef().MaxJavaVersion))
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_LineSeparator(t *testing.T) {
	propObj := object.StringObjectFromGoString("line.separator")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	var expected string
	if runtime.GOOS == "windows" {
		expected = "\\r\\n"
	} else {
		expected = "\\n"
	}
	if object.GoStringFromStringObject(result.(*object.Object)) != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_NativeEncoding(t *testing.T) {
	propObj := object.StringObjectFromGoString("native.encoding")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString(globals.GetCharsetName())
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_OsArch(t *testing.T) {
	propObj := object.StringObjectFromGoString("os.arch")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString(runtime.GOARCH)
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_OsName(t *testing.T) {
	propObj := object.StringObjectFromGoString("os.name")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString(runtime.GOOS)
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_OsVersion(t *testing.T) {
	propObj := object.StringObjectFromGoString("os.version")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString("not yet available")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_PathSeparator(t *testing.T) {
	propObj := object.StringObjectFromGoString("path.separator")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected := object.StringObjectFromGoString(string(os.PathSeparator))
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_UserDir(t *testing.T) {
	propObj := object.StringObjectFromGoString("user.dir")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	expected, _ := os.Getwd()
	if object.GoStringFromStringObject(result.(*object.Object)) != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_UserHome(t *testing.T) {
	propObj := object.StringObjectFromGoString("user.home")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	currentUser, _ := user.Current()
	expected := currentUser.HomeDir
	if object.GoStringFromStringObject(result.(*object.Object)) != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_UserName(t *testing.T) {
	propObj := object.StringObjectFromGoString("user.name")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	currentUser, _ := user.Current()
	expected := currentUser.Name
	if object.GoStringFromStringObject(result.(*object.Object)) != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_UserTimezone(t *testing.T) {
	propObj := object.StringObjectFromGoString("user.timezone")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	now := time.Now()
	expected, _ := now.Zone()
	if object.GoStringFromStringObject(result.(*object.Object)) != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_Default(t *testing.T) {
	propObj := object.StringObjectFromGoString("unknown.property")
	params := []interface{}{propObj}
	result := systemGetProperty(params)
	if result != object.Null {
		t.Errorf("Expected null, got %v", result)
	}
}
