/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023-6 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
)

// Implement java/util/Locale (Java 21).
//
// Strategy: the Go type ghelpers.Loc (see localeHelpers.go) is the source of truth.
// A Locale object is a carrier of four String fields: language, script, country, variant.
// Every gfunction converts with ghelpers.LocFromObject / ghelpers.LocToObject and works with Go values.
// Locale extensions (-u-, -x-) are not supported yet.

const localeClassPrefix = "java/util/Locale."
const localeStaticType = "Ljava/util/Locale;"

func Load_Util_Locale() {

	// On Windows, the initial default locale comes from the platform rather than from environment variables.
	if globals.OnWindows {
		ghelpers.SetLocPlatformDefaultSource(_getDefaultLanguageFromWindows)
	}

	ghelpers.MethodSignatures["java/util/Locale.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeClinit,
		}

	ghelpers.MethodSignatures["java/util/Locale.<init>(Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localeInit,
		}

	ghelpers.MethodSignatures["java/util/Locale.<init>(Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localeInit,
		}

	ghelpers.MethodSignatures["java/util/Locale.<init>(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  localeInit,
		}

	ghelpers.MethodSignatures["java/util/Locale.clone()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeClone,
		}

	ghelpers.MethodSignatures["java/util/Locale.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localeEquals,
		}

	ghelpers.MethodSignatures["java/util/Locale.forLanguageTag(Ljava/lang/String;)Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localeForLanguageTag,
		}

	ghelpers.MethodSignatures["java/util/Locale.getAvailableLocales()[Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeGetAvailableLocales,
		}

	ghelpers.MethodSignatures["java/util/Locale.getCountry()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeComponentGetter(func(loc ghelpers.Loc) string { return loc.Country }),
		}

	ghelpers.MethodSignatures["java/util/Locale.getDefault()Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeGetDefault,
		}

	ghelpers.MethodSignatures["java/util/Locale.getDefault(Ljava/util/Locale$Category;)Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localeGetDefaultForCategory,
		}

	ghelpers.MethodSignatures["java/util/Locale.getDisplayCountry()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeGetDisplayCountry,
		}

	ghelpers.MethodSignatures["java/util/Locale.getDisplayCountry(Ljava/util/Locale;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localeGetDisplayCountryInLocale,
		}

	ghelpers.MethodSignatures["java/util/Locale.getDisplayLanguage()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeGetDisplayLanguage,
		}

	ghelpers.MethodSignatures["java/util/Locale.getDisplayLanguage(Ljava/util/Locale;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localeGetDisplayLanguageInLocale,
		}

	ghelpers.MethodSignatures["java/util/Locale.getDisplayName()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeGetDisplayName,
		}

	ghelpers.MethodSignatures["java/util/Locale.getDisplayName(Ljava/util/Locale;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localeGetDisplayNameInLocale,
		}

	ghelpers.MethodSignatures["java/util/Locale.getDisplayScript()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeGetDisplayScript,
		}

	ghelpers.MethodSignatures["java/util/Locale.getDisplayScript(Ljava/util/Locale;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localeGetDisplayScriptInLocale,
		}

	ghelpers.MethodSignatures["java/util/Locale.getDisplayVariant()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeGetDisplayVariant,
		}

	ghelpers.MethodSignatures["java/util/Locale.getDisplayVariant(Ljava/util/Locale;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localeGetDisplayVariantInLocale,
		}

	ghelpers.MethodSignatures["java/util/Locale.getExtension(C)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localeGetExtension,
		}

	ghelpers.MethodSignatures["java/util/Locale.getExtensionKeys()Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeGetExtensionKeys,
		}

	ghelpers.MethodSignatures["java/util/Locale.getISO3Country()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeGetISO3Country,
		}

	ghelpers.MethodSignatures["java/util/Locale.getISO3Language()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeGetISO3Language,
		}

	ghelpers.MethodSignatures["java/util/Locale.getISOCountries()[Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeGetISOCountries,
		}

	ghelpers.MethodSignatures["java/util/Locale.getISOCountries(Ljava/util/Locale$IsoCountryCode;)Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localeGetISOCountriesForIsoCountryCode,
		}

	ghelpers.MethodSignatures["java/util/Locale.getISOLanguages()[Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeGetISOLanguages,
		}

	ghelpers.MethodSignatures["java/util/Locale.getInstance(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  localeGetInstanceLanguageCountryVariant,
		}

	ghelpers.MethodSignatures["java/util/Locale.getInstance(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Lsun/util/locale/LocaleExtensions;)Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  localeGetInstanceWithScript,
		}

	ghelpers.MethodSignatures["java/util/Locale.getInstance(Lsun/util/locale/BaseLocale;Lsun/util/locale/LocaleExtensions;)Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Locale.getLanguage()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeComponentGetter(func(loc ghelpers.Loc) string { return loc.Language }),
		}

	ghelpers.MethodSignatures["java/util/Locale.getScript()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeComponentGetter(func(loc ghelpers.Loc) string { return loc.Script }),
		}

	ghelpers.MethodSignatures["java/util/Locale.getUnicodeLocaleAttributes()Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeGetUnicodeLocaleAttributes,
		}

	ghelpers.MethodSignatures["java/util/Locale.getUnicodeLocaleKeys()Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeGetUnicodeLocaleKeys,
		}

	ghelpers.MethodSignatures["java/util/Locale.getUnicodeLocaleType(Ljava/lang/String;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localeGetUnicodeLocaleType,
		}

	ghelpers.MethodSignatures["java/util/Locale.getVariant()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeComponentGetter(func(loc ghelpers.Loc) string { return loc.Variant }),
		}

	ghelpers.MethodSignatures["java/util/Locale.hasExtensions()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeHasExtensions,
		}

	ghelpers.MethodSignatures["java/util/Locale.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeHashCode,
		}

	ghelpers.MethodSignatures["java/util/Locale.of(Ljava/lang/String;)Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localeOf,
		}

	ghelpers.MethodSignatures["java/util/Locale.of(Ljava/lang/String;Ljava/lang/String;)Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localeOf,
		}

	ghelpers.MethodSignatures["java/util/Locale.of(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  localeOf,
		}

	ghelpers.MethodSignatures["java/util/Locale.setDefault(Ljava/util/Locale;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  localeSetDefault,
		}

	ghelpers.MethodSignatures["java/util/Locale.setDefault(Ljava/util/Locale$Category;Ljava/util/Locale;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  localeSetDefaultForCategory,
		}

	ghelpers.MethodSignatures["java/util/Locale.stripExtensions()Ljava/util/Locale;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeStripExtensions,
		}

	ghelpers.MethodSignatures["java/util/Locale.toLanguageTag()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeToLanguageTag,
		}

	ghelpers.MethodSignatures["java/util/Locale.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  localeToString,
		}
}

