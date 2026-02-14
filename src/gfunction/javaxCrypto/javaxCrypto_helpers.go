/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2026 by the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package javaxCrypto

// CipherTransformation entry definition
type CipherTransformation struct {
	Name            string
	Enabled         bool
	Category        TransformationCategory
	Algorithm       string // For triples: the base algorithm (AES, DES, etc.)
	Mode            string // For triples: ECB, CBC, GCM, etc.
	Padding         string // For triples: NoPadding, PKCS5Padding, etc.
	NeedsIV         bool
	IVLength        int // in bytes, 0 if variable or not applicable
	NeedsTagLength  bool
	TagLength       int // in bits, 0 if variable or not applicable
	NeedsSalt       bool
	NeedsIterations bool
	KeyDerivation   bool // true for PBE algorithms
	IsAEAD          bool // Authenticated Encryption with Associated Data
	Notes           string
}

type TransformationCategory int

const (
	CategorySelfContained TransformationCategory = iota // Category 1: Algorithm name that does not appear as a triple
	CategoryTriple                                      // Category 2: Algorithm name form: Algorithm/Mode/Padding
)

/*
**
CipherConfigTable contains all valid Java cipher transformations and includes:
- All self-contained transformations (Category 1)
- All Algorithm/Mode/Padding triples (Category 2) for AES, DES, DESede, Blowfish, RC2, RC5, and RSA
- Detailed metadata about each transformation's requirements
- Helper functions for validation and parameter discovery
- Notes about security considerations and usage
**
*/
var CipherConfigTable = map[string]CipherTransformation{
	// ============================================================================
	// CATEGORY 1: SELF-CONTAINED TRANSFORMATION STRINGS
	// ============================================================================

	// Stream Ciphers
	"ChaCha20": {
		Name:           "ChaCha20",
		Enabled:        true,
		Category:       CategorySelfContained,
		NeedsIV:        true,
		IVLength:       12, // 96-bit nonce
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Stream cipher, uses nonce/counter via ChaCha20ParameterSpec or IvParameterSpec",
	},
	"ChaCha20-Poly1305": {
		Name:           "ChaCha20-Poly1305",
		Enabled:        true,
		Category:       CategorySelfContained,
		NeedsIV:        true,
		IVLength:       12, // 96-bit nonce
		NeedsTagLength: false,
		IsAEAD:         true,
		Notes:          "AEAD stream cipher, authentication tag handled automatically",
	},
	"RC4": {
		Name:           "RC4",
		Enabled:        true,
		Category:       CategorySelfContained,
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Stream cipher, deprecated due to vulnerabilities",
	},
	"ARCFOUR": {
		Name:           "ARCFOUR",
		Enabled:        true,
		Category:       CategorySelfContained,
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Alias for RC4, deprecated",
	},

	// Key Wrap Algorithms
	"AESWrap": {
		Name:           "AESWrap",
		Enabled:        true,
		Category:       CategorySelfContained,
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "RFC 3394 AES Key Wrap, used for wrapping cryptographic keys",
	},
	"AESWrapPad": {
		Name:           "AESWrapPad",
		Enabled:        true,
		Category:       CategorySelfContained,
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "RFC 5649 AES Key Wrap with Padding",
	},
	"DESedeWrap": {
		Name:           "DESedeWrap",
		Enabled:        true,
		Category:       CategorySelfContained,
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Triple DES Key Wrap",
	},
	"RC2Wrap": {
		Name:           "RC2Wrap",
		Enabled:        true,
		Category:       CategorySelfContained,
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "RC2 Key Wrap",
	},

	// Password-Based Encryption (PBE) Algorithms
	"PBEWithMD5AndDES": {
		Name:            "PBEWithMD5AndDES",
		Enabled:         true,
		Category:        CategorySelfContained,
		NeedsSalt:       true,
		NeedsIterations: true,
		KeyDerivation:   true,
		IsAEAD:          false,
		Notes:           "Legacy PBE, uses PBEParameterSpec(salt, iterations), deprecated",
	},
	"PBEWithMD5AndTripleDES": {
		Name:            "PBEWithMD5AndTripleDES",
		Enabled:         true,
		Category:        CategorySelfContained,
		NeedsSalt:       true,
		NeedsIterations: true,
		KeyDerivation:   true,
		IsAEAD:          false,
		Notes:           "Legacy PBE with Triple DES",
	},
	"PBEWithSHA1AndDESede": {
		Name:            "PBEWithSHA1AndDESede",
		Enabled:         true,
		Category:        CategorySelfContained,
		NeedsSalt:       true,
		NeedsIterations: true,
		KeyDerivation:   true,
		IsAEAD:          false,
		Notes:           "PBE with SHA1 and Triple DES",
	},
	"PBEWithSHA1AndRC2_40": {
		Name:            "PBEWithSHA1AndRC2_40",
		Enabled:         true,
		Category:        CategorySelfContained,
		NeedsSalt:       true,
		NeedsIterations: true,
		KeyDerivation:   true,
		IsAEAD:          false,
		Notes:           "PBE with SHA1 and 40-bit RC2",
	},
	"PBEWithSHA1AndRC2_128": {
		Name:            "PBEWithSHA1AndRC2_128",
		Enabled:         true,
		Category:        CategorySelfContained,
		NeedsSalt:       true,
		NeedsIterations: true,
		KeyDerivation:   true,
		IsAEAD:          false,
		Notes:           "PBE with SHA1 and 128-bit RC2",
	},
	"PBEWithSHA1AndRC4_40": {
		Name:            "PBEWithSHA1AndRC4_40",
		Enabled:         true,
		Category:        CategorySelfContained,
		NeedsSalt:       true,
		NeedsIterations: true,
		KeyDerivation:   true,
		IsAEAD:          false,
		Notes:           "PBE with SHA1 and 40-bit RC4",
	},
	"PBEWithSHA1AndRC4_128": {
		Name:            "PBEWithSHA1AndRC4_128",
		Enabled:         true,
		Category:        CategorySelfContained,
		NeedsSalt:       true,
		NeedsIterations: true,
		KeyDerivation:   true,
		IsAEAD:          false,
		Notes:           "PBE with SHA1 and 128-bit RC4",
	},
	"PBEWithHmacSHA1AndAES_128": {
		Name:            "PBEWithHmacSHA1AndAES_128",
		Enabled:         true,
		Category:        CategorySelfContained,
		NeedsSalt:       true,
		NeedsIterations: true,
		KeyDerivation:   true,
		IsAEAD:          false,
		Notes:           "PBKDF2 with HMAC-SHA1 and AES-128",
	},
	"PBEWithHmacSHA224AndAES_128": {
		Name:            "PBEWithHmacSHA224AndAES_128",
		Enabled:         true,
		Category:        CategorySelfContained,
		NeedsSalt:       true,
		NeedsIterations: true,
		KeyDerivation:   true,
		IsAEAD:          false,
		Notes:           "PBKDF2 with HMAC-SHA224 and AES-128",
	},
	"PBEWithHmacSHA256AndAES_128": {
		Name:            "PBEWithHmacSHA256AndAES_128",
		Enabled:         true,
		Category:        CategorySelfContained,
		NeedsSalt:       true,
		NeedsIterations: true,
		KeyDerivation:   true,
		IsAEAD:          false,
		Notes:           "PBKDF2 with HMAC-SHA256 and AES-128",
	},
	"PBEWithHmacSHA384AndAES_128": {
		Name:            "PBEWithHmacSHA384AndAES_128",
		Enabled:         true,
		Category:        CategorySelfContained,
		NeedsSalt:       true,
		NeedsIterations: true,
		KeyDerivation:   true,
		IsAEAD:          false,
		Notes:           "PBKDF2 with HMAC-SHA384 and AES-128",
	},
	"PBEWithHmacSHA512AndAES_128": {
		Name:            "PBEWithHmacSHA512AndAES_128",
		Enabled:         true,
		Category:        CategorySelfContained,
		NeedsSalt:       true,
		NeedsIterations: true,
		KeyDerivation:   true,
		IsAEAD:          false,
		Notes:           "PBKDF2 with HMAC-SHA512 and AES-128",
	},
	"PBEWithHmacSHA1AndAES_256": {
		Name:            "PBEWithHmacSHA1AndAES_256",
		Enabled:         true,
		Category:        CategorySelfContained,
		NeedsSalt:       true,
		NeedsIterations: true,
		KeyDerivation:   true,
		IsAEAD:          false,
		Notes:           "PBKDF2 with HMAC-SHA1 and AES-256",
	},
	"PBEWithHmacSHA224AndAES_256": {
		Name:            "PBEWithHmacSHA224AndAES_256",
		Enabled:         true,
		Category:        CategorySelfContained,
		NeedsSalt:       true,
		NeedsIterations: true,
		KeyDerivation:   true,
		IsAEAD:          false,
		Notes:           "PBKDF2 with HMAC-SHA224 and AES-256",
	},
	"PBEWithHmacSHA256AndAES_256": {
		Name:            "PBEWithHmacSHA256AndAES_256",
		Enabled:         true,
		Category:        CategorySelfContained,
		NeedsSalt:       true,
		NeedsIterations: true,
		KeyDerivation:   true,
		IsAEAD:          false,
		Notes:           "PBKDF2 with HMAC-SHA256 and AES-256",
	},
	"PBEWithHmacSHA384AndAES_256": {
		Name:            "PBEWithHmacSHA384AndAES_256",
		Enabled:         true,
		Category:        CategorySelfContained,
		NeedsSalt:       true,
		NeedsIterations: true,
		KeyDerivation:   true,
		IsAEAD:          false,
		Notes:           "PBKDF2 with HMAC-SHA384 and AES-256",
	},
	"PBEWithHmacSHA512AndAES_256": {
		Name:            "PBEWithHmacSHA512AndAES_256",
		Enabled:         true,
		Category:        CategorySelfContained,
		NeedsSalt:       true,
		NeedsIterations: true,
		KeyDerivation:   true,
		IsAEAD:          false,
		Notes:           "PBKDF2 with HMAC-SHA512 and AES-256",
	},

	// Other
	"ECIES": {
		Name:           "ECIES",
		Enabled:        true,
		Category:       CategorySelfContained,
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Elliptic Curve Integrated Encryption Scheme",
	},

	// ============================================================================
	// CATEGORY 2: ALGORITHM/MODE/PADDING TRIPLES - AES
	// ============================================================================

	"AES/ECB/NoPadding": {
		Name:           "AES/ECB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "AES",
		Mode:           "ECB",
		Padding:        "NoPadding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Insecure: ECB mode, no IV. Input must be multiple of 16 bytes.",
	},
	"AES/ECB/PKCS5Padding": {
		Name:           "AES/ECB/PKCS5Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "AES",
		Mode:           "ECB",
		Padding:        "PKCS5Padding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Insecure: ECB mode, no IV. Avoid in production.",
	},
	"AES/ECB/ISO10126Padding": {
		Name:           "AES/ECB/ISO10126Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "AES",
		Mode:           "ECB",
		Padding:        "ISO10126Padding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Insecure: ECB mode, no IV",
	},
	"AES/CBC/NoPadding": {
		Name:           "AES/CBC/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "AES",
		Mode:           "CBC",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       16,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec. Input must be multiple of 16 bytes.",
	},
	"AES/CBC/PKCS5Padding": {
		Name:           "AES/CBC/PKCS5Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "AES",
		Mode:           "CBC",
		Padding:        "PKCS5Padding",
		NeedsIV:        true,
		IVLength:       16,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec. Common choice but not authenticated.",
	},
	"AES/CBC/ISO10126Padding": {
		Name:           "AES/CBC/ISO10126Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "AES",
		Mode:           "CBC",
		Padding:        "ISO10126Padding",
		NeedsIV:        true,
		IVLength:       16,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec",
	},
	"AES/CFB/NoPadding": {
		Name:           "AES/CFB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "AES",
		Mode:           "CFB",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       16,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec. Stream mode, padding not needed.",
	},
	"AES/CFB/PKCS5Padding": {
		Name:           "AES/CFB/PKCS5Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "AES",
		Mode:           "CFB",
		Padding:        "PKCS5Padding",
		NeedsIV:        true,
		IVLength:       16,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec. Padding unnecessary for CFB.",
	},
	"AES/OFB/NoPadding": {
		Name:           "AES/OFB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "AES",
		Mode:           "OFB",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       16,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec. Stream mode, padding not needed.",
	},
	"AES/OFB/PKCS5Padding": {
		Name:           "AES/OFB/PKCS5Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "AES",
		Mode:           "OFB",
		Padding:        "PKCS5Padding",
		NeedsIV:        true,
		IVLength:       16,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec. Padding unnecessary for OFB.",
	},
	"AES/CTR/NoPadding": {
		Name:           "AES/CTR/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "AES",
		Mode:           "CTR",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       16,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec. Stream mode, no padding needed.",
	},
	"AES/GCM/NoPadding": {
		Name:           "AES/GCM/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "AES",
		Mode:           "GCM",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       12, // 96-bit recommended
		NeedsTagLength: true,
		TagLength:      128, // 128-bit typical
		IsAEAD:         true,
		Notes:          "Uses GCMParameterSpec(tagLength, iv). Recommended AEAD mode.",
	},
	"AES/CCM/NoPadding": {
		Name:           "AES/CCM/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "AES",
		Mode:           "CCM",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       0, // Variable: 7-13 bytes
		NeedsTagLength: true,
		TagLength:      0, // Variable: 32-128 bits
		IsAEAD:         true,
		Notes:          "Uses CCMParameterSpec. Limited support, nonce length varies.",
	},

	// ============================================================================
	// CATEGORY 2: ALGORITHM/MODE/PADDING TRIPLES - DES
	// ============================================================================

	"DES/ECB/NoPadding": {
		Name:           "DES/ECB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "DES",
		Mode:           "ECB",
		Padding:        "NoPadding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Deprecated. Insecure: 56-bit key, ECB mode. Input multiple of 8 bytes.",
	},
	"DES/ECB/PKCS5Padding": {
		Name:           "DES/ECB/PKCS5Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "DES",
		Mode:           "ECB",
		Padding:        "PKCS5Padding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Deprecated. Insecure: 56-bit key, ECB mode.",
	},
	"DES/CBC/NoPadding": {
		Name:           "DES/CBC/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "DES",
		Mode:           "CBC",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Deprecated. Uses IvParameterSpec. Input multiple of 8 bytes.",
	},
	"DES/CBC/PKCS5Padding": {
		Name:           "DES/CBC/PKCS5Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "DES",
		Mode:           "CBC",
		Padding:        "PKCS5Padding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Deprecated. Uses IvParameterSpec.",
	},
	"DES/CFB/NoPadding": {
		Name:           "DES/CFB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "DES",
		Mode:           "CFB",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Deprecated. Uses IvParameterSpec.",
	},
	"DES/OFB/NoPadding": {
		Name:           "DES/OFB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "DES",
		Mode:           "OFB",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Deprecated. Uses IvParameterSpec.",
	},
	"DES/CTR/NoPadding": {
		Name:           "DES/CTR/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "DES",
		Mode:           "CTR",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Deprecated. Uses IvParameterSpec.",
	},

	// ============================================================================
	// CATEGORY 2: ALGORITHM/MODE/PADDING TRIPLES - DESede (Triple DES)
	// ============================================================================

	"DESede/ECB/NoPadding": {
		Name:           "DESede/ECB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "DESede",
		Mode:           "ECB",
		Padding:        "NoPadding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Triple DES. Insecure: ECB mode. Input multiple of 8 bytes.",
	},
	"DESede/ECB/PKCS5Padding": {
		Name:           "DESede/ECB/PKCS5Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "DESede",
		Mode:           "ECB",
		Padding:        "PKCS5Padding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Triple DES. Insecure: ECB mode.",
	},
	"DESede/CBC/NoPadding": {
		Name:           "DESede/CBC/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "DESede",
		Mode:           "CBC",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Triple DES. Uses IvParameterSpec. Input multiple of 8 bytes.",
	},
	"DESede/CBC/PKCS5Padding": {
		Name:           "DESede/CBC/PKCS5Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "DESede",
		Mode:           "CBC",
		Padding:        "PKCS5Padding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Triple DES. Uses IvParameterSpec.",
	},
	"DESede/CFB/NoPadding": {
		Name:           "DESede/CFB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "DESede",
		Mode:           "CFB",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Triple DES. Uses IvParameterSpec.",
	},
	"DESede/OFB/NoPadding": {
		Name:           "DESede/OFB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "DESede",
		Mode:           "OFB",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Triple DES. Uses IvParameterSpec.",
	},
	"DESede/CTR/NoPadding": {
		Name:           "DESede/CTR/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "DESede",
		Mode:           "CTR",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Triple DES. Uses IvParameterSpec.",
	},

	// ============================================================================
	// CATEGORY 2: ALGORITHM/MODE/PADDING TRIPLES - Blowfish
	// ============================================================================

	"Blowfish/ECB/NoPadding": {
		Name:           "Blowfish/ECB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "Blowfish",
		Mode:           "ECB",
		Padding:        "NoPadding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Insecure: ECB mode. Input multiple of 8 bytes.",
	},
	"Blowfish/ECB/PKCS5Padding": {
		Name:           "Blowfish/ECB/PKCS5Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "Blowfish",
		Mode:           "ECB",
		Padding:        "PKCS5Padding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Insecure: ECB mode.",
	},
	"Blowfish/CBC/NoPadding": {
		Name:           "Blowfish/CBC/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "Blowfish",
		Mode:           "CBC",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec. Input multiple of 8 bytes.",
	},
	"Blowfish/CBC/PKCS5Padding": {
		Name:           "Blowfish/CBC/PKCS5Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "Blowfish",
		Mode:           "CBC",
		Padding:        "PKCS5Padding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec.",
	},
	"Blowfish/CFB/NoPadding": {
		Name:           "Blowfish/CFB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "Blowfish",
		Mode:           "CFB",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec.",
	},
	"Blowfish/OFB/NoPadding": {
		Name:           "Blowfish/OFB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "Blowfish",
		Mode:           "OFB",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec.",
	},
	"Blowfish/CTR/NoPadding": {
		Name:           "Blowfish/CTR/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "Blowfish",
		Mode:           "CTR",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec.",
	},

	// ============================================================================
	// CATEGORY 2: ALGORITHM/MODE/PADDING TRIPLES - RC2
	// ============================================================================

	"RC2/ECB/NoPadding": {
		Name:           "RC2/ECB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RC2",
		Mode:           "ECB",
		Padding:        "NoPadding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Insecure: ECB mode. Input multiple of 8 bytes.",
	},
	"RC2/ECB/PKCS5Padding": {
		Name:           "RC2/ECB/PKCS5Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RC2",
		Mode:           "ECB",
		Padding:        "PKCS5Padding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Insecure: ECB mode.",
	},
	"RC2/CBC/NoPadding": {
		Name:           "RC2/CBC/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RC2",
		Mode:           "CBC",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec. Input multiple of 8 bytes.",
	},
	"RC2/CBC/PKCS5Padding": {
		Name:           "RC2/CBC/PKCS5Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RC2",
		Mode:           "CBC",
		Padding:        "PKCS5Padding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec.",
	},
	"RC2/CFB/NoPadding": {
		Name:           "RC2/CFB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RC2",
		Mode:           "CFB",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec.",
	},
	"RC2/OFB/NoPadding": {
		Name:           "RC2/OFB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RC2",
		Mode:           "OFB",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec.",
	},
	"RC2/CTR/NoPadding": {
		Name:           "RC2/CTR/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RC2",
		Mode:           "CTR",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec.",
	},

	// ============================================================================
	// CATEGORY 2: ALGORITHM/MODE/PADDING TRIPLES - RC5
	// ============================================================================

	"RC5/ECB/NoPadding": {
		Name:           "RC5/ECB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RC5",
		Mode:           "ECB",
		Padding:        "NoPadding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Insecure: ECB mode. Input multiple of 8 bytes.",
	},
	"RC5/ECB/PKCS5Padding": {
		Name:           "RC5/ECB/PKCS5Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RC5",
		Mode:           "ECB",
		Padding:        "PKCS5Padding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Insecure: ECB mode.",
	},
	"RC5/CBC/NoPadding": {
		Name:           "RC5/CBC/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RC5",
		Mode:           "CBC",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec. Input multiple of 8 bytes.",
	},
	"RC5/CBC/PKCS5Padding": {
		Name:           "RC5/CBC/PKCS5Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RC5",
		Mode:           "CBC",
		Padding:        "PKCS5Padding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec.",
	},
	"RC5/CFB/NoPadding": {
		Name:           "RC5/CFB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RC5",
		Mode:           "CFB",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec.",
	},
	"RC5/OFB/NoPadding": {
		Name:           "RC5/OFB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RC5",
		Mode:           "OFB",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec.",
	},
	"RC5/CTR/NoPadding": {
		Name:           "RC5/CTR/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RC5",
		Mode:           "CTR",
		Padding:        "NoPadding",
		NeedsIV:        true,
		IVLength:       8,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Uses IvParameterSpec.",
	},

	// ============================================================================
	// CATEGORY 2: ALGORITHM/MODE/PADDING TRIPLES - RSA
	// ============================================================================

	"RSA/ECB/NoPadding": {
		Name:           "RSA/ECB/NoPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RSA",
		Mode:           "ECB",
		Padding:        "NoPadding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Asymmetric cipher. Raw RSA without padding - insecure for most uses.",
	},
	"RSA/ECB/PKCS1Padding": {
		Name:           "RSA/ECB/PKCS1Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RSA",
		Mode:           "ECB",
		Padding:        "PKCS1Padding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Asymmetric cipher. PKCS#1 v1.5 padding. Vulnerable to padding oracle attacks.",
	},
	"RSA/ECB/OAEPPadding": {
		Name:           "RSA/ECB/OAEPPadding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RSA",
		Mode:           "ECB",
		Padding:        "OAEPPadding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Asymmetric cipher. OAEP with SHA-1 and MGF1. May use OAEPParameterSpec.",
	},
	"RSA/ECB/OAEPWithSHA-1AndMGF1Padding": {
		Name:           "RSA/ECB/OAEPWithSHA-1AndMGF1Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RSA",
		Mode:           "ECB",
		Padding:        "OAEPWithSHA-1AndMGF1Padding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Asymmetric cipher. OAEP explicitly with SHA-1.",
	},
	"RSA/ECB/OAEPWithSHA-256AndMGF1Padding": {
		Name:           "RSA/ECB/OAEPWithSHA-256AndMGF1Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RSA",
		Mode:           "ECB",
		Padding:        "OAEPWithSHA-256AndMGF1Padding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Asymmetric cipher. OAEP with SHA-256. Recommended for RSA encryption.",
	},
	"RSA/ECB/OAEPWithSHA-384AndMGF1Padding": {
		Name:           "RSA/ECB/OAEPWithSHA-384AndMGF1Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RSA",
		Mode:           "ECB",
		Padding:        "OAEPWithSHA-384AndMGF1Padding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Asymmetric cipher. OAEP with SHA-384.",
	},
	"RSA/ECB/OAEPWithSHA-512AndMGF1Padding": {
		Name:           "RSA/ECB/OAEPWithSHA-512AndMGF1Padding",
		Enabled:        true,
		Category:       CategoryTriple,
		Algorithm:      "RSA",
		Mode:           "ECB",
		Padding:        "OAEPWithSHA-512AndMGF1Padding",
		NeedsIV:        false,
		NeedsTagLength: false,
		IsAEAD:         false,
		Notes:          "Asymmetric cipher. OAEP with SHA-512.",
	},
}

