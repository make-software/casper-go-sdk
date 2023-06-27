package ed25519

import (
	"crypto/ed25519"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

const PemFramePublicKeyPrefixSize = 12
const PublicKeySize = 32

type PublicKey ed25519.PublicKey

func (v PublicKey) Bytes() []byte {
	return []byte(v)
}

func (v PublicKey) VerifySignature(message []byte, sig []byte) bool {
	return ed25519.Verify(ed25519.PublicKey(v), message, sig)
}

func NewPublicKey(data []byte) (PublicKey, error) {
	if len(data) != PublicKeySize {
		return nil, fmt.Errorf("can't parse wrong size of public key: %d", len(data))
	}
	return PublicKey(data), nil
}

func ParsePublicKeyFile(path string) (PublicKey, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("can't read file")
	}

	block, _ := pem.Decode(content)
	if err != nil {
		return PublicKey{}, err
	}
	// Trim PEM prefix
	data := block.Bytes[PemFramePublicKeyPrefixSize:]

	return NewPublicKey(data)
}
