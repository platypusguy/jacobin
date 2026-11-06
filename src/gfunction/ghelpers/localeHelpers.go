/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package ghelpers

import (
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
	"jacobin/src/util"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"unicode/utf16"
)

// Go-side locale core.
//
// Strategy: the Go type Loc is the source of truth for a locale. A java/util/Locale object is only a
// carrier for four String fields (language, script, country, variant). Gfunctions in java/util/Locale,
// java/time/*, java/text/*, etc. convert the Java object to a Loc once (LocFromObject) and then work
// exclusively with Go values.
//
// Target: Java 21 (CLDR 42 locale data).

// LocaleDefaultCategory mirrors java.util.Locale.Category: the PURPOSE for which a default locale is used.
// The ordinals match the Java enum (DISPLAY = 0, FORMAT = 1).
type LocaleDefaultCategory int

const (
	LocaleCategoryDisplay LocaleDefaultCategory = iota // default locale for user-interface text, such as display names
	LocaleCategoryFormat                               // default locale for formatting dates, times, numbers, and currencies
)

const localeDefaultCategoryCount = 2

// Names of the four String fields of a java/util/Locale object as represented in Jacobin.
const (
	LocFieldLanguage = "language"
	LocFieldScript   = "script"
	LocFieldCountry  = "country"
	LocFieldVariant  = "variant"
)

var (
	classNameHashSet = "java/util/HashSet"
	fieldNameMap     = "map"
)

// Loc is a normalized locale: language lower case, script title case, country upper case,
// variant exactly as given. The zero value is the root locale (Java: Locale.ROOT).
type Loc struct {
	Language string
	Script   string
	Country  string
	Variant  string
}

// RootLoc is the root locale (Java: Locale.ROOT).
var RootLoc = Loc{}

// NamedLoc associates a static field name of java/util/Locale with its locale.
type NamedLoc struct {
	FieldName string
	Value     Loc
}

// StandardLocConstants lists the public static Locale constants of Java 21, in declaration-neutral order.
var StandardLocConstants = []NamedLoc{
	{"ROOT", RootLoc},
	{"ENGLISH", Loc{Language: "en"}},
	{"FRENCH", Loc{Language: "fr"}},
	{"GERMAN", Loc{Language: "de"}},
	{"ITALIAN", Loc{Language: "it"}},
	{"JAPANESE", Loc{Language: "ja"}},
	{"KOREAN", Loc{Language: "ko"}},
	{"CHINESE", Loc{Language: "zh"}},
	{"SIMPLIFIED_CHINESE", Loc{Language: "zh", Country: "CN"}},
	{"TRADITIONAL_CHINESE", Loc{Language: "zh", Country: "TW"}},
	{"FRANCE", Loc{Language: "fr", Country: "FR"}},
	{"GERMANY", Loc{Language: "de", Country: "DE"}},
	{"ITALY", Loc{Language: "it", Country: "IT"}},
	{"JAPAN", Loc{Language: "ja", Country: "JP"}},
	{"KOREA", Loc{Language: "ko", Country: "KR"}},
	{"UK", Loc{Language: "en", Country: "GB"}},
	{"US", Loc{Language: "en", Country: "US"}},
	{"CANADA", Loc{Language: "en", Country: "CA"}},
	{"CANADA_FRENCH", Loc{Language: "fr", Country: "CA"}},
	{"CHINA", Loc{Language: "zh", Country: "CN"}},
	{"PRC", Loc{Language: "zh", Country: "CN"}},
	{"TAIWAN", Loc{Language: "zh", Country: "TW"}},
}

// Since Java 17, the legacy ISO 639 codes are converted to the current ones.
var legacyLanguageCodes = map[string]string{
	"iw": "he",
	"ji": "yi",
	"in": "id",
}

// ---------------------------------------------------------------------------------------------
// Construction, parsing, formatting
// ---------------------------------------------------------------------------------------------

// NormalizeLoc builds a Loc the way Locale.of(language, country, variant) does.
func NormalizeLoc(language, country, variant string) Loc {
	return NormalizeLocWithScript(language, "", country, variant)
}

// NormalizeLocWithScript is NormalizeLoc with an ISO 15924 script (e.g., "Hans").
func NormalizeLocWithScript(language, script, country, variant string) Loc {
	language = strings.ToLower(language)
	if replacement, found := legacyLanguageCodes[language]; found {
		language = replacement
	}
	return Loc{
		Language: language,
		Script:   titleCaseScript(script),
		Country:  strings.ToUpper(country),
		Variant:  variant,
	}
}

func titleCaseScript(script string) string {
	if script == "" {
		return script
	}
	return strings.ToUpper(script[:1]) + strings.ToLower(script[1:])
}

