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
	"sync"
)

var classNameDecimalStyle = "java/time/format/DecimalStyle"
var classNameDecimalStyleHashSet = "java/util/HashSet"

var decimalStyleMutex sync.Mutex
var decimalStyleInitialized bool = false
var decimalStyleStandard *object.Object

func ensureDecimalStyleInited() {
	decimalStyleMutex.Lock()
	defer decimalStyleMutex.Unlock()
	if decimalStyleInitialized {
		return
	}

	decimalStyleStandard = createDecimalStyle('0', '+', '-', '.')
	_ = statics.AddStatic(classNameDecimalStyle+".STANDARD", statics.Static{
		Type:  "Ljava/time/format/DecimalStyle;",
		Value: decimalStyleStandard,
	})

	decimalStyleInitialized = true
}

func createDecimalStyle(zeroDigit, positiveSign, negativeSign, decimalSeparator rune) *object.Object {
	styleObj := object.MakeEmptyObjectWithClassName(&classNameDecimalStyle)
	styleObj.FieldTable["zeroDigit"] = object.Field{
		Ftype:  types.Char,
		Fvalue: int64(zeroDigit),
	}
	styleObj.FieldTable["positiveSign"] = object.Field{
		Ftype:  types.Char,
		Fvalue: int64(positiveSign),
	}
	styleObj.FieldTable["negativeSign"] = object.Field{
		Ftype:  types.Char,
		Fvalue: int64(negativeSign),
	}
	styleObj.FieldTable["decimalSeparator"] = object.Field{
		Ftype:  types.Char,
		Fvalue: int64(decimalSeparator),
	}
	return styleObj
}

func getDecimalStyleFields(self *object.Object) (rune, rune, rune, rune, bool) {
	if self == nil || self.FieldTable == nil {
		return '0', '+', '-', '.', false
	}
	zd := rune('0')
	ps := rune('+')
	ns := rune('-')
	ds := rune('.')
	if f, ok := self.FieldTable["zeroDigit"]; ok {
		if v, ok := ghelpers.GetInt64Val(f.Fvalue); ok {
			zd = rune(v)
		}
	}
	if f, ok := self.FieldTable["positiveSign"]; ok {
		if v, ok := ghelpers.GetInt64Val(f.Fvalue); ok {
			ps = rune(v)
		}
	}
	if f, ok := self.FieldTable["negativeSign"]; ok {
		if v, ok := ghelpers.GetInt64Val(f.Fvalue); ok {
			ns = rune(v)
		}
	}
	if f, ok := self.FieldTable["decimalSeparator"]; ok {
		if v, ok := ghelpers.GetInt64Val(f.Fvalue); ok {
			ds = rune(v)
		}
	}
	return zd, ps, ns, ds, true
}

// Load_Time_Format_DecimalStyle registers all java/time/format/DecimalStyle methods.
func Load_Time_Format_DecimalStyle() {
	ghelpers.MethodSignatures["java/time/format/DecimalStyle.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  decimalStyleClinit,
		}

	ghelpers.MethodSignatures["java/time/format/DecimalStyle.<init>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  decimalStyleInitDefault,
		}

	ghelpers.MethodSignatures["java/time/format/DecimalStyle.<init>(CCCC)V"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  decimalStyleInit4,
		}

	ghelpers.MethodSignatures["java/time/format/DecimalStyle.getAvailableLocales()Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  decimalStyleGetAvailableLocales,
		}

	ghelpers.MethodSignatures["java/time/format/DecimalStyle.ofDefaultLocale()Ljava/time/format/DecimalStyle;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  decimalStyleOfDefaultLocale,
		}

	ghelpers.MethodSignatures["java/time/format/DecimalStyle.of(Ljava/util/Locale;)Ljava/time/format/DecimalStyle;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  decimalStyleOf,
		}

	ghelpers.MethodSignatures["java/time/format/DecimalStyle.getZeroDigit()C"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  decimalStyleGetZeroDigit,
		}

	ghelpers.MethodSignatures["java/time/format/DecimalStyle.withZeroDigit(C)Ljava/time/format/DecimalStyle;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  decimalStyleWithZeroDigit,
		}

	ghelpers.MethodSignatures["java/time/format/DecimalStyle.getPositiveSign()C"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  decimalStyleGetPositiveSign,
		}

	ghelpers.MethodSignatures["java/time/format/DecimalStyle.withPositiveSign(C)Ljava/time/format/DecimalStyle;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  decimalStyleWithPositiveSign,
		}

	ghelpers.MethodSignatures["java/time/format/DecimalStyle.getNegativeSign()C"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  decimalStyleGetNegativeSign,
		}

	ghelpers.MethodSignatures["java/time/format/DecimalStyle.withNegativeSign(C)Ljava/time/format/DecimalStyle;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  decimalStyleWithNegativeSign,
		}

	ghelpers.MethodSignatures["java/time/format/DecimalStyle.getDecimalSeparator()C"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  decimalStyleGetDecimalSeparator,
		}

	ghelpers.MethodSignatures["java/time/format/DecimalStyle.withDecimalSeparator(C)Ljava/time/format/DecimalStyle;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  decimalStyleWithDecimalSeparator,
		}

	ghelpers.MethodSignatures["java/time/format/DecimalStyle.convertToDigit(C)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  decimalStyleConvertToDigit,
		}

	ghelpers.MethodSignatures["java/time/format/DecimalStyle.convertNumberToI18N(Ljava/lang/String;)Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  decimalStyleConvertNumberToI18N,
		}

	ghelpers.MethodSignatures["java/time/format/DecimalStyle.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  decimalStyleEquals,
		}

	ghelpers.MethodSignatures["java/time/format/DecimalStyle.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  decimalStyleHashCode,
		}

	ghelpers.MethodSignatures["java/time/format/DecimalStyle.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  decimalStyleToString,
		}
}

