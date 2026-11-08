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

func zoneRulesProviderTestInit() {
	globals.InitGlobals("test")
	Load_Time_ZoneOffset()
	Load_Time_Zone_ZoneRules()
	Load_Time_Zone_ZoneRulesException()
	Load_Time_Zone_ZoneRulesProvider()
	Load_Time_ZoneId()
	Load_Time_ZoneRegion()
}

func TestZoneRulesProviderRegistration(t *testing.T) {
	zoneRulesProviderTestInit()

	expectedSignatures := []string{
		"java/time/zone/ZoneRulesProvider.<clinit>()V",
		"java/time/zone/ZoneRulesProvider.getAvailableZoneIds()Ljava/util/Set;",
		"java/time/zone/ZoneRulesProvider.getRules(Ljava/lang/String;Z)Ljava/time/zone/ZoneRules;",
		"java/time/zone/ZoneRulesProvider.getVersions(Ljava/lang/String;)Ljava/util/NavigableMap;",
		"java/time/zone/ZoneRulesProvider.registerProvider(Ljava/time/zone/ZoneRulesProvider;)V",
		"java/time/zone/ZoneRulesProvider.refresh()Z",
		"java/time/zone/ZoneRulesProvider.<init>()V",
		"java/time/zone/ZoneRulesProvider.provideZoneIds()Ljava/util/Set;",
		"java/time/zone/ZoneRulesProvider.provideRules(Ljava/lang/String;Z)Ljava/time/zone/ZoneRules;",
		"java/time/zone/ZoneRulesProvider.provideVersions(Ljava/lang/String;)Ljava/util/NavigableMap;",
		"java/time/zone/ZoneRulesProvider.provideRefresh()Z",
		"java/time/zone/TzdbZoneRulesProvider.<clinit>()V",
		"java/time/zone/TzdbZoneRulesProvider.<init>()V",
		"java/time/zone/TzdbZoneRulesProvider.provideZoneIds()Ljava/util/Set;",
		"java/time/zone/TzdbZoneRulesProvider.provideRules(Ljava/lang/String;Z)Ljava/time/zone/ZoneRules;",
		"java/time/zone/TzdbZoneRulesProvider.provideVersions(Ljava/lang/String;)Ljava/util/NavigableMap;",
		"java/time/zone/TzdbZoneRulesProvider.provideRefresh()Z",
		"java/time/zone/TzdbZoneRulesProvider.toString()Ljava/lang/String;",
	}

	for _, sig := range expectedSignatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Method signature %s not registered", sig)
		}
	}
}

func TestZoneRulesProviderGetAvailableZoneIds(t *testing.T) {
	zoneRulesProviderTestInit()

	res := zoneRulesProviderGetAvailableZoneIds([]any{})
	if res == nil {
		t.Fatalf("zoneRulesProviderGetAvailableZoneIds returned nil")
	}
	setObj, ok := res.(*object.Object)
	if !ok || setObj == nil {
		t.Fatalf("Expected Set object, got %T", res)
	}

	mapFld, ok := setObj.FieldTable["map"]
	if !ok {
		t.Fatalf("Set object has no 'map' field")
	}
	hm, ok := mapFld.Fvalue.(types.DefHashMap)
	if !ok {
		t.Fatalf("Set field 'map' is not DefHashMap")
	}
	if len(hm) == 0 {
		t.Errorf("Available zone IDs set is empty")
	}
	hasUT := false
	hasUTC := false
	for _, v := range hm {
		if strObj, ok := v.(*object.Object); ok {
			str := object.GoStringFromStringObject(strObj)
			if str == "UT" {
				hasUT = true
			}
			if str == "UTC" {
				hasUTC = true
			}
		}
	}
	if hasUT {
		t.Errorf("Available zone IDs set should NOT contain 'UT'")
	}
	if !hasUTC {
		t.Errorf("Available zone IDs set must contain 'UTC'")
	}
}

func TestZoneRulesProviderGetRules(t *testing.T) {
	zoneRulesProviderTestInit()

	// Null parameter -> NPE
	resNull := zoneRulesProviderGetRules([]any{object.Null, types.JavaBoolFalse})
	if err, isErr := resNull.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("Expected NPE for null zoneId, got %v", resNull)
	} else if err.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NPE, got %v", err.ExceptionType)
	}

	// Valid UTC
	resUTC := zoneRulesProviderGetRules([]any{object.StringObjectFromGoString("UTC"), types.JavaBoolFalse})
	if err, isErr := resUTC.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("Error getting rules for UTC: %v", err)
	}
	utcObj := resUTC.(*object.Object)
	if !isZoneRulesFixed(utcObj) {
		t.Errorf("Expected UTC rules to be fixed offset")
	}

	// Valid GMT+03:00
	resOffset := zoneRulesProviderGetRules([]any{object.StringObjectFromGoString("+03:00"), types.JavaBoolFalse})
	if err, isErr := resOffset.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("Error getting rules for +03:00: %v", err)
	}
	offsetObj := resOffset.(*object.Object)
	if !isZoneRulesFixed(offsetObj) {
		t.Errorf("Expected +03:00 rules to be fixed offset")
	}

	// Valid region ID
	resNY := zoneRulesProviderGetRules([]any{object.StringObjectFromGoString("America/New_York"), types.JavaBoolFalse})
	if err, isErr := resNY.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("Error getting rules for America/New_York: %v", err)
	}
	nyObj := resNY.(*object.Object)
	if isZoneRulesFixed(nyObj) {
		t.Errorf("Expected America/New_York rules not to be fixed offset")
	}

	// Unknown zone ID -> DateTimeException
	resUnknown := zoneRulesProviderGetRules([]any{object.StringObjectFromGoString("Invalid/Unknown_Zone_12345"), types.JavaBoolFalse})
	if err, isErr := resUnknown.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("Expected DateTimeException for unknown zone, got %v", resUnknown)
	} else if err.ExceptionType != excNames.DateTimeException {
		t.Errorf("Expected DateTimeException, got %v", err.ExceptionType)
	}
}

