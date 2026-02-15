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

func TestDHPrivateGetX(t *testing.T) {
	globals.InitGlobals("test")

	t.Run("Success", func(t *testing.T) {
		xValue := big.NewInt(12345)
		innerObj := object.MakeEmptyObject()
		innerObj.FieldTable["x"] = object.Field{Ftype: types.BigInteger, Fvalue: xValue}

		thisObj := object.MakeEmptyObject()
		thisObj.FieldTable["value"] = object.Field{Ftype: types.Ref, Fvalue: innerObj}

		result := dhPrivateGetX([]any{thisObj})
		resObj, ok := result.(*object.Object)
		if !ok {
			t.Fatalf("Expected *object.Object, got %T", result)
		}

		if resObj.FieldTable["value"].Fvalue.(*big.Int).Cmp(xValue) != 0 {
			t.Errorf("Expected %v, got %v", xValue, resObj.FieldTable["value"].Fvalue)
		}
	})

	t.Run("WrongNumberOfParams", func(t *testing.T) {
		result := dhPrivateGetX([]any{})
		if _, ok := result.(*ghelpers.GErrBlk); !ok {
			t.Errorf("Expected *ghelpers.GErrBlk, got %T", result)
		}
	})

	t.Run("NotAnObject", func(t *testing.T) {
		result := dhPrivateGetX([]any{"not an object"})
		if _, ok := result.(*ghelpers.GErrBlk); !ok {
			t.Errorf("Expected *ghelpers.GErrBlk, got %T", result)
		}
	})

	t.Run("ExtractionFailed", func(t *testing.T) {
		thisObj := object.MakeEmptyObject()
		thisObj.FieldTable["value"] = object.Field{Ftype: types.Int, Fvalue: int64(123)} // Wrong type

		result := dhPrivateGetX([]any{thisObj})
		gerr, ok := result.(*ghelpers.GErrBlk)
		if !ok {
			t.Fatalf("Expected *ghelpers.GErrBlk, got %T", result)
		}
		if gerr.ErrMsg != "dhPrivateGetX: DH private key extraction failed" {
			t.Errorf("Unexpected error message: %s", gerr.ErrMsg)
		}
	})

	t.Run("FieldExtractionFailed", func(t *testing.T) {
		innerObj := object.MakeEmptyObject()
		// "x" field is missing or wrong type
		innerObj.FieldTable["x"] = object.Field{Ftype: types.Int, Fvalue: int64(123)}

		thisObj := object.MakeEmptyObject()
		thisObj.FieldTable["value"] = object.Field{Ftype: types.Ref, Fvalue: innerObj}

		result := dhPrivateGetX([]any{thisObj})
		gerr, ok := result.(*ghelpers.GErrBlk)
		if !ok {
			t.Fatalf("Expected *ghelpers.GErrBlk, got %T", result)
		}
		if gerr.ErrMsg != "dhPrivateKeyGetX: DH private key x-field extraction failed" {
			t.Errorf("Unexpected error message: %s", gerr.ErrMsg)
		}
	})
}

func TestDHPublicGetY(t *testing.T) {
	globals.InitGlobals("test")

	t.Run("Success", func(t *testing.T) {
		yValue := big.NewInt(67890)
		innerObj := object.MakeEmptyObject()
		innerObj.FieldTable["y"] = object.Field{Ftype: types.BigInteger, Fvalue: yValue}

		thisObj := object.MakeEmptyObject()
		thisObj.FieldTable["value"] = object.Field{Ftype: types.Ref, Fvalue: innerObj}

		result := dhPublicKeyGetY([]any{thisObj})
		resObj, ok := result.(*object.Object)
		if !ok {
			t.Fatalf("Expected *object.Object, got %T", result)
		}

		if resObj.FieldTable["value"].Fvalue.(*big.Int).Cmp(yValue) != 0 {
			t.Errorf("Expected %v, got %v", yValue, resObj.FieldTable["value"].Fvalue)
		}
	})

	t.Run("ExtractionFailed", func(t *testing.T) {
		thisObj := object.MakeEmptyObject()
		// value is missing or not an object
		result := dhPublicKeyGetY([]any{thisObj})
		gerr, ok := result.(*ghelpers.GErrBlk)
		if !ok {
			t.Fatalf("Expected *ghelpers.GErrBlk, got %T", result)
		}
		if gerr.ErrMsg != "dhPublicKeyGetY: DH public key extraction failed" {
			t.Errorf("Unexpected error message: %s", gerr.ErrMsg)
		}
	})

	t.Run("FieldExtractionFailed", func(t *testing.T) {
		innerObj := object.MakeEmptyObject()
		// "y" field is missing or wrong type
		innerObj.FieldTable["y"] = object.Field{Ftype: types.Int, Fvalue: int64(123)}

		thisObj := object.MakeEmptyObject()
		thisObj.FieldTable["value"] = object.Field{Ftype: types.Ref, Fvalue: innerObj}

		result := dhPublicKeyGetY([]any{thisObj})
		gerr, ok := result.(*ghelpers.GErrBlk)
		if !ok {
			t.Fatalf("Expected *ghelpers.GErrBlk, got %T", result)
		}
		if gerr.ErrMsg != "dhPublicKeyGetY: DH public key y-field extraction failed" {
			t.Errorf("Unexpected error message: %s", gerr.ErrMsg)
		}
	})
}
