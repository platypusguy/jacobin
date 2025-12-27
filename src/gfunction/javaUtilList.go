/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2025 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"jacobin/src/excNames"
	"jacobin/src/object"
	"jacobin/src/types"
)

func Load_Util_List() {

	MethodSignatures["java/util/List.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/util/List.of()Ljava/util/List;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  listOf,
		}

	MethodSignatures["java/util/List.of(Ljava/lang/Object;)Ljava/util/List;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  listOf,
		}

	MethodSignatures["java/util/List.of(Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/List;"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  listOf,
		}

	MethodSignatures["java/util/List.of(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/List;"] =
		GMeth{
			ParamSlots: 3,
			GFunction:  listOf,
		}

	MethodSignatures["java/util/List.of(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/List;"] =
		GMeth{
			ParamSlots: 4,
			GFunction:  listOf,
		}

	MethodSignatures["java/util/List.of(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/List;"] =
		GMeth{
			ParamSlots: 5,
			GFunction:  listOf,
		}

	MethodSignatures["java/util/List.of(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/List;"] =
		GMeth{
			ParamSlots: 6,
			GFunction:  listOf,
		}

	MethodSignatures["java/util/List.of(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/List;"] =
		GMeth{
			ParamSlots: 7,
			GFunction:  listOf,
		}

	MethodSignatures["java/util/List.of(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/List;"] =
		GMeth{
			ParamSlots: 8,
			GFunction:  listOf,
		}

	MethodSignatures["java/util/List.of(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/List;"] =
		GMeth{
			ParamSlots: 9,
			GFunction:  listOf,
		}

	MethodSignatures["java/util/List.of(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Ljava/util/List;"] =
		GMeth{
			ParamSlots: 10,
			GFunction:  listOf,
		}

	MethodSignatures["java/util/List.of([Ljava/lang/Object;)Ljava/util/List;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  listOfVarargs,
		}

	// Traps for functions that reference forbidden types: Collection, Consumer, ListIterator, Spliterator, UnaryOperator, Comparator.

	MethodSignatures["java/util/List.addAll(Ljava/util/Collection;)Z"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/List.addAll(ILjava/util/Collection;)Z"] = GMeth{ParamSlots: 2, GFunction: trapFunction}
	MethodSignatures["java/util/List.containsAll(Ljava/util/Collection;)Z"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/List.removeAll(Ljava/util/Collection;)Z"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/List.retainAll(Ljava/util/Collection;)Z"] = GMeth{ParamSlots: 1, GFunction: trapFunction}
	MethodSignatures["java/util/List.copyOf(Ljava/util/Collection;)Ljava/util/List;"] = GMeth{ParamSlots: 1, GFunction: trapFunction}

	MethodSignatures["java/util/List.forEach(Ljava/util/function/Consumer;)V"] = GMeth{ParamSlots: 1, GFunction: trapFunction}

	MethodSignatures["java/util/List.listIterator()Ljava/util/ListIterator;"] = GMeth{
		ParamSlots: 0,
		GFunction:  listListIterator,
	}
	MethodSignatures["java/util/List.listIterator(I)Ljava/util/ListIterator;"] = GMeth{
		ParamSlots: 1,
		GFunction:  listListIteratorWithIndex,
	}

	MethodSignatures["java/util/List.spliterator()Ljava/util/Spliterator;"] = GMeth{ParamSlots: 0, GFunction: trapFunction}

	MethodSignatures["java/util/List.replaceAll(Ljava/util/function/UnaryOperator;)V"] = GMeth{ParamSlots: 1, GFunction: trapFunction}

	MethodSignatures["java/util/List.iterator()Ljava/util/Iterator;"] = GMeth{
		ParamSlots: 0,
		GFunction:  listIterator,
	}

	MethodSignatures["java/util/List.sort(Ljava/util/Comparator;)V"] = GMeth{ParamSlots: 1, GFunction: trapFunction}

}

func listIterator(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "listIterator: Invalid self argument")
	}
	return NewIterator(self)
}

func listListIterator(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "listIterator: Invalid self argument")
	}
	return NewListIterator(self, 0)
}

func listListIteratorWithIndex(params []interface{}) interface{} {
	self, ok := params[0].(*object.Object)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "listIterator: Invalid self argument")
	}
	index, ok := params[1].(int64)
	if !ok {
		return getGErrBlk(excNames.IllegalArgumentException, "listIterator: Invalid index argument")
	}
	return NewListIterator(self, int(index))
}

func listOf(params []interface{}) interface{} {
	// Java List.of(...) returns an unmodifiable list.
	// In Java, List.of(...) also forbids null elements.
	for _, p := range params {
		if p == nil || p == object.Null {
			return getGErrBlk(excNames.NullPointerException, "List.of: null element")
		}
	}

	// We'll return an ArrayList for now as a minimal implementation.
	// Copy params to a new slice to ensure it's independent.
	list := make([]interface{}, len(params))
	copy(list, params)

	listObj := object.MakePrimitiveObject("java/util/ArrayList", types.ArrayList, list)
	return listObj
}

func listOfVarargs(params []interface{}) interface{} {
	if len(params) == 0 {
		return listOf([]interface{}{})
	}

	arrayObj, ok := params[0].(*object.Object)
	if !ok || arrayObj == nil {
		return getGErrBlk(excNames.NullPointerException, "List.of: array is null")
	}

	field, ok := arrayObj.FieldTable["value"]
	if !ok {
		return listOf([]interface{}{})
	}

	// Reference arrays store elements as []*object.Object
	elements, ok := field.Fvalue.([]*object.Object)
	if !ok {
		// Try []interface{} just in case
		elements2, ok2 := field.Fvalue.([]interface{})
		if !ok2 {
			return listOf([]interface{}{})
		}
		return listOf(elements2)
	}

	// Convert []*object.Object to []interface{}
	ifaceElements := make([]interface{}, len(elements))
	for i, e := range elements {
		ifaceElements[i] = e
	}

	return listOf(ifaceElements)
}
