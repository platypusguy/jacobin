/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"reflect"
	"testing"
)

func TestLocaleRegistration(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()
	Load_Util_Locale()

	expectedMethods := []string{
		"java/util/Locale.<clinit>()V",
		"java/util/Locale.<init>(Ljava/lang/String;)V",
		"java/util/Locale.<init>(Ljava/lang/String;Ljava/lang/String;)V",
		"java/util/Locale.<init>(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V",
		"java/util/Locale.clone()Ljava/lang/Object;",
		"java/util/Locale.equals(Ljava/lang/Object;)Z",
		"java/util/Locale.forLanguageTag(Ljava/lang/String;)Ljava/util/Locale;",
		"java/util/Locale.getAvailableLocales()[Ljava/util/Locale;",
		"java/util/Locale.getCountry()Ljava/lang/String;",
		"java/util/Locale.getDefault()Ljava/util/Locale;",
		"java/util/Locale.getDefault(Ljava/util/Locale$Category;)Ljava/util/Locale;",
		"java/util/Locale.getDisplayCountry()Ljava/lang/String;",
		"java/util/Locale.getDisplayCountry(Ljava/util/Locale;)Ljava/lang/String;",
		"java/util/Locale.getDisplayLanguage()Ljava/lang/String;",
		"java/util/Locale.getDisplayLanguage(Ljava/util/Locale;)Ljava/lang/String;",
		"java/util/Locale.getDisplayName()Ljava/lang/String;",
		"java/util/Locale.getDisplayName(Ljava/util/Locale;)Ljava/lang/String;",
		"java/util/Locale.getDisplayScript()Ljava/lang/String;",
		"java/util/Locale.getDisplayScript(Ljava/util/Locale;)Ljava/lang/String;",
		"java/util/Locale.getDisplayVariant()Ljava/lang/String;",
		"java/util/Locale.getDisplayVariant(Ljava/util/Locale;)Ljava/lang/String;",
		"java/util/Locale.getExtension(C)Ljava/lang/String;",
		"java/util/Locale.getExtensionKeys()Ljava/util/Set;",
		"java/util/Locale.getISO3Country()Ljava/lang/String;",
		"java/util/Locale.getISO3Language()Ljava/lang/String;",
		"java/util/Locale.getISOCountries()[Ljava/lang/String;",
		"java/util/Locale.getISOCountries(Ljava/util/Locale$IsoCountryCode;)Ljava/util/Set;",
		"java/util/Locale.getISOLanguages()[Ljava/lang/String;",
		"java/util/Locale.getInstance(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)Ljava/util/Locale;",
		"java/util/Locale.getInstance(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Lsun/util/locale/LocaleExtensions;)Ljava/util/Locale;",
		"java/util/Locale.getInstance(Lsun/util/locale/BaseLocale;Lsun/util/locale/LocaleExtensions;)Ljava/util/Locale;",
		"java/util/Locale.getLanguage()Ljava/lang/String;",
		"java/util/Locale.getScript()Ljava/lang/String;",
		"java/util/Locale.getUnicodeLocaleAttributes()Ljava/util/Set;",
		"java/util/Locale.getUnicodeLocaleKeys()Ljava/util/Set;",
		"java/util/Locale.getUnicodeLocaleType(Ljava/lang/String;)Ljava/lang/String;",
		"java/util/Locale.getVariant()Ljava/lang/String;",
		"java/util/Locale.hasExtensions()Z",
		"java/util/Locale.hashCode()I",
		"java/util/Locale.of(Ljava/lang/String;)Ljava/util/Locale;",
		"java/util/Locale.of(Ljava/lang/String;Ljava/lang/String;)Ljava/util/Locale;",
		"java/util/Locale.of(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)Ljava/util/Locale;",
		"java/util/Locale.setDefault(Ljava/util/Locale;)V",
		"java/util/Locale.setDefault(Ljava/util/Locale$Category;Ljava/util/Locale;)V",
		"java/util/Locale.stripExtensions()Ljava/util/Locale;",
		"java/util/Locale.toLanguageTag()Ljava/lang/String;",
		"java/util/Locale.toString()Ljava/lang/String;",
	}

	for _, sig := range expectedMethods {
		gmeth, ok := ghelpers.MethodSignatures[sig]
		if !ok {
			t.Errorf("Method signature %s not registered", sig)
			continue
		}
		if sig == "java/util/Locale.getInstance(Lsun/util/locale/BaseLocale;Lsun/util/locale/LocaleExtensions;)Ljava/util/Locale;" {
			if reflectSameFunction(gmeth.GFunction, ghelpers.TrapFunction) == false {
				t.Errorf("getInstance(BaseLocale, LocaleExtensions) should remain a trap")
			}
		} else {
			if reflectSameFunction(gmeth.GFunction, ghelpers.TrapFunction) {
				t.Errorf("Method %s is still a trap", sig)
			}
		}
	}
}

