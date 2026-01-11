/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaMath

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"sync"
)

func Load_Math_Rounding_Mode() {

	ghelpers.MethodSignatures["java/math/RoundingMode.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rmodeClinit,
		}

	ghelpers.MethodSignatures["java/math/RoundingMode.clone()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/RoundingMode.compareTo(Ljava/lang/Enum;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/RoundingMode.compareTo(Ljava/lang/Object;)I"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/RoundingMode.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  rmodeEquals,
		}

	ghelpers.MethodSignatures["java/math/RoundingMode.finalize()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/RoundingMode.getClass()Ljava/lang/Class;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/RoundingMode.getDeclaringClass()Ljava/lang/Class;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/RoundingMode.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/RoundingMode.name()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rmodeName,
		}

	ghelpers.MethodSignatures["java/math/RoundingMode.notify()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/RoundingMode.notifyAll()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/RoundingMode.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/RoundingMode.valueOf(I)Ljava/math/RoundingMode;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  rmodeValueOfInt,
		}

	ghelpers.MethodSignatures["java/math/RoundingMode.valueOf(Ljava/lang/String;)Ljava/math/RoundingMode;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  rmodeValueOfString,
		}

	ghelpers.MethodSignatures["java/math/RoundingMode.values()[Ljava/math/RoundingMode;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  rmodeValues,
		}

	ghelpers.MethodSignatures["java/math/RoundingMode.wait()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/RoundingMode.wait(J)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/math/RoundingMode.wait(JI)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}
}

// Mutex for protecting the Log function during multithreading.
var roundingModeMutex = sync.Mutex{}
var rmodeOnceInitialized bool = false
var rmodeClassName = "java/math/RoundingMode"
var rmodeNames = []string{"UP", "DOWN", "CEILING", "FLOOR", "HALF_UP", "HALF_DOWN", "HALF_EVEN", "UNNECESSARY"}
var rmodeInstances []*object.Object // length 8

// ensureRoundingModeInited lazily creates singleton instances for all RoundingMode constants.
func ensureRoundingModeInited() {
	roundingModeMutex.Lock()
	defer roundingModeMutex.Unlock()
	if rmodeOnceInitialized {
		return
	}
	// Create instances in declaration order with ordinal and name fields.
	rmodeInstances = make([]*object.Object, len(rmodeNames))
	for i, nm := range rmodeNames {
		obj := object.MakeEmptyObjectWithClassName(&rmodeClassName)
		// Set minimal fields commonly present on enums: name (String) and ordinal (int)
		obj.FieldTable["name"] = object.Field{Ftype: types.StringClassRef, Fvalue: object.StringObjectFromGoString(nm)}
		obj.FieldTable["ordinal"] = object.Field{Ftype: types.Int, Fvalue: int64(i)}
		rmodeInstances[i] = obj
		// Register static field for enum constant so getstatic returns the singleton
		_ = statics.AddStatic(rmodeClassName+"."+nm, statics.Static{Type: "Ljava/math/RoundingMode;", Value: obj})
	}
	rmodeOnceInitialized = true
}

// rmodeClinit is the <clinit> for RoundingMode; it ensures constants are available for getstatic
func rmodeClinit([]interface{}) interface{} {
	ensureRoundingModeInited()
	return nil
}

// valueOf(int): Maps legacy BigDecimal rounding constants to RoundingMode enum instances
// 0:UP 1:DOWN 2:CEILING 3:FLOOR 4:HALF_UP 5:HALF_DOWN 6:HALF_EVEN 7:UNNECESSARY
func rmodeValueOfInt(params []interface{}) interface{} {
	ensureRoundingModeInited()
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "RoundingMode.valueOf(int): missing argument")
	}
	code, ok := params[0].(int64)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "RoundingMode.valueOf(int): argument is not an int")
	}
	if code < 0 || int(code) >= len(rmodeInstances) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "RoundingMode.valueOf(int): invalid rounding mode")
	}
	return rmodeInstances[code]
}

// valueOf(String): Standard Enum.valueOf behavior for RoundingMode
func rmodeValueOfString(params []interface{}) interface{} {
	ensureRoundingModeInited()
	if len(params) < 1 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "RoundingMode.valueOf(String): missing argument")
	}
	// Type and null checks
	strObj, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "RoundingMode.valueOf(String): argument is not a String")
	}
	if strObj == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "RoundingMode.valueOf(String): name is null")
	}
	// Ensure it's actually a Java String object (not just any empty object)
	if !object.IsStringObject(strObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "RoundingMode.valueOf(String): argument is not a String")
	}
	name := object.GoStringFromStringObject(strObj)
	// Match exactly (case-sensitive) like Enum.valueOf
	for i, nm := range rmodeNames {
		if nm == name {
			return rmodeInstances[i]
		}
	}
	return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "RoundingMode.valueOf(String): no enum constant")
}

// values(): returns array of all constants in declaration order
func rmodeValues(params []interface{}) interface{} {
	ensureRoundingModeInited()
	arr := object.Make1DimRefArray("Ljava/math/RoundingMode;", int64(len(rmodeInstances)))
	slot := arr.FieldTable["value"].Fvalue.([]*object.Object)
	copy(slot, rmodeInstances)
	arr.FieldTable["value"] = object.Field{Ftype: types.RefArray + "Ljava/math/RoundingMode;", Fvalue: slot}
	return arr
}

// rmodeName implements RoundingMode.name(): returns the enum constant name as a Java String
func rmodeName(params []interface{}) interface{} {
	ensureRoundingModeInited()
	if len(params) == 0 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "RoundingMode.name(): missing receiver")
	}
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "RoundingMode.name(): null receiver")
	}
	fld, ok := self.FieldTable["name"]
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "RoundingMode.name(): missing name field")
	}
	strObj, ok := fld.Fvalue.(*object.Object)
	if !ok || !object.IsStringObject(strObj) {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "RoundingMode.name(): name field is not a String")
	}
	return strObj
}

// rmodeEquals implements RoundingMode.equals(Object): reference identity for enum singletons
func rmodeEquals(params []interface{}) interface{} {
	ensureRoundingModeInited()
	// Expect: params[0] = this, params[1] = other
	if len(params) < 2 {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "RoundingMode.equals(): missing arguments")
	}
	self, ok := params[0].(*object.Object)
	if !ok || object.IsNull(self) {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "RoundingMode.equals(): null receiver")
	}
	other := params[1]
	if object.IsNull(other) {
		return types.JavaBoolFalse
	}
	otherObj, ok := other.(*object.Object)
	if !ok {
		return types.JavaBoolFalse
	}
	// Identity check: true iff same instance
	if self == otherObj {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}
