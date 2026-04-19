package javaxCrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"errors"
	"fmt"
)

func performCipher(config CipherTransformation, opmode int64, key []byte, iv []byte, input []byte) ([]byte, error) {
	var block cipher.Block
	var err error

	switch config.Algorithm {
	case "AES":
		block, err = aes.NewCipher(key)
	case "DES":
		block, err = des.NewCipher(key)
	case "DESede":
		block, err = des.NewTripleDESCipher(key)
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", config.Algorithm)
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
				return nil, errors.New("input length must be multiple of block size for CBC decryption")
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
			return nil, fmt.Errorf("invalid IV length: expected %d, got %d", block.BlockSize(), len(iv))
		}
		result = make([]byte, len(input))
		mode := cipher.NewCTR(block, iv)
		mode.XORKeyStream(result, input)
		// No padding for CTR
	default:
		return nil, fmt.Errorf("unsupported mode: %s", config.Mode)
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
