/*
	Standard (Java 21):
	https://docs.oracle.com/en/java/javase/21/docs/specs/security/standard-names.html
*/

package security

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"hash"
	"math/big"
	"strings"
)

// ---------------------------
// KeyPair struct
// ---------------------------
type KeyPair struct {
	PrivateKey crypto.PrivateKey
	PublicKey  crypto.PublicKey
	Algorithm  string
}

// ---------------------------
// KeyPairGenerator struct
// ---------------------------
type KeyPairGenerator struct {
	Algorithm string
	KeySize   int
	CurveName string
}

// ---------------------------
// Key pair constructors
// ---------------------------

func NewKeyPairGenerator(algorithm string) (*KeyPairGenerator, error) {
	switch algorithm {
	case "RSA", "EC", "Ed25519", "DSA":
		return &KeyPairGenerator{Algorithm: algorithm}, nil
	default:
		return nil, errors.New("unsupported algorithm: " + algorithm)
	}
}

// Initialize key size
func (kpg *KeyPairGenerator) Initialize(keySize int) error {
	kpg.KeySize = keySize
	return nil
}

// ---------------------------
// Generate KeyPair
// ---------------------------
func (kpg *KeyPairGenerator) GenerateKeyPair() (*KeyPair, error) {
	switch kpg.Algorithm {
	case "RSA":
		if kpg.KeySize == 0 {
			kpg.KeySize = 2048
		}
		priv, err := rsa.GenerateKey(rand.Reader, kpg.KeySize)
		if err != nil {
			return nil, err
		}
		return &KeyPair{PrivateKey: priv, PublicKey: &priv.PublicKey, Algorithm: "RSA"}, nil

	case "EC":
		curve, err := selectCurve(kpg)
		if err != nil {
			return nil, err
		}
		priv, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			return nil, err
		}
		return &KeyPair{PrivateKey: priv, PublicKey: &priv.PublicKey, Algorithm: "EC"}, nil

	case "Ed25519":
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		return &KeyPair{PrivateKey: priv, PublicKey: pub, Algorithm: "Ed25519"}, nil

	case "DSA":
		var params dsa.Parameters
		switch {
		case kpg.KeySize >= 3072:
			if err := dsa.GenerateParameters(&params, rand.Reader, dsa.L3072N256); err != nil {
				return nil, err
			}
		case kpg.KeySize >= 2048:
			if err := dsa.GenerateParameters(&params, rand.Reader, dsa.L2048N256); err != nil {
				return nil, err
			}
		default:
			if err := dsa.GenerateParameters(&params, rand.Reader, dsa.L1024N160); err != nil {
				return nil, err
			}
		}
		priv := new(dsa.PrivateKey)
		priv.Parameters = params
		if err := dsa.GenerateKey(priv, rand.Reader); err != nil {
			return nil, err
		}
		return &KeyPair{PrivateKey: priv, PublicKey: &priv.PublicKey, Algorithm: "DSA"}, nil
	}

	return nil, errors.New("unsupported algorithm: " + kpg.Algorithm)
}

// ---------------------------
// Java-style signature parsing
// ---------------------------
func parseSignatureAlgorithm(alg string) (crypto.Hash, string, error) {
	alg = strings.ToUpper(alg)
	switch {
	case strings.HasSuffix(alg, "WITHRSA"):
		switch {
		case strings.HasPrefix(alg, "SHA256"):
			return crypto.SHA256, "RSA", nil
		case strings.HasPrefix(alg, "SHA512"):
			return crypto.SHA512, "RSA", nil
		default:
			return crypto.SHA256, "RSA", nil
		}
	case strings.HasSuffix(alg, "WITHECDSA"):
		switch {
		case strings.HasPrefix(alg, "SHA256"):
			return crypto.SHA256, "ECDSA", nil
		case strings.HasPrefix(alg, "SHA384"):
			return crypto.SHA384, "ECDSA", nil
		case strings.HasPrefix(alg, "SHA512"):
			return crypto.SHA512, "ECDSA", nil
		default:
			return crypto.SHA256, "ECDSA", nil
		}
	case strings.HasSuffix(alg, "WITHDSA"):
		switch {
		case strings.HasPrefix(alg, "SHA1"):
			return crypto.SHA1, "DSA", nil
		case strings.HasPrefix(alg, "SHA256"):
			return crypto.SHA256, "DSA", nil
		default:
			return crypto.SHA1, "DSA", nil
		}
	case alg == "ED25519":
		return 0, "Ed25519", nil
	default:
		return 0, "", errors.New("unsupported signature algorithm: " + alg)
	}
}

