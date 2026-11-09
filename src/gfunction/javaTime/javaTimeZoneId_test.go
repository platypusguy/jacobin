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
	"jacobin/src/statics"
	"jacobin/src/types"
	"testing"
)

func zoneIdTestInit() {
	globals.InitGlobals("test")
	Load_Time_ZoneOffset()
	Load_Time_Zone_ZoneRules()
	Load_Time_ZoneId()
}

func TestZoneIdClinitAndStaticFields(t *testing.T) {
	zoneIdTestInit()
	ensureZoneIdInited()

	mapVal := statics.GetStaticValue(classNameZoneIdVal, "SHORT_IDS")
	if mapVal == nil {
		t.Fatalf("ZoneId.SHORT_IDS static field not found")
	}
	mapObj, ok := mapVal.(*object.Object)
	if !ok || mapObj == nil {
		t.Fatalf("ZoneId.SHORT_IDS is not an Object")
	}
	fld, ok := mapObj.FieldTable["map"]
	if !ok {
		t.Fatalf("ZoneId.SHORT_IDS missing map field")
	}
	defMap, ok := fld.Fvalue.(types.DefHashMap)
	if !ok {
		t.Fatalf("ZoneId.SHORT_IDS map field is not DefHashMap")
	}
	if len(defMap) < 25 {
		t.Fatalf("ZoneId.SHORT_IDS expected >= 25 entries, got %d", len(defMap))
	}

	estVal, ok := defMap["EST"]
	if !ok {
		t.Fatalf("ZoneId.SHORT_IDS missing EST")
	}
	estObj, ok := estVal.(*object.Object)
	if !ok || object.GoStringFromStringObject(estObj) != "-05:00" {
		t.Fatalf("ZoneId.SHORT_IDS['EST'] = %v, want -05:00", estVal)
	}

	actVal, ok := defMap["ACT"]
	if !ok {
		t.Fatalf("ZoneId.SHORT_IDS missing ACT")
	}
	actObj, ok := actVal.(*object.Object)
	if !ok || object.GoStringFromStringObject(actObj) != "Australia/Darwin" {
		t.Fatalf("ZoneId.SHORT_IDS['ACT'] = %v, want Australia/Darwin", actVal)
	}
}

func TestZoneIdOfOffsets(t *testing.T) {
	zoneIdTestInit()

	// "Z"
	zRes := zoneIdOf([]any{object.StringObjectFromGoString("Z")})
	if err, isErr := zRes.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneId.of('Z') error: %v", err)
	}
	zObj := zRes.(*object.Object)
	if getZoneIdString(zObj) != "Z" {
		t.Fatalf("ZoneId.of('Z').getId() = %v, want Z", getZoneIdString(zObj))
	}

	// "+02:00"
	p2Res := zoneIdOf([]any{object.StringObjectFromGoString("+02:00")})
	if err, isErr := p2Res.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneId.of('+02:00') error: %v", err)
	}
	p2Obj := p2Res.(*object.Object)
	if getZoneIdString(p2Obj) != "+02:00" {
		t.Fatalf("ZoneId.of('+02:00').getId() = %v, want +02:00", getZoneIdString(p2Obj))
	}

	// "UTC"
	utcRes := zoneIdOf([]any{object.StringObjectFromGoString("UTC")})
	if err, isErr := utcRes.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneId.of('UTC') error: %v", err)
	}
	utcObj := utcRes.(*object.Object)
	if getZoneIdString(utcObj) != "UTC" {
		t.Fatalf("ZoneId.of('UTC').getId() = %v, want UTC", getZoneIdString(utcObj))
	}

	// "GMT+05:30"
	gmtRes := zoneIdOf([]any{object.StringObjectFromGoString("GMT+05:30")})
	if err, isErr := gmtRes.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneId.of('GMT+05:30') error: %v", err)
	}
	gmtObj := gmtRes.(*object.Object)
	if getZoneIdString(gmtObj) != "GMT+05:30" {
		t.Fatalf("ZoneId.of('GMT+05:30').getId() = %v, want GMT+05:30", getZoneIdString(gmtObj))
	}

	// "UTC+00:00" should normalize ID to "UTC"
	utcZeroRes := zoneIdOf([]any{object.StringObjectFromGoString("UTC+00:00")})
	if err, isErr := utcZeroRes.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneId.of('UTC+00:00') error: %v", err)
	}
	utcZeroObj := utcZeroRes.(*object.Object)
	if getZoneIdString(utcZeroObj) != "UTC" {
		t.Fatalf("ZoneId.of('UTC+00:00').getId() = %v, want UTC", getZoneIdString(utcZeroObj))
	}

	// Invalid format
	invRes := zoneIdOf([]any{object.StringObjectFromGoString("+19:00")})
	if _, isErr := invRes.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("ZoneId.of('+19:00') expected error, got %v", invRes)
	}
}