func reflectSameFunction(f1, f2 func([]interface{}) interface{}) bool {
	if f1 == nil || f2 == nil {
		return f1 == nil && f2 == nil
	}
	return reflect.ValueOf(f1).Pointer() == reflect.ValueOf(f2).Pointer()
}

func TestLocaleAvailableLocalesAndISOCodes(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()

	// getAvailableLocales
	res := localeGetAvailableLocales([]interface{}{})
	arrObj, ok := res.(*object.Object)
	if !ok || arrObj == nil {
		t.Fatalf("localeGetAvailableLocales failed: %v", res)
	}
	locArr, ok := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
	if !ok || len(locArr) == 0 {
		t.Fatalf("localeGetAvailableLocales returned empty or invalid array")
	}

	// getISOCountries
	res = localeGetISOCountries([]interface{}{})
	arrObj = res.(*object.Object)
	countriesArr := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
	if len(countriesArr) == 0 {
		t.Errorf("localeGetISOCountries returned empty array")
	}
	firstCountry := object.GoStringFromStringObject(countriesArr[0])
	if firstCountry != "AD" {
		t.Errorf("localeGetISOCountries expected first country AD, got %s", firstCountry)
	}

	// getISOLanguages
	res = localeGetISOLanguages([]interface{}{})
	arrObj = res.(*object.Object)
	langsArr := arrObj.FieldTable["value"].Fvalue.([]*object.Object)
	if len(langsArr) == 0 {
		t.Errorf("localeGetISOLanguages returned empty array")
	}
	firstLang := object.GoStringFromStringObject(langsArr[0])
	if firstLang != "aa" {
		t.Errorf("localeGetISOLanguages expected first lang aa, got %s", firstLang)
	}

	// getISOCountries(IsoCountryCode)
	// 1. null param -> NPE
	if err := localeGetISOCountriesForIsoCountryCode([]interface{}{nil}); err == nil {
		t.Errorf("expected NPE for null IsoCountryCode")
	}

	// 2. PART1_ALPHA2
	codeObj := object.MakeEmptyObject()
	codeObj.FieldTable["name"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: object.StringObjectFromGoString("PART1_ALPHA2"),
	}
	setRes := localeGetISOCountriesForIsoCountryCode([]interface{}{codeObj})
	setObj, ok := setRes.(*object.Object)
	if !ok {
		t.Fatalf("expected Set object, got %v", setRes)
	}
	if hm, ok := setObj.FieldTable["map"].Fvalue.(types.DefHashMap); !ok || len(hm) == 0 {
		t.Errorf("expected non-empty set for PART1_ALPHA2")
	}

	// 3. PART1_ALPHA3
	codeObj3 := object.MakeEmptyObject()
	codeObj3.FieldTable["name"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: object.StringObjectFromGoString("PART1_ALPHA3"),
	}
	setRes3 := localeGetISOCountriesForIsoCountryCode([]interface{}{codeObj3})
	setObj3 := setRes3.(*object.Object)
	if hm, ok := setObj3.FieldTable["map"].Fvalue.(types.DefHashMap); !ok || len(hm) == 0 {
		t.Errorf("expected non-empty set for PART1_ALPHA3")
	}

	// 4. PART3
	codeObjP3 := object.MakeEmptyObject()
	codeObjP3.FieldTable["ordinal"] = object.Field{
		Ftype:  types.Int,
		Fvalue: int64(2),
	}
	setResP3 := localeGetISOCountriesForIsoCountryCode([]interface{}{codeObjP3})
	setObjP3 := setResP3.(*object.Object)
	if hm, ok := setObjP3.FieldTable["map"].Fvalue.(types.DefHashMap); !ok || len(hm) != 0 {
		t.Errorf("expected empty set for PART3")
	}

	// 5. Invalid IsoCountryCode
	codeObjInvalid := object.MakeEmptyObject()
	codeObjInvalid.FieldTable["name"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: object.StringObjectFromGoString("INVALID"),
	}
	if err := localeGetISOCountriesForIsoCountryCode([]interface{}{codeObjInvalid}); err == nil {
		t.Errorf("expected error for invalid IsoCountryCode")
	}
}

