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
	"jacobin/globals"
	"jacobin/log"
	"jacobin/util"
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
	gl := *globals.GetGlobalRef()
	jh := gl.JavaHome
	dllList := util.SearchDirByFileExtension(jh, "dll")
	dllListSize := len(*dllList)
	var functionListSize = 0

	for _, dllFile := range *dllList {
		pefile, err := pe.NewPEFile(dllFile)

		if err != nil {
			errMsg := fmt.Sprintf("error parsing DLL file %s", filename)
			exceptions.ThrowEx(excNames.FileNotFoundException, errMsg, nil)
			return errors.New(errMsg)
		}

		for _, entry := range pefile.ExportDirectory.Exports {
			if entry.Name[0] != '?' { // guessing that natives that start wih ? are not real functions
				result := fmt.Sprintf("%s,%s", dllFile, entry.Name)
				fmt.Println(result)
				functionListSize += 1
			}
		}

	}

	summary := fmt.Sprintf(
		"Native function table for Windows created: %d native functions in %d .dll files",
		functionListSize, dllListSize)
	log.Log(summary, log.FINE)
	return nil
}
