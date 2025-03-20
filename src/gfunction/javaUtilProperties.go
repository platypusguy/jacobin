/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/excNames"
	"jacobin/globals"
	"jacobin/object"
	"jacobin/types"
	"sync"
)

// Implementation of some of the functions in Java/util/Locale.
// Strategy: Locale = jacobin Object wrapping a Go string.

func Load_Util_Properties() {

	MethodSignatures["java/util/Properties.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/util/Properties.<init>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  propertiesInit,
		}

	MethodSignatures["java/util/Properties.<init>(I)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  propertiesInit,
		}

	MethodSignatures["java/util/Properties.<init>(Ljava/util/Properties;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Properties.clear()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  propertiesInit,
		}

	MethodSignatures["java/util/Properties.getProperty(Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  propertiesGetProperty,
		}

	MethodSignatures["java/util/Properties.getProperty(Ljava/lang/String;Ljava/lang/String;)Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  propertiesGetProperty,
		}

	MethodSignatures["java/util/Properties.list(Ljava/io/PrintStream;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Properties.list(Ljava/io/PrintWriter;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Properties.load(Ljava/io/InputStream;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Properties.load(Ljava/io/Reader;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Properties.loadFromXML(Ljava/io/InputStream;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Properties.propertyNames()Ljava/util/Enumeration;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Properties.remove(Ljava/lang/Object;)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  propertiesRemove,
		}

	MethodSignatures["java/util/Properties.save(Ljava/io/OutputStream;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Properties.setProperty(Ljava/lang/String;Ljava/lang/String;)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  propertiesSetProperty,
		}

	MethodSignatures["java/util/Properties.size()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  propertiesSize,
		}

	MethodSignatures["java/util/Properties.store(Ljava/io/OutputStream;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Properties.store(Ljava/io/Writer;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Properties.storeToXML(Ljava/io/OutputStream;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Properties.storeToXML(Ljava/io/OutputStream;Ljava/lang/String;Ljava/lang/String;)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Properties.storeToXML(Ljava/io/OutputStream;Ljava/lang/String;Ljava/nio/charset/Charset;)V"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Properties.stringPropertyNames()Ljava/util/Set;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Properties.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  propertiesToString,
		}

}

var classNameProperties = "java/util/Properties"
var fieldNameProperties = "map"
var propertiesMutex = sync.RWMutex{}

func propertiesInit(params []interface{}) interface{} {
	propertiesMutex.Lock()
	defer propertiesMutex.Unlock()

	nilMap := make(types.DefProperties)
	obj, ok := params[0].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("Properties.<init>: Properties object is invalid: %T", params[0])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	object.ClearFieldTable(obj)
	fld := obj.FieldTable[fieldNameProperties]
	fld.Ftype = types.Properties
	fld.Fvalue = nilMap
	obj.FieldTable[fieldNameProperties] = fld
	return nil
}

