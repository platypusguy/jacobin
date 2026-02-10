package javaSecurity

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"jacobin/src/object"
)

// keyExchangeEncryptRSA encrypts a shared secret using RSA public key
func keyExchangeEncryptRSA(publicKeyObj *object.Object, sharedSecret []byte) ([]byte, error) {
	// Extract RSA public key
	pubKeyValue, ok := publicKeyObj.FieldTable["value"].Fvalue.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("keyExchangeEncryptRSA: invalid public key object")
	}

	// Encrypt the shared secret using RSA-OAEP
	// OAEP is recommended over PKCS#1 v1.5 for security
	rng := rand.Reader
	label := []byte("") // Optional label, usually empty

	ciphertext, err := rsa.EncryptOAEP(
		sha256.New(),
		rng,
		pubKeyValue,
		sharedSecret,
		label,
	)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeEncryptRSA: encryption failed: %w", err)
	}

	return ciphertext, nil
}

// keyExchangeDecryptRSA decrypts a shared secret using RSA private key
func keyExchangeDecryptRSA(privateKeyObj *object.Object, ciphertext []byte) ([]byte, error) {
	// Extract RSA private key
	privKeyValue, ok := privateKeyObj.FieldTable["value"].Fvalue.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("keyExchangeDecryptRSA: invalid private key object")
	}

	// Decrypt the shared secret using RSA-OAEP
	rng := rand.Reader
	label := []byte("") // Must match the label used during encryption

	plaintext, err := rsa.DecryptOAEP(
		sha256.New(),
		rng,
		privKeyValue,
		ciphertext,
		label,
	)
	if err != nil {
		return nil, fmt.Errorf("keyExchangeDecryptRSA: decryption failed: %w", err)
	}

	return plaintext, nil
}
