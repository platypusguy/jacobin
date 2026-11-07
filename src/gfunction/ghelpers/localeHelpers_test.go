/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package ghelpers

import (
	"jacobin/src/excNames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	globals.InitGlobals("test")
	globals.InitStringPool()
	os.Exit(m.Run())
}

func TestStandardLocConstants(t *testing.T) {
	if !RootLoc.IsRootLoc() {
		t.Errorf("RootLoc.IsRootLoc() should be true")
	}
	if RootLoc != (Loc{}) {
		t.Errorf("RootLoc should be empty Loc, got %+v", RootLoc)
	}

	foundRoot := false
	foundUS := false
	foundEnglish := false
	for _, named := range StandardLocConstants {
		if named.FieldName == "ROOT" {
			foundRoot = true
			if named.Value != RootLoc {
				t.Errorf("StandardLocConstants ROOT = %+v, expected RootLoc", named.Value)
			}
		}
		if named.FieldName == "US" {
			foundUS = true
			if named.Value != (Loc{Language: "en", Country: "US"}) {
				t.Errorf("StandardLocConstants US = %+v, expected en_US", named.Value)
			}
		}
		if named.FieldName == "ENGLISH" {
			foundEnglish = true
			if named.Value != (Loc{Language: "en"}) {
				t.Errorf("StandardLocConstants ENGLISH = %+v, expected en", named.Value)
			}
		}
	}
	if !foundRoot || !foundUS || !foundEnglish {
		t.Errorf("Missing expected standard locale constants: root=%v, us=%v, english=%v", foundRoot, foundUS, foundEnglish)
	}
}

func TestNormalizeLoc(t *testing.T) {
	tests := []struct {
		lang     string
		country  string
		variant  string
		expected Loc
	}{
		{"EN", "us", "posix", Loc{Language: "en", Country: "US", Variant: "posix"}},
		{"IW", "IL", "", Loc{Language: "he", Country: "IL"}}, // legacy code conversion
		{"in", "id", "", Loc{Language: "id", Country: "ID"}}, // legacy code in -> id
		{"ji", "ua", "", Loc{Language: "yi", Country: "UA"}}, // legacy code ji -> yi
		{"fr", "FR", "VAR", Loc{Language: "fr", Country: "FR", Variant: "VAR"}},
	}

	for _, tt := range tests {
		result := NormalizeLoc(tt.lang, tt.country, tt.variant)
		if result != tt.expected {
			t.Errorf("NormalizeLoc(%q, %q, %q) = %+v, want %+v", tt.lang, tt.country, tt.variant, result, tt.expected)
		}
	}
}

func TestNormalizeLocWithScript(t *testing.T) {
	tests := []struct {
		lang     string
		script   string
		country  string
		variant  string
		expected Loc
	}{
		{"ZH", "hans", "cn", "", Loc{Language: "zh", Script: "Hans", Country: "CN"}},
		{"ZH", "HANT", "TW", "WIN", Loc{Language: "zh", Script: "Hant", Country: "TW", Variant: "WIN"}},
		{"en", "", "US", "", Loc{Language: "en", Script: "", Country: "US"}},
		{"IW", "hebr", "il", "", Loc{Language: "he", Script: "Hebr", Country: "IL"}},
	}

	for _, tt := range tests {
		result := NormalizeLocWithScript(tt.lang, tt.script, tt.country, tt.variant)
		if result != tt.expected {
			t.Errorf("NormalizeLocWithScript(%q, %q, %q, %q) = %+v, want %+v",
				tt.lang, tt.script, tt.country, tt.variant, result, tt.expected)
		}
	}
}

