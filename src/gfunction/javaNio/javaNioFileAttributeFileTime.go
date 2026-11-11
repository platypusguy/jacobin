/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaNio

import (
	"fmt"
	"math"
	"math/big"
	"time"

	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/gfunction/javaTime"
	"jacobin/src/object"
	"jacobin/src/types"
)

var fileTimeClassName = "java/nio/file/attribute/FileTime"

const (
	instantMinSecond int64 = -31557014167219200
	instantMaxSecond int64 = 31556889864403199
)

var timeUnitOrdinalNames = []string{
	"NANOSECONDS",
	"MICROSECONDS",
	"MILLISECONDS",
	"SECONDS",
	"MINUTES",
	"HOURS",
	"DAYS",
}

var timeUnitScaleInNanos = map[string]int64{
	"NANOSECONDS":  1,
	"MICROSECONDS": 1000,
	"MILLISECONDS": 1000000,
	"SECONDS":      1000000000,
	"MINUTES":      60000000000,
	"HOURS":        3600000000000,
	"DAYS":         86400000000000,
}

// Load_Nio_File_Attribute_FileTime registers all java.nio.file.attribute.FileTime methods.
func Load_Nio_File_Attribute_FileTime() {
	// <clinit>()V
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.<clinit>()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.ClinitGeneric}

	// from(JLjava/util/concurrent/TimeUnit;)Ljava/nio/file/attribute/FileTime;
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.from(JLjava/util/concurrent/TimeUnit;)Ljava/nio/file/attribute/FileTime;"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: fileTimeFromUnit}

	// from(Ljava/time/Instant;)Ljava/nio/file/attribute/FileTime;
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.from(Ljava/time/Instant;)Ljava/nio/file/attribute/FileTime;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: fileTimeFromInstant}

	// fromMillis(J)Ljava/nio/file/attribute/FileTime;
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.fromMillis(J)Ljava/nio/file/attribute/FileTime;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: fileTimeFromMillis}

	// to(Ljava/util/concurrent/TimeUnit;)J
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.to(Ljava/util/concurrent/TimeUnit;)J"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: fileTimeToUnit}

	// toMillis()J
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.toMillis()J"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: fileTimeToMillis}

	// toInstant()Ljava/time/Instant;
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.toInstant()Ljava/time/Instant;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: fileTimeToInstant}

	// compareTo(Ljava/nio/file/attribute/FileTime;)I
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.compareTo(Ljava/nio/file/attribute/FileTime;)I"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: fileTimeCompareTo}

	// compareTo(Ljava/lang/Object;)I
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: fileTimeCompareToObject}

	// equals(Ljava/lang/Object;)Z
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: fileTimeEquals}

	// hashCode()I
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.hashCode()I"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: fileTimeHashCode}

	// toString()Ljava/lang/String;
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: fileTimeToString}

	// Standard Object lifecycle & concurrency methods
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.clone()Ljava/lang/Object;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.finalize()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.getClass()Ljava/lang/Class;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.notify()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.notifyAll()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.wait()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.wait(J)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: ghelpers.TrapFunction}
	ghelpers.MethodSignatures["java/nio/file/attribute/FileTime.wait(JI)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: ghelpers.TrapFunction}
}

func makeFileTimeFromUnit(val int64, unit string) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&fileTimeClassName)
	obj.FieldTable["value"] = object.Field{Ftype: types.Long, Fvalue: val}
	obj.FieldTable["unit"] = object.Field{Ftype: types.StringClassRef, Fvalue: object.StringObjectFromGoString(unit)}
	obj.FieldTable["unitName"] = object.Field{Ftype: "string", Fvalue: unit}
	obj.FieldTable["fromInstant"] = object.Field{Ftype: types.Bool, Fvalue: types.JavaBoolFalse}
	return obj
}

