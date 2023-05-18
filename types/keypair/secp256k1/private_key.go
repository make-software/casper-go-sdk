package secp256k1

import (
	"errors"
	"os"

	"github.com/btcsuite/btcd/btcec"
)

type PrivateKey struct {
	key *btcec.PrivateKey
}

func (v PrivateKey) PublicKeyBytes() []byte {
	return v.key.PubKey().SerializeCompressed()
}

func (v PrivateKey) Sign(mes []byte) ([]byte, error) {
	val, err := v.key.Sign(mes)
	if err != nil {
		return nil, err
	}
	return val.Serialize(), nil
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