// ValidateCipherTransformation checks if a transformation string is valid
func ValidateCipherTransformation(transformation string) (CipherTransformation, bool) {
	config, exists := CipherConfigTable[transformation]
	if !exists {
		return CipherTransformation{}, false
	}
	return config, config.Enabled
}

// GetRequiredParameters returns what parameters are needed for initialization
func (ct CipherTransformation) GetRequiredParameters() []string {
	params := []string{}

	if ct.KeyDerivation {
		params = append(params, "password")
	} else {
		params = append(params, "key")
	}

	if ct.NeedsIV {
		params = append(params, "iv")
	}

	if ct.NeedsTagLength {
		params = append(params, "tagLength")
	}

	if ct.NeedsSalt {
		params = append(params, "salt")
	}

	if ct.NeedsIterations {
		params = append(params, "iterations")
	}

	return params
}

// GetEnabledTransformations returns all enabled cipher transformations
func GetEnabledTransformations() map[string]CipherTransformation {
	enabled := make(map[string]CipherTransformation)
	for name, config := range CipherConfigTable {
		if config.Enabled {
			enabled[name] = config
		}
	}
	return enabled
}

// DisableTransformation marks a cipher transformation as disabled
func DisableTransformation(name string) bool {
	if config, exists := CipherConfigTable[name]; exists {
		config.Enabled = false
		CipherConfigTable[name] = config
		return true
	}
	return false
}

// EnableTransformation marks a cipher transformation as enabled
func EnableTransformation(name string) bool {
	if config, exists := CipherConfigTable[name]; exists {
		config.Enabled = true
		CipherConfigTable[name] = config
		return true
	}
	return false
}
