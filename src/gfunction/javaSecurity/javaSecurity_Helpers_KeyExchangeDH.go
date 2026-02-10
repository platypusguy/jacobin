package javaSecurity

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"jacobin/src/object"
	"math/big"
)

// keyExchangeEncryptDH encrypts a shared secret using DH public key
// This generates an ephemeral DH key pair, computes the shared secret,
// and encrypts the provided sharedSecret data
func keyExchangeEncryptDH(publicKeyObj *object.Object, sharedSecret []byte) ([]byte, error) {
	// Extract DH public key and parameters from object
	pubKeyValue := publicKeyObj.FieldTable["value"].Fvalue.(*big.Int)
	paramsObj := publicKeyObj.FieldTable["params"].Fvalue.(*object.Object)
	p := paramsObj.FieldTable["p"].Fvalue.(*big.Int)
	g := paramsObj.FieldTable["g"].Fvalue.(*big.Int)

	// Generate ephemeral key pair
	privateKeySize := 256 // bits
	if p.BitLen() < privateKeySize {
		privateKeySize = p.BitLen() - 1
	}

	ephemeralPrivate, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), uint(privateKeySize)))
	if err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptDH: failed to generate ephemeral key: %w", err)
	}

	// Ensure private key is at least 2
	two := big.NewInt(2)
	if ephemeralPrivate.Cmp(two) < 0 {
		ephemeralPrivate = two
	}

	// Generate ephemeral public key: g^private mod p
	ephemeralPublic := new(big.Int).Exp(g, ephemeralPrivate, p)

	// Validate recipient's public key
	pMinusTwo := new(big.Int).Sub(p, two)
	if pubKeyValue.Cmp(two) < 0 || pubKeyValue.Cmp(pMinusTwo) > 0 {
		return nil, fmt.Errorf("keyExchangeEncryptDH: invalid public key: out of range")
	}

	one := big.NewInt(1)
	pMinusOne := new(big.Int).Sub(p, one)
	if pubKeyValue.Cmp(one) == 0 || pubKeyValue.Cmp(pMinusOne) == 0 {
		return nil, fmt.Errorf("keyExchangeEncryptDH: invalid public key: weak value")
	}

	// Compute shared key: (their_public ^ our_ephemeral_private) mod p
	dhSharedSecret := new(big.Int).Exp(pubKeyValue, ephemeralPrivate, p)

	// Derive encryption key using SHA-256
	hash := sha256.Sum256(dhSharedSecret.Bytes())
	encryptionKey := hash[:]

	// Encrypt the shared secret using AES-GCM
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptDH: cipher creation failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptDH: GCM creation failed: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptDH: nonce generation failed: %w", err)
	}

	// Encrypt the shared secret
	ciphertext := gcm.Seal(nonce, nonce, sharedSecret, nil)

	// Prepend ephemeral public key to ciphertext
	// Format: [ephemeral_public_key_length (4 bytes)][ephemeral_public_key][ciphertext]
	ephemeralPubBytes := ephemeralPublic.Bytes()
	result := make([]byte, 4+len(ephemeralPubBytes)+len(ciphertext))

	// Write ephemeral public key length (big-endian)
	result[0] = byte(len(ephemeralPubBytes) >> 24)
	result[1] = byte(len(ephemeralPubBytes) >> 16)
	result[2] = byte(len(ephemeralPubBytes) >> 8)
	result[3] = byte(len(ephemeralPubBytes))

	// Write ephemeral public key
	copy(result[4:], ephemeralPubBytes)

	// Write ciphertext
	copy(result[4+len(ephemeralPubBytes):], ciphertext)

	return result, nil
}

// keyExchangeDecryptDH decrypts a shared secret using DH private key
func keyExchangeDecryptDH(privateKeyObj *object.Object, ciphertext []byte) ([]byte, error) {
	// Extract DH private key and parameters from object
	privKeyValue := privateKeyObj.FieldTable["value"].Fvalue.(*big.Int)
	paramsObj := privateKeyObj.FieldTable["params"].Fvalue.(*object.Object)
	p := paramsObj.FieldTable["p"].Fvalue.(*big.Int)

	// Extract ephemeral public key from ciphertext
	if len(ciphertext) < 4 {
		return nil, fmt.Errorf("keyExchangeDecryptDH: ciphertext too short")
	}

	// Read ephemeral public key length (big-endian)
	ephemeralPubLen := int(ciphertext[0])<<24 | int(ciphertext[1])<<16 | int(ciphertext[2])<<8 | int(ciphertext[3])

	if len(ciphertext) < 4+ephemeralPubLen {
		return nil, fmt.Errorf("keyExchangeDecryptDH: ciphertext too short for ephemeral key")
	}

	// Extract ephemeral public key
	ephemeralPubBytes := ciphertext[4 : 4+ephemeralPubLen]
	ephemeralPublic := new(big.Int).SetBytes(ephemeralPubBytes)

	// Extract actual ciphertext
	actualCiphertext := ciphertext[4+ephemeralPubLen:]

	// Validate ephemeral public key
	two := big.NewInt(2)
	pMinusTwo := new(big.Int).Sub(p, two)
	if ephemeralPublic.Cmp(two) < 0 || ephemeralPublic.Cmp(pMinusTwo) > 0 {
		return nil, fmt.Errorf("keyExchangeDecryptDH: invalid ephemeral public key: out of range")
	}

	one := big.NewInt(1)
	pMinusOne := new(big.Int).Sub(p, one)
	if ephemeralPublic.Cmp(one) == 0 || ephemeralPublic.Cmp(pMinusOne) == 0 {
		return nil, fmt.Errorf("invalid ephemeral public key: weak value")
	}

	// Compute shared key: (ephemeral_public ^ our_private) mod p
	dhSharedSecret := new(big.Int).Exp(ephemeralPublic, privKeyValue, p)

	// Derive decryption key using SHA-256
	hash := sha256.Sum256(dhSharedSecret.Bytes())
	decryptionKey := hash[:]

	// Decrypt using AES-GCM
	block, err := aes.NewCipher(decryptionKey)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeDecryptDH: cipher creation failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("GCM creation failed: %w", err)
	}

	// Extract nonce from ciphertext
	nonceSize := gcm.NonceSize()
	if len(actualCiphertext) < nonceSize {
		return nil, fmt.Errorf("keyExchangeDecryptDH: ciphertext too short for nonce")
	}

	nonce, encryptedData := actualCiphertext[:nonceSize], actualCiphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeDecryptDH: decryption failed: %w", err)
	}

	return plaintext, nil
}