func makeFileTimeFromInstant(inst *object.Object) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&fileTimeClassName)
	sec, nano := extractInstantSecondsAndNanos(inst)
	obj.FieldTable["instant"] = object.Field{Ftype: "Ljava/time/Instant;", Fvalue: inst}
	obj.FieldTable["seconds"] = object.Field{Ftype: types.Long, Fvalue: sec}
	obj.FieldTable["nanos"] = object.Field{Ftype: types.Long, Fvalue: nano}
	obj.FieldTable["fromInstant"] = object.Field{Ftype: types.Bool, Fvalue: types.JavaBoolTrue}
	return obj
}

func newFileTime(t time.Time) *object.Object {
	return makeFileTimeFromUnit(t.UnixNano(), "NANOSECONDS")
}

func parseTimeUnit(arg any) (string, *ghelpers.GErrBlk) {
	if arg == nil {
		return "", ghelpers.GetGErrBlk(excNames.NullPointerException, "FileTime: unit is null")
	}
	if obj, ok := arg.(*object.Object); ok {
		if object.IsNull(obj) {
			return "", ghelpers.GetGErrBlk(excNames.NullPointerException, "FileTime: unit is null")
		}
		if object.IsStringObject(obj) {
			unitStr := object.GoStringFromStringObject(obj)
			if _, exists := timeUnitScaleInNanos[unitStr]; exists {
				return unitStr, nil
			}
			return "", ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "FileTime: invalid TimeUnit: "+unitStr)
		}
		if fld, exists := obj.FieldTable["name"]; exists {
			var name string
			if sObj, ok := fld.Fvalue.(*object.Object); ok {
				name = object.GoStringFromStringObject(sObj)
			} else if s, ok := fld.Fvalue.(string); ok {
				name = s
			}
			if _, exists := timeUnitScaleInNanos[name]; exists {
				return name, nil
			}
		}
		if fld, exists := obj.FieldTable["ordinal"]; exists {
			if ord, ok := fld.Fvalue.(int64); ok && ord >= 0 && int(ord) < len(timeUnitOrdinalNames) {
				return timeUnitOrdinalNames[ord], nil
			}
		}
	} else if str, ok := arg.(string); ok {
		if _, exists := timeUnitScaleInNanos[str]; exists {
			return str, nil
		}
		return "", ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "FileTime: invalid TimeUnit: "+str)
	}
	return "", ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "FileTime: invalid TimeUnit argument")
}

func isFromInstant(obj *object.Object) bool {
	if obj == nil || obj.FieldTable == nil {
		return false
	}
	if fld, ok := obj.FieldTable["fromInstant"]; ok {
		if val, ok2 := ghelpers.GetInt64Val(fld.Fvalue); ok2 {
			return val == types.JavaBoolTrue
		}
	}
	_, hasInst := obj.FieldTable["instant"]
	_, hasSec := obj.FieldTable["seconds"]
	return hasInst || hasSec
}

func getFileTimeUnit(obj *object.Object) string {
	if obj == nil || obj.FieldTable == nil {
		return "MILLISECONDS"
	}
	if fld, ok := obj.FieldTable["unitName"]; ok {
		if s, ok2 := fld.Fvalue.(string); ok2 && s != "" {
			return s
		}
	}
	if fld, ok := obj.FieldTable["unit"]; ok {
		if sObj, ok2 := fld.Fvalue.(*object.Object); ok2 {
			s := object.GoStringFromStringObject(sObj)
			if s != "" {
				return s
			}
		}
	}
	return "MILLISECONDS"
}

func getFileTimeValue(obj *object.Object) int64 {
	if obj == nil || obj.FieldTable == nil {
		return 0
	}
	if fld, ok := obj.FieldTable["value"]; ok {
		if v, ok2 := ghelpers.GetInt64Val(fld.Fvalue); ok2 {
			return v
		}
	}
	return 0
}

