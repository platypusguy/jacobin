/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"golang.org/x/term"
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/statics"
	"jacobin/types"
	"os"
	"syscall"
)

// Implementation of some of the functions in Java/lang/Class.

func Load_Io_Console() {

	// Class initialisation for Console.
	MethodSignatures["java/io/Console.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  consoleClinit,
		}

	// Returns the Charset object used for the Console.
	MethodSignatures["java/io/Console.charset()Ljava/nio/charset/Charset;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	// Flush java/lang/System.in/out/err.
	MethodSignatures["java/io/Console.flush()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  consoleFlush,
		}

	// Console format.
	MethodSignatures["java/io/Console.format(Ljava/lang/String;[Ljava/lang/Object;)Ljava/io/Console;"] =
		GMeth{
			ParamSlots: 2, // the format string, the parameters (if any)
			GFunction:  consolePrintf,
		}

	// Console Printf.
	MethodSignatures["java/io/Console.printf(Ljava/lang/String;[Ljava/lang/Object;)Ljava/io/Console;"] =
		GMeth{
			ParamSlots: 2, // the format string, the parameters (if any)
			GFunction:  consolePrintf,
		}

	// Retrieves the unique Reader object associated with this console.
	MethodSignatures["java/io/Console.reader()Ljava/io/Reader;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	// Reads a single line of text from the console.
	MethodSignatures["java/io/Console.readLine()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  consoleReadLine,
		}

	// Provides a formatted prompt, then reads a single line of text from the console.
	MethodSignatures["java/io/Console.readLine(Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  consolePrintfReadLine,
		}

	// Reads a password or passphrase from the console with echoing disabled.
	MethodSignatures["java/io/Console.readPassword()[C"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  consoleReadPassword,
		}

	// Provides a formatted prompt, then reads a password or passphrase from the console.
	MethodSignatures["java/io/Console.readPassword(Ljava/lang/String;[Ljava/lang/Object;)[C"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  consolePrintfReadPassword,
		}

	// Retrieves the unique PrintWriter object associated with this console.
	MethodSignatures["java/io/Console.writer()Ljava/io/PrintWriter;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

}

// "java/io/Console.<clinit>()V" - Initialise class Console.
func consoleClinit([]interface{}) interface{} {
	klass := classloader.MethAreaFetch("java/io/Console")
	if klass == nil {
		errMsg := "Console <clinit>: Could not find java/io/Console in the MethodArea"
		return getGErrBlk(excNames.ClassNotLoadedException, errMsg)
	}
	klass.Data.ClInit = types.ClInitRun // just mark that String.<clinit>() has been run
	return nil
}

// Flush java/lang/System.in/out/err.
// "java/io/Console.flush()V"
func consoleFlush([]interface{}) interface{} {
	stdinout := statics.GetStaticValue("java/lang/System", "in").(*os.File)
	_ = stdinout.Sync()
	stdinout = statics.GetStaticValue("java/lang/System", "out").(*os.File)
	_ = stdinout.Sync()
	// Note: java/lang/System.err is not associated with the system console.
	return nil
}

// Printf -- handle the variable args and then call golang's own printf function
// "java/io/Console.format(Ljava/lang/String;[Ljava/lang/Object;)Ljava/io/Console;"
func consolePrintf(params []interface{}) interface{} {
	var intfSprintf = new([]interface{})
	*intfSprintf = append(*intfSprintf, params[1])
	*intfSprintf = append(*intfSprintf, params[2])
	retval := StringFormatter(*intfSprintf)
	switch retval.(type) {
	case *object.Object:
	default:
		return retval
	}
	objPtr := retval.(*object.Object)
	str := object.GoStringFromStringObject(objPtr)
	stdout := statics.GetStaticValue("java/lang/System", "out").(*os.File)
	_, _ = fmt.Fprint(stdout, str)
	return stdout // Return the *os.File

}

// Reads a single line of text from the console.
// "java/io/Console.readLine()Ljava/lang/String;"
func consoleReadLine([]interface{}) interface{} {
	var bytes []byte
	var bb = []byte{0x00}
	var nbytes int
	var err error
	stdin := statics.GetStaticValue("java/lang/System", "in").(*os.File)
	for {
		nbytes, err = stdin.Read(bb)
		if nbytes == 0 {
			break
		}
		if err != nil {
			errMsg := fmt.Sprintf("stdin.Read error: %s", err.Error())
			return getGErrBlk(excNames.IOException, errMsg)
		}
		if bb[0] == '\n' {
			break
		}
		bytes = append(bytes, bb[0])
	}
	str := string(bytes)
	return object.StringObjectFromGoString(str)
}

// Provides a formatted prompt, then reads a single line of text from the console.
// "java/io/Console.readLine(Ljava/lang/String;[Ljava/lang/Object;)Ljava/lang/String;"
func consolePrintfReadLine(params []interface{}) interface{} {
	_ = consolePrintf(params)
	objPtr := consoleReadLine(params)
	return objPtr
}

// Read a password from console.
// "java/io/Console.readPassword()[C"
func consoleReadPassword([]interface{}) interface{} {
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		errMsg := fmt.Sprintf("stdin.ReadPassword failed, reason: %s", err.Error())
		return getGErrBlk(excNames.IOException, errMsg)
	}
	stdout := statics.GetStaticValue("java/lang/System", "out").(*os.File)
	_, _ = fmt.Fprint(stdout, "\n")

	// Convert password to int64 array, insert into an object, and return to caller
	var iArray []int64
	for _, bb := range password {
		iArray = append(iArray, int64(bb))
	}
	return populator("[C", types.IntArray, iArray)
}

// Provides a formatted prompt, then reads a password or passphrase from the console.
// "java/io/Console.readPassword(Ljava/lang/String;[Ljava/lang/Object;)[C"
func consolePrintfReadPassword(params []interface{}) interface{} {
	_ = consolePrintf(params)
	return consoleReadPassword(params)
}
