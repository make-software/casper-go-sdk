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

// NewPrivateKeyFromSeedBytes creates an ED25519 PrivateKey from raw seed bytes
func NewPrivateKeyFromSeedBytes(seed []byte) (PrivateKey, error) {
	// Check if the seed size matches the expected ED25519 seed size
	if len(seed) != ed25519.SeedSize {
		return PrivateKey{}, fmt.Errorf("wrong seed size: expected %v bytes, got %v bytes", ed25519.SeedSize, len(seed))
	}

	privateKey := ed25519.NewKeyFromSeed(seed)
	return PrivateKey{
		key: privateKey,
	}, nil
}

// NewPrivateKeyFromSeedHex creates an ED25519 PrivateKey from a hex string seed
func NewPrivateKeyFromSeedHex(seed string) (PrivateKey, error) {
	// Validate hex string length (64 hex characters = 32 bytes for ED25519 seed)
	if len(seed) != ed25519.SeedSize*2 {
		return PrivateKey{}, fmt.Errorf("invalid hex string length: expected %v characters, got %v", ed25519.SeedSize*2, len(seed))
	}

	b, err := hex.DecodeString(seed)
	if err != nil {
		return PrivateKey{}, fmt.Errorf("failed to decode hex string: %v", err)
	}

	return NewPrivateKeyFromSeedBytes(b)
}