func extractInstantSecondsAndNanos(obj *object.Object) (int64, int64) {
	if obj == nil || obj.FieldTable == nil {
		return 0, 0
	}
	var seconds, nanos int64
	if sFld, ok := obj.FieldTable["seconds"]; ok {
		if s, ok2 := ghelpers.GetInt64Val(sFld.Fvalue); ok2 {
			seconds = s
		}
	} else if sFld, ok := obj.FieldTable["epochSecond"]; ok {
		if s, ok2 := ghelpers.GetInt64Val(sFld.Fvalue); ok2 {
			seconds = s
		}
	}
	if nFld, ok := obj.FieldTable["nanos"]; ok {
		if n, ok2 := ghelpers.GetInt64Val(nFld.Fvalue); ok2 {
			nanos = n
		}
	} else if nFld, ok := obj.FieldTable["nano"]; ok {
		if n, ok2 := ghelpers.GetInt64Val(nFld.Fvalue); ok2 {
			nanos = n
		}
	}
	return seconds, nanos
}

func fileTimeToSecondsAndNanos(obj *object.Object) (int64, int64, *big.Int) {
	if isFromInstant(obj) {
		var sec, nano int64
		if instFld, ok := obj.FieldTable["instant"]; ok {
			if inst, ok2 := instFld.Fvalue.(*object.Object); ok2 && !object.IsNull(inst) {
				sec, nano = extractInstantSecondsAndNanos(inst)
			} else {
				sec, nano = extractInstantSecondsAndNanos(obj)
			}
		} else {
			sec, nano = extractInstantSecondsAndNanos(obj)
		}
		return sec, nano, big.NewInt(sec)
	}

	val := getFileTimeValue(obj)
	unit := getFileTimeUnit(obj)

	switch unit {
	case "NANOSECONDS":
		sec := floorDiv(val, 1000000000)
		nano := floorMod(val, 1000000000)
		return sec, nano, big.NewInt(sec)
	case "MICROSECONDS":
		sec := floorDiv(val, 1000000)
		nano := floorMod(val, 1000000) * 1000
		return sec, nano, big.NewInt(sec)
	case "MILLISECONDS":
		sec := floorDiv(val, 1000)
		nano := floorMod(val, 1000) * 1000000
		return sec, nano, big.NewInt(sec)
	case "SECONDS":
		return val, 0, big.NewInt(val)
	case "MINUTES":
		bigSec := new(big.Int).Mul(big.NewInt(val), big.NewInt(60))
		return saturateBigInt(bigSec), 0, bigSec
	case "HOURS":
		bigSec := new(big.Int).Mul(big.NewInt(val), big.NewInt(3600))
		return saturateBigInt(bigSec), 0, bigSec
	case "DAYS":
		bigSec := new(big.Int).Mul(big.NewInt(val), big.NewInt(86400))
		return saturateBigInt(bigSec), 0, bigSec
	default:
		sec := floorDiv(val, 1000)
		nano := floorMod(val, 1000) * 1000000
		return sec, nano, big.NewInt(sec)
	}
}

func saturateBigInt(b *big.Int) int64 {
	minBig := big.NewInt(math.MinInt64)
	maxBig := big.NewInt(math.MaxInt64)
	if b.Cmp(minBig) < 0 {
		return math.MinInt64
	}
	if b.Cmp(maxBig) > 0 {
		return math.MaxInt64
	}
	return b.Int64()
}

func convertTimeUnit(val int64, srcUnit, dstUnit string) int64 {
	if srcUnit == dstUnit {
		return val
	}
	srcScale := timeUnitScaleInNanos[srcUnit]
	dstScale := timeUnitScaleInNanos[dstUnit]
	if srcScale > dstScale {
		// Down-conversion: multiply by scale factor with saturation
		scale := srcScale / dstScale
		if val > math.MaxInt64/scale {
			return math.MaxInt64
		}
		if val < math.MinInt64/scale {
			return math.MinInt64
		}
		return val * scale
	}
	// Up-conversion: divide by scale factor
	scale := dstScale / srcScale
	return val / scale
}

