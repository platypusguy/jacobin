package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"reflect"
	"testing"
)

func newLogRecordObj() *object.Object {
	return object.MakeEmptyObjectWithClassName(&logRecordClassName)
}

func TestLoggingLogRecordInitAndGetters(t *testing.T) {
	globals.InitStringPool()

	rec := newLogRecordObj()
	level := makeLevelObject("WARNING", standardLevels["WARNING"], "")
	msg := object.StringObjectFromGoString("hello world")

	ret := loggingLogRecordInit([]interface{}{rec, level, msg})
	if ret != nil {
		t.Fatalf("loggingLogRecordInit returned error: %v", ret)
	}

	if got := loggingLogRecordGetLevel([]interface{}{rec}); got != level {
		t.Fatalf("expected getLevel to return the initial level")
	}
	if got := loggingLogRecordGetMessage([]interface{}{rec}); got != msg {
		t.Fatalf("expected getMessage to return the initial message")
	}
	if got := loggingLogRecordGetLoggerName([]interface{}{rec}); got != object.Null {
		t.Fatalf("expected getLoggerName to default to null")
	}
	if got := loggingLogRecordGetSequenceNumber([]interface{}{rec}).(int64); got != 0 {
		t.Fatalf("expected default sequence number 0, got %d", got)
	}
	if got := loggingLogRecordGetSourceClassName([]interface{}{rec}); got != object.Null {
		t.Fatalf("expected getSourceClassName to default to null")
	}
	if got := loggingLogRecordGetSourceMethodName([]interface{}{rec}); got != object.Null {
		t.Fatalf("expected getSourceMethodName to default to null")
	}
	if got := loggingLogRecordGetParameters([]interface{}{rec}); got != object.Null {
		t.Fatalf("expected getParameters to default to null")
	}
	if got := loggingLogRecordGetThreadID([]interface{}{rec}).(int64); got != 0 {
		t.Fatalf("expected default threadID 0, got %d", got)
	}
	if got := loggingLogRecordGetThrown([]interface{}{rec}); got != object.Null {
		t.Fatalf("expected getThrown to default to null")
	}
	if got := loggingLogRecordGetResourceBundleName([]interface{}{rec}); got != object.Null {
		t.Fatalf("expected getResourceBundleName to default to null")
	}
	if got := loggingLogRecordGetResourceBundle([]interface{}{rec}); got != object.Null {
		t.Fatalf("expected getResourceBundle to default to null")
	}
}

