/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaIo

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/util"
	"strings"
	"testing"
)

func loadAllGFunctions() {
	globals.InitGlobals("test")
	Load_Io_BufferedInputStream()
	Load_Io_BufferedOutputStream()
	Load_Io_BufferedReader()
	Load_Io_BufferedWriter()
	Load_Io_ByteArrayInputStream()
	Load_Io_ByteArrayOutputStream()
	Load_Io_Console()
	Load_Io_File()
	Load_Io_FileInputStream()
	Load_Io_FileOutputStream()
	Load_Io_FileReader()
	Load_Io_FileWriter()
	Load_Io_FilterInputStream()
	Load_Io_FilterOutputStream()
	Load_Io_InputStream()
	Load_Io_InputStreamReader()
	Load_Io_OutputStream()
	Load_Io_OutputStreamWriter()
	Load_Io_PrintStream()
	Load_Io_RandomAccessFile()
}

const classPrefix = "java/io"

func TestJavaIoParamSlotsConsistency(t *testing.T) {
	loadAllGFunctions()

	var mismatches []string
	for k, v := range ghelpers.MethodSignatures {
		if !strings.HasPrefix(k, classPrefix) {
			continue
		}
		_, _, methType := util.ParseFQN(k)
		params := util.ParseIncomingParamsFromMethTypeString(methType)
		expectedSlots := len(params)
		if v.ParamSlots != expectedSlots {
			t.Errorf("Signature mismatch for %s: ParamSlots is %d, expected %d (params: %v)", k, v.ParamSlots, expectedSlots, params)
			mismatches = append(mismatches, k)
		}
	}
	t.Logf("Total %s signatures checked: %d, mismatches: %d", classPrefix, len(ghelpers.MethodSignatures), len(mismatches))
	if len(mismatches) > 0 {
		t.Fatalf("Found %d ParamSlots mismatches", len(mismatches))
	}
}
