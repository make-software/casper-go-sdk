package ed25519

import (
	"crypto/ed25519"
	"encoding/pem"
	"errors"
	"os"
)

const PemFramePublicKeyPrefixSize = 12

type PublicKey ed25519.PublicKey

func (v PublicKey) Bytes() []byte {
	return []byte(v)
}

func (v PublicKey) VerifySignature(message []byte, sig []byte) bool {
	return ed25519.Verify(ed25519.PublicKey(v), message, sig)
}

func NewPublicKey(data []byte) PublicKey {
	return PublicKey(data)
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

	return PublicKey(data), nil
}