func decimalStyleClinit([]any) any {
	ensureDecimalStyleInited()
	return nil
}

func decimalStyleInitDefault(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DecimalStyle.<init>: self is null")
	}
	self := params[0].(*object.Object)
	self.FieldTable["zeroDigit"] = object.Field{Ftype: types.Char, Fvalue: int64('0')}
	self.FieldTable["positiveSign"] = object.Field{Ftype: types.Char, Fvalue: int64('+')}
	self.FieldTable["negativeSign"] = object.Field{Ftype: types.Char, Fvalue: int64('-')}
	self.FieldTable["decimalSeparator"] = object.Field{Ftype: types.Char, Fvalue: int64('.')}
	return nil
}

func decimalStyleInit4(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DecimalStyle.<init>: self is null")
	}
	self := params[0].(*object.Object)
	zd := rune('0')
	ps := rune('+')
	ns := rune('-')
	ds := rune('.')
	if len(params) > 1 {
		if v, ok := ghelpers.GetInt64Val(params[1]); ok {
			zd = rune(v)
		}
	}
	if len(params) > 2 {
		if v, ok := ghelpers.GetInt64Val(params[2]); ok {
			ps = rune(v)
		}
	}
	if len(params) > 3 {
		if v, ok := ghelpers.GetInt64Val(params[3]); ok {
			ns = rune(v)
		}
	}
	if len(params) > 4 {
		if v, ok := ghelpers.GetInt64Val(params[4]); ok {
			ds = rune(v)
		}
	}
	self.FieldTable["zeroDigit"] = object.Field{Ftype: types.Char, Fvalue: int64(zd)}
	self.FieldTable["positiveSign"] = object.Field{Ftype: types.Char, Fvalue: int64(ps)}
	self.FieldTable["negativeSign"] = object.Field{Ftype: types.Char, Fvalue: int64(ns)}
	self.FieldTable["decimalSeparator"] = object.Field{Ftype: types.Char, Fvalue: int64(ds)}
	return nil
}

func decimalStyleGetAvailableLocales([]any) any {
	ensureDecimalStyleInited()
	setObj := object.MakeEmptyObjectWithClassName(&classNameDecimalStyleHashSet)
	setObj.FieldTable["map"] = object.Field{
		Ftype:  types.Ref,
		Fvalue: make(types.DefHashMap),
	}
	return setObj
}

func decimalStyleOfDefaultLocale([]any) any {
	ensureDecimalStyleInited()
	return decimalStyleStandard
}

func decimalStyleOf(params []any) any {
	ensureDecimalStyleInited()
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DecimalStyle.of: locale is null")
	}
	return decimalStyleStandard
}

func decimalStyleGetZeroDigit(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DecimalStyle.getZeroDigit: self is null")
	}
	self := params[0].(*object.Object)
	zd, _, _, _, ok := getDecimalStyleFields(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DecimalStyle.getZeroDigit: invalid instance")
	}
	return int64(zd)
}

func decimalStyleWithZeroDigit(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DecimalStyle.withZeroDigit: self is null")
	}
	self := params[0].(*object.Object)
	zd, ps, ns, ds, ok := getDecimalStyleFields(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DecimalStyle.withZeroDigit: invalid instance")
	}
	newZdVal, okVal := ghelpers.GetInt64Val(params[1])
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DecimalStyle.withZeroDigit: invalid char")
	}
	newZd := rune(newZdVal)
	if newZd == zd {
		return self
	}
	return createDecimalStyle(newZd, ps, ns, ds)
}

func decimalStyleGetPositiveSign(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DecimalStyle.getPositiveSign: self is null")
	}
	self := params[0].(*object.Object)
	_, ps, _, _, ok := getDecimalStyleFields(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DecimalStyle.getPositiveSign: invalid instance")
	}
	return int64(ps)
}

