/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package javaLang

import (
	"bytes"
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaUtil"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/stringPool"
	"jacobin/src/types"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestArrayCopyNonOverlapping(t *testing.T) {
	globals.InitGlobals("test")

	src := object.Make1DimArray(object.T_INT, 10)
	dest := object.Make1DimArray(object.T_INT, 10)

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

	err := systemArrayCopy(params)

	if err != nil {
		e := err.(error)
		t.Errorf("Unexpected error in test of systemArrayCopy(): %s", error.Error(e))
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

	src := object.Make1DimArray(object.T_BYTE, 10)
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

	err := systemArrayCopy(params)

	if err != nil {
		e := err.(error)
		t.Errorf("Unexpected error in test of systemArrayCopy(): %s", error.Error(e))
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

	src := object.Make1DimArray(object.T_INT, 10)
	dest := object.Make1DimArray(object.T_INT, 10)

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

	err := systemArrayCopy(params)

	if err == nil {
		t.Errorf("Expecting error, but got none")
	}

	errMsg := err.(*ghelpers.GErrBlk).ErrMsg
	if !strings.Contains(errMsg, "Expected 5 parameters") {
		t.Errorf("Expected error re 5 parameters, got %s", errMsg)
	}
}

func TestArrayCopyInvalidPos(t *testing.T) {
	globals.InitGlobals("test")

	src := object.Make1DimArray(object.T_INT, 10)
	dest := object.Make1DimArray(object.T_INT, 10)

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

	err := systemArrayCopy(params)

	if err == nil {
		t.Errorf("Exoected an error message, but got none")
	}

	errMsg := err.(*ghelpers.GErrBlk).ErrMsg
	if !strings.Contains(errMsg, "Negative position") {
		t.Errorf("Expected error re invalid position, got %s", errMsg)
	}
}

func TestArrayCopyNullArray(t *testing.T) {
	globals.InitGlobals("test")

	src := object.Make1DimArray(object.T_INT, 10)
	dest := object.Make1DimArray(object.T_INT, 10)

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

	err := systemArrayCopy(params)

	if err == nil {
		t.Errorf("Exoected an error message, but got none")
	}

	errMsg := err.(*ghelpers.GErrBlk).ErrMsg
	if !strings.Contains(errMsg, "null src or dest") {
		t.Errorf("Expected error re null array, got %s", errMsg)
	}
}

func TestArrayCopyInvalidObject(t *testing.T) {
	globals.InitGlobals("test")

	src := object.Make1DimArray(object.T_INT, 10)
	dest := object.CloneObject(src)

	objType := "invalid object"
	src.KlassName = stringPool.GetStringIndex(&objType)

	params := make([]interface{}, 5)
	params[0] = src
	params[1] = int64(2)
	params[2] = dest
	params[3] = int64(0)
	params[4] = int64(5)

	err := systemArrayCopy(params)

	if err == nil {
		t.Errorf("Exoected an error message, but got none")
	}

	errMsg := err.(*ghelpers.GErrBlk).ErrMsg
	if !strings.Contains(errMsg, "wrong srcType") {
		t.Errorf("Expected error re wrong srcType, got %s", errMsg)
	}
}

func TestArrayCopyInvalidLength(t *testing.T) {
	globals.InitGlobals("test")

	src := object.Make1DimArray(object.T_INT, 10)
	dest := object.Make1DimArray(object.T_INT, 10)

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

	err := systemArrayCopy(params)

	if err == nil {
		t.Errorf("Exoected an error message, but got none")
	}

	errMsg := err.(*ghelpers.GErrBlk).ErrMsg
	if !strings.Contains(errMsg, "Array position + length exceeds array size") {
		t.Errorf("Expected error re invalid length, got %s", errMsg)
	}
}

func TestGetMilliTime(t *testing.T) {
	globals.InitGlobals("test")
	ret := systemCurrentTimeMillis(nil).(int64)
	if ret < 1739512706877 { // milli time on 13 Feb 2025 at roughtly 10PM PST
		t.Errorf("Expected a greater value from nanoTime(), got %d", ret)
	}
}

func TestGetNanoTime(t *testing.T) {
	globals.InitGlobals("test")
	ret := systemNanoTime(nil).(int64)
	if ret < 1739512706877498200 { // nanotime on 13 Feb 2025 at roughtly 10PM PST
		t.Errorf("Expected a greater value from nanoTime(), got %d", ret)
	}
}

func TestExitI(t *testing.T) {
	globals.InitGlobals("test")
	ret := systemExitI([]interface{}{int64(17)})
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

	ret := systemConsole(nil)
	if ret.(*os.File) != os.Stdin {
		t.Errorf("Expected getConsole() to return stdin, got %v", ret)
	}
}

// the various property retrievals tested next

func TestGetProperty_FileEncoding(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("file.encoding")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString(globals.GetGlobalRef().FileEncoding)
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_FileSeparator(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("file.separator")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString(string(os.PathSeparator))
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaClassPath(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.class.path")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString(".")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaCompiler(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.compiler")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString("no JIT")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaHome(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.home")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString(globals.GetGlobalRef().JavaHome)
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaIoTmpdir(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.io.tmpdir")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString(os.TempDir())
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaLibraryPath(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.library.path")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString(globals.GetGlobalRef().JavaHome)
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVendor(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.vendor")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString("Jacobin")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVendorUrl(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.vendor.url")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString("https://jacobin.org")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVendorVersion(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.vendor.version")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString(globals.GetGlobalRef().Version)
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVersion(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.version")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	_, verString := globals.GetJDKmajorVersion()
	expected := object.StringObjectFromGoString(verString)
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JDKmajorVersion(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("jdk.major.version")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	_, verString := globals.GetJDKmajorVersion()
	expected := object.StringObjectFromGoString(verString)
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVmName(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.vm.name")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString(
		fmt.Sprintf("Jacobin VM v. %s (Java %d) 64-bit VM", globals.GetGlobalRef().Version, globals.GetGlobalRef().MaxJavaVersion))
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVmSpecificationName(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.vm.specification.name")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString("Java Virtual Machine Specification")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVmSpecificationVendor(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.vm.specification.vendor")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString("Oracle and Jacobin")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVmSpecificationVersion(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.vm.specification.version")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString(strconv.Itoa(globals.GetGlobalRef().MaxJavaVersion))
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVmVendor(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.vm.vendor")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString("Jacobin")
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_JavaVmVersion(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.vm.version")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString(strconv.Itoa(globals.GetGlobalRef().MaxJavaVersion))
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_LineSeparator(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("line.separator")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	var expected string
	if runtime.GOOS == "windows" {
		expected = "\r\n"
	} else {
		expected = "\n"
	}
	if object.GoStringFromStringObject(result.(*object.Object)) != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_NativeEncoding(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("native.encoding")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString(globals.GetCharsetName())
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_OsArch(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("os.arch")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString(runtime.GOARCH)
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_OsName(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("os.name")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString(runtime.GOOS)
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_OsVersion(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("os.version")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	observed := object.GoStringFromStringObject(result.(*object.Object))
	errPhrase := "cmd.CombinedOutput() failed"
	if strings.Contains(observed, errPhrase) {
		t.Errorf("Was not expecting: %s", observed)
	}
}

func TestGetProperty_PathSeparator(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("path.separator")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected := object.StringObjectFromGoString(string(os.PathListSeparator))
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_UserDir(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("user.dir")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	expected, _ := os.Getwd()
	if object.GoStringFromStringObject(result.(*object.Object)) != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_UserHome(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("user.home")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	currentUser, _ := user.Current()
	expected := currentUser.HomeDir
	if object.GoStringFromStringObject(result.(*object.Object)) != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_UserName(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("user.name")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	currentUser, _ := user.Current()
	expected := currentUser.Name
	if object.GoStringFromStringObject(result.(*object.Object)) != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_UserTimezone(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("user.timezone")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	now := time.Now()
	expected, _ := now.Zone()
	if object.GoStringFromStringObject(result.(*object.Object)) != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetProperty_DefaultNull(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("unknown.property")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	if result != object.Null {
		t.Errorf("Expected null, got %v", result)
	}
}

func TestGetProperty_DefaultValue(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("hokey")
	dfltObj := object.StringObjectFromGoString("pokey")
	params := []interface{}{propObj, dfltObj}
	result := SystemGetProperty(params)
	expected := dfltObj
	if object.GoStringFromStringObject(result.(*object.Object)) != object.GoStringFromStringObject(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestSetProperty_JavaIoTmpdir(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.io.tmpdir")
	params := []interface{}{propObj}
	result := SystemGetProperty(params)
	observed := object.GoStringFromStringObject(result.(*object.Object))
	expected := os.TempDir()
	if observed != expected {
		t.Errorf("1. Expected %s, observed %s", expected, observed)
		return
	}

	newValueObj := object.StringObjectFromGoString("rubbish")
	params = []interface{}{propObj, newValueObj}
	result = systemSetProperty(params)
	observed = object.GoStringFromStringObject(result.(*object.Object))
	if observed != expected {
		t.Errorf("2. Expected %s, observed %s", expected, observed)
		return
	}
	params = []interface{}{propObj}
	result = SystemGetProperty(params)
	observed = object.GoStringFromStringObject(result.(*object.Object))
	expected = "rubbish"
	if observed != expected {
		t.Errorf("3. Expected %s, observed %s", expected, observed)
		return
	}

}

func TestClearProperty_JavaIoTmpdir(t *testing.T) {
	globals.InitGlobals("test")
	propObj := object.StringObjectFromGoString("java.io.tmpdir")
	params := []interface{}{propObj}
	result := systemClearProperty(params)
	observed := object.GoStringFromStringObject(result.(*object.Object))
	expected := os.TempDir()
	if observed != expected {
		t.Errorf("1. Expected %s, observed %s", expected, observed)
		return
	}

	result = SystemGetProperty(params)
	if !object.IsNull(result) {
		observed := object.GoStringFromStringObject(result.(*object.Object))
		if observed != "" {
			t.Errorf("Expected null, observed %s", observed)
			return
		}
	}

}

func TestSetProperties(t *testing.T) {
	globals.InitGlobals("test")

	// Create a Properties object.
	propsObj := object.MakeEmptyObjectWithClassName(&types.ClassNameProperties)
	params := []interface{}{propsObj}
	javaUtil.PropertiesInit(params)

	// Add 2 elements to the Properties object.
	key1Obj := object.StringObjectFromGoString("color")
	value1Obj := object.StringObjectFromGoString("green")
	params = []interface{}{propsObj, key1Obj, value1Obj}
	_ = javaUtil.PropertiesSetProperty(params)
	key2Obj := object.StringObjectFromGoString("sound")
	value2Obj := object.StringObjectFromGoString("soft")
	params = []interface{}{propsObj, key2Obj, value2Obj}
	_ = javaUtil.PropertiesSetProperty(params)

	// Replace the system properties.
	propsMap := propsObj.FieldTable[types.FieldNameProperties].Fvalue.(types.DefProperties)
	globals.ReplaceSystemProperties(propsMap)

	// Try to get tmp dir.
	propObj := object.StringObjectFromGoString("java.io.tmpdir")
	params = []interface{}{propObj}
	result := SystemGetProperty(params)
	if result != object.Null {
		t.Errorf("Expected object.Null, observed %v", result)
		return
	}

	// Try to get "color".
	propObj = object.StringObjectFromGoString("color")
	params = []interface{}{propObj}
	result = SystemGetProperty(params)
	observed := object.GoStringFromStringObject(result.(*object.Object))
	expected := "green"
	if observed != expected {
		t.Errorf("3. Expected %s, observed %s", expected, observed)
		return
	}

	// Try to get "sound".
	propObj = object.StringObjectFromGoString("sound")
	params = []interface{}{propObj}
	result = SystemGetProperty(params)
	observed = object.GoStringFromStringObject(result.(*object.Object))
	expected = "soft"
	if observed != expected {
		t.Errorf("3. Expected %s, observed %s", expected, observed)
		return
	}
}

func TestRunGc(t *testing.T) {
	globals.InitGlobals("test")
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	initialGcCount := memStats.NumGC

	systemForceGC(nil)

	runtime.ReadMemStats(&memStats)
	finalGcCount := memStats.NumGC

	if finalGcCount <= initialGcCount {
		t.Errorf("TestRunGc: GC Count does not reflect a run of GC. Initial: %d, Final: %d",
			initialGcCount, finalGcCount)
	}
}

func TestSystemGetProperty_CorrectnessAndEdgeCases(t *testing.T) {
	globals.InitGlobals("test")

	// 1. Existing property
	res := SystemGetProperty([]interface{}{object.StringObjectFromGoString("java.vendor")})
	if object.IsNull(res) {
		t.Fatalf("expected non-null for java.vendor")
	}
	if str := object.GoStringFromStringObject(res.(*object.Object)); str != "Jacobin" {
		t.Errorf("expected 'Jacobin', got '%s'", str)
	}

	// 2. Non-existent property without default -> returns object.Null
	res = SystemGetProperty([]interface{}{object.StringObjectFromGoString("non.existent.prop")})
	if res != object.Null {
		t.Errorf("expected object.Null for missing property, got %v", res)
	}

	// 3. Non-existent property with default -> returns default object
	defObj := object.StringObjectFromGoString("myDefaultValue")
	res = SystemGetProperty([]interface{}{object.StringObjectFromGoString("non.existent.prop"), defObj})
	if res != defObj {
		t.Errorf("expected default object %v, got %v", defObj, res)
	}

	// 4. Existing property with default -> returns existing value, not default
	res = SystemGetProperty([]interface{}{object.StringObjectFromGoString("java.vendor"), defObj})
	if res == defObj {
		t.Errorf("expected existing value, but got default object")
	}
	if str := object.GoStringFromStringObject(res.(*object.Object)); str != "Jacobin" {
		t.Errorf("expected 'Jacobin', got '%s'", str)
	}

	// 5. Non-existent property with null default -> returns object.Null
	res = SystemGetProperty([]interface{}{object.StringObjectFromGoString("non.existent.prop"), object.Null})
	if res != object.Null {
		t.Errorf("expected object.Null when def is object.Null, got %v", res)
	}
	res = SystemGetProperty([]interface{}{object.StringObjectFromGoString("non.existent.prop"), nil})
	if res != object.Null {
		t.Errorf("expected object.Null when def is nil, got %v", res)
	}

	// 6. Nil or Null key -> NullPointerException
	res = SystemGetProperty([]interface{}{nil})
	if errBlk, ok := res.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for nil key, got %v", res)
	}
	res = SystemGetProperty([]interface{}{object.Null})
	if errBlk, ok := res.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for object.Null key, got %v", res)
	}
	res = SystemGetProperty([]interface{}{})
	if errBlk, ok := res.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NullPointerException for empty params, got %v", res)
	}

	// 7. Empty key string -> IllegalArgumentException
	res = SystemGetProperty([]interface{}{object.StringObjectFromGoString("")})
	if errBlk, ok := res.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException for empty key, got %v", res)
	}
	res = SystemGetProperty([]interface{}{object.StringObjectFromGoString(""), defObj})
	if errBlk, ok := res.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("expected IllegalArgumentException for empty key with def, got %v", res)
	}

	// 8. systemGetPropertyDefault alias
	res = systemGetPropertyDefault([]interface{}{object.StringObjectFromGoString("non.existent.prop"), defObj})
	if res != defObj {
		t.Errorf("expected default object from systemGetPropertyDefault, got %v", res)
	}
}

func TestSystemSetOutAndSetErr(t *testing.T) {
	globals.InitGlobals("test")
	Load_Lang_System()

	// Verify method registrations
	setInSig := "java/lang/System.setIn(Ljava/io/InputStream;)V"
	setOutSig := "java/lang/System.setOut(Ljava/io/PrintStream;)V"
	setErrSig := "java/lang/System.setErr(Ljava/io/PrintStream;)V"

	gmIn, ok := ghelpers.MethodSignatures[setInSig]
	if !ok {
		t.Fatalf("Method %s not found in MethodSignatures", setInSig)
	}
	if gmIn.ParamSlots != 1 {
		t.Errorf("Expected ParamSlots 1 for %s, got %d", setInSig, gmIn.ParamSlots)
	}

	gmOut, ok := ghelpers.MethodSignatures[setOutSig]
	if !ok {
		t.Fatalf("Method %s not found in MethodSignatures", setOutSig)
	}
	if gmOut.ParamSlots != 1 {
		t.Errorf("Expected ParamSlots 1 for %s, got %d", setOutSig, gmOut.ParamSlots)
	}

	gmErr, ok := ghelpers.MethodSignatures[setErrSig]
	if !ok {
		t.Fatalf("Method %s not found in MethodSignatures", setErrSig)
	}
	if gmErr.ParamSlots != 1 {
		t.Errorf("Expected ParamSlots 1 for %s, got %d", setErrSig, gmErr.ParamSlots)
	}

	// Save original statics to restore later
	origIn := statics.GetStaticValue("java/lang/System", "in")
	origOut := statics.GetStaticValue("java/lang/System", "out")
	origErr := statics.GetStaticValue("java/lang/System", "err")
	defer func() {
		_ = statics.AddStatic("java/lang/System.in", statics.Static{Type: "GS", Value: origIn})
		_ = statics.AddStatic("java/lang/System.out", statics.Static{Type: "GS", Value: origOut})
		_ = statics.AddStatic("java/lang/System.err", statics.Static{Type: "GS", Value: origErr})
	}()

	// Test missing parameter error handling
	res := systemSetIn([]any{})
	if errBlk, ok := res.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for systemSetIn with empty params, got %v", res)
	}

	res = systemSetOut([]any{})
	if errBlk, ok := res.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for systemSetOut with empty params, got %v", res)
	}

	res = systemSetErr([]any{})
	if errBlk, ok := res.(*ghelpers.GErrBlk); !ok || errBlk.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("Expected IllegalArgumentException for systemSetErr with empty params, got %v", res)
	}

	// Test setting in to a new reader / object
	bufIn := bytes.NewReader([]byte("test input"))
	res = systemSetIn([]any{bufIn})
	if res != nil {
		t.Errorf("systemSetIn returned unexpected result: %v", res)
	}
	newIn := statics.GetStaticValue("java/lang/System", "in")
	if newIn != bufIn {
		t.Errorf("Expected System.in to be %v, got %v", bufIn, newIn)
	}

	// Test setting out to a new writer / object
	bufOut := new(bytes.Buffer)
	res = systemSetOut([]any{bufOut})
	if res != nil {
		t.Errorf("systemSetOut returned unexpected result: %v", res)
	}
	newOut := statics.GetStaticValue("java/lang/System", "out")
	if newOut != bufOut {
		t.Errorf("Expected System.out to be %v, got %v", bufOut, newOut)
	}

	// Test setting err to a new writer / object
	bufErr := new(bytes.Buffer)
	res = systemSetErr([]any{bufErr})
	if res != nil {
		t.Errorf("systemSetErr returned unexpected result: %v", res)
	}
	newErr := statics.GetStaticValue("java/lang/System", "err")
	if newErr != bufErr {
		t.Errorf("Expected System.err to be %v, got %v", bufErr, newErr)
	}

	// Test with *object.Object
	isClassName := "java/io/InputStream"
	isObj := object.MakeEmptyObjectWithClassName(&isClassName)

	res = gmIn.GFunction([]any{isObj})
	if res != nil {
		t.Errorf("Invoking setIn via GFunction returned error: %v", res)
	}
	if statics.GetStaticValue("java/lang/System", "in") != isObj {
		t.Errorf("Expected System.in to be isObj")
	}

	psClassName := "java/io/PrintStream"
	psObj := object.MakeEmptyObjectWithClassName(&psClassName)

	res = gmOut.GFunction([]any{psObj})
	if res != nil {
		t.Errorf("Invoking setOut via GFunction returned error: %v", res)
	}
	if statics.GetStaticValue("java/lang/System", "out") != psObj {
		t.Errorf("Expected System.out to be psObj")
	}

	res = gmErr.GFunction([]any{psObj})
	if res != nil {
		t.Errorf("Invoking setErr via GFunction returned error: %v", res)
	}
	if statics.GetStaticValue("java/lang/System", "err") != psObj {
		t.Errorf("Expected System.err to be psObj")
	}
}