func TestZoneIdOfRegions(t *testing.T) {
	zoneIdTestInit()

	// "Europe/London"
	londonRes := zoneIdOf([]any{object.StringObjectFromGoString("Europe/London")})
	if err, isErr := londonRes.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneId.of('Europe/London') error: %v", err)
	}
	londonObj := londonRes.(*object.Object)
	if getZoneIdString(londonObj) != "Europe/London" {
		t.Fatalf("ZoneId.of('Europe/London').getId() = %v, want Europe/London", getZoneIdString(londonObj))
	}

	// Invalid characters in region ID
	invRes := zoneIdOf([]any{object.StringObjectFromGoString("Invalid@Zone")})
	if err, isErr := invRes.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("ZoneId.of('Invalid@Zone') expected error, got %v", invRes)
	} else if err.ExceptionType != excNames.DateTimeException {
		t.Fatalf("ZoneId.of('Invalid@Zone') expected DateTimeException, got %v", err.ExceptionType)
	}

	// Unknown region ID
	unknownRes := zoneIdOf([]any{object.StringObjectFromGoString("NonExistent/Zone_12345")})
	if err, isErr := unknownRes.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("ZoneId.of('NonExistent/Zone_12345') expected error, got %v", unknownRes)
	} else if err.ExceptionType != excNames.DateTimeException {
		t.Fatalf("ZoneId.of('NonExistent/Zone_12345') expected DateTimeException, got %v", err.ExceptionType)
	}
}

func TestZoneIdOfWithMap(t *testing.T) {
	zoneIdTestInit()
	ensureZoneIdInited()

	mapVal := statics.GetStaticValue(classNameZoneIdVal, "SHORT_IDS")
	mapObj := mapVal.(*object.Object)

	res := zoneIdOfWithMap([]any{object.StringObjectFromGoString("EST"), mapObj})
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneId.of('EST', SHORT_IDS) error: %v", err)
	}
	estObj := res.(*object.Object)
	if getZoneIdString(estObj) != "-05:00" {
		t.Fatalf("ZoneId.of('EST', SHORT_IDS).getId() = %v, want -05:00", getZoneIdString(estObj))
	}
}

func TestZoneIdOfOffset(t *testing.T) {
	zoneIdTestInit()

	off2 := createZoneOffset(7200) // +02:00
	off0 := createZoneOffset(0)    // Z

	// prefix = "GMT", offset = +02:00
	res := zoneIdOfOffset([]any{object.StringObjectFromGoString("GMT"), off2})
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneId.ofOffset('GMT', +02:00) error: %v", err)
	}
	gmt2Obj := res.(*object.Object)
	if getZoneIdString(gmt2Obj) != "GMT+02:00" {
		t.Fatalf("ZoneId.ofOffset('GMT', +02:00).getId() = %v, want GMT+02:00", getZoneIdString(gmt2Obj))
	}

	// prefix = "", offset = +02:00 -> returns ZoneOffset itself
	resEmpty := zoneIdOfOffset([]any{object.StringObjectFromGoString(""), off2})
	if err, isErr := resEmpty.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneId.ofOffset('', +02:00) error: %v", err)
	}
	if resEmpty.(*object.Object) != off2 {
		t.Fatalf("ZoneId.ofOffset('', +02:00) expected offset object itself")
	}

	// prefix = "UTC", offset = Z -> "UTC"
	resUTC := zoneIdOfOffset([]any{object.StringObjectFromGoString("UTC"), off0})
	if err, isErr := resUTC.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneId.ofOffset('UTC', Z) error: %v", err)
	}
	if getZoneIdString(resUTC.(*object.Object)) != "UTC" {
		t.Fatalf("ZoneId.ofOffset('UTC', Z).getId() = %v, want UTC", getZoneIdString(resUTC.(*object.Object)))
	}

	// prefix = "INVALID" -> IllegalArgumentException
	resInv := zoneIdOfOffset([]any{object.StringObjectFromGoString("INVALID"), off2})
	if err, isErr := resInv.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("ZoneId.ofOffset('INVALID', +02:00) expected error, got %v", resInv)
	} else if err.ExceptionType != excNames.IllegalArgumentException {
		t.Fatalf("ZoneId.ofOffset('INVALID', +02:00) expected IllegalArgumentException, got %v", err.ExceptionType)
	}
}

func TestZoneIdFromTemporal(t *testing.T) {
	zoneIdTestInit()
	Load_Time_Temporal_TemporalQueries()

	zoneLondon := makeZoneIdFromID("Europe/London")
	tempObj := object.MakeEmptyObjectWithClassName(&classNameZoneIdVal)
	tempObj.FieldTable["zoneId"] = object.Field{Ftype: types.Ref, Fvalue: zoneLondon}

	res := zoneIdFrom([]any{tempObj})
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneId.from(tempObj) error: %v", err)
	}
	if res.(*object.Object) != zoneLondon {
		t.Fatalf("ZoneId.from(tempObj) expected zoneLondon object")
	}

	// Missing zone
	emptyTemp := object.MakeEmptyObjectWithClassName(&classNameZoneIdVal)
	resEmpty := zoneIdFrom([]any{emptyTemp})
	if err, isErr := resEmpty.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("ZoneId.from(emptyTemp) expected error, got %v", resEmpty)
	} else if err.ExceptionType != excNames.DateTimeException {
		t.Fatalf("ZoneId.from(emptyTemp) expected DateTimeException, got %v", err.ExceptionType)
	}
}

