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
	"jacobin/exceptions"
	"jacobin/object"
	"jacobin/statics"
	"jacobin/types"
	"os"
	"syscall"
)

// Implementation of some of the functions in Java/lang/Class.

func Load_Io_Console() map[string]GMeth {

	// Class initialisation for Console.
	MethodSignatures["java/io/Console.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  consoleClinit,
		}

	// Flushes the console and forces any buffered output to be written immediately.
	MethodSignatures["java/io/Console.charset()Ljava/nio/charset/Charset;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  noSupportYetInConsole,
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
			GFunction:  noSupportYetInConsole,
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
			GFunction:  noSupportYetInConsole,
		}

	return MethodSignatures
}

// Initialise class Console.
func consoleClinit([]interface{}) interface{} {
	klass := classloader.MethAreaFetch("java/io/Console")
	if klass == nil {
		errMsg := "consoleClinit: Could not find java/io/Console in the MethodArea"
		return getGErrBlk(exceptions.ClassNotLoadedException, errMsg)
	}
	klass.Data.ClInit = types.ClInitRun // just mark that String.<clinit>() has been run
	return nil
}

// No support YET for references to Charset objects nor for Unicode code point arrays
func noSupportYetInConsole([]interface{}) interface{} {
	errMsg := "No support yet for Reader/PrintWriter/Charset in class Console"
	return getGErrBlk(exceptions.UnsupportedOperationException, errMsg)
}

// Flush java/lang/System.in/out/err.
func consoleFlush([]interface{}) interface{} {
	stdinout := statics.GetStaticValue("java/lang/System", "in").(*os.File)
	_ = stdinout.Sync()
	stdinout = statics.GetStaticValue("java/lang/System", "out").(*os.File)
	_ = stdinout.Sync()
	// Note: java/lang/System.err is not associated with the system console.
	return nil
}

// Printf -- handle the variable args and then call golang's own printf function
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
	str := object.GetGoStringFromObject(objPtr)
	stdout := statics.GetStaticValue("java/lang/System", "out").(*os.File)
	_, _ = fmt.Fprint(stdout, str)
	return stdout // Return the *os.File

}

// Reads a single line of text from the console.
func consoleReadLine([]interface{}) interface{} {
	var bytes []byte
	var bite = []byte{0x00}
	var nbytes int
	var err error
	stdin := statics.GetStaticValue("java/lang/System", "in").(*os.File)
	for {
		nbytes, err = stdin.Read(bite)
		if nbytes == 0 {
			break
		}
		if err != nil {
			errMsg := fmt.Sprintf("consoleReadLine stdin.Read: %s", err.Error())
			return getGErrBlk(exceptions.IOException, errMsg)
		}
		if bite[0] == '\n' {
			break
		}
		bytes = append(bytes, bite[0])
	}
	str := string(bytes)
	return object.NewPoolStringFromGoString(str)
}

// Provides a formatted prompt, then reads a single line of text from the console.
func consolePrintfReadLine(params []interface{}) interface{} {
	_ = consolePrintf(params)
	objPtr := consoleReadLine(params)
	return objPtr
}

func consoleReadPassword([]interface{}) interface{} {
	password, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		errMsg := fmt.Sprintf("consoleReadPassword term.ReadPassword: %s", err.Error())
		return getGErrBlk(exceptions.IOException, errMsg)
	}
	stdout := statics.GetStaticValue("java/lang/System", "out").(*os.File)
	_, _ = fmt.Fprint(stdout, "\n")
	return password
}

// Provides a formatted prompt, then reads a password or passphrase from the console.
func consolePrintfReadPassword(params []interface{}) interface{} {
	_ = consolePrintf(params)
	return consoleReadPassword(params)
}
