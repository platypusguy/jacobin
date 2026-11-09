/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"strings"
	"sync"
)

var (
	classNameZoneOffsetVal = "java/time/ZoneOffset"
	classNameZoneOffset    = "java/time/ZoneOffset"
	classNameZoneRulesVal  = "java/time/zone/ZoneRules"
)

const (
	zoneOffsetMinSeconds = -64800 // -18:00
	zoneOffsetMaxSeconds = 64800  // +18:00
)

var (
	zoneOffsetMutex       sync.Mutex
	zoneOffsetInitialized bool = false
	zoneOffsetUTC         *object.Object
	zoneOffsetMin         *object.Object
	zoneOffsetMax         *object.Object
	zoneOffsetCache       map[int64]*object.Object
)

func ensureZoneOffsetInited() {
	zoneOffsetMutex.Lock()
	defer zoneOffsetMutex.Unlock()
	if zoneOffsetInitialized {
		return
	}
	zoneOffsetCache = make(map[int64]*object.Object)

	zoneOffsetUTC = makeZoneOffsetInternal(0, "Z")
	zoneOffsetMin = makeZoneOffsetInternal(zoneOffsetMinSeconds, "-18:00")
	zoneOffsetMax = makeZoneOffsetInternal(zoneOffsetMaxSeconds, "+18:00")

	zoneOffsetCache[0] = zoneOffsetUTC
	zoneOffsetCache[zoneOffsetMinSeconds] = zoneOffsetMin
	zoneOffsetCache[zoneOffsetMaxSeconds] = zoneOffsetMax

	_ = statics.AddStatic(classNameZoneOffsetVal+".UTC", statics.Static{
		Type:  "Ljava/time/ZoneOffset;",
		Value: zoneOffsetUTC,
	})
	_ = statics.AddStatic(classNameZoneOffsetVal+".MIN", statics.Static{
		Type:  "Ljava/time/ZoneOffset;",
		Value: zoneOffsetMin,
	})
	_ = statics.AddStatic(classNameZoneOffsetVal+".MAX", statics.Static{
		Type:  "Ljava/time/ZoneOffset;",
		Value: zoneOffsetMax,
	})

	zoneOffsetInitialized = true
}

func makeZoneOffsetInternal(totalSeconds int64, id string) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&classNameZoneOffsetVal)
	idObj := object.StringObjectFromGoString(id)
	obj.FieldTable["totalSeconds"] = object.Field{Ftype: types.Int, Fvalue: totalSeconds}
	obj.FieldTable["id"] = object.Field{Ftype: types.Ref, Fvalue: idObj}
	obj.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: obj}
	obj.FieldTable["zone"] = object.Field{Ftype: types.Ref, Fvalue: obj}
	obj.FieldTable["zoneId"] = object.Field{Ftype: types.Ref, Fvalue: obj}
	return obj
}

func createZoneOffset(totalSeconds int64) *object.Object {
	ensureZoneOffsetInited()
	if totalSeconds < zoneOffsetMinSeconds || totalSeconds > zoneOffsetMaxSeconds {
		return nil
	}
	zoneOffsetMutex.Lock()
	defer zoneOffsetMutex.Unlock()
	if cached, ok := zoneOffsetCache[totalSeconds]; ok {
		return cached
	}
	idStr := formatZoneOffsetID(totalSeconds)
	obj := makeZoneOffsetInternal(totalSeconds, idStr)
	zoneOffsetCache[totalSeconds] = obj
	return obj
}

func formatZoneOffsetID(totalSeconds int64) string {
	if totalSeconds == 0 {
		return "Z"
	}
	absSec := totalSeconds
	sign := "+"
	if absSec < 0 {
		sign = "-"
		absSec = -absSec
	}
	hours := absSec / 3600
	minutes := (absSec % 3600) / 60
	secs := absSec % 60
	if secs == 0 {
		return fmt.Sprintf("%s%02d:%02d", sign, hours, minutes)
	}
	return fmt.Sprintf("%s%02d:%02d:%02d", sign, hours, minutes, secs)
}