func TestZoneIdGetRulesAndNormalized(t *testing.T) {
	zoneIdTestInit()

	// "UTC"
	utcZone := makeZoneIdFromID("UTC")
	rulesRes := zoneIdGetRules([]any{utcZone})
	if err, isErr := rulesRes.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneId.getRules('UTC') error: %v", err)
	}
	rulesObj := rulesRes.(*object.Object)
	if !isZoneRulesFixed(rulesObj) {
		t.Fatalf("ZoneId.getRules('UTC') expected fixed rules")
	}

	normRes := zoneIdNormalized([]any{utcZone})
	if err, isErr := normRes.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneId.normalized('UTC') error: %v", err)
	}
	normObj := normRes.(*object.Object)
	if getZoneIdString(normObj) != "Z" {
		t.Fatalf("ZoneId.normalized('UTC').getId() = %v, want Z", getZoneIdString(normObj))
	}

	// "GMT+03:00"
	gmt3 := makeZoneIdFromID("GMT+03:00")
	norm3Res := zoneIdNormalized([]any{gmt3})
	if err, isErr := norm3Res.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneId.normalized('GMT+03:00') error: %v", err)
	}
	norm3Obj := norm3Res.(*object.Object)
	if getZoneIdString(norm3Obj) != "+03:00" {
		t.Fatalf("ZoneId.normalized('GMT+03:00').getId() = %v, want +03:00", getZoneIdString(norm3Obj))
	}

	// "Europe/London" (dynamic zone -> normalized returns self)
	london := makeZoneIdFromID("Europe/London")
	normLondonRes := zoneIdNormalized([]any{london})
	if err, isErr := normLondonRes.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneId.normalized('Europe/London') error: %v", err)
	}
	if normLondonRes.(*object.Object) != london {
		t.Fatalf("ZoneId.normalized('Europe/London') expected self")
	}
}

func TestZoneIdEqualsHashCodeToString(t *testing.T) {
	zoneIdTestInit()

	z1 := makeZoneIdFromID("Europe/London")
	z2 := makeZoneIdFromID("Europe/London")
	z3 := makeZoneIdFromID("UTC")

	eq1 := zoneIdEquals([]any{z1, z2})
	if eq1 != types.JavaBoolTrue {
		t.Fatalf("ZoneId.equals(z1, z2) = %v, want true", eq1)
	}

	eq2 := zoneIdEquals([]any{z1, z3})
	if eq2 != types.JavaBoolFalse {
		t.Fatalf("ZoneId.equals(z1, z3) = %v, want false", eq2)
	}

	eqNil := zoneIdEquals([]any{z1, nil})
	if eqNil != types.JavaBoolFalse {
		t.Fatalf("ZoneId.equals(z1, nil) = %v, want false", eqNil)
	}

	h1 := zoneIdHashCode([]any{z1}).(int64)
	h2 := zoneIdHashCode([]any{z2}).(int64)
	if h1 != h2 {
		t.Fatalf("ZoneId.hashCode(z1) = %d != hashCode(z2) = %d", h1, h2)
	}

	strRes := zoneIdToString([]any{z1}).(*object.Object)
	if object.GoStringFromStringObject(strRes) != "Europe/London" {
		t.Fatalf("ZoneId.toString() = %v, want Europe/London", object.GoStringFromStringObject(strRes))
	}

	dispRes := zoneIdGetDisplayName([]any{z1, object.Null, object.Null})
	if _, isErr := dispRes.(*ghelpers.GErrBlk); !isErr {
		t.Fatalf("ZoneId.getDisplayName(null, null) expected NPE")
	}
}

func TestZoneIdAvailableZoneIdsAndSystemDefault(t *testing.T) {
	zoneIdTestInit()

	availRes := zoneIdGetAvailableZoneIds([]any{})
	if err, isErr := availRes.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneId.getAvailableZoneIds() error: %v", err)
	}
	availObj, ok := availRes.(*object.Object)
	if !ok || availObj == nil {
		t.Fatalf("ZoneId.getAvailableZoneIds() returned non-object")
	}
	mapFld, ok := availObj.FieldTable["map"]
	if !ok {
		t.Fatalf("Set object has no 'map' field")
	}
	hm, ok := mapFld.Fvalue.(types.DefHashMap)
	if !ok {
		t.Fatalf("Set field 'map' is not DefHashMap")
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
		t.Errorf("ZoneId.getAvailableZoneIds() should NOT contain 'UT'")
	}
	if !hasUTC {
		t.Errorf("ZoneId.getAvailableZoneIds() must contain 'UTC'")
	}

	defRes := zoneIdSystemDefault([]any{})
	if err, isErr := defRes.(*ghelpers.GErrBlk); isErr {
		t.Fatalf("ZoneId.systemDefault() error: %v", err)
	}
	defObj, ok := defRes.(*object.Object)
	if !ok || defObj == nil {
		t.Fatalf("ZoneId.systemDefault() returned non-object")
	}
}
