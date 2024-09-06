/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2022 by Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)
 */

package jvm

import (
	"io"
	"jacobin/classloader"
	"jacobin/gfunction"
	"jacobin/globals"
	"jacobin/log"
	"jacobin/thread"
	"os"
	"strconv"
	"strings"
	"testing"
)

// These tests use the byte array corresponding to Hello2.class, which computes a series of small
// numbers and prints them to stdout.

var Hello2Bytes = []byte{
	0xCA, 0xFE, 0xBA, 0xBE, 0x00, 0x00, 0x00, 0x37, 0x00, 0x2B, 0x07, 0x00, 0x02, 0x01, 0x00, 0x06,
	0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x32, 0x07, 0x00, 0x04, 0x01, 0x00, 0x10, 0x6A, 0x61, 0x76, 0x61,
	0x2F, 0x6C, 0x61, 0x6E, 0x67, 0x2F, 0x4F, 0x62, 0x6A, 0x65, 0x63, 0x74, 0x01, 0x00, 0x06, 0x3C,
	0x69, 0x6E, 0x69, 0x74, 0x3E, 0x01, 0x00, 0x03, 0x28, 0x29, 0x56, 0x01, 0x00, 0x04, 0x43, 0x6F,
	0x64, 0x65, 0x0A, 0x00, 0x03, 0x00, 0x09, 0x0C, 0x00, 0x05, 0x00, 0x06, 0x01, 0x00, 0x0F, 0x4C,
	0x69, 0x6E, 0x65, 0x4E, 0x75, 0x6D, 0x62, 0x65, 0x72, 0x54, 0x61, 0x62, 0x6C, 0x65, 0x01, 0x00,
	0x12, 0x4C, 0x6F, 0x63, 0x61, 0x6C, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6C, 0x65, 0x54, 0x61,
	0x62, 0x6C, 0x65, 0x01, 0x00, 0x04, 0x74, 0x68, 0x69, 0x73, 0x01, 0x00, 0x08, 0x4C, 0x48, 0x65,
	0x6C, 0x6C, 0x6F, 0x32, 0x3B, 0x01, 0x00, 0x04, 0x6D, 0x61, 0x69, 0x6E, 0x01, 0x00, 0x16, 0x28,
	0x5B, 0x4C, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x6C, 0x61, 0x6E, 0x67, 0x2F, 0x53, 0x74, 0x72, 0x69,
	0x6E, 0x67, 0x3B, 0x29, 0x56, 0x0A, 0x00, 0x01, 0x00, 0x11, 0x0C, 0x00, 0x12, 0x00, 0x13, 0x01,
	0x00, 0x06, 0x61, 0x64, 0x64, 0x54, 0x77, 0x6F, 0x01, 0x00, 0x05, 0x28, 0x49, 0x49, 0x29, 0x49,
	0x09, 0x00, 0x15, 0x00, 0x17, 0x07, 0x00, 0x16, 0x01, 0x00, 0x10, 0x6A, 0x61, 0x76, 0x61, 0x2F,
	0x6C, 0x61, 0x6E, 0x67, 0x2F, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6D, 0x0C, 0x00, 0x18, 0x00, 0x19,
	0x01, 0x00, 0x03, 0x6F, 0x75, 0x74, 0x01, 0x00, 0x15, 0x4C, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x69,
	0x6F, 0x2F, 0x50, 0x72, 0x69, 0x6E, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6D, 0x3B, 0x0A, 0x00,
	0x1B, 0x00, 0x1D, 0x07, 0x00, 0x1C, 0x01, 0x00, 0x13, 0x6A, 0x61, 0x76, 0x61, 0x2F, 0x69, 0x6F,
	0x2F, 0x50, 0x72, 0x69, 0x6E, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6D, 0x0C, 0x00, 0x1E, 0x00,
	0x1F, 0x01, 0x00, 0x07, 0x70, 0x72, 0x69, 0x6E, 0x74, 0x6C, 0x6E, 0x01, 0x00, 0x04, 0x28, 0x49,
	0x29, 0x56, 0x01, 0x00, 0x04, 0x61, 0x72, 0x67, 0x73, 0x01, 0x00, 0x13, 0x5B, 0x4C, 0x6A, 0x61,
	0x76, 0x61, 0x2F, 0x6C, 0x61, 0x6E, 0x67, 0x2F, 0x53, 0x74, 0x72, 0x69, 0x6E, 0x67, 0x3B, 0x01,
	0x00, 0x01, 0x78, 0x01, 0x00, 0x01, 0x49, 0x01, 0x00, 0x01, 0x69, 0x01, 0x00, 0x0D, 0x53, 0x74,
	0x61, 0x63, 0x6B, 0x4D, 0x61, 0x70, 0x54, 0x61, 0x62, 0x6C, 0x65, 0x07, 0x00, 0x21, 0x01, 0x00,
	0x01, 0x6A, 0x01, 0x00, 0x01, 0x6B, 0x01, 0x00, 0x0A, 0x53, 0x6F, 0x75, 0x72, 0x63, 0x65, 0x46,
	0x69, 0x6C, 0x65, 0x01, 0x00, 0x0B, 0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x32, 0x2E, 0x6A, 0x61, 0x76,
	0x61, 0x00, 0x20, 0x00, 0x01, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00,
	0x05, 0x00, 0x06, 0x00, 0x01, 0x00, 0x07, 0x00, 0x00, 0x00, 0x2F, 0x00, 0x01, 0x00, 0x01, 0x00,
	0x00, 0x00, 0x05, 0x2A, 0xB7, 0x00, 0x08, 0xB1, 0x00, 0x00, 0x00, 0x02, 0x00, 0x0A, 0x00, 0x00,
	0x00, 0x06, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x0B, 0x00, 0x00, 0x00, 0x0C, 0x00, 0x01,
	0x00, 0x00, 0x00, 0x05, 0x00, 0x0C, 0x00, 0x0D, 0x00, 0x00, 0x00, 0x09, 0x00, 0x0E, 0x00, 0x0F,
	0x00, 0x01, 0x00, 0x07, 0x00, 0x00, 0x00, 0x81, 0x00, 0x03, 0x00, 0x03, 0x00, 0x00, 0x00, 0x1E,
	0x03, 0x3D, 0xA7, 0x00, 0x15, 0x1C, 0x1C, 0x04, 0x64, 0xB8, 0x00, 0x10, 0x3C, 0xB2, 0x00, 0x14,
	0x1B, 0xB6, 0x00, 0x1A, 0x84, 0x02, 0x01, 0x1C, 0x10, 0x0A, 0xA1, 0xFF, 0xEB, 0xB1, 0x00, 0x00,
	0x00, 0x03, 0x00, 0x0A, 0x00, 0x00, 0x00, 0x16, 0x00, 0x05, 0x00, 0x00, 0x00, 0x06, 0x00, 0x05,
	0x00, 0x07, 0x00, 0x0D, 0x00, 0x08, 0x00, 0x14, 0x00, 0x06, 0x00, 0x1D, 0x00, 0x0A, 0x00, 0x0B,
	0x00, 0x00, 0x00, 0x20, 0x00, 0x03, 0x00, 0x00, 0x00, 0x1E, 0x00, 0x20, 0x00, 0x21, 0x00, 0x00,
	0x00, 0x0D, 0x00, 0x0A, 0x00, 0x22, 0x00, 0x23, 0x00, 0x01, 0x00, 0x02, 0x00, 0x1B, 0x00, 0x24,
	0x00, 0x23, 0x00, 0x02, 0x00, 0x25, 0x00, 0x00, 0x00, 0x0F, 0x00, 0x02, 0xFF, 0x00, 0x05, 0x00,
	0x03, 0x07, 0x00, 0x26, 0x00, 0x01, 0x00, 0x00, 0x11, 0x00, 0x08, 0x00, 0x12, 0x00, 0x13, 0x00,
	0x01, 0x00, 0x07, 0x00, 0x00, 0x00, 0x38, 0x00, 0x02, 0x00, 0x02, 0x00, 0x00, 0x00, 0x04, 0x1A,
	0x1B, 0x60, 0xAC, 0x00, 0x00, 0x00, 0x02, 0x00, 0x0A, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00,
	0x00, 0x00, 0x0D, 0x00, 0x0B, 0x00, 0x00, 0x00, 0x16, 0x00, 0x02, 0x00, 0x00, 0x00, 0x04, 0x00,
	0x27, 0x00, 0x23, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x28, 0x00, 0x23, 0x00, 0x01, 0x00,
	0x01, 0x00, 0x29, 0x00, 0x00, 0x00, 0x02, 0x00, 0x2A,
}

