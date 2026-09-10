/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by the Jacobin authors. Consult jacobin.org.
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
	"os"
	"strings"
	"sync"
)

// Implementation of java/util/logging/Logger.
// Strategy: Logger = jacobin Object with fields:
//   name:               string (java/lang/String) -- the logger's name, or "" for anonymous loggers
//   level:              *object.Object (java/util/logging/Level) -- nil/object.Null means "inherit from parent"
//   parent:             *object.Object (java/util/logging/Logger)
//   filter:             *object.Object (java/util/logging/Filter)
//   useParentHandlers:  int64 (boolean) -- whether to also send output to the parent's handlers
//   handlers:           []*object.Object (java/util/logging/Handler)
//   resourceBundleName: string
//   resourceBundle:     *object.Object (java/util/ResourceBundle)
// Named loggers are cached in loggerRegistry so that getLogger(name) always
// returns the same instance for the same name, matching real Logger semantics.
// The message-formatting methods (config/info/severe/log/logp/logrb/etc.) write
// directly to System.err, following the same "publish to System.err" pattern
// used by Handler/ConsoleHandler in this package.

var loggerClassName = "java/util/logging/Logger"
var fieldNameLoggerName = "name"
var fieldNameLoggerParent = "parent"
var fieldNameLoggerFilter = "filter"
var fieldNameLoggerUseParentHandlers = "useParentHandlers"
var fieldNameLoggerHandlers = "handlers"
var fieldNameLoggerResourceBundleName = "resourceBundleName"
var fieldNameLoggerResourceBundle = "resourceBundle"

var loggerRegistry = map[string]*object.Object{}
var loggerRegistryMutex sync.Mutex
var globalLoggerObj *object.Object
var globalLoggerMutex sync.Mutex

