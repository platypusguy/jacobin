/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

import (
	"math/big"
	"testing"

	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
)

func makeBigIntegerObject(val int64) *object.Object {
	obj := object.MakeEmptyObjectWithClassName(&types.ClassNameBigInteger)
	ghelpers.InitBigIntegerField(obj, val)
	return obj
}

func TestDHParameterSpecInit(t *testing.T) {
	globals.InitGlobals("test")

	t.Run("Success 2-args", func(t *testing.T) {
		pObj := makeBigIntegerObject(23)
		gObj := makeBigIntegerObject(5)
		specObj := object.MakeEmptyObjectWithClassName(&types.ClassNameDHParameterSpec)

		args := []any{specObj, pObj, gObj}
		result := dhparameterspecInit(args)

		if result != nil {
			t.Fatalf("Expected nil result, got %v", result)
		}

		p := specObj.FieldTable["p"].Fvalue.(*big.Int)
		if p.Int64() != 23 {
			t.Errorf("Expected p=23, got %v", p.Int64())
		}
		g := specObj.FieldTable["g"].Fvalue.(*big.Int)
		if g.Int64() != 5 {
			t.Errorf("Expected g=5, got %v", g.Int64())
		}
		l := specObj.FieldTable["l"].Fvalue.(int64)
		if l != 0 {
			t.Errorf("Expected l=0, got %v", l)
		}
	})

	t.Run("Success 3-args", func(t *testing.T) {
		pObj := makeBigIntegerObject(23)
		gObj := makeBigIntegerObject(5)
		specObj := object.MakeEmptyObjectWithClassName(&types.ClassNameDHParameterSpec)

		args := []any{specObj, pObj, gObj, int64(10)}
		result := dhparameterspecInit(args)

		if result != nil {
			t.Fatalf("Expected nil result, got %v", result)
		}

		l := specObj.FieldTable["l"].Fvalue.(int64)
		if l != 10 {
			t.Errorf("Expected l=10, got %v", l)
		}
	})

	t.Run("Wrong number of arguments", func(t *testing.T) {
		args := []any{object.MakeEmptyObject()}
		result := dhparameterspecInit(args)
		if _, ok := result.(*ghelpers.GErrBlk); !ok {
			t.Errorf("Expected GErrBlk, got %T", result)
		}
	})

	t.Run("Receiver not an object", func(t *testing.T) {
		args := []any{"not an object", makeBigIntegerObject(23), makeBigIntegerObject(5)}
		result := dhparameterspecInit(args)
		if _, ok := result.(*ghelpers.GErrBlk); !ok {
			t.Errorf("Expected GErrBlk, got %T", result)
		}
	})

	t.Run("p not BigInteger object", func(t *testing.T) {
		specObj := object.MakeEmptyObject()
		args := []any{specObj, "not an object", makeBigIntegerObject(5)}
		result := dhparameterspecInit(args)
		err, ok := result.(*ghelpers.GErrBlk)
		if !ok || err.ErrMsg != "dhparameterspecInit: p must be a BigInteger object" {
			t.Errorf("Expected specific error, got %v", result)
		}
	})

	t.Run("g not BigInteger object", func(t *testing.T) {
		specObj := object.MakeEmptyObject()
		args := []any{specObj, makeBigIntegerObject(23), "not an object"}
		result := dhparameterspecInit(args)
		err, ok := result.(*ghelpers.GErrBlk)
		if !ok || err.ErrMsg != "dhparameterspecInit: g must be a BigInteger object" {
			t.Errorf("Expected specific error, got %v", result)
		}
	})

	t.Run("l not int64", func(t *testing.T) {
		specObj := object.MakeEmptyObject()
		args := []any{specObj, makeBigIntegerObject(23), makeBigIntegerObject(5), "not int64"}
		result := dhparameterspecInit(args)
		err, ok := result.(*ghelpers.GErrBlk)
		if !ok || err.ErrMsg != "dhparameterspecInit: l must be int64" {
			t.Errorf("Expected specific error, got %v", result)
		}
	})
}

func TestDHParameterSpecGetters(t *testing.T) {
	globals.InitGlobals("test")

	pVal := big.NewInt(23)
	gVal := big.NewInt(5)
	lVal := int64(10)

	specObj := object.MakeEmptyObject()
	specObj.FieldTable["p"] = object.Field{Ftype: types.BigInteger, Fvalue: pVal}
	specObj.FieldTable["g"] = object.Field{Ftype: types.BigInteger, Fvalue: gVal}
	specObj.FieldTable["l"] = object.Field{Ftype: types.Int, Fvalue: lVal}

	t.Run("GetP Success", func(t *testing.T) {
		result := dhparameterspecGetP([]any{specObj})
		obj, ok := result.(*object.Object)
		if !ok {
			t.Fatalf("Expected *object.Object, got %T", result)
		}
		val := obj.FieldTable["value"].Fvalue.(*big.Int)
		if val.Cmp(pVal) != 0 {
			t.Errorf("Expected p=%v, got %v", pVal, val)
		}
	})

	t.Run("GetG Success", func(t *testing.T) {
		result := dhparameterspecGetG([]any{specObj})
		obj, ok := result.(*object.Object)
		if !ok {
			t.Fatalf("Expected *object.Object, got %T", result)
		}
		val := obj.FieldTable["value"].Fvalue.(*big.Int)
		if val.Cmp(gVal) != 0 {
			t.Errorf("Expected g=%v, got %v", gVal, val)
		}
	})

	t.Run("GetL Success", func(t *testing.T) {
		result := dhparameterspecGetL([]any{specObj})
		val, ok := result.(int64)
		if !ok {
			t.Fatalf("Expected int64, got %T", result)
		}
		if val != lVal {
			t.Errorf("Expected l=%v, got %v", lVal, val)
		}
	})

	t.Run("GetP Error Not Initialized", func(t *testing.T) {
		emptyObj := object.MakeEmptyObject()
		result := dhparameterspecGetP([]any{emptyObj})
		if _, ok := result.(*ghelpers.GErrBlk); !ok {
			t.Errorf("Expected GErrBlk, got %T", result)
		}
	})

	t.Run("GetG Error Not Initialized", func(t *testing.T) {
		emptyObj := object.MakeEmptyObject()
		result := dhparameterspecGetG([]any{emptyObj})
		if _, ok := result.(*ghelpers.GErrBlk); !ok {
			t.Errorf("Expected GErrBlk, got %T", result)
		}
	})

	t.Run("GetL Default Behavior", func(t *testing.T) {
		emptyObj := object.MakeEmptyObject()
		result := dhparameterspecGetL([]any{emptyObj})
		val, ok := result.(int64)
		if !ok || val != 0 {
			t.Errorf("Expected int64(0), got %v (%T)", result, result)
		}
	})
}