func TestLoggingLogRecordSettersAndGetters(t *testing.T) {
	globals.InitStringPool()

	rec := newLogRecordObj()
	level := makeLevelObject("INFO", standardLevels["INFO"], "")
	msg := object.StringObjectFromGoString("m1")
	_ = loggingLogRecordInit([]interface{}{rec, level, msg})

	// setLevel/getLevel
	newLevel := makeLevelObject("SEVERE", standardLevels["SEVERE"], "")
	_ = loggingLogRecordSetLevel([]interface{}{rec, newLevel})
	if got := loggingLogRecordGetLevel([]interface{}{rec}); got != newLevel {
		t.Fatalf("expected getLevel to return the set level")
	}

	// setMessage/getMessage
	newMsg := object.StringObjectFromGoString("m2")
	_ = loggingLogRecordSetMessage([]interface{}{rec, newMsg})
	if got := loggingLogRecordGetMessage([]interface{}{rec}); got != newMsg {
		t.Fatalf("expected getMessage to return the set message")
	}

	// setLoggerName/getLoggerName
	loggerName := object.StringObjectFromGoString("com.foo.Bar")
	_ = loggingLogRecordSetLoggerName([]interface{}{rec, loggerName})
	if got := loggingLogRecordGetLoggerName([]interface{}{rec}); got != loggerName {
		t.Fatalf("expected getLoggerName to return the set name")
	}

	// setSequenceNumber/getSequenceNumber
	_ = loggingLogRecordSetSequenceNumber([]interface{}{rec, int64(42)})
	if got := loggingLogRecordGetSequenceNumber([]interface{}{rec}).(int64); got != 42 {
		t.Fatalf("expected sequence number 42, got %d", got)
	}

	// setSourceClassName/getSourceClassName
	className := object.StringObjectFromGoString("com.foo.Bar")
	_ = loggingLogRecordSetSourceClassName([]interface{}{rec, className})
	if got := loggingLogRecordGetSourceClassName([]interface{}{rec}); got != className {
		t.Fatalf("expected getSourceClassName to return the set name")
	}

	// setSourceMethodName/getSourceMethodName
	methodName := object.StringObjectFromGoString("doIt")
	_ = loggingLogRecordSetSourceMethodName([]interface{}{rec, methodName})
	if got := loggingLogRecordGetSourceMethodName([]interface{}{rec}); got != methodName {
		t.Fatalf("expected getSourceMethodName to return the set name")
	}

	// setParameters/getParameters
	params := object.MakeEmptyObject()
	_ = loggingLogRecordSetParameters([]interface{}{rec, params})
	if got := loggingLogRecordGetParameters([]interface{}{rec}); got != params {
		t.Fatalf("expected getParameters to return the set parameters")
	}

	// setThreadID/getThreadID
	_ = loggingLogRecordSetThreadID([]interface{}{rec, int64(7)})
	if got := loggingLogRecordGetThreadID([]interface{}{rec}).(int64); got != 7 {
		t.Fatalf("expected threadID 7, got %d", got)
	}

	// setThrown/getThrown
	thrown := object.MakeEmptyObject()
	_ = loggingLogRecordSetThrown([]interface{}{rec, thrown})
	if got := loggingLogRecordGetThrown([]interface{}{rec}); got != thrown {
		t.Fatalf("expected getThrown to return the set throwable")
	}

	// setResourceBundleName/getResourceBundleName
	rbName := object.StringObjectFromGoString("MyBundle")
	_ = loggingLogRecordSetResourceBundleName([]interface{}{rec, rbName})
	if got := loggingLogRecordGetResourceBundleName([]interface{}{rec}); got != rbName {
		t.Fatalf("expected getResourceBundleName to return the set name")
	}

	// setResourceBundle/getResourceBundle
	rb := object.MakeEmptyObject()
	_ = loggingLogRecordSetResourceBundle([]interface{}{rec, rb})
	if got := loggingLogRecordGetResourceBundle([]interface{}{rec}); got != rb {
		t.Fatalf("expected getResourceBundle to return the set bundle")
	}
}

func TestLoggingLogRecordErrorPaths(t *testing.T) {
	globals.InitStringPool()

	notAnObject := []interface{}{42}

	fns := []func([]interface{}) interface{}{
		loggingLogRecordInit,
		loggingLogRecordGetLevel,
		loggingLogRecordSetLevel,
		loggingLogRecordGetMessage,
		loggingLogRecordSetMessage,
		loggingLogRecordGetLoggerName,
		loggingLogRecordSetLoggerName,
		loggingLogRecordGetSequenceNumber,
		loggingLogRecordSetSequenceNumber,
		loggingLogRecordGetSourceClassName,
		loggingLogRecordSetSourceClassName,
		loggingLogRecordGetSourceMethodName,
		loggingLogRecordSetSourceMethodName,
		loggingLogRecordGetParameters,
		loggingLogRecordSetParameters,
		loggingLogRecordGetThreadID,
		loggingLogRecordSetThreadID,
		loggingLogRecordGetThrown,
		loggingLogRecordSetThrown,
		loggingLogRecordGetResourceBundleName,
		loggingLogRecordSetResourceBundleName,
		loggingLogRecordGetResourceBundle,
		loggingLogRecordSetResourceBundle,
	}

	for _, fn := range fns {
		params := append(notAnObject, nil, nil)
		ret := fn(params)
		if _, ok := ret.(*ghelpers.GErrBlk); !ok {
			t.Fatalf("expected error block for non-object first param, got %v (%T)", ret, ret)
		}
	}
}

