/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func zoneRegionTestInit() {
	globals.InitGlobals("test")
	Load_Time_ZoneOffset()
	Load_Time_Zone_ZoneRules()
	Load_Time_ZoneId()
	Load_Time_ZoneRegion()
}

func TestZoneRegionRegistration(t *testing.T) {
	zoneRegionTestInit()

	expectedSignatures := []string{
		"java/time/ZoneRegion.<clinit>()V",
		"java/time/ZoneRegion.<init>()V",
		"java/time/ZoneRegion.<init>(Ljava/lang/String;Ljava/time/zone/ZoneRules;)V",
		"java/time/ZoneRegion.ofId(Ljava/lang/String;Z)Ljava/time/ZoneRegion;",
		"java/time/ZoneRegion.ofWithPrefix(Ljava/lang/String;IZ)Ljava/time/ZoneRegion;",
		"java/time/ZoneRegion.checkName(Ljava/lang/String;)V",
		"java/time/ZoneRegion.getId()Ljava/lang/String;",
		"java/time/ZoneRegion.getRules()Ljava/time/zone/ZoneRules;",
		"java/time/ZoneRegion.getDisplayName(Ljava/time/format/TextStyle;Ljava/util/Locale;)Ljava/lang/String;",
		"java/time/ZoneRegion.normalized()Ljava/time/ZoneId;",
		"java/time/ZoneRegion.equals(Ljava/lang/Object;)Z",
		"java/time/ZoneRegion.hashCode()I",
		"java/time/ZoneRegion.toString()Ljava/lang/String;",
		"java/time/ZoneRegion.write(Ljava/io/DataOutput;)V",
		"java/time/ZoneRegion.writeExternal(Ljava/io/DataOutput;)V",
		"java/time/ZoneRegion.readExternal(Ljava/io/DataInput;)Ljava/time/ZoneId;",
	}

	for _, sig := range expectedSignatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Method signature %s not registered", sig)
		}
	}
}

func TestZoneRegionOfId(t *testing.T) {
	zoneRegionTestInit()

	// Valid region
	res := zoneRegionOfId([]any{object.StringObjectFromGoString("Europe/London"), int64(1)})
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneRegion.ofId('Europe/London', true) error: %v", err)
	}
	obj := res.(*object.Object)
	if getZoneIdString(obj) != "Europe/London" {
		t.Fatalf("ZoneRegion.ofId('Europe/London').getId() = %v, want Europe/London", getZoneIdString(obj))
	}

	// Null zoneId -> NPE
	resNull := zoneRegionOfId([]any{object.Null, int64(1)})
	if err, isErr := resNull.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("ZoneRegion.ofId(null) expected NPE, got %v", resNull)
	} else if err.ExceptionType != excNames.NullPointerException {
		t.Fatalf("ZoneRegion.ofId(null) expected NPE, got %v", err.ExceptionType)
	}

	// Invalid name (too short) -> DateTimeException
	resShort := zoneRegionOfId([]any{object.StringObjectFromGoString("A"), int64(1)})
	if err, isErr := resShort.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("ZoneRegion.ofId('A') expected DateTimeException, got %v", resShort)
	} else if err.ExceptionType != excNames.DateTimeException {
		t.Fatalf("ZoneRegion.ofId('A') expected DateTimeException, got %v", err.ExceptionType)
	}

	// Invalid characters -> DateTimeException
	resInv := zoneRegionOfId([]any{object.StringObjectFromGoString("Invalid@Zone"), int64(1)})
	if err, isErr := resInv.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("ZoneRegion.ofId('Invalid@Zone') expected DateTimeException, got %v", resInv)
	} else if err.ExceptionType != excNames.DateTimeException {
		t.Fatalf("ZoneRegion.ofId('Invalid@Zone') expected DateTimeException, got %v", err.ExceptionType)
	}

	// Unknown region with checkAvailable=true -> DateTimeException
	resUnknown := zoneRegionOfId([]any{object.StringObjectFromGoString("NonExistent/Zone_12345"), int64(1)})
	if err, isErr := resUnknown.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("ZoneRegion.ofId('NonExistent/Zone_12345', true) expected DateTimeException, got %v", resUnknown)
	} else if err.ExceptionType != excNames.DateTimeException {
		t.Fatalf("ZoneRegion.ofId('NonExistent/Zone_12345', true) expected DateTimeException, got %v", err.ExceptionType)
	}

	// Unknown region with checkAvailable=false -> creates ZoneRegion
	resUnknownNoCheck := zoneRegionOfId([]any{object.StringObjectFromGoString("NonExistent/Zone_12345"), int64(0)})
	if err, isErr := resUnknownNoCheck.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneRegion.ofId('NonExistent/Zone_12345', false) error: %v", err)
	}
	objNoCheck := resUnknownNoCheck.(*object.Object)
	if getZoneIdString(objNoCheck) != "NonExistent/Zone_12345" {
		t.Fatalf("ZoneRegion.ofId('NonExistent/Zone_12345', false).getId() = %v", getZoneIdString(objNoCheck))
	}
}

