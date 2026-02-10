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

// keyExchangeEncryptECDH encrypts a shared secret using ECDH public key
// This generates an ephemeral ECDH key pair, computes the shared secret,
// and encrypts the provided sharedSecret data
func keyExchangeEncryptECDH(publicKeyObj *object.Object, sharedSecret []byte) ([]byte, error) {
	// Extract ECDH public key
	pubKeyValue, ok := publicKeyObj.FieldTable["value"].Fvalue.(*ecdh.PublicKey)
	if !ok {
		return nil, fmt.Errorf("keyExchangeEncryptECDH: invalid public key object")
	}

	// Get the curve
	curve := pubKeyValue.Curve()

	// Generate ephemeral key pair on the same curve
	ephemeralPrivate, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptECDH: failed to generate ephemeral key: %w", err)
	}

	// Compute shared secret: ephemeral_private * recipient_public
	dhSharedSecret, err := ephemeralPrivate.ECDH(pubKeyValue)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptECDH: ECDH failed: %w", err)
	}

	// Derive encryption key using SHA-256
	hash := sha256.Sum256(dhSharedSecret)
	encryptionKey := hash[:]

	// Encrypt the shared secret using AES-GCM
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptECDH: cipher creation failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptECDH: GCM creation failed: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptECDH: nonce generation failed: %w", err)
	}

	// Encrypt the shared secret
	ciphertext := gcm.Seal(nonce, nonce, sharedSecret, nil)

	// Serialize ephemeral public key using Bytes()
	ephemeralPubBytes := ephemeralPrivate.PublicKey().Bytes()

	// Prepend ephemeral public key to ciphertext
	// Format: [ephemeral_public_key_length (2 bytes)][ephemeral_public_key][ciphertext]
	result := make([]byte, 2+len(ephemeralPubBytes)+len(ciphertext))

	// Write ephemeral public key length (big-endian, 2 bytes sufficient for EC keys)
	result[0] = byte(len(ephemeralPubBytes) >> 8)
	result[1] = byte(len(ephemeralPubBytes))

	// Write ephemeral public key
	copy(result[2:], ephemeralPubBytes)

	// Write ciphertext
	copy(result[2+len(ephemeralPubBytes):], ciphertext)

	return result, nil
}

// keyExchangeDecryptECDH decrypts a shared secret using ECDH private key
func keyExchangeDecryptECDH(privateKeyObj *object.Object, ciphertext []byte) ([]byte, error) {
	// Extract ECDH private key
	privKeyValue, ok := privateKeyObj.FieldTable["value"].Fvalue.(*ecdh.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("keyExchangeDecryptECDH: invalid private key object")
	}

	// Get the curve
	curve := privKeyValue.Curve()

	// Extract ephemeral public key from ciphertext
	if len(ciphertext) < 2 {
		return nil, fmt.Errorf("keyExchangeDecryptECDH: ciphertext too short")
	}

	// Read ephemeral public key length (big-endian)
	ephemeralPubLen := int(ciphertext[0])<<8 | int(ciphertext[1])

	if len(ciphertext) < 2+ephemeralPubLen {
		return nil, fmt.Errorf("keyExchangeDecryptECDH: ciphertext too short for ephemeral key")
	}

	// Extract ephemeral public key bytes
	ephemeralPubBytes := ciphertext[2 : 2+ephemeralPubLen]

	// Deserialize ephemeral public key using NewPublicKey
	ephemeralPublic, err := curve.NewPublicKey(ephemeralPubBytes)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeDecryptECDH: failed to parse ephemeral public key: %w", err)
	}

	// Extract actual ciphertext
	actualCiphertext := ciphertext[2+ephemeralPubLen:]

	// Compute shared secret: our_private * ephemeral_public
	dhSharedSecret, err := privKeyValue.ECDH(ephemeralPublic)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeDecryptECDH: ECDH failed: %w", err)
	}

	// Derive decryption key using SHA-256
	hash := sha256.Sum256(dhSharedSecret)
	decryptionKey := hash[:]

	// Decrypt using AES-GCM
	block, err := aes.NewCipher(decryptionKey)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeDecryptECDH: cipher creation failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeDecryptECDH: GCM creation failed: %w", err)
	}

	// Extract nonce from ciphertext
	nonceSize := gcm.NonceSize()
	if len(actualCiphertext) < nonceSize {
		return nil, fmt.Errorf("keyExchangeDecryptECDH: ciphertext too short for nonce")
	}

	nonce, encryptedData := actualCiphertext[:nonceSize], actualCiphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeDecryptECDH: decryption failed: %w", err)
	}

	return plaintext, nil
}
