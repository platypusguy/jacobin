/*
	Standard (Java 21):
	https://docs.oracle.com/en/java/javase/21/docs/specs/security/standard-names.html
*/

package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"strings"
)

// ---------------------------
// RSAEncrypt / RSADecrypt
// ---------------------------

func (kp *KeyPair) RSAEncrypt(plain []byte, alg string) ([]byte, error) {
	alg = strings.ToUpper(alg)
	switch alg {
	case "RSA", "RSA/ECB/PKCS1PADDING":
		pub, ok := kp.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("public key is not RSA")
		}
		return rsa.EncryptPKCS1v15(rand.Reader, pub, plain)
	case "RSA/OAEP":
		pub, ok := kp.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("public key is not RSA")
		}
		h := sha256.New()
		return rsa.EncryptOAEP(h, rand.Reader, pub, plain, nil)
	default:
		return nil, errors.New("unsupported encryption algorithm: " + alg)
	}
}

func (kp *KeyPair) RSADecrypt(ciphertext []byte, alg string) ([]byte, error) {
	alg = strings.ToUpper(alg)
	switch alg {
	case "RSA", "RSA/ECB/PKCS1PADDING":
		priv, ok := kp.PrivateKey.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("private key is not RSA")
		}
		return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	case "RSA/OAEP":
		priv, ok := kp.PrivateKey.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("private key is not RSA")
		}
		h := sha256.New()
		return rsa.DecryptOAEP(h, rand.Reader, priv, ciphertext, nil)
	default:
		return nil, errors.New("unsupported decryption algorithm: " + alg)
	}
}
