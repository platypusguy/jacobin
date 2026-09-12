/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/statics"
	"jacobin/src/types"
	"math"
	"strconv"
)

// Implementation of java/util/logging/Level.
// Strategy: Level = jacobin Object with three fields:
//   name:               Go string  -- the name of the level (e.g. "INFO")
//   value:              int64      -- the integer value of the level
//   resourceBundleName: Go string  -- the (possibly empty) resource bundle name

var levelClassName = "java/util/logging/Level"
var fieldNameLevelName = "name"
var fieldNameLevelValue = "value"
var fieldNameLevelResourceBundleName = "resourceBundleName"

// standardLevels maps the standard Level names to their integer values, used by parse().
var standardLevels = map[string]int64{
	"OFF":     int64(math.MaxInt32),
	"SEVERE":  1000,
	"WARNING": 900,
	"INFO":    800,
	"CONFIG":  700,
	"FINE":    500,
	"FINER":   400,
	"FINEST":  300,
	"ALL":     int64(math.MinInt32),
}

func Load_Util_Logging_Level() {

	ghelpers.MethodSignatures["java/util/logging/Level.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingLevelClinit,
		}

	ghelpers.MethodSignatures["java/util/logging/Level.<init>(Ljava/lang/String;I)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  loggingLevelInit,
		}

	ghelpers.MethodSignatures["java/util/logging/Level.<init>(Ljava/lang/String;ILjava/lang/String;)V"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  loggingLevelInitWithResourceBundle,
		}

	ghelpers.MethodSignatures["java/util/logging/Level.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingLevelEquals,
		}

	ghelpers.MethodSignatures["java/util/logging/Level.getLocalizedName()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingLevelGetLocalizedName,
		}

	ghelpers.MethodSignatures["java/util/logging/Level.getName()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingLevelGetName,
		}

	ghelpers.MethodSignatures["java/util/logging/Level.getResourceBundleName()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingLevelGetResourceBundleName,
		}

	ghelpers.MethodSignatures["java/util/logging/Level.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingLevelHashCode,
		}

	ghelpers.MethodSignatures["java/util/logging/Level.intValue()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingLevelIntValue,
		}

	ghelpers.MethodSignatures["java/util/logging/Level.parse(Ljava/lang/String;)Ljava/util/logging/Level;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  loggingLevelParse,
		}

	ghelpers.MethodSignatures["java/util/logging/Level.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  loggingLevelToString,
		}
}

// loggingLevelClinit is the static initializer for java/util/logging/Level. It creates
// the standard Level instances (OFF, SEVERE, WARNING, INFO, CONFIG, FINE, FINER, FINEST, ALL)
// and registers them as public static final fields, matching the real java.util.logging.Level.
// This is necessary because other classes (e.g. FileHandler, StreamHandler) access these
// fields directly via GETSTATIC.
func loggingLevelClinit(_ []interface{}) interface{} {
	for name, value := range standardLevels {
		_ = statics.AddStatic(levelClassName+"."+name, statics.Static{
			Type:  types.Ref,
			Value: makeLevelObject(name, value, ""),
		})
	}
	return nil
}

// makeLevelObject creates a Level object with the given name, value, and resource bundle name.
func makeLevelObject(name string, value int64, resourceBundleName string) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&levelClassName)
	obj.FieldTable[fieldNameLevelName] = object.Field{Ftype: types.StringClassRef, Fvalue: name}
	obj.FieldTable[fieldNameLevelValue] = object.Field{Ftype: types.Int, Fvalue: value}
	obj.FieldTable[fieldNameLevelResourceBundleName] = object.Field{Ftype: types.StringClassRef, Fvalue: resourceBundleName}
	return obj
}

// "java/util/logging/Level.<init>(Ljava/lang/String;I)V"
func loggingLevelInit(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLevelInit: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	nameObj, ok := params[1].(*object.Object)
	if !ok || nameObj == nil {
		errMsg := "loggingLevelInit: The name parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	name := object.GoStringFromStringObject(nameObj)
	value := int64(int32(params[2].(int64)))

	if obj.KlassName == 0 || obj.KlassName == types.InvalidStringIndex {
		obj.KlassName = object.StringPoolIndexFromGoString(levelClassName)
	}

	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameLevelName] = object.Field{Ftype: types.StringClassRef, Fvalue: name}
	obj.FieldTable[fieldNameLevelValue] = object.Field{Ftype: types.Int, Fvalue: value}
	obj.FieldTable[fieldNameLevelResourceBundleName] = object.Field{Ftype: types.StringClassRef, Fvalue: ""}
	return nil
}

