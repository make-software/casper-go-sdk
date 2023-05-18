package secp256k1

import (
	"crypto/sha256"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
)

type PublicKey btcec.PublicKey

func (v PublicKey) Bytes() []byte {
	key := btcec.PublicKey(v)
	return key.SerializeCompressed()
}

func (v PublicKey) VerifySignature(message []byte, sig []byte) bool {
	if len(sig) != 64 {
		return false
	}

	signature := signatureFromBytes(sig)

	key := btcec.PublicKey(v)
	sum256 := sha256.Sum256(message)
	return signature.Verify(sum256[:], &key)
}

func NewPublicKey(data []byte) (PublicKey, error) {
	key, err := btcec.ParsePubKey(data, btcec.S256())
	if err != nil {
		return PublicKey{}, err
	}
	return PublicKey(*key), err
}

// Read Signature struct from R || S. Caller needs to ensure
// that len(sigStr) == 64.
func signatureFromBytes(sigStr []byte) *btcec.Signature {
	return &btcec.Signature{
		R: new(big.Int).SetBytes(sigStr[:32]),
		S: new(big.Int).SetBytes(sigStr[32:64]),
	}
}
