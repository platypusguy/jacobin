/*
	Standard (Java 21):
	https://docs.oracle.com/en/java/javase/21/docs/specs/security/standard-names.html
*/

package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"math/big"
	"strings"
)

// Initialize EC curve name
func (kpg *KeyPairGenerator) InitializeCurve(curveName string) {
	kpg.CurveName = curveName
}

// ---------------------------
// EC Curve selection
// ---------------------------
func selectCurve(kpg *KeyPairGenerator) (elliptic.Curve, error) {
	if kpg.CurveName != "" {
		switch strings.ToUpper(kpg.CurveName) {
		case "P-256", "SECP256R1":
			return elliptic.P256(), nil
		case "P-384", "SECP384R1":
			return elliptic.P384(), nil
		case "P-521", "SECP521R1":
			return elliptic.P521(), nil
		default:
			return nil, errors.New("unsupported EC curve: " + kpg.CurveName)
		}
	}

	switch kpg.KeySize {
	case 384:
		return elliptic.P384(), nil
	case 521:
		return elliptic.P521(), nil
	default:
		return elliptic.P256(), nil
	}
}

// ---------------------------
// KeyPairGenerator for EC
// ---------------------------
func (kpg *KeyPairGenerator) GenerateECKeyPair() (*KeyPair, error) {
	var curve elliptic.Curve
	switch strings.ToUpper(kpg.CurveName) {
	case "P-256", "SECP256R1":
		curve = elliptic.P256()
	case "P-384", "SECP384R1":
		curve = elliptic.P384()
	case "P-521", "SECP521R1":
		curve = elliptic.P521()
	default:
		return nil, errors.New("unsupported EC curve: " + kpg.CurveName)
	}

	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	return &KeyPair{PrivateKey: priv, PublicKey: &priv.PublicKey, Algorithm: "EC"}, nil
}

// ---------------------------
// ECDSA Signing / Verification
// ---------------------------
func (kp *KeyPair) SignECDSA(message []byte) ([]byte, error) {
	priv, ok := kp.PrivateKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not ECDSA")
	}

	// Choose hashish based on curve
	var hashish []byte
	switch priv.Curve.Params().BitSize {
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
		return nil, errors.New("unsupported ECDSA curve")
	}

	// Sign
	r, s, err := ecdsa.Sign(rand.Reader, priv, hashish)
	if err != nil {
		return nil, err
	}

	// Fixed-length serialization
	byteLen := (priv.Curve.Params().BitSize + 7) / 8
	rBytes := r.Bytes()
	sBytes := s.Bytes()

	rPadded := make([]byte, byteLen)
	sPadded := make([]byte, byteLen)
	copy(rPadded[byteLen-len(rBytes):], rBytes)
	copy(sPadded[byteLen-len(sBytes):], sBytes)

	sig := append(rPadded, sPadded...)
	return sig, nil
}

func (kp *KeyPair) VerifyECDSA(message, sig []byte) (bool, error) {
	pub, ok := kp.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return false, errors.New("public key is not ECDSA")
	}
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

	half := len(sig) / 2
	if half == 0 {
		return false, errors.New("invalid signature length")
	}
	r := new(big.Int).SetBytes(sig[:half])
	s := new(big.Int).SetBytes(sig[half:])
	valid := ecdsa.Verify(pub, hashish[:], r, s)
	return valid, nil
}

// ---------------------------
// Serialize / Deserialize EC Points
// ---------------------------
func serializePoint(curve elliptic.Curve, x, y *big.Int) []byte {
	byteLen := (curve.Params().BitSize + 7) / 8
	xb := x.Bytes()
	yb := y.Bytes()
	if len(xb) < byteLen {
		xb = append(make([]byte, byteLen-len(xb)), xb...)
	}
	if len(yb) < byteLen {
		yb = append(make([]byte, byteLen-len(yb)), yb...)
	}
	return append(xb, yb...)
}

func deserializePoint(curve elliptic.Curve, data []byte) (*big.Int, *big.Int, error) {
	byteLen := (curve.Params().BitSize + 7) / 8
	if len(data) != 2*byteLen {
		return nil, nil, errors.New("invalid ephemeral public key length")
	}
	x := new(big.Int).SetBytes(data[:byteLen])
	y := new(big.Int).SetBytes(data[byteLen:])
	return x, y, nil
}

// ---------------------------
// ECIES Encryption
// ---------------------------
func (kp *KeyPair) EncryptECIES(recipientPub *ecdsa.PublicKey, plaintext []byte) ([]byte, []byte, error) {
	// Generate ephemeral key
	ephemeral, err := ecdsa.GenerateKey(recipientPub.Curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// Compute shared secret: ephemeralPriv * recipientPub
	x, _ := recipientPub.Curve.ScalarMult(recipientPub.X, recipientPub.Y, ephemeral.D.Bytes())
	if x == nil {
		return nil, nil, errors.New("failed to compute shared secret")
	}

	// AES-256 key derived from shared secret
	key := sha256.Sum256(x.Bytes())

	// AES-GCM encryption
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, nil, err
	}
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	// Serialize ephemeral public key
	ephemeralPub := serializePoint(ephemeral.Curve, ephemeral.X, ephemeral.Y)
	return ephemeralPub, append(nonce, ciphertext...), nil
}

// ---------------------------
// ECIES Decryption
// ---------------------------
func (kp *KeyPair) DecryptECIES(ephemeralPubBytes, ciphertext []byte) ([]byte, error) {
	priv, ok := kp.PrivateKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not ECDSA")
	}

	// Deserialize ephemeral public key
	x, y, err := deserializePoint(priv.Curve, ephemeralPubBytes)
	if err != nil {
		return nil, err
	}

	// Compute shared secret
	secretX, _ := priv.Curve.ScalarMult(x, y, priv.D.Bytes())
	if secretX == nil {
		return nil, errors.New("failed to compute shared secret")
	}
	key := sha256.Sum256(secretX.Bytes())

	// AES-GCM decryption
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aesgcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}
	nonce := ciphertext[:aesgcm.NonceSize()]
	ciphertextData := ciphertext[aesgcm.NonceSize():]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertextData, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