func TestHexHello2ValidClass(t *testing.T) {

	if testing.Short() { // don't run if running quick tests only. (Used primarily so GitHub doesn't run and bork)
		t.Skip()
	}

	redirecting := true
	var normalStderr, rerr, werr *os.File
	var normalStdout, rout, wout *os.File
	var err error

	// redirect stderr & stdout to capture results from stderr
	if redirecting {
		// stderr
		normalStderr = os.Stderr
		rerr, werr, err = os.Pipe()
		if err != nil {
			t.Errorf("os.Pipe returned an error: %s", err.Error())
			return
		}
		os.Stderr = werr
		// stdout
		normalStdout = os.Stdout
		rout, wout, _ = os.Pipe()
		os.Stdout = wout
	}

	// Initialise global, logging, classloader
	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)
	t.Logf("globals.InitGlobals and log.Init ok\n")

	// Initialise classloader
	err = classloader.Init()
	if err != nil {
		t.Errorf("classloader.Init returned an error: %s\n", err.Error())
		return
	}
	t.Logf("classloader.Init ok\n")

	// Load the base classes
	classloader.LoadBaseClasses()
	t.Logf("LoadBaseClasses ok\n")

	// Show the map size and check it for java/lang/System
	mapSize := classloader.JmodMapSize()
	if mapSize < 1 {
		t.Errorf("map size < 1 (fatal error)")
		return
	}
	t.Logf("Map size is %d\n", mapSize)

	// Set up MethArea for Hello2
	eKI := classloader.Klass{
		Status: 'I', // I = initializing the load
		Loader: "",
		Data:   nil,
	}
	classloader.MethAreaInsert("Hello2", &eKI)

	// Load bytes for Hello2
	_, _, err = classloader.ParseAndPostClass(&classloader.BootstrapCL, "Hello2.class", Hello2Bytes)
	if err != nil {
		t.Errorf("Got error from classloader.ParseAndPostCLass: %s", error.Error(err))
		return
	} else {
		t.Logf("Loaded class Hello2\n")
	}

	// Run class Hello2
	classloader.MTable = make(map[string]classloader.MTentry)
	gfunction.MTableLoadGFunctions(&classloader.MTable)
	mainThread := thread.CreateThread()
	err = StartExec("Hello2", &mainThread, globals.GetGlobalRef())
	if err != nil {
		t.Errorf("Got error from StartExec(): %s", error.Error(err))
		return
	}

	t.Logf("StartExec(Hello2) succeeded\n")

	if redirecting {
		_ = werr.Close()
		_ = wout.Close()
		msgStderr, _ := io.ReadAll(rerr)
		msgStdout, _ := io.ReadAll(rout)
		os.Stderr = normalStderr
		os.Stdout = normalStdout
		if !strings.Contains(string(msgStdout), "-1") {
			t.Errorf("Error in output: expected to contain in part '-1', but saw stdout & stderr as follows:\nstdout: %s\nstderr: %s\n",
				string(msgStdout), string(msgStderr))
		}
	}
}

