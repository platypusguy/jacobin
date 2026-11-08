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
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	classNameZoneRulesProvider     = "java/time/zone/ZoneRulesProvider"
	classNameTzdbZoneRulesProvider = "java/time/zone/TzdbZoneRulesProvider"
	classNameTreeMapZone           = "java/util/TreeMap"

	zoneRulesProviderMutex       sync.RWMutex
	zoneRulesProviderInitialized bool = false

	// Cached list of available zone IDs
	availableZoneIdsList []string
	availableZoneIdsSet  map[string]bool
)

func ensureZoneRulesProviderInited() {
	zoneRulesProviderMutex.Lock()
	defer zoneRulesProviderMutex.Unlock()
	if zoneRulesProviderInitialized {
		return
	}

	ensureZoneOffsetInited()

	availableZoneIdsSet = make(map[string]bool)
	ids, err := ghelpers.GetZoneIDs()
	if err == nil && len(ids) > 0 {
		for _, id := range ids {
			if id != "" {
				availableZoneIdsSet[id] = true
			}
		}
	} else {
		// Fallback standard zone IDs
		standardIDs := []string{
			"UTC", "GMT", "Universal", "Zulu", "Etc/UTC", "Etc/GMT",
		}
		for _, id := range standardIDs {
			availableZoneIdsSet[id] = true
		}
	}

	// Always ensure standard UTC and GMT are present
	availableZoneIdsSet["UTC"] = true
	availableZoneIdsSet["GMT"] = true

	availableZoneIdsList = make([]string, 0, len(availableZoneIdsSet))
	for id := range availableZoneIdsSet {
		availableZoneIdsList = append(availableZoneIdsList, id)
	}
	sort.Strings(availableZoneIdsList)

	zoneRulesProviderInitialized = true
}

// Load_Time_Zone_ZoneRulesProvider registers all java/time/zone/ZoneRulesProvider
// and java/time/zone/TzdbZoneRulesProvider method signatures.
func Load_Time_Zone_ZoneRulesProvider() {
	ensureZoneRulesProviderInited()

	// java/time/zone/ZoneRulesProvider static methods
	ghelpers.MethodSignatures["java/time/zone/ZoneRulesProvider.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneRulesProviderClinit,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRulesProvider.getAvailableZoneIds()Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneRulesProviderGetAvailableZoneIds,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRulesProvider.getRules(Ljava/lang/String;Z)Ljava/time/zone/ZoneRules;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zoneRulesProviderGetRules,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRulesProvider.getVersions(Ljava/lang/String;)Ljava/util/NavigableMap;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRulesProviderGetVersions,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRulesProvider.registerProvider(Ljava/time/zone/ZoneRulesProvider;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRulesProviderRegisterProvider,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRulesProvider.refresh()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneRulesProviderRefresh,
		}

	// java/time/zone/ZoneRulesProvider instance methods
	ghelpers.MethodSignatures["java/time/zone/ZoneRulesProvider.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneRulesProviderInit,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRulesProvider.provideZoneIds()Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneRulesProviderProvideZoneIds,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRulesProvider.provideRules(Ljava/lang/String;Z)Ljava/time/zone/ZoneRules;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zoneRulesProviderProvideRules,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRulesProvider.provideVersions(Ljava/lang/String;)Ljava/util/NavigableMap;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRulesProviderProvideVersions,
		}

	ghelpers.MethodSignatures["java/time/zone/ZoneRulesProvider.provideRefresh()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneRulesProviderProvideRefresh,
		}

	// java/time/zone/TzdbZoneRulesProvider methods
	ghelpers.MethodSignatures["java/time/zone/TzdbZoneRulesProvider.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneRulesProviderClinit,
		}

	ghelpers.MethodSignatures["java/time/zone/TzdbZoneRulesProvider.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneRulesProviderInit,
		}

	ghelpers.MethodSignatures["java/time/zone/TzdbZoneRulesProvider.provideZoneIds()Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneRulesProviderProvideZoneIds,
		}

	ghelpers.MethodSignatures["java/time/zone/TzdbZoneRulesProvider.provideRules(Ljava/lang/String;Z)Ljava/time/zone/ZoneRules;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  zoneRulesProviderProvideRules,
		}

	ghelpers.MethodSignatures["java/time/zone/TzdbZoneRulesProvider.provideVersions(Ljava/lang/String;)Ljava/util/NavigableMap;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  zoneRulesProviderProvideVersions,
		}

	ghelpers.MethodSignatures["java/time/zone/TzdbZoneRulesProvider.provideRefresh()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  zoneRulesProviderProvideRefresh,
		}

	ghelpers.MethodSignatures["java/time/zone/TzdbZoneRulesProvider.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  tzdbZoneRulesProviderToString,
		}
}

