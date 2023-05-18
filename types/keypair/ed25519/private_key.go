package ed25519

import (
	"crypto/ed25519"
	"encoding/pem"
	"errors"
	"os"
)

const PemFramePrivateKeyPrefixSize = 16

type PrivateKey struct {
	key ed25519.PrivateKey
}

func (k PrivateKey) PublicKeyBytes() []byte {
	return k.key.Public().(ed25519.PublicKey)
}

func (k PrivateKey) Sign(mes []byte) ([]byte, error) {
	return ed25519.Sign(k.key, mes), nil
}

func GeneratePrivateKey() (PrivateKey, error) {
	_, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}
	return PrivateKey{
		key: priv,
	}, nil
}

func NewPrivateKeyFromPEMFile(path string) (PrivateKey, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return PrivateKey{}, errors.New("can't read file")
	}

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