func decimalStyleWithPositiveSign(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DecimalStyle.withPositiveSign: self is null")
	}
	self := params[0].(*object.Object)
	zd, ps, ns, ds, ok := getDecimalStyleFields(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DecimalStyle.withPositiveSign: invalid instance")
	}
	newPsVal, okVal := ghelpers.GetInt64Val(params[1])
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DecimalStyle.withPositiveSign: invalid char")
	}
	newPs := rune(newPsVal)
	if newPs == ps {
		return self
	}
	return createDecimalStyle(zd, newPs, ns, ds)
}

func decimalStyleGetNegativeSign(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DecimalStyle.getNegativeSign: self is null")
	}
	self := params[0].(*object.Object)
	_, _, ns, _, ok := getDecimalStyleFields(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DecimalStyle.getNegativeSign: invalid instance")
	}
	return int64(ns)
}

func decimalStyleWithNegativeSign(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DecimalStyle.withNegativeSign: self is null")
	}
	self := params[0].(*object.Object)
	zd, ps, ns, ds, ok := getDecimalStyleFields(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DecimalStyle.withNegativeSign: invalid instance")
	}
	newNsVal, okVal := ghelpers.GetInt64Val(params[1])
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DecimalStyle.withNegativeSign: invalid char")
	}
	newNs := rune(newNsVal)
	if newNs == ns {
		return self
	}
	return createDecimalStyle(zd, ps, newNs, ds)
}

func decimalStyleGetDecimalSeparator(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DecimalStyle.getDecimalSeparator: self is null")
	}
	self := params[0].(*object.Object)
	_, _, _, ds, ok := getDecimalStyleFields(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DecimalStyle.getDecimalSeparator: invalid instance")
	}
	return int64(ds)
}

func decimalStyleWithDecimalSeparator(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DecimalStyle.withDecimalSeparator: self is null")
	}
	self := params[0].(*object.Object)
	zd, ps, ns, ds, ok := getDecimalStyleFields(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DecimalStyle.withDecimalSeparator: invalid instance")
	}
	newDsVal, okVal := ghelpers.GetInt64Val(params[1])
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DecimalStyle.withDecimalSeparator: invalid char")
	}
	newDs := rune(newDsVal)
	if newDs == ds {
		return self
	}
	return createDecimalStyle(zd, ps, ns, newDs)
}

func decimalStyleConvertToDigit(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DecimalStyle.convertToDigit: self is null")
	}
	self := params[0].(*object.Object)
	zd, _, _, _, ok := getDecimalStyleFields(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DecimalStyle.convertToDigit: invalid instance")
	}
	chVal, okVal := ghelpers.GetInt64Val(params[1])
	if !okVal {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DecimalStyle.convertToDigit: invalid char")
	}
	ch := rune(chVal)
	val := ch - zd
	if val >= 0 && val <= 9 {
		return int64(val)
	}
	return int64(-1)
}

func decimalStyleConvertNumberToI18N(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DecimalStyle.convertNumberToI18N: null parameter")
	}
	self := params[0].(*object.Object)
	textObj := params[1].(*object.Object)
	if !object.IsStringObject(textObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DecimalStyle.convertNumberToI18N: text is not a String")
	}

	zd, _, _, _, ok := getDecimalStyleFields(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DecimalStyle.convertNumberToI18N: invalid instance")
	}

	if zd == '0' {
		return textObj
	}

	text := object.GoStringFromStringObject(textObj)
	diff := zd - '0'
	runes := []rune(text)
	for i := range runes {
		runes[i] = runes[i] + diff
	}
	return object.StringObjectFromGoString(string(runes))
}

func decimalStyleEquals(params []any) any {
	if len(params) < 2 || params[0] == nil || params[1] == nil ||
		object.IsNull(params[0].(*object.Object)) || object.IsNull(params[1].(*object.Object)) {
		return types.JavaBoolFalse
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	if self == other {
		return types.JavaBoolTrue
	}

	zd1, ps1, ns1, ds1, ok1 := getDecimalStyleFields(self)
	zd2, ps2, ns2, ds2, ok2 := getDecimalStyleFields(other)
	if ok1 && ok2 && zd1 == zd2 && ps1 == ps2 && ns1 == ns2 && ds1 == ds2 {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func decimalStyleHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DecimalStyle.hashCode: self is null")
	}
	self := params[0].(*object.Object)
	zd, ps, ns, ds, ok := getDecimalStyleFields(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DecimalStyle.hashCode: invalid instance")
	}
	return int64(zd + ps + ns + ds)
}

func decimalStyleToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "DecimalStyle.toString: self is null")
	}
	self := params[0].(*object.Object)
	zd, ps, ns, ds, ok := getDecimalStyleFields(self)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "DecimalStyle.toString: invalid instance")
	}
	str := fmt.Sprintf("DecimalStyle[%c%c%c%c]", zd, ps, ns, ds)
	return object.StringObjectFromGoString(str)
}
