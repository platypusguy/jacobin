package javaSecurity

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"jacobin/src/object"
)

func TestKeyExchangeRSA_Roundtrip(t *testing.T) {
	// Generate RSA key pair
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}
	pubKey := &privKey.PublicKey

	// Wrap in object.Object
	privKeyObj := object.MakeEmptyObject()
	privKeyObj.FieldTable["value"] = object.Field{Fvalue: privKey}

	pubKeyObj := object.MakeEmptyObject()
	pubKeyObj.FieldTable["value"] = object.Field{Fvalue: pubKey}

	sharedSecret := []byte("this is a top secret RSA shared secret")

	// Encrypt
	ciphertext, err := keyExchangeEncryptRSA(pubKeyObj, sharedSecret)
	if err != nil {
		t.Fatalf("keyExchangeEncryptRSA failed: %v", err)
	}

	// Decrypt
	decrypted, err := keyExchangeDecryptRSA(privKeyObj, ciphertext)
	if err != nil {
		t.Fatalf("keyExchangeDecryptRSA failed: %v", err)
	}

	if !bytes.Equal(sharedSecret, decrypted) {
		t.Errorf("Decrypted secret does not match. Expected %v, got %v", sharedSecret, decrypted)
	}
}

func TestKeyExchangeRSA_InvalidKeys(t *testing.T) {
	// Test with invalid public key type
	pubKeyObj := object.MakeEmptyObject()
	pubKeyObj.FieldTable["value"] = object.Field{Fvalue: "not a public key"}

	_, err := keyExchangeEncryptRSA(pubKeyObj, []byte("secret"))
	if err == nil {
		t.Error("Expected error for invalid public key type, got nil")
	}

	// Test with invalid private key type
	privKeyObj := object.MakeEmptyObject()
	privKeyObj.FieldTable["value"] = object.Field{Fvalue: "not a private key"}

	_, err = keyExchangeDecryptRSA(privKeyObj, []byte("ciphertext"))
	if err == nil {
		t.Error("Expected error for invalid private key type, got nil")
	}
}

func TestKeyExchangeRSA_DecryptionFailed(t *testing.T) {
	privKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	pubKey := &privKey.PublicKey

	pubKeyObj := object.MakeEmptyObject()
	pubKeyObj.FieldTable["value"] = object.Field{Fvalue: pubKey}
	privKeyObj := object.MakeEmptyObject()
	privKeyObj.FieldTable["value"] = object.Field{Fvalue: privKey}

	sharedSecret := []byte("secret")
	ciphertext, _ := keyExchangeEncryptRSA(pubKeyObj, sharedSecret)

	// Tamper with ciphertext
	ciphertext[len(ciphertext)-1] ^= 0xFF

	_, err := keyExchangeDecryptRSA(privKeyObj, ciphertext)
	if err == nil {
		t.Error("Expected decryption failure for tampered ciphertext, got nil")
	}
}
