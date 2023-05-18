package secp256k1

import (
	"math/big"

	"github.com/btcsuite/btcd/btcec"
)

type PublicKey btcec.PublicKey

func (v PublicKey) Bytes() []byte {
	key := btcec.PublicKey(v)
	return key.SerializeCompressed()
}

func (v PublicKey) VerifySignature(message []byte, sig []byte) bool {
	signature, err := btcec.ParseSignature(sig, btcec.S256())
	if err != nil {
		return false
	}

	key := btcec.PublicKey(v)
	verify := signature.Verify(message, &key)
	return verify
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
