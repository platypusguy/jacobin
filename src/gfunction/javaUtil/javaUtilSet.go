/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
)

func Load_Util_Set() {

	ghelpers.MethodSignatures["java/util/Set.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/Set.add(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  hashsetAdd,
		}

	ghelpers.MethodSignatures["java/util/Set.addAll(Ljava/util/Collection;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Set.clear()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  hashmapInit,
		}

	ghelpers.MethodSignatures["java/util/Set.contains(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  hashsetContains,
		}

	ghelpers.MethodSignatures["java/util/Set.containsAll(Ljava/util/Collection;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Set.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Set.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Set.isEmpty()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  hashsetIsEmpty,
		}

	ghelpers.MethodSignatures["java/util/Set.iterator()Ljava/util/Iterator;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  setIterator,
		}

	ghelpers.MethodSignatures["java/util/Set.remove(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  hashsetRemove,
		}

	ghelpers.MethodSignatures["java/util/Set.removeAll(Ljava/util/Collection;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Set.retainAll(Ljava/util/Collection;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Set.size()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  hashmapSize,
		}

	ghelpers.MethodSignatures["java/util/Set.spliterator()Ljava/util/Spliterator;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Set.toArray()[Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  hashsetToArray,
		}

	ghelpers.MethodSignatures["java/util/Set.toArray([Ljava/lang/Object;)[Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Set.of()Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  setOf,
		}

	ghelpers.MethodSignatures["java/util/Set.of(Ljava/lang/Object;)Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  setOf,
		}

	ghelpers.MethodSignatures["java/util/Set.of(Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  setOf,
		}

	ghelpers.MethodSignatures["java/util/Set.of(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 3,
			GFunction:  setOf,
		}

	ghelpers.MethodSignatures["java/util/Set.of(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 4,
			GFunction:  setOf,
		}

	ghelpers.MethodSignatures["java/util/Set.of(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 5,
			GFunction:  setOf,
		}

	ghelpers.MethodSignatures["java/util/Set.of(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 6,
			GFunction:  setOf,
		}

	ghelpers.MethodSignatures["java/util/Set.of(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 7,
			GFunction:  setOf,
		}

	ghelpers.MethodSignatures["java/util/Set.of(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 8,
			GFunction:  setOf,
		}

	ghelpers.MethodSignatures["java/util/Set.of(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 9,
			GFunction:  setOf,
		}

	ghelpers.MethodSignatures["java/util/Set.of(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 10,
			GFunction:  setOf,
		}

	ghelpers.MethodSignatures["java/util/Set.of([Ljava/lang/Object;)Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  setOfVarargs,
		}

	ghelpers.MethodSignatures["java/util/Set.copyOf(Ljava/util/Collection;)Ljava/util/Set;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}
}

func setIterator(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "setIterator: Invalid self argument")
	}
	return NewIterator(self)
}

func setOf(params []interface{}) interface{} {
	// Java Set.of(...) returns an unmodifiable set.
	// In Java, Set.of(...) also forbids null elements and duplicates.
	seen := make(map[interface{}]bool)
	for _, p := range params {
		if p == nil || p == object.Null {
			return ghelpers.GetGErrBlk(excNames.NullPointerException, "Set.of: null element")
		}
		// In a real VM we should check for value equality.
		// For now, pointer equality for objects.
		if seen[p] {
			return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, "Set.of: duplicate element")
		}
		seen[p] = true
	}

	// We'll return a HashSet for now as a minimal implementation.
	hs := object.MakeEmptyObjectWithClassName(&classNameHashSet)
	if ret := hashmapInit([]interface{}{hs}); ret != nil {
		return ret
	}

	for _, p := range params {
		if ret := hashsetAdd([]interface{}{hs, p}); ret != nil {
			if _, ok := ret.(*ghelpers.GErrBlk); ok {
				return ret
			}
		}
	}

	return hs
}

func setOfVarargs(params []interface{}) interface{} {
	if len(params) == 0 {
		return setOf([]interface{}{})
	}

	arrayObj, ok := params[0].(*object.Object)
	if !ok || arrayObj == nil {
		return ghelpers.GetGErrBlk(excNames.NullPointerException, "Set.of: array is null")
	}

	field, ok := arrayObj.FieldTable["value"]
	if !ok {
		return setOf([]interface{}{})
	}

	elements, ok := field.Fvalue.([]*object.Object)
	if !ok {
		elements2, ok2 := field.Fvalue.([]interface{})
		if !ok2 {
			return setOf([]interface{}{})
		}
		return setOf(elements2)
	}

	ifaceElements := make([]interface{}, len(elements))
	for i, e := range elements {
		ifaceElements[i] = e
	}

	return setOf(ifaceElements)
}
