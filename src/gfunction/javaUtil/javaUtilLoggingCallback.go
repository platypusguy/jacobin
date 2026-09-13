/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"container/list"
	"jacobin/src/classloader"
	"jacobin/src/frames"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
)

// This file provides support for dispatching to user-supplied (bytecode)
// subclasses of java.util.logging.Formatter and java.util.logging.Filter,
// so that e.g. a custom Formatter.format(LogRecord) or Filter.isLoggable(LogRecord)
// actually gets executed, rather than always falling back to jacobin's
// built-in default behavior. This mirrors the real JDK, where Handler.publish()
// always calls through getFormatter().format(record) and, if a filter is set,
// filter.isLoggable(record) -- both of which are ordinary (possibly overridden)
// Java method calls.

// loggingExtractFsAndArgs inspects the leading element of params: if it is a
// frame stack (*list.List), it is treated as call context prepended by the
// caller (see loggingLoggerPublishToHandler) and is stripped off; otherwise
// there is no available frame stack (e.g. when called directly as a GFunction
// from bytecode) and nil is returned for it.
func loggingExtractFsAndArgs(params []interface{}) (*list.List, []interface{}) {
	if len(params) > 0 {
		if fs, ok := params[0].(*list.List); ok {
			return fs, params[1:]
		}
	}
	return nil, params
}

// loggingInvokeUserJavaMethod attempts to dynamically dispatch methodName/methodType
// on obj's concrete class (and its ancestors), but only if the found implementation
// is real Java bytecode (MType == 'J'); i.e., user-supplied code, as opposed to a
// jacobin-native (GFunction) implementation, which callers already know how to
// handle themselves. Returns (result, true) on success, or (nil, false) if no
// frame stack is available, obj is null, or no Java bytecode implementation is found.
func loggingInvokeUserJavaMethod(fs *list.List, obj *object.Object, methodName, methodType string, args ...interface{}) (interface{}, bool) {
	if fs == nil || obj == nil || object.IsNull(obj) {
		return nil, false
	}

	currClass := object.GoStringFromStringPoolIndex(obj.KlassName)
	for currClass != "" {
		specificFQN := currClass + "." + methodName + methodType

		mtEntry := classloader.GetMtableEntry(specificFQN)
		if mtEntry.Meth == nil {
			mtEntry, _ = classloader.FetchMethodAndCP(currClass, methodName, methodType)
		}

		if mtEntry.Meth != nil {
			if mtEntry.MType != 'J' {
				// A jacobin-native (GFunction) implementation was found -- this is
				// not user-supplied code, so let the caller apply its own fallback.
				return nil, false
			}

			callArgs := make([]interface{}, 0, len(args)+1)
			callArgs = append(callArgs, obj)
			callArgs = append(callArgs, args...)
			globals.GetGlobalRef().FuncRunJavaFromG(fs, currClass, methodName, methodType, callArgs...)

			fr, ok := fs.Front().Value.(*frames.Frame)
			if !ok || fr.TOS < 0 {
				return nil, true
			}
			res := fr.OpStack[fr.TOS]
			fr.TOS--
			return res, true
		}

		klass := classloader.MethAreaFetch(currClass)
		if klass == nil || klass.Data.SuperclassIndex == types.InvalidStringIndex {
			break
		}
		currClass = *stringPool.GetStringPointer(uint32(klass.Data.SuperclassIndex))
	}
	return nil, false
}