// ParseLoc accepts any of:
//
//	POSIX environment style   de_DE.UTF-8@euro   (encoding and modifier are discarded)
//	Locale.toString() style   en_US_POSIX, zh_CN_#Hans
//	BCP 47 language tag       zh-Hans-CN, en-US   ("und" is the root locale; extensions are ignored)
//
// C, POSIX, and the empty string yield RootLoc.
func ParseLoc(tag string) Loc {
	tag = strings.TrimSpace(tag)
	if cutIndex := strings.IndexAny(tag, ".@"); cutIndex >= 0 {
		tag = tag[:cutIndex]
	}
	if tag == "" || tag == "C" || tag == "POSIX" {
		return RootLoc
	}

	scriptAfterHash := ""
	if hashIndex := strings.Index(tag, "#"); hashIndex >= 0 {
		scriptAfterHash = tag[hashIndex+1:]
		tag = tag[:hashIndex]
	}

	separator := "_"
	if strings.Contains(tag, "-") {
		separator = "-"
	}
	parts := strings.Split(tag, separator)

	var parsed Loc
	parsed.Language = strings.ToLower(parts[0])
	if replacement, found := legacyLanguageCodes[parsed.Language]; found {
		parsed.Language = replacement
	}
	if parsed.Language == "und" {
		parsed.Language = ""
	}

	for _, part := range parts[1:] {
		if part == "" {
			continue
		}
		if len(part) == 1 {
			break // BCP 47 singleton: extensions or private use follow; not supported
		}
		noSubtagsYet := parsed.Script == "" && parsed.Country == "" && parsed.Variant == ""
		switch {
		case noSubtagsYet && len(part) == 4 && isAsciiLetters(part):
			parsed.Script = titleCaseScript(part)
		case parsed.Country == "" && parsed.Variant == "" &&
			((len(part) == 2 && isAsciiLetters(part)) || (len(part) == 3 && isAsciiDigits(part))):
			parsed.Country = strings.ToUpper(part)
		default:
			if parsed.Variant == "" {
				parsed.Variant = part
			} else {
				parsed.Variant += "_" + part
			}
		}
	}

	if scriptAfterHash != "" {
		parsed.Script = titleCaseScript(scriptAfterHash)
	}
	return parsed
}

func isAsciiLetters(text string) bool {
	if text == "" {
		return false
	}
	for _, character := range text {
		isLetter := (character >= 'a' && character <= 'z') || (character >= 'A' && character <= 'Z')
		if !isLetter {
			return false
		}
	}
	return true
}

func isAsciiDigits(text string) bool {
	if text == "" {
		return false
	}
	for _, character := range text {
		if character < '0' || character > '9' {
			return false
		}
	}
	return true
}

// IsRootLoc reports whether the locale is the root locale.
func (loc Loc) IsRootLoc() bool {
	return loc == RootLoc
}

// ToJavaLocaleString renders the locale exactly as Locale.toString() does:
// en, en_US, en__POSIX, _US, zh_CN_#Hans, zh__#Hans. The root locale yields "".
func (loc Loc) ToJavaLocaleString() string {
	hasLanguage := loc.Language != ""
	hasScript := loc.Script != ""
	hasCountry := loc.Country != ""
	hasVariant := loc.Variant != ""

	text := loc.Language
	if hasCountry || (hasLanguage && (hasVariant || hasScript)) {
		text += "_" + loc.Country
	}
	if hasVariant && (hasLanguage || hasCountry) {
		text += "_" + loc.Variant
	}
	if hasScript && (hasLanguage || hasCountry) {
		text += "_#" + loc.Script
	}
	return text
}

// ToLanguageTag renders the locale as Locale.toLanguageTag() does (BCP 47): en-US, zh-Hans-CN, und.
// Variants that are not well-formed BCP 47 variants go into a private-use "-x-lvariant-" section.
func (loc Loc) ToLanguageTag() string {
	subtags := make([]string, 0, 5)
	if loc.Language == "" {
		subtags = append(subtags, "und")
	} else {
		subtags = append(subtags, loc.Language)
	}
	if loc.Script != "" {
		subtags = append(subtags, loc.Script)
	}
	if loc.Country != "" {
		subtags = append(subtags, loc.Country)
	}

	var invalidVariants []string
	if loc.Variant != "" {
		for _, variantPart := range strings.Split(loc.Variant, "_") {
			if isWellFormedBcp47Variant(variantPart) {
				subtags = append(subtags, variantPart)
			} else if variantPart != "" {
				invalidVariants = append(invalidVariants, variantPart)
			}
		}
	}
	if len(invalidVariants) > 0 {
		subtags = append(subtags, "x", "lvariant")
		subtags = append(subtags, invalidVariants...)
	}
	return strings.Join(subtags, "-")
}

func isWellFormedBcp47Variant(variant string) bool {
	if !isAsciiAlphanumeric(variant) {
		return false
	}
	if len(variant) >= 5 && len(variant) <= 8 {
		return true
	}
	return len(variant) == 4 && variant[0] >= '0' && variant[0] <= '9'
}

func isAsciiAlphanumeric(text string) bool {
	if text == "" {
		return false
	}
	for _, character := range text {
		isLetter := (character >= 'a' && character <= 'z') || (character >= 'A' && character <= 'Z')
		isDigit := character >= '0' && character <= '9'
		if !isLetter && !isDigit {
			return false
		}
	}
	return true
}

// CandidateLocs returns the lookup chain for locale-sensitive data, most specific first, ending with RootLoc,
// without duplicates. This is a deliberate simplification: CLDR also defines explicit parent locales
// (for example en_GB -> en_001 -> en, and zh_Hant_HK -> zh_Hant). The CLDR data generator should emit that
// parentLocales table, and this function should consult it before falling back to truncation.
func (loc Loc) CandidateLocs() []Loc {
	candidates := make([]Loc, 0, 6)
	appendIfNew := func(candidate Loc) {
		for _, existing := range candidates {
			if existing == candidate {
				return
			}
		}
		candidates = append(candidates, candidate)
	}
	appendIfNew(loc)
	appendIfNew(Loc{Language: loc.Language, Script: loc.Script, Country: loc.Country})
	appendIfNew(Loc{Language: loc.Language, Country: loc.Country})
	appendIfNew(Loc{Language: loc.Language, Script: loc.Script})
	appendIfNew(Loc{Language: loc.Language})
	appendIfNew(RootLoc)
	return candidates
}

