/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaSecurity

import (
	"jacobin/src/globals"
	"jacobin/src/object"
	"jacobin/src/types"
	"testing"
)

func makeSecureRandomObj() *object.Object {
	return object.MakeEmptyObjectWithClassName(&secureRandomClassName)
}

func TestSecureRandom_Init_And_SetSeed(t *testing.T) {
	globals.InitStringPool()

	sr := makeSecureRandomObj()
	if ret := secureRandomInit([]interface{}{sr}); ret != nil {
		t.Fatalf("secureRandomInit returned error: %v", ret)
	}
	// seed set by init
	if _, ok := sr.FieldTable["seed"]; !ok {
		t.Fatalf("expected seed field after init")
	}

	// setSeed with int64
	want := int64(0x1122334455667788)
	if ret := secureRandomSetSeed([]interface{}{sr, want}); ret != nil {
		t.Fatalf("setSeed(int64) returned error: %v", ret)
	}
	got := sr.FieldTable["seed"].Fvalue.(int64)
	if got != want {
		t.Fatalf("setSeed(int64) mismatch: expected %d, got %d", want, got)
	}

	// setSeed with byte array (JavaByte array)
	bytes := []byte{0xAA, 0xBB, 0xCC, 0xDD, 0x01, 0x02, 0x03, 0x04}
	jb := object.JavaByteArrayFromGoByteArray(bytes)
	byteArrObj := object.StringObjectFromJavaByteArray(jb)
	if ret := secureRandomSetSeed([]interface{}{sr, byteArrObj}); ret != nil {
		t.Fatalf("setSeed([B) returned error: %v", ret)
	}
	// Seed becomes big-endian int64 of bytes
	want64 := types.BytesToInt64BE(bytes)
	got64 := sr.FieldTable["seed"].Fvalue.(int64)
	if got64 != want64 {
		t.Fatalf("setSeed([B) mismatch: expected %d, got %d", want64, got64)
	}
}

func TestSecureRandom_NextBytes(t *testing.T) {
	globals.InitStringPool()

	sr := makeSecureRandomObj()
	_ = secureRandomInit([]interface{}{sr})

	// Prepare a byte array object of length 16
	buf := make([]byte, 16)
	jb := object.JavaByteArrayFromGoByteArray(buf)
	baObj := object.StringObjectFromJavaByteArray(jb)

	if ret := secureRandomNextBytes([]interface{}{sr, baObj}); ret != nil {
		t.Fatalf("nextBytes returned error: %v", ret)
	}

	gotJB := object.JavaByteArrayFromStringObject(baObj)
	if len(gotJB) != 16 {
		t.Fatalf("nextBytes length mismatch: expected 16, got %d", len(gotJB))
	}
	// Ensure not all zeros
	allZero := true
	for _, b := range gotJB {
		if b != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		t.Fatalf("nextBytes produced all zeros, unlikely for secure RNG")
	}

	// Call again and expect different content (very high probability)
	prev := make([]types.JavaByte, len(gotJB))
	copy(prev, gotJB)
	if ret := secureRandomNextBytes([]interface{}{sr, baObj}); ret != nil {
		t.Fatalf("nextBytes (second) returned error: %v", ret)
	}
	gotJB2 := object.JavaByteArrayFromStringObject(baObj)
	same := true
	for i := range gotJB2 {
		if gotJB2[i] != prev[i] {
			same = false
			break
		}
	}
	if same {
		t.Fatalf("two nextBytes calls yielded identical output; extremely unlikely")
	}
}

