/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package sunSecurity

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"testing"
)

func TestLoad_Sun_Security_Jca_ProviderList(t *testing.T) {
	globals.InitGlobals("test")
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Sun_Security_Jca_ProviderList()

	expectedSignatures := []struct {
		sig   string
		slots int
	}{
		{"sun/security/jca/ProviderList.<clinit>()V", 0},
		{"sun/security/jca/ProviderList.<init>()V", 0},
		{"sun/security/jca/ProviderList.add(Lsun/security/jca/ProviderConfig;)Lsun/security/jca/ProviderList;", 1},
		{"sun/security/jca/ProviderList.fromSecurityProperties()Lsun/security/jca/ProviderList;", 0},
		{"sun/security/jca/ProviderList.getDefault()Lsun/security/jca/ProviderList;", 0},
		{"sun/security/jca/ProviderList.getJarList(Ljava/lang/String;)Ljava/util/List;", 1},
		{"sun/security/jca/ProviderList.getProvider(Ljava/lang/String;)Ljava/security/Provider;", 1},
		{"sun/security/jca/ProviderList.getProviderConfig(Ljava/lang/String;)Lsun/security/jca/ProviderConfig;", 1},
		{"sun/security/jca/ProviderList.getProviderConfigs()Ljava/util/List;", 0},
		{"sun/security/jca/ProviderList.getService(Ljava/lang/String;Ljava/lang/String;)Ljava/security/Provider$Service;", 2},
		{"sun/security/jca/ProviderList.insertAt(Lsun/security/jca/ProviderConfig;I)Lsun/security/jca/ProviderList;", 2},
		{"sun/security/jca/ProviderList.isEmpty()Z", 0},
		{"sun/security/jca/ProviderList.loadAll()V", 0},
		{"sun/security/jca/ProviderList.newList(Lsun/security/jca/ProviderConfig;)Lsun/security/jca/ProviderList;", 1},
		{"sun/security/jca/ProviderList.providers()Ljava/util/List;", 0},
		{"sun/security/jca/ProviderList.remove(Ljava/lang/String;)Lsun/security/jca/ProviderList;", 1},
		{"sun/security/jca/ProviderList.removeInvalid()Lsun/security/jca/ProviderList;", 0},
		{"sun/security/jca/ProviderList.size()I", 0},
		{"sun/security/jca/ProviderList.toArray()[Ljava/security/Provider;", 0},
		{"sun/security/jca/ProviderList.toString()Ljava/lang/String;", 0},
		{"sun/security/jca/ProviderList$ServiceList.tryGet(I)Ljava/security/Provider$Service;", 1},
	}

	for _, m := range expectedSignatures {
		gm, ok := ghelpers.MethodSignatures[m.sig]
		if !ok {
			t.Errorf("Method signature not registered: %s", m.sig)
			continue
		}
		if gm.ParamSlots != m.slots {
			t.Errorf("%s: expected %d param slots, got %d", m.sig, m.slots, gm.ParamSlots)
		}
	}
}

func TestProviderList_clinit(t *testing.T) {
	globals.InitGlobals("test")
	// Clear statics manually since ClearStatics doesn't exist
	statics.Statics = make(map[string]statics.Static)

	clinitProviderList(nil)

	thisClassName := "sun/security/jca/ProviderList"

	// Check debug
	debug := statics.GetStaticValue(thisClassName, "debug")
	if !object.IsNull(debug) {
		t.Errorf("Expected debug to be null")
	}

	// Check PC0
	pc0 := statics.GetStaticValue("L"+thisClassName, "PC0")
	if pc0 == nil || object.IsNull(pc0) {
		t.Errorf("PC0 static not found or null")
	}

	// Check P0
	p0 := statics.GetStaticValue("L"+thisClassName, "P0")
	if p0 == nil || object.IsNull(p0) {
		t.Errorf("P0 static not found or null")
	}

	// Check EMPTY
	emptyPL := statics.GetStaticValue(thisClassName, "EMPTY")
	if emptyPL == nil || object.IsNull(emptyPL) {
		t.Errorf("EMPTY static not found or null")
	}
}

func TestProviderList_StaticFactories(t *testing.T) {
	globals.InitGlobals("test")
	statics.Statics = make(map[string]statics.Static)
	clinitProviderList(nil)

	// fromSecurityProperties
	ret := providerListFromSecurityProperties(nil)
	pl, ok := ret.(*object.Object)
	if !ok {
		t.Fatalf("fromSecurityProperties returned %T, expected *object.Object", ret)
	}
	if pl.FieldTable["providers"].Fvalue == nil {
		t.Error("providers field is nil")
	}

	// getDefault
	ret = providerListGetDefault(nil)
	pl, ok = ret.(*object.Object)
	if !ok {
		t.Fatalf("getDefault returned %T, expected *object.Object", ret)
	}
	if pl.FieldTable["providers"].Fvalue == nil {
		t.Error("providers field is nil")
	}
}

func TestProviderList_InstanceMethods(t *testing.T) {
	globals.InitGlobals("test")
	statics.Statics = make(map[string]statics.Static)
	clinitProviderList(nil)

	thisClassName := "sun/security/jca/ProviderList"
	plObj := object.MakeEmptyObjectWithClassName(&thisClassName)

	p0 := statics.GetStaticValue("L"+thisClassName, "P0").(*object.Object)
	plObj.FieldTable["providers"] = object.Field{Ftype: types.Ref, Fvalue: p0}

	// providers()
	ret := providerListProviders([]any{plObj})
	providersArr, ok := ret.(*object.Object)
	if !ok {
		t.Fatalf("providers() returned %T, expected *object.Object", ret)
	}
	if providersArr != p0 {
		t.Error("providers() mismatch")
	}

	// size()
	size := providerListSize(nil).(int64)
	if size != 1 {
		t.Errorf("Expected size 1, got %d", size)
	}

	// tryGet()
	ret = providerListTryGet([]any{plObj, int64(0)})
	if !object.IsNull(ret) {
		t.Error("tryGet() should return null")
	}
}
