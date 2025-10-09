package security

import (
	"bytes"
	"crypto/ecdsa"
	"testing"
)

/*
 * Given the algorithm and key size, generated and return a KeyPair struct.
 */
func setUpKeyPairs(t *testing.T, algo string, keySize int, curveName string) *KeyPair {
	kpgen, err := NewKeyPairGenerator(algo)
	if err != nil {
		t.Fatalf("*** NewKeyPairGenerator(%s) failed, err: %s", algo, err.Error())
	}
	err = kpgen.Initialize(keySize)
	if err != nil {
		t.Fatalf("*** kpgen.Initialize(%s-%d) failed, err: %s", algo, keySize, err)
	}
	kpgen.InitializeCurve(curveName)
	keyPair, err := kpgen.GenerateKeyPair()
	if err != nil {
		t.Fatalf("*** kpgen.GenerateKeyPair(%s-%d) failed, err: %s", algo, keySize, err)
	}

	return keyPair
}

func TestRsaCrypto(t *testing.T) {
	plaintext := []byte("Hello RSA encryption!")

	// --- RSA-2048 key pair ---
	rsaKP := setUpKeyPairs(t, "RSA", 2048, "dummy")

	// RSAEncrypt/RSADecrypt using RSA PKCS1
	ciphertext, err := rsaKP.RSAEncrypt(plaintext, "RSA")
	if err != nil {
		t.Fatal("*** rsaKP.RSAEncrypt(RSA) failed, err: " + err.Error())
	}
	decrypted, err := rsaKP.RSADecrypt(ciphertext, "RSA")
	if err != nil {
		t.Fatal("*** rsaKP.RSADecrypt(RSA) failed, err: " + err.Error())
	}

	if string(plaintext) != string(decrypted) {
		t.Error("*** RSA-2048 plaintext: ", string(plaintext))
		t.Error("*** RSA-2048 decrypted: ", string(decrypted), " :: no match!")
	}
}

func TestDiffieHillmanCrypto(t *testing.T) {
	params, err := NewDHParameters(2048)
	if err != nil {
		t.Fatal("*** NewDHParameters(2048) failed, err: " + err.Error())
	}
	alice, err := GenerateDHKeyPair(params)
	if err != nil {
		t.Fatal("*** GenerateDHKeyPair(alice) failed, err: " + err.Error())
	}
	bob, err := GenerateDHKeyPair(params)
	if err != nil {
		t.Fatal("*** GenerateDHKeyPair(bob) failed, err: " + err.Error())
	}

	secretAlice, err := alice.ComputeShared(bob.PublicKey)
	if err != nil {
		t.Fatal("*** alice.ComputeShared(bob.PublicKey) failed, err: " + err.Error())
	}
	secretBob, err := bob.ComputeShared(alice.PublicKey)
	if err != nil {
		t.Fatal("*** bob.ComputeShared(alice.PublicKey) failed, err: " + err.Error())
	}

	if secretAlice.Cmp(secretBob) != 0 {
		t.Error("*** DH secretAlice.Cmp(secretBob) != 0 :: no match!")
	}

	// Convert secret to 32-byte (256 bits) AES key
	key := secretAlice.Bytes()
	if len(key) < 32 {
		padding := make([]byte, 32-len(key))
		key = append(padding, key...)
	} else if len(key) > 32 {
		key = key[:32]
	}

	plaintext := []byte("Hello AES-GCM encryption!")

	// RSAEncrypt/RSADecrypt via AES-GCM
	nonce, ciphertext, err := EncryptWithAESGCM(key, plaintext)
	if err != nil {
		t.Fatal("*** EncryptWithAESGCM failed, err: " + err.Error())
	}
	decrypted, err := DecryptWithAESGCM(key, nonce, ciphertext)
	if err != nil {
		t.Fatal("*** DecryptWithAESGCM failed, err: " + err.Error())
	}

	if string(plaintext) != string(decrypted) {
		t.Error("*** AES-GCM plaintext: ", string(plaintext))
		t.Error("*** AES-GCM decrypted: ", string(decrypted), " :: no match!")
	}
}