func TestLoggingLogRecordMethodSignaturesRegistered(t *testing.T) {
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)
	Load_Util_Logging_LogRecord()

	expected := map[string]int{
		"java/util/logging/LogRecord.<init>(Ljava/util/logging/Level;Ljava/lang/String;)V": 2,
		"java/util/logging/LogRecord.getLevel()Ljava/util/logging/Level;":                  0,
		"java/util/logging/LogRecord.setLevel(Ljava/util/logging/Level;)V":                 1,
		"java/util/logging/LogRecord.getMessage()Ljava/lang/String;":                       0,
		"java/util/logging/LogRecord.setMessage(Ljava/lang/String;)V":                      1,
		"java/util/logging/LogRecord.getLoggerName()Ljava/lang/String;":                    0,
		"java/util/logging/LogRecord.setLoggerName(Ljava/lang/String;)V":                   1,
		"java/util/logging/LogRecord.getSequenceNumber()J":                                 0,
		"java/util/logging/LogRecord.setSequenceNumber(J)V":                                1,
		"java/util/logging/LogRecord.getSourceClassName()Ljava/lang/String;":               0,
		"java/util/logging/LogRecord.setSourceClassName(Ljava/lang/String;)V":              1,
		"java/util/logging/LogRecord.getSourceMethodName()Ljava/lang/String;":              0,
		"java/util/logging/LogRecord.setSourceMethodName(Ljava/lang/String;)V":             1,
		"java/util/logging/LogRecord.getParameters()[Ljava/lang/Object;":                   0,
		"java/util/logging/LogRecord.setParameters([Ljava/lang/Object;)V":                  1,
		"java/util/logging/LogRecord.getThreadID()I":                                       0,
		"java/util/logging/LogRecord.setThreadID(I)V":                                      1,
		"java/util/logging/LogRecord.getMillis()J":                                         0,
		"java/util/logging/LogRecord.setMillis(J)V":                                        1,
		"java/util/logging/LogRecord.getThrown()Ljava/lang/Throwable;":                     0,
		"java/util/logging/LogRecord.setThrown(Ljava/lang/Throwable;)V":                    1,
		"java/util/logging/LogRecord.getResourceBundleName()Ljava/lang/String;":            0,
		"java/util/logging/LogRecord.setResourceBundleName(Ljava/lang/String;)V":           1,
		"java/util/logging/LogRecord.getResourceBundle()Ljava/util/ResourceBundle;":        0,
		"java/util/logging/LogRecord.setResourceBundle(Ljava/util/ResourceBundle;)V":       1,
	}

	for key, slots := range expected {
		gm, ok := ghelpers.MethodSignatures[key]
		if !ok {
			t.Errorf("missing MethodSignatures entry for %s", key)
			continue
		}
		if gm.ParamSlots != slots {
			t.Errorf("%s: expected ParamSlots %d, got %d", key, slots, gm.ParamSlots)
		}
	}

	// getMillis/setMillis are deprecated as of Java 9 (in favor of
	// getInstant()/setInstant()), but real JDK still runs them (they are not
	// unsupported), so they must remain functional rather than trapped.
	functional := map[string]interface{}{
		"java/util/logging/LogRecord.getMillis()J":  loggingLogRecordGetMillis,
		"java/util/logging/LogRecord.setMillis(J)V": loggingLogRecordSetMillis,
	}
	for key, want := range functional {
		gm := ghelpers.MethodSignatures[key]
		if reflect.ValueOf(gm.GFunction).Pointer() != reflect.ValueOf(want).Pointer() {
			t.Errorf("%s: expected GFunction to be functional (not trapped)", key)
		}
	}
}