func TestZoneRulesProviderGetVersions(t *testing.T) {
	zoneRulesProviderTestInit()

	// Null parameter -> NPE
	resNull := zoneRulesProviderGetVersions([]any{object.Null})
	if err, isErr := resNull.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("Expected NPE for null zoneId, got %v", resNull)
	} else if err.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NPE, got %v", err.ExceptionType)
	}

	// Valid zone ID -> returns NavigableMap
	res := zoneRulesProviderGetVersions([]any{object.StringObjectFromGoString("Europe/London")})
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("Error getting versions for Europe/London: %v", err)
	}
	mapObj, ok := res.(*object.Object)
	if !ok || mapObj == nil {
		t.Fatalf("Expected TreeMap object, got %T", res)
	}
	fld, ok := mapObj.FieldTable["map"]
	if !ok {
		t.Fatalf("TreeMap object has no 'map' field")
	}
	hm, ok := fld.Fvalue.(types.DefHashMap)
	if !ok || len(hm) == 0 {
		t.Fatalf("TreeMap 'map' field is empty or not DefHashMap")
	}
	if _, hasTZDB := hm["TZDB"]; !hasTZDB {
		t.Errorf("Expected versions map to contain 'TZDB' key")
	}

	// Unknown zone ID -> DateTimeException
	resUnknown := zoneRulesProviderGetVersions([]any{object.StringObjectFromGoString("Unknown/Invalid_Zone")})
	if err, isErr := resUnknown.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("Expected DateTimeException for unknown zone, got %v", resUnknown)
	} else if err.ExceptionType != excNames.DateTimeException {
		t.Errorf("Expected DateTimeException, got %v", err.ExceptionType)
	}
}

func TestZoneRulesProviderRegisterProvider(t *testing.T) {
	zoneRulesProviderTestInit()

	// Null provider -> NPE
	resNull := zoneRulesProviderRegisterProvider([]any{object.Null})
	if err, isErr := resNull.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("Expected NPE for null provider, got %v", resNull)
	} else if err.ExceptionType != excNames.NullPointerException {
		t.Errorf("Expected NPE, got %v", err.ExceptionType)
	}

	// Register valid custom provider with new zones
	customProvider := object.MakeEmptyObjectWithClassName(&classNameZoneRulesProvider)
	customProvider.FieldTable["zoneIds"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: []string{"Custom/Mars_North", "Custom/Mars_South"},
	}
	resReg := zoneRulesProviderRegisterProvider([]any{customProvider})
	if err, isErr := resReg.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("Failed to register custom provider: %v", err)
	}

	// Confirm new zones are now available
	if !isZoneIdAvailableInternal("Custom/Mars_North") {
		t.Errorf("Expected Custom/Mars_North to be available after registration")
	}

	// Register conflicting provider with existing zone -> DateTimeException
	conflictProvider := object.MakeEmptyObjectWithClassName(&classNameZoneRulesProvider)
	conflictProvider.FieldTable["zoneIds"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: []string{"Custom/Mars_North"},
	}
	resConflict := zoneRulesProviderRegisterProvider([]any{conflictProvider})
	if err, isErr := resConflict.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("Expected exception for conflicting zone registration, got %v", resConflict)
	} else if err.ExceptionType != excNames.DateTimeException {
		t.Errorf("Expected DateTimeException, got %v", err.ExceptionType)
	}
}

func TestZoneRulesProviderRefresh(t *testing.T) {
	zoneRulesProviderTestInit()

	res := zoneRulesProviderRefresh([]any{})
	if res != types.JavaBoolFalse {
		t.Errorf("Expected refresh() to return false, got %v", res)
	}
}

func TestZoneRulesProviderInstanceMethods(t *testing.T) {
	zoneRulesProviderTestInit()

	self := object.MakeEmptyObjectWithClassName(&classNameTzdbZoneRulesProvider)

	// provideZoneIds
	resIds := zoneRulesProviderProvideZoneIds([]any{self})
	if resIds == nil {
		t.Fatalf("provideZoneIds returned nil")
	}

	// provideRules
	resRules := zoneRulesProviderProvideRules([]any{self, object.StringObjectFromGoString("UTC"), types.JavaBoolFalse})
	if err, isErr := resRules.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("provideRules error: %v", err)
	}

	// provideVersions
	resVersions := zoneRulesProviderProvideVersions([]any{self, object.StringObjectFromGoString("UTC")})
	if err, isErr := resVersions.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("provideVersions error: %v", err)
	}

	// provideRefresh
	resRefresh := zoneRulesProviderProvideRefresh([]any{self})
	if resRefresh != types.JavaBoolFalse {
		t.Errorf("provideRefresh expected false, got %v", resRefresh)
	}

	// toString
	resStr := tzdbZoneRulesProviderToString([]any{self})
	strObj, ok := resStr.(*object.Object)
	if !ok || object.GoStringFromStringObject(strObj) != "TZDB" {
		t.Errorf("Expected 'TZDB', got %v", resStr)
	}
}