func TestSignVerifyRSA(t *testing.T) {
	message := []byte("Test message for signing")

	// --- RSA Signing ---
	rsaKP := setUpKeyPairs(t, "RSA", 2048, "dummy")

	signatureRSA, err := rsaKP.SignWithAlgorithm(message, "SHA256WITHRSA")
	if err != nil {
		t.Fatal("*** rsaKP.SignWithAlgorithm failed:", err)
	}
	success, err := rsaKP.VerifyWithAlgorithm(message, signatureRSA, "SHA256WITHRSA")
	if err != nil {
		t.Fatal("*** rsaKP.VerifyWithAlgorithm failed:", err)
	}
	if !success {
		t.Error("*** rsaKP.VerifyWithAlgorithm failed to match the message")
	}
}

func TestSignVerifyDSA(t *testing.T) {
	message := []byte("Test message for signing")

	dsaKP := setUpKeyPairs(t, "DSA", 1024, "dummy")

	signatureDSA, err := dsaKP.SignWithAlgorithm(message, "SHA256WITHDSA")
	if err != nil {
		t.Fatal("*** dsaKP.SignWithAlgorithm failed, err: ", err)
	}
	success, err := dsaKP.VerifyWithAlgorithm(message, signatureDSA, "SHA256WITHDSA")
	if err != nil {
		t.Fatal("*** dsaKP.VerifyWithAlgorithm failed, err: ", err)
	}
	if !success {
		t.Error("*** dsaKP.VerifyWithAlgorithm failed to match message")
	}
}

func TestSignVerifyEd25519(t *testing.T) {
	message := []byte("Test message for signing")
	edKP := setUpKeyPairs(t, "Ed25519", 2048, "dummy")

	signatureEd, err := edKP.SignWithAlgorithm(message, "Ed25519")
	if err != nil {
		t.Fatal("*** edKP.SignWithAlgorithm failed, err: ", err)
	}
	success, err := edKP.VerifyWithAlgorithm(message, signatureEd, "Ed25519")
	if err != nil {
		t.Fatal("*** edKP.VerifyWithAlgorithm failed, err: ", err.Error())
	}
	if !success {
		t.Error("*** edKP.VerifyWithAlgorithm failed to match message")
	}
}

func TestECTableDriven(t *testing.T) {
	curves := []string{"P-256", "P-384", "P-521"}
	messages := [][]byte{
		[]byte("Short message"),
		[]byte("Mary had a little lamb whose fleece was white as snow"),
		bytes.Repeat([]byte("a whole lotta bytes!"), 1000), // 20,000 bytes
	}

	for _, curve := range curves {
		kp := setUpKeyPairs(t, "EC", 2048, curve)
		for idx, msg := range messages {
			// --- ECDSA Sign/Verify ---
			sig, err := kp.SignECDSA(msg)
			if err != nil {
				t.Fatalf("*** [%s][msg %d] SignECDSA failed: %v", curve, idx, err)
			}
			valid, err := kp.VerifyECDSA(msg, sig)
			if err != nil {
				t.Fatalf("*** [%s][msg %d] VerifyECDSA failed: %v", curve, idx, err)
			}
			if !valid {
				t.Errorf("*** [%s][msg %d] VerifyECDSA failed to verify signature", curve, idx)
			}

			// --- ECIES RSAEncrypt/RSADecrypt ---
			// Generate a recipient key for encryption test
			kpg, err := NewKeyPairGenerator("EC")
			if err != nil {
				t.Fatalf("*** NewKeyPairGenerator(\"EC\", recipient) failed: %s", err.Error())
			}
			kpg.InitializeCurve(curve)
			kpg.Initialize(2048)
			recipientKP, _ := kpg.GenerateECKeyPair()
			ephemeralPub, ciphertext, err := kp.EncryptECIES(recipientKP.PublicKey.(*ecdsa.PublicKey), msg)
			if err != nil {
				t.Fatalf("*** [%s][msg %d] ECIES RSAEncrypt failed: %v", curve, idx, err)
			}

			plaintext, err := recipientKP.DecryptECIES(ephemeralPub, ciphertext)
			if err != nil {
				t.Fatalf("*** [%s][msg %d] ECIES RSADecrypt failed: %v", curve, idx, err)
			}

			if !bytes.Equal(msg, plaintext) {
				t.Errorf("*** [%s][msg %d] ECIES decryption mismatch\nGot: %d bytes\nWant: %d bytes",
					curve, idx, len(plaintext), len(msg))
			}

			// --- Negative test: wrong key ---
			otherKP, _ := kpg.GenerateECKeyPair()
			_, err = otherKP.DecryptECIES(ephemeralPub, ciphertext)
			if err == nil {
				t.Errorf("*** [%s][msg %d] ECIES decryption should fail with wrong private key", curve, idx)
			}
		}
	}
}