func TestHexHello2InvalidMagicNumber(t *testing.T) {

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)
	err := classloader.Init()

	testBytes := Hello2Bytes[0:2]

	_, _, err = classloader.ParseAndPostClass(&classloader.BootstrapCL, "Hello2", testBytes)
	if err == nil {
		t.Error("Expected an error, but got none.")
	}

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout

	if !strings.Contains(string(msg), "invalid magic number") {
		t.Errorf("Expected error message to contain in part 'invalid magic number', got: %s", string(msg))
	}
}

func TestHexHello2InvalidJavaVersion(t *testing.T) {

	normalStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	normalStdout := os.Stdout
	_, wout, _ := os.Pipe()
	os.Stdout = wout

	globals.InitGlobals("test")
	log.Init()
	_ = log.SetLogLevel(log.WARNING)
	err := classloader.Init()

	var testBytes = Hello2Bytes[0:8]
	testBytes[7] = 99 // change class to unsupported version of Java class files

	_, _, err = classloader.ParseAndPostClass(&classloader.BootstrapCL, "Hello2", testBytes)
	if err == nil {
		t.Error("Expected an error, but got none.")
	}

	_ = w.Close()
	msg, _ := io.ReadAll(r)
	os.Stderr = normalStderr

	_ = wout.Close()
	os.Stdout = normalStdout

	versionString := "Java " + strconv.Itoa(globals.GetGlobalRef().MaxJavaVersion)
	if !strings.Contains(string(msg), "supports only Java versions through "+versionString) {
		t.Errorf("Expected error message to contain 'supports only Java versions through %s', got: %s",
			versionString, string(msg))
	}
}
