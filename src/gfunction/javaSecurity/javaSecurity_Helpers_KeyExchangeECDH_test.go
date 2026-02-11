package javaSecurity

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"testing"

	"jacobin/src/object"
)

func TestKeyExchangeECDH_Roundtrip(t *testing.T) {
	curves := []struct {
		name  string
		curve ecdh.Curve
	}{
		{"P256", ecdh.P256()},
		{"P384", ecdh.P384()},
		{"P521", ecdh.P521()},
		{"X25519", ecdh.X25519()},
	}

	for _, tc := range curves {
		t.Run(tc.name, func(t *testing.T) {
			// Generate Recipient's Key Pair
			privKey, err := tc.curve.GenerateKey(rand.Reader)
			if err != nil {
				t.Fatalf("Failed to generate key for %s: %v", tc.name, err)
			}
			pubKey := privKey.PublicKey()

			// Wrap in object.Object
			privKeyObj := object.MakeEmptyObject()
			privKeyObj.FieldTable["value"] = object.Field{Fvalue: privKey}

			pubKeyObj := object.MakeEmptyObject()
			pubKeyObj.FieldTable["value"] = object.Field{Fvalue: pubKey}

			sharedSecret := []byte("this is a secret shared message for ECDH " + tc.name)

			// Encrypt
			ciphertext, err := keyExchangeEncryptECDH(pubKeyObj, sharedSecret)
			if err != nil {
				t.Fatalf("keyExchangeEncryptECDH failed: %v", err)
			}

			// Decrypt
			decrypted, err := keyExchangeDecryptECDH(privKeyObj, ciphertext)
			if err != nil {
				t.Fatalf("keyExchangeDecryptECDH failed: %v", err)
			}

			if !bytes.Equal(sharedSecret, decrypted) {
				t.Errorf("Decrypted secret does not match. Expected %s, got %s", sharedSecret, decrypted)
			}
		})
	}
}

func TestKeyExchangeECDH_InvalidKeys(t *testing.T) {
	// Test with nil or wrong type in FieldTable
	pubKeyObj := object.MakeEmptyObject()
	pubKeyObj.FieldTable["value"] = object.Field{Fvalue: "not a public key"}

	_, err := keyExchangeEncryptECDH(pubKeyObj, []byte("secret"))
	if err == nil {
		t.Error("Expected error for invalid public key type, got nil")
	}

	privKeyObj := object.MakeEmptyObject()
	privKeyObj.FieldTable["value"] = object.Field{Fvalue: "not a private key"}

	_, err = keyExchangeDecryptECDH(privKeyObj, []byte{0, 0, 0, 0})
	if err == nil {
		t.Error("Expected error for invalid private key type, got nil")
	}
}

func TestKeyExchangeECDH_InvalidCiphertext(t *testing.T) {
	privKey, _ := ecdh.P256().GenerateKey(rand.Reader)
	privKeyObj := object.MakeEmptyObject()
	privKeyObj.FieldTable["value"] = object.Field{Fvalue: privKey}

	// Too short for length prefix
	_, err := keyExchangeDecryptECDH(privKeyObj, []byte{0})
	if err == nil {
		t.Error("Expected error for too short ciphertext (length prefix), got nil")
	}

	// Too short for declared ephemeral key length
	_, err = keyExchangeDecryptECDH(privKeyObj, []byte{0, 65})
	if err == nil {
		t.Error("Expected error for too short ciphertext (key data), got nil")
	}

	// Invalid ephemeral key bytes
	invalidKeyData := make([]byte, 67)
	invalidKeyData[1] = 65 // length 65
	_, err = keyExchangeDecryptECDH(privKeyObj, invalidKeyData)
	if err == nil {
		t.Error("Expected error for invalid ephemeral key data, got nil")
	}
}

func TestKeyExchangeECDH_DecryptionFailed(t *testing.T) {
	privKey, _ := ecdh.P256().GenerateKey(rand.Reader)
	pubKey := privKey.PublicKey()

	privKeyObj := object.MakeEmptyObject()
	privKeyObj.FieldTable["value"] = object.Field{Fvalue: privKey}
	pubKeyObj := object.MakeEmptyObject()
	pubKeyObj.FieldTable["value"] = object.Field{Fvalue: pubKey}

	sharedSecret := []byte("secret")
	ciphertext, _ := keyExchangeEncryptECDH(pubKeyObj, sharedSecret)

	// Tamper with ciphertext (after length and ephemeral key)
	// Format: [2 bytes len][ephemeral key][nonce][encrypted data + tag]
	// For P256, ephemeral key length is usually 65 (uncompressed).
	// So tamper at index 2 + 65 + something.
	tamperIdx := len(ciphertext) - 1
	ciphertext[tamperIdx] ^= 0xFF

	_, err := keyExchangeDecryptECDH(privKeyObj, ciphertext)
	if err == nil {
		t.Error("Expected decryption failure for tampered ciphertext, got nil")
	}
}