// "java/util/logging/Level.<init>(Ljava/lang/String;ILjava/lang/String;)V"
func loggingLevelInitWithResourceBundle(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLevelInitWithResourceBundle: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	nameObj, ok := params[1].(*object.Object)
	if !ok || nameObj == nil {
		errMsg := "loggingLevelInitWithResourceBundle: The name parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	name := object.GoStringFromStringObject(nameObj)
	value := int64(int32(params[2].(int64)))

	resourceBundleName := ""
	if rbObj, ok := params[3].(*object.Object); ok && rbObj != nil && !object.IsNull(rbObj) {
		resourceBundleName = object.GoStringFromStringObject(rbObj)
	}

	if obj.KlassName == 0 || obj.KlassName == types.InvalidStringIndex {
		obj.KlassName = object.StringPoolIndexFromGoString(levelClassName)
	}

	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameLevelName] = object.Field{Ftype: types.StringClassRef, Fvalue: name}
	obj.FieldTable[fieldNameLevelValue] = object.Field{Ftype: types.Int, Fvalue: value}
	obj.FieldTable[fieldNameLevelResourceBundleName] = object.Field{Ftype: types.StringClassRef, Fvalue: resourceBundleName}
	return nil
}

// "java/util/logging/Level.getName()Ljava/lang/String;"
func loggingLevelGetName(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLevelGetName: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	name, _ := obj.FieldTable[fieldNameLevelName].Fvalue.(string)
	return object.StringObjectFromGoString(name)
}

// "java/util/logging/Level.getLocalizedName()Ljava/lang/String;"
// No localization support: the localized name is identical to the level's name.
func loggingLevelGetLocalizedName(params []interface{}) interface{} {
	return loggingLevelGetName(params)
}

// "java/util/logging/Level.getResourceBundleName()Ljava/lang/String;"
func loggingLevelGetResourceBundleName(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLevelGetResourceBundleName: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	rbName, _ := obj.FieldTable[fieldNameLevelResourceBundleName].Fvalue.(string)
	if rbName == "" {
		return object.Null
	}
	return object.StringObjectFromGoString(rbName)
}

// "java/util/logging/Level.intValue()I"
func loggingLevelIntValue(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := fmt.Sprintf("loggingLevelIntValue: The first parameter is not an object, observed type: %T", params[0])
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	value, _ := obj.FieldTable[fieldNameLevelValue].Fvalue.(int64)
	return value
}

// "java/util/logging/Level.toString()Ljava/lang/String;"
// Level.toString() returns the same value as getName().
func loggingLevelToString(params []interface{}) interface{} {
	return loggingLevelGetName(params)
}

// "java/util/logging/Level.hashCode()I"
// Level.hashCode() returns the same value as intValue().
func loggingLevelHashCode(params []interface{}) interface{} {
	return loggingLevelIntValue(params)
}

// "java/util/logging/Level.equals(Ljava/lang/Object;)Z"
func loggingLevelEquals(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLevelEquals: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	other, ok := params[1].(*object.Object)
	if !ok || other == nil || object.IsNull(other) {
		return types.JavaBoolFalse
	}

	otherValue, ok := other.FieldTable[fieldNameLevelValue].Fvalue.(int64)
	if !ok {
		return types.JavaBoolFalse
	}

	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	thisValue, _ := obj.FieldTable[fieldNameLevelValue].Fvalue.(int64)

	if thisValue == otherValue {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// "java/util/logging/Level.parse(Ljava/lang/String;)Ljava/util/logging/Level;"
func loggingLevelParse(params []interface{}) interface{} {
	nameObj, ok := params[0].(*object.Object)
	if !ok || nameObj == nil {
		errMsg := "loggingLevelParse: The parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	name := object.GoStringFromStringObject(nameObj)

	if value, found := standardLevels[name]; found {
		return makeLevelObject(name, value, "")
	}

	// Not a standard level name: try to parse it as an integer value.
	value, err := strconv.ParseInt(name, 10, 64)
	if err != nil {
		errMsg := "loggingLevelParse: Bad level \"" + name + "\""
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	return makeLevelObject(name, value, "")
}