func TestParseLoc(t *testing.T) {
	tests := []struct {
		input    string
		expected Loc
	}{
		{"", RootLoc},
		{"   ", RootLoc},
		{"C", RootLoc},
		{"POSIX", RootLoc},
		{"de_DE.UTF-8@euro", Loc{Language: "de", Country: "DE"}},
		{"en_US.UTF-8", Loc{Language: "en", Country: "US"}},
		{"zh_CN_#Hans", Loc{Language: "zh", Script: "Hans", Country: "CN"}},
		{"zh__#Hans", Loc{Language: "zh", Script: "Hans"}},
		{"zh-Hans-CN", Loc{Language: "zh", Script: "Hans", Country: "CN"}},
		{"en-US", Loc{Language: "en", Country: "US"}},
		{"und", RootLoc},
		{"und-US", Loc{Country: "US"}},
		{"iw-IL", Loc{Language: "he", Country: "IL"}},
		{"en_US_POSIX", Loc{Language: "en", Country: "US", Variant: "POSIX"}},
		{"en_US_VAR1_VAR2", Loc{Language: "en", Country: "US", Variant: "VAR1_VAR2"}},
		{"_US", Loc{Country: "US"}},
		{"en-001", Loc{Language: "en", Country: "001"}},
		{"en-US-u-ca-japanese", Loc{Language: "en", Country: "US"}}, // singleton 'u' stops subtags
	}

	for _, tt := range tests {
		got := ParseLoc(tt.input)
		if got != tt.expected {
			t.Errorf("ParseLoc(%q) = %+v, want %+v", tt.input, got, tt.expected)
		}
	}
}

func TestToJavaLocaleString(t *testing.T) {
	tests := []struct {
		loc      Loc
		expected string
	}{
		{RootLoc, ""},
		{Loc{Language: "en"}, "en"},
		{Loc{Language: "en", Country: "US"}, "en_US"},
		{Loc{Language: "en", Variant: "POSIX"}, "en__POSIX"},
		{Loc{Country: "US"}, "_US"},
		{Loc{Language: "zh", Script: "Hans", Country: "CN"}, "zh_CN_#Hans"},
		{Loc{Language: "zh", Script: "Hans"}, "zh__#Hans"},
		{Loc{Language: "en", Script: "Latn", Country: "US", Variant: "POSIX"}, "en_US_POSIX_#Latn"},
		{Loc{Country: "US", Variant: "WIN"}, "_US_WIN"},
	}

	for _, tt := range tests {
		got := tt.loc.ToJavaLocaleString()
		if got != tt.expected {
			t.Errorf("Loc%+v.ToJavaLocaleString() = %q, want %q", tt.loc, got, tt.expected)
		}
	}
}

func TestToLanguageTag(t *testing.T) {
	tests := []struct {
		loc      Loc
		expected string
	}{
		{RootLoc, "und"},
		{Loc{Language: "en"}, "en"},
		{Loc{Language: "en", Country: "US"}, "en-US"},
		{Loc{Language: "zh", Script: "Hans", Country: "CN"}, "zh-Hans-CN"},
		{Loc{Language: "en", Country: "US", Variant: "POSIX"}, "en-US-POSIX"},
		{Loc{Language: "en", Country: "US", Variant: "1996"}, "en-US-1996"},
		{Loc{Language: "en", Country: "US", Variant: "foo"}, "en-US-x-lvariant-foo"}, // 3 chars is invalid BCP47 variant
		{Loc{Language: "en", Country: "US", Variant: "12"}, "en-US-x-lvariant-12"},
		{Loc{Language: "de", Variant: "POSIX_inv"}, "de-POSIX-x-lvariant-inv"},
	}

	for _, tt := range tests {
		got := tt.loc.ToLanguageTag()
		if got != tt.expected {
			t.Errorf("Loc%+v.ToLanguageTag() = %q, want %q", tt.loc, got, tt.expected)
		}
	}
}

func TestCandidateLocs(t *testing.T) {
	loc := Loc{Language: "zh", Script: "Hans", Country: "CN", Variant: "FOO"}
	candidates := loc.CandidateLocs()
	expected := []Loc{
		{Language: "zh", Script: "Hans", Country: "CN", Variant: "FOO"},
		{Language: "zh", Script: "Hans", Country: "CN"},
		{Language: "zh", Country: "CN"},
		{Language: "zh", Script: "Hans"},
		{Language: "zh"},
		RootLoc,
	}

	if !reflect.DeepEqual(candidates, expected) {
		t.Errorf("CandidateLocs() = %+v, want %+v", candidates, expected)
	}

	rootCandidates := RootLoc.CandidateLocs()
	expectedRoot := []Loc{RootLoc}
	if !reflect.DeepEqual(rootCandidates, expectedRoot) {
		t.Errorf("RootLoc.CandidateLocs() = %+v, want %+v", rootCandidates, expectedRoot)
	}
}