func Load_Util_Logging_Logger() {

	ghelpers.MethodSignatures["java/util/logging/Logger.<clinit>()V"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: ghelpers.ClinitGeneric}

	ghelpers.MethodSignatures["java/util/logging/Logger.<init>(Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: loggingLoggerInit}

	ghelpers.MethodSignatures["java/util/logging/Logger.addHandler(Ljava/util/logging/Handler;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: loggingLoggerAddHandler}

	ghelpers.MethodSignatures["java/util/logging/Logger.config(Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: loggingLoggerConfig}

	ghelpers.MethodSignatures["java/util/logging/Logger.entering(Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: loggingLoggerEntering}

	ghelpers.MethodSignatures["java/util/logging/Logger.entering(Ljava/lang/String;Ljava/lang/String;Ljava/lang/Object;)V"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: loggingLoggerEnteringWithParam}

	ghelpers.MethodSignatures["java/util/logging/Logger.exiting(Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: loggingLoggerExiting}

	ghelpers.MethodSignatures["java/util/logging/Logger.exiting(Ljava/lang/String;Ljava/lang/String;Ljava/lang/Object;)V"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: loggingLoggerExitingWithParam}

	ghelpers.MethodSignatures["java/util/logging/Logger.fine(Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: loggingLoggerFine}

	ghelpers.MethodSignatures["java/util/logging/Logger.finer(Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: loggingLoggerFiner}

	ghelpers.MethodSignatures["java/util/logging/Logger.finest(Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: loggingLoggerFinest}

	ghelpers.MethodSignatures["java/util/logging/Logger.getAnonymousLogger()Ljava/util/logging/Logger;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: loggingLoggerGetAnonymousLogger}

	ghelpers.MethodSignatures["java/util/logging/Logger.getAnonymousLogger(Ljava/lang/String;)Ljava/util/logging/Logger;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: loggingLoggerGetAnonymousLoggerWithBundle}

	ghelpers.MethodSignatures["java/util/logging/Logger.getFilter()Ljava/util/logging/Filter;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: loggingLoggerGetFilter}

	ghelpers.MethodSignatures["java/util/logging/Logger.getGlobal()Ljava/util/logging/Logger;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: loggingLoggerGetGlobal}

	ghelpers.MethodSignatures["java/util/logging/Logger.getHandlers()[Ljava/util/logging/Handler;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: loggingLoggerGetHandlers}

	ghelpers.MethodSignatures["java/util/logging/Logger.getLevel()Ljava/util/logging/Level;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: loggingLoggerGetLevel}

	ghelpers.MethodSignatures["java/util/logging/Logger.getLogger(Ljava/lang/String;)Ljava/util/logging/Logger;"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: loggingLoggerGetLogger}

	ghelpers.MethodSignatures["java/util/logging/Logger.getLogger(Ljava/lang/String;Ljava/lang/String;)Ljava/util/logging/Logger;"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: loggingLoggerGetLoggerWithBundle}

	ghelpers.MethodSignatures["java/util/logging/Logger.getName()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: loggingLoggerGetName}

	ghelpers.MethodSignatures["java/util/logging/Logger.getParent()Ljava/util/logging/Logger;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: loggingLoggerGetParent}

	ghelpers.MethodSignatures["java/util/logging/Logger.getResourceBundle()Ljava/util/ResourceBundle;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: loggingLoggerGetResourceBundle}

	ghelpers.MethodSignatures["java/util/logging/Logger.getResourceBundleName()Ljava/lang/String;"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: loggingLoggerGetResourceBundleName}

	ghelpers.MethodSignatures["java/util/logging/Logger.getUseParentHandlers()Z"] =
		ghelpers.GMeth{ParamSlots: 0, GFunction: loggingLoggerGetUseParentHandlers}

	ghelpers.MethodSignatures["java/util/logging/Logger.info(Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: loggingLoggerInfo}

	ghelpers.MethodSignatures["java/util/logging/Logger.isLoggable(Ljava/util/logging/Level;)Z"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: loggingLoggerIsLoggable}

	ghelpers.MethodSignatures["java/util/logging/Logger.log(Ljava/util/logging/Level;Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: loggingLoggerLog}

	ghelpers.MethodSignatures["java/util/logging/Logger.log(Ljava/util/logging/Level;Ljava/lang/String;[Ljava/lang/Object;)V"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: loggingLoggerLogWithParams}

	ghelpers.MethodSignatures["java/util/logging/Logger.logp(Ljava/util/logging/Level;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 4, GFunction: loggingLoggerLogp}

	ghelpers.MethodSignatures["java/util/logging/Logger.logp(Ljava/util/logging/Level;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;[Ljava/lang/Object;)V"] =
		ghelpers.GMeth{ParamSlots: 5, GFunction: loggingLoggerLogpWithParams}

	ghelpers.MethodSignatures["java/util/logging/Logger.logrb(Ljava/util/logging/Level;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 5, GFunction: loggingLoggerLogrb}

	ghelpers.MethodSignatures["java/util/logging/Logger.logrb(Ljava/util/logging/Level;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;[Ljava/lang/Object;)V"] =
		ghelpers.GMeth{ParamSlots: 6, GFunction: loggingLoggerLogrbWithParams}

	ghelpers.MethodSignatures["java/util/logging/Logger.removeHandler(Ljava/util/logging/Handler;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: loggingLoggerRemoveHandler}

	ghelpers.MethodSignatures["java/util/logging/Logger.setFilter(Ljava/util/logging/Filter;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: loggingLoggerSetFilter}

	ghelpers.MethodSignatures["java/util/logging/Logger.setLevel(Ljava/util/logging/Level;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: loggingLoggerSetLevel}

	ghelpers.MethodSignatures["java/util/logging/Logger.setParent(Ljava/util/logging/Logger;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: loggingLoggerSetParent}

	ghelpers.MethodSignatures["java/util/logging/Logger.setUseParentHandlers(Z)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: loggingLoggerSetUseParentHandlers}

	ghelpers.MethodSignatures["java/util/logging/Logger.severe(Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: loggingLoggerSevere}

	ghelpers.MethodSignatures["java/util/logging/Logger.throwing(Ljava/lang/String;Ljava/lang/String;Ljava/lang/Throwable;)V"] =
		ghelpers.GMeth{ParamSlots: 3, GFunction: loggingLoggerThrowing}

	ghelpers.MethodSignatures["java/util/logging/Logger.warning(Ljava/lang/String;)V"] =
		ghelpers.GMeth{ParamSlots: 1, GFunction: loggingLoggerWarning}

	ghelpers.MethodSignatures["java/util/logging/Logger.warning(Ljava/lang/String;[Ljava/lang/Object;)V"] =
		ghelpers.GMeth{ParamSlots: 2, GFunction: loggingLoggerWarningWithParams}
}

// makeLoggerObject creates a new Logger object with the given name and resource
// bundle name, and default state (no level set, no parent, no filter, useParentHandlers
// true, no handlers, no resource bundle).
func makeLoggerObject(name, resourceBundleName string) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&loggerClassName)
	obj.FieldTable[fieldNameLoggerName] = object.Field{Ftype: types.StringClassRef, Fvalue: name}
	obj.FieldTable[fieldNameHandlerLevel] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameLoggerParent] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameLoggerFilter] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameLoggerUseParentHandlers] = object.Field{Ftype: types.Int, Fvalue: types.JavaBoolTrue}
	obj.FieldTable[fieldNameLoggerHandlers] = object.Field{Ftype: types.Ref, Fvalue: []*object.Object{}}
	obj.FieldTable[fieldNameLoggerResourceBundleName] = object.Field{Ftype: types.StringClassRef, Fvalue: resourceBundleName}
	obj.FieldTable[fieldNameLoggerResourceBundle] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	return obj
}

// "java/util/logging/Logger.<init>(Ljava/lang/String;Ljava/lang/String;)V"
// Protected constructor, used by subclasses. name may be null (anonymous logger).
func loggingLoggerInit(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLoggerInit: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	name := loggingLoggerStringArg(params[1])
	resourceBundleName := loggingLoggerStringArg(params[2])

	if obj.KlassName == 0 || obj.KlassName == types.InvalidStringIndex {
		obj.KlassName = object.StringPoolIndexFromGoString(loggerClassName)
	}

	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameLoggerName] = object.Field{Ftype: types.StringClassRef, Fvalue: name}
	obj.FieldTable[fieldNameHandlerLevel] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameLoggerParent] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameLoggerFilter] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	obj.FieldTable[fieldNameLoggerUseParentHandlers] = object.Field{Ftype: types.Int, Fvalue: types.JavaBoolTrue}
	obj.FieldTable[fieldNameLoggerHandlers] = object.Field{Ftype: types.Ref, Fvalue: []*object.Object{}}
	obj.FieldTable[fieldNameLoggerResourceBundleName] = object.Field{Ftype: types.StringClassRef, Fvalue: resourceBundleName}
	obj.FieldTable[fieldNameLoggerResourceBundle] = object.Field{Ftype: types.Ref, Fvalue: object.Null}
	return nil
}

// "java/util/logging/Logger.getLogger(Ljava/lang/String;)Ljava/util/logging/Logger;"
func loggingLoggerGetLogger(params []interface{}) interface{} {
	return loggingLoggerGetOrCreate(loggingLoggerStringArg(params[0]), "")
}

// "java/util/logging/Logger.getLogger(Ljava/lang/String;Ljava/lang/String;)Ljava/util/logging/Logger;"
func loggingLoggerGetLoggerWithBundle(params []interface{}) interface{} {
	return loggingLoggerGetOrCreate(loggingLoggerStringArg(params[0]), loggingLoggerStringArg(params[1]))
}

// loggingLoggerGetOrCreate returns the cached Logger for name, creating it (and
// registering it with the root/global logger as parent) if it doesn't yet exist.
func loggingLoggerGetOrCreate(name, resourceBundleName string) *object.Object {
	loggerRegistryMutex.Lock()
	defer loggerRegistryMutex.Unlock()

	if existing, found := loggerRegistry[name]; found {
		return existing
	}

	logger := makeLoggerObject(name, resourceBundleName)
	logger.FieldTable[fieldNameLoggerParent] = object.Field{Ftype: types.Ref, Fvalue: loggingLoggerGlobal()}
	loggerRegistry[name] = logger
	return logger
}

// loggingLoggerGlobal lazily creates and returns the shared global Logger instance.
func loggingLoggerGlobal() *object.Object {
	globalLoggerMutex.Lock()
	defer globalLoggerMutex.Unlock()
	if globalLoggerObj == nil {
		globalLoggerObj = makeLoggerObject("global", "")
		globalLoggerObj.FieldTable[fieldNameHandlerLevel] = object.Field{Ftype: types.Ref, Fvalue: makeLevelObject("INFO", standardLevels["INFO"], "")}
	}
	return globalLoggerObj
}

// "java/util/logging/Logger.getGlobal()Ljava/util/logging/Logger;"
func loggingLoggerGetGlobal([]interface{}) interface{} {
	return loggingLoggerGlobal()
}

// "java/util/logging/Logger.getAnonymousLogger()Ljava/util/logging/Logger;"
func loggingLoggerGetAnonymousLogger([]interface{}) interface{} {
	logger := makeLoggerObject("", "")
	logger.FieldTable[fieldNameLoggerParent] = object.Field{Ftype: types.Ref, Fvalue: loggingLoggerGlobal()}
	return logger
}

// "java/util/logging/Logger.getAnonymousLogger(Ljava/lang/String;)Ljava/util/logging/Logger;"
func loggingLoggerGetAnonymousLoggerWithBundle(params []interface{}) interface{} {
	logger := makeLoggerObject("", loggingLoggerStringArg(params[0]))
	logger.FieldTable[fieldNameLoggerParent] = object.Field{Ftype: types.Ref, Fvalue: loggingLoggerGlobal()}
	return logger
}

// "java/util/logging/Logger.getName()Ljava/lang/String;"
func loggingLoggerGetName(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLoggerGetName: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	name, _ := obj.FieldTable[fieldNameLoggerName].Fvalue.(string)
	if name == "" {
		return object.Null
	}
	return object.StringObjectFromGoString(name)
}

// "java/util/logging/Logger.getLevel()Ljava/util/logging/Level;"
func loggingLoggerGetLevel(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLoggerGetLevel: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	level, ok := obj.FieldTable[fieldNameHandlerLevel].Fvalue.(*object.Object)
	if !ok || level == nil {
		return object.Null
	}
	return level
}

// "java/util/logging/Logger.setLevel(Ljava/util/logging/Level;)V"
func loggingLoggerSetLevel(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLoggerSetLevel: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameHandlerLevel] = object.Field{Ftype: types.Ref, Fvalue: params[1]}
	return nil
}

// loggingLoggerEffectiveLevel walks up the parent chain to find the first
// explicitly-set Level, defaulting to Level.INFO if none is found.
func loggingLoggerEffectiveLevel(obj *object.Object) *object.Object {
	current := obj
	for i := 0; current != nil && !object.IsNull(current) && i < 64; i++ {
		current.ThMutex.RLock()
		level, hasLevel := current.FieldTable[fieldNameHandlerLevel].Fvalue.(*object.Object)
		parent, _ := current.FieldTable[fieldNameLoggerParent].Fvalue.(*object.Object)
		current.ThMutex.RUnlock()
		if hasLevel && level != nil && !object.IsNull(level) {
			return level
		}
		current = parent
	}
	return makeLevelObject("INFO", standardLevels["INFO"], "")
}

// "java/util/logging/Logger.getParent()Ljava/util/logging/Logger;"
func loggingLoggerGetParent(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLoggerGetParent: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	parent, ok := obj.FieldTable[fieldNameLoggerParent].Fvalue.(*object.Object)
	if !ok || parent == nil {
		return object.Null
	}
	return parent
}

// "java/util/logging/Logger.setParent(Ljava/util/logging/Logger;)V"
func loggingLoggerSetParent(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLoggerSetParent: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameLoggerParent] = object.Field{Ftype: types.Ref, Fvalue: params[1]}
	return nil
}

// "java/util/logging/Logger.getFilter()Ljava/util/logging/Filter;"
func loggingLoggerGetFilter(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLoggerGetFilter: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	filter, ok := obj.FieldTable[fieldNameLoggerFilter].Fvalue.(*object.Object)
	if !ok || filter == nil {
		return object.Null
	}
	return filter
}

// "java/util/logging/Logger.setFilter(Ljava/util/logging/Filter;)V"
func loggingLoggerSetFilter(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLoggerSetFilter: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameLoggerFilter] = object.Field{Ftype: types.Ref, Fvalue: params[1]}
	return nil
}

// "java/util/logging/Logger.getUseParentHandlers()Z"
func loggingLoggerGetUseParentHandlers(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLoggerGetUseParentHandlers: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	use, ok := obj.FieldTable[fieldNameLoggerUseParentHandlers].Fvalue.(int64)
	if !ok {
		return types.JavaBoolTrue
	}
	return use
}

// "java/util/logging/Logger.setUseParentHandlers(Z)V"
func loggingLoggerSetUseParentHandlers(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLoggerSetUseParentHandlers: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	use, _ := params[1].(int64)
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	obj.FieldTable[fieldNameLoggerUseParentHandlers] = object.Field{Ftype: types.Int, Fvalue: use}
	return nil
}

// "java/util/logging/Logger.getResourceBundleName()Ljava/lang/String;"
func loggingLoggerGetResourceBundleName(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLoggerGetResourceBundleName: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	rbName, _ := obj.FieldTable[fieldNameLoggerResourceBundleName].Fvalue.(string)
	if rbName == "" {
		return object.Null
	}
	return object.StringObjectFromGoString(rbName)
}

// "java/util/logging/Logger.getResourceBundle()Ljava/util/ResourceBundle;"
func loggingLoggerGetResourceBundle(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLoggerGetResourceBundle: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	rb, ok := obj.FieldTable[fieldNameLoggerResourceBundle].Fvalue.(*object.Object)
	if !ok || rb == nil {
		return object.Null
	}
	return rb
}

// "java/util/logging/Logger.addHandler(Ljava/util/logging/Handler;)V"
func loggingLoggerAddHandler(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLoggerAddHandler: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	handler, ok := params[1].(*object.Object)
	if !ok || handler == nil {
		errMsg := "loggingLoggerAddHandler: The handler parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	handlers, _ := obj.FieldTable[fieldNameLoggerHandlers].Fvalue.([]*object.Object)
	handlers = append(handlers, handler)
	obj.FieldTable[fieldNameLoggerHandlers] = object.Field{Ftype: types.Ref, Fvalue: handlers}
	return nil
}

// "java/util/logging/Logger.removeHandler(Ljava/util/logging/Handler;)V"
func loggingLoggerRemoveHandler(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLoggerRemoveHandler: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	handler, ok := params[1].(*object.Object)
	if !ok || handler == nil {
		return nil
	}
	obj.ThMutex.Lock()
	defer obj.ThMutex.Unlock()
	handlers, _ := obj.FieldTable[fieldNameLoggerHandlers].Fvalue.([]*object.Object)
	filtered := make([]*object.Object, 0, len(handlers))
	for _, h := range handlers {
		if h != handler {
			filtered = append(filtered, h)
		}
	}
	obj.FieldTable[fieldNameLoggerHandlers] = object.Field{Ftype: types.Ref, Fvalue: filtered}
	return nil
}

// "java/util/logging/Logger.getHandlers()[Ljava/util/logging/Handler;"
func loggingLoggerGetHandlers(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLoggerGetHandlers: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	obj.ThMutex.RLock()
	defer obj.ThMutex.RUnlock()
	handlers, _ := obj.FieldTable[fieldNameLoggerHandlers].Fvalue.([]*object.Object)
	result := object.MakeEmptyObject()
	result.FieldTable["value"] = object.Field{Ftype: types.Ref, Fvalue: handlers}
	return result
}

// "java/util/logging/Logger.isLoggable(Ljava/util/logging/Level;)Z"
func loggingLoggerIsLoggable(params []interface{}) interface{} {
	obj, ok := params[0].(*object.Object)
	if !ok || obj == nil {
		errMsg := "loggingLoggerIsLoggable: The first parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	msgLevel, ok := params[1].(*object.Object)
	if !ok || msgLevel == nil || object.IsNull(msgLevel) {
		return types.JavaBoolFalse
	}
	msgValue, _ := msgLevel.FieldTable[fieldNameLevelValue].Fvalue.(int64)

	effective := loggingLoggerEffectiveLevel(obj)
	effectiveValue, _ := effective.FieldTable[fieldNameLevelValue].Fvalue.(int64)

	if msgValue >= effectiveValue {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

// "java/util/logging/Logger.fine(Ljava/lang/String;)V"
func loggingLoggerFine(params []interface{}) interface{} {
	return loggingLoggerWrite("FINE", loggingLoggerStringArg(params[1]))
}

// loggingLoggerStringArg extracts a Go string from a java/lang/String parameter.
func loggingLoggerStringArg(param interface{}) string {
	strObj, ok := param.(*object.Object)
	if !ok || strObj == nil || object.IsNull(strObj) {
		return ""
	}
	return object.GoStringFromStringObject(strObj)
}

// loggingLoggerWrite formats a level-tagged message and writes it to System.err.
func loggingLoggerWrite(levelPrefix, msg string) interface{} {
	stderr, ok := statics.GetStaticValue("java/lang/System", "err").(*os.File)
	if !ok || stderr == nil {
		errMsg := "loggingLoggerWrite: could not obtain System.err"
		return ghelpers.GetGErrBlk(excNames.IOException, errMsg)
	}
	_, _ = fmt.Fprintf(stderr, "%s: %s\n", levelPrefix, msg)
	return nil
}

// loggingLoggerFormatParams renders an array of Java Objects as a comma-separated
// string, for use with the log/logp/logrb/warning variants that take Object[] params.
func loggingLoggerFormatParams(param interface{}) string {
	paramsObj, ok := param.(*object.Object)
	if !ok || paramsObj == nil || object.IsNull(paramsObj) {
		return ""
	}
	arrField, exists := paramsObj.FieldTable["value"]
	if !exists {
		return ""
	}
	arr, ok := arrField.Fvalue.([]*object.Object)
	if !ok {
		return ""
	}
	parts := make([]string, 0, len(arr))
	for _, elem := range arr {
		if elem == nil || object.IsNull(elem) {
			parts = append(parts, "null")
			continue
		}
		parts = append(parts, object.GoStringFromStringObject(elem))
	}
	return strings.Join(parts, ", ")
}

// "java/util/logging/Logger.config(Ljava/lang/String;)V"
func loggingLoggerConfig(params []interface{}) interface{} {
	return loggingLoggerWrite("CONFIG", loggingLoggerStringArg(params[1]))
}

// "java/util/logging/Logger.entering(Ljava/lang/String;Ljava/lang/String;)V"
func loggingLoggerEntering(params []interface{}) interface{} {
	sourceClass := loggingLoggerStringArg(params[1])
	sourceMethod := loggingLoggerStringArg(params[2])
	return loggingLoggerWrite("FINER", fmt.Sprintf("ENTRY %s %s", sourceClass, sourceMethod))
}

// "java/util/logging/Logger.entering(Ljava/lang/String;Ljava/lang/String;Ljava/lang/Object;)V"
func loggingLoggerEnteringWithParam(params []interface{}) interface{} {
	sourceClass := loggingLoggerStringArg(params[1])
	sourceMethod := loggingLoggerStringArg(params[2])
	param := loggingLoggerFormatParams(params[3])
	return loggingLoggerWrite("FINER", fmt.Sprintf("ENTRY %s %s %s", sourceClass, sourceMethod, param))
}

// "java/util/logging/Logger.exiting(Ljava/lang/String;Ljava/lang/String;)V"
func loggingLoggerExiting(params []interface{}) interface{} {
	sourceClass := loggingLoggerStringArg(params[1])
	sourceMethod := loggingLoggerStringArg(params[2])
	return loggingLoggerWrite("FINER", fmt.Sprintf("RETURN %s %s", sourceClass, sourceMethod))
}

// "java/util/logging/Logger.exiting(Ljava/lang/String;Ljava/lang/String;Ljava/lang/Object;)V"
func loggingLoggerExitingWithParam(params []interface{}) interface{} {
	sourceClass := loggingLoggerStringArg(params[1])
	sourceMethod := loggingLoggerStringArg(params[2])
	result := loggingLoggerStringArg(params[3])
	return loggingLoggerWrite("FINER", fmt.Sprintf("RETURN %s %s %s", sourceClass, sourceMethod, result))
}

// "java/util/logging/Logger.finest(Ljava/lang/String;)V"
func loggingLoggerFinest(params []interface{}) interface{} {
	return loggingLoggerWrite("FINEST", loggingLoggerStringArg(params[1]))
}

// "java/util/logging/Logger.finer(Ljava/lang/String;)V"
func loggingLoggerFiner(params []interface{}) interface{} {
	return loggingLoggerWrite("FINER", loggingLoggerStringArg(params[1]))
}

// "java/util/logging/Logger.info(Ljava/lang/String;)V"
func loggingLoggerInfo(params []interface{}) interface{} {
	return loggingLoggerWrite("INFO", loggingLoggerStringArg(params[1]))
}

// loggingLoggerLevelName returns the name of the given Level object, or "INFO" if unavailable.
func loggingLoggerLevelName(param interface{}) string {
	levelObj, ok := param.(*object.Object)
	if !ok || levelObj == nil || object.IsNull(levelObj) {
		return "INFO"
	}
	name, ok := levelObj.FieldTable[fieldNameLevelName].Fvalue.(string)
	if !ok || name == "" {
		return "INFO"
	}
	return name
}

// "java/util/logging/Logger.log(Ljava/util/logging/Level;Ljava/lang/String;)V"
func loggingLoggerLog(params []interface{}) interface{} {
	levelName := loggingLoggerLevelName(params[1])
	msg := loggingLoggerStringArg(params[2])
	return loggingLoggerWrite(levelName, msg)
}

// "java/util/logging/Logger.log(Ljava/util/logging/Level;Ljava/lang/String;[Ljava/lang/Object;)V"
func loggingLoggerLogWithParams(params []interface{}) interface{} {
	levelName := loggingLoggerLevelName(params[1])
	msg := loggingLoggerStringArg(params[2])
	extra := loggingLoggerFormatParams(params[3])
	if extra != "" {
		msg = msg + " [" + extra + "]"
	}
	return loggingLoggerWrite(levelName, msg)
}

// "java/util/logging/Logger.logp(Ljava/util/logging/Level;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V"
func loggingLoggerLogp(params []interface{}) interface{} {
	levelName := loggingLoggerLevelName(params[1])
	sourceClass := loggingLoggerStringArg(params[2])
	sourceMethod := loggingLoggerStringArg(params[3])
	msg := loggingLoggerStringArg(params[4])
	return loggingLoggerWrite(levelName, fmt.Sprintf("%s %s: %s", sourceClass, sourceMethod, msg))
}

// "java/util/logging/Logger.logp(Ljava/util/logging/Level;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;[Ljava/lang/Object;)V"
func loggingLoggerLogpWithParams(params []interface{}) interface{} {
	levelName := loggingLoggerLevelName(params[1])
	sourceClass := loggingLoggerStringArg(params[2])
	sourceMethod := loggingLoggerStringArg(params[3])
	msg := loggingLoggerStringArg(params[4])
	extra := loggingLoggerFormatParams(params[5])
	if extra != "" {
		msg = msg + " [" + extra + "]"
	}
	return loggingLoggerWrite(levelName, fmt.Sprintf("%s %s: %s", sourceClass, sourceMethod, msg))
}

// "java/util/logging/Logger.logrb(Ljava/util/logging/Level;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V"
// The resource bundle name is accepted but ignored, since there is no resource bundle lookup support here.
func loggingLoggerLogrb(params []interface{}) interface{} {
	levelName := loggingLoggerLevelName(params[1])
	sourceClass := loggingLoggerStringArg(params[2])
	sourceMethod := loggingLoggerStringArg(params[3])
	msg := loggingLoggerStringArg(params[5])
	return loggingLoggerWrite(levelName, fmt.Sprintf("%s %s: %s", sourceClass, sourceMethod, msg))
}

// "java/util/logging/Logger.logrb(Ljava/util/logging/Level;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;[Ljava/lang/Object;)V"
func loggingLoggerLogrbWithParams(params []interface{}) interface{} {
	levelName := loggingLoggerLevelName(params[1])
	sourceClass := loggingLoggerStringArg(params[2])
	sourceMethod := loggingLoggerStringArg(params[3])
	msg := loggingLoggerStringArg(params[5])
	extra := loggingLoggerFormatParams(params[6])
	if extra != "" {
		msg = msg + " [" + extra + "]"
	}
	return loggingLoggerWrite(levelName, fmt.Sprintf("%s %s: %s", sourceClass, sourceMethod, msg))
}

// "java/util/logging/Logger.severe(Ljava/lang/String;)V"
func loggingLoggerSevere(params []interface{}) interface{} {
	return loggingLoggerWrite("SEVERE", loggingLoggerStringArg(params[1]))
}

// "java/util/logging/Logger.throwing(Ljava/lang/String;Ljava/lang/String;Ljava/lang/Throwable;)V"
func loggingLoggerThrowing(params []interface{}) interface{} {
	sourceClass := loggingLoggerStringArg(params[1])
	sourceMethod := loggingLoggerStringArg(params[2])
	return loggingLoggerWrite("FINER", fmt.Sprintf("THROW %s %s", sourceClass, sourceMethod))
}

// "java/util/logging/Logger.warning(Ljava/lang/String;)V"
func loggingLoggerWarning(params []interface{}) interface{} {
	return loggingLoggerWrite("WARNING", loggingLoggerStringArg(params[1]))
}

// "java/util/logging/Logger.warning(Ljava/lang/String;[Ljava/lang/Object;)V"
func loggingLoggerWarningWithParams(params []interface{}) interface{} {
	msg := loggingLoggerStringArg(params[1])
	extra := loggingLoggerFormatParams(params[2])
	if extra != "" {
		msg = msg + " [" + extra + "]"
	}
	return loggingLoggerWrite("WARNING", msg)
}