func zoneRulesProviderClinit([]any) any {
	ensureZoneRulesProviderInited()
	return nil
}

func zoneRulesProviderInit([]any) any {
	return nil
}

func createNavigableMap(versions map[string]*object.Object) *object.Object {
	tm := object.MakeEmptyObjectWithClassName(&classNameTreeMapZone)
	defMap := make(types.DefHashMap)
	for k, v := range versions {
		defMap[k] = v
	}
	tm.FieldTable["map"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: defMap,
	}
	return tm
}

func zoneRulesProviderGetAvailableZoneIds([]any) any {
	ensureZoneRulesProviderInited()
	zoneRulesProviderMutex.RLock()
	idsCopy := make([]string, len(availableZoneIdsList))
	copy(idsCopy, availableZoneIdsList)
	zoneRulesProviderMutex.RUnlock()

	return ghelpers.CreateStringSet(idsCopy)
}

func isZoneIdAvailableInternal(zoneId string) bool {
	ensureZoneRulesProviderInited()
	zoneRulesProviderMutex.RLock()
	defer zoneRulesProviderMutex.RUnlock()

	if availableZoneIdsSet[zoneId] {
		return true
	}
	if zoneId == "UTC" || zoneId == "GMT" || zoneId == "UT" {
		return true
	}
	for _, prefix := range []string{"UTC+", "UTC-", "GMT+", "GMT-", "UT+", "UT-"} {
		if strings.HasPrefix(zoneId, prefix) {
			return true
		}
	}
	// Check standard location loader as fallback
	if _, err := time.LoadLocation(zoneId); err == nil {
		return true
	}
	return false
}

func getZoneRulesForZoneIdInternal(zoneId string) (*object.Object, *ghelpers.GErrBlk) {
	if zoneId == "" {
		return nil, ghelpers.GetGErrBlk(excNames.DateTimeException, "Unknown time-zone ID: empty string")
	}

	if zoneId == "UTC" || zoneId == "GMT" || zoneId == "UT" {
		return createZoneRulesForOffset(createZoneOffset(0)), nil
	}

	for _, prefix := range []string{"UTC+", "UTC-", "GMT+", "GMT-", "UT+", "UT-"} {
		if strings.HasPrefix(zoneId, prefix) {
			pfx := prefix[:len(prefix)-1]
			offsetStr := zoneId[len(pfx):]
			offsetRes := zoneOffsetOf([]any{offsetStr})
			if err, isErr := offsetRes.(*ghelpers.GErrBlk); isErr {
				return nil, err
			}
			return createZoneRulesForOffset(offsetRes.(*object.Object)), nil
		}
	}

	if strings.HasPrefix(zoneId, "+") || strings.HasPrefix(zoneId, "-") {
		offsetRes := zoneOffsetOf([]any{zoneId})
		if err, isErr := offsetRes.(*ghelpers.GErrBlk); isErr {
			return nil, err
		}
		return createZoneRulesForOffset(offsetRes.(*object.Object)), nil
	}

	if !isZoneIdAvailableInternal(zoneId) {
		return nil, ghelpers.GetGErrBlk(excNames.DateTimeException, "Unknown time-zone ID: "+zoneId)
	}

	return createZoneRulesForZoneId(zoneId), nil
}