func convertInstantToUnit(sec int64, nano int64, dstUnit string) int64 {
	switch dstUnit {
	case "NANOSECONDS":
		if sec > math.MaxInt64/1000000000 {
			return math.MaxInt64
		}
		if sec < math.MinInt64/1000000000 {
			return math.MinInt64
		}
		secNanos := sec * 1000000000
		if secNanos > 0 && math.MaxInt64-secNanos < nano {
			return math.MaxInt64
		}
		return secNanos + nano
	case "MICROSECONDS":
		if sec > math.MaxInt64/1000000 {
			return math.MaxInt64
		}
		if sec < math.MinInt64/1000000 {
			return math.MinInt64
		}
		secMicros := sec * 1000000
		if secMicros > 0 && math.MaxInt64-secMicros < (nano/1000) {
			return math.MaxInt64
		}
		return secMicros + nano/1000
	case "MILLISECONDS":
		if sec > math.MaxInt64/1000 {
			return math.MaxInt64
		}
		if sec < math.MinInt64/1000 {
			return math.MinInt64
		}
		secMillis := sec * 1000
		if secMillis > 0 && math.MaxInt64-secMillis < (nano/1000000) {
			return math.MaxInt64
		}
		return secMillis + nano/1000000
	case "SECONDS":
		return sec
	case "MINUTES":
		return sec / 60
	case "HOURS":
		return sec / 3600
	case "DAYS":
		return sec / 86400
	default:
		return 0
	}
}

// fileTimeFromUnit: static FileTime.from(long value, TimeUnit unit)
func fileTimeFromUnit(params []any) any {
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "FileTime.from: missing arguments")
	}
	val, ok := ghelpers.GetInt64Val(params[0])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "FileTime.from: invalid value argument")
	}
	unit, errBlk := parseTimeUnit(params[1])
	if errBlk != nil {
		return errBlk
	}
	return makeFileTimeFromUnit(val, unit)
}

// fileTimeFromInstant: static FileTime.from(Instant instant)
func fileTimeFromInstant(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "FileTime.from: instant is null")
	}
	inst, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "FileTime.from: argument is not an Instant")
	}
	return makeFileTimeFromInstant(inst)
}

// fileTimeFromMillis: static FileTime.fromMillis(long value)
func fileTimeFromMillis(params []any) any {
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "FileTime.fromMillis: missing argument")
	}
	val, ok := ghelpers.GetInt64Val(params[0])
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "FileTime.fromMillis: invalid value argument")
	}
	return makeFileTimeFromUnit(val, "MILLISECONDS")
}

// fileTimeToUnit: FileTime.to(TimeUnit unit)
func fileTimeToUnit(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "FileTime.to: self is null")
	}
	self := params[0].(*object.Object)
	dstUnit, errBlk := parseTimeUnit(params[1])
	if errBlk != nil {
		return errBlk
	}

	if isFromInstant(self) {
		sec, nano, _ := fileTimeToSecondsAndNanos(self)
		return convertInstantToUnit(sec, nano, dstUnit)
	}

	srcVal := getFileTimeValue(self)
	srcUnit := getFileTimeUnit(self)
	return convertTimeUnit(srcVal, srcUnit, dstUnit)
}

// fileTimeToMillis: FileTime.toMillis()
func fileTimeToMillis(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "FileTime.toMillis: self is null")
	}
	self := params[0].(*object.Object)
	if isFromInstant(self) {
		sec, nano, _ := fileTimeToSecondsAndNanos(self)
		return convertInstantToUnit(sec, nano, "MILLISECONDS")
	}
	srcVal := getFileTimeValue(self)
	srcUnit := getFileTimeUnit(self)
	return convertTimeUnit(srcVal, srcUnit, "MILLISECONDS")
}

// fileTimeToInstant: FileTime.toInstant()
func fileTimeToInstant(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "FileTime.toInstant: self is null")
	}
	self := params[0].(*object.Object)
	return fileTimeToInstantInternal(self)
}