// ---------------------------------------------------------------------------------------------
// Class initialization
// ---------------------------------------------------------------------------------------------

// Populate the public static constants: Locale.ROOT, Locale.US, Locale.ENGLISH, and so on.
// "java/util/Locale.<clinit>()V"
func localeClinit([]interface{}) interface{} {
	for _, constant := range ghelpers.StandardLocConstants {
		statics.AddStatic(localeClassPrefix+constant.FieldName,
			statics.Static{Type: localeStaticType, Value: ghelpers.LocToObject(constant.Value)})
	}
	return nil
}

// ---------------------------------------------------------------------------------------------
// Constructors and factories
// ---------------------------------------------------------------------------------------------

// Convert Locale gfunction string arguments to Go strings; a null argument is a NullPointerException.
func localeStringArguments(functionName string, arguments []interface{}) ([]string, *ghelpers.GErrBlk) {
	goStrings := make([]string, len(arguments))
	for index, argument := range arguments {
		stringObj, ok := argument.(*object.Object)
		if !ok || object.IsNull(stringObj) {
			return nil, ghelpers.GetGErrBlk(excNames.NullPointerException,
				fmt.Sprintf("%s: argument %d is null", functionName, index+1))
		}
		goStrings[index] = object.GoStringFromStringObject(stringObj)
	}
	return goStrings, nil
}

// language, country, variant; missing trailing values are empty.
func localeFromLanguageCountryVariant(goStrings []string) ghelpers.Loc {
	var padded [3]string
	copy(padded[:], goStrings)
	return ghelpers.NormalizeLoc(padded[0], padded[1], padded[2])
}

// <init>(String language)
// <init>(String language, String country)
// <init>(String language, String country, String variant)
func localeInit(params []interface{}) interface{} {
	receiver, ok := params[0].(*object.Object)
	if !ok || object.IsNull(receiver) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Locale.<init>: this is null")
	}
	goStrings, errBlk := localeStringArguments("Locale.<init>", params[1:])
	if errBlk != nil {
		return errBlk
	}
	ghelpers.SetLocFields(receiver, localeFromLanguageCountryVariant(goStrings))
	return nil
}

