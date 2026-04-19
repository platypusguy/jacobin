/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

import (
	"bytes"
	"testing"
)

func TestApplyPadding(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		blockSize int
		padding   string
		expected  []byte
	}{
		{
			name:      "PKCS5Padding full block",
			input:     []byte{1, 2, 3, 4, 5, 6, 7, 8},
			blockSize: 8,
			padding:   "PKCS5Padding",
			expected:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 8, 8, 8, 8, 8, 8, 8, 8},
		},
		{
			name:      "PKCS5Padding partial block",
			input:     []byte{1, 2, 3, 4, 5},
			blockSize: 8,
			padding:   "PKCS5Padding",
			expected:  []byte{1, 2, 3, 4, 5, 3, 3, 3},
		},
		{
			name:      "NoPadding",
			input:     []byte{1, 2, 3, 4, 5},
			blockSize: 8,
			padding:   "NoPadding",
			expected:  []byte{1, 2, 3, 4, 5},
		},
		{
			name:      "Unknown padding (default)",
			input:     []byte{1, 2, 3, 4, 5},
			blockSize: 8,
			padding:   "Unknown",
			expected:  []byte{1, 2, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := applyPadding(tt.input, tt.blockSize, tt.padding)
			if !bytes.Equal(got, tt.expected) {
				t.Errorf("applyPadding() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRemovePadding(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		blockSize int
		padding   string
		expected  []byte
		wantErr   bool
	}{
		{
			name:      "PKCS5Padding full block",
			input:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 8, 8, 8, 8, 8, 8, 8, 8},
			blockSize: 8,
			padding:   "PKCS5Padding",
			expected:  []byte{1, 2, 3, 4, 5, 6, 7, 8},
			wantErr:   false,
		},
		{
			name:      "PKCS5Padding partial block",
			input:     []byte{1, 2, 3, 4, 5, 3, 3, 3},
			blockSize: 8,
			padding:   "PKCS5Padding",
			expected:  []byte{1, 2, 3, 4, 5},
			wantErr:   false,
		},
		{
			name:      "NoPadding",
			input:     []byte{1, 2, 3, 4, 5},
			blockSize: 8,
			padding:   "NoPadding",
			expected:  []byte{1, 2, 3, 4, 5},
			wantErr:   false,
		},
		{
			name:      "Empty input",
			input:     []byte{},
			blockSize: 8,
			padding:   "PKCS5Padding",
			expected:  []byte{},
			wantErr:   false,
		},
		{
			name:      "Invalid input length",
			input:     []byte{1, 2, 3, 4, 5},
			blockSize: 8,
			padding:   "PKCS5Padding",
			expected:  nil,
			wantErr:   true,
		},
		{
			name:      "Invalid padding length (0)",
			input:     []byte{1, 2, 3, 4, 5, 6, 7, 0},
			blockSize: 8,
			padding:   "PKCS5Padding",
			expected:  nil,
			wantErr:   true,
		},
		{
			name:      "Invalid padding length (too large)",
			input:     []byte{1, 2, 3, 4, 5, 6, 7, 9},
			blockSize: 8,
			padding:   "PKCS5Padding",
			expected:  nil,
			wantErr:   true,
		},
		{
			name:      "Invalid padding bytes",
			input:     []byte{1, 2, 3, 4, 5, 1, 2, 3},
			blockSize: 8,
			padding:   "PKCS5Padding",
			expected:  nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := removePadding(tt.input, tt.blockSize, tt.padding)
			if (err != nil) != tt.wantErr {
				t.Errorf("removePadding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !bytes.Equal(got, tt.expected) {
				t.Errorf("removePadding() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPerformCipherAES(t *testing.T) {
	key := []byte("1234567812345678") // 16 bytes
	iv := []byte("iviviviviviviviv")  // 16 bytes
	input := []byte("Hello Jacobin!")

	tests := []struct {
		name           string
		transformation CipherTransformation
	}{
		{
			name: "AES/CBC/PKCS5Padding",
			transformation: CipherTransformation{
				Algorithm: "AES",
				Mode:      "CBC",
				Padding:   "PKCS5Padding",
			},
		},
		{
			name: "AES/ECB/PKCS5Padding",
			transformation: CipherTransformation{
				Algorithm: "AES",
				Mode:      "ECB",
				Padding:   "PKCS5Padding",
			},
		},
		{
			name: "AES/CTR/NoPadding",
			transformation: CipherTransformation{
				Algorithm: "AES",
				Mode:      "CTR",
				Padding:   "NoPadding",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := performCipher(tt.transformation, 1, key, iv, input)
			if err != nil {
				t.Fatalf("Encryption failed: %v", err)
			}

			// Decrypt
			decrypted, err := performCipher(tt.transformation, 2, key, iv, encrypted)
			if err != nil {
				t.Fatalf("Decryption failed: %v", err)
			}

			if !bytes.Equal(decrypted, input) {
				t.Errorf("Decrypted data does not match input. Got %v, want %v", decrypted, input)
			}
		})
	}
}

func TestPerformCipherDES(t *testing.T) {
	key := []byte("12345678") // 8 bytes for DES
	iv := []byte("iviviviv")  // 8 bytes
	input := []byte("DES test")

	transformation := CipherTransformation{
		Algorithm: "DES",
		Mode:      "CBC",
		Padding:   "PKCS5Padding",
	}

	// Encrypt
	encrypted, err := performCipher(transformation, 1, key, iv, input)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Decrypt
	decrypted, err := performCipher(transformation, 2, key, iv, encrypted)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if !bytes.Equal(decrypted, input) {
		t.Errorf("Decrypted data does not match input. Got %v, want %v", decrypted, input)
	}
}

func TestPerformCipherDESede(t *testing.T) {
	key := []byte("123456781234567812345678") // 24 bytes for DESede
	iv := []byte("iviviviv")                  // 8 bytes
	input := []byte("TripleDES test")

	transformation := CipherTransformation{
		Algorithm: "DESede",
		Mode:      "CBC",
		Padding:   "PKCS5Padding",
	}

	// Encrypt
	encrypted, err := performCipher(transformation, 1, key, iv, input)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Decrypt
	decrypted, err := performCipher(transformation, 2, key, iv, encrypted)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if !bytes.Equal(decrypted, input) {
		t.Errorf("Decrypted data does not match input. Got %v, want %v", decrypted, input)
	}
}

func TestPerformCipherErrors(t *testing.T) {
	key := []byte("1234567812345678")
	iv := []byte("iviviviviviviviv")

	tests := []struct {
		name           string
		transformation CipherTransformation
		key            []byte
		iv             []byte
		input          []byte
		opmode         int64
		wantErr        bool
	}{
		{
			name: "Unsupported algorithm",
			transformation: CipherTransformation{
				Algorithm: "UNSUPPORTED",
			},
			key:     key,
			iv:      iv,
			input:   []byte("test"),
			opmode:  1,
			wantErr: true,
		},
		{
			name: "Invalid key size for AES",
			transformation: CipherTransformation{
				Algorithm: "AES",
			},
			key:     []byte("short"),
			iv:      iv,
			input:   []byte("test"),
			opmode:  1,
			wantErr: true,
		},
		{
			name: "Unsupported mode",
			transformation: CipherTransformation{
				Algorithm: "AES",
				Mode:      "UNSUPPORTED",
			},
			key:     key,
			iv:      iv,
			input:   []byte("test"),
			opmode:  1,
			wantErr: true,
		},
		{
			name: "Invalid IV length for CBC",
			transformation: CipherTransformation{
				Algorithm: "AES",
				Mode:      "CBC",
			},
			key:     key,
			iv:      []byte("short"),
			input:   []byte("test"),
			opmode:  1,
			wantErr: true,
		},
		{
			name: "Invalid input length for ECB decryption",
			transformation: CipherTransformation{
				Algorithm: "AES",
				Mode:      "ECB",
				Padding:   "NoPadding",
			},
			key:     key,
			iv:      iv,
			input:   []byte("not a block"),
			opmode:  2,
			wantErr: true,
		},
		{
			name: "Invalid input length for CBC decryption",
			transformation: CipherTransformation{
				Algorithm: "AES",
				Mode:      "CBC",
				Padding:   "NoPadding",
			},
			key:     key,
			iv:      iv,
			input:   []byte("not a block"),
			opmode:  2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := performCipher(tt.transformation, tt.opmode, tt.key, tt.iv, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("performCipher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
