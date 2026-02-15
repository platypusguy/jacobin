/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"crypto/rand"
	"errors"
	"math/big"

	"jacobin/src/object"
	"jacobin/src/types"
)

// Generate DH key pair
func generateDHKeyPair(kpgObj *object.Object) (*object.Object, *object.Object, error) {

	// Get p.
	pField, ok := kpgObj.FieldTable["p"]
	if !ok {
		return nil, nil, errors.New("generateDHKeyPair: missing field p")
	}
	pBI, ok := pField.Fvalue.(*big.Int)
	if !ok {
		return nil, nil, errors.New("generateDHKeyPair: p value is not *big.Int")
	}

	// Get g.
	gField, ok := kpgObj.FieldTable["g"]
	if !ok {
		return nil, nil, errors.New("generateDHKeyPair: missing field g")
	}
	gBI, ok := gField.Fvalue.(*big.Int)
	if !ok {
		return nil, nil, errors.New("generateDHKeyPair: g value is not *big.Int")
	}

	// Get l.
	lField, ok := kpgObj.FieldTable["l"]
	if !ok {
		return nil, nil, errors.New("generateDHKeyPair: missing field l")
	}
	lValue, ok := lField.Fvalue.(int64)
	if !ok {
		return nil, nil, errors.New("generateDHKeyPair: l is not int64")
	}

	// Compute private exponent bit length.
	var xValue *big.Int
	var err error
	if lValue > 0 {
		xValue, err = randomBigInt(int(lValue))
	} else {
		xValue, err = randomBigInt(pBI.BitLen() - 1)
	}
	if err != nil {
		return nil, nil, err
	}

	// Compute public value y = g^x mod p.
	yValue := new(big.Int).Exp(gBI, xValue, pBI)

	// --- Construct DHPrivateKey object ---
	privKey := object.MakeEmptyObjectWithClassName(&types.ClassNameDHPrivateKey)
	privKey.FieldTable["x"] = object.Field{Ftype: types.BigInteger, Fvalue: xValue}
	privKey.FieldTable["p"] = object.Field{Ftype: types.BigInteger, Fvalue: pBI}
	privKey.FieldTable["g"] = object.Field{Ftype: types.BigInteger, Fvalue: gBI}
	privKey.FieldTable["l"] = object.Field{Ftype: types.Int, Fvalue: lValue}
	privKey.FieldTable["params"] = object.Field{Ftype: types.Ref, Fvalue: kpgObj.FieldTable["paramSpec"].Fvalue}
	privKey.FieldTable["algorithm"] = object.Field{Ftype: types.Ref, Fvalue: kpgObj.FieldTable["algorithm"].Fvalue}

	// --- Construct DHPublicKey object ---
	pubKey := object.MakeEmptyObjectWithClassName(&types.ClassNameDHPublicKey)
	pubKey.FieldTable["y"] = object.Field{Ftype: types.BigInteger, Fvalue: yValue}
	pubKey.FieldTable["p"] = object.Field{Ftype: types.BigInteger, Fvalue: pBI}
	pubKey.FieldTable["g"] = object.Field{Ftype: types.BigInteger, Fvalue: gBI}
	pubKey.FieldTable["l"] = object.Field{Ftype: types.Int, Fvalue: lValue}
	pubKey.FieldTable["params"] = object.Field{Ftype: types.Ref, Fvalue: kpgObj.FieldTable["paramSpec"].Fvalue}
	pubKey.FieldTable["algorithm"] = object.Field{Ftype: types.Ref, Fvalue: kpgObj.FieldTable["algorithm"].Fvalue}

	return privKey, pubKey, nil
}

// Helper: generate random integer in [1, max-1]
func randomBigInt(bits int) (*big.Int, error) {
	if bits <= 0 {
		return nil, errors.New("invalid bit size")
	}
	max := new(big.Int).Lsh(big.NewInt(1), uint(bits))
	max.Sub(max, big.NewInt(1)) // max = 2^bits - 1
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, err
	}
	// ensure >=1
	if n.Cmp(big.NewInt(0)) == 0 {
		n.Add(n, big.NewInt(1))
	}
	return n, nil
}
