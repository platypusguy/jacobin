package javaSecurity

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"jacobin/src/object"

	"github.com/cloudflare/circl/dh/x448"
)

// keyExchangeEncryptX448 encrypts a shared secret using X448 public key
// This generates an ephemeral X448 key pair, computes the shared secret,
// and encrypts the provided sharedSecret data
func keyExchangeEncryptX448(publicKeyObj *object.Object, sharedSecret []byte) ([]byte, error) {
	// Extract X448 public key bytes (56 bytes)
	pubKeyBytes, ok := publicKeyObj.FieldTable["value"].Fvalue.([]byte)
	if !ok {
		return nil, fmt.Errorf("keyExchangeEncryptX448: invalid public key object: expected []byte")
	}

	// Validate public key length
	if len(pubKeyBytes) != x448.Size {
		return nil, fmt.Errorf("keyExchangeEncryptX448: invalid public key length: expected %d bytes, got %d", x448.Size, len(pubKeyBytes))
	}

	// Parse recipient's public key
	var recipientPubKey x448.Key
	copy(recipientPubKey[:], pubKeyBytes)

	// Generate ephemeral X448 key pair
	var ephemeralPrivate x448.Key
	var ephemeralPublic x448.Key

	// Generate random private key
	if _, err := io.ReadFull(rand.Reader, ephemeralPrivate[:]); err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptX448: failed to generate ephemeral private key: %w", err)
	}

	// Compute ephemeral public key
	x448.KeyGen(&ephemeralPublic, &ephemeralPrivate)

	// Compute shared secret using X448
	var dhSharedSecret x448.Key
	x448.Shared(&dhSharedSecret, &ephemeralPrivate, &recipientPubKey)

	// Check for low-order points (all-zero shared secret is invalid)
	allZero := true
	for _, b := range dhSharedSecret {
		if b != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		return nil, fmt.Errorf("keyExchangeEncryptX448: invalid public key: contributes to weak shared secret")
	}

	// Derive encryption key using SHA-256
	hash := sha256.Sum256(dhSharedSecret[:])
	encryptionKey := hash[:]

	// Encrypt the shared secret using AES-GCM
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptX448: cipher creation failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptX448: GCM creation failed: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptX448: nonce generation failed: %w", err)
	}

	// Encrypt the shared secret
	ciphertext := gcm.Seal(nonce, nonce, sharedSecret, nil)

	// Prepend ephemeral public key to ciphertext
	// Format: [ephemeral_public_key (56 bytes)][ciphertext]
	result := make([]byte, x448.Size+len(ciphertext))

	// Write ephemeral public key
	copy(result[:x448.Size], ephemeralPublic[:])

	// Write ciphertext
	copy(result[x448.Size:], ciphertext)

	return result, nil
}

// keyExchangeDecryptX448 decrypts a shared secret using X448 private key
func keyExchangeDecryptX448(privateKeyObj *object.Object, ciphertext []byte) ([]byte, error) {
	// Extract X448 private key bytes (56 bytes)
	privKeyBytes, ok := privateKeyObj.FieldTable["value"].Fvalue.([]byte)
	if !ok {
		return nil, fmt.Errorf("keyExchangeDecryptX448: invalid private key object: expected []byte")
	}

	// Validate private key length
	if len(privKeyBytes) != x448.Size {
		return nil, fmt.Errorf("keyExchangeDecryptX448: invalid private key length: expected %d bytes, got %d", x448.Size, len(privKeyBytes))
	}

	// Validate ciphertext length
	if len(ciphertext) < x448.Size {
		return nil, fmt.Errorf("keyExchangeDecryptX448: ciphertext too short: must be at least %d bytes", x448.Size)
	}

	// Parse our private key
	var ourPrivateKey x448.Key
	copy(ourPrivateKey[:], privKeyBytes)

	// Extract ephemeral public key (first 56 bytes)
	var ephemeralPublic x448.Key
	copy(ephemeralPublic[:], ciphertext[:x448.Size])

	// Extract actual ciphertext
	actualCiphertext := ciphertext[x448.Size:]

	// Compute shared secret using X448
	var dhSharedSecret x448.Key
	x448.Shared(&dhSharedSecret, &ourPrivateKey, &ephemeralPublic)

	// Check for low-order points (all-zero shared secret is invalid)
	allZero := true
	for _, b := range dhSharedSecret {
		if b != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		return nil, fmt.Errorf("keyExchangeDecryptX448: invalid ephemeral public key: contributes to weak shared secret")
	}

	// Derive decryption key using SHA-256
	hash := sha256.Sum256(dhSharedSecret[:])
	decryptionKey := hash[:]

	// Decrypt using AES-GCM
	block, err := aes.NewCipher(decryptionKey)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeDecryptX448: cipher creation failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeDecryptX448: GCM creation failed: %w", err)
	}

	// Extract nonce from ciphertext
	nonceSize := gcm.NonceSize()
	if len(actualCiphertext) < nonceSize {
		return nil, fmt.Errorf("keyExchangeDecryptX448: ciphertext too short for nonce")
	}

	nonce, encryptedData := actualCiphertext[:nonceSize], actualCiphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeDecryptX448: decryption failed: %w", err)
	}

	return plaintext, nil
}