// JavaHashCode reproduces Locale.hashCode() for a locale without extensions:
// 31-based hash over language, script, region, and variant (BaseLocale.hashCode).
func (loc Loc) JavaHashCode() int32 {
	var hash int32
	for _, component := range []string{loc.Language, loc.Script, loc.Country, loc.Variant} {
		hash = 31*hash + javaStringHashCode(component)
	}
	return hash
}

// javaStringHashCode reproduces String.hashCode(), which works on UTF-16 code units.
func javaStringHashCode(text string) int32 {
	var hash int32
	for _, codeUnit := range utf16.Encode([]rune(text)) {
		hash = 31*hash + int32(codeUnit)
	}
	return hash
}

// ---------------------------------------------------------------------------------------------
// Default locales (one per LocaleDefaultCategory), safe for concurrent use
// ---------------------------------------------------------------------------------------------

var (
	defaultLocMutex          sync.RWMutex
	defaultLocByCategory     [localeDefaultCategoryCount]Loc
	defaultLocInitialized    bool
	locPlatformDefaultSource func() string
)

// SetLocPlatformDefaultSource registers the platform-specific provider of the initial default locale
// (for example, a Windows API wrapper). The provider returns text that ParseLoc understands.
// Without a provider, the LC_ALL, LC_MESSAGES, and LANG environment variables are consulted.
// Call this before the first use of DefaultLoc.
func SetLocPlatformDefaultSource(source func() string) {
	defaultLocMutex.Lock()
	defer defaultLocMutex.Unlock()
	locPlatformDefaultSource = source
}

// DefaultLoc returns the default locale for the category (Java: Locale.getDefault(Category)).
func DefaultLoc(category LocaleDefaultCategory) Loc {
	category = clampLocaleDefaultCategory(category)

	defaultLocMutex.RLock()
	if defaultLocInitialized {
		found := defaultLocByCategory[category]
		defaultLocMutex.RUnlock()
		return found
	}
	defaultLocMutex.RUnlock()

	defaultLocMutex.Lock()
	defer defaultLocMutex.Unlock()
	initializeDefaultLocsLocked()
	return defaultLocByCategory[category]
}

// SetDefaultLoc sets the default locale for every category (Java: Locale.setDefault(Locale)).
func SetDefaultLoc(loc Loc) {
	defaultLocMutex.Lock()
	defer defaultLocMutex.Unlock()
	for index := range defaultLocByCategory {
		defaultLocByCategory[index] = loc
	}
	defaultLocInitialized = true
}

// SetDefaultLocForCategory sets the default locale for one category
// (Java: Locale.setDefault(Locale.Category, Locale)).
func SetDefaultLocForCategory(category LocaleDefaultCategory, loc Loc) {
	category = clampLocaleDefaultCategory(category)
	defaultLocMutex.Lock()
	defer defaultLocMutex.Unlock()
	initializeDefaultLocsLocked()
	defaultLocByCategory[category] = loc
}

func clampLocaleDefaultCategory(category LocaleDefaultCategory) LocaleDefaultCategory {
	if category < 0 || category >= localeDefaultCategoryCount {
		return LocaleCategoryFormat
	}
	return category
}

// initializeDefaultLocsLocked requires defaultLocMutex to be held for writing.
func initializeDefaultLocsLocked() {
	if defaultLocInitialized {
		return
	}
	initial := Loc{Language: "en", Country: "US"}
	if locPlatformDefaultSource != nil {
		if fromPlatform := ParseLoc(locPlatformDefaultSource()); !fromPlatform.IsRootLoc() {
			initial = fromPlatform
		}
	} else if fromEnvironment, found := defaultLocFromEnvironment(); found {
		initial = fromEnvironment
	}
	for index := range defaultLocByCategory {
		defaultLocByCategory[index] = initial
	}
	defaultLocInitialized = true
}

func defaultLocFromEnvironment() (Loc, bool) {
	for _, variableName := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if value := os.Getenv(variableName); value != "" {
			if parsed := ParseLoc(value); !parsed.IsRootLoc() {
				return parsed, true
			}
		}
	}
	return RootLoc, false
}

// ---------------------------------------------------------------------------------------------
// Conversion between Loc and java/util/Locale objects
// ---------------------------------------------------------------------------------------------

// SetLocFields stores the four locale components as String fields of a java/util/Locale object.
func SetLocFields(localeObj *object.Object, loc Loc) {
	storeStringField := func(fieldName, value string) {
		localeObj.FieldTable[fieldName] = object.Field{
			Ftype:  types.Ref,
			Fvalue: object.StringObjectFromGoString(value),
		}
	}
	storeStringField(LocFieldLanguage, loc.Language)
	storeStringField(LocFieldScript, loc.Script)
	storeStringField(LocFieldCountry, loc.Country)
	storeStringField(LocFieldVariant, loc.Variant)
}

// LocToObject creates a new java/util/Locale object for the locale.
func LocToObject(loc Loc) *object.Object {
	localeObj := object.MakeEmptyObjectWithClassName(&types.ClassNameLocale)
	SetLocFields(localeObj, loc)
	return localeObj
}