func TestZoneRegionOfWithPrefix(t *testing.T) {
	zoneRegionTestInit()

	// "GMT+02:00", prefixLength = 3
	res1 := zoneRegionOfWithPrefix([]any{object.StringObjectFromGoString("GMT+02:00"), int64(3), int64(1)})
	if err, isErr := res1.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneRegion.ofWithPrefix('GMT+02:00', 3, true) error: %v", err)
	}
	obj1 := res1.(*object.Object)
	if getZoneIdString(obj1) != "GMT+02:00" {
		t.Fatalf("ZoneRegion.ofWithPrefix('GMT+02:00').getId() = %v, want GMT+02:00", getZoneIdString(obj1))
	}

	// "UTC", prefixLength = 3 (length == prefixLength)
	res2 := zoneRegionOfWithPrefix([]any{object.StringObjectFromGoString("UTC"), int64(3), int64(1)})
	if err, isErr := res2.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneRegion.ofWithPrefix('UTC', 3, true) error: %v", err)
	}
	obj2 := res2.(*object.Object)
	if getZoneIdString(obj2) != "UTC" {
		t.Fatalf("ZoneRegion.ofWithPrefix('UTC').getId() = %v, want UTC", getZoneIdString(obj2))
	}

	// "GMT+00:00", prefixLength = 3 -> offset is 0, returns ofId("GMT")
	res3 := zoneRegionOfWithPrefix([]any{object.StringObjectFromGoString("GMT+00:00"), int64(3), int64(1)})
	if err, isErr := res3.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneRegion.ofWithPrefix('GMT+00:00', 3, true) error: %v", err)
	}
	obj3 := res3.(*object.Object)
	if getZoneIdString(obj3) != "GMT" {
		t.Fatalf("ZoneRegion.ofWithPrefix('GMT+00:00').getId() = %v, want GMT", getZoneIdString(obj3))
	}

	// "GMT0", prefixLength = 3 (char at prefixLength is not '+' or '-') -> ofId("GMT0")
	res4 := zoneRegionOfWithPrefix([]any{object.StringObjectFromGoString("GMT0"), int64(3), int64(0)})
	if err, isErr := res4.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneRegion.ofWithPrefix('GMT0', 3, false) error: %v", err)
	}
	obj4 := res4.(*object.Object)
	if getZoneIdString(obj4) != "GMT0" {
		t.Fatalf("ZoneRegion.ofWithPrefix('GMT0').getId() = %v, want GMT0", getZoneIdString(obj4))
	}

	// Invalid offset part: "GMT+99:99" -> DateTimeException
	resInv := zoneRegionOfWithPrefix([]any{object.StringObjectFromGoString("GMT+99:99"), int64(3), int64(1)})
	if err, isErr := resInv.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("ZoneRegion.ofWithPrefix('GMT+99:99') expected DateTimeException, got %v", resInv)
	} else if err.ExceptionType != excNames.DateTimeException {
		t.Fatalf("ZoneRegion.ofWithPrefix('GMT+99:99') expected DateTimeException, got %v", err.ExceptionType)
	}

	// Null zoneId -> NPE
	resNull := zoneRegionOfWithPrefix([]any{object.Null, int64(3), int64(1)})
	if err, isErr := resNull.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("ZoneRegion.ofWithPrefix(null) expected NPE, got %v", resNull)
	} else if err.ExceptionType != excNames.NullPointerException {
		t.Fatalf("ZoneRegion.ofWithPrefix(null) expected NPE, got %v", err.ExceptionType)
	}

	// Invalid prefixLength out of bounds -> DateTimeException
	resBounds := zoneRegionOfWithPrefix([]any{object.StringObjectFromGoString("UTC"), int64(10), int64(1)})
	if err, isErr := resBounds.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("ZoneRegion.ofWithPrefix(len 10) expected DateTimeException, got %v", resBounds)
	} else if err.ExceptionType != excNames.DateTimeException {
		t.Fatalf("ZoneRegion.ofWithPrefix(len 10) expected DateTimeException, got %v", err.ExceptionType)
	}
}

