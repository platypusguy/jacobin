package javaSecurity

import (
	"crypto/ecdh"
	"crypto/rand"
	"jacobin/src/object"
	"testing"
)

func TestKeyExchangeX25519_Roundtrip(t *testing.T) {
	// Generate a recipient X25519 key pair
	recipientPrivate, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate recipient key: %v", err)
	}
	recipientPublic := recipientPrivate.PublicKey()

	// Prepare public key object
	pubObj := object.MakeEmptyObject()
	pubObj.FieldTable["value"] = object.Field{Fvalue: recipientPublic.Bytes()}

	// Prepare private key object
	privObj := object.MakeEmptyObject()
	privObj.FieldTable["value"] = object.Field{Fvalue: recipientPrivate.Bytes()}

	sharedSecret := []byte("this is a highly secret shared key")

	// Encrypt
	ciphertext, err := keyExchangeEncryptX25519(pubObj, sharedSecret)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Decrypt
	plaintext, err := keyExchangeDecryptX25519(privObj, ciphertext)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if string(plaintext) != string(sharedSecret) {
		t.Errorf("Plaintext mismatch: expected %s, got %s", sharedSecret, plaintext)
	}
}

func TestKeyExchangeX25519_InvalidKeys(t *testing.T) {
	sharedSecret := []byte("secret")

	t.Run("InvalidPublicKeyType", func(t *testing.T) {
		pubObj := object.MakeEmptyObject()
		pubObj.FieldTable["value"] = object.Field{Fvalue: "not bytes"}
		_, err := keyExchangeEncryptX25519(pubObj, sharedSecret)
		if err == nil {
			t.Error("Expected error for invalid public key type, got nil")
		}
	})

	t.Run("InvalidPublicKeyLength", func(t *testing.T) {
		pubObj := object.MakeEmptyObject()
		pubObj.FieldTable["value"] = object.Field{Fvalue: []byte{1, 2, 3}}
		_, err := keyExchangeEncryptX25519(pubObj, sharedSecret)
		if err == nil {
			t.Error("Expected error for invalid public key length, got nil")
		}
	})

	t.Run("InvalidPrivateKeyType", func(t *testing.T) {
		privObj := object.MakeEmptyObject()
		privObj.FieldTable["value"] = object.Field{Fvalue: 12345}
		_, err := keyExchangeDecryptX25519(privObj, []byte("too short"))
		if err == nil {
			t.Error("Expected error for invalid private key type, got nil")
		}
	})

	t.Run("InvalidPrivateKeyLength", func(t *testing.T) {
		privObj := object.MakeEmptyObject()
		privObj.FieldTable["value"] = object.Field{Fvalue: []byte{1, 2, 3}}
		_, err := keyExchangeDecryptX25519(privObj, make([]byte, 32))
		if err == nil {
			t.Error("Expected error for invalid private key length, got nil")
		}
	})
}

func TestKeyExchangeX25519_InvalidCiphertext(t *testing.T) {
	recipientPrivate, _ := ecdh.X25519().GenerateKey(rand.Reader)
	privObj := object.MakeEmptyObject()
	privObj.FieldTable["value"] = object.Field{Fvalue: recipientPrivate.Bytes()}

	t.Run("TooShort", func(t *testing.T) {
		_, err := keyExchangeDecryptX25519(privObj, []byte("short"))
		if err == nil {
			t.Error("Expected error for short ciphertext, got nil")
		}
	})

	t.Run("NoNonce", func(t *testing.T) {
		// 32 bytes ephemeral public key, but nothing else
		_, err := keyExchangeDecryptX25519(privObj, make([]byte, 32))
		if err == nil {
			t.Error("Expected error for ciphertext missing nonce, got nil")
		}
	})
}

func TestKeyExchangeX25519_DecryptionFailed(t *testing.T) {
	recipientPrivate, _ := ecdh.X25519().GenerateKey(rand.Reader)
	recipientPublic := recipientPrivate.PublicKey()

	pubObj := object.MakeEmptyObject()
	pubObj.FieldTable["value"] = object.Field{Fvalue: recipientPublic.Bytes()}
	privObj := object.MakeEmptyObject()
	privObj.FieldTable["value"] = object.Field{Fvalue: recipientPrivate.Bytes()}

	sharedSecret := []byte("secret")
	ciphertext, _ := keyExchangeEncryptX25519(pubObj, sharedSecret)

	// Tamper with ciphertext (after the 32-byte ephemeral public key)
	ciphertext[35] ^= 0xFF

	_, err := keyExchangeDecryptX25519(privObj, ciphertext)
	if err == nil {
		t.Error("Expected decryption error for tampered ciphertext, got nil")
	}
}
