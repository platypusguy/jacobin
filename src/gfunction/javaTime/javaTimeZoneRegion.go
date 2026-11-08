/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/types"
	"time"
)

var (
	zoneRegionInited bool
)

func ensureZoneRegionInited() {
	if zoneRegionInited {
		return
	}
	ensureZoneIdInited()
	ensureZoneOffsetInited()
	zoneRegionInited = true
}

// Load_Time_ZoneRegion registers all java/time/ZoneRegion method signatures.
func Load_Time_ZoneRegion() {
	ensureZoneRegionInited()

	ghelpers.MethodSignatures["java/time/ZoneRegion.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneRegionClinit,
		}

	ghelpers.MethodSignatures["java/time/ZoneRegion.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneRegionInit,
		}

	ghelpers.MethodSignatures["java/time/ZoneRegion.<init>(Ljava/lang/String;Ljava/time/zone/ZoneRules;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zoneRegionInitWithIdRules,
		}

	ghelpers.MethodSignatures["java/time/ZoneRegion.ofId(Ljava/lang/String;Z)Ljava/time/ZoneRegion;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zoneRegionOfId,
		}

	ghelpers.MethodSignatures["java/time/ZoneRegion.ofWithPrefix(Ljava/lang/String;IZ)Ljava/time/ZoneRegion;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  zoneRegionOfWithPrefix,
		}

	ghelpers.MethodSignatures["java/time/ZoneRegion.checkName(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRegionCheckName,
		}

	ghelpers.MethodSignatures["java/time/ZoneRegion.getId()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneIdGetId,
		}

	ghelpers.MethodSignatures["java/time/ZoneRegion.getRules()Ljava/time/zone/ZoneRules;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneIdGetRules,
		}

	ghelpers.MethodSignatures["java/time/ZoneRegion.getDisplayName(Ljava/time/format/TextStyle;Ljava/util/Locale;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zoneIdGetDisplayName,
		}

	ghelpers.MethodSignatures["java/time/ZoneRegion.normalized()Ljava/time/ZoneId;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneIdNormalized,
		}

	ghelpers.MethodSignatures["java/time/ZoneRegion.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneIdEquals,
		}

	ghelpers.MethodSignatures["java/time/ZoneRegion.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneIdHashCode,
		}

	ghelpers.MethodSignatures["java/time/ZoneRegion.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneIdToString,
		}

	ghelpers.MethodSignatures["java/time/ZoneRegion.write(Ljava/io/DataOutput;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRegionWrite,
		}

	ghelpers.MethodSignatures["java/time/ZoneRegion.writeExternal(Ljava/io/DataOutput;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRegionWriteExternal,
		}

	ghelpers.MethodSignatures["java/time/ZoneRegion.readExternal(Ljava/io/DataInput;)Ljava/time/ZoneId;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRegionReadExternal,
		}
}

func zoneRegionClinit([]any) any {
	ensureZoneRegionInited()
	return nil
}

func zoneRegionInit([]any) any {
	return nil
}

func zoneRegionInitWithIdRules(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRegion.<init>: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRegion.<init>: zoneId is null")
	}
	self := params[0].(*object.Object)
	idObj := params[1].(*object.Object)
	self.FieldTable["id"] = object.Field{Ftype: types.Ref, Fvalue: idObj}
	if len(params) >= 3 && params[2] != nil && !object.IsNull(params[2].(*object.Object)) {
		rulesObj := params[2].(*object.Object)
		self.FieldTable["rules"] = object.Field{Ftype: types.Ref, Fvalue: rulesObj}
	}
	self.FieldTable["zoneId"] = object.Field{Ftype: types.Ref, Fvalue: self}
	self.FieldTable["zone"] = object.Field{Ftype: types.Ref, Fvalue: self}
	return nil
}

func zoneRegionCheckNameInternal(zoneIdStr string) *ghelpers.GErrBlk {
	n := len(zoneIdStr)
	if n < 2 {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Invalid ID for region-based ZoneId, invalid format: "+zoneIdStr)
	}
	for i := 0; i < n; i++ {
		c := zoneIdStr[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			continue
		}
		if c == '/' && i != 0 {
			continue
		}
		if (c >= '0' && c <= '9' || c == '~' || c == '.' || c == '_' || c == '+' || c == '-') && i != 0 {
			continue
		}
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Invalid ID for region-based ZoneId, invalid format: "+zoneIdStr)
	}
	return nil
}