func getZoneOffsetTotalSecondsVal(obj *object.Object) (int64, bool) {
	if obj == nil || obj.FieldTable == nil {
		return 0, false
	}
	if fld, ok := obj.FieldTable["totalSeconds"]; ok {
		if val, okVal := getInt64Val(fld.Fvalue); okVal {
			return val, true
		}
	}
	if fld, ok := obj.FieldTable["seconds"]; ok {
		if val, okVal := getInt64Val(fld.Fvalue); okVal {
			return val, true
		}
	}
	return 0, false
}

func getZoneOffsetID(obj *object.Object) string {
	if obj == nil || obj.FieldTable == nil {
		return "Z"
	}
	if fld, ok := obj.FieldTable["id"]; ok {
		if strObj, okStr := fld.Fvalue.(*object.Object); okStr && object.IsStringObject(strObj) {
			return object.GoStringFromStringObject(strObj)
		}
		if s, okS := fld.Fvalue.(string); okS {
			return s
		}
	}
	if sec, ok := getZoneOffsetTotalSecondsVal(obj); ok {
		return formatZoneOffsetID(sec)
	}
	return "Z"
}

func isZoneOffsetObject(obj *object.Object) bool {
	if obj == nil || obj.FieldTable == nil {
		return false
	}
	klass := object.GoStringFromStringPoolIndex(obj.KlassName)
	if klass == classNameZoneOffsetVal {
		return true
	}
	_, hasSeconds := obj.FieldTable["totalSeconds"]
	return hasSeconds
}

func getChronoFieldForZoneOffset(name string) *object.Object {
	ensureChronoFieldInited()
	for index, definition := range chronoFieldDefinitions {
		if definition.enumName == name {
			return chronoFieldInstances[index]
		}
	}
	return nil
}

// Load_Time_ZoneOffset registers all java/time/ZoneOffset methods.
func Load_Time_ZoneOffset() {
	ghelpers.MethodSignatures["java/time/ZoneOffset.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetClinit,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetInit,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.<init>(I)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetInit,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.adjustInto(Ljava/time/temporal/Temporal;)Ljava/time/temporal/Temporal;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetAdjustInto,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetCompareTo,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.compareTo(Ljava/time/ZoneOffset;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetCompareTo,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetEquals,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.from(Ljava/time/temporal/TemporalAccessor;)Ljava/time/ZoneOffset;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetFrom,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.get(Ljava/time/temporal/TemporalField;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetGet,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.getId()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetGetId,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.getLong(Ljava/time/temporal/TemporalField;)J"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetGetLong,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.getRules()Ljava/time/zone/ZoneRules;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetGetRules,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.getTotalSeconds()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetGetTotalSeconds,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetHashCode,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.isSupported(Ljava/time/temporal/TemporalField;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetIsSupported,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.normalized()Ljava/time/ZoneId;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetNormalized,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.normalized()Ljava/time/ZoneOffset;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetNormalized,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.of(Ljava/lang/String;)Ljava/time/ZoneOffset;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetOf,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.ofHours(I)Ljava/time/ZoneOffset;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetOfHours,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.ofHoursMinutes(II)Ljava/time/ZoneOffset;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zoneOffsetOfHoursMinutes,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.ofHoursMinutesSeconds(III)Ljava/time/ZoneOffset;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  zoneOffsetOfHoursMinutesSeconds,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.ofTotalSeconds(I)Ljava/time/ZoneOffset;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetOfTotalSeconds,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.query(Ljava/time/temporal/TemporalQuery;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetQuery,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.range(Ljava/time/temporal/TemporalField;)Ljava/time/temporal/ValueRange;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneOffsetRange,
		}

	ghelpers.MethodSignatures["java/time/ZoneOffset.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneOffsetToString,
		}
}

func zoneOffsetClinit([]any) any {
	ensureZoneOffsetInited()
	return nil
}

func zoneOffsetInit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.<init>: self is null")
	}
	self := params[0].(*object.Object)
	var totalSeconds int64 = 0
	if len(params) > 1 && params[1] != nil {
		if sec, ok := getInt64Val(params[1]); ok {
			totalSeconds = sec
		}
	}
	idStr := formatZoneOffsetID(totalSeconds)
	idObj := object.StringObjectFromGoString(idStr)
	self.FieldTable["totalSeconds"] = object.Field{Ftype: types.Int, Fvalue: totalSeconds}
	self.FieldTable["id"] = object.Field{Ftype: types.Ref, Fvalue: idObj}
	self.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: self}
	self.FieldTable["zone"] = object.Field{Ftype: types.Ref, Fvalue: self}
	self.FieldTable["zoneId"] = object.Field{Ftype: types.Ref, Fvalue: self}
	return nil
}