func fileTimeToInstantInternal(self *object.Object) *object.Object {
	if isFromInstant(self) {
		if instFld, ok := self.FieldTable["instant"]; ok {
			if inst, ok2 := instFld.Fvalue.(*object.Object); ok2 && !object.IsNull(inst) {
				return inst
			}
		}
		sec, nano, _ := fileTimeToSecondsAndNanos(self)
		return javaTime.CreateInstant(sec, nano)
	}

	_, nano, bigSec := fileTimeToSecondsAndNanos(self)
	minBig := big.NewInt(instantMinSecond)
	maxBig := big.NewInt(instantMaxSecond)

	if bigSec.Cmp(minBig) < 0 {
		return javaTime.CreateInstant(instantMinSecond, 0)
	}
	if bigSec.Cmp(maxBig) > 0 {
		return javaTime.CreateInstant(instantMaxSecond, 999999999)
	}
	return javaTime.CreateInstant(bigSec.Int64(), nano)
}

// fileTimeCompareTo: FileTime.compareTo(FileTime other)
func fileTimeCompareTo(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "FileTime.compareTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "FileTime.compareTo: other is null")
	}
	self := params[0].(*object.Object)
	other := params[1].(*object.Object)

	return compareFileTimes(self, other)
}

// fileTimeCompareToObject: FileTime.compareTo(Object other)
func fileTimeCompareToObject(params []any) any {
	if len(params) < 2 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "FileTime.compareTo: self is null")
	}
	if params[1] == nil || object.IsNull(params[1].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "FileTime.compareTo: other is null")
	}
	other := params[1].(*object.Object)
	kName := object.GoStringFromStringPoolIndex(other.KlassName)
	if kName != fileTimeClassName {
		return ghelpers.GetGErrBlk(excNames.ClassCastException, "FileTime.compareTo: cannot compare to "+kName)
	}
	return fileTimeCompareTo(params)
}

func compareFileTimes(ft1, ft2 *object.Object) int64 {
	u1 := getFileTimeUnit(ft1)
	u2 := getFileTimeUnit(ft2)
	fromInst1 := isFromInstant(ft1)
	fromInst2 := isFromInstant(ft2)

	if !fromInst1 && !fromInst2 && u1 == u2 {
		v1 := getFileTimeValue(ft1)
		v2 := getFileTimeValue(ft2)
		if v1 < v2 {
			return -1
		}
		if v1 > v2 {
			return 1
		}
		return 0
	}

	_, n1, bigSec1 := fileTimeToSecondsAndNanos(ft1)
	_, n2, bigSec2 := fileTimeToSecondsAndNanos(ft2)

	secCmp := bigSec1.Cmp(bigSec2)
	if secCmp < 0 {
		return -1
	}
	if secCmp > 0 {
		return 1
	}
	if n1 < n2 {
		return -1
	}
	if n1 > n2 {
		return 1
	}
	return 0
}

// fileTimeEquals: FileTime.equals(Object other)
func fileTimeEquals(params []any) any {
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

	kName := object.GoStringFromStringPoolIndex(other.KlassName)
	if kName != fileTimeClassName {
		return types.JavaBoolFalse
	}

	cmp := compareFileTimes(self, other)
	return types.ConvertGoBoolToJavaBool(cmp == 0)
}

// fileTimeHashCode: FileTime.hashCode()
func fileTimeHashCode(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "FileTime.hashCode: self is null")
	}
	self := params[0].(*object.Object)
	inst := fileTimeToInstantInternal(self)
	sec, nano := extractInstantSecondsAndNanos(inst)
	uSec := uint64(sec)
	h := int32(uSec^(uSec>>32)) + int32(51*nano)
	return int64(h)
}

// fileTimeToString: FileTime.toString()
func fileTimeToString(params []any) any {
	if len(params) < 1 || params[0] == nil || object.IsNull(params[0].(*object.Object)) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "FileTime.toString: self is null")
	}
	self := params[0].(*object.Object)
	res := fileTimeToStringInternal(self)
	return object.StringObjectFromGoString(res)
}

