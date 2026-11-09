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
	"jacobin/src/types"
	"time"
)

var (
	classNameZoneRules = "java/time/zone/ZoneRules"
)

// Load_Time_Zone_ZoneRules registers all java/time/zone/ZoneRules methods.
func Load_Time_Zone_ZoneRules() {
	ghelpers.MethodSignatures["java/time/zone/ZoneRules.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRules.of(Ljava/time/ZoneOffset;)Ljava/time/zone/ZoneRules;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRulesOfOffset,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRules.of(Ljava/time/ZoneOffset;Ljava/time/ZoneOffset;Ljava/util/List;Ljava/util/List;Ljava/util/List;)Ljava/time/zone/ZoneRules;"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  zoneRulesOfFull,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRules.isFixedOffset()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneRulesIsFixedOffset,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRules.getOffset(Ljava/time/Instant;)Ljava/time/ZoneOffset;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRulesGetOffsetInstant,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRules.getOffset(Ljava/time/LocalDateTime;)Ljava/time/ZoneOffset;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRulesGetOffsetLocalDateTime,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRules.getValidOffsets(Ljava/time/LocalDateTime;)Ljava/util/List;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRulesGetValidOffsets,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRules.getTransition(Ljava/time/LocalDateTime;)Ljava/time/zone/ZoneOffsetTransition;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRulesGetTransition,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRules.getStandardOffset(Ljava/time/Instant;)Ljava/time/ZoneOffset;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRulesGetStandardOffset,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRules.getDaylightSavings(Ljava/time/Instant;)Ljava/time/Duration;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRulesGetDaylightSavings,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRules.isDaylightSavings(Ljava/time/Instant;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRulesIsDaylightSavings,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRules.isValidOffset(Ljava/time/LocalDateTime;Ljava/time/ZoneOffset;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zoneRulesIsValidOffset,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRules.nextTransition(Ljava/time/Instant;)Ljava/time/zone/ZoneOffsetTransition;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRulesNextTransition,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRules.previousTransition(Ljava/time/Instant;)Ljava/time/zone/ZoneOffsetTransition;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRulesPreviousTransition,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRules.getTransitions()Ljava/util/List;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneRulesGetTransitions,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRules.getTransitionRules()Ljava/util/List;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneRulesGetTransitionRules,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRules.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRulesEquals,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRules.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneRulesHashCode,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRules.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneRulesToString,
		}
}

func createZoneRulesForOffset(offsetObj *object.Object) *object.Object {
	rules := object.MakeEmptyObjectWithClassName(&classNameZoneRules)
	rules.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: offsetObj}
	rules.FieldTable["standardOffset"] = object.Field{Ftype: types.Ref, Fvalue: offsetObj}
	rules.FieldTable["wallOffset"] = object.Field{Ftype: types.Ref, Fvalue: offsetObj}
	rules.FieldTable["fixedOffset"] = object.Field{Ftype: types.Int, Fvalue: int64(1)}
	rules.FieldTable["isFixedOffset"] = object.Field{Ftype: types.Int, Fvalue: int64(1)}
	return rules
}

func createZoneRulesForZoneId(zoneID string) *object.Object {
	rules := object.MakeEmptyObjectWithClassName(&classNameZoneRules)
	rules.FieldTable["zoneId"] = object.Field{Ftype: types.Ref, Fvalue: object.StringObjectFromGoString(zoneID)}
	rules.FieldTable["fixedOffset"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
	rules.FieldTable["isFixedOffset"] = object.Field{Ftype: types.Int, Fvalue: int64(0)}
	return rules
}

func isZoneRulesFixed(rules *object.Object) bool {
	if rules == nil || rules.FieldTable == nil {
		return true
	}
	if fld, ok := rules.FieldTable["isFixedOffset"]; ok {
		if val, okVal := getInt64Val(fld.Fvalue); okVal {
			return val != 0
		}
	}
	if fld, ok := rules.FieldTable["fixedOffset"]; ok {
		if val, okVal := getInt64Val(fld.Fvalue); okVal {
			return val != 0
		}
	}
	_, hasZoneId := rules.FieldTable["zoneId"]
	return !hasZoneId
}

func getZoneRulesZoneIdString(rules *object.Object) string {
	if rules == nil || rules.FieldTable == nil {
		return ""
	}
	if fld, ok := rules.FieldTable["zoneId"]; ok {
		if strObj, okStr := fld.Fvalue.(*object.Object); okStr && object.IsStringObject(strObj) {
			return object.GoStringFromStringObject(strObj)
		}
		if s, okS := fld.Fvalue.(string); okS {
			return s
		}
	}
	return ""
}

func getZoneRulesOffsetObject(rules *object.Object) *object.Object {
	if rules == nil || rules.FieldTable == nil {
		return createZoneOffset(0)
	}
	if fld, ok := rules.FieldTable["offset"]; ok {
		if offObj, okObj := fld.Fvalue.(*object.Object); okObj && offObj != nil {
			return offObj
		}
	}
	if fld, ok := rules.FieldTable["standardOffset"]; ok {
		if offObj, okObj := fld.Fvalue.(*object.Object); okObj && offObj != nil {
			return offObj
		}
	}
	return createZoneOffset(0)
}

func zoneRulesOfOffset(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.of: offset is null")
	}
	offsetObj := params[0].(*object.Object)
	return createZoneRulesForOffset(offsetObj)
}

func zoneRulesOfFull(params []any) any {
	if len(params) < 5 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneRules.of: invalid parameter count")
	}
	baseStandardOffset := params[0].(*object.Object)
	baseWallOffset := params[1].(*object.Object)
	standardOffsetTransitionList := params[2].(*object.Object)
	transitionList := params[3].(*object.Object)
	lastRules := params[4].(*object.Object)

	if baseStandardOffset == nil || object.IsNull(baseStandardOffset) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.of: baseStandardOffset is null")
	}
	if baseWallOffset == nil || object.IsNull(baseWallOffset) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.of: baseWallOffset is null")
	}

	rules := object.MakeEmptyObjectWithClassName(&classNameZoneRules)
	rules.FieldTable["standardOffset"] = object.Field{Ftype: types.Ref, Fvalue: baseStandardOffset}
	rules.FieldTable["wallOffset"] = object.Field{Ftype: types.Ref, Fvalue: baseWallOffset}
	rules.FieldTable["offset"] = object.Field{Ftype: types.Ref, Fvalue: baseWallOffset}
	rules.FieldTable["standardTransitions"] = object.Field{Ftype: types.Ref, Fvalue: standardOffsetTransitionList}
	rules.FieldTable["transitions"] = object.Field{Ftype: types.Ref, Fvalue: transitionList}
	rules.FieldTable["lastRules"] = object.Field{Ftype: types.Ref, Fvalue: lastRules}

	isFixed := int64(0)
	if (standardOffsetTransitionList == nil || object.IsNull(standardOffsetTransitionList)) &&
		(transitionList == nil || object.IsNull(transitionList)) &&
		(lastRules == nil || object.IsNull(lastRules)) {
		isFixed = 1
	}
	rules.FieldTable["fixedOffset"] = object.Field{Ftype: types.Int, Fvalue: isFixed}
	rules.FieldTable["isFixedOffset"] = object.Field{Ftype: types.Int, Fvalue: isFixed}
	return rules
}