func zoneOffsetOf(params []any) any {
	if len(params) < 1 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.of: offsetId is null")
	}
	if obj, ok := params[0].(*object.Object); ok && object.IsNull(obj) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.of: offsetId is null")
	}
	ensureZoneOffsetInited()
	var offsetId string
	switch v := params[0].(type) {
	case *object.Object:
		if object.IsStringObject(v) {
			offsetId = object.GoStringFromStringObject(v)
		} else {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneOffset.of: expected String")
		}
	case string:
		offsetId = v
	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneOffset.of: invalid argument type")
	}

	if offsetId == "Z" {
		return zoneOffsetUTC
	}

	var hours, minutes, seconds int64
	length := len(offsetId)
	if length == 2 {
		offsetId = string(offsetId[0]) + "0" + string(offsetId[1])
		length = 3
	}

	switch length {
	case 3:
		h, err := parseZoneOffsetNumber(offsetId, 1, false)
		if err != nil {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
		}
		hours = h
	case 5:
		h, err := parseZoneOffsetNumber(offsetId, 1, false)
		if err != nil {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
		}
		m, err := parseZoneOffsetNumber(offsetId, 3, false)
		if err != nil {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
		}
		hours, minutes = h, m
	case 6:
		h, err := parseZoneOffsetNumber(offsetId, 1, false)
		if err != nil {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
		}
		m, err := parseZoneOffsetNumber(offsetId, 4, true)
		if err != nil {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
		}
		hours, minutes = h, m
	case 7:
		h, err := parseZoneOffsetNumber(offsetId, 1, false)
		if err != nil {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
		}
		m, err := parseZoneOffsetNumber(offsetId, 3, false)
		if err != nil {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
		}
		s, err := parseZoneOffsetNumber(offsetId, 5, false)
		if err != nil {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
		}
		hours, minutes, seconds = h, m, s
	case 9:
		h, err := parseZoneOffsetNumber(offsetId, 1, false)
		if err != nil {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
		}
		m, err := parseZoneOffsetNumber(offsetId, 4, true)
		if err != nil {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
		}
		s, err := parseZoneOffsetNumber(offsetId, 7, true)
		if err != nil {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, err.Error())
		}
		hours, minutes, seconds = h, m, s
	default:
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Invalid ID for ZoneOffset, invalid format: "+offsetId)
	}

	first := offsetId[0]
	if first != '+' && first != '-' {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Invalid ID for ZoneOffset, plus/minus not found when expected: "+offsetId)
	}

	if first == '-' {
		return zoneOffsetOfHoursMinutesSeconds([]any{-hours, -minutes, -seconds})
	}
	return zoneOffsetOfHoursMinutesSeconds([]any{hours, minutes, seconds})
}

func parseZoneOffsetNumber(offsetId string, pos int, precededByColon bool) (int64, error) {
	if precededByColon && offsetId[pos-1] != ':' {
		return 0, fmt.Errorf("Invalid ID for ZoneOffset, colon not found when expected: %s", offsetId)
	}
	ch1 := offsetId[pos]
	ch2 := offsetId[pos+1]
	if ch1 < '0' || ch1 > '9' || ch2 < '0' || ch2 > '9' {
		return 0, fmt.Errorf("Invalid ID for ZoneOffset, non numeric characters found: %s", offsetId)
	}
	return int64((ch1-'0')*10 + (ch2 - '0')), nil
}

func zoneOffsetOfHours(params []any) any {
	if len(params) < 1 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneOffset.ofHours: missing argument")
	}
	hours, ok := getInt64Val(params[0])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneOffset.ofHours: invalid hours")
	}
	return zoneOffsetOfHoursMinutesSeconds([]any{hours, int64(0), int64(0)})
}

func zoneOffsetOfHoursMinutes(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneOffset.ofHoursMinutes: missing argument")
	}
	hours, okH := getInt64Val(params[0])
	minutes, okM := getInt64Val(params[1])
	if !okH || !okM {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneOffset.ofHoursMinutes: invalid argument")
	}
	return zoneOffsetOfHoursMinutesSeconds([]any{hours, minutes, int64(0)})
}

