package ed25519

import (
	"crypto/ed25519"
	"encoding/pem"
)

const PemFramePrivateKeyPrefixSize = 16

func PrivateKeyToPem(priv ed25519.PrivateKey) ([]byte, error) {
	// Extract the seed (first 32 bytes of the private key)
	seed := priv[:32]

	// Create the prefix for the seed
	prefix := make([]byte, PemFramePrivateKeyPrefixSize)

	// Combine the prefix and the seed to match the original format
	fullKey := append(prefix, seed...)

	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: fullKey,
		},
	), nil
}

func NewPrivateKeyFromPEM(content []byte) (PrivateKey, error) {
	block, _ := pem.Decode(content)
	// Trim PEM prefix
	data := block.Bytes[PemFramePrivateKeyPrefixSize:]
	// Use only last 32 bytes of private Key
	privateKeySize := len(data)
	if privateKeySize > 32 {
		data = data[privateKeySize%32:]
	}

	privateEdDSA := ed25519.NewKeyFromSeed(data)

	return PrivateKey{
		key: privateEdDSA,
	}, nil
}
