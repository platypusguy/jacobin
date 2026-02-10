package javaSecurity

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"jacobin/src/object"
)

// keyExchangeEncryptX25519 encrypts a shared secret using X25519 public key
// This generates an ephemeral X25519 key pair, computes the shared secret,
// and encrypts the provided sharedSecret data
func keyExchangeEncryptX25519(publicKeyObj *object.Object, sharedSecret []byte) ([]byte, error) {
	// Extract X25519 public key bytes (32 bytes)
	pubKeyBytes, ok := publicKeyObj.FieldTable["value"].Fvalue.([]byte)
	if !ok {
		return nil, fmt.Errorf("keyExchangeEncryptX25519: invalid public key object: expected []byte")
	}

	// Validate public key length
	if len(pubKeyBytes) != 32 {
		return nil, fmt.Errorf("keyExchangeEncryptX25519: invalid public key length: expected 32 bytes, got %d", len(pubKeyBytes))
	}

	// Parse the public key
	recipientPubKey, err := ecdh.X25519().NewPublicKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptX25519: failed to parse public key: %w", err)
	}

	// Generate ephemeral X25519 key pair
	ephemeralPrivate, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptX25519: failed to generate ephemeral key: %w", err)
	}

	// Compute shared secret using ECDH
	dhSharedSecret, err := ephemeralPrivate.ECDH(recipientPubKey)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptX25519: ECDH failed: %w", err)
	}

	// Derive encryption key using SHA-256
	hash := sha256.Sum256(dhSharedSecret)
	encryptionKey := hash[:]

	// Encrypt the shared secret using AES-GCM
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptX25519: cipher creation failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptX25519: GCM creation failed: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptX25519: nonce generation failed: %w", err)
	}

	// Encrypt the shared secret
	ciphertext := gcm.Seal(nonce, nonce, sharedSecret, nil)

	// Get ephemeral public key bytes
	ephemeralPublicBytes := ephemeralPrivate.PublicKey().Bytes()

	// Prepend ephemeral public key to ciphertext
	// Format: [ephemeral_public_key (32 bytes)][ciphertext]
	result := make([]byte, 32+len(ciphertext))

	// Write ephemeral public key
	copy(result[:32], ephemeralPublicBytes)

	// Write ciphertext
	copy(result[32:], ciphertext)

	return result, nil
}

// keyExchangeDecryptX25519 decrypts a shared secret using X25519 private key
func keyExchangeDecryptX25519(privateKeyObj *object.Object, ciphertext []byte) ([]byte, error) {
	// Extract X25519 private key bytes (32 bytes)
	privKeyBytes, ok := privateKeyObj.FieldTable["value"].Fvalue.([]byte)
	if !ok {
		return nil, fmt.Errorf("keyExchangeDecryptX25519: invalid private key object: expected []byte")
	}

	// Validate private key length
	if len(privKeyBytes) != 32 {
		return nil, fmt.Errorf("keyExchangeDecryptX25519: invalid private key length: expected 32 bytes, got %d", len(privKeyBytes))
	}

	// Validate ciphertext length
	if len(ciphertext) < 32 {
		return nil, fmt.Errorf("keyExchangeDecryptX25519: ciphertext too short: must be at least 32 bytes")
	}

	// Extract ephemeral public key (first 32 bytes)
	ephemeralPublicBytes := ciphertext[:32]

	// Extract actual ciphertext
	actualCiphertext := ciphertext[32:]

	// Parse our private key
	ourPrivateKey, err := ecdh.X25519().NewPrivateKey(privKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeDecryptX25519: failed to parse private key: %w", err)
	}

	// Parse ephemeral public key
	ephemeralPublic, err := ecdh.X25519().NewPublicKey(ephemeralPublicBytes)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeDecryptX25519: failed to parse ephemeral public key: %w", err)
	}

	// Compute shared secret using ECDH
	dhSharedSecret, err := ourPrivateKey.ECDH(ephemeralPublic)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeDecryptX25519: ECDH failed: %w", err)
	}

	// Derive decryption key using SHA-256
	hash := sha256.Sum256(dhSharedSecret)
	decryptionKey := hash[:]

	// Decrypt using AES-GCM
	block, err := aes.NewCipher(decryptionKey)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeDecryptX25519: cipher creation failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeDecryptX25519: GCM creation failed: %w", err)
	}

	// Extract nonce from ciphertext
	nonceSize := gcm.NonceSize()
	if len(actualCiphertext) < nonceSize {
		return nil, fmt.Errorf("keyExchangeDecryptX25519: ciphertext too short for nonce")
	}

	nonce, encryptedData := actualCiphertext[:nonceSize], actualCiphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeDecryptX25519: decryption failed: %w", err)
	}

	return plaintext, nil
}
