package ed25519

import (
	"crypto/ed25519"
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

func (k PrivateKey) PublicKeyBytes() []byte {
	return k.key.Public().(ed25519.PublicKey)
}

func (k PrivateKey) Sign(mes []byte) ([]byte, error) {
	return ed25519.Sign(k.key, mes), nil
}

func (k PrivateKey) ToPem() ([]byte, error) {
	return PrivateKeyToPem(k.key)
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