// LocFromObject extracts the Loc from a java/util/Locale object. Use this in every gfunction that takes
// a Locale parameter or receiver. A null reference yields NullPointerException; an object without the
// locale fields yields IllegalArgumentException.
func LocFromObject(value interface{}) (Loc, *GErrBlk) {
	localeObj, ok := value.(*object.Object)
	if !ok || object.IsNull(localeObj) {
		return RootLoc, GetGErrBlk(excNames.NullPointerException, "LocFromObject: Locale reference is null")
	}
	if _, present := localeObj.FieldTable[LocFieldLanguage]; !present {
		return RootLoc, GetGErrBlk(excNames.IllegalArgumentException,
			"LocFromObject: object is not a Locale (no language field)")
	}

	readStringField := func(fieldName string) string {
		stringObj, isObject := localeObj.FieldTable[fieldName].Fvalue.(*object.Object)
		if !isObject || object.IsNull(stringObj) {
			return ""
		}
		return object.GoStringFromStringObject(stringObj)
	}
	return Loc{
		Language: readStringField(LocFieldLanguage),
		Script:   readStringField(LocFieldScript),
		Country:  readStringField(LocFieldCountry),
		Variant:  readStringField(LocFieldVariant),
	}, nil
}

// LocCategoryFromObject converts a java/util/Locale$Category enum object into a LocaleDefaultCategory.
// It looks at the enum's "name" field first and falls back to its "ordinal" field.
func LocCategoryFromObject(value interface{}) (LocaleDefaultCategory, *GErrBlk) {
	categoryObj, ok := value.(*object.Object)
	if !ok || object.IsNull(categoryObj) {
		return LocaleCategoryFormat, GetGErrBlk(excNames.NullPointerException,
			"LocCategoryFromObject: Locale.Category reference is null")
	}

	if nameField, present := categoryObj.FieldTable["name"]; present {
		if nameObj, isObject := nameField.Fvalue.(*object.Object); isObject && !object.IsNull(nameObj) {
			switch object.GoStringFromStringObject(nameObj) {
			case "DISPLAY":
				return LocaleCategoryDisplay, nil
			case "FORMAT":
				return LocaleCategoryFormat, nil
			}
		}
	}
	if ordinalField, present := categoryObj.FieldTable["ordinal"]; present {
		if ordinal, isInt := ordinalField.Fvalue.(int64); isInt {
			switch ordinal {
			case int64(LocaleCategoryDisplay):
				return LocaleCategoryDisplay, nil
			case int64(LocaleCategoryFormat):
				return LocaleCategoryFormat, nil
			}
		}
	}
	return LocaleCategoryFormat, GetGErrBlk(excNames.IllegalArgumentException,
		"LocCategoryFromObject: unrecognized Locale.Category object")
}

// ---------------------------------------------------------------------------------------------
// ISO codes and Display names
// ---------------------------------------------------------------------------------------------

