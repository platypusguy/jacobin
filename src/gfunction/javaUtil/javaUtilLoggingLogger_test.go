package javaUtil

import (
	"io"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"os"
	"strings"
	"testing"
)

// captureLoggerOutput temporarily redirects java/lang/System.err to a pipe,
// runs fn, and returns everything written during that call.
func captureLoggerOutput(t *testing.T, fn func()) string {
	t.Helper()
	globals.InitStringPool()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	defer r.Close()

	defer func() {
		_ = statics.AddStatic("java/lang/System.err", statics.Static{Type: "GS", Value: os.Stderr})
	}()
	_ = statics.AddStatic("java/lang/System.err", statics.Static{Type: "GS", Value: w})

	fn()

	_ = w.Close()
	buf, _ := io.ReadAll(r)
	return string(buf)
}

func TestLoggingLoggerConfig(t *testing.T) {
	out := captureLoggerOutput(t, func() {
		ret := loggingLoggerConfig([]interface{}{nil, nil, object.StringObjectFromGoString("hello config")})
		if ret != nil {
			t.Fatalf("loggingLoggerConfig returned error: %v", ret)
		}
	})
	if out != "CONFIG: hello config\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggingLoggerEntering(t *testing.T) {
	out := captureLoggerOutput(t, func() {
		ret := loggingLoggerEntering([]interface{}{nil,
			object.StringObjectFromGoString("com.foo.Bar"),
			object.StringObjectFromGoString("doIt")})
		if ret != nil {
			t.Fatalf("loggingLoggerEntering returned error: %v", ret)
		}
	})
	if out != "FINER: ENTRY com.foo.Bar doIt\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggingLoggerEnteringWithParam(t *testing.T) {
	out := captureLoggerOutput(t, func() {
		ret := loggingLoggerEnteringWithParam([]interface{}{nil,
			object.StringObjectFromGoString("com.foo.Bar"),
			object.StringObjectFromGoString("doIt"),
			object.StringObjectFromGoString("arg1")})
		if ret != nil {
			t.Fatalf("loggingLoggerEnteringWithParam returned error: %v", ret)
		}
	})
	if !strings.HasPrefix(out, "FINER: ENTRY com.foo.Bar doIt") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggingLoggerExiting(t *testing.T) {
	out := captureLoggerOutput(t, func() {
		ret := loggingLoggerExiting([]interface{}{nil,
			object.StringObjectFromGoString("com.foo.Bar"),
			object.StringObjectFromGoString("doIt")})
		if ret != nil {
			t.Fatalf("loggingLoggerExiting returned error: %v", ret)
		}
	})
	if out != "FINER: RETURN com.foo.Bar doIt\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggingLoggerExitingWithParam(t *testing.T) {
	out := captureLoggerOutput(t, func() {
		ret := loggingLoggerExitingWithParam([]interface{}{nil,
			object.StringObjectFromGoString("com.foo.Bar"),
			object.StringObjectFromGoString("doIt"),
			object.StringObjectFromGoString("result1")})
		if ret != nil {
			t.Fatalf("loggingLoggerExitingWithParam returned error: %v", ret)
		}
	})
	if out != "FINER: RETURN com.foo.Bar doIt result1\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggingLoggerFinestFinerInfoSevereWarning(t *testing.T) {
	cases := []struct {
		fn       func(params []interface{}) interface{}
		prefix   string
		expected string
	}{
		{loggingLoggerFinest, "FINEST", "FINEST: msg1\n"},
		{loggingLoggerFiner, "FINER", "FINER: msg1\n"},
		{loggingLoggerInfo, "INFO", "INFO: msg1\n"},
		{loggingLoggerSevere, "SEVERE", "SEVERE: msg1\n"},
		{loggingLoggerWarning, "WARNING", "WARNING: msg1\n"},
	}
	for _, c := range cases {
		out := captureLoggerOutput(t, func() {
			ret := c.fn([]interface{}{nil, nil, object.StringObjectFromGoString("msg1")})
			if ret != nil {
				t.Fatalf("%s returned error: %v", c.prefix, ret)
			}
		})
		if out != c.expected {
			t.Fatalf("%s: unexpected output: %q", c.prefix, out)
		}
	}
}

func TestLoggingLoggerWarningWithParams(t *testing.T) {
	out := captureLoggerOutput(t, func() {
		arr := object.MakeEmptyObject()
		arr.FieldTable["value"] = object.Field{Fvalue: []*object.Object{
			object.StringObjectFromGoString("a"),
			object.StringObjectFromGoString("b"),
		}}
		ret := loggingLoggerWarningWithParams([]interface{}{nil, nil,
			object.StringObjectFromGoString("msg2"), arr})
		if ret != nil {
			t.Fatalf("loggingLoggerWarningWithParams returned error: %v", ret)
		}
	})
	if out != "WARNING: msg2 [a, b]\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggingLoggerLog(t *testing.T) {
	globals.InitStringPool()
	level := makeLevelObject("SEVERE", standardLevels["SEVERE"], "")

	out := captureLoggerOutput(t, func() {
		ret := loggingLoggerLog([]interface{}{nil, nil, level, object.StringObjectFromGoString("bad thing")})
		if ret != nil {
			t.Fatalf("loggingLoggerLog returned error: %v", ret)
		}
	})
	if out != "SEVERE: bad thing\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggingLoggerLogWithParams(t *testing.T) {
	globals.InitStringPool()
	level := makeLevelObject("INFO", standardLevels["INFO"], "")

	out := captureLoggerOutput(t, func() {
		arr := object.MakeEmptyObject()
		arr.FieldTable["value"] = object.Field{Fvalue: []*object.Object{
			object.StringObjectFromGoString("x"),
		}}
		ret := loggingLoggerLogWithParams([]interface{}{nil, nil, level, object.StringObjectFromGoString("msg"), arr})
		if ret != nil {
			t.Fatalf("loggingLoggerLogWithParams returned error: %v", ret)
		}
	})
	if out != "INFO: msg [x]\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggingLoggerLogp(t *testing.T) {
	globals.InitStringPool()
	level := makeLevelObject("WARNING", standardLevels["WARNING"], "")

	out := captureLoggerOutput(t, func() {
		ret := loggingLoggerLogp([]interface{}{nil, level,
			object.StringObjectFromGoString("com.foo.Bar"),
			object.StringObjectFromGoString("doIt"),
			object.StringObjectFromGoString("careful")})
		if ret != nil {
			t.Fatalf("loggingLoggerLogp returned error: %v", ret)
		}
	})
	if out != "WARNING: com.foo.Bar doIt: careful\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggingLoggerLogpWithParams(t *testing.T) {
	globals.InitStringPool()
	level := makeLevelObject("CONFIG", standardLevels["CONFIG"], "")

	out := captureLoggerOutput(t, func() {
		arr := object.MakeEmptyObject()
		arr.FieldTable["value"] = object.Field{Fvalue: []*object.Object{
			object.StringObjectFromGoString("y"),
		}}
		ret := loggingLoggerLogpWithParams([]interface{}{nil, level,
			object.StringObjectFromGoString("com.foo.Bar"),
			object.StringObjectFromGoString("doIt"),
			object.StringObjectFromGoString("msg"), arr})
		if ret != nil {
			t.Fatalf("loggingLoggerLogpWithParams returned error: %v", ret)
		}
	})
	if out != "CONFIG: com.foo.Bar doIt: msg [y]\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggingLoggerLogrb(t *testing.T) {
	globals.InitStringPool()
	level := makeLevelObject("SEVERE", standardLevels["SEVERE"], "")

	out := captureLoggerOutput(t, func() {
		ret := loggingLoggerLogrb([]interface{}{nil, level,
			object.StringObjectFromGoString("com.foo.Bar"),
			object.StringObjectFromGoString("doIt"),
			object.StringObjectFromGoString("some.bundle"),
			object.StringObjectFromGoString("boom")})
		if ret != nil {
			t.Fatalf("loggingLoggerLogrb returned error: %v", ret)
		}
	})
	if out != "SEVERE: com.foo.Bar doIt: boom\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggingLoggerLogrbWithParams(t *testing.T) {
	globals.InitStringPool()
	level := makeLevelObject("INFO", standardLevels["INFO"], "")

	out := captureLoggerOutput(t, func() {
		arr := object.MakeEmptyObject()
		arr.FieldTable["value"] = object.Field{Fvalue: []*object.Object{
			object.StringObjectFromGoString("z"),
		}}
		ret := loggingLoggerLogrbWithParams([]interface{}{nil, level,
			object.StringObjectFromGoString("com.foo.Bar"),
			object.StringObjectFromGoString("doIt"),
			object.StringObjectFromGoString("some.bundle"),
			object.StringObjectFromGoString("msg"), arr})
		if ret != nil {
			t.Fatalf("loggingLoggerLogrbWithParams returned error: %v", ret)
		}
	})
	if out != "INFO: com.foo.Bar doIt: msg [z]\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggingLoggerThrowing(t *testing.T) {
	out := captureLoggerOutput(t, func() {
		ret := loggingLoggerThrowing([]interface{}{nil,
			object.StringObjectFromGoString("com.foo.Bar"),
			object.StringObjectFromGoString("doIt"),
			object.MakeEmptyObject()})
		if ret != nil {
			t.Fatalf("loggingLoggerThrowing returned error: %v", ret)
		}
	})
	if out != "FINER: THROW com.foo.Bar doIt\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggingLoggerLog_NilLevelDefaultsToInfo(t *testing.T) {
	globals.InitStringPool()

	out := captureLoggerOutput(t, func() {
		ret := loggingLoggerLog([]interface{}{nil, nil, object.Null, object.StringObjectFromGoString("msg")})
		if ret != nil {
			t.Fatalf("loggingLoggerLog returned error: %v", ret)
		}
	})
	if out != "INFO: msg\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggingLoggerFine(t *testing.T) {
	out := captureLoggerOutput(t, func() {
		ret := loggingLoggerFine([]interface{}{nil, nil, object.StringObjectFromGoString("fine msg")})
		if ret != nil {
			t.Fatalf("loggingLoggerFine returned error: %v", ret)
		}
	})
	if out != "FINE: fine msg\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestLoggingLoggerGetLoggerReturnsSameInstance(t *testing.T) {
	globals.InitStringPool()
	name := object.StringObjectFromGoString("com.example.TestLoggerA")

	l1 := loggingLoggerGetLogger([]interface{}{name}).(*object.Object)
	l2 := loggingLoggerGetLogger([]interface{}{name}).(*object.Object)
	if l1 != l2 {
		t.Fatalf("getLogger should return the same instance for the same name")
	}

	gotName := loggingLoggerGetName([]interface{}{l1})
	if object.GoStringFromStringObject(gotName.(*object.Object)) != "com.example.TestLoggerA" {
		t.Fatalf("unexpected logger name: %v", gotName)
	}

	parent := loggingLoggerGetParent([]interface{}{l1})
	if parent == nil || object.IsNull(parent.(*object.Object)) {
		t.Fatalf("expected non-null parent (the global logger)")
	}
}

func TestLoggingLoggerGetLoggerWithBundle(t *testing.T) {
	globals.InitStringPool()
	name := object.StringObjectFromGoString("com.example.TestLoggerBundle")
	bundle := object.StringObjectFromGoString("some.bundle")

	logger := loggingLoggerGetLoggerWithBundle([]interface{}{name, bundle}).(*object.Object)
	rbName := loggingLoggerGetResourceBundleName([]interface{}{logger})
	if object.GoStringFromStringObject(rbName.(*object.Object)) != "some.bundle" {
		t.Fatalf("unexpected resource bundle name: %v", rbName)
	}
}

func TestLoggingLoggerGetGlobal(t *testing.T) {
	globals.InitStringPool()
	g1 := loggingLoggerGetGlobal(nil).(*object.Object)
	g2 := loggingLoggerGetGlobal(nil).(*object.Object)
	if g1 != g2 {
		t.Fatalf("getGlobal should always return the same instance")
	}
}

func TestLoggingLoggerGetAnonymousLogger(t *testing.T) {
	globals.InitStringPool()
	l := loggingLoggerGetAnonymousLogger(nil).(*object.Object)
	name := loggingLoggerGetName([]interface{}{l})
	if name != object.Null {
		t.Fatalf("expected anonymous logger to have a null name, got: %v", name)
	}

	l2 := loggingLoggerGetAnonymousLoggerWithBundle([]interface{}{object.StringObjectFromGoString("bundleX")}).(*object.Object)
	rbName := loggingLoggerGetResourceBundleName([]interface{}{l2})
	if object.GoStringFromStringObject(rbName.(*object.Object)) != "bundleX" {
		t.Fatalf("unexpected resource bundle name: %v", rbName)
	}
}

func TestLoggingLoggerLevelGetSet(t *testing.T) {
	globals.InitStringPool()
	l := makeLoggerObject("levelTest", "")

	if got := loggingLoggerGetLevel([]interface{}{l}); got != object.Null {
		t.Fatalf("expected null level by default, got: %v", got)
	}

	warnLevel := makeLevelObject("WARNING", standardLevels["WARNING"], "")
	if ret := loggingLoggerSetLevel([]interface{}{l, warnLevel}); ret != nil {
		t.Fatalf("setLevel returned error: %v", ret)
	}
	got := loggingLoggerGetLevel([]interface{}{l})
	if got.(*object.Object) != warnLevel {
		t.Fatalf("getLevel did not return the level that was set")
	}
}

func TestLoggingLoggerIsLoggable(t *testing.T) {
	globals.InitStringPool()
	l := makeLoggerObject("isLoggableTest", "")
	_ = loggingLoggerSetLevel([]interface{}{l, makeLevelObject("WARNING", standardLevels["WARNING"], "")})

	infoLevel := makeLevelObject("INFO", standardLevels["INFO"], "")
	if got := loggingLoggerIsLoggable([]interface{}{l, infoLevel}); got != types.JavaBoolFalse {
		t.Fatalf("expected INFO not loggable when level is WARNING, got: %v", got)
	}

	severeLevel := makeLevelObject("SEVERE", standardLevels["SEVERE"], "")
	if got := loggingLoggerIsLoggable([]interface{}{l, severeLevel}); got != types.JavaBoolTrue {
		t.Fatalf("expected SEVERE loggable when level is WARNING, got: %v", got)
	}
}

func TestLoggingLoggerIsLoggableInheritsFromParent(t *testing.T) {
	globals.InitStringPool()
	parent := makeLoggerObject("parentLogger", "")
	_ = loggingLoggerSetLevel([]interface{}{parent, makeLevelObject("SEVERE", standardLevels["SEVERE"], "")})

	child := makeLoggerObject("childLogger", "")
	_ = loggingLoggerSetParent([]interface{}{child, parent})

	warnLevel := makeLevelObject("WARNING", standardLevels["WARNING"], "")
	if got := loggingLoggerIsLoggable([]interface{}{child, warnLevel}); got != types.JavaBoolFalse {
		t.Fatalf("expected WARNING not loggable when effective (parent) level is SEVERE, got: %v", got)
	}

	gotParent := loggingLoggerGetParent([]interface{}{child})
	if gotParent.(*object.Object) != parent {
		t.Fatalf("getParent did not return the parent that was set")
	}
}

func TestLoggingLoggerFilterGetSet(t *testing.T) {
	globals.InitStringPool()
	l := makeLoggerObject("filterTest", "")

	if got := loggingLoggerGetFilter([]interface{}{l}); got != object.Null {
		t.Fatalf("expected null filter by default, got: %v", got)
	}

	filterObj := object.MakeEmptyObject()
	if ret := loggingLoggerSetFilter([]interface{}{l, filterObj}); ret != nil {
		t.Fatalf("setFilter returned error: %v", ret)
	}
	if got := loggingLoggerGetFilter([]interface{}{l}); got.(*object.Object) != filterObj {
		t.Fatalf("getFilter did not return the filter that was set")
	}
}

func TestLoggingLoggerUseParentHandlersGetSet(t *testing.T) {
	globals.InitStringPool()
	l := makeLoggerObject("useParentHandlersTest", "")

	if got := loggingLoggerGetUseParentHandlers([]interface{}{l}); got != types.JavaBoolTrue {
		t.Fatalf("expected useParentHandlers to default to true, got: %v", got)
	}

	if ret := loggingLoggerSetUseParentHandlers([]interface{}{l, types.JavaBoolFalse}); ret != nil {
		t.Fatalf("setUseParentHandlers returned error: %v", ret)
	}
	if got := loggingLoggerGetUseParentHandlers([]interface{}{l}); got != types.JavaBoolFalse {
		t.Fatalf("expected useParentHandlers to be false after setting, got: %v", got)
	}
}

func TestLoggingLoggerHandlers(t *testing.T) {
	globals.InitStringPool()
	l := makeLoggerObject("handlersTest", "")

	initial := loggingLoggerGetHandlers([]interface{}{l}).(*object.Object)
	if arr, _ := initial.FieldTable["value"].Fvalue.([]*object.Object); len(arr) != 0 {
		t.Fatalf("expected no handlers initially, got: %v", arr)
	}

	h1 := object.MakeEmptyObject()
	h2 := object.MakeEmptyObject()
	_ = loggingLoggerAddHandler([]interface{}{l, h1})
	_ = loggingLoggerAddHandler([]interface{}{l, h2})

	after := loggingLoggerGetHandlers([]interface{}{l}).(*object.Object)
	arr, _ := after.FieldTable["value"].Fvalue.([]*object.Object)
	if len(arr) != 2 {
		t.Fatalf("expected 2 handlers, got: %d", len(arr))
	}

	_ = loggingLoggerRemoveHandler([]interface{}{l, h1})
	afterRemove := loggingLoggerGetHandlers([]interface{}{l}).(*object.Object)
	arr2, _ := afterRemove.FieldTable["value"].Fvalue.([]*object.Object)
	if len(arr2) != 1 || arr2[0] != h2 {
		t.Fatalf("expected only h2 to remain after removeHandler, got: %v", arr2)
	}
}

func TestLoggingLoggerResourceBundle(t *testing.T) {
	globals.InitStringPool()
	l := makeLoggerObject("resourceBundleTest", "")
	if got := loggingLoggerGetResourceBundle([]interface{}{l}); got != object.Null {
		t.Fatalf("expected null resource bundle by default, got: %v", got)
	}
	if got := loggingLoggerGetResourceBundleName([]interface{}{l}); got != object.Null {
		t.Fatalf("expected null resource bundle name by default, got: %v", got)
	}
}

func TestLoggingLoggerInit(t *testing.T) {
	globals.InitStringPool()
	obj := object.MakeEmptyObject()
	ret := loggingLoggerInit([]interface{}{obj,
		object.StringObjectFromGoString("myLogger"),
		object.StringObjectFromGoString("myBundle")})
	if ret != nil {
		t.Fatalf("loggingLoggerInit returned error: %v", ret)
	}

	name := loggingLoggerGetName([]interface{}{obj})
	if object.GoStringFromStringObject(name.(*object.Object)) != "myLogger" {
		t.Fatalf("unexpected name after init: %v", name)
	}
	rbName := loggingLoggerGetResourceBundleName([]interface{}{obj})
	if object.GoStringFromStringObject(rbName.(*object.Object)) != "myBundle" {
		t.Fatalf("unexpected resource bundle name after init: %v", rbName)
	}
}

func TestLoggingLoggerErrorPaths(t *testing.T) {
	globals.InitStringPool()
	badParam := []interface{}{"not an object"}

	fns := map[string]func([]interface{}) interface{}{
		"getName":               loggingLoggerGetName,
		"getLevel":              loggingLoggerGetLevel,
		"getParent":             loggingLoggerGetParent,
		"getFilter":             loggingLoggerGetFilter,
		"getUseParentHandlers":  loggingLoggerGetUseParentHandlers,
		"getResourceBundleName": loggingLoggerGetResourceBundleName,
		"getResourceBundle":     loggingLoggerGetResourceBundle,
		"getHandlers":           loggingLoggerGetHandlers,
		"isLoggable":            loggingLoggerIsLoggable,
	}
	for name, fn := range fns {
		if _, ok := fn(badParam).(*ghelpers.GErrBlk); !ok {
			t.Fatalf("%s: expected GErrBlk for non-object receiver", name)
		}
	}

	twoArgFns := map[string]func([]interface{}) interface{}{
		"setLevel":             loggingLoggerSetLevel,
		"setParent":            loggingLoggerSetParent,
		"setFilter":            loggingLoggerSetFilter,
		"setUseParentHandlers": loggingLoggerSetUseParentHandlers,
		"addHandler":           loggingLoggerAddHandler,
		"removeHandler":        loggingLoggerRemoveHandler,
	}
	for name, fn := range twoArgFns {
		if _, ok := fn([]interface{}{"not an object", nil}).(*ghelpers.GErrBlk); !ok {
			t.Fatalf("%s: expected GErrBlk for non-object receiver", name)
		}
	}
}