func zoneRulesIsFixedOffset(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.isFixedOffset: self is null")
	}
	self := params[0].(*object.Object)
	if isZoneRulesFixed(self) {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func zoneRulesGetOffsetInstant(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.getOffset: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.getOffset: instant is null")
	}
	self := params[0].(*object.Object)
	instantObj := params[1].(*object.Object)

	if isZoneRulesFixed(self) {
		return getZoneRulesOffsetObject(self)
	}

	zoneID := getZoneRulesZoneIdString(self)
	if zoneID != "" {
		epochSec, _ := getInstantEpochSecond(instantObj)
		if loc, locErr := time.LoadLocation(zoneID); locErr == nil {
			t := time.Unix(epochSec, 0).In(loc)
			_, offSec := t.Zone()
			return createZoneOffset(int64(offSec))
		}
		epochMillis := epochSec * 1000
		_, totalOffsetMillis, err := ghelpers.GetZoneOffsets(zoneID, epochMillis)
		if err == nil {
			totalSec := totalOffsetMillis / 1000
			return createZoneOffset(totalSec)
		}
	}
	return getZoneRulesOffsetObject(self)
}

func getOffsetInfoForLDT(self, ldtObj *object.Object) (validOffsets []*object.Object, trans *object.Object, errBlk *ghelpers.GErrBlk) {
	if isZoneRulesFixed(self) {
		off := getZoneRulesOffsetObject(self)
		return []*object.Object{off}, nil, nil
	}

	zoneID := getZoneRulesZoneIdString(self)
	if zoneID == "" {
		off := getZoneRulesOffsetObject(self)
		return []*object.Object{off}, nil, nil
	}

	loc, err := time.LoadLocation(zoneID)
	if err != nil {
		off := getZoneRulesOffsetObject(self)
		return []*object.Object{off}, nil, nil
	}

	dateObj := getLocalDateTimeDate(ldtObj)
	timeObj := getLocalDateTimeTime(ldtObj)
	y, m, d, _ := getLocalDateYMD(dateObj)
	h, min, s, nano := getLocalTimeHMSN(timeObj)
	epochDay := ymdToEpochDay(y, m, d)
	secOfDay := h*3600 + min*60 + s
	wallSec := epochDay*86400 + secOfDay

	// Sample candidate offsets across a 2-day window around wallSec
	candidateOffsetMap := make(map[int]bool)
	var candidateOffsets []int
	for delta := -86400; delta <= 86400; delta += 14400 {
		tSample := time.Unix(wallSec+int64(delta), 0).In(loc)
		_, offSec := tSample.Zone()
		if !candidateOffsetMap[offSec] {
			candidateOffsetMap[offSec] = true
			candidateOffsets = append(candidateOffsets, offSec)
		}
	}

	for _, offSec := range candidateOffsets {
		instSec := wallSec - int64(offSec)
		t := time.Unix(instSec, nano).In(loc)
		_, actualOffSec := t.Zone()
		tY, tM, tD := t.Date()
		tH, tMin, tS := t.Clock()
		if actualOffSec == offSec &&
			int64(tY) == y && int64(tM) == m && int64(tD) == d &&
			int64(tH) == h && int64(tMin) == min && int64(tS) == s {
			offObj := createZoneOffset(int64(offSec))
			validOffsets = append(validOffsets, offObj)
		}
	}

	if len(validOffsets) == 1 {
		return validOffsets, nil, nil
	}

	if len(validOffsets) == 0 {
		// GAP: Clocks jump forward from offsetBefore to offsetAfter.
		// Find offsets before and after using candidate sample instants.
		tBeforeSample := time.Unix(wallSec+86400, 0).In(loc)
		tAfterSample := time.Unix(wallSec-86400, 0).In(loc)
		if len(candidateOffsets) >= 2 {
			tBeforeSample = time.Unix(wallSec-int64(candidateOffsets[0]), 0).In(loc)
			tAfterSample = time.Unix(wallSec-int64(candidateOffsets[1]), 0).In(loc)
		}
		_, off1 := tBeforeSample.Zone()
		_, off2 := tAfterSample.Zone()

		offBeforeSec := off1
		offAfterSec := off2
		if len(candidateOffsets) == 2 {
			// In a spring-forward gap, offsetBefore is smaller (e.g. -18000 vs -14400)
			if candidateOffsets[0] < candidateOffsets[1] {
				offBeforeSec = candidateOffsets[0]
				offAfterSec = candidateOffsets[1]
			} else {
				offBeforeSec = candidateOffsets[1]
				offAfterSec = candidateOffsets[0]
			}
		}

		offsetBefore := createZoneOffset(int64(offBeforeSec))
		offsetAfter := createZoneOffset(int64(offAfterSec))

		// Find exact transition point via binary search in UTC seconds
		low := wallSec - int64(offAfterSec) - 86400
		high := wallSec - int64(offBeforeSec) + 86400
		for low < high {
			mid := low + (high-low)/2
			_, offMid := time.Unix(mid, 0).In(loc).Zone()
			if offMid == offBeforeSec {
				low = mid + 1
			} else {
				high = mid
			}
		}
		transInstant := low
		transLdtSec := transInstant + int64(offBeforeSec)
		transEpochDay := floorDiv(transLdtSec, 86400)
		transSecOfDay := floorMod(transLdtSec, 86400)
		transY, transM, transD := epochDayToYMD(transEpochDay)
		transH := transSecOfDay / 3600
		transMin := (transSecOfDay % 3600) / 60
		transS := transSecOfDay % 60
		transLdt := createLocalDateTime(CreateLocalDate(transY, transM, transD), createLocalTime(transH, transMin, transS, 0))
		trans = createZoneOffsetTransition(transLdt, offsetBefore, offsetAfter)
		return validOffsets, trans, nil
	}

	if len(validOffsets) >= 2 {
		// OVERLAP: Earlier offset (larger, e.g. -04:00 / -14400) was in effect before transition to -05:00 / -18000
		offset1 := validOffsets[0]
		offset2 := validOffsets[1]
		off1 := getZoneOffsetTotalSeconds(offset1)
		off2 := getZoneOffsetTotalSeconds(offset2)

		var offsetBefore, offsetAfter *object.Object
		var offBeforeSec int64
		if off1 > off2 {
			offsetBefore = offset1
			offsetAfter = offset2
			offBeforeSec = off1
		} else {
			offsetBefore = offset2
			offsetAfter = offset1
			offBeforeSec = off2
		}

		// Find exact transition point via binary search
		low := wallSec - offBeforeSec - 86400
		high := wallSec - offBeforeSec + 86400
		for low < high {
			mid := low + (high-low)/2
			_, offMid := time.Unix(mid, 0).In(loc).Zone()
			if int64(offMid) == offBeforeSec {
				low = mid + 1
			} else {
				high = mid
			}
		}
		transInstant := low
		transLdtSec := transInstant + offBeforeSec
		transEpochDay := floorDiv(transLdtSec, 86400)
		transSecOfDay := floorMod(transLdtSec, 86400)
		transY, transM, transD := epochDayToYMD(transEpochDay)
		transH := transSecOfDay / 3600
		transMin := (transSecOfDay % 3600) / 60
		transS := transSecOfDay % 60
		transLdt := createLocalDateTime(CreateLocalDate(transY, transM, transD), createLocalTime(transH, transMin, transS, 0))
		trans = createZoneOffsetTransition(transLdt, offsetBefore, offsetAfter)
		return validOffsets, trans, nil
	}

	return validOffsets, nil, nil
}