var (
	isoCountries2 = []string{
		"AD", "AE", "AF", "AG", "AI", "AL", "AM", "AO", "AQ", "AR", "AS", "AT", "AU", "AW", "AX", "AZ",
		"BA", "BB", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BL", "BM", "BN", "BO", "BQ", "BR", "BS",
		"BT", "BV", "BW", "BY", "BZ", "CA", "CC", "CD", "CF", "CG", "CH", "CI", "CK", "CL", "CM", "CN",
		"CO", "CR", "CU", "CV", "CW", "CX", "CY", "CZ", "DE", "DJ", "DK", "DM", "DO", "DZ", "EC", "EE",
		"EG", "EH", "ER", "ES", "ET", "FI", "FJ", "FK", "FM", "FO", "FR", "GA", "GB", "GD", "GE", "GF",
		"GG", "GH", "GI", "GL", "GM", "GN", "GP", "GQ", "GR", "GS", "GT", "GU", "GW", "GY", "HK", "HM",
		"HN", "HR", "HT", "HU", "ID", "IE", "IL", "IM", "IN", "IO", "IQ", "IR", "IS", "IT", "JE", "JM",
		"JO", "JP", "KE", "KG", "KH", "KI", "KM", "KN", "KP", "KR", "KW", "KY", "KZ", "LA", "LB", "LC",
		"LI", "LK", "LR", "LS", "LT", "LU", "LV", "LY", "MA", "MC", "MD", "ME", "MF", "MG", "MH", "MK",
		"ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS", "MT", "MU", "MV", "MW", "MX", "MY", "MZ", "NA",
		"NC", "NE", "NF", "NG", "NI", "NL", "NO", "NP", "NR", "NU", "NZ", "OM", "PA", "PE", "PF", "PG",
		"PH", "PK", "PL", "PM", "PN", "PR", "PS", "PT", "PW", "PY", "QA", "RE", "RO", "RS", "RU", "RW",
		"SA", "SB", "SC", "SD", "SE", "SG", "SH", "SI", "SJ", "SK", "SL", "SM", "SN", "SO", "SR", "SS",
		"ST", "SV", "SX", "SY", "SZ", "TC", "TD", "TF", "TG", "TH", "TJ", "TK", "TL", "TM", "TN", "TO",
		"TR", "TT", "TV", "TW", "TZ", "UA", "UG", "UM", "US", "UY", "UZ", "VA", "VC", "VE", "VG", "VI",
		"VN", "VU", "WF", "WS", "YE", "YT", "ZA", "ZM", "ZW",
	}

	isoCountries3 = []string{
		"ABW", "AFG", "AGO", "AIA", "ALA", "ALB", "AND", "ARE", "ARG", "ARM", "ASM", "ATA", "ATF", "ATG",
		"AUS", "AUT", "AZE", "BDI", "BEL", "BEN", "BES", "BFA", "BGD", "BGR", "BHR", "BHS", "BIH", "BLM",
		"BLR", "BLZ", "BMU", "BOL", "BRA", "BRB", "BRN", "BTN", "BVT", "BWA", "CAF", "CAN", "CCK", "CHE",
		"CHL", "CHN", "CIV", "CMR", "COD", "COG", "COK", "COL", "COM", "CPV", "CRI", "CUB", "CUW", "CXR",
		"CYM", "CYP", "CZE", "DEU", "DJI", "DMA", "DNK", "DOM", "DZA", "ECU", "EGY", "ERI", "ESH", "ESP",
		"EST", "ETH", "FIN", "FJI", "FLK", "FRA", "FRO", "FSM", "GAB", "GBR", "GEO", "GGY", "GHA", "GIB",
		"GIN", "GLP", "GMB", "GNB", "GNQ", "GRC", "GRD", "GRL", "GTM", "GUF", "GUM", "GUY", "HKG", "HMD",
		"HND", "HRV", "HTI", "HUN", "IDN", "IMN", "IND", "IOT", "IRL", "IRN", "IRQ", "ISL", "ISR", "ITA",
		"JAM", "JEY", "JOR", "JPN", "KAZ", "KEN", "KGZ", "KHM", "KIR", "KNA", "KOR", "KWT", "LAO", "LBN",
		"LBR", "LBY", "LCA", "LIE", "LKA", "LSO", "LTU", "LUX", "LVA", "MAC", "MAF", "MAR", "MCO", "MDA",
		"MDG", "MDV", "MEX", "MHL", "MKD", "MLI", "MLT", "MMR", "MNE", "MNG", "MNP", "MOZ", "MRT", "MSR",
		"MTQ", "MUS", "MWI", "MYS", "MYT", "NAM", "NCL", "NER", "NFK", "NGA", "NIC", "NIU", "NLD", "NOR",
		"NPL", "NRU", "NZL", "OMN", "PAK", "PAN", "PCN", "PER", "PHL", "PLW", "PNG", "POL", "PRI", "PRK",
		"PRT", "PRY", "PSE", "PYF", "QAT", "REU", "ROU", "RUS", "RWA", "SAU", "SDN", "SEN", "SGP", "SGS",
		"SHN", "SJM", "SLB", "SLE", "SLV", "SMR", "SOM", "SPM", "SRB", "SSD", "STP", "SUR", "SVK", "SVN",
		"SWE", "SWZ", "SXM", "SYC", "SYR", "TCA", "TCD", "TGO", "THA", "TJK", "TKL", "TKM", "TLS", "TON",
		"TTO", "TUN", "TUR", "TUV", "TWN", "TZA", "UGA", "UKR", "UMI", "URY", "USA", "UZB", "VAT", "VCT",
		"VEN", "VGB", "VIR", "VNM", "VUT", "WLF", "WSM", "YEM", "ZAF", "ZMB", "ZWE",
	}

	isoLanguages2 = []string{
		"aa", "ab", "ae", "af", "ak", "am", "an", "ar", "as", "av", "ay", "az", "ba", "be", "bg", "bh",
		"bi", "bm", "bn", "bo", "br", "bs", "ca", "ce", "ch", "co", "cr", "cs", "cu", "cv", "cy", "da",
		"de", "dv", "dz", "ee", "el", "en", "eo", "es", "et", "eu", "fa", "ff", "fi", "fj", "fo", "fr",
		"fy", "ga", "gd", "gl", "gn", "gu", "gv", "ha", "he", "hi", "ho", "hr", "ht", "hu", "hy", "hz",
		"ia", "id", "ie", "ig", "ii", "ik", "io", "is", "it", "iu", "ja", "jv", "ka", "kg", "ki", "kj",
		"kk", "kl", "km", "kn", "ko", "kr", "ks", "ku", "kv", "kw", "ky", "la", "lb", "lg", "li", "ln",
		"lo", "lt", "lu", "lv", "mg", "mh", "mi", "mk", "ml", "mn", "mr", "ms", "mt", "my", "na", "nb",
		"nd", "ne", "ng", "nl", "nn", "no", "nr", "nv", "ny", "oc", "oj", "om", "or", "os", "pa", "pi",
		"pl", "ps", "pt", "qu", "rm", "rn", "ro", "ru", "rw", "sa", "sc", "sd", "se", "sg", "si", "sk",
		"sl", "sm", "sn", "so", "sq", "sr", "ss", "st", "su", "sv", "sw", "ta", "te", "tg", "th", "ti",
		"tk", "tl", "tn", "to", "tr", "ts", "tt", "tw", "ty", "ug", "uk", "ur", "uz", "ve", "vi", "vo",
		"wa", "wo", "xh", "yi", "yo", "za", "zh", "zu",
	}

	iso2ToIso3CountryMap = map[string]string{
		"AD": "AND", "AE": "ARE", "AF": "AFG", "AG": "ATG", "AI": "AIA", "AL": "ALB", "AM": "ARM", "AO": "AGO",
		"AQ": "ATA", "AR": "ARG", "AS": "ASM", "AT": "AUT", "AU": "AUS", "AW": "ABW", "AX": "ALA", "AZ": "AZE",
		"BA": "BIH", "BB": "BRB", "BD": "BGD", "BE": "BEL", "BF": "BFA", "BG": "BGR", "BH": "BHR", "BI": "BDI",
		"BJ": "BEN", "BL": "BLM", "BM": "BMU", "BN": "BRN", "BO": "BOL", "BQ": "BES", "BR": "BRA", "BS": "BHS",
		"BT": "BTN", "BV": "BVT", "BW": "BWA", "BY": "BLR", "BZ": "BLZ", "CA": "CAN", "CC": "CCK", "CD": "COD",
		"CF": "CAF", "CG": "COG", "CH": "CHE", "CI": "CIV", "CK": "COK", "CL": "CHL", "CM": "CMR", "CN": "CHN",
		"CO": "COL", "CR": "CRI", "CU": "CUB", "CV": "CPV", "CW": "CUW", "CX": "CXR", "CY": "CYP", "CZ": "CZE",
		"DE": "DEU", "DJ": "DJI", "DK": "DNK", "DM": "DMA", "DO": "DOM", "DZ": "DZA", "EC": "ECU", "EE": "EST",
		"EG": "EGY", "EH": "ESH", "ER": "ERI", "ES": "ESP", "ET": "ETH", "FI": "FIN", "FJ": "FJI", "FK": "FLK",
		"FM": "FSM", "FO": "FRO", "FR": "FRA", "GA": "GAB", "GB": "GBR", "GD": "GRD", "GE": "GEO", "GF": "GUF",
		"GG": "GGY", "GH": "GHA", "GI": "GIB", "GL": "GRL", "GM": "GMB", "GN": "GIN", "GP": "GLP", "GQ": "GNQ",
		"GR": "GRC", "GS": "SGS", "GT": "GTM", "GU": "GUM", "GW": "GNB", "GY": "GUY", "HK": "HKG", "HM": "HMD",
		"HN": "HND", "HR": "HRV", "HT": "HTI", "HU": "HUN", "ID": "IDN", "IE": "IRL", "IL": "ISR", "IM": "IMN",
		"IN": "IND", "IO": "IOT", "IQ": "IRQ", "IR": "IRN", "IS": "ISL", "IT": "ITA", "JE": "JEY", "JM": "JAM",
		"JO": "JOR", "JP": "JPN", "KE": "KEN", "KG": "KGZ", "KH": "KHM", "KI": "KIR", "KM": "COM", "KN": "KNA",
		"KP": "PRK", "KR": "KOR", "KW": "KWT", "KY": "CYM", "KZ": "KAZ", "LA": "LAO", "LB": "LBN", "LC": "LCA",
		"LI": "LIE", "LK": "LKA", "LR": "LBR", "LS": "LSO", "LT": "LTU", "LU": "LUX", "LV": "LVA", "LY": "LBY",
		"MA": "MAR", "MC": "MCO", "MD": "MDA", "ME": "MNE", "MF": "MAF", "MG": "MDG", "MH": "MHL", "MK": "MKD",
		"ML": "MLI", "MM": "MMR", "MN": "MNG", "MO": "MAC", "MP": "MNP", "MQ": "MTQ", "MR": "MRT", "MS": "MSR",
		"MT": "MLT", "MU": "MUS", "MV": "MDV", "MW": "MWI", "MX": "MEX", "MY": "MYS", "MZ": "MOZ", "NA": "NAM",
		"NC": "NCL", "NE": "NER", "NF": "NFK", "NG": "NGA", "NI": "NIC", "NL": "NLD", "NO": "NOR", "NP": "NPL",
		"NR": "NRU", "NU": "NIU", "NZ": "NZL", "OM": "OMN", "PA": "PAN", "PE": "PER", "PF": "PYF", "PG": "PNG",
		"PH": "PHL", "PK": "PAK", "PL": "POL", "PM": "SPM", "PN": "PCN", "PR": "PRI", "PS": "PSE", "PT": "PRT",
		"PW": "PLW", "PY": "PRY", "QA": "QAT", "RE": "REU", "RO": "ROU", "RS": "SRB", "RU": "RUS", "RW": "RWA",
		"SA": "SAU", "SB": "SLB", "SC": "SYC", "SD": "SDN", "SE": "SWE", "SG": "SGP", "SH": "SHN", "SI": "SVN",
		"SJ": "SJM", "SK": "SVK", "SL": "SLE", "SM": "SMR", "SN": "SEN", "SO": "SOM", "SR": "SUR", "SS": "SSD",
		"ST": "STP", "SV": "SLV", "SX": "SXM", "SY": "SYR", "SZ": "SWZ", "TC": "TCA", "TD": "TCD", "TF": "ATF",
		"TG": "TGO", "TH": "THA", "TJ": "TJK", "TK": "TKL", "TL": "TLS", "TM": "TKM", "TN": "TUN", "TO": "TON",
		"TR": "TUR", "TT": "TTO", "TV": "TUV", "TW": "TWN", "TZ": "TZA", "UA": "UKR", "UG": "UGA", "UM": "UMI",
		"US": "USA", "UY": "URY", "UZ": "UZB", "VA": "VAT", "VC": "VCT", "VE": "VEN", "VG": "VGB", "VI": "VIR",
		"VN": "VNM", "VU": "VUT", "WF": "WLF", "WS": "WSM", "YE": "YEM", "YT": "MYT", "ZA": "ZAF", "ZM": "ZMB",
		"ZW": "ZWE",
	}

	iso2ToIso3LanguageMap = map[string]string{
		"aa": "aar", "ab": "abk", "ae": "ave", "af": "afr", "ak": "aka", "am": "amh", "an": "arg", "ar": "ara",
		"as": "asm", "av": "ava", "ay": "aym", "az": "aze", "ba": "bak", "be": "bel", "bg": "bul", "bh": "bih",
		"bi": "bis", "bm": "bam", "bn": "ben", "bo": "bod", "br": "bre", "bs": "bos", "ca": "cat", "ce": "che",
		"ch": "cha", "co": "cos", "cr": "cre", "cs": "ces", "cu": "chu", "cv": "chv", "cy": "cym", "da": "dan",
		"de": "deu", "dv": "div", "dz": "dzo", "ee": "ewe", "el": "ell", "en": "eng", "eo": "epo", "es": "spa",
		"et": "est", "eu": "eus", "fa": "fas", "ff": "ful", "fi": "fin", "fj": "fij", "fo": "fao", "fr": "fra",
		"fy": "fry", "ga": "gle", "gd": "gla", "gl": "glg", "gn": "grn", "gu": "guj", "gv": "glv", "ha": "hau",
		"he": "heb", "hi": "hin", "ho": "hmo", "hr": "hrv", "ht": "hat", "hu": "hun", "hy": "hye", "hz": "her",
		"ia": "ina", "id": "ind", "ie": "ile", "ig": "ibo", "ii": "iii", "ik": "ipk", "io": "ido", "is": "isl",
		"it": "ita", "iu": "iku", "ja": "jpn", "jv": "jav", "ka": "kat", "kg": "kon", "ki": "kik", "kj": "kua",
		"kk": "kaz", "kl": "kal", "km": "khm", "kn": "kan", "ko": "kor", "kr": "kau", "ks": "kas", "ku": "kur",
		"kv": "kom", "kw": "cor", "ky": "kir", "la": "lat", "lb": "ltz", "lg": "lug", "li": "lim", "ln": "lin",
		"lo": "lao", "lt": "lit", "lu": "lub", "lv": "lav", "mg": "mlg", "mh": "mah", "mi": "mri", "mk": "mkd",
		"ml": "mal", "mn": "mon", "mr": "mar", "ms": "msa", "mt": "mlt", "my": "mya", "na": "nau", "nb": "nob",
		"nd": "nde", "ne": "nep", "ng": "ndo", "nl": "nld", "nn": "nno", "no": "nor", "nr": "nbl", "nv": "nav",
		"ny": "nya", "oc": "oci", "oj": "oji", "om": "orm", "or": "ori", "os": "oss", "pa": "pan", "pi": "pli",
		"pl": "pol", "ps": "pus", "pt": "por", "qu": "que", "rm": "roh", "rn": "run", "ro": "ron", "ru": "rus",
		"rw": "kin", "sa": "san", "sc": "srd", "sd": "snd", "se": "sme", "sg": "sag", "si": "sin", "sk": "slk",
		"sl": "slv", "sm": "smo", "sn": "sna", "so": "som", "sq": "sqi", "sr": "srp", "ss": "ssw", "st": "sot",
		"su": "sun", "sv": "swe", "sw": "swa", "ta": "tam", "te": "tel", "tg": "tgk", "th": "tha", "ti": "tir",
		"tk": "tuk", "tl": "tgl", "tn": "tsn", "to": "ton", "tr": "tur", "ts": "tso", "tt": "tat", "tw": "twi",
		"ty": "tah", "ug": "uig", "uk": "ukr", "ur": "urd", "uz": "uzb", "ve": "ven", "vi": "vie", "vo": "vol",
		"wa": "wln", "wo": "wol", "xh": "xho", "yi": "yid", "yo": "yor", "za": "zha", "zh": "zho", "zu": "zul",
		"in": "ind", "iw": "heb", "ji": "yid",
	}

	displayLanguageNamesEN = map[string]string{
		"en": "English", "fr": "French", "de": "German", "it": "Italian", "ja": "Japanese",
		"ko": "Korean", "zh": "Chinese", "es": "Spanish", "ru": "Russian", "pt": "Portuguese",
		"ar": "Arabic", "he": "Hebrew", "hi": "Hindi", "id": "Indonesian", "nl": "Dutch",
		"pl": "Polish", "sv": "Swedish", "tr": "Turkish", "vi": "Vietnamese", "th": "Thai",
		"cs": "Czech", "el": "Greek", "da": "Danish", "fi": "Finnish", "no": "Norwegian",
		"hu": "Hungarian", "ro": "Romanian", "uk": "Ukrainian", "fa": "Persian", "bg": "Bulgarian",
		"hr": "Croatian", "sk": "Slovak", "sl": "Slovenian", "lt": "Lithuanian", "lv": "Latvian",
		"et": "Estonian", "ms": "Malay", "ca": "Catalan", "sr": "Serbian", "is": "Icelandic",
	}

	displayCountryNamesEN = map[string]string{
		"US": "United States", "GB": "United Kingdom", "CA": "Canada", "FR": "France", "DE": "Germany",
		"IT": "Italy", "JP": "Japan", "KR": "South Korea", "CN": "China", "TW": "Taiwan",
		"ES": "Spain", "RU": "Russia", "BR": "Brazil", "MX": "Mexico", "AU": "Australia",
		"IN": "India", "IL": "Israel", "ID": "Indonesia", "NL": "Netherlands", "PL": "Poland",
		"SE": "Sweden", "TR": "Turkey", "VN": "Vietnam", "TH": "Thailand", "CH": "Switzerland",
		"AT": "Austria", "BE": "Belgium", "CZ": "Czechia", "DK": "Denmark", "FI": "Finland",
		"GR": "Greece", "HU": "Hungary", "IE": "Ireland", "NO": "Norway", "NZ": "New Zealand",
		"PT": "Portugal", "RO": "Romania", "SA": "Saudi Arabia", "SG": "Singapore", "ZA": "South Africa",
		"UA": "Ukraine", "HK": "Hong Kong", "EG": "Egypt", "AR": "Argentina", "CL": "Chile",
		"CO": "Colombia", "PE": "Peru", "PH": "Philippines", "PK": "Pakistan", "NG": "Nigeria",
	}

	displayScriptNamesEN = map[string]string{
		"Latn": "Latin", "Hans": "Simplified", "Hant": "Traditional", "Cyrl": "Cyrillic",
		"Arab": "Arabic", "Deva": "Devanagari", "Hebr": "Hebrew", "Jpan": "Japanese",
		"Kore": "Korean", "Grek": "Greek", "Thai": "Thai",
	}
)

