package secp256k1

import (
	"encoding/hex"
	"fmt"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// NewPrivateKeyFromBytes creates a secp256k1 PrivateKey from raw bytes
func NewPrivateKeyFromBytes(key []byte) (PrivateKey, error) {
	// Validate key length (32 bytes for secp256k1)
	if len(key) != 32 {
		return PrivateKey{}, fmt.Errorf("invalid private key length: expected 32 bytes, got %d", len(key))
	}

	privateKey := secp256k1.PrivKeyFromBytes(key)

	return PrivateKey{
		key: privateKey,
	}, nil
}

// NewPrivateKeyFromHex creates a secp256k1 PrivateKey from a hex string
func NewPrivateKeyFromHex(key string) (PrivateKey, error) {
	// Validate hex string length (64 hex characters for secp256k1)
	if len(key) != 64 {
		return PrivateKey{}, fmt.Errorf("invalid private key hex length: expected 64 characters, got %d", len(key))
	}

	// Decode the hex string into bytes
	b, err := hex.DecodeString(key)
	if err != nil {
		return PrivateKey{}, fmt.Errorf("failed to decode hex: %v", err)
	}

	return NewPrivateKeyFromBytes(b)
}