func zoneRulesGetOffsetLocalDateTime(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.getOffset: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.getOffset: localDateTime is null")
	}
	self := params[0].(*object.Object)
	ldtObj := params[1].(*object.Object)

	validOffsets, trans, err := getOffsetInfoForLDT(self, ldtObj)
	if err != nil {
		return err
	}
	if len(validOffsets) > 0 {
		return validOffsets[0]
	}
	if trans != nil {
		return getTransitionOffsetBefore(trans)
	}
	return getZoneRulesOffsetObject(self)
}

func zoneRulesGetValidOffsets(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.getValidOffsets: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.getValidOffsets: localDateTime is null")
	}
	self := params[0].(*object.Object)
	ldtObj := params[1].(*object.Object)

	validOffsets, _, err := getOffsetInfoForLDT(self, ldtObj)
	if err != nil {
		return err
	}
	listObj := object.MakeEmptyObjectWithClassName(&classNameArrayList)
	elements := make([]any, len(validOffsets))
	for i, off := range validOffsets {
		elements[i] = off
	}
	listObj.FieldTable["elements"] = object.Field{Ftype: types.Ref, Fvalue: elements}
	return listObj
}

func zoneRulesGetTransition(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.getTransition: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.getTransition: localDateTime is null")
	}
	self := params[0].(*object.Object)
	ldtObj := params[1].(*object.Object)

	_, trans, err := getOffsetInfoForLDT(self, ldtObj)
	if err != nil {
		return err
	}
	return trans
}

