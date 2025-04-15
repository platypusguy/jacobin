/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/classloader"
	"jacobin/excNames"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/shutdown"
	"jacobin/statics"
	"jacobin/stringPool"
	"jacobin/trace"
	"jacobin/types"
	"os"
	"runtime"
	"strings"
	"time"
)

/*
 Each object or library that has golang methods contains a reference to MethodSignatures,
 which contain data needed to insert the golang method into the MTable of the currently
 executing JVM. MethodSignatures is a map whose key is the fully qualified name and
 type of the method (that is, the method's full signature) and a value consisting of
 a struct of an int (the number of slots to pop off the caller's operand stack when
 creating the new frame and a function). All methods have the same signature, regardless
 of the signature of their Java counterparts. That signature is that it accepts a slice
 of interface{} and returns an interface{}. The accepted slice can be empty and the
 return interface can be nil. This covers all Java functions. (Objects are returned
 as a 64-bit address in this scheme--as they are in the JVM).

 The passed-in slice contains one entry for every parameter passed to the method (which
 could mean an empty slice). Note: longs and doubles use only one parameter entry each.
*/

func Load_Lang_System() {

	MethodSignatures["java/lang/System.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  systemClinit,
		}

	MethodSignatures["java/lang/System.allowSecurityManager()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  returnFalse,
		}

	MethodSignatures["java/lang/System.arraycopy(Ljava/lang/Object;ILjava/lang/Object;II)V"] = // copy array (full or partial)
		GMeth{
			ParamSlots: 5,
			GFunction:  systemArrayCopy,
		}

	MethodSignatures["java/lang/System.clearProperty(Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  systemClearProperty,
		}

	MethodSignatures["java/lang/System.console()Ljava/io/Console;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  systemConsole,
		}

	MethodSignatures["java/lang/System.currentTimeMillis()J"] = // get time in ms since Jan 1, 1970, returned as long
		GMeth{
			ParamSlots: 0,
			GFunction:  systemCurrentTimeMillis,
		}

	MethodSignatures["java/lang/System.exit(I)V"] = // shutdown the app
		GMeth{
			ParamSlots: 1,
			GFunction:  systemExitI,
		}

	MethodSignatures["java/lang/System.gc()V"] = // for a GC cycle
		GMeth{
			ParamSlots: 0,
			GFunction:  systemForceGC,
		}

	MethodSignatures["java/lang/System.getenv()Ljava/util/Map;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/System.getenv(Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  systemGetenv,
		}

	MethodSignatures["java/lang/System.getLogger(Ljava/lang/String;)Ljava/lang/System/Logger;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/System.getLogger(Ljava/lang/String;Ljava/util/ResourceBundle;)Ljava/lang/System/Logger;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/System.getProperties()Ljava/util/Properties;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  systemGetProperties,
		}

	MethodSignatures["java/lang/System.getProperty(Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  systemGetProperty,
		}

	MethodSignatures["java/lang/System.getProperty(Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  systemGetPropertyDefault,
		}

	MethodSignatures["java/lang/System.getSecurityManager()Ljava/lang/SecurityManager;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/lang/System.identityHashCode(Ljava/lang/Object;)I"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/System.inheritedChannel()Ljava/nio/channels/Channel;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/System.lineSeparator()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  systemLineSeparator,
		}

	MethodSignatures["java/lang/System.load(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/System.loadLibrary(Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/System.mapLibraryName(Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/System.nanoTime()J"] = // get nanoseconds time, returned as long
		GMeth{
			ParamSlots: 0,
			GFunction:  systemNanoTime,
		}

	MethodSignatures["java/lang/System.registerNatives()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  justReturn,
		}

	MethodSignatures["java/lang/System.runFinalization()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapDeprecated,
		}

	MethodSignatures["java/lang/System.setErr(Ljava/io/PrintStream;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/System.setIn(Ljava/io/InputStream;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/System.setOut(Ljava/io/PrintStream;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/lang/System.setProperties(Ljava/util/Properties;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  systemSetProperties,
		}

	MethodSignatures["java/lang/System.setProperty(Ljava/lang/String;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  systemSetProperty,
		}

	MethodSignatures["java/lang/System.setSecurityManager(Ljava/lang/SecurityManager;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapDeprecated,
		}

}

