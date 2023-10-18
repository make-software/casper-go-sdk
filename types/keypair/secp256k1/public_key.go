package secp256k1

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
)

const PublicKeySize = 33

type PublicKey struct {
	key *secp256k1.PublicKey
}

func (v PublicKey) Bytes() []byte {
	return v.key.SerializeCompressed()
}

// VerifySignature verifies a signature of the form R || S.
// It rejects signatures which are not in lower-S form.
func (v PublicKey) VerifySignature(msg []byte, sigStr []byte) bool {
	var signature *ecdsa.Signature
	var err error
	// if old signature len = 64, parse it as raw signature
	if len(sigStr) == 64 {
		// Split the signature bytes into r and s values and parse them into ModNScalar
		var r, s secp256k1.ModNScalar
		var bytesR [32]byte
		var bytesS [32]byte
		copy(bytesR[:], sigStr[:32])
		copy(bytesS[:], sigStr[32:])
		r.SetBytes(&bytesR)
		s.SetBytes(&bytesS)

		signature = ecdsa.NewSignature(&r, &s)
	} else {
		signature, err = ecdsa.ParseDERSignature(sigStr)
		if err != nil {
			log.Println(err)
			return false
		}
	}
	hash := sha256.Sum256(msg)
	return signature.Verify(hash[:], v.key)
}

func NewPublicKey(data []byte) (PublicKey, error) {
	if len(data) != PublicKeySize {
		return PublicKey{}, fmt.Errorf("can't parse wrong size of public key: %d", len(data))
	}
	key, err := secp256k1.ParsePubKey(data)
	if err != nil {
		return PublicKey{}, err
	}
	return PublicKey{key: key}, err
}