// ---------------------------
// Hash function mapping
// ---------------------------
func getHashFunc(h crypto.Hash) (hash.Hash, error) {
	switch h {
	case crypto.SHA256:
		return sha256.New(), nil
	case crypto.SHA384:
		return sha512.New384(), nil
	case crypto.SHA512:
		return sha512.New(), nil
	case 0:
		return nil, nil
	default:
		return nil, errors.New("unsupported hash function")
	}
}

// ---------------------------
// Sign / Verify
// ---------------------------
func (kp *KeyPair) SignWithAlgorithm(data []byte, sigAlg string) ([]byte, error) {
	hashAlg, keyType, err := parseSignatureAlgorithm(sigAlg)
	if err != nil {
		return nil, err
	}

	switch keyType {
	case "RSA":
		priv, ok := kp.PrivateKey.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("private key is not RSA")
		}
		hashFunc, err := getHashFunc(hashAlg)
		if err != nil {
			return nil, err
		}
		hashFunc.Write(data)
		return rsa.SignPKCS1v15(rand.Reader, priv, hashAlg, hashFunc.Sum(nil))

	case "ECDSA":
		priv, ok := kp.PrivateKey.(*ecdsa.PrivateKey)
		if !ok {
			return nil, errors.New("private key is not ECDSA")
		}
		hashFunc, err := getHashFunc(hashAlg)
		if err != nil {
			return nil, err
		}
		hashFunc.Write(data)
		return ecdsa.SignASN1(rand.Reader, priv, hashFunc.Sum(nil))

	case "Ed25519":
		priv, ok := kp.PrivateKey.(ed25519.PrivateKey)
		if !ok {
			return nil, errors.New("private key is not Ed25519")
		}
		return ed25519.Sign(priv, data), nil

	case "DSA":
		priv, ok := kp.PrivateKey.(*dsa.PrivateKey)
		if !ok {
			return nil, errors.New("private key is not DSA")
		}
		hashFunc, err := getHashFunc(hashAlg)
		if err != nil {
			return nil, err
		}
		hashFunc.Write(data)
		r, s, err := dsa.Sign(rand.Reader, priv, hashFunc.Sum(nil))
		if err != nil {
			return nil, err
		}
		return marshalDSASignature(r, s, priv.Parameters.Q.BitLen()), nil

	default:
		return nil, errors.New("unsupported key type for signing")
	}
}