// of(String language)
// of(String language, String country)
// of(String language, String country, String variant)
func localeOf(params []interface{}) interface{} {
	if len(params) < 1 || len(params) > 3 {
		// Cannot happen if the MethodSignatures table is right.
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException,
			fmt.Sprintf("Locale.of: unexpected parameter count %d", len(params)))
	}
	goStrings, errBlk := localeStringArguments("Locale.of", params)
	if errBlk != nil {
		return errBlk
	}
	return ghelpers.LocToObject(localeFromLanguageCountryVariant(goStrings))
}

// forLanguageTag(String languageTag)
func localeForLanguageTag(params []interface{}) interface{} {
	goStrings, errBlk := localeStringArguments("Locale.forLanguageTag", params[:1])
	if errBlk != nil {
		return errBlk
	}
	return ghelpers.LocToObject(ghelpers.ParseLoc(goStrings[0]))
}

// getInstance(String language, String country, String variant)
func localeGetInstanceLanguageCountryVariant(params []interface{}) interface{} {
	goStrings, errBlk := localeStringArguments("Locale.getInstance", params[:3])
	if errBlk != nil {
		return errBlk
	}
	return ghelpers.LocToObject(localeFromLanguageCountryVariant(goStrings))
}

// getInstance(String language, String script, String country, String variant, LocaleExtensions extensions)
func localeGetInstanceWithScript(params []interface{}) interface{} {
	goStrings, errBlk := localeStringArguments("Locale.getInstance", params[:4])
	if errBlk != nil {
		return errBlk
	}
	return ghelpers.LocToObject(ghelpers.NormalizeLocWithScript(goStrings[0], goStrings[1], goStrings[2], goStrings[3]))
}

// ---------------------------------------------------------------------------------------------
// Default locale
// ---------------------------------------------------------------------------------------------

// getDefault()
// Locale.getDefault() answers for the FORMAT category, as in the JDK.
func localeGetDefault([]interface{}) interface{} {
	return ghelpers.LocToObject(ghelpers.DefaultLoc(ghelpers.LocaleCategoryFormat))
}

// getDefault(Locale.Category category)
func localeGetDefaultForCategory(params []interface{}) interface{} {
	category, errBlk := ghelpers.LocCategoryFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	return ghelpers.LocToObject(ghelpers.DefaultLoc(category))
}

// setDefault(Locale newLocale): sets every category.
func localeSetDefault(params []interface{}) interface{} {
	loc, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	ghelpers.SetDefaultLoc(loc)
	return nil
}

// setDefault(Locale.Category category, Locale newLocale)
func localeSetDefaultForCategory(params []interface{}) interface{} {
	category, errBlk := ghelpers.LocCategoryFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	loc, errBlk := ghelpers.LocFromObject(params[1])
	if errBlk != nil {
		return errBlk
	}
	ghelpers.SetDefaultLocForCategory(category, loc)
	return nil
}

// ---------------------------------------------------------------------------------------------
// Instance methods
// ---------------------------------------------------------------------------------------------

// Build a gfunction that returns one component (language, script, country, or variant) of the receiver.
func localeComponentGetter(selectComponent func(loc ghelpers.Loc) string) func([]interface{}) interface{} {
	return func(params []interface{}) interface{} {
		loc, errBlk := ghelpers.LocFromObject(params[0])
		if errBlk != nil {
			return errBlk
		}
		return object.StringObjectFromGoString(selectComponent(loc))
	}
}

// toString()
func LocaleToString(params []interface{}) interface{} {
	return localeToString(params)
}

func localeToString(params []interface{}) interface{} {
	loc, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	return object.StringObjectFromGoString(loc.ToJavaLocaleString())
}

// toLanguageTag()
func localeToLanguageTag(params []interface{}) interface{} {
	loc, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	return object.StringObjectFromGoString(loc.ToLanguageTag())
}

// equals(Object other): false for null and for objects that are not Locales.
func localeEquals(params []interface{}) interface{} {
	receiverLoc, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	otherLoc, otherErrBlk := ghelpers.LocFromObject(params[1])
	if otherErrBlk != nil {
		return localeBooleanResult(false)
	}
	return localeBooleanResult(receiverLoc == otherLoc)
}

// hashCode()
func localeHashCode(params []interface{}) interface{} {
	loc, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	return int64(loc.JavaHashCode())
}

// clone()
func localeClone(params []interface{}) interface{} {
	loc, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	return ghelpers.LocToObject(loc)
}

// hasExtensions()
func localeHasExtensions(params []interface{}) interface{} {
	_, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	return localeBooleanResult(false)
}