func zoneRulesProviderGetRules(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRulesProvider.getRules: zoneId is null")
	}

	zoneIdObj := params[0].(*object.Object)
	zoneIdStr := object.GoStringFromStringObject(zoneIdObj)

	rules, err := getZoneRulesForZoneIdInternal(zoneIdStr)
	if err != nil {
		return err
	}
	return rules
}

func zoneRulesProviderGetVersions(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRulesProvider.getVersions: zoneId is null")
	}

	zoneIdObj := params[0].(*object.Object)
	zoneIdStr := object.GoStringFromStringObject(zoneIdObj)

	rules, err := getZoneRulesForZoneIdInternal(zoneIdStr)
	if err != nil {
		return err
	}

	versionsMap := map[string]*object.Object{
		"TZDB": rules,
	}
	return createNavigableMap(versionsMap)
}

func zoneRulesProviderRegisterProvider(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRulesProvider.registerProvider: provider is null")
	}

	providerObj := params[0].(*object.Object)
	ensureZoneRulesProviderInited()

	// Query provider's zone IDs if available
	zoneRulesProviderMutex.Lock()
	defer zoneRulesProviderMutex.Unlock()

	// If provider has custom IDs in fieldTable, check duplicates
	if providerObj.FieldTable != nil {
		if idsFld, ok := providerObj.FieldTable["zoneIds"]; ok {
			if idList, isList := idsFld.Fvalue.([]string); isList {
				for _, id := range idList {
					if availableZoneIdsSet[id] {
						return ghelpers.GetGErrBlk(excNames.DateTimeException, "Unable to register zone as one already registered with that ID: "+id)
					}
				}
				for _, id := range idList {
					availableZoneIdsSet[id] = true
					availableZoneIdsList = append(availableZoneIdsList, id)
				}
				sort.Strings(availableZoneIdsList)
			}
		}
	}

	return nil
}

func zoneRulesProviderRefresh([]any) any {
	return types.JavaBoolFalse
}

func zoneRulesProviderProvideZoneIds([]any) any {
	return zoneRulesProviderGetAvailableZoneIds(nil)
}

func zoneRulesProviderProvideRules(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRulesProvider.provideRules: self is null")
	}
	// For instance call, params[0] is self, params[1] is zoneId
	var zoneIdObj *object.Object
	if len(params) >= 2 && params[1] != nil && !object.IsNull(params[1].(*object.Object)) {
		zoneIdObj = params[1].(*object.Object)
	} else {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRulesProvider.provideRules: zoneId is null")
	}

	zoneIdStr := object.GoStringFromStringObject(zoneIdObj)
	rules, err := getZoneRulesForZoneIdInternal(zoneIdStr)
	if err != nil {
		return err
	}
	return rules
}

func zoneRulesProviderProvideVersions(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRulesProvider.provideVersions: self is null")
	}
	var zoneIdObj *object.Object
	if len(params) >= 2 && params[1] != nil && !object.IsNull(params[1].(*object.Object)) {
		zoneIdObj = params[1].(*object.Object)
	} else {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "ZoneRulesProvider.provideVersions: zoneId is null")
	}

	zoneIdStr := object.GoStringFromStringObject(zoneIdObj)
	rules, err := getZoneRulesForZoneIdInternal(zoneIdStr)
	if err != nil {
		return err
	}

	versionsMap := map[string]*object.Object{
		"TZDB": rules,
	}
	return createNavigableMap(versionsMap)
}

func zoneRulesProviderProvideRefresh([]any) any {
	return types.JavaBoolFalse
}

func tzdbZoneRulesProviderToString([]any) any {
	return object.StringObjectFromGoString("TZDB")
}
