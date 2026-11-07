/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaTime

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"testing"
)

func isDecimalStyleError(result any) bool {
	if result == nil {
		return false
	}
	_, ok := result.(*ghelpers.GErrBlk)
	return ok
}

func TestDecimalStyle(t *testing.T) {
	globals.InitGlobals("test")
	Load_Time_Format_DecimalStyle()

	// 1. Verify method registration
	expectedSignatures := []string{
		"java/time/format/DecimalStyle.<clinit>()V",
		"java/time/format/DecimalStyle.<init>()V",
		"java/time/format/DecimalStyle.<init>(CCCC)V",
		"java/time/format/DecimalStyle.getAvailableLocales()Ljava/util/Set;",
		"java/time/format/DecimalStyle.ofDefaultLocale()Ljava/time/format/DecimalStyle;",
		"java/time/format/DecimalStyle.of(Ljava/util/Locale;)Ljava/time/format/DecimalStyle;",
		"java/time/format/DecimalStyle.getZeroDigit()C",
		"java/time/format/DecimalStyle.withZeroDigit(C)Ljava/time/format/DecimalStyle;",
		"java/time/format/DecimalStyle.getPositiveSign()C",
		"java/time/format/DecimalStyle.withPositiveSign(C)Ljava/time/format/DecimalStyle;",
		"java/time/format/DecimalStyle.getNegativeSign()C",
		"java/time/format/DecimalStyle.withNegativeSign(C)Ljava/time/format/DecimalStyle;",
		"java/time/format/DecimalStyle.getDecimalSeparator()C",
		"java/time/format/DecimalStyle.withDecimalSeparator(C)Ljava/time/format/DecimalStyle;",
		"java/time/format/DecimalStyle.convertToDigit(C)I",
		"java/time/format/DecimalStyle.convertNumberToI18N(Ljava/lang/String;)Ljava/lang/String;",
		"java/time/format/DecimalStyle.equals(Ljava/lang/Object;)Z",
		"java/time/format/DecimalStyle.hashCode()I",
		"java/time/format/DecimalStyle.toString()Ljava/lang/String;",
	}

	for _, sig := range expectedSignatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Method signature %s not registered", sig)
		}
	}

	// 2. clinit and STANDARD static field
	decimalStyleClinit(nil)
	s, ok := statics.QueryStatic("java/time/format/DecimalStyle", "STANDARD")
	if !ok || s.Value == nil {
		t.Fatalf("Static java/time/format/DecimalStyle.STANDARD not found")
	}
	standardObj := s.Value.(*object.Object)

	// 3. Check STANDARD getters
	zd := decimalStyleGetZeroDigit([]any{standardObj})
	if zd != int64('0') {
		t.Errorf("STANDARD zeroDigit expected '0', got %v", zd)
	}
	ps := decimalStyleGetPositiveSign([]any{standardObj})
	if ps != int64('+') {
		t.Errorf("STANDARD positiveSign expected '+', got %v", ps)
	}
	ns := decimalStyleGetNegativeSign([]any{standardObj})
	if ns != int64('-') {
		t.Errorf("STANDARD negativeSign expected '-', got %v", ns)
	}
	ds := decimalStyleGetDecimalSeparator([]any{standardObj})
	if ds != int64('.') {
		t.Errorf("STANDARD decimalSeparator expected '.', got %v", ds)
	}

	// 4. Test toString
	ts := decimalStyleToString([]any{standardObj})
	tsStrObj, okStr := ts.(*object.Object)
	if !okStr || object.GoStringFromStringObject(tsStrObj) != "DecimalStyle[0+-.]" {
		t.Errorf("STANDARD toString expected DecimalStyle[0+-.] , got %v", object.GoStringFromStringObject(tsStrObj))
	}

	// 5. Test ofDefaultLocale and of(Locale)
	defStyle := decimalStyleOfDefaultLocale(nil).(*object.Object)
	if defStyle != standardObj {
		t.Errorf("Expected ofDefaultLocale to return standardObj, got %v", defStyle)
	}
	localeObj := object.MakeEmptyObjectWithClassName(&classNameDecimalStyle)
	ofLocale := decimalStyleOf([]any{localeObj}).(*object.Object)
	if ofLocale != standardObj {
		t.Errorf("Expected of(locale) to return standardObj, got %v", ofLocale)
	}
	// Null locale
	errLocale := decimalStyleOf([]any{nil})
	if !isDecimalStyleError(errLocale) {
		t.Errorf("Expected NPE for of(null), got %v", errLocale)
	}

	// 6. Test getAvailableLocales
	localesSet := decimalStyleGetAvailableLocales(nil).(*object.Object)
	if object.IsNull(localesSet) {
		t.Fatalf("Expected non-null locales set")
	}

	// 7. Test withZeroDigit
	sameZd := decimalStyleWithZeroDigit([]any{standardObj, int64('0')}).(*object.Object)
	if sameZd != standardObj {
		t.Errorf("withZeroDigit('0') on STANDARD should return same object")
	}
	arabicStyle := decimalStyleWithZeroDigit([]any{standardObj, int64('\u0660')}).(*object.Object)
	if arabicStyle == standardObj {
		t.Errorf("withZeroDigit('\\u0660') should return new object")
	}
	if decimalStyleGetZeroDigit([]any{arabicStyle}) != int64('\u0660') {
		t.Errorf("Expected arabicStyle zeroDigit to be \\u0660")
	}

	// 8. Test withPositiveSign
	samePs := decimalStyleWithPositiveSign([]any{standardObj, int64('+')}).(*object.Object)
	if samePs != standardObj {
		t.Errorf("withPositiveSign('+') on STANDARD should return same object")
	}
	customPs := decimalStyleWithPositiveSign([]any{standardObj, int64('*')}).(*object.Object)
	if decimalStyleGetPositiveSign([]any{customPs}) != int64('*') {
		t.Errorf("Expected positiveSign to be '*'")
	}

	// 9. Test withNegativeSign
	sameNs := decimalStyleWithNegativeSign([]any{standardObj, int64('-')}).(*object.Object)
	if sameNs != standardObj {
		t.Errorf("withNegativeSign('-') on STANDARD should return same object")
	}
	customNs := decimalStyleWithNegativeSign([]any{standardObj, int64('~')}).(*object.Object)
	if decimalStyleGetNegativeSign([]any{customNs}) != int64('~') {
		t.Errorf("Expected negativeSign to be '~'")
	}

	// 10. Test withDecimalSeparator
	sameDs := decimalStyleWithDecimalSeparator([]any{standardObj, int64('.')}).(*object.Object)
	if sameDs != standardObj {
		t.Errorf("withDecimalSeparator('.') on STANDARD should return same object")
	}
	customDs := decimalStyleWithDecimalSeparator([]any{standardObj, int64(',')}).(*object.Object)
	if decimalStyleGetDecimalSeparator([]any{customDs}) != int64(',') {
		t.Errorf("Expected decimalSeparator to be ','")
	}

	// 11. Test convertToDigit
	for i := range 10 {
		digitChar := int64('0' + i)
		val := decimalStyleConvertToDigit([]any{standardObj, digitChar})
		if val != int64(i) {
			t.Errorf("convertToDigit('%c') expected %d, got %v", rune(digitChar), i, val)
		}
	}
	invalidDigit := decimalStyleConvertToDigit([]any{standardObj, int64('a')})
	if invalidDigit != int64(-1) {
		t.Errorf("convertToDigit('a') expected -1, got %v", invalidDigit)
	}

	// With arabic digits
	for i := range 10 {
		digitChar := int64('\u0660' + rune(i))
		val := decimalStyleConvertToDigit([]any{arabicStyle, digitChar})
		if val != int64(i) {
			t.Errorf("arabicStyle convertToDigit expected %d, got %v", i, val)
		}
	}

	// 12. Test convertNumberToI18N
	standardText := object.StringObjectFromGoString("12345")
	convertedStandard := decimalStyleConvertNumberToI18N([]any{standardObj, standardText}).(*object.Object)
	if object.GoStringFromStringObject(convertedStandard) != "12345" {
		t.Errorf("convertNumberToI18N with standard style expected 12345, got %s", object.GoStringFromStringObject(convertedStandard))
	}

	convertedArabic := decimalStyleConvertNumberToI18N([]any{arabicStyle, standardText}).(*object.Object)
	expectedArabic := string([]rune{'\u0661', '\u0662', '\u0663', '\u0664', '\u0665'})
	if object.GoStringFromStringObject(convertedArabic) != expectedArabic {
		t.Errorf("convertNumberToI18N with arabic style expected %s, got %s", expectedArabic, object.GoStringFromStringObject(convertedArabic))
	}

	// 13. Test equals and hashCode
	eqSame := decimalStyleEquals([]any{standardObj, standardObj})
	if eqSame != types.JavaBoolTrue {
		t.Errorf("equals(same) expected true")
	}
	cloneStandard := createDecimalStyle('0', '+', '-', '.')
	eqClone := decimalStyleEquals([]any{standardObj, cloneStandard})
	if eqClone != types.JavaBoolTrue {
		t.Errorf("equals(clone) expected true")
	}
	eqDiff := decimalStyleEquals([]any{standardObj, arabicStyle})
	if eqDiff != types.JavaBoolFalse {
		t.Errorf("equals(diff) expected false")
	}

	hc1 := decimalStyleHashCode([]any{standardObj}).(int64)
	hc2 := decimalStyleHashCode([]any{cloneStandard}).(int64)
	if hc1 != hc2 {
		t.Errorf("hashCode of equal objects should match: %d vs %d", hc1, hc2)
	}

	// 14. Test constructors
	customObj := object.MakeEmptyObjectWithClassName(&classNameDecimalStyle)
	decimalStyleInitDefault([]any{customObj})
	if decimalStyleGetZeroDigit([]any{customObj}) != int64('0') {
		t.Errorf("InitDefault failed")
	}

	customObj4 := object.MakeEmptyObjectWithClassName(&classNameDecimalStyle)
	decimalStyleInit4([]any{customObj4, int64('0'), int64('+'), int64('-'), int64(',')})
	if decimalStyleGetDecimalSeparator([]any{customObj4}) != int64(',') {
		t.Errorf("Init4 failed to set decimal separator")
	}

	// 15. Error cases
	if !isDecimalStyleError(decimalStyleGetZeroDigit([]any{nil})) {
		t.Errorf("Expected error for getZeroDigit(null)")
	}
	if !isDecimalStyleError(decimalStyleWithZeroDigit([]any{nil, int64('0')})) {
		t.Errorf("Expected error for withZeroDigit(null, '0')")
	}
	if !isDecimalStyleError(decimalStyleConvertNumberToI18N([]any{standardObj, nil})) {
		t.Errorf("Expected error for convertNumberToI18N(standard, null)")
	}
}
