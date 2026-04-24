package javaxCrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rc4"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/blowfish"
	"golang.org/x/crypto/chacha20"
	"golang.org/x/crypto/chacha20poly1305"
)

func performCipher(config CipherTransformation, opmode int64, key []byte, iv []byte, input []byte) ([]byte, error) {
	var block cipher.Block
	var err error

	switch config.Algorithm {
	case "AES", "AESWrap", "AESWrapPad":
		block, err = aes.NewCipher(key)
	case "DES", "PBEWithMD5AndDES":
		block, err = des.NewCipher(key)
	case "DESede", "DESedeWrap", "PBEWithMD5AndTripleDES", "PBEWithSHA1AndDESede":
		block, err = des.NewTripleDESCipher(key)
	case "Blowfish":
		block, err = blowfish.NewCipher(key)
	case "RC2", "RC2Wrap", "PBEWithSHA1AndRC2_40", "PBEWithSHA1AndRC2_128":
		bits := len(key) * 8
		if strings.Contains(config.Name, "40") {
			bits = 40
		} else if strings.Contains(config.Name, "128") {
			bits = 128
		}
		block, err = newRC2Cipher(key, bits)
	case "ARCFOUR", "RC4":
		// These are stream ciphers, handled below in the mode switch
		block = nil
	case "":
		return nil, fmt.Errorf("performCipher: nil algorithm")
	default:
		return nil, fmt.Errorf("performCipher: unsupported algorithm: %s", config.Algorithm)
	}

	if err != nil {
		return nil, err
	}

	isEncrypt := opmode == 1 // ENCRYPT_MODE

	var result []byte
	switch config.Mode {
	case "ECB":
		if isEncrypt {
			input = applyPadding(input, block.BlockSize(), config.Padding)
			result = make([]byte, len(input))
			for i := 0; i < len(input); i += block.BlockSize() {
				block.Encrypt(result[i:i+block.BlockSize()], input[i:i+block.BlockSize()])
			}
		} else {
			if len(input)%block.BlockSize() != 0 {
				return nil, errors.New("input length must be multiple of block size for ECB decryption")
			}
			result = make([]byte, len(input))
			for i := 0; i < len(input); i += block.BlockSize() {
				block.Decrypt(result[i:i+block.BlockSize()], input[i:i+block.BlockSize()])
			}
			result, err = removePadding(result, block.BlockSize(), config.Padding)
			if err != nil {
				return nil, err
			}
		}
	case "CBC":
		if len(iv) == 0 && strings.Contains(config.Name, "PBEWith") {
			// Legacy PBE usually defaults to CBC
			if !strings.Contains(config.Name, "Hmac") && !strings.Contains(config.Name, "HMAC") {
				// For legacy PBE, if IV is not provided, extract it from key material
				if strings.Contains(config.Name, "DESede") || strings.Contains(config.Name, "TripleDES") {
					if len(key) >= 32 {
						iv = key[24:32]
					}
				} else if strings.Contains(config.Name, "DES") {
					if len(key) >= 16 {
						iv = key[8:16]
					}
				} else if strings.Contains(config.Name, "RC2") {
					bits := 128
					if strings.Contains(config.Name, "40") {
						bits = 40
					}
					if len(key) >= (bits/8)+8 {
						iv = key[bits/8 : (bits/8)+8]
					}
				}
			}
			if len(iv) == 0 {
				iv = make([]byte, block.BlockSize())
			}
		}

		if len(iv) != block.BlockSize() {
			return nil, fmt.Errorf("invalid IV length: expected %d, got %d", block.BlockSize(), len(iv))
		}
		if isEncrypt {
			input = applyPadding(input, block.BlockSize(), config.Padding)
			result = make([]byte, len(input))
			mode := cipher.NewCBCEncrypter(block, iv)
			mode.CryptBlocks(result, input)
		} else {
			if len(input)%block.BlockSize() != 0 {
				return nil, errors.New("performCipher: input length must be multiple of block size for CBC decryption")
			}
			result = make([]byte, len(input))
			mode := cipher.NewCBCDecrypter(block, iv)
			mode.CryptBlocks(result, input)
			result, err = removePadding(result, block.BlockSize(), config.Padding)
			if err != nil {
				return nil, err
			}
		}
	case "CTR":
		if len(iv) != block.BlockSize() {
			return nil, fmt.Errorf("performCipher: invalid IV length: expected %d, got %d", block.BlockSize(), len(iv))
		}
		result = make([]byte, len(input))
		mode := cipher.NewCTR(block, iv)
		mode.XORKeyStream(result, input)
		// No padding for CTR
	case "GCM":
		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		if len(iv) != aesgcm.NonceSize() {
			return nil, fmt.Errorf("performCipher: invalid nonce length: expected %d, got %d", aesgcm.NonceSize(), len(iv))
		}
		if isEncrypt {
			result = aesgcm.Seal(nil, iv, input, nil)
		} else {
			result, err = aesgcm.Open(nil, iv, input, nil)
			if err != nil {
				return nil, err
			}
		}
	case "ChaCha20":
		// ChaCha20 uses a 96-bit nonce (12 bytes) and a 32-bit counter
		if len(iv) != 12 {
			return nil, fmt.Errorf("performCipher: invalid IV length for ChaCha20: expected 12, got %d", len(iv))
		}
		result = make([]byte, len(input))
		c, err := chacha20.NewUnauthenticatedCipher(key, iv)
		if err != nil {
			return nil, err
		}
		c.XORKeyStream(result, input)
	case "ChaCha20-Poly1305":
		aead, err := chacha20poly1305.New(key)
		if err != nil {
			return nil, err
		}
		if len(iv) != aead.NonceSize() {
			return nil, fmt.Errorf("performCipher: invalid nonce length for ChaCha20-Poly1305: expected %d, got %d", aead.NonceSize(), len(iv))
		}
		if isEncrypt {
			result = aead.Seal(nil, iv, input, nil)
		} else {
			result, err = aead.Open(nil, iv, input, nil)
			if err != nil {
				return nil, err
			}
		}
	case "RC4", "ARCFOUR":
		c, err := rc4.NewCipher(key)
		if err != nil {
			return nil, err
		}
		result = make([]byte, len(input))
		c.XORKeyStream(result, input)
		return result, nil
	case "AESWrap", "AESWrapPad", "DESedeWrap":
		// Key wrap algorithms.
		// For now, let's treat them as ECB for testing purposes if they are simple,
		// but properly they need RFC 3394/5649.
		// Go's crypto doesn't have a standard keywrap yet.
		return nil, fmt.Errorf("performCipher: %s not yet supported", config.Mode)
	case "":
		// Some algorithms (like PBE) might not specify mode in the transformation string
		// but have a default.
		if strings.Contains(config.Name, "PBEWith") {
			// Legacy PBE usually defaults to CBC
			if len(iv) == 0 {
				if !strings.Contains(config.Name, "Hmac") && !strings.Contains(config.Name, "HMAC") {
					// For legacy PBE, if IV is not provided, extract it from key material
					if strings.Contains(config.Name, "DESede") || strings.Contains(config.Name, "TripleDES") {
						if len(key) >= 32 {
							iv = key[24:32]
						}
					} else if strings.Contains(config.Name, "DES") {
						if len(key) >= 16 {
							iv = key[8:16]
						}
					} else if strings.Contains(config.Name, "RC2") {
						bits := 128
						if strings.Contains(config.Name, "40") {
							bits = 40
						}
						if len(key) >= (bits/8)+8 {
							iv = key[bits/8 : (bits/8)+8]
						}
					}
				}
				if len(iv) == 0 {
					iv = make([]byte, block.BlockSize())
				}
			}
			if len(iv) != block.BlockSize() {
				return nil, fmt.Errorf("invalid IV length for PBE: expected %d, got %d", block.BlockSize(), len(iv))
			}
			if isEncrypt {
				input = applyPadding(input, block.BlockSize(), "PKCS5Padding")
				result = make([]byte, len(input))
				mode := cipher.NewCBCEncrypter(block, iv)
				mode.CryptBlocks(result, input)
			} else {
				if len(input)%block.BlockSize() != 0 {
					return nil, errors.New("performCipher: input length must be multiple of block size for CBC decryption")
				}
				result = make([]byte, len(input))
				mode := cipher.NewCBCDecrypter(block, iv)
				mode.CryptBlocks(result, input)
				result, err = removePadding(result, block.BlockSize(), "PKCS5Padding")
				if err != nil {
					return nil, err
				}
			}
			return result, nil
		}
		return nil, fmt.Errorf("performCipher: nil mode")
	default:
		return nil, fmt.Errorf("performCipher: unsupported mode: %s", config.Mode)
	}

	return result, nil
}

func applyPadding(input []byte, blockSize int, padding string) []byte {
	switch padding {
	case "PKCS5Padding", "PKCS7Padding":
		padLen := blockSize - (len(input) % blockSize)
		padding := make([]byte, padLen)
		for i := range padding {
			padding[i] = byte(padLen)
		}
		return append(input, padding...)
	case "NoPadding":
		return input
	default:
		return input
	}
}

func removePadding(input []byte, blockSize int, padding string) ([]byte, error) {
	if len(input) == 0 {
		return input, nil
	}
	switch padding {
	case "PKCS5Padding", "PKCS7Padding":
		if len(input)%blockSize != 0 {
			return nil, errors.New("invalid input length for padding removal")
		}
		padLen := int(input[len(input)-1])
		if padLen <= 0 || padLen > blockSize {
			return nil, errors.New("invalid padding length")
		}
		for i := len(input) - padLen; i < len(input); i++ {
			if input[i] != byte(padLen) {
				return nil, errors.New("invalid padding bytes")
			}
		}
		return input[:len(input)-padLen], nil
	case "NoPadding":
		return input, nil
	default:
		return input, nil
	}
}