func (kp *KeyPair) VerifyWithAlgorithm(message []byte, sig []byte, sigAlg string) (bool, error) {
	switch strings.ToUpper(sigAlg) {

	// -------------------- ECDSA --------------------
	case "ECDSA", "SHA256WITHECDSA":
		pub, ok := kp.PublicKey.(*ecdsa.PublicKey)
		if !ok {
			return false, errors.New("public key is not ECDSA")
		}

		byteLen := (pub.Curve.Params().BitSize + 7) / 8
		if len(sig) != 2*byteLen {
			return false, errors.New("invalid ECDSA signature length")
		}

		r := new(big.Int).SetBytes(sig[:byteLen])
		s := new(big.Int).SetBytes(sig[byteLen:])

		// Choose hashish based on curve
		var hashish []byte
		switch pub.Curve.Params().BitSize {
		case 256:
			h := sha256.Sum256(message)
			hashish = h[:]
		case 384:
			h := sha512.Sum384(message)
			hashish = h[:]
		case 521:
			h := sha512.Sum512(message)
			hashish = h[:]
		default:
			return false, errors.New("unsupported ECDSA curve")
		}

		valid := ecdsa.Verify(pub, hashish, r, s)
		return valid, nil

	// -------------------- RSA --------------------
	case "RSA", "SHA256WITHRSA":
		pub, ok := kp.PublicKey.(*rsa.PublicKey)
		if !ok {
			return false, errors.New("public key is not RSA")
		}

		hashish := sha256.Sum256(message)
		err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashish[:], sig)
		if err != nil {
			return false, err
		}
		return true, nil

	// -------------------- DSA --------------------
	case "DSA", "SHA256WITHDSA":
		pub, ok := kp.PublicKey.(*dsa.PublicKey)
		if !ok {
			return false, errors.New("public key is not DSA")
		}

		qLen := (pub.Q.BitLen() + 7) / 8
		if len(sig) != 2*qLen {
			return false, errors.New("invalid DSA signature length")
		}

		r := new(big.Int).SetBytes(sig[:qLen])
		s := new(big.Int).SetBytes(sig[qLen:])

		hashish := sha256.Sum256(message)
		valid := dsa.Verify(pub, hashish[:], r, s)
		return valid, nil

	// -------------------- Ed25519 --------------------
	case "ED25519":
		pub, ok := kp.PublicKey.(ed25519.PublicKey)
		if !ok {
			return false, errors.New("public key is not Ed25519")
		}
		valid := ed25519.Verify(pub, message, sig)
		return valid, nil

	// -------------------- Unsupported --------------------
	default:
		return false, errors.New("unsupported signature algorithm: " + sigAlg)
	}
}

// deserializeECDSASignature splits a fixed-length ECDSA signature into r and s
func deserializeECDSASignature(curve elliptic.Curve, sig []byte) (*big.Int, *big.Int, error) {
	byteLen := (curve.Params().BitSize + 7) / 8
	if len(sig) != 2*byteLen {
		return nil, nil, errors.New("invalid ECDSA signature length")
	}
	r := new(big.Int).SetBytes(sig[:byteLen])
	s := new(big.Int).SetBytes(sig[byteLen:])
	return r, s, nil
}

// deserializeDSASignature splits a fixed-length DSA signature into r and s.
// qLen is the length in bytes of the DSA subgroup order Q (pub.Q.BitLen() / 8).
func deserializeDSASignature(sig []byte, qLen int) (*big.Int, *big.Int, error) {
	if len(sig) != 2*qLen {
		return nil, nil, errors.New("invalid DSA signature length")
	}

	r := new(big.Int).SetBytes(sig[:qLen])
	s := new(big.Int).SetBytes(sig[qLen:])
	return r, s, nil
}

// ---------------------------
// DSA helper functions
// ---------------------------
func marshalDSASignature(r, s *big.Int, qBitLen int) []byte {
	size := (qBitLen + 7) / 8
	buf := make([]byte, size*2)
	r.FillBytes(buf[:size])
	s.FillBytes(buf[size:])
	return buf
}

func unmarshalDSASignature(sig []byte) (*big.Int, *big.Int, error) {
	if len(sig)%2 != 0 {
		return nil, nil, errors.New("invalid DSA signature length")
	}
	size := len(sig) / 2
	r := new(big.Int).SetBytes(sig[:size])
	s := new(big.Int).SetBytes(sig[size:])
	return r, s, nil
}

// ---------------------------
// Symmetric Encryption/Decryption (AES-GCM)
// ---------------------------
func EncryptWithAESGCM(key, plaintext []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	nonce := make([]byte, 12)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return nonce, ciphertext, nil
}

func DecryptWithAESGCM(key, nonce, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