func TestSecureRandom_NextIntFloatBoolean(t *testing.T) {
	globals.InitStringPool()

	sr := makeSecureRandomObj()
	_ = secureRandomInit([]interface{}{sr})

	// nextInt (and nextLong maps to same impl)
	vi := secureRandomNextInt([]interface{}{sr})
	if _, ok := vi.(int64); !ok {
		t.Fatalf("nextInt did not return int64, got %T", vi)
	}

	// nextFloat used for nextFloat and nextDouble
	vf := secureRandomNextFloat([]interface{}{sr})
	f64, ok := vf.(float64)
	if !ok {
		t.Fatalf("nextFloat did not return float64, got %T", vf)
	}
	if !(f64 >= 0.0 && f64 < 1.0) {
		t.Fatalf("nextFloat not in [0,1): %v", f64)
	}

	// nextBoolean returns Java boolean constants
	vb := secureRandomNextBoolean([]interface{}{sr})
	if vb != types.JavaBoolTrue && vb != types.JavaBoolFalse {
		t.Fatalf("nextBoolean returned invalid value: %v", vb)
	}
}

func TestSecureRandom_GetSeed_And_GenerateSeed(t *testing.T) {
	globals.InitStringPool()

	// getSeed(size) is static-like and takes only the size in params
	out := secureRandomGetSeed([]interface{}{int64(8)})
	arrObj, ok := out.(*object.Object)
	if !ok {
		t.Fatalf("getSeed did not return object, got %T", out)
	}
	jb := object.JavaByteArrayFromStringObject(arrObj)
	if len(jb) != 8 {
		t.Fatalf("getSeed length mismatch: expected 8, got %d", len(jb))
	}

	// generateSeed on an instance
	sr := makeSecureRandomObj()
	_ = secureRandomInit([]interface{}{sr})

	out2 := secureRandomGenerateSeed([]interface{}{sr, int64(8)})
	arrObj2, ok := out2.(*object.Object)
	if !ok {
		t.Fatalf("generateSeed did not return object, got %T", out2)
	}
	jb2 := object.JavaByteArrayFromStringObject(arrObj2)
	if len(jb2) != 8 {
		t.Fatalf("generateSeed length mismatch: expected 8, got %d", len(jb2))
	}
	// Expect different output between calls with very high probability
	same := true
	for i := range jb2 {
		if jb2[i] != jb[i] {
			same = false
			break
		}
	}
	if same {
		t.Fatalf("getSeed and generateSeed returned identical bytes; extremely unlikely")
	}

	// invalid generateSeed size
	if err := secureRandomGenerateSeed([]interface{}{sr, int64(0)}); err == nil {
		t.Fatalf("expected error for non-positive generateSeed size")
	}
}

func TestSecureRandom_GetAlgorithm_ToString_GetInstance(t *testing.T) {
	globals.InitStringPool()

	// getAlgorithm and toString should return fixed string
	alg := secureRandomGetAlgorithm([]interface{}{})
	if s := object.GoStringFromStringObject(alg.(*object.Object)); s != "go/crypto/rand" {
		t.Fatalf("getAlgorithm mismatch: %q", s)
	}
	ts := secureRandomToString([]interface{}{})
	if s := object.GoStringFromStringObject(ts.(*object.Object)); s != "go/crypto/rand" {
		t.Fatalf("toString mismatch: %q", s)
	}

	// getInstance returns the same object passed in (reseeded)
	sr := makeSecureRandomObj()
	got := secureRandomGetInstance([]interface{}{sr})
	if got != sr {
		t.Fatalf("getInstance did not return the provided object")
	}

	// getInstanceStrong returns a new object
	strong := secureRandomGetInstanceStrong([]interface{}{})
	if _, ok := strong.(*object.Object); !ok {
		t.Fatalf("getInstanceStrong did not return an object, got %T", strong)
	}
}

func TestSecureRandom_Reseed_NoArg(t *testing.T) {
	globals.InitStringPool()

	sr := makeSecureRandomObj()
	_ = secureRandomInit([]interface{}{sr})

	before := sr.FieldTable["seed"].Fvalue.(int64)
	// reseed()V overload mapped to same function; we still pass the object as param
	if ret := secureRandomReseed([]interface{}{sr}); ret != nil {
		t.Fatalf("reseed returned error: %v", ret)
	}
	after := sr.FieldTable["seed"].Fvalue.(int64)
	if before == after {
		t.Fatalf("reseed did not change the seed value")
	}
}