func TestLocaleExtensionsAndUnicode(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()

	loc := ghelpers.Loc{Language: "en", Country: "US"}
	locObj := ghelpers.LocToObject(loc)

	// getExtensionKeys
	res := localeGetExtensionKeys([]interface{}{locObj})
	setObj, ok := res.(*object.Object)
	if !ok || len(setObj.FieldTable["map"].Fvalue.(types.DefHashMap)) != 0 {
		t.Errorf("expected empty set for getExtensionKeys")
	}

	// getUnicodeLocaleAttributes
	res = localeGetUnicodeLocaleAttributes([]interface{}{locObj})
	setObj = res.(*object.Object)
	if len(setObj.FieldTable["map"].Fvalue.(types.DefHashMap)) != 0 {
		t.Errorf("expected empty set for getUnicodeLocaleAttributes")
	}

	// getUnicodeLocaleKeys
	res = localeGetUnicodeLocaleKeys([]interface{}{locObj})
	setObj = res.(*object.Object)
	if len(setObj.FieldTable["map"].Fvalue.(types.DefHashMap)) != 0 {
		t.Errorf("expected empty set for getUnicodeLocaleKeys")
	}

	// getExtension
	res = localeGetExtension([]interface{}{locObj, int64('u')})
	if res != object.Null {
		t.Errorf("expected object.Null for getExtension, got %v", res)
	}

	// getUnicodeLocaleType
	res = localeGetUnicodeLocaleType([]interface{}{locObj, object.StringObjectFromGoString("ca")})
	if res != object.Null {
		t.Errorf("expected object.Null for getUnicodeLocaleType, got %v", res)
	}

	// getUnicodeLocaleType with null key
	res = localeGetUnicodeLocaleType([]interface{}{locObj, nil})
	if geb, ok := res.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.NullPointerException {
		t.Errorf("expected NPE for getUnicodeLocaleType with null key")
	}
}