func zoneRulesGetStandardOffset(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.getStandardOffset: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.getStandardOffset: instant is null")
	}
	self := params[0].(*object.Object)
	instantObj := params[1].(*object.Object)

	if isZoneRulesFixed(self) {
		return getZoneRulesOffsetObject(self)
	}

	zoneID := getZoneRulesZoneIdString(self)
	if zoneID != "" {
		epochSec, _ := getInstantEpochSecond(instantObj)
		epochMillis := epochSec * 1000
		standardOffsetMillis, _, err := ghelpers.GetZoneOffsets(zoneID, epochMillis)
		if err == nil {
			stdSec := standardOffsetMillis / 1000
			return createZoneOffset(stdSec)
		}
	}
	return getZoneRulesOffsetObject(self)
}

func zoneRulesGetDaylightSavings(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.getDaylightSavings: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.getDaylightSavings: instant is null")
	}
	self := params[0].(*object.Object)
	instantObj := params[1].(*object.Object)

	if isZoneRulesFixed(self) {
		return createDuration(0, 0)
	}

	zoneID := getZoneRulesZoneIdString(self)
	if zoneID != "" {
		epochSec, _ := getInstantEpochSecond(instantObj)
		epochMillis := epochSec * 1000
		stdMillis, totMillis, err := ghelpers.GetZoneOffsets(zoneID, epochMillis)
		if err == nil {
			diffSec := (totMillis - stdMillis) / 1000
			if diffSec < 0 {
				diffSec = 0
			}
			return createDuration(diffSec, 0)
		}
	}
	return createDuration(0, 0)
}

