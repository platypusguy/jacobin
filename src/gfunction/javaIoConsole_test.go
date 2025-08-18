/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin Authors.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package gfunction

import (
	"io"
	"jacobin/src/classloader"
	"os"
	"reflect"
	"testing"

	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
)

func TestLoad_Io_Console_RegistersMethods(t *testing.T) {
	// Save and restore the global MethodSignatures map to avoid test pollution
	saved := MethodSignatures
	defer func() { MethodSignatures = saved }()

	MethodSignatures = make(map[string]GMeth)

	Load_Io_Console()

	// Expected keys and their characteristics
	checks := []struct {
		key   string
		slots int
		fn    func([]interface{}) interface{}
	}{
		{"java/io/Console.<clinit>()V", 0, consoleClinit},
		{"java/io/Console.charset()Ljava/nio/charset/Charset;", 0, trapFunction},
		{"java/io/Console.flush()V", 0, consoleFlush},
		{"java/io/Console.format(Ljava/lang/String;[Ljava/lang/Object;)Ljava/io/Console;", 2, consolePrintf},
		{"java/io/Console.printf(Ljava/lang/String;[Ljava/lang/Object;)Ljava/io/Console;", 2, consolePrintf},
		{"java/io/Console.reader()Ljava/io/Reader;", 0, trapFunction},
		{"java/io/Console.readLine()Ljava/lang/String;", 0, consoleReadLine},
		{"java/io/Console.readLine(Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/String;", 2, consolePrintfReadLine},
		{"java/io/Console.readPassword()[C", 0, consoleReadPassword},
		{"java/io/Console.readPassword(Ljava/lang/String;[Ljava/lang/Object;)[C", 2, consolePrintfReadPassword},
		{"java/io/Console.writer()Ljava/io/PrintWriter;", 0, trapFunction},
	}

	for _, c := range checks {
		got, ok := MethodSignatures[c.key]
		if !ok {
			t.Fatalf("missing MethodSignatures entry for %s", c.key)
		}
		if got.ParamSlots != c.slots {
			t.Fatalf("%s ParamSlots expected %d, got %d", c.key, c.slots, got.ParamSlots)
		}
		if got.GFunction == nil {
			t.Fatalf("%s GFunction expected non-nil", c.key)
		}
		// function identity check via function pointer
		if reflect.ValueOf(got.GFunction).Pointer() != reflect.ValueOf(c.fn).Pointer() {
			t.Fatalf("%s GFunction mismatch", c.key)
		}
	}
}

func TestConsoleClinit_NoClass_ReturnsErrBlk(t *testing.T) {
	globals.InitGlobals("test")
	_ = classloader.Init()
	classloader.LoadBaseClasses()

	ret := consoleClinit(nil)
	blk, ok := ret.(*GErrBlk)
	if !ok {
		t.Fatalf("expected *GErrBlk, got %T", ret)
	}
	if blk.ExceptionType != excNames.ClassNotLoadedException {
		t.Fatalf("unexpected exception type: %d", blk.ExceptionType)
	}
}

func TestConsolePrintf_WritesToSystemOut(t *testing.T) {
	globals.InitGlobals("test")

	// Pipe to capture writes to System.out
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	// Save and restore System.out static
	// We will restore to os.Stdout to keep environment sane for other tests.
	defer func() {
		_ = statics.AddStatic("java/lang/System.out", statics.Static{Type: "GS", Value: os.Stdout})
	}()
	_ = statics.AddStatic("java/lang/System.out", statics.Static{Type: "GS", Value: w})

	// Build params: [this, fmtString, argsArray]
	fmtObj := object.StringObjectFromGoString("Hello %s!")
	arg := object.StringObjectFromGoString("World")
	argsArr := makeObjectRefArray(arg)

	params := []interface{}{object.MakeEmptyObject(), fmtObj, argsArr}

	ret := consolePrintf(params)
	// Return is *os.File (stdout), but we care about content
	if _, ok := ret.(*os.File); !ok {
		t.Fatalf("consolePrintf expected to return *os.File, got %T", ret)
	}

 // Read what was written (Close signals EOF for the reader; Sync on pipes is not portable)
 _ = w.Close()
 buf, _ := io.ReadAll(r)
 if string(buf) != "Hello World!" {
     t.Fatalf("unexpected output: %q", string(buf))
 }
}

func TestConsoleReadLine_ReadsUntilNewline(t *testing.T) {
	globals.InitGlobals("test")

	// Pipe to feed input to System.in
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	defer r.Close()
	defer w.Close()

	defer func() {
		_ = statics.AddStatic("java/lang/System.in", statics.Static{Type: "GS", Value: os.Stdin})
	}()
	_ = statics.AddStatic("java/lang/System.in", statics.Static{Type: "GS", Value: r})

	// Write a line and close writer to simulate EOF after newline
	if _, err := w.Write([]byte("abc\n")); err != nil {
		t.Fatalf("write: %v", err)
	}
	_ = w.Close()

	ret := consoleReadLine(nil)
	sObj, ok := ret.(*object.Object)
	if !ok {
		t.Fatalf("expected *object.Object, got %T", ret)
	}
	got := object.GoStringFromStringObject(sObj)
	if got != "abc" {
		t.Fatalf("unexpected readLine content %q", got)
	}
}

func TestConsolePrintfReadLine_PromptAndRead(t *testing.T) {
	globals.InitGlobals("test")

	// Setup System.out
	rout, wout, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe out: %v", err)
	}
	defer rout.Close()
	defer wout.Close()
	defer func() {
		_ = statics.AddStatic("java/lang/System.out", statics.Static{Type: "GS", Value: os.Stdout})
	}()
	_ = statics.AddStatic("java/lang/System.out", statics.Static{Type: "GS", Value: wout})

	// Setup System.in
	rin, win, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe in: %v", err)
	}
	defer rin.Close()
	defer win.Close()
	defer func() {
		_ = statics.AddStatic("java/lang/System.in", statics.Static{Type: "GS", Value: os.Stdin})
	}()
	_ = statics.AddStatic("java/lang/System.in", statics.Static{Type: "GS", Value: rin})

	// Prepare params: [this, fmtString, argsArray(empty)]
	fmtObj := object.StringObjectFromGoString("Enter: ")
	emptyArgs := makeObjectRefArray() // zero-length array
	params := []interface{}{object.MakeEmptyObject(), fmtObj, emptyArgs}

	// Provide input line
	if _, err := win.Write([]byte("xyz\n")); err != nil {
		t.Fatalf("write: %v", err)
	}
	_ = win.Close()

	ret := consolePrintfReadLine(params)
	sObj := ret.(*object.Object)
	got := object.GoStringFromStringObject(sObj)
	if got != "xyz" {
		t.Fatalf("unexpected consolePrintfReadLine string: %q", got)
	}

	// Validate prompt was written
	_ = wout.Close()
	outBytes, _ := io.ReadAll(rout)
	if string(outBytes) != "Enter: " {
		t.Fatalf("unexpected prompt output: %q", string(outBytes))
	}
}