// stripExtensions()
func localeStripExtensions(params []interface{}) interface{} {
	_, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	return params[0]
}

// Java boolean results travel as int64 on the operand stack.
func localeBooleanResult(value bool) interface{} {
	if value {
		return int64(1)
	}
	return int64(0)
}

// ---------------------------------------------------------------------------------------------
// Newly implemented Locale static and instance methods
// ---------------------------------------------------------------------------------------------

// getAvailableLocales()[Ljava/util/Locale;
func localeGetAvailableLocales([]interface{}) interface{} {
	locales := ghelpers.GetAvailableLocales()
	elements := make([]*object.Object, len(locales))
	for i, loc := range locales {
		elements[i] = ghelpers.LocToObject(loc)
	}
	return ghelpers.MakeObjectArray("[Ljava/util/Locale;", elements)
}

// getISOCountries()[Ljava/lang/String;
func localeGetISOCountries([]interface{}) interface{} {
	countries := ghelpers.GetISOCountries()
	elements := make([]*object.Object, len(countries))
	for i, country := range countries {
		elements[i] = object.StringObjectFromGoString(country)
	}
	return ghelpers.MakeObjectArray("[Ljava/lang/String;", elements)
}

// getISOCountries(Locale.IsoCountryCode type)Ljava/util/Set;
func localeGetISOCountriesForIsoCountryCode(params []interface{}) interface{} {
	if len(params) == 0 || params[0] == nil || params[0] == object.Null {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Locale.getISOCountries: isoCountryCode is null")
	}
	codeObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Locale.getISOCountries: invalid isoCountryCode")
	}

	name := ""
	if nameField, present := codeObj.FieldTable["name"]; present {
		if nameObj, isObj := nameField.Fvalue.(*object.Object); isObj && !object.IsNull(nameObj) {
			name = object.GoStringFromStringObject(nameObj)
		}
	}
	if name == "" {
		if ordField, present := codeObj.FieldTable["ordinal"]; present {
			if ord, isInt := ordField.Fvalue.(int64); isInt {
				switch ord {
				case 0:
					name = "PART1_ALPHA2"
				case 1:
					name = "PART1_ALPHA3"
				case 2:
					name = "PART3"
				}
			}
		}
	}

	var countries []string
	switch name {
	case "PART1_ALPHA2":
		countries = ghelpers.GetISOCountries()
	case "PART1_ALPHA3":
		countries = ghelpers.GetISOCountriesAlpha3()
	case "PART3":
		countries = []string{} // Part 3 is deleted/historical codes, can be empty or unmodifiable set
	default:
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Locale.getISOCountries: unsupported IsoCountryCode "+name)
	}

	return ghelpers.CreateStringSet(countries)
}

// getISOLanguages()[Ljava/lang/String;
func localeGetISOLanguages([]interface{}) interface{} {
	languages := ghelpers.GetISOLanguages()
	elements := make([]*object.Object, len(languages))
	for i, lang := range languages {
		elements[i] = object.StringObjectFromGoString(lang)
	}
	return ghelpers.MakeObjectArray("[Ljava/lang/String;", elements)
}

// getExtensionKeys()Ljava/util/Set;
func localeGetExtensionKeys(params []interface{}) interface{} {
	if _, errBlk := ghelpers.LocFromObject(params[0]); errBlk != nil {
		return errBlk
	}
	return ghelpers.CreateEmptySet()
}

// getUnicodeLocaleAttributes()Ljava/util/Set;
func localeGetUnicodeLocaleAttributes(params []interface{}) interface{} {
	if _, errBlk := ghelpers.LocFromObject(params[0]); errBlk != nil {
		return errBlk
	}
	return ghelpers.CreateEmptySet()
}

// getUnicodeLocaleKeys()Ljava/util/Set;
func localeGetUnicodeLocaleKeys(params []interface{}) interface{} {
	if _, errBlk := ghelpers.LocFromObject(params[0]); errBlk != nil {
		return errBlk
	}
	return ghelpers.CreateEmptySet()
}

// getExtension(C)Ljava/lang/String;
func localeGetExtension(params []interface{}) interface{} {
	if _, errBlk := ghelpers.LocFromObject(params[0]); errBlk != nil {
		return errBlk
	}
	return object.Null
}

// getUnicodeLocaleType(String key)Ljava/lang/String;
func localeGetUnicodeLocaleType(params []interface{}) interface{} {
	if _, errBlk := ghelpers.LocFromObject(params[0]); errBlk != nil {
		return errBlk
	}
	if len(params) < 2 || params[1] == nil || params[1] == object.Null {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Locale.getUnicodeLocaleType: key is null")
	}
	return object.Null
}