func TestZoneRegionCheckName(t *testing.T) {
	zoneRegionTestInit()

	// Valid names
	validNames := []string{"Europe/London", "UTC", "America/New_York", "Etc/GMT+1", "GMT-8", "A/B~._+-1"}
	for _, name := range validNames {
		res := zoneRegionCheckName([]any{object.StringObjectFromGoString(name)})
		if res != nil {
			t.Errorf("ZoneRegion.checkName(%q) returned unexpected error: %v", name, res)
		}
	}

	// Invalid names
	invalidNames := []string{"A", "", "/London", "123", "Invalid@Zone", "+02:00", "-05:00"}
	for _, name := range invalidNames {
		res := zoneRegionCheckName([]any{object.StringObjectFromGoString(name)})
		if err, isErr := res.(*ghelpers.GErrBlk); !isErr {
			t.Errorf("ZoneRegion.checkName(%q) expected error, got nil", name)
		} else if err.ExceptionType != excNames.DateTimeException {
			t.Errorf("ZoneRegion.checkName(%q) expected DateTimeException, got %v", name, err.ExceptionType)
		}
	}

	// Null -> NPE
	resNull := zoneRegionCheckName([]any{object.Null})
	if err, isErr := resNull.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("ZoneRegion.checkName(null) expected NPE, got %v", resNull)
	} else if err.ExceptionType != excNames.NullPointerException {
		t.Fatalf("ZoneRegion.checkName(null) expected NPE, got %v", err.ExceptionType)
	}
}

func TestZoneRegionInitAndMethods(t *testing.T) {
	zoneRegionTestInit()

	// <init>()V
	resInit := zoneRegionInit([]any{})
	if resInit != nil {
		t.Fatalf("ZoneRegion.<init>() returned %v", resInit)
	}

	// <init>(String, ZoneRules)
	zrObj := object.MakeEmptyObjectWithClassName(&classNameZoneRegionVal)
	rulesObj := createZoneRulesForZoneId("Europe/London")
	idObj := object.StringObjectFromGoString("Europe/London")
	resInit2 := zoneRegionInitWithIdRules([]any{zrObj, idObj, rulesObj})
	if resInit2 != nil {
		t.Fatalf("ZoneRegion.<init>(id, rules) returned error: %v", resInit2)
	}
	if getZoneIdString(zrObj) != "Europe/London" {
		t.Fatalf("ZoneRegion.getId() = %v, want Europe/London", getZoneIdString(zrObj))
	}

	// getRules()
	rulesRes := zoneIdGetRules([]any{zrObj})
	if err, isErr := rulesRes.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneRegion.getRules() error: %v", err)
	}
	if rulesRes.(*object.Object) != rulesObj {
		t.Fatalf("ZoneRegion.getRules() returned unexpected object")
	}

	// normalized() on Europe/London returns self
	normRes := zoneIdNormalized([]any{zrObj})
	if normRes.(*object.Object) != zrObj {
		t.Fatalf("ZoneRegion.normalized() on dynamic region expected self")
	}

	// equals, hashCode, toString
	zrObj2 := createZoneRegion("Europe/London", rulesObj)
	if eq := zoneIdEquals([]any{zrObj, zrObj2}); eq != types.JavaBoolTrue {
		t.Fatalf("ZoneRegion.equals(zrObj, zrObj2) = %v, want true", eq)
	}
	if h1, h2 := zoneIdHashCode([]any{zrObj}).(int64), zoneIdHashCode([]any{zrObj2}).(int64); h1 != h2 {
		t.Fatalf("ZoneRegion hashCode mismatch: %d != %d", h1, h2)
	}
	strRes := zoneIdToString([]any{zrObj}).(*object.Object)
	if object.GoStringFromStringObject(strRes) != "Europe/London" {
		t.Fatalf("ZoneRegion.toString() = %v, want Europe/London", object.GoStringFromStringObject(strRes))
	}

	// write, writeExternal, readExternal
	wRes := zoneRegionWrite([]any{zrObj, zrObj})
	if wRes != nil {
		t.Fatalf("ZoneRegion.write() returned %v", wRes)
	}
	weRes := zoneRegionWriteExternal([]any{zrObj, zrObj})
	if weRes != nil {
		t.Fatalf("ZoneRegion.writeExternal() returned %v", weRes)
	}
	reRes := zoneRegionReadExternal([]any{zrObj})
	if _, ok := reRes.(*object.Object); !ok {
		t.Fatalf("ZoneRegion.readExternal() returned non-object: %v", reRes)
	}
}