func TestLocaleDisplayNamesAndISO3(t *testing.T) {
	globals.InitGlobals("test")
	globals.InitStringPool()

	loc := ghelpers.Loc{Language: "fr", Script: "Latn", Country: "CA", Variant: "POSIX"}
	locObj := ghelpers.LocToObject(loc)
	inLoc := ghelpers.Loc{Language: "en", Country: "US"}
	inLocObj := ghelpers.LocToObject(inLoc)

	// getDisplayLanguage
	res := localeGetDisplayLanguage([]interface{}{locObj})
	if str := object.GoStringFromStringObject(res.(*object.Object)); str != "French" {
		t.Errorf("getDisplayLanguage expected French, got %s", str)
	}
	res = localeGetDisplayLanguageInLocale([]interface{}{locObj, inLocObj})
	if str := object.GoStringFromStringObject(res.(*object.Object)); str != "French" {
		t.Errorf("getDisplayLanguageInLocale expected French, got %s", str)
	}

	// getDisplayScript
	res = localeGetDisplayScript([]interface{}{locObj})
	if str := object.GoStringFromStringObject(res.(*object.Object)); str != "Latin" {
		t.Errorf("getDisplayScript expected Latin, got %s", str)
	}
	res = localeGetDisplayScriptInLocale([]interface{}{locObj, inLocObj})
	if str := object.GoStringFromStringObject(res.(*object.Object)); str != "Latin" {
		t.Errorf("getDisplayScriptInLocale expected Latin, got %s", str)
	}

	// getDisplayCountry
	res = localeGetDisplayCountry([]interface{}{locObj})
	if str := object.GoStringFromStringObject(res.(*object.Object)); str != "Canada" {
		t.Errorf("getDisplayCountry expected Canada, got %s", str)
	}
	res = localeGetDisplayCountryInLocale([]interface{}{locObj, inLocObj})
	if str := object.GoStringFromStringObject(res.(*object.Object)); str != "Canada" {
		t.Errorf("getDisplayCountryInLocale expected Canada, got %s", str)
	}

	// getDisplayVariant
	res = localeGetDisplayVariant([]interface{}{locObj})
	if str := object.GoStringFromStringObject(res.(*object.Object)); str != "POSIX" {
		t.Errorf("getDisplayVariant expected POSIX, got %s", str)
	}
	res = localeGetDisplayVariantInLocale([]interface{}{locObj, inLocObj})
	if str := object.GoStringFromStringObject(res.(*object.Object)); str != "POSIX" {
		t.Errorf("getDisplayVariantInLocale expected POSIX, got %s", str)
	}

	// getDisplayName
	res = localeGetDisplayName([]interface{}{locObj})
	if str := object.GoStringFromStringObject(res.(*object.Object)); str != "French (Latin, Canada, POSIX)" {
		t.Errorf("getDisplayName expected French (Latin, Canada, POSIX), got %s", str)
	}
	res = localeGetDisplayNameInLocale([]interface{}{locObj, inLocObj})
	if str := object.GoStringFromStringObject(res.(*object.Object)); str != "French (Latin, Canada, POSIX)" {
		t.Errorf("getDisplayNameInLocale expected French (Latin, Canada, POSIX), got %s", str)
	}

	// Root display name
	rootObj := ghelpers.LocToObject(ghelpers.RootLoc)
	res = localeGetDisplayName([]interface{}{rootObj})
	if str := object.GoStringFromStringObject(res.(*object.Object)); str != "" {
		t.Errorf("root getDisplayName expected empty string, got %s", str)
	}

	// Only country display name
	onlyCountryObj := ghelpers.LocToObject(ghelpers.Loc{Country: "US"})
	res = localeGetDisplayName([]interface{}{onlyCountryObj})
	if str := object.GoStringFromStringObject(res.(*object.Object)); str != "(United States)" {
		t.Errorf("only country getDisplayName expected (United States), got %s", str)
	}

	// getISO3Language
	res = localeGetISO3Language([]interface{}{locObj})
	if str := object.GoStringFromStringObject(res.(*object.Object)); str != "fra" {
		t.Errorf("getISO3Language expected fra, got %s", str)
	}

	// getISO3Country
	res = localeGetISO3Country([]interface{}{locObj})
	if str := object.GoStringFromStringObject(res.(*object.Object)); str != "CAN" {
		t.Errorf("getISO3Country expected CAN, got %s", str)
	}

	// getISO3Language for root -> empty string
	res = localeGetISO3Language([]interface{}{rootObj})
	if str := object.GoStringFromStringObject(res.(*object.Object)); str != "" {
		t.Errorf("getISO3Language root expected empty string, got %s", str)
	}

	// getISO3Country for root -> empty string
	res = localeGetISO3Country([]interface{}{rootObj})
	if str := object.GoStringFromStringObject(res.(*object.Object)); str != "" {
		t.Errorf("getISO3Country root expected empty string, got %s", str)
	}

	// Unknown 2-letter lang -> MissingResourceException
	unknownLangObj := ghelpers.LocToObject(ghelpers.Loc{Language: "zz"})
	res = localeGetISO3Language([]interface{}{unknownLangObj})
	if geb, ok := res.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.MissingResourceException {
		t.Errorf("expected MissingResourceException for unknown language, got %v", res)
	}

	// Unknown 2-letter country -> MissingResourceException
	unknownCountryObj := ghelpers.LocToObject(ghelpers.Loc{Country: "ZZ"})
	res = localeGetISO3Country([]interface{}{unknownCountryObj})
	if geb, ok := res.(*ghelpers.GErrBlk); !ok || geb.ExceptionType != excNames.MissingResourceException {
		t.Errorf("expected MissingResourceException for unknown country, got %v", res)
	}
}