// Given a properties table and a key, retrieve the associated value.
func propertiesGetProperty(params []interface{}) interface{} {
	// Get properties table.
	this, ok := params[0].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("propertiesGetProperty: Properties object is invalid: %T", params[0])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	properties, ok := this.FieldTable[fieldNameProperties].Fvalue.(types.DefProperties)
	if !ok {
		errMsg := "propertiesGetProperty: Properties table is missing or invalid"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get key.
	keyObj, ok := params[1].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("propertiesGetProperty: Key object is invalid: %T", params[1])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	key := object.GoStringFromStringObject(keyObj)
	if len(key) == 0 {
		errMsg := "propertiesGetProperty: Key parameter is not String or is null"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get default value.
	flagDefault := false
	var dfltValue string
	if len(params) > 2 {
		dfltObj, ok := params[2].(*object.Object)
		if ok {
			dfltValue = object.GoStringFromStringObject(dfltObj)
			if len(key) == 0 {
				errMsg := "propertiesGetProperty: Default value parameter is not String or is null"
				return getGErrBlk(excNames.IllegalArgumentException, errMsg)
			}
			flagDefault = true
		}
	}

	// Get value associated with key.
	// If present, return it as a JavaString.
	value, ok := properties[key]
	if ok {
		return object.StringObjectFromGoString(value)
	}

	// Value is not present.
	// If default value supplied, return it.
	// Otherwise, return null.
	if flagDefault {
		return object.StringObjectFromGoString(dfltValue)
	}
	return object.Null
}

// Given a properties table and a key, remove this entry.
func propertiesRemove(params []interface{}) interface{} {
	propertiesMutex.Lock()
	defer propertiesMutex.Unlock()

	// Get properties table.
	this, ok := params[0].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("propertiesRemove: Properties object is invalid: %T", params[0])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	fld, ok := this.FieldTable[fieldNameProperties]
	if !ok {
		errMsg := "propertiesRemove: Properties table is missing or invalid"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	properties := fld.Fvalue.(types.DefProperties)

	// Get key.
	keyObj, ok := params[1].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("propertiesRemove: Key object is invalid: %T", params[1])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	key := object.GoStringFromStringObject(keyObj)
	if len(key) == 0 {
		errMsg := "propertiesRemove: Key parameter is not String or is null"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Return value.
	flagReturnValue := false
	oldValue, ok := properties[key]
	if ok {
		flagReturnValue = true
	}

	// Remove entry associated with key.
	delete(properties, key)
	fld.Fvalue = properties
	this.FieldTable[fieldNameProperties] = fld

	// If there was an old value, return it as a Java String.
	// Otherwise, return null.
	if flagReturnValue {
		object.StringObjectFromGoString(oldValue)
	}
	return object.Null
}

// Given a properties table and a key, set its entry ith the specified value.
func propertiesSetProperty(params []interface{}) interface{} {
	propertiesMutex.Lock()
	defer propertiesMutex.Unlock()

	// Get properties table.
	this, ok := params[0].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("propertiesSetProperty: Properties object is invalid: %T", params[0])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	fld, ok := this.FieldTable[fieldNameProperties]
	if !ok {
		errMsg := "propertiesSetProperty: properties table is missing or invalid"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	properties := fld.Fvalue.(types.DefProperties)

	// Get key.
	keyObj, ok := params[1].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("propertiesSetProperty: Key object is invalid: %T", params[1])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	key := object.GoStringFromStringObject(keyObj)
	if len(key) == 0 {
		errMsg := "propertiesGetProperty: Key parameter is not String or is null"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get new value.
	valueObj, ok := params[2].(*object.Object)
	if !ok {
		errMsg := "propertiesSetProperty: Value parameter is not an object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	value := object.GoStringFromStringObject(valueObj)
	if len(value) == 0 {
		errMsg := "propertiesSetProperty: Value parameter is not String or is null"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get old value if present.
	oldValue, oldPresent := properties[key]

	// Set entry key, value.
	properties[key] = value
	fld.Fvalue = properties
	this.FieldTable[fieldNameProperties] = fld

	// If there was an old value, return it. Otherwise, return Null.
	if oldPresent {
		return object.StringObjectFromGoString(oldValue)
	}
	return object.Null
}

// Given a properties table, compute the number of entries.
func propertiesSize(params []interface{}) interface{} {
	// Get properties table.
	this, ok := params[0].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("propertiesSize: Properties object is invalid: %T", params[0])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	properties, ok := this.FieldTable[fieldNameProperties].Fvalue.(types.DefProperties)
	if !ok {
		errMsg := "propertiesGetProperty: properties table is missing or invalid"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Return the number of entries.
	return int64(len(properties))
}

// Given a properties table, return a string representation of this Properties object
// in the form of a set of entries, enclosed in braces and separated by the ASCII characters " , " (comma and space).
// Each entry is rendered as the key, an equals sign =, and the associated element.
func propertiesToString(params []interface{}) interface{} {
	// Get properties table.
	this, ok := params[0].(*object.Object)
	if !ok {
		errMsg := fmt.Sprintf("propertiesToString: Properties object is invalid: %T", params[0])
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	properties, ok := this.FieldTable[fieldNameProperties].Fvalue.(types.DefProperties)
	if !ok {
		errMsg := "propertiesToString: properties table is missing or invalid"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Create a slice of keys.
	keys := make([]string, 0, len(properties))
	for key := range properties {
		keys = append(keys, key)
	}

	// Sort the keys, case-insensitive.
	globals.SortCaseInsensitive(&keys)

	// Build longString, consisting of key-value pairs.
	var longString = "{"
	for _, key := range keys {
		value := properties[key]
		longString += key + "=" + value + ", "
	}
	longString = longString[:len(longString)-2] + "}"

	// Return longString as a Java String.
	return object.StringObjectFromGoString(longString)
}