func TestJavaHashCode(t *testing.T) {
	loc := Loc{Language: "en", Country: "US"}
	hash := loc.JavaHashCode()
	if hash == 0 {
		t.Errorf("Expected non-zero hash code for en_US, got 0")
	}

	// Test deterministic computation
	if hash != loc.JavaHashCode() {
		t.Errorf("JavaHashCode not deterministic")
	}

	// RootLoc hash
	rootHash := RootLoc.JavaHashCode()
	if rootHash != 0 {
		t.Errorf("RootLoc hash code expected 0, got %d", rootHash)
	}

	// Non-BMP unicode (surrogate pair)
	locEmoji := Loc{Language: "en", Variant: "\U0001F600"}
	if locEmoji.JavaHashCode() == 0 {
		t.Errorf("Expected non-zero hash for emoji variant")
	}
}

func TestDefaultLocManagement(t *testing.T) {
	// Set default locale
	customLoc := Loc{Language: "fr", Country: "CA"}
	SetDefaultLoc(customLoc)

	if got := DefaultLoc(LocaleCategoryDisplay); got != customLoc {
		t.Errorf("DefaultLoc(Display) = %+v, want %+v", got, customLoc)
	}
	if got := DefaultLoc(LocaleCategoryFormat); got != customLoc {
		t.Errorf("DefaultLoc(Format) = %+v, want %+v", got, customLoc)
	}

	// Set category specific
	displayLoc := Loc{Language: "de", Country: "DE"}
	SetDefaultLocForCategory(LocaleCategoryDisplay, displayLoc)

	if got := DefaultLoc(LocaleCategoryDisplay); got != displayLoc {
		t.Errorf("DefaultLoc(Display) = %+v, want %+v", got, displayLoc)
	}
	if got := DefaultLoc(LocaleCategoryFormat); got != customLoc {
		t.Errorf("DefaultLoc(Format) = %+v, want %+v", got, customLoc)
	}

	// Category out-of-range clamping -> falls back to Format category
	if got := DefaultLoc(LocaleDefaultCategory(99)); got != customLoc {
		t.Errorf("DefaultLoc(99) = %+v, want %+v", got, customLoc)
	}
	if got := DefaultLoc(LocaleDefaultCategory(-5)); got != customLoc {
		t.Errorf("DefaultLoc(-5) = %+v, want %+v", got, customLoc)
	}

	SetDefaultLocForCategory(LocaleDefaultCategory(-1), Loc{Language: "it", Country: "IT"})
	if got := DefaultLoc(LocaleCategoryFormat); got != (Loc{Language: "it", Country: "IT"}) {
		t.Errorf("DefaultLoc(Format) after clamping set = %+v, want it_IT", got)
	}
}

func TestPlatformAndEnvironmentDefaults(t *testing.T) {
	// Reset initialized state for test
	defaultLocMutex.Lock()
	defaultLocInitialized = false
	locPlatformDefaultSource = nil
	defaultLocMutex.Unlock()

	// Test platform source
	SetLocPlatformDefaultSource(func() string {
		return "ja_JP"
	})
	if got := DefaultLoc(LocaleCategoryDisplay); got != (Loc{Language: "ja", Country: "JP"}) {
		t.Errorf("DefaultLoc with platform source = %+v, want ja_JP", got)
	}

	// Reset initialized state again for env var testing
	defaultLocMutex.Lock()
	defaultLocInitialized = false
	locPlatformDefaultSource = nil
	defaultLocMutex.Unlock()

	origLang := os.Getenv("LANG")
	defer func() {
		_ = os.Setenv("LANG", origLang)
	}()
	_ = os.Setenv("LANG", "es_ES.UTF-8")

	if got := DefaultLoc(LocaleCategoryDisplay); got != (Loc{Language: "es", Country: "ES"}) {
		t.Errorf("DefaultLoc from LANG env = %+v, want es_ES", got)
	}
}