// getDisplayLanguage()Ljava/lang/String;
func localeGetDisplayLanguage(params []interface{}) interface{} {
	loc, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	inLocale := ghelpers.DefaultLoc(ghelpers.LocaleCategoryDisplay)
	return object.StringObjectFromGoString(ghelpers.GetDisplayLanguage(loc, inLocale))
}

// getDisplayLanguage(Locale inLocale)Ljava/lang/String;
func localeGetDisplayLanguageInLocale(params []interface{}) interface{} {
	loc, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	inLocale, inErr := ghelpers.LocFromObject(params[1])
	if inErr != nil {
		return inErr
	}
	return object.StringObjectFromGoString(ghelpers.GetDisplayLanguage(loc, inLocale))
}

// getDisplayScript()Ljava/lang/String;
func localeGetDisplayScript(params []interface{}) interface{} {
	loc, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	inLocale := ghelpers.DefaultLoc(ghelpers.LocaleCategoryDisplay)
	return object.StringObjectFromGoString(ghelpers.GetDisplayScript(loc, inLocale))
}

// getDisplayScript(Locale inLocale)Ljava/lang/String;
func localeGetDisplayScriptInLocale(params []interface{}) interface{} {
	loc, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	inLocale, inErr := ghelpers.LocFromObject(params[1])
	if inErr != nil {
		return inErr
	}
	return object.StringObjectFromGoString(ghelpers.GetDisplayScript(loc, inLocale))
}

// getDisplayCountry()Ljava/lang/String;
func localeGetDisplayCountry(params []interface{}) interface{} {
	loc, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	inLocale := ghelpers.DefaultLoc(ghelpers.LocaleCategoryDisplay)
	return object.StringObjectFromGoString(ghelpers.GetDisplayCountry(loc, inLocale))
}

// getDisplayCountry(Locale inLocale)Ljava/lang/String;
func localeGetDisplayCountryInLocale(params []interface{}) interface{} {
	loc, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	inLocale, inErr := ghelpers.LocFromObject(params[1])
	if inErr != nil {
		return inErr
	}
	return object.StringObjectFromGoString(ghelpers.GetDisplayCountry(loc, inLocale))
}

// getDisplayVariant()Ljava/lang/String;
func localeGetDisplayVariant(params []interface{}) interface{} {
	loc, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	inLocale := ghelpers.DefaultLoc(ghelpers.LocaleCategoryDisplay)
	return object.StringObjectFromGoString(ghelpers.GetDisplayVariant(loc, inLocale))
}

// getDisplayVariant(Locale inLocale)Ljava/lang/String;
func localeGetDisplayVariantInLocale(params []interface{}) interface{} {
	loc, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	inLocale, inErr := ghelpers.LocFromObject(params[1])
	if inErr != nil {
		return inErr
	}
	return object.StringObjectFromGoString(ghelpers.GetDisplayVariant(loc, inLocale))
}

// getDisplayName()Ljava/lang/String;
func localeGetDisplayName(params []interface{}) interface{} {
	loc, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	inLocale := ghelpers.DefaultLoc(ghelpers.LocaleCategoryDisplay)
	return object.StringObjectFromGoString(ghelpers.GetDisplayName(loc, inLocale))
}

// getDisplayName(Locale inLocale)Ljava/lang/String;
func localeGetDisplayNameInLocale(params []interface{}) interface{} {
	loc, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	inLocale, inErr := ghelpers.LocFromObject(params[1])
	if inErr != nil {
		return inErr
	}
	return object.StringObjectFromGoString(ghelpers.GetDisplayName(loc, inLocale))
}

// getISO3Language()Ljava/lang/String;
func localeGetISO3Language(params []interface{}) interface{} {
	loc, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	iso3, ok := ghelpers.GetISO3Language(loc.Language)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.MissingResourceException,
			fmt.Sprintf("Couldn't find 3-letter language code for %s", loc.Language))
	}
	return object.StringObjectFromGoString(iso3)
}

// getISO3Country()Ljava/lang/String;
func localeGetISO3Country(params []interface{}) interface{} {
	loc, errBlk := ghelpers.LocFromObject(params[0])
	if errBlk != nil {
		return errBlk
	}
	iso3, ok := ghelpers.GetISO3Country(loc.Country)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.MissingResourceException,
			fmt.Sprintf("Couldn't find 3-letter country code for %s", loc.Country))
	}
	return object.StringObjectFromGoString(iso3)
}
