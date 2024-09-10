//go:build windows

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) Consult jacobin.org.
 */

package native

import (
	"errors"
	"fmt"
	"github.com/omarghader/pefile-go/pe"
	"jacobin/excNames"
	"jacobin/exceptions"
	"runtime"
	"syscall"
)

// load a DLL into memory, after which we'll extract the function names
func loadDll(path string) *syscall.DLL {
	if runtime.GOOS != "windows" {
		errMsg := "DLL operations are only supported on Windows"
		exceptions.ThrowEx(excNames.UnsupportedOperationException, errMsg, nil)
		return nil // only occurs in testing
	}
	if path == "" {
		errMsg := "empty dll path in loadDll()"
		exceptions.ThrowEx(excNames.IOException, errMsg, nil)
		return nil // only occurs in testing
	}

	dllPtr, err := syscall.LoadDLL(path)
	if err != nil {
		errMsg := fmt.Sprintf("in loadDll() could not load DLL %s", path)
		exceptions.ThrowEx(excNames.IOException, errMsg, nil)
		return nil // only occurs in testing
	}

	return dllPtr
}

// deletes the loded DLL from memory
func unloadDll(dllPtr *syscall.DLL) error {
	if runtime.GOOS != "windows" {
		errMsg := "DLL operations are only supported on Windows"
		exceptions.ThrowEx(excNames.UnsupportedOperationException, errMsg, nil)
		return errors.New(errMsg) // only occurs in testing
	}

	if dllPtr == nil {
		errMsg := "empty dll path in unloadDll()"
		exceptions.ThrowEx(excNames.IOException, errMsg, nil)
		return errors.New(errMsg) // only occurs in testing
	}

	err := dllPtr.Release()
	if err != nil {
		errMsg := fmt.Sprintf("error releasing DLL %s", dllPtr.Name)
		exceptions.ThrowEx(excNames.IOException, errMsg, nil)
		return errors.New(errMsg) // only occurs in testing
	}

	return nil
}

func CreateNativeFunctionTable(filename string) error {
	pefile, err := pe.NewPEFile(filename)
	if err != nil {
		errMsg := fmt.Sprintf("error parsing DLL file %s", filename)
		exceptions.ThrowEx(excNames.FileNotFoundException, errMsg, nil)
		return errors.New(errMsg)
	}

	for _, entry := range pefile.ExportDirectory.Exports {
		result := fmt.Sprintf("%s,%s", filename, entry.Name)
		fmt.Println(result)
	}
	return nil
}
