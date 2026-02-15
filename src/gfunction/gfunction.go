/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-5 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/src/classloader"
	"jacobin/src/excNames"
	"jacobin/src/exceptions"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaIo"
	"jacobin/src/gfunction/javaLang"
	"jacobin/src/gfunction/javaMath"
	"jacobin/src/gfunction/javaNio"
	"jacobin/src/gfunction/javaSecurity"
	"jacobin/src/gfunction/javaText"
	"jacobin/src/gfunction/javaUtil"
	"jacobin/src/gfunction/javaxCrypto"
	"jacobin/src/gfunction/misc"
	"jacobin/src/gfunction/sunSecurity"
	"jacobin/src/trace"
	"strings"
)

// MTableLoadGFunctions loads the Go methods from files that contain them. It does this
// by calling the Load_* function in each of those files to load whatever Go functions
// they make available.
func MTableLoadGFunctions(MTable *classloader.MT) {

	// Load traps first, then override with our implementations for Files
	// so any methods we implement replace trapped entries.
	ghelpers.Load_Traps()
	ghelpers.Load_Traps_Java_Io()
	ghelpers.Load_Traps_Java_Nio()
	ghelpers.Load_Traps_Java_Security()

	// java/awt/*
	misc.Load_Awt_Graphics_Environment()

	// java/io/*
	javaIo.Load_Io_BufferedReader()
	javaIo.Load_Io_BufferedWriter()
	javaIo.Load_Io_Console()
	javaIo.Load_Io_File()
	javaIo.Load_Io_FileInputStream()
	javaIo.Load_Io_FileOutputStream()
	javaIo.Load_Io_FileReader()
	javaIo.Load_Io_FileWriter()
	javaIo.Load_Io_FilterInputStream()
	javaIo.Load_Io_InputStreamReader()
	javaIo.Load_Io_OutputStreamWriter()
	javaIo.Load_Io_PrintStream()
	javaIo.Load_Io_RandomAccessFile()

	// java/lang/*
	javaLang.ClassClinitIsh() // Special case clinit for java/lang/Class.
	javaLang.Load_Lang_Boolean()
	javaLang.Load_Lang_Byte()
	javaLang.Load_Lang_Character()
	javaLang.Load_Lang_CharSequence()
	javaLang.Load_Lang_Class()
	javaLang.Load_Lang_Double()
	javaLang.Load_Lang_Float()
	javaLang.Load_Lang_Integer()
	javaLang.Load_Lang_Long()
	javaLang.Load_Lang_Math()
	javaLang.Load_Lang_Object()
	javaLang.Load_Lang_Process()
	javaLang.Load_Lang_Process_Builder()
	javaLang.Load_Lang_Process_Handle_Impl()
	javaLang.Load_Lang_Reflect_Modifier()
	javaLang.Load_Lang_Runtime()
	javaLang.Load_Lang_SecurityManager()
	javaLang.Load_Lang_Short()
	javaLang.Load_Lang_StackTraceELement()
	javaLang.Load_Lang_String()
	javaLang.Load_Lang_StringBuffer()
	javaLang.Load_Lang_StringBuilder()
	javaLang.Load_Lang_System()
	javaLang.Load_Lang_Thread()
	javaLang.Load_Lang_Thread_Group()
	javaLang.Load_Lang_Thread_State()
	javaLang.Load_Lang_Throwable()
	javaLang.Load_Lang_UTF16()

	// java/math/*
	javaMath.Load_Math_Big_Decimal()
	javaMath.Load_Math_Big_Integer()
	javaMath.Load_Math_Math_Context()
	javaMath.Load_Math_Rounding_Mode()

	// java/nio/*
	javaNio.Load_Nio_File_Files()
	javaNio.Load_Nio_File_Path()
	javaNio.Load_Nio_File_Paths()

	// java/text/*
	javaText.Load_Math_SimpleDateFormat()

	// java/security/*
	javaSecurity.Load_ECFieldAndPoint()
	javaSecurity.Load_Security_Interfaces_EC_Keys()
	javaSecurity.Load_ECParameterSpec()
	javaSecurity.Load_EllipticCurve()
	javaSecurity.Load_Security_KeyPair()
	javaSecurity.Load_KeyPairGenerator()
	javaSecurity.Load_PublicAndPrivateKeys()
	javaSecurity.Load_Security()
	javaSecurity.Load_Security_Interfaces_DSA_Keys()
	javaSecurity.Load_Security_Interfaces_EC_Keys()
	javaSecurity.Load_Security_Interfaces_RSA_Keys()
	javaSecurity.Load_Security_Key()
	javaSecurity.Load_Security_MessageDigest()
	javaSecurity.Load_Security_Provider()
	javaSecurity.Load_Security_Provider_Service()
	javaSecurity.Load_Security_SecureRandom()
	javaSecurity.Load_Security_Signature()
	javaSecurity.Load_Security_Spec_NamedParameterSpec()
	javaSecurity.Load_Security_Spec_AlgorithmParameterSpec()

	// javax/crypto/*
	javaxCrypto.Load_Crypto_Interfaces_DH_Keys()
	javaxCrypto.Load_Crypto_KeyAgreement()
	javaxCrypto.Load_Crypto_Spec_SecretKeySpec()
	javaxCrypto.Load_Crypto_Spec_DHParameterSpec()

	// java/util/*
	javaUtil.Load_Util_ArrayList()
	javaUtil.Load_Util_Arrays()
	javaUtil.Load_Util_Base64()
	javaUtil.Load_Util_Concurrent_Atomic_AtomicInteger()
	javaUtil.Load_Util_Concurrent_Atomic_Atomic_Long()
	javaUtil.Load_Util_Date()
	javaUtil.Load_Util_Iterator()
	javaUtil.Load_Util_List()
	javaUtil.Load_Util_ListIterator()
	javaUtil.Load_Util_Hash_Map()
	javaUtil.Load_Util_Hash_Set()
	javaUtil.Load_Util_HexFormat()
	javaUtil.Load_Util_LinkedList()
	javaUtil.Load_Util_Locale()
	javaUtil.Load_Util_Logging_Logger()
	javaUtil.Load_Util_Map()
	javaUtil.Load_Util_Properties()
	javaUtil.Load_Util_Objects()
	javaUtil.Load_Util_Optional()
	javaUtil.Load_Util_Random()
	javaUtil.Load_Util_TimeZone()
	javaUtil.Load_Util_Vector()
	javaUtil.Load_Util_Zip_Adler32()
	javaUtil.Load_Util_Zip_Crc32_Crc32c()

	// javax.*
	misc.Load_Javax_Net_Ssl_SSLContext()

	// jdk/internal/misc/*
	misc.Load_Jdk_Internal_Misc_Unsafe()
	misc.Load_Jdk_Internal_Misc_ScopedMemoryAccess()

	// Sun
	sunSecurity.Load_Sun_Security_Action_GetPropertyAction()
	sunSecurity.Load_Sun_Security_Jca_ProviderList()

	// Load functions that invoke ghelpers.ClinitGeneric() and do nothing else.
	Load_Other_Methods()

	// Load diagnostic helper functions.
	misc.Load_jj()

	//	now, with the accumulated ghelpers.MethodSignatures maps, load MTable.
	loadlib(MTable, ghelpers.MethodSignatures)
	ghelpers.TestGfunctionsLoaded = true
}