func zoneOffsetOfHoursMinutesSeconds(params []any) any {
	if len(params) < 3 || params[0] == nil || params[1] == nil || params[2] == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneOffset.ofHoursMinutesSeconds: missing argument")
	}
	hours, okH := getInt64Val(params[0])
	minutes, okM := getInt64Val(params[1])
	seconds, okS := getInt64Val(params[2])
	if !okH || !okM || !okS {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneOffset.ofHoursMinutesSeconds: invalid argument")
	}

	if hours < -18 || hours > 18 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Zone offset hours not in valid range: value %d is not in the range -18 to 18", hours))
	}
	if hours > 0 {
		if minutes < 0 || seconds < 0 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "Zone offset minutes and seconds must be positive because hours is positive")
		}
	} else if hours < 0 {
		if minutes > 0 || seconds > 0 {
			return ghelpers.GetGErrBlk(excNames.DateTimeException, "Zone offset minutes and seconds must be negative because hours is negative")
		}
	} else if (minutes > 0 && seconds < 0) || (minutes < 0 && seconds > 0) {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Zone offset minutes and seconds must have the same sign")
	}

	absMinutes := minutes
	if absMinutes < 0 {
		absMinutes = -absMinutes
	}
	if absMinutes > 59 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Zone offset minutes not in valid range: value %d is not in the range 0 to 59", absMinutes))
	}

	absSeconds := seconds
	if absSeconds < 0 {
		absSeconds = -absSeconds
	}
	if absSeconds > 59 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Zone offset seconds not in valid range: value %d is not in the range 0 to 59", absSeconds))
	}

	absHours := hours
	if absHours < 0 {
		absHours = -absHours
	}
	if absHours == 18 && (absMinutes > 0 || absSeconds > 0) {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Zone offset not in valid range: -18:00 to +18:00")
	}

	totalSec := hours*3600 + minutes*60 + seconds
	return zoneOffsetOfTotalSeconds([]any{totalSec})
}

func zoneOffsetOfTotalSeconds(params []any) any {
	if len(params) < 1 || params[0] == nil {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneOffset.ofTotalSeconds: missing argument")
	}
	totalSeconds, ok := getInt64Val(params[0])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneOffset.ofTotalSeconds: invalid totalSeconds")
	}
	if totalSeconds < zoneOffsetMinSeconds || totalSeconds > zoneOffsetMaxSeconds {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Zone offset not in valid range: -18:00 to +18:00")
	}
	return createZoneOffset(totalSeconds)
}

func zoneOffsetFrom(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.from: temporal is null")
	}
	temporal := params[0].(*object.Object)
	if isZoneOffsetObject(temporal) {
		return temporal
	}

	ensureTemporalQueriesInited()
	res := executeTemporalQuery(offsetQueryInstance, temporal)
	if resObj, ok := res.(*object.Object); ok && !object.IsNull(resObj) && isZoneOffsetObject(resObj) {
		return resObj
	}

	if temporal.FieldTable != nil {
		if offFld, exists := temporal.FieldTable["offset"]; exists {
			if offObj, ok := offFld.Fvalue.(*object.Object); ok && !object.IsNull(offObj) && isZoneOffsetObject(offObj) {
				return offObj
			}
		}
		if offFld, exists := temporal.FieldTable["totalSeconds"]; exists {
			if sec, ok := getInt64Val(offFld.Fvalue); ok {
				return createZoneOffset(sec)
			}
		}
	}
	return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unable to obtain ZoneOffset from TemporalAccessor")
}

func zoneOffsetGet(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.get: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.get: field is null")
	}
	self := params[0].(*object.Object)
	field := params[1].(*object.Object)
	fieldName := getTemporalFieldName(field)
	if strings.ToUpper(fieldName) == "OFFSET_SECONDS" {
		sec, _ := getZoneOffsetTotalSecondsVal(self)
		return sec
	}
	if object.GoStringFromStringPoolIndex(field.KlassName) == classNameChronoField {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
	}
	return temporalAccessorGet(params)
}