func GetISOCountries() []string {
	return slices.Clone(isoCountries2)
}

func GetISOCountriesAlpha3() []string {
	return slices.Clone(isoCountries3)
}

func GetISOLanguages() []string {
	return slices.Clone(isoLanguages2)
}

func GetAvailableLocales() []Loc {
	locs := make([]Loc, len(StandardLocConstants))
	for i, c := range StandardLocConstants {
		locs[i] = c.Value
	}
	return locs
}

func GetISO3Language(lang string) (string, bool) {
	lang = strings.ToLower(lang)
	if lang == "" {
		return "", true
	}
	if len(lang) == 3 {
		return lang, true
	}
	code, ok := iso2ToIso3LanguageMap[lang]
	return code, ok
}

func GetISO3Country(country string) (string, bool) {
	country = strings.ToUpper(country)
	if country == "" {
		return "", true
	}
	if len(country) == 3 {
		return country, true
	}
	code, ok := iso2ToIso3CountryMap[country]
	return code, ok
}

func GetDisplayLanguage(loc Loc, inLocale Loc) string {
	if loc.Language == "" {
		return ""
	}
	if name, ok := displayLanguageNamesEN[loc.Language]; ok {
		return name
	}
	return loc.Language
}

func GetDisplayScript(loc Loc, inLocale Loc) string {
	if loc.Script == "" {
		return ""
	}
	if name, ok := displayScriptNamesEN[loc.Script]; ok {
		return name
	}
	return loc.Script
}