func TestLocToObjectAndLocFromObject(t *testing.T) {
	loc := Loc{Language: "en", Script: "Latn", Country: "GB", Variant: "POSIX"}
	obj := LocToObject(loc)

	if obj == nil || obj.FieldTable == nil {
		t.Fatalf("LocToObject returned nil or empty field table")
	}

	extracted, gerr := LocFromObject(obj)
	if gerr != nil {
		t.Fatalf("LocFromObject failed: %v", gerr)
	}
	if extracted != loc {
		t.Errorf("LocFromObject = %+v, want %+v", extracted, loc)
	}

	// Error cases for LocFromObject
	// 1. Nil value
	if _, err := LocFromObject(nil); err == nil || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("LocFromObject(nil) expected NPE, got %v", err)
	}

	// 2. Null object
	if _, err := LocFromObject(object.Null); err == nil || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("LocFromObject(object.Null) expected NPE, got %v", err)
	}

	// 3. Object missing language field
	invalidObj := object.MakeEmptyObject()
	if _, err := LocFromObject(invalidObj); err == nil || err.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("LocFromObject(invalidObj) expected IllegalArgumentException, got %v", err)
	}

	// 4. Object with null field value
	partialObj := object.MakeEmptyObject()
	partialObj.FieldTable[LocFieldLanguage] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	partialObj.FieldTable[LocFieldScript] = object.Field{Ftype: types.Ref, Fvalue: nil}
	partialObj.FieldTable[LocFieldCountry] = object.Field{Ftype: types.Ref, Fvalue: nil}
	partialObj.FieldTable[LocFieldVariant] = object.Field{Ftype: types.Ref, Fvalue: nil}
	extPartial, err := LocFromObject(partialObj)
	if err != nil {
		t.Errorf("LocFromObject with null field values failed: %v", err)
	}
	if extPartial != RootLoc {
		t.Errorf("LocFromObject with null field values = %+v, want RootLoc", extPartial)
	}
}

func TestLocCategoryFromObject(t *testing.T) {
	// Null cases
	if _, err := LocCategoryFromObject(nil); err == nil || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("LocCategoryFromObject(nil) expected NPE, got %v", err)
	}
	if _, err := LocCategoryFromObject(object.Null); err == nil || err.ExceptionType != excNames.NullPointerException {
		t.Errorf("LocCategoryFromObject(null) expected NPE, got %v", err)
	}

	// Category with name field "DISPLAY"
	displayObj := object.MakeEmptyObject()
	displayObj.FieldTable["name"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: object.StringObjectFromGoString("DISPLAY"),
	}
	cat, err := LocCategoryFromObject(displayObj)
	if err != nil || cat != LocaleCategoryDisplay {
		t.Errorf("LocCategoryFromObject(DISPLAY) = (%v, %v), want (LocaleCategoryDisplay, nil)", cat, err)
	}

	// Category with name field "FORMAT"
	formatObj := object.MakeEmptyObject()
	formatObj.FieldTable["name"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: object.StringObjectFromGoString("FORMAT"),
	}
	cat, err = LocCategoryFromObject(formatObj)
	if err != nil || cat != LocaleCategoryFormat {
		t.Errorf("LocCategoryFromObject(FORMAT) = (%v, %v), want (LocaleCategoryFormat, nil)", cat, err)
	}

	// Category with ordinal field 0
	ord0Obj := object.MakeEmptyObject()
	ord0Obj.FieldTable["ordinal"] = object.Field{
		Ftype:  types.Int,
		Fvalue: int64(0),
	}
	cat, err = LocCategoryFromObject(ord0Obj)
	if err != nil || cat != LocaleCategoryDisplay {
		t.Errorf("LocCategoryFromObject(ord 0) = (%v, %v), want (LocaleCategoryDisplay, nil)", cat, err)
	}

	// Category with ordinal field 1
	ord1Obj := object.MakeEmptyObject()
	ord1Obj.FieldTable["ordinal"] = object.Field{
		Ftype:  types.Int,
		Fvalue: int64(1),
	}
	cat, err = LocCategoryFromObject(ord1Obj)
	if err != nil || cat != LocaleCategoryFormat {
		t.Errorf("LocCategoryFromObject(ord 1) = (%v, %v), want (LocaleCategoryFormat, nil)", cat, err)
	}

	// Unrecognized Category
	unknownObj := object.MakeEmptyObject()
	unknownObj.FieldTable["name"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: object.StringObjectFromGoString("OTHER"),
	}
	unknownObj.FieldTable["ordinal"] = object.Field{
		Ftype:  types.Int,
		Fvalue: int64(99),
	}
	cat, err = LocCategoryFromObject(unknownObj)
	if err == nil || err.ExceptionType != excNames.IllegalArgumentException {
		t.Errorf("LocCategoryFromObject(unknown) expected IllegalArgumentException, got (%v, %v)", cat, err)
	}
}

