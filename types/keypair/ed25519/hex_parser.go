package ed25519

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
)

// NewPrivateKeyFromBytes creates an ED25519 PrivateKey from raw bytes
func NewPrivateKeyFromBytes(key []byte) (PrivateKey, error) {
	// Check if the key size matches the expected ED25519 private key size
	if len(key) != ed25519.PrivateKeySize {
		return PrivateKey{}, fmt.Errorf("wrong key size: expected %v bytes, got %v bytes", ed25519.PrivateKeySize, len(key))
	}

	privateKey := ed25519.PrivateKey(key)
	return PrivateKey{
		key: privateKey,
	}, nil
}

// NewPrivateKeyFromHex creates an ED25519 PrivateKey from a hex string
func NewPrivateKeyFromHex(key string) (PrivateKey, error) {
	// Validate hex string length (128 hex characters = 64 bytes for ED25519 private key)
	if len(key) != ed25519.PrivateKeySize*2 {
		return PrivateKey{}, fmt.Errorf("invalid hex string length: expected %v characters, got %v", ed25519.PrivateKeySize*2, len(key))
	}

	b, err := hex.DecodeString(key)
	if err != nil {
		return PrivateKey{}, fmt.Errorf("failed to decode hex string: %v", err)
	}

	return NewPrivateKeyFromBytes(b)
}