func zoneOffsetGetLong(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.getLong: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.getLong: field is null")
	}
	self := params[0].(*object.Object)
	field := params[1].(*object.Object)
	fieldName := getTemporalFieldName(field)
	if strings.ToUpper(fieldName) == "OFFSET_SECONDS" {
		sec, _ := getZoneOffsetTotalSecondsVal(self)
		return sec
	}
	if object.GoStringFromStringPoolIndex(field.KlassName) == classNameChronoField {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
	}
	return temporalAccessorGetLong(params)
}

func zoneOffsetIsSupported(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.isSupported: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	field := params[1].(*object.Object)
	fieldName := getTemporalFieldName(field)
	if strings.ToUpper(fieldName) == "OFFSET_SECONDS" {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func zoneOffsetRange(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.range: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.range: field is null")
	}
	field := params[1].(*object.Object)
	fieldName := getTemporalFieldName(field)
	if strings.ToUpper(fieldName) == "OFFSET_SECONDS" {
		return createValueRange(-64800, -64800, 64800, 64800)
	}
	if object.GoStringFromStringPoolIndex(field.KlassName) == classNameChronoField {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, fmt.Sprintf("Unsupported field: %s", fieldName))
	}
	return temporalAccessorRange(params)
}

func zoneOffsetQuery(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.query: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.query: query is null")
	}
	self := params[0].(*object.Object)
	query := params[1].(*object.Object)

	queryName := ""
	if query.FieldTable != nil {
		if nameFld, okName := query.FieldTable["name"]; okName {
			if strObj, okStr := nameFld.Fvalue.(*object.Object); okStr && object.IsStringObject(strObj) {
				queryName = object.GoStringFromStringObject(strObj)
			}
		}
	}

	switch queryName {
	case "offset", "zone":
		return self
	case "zoneId", "chronology", "precision", "localDate", "localTime":
		return nil
	default:
		return executeTemporalQuery(query, self)
	}
}

func zoneOffsetAdjustInto(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.adjustInto: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.adjustInto: temporal is null")
	}
	self := params[0].(*object.Object)
	temporal := params[1].(*object.Object)
	totalSeconds, ok := getZoneOffsetTotalSecondsVal(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "ZoneOffset.adjustInto: invalid state")
	}
	cfOffsetSeconds := getChronoFieldForZoneOffset("OFFSET_SECONDS")
	if cfOffsetSeconds == nil {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "ZoneOffset.adjustInto: OFFSET_SECONDS unavailable")
	}
	return temporalWithField([]any{temporal, cfOffsetSeconds, totalSeconds})
}

func zoneOffsetCompareTo(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.compareTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.compareTo: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)
	s1, ok1 := getZoneOffsetTotalSecondsVal(self)
	s2, ok2 := getZoneOffsetTotalSecondsVal(other)
	if !ok1 || !ok2 {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "ZoneOffset.compareTo: cannot compare")
	}
	if s2 < s1 {
		return int64(-1)
	} else if s2 > s1 {
		return int64(1)
	}
	return int64(0)
}

func zoneOffsetEquals(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return types.JavaBoolFalse
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)
	if self == other {
		return types.JavaBoolTrue
	}
	s1, ok1 := getZoneOffsetTotalSecondsVal(self)
	s2, ok2 := getZoneOffsetTotalSecondsVal(other)
	if ok1 && ok2 && s1 == s2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func zoneOffsetHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.hashCode: self is null")
	}
	self := params[0].(*object.Object)
	sec, _ := getZoneOffsetTotalSecondsVal(self)
	return sec
}

func zoneOffsetToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.toString: self is null")
	}
	self := params[0].(*object.Object)
	idStr := getZoneOffsetID(self)
	return object.StringObjectFromGoString(idStr)
}

func zoneOffsetGetId(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.getId: self is null")
	}
	self := params[0].(*object.Object)
	idStr := getZoneOffsetID(self)
	return object.StringObjectFromGoString(idStr)
}

func zoneOffsetGetRules(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.getRules: self is null")
	}
	self := params[0].(*object.Object)
	return createZoneRulesForOffset(self)
}

func zoneOffsetGetTotalSeconds(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.getTotalSeconds: self is null")
	}
	self := params[0].(*object.Object)
	sec, _ := getZoneOffsetTotalSecondsVal(self)
	return sec
}

func zoneOffsetNormalized(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneOffset.normalized: self is null")
	}
	return params[0].(*object.Object)
}
