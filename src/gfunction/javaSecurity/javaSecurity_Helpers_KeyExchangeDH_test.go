package javaSecurity

import (
	"bytes"
	"math/big"
	"testing"

	"jacobin/src/object"
)

func TestKeyExchangeDH_Roundtrip(t *testing.T) {
	// Standard DH parameters (Group 14 - 2048-bit MODP Group)
	// For testing, we can use smaller ones for speed, but let's use something reasonable.
	p, _ := new(big.Int).SetString("FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1"+
		"29024E088A67CC74020BBEA63B139B22514A08798E3404DD"+
		"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245"+
		"E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED"+
		"EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3D"+
		"C2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F"+
		"83655D23DCA3AD961C62F356208552BB9ED529077096966D"+
		"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3B"+
		"E39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9"+
		"DE2BCBF6955817183995497CEA956AE515D2261898FA0510"+
		"15728E5A8AACAA68FFFFFFFFFFFFFFFF", 16)
	g := big.NewInt(2)

	// Recipient's Private Key
	privValue := big.NewInt(123456789)
	
	paramsObj := object.MakeEmptyObject()
	paramsObj.FieldTable["p"] = object.Field{Fvalue: p}
	paramsObj.FieldTable["g"] = object.Field{Fvalue: g}

	privKeyObj := object.MakeEmptyObject()
	privKeyObj.FieldTable["value"] = object.Field{Fvalue: privValue}
	privKeyObj.FieldTable["params"] = object.Field{Fvalue: paramsObj}

	// Recipient's Public Key: g^priv mod p
	pubValue := new(big.Int).Exp(g, privValue, p)
	pubKeyObj := object.MakeEmptyObject()
	pubKeyObj.FieldTable["value"] = object.Field{Fvalue: pubValue}
	pubKeyObj.FieldTable["params"] = object.Field{Fvalue: paramsObj}

	sharedSecret := []byte("this is a shared secret")

	// Encrypt
	ciphertext, err := keyExchangeEncryptDH(pubKeyObj, sharedSecret)
	if err != nil {
		t.Fatalf("keyExchangeEncryptDH failed: %v", err)
	}

	// Decrypt
	decrypted, err := keyExchangeDecryptDH(privKeyObj, ciphertext)
	if err != nil {
		t.Fatalf("keyExchangeDecryptDH failed: %v", err)
	}

	if !bytes.Equal(sharedSecret, decrypted) {
		t.Errorf("Decrypted secret does not match. Expected %v, got %v", sharedSecret, decrypted)
	}
}

func TestKeyExchangeDH_InvalidPublicKey(t *testing.T) {
	p := big.NewInt(23)
	g := big.NewInt(5)
	
	paramsObj := object.MakeEmptyObject()
	paramsObj.FieldTable["p"] = object.Field{Fvalue: p}
	paramsObj.FieldTable["g"] = object.Field{Fvalue: g}

	// Weak public key: 1
	pubKeyObj := object.MakeEmptyObject()
	pubKeyObj.FieldTable["value"] = object.Field{Fvalue: big.NewInt(1)}
	pubKeyObj.FieldTable["params"] = object.Field{Fvalue: paramsObj}

	sharedSecret := []byte("secret")
	_, err := keyExchangeEncryptDH(pubKeyObj, sharedSecret)
	if err == nil {
		t.Error("Expected error for weak public key (1), got nil")
	}

	// Weak public key: p-1
	pubKeyObj.FieldTable["value"] = object.Field{Fvalue: big.NewInt(22)}
	_, err = keyExchangeEncryptDH(pubKeyObj, sharedSecret)
	if err == nil {
		t.Error("Expected error for weak public key (p-1), got nil")
	}

	// Out of range: p
	pubKeyObj.FieldTable["value"] = object.Field{Fvalue: big.NewInt(23)}
	_, err = keyExchangeEncryptDH(pubKeyObj, sharedSecret)
	if err == nil {
		t.Error("Expected error for out of range public key (p), got nil")
	}
}

func TestKeyExchangeDH_InvalidCiphertext(t *testing.T) {
	p := big.NewInt(23)
	privKeyObj := object.MakeEmptyObject()
	privKeyObj.FieldTable["value"] = object.Field{Fvalue: big.NewInt(5)}
	paramsObj := object.MakeEmptyObject()
	paramsObj.FieldTable["p"] = object.Field{Fvalue: p}
	privKeyObj.FieldTable["params"] = object.Field{Fvalue: paramsObj}

	// Too short
	_, err := keyExchangeDecryptDH(privKeyObj, []byte{0, 0, 0})
	if err == nil {
		t.Error("Expected error for too short ciphertext, got nil")
	}

	// Wrong length for ephemeral key
	_, err = keyExchangeDecryptDH(privKeyObj, []byte{0, 0, 0, 10})
	if err == nil {
		t.Error("Expected error for inconsistent ciphertext length, got nil")
	}
}