func zoneRulesIsDaylightSavings(params []any) any {
	durRes := zoneRulesGetDaylightSavings(params)
	if err, isErr := durRes.(*ghelpers.GErrBlk); isErr {
		return err
	}
	durObj := durRes.(*object.Object)
	if durObj != nil && durObj.FieldTable != nil {
		if secFld, ok := durObj.FieldTable["seconds"]; ok {
			if sec, okSec := getInt64Val(secFld.Fvalue); okSec && sec > 0 {
				return types.JavaBoolTrue
			}
		}
	}
	return types.JavaBoolFalse
}

func zoneRulesIsValidOffset(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.isValidOffset: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	if len(params) < 3 || params[2] == nil || object.IsNull(params[2].(*object.Object)) {
		return types.JavaBoolFalse
	}
	self := params[0].(*object.Object)
	ldtObj := params[1].(*object.Object)
	offsetObj := params[2].(*object.Object)

	validOffsets, _, err := getOffsetInfoForLDT(self, ldtObj)
	if err != nil {
		return types.JavaBoolFalse
	}
	targetSec := getZoneOffsetTotalSeconds(offsetObj)
	for _, off := range validOffsets {
		if getZoneOffsetTotalSeconds(off) == targetSec {
			return types.JavaBoolTrue
		}
	}
	return types.JavaBoolFalse
}

func zoneRulesNextTransition(params []any) any {
	return nil
}

func zoneRulesPreviousTransition(params []any) any {
	return nil
}

func zoneRulesGetTransitions(params []any) any {
	listObj := object.MakeEmptyObjectWithClassName(&classNameArrayList)
	listObj.FieldTable["elements"] = object.Field{Ftype: types.Ref, Fvalue: []any{}}
	return listObj
}

func zoneRulesGetTransitionRules(params []any) any {
	listObj := object.MakeEmptyObjectWithClassName(&classNameArrayList)
	listObj.FieldTable["elements"] = object.Field{Ftype: types.Ref, Fvalue: []any{}}
	return listObj
}

func zoneRulesEquals(params []any) any {
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
	f1 := isZoneRulesFixed(self)
	f2 := isZoneRulesFixed(other)
	if f1 && f2 {
		o1 := getZoneRulesOffsetObject(self)
		o2 := getZoneRulesOffsetObject(other)
		return zoneOffsetEquals([]any{o1, o2})
	}
	z1 := getZoneRulesZoneIdString(self)
	z2 := getZoneRulesZoneIdString(other)
	if z1 != "" && z1 == z2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func zoneRulesHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.hashCode: self is null")
	}
	self := params[0].(*object.Object)
	if isZoneRulesFixed(self) {
		off := getZoneRulesOffsetObject(self)
		return zoneOffsetHashCode([]any{off})
	}
	z := getZoneRulesZoneIdString(self)
	var h int64 = 0
	for _, c := range z {
		h = 31*h + int64(c)
	}
	return h
}

func zoneRulesToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRules.toString: self is null")
	}
	self := params[0].(*object.Object)
	if isZoneRulesFixed(self) {
		off := getZoneRulesOffsetObject(self)
		offStr := getZoneOffsetID(off)
		return object.StringObjectFromGoString(fmt.Sprintf("ZoneRules[currentStandardOffset=%s]", offStr))
	}
	z := getZoneRulesZoneIdString(self)
	if z != "" {
		return object.StringObjectFromGoString(fmt.Sprintf("ZoneRules[%s]", z))
	}
	return object.StringObjectFromGoString("ZoneRules[]")
}

func getInstantEpochSecond(obj *object.Object) (int64, bool) {
	if obj == nil || obj.FieldTable == nil {
		return 0, false
	}
	if fld, ok := obj.FieldTable["seconds"]; ok {
		if val, okVal := getInt64Val(fld.Fvalue); okVal {
			return val, true
		}
	}
	if fld, ok := obj.FieldTable["epochSecond"]; ok {
		if val, okVal := getInt64Val(fld.Fvalue); okVal {
			return val, true
		}
	}
	return 0, false
}