func systemClinit([]interface{}) interface{} {
	klass := classloader.MethAreaFetch("java/lang/System")
	if klass == nil {
		errMsg := "systemClinit: Expected java/lang/System to be in the MethodArea, but it was not"
		trace.Error(errMsg)
		return getGErrBlk(excNames.ClassNotLoadedException, errMsg)
	}
	if klass.Data.ClInit != types.ClInitRun {
		_ = statics.AddStatic("java/lang/System.in", statics.Static{Type: "GS", Value: os.Stdin})
		_ = statics.AddStatic("java/lang/System.err", statics.Static{Type: "GS", Value: os.Stderr})
		_ = statics.AddStatic("java/lang/System.out", statics.Static{Type: "GS", Value: os.Stdout})
		klass.Data.ClInit = types.ClInitRun
	}
	return nil
}

// systemArrayCopy copies an array or subarray from one array to another, both of which must exist.
// It is a complex native function in the JDK. Javadoc here:
// docs.oracle.com/en/java/javase/17/docs/api/java.base/java/lang/System.html#arraycopy(java.lang.Object,int,java.lang.Object,int,int)
func systemArrayCopy(params []interface{}) interface{} {
	if len(params) != 5 {
		errMsg := fmt.Sprintf("systemArrayCopy: Expected 5 parameters, got %d", len(params))
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	src := params[0].(*object.Object)
	srcPos := params[1].(int64)
	dest := params[2].(*object.Object)
	destPos := params[3].(int64)
	length := params[4].(int64)

	if object.IsNull(src) || object.IsNull(dest) {
		errMsg := fmt.Sprintf("systemArrayCopy: null src or dest")
		return getGErrBlk(excNames.NullPointerException, errMsg)
	}

	if srcPos < 0 || destPos < 0 || length < 0 {
		errMsg := fmt.Sprintf(
			"systemArrayCopy: Negative position in: srcPose=%d, destPos=%d, or length=%d", srcPos, destPos, length)
		return getGErrBlk(excNames.ArrayIndexOutOfBoundsException, errMsg)
	}

	srcType := *(stringPool.GetStringPointer(src.KlassName))
	destType := *(stringPool.GetStringPointer(dest.KlassName))

	if !strings.HasPrefix(srcType, types.Array) || !strings.HasPrefix(destType, types.Array) || srcType != destType {
		errMsg := fmt.Sprintf("systemArrayCopy: invalid src or dest array")
		return getGErrBlk(excNames.ArrayStoreException, errMsg)
	}

	srcLen := object.ArrayLength(src)
	destLen := object.ArrayLength(dest)

	if srcPos+length > srcLen || destPos+length > destLen {
		errMsg := fmt.Sprintf("systemArrayCopy: Array position + length exceeds array size")
		return getGErrBlk(excNames.ArrayIndexOutOfBoundsException, errMsg)
	}

	s := srcPos
	d := destPos

	if (src != dest) || ((src == dest) && (srcPos+length < destPos)) {
		// non-overlapping copy of identical items
		switch srcType {
		case types.ByteArray:
			switch src.FieldTable["value"].Fvalue.(type) {
			case []types.JavaByte:
				sArr := src.FieldTable["value"].Fvalue.([]types.JavaByte)
				dArr := dest.FieldTable["value"].Fvalue.([]types.JavaByte)
				for i := int64(0); i < length; i++ {
					dArr[d] = sArr[s]
					d += 1
					s += 1
				}
			case []byte:
				sArr := src.FieldTable["value"].Fvalue.([]byte)
				dArr := dest.FieldTable["value"].Fvalue.([]byte)
				for i := int64(0); i < length; i++ {
					dArr[d] = sArr[s]
					d += 1
					s += 1
				}
			}

		case types.RefArray: // TODO: make sure refs are to the same object types
			sArr := src.FieldTable["value"].Fvalue.([]*object.Object)
			dArr := dest.FieldTable["value"].Fvalue.([]*object.Object)
			for i := int64(0); i < length; i++ {
				dArr[d] = sArr[s]
				d += 1
				s += 1
			}
		case types.FloatArray:
			sArr := src.FieldTable["value"].Fvalue.([]float64)
			dArr := dest.FieldTable["value"].Fvalue.([]float64)
			for i := int64(0); i < length; i++ {
				dArr[d] = sArr[s]
				d += 1
				s += 1
			}
		case types.IntArray:
			sArr := src.FieldTable["value"].Fvalue.([]int64)
			dArr := dest.FieldTable["value"].Fvalue.([]int64)
			for i := int64(0); i < length; i++ {
				dArr[d] = sArr[s]
				d += 1
				s += 1
			}
		}
	} else { // overlapping copy uses a temporary array
		tempArray := make([]interface{}, length)

		switch srcType {
		case types.ByteArray:
			sArr := src.FieldTable["value"].Fvalue.([]types.JavaByte)
			dArr := dest.FieldTable["value"].Fvalue.([]types.JavaByte)
			for i := int64(0); i < length; i++ {
				tempArray[i] = sArr[s]
				s += 1
			}
			for i := int64(0); i < length; i++ {
				dArr[d] = tempArray[i].(types.JavaByte)
				d += 1
			}

		case types.RefArray: // TODO: make sure refs are to the same object types
			sArr := src.FieldTable["value"].Fvalue.([]*object.Object)
			dArr := dest.FieldTable["value"].Fvalue.([]*object.Object)
			for i := int64(0); i < length; i++ {
				tempArray[i] = sArr[s]
				s += 1
			}
			for i := int64(0); i < length; i++ {
				dArr[d] = tempArray[i].(*object.Object)
				d += 1
			}

		case types.FloatArray:
			sArr := src.FieldTable["value"].Fvalue.([]float64)
			dArr := dest.FieldTable["value"].Fvalue.([]float64)
			for i := int64(0); i < length; i++ {
				tempArray[i] = sArr[s]
				s += 1
			}
			for i := int64(0); i < length; i++ {
				dArr[d] = tempArray[i].(float64)
				d += 1
			}

		case types.IntArray:
			sArr := src.FieldTable["value"].Fvalue.([]int64)
			dArr := dest.FieldTable["value"].Fvalue.([]int64)
			for i := int64(0); i < length; i++ {
				tempArray[i] = sArr[s]
				s += 1
			}
			for i := int64(0); i < length; i++ {
				dArr[d] = tempArray[i].(int64)
				d += 1
			}

		}
	}

	return nil
}

// Return the system input console as a *os.File.
func systemConsole([]interface{}) interface{} {
	return statics.GetStaticValue("java/lang/System", "in")
}

// Return time in milliseconds, measured since midnight of Jan 1, 1970
func systemCurrentTimeMillis([]interface{}) interface{} {
	return time.Now().UnixMilli() // is int64
}

// Return time in nanoseconds. Note that in golang this function has a lower (that is, less good)
// resolution than Java: two successive calls often return the same value.
func systemNanoTime([]interface{}) interface{} {
	return time.Now().UnixNano() // is int64
}

// Exits the program directly, returning the passed in value
// exit is a static function, so no object ref and exit value is in params[0]
func systemExitI(params []interface{}) interface{} {
	exitCode := params[0].(int64)
	var exitStatus = int(exitCode)
	shutdown.Exit(exitStatus)
	return exitCode // this code is not executed as previous line ends Jacobin
}

// Force a garbage collection cycle.
func systemForceGC([]interface{}) interface{} {
	runtime.GC()
	return nil
}

// Get an environment variable string.
func systemGetenv(params []interface{}) interface{} {
	key := object.GoStringFromStringObject(params[0].(*object.Object))
	return object.StringObjectFromGoString(os.Getenv(key))
}

// Get a system property - high level function.
func systemClearProperty(params []interface{}) interface{} {
	propObj := params[0].(*object.Object) // string
	propStr := object.GoStringFromStringObject(propObj)

	value := globals.GetSystemProperty(propStr)
	if value == "" {
		return object.Null
	}

	globals.RemoveSystemProperty(propStr)

	return object.StringObjectFromGoString(value)
}

// Get a system property - high level function.
func systemGetProperty(params []interface{}) interface{} {
	propObj := params[0].(*object.Object)
	propStr := object.GoStringFromStringObject(propObj)

	value := globals.GetSystemProperty(propStr)
	if value == "" {
		return object.Null
	}
	return object.StringObjectFromGoString(value)
}

// Get a system property - high level function.
func systemGetPropertyDefault(params []interface{}) interface{} {
	propObj := params[0].(*object.Object)
	propStr := object.GoStringFromStringObject(propObj)

	value := globals.GetSystemProperty(propStr)
	if value == "" {
		return params[1].(*object.Object)
	}
	return object.StringObjectFromGoString(value)
}

// systemGetProperties: Create a Properties object and set its map elements to system properties.
func systemGetProperties([]interface{}) interface{} {
	var propMap types.DefProperties
	propMap = make(types.DefProperties)

	propMap["file.encoding"] = globals.GetSystemProperty("file.encoding")
	propMap["file.separator"] = globals.GetSystemProperty("file.separator")
	propMap["java.compiler"] = globals.GetSystemProperty("java.compiler")
	propMap["java.home"] = globals.GetSystemProperty("java.home")
	propMap["java.io.tmpdir"] = globals.GetSystemProperty("java.io.tmpdir")
	propMap["java.library.path"] = globals.GetSystemProperty("java.library.path")
	propMap["java.vendor"] = globals.GetSystemProperty("java.vendor")
	propMap["java.vendor.url"] = globals.GetSystemProperty("java.vendor.url")
	propMap["java.vendor.version"] = globals.GetSystemProperty("java.vendor.version")
	propMap["java.version"] = globals.GetSystemProperty("java.version")
	propMap["java.vm.name"] = globals.GetSystemProperty("java.vm.name")
	propMap["java.vm.specification.name"] = globals.GetSystemProperty("java.vm.specification.name")
	propMap["java.vm.specification.vendor"] = globals.GetSystemProperty("java.vm.specification.vendor")
	propMap["java.vm.specification.version"] = globals.GetSystemProperty("java.vm.specification.version")
	propMap["java.vm.vendor"] = globals.GetSystemProperty("java.vm.vendor")
	propMap["java.vm.version"] = globals.GetSystemProperty("java.vm.version")
	propMap["line.separator"] = globals.GetSystemProperty("line.separator")
	propMap["native.encoding"] = globals.GetSystemProperty("native.encoding")
	propMap["os.arch"] = globals.GetSystemProperty("os.arch")
	propMap["os.name"] = globals.GetSystemProperty("os.name")
	propMap["os.version"] = globals.GetSystemProperty("os.version")
	propMap["path.separator"] = globals.GetSystemProperty("path.separator")
	propMap["stdout.encoding"] = globals.GetSystemProperty("stdout.encoding")
	propMap["stderr.encoding"] = globals.GetSystemProperty("stderr.encoding")
	propMap["user.dir"] = globals.GetSystemProperty("user.dir")
	propMap["user.home"] = globals.GetSystemProperty("user.home")
	propMap["user.name"] = globals.GetSystemProperty("user.name")
	propMap["user.timezone"] = globals.GetSystemProperty("user.timezone")

	return object.MakeOneFieldObject(classNameProperties, fieldNameProperties, types.Properties, propMap)

}

// Get the system line separator.
func systemLineSeparator([]interface{}) interface{} {
	str := globals.GetSystemProperty("line.separator")
	return object.StringObjectFromGoString(str)
}

// Set a system property.
func systemSetProperties(params []interface{}) interface{} {
	propertiesObj := params[0].(*object.Object)
	newMap := propertiesObj.FieldTable[fieldNameProperties].Fvalue.(types.DefProperties)

	globals.ReplaceSystemProperties(newMap)

	return nil
}

// Set a system property.
func systemSetProperty(params []interface{}) interface{} {
	keyObj := params[0].(*object.Object)
	keyStr := object.GoStringFromStringObject(keyObj)
	valueObj := params[1].(*object.Object)
	valueStr := object.GoStringFromStringObject(valueObj)

	value := globals.GetSystemProperty(keyStr)
	globals.SetSystemProperty(keyStr, valueStr)

	return object.StringObjectFromGoString(value)
}
