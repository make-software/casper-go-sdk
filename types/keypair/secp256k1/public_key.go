package secp256k1

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
)

const PublicKeySize = 33

// used to reject malleable signatures
// see:
//   - https://github.com/ethereum/go-ethereum/blob/f9401ae011ddf7f8d2d95020b7446c17f8d98dc1/crypto/signature_nocgo.go#L90-L93
//   - https://github.com/ethereum/go-ethereum/blob/f9401ae011ddf7f8d2d95020b7446c17f8d98dc1/crypto/crypto.go#L39
var secp256k1halfN = new(big.Int).Rsh(btcec.S256().N, 1)

type PublicKey btcec.PublicKey

func (v PublicKey) Bytes() []byte {
	key := btcec.PublicKey(v)
	return key.SerializeCompressed()
}

// VerifySignature verifies a signature of the form R || S.
// It rejects signatures which are not in lower-S form.
func (v PublicKey) VerifySignature(msg []byte, sigStr []byte) bool {
	if len(sigStr) != 64 {
		return false
	}

	pub, err := btcec.ParsePubKey(v.Bytes(), btcec.S256())
	if err != nil {
		return false
	}

	// parse the signature:
	signature := signatureFromBytes(sigStr)
	// Reject malleable signatures. libsecp256k1 does this check but btcec doesn't.
	// see: https://github.com/ethereum/go-ethereum/blob/f9401ae011ddf7f8d2d95020b7446c17f8d98dc1/crypto/signature_nocgo.go#L90-L93
	if signature.S.Cmp(secp256k1halfN) > 0 {
		return false
	}

	sum256 := sha256.Sum256(msg)
	return signature.Verify(sum256[:], pub)
}

// Read Signature struct from R || S. Caller needs to ensure
// that len(sigStr) == 64.
func signatureFromBytes(sigStr []byte) *btcec.Signature {
	return &btcec.Signature{
		R: new(big.Int).SetBytes(sigStr[:32]),
		S: new(big.Int).SetBytes(sigStr[32:64]),
	}
}

func NewPublicKey(data []byte) (PublicKey, error) {
	if len(data) != PublicKeySize {
		return PublicKey{}, fmt.Errorf("can't parse wrong size of public key: %d", len(data))
	}
	key, err := btcec.ParsePubKey(data, btcec.S256())
	if err != nil {
		return PublicKey{}, err
	}
	return PublicKey(*key), err
}
