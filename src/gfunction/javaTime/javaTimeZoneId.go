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
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	classNameZoneIdVal     = "java/time/ZoneId"
	classNameZoneId        = "java/time/ZoneId"
	classNameZoneRegionVal = "java/time/ZoneRegion"
	classNameZoneRegion    = "java/time/ZoneRegion"
	classNameHashMapZone   = "java/util/HashMap"
)

var (
	zoneIdMutex       sync.Mutex
	zoneIdInitialized bool = false
	zoneIdShortIDsMap      = map[string]string{
		"ACT": "Australia/Darwin",
		"AET": "Australia/Sydney",
		"AGT": "America/Argentina/Buenos_Aires",
		"ART": "Africa/Cairo",
		"AST": "America/Anchorage",
		"BET": "America/Sao_Paulo",
		"BST": "Asia/Dhaka",
		"CAT": "Africa/Harare",
		"CNT": "America/St_Johns",
		"CST": "America/Chicago",
		"CTT": "Asia/Shanghai",
		"EAT": "Africa/Addis_Ababa",
		"ECT": "Europe/Paris",
		"IET": "America/Indiana/Indianapolis",
		"IST": "Asia/Kolkata",
		"JST": "Asia/Tokyo",
		"MIT": "Pacific/Apia",
		"NET": "Asia/Yerevan",
		"NST": "Pacific/Auckland",
		"PLT": "Asia/Karachi",
		"PNT": "America/Phoenix",
		"PRT": "America/Puerto_Rico",
		"PST": "America/Los_Angeles",
		"SST": "Pacific/Guadalcanal",
		"VST": "Asia/Ho_Chi_Minh",
		"EST": "-05:00",
		"MST": "-07:00",
		"HST": "-10:00",
	}
	zoneRegionPattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9~/._+-]*$`)
)

func ensureZoneIdInited() {
	zoneIdMutex.Lock()
	defer zoneIdMutex.Unlock()
	if zoneIdInitialized {
		return
	}

	mapObj := object.MakeEmptyObjectWithClassName(&classNameHashMapZone)
	defMap := make(types.DefHashMap)
	for k, v := range zoneIdShortIDsMap {
		defMap[k] = object.StringObjectFromGoString(v)
	}
	mapObj.FieldTable["map"] = object.Field{Ftype: types.Ref, Fvalue: defMap}

	_ = statics.AddStatic(classNameZoneIdVal+".SHORT_IDS", statics.Static{
		Type:  "Ljava/util/Map;",
		Value: mapObj,
	})

	zoneIdInitialized = true
}

// Load_Time_ZoneId registers all java/time/ZoneId and java/time/ZoneRegion methods.
func Load_Time_ZoneId() {

	// java/time/ZoneId
	ghelpers.MethodSignatures["java/time/ZoneId.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneIdClinit,
		}

	ghelpers.MethodSignatures["java/time/ZoneId.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneIdInit,
		}

	ghelpers.MethodSignatures["java/time/ZoneId.systemDefault()Ljava/time/ZoneId;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneIdSystemDefault,
		}

	ghelpers.MethodSignatures["java/time/ZoneId.getAvailableZoneIds()Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneIdGetAvailableZoneIds,
		}

	ghelpers.MethodSignatures["java/time/ZoneId.of(Ljava/lang/String;)Ljava/time/ZoneId;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneIdOf,
		}

	ghelpers.MethodSignatures["java/time/ZoneId.of(Ljava/lang/String;Ljava/util/Map;)Ljava/time/ZoneId;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zoneIdOfWithMap,
		}

	ghelpers.MethodSignatures["java/time/ZoneId.ofOffset(Ljava/lang/String;Ljava/time/ZoneOffset;)Ljava/time/ZoneId;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zoneIdOfOffset,
		}

	ghelpers.MethodSignatures["java/time/ZoneId.from(Ljava/time/temporal/TemporalAccessor;)Ljava/time/ZoneId;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneIdFrom,
		}

	ghelpers.MethodSignatures["java/time/ZoneId.getId()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneIdGetId,
		}

	ghelpers.MethodSignatures["java/time/ZoneId.getDisplayName(Ljava/time/format/TextStyle;Ljava/util/Locale;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zoneIdGetDisplayName,
		}

	ghelpers.MethodSignatures["java/time/ZoneId.getRules()Ljava/time/zone/ZoneRules;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneIdGetRules,
		}

	ghelpers.MethodSignatures["java/time/ZoneId.normalized()Ljava/time/ZoneId;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneIdNormalized,
		}

	ghelpers.MethodSignatures["java/time/ZoneId.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneIdEquals,
		}

	ghelpers.MethodSignatures["java/time/ZoneId.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneIdHashCode,
		}

	ghelpers.MethodSignatures["java/time/ZoneId.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneIdToString,
		}

	// java/time/ZoneRegion
	Load_Time_ZoneRegion()
}

func zoneIdClinit([]any) any {
	ensureZoneIdInited()
	return nil
}

func zoneIdInit(params []any) any {
	return nil
}

func createZoneId(id string) *object.Object {
	return makeZoneIdFromID(id)
}

func createZoneRegion(id string, rules *object.Object) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameZoneRegionVal)
	idObj := object.StringObjectFromGoString(id)
	obj.FieldTable["id"] = object.Field{Ftype: types.Ref, Fvalue: idObj}
	if rules != nil {
		obj.FieldTable["rules"] = object.Field{Ftype: types.Ref, Fvalue: rules}
	}
	obj.FieldTable["zoneId"] = object.Field{Ftype: types.Ref, Fvalue: obj}
	obj.FieldTable["zone"] = object.Field{Ftype: types.Ref, Fvalue: obj}
	return obj
}

func makeZoneIdFromID(id string) *object.Object {
	res := zoneIdOfInternal(id, true)
	if err, isErr := res.(*ghelpers.GErrBlk); isErr {
		_ = err
		// fallback to simple region object
		return createZoneRegion(id, createZoneRulesForZoneId(id))
	}
	return res.(*object.Object)
}

func zoneIdSystemDefault([]any) any {
	tz := globals.GetSystemProperty("user.timezone")
	if tz != "" {
		res := zoneIdOfInternal(tz, false)
		if _, isErr := res.(*ghelpers.GErrBlk); !isErr {
			return res
		}
	}

	hostTz := time.Now().Location().String()
	if hostTz != "" && hostTz != "Local" {
		res := zoneIdOfInternal(hostTz, false)
		if _, isErr := res.(*ghelpers.GErrBlk); !isErr {
			return res
		}
	}

	return createZoneOffset(0)
}

func zoneIdGetAvailableZoneIds([]any) any {
	return zoneRulesProviderGetAvailableZoneIds(nil)
}

func zoneIdOf(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneId.of: zoneId is null")
	}
	zoneIdObj := params[0].(*object.Object)
	zoneIdStr := object.GoStringFromStringObject(zoneIdObj)
	return zoneIdOfInternal(zoneIdStr, true)
}

func zoneIdOfWithMap(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneId.of: zoneId is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneId.of: aliasMap is null")
	}

	zoneIdObj := params[0].(*object.Object)
	zoneIdStr := object.GoStringFromStringObject(zoneIdObj)
	mapObj := params[1].(*object.Object)

	if mapObj.FieldTable != nil {
		if mapFld, exists := mapObj.FieldTable["map"]; exists {
			if defMap, isDefMap := mapFld.Fvalue.(types.DefHashMap); isDefMap {
				if mappedVal, found := defMap[zoneIdStr]; found && mappedVal != nil {
					if mappedObj, isObj := mappedVal.(*object.Object); isObj && object.IsStringObject(mappedObj) {
						zoneIdStr = object.GoStringFromStringObject(mappedObj)
					} else if s, isStr := mappedVal.(string); isStr {
						zoneIdStr = s
					}
				}
			}
		}
	}

	return zoneIdOfInternal(zoneIdStr, true)
}

func zoneIdOfInternal(zoneIdStr string, checkAvailable bool) any {
	if zoneIdStr == "" {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Invalid ID for ZoneId: empty string")
	}

	if zoneIdStr == "Z" {
		return createZoneOffset(0)
	}

	if len(zoneIdStr) == 1 {
		c := zoneIdStr[0]
		if c == '+' || c == '-' {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "Invalid ID for ZoneOffset: "+zoneIdStr)
		}
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Invalid ID for ZoneId: "+zoneIdStr)
	}

	if strings.HasPrefix(zoneIdStr, "+") || strings.HasPrefix(zoneIdStr, "-") {
		offsetRes := zoneOffsetOf([]any{zoneIdStr})
		return offsetRes
	}

	if zoneIdStr == "UTC" || zoneIdStr == "GMT" || zoneIdStr == "UT" {
		return createZoneRegion(zoneIdStr, createZoneRulesForOffset(createZoneOffset(0)))
	}

	for _, prefix := range []string{"UTC+", "UTC-", "GMT+", "GMT-", "UT+", "UT-"} {
		if strings.HasPrefix(zoneIdStr, prefix) {
			pfx := prefix[:len(prefix)-1]
			offsetStr := zoneIdStr[len(pfx):]
			offsetRes := zoneOffsetOf([]any{offsetStr})
			if err, isErr := offsetRes.(*ghelpers.GErrBlk); isErr {
				return err
			}
			offsetObj := offsetRes.(*object.Object)
			totalSec := getZoneOffsetTotalSeconds(offsetObj)
			var idStr string
			if totalSec == 0 {
				idStr = pfx
			} else {
				idStr = pfx + getZoneOffsetID(offsetObj)
			}
			return createZoneRegion(idStr, createZoneRulesForOffset(offsetObj))
		}
	}

	// Region-based zone ID
	if !zoneRegionPattern.MatchString(zoneIdStr) || len(zoneIdStr) < 2 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Invalid ID for region-based ZoneId, invalid format: "+zoneIdStr)
	}

	if checkAvailable {
		if !isZoneIdAvailableInternal(zoneIdStr) {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unknown time-zone ID: "+zoneIdStr)
		}
	}

	return createZoneRegion(zoneIdStr, createZoneRulesForZoneId(zoneIdStr))
}

func zoneIdOfOffset(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneId.ofOffset: prefix is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneId.ofOffset: offset is null")
	}

	prefixObj := params[0].(*object.Object)
	offsetObj := params[1].(*object.Object)
	prefixStr := object.GoStringFromStringObject(prefixObj)

	if prefixStr == "" {
		return offsetObj
	}

	if prefixStr != "GMT" && prefixStr != "UTC" && prefixStr != "UT" {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneId.ofOffset: prefix should be GMT, UTC, or UT, or empty, but was: "+prefixStr)
	}

	totalSec := getZoneOffsetTotalSeconds(offsetObj)
	if totalSec == 0 {
		return createZoneRegion(prefixStr, createZoneRulesForOffset(offsetObj))
	}

	idStr := prefixStr + getZoneOffsetID(offsetObj)
	return createZoneRegion(idStr, createZoneRulesForOffset(offsetObj))
}

func zoneIdFrom(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneId.from: temporal is null")
	}

	temporalObj := params[0].(*object.Object)
	if isZoneIdObject(temporalObj) {
		return temporalObj
	}

	ensureTemporalQueriesInited()

	zoneIdRes := executeTemporalQuery(zoneIdQueryInstance, temporalObj)
	if zoneIdObj, ok := zoneIdRes.(*object.Object); ok && zoneIdObj != nil && !object.IsNull(zoneIdObj) {
		return zoneIdObj
	}

	zoneRes := executeTemporalQuery(zoneQueryInstance, temporalObj)
	if zoneObj, ok := zoneRes.(*object.Object); ok && zoneObj != nil && !object.IsNull(zoneObj) {
		return zoneObj
	}

	return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unable to obtain ZoneId from TemporalAccessor: "+object.GoStringFromStringPoolIndex(temporalObj.KlassName))
}

func isZoneIdObject(obj *object.Object) bool {
	if obj == nil || obj.FieldTable == nil {
		return false
	}
	_, hasId := obj.FieldTable["id"]
	_, hasTotalSec := obj.FieldTable["totalSeconds"]
	return hasId || hasTotalSec
}

func getZoneIdString(obj *object.Object) string {
	if obj == nil || obj.FieldTable == nil {
		return ""
	}
	if fld, ok := obj.FieldTable["id"]; ok {
		if strObj, okStr := fld.Fvalue.(*object.Object); okStr && object.IsStringObject(strObj) {
			return object.GoStringFromStringObject(strObj)
		}
		if s, okS := fld.Fvalue.(string); okS {
			return s
		}
	}
	if fld, ok := obj.FieldTable["totalSeconds"]; ok {
		if sec, okSec := getInt64Val(fld.Fvalue); okSec {
			return formatZoneOffsetID(sec)
		}
	}
	return ""
}

func zoneIdGetId(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneId.getId: self is null")
	}
	self := params[0].(*object.Object)
	idStr := getZoneIdString(self)
	return object.StringObjectFromGoString(idStr)
}

func zoneIdGetDisplayName(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneId.getDisplayName: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneId.getDisplayName: style is null")
	}
	if len(params) < 3 || params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneId.getDisplayName: locale is null")
	}
	self := params[0].(*object.Object)
	idStr := getZoneIdString(self)
	return object.StringObjectFromGoString(idStr)
}

func zoneIdGetRules(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneId.getRules: self is null")
	}
	self := params[0].(*object.Object)
	if fld, ok := self.FieldTable["rules"]; ok {
		if rulesObj, okRules := fld.Fvalue.(*object.Object); okRules && rulesObj != nil && !object.IsNull(rulesObj) {
			return rulesObj
		}
	}
	if _, isOffset := self.FieldTable["totalSeconds"]; isOffset {
		return createZoneRulesForOffset(self)
	}

	idStr := getZoneIdString(self)
	if idStr == "UTC" || idStr == "GMT" || idStr == "UT" {
		return createZoneRulesForOffset(createZoneOffset(0))
	}
	for _, prefix := range []string{"UTC+", "UTC-", "GMT+", "GMT-", "UT+", "UT-"} {
		if strings.HasPrefix(idStr, prefix) {
			pfx := prefix[:len(prefix)-1]
			offsetStr := idStr[len(pfx):]
			offsetRes := zoneOffsetOf([]any{offsetStr})
			if err, isErr := offsetRes.(*ghelpers.GErrBlk); isErr {
				return err
			}
			return createZoneRulesForOffset(offsetRes.(*object.Object))
		}
	}

	return createZoneRulesForZoneId(idStr)
}

func zoneIdNormalized(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneId.normalized: self is null")
	}
	self := params[0].(*object.Object)
	if _, isOffset := self.FieldTable["totalSeconds"]; isOffset {
		return self
	}

	rulesRes := zoneIdGetRules([]any{self})
	if err, isErr := rulesRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	rulesObj := rulesRes.(*object.Object)

	if isZoneRulesFixed(rulesObj) {
		return getZoneRulesOffsetObject(rulesObj)
	}

	return self
}

func zoneIdEquals(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneId.equals: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)
	if self == other {
		return types.JavaBoolTrue
	}
	if !isZoneIdObject(other) {
		return types.JavaBoolFalse
	}
	id1 := getZoneIdString(self)
	id2 := getZoneIdString(other)
	if id1 == id2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func zoneIdHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneId.hashCode: self is null")
	}
	self := params[0].(*object.Object)
	idStr := getZoneIdString(self)
	var h int32 = 0
	for i := 0; i < len(idStr); i++ {
		h = 31*h + int32(idStr[i])
	}
	return int64(h)
}

func zoneIdToString(params []any) any {
	return zoneIdGetId(params)
}
