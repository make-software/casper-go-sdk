package secp256k1

import (
	"crypto/sha256"
	"errors"
	"os"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
)

type PrivateKey struct {
	key *secp256k1.PrivateKey
}

func (v PrivateKey) PublicKeyBytes() []byte {
	return v.key.PubKey().SerializeCompressed()
}

func (v PrivateKey) Sign(mes []byte) ([]byte, error) {
	hash := sha256.Sum256(mes)
	// Return the signature as a concatenation of the R and S values in big-endian to match the old signature format.
	return ecdsa.SignCompact(v.key, hash[:], false)[1:], nil
}

func NewPrivateKeyFromPemFile(path string) (PrivateKey, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return PrivateKey{}, errors.New("can't read file")
	}
	return NewPrivateKeyFromPem(content)
}

func NewPrivateKeyFromPem(content []byte) (PrivateKey, error) {
	private, err := PemToPrivateKey(content)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		key: private,
	}, nil
}
