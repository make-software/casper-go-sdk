package secp256k1

import (
	"crypto/sha256"
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
	hash := sha256.Sum256(mes)
	sig, err := v.key.Sign(hash[:])
	if err != nil {
		return nil, err
	}

	return serializeSig(sig), nil
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

// Serialize signature to R || S.
// R, S are padded to 32 bytes respectively.
func serializeSig(sig *btcec.Signature) []byte {
	rBytes := sig.R.Bytes()
	sBytes := sig.S.Bytes()
	sigBytes := make([]byte, 64)
	// 0 pad the byte arrays from the left if they aren't big enough.
	copy(sigBytes[32-len(rBytes):32], rBytes)
	copy(sigBytes[64-len(sBytes):64], sBytes)
	return sigBytes
}