func TestISOCodesAndDisplayHelpers(t *testing.T) {
	countries := GetISOCountries()
	if len(countries) == 0 || countries[0] != "AD" {
		t.Errorf("GetISOCountries failed")
	}

	countries3 := GetISOCountriesAlpha3()
	if len(countries3) == 0 || countries3[0] != "ABW" {
		t.Errorf("GetISOCountriesAlpha3 failed")
	}

	languages := GetISOLanguages()
	if len(languages) == 0 || languages[0] != "aa" {
		t.Errorf("GetISOLanguages failed")
	}

	available := GetAvailableLocales()
	if len(available) == 0 || available[0] != RootLoc {
		t.Errorf("GetAvailableLocales failed")
	}

	// 3-letter codes
	l3, ok := GetISO3Language("en")
	if !ok || l3 != "eng" {
		t.Errorf("GetISO3Language(en) = (%s, %v), want (eng, true)", l3, ok)
	}
	l3, ok = GetISO3Language("fra")
	if !ok || l3 != "fra" {
		t.Errorf("GetISO3Language(fra) = (%s, %v), want (fra, true)", l3, ok)
	}
	l3, ok = GetISO3Language("")
	if !ok || l3 != "" {
		t.Errorf("GetISO3Language(\"\") = (%s, %v), want (\"\", true)", l3, ok)
	}

	c3, ok := GetISO3Country("US")
	if !ok || c3 != "USA" {
		t.Errorf("GetISO3Country(US) = (%s, %v), want (USA, true)", c3, ok)
	}
	c3, ok = GetISO3Country("CAN")
	if !ok || c3 != "CAN" {
		t.Errorf("GetISO3Country(CAN) = (%s, %v), want (CAN, true)", c3, ok)
	}
	c3, ok = GetISO3Country("")
	if !ok || c3 != "" {
		t.Errorf("GetISO3Country(\"\") = (%s, %v), want (\"\", true)", c3, ok)
	}

	// Display names
	loc := Loc{Language: "en", Script: "Latn", Country: "US", Variant: "POSIX"}
	inLoc := Loc{Language: "en", Country: "US"}
	if dl := GetDisplayLanguage(loc, inLoc); dl != "English" {
		t.Errorf("GetDisplayLanguage = %s, want English", dl)
	}
	if ds := GetDisplayScript(loc, inLoc); ds != "Latin" {
		t.Errorf("GetDisplayScript = %s, want Latin", ds)
	}
	if dc := GetDisplayCountry(loc, inLoc); dc != "United States" {
		t.Errorf("GetDisplayCountry = %s, want United States", dc)
	}
	if dv := GetDisplayVariant(loc, inLoc); dv != "POSIX" {
		t.Errorf("GetDisplayVariant = %s, want POSIX", dv)
	}
	if dn := GetDisplayName(loc, inLoc); dn != "English (Latin, United States, POSIX)" {
		t.Errorf("GetDisplayName = %s, want English (Latin, United States, POSIX)", dn)
	}

	// Fallback when name unknown
	unknownLoc := Loc{Language: "xyz", Script: "Xyz", Country: "ZZ"}
	if dl := GetDisplayLanguage(unknownLoc, inLoc); dl != "xyz" {
		t.Errorf("GetDisplayLanguage(unknown) = %s, want xyz", dl)
	}
	if ds := GetDisplayScript(unknownLoc, inLoc); ds != "Xyz" {
		t.Errorf("GetDisplayScript(unknown) = %s, want Xyz", ds)
	}
	if dc := GetDisplayCountry(unknownLoc, inLoc); dc != "ZZ" {
		t.Errorf("GetDisplayCountry(unknown) = %s, want ZZ", dc)
	}

	// Set and array creation helpers
	setObj := CreateStringSet([]string{"foo", "bar"})
	if setObj == nil {
		t.Fatalf("CreateStringSet returned nil")
	}

	emptySet := CreateEmptySet()
	if emptySet == nil {
		t.Fatalf("CreateEmptySet returned nil")
	}

	arrObj := MakeObjectArray("java/lang/String", []*object.Object{object.StringObjectFromGoString("test")})
	if arrObj == nil {
		t.Fatalf("MakeObjectArray returned nil")
	}
}
