/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/util"
	"strings"
	"testing"
)

func loadAllJavaUtil() {
	globals.InitGlobals("test")
	Load_Util_ArrayList()
	Load_Util_Arrays()
	Load_Util_Base64()
	Load_Util_BitSet()
	Load_Util_Collection()
	Load_Util_Concurrent_Atomic_AtomicInteger()
	Load_Util_Concurrent_Atomic_Atomic_Long()
	Load_Util_Concurrent_CyclicBarrier()
	Load_Util_Concurrent_Hash_Map()
	Load_Util_Concurrent_TimeUnit()
	Load_Util_Date()
	Load_Util_Function_Supplier()
	Load_Util_Iterator()
	Load_Util_List()
	Load_Util_ListIterator()
	Load_Util_Set()
	Load_Util_StringTokenizer()
	Load_Util_SequencedCollection()
	Load_Util_SequencedMap()
	Load_Util_SequencedSet()
	Load_Util_Hash_Map()
	Load_Util_Hash_Set()
	Load_Util_Tree_Set()
	Load_Util_HexFormat()
	Load_Util_LinkedList()
	Load_Util_Locale()
	Load_Util_LongSummaryStatistics()
	Load_Util_Logging_ConsoleHandler()
	Load_Util_Logging_FileHandler()
	Load_Util_Logging_Filter()
	Load_Util_Logging_Formatter()
	Load_Util_Logging_Handler()
	Load_Util_Logging_Level()
	Load_Util_Logging_Logger()
	Load_Util_Logging_LogManager()
	Load_Util_Logging_LogRecord()
	Load_Util_Logging_MemoryHandler()
	Load_Util_Logging_SimpleFormatter()
	Load_Util_Logging_SocketHandler()
	Load_Util_Logging_StreamHandler()
	Load_Util_Map()
	Load_Util_Properties()
	Load_Util_Objects()
	Load_Util_Optional()
	Load_Util_Random()
	Load_Util_ResourceBundle()
	Load_Util_TimeZone()
	Load_Util_Vector()
	Load_Util_Zip_Adler32()
	Load_Util_Zip_CheckedInputStream()
	Load_Util_Zip_Crc32_Crc32c()
	Load_Util_Zip_ZipEntry()
	Load_Util_Zip_ZipInputStream()
}

func TestJavaUtilParamSlotsConsistency(t *testing.T) {
	// Save and restore the global ghelpers.MethodSignatures map to avoid test pollution
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	loadAllJavaUtil()

	var mismatches []string
	for k, v := range ghelpers.MethodSignatures {
		if !strings.HasPrefix(k, "java/util/") {
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
	t.Logf("Total java/util/ signatures checked: %d, mismatches: %d", len(ghelpers.MethodSignatures), len(mismatches))
	if len(mismatches) > 0 {
		t.Fatalf("Found %d ParamSlots mismatches", len(mismatches))
	}
}