func zoneRegionCheckName(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRegion.checkName: zoneId is null")
	}
	zoneIdObj := params[0].(*object.Object)
	zoneIdStr := object.GoStringFromStringObject(zoneIdObj)
	if err := zoneRegionCheckNameInternal(zoneIdStr); err != nil {
		return err
	}
	return nil
}

func zoneRegionOfId(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRegion.ofId: zoneId is null")
	}
	zoneIdObj := params[0].(*object.Object)
	zoneIdStr := object.GoStringFromStringObject(zoneIdObj)

	checkAvailable := true
	if len(params) >= 2 && params[1] != nil {
		if b, ok := getInt64Val(params[1]); ok {
			checkAvailable = b != types.JavaBoolFalse
		}
	}

	if err := zoneRegionCheckNameInternal(zoneIdStr); err != nil {
		return err
	}

	if checkAvailable {
		ids, err := ghelpers.GetZoneIDs()
		if err == nil && len(ids) > 0 {
			found := false
			for _, id := range ids {
				if id == zoneIdStr {
					found = true
					break
				}
			}
			if !found {
				return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unknown time-zone ID: "+zoneIdStr)
			}
		} else {
			if _, locErr := time.LoadLocation(zoneIdStr); locErr != nil {
				return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unknown time-zone ID: "+zoneIdStr)
			}
		}
	}

	return createZoneRegion(zoneIdStr, createZoneRulesForZoneId(zoneIdStr))
}

func zoneRegionOfWithPrefix(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRegion.ofWithPrefix: zoneId is null")
	}
	if len(params) < 2 || params[1] == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRegion.ofWithPrefix: prefixLength is null")
	}

	zoneIdObj := params[0].(*object.Object)
	zoneIdStr := object.GoStringFromStringObject(zoneIdObj)

	var prefixLength int
	if pl, ok := getInt64Val(params[1]); ok {
		prefixLength = int(pl)
	} else if pli, ok2 := params[1].(int); ok2 {
		prefixLength = pli
	} else {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "ZoneRegion.ofWithPrefix: invalid prefixLength")
	}

	if prefixLength < 0 || prefixLength > len(zoneIdStr) {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Invalid prefix length for ZoneRegion: "+zoneIdStr)
	}

	checkAvailable := true
	if len(params) >= 3 && params[2] != nil {
		if b, ok := getInt64Val(params[2]); ok {
			checkAvailable = b != types.JavaBoolFalse
		}
	}

	prefix := zoneIdStr[:prefixLength]
	if len(zoneIdStr) == prefixLength {
		return zoneRegionOfId([]any{object.StringObjectFromGoString(prefix), types.ConvertGoBoolToJavaBool(checkAvailable)})
	}

	if zoneIdStr[prefixLength] != '+' && zoneIdStr[prefixLength] != '-' {
		return zoneRegionOfId([]any{zoneIdObj, types.ConvertGoBoolToJavaBool(checkAvailable)})
	}

	offsetStr := zoneIdStr[prefixLength:]
	offsetRes := zoneOffsetOf([]any{object.StringObjectFromGoString(offsetStr)})
	if _, isErr := offsetRes.(*ghelpers.GErrBlk); isErr {
		return ghelpers.GetGErrBlk(excNames.DateTimeException, "Invalid ID for offset-based ZoneId: "+zoneIdStr)
	}

	offsetObj := offsetRes.(*object.Object)
	totalSec := getZoneOffsetTotalSeconds(offsetObj)
	if totalSec == 0 {
		return zoneRegionOfId([]any{object.StringObjectFromGoString(prefix), types.ConvertGoBoolToJavaBool(checkAvailable)})
	}

	fullId := prefix + getZoneOffsetID(offsetObj)
	rules := createZoneRulesForOffset(offsetObj)
	return createZoneRegion(fullId, rules)
}

func zoneRegionWrite(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRegion.write: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRegion.write: out is null")
	}
	return nil
}

func zoneRegionWriteExternal(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRegion.writeExternal: self is null")
	}
	if len(params) < 2 || params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRegion.writeExternal: out is null")
	}
	return nil
}

func zoneRegionReadExternal(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRegion.readExternal: in is null")
	}
	return makeZoneIdFromID("UTC")
}
