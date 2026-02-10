/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"jacobin/src/gfunction/ghelpers"
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"math/big"
	"testing"
)

func TestLoad_ECParameterSpec(t *testing.T) {
	globals.InitGlobals("test")
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_ECParameterSpec()

	expectedSignatures := []string{
		"java/security/spec/ECGenParameterSpec.<init>(Ljava/lang/String;)V",
		"java/security/spec/ECParameterSpec.<init>(Ljava/security/spec/EllipticCurve;Ljava/security/spec/ECPoint;Ljava/math/BigInteger;I)V",
		"java/security/spec/ECParameterSpec.getCurve()Ljava/security/spec/EllipticCurve;",
		"java/security/spec/ECParameterSpec.getGenerator()Ljava/security/spec/ECPoint;",
		"java/security/spec/ECParameterSpec.getOrder()Ljava/math/BigInteger;",
		"java/security/spec/ECParameterSpec.getCofactor()I",
	}

	for _, sig := range expectedSignatures {
		if _, ok := ghelpers.MethodSignatures[sig]; !ok {
			t.Errorf("Expected signature %s not registered", sig)
		}
	}
}

func TestECParameterGenSpecInitStringDetailed(t *testing.T) {
	globals.InitGlobals("test")

	tests := []struct {
		name      string
		curveName string
		wantErr   bool
	}{
		{"P-256", "secp256r1", false},
		{"P-224", "P-224", false},
		{"P-384", "secp384r1", false},
		{"P-521", "P-521", false},
		{"Invalid", "invalid-curve", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thisObj := object.MakeEmptyObjectWithClassName(&types.ClassNameECGenParameterSpec)
			curveNameObj := object.StringObjectFromGoString(tt.curveName)

			params := []any{thisObj, curveNameObj}
			res := ecparameterGenSpecInitString(params)

			if tt.wantErr {
				if _, ok := res.(*ghelpers.GErrBlk); !ok {
					t.Errorf("Expected error for curve %s, got %v", tt.curveName, res)
				}
			} else {
				if res != nil {
					t.Errorf("Expected nil return for curve %s, got %v", tt.curveName, res)
				}
				// Verify fields were populated
				if _, ok := thisObj.FieldTable["curve"]; !ok {
					t.Error("curve field not set")
				}
				if _, ok := thisObj.FieldTable["g"]; !ok {
					t.Error("g (generator) field not set")
				}
				if _, ok := thisObj.FieldTable["n"]; !ok {
					t.Error("n (order) field not set")
				}
				if _, ok := thisObj.FieldTable["h"]; !ok {
					t.Error("h (cofactor) field not set")
				}
			}
		})
	}
}

func TestECParameterSpecInit(t *testing.T) {
	globals.InitGlobals("test")

	thisObj := object.MakeEmptyObjectWithClassName(&types.ClassNameECParameterSpec)
	curveObj := object.MakeEmptyObjectWithClassName(&types.ClassNameEllipticCurve)
	gObj := object.MakeEmptyObjectWithClassName(&types.ClassNameECPoint)
	nObj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, big.NewInt(12345))
	hVal := int64(1)

	params := []any{thisObj, curveObj, gObj, nObj, hVal}
	res := ecparameterSpecInit(params)

	if res != nil {
		t.Errorf("Expected nil return, got %v", res)
	}

	if val, ok := thisObj.FieldTable["curve"]; !ok || val.Fvalue != curveObj {
		t.Error("curve field not set correctly")
	}
	if val, ok := thisObj.FieldTable["g"]; !ok || val.Fvalue != gObj {
		t.Error("g field not set correctly")
	}
	if val, ok := thisObj.FieldTable["n"]; !ok || val.Fvalue != nObj {
		t.Error("n field not set correctly")
	}
	if val, ok := thisObj.FieldTable["h"]; !ok || val.Fvalue != hVal {
		t.Error("h field not set correctly")
	}
}

func TestECParameterSpecGetters(t *testing.T) {
	globals.InitGlobals("test")

	thisObj := object.MakeEmptyObjectWithClassName(&types.ClassNameECParameterSpec)
	curveObj := object.MakeEmptyObjectWithClassName(&types.ClassNameEllipticCurve)
	gObj := object.MakeEmptyObjectWithClassName(&types.ClassNameECPoint)
	nObj := object.MakePrimitiveObject(types.ClassNameBigInteger, types.BigInteger, big.NewInt(12345))
	hVal := int64(1)

	thisObj.FieldTable = map[string]object.Field{
		"curve": {Ftype: types.Ref, Fvalue: curveObj},
		"g":     {Ftype: types.Ref, Fvalue: gObj},
		"n":     {Ftype: types.Ref, Fvalue: nObj},
		"h":     {Ftype: types.Int, Fvalue: hVal},
	}

	// Test getCurve
	res := ecparameterSpecGetCurve([]any{thisObj})
	if res != curveObj {
		t.Errorf("getCurve: expected %v, got %v", curveObj, res)
	}

	// Test getGenerator
	res = ecparameterSpecGetGenerator([]any{thisObj})
	if res != gObj {
		t.Errorf("getGenerator: expected %v, got %v", gObj, res)
	}

	// Test getOrder
	res = ecparameterSpecGetOrder([]any{thisObj})
	if res != nObj {
		t.Errorf("getOrder: expected %v, got %v", nObj, res)
	}

	// Test getCofactor
	res = ecparameterSpecGetCofactor([]any{thisObj})
	if res != hVal {
		t.Errorf("getCofactor: expected %v, got %v", hVal, res)
	}
}

func TestECParameterSpecInvalidParams(t *testing.T) {
	globals.InitGlobals("test")

	t.Run("ecparameterGenSpecInitString invalid count", func(t *testing.T) {
		res := ecparameterGenSpecInitString([]any{})
		if _, ok := res.(*ghelpers.GErrBlk); !ok {
			t.Error("Expected error for empty params")
		}
	})

	t.Run("ecparameterSpecInit invalid count", func(t *testing.T) {
		res := ecparameterSpecInit([]any{})
		if _, ok := res.(*ghelpers.GErrBlk); !ok {
			t.Error("Expected error for empty params")
		}
	})

	t.Run("Getters invalid count", func(t *testing.T) {
		if _, ok := ecparameterSpecGetCurve([]any{}).(*ghelpers.GErrBlk); !ok {
			t.Error("getCurve should fail with no params")
		}
		if _, ok := ecparameterSpecGetGenerator([]any{}).(*ghelpers.GErrBlk); !ok {
			t.Error("getGenerator should fail with no params")
		}
		if _, ok := ecparameterSpecGetOrder([]any{}).(*ghelpers.GErrBlk); !ok {
			t.Error("getOrder should fail with no params")
		}
		if _, ok := ecparameterSpecGetCofactor([]any{}).(*ghelpers.GErrBlk); !ok {
			t.Error("getCofactor should fail with no params")
		}
	})
}