// load the test gfunctions in testGfunctions.go
func LoadTestGfunctions(MTable *classloader.MT) {
	Load_TestGfunctions()
	loadlib(MTable, ghelpers.TestMethodSignatures)
	ghelpers.TestGfunctionsLoaded = true
}

func checkKey(key string) bool {
	if strings.Index(key, ".") == -1 || strings.Index(key, "(") == -1 || strings.Index(key, ")") == -1 {
		return false
	}
	if strings.HasSuffix(key, ")") {
		return false
	}
	return true
}

func loadlib(tbl *classloader.MT, libMeths map[string]ghelpers.GMeth) {
	ok := true
	for key, val := range libMeths {
		if !checkKey(key) {
			errMsg := fmt.Sprintf("loadlib: Invalid key=%s", key)
			trace.Error(errMsg)
			ok = false
		}
		gme := ghelpers.GMeth{}
		gme.ParamSlots = val.ParamSlots
		gme.GFunction = val.GFunction
		gme.NeedsContext = val.NeedsContext

		tableEntry := classloader.MTentry{
			MType: 'G',
			Meth:  gme,
		}

		classloader.AddEntry(tbl, key, tableEntry)
	}
	if !ok {
		exceptions.ThrowExNil(excNames.InternalException, "loadlib: at least one key was invalid")
	}
}
