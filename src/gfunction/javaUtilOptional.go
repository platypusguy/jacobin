/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package gfunction

import (
	"fmt"
	"jacobin/excNames"
	"jacobin/object"
	"jacobin/stringPool"
	"jacobin/types"
)

func Load_Util_Optional() {
	MethodSignatures["java/util/Optional.<clinit>()V"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  clinitGeneric,
		}

	MethodSignatures["java/util/Optional.empty()Ljava/util/Optional;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  optionalEmpty,
		}

	MethodSignatures["java/util/Optional.equals(Ljava/lang/Object;)Z"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  optionalEquals,
		}

	MethodSignatures["java/util/Optional.filter(Ljava/util/function/Predicate;)Ljava/util/Optional;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Optional.flatMap(Ljava/util/function/Function;)Ljava/util/Optional;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Optional.get()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  optionalGet,
		}

	MethodSignatures["java/util/Optional.hashCode()I"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Optional.ifPresent(Ljava/util/function/Consumer;)V"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Optional.ifPresentOrElse(Ljava/util/function/Consumer;Ljava/lang/Runnable;)V"] =
		GMeth{
			ParamSlots: 2,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Optional.isEmpty()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  optionalIsEmpty,
		}

	MethodSignatures["java/util/Optional.isPresent()Z"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  optionalIsPresent,
		}

	MethodSignatures["java/util/Optional.map(Ljava/util/function/Function;)Ljava/util/Optional;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Optional.of(Ljava/lang/Object;)Ljava/util/Optional;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Optional.ofNullable(Ljava/lang/Object;)Ljava/util/Optional;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Optional.or(Ljava/util/function/Supplier;)Ljava/util/Optional;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Optional.orElse(Ljava/lang/Object;)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  optionalOrElse,
		}

	MethodSignatures["java/util/Optional.orElseGet(Ljava/util/function/Supplier;)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Optional.orElseThrow()Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  optionalOrElseThrow,
		}

	MethodSignatures["java/util/Optional.orElseThrow(Ljava/util/function/Supplier;)Ljava/lang/Object;"] =
		GMeth{
			ParamSlots: 1,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Optional.stream()Ljava/util/Stream;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  trapFunction,
		}

	MethodSignatures["java/util/Optional.toString()Ljava/lang/String;"] =
		GMeth{
			ParamSlots: 0,
			GFunction:  optionalToString,
		}
}

var classNameOptional string = "java/util/Optional"

func optionalEmpty(params []interface{}) interface{} {
	return object.MakeEmptyObjectWithClassName(&classNameOptional)
}

func optionalEquals(params []interface{}) interface{} {
	var thisFvalue any
	var thatFvalue any

	// Get this object.
	this, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "optionalEquals: This is not an Optional object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get this field value.
	thisFvalue = this.FieldTable["value"].Fvalue

	// Get argument object.
	that, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "optionalEquals: Parameter is not an object"
		return getGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	if *stringPool.GetStringPointer(this.KlassName) != classNameOptional {
		return types.JavaBoolFalse
	}

	// Get that field value.
	thatFvalue = that.FieldTable["value"].Fvalue

	// Are they both nil?
	if (thisFvalue == nil) && (thatFvalue == nil) {
		return types.JavaBoolTrue // yes
	}

	// Are they equal?
	if thatFvalue != thisFvalue {
		return types.JavaBoolFalse // no
	}

	// They are equal.
	return types.JavaBoolTrue
}

func optionalGet(params []interface{}) interface{} {
	// Get this object.
	this, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "optionalGet: This is not an Optional object"
		return getGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// Get this field value.
	thisFvalue := this.FieldTable["value"].Fvalue
	if thisFvalue == nil {
		errMsg := "optionalGet: Value field not present"
		return getGErrBlk(excNames.NoSuchElementException, errMsg)
	}

	return thisFvalue
}

func optionalIsEmpty(params []interface{}) interface{} {
	// Get this object.
	this, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "optionalIsEmpty: This is not an Optional object"
		return getGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// Get this field value.
	thisFvalue := this.FieldTable["value"].Fvalue
	if thisFvalue == nil {
		return types.JavaBoolTrue
	}
	return types.JavaBoolFalse
}

func optionalIsPresent(params []interface{}) interface{} {
	// Get this object.
	this, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "optionalIsPresent: This is not an Optional object"
		return getGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// Get this field value.
	thisFvalue := this.FieldTable["value"].Fvalue
	if thisFvalue == nil {
		return types.JavaBoolFalse
	}
	return types.JavaBoolTrue
}

func optionalToString(params []interface{}) interface{} {
	// Get this object.
	this, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "optionalIsPresent: This is not an Optional object"
		return getGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// Get this field value.
	thisFvalue := this.FieldTable["value"].Fvalue
	if thisFvalue == nil {
		return object.StringObjectFromGoString("empty")
	}

	// Return stringified value.
	return object.StringObjectFromGoString(fmt.Sprintf("%T :: %v", thisFvalue, thisFvalue))

}

func optionalOrElse(params []interface{}) interface{} {
	// Get this object.
	this, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "optionalOrElse: This is not an object"
		return getGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// Get this object.
	that, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "optionalOrElse: Parameter is not an object"
		return getGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// Get the field values.
	thisFvalue := this.FieldTable["value"].Fvalue
	thatFvalue := that.FieldTable["value"].Fvalue

	// If this field value is present, return it.
	// Else return that object's field value.
	if thisFvalue == nil {
		return thatFvalue
	}
	return thisFvalue

}

func optionalOrElseThrow(params []interface{}) interface{} {
	// Get this object.
	this, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "optionalOrElseThrow: This is not an object"
		return getGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// If this field value is present, return it.
	// Else throw a NoSuchElementException.
	thisFvalue := this.FieldTable["value"].Fvalue
	if thisFvalue == nil {
		errMsg := "optionalOrElseThrow: Value field not present"
		return getGErrBlk(excNames.NoSuchElementException, errMsg)
	}
	return thisFvalue

}
