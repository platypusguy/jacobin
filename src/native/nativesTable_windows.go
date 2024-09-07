//go:build windows
// +build windows

/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) Consult jacobin.org.
 */

package native

import (
	"errors"
	"fmt"
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
