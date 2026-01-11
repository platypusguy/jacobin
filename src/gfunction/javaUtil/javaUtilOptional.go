/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaUtil

import (
	"fmt"
	"jacobin/src/excNames"
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/object"
	"jacobin/src/stringPool"
	"jacobin/src/types"
)

func Load_Util_Optional() {
	ghelpers.MethodSignatures["java/util/Optional.<clinit>()V"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.ClinitGeneric,
		}

	ghelpers.MethodSignatures["java/util/Optional.empty()Ljava/util/Optional;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  optionalEmpty,
		}

	ghelpers.MethodSignatures["java/util/Optional.equals(Ljava/lang/Object;)Z"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  optionalEquals,
		}

	ghelpers.MethodSignatures["java/util/Optional.filter(Ljava/util/function/Predicate;)Ljava/util/Optional;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Optional.flatMap(Ljava/util/function/Function;)Ljava/util/Optional;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Optional.get()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  optionalGet,
		}

	ghelpers.MethodSignatures["java/util/Optional.hashCode()I"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Optional.ifPresent(Ljava/util/function/Consumer;)V"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Optional.ifPresentOrElse(Ljava/util/function/Consumer;Ljava/lang/Runnable;)V"] =
		ghelpers.GMeth{
			ParamSlots: 2,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Optional.isEmpty()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  optionalIsEmpty,
		}

	ghelpers.MethodSignatures["java/util/Optional.isPresent()Z"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  optionalIsPresent,
		}

	ghelpers.MethodSignatures["java/util/Optional.map(Ljava/util/function/Function;)Ljava/util/Optional;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Optional.of(Ljava/lang/Object;)Ljava/util/Optional;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Optional.ofNullable(Ljava/lang/Object;)Ljava/util/Optional;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Optional.or(Ljava/util/function/Supplier;)Ljava/util/Optional;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Optional.orElse(Ljava/lang/Object;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  optionalOrElse,
		}

	ghelpers.MethodSignatures["java/util/Optional.orElseGet(Ljava/util/function/Supplier;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Optional.orElseThrow()Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  optionalOrElseThrow,
		}

	ghelpers.MethodSignatures["java/util/Optional.orElseThrow(Ljava/util/function/Supplier;)Ljava/lang/Object;"] =
		ghelpers.GMeth{
			ParamSlots: 1,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Optional.stream()Ljava/util/Stream;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  ghelpers.TrapFunction,
		}

	ghelpers.MethodSignatures["java/util/Optional.toString()Ljava/lang/String;"] =
		ghelpers.GMeth{
			ParamSlots: 0,
			GFunction:  optionalToString,
		}
}

func optionalEmpty(params []interface{}) interface{} {
	return object.MakeEmptyObjectWithClassName(&types.ClassNameOptional)
}

func optionalEquals(params []interface{}) interface{} {
	var thisFvalue any
	var thatFvalue any

	// Get this object.
	this, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "optionalEquals: This is not an Optional object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}

	// Get this field value.
	thisFvalue = this.FieldTable["value"].Fvalue

	// Get argument object.
	that, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "optionalEquals: Parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.IllegalArgumentException, errMsg)
	}
	if *stringPool.GetStringPointer(this.KlassName) != types.ClassNameOptional {
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
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// Get this field value.
	thisFvalue := this.FieldTable["value"].Fvalue
	if thisFvalue == nil {
		errMsg := "optionalGet: Value field not present"
		return ghelpers.GetGErrBlk(excNames.NoSuchElementException, errMsg)
	}

	return thisFvalue
}

func optionalIsEmpty(params []interface{}) interface{} {
	// Get this object.
	this, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "optionalIsEmpty: This is not an Optional object"
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
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
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
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
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// Get this field value.
	thisFvalue := this.FieldTable["value"].Fvalue
	if thisFvalue == nil {
		return object.StringObjectFromGoString("empty")
	}

	// Return stringified value.
	return object.StringObjectFromGoString(fmt.Sprintf("Optional[%T :: %v]", thisFvalue, thisFvalue))

}

func optionalOrElse(params []interface{}) interface{} {
	// Get this object.
	this, ok := params[0].(*object.Object)
	if !ok {
		errMsg := "optionalOrElse: This is not an object"
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// Get this object.
	that, ok := params[1].(*object.Object)
	if !ok {
		errMsg := "optionalOrElse: Parameter is not an object"
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
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
		return ghelpers.GetGErrBlk(excNames.VirtualMachineError, errMsg)
	}

	// If this field value is present, return it.
	// Else throw a NoSuchElementException.
	thisFvalue := this.FieldTable["value"].Fvalue
	if thisFvalue == nil {
		errMsg := "optionalOrElseThrow: Value field not present"
		return ghelpers.GetGErrBlk(excNames.NoSuchElementException, errMsg)
	}
	return thisFvalue

}