func GetDisplayCountry(loc Loc, inLocale Loc) string {
	if loc.Country == "" {
		return ""
	}
	if name, ok := displayCountryNamesEN[loc.Country]; ok {
		return name
	}
	return loc.Country
}

func GetDisplayVariant(loc Loc, inLocale Loc) string {
	return loc.Variant
}

func GetDisplayName(loc Loc, inLocale Loc) string {
	dispLang := GetDisplayLanguage(loc, inLocale)
	dispScript := GetDisplayScript(loc, inLocale)
	dispCountry := GetDisplayCountry(loc, inLocale)
	dispVar := GetDisplayVariant(loc, inLocale)

	var details []string
	if dispScript != "" {
		details = append(details, dispScript)
	}
	if dispCountry != "" {
		details = append(details, dispCountry)
	}
	if dispVar != "" {
		details = append(details, dispVar)
	}

	if len(details) == 0 {
		return dispLang
	}
	if dispLang == "" {
		return "(" + strings.Join(details, ", ") + ")"
	}
	return dispLang + " (" + strings.Join(details, ", ") + ")"
}

// MakeObjectArray wraps an array of object pointers into a JVM array object.
func MakeObjectArray(elemClassName string, elements []*object.Object) *object.Object {
	className := elemClassName
	if !strings.HasPrefix(className, "[") {
		className = "[L" + className + ";"
	}
	arrObj := object.MakeEmptyObjectWithClassName(&className)
	arrObj.FieldTable["value"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: elements,
	}
	return arrObj
}

// CreateStringSet creates a java/util/HashSet containing the given Go strings.
func CreateStringSet(items []string) *object.Object {
	hs := object.MakeEmptyObjectWithClassName(&classNameHashSet)
	hs.FieldTable[fieldNameMap] = object.Field{
		Ftype:  types.Ref,
		Fvalue: make(types.DefHashMap),
	}
	hm := hs.FieldTable[fieldNameMap].Fvalue.(types.DefHashMap)

	for _, item := range items {
		strObj := object.StringObjectFromGoString(item)
		if fld, ok := strObj.FieldTable["value"]; ok {
			if hashedUint64, err := util.HashAnything(fld.Fvalue); err == nil {
				keyString := "java/lang/String:" + strconv.FormatUint(hashedUint64, 10)
				hm[keyString] = strObj
			}
		}
	}
	return hs
}

// CreateEmptySet creates an empty java/util/HashSet.
func CreateEmptySet() *object.Object {
	hs := object.MakeEmptyObjectWithClassName(&classNameHashSet)
	hs.FieldTable[fieldNameMap] = object.Field{
		Ftype:  types.Ref,
		Fvalue: make(types.DefHashMap),
	}
	return hs
}
