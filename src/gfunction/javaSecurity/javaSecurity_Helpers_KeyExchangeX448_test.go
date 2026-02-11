package javaSecurity

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"

	"jacobin/src/object"

	"github.com/cloudflare/circl/dh/x448"
)

func TestKeyExchangeX448_Roundtrip(t *testing.T) {
	// Generate Recipient's Key Pair
	var privKey, pubKey x448.Key
	if _, err := io.ReadFull(rand.Reader, privKey[:]); err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}
	x448.KeyGen(&pubKey, &privKey)

	// Wrap in object.Object
	privKeyObj := object.MakeEmptyObject()
	privKeyObj.FieldTable["value"] = object.Field{Fvalue: privKey[:]}

	pubKeyObj := object.MakeEmptyObject()
	pubKeyObj.FieldTable["value"] = object.Field{Fvalue: pubKey[:]}

	sharedSecret := []byte("this is a secret shared message for X448")

	// Encrypt
	ciphertext, err := keyExchangeEncryptX448(pubKeyObj, sharedSecret)
	if err != nil {
		t.Fatalf("keyExchangeEncryptX448 failed: %v", err)
	}

	// Decrypt
	decrypted, err := keyExchangeDecryptX448(privKeyObj, ciphertext)
	if err != nil {
		t.Fatalf("keyExchangeDecryptX448 failed: %v", err)
	}

	if !bytes.Equal(sharedSecret, decrypted) {
		t.Errorf("Decrypted secret does not match. Expected %s, got %s", sharedSecret, decrypted)
	}
}

func TestKeyExchangeX448_InvalidKeys(t *testing.T) {
	// Test with invalid public key type
	pubKeyObj := object.MakeEmptyObject()
	pubKeyObj.FieldTable["value"] = object.Field{Fvalue: "not a byte slice"}

	_, err := keyExchangeEncryptX448(pubKeyObj, []byte("secret"))
	if err == nil {
		t.Error("Expected error for invalid public key type, got nil")
	}

	// Test with invalid public key length
	pubKeyObj.FieldTable["value"] = object.Field{Fvalue: []byte{1, 2, 3}}
	_, err = keyExchangeEncryptX448(pubKeyObj, []byte("secret"))
	if err == nil {
		t.Error("Expected error for invalid public key length, got nil")
	}

	// Test with invalid private key type
	privKeyObj := object.MakeEmptyObject()
	privKeyObj.FieldTable["value"] = object.Field{Fvalue: "not a byte slice"}

	_, err = keyExchangeDecryptX448(privKeyObj, make([]byte, x448.Size+12+16)) // enough for min ciphertext
	if err == nil {
		t.Error("Expected error for invalid private key type, got nil")
	}

	// Test with invalid private key length
	privKeyObj.FieldTable["value"] = object.Field{Fvalue: []byte{1, 2, 3}}
	_, err = keyExchangeDecryptX448(privKeyObj, make([]byte, x448.Size+12+16))
	if err == nil {
		t.Error("Expected error for invalid private key length, got nil")
	}
}

func TestKeyExchangeX448_InvalidCiphertext(t *testing.T) {
	var privKey, pubKey x448.Key
	io.ReadFull(rand.Reader, privKey[:])
	x448.KeyGen(&pubKey, &privKey)

	privKeyObj := object.MakeEmptyObject()
	privKeyObj.FieldTable["value"] = object.Field{Fvalue: privKey[:]}

	// Too short for ephemeral public key
	_, err := keyExchangeDecryptX448(privKeyObj, make([]byte, x448.Size-1))
	if err == nil {
		t.Error("Expected error for too short ciphertext (no ephemeral key), got nil")
	}

	// Too short for AES-GCM nonce
	_, err = keyExchangeDecryptX448(privKeyObj, make([]byte, x448.Size+1))
	if err == nil {
		t.Error("Expected error for too short ciphertext (no nonce), got nil")
	}
}

func TestKeyExchangeX448_DecryptionFailed(t *testing.T) {
	var privKey, pubKey x448.Key
	io.ReadFull(rand.Reader, privKey[:])
	x448.KeyGen(&pubKey, &privKey)

	privKeyObj := object.MakeEmptyObject()
	privKeyObj.FieldTable["value"] = object.Field{Fvalue: privKey[:]}
	pubKeyObj := object.MakeEmptyObject()
	pubKeyObj.FieldTable["value"] = object.Field{Fvalue: pubKey[:]}

	sharedSecret := []byte("secret")
	ciphertext, _ := keyExchangeEncryptX448(pubKeyObj, sharedSecret)

	// Tamper with ciphertext (last byte)
	ciphertext[len(ciphertext)-1] ^= 0xFF

	_, err := keyExchangeDecryptX448(privKeyObj, ciphertext)
	if err == nil {
		t.Error("Expected decryption failure for tampered ciphertext, got nil")
	}
}

func TestKeyExchangeX448_WeakSharedSecret(t *testing.T) {
	// X448 low-order points or invalid public keys can lead to all-zero shared secret.
	// The implementation checks for this.

	// An all-zero public key is one such input.
	pubKey := make([]byte, x448.Size)
	pubKeyObj := object.MakeEmptyObject()
	pubKeyObj.FieldTable["value"] = object.Field{Fvalue: pubKey}

	_, err := keyExchangeEncryptX448(pubKeyObj, []byte("secret"))
	if err == nil {
		t.Error("Expected error for weak public key, got nil")
	}

	// Also test decryption with weak ephemeral public key
	var privKey x448.Key
	io.ReadFull(rand.Reader, privKey[:])
	privKeyObj := object.MakeEmptyObject()
	privKeyObj.FieldTable["value"] = object.Field{Fvalue: privKey[:]}

	ciphertext := make([]byte, x448.Size+32) // enough room
	// ephemeral public key is all zeros at the beginning
	_, err = keyExchangeDecryptX448(privKeyObj, ciphertext)
	if err == nil {
		t.Error("Expected error for weak ephemeral public key in ciphertext, got nil")
	}
}