func fileTimeToStringInternal(self *object.Object) string {
	fromInst := isFromInstant(self)
	if fromInst {
		sec, nano, _ := fileTimeToSecondsAndNanos(self)
		return formatIso8601(sec, nano, true, "")
	}

	unit := getFileTimeUnit(self)
	val := getFileTimeValue(self)

	var sec int64
	var fracStr string
	switch unit {
	case "NANOSECONDS":
		sec = floorDiv(val, 1000000000)
		nanos := floorMod(val, 1000000000)
		fracStr = fmt.Sprintf(".%09d", nanos)
	case "MICROSECONDS":
		sec = floorDiv(val, 1000000)
		micros := floorMod(val, 1000000)
		fracStr = fmt.Sprintf(".%06d", micros)
	case "MILLISECONDS":
		sec = floorDiv(val, 1000)
		millis := floorMod(val, 1000)
		fracStr = fmt.Sprintf(".%03d", millis)
	case "SECONDS":
		sec = val
		fracStr = ""
	case "MINUTES":
		sec = val * 60
		fracStr = ""
	case "HOURS":
		sec = val * 3600
		fracStr = ""
	case "DAYS":
		sec = val * 86400
		fracStr = ""
	default:
		sec = floorDiv(val, 1000)
		millis := floorMod(val, 1000)
		fracStr = fmt.Sprintf(".%03d", millis)
	}

	return formatIso8601(sec, 0, false, fracStr)
}

func formatIso8601(sec int64, nano int64, fromInst bool, explicitFrac string) string {
	epochDay := floorDiv(sec, 86400)
	secOfDay := floorMod(sec, 86400)
	y, m, d := epochDayToYMD(epochDay)
	h := secOfDay / 3600
	min := (secOfDay % 3600) / 60
	s := secOfDay % 60

	var yearStr string
	if y >= 0 && y <= 9999 {
		yearStr = fmt.Sprintf("%04d", y)
	} else if y > 9999 {
		yearStr = fmt.Sprintf("+%d", y)
	} else {
		yearStr = fmt.Sprintf("%05d", y)
	}

	var fracStr string
	if fromInst {
		if nano > 0 {
			if nano%1000000 == 0 {
				fracStr = fmt.Sprintf(".%03d", nano/1000000)
			} else if nano%1000 == 0 {
				fracStr = fmt.Sprintf(".%06d", nano/1000)
			} else {
				fracStr = fmt.Sprintf(".%09d", nano)
			}
		}
	} else {
		fracStr = explicitFrac
	}

	return fmt.Sprintf("%s-%02d-%02dT%02d:%02d:%02d%sZ", yearStr, m, d, h, min, s, fracStr)
}

func floorDiv(x, y int64) int64 {
	r := x / y
	if (x^y) < 0 && (r*y != x) {
		r--
	}
	return r
}

func floorMod(x, y int64) int64 {
	return x - floorDiv(x, y)*y
}

func epochDayToYMD(epochDay int64) (int64, int64, int64) {
	var zeroDay = epochDay + 719528 - 60
	var adjust int64 = 0
	if zeroDay < 0 {
		var adjustCycles = (zeroDay+1)/146097 - 1
		adjust = adjustCycles * 400
		zeroDay -= adjustCycles * 146097
	}
	var yearEst = (400*zeroDay + 591) / 146097
	var doyEst = zeroDay - (365*yearEst + yearEst/4 - yearEst/100 + yearEst/400)
	if doyEst < 0 {
		yearEst--
		doyEst = zeroDay - (365*yearEst + yearEst/4 - yearEst/100 + yearEst/400)
	}
	yearEst += adjust

	var marchDay0 = doyEst
	var marchMonth0 = (marchDay0*5 + 2) / 153
	var month = (marchMonth0+2)%12 + 1
	var dom = marchDay0 - (marchMonth0*306+5)/10 + 1
	yearEst += (marchMonth0 + 2) / 12

	return yearEst, month, dom
}
